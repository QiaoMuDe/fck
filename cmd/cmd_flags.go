package cmd

import "flag"

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

	// fck size 子命令
	sizeCmd      = flag.NewFlagSet("size", flag.ExitOnError)
	sizeCmdHelp  = sizeCmd.Bool("h", false, "打印帮助信息并退出")
	sizeCmdColor = sizeCmd.Bool("c", false, "启用颜色输出")

	// fck check 子命令
	checkCmd      = flag.NewFlagSet("check", flag.ExitOnError)
	checkCmdHelp  = checkCmd.Bool("h", false, "打印帮助信息并退出")
	checkCmdFile  = checkCmd.String("f", "", "指定用于校验的哈希值文件，程序将依据该文件中的哈希值进行校验操作")
	checkCmdDirs  = checkCmd.String("d", "", "指定需要根据哈希值文件进行校验的目标目录")
	checkCmdDirA  = checkCmd.String("a", "", "指定要校验的目录A")
	checkCmdDirB  = checkCmd.String("b", "", "指定要校验的目录B")
	checkCmdType  = checkCmd.String("t", "md5", "指定哈希算法，支持 md5、sha1、sha256、sha512")
	checkCmdWrite = checkCmd.Bool("w", false, "将校验结果写入文件, 文件名为check_dir.check")

	// fck find 子命令
	findCmd           = flag.NewFlagSet("find", flag.ExitOnError)
	findCmdHelp       = findCmd.Bool("h", false, "打印帮助信息并退出")
	findCmdPath       = findCmd.String("p", "", "指定要查找的路径")
	findCmdKeyword    = findCmd.String("k", "", "指定要查找的关键字或正则表达式")
	findCmdMaxDepth   = findCmd.Int("m", -1, "指定查找的最大深度, -1 表示不限制")
	findCmdFile       = findCmd.Bool("f", false, "限制只查找文件")
	findCmdDir        = findCmd.Bool("d", false, "限制只查找目录")
	findCmdSymlink    = findCmd.Bool("l", false, "限制只查找软链接")
	findCmdReadOnly   = findCmd.Bool("ro", false, "限制只查找只读文件")
	findCmdHiddenOnly = findCmd.Bool("ho", false, "限制只显示隐藏文件或目录")
	findCmdSize       = findCmd.String("size", "", "按文件大小过滤, 格式如+5M(大于5M)或-5M(小于5M), 支持单位B/K/M/G")
	findCmdModTime    = findCmd.String("mtime", "", "按修改时间查找, 格式如+5(5天前)或-5(5天内)")
	findCmdCase       = findCmd.Bool("C", false, "启用大小写敏感匹配, 默认不区分大小写")
	findCmdFullPath   = findCmd.Bool("full", false, "是否显示完整路径, 默认显示匹配到的路径")
	findCmdHidden     = findCmd.Bool("hidden", false, "是否显示隐藏文件, 默认不显示隐藏文件")
	findCmdColor      = findCmd.Bool("c", false, "启用颜色输出")
	findCmdRegex      = findCmd.Bool("regex", false, "启用正则表达式匹配, 默认不启用")
)
