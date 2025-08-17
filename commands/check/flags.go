// Package check 定义了文件校验命令的标志和参数配置。
// 该文件负责初始化 check 子命令的命令行参数解析和帮助信息设置。
package check

import (
	"flag"

	"gitee.com/MM-Q/qflag"
	"gitee.com/MM-Q/qflag/cmd"
)

var (
	// fck check 子命令
	checkCmd *cmd.Cmd
)

func InitCheckCmd() *cmd.Cmd {
	// fck check 子命令
	checkCmd = qflag.NewCmd("check", "c", flag.ExitOnError).
		WithUseChinese(true).
		WithDescription("文件校验工具, 对比指定目录A和目录B的文件差异, 并支持指定校验类型").
		WithNote("校验文件必须包含有效的头信息").
		WithNote("校验时会自动跳过空行和注释行(以#开头的行)")

	// 创建并返回一个命令对象
	return checkCmd
}
