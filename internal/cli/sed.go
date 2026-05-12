package cli

import (
	"fmt"
	"os"

	"gitee.com/MM-Q/fck/internal/commands/sed"
	"gitee.com/MM-Q/fck/internal/types"
	"gitee.com/MM-Q/go-kit/fs"
	"gitee.com/MM-Q/go-kit/term"
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

	// 二进制文件处理
	sedText         *qflag.BoolFlag // -a, --text 强制将二进制文件视为文本处理
	sedIgnoreBinary *qflag.BoolFlag // --ignore-binary 忽略二进制文件（不输出提示）
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
	sedIgnoreCase = SedCmd.Bool("ignore-case", "ic", "忽略大小写", false)
	sedMaxBuffer = SedCmd.Size("buffer-size", "bs", "最大行缓冲区大小", types.DefaultMaxBufferSize)

	// 二进制文件处理标志
	sedText = SedCmd.Bool("text", "a", "强制将二进制文件视为文本处理", false)
	sedIgnoreBinary = SedCmd.Bool("ignore-binary", "I", "完全忽略二进制文件，不输出提示", false)

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
			"忽略大小写":      fmt.Sprintf("%s sed -ic -p \"hello\" -r \"world\" file.txt", qflag.Root.Name()),
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

	// 检查缓冲区大小（不能小于等于初始缓冲区大小）
	maxBuffer := sedMaxBuffer.Get()
	if maxBuffer <= types.InitialBufferSize {
		return fmt.Errorf("buffer size must be greater than %d bytes", types.InitialBufferSize)
	}

	// 检测是否为管道输入
	isPipe := term.IsStdinPipe()

	// === 管道输入模式 ===
	if isPipe {
		// 检查冲突标志
		if sedInPlace.Get() {
			return fmt.Errorf("cannot use -i flag when reading from stdin")
		}

		// 准备配置
		config := sed.SedConfig{
			Target:       "-", // 使用 "-" 表示 stdin
			Pattern:      sedPattern.Get(),
			Replacement:  sedReplacement.Get(),
			Regexp:       sedRegexp.Get(),
			LineRange:    sedLineRange.Get(),
			InPlace:      false, // 管道模式强制预览
			Backup:       false,
			MaxCount:     sedMaxCount.Get(),
			IgnoreCase:   sedIgnoreCase.Get(),
			MaxBuffer:    maxBuffer,
			Text:         sedText.Get(),
			IgnoreBinary: sedIgnoreBinary.Get(),
		}

		return sed.SedCmdMain(config)
	}

	// === 文件模式 ===
	// 检查文件参数（非管道模式下需要）
	if len(args) < 1 {
		return fmt.Errorf("no target file specified")
	}

	// 展开通配符并收集所有文件
	targets, err := fs.Expand(args)
	if err != nil {
		return err
	}

	// 准备基础配置（所有文件共用）
	baseConfig := sed.SedConfig{
		Pattern:      sedPattern.Get(),
		Replacement:  sedReplacement.Get(),
		Regexp:       sedRegexp.Get(),
		LineRange:    sedLineRange.Get(),
		InPlace:      sedInPlace.Get(),
		Backup:       sedBackup.Get(),
		MaxCount:     sedMaxCount.Get(),
		IgnoreCase:   sedIgnoreCase.Get(),
		MaxBuffer:    maxBuffer,
		Text:         sedText.Get(),
		IgnoreBinary: sedIgnoreBinary.Get(),
	}

	// 处理多个文件
	var hasError bool
	for _, target := range targets {
		config := baseConfig
		config.Target = target

		// 多文件时打印文件名（预览模式）
		if len(targets) > 1 && !config.InPlace {
			fmt.Printf("==> %s <==\n", target)
		}

		if err := sed.SedCmdMain(config); err != nil {
			fmt.Fprintf(os.Stderr, "error processing %s: %v\n", target, err)
			hasError = true
			// 继续处理下一个文件
		}

		// 多文件预览时添加空行分隔
		if len(targets) > 1 && !config.InPlace {
			fmt.Println()
		}
	}

	if hasError {
		return fmt.Errorf("some files failed to process")
	}
	return nil
}
