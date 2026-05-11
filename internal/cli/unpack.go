package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/unpack"
	"gitee.com/MM-Q/fck/internal/types"
	"gitee.com/MM-Q/qflag"
)

var UnpackCmd *qflag.Cmd

var (
	unpackIncludePatterns *qflag.StringSliceFlag // 包含模式
	unpackExcludePatterns *qflag.StringSliceFlag // 排除模式
	unpackMinSize         *qflag.Int64Flag       // 最小文件大小限制
	unpackMaxSize         *qflag.Int64Flag       // 最大文件大小限制
	unpackOverwrite       *qflag.BoolFlag        // 覆盖已存在文件
	unpackProgress        *qflag.BoolFlag        // 启用进度显示
	unpackProgressStyle   *qflag.EnumFlag        // 进度条样式
	unpackNoValidate      *qflag.BoolFlag        // 禁用路径验证
)

func init() {
	UnpackCmd = qflag.NewCmd("unpack", "upk", qflag.ExitOnError)

	unpackIncludePatterns = UnpackCmd.StringSlice("include", "i", "包含的文件模式(支持glob语法)", []string{})
	unpackExcludePatterns = UnpackCmd.StringSlice("exclude", "e", "排除的文件模式(支持glob语法)", []string{})
	unpackMinSize = UnpackCmd.Int64("min-size", "ms", "最小文件大小限制(单位字节, 0表示不限制)", 0)
	unpackMaxSize = UnpackCmd.Int64("max-size", "mx", "最大文件大小限制(单位字节, 0表示不限制)", 0)

	unpackOverwrite = UnpackCmd.Bool("overwrite", "f", "覆盖已存在的文件", false)
	unpackProgress = UnpackCmd.Bool("progress", "p", "显示解压进度", false)
	unpackProgressStyle = UnpackCmd.Enum("progress-style", "ps", "进度条样式，支持以下选项：\n"+
		"\t\t\t\t\t[text   ] - 文本样式\n"+
		"\t\t\t\t\t[default] - 默认样式\n"+
		"\t\t\t\t\t[unicode] - unicode 样式\n"+
		"\t\t\t\t\t[ascii  ] - ascii 样式", types.ProgressStyleAscii, types.SupportedProgressStyles)
	unpackNoValidate = UnpackCmd.Bool("no-validate", "nv", "禁用路径验证", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "智能解压缩工具",
		Notes:       []string{"支持的格式有: .zip, .tar, .tar.gz, .tgz, .gz, .bz2, .bzip2, .zlib"},
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s unpack [options] <archive> [dst]", qflag.Root.Name()),
		Examples: map[string]string{
			"解压压缩包":      fmt.Sprintf("%s unpack archive.zip", qflag.Root.Name()),
			"解压压缩包到指定目录": fmt.Sprintf("%s unpack archive.zip /path/to/dst", qflag.Root.Name()),
		},
	}

	if err := UnpackCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	UnpackCmd.SetRun(runUnpack)
}

func runUnpack(cmd qflag.Command) error {
	args := cmd.Args()
	if len(args) < 1 {
		return fmt.Errorf("archive path is required")
	}

	packPath := args[0] // 压缩包路径
	dstPath := ""       // 解压目标路径
	if len(args) > 1 {
		dstPath = args[1]
	}

	config := unpack.UnpackConfig{
		PackPath:        packPath,
		DstPath:         dstPath,
		IncludePatterns: unpackIncludePatterns.Get(),
		ExcludePatterns: unpackExcludePatterns.Get(),
		MinSize:         unpackMinSize.Get(),
		MaxSize:         unpackMaxSize.Get(),
		Overwrite:       unpackOverwrite.Get(),
		Progress:        unpackProgress.Get(),
		ProgressStyle:   unpackProgressStyle.Get(),
		NoValidate:      unpackNoValidate.Get(),
	}

	return unpack.UnpackCmdMain(config)
}
