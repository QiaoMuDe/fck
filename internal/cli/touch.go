package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/touch"
	"gitee.com/MM-Q/qflag"
)

var TouchCmd *qflag.Cmd

var (
	touchNoCreate *qflag.BoolFlag   // 不创建文件
	touchTime     *qflag.StringFlag // 设置时间
	touchAccess   *qflag.BoolFlag   // 只更新访问时间
	touchModify   *qflag.BoolFlag   // 只更新修改时间
	touchVerbose  *qflag.BoolFlag   // 显示操作的文件
)

func init() {
	TouchCmd = qflag.NewCmd("touch", "t", qflag.ExitOnError)

	touchNoCreate = TouchCmd.Bool("no-create", "c", "不创建文件（仅更新时间戳）", false)
	touchTime = TouchCmd.String("time", "t", "设置时间（格式：YYYYMMDDHHMM.SS）", "")
	touchAccess = TouchCmd.Bool("access", "a", "只更新访问时间", false)
	touchModify = TouchCmd.Bool("modify", "m", "只更新修改时间", false)
	touchVerbose = TouchCmd.Bool("verbose", "v", "显示操作的文件", false)

	cmdOpts := &qflag.CmdOpts{
		Desc: "文件创建和时间戳更新工具, 支持创建空文件、更新文件时间戳、设置特定时间等功能",
		Notes: []string{
			"目标通过位置参数传递，支持多个文件",
			"默认创建空文件，使用 -c 选项仅更新时间戳",
			"时间格式为 YYYYMMDDHHMM.SS，如 202401011230.00",
			"使用 -a 选项只更新访问时间，-m 选项只更新修改时间",
		},
		UseChinese: true,
	}

	if err := TouchCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	TouchCmd.SetRun(runTouch)
}

func runTouch(cmd qflag.Command) error {
	config := touch.TouchConfig{
		Targets:  cmd.Args(),
		NoCreate: touchNoCreate.Get(),
		Time:     touchTime.Get(),
		Access:   touchAccess.Get(),
		Modify:   touchModify.Get(),
		Verbose:  touchVerbose.Get(),
	}

	return touch.TouchCmdMain(config)
}
