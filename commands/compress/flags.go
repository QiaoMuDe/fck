package compress

import (
	"flag"
	"fmt"

	"gitee.com/MM-Q/qflag"
	"gitee.com/MM-Q/qflag/cmd"
)

var (
	compressCmd *cmd.Cmd
)

func InitCompressCmd() *cmd.Cmd {
	compressCmd = qflag.NewCmd("compress", "", flag.ExitOnError).
		WithUseChinese(true).
		WithUsageSyntax(fmt.Sprint(qflag.LongName(), " compress [options] <file> [file...]")).
		WithDescription("智能压缩打包工具, 智能识别文件类型并压缩打包")

	return compressCmd
}
