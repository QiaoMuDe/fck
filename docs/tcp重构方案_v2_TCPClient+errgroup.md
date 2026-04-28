# TCP 模块重构方案 v2 - TCPClient + errgroup

## 1. 方案概述

使用 `TCPClient` 结构体封装 TCP 连接操作，结合 `errgroup` 管理并发 goroutine，实现三种发送模式的统一处理。

### 核心设计
- **TCPClient**: 封装连接管理和基础收发操作
- **errgroup**: 管理接收/发送 goroutine 的生命周期和错误处理
- **三种模式**: 基于 TCPClient 实现各自的发送接收策略

---

## 2. 文件结构

```
internal/commands/tcp/
├── client.go          # TCPClient 结构体（新增）
├── modes.go           # 三种模式实现（新增）
├── interactive.go     # 交互式模式（改造）
├── cmd_tcp.go         # 主入口（简化）
├── listen.go          # 监听模式（保持不变）
├── scan.go            # 扫描模式（保持不变）
├── banner.go          # Banner 模式（保持不变）
├── path.go            # 路径解析（保持不变）
└── utils.go           # 工具函数（保持不变）
```

---

## 3. 详细设计

### 3.1 client.go - TCPClient 结构体

```go
package tcp

import (
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"
)

// TCPClient TCP 客户端封装
type TCPClient struct {
	conn        net.Conn
	remoteAddr  string
	config      TcpConfig
	
	// 统计
	bytesSent   int64
	bytesRecv   int64
	sendSeq     int32
	recvSeq     int32
	
	// 通道
	recvChan    chan []byte
	sendChan    chan []byte
	errChan     chan error
	
	// 控制
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
}

// NewTCPClient 创建 TCP 客户端
func NewTCPClient(config TcpConfig) *TCPClient {
	ctx, cancel := context.WithCancel(context.Background())
	return &TCPClient{
		config:   config,
		recvChan: make(chan []byte, 100),
		sendChan: make(chan []byte, 100),
		errChan:  make(chan error, 1),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Connect 建立连接
func (c *TCPClient) Connect(host string, port int) error {
	target := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", target, c.config.Timeout)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	c.conn = conn
	c.remoteAddr = conn.RemoteAddr().String()
	return nil
}

// Close 关闭连接
func (c *TCPClient) Close() error {
	c.cancel()
	if c.conn != nil {
		_ = c.conn.Close()
	}
	close(c.recvChan)
	close(c.sendChan)
	c.wg.Wait()
	return nil
}

// Send 发送数据（同步）
func (c *TCPClient) Send(data []byte) error {
	_, err := c.conn.Write(data)
	if err != nil {
		return err
	}
	atomic.AddInt64(&c.bytesSent, int64(len(data)))
	atomic.AddInt32(&c.sendSeq, 1)
	return nil
}

// SendString 发送字符串
func (c *TCPClient) SendString(s string) error {
	return c.Send([]byte(s))
}

// StartReceive 启动接收循环（在 errgroup 中运行）
func (c *TCPClient) StartReceive(g *errgroup.Group, timeout time.Duration) {
	g.Go(func() error {
		return c.receiveLoop(timeout)
	})
}

// receiveLoop 接收循环
func (c *TCPClient) receiveLoop(timeout time.Duration) error {
	buf := make([]byte, 4096)
	
	// 默认超时
	if timeout <= 0 {
		timeout = 100 * time.Millisecond
	}
	
	for {
		select {
		case <-c.ctx.Done():
			return nil
		default:
		}
		
		_ = c.conn.SetReadDeadline(time.Now().Add(timeout))
		n, err := c.conn.Read(buf)
		
		if n > 0 {
			data := make([]byte, n)
			copy(data, buf[:n])
			atomic.AddInt64(&c.bytesRecv, int64(n))
			atomic.AddInt32(&c.recvSeq, 1)
			
			select {
			case c.recvChan <- data:
			case <-c.ctx.Done():
				return nil
			}
		}
		
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			if err.Error() == "EOF" {
				return nil
			}
			return err
		}
	}
}

// GetRecvChan 获取接收通道
func (c *TCPClient) GetRecvChan() <-chan []byte {
	return c.recvChan
}

// GetStats 获取统计信息
func (c *TCPClient) GetStats() (sent, recv int64, sendSeq, recvSeq int32) {
	return atomic.LoadInt64(&c.bytesSent),
		atomic.LoadInt64(&c.bytesRecv),
		atomic.LoadInt32(&c.sendSeq),
		atomic.LoadInt32(&c.recvSeq)
}
```

### 3.2 modes.go - 三种发送模式

```go
package tcp

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/sync/errgroup"
)

// RunStringMode 单次字符串发送模式
// 发送数据后，持续接收响应直到超时
func RunStringMode(config TcpConfig, ipAddr net.IP) error {
	client := NewTCPClient(config)
	if err := client.Connect(ipAddr.String(), config.Port); err != nil {
		return err
	}
	defer client.Close()

	// 发送数据
	data := []byte(config.Data)
	if err := client.Send(data); err != nil {
		return fmt.Errorf("send failed: %w", err)
	}

	if !config.Quiet && !config.Json {
		fmt.Printf("[→] Sent %d bytes\n", len(data))
	}

	// 如果没有设置等待时间，直接返回
	if config.Wait <= 0 {
		return nil
	}

	// 使用 errgroup 管理接收
	g, ctx := errgroup.WithContext(context.Background())
	client.StartReceive(g, config.Wait)

	// 主循环：处理接收数据
	var allData []byte
	done := make(chan struct{})
	
	go func() {
		for data := range client.GetRecvChan() {
			allData = append(allData, data...)
			if !config.Quiet && !config.Json {
				if config.Hex {
					fmt.Printf("[←] Received %d bytes:\n%s\n", len(data), hexDump(data))
				} else {
					fmt.Printf("[←] Received %d bytes:\n%s\n", len(data), string(data))
				}
			}
		}
		close(done)
	}()

	// 等待超时或错误
	select {
	case <-done:
	case <-time.After(config.Wait):
		client.Close()
	}

	if err := g.Wait(); err != nil {
		return err
	}

	// 输出结果
	if config.Json {
		outputJSON(map[string]interface{}{
			"sent":     len(data),
			"received": len(allData),
			"data":     string(allData),
		})
	} else if !config.Quiet {
		fmt.Printf("\nTotal received: %d bytes\n", len(allData))
		if config.Hex {
			fmt.Printf("Response (hex):\n%s\n", hexDump(allData))
		} else {
			fmt.Printf("Response:\n%s\n", string(allData))
		}
	}

	return nil
}

// RunPathMode 路径发送模式
// 发送多个文件，每个文件发送后等待响应
func RunPathMode(config TcpConfig, ipAddr net.IP) error {
	client := NewTCPClient(config)
	if err := client.Connect(ipAddr.String(), config.Port); err != nil {
		return err
	}
	defer client.Close()

	// 解析路径
	pathInfo, err := resolvePath(config.Path, config.MaxFileSize)
	if err != nil {
		return err
	}

	// 使用 errgroup 管理接收
	g, ctx := errgroup.WithContext(context.Background())
	client.StartReceive(g, 0) // 0 表示使用默认短超时

	// 发送每个文件并等待响应
	results := make([]map[string]interface{}, 0, len(pathInfo.Paths))
	
	for _, filePath := range pathInfo.Paths {
		data, err := os.ReadFile(filePath)
		if err != nil {
			if !config.Quiet {
				fmt.Printf("Warning: failed to read %s: %v\n", filePath, err)
			}
			continue
		}

		// 发送文件
		if err := client.Send(data); err != nil {
			return fmt.Errorf("send %s failed: %w", filePath, err)
		}

		if !config.Quiet && !config.Json {
			fmt.Printf("[→] Sent %s (%d bytes)\n", filePath, len(data))
		}

		// 等待该文件的响应（使用 --wait 指定的时间）
		fileResult := map[string]interface{}{
			"file": filePath,
			"sent": len(data),
		}
		
		var response []byte
		// 使用用户指定的 --wait 时间，如果为 0 则使用默认 500ms
		responseTimeout := config.Wait
		if responseTimeout <= 0 {
			responseTimeout = 500 * time.Millisecond
		}

		select {
		case recvData := <-client.GetRecvChan():
			response = recvData
			fileResult["received"] = len(recvData)
			if !config.Quiet && !config.Json {
				if config.Hex {
					fmt.Printf("[←] Received %d bytes:\n%s\n", len(recvData), hexDump(recvData))
				} else {
					fmt.Printf("[←] Received %d bytes:\n%s\n", len(recvData), string(recvData))
				}
			}
		case <-time.After(responseTimeout):
			fileResult["timeout"] = true
			if !config.Quiet && !config.Json {
				fmt.Printf("[←] Timeout (no response within %v)\n", responseTimeout)
			}
		case <-ctx.Done():
			return ctx.Err()
		}

		if len(response) > 0 {
			fileResult["response"] = string(response)
		}
		results = append(results, fileResult)

		// 文件间短暂延迟
		time.Sleep(10 * time.Millisecond)
	}

	client.Close()
	if err := g.Wait(); err != nil {
		return err
	}

	// 输出结果
	if config.Json {
		outputJSON(map[string]interface{}{
			"files":   results,
			"total":   len(results),
		})
	}

	return nil
}

// RunInteractiveMode 交互式发送模式
// 发送和接收完全并行
func RunInteractiveMode(config TcpConfig, ipAddr net.IP) error {
	// 交互式模式不支持 JSON 输出
	if config.Json {
		return fmt.Errorf("interactive mode does not support JSON output")
	}

	client := NewTCPClient(config)
	if err := client.Connect(ipAddr.String(), config.Port); err != nil {
		return err
	}
	defer client.Close()

	if !config.Quiet {
		fmt.Printf("Connected to %s\n", client.remoteAddr)
		fmt.Println("Type 'help' for available commands, 'quit' to exit")
		fmt.Println()
	}

	// 配置 readline
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "tcp> ",
		AutoComplete:    tcpCompleter,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		return fmt.Errorf("failed to initialize readline: %w", err)
	}
	defer rl.Close()

	// 使用 errgroup 管理接收和输入
	g, ctx := errgroup.WithContext(context.Background())

	// 启动接收
	client.StartReceive(g, 0)

	// 输入 channel
	inputChan := make(chan string)
	inputErrChan := make(chan error)

	// 启动输入读取 goroutine
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return nil
			default:
			}

			line, err := rl.Readline()
			if err != nil {
				inputErrChan <- err
				return nil
			}
			inputChan <- line
		}
	})

	// 主循环
	for {
		select {
		case data := <-client.GetRecvChan():
			// 立即显示接收数据
			if !config.Quiet {
				fmt.Println()
				if config.Hex {
					fmt.Printf("[←] Received %d bytes:\n%s\n", 
						len(data), hexDump(data))
				} else {
					fmt.Printf("[←] Received %d bytes:\n%s\n", 
						len(data), string(data))
				}
			}
			// 刷新提示符
			rl.Refresh()

		case line := <-inputChan:
			// 处理输入
			if err := handleInteractiveLine(client, line); err != nil {
				if err.Error() == "quit" {
					return nil
				}
				fmt.Printf("Error: %v\n", err)
			}

		case err := <-inputErrChan:
			// EOF 或中断
			return nil

		case <-ctx.Done():
			return nil
		}
	}
}

// handleInteractiveLine 处理交互式输入
func handleInteractiveLine(client *TCPClient, line string) error {
	line = trimSpace(line)

	if line == "" {
		return client.SendString("\n")
	}

	parts := splitFields(line)
	cmd := toLower(parts[0])

	switch cmd {
	case "quit", "exit":
		return fmt.Errorf("quit")

	case "help":
		printInteractiveHelp()
		return nil

	case "close":
		client.Close()
		return fmt.Errorf("quit")

	case "status":
		sent, recv, sendSeq, recvSeq := client.GetStats()
		fmt.Printf("Bytes sent: %d (seq=%d), received: %d (seq=%d)\n",
			sent, sendSeq, recv, recvSeq)
		return nil

	default:
		return client.SendString(line + "\n")
	}
}
```

### 3.3 cmd_tcp.go - 简化入口

```go
package tcp

import (
	"fmt"
	"net"
)

// TcpCmdMain TCP 命令主入口（简化版）
func TcpCmdMain(config TcpConfig) error {
	// 监听模式
	if config.Listen {
		return runListenMode(config)
	}

	// 扫描模式
	if config.Scan || config.Range != "" {
		return runScanMode(config)
	}

	// 解析目标地址
	ipAddr, err := resolveHost(config.Host)
	if err != nil {
		return fmt.Errorf("failed to resolve host: %w", err)
	}

	// Banner 模式
	if config.Banner {
		return runBannerMode(config, ipAddr)
	}

	// 交互式模式
	if config.Interactive {
		return RunInteractiveMode(config, ipAddr)
	}

	// 路径发送模式
	if config.Path != "" {
		return RunPathMode(config, ipAddr)
	}

	// 字符串发送模式（默认）
	if config.Data != "" {
		return RunStringMode(config, ipAddr)
	}

	// 默认连接测试
	return runConnectMode(config, ipAddr)
}
```

---

## 4. 关键改进点

### 4.1 统一封装
- `TCPClient` 封装连接、发送、接收、统计
- 三种模式复用相同的基础功能

### 4.2 并发管理
- `errgroup` 管理 goroutine 生命周期
- 统一错误处理和取消机制

### 4.3 三种模式差异

| 模式 | 发送策略 | 接收策略 |
|------|----------|----------|
| 字符串 | 单次发送 | 持续接收直到超时 |
| 路径 | 逐个文件发送 | 每个文件后等待响应 |
| 交互式 | 用户输入驱动 | 并行持续接收 |

---

## 5. 依赖

```go
// go.mod
go get golang.org/x/sync/errgroup
```

---

## 6. 测试用例

```bash
# 1. 字符串发送模式
fck tcp -d "hello" -w 3s example.com 80

# 2. 路径发送模式
fck tcp -p ./requests/ -w 1s example.com 80

# 3. 交互式模式
fck tcp -I example.com 80
```

---

**文档版本**: v2.0  
**创建日期**: 2026-04-27  
**作者**: AI Assistant
