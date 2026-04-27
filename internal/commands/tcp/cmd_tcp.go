// Package tcp 实现了 TCP 连接测试和端口扫描功能
package tcp

import (
	"context"
	"fmt"
	"math"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

// maxDuration 用于初始化 MinTime 的最大值
// 使用 math.MaxInt64 作为 time.Duration 的初始值,表示尚未设置
const maxDuration = time.Duration(math.MaxInt64)

// TcpConfig 配置结构体
type TcpConfig struct {
	Host        string        // 目标主机
	Port        int           // 目标端口
	Timeout     time.Duration // 连接超时时间
	Count       int           // 连接次数
	Interval    time.Duration // 连接间隔
	Scan        bool          // 扫描模式
	Range       string        // 端口范围
	OpenOnly    bool          // 仅显示开放端口
	Banner      bool          // 获取 banner
	Data        string        // 发送数据
	Path        string        // 文件/目录/通配符路径
	MaxFileSize int64         // 最大文件大小限制（字节）
	Wait        time.Duration // 等待响应时间
	Quiet       bool          // 静默模式
	Json        bool          // JSON 输出
	Listen      bool          // 监听模式
	Interactive bool          // 交互式模式
	Hex         bool          // 十六进制显示
}

// PortResult 端口扫描结果
type PortResult struct {
	Port    int     `json:"port"`    // 端口号
	State   string  `json:"state"`   // 状态: open/closed/filtered
	Service string  `json:"service"` // 服务名称
	TimeMs  float64 `json:"time_ms"` // 响应时间(毫秒)
}

// ScanResult 扫描结果
type ScanResult struct {
	Host      string       `json:"host"`       // 目标主机
	IP        string       `json:"ip"`         // IP 地址
	Ports     []PortResult `json:"ports"`      // 端口结果
	Total     int          `json:"total"`      // 扫描总数
	Open      int          `json:"open"`       // 开放数量
	Closed    int          `json:"closed"`     // 关闭数量
	TimeTaken string       `json:"time_taken"` // 耗时
}

// ConnectionStats 连接统计
type ConnectionStats struct {
	Attempted int           // 尝试连接数
	Succeeded int           // 成功连接数
	Failed    int           // 失败连接数
	MinTime   time.Duration // 最小耗时
	MaxTime   time.Duration // 最大耗时
	TotalTime time.Duration // 总耗时
}

// ConnectResult 单次连接结果
type ConnectResult struct {
	Seq    int     `json:"seq"`     // 序号
	State  string  `json:"state"`   // 状态
	TimeMs float64 `json:"time_ms"` // 耗时(毫秒)
}

// TcpCmdMain 执行 TCP 命令
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - error: 执行错误
func TcpCmdMain(config TcpConfig) error {
	// 监听模式不需要解析主机
	if config.Listen {
		return runListenMode(config)
	}

	// 解析目标地址
	ipAddr, err := resolveHost(config.Host)
	if err != nil {
		return fmt.Errorf("failed to resolve host: %w", err)
	}

	// 根据模式执行不同逻辑
	if config.Interactive {
		return runInteractiveMode(config, ipAddr)
	}

	if config.Scan || config.Range != "" {
		return runScanMode(config, ipAddr)
	}

	if config.Banner {
		return runBannerMode(config, ipAddr)
	}

	if config.Data != "" || config.Path != "" {
		return runDataMode(config, ipAddr)
	}

	// 默认连接测试模式
	return runConnectMode(config, ipAddr)
}

// runConnectMode 运行连接测试模式
// 该模式用于测试单个端口的连通性,支持单次或多次连接测试
//
// 参数:
//   - config: TCP配置,包含目标主机、端口、超时时间、连接次数等
//   - ipAddr: 解析后的目标IP地址
//
// 返回值:
//   - error: 执行过程中的错误,如端口无效、连接失败等
//
// 功能特性:
//   - 支持无限循环模式(-c 0),可通过Ctrl+C中断
//   - 支持自定义连接间隔
//   - 实时显示每次连接结果
//   - 多次连接时自动汇总统计信息
//   - 支持JSON格式输出
func runConnectMode(config TcpConfig, ipAddr net.IP) error {
	if config.Port < 1 || config.Port > 65535 {
		return fmt.Errorf("invalid port: %d, must be 1-65535", config.Port)
	}

	stats := &ConnectionStats{
		MinTime: maxDuration,
	}

	// 存储每次连接的结果用于 JSON 输出
	var results []ConnectResult

	target := net.JoinHostPort(ipAddr.String(), strconv.Itoa(config.Port))
	targetStr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	// 设置中断信号处理 (Ctrl+C)
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interruptChan)

	count := config.Count
	infiniteMode := count == 0

	for i := 0; infiniteMode || i < count; i++ {
		// 检查是否收到中断信号
		select {
		case <-interruptChan:
			goto finish
		default:
		}

		start := time.Now()
		conn, err := net.DialTimeout("tcp", target, config.Timeout)
		duration := time.Since(start)
		timeMs := float64(duration.Microseconds()) / 1000.0

		stats.Attempted++
		stats.TotalTime += duration

		var state string
		if err != nil {
			stats.Failed++
			state = formatErrorState(err)
		} else {
			stats.Succeeded++
			state = "open"
			if duration < stats.MinTime {
				stats.MinTime = duration
			}
			if duration > stats.MaxTime {
				stats.MaxTime = duration
			}
			_ = conn.Close()
		}

		// 记录结果
		results = append(results, ConnectResult{
			Seq:    i + 1,
			State:  state,
			TimeMs: timeMs,
		})

		// 打印单行结果
		if !config.Quiet && !config.Json {
			if infiniteMode {
				fmt.Printf("%s  %s  %.3fms  [#%d]\n", targetStr, state, timeMs, i+1)
			} else {
				fmt.Printf("%s  %s  %.3fms\n", targetStr, state, timeMs)
			}
		}

		// 非无限模式且不是最后一次，或者无限模式下，需要间隔
		if (!infiniteMode && i < count-1 || infiniteMode) && config.Interval > 0 {
			time.Sleep(config.Interval)
		}
	}

finish:

	if stats.MinTime == maxDuration {
		stats.MinTime = 0
	}

	if config.Json {
		return outputConnectJSON(config, stats, results)
	}

	// 只有多次连接时才打印统计信息
	if !config.Quiet && stats.Attempted > 1 {
		printConnectStats(stats)
	}

	return nil
}

// runScanMode 运行端口扫描模式
// 该模式用于批量扫描多个端口的状态,支持并发扫描
//
// 参数:
//   - config: TCP配置,包含目标主机、端口范围、超时时间等
//   - ipAddr: 解析后的目标IP地址
//
// 返回值:
//   - error: 执行过程中的错误,如端口范围无效等
//
// 功能特性:
//   - 支持多种端口范围格式(如: 80, 80-100, 80,443,8080)
//   - 自动去重和排序
//   - 并发扫描,并发数根据CPU核心数动态调整
//   - 支持仅显示开放端口
//   - 支持JSON格式输出
func runScanMode(config TcpConfig, ipAddr net.IP) error {
	ports, err := parsePortRange(config.Range)
	if err != nil {
		return fmt.Errorf("invalid port range: %w", err)
	}

	if len(ports) == 0 {
		return fmt.Errorf("no ports to scan")
	}

	if !config.Quiet && !config.Json {
		fmt.Printf("Scanning %s (%s) ports %s...\n", config.Host, ipAddr.String(), config.Range)
	}

	start := time.Now()
	results := scanPorts(ipAddr, ports, config)
	duration := time.Since(start)

	if config.Json {
		return outputScanJSON(config, ipAddr, results, duration)
	}

	printScanResults(config, results, duration)
	return nil
}

// runBannerMode 运行 Banner 获取模式
// 该模式用于连接目标端口并读取服务返回的欢迎信息(Banner)
//
// 参数:
//   - config: TCP配置,包含目标主机、端口、超时时间、等待时间等
//   - ipAddr: 解析后的目标IP地址
//
// 返回值:
//   - error: 执行过程中的错误,如连接失败等
//
// 功能特性:
//   - 建立TCP连接后不发送数据,仅接收
//   - 支持设置等待时间(-w),控制读取超时
//   - 常用于识别服务类型和版本
//   - 适用于SSH、SMTP等主动发送Banner的服务
func runBannerMode(config TcpConfig, ipAddr net.IP) error {
	if config.Port < 1 || config.Port > 65535 {
		return fmt.Errorf("invalid port: %d, must be 1-65535", config.Port)
	}

	target := net.JoinHostPort(ipAddr.String(), strconv.Itoa(config.Port))

	if !config.Quiet && !config.Json {
		fmt.Printf("Connecting to %s (%s):%d...\n", config.Host, ipAddr.String(), config.Port)
	}

	start := time.Now()
	conn, err := net.DialTimeout("tcp", target, config.Timeout)
	duration := time.Since(start)
	if err != nil {
		if config.Json {
			return outputBannerJSON(config, ipAddr, "", formatErrorState(err), duration)
		}
		return fmt.Errorf("connection failed: %w", err)
	}
	defer func() { _ = conn.Close() }()

	if !config.Quiet && !config.Json {
		fmt.Printf("Connected\n")
	}

	// 如果指定了等待时间，读取响应
	var banner string
	if config.Wait > 0 {
		_ = conn.SetReadDeadline(time.Now().Add(config.Wait))
		buffer := make([]byte, 4096)
		n, err := conn.Read(buffer)
		if err == nil && n > 0 {
			banner = string(buffer[:n])
			if !config.Quiet && !config.Json {
				fmt.Printf("\n--- Received Banner ---\n%s\n", banner)
			}
		}
	} else if !config.Quiet && !config.Json {
		fmt.Printf("\nTip: use -w/--wait to specify time to wait for banner\n")
	}

	if config.Json {
		return outputBannerJSON(config, ipAddr, banner, "success", duration)
	}

	return nil
}

// runDataMode 运行数据发送模式
// 该模式用于建立TCP连接后发送自定义数据,并可选接收响应
//
// 参数:
//   - config: TCP配置,包含目标主机、端口、超时时间、发送数据等
//   - ipAddr: 解析后的目标IP地址
//
// 返回值:
//   - error: 执行过程中的错误,如连接失败、文件读取失败等
//
// 功能特性:
//   - 支持从命令行参数(-d)直接发送字符串
//   - 支持从文件(-f)读取二进制或文本数据发送
//   - 支持设置等待时间(-w)接收响应
//   - 适用于测试自定义协议、发送HTTP请求等场景
func runDataMode(config TcpConfig, ipAddr net.IP) error {
	if config.Port < 1 || config.Port > 65535 {
		return fmt.Errorf("invalid port: %d, must be 1-65535", config.Port)
	}

	// 解析路径，支持文件、目录、通配符
	var files []string
	var dataSource string
	if config.Path != "" {
		pathInfo, err := resolvePath(config.Path, config.MaxFileSize)
		if err != nil {
			return err
		}
		files = pathInfo.Paths
		dataSource = config.Path
	} else {
		// 直接发送数据内容
		files = nil
		dataSource = "inline"
	}

	target := net.JoinHostPort(ipAddr.String(), strconv.Itoa(config.Port))

	if !config.Quiet && !config.Json {
		fmt.Printf("Connecting to %s (%s):%d...\n", config.Host, ipAddr.String(), config.Port)
	}

	start := time.Now()
	conn, err := net.DialTimeout("tcp", target, config.Timeout)
	dialDuration := time.Since(start)
	if err != nil {
		if config.Json {
			return outputDataJSON(config, ipAddr, []byte(config.Data), dataSource, "", formatErrorState(err), dialDuration, 0)
		}
		return fmt.Errorf("connection failed: %w", err)
	}
	defer func() { _ = conn.Close() }()

	// 发送数据并统计耗时
	sendStart := time.Now()

	if config.Data != "" {
		// 发送内联数据
		if err := sendData(conn, config, []byte(config.Data), "inline"); err != nil {
			return err
		}
	}

	// 发送文件
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			if !config.Quiet {
				fmt.Printf("Warning: failed to read file %s: %v\n", file, err)
			}
			continue
		}
		if err := sendFileData(conn, config, data, file); err != nil {
			return err
		}
	}

	sendDuration := time.Since(sendStart)

	// 等待响应
	var response string
	if config.Wait > 0 {
		_ = conn.SetReadDeadline(time.Now().Add(config.Wait))
		buffer := make([]byte, 4096)
		n, err := conn.Read(buffer)
		if err == nil && n > 0 {
			response = string(buffer[:n])
			if !config.Quiet && !config.Json {
				fmt.Printf("\n--- Received Response (%d bytes) ---\n%s\n", n, response)
			}
		} else if err != nil && !config.Quiet && !config.Json {
			fmt.Printf("No response received\n")
		}
	}

	if config.Json {
		return outputDataJSON(config, ipAddr, []byte(config.Data), dataSource, response, "success", dialDuration, sendDuration)
	}

	return nil
}

// sendData 发送数据
func sendData(conn net.Conn, config TcpConfig, data []byte, source string) error {
	if !config.Quiet && !config.Json {
		fmt.Printf("Sending %d bytes from %s...\n", len(data), source)
	}

	_, err := conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send data: %w", err)
	}

	return nil
}

// sendFileData 发送文件数据
func sendFileData(conn net.Conn, config TcpConfig, data []byte, filename string) error {
	if !config.Quiet && !config.Json {
		fmt.Printf("Sending file: %s (%d bytes)...\n", filename, len(data))
	}

	_, err := conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send file %s: %w", filename, err)
	}

	return nil
}

// runListenMode 运行监听模式
// 在指定端口监听TCP连接,接收并显示数据
//
// 参数:
//   - config: TCP配置,包含监听端口、超时时间等
//
// 返回值:
//   - error: 执行过程中的错误,如端口被占用等
//
// 功能特性:
//   - 支持多客户端并发连接
//   - 实时显示接收到的数据
//   - 显示连接和断开事件
//   - 支持Ctrl+C优雅退出
//   - 支持JSON格式输出
func runListenMode(config TcpConfig) error {
	if config.Port < 1 || config.Port > 65535 {
		return fmt.Errorf("invalid port: %d, must be 1-65535", config.Port)
	}

	addr := fmt.Sprintf(":%d", config.Port)

	// 创建监听器
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %w", config.Port, err)
	}
	defer func() { _ = listener.Close() }()

	if !config.Quiet && !config.Json {
		fmt.Printf("Listening on %s...\n", listener.Addr().String())
		fmt.Println("Press Ctrl+C to stop")
	}

	// 设置中断信号处理
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interruptChan)

	// 用于优雅关闭的context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 等待中断信号
	go func() {
		<-interruptChan
		if !config.Quiet && !config.Json {
			fmt.Println("\nShutting down...")
		}
		cancel()
		_ = listener.Close()
	}()

	// 连接计数器
	var connID int64

	for {
		// 使用带超时的accept来定期检查context
		_ = listener.(*net.TCPListener).SetDeadline(time.Now().Add(100 * time.Millisecond))
		conn, err := listener.Accept()

		select {
		case <-ctx.Done():
			return nil
		default:
		}

		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			if !config.Quiet && !config.Json {
				fmt.Printf("Accept error: %v\n", err)
			}
			continue
		}

		connID++
		currentID := connID

		// 处理连接
		go handleListenConnection(ctx, conn, currentID, config)
	}
}

// handleListenConnection 处理监听模式的单个连接
//
// 参数:
//   - ctx: 上下文,用于检测关闭信号
//   - conn: TCP连接
//   - connID: 连接ID
//   - config: TCP配置
func handleListenConnection(ctx context.Context, conn net.Conn, connID int64, config TcpConfig) {
	defer func() { _ = conn.Close() }()

	remoteAddr := conn.RemoteAddr().String()
	startTime := time.Now()

	if !config.Quiet && !config.Json {
		fmt.Printf("\n[%s] [#%d] Connection from %s\n", startTime.Format("2006-01-02 15:04:05"), connID, remoteAddr)
	}

	// 用于JSON输出的数据收集
	var receivedData []string
	totalBytes := 0

	// 读取数据循环
	buffer := make([]byte, 4096)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		// 设置读取超时
		if config.Wait > 0 {
			_ = conn.SetReadDeadline(time.Now().Add(config.Wait))
		}

		n, err := conn.Read(buffer)
		if n > 0 {
			data := buffer[:n]
			totalBytes += n

			if config.Json {
				receivedData = append(receivedData, string(data))
			} else if !config.Quiet {
				fmt.Printf("[%s] [#%d] Received %d bytes:\n%s\n",
					time.Now().Format("2006-01-02 15:04:05"),
					connID,
					n,
					string(data))
			}

			// 发送确认回应
			ackMsg := fmt.Sprintf("[+ACK] Received %d bytes at %s\n", n, time.Now().Format("2006-01-02 15:04:05"))
			_, writeErr := conn.Write([]byte(ackMsg))
			if writeErr != nil && !config.Quiet {
				fmt.Printf("[%s] [#%d] Write error: %v\n", time.Now().Format("2006-01-02 15:04:05"), connID, writeErr)
			}
		}

		if err != nil {
			if config.Json {
				_ = outputListenJSON(connID, remoteAddr, startTime, time.Since(startTime), totalBytes, receivedData, err != nil)
			} else if !config.Quiet {
				if err.Error() != "EOF" {
					fmt.Printf("[%s] [#%d] Read error: %v\n", time.Now().Format("2006-01-02 15:04:05"), connID, err)
				}
				fmt.Printf("[%s] [#%d] Connection closed (total: %d bytes, duration: %s)\n",
					time.Now().Format("2006-01-02 15:04:05"),
					connID,
					totalBytes,
					time.Since(startTime).Round(time.Millisecond))
			}
			return
		}
	}
}

// scanPorts 并发扫描多个端口
// 使用goroutine池控制并发数,提高扫描效率
//
// 参数:
//   - ipAddr: 目标IP地址
//   - ports: 待扫描的端口号列表
//   - config: TCP配置,包含超时时间等
//
// 返回值:
//   - []PortResult: 每个端口的扫描结果
//
// 并发控制:
//   - 并发数 = CPU核心数 * 2,最小为4
//   - 使用channel实现信号量控制
//   - 使用WaitGroup等待所有扫描完成
func scanPorts(ipAddr net.IP, ports []int, config TcpConfig) []PortResult {
	var wg sync.WaitGroup
	results := make([]PortResult, len(ports))

	// 限制并发数为 CPU 核心数 * 2
	concurrency := runtime.NumCPU() * 2
	if concurrency < 4 {
		concurrency = 4 // 最小并发数
	}
	semaphore := make(chan struct{}, concurrency)

	for i, port := range ports {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(idx, p int) {
			defer wg.Done()
			defer func() { <-semaphore }()

			results[idx] = scanPort(ipAddr, p, config.Timeout)
		}(i, port)
	}

	wg.Wait()
	return results
}

// scanPort 扫描单个端口
// 尝试建立TCP连接,根据结果判断端口状态
//
// 参数:
//   - ipAddr: 目标IP地址
//   - port: 目标端口号
//   - timeout: 连接超时时间
//
// 返回值:
//   - PortResult: 端口扫描结果,包含状态、服务猜测、响应时间等
//
// 状态说明:
//   - open: 连接成功,端口开放
//   - closed: 连接被拒绝,端口关闭
//   - filtered: 连接超时或被过滤
func scanPort(ipAddr net.IP, port int, timeout time.Duration) PortResult {
	target := net.JoinHostPort(ipAddr.String(), strconv.Itoa(port))
	start := time.Now()
	conn, err := net.DialTimeout("tcp", target, timeout)
	duration := time.Since(start)

	result := PortResult{
		Port:    port,
		Service: guessService(port),
		TimeMs:  float64(duration.Microseconds()) / 1000.0,
	}

	if err != nil {
		if strings.Contains(err.Error(), "refused") {
			result.State = "closed"
		} else {
			result.State = "filtered"
		}
	} else {
		result.State = "open"
		_ = conn.Close()
	}

	return result
}
