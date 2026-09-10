package cli

import (
	"fmt"

	"gitee.com/MM-Q/color"
	"gitee.com/MM-Q/fck/internal/commands/size"
	"gitee.com/MM-Q/fck/internal/types"
	"gitee.com/MM-Q/go-kit/fs"
	"gitee.com/MM-Q/qflag"
)

var SizeCmd *qflag.Cmd

var (
	sizeNoColor       *qflag.BoolFlag // 禁用颜色输出
	sizeTableStyle    *qflag.EnumFlag // 指定表格样式
	sizeHidden        *qflag.BoolFlag // 包含隐藏文件或目录进行大小计算，默认过滤
	sizeHuman         *qflag.BoolFlag // 人类可读格式显示大小
	sizeFollowSymlink *qflag.BoolFlag // 跟随符号链接
)

func init() {
	SizeCmd = qflag.NewCmd("size", "", qflag.ExitOnError)

	sizeNoColor = SizeCmd.Bool("no-color", "n", "禁用颜色输出", false)
	sizeHidden = SizeCmd.Bool("hidden", "H", "包含隐藏文件或目录进行大小计算，默认过滤", false)
	sizeHuman = SizeCmd.Bool("human", "u", "以人类可读格式显示大小(如KB/MB/GB)", false)
	sizeFollowSymlink = SizeCmd.Bool("follow-symlinks", "L", "跟随符号链接计算目标大小", false)
	sizeTableStyle = SizeCmd.Enum("table-style", "ts", qflag.EnumHelp("指定表格样式，支持以下选项:", types.TableStyleOptions, "\t\t\t\t\t"), "def", types.TableStyles)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "文件目录大小计算工具",
		Notes:       []string{"默认显示字节数，使用 -u/--human 转换为可读格式", "默认不跟随符号链接，使用 -L/--follow-symlinks 跟随符号链接"},
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s size [options] <target>...", qflag.Root.Name()),
		Examples: map[string]string{
			"查看当前目录": fmt.Sprintf("%s size", qflag.Root.Name()),
			"查看指定目录": fmt.Sprintf("%s size /path/to/dir", qflag.Root.Name()),
			"人类可读格式": fmt.Sprintf("%s size -u /path/to/dir", qflag.Root.Name()),
			"包含隐藏文件": fmt.Sprintf("%s size -H /path/to/dir", qflag.Root.Name()),
			"跟随符号链接": fmt.Sprintf("%s size -L /path/to/dir", qflag.Root.Name()),
			"使用通配符":  fmt.Sprintf("%s size *.txt", qflag.Root.Name()),
		},
	}

	if err := SizeCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	SizeCmd.SetRun(runSize)
}

func runSize(cmd qflag.Command) error {
	cl := color.G()

	args := cmd.Args()

	// 展开通配符
	targets, err := fs.ExpandFiles(args)
	if err != nil {
		return err
	}

	config := size.SizeConfig{
		Args:          targets,
		NoColor:       sizeNoColor.Get(),
		TableStyle:    sizeTableStyle.Get(),
		Hidden:        sizeHidden.Get(),
		Human:         sizeHuman.Get(),
		FollowSymlink: sizeFollowSymlink.Get(),
	}

	return size.SizeCmdMain(cl, config)
}
