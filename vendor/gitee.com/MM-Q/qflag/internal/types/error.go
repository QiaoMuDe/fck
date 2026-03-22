// Package types 定义了qflag项目的核心类型和接口
//
// types 包提供了整个项目的基础类型定义, 包括:
//   - 标志类型和接口定义
//   - 命令接口定义
//   - 注册表接口定义
//   - 错误处理类型
//
// 这些类型和接口构成了整个框架的核心抽象层,
// 为具体的实现提供了统一的规范和契约。
package types

import (
	"errors"
	"fmt"
)

// Error 错误类型
//
// Error 是qflag项目的标准错误类型, 提供了结构化的错误信息。
// 包含错误码、错误消息和原始错误, 便于错误分类和处理。
//
// 字段说明:
//   - Code: 错误码, 用于错误分类和程序化处理
//   - Message: 错误消息, 面向用户的描述信息
//   - Cause: 原始错误, 包装的底层错误
//
// 特性:
//   - 实现error接口
//   - 支持错误链 (errors.Unwrap)
//   - 支持错误比较 (errors.Is)
//   - 提供错误码匹配
type Error struct {
	Code    string // 错误码, 用于错误分类
	Message string // 错误消息, 面向用户
	Cause   error  // 原始错误, 底层错误原因
}

// NewError 创建新的错误
//
// 参数:
//   - code: 错误码, 用于错误分类和识别
//   - message: 错误消息, 面向用户的描述信息
//   - cause: 原始错误, 可以为nil
//
// 返回值:
//   - *Error: 新创建的错误实例
//
// 功能说明:
//   - 创建结构化的错误实例
//   - 保留原始错误信息
//   - 提供错误分类能力
func NewError(code, message string, cause error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// Error 实现 error 接口
//
// 返回值:
//   - string: 格式化的错误字符串
//
// 功能说明:
//   - 返回用户友好的错误信息
//   - 包含原始错误信息 (如果有)
//   - 格式: 消息 + ": " + 原始错误
func (e *Error) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

// Unwrap 实现 errors.Unwrap 接口
//
// 返回值:
//   - error: 原始错误
//
// 功能说明:
//   - 支持错误链操作
//   - 允许使用errors.Unwrap获取底层错误
//   - 支持errors.As和errors.Is
func (e *Error) Unwrap() error {
	return e.Cause
}

// Is 判断错误是否相同
//
// 参数:
//   - target: 要比较的目标错误
//
// 返回值:
//   - bool: 是否相同, true表示相同
//
// 功能说明:
//   - 基于错误码进行比较
//   - 支持errors.Is函数
//   - 忽略错误消息和原始错误
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// 预定义错误码
//
// 以下是项目中常用的预定义错误, 可以直接使用或作为参考。
// 所有预定义错误都使用NewError创建, 保持一致的错误结构。
var (
	// ErrInvalidFlagType 无效的标志类型错误
	//
	// 使用场景:
	//   - 传入不支持的标志类型
	//   - 标志类型转换失败
	ErrInvalidFlagType = NewError("INVALID_FLAG_TYPE", "invalid flag type", nil)

	// ErrFlagNotFound 标志不存在错误
	//
	// 使用场景:
	//   - 查找不存在的标志
	//   - 引用未注册的标志
	ErrFlagNotFound = NewError("FLAG_NOT_FOUND", "flag not found", nil)

	// ErrCmdNotFound 命令不存在错误
	//
	// 使用场景:
	//   - 查找不存在的命令
	//   - 引用未注册的命令
	ErrCmdNotFound = NewError("COMMAND_NOT_FOUND", "cmd not found", nil)

	// ErrFlagAlreadyExists 标志已存在错误
	//
	// 使用场景:
	//   - 注册重复的标志
	//   - 标志名称冲突
	ErrFlagAlreadyExists = NewError("FLAG_ALREADY_EXISTS", "flag already exists", nil)

	// ErrCmdAlreadyExists 命令已存在错误
	//
	// 使用场景:
	//   - 注册重复的命令
	//   - 命令名称冲突
	ErrCmdAlreadyExists = NewError("COMMAND_ALREADY_EXISTS", "cmd already exists", nil)

	// ErrParseFailed 解析失败错误
	//
	// 使用场景:
	//   - 命令行参数解析失败
	//   - 配置文件解析失败
	ErrParseFailed = NewError("PARSE_FAILED", "parse failed", nil)

	// ErrValidationFailed 验证失败错误
	//
	// 使用场景:
	//   - 标志值验证失败
	//   - 业务规则验证失败
	ErrValidationFailed = NewError("VALIDATION_FAILED", "validation failed", nil)

	// ErrRequiredFlag 必填标志缺失错误
	//
	// 使用场景:
	//   - 必填标志未提供
	//   - 必填标志值为空
	ErrRequiredFlag = NewError("REQUIRED_FLAG", "required flag is missing", nil)

	// ErrInvalidValue 无效值错误
	//
	// 使用场景:
	//   - 标志值格式错误
	//   - 标志值超出范围
	ErrInvalidValue = NewError("INVALID_VALUE", "invalid flag value", nil)
)

// WrapError 包装错误
//
// 参数:
//   - err: 要包装的原始错误
//   - code: 新的错误码
//   - message: 新的错误消息
//
// 返回值:
//   - *Error: 包装后的错误
//
// 功能说明:
//   - 为现有错误添加上下文信息
//   - 保持原始错误链
//   - 提供新的错误分类
//
// 使用场景:
//   - 为底层错误添加业务上下文
//   - 统一错误处理格式
//   - 错误转换和适配
func WrapError(err error, code, message string) *Error {
	return NewError(code, message, err)
}

// WrapParseError 包装解析错误, 专门用于标志解析场景
//
// 参数:
//   - err: 原始解析错误
//   - flagType: 标志类型描述
//   - value: 解析失败的值
//
// 返回值:
//   - *Error: 包装后的解析错误
//
// 功能说明:
//   - 专门用于标志解析错误
//   - 自动生成描述性错误消息
//   - 保留原始错误信息
//
// 使用场景:
//   - 标志值解析失败
//   - 类型转换错误
//   - 格式验证错误
func WrapParseError(err error, flagType, value string) *Error {
	if err == nil {
		return nil
	}
	return NewError("PARSE_ERROR",
		fmt.Sprintf("failed to parse %s value: %s", flagType, value),
		err)
}

// IsNotFoundError 判断是否为"未找到"错误
//
// 参数:
//   - err: 要检查的错误
//
// 返回值:
//   - bool: 是否为未找到错误, true表示是
//
// 功能说明:
//   - 检查错误码是否为FLAG_NOT_FOUND或COMMAND_NOT_FOUND
//   - 支持错误链检查
//   - 便于统一处理未找到类型的错误
//
// 使用场景:
//   - 统一处理资源不存在的情况
//   - 区分未找到错误和其他错误
//   - 简化错误处理逻辑
func IsNotFoundError(err error) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Code == ErrFlagNotFound.Code || e.Code == ErrCmdNotFound.Code
	}
	return false
}
