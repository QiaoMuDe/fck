package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/tr"
	"gitee.com/MM-Q/qflag"
)

var TrCmd *qflag.Cmd

var (
	trDelete         *qflag.BoolFlag // -d, --delete
	trSqueezeRepeats *qflag.BoolFlag // -s, --squeeze-repeats
	trComplement     *qflag.BoolFlag // -c, --complement
	trTruncateSet1   *qflag.BoolFlag // -t, --truncate-set1
)

func init() {
	TrCmd = qflag.NewCmd("tr", "", qflag.ExitOnError)

	trDelete = TrCmd.Bool("delete", "d", "删除 set1 中的字符", false)
	trSqueezeRepeats = TrCmd.Bool("squeeze-repeats", "s", "压缩 set1 中连续重复的字符", false)
	trComplement = TrCmd.Bool("complement", "c", "使用 set1 的补集", false)
	trTruncateSet1 = TrCmd.Bool("truncate-set1", "t", "将 set1 截断为 set2 的长度", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "字符转换工具",
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s tr [options] <set1> [set2]", qflag.Root.Name()),
		Notes: []string{
			"从标准输入读取数据，转换后输出到标准输出",
			"set1 和 set2 支持以下格式：\n\t- 单个字符: 'abc', '123' \n\t- 字符范围: 'a-z', 'A-Z', '0-9' \n\t- 转义字符: \\n (换行), \\t (制表符), \\\\ (反斜杠)",
			"转换模式需要同时指定 set1 和 set2",
			"删除模式 (-d) 和压缩模式 (-s) 只需要 set1",
		},
		Examples: map[string]string{
			"大小写转换": fmt.Sprintf("echo 'hello' | %s tr 'a-z' 'A-Z'", qflag.Root.Name()),
			"删除数字":  fmt.Sprintf("echo 'abc123' | %s tr -d '0-9'", qflag.Root.Name()),
			"压缩空格":  fmt.Sprintf("echo 'a    b' | %s tr -s ' '", qflag.Root.Name()),
			"数字转X":  fmt.Sprintf("echo 'abc123' | %s tr '0-9' 'X'", qflag.Root.Name()),
			"删除非字母": fmt.Sprintf("echo 'a1b2c3' | %s tr -c -d 'a-zA-Z'", qflag.Root.Name()),
		},
		MutexGroups: []qflag.MutexGroup{
			{
				Name: "operation-mode",
				Flags: []string{
					"delete",
					"squeeze-repeats",
				},
				AllowNone: true,
			},
			{
				Name: "truncate-with-mode",
				Flags: []string{
					"truncate-set1",
					"delete",
					"squeeze-repeats",
				},
				AllowNone: true,
			},
		},
	}

	if err := TrCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	TrCmd.SetRun(runTr)
}

func runTr(cmd qflag.Command) error {
	args := cmd.Args()
	if len(args) < 1 {
		return fmt.Errorf("set1 is required")
	}

	config := tr.TrConfig{
		Set1:           args[0],
		Delete:         trDelete.Get(),
		SqueezeRepeats: trSqueezeRepeats.Get(),
		Complement:     trComplement.Get(),
		TruncateSet1:   trTruncateSet1.Get(),
	}

	if len(args) >= 2 {
		config.Set2 = args[1]
	}

	return tr.TrCmdMain(config)
}
