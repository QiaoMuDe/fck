// Package flags 标志类型定义和接口
// 本文件定义了所有标志类型的通用接口和基础标志结构体，包括标志类型枚举、
// 验证器接口、标志接口等核心定义，为整个标志系统提供基础类型支持。
package flags

// 定义非法字符集常量, 防止非法的标志名称
const InvalidFlagChars = " !@#$%^&*(){}[]|\\;:'\"<>,.?/"

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
