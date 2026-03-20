package cli

import (
	"fmt"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/internal/commands/check"
	"gitee.com/MM-Q/qflag"
)

var CheckCmd *qflag.Cmd

var (
	checkFile    *qflag.StringFlag
	checkBaseDir *qflag.StringFlag
	checkQuiet   *qflag.BoolFlag
	checkColor   *qflag.BoolFlag
)

func init() {
	CheckCmd = qflag.NewCmd("check", "c", qflag.ExitOnError)

	checkFile = CheckCmd.String("file", "f", "指定校验文件路径(默认为checksum.hash)", "")
	checkBaseDir = CheckCmd.String("base-dir", "b", "手动指定校验基准目录(覆盖自动检测)", "")
	checkQuiet = CheckCmd.Bool("quiet", "q", "是否静默模式, 不输出校验通过的信息避免噪音", false)
	checkColor = CheckCmd.Bool("color", "c", "是否启用颜色输出", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:       "文件校验工具, 对比指定目录A和目录B的文件差异, 并支持指定校验类型",
		Notes:      []string{"校验文件必须包含有效的头信息", "校验时会自动跳过空行和注释行(以#开头的行)"},
		UseChinese: true,
	}

	if err := CheckCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	CheckCmd.SetRun(runCheck)
}

func runCheck(cmd qflag.Command) error {
	cl := colorlib.NewColorLib()

	config := check.CheckConfig{
		CheckFile: checkFile.Get(),
		BaseDir:   checkBaseDir.Get(),
		Quiet:     checkQuiet.Get(),
		Color:     checkColor.Get(),
	}

	return check.CheckCmdMain(cl, config)
}
