// Package cli 提供命令行接口
package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/newline"
	"gitee.com/MM-Q/fck/internal/types"
	"gitee.com/MM-Q/qflag"
)

var NewlineCmd *qflag.Cmd

var (
	newlineTo     *qflag.EnumFlag   // -t 目标换行符格式
	newlineWrite  *qflag.BoolFlag   // -w 原地写入
	newlineBackup *qflag.BoolFlag   // -b 创建备份
	newlineOutput *qflag.StringFlag // -o 输出路径
	newlineList   *qflag.BoolFlag   // -l 列出支持的换行符类型
	newlineQuiet  *qflag.BoolFlag   // -q 静默模式
	newlineForce  *qflag.BoolFlag   // -f 强制转换
)

func init() {
	NewlineCmd = qflag.NewCmd("newline", "", qflag.ExitOnError)

	// 使用 types 包中定义的换行符类型列表
	newlineTo = NewlineCmd.Enum("to", "t", "目标换行符格式", types.NewlineNone, types.TargetNewlines)
	newlineWrite = NewlineCmd.Bool("write", "w", "原地写入 (覆盖原文件)", false)
	newlineBackup = NewlineCmd.Bool("backup", "b", "原地写入时创建备份 (.bak)", false)
	newlineOutput = NewlineCmd.String("output", "o", "输出目录或文件路径 (默认生成 .unix/.win/.mac 后缀文件)", ".")
	newlineList = NewlineCmd.Bool("list", "l", "列出支持的换行符类型", false)
	newlineQuiet = NewlineCmd.Bool("quiet", "q", "静默模式", false)
	newlineForce = NewlineCmd.Bool("force", "f", "强制转换 (即使已经是目标格式)", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "文件换行符检测与转换工具",
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s newline [options] <file...>", qflag.Root.Name()),
		Examples: map[string]string{
			"检测文件换行符":        fmt.Sprintf("%s newline file.txt", qflag.Root.Name()),
			"检测多个文件":         fmt.Sprintf("%s newline file1.txt file2.txt", qflag.Root.Name()),
			"转换为 Unix 格式":    fmt.Sprintf("%s newline -t lf file.txt", qflag.Root.Name()),
			"转换为 Windows 格式": fmt.Sprintf("%s newline -t crlf file.txt", qflag.Root.Name()),
			"原地转换":           fmt.Sprintf("%s newline -t lf -w file.txt", qflag.Root.Name()),
			"原地转换并备份":        fmt.Sprintf("%s newline -t lf -w -b file.txt", qflag.Root.Name()),
			"批量转换":           fmt.Sprintf("%s newline -t lf *.txt", qflag.Root.Name()),
			"列出支持的类型":        fmt.Sprintf("%s newline -l", qflag.Root.Name()),
		},
		Notes: []string{
			"-t 默认为 none (只检测不转换), 指定目标格式时执行转换",
			"使用 -w 原地写入会直接覆盖原文件",
			"建议重要文件使用 -b 创建备份",
			"Mixed 表示文件中存在多种换行符",
		},
		MutexGroups: []qflag.MutexGroup{
			{
				Name:      "to-list",
				Flags:     []string{"list", "to"},
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
			{
				Name:    "write-to",
				Trigger: "write",
				Targets: []string{"to"},
				Type:    qflag.DepRequired,
			},
		},
	}

	if err := NewlineCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	NewlineCmd.SetRun(runNewline)
}

// runNewline 执行 newline 命令
func runNewline(cmd qflag.Command) error {
	// 列出支持的换行符类型
	if newlineList.Get() {
		return newline.ListTypes()
	}

	// 检查文件参数
	files := cmd.Args()
	if len(files) == 0 {
		return fmt.Errorf("no input files specified")
	}

	config := newline.Config{
		Files:     files,
		ToNewline: newlineTo.Get(),
		Write:     newlineWrite.Get(),
		Backup:    newlineBackup.Get(),
		Output:    newlineOutput.Get(),
		Quiet:     newlineQuiet.Get(),
		Force:     newlineForce.Get(),
	}

	return newline.CmdMain(config)
}
