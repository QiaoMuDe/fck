package tcp

import (
	"fmt"
	"time"

	"gitee.com/MM-Q/fck/internal/commands/tcp"
	"gitee.com/MM-Q/fck/internal/types"
	"gitee.com/MM-Q/qflag"
)

// ScanCmd tcp scan 命令
var ScanCmd *qflag.Cmd

var (
	scanTimeout    *qflag.DurationFlag
	scanConcurrent *qflag.IntFlag
	scanShowClosed *qflag.BoolFlag
	scanOutput     *qflag.StringFlag
	scanFormat     *qflag.EnumFlag
	scanProgress   *qflag.BoolFlag
	scanSummary    *qflag.BoolFlag
)

func init() {
	ScanCmd = qflag.NewCmd("scan", "sc", qflag.ExitOnError)

	scanTimeout = ScanCmd.Duration("timeout", "t", "单个端口扫描超时时间", 2*time.Second)
	scanConcurrent = ScanCmd.Int("concurrent", "c", "并发扫描数量", tcp.GetDefaultConcurrent())
	scanShowClosed = ScanCmd.Bool("show-closed", "s", "显示关闭的端口", false)
	scanOutput = ScanCmd.String("output", "o", "输出结果到文件", "")
	scanFormat = ScanCmd.Enum("format", "f", "输出格式", types.TCPScanFormatDefault, types.TCPScanFormatOptions)
	scanProgress = ScanCmd.Bool("progress", "P", "显示扫描进度", false)
	scanSummary = ScanCmd.Bool("summary", "S", "显示扫描统计摘要", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "TCP 端口扫描工具",
		UsageSyntax: fmt.Sprintf("%s tcp scan [options] <target> <ports>", qflag.Root.Name()),
		UseChinese:  true,
		Notes: []string{
			"必须指定目标主机和端口",
			"支持域名和 IP 地址",
			"高并发扫描可能会被防火墙拦截",
			"端口参数支持: 80, 1-1024, 80,443,8080",
		},
		Examples: map[string]string{
			"扫描单个端口":     fmt.Sprintf("%s tcp scan 192.168.1.1 80", qflag.Root.Name()),
			"扫描端口范围":     fmt.Sprintf("%s tcp scan 192.168.1.1 1-1024", qflag.Root.Name()),
			"扫描多个端口":     fmt.Sprintf("%s tcp scan example.com 22,80,443", qflag.Root.Name()),
			"高并发快速扫描":    fmt.Sprintf("%s tcp scan -c 500 -t 1s 192.168.1.1 1-65535", qflag.Root.Name()),
			"导出 JSON 结果": fmt.Sprintf("%s tcp scan -f json -o result.json 192.168.1.1 1-1024", qflag.Root.Name()),
		},
	}

	if err := ScanCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	ScanCmd.SetRun(runScan)
}

// runScan 运行 scan 命令
//
// 参数:
//   - cmd: 命令接口
//
// 返回值:
//   - error: 错误
func runScan(cmd qflag.Command) error {
	args := cmd.Args()
	if len(args) < 2 {
		return fmt.Errorf("usage: fck tcp scan <target> <ports>")
	}

	config := tcp.ScanConfig{
		Target:     args[0],
		Ports:      args[1],
		Timeout:    scanTimeout.Get(),
		Concurrent: scanConcurrent.Get(),
		ShowClosed: scanShowClosed.Get(),
		Output:     scanOutput.Get(),
		Format:     scanFormat.Get(),
		Progress:   scanProgress.Get(),
		Summary:    scanSummary.Get(),
	}

	return tcp.ScanCmdMain(config)
}
