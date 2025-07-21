// flags 定义了所有标志类型的通用接口和基础标志结构体
package flags

// 定义非法字符集常量, 防止非法的标志名称
const InvalidFlagChars = " !@#$%^&*(){}[]|\\;:'\"<>,.?/"

// 标志类型
type FlagType int

const (
	FlagTypeUnknown  FlagType = iota // 未知类型
	FlagTypeInt                      // 整数类型
	FlagTypeInt64                    // 64位整数类型
	FlagTypeUint16                   // 16位无符号整数类型
	FlagTypeUint32                   // 32位无符号整数类型
	FlagTypeUint64                   // 64位无符号整数类型
	FlagTypeString                   // 字符串类型
	FlagTypeBool                     // 布尔类型
	FlagTypeFloat64                  // 64位浮点数类型
	FlagTypeEnum                     // 枚举类型
	FlagTypeDuration                 // 时间间隔类型
	FlagTypeSlice                    // 切片类型
	FlagTypeTime                     // 时间类型
	FlagTypeMap                      // 映射类型
	FlagTypePath                     // 路径类型
	FlagTypeIP4                      // IPv4地址类型
	FlagTypeIP6                      // IPv6地址类型
	FlagTypeURL                      // URL类型
)

// 内置标志名称
var (
	HelpFlagName                 = "help"                      // 帮助标志名称
	HelpFlagShortName            = "h"                         // 帮助标志短名称
	VersionFlagLongName          = "version"                   // 版本标志名称
	VersionFlagShortName         = "v"                         // 版本标志短名称
	CompletionShellFlagLongName  = "generate-shell-completion" // 生成shell补全标志长名称
	CompletionShellFlagShortName = "gsc"                       // 生成shell补全标志短名称
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
	HelpFlagUsageEn    = "Show help"                       // 帮助标志英文使用说明
	VersionFlagUsageEn = "Show the version of the program" // 版本标志英文使用说明
	VersionFlagUsageZh = "显示程序的版本信息"                       // 版本标志中文使用说明
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
	// value: 待验证的参数值
	// 返回值: 验证通过返回nil, 否则返回错误信息
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

// FlagTypeToString 将FlagType转换为字符串
//
// 参数:
//   - flagType: 需要转换的FlagType枚举值
//   - withBrackets: 是否在返回的字符串中包含尖括号
//     如果为true且flagType为bool时返回空字符串
//
// 返回值:
//   - 对应的类型字符串，如果类型未知则返回"unknown"或"<unknown>"
func FlagTypeToString(flagType FlagType, withBrackets bool) string {
	var result string

	switch flagType {
	case FlagTypeInt:
		result = "int"
	case FlagTypeInt64:
		result = "int64"
	case FlagTypeUint16:
		result = "uint16"
	case FlagTypeUint32:
		result = "uint32"
	case FlagTypeUint64:
		result = "uint64"
	case FlagTypeString:
		result = "string"
	case FlagTypeBool:
		// 布尔类型特殊处理
		if withBrackets {
			return ""
		}
		return "bool"
	case FlagTypeFloat64:
		result = "float64"
	case FlagTypeEnum:
		result = "enum"
	case FlagTypeDuration:
		result = "duration"
	case FlagTypeTime:
		result = "time"
	case FlagTypeMap:
		result = "map"
	case FlagTypePath:
		result = "path"
	case FlagTypeSlice:
		result = "slice"
	case FlagTypeIP4:
		result = "ip4"
	case FlagTypeIP6:
		result = "ip6"
	case FlagTypeURL:
		result = "url"
	default:
		result = "unknown"
	}

	if withBrackets {
		return "<" + result + ">"
	}
	return result
}
