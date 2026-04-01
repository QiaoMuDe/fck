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
	watchInterval   *qflag.DurationFlag
	watchMaxCount   *qflag.IntFlag
	watchExitErr    *qflag.BoolFlag
	watchNoHeader   *qflag.BoolFlag
	watchTimeout    *qflag.DurationFlag
	watchClearLines *qflag.IntFlag
	watchQuiet      *qflag.BoolFlag
)

func init() {
	WatchCmd = qflag.NewCmd("watch", "w", qflag.ExitOnError)

	watchInterval = WatchCmd.Duration("interval", "i", "执行间隔时间(秒), 默认1秒", 1*time.Second)
	watchMaxCount = WatchCmd.Int("count", "n", "执行次数限制, -1表示无限制(默认)", -1)
	watchExitErr = WatchCmd.Bool("exit-on-error", "e", "命令执行失败时退出", false)
	watchNoHeader = WatchCmd.Bool("no-header", "nh", "轻度静默模式, 不显示标题栏和换行符, 但显示命令输出", false)
	watchTimeout = WatchCmd.Duration("timeout", "t", "单次命令执行超时时间(秒), 默认30秒", 30*time.Second)
	watchClearLines = WatchCmd.Int("clear-line", "cl", "每次执行前打印指定数量的换行符进行清屏, 0表示不清屏(默认)", 20)
	watchQuiet = WatchCmd.Bool("quiet", "q", "完全静默模式, 不显示标题栏、换行符和命令输出", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:       "命令监控工具, 周期性执行指定命令并显示输出结果",
		Notes:      []string{"如果不指定命令, 将提示输入要监控的命令", "使用 Ctrl+C 可以随时停止监控", "命令执行失败时默认继续监控, 除非使用 -e 标志"},
		UseChinese: true,
	}

	if err := WatchCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	WatchCmd.SetRun(runWatch)
}

func runWatch(cmd qflag.Command) error {
	args := cmd.Args()
	if len(args) < 1 {
		return fmt.Errorf("缺少要监控的命令参数")
	}

	var command string
	if len(args) == 1 {
		command = args[0]
	} else {
		command = strings.Join(args, " ")
	}

	config := watch.WatchConfig{
		Command:     command,
		Interval:    watchInterval.Get(),
		MaxCount:    watchMaxCount.Get(),
		ExitOnError: watchExitErr.Get(),
		NoHeader:    watchNoHeader.Get(),
		Timeout:     watchTimeout.Get(),
		ClearLines:  watchClearLines.Get(),
		Quiet:       watchQuiet.Get(),
	}

	return watch.WatchCmdMain(config)
}
