package cli

import (
	"fmt"
	"path/filepath"

	"gitee.com/MM-Q/fck/internal/commands/rm"
	"gitee.com/MM-Q/qflag"
)

var RMCmd *qflag.Cmd

var (
	rmRecursive *qflag.BoolFlag // 递归删除目录及其内容
	rmForce     *qflag.BoolFlag // 强制删除（删除只读文件，跳过确认）
	rmVerbose   *qflag.BoolFlag // 显示删除的文件/目录
)

func init() {
	RMCmd = qflag.NewCmd("rm", "", qflag.ExitOnError)

	rmRecursive = RMCmd.Bool("recursive", "r", "递归删除目录及其内容", false)
	rmForce = RMCmd.Bool("force", "f", "强制删除（删除只读文件，跳过确认）", false)
	rmVerbose = RMCmd.Bool("verbose", "v", "显示删除的文件/目录", false)

	cmdOpts := &qflag.CmdOpts{
		Desc: "文件删除工具",
		Notes: []string{
			"目标通过位置参数传递，支持多个文件/目录",
			"默认删除前需要确认，使用 -f 选项强制删除",
			"递归删除时，会删除目录及其所有内容",
			"确认提示选项: y(是), n(否), a(全部是), q(退出)",
		},
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s rm [options] <target>...", qflag.Root.Name()),
		Examples: map[string]string{
			"删除文件": fmt.Sprintf("%s rm file.txt", qflag.Root.Name()),
			"删除目录": fmt.Sprintf("%s rm -r dir/", qflag.Root.Name()),
		},
	}

	if err := RMCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	RMCmd.SetRun(runRM)
}

func runRM(cmd qflag.Command) error {
	args := cmd.Args()
	targets := []string{}

	for _, arg := range args {
		matches, err := filepath.Glob(arg)
		if err != nil {
			return fmt.Errorf("glob pattern expansion failed: %w", err)
		}
		if len(matches) == 0 {
			targets = append(targets, arg)
		} else {
			targets = append(targets, matches...)
		}
	}

	config := rm.RMConfig{
		Targets:   targets,
		Recursive: rmRecursive.Get(),
		Force:     rmForce.Get(),
		Verbose:   rmVerbose.Get(),
	}

	return rm.RMCmdMain(config)
}
