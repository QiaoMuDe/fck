package tcp

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// BuiltinCommand 内置命令定义
type BuiltinCommand struct {
	Name        string
	Aliases     []string
	Description string
	Handler     func(*InteractiveSession, []string) error
}

// builtinCommands 内置命令注册表
var builtinCommands = []*BuiltinCommand{
	{
		Name:        "quit",
		Aliases:     []string{"q", "exit"},
		Description: "退出交互模式",
		Handler:     cmdQuit,
	},
	{
		Name:        "history",
		Aliases:     []string{"h"},
		Description: "显示历史记录",
		Handler:     cmdHistory,
	},
	{
		Name:        "clear",
		Aliases:     []string{"c"},
		Description: "清屏",
		Handler:     cmdClear,
	},
	{
		Name:        "info",
		Aliases:     []string{"i"},
		Description: "显示连接信息和统计",
		Handler:     cmdInfo,
	},
	{
		Name:        "ping",
		Aliases:     []string{},
		Description: "发送测试消息并计算延迟",
		Handler:     cmdPing,
	},
	{
		Name:        "help",
		Aliases:     []string{"?"},
		Description: "显示帮助信息",
		Handler:     cmdHelp,
	},
	{
		Name:        "send",
		Aliases:     []string{},
		Description: "发送文件内容",
		Handler:     cmdSendFile,
	},
	{
		Name:        "hex",
		Aliases:     []string{},
		Description: "发送十六进制数据",
		Handler:     cmdSendHex,
	},
	{
		Name:        "interval",
		Aliases:     []string{"i"},
		Description: "设置发送间隔(毫秒)",
		Handler:     cmdInterval,
	},
	{
		Name:        "repeat",
		Aliases:     []string{"r"},
		Description: "重复发送上一次内容",
		Handler:     cmdRepeat,
	},
}

// handleBuiltinCommand 处理内置命令
//
// 参数:
//   - line: 输入行
//
// 返回值:
//   - error: 错误信息
func (s *InteractiveSession) handleBuiltinCommand(line string) error {
	// 解析命令和参数
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return nil
	}

	cmdName := strings.TrimPrefix(parts[0], "/")
	args := parts[1:]

	// 查找命令
	for _, cmd := range builtinCommands {
		if cmd.Name == cmdName || contains(cmd.Aliases, cmdName) {
			return cmd.Handler(s, args)
		}
	}

	return fmt.Errorf("unknown command: %s, type /help for available commands", cmdName)
}

// contains 检查字符串是否在切片中
//
// 参数:
//   - slice: 字符串切片
//   - str: 要查找的字符串
//
// 返回值:
//   - bool: 是否包含
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// cmdQuit 退出交互模式
//
// 参数:
//   - s: 交互式会话
//   - args: 参数
//
// 返回值:
//   - error: 错误信息
func cmdQuit(s *InteractiveSession, args []string) error {
	return io.EOF
}

// cmdHistory 显示历史记录
//
// 参数:
//   - s: 交互式会话
//   - args: 参数
//
// 返回值:
//   - error: 错误信息
func cmdHistory(s *InteractiveSession, args []string) error {
	if len(s.history) == 0 {
		fmt.Println("No history")
		fmt.Println()
		return nil
	}
	for i, h := range s.history {
		fmt.Printf("%3d: %s\n", i+1, h)
	}
	fmt.Println()
	return nil
}

// cmdClear 清屏
//
// 参数:
//   - s: 交互式会话
//   - args: 参数
//
// 返回值:
//   - error: 错误信息
func cmdClear(s *InteractiveSession, args []string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// cmdInfo 显示连接信息和统计
//
// 参数:
//   - s: 交互式会话
//   - args: 参数
//
// 返回值:
//   - error: 错误信息
func cmdInfo(s *InteractiveSession, args []string) error {
	// 使用 ping 检测连接状态
	connStatus := "OK"
	start := time.Now()
	_, _, err := s.sendMessage("PING", true) // 静默模式
	if err != nil {
		connStatus = "Error"
	}
	latency := time.Since(start)

	fmt.Printf("Local Address:    %s\n", s.conn.LocalAddr())
	fmt.Printf("Remote Address:   %s\n", s.conn.RemoteAddr())
	fmt.Printf("Connection:       %s (%v)\n", connStatus, latency)
	fmt.Printf("Duration:         %v\n", time.Since(s.stats.StartTime))
	fmt.Printf("Messages Sent:    %d\n", s.msgCount)
	fmt.Printf("Bytes Sent:       %d\n", s.stats.BytesSent)
	fmt.Printf("Bytes Received:   %d\n", s.stats.BytesReceived)
	fmt.Printf("Last Activity:    %v ago\n", time.Since(s.lastActivity).Round(time.Second))
	fmt.Printf("Interval:         %v\n", s.interval)
	fmt.Printf("Timeout:          %v\n", s.config.Timeout)
	fmt.Printf("Buffer Size:      %d\n", s.config.BufferSize)
	fmt.Println()
	return nil
}

// cmdPing 发送测试消息并计算延迟
//
// 参数:
//   - s: 交互式会话
//   - args: 参数
//
// 返回值:
//   - error: 错误信息
func cmdPing(s *InteractiveSession, args []string) error {
	start := time.Now()
	msg := "PING"
	sent, received, err := s.sendMessage(msg, false)
	if err != nil {
		return err
	}
	s.stats.BytesSent += sent
	s.stats.BytesReceived += received
	latency := time.Since(start)
	fmt.Printf("Latency: %v\n", latency)
	fmt.Println()
	return nil
}

// cmdHelp 显示帮助信息
//
// 参数:
//   - s: 交互式会话
//   - args: 参数
//
// 返回值:
//   - error: 错误信息
func cmdHelp(s *InteractiveSession, args []string) error {
	fmt.Println("Available commands:")
	fmt.Println()
	fmt.Println("  /quit (q)      - 退出交互模式")
	fmt.Println("  /history (h)   - 显示历史记录")
	fmt.Println("  /clear (c)     - 清屏")
	fmt.Println("  /info (i)      - 显示连接信息和统计")
	fmt.Println("  /ping          - 发送测试消息并计算延迟")
	fmt.Println("  /help (?)      - 显示帮助信息")
	fmt.Println("  /send <file>   - 发送文件内容")
	fmt.Println("  /hex <data>    - 发送十六进制数据")
	fmt.Println("  /interval <ms> - 设置发送间隔(毫秒)")
	fmt.Println("  /repeat <n>    - 重复发送上一次内容")
	fmt.Println()
	fmt.Println("Any other input will be sent as a message to the server.")
	fmt.Println()
	return nil
}

// cmdSendFile 发送文件内容
//
// 参数:
//   - s: 交互式会话
//   - args: 参数
//
// 返回值:
//   - error: 错误信息
func cmdSendFile(s *InteractiveSession, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: /send <file>")
	}

	filename := args[0]
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// 发送文件内容
	n, err := s.conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send file: %w", err)
	}
	s.stats.BytesSent += int64(n)
	fmt.Printf("Sent file %s (%d bytes)\n", filename, len(data))

	// 接收响应
	if !s.config.NoResponse {
		_, err := s.receiveResponse(false)
		if err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

// cmdSendHex 发送十六进制数据
//
// 参数:
//   - s: 交互式会话
//   - args: 参数
//
// 返回值:
//   - error: 错误信息
func cmdSendHex(s *InteractiveSession, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: /hex <hex_string>, e.g., /hex 48656c6c6f")
	}

	hexStr := strings.Join(args, "")
	// 去除可能的空格和0x前缀
	hexStr = strings.ReplaceAll(hexStr, " ", "")
	hexStr = strings.TrimPrefix(hexStr, "0x")

	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return fmt.Errorf("invalid hex string: %w", err)
	}

	// 发送十六进制解码后的数据
	n, err := s.conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send hex data: %w", err)
	}
	s.stats.BytesSent += int64(n)
	fmt.Printf("Sent hex data (%d bytes)\n", len(data))

	// 接收响应
	if !s.config.NoResponse {
		_, err := s.receiveResponse(false)
		if err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

// cmdInterval 设置发送间隔
//
// 参数:
//   - s: 交互式会话
//   - args: 参数
//
// 返回值:
//   - error: 错误信息
func cmdInterval(s *InteractiveSession, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: /interval <milliseconds>, e.g., /interval 1000")
	}

	ms, err := strconv.Atoi(args[0])
	if err != nil || ms < 0 {
		return fmt.Errorf("invalid interval: %s (must be a non-negative integer)", args[0])
	}

	s.interval = time.Duration(ms) * time.Millisecond
	fmt.Printf("Interval set to %v\n", s.interval)
	fmt.Println()
	return nil
}

// cmdRepeat 重复发送上一次内容
//
// 参数:
//   - s: 交互式会话
//   - args: 参数
//
// 返回值:
//   - error: 错误信息
func cmdRepeat(s *InteractiveSession, args []string) error {
	if s.lastSend == "" {
		return fmt.Errorf("no previous message to repeat, send a message first")
	}

	count := 1
	if len(args) > 0 {
		c, err := strconv.Atoi(args[0])
		if err != nil || c < 1 {
			return fmt.Errorf("invalid count: %s (must be a positive integer)", args[0])
		}
		count = c
	}

	// 批量发送时使用静默模式，只显示汇总
	var totalSent, totalReceived int64
	for i := 0; i < count; i++ {
		sent, received, err := s.sendMessage(s.lastSend, true) // 静默模式
		if err != nil {
			return err
		}
		totalSent += sent
		totalReceived += received
		s.stats.BytesSent += sent
		s.stats.BytesReceived += received
		if s.interval > 0 && i < count-1 {
			time.Sleep(s.interval)
		}
	}

	fmt.Printf("Sent %d messages, received %d responses\n", count, count)
	fmt.Printf("Total: %d bytes sent, %d bytes received\n", totalSent, totalReceived)
	fmt.Println()
	return nil
}
