package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/md"
	"gitee.com/MM-Q/fck/internal/utils"
	"gitee.com/MM-Q/qflag"
)

var MdCmd *qflag.Cmd

var (
	mdLess    *qflag.BoolFlag // -l, --less 使用分页器查看
	mdRaw     *qflag.BoolFlag // -r, --raw 分页器中显示原始文件
	mdStyle   *qflag.EnumFlag // -s, --style 渲染样式
	mdWidth   *qflag.IntFlag  // -w, --width 换行宽度
	mdMaxSize *qflag.SizeFlag // -S, --max-size 最大文件大小 (默认100MB)
)

func init() {
	MdCmd = qflag.NewCmd("md", "", qflag.ExitOnError)

	mdLess = MdCmd.Bool("less", "l", "使用分页器查看 (默认直接输出到终端)", false)
	mdRaw = MdCmd.Bool("raw", "r", "分页器中同时加载原始文件 (按 ] 切换)", false)
	mdStyle = MdCmd.Enum("style", "s", "渲染样式 (auto/ascii/dark/light/dracula/tokyo-night/pink/notty)", "auto", []string{"auto", "ascii", "dark", "light", "dracula", "tokyo-night", "pink", "notty"})
	mdWidth = MdCmd.Int("width", "w", "换行宽度 (0 表示自动)", utils.GetSafeTerminalWidth())
	mdMaxSize = MdCmd.Size("max-size", "S", "最大文件大小 (如: 50m, 100m, 1g)", 100*1024*1024)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "预览 Markdown 文件",
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s md [options] <path>", qflag.Root.Name()),
		Examples: map[string]string{
			"预览文件":     fmt.Sprintf("%s md README.md", qflag.Root.Name()),
			"使用分页器":    fmt.Sprintf("%s md -l README.md", qflag.Root.Name()),
			"指定样式":     fmt.Sprintf("%s md -s dark README.md", qflag.Root.Name()),
			"指定宽度":     fmt.Sprintf("%s md -w 100 README.md", qflag.Root.Name()),
			"分页器+暗色主题": fmt.Sprintf("%s md -l -s dark README.md", qflag.Root.Name()),
			"管道输入":     fmt.Sprintf("echo '# Hello' | %s md", qflag.Root.Name()),
			"管道+分页器":   fmt.Sprintf("cat README.md | %s md -l", qflag.Root.Name()),
		},
		Notes: []string{
			"支持 glamour 的所有内置样式: auto, ascii, dark, light, dracula, tokyo-night, pink, notty",
			"默认使用 auto 样式，根据终端背景自动选择",
			"默认直接输出到终端，使用 -l 启用分页器查看",
			"使用 -r 标志在分页器中同时加载原始文件，按 ] 切换视图",
			"支持从管道读取内容，此时不能指定文件参数",
			"文件大小超过限制时，提示使用 -l 标志通过分页器查看",
		},
	}

	if err := MdCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	MdCmd.SetRun(runMd)
}

func runMd(cmd qflag.Command) error {
	file := cmd.Arg(0)

	// 检测管道输入
	isPipe := utils.IsStdinPipe()

	if isPipe {
		// 管道模式：不能指定文件参数
		if file != "" {
			return fmt.Errorf("cannot specify file when reading from pipe")
		}
		file = "stdin"
	} else {
		// 文件模式：必须指定文件
		if file == "" {
			return fmt.Errorf("no markdown file specified")
		}
	}

	config := md.MdConfig{
		File:      file,
		UsePipe:   isPipe,
		UsePager:  mdLess.Get(),
		ShowRaw:   mdRaw.Get(),
		Style:     mdStyle.Get(),
		WordWidth: mdWidth.Get(),
		MaxSize:   mdMaxSize.Get(),
	}

	return md.MdCmdMain(config)
}
