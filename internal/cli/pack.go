package cli

import (
	"fmt"
	"path/filepath"

	"gitee.com/MM-Q/fck/internal/commands/pack"
	"gitee.com/MM-Q/fck/internal/types"
	"gitee.com/MM-Q/qflag"
)

var PackCmd *qflag.Cmd

var (
	// 过滤器配置标志
	packIncludePatterns *qflag.StringSliceFlag // 包含模式
	packExcludePatterns *qflag.StringSliceFlag // 排除模式
	packMinSize         *qflag.Int64Flag       // 最小文件大小限制
	packMaxSize         *qflag.Int64Flag       // 最大文件大小限制

	// 压缩配置标志
	packCompressionLevel *qflag.EnumFlag   // 压缩级别
	packOverwrite        *qflag.BoolFlag   // 覆盖已存在文件
	packProgress         *qflag.BoolFlag   // 启用进度显示
	packProgressStyle    *qflag.EnumFlag   // 进度条样式
	packNoValidate       *qflag.BoolFlag   // 禁用路径验证
	packOutput           *qflag.StringFlag // 压缩包输出路径
	packTimestamp        *qflag.BoolFlag   // 添加时间戳到压缩包名
)

func init() {
	PackCmd = qflag.NewCmd("pack", "pk", qflag.ExitOnError)

	packIncludePatterns = PackCmd.StringSlice("include", "i", "包含的文件模式(支持glob语法)", []string{})
	packExcludePatterns = PackCmd.StringSlice("exclude", "e", "排除的文件模式(支持glob语法)", []string{})
	packMinSize = PackCmd.Int64("min-size", "ms", "最小文件大小限制(单位字节, 0表示不限制)", 0)
	packMaxSize = PackCmd.Int64("max-size", "mx", "最大文件大小限制(单位字节, 0表示不限制)", 0)

	packCompressionLevel = PackCmd.Enum("compression", "c", "压缩级别，支持以下选项：\n"+
		"\t\t\t\t\t[default ] - 默认压缩级别\n"+
		"\t\t\t\t\t[none    ] - 不压缩\n"+
		"\t\t\t\t\t[fast    ] - 快速压缩\n"+
		"\t\t\t\t\t[best    ] - 最佳压缩\n"+
		"\t\t\t\t\t[huffman ] - huffman 压缩", types.CompressionLevelDefault, types.SupportedCompressionLevels)
	packOverwrite = PackCmd.Bool("overwrite", "f", "覆盖已存在的压缩文件", false)
	packProgress = PackCmd.Bool("progress", "p", "显示压缩进度", false)
	packProgressStyle = PackCmd.Enum("progress-style", "ps", "进度条样式，支持以下选项：\n"+
		"\t\t\t\t\t[text   ] - 文本样式\n"+
		"\t\t\t\t\t[default] - 默认样式\n"+
		"\t\t\t\t\t[unicode] - unicode 样式\n"+
		"\t\t\t\t\t[ascii  ] - ascii 样式", types.ProgressStyleAscii, types.SupportedProgressStyles)
	packNoValidate = PackCmd.Bool("no-validate", "nv", "禁用路径验证", false)
	packOutput = PackCmd.String("output", "o", "压缩包输出路径，默认为: 源文件名+.zip", "")
	packTimestamp = PackCmd.Bool("timestamp", "t", "在压缩包名中添加时间戳(格式: YYYYMMDDHHMMSS)", false)

	cmdOpts := &qflag.CmdOpts{
		Desc: "智能打包压缩工具",
		Notes: []string{
			"支持的格式有: .zip, .tar, .tar.gz, .tgz, .gz, .bz2, .bzip2, .zlib",
			"通过文件后缀指定压缩格式",
			"不指定 -o 时，自动生成源文件名+.zip",
			"使用 -t 可在压缩包名中添加时间戳，避免覆盖",
		},
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s pack [options] <source>", qflag.Root.Name()),
		Examples: map[string]string{
			"快速打包文件":   fmt.Sprintf("%s pack source.txt", qflag.Root.Name()),
			"快速打包目录":   fmt.Sprintf("%s pack source/", qflag.Root.Name()),
			"指定压缩包名称":  fmt.Sprintf("%s pack -o backup.zip source.txt", qflag.Root.Name()),
			"指定压缩格式":   fmt.Sprintf("%s pack -o backup.tar.gz source/", qflag.Root.Name()),
			"添加时间戳":    fmt.Sprintf("%s pack -t source.txt", qflag.Root.Name()),
			"指定名称+时间戳": fmt.Sprintf("%s pack -t -o backup.zip source.txt", qflag.Root.Name()),
		},
	}

	if err := PackCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	PackCmd.SetRun(runPack)
}

func runPack(cmd qflag.Command) error {
	args := cmd.Args()
	if len(args) < 1 {
		return fmt.Errorf("source path is required")
	}

	srcPath := filepath.Clean(args[0]) // 源路径

	config := pack.PackConfig{
		SrcPath:          srcPath,
		PackPath:         packOutput.Get(),
		IncludePatterns:  packIncludePatterns.Get(),
		ExcludePatterns:  packExcludePatterns.Get(),
		MinSize:          packMinSize.Get(),
		MaxSize:          packMaxSize.Get(),
		CompressionLevel: packCompressionLevel.Get(),
		Overwrite:        packOverwrite.Get(),
		Progress:         packProgress.Get(),
		ProgressStyle:    packProgressStyle.Get(),
		NoValidate:       packNoValidate.Get(),
		Timestamp:        packTimestamp.Get(),
	}

	return pack.PackCmdMain(config)
}
