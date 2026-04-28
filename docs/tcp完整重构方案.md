# TCP 模块完整重构方案

## 1. 重构目标

### 1.1 核心目标
- 统一三种发送模式（字符串、路径、交互式）的底层实现
- 解决交互式模式接收阻塞问题
- 支持持续接收服务端响应（而非单次读取）
- 代码结构清晰，易于维护和扩展

### 1.2 设计原则
- **单一职责**：每个文件只负责一个功能模块
- **并发安全**：使用 channel 和 goroutine 实现真正的并行
- **向后兼容**：CLI 接口保持不变
- **可测试性**：模块间通过接口解耦

---

## 2. 新文件结构

```
internal/commands/tcp/
├── types.go           # 类型定义（配置、接口等）
├── session.go         # TCP 会话管理（核心）
├── sender.go          # 发送模块
├── receiver.go        # 接收模块
├── handler.go         # 数据处理回调
├── modes.go           # 三种模式统一入口
│   ├── runStringMode()
│   ├── runPathMode()
│   └── runInteractiveMode()
├── cmd_tcp.go         # 主命令（简化后的入口）
├── listen.go          # 监听模式（独立）
├── scan.go            # 扫描模式（独立）
├── banner.go          # Banner 模式（独立）
├── utils.go           # 工具函数
└── output.go          # 输出处理
```

---

## 3. 详细设计

### 3.1 types.go - 类型定义

```go
package tcp

import (
	"net"
	"time"
)

// Mode 发送模式类型
type Mode int

const (
	ModeString      Mode = iota // 字符串发送模式
	ModePath                    // 路径发送模式
	ModeInteractive             // 交互式模式
)

// TcpConfig 配置结构体（保持向后兼容）
type TcpConfig struct {
	Host        string
	Port        int
	Timeout     time.Duration
	Count       int
	Interval    time.Duration
	Scan        bool
	Range       string
	OpenOnly    bool
	Banner      bool
	Data        string
	Path        string
	MaxFileSize int64
	Wait        time.Duration
	Quiet       bool
	Json        bool
	Listen      bool
	Interactive bool
	Hex         bool
}

// DataHandler 数据处理回调接口
type DataHandler interface {
	// OnConnect 连接建立时调用
	OnConnect(remoteAddr string)
	
	// OnSend 数据发送时调用
	OnSend(data []byte, seq int)
	
	// OnReceive 接收到数据时调用
	OnReceive(data []byte, seq int) bool // 返回 false 停止接收
	
	// OnError 发生错误时调用
	OnError(err error)
	
	// OnClose 连接关闭时调用
	OnClose(totalSent, totalRecv int64, duration time.Duration)
}

// SendRequest 发送请求
type SendRequest struct {
	Data     []byte
	FilePath string // 可选，用于标识数据来源
}

// ReceiveResult 接收结果
type ReceiveResult struct {
	Data []byte
	Seq  int
}
```

### 3.2 session.go - TCP 会话管理（核心）

```go
package tcp

import (
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// Session TCP 会话
type Session struct {
	// 配置
	Config TcpConfig
	
	// 连接
	Conn       net.Conn
	RemoteAddr string
	
	// 通道
	SendChan   chan *SendRequest  // 发送队列
	RecvChan   chan []byte        // 接收数据
	ErrChan    chan error         // 错误
	Done       chan struct{}      // 退出信号
	
	// 状态
	StartTime   time.Time
	BytesSent   int64
	BytesRecv   int64
	SendSeq     int32
	RecvSeq     int32
	IsConnected atomic.Bool
	
	// 处理器
	Handler DataHandler
	
	// 内部控制
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewSession 创建新会话
func NewSession(config TcpConfig, handler DataHandler) *Session {
	ctx, cancel := context.WithCancel(context.Background())
	return &Session{
		Config:     config,
		SendChan:   make(chan *SendRequest, 100),
		RecvChan:   make(chan []byte, 100),
		ErrChan:    make(chan error, 1),
		Done:       make(chan struct{}),
		Handler:    handler,
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Connect 建立连接
func (s *Session) Connect(host string, port int) error {
	target := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	
	conn, err := net.DialTimeout("tcp", target, s.Config.Timeout)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	
	s.Conn = conn
	s.RemoteAddr = conn.RemoteAddr().String()
	s.StartTime = time.Now()
	s.IsConnected.Store(true)
	
	if s.Handler != nil {
		s.Handler.OnConnect(s.RemoteAddr)
	}
	
	return nil
}

// Start 启动会话（启动发送和接收 goroutine）
func (s *Session) Start() {
	s.wg.Add(2)
	go s.sendLoop()
	go s.recvLoop()
}

// sendLoop 发送循环
func (s *Session) sendLoop() {
	defer s.wg.Done()
	
	for {
		select {
		case <-s.ctx.Done():
			return
		case req := <-s.SendChan:
			if req == nil {
				return
			}
			
			n, err := s.Conn.Write(req.Data)
			if err != nil {
				select {
				case s.ErrChan <- fmt.Errorf("send failed: %w", err):
				case <-s.ctx.Done():
				}
				return
			}
			
			atomic.AddInt64(&s.BytesSent, int64(n))
			seq := atomic.AddInt32(&s.SendSeq, 1)
			
			if s.Handler != nil {
				s.Handler.OnSend(req.Data, int(seq))
			}
		}
	}
}

// recvLoop 接收循环
func (s *Session) recvLoop() {
	defer s.wg.Done()
	
	buf := make([]byte, 4096)
	
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
		}
		
		// 设置读取超时
		readTimeout := s.Config.Wait
		if s.Config.Interactive && readTimeout <= 0 {
			readTimeout = 100 * time.Millisecond
		}
		
		if readTimeout > 0 {
			_ = s.Conn.SetReadDeadline(time.Now().Add(readTimeout))
		}
		
		n, err := s.Conn.Read(buf)
		
		if n > 0 {
			data := make([]byte, n)
			copy(data, buf[:n])
			
			atomic.AddInt64(&s.BytesRecv, int64(n))
			seq := atomic.AddInt32(&s.RecvSeq, 1)
			
			// 发送到 RecvChan
			select {
			case s.RecvChan <- data:
			case <-s.ctx.Done():
				return
			}
			
			// 调用处理器
			if s.Handler != nil {
				if !s.Handler.OnReceive(data, int(seq)) {
					// 处理器要求停止
					return
				}
			}
		}
		
		if err != nil {
			// 区分超时和连接错误
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// 超时继续
				continue
			}
			
			// 连接关闭或其他错误
			if err.Error() != "EOF" {
				select {
				case s.ErrChan <- fmt.Errorf("receive failed: %w", err):
				case <-s.ctx.Done():
				}
			}
			return
		}
	}
}

// Send 发送数据（非阻塞）
func (s *Session) Send(data []byte, filePath string) error {
	select {
	case s.SendChan <- &SendRequest{Data: data, FilePath: filePath}:
		return nil
	case <-s.ctx.Done():
		return fmt.Errorf("session closed")
	}
}

// SendString 发送字符串
func (s *Session) SendString(data string) error {
	return s.Send([]byte(data), "")
}

// Close 关闭会话
func (s *Session) Close() error {
	s.cancel()
	s.IsConnected.Store(false)
	
	close(s.SendChan)
	
	if s.Conn != nil {
		_ = s.Conn.Close()
	}
	
	// 等待 goroutine 结束
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()
	
	select {
	case <-done:
	case <-time.After(5 * time.Second):
	}
	
	close(s.Done)
	
	// 调用关闭回调
	if s.Handler != nil {
		duration := time.Since(s.StartTime)
		s.Handler.OnClose(
			atomic.LoadInt64(&s.BytesSent),
			atomic.LoadInt64(&s.BytesRecv),
			duration,
		)
	}
	
	return nil
}

// Wait 等待会话结束
func (s *Session) Wait() error {
	select {
	case <-s.Done:
		return nil
	case err := <-s.ErrChan:
		return err
	}
}
```

### 3.3 handler.go - 默认处理器实现

```go
package tcp

import (
	"fmt"
	"strings"
	"time"
)

// ConsoleHandler 控制台输出处理器
type ConsoleHandler struct {
	Config      TcpConfig
	PrintPrefix bool
}

// NewConsoleHandler 创建控制台处理器
func NewConsoleHandler(config TcpConfig) *ConsoleHandler {
	return &ConsoleHandler{
		Config:      config,
		PrintPrefix: true,
	}
}

func (h *ConsoleHandler) OnConnect(remoteAddr string) {
	if h.Config.Quiet || h.Config.Json {
		return
	}
	fmt.Printf("Connected to %s\n", remoteAddr)
}

func (h *ConsoleHandler) OnSend(data []byte, seq int) {
	if h.Config.Quiet || h.Config.Json {
		return
	}
	if h.PrintPrefix {
		fmt.Printf("[→] Sent %d bytes (seq=%d)\n", len(data), seq)
	}
}

func (h *ConsoleHandler) OnReceive(data []byte, seq int) bool {
	if h.Config.Quiet {
		return true
	}
	
	if h.Config.Json {
		// JSON 模式下累积数据，不直接打印
		return true
	}
	
	if h.Config.Hex {
		fmt.Printf("[←] Received %d bytes (seq=%d):\n%s\n", 
			len(data), seq, hexDump(data))
	} else {
		content := string(data)
		// 去除末尾换行，避免双换行
		content = strings.TrimRight(content, "\n\r")
		fmt.Printf("[←] Received %d bytes (seq=%d):\n%s\n", 
			len(data), seq, content)
	}
	
	return true // 继续接收
}

func (h *ConsoleHandler) OnError(err error) {
	if h.Config.Quiet {
		return
	}
	fmt.Printf("[Error] %v\n", err)
}

func (h *ConsoleHandler) OnClose(totalSent, totalRecv int64, duration time.Duration) {
	if h.Config.Quiet || h.Config.Json {
		return
	}
	fmt.Printf("\nConnection closed.\n")
	fmt.Printf("Total sent: %d bytes, received: %d bytes, duration: %s\n",
		totalSent, totalRecv, duration.Round(time.Millisecond))
}

// JSONHandler JSON 输出处理器
type JSONHandler struct {
	Config       TcpConfig
	ConnectInfo  map[string]interface{}
	SentRecords  []map[string]interface{}
	RecvRecords  []map[string]interface{}
}

func NewJSONHandler(config TcpConfig) *JSONHandler {
	return &JSONHandler{
		Config:      config,
		ConnectInfo: make(map[string]interface{}),
		SentRecords: make([]map[string]interface{}, 0),
		RecvRecords: make([]map[string]interface{}, 0),
	}
}

func (h *JSONHandler) OnConnect(remoteAddr string) {
	h.ConnectInfo["remote_addr"] = remoteAddr
	h.ConnectInfo["connect_time"] = time.Now().Format(time.RFC3339)
}

func (h *JSONHandler) OnSend(data []byte, seq int) {
	h.SentRecords = append(h.SentRecords, map[string]interface{}{
		"seq":        seq,
		"bytes":      len(data),
		"timestamp":  time.Now().Format(time.RFC3339),
	})
}

func (h *JSONHandler) OnReceive(data []byte, seq int) bool {
	h.RecvRecords = append(h.RecvRecords, map[string]interface{}{
		"seq":       seq,
		"bytes":     len(data),
		"data":      string(data),
		"timestamp": time.Now().Format(time.RFC3339),
	})
	return true
}

func (h *JSONHandler) OnError(err error) {
	h.ConnectInfo["error"] = err.Error()
}

func (h *JSONHandler) OnClose(totalSent, totalRecv int64, duration time.Duration) {
	h.ConnectInfo["total_sent"] = totalSent
	h.ConnectInfo["total_recv"] = totalRecv
	h.ConnectInfo["duration_ms"] = duration.Milliseconds()
	h.ConnectInfo["sent_records"] = h.SentRecords
	h.ConnectInfo["recv_records"] = h.RecvRecords
	
	// 输出 JSON
	outputJSON(h.ConnectInfo)
}

// InteractiveHandler 交互式处理器（特殊处理）
type InteractiveHandler struct {
	*ConsoleHandler
	Prompt      string
}

func NewInteractiveHandler(config TcpConfig) *InteractiveHandler {
	return &InteractiveHandler{
		ConsoleHandler: NewConsoleHandler(config),
		Prompt:         "tcp> ",
	}
}

func (h *InteractiveHandler) OnReceive(data []byte, seq int) bool {
	// 先打印换行，避免和输入行混在一起
	fmt.Println()
	return h.ConsoleHandler.OnReceive(data, seq)
}
```

### 3.4 modes.go - 三种模式统一实现

```go
package tcp

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chzyer/readline"
)

// RunStringMode 运行字符串发送模式
func RunStringMode(config TcpConfig, ipAddr net.IP) error {
	handler := NewConsoleHandler(config)
	if config.Json {
		handler = NewJSONHandler(config)
	}
	
	session := NewSession(config, handler)
	
	// 连接
	if err := session.Connect(ipAddr.String(), config.Port); err != nil {
		return err
	}
	defer session.Close()
	
	// 启动会话
	session.Start()
	
	// 发送数据
	data := config.Data
	if !strings.HasSuffix(data, "\n") {
		data += "\n"
	}
	
	if err := session.SendString(data); err != nil {
		return err
	}
	
	// 等待接收（如果有 -w）
	if config.Wait > 0 {
		time.Sleep(config.Wait)
	}
	
	return nil
}

// RunPathMode 运行路径发送模式
func RunPathMode(config TcpConfig, ipAddr net.IP) error {
	handler := NewConsoleHandler(config)
	if config.Json {
		handler = NewJSONHandler(config)
	}
	
	session := NewSession(config, handler)
	
	if err := session.Connect(ipAddr.String(), config.Port); err != nil {
		return err
	}
	defer session.Close()
	
	session.Start()
	
	// 解析路径
	pathInfo, err := resolvePath(config.Path, config.MaxFileSize)
	if err != nil {
		return err
	}
	
	// 发送所有文件
	for _, filePath := range pathInfo.Paths {
		data, err := os.ReadFile(filePath)
		if err != nil {
			if !config.Quiet {
				fmt.Printf("Warning: failed to read %s: %v\n", filePath, err)
			}
			continue
		}
		
		if err := session.Send(data, filePath); err != nil {
			return err
		}
		
		// 文件间短暂延迟，避免发送过快
		time.Sleep(10 * time.Millisecond)
	}
	
	// 等待接收
	if config.Wait > 0 {
		time.Sleep(config.Wait)
	}
	
	return nil
}

// RunInteractiveMode 运行交互式模式
func RunInteractiveMode(config TcpConfig, ipAddr net.IP) error {
	handler := NewInteractiveHandler(config)
	session := NewSession(config, handler)
	
	if err := session.Connect(ipAddr.String(), config.Port); err != nil {
		return err
	}
	defer session.Close()
	
	session.Start()
	
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
	
	// 输入 channel
	inputChan := make(chan string)
	inputErrChan := make(chan error)
	
	// 启动输入读取 goroutine
	go func() {
		for {
			line, err := rl.Readline()
			if err != nil {
				inputErrChan <- err
				return
			}
			inputChan <- line
		}
	}()
	
	// 主循环：真正并行处理接收和输入
	for {
		select {
		case data := <-session.RecvChan:
			// 立即显示接收的数据
			handler.OnReceive(data, int(session.RecvSeq))
			// 刷新提示符
			rl.Refresh()
			
		case line := <-inputChan:
			// 处理用户输入
			if err := handleInteractiveLine(session, line); err != nil {
				if err.Error() == "quit" {
					return nil
				}
				fmt.Printf("Error: %v\n", err)
			}
			
		case err := <-inputErrChan:
			// EOF 或中断
			if err.Error() == "EOF" {
				return nil
			}
			return err
			
		case err := <-session.ErrChan:
			return err
			
		case <-session.Done:
			return nil
		}
	}
}

// handleInteractiveLine 处理交互式输入行
func handleInteractiveLine(session *Session, line string) error {
	line = strings.TrimSpace(line)
	
	if line == "" {
		// 空行发送换行符
		return session.SendString("\n")
	}
	
	// 处理命令
	parts := strings.Fields(line)
	cmd := strings.ToLower(parts[0])
	
	switch cmd {
	case "quit", "exit":
		return fmt.Errorf("quit")
		
	case "help":
		printInteractiveHelp()
		return nil
		
	case "close":
		session.Close()
		return fmt.Errorf("quit")
		
	case "status":
		fmt.Printf("Connected: %v\n", session.IsConnected.Load())
		fmt.Printf("Bytes sent: %d, received: %d\n", 
			session.BytesSent, session.BytesRecv)
		return nil
		
	case "path":
		if len(parts) < 2 {
			return fmt.Errorf("usage: path <filepath>")
		}
		return sendPathInteractive(session, parts[1])
		
	default:
		// 普通数据
		return session.SendString(line + "\n")
	}
}

// sendPathInteractive 交互式发送文件
func sendPathInteractive(session *Session, path string) error {
	pathInfo, err := resolvePath(path, session.Config.MaxFileSize)
	if err != nil {
		return err
	}
	
	for _, filePath := range pathInfo.Paths {
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Warning: failed to read %s: %v\n", filePath, err)
			continue
		}
		
		if err := session.Send(data, filePath); err != nil {
			return err
		}
		fmt.Printf("Sent file: %s (%d bytes)\n", filePath, len(data))
	}
	
	return nil
}

// printInteractiveHelp 打印交互式帮助
func printInteractiveHelp() {
	fmt.Println("Interactive mode commands:")
	fmt.Println("  help        Show this help")
	fmt.Println("  quit/exit   Exit interactive mode")
	fmt.Println("  close       Close connection")
	fmt.Println("  status      Show connection status")
	fmt.Println("  path <file> Send file")
	fmt.Println("")
	fmt.Println("Any other text will be sent to the server.")
}
```

### 3.5 cmd_tcp.go - 简化后的主入口

```go
package tcp

import (
	"fmt"
	"net"
)

// TcpCmdMain 执行 TCP 命令（简化版）
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
	
	// 默认连接测试模式
	return runConnectMode(config, ipAddr)
}
```

---

## 4. 关键改进点

### 4.1 真正的并行处理
```
旧模式：
  主循环 → 检查接收（非阻塞）→ 阻塞等待输入 → 处理输入 → 循环

新模式：
  发送 goroutine ─┐
                  ├→ 网络连接
  接收 goroutine ─┘      ↑
                         ↓
  主循环: select { RecvChan, InputChan, ErrChan }
```

### 4.2 统一的数据处理
- 三种模式使用相同的 `Session` 和 `DataHandler` 接口
- 输出格式统一（控制台/JSON/十六进制）
- 统计信息统一（发送/接收字节数、序列号等）

### 4.3 交互式模式改进
- 接收数据立即显示，无需等待回车
- 使用 `rl.Refresh()` 刷新提示符
- 支持同时显示接收数据和用户输入

---

## 5. 测试用例

### 5.1 字符串发送模式
```bash
# 基本发送
fck tcp -d "hello" -w 2s example.com 80

# JSON 输出
fck tcp -d "hello" -w 2s -j example.com 80

# 十六进制显示
fck tcp -d "hello" -w 2s --hex example.com 80
```

### 5.2 路径发送模式
```bash
# 发送单个文件
fck tcp -p ./request.txt -w 3s example.com 80

# 发送目录
fck tcp -p ./requests/ -w 5s example.com 80

# 通配符
fck tcp -p "./logs/*.log" -w 10s example.com 8080
```

### 5.3 交互式模式
```bash
# 基本交互
fck tcp -I example.com 80

# 带超时
fck tcp -I -w 500ms example.com 80

# 十六进制显示
fck tcp -I --hex example.com 80
```

---

## 6. 实现优先级

| 优先级 | 文件 | 说明 |
|--------|------|------|
| P0 | types.go | 基础类型定义 |
| P0 | session.go | 核心会话管理 |
| P0 | handler.go | 处理器实现 |
| P0 | modes.go | 三种模式实现 |
| P1 | cmd_tcp.go | 简化入口 |
| P1 | utils.go | 工具函数迁移 |
| P2 | output.go | 输出统一处理 |
| P2 | 测试 | 三种模式全面测试 |

---

## 7. 向后兼容性

- CLI 参数完全兼容
- JSON 输出格式保持兼容
- 配置文件兼容
- 行为兼容（默认模式不变）

---

**文档版本**: v1.0  
**创建日期**: 2026-04-27  
**作者**: AI Assistant
