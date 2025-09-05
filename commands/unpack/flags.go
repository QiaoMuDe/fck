package unpack

import (
	"flag"
	"fmt"

	"gitee.com/MM-Q/fck/commands/internal/types"
	"gitee.com/MM-Q/qflag"
	"gitee.com/MM-Q/qflag/cmd"
)

var (
	unpackCmd *cmd.Cmd // unpack 命令

	// 过滤器配置标志
	includePatterns *qflag.StringSliceFlag // 包含模式
	excludePatterns *qflag.StringSliceFlag // 排除模式
	minSize         *qflag.Int64Flag       // 最小文件大小
	maxSize         *qflag.Int64Flag       // 最大文件大小

	// 压缩配置标志
	overwrite     *qflag.BoolFlag // 覆盖已存在文件
	progress      *qflag.BoolFlag // 启用进度显示
	progressStyle *qflag.EnumFlag // 进度条样式
	noValidate    *qflag.BoolFlag // 禁用路径验证
)

// InitUnpackCmd 初始化unpack命令及其所有标志
func InitUnpackCmd() *cmd.Cmd {
	unpackCmd = qflag.NewCmd("unpack", "up", flag.ExitOnError).
		WithUseChinese(true).
		WithUsageSyntax(fmt.Sprint(qflag.LongName(), " unpack [options] <archive> [dst]")).
		WithDescription("智能解压缩工具, 智能识别压缩文件格式并解压")
	unpackCmd.AddNote("支持的格式有: .zip, .tar, .tar.gz, .tgz, .gz, .bz2, .bzip2, .zlib")

	// 添加过滤器配置标志
	includePatterns = unpackCmd.StringSlice("include", "i", []string{}, "包含的文件模式(支持glob语法)")
	excludePatterns = unpackCmd.StringSlice("exclude", "e", []string{}, "排除的文件模式(支持glob语法)")
	minSize = unpackCmd.Int64("min-size", "ms", 0, "最小文件大小限制(单位字节, 0表示不限制)")
	maxSize = unpackCmd.Int64("max-size", "mx", 0, "最大文件大小限制(单位字节, 0表示不限制)")

	// 添加压缩配置标志
	overwrite = unpackCmd.Bool("overwrite", "f", false, "覆盖已存在的文件")
	progress = unpackCmd.Bool("progress", "p", false, "显示解压进度")
	progressStyle = unpackCmd.Enum("progress-style", "ps", types.ProgressStyleAscii, "进度条样式，支持以下选项：\n"+
		"\t\t\t\t\t[text   ] - 文本样式\n"+
		"\t\t\t\t\t[default] - 默认样式\n"+
		"\t\t\t\t\t[unicode] - unicode 样式\n"+
		"\t\t\t\t\t[ascii  ] - ascii 样式", types.SupportedProgressStyles)
	noValidate = unpackCmd.Bool("no-validate", "nv", false, "禁用路径验证")

	return unpackCmd
}
