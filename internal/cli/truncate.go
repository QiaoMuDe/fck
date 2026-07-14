package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/truncate"
	"gitee.com/MM-Q/go-kit/fs"
	"gitee.com/MM-Q/qflag"
)

var TruncateCmd *qflag.Cmd

var (
	truncateSize    *qflag.SizeFlag   // 设置文件大小
	truncateCreate  *qflag.BoolFlag   // 如果文件不存在则创建
	truncateRef     *qflag.StringFlag // 参考文件大小
	truncateVerbose *qflag.BoolFlag   // 显示操作的文件
)

func init() {
	TruncateCmd = qflag.NewCmd("truncate", "", qflag.ExitOnError)

	truncateSize = TruncateCmd.Size("size", "s", "设置文件大小 (支持单位: B、K、M、G、T)", 0)
	truncateCreate = TruncateCmd.Bool("create", "c", "如果文件不存在则创建", false)
	truncateRef = TruncateCmd.String("reference", "r", "参考文件大小", "")
	truncateVerbose = TruncateCmd.Bool("verbose", "v", "显示操作的文件", false)

	cmdOpts := &qflag.CmdOpts{
		Desc: "文件截断工具",
		Notes: []string{
			"文件大小通过 -s 选项设置, 支持单位: B、K、M、G、T",
			"目标通过位置参数传递, 支持多个文件",
		},
		UseChinese: true,
		Examples: map[string]string{
			"截断文件大小":   fmt.Sprintf("%s truncate -s 100K file.txt", qflag.Root.Name()),
			"参考其他文件大小": fmt.Sprintf("%s truncate -r file.txt", qflag.Root.Name()),
		},
	}

	if err := TruncateCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	TruncateCmd.SetRun(runTruncate)
}

func runTruncate(cmd qflag.Command) error {
	args := cmd.Args()

	// 展开通配符
	targets, err := fs.ExpandFiles(args)
	if err != nil {
		return err
	}

	config := truncate.TruncateConfig{
		Size:      truncateSize.Get(),
		Reference: truncateRef.Get(),
		Targets:   targets,
		Create:    truncateCreate.Get(),
		Verbose:   truncateVerbose.Get(),
	}

	return truncate.TruncateCmdMain(config)
}
