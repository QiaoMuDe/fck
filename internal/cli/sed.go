package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/sed"
	"gitee.com/MM-Q/fck/internal/types"
	"gitee.com/MM-Q/qflag"
)

// SedCmd sed 命令
var SedCmd *qflag.Cmd

// sed 命令标志
var (
	sedPattern     *qflag.StringFlag // -p 匹配模式
	sedReplacement *qflag.StringFlag // -r 替换内容
	sedRegexp      *qflag.BoolFlag   // -E 使用正则表达式
	sedLineRange   *qflag.StringFlag // -n 行号范围
	sedInPlace     *qflag.BoolFlag   // -i 原地修改
	sedBackup      *qflag.BoolFlag   // -b 创建备份
	sedMaxCount    *qflag.IntFlag    // -c 最大替换次数
	sedIgnoreCase  *qflag.BoolFlag   // -I 忽略大小写
	sedMaxBuffer   *qflag.SizeFlag   // --buffer-size 最大行缓冲区
)

func init() {
	SedCmd = qflag.NewCmd("sed", "", qflag.ExitOnError)

	// 定义标志
	sedPattern = SedCmd.String("pattern", "p", "匹配模式 (固定字符串或正则) ", "")
	sedReplacement = SedCmd.String("replacement", "r", "替换内容", "")
	sedRegexp = SedCmd.Bool("regexp", "E", "使用正则表达式匹配", false)
	sedLineRange = SedCmd.String("line", "n", "指定行号范围 (如: 5 或 1,10) ", "")
	sedInPlace = SedCmd.Bool("in-place", "i", "原地修改文件 (默认输出到屏幕) ", false)
	sedBackup = SedCmd.Bool("backup", "b", "修改前创建 .bak 备份", false)
	sedMaxCount = SedCmd.Int("count", "c", "最大替换次数 (0表示无限制) ", 0)
	sedIgnoreCase = SedCmd.Bool("ignore-case", "I", "忽略大小写", false)
	sedMaxBuffer = SedCmd.Size("buffer-size", "bs", "最大行缓冲区大小", types.DefaultMaxBufferSize)

	// 命令配置
	cmdOpts := &qflag.CmdOpts{
		Desc:        "文本替换工具",
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s sed [options] <file>", qflag.Root.Name()),
		Examples: map[string]string{
			"基础替换 (预览) ": fmt.Sprintf("%s sed -p \"old\" -r \"new\" file.txt", qflag.Root.Name()),
			"原地修改文件":     fmt.Sprintf("%s sed -p \"old\" -r \"new\" -i file.txt", qflag.Root.Name()),
			"修改并备份":      fmt.Sprintf("%s sed -p \"old\" -r \"new\" -i -b file.txt", qflag.Root.Name()),
			"使用正则表达式":    fmt.Sprintf("%s sed -E -p \"foo\\d+\" -r \"bar\" file.txt", qflag.Root.Name()),
			"只替换第 5 行":   fmt.Sprintf("%s sed -p \"old\" -r \"new\" -n 5 file.txt", qflag.Root.Name()),
			"替换 1-10 行":  fmt.Sprintf("%s sed -p \"old\" -r \"new\" -n 1,10 file.txt", qflag.Root.Name()),
			"只替换前 3 处":   fmt.Sprintf("%s sed -p \"old\" -r \"new\" -c 3 file.txt", qflag.Root.Name()),
			"忽略大小写":      fmt.Sprintf("%s sed -I -p \"hello\" -r \"world\" file.txt", qflag.Root.Name()),
		},
		Notes: []string{
			"默认输出到屏幕预览, 使用 -i 才修改原文件",
			"使用 -i -b 组合可在修改前自动创建 .bak 备份",
			"行号范围格式: 5 (单行) 、1,10 (范围) 、5, (第5行到末尾) ",
			"正则模式使用 Go 正则语法",
		},
	}

	if err := SedCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	SedCmd.SetRun(runSed)
}

// runSed 执行 sed 命令
//
// 参数:
//   - cmd: qflag 命令对象
//
// 返回:
//   - error: 执行错误
func runSed(cmd qflag.Command) error {
	args := cmd.Args()

	// 检查文件参数
	if len(args) < 1 {
		return fmt.Errorf("no target file specified")
	}

	// 检查缓冲区大小（不能小于等于初始缓冲区大小）
	maxBuffer := sedMaxBuffer.Get()
	if maxBuffer <= types.InitialBufferSize {
		return fmt.Errorf("buffer size must be greater than %d bytes", types.InitialBufferSize)
	}

	config := sed.SedConfig{
		Target:      args[0],
		Pattern:     sedPattern.Get(),
		Replacement: sedReplacement.Get(),
		Regexp:      sedRegexp.Get(),
		LineRange:   sedLineRange.Get(),
		InPlace:     sedInPlace.Get(),
		Backup:      sedBackup.Get(),
		MaxCount:    sedMaxCount.Get(),
		IgnoreCase:  sedIgnoreCase.Get(),
		MaxBuffer:   maxBuffer,
	}

	return sed.SedCmdMain(config)
}
