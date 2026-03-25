package cli

import (
	"fmt"
	"runtime/debug"

	"gitee.com/MM-Q/fck/internal/types"
	"gitee.com/MM-Q/qflag"
	"gitee.com/MM-Q/verman"
)

func InitAndRun() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v\n%s", r, debug.Stack())
		}
	}()

	cmdOpts := &qflag.CmdOpts{
		Desc:       "多功能文件处理工具集, 提供文件哈希计算、大小统计、查找和校验等实用功能",
		UseChinese: true,
		Completion: true,
		Version:    verman.V.Version(),
		RunFunc:    run,
		LogoText:   types.FckHelpLogo,
		Notes: []string{
			"各子命令有独立帮助文档，可通过-h参数查看, 例如 'fck <子命令> -h' 查看各子命令详细帮助",
		},
		SubCmds: []qflag.Command{
			PackCmd,
			PreviewCmd,
			UnpackCmd,
			WatchCmd,
			CheckCmd,
			FindCmd,
			HashCmd,
			ListCmd,
			SizeCmd,
			DateCmd,
			EchoCmd,
			RMCmd,
			MkdirCmd,
			TouchCmd,
			ChmodCmd,
			PwdCmd,
			HomeCmd,
			TruncateCmd,
			CpCmd,
			MvCmd,
			TestCmd,
			GmCmd,
			XargsCmd,
		},
	}

	if err := qflag.ApplyOpts(cmdOpts); err != nil {
		err = fmt.Errorf("apply opts error: %v", err)
		return err
	}

	if err := qflag.ParseAndRoute(); err != nil {
		err = fmt.Errorf("parse and route error: %v", err)
		return err
	}

	return nil
}

// 根命令运行函数
func run(cmd qflag.Command) error {
	qflag.Root.PrintHelp()
	return nil
}
