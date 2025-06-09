// Package cmd 包含了命令行工具相关的功能和函数
package cmd

import (
	"flag"
	"fmt"
	"os"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/globals"
	"gitee.com/MM-Q/verman"
)

func Run() {
	defer func() {
		if err := recover(); err != nil {
			// 打印错误信息并退出
			fmt.Printf("fck在运行过程中发生了错误: %v\n", err)
			os.Exit(1)
		}
	}()

	// 解析命令行参数
	flag.Parse()

	// 如果是 -v 或 version, 则打印版本信息并退出
	if *versionF || flag.Arg(0) == "version" {
		// 打印版本信息并退出
		version := verman.Get()
		fmt.Printf("%s %s\n", version.AppName, version.GitVersion)
		os.Exit(0)
	}

	// 如果是 -h 或 没有参数或者第一个参数是 help, 则打印帮助信息并退出
	if *helpF || flag.NArg() == 0 || flag.Arg(0) == "help" {
		fmt.Println(globals.FckHelp)
		os.Exit(0)
	}

	// 获取子命令专用cl
	cmdCl := colorlib.NewColorLib()

	// 执行子命令
	switch flag.Arg(0) {
	case "hash":
		// 解析 hash 子命令的参数
		if err := hashCmd.Parse(flag.Args()[1:]); err != nil {
			fmt.Printf("解析hash子命令的参数时发生了错误: %v\n", err)
			os.Exit(1)
		}

		// 如果是 -h 或 help, 则打印帮助信息并退出
		if *hashCmdHelp {
			fmt.Println(globals.HashHelp)
			os.Exit(0)
		}

		// 执行 hash 子命令
		if err := hashCmdMain(hashCmd, cmdCl); err != nil {
			fmt.Printf("执行hash子命令时发生了错误: %v\n", err)
			os.Exit(1)
		}
	case "h":
		// 解析 hash 子命令的参数
		if err := hashCmd.Parse(flag.Args()[1:]); err != nil {
			fmt.Printf("解析hash子命令的参数时发生了错误: %v\n", err)
			os.Exit(1)
		}

		// 如果是 -h 或 help, 则打印帮助信息并退出
		if *hashCmdHelp {
			fmt.Println(globals.HashHelp)
			os.Exit(0)
		}

		// 执行 hash 子命令
		if err := hashCmdMain(hashCmd, cmdCl); err != nil {
			fmt.Printf("执行hash子命令时发生了错误: %v\n", err)
			os.Exit(1)
		}
	case "size":
		// 解析 size 子命令的参数
		if err := sizeCmd.Parse(flag.Args()[1:]); err != nil {
			fmt.Printf("解析size子命令的参数时发生了错误: %v\n", err)
			os.Exit(1)
		}

		// 如果是 -h 或 help, 则打印帮助信息并退出
		if *sizeCmdHelp {
			fmt.Println(globals.SizeHelp)
			os.Exit(0)
		}

		// 执行 size 子命令
		if err := sizeCmdMain(sizeCmd, cmdCl); err != nil {
			fmt.Printf("执行size子命令时发生了错误: %v\n", err)
			os.Exit(1)
		}
	case "s":
		// 解析 size 子命令的参数
		if err := sizeCmd.Parse(flag.Args()[1:]); err != nil {
			fmt.Printf("解析size子命令的参数时发生了错误: %v\n", err)
			os.Exit(1)
		}

		// 如果是 -h 或 help, 则打印帮助信息并退出
		if *sizeCmdHelp {
			fmt.Println(globals.SizeHelp)
			os.Exit(0)
		}

		// 执行 size 子命令
		if err := sizeCmdMain(sizeCmd, cmdCl); err != nil {
			fmt.Printf("执行size子命令时发生了错误: %v\n", err)
			os.Exit(1)
		}
	case "diff":
		// 解析 diff 子命令的参数
		if err := diffCmd.Parse(flag.Args()[1:]); err != nil {
			fmt.Printf("解析diff子命令的参数时发生了错误: %v\n", err)
			os.Exit(1)
		}

		// 如果是 -h 或 help, 则打印帮助信息并退出
		if *diffCmdHelp {
			fmt.Println(globals.DiffHelp)
			os.Exit(0)
		}

		// 执行 diff 子命令
		if err := diffCmdMain(cmdCl); err != nil {
			fmt.Printf("执行diff子命令时发生了错误: %v\n", err)
			os.Exit(1)
		}
	case "d":
		// 解析 diff 子命令的参数
		if err := diffCmd.Parse(flag.Args()[1:]); err != nil {
			fmt.Printf("解析diff子命令的参数时发生了错误: %v\n", err)
			os.Exit(1)
		}

		// 如果是 -h 或 help, 则打印帮助信息并退出
		if *diffCmdHelp {
			fmt.Println(globals.DiffHelp)
			os.Exit(0)
		}

		// 执行 diff 子命令
		if err := diffCmdMain(cmdCl); err != nil {
			fmt.Printf("执行diff子命令时发生了错误: %v\n", err)
			os.Exit(1)
		}
	case "find":
		// 解析 find 子命令的参数
		if err := findCmd.Parse(flag.Args()[1:]); err != nil {
			fmt.Printf("解析find子命令的参数时发生了错误: %v\n", err)
			os.Exit(1)
		}

		// 如果是 -h 或 help, 则打印帮助信息并退出
		if *findCmdHelp {
			fmt.Println(globals.FindHelp)
			os.Exit(0)
		}

		// 执行 find 子命令
		if err := findCmdMain(cmdCl, findCmd); err != nil {
			fmt.Printf("执行find子命令时发生了错误: %v\n", err)
			os.Exit(1)
		}
	case "f":
		// 解析 find 子命令的参数
		if err := findCmd.Parse(flag.Args()[1:]); err != nil {
			fmt.Printf("解析find子命令的参数时发生了错误: %v\n", err)
			os.Exit(1)
		}

		// 如果是 -h 或 help, 则打印帮助信息并退出
		if *findCmdHelp {
			fmt.Println(globals.FindHelp)
			os.Exit(0)
		}

		// 执行 find 子命令
		if err := findCmdMain(cmdCl, findCmd); err != nil {
			fmt.Printf("执行find子命令时发生了错误: %v\n", err)
			os.Exit(1)
		}
	case "list":
		// 解析 list 子命令的参数
		if err := listCmd.Parse(flag.Args()[1:]); err != nil {
			fmt.Printf("解析list子命令的参数时发生了错误: %v\n", err)
			os.Exit(1)
		}
		// 如果是 -h 或 help, 则打印帮助信息并退出
		if *listCmdHelp {
			fmt.Println(globals.ListHelp)
			os.Exit(0)
		}
		// 执行 list 子命令
		if err := listCmdMain(cmdCl, listCmd); err != nil {
			fmt.Printf("执行list子命令时发生了错误: %v\n", err)
			os.Exit(1)
		}
	case "ls":
		// 解析 list 子命令的参数
		if err := listCmd.Parse(flag.Args()[1:]); err != nil {
			fmt.Printf("解析list子命令的参数时发生了错误: %v\n", err)
			os.Exit(1)
		}
		// 如果是 -h 或 help, 则打印帮助信息并退出
		if *listCmdHelp {
			fmt.Println(globals.ListHelp)
			os.Exit(0)
		}
		// 执行 list 子命令
		if err := listCmdMain(cmdCl, listCmd); err != nil {
			fmt.Printf("执行list子命令时发生了错误: %v\n", err)
			os.Exit(1)
		}
	default:
		// 如果是未知的子命令, 则打印帮助信息并退出
		fmt.Println(globals.FckHelp)
		os.Exit(0)
	}

	os.Exit(0)
}
