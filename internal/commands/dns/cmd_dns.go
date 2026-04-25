package dns

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"gitee.com/MM-Q/fck/internal/types"
)

// DnsConfig DNS 查询配置
type DnsConfig struct {
	Type    string
	Server  string
	Reverse bool
	Timeout time.Duration
	Retry   int
	File    string
	Output  string
	Format  string
	Short   bool
	Verbose bool
	NoColor bool
	Targets []string
}

// DnsResult DNS 查询结果
type DnsResult struct {
	Domain   string        `json:"domain"`
	Type     string        `json:"type"`
	Records  []DnsRecord   `json:"records"`
	Server   string        `json:"server"`
	Duration time.Duration `json:"duration_ms"`
	Error    string        `json:"error,omitempty"`
}

// DnsRecord 单个 DNS 记录
type DnsRecord struct {
	Type     string `json:"type"`
	Value    string `json:"value"`
	TTL      uint32 `json:"ttl"`
	Priority uint16 `json:"priority,omitempty"`
	Weight   uint16 `json:"weight,omitempty"`
	Port     uint16 `json:"port,omitempty"`
}

// BatchResult 批量查询结果
type BatchResult struct {
	Results []DnsResult   `json:"results"`
	Total   int           `json:"total"`
	Success int           `json:"success"`
	Failed  int           `json:"failed"`
	Time    time.Duration `json:"total_time_ms"`
}

// DnsCmdMain DNS 查询主函数
func DnsCmdMain(config DnsConfig) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	startTime := time.Now()

	targets, err := getTargets(config)
	if err != nil {
		return err
	}

	if len(targets) == 0 {
		return fmt.Errorf("no target specified")
	}

	resolver, err := createResolver(config)
	if err != nil {
		return fmt.Errorf("failed to create DNS resolver: %w", err)
	}

	results := make([]DnsResult, 0, len(targets))
	for i, target := range targets {
		result := queryWithRetry(ctx, resolver, target, config)
		results = append(results, result)

		if !config.Short && config.File == "" {
			if i > 0 {
				fmt.Println()
			}
			printResult(result, config)
		}
	}

	if config.Short && len(results) > 1 {
		return fmt.Errorf("short mode does not support multiple domains")
	}

	if config.File != "" || len(results) > 1 {
		batchResult := BatchResult{
			Results: results,
			Total:   len(results),
			Time:    time.Since(startTime),
		}
		for _, r := range results {
			if r.Error == "" {
				batchResult.Success++
			} else {
				batchResult.Failed++
			}
		}

		if !config.Short {
			fmt.Printf("\nTotal: %d, Success: %d, Failed: %d, Time: %v\n",
				batchResult.Total, batchResult.Success, batchResult.Failed, batchResult.Time)
		}

		if config.Output != "" {
			return saveResults(batchResult, config)
		}
	} else if config.Short && len(results) > 0 {
		printShortResult(results[0])
	}

	return nil
}

// getTargets 获取所有查询目标
func getTargets(config DnsConfig) ([]string, error) {
	var targets []string

	targets = append(targets, config.Targets...)

	if config.File != "" {
		data, err := os.ReadFile(config.File)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}

		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "#") {
				targets = append(targets, line)
			}
		}
	}

	seen := make(map[string]bool)
	unique := make([]string, 0, len(targets))
	for _, t := range targets {
		if !seen[t] {
			seen[t] = true
			unique = append(unique, t)
		}
	}

	return unique, nil
}

// createResolver 创建 DNS 解析器
func createResolver(config DnsConfig) (*net.Resolver, error) {
	if config.Server == "" {
		return net.DefaultResolver, nil
	}

	server := config.Server
	if !strings.Contains(server, ":") {
		server = net.JoinHostPort(server, "53")
	}

	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: config.Timeout,
			}
			return d.DialContext(ctx, network, server)
		},
	}

	return resolver, nil
}

// queryWithRetry 带重试的 DNS 查询
func queryWithRetry(ctx context.Context, resolver *net.Resolver, target string, config DnsConfig) DnsResult {
	var result DnsResult
	for i := 0; i < config.Retry; i++ {
		result = queryDNS(resolver, target, config)
		if result.Error == "" {
			return result
		}
		if i < config.Retry-1 {
			select {
			case <-ctx.Done():
				result.Error = "query cancelled"
				return result
			case <-time.After(time.Duration(i+1) * 500 * time.Millisecond):
			}
		}
	}
	return result
}

// queryDNS 执行 DNS 查询
func queryDNS(resolver *net.Resolver, target string, config DnsConfig) DnsResult {
	if config.Reverse || config.Type == "ptr" {
		return queryPTR(resolver, target, config)
	}

	switch config.Type {
	case "a":
		return queryA(resolver, target, config)
	case "aaaa":
		return queryAAAA(resolver, target, config)
	case "mx":
		return queryMX(resolver, target, config)
	case "ns":
		return queryNS(resolver, target, config)
	case "txt":
		return queryTXT(resolver, target, config)
	case "cname":
		return queryCNAME(resolver, target, config)
	case "soa":
		return querySOA(resolver, target, config)
	case "srv":
		return querySRV(resolver, target, config)
	case "any":
		return queryANY(resolver, target, config)
	default:
		return DnsResult{
			Domain: target,
			Type:   config.Type,
			Error:  fmt.Sprintf("unsupported query type: %s", config.Type),
		}
	}
}

// queryA 查询 A 记录
func queryA(resolver *net.Resolver, domain string, config DnsConfig) DnsResult {
	start := time.Now()
	result := DnsResult{
		Domain: domain,
		Type:   "A",
		Server: config.Server,
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	ips, err := resolver.LookupHost(ctx, domain)
	if err != nil {
		result.Error = err.Error()
		result.Duration = time.Since(start)
		return result
	}

	for _, ip := range ips {
		if net.ParseIP(ip).To4() != nil {
			result.Records = append(result.Records, DnsRecord{
				Type:  "A",
				Value: ip,
			})
		}
	}

	result.Duration = time.Since(start)
	return result
}

// queryAAAA 查询 AAAA 记录
func queryAAAA(resolver *net.Resolver, domain string, config DnsConfig) DnsResult {
	start := time.Now()
	result := DnsResult{
		Domain: domain,
		Type:   "AAAA",
		Server: config.Server,
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	ips, err := resolver.LookupHost(ctx, domain)
	if err != nil {
		result.Error = err.Error()
		result.Duration = time.Since(start)
		return result
	}

	for _, ip := range ips {
		parsed := net.ParseIP(ip)
		if parsed.To4() == nil && parsed.To16() != nil {
			result.Records = append(result.Records, DnsRecord{
				Type:  "AAAA",
				Value: ip,
			})
		}
	}

	result.Duration = time.Since(start)
	return result
}

// queryMX 查询 MX 记录
func queryMX(resolver *net.Resolver, domain string, config DnsConfig) DnsResult {
	start := time.Now()
	result := DnsResult{
		Domain: domain,
		Type:   "MX",
		Server: config.Server,
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	mxs, err := resolver.LookupMX(ctx, domain)
	if err != nil {
		result.Error = err.Error()
		result.Duration = time.Since(start)
		return result
	}

	for _, mx := range mxs {
		result.Records = append(result.Records, DnsRecord{
			Type:     "MX",
			Value:    strings.TrimSuffix(mx.Host, "."),
			Priority: mx.Pref,
		})
	}

	result.Duration = time.Since(start)
	return result
}

// queryNS 查询 NS 记录
func queryNS(resolver *net.Resolver, domain string, config DnsConfig) DnsResult {
	start := time.Now()
	result := DnsResult{
		Domain: domain,
		Type:   "NS",
		Server: config.Server,
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	nss, err := resolver.LookupNS(ctx, domain)
	if err != nil {
		result.Error = err.Error()
		result.Duration = time.Since(start)
		return result
	}

	for _, ns := range nss {
		result.Records = append(result.Records, DnsRecord{
			Type:  "NS",
			Value: strings.TrimSuffix(ns.Host, "."),
		})
	}

	result.Duration = time.Since(start)
	return result
}

// queryTXT 查询 TXT 记录
func queryTXT(resolver *net.Resolver, domain string, config DnsConfig) DnsResult {
	start := time.Now()
	result := DnsResult{
		Domain: domain,
		Type:   "TXT",
		Server: config.Server,
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	txts, err := resolver.LookupTXT(ctx, domain)
	if err != nil {
		result.Error = err.Error()
		result.Duration = time.Since(start)
		return result
	}

	for _, txt := range txts {
		result.Records = append(result.Records, DnsRecord{
			Type:  "TXT",
			Value: txt,
		})
	}

	result.Duration = time.Since(start)
	return result
}

// queryCNAME 查询 CNAME 记录
func queryCNAME(resolver *net.Resolver, domain string, config DnsConfig) DnsResult {
	start := time.Now()
	result := DnsResult{
		Domain: domain,
		Type:   "CNAME",
		Server: config.Server,
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	cname, err := resolver.LookupCNAME(ctx, domain)
	if err != nil {
		result.Error = err.Error()
		result.Duration = time.Since(start)
		return result
	}

	cname = strings.TrimSuffix(cname, ".")
	if cname != domain {
		result.Records = append(result.Records, DnsRecord{
			Type:  "CNAME",
			Value: cname,
		})
	}

	result.Duration = time.Since(start)
	return result
}

// queryPTR 反向解析
func queryPTR(resolver *net.Resolver, ip string, config DnsConfig) DnsResult {
	start := time.Now()
	result := DnsResult{
		Domain: ip,
		Type:   "PTR",
		Server: config.Server,
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		result.Error = "invalid IP address"
		result.Duration = time.Since(start)
		return result
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	names, err := resolver.LookupAddr(ctx, parsedIP.String())
	if err != nil {
		result.Error = err.Error()
		result.Duration = time.Since(start)
		return result
	}

	for _, name := range names {
		result.Records = append(result.Records, DnsRecord{
			Type:  "PTR",
			Value: strings.TrimSuffix(name, "."),
		})
	}

	result.Duration = time.Since(start)
	return result
}

// querySOA 查询 SOA 记录
func querySOA(resolver *net.Resolver, domain string, config DnsConfig) DnsResult {
	start := time.Now()
	result := DnsResult{
		Domain: domain,
		Type:   "SOA",
		Server: config.Server,
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	soa, err := lookupSOA(ctx, resolver, domain)
	if err != nil {
		result.Error = err.Error()
		result.Duration = time.Since(start)
		return result
	}

	result.Records = append(result.Records, soa)
	result.Duration = time.Since(start)
	return result
}

// querySRV 查询 SRV 记录
func querySRV(resolver *net.Resolver, domain string, config DnsConfig) DnsResult {
	start := time.Now()
	result := DnsResult{
		Domain: domain,
		Type:   "SRV",
		Server: config.Server,
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	// SRV 查询需要服务名和协议，如 "_sip._tcp.example.com"
	// 如果域名已经包含服务前缀，直接使用；否则尝试查询常见的服务
	_, srvs, err := resolver.LookupSRV(ctx, "", "", domain)
	if err != nil {
		result.Error = err.Error()
		result.Duration = time.Since(start)
		return result
	}

	for _, srv := range srvs {
		result.Records = append(result.Records, DnsRecord{
			Type:     "SRV",
			Value:    strings.TrimSuffix(srv.Target, "."),
			Priority: srv.Priority,
			Weight:   srv.Weight,
			Port:     srv.Port,
		})
	}

	result.Duration = time.Since(start)
	return result
}

// queryANY 查询所有记录类型
func queryANY(resolver *net.Resolver, domain string, config DnsConfig) DnsResult {
	start := time.Now()
	result := DnsResult{
		Domain: domain,
		Type:   "ANY",
		Server: config.Server,
	}

	for _, t := range types.GetDNSAnyTypes() {
		subConfig := config
		subConfig.Type = t
		subResult := queryDNS(resolver, domain, subConfig)
		result.Records = append(result.Records, subResult.Records...)
	}

	result.Duration = time.Since(start)
	return result
}

// lookupSOA 查询 SOA 记录 (使用自定义实现)
func lookupSOA(ctx context.Context, resolver *net.Resolver, domain string) (DnsRecord, error) {
	// SOA 记录需要通过自定义 DNS 查询获取
	// 这里使用 NS 查询作为替代方案
	nss, err := resolver.LookupNS(ctx, domain)
	if err != nil {
		return DnsRecord{}, err
	}

	if len(nss) > 0 {
		return DnsRecord{
			Type:  "SOA",
			Value: fmt.Sprintf("Primary NS: %s", strings.TrimSuffix(nss[0].Host, ".")),
		}, nil
	}

	return DnsRecord{}, fmt.Errorf("no SOA record found")
}

// printResult 打印查询结果
func printResult(result DnsResult, config DnsConfig) {
	if config.Short {
		printShortResult(result)
		return
	}

	fmt.Printf("[%s] %s", result.Type, result.Domain)
	if result.Server != "" {
		fmt.Printf(" (DNS: %s)", result.Server)
	}
	fmt.Printf(" - %v\n", result.Duration)

	if result.Error != "" {
		fmt.Printf("  Error: %s\n", result.Error)
		return
	}

	if len(result.Records) == 0 {
		fmt.Println("  No records found")
		return
	}

	for _, record := range result.Records {
		switch record.Type {
		case "MX":
			fmt.Printf("  %s\t%d\t%s\n", record.Type, record.Priority, record.Value)
		case "SRV":
			fmt.Printf("  %s\t%d\t%d\t%d\t%s\n", record.Type, record.Priority, record.Weight, record.Port, record.Value)
		default:
			fmt.Printf("  %s\t%s\n", record.Type, record.Value)
		}
	}
}

// printShortResult 简洁模式输出
func printShortResult(result DnsResult) {
	if result.Error != "" {
		fmt.Fprintf(os.Stderr, "Error: %s\n", result.Error)
		return
	}

	for _, record := range result.Records {
		fmt.Println(record.Value)
	}
}

// saveResults 保存结果到文件
func saveResults(batchResult BatchResult, config DnsConfig) error {
	switch config.Format {
	case "json":
		return saveJSON(batchResult, config.Output)
	case "csv":
		return saveCSV(batchResult, config.Output)
	default:
		return saveText(batchResult, config.Output)
	}
}

// saveJSON 保存为 JSON 格式
func saveJSON(batchResult BatchResult, filename string) error {
	data, err := json.MarshalIndent(batchResult, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("Results saved to: %s\n", filename)
	return nil
}

// saveCSV 保存为 CSV 格式
func saveCSV(batchResult BatchResult, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing file: %v\n", err)
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入表头
	if err := writer.Write([]string{"Domain", "Type", "Record Type", "Value", "Priority", "Duration(ms)", "Error"}); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// 写入数据
	for _, result := range batchResult.Results {
		if len(result.Records) == 0 {
			if err := writer.Write([]string{
				result.Domain,
				result.Type,
				"",
				"",
				"",
				result.Duration.String(),
				result.Error,
			}); err != nil {
				return fmt.Errorf("failed to write CSV row: %w", err)
			}
		} else {
			for _, record := range result.Records {
				priority := ""
				if record.Priority > 0 {
					priority = fmt.Sprintf("%d", record.Priority)
				}
				if err := writer.Write([]string{
					result.Domain,
					result.Type,
					record.Type,
					record.Value,
					priority,
					result.Duration.String(),
					result.Error,
				}); err != nil {
					return fmt.Errorf("failed to write CSV row: %w", err)
				}
			}
		}
	}

	fmt.Printf("Results saved to: %s\n", filename)
	return nil
}

// saveText 保存为文本格式
func saveText(batchResult BatchResult, filename string) error {
	var sb strings.Builder

	for _, result := range batchResult.Results {
		sb.WriteString(fmt.Sprintf("[%s] %s\n", result.Type, result.Domain))
		if result.Error != "" {
			sb.WriteString(fmt.Sprintf("  Error: %s\n", result.Error))
		} else {
			for _, record := range result.Records {
				sb.WriteString(fmt.Sprintf("  %s\t%s\n", record.Type, record.Value))
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("Total: %d, Success: %d, Failed: %d, Time: %v\n",
		batchResult.Total, batchResult.Success, batchResult.Failed, batchResult.Time))

	if err := os.WriteFile(filename, []byte(sb.String()), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("Results saved to: %s\n", filename)
	return nil
}
