package extract

import (
	"flag"
	"fmt"

	"gitee.com/MM-Q/qflag"
	"gitee.com/MM-Q/qflag/cmd"
)

var (
	extractCmd *cmd.Cmd
)

func InitExtractCmd() *cmd.Cmd {
	extractCmd = qflag.NewCmd("extract", "", flag.ExitOnError).
		WithUseChinese(true).
		WithUsageSyntax(fmt.Sprint(qflag.LongName(), " extract [options] <archive> [destination]")).
		WithDescription("智能解压缩工具, 智能识别压缩文件格式并解压")

	return extractCmd
}
