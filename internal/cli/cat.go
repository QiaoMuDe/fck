package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/cat"
	"gitee.com/MM-Q/qflag"
)

var CatCmd *qflag.Cmd

var (
	catShowLineNum  *qflag.BoolFlag // -n 显示所有行号
	catShowNonBlank *qflag.BoolFlag // -b 显示非空行行号
	catShowEnd      *qflag.BoolFlag // -E 显示行尾$
	catShowTabs     *qflag.BoolFlag // -T 显示制表符为^I
	catShowAll      *qflag.BoolFlag // -A 等价于 -ET
	catHeadLines    *qflag.IntFlag  // -u 显示前N行
	catTailLines    *qflag.IntFlag  // -d 显示后N行
	catQuiet        *qflag.BoolFlag // -q 静默模式

	// 二进制文件处理
	catText         *qflag.BoolFlag // -a, --text 强制将二进制文件视为文本处理
	catIgnoreBinary *qflag.BoolFlag // -I, --ignore-binary 忽略二进制文件

	// 分页查看
	catUseLess *qflag.BoolFlag // -l, --less 使用分页器查看文件内容

	// 语法高亮
	catHighlight *qflag.BoolFlag // -H, --highlight 启用语法高亮

	// 文件大小限制
	catMaxSize *qflag.SizeFlag // -S, --max-size 最大文件大小 (默认100MB)
)

func init() {
	CatCmd = qflag.NewCmd("cat", "", qflag.ExitOnError)

	catShowLineNum = CatCmd.Bool("number", "n", "显示所有行号", false)
	catShowNonBlank = CatCmd.Bool("number-nonblank", "b", "仅显示非空行行号", false)
	catShowEnd = CatCmd.Bool("show-ends", "E", "在每行末尾显示$", false)
	catShowTabs = CatCmd.Bool("show-tabs", "T", "将制表符显示为^I", false)
	catShowAll = CatCmd.Bool("show-all", "A", "等价于 -ET", false)
	catHeadLines = CatCmd.Int("head", "u", "显示前N行 (0表示全部)", 0)
	catTailLines = CatCmd.Int("tail", "d", "显示后N行 (0表示全部)", 0)
	catQuiet = CatCmd.Bool("quiet", "q", "静默模式 (不显示错误)", false)

	// 二进制文件处理标志
	catText = CatCmd.Bool("text", "a", "强制将二进制文件视为文本处理", false)
	catIgnoreBinary = CatCmd.Bool("ignore-binary", "I", "完全忽略二进制文件，不输出提示", false)

	// 分页查看标志
	catUseLess = CatCmd.Bool("less", "l", "使用分页器查看文件内容", false)

	// 语法高亮标志
	catHighlight = CatCmd.Bool("highlight", "H", "启用语法高亮", false)

	// 文件大小限制标志
	catMaxSize = CatCmd.Size("max-size", "S", "最大文件大小 (如: 50m, 100m, 1g)", 100*1024*1024)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "显示文件内容",
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s cat [options] <file>", qflag.Root.Name()),
		Examples: map[string]string{
			"显示文件内容":    fmt.Sprintf("%s cat file.txt", qflag.Root.Name()),
			"显示行号":      fmt.Sprintf("%s cat -n file.txt", qflag.Root.Name()),
			"显示前10行":    fmt.Sprintf("%s cat -u 10 file.txt", qflag.Root.Name()),
			"显示后5行":     fmt.Sprintf("%s cat -d 5 file.txt", qflag.Root.Name()),
			"静默跳过二进制文件": fmt.Sprintf("%s cat -I file.bin", qflag.Root.Name()),
			"强制显示二进制内容": fmt.Sprintf("%s cat -a file.bin", qflag.Root.Name()),
			"使用分页器查看":   fmt.Sprintf("%s cat -l file.txt", qflag.Root.Name()),
			"启用语法高亮":    fmt.Sprintf("%s cat -H main.go", qflag.Root.Name()),
			"分页器+语法高亮":  fmt.Sprintf("%s cat -l -H main.go", qflag.Root.Name()),
		},
		Notes: []string{
			"默认检测到二进制文件会输出提示，使用 -a 强制处理，使用 -I 静默跳过",
		},
		MutexGroups: []qflag.MutexGroup{
			{
				Name:      "head-tail",
				Flags:     []string{"head", "tail"},
				AllowNone: true,
			},
			{
				Name:      "line-number",
				Flags:     []string{"number", "number-nonblank"},
				AllowNone: true,
			},
			{
				Name:      "highlight-show-all",
				Flags:     []string{"highlight", "show-all"},
				AllowNone: true,
			},
			{
				Name:      "highlight-show-ends",
				Flags:     []string{"highlight", "show-ends"},
				AllowNone: true,
			},
			{
				Name:      "highlight-show-tabs",
				Flags:     []string{"highlight", "show-tabs"},
				AllowNone: true,
			},
		},
	}

	if err := CatCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	CatCmd.SetRun(runCat)
}

func runCat(cmd qflag.Command) error {
	args := cmd.Args()
	// 只保留第一个文件参数
	if len(args) > 1 {
		args = args[:1]
	}

	config := cat.CatConfig{
		Targets:      args,
		ShowLineNum:  catShowLineNum.Get(),
		ShowNonBlank: catShowNonBlank.Get(),
		ShowEnd:      catShowEnd.Get(),
		ShowTabs:     catShowTabs.Get(),
		ShowAll:      catShowAll.Get(),
		HeadLines:    catHeadLines.Get(),
		MaxSize:      catMaxSize.Get(),
		TailLines:    catTailLines.Get(),
		Quiet:        catQuiet.Get(),
		Text:         catText.Get(),
		IgnoreBinary: catIgnoreBinary.Get(),
		UseLess:      catUseLess.Get(),
		Highlight:    catHighlight.Get(),
	}

	return cat.CatCmdMain(config)
}
