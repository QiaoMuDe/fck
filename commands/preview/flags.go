package preview

import (
	"flag"
	"fmt"

	"gitee.com/MM-Q/qflag"
	"gitee.com/MM-Q/qflag/cmd"
)

var (
	previewCmd *cmd.Cmd // Preview command

	// 预览配置标志
	infoFlag  *qflag.BoolFlag // 打印压缩包信息
	lsFlag    *qflag.BoolFlag // 简洁文件列表
	llFlag    *qflag.BoolFlag // 详细文件列表
	limitFlag *qflag.IntFlag  // 限制文件数量
)

// 初始化预览命令
func InitPreviewCmd() *cmd.Cmd {
	previewCmd = cmd.NewCmd("preview", "pv", flag.ExitOnError).
		WithChinese(true).
		WithUsage(fmt.Sprint(qflag.Root.LongName(), " preview [options] <archive>")).
		WithDesc("压缩包预览工具, 查看压缩包信息和文件列表")
	previewCmd.AddNote("支持的格式有: .zip, .tar, .tar.gz, .tgz, .gz, .bz2, .bzip2, .zlib")

	// 添加预览配置标志
	infoFlag = previewCmd.Bool("info", "i", false, "打印压缩包基本信息")
	lsFlag = previewCmd.Bool("list", "ls", false, "以简洁的方式打印文件列表")
	llFlag = previewCmd.Bool("list-long", "ll", false, "以详细的方式打印文件列表")
	limitFlag = previewCmd.Int("limit", "l", 0, "限制显示的文件数量(0表示不限制)")

	return previewCmd
}
