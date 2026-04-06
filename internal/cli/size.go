package cli

import (
	"fmt"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/internal/commands/size"
	"gitee.com/MM-Q/fck/internal/types"
	"gitee.com/MM-Q/qflag"
)

var SizeCmd *qflag.Cmd

var (
	sizeColor      *qflag.BoolFlag // 启用颜色输出
	sizeTableStyle *qflag.EnumFlag // 指定表格样式
	sizeHidden     *qflag.BoolFlag // 包含隐藏文件或目录进行大小计算，默认过滤
	sizeHuman      *qflag.BoolFlag // 人类可读格式显示大小
)

func init() {
	SizeCmd = qflag.NewCmd("size", "s", qflag.ExitOnError)

	sizeColor = SizeCmd.Bool("color", "c", "启用颜色输出", false)
	sizeHidden = SizeCmd.Bool("hidden", "H", "包含隐藏文件或目录进行大小计算，默认过滤", false)
	sizeHuman = SizeCmd.Bool("human", "u", "以人类可读格式显示大小(如KB/MB/GB)", false)
	sizeTableStyle = SizeCmd.Enum("table-style", "ts", "指定表格样式，支持以下选项：\n"+
		"\t\t\t\t\t[def ]   - 默认样式\n"+
		"\t\t\t\t\t[l   ]   - 浅色样式\n"+
		"\t\t\t\t\t[r   ]   - 圆角样式\n"+
		"\t\t\t\t\t[bd  ]   - 粗体样式\n"+
		"\t\t\t\t\t[cb  ]   - 亮色彩色样式\n"+
		"\t\t\t\t\t[cd  ]   - 暗色彩色样式\n"+
		"\t\t\t\t\t[db  ]   - 双线样式\n"+
		"\t\t\t\t\t[cbb ]   - 黑色背景蓝色字体\n"+
		"\t\t\t\t\t[cbc ]   - 青色背景蓝色字体\n"+
		"\t\t\t\t\t[cbg ]   - 绿色背景蓝色字体\n"+
		"\t\t\t\t\t[cbm ]   - 紫色背景蓝色字体\n"+
		"\t\t\t\t\t[cby ]   - 黄色背景蓝色字体\n"+
		"\t\t\t\t\t[cbr ]   - 红色背景蓝色字体\n"+
		"\t\t\t\t\t[cwb ]   - 蓝色背景白色字体\n"+
		"\t\t\t\t\t[ccw ]   - 青色背景白色字体\n"+
		"\t\t\t\t\t[cgw ]   - 绿色背景白色字体\n"+
		"\t\t\t\t\t[cmw ]   - 紫色背景白色字体\n"+
		"\t\t\t\t\t[crw ]   - 红色背景白色字体\n"+
		"\t\t\t\t\t[cyw ]   - 黄色背景白色字体\n"+
		"\t\t\t\t\t[none]   - 禁用表格样式", "def", types.TableStyles)

	cmdOpts := &qflag.CmdOpts{
		Desc:       "文件目录大小计算工具",
		Notes:      []string{"默认显示字节数，使用 -u/--human 转换为可读格式"},
		UseChinese: true,
	}

	if err := SizeCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	SizeCmd.SetRun(runSize)
}

func runSize(cmd qflag.Command) error {
	cl := colorlib.NewColorLib()

	config := size.SizeConfig{
		Args:       cmd.Args(),
		Color:      sizeColor.Get(),
		TableStyle: sizeTableStyle.Get(),
		Hidden:     sizeHidden.Get(),
		Human:      sizeHuman.Get(),
	}

	return size.SizeCmdMain(cl, config)
}
