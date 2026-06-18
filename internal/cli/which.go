package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/which"
	"gitee.com/MM-Q/qflag"
)

var WhichCmd *qflag.Cmd

var (
	whichAll    *qflag.BoolFlag // 显示所有匹配的路径
	whichSilent *qflag.BoolFlag // 静默模式，不输出任何内容
)

func init() {
	WhichCmd = qflag.NewCmd("which", "", qflag.ExitOnError)

	whichAll = WhichCmd.Bool("all", "a", "显示所有匹配的路径 (而不仅是第一个)", false)
	whichSilent = WhichCmd.Bool("silent", "s", "静默模式，不输出任何内容，只返回退出码", false)

	cmdOpts := &qflag.CmdOpts{
		Desc: "在环境变量 PATH 中查找命令的可执行文件路径",
		Notes: []string{
			"如果找到命令，输出其完整路径并返回退出码 0",
			"如果未找到命令，返回错误信息并返回退出码 1",
			"使用 -a 选项可以显示所有匹配的路径",
			"使用 -s 选项可以在脚本中静默检查命令是否存在",
		},
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s which [options] <command> [command...]", qflag.Root.Name()),
		Examples: map[string]string{
			"查找单个命令":       fmt.Sprintf("%s which go", qflag.Root.Name()),
			"查找多个命令":       fmt.Sprintf("%s which go git node", qflag.Root.Name()),
			"显示所有匹配路径":     fmt.Sprintf("%s which -a python", qflag.Root.Name()),
			"静默模式检查（用于脚本）": fmt.Sprintf("%s which -s go && echo 'go is installed'", qflag.Root.Name()),
		},
	}

	if err := WhichCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	WhichCmd.SetRun(runWhich)
}

func runWhich(cmd qflag.Command) error {
	config := which.WhichConfig{
		Commands: cmd.Args(),
		All:      whichAll.Get(),
		Silent:   whichSilent.Get(),
	}

	return which.WhichCmdMain(config)
}
