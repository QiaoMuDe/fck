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

	"gitee.com/MM-Q/shellx/shx"
)

// WatchConfig 监控配置结构体
type WatchConfig struct {
	Command     string        // 执行的命令
	Interval    time.Duration // 间隔时间
	MaxCount    int           // 最大执行次数
	ExitOnError bool          // 是否在错误时退出
	NoHeader    bool          // 轻度静默模式(不显示标题栏和换行符)
	Timeout     time.Duration // 超时时间
	ClearLines  int           // 清屏行数
	Quiet       bool          // 静默模式
}

// WatchCmdMain 执行监控命令
//
// 参数:
//   - config: 监控配置
//
// 返回值:
//   - error: 监控过程中可能发生的错误
func WatchCmdMain(config WatchConfig) error {
	command := config.Command
	interval := config.Interval
	maxCount := config.MaxCount
	exitOnError := config.ExitOnError
	noHeader := config.NoHeader
	timeout := config.Timeout
	clearLines := config.ClearLines
	quiet := config.Quiet

	if command == "" {
		return errors.New("command is empty")
	}

	if maxCount < -1 || maxCount == 0 {
		return errors.New("maxCount must be -1 (unlimited) or a positive number")
	}

	if timeout <= 0 {
		return errors.New("timeout must be greater than 0")
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-sigChan
		cancel()
	}()

	executionCount := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		if maxCount > 0 && executionCount >= maxCount {
			break
		}
		executionCount++

		if !noHeader && !quiet {
			fmt.Printf("%s\n", strings.Repeat("-", 50))
			fmt.Printf("Every %gs: %s [%s]\n\n", interval.Seconds(), time.Now().Format("2006-01-02 15:04:05"), command)
		}

		cmd := shx.New(command).WithTimeout(timeout)
		if !quiet {
			cmd = cmd.WithStderr(os.Stderr).WithStdout(os.Stdout)
		}
		err := cmd.Exec()
		if err != nil {
			fmt.Println(err)
			if exitOnError {
				break
			}
		}

		if clearLines > 0 && !noHeader && !quiet {
			fmt.Print(strings.Repeat("\n", clearLines))
		}

		if interval > 0 && (maxCount <= 0 || executionCount < maxCount) {
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(interval):
			}
		}
	}

	return nil
}
