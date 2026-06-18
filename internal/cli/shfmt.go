package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/shfmt"
	"gitee.com/MM-Q/qflag"
)

var ShfmtCmd *qflag.Cmd

var (
	shfmtWrite  *qflag.BoolFlag // -w, --write   原地写入
	shfmtBackup *qflag.BoolFlag // -b, --backup  写入前备份
)

func init() {
	ShfmtCmd = qflag.NewCmd("shfmt", "", qflag.ExitOnError)

	shfmtWrite = ShfmtCmd.Bool("write", "w", "原地写入格式化后的内容到原文件", false)
	shfmtBackup = ShfmtCmd.Bool("backup", "b", "写入前创建备份文件 (.bak后缀)", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "Shell 脚本格式化工具",
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s shfmt [options] [files...]", qflag.Root.Name()),
		Notes: []string{
			"输入方式: 指定文件路径（支持通配符）或通过管道传递脚本内容",
			"支持通配符模式，如: *.sh, **/*.sh",
			"使用 -w 原地写入时，格式化内容直接写回原文件",
			"使用 -b 备份时，会创建 .bak 后缀的备份文件",
		},
		Examples: map[string]string{
			"格式化单个文件":   fmt.Sprintf("%s shfmt script.sh", qflag.Root.Name()),
			"通配符批量格式化":  fmt.Sprintf("%s shfmt -w *.sh", qflag.Root.Name()),
			"递归格式化所有脚本": fmt.Sprintf("%s shfmt -w **/*.sh", qflag.Root.Name()),
			"原地写入并备份":   fmt.Sprintf("%s shfmt -w -b script.sh", qflag.Root.Name()),
			"管道输入格式化":   fmt.Sprintf("cat script.sh | %s shfmt", qflag.Root.Name()),
		},
		FlagDependencies: []qflag.FlagDependency{
			{
				Name:    "backup-write",
				Trigger: "backup",
				Targets: []string{"write"},
				Type:    qflag.DepRequired,
			},
		},
	}

	if err := ShfmtCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	ShfmtCmd.SetRun(runShfmt)
}

// runShfmt 执行 shfmt 命令
func runShfmt(cmd qflag.Command) error {
	config := shfmt.ShfmtConfig{
		Write:  shfmtWrite.Get(),
		Backup: shfmtBackup.Get(),
		Files:  cmd.Args(),
	}

	return shfmt.ShfmtCmdMain(config)
}
