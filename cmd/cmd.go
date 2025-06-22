// Package cmd 包含了命令行工具相关的功能和函数
package cmd

import (
	"flag"
	"fmt"
	"os"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/globals"
)

func Run() {
	defer func() {
		if err := recover(); err != nil {
			// 打印错误信息并退出
			fmt.Printf("err: %v\n", err)
			os.Exit(1)
		}
	}()

	// 解析命令行参数
	flag.Parse()

	// 获取子命令专用cl
	cmdCl := colorlib.NewColorLib()

	// 执行子命令
	switch flag.Arg(0) {
	case hashCmd.LongName(), hashCmd.ShortName(): // hash 子命令
		// 解析 hash 子命令的参数
		if err := hashCmd.Parse(flag.Args()[1:]); err != nil {
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
		if err := sizeCmd.Parse(flag.Args()[1:]); err != nil {
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
		if err := diffCmd.Parse(flag.Args()[1:]); err != nil {
			fmt.Printf("err: %v\n", err)
			os.Exit(1)
		}

		// 执行 diff 子命令
		if err := diffCmdMain(cmdCl); err != nil {
			fmt.Printf("err: %v\n", err)
			os.Exit(1)
		}
	case "find":
		// 解析 find 子命令的参数
		if err := findCmd.Parse(flag.Args()[1:]); err != nil {
			fmt.Printf("err: %v\n", err)
			os.Exit(1)
		}

		// 如果是 -h 或 help, 则打印帮助信息并退出
		if *findCmdHelp {
			fmt.Println(globals.FindHelp)
			os.Exit(0)
		}

		// 执行 find 子命令
		if err := findCmdMain(cmdCl, findCmd); err != nil {
			fmt.Printf("err: %v\n", err)
			os.Exit(1)
		}
	case "f":
		// 解析 find 子命令的参数
		if err := findCmd.Parse(flag.Args()[1:]); err != nil {
			fmt.Printf("err: %v\n", err)
			os.Exit(1)
		}

		// 如果是 -h 或 help, 则打印帮助信息并退出
		if *findCmdHelp {
			fmt.Println(globals.FindHelp)
			os.Exit(0)
		}

		// 执行 find 子命令
		if err := findCmdMain(cmdCl, findCmd); err != nil {
			fmt.Printf("err: %v\n", err)
			os.Exit(1)
		}
	case "list":
		// 解析 list 子命令的参数
		if err := listCmd.Parse(flag.Args()[1:]); err != nil {
			fmt.Printf("err: %v\n", err)
			os.Exit(1)
		}
		// 如果是 -h 或 help, 则打印帮助信息并退出
		if *listCmdHelp {
			fmt.Println(globals.ListHelp)
			os.Exit(0)
		}
		// 执行 list 子命令
		if err := listCmdMain(cmdCl, listCmd); err != nil {
			fmt.Printf("err: %v\n", err)
			os.Exit(1)
		}
	case "ls":
		// 解析 list 子命令的参数
		if err := listCmd.Parse(flag.Args()[1:]); err != nil {
			fmt.Printf("err: %v\n", err)
			os.Exit(1)
		}
		// 如果是 -h 或 help, 则打印帮助信息并退出
		if *listCmdHelp {
			fmt.Println(globals.ListHelp)
			os.Exit(0)
		}
		// 执行 list 子命令
		if err := listCmdMain(cmdCl, listCmd); err != nil {
			fmt.Printf("err: %v\n", err)
			os.Exit(1)
		}
	default:
		// 如果是未知的子命令, 则打印帮助信息并退出
		fmt.Println(globals.FckHelp)
		os.Exit(0)
	}

	os.Exit(0)
}
