package qerr

import (
	"errors"
	"fmt"
	"strings"
)

// 命令行解析相关错误变量
var (
	ErrFlagParseFailed       = errors.New("Parameter parsing error")  // 全局实例标志解析错误
	ErrSubCommandParseFailed = errors.New("Subcommand parsing error") // 子命令标志解析错误
	ErrPanicRecovered        = errors.New("panic recovered")          // 恐慌捕获错误
	ErrValidationFailed      = errors.New("Validation failed")        // 参数验证失败错误
)

// NewValidationError 创建一个新的验证错误
func NewValidationError(message string) error {
	return fmt.Errorf("%s: %s", ErrValidationFailed, message)
}

// NewValidationErrorf 创建一个格式化的验证错误
func NewValidationErrorf(format string, v ...interface{}) error {
	return fmt.Errorf("%s: %s", ErrValidationFailed, fmt.Sprintf(format, v...))
}

// JoinErrors 将错误切片合并为单个错误，并去除重复错误
func JoinErrors(errors []error) error {
	if len(errors) == 0 {
		return nil
	}
	if len(errors) == 1 {
		return errors[0]
	}

	// 使用切片和map保持插入顺序并去重
	seen := make(map[string]bool)
	uniqueErrors := make([]error, 0, len(errors))
	for _, err := range errors {
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
	return fmt.Errorf("Merged error message:\n%s", b.String())
}
