package main

import (
	"fmt"
	"os"

	"gitee.com/MM-Q/fck/internal/cli"
)

func main() {
	if err := cli.InitAndRun(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
