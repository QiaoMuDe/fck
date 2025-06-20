package cmd

import (
	"flag"

	"gitee.com/MM-Q/fck/globals"
)

var (
	// versionF 版本信息
	versionF = flag.Bool("v", false, "打印版本信息并退出")
	// helpF 帮助信息
	helpF = flag.Bool("h", false, "打印帮助信息并退出")

	// fck hash 子命令
	hashCmd          = flag.NewFlagSet("hash", flag.ExitOnError)
	hashCmdHelp      = hashCmd.Bool("h", false, "打印帮助信息并退出")
	hashCmdType      = hashCmd.String("t", "md5", "指定哈希算法，支持 md5、sha1、sha256、sha512")
	hashCmdRecursion = hashCmd.Bool("r", false, "递归处理目录")
	hashCmdJob       = hashCmd.Int("j", -1, "指定并发数量, 默认为-1表示根据CPU核心数自动设置, 其余整数表示并发任务数")
	hashCmdWrite     = hashCmd.Bool("w", false, "将哈希值写入文件, 文件名为checksum.hash")
	hashCmdHidden    = hashCmd.Bool("H", false, "启用计算隐藏文件/目录的哈希值，默认跳过")

	// fck size 子命令
	sizeCmd           = flag.NewFlagSet("size", flag.ExitOnError)
	sizeCmdHelp       = sizeCmd.Bool("h", false, "打印帮助信息并退出")
	sizeCmdColor      = sizeCmd.Bool("c", false, "启用颜色输出")
	sizeCmdJob        = sizeCmd.Int("j", -1, "指定并发数量, 默认为-1表示根据CPU核心数自动设置, 其余整数表示并发任务数")
	sizeCmdTableStyle = sizeCmd.String("ts", "", "指定表格样式，支持以下选项：\n"+
		"  default - 默认样式\n"+
		"  l      - 浅色样式\n"+
		"  r      - 圆角样式\n"+
		"  bd     - 粗体样式\n"+
		"  cb     - 亮色彩色样式\n"+
		"  cd     - 暗色彩色样式\n"+
		"  db     - 双线样式\n"+
		"  cbb    - 黑色背景蓝色字体\n"+
		"  cbc    - 青色背景蓝色字体\n"+
		"  cbg    - 绿色背景蓝色字体\n"+
		"  cbm    - 紫色背景蓝色字体\n"+
		"  cby    - 黄色背景蓝色字体\n"+
		"  cbr    - 红色背景蓝色字体\n"+
		"  cwb    - 蓝色背景白色字体\n"+
		"  ccw    - 青色背景白色字体\n"+
		"  cgw    - 绿色背景白色字体\n"+
		"  cmw    - 紫色背景白色字体\n"+
		"  crw    - 红色背景白色字体\n"+
		"  cyw    - 黄色背景白色字体\n"+
		"  none   - 禁用表格样式")
	sizeCmdHidden = sizeCmd.Bool("H", false, "包含隐藏文件或目录进行大小计算，默认过滤")

	// fck diff 子命令
	diffCmd      = flag.NewFlagSet("diff", flag.ExitOnError)
	diffCmdHelp  = diffCmd.Bool("h", false, "打印帮助信息并退出")
	diffCmdFile  = diffCmd.String("f", globals.OutputFileName, "指定用于校验的哈希值文件，程序将依据该文件中的哈希值进行校验操作")
	diffCmdDirs  = diffCmd.String("d", "", "指定需要根据哈希值文件进行校验的目标目录")
	diffCmdDirA  = diffCmd.String("a", "", "指定要对比的目录A")
	diffCmdDirB  = diffCmd.String("b", "", "指定要对比的目录B")
	diffCmdType  = diffCmd.String("t", "md5", "指定哈希算法，支持 md5、sha1、sha256、sha512")
	diffCmdWrite = diffCmd.Bool("w", false, "将校验结果写入文件, 文件名为check_dir.check")

	// fck find 子命令
	findCmd              = flag.NewFlagSet("find", flag.ExitOnError)
	findCmdHelp          = findCmd.Bool("h", false, "打印帮助信息并退出")
	findCmdName          = findCmd.String("n", "", "指定要查找的文件或目录名")
	findCmdPath          = findCmd.String("p", "", "指定要查找的路径")
	findCmdExt           = findCmd.String("ext", "", "按文件扩展名查找(支持多个扩展名，如 .txt,.go)")
	findCmdMaxDepth      = findCmd.Int("m", -1, "指定查找的最大深度, -1 表示不限制")
	findCmdSize          = findCmd.String("size", "", "按文件大小过滤, 格式如+5M(大于5M)或-5M(小于5M), 支持单位B/K/M/G")
	findCmdModTime       = findCmd.String("mtime", "", "按修改时间查找, 格式如+5(5天前)或-5(5天内)")
	findCmdCase          = findCmd.Bool("C", false, "启用大小写敏感匹配, 默认不区分大小写")
	findCmdFullPath      = findCmd.Bool("F", false, "是否显示完整路径, 默认显示匹配到的路径")
	findCmdHidden        = findCmd.Bool("H", false, "显示隐藏文件和目录，默认过滤隐藏项")
	findCmdColor         = findCmd.Bool("c", false, "启用颜色输出")
	findCmdRegex         = findCmd.Bool("R", false, "启用正则表达式匹配, 默认不启用")
	findCmdExcludeName   = findCmd.String("en", "", "指定要排除的文件或目录名")
	findCmdExcludePath   = findCmd.String("ep", "", "指定要排除的路径")
	findCmdExec          = findCmd.String("exec", "", "对匹配的每个路径执行指定命令，使用{}作为占位符")
	findCmdPrintCmd      = findCmd.Bool("print-cmd", false, "在执行-exec命令前打印将要执行的命令")
	findCmdDelete        = findCmd.Bool("delete", false, "删除匹配的文件或目录")
	findCmdPrintDelete   = findCmd.Bool("print-del", false, "在删除前打印将要删除的文件或目录")
	findCmdMove          = findCmd.String("mv", "", "将匹配项移动到指定的路径")
	findCmdPrintMove     = findCmd.Bool("print-mv", false, "在移动前打印 old -> new 的映射")
	findCmdAnd           = findCmd.Bool("and", true, "用于在-n和-p参数中组合条件, 默认为true, 表示所有条件必须满足")
	findCmdOr            = findCmd.Bool("or", false, "用于在-n和-p参数中组合条件, 默认为false, 表示只要满足任一条件即可")
	findCmdMaxDepthLimit = findCmd.Int("max-depth", 32, "指定软连接最大解析深度, 默认为32, 超过该深度将停止解析")
	findCmdCount         = findCmd.Bool("count", false, "仅统计匹配项的数量而不显示具体路径")
	findCmdX             = findCmd.Bool("X", false, "启用并发模式")
	findCmdType          = findCmd.String("type", "all", "指定要查找的类型，支持以下选项：\n"+
		" [f|file]       - 只查找文件\n"+
		" [d|dir]        - 只查找目录\n"+
		" [l|symlink]    - 只查找软链接\n"+
		" [r|readonly]   - 只查找只读文件\n"+
		" [h|hidden]     - 只显示隐藏文件或目录\n"+
		" [e|empty]      - 只查找空文件或目录\n"+
		" [x|executable] - 只查找可执行文件\n"+
		" [s|socket]     - 只查找socket文件\n"+
		" [p|pipe]       - 只查找管道文件\n"+
		" [b|block]      - 只查找块设备文件\n"+
		" [c|char]       - 只查找字符设备文件\n"+
		" [a|append]     - 只查找追加模式文件\n"+
		" [n|nonappend]  - 只查找非追加模式文件\n"+
		" [u|exclusive]  - 只查找独占模式文件")
	findCmdWholeWord = findCmd.Bool("W", false, "匹配完整关键字")

	// fck list 子命令
	listCmd              = flag.NewFlagSet("list", flag.ExitOnError)
	listCmdHelp          = listCmd.Bool("h", false, "打印帮助信息并退出")
	listCmdAll           = listCmd.Bool("a", false, "列出所有文件和目录，包括隐藏项")
	listCmdColor         = listCmd.Bool("c", false, "启用颜色输出")
	listCmdSortByTime    = listCmd.Bool("t", false, "按修改时间排序")
	listCmdSortBySize    = listCmd.Bool("s", false, "按文件大小排序")
	listCmdSortByName    = listCmd.Bool("n", false, "按文件名排序")
	listCmdDirEctory     = listCmd.Bool("D", false, "列出目录本身，而不是其内容")
	listCmdFileOnly      = listCmd.Bool("f", false, "只列出文件，不列出目录")
	listCmdDirOnly       = listCmd.Bool("d", false, "只列出目录，不列出文件")
	listCmdSymlink       = listCmd.Bool("L", false, "只列出软链接，不列出其他类型的文件")
	listCmdReadOnly      = listCmd.Bool("ro", false, "只列出只读文件")
	listCmdHiddenOnly    = listCmd.Bool("ho", false, "只列出隐藏文件或目录")
	listCmdLongFormat    = listCmd.Bool("l", false, "使用长格式显示文件信息，包括权限、所有者、大小等")
	listCmdReverseSort   = listCmd.Bool("r", false, "反向排序")
	listCmdQuoteNames    = listCmd.Bool("q", false, "在输出时用双引号包裹条目")
	listCmdRecursion     = listCmd.Bool("R", false, "递归列出目录及其子目录的内容")
	listCmdShowUserGroup = listCmd.Bool("u", false, "显示文件的用户和组信息")
	listCmdTableStyle    = listCmd.String("ts", "none", "指定表格样式，支持以下选项：\n"+
		"  default - 默认样式\n"+
		"  l      - 浅色样式\n"+
		"  r      - 圆角样式\n"+
		"  bd     - 粗体样式\n"+
		"  cb     - 亮色彩色样式\n"+
		"  cd     - 暗色彩色样式\n"+
		"  db     - 双线样式\n"+
		"  cbb    - 黑色背景蓝色字体\n"+
		"  cbc    - 青色背景蓝色字体\n"+
		"  cbg    - 绿色背景蓝色字体\n"+
		"  cbm    - 紫色背景蓝色字体\n"+
		"  cby    - 黄色背景蓝色字体\n"+
		"  cbr    - 红色背景蓝色字体\n"+
		"  cwb    - 蓝色背景白色字体\n"+
		"  ccw    - 青色背景白色字体\n"+
		"  cgw    - 绿色背景白色字体\n"+
		"  cmw    - 紫色背景白色字体\n"+
		"  crw    - 红色背景白色字体\n"+
		"  cyw    - 黄色背景白色字体\n"+
		"  none   - 禁用边框样式")
	listCmdDevColor      = listCmd.Bool("dev-color", false, "启用开发环境下的颜色输出。注意：此选项需配合颜色输出选项 -c 一同使用")
	listCmdDevColorShort = listCmd.Bool("dc", false, "启用开发环境下的颜色输出。注意：此选项需配合颜色输出选项 -c 一同使用")
)
