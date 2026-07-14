package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/tail"
	"gitee.com/MM-Q/go-kit/fs"
	"gitee.com/MM-Q/qflag"
)

var TailCmd *qflag.Cmd

var (
	tailLines   *qflag.IntFlag   // -n, --lines    行数
	tailBytes   *qflag.Int64Flag // -c, --bytes    字节数
	tailFollow  *qflag.BoolFlag  // -f, --follow   实时追踪
	tailQuiet   *qflag.BoolFlag  // -q, --quiet    不显示文件名
	tailVerbose *qflag.BoolFlag  // -v, --verbose  总是显示文件名
)

func init() {
	TailCmd = qflag.NewCmd("tail", "", qflag.ExitOnError)

	tailLines = TailCmd.Int("lines", "n", "显示后N行", 10)
	tailBytes = TailCmd.Int64("bytes", "c", "显示后N字节", 0)
	tailFollow = TailCmd.Bool("follow", "f", "实时追踪文件变化", false)
	tailQuiet = TailCmd.Bool("quiet", "q", "不显示文件名标题", false)
	tailVerbose = TailCmd.Bool("verbose", "v", "总是显示文件名标题", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "显示文件结尾内容",
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s tail [options] <target>...", qflag.Root.Name()),
		Examples: map[string]string{
			"显示后10行":    fmt.Sprintf("%s tail file.txt", qflag.Root.Name()),
			"显示后20行":    fmt.Sprintf("%s tail -n 20 file.txt", qflag.Root.Name()),
			"显示后100字节":  fmt.Sprintf("%s tail -c 100 file.txt", qflag.Root.Name()),
			"实时追踪日志":    fmt.Sprintf("%s tail -f /var/log/app.log", qflag.Root.Name()),
			"追踪并显示100行": fmt.Sprintf("%s tail -n 100 -f /var/log/app.log", qflag.Root.Name()),
			"多个文件":      fmt.Sprintf("%s tail file1.txt file2.txt", qflag.Root.Name()),
			"从管道读取":     fmt.Sprintf("cat big.log | %s tail -n 50", qflag.Root.Name()),
		},
		Notes: []string{
			"默认显示后10行",
			"-n 和 -c 不能同时使用",
			"多个文件时自动显示文件名标题",
			"-f 模式会持续运行，按 Ctrl+C 退出",
		},
		MutexGroups: []qflag.MutexGroup{
			{
				Name:      "count-mode",
				Flags:     []string{"lines", "bytes"},
				AllowNone: true,
			},
		},
	}

	if err := TailCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	TailCmd.SetRun(runTail)
}

func runTail(cmd qflag.Command) error {
	args := cmd.Args()

	// 展开通配符，只保留文件
	targets, err := fs.ExpandFiles(args)
	if err != nil {
		return err
	}

	config := tail.TailConfig{
		Targets: targets,
		Lines:   tailLines.Get(),
		Bytes:   tailBytes.Get(),
		Follow:  tailFollow.Get(),
		Quiet:   tailQuiet.Get(),
		Verbose: tailVerbose.Get(),
	}

	return tail.TailCmdMain(config)
}
