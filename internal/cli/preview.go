package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/preview"
	"gitee.com/MM-Q/qflag"
)

var PreviewCmd *qflag.Cmd

var (
	previewInfo  *qflag.BoolFlag // 打印压缩包信息
	previewLs    *qflag.BoolFlag // 简洁文件列表
	previewLl    *qflag.BoolFlag // 详细文件列表
	previewLimit *qflag.IntFlag  // 限制文件数量
)

func init() {
	PreviewCmd = qflag.NewCmd("preview", "pv", qflag.ExitOnError)

	previewInfo = PreviewCmd.Bool("info", "i", "打印压缩包基本信息", false)
	previewLs = PreviewCmd.Bool("list", "ls", "以简洁的方式打印文件列表", false)
	previewLl = PreviewCmd.Bool("list-long", "ll", "以详细的方式打印文件列表", false)
	previewLimit = PreviewCmd.Int("limit", "l", "限制显示的文件数量(0表示不限制)", 0)

	cmdOpts := &qflag.CmdOpts{
		Desc:       "压缩包预览工具, 查看压缩包信息和文件列表",
		Notes:      []string{"支持的格式有: .zip, .tar, .tar.gz, .tgz, .gz, .bz2, .bzip2, .zlib"},
		UseChinese: true,
	}

	if err := PreviewCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	PreviewCmd.SetRun(runPreview)
}

func runPreview(cmd qflag.Command) error {
	args := cmd.Args()
	if len(args) < 1 {
		return fmt.Errorf("缺少压缩包路径参数")
	}

	packPath := args[0]

	config := preview.PreviewConfig{
		PackPath: packPath,
		Info:     previewInfo.Get(),
		Ls:       previewLs.Get(),
		Ll:       previewLl.Get(),
		Limit:    previewLimit.Get(),
	}

	return preview.PreviewCmdMain(config)
}
