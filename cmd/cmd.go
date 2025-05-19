package cmd

import (
	"flag"
	"fmt"
	"os"

	"gitee.com/MM-Q/colorlib"
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
		//
		return nil
	}

	// 执行子命令
	switch flag.Arg(0) {
	case "hash":
		hashCmd.Parse(flag.Args()[1:])
		// 如果是 -h 或 help，则打印帮助信息并退出
		if *hashCmdHelp {
			//
			return nil
		}
		// 执行 hash 子命令
		if err := hashCmdMain(hashCmd, cl); err != nil {
			// 打印错误信息并退出
			return fmt.Errorf("执行hash子命令时发生了错误: %v", err)
		}
	case "size":
		sizeCmd.Parse(flag.Args()[1:])
		// 如果是 -h 或 help，则打印帮助信息并退出
		if *sizeCmdHelp {
			//
			return nil
		}
		// 执行 size 子命令
		//
	case "check":
		checkCmd.Parse(flag.Args()[1:])
		// 如果是 -h 或 help，则打印帮助信息并退出
		if *checkCmdHelp {
			//
			return nil
		}
	case "find":
		findCmd.Parse(flag.Args()[1:])
		// 如果是 -h 或 help，则打印帮助信息并退出
		if *findCmdHelp {
			//
			return nil
		}
		// 执行 find 子命令
		//
	default:
		// 如果是未知的子命令，则打印帮助信息并退出
		//
		return nil
	}

	return nil
}
