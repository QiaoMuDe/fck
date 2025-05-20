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
	sizeCmd     = flag.NewFlagSet("size", flag.ExitOnError)
	sizeCmdHelp = sizeCmd.Bool("h", false, "打印帮助信息并退出")

	// fck check 子命令
	checkCmd           = flag.NewFlagSet("check", flag.ExitOnError)
	checkCmdHelp       = checkCmd.Bool("h", false, "打印帮助信息并退出")
	checkCmdFile       = checkCmd.String("f", "", "指定校验值文件, 根据文件中的哈希值进行校验")
	checkCmdSourcePath = checkCmd.String("s", "", "指定要校验的源目录")
	checkCmdTargetPath = checkCmd.String("t", "", "指定要校验的目标目录")

	// fck find 子命令
	findCmd         = flag.NewFlagSet("find", flag.ExitOnError)
	findCmdHelp     = findCmd.Bool("h", false, "打印帮助信息并退出")
	findCmdPath     = findCmd.String("p", "", "指定要查找的路径")
	findCmdKeyword  = findCmd.String("k", "", "指定要查找的关键字或正则表达式")
	findCmdMaxDepth = findCmd.Int("m", -1, "指定查找的最大深度, -1 表示不限制")
	findCmdFile     = findCmd.Bool("f", false, "限制只查找文件")
	findCmdDir      = findCmd.Bool("d", false, "限制只查找目录")
	findCmdSymlink  = findCmd.Bool("l", false, "限制只查找软链接")
	findCmdSize     = findCmd.String("s", "", "按文件大小过滤，格式如+5M(大于5M)或-5M(小于5M), 支持单位B/K/M/G")
	findCmdModTime  = findCmd.String("mtime", "", "按修改时间查找, 格式如+5(大于5天)或-5(小于5天)")
)
