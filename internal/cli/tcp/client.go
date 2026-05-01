package tcp

import (
	"fmt"
	"time"

	"gitee.com/MM-Q/fck/internal/commands/tcp"
	"gitee.com/MM-Q/qflag"
)

// ClientCmd tcp client 命令
var ClientCmd *qflag.Cmd

var (
	clientMessage    *qflag.StringFlag
	clientTimeout    *qflag.DurationFlag
	clientBufferSize *qflag.SizeFlag
	clientNoResponse *qflag.BoolFlag
	clientDelimiter  *qflag.StringFlag
	clientHex        *qflag.BoolFlag
)

func init() {
	ClientCmd = qflag.NewCmd("client", "c", qflag.ExitOnError)

	clientMessage = ClientCmd.String("message", "m", "要发送的字符串消息", "")
	clientTimeout = ClientCmd.Duration("timeout", "t", "连接和传输超时时间", 3*time.Second)
	clientBufferSize = ClientCmd.Size("buffer-size", "b", "发送缓冲区大小", 4*1024)
	clientNoResponse = ClientCmd.Bool("no-response", "n", "不等待服务器响应", false)
	clientDelimiter = ClientCmd.String("delimiter", "D", "消息分隔符 (交互式模式使用) ", "EOF")
	clientHex = ClientCmd.Bool("hex", "x", "将消息作为十六进制数据发送", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "TCP 客户端工具",
		UsageSyntax: fmt.Sprintf("%s tcp client [options] <address:port>", qflag.Root.Name()),
		UseChinese:  true,
		Notes: []string{
			"支持三种发送模式（按优先级）：",
			"  1. 管道输入: echo 'data' | fck tcp client <address:port>",
			"  2. 字符串模式: -m '消息内容'（与管道互斥）",
			"  3. 交互模式: 无标志时默认进入（使用分隔符 EOF 退出）",
			"管道模式自动检测，无需额外标志",
		},
		Examples: map[string]string{
			"管道发送":     fmt.Sprintf("echo 'Hello Server' | %s tcp client 192.168.1.1:8080", qflag.Root.Name()),
			"文件内容管道发送": fmt.Sprintf("cat data.txt | %s tcp client 192.168.1.1:8080", qflag.Root.Name()),
			"发送字符串":    fmt.Sprintf("%s tcp client -m 'Hello Server' 192.168.1.1:8080", qflag.Root.Name()),
			"发送十六进制":   fmt.Sprintf("%s tcp client -m '48656c6c6f' -x 192.168.1.1:8080", qflag.Root.Name()),
			"交互式模式":    fmt.Sprintf("%s tcp client 192.168.1.1:8080", qflag.Root.Name()),
			"不等待响应":    fmt.Sprintf("%s tcp client -m 'hello' -n 192.168.1.1:8080", qflag.Root.Name()),
		},
	}

	if err := ClientCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	ClientCmd.SetRun(runClient)
}

// runClient 运行 client 命令
//
// 参数:
//   - cmd: 命令接口
//
// 返回值:
//   - error: 错误
func runClient(cmd qflag.Command) error {
	args := cmd.Args()
	if len(args) == 0 {
		return fmt.Errorf("server address not specified")
	}

	config := tcp.ClientConfig{
		Address:    args[0],
		Message:    clientMessage.Get(),
		Timeout:    clientTimeout.Get(),
		BufferSize: int(clientBufferSize.Get()),
		NoResponse: clientNoResponse.Get(),
		Delimiter:  clientDelimiter.Get(),
		Hex:        clientHex.Get(),
	}

	return tcp.ClientCmdMain(config)
}
