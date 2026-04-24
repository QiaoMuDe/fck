package tcp

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chzyer/readline"
)

// 交互式模式命令常量
const (
	CmdHelp   = "help"   // 显示帮助信息
	CmdQuit   = "quit"   // 退出交互模式
	CmdExit   = "exit"   // 退出交互模式（quit的别名）
	CmdClose  = "close"  // 关闭当前连接
	CmdHex    = "hex"    // 切换十六进制显示模式
	CmdStatus = "status" // 显示连接状态
	CmdFile   = "file"   // 发送文件内容
)

// InteractiveSession 交互式会话状态
type InteractiveSession struct {
	Conn      net.Conn
	Config    TcpConfig
	HexMode   bool
	StartTime time.Time
	BytesSent int64
	BytesRecv int64
	RL        *readline.Instance
}

// tcpCompleter Tab 补全器
var tcpCompleter = readline.NewPrefixCompleter(
	readline.PcItem(CmdHelp),
	readline.PcItem(CmdQuit),
	readline.PcItem(CmdExit),
	readline.PcItem(CmdClose),
	readline.PcItem(CmdHex),
	readline.PcItem(CmdStatus),
	readline.PcItem(CmdFile, readline.PcItemDynamic(listFiles)),
)

// listFiles 返回当前目录下的文件列表用于补全
// line 参数是用户当前输入的整行内容，如 "file " 或 "file READ"
func listFiles(line string) []string {
	// 提取 file 命令后的路径前缀
	prefix := ""
	parts := strings.SplitN(line, " ", 2)
	if len(parts) > 1 {
		prefix = strings.TrimSpace(parts[1])
	}

	dir := "."
	filePrefix := ""

	// 解析路径和文件名前缀
	if prefix != "" {
		// 统一使用 / 作为分隔符处理
		normalizedPath := strings.ReplaceAll(prefix, "\\", "/")
		if idx := strings.LastIndex(normalizedPath, "/"); idx >= 0 {
			dir = normalizedPath[:idx]
			if dir == "" {
				dir = "/"
			}
			filePrefix = normalizedPath[idx+1:]
		} else {
			filePrefix = normalizedPath
		}
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var files []string
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(filePrefix)) {
			if entry.IsDir() {
				files = append(files, name+"/")
			} else {
				files = append(files, name)
			}
		}
	}
	return files
}

// runInteractiveMode 运行交互式模式
//
// 参数:
//   - config: TCP配置
//   - ipAddr: 目标IP地址
//
// 返回值:
//   - error: 执行错误
func runInteractiveMode(config TcpConfig, ipAddr net.IP) error {
	target := net.JoinHostPort(ipAddr.String(), strconv.Itoa(config.Port))

	if !config.Quiet && !config.Json {
		fmt.Printf("Connecting to %s (%s):%d...\n", config.Host, ipAddr.String(), config.Port)
	}

	// 建立连接
	conn, err := net.DialTimeout("tcp", target, config.Timeout)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer func() { _ = conn.Close() }()

	if !config.Quiet && !config.Json {
		fmt.Printf("Connected to %s:%d\n", config.Host, config.Port)
		fmt.Println("Type 'help' for available commands, 'quit' to exit")
		fmt.Println()
	}

	// 配置 readline（不保存历史记录到文件）
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "tcp> ",
		AutoComplete:    tcpCompleter,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		return fmt.Errorf("failed to initialize readline: %w", err)
	}
	defer func() { _ = rl.Close() }()

	// 创建会话
	session := &InteractiveSession{
		Conn:      conn,
		Config:    config,
		HexMode:   config.Hex,
		StartTime: time.Now(),
		RL:        rl,
	}

	// 启动接收 goroutine
	recvChan := make(chan []byte, 10)
	errChan := make(chan error, 1)
	done := make(chan bool)
	go receiveLoopInteractive(session, recvChan, errChan, done)

	// 主循环
	for {
		// 检查是否有接收到的数据
		select {
		case data := <-recvChan:
			printReceivedDataInteractive(data, session.HexMode, session.Config.Quiet)
		case err := <-errChan:
			if err != nil {
				if !config.Quiet && !config.Json {
					fmt.Printf("\n[Error] %v\n", err)
				}
			}
			return nil
		default:
			// 继续读取输入
		}

		// 读取用户输入
		line, err := rl.Readline()
		if err != nil {
			// EOF 或中断
			break
		}

		// 处理输入
		if err := handleInteractiveInput(session, line); err != nil {
			if err.Error() == "quit" {
				break
			}
			if !config.Quiet {
				fmt.Printf("Error: %v\n", err)
			}
		}
	}

	close(done)

	if !config.Quiet && !config.Json {
		fmt.Println("Connection closed.")
	}

	return nil
}

// handleInteractiveInput 处理用户输入
//
// 参数:
//   - session: 交互式会话
//   - input: 用户输入
//
// 返回值:
//   - error: 处理错误,返回 "quit" 表示退出
func handleInteractiveInput(session *InteractiveSession, input string) error {
	input = strings.TrimSpace(input)

	// 空行处理
	if input == "" {
		// 发送换行符
		_, err := session.Conn.Write([]byte("\n"))
		if err != nil {
			return fmt.Errorf("failed to send: %w", err)
		}
		session.BytesSent += 1
		return nil
	}

	// 处理特殊命令
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return nil
	}

	cmd := strings.ToLower(parts[0])

	switch cmd {
	case CmdQuit, CmdExit:
		return fmt.Errorf("quit")

	case CmdHelp:
		printInteractiveHelp()
		return nil

	case CmdClose:
		_ = session.Conn.Close()
		return fmt.Errorf("connection closed")

	case CmdHex:
		session.HexMode = !session.HexMode
		if !session.Config.Quiet {
			if session.HexMode {
				fmt.Println("Hex mode: ON")
			} else {
				fmt.Println("Hex mode: OFF")
			}
		}
		return nil

	case CmdStatus:
		printInteractiveStatus(session)
		return nil

	case CmdFile:
		if len(parts) < 2 {
			return fmt.Errorf("usage: %s <path>", CmdFile)
		}
		return sendFileInteractive(session, parts[1])

	default:
		// 普通数据,发送给服务器
		data := []byte(input + "\n")
		n, err := session.Conn.Write(data)
		if err != nil {
			return fmt.Errorf("failed to send: %w", err)
		}
		session.BytesSent += int64(n)

		if !session.Config.Quiet && !session.Config.Json {
			fmt.Printf("[→] Sent %d bytes\n", n)
		}
	}

	return nil
}

// receiveLoopInteractive 接收数据循环
//
// 参数:
//   - session: 交互式会话
//   - recvChan: 接收数据通道
//   - errChan: 错误通道
//   - done: 退出信号通道
func receiveLoopInteractive(session *InteractiveSession, recvChan chan<- []byte, errChan chan<- error, done <-chan bool) {
	buf := make([]byte, 4096)

	for {
		select {
		case <-done:
			return
		default:
		}

		// 设置读取超时
		if session.Config.Wait > 0 {
			_ = session.Conn.SetReadDeadline(time.Now().Add(session.Config.Wait))
		}

		n, err := session.Conn.Read(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// 超时,继续
				continue
			}
			// 连接关闭或其他错误
			errChan <- err
			return
		}

		if n > 0 {
			data := make([]byte, n)
			copy(data, buf[:n])
			session.BytesRecv += int64(n)

			select {
			case recvChan <- data:
			case <-done:
				return
			}
		}
	}
}

// printReceivedDataInteractive 打印接收到的数据
//
// 参数:
//   - data: 接收到的数据
//   - hexMode: 是否以十六进制格式显示
//   - quiet: 是否静默模式
func printReceivedDataInteractive(data []byte, hexMode, quiet bool) {
	if quiet {
		return
	}

	if hexMode {
		fmt.Printf("[←] Received (%d bytes):\n", len(data))
		printHexDump(data)
	} else {
		// 尝试作为文本显示
		text := string(data)
		// 去除末尾的换行符
		text = strings.TrimRight(text, "\r\n")
		fmt.Printf("[←] Received (%d bytes): %s\n", len(data), text)
	}
}

// printHexDump 打印十六进制转储
//
// 参数:
//   - data: 要打印的数据
func printHexDump(data []byte) {
	const bytesPerLine = 16

	for i := 0; i < len(data); i += bytesPerLine {
		end := i + bytesPerLine
		if end > len(data) {
			end = len(data)
		}

		// 打印偏移量
		fmt.Printf("%08x  ", i)

		// 打印十六进制
		for j := i; j < i+bytesPerLine; j++ {
			if j < end {
				fmt.Printf("%02x ", data[j])
			} else {
				fmt.Print("   ")
			}
			if j == i+bytesPerLine/2-1 {
				fmt.Print(" ")
			}
		}

		// 打印 ASCII
		fmt.Print(" |")
		for j := i; j < end; j++ {
			if data[j] >= 32 && data[j] < 127 {
				fmt.Printf("%c", data[j])
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println("|")
	}
}

// printInteractiveHelp 打印交互式模式帮助信息
func printInteractiveHelp() {
	fmt.Println("Available commands:")
	fmt.Printf("  %-17s Show this help message\n", CmdHelp)
	fmt.Printf("  %-17s Exit interactive mode\n", CmdQuit+", "+CmdExit)
	fmt.Printf("  %-17s Close current connection\n", CmdClose)
	fmt.Printf("  %-17s Toggle hex display mode\n", CmdHex)
	fmt.Printf("  %-17s Show connection status\n", CmdStatus)
	fmt.Printf("  %-17s Send file contents\n", CmdFile+" <path>")
	fmt.Println("")
	fmt.Println("Any other text will be sent to the server.")
	fmt.Println("Use Tab for command completion, ↑↓ for history.")
}

// printInteractiveStatus 打印交互式会话状态
//
// 参数:
//   - session: 交互式会话
func printInteractiveStatus(session *InteractiveSession) {
	duration := time.Since(session.StartTime)
	fmt.Println("Connection Status:")
	fmt.Printf("  Remote:    %s\n", session.Conn.RemoteAddr())
	fmt.Printf("  Local:     %s\n", session.Conn.LocalAddr())
	fmt.Printf("  Duration:  %s\n", duration.Round(time.Second))
	fmt.Printf("  Sent:      %d bytes\n", session.BytesSent)
	fmt.Printf("  Received:  %d bytes\n", session.BytesRecv)
	fmt.Printf("  Hex mode:  %v\n", session.HexMode)
}

// sendFileInteractive 发送文件内容
//
// 参数:
//   - session: 交互式会话
//   - path: 文件路径
//
// 返回值:
//   - error: 发送错误
func sendFileInteractive(session *InteractiveSession, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	n, err := session.Conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send: %w", err)
	}

	session.BytesSent += int64(n)

	if !session.Config.Quiet && !session.Config.Json {
		fmt.Printf("[→] Sent file '%s' (%d bytes)\n", path, n)
	}

	return nil
}
