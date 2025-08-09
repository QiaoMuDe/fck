// package commands 包含了命令行工具相关的功能和函数
package commands

import (
	"fmt"
	"os"
	"runtime/debug"

	"gitee.com/MM-Q/colorlib"
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

	// 获取子命令专用cl
	cmdCl := colorlib.NewColorLib()

	// 执行子命令
	switch qflag.Arg(0) {
	case hashCmd.LongName(), hashCmd.ShortName(): // hash 子命令
		// 解析 hash 子命令的参数
		if err := hashCmd.Parse(qflag.Args()[1:]); err != nil {
			fmt.Printf("err: %v\n", err)
			os.Exit(1)
		}

		// 执行 hash 子命令
		if err := hashCmdMain(cmdCl); err != nil {
			fmt.Printf("err: %v\n", err)
			os.Exit(1)
		}
	case sizeCmd.LongName(), sizeCmd.ShortName(): // size 子命令
		// 解析 size 子命令的参数
		if err := sizeCmd.Parse(qflag.Args()[1:]); err != nil {
			fmt.Printf("err: %v\n", err)
			os.Exit(1)
		}

		// 执行 size 子命令
		if err := sizeCmdMain(cmdCl); err != nil {
			fmt.Printf("err: %v\n", err)
			os.Exit(1)
		}
	case diffCmd.LongName(), diffCmd.ShortName(): // diff 子命令
		// 解析 diff 子命令的参数
		if err := diffCmd.Parse(qflag.Args()[1:]); err != nil {
			fmt.Printf("err: %v\n", err)
			os.Exit(1)
		}

		// 执行 diff 子命令
		if err := diffCmdMain(cmdCl); err != nil {
			fmt.Printf("err: %v\n", err)
			os.Exit(1)
		}
	case findCmd.LongName(), findCmd.ShortName(): // find 子命令
		// 解析 find 子命令的参数
		if err := findCmd.Parse(qflag.Args()[1:]); err != nil {
			fmt.Printf("err: %v\n", err)
			os.Exit(1)
		}

		// 执行 find 子命令
		if err := findCmdMain(cmdCl); err != nil {
			fmt.Printf("err: %v\n", err)
			os.Exit(1)
		}
	case listCmd.LongName(), listCmd.ShortName(): // list 子命令
		// 解析 list 子命令的参数
		if err := listCmd.Parse(qflag.Args()[1:]); err != nil {
			fmt.Printf("err: %v\n", err)
			os.Exit(1)
		}
		// 执行 list 子命令
		if err := listCmdMain(cmdCl); err != nil {
			fmt.Printf("err: %v\n", err)
			os.Exit(1)
		}
	default:
		// 如果是未知的子命令, 则打印帮助信息并退出
		fmt.Printf("err: 未知的子命令 %s\n", qflag.Arg(0))
		qflag.PrintHelp()
		os.Exit(0)
	}

	os.Exit(0)
}
