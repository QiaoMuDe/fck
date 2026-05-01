package tcp

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync/atomic"
	"syscall"
	"time"
)

// connectionStats 连接统计
type connectionStats struct {
	active int32 // 当前活跃连接数
	total  int32 // 总连接数
}

var stats connectionStats

// logType 日志类型
type logType string

const (
	logConn  logType = "CONN" // 连接建立
	logRecv  logType = "RECV" // 接收数据
	logSend  logType = "SEND" // 发送数据
	logDisc  logType = "DISC" // 断开连接
	logError logType = "ERR"  // 错误
)

// timeFormat 时间格式常量
const timeFormat = "2006-01-02 15:04:05"

// serverLogger 服务端日志记录器
type serverLogger struct {
	clientAddr string
}

// newServerLogger 创建新的日志记录器
//
// 参数:
//   - clientAddr: 客户端地址
//
// 返回值:
//   - *serverLogger: 日志记录器
func newServerLogger(clientAddr string) *serverLogger {
	return &serverLogger{clientAddr: clientAddr}
}

// log 打印日志
//
// 参数:
//   - logType: 日志类型
//   - detail: 详情
func (l *serverLogger) log(logType logType, detail string) {
	timestamp := time.Now().Format(timeFormat)
	fmt.Printf("[%s] [%s] [%s] %s\n", timestamp, logType, l.clientAddr, detail)
}

// logConn 记录连接建立
func (l *serverLogger) logConn() {
	active := atomic.AddInt32(&stats.active, 1)
	total := atomic.AddInt32(&stats.total, 1)
	l.log(logConn, fmt.Sprintf("Connected (active: %d, total: %d)", active, total))
}

// logDisc 记录断开连接
//
// 参数:
//   - duration: 连接持续时间
func (l *serverLogger) logDisc(duration time.Duration) {
	active := atomic.AddInt32(&stats.active, -1)
	l.log(logDisc, fmt.Sprintf("Disconnected (duration: %.3fs, active: %d)", duration.Seconds(), active))
}

// logRecv 记录接收数据
//
// 参数:
//   - data: 接收的数据
//   - size: 数据大小
func (l *serverLogger) logRecv(data []byte, size int) {
	l.log(logRecv, fmt.Sprintf("Received data, size: %d bytes", size))
	content := formatDataContent(data)
	fmt.Printf("%s\n", content)
}

// logSend 记录发送数据
//
// 参数:
//   - respType: 响应类型
//   - size: 数据大小
func (l *serverLogger) logSend(respType string, size int) {
	l.log(logSend, fmt.Sprintf("Sent %s response (%d bytes)", respType, size))
}

// logError 记录错误
//
// 参数:
//   - err: 错误信息
func (l *serverLogger) logError(err string) {
	l.log(logError, err)
}

// formatDataContent 格式化数据内容用于显示
//
// 参数:
//   - data: 原始数据
//
// 返回值:
//   - string: 格式化后的字符串
func formatDataContent(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	const maxLen = 100

	// 尝试作为UTF-8解码
	str := string(data)
	if len(str) > maxLen {
		return str[:maxLen] + fmt.Sprintf("... (total: %d bytes)", len(data))
	}
	return str
}

// ServerCmdMain TCP 服务端主入口
//
// 参数:
//   - config: 服务端配置
//
// 返回值:
//   - error: 执行错误
func ServerCmdMain(config ServerConfig) error {
	// 创建可取消的上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 监听系统信号（Ctrl+C）
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 创建监听器
	address := fmt.Sprintf("%s:%d", config.Address, config.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to start server on %s: %w", address, err)
	}
	defer func() {
		if err := listener.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close listener: %v\n", err)
		}
	}()

	fmt.Printf("TCP Server listening on %s\n", address)
	fmt.Printf("Press Ctrl+C to stop\n\n")

	// 创建连接限制信号量
	sem := make(chan struct{}, config.MaxConn)

	// 在后台监听信号
	go func() {
		<-sigChan
		fmt.Println("\nReceived shutdown signal, stopping server...")
		cancel()
		if err := listener.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close listener: %v\n", err)
		}
	}()

	// 接受连接循环
	for {
		select {
		case <-ctx.Done():
			// 等待所有连接处理完成
			for i := 0; i < config.MaxConn; i++ {
				sem <- struct{}{}
			}
			fmt.Println("Server stopped gracefully")
			return nil
		default:
		}

		// 设置接受超时，以便定期检查 ctx.Done()
		if tcpListener, ok := listener.(*net.TCPListener); ok {
			if err := tcpListener.SetDeadline(time.Now().Add(1 * time.Second)); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to set deadline: %v\n", err)
				continue
			}
		}

		conn, err := listener.Accept()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				// 超时，继续循环检查 ctx.Done()
				continue
			}
			select {
			case <-ctx.Done():
				return nil
			default:
				fmt.Fprintf(os.Stderr, "Failed to accept connection: %v\n", err)
				continue
			}
		}

		// 获取信号量（限制并发）
		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			if err := conn.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to close connection: %v\n", err)
			}
			return nil
		}

		// 处理连接
		go func(c net.Conn) {
			defer func() { <-sem }()
			handleConnection(c, config)
		}(conn)
	}
}

// handleConnection 处理单个客户端连接
//
// 参数:
//   - conn: 客户端连接
//   - config: 服务端配置
func handleConnection(conn net.Conn, config ServerConfig) {
	clientAddr := conn.RemoteAddr().String()
	startTime := time.Now()
	logger := newServerLogger(clientAddr)

	// 确保连接最终被关闭
	defer func() {
		if err := conn.Close(); err != nil {
			logger.logError(fmt.Sprintf("Failed to close connection: %v", err))
		}
		// 记录断开连接
		logger.logDisc(time.Since(startTime))
	}()

	// 记录连接建立
	logger.logConn()

	// 使用 bufio.Reader 按行读取，支持即时响应
	reader := bufio.NewReader(conn)

	for {
		// 读取一行数据（以 \n 为分隔符）
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				logger.logError(err.Error())
			}
			// 处理最后可能剩余的数据（没有换行符结尾）
			if len(line) > 0 {
				handleMessage(conn, line, config, logger, startTime)
			}
			break
		}

		// 立即处理收到的消息并响应
		handleMessage(conn, line, config, logger, startTime)
	}
}

// handleMessage 处理单条消息并发送响应
//
// 参数:
//   - conn: 客户端连接
//   - data: 接收到的数据
//   - config: 服务端配置
//   - logger: 日志记录器
//   - connStartTime: 连接开始时间
func handleMessage(conn net.Conn, data string, config ServerConfig, logger *serverLogger, connStartTime time.Time) {
	// 去除换行符
	data = strings.TrimRight(data, "\r\n")

	// 空消息不处理
	if len(data) == 0 {
		return
	}

	// 记录接收的数据
	logger.logRecv([]byte(data), len(data))

	// 保存到文件（如果配置了输出目录）
	if config.OutputDir != "" {
		saveReceivedData(config.OutputDir, conn.RemoteAddr().String(), data, connStartTime)
	}

	// 构建并发送响应
	var response string
	if config.Echo {
		// Echo 模式：返回接收到的数据
		response = formatResponse("ECHO", len(data))
	} else if config.Response != "" {
		// 自定义响应内容
		response = config.Response
	} else {
		// 默认响应格式
		response = formatResponse("ACK", len(data))
	}

	_, err := conn.Write([]byte(response))
	if err != nil {
		logger.logError(fmt.Sprintf("Failed to send response: %v", err))
	} else {
		// 记录发送的数据
		respType := "ACK"
		if config.Echo {
			respType = "ECHO"
		}
		logger.logSend(respType, len(response))
	}
}

// saveReceivedData 保存接收到的数据（所有连接每小时一个文件）
//
// 参数:
//   - outputDir: 输出目录
//   - clientAddr: 客户端地址
//   - data: 接收到的数据
//   - connStartTime: 连接开始时间（用于确定文件）
func saveReceivedData(outputDir string, clientAddr string, data string, connStartTime time.Time) {
	// 确保输出目录存在
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create output directory: %v\n", err)
		return
	}

	// 生成文件名：tcp_年月日小时.txt（所有连接共享）
	hourTimestamp := connStartTime.Format("2006010215")
	filename := fmt.Sprintf("tcp_%s.txt", hourTimestamp)
	filePath := filepath.Join(outputDir, filename)

	// 追加写入原始数据（带客户端地址前缀）
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close file: %v\n", err)
		}
	}()

	// 写入数据：[客户端地址] 数据 + 换行符 + 空行
	record := fmt.Sprintf("[%s] %s\n\n", clientAddr, data)
	if _, err := file.WriteString(record); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to file: %v\n", err)
		return
	}
}

// formatResponse 格式化标准响应
//
// 参数:
//   - respType: 响应类型 (ACK, ECHO)
//   - bytes: 接收字节数
//
// 返回值:
//   - string: 格式化后的响应
//
// 格式: [YYYY-MM-DD HH:MM:SS] [OK] [TYPE] Received X bytes\ncontent
func formatResponse(respType string, bytes int) string {
	timestamp := time.Now().Format(timeFormat)
	return fmt.Sprintf("[%s] [OK] [%s] Received %d bytes", timestamp, respType, bytes)
}
