package cli

import (
	"fmt"
	"strconv"
	"time"

	"gitee.com/MM-Q/fck/internal/commands/tcp"
	"gitee.com/MM-Q/fck/internal/types"
	"gitee.com/MM-Q/qflag"
)

var TcpCmd *qflag.Cmd

var (
	tcpTimeout     *qflag.DurationFlag // -t, --timeout  连接超时时间
	tcpCount       *qflag.IntFlag      // -c, --count    连接次数
	tcpInterval    *qflag.DurationFlag // -i, --interval 连接间隔
	tcpScan        *qflag.BoolFlag     // -s, --scan     扫描模式
	tcpRange       *qflag.StringFlag   // -r, --range    端口范围
	tcpOpenOnly    *qflag.BoolFlag     // -o, --open     仅显示开放端口
	tcpBanner      *qflag.BoolFlag     // -b, --banner   获取 banner
	tcpData        *qflag.StringFlag   // -d, --data     发送数据
	tcpPath        *qflag.StringFlag   // -p, --path     文件/目录/通配符路径
	tcpMaxFileSize *qflag.SizeFlag     // --max-file-size 最大文件大小限制
	tcpWait        *qflag.DurationFlag // -w, --wait     等待响应时间
	tcpQuiet       *qflag.BoolFlag     // -q, --quiet       静默模式
	tcpJson        *qflag.BoolFlag     // -j, --json        JSON 输出
	tcpListen      *qflag.BoolFlag     // -l, --listen      监听模式
	tcpInteractive *qflag.BoolFlag     // -I, --interactive 交互式模式
	tcpHex         *qflag.BoolFlag     // --hex             十六进制显示
)

func init() {
	TcpCmd = qflag.NewCmd("tcp", "", qflag.ExitOnError)

	tcpTimeout = TcpCmd.Duration("timeout", "t", "连接超时时间", time.Second*5)
	tcpCount = TcpCmd.Int("count", "c", "连接次数, 0 表示无限", 1)
	tcpInterval = TcpCmd.Duration("interval", "i", "连接间隔, 如: 1s, 500ms", time.Second)
	tcpScan = TcpCmd.Bool("scan", "s", "端口扫描模式", false)
	tcpRange = TcpCmd.String("range", "r", "端口范围, 如: 80-100,443,8080", "")
	tcpOpenOnly = TcpCmd.Bool("open", "o", "仅显示开放的端口", false)
	tcpBanner = TcpCmd.Bool("banner", "b", "尝试获取服务 banner", false)
	tcpData = TcpCmd.String("data", "d", "发送的数据内容", "")
	tcpPath = TcpCmd.String("path", "p", "文件/目录/通配符路径,支持发送单个文件、目录下所有文件或通配符匹配的文件", "")
	tcpMaxFileSize = TcpCmd.Size("max-file-size", "mfs", "最大文件大小限制,支持单位如: 10MB, 1GB", types.DefaultMaxBufferSize)
	tcpWait = TcpCmd.Duration("wait", "w", "等待响应时间", 0)
	tcpQuiet = TcpCmd.Bool("quiet", "q", "静默模式", false)
	tcpJson = TcpCmd.Bool("json", "j", "JSON 格式输出", false)
	tcpListen = TcpCmd.Bool("listen", "l", "监听模式,在指定端口接收数据", false)
	tcpInteractive = TcpCmd.Bool("interactive", "I", "交互式模式,支持实时发送和接收数据", false)
	tcpHex = TcpCmd.Bool("hex", "", "以十六进制格式显示接收的数据", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "TCP 连接测试和端口扫描工具",
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s tcp [options] <host> [port]", qflag.Root.Name()),
		Notes: []string{
			"host 可以是 IP 地址或域名",
			"port 为端口号, 范围 1-65535",
			"扫描模式(-s)需要指定 -r 端口范围",
			"非扫描模式需要指定 port 参数",
			"使用 -j 获取 JSON 格式输出便于脚本处理",
		},
		Examples: map[string]string{
			"测试单个端口":      fmt.Sprintf("%s tcp baidu.com 80", qflag.Root.Name()),
			"指定超时时间":      fmt.Sprintf("%s tcp -t 10s baidu.com 443", qflag.Root.Name()),
			"多次连接测试":      fmt.Sprintf("%s tcp -c 5 -i 2s baidu.com 80", qflag.Root.Name()),
			"扫描端口范围":      fmt.Sprintf("%s tcp -s -r 1-1000 baidu.com", qflag.Root.Name()),
			"扫描指定端口":      fmt.Sprintf("%s tcp -s -r 80,443,8080 baidu.com", qflag.Root.Name()),
			"仅显示开放端口":     fmt.Sprintf("%s tcp -s -r 1-100 -o baidu.com", qflag.Root.Name()),
			"获取服务 banner": fmt.Sprintf("%s tcp -b -w 2s baidu.com 80", qflag.Root.Name()),
			"发送数据":        fmt.Sprintf("%s tcp -d \"hello\" -w 5s target.com 8080", qflag.Root.Name()),
			"发送文件":        fmt.Sprintf("%s tcp -p ./data.txt target.com 8080", qflag.Root.Name()),
			"发送目录":        fmt.Sprintf("%s tcp -p ./configs/ target.com 8080", qflag.Root.Name()),
			"发送通配符匹配":     fmt.Sprintf("%s tcp -p \"./logs/*.log\" target.com 8080", qflag.Root.Name()),
			"JSON 格式扫描":   fmt.Sprintf("%s tcp -s -r 22,80,443 -j target.com", qflag.Root.Name()),
			"监听端口接收数据":    fmt.Sprintf("%s tcp -l 8080", qflag.Root.Name()),
			"交互式模式":       fmt.Sprintf("%s tcp -I target.com 8080", qflag.Root.Name()),
			"交互式十六进制显示":   fmt.Sprintf("%s tcp -I --hex target.com 8080", qflag.Root.Name()),
		},
		MutexGroups: []qflag.MutexGroup{
			{
				Name:      "mode",
				Flags:     []string{"scan", "banner", "data", "path", "listen", "interactive"},
				AllowNone: true,
			},
		},
	}

	if err := TcpCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts error: %w", err))
	}

	TcpCmd.SetRun(runTcp)
}

func runTcp(cmd qflag.Command) error {
	args := cmd.Args()
	if len(args) == 0 {
		return fmt.Errorf("please specify a host")
	}

	host := args[0]

	// 解析端口
	// 检查参数
	scanMode := tcpScan.Get() || tcpRange.Get() != ""
	listenMode := tcpListen.Get()

	port := 0
	if listenMode {
		// 监听模式: port 是第一个参数
		if len(args) >= 1 {
			p, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid port: %s", args[0])
			}
			port = p
		}
		if port == 0 {
			return fmt.Errorf("listen mode requires a port")
		}
	} else {
		// 其他模式: port 是第二个参数
		if len(args) >= 2 {
			p, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid port: %s", args[1])
			}
			port = p
		}
		if !scanMode && port == 0 {
			return fmt.Errorf("please specify a port (or use -s for scan mode)")
		}
	}

	// 检查扫描模式是否指定了端口范围
	if scanMode && tcpRange.Get() == "" {
		return fmt.Errorf("scan mode requires -r/--range to specify port range")
	}

	// 获取最大文件大小，如果小于等于0则使用默认值10MB
	maxFileSize := tcpMaxFileSize.Get()
	if maxFileSize <= 0 {
		maxFileSize = types.DefaultMaxBufferSize
	}

	config := tcp.TcpConfig{
		Host:        host,
		Port:        port,
		Timeout:     tcpTimeout.Get(),
		Count:       tcpCount.Get(),
		Interval:    tcpInterval.Get(),
		Scan:        tcpScan.Get(),
		Range:       tcpRange.Get(),
		OpenOnly:    tcpOpenOnly.Get(),
		Banner:      tcpBanner.Get(),
		Data:        tcpData.Get(),
		Path:        tcpPath.Get(),
		MaxFileSize: maxFileSize,
		Wait:        tcpWait.Get(),
		Quiet:       tcpQuiet.Get(),
		Json:        tcpJson.Get(),
		Listen:      tcpListen.Get(),
		Interactive: tcpInteractive.Get(),
		Hex:         tcpHex.Get(),
	}

	return tcp.TcpCmdMain(config)
}
