// Package size 定义了 size 子命令的命令行标志和参数配置。
// 该文件包含 size 命令支持的所有选项，如颜色输出、表格样式、隐藏文件处理等。
package size

import (
	"flag"
	"fmt"

	"gitee.com/MM-Q/fck/commands/internal/types"
	"gitee.com/MM-Q/qflag"
	"gitee.com/MM-Q/qflag/cmd"
)

var (
	// fck size 子命令
	sizeCmd           *cmd.Cmd
	sizeCmdColor      *qflag.BoolFlag // color 标志
	sizeCmdTableStyle *qflag.EnumFlag // ts 标志
	sizeCmdHidden     *qflag.BoolFlag // hidden 标志
)

// 初始化
func InitSizeCmd() *cmd.Cmd {
	// fck size 子命令
	sizeCmd = qflag.NewCmd("size", "s", flag.ExitOnError).
		WithUsage(fmt.Sprint(qflag.Root.LongName(), " size [options] <path>...\n")).
		WithChinese(true).
		WithNote("大小单位会自动选择最合适的(B/KB/MB/GB/TB)").
		WithDesc("文件目录大小计算工具, 计算指定文件或目录的大小，并以人类可读格式(B/KB/MB/GB/TB)显示")

	// 标志定义
	sizeCmdColor = sizeCmd.Bool("color", "c", false, "启用颜色输出")
	sizeCmdHidden = sizeCmd.Bool("hidden", "H", false, "包含隐藏文件或目录进行大小计算，默认过滤")
	sizeCmdTableStyle = sizeCmd.Enum("table-style", "ts", "def", "指定表格样式，支持以下选项：\n"+
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
		"\t\t\t\t\t[none]   - 禁用表格样式", types.TableStyles)

	// 返回子命令
	return sizeCmd
}
