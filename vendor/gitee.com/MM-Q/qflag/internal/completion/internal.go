// Package completion 自动补全内部实现
// 本文件包含了自动补全功能的内部实现逻辑，提供补全算法、
// 匹配策略等核心功能的底层支持。
package completion

import (
	"fmt"
	"strings"
	"sync"

	"gitee.com/MM-Q/qflag/flags"
	"gitee.com/MM-Q/qflag/internal/types"
)

// 通用字符串构建器对象池 - 供bash和pwsh共同使用
var stringBuilderPool = sync.Pool{
	New: func() interface{} {
		builder := &strings.Builder{}
		builder.Grow(512) // 预分配512字节容量
		return builder
	},
}

// buildString 使用对象池构建字符串的通用辅助函数
// 供bash_completion.go和pwsh_completion.go共同使用
func buildString(fn func(*strings.Builder)) string {
	builder := stringBuilderPool.Get().(*strings.Builder)
	defer func() {
		// 如果容量过大则不回收，避免内存浪费
		if builder.Cap() <= 8192 {
			builder.Reset()
			stringBuilderPool.Put(builder)
		}
	}()

	fn(builder)
	return builder.String()
}

// getValueTypeByFlagType 根据标志类型获取值类型
//
// 参数:
//   - flagType - 标志类型
//
// 返回值:
//   - string: 值类型
func getValueTypeByFlagType(flagType flags.FlagType) string {
	switch flagType {
	case flags.FlagTypeBool:
		return "bool"
	case flags.FlagTypeEnum:
		return "enum"
	default:
		return "string"
	}
}

// getParamTypeByFlagType 根据标志类型获取参数需求类型
//
// 参数:
//   - flagType - 标志类型
//
// 返回值:
//   - string: 参数需求类型
func getParamTypeByFlagType(flagType flags.FlagType) string {
	if flagType == flags.FlagTypeBool {
		return "none"
	}
	return "required"
}

// validateCompletionGeneration 验证补全脚本生成所需的命令状态
//
// 参数:
//   - c: 命令实例
//
// 返回值:
//   - error: 如果验证失败, 返回相应的错误信息; 否则返回nil
func validateCompletionGeneration(ctx *types.CmdContext) error {
	if ctx == nil {
		return fmt.Errorf("command instance is nil")
	}
	if ctx.Parent != nil {
		return fmt.Errorf("invalid command state: not a root command")
	}
	if ctx.FlagRegistry == nil {
		return fmt.Errorf("invalid command state: flag registry is nil")
	}
	return nil
}
