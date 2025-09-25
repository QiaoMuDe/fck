// Package list 定义了 list 子命令的命令行标志和参数配置。
// 该文件包含所有 list 命令支持的选项，如排序方式、显示格式、过滤条件等。
package list

import (
	"flag"
	"fmt"

	"gitee.com/MM-Q/fck/commands/internal/types"
	"gitee.com/MM-Q/qflag"
	"gitee.com/MM-Q/qflag/cmd"
)

var (
	// fck list 子命令
	listCmd              *cmd.Cmd
	listCmdAll           *qflag.BoolFlag // all 标志
	listCmdColor         *qflag.BoolFlag // color 标志
	listCmdSortByName    *qflag.BoolFlag // sort-by-name 标志
	listCmdSortBySize    *qflag.BoolFlag // sort-by-size 标志
	listCmdSortByTime    *qflag.BoolFlag // sort-by-time 标志
	listCmdDirItself     *qflag.BoolFlag // D 标志
	listCmdLongFormat    *qflag.BoolFlag // l 标志
	listCmdReverseSort   *qflag.BoolFlag // r 标志
	listCmdQuoteNames    *qflag.BoolFlag // q 标志
	listCmdRecursion     *qflag.BoolFlag // R 标志
	listCmdShowUserGroup *qflag.BoolFlag // u 标志
	listCmdTableStyle    *qflag.EnumFlag // ts 标志
	listCmdType          *qflag.EnumFlag // type 标志
)

func InitListCmd() *cmd.Cmd {
	// fck list 子命令
	listCmd = qflag.NewCmd("list", "ls", flag.ExitOnError).
		WithDesc("文件目录列表工具, 列出指定目录中的文件和目录，并支持多种排序和过滤选项").
		WithChinese(true).
		WithUsage(fmt.Sprint(qflag.LongName(), " list [options] <path>\n"))
	listCmd.AddNote("如果不指定路径，默认为当前目录")
	listCmd.AddNote("排序选项(-t, -s, -n)不能同时使用, 后指定的选项会覆盖前一个")

	// 添加标志
	listCmdAll = listCmd.Bool("all", "a", false, "列出所有文件和目录，包括隐藏文件和目录")
	listCmdColor = listCmd.Bool("color", "c", false, "启用颜色输出")
	listCmdSortByTime = listCmd.Bool("time", "t", false, "按修改时间排序")
	listCmdSortBySize = listCmd.Bool("size", "s", false, "按文件大小排序")
	listCmdSortByName = listCmd.Bool("name", "n", false, "按文件名排序")
	listCmdDirItself = listCmd.Bool("dir-itself", "D", false, "列出目录本身，而不是文件")
	listCmdType = listCmd.Enum("type", "ty", types.FindTypeAll, "指定要查找的类型，支持以下选项: \n"+
		"\t\t\t\t\t[f | file]       - 只查找文件\n"+
		"\t\t\t\t\t[d | dir]        - 只查找目录\n"+
		"\t\t\t\t\t[l | symlink]    - 只查找软链接\n"+
		"\t\t\t\t\t[r | readonly]   - 只查找只读文件\n"+
		"\t\t\t\t\t[h | hidden]     - 只显示隐藏文件或目录", types.ListTypeLimits)
	listCmdLongFormat = listCmd.Bool("long", "l", false, "使用长格式显示文件信息，包括权限、所有者、大小等")
	listCmdReverseSort = listCmd.Bool("reverse", "r", false, "反向排序")
	listCmdQuoteNames = listCmd.Bool("quote-names", "q", false, "在输出时用双引号包裹条目")
	listCmdRecursion = listCmd.Bool("recursion", "R", false, "递归列出目录及其子目录的内容")
	listCmdShowUserGroup = listCmd.Bool("user-group", "u", false, "显示文件的用户和组信息")
	listCmdTableStyle = listCmd.Enum("table-style", "ts", "none", "指定表格样式，支持以下选项：\n"+
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
		"\t\t\t\t\t[none]   - 禁用边框样式", types.TableStyles)

	// 返回子命令
	return listCmd
}
