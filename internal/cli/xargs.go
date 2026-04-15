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
	xargsMaxLines     *qflag.IntFlag    // -L, --max-lines  每批最大行数
	xargsMaxChars     *qflag.IntFlag    // -s, --max-chars  命令最大长度
	xargsReplaceStr   *qflag.BoolFlag   // -i, --replace    启用占位符替换模式
	xargsReplaceDelim *qflag.StringFlag // -I, --replace-delim  自定义占位符字符串
	xargsMaxProcs     *qflag.IntFlag    // -P, --max-procs  并行进程数
	xargsNoRunIfEmpty *qflag.BoolFlag   // -r, --no-run-if-empty  空输入不执行
	xargsInteractive  *qflag.BoolFlag   // -p, --interactive 执行前确认
	xargsVerbose      *qflag.BoolFlag   // -t, --verbose    打印执行的命令
	xargsExitOnError  *qflag.BoolFlag   // -e, --exit-on-error  出错立即停止
)

func init() {
	XargsCmd = qflag.NewCmd("xargs", "", qflag.ExitOnError)

	xargsDelimiter = XargsCmd.String("delimiter", "d", "输入分隔符 (默认空格/换行/制表)", "")
	xargsNull = XargsCmd.Bool("null", "0", "使用 \\0 作为分隔符", false)
	xargsArgFile = XargsCmd.String("arg-file", "a", "从文件读取参数", "")
	xargsMaxArgs = XargsCmd.Int("max-args", "n", "每批最大参数个数", 0)
	xargsMaxLines = XargsCmd.Int("max-lines", "L", "每批最大行数", 0)
	xargsMaxChars = XargsCmd.Int("max-chars", "s", "命令最大长度", 0)
	xargsReplaceStr = XargsCmd.Bool("replace", "i", "启用占位符替换模式 (默认占位符 {})", false)
	xargsReplaceDelim = XargsCmd.String("replace-delim", "I", "自定义占位符字符串", "")
	xargsMaxProcs = XargsCmd.Int("max-procs", "P", "并行进程数", 1)
	xargsNoRunIfEmpty = XargsCmd.Bool("no-run-if-empty", "r", "空输入不执行", false)
	xargsInteractive = XargsCmd.Bool("interactive", "p", "执行前确认", false)
	xargsVerbose = XargsCmd.Bool("verbose", "t", "打印执行的命令", false)
	xargsExitOnError = XargsCmd.Bool("exit-on-error", "e", "出错立即停止", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:       "从标准输入或文件读取参数，批量执行指定命令",
		UseChinese: true,
		Examples: map[string]string{
			"基础用法":   `echo a b c | fck xargs echo`,
			"参数替换":   `find . -name "*.log" | fck xargs -i {} mv {} {}.bak`,
			"自定义占位符": `echo a | fck xargs -i -I %% echo %%`,
			"批量处理":   `echo a b c d | fck xargs -n 2 echo`,
			"并行执行":   `cat urls.txt | fck xargs -P 4 curl -O`,
			"从文件读取":  `fck xargs -a files.txt rm`,
			"空分隔符":   `find . -print0 | fck xargs -0 rm`,
			"执行前确认":  `echo file1 | fck xargs -p rm`,
			"限制命令长度": `fck xargs -s 1024 echo`,
			"组合使用":   `find . -name "*.tmp" | fck xargs -0 -i -P 4 mv {} /backup/`,
		},
		Notes: []string{
			"默认模式：所有参数追加到命令后，执行一次",
			"-i 模式：每个参数单独替换占位符，执行多次",
			"使用 -0 可以处理包含空格或特殊字符的文件名",
			"使用 -P 可以并行执行，提高处理速度",
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
		return fmt.Errorf("请指定要执行的命令")
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
		MaxLines:     xargsMaxLines.Get(),
		MaxChars:     xargsMaxChars.Get(),
		ReplaceStr:   replaceStr,
		ReplaceDelim: replaceDelim,
		MaxProcs:     xargsMaxProcs.Get(),
		NoRunIfEmpty: xargsNoRunIfEmpty.Get(),
		Interactive:  xargsInteractive.Get(),
		Verbose:      xargsVerbose.Get(),
		ExitOnError:  xargsExitOnError.Get(),
		Command:      args[0],
		CommandArgs:  args[1:],
	}

	return xargs.XargsCmdMain(config)
}
