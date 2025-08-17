// Package main 是 fck 多功能文件处理工具的程序入口点。
// 该文件包含应用程序的主函数，负责启动命令行工具并调用相应的子命令处理逻辑。
package main

import (
	"gitee.com/MM-Q/fck/commands"
)

// 主函数
func main() {
	// 运行命令
	commands.Run()
}
