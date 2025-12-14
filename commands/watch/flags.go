// Package watch 定义了 watch 子命令的命令行标志和参数配置。
// 该文件包含所有 watch 命令支持的选项, 如执行间隔、次数限制、颜色输出、超时设置等。
package watch

import (
	"flag"
	"fmt"
	"time"

	"gitee.com/MM-Q/qflag"
	"gitee.com/MM-Q/qflag/cmd"
)

var (
	// fck watch 子命令
	watchCmd           *cmd.Cmd
	watchCmdInterval   *qflag.DurationFlag // interval 标志
	watchCmdMaxCount   *qflag.IntFlag      // count 标志
	watchCmdExitErr    *qflag.BoolFlag     // exit-on-error 标志
	watchCmdNoHeader   *qflag.BoolFlag     // no-header 标志
	watchCmdTimeout    *qflag.DurationFlag // timeout 标志
	watchCmdShell      *qflag.EnumFlag     // shell 标志
	watchCmdClearLines *qflag.IntFlag      // clear-line 标志
	watchCmdQuiet      *qflag.BoolFlag     // quiet 标志
)

func InitWatchCmd() *cmd.Cmd {
	// fck watch 子命令
	watchCmd = qflag.NewCmd("watch", "w", flag.ExitOnError).
		WithDesc("命令监控工具, 周期性执行指定命令并显示输出结果").
		WithChinese(true).
		WithUsage(fmt.Sprint(qflag.Root.LongName(), " watch [options] command\n"))
	watchCmd.AddNote("如果不指定命令, 将提示输入要监控的命令")
	watchCmd.AddNote("使用 Ctrl+C 可以随时停止监控")
	watchCmd.AddNote("命令执行失败时默认继续监控, 除非使用 -e 标志")

	// 添加标志
	watchCmdInterval = watchCmd.Duration("interval", "i", 1*time.Second, "执行间隔时间(秒), 默认1秒")
	watchCmdMaxCount = watchCmd.Int("count", "n", -1, "执行次数限制, -1表示无限制(默认)")
	watchCmdExitErr = watchCmd.Bool("exit-on-error", "e", false, "命令执行失败时退出")
	watchCmdNoHeader = watchCmd.Bool("no-header", "nh", false, "轻度静默模式, 不显示标题栏和换行符, 但显示命令输出")
	watchCmdTimeout = watchCmd.Duration("timeout", "t", 30*time.Second, "单次命令执行超时时间(秒), 默认30秒")
	watchCmdShell = watchCmd.Enum("shell", "s", "def1", "指定使用的shell, 默认使用系统默认shell, 可选值:\n"+
		"\t\t\t\t   [def1      ] - 默认值, 使用系统默认shell(win系统默认使用cmd, linux系统默认使用sh)\n"+
		"\t\t\t\t   [def2      ] - 使用系统默认shell(win系统默认使用powershell, linux系统默认使用sh)\n"+
		"\t\t\t\t   [bash      ] - bash shell\n"+
		"\t\t\t\t   [cmd       ] - cmd shell\n"+
		"\t\t\t\t   [pwsh      ] - pwsh shell\n"+
		"\t\t\t\t   [powershell] - powershell shell\n"+
		"\t\t\t\t   [sh        ] - sh shell\n"+
		"\t\t\t\t   [none      ] - 不使用shell, 直接执行命令", supportedShells)
	watchCmdClearLines = watchCmd.Int("clear-line", "cl", 20, "每次执行前打印指定数量的换行符进行清屏, 0表示不清屏(默认)")
	watchCmdQuiet = watchCmd.Bool("quiet", "q", false, "完全静默模式, 不显示标题栏、换行符和命令输出")
	// 返回子命令
	return watchCmd
}
