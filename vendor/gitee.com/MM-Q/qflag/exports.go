// package qflag 根包统一导出入口
// 本文件用于将各子包的核心功能导出到根包，简化外部使用。
// 通过类型别名和变量导出的方式，为用户提供统一的API接口。
package qflag

import (
	"gitee.com/MM-Q/qflag/flags"
	"gitee.com/MM-Q/qflag/internal/types"
)

/*
项目地址: https://gitee.com/MM-Q/qflag
*/

// ErrorHandling 错误处理策略
type ErrorHandling = flags.ErrorHandling

// ErrorHandling 错误处理策略常量
var (
	// ContinueOnError 解析错误时继续解析并返回错误
	ContinueOnError ErrorHandling = flags.ContinueOnError
	// ExitOnError 解析错误时退出程序
	ExitOnError ErrorHandling = flags.ExitOnError
	// PanicOnError 解析错误时触发panic
	PanicOnError ErrorHandling = flags.PanicOnError
)

// CmdConfig 导出cmd包中的CmdConfig结构体
type CmdConfig = types.CmdConfig

// ExampleInfo 导出示例信息类型
type ExampleInfo = types.ExampleInfo

// Flag 导出flag包中的Flag结构体
type Flag = flags.Flag

// FlagRegistry 导出flag包中的FlagRegistry结构体
type FlagRegistry = flags.FlagRegistry

// StringFlag 导出flag包中的StringFlag结构体
type StringFlag = flags.StringFlag

// IntFlag 导出flag包中的IntFlag结构体
type IntFlag = flags.IntFlag

// BoolFlag 导出flag包中的BoolFlag结构体
type BoolFlag = flags.BoolFlag

// DurationFlag 导出flag包中的DurationFlag结构体
type DurationFlag = flags.DurationFlag

// Float64Flag 导出flag包中的Float64Flag结构体
type Float64Flag = flags.Float64Flag

// Int64Flag 导出flag包中的Int64Flag结构体
type Int64Flag = flags.Int64Flag

// StringSliceFlag 导出flag包中的StringSliceFlag结构体
type StringSliceFlag = flags.StringSliceFlag

// IntSliceFlag 导出flag包中的IntSliceFlag结构体
type IntSliceFlag = flags.IntSliceFlag

// Int64SliceFlag 导出flag包中的Int64SliceFlag结构体
type Int64SliceFlag = flags.Int64SliceFlag

// EnumFlag 导出flag包中的EnumFlag结构体
type EnumFlag = flags.EnumFlag

// MapFlag 导出flag包中的MapFlag结构体
type MapFlag = flags.MapFlag

// TimeFlag 导出flag包中的TimeFlag结构体
type TimeFlag = flags.TimeFlag

// Uint16Flag 导出flag包中的UintFlag结构体
type Uint16Flag = flags.Uint16Flag

// Uint32Flag 导出flag包中的Uint32Flag结构体
type Uint32Flag = flags.Uint32Flag

// Uint64Flag 导出flag包中的Uint64Flag结构体
type Uint64Flag = flags.Uint64Flag

// SizeFlag 导出flag包中的SizeFlag结构体
type SizeFlag = flags.SizeFlag
