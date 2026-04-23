package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/seq"
	"gitee.com/MM-Q/qflag"
)

var SeqCmd *qflag.Cmd

var (
	seqSeparator  *qflag.StringFlag // 分隔符
	seqTerminator *qflag.StringFlag // 终止符
	seqEqualWidth *qflag.BoolFlag   // 等宽输出
)

func init() {
	SeqCmd = qflag.NewCmd("seq", "", qflag.ExitOnError)

	seqSeparator = SeqCmd.String("separator", "s", "指定分隔符（支持 \\n \\t 等转义字符）", "\n")
	seqTerminator = SeqCmd.String("terminator", "t", "指定终止符（支持 \\n \\t 等转义字符）", "\n")
	seqEqualWidth = SeqCmd.Bool("equal-width", "w", "等宽输出（自动补零对齐）", false)

	cmdOpts := &qflag.CmdOpts{
		Desc: "生成数字序列",
		Notes: []string{
			"支持 1-3 个位置参数：",
			"  1个参数: seq LAST          (从 1 到 LAST)",
			"  2个参数: seq FIRST LAST    (从 FIRST 到 LAST，步长为 1)",
			"  3个参数: seq FIRST STEP LAST (从 FIRST 到 LAST，步长为 STEP)",
			"",
			"支持小数和负数，支持反向序列（负步长）",
			"",
			"注意：使用负数参数时，需要在选项后添加 -- 停止标志解析",
			"  例如: seq -- 5 -1 1",
		},
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s seq [options] <last> | <first> <last> | <first> <step> <last>", qflag.Root.Name()),
		Examples: map[string]string{
			"生成 1 到 5":       fmt.Sprintf("%s seq 5", qflag.Root.Name()),
			"生成 2 到 5":       fmt.Sprintf("%s seq 2 5", qflag.Root.Name()),
			"生成 0 到 10，步长 2": fmt.Sprintf("%s seq 0 2 10", qflag.Root.Name()),
			"小数序列":           fmt.Sprintf("%s seq 0 0.5 2", qflag.Root.Name()),
			"反向序列":           fmt.Sprintf("%s seq -- 5 -1 1", qflag.Root.Name()),
			"空格分隔":           fmt.Sprintf("%s seq -s \" \" 5", qflag.Root.Name()),
			"制表符分隔":          fmt.Sprintf("%s seq -s \"\\t\" 5", qflag.Root.Name()),
			"等宽输出":           fmt.Sprintf("%s seq -w 99 101", qflag.Root.Name()),
		},
	}

	if err := SeqCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	SeqCmd.SetRun(runSeq)
}

func runSeq(cmd qflag.Command) error {
	args := cmd.Args()
	if len(args) == 0 {
		return fmt.Errorf("requires at least 1 argument")
	}

	// 解析位置参数
	start, end, step, err := seq.ParseArgs(args)
	if err != nil {
		return err
	}

	config := seq.SeqConfig{
		Start:      start,
		End:        end,
		Step:       step,
		Separator:  seqSeparator.Get(),
		Terminator: seqTerminator.Get(),
		EqualWidth: seqEqualWidth.Get(),
	}

	return seq.SeqCmdMain(config)
}
