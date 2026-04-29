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
	clientMessage     *qflag.StringFlag
	clientPath        *qflag.StringFlag
	clientInteractive *qflag.BoolFlag
	clientTimeout     *qflag.DurationFlag
	clientBufferSize  *qflag.SizeFlag
	clientNoResponse  *qflag.BoolFlag
	clientDelimiter   *qflag.StringFlag
)

func init() {
	ClientCmd = qflag.NewCmd("client", "c", qflag.ExitOnError)

	clientMessage = ClientCmd.String("message", "m", "要发送的字符串消息", "")
	clientPath = ClientCmd.String("path", "p", "要发送的文件/目录路径, 支持通配符 (如 *.txt) ", "")
	clientInteractive = ClientCmd.Bool("interactive", "i", "交互式模式, 持续输入发送", false)
	clientTimeout = ClientCmd.Duration("timeout", "t", "连接和传输超时时间", 10*time.Second)
	clientBufferSize = ClientCmd.Size("buffer-size", "b", "发送缓冲区大小", 4*1024)
	clientNoResponse = ClientCmd.Bool("no-response", "n", "不等待服务器响应", false)
	clientDelimiter = ClientCmd.String("delimiter", "D", "消息分隔符 (交互式模式使用) ", "EOF")

	cmdOpts := &qflag.CmdOpts{
		Desc:        "TCP 客户端工具",
		UsageSyntax: fmt.Sprintf("%s tcp client [options] <address:port>", qflag.Root.Name()),
		UseChinese:  true,
		Notes: []string{
			"-m, -p, -i 三个选项互斥, 必须指定其中一个",
			"交互式模式下使用分隔符退出 (默认 EOF) ",
			"路径参数支持文件、目录 (不递归子目录) 和通配符",
			"三种发送模式都会处理服务端返回的数据包",
		},
		Examples: map[string]string{
			"发送字符串":     fmt.Sprintf("%s tcp client -m 'Hello Server' 192.168.1.1:8080", qflag.Root.Name()),
			"发送单个文件":    fmt.Sprintf("%s tcp client -p /path/to/file.txt 192.168.1.1:8080", qflag.Root.Name()),
			"发送目录下所有文件": fmt.Sprintf("%s tcp client -p /path/to/dir 192.168.1.1:8080", qflag.Root.Name()),
			"发送通配符匹配文件": fmt.Sprintf("%s tcp client -p '/path/to/*.txt' 192.168.1.1:8080", qflag.Root.Name()),
			"交互式模式":     fmt.Sprintf("%s tcp client -i 192.168.1.1:8080", qflag.Root.Name()),
			"不等待响应":     fmt.Sprintf("%s tcp client -m 'hello' -n 192.168.1.1:8080", qflag.Root.Name()),
		},
		MutexGroups: []qflag.MutexGroup{
			{
				Name:      "send-mode",
				Flags:     []string{"message", "path", "interactive"},
				AllowNone: true,
			},
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

	// 检查是否指定了发送模式
	if clientMessage.Get() == "" && clientPath.Get() == "" && !clientInteractive.Get() {
		return fmt.Errorf("must specify send mode: -m (message), -p (path) or -i (interactive)")
	}

	config := tcp.ClientConfig{
		Address:     args[0],
		Message:     clientMessage.Get(),
		Path:        clientPath.Get(),
		Interactive: clientInteractive.Get(),
		Timeout:     clientTimeout.Get(),
		BufferSize:  int(clientBufferSize.Get()),
		NoResponse:  clientNoResponse.Get(),
		Delimiter:   clientDelimiter.Get(),
	}

	return tcp.ClientCmdMain(config)
}
