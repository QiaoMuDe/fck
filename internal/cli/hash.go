package cli

import (
	"fmt"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/internal/commands/hash"
	"gitee.com/MM-Q/qflag"
)

var HashCmd *qflag.Cmd

var (
	hashType      *qflag.EnumFlag
	hashRecursion *qflag.BoolFlag
	hashWrite     *qflag.BoolFlag
	hashHidden    *qflag.BoolFlag
	hashProgress  *qflag.BoolFlag
	hashLocal     *qflag.BoolFlag
	hashBasePath  *qflag.StringFlag
	noColor       *qflag.BoolFlag
)

func init() {
	HashCmd = qflag.NewCmd("hash", "h", qflag.ExitOnError)

	hashType = HashCmd.Enum("type", "t", "指定哈希算法，支持 md5、sha1、sha256、sha512", "md5", []string{"md5", "sha1", "sha256", "sha512"})
	hashRecursion = HashCmd.Bool("recursion", "r", "递归处理目录", false)
	hashWrite = HashCmd.Bool("write", "w", "将哈希值写入文件, 文件名为checksum.hash", false)
	hashHidden = HashCmd.Bool("hidden", "H", "启用计算隐藏文件/目录的哈希值，默认跳过", false)
	hashProgress = HashCmd.Bool("progress", "p", "显示文件哈希计算进度条, 推荐在大文件处理时使用", false)
	hashLocal = HashCmd.Bool("local", "l", "生成本地模式校验文件，记录绝对路径和基准目录", false)
	hashBasePath = HashCmd.String("base-path", "b", "指定基准路径(默认为当前工作目录)", "")
	noColor = HashCmd.Bool("no-color", "n", "禁用颜色输出", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:       "文件哈希计算工具, 计算指定文件或目录的哈希值，支持多种哈希算法和并发处理",
		Notes:      []string{"哈希值计算基于文件内容，不包括元数据"},
		UseChinese: true,
	}

	if err := HashCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	HashCmd.SetRun(runHash)
}

func runHash(cmd qflag.Command) error {
	cl := colorlib.NewColorLib()

	// 禁用颜色
	if noColor.Get() {
		cl.SetColor(!noColor.Get())
	}

	config := hash.HashConfig{
		TargetPaths: cmd.Args(),
		Type:        hashType.Get(),
		Recursion:   hashRecursion.Get(),
		Write:       hashWrite.Get(),
		Hidden:      hashHidden.Get(),
		Progress:    hashProgress.Get(),
		Local:       hashLocal.Get(),
		BasePath:    hashBasePath.Get(),
	}

	return hash.HashCmdMain(cl, config)
}
