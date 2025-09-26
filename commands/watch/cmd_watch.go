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
	command := watchCmdCommand.Get()     // 命令
	interval := watchCmdInterval.Get()   // 间隔
	times := watchCmdTimes.Get()         // 运行次数
	exitOnError := watchCmdExitErr.Get() // 是否在错误时退出
	noTitle := watchCmdNoTitle.Get()     // 是否禁用标题
	timeout := watchCmdTimeout.Get()     // 超时时间
	shell := watchCmdShell.Get()         // shell类型

	// 验证命令参数
	if command == "" {
		return errors.New("command is empty")
	}

	// 验证运行次数参数
	if times < 0 {
		return errors.New("times must be greater than or equal to 0")
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

		// 检查执行次数限制
		if times > 0 && executionCount >= times {
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
		fmt.Printf("\n\n\n\n\n")

		// 等待间隔时间(如果不是最后一次执行)
		if interval > 0 && (times <= 0 || executionCount < times) {
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(interval): // 等待指定时间
			}
		}
	}

	return nil
}
