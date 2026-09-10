package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/wc"
	"gitee.com/MM-Q/fck/internal/types"
	"gitee.com/MM-Q/qflag"
)

var WcCmd *qflag.Cmd

var (
	wcLines      *qflag.BoolFlag // 显示行数
	wcWords      *qflag.BoolFlag // 显示单词数
	wcBytes      *qflag.BoolFlag // 显示字节数
	wcChars      *qflag.BoolFlag // 显示字符数
	wcMaxLine    *qflag.BoolFlag // 显示最大行长度
	wcAll        *qflag.BoolFlag // 显示所有统计项
	wcTableStyle *qflag.EnumFlag // 表格样式
)

func init() {
	WcCmd = qflag.NewCmd("wc", "", qflag.ExitOnError)

	wcLines = WcCmd.Bool("lines", "l", "仅统计行数", false)
	wcWords = WcCmd.Bool("words", "w", "仅统计单词数", false)
	wcBytes = WcCmd.Bool("bytes", "c", "仅统计字节数", false)
	wcChars = WcCmd.Bool("chars", "m", "仅统计字符数", false)
	wcMaxLine = WcCmd.Bool("max-line-length", "L", "仅统计最大行长度", false)
	wcAll = WcCmd.Bool("all", "a", "显示所有统计项 (行数、单词数、字符数、字节数、最大行宽)", false)
	wcTableStyle = WcCmd.Enum("table-style", "ts", qflag.EnumHelp("指定表格样式，支持以下选项:", types.TableStyleOptions, "\t\t\t\t\t"), "none", types.TableStyles)

	cmdOpts := &qflag.CmdOpts{
		Desc:       "统计文件的行数、单词数、字节数和字符数",
		UseChinese: true,
		Examples: map[string]string{
			"默认统计 (行数、单词、字节)": fmt.Sprintf("%s wc file.txt", qflag.Root.Name()),
			"仅统计行数":           fmt.Sprintf("%s wc -l file.txt", qflag.Root.Name()),
			"统计行数和单词数":        fmt.Sprintf("%s wc -lw file.txt", qflag.Root.Name()),
			"统计字符数 (Unicode)": fmt.Sprintf("%s wc -m file.txt", qflag.Root.Name()),
			"统计最大行长度":         fmt.Sprintf("%s wc -L file.txt", qflag.Root.Name()),
			"多文件统计":           fmt.Sprintf("%s wc file1.txt file2.txt", qflag.Root.Name()),
			"通配符支持":           fmt.Sprintf("%s wc *.txt", qflag.Root.Name()),
			"管道输入":            fmt.Sprintf("%s wc file.txt | %s wc", qflag.Root.Name(), qflag.Root.Name()),
			"统计代码行数":          fmt.Sprintf("%s wc -l *.go", qflag.Root.Name()),
		},
		Notes: []string{
			"无标志时默认输出：行数、单词数、字节数",
			"使用 -a/--all 可显示所有统计项",
			"多标志可同时使用，按固定顺序输出：行数、单词、字符、字节、最大行宽",
			"支持通配符展开 (跨平台兼容)",
			"管道输入时优先处理管道，忽略文件参数",
			"字符数 (-m) 统计 Unicode 字符，与字节数 (-c)不同",
		},
	}

	if err := WcCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	WcCmd.SetRun(runWc)
}

func runWc(cmd qflag.Command) error {
	config := wc.WcConfig{
		ShowLines:   wcLines.Get(),
		ShowWords:   wcWords.Get(),
		ShowBytes:   wcBytes.Get(),
		ShowChars:   wcChars.Get(),
		ShowMaxLine: wcMaxLine.Get(),
		ShowAll:     wcAll.Get(),
		TableStyle:  wcTableStyle.Get(),
		Files:       cmd.Args(),
	}

	return wc.WcCmdMain(config)
}
