# TCP 发送模式接收回应重构方案（方案B）

## 1. 方案概述

### 1.1 目标
- 统一单次发送、路径发送、交互式三种模式的接收处理逻辑
- 复用交互式模式的成熟接收循环代码
- 支持持续接收直到超时或连接关闭
- 保持向后兼容（`-w` 参数行为不变）

### 1.2 核心思想
将 `receiveLoopInteractive` 改造为通用接收循环，通过回调函数处理不同模式的数据输出需求。

---

## 2. 当前代码结构分析

### 2.1 交互式模式（已实现）
```go
// interactive.go
func receiveLoopInteractive(session *InteractiveSession, recvChan chan<- []byte, errChan chan<- error, done <-chan bool) {
    buf := make([]byte, 4096)
    for {
        select {
        case <-done:
            return
        default:
        }

        if session.Config.Wait > 0 {
            _ = session.Conn.SetReadDeadline(time.Now().Add(session.Config.Wait))
        }

        n, err := session.Conn.Read(buf)
        if err != nil {
            errChan <- err
            return
        }
        if n > 0 {
            recvChan <- buf[:n]
        }
    }
}
```

### 2.2 单次/路径发送模式（需要改造）
```go
// cmd_tcp.go - runDataMode
// 当前只读取一次，需要改为循环读取
if config.Wait > 0 {
    _ = conn.SetReadDeadline(time.Now().Add(config.Wait))
    buffer := make([]byte, 4096)
    n, err := conn.Read(buffer)  // ← 只读一次
    // ...
}
```

---

## 3. 重构方案设计

### 3.1 文件结构调整

```
internal/commands/tcp/
├── cmd_tcp.go           # 主命令入口（保持不变）
├── interactive.go       # 交互式模式（改造）
├── path.go             # 路径处理（保持不变）
├── output.go           # 输出处理（新增/改造）
└── receiver.go         # 通用接收模块（新增） ← 核心
```

### 3.2 新增文件：receiver.go

```go
// Package tcp 实现了 TCP 连接测试和端口扫描功能
package tcp

import (
	"net"
	"time"
)

// ReceiveHandler 接收数据处理回调函数
// 参数:
//   - data: 接收到的原始数据
//   - seq: 数据包序号（从1开始）
// 返回值:
//   - bool: 是否继续接收，返回 false 则停止接收循环
//   - error: 处理错误，非 nil 则停止接收循环

type ReceiveHandler func(data []byte, seq int) (bool, error)

// ReceiveOptions 接收选项
type ReceiveOptions struct {
	// 总超时时间，0 表示不限制（但会使用默认小超时防止阻塞）
	Timeout time.Duration

	// 单次读取超时，默认 100ms
	ReadTimeout time.Duration

	// 缓冲区大小，默认 4096
	BufferSize int

	// 最大接收数据量，0 表示不限制
	MaxBytes int64

	// 是否以十六进制格式处理数据
	HexMode bool
}

// DefaultReceiveOptions 返回默认接收选项
func DefaultReceiveOptions() *ReceiveOptions {
	return &ReceiveOptions{
		Timeout:     0,
		ReadTimeout: 100 * time.Millisecond,
		BufferSize:  4096,
		MaxBytes:    0,
		HexMode:     false,
	}
}

// ReceiveLoop 通用接收循环
//
// 参数:
//   - conn: TCP 连接
//   - opts: 接收选项
//   - handler: 数据处理回调函数
//
// 返回值:
//   - []byte: 所有接收数据的拼接结果（如果 handler 返回 nil）
//   - int: 接收数据包数量
//   - error: 接收过程中的错误
//
// 功能特性:
//   - 支持超时控制（总超时和单次读取超时）
//   - 支持数据量限制
//   - 支持通过 handler 控制是否继续接收
//   - 自动处理超时错误，区分临时超时和连接错误
func ReceiveLoop(conn net.Conn, opts *ReceiveOptions, handler ReceiveHandler) ([]byte, int, error) {
	if opts == nil {
		opts = DefaultReceiveOptions()
	}

	buf := make([]byte, opts.BufferSize)
	var result []byte
	var totalBytes int64
	seq := 0
	startTime := time.Now()

	for {
		// 检查总超时
		if opts.Timeout > 0 && time.Since(startTime) >= opts.Timeout {
			return result, seq, nil // 正常结束，非错误
		}

		// 设置单次读取超时
		readTimeout := opts.ReadTimeout
		if opts.Timeout > 0 {
			remaining := opts.Timeout - time.Since(startTime)
			if remaining <= 0 {
				return result, seq, nil
			}
			if remaining < readTimeout {
				readTimeout = remaining
			}
		}

		_ = conn.SetReadDeadline(time.Now().Add(readTimeout))
		n, err := conn.Read(buf)

		if n > 0 {
			seq++
			data := buf[:n]
			totalBytes += int64(n)

			// 检查数据量限制
			if opts.MaxBytes > 0 && totalBytes > opts.MaxBytes {
				return result, seq, fmt.Errorf("max bytes limit exceeded: %d", opts.MaxBytes)
			}

			// 累积数据
			if result == nil {
				result = make([]byte, 0, opts.BufferSize*2)
			}
			result = append(result, data...)

			// 调用处理回调
			if handler != nil {
				continueRecv, handlerErr := handler(data, seq)
				if handlerErr != nil {
					return result, seq, handlerErr
				}
				if !continueRecv {
					return result, seq, nil
				}
			}
		}

		if err != nil {
			// 区分超时错误和连接错误
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// 临时超时，继续尝试
				continue
			}
			// 连接关闭或其他错误
			if err.Error() == "EOF" {
				return result, seq, nil // 正常结束
			}
			return result, seq, err
		}
	}
}

// ReceiveAll 简单接收所有数据（无回调版本）
//
// 参数:
//   - conn: TCP 连接
//   - timeout: 总超时时间
//
// 返回值:
//   - []byte: 所有接收的数据
//   - error: 接收错误
func ReceiveAll(conn net.Conn, timeout time.Duration) ([]byte, error) {
	opts := DefaultReceiveOptions()
	opts.Timeout = timeout
	data, _, err := ReceiveLoop(conn, opts, nil)
	return data, err
}

// ReceiveWithPrint 接收并打印数据（用于单次/路径发送模式）
//
// 参数:
//   - conn: TCP 连接
//   - timeout: 总超时时间
//   - quiet: 是否静默模式
//   - hexMode: 是否十六进制显示
//
// 返回值:
//   - []byte: 所有接收的数据
//   - error: 接收错误
func ReceiveWithPrint(conn net.Conn, timeout time.Duration, quiet, hexMode bool) ([]byte, error) {
	if timeout <= 0 {
		return nil, nil
	}

	opts := DefaultReceiveOptions()
	opts.Timeout = timeout
	opts.HexMode = hexMode

	var allData []byte
	handler := func(data []byte, seq int) (bool, error) {
		if !quiet {
			if hexMode {
				fmt.Printf("[←] Received packet #%d (%d bytes):\n%s\n", seq, len(data), hexDump(data))
			} else {
				fmt.Printf("[←] Received packet #%d (%d bytes):\n%s\n", seq, len(data), string(data))
			}
		}
		return true, nil // 继续接收
	}

	return ReceiveLoop(conn, opts, handler)
}

// hexDump 将数据转换为十六进制字符串
func hexDump(data []byte) string {
	var result strings.Builder
	for i := 0; i < len(data); i += 16 {
		// 地址偏移
		result.WriteString(fmt.Sprintf("%04x  ", i))

		// 十六进制部分
		for j := 0; j < 16; j++ {
			if i+j < len(data) {
				result.WriteString(fmt.Sprintf("%02x ", data[i+j]))
			} else {
				result.WriteString("   ")
			}
			if j == 7 {
				result.WriteString(" ")
			}
		}

		// ASCII 部分
		result.WriteString(" |")
		for j := 0; j < 16 && i+j < len(data); j++ {
			b := data[i+j]
			if b >= 32 && b < 127 {
				result.WriteByte(b)
			} else {
				result.WriteByte('.')
			}
		}
		result.WriteString("|\n")
	}
	return result.String()
}
```

### 3.3 改造：interactive.go

将原来的 `receiveLoopInteractive` 改为使用通用 `ReceiveLoop`：

```go
// receiveLoopInteractive 交互式接收循环（改造后）
//
// 参数:
//   - session: 交互式会话
//   - recvChan: 接收数据通道
//   - errChan: 错误通道
//   - done: 退出信号通道
func receiveLoopInteractive(session *InteractiveSession, recvChan chan<- []byte, errChan chan<- error, done <-chan bool) {
	opts := DefaultReceiveOptions()
	opts.ReadTimeout = session.Config.Wait
	if opts.ReadTimeout <= 0 {
		opts.ReadTimeout = 100 * time.Millisecond
	}

	handler := func(data []byte, seq int) (bool, error) {
		select {
		case <-done:
			return false, nil
		case recvChan <- data:
			return true, nil
		}
	}

	_, _, err := ReceiveLoop(session.Conn, opts, handler)
	if err != nil {
		select {
		case errChan <- err:
		case <-done:
		}
	}
}
```

### 3.4 改造：cmd_tcp.go - runDataMode

```go
// runDataMode 运行数据发送模式（改造后）
func runDataMode(config TcpConfig, ipAddr net.IP) error {
	// ... 连接建立代码保持不变 ...

	// 发送数据...

	// 改造：使用通用接收函数
	var response []byte
	if config.Wait > 0 {
		if !config.Quiet && !config.Json {
			fmt.Println("\nWaiting for response...")
		}

		var err error
		response, err = ReceiveWithPrint(conn, config.Wait, config.Quiet, config.Hex)
		if err != nil && !config.Quiet {
			fmt.Printf("Receive error: %v\n", err)
		}
	}

	// ... JSON 输出等后续代码 ...
}
```

### 3.5 改造：cmd_tcp.go - runBannerMode

```go
// runBannerMode 运行 Banner 获取模式（改造后）
func runBannerMode(config TcpConfig, ipAddr net.IP) error {
	// ... 连接建立代码 ...

	// 改造：使用通用接收函数
	var banner string
	if config.Wait > 0 {
		data, _, err := ReceiveLoop(conn, &ReceiveOptions{
			Timeout:     config.Wait,
			ReadTimeout: 100 * time.Millisecond,
			BufferSize:  4096,
		}, nil) // 无回调，直接累积数据

		if err == nil && len(data) > 0 {
			banner = string(data)
		}
	}

	// ... 输出代码 ...
}
```

---

## 4. 输出处理模块（output.go）

统一处理各种输出格式：

```go
// output.go
package tcp

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// ResponseOutput 响应输出接口
type ResponseOutput interface {
	PrintReceive(data []byte, seq int, hexMode bool)
	PrintError(err error)
	PrintSummary()
}

// ConsoleOutput 控制台输出
type ConsoleOutput struct {
	Quiet   bool
	HexMode bool
}

func (c *ConsoleOutput) PrintReceive(data []byte, seq int, hexMode bool) {
	if c.Quiet {
		return
	}
	if hexMode {
		fmt.Printf("[←] Packet #%d (%d bytes):\n%s\n", seq, len(data), hexDump(data))
	} else {
		fmt.Printf("[←] Packet #%d (%d bytes):\n%s\n", seq, len(data), string(data))
	}
}

// JSONOutput JSON 格式输出
type JSONOutput struct {
	Data []string `json:"responses"`
}

func (j *JSONOutput) PrintReceive(data []byte, seq int, hexMode bool) {
	j.Data = append(j.Data, string(data))
}

func (j *JSONOutput) Flush() {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(j)
}
```

---

## 5. 使用示例

### 5.1 单次发送模式
```bash
# 发送并持续接收 5 秒
fck tcp -d "GET / HTTP/1.1\r\nHost: example.com\r\n\r\n" -w 5s example.com 80

# 输出示例：
# Connecting to example.com (93.184.216.34):80...
# Connected, sending 47 bytes...
# 
# Waiting for response...
# [←] Packet #1 (645 bytes):
# HTTP/1.1 200 OK
# ...
# [←] Packet #2 (1024 bytes):
# ...
```

### 5.2 路径发送模式
```bash
# 发送多个文件，最后统一接收响应
fck tcp -p ./requests/ -w 3s example.com 8080
```

### 5.3 交互式模式（保持不变）
```bash
fck tcp -I example.com 8080
```

---

## 6. 边缘情况处理

| 场景 | 处理方案 |
|------|----------|
| 服务端立即关闭连接 | `ReceiveLoop` 返回 EOF，正常结束 |
| 服务端不发送数据 | 等待超时后正常返回空数据 |
| 服务端分多次发送 | 循环读取直到超时，累积所有数据 |
| 数据量超过缓冲区 | 多次读取，自动扩容 result 切片 |
| 网络中断 | 返回网络错误，已接收数据一并返回 |
| 用户按 Ctrl+C | 通过 context 取消（交互式）或等待超时 |

---

## 7. 实现优先级

| 优先级 | 任务 | 说明 |
|--------|------|------|
| P0 | 创建 receiver.go | 核心模块 |
| P0 | 改造 runDataMode | 单次/路径发送模式 |
| P1 | 改造 runBannerMode | Banner 模式 |
| P1 | 改造 interactive.go | 复用通用循环 |
| P2 | 创建 output.go | 统一输出处理 |
| P2 | 添加十六进制显示 | `--hex` 参数支持 |

---

## 8. 向后兼容性

- `-w` 参数行为保持不变（总等待时间）
- 静默模式 `-q` 仍然有效
- JSON 输出格式保持不变
- 交互式模式用户无感知

---

**文档版本**: v1.0  
**创建日期**: 2026-01-24  
**作者**: AI Assistant
