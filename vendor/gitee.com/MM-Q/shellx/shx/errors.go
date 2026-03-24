package shx

import (
	"context"
	"errors"
	"fmt"
	"time"

	"mvdan.cc/sh/v3/interp"
)

// 预定义错误
var (
	// ErrAlreadyExecuted 表示命令已经执行过
	ErrAlreadyExecuted = errors.New("command has already been executed")

	// ErrNilContext 表示上下文为 nil
	ErrNilContext = errors.New("context cannot be nil")

	// ErrNilReader 表示 reader 为 nil
	ErrNilReader = errors.New("reader cannot be nil")

	// ErrNilWriter 表示 writer 为 nil
	ErrNilWriter = errors.New("writer cannot be nil")
)

// handleError 处理执行错误
//
// 参数：
//   - err: 原始错误
//   - cmdStr: 命令字符串（用于错误信息）
//   - timeout: 超时时间
//
// 返回：
//   - 处理后的错误
func handleError(err error, cmdStr string, timeout time.Duration) error {
	if err == nil {
		return nil
	}

	// 检查是否是上下文取消
	if errors.Is(err, context.Canceled) {
		return fmt.Errorf("command canceled: %s", cmdStr)
	}

	// 检查是否是超时
	if errors.Is(err, context.DeadlineExceeded) {
		if timeout > 0 {
			return fmt.Errorf("command timed out after %v: %s", timeout, cmdStr)
		}
		return fmt.Errorf("command timed out: %s", cmdStr)
	}

	// 检查是否是退出状态错误
	var exitStatus interp.ExitStatus
	if errors.As(err, &exitStatus) {
		// 退出码错误不包装，由调用方处理
		return ExitStatus{Code: uint8(exitStatus)}
	}

	return fmt.Errorf("command failed: %s: %w", cmdStr, err)
}

// IsExitStatus 检查错误是否是退出状态错误
//
// 参数：
//   - err: 错误对象
//
// 返回：
//   - uint8: 退出码
//   - bool: 是否是退出状态错误
func IsExitStatus(err error) (uint8, bool) {
	// 先检查自定义ExitStatus
	var es ExitStatus
	if errors.As(err, &es) {
		return es.Code, true
	}

	// 再检查原生interp.ExitStatus（它是uint8类型）
	var is interp.ExitStatus
	if errors.As(err, &is) {
		return uint8(is), true
	}

	return 0, false
}
