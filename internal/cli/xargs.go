package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/xargs"
	"gitee.com/MM-Q/qflag"
)

var XargsCmd *qflag.Cmd

var (
	xargsDelimiter    *qflag.StringFlag // -d, --delimiter  输入分隔符
	xargsNull         *qflag.BoolFlag   // -0, --null       使用 \0 作为分隔符
	xargsArgFile      *qflag.StringFlag // -a, --arg-file   从文件读取参数
	xargsMaxArgs      *qflag.IntFlag    // -n, --max-args   每批最大参数个数
	xargsMaxChars     *qflag.IntFlag    // -s, --max-chars  命令最大长度
	xargsReplaceStr   *qflag.BoolFlag   // -i, --replace    启用占位符替换模式
	xargsReplaceDelim *qflag.StringFlag // -I, --replace-delim  自定义占位符字符串
	xargsMaxProcs     *qflag.IntFlag    // -P, --max-procs  并行进程数
	xargsNoRunIfEmpty *qflag.BoolFlag   // -r, --no-run-if-empty  空输入不执行
	xargsVerbose      *qflag.BoolFlag   // -t, --verbose    打印执行的命令
	xargsExitOnError  *qflag.BoolFlag   // -e, --exit-on-error  出错立即停止
	xargsShell        *qflag.BoolFlag   // --shell 通过 shell 执行命令
)

func init() {
	XargsCmd = qflag.NewCmd("xargs", "", qflag.ExitOnError)

	xargsDelimiter = XargsCmd.String("delimiter", "d", "输入分隔符 (默认空格/换行/制表)", "")
	xargsNull = XargsCmd.Bool("null", "0", "使用 \\0 作为分隔符", false)
	xargsArgFile = XargsCmd.String("arg-file", "a", "从文件读取参数", "")
	xargsMaxArgs = XargsCmd.Int("max-args", "n", "每批最大参数个数", 0)
	xargsMaxChars = XargsCmd.Int("max-chars", "s", "命令最大长度", 0)
	xargsReplaceStr = XargsCmd.Bool("replace", "i", "启用占位符替换模式 (默认占位符 {})", false)
	xargsReplaceDelim = XargsCmd.String("replace-delim", "I", "自定义占位符字符串", "")
	xargsMaxProcs = XargsCmd.Int("max-procs", "P", "并行进程数", 1)
	xargsNoRunIfEmpty = XargsCmd.Bool("no-run-if-empty", "r", "空输入不执行", false)
	xargsVerbose = XargsCmd.Bool("verbose", "t", "打印执行的命令", false)
	xargsExitOnError = XargsCmd.Bool("exit-on-error", "e", "出错立即停止", false)
	xargsShell = XargsCmd.Bool("shell", "S", "通过 shell 执行命令（支持管道、重定向等）", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:       "从标准输入或文件读取参数，批量执行指定命令",
		UseChinese: true,
		Examples: map[string]string{
			"基础用法":   fmt.Sprintf("echo a b c | %s xargs echo", qflag.Root.Name()),
			"参数替换":   fmt.Sprintf("find . -name \"*.log\" | %s xargs -i mv {} {}.bak", qflag.Root.Name()),
			"自定义占位符": fmt.Sprintf("echo a | %s xargs -i -I %% echo %%", qflag.Root.Name()),
			"批量处理":   fmt.Sprintf("echo a b c d | %s xargs -n 2 echo", qflag.Root.Name()),
			"并行执行":   fmt.Sprintf("cat urls.txt | %s xargs -P 4 curl -O", qflag.Root.Name()),
			"从文件读取":  fmt.Sprintf("%s xargs -a files.txt rm", qflag.Root.Name()),
			"空分隔符":   fmt.Sprintf("find . -print0 | %s xargs -0 rm", qflag.Root.Name()),
			"限制命令长度": fmt.Sprintf("%s xargs -s 1024 echo", qflag.Root.Name()),
			"组合使用":   fmt.Sprintf("find . -name \"*.tmp\" | %s xargs -0 -i -P 4 mv {} /backup/", qflag.Root.Name()),
			"使用管道":   fmt.Sprintf("echo \"file.txt\" | %s xargs -i --shell \"cat {} | grep pattern\"", qflag.Root.Name()),
			"使用重定向":  fmt.Sprintf("echo \"file.txt\" | %s xargs -i --shell \"cat {} > output.txt\"", qflag.Root.Name()),
		},
		Notes: []string{
			"默认模式：所有参数追加到命令后，执行一次",
			"-i 模式：每个参数单独替换占位符，执行多次",
			"使用 -0 可以处理包含空格或特殊字符的文件名",
			"使用 -P 可以并行执行，提高处理速度",
			"默认使用直接执行模式，避免 shell 注入风险",
			"如需使用管道、重定向等 shell 特性，请添加 --shell/-S 标志",
		},
		MutexGroups: []qflag.MutexGroup{
			{
				Name: "replace-max-args",
				Flags: []string{
					"replace",
					"max-args",
				},
				AllowNone: true,
			},
			{
				Name: "replace-delim-max-args",
				Flags: []string{
					"replace-delim",
					"max-args",
				},
				AllowNone: true,
			},
		},
	}

	if err := XargsCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	XargsCmd.SetRun(runXargs)
}

func runXargs(cmd qflag.Command) error {
	args := cmd.Args()
	if len(args) == 0 {
		return fmt.Errorf("missing command argument")
	}

	// 解析占位符
	replaceStr := ""
	if xargsReplaceStr.Get() {
		replaceStr = "{}"
	}
	replaceDelim := xargsReplaceDelim.Get()
	if replaceDelim != "" {
		replaceStr = replaceDelim
	}

	config := xargs.XargsConfig{
		Delimiter:    xargsDelimiter.Get(),
		NullMode:     xargsNull.Get(),
		ArgFile:      xargsArgFile.Get(),
		MaxArgs:      xargsMaxArgs.Get(),
		MaxChars:     xargsMaxChars.Get(),
		ReplaceStr:   replaceStr,
		ReplaceDelim: replaceDelim,
		MaxProcs:     xargsMaxProcs.Get(),
		NoRunIfEmpty: xargsNoRunIfEmpty.Get(),
		Verbose:      xargsVerbose.Get(),
		ExitOnError:  xargsExitOnError.Get(),
		Shell:        xargsShell.Get(),
		Command:      args[0],  // 第一个参数作为命令
		CommandArgs:  args[1:], // 其他参数作为命令参数
	}

	return xargs.XargsCmdMain(config)
}
