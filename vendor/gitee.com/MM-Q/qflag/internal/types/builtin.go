package types

import "runtime"

// BuiltinFlagType 内置标志类型
//
// 内置标志是系统自动处理的特殊标志, 如帮助标志和版本标志。
// 这些标志在解析完成后会被自动检查, 如果被设置则执行相应的操作。
type BuiltinFlagType int

const (
	// HelpFlag 帮助标志
	//
	// 帮助标志用于显示命令的帮助信息, 包括用法、选项、子命令等。
	// 总是会被注册, 因为所有命令都应该支持帮助功能。
	HelpFlag BuiltinFlagType = iota

	// VersionFlag 版本标志
	//
	// 版本标志用于显示命令的版本信息。
	// 只有在命令设置了版本信息时才会被注册。
	VersionFlag

	// CompletionFlag 补全标志
	//
	// 补全标志用于生成Shell自动补全脚本。
	// 总是会被注册, 支持bash和pwsh两种Shell类型。
	CompletionFlag

	// 可以继续添加其他内置标志
	// 例如: ConfigFlag, VerboseFlag 等
)

// 内置标志名称常量
const (
	// HelpFlagName 帮助标志名称
	HelpFlagName = "help"

	// HelpFlagShortName 帮助标志短名称
	HelpFlagShortName = "h"

	// VersionFlagName 版本标志名称
	VersionFlagName = "version"

	// VersionFlagShortName 版本标志短名称
	VersionFlagShortName = "v"

	// CompletionFlagName 补全标志名称
	CompletionFlagName = "completion"

	// // CompletionFlagShortName 补全标志短名称
	// CompletionFlagShortName = "c"
)

// BuiltinFlagHandler 内置标志处理器接口
//
// 内置标志处理器负责处理特定类型的内置标志。
// 每种内置标志类型都应该有一个对应的处理器实现。
type BuiltinFlagHandler interface {
	// Handle 处理内置标志
	//
	// 参数:
	//   - cmd: 要处理的命令
	//
	// 返回值:
	//   - error: 处理失败时返回错误
	//
	// 功能说明:
	//   - 执行内置标志的特定操作
	//   - 通常在执行后会调用 os.Exit 退出程序
	//   - 例如: 帮助标志会打印帮助信息并退出
	Handle(cmd Command) error

	// Type 返回标志类型
	//
	// 返回值:
	//   - BuiltinFlagType: 标志类型
	//
	// 功能说明:
	//   - 用于标识处理器处理的标志类型
	//   - 在注册和管理时使用
	Type() BuiltinFlagType

	// ShouldRegister 判断是否应该注册此标志
	//
	// 参数:
	//   - cmd: 要检查的命令
	//
	// 返回值:
	//   - bool: 是否应该注册
	//
	// 功能说明:
	//   - 根据命令的配置决定是否注册此标志
	//   - 例如: 版本标志只有在设置了版本信息时才注册
	//   - 帮助标志总是注册
	ShouldRegister(cmd Command) bool

	// ShouldSkipRegistration 判断是否应该跳过注册
	//
	// 参数:
	//   - cmd: 要检查的命令
	//
	// 返回值:
	//   - bool: 是否应该跳过注册
	//
	// 功能说明:
	//   - 检查标志是否已经被注册（避免重复注册）
	//   - 支持重复解析场景
	//   - 由每个处理器自己实现检查逻辑
	ShouldSkipRegistration(cmd Command) bool
}

const (
	// BashShell bash shell
	BashShell = "bash"

	// PwshShell pwsh shell
	PwshShell = "pwsh"

	// PowershellShell powershell shell
	PowershellShell = "powershell"
)

// Shell切片, 用于存储支持的Shell类型
var SupportedShells = []string{
	BashShell,
	PwshShell,
	PowershellShell,
}

// CurrentShell 根据当前平台返回对应的Shell类型
//
// 返回值:
//   - string: 当前Shell类型
func CurrentShell() string {
	if runtime.GOOS == "windows" {
		return PwshShell
	}
	return BashShell
}

// 补全脚本生成相关常量定义
const (
	// DefaultFlagParamsCapacity 预估的标志参数初始容量
	// 基于常见CLI工具分析, 大多数工具的标志数量在100-500之间
	DefaultFlagParamsCapacity = 256

	// NamesPerItem 每个标志/命令的名称数量(长名+短名)
	NamesPerItem = 2

	// MaxTraverseDepth 命令树遍历的最大深度限制
	// 防止循环引用导致的无限递归, 一般CLI工具很少超过20层
	MaxTraverseDepth = 50
)
