package cli

import (
	"fmt"
	"strings"

	"gitee.com/MM-Q/fck/internal/commands/shx"
	"gitee.com/MM-Q/go-kit/term"
	"gitee.com/MM-Q/qflag"
)

var ShxCmd *qflag.Cmd

var (
	shxTimeout *qflag.DurationFlag    // -t, --timeout    超时时间
	shxDir     *qflag.StringFlag      // -d, --dir        工作目录
	shxEnv     *qflag.StringSliceFlag // -e, --env     环境变量
)

func init() {
	ShxCmd = qflag.NewCmd("shx", "", qflag.ExitOnError)

	shxTimeout = ShxCmd.Duration("timeout", "t", "命令执行超时时间 (如: 5s, 1m)", 0)
	shxDir = ShxCmd.String("dir", "d", "指定工作目录", "")
	shxEnv = ShxCmd.StringSlice("env", "e", "设置环境变量 (可多次指定, 如: -e KEY=VALUE)", nil)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "Shell 命令/脚本执行工具",
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s shx [options] <command|script>", qflag.Root.Name()),
		Notes: []string{
			"执行命令字符串: 直接传入命令，如 \"echo hello\"",
			"执行脚本文件: 传入 .sh 文件路径，如 deploy.sh",
			"脚本文件需以 .sh 结尾且文件存在，否则视为命令字符串",
			"支持管道传递 stdin 到执行的命令",
		},
		Examples: map[string]string{
			"执行命令":       fmt.Sprintf("%s shx \"echo hello\"", qflag.Root.Name()),
			"执行脚本":       fmt.Sprintf("%s shx deploy.sh", qflag.Root.Name()),
			"超时控制":       fmt.Sprintf("%s shx -t 5s \"sleep 10\"", qflag.Root.Name()),
			"指定工作目录":     fmt.Sprintf("%s shx -d /tmp pwd", qflag.Root.Name()),
			"设置环境变量":     fmt.Sprintf("%s shx -e FOO=bar \"echo $FOO\"", qflag.Root.Name()),
			"管道传递 stdin": fmt.Sprintf("echo data | %s shx cat", qflag.Root.Name()),
		},
	}

	if err := ShxCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	ShxCmd.SetRun(runShx)
}

// runShx 执行 shx 命令
func runShx(cmd qflag.Command) error {
	args := cmd.Args()
	if len(args) == 0 {
		return fmt.Errorf("command or script path is required")
	}

	// 将位置参数拼接为命令字符串
	cmdStr := strings.Join(args, " ")

	// 检测是否有管道输入
	hasStdin := term.IsStdinPipe()

	config := shx.ShxConfig{
		Timeout: shxTimeout.Get(),
		Dir:     shxDir.Get(),
		Env:     shxEnv.Get(),
		Cmd:     cmdStr,
		Stdin:   hasStdin,
	}

	return shx.ShxCmdMain(config)
}
