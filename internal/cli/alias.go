package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/alias"
	"gitee.com/MM-Q/qflag"
)

var AliasCmd *qflag.Cmd

var (
	aliasType *qflag.EnumFlag // -t, --type  Shell 类型
)

func init() {
	AliasCmd = qflag.NewCmd("alias", "", qflag.ExitOnError)

	aliasType = AliasCmd.Enum("type", "t", "Shell 类型 (bash/pwsh)", "none", []string{"bash", "pwsh", "none"})

	cmdOpts := &qflag.CmdOpts{
		Desc:        "生成 FCK 子命令的 Shell 别名定义",
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s alias [options]", qflag.Root.Name()),
		Examples: map[string]string{
			"生成 Bash 别名":       fmt.Sprintf("%s alias --type bash", qflag.Root.Name()),
			"生成 PowerShell 别名": fmt.Sprintf("%s alias --type pwsh", qflag.Root.Name()),
		},
		Notes: []string{
			"Bash 用户: 将输出添加到 ~/.bashrc 或 ~/.bash_aliases",
			"PowerShell 用户: 将输出添加到 $PROFILE 文件",
		},
	}

	if err := AliasCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	AliasCmd.SetRun(runAlias)
}

func runAlias(cmd qflag.Command) error {
	shellType := aliasType.Get()

	// 未指定类型时打印帮助信息
	if shellType == "none" {
		AliasCmd.PrintHelp()
		return nil
	}

	config := alias.AliasConfig{
		Type: shellType,
	}

	return alias.AliasCmdMain(config)
}
