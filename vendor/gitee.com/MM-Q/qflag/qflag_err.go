// 定义错误常量
package qflag

import "fmt"

// 命令行解析相关错误常量
const (
	ErrFlagParseFailed       = "Parameter parsing error"  // 全局实例标志解析错误
	ErrSubCommandParseFailed = "Subcommand parsing error" // 子命令标志解析错误
	ErrPanicRecovered        = "panic recovered"          // 恐慌捕获错误
	ErrValidationFailed      = "Validation failed"        // 参数验证失败错误
)

// NewValidationError 创建一个新的验证错误
func NewValidationError(message string) error {
	return fmt.Errorf("%s: %s", ErrValidationFailed, message)
}

// NewValidationErrorf 创建一个格式化的验证错误
func NewValidationErrorf(format string, v ...interface{}) error {
	return fmt.Errorf("%s: %s", ErrValidationFailed, fmt.Sprintf(format, v...))
}
