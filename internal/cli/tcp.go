package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/cli/tcp"
	"gitee.com/MM-Q/qflag"
)

// TcpCmd tcp 命令
var TcpCmd *qflag.Cmd

func init() {
	TcpCmd = qflag.NewCmd("tcp", "", qflag.ExitOnError)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "TCP 网络工具，支持端口扫描、客户端通信和服务端监听",
		UsageSyntax: fmt.Sprintf("%s tcp [command] [options]", qflag.Root.Name()),
		UseChinese:  true,
		Notes: []string{
			"支持三种模式: scan(扫描), client(客户端), server(服务端)",
			"使用 'fck tcp [command] -h' 查看子命令详细帮助",
		},
		Examples: map[string]string{
			"端口扫描":  fmt.Sprintf("%s tcp scan 192.168.1.1 1-1024", qflag.Root.Name()),
			"发送消息":  fmt.Sprintf("%s tcp client -m 'hello' 192.168.1.1:8080", qflag.Root.Name()),
			"启动服务端": fmt.Sprintf("%s tcp server -p 8080", qflag.Root.Name()),
		},
		SubCmds: []qflag.Command{
			tcp.ScanCmd,
			tcp.ClientCmd,
			tcp.ServerCmd,
		},
	}

	if err := TcpCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	TcpCmd.SetRun(runTcp)
}

// runTcp 运行 tcp 命令
//
// 参数:
//   - cmd: 命令接口
//
// 返回值:
//   - error: 错误
func runTcp(cmd qflag.Command) error {
	cmd.PrintHelp()
	return nil
}
