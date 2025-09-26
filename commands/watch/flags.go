// Package watch 定义了 watch 子命令的命令行标志和参数配置。
// 该文件包含所有 watch 命令支持的选项, 如执行间隔、次数限制、颜色输出、超时设置等。
package watch

import (
	"flag"
	"time"

	"gitee.com/MM-Q/qflag"
	"gitee.com/MM-Q/qflag/cmd"
)

var (
	// fck watch 子命令
	watchCmd         *cmd.Cmd
	watchCmdInterval *qflag.DurationFlag // interval 标志
	watchCmdTimes    *qflag.IntFlag      // times 标志
	watchCmdExitErr  *qflag.BoolFlag     // exit-on-error 标志
	watchCmdNoTitle  *qflag.BoolFlag     // no-title 标志
	watchCmdTimeout  *qflag.DurationFlag // timeout 标志
	watchCmdShell    *qflag.EnumFlag     // shell 标志
	watchCmdCommand  *qflag.StringFlag   // command 标志
)

func InitWatchCmd() *cmd.Cmd {
	// fck watch 子命令
	watchCmd = qflag.NewCmd("watch", "w", flag.ExitOnError).
		WithDesc("命令监控工具, 周期性执行指定命令并显示输出结果").
		WithChinese(true)
	watchCmd.AddNote("如果不指定命令, 将提示输入要监控的命令")
	watchCmd.AddNote("使用 Ctrl+C 可以随时停止监控")
	watchCmd.AddNote("命令执行失败时默认继续监控, 除非使用 -e 标志")

	// 添加标志
	watchCmdInterval = watchCmd.Duration("interval", "i", 1*time.Second, "执行间隔时间(秒), 默认1秒")
	watchCmdTimes = watchCmd.Int("times", "n", -1, "执行次数限制, -1表示无限制(默认)")
	watchCmdExitErr = watchCmd.Bool("exit-on-error", "e", false, "命令执行失败时退出")
	watchCmdNoTitle = watchCmd.Bool("no-title", "nt", false, "不显示标题栏")
	watchCmdTimeout = watchCmd.Duration("timeout", "t", 30*time.Second, "单次命令执行超时时间(秒), 默认30秒")
	watchCmdShell = watchCmd.Enum("shell", "s", "def1", "指定使用的shell, 默认使用系统默认shell, 可选值:"+
		"\t\t\t\t\t[def1      ] - 默认值, 使用系统默认shell(win系统默认使用cmd, linux系统默认使用sh)"+
		"\t\t\t\t\t[def2      ] - 使用系统默认shell(win系统默认使用powershell, linux系统默认使用sh)"+
		"\t\t\t\t\t[bash      ] - bash shell"+
		"\t\t\t\t\t[cmd       ] - cmd shell"+
		"\t\t\t\t\t[pwsh      ] - pwsh shell"+
		"\t\t\t\t\t[powershell] - powershell shell"+
		"\t\t\t\t\t[sh        ] - sh shell"+
		"\t\t\t\t\t[none      ] - 不使用shell, 直接执行命令", supportedShells)
	watchCmdCommand = watchCmd.String("command", "cmd", "", "指定要监控执行的命令")

	// 返回子命令
	return watchCmd
}
