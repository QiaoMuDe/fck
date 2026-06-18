package main

import (
	"os"

	"gitee.com/MM-Q/color"

	"gitee.com/MM-Q/fck/internal/cli"
	"gitee.com/MM-Q/shx"
)

func main() {
	if err := cli.InitAndRun(); err != nil {
		_, _ = color.New(color.FgRed, color.Bold).Fprintln(os.Stderr, err)

		// shx 子命令的退出码透传
		if code, ok := shx.IsExitStatus(err); ok {
			os.Exit(int(code))
		}

		// 其他错误
		os.Exit(1)
	}
}
