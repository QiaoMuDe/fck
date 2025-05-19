package main

import (
	"os"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/cmd"
)

// 主函数
func main() {
	// 初始化颜色库
	cl := colorlib.NewColorLib()
	
	// 运行命令
	if runErr := cmd.Run(cl) ; runErr != nil {
		cl.PrintError(runErr.Error())
		os.Exit(1)
	}
}
