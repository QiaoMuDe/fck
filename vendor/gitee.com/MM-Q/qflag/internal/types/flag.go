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
	"flag"
	"fmt"
)

// ErrorHandling 错误处理方式枚举
//
// ErrorHandling 定义了解析错误时的处理策略, 直接使用标准库
// flag包的错误处理方式, 保持兼容性。
//
// 可选值:
//   - ContinueOnError: 解析错误时继续解析并返回错误
//   - ExitOnError: 解析错误时退出程序
//   - PanicOnError: 解析错误时触发panic
type ErrorHandling = flag.ErrorHandling

var (
	// ContinueOnError 解析错误时继续解析并返回错误
	//
	// 使用场景:
	//   - 需要收集所有错误
	//   - 自定义错误处理逻辑
	//   - 交互式应用
	ContinueOnError ErrorHandling = flag.ContinueOnError

	// ExitOnError 解析错误时退出程序
	//
	// 使用场景:
	//   - 简单命令行工具
	//   - 错误即致命的应用
	//   - 脚本和自动化工具
	ExitOnError ErrorHandling = flag.ExitOnError

	// PanicOnError 解析错误时触发panic
	//
	// 使用场景:
	//   - 开发和测试环境
	//   - 需要快速失败的场景
	//   - 调试和诊断
	PanicOnError ErrorHandling = flag.PanicOnError
)

// FlagType 标志类型枚举
//
// FlagType 定义了所有支持的标志类型, 用于类型识别和
// 特定处理逻辑的实现。
//
// 设计原则:
//   - 每种类型对应一种数据格式
//   - 支持基础类型和复合类型
//   - 便于类型检查和转换
type FlagType int

const (
	FlagTypeUnknown FlagType = iota // 未知标志类型, 用于错误处理

	// 基础类型
	FlagTypeString  // 字符串标志, 存储任意文本
	FlagTypeInt     // 整数标志, 平台相关int类型
	FlagTypeInt64   // 64位整数标志, 固定64位整数
	FlagTypeUint    // 无符号整数标志, 平台相关uint类型
	FlagTypeUint8   // 8位无符号整数标志, 0-255
	FlagTypeUint16  // 16位无符号整数标志, 0-65535
	FlagTypeUint32  // 32位无符号整数标志, 0-4294967295
	FlagTypeUint64  // 64位无符号整数标志, 0-18446744073709551615
	FlagTypeFloat64 // 64位浮点数标志, IEEE 754双精度
	FlagTypeBool    // 布尔标志, true/false值

	// 特殊类型
	FlagTypeEnum // 枚举标志, 限制为预定义值集合

	// 时间和大小类型
	FlagTypeDuration // 持续时间标志, 支持时间单位解析
	FlagTypeTime     // 时间标志, 支持多种时间格式
	FlagTypeSize     // 大小标志, 支持存储单位解析

	// 集合类型
	FlagTypeMap         // 映射标志, 键值对集合
	FlagTypeStringSlice // 字符串切片标志, 字符串数组
	FlagTypeIntSlice    // 整数切片标志, 整数数组
	FlagTypeInt64Slice  // 64位整数切片标志, 64位整数数组
)

// String 返回标志类型的字符串表示
//
// 返回值:
//   - string: 类型的可读字符串表示
//
// 功能说明:
//   - 提供人类可读的类型名称
//   - 用于错误消息和日志
//   - 未知类型返回格式化字符串
//
// 示例:
//   - FlagTypeString -> "string"
//   - FlagTypeIntSlice -> "[]int"
//   - FlagType(999) -> "FlagType(999)"
func (t FlagType) String() string {
	switch t {
	case FlagTypeUnknown:
		return "unknown"
	case FlagTypeString:
		return "string"
	case FlagTypeInt:
		return "int"
	case FlagTypeInt64:
		return "int64"
	case FlagTypeUint:
		return "uint"
	case FlagTypeUint8:
		return "uint8"
	case FlagTypeUint16:
		return "uint16"
	case FlagTypeUint32:
		return "uint32"
	case FlagTypeUint64:
		return "uint64"
	case FlagTypeFloat64:
		return "float64"
	case FlagTypeBool:
		return "bool"
	case FlagTypeEnum:
		return "enum"
	case FlagTypeDuration:
		return "duration"
	case FlagTypeTime:
		return "time"
	case FlagTypeMap:
		return "map"
	case FlagTypeStringSlice:
		return "[]string"
	case FlagTypeIntSlice:
		return "[]int"
	case FlagTypeInt64Slice:
		return "[]int64"
	case FlagTypeSize:
		return "size"
	default:
		return fmt.Sprintf("FlagType(%d)", t)
	}
}

// Validator 验证器函数类型
//
// Validator 是一个泛型函数类型，用于验证标志值的有效性。
// 验证器接收一个类型为 T 的值，返回错误信息。
//
// 参数:
//   - value: 要验证的值
//
// 返回值:
//   - error: 验证失败时返回错误，验证通过返回 nil
//
// 功能说明:
//   - 验证器在标志的 Set 方法中被调用
//   - 在解析完值后、设置值之前执行验证
//   - 如果验证失败，Set 方法会返回错误，标志值不会被设置
//   - 验证器是可选的，未设置时跳过验证
//   - 重复设置验证器会覆盖之前的验证器
//
// 空值处理:
//   - StringFlag: 空字符串不经过验证器，直接设置
//   - BoolFlag: 不经过验证器（无空值概念）
//   - 集合类型 (MapFlag, StringSliceFlag, IntSliceFlag, Int64SliceFlag): 空字符串不经过验证器，创建空集合
//   - 其他类型: 空字符串直接返回错误，不经过验证器
//
// 使用示例:
//   - port.SetValidator(func(value int) error {
//     if value < 1 || value > 65535 {
//     return fmt.Errorf("端口 %d 超出范围 [1, 65535]", value)
//     }
//     return nil
//     })
type Validator[T any] func(value T) error

// IsValid 检查标志类型是否有效
//
// 返回值:
//   - bool: 是否有效, true表示有效
//
// 功能说明:
//   - 排除未知类型
//   - 用于类型验证
//   - 确保类型在预定义范围内
func (t FlagType) IsValid() bool {
	return t != FlagTypeUnknown
}

// IsSliceType 检查是否为切片类型
//
// 返回值:
//   - bool: 是否为切片类型, true表示是
//
// 功能说明:
//   - 识别所有切片类型的标志
//   - 用于特殊处理逻辑
//   - 支持多值输入的标志
func (t FlagType) IsSliceType() bool {
	return t == FlagTypeStringSlice || t == FlagTypeIntSlice || t == FlagTypeInt64Slice
}

// IsNumericType 检查是否为数值类型
//
// 返回值:
//   - bool: 是否为数值类型, true表示是
//
// 功能说明:
//   - 识别所有数值类型的标志
//   - 包括整数、浮点数和大小类型
//   - 用于数值范围验证
func (t FlagType) IsNumericType() bool {
	switch t {
	case FlagTypeInt, FlagTypeInt64, FlagTypeUint, FlagTypeUint16, FlagTypeUint32, FlagTypeUint64,
		FlagTypeFloat64, FlagTypeSize:
		return true
	default:
		return false
	}
}

// Flag 接口定义了标志的核心行为
//
// Flag 是所有标志类型必须实现的基础接口, 定义了标志的
// 基本操作和属性。所有具体标志类型都应实现此接口。
//
// 设计原则:
//   - 提供统一的标志操作接口
//   - 支持多种数据类型
//   - 支持验证和环境变量绑定
//   - 提供完整的生命周期管理
type Flag interface {
	// Name 获取标志名称
	//
	// 返回值:
	//   - string: 标志的完整名称
	//
	// 功能说明:
	//   - 返回标志的主要标识符
	//   - 用于命令行参数和查找
	//   - 名称在注册表中必须唯一
	Name() string

	// LongName 获取标志长名称
	//
	// 返回值:
	//   - string: 标志的长名称
	//
	// 功能说明:
	//   - 与Name方法功能相同
	//   - 提供语义明确的方法名
	//   - 保持接口一致性
	LongName() string

	// ShortName 获取标志短名称
	//
	// 返回值:
	//   - string: 标志的短名称
	//
	// 功能说明:
	//   - 返回标志的简短形式
	//   - 通常为单个字符
	//   - 可能为空字符串
	ShortName() string

	// Desc 获取标志描述
	//
	// 返回值:
	//   - string: 标志的描述信息
	//
	// 功能说明:
	//   - 返回标志的功能说明
	//   - 用于帮助信息生成
	//   - 应简洁明了地描述标志用途
	Desc() string

	// Type 获取标志类型
	//
	// 返回值:
	//   - FlagType: 标志的类型枚举
	//
	// 功能说明:
	//   - 返回标志的数据类型
	//   - 用于类型检查和转换
	//   - 决定解析和验证逻辑
	Type() FlagType

	// Set 设置标志值
	//
	// 参数:
	//   - value: 要设置的字符串值
	//
	// 返回值:
	//   - error: 设置失败时返回错误
	//
	// 功能说明:
	//   - 从字符串解析并设置值
	//   - 自动进行类型转换
	//   - 更新内部状态和标记
	Set(value string) error

	// GetDef 获取默认值
	//
	// 返回值:
	//   - any: 标志的默认值
	//
	// 功能说明:
	//   - 返回初始化时设置的默认值
	//   - 用于帮助信息显示
	//   - 用户未设置值时使用此值
	GetDef() any

	// GetStr 获取标志当前值的字符串表示
	//
	// 返回值:
	//   - string: 标志当前值的字符串表示
	//
	// 功能说明:
	//   - 获取标志当前值的字符串表示
	//   - 与String()方法不同, 此方法专注于值本身
	//   - 用于内置标志处理中获取标志值
	GetStr() string

	// IsSet 检查标志是否被用户设置
	//
	// 返回值:
	//   - bool: 是否被设置, true表示已设置
	//
	// 功能说明:
	//   - 区分默认值和用户设置值
	//   - 用于必填标志检查
	//   - 影响某些标志的行为逻辑
	IsSet() bool

	// Reset 重置标志到默认状态
	//
	// 功能说明:
	//   - 清除用户设置的值
	//   - 恢复到默认值
	//   - 重置设置状态标记
	Reset()

	// String 获取标志值的字符串表示
	//
	// 返回值:
	//   - string: 值的字符串表示
	//
	// 功能说明:
	//   - 实现fmt.Stringer接口
	//   - 用于显示和日志输出
	//   - 格式应与输入格式兼容
	String() string

	// BindEnv 绑定环境变量
	//
	// 参数:
	//   - name: 环境变量名称
	//
	// 功能说明:
	//   - 从环境变量读取默认值
	//   - 优先级低于命令行参数
	//   - 支持配置文件和环境变量
	BindEnv(name string)

	// GetEnvVar 获取绑定的环境变量名
	//
	// 返回值:
	//   - string: 环境变量名称
	//
	// 功能说明:
	//   - 返回当前绑定的环境变量
	//   - 未绑定时返回空字符串
	//   - 用于调试和配置管理
	GetEnvVar() string

	// EnumValues 获取枚举类型的可选值
	//
	// 返回值:
	//   - []string: 枚举类型的可选值列表
	//
	// 功能说明:
	//   - 非枚举类型返回空切片
	//   - 枚举类型返回所有可选值
	//   - 用于补全脚本生成和验证
	EnumValues() []string
}
