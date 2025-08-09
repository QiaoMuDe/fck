// package commands 包含了命令行工具相关的功能和函数
package commands

import (
	"fmt"
	"os"
	"runtime/debug"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/commands/diff"
	"gitee.com/MM-Q/fck/commands/find"
	"gitee.com/MM-Q/fck/commands/hash"
	"gitee.com/MM-Q/fck/commands/list"
	"gitee.com/MM-Q/fck/commands/size"
	"gitee.com/MM-Q/qflag"
)

func Run() {
	defer func() {
		if err := recover(); err != nil {
			// 打印错误信息并退出
			fmt.Printf("err: %v\nstack: %s\n", err, debug.Stack())
			os.Exit(1)
		}
	}()

	// 初始化主命令
	InItMainCmd()

	// 获取子命令专用CL
	cmdCL := colorlib.NewColorLib()

	// 获取sizeCmd子命令
	sizeCmd := size.InitSizeCmd()

	// 获取listCmd子命令
	listCmd := list.InitListCmd()

	// 获取diffCmd子命令
	diffCmd := diff.InitDiffCmd()

	// 获取hashCmd子命令
	hashCmd := hash.InitHashCmd()

	// 获取findCmd子命令
	findCmd := find.InitFindCmd()

	// 添加子命令到全局根命令
	if addCmdErr := qflag.AddSubCmd(sizeCmd, listCmd, diffCmd, hashCmd, findCmd); addCmdErr != nil {
		fmt.Printf("err: %v\n", addCmdErr)
		os.Exit(1)
	}

	// 解析参数
	if parseErr := qflag.Parse(); parseErr != nil {
		fmt.Printf("err: %v\n", parseErr)
		os.Exit(1)
	}

	// 执行子命令
	switch qflag.Arg(0) {
	case hashCmd.LongName(), hashCmd.ShortName(): // hash 子命令
		// 执行 hash 子命令
		if err := hash.HashCmdMain(cmdCL); err != nil {
			panic(err)
		}
	case sizeCmd.LongName(), sizeCmd.ShortName(): // size 子命令
		// 执行 size 子命令
		if err := size.SizeCmdMain(cmdCL); err != nil {
			panic(err)
		}
	case diffCmd.LongName(), diffCmd.ShortName(): // diff 子命令
		// 执行 diff 子命令
		if err := diff.DiffCmdMain(cmdCL); err != nil {
			panic(err)
		}
	case findCmd.LongName(), findCmd.ShortName(): // find 子命令
		// 执行 find 子命令
		if err := find.FindCmdMain(cmdCL); err != nil {
			panic(err)
		}
	case listCmd.LongName(), listCmd.ShortName(): // list 子命令
		// 执行 list 子命令
		if err := list.ListCmdMain(cmdCL); err != nil {
			panic(err)
		}
	default:
		// 如果是未知的子命令, 则打印帮助信息并退出
		fmt.Printf("err: 未知的子命令 %s\n", qflag.Arg(0))
		qflag.PrintHelp()
	}

	os.Exit(0)
}
