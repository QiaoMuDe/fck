package tcp

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
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

	// 根据模式执行发送
	var stats *TransferStats
	switch {
	case config.Message != "":
		stats, err = sendString(conn, config)
	case config.Path != "":
		stats, err = sendPath(conn, config)
	case config.Interactive:
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

// sendPath 发送路径（文件、目录或通配符）
//
// 参数:
//   - conn: TCP 连接
//   - config: 客户端配置
//
// 返回值:
//   - *TransferStats: 传输统计
//   - error: 错误
func sendPath(conn net.Conn, config ClientConfig) (*TransferStats, error) {
	// 解析路径
	files, err := resolvePath(config.Path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no files found for path: %s", config.Path)
	}

	startTime := time.Now()
	stats := &TransferStats{}

	// 发送每个文件，文件间用换行分隔
	for i, file := range files {
		// 发送文件内容
		sent, err := sendRawFileData(conn, file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to send %s: %v\n", file, err)
			continue
		}
		stats.BytesSent += sent
		stats.FilesSent++

		// 除了最后一个文件，都追加换行符
		if i < len(files)-1 {
			if _, err := conn.Write([]byte("\n")); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to send newline: %v\n", err)
			} else {
				stats.BytesSent += 1
			}
		}
	}

	stats.Duration = time.Since(startTime)
	return stats, nil
}

// resolvePath 解析路径（支持文件、目录、通配符）
//
// 参数:
//   - path: 路径字符串
//
// 返回值:
//   - []string: 文件列表
//   - error: 错误
func resolvePath(path string) ([]string, error) {
	// 检查是否是通配符
	if strings.Contains(path, "*") || strings.Contains(path, "?") {
		matches, err := filepath.Glob(path)
		if err != nil {
			return nil, fmt.Errorf("invalid glob pattern: %w", err)
		}
		return matches, nil
	}

	// 检查路径是否存在
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to access path: %w", err)
	}

	// 如果是目录，读取目录下所有文件（不递归）
	if info.IsDir() {
		entries, err := os.ReadDir(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read directory: %w", err)
		}

		var files []string
		for _, entry := range entries {
			if !entry.IsDir() {
				files = append(files, filepath.Join(path, entry.Name()))
			}
		}
		return files, nil
	}

	// 普通文件
	return []string{path}, nil
}

// sendRawFileData 发送文件原始字节（无协议头）
//
// 参数:
//   - conn: TCP 连接
//   - filePath: 文件路径
//
// 返回值:
//   - int64: 发送字节数
//   - error: 错误
func sendRawFileData(conn net.Conn, filePath string) (int64, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close file: %v\n", err)
		}
	}()

	// 直接发送文件内容（无协议头）
	sent, err := io.Copy(conn, file)
	if err != nil {
		return sent, fmt.Errorf("failed to send file content: %w", err)
	}

	return sent, nil
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

		// 接收响应
		if !config.NoResponse {
			response, err := readResponse(conn, config.Timeout, config.BufferSize)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to read response: %v\n", err)
			} else {
				stats.BytesReceived += int64(len(response))
				fmt.Printf("< %s\n", string(response))
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
