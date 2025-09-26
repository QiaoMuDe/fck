// Package watch 实现了命令监控功能。
// 该文件提供了周期性执行指定命令并显示输出结果的核心功能，支持间隔设置、次数限制、颜色输出等。
package watch

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"gitee.com/MM-Q/shellx"
)

// 定义一个支持的shell类型列表
var supportedShells = []string{
	"bash",
	"cmd",
	"pwsh",
	"powershell",
	"sh",
	"none",
	"def1",
	"def2",
}

// 定义一个支持的shell类型映射
var shellMap = map[string]shellx.ShellType{
	"bash":       shellx.ShellBash,
	"cmd":        shellx.ShellCmd,
	"pwsh":       shellx.ShellPwsh,
	"powershell": shellx.ShellPowerShell,
	"sh":         shellx.ShellSh,
	"none":       shellx.ShellNone,
	"def1":       shellx.ShellDef1,
	"def2":       shellx.ShellDef2,
}

// WatchCmdMain 是 watch 子命令的主函数
//
// 参数:
//   - cl: 用于打印输出的 ColorLib 对象
//
// 返回:
//   - error: 如果发生错误，返回错误信息，否则返回 nil
func WatchCmdMain() error {
	// 获取命令参数
	args := watchCmd.Args()                // 执行的命令
	interval := watchCmdInterval.Get()     // 间隔时间
	maxCount := watchCmdMaxCount.Get()     // 最大执行次数
	exitOnError := watchCmdExitErr.Get()   // 是否在错误时退出
	noTitle := watchCmdNoTitle.Get()       // 是否禁用标题
	timeout := watchCmdTimeout.Get()       // 超时时间
	shell := watchCmdShell.Get()           // shell类型
	clearLines := watchCmdClearLines.Get() // 清屏行数

	// 验证命令参数
	var command string
	if len(args) == 0 {
		return errors.New("command is empty")
	} else if len(args) == 1 {
		command = args[0]
	} else {
		command = strings.Join(args, " ")
	}

	// 验证最大执行次数参数
	if maxCount < -1 || maxCount == 0 {
		return errors.New("maxCount must be -1 (unlimited) or a positive number")
	}

	// 验证超时时间参数
	if timeout <= 0 {
		return errors.New("timeout must be greater than 0")
	}

	// 设置信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动信号监听
	go func() {
		<-sigChan
		cancel()
	}()

	// 执行计数器
	executionCount := 0

	// 主监控循环
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		// 检查最大执行次数限制
		if maxCount > 0 && executionCount >= maxCount {
			break
		}
		executionCount++

		// 显示标题(如果未禁用)
		if !noTitle {
			fmt.Printf("%s\n", strings.Repeat("-", 50))
			fmt.Printf("Every %gs: %s [%s]\n\n", interval.Seconds(), time.Now().Format("2006-01-02 15:04:05"), command)
		}

		// 执行命令
		err := shellx.NewCmdStr(command).WithShell(shellMap[strings.ToLower(shell)]).WithTimeout(timeout).WithStderr(os.Stderr).WithStdout(os.Stdout).Exec()
		if err != nil {
			fmt.Println(err)
			if exitOnError {
				break
			}
		}

		// 清屏(如果指定了清屏行数)
		if clearLines > 0 {
			fmt.Print(strings.Repeat("\n", clearLines))
		}

		// 等待间隔时间(如果不是最后一次执行)
		if interval > 0 && (maxCount <= 0 || executionCount < maxCount) {
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(interval): // 等待指定时间
			}
		}
	}

	return nil
}
