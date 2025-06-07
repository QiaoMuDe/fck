package main

import (
	"os"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/cmd"
)

// 主函数
func main() {
	// 使用线程安全方式获取CL
	cl := colorlib.GetCL()

	// 运行命令
	if runErr := cmd.Run(cl); runErr != nil {
		cl.PrintErr(runErr.Error())
		os.Exit(1)
	}
}
