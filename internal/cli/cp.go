package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/cp"
	"gitee.com/MM-Q/qflag"
)

var CpCmd *qflag.Cmd

var (
	cpForce       *qflag.BoolFlag // 强制覆盖
	cpInteractive *qflag.BoolFlag // 交互式覆盖
	cpVerbose     *qflag.BoolFlag // 显示复制的文件/目录
)

func init() {
	CpCmd = qflag.NewCmd("cp", "", qflag.ExitOnError)

	cpForce = CpCmd.Bool("force", "f", "强制覆盖已存在的文件", false)
	cpInteractive = CpCmd.Bool("interactive", "i", "交互式覆盖（覆盖前提示）", false)
	cpVerbose = CpCmd.Bool("verbose", "v", "显示复制的文件/目录", false)

	cmdOpts := &qflag.CmdOpts{
		Desc: "文件复制工具, 支持复制文件和目录, 自动递归复制目录并保留文件属性",
		Notes: []string{
			"源文件通过位置参数传递，支持多个文件",
			"目标通过最后一个位置参数传递",
			"使用 -f 选项强制覆盖已存在的文件",
			"使用 -i 选项在覆盖前提示确认",
			"使用 -v 选项显示复制的文件/目录",
			"自动递归复制目录并保留文件属性（权限、时间戳）",
		},
		UseChinese: true,
	}

	if err := CpCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	CpCmd.SetRun(runCp)
}

func runCp(cmd qflag.Command) error {
	args := cmd.Args()
	if len(args) < 2 {
		return fmt.Errorf("参数不足，至少需要指定源文件和目标路径")
	}

	config := cp.CpConfig{
		Force:       cpForce.Get(),
		Interactive: cpInteractive.Get(),
		Verbose:     cpVerbose.Get(),
		Sources:     args[:len(args)-1],
		Target:      args[len(args)-1],
	}

	return cp.CpCmdMain(config)
}
