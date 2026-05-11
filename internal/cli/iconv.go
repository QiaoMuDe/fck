// Package cli 提供命令行接口
package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/iconv"
	"gitee.com/MM-Q/fck/internal/types"
	"gitee.com/MM-Q/qflag"
)

var IconvCmd *qflag.Cmd

var (
	iconvFrom   *qflag.EnumFlag   // -f 源编码
	iconvTo     *qflag.EnumFlag   // -t 目标编码
	iconvWrite  *qflag.BoolFlag   // -w 原地写入
	iconvBackup *qflag.BoolFlag   // -b 创建备份
	iconvOutput *qflag.StringFlag // -o 输出路径
	iconvList   *qflag.BoolFlag   // -l 列出支持的编码
	iconvQuiet  *qflag.BoolFlag   // -q 静默模式
	iconvShort  *qflag.BoolFlag   // -s 短输出模式（只输出编码，用于脚本）
)

func init() {
	IconvCmd = qflag.NewCmd("iconv", "", qflag.ExitOnError)

	// 使用 types 包中定义的编码列表
	iconvFrom = IconvCmd.Enum("from", "f", "源文件编码", types.EncodingAuto, types.SupportedEncodings)
	iconvTo = IconvCmd.Enum("to", "t", "目标编码", types.EncodingNone, types.TargetEncodings)
	iconvWrite = IconvCmd.Bool("write", "w", "原地写入 (覆盖原文件) ", false)
	iconvBackup = IconvCmd.Bool("backup", "b", "原地写入时创建备份 (.bak) ", false)
	iconvOutput = IconvCmd.String("output", "o", "输出路径 (默认生成 .converted 后缀文件) ", ".")
	iconvList = IconvCmd.Bool("list", "l", "列出支持的编码格式", false)
	iconvQuiet = IconvCmd.Bool("quiet", "q", "静默模式", false)
	iconvShort = IconvCmd.Bool("short", "s", "短输出模式 (只输出编码, 用于脚本)", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "文件编码转换工具",
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s iconv [options] <file...>", qflag.Root.Name()),
		Examples: map[string]string{
			"检测文件编码":   fmt.Sprintf("%s iconv file.txt", qflag.Root.Name()),
			"检测多个文件编码": fmt.Sprintf("%s iconv file1.txt file2.txt", qflag.Root.Name()),
			"指定编码转换":   fmt.Sprintf("%s iconv -f GBK -t UTF-8 file.txt", qflag.Root.Name()),
			"自动检测并转换":  fmt.Sprintf("%s iconv -t UTF-8 file.txt", qflag.Root.Name()),
			"原地转换":     fmt.Sprintf("%s iconv -t UTF-8 -w file.txt", qflag.Root.Name()),
			"原地转换并备份":  fmt.Sprintf("%s iconv -t UTF-8 -w -b file.txt", qflag.Root.Name()),
			"批量转换":     fmt.Sprintf("%s iconv -t UTF-8 *.txt", qflag.Root.Name()),
			"输出到指定目录":  fmt.Sprintf("%s iconv -t UTF-8 -o output/ file.txt", qflag.Root.Name()),
			"输出到指定文件":  fmt.Sprintf("%s iconv -t UTF-8 -o out.txt file.txt", qflag.Root.Name()),
			"列出支持的编码":  fmt.Sprintf("%s iconv -l", qflag.Root.Name()),
		},
		Notes: []string{
			"-t 默认为 none (只检测不转换) , 指定目标编码时执行转换",
			"-f 默认为 auto (自动检测源编码) ",
			"使用 -w 原地写入会直接覆盖原文件",
			"使用 -o 可指定输出路径",
			"建议重要文件使用 -b 创建备份",
			"不支持递归目录处理, 请使用通配符或管道",
		},
		MutexGroups: []qflag.MutexGroup{
			{
				Name:      "to-list", // 目标编码和列出编码互斥
				Flags:     []string{"list", "to"},
				AllowNone: true,
			},
		},
		FlagDependencies: []qflag.FlagDependency{
			{
				Name:    "backup-write", // 使用 -b 时, -w 必须指定
				Trigger: "backup",
				Targets: []string{"write"},
				Type:    qflag.DepRequired,
			},
			{
				Name:    "write-to", // 使用 -w 时, -t 必须指定
				Trigger: "write",
				Targets: []string{"to"},
				Type:    qflag.DepRequired,
			},
		},
	}

	if err := IconvCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	IconvCmd.SetRun(runIconv)
}

// runIconv 执行 iconv 命令
func runIconv(cmd qflag.Command) error {
	// 列出支持的编码
	if iconvList.Get() {
		return iconv.ListEncodings()
	}

	// 检查文件参数
	files := cmd.Args()
	if len(files) == 0 {
		return fmt.Errorf("no input files specified")
	}

	config := iconv.IconvConfig{
		Files:        files,
		FromEncoding: iconvFrom.Get(),
		ToEncoding:   iconvTo.Get(),
		Write:        iconvWrite.Get(),
		Backup:       iconvBackup.Get(),
		Output:       iconvOutput.Get(),
		Quiet:        iconvQuiet.Get(),
		Short:        iconvShort.Get(),
	}

	return iconv.IconvCmdMain(config)
}
