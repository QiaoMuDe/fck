package cmd

import (
	"flag"
	"fmt"
	"os"

	"gitee.com/MM-Q/fck/globals"
	"gitee.com/MM-Q/qflag"
	"gitee.com/MM-Q/verman"
)

var (
	// fck hash 子命令
	hashCmd          *qflag.Cmd
	hashCmdType      *qflag.EnumFlag // type 标志
	hashCmdRecursion *qflag.BoolFlag // recursion 标志
	hashCmdJob       *qflag.IntFlag  // job 标志
	hashCmdWrite     *qflag.BoolFlag // write 标志
	hashCmdHidden    *qflag.BoolFlag // hidden 标志

	// fck size 子命令
	sizeCmd           *qflag.Cmd
	sizeCmdColor      *qflag.BoolFlag   // color 标志
	sizeCmdJob        *qflag.IntFlag    // job 标志
	sizeCmdTableStyle *qflag.StringFlag // ts 标志
	sizeCmdHidden     *qflag.BoolFlag   // hidden 标志

	// fck diff 子命令
	diffCmd      *qflag.Cmd
	diffCmdFile  *qflag.StringFlag // file 标志
	diffCmdDirs  *qflag.StringFlag // dirs 标志
	diffCmdDirA  *qflag.StringFlag // dirA 标志
	diffCmdDirB  *qflag.StringFlag // dirB 标志
	diffCmdType  *qflag.EnumFlag   // type 标志
	diffCmdWrite *qflag.BoolFlag   // write 标志

	// fck find 子命令
	findCmd              *qflag.Cmd
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
	findCmdType          *qflag.StringFlag // type 标志
	findCmdWholeWord     *qflag.BoolFlag   // whole-word 标志

	// fck list 子命令
	listCmd              *qflag.Cmd
	listCmdAll           *qflag.BoolFlag   // all 标志
	listCmdColor         *qflag.BoolFlag   // color 标志
	listCmdSortByName    *qflag.BoolFlag   // sort-by-name 标志
	listCmdSortBySize    *qflag.BoolFlag   // sort-by-size 标志
	listCmdSortByTime    *qflag.BoolFlag   // sort-by-time 标志
	listCmdDirEctory     *qflag.BoolFlag   // D 标志
	listCmdFileOnly      *qflag.BoolFlag   // f 标志
	listCmdDirOnly       *qflag.BoolFlag   // d 标志
	listCmdSymlink       *qflag.BoolFlag   // l 标志
	listCmdReadOnly      *qflag.BoolFlag   // ro 标志
	listCmdHiddenOnly    *qflag.BoolFlag   // ho 标志
	listCmdLongFormat    *qflag.BoolFlag   // l 标志
	listCmdReverseSort   *qflag.BoolFlag   // r 标志
	listCmdQuoteNames    *qflag.BoolFlag   // q 标志
	listCmdRecursion     *qflag.BoolFlag   // R 标志
	listCmdShowUserGroup *qflag.BoolFlag   // u 标志
	listCmdTableStyle    *qflag.StringFlag // ts 标志
	listCmdDevColor      *qflag.BoolFlag   // dev-color 标志
)

func init() {
	defer func() {
		if err := recover(); err != nil {
			// 打印错误信息并退出
			fmt.Printf("err: %v\n", err)
			os.Exit(1)
		}
	}()

	// 注册版本信息标志
	v := verman.Get()                                               // 获取版本信息
	qflag.SetVersion(fmt.Sprintf("%s %s", v.AppName, v.GitVersion)) // 设置版本信息
	qflag.SetUseChinese(true)                                       // 启用中文帮助信息
	qflag.SetDescription("多功能文件处理工具集, 提供文件哈希计算、大小统计、查找和校验等实用功能")    // 设置命令行描述
	qflag.AddNote("各子命令有独立帮助文档，可通过-h参数查看, 例如 'fck <子命令> -h' 查看各子命令详细帮助")
	qflag.AddNote("所有路径参数支持Windows和Unix风格")
	qflag.SetLogoText(globals.FckHelpLogo) // 设置命令行logo

	// fck hash 子命令
	hashCmd = qflag.NewCmd("hash", "h", flag.ExitOnError)
	hashCmd.SetUsageSyntax(fmt.Sprint(qflag.LongName(), " hash [options] <path>\n\n"))
	hashCmd.SetUseChinese(true)                                     // 启用中文帮助信息
	hashCmd.SetDescription("文件哈希计算工具, 计算指定文件或目录的哈希值，支持多种哈希算法和并发处理") // 设置命令行描述
	hashCmdType = hashCmd.Enum("type", "t", "md5", "指定哈希算法，支持 md5、sha1、sha256、sha512", []string{"md5", "sha1", "sha256", "sha512"})
	hashCmdRecursion = hashCmd.Bool("recursion", "r", false, "递归处理目录")
	hashCmdJob = hashCmd.Int("job", "j", -1, "指定并发数量, 默认为-1表示根据CPU核心数自动设置, 其余整数表示并发任务数")
	hashCmdWrite = hashCmd.Bool("write", "w", false, "将哈希值写入文件, 文件名为checksum.hash")
	hashCmdHidden = hashCmd.Bool("hidden", "H", false, "启用计算隐藏文件/目录的哈希值，默认跳过")

	// fck size 子命令
	sizeCmd = qflag.NewCmd("size", "s", flag.ExitOnError)
	sizeCmd.SetUsageSyntax(fmt.Sprint(qflag.LongName(), " size [options] <path>\n\n"))
	sizeCmd.SetUseChinese(true) // 启用中文帮助信息
	sizeCmd.AddNote("大小单位会自动选择最合适的(B/KB/MB/GB/TB)")
	sizeCmd.SetDescription("文件目录大小计算工具, 计算指定文件或目录的大小，并以人类可读格式(B/KB/MB/GB/TB)显示")
	sizeCmdColor = sizeCmd.Bool("color", "c", false, "启用颜色输出")
	sizeCmdJob = sizeCmd.Int("job", "j", -1, "指定并发数量, 默认为-1表示根据CPU核心数自动设置, 其余整数表示并发任务数")
	sizeCmdHidden = sizeCmd.Bool("hidden", "H", false, "包含隐藏文件或目录进行大小计算，默认过滤")
	sizeCmdTableStyle = sizeCmd.String("table-style", "ts", "none", "指定表格样式，支持以下选项：\n"+
		"\t\t\t\t\t[default] - 默认样式\n"+
		"\t\t\t\t\t[l  ]     - 浅色样式\n"+
		"\t\t\t\t\t[r  ]     - 圆角样式\n"+
		"\t\t\t\t\t[bd ]     - 粗体样式\n"+
		"\t\t\t\t\t[cb ]     - 亮色彩色样式\n"+
		"\t\t\t\t\t[cd ]     - 暗色彩色样式\n"+
		"\t\t\t\t\t[db ]     - 双线样式\n"+
		"\t\t\t\t\t[cbb]     - 黑色背景蓝色字体\n"+
		"\t\t\t\t\t[cbc]     - 青色背景蓝色字体\n"+
		"\t\t\t\t\t[cbg]     - 绿色背景蓝色字体\n"+
		"\t\t\t\t\t[cbm]     - 紫色背景蓝色字体\n"+
		"\t\t\t\t\t[cby]     - 黄色背景蓝色字体\n"+
		"\t\t\t\t\t[cbr]     - 红色背景蓝色字体\n"+
		"\t\t\t\t\t[cwb]     - 蓝色背景白色字体\n"+
		"\t\t\t\t\t[ccw]     - 青色背景白色字体\n"+
		"\t\t\t\t\t[cgw]     - 绿色背景白色字体\n"+
		"\t\t\t\t\t[cmw]     - 紫色背景白色字体\n"+
		"\t\t\t\t\t[crw]     - 红色背景白色字体\n"+
		"\t\t\t\t\t[cyw]     - 黄色背景白色字体\n"+
		"\t\t\t\t\t[none]    - 禁用表格样式")

	// fck diff 子命令
	diffCmd = qflag.NewCmd("diff", "d", flag.ExitOnError)
	diffCmd.SetUseChinese(true) // 启用中文帮助信息
	diffCmd.AddNote("校验文件必须包含有效的头信息")
	diffCmd.AddNote("校验时会自动跳过空行和注释行（以#开头的行）")
	diffCmd.SetDescription("文件校验工具, 对比指定目录A和目录B的文件差异, 并支持指定校验类型")
	diffCmdFile = diffCmd.String("file", "f", globals.OutputFileName, "指定用于校验的哈希值文件，程序将依据该文件中的哈希值进行校验操作")
	diffCmdDirs = diffCmd.String("dir", "d", "", "指定需要根据哈希值文件进行校验的目标目录")
	diffCmdDirA = diffCmd.String("dirA", "a", "", "指定要对比的目录A")
	diffCmdDirB = diffCmd.String("dirB", "b", "", "指定要对比的目录B")
	diffCmdWrite = diffCmd.Bool("write", "w", false, "将校验结果写入文件, 文件名为check_dir.check")
	diffCmdType = diffCmd.Enum("type", "t", "md5", "指定哈希算法，支持 md5、sha1、sha256、sha512", []string{"md5", "sha1", "sha256", "sha512"})

	// fck find 子命令
	findCmd = qflag.NewCmd("find", "f", flag.ExitOnError)
	findCmd.AddNote("大小单位支持B/K/M/G/b/k/m/g")
	findCmd.AddNote("时间参数以天为单位")
	findCmd.AddNote("不能同时指定-f、-d和-l标志")
	findCmd.AddNote("不能同时执行-exec和-delete标志")
	findCmd.AddNote("如果不指定路径，默认为当前目录")
	findCmd.SetUseChinese(true)                                                        // 启用中文帮助信息
	findCmd.SetDescription("文件目录查找工具, 在指定目录及其子目录中按照多种条件查找文件和目录")                       // 设置命令行描述
	findCmd.SetUsageSyntax(fmt.Sprint(qflag.LongName(), " find [options] <path>\n\n")) // 设置自定义使用说明
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
	findCmdType = findCmd.String("type", "t", "all", "指定要查找的类型，支持以下选项：\n"+
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
		"\t\t\t\t\t[u | exclusive]  - 只查找独占模式文件")
	findCmdWholeWord = findCmd.Bool("whole-word", "W", false, "匹配完整关键字")

	// fck list 子命令
	listCmd = qflag.NewCmd("list", "ls", flag.ExitOnError)
	listCmd.AddNote("如果不指定路径，默认为当前目录")
	listCmd.AddNote("排序选项(-t, -s, -n)不能同时使用, 后指定的选项会覆盖前一个")
	listCmd.SetUseChinese(true)                                                        // 启用中文帮助信息
	listCmd.SetDescription("文件目录列表工具, 列出指定目录中的文件和目录，并支持多种排序和过滤选项")                     // 设置命令行描述
	listCmd.SetUsageSyntax(fmt.Sprint(qflag.LongName(), " list [options] <path>\n\n")) // 设置自定义使用说明
	listCmdAll = listCmd.Bool("all", "a", false, "列出所有文件和目录，包括隐藏文件和目录")
	listCmdColor = listCmd.Bool("color", "c", false, "启用颜色输出")
	listCmdSortByTime = listCmd.Bool("time", "t", false, "按修改时间排序")
	listCmdSortBySize = listCmd.Bool("size", "s", false, "按文件大小排序")
	listCmdSortByName = listCmd.Bool("name", "n", false, "按文件名排序")
	listCmdDirEctory = listCmd.Bool("directory", "D", false, "列出目录本身，而不是文件")
	listCmdFileOnly = listCmd.Bool("file", "f", false, "只列出文件，不列出目录")
	listCmdDirOnly = listCmd.Bool("dir", "d", false, "只列出目录，不列出文件")
	listCmdSymlink = listCmd.Bool("symlink", "L", false, "只列出软链接，不列出其他类型的文件")
	listCmdReadOnly = listCmd.Bool("readonly", "ro", false, "只列出只读文件")
	listCmdHiddenOnly = listCmd.Bool("hidden", "ho", false, "只列出隐藏文件或目录")
	listCmdLongFormat = listCmd.Bool("long", "l", false, "使用长格式显示文件信息，包括权限、所有者、大小等")
	listCmdReverseSort = listCmd.Bool("reverse", "r", false, "反向排序")
	listCmdQuoteNames = listCmd.Bool("quote-names", "q", false, "在输出时用双引号包裹条目")
	listCmdRecursion = listCmd.Bool("recursion", "R", false, "递归列出目录及其子目录的内容")
	listCmdShowUserGroup = listCmd.Bool("user-group", "u", false, "显示文件的用户和组信息")
	listCmdTableStyle = listCmd.String("table-style", "ts", "none", "指定表格样式，支持以下选项：\n"+
		"\t\t\t\t\t[default] - 默认样式\n"+
		"\t\t\t\t\t[l  ]     - 浅色样式\n"+
		"\t\t\t\t\t[r  ]     - 圆角样式\n"+
		"\t\t\t\t\t[bd ]     - 粗体样式\n"+
		"\t\t\t\t\t[cb ]     - 亮色彩色样式\n"+
		"\t\t\t\t\t[cd ]     - 暗色彩色样式\n"+
		"\t\t\t\t\t[db ]     - 双线样式\n"+
		"\t\t\t\t\t[cbb]     - 黑色背景蓝色字体\n"+
		"\t\t\t\t\t[cbc]     - 青色背景蓝色字体\n"+
		"\t\t\t\t\t[cbg]     - 绿色背景蓝色字体\n"+
		"\t\t\t\t\t[cbm]     - 紫色背景蓝色字体\n"+
		"\t\t\t\t\t[cby]     - 黄色背景蓝色字体\n"+
		"\t\t\t\t\t[cbr]     - 红色背景蓝色字体\n"+
		"\t\t\t\t\t[cwb]     - 蓝色背景白色字体\n"+
		"\t\t\t\t\t[ccw]     - 青色背景白色字体\n"+
		"\t\t\t\t\t[cgw]     - 绿色背景白色字体\n"+
		"\t\t\t\t\t[cmw]     - 紫色背景白色字体\n"+
		"\t\t\t\t\t[crw]     - 红色背景白色字体\n"+
		"\t\t\t\t\t[cyw]     - 黄色背景白色字体\n"+
		"\t\t\t\t\t[none]    - 禁用边框样式")
	listCmdDevColor = listCmd.Bool("dev-color", "dc", false, "启用开发环境下的颜色输出。注意：此选项需配合颜色输出选项 -c 一同使用")

	// 添加子命令
	if addErr := qflag.AddSubCmd(hashCmd, sizeCmd, diffCmd, findCmd, listCmd); addErr != nil {
		fmt.Printf("err: %v\n", addErr)
		os.Exit(1)
	}

	// 解析全局参数
	if err := qflag.Parse(); err != nil {
		fmt.Printf("err: %v\n", err)
		os.Exit(1)
	}
}
