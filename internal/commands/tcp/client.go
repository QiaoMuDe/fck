package tcp

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"gitee.com/MM-Q/fck/internal/utils"
)

// ClientCmdMain TCP 客户端主入口
//
// 参数:
//   - config: 客户端配置
//
// 返回值:
//   - error: 执行错误
func ClientCmdMain(config ClientConfig) error {

	// 建立连接
	conn, err := net.DialTimeout("tcp", config.Address, config.Timeout)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", config.Address, err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close connection: %v\n", err)
		}
	}()

	// 根据模式执行发送（优先级：管道 > 字符串 > 交互式）
	var stats *TransferStats
	switch {
	case utils.IsStdinPipe():
		stats, err = sendStdin(conn, config)

	case config.Message != "":
		stats, err = sendString(conn, config)

	default:
		stats, err = sendInteractive(conn, config)
	}

	if err != nil {
		return err
	}

	// 输出统计
	printClientStats(stats)
	return nil
}

// sendString 发送字符串消息
//
// 参数:
//   - conn: TCP 连接
//   - config: 客户端配置
//
// 返回值:
//   - *TransferStats: 传输统计
//   - error: 错误
func sendString(conn net.Conn, config ClientConfig) (*TransferStats, error) {
	startTime := time.Now()
	stats := &TransferStats{}

	// 发送数据
	data := []byte(config.Message)
	n, err := conn.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}
	stats.BytesSent = int64(n)

	// 关闭写入端，通知服务端数据发送完毕
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		if err := tcpConn.CloseWrite(); err != nil {
			return nil, fmt.Errorf("failed to close write: %w", err)
		}
	}

	// 接收响应（如果不禁用）
	if !config.NoResponse {
		response, err := readResponse(conn, config.Timeout, config.BufferSize)
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}
		stats.BytesReceived = int64(len(response))
		fmt.Printf("Response: %s\n", string(response))
	}

	stats.Duration = time.Since(startTime)
	return stats, nil
}

// sendStdin 从标准输入（管道）读取并发送数据
//
// 参数:
//   - conn: TCP 连接
//   - config: 客户端配置
//
// 返回值:
//   - *TransferStats: 传输统计
//   - error: 错误
func sendStdin(conn net.Conn, config ClientConfig) (*TransferStats, error) {
	startTime := time.Now()
	stats := &TransferStats{}

	// 读取全部 stdin 内容
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("failed to read from stdin: %w", err)
	}

	// 发送数据
	n, err := conn.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed to send data: %w", err)
	}
	stats.BytesSent = int64(n)

	// 关闭写入端，通知服务端数据发送完毕
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		if err := tcpConn.CloseWrite(); err != nil {
			return nil, fmt.Errorf("failed to close write: %w", err)
		}
	}

	// 接收响应（如果不禁用）
	if !config.NoResponse {
		response, err := readResponse(conn, config.Timeout, config.BufferSize)
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}
		stats.BytesReceived = int64(len(response))
		fmt.Printf("Response: %s\n", string(response))
	}

	stats.Duration = time.Since(startTime)
	return stats, nil
}

// sendInteractive 交互式发送模式
//
// 参数:
//   - conn: TCP 连接
//   - config: 客户端配置
//
// 返回值:
//   - *TransferStats: 传输统计
//   - error: 错误
func sendInteractive(conn net.Conn, config ClientConfig) (*TransferStats, error) {
	startTime := time.Now()
	stats := &TransferStats{}
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Interactive mode started. Type your message and press Enter to send.")
	fmt.Printf("Use delimiter '%s' on a separate line to exit.\n", config.Delimiter)
	fmt.Println("----------------------------------------")

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to read input: %w", err)
		}

		// 去除换行符
		line = strings.TrimRight(line, "\r\n")

		// 检查分隔符
		if line == config.Delimiter {
			fmt.Println("Exiting interactive mode...")
			break
		}

		// 发送消息
		data := []byte(line + "\n")
		n, err := conn.Write(data)
		if err != nil {
			return nil, fmt.Errorf("failed to send message: %w", err)
		}
		stats.BytesSent += int64(n)

		// 接收响应（交互模式下只读取一次，不等待连接关闭）
		if !config.NoResponse {
			// 设置读取超时
			if config.Timeout > 0 {
				if err := conn.SetReadDeadline(time.Now().Add(config.Timeout)); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to set read deadline: %v\n", err)
				}
			}

			// 单次读取响应（不循环等待 EOF）
			buffer := make([]byte, config.BufferSize)
			n, err := conn.Read(buffer)
			if err != nil {
				if err != io.EOF {
					fmt.Fprintf(os.Stderr, "Warning: failed to read response: %v\n", err)
				}
			} else if n > 0 {
				stats.BytesReceived += int64(n)
				fmt.Printf("< %s\n\n", string(buffer[:n]))
			}

			// 清除读取超时，准备下一次交互
			if err := conn.SetReadDeadline(time.Time{}); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to clear read deadline: %v\n", err)
			}
		}
	}

	stats.Duration = time.Since(startTime)
	return stats, nil
}

// readResponse 读取服务器响应
//
// 参数:
//   - conn: TCP 连接
//   - timeout: 超时时间
//   - bufferSize: 缓冲区大小
//
// 返回值:
//   - []byte: 响应数据
//   - error: 错误
func readResponse(conn net.Conn, timeout time.Duration, bufferSize int) ([]byte, error) {
	// 设置读取超时
	if timeout > 0 {
		if err := conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
			return nil, fmt.Errorf("failed to set read deadline: %w", err)
		}
	}

	// 循环读取所有数据直到 EOF 或超时
	buffer := make([]byte, bufferSize)
	var response []byte

	for {
		n, err := conn.Read(buffer)
		if n > 0 {
			response = append(response, buffer[:n]...)
		}

		if err != nil {
			if err == io.EOF {
				// 服务端关闭连接，返回已读取的数据
				return response, nil
			}
			return nil, err
		}
	}
}

// printClientStats 打印客户端统计信息
//
// 参数:
//   - stats: 传输统计
func printClientStats(stats *TransferStats) {
	fmt.Printf("\nTransfer complete: %d bytes sent, %d bytes received, %d files, duration: %v\n",
		stats.BytesSent, stats.BytesReceived, stats.FilesSent, stats.Duration.Round(time.Millisecond))
}
