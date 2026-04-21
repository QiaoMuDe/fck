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
	checkVerbose *qflag.BoolFlag
	checkNoColor *qflag.BoolFlag
)

func init() {
	CheckCmd = qflag.NewCmd("check", "", qflag.ExitOnError)

	checkFile = CheckCmd.String("file", "f", "指定校验文件路径(默认为checksum.hash)", "")
	checkBaseDir = CheckCmd.String("base-dir", "b", "手动指定校验基准目录(覆盖自动检测)", "")
	checkVerbose = CheckCmd.Bool("verbose", "v", "显示详细信息，包括校验通过的文件", false)
	checkNoColor = CheckCmd.Bool("no-color", "n", "是否禁用颜色输出", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "校验指定文件的哈希值是否与校验文件中的记录一致",
		UsageSyntax: fmt.Sprintf("%s check [options]", qflag.Root.Name()),
		Notes:       []string{"校验文件必须包含有效的头信息", "校验时会自动跳过空行和注释行(以#开头的行)"},
		UseChinese:  true,
		Examples: map[string]string{
			"校验当前目录": fmt.Sprintf("%s check", qflag.Root.Name()),
			"指定校验文件": fmt.Sprintf("%s check -f checksum.hash", qflag.Root.Name()),
		},
	}

	if err := CheckCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	CheckCmd.SetRun(runCheck)
}

func runCheck(cmd qflag.Command) error {
	cl := colorlib.NewColorLib()
	if checkNoColor.Get() {
		cl.SetColor(!checkNoColor.Get())
	}

	config := check.CheckConfig{
		CheckFile: checkFile.Get(),
		BaseDir:   checkBaseDir.Get(),
		Verbose:   checkVerbose.Get(),
	}

	return check.CheckCmdMain(cl, config)
}
