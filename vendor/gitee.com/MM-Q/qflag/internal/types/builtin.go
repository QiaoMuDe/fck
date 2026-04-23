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

	// InstallCompletionFlag 安装补全标志
	//
	// 安装补全标志用于自动生成补全脚本并配置到Shell。
	// 只在根命令中注册，需要启用补全功能。
	InstallCompletionFlag

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

	// InstallCompletionFlagName 安装补全标志名称
	InstallCompletionFlagName = "install-completion"

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

// __complete 子命令的指令常量定义
const (
	// InstructionFuzzy 模糊匹配指令
	// 用法: __complete fuzzy <模式> <候选1> [候选2] ...
	// 输出: 每行一个匹配结果（按匹配质量降序）
	InstructionFuzzy = "fuzzy"

	// InstructionContext 上下文计算指令
	// 用法: __complete context <arg0> [arg1] ...
	// 输出: 上下文路径，如 "/server/start/"
	InstructionContext = "context"

	// InstructionCandidates 候选选项获取指令
	// 用法: __complete candidates <context>
	// 输出: 空格分隔的候选选项列表
	InstructionCandidates = "candidates"

	// InstructionEnum 枚举值获取指令
	// 用法: __complete enum <context> <flag-name>
	// 输出: 空格分隔的枚举值列表
	InstructionEnum = "enum"

	// InstructionAll 统一补全指令
	// 用法: __complete all <cur> <prev> [cmd_args...]
	// 输出: 多行格式，包含 CONTEXT, CUR, PREV, CANDIDATES, ENUM, MATCHES, IS_FLAG
	InstructionAll = "all"
)

// CompleteCmdName 补全命令名称
const CompleteCmdName = "__complete"

// 补全安装相关常量
const (
	// CompletionsDirName 补全脚本存放目录名
	CompletionsDirName = ".qflag_completions"

	// CompletionScriptComment 补全脚本注释模板
	CompletionScriptComment = "# qflag completion for %s\n"

	// PwshCompletionScriptExt PowerShell 补全脚本扩展名
	PwshCompletionScriptExt = ".ps1"

	// BashCompletionScriptExt Bash 补全脚本扩展名
	BashCompletionScriptExt = ".sh"

	// PwshProfileDirWindows Windows PowerShell 配置文件目录
	PwshProfileDirWindows = "Documents/PowerShell"

	// PwshProfileFileName PowerShell 配置文件名
	PwshProfileFileName = "Microsoft.PowerShell_profile.ps1"

	// PwshProfileDirUnix Unix PowerShell 配置文件目录
	PwshProfileDirUnix = ".config/powershell"

	// BashProfileFileNameDarwin macOS Bash 配置文件名
	BashProfileFileNameDarwin = ".bash_profile"

	// BashProfileFileNameLinux Linux Bash 配置文件名
	BashProfileFileNameLinux = ".bashrc"
)

// 补全加载命令模板
const (
	// PwshLoadCommandTemplate PowerShell 加载命令模板
	// 参数: 程序名（3次）、脚本路径（1次）
	PwshLoadCommandTemplate = "$__qflag_comp_%s = '%s'; if (Test-Path $__qflag_comp_%s) { . $__qflag_comp_%s }"

	// BashLoadCommandTemplate Bash 加载命令模板
	// 参数: 脚本路径（2次）
	BashLoadCommandTemplate = "[ -f '%s' ] && source '%s'"
)

// 补全安装成功信息 - 中文
const (
	// InstallSuccessScriptPathCN 脚本安装路径提示（中文）
	InstallSuccessScriptPathCN = "✓ 补全脚本已安装: %s"

	// InstallSuccessProfilePathCN 配置文件路径提示（中文）
	InstallSuccessProfilePathCN = "✓ 加载命令已添加到: %s"

	// InstallSuccessHintCN 重启提示（中文）
	InstallSuccessHintCN = "\n请重启终端或运行以下命令启用补全:"
)

// 补全安装成功信息 - 英文
const (
	// InstallSuccessScriptPathEN 脚本安装路径提示（英文）
	InstallSuccessScriptPathEN = "✓ Completion script installed: %s"

	// InstallSuccessProfilePathEN 配置文件路径提示（英文）
	InstallSuccessProfilePathEN = "✓ Load command added to: %s"

	// InstallSuccessHintEN 重启提示（英文）
	InstallSuccessHintEN = "\nPlease restart your terminal or run the following command to enable completions:"
)

// 补全执行命令（Shell 命令本身不需要翻译）
const (
	// InstallSuccessBashCmd Bash 执行命令
	InstallSuccessBashCmd = "  source %s"

	// InstallSuccessPwshCmd PowerShell 执行命令
	InstallSuccessPwshCmd = "  . %s"
)

// 内置标志描述 - 中文
const (
	// HelpFlagDescCN 帮助标志描述（中文）
	HelpFlagDescCN = "显示帮助信息"

	// VersionFlagDescCN 版本标志描述（中文）
	VersionFlagDescCN = "显示版本信息"

	// CompletionFlagDescCN 补全标志描述（中文）
	CompletionFlagDescCN = "生成Shell自动补全脚本, 支持的Shell: %v"

	// InstallCompletionFlagDescCN 安装补全标志描述（中文）
	InstallCompletionFlagDescCN = "安装Shell自动补全脚本到系统, 支持的Shell: %v"
)

// 内置标志描述 - 英文
const (
	// HelpFlagDescEN 帮助标志描述（英文）
	HelpFlagDescEN = "Show help information"

	// VersionFlagDescEN 版本标志描述（英文）
	VersionFlagDescEN = "Show version information"

	// CompletionFlagDescEN 补全标志描述（英文）
	CompletionFlagDescEN = "Generate shell completion script. Supported shells: %v"

	// InstallCompletionFlagDescEN 安装补全标志描述（英文）
	InstallCompletionFlagDescEN = "Install shell completion script to system. Supported shells: %v"
)
