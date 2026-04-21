package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"gitee.com/MM-Q/fck/internal/commands/grep"
	"gitee.com/MM-Q/qflag"
)

var GrepCmd *qflag.Cmd

var (
	grepIgnoreCase    *qflag.BoolFlag // -i 忽略大小写
	grepInvertMatch   *qflag.BoolFlag // -v 反向匹配
	grepLineNumber    *qflag.BoolFlag // -n 显示行号
	grepCount         *qflag.BoolFlag // -c 只显示计数
	grepWithFilename  *qflag.BoolFlag // -H 显示文件名
	grepRegexp        *qflag.BoolFlag // -E 使用正则表达式
	grepMaxCount      *qflag.IntFlag  // -m 最大匹配数
	grepBeforeContext *qflag.IntFlag  // -B 前文行数
	grepAfterContext  *qflag.IntFlag  // -A 后文行数
	grepContext       *qflag.IntFlag  // -C 上下文行数
	grepNoColor       *qflag.BoolFlag // --no-color 禁用颜色

	// 递归搜索相关
	grepRecursive     *qflag.BoolFlag        // -r 递归搜索
	grepFollowSymlink *qflag.BoolFlag        // -R 递归并跟随符号链接
	grepInclude       *qflag.StringSliceFlag // --include 包含文件模式
	grepExclude       *qflag.StringSliceFlag // --exclude 排除文件模式
	grepIncludeDir    *qflag.StringSliceFlag // --include-dir 包含目录模式
	grepExcludeDir    *qflag.StringSliceFlag // --exclude-dir 排除目录模式
	grepNoMessages    *qflag.BoolFlag        // -s 静默模式

	// 隐藏文件支持
	grepHidden *qflag.BoolFlag // --hidden 处理隐藏文件/目录（默认跳过）

	// 二进制文件处理
	grepText         *qflag.BoolFlag // -a, --text 强制将二进制文件视为文本处理
	grepIgnoreBinary *qflag.BoolFlag // -I 忽略二进制文件（不输出提示）
)

func init() {
	GrepCmd = qflag.NewCmd("grep", "", qflag.ExitOnError)

	grepIgnoreCase = GrepCmd.Bool("ignore-case", "i", "忽略大小写", false)
	grepInvertMatch = GrepCmd.Bool("invert-match", "v", "显示不匹配的行", false)
	grepLineNumber = GrepCmd.Bool("line-number", "n", "显示行号", false)
	grepCount = GrepCmd.Bool("count", "c", "只显示匹配行数", false)
	grepWithFilename = GrepCmd.Bool("with-filename", "H", "显示文件名", false)
	grepRegexp = GrepCmd.Bool("regexp", "E", "使用正则表达式匹配", false)
	grepMaxCount = GrepCmd.Int("max-count", "m", "最大匹配行数 (0表示无限制)", 0)
	grepBeforeContext = GrepCmd.Int("before-context", "B", "显示匹配行前N行", 0)
	grepAfterContext = GrepCmd.Int("after-context", "A", "显示匹配行后N行", 0)
	grepContext = GrepCmd.Int("context", "C", "显示匹配行前后N行", 0)
	grepNoColor = GrepCmd.Bool("no-color", "nc", "禁用颜色高亮", false)
	grepNoMessages = GrepCmd.Bool("no-messages", "s", "静默模式，不显示错误信息", false)

	// 递归搜索相关标志
	grepRecursive = GrepCmd.Bool("recursive", "r", "递归搜索目录", false)
	grepFollowSymlink = GrepCmd.Bool("follow-symlink", "R", "递归搜索并跟随符号链接", false)
	grepInclude = GrepCmd.StringSlice("include", "", "只搜索匹配模式的文件，多个模式用逗号分隔 (如 *.go,*.js)", []string{})
	grepExclude = GrepCmd.StringSlice("exclude", "", "排除匹配模式的文件，多个模式用逗号分隔", []string{})
	grepIncludeDir = GrepCmd.StringSlice("include-dir", "", "只进入匹配模式的目录，多个模式用逗号分隔", []string{})
	grepExcludeDir = GrepCmd.StringSlice("exclude-dir", "", "排除匹配模式的目录，多个模式用逗号分隔", []string{})

	// 隐藏文件相关标志
	grepHidden = GrepCmd.Bool("hidden", "hd", "处理隐藏文件/目录，默认跳过隐藏项", false)

	// 二进制文件处理标志
	grepText = GrepCmd.Bool("text", "a", "强制将二进制文件视为文本处理", false)
	grepIgnoreBinary = GrepCmd.Bool("ignore-binary", "I", "完全忽略二进制文件，不输出提示", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "文本搜索工具",
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s grep [options] <pattern> [file...]", qflag.Root.Name()),
		Examples: map[string]string{
			"基础搜索":       fmt.Sprintf("%s grep \"error\" log.txt", qflag.Root.Name()),
			"忽略大小写":      fmt.Sprintf("%s grep -i \"error\" log.txt", qflag.Root.Name()),
			"显示行号":       fmt.Sprintf("%s grep -n \"func\" main.go", qflag.Root.Name()),
			"显示文件名和行号":   fmt.Sprintf("%s grep -n -H \"TODO\" main.go", qflag.Root.Name()),
			"只显示计数":      fmt.Sprintf("%s grep -c \"http\" access.log", qflag.Root.Name()),
			"显示不匹配的行":    fmt.Sprintf("%s grep -v \"^#\" config.txt", qflag.Root.Name()),
			"正则表达式匹配":    fmt.Sprintf("%s grep -E \"^error|fatal\" log.txt", qflag.Root.Name()),
			"限制匹配数量":     fmt.Sprintf("%s grep -m 10 \"http\" access.log", qflag.Root.Name()),
			"显示上下文":      fmt.Sprintf("%s grep -C 3 \"exception\" log.txt", qflag.Root.Name()),
			"结合 find 使用": fmt.Sprintf("%s find . -name \"*.go\" -exec %s grep -n \"TODO\" {} \\;", qflag.Root.Name(), qflag.Root.Name()),
		},
		Notes: []string{
			"默认使用固定字符串匹配，更轻量高效",
			"默认高亮显示匹配关键字，使用 --no-color 禁用",
			"不指定 file 时从标准输入读取",
			"默认跳过隐藏文件/目录，使用 --hidden/-hd 包含",
			"默认跳过二进制文件并输出提示，使用 -a 强制处理，使用 -I 静默跳过",
		},
	}

	if err := GrepCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	GrepCmd.SetRun(runGrep)
}

// runGrep 执行 grep 命令
//
// 参数:
//   - cmd: qflag 命令对象
//
// 返回:
//   - error: 执行错误 (如果有)
func runGrep(cmd qflag.Command) error {
	args := cmd.Args()

	// 检查参数
	if len(args) < 1 {
		return fmt.Errorf("no pattern specified")
	}

	// 提取模式参数
	pattern := args[0]

	// 收集文件参数（从 args[1] 开始）
	var fileArgs []string
	if len(args) >= 2 {
		fileArgs = args[1:]
	}

	// 展开通配符
	var targets []string
	for _, arg := range fileArgs {
		matches, err := filepath.Glob(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error globbing %s: %v\n", arg, err)
			continue
		}
		if len(matches) == 0 {
			// 没有匹配，保留原参数（让 core 层处理文件不存在错误）
			targets = append(targets, arg)
		} else {
			targets = append(targets, matches...)
		}
	}

	// 准备基础配置
	baseConfig := grep.GrepConfig{
		Pattern:       pattern,
		IgnoreCase:    grepIgnoreCase.Get(),
		InvertMatch:   grepInvertMatch.Get(),
		LineNumber:    grepLineNumber.Get(),
		Count:         grepCount.Get(),
		WithFilename:  grepWithFilename.Get(),
		Regexp:        grepRegexp.Get(),
		MaxCount:      grepMaxCount.Get(),
		BeforeContext: grepBeforeContext.Get(),
		AfterContext:  grepAfterContext.Get(),
		Context:       grepContext.Get(),
		NoColor:       grepNoColor.Get(),
		Recursive:     grepRecursive.Get(),
		FollowSymlink: grepFollowSymlink.Get(),
		Include:       grepInclude.Get(),
		Exclude:       grepExclude.Get(),
		IncludeDir:    grepIncludeDir.Get(),
		ExcludeDir:    grepExcludeDir.Get(),
		NoMessages:    grepNoMessages.Get(),
		Hidden:        grepHidden.Get(),
		Text:          grepText.Get(),
		IgnoreBinary:  grepIgnoreBinary.Get(),
	}

	// 无文件参数：从 stdin 读取
	if len(targets) == 0 {
		baseConfig.Target = ""
		return grep.GrepCmdMain(baseConfig)
	}

	// 单文件：保持原有行为
	if len(targets) == 1 {
		baseConfig.Target = targets[0]
		return grep.GrepCmdMain(baseConfig)
	}

	// 多目标：循环处理（支持文件和目录共存）
	var hasError bool
	for _, target := range targets {
		config := baseConfig
		config.Target = target
		// 多文件时强制显示文件名
		config.WithFilename = true

		if err := grep.GrepCmdMain(config); err != nil {
			fmt.Fprintf(os.Stderr, "error processing %s: %v\n", target, err)
			hasError = true
		}
	}

	if hasError {
		return fmt.Errorf("some files failed to process")
	}
	return nil
}
