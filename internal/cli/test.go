package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/testcmd"
	"gitee.com/MM-Q/qflag"
)

var TestCmd *qflag.Cmd

var (
	testFile *qflag.BoolFlag // 检查指定文件是否存在
	testDir  *qflag.BoolFlag // 检查指定目录是否存在
	testLink *qflag.BoolFlag // 检查指定符号链接是否存在
)

func init() {
	TestCmd = qflag.NewCmd("test", "", qflag.ExitOnError)

	testFile = TestCmd.Bool("file", "f", "检查指定文件是否存在", false)
	testDir = TestCmd.Bool("dir", "d", "检查指定目录是否存在", false)
	testLink = TestCmd.Bool("link", "l", "检查指定符号链接是否存在", false)

	cmdOpts := &qflag.CmdOpts{
		Desc: "路径检测工具",
		Notes: []string{
			"路径通过位置参数传递",
			"不指定标志时，检测路径是否存在",
			"使用 -f 选项检测文件是否存在",
			"使用 -d 选项检测目录是否存在",
			"使用 -l 选项检测符号链接是否存在",
			"适合在脚本中使用，返回正确的退出码",
		},
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s test [options] <path>", qflag.Root.Name()),
		Examples: map[string]string{
			"检测文件是否存在":   fmt.Sprintf("%s test -f /path/to/file", qflag.Root.Name()),
			"检测目录是否存在":   fmt.Sprintf("%s test -d /path/to/dir", qflag.Root.Name()),
			"检测符号链接是否存在": fmt.Sprintf("%s test -l /path/to/link", qflag.Root.Name()),
		},
	}

	if err := TestCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	TestCmd.SetRun(runTest)
}

func runTest(cmd qflag.Command) error {
	args := cmd.Args()
	if len(args) == 0 {
		return fmt.Errorf("missing required path argument")
	}

	config := testcmd.TestConfig{
		CheckFile: testFile.Get(),
		CheckDir:  testDir.Get(),
		CheckLink: testLink.Get(),
		Path:      args[0],
	}

	return testcmd.TestCmdMain(config)
}
