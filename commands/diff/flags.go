package diff

import (
	"flag"

	"gitee.com/MM-Q/fck/commands/internal/types"
	"gitee.com/MM-Q/qflag"
	"gitee.com/MM-Q/qflag/cmd"
)

var (
	// fck diff 子命令
	diffCmd     *cmd.Cmd
	diffCmdFile *qflag.StringFlag // file 标志
)

func InitDiffCmd() *cmd.Cmd {
	// fck diff 子命令
	diffCmd = qflag.NewCmd("diff", "d", flag.ExitOnError).
		WithUseChinese(true).
		WithDescription("文件校验工具, 对比指定目录A和目录B的文件差异, 并支持指定校验类型").
		WithNote("校验文件必须包含有效的头信息").
		WithNote("校验时会自动跳过空行和注释行(以#开头的行)")
	diffCmdFile = diffCmd.String("file", "f", types.OutputFileName, "指定用于校验的哈希值文件，程序将依据该文件中的哈希值进行校验操作")

	// 创建并返回一个命令对象
	return diffCmd
}
