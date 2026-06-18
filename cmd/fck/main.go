package main

import (
	"os"

	"gitee.com/MM-Q/color"

	"gitee.com/MM-Q/fck/internal/cli"
)

func main() {
	if err := cli.InitAndRun(); err != nil {
		_, _ = color.New(color.FgRed, color.Bold).Println(err)
		os.Exit(1)
	}
}
