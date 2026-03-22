package qflag

// 导出公共接口和类型
import (
	"gitee.com/MM-Q/qflag/internal/cmd"
	"gitee.com/MM-Q/qflag/internal/completion"
	"gitee.com/MM-Q/qflag/internal/flag"
	"gitee.com/MM-Q/qflag/internal/types"
)

// Command 定义了命令行工具中命令的基本接口, 包括标志管理、参数解析、子命令管理等功能
// 实现此接口的类型可以作为命令行工具的命令使用
type Command = types.Command

// MutexGroup 定义了一组互斥的标志, 其中最多只能有一个被设置
// 当用户设置了互斥组中的多个标志时, 解析器会返回错误
type MutexGroup = types.MutexGroup

// RequiredGroup 定义了一组必需的标志，其中所有标志都必须被设置
// 当用户没有设置必需组中的某些标志时，解析器会返回错误
type RequiredGroup = types.RequiredGroup

// CmdConfig 包含了命令的各种配置选项, 用于自定义命令的行为和外观
// 这些配置会影响命令的帮助信息显示、环境变量处理、错误提示等
type CmdConfig = types.CmdConfig

// Error 是qflag项目的标准错误类型, 提供了结构化的错误信息。
// 包含错误码、错误消息和原始错误, 便于错误分类和处理。
//
// 字段说明:
//   - Code: 错误码, 用于错误分类和程序化处理
//   - Message: 错误消息, 面向用户的描述信息
//   - Cause: 原始错误, 包装的底层错误
//
// 特性:
//   - 实现error接口
//   - 支持错误链 (errors.Unwrap)
//   - 支持错误比较 (errors.Is)
//   - 提供错误码匹配
type Error = types.Error

// NewError 创建新的错误
//
// 参数:
//   - code: 错误码, 用于错误分类和识别
//   - message: 错误消息, 面向用户的描述信息
//   - cause: 原始错误, 可以为nil
//
// 返回值:
//   - *Error: 新创建的错误实例
//
// 功能说明:
//   - 创建结构化的错误实例
//   - 保留原始错误信息
//   - 提供错误分类能力
var NewError = types.NewError

// WrapError 包装错误
//
// 参数:
//   - err: 要包装的原始错误
//   - code: 新的错误码
//   - message: 新的错误消息
//
// 返回值:
//   - *Error: 包装后的错误
//
// 功能说明:
//   - 为现有错误添加上下文信息
//   - 保持原始错误链
//   - 提供新的错误分类
//
// 使用场景:
//   - 为底层错误添加业务上下文
//   - 统一错误处理格式
//   - 错误转换和适配
var WrapError = types.WrapError

// WrapParseError 包装解析错误, 专门用于标志解析场景
//
// 参数:
//   - err: 原始解析错误
//   - flagType: 标志类型描述
//   - value: 解析失败的值
//
// 返回值:
//   - *Error: 包装后的解析错误
//
// 功能说明:
//   - 专门用于标志解析错误
//   - 自动生成描述性错误消息
//   - 保留原始错误信息
//
// 使用场景:
//   - 标志值解析失败
//   - 类型转换错误
//   - 格式验证错误
var WrapParseError = types.WrapParseError

// IsNotFoundError 判断是否为"未找到"错误
//
// 参数:
//   - err: 要检查的错误
//
// 返回值:
//   - bool: 是否为未找到错误, true表示是
//
// 功能说明:
//   - 检查错误码是否为FLAG_NOT_FOUND或COMMAND_NOT_FOUND
//   - 支持错误链检查
//   - 便于统一处理未找到类型的错误
//
// 使用场景:
//   - 统一处理资源不存在的情况
//   - 区分未找到错误和其他错误
//   - 简化错误处理逻辑
var IsNotFoundError = types.IsNotFoundError

// 预定义错误变量
var (
	// ErrInvalidFlagType 无效的标志类型错误
	//
	// 使用场景:
	//   - 传入不支持的标志类型
	//   - 标志类型转换失败
	ErrInvalidFlagType = types.ErrInvalidFlagType

	// ErrFlagNotFound 标志不存在错误
	//
	// 使用场景:
	//   - 查找不存在的标志
	//   - 引用未注册的标志
	ErrFlagNotFound = types.ErrFlagNotFound

	// ErrCmdNotFound 命令不存在错误
	//
	// 使用场景:
	//   - 查找不存在的命令
	//   - 引用未注册的命令
	ErrCmdNotFound = types.ErrCmdNotFound

	// ErrFlagAlreadyExists 标志已存在错误
	//
	// 使用场景:
	//   - 注册重复的标志
	//   - 标志名称冲突
	ErrFlagAlreadyExists = types.ErrFlagAlreadyExists

	// ErrCmdAlreadyExists 命令已存在错误
	//
	// 使用场景:
	//   - 注册重复的命令
	//   - 命令名称冲突
	ErrCmdAlreadyExists = types.ErrCmdAlreadyExists

	// ErrParseFailed 解析失败错误
	//
	// 使用场景:
	//   - 命令行参数解析失败
	//   - 配置文件解析失败
	ErrParseFailed = types.ErrParseFailed

	// ErrValidationFailed 验证失败错误
	//
	// 使用场景:
	//   - 标志值验证失败
	//   - 业务规则验证失败
	ErrValidationFailed = types.ErrValidationFailed

	// ErrRequiredFlag 必填标志缺失错误
	//
	// 使用场景:
	//   - 必填标志未提供
	//   - 必填标志值为空
	ErrRequiredFlag = types.ErrRequiredFlag

	// ErrInvalidValue 无效值错误
	//
	// 使用场景:
	//   - 标志值格式错误
	//   - 标志值超出范围
	ErrInvalidValue = types.ErrInvalidValue
)

// ErrorHandling 错误处理方式枚举
// ErrorHandling 定义了解析错误时的处理策略, 直接使用标准库
// flag包的错误处理方式, 保持兼容性。
//
// 可选值:
//   - ContinueOnError: 解析错误时继续解析并返回错误
//   - ExitOnError: 解析错误时退出程序
//   - PanicOnError: 解析错误时触发panic
type ErrorHandling = types.ErrorHandling

// 错误处理策略常量
var (
	// ContinueOnError 解析错误时继续解析并返回错误
	//
	// 使用场景:
	//   - 需要收集所有错误
	//   - 自定义错误处理逻辑
	//   - 交互式应用
	ContinueOnError = types.ContinueOnError

	// ExitOnError 解析错误时退出程序
	//
	// 使用场景:
	//   - 简单命令行工具
	//   - 错误即致命的应用
	//   - 脚本和自动化工具
	ExitOnError = types.ExitOnError

	// PanicOnError 解析错误时触发panic
	//
	// 使用场景:
	//   - 开发和测试环境
	//   - 需要快速失败的场景
	//   - 调试和诊断
	PanicOnError = types.PanicOnError
)

// FlagType 标志类型枚举
// FlagType 定义了所有支持的标志类型, 用于类型识别和
// 特定处理逻辑的实现。
//
// 设计原则:
//   - 每种类型对应一种数据格式
//   - 支持基础类型和复合类型
//   - 便于类型检查和转换
type FlagType = types.FlagType

// 标志类型常量
const (
	FlagTypeUnknown FlagType = types.FlagTypeUnknown // 未知标志类型, 用于错误处理

	// 基础类型
	FlagTypeString  FlagType = types.FlagTypeString  // 字符串标志, 存储任意文本
	FlagTypeInt     FlagType = types.FlagTypeInt     // 整数标志, 平台相关int类型
	FlagTypeInt64   FlagType = types.FlagTypeInt64   // 64位整数标志, 固定64位整数
	FlagTypeUint    FlagType = types.FlagTypeUint    // 无符号整数标志, 平台相关uint类型
	FlagTypeUint8   FlagType = types.FlagTypeUint8   // 8位无符号整数标志, 0-255
	FlagTypeUint16  FlagType = types.FlagTypeUint16  // 16位无符号整数标志, 0-65535
	FlagTypeUint32  FlagType = types.FlagTypeUint32  // 32位无符号整数标志, 0-4294967295
	FlagTypeUint64  FlagType = types.FlagTypeUint64  // 64位无符号整数标志, 0-18446744073709551615
	FlagTypeFloat64 FlagType = types.FlagTypeFloat64 // 64位浮点数标志, IEEE 754双精度
	FlagTypeBool    FlagType = types.FlagTypeBool    // 布尔标志, true/false值

	// 特殊类型
	FlagTypeEnum FlagType = types.FlagTypeEnum // 枚举标志, 限制为预定义值集合

	// 时间和大小类型
	FlagTypeDuration FlagType = types.FlagTypeDuration // 持续时间标志, 支持时间单位解析
	FlagTypeTime     FlagType = types.FlagTypeTime     // 时间标志, 支持多种时间格式
	FlagTypeSize     FlagType = types.FlagTypeSize     // 大小标志, 支持存储单位解析

	// 集合类型
	FlagTypeMap         FlagType = types.FlagTypeMap         // 映射标志, 键值对集合
	FlagTypeStringSlice FlagType = types.FlagTypeStringSlice // 字符串切片标志, 字符串数组
	FlagTypeIntSlice    FlagType = types.FlagTypeIntSlice    // 整数切片标志, 整数数组
	FlagTypeInt64Slice  FlagType = types.FlagTypeInt64Slice  // 64位整数切片标志, 64位整数数组
)

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
//
//	// 端口号验证：1-65535
//	port.SetValidator(func(value int) error {
//	    if value < 1 || value > 65535 {
//	        return fmt.Errorf("端口 %d 超出范围 [1, 65535]", value)
//	    }
//	    return nil
//	})
//
//	// 字符串长度验证：3-20个字符
//	username.SetValidator(func(value string) error {
//	    if len(value) < 3 || len(value) > 20 {
//	        return fmt.Errorf("用户名长度 %d 超出范围 [3, 20]", len(value))
//	    }
//	    return nil
//	})
//
//	// 邮箱格式验证
//	email.SetValidator(func(value string) error {
//	    if !isValidEmail(value) {
//	        return fmt.Errorf("邮箱格式无效: %s", value)
//	    }
//	    return nil
//	})
//
// 注意事项:
//   - 验证器应该快速执行，避免耗时操作
//   - 验证器返回的错误应该清晰描述失败原因
//   - 验证器执行时已经持有锁，验证器本身不需要处理并发
//   - 验证器可以随时通过 ClearValidator 清除
type Validator[T any] = types.Validator[T]

// Flag 接口定义了标志的核心行为
// Flag 是所有标志类型必须实现的基础接口, 定义了标志的
// 基本操作和属性。所有具体标志类型都应实现此接口。
//
// 设计原则:
//   - 提供统一的标志操作接口
//   - 支持多种数据类型
//   - 支持验证和环境变量绑定
//   - 提供完整的生命周期管理
type Flag = types.Flag

// Parser 解析器接口
// Parser 定义了命令行参数解析的标准接口, 提供了不同层次的
// 解析功能, 从简单的参数解析到完整的命令路由执行。
//
// 设计理念:
//   - 分层设计: 提供不同层次的解析功能
//   - 灵活性: 支持仅解析、解析+路由等多种使用模式
//   - 可扩展性: 接口设计允许不同的解析策略实现
//
// 使用场景:
//   - 命令行工具的参数解析
//   - 子命令系统的路由管理
//   - 配置管理和参数验证
type Parser = types.Parser

// FlagRegistry 标志注册表接口
// FlagRegistry 定义了标志注册和管理的标准接口, 提供了
// 标志的完整生命周期管理功能。
//
// 核心功能:
//   - 标志的注册和注销
//   - 基于名称的查找和检索
//   - 批量操作和遍历支持
//   - 存在性检查和计数
//
// 设计特点:
//   - 支持长名称和短名称查找
//   - 提供统一的错误处理
//   - 支持别名管理 (通过具体实现)
//   - 线程安全由具体实现保证
type FlagRegistry = types.FlagRegistry

// CmdRegistry 命令注册表接口
// CmdRegistry 定义了命令注册和管理的标准接口, 提供了
// 命令的完整生命周期管理功能。
//
// 核心功能:
//   - 命令的注册和注销
//   - 基于名称的查找和检索
//   - 批量操作和遍历支持
//   - 存在性检查和计数
//
// 设计特点:
//   - 支持长名称和短名称查找
//   - 提供统一的错误处理
//   - 支持别名管理 (通过具体实现)
//   - 线程安全由具体实现保证
type CmdRegistry = types.CmdRegistry

// Cmd 命令结构体类型
// Cmd 是一个命令结构体, 实现了 types.Command 接口
// 提供了完整的命令行命令实现, 支持标志管理、子命令、
// 参数解析和执行等功能。使用读写锁保证并发安全。
type Cmd = cmd.Cmd

// NewCmd 创建新的命令实例
//
// 参数:
//   - longName: 命令的长名称
//   - shortName: 命令的短名称
//   - errorHandling: 错误处理策略
//
// 返回值:
//   - *Cmd: 初始化完成的命令实例
//
// 功能说明:
//   - 创建命令并初始化基本字段
//   - 创建标志和子命令注册器
//   - 设置默认解析器
//   - 初始化配置选项
var NewCmd = cmd.NewCmd

// CmdSpec 命令规格结构体
// CmdSpec 提供了通过规格创建命令的方式, 包含命令的所有属性。
// 这种方式比函数式配置更加直观和集中。
type CmdSpec = cmd.CmdSpec

// CmdOpts 命令选项结构体
// CmdOpts 提供了配置现有命令的方式，包含命令的所有可配置属性。
// 与 CmdSpec 不同，CmdOpts 用于配置已存在的命令，而不是创建新命令。
type CmdOpts = cmd.CmdOpts

// NewCmdSpec 创建新的命令规格
//
// 参数:
//   - longName: 命令长名称
//   - shortName: 命令短名称
//
// 返回值:
//   - *CmdSpec: 初始化的命令规格
//
// 功能说明:
//   - 创建基本命令规格
//   - 设置默认值
//   - 初始化所有字段
var NewCmdSpec = cmd.NewCmdSpec

// NewCmdOpts 创建新的命令选项
//
// 返回值:
//   - *CmdOpts: 初始化的命令选项
//
// 功能说明:
//   - 创建基本命令选项
//   - 初始化所有字段为零值
//   - 初始化 map 和 slice 避免空指针
var NewCmdOpts = cmd.NewCmdOpts

// NewCmdFromSpec 从规格创建命令
//
// 参数:
//   - spec: 命令规格结构体
//
// 返回值:
//   - *Cmd: 创建的命令实例
//   - error: 创建失败时返回错误
//
// 功能说明:
//   - 根据规格结构体创建命令
//   - 自动设置所有属性和配置
//   - 递归创建子命令
//   - 支持默认值处理
//   - 使用defer捕获panic, 转换为错误返回
var NewCmdFromSpec = cmd.NewCmdFromSpec

// GenerateCompletion 生成补全脚本
//
// 参数:
//   - cmd: 要生成补全脚本的命令
//   - shellType: Shell类型 (bash, pwsh, powershell)
//
// 返回值:
//   - string: 生成的补全脚本
//   - error: 生成失败时返回错误
//
// 功能说明:
//   - 为指定命令生成自动补全脚本
//   - 支持多种Shell类型
//   - 包含完整的命令树和标志信息
var GenerateCompletion = completion.Generate

// GenAndPrintCompletion 生成并打印补全脚本
//
// 参数:
//   - cmd: 要生成补全脚本的命令
//   - shellType: Shell类型 (bash, pwsh, powershell)
//
// 功能说明:
//   - 生成自动补全脚本
//   - 直接输出到标准输出
//   - 便于在命令行中直接使用
var GenAndPrintCompletion = completion.GenAndPrint

// StringFlag 字符串标志
// StringFlag 用于处理字符串类型的命令行参数。
// 它接受任何字符串值, 包括空字符串。
type StringFlag = flag.StringFlag

// BoolFlag 布尔标志
// BoolFlag 用于处理布尔类型的命令行参数。
// 它接受多种布尔值表示形式, 包括 "true", "false", "1", "0", "t", "f", "TRUE", "FALSE" 等。
type BoolFlag = flag.BoolFlag

// IntFlag 整数标志
// IntFlag 用于处理整数类型的命令行参数。
// 使用平台相关的int类型, 在32位系统上为32位整数, 在64位系统上为64位整数。
//
// 注意事项:
//   - 支持正数和负数
//   - 支持十进制格式
//   - 超出平台int范围会返回错误
type IntFlag = flag.IntFlag

// Int64Flag 64位整数标志
// Int64Flag 用于处理64位整数类型的命令行参数。
// 在所有平台上都使用固定的64位整数, 提供一致的行为。
//
// 注意事项:
//   - 支持正数和负数
//   - 支持十进制格式
//   - 范围: -9,223,372,036,854,775,808 到 9,223,372,036,854,775,807
type Int64Flag = flag.Int64Flag

// UintFlag 无符号整数标志
// UintFlag 用于处理无符号整数类型的命令行参数。
// 使用平台相关的uint类型, 在32位系统上为32位无符号整数, 在64位系统上为64位无符号整数。
//
// 注意事项:
//   - 只支持非负数
//   - 支持十进制格式
//   - 超出平台uint范围会返回错误
type UintFlag = flag.UintFlag

// Uint8Flag 8位无符号整数标志
// Uint8Flag 用于处理8位无符号整数类型的命令行参数。
// 适用于处理字节值、小范围计数器等场景。
//
// 注意事项:
//   - 只支持非负数
//   - 支持十进制格式
//   - 范围: 0 到 255
type Uint8Flag = flag.Uint8Flag

// Uint16Flag 16位无符号整数标志
// Uint16Flag 用于处理16位无符号整数类型的命令行参数。
// 适用于处理端口号、短范围计数器等场景。
//
// 注意事项:
//   - 只支持非负数
//   - 支持十进制格式
//   - 范围: 0 到 65,535
type Uint16Flag = flag.Uint16Flag

// Uint32Flag 32位无符号整数标志
// Uint32Flag 用于处理32位无符号整数类型的命令行参数。
// 适用于处理IP地址、大范围计数器等场景。
//
// 注意事项:
//   - 只支持非负数
//   - 支持十进制格式
//   - 范围: 0 到 4,294,967,295
type Uint32Flag = flag.Uint32Flag

// Uint64Flag 64位无符号整数标志
// Uint64Flag 用于处理64位无符号整数类型的命令行参数。
// 在所有平台上都使用固定的64位无符号整数, 提供一致的行为。
//
// 注意事项:
//   - 只支持非负数
//   - 支持十进制格式
//   - 范围: 0 到 18,446,744,073,709,551,615
type Uint64Flag = flag.Uint64Flag

// Float64Flag 64位浮点数标志
// Float64Flag 用于处理64位浮点数类型的命令行参数。
// 支持整数、小数和科学计数法表示的数值。
//
// 注意事项:
//   - 支持正数和负数
//   - 支持十进制格式和科学计数法
//   - 支持特殊值: NaN、+Inf、-Inf
//   - 精度遵循IEEE 754双精度浮点数标准
type Float64Flag = flag.Float64Flag

// EnumFlag 枚举标志
// EnumFlag 用于处理枚举类型的命令行参数, 限制输入值必须在预定义的允许值列表中。
// 使用映射表(map)实现O(1)时间复杂度的值查找, 提高性能。
//
// 特性:
//   - 使用映射表进行快速值验证
//   - 不允许空字符串作为枚举值
//   - 默认值必须在允许值列表中
//   - 不允许设置空值
type EnumFlag = flag.EnumFlag

// DurationFlag 持续时间标志
// DurationFlag 用于处理时间间隔类型的命令行参数。
// 支持Go标准库time.ParseDuration所支持的所有格式, 如 "300ms", "-1.5h", "2h45m" 等。
//
// 支持的格式:
//   - "ns": 纳秒
//   - "us" (或 "µs"): 微秒
//   - "ms": 毫秒
//   - "s": 秒
//   - "m": 分钟
//   - "h": 小时
//
// 注意事项:
//   - 支持负数表示负时间间隔
//   - 支持小数表示部分时间单位
//   - 可以组合多个单位, 如 "1h30m"
type DurationFlag = flag.DurationFlag

// TimeFlag 时间标志
// TimeFlag 用于处理时间类型的命令行参数。
// 支持自动检测多种常见时间格式, 也支持指定特定格式进行解析。
//
// 特性:
//   - 自动检测常见时间格式
//   - 支持自定义格式解析
//   - 记录当前使用的格式
//   - 线程安全的格式存储
//
// 常见支持格式:
//   - RFC3339: "2006-01-02T15:04:05Z07:00"
//   - RFC1123: "Mon, 02 Jan 2006 15:04:05 MST"
//   - 日期格式: "2006-01-02", "2006/01/02"
//   - 时间格式: "15:04:05", "15:04"
//   - 其他常见格式
type TimeFlag = flag.TimeFlag

// SizeFlag 大小标志 (支持KB、MB、GB等单位)
// SizeFlag 用于处理大小类型的命令行参数, 支持多种大小单位。
// 可以解析带有单位的大小值, 并将其转换为字节数。
//
// 支持的单位:
//   - B/b: 字节
//   - KB/kb/K/k: 千字节 (1024字节)
//   - MB/mb/M/m: 兆字节 (1024^2字节)
//   - GB/gb/G/g: 吉字节 (1024^3字节)
//   - TB/tb/T/t: 太字节 (1024^4字节)
//   - PB/pb/P/p: 拍字节 (1024^5字节)
//   - KiB/kib: 二进制千字节 (1024字节)
//   - MiB/mib: 二进制兆字节 (1024^2字节)
//   - GiB/gib: 二进制吉字节 (1024^3字节)
//   - TiB/tib: 二进制太字节 (1024^4字节)
//   - PiB/pib: 二进制拍字节 (1024^5字节)
//
// 注意事项:
//   - 支持小数, 如 "1.5MB"
//   - 不支持负数
//   - 默认单位为字节(B)
//   - 大小写不敏感
type SizeFlag = flag.SizeFlag

// MapFlag 用于处理键值对映射类型的命令行参数。
// 支持的格式: key1=value1,key2=value2
//
// 空值处理:
//   - 空字符串 "" 表示创建空映射
//   - ",,," 中的空对会被跳过
//   - 使用 SetKV 方法设置键值对时, 键不能为空
//   - 使用 Clear 方法可以清空映射
type MapFlag = flag.MapFlag

// StringSliceFlag 字符串切片标志
type StringSliceFlag = flag.StringSliceFlag

// IntSliceFlag 整数切片标志
type IntSliceFlag = flag.IntSliceFlag

// Int64SliceFlag 64位整数切片标志
type Int64SliceFlag = flag.Int64SliceFlag
