package tcp

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chzyer/readline"
)

// InteractiveSession 交互式会话状态
type InteractiveSession struct {
	conn         net.Conn       // TCP 连接
	config       ClientConfig   // 客户端配置
	stats        *TransferStats // 传输统计
	rl           *readline.Instance
	history      []string      // 内存中的历史记录
	lastSend     string        // 上一次发送的内容
	interval     time.Duration // 发送间隔
	msgCount     int64         // 消息计数
	lastActivity time.Time     // 最后活动时间
}

// runInteractiveMode 启动 readline 交互模式
//
// 参数:
//   - conn: TCP 连接
//   - config: 客户端配置
//
// 返回值:
//   - error: 错误信息
func runInteractiveMode(conn net.Conn, config ClientConfig) (*TransferStats, error) {
	startTime := time.Now()
	stats := &TransferStats{StartTime: startTime}

	// 创建 readline 配置（不设置 HistoryFile，历史仅内存存储）
	rlConfig := &readline.Config{
		Prompt:          "tcp> ",
		AutoComplete:    buildCompleter(),
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	}

	// 创建 readline 实例
	rl, err := readline.NewEx(rlConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create readline: %w", err)
	}
	defer func() {
		if err := rl.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to close readline: %v\n", err)
		}
	}()

	// 创建会话
	session := &InteractiveSession{
		conn:         conn,
		config:       config,
		stats:        stats,
		rl:           rl,
		interval:     0,
		lastActivity: startTime,
	}

	// 显示欢迎信息
	printWelcome(conn)

	// 主循环
	for {
		line, err := rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				continue
			}
			if err == io.EOF {
				fmt.Println("\nExiting interactive mode...")
				break
			}
			return nil, err
		}

		// 处理输入
		if err := session.handleInput(line); err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	}

	stats.Duration = time.Since(startTime)
	return stats, nil
}

// printWelcome 显示欢迎信息
//
// 参数:
//   - conn: TCP 连接
func printWelcome(conn net.Conn) {
	fmt.Println("TCP Interactive Client")
	fmt.Printf("Connected to %s\n", conn.RemoteAddr())
	fmt.Println("Type /help for available commands")
	fmt.Println()
}

// handleInput 处理用户输入
//
// 参数:
//   - line: 输入行
//
// 返回值:
//   - error: 错误信息
func (s *InteractiveSession) handleInput(line string) error {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}

	// 添加到历史记录
	s.addHistory(line)

	// 检查是否是内置命令（以 / 开头）
	if strings.HasPrefix(line, "/") {
		return s.handleBuiltinCommand(line)
	}

	// 普通消息发送
	s.lastSend = line
	sent, received, err := s.sendMessage(line, false)
	s.stats.BytesSent += sent
	s.stats.BytesReceived += received
	return err
}

// addHistory 添加历史记录
//
// 参数:
//   - line: 输入行
func (s *InteractiveSession) addHistory(line string) {
	// 避免重复添加相同的连续记录
	if len(s.history) > 0 && s.history[len(s.history)-1] == line {
		return
	}
	s.history = append(s.history, line)
	// 限制历史记录数量（最多 1000 条）
	if len(s.history) > 1000 {
		s.history = s.history[1:]
	}
}

// sendMessage 发送消息
//
// 参数:
//   - msg: 消息内容
//   - silent: 是否静默模式（不输出响应）
//
// 返回值:
//   - int64: 发送的字节数
//   - int64: 接收的字节数
//   - error: 错误信息
func (s *InteractiveSession) sendMessage(msg string, silent bool) (int64, int64, error) {
	data := []byte(msg + "\n")
	n, err := s.conn.Write(data)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to send message: %w", err)
	}
	sentBytes := int64(n)

	// 更新消息计数和最后活动时间
	s.msgCount++
	s.lastActivity = time.Now()

	// 接收响应
	var receivedBytes int64
	if !s.config.NoResponse {
		var err error
		receivedBytes, err = s.receiveResponse(silent)
		if err != nil {
			return sentBytes, 0, err
		}
	}

	return sentBytes, receivedBytes, nil
}

// receiveResponse 接收服务器响应
//
// 参数:
//   - silent: 是否静默模式（不输出响应）
//
// 返回值:
//   - int64: 接收的字节数
//   - error: 错误信息
func (s *InteractiveSession) receiveResponse(silent bool) (int64, error) {
	// 设置读取超时
	if s.config.Timeout > 0 {
		if err := s.conn.SetReadDeadline(time.Now().Add(s.config.Timeout)); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to set read deadline: %v\n", err)
		}
	}

	// 单次读取响应
	buffer := make([]byte, s.config.BufferSize)
	n, err := s.conn.Read(buffer)
	if err != nil {
		if err != io.EOF {
			fmt.Fprintf(os.Stderr, "Warning: failed to read response: %v\n", err)
		}
	} else if n > 0 {
		s.stats.BytesReceived += int64(n)
		if !silent {
			fmt.Printf("└ %s\n\n", string(buffer[:n]))
		}
		return int64(n), nil
	}

	// 清除读取超时
	if err := s.conn.SetReadDeadline(time.Time{}); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to clear read deadline: %v\n", err)
	}

	return int64(n), nil
}

// buildCompleter 构建自动补全
//
// 返回值:
//   - *readline.PrefixCompleter: 补全器
func buildCompleter() *readline.PrefixCompleter {
	return readline.NewPrefixCompleter(
		readline.PcItem("/quit",
			readline.PcItem("-h"),
		),
		readline.PcItem("/history"),
		readline.PcItem("/clear"),
		readline.PcItem("/info"),
		readline.PcItem("/ping"),
		readline.PcItem("/help"),
		readline.PcItem("/send",
			readline.PcItemDynamic(listFiles),
		),
		readline.PcItem("/hex"),
		readline.PcItem("/interval"),
		readline.PcItem("/repeat"),
	)
}

// listFiles 列出当前目录文件（用于自动补全）
//
// 参数:
//   - prefix: 前缀
//
// 返回值:
//   - []string: 文件列表
func listFiles(prefix string) []string {
	files, _ := filepath.Glob(prefix + "*")
	return files
}
