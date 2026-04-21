package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/mkdir"
	"gitee.com/MM-Q/qflag"
)

var MkdirCmd *qflag.Cmd

var (
	mkdirParents *qflag.BoolFlag   // 递归创建父目录
	mkdirMode    *qflag.StringFlag // 设置目录权限
	mkdirVerbose *qflag.BoolFlag   // 显示创建的目录
)

func init() {
	MkdirCmd = qflag.NewCmd("mkdir", "", qflag.ExitOnError)

	mkdirParents = MkdirCmd.Bool("parents", "p", "递归创建父目录", false)
	mkdirMode = MkdirCmd.String("mode", "m", "设置目录权限 (八进制，如 755)", "0755")
	mkdirVerbose = MkdirCmd.Bool("verbose", "v", "显示创建的目录", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "目录创建工具",
		UsageSyntax: fmt.Sprintf("%s mkdir [options] <path...>", qflag.Root.Name()),
		Notes: []string{
			"目标通过位置参数传递，支持多个目录",
			"默认权限为 0755, 可使用 -m 选项自定义",
			"使用 -p 选项递归创建父目录，跳过已存在的目录",
			"权限格式为八进制，如 755、777、700 等",
		},
		UseChinese: true,
		Examples: map[string]string{
			"创建目录":         fmt.Sprintf("%s mkdir dir", qflag.Root.Name()),
			"创建目录并递归创建父目录": fmt.Sprintf("%s mkdir -p dir/subdir", qflag.Root.Name()),
		},
	}

	if err := MkdirCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	MkdirCmd.SetRun(runMkdir)
}

func runMkdir(cmd qflag.Command) error {
	mode, err := mkdir.ParseMode(mkdirMode.Get())
	if err != nil {
		return err
	}

	config := mkdir.MkdirConfig{
		Targets: cmd.Args(),
		Parents: mkdirParents.Get(),
		Mode:    mode,
		Verbose: mkdirVerbose.Get(),
	}

	return mkdir.MkdirCmdMain(config)
}
