package tcp

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"gitee.com/MM-Q/fck/internal/types"
	"github.com/jedib0t/go-pretty/v6/table"
)

// commonServices 常见端口服务映射
var commonServices = map[int]string{
	// 文件传输
	20: "ftp-data",
	21: "ftp",
	22: "ssh",
	23: "telnet",
	69: "tftp",
	// 邮件服务
	25:  "smtp",
	110: "pop3",
	143: "imap",
	465: "smtps",
	587: "submission",
	993: "imaps",
	995: "pop3s",
	// DNS
	53: "dns",
	// Web服务
	80:   "http",
	443:  "https",
	3000: "http-alt",
	3128: "squid-http",
	4443: "https-alt",
	5000: "http-alt",
	8000: "http-alt",
	8008: "http",
	8080: "http-proxy",
	8081: "http-alt",
	8443: "https-alt",
	8888: "http-alt",
	9000: "http-alt",
	9090: "http-alt",
	// 数据库
	1433:  "mssql",
	1521:  "oracle",
	1830:  "oracle",
	3306:  "mysql",
	5432:  "postgresql",
	6379:  "redis",
	6380:  "redis",
	27017: "mongodb",
	27018: "mongodb-shard",
	28017: "mongodb-web",
	// 远程桌面
	3389: "ms-wbt-server",
	5900: "vnc",
	5901: "vnc-1",
	5902: "vnc-2",
	5903: "vnc-3",
	6000: "x11",
	6001: "x11-1",
	// Windows服务
	135: "msrpc",
	139: "netbios-ssn",
	445: "microsoft-ds",
	// 网络服务
	111: "rpcbind",
	161: "snmp",
	162: "snmptrap",
	// VPN
	500:  "isakmp",
	1723: "pptp",
	4500: "ipsec-nat-t",
	1194: "openvpn",
	// 消息队列/缓存
	5672:  "amqp",
	5671:  "amqps",
	11211: "memcached",
	2181:  "zookeeper",
	9092:  "kafka",
	// 搜索引擎
	9200: "elasticsearch",
	9300: "elasticsearch-transport",
	// 版本控制
	9418: "git",
	3690: "svn",
	// 其他常见服务
	88:    "kerberos",
	389:   "ldap",
	636:   "ldaps",
	1080:  "socks",
	2049:  "nfs",
	3260:  "iscsi",
	4369:  "epmd",
	5222:  "xmpp",
	5269:  "xmpp-server",
	6667:  "irc",
	7001:  "weblogic",
	8009:  "ajp13",
	8161:  "activemq",
	10000: "webmin",
	27015: "steam",
	50070: "hdfs-namenode",
	50075: "hdfs-datanode",
	8088:  "yarn",
}

// ScanCmdMain 端口扫描主入口
//
// 参数:
//   - config: 扫描配置
//
// 返回值:
//   - error: 执行错误
func ScanCmdMain(config ScanConfig) error {
	// 解析端口范围
	ports, err := parsePortRange(config.Ports)
	if err != nil {
		return fmt.Errorf("invalid port range: %w", err)
	}

	// 解析目标地址
	target, err := resolveTarget(config.Target)
	if err != nil {
		return fmt.Errorf("failed to resolve target: %w", err)
	}

	// 执行扫描
	result := &ScanResult{
		Target:    target,
		StartTime: time.Now(),
		Results:   make([]PortResult, 0, len(ports)),
	}

	// 使用工作池控制并发
	results := scanPorts(target, ports, config)
	result.Results = results
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	// 统计结果
	result.Stats = calculateStats(result.Results)

	// 输出结果
	return outputScanResult(result, config)
}

// parsePortRange 解析端口范围字符串
//
// 参数:
//   - portRange: 端口范围字符串，如 "80", "1-1024", "80,443,8080"
//
// 返回值:
//   - []int: 端口列表
//   - error: 错误
func parsePortRange(portRange string) ([]int, error) {
	var ports []int
	portMap := make(map[int]bool)

	// 按逗号分割
	parts := strings.Split(portRange, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// 检查是否是范围
		if strings.Contains(part, "-") {
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
				if !portMap[i] {
					portMap[i] = true
					ports = append(ports, i)
				}
			}
		} else {
			// 单个端口
			port, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid port number: %s", part)
			}

			if port < 1 || port > 65535 {
				return nil, fmt.Errorf("port number out of range (1-65535): %d", port)
			}

			if !portMap[port] {
				portMap[port] = true
				ports = append(ports, port)
			}
		}
	}

	// 排序端口
	sort.Ints(ports)

	return ports, nil
}

// resolveTarget 解析目标地址
//
// 参数:
//   - target: 目标主机名或IP
//
// 返回值:
//   - string: 解析后的地址
//   - error: 错误
func resolveTarget(target string) (string, error) {
	// 检查是否是IP地址
	if net.ParseIP(target) != nil {
		return target, nil
	}

	// 尝试解析域名
	ips, err := net.LookupIP(target)
	if err != nil {
		return "", fmt.Errorf("failed to resolve hostname: %w", err)
	}

	if len(ips) == 0 {
		return "", fmt.Errorf("no IP addresses found for hostname: %s", target)
	}

	// 优先返回IPv4地址
	for _, ip := range ips {
		if ip.To4() != nil {
			return ip.String(), nil
		}
	}

	// 如果没有IPv4，返回第一个IPv6
	return ips[0].String(), nil
}

// scanPorts 并发扫描端口
//
// 参数:
//   - target: 目标地址
//   - ports: 端口列表
//   - config: 扫描配置
//
// 返回值:
//   - []PortResult: 扫描结果列表
func scanPorts(target string, ports []int, config ScanConfig) []PortResult {
	// 使用有缓冲通道作为工作队列
	portChan := make(chan int, config.Concurrent)
	resultChan := make(chan PortResult, len(ports))

	var wg sync.WaitGroup

	// 启动工作协程
	for i := 0; i < config.Concurrent; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range portChan {
				result := scanSinglePort(target, port, config.Timeout)
				resultChan <- result
			}
		}()
	}

	// 分发任务
	go func() {
		for _, port := range ports {
			portChan <- port
		}
		close(portChan)
	}()

	// 等待所有扫描完成
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 收集结果
	var results []PortResult
	completed := 0
	total := len(ports)

	for result := range resultChan {
		results = append(results, result)
		completed++

		// 如果启用了进度显示，实时输出结果
		if config.Progress {
			printProgressResult(result, completed, total, target)
		}
	}

	return results
}

// printProgressResult 打印进度结果
//
// 参数:
//   - result: 端口扫描结果
//   - completed: 已完成数量
//   - total: 总数
//   - target: 目标地址
func printProgressResult(result PortResult, completed, total int, target string) {
	// 只显示开放的端口，或者显示关闭的端口（如果配置了）
	// 但总是显示进度信息
	percent := float64(completed) * 100 / float64(total)

	// 使用 \r 回到行首，实现进度条效果
	fmt.Printf("\r[%3.0f%%] %d/%d ", percent, completed, total)

	// 如果端口是开放的，额外显示信息
	if result.Status == PortOpen {
		fmt.Printf("-> %s:%d %s %v", target, result.Port, result.Status, result.Response.Round(time.Millisecond))
	}

	// 如果全部完成，清除进度行
	if completed == total {
		// 使用 ANSI 转义序列清除当前行
		fmt.Printf("\r\033[K")
	}
}

// scanSinglePort 扫描单个端口
//
// 参数:
//   - target: 目标地址
//   - port: 端口号
//   - timeout: 超时时间
//
// 返回值:
//   - PortResult: 扫描结果
func scanSinglePort(target string, port int, timeout time.Duration) PortResult {
	address := net.JoinHostPort(target, strconv.Itoa(port))
	start := time.Now()

	conn, err := net.DialTimeout("tcp", address, timeout)
	duration := time.Since(start)

	result := PortResult{
		Port:     port,
		Response: duration,
	}

	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "timeout") || strings.Contains(errStr, "timed out") {
			result.Status = PortTimeout
		} else if strings.Contains(errStr, "refused") || strings.Contains(errStr, "connection refused") {
			result.Status = PortClosed
		} else {
			result.Status = PortFiltered
		}
		return result
	}

	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close connection: %v\n", err)
		}
	}()
	result.Status = PortOpen
	result.Service = guessService(port)
	return result
}

// guessService 猜测端口对应的服务
//
// 参数:
//   - port: 端口号
//
// 返回值:
//   - string: 服务名称
func guessService(port int) string {
	if service, ok := commonServices[port]; ok {
		return service
	}
	return "unknown"
}

// calculateStats 计算扫描统计
//
// 参数:
//   - results: 扫描结果列表
//
// 返回值:
//   - ScanStats: 统计信息
func calculateStats(results []PortResult) ScanStats {
	stats := ScanStats{
		Total: len(results),
	}

	for _, result := range results {
		switch result.Status {
		case PortOpen:
			stats.Open++
		case PortClosed:
			stats.Closed++
		case PortFiltered:
			stats.Filtered++
		case PortTimeout:
			stats.Timeout++
		}
	}

	return stats
}

// outputScanResult 输出扫描结果
//
// 参数:
//   - result: 扫描结果
//   - config: 扫描配置
//
// 返回值:
//   - error: 错误
func outputScanResult(result *ScanResult, config ScanConfig) error {
	// 根据格式输出
	switch config.Format {
	case types.TCPScanFormatTable:
		return outputTable(result, config)
	case types.TCPScanFormatJSON:
		return outputJSON(result, config)
	case types.TCPScanFormatCSV:
		return outputCSV(result, config)
	default:
		return outputDefault(result, config)
	}
}

// outputDefault 以默认简洁格式输出结果
//
// 参数:
//   - result: 扫描结果
//   - config: 扫描配置
//
// 返回值:
//   - error: 错误
func outputDefault(result *ScanResult, config ScanConfig) error {
	// 按端口排序
	sort.Slice(result.Results, func(i, j int) bool {
		return result.Results[i].Port < result.Results[j].Port
	})

	// 打印每个端口结果
	for _, r := range result.Results {
		// 只显示开放的端口，或者显示关闭的端口（如果配置了）
		if r.Status == PortOpen || config.ShowClosed {
			fmt.Printf("%s:%d  %s  %v\n", result.Target, r.Port, r.Status, r.Response.Round(time.Millisecond))
		}
	}

	// 只有在启用摘要且扫描多个端口时才打印统计行
	if config.Summary && result.Stats.Total > 1 {
		fmt.Printf("TOTAL:%d  OPEN:%d  %.3fs\n", result.Stats.Total, result.Stats.Open, result.Duration.Seconds())
	}

	// 扫描单个端口时，如果没有开放端口且没有配置显示关闭端口，也显示结果
	if result.Stats.Total == 1 && result.Stats.Open == 0 && !config.ShowClosed {
		for _, r := range result.Results {
			fmt.Printf("%s:%d  %s  %v\n", result.Target, r.Port, r.Status, r.Response.Round(time.Millisecond))
		}
	}

	return nil
}

// outputTable 以表格格式输出结果
//
// 参数:
//   - result: 扫描结果
//   - config: 扫描配置
//
// 返回值:
//   - error: 错误
func outputTable(result *ScanResult, config ScanConfig) error {
	// 按端口排序
	sort.Slice(result.Results, func(i, j int) bool {
		return result.Results[i].Port < result.Results[j].Port
	})

	// 过滤结果
	var filteredResults []PortResult
	for _, r := range result.Results {
		if r.Status == PortOpen || config.ShowClosed {
			filteredResults = append(filteredResults, r)
		}
	}

	// 如果没有结果，显示所有结果（单个端口情况）
	if len(filteredResults) == 0 && result.Stats.Total == 1 {
		filteredResults = result.Results
	}

	if len(filteredResults) == 0 {
		fmt.Println("No ports to display")
		return nil
	}

	// 使用 go-pretty/table 创建表格
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// 设置表头
	t.AppendHeader(table.Row{"PORT", "STATUS", "SERVICE", "RESPONSE TIME", "TARGET"})

	// 添加数据行
	for _, r := range filteredResults {
		service := r.Service
		if service == "" {
			service = "unknown"
		}
		t.AppendRow(table.Row{
			r.Port,
			r.Status,
			service,
			r.Response.Round(time.Millisecond),
			result.Target,
		})
	}

	// 添加统计信息（如果启用摘要）
	if config.Summary {
		t.AppendSeparator()
		t.AppendFooter(table.Row{
			fmt.Sprintf("Total: %d", result.Stats.Total),
			fmt.Sprintf("Open: %d", result.Stats.Open),
			"",
			fmt.Sprintf("Duration: %.3fs", result.Duration.Seconds()),
			"",
		})
	}

	t.Render()
	return nil
}

// outputJSON 以JSON格式输出结果
//
// 参数:
//   - result: 扫描结果
//   - config: 扫描配置
//
// 返回值:
//   - error: 错误
func outputJSON(result *ScanResult, config ScanConfig) error {
	// 过滤结果
	var filteredResults []PortResult
	for _, r := range result.Results {
		if r.Status == PortOpen || config.ShowClosed {
			filteredResults = append(filteredResults, r)
		}
	}

	output := struct {
		Target    string       `json:"target"`
		StartTime time.Time    `json:"start_time"`
		EndTime   time.Time    `json:"end_time"`
		Duration  string       `json:"duration"`
		Results   []PortResult `json:"results"`
		Stats     ScanStats    `json:"stats"`
	}{
		Target:    result.Target,
		StartTime: result.StartTime,
		EndTime:   result.EndTime,
		Duration:  result.Duration.String(),
		Results:   filteredResults,
		Stats:     result.Stats,
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(output); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	if config.Output != "" {
		if err := os.WriteFile(config.Output, buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		fmt.Printf("Results saved to: %s\n", config.Output)
	} else {
		fmt.Print(buf.String())
	}

	return nil
}

// outputCSV 以CSV格式输出结果
//
// 参数:
//   - result: 扫描结果
//   - config: 扫描配置
//
// 返回值:
//   - error: 错误
func outputCSV(result *ScanResult, config ScanConfig) error {
	// 过滤结果
	var filteredResults []PortResult
	for _, r := range result.Results {
		if r.Status == PortOpen || config.ShowClosed {
			filteredResults = append(filteredResults, r)
		}
	}

	// 按端口排序
	sort.Slice(filteredResults, func(i, j int) bool {
		return filteredResults[i].Port < filteredResults[j].Port
	})

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// 写入表头
	if err := writer.Write([]string{"Port", "Status", "Service", "Response Time (ms)"}); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// 写入数据
	for _, r := range filteredResults {
		record := []string{
			strconv.Itoa(r.Port),
			r.Status,
			r.Service,
			strconv.FormatInt(r.Response.Milliseconds(), 10),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return fmt.Errorf("CSV writer error: %w", err)
	}

	if config.Output != "" {
		if err := os.WriteFile(config.Output, buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		fmt.Printf("Results saved to: %s\n", config.Output)
	} else {
		fmt.Print(buf.String())
	}

	return nil
}

// GetDefaultConcurrent 获取默认并发数
//
// 返回值:
//   - int: CPU核心数*2
func GetDefaultConcurrent() int {
	return runtime.NumCPU() * 2
}
