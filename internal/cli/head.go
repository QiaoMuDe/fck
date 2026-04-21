package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/head"
	"gitee.com/MM-Q/qflag"
)

var HeadCmd *qflag.Cmd

var (
	headLines   *qflag.IntFlag   // -n, --lines    行数
	headBytes   *qflag.Int64Flag // -c, --bytes    字节数
	headQuiet   *qflag.BoolFlag  // -q, --quiet    不显示文件名
	headVerbose *qflag.BoolFlag  // -v, --verbose  总是显示文件名
)

func init() {
	HeadCmd = qflag.NewCmd("head", "", qflag.ExitOnError)

	headLines = HeadCmd.Int("lines", "n", "显示前N行", 10)
	headBytes = HeadCmd.Int64("bytes", "c", "显示前N字节", 0)
	headQuiet = HeadCmd.Bool("quiet", "q", "不显示文件名标题", false)
	headVerbose = HeadCmd.Bool("verbose", "v", "总是显示文件名标题", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:       "显示文件开头内容",
		UseChinese: true,
		Examples: map[string]string{
			"显示前10行":   fmt.Sprintf("%s head file.txt", qflag.Root.Name()),
			"显示前20行":   fmt.Sprintf("%s head -n 20 file.txt", qflag.Root.Name()),
			"显示前100字节": fmt.Sprintf("%s head -c 100 file.txt", qflag.Root.Name()),
			"多个文件":     fmt.Sprintf("%s head file1.txt file2.txt", qflag.Root.Name()),
			"从管道读取":    fmt.Sprintf("%s cat big.log | %s head -n 100", qflag.Root.Name(), qflag.Root.Name()),
			"简洁模式":     fmt.Sprintf("%s head -q file.txt", qflag.Root.Name()),
		},
		Notes: []string{
			"默认显示前10行",
			"-n 和 -c 不能同时使用",
			"多个文件时自动显示文件名标题",
		},
		MutexGroups: []qflag.MutexGroup{
			{
				Name:      "count-mode",
				Flags:     []string{"lines", "bytes"},
				AllowNone: true,
			},
		},
	}

	if err := HeadCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	HeadCmd.SetRun(runHead)
}

func runHead(cmd qflag.Command) error {
	config := head.HeadConfig{
		Targets: cmd.Args(),
		Lines:   headLines.Get(),
		Bytes:   headBytes.Get(),
		Quiet:   headQuiet.Get(),
		Verbose: headVerbose.Get(),
	}

	return head.HeadCmdMain(config)
}
