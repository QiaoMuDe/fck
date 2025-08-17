// Package commands 实现了 fck 命令行工具的主要入口和子命令调度功能。
// 该包负责初始化各个子命令（size、list、check、hash、find），
// 解析命令行参数，并根据用户输入调度到相应的子命令执行器。
package commands

import (
	"fmt"
	"os"
	"runtime/debug"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/commands/check"
	"gitee.com/MM-Q/fck/commands/find"
	"gitee.com/MM-Q/fck/commands/hash"
	"gitee.com/MM-Q/fck/commands/list"
	"gitee.com/MM-Q/fck/commands/size"
	"gitee.com/MM-Q/qflag"
)

// Run 运行命令行工具
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
	checkCmd := check.InitCheckCmd()

	// 获取hashCmd子命令
	hashCmd := hash.InitHashCmd()

	// 获取findCmd子命令
	findCmd := find.InitFindCmd()

	// 添加子命令到全局根命令
	if addCmdErr := qflag.AddSubCmd(sizeCmd, listCmd, checkCmd, hashCmd, findCmd); addCmdErr != nil {
		fmt.Printf("err: %v\n", addCmdErr)
		os.Exit(1)
	}

	// 解析参数
	if parseErr := qflag.Parse(); parseErr != nil {
		fmt.Printf("err: %v\n", parseErr)
		os.Exit(1)
	}

	// 获取子命令名字
	subCmdName := qflag.Arg(0)

	// 如果没有指定子命令，则打印帮助信息
	if subCmdName == "" {
		// 打印帮助信息并退出
		qflag.PrintHelp()
		os.Exit(0)
	}

	// 执行子命令
	switch subCmdName {
	case hashCmd.LongName(), hashCmd.ShortName(): // hash 子命令
		// 执行 hash 子命令
		if err := hash.HashCmdMain(cmdCL); err != nil {
			fmt.Printf("err: %v\n", err)
			os.Exit(1)
		}
	case sizeCmd.LongName(), sizeCmd.ShortName(): // size 子命令
		// 执行 size 子命令
		if err := size.SizeCmdMain(cmdCL); err != nil {
			fmt.Printf("err: %v\n", err)
			os.Exit(1)
		}
	case checkCmd.LongName(), checkCmd.ShortName(): // check 子命令
		// 执行 check 子命令
		if err := check.CheckCmdMain(cmdCL); err != nil {
			fmt.Printf("err: %v\n", err)
			os.Exit(1)
		}
	case findCmd.LongName(), findCmd.ShortName(): // find 子命令
		// 执行 find 子命令
		if err := find.FindCmdMain(cmdCL); err != nil {
			fmt.Printf("err: %v\n", err)
			os.Exit(1)
		}
	case listCmd.LongName(), listCmd.ShortName(): // list 子命令
		// 执行 list 子命令
		if err := list.ListCmdMain(cmdCL); err != nil {
			fmt.Printf("err: %v\n", err)
			os.Exit(1)
		}
	default:
		// 如果是未知的子命令, 则打印帮助信息并退出
		fmt.Printf("err: 未知的子命令 %s\n", subCmdName)
		qflag.PrintHelp()
		os.Exit(1)
	}

	os.Exit(0)
}
