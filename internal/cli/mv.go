package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/mv"
	"gitee.com/MM-Q/qflag"
)

var MvCmd *qflag.Cmd

var (
	mvForce       *qflag.BoolFlag // 强制覆盖
	mvInteractive *qflag.BoolFlag // 交互式覆盖
	mvVerbose     *qflag.BoolFlag // 显示移动的文件/目录
)

func init() {
	MvCmd = qflag.NewCmd("mv", "", qflag.ExitOnError)

	mvForce = MvCmd.Bool("force", "f", "强制覆盖已存在的文件", false)
	mvInteractive = MvCmd.Bool("interactive", "i", "交互式覆盖（覆盖前提示）", false)
	mvVerbose = MvCmd.Bool("verbose", "v", "显示移动的文件/目录", false)

	cmdOpts := &qflag.CmdOpts{
		Desc: "文件移动工具, 支持移动文件和目录, 支持文件重命名",
		Notes: []string{
			"源文件通过位置参数传递，支持多个文件",
			"目标通过最后一个位置参数传递",
			"使用 -f 选项强制覆盖已存在的文件",
			"使用 -i 选项在覆盖前提示确认",
			"使用 -v 选项显示移动的文件/目录",
			"支持跨设备移动（自动处理）",
		},
		UseChinese: true,
	}

	if err := MvCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	MvCmd.SetRun(runMv)
}

func runMv(cmd qflag.Command) error {
	args := cmd.Args()
	if len(args) < 2 {
		return fmt.Errorf("参数不足，至少需要指定源文件和目标路径")
	}

	config := mv.MvConfig{
		Force:       mvForce.Get(),
		Interactive: mvInteractive.Get(),
		Verbose:     mvVerbose.Get(),
		Sources:     args[:len(args)-1],
		Target:      args[len(args)-1],
	}

	return mv.MvCmdMain(config)
}
