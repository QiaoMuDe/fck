// flags 定义了所有标志类型的通用接口和基础标志结构体
package flags

// 定义非法字符集常量, 防止非法的标志名称
const InvalidFlagChars = " !@#$%^&*(){}[]|\\;:'\"<>,.?/"

// 标志类型
type FlagType int

const (
	FlagTypeInt      FlagType = iota + 1 // 整数类型
	FlagTypeInt64                        // 64位整数类型
	FlagTypeUint16                       // 16位无符号整数类型
	FlagTypeUint32                       // 32位无符号整数类型
	FlagTypeUint64                       // 64位无符号整数类型
	FlagTypeString                       // 字符串类型
	FlagTypeBool                         // 布尔类型
	FlagTypeFloat64                      // 64位浮点数类型
	FlagTypeEnum                         // 枚举类型
	FlagTypeDuration                     // 时间间隔类型
	FlagTypeSlice                        // 切片类型
	FlagTypeTime                         // 时间类型
	FlagTypeMap                          // 映射类型
	FlagTypePath                         // 路径类型
	FlagTypeIP4                          // IPv4地址类型
	FlagTypeIP6                          // IPv6地址类型
	FlagTypeURL                          // URL类型
)

// 内置标志名称
var (
	HelpFlagName            = "help"    // 帮助标志名称
	HelpFlagShortName       = "h"       // 帮助标志短名称
	ShowInstallPathFlagName = "sip"     // 显示安装路径标志名称
	VersionFlagLongName     = "version" // 版本标志名称
	VersionFlagShortName    = "v"       // 版本标志短名称
)

// 内置标志使用说明
var (
	HelpFlagUsageEn            = "Show help information"                     // 帮助标志英文使用说明
	HelpFlagUsageZh            = "显示帮助信息"                                    // 帮助标志中文使用说明
	ShowInstallPathFlagUsageEn = "Show the installation path of the program" // 安装路径标志英文使用说明
	ShowInstallPathFlagUsageZh = "显示程序的安装路径"                                 // 安装路径标志中文使用说明
	VersionFlagUsageEn         = "Show the version of the program"           // 版本标志英文使用说明
	VersionFlagUsageZh         = "显示程序的版本信息"                                 // 版本标志中文使用说明
)

// 定义标志的分隔符常量
const (
	FlagSplitComma     = "," // 逗号
	FlagSplitSemicolon = ";" // 分号
	FlagSplitPipe      = "|" // 竖线
	FlagKVColon        = ":" // 冒号
	FlagKVEqual        = "=" // 等号
)

// Flag支持的标志分隔符切片
var FlagSplitSlice = []string{
	FlagSplitComma,     // 逗号
	FlagSplitSemicolon, // 分号
	FlagSplitPipe,      // 竖线
	FlagKVColon,        // 冒号
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
}

// TypedFlag 所有标志类型的通用接口,定义了标志的元数据访问方法和默认值访问方法
type TypedFlag[T any] interface {
	Flag                    // 继承标志接口
	GetDefault() T          // 获取标志的具体类型默认值
	Get() T                 // 获取标志的具体类型值
	GetPointer() *T         // 获取标志值的指针
	Set(T) error            // 设置标志的具体类型值
	SetValidator(Validator) // 设置标志的验证器
}

// FlagTypeToString 将FlagType转换为字符串
func FlagTypeToString(flagType FlagType) string {
	switch flagType {
	case FlagTypeInt:
		return "<int>"
	case FlagTypeInt64:
		return "<int64>"
	case FlagTypeUint16:
		return "<uint16>"
	case FlagTypeUint32:
		return "<uint32>"
	case FlagTypeUint64:
		return "<uint64>"
	case FlagTypeString:
		return "<string>"
	case FlagTypeBool:
		// 布尔类型没有参数类型字符串
		return ""
	case FlagTypeFloat64:
		return "<float64>"
	case FlagTypeEnum:
		return "<enum>"
	case FlagTypeDuration:
		return "<duration>"
	case FlagTypeTime:
		return "<time>"
	case FlagTypeMap:
		return "<map>"
	case FlagTypePath:
		return "<path>"
	case FlagTypeIP4:
		return "<ipv4>"
	case FlagTypeIP6:
		return "<ipv6>"
	case FlagTypeURL:
		return "<url>"
	default:
		return "<unknown>"
	}
}
