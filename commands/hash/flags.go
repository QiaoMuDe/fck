// Package hash 定义了文件哈希计算命令的标志和参数配置。
// 该文件负责初始化 hash 子命令的命令行参数解析和帮助信息设置。
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
	hashCmdType      *qflag.EnumFlag   // type 标志
	hashCmdRecursion *qflag.BoolFlag   // recursion 标志
	hashCmdWrite     *qflag.BoolFlag   // write 标志
	hashCmdHidden    *qflag.BoolFlag   // hidden 标志
	hashCmdProgress  *qflag.BoolFlag   // progress 标志
	hashCmdLocal     *qflag.BoolFlag   // local 标志
	hashCmdBasePath  *qflag.StringFlag // base-path 标志
)

func InitHashCmd() *cmd.Cmd {
	// fck hash 子命令
	hashCmd = qflag.NewCmd("hash", "h", flag.ExitOnError).
		WithUsage(fmt.Sprint(qflag.LongName(), " hash [options] <path>\n")).
		WithChinese(true).
		WithDesc("文件哈希计算工具, 计算指定文件或目录的哈希值，支持多种哈希算法和并发处理")
	hashCmdType = hashCmd.Enum("type", "t", "md5", "指定哈希算法，支持 md5、sha1、sha256、sha512", []string{"md5", "sha1", "sha256", "sha512"})
	hashCmdRecursion = hashCmd.Bool("recursion", "r", false, "递归处理目录")
	hashCmdWrite = hashCmd.Bool("write", "w", false, "将哈希值写入文件, 文件名为checksum.hash")
	hashCmdHidden = hashCmd.Bool("hidden", "H", false, "启用计算隐藏文件/目录的哈希值，默认跳过")
	hashCmdProgress = hashCmd.Bool("progress", "p", false, "显示文件哈希计算进度条, 推荐在大文件处理时使用")
	hashCmdLocal = hashCmd.Bool("local", "l", false, "生成本地模式校验文件，记录绝对路径和基准目录")
	hashCmdBasePath = hashCmd.String("base-path", "b", "", "指定基准路径(默认为当前工作目录)")

	return hashCmd
}
