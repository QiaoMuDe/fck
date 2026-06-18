package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/shfmt"
	"gitee.com/MM-Q/qflag"
	"gitee.com/MM-Q/shx"
)

var ShfmtCmd *qflag.Cmd

var (
	shfmtWrite            *qflag.BoolFlag // -w, --write  原地写入
	shfmtBackup           *qflag.BoolFlag // -b, --backup  写入前备份
	shfmtList             *qflag.BoolFlag // -l, --list  列表模式
	shfmtIndent           *qflag.UintFlag // -i, --indent  缩进空格数
	shfmtSwitchCaseIndent *qflag.BoolFlag // --switch-case-indent  case 语句体缩进
	shfmtKeepComments     *qflag.BoolFlag // --keep-comments  保留注释
	shfmtBinaryNextLine   *qflag.BoolFlag // --binary-next-line  二元操作符换行
	shfmtFunctionNextLine *qflag.BoolFlag // --function-next-line  函数体换行
	shfmtSpaceRedirects   *qflag.BoolFlag // --space-redirects  重定向符空格
	shfmtSingleLine       *qflag.BoolFlag // --single-line  单行输出
	shfmtMinify           *qflag.BoolFlag // -m, --minify  最小化输出
)

func init() {
	ShfmtCmd = qflag.NewCmd("shfmt", "", qflag.ExitOnError)

	// 基础标志
	shfmtWrite = ShfmtCmd.Bool("write", "w", "原地写入格式化后的内容到原文件", false)
	shfmtBackup = ShfmtCmd.Bool("backup", "b", "写入前创建备份文件 (.bak后缀)", false)
	shfmtList = ShfmtCmd.Bool("list", "l", "列表模式: 仅列出需要格式化的文件，不修改文件", false)

	// 格式化选项标志（默认值从 shx 库获取，避免硬编码）
	def := shx.DefaultFormatOptions()
	shfmtIndent = ShfmtCmd.Uint("indent", "i", "缩进空格数 (0 表示使用 tab)", def.Indent)
	shfmtSwitchCaseIndent = ShfmtCmd.Bool("switch-case-indent", "", "case 语句体是否缩进", def.SwitchCaseIndent)
	shfmtKeepComments = ShfmtCmd.Bool("keep-comments", "", "是否保留注释", def.KeepComments)
	shfmtBinaryNextLine = ShfmtCmd.Bool("binary-next-line", "", "二元操作符（&&、||）是否换行显示", false)
	shfmtFunctionNextLine = ShfmtCmd.Bool("function-next-line", "", "函数体 { 是否换行", false)
	shfmtSpaceRedirects = ShfmtCmd.Bool("space-redirects", "", "重定向符前后是否加空格", false)
	shfmtSingleLine = ShfmtCmd.Bool("single-line", "", "是否单行输出", false)
	shfmtMinify = ShfmtCmd.Bool("minify", "m", "最小化输出 (压缩模式)", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "Shell 脚本格式化工具",
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s shfmt [options] [files...]", qflag.Root.Name()),
		Notes: []string{
			"输入方式: 指定文件路径（支持通配符）或通过管道传递脚本内容",
			"支持通配符模式，如: *.sh, **/*.sh",
		},
		Examples: map[string]string{
			"格式化单个文件":    fmt.Sprintf("%s shfmt script.sh", qflag.Root.Name()),
			"通配符批量格式化":   fmt.Sprintf("%s shfmt -w *.sh", qflag.Root.Name()),
			"递归格式化所有脚本":  fmt.Sprintf("%s shfmt -w **/*.sh", qflag.Root.Name()),
			"原地写入并备份":    fmt.Sprintf("%s shfmt -w -b script.sh", qflag.Root.Name()),
			"管道输入格式化":    fmt.Sprintf("cat script.sh | %s shfmt", qflag.Root.Name()),
			"使用 4 空格缩进":  fmt.Sprintf("%s shfmt -i 4 script.sh", qflag.Root.Name()),
			"使用 tab 缩进":  fmt.Sprintf("%s shfmt -i 0 script.sh", qflag.Root.Name()),
			"最小化输出并原地写入": fmt.Sprintf("%s shfmt -m -w script.sh", qflag.Root.Name()),
			"列表模式":       fmt.Sprintf("%s shfmt -l *.sh", qflag.Root.Name()),
		},
		MutexGroups: []qflag.MutexGroup{
			{
				Name:      "mode",
				Flags:     []string{"list", "write"},
				AllowNone: true,
			},
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
	// 构建格式化选项
	opts := shx.FormatOptions{
		Indent:           shfmtIndent.Get(),
		SwitchCaseIndent: shfmtSwitchCaseIndent.Get(),
		KeepComments:     shfmtKeepComments.Get(),
		BinaryNextLine:   shfmtBinaryNextLine.Get(),
		FunctionNextLine: shfmtFunctionNextLine.Get(),
		SpaceRedirects:   shfmtSpaceRedirects.Get(),
		SingleLine:       shfmtSingleLine.Get(),
		Minify:           shfmtMinify.Get(),
	}

	config := shfmt.ShfmtConfig{
		Write:         shfmtWrite.Get(),
		Backup:        shfmtBackup.Get(),
		List:          shfmtList.Get(),
		Files:         cmd.Args(),
		FormatOptions: opts,
	}

	return shfmt.ShfmtCmdMain(config)
}
