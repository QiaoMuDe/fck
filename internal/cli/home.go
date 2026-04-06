package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/home"
	"gitee.com/MM-Q/qflag"
)

var HomeCmd *qflag.Cmd

func init() {
	HomeCmd = qflag.NewCmd("home", "", qflag.ExitOnError)

	cmdOpts := &qflag.CmdOpts{
		Desc: "显示用户主目录",
		Notes: []string{
			"直接显示用户主目录路径，无需任何选项",
			"Windows 平台: C:\\Users\\username",
			"Linux 平台: /home/username",
			"macOS 平台: /Users/username",
		},
		UseChinese: true,
	}

	if err := HomeCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	HomeCmd.SetRun(runHome)
}

func runHome(cmd qflag.Command) error {
	config := home.HomeConfig{}

	return home.HomeCmdMain(config)
}
