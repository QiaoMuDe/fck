package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/base64"
	"gitee.com/MM-Q/qflag"
)

var Base64Cmd *qflag.Cmd

var (
	base64Decode    *qflag.BoolFlag   // 解码模式
	base64File      *qflag.StringFlag // 输入文件路径
	base64Output    *qflag.StringFlag // 输出文件路径
	base64Wrap      *qflag.IntFlag    // 每行字符数
	base64URLSafe   *qflag.BoolFlag   // URL 安全变体
	base64NoPadding *qflag.BoolFlag   // 禁用填充
)

func init() {
	Base64Cmd = qflag.NewCmd("base64", "", qflag.ExitOnError)

	base64Decode = Base64Cmd.Bool("decode", "d", "解码模式", false)
	base64File = Base64Cmd.String("file", "f", "从文件读取输入", "")
	base64Output = Base64Cmd.String("output", "o", "输出文件路径", "")
	base64Wrap = Base64Cmd.Int("wrap", "w", "每行最大字符数 (0=不换行) ", 76)
	base64URLSafe = Base64Cmd.Bool("url-safe", "u", "使用 URL 安全变体 (- 替换 +, _ 替换 /) ", false)
	base64NoPadding = Base64Cmd.Bool("no-padding", "p", "禁用填充 (=) ", false)

	cmdOpts := &qflag.CmdOpts{
		Desc: "Base64 编解码工具",
		Notes: []string{
			"优先从管道/重定向读取输入",
			"其次使用 -f 选项从文件读取",
			"最后编码位置参数字符串",
			"使用 -d 选项切换到解码模式",
			"URL 安全变体将 + 替换为 -, / 替换为 _",
		},
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s base64 [options] [string...]", qflag.Root.Name()),
		Examples: map[string]string{
			"编码字符串":    fmt.Sprintf("%s base64 'hello world'", qflag.Root.Name()),
			"编码多个词":    fmt.Sprintf("%s base64 hello world", qflag.Root.Name()),
			"编码文件":     fmt.Sprintf("%s base64 -f file.txt", qflag.Root.Name()),
			"解码字符串":    fmt.Sprintf("%s base64 -d 'aGVsbG8='", qflag.Root.Name()),
			"解码文件":     fmt.Sprintf("%s base64 -d -f encoded.txt", qflag.Root.Name()),
			"不换行输出":    fmt.Sprintf("%s base64 -w 0 'hello'", qflag.Root.Name()),
			"URL 安全编码": fmt.Sprintf("%s base64 -u 'hello'", qflag.Root.Name()),
			"输出到文件":    fmt.Sprintf("%s base64 -o out.txt 'hello'", qflag.Root.Name()),
			"管道使用":     fmt.Sprintf("echo 'hello' | %s base64", qflag.Root.Name()),
		},
	}

	if err := Base64Cmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	Base64Cmd.SetRun(runBase64)
}

func runBase64(cmd qflag.Command) error {
	config := base64.Base64Config{
		Decode:    base64Decode.Get(),
		Strings:   cmd.Args(),
		FilePath:  base64File.Get(),
		Output:    base64Output.Get(),
		Wrap:      base64Wrap.Get(),
		URLSafe:   base64URLSafe.Get(),
		NoPadding: base64NoPadding.Get(),
	}

	return base64.Base64CmdMain(config)
}
