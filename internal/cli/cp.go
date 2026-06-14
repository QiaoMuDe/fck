package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/cp"
	"gitee.com/MM-Q/go-kit/fs"
	"gitee.com/MM-Q/qflag"
)

var CpCmd *qflag.Cmd

var (
	cpForce       *qflag.BoolFlag // 强制覆盖
	cpInteractive *qflag.BoolFlag // 交互式覆盖
	cpVerbose     *qflag.BoolFlag // 显示复制的文件/目录
	cpRecursive   *qflag.BoolFlag // 递归复制目录
)

func init() {
	CpCmd = qflag.NewCmd("cp", "", qflag.ExitOnError)

	cpForce = CpCmd.Bool("force", "f", "强制覆盖已存在的文件", false)
	cpInteractive = CpCmd.Bool("interactive", "i", "交互式覆盖（覆盖前提示）", false)
	cpVerbose = CpCmd.Bool("verbose", "v", "显示复制的文件/目录", false)
	cpRecursive = CpCmd.Bool("recursive", "r", "递归复制目录", false)

	cmdOpts := &qflag.CmdOpts{
		Desc: "文件目录复制工具",
		Notes: []string{
			"源文件通过位置参数传递，支持多个文件",
			"目标通过最后一个位置参数传递",
			"复制目录需要使用 -r 选项递归复制",
			"自动保留文件属性 (权限、时间戳)",
		},
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s cp [options] <source>... <target>", qflag.Root.Name()),
		Examples: map[string]string{
			"复制文件":   fmt.Sprintf("%s cp file.txt /path/dest.txt", qflag.Root.Name()),
			"复制目录":   fmt.Sprintf("%s cp -r src/ dest/", qflag.Root.Name()),
			"强制覆盖":   fmt.Sprintf("%s cp -f file.txt /path/dest.txt", qflag.Root.Name()),
			"交互式覆盖":  fmt.Sprintf("%s cp -i file.txt /path/dest.txt", qflag.Root.Name()),
			"显示详细信息": fmt.Sprintf("%s cp -v file.txt /path/dest.txt", qflag.Root.Name()),
		},
	}

	if err := CpCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	CpCmd.SetRun(runCp)
}

func runCp(cmd qflag.Command) error {
	args := cmd.Args()
	if len(args) < 2 {
		return fmt.Errorf("insufficient arguments: at least source and destination required")
	}

	sources := args[:len(args)-1]

	// 展开通配符
	expanded, err := fs.Expand(sources)
	if err != nil {
		return err
	}

	config := cp.CpConfig{
		Force:       cpForce.Get(),
		Interactive: cpInteractive.Get(),
		Verbose:     cpVerbose.Get(),
		Recursive:   cpRecursive.Get(),
		Sources:     expanded,
		Target:      args[len(args)-1], // 目标路径（最后一个参数）
	}

	return cp.CpCmdMain(config)
}
