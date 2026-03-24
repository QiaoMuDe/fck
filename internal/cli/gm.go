package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/gm"
	"gitee.com/MM-Q/qflag"
)

var GmCmd *qflag.Cmd

var (
	gmCheck   *qflag.BoolFlag
	gmVersion *qflag.BoolFlag
	gmHash    *qflag.BoolFlag
	gmTime    *qflag.BoolFlag
	gmStatus  *qflag.BoolFlag
	gmAbbrev  *qflag.IntFlag
	gmFormat  *qflag.StringFlag
	gmJSON    *qflag.BoolFlag
)

func init() {
	GmCmd = qflag.NewCmd("gm", "", qflag.ContinueOnError)

	gmCheck = GmCmd.Bool("check", "c", "检查当前目录是否为Git仓库", false)
	gmVersion = GmCmd.Bool("version", "v", "获取Git仓库的版本号 (标签或提交描述) ", false)
	gmHash = GmCmd.Bool("hash", "H", "获取当前提交的哈希值", false)
	gmTime = GmCmd.Bool("time", "t", "获取最新提交的时间", false)
	gmStatus = GmCmd.Bool("status", "s", "获取Git仓库的状态", false)
	gmAbbrev = GmCmd.Int("abbrev", "a", "哈希缩写长度 (配合--hash使用, 0表示完整哈希) ", 7)
	gmFormat = GmCmd.String("format", "f", "时间格式 (配合--time使用, Go时间格式) ", "2006-01-02 15:04:05")
	gmJSON = GmCmd.Bool("json", "j", "以JSON格式输出结果", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:       "获取Git仓库的各种元数据信息, 包括版本号、提交哈希、提交时间、仓库状态等",
		Notes:      []string{"不指定任何功能标志时, 默认获取所有元数据信息"},
		UseChinese: true,
	}

	if err := GmCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	GmCmd.SetRun(runGm)
}

func runGm(cmd qflag.Command) error {
	// 获取位置参数 (工作目录)
	args := cmd.Args()
	workDir := ""
	if len(args) > 0 {
		workDir = args[0]
	}

	config := gm.GmConfig{
		WorkDir: workDir,
		Check:   gmCheck.Get(),
		Version: gmVersion.Get(),
		Hash:    gmHash.Get(),
		Time:    gmTime.Get(),
		Status:  gmStatus.Get(),
		Abbrev:  gmAbbrev.Get(),
		Format:  gmFormat.Get(),
		JSON:    gmJSON.Get(),
	}

	return gm.GmCmdMain(config)
}
