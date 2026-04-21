package cli

import (
	"fmt"
	"strings"
	"time"

	"gitee.com/MM-Q/fck/internal/commands/watch"
	"gitee.com/MM-Q/qflag"
)

var WatchCmd *qflag.Cmd

var (
	watchInterval *qflag.DurationFlag
	watchMaxCount *qflag.IntFlag
	watchExitErr  *qflag.BoolFlag
	watchNoHeader *qflag.BoolFlag
	watchTimeout  *qflag.DurationFlag
	watchQuiet    *qflag.BoolFlag
	watchDiff     *qflag.BoolFlag
	watchPrecise  *qflag.BoolFlag
	watchNoColor  *qflag.BoolFlag
)

func init() {
	WatchCmd = qflag.NewCmd("watch", "", qflag.ExitOnError)

	// 基本配置
	watchInterval = WatchCmd.Duration("interval", "n", "执行间隔时间，默认2秒", 2*time.Second)
	watchMaxCount = WatchCmd.Int("count", "c", "执行次数限制，-1表示无限制(默认)", -1)
	watchTimeout = WatchCmd.Duration("timeout", "t", "单次命令执行超时时间，默认30秒", 30*time.Second)

	// 显示控制
	watchNoHeader = WatchCmd.Bool("no-title", "", "不显示标题栏", false)
	watchQuiet = WatchCmd.Bool("quiet", "q", "静默模式，不显示标题栏和命令输出", false)
	watchNoColor = WatchCmd.Bool("no-color", "", "禁用颜色输出", false)

	// 差异高亮
	watchDiff = WatchCmd.Bool("differences", "d", "高亮显示变化的行", false)

	// 执行控制
	watchExitErr = WatchCmd.Bool("errexit", "e", "命令执行失败时退出", false)
	watchPrecise = WatchCmd.Bool("precise", "p", "精确计时模式，补偿命令执行耗时", false)

	cmdOpts := &qflag.CmdOpts{
		Desc: "命令监控工具 - 周期性执行命令并显示输出",
		Notes: []string{
			"类似于 Linux 系统的 watch 命令",
			"使用 Ctrl+C 可以随时停止监控",
			"命令执行失败时默认继续监控，除非使用 -e 标志",
			"如果监控的命令包含以 - 开头的参数，请使用 -- 分隔，如: fck watch -- ls -la",
		},
		UseChinese: true,
		Examples: map[string]string{
			"每隔 2 秒执行 ls -la": fmt.Sprintf("%s watch -n 2 -- ls -la", qflag.Root.Name()),
			"高亮显示变化的行":        fmt.Sprintf("%s watch -d date", qflag.Root.Name()),
			"精确计时模式 (补偿执行耗时)": fmt.Sprintf("%s watch -p -n 1 ./benchmark.sh", qflag.Root.Name()),
			"执行 10 次后退出":      fmt.Sprintf("%s watch -c 10 -n 5 -- curl -s http://api.example.com/status", qflag.Root.Name()),
			"出错时立即退出":         fmt.Sprintf("%s watch -e kubectl get pods", qflag.Root.Name()),
			"静默模式 (无标题栏)":     fmt.Sprintf("%s watch -q ./monitor.sh", qflag.Root.Name()),
		},
		UsageSyntax: fmt.Sprintf("%s watch [options] <command>", qflag.Root.Name()),
	}

	if err := WatchCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	WatchCmd.SetRun(runWatch)
}

func runWatch(cmd qflag.Command) error {
	args := cmd.Args()
	if len(args) < 1 {
		return fmt.Errorf("missing command argument")
	}

	// 组合命令参数
	var command string
	if len(args) == 1 {
		command = args[0]
	} else {
		command = strings.Join(args, " ")
	}

	// 确定差异模式
	diffMode := watch.DiffModeNone
	if watchDiff.Get() {
		diffMode = watch.DiffModeLine
	}

	// 静默模式覆盖其他显示选项
	noHeader := watchNoHeader.Get()
	quiet := watchQuiet.Get()
	if quiet {
		noHeader = true
	}

	config := watch.WatchConfig{
		Command:     command,
		Interval:    watchInterval.Get(),
		MaxCount:    watchMaxCount.Get(),
		ExitOnError: watchExitErr.Get(),
		NoHeader:    noHeader,
		NoColor:     watchNoColor.Get() || quiet,
		ClearScreen: true,
		Diff:        diffMode,
		Timeout:     watchTimeout.Get(),
		Precise:     watchPrecise.Get(),
		Quiet:       quiet,
	}

	return watch.WatchCmdMain(config)
}
