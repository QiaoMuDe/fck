# Package flags

Package flags 提供了命令行标志的定义和解析功能，支持多种数据类型和复杂的标志结构。

## 常量定义

### 补全标志说明

```go
const (
    CompletionShellDescCN = "生成指定的 Shell 补全脚本, 可选类型: %v"
    CompletionShellDescEN = "Generate the specified Shell completion script, optional types: %v"
)
```

定义中英文的补全标志使用说明。

### 支持的 Shell 类型

```go
const (
    ShellBash       = "bash"       // bash shell
    ShellPowershell = "powershell" // powershell shell
    ShellPwsh       = "pwsh"       // pwsh shell
    ShellNone       = "none"       // 无 shell
)
```

定义支持的 Shell 类型。

### 标志分隔符

```go
const (
    FlagSplitComma     = ","         // 逗号
    FlagSplitSemicolon = ";"         // 分号
    FlagSplitPipe      = "|"         // 竖线
    FlagKVColon        = ":"         // 冒号
    FlagKVEqual        = "="         // 等号
)
```

定义标志的分隔符常量。

### 非法字符集

```go
const InvalidFlagChars = " !@#$%^&*(){}[]|\\;:'\"<>,.?/"
```

定义非法字符集常量，防止非法的标志名称。

## 变量定义

### 内置标志名称

```go
var (
    HelpFlagName                 = "help"                      // 帮助标志名称
    HelpFlagShortName            = "h"                         // 帮助标志短名称
    VersionFlagLongName          = "version"                   // 版本标志名称
    VersionFlagShortName         = "v"                         // 版本标志短名称
    CompletionShellFlagLongName  = "completion"                // 生成 shell 补全标志长名称
)
```

定义内置标志的名称。

### 内置标志使用说明

```go
var (
    HelpFlagUsage    = "Show help"    // 帮助标志使用说明
    VersionFlagUsage = "Show version" // 版本标志使用说明
)
```

定义内置标志的使用说明。

### 标志分隔符切片

```go
var FlagSplitSlice = []string{
    FlagSplitComma,
    FlagSplitSemicolon,
    FlagSplitPipe,
    FlagKVColon,
}
```

定义支持的标志分隔符切片。

### Shell 类型切片

```go
var ShellSlice = []string{ShellNone, ShellBash, ShellPowershell, ShellPwsh}
```

定义支持的 Shell 类型切片。

## 函数定义

### FlagType 转换

```go
func FlagTypeToString(flagType FlagType) string
```

将 FlagType 转换为带语义信息的字符串。

- 参数: `flagType` - 需要转换的 FlagType 枚举值
- 返回值: 带语义信息的类型字符串，用于命令行帮助信息显示

## 类型定义

### BaseFlag 泛型基础标志结构体

```go
type BaseFlag[T any] struct {
    // Has unexported fields.
}
```

BaseFlag 是泛型基础标志结构体，封装所有标志的通用字段和方法。

#### 方法

```go
func (f *BaseFlag[T]) BindEnv(envName string) *BaseFlag[T]
```

绑定环境变量到标志。

- 参数: `envName` - 环境变量名
- 返回值: 标志对象本身，支持链式调用

```go
func (f *BaseFlag[T]) Get() T
```

获取标志的实际值（泛型类型）。

- 返回值: T - 标志值

```go
func (f *BaseFlag[T]) GetDefault() T
```

获取标志的初始默认值（泛型类型）。

- 返回值: T - 初始默认值

```go
func (f *BaseFlag[T]) GetDefaultAny() any
```

获取标志的初始默认值（any 类型）。

- 返回值: any - 初始默认值

```go
func (f *BaseFlag[T]) GetEnvVar() string
```

获取绑定的环境变量名。

- 返回值: string - 环境变量名

```go
func (f *BaseFlag[T]) Init(longName, shortName string, usage string, value *T) error
```

初始化标志的元数据和值指针。

- 参数:
  - `longName` - 长标志名称
  - `shortName` - 短标志字符
  - `usage` - 帮助说明
  - `value` - 标志值指针
- 返回值: error - 初始化错误信息

```go
func (f *BaseFlag[T]) IsSet() bool
```

判断标志是否已被设置值。

- 返回值: bool - true 表示已设置值，false 表示未设置

```go
func (f *BaseFlag[T]) LongName() string
```

获取标志的长名称。

- 返回值: string - 长标志名称

```go
func (f *BaseFlag[T]) Name() string
```

获取标志的名称。

- 返回值: string - 标志名称，优先返回长名称，如果长名称为空则返回短名称

```go
func (f *BaseFlag[T]) Reset()
```

将标志重置为初始默认值。

```go
func (f *BaseFlag[T]) Set(value T) error
```

设置标志的值（泛型类型）。

- 参数: `value` - 标志值
- 返回值: error - 错误信息

```go
func (f *BaseFlag[T]) SetValidator(validator Validator)
```

设置标志的验证器。

- 参数: `validator` - 验证器接口

```go
func (f *BaseFlag[T]) ShortName() string
```

获取标志的短名称。

- 返回值: string - 短标志字符

```go
func (f *BaseFlag[T]) String() string
```

返回标志的字符串表示。

```go
func (f *BaseFlag[T]) Type() FlagType
```

返回标志类型。

- 注意: 具体标志类型需要重写此方法返回正确的 FlagType。

```go
func (f *BaseFlag[T]) Usage() string
```

获取标志的用法说明。

- 返回值: string - 用法说明

### BoolFlag 布尔类型标志结构体

```go
type BoolFlag struct {
    BaseFlag[bool]
    // Has unexported fields.
}
```

BoolFlag 是布尔类型标志结构体，继承 BaseFlag[bool] 泛型结构体，实现 Flag 接口。

#### 方法

```go
func (f *BoolFlag) IsBoolFlag() bool
```

实现 flag.boolFlag 接口，返回 true。

```go
func (f *BoolFlag) Set(value string) error
```

实现 flag.Value 接口，解析并设置布尔值。

- 支持以下布尔值格式（大小写不敏感）:
  - 真值: "true", "1", "t", "T", "TRUE", "True"
  - 假值: "false", "0", "f", "F", "FALSE", "False"
- 参数: `value` - 待设置的布尔值字符串
- 返回值: error - 解析或验证失败时返回错误信息
- 示例:
  - `flag.Set("true")` // ✅ 成功，值为 true
  - `flag.Set("1")` // ✅ 成功，值为 true
  - `flag.Set("FALSE")` // ✅ 成功，值为 false
  - `flag.Set("0")` // ✅ 成功，值为 false
  - `flag.Set("yes")` // ❌ 失败，返回解析错误

```go
func (f *BoolFlag) String() string
```

实现 flag.Value 接口，返回布尔值字符串。

- 返回值: string - 布尔值字符串

```go
func (f *BoolFlag) Type() FlagType
```

返回标志类型。

- 返回值: FlagType - 标志类型枚举值

### DurationFlag 时间间隔类型标志结构体

```go
type DurationFlag struct {
    BaseFlag[time.Duration]
    // Has unexported fields.
}
```

DurationFlag 是时间间隔类型标志结构体，继承 BaseFlag[time.Duration] 泛型结构体，实现 Flag 接口。

#### 方法

```go
func (f *DurationFlag) Set(value string) error
```

实现 flag.Value 接口，解析并设置时间间隔值。

- 参数: `value` - 待设置的值
- 返回值: error - 解析或验证失败时返回错误信息

```go
func (f *DurationFlag) String() string
```

实现 flag.Value 接口，返回当前值的字符串表示。

- 返回值: string - 当前值的字符串表示

```go
func (f *DurationFlag) Type() FlagType
```

返回标志类型。

- 返回值: FlagType - 标志类型枚举值

### EnumFlag 枚举类型标志结构体

```go
type EnumFlag struct {
    BaseFlag[string]

    // Has unexported fields.
}
```

EnumFlag 是枚举类型标志结构体，继承 BaseFlag[string] 泛型结构体，增加枚举特有的选项验证。

#### 方法

```go
func (f *EnumFlag) GetOptions() []string
```

返回枚举的所有可选值。

- 返回值: []string - 枚举的所有可选值

```go
func (f *EnumFlag) Init(longName, shortName string, defValue string, usage string, options []string) error
```

初始化枚举类型标志。

- 参数:
  - `longName` - 长标志名称
  - `shortName` - 短标志字符
  - `defValue` - 默认值
  - `usage` - 帮助说明
  - `options` - 枚举可选值列表
- 返回值: error - 初始化错误信息

```go
func (f *EnumFlag) IsCheck(value string) error
```

检查枚举值是否有效。

- 参数: `value` - 待检查的枚举值
- 返回值: error - 为 nil 说明值有效，否则返回错误信息

```go
func (f *EnumFlag) Set(value string) error
```

实现 flag.Value 接口，解析并设置枚举值。

- 参数: `value` - 待设置的值
- 返回值: error - 解析或验证失败时返回错误信息

```go
func (f *EnumFlag) SetCaseSensitive(sensitive bool) *EnumFlag
```

设置枚举值是否区分大小写。

- 参数: `sensitive` - true 表示区分大小写，false 表示不区分（默认）
- 返回值: *EnumFlag - 返回自身以支持链式调用

```go
func (f *EnumFlag) String() string
```

实现 flag.Value 接口，返回当前值的字符串表示。

- 返回值: string - 当前值的字符串表示

```go
func (f *EnumFlag) Type() FlagType
```

返回标志类型。

- 返回值: FlagType - 标志类型枚举值

### Flag 接口

```go
type Flag interface {
    LongName() string   // 获取标志的长名称
    ShortName() string  // 获取标志的短名称
    Usage() string      // 获取标志的用法
    Type() FlagType     // 获取标志类型
    GetDefaultAny() any // 获取标志的默认值 (any 类型)
    String() string     // 获取标志的字符串表示
    IsSet() bool        // 判断标志是否已设置值
    Reset()             // 重置标志值为默认值
    GetEnvVar() string  // 获取标志绑定的环境变量名称
}
```

Flag 是所有标志类型的通用接口，定义了标志的元数据访问方法。

### FlagMeta 标志元数据

```go
type FlagMeta struct {
    Flag Flag // 标志对象
}
```

FlagMeta 统一存储标志的完整元数据。

#### 方法

```go
func (m *FlagMeta) GetDefault() any
```

获取标志的默认值。

```go
func (m *FlagMeta) GetFlag() Flag
```

获取标志对象。

```go
func (m *FlagMeta) GetFlagType() FlagType
```

获取标志的类型。

```go
func (m *FlagMeta) GetLongName() string
```

获取标志的长名称。

```go
func (m *FlagMeta) GetName() string
```

获取标志的名称。

- 优先返回长名称，如果长名称为空，则返回短名称。

```go
func (m *FlagMeta) GetShortName() string
```

获取标志的短名称。

```go
func (m *FlagMeta) GetUsage() string
```

获取标志的用法描述。

### FlagRegistry 标志注册表

```go
type FlagRegistry struct {
    // Has unexported fields.
}
```

FlagRegistry 集中管理所有标志元数据及索引。

#### 方法

```go
func NewFlagRegistry() *FlagRegistry
```

创建一个空的标志注册表。

- 返回值: *FlagRegistry - 创建的标志注册表指针

```go
func (r *FlagRegistry) GetAllFlagsCount() int
```

获取所有标志数量（长标志 + 短标志）。

- 返回值: int - 所有标志的数量

```go
func (r *FlagRegistry) GetByLong(longName string) (*FlagMeta, bool)
```

通过长标志名称查找对应的标志元数据。

- 参数: `longName` - 标志的长名称（如 "help"）
- 返回值:
  - *FlagMeta - 找到的标志元数据指针，未找到时为 nil
  - bool - 是否找到标志，true 表示找到

```go
func (r *FlagRegistry) GetByName(name string) (*FlagMeta, bool)
```

通过标志名称查找标志元数据。

- 参数: `name` - 可以是长名称（如 "help"）或短名称（如 "h"）
- 返回值:
  - *FlagMeta - 找到的标志元数据指针，未找到时为 nil
  - bool - 是否找到标志，true 表示找到

```go
func (r *FlagRegistry) GetByShort(shortName string) (*FlagMeta, bool)
```

通过短标志名称查找对应的标志元数据。

- 参数: `shortName` - 标志的短名称（如 "h" 对应 "help"）
- 返回值:
  - *FlagMeta - 找到的标志元数据指针，未找到时为 nil
  - bool - 是否找到标志，true 表示找到

```go
func (r *FlagRegistry) GetFlagMetaCount() int
```

获取标志元数据数量。

- 返回值: int - 标志元数据的数量

```go
func (r *FlagRegistry) GetFlagMetaList() []*FlagMeta
```

获取所有标志元数据列表。

- 返回值: []*FlagMeta - 所有标志元数据的切片

```go
func (r *FlagRegistry) GetFlagNameMap() map[string]*FlagMeta
```

获取所有标志映射（长标志 + 短标志）。

- 返回值: map[string]*FlagMeta - 长短标志名称到标志元数据的映射

```go
func (r *FlagRegistry) GetLongFlagMap() map[string]*FlagMeta
```

获取长标志映射。

- 返回值: map[string]*FlagMeta - 长标志名称到标志元数据的映射

```go
func (r *FlagRegistry) GetLongFlagsCount() int
```

获取长标志数量。

- 返回值: int - 长标志的数量

```go
func (r *FlagRegistry) GetShortFlagMap() map[string]*FlagMeta
```

获取短标志映射。

- 返回值: map[string]*FlagMeta - 短标志名称到标志元数据的映射

```go
func (r *FlagRegistry) GetShortFlagsCount() int
```

获取短标志数量。

- 返回值: int - 短标志的数量

```go
func (r *FlagRegistry) RegisterFlag(meta *FlagMeta) error
```

注册一个新的标志元数据到注册表中。

- 参数: `meta` - 要注册的标志元数据
- 该方法会执行以下操作:
  1. 检查长名称和短名称是否已存在
  2. 将标志添加到长名称索引
  3. 将标志添加到短名称索引
  4. 将标志添加到所有标志列表
- 返回值: error - 错误信息，无错误时为 nil

### FlagType 标志类型

```go
type FlagType int
```

标志类型。

#### 常量

```go
const (
    FlagTypeUnknown  FlagType = iota // 未知类型
    FlagTypeInt                      // 整数类型
    FlagTypeInt64                    // 64 位整数类型
    FlagTypeUint16                   // 16 位无符号整数类型
    FlagTypeUint32                   // 32 位无符号整数类型
    FlagTypeUint64                   // 64 位无符号整数类型
    FlagTypeString                   // 字符串类型
    FlagTypeBool                     // 布尔类型
    FlagTypeEnum                     // 枚举类型
    FlagTypeDuration                 // 时间间隔类型
    FlagTypeStringSlice              // 字符串切片类型
    FlagTypeIntSlice                 // 整数切片类型
    FlagTypeInt64Slice               // 64位整数切片类型
    FlagTypeSize                     // 大小类型
    FlagTypeTime                     // 时间类型
    FlagTypeMap                      // 映射类型
)
```

定义各种标志类型。

### Float64Flag 浮点型标志结构体

```go
type Float64Flag struct {
    BaseFlag[float64]
    // Has unexported fields.
}
```

Float64Flag 是浮点型标志结构体，继承 BaseFlag[float64] 泛型结构体，实现 Flag 接口。

#### 方法

```go
func (f *Float64Flag) Set(value string) error
```

实现 flag.Value 接口，解析并设置浮点值。

- 参数: `value` - 待解析的浮点值
- 返回值: error - 解析错误或验证错误

```go
func (f *Float64Flag) Type() FlagType
```

返回标志类型。

- 返回值: FlagType - 标志类型枚举值

### Int64Flag 64 位整数类型标志结构体

```go
type Int64Flag struct {
    BaseFlag[int64]
    // Has unexported fields.
}
```

Int64Flag 是 64 位整数类型标志结构体，继承 BaseFlag[int64] 泛型结构体，实现 Flag 接口。

#### 方法

```go
func (f *Int64Flag) Set(value string) error
```

实现 flag.Value 接口，解析并设置 64 位整数值。

- 参数: `value` - 待解析的 64 位整数值
- 返回值: error - 解析错误或验证错误

```go
func (f *Int64Flag) SetRange(min, max int64)
```

设置 64 位整数的有效范围。

- 参数:
  - `min` - 最小值
  - `max` - 最大值

```go
func (f *Int64Flag) Type() FlagType
```

返回标志类型。

- 返回值: FlagType - 标志类型枚举值

### IntFlag 整数类型标志结构体

```go
type IntFlag struct {
    BaseFlag[int]
    // Has unexported fields.
}
```

IntFlag 是整数类型标志结构体，继承 BaseFlag[int] 泛型结构体，实现 Flag 接口。

#### 方法

```go
func (f *IntFlag) Set(value string) error
```

实现 flag.Value 接口，解析并验证整数值。

- 参数: `value` - 待解析的整数值
- 返回值: error - 解析错误或验证错误

```go
func (f *IntFlag) SetRange(min, max int)
```

设置整数的有效范围。

- 参数:
  - `min` - 最小值
  - `max` - 最大值

```go
func (f *IntFlag) String() string
```

实现 flag.Value 接口，返回当前整数值的字符串表示。

- 返回值: string - 当前整数值的字符串表示

```go
func (f *IntFlag) Type() FlagType
```

返回标志类型。

- 返回值: FlagType - 标志类型枚举值

### MapFlag 键值对类型标志结构体

```go
type MapFlag struct {
    BaseFlag[map[string]string]

    // Has unexported fields.
}
```

MapFlag 是键值对类型标志结构体，继承 BaseFlag[map[string]string] 泛型结构体，实现 Flag 接口。

#### 方法

```go
func (f *MapFlag) Set(value string) error
```

实现 flag.Value 接口，解析并设置键值对。

- 参数: `value` - 待设置的值
- 返回值: error - 解析或验证失败时返回错误信息

```go
func (f *MapFlag) SetDelimiters(keyDelimiter, valueDelimiter string)
```

设置键值对分隔符。

- 参数:
  - `keyDelimiter` - 键值对分隔符
  - `valueDelimiter` - 键值分隔符

```go
func (f *MapFlag) SetIgnoreCase(enable bool)
```

设置是否忽略键的大小写。

- 参数: `enable` - 是否忽略键的大小写
- 注意: 当 enable 为 true 时，所有键将转换为小写进行存储和比较

```go
func (f *MapFlag) String() string
```

实现 flag.Value 接口，返回当前值的字符串表示。

- 返回值: string - 当前值的字符串表示

```go
func (f *MapFlag) Type() FlagType
```

返回标志类型。

- 返回值: FlagType - 标志类型枚举值

### StringSliceFlag 字符串切片类型标志结构体

```go
type StringSliceFlag struct {
    BaseFlag[[]string] // 基类

    // Has unexported fields.
}
```

StringSliceFlag 是字符串切片类型标志结构体，继承 BaseFlag[[]string] 泛型结构体，实现 Flag 接口。

### IntSliceFlag 整数切片类型标志结构体

```go
type IntSliceFlag struct {
    BaseFlag[[]int] // 基类

    // Has unexported fields.
}
```

IntSliceFlag 是整数切片类型标志结构体，继承 BaseFlag[[]int] 泛型结构体，实现 Flag 接口。

### Int64SliceFlag 64位整数切片类型标志结构体

```go
type Int64SliceFlag struct {
    BaseFlag[[]int64] // 基类

    // Has unexported fields.
}
```

Int64SliceFlag 是64位整数切片类型标志结构体，继承 BaseFlag[[]int64] 泛型结构体，实现 Flag 接口。

### SizeFlag 大小类型标志结构体

```go
type SizeFlag struct {
    BaseFlag[int64]
    // Has unexported fields.
}
```

SizeFlag 是大小类型标志结构体，继承 BaseFlag[int64] 泛型结构体，实现 Flag 接口。

#### 方法

```go
func (f *SizeFlag) Init(longName, shortName string, defValue int64, usage string) error
```

初始化大小标志。

- 参数:
  - `longName` - 长标志名称
  - `shortName` - 短标志字符
  - `defValue` - 默认值（字节数）
  - `usage` - 帮助说明
- 返回值: `error` - 初始化错误信息

```go
func (f *SizeFlag) Set(value string) error
```

实现 `flag.Value` 接口，解析并设置大小值。

- 参数: `value` - 待设置的值
- 返回值: `error` - 解析或验证失败时返回错误信息

```go
func (f *SizeFlag) String() string
```

实现 `flag.Value` 接口，返回当前值的字符串表示。

- 返回值: `string` - 当前值的字符串表示

```go
func (f *SizeFlag) Type() FlagType
```

返回标志类型。

- 返回值: `FlagType` - 标志类型枚举值

```go
func (f *SizeFlag) SetAllowDecimal(allow bool) *SizeFlag
```

设置是否允许小数。

- 参数: `allow` - `true` 表示允许小数，`false` 表示不允许
- 返回值: `*SizeFlag` - 返回自身以支持链式调用

```go
func (f *SizeFlag) SetAllowNegative(allow bool) *SizeFlag
```

设置是否允许负数。

- 参数: `allow` - `true` 表示允许负数，`false` 表示不允许
- 返回值: `*SizeFlag` - 返回自身以支持链式调用

```go
func (f *SizeFlag) GetBytes() int64
```

获取字节数。

- 返回值: `int64` - 字节数

```go
func (f *SizeFlag) GetKiB() float64
```

获取KiB数。

- 返回值: `float64` - KiB数

```go
func (f *SizeFlag) GetMiB() float64
```

获取MiB数。

- 返回值: `float64` - MiB数

```go
func (f *SizeFlag) GetGiB() float64
```

获取GiB数。

- 返回值: `float64` - GiB数

```go
func (f *SizeFlag) GetTiB() float64
```

获取TiB数。

- 返回值: `float64` - TiB数

```go
func (f *SizeFlag) GetPiB() float64
```

获取PiB数。

- 返回值: `float64` - PiB数

```go
func (f *SizeFlag) IsZero() bool
```

检查是否为零。

- 返回值: `bool` - `true` 表示为零，`false` 表示非零

```go
func (f *SizeFlag) IsPositive() bool
```

检查是否为正数。

- 返回值: `bool` - `true` 表示为正数，`false` 表示非正数

```go
func (f *SizeFlag) IsNegative() bool
```

检查是否为负数。

- 返回值: `bool` - `true` 表示为负数，`false` 表示非负数

```go
func (f *SizeFlag) GetAllowDecimal() bool
```

获取是否允许小数设置。

- 返回值: `bool` - `true` 表示允许小数，`false` 表示不允许

```go
func (f *SizeFlag) GetAllowNegative() bool
```

获取是否允许负数设置。

- 返回值: `bool` - `true` 表示允许负数，`false` 表示不允许

#### StringSliceFlag 方法

```go
func (f *StringSliceFlag) Clear() error
```

清空字符串切片所有元素。

- 返回值: 操作成功返回 nil，否则返回错误信息
- 注意: 该方法会改变切片的指针

```go
func (f *StringSliceFlag) Contains(element string) bool
```

检查字符串切片是否包含指定元素。

- 参数: `element` - 待检查的元素
- 返回值: 若切片包含指定元素，返回 true，否则返回 false
- 注意: 当切片未设置值时，将使用默认值进行检查

```go
func (f *StringSliceFlag) GetDelimiters() []string
```

获取当前分隔符列表。

```go
func (f *StringSliceFlag) Init(longName, shortName string, defValue []string, usage string) error
```

初始化字符串切片类型标志。

- 参数:
  - `longName` - 长标志名称
  - `shortName` - 短标志字符
  - `defValue` - 默认值（字符串切片类型）
  - `usage` - 帮助说明
- 返回值: error - 初始化错误信息

```go
func (f *StringSliceFlag) Len() int
```

获取字符串切片长度。

- 返回值: 获取切片长度

```go
func (f *StringSliceFlag) Remove(element string) error
```

从字符串切片中移除指定元素（支持移除空字符串元素）。

- 参数: `element` - 待移除的元素（支持空字符串）
- 返回值: 操作成功返回 nil，否则返回错误信息

```go
func (f *StringSliceFlag) Set(value string) error
```

实现 flag.Value 接口，解析并设置字符串切片值。

- 参数: `value` - 待解析的切片值
- 注意: 如果切片中包含分隔符，则根据分隔符进行分割，否则将整个值作为单个元素。例如: "a,b,c" -> ["a", "b", "c"]

```go
func (f *StringSliceFlag) SetDelimiters(delimiters []string)
```

设置字符串切片解析的分隔符列表。

- 参数: `delimiters` - 分隔符列表

```go
func (f *StringSliceFlag) SetSkipEmpty(skip bool)
```

设置是否跳过空元素。

- 参数: `skip` - 为 true 时跳过空元素，为 false 时保留空元素
- 线程安全的空元素跳过更新

```go
func (f *StringSliceFlag) Sort() error
```

对字符串切片进行排序。

- 对当前切片标志的值进行原地排序，修改原切片内容
- 采用 Go 标准库的 sort.Strings() 函数进行字典序排序（按 Unicode 代码点升序排列）
- 注意:
  - 排序会直接修改当前标志的值，而非返回新切片
  - 排序区分大小写，遵循 Unicode 代码点比较规则（如 'A' < 'a' < 'z'）
  - 若切片未设置值，将使用默认值进行排序
- 返回值: 排序成功返回 nil，若排序过程中发生错误则返回错误信息

```go
func (f *StringSliceFlag) String() string
```

实现 flag.Value 接口，返回当前值的字符串表示。

```go
func (f *StringSliceFlag) Type() FlagType
```

返回标志类型。

#### IntSliceFlag 方法

```go
func (f *IntSliceFlag) Clear() error
```

清空整数切片所有元素。

- 返回值: 操作成功返回 nil，否则返回错误信息
- 注意: 该方法会改变切片的指针

```go
func (f *IntSliceFlag) Contains(element int) bool
```

检查整数切片是否包含指定元素。

- 参数: `element` - 待检查的元素
- 返回值: 若切片包含指定元素，返回 true，否则返回 false
- 注意: 当切片未设置值时，将使用默认值进行检查

```go
func (f *IntSliceFlag) GetDelimiters() []string
```

获取当前分隔符列表。

```go
func (f *IntSliceFlag) Init(longName, shortName string, defValue []int, usage string) error
```

初始化整数切片类型标志。

- 参数:
  - `longName` - 长标志名称
  - `shortName` - 短标志字符
  - `defValue` - 默认值（整数切片类型）
  - `usage` - 帮助说明
- 返回值: error - 初始化错误信息

```go
func (f *IntSliceFlag) Len() int
```

获取整数切片长度。

- 返回值: 获取切片长度

```go
func (f *IntSliceFlag) Remove(element int) error
```

从整数切片中移除指定元素。

- 参数: `element` - 待移除的元素
- 返回值: 操作成功返回 nil，否则返回错误信息

```go
func (f *IntSliceFlag) Set(value string) error
```

实现 flag.Value 接口，解析并设置整数切片值。

- 参数: `value` - 待解析的切片值
- 注意: 如果切片中包含分隔符，则根据分隔符进行分割，否则将整个值作为单个元素。例如: "1,2,3" -> [1, 2, 3]

```go
func (f *IntSliceFlag) SetDelimiters(delimiters []string)
```

设置整数切片解析的分隔符列表。

- 参数: `delimiters` - 分隔符列表

```go
func (f *IntSliceFlag) SetSkipEmpty(skip bool)
```

设置是否跳过空元素。

- 参数: `skip` - 为 true 时跳过空元素，为 false 时保留空元素
- 线程安全的空元素跳过更新

```go
func (f *IntSliceFlag) Sort() error
```

对整数切片进行排序。

- 对当前切片标志的值进行原地排序，修改原切片内容
- 采用 Go 标准库的 sort.Ints() 函数进行数值升序排序
- 注意:
  - 排序会直接修改当前标志的值，而非返回新切片
  - 若切片未设置值，将使用默认值进行排序
- 返回值: 排序成功返回 nil，若排序过程中发生错误则返回错误信息

```go
func (f *IntSliceFlag) String() string
```

实现 flag.Value 接口，返回当前值的字符串表示。

```go
func (f *IntSliceFlag) Type() FlagType
```

返回标志类型。

#### Int64SliceFlag 方法

```go
func (f *Int64SliceFlag) Clear() error
```

清空64位整数切片所有元素。

- 返回值: 操作成功返回 nil，否则返回错误信息
- 注意: 该方法会改变切片的指针

```go
func (f *Int64SliceFlag) Contains(element int64) bool
```

检查64位整数切片是否包含指定元素。

- 参数: `element` - 待检查的元素
- 返回值: 若切片包含指定元素，返回 true，否则返回 false
- 注意: 当切片未设置值时，将使用默认值进行检查

```go
func (f *Int64SliceFlag) GetDelimiters() []string
```

获取当前分隔符列表。

```go
func (f *Int64SliceFlag) Init(longName, shortName string, defValue []int64, usage string) error
```

初始化64位整数切片类型标志。

- 参数:
  - `longName` - 长标志名称
  - `shortName` - 短标志字符
  - `defValue` - 默认值（64位整数切片类型）
  - `usage` - 帮助说明
- 返回值: error - 初始化错误信息

```go
func (f *Int64SliceFlag) Len() int
```

获取64位整数切片长度。

- 返回值: 获取切片长度

```go
func (f *Int64SliceFlag) Remove(element int64) error
```

从64位整数切片中移除指定元素。

- 参数: `element` - 待移除的元素
- 返回值: 操作成功返回 nil，否则返回错误信息

```go
func (f *Int64SliceFlag) Set(value string) error
```

实现 flag.Value 接口，解析并设置64位整数切片值。

- 参数: `value` - 待解析的切片值
- 注意: 如果切片中包含分隔符，则根据分隔符进行分割，否则将整个值作为单个元素。例如: "1,2,3" -> [1, 2, 3]

```go
func (f *Int64SliceFlag) SetDelimiters(delimiters []string)
```

设置64位整数切片解析的分隔符列表。

- 参数: `delimiters` - 分隔符列表

```go
func (f *Int64SliceFlag) SetSkipEmpty(skip bool)
```

设置是否跳过空元素。

- 参数: `skip` - 为 true 时跳过空元素，为 false 时保留空元素
- 线程安全的空元素跳过更新

```go
func (f *Int64SliceFlag) Sort() error
```

对64位整数切片进行排序。

- 对当前切片标志的值进行原地排序，修改原切片内容
- 采用自定义排序函数进行数值升序排序
- 注意:
  - 排序会直接修改当前标志的值，而非返回新切片
  - 若切片未设置值，将使用默认值进行排序
- 返回值: 排序成功返回 nil，若排序过程中发生错误则返回错误信息

```go
func (f *Int64SliceFlag) String() string
```

实现 flag.Value 接口，返回当前值的字符串表示。

```go
func (f *Int64SliceFlag) Type() FlagType
```

返回标志类型。

### StringFlag 字符串类型标志结构体

```go
type StringFlag struct {
    BaseFlag[string]
}
```

StringFlag 是字符串类型标志结构体，继承 BaseFlag[string] 泛型结构体，实现 Flag 接口。

#### 方法

```go
func (f *StringFlag) Contains(substr string) bool
```

检查字符串是否包含指定子串。

- 参数: `substr` - 子串
- 返回值: bool - 如果包含子串则返回 true，否则返回 false

```go
func (f *StringFlag) Len() int
```

获取字符串标志的长度。

- 返回值: 字符串的字符数（按 UTF-8 编码计算）

```go
func (f *StringFlag) Set(value string) error
```

实现 flag.Value 接口的 Set 方法，将字符串值解析并设置到标志中。

- 参数: `value` - 待设置的字符串值
- 返回值: error - 设置失败时返回错误信息

```go
func (f *StringFlag) String() string
```

返回带引号的字符串值。

- 返回值: string - 带引号的字符串值

```go
func (f *StringFlag) ToLower() string
```

将字符串标志值转换为小写。

```go
func (f *StringFlag) ToUpper() string
```

将字符串标志值转换为大写。

```go
func (f *StringFlag) Type() FlagType
```

返回标志类型。

- 返回值: FlagType - 标志类型枚举值

### TimeFlag 时间类型标志结构体

```go
type TimeFlag struct {
    BaseFlag[time.Time]

    // Has unexported fields.
}
```

TimeFlag 是时间类型标志结构体，继承 BaseFlag[time.Time] 泛型结构体，实现 Flag 接口。

#### 方法

```go
func (f *TimeFlag) Init(longName, shortName string, defValue string, usage string) error
```

初始化时间类型标志，支持字符串类型默认值。

- 参数:
  - `longName` - 长标志名称
  - `shortName` - 短标志字符
  - `defValue` - 默认值（字符串格式，支持多种时间表达）
  - `usage` - 帮助说明
- 返回值: error - 初始化错误信息
- 支持的默认值格式:
  - "now" 或 "" : 当前时间
  - "zero" : 零时间 (time.Time{})
  - "1h", "30m", "-2h" : 相对时间（基于当前时间的偏移）
  - "2006-01-02", "2006-01-02 15:04:05" : 绝对时间格式
  - RFC3339 等标准格式

```go
func (f *TimeFlag) Set(value string) error
```

实现 flag.Value 接口，解析并设置时间值。

- 参数: `value` - 待解析的时间字符串
- 返回值: error - 解析或验证失败时返回错误信息

```go
func (f *TimeFlag) SetOutputFormat(format string)
```

设置时间输出格式。

- 参数: `format` - 时间格式化字符串
- 注意: 此方法线程安全

```go
func (f *TimeFlag) String() string
```

实现 flag.Value 接口，返回当前时间的字符串表示。

- 返回值: string - 格式化后的时间字符串
- 注意: 加锁保证 outputFormat 和 value 的并发安全访问

```go
func (f *TimeFlag) Type() FlagType
```

返回标志类型。

- 返回值: FlagType - 标志类型枚举值

### TypedFlag 泛型标志接口

```go
type TypedFlag[T any] interface {
    Flag                           // 继承标志接口
    GetDefault() T                 // 获取标志的具体类型默认值
    Get() T                        // 获取标志的具体类型值
    GetPointer() *T                // 获取标志值的指针
    Set(T) error                   // 设置标志的具体类型值
    SetValidator(Validator)        // 设置标志的验证器
    BindEnv(envName string) *BaseFlag[T] // 绑定环境变量
}
```

TypedFlag 是所有标志类型的通用接口，定义了标志的元数据访问方法和默认值访问方法。

### Uint16Flag 16 位无符号整数类型标志结构体

```go
type Uint16Flag struct {
    BaseFlag[uint16] // 基类
    // Has unexported fields.
}
```

Uint16Flag 是 16 位无符号整数类型标志结构体，继承 BaseFlag[uint16] 泛型结构体，实现 Flag 接口。

#### 方法

```go
func (f *Uint16Flag) Set(value string) error
```

实现 flag.Value 接口，解析并设置 16 位无符号整数值。

- 验证值是否在 uint16 范围内 (0-65535)
- 参数: `value` - 待设置的值 (0-65535)
- 返回值: error - 解析或验证失败时返回错误信息

```go
func (f *Uint16Flag) String() string
```

实现 flag.Value 接口，返回当前值的字符串表示。

- 返回值: string - 当前值的字符串表示

```go
func (f *Uint16Flag) Type() FlagType
```

返回标志类型。

- 返回值: FlagType - 标志类型枚举值

### Uint32Flag 32 位无符号整数类型标志结构体

```go
type Uint32Flag struct {
    BaseFlag[uint32] // 基类
    // Has unexported fields.
}
```

Uint32Flag 是 32 位无符号整数类型标志结构体，继承 BaseFlag[uint32] 泛型结构体，实现 Flag 接口。

#### 方法

```go
func (f *Uint32Flag) Set(value string) error
```

实现 flag.Value 接口，解析并设置 32 位无符号整数值。

- 验证值是否在 uint32 范围内 (0-4294967295)
- 参数: `value` - 待设置的值 (0-4294967295)
- 返回值: error - 解析或验证失败时返回错误信息

```go
func (f *Uint32Flag) String() string
```

实现 flag.Value 接口，返回当前值的字符串表示。

```go
func (f *Uint32Flag) Type() FlagType
```

返回标志类型。

### Uint64Flag 64 位无符号整数类型标志结构体

```go
type Uint64Flag struct {
    BaseFlag[uint64] // 基类
    // Has unexported fields.
}
```

Uint64Flag 是 64 位无符号整数类型标志结构体，继承 BaseFlag[uint64] 泛型结构体，实现 Flag 接口。

#### 方法

```go
func (f *Uint64Flag) Set(value string) error
```

实现 flag.Value 接口，解析并设置 64 位无符号整数值。

- 验证值是否在 uint64 范围内 (0-18446744073709551615)
- 参数: `value` - 待设置的值 (0-18446744073709551615)
- 返回值: error - 解析或验证失败时返回错误信息

```go
func (f *Uint64Flag) String() string
```

实现 flag.Value 接口，返回当前值的字符串表示。

```go
func (f *Uint64Flag) Type() FlagType
```

返回标志类型。

### Validator 验证器接口

```go
type Validator interface {
    // Validate 验证参数值是否合法
    // value: 待验证的参数值
    // 返回值: 验证通过返回 nil，否则返回错误信息
    Validate(value any) error
}
```

Validator 是验证器接口，所有自定义验证器需实现此接口。