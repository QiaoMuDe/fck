package hash

import (
	"flag"
	"fmt"

	"gitee.com/MM-Q/qflag"
	"gitee.com/MM-Q/qflag/cmd"
)

var (
	// fck hash 子命令
	hashCmd          *cmd.Cmd
	hashCmdType      *qflag.EnumFlag // type 标志
	hashCmdRecursion *qflag.BoolFlag // recursion 标志
	hashCmdWrite     *qflag.BoolFlag // write 标志
	hashCmdHidden    *qflag.BoolFlag // hidden 标志
	hashCmdProgress  *qflag.BoolFlag // progress 标志
)

func InitHashCmd() *cmd.Cmd {
	// fck hash 子命令
	hashCmd = qflag.NewCmd("hash", "h", flag.ExitOnError).
		WithUsageSyntax(fmt.Sprint(qflag.LongName(), " hash [options] <path>\n")).
		WithUseChinese(true).
		WithDescription("文件哈希计算工具, 计算指定文件或目录的哈希值，支持多种哈希算法和并发处理")
	hashCmdType = hashCmd.Enum("type", "t", "md5", "指定哈希算法，支持 md5、sha1、sha256、sha512", []string{"md5", "sha1", "sha256", "sha512"})
	hashCmdRecursion = hashCmd.Bool("recursion", "r", false, "递归处理目录")
	hashCmdWrite = hashCmd.Bool("write", "w", false, "将哈希值写入文件, 文件名为checksum.hash")
	hashCmdHidden = hashCmd.Bool("hidden", "H", false, "启用计算隐藏文件/目录的哈希值，默认跳过")
	hashCmdProgress = hashCmd.Bool("progress", "p", false, "显示文件哈希计算进度条, 推荐在大文件处理时使用")

	return hashCmd
}
