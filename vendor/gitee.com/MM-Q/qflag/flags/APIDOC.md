# Package flags

Package flags 定义了所有标志类型的通用接口和基础标志结构体。

## CONSTANTS

### 补全标志的使用说明

```go
const (
    CompletionShellDescCN = "生成指定的 Shell 补全脚本, 可选类型: %v"
    CompletionShellDescEN = "Generate the specified Shell completion script, optional types: %v"
)
```

### 支持的 Shell 类型

```go
const (
    ShellBash       = "bash"
    ShellPowershell = "powershell"
    ShellPwsh       = "pwsh"
    ShellNone       = "none"
)
```

### 标志的分隔符常量

```go
const (
    FlagSplitComma     = "," // 逗号
    FlagSplitSemicolon = ";" // 分号
    FlagSplitPipe      = "|" // 竖线
    FlagKVColon        = ":" // 冒号
    FlagKVEqual        = "=" // 等号
)
```

### 非法字符集常量

```go
const InvalidFlagChars = " !@#$%^&*(){}[]|\\;:'\"<>,.?"
```

## VARIABLES

### 内置标志名称

```go
var (
    HelpFlagName                 = "help"
    HelpFlagShortName            = "h"
    VersionFlagLongName          = "version"
    VersionFlagShortName         = "v"
    CompletionShellFlagLongName  = "generate-shell-completion"
    CompletionShellFlagShortName = "gsc"
)
```

### 内置标志使用说明

```go
var (
    HelpFlagUsageEn    = "Show help"
    VersionFlagUsageEn = "Show the version of the program"
    VersionFlagUsageZh = "显示程序的版本信息"
)
```

### 标志分隔符切片

```go
var FlagSplitSlice = []string{
    FlagSplitComma,
    FlagSplitSemicolon,
    FlagSplitPipe,
    FlagKVColon,
}
```

### 支持的 Shell 类型切片

```go
var ShellSlice = []string{ShellNone, ShellBash, ShellPowershell, ShellPwsh}
```

## FUNCTIONS

### FlagTypeToString

```go
func FlagTypeToString(flagType FlagType, withBrackets bool) string
```

将 `FlagType` 转换为字符串。

- 参数：
  - `flagType`: 需要转换的 `FlagType` 枚举值。
  - `withBrackets`: 是否在返回的字符串中包含尖括号。如果为 `true` 且 `flagType` 为 `bool` 时返回空字符串。

- 返回值：
  - 对应的类型字符串，如果类型未知则返回 `"unknown"` 或 `"<unknown>"`。

## TYPES

### BaseFlag

```go
type BaseFlag[T any] struct {
    // Has unexported fields.
}
```

`BaseFlag` 泛型基础标志结构体，封装所有标志的通用字段和方法。

#### BindEnv

```go
func (f *BaseFlag[T]) BindEnv(envName string) *BaseFlag[T]
```

绑定环境变量到标志。

- 参数：
  - `envName`: 环境变量名。

- 返回值：
  - 标志对象本身，支持链式调用。

#### Get

```go
func (f *BaseFlag[T]) Get() T
```

获取标志的实际值（泛型类型）。

- 返回值：
  - `T`: 标志值。

#### GetDefault

```go
func (f *BaseFlag[T]) GetDefault() T
```

获取标志的初始默认值（泛型类型）。

- 返回值：
  - `T`: 初始默认值。

#### GetDefaultAny

```go
func (f *BaseFlag[T]) GetDefaultAny() any
```

获取标志的初始默认值（`any` 类型）。

- 返回值：
  - `any`: 初始默认值。

#### GetEnvVar

```go
func (f *BaseFlag[T]) GetEnvVar() string
```

获取绑定的环境变量名。

- 返回值：
  - `string`: 环境变量名。

#### GetPointer

```go
func (f *BaseFlag[T]) GetPointer() *T
```

返回标志值的指针（泛型类型）。

- 返回值：
  - `*T`: 标志值的指针。

- 注意：
  1. 获取指针过程受锁保护，但直接修改指针指向的值仍会绕过验证机制。
  2. 多线程环境下修改时需额外同步措施，建议优先使用 `Set()` 方法。

#### Init

```go
func (f *BaseFlag[T]) Init(longName, shortName string, usage string, value *T) error
```

初始化标志的元数据和值指针，无需显式调用，仅在创建标志对象时自动调用。

- 参数：
  - `longName`: 长标志名称。
  - `shortName`: 短标志字符。
  - `usage`: 帮助说明。
  - `value`: 标志值指针。

- 返回值：
  - `error`: 初始化错误信息。

#### IsSet

```go
func (f *BaseFlag[T]) IsSet() bool
```

判断标志是否已被设置值。

- 返回值：
  - `bool`: `true` 表示已设置值，`false` 表示未设置。

#### LongName

```go
func (f *BaseFlag[T]) LongName() string
```

获取标志的长名称。

- 返回值：
  - `string`: 长标志名称。

#### Reset

```go
func (f *BaseFlag[T]) Reset()
```

将标志重置为初始默认值。

#### Set

```go
func (f *BaseFlag[T]) Set(value T) error
```

设置标志的值（泛型类型）。

- 参数：
  - `value T`: 标志值。

- 返回值：
  - `error`: 错误信息。

#### SetValidator

```go
func (f *BaseFlag[T]) SetValidator(validator Validator)
```

设置标志的验证器（泛型类型）。

- 参数：
  - `validator Validator`: 验证器接口。

#### ShortName

```go
func (f *BaseFlag[T]) ShortName() string
```

获取标志的短名称。

- 返回值：
  - `string`: 短标志字符。

#### String

```go
func (f *BaseFlag[T]) String() string
```

返回标志的字符串表示。

#### Type

```go
func (f *BaseFlag[T]) Type() FlagType
```

返回标志类型，默认实现返回 `0`，需要子类重写。

- 注意：
  - 具体标志类型需要重写此方法返回正确的 `FlagType`。

#### Usage

```go
func (f *BaseFlag[T]) Usage() string
```

获取标志的用法说明。

- 返回值：
  - `string`: 用法说明。

### BoolFlag

```go
type BoolFlag struct {
    BaseFlag[bool]
    // Has unexported fields.
}
```

`BoolFlag` 布尔类型标志结构体，继承 `BaseFlag[bool]` 泛型结构体，实现 `Flag` 接口。

#### IsBoolFlag

```go
func (f *BoolFlag) IsBoolFlag() bool
```

实现 `flag.boolFlag` 接口，返回 `true`。

#### Set

```go
func (f *BoolFlag) Set(value string) error
```

实现 `flag.Value` 接口，解析并设置布尔值。

- 参数：
  - `value`: 待设置的值。

- 返回值：
  - `error`: 解析或验证失败时返回错误信息。

#### String

```go
func (f *BoolFlag) String() string
```

实现 `flag.Value` 接口，返回布尔值字符串。

- 返回值：
  - `string`: 布尔值字符串。

#### Type

```go
func (f *BoolFlag) Type() FlagType
```

返回标志类型。

- 返回值：
  - `FlagType`: 标志类型枚举值。

### DurationFlag

```go
type DurationFlag struct {
    BaseFlag[time.Duration]
    // Has unexported fields.
}
```

`DurationFlag` 时间间隔类型标志结构体，继承 `BaseFlag[time.Duration]` 泛型结构体，实现 `Flag` 接口。

#### Set

```go
func (f *DurationFlag) Set(value string) error
```

实现 `flag.Value` 接口，解析并设置时间间隔值。

- 参数：
  - `value`: 待设置的值。

- 返回值：
  - `error`: 解析或验证失败时返回错误信息。

#### String

```go
func (f *DurationFlag) String() string
```

实现 `flag.Value` 接口，返回当前值的字符串表示。

- 返回值：
  - `string`: 当前值的字符串表示。

#### Type

```go
func (f *DurationFlag) Type() FlagType
```

返回标志类型。

- 返回值：
  - `FlagType`: 标志类型枚举值。

### EnumFlag

```go
type EnumFlag struct {
    BaseFlag[string]

    // Has unexported fields.
}
```

`EnumFlag` 枚举类型标志结构体，继承 `BaseFlag[string]` 泛型结构体，增加枚举特有的选项验证。

#### GetOptions

```go
func (f *EnumFlag) GetOptions() []string
```

返回枚举的所有可选值。

- 返回值：
  - `[]string`: 枚举的所有可选值。

#### Init

```go
func (f *EnumFlag) Init(longName, shortName string, defValue string, usage string, options []string) error
```

初始化枚举类型标志，无需显式调用，仅在创建标志对象时自动调用。

- 参数：
  - `longName`: 长标志名称。
  - `shortName`: 短标志字符。
  - `defValue`: 默认值。
  - `usage`: 帮助说明。
  - `options`: 枚举可选值列表。

- 返回值：
  - `error`: 初始化错误信息。

#### IsCheck

```go
func (f *EnumFlag) IsCheck(value string) error
```

检查枚举值是否有效。

- 参数：
  - `value`: 待检查的枚举值。

- 返回值：
  - `error`: 为 `nil` 说明值有效，否则返回错误信息。

#### Set

```go
func (f *EnumFlag) Set(value string) error
```

实现 `flag.Value` 接口，解析并设置枚举值。

- 参数：
  - `value`: 待设置的值。

- 返回值：
  - `error`: 解析或验证失败时返回错误信息。

#### SetCaseSensitive

```go
func (f *EnumFlag) SetCaseSensitive(sensitive bool) *EnumFlag
```

设置枚举值是否区分大小写。

- 参数：
  - `sensitive`: `true` 表示区分大小写，默认为 `false`。

- 返回值：
  - `*EnumFlag`: 返回自身以支持链式调用。

#### String

```go
func (f *EnumFlag) String() string
```

实现 `flag.Value` 接口，返回当前值的字符串表示。

- 返回值：
  - `string`: 当前值的字符串表示。

#### Type

```go
func (f *EnumFlag) Type() FlagType
```

返回标志类型。

- 返回值：
  - `FlagType`: 标志类型枚举值。

### Flag

```go
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
```

`Flag` 所有标志类型的通用接口，定义了标志的元数据访问方法。

### FlagMeta

```go
type FlagMeta struct {
    Flag Flag // 标志对象
}
```

`FlagMeta` 统一存储标志的完整元数据。

#### GetDefault

```go
func (m *FlagMeta) GetDefault() any
```

获取标志的默认值。

#### GetFlag

```go
func (m *FlagMeta) GetFlag() Flag
```

获取标志对象。

#### GetFlagType

```go
func (m *FlagMeta) GetFlagType() FlagType
```

获取标志的类型。

#### GetLongName

```go
func (m *FlagMeta) GetLongName() string
```

获取标志的长名称。

#### GetName

```go
func (m *FlagMeta) GetName() string
```

获取标志的名称。

- 优先返回长名称，如果长名称为空，则返回短名称。

#### GetShortName

```go
func (m *FlagMeta) GetShortName() string
```

获取标志的短名称。

#### GetUsage

```go
func (m *FlagMeta) GetUsage() string
```

获取标志的用法描述。

### FlagMetaInterface

```go
type FlagMetaInterface interface {
    GetFlagType() FlagType // 获取标志类型
    GetFlag() Flag         // 获取标志对象
    GetLongName() string   // 获取标志的长名称
    GetShortName() string  // 获取标志的短名称
    GetName() string       // 获取标志的名称
    GetUsage() string      // 获取标志的用法描述
    GetDefault() any       // 获取标志的默认值
    GetValue() any         // 获取标志的当前值
}
```

`FlagMetaInterface` 标志元数据接口，定义了标志元数据的获取方法。

### FlagRegistry

```go
type FlagRegistry struct {
    // Has unexported fields.
}
```

`FlagRegistry` 集中管理所有标志元数据及索引。

#### NewFlagRegistry

```go
func NewFlagRegistry() *FlagRegistry
```

创建一个空的标志注册表。

- 返回值：
  - `*FlagRegistry`: 创建的标志注册表指针。

#### GetALLFlagsCount

```go
func (r *FlagRegistry) GetALLFlagsCount() int
```

获取所有标志数量（长标志 + 短标志）。

- 返回值：
  - `int`: 所有标志的数量。

#### GetAllFlagMetas

```go
func (r *FlagRegistry) GetAllFlagMetas() []*FlagMeta
```

获取所有标志元数据列表。

- 返回值：
  - `[]*FlagMeta`: 所有标志元数据的切片。

#### GetAllFlags

```go
func (r *FlagRegistry) GetAllFlags() map[string]*FlagMeta
```

获取所有标志映射（长标志 + 短标志）。

- 返回值：
  - `map[string]*FlagMeta`: 长短标志名称到标志元数据的映射。

#### GetByLong

```go
func (r *FlagRegistry) GetByLong(longName string) (*FlagMeta, bool)
```

通过长标志名称查找对应的标志元数据。

- 参数：
  - `longName`: 标志的长名称。

- 返回值：
  - `*FlagMeta`: 找到的标志元数据指针，未找到时为 `nil`。
  - `bool`: 是否找到标志，`true` 表示找到。

#### GetByName

```go
func (r *FlagRegistry) GetByName(name string) (*FlagMeta, bool)
```

通过标志名称查找标志元数据。

- 参数：
  - `name`: 可以是长名称或短名称。

- 返回值：
  - `*FlagMeta`: 找到的标志元数据指针，未找到时为 `nil`。
  - `bool`: 是否找到标志，`true` 表示找到。

#### GetByShort

```go
func (r *FlagRegistry) GetByShort(shortName string) (*FlagMeta, bool)
```

通过短标志名称查找对应的标志元数据。

- 参数：
  - `shortName`: 标志的短名称。

- 返回值：
  - `*FlagMeta`: 找到的标志元数据指针，未找到时为 `nil`。
  - `bool`: 是否找到标志，`true` 表示找到。

#### GetLongFlags

```go
func (r *FlagRegistry) GetLongFlags() map[string]*FlagMeta
```

获取长标志映射。

- 返回值：
  - `map[string]*FlagMeta`: 长标志名称到标志元数据的映射。

#### GetLongFlagsCount

```go
func (r *FlagRegistry) GetLongFlagsCount() int
```

获取长标志数量。

- 返回值：
  - `int`: 长标志的数量。

#### GetShortFlags

```go
func (r *FlagRegistry) GetShortFlags() map[string]*FlagMeta
```

获取短标志映射。

- 返回值：
  - `map[string]*FlagMeta`: 短标志名称到标志元数据的映射。

#### GetShortFlagsCount

```go
func (r *FlagRegistry) GetShortFlagsCount() int
```

获取短标志数量。

- 返回值：
  - `int`: 短标志的数量。

#### RegisterFlag

```go
func (r *FlagRegistry) RegisterFlag(meta *FlagMeta) error
```

注册一个新的标志元数据到注册表中。

- 参数：
  - `meta`: 要注册的标志元数据。

- 返回值：
  - `error`: 注册错误信息。

- 注意：
  - 该方法线程安全，但发现重复标志时会 panic。

### FlagRegistryInterface

```go
type FlagRegistryInterface interface {
    GetAllFlagMetas() []*FlagMeta                  // 获取所有标志元数据列表
    GetLongFlags() map[string]*FlagMeta            // 获取长标志映射
    GetShortFlags() map[string]*FlagMeta           // 获取短标志映射
    GetAllFlags() map[string]*FlagMeta             // 获取所有标志映射（长+短）
    GetLongFlagsCount() int                        // 获取长标志数量
    GetShortFlagsCount() int                       // 获取短标志数量
    GetALLFlagsCount() int                         // 获取所有标志数量（长+短）
    RegisterFlag(meta *FlagMeta) error             // 注册一个新的标志元数据到注册表中
    GetByLong(longName string) (*FlagMeta, bool)   // 通过长标志名称查找对应的标志元数据
    GetByShort(shortName string) (*FlagMeta, bool) // 通过短标志名称查找对应的标志元数据
    GetByName(name string) (*FlagMeta, bool)       // 通过标志名称查找标志元数据
}
```

`FlagRegistryInterface` 标志注册表接口，定义了标志元数据的增删改查操作。

### FlagType

```go
type FlagType int
```

标志类型。

```go
const (
    FlagTypeUnknown FlagType = iota
    FlagTypeInt
    FlagTypeInt64
    FlagTypeUint16
    FlagTypeUint32
    FlagTypeUint64
    FlagTypeString
    FlagTypeBool
    FlagTypeFloat64
    FlagTypeEnum
    FlagTypeDuration
    FlagTypeSlice
    FlagTypeTime
    FlagTypeMap
    FlagTypePath
    FlagTypeIP4
    FlagTypeIP6
    FlagTypeURL
)
```

### Float64Flag

```go
type Float64Flag struct {
    BaseFlag[float64]
    // Has unexported fields.
}
```

`Float64Flag` 浮点型标志结构体，继承 `BaseFlag[float64]` 泛型结构体，实现 `Flag` 接口。

#### Set

```go
func (f *Float64Flag) Set(value string) error
```

实现 `flag.Value` 接口，解析并设置浮点值。

- 参数：
  - `value`: 待解析的浮点值。

- 返回值：
  - `error`: 解析错误或验证错误。

#### Type

```go
func (f *Float64Flag) Type() FlagType
```

返回标志类型。

- 返回值：
  - `FlagType`: 标志类型枚举值。

### IP4Flag

```go
type IP4Flag struct {
    BaseFlag[string]
    // Has unexported fields.
}
```

`IP4Flag` IPv4 地址类型标志结构体，继承 `BaseFlag[string]` 泛型结构体，实现 `Flag` 接口。

#### Set

```go
func (f *IP4Flag) Set(value string) error
```

实现 `flag.Value` 接口，解析并验证 IPv4 地址。

- 参数：
  - `value`: 待解析的 IPv4 地址值。

- 返回值：
  - `error`: 解析或验证错误。

#### String

```go
func (f *IP4Flag) String() string
```

实现 `flag.Value` 接口，返回当前值的字符串表示。

- 返回值：
  - `string`: 当前值的字符串表示。

#### Type

```go
func (f *IP4Flag) Type() FlagType
```

返回标志类型。

- 返回值：
  - `FlagType`: 标志类型枚举值。

### IP6Flag

```go
type IP6Flag struct {
    BaseFlag[string]
    // Has unexported fields.
}
```

`IP6Flag` IPv6 地址类型标志结构体，继承 `BaseFlag[string]` 泛型结构体，实现 `Flag` 接口。

#### Set

```go
func (f *IP6Flag) Set(value string) error
```

实现 `flag.Value` 接口，解析并验证 IPv6 地址。

- 参数：
  - `value`: 待解析的 IPv6 地址值。

- 返回值：
  - `error`: 解析或验证错误。

#### String

```go
func (f *IP6Flag) String() string
```

实现 `flag.Value` 接口，返回当前值的字符串表示。

- 返回值：
  - `string`: 当前值的字符串表示。

#### Type

```go
func (f *IP6Flag) Type() FlagType
```

返回标志类型。

- 返回值：
  - `FlagType`: 标志类型枚举值。

### Int64Flag

```go
type Int64Flag struct {
    BaseFlag[int64]
    // Has unexported fields.
}
```

`Int64Flag` 64 位整数类型标志结构体，继承 `BaseFlag[int64]` 泛型结构体，实现 `Flag` 接口。

#### Set

```go
func (f *Int64Flag) Set(value string) error
```

实现 `flag.Value` 接口，解析并设置 64 位整数值。

- 参数：
  - `value`: 待解析的 64 位整数值。

- 返回值：
  - `error`: 解析错误或验证错误。

#### SetRange

```go
func (f *Int64Flag) SetRange(min, max int64)
```

设置 64 位整数的有效范围。

- 参数：
  - `min`: 最小值。
  - `max`: 最大值。

#### Type

```go
func (f *Int64Flag) Type() FlagType
```

返回标志类型。

- 返回值：
  - `FlagType`: 标志类型枚举值。

### IntFlag

```go
type IntFlag struct {
    BaseFlag[int]
    // Has unexported fields.
}
```

`IntFlag` 整数类型标志结构体，继承 `BaseFlag[int]` 泛型结构体，实现 `Flag` 接口。

#### Set

```go
func (f *IntFlag) Set(value string) error
```

实现 `flag.Value` 接口，解析并验证整数值。

- 参数：
  - `value`: 待解析的整数值。

- 返回值：
  - `error`: 解析错误或验证错误。

#### SetRange

```go
func (f *IntFlag) SetRange(min, max int)
```

设置整数的有效范围。

- 参数：
  - `min`: 最小值。
  - `max`: 最大值。

#### String

```go
func (f *IntFlag) String() string
```

实现 `flag.Value` 接口，返回当前整数值的字符串表示。

- 返回值：
  - `string`: 当前整数值的字符串表示。

#### Type

```go
func (f *IntFlag) Type() FlagType
```

返回标志类型。

- 返回值：
  - `FlagType`: 标志类型枚举值。

### MapFlag

```go
type MapFlag struct {
    BaseFlag[map[string]string]

    // Has unexported fields.
}
```

`MapFlag` 键值对类型标志结构体，继承 `BaseFlag[map[string]string]` 泛型结构体，实现 `Flag` 接口。

#### Set

```go
func (f *MapFlag) Set(value string) error
```

实现 `flag.Value` 接口，解析并设置键值对。

- 参数：
  - `value`: 待设置的值。

- 返回值：
  - `error`: 解析或验证失败时返回错误信息。

#### SetDelimiters

```go
func (f *MapFlag) SetDelimiters(keyDelimiter, valueDelimiter string)
```

设置键值对分隔符。

- 参数：
  - `keyDelimiter`: 键值对分隔符。
  - `valueDelimiter`: 键值分隔符。

#### SetIgnoreCase

```go
func (f *MapFlag) SetIgnoreCase(enable bool)
```

设置是否忽略键的大小写。

- 参数：
  - `enable`: 是否忽略键的大小写。

- 注意：
  - 当 `enable` 为 `true` 时，所有键将转换为小写进行存储和比较。

#### String

```go
func (f *MapFlag) String() string
```

实现 `flag.Value` 接口，返回当前值的字符串表示。

- 返回值：
  - `string`: 当前值的字符串表示。

#### Type

```go
func (f *MapFlag) Type() FlagType
```

返回标志类型。

- 返回值：
  - `FlagType`: 标志类型枚举值。

### PathFlag

```go
type PathFlag struct {
    BaseFlag[string]
    // Has unexported fields.
}
```

`PathFlag` 路径类型标志结构体，继承 `BaseFlag[string]` 泛型结构体，实现 `Flag` 接口。

#### Init

```go
func (f *PathFlag) Init(longName, shortName string, defValue string, usage string) error
```

初始化路径标志。

- 参数：
  - `longName`: 长名称。
  - `shortName`: 短名称。
  - `defValue`: 默认值。
  - `usage`: 使用说明。

- 返回值：
  - `error`: 初始化错误。

#### IsDirectory

```go
func (f *PathFlag) IsDirectory(isDir bool) *PathFlag
```

设置路径是否必须是目录。

- 参数：
  - `isDir`: 是否必须是目录。

- 返回值：
  - `*PathFlag`: 当前路径标志对象。

- 示例：
  ```go
  cmd.Path("log-dir", "l", "/var/log/app", "日志目录").IsDirectory(true)
  ```

#### MustExist

```go
func (f *PathFlag) MustExist(mustExist bool) *PathFlag
```

设置路径是否必须存在。

- 参数：
  - `mustExist`: 是否必须存在。

- 返回值：
  - `*PathFlag`: 当前路径标志对象。

- 示例：
  ```go
  cmd.Path("output", "o", "/tmp/output", "输出目录").MustExist(false)
  ```

#### Set

```go
func (f *PathFlag) Set(value string) error
```

实现 `flag.Value` 接口，解析并验证路径。

- 参数：
  - `value`: 待解析的路径值。

- 返回值：
  - `error`: 解析错误或验证错误。

#### String

```go
func (f *PathFlag) String() string
```

实现 `flag.Value` 接口，返回当前值的字符串表示。

#### Type

```go
func (f *PathFlag) Type() FlagType
```

返回标志类型。

### SliceFlag

```go
type SliceFlag struct {
    BaseFlag[[]string] // 基类

    // Has unexported fields.
}
```

`SliceFlag` 切片类型标志结构体，继承 `BaseFlag[[]string]` 泛型结构体，实现 `Flag` 接口。

#### Clear

```go
func (f *SliceFlag) Clear() error
```

清空切片所有元素。

- 返回值：
  - 操作成功返回 `nil`，否则返回错误信息。

- 注意：
  - 该方法会改变切片的指针。

#### Contains

```go
func (f *SliceFlag) Contains(element string) bool
```

检查切片是否包含指定元素。

- 参数：
  - `element`: 待检查的元素。

- 返回值：
  - 若切片包含指定元素，返回 `true`，否则返回 `false`。

- 注意：
  - 当切片未设置值时，将使用默认值进行检查。

#### GetDelimiters

```go
func (f *SliceFlag) GetDelimiters() []string
```

获取当前分隔符列表。

#### Init

```go
func (f *SliceFlag) Init(longName, shortName string, defValue []string, usage string) error
```

初始化切片类型标志。

- 参数：
  - `longName`: 长标志名称。
  - `shortName`: 短标志字符。
  - `defValue`: 默认值（切片类型）。
  - `usage`: 帮助说明。

- 返回值：
  - `error`: 初始化错误信息。

#### Len

```go
func (f *SliceFlag) Len() int
```

获取切片长度。

- 返回值：
  - 获取切片长度。

#### Remove

```go
func (f *SliceFlag) Remove(element string) error
```

从切片中移除指定元素（支持移除空字符串元素）。

- 参数：
  - `element`: 待移除的元素（支持空字符串）。

- 返回值：
  - 操作成功返回 `nil`，否则返回错误信息。

#### Set

```go
func (f *SliceFlag) Set(value string) error
```

实现 `flag.Value` 接口，解析并设置切片值。

- 参数：
  - `value`: 待解析的切片值。

- 注意：
  - 如果切片中包含分隔符，则根据分隔符进行分割，否则将整个值作为单个元素。
  - 例如：`"a,b,c" -> ["a", "b", "c"]`。

#### SetDelimiters

```go
func (f *SliceFlag) SetDelimiters(delimiters []string)
```

设置切片解析的分隔符列表。

- 参数：
  - `delimiters`: 分隔符列表。

#### SetSkipEmpty

```go
func (f *SliceFlag) SetSkipEmpty(skip bool)
```

设置是否跳过空元素。

- 参数：
  - `skip`: 为 `true` 时跳过空元素，为 `false` 时保留空元素。

- 线程安全的空元素跳过更新。

#### Sort

```go
func (f *SliceFlag) Sort() error
```

对切片进行排序，对当前切片标志的值进行原地排序，修改原切片内容。采用 Go 标准库的 `sort.Strings()` 函数进行字典序排序（按 Unicode 代码点升序排列）。

- 注意：
  - 排序会直接修改当前标志的值，而非返回新切片。
  - 排序区分大小写，遵循 Unicode 代码点比较规则（如 `'A' < 'a' < 'z'`）。
  - 若切片未设置值，将使用默认值进行排序。

- 返回值：
  - 排序成功返回 `nil`，若排序过程中发生错误则返回错误信息。

#### String

```go
func (f *SliceFlag) String() string
```

实现 `flag.Value` 接口，返回当前值的字符串表示。

#### Type

```go
func (f *SliceFlag) Type() FlagType
```

返回标志类型。

### StringFlag

```go
type StringFlag struct {
    BaseFlag[string]
}
```

`StringFlag` 字符串类型标志结构体，继承 `BaseFlag[string]` 泛型结构体，实现 `Flag` 接口。

#### Contains

```go
func (f *StringFlag) Contains(substr string) bool
```

检查字符串是否包含指定子串。

- 参数：
  - `substr`: 子串。

- 返回值：
  - `bool`: 如果包含子串则返回 `true`，否则返回 `false`。

#### Len

```go
func (f *StringFlag) Len() int
```

获取字符串标志的长度。

- 返回值：
  - 字符串的字符数（按 UTF-8 编码计算）。

#### Set

```go
func (f *StringFlag) Set(value string) error
```

实现 `flag.Value` 接口的 `Set` 方法，将字符串值解析并设置到标志中。

- 参数：
  - `value`: 待设置的字符串值。

- 返回值：
  - `error`: 设置失败时返回错误信息。

#### String

```go
func (f *StringFlag) String() string
```

返回带引号的字符串值。

- 返回值：
  - `string`: 带引号的字符串值。

#### ToLower

```go
func (f *StringFlag) ToLower() string
```

将字符串标志值转换为小写。

#### ToUpper

```go
func (f *StringFlag) ToUpper() string
```

将字符串标志值转换为大写。

#### Type

```go
func (f *StringFlag) Type() FlagType
```

返回标志类型。

- 返回值：
  - `FlagType`: 标志类型枚举值。

### TimeFlag

```go
type TimeFlag struct {
    BaseFlag[time.Time]

    // Has unexported fields.
}
```

`TimeFlag` 时间类型标志结构体，继承 `BaseFlag[time.Time]` 泛型结构体，实现 `Flag` 接口。

#### Set

```go
func (f *TimeFlag) Set(value string) error
```

实现 `flag.Value` 接口，解析并设置时间值。

- 参数：
  - `value`: 待解析的时间字符串。

- 返回值：
  - `error`: 解析或验证失败时返回错误信息。

#### SetOutputFormat

```go
func (f *TimeFlag) SetOutputFormat(format string)
```

设置时间输出格式。

- 参数：
  - `format`: 时间格式化字符串。

- 注意：
  - 此方法线程安全。

#### String

```go
func (f *TimeFlag) String() string
```

实现 `flag.Value` 接口，返回当前时间的字符串表示。

- 返回值：
  - `string`: 格式化后的时间字符串。

- 注意：
  - 加锁保证 `outputFormat` 和 `value` 的并发安全访问。

#### Type

```go
func (f *TimeFlag) Type() FlagType
```

返回标志类型。

- 返回值：
  - `FlagType`: 标志类型枚举值。

### TypedFlag

```go
type TypedFlag[T any] interface {
    Flag            // 继承标志接口
    GetDefault() T  // 获取标志的具体类型默认值
    Get() T         // 获取标志的具体类型值
    GetPointer() *T // 获取标志值的指针
    Set(T) error    // 设置标志的具体类型值
    SetValidator(Validator) // 设置标志的验证器
    BindEnv(envName string) *BaseFlag[T] // 绑定环境变量
}
```

`TypedFlag` 所有标志类型的通用接口，定义了标志的元数据访问方法和默认值访问方法。

### URLFlag

```go
type URLFlag struct {
    BaseFlag[string]
    // Has unexported fields.
}
```

`URLFlag` URL 类型标志结构体，继承 `BaseFlag[string]` 泛型结构体，实现 `Flag` 接口。

#### Set

```go
func (f *URLFlag) Set(value string) error
```

实现 `flag.Value` 接口，解析并验证 URL 格式。

- 参数：
  - `value`: 待解析的 URL 字符串。

- 返回值：
  - `error`: 解析或验证失败时返回错误信息。

#### String

```go
func (f *URLFlag) String() string
```

实现 `flag.Value` 接口，返回当前值的字符串表示。

- 返回值：
  - `string`: 当前 URL 值的字符串表示。

#### Type

```go
func (f *URLFlag) Type() FlagType
```

返回标志类型。

- 返回值：
  - `FlagType`: 标志类型枚举值。

### Uint16Flag

```go
type Uint16Flag struct {
    BaseFlag[uint16] // 基类
    // Has unexported fields.
}
```

`Uint16Flag` 16 位无符号整数类型标志结构体，继承 `BaseFlag[uint16]` 泛型结构体，实现 `Flag` 接口。

#### Set

```go
func (f *Uint16Flag) Set(value string) error
```

实现 `flag.Value` 接口，解析并设置 16 位无符号整数值。验证值是否在 uint16 范围内（`0-65535`）。

- 参数：
  - `value`: 待设置的值（`0-65535`）。

- 返回值：
  - `error`: 解析或验证失败时返回错误信息。

#### String

```go
func (f *Uint16Flag) String() string
```

实现 `flag.Value` 接口，返回当前值的字符串表示。

- 返回值：
  - `string`: 当前值的字符串表示。

#### Type

```go
func (f *Uint16Flag) Type() FlagType
```

返回标志类型。

- 返回值：
  - `FlagType`: 标志类型枚举值。

### Uint32Flag

```go
type Uint32Flag struct {
    BaseFlag[uint32] // 基类
    // Has unexported fields.
}
```

`Uint32Flag` 32 位无符号整数类型标志结构体，继承 `BaseFlag[uint32]` 泛型结构体，实现 `Flag` 接口。

#### Set

```go
func (f *Uint32Flag) Set(value string) error
```

实现 `flag.Value` 接口，解析并设置 32 位无符号整数值。验证值是否在 uint32 范围内（`0-4294967295`）。

- 参数：
  - `value`: 待设置的值（`0-4294967295`）。

- 返回值：
  - `error`: 解析或验证失败时返回错误信息。

#### String

```go
func (f *Uint32Flag) String() string
```

实现 `flag.Value` 接口，返回当前值的字符串表示。

#### Type

```go
func (f *Uint32Flag) Type() FlagType
```

返回标志类型。

### Uint64Flag

```go
type Uint64Flag struct {
    BaseFlag[uint64] // 基类
    // Has unexported fields.
}
```

`Uint64Flag` 64 位无符号整数类型标志结构体，继承 `BaseFlag[uint64]` 泛型结构体，实现 `Flag` 接口。

#### Set

```go
func (f *Uint64Flag) Set(value string) error
```

实现 `flag.Value` 接口，解析并设置 64 位无符号整数值。验证值是否在 uint64 范围内（`0-18446744073709551615`）。

- 参数：
  - `value`: 待设置的值（`0-18446744073709551615`）。

- 返回值：
  - `error`: 解析或验证失败时返回错误信息。

#### String

```go
func (f *Uint64Flag) String() string
```

实现 `flag.Value` 接口，返回当前值的字符串表示。

#### Type

```go
func (f *Uint64Flag) Type() FlagType
```

返回标志类型。

### Validator

```go
type Validator interface {
    Validate(value any) error
}
```

验证器接口，所有自定义验证器需实现此接口。

- 方法：
  - `Validate`: 验证参数值是否合法。
    - 参数：
      - `value`: 待验证的参数值。
    - 返回值：
      - `error`: 验证通过返回 `nil`，否则返回错误信息。