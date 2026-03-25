package cli

import (
	"fmt"
	"os"

	"gitee.com/MM-Q/fck/internal/commands/xargs"
	"gitee.com/MM-Q/qflag"
)

// ============================================
// 1. 全局命令变量
// ============================================
var (
	XargsCmd *qflag.Cmd
)

// ============================================
// 2. 全局标志变量
// ============================================
var (
	xargsMaxArgs    *qflag.IntFlag
	xargsMaxProcs   *qflag.IntFlag
	xargsReplace    *qflag.StringFlag
	xargsNoRunEmpty *qflag.BoolFlag
	xargsMaxLines   *qflag.IntFlag
	xargsVerbose    *qflag.BoolFlag
	xargsDelimiter  *qflag.StringFlag
	xargsNull       *qflag.BoolFlag
	xargsEof        *qflag.StringFlag
)

// ============================================
// 3. init() 初始化
// ============================================
func init() {
	XargsCmd = qflag.NewCmd("xargs", "x", qflag.ExitOnError)

	// 定义标志
	xargsMaxArgs = XargsCmd.Int("max-args", "n", "每次执行传递的最大参数数量", 1)
	xargsMaxProcs = XargsCmd.Int("max-procs", "P", "最大并行进程数", 1)
	xargsReplace = XargsCmd.String("replace", "I", "占位符（用于替换命令中的参数位置）", "{}")
	xargsNoRunEmpty = XargsCmd.Bool("no-run-if-empty", "r", "如果输入为空，不执行命令", false)
	xargsMaxLines = XargsCmd.Int("max-lines", "L", "最大执行次数（0表示无限制）", 0)
	xargsVerbose = XargsCmd.Bool("verbose", "t", "打印要执行的命令", false)
	xargsDelimiter = XargsCmd.String("delimiter", "d", "输入分隔符", "\n")
	xargsNull = XargsCmd.Bool("null", "0", "使用空字符（\\0）作为分隔符", false)
	xargsEof = XargsCmd.String("eof", "e", "指定结束标志，遇到此行停止读取", "")

	// 应用命令配置
	cmdOpts := &qflag.CmdOpts{
		Desc:        "从标准输入读取数据并执行命令",
		UsageSyntax: fmt.Sprintf("%s xargs [选项] 命令 [初始参数]", qflag.Root.Name()),
		UseChinese:  true,
		Notes: []string{
			"从标准输入读取数据，将每行作为参数传递给指定命令",
			"如果不指定命令，默认使用 echo",
			"默认使用 {} 作为占位符，可通过 -I 选项自定义",
		},
		Examples: map[string]string{
			"基本用法":  fmt.Sprintf(`echo "file.txt" | %s xargs cat`, qflag.Root.Name()),
			"批量处理":  fmt.Sprintf(`ls *.txt | %s xargs -n 2 cat`, qflag.Root.Name()),
			"并行执行":  fmt.Sprintf(`cat urls.txt | %s xargs -P 4 curl -O`, qflag.Root.Name()),
			"使用占位符": fmt.Sprintf(`echo "file.txt" | %s xargs -I {} cp {} {}.bak`, qflag.Root.Name()),
		},
	}

	if err := XargsCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	XargsCmd.SetRun(runXargs)
}

// ============================================
// 4. run() 函数
// ============================================
func runXargs(cmd qflag.Command) error {
	// 获取标志值
	config := xargs.XargsConfig{
		MaxArgs:      xargsMaxArgs.Get(),
		MaxProcs:     xargsMaxProcs.Get(),
		Replace:      xargsReplace.Get(),
		NoRunIfEmpty: xargsNoRunEmpty.Get(),
		MaxLines:     xargsMaxLines.Get(),
		Verbose:      xargsVerbose.Get(),
		Delimiter:    xargsDelimiter.Get(),
		UseNull:      xargsNull.Get(),
		Eof:          xargsEof.Get(),
	}

	// 获取要执行的命令和初始参数
	args := cmd.Args()
	if len(args) == 0 {
		config.TargetCmd = "echo" // 默认命令
	} else {
		config.TargetCmd = args[0]
		config.InitialArgs = args[1:]
	}

	// 调用业务逻辑
	return xargs.XargsCmdMain(config, os.Stdin)
}
