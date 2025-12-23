package pack

import (
	"flag"
	"fmt"

	"gitee.com/MM-Q/fck/commands/internal/types"
	"gitee.com/MM-Q/qflag"
)

var (
	packCmd *qflag.Cmd // pack 命令

	// 过滤器配置标志
	includePatterns *qflag.StringSliceFlag // 包含模式
	excludePatterns *qflag.StringSliceFlag // 排除模式
	minSize         *qflag.Int64Flag       // 最小文件大小
	maxSize         *qflag.Int64Flag       // 最大文件大小

	// 压缩配置标志
	compressionLevel *qflag.EnumFlag // 压缩级别
	overwrite        *qflag.BoolFlag // 覆盖已存在文件
	progress         *qflag.BoolFlag // 启用进度显示
	progressStyle    *qflag.EnumFlag // 进度条样式
	noValidate       *qflag.BoolFlag // 禁用路径验证
)

func InitPackCmd() *qflag.Cmd {
	packCmd = qflag.NewCmd("pack", "p", flag.ExitOnError)

	packCmdCfg := qflag.CmdConfig{
		UseChinese:  true,
		Desc:        "智能压缩打包工具, 智能识别文件类型并压缩打包",
		Notes:       []string{"支持的格式有: .zip, .tar, .tar.gz, .tgz, .gz, .bz2, .bzip2, .zlib"},
		UsageSyntax: fmt.Sprintf("%s pack [options] <archive> [src]\n", qflag.Root.LongName()),
	}

	packCmd.ApplyConfig(packCmdCfg)

	// 添加过滤器配置标志
	includePatterns = packCmd.StringSlice("include", "i", []string{}, "包含的文件模式(支持glob语法)")
	excludePatterns = packCmd.StringSlice("exclude", "e", []string{}, "排除的文件模式(支持glob语法)")
	minSize = packCmd.Int64("min-size", "ms", 0, "最小文件大小限制(单位字节, 0表示不限制)")
	maxSize = packCmd.Int64("max-size", "mx", 0, "最大文件大小限制(单位字节, 0表示不限制)")

	// 添加压缩配置标志
	compressionLevel = packCmd.Enum("compression", "c", types.CompressionLevelDefault, "压缩级别，支持以下选项：\n"+
		"\t\t\t\t\t[default ] - 默认压缩级别\n"+
		"\t\t\t\t\t[none    ] - 不压缩\n"+
		"\t\t\t\t\t[fast    ] - 快速压缩\n"+
		"\t\t\t\t\t[best    ] - 最佳压缩\n"+
		"\t\t\t\t\t[huffman ] - huffman 压缩", types.SupportedCompressionLevels)
	overwrite = packCmd.Bool("overwrite", "f", false, "覆盖已存在的压缩文件")
	progress = packCmd.Bool("progress", "p", false, "显示压缩进度")
	progressStyle = packCmd.Enum("progress-style", "ps", types.ProgressStyleAscii, "进度条样式，支持以下选项：\n"+
		"\t\t\t\t\t[text   ] - 文本样式\n"+
		"\t\t\t\t\t[default] - 默认样式\n"+
		"\t\t\t\t\t[unicode] - unicode 样式\n"+
		"\t\t\t\t\t[ascii  ] - ascii 样式", types.SupportedProgressStyles)
	noValidate = packCmd.Bool("no-validate", "nv", false, "禁用路径验证")

	return packCmd
}
