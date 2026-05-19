package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/hex2str"
	"gitee.com/MM-Q/qflag"
)

var Hex2strCmd *qflag.Cmd

var (
	hex2strEncode *qflag.BoolFlag // -e 编码模式（字符串→十六进制）
	hex2strUpper  *qflag.BoolFlag // -u 编码输出使用大写字母
)

func init() {
	Hex2strCmd = qflag.NewCmd("hex2str", "h2s", qflag.ExitOnError)

	hex2strEncode = Hex2strCmd.Bool("encode", "e", "编码模式：将字符串转换为十六进制", false)
	hex2strUpper = Hex2strCmd.Bool("upper", "u", "编码输出使用大写字母 (A-F)", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "十六进制与字符串相互转换工具",
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s hex2str [options] [data]", qflag.Root.Name()),
		Examples: map[string]string{
			"十六进制转字符串":  fmt.Sprintf("%s hex2str 48656c6c6f", qflag.Root.Name()),
			"字符串转十六进制":  fmt.Sprintf("%s hex2str -e Hello", qflag.Root.Name()),
			"大写十六进制输出":  fmt.Sprintf("%s hex2str -e -u Hello", qflag.Root.Name()),
			"管道解码":      fmt.Sprintf("echo '48656c6c6f' | %s hex2str", qflag.Root.Name()),
			"管道编码":      fmt.Sprintf("echo 'Hello' | %s hex2str -e", qflag.Root.Name()),
			"处理带空格十六进制": fmt.Sprintf("%s hex2str '48 65 6C 6C 6F'", qflag.Root.Name()),
			"处理 0x 前缀":  fmt.Sprintf("%s hex2str '0x48656c6c6f'", qflag.Root.Name()),
		},
		Notes: []string{
			"默认模式：十六进制→字符串解码",
			"使用 -e 切换到编码模式：字符串→十六进制",
			"支持从管道/stdin 读取输入",
			"解码时自动处理空格、0x 前缀",
		},
	}

	if err := Hex2strCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	Hex2strCmd.SetRun(runHex2str)
}

// runHex2str 执行 hex2str 命令
//
// 参数:
//   - cmd: qflag 命令对象
//
// 返回值:
//   - error: 执行错误
func runHex2str(cmd qflag.Command) error {
	// 获取位置参数
	args := cmd.Args()
	input := ""
	if len(args) > 0 {
		input = args[0]
	}

	config := hex2str.Hex2strConfig{
		Encode: hex2strEncode.Get(),
		Upper:  hex2strUpper.Get(),
		Input:  input,
	}

	return hex2str.Hex2strCmdMain(config)
}
