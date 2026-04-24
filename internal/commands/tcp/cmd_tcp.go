// Package tcp 实现了 TCP 连接测试和端口扫描功能
package tcp

import (
	"encoding/json"
	"fmt"
	"math"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

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
func runConnectMode(config TcpConfig, ipAddr net.IP) error {
	if config.Port < 1 || config.Port > 65535 {
		return fmt.Errorf("invalid port: %d, must be 1-65535", config.Port)
	}

	stats := &ConnectionStats{
		MinTime: time.Duration(math.MaxInt64),
	}

	target := net.JoinHostPort(ipAddr.String(), strconv.Itoa(config.Port))

	if !config.Quiet && !config.Json {
		fmt.Printf("Connecting to %s (%s):%d...\n", config.Host, ipAddr.String(), config.Port)
	}

	count := config.Count
	if count == 0 {
		count = 1
	}

	for i := 0; i < count; i++ {
		start := time.Now()
		conn, err := net.DialTimeout("tcp", target, config.Timeout)
		duration := time.Since(start)

		stats.Attempted++
		stats.TotalTime += duration

		if err != nil {
			stats.Failed++
			if !config.Quiet && !config.Json {
				fmt.Printf("Connection %d failed: %v\n", i+1, err)
			}
		} else {
			stats.Succeeded++
			if duration < stats.MinTime {
				stats.MinTime = duration
			}
			if duration > stats.MaxTime {
				stats.MaxTime = duration
			}
			if !config.Quiet && !config.Json {
				fmt.Printf("Connected to %s:%d (time=%.3fms)\n", config.Host, config.Port, float64(duration.Microseconds())/1000.0)
			}
			_ = conn.Close()
		}

		if i < count-1 && config.Interval > 0 {
			time.Sleep(config.Interval)
		}
	}

	if stats.MinTime == time.Duration(1<<63-1) {
		stats.MinTime = 0
	}

	if config.Json {
		return outputConnectJSON(config, stats)
	}

	if !config.Quiet {
		printConnectStats(config, stats)
	}

	return nil
}

// runScanMode 运行端口扫描模式
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
	results := scanPorts(config.Host, ipAddr, ports, config)
	duration := time.Since(start)

	if config.Json {
		return outputScanJSON(config, ipAddr, results, duration)
	}

	printScanResults(config, results, duration)
	return nil
}

// runBannerMode 运行 Banner 获取模式
func runBannerMode(config TcpConfig, ipAddr net.IP) error {
	if config.Port < 1 || config.Port > 65535 {
		return fmt.Errorf("invalid port: %d, must be 1-65535", config.Port)
	}

	target := net.JoinHostPort(ipAddr.String(), strconv.Itoa(config.Port))

	if !config.Quiet {
		fmt.Printf("Connecting to %s (%s):%d...\n", config.Host, ipAddr.String(), config.Port)
	}

	conn, err := net.DialTimeout("tcp", target, config.Timeout)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer func() { _ = conn.Close() }()

	if !config.Quiet {
		fmt.Printf("Connected\n")
	}

	// 如果指定了等待时间，读取响应
	if config.Wait > 0 {
		_ = conn.SetReadDeadline(time.Now().Add(config.Wait))
		buffer := make([]byte, 4096)
		n, err := conn.Read(buffer)
		if err == nil && n > 0 {
			if !config.Quiet {
				fmt.Printf("\n--- Received Banner ---\n%s\n", string(buffer[:n]))
			}
		}
	} else if !config.Quiet {
		fmt.Printf("\nTip: use -w/--wait to specify time to wait for banner\n")
	}

	return nil
}

// runDataMode 运行数据发送模式
func runDataMode(config TcpConfig, ipAddr net.IP) error {
	if config.Port < 1 || config.Port > 65535 {
		return fmt.Errorf("invalid port: %d, must be 1-65535", config.Port)
	}

	var data []byte
	if config.File != "" {
		content, err := os.ReadFile(config.File)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		data = content
	} else {
		data = []byte(config.Data)
	}

	target := net.JoinHostPort(ipAddr.String(), strconv.Itoa(config.Port))

	if !config.Quiet {
		fmt.Printf("Connecting to %s (%s):%d...\n", config.Host, ipAddr.String(), config.Port)
	}

	conn, err := net.DialTimeout("tcp", target, config.Timeout)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer func() { _ = conn.Close() }()

	if !config.Quiet {
		fmt.Printf("Connected, sending %d bytes...\n", len(data))
	}

	// 发送数据
	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send data: %w", err)
	}

	// 等待响应
	if config.Wait > 0 {
		_ = conn.SetReadDeadline(time.Now().Add(config.Wait))
		buffer := make([]byte, 4096)
		n, err := conn.Read(buffer)
		if err == nil && n > 0 {
			if !config.Quiet {
				fmt.Printf("\n--- Received Response (%d bytes) ---\n%s\n", n, string(buffer[:n]))
			}
		} else if err != nil && !config.Quiet {
			fmt.Printf("No response received\n")
		}
	}

	return nil
}

// scanPorts 扫描端口
func scanPorts(host string, ipAddr net.IP, ports []int, config TcpConfig) []PortResult {
	var wg sync.WaitGroup
	results := make([]PortResult, len(ports))

	// 限制并发数
	semaphore := make(chan struct{}, 100)

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

// parsePortRange 解析端口范围
// 支持格式: "80", "80-100", "80,443,8080", "1-100,443,8080-8090"
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

// guessService 根据端口猜测服务
func guessService(port int) string {
	services := map[int]string{
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

	if service, ok := services[port]; ok {
		return service
	}
	return "unknown"
}

// printConnectStats 打印连接统计
func printConnectStats(config TcpConfig, stats *ConnectionStats) {
	fmt.Printf("\n--- %s:%d tcp statistics ---\n", config.Host, config.Port)
	fmt.Printf("%d connections attempted, %d succeeded, %d failed\n",
		stats.Attempted, stats.Succeeded, stats.Failed)

	if stats.Succeeded > 0 {
		avgTime := stats.TotalTime / time.Duration(stats.Succeeded)
		fmt.Printf("rtt min/avg/max = %.3f/%.3f/%.3f ms\n",
			float64(stats.MinTime.Microseconds())/1000.0,
			float64(avgTime.Microseconds())/1000.0,
			float64(stats.MaxTime.Microseconds())/1000.0)
	}
}

// printScanResults 打印扫描结果
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

// outputConnectJSON 输出连接测试 JSON 结果
func outputConnectJSON(config TcpConfig, stats *ConnectionStats) error {
	avgTime := time.Duration(0)
	if stats.Succeeded > 0 {
		avgTime = stats.TotalTime / time.Duration(stats.Succeeded)
	}

	result := map[string]interface{}{
		"host":      config.Host,
		"port":      config.Port,
		"attempted": stats.Attempted,
		"succeeded": stats.Succeeded,
		"failed":    stats.Failed,
		"rtt_ms": map[string]float64{
			"min": float64(stats.MinTime.Microseconds()) / 1000.0,
			"avg": float64(avgTime.Microseconds()) / 1000.0,
			"max": float64(stats.MaxTime.Microseconds()) / 1000.0,
		},
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

// outputScanJSON 输出扫描 JSON 结果
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
