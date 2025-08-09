// Package qerr 错误处理和自定义错误类型定义
// 本文件定义了qflag包使用的各种自定义错误类型，提供统一的错误处理机制，
// 包括验证错误、解析错误等不同类型的错误定义和处理方法。
package qerr

import (
	"errors"
	"fmt"
	"strings"
)

// 命令行解析相关错误变量
var (
	ErrFlagParseFailed       = errors.New("parameter parsing error")             // 全局实例标志解析错误
	ErrSubCommandParseFailed = errors.New("subcommand parsing error")            // 子命令标志解析错误
	ErrPanicRecovered        = errors.New("panic recovered")                     // 恐慌捕获错误
	ErrValidationFailed      = errors.New("validation failed")                   // 参数验证失败错误
	ErrEnvLoadFailed         = errors.New("environment variable loading failed") // 环境变量加载失败错误
	ErrAddSubCommandFailed   = errors.New("add subcommand failed")               // 添加子命令失败错误
)

// NewValidationError 创建一个新的验证错误
//
// 参数值:
//   - message string: 错误消息
//
// 返回值:
//   - error: 验证错误
func NewValidationError(message string) error {
	return fmt.Errorf("%s: %s", ErrValidationFailed, message)
}

// NewValidationErrorf 创建一个格式化的验证错误
//
// 参数值:
//   - format string: 格式化字符串
//   - v ...interface{}: 格式化参数
//
// 返回值:
//   - error: 验证错误
func NewValidationErrorf(format string, v ...interface{}) error {
	return fmt.Errorf("%s: %s", ErrValidationFailed, fmt.Sprintf(format, v...))
}

// JoinErrors 将错误切片合并为单个错误，并去除重复错误
//
// 参数值:
//   - errors []error: 错误切片
//
// 返回值:
//   - error: 合并后的错误
func JoinErrors(errors []error) error {
	if len(errors) == 0 {
		return nil
	}
	// 过滤nil错误
	nonNilErrors := make([]error, 0, len(errors))
	for _, err := range errors {
		if err != nil {
			nonNilErrors = append(nonNilErrors, err)
		}
	}
	if len(nonNilErrors) == 0 {
		return nil
	}
	if len(nonNilErrors) == 1 {
		return nonNilErrors[0]
	}

	// 使用切片和map保持插入顺序并去重
	seen := make(map[string]bool)
	uniqueErrors := make([]error, 0, len(nonNilErrors))
	for _, err := range nonNilErrors {
		errStr := err.Error()
		if !seen[errStr] {
			seen[errStr] = true
			uniqueErrors = append(uniqueErrors, err)
		}
	}

	// 构建错误信息
	var b strings.Builder
	b.WriteString(fmt.Sprintf("A total of %d unique errors:\n", len(uniqueErrors)))
	i := 1
	for _, err := range uniqueErrors {
		b.WriteString(fmt.Sprintf("  %d. %v\n", i, err))
		i++
	}

	// 使用常量格式字符串，将错误信息作为参数传入
	return fmt.Errorf("merged error message:\n%s", b.String())
}
