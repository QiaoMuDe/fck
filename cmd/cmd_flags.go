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
	hashCmdJob       = hashCmd.Int("j", 1, "指定并发数量")
	hashCmdWrite     = hashCmd.Bool("w", false, "将哈希值写入文件, 文件名为checksum.hash")
	hashCmdHidden    = hashCmd.Bool("H", false, "启用计算隐藏文件/目录的哈希值，默认跳过")

	// fck size 子命令
	sizeCmd      = flag.NewFlagSet("size", flag.ExitOnError)
	sizeCmdHelp  = sizeCmd.Bool("h", false, "打印帮助信息并退出")
	sizeCmdColor = sizeCmd.Bool("c", false, "启用颜色输出")

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
	findCmd            = flag.NewFlagSet("find", flag.ExitOnError)
	findCmdHelp        = findCmd.Bool("h", false, "打印帮助信息并退出")
	findCmdName        = findCmd.String("n", "", "指定要查找的文件或目录名")
	findCmdPath        = findCmd.String("p", "", "指定要查找的路径")
	findCmdMaxDepth    = findCmd.Int("m", -1, "指定查找的最大深度, -1 表示不限制")
	findCmdFile        = findCmd.Bool("f", false, "限制只查找文件")
	findCmdDir         = findCmd.Bool("d", false, "限制只查找目录")
	findCmdSymlink     = findCmd.Bool("l", false, "限制只查找软链接")
	findCmdReadOnly    = findCmd.Bool("ro", false, "限制只查找只读文件")
	findCmdHiddenOnly  = findCmd.Bool("ho", false, "限制只显示隐藏文件或目录")
	findCmdSize        = findCmd.String("size", "", "按文件大小过滤, 格式如+5M(大于5M)或-5M(小于5M), 支持单位B/K/M/G")
	findCmdModTime     = findCmd.String("mtime", "", "按修改时间查找, 格式如+5(5天前)或-5(5天内)")
	findCmdCase        = findCmd.Bool("C", false, "启用大小写敏感匹配, 默认不区分大小写")
	findCmdFullPath    = findCmd.Bool("F", false, "是否显示完整路径, 默认显示匹配到的路径")
	findCmdHidden      = findCmd.Bool("H", false, "显示隐藏文件和目录，默认过滤隐藏项")
	findCmdColor       = findCmd.Bool("c", false, "启用颜色输出")
	findCmdRegex       = findCmd.Bool("R", false, "启用正则表达式匹配, 默认不启用")
	findCmdExcludeName = findCmd.String("en", "", "指定要排除的文件或目录名")
	findCmdExcludePath = findCmd.String("ep", "", "指定要排除的路径")
	findCmdExec        = findCmd.String("exec", "", "对匹配的每个路径执行指定命令，使用{}作为占位符")
	findCmdPrintCmd    = findCmd.Bool("print-cmd", false, "在执行-exec命令前打印将要执行的命令")
	findCmdDelete      = findCmd.Bool("delete", false, "删除匹配的文件或目录")
	findCmdMove        = findCmd.String("mv", "", "将匹配项移动到指定的路径")
	findCmdPrintMove   = findCmd.Bool("print-mv", false, "在移动前打印 old -> new 的映射")
	findCmdAnd         = findCmd.Bool("and", true, "用于在-n和-p参数中组合条件, 默认为true, 表示所有条件必须满足")
	findCmdOr          = findCmd.Bool("or", false, "用于在-n和-p参数中组合条件, 默认为false, 表示只要满足任一条件即可")

	// fck list 子命令
	listCmd            = flag.NewFlagSet("list", flag.ExitOnError)
	listCmdHelp        = listCmd.Bool("h", false, "打印帮助信息并退出")
	listCmdAll         = listCmd.Bool("a", false, "列出所有文件和目录，包括隐藏项")
	listCmdColor       = listCmd.Bool("c", false, "启用颜色输出")
	listCmdSortByTime  = listCmd.Bool("t", false, "按修改时间排序")
	listCmdSortBySize  = listCmd.Bool("s", false, "按文件大小排序")
	listCmdSortByName  = listCmd.Bool("n", false, "按文件名排序")
	listCmdDirEctory   = listCmd.Bool("D", false, "列出目录本身，而不是其内容")
	listCmdFileOnly    = listCmd.Bool("f", false, "只列出文件，不列出目录")
	listCmdDirOnly     = listCmd.Bool("d", false, "只列出目录，不列出文件")
	listCmdSymlink     = listCmd.Bool("L", false, "只列出软链接，不列出其他类型的文件")
	listCmdReadOnly    = listCmd.Bool("ro", false, "只列出只读文件")
	listCmdHiddenOnly  = listCmd.Bool("ho", false, "只列出隐藏文件或目录")
	listCmdLongFormat  = listCmd.Bool("l", false, "使用长格式显示文件信息，包括权限、所有者、大小等")
	listCmdReverseSort = listCmd.Bool("r", false, "反向排序")
	listCmdQuoteNames  = listCmd.Bool("q", false, "在输出时用双引号包裹条目")
	// listCmdRecursion     = listCmd.Bool("R", false, "递归列出目录及其子目录的内容")
	listCmdShowUserGroup = listCmd.Bool("u", false, "显示文件的用户和组信息")
)
