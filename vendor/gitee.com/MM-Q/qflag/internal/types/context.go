// Package types 命令上下文和状态管理
// 本文件定义了命令上下文结构体，用于管理命令的状态、子命令、标志注册表等信息，
// 提供命令执行过程中的状态维护和数据共享功能。
package types

import (
	"flag"
	"sync"
	"sync/atomic"

	"gitee.com/MM-Q/qflag/flags"
)

// CmdContext 命令上下文，包含所有必要的状态信息
// 这是所有函数操作的核心数据结构
type CmdContext struct {
	// 长命令名称
	LongName string
	// 短命令名称
	ShortName string

	// 标志注册表, 统一管理标志的元数据
	FlagRegistry *flags.FlagRegistry
	// 底层flag集合, 处理参数解析
	FlagSet *flag.FlagSet

	// 命令行参数(非标志参数)
	Args []string
	// 是否已经解析过参数
	Parsed atomic.Bool
	// 用于确保参数解析只执行一次
	ParseOnce sync.Once
	// 读写锁
	Mutex sync.RWMutex

	// 子命令上下文切片
	SubCmds []*CmdContext
	// 子命令映射表
	SubCmdMap map[string]*CmdContext
	// 父命令上下文
	Parent *CmdContext

	// 配置信息
	Config *CmdConfig

	// 内置标志结构体
	BuiltinFlags *BuiltinFlags

	// ParseHook 解析阶段钩子函数
	// 在标志解析完成后、子命令参数处理后调用
	//
	// 参数:
	//   - 当前命令上下文
	//
	// 返回值:
	//   - error: 错误信息, 非nil时会中断解析流程
	//   - bool: 是否需要退出程序
	ParseHook func(*CmdContext) (error, bool)
}

// NewCmdContext 创建新的命令上下文
//
// 参数:
//   - longName: 长命令名称
//   - shortName: 短命令名称
//   - errorHandling: 错误处理方式
//
// 返回值:
//   - *CmdContext: 新创建的命令上下文
//
// errorHandling可选参数:
//   - flag.ContinueOnError: 解析标志时遇到错误继续解析, 并返回错误信息
//   - flag.ExitOnError: 解析标志时遇到错误立即退出程序, 并返回错误信息
//   - flag.PanicOnError: 解析标志时遇到错误立即触发panic
func NewCmdContext(longName, shortName string, errorHandling flag.ErrorHandling) *CmdContext {
	if longName == "" && shortName == "" {
		panic("cmd long name and short name cannot both be empty")
	}

	// 优先使用长名称, 如果长名称为空则使用短名称
	cmdName := longName
	if cmdName == "" {
		cmdName = shortName
	}

	return &CmdContext{
		LongName:     longName,                                // 长名称
		ShortName:    shortName,                               // 短名称
		FlagSet:      flag.NewFlagSet(cmdName, errorHandling), // 创建新的flag集
		FlagRegistry: flags.NewFlagRegistry(),                 // 创建新的标志注册表
		Args:         []string{},                              // 命令行参数
		SubCmds:      []*CmdContext{},                         // 子命令上下文切片
		SubCmdMap:    make(map[string]*CmdContext),            // 子命令映射表
		Config:       NewCmdConfig(),                          // 创建新的命令配置
		BuiltinFlags: NewBuiltinFlags(),                       // 创建新的内置标志结构体
	}
}

// GetName 获取命令名称
// 如果长命令名称不为空则返回长命令名称, 否则返回短命令名称
//
// 返回值:
//   - string: 命令名称
func (ctx *CmdContext) GetName() string {
	if ctx.LongName != "" {
		return ctx.LongName
	}
	return ctx.ShortName
}
