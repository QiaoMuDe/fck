package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/tail"
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
		Desc:       "显示文件结尾内容",
		UseChinese: true,
		Examples: map[string]string{
			"显示后10行":    "fck tail file.txt",
			"显示后20行":    "fck tail -n 20 file.txt",
			"显示后100字节":  "fck tail -c 100 file.txt",
			"实时追踪日志":    "fck tail -f /var/log/app.log",
			"追踪并显示100行": "fck tail -n 100 -f /var/log/app.log",
			"多个文件":      "fck tail file1.txt file2.txt",
			"从管道读取":     "cat big.log | fck tail -n 50",
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
	config := tail.TailConfig{
		Targets: cmd.Args(),
		Lines:   tailLines.Get(),
		Bytes:   tailBytes.Get(),
		Follow:  tailFollow.Get(),
		Quiet:   tailQuiet.Get(),
		Verbose: tailVerbose.Get(),
	}

	return tail.TailCmdMain(config)
}
