package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/shck"
	"gitee.com/MM-Q/qflag"
)

var ShckCmd *qflag.Cmd
var shckQuiet *qflag.BoolFlag

func init() {
	ShckCmd = qflag.NewCmd("shck", "", qflag.ExitOnError)
	shckQuiet = ShckCmd.Bool("quiet", "q", "静默模式: 仅输出语法错误，不显示成功信息", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "Shell 脚本语法检查工具",
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s shck [files...]", qflag.Root.Name()),
		Notes: []string{
			"输入方式: 指定文件路径（支持通配符）或通过管道传递脚本内容",
			"支持通配符模式，如: *.sh, src/*.sh",
			"语法错误时显示行号、列号和错误描述",
		},
		Examples: map[string]string{
			"检查单个文件":   fmt.Sprintf("%s shck script.sh", qflag.Root.Name()),
			"通配符批量检查":  fmt.Sprintf("%s shck *.sh", qflag.Root.Name()),
			"递归检查所有脚本": fmt.Sprintf("%s shck **/*.sh", qflag.Root.Name()),
			"管道输入检查":   fmt.Sprintf("cat script.sh | %s shck", qflag.Root.Name()),
		},
	}

	if err := ShckCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	ShckCmd.SetRun(runShck)
}

// runShck 执行 shck 命令
func runShck(cmd qflag.Command) error {
	config := shck.ShckConfig{
		Files: cmd.Args(),
		Quiet: shckQuiet.Get(),
	}

	return shck.ShckCmdMain(config)
}
