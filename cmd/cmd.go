package cmd

import (
	"flag"
	"fmt"
	"os"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/globals"
	"gitee.com/MM-Q/verman"
)

func Run(cl *colorlib.ColorLib) error {
	defer func() {
		if err := recover(); err != nil {
			// 打印错误信息并退出
			cl.PrintErrf("fck在运行过程中发生了错误: %v\n", err)
			os.Exit(1)
		}
	}()

	// 解析命令行参数
	flag.Parse()

	// 如果是 -v 或 version，则打印版本信息并退出
	if *versionF || flag.Arg(0) == "version" {
		// 打印版本信息并退出
		version := verman.Get()
		cl.Greenf("%s %s", version.AppName, version.GitVersion)
		return nil
	}

	// 如果是 -h 或 没有参数或者第一个参数是 help，则打印帮助信息并退出
	if *helpF || flag.NArg() == 0 || flag.Arg(0) == "help" {
		fmt.Println(globals.FckHelp)
		return nil
	}

	// 执行子命令
	switch flag.Arg(0) {
	case "hash":
		hashCmd.Parse(flag.Args()[1:])
		// 如果是 -h 或 help，则打印帮助信息并退出
		if *hashCmdHelp {
			fmt.Println(globals.HashHelp)
			return nil
		}
		// 执行 hash 子命令
		if err := hashCmdMain(hashCmd, cl); err != nil {
			return fmt.Errorf("执行hash子命令时发生了错误: %v", err)
		}
	case "h":
		hashCmd.Parse(flag.Args()[1:])
		// 如果是 -h 或 help，则打印帮助信息并退出
		if *hashCmdHelp {
			fmt.Println(globals.HashHelp)
			return nil
		}
		// 执行 hash 子命令
		if err := hashCmdMain(hashCmd, cl); err != nil {
			return fmt.Errorf("执行hash子命令时发生了错误: %v", err)
		}
	case "size":
		sizeCmd.Parse(flag.Args()[1:])
		// 如果是 -h 或 help，则打印帮助信息并退出
		if *sizeCmdHelp {
			fmt.Println(globals.SizeHelp)
			return nil
		}
		// 执行 size 子命令
		if err := sizeCmdMain(sizeCmd, cl); err != nil {
			return fmt.Errorf("执行size子命令时发生了错误: %v", err)
		}
	case "s":
		sizeCmd.Parse(flag.Args()[1:])
		// 如果是 -h 或 help，则打印帮助信息并退出
		if *sizeCmdHelp {
			fmt.Println(globals.SizeHelp)
			return nil
		}
		// 执行 size 子命令
		if err := sizeCmdMain(sizeCmd, cl); err != nil {
			return fmt.Errorf("执行size子命令时发生了错误: %v", err)
		}
	case "check":
		checkCmd.Parse(flag.Args()[1:])
		// 如果是 -h 或 help，则打印帮助信息并退出
		if *checkCmdHelp {
			fmt.Println(globals.CheckHelp)
			return nil
		}

		// 执行 check 子命令
		if err := checkCmdMain(cl); err != nil {
			return fmt.Errorf("执行check子命令时发生了错误: %v", err)
		}
	case "c":
		checkCmd.Parse(flag.Args()[1:])
		// 如果是 -h 或 help，则打印帮助信息并退出
		if *checkCmdHelp {
			fmt.Println(globals.CheckHelp)
			return nil
		}

		// 执行 check 子命令
		if err := checkCmdMain(cl); err != nil {
			return fmt.Errorf("执行check子命令时发生了错误: %v", err)
		}
	case "find":
		findCmd.Parse(flag.Args()[1:])
		// 如果是 -h 或 help，则打印帮助信息并退出
		if *findCmdHelp {
			fmt.Println(globals.FindHelp)
			return nil
		}
		// 执行 find 子命令
		if err := findCmdMain(); err != nil {
			return fmt.Errorf("执行find子命令时发生了错误: %v", err)
		}
	case "f":
		findCmd.Parse(flag.Args()[1:])
		// 如果是 -h 或 help，则打印帮助信息并退出
		if *findCmdHelp {
			fmt.Println(globals.FindHelp)
			return nil
		}
		// 执行 find 子命令
		if err := findCmdMain(); err != nil {
			return fmt.Errorf("执行find子命令时发生了错误: %v", err)
		}
	default:
		// 如果是未知的子命令，则打印帮助信息并退出
		//
		return nil
	}

	return nil
}
