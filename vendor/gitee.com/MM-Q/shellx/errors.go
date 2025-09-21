// Package shellx 错误处理模块
// 本文件定义了 shellx 包中的错误类型、错误变量和错误处理函数，包括：
//   - 预定义的错误变量（超时、取消、未启动等）
//   - 错误消息常量定义
//   - 智能错误判断和分类函数 judgeError
//
// 提供统一的错误处理机制，能够准确识别和格式化各种命令执行错误。
package shellx

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
)

// 预定义的错误变量
var (
	// ErrAlreadyExecuted 表示命令已经执行过
	ErrAlreadyExecuted = errors.New("command has already been executed")
	// ErrNotStarted 表示命令尚未启动
	ErrNotStarted = errors.New("command has not been started")
	// ErrNoProcess 表示没有进程可操作
	ErrNoProcess = errors.New("no process to operate")
)

// 错误消息常量
const (
	// 超时和取消错误消息
	msgTimeoutExceeded = "command execution timeout: %s exceeded the %v time limit" // 超时错误消息
	msgCanceled        = "command execution canceled: %s was interrupted"           // 取消错误消息

	// exec 包错误消息
	msgErrDot       = "cannot execute current directory (security restriction): %s" // 执行当前目录错误消息
	msgErrNotFound  = "command not found: %s is not a valid command or executable"  // 命令未找到错误消息
	msgErrWaitDelay = "command execution failed: %s process wait timeout occurred"  // 执行等待超时错误消息

	// 退出码和系统错误消息
	msgExitCode    = "command execution failed: %s exited with code %d"      // 退出码错误消息
	msgSystemError = "system error: %s encountered an unexpected error - %v" // 系统错误消息
)

// judgeError 判断错误类型并返回对应的错误信息
//
// 参数:
//   - err: 错误对象
//   - c: Command 对象，用于获取用户上下文
//
// 返回值:
//   - error: 返回对应的错误信息
func judgeError(err error, c *Command) error {
	if err == nil {
		return nil
	}

	// 获取命令字符串用于错误提示
	var cmdStr string
	if c != nil {
		cmdStr = c.CmdStr()
	}

	// 检查是否为用户取消错误或超时错误
	if c != nil && c.userCtx != nil {
		ctxErr := c.userCtx.Err()
		switch {
		case errors.Is(ctxErr, context.DeadlineExceeded): // 超时错误
			return fmt.Errorf(msgTimeoutExceeded, cmdStr, c.getEffectiveTimeout())

		case errors.Is(ctxErr, context.Canceled): // 上下文取消错误
			return fmt.Errorf(msgCanceled, cmdStr)
		}
	}

	// 检查是否为 exec 包错误
	switch {
	case errors.Is(err, exec.ErrDot): // 无法执行当前目录
		return fmt.Errorf(msgErrDot, cmdStr)

	case errors.Is(err, exec.ErrNotFound): // 命令未找到
		return fmt.Errorf(msgErrNotFound, cmdStr)

	case errors.Is(err, exec.ErrWaitDelay): // 管道关闭延迟错误
		return fmt.Errorf(msgErrWaitDelay, cmdStr)
	}

	// 检查退出码错误
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode := exitErr.ExitCode()
		return fmt.Errorf(msgExitCode, cmdStr, exitCode)
	}

	// 其他系统错误
	return fmt.Errorf(msgSystemError, cmdStr, err)
}
