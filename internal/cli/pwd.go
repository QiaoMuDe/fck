package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/pwd"
	"gitee.com/MM-Q/qflag"
)

var PwdCmd *qflag.Cmd

var (
	pwdPhysical *qflag.BoolFlag // 显示物理路径（解析符号链接）
	pwdLogical  *qflag.BoolFlag // 显示逻辑路径（不解析符号链接）
)

func init() {
	PwdCmd = qflag.NewCmd("pwd", "", qflag.ExitOnError)

	pwdPhysical = PwdCmd.Bool("physical", "P", "显示物理路径（解析符号链接）", false)
	pwdLogical = PwdCmd.Bool("logical", "L", "显示逻辑路径（不解析符号链接）", true)

	cmdOpts := &qflag.CmdOpts{
		Desc: "打印当前工作目录",
		Notes: []string{
			"默认显示逻辑路径（不解析符号链接）",
			"使用 -P 选项显示物理路径（解析符号链接）",
			"使用 -L 选项显式显示逻辑路径（不解析符号链接）",
		},
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s pwd [options]", qflag.Root.Name()),
		Examples: map[string]string{
			"显示当前工作目录": fmt.Sprintf("%s pwd", qflag.Root.Name()),
		},
	}

	if err := PwdCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	PwdCmd.SetRun(runPwd)
}

func runPwd(cmd qflag.Command) error {
	config := pwd.PwdConfig{
		Physical: pwdPhysical.Get(),
		Logical:  pwdLogical.Get(),
	}

	return pwd.PwdCmdMain(config)
}
