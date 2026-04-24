package cli

import (
	"fmt"
	"time"

	"gitee.com/MM-Q/fck/internal/commands/ping"
	"gitee.com/MM-Q/qflag"
)

var PingCmd *qflag.Cmd

var (
	pingCount    *qflag.IntFlag      // -c, --count      发送包数量
	pingInterval *qflag.DurationFlag // -i, --interval   发送间隔
	pingTimeout  *qflag.DurationFlag // -W, --timeout    总超时时间
	pingSize     *qflag.IntFlag      // -s, --size       数据包大小
	pingTTL      *qflag.IntFlag      // -t, --ttl        TTL 值
	pingQuiet    *qflag.BoolFlag     // -q, --quiet      静默模式
)

func init() {
	PingCmd = qflag.NewCmd("ping", "", qflag.ExitOnError)

	pingCount = PingCmd.Int("count", "c", "发送包数量, 0 表示无限", 4)
	pingInterval = PingCmd.Duration("interval", "i", "发送间隔, 如: 1s, 500ms", time.Second)
	pingTimeout = PingCmd.Duration("timeout", "W", "总超时时间", time.Second*5)
	pingSize = PingCmd.Int("size", "s", "数据包大小(字节)", 56)
	pingTTL = PingCmd.Int("ttl", "t", "TTL 生存时间", 64)
	pingQuiet = PingCmd.Bool("quiet", "q", "静默模式, 只显示统计结果", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "测试网络连通性",
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s ping [options] <host>", qflag.Root.Name()),
		Notes: []string{
			"host 可以是 IP 地址或域名",
			"需要管理员/root 权限发送 ICMP 包",
			"Windows 上可能需要关闭防火墙",
			"按 Ctrl+C 可以中断 ping",
		},
		Examples: map[string]string{
			"ping 百度":    fmt.Sprintf("%s ping baidu.com", qflag.Root.Name()),
			"ping 指定次数":  fmt.Sprintf("%s ping -c 10 baidu.com", qflag.Root.Name()),
			"ping 指定间隔":  fmt.Sprintf("%s ping -i 2s baidu.com", qflag.Root.Name()),
			"ping 大包测试":  fmt.Sprintf("%s ping -s 1400 baidu.com", qflag.Root.Name()),
			"ping IP 地址": fmt.Sprintf("%s ping 8.8.8.8", qflag.Root.Name()),
			"静默模式":       fmt.Sprintf("%s ping -q baidu.com", qflag.Root.Name()),
			"快速 ping":    fmt.Sprintf("%s ping -i 100ms -c 100 baidu.com", qflag.Root.Name()),
		},
	}

	if err := PingCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts error: %w", err))
	}

	PingCmd.SetRun(runPing)
}

func runPing(cmd qflag.Command) error {
	targets := cmd.Args()
	if len(targets) == 0 {
		return fmt.Errorf("please specify a host to ping")
	}
	if len(targets) > 1 {
		return fmt.Errorf("only one host can be pinged at a time")
	}

	config := ping.PingConfig{
		Host:     targets[0],
		Count:    pingCount.Get(),
		Interval: pingInterval.Get(),
		Timeout:  pingTimeout.Get(),
		Size:     pingSize.Get(),
		TTL:      pingTTL.Get(),
		Quiet:    pingQuiet.Get(),
	}

	return ping.PingCmdMain(config)
}
