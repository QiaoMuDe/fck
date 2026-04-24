// Package tcp 实现了 TCP 连接测试和端口扫描功能
package tcp

import (
	"encoding/json"
	"fmt"
	"math"
	"net"
	"os"
	"os/signal"
	"runtime"
	"sort"
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
	Host     string        // 目标主机
	Port     int           // 目标端口
	Timeout  time.Duration // 连接超时时间
	Count    int           // 连接次数
	Interval time.Duration // 连接间隔
	Scan     bool          // 扫描模式
	Range    string        // 端口范围
	OpenOnly bool          // 仅显示开放端口
	Banner   bool          // 获取 banner
	Data     string        // 发送数据
	File     string        // 数据文件路径
	Wait     time.Duration // 等待响应时间
	Quiet    bool          // 静默模式
	Json     bool          // JSON 输出
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
	// 解析目标地址
	ipAddr, err := resolveHost(config.Host)
	if err != nil {
		return fmt.Errorf("failed to resolve host: %w", err)
	}

	// 根据模式执行不同逻辑
	if config.Scan || config.Range != "" {
		return runScanMode(config, ipAddr)
	}

	if config.Banner {
		return runBannerMode(config, ipAddr)
	}

	if config.Data != "" || config.File != "" {
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

	var data []byte
	var dataSource string
	if config.File != "" {
		content, err := os.ReadFile(config.File)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		data = content
		dataSource = config.File
	} else {
		data = []byte(config.Data)
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
			return outputDataJSON(config, ipAddr, data, dataSource, "", formatErrorState(err), dialDuration, 0)
		}
		return fmt.Errorf("connection failed: %w", err)
	}
	defer func() { _ = conn.Close() }()

	if !config.Quiet && !config.Json {
		fmt.Printf("Connected, sending %d bytes...\n", len(data))
	}

	// 发送数据
	sendStart := time.Now()
	_, err = conn.Write(data)
	if err != nil {
		if config.Json {
			return outputDataJSON(config, ipAddr, data, dataSource, "", "send_failed", dialDuration, time.Since(sendStart))
		}
		return fmt.Errorf("failed to send data: %w", err)
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
		return outputDataJSON(config, ipAddr, data, dataSource, response, "success", dialDuration, sendDuration)
	}

	return nil
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

// parsePortRange 解析端口范围字符串
// 支持多种格式,自动去重和排序
//
// 参数:
//   - rangeStr: 端口范围字符串
//
// 返回值:
//   - []int: 解析后的端口号列表(已排序)
//   - error: 解析错误,如格式无效、端口越界等
//
// 支持的格式:
//   - 单个端口: "80"
//   - 范围: "80-100"
//   - 多个端口: "80,443,8080"
//   - 混合: "1-100,443,8080-8090"
//
// 注意事项:
//   - 端口范围: 1-65535
//   - 自动去重
//   - 结果按升序排列
func parsePortRange(rangeStr string) ([]int, error) {
	if rangeStr == "" {
		return nil, fmt.Errorf("empty port range")
	}

	ports := make(map[int]bool)
	parts := strings.Split(rangeStr, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if strings.Contains(part, "-") {
			// 范围格式: 80-100
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid port range: %s", part)
			}

			start, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
			if err != nil {
				return nil, fmt.Errorf("invalid port number: %s", rangeParts[0])
			}

			end, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
			if err != nil {
				return nil, fmt.Errorf("invalid port number: %s", rangeParts[1])
			}

			if start < 1 || end > 65535 || start > end {
				return nil, fmt.Errorf("invalid port range: %d-%d", start, end)
			}

			for i := start; i <= end; i++ {
				ports[i] = true
			}
		} else {
			// 单个端口
			port, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid port number: %s", part)
			}
			if port < 1 || port > 65535 {
				return nil, fmt.Errorf("port out of range: %d", port)
			}
			ports[port] = true
		}
	}

	// 转换为切片
	result := make([]int, 0, len(ports))
	for port := range ports {
		result = append(result, port)
	}

	// 排序
	sort.Ints(result)

	return result, nil
}

// resolveHost 解析主机地址
// 支持IP地址和域名,优先返回IPv4地址
//
// 参数:
//   - host: 主机名或IP地址
//
// 返回值:
//   - net.IP: 解析后的IP地址
//   - error: 解析错误,如DNS解析失败等
//
// 解析逻辑:
//   - 如果是有效IP地址,直接返回
//   - 否则进行DNS解析
//   - 优先返回IPv4地址
//   - 如无IPv4,返回第一个可用地址
func resolveHost(host string) (net.IP, error) {
	// 尝试解析为 IP
	ip := net.ParseIP(host)
	if ip != nil {
		return ip, nil
	}

	// DNS 解析
	addrs, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}

	// 优先使用 IPv4
	for _, addr := range addrs {
		if addr.To4() != nil {
			return addr, nil
		}
	}

	// 如果没有 IPv4，使用第一个
	if len(addrs) > 0 {
		return addrs[0], nil
	}

	return nil, fmt.Errorf("no IP address found for host: %s", host)
}

// commonServices 常见端口到服务名称的映射
// 基于IANA标准端口分配和常见服务端口
// 作为包级变量避免每次调用重复分配内存
var commonServices = map[int]string{
	20:    "ftp-data",
	21:    "ftp",
	22:    "ssh",
	23:    "telnet",
	25:    "smtp",
	53:    "dns",
	67:    "dhcp",
	68:    "dhcp",
	80:    "http",
	110:   "pop3",
	111:   "rpcbind",
	113:   "ident",
	119:   "nntp",
	123:   "ntp",
	135:   "msrpc",
	139:   "netbios-ssn",
	143:   "imap",
	161:   "snmp",
	162:   "snmptrap",
	179:   "bgp",
	194:   "irc",
	389:   "ldap",
	443:   "https",
	445:   "smb",
	465:   "smtps",
	514:   "syslog",
	515:   "printer",
	587:   "submission",
	631:   "ipp",
	636:   "ldaps",
	873:   "rsync",
	989:   "ftps-data",
	990:   "ftps",
	993:   "imaps",
	995:   "pop3s",
	1080:  "socks",
	1194:  "openvpn",
	1433:  "mssql",
	1434:  "ms-sql-m",
	1521:  "oracle",
	1701:  "l2tp",
	1723:  "pptp",
	1883:  "mqtt",
	2049:  "nfs",
	2082:  "cpanel",
	2083:  "cpanel-ssl",
	2086:  "whm",
	2087:  "whm-ssl",
	2222:  "directadmin",
	2375:  "docker",
	2376:  "docker-ssl",
	3000:  "grafana",
	3306:  "mysql",
	3389:  "rdp",
	3690:  "svn",
	4333:  "mimer",
	4444:  "metasploit",
	4500:  "ipsec-nat-t",
	5000:  "upnp",
	5001:  "iperf",
	5060:  "sip",
	5061:  "sips",
	5432:  "postgresql",
	5900:  "vnc",
	5984:  "couchdb",
	5985:  "winrm",
	5986:  "winrm-ssl",
	6379:  "redis",
	6443:  "kubernetes",
	7001:  "weblogic",
	8000:  "http-alt",
	8008:  "http",
	8080:  "http-proxy",
	8086:  "influxdb",
	8443:  "https-alt",
	8888:  "http-alt",
	9000:  "sonarqube",
	9090:  "prometheus",
	9092:  "kafka",
	9200:  "elasticsearch",
	9300:  "elasticsearch-transport",
	9418:  "git",
	9999:  "abyss",
	10000: "webmin",
	11211: "memcached",
	27017: "mongodb",
	27018: "mongodb-shard",
	27019: "mongodb-config",
	28017: "mongodb-web",
	50000: "sap",
	50070: "hadoop-namenode",
}

// guessService 根据端口号猜测服务名称
// 基于IANA标准端口分配和常见服务端口
//
// 参数:
//   - port: 端口号
//
// 返回值:
//   - string: 服务名称,如"http"、"ssh"等;未知端口返回"unknown"
//
// 说明:
//   - 包含常见服务的知名端口映射
//   - 仅供参考,实际服务可能不同
//   - 涵盖Web、数据库、消息队列等常见服务
func guessService(port int) string {
	if service, ok := commonServices[port]; ok {
		return service
	}
	return "unknown"
}

// formatErrorState 格式化错误状态为简洁的字符串
//
// 参数:
//   - err: 连接错误
//
// 返回值:
//   - string: 状态字符串 (timeout/refused/filtered/error)
func formatErrorState(err error) string {
	errStr := err.Error()
	if strings.Contains(errStr, "timeout") || strings.Contains(errStr, "timed out") {
		return "timeout"
	}
	if strings.Contains(errStr, "refused") || strings.Contains(errStr, "connection refused") {
		return "refused"
	}
	if strings.Contains(errStr, "no such host") || strings.Contains(errStr, "i/o timeout") {
		return "timeout"
	}
	if strings.Contains(errStr, "network is unreachable") || strings.Contains(errStr, "no route to host") {
		return "unreachable"
	}
	return "error"
}

// printConnectStats 打印连接测试统计信息
// 显示连接次数、成功/失败数量、响应时间统计
//
// 参数:
//   - stats: 连接统计信息
//
// 输出格式:
//
//	--- Summary ---
//	N attempts, M success, K failed
//	min/avg/max = x.xxx/y.yyy/z.zzz ms
//
// 说明:
//   - 仅在多次连接时调用
//   - 仅统计成功的连接计算平均时间
func printConnectStats(stats *ConnectionStats) {
	fmt.Printf("\n--- Summary ---\n")
	fmt.Printf("%d attempts, %d success, %d failed\n",
		stats.Attempted, stats.Succeeded, stats.Failed)

	if stats.Succeeded > 0 {
		avgTime := stats.TotalTime / time.Duration(stats.Succeeded)
		fmt.Printf("min/avg/max = %.3f/%.3f/%.3f ms\n",
			float64(stats.MinTime.Microseconds())/1000.0,
			float64(avgTime.Microseconds())/1000.0,
			float64(stats.MaxTime.Microseconds())/1000.0)
	}
}

// printScanResults 打印端口扫描结果
// 以表格形式显示每个端口的状态、服务猜测和响应时间
//
// 参数:
//   - config: TCP配置,控制输出选项
//   - results: 端口扫描结果列表
//   - duration: 总扫描耗时
//
// 输出格式:
//
//	PORT     STATE    SERVICE    TIME(ms)
//	80       open     http       15.234
//	443      open     https      18.567
//
// 功能特性:
//   - 支持仅显示开放端口(-o选项)
//   - 自动统计开放/关闭端口数量
//   - 显示总扫描耗时
func printScanResults(config TcpConfig, results []PortResult, duration time.Duration) {
	if !config.Quiet {
		fmt.Printf("\nPORT     STATE    SERVICE    TIME(ms)\n")
	}

	open := 0
	closed := 0

	for _, result := range results {
		if result.State == "open" {
			open++
			if !config.OpenOnly {
				fmt.Printf("%-8d %-8s %-10s %.3f\n",
					result.Port, result.State, result.Service,
					result.TimeMs)
			}
		} else {
			closed++
			if !config.OpenOnly {
				fmt.Printf("%-8d %-8s %-10s %.3f\n",
					result.Port, result.State, result.Service,
					result.TimeMs)
			}
		}
	}

	if config.OpenOnly {
		for _, result := range results {
			if result.State == "open" {
				fmt.Printf("%-8d %-8s %-10s %.3f\n",
					result.Port, result.State, result.Service,
					result.TimeMs)
			}
		}
	}

	fmt.Printf("\nScan completed: %d ports scanned, %d open, %d closed/filtered\n",
		len(results), open, closed)
	fmt.Printf("Time taken: %.3fs\n", duration.Seconds())
}

// outputConnectJSON 以JSON格式输出连接测试结果
// 包含每次连接的详细结果和汇总统计
//
// 参数:
//   - config: TCP配置,包含目标主机和端口
//   - stats: 连接统计信息
//   - results: 每次连接的详细结果
//
// 返回值:
//   - error: JSON编码错误
//
// 输出示例:
//
//	{
//	  "host": "baidu.com",
//	  "port": 80,
//	  "results": [...],
//	  "summary": {...}
//	}
func outputConnectJSON(config TcpConfig, stats *ConnectionStats, results []ConnectResult) error {
	avgTime := time.Duration(0)
	if stats.Succeeded > 0 {
		avgTime = stats.TotalTime / time.Duration(stats.Succeeded)
	}

	result := map[string]interface{}{
		"host":    config.Host,
		"port":    config.Port,
		"results": results,
		"summary": map[string]interface{}{
			"attempted": stats.Attempted,
			"success":   stats.Succeeded,
			"failed":    stats.Failed,
			"min_ms":    float64(stats.MinTime.Microseconds()) / 1000.0,
			"avg_ms":    float64(avgTime.Microseconds()) / 1000.0,
			"max_ms":    float64(stats.MaxTime.Microseconds()) / 1000.0,
		},
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

// outputBannerJSON 以JSON格式输出Banner获取结果
//
// 参数:
//   - config: TCP配置,包含目标主机和端口
//   - ipAddr: 目标IP地址
//   - banner: 获取到的Banner内容
//   - state: 状态(success/timeout/refused等)
//   - duration: 连接耗时
//
// 返回值:
//   - error: JSON编码错误
func outputBannerJSON(config TcpConfig, ipAddr net.IP, banner string, state string, duration time.Duration) error {
	result := map[string]interface{}{
		"host":      config.Host,
		"ip":        ipAddr.String(),
		"port":      config.Port,
		"state":     state,
		"time_ms":   float64(duration.Microseconds()) / 1000.0,
		"banner":    banner,
		"wait_time": config.Wait.String(),
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

// outputDataJSON 以JSON格式输出数据发送结果
//
// 参数:
//   - config: TCP配置,包含目标主机和端口
//   - ipAddr: 目标IP地址
//   - data: 发送的数据内容
//   - dataSource: 数据来源(file路径或"inline")
//   - response: 接收到的响应
//   - state: 状态(success/send_failed/timeout等)
//   - dialDuration: 连接耗时
//   - sendDuration: 发送耗时
//
// 返回值:
//   - error: JSON编码错误
func outputDataJSON(config TcpConfig, ipAddr net.IP, data []byte, dataSource string, response string, state string, dialDuration time.Duration, sendDuration time.Duration) error {
	result := map[string]interface{}{
		"host":        config.Host,
		"ip":          ipAddr.String(),
		"port":        config.Port,
		"state":       state,
		"data_source": dataSource,
		"data_size":   len(data),
		"dial_ms":     float64(dialDuration.Microseconds()) / 1000.0,
		"send_ms":     float64(sendDuration.Microseconds()) / 1000.0,
		"wait_time":   config.Wait.String(),
		"response":    response,
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

// outputScanJSON 以JSON格式输出端口扫描结果
// 包含每个端口的详细结果和汇总统计
//
// 参数:
//   - config: TCP配置,包含目标主机
//   - ipAddr: 目标IP地址
//   - results: 端口扫描结果列表
//   - duration: 总扫描耗时
//
// 返回值:
//   - error: JSON编码错误
//
// 输出示例:
//
//	{
//	  "host": "target.com",
//	  "ip": "192.168.1.1",
//	  "ports": [...],
//	  "total": 100,
//	  "open": 3,
//	  ...
//	}
func outputScanJSON(config TcpConfig, ipAddr net.IP, results []PortResult, duration time.Duration) error {
	open := 0
	for _, r := range results {
		if r.State == "open" {
			open++
		}
	}

	result := ScanResult{
		Host:      config.Host,
		IP:        ipAddr.String(),
		Ports:     results,
		Total:     len(results),
		Open:      open,
		Closed:    len(results) - open,
		TimeTaken: fmt.Sprintf("%.3fs", duration.Seconds()),
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}
