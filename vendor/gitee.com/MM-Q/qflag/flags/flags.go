// Package flags 标志类型定义和接口
// 本文件定义了所有标志类型的通用接口和基础标志结构体，包括标志类型枚举、
// 验证器接口、标志接口等核心定义，为整个标志系统提供基础类型支持。
package flags

import (
	"flag"
	"fmt"
	"sort"
	"strings"
	"time"
)

// 定义非法字符集常量, 防止非法的标志名称
const InvalidFlagChars = " !@#$%^&*(){}[]|\\;:'\"<>,.?/"

// ErrorHandling 错误处理策略
// 便于使用, 是flag包中的ErrorHandling类型的别名
type ErrorHandling = flag.ErrorHandling

// ErrorHandling 错误处理策略常量
var (
	// ContinueOnError 解析错误时继续解析并返回错误
	ContinueOnError ErrorHandling = flag.ContinueOnError
	// ExitOnError 解析错误时退出程序
	ExitOnError ErrorHandling = flag.ExitOnError
	// PanicOnError 解析错误时触发panic
	PanicOnError ErrorHandling = flag.PanicOnError
)

// 标志类型
type FlagType int

const (
	FlagTypeUnknown     FlagType = iota // 未知类型
	FlagTypeInt                         // 整数类型
	FlagTypeInt64                       // 64位整数类型
	FlagTypeUint16                      // 16位无符号整数类型
	FlagTypeUint32                      // 32位无符号整数类型
	FlagTypeUint64                      // 64位无符号整数类型
	FlagTypeString                      // 字符串类型
	FlagTypeBool                        // 布尔类型
	FlagTypeFloat64                     // 64位浮点数类型
	FlagTypeEnum                        // 枚举类型
	FlagTypeDuration                    // 时间间隔类型
	FlagTypeTime                        // 时间类型
	FlagTypeMap                         // 映射类型
	FlagTypeStringSlice                 // 字符串切片类型
	FlagTypeIntSlice                    // []int 切片类型
	FlagTypeInt64Slice                  // []int64 切片类型
	FlagTypeSize                        // 大小类型
)

// 内置标志名称
var (
	HelpFlagName                = "help"       // 帮助标志名称
	HelpFlagShortName           = "h"          // 帮助标志短名称
	VersionFlagLongName         = "version"    // 版本标志名称
	VersionFlagShortName        = "v"          // 版本标志短名称
	CompletionShellFlagLongName = "completion" // 生成shell补全标志长名称
)

// 定义中英文的补全标志的使用说明
const (
	CompletionShellDescCN = "生成指定的 Shell 补全脚本, 可选类型: %v"
	CompletionShellDescEN = "Generate the specified Shell completion script, optional types: %v"
)

// 支持的Shell类型切片
var ShellSlice = []string{ShellNone, ShellBash, ShellPowershell, ShellPwsh}

// 支持的Shell类型
const (
	ShellBash       = "bash"       // bash shell
	ShellPowershell = "powershell" // powershell shell
	ShellPwsh       = "pwsh"       // pwsh shell
	ShellNone       = "none"       // 无shell
)

// 内置标志使用说明
var (
	HelpFlagUsage    = "Show help"    // 帮助标志使用说明
	VersionFlagUsage = "Show version" // 版本标志使用说明
)

// 定义标志的分隔符常量
const (
	// 逗号
	FlagSplitComma = ","

	// 分号
	FlagSplitSemicolon = ";"

	// 竖线
	FlagSplitPipe = "|"

	// 冒号
	FlagKVColon = ":"

	// 等号
	FlagKVEqual = "="
)

// Flag支持的标志分隔符切片
var FlagSplitSlice = []string{
	// 逗号
	FlagSplitComma,

	// 分号
	FlagSplitSemicolon,

	// 竖线
	FlagSplitPipe,

	// 冒号
	FlagKVColon,
}

// Validator 验证器接口, 所有自定义验证器需实现此接口
type Validator interface {
	// Validate 验证参数值是否合法
	//
	// 参数:
	//   - value: 需要验证的参数值
	//
	// 返回:
	//   - error: 验证失败时返回错误信息, 验证成功时返回nil
	Validate(value any) error
}

// Flag 所有标志类型的通用接口,定义了标志的元数据访问方法
type Flag interface {
	LongName() string   // 获取标志的长名称
	ShortName() string  // 获取标志的短名称
	Usage() string      // 获取标志的用法
	Type() FlagType     // 获取标志类型
	GetDefaultAny() any // 获取标志的默认值(any类型)
	String() string     // 获取标志的字符串表示
	IsSet() bool        // 判断标志是否已设置值
	Reset()             // 重置标志值为默认值
	GetEnvVar() string  // 获取标志绑定的环境变量名称
}

// TypedFlag 所有标志类型的通用接口,定义了标志的元数据访问方法和默认值访问方法
type TypedFlag[T any] interface {
	Flag                                 // 继承标志接口
	GetDefault() T                       // 获取标志的具体类型默认值
	Get() T                              // 获取标志的具体类型值
	GetPointer() *T                      // 获取标志值的指针
	Set(T) error                         // 设置标志的具体类型值
	SetValidator(Validator)              // 设置标志的验证器
	BindEnv(envName string) *BaseFlag[T] // 绑定环境变量
}

// FlagTypeToString 将FlagType转换为带语义信息的字符串
//
// 参数:
//   - flagType: 需要转换的FlagType枚举值
//
// 返回值:
//   - 带语义信息的类型字符串，用于命令行帮助信息显示
func FlagTypeToString(flagType FlagType) string {
	switch flagType {
	case FlagTypeInt: // 整数类型

		return "<int>"

	case FlagTypeInt64: // 64位整数类型
		return "<int64>"

	case FlagTypeUint16: // 16位无符号整数类型
		return "<0-65535>"

	case FlagTypeUint32: // 32位无符号整数类型
		return "<uint32>"

	case FlagTypeUint64: // 64位无符号整数类型
		return "<uint64>"

	case FlagTypeString: // 字符串类型
		return "<string>"

	case FlagTypeBool: // 布尔类型
		// 布尔类型特殊处理
		return ""

	case FlagTypeFloat64: // 64位浮点数类型
		return "<float64>"

	case FlagTypeEnum: // 枚举类型
		return "<enum>"

	case FlagTypeDuration: // 时间间隔类型
		return "<duration>"

	case FlagTypeTime: // 时间类型
		return "<time>"

	case FlagTypeMap: // 映射类型
		return "<k=v,k=v,...>"

	case FlagTypeStringSlice, FlagTypeIntSlice, FlagTypeInt64Slice: // 切片类型
		return "<value,value,...>"

	case FlagTypeSize: // 大小类型
		return "<size+unit>"

	default:
		return "<value>"
	}
}

// FormatDefaultValue 根据标志类型格式化默认值为人类可读的字符串
//
// 参数:
//   - flagType: 标志类型
//   - defaultValue: 默认值
//
// 返回值:
//   - 格式化后的默认值字符串
func FormatDefaultValue(flagType FlagType, defaultValue any) string {
	// 如果默认值为nil，返回空字符串
	if defaultValue == nil {
		return ""
	}

	switch flagType {
	case FlagTypeString:
		if str, ok := defaultValue.(string); ok {
			return str
		}
		return fmt.Sprintf("%v", defaultValue)

	case FlagTypeBool:
		if b, ok := defaultValue.(bool); ok {
			return fmt.Sprintf("%t", b)
		}
		return fmt.Sprintf("%v", defaultValue)

	case FlagTypeInt:
		if i, ok := defaultValue.(int); ok {
			return fmt.Sprintf("%d", i)
		}
		return fmt.Sprintf("%v", defaultValue)

	case FlagTypeInt64:
		if i, ok := defaultValue.(int64); ok {
			return fmt.Sprintf("%d", i)
		}
		return fmt.Sprintf("%v", defaultValue)

	case FlagTypeUint16:
		if i, ok := defaultValue.(uint16); ok {
			return fmt.Sprintf("%d", i)
		}
		return fmt.Sprintf("%v", defaultValue)

	case FlagTypeUint32:
		if i, ok := defaultValue.(uint32); ok {
			return fmt.Sprintf("%d", i)
		}
		return fmt.Sprintf("%v", defaultValue)

	case FlagTypeUint64:
		if i, ok := defaultValue.(uint64); ok {
			return fmt.Sprintf("%d", i)
		}
		return fmt.Sprintf("%v", defaultValue)

	case FlagTypeFloat64:
		if f, ok := defaultValue.(float64); ok {
			return fmt.Sprintf("%g", f)
		}
		return fmt.Sprintf("%v", defaultValue)

	case FlagTypeDuration:
		if d, ok := defaultValue.(time.Duration); ok {
			return d.String()
		}
		return fmt.Sprintf("%v", defaultValue)

	case FlagTypeTime:
		if t, ok := defaultValue.(time.Time); ok {
			// 根据时间的特点智能选择格式
			if t.IsZero() {
				return "" // 零值时间返回空字符串
			}

			// 如果时间只有日期部分（时分秒都是0），只显示日期
			if t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0 && t.Nanosecond() == 0 {
				return t.Format("2006-01-02")
			}

			// 默认显示到秒
			return t.Format("2006-01-02 15:04:05")
		}
		return fmt.Sprintf("%v", defaultValue)

	case FlagTypeMap:
		if m, ok := defaultValue.(map[string]string); ok {
			var entries []string
			for k, v := range m {
				entries = append(entries, fmt.Sprintf("%s=%s", k, v))
			}
			// 对键值对进行排序，确保输出一致
			sort.Strings(entries)
			return strings.Join(entries, ",")
		}
		return fmt.Sprintf("%v", defaultValue)

	case FlagTypeStringSlice:
		if slice, ok := defaultValue.([]string); ok {
			return strings.Join(slice, ",")
		}
		return fmt.Sprintf("%v", defaultValue)

	case FlagTypeIntSlice:
		if slice, ok := defaultValue.([]int); ok {
			var strValues []string
			for _, v := range slice {
				strValues = append(strValues, fmt.Sprintf("%d", v))
			}
			return strings.Join(strValues, ",")
		}
		return fmt.Sprintf("%v", defaultValue)

	case FlagTypeInt64Slice:
		if slice, ok := defaultValue.([]int64); ok {
			var strValues []string
			for _, v := range slice {
				strValues = append(strValues, fmt.Sprintf("%d", v))
			}
			return strings.Join(strValues, ",")
		}
		return fmt.Sprintf("%v", defaultValue)

	case FlagTypeSize:
		// Size类型通常是一个带有单位的数值，这里简单处理
		return fmt.Sprintf("%v", defaultValue)

	case FlagTypeEnum:
		// 枚举类型直接显示值
		return fmt.Sprintf("%v", defaultValue)

	default:
		// 未知类型使用默认格式化
		return fmt.Sprintf("%v", defaultValue)
	}
}
