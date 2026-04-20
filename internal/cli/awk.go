package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/awk"
	"gitee.com/MM-Q/qflag"
)

var AwkCmd *qflag.Cmd

var (
	awkPattern   *qflag.StringFlag   // -p, --pattern  正则匹配模式
	awkField     *qflag.IntSliceFlag // -f, --field    字段索引
	awkChars     *qflag.StringFlag   // -c, --chars    字符范围
	awkSeparator *qflag.StringFlag   // -F, --separator 输入分隔符
	awkOutputSep *qflag.StringFlag   // -O, --output-separator 输出分隔符
	awkLineNum   *qflag.BoolFlag     // -n, --line-number 显示行号
)

func init() {
	AwkCmd = qflag.NewCmd("awk", "", qflag.ExitOnError)

	awkPattern = AwkCmd.String("pattern", "p", "正则匹配模式，只处理匹配的行", "")
	awkField = AwkCmd.IntSlice("field", "f", "要输出的字段索引 (1-based, 0=整行, -1=最后一个字段)", []int{})
	awkChars = AwkCmd.String("chars", "c", "要输出的字符范围 (如: 1-10,5-,-20)", "")
	awkSeparator = AwkCmd.String("separator", "F", "输入字段分隔符（默认任意空白）", "")
	awkOutputSep = AwkCmd.String("output-separator", "O", "输出字段分隔符", " ")
	awkLineNum = AwkCmd.Bool("line-number", "n", "显示行号", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:       "文本字段处理工具",
		UseChinese: true,
		Examples: map[string]string{
			"打印第1列":      "fck awk -f 1 file.txt",
			"打印第1和第3列":   "fck awk -f 1,3 file.txt",
			"指定分隔符":      "fck awk -F: -f 1,3 /etc/passwd",
			"打印最后一列":     "fck awk -f -1 file.txt",
			"正则匹配并打印":    "fck awk -p \"error\" -f 0 log.txt",
			"显示行号":       "fck awk -n -f 1 file.txt",
			"CSV处理":      "fck awk -F, -f 1,3 data.csv",
			"制表符分隔":      "fck awk -F\"\\t\" -f 2 file.tsv",
			"自定义输出格式":    "fck awk -F: -O\" -> \" -f 1,6 /etc/passwd",
			"管道输入":       "cat file.txt | fck awk -f 1",
			"提取前10个字符":   "fck awk -c 1-10 file.txt",
			"提取第5个字符到行尾": "fck awk -c 5- file.txt",
			"提取多个字符范围":   "fck awk -c 1,5-10,15 file.txt",
		},
		Notes: []string{
			"字段索引从1开始, 0表示整行, -1表示最后一个字段",
			"字符范围格式: N(单个), N-M(范围), N-(从N到结尾), -N(前N个)",
			"-f 和 -c 参数互斥，不能同时使用",
			"支持从管道读取输入, 此时不需要指定文件参数",
			"正则匹配使用 Go 正则语法",
			"默认分隔符为空格, 支持任意字符串作为分隔符",
		},
		MutexGroups: []qflag.MutexGroup{
			{
				Name:      "extract",
				Flags:     []string{"field", "chars"},
				AllowNone: true,
			},
		},
	}

	if err := AwkCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	AwkCmd.SetRun(runAwk)
}

func runAwk(cmd qflag.Command) error {
	config := awk.AwkConfig{
		Pattern:     awkPattern.Get(),
		Fields:      awkField.Get(),
		FieldSep:    awkSeparator.Get(),
		OutputSep:   awkOutputSep.Get(),
		ShowLineNum: awkLineNum.Get(),
		Files:       cmd.Args(),
	}

	// 如果没有指定字段，默认输出整行
	if len(config.Fields) == 0 {
		config.Fields = []int{0}
	}

	// 解析字符范围
	charsStr := awkChars.Get()
	if charsStr != "" {
		ranges, err := awk.ParseCharRanges(charsStr)
		if err != nil {
			return fmt.Errorf("invalid char ranges: %w", err)
		}
		config.CharRanges = ranges
	}

	return awk.AwkCmdMain(config)
}
