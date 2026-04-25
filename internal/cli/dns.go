package cli

import (
	"fmt"
	"time"

	"gitee.com/MM-Q/fck/internal/commands/dns"
	"gitee.com/MM-Q/fck/internal/types"
	"gitee.com/MM-Q/qflag"
)

// DnsCmd dns 命令
var DnsCmd *qflag.Cmd

var (
	dnsType    *qflag.EnumFlag
	dnsServer  *qflag.StringFlag
	dnsReverse *qflag.BoolFlag
	dnsAll     *qflag.BoolFlag
	dnsTimeout *qflag.DurationFlag
	dnsRetry   *qflag.IntFlag
	dnsFile    *qflag.StringFlag
	dnsOutput  *qflag.StringFlag
	dnsFormat  *qflag.EnumFlag
	dnsShort   *qflag.BoolFlag
	dnsVerbose *qflag.BoolFlag
	dnsNoColor *qflag.BoolFlag
)

var dnsFormatOptions = []string{
	"text", "json", "csv",
}

func init() {
	DnsCmd = qflag.NewCmd("dns", "dn", qflag.ExitOnError)

	dnsType = DnsCmd.Enum("type", "t", fmt.Sprintf("查询类型, 支持 %v", types.DNSQueryTypes), "a", types.DNSQueryTypes)
	dnsServer = DnsCmd.String("server", "s", "指定 DNS 服务器", "")
	dnsReverse = DnsCmd.Bool("reverse", "r", "反向解析模式", false)
	dnsAll = DnsCmd.Bool("all", "a", "查询所有记录类型 (等同于 -t any)", false)
	dnsTimeout = DnsCmd.Duration("timeout", "T", "查询超时时间", time.Second*5)
	dnsRetry = DnsCmd.Int("retry", "R", "查询失败重试次数", 3)
	dnsFile = DnsCmd.String("file", "f", "从文件读取域名列表", "")
	dnsOutput = DnsCmd.String("output", "o", "输出文件", "")
	dnsFormat = DnsCmd.Enum("format", "", fmt.Sprintf("输出格式, 支持 %v", dnsFormatOptions), "text", dnsFormatOptions)
	dnsShort = DnsCmd.Bool("short", "S", "简洁输出模式", false)
	dnsVerbose = DnsCmd.Bool("verbose", "v", "详细输出", false)
	dnsNoColor = DnsCmd.Bool("no-color", "", "禁用颜色输出", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "DNS 查询工具",
		UsageSyntax: fmt.Sprintf("%s dns [选项] <域名或IP>", qflag.Root.Name()),
		UseChinese:  true,
		Notes: []string{
			"支持查询 A、AAAA、MX、NS、TXT、CNAME、PTR、SOA 等记录类型",
			"使用 -r 参数进行反向解析 (IP 查域名)",
			"使用 -s 指定 DNS 服务器，格式为 IP:端口 (默认端口 53)",
			"使用 -f 从文件批量查询，每行一个域名",
			"使用 --short 获取简洁输出，适合脚本处理",
			"使用 -a 或 --all 查询所有记录类型 (与 -t 互斥)",
		},
		Examples: map[string]string{
			"查询 A 记录":    fmt.Sprintf("%s dns www.example.com", qflag.Root.Name()),
			"查询 MX 记录":   fmt.Sprintf("%s dns -t mx gmail.com", qflag.Root.Name()),
			"查询 NS 记录":   fmt.Sprintf("%s dns -t ns baidu.com", qflag.Root.Name()),
			"查询 TXT 记录":  fmt.Sprintf("%s dns -t txt example.com", qflag.Root.Name()),
			"查询所有记录":     fmt.Sprintf("%s dns -a www.example.com", qflag.Root.Name()),
			"反向解析":       fmt.Sprintf("%s dns -r 8.8.8.8", qflag.Root.Name()),
			"指定 DNS 服务器": fmt.Sprintf("%s dns -s 114.114.114.114 www.baidu.com", qflag.Root.Name()),
			"批量查询":       fmt.Sprintf("%s dns -f domains.txt", qflag.Root.Name()),
			"简洁输出":       fmt.Sprintf("%s dns -S www.example.com", qflag.Root.Name()),
			"导出 JSON":    fmt.Sprintf("%s dns --format json -o result.json www.example.com", qflag.Root.Name()),
			"导出 CSV":     fmt.Sprintf("%s dns --format csv -o result.csv -f domains.txt", qflag.Root.Name()),
		},
		MutexGroups: []qflag.MutexGroup{
			{
				Name:      "query-type",
				Flags:     []string{"type", "all"},
				AllowNone: true,
			},
		},
	}

	if err := DnsCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	DnsCmd.SetRun(runDns)
}

// runDns 运行 dns 命令
func runDns(cmd qflag.Command) error {
	args := cmd.Args()

	if len(args) == 0 && dnsFile.Get() == "" {
		cmd.PrintHelp()
		return nil
	}

	queryType := dnsType.Get()
	if dnsAll.Get() {
		queryType = "any"
	}

	config := dns.DnsConfig{
		Type:    queryType,
		Server:  dnsServer.Get(),
		Reverse: dnsReverse.Get(),
		Timeout: dnsTimeout.Get(),
		Retry:   dnsRetry.Get(),
		File:    dnsFile.Get(),
		Output:  dnsOutput.Get(),
		Format:  dnsFormat.Get(),
		Short:   dnsShort.Get(),
		Verbose: dnsVerbose.Get(),
		NoColor: dnsNoColor.Get(),
		Targets: args,
	}

	return dns.DnsCmdMain(config)
}
