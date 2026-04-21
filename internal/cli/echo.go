package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/echo"
	"gitee.com/MM-Q/qflag"
)

var EchoCmd *qflag.Cmd

var (
	echoNoNewline *qflag.BoolFlag   // 不在末尾添加换行符
	echoRaw       *qflag.BoolFlag   // 原始输出，不解析转义字符
	echoTrim      *qflag.BoolFlag   // 去除首尾空白字符
	echoUpper     *qflag.BoolFlag   // 转换为大写
	echoLower     *qflag.BoolFlag   // 转换为小写
	echoRepeat    *qflag.IntFlag    // 重复输出次数
	echoColor     *qflag.StringFlag // 颜色输出
)

func init() {
	EchoCmd = qflag.NewCmd("echo", "", qflag.ExitOnError)

	echoNoNewline = EchoCmd.Bool("no-newline", "n", "不在末尾添加换行符", false)
	echoRaw = EchoCmd.Bool("raw", "r", "原始输出，不解析转义字符", false)
	echoTrim = EchoCmd.Bool("trim", "T", "去除首尾空白字符", false)
	echoUpper = EchoCmd.Bool("upper", "u", "转换为大写", false)
	echoLower = EchoCmd.Bool("lower", "l", "转换为小写", false)
	echoRepeat = EchoCmd.Int("repeat", "R", "重复输出次数", 1)
	echoColor = EchoCmd.String("color", "c", "颜色输出（red, green, yellow, blue等）", "")

	cmdOpts := &qflag.CmdOpts{
		Desc:        "文本输出工具",
		UsageSyntax: fmt.Sprintf("%s echo [options] [text...]", qflag.Root.Name()),
		Notes: []string{
			"支持转义字符: \\n(换行), \\t(制表符), \\r(回车), \\\\(反斜杠), \\\"(双引号)",
			"支持颜色: black, red, green, yellow, blue, magenta, cyan, white",
			"支持亮色: bright-black, bright-red, bright-green, bright-yellow, bright-blue, bright-magenta, bright-cyan, bright-white",
		},
		UseChinese: true,
		Examples: map[string]string{
			"输出文本": fmt.Sprintf("%s echo \"Hello, World!\"", qflag.Root.Name()),
		},
	}

	if err := EchoCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	EchoCmd.SetRun(runEcho)
}

func runEcho(cmd qflag.Command) error {
	args := cmd.Args()
	text := ""
	for i, arg := range args {
		if i > 0 {
			text += " "
		}
		text += arg
	}

	config := echo.EchoConfig{
		Text:      text,
		NoNewline: echoNoNewline.Get(),
		Raw:       echoRaw.Get(),
		Trim:      echoTrim.Get(),
		Upper:     echoUpper.Get(),
		Lower:     echoLower.Get(),
		Repeat:    echoRepeat.Get(),
		Color:     echoColor.Get(),
	}

	return echo.EchoCmdMain(config)
}
