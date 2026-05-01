package tcp

import (
	"fmt"
	"io"
	"net"
	"os"
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

// sendInteractive 发送交互式消息（使用 readline）
//
// 参数:
//   - conn: TCP 连接
//   - config: 客户端配置
//
// 返回值:
//   - *TransferStats: 传输统计
//   - error: 错误
func sendInteractive(conn net.Conn, config ClientConfig) (*TransferStats, error) {
	return runInteractiveMode(conn, config)
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
