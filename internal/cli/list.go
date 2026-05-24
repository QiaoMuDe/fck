package cli

import (
	"fmt"
	"os"

	"gitee.com/MM-Q/color"
	"gitee.com/MM-Q/fck/internal/commands/list"
	"gitee.com/MM-Q/fck/internal/types"
	"gitee.com/MM-Q/qflag"
	"golang.org/x/term"
)

var ListCmd *qflag.Cmd

var (
	listAll           *qflag.BoolFlag // all 标志
	listColor         *qflag.BoolFlag // color 标志
	listSortByTime    *qflag.BoolFlag // time 标志
	listSortBySize    *qflag.BoolFlag // size 标志
	listSortByName    *qflag.BoolFlag // name 标志
	listDirItself     *qflag.BoolFlag // dir-itself 标志
	listLongFormat    *qflag.BoolFlag // long 标志
	listReverseSort   *qflag.BoolFlag // reverse 标志
	listQuoteNames    *qflag.BoolFlag // quote-names 标志
	listRecursion     *qflag.BoolFlag // recursion 标志
	listShowUserGroup *qflag.BoolFlag // user-group 标志
	listTableStyle    *qflag.EnumFlag // table-style 标志
	listType          *qflag.EnumFlag // type 标志
	listIcon          *qflag.BoolFlag // icon 标志
	listDisableIndex  *qflag.BoolFlag // disable-index 标志
)

func init() {
	ListCmd = qflag.NewCmd("list", "ls", qflag.ExitOnError)

	listAll = ListCmd.Bool("all", "a", "列出所有文件和目录，包括隐藏文件和目录", false)
	listColor = ListCmd.Bool("color", "c", "启用颜色输出", false)
	listSortByTime = ListCmd.Bool("time", "t", "按修改时间排序", false)
	listSortBySize = ListCmd.Bool("size", "s", "按文件大小排序", false)
	listSortByName = ListCmd.Bool("name", "n", "按文件名排序", false)
	listDirItself = ListCmd.Bool("dir-itself", "D", "列出目录本身，而不是文件", false)
	listType = ListCmd.Enum("type", "ty", "指定要查找的类型，支持以下选项: \n"+
		"\t\t\t\t\t[f | file]       - 只查找文件\n"+
		"\t\t\t\t\t[d | dir]        - 只查找目录\n"+
		"\t\t\t\t\t[l | symlink]    - 只查找软链接\n"+
		"\t\t\t\t\t[r | readonly]   - 只查找只读文件\n"+
		"\t\t\t\t\t[h | hidden]     - 只显示隐藏文件或目录", types.FindTypeAll, types.ListTypeLimits)
	listLongFormat = ListCmd.Bool("long", "l", "使用长格式显示文件信息，包括权限、所有者、大小等", false)
	listReverseSort = ListCmd.Bool("reverse", "r", "反向排序", false)
	listQuoteNames = ListCmd.Bool("quote-names", "q", "在输出时用双引号包裹条目", false)
	listRecursion = ListCmd.Bool("recursion", "R", "递归列出目录及其子目录的内容", false)
	listShowUserGroup = ListCmd.Bool("user-group", "u", "显示文件的用户和组信息", false)
	listTableStyle = ListCmd.Enum("table-style", "ts", "指定表格样式，支持以下选项：\n"+
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
		"\t\t\t\t\t[none]   - 禁用边框样式", "none", types.TableStyles)
	listIcon = ListCmd.Bool("icon", "i", "显示文件图标", false)
	listDisableIndex = ListCmd.Bool("disable-index", "di", "禁用索引序号", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "文件目录列表工具",
		Notes:       []string{"如果不指定路径，默认为当前目录", "排序选项(-t, -s, -n)不能同时使用, 后指定的选项会覆盖前一个"},
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s list [options] <path...>", qflag.Root.Name()),
		Examples: map[string]string{
			"显示当前目录": fmt.Sprintf("%s list", qflag.Root.Name()),
			"显示指定目录": fmt.Sprintf("%s list /path/to/directory", qflag.Root.Name()),
			"查看长格式":  fmt.Sprintf("%s list -l", qflag.Root.Name()),
		},
	}

	if err := ListCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	ListCmd.SetRun(runList)
}

// runList 执行 list 命令
//
// 参数:
//   - cmd: qflag 命令对象
//
// 返回:
//   - error: 执行错误 (如果有)
func runList(cmd qflag.Command) error {
	cl := color.G()

	// 检测是否在管道中（标准输出不是终端）
	// 如果在管道中，禁用图标显示，避免图标字符干扰管道数据处理
	isPipe := !term.IsTerminal(int(os.Stdout.Fd()))

	config := list.ListConfig{
		Args:          cmd.Args(),
		All:           listAll.Get(),
		Color:         listColor.Get(),
		SortByTime:    listSortByTime.Get(),
		SortBySize:    listSortBySize.Get(),
		SortByName:    listSortByName.Get(),
		DirItself:     listDirItself.Get(),
		LongFormat:    listLongFormat.Get(),
		ReverseSort:   listReverseSort.Get(),
		QuoteNames:    listQuoteNames.Get(),
		Recursion:     listRecursion.Get(),
		ShowUserGroup: listShowUserGroup.Get(),
		TableStyle:    listTableStyle.Get(),
		Type:          listType.Get(),
		Icon:          listIcon.Get() && !isPipe, // 管道环境下强制禁用图标，避免图标字符干扰后续处理
		DisableIndex:  listDisableIndex.Get(),
	}

	return list.ListCmdMain(cl, config)
}
