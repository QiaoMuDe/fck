package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/tee"
	"gitee.com/MM-Q/qflag"
)

var TeeCmd *qflag.Cmd

var (
	teeAppend           *qflag.BoolFlag // 追加模式
	teeIgnoreInterrupts *qflag.BoolFlag // 忽略中断信号
)

func init() {
	TeeCmd = qflag.NewCmd("tee", "", qflag.ExitOnError)

	teeAppend = TeeCmd.Bool("append", "a", "追加到文件 (默认覆盖)", false)
	teeIgnoreInterrupts = TeeCmd.Bool("ignore-interrupts", "i", "忽略中断信号 (SIGINT)", false)

	cmdOpts := &qflag.CmdOpts{
		Desc: "从标准输入读取并输出到多个目标",
		Notes: []string{
			"同时输出到标准输出和指定的文件",
			"支持追加模式 (-a)",
			"支持多文件同时写入",
			"使用 -i 忽略中断信号, 确保数据完整",
		},
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s tee [options] [file...]", qflag.Root.Name()),
		Examples: map[string]string{
			"写入文件":   fmt.Sprintf("echo 'hello' | %s tee file.txt", qflag.Root.Name()),
			"追加到文件":  fmt.Sprintf("echo 'line 2' | %s tee -a file.txt", qflag.Root.Name()),
			"多文件输出":  fmt.Sprintf("echo 'data' | %s tee file1.txt file2.txt", qflag.Root.Name()),
			"保存构建日志": fmt.Sprintf("make 2>&1 | %s tee build.log", qflag.Root.Name()),
			"忽略中断信号": fmt.Sprintf("long_cmd | %s tee -i output.log", qflag.Root.Name()),
			"管道链中使用": fmt.Sprintf("cat data | grep 'err' | %s tee errors.log | wc -l", qflag.Root.Name()),
		},
	}

	if err := TeeCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	TeeCmd.SetRun(runTee)
}

func runTee(cmd qflag.Command) error {
	config := tee.TeeConfig{
		Append:           teeAppend.Get(),
		IgnoreInterrupts: teeIgnoreInterrupts.Get(),
		Files:            cmd.Args(),
	}

	return tee.TeeCmdMain(config)
}
