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
	// ErrNilContext 表示上下文为 nil
	ErrNilContext = errors.New("context cannot be nil")
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
		// 退出码错误不包装额外信息，但保留原始错误链
		return ExitStatus{Code: uint8(exitStatus), err: err}
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

// UnclosedQuoteError 表示命令字符串中存在未闭合的引号
type UnclosedQuoteError struct {
	QuoteType rune // 未闭合的引号类型 (', ", 或 `)
}

func (e *UnclosedQuoteError) Error() string {
	return fmt.Sprintf("unclosed quote in command string: %c", e.QuoteType)
}

// QuoteType 返回未闭合的引号类型
func (e *UnclosedQuoteError) GetQuoteType() rune {
	return e.QuoteType
}

// SyntaxError 表示 shell 语法错误
//
// 当解析 shell 脚本遇到语法错误时，返回此类型，包含错误位置和描述信息。
type SyntaxError struct {
	File    string // 文件名（如果来自文件则为文件路径，否则为空字符串）
	Line    int    // 错误行号（从 1 开始）
	Column  int    // 错误列号（从 1 开始）
	Message string // 错误描述
}

// Error 实现 error 接口
func (e *SyntaxError) Error() string {
	if e.File != "" {
		return e.File + ":" + e.Message
	}
	if e.Line > 0 && e.Column > 0 {
		return fmt.Sprintf("line %d:%d: %s", e.Line, e.Column, e.Message)
	}
	return e.Message
}
