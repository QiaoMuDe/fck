# DNS 命令设计方案

## 命令概述

DNS 查询工具，用于查询域名的各种 DNS 记录类型，支持正向解析、反向解析和指定 DNS 服务器查询。

## 命令结构

```
fck dns [选项] <域名或IP>
```

## 功能特性

### 核心功能
1. **A 记录查询** - 域名解析到 IPv4 地址
2. **AAAA 记录查询** - 域名解析到 IPv6 地址
3. **MX 记录查询** - 邮件服务器记录
4. **NS 记录查询** - 域名服务器记录
5. **TXT 记录查询** - 文本记录
6. **CNAME 记录查询** - 别名记录
7. **PTR 反向解析** - IP 反查域名
8. **SOA 记录查询** - 授权信息记录
9. **SRV 记录查询** - 服务定位记录
10. **CAA 记录查询** - 证书颁发机构授权
11. **ANY 记录查询** - 查询所有记录类型
12. **批量查询** - 从文件读取域名列表批量查询

### 高级功能
- **指定 DNS 服务器** - 使用自定义 DNS 服务器进行查询
- **超时控制** - 设置查询超时时间
- **重试机制** - 查询失败自动重试
- **多种输出格式** - text/json/csv
- **结果缓存** - 可选启用查询结果缓存

## CLI 定义

### 文件位置
`internal/cli/dns.go`

### 标志定义

| 长名称 | 短名称 | 类型 | 默认值 | 说明 |
|--------|--------|------|--------|------|
| type | t | enum | a | 查询类型 (a/aaaa/mx/ns/txt/cname/ptr/soa/srv/caa/any) |
| server | s | string | "" | 指定 DNS 服务器 (如 8.8.8.8:53) |
| reverse | r | bool | false | 反向解析模式 (IP 查域名) |
| timeout | T | duration | 5s | 查询超时时间 |
| retry | R | int | 3 | 查询失败重试次数 |
| file | f | string | "" | 从文件读取域名列表 |
| output | o | string | "" | 输出文件 |
| format | | enum | text | 输出格式 (text/json/csv) |
| short | S | bool | false | 简洁输出模式 |
| verbose | v | bool | false | 详细输出 |
| no-color | | bool | false | 禁用颜色输出 |

### 代码实现

```go
package cli

import (
	"fmt"
	"time"

	"gitee.com/MM-Q/fck/internal/commands/dns"
	"gitee.com/MM-Q/qflag"
)

var DnsCmd *qflag.Cmd

var (
	dnsType     *qflag.EnumFlag
	dnsServer   *qflag.StringFlag
	dnsReverse  *qflag.BoolFlag
	dnsTimeout  *qflag.DurationFlag
	dnsRetry    *qflag.IntFlag
	dnsFile     *qflag.StringFlag
	dnsOutput   *qflag.StringFlag
	dnsFormat   *qflag.EnumFlag
	dnsShort    *qflag.BoolFlag
	dnsVerbose  *qflag.BoolFlag
	dnsNoColor  *qflag.BoolFlag
)

func init() {
	DnsCmd = qflag.NewCmd("dns", "dn", qflag.ExitOnError)

	dnsType = DnsCmd.Enum("type", "t", "查询类型", "a", []string{
		"a", "aaaa", "mx", "ns", "txt", "cname", "ptr", "soa", "srv", "caa", "any",
	})
	dnsServer = DnsCmd.String("server", "s", "指定 DNS 服务器", "")
	dnsReverse = DnsCmd.Bool("reverse", "r", "反向解析模式", false)
	dnsTimeout = DnsCmd.Duration("timeout", "T", "查询超时时间", time.Second*5)
	dnsRetry = DnsCmd.Int("retry", "R", "查询失败重试次数", 3)
	dnsFile = DnsCmd.String("file", "f", "从文件读取域名列表", "")
	dnsOutput = DnsCmd.String("output", "o", "输出文件", "")
	dnsFormat = DnsCmd.Enum("format", "", "输出格式", "text", []string{"text", "json", "csv"})
	dnsShort = DnsCmd.Bool("short", "S", "简洁输出模式", false)
	dnsVerbose = DnsCmd.Bool("verbose", "v", "详细输出", false)
	dnsNoColor = DnsCmd.Bool("no-color", "", "禁用颜色输出", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "DNS 查询工具",
		UsageSyntax: fmt.Sprintf("%s dns [选项] <域名或IP>", qflag.Root.Name()),
		UseChinese:  true,
		Notes: []string{
			"支持查询 A、AAAA、MX、NS、TXT、CNAME、PTR、SOA、SRV、CAA 等记录类型",
			"使用 -r 参数进行反向解析 (IP 查域名)",
			"使用 -s 指定 DNS 服务器，格式为 IP:端口 (默认端口 53)",
			"使用 -f 从文件批量查询，每行一个域名",
			"使用 --short 获取简洁输出，适合脚本处理",
		},
		Examples: map[string]string{
			"查询 A 记录":           fmt.Sprintf("%s dns www.example.com", qflag.Root.Name()),
			"查询 MX 记录":          fmt.Sprintf("%s dns -t mx gmail.com", qflag.Root.Name()),
			"查询 NS 记录":          fmt.Sprintf("%s dns -t ns baidu.com", qflag.Root.Name()),
			"查询 TXT 记录":         fmt.Sprintf("%s dns -t txt example.com", qflag.Root.Name()),
			"反向解析":             fmt.Sprintf("%s dns -r 8.8.8.8", qflag.Root.Name()),
			"指定 DNS 服务器":       fmt.Sprintf("%s dns -s 114.114.114.114 www.baidu.com", qflag.Root.Name()),
			"查询所有记录":         fmt.Sprintf("%s dns -t any example.com", qflag.Root.Name()),
			"批量查询":             fmt.Sprintf("%s dns -f domains.txt", qflag.Root.Name()),
			"简洁输出":             fmt.Sprintf("%s dns -S www.example.com", qflag.Root.Name()),
			"导出 JSON":            fmt.Sprintf("%s dns --format json -o result.json www.example.com", qflag.Root.Name()),
			"导出 CSV":             fmt.Sprintf("%s dns --format csv -o result.csv -f domains.txt", qflag.Root.Name()),
		},
	}

	if err := DnsCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	DnsCmd.SetRun(runDns)
}

func runDns(cmd qflag.Command) error {
	args := cmd.Args()
	
	// 如果没有参数且没有指定文件，显示帮助
	if len(args) == 0 && dnsFile.Get() == "" {
		cmd.PrintHelp()
		return nil
	}

	config := dns.DnsConfig{
		Type:     dnsType.Get(),
		Server:   dnsServer.Get(),
		Reverse:  dnsReverse.Get(),
		Timeout:  dnsTimeout.Get(),
		Retry:    dnsRetry.Get(),
		File:     dnsFile.Get(),
		Output:   dnsOutput.Get(),
		Format:   dnsFormat.Get(),
		Short:    dnsShort.Get(),
		Verbose:  dnsVerbose.Get(),
		NoColor:  dnsNoColor.Get(),
		Targets:  args,
	}

	return dns.DnsCmdMain(config)
}
```

## 业务逻辑

### 文件位置
`internal/commands/dns/cmd_dns.go`

### 配置结构体

```go
package dns

import "time"

// DnsConfig DNS 查询配置
type DnsConfig struct {
	Type     string
	Server   string
	Reverse  bool
	Timeout  time.Duration
	Retry    int
	File     string
	Output   string
	Format   string
	Short    bool
	Verbose  bool
	NoColor  bool
	Targets  []string
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
	Priority uint16 `json:"priority,omitempty"` // MX、SRV 记录使用
	Weight   uint16 `json:"weight,omitempty"`   // SRV 记录使用
	Port     uint16 `json:"port,omitempty"`     // SRV 记录使用
}

// BatchResult 批量查询结果
type BatchResult struct {
	Results []DnsResult   `json:"results"`
	Total   int           `json:"total"`
	Success int           `json:"success"`
	Failed  int           `json:"failed"`
	Time    time.Duration `json:"total_time_ms"`
}
```

### 主函数

```go
// DnsCmdMain DNS 查询主函数
func DnsCmdMain(config DnsConfig) error {
	startTime := time.Now()

	// 获取所有查询目标
	targets, err := getTargets(config)
	if err != nil {
		return err
	}

	if len(targets) == 0 {
		return fmt.Errorf("未指定查询目标")
	}

	// 创建 DNS 解析器
	resolver, err := createResolver(config)
	if err != nil {
		return fmt.Errorf("创建 DNS 解析器失败: %w", err)
	}

	// 执行查询
	results := make([]DnsResult, 0, len(targets))
	for _, target := range targets {
		result := queryDNS(resolver, target, config)
		results = append(results, result)

		// 如果不是简洁模式，实时输出结果
		if !config.Short && config.File == "" {
			printResult(result, config)
		}
	}

	// 输出结果
	if config.File != "" || len(results) > 1 {
		// 批量查询结果
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
			fmt.Printf("\n总计: %d, 成功: %d, 失败: %d, 耗时: %v\n",
				batchResult.Total, batchResult.Success, batchResult.Failed, batchResult.Time)
		}

		// 保存结果
		if config.Output != "" {
			return saveResults(batchResult, config)
		}
	} else if config.Short && len(results) > 0 {
		// 简洁模式输出
		printShortResult(results[0], config)
	}

	return nil
}
```

## 核心功能实现

### 1. 获取查询目标

```go
// getTargets 获取所有查询目标
func getTargets(config DnsConfig) ([]string, error) {
	var targets []string

	// 从命令行参数获取
	targets = append(targets, config.Targets...)

	// 从文件读取
	if config.File != "" {
		data, err := os.ReadFile(config.File)
		if err != nil {
			return nil, fmt.Errorf("读取文件失败: %w", err)
		}

		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "#") {
				targets = append(targets, line)
			}
		}
	}

	// 去重
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
```

### 2. 创建 DNS 解析器

```go
// createResolver 创建 DNS 解析器
func createResolver(config DnsConfig) (*net.Resolver, error) {
	if config.Server == "" {
		// 使用系统默认 DNS
		return net.DefaultResolver, nil
	}

	// 解析 DNS 服务器地址
	server := config.Server
	if !strings.Contains(server, ":") {
		server = net.JoinHostPort(server, "53")
	}

	// 创建自定义解析器
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
```

### 3. DNS 查询实现

```go
// queryDNS 执行 DNS 查询
func queryDNS(resolver *net.Resolver, target string, config DnsConfig) DnsResult {
	start := time.Now()
	result := DnsResult{
		Domain: target,
		Type:   config.Type,
		Server: config.Server,
	}

	// 如果是反向解析
	if config.Reverse || config.Type == "ptr" {
		return queryPTR(resolver, target, config)
	}

	// 根据类型执行查询
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
	case "any":
		return queryANY(resolver, target, config)
	default:
		result.Error = fmt.Sprintf("不支持的查询类型: %s", config.Type)
	}

	result.Duration = time.Since(start)
	return result
}

// queryA 查询 A 记录
func queryA(resolver *net.Resolver, domain string, config DnsConfig) DnsResult {
	start := time.Now()
	result := DnsResult{
		Domain: domain,
		Type:   "A",
	}

	ips, err := resolver.LookupHost(context.Background(), domain)
	if err != nil {
		result.Error = err.Error()
		result.Duration = time.Since(start)
		return result
	}

	for _, ip := range ips {
		// 只保留 IPv4 地址
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

// queryMX 查询 MX 记录
func queryMX(resolver *net.Resolver, domain string, config DnsConfig) DnsResult {
	start := time.Now()
	result := DnsResult{
		Domain: domain,
		Type:   "MX",
	}

	mxs, err := resolver.LookupMX(context.Background(), domain)
	if err != nil {
		result.Error = err.Error()
		result.Duration = time.Since(start)
		return result
	}

	for _, mx := range mxs {
		result.Records = append(result.Records, DnsRecord{
			Type:     "MX",
			Value:    mx.Host,
			Priority: mx.Pref,
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
	}

	// 验证 IP 地址
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		result.Error = "无效的 IP 地址"
		result.Duration = time.Since(start)
		return result
	}

	names, err := resolver.LookupAddr(context.Background(), ip)
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
```

### 4. 输出格式化

```go
// printResult 打印查询结果
func printResult(result DnsResult, config DnsConfig) {
	if config.Short {
		printShortResult(result, config)
		return
	}

	// 打印标题
	fmt.Printf("\n[%s] %s", result.Type, result.Domain)
	if result.Server != "" {
		fmt.Printf(" (DNS: %s)", result.Server)
	}
	fmt.Printf(" - %v\n", result.Duration)

	// 打印错误
	if result.Error != "" {
		fmt.Printf("  错误: %s\n", result.Error)
		return
	}

	// 打印记录
	if len(result.Records) == 0 {
		fmt.Println("  未找到记录")
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
func printShortResult(result DnsResult, config DnsConfig) {
	if result.Error != "" {
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
```

## 使用示例

```bash
# 查询 A 记录
fck dns www.example.com

# 查询 AAAA 记录
fck dns -t aaaa www.google.com

# 查询 MX 记录
fck dns -t mx gmail.com

# 查询 NS 记录
fck dns -t ns baidu.com

# 查询 TXT 记录
fck dns -t txt example.com

# 反向解析
fck dns -r 8.8.8.8

# 指定 DNS 服务器
fck dns -s 114.114.114.114 www.baidu.com
fck dns -s 8.8.8.8:53 www.google.com

# 批量查询
fck dns -f domains.txt

# 简洁输出 (适合脚本)
fck dns -S www.example.com

# 导出 JSON
fck dns --format json -o result.json www.example.com

# 导出 CSV
fck dns --format csv -o result.csv -f domains.txt

# 查询所有记录
fck dns -t any example.com
```

## 注册到根命令

在 `internal/cli/root.go` 的 `SubCmds` 中添加：

```go
SubCmds: []qflag.Command{
	// ... 其他命令
	DnsCmd,
},
```

## 实现步骤

1. 创建业务逻辑目录 `internal/commands/dns/`
2. 创建 `internal/commands/dns/cmd_dns.go` - 主逻辑
3. 创建 `internal/commands/dns/query.go` - DNS 查询实现
4. 创建 `internal/commands/dns/output.go` - 输出格式化
5. 创建 `internal/cli/dns.go` - CLI 定义
6. 在 `internal/cli/root.go` 中注册命令
7. 测试编译和功能验证

## 复杂度评估

| 方面 | 评估 | 说明 |
|------|------|------|
| 代码量 | 低 | 约 400-500 行 |
| 依赖 | 低 | 仅使用 Go 标准库 |
| 难度 | 低 | 使用 net.Resolver 标准接口 |
| 文件数 | 少 | 3-4 个文件 |
| 测试 | 简单 | 有标准库支持，易于测试 |

## 注意事项

1. **Go 标准库限制**: `net.Resolver` 不支持所有 DNS 记录类型 (如 SRV、CAA 需要额外实现)
2. **超时控制**: 需要处理 DNS 查询超时
3. **并发查询**: 批量查询时可考虑并发优化
4. **错误处理**: DNS 查询可能失败，需要友好的错误提示
