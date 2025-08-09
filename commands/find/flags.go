package find

import (
	"flag"
	"fmt"

	"gitee.com/MM-Q/fck/commands/internal/types"
	"gitee.com/MM-Q/qflag"
	"gitee.com/MM-Q/qflag/cmd"
)

var (
	// fck find 子命令
	findCmd              *cmd.Cmd
	findCmdName          *qflag.StringFlag // name 标志
	findCmdPath          *qflag.StringFlag // path 标志
	findCmdExt           *qflag.SliceFlag  // ext 标志
	findCmdMaxDepth      *qflag.IntFlag    // max-depth 标志
	findCmdSize          *qflag.StringFlag // size 标志
	findCmdModTime       *qflag.StringFlag // mod-time 标志
	findCmdCase          *qflag.BoolFlag   // case 标志
	findCmdFullPath      *qflag.BoolFlag   // full-path 标志
	findCmdHidden        *qflag.BoolFlag   // hidden 标志
	findCmdColor         *qflag.BoolFlag   // color 标志
	findCmdRegex         *qflag.BoolFlag   // regex 标志
	findCmdExcludeName   *qflag.StringFlag // exclude-name 标志
	findCmdExcludePath   *qflag.StringFlag // exclude-path 标志
	findCmdExec          *qflag.StringFlag // exec 标志
	findCmdPrintCmd      *qflag.BoolFlag   // print-cmd 标志
	findCmdDelete        *qflag.BoolFlag   // delete 标志
	findCmdPrintDelete   *qflag.BoolFlag   // print-delete 标志
	findCmdMove          *qflag.StringFlag // move 标志
	findCmdPrintMove     *qflag.BoolFlag   // print-move 标志
	findCmdAnd           *qflag.BoolFlag   // and 标志
	findCmdOr            *qflag.BoolFlag   // or 标志
	findCmdMaxDepthLimit *qflag.IntFlag    // max-depth-limit 标志
	findCmdCount         *qflag.BoolFlag   // count 标志
	findCmdX             *qflag.BoolFlag   // x 标志
	findCmdType          *qflag.EnumFlag   // type 标志
	findCmdWholeWord     *qflag.BoolFlag   // whole-word 标志
)

func InitFindCmd() *cmd.Cmd {
	// fck find 子命令
	findCmd = qflag.NewCmd("find", "f", flag.ExitOnError).
		WithUsageSyntax(fmt.Sprint(qflag.LongName(), " find [options] <path>\n")).
		WithUseChinese(true).
		WithDescription("文件目录查找工具, 在指定目录及其子目录中按照多种条件查找文件和目录")
	findCmd.AddNote("大小单位支持B/K/M/G/b/k/m/g")
	findCmd.AddNote("时间参数以天为单位")
	findCmd.AddNote("不能同时指定-f、-d和-l标志")
	findCmd.AddNote("不能同时执行-exec和-delete标志")
	findCmd.AddNote("如果不指定路径，默认为当前目录")
	findCmdName = findCmd.String("name", "n", "", "指定要查找的文件或目录名")
	findCmdPath = findCmd.String("path", "p", "", "指定要查找的路径")
	findCmdExt = findCmd.Slice("ext", "e", []string{}, "按文件扩展名查找(支持多个扩展名，如 '.txt,.go', '.txt|.go', '.txt;.go')")
	findCmdMaxDepth = findCmd.Int("max-depth", "m", -1, "指定查找的最大深度, -1 表示不限制")
	findCmdSize = findCmd.String("size", "s", "", "按文件大小过滤, 格式如+5M(大于5M)或-5M(小于5M), 支持单位B/K/M/G")
	findCmdModTime = findCmd.String("mtime", "mt", "", "按修改时间过滤, 默认格式如+5(5天前)或-5(5天内)")
	findCmdCase = findCmd.Bool("case", "C", false, "启用大小写敏感匹配, 默认不区分大小写")
	findCmdFullPath = findCmd.Bool("full-path", "F", false, "是否显示完整路径, 默认显示匹配到的路径")
	findCmdHidden = findCmd.Bool("hidden", "H", false, "显示隐藏文件和目录，默认过滤隐藏项")
	findCmdColor = findCmd.Bool("color", "c", false, "启用颜色输出")
	findCmdRegex = findCmd.Bool("regex", "R", false, "启用正则表达式匹配, 默认不启用")
	findCmdExcludeName = findCmd.String("exclude-name", "en", "", "指定要排除的文件或目录名")
	findCmdExcludePath = findCmd.String("exclude-path", "ep", "", "指定要排除的路径")
	findCmdExec = findCmd.String("exec", "ex", "", "对匹配的每个路径执行指定命令，使用{}作为占位符")
	findCmdPrintCmd = findCmd.Bool("print-cmd", "pc", false, "在执行-exec命令前打印将要执行的命令")
	findCmdDelete = findCmd.Bool("delete", "d", false, "删除匹配的文件或目录")
	findCmdPrintDelete = findCmd.Bool("print-del", "pd", false, "在删除前打印将要删除的文件或目录")
	findCmdMove = findCmd.String("move", "mv", "", "将匹配项移动到指定的路径")
	findCmdPrintMove = findCmd.Bool("print-mv", "pm", false, "在移动前打印 old -> new")
	findCmdAnd = findCmd.Bool("and", "", true, "用于在-n和-p参数中组合条件, 默认为true, 表示所有条件必须满足")
	findCmdOr = findCmd.Bool("or", "", false, "用于在-n和-p参数中组合条件, 默认为false, 表示只要满足任一条件即可")
	findCmdMaxDepthLimit = findCmd.Int("max-depth-limit", "mdl", 32, "指定软连接最大解析深度, 默认为32, 超过该深度将停止解析")
	findCmdCount = findCmd.Bool("count", "ct", false, "仅统计匹配项的数量而不显示具体路径")
	findCmdX = findCmd.Bool("xmode", "X", false, "启用并发模式")
	findCmdType = findCmd.Enum("type", "t", "all", "指定要查找的类型，支持以下选项：\n"+
		"\t\t\t\t\t[f | file]       - 只查找文件\n"+
		"\t\t\t\t\t[d | dir]        - 只查找目录\n"+
		"\t\t\t\t\t[l | symlink]    - 只查找软链接\n"+
		"\t\t\t\t\t[r | readonly]   - 只查找只读文件\n"+
		"\t\t\t\t\t[h | hidden]     - 只显示隐藏文件或目录\n"+
		"\t\t\t\t\t[e | empty]      - 只查找空文件或目录\n"+
		"\t\t\t\t\t[x | executable] - 只查找可执行文件\n"+
		"\t\t\t\t\t[s | socket]     - 只查找socket文件\n"+
		"\t\t\t\t\t[p | pipe]       - 只查找管道文件\n"+
		"\t\t\t\t\t[b | block]      - 只查找块设备文件\n"+
		"\t\t\t\t\t[c | char]       - 只查找字符设备文件\n"+
		"\t\t\t\t\t[a | append]     - 只查找追加模式文件\n"+
		"\t\t\t\t\t[n | nonappend]  - 只查找非追加模式文件\n"+
		"\t\t\t\t\t[u | exclusive]  - 只查找独占模式文件", types.FindTypeLimits)
	findCmdWholeWord = findCmd.Bool("whole-word", "W", false, "匹配完整关键字")

	return findCmd
}
