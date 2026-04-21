package cli

import (
	"fmt"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/internal/commands/find"
	"gitee.com/MM-Q/fck/internal/types"
	"gitee.com/MM-Q/qflag"
)

var FindCmd *qflag.Cmd

var (
	findName          *qflag.StringFlag      // 文件名模式
	findPath          *qflag.StringFlag      // 路径模式
	findExt           *qflag.StringSliceFlag // 扩展名切片
	findMaxDepth      *qflag.IntFlag         // 最大深度
	findSize          *qflag.StringFlag      // 文件大小
	findModTime       *qflag.StringFlag      // 修改时间模式
	findCase          *qflag.BoolFlag        // 区分大小写
	findFullPath      *qflag.BoolFlag        // 完整路径
	findHidden        *qflag.BoolFlag        // 隐藏文件
	findColor         *qflag.BoolFlag        // 颜色输出
	findRegex         *qflag.BoolFlag        // 正则表达式
	findExcludeName   *qflag.StringFlag      // 排除文件名
	findExcludePath   *qflag.StringFlag      // 排除路径
	findExec          *qflag.StringFlag      // 执行命令
	findDelete        *qflag.BoolFlag        // 删除
	findMove          *qflag.StringFlag      // 移动路径
	findPrintActions  *qflag.BoolFlag        // 打印操作
	findMaxDepthLimit *qflag.IntFlag         // 最大深度限制
	findCount         *qflag.BoolFlag        // 统计数量
	findType          *qflag.EnumFlag        // 类型
	findWholeWord     *qflag.BoolFlag        // 完整关键字
	findQuiet         *qflag.BoolFlag        // 静默模式
)

func init() {
	FindCmd = qflag.NewCmd("find", "", qflag.ExitOnError)

	findName = FindCmd.String("name", "n", "指定要查找的文件或目录名", "")
	findPath = FindCmd.String("path", "p", "指定要查找的路径", "")
	findExt = FindCmd.StringSlice("ext", "e", "按文件扩展名查找", []string{})
	findMaxDepth = FindCmd.Int("max-depth", "m", "指定查找的最大深度, -1 表示不限制", -1)
	findSize = FindCmd.String("size", "s", "按文件大小过滤, 格式如+5M(大于5M)或-5M(小于5M), 支持单位B/K/M/G", "")
	findModTime = FindCmd.String("mtime", "mt", "按修改时间过滤, 默认格式如+5(5天前)或-5(5天内)", "")
	findCase = FindCmd.Bool("case", "C", "启用大小写敏感匹配, 默认不区分大小写", false)
	findFullPath = FindCmd.Bool("full-path", "F", "是否显示完整路径, 默认显示匹配到的路径", false)
	findHidden = FindCmd.Bool("hidden", "H", "显示隐藏文件和目录，默认过滤隐藏项", false)
	findColor = FindCmd.Bool("color", "c", "启用颜色输出", false)
	findRegex = FindCmd.Bool("regex", "R", "启用正则表达式匹配, 默认不启用", false)
	findExcludeName = FindCmd.String("exclude-name", "en", "指定要排除的文件或目录名", "")
	findExcludePath = FindCmd.String("exclude-path", "ep", "指定要排除的路径", "")
	findExec = FindCmd.String("exec", "ex", "对匹配的每个路径执行指定命令，使用{}作为占位符", "")
	findDelete = FindCmd.Bool("delete", "d", "删除匹配的文件或目录", false)
	findMove = FindCmd.String("move", "mv", "将匹配项移动到指定的路径", "")
	findPrintActions = FindCmd.Bool("print-actions", "pa", "打印执行的操作详情(exec/delete/move)", false)
	findMaxDepthLimit = FindCmd.Int("max-depth-limit", "mdl", "指定软连接最大解析深度, 默认为32, 超过该深度将停止解析", 32)
	findCount = FindCmd.Bool("count", "ct", "仅统计匹配项的数量而不显示具体路径", false)
	findType = FindCmd.Enum("type", "t", "指定要查找的类型，支持以下选项：\n"+
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
		"\t\t\t\t\t[u | exclusive]  - 只查找独占模式文件", "all", types.FindTypeLimits)
	findWholeWord = FindCmd.Bool("whole-word", "W", "匹配完整关键字", false)
	findQuiet = FindCmd.Bool("quiet", "q", "静默模式，不显示权限错误和警告信息", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "文件目录查找工具",
		Notes:       []string{"大小单位支持B/K/M/G/b/k/m/g", "时间参数以天为单位", "不能同时执行-exec和-delete以及-move标志", "如果不指定路径，默认为当前目录"},
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s find [options] <path>", qflag.Root.Name()),
		Examples: map[string]string{
			"查找所有文件":             fmt.Sprintf("%s find", qflag.Root.Name()),
			"查找指定路径下的所有文件":       fmt.Sprintf("%s find /opt/project", qflag.Root.Name()),
			"查找所有的扩展名为py和go的文件":  fmt.Sprintf("%s find -e py,go", qflag.Root.Name()),
			"删除当前目录下所有log扩展名的文件": fmt.Sprintf("%s find -e log -d", qflag.Root.Name()),
		},
		MutexGroups: []qflag.MutexGroup{
			{
				Name:      "action",
				Flags:     []string{"exec", "delete", "move", "count"},
				AllowNone: true,
			},
		},
	}

	if err := FindCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	FindCmd.SetRun(runFind)
}

func runFind(cmd qflag.Command) error {
	cl := colorlib.NewColorLib()

	config := find.FindConfig{
		FindPath:       cmd.Arg(0),
		NamePattern:    findName.Get(),
		PathPattern:    findPath.Get(),
		ExtSlice:       findExt.Get(),
		MaxDepth:       findMaxDepth.Get(),
		SizePattern:    findSize.Get(),
		ModTimePattern: findModTime.Get(),
		CaseSensitive:  findCase.Get(),
		FullPath:       findFullPath.Get(),
		Hidden:         findHidden.Get(),
		Color:          findColor.Get(),
		Regex:          findRegex.Get(),
		ExcludeName:    findExcludeName.Get(),
		ExcludePath:    findExcludePath.Get(),
		ExecCmd:        findExec.Get(),
		Delete:         findDelete.Get(),
		MovePath:       findMove.Get(),
		PrintActions:   findPrintActions.Get(),
		MaxDepthLimit:  findMaxDepthLimit.Get(),
		Count:          findCount.Get(),
		Type:           findType.Get(),
		WholeWord:      findWholeWord.Get(),
		Quiet:          findQuiet.Get(),
	}

	return find.FindCmdMain(cl, &config)
}
