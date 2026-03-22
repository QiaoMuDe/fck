package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/chmod"
	"gitee.com/MM-Q/qflag"
)

var ChmodCmd *qflag.Cmd

var (
	chmodRecursive *qflag.BoolFlag // 递归修改目录及其内容
	chmodVerbose   *qflag.BoolFlag // 显示修改的文件/目录
	chmodChanges   *qflag.BoolFlag // 只显示有变化的文件/目录
	chmodQuiet     *qflag.BoolFlag // 抑制错误信息
)

func init() {
	ChmodCmd = qflag.NewCmd("chmod", "", qflag.ExitOnError)

	chmodRecursive = ChmodCmd.Bool("recursive", "R", "递归修改目录及其内容", false)
	chmodVerbose = ChmodCmd.Bool("verbose", "v", "显示修改的文件/目录", false)
	chmodChanges = ChmodCmd.Bool("changes", "c", "只显示有变化的文件/目录", false)
	chmodQuiet = ChmodCmd.Bool("quiet", "f", "抑制错误信息", false)

	cmdOpts := &qflag.CmdOpts{
		Desc: "文件权限修改工具, 支持递归修改、权限模式设置、详细输出等功能",
		Notes: []string{
			"权限模式通过第一个位置参数传递，目标通过后续位置参数传递",
			"权限格式为八进制，如 755、644、777 等",
			"使用 -R 选项递归修改目录及其所有内容",
			"使用 -v 选项显示修改的文件/目录",
			"使用 -c 选项只显示有变化的文件/目录",
			"使用 -f 选项抑制错误信息",
		},
		UseChinese: true,
	}

	if err := ChmodCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	ChmodCmd.SetRun(runChmod)
}

func runChmod(cmd qflag.Command) error {
	args := cmd.Args()
	if len(args) == 0 {
		return fmt.Errorf("未指定权限模式")
	}

	config := chmod.ChmodConfig{
		Mode:      args[0],
		Targets:   args[1:],
		Recursive: chmodRecursive.Get(),
		Verbose:   chmodVerbose.Get(),
		Changes:   chmodChanges.Get(),
		Quiet:     chmodQuiet.Get(),
	}

	return chmod.ChmodCmdMain(config)
}
