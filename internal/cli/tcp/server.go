package tcp

import (
	"fmt"
	"time"

	"gitee.com/MM-Q/fck/internal/commands/tcp"
	"gitee.com/MM-Q/qflag"
)

// ServerCmd tcp server 命令
var ServerCmd *qflag.Cmd

var (
	serverPort       *qflag.IntFlag
	serverAddress    *qflag.StringFlag
	serverTimeout    *qflag.DurationFlag
	serverMaxConn    *qflag.IntFlag
	serverBufferSize *qflag.SizeFlag
	serverResponse   *qflag.StringFlag
	serverOutput     *qflag.StringFlag
	serverEcho       *qflag.BoolFlag
)

func init() {
	ServerCmd = qflag.NewCmd("server", "s", qflag.ExitOnError)

	serverPort = ServerCmd.Int("port", "p", "监听端口", 8888)
	serverAddress = ServerCmd.String("address", "a", "监听地址", "0.0.0.0")
	serverTimeout = ServerCmd.Duration("timeout", "t", "连接超时时间", 30*time.Second)
	serverMaxConn = ServerCmd.Int("max-conn", "m", "最大并发连接数", tcp.GetDefaultConcurrent())
	serverBufferSize = ServerCmd.Size("buffer-size", "b", "接收缓冲区大小", 4*1024)
	serverResponse = ServerCmd.String("response", "r", "响应消息内容", "")
	serverOutput = ServerCmd.String("output", "o", "接收数据保存目录", "")
	serverEcho = ServerCmd.Bool("echo", "e", "回声模式（将接收的数据原样返回）", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "TCP 服务端工具",
		UsageSyntax: fmt.Sprintf("%s tcp server [options]", qflag.Root.Name()),
		UseChinese:  true,
		Notes: []string{
			"默认监听 0.0.0.0:8888",
			"回声模式会将接收的数据原样返回",
			"使用 Ctrl+C 停止服务端",
			"接收的文件会保存在指定的输出目录中",
		},
		Examples: map[string]string{
			"启动默认服务端": fmt.Sprintf("%s tcp server", qflag.Root.Name()),
			"指定端口和地址": fmt.Sprintf("%s tcp server -p 9090 -a 127.0.0.1", qflag.Root.Name()),
			"回声模式":    fmt.Sprintf("%s tcp server -e -p 8888", qflag.Root.Name()),
			"保存接收的数据": fmt.Sprintf("%s tcp server -o /tmp/received -p 8888", qflag.Root.Name()),
			"自定义响应消息": fmt.Sprintf("%s tcp server -r 'Received OK' -p 8888", qflag.Root.Name()),
		},
	}

	if err := ServerCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	ServerCmd.SetRun(runServer)
}

// runServer 运行 server 命令
//
// 参数:
//   - cmd: 命令接口
//
// 返回值:
//   - error: 错误
func runServer(cmd qflag.Command) error {
	config := tcp.ServerConfig{
		Address:    serverAddress.Get(),
		Port:       serverPort.Get(),
		Timeout:    serverTimeout.Get(),
		MaxConn:    serverMaxConn.Get(),
		BufferSize: int(serverBufferSize.Get()),
		Response:   serverResponse.Get(),
		OutputDir:  serverOutput.Get(),
		Echo:       serverEcho.Get(),
	}

	return tcp.ServerCmdMain(config)
}
