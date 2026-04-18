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
		Desc:              "一站式文件与系统管理工具集，集成文件操作、内容查看、压缩解压、哈希校验、系统信息等多种实用功能",
		UseChinese:        true,
		Completion:        true,
		DynamicCompletion: true,
		Version:           verman.V.Version(),
		RunFunc:           run,
		LogoText:          types.FckHelpLogo,
		Notes: []string{
			fmt.Sprintf("各子命令有独立帮助文档，可通过 --help/-h 参数查看, 例如 '%s <子命令> -h' 查看各子命令详细帮助", qflag.Root.Name()),
		},
		SubCmds: []qflag.Command{
			PackCmd,
			PreviewCmd,
			UnpackCmd,
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
			PwdCmd,
			HomeCmd,
			TruncateCmd,
			CpCmd,
			MvCmd,
			TestCmd,
			GmCmd,
			CatCmd,
			GrepCmd,
			SedCmd,
			AwkCmd,
			XargsCmd,
			AliasCmd,
			MdCmd,
			WatchCmd,
		},
	}

	if err := qflag.ApplyOpts(cmdOpts); err != nil {
		err = fmt.Errorf("err: %v", err)
		return err
	}

	if err := qflag.ParseAndRoute(); err != nil {
		err = fmt.Errorf("err: %v", err)
		return err
	}

	return nil
}

// 根命令运行函数
func run(cmd qflag.Command) error {
	qflag.Root.PrintHelp()
	return nil
}
