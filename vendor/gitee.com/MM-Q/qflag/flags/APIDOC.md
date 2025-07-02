# Package flags

flags 定义了所有标志类型的通用接口和基础标志结构体

## CONSTANTS

```go
const (
	FlagSplitComma     = "," // 逗号
	FlagSplitSemicolon = ";" // 分号
	FlagSplitPipe      = "|" // 竖线
	FlagKVColon        = ":" // 冒号
	FlagKVEqual        = "=" // 等号
)
```

定义标志的分隔符常量。

```go
const InvalidFlagChars = " !@#$%^&*(){}[]|\\;:'\"<>,.?/"
```

定义非法字符集常量，防止非法的标志名称。

## VARIABLES

```go
var (
	HelpFlagName            = "help"    // 帮助标志名称
	HelpFlagShortName       = "h"       // 帮助标志短名称
	ShowInstallPathFlagName = "sip"     // 显示安装路径标志名称
	VersionFlagLongName     = "version" // 版本标志名称
	VersionFlagShortName    = "v"       // 版本标志短名称
)
```

内置标志名称。

```go
var FlagSplitSlice = []string{
	FlagSplitComma,
	FlagSplitSemicolon,
	FlagSplitPipe,
	FlagKVColon,
}
```

Flag支持的标志分隔符切片。

## FUNCTIONS

```go
func FlagTypeToString(flagType FlagType) string
```

FlagTypeToString 将FlagType转换为字符串。

## TYPES

```go
type BaseFlag[T any] struct {
	// Has unexported fields.
}
```

BaseFlag 泛型基础标志结构体，封装所有标志的通用字段和方法。

```go
func (f *BaseFlag[T]) Get() T
```

Get 获取标志的实际值 优先级：已设置的值 > 默认值 线程安全：使用互斥锁保证并发访问安全。

```go
func (f *BaseFlag[T]) GetDefault() T
```

GetDefault 获取标志的默认值。

```go
func (f *BaseFlag[T]) GetDefaultAny() any
```

GetDefaultAny 获取标志的默认值(any类型)。

```go
func (f *BaseFlag[T]) GetPointer() *T
```

GetPointer 返回标志值的指针。

注意:

1. 获取指针过程受锁保护, 但直接修改指针指向的值仍会绕过验证机制
2. 多线程环境下修改时需额外同步措施, 建议优先使用Set()方法

```go
func (f *BaseFlag[T]) Init(longName, shortName string, defValue T, usage string, value *T) error
```

Init 初始化标志的元数据和值指针, 无需显式调用, 仅在创建标志对象时自动调用。

参数:

  - longName: 长标志名称
  - shortName: 短标志字符
  - defValue: 默认值
  - usage: 帮助说明
  - value: 标志值指针

返回值:

  - error: 初始化错误信息

```go
func (f *BaseFlag[T]) IsSet() bool
```

IsSet 判断标志是否已被设置值。

返回值: true表示已设置值, false表示未设置。

```go
func (f *BaseFlag[T]) LongName() string
```

LongName 获取标志的长名称。

```go
func (f *BaseFlag[T]) Reset()
```

Reset 将标志值重置为默认值 线程安全：使用互斥锁保证并发安全。

```go
func (f *BaseFlag[T]) Set(value T) error
```

Set 设置标志的值。

参数: value 标志值。

返回: 错误信息。

```go
func (f *BaseFlag[T]) SetValidator(validator Validator)
```

SetValidator 设置标志的验证器。

参数: validator 验证器接口。

```go
func (f *BaseFlag[T]) ShortName() string
```

ShortName 获取标志的短名称。

```go
func (f *BaseFlag[T]) String() string
```

String 返回标志的字符串表示。

```go
func (f *BaseFlag[T]) Usage() string
```

Usage 获取标志的用法说明。

```go
type BoolFlag struct {
	BaseFlag[bool]
}
```

BoolFlag 布尔类型标志结构体 继承BaseFlag[bool]泛型结构体,实现Flag接口。

```go
func (f *BoolFlag) Type() FlagType
```

Type 返回标志类型。

```go
type DurationFlag struct {
	BaseFlag[time.Duration]
}
```

DurationFlag 时间间隔类型标志结构体 继承BaseFlag[time.Duration]泛型结构体,实现Flag接口。

```go
func (f *DurationFlag) Set(value string) error
```

Set 实现flag.Value接口, 解析并设置时间间隔值。

```go
func (f *DurationFlag) String() string
```

String 实现flag.Value接口, 返回当前值的字符串表示。

```go
func (f *DurationFlag) Type() FlagType
```

Type 返回标志类型。

```go
type EnumFlag struct {
	BaseFlag[string]
	// Has unexported fields.
}
```

EnumFlag 枚举类型标志结构体 继承BaseFlag[string]泛型结构体,增加枚举特有的选项验证。

```go
func (f *EnumFlag) Init(longName, shortName string, defValue string, usage string, options []string) error
```

Init 初始化枚举类型标志, 无需显式调用, 仅在创建标志对象时自动调用。

参数:

  - longName: 长标志名称
  - shortName: 短标志字符
  - defValue: 默认值
  - usage: 帮助说明
  - options: 枚举可选值列表

返回值:

  - error: 初始化错误信息

```go
func (f *EnumFlag) IsCheck(value string) error
```

IsCheck 检查枚举值是否有效 返回值: 为nil, 说明值有效,否则返回错误信息。

```go
func (f *EnumFlag) Set(value string) error
```

Set 实现flag.Value接口, 解析并设置枚举值。

```go
func (f *EnumFlag) String() string
```

String 实现flag.Value接口, 返回当前值的字符串表示。

```go
func (f *EnumFlag) Type() FlagType
```

实现Flag接口。

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
}
```

Flag 所有标志类型的通用接口,定义了标志的元数据访问方法。

```go
type FlagMeta struct {
	Flag Flag // 标志对象
}
```

FlagMeta 统一存储标志的完整元数据。

```go
func (m *FlagMeta) GetDefault() any
```

GetDefault 获取标志的默认值。

```go
func (m *FlagMeta) GetFlag() Flag
```

GetFlag 获取标志对象。

```go
func (m *FlagMeta) GetFlagType() FlagType
```

GetFlagType 获取标志的类型。

```go
func (m *FlagMeta) GetLongName() string
```

GetLongName 获取标志的长名称。

```go
func (m *FlagMeta) GetShortName() string
```

GetShortName 获取标志的短名称。

```go
func (m *FlagMeta) GetUsage() string
```

GetUsage 获取标志的用法描述。

```go
type FlagMetaInterface interface {
	GetFlagType() FlagType // 获取标志类型
	GetFlag() Flag         // 获取标志对象
	GetLongName() string   // 获取标志的长名称
	GetShortName() string  // 获取标志的短名称
	GetUsage() string      // 获取标志的用法描述
	GetDefault() any       // 获取标志的默认值
	GetValue() any         // 获取标志的当前值
}
```

FlagMetaInterface 标志元数据接口, 定义了标志元数据的获取方法。

```go
type FlagRegistry struct {
	// Has unexported fields.
}
```

FlagRegistry 集中管理所有标志元数据及索引。

```go
func NewFlagRegistry() *FlagRegistry
```

创建一个空的标志注册表。

```go
func (r *FlagRegistry) GetAllFlags() []*FlagMeta
```

GetAllFlags 获取所有标志元数据列表 返回值:

  - []*FlagMeta: 所有标志元数据的切片

```go
func (r *FlagRegistry) GetByLong(longName string) (*FlagMeta, bool)
```

GetByLong 通过长标志名称查找对应的标志元数据 参数:

  - longName: 标志的长名称(如"help")

返回值:

  - *FlagMeta: 找到的标志元数据指针, 未找到时为nil
  - bool: 是否找到标志, true表示找到

```go
func (r *FlagRegistry) GetByName(name string) (*FlagMeta, bool)
```

GetByName 通过标志名称查找标志元数据 参数name可以是长名称(如"help")或短名称(如"h") 返回值:

  - *FlagMeta: 找到的标志元数据指针, 未找到时为nil
  - bool: 是否找到标志, true表示找到

```go
func (r *FlagRegistry) GetByShort(shortName string) (*FlagMeta, bool)
```

GetByShort 通过短标志名称查找对应的标志元数据 参数:

  - shortName: 标志的短名称(如"h"对应"help")

返回值:

  - *FlagMeta: 找到的标志元数据指针, 未找到时为nil
  - bool: 是否找到标志, true表示找到

```go
func (r *FlagRegistry) GetLongFlags() map[string]*FlagMeta
```

GetLongFlags 获取长标志映射 返回值:

  - map[string]*FlagMeta: 长标志名称到标志元数据的映射

```go
func (r *FlagRegistry) GetShortFlags() map[string]*FlagMeta
```

GetShortFlags 获取短标志映射 返回值:

  - map[string]*FlagMeta: 短标志名称到标志元数据的映射

```go
func (r *FlagRegistry) RegisterFlag(meta *FlagMeta) error
```

RegisterFlag 注册一个新的标志元数据到注册表中 该方法会执行以下操作: 1. 检查长名称和短名称是否已存在 2. 将标志添加到长名称索引
3. 将标志添加到短名称索引 4. 将标志添加到所有标志列表 注意: 该方法线程安全, 但发现重复标志时会panic。

```go
type FlagRegistryInterface interface {
	GetAllFlags() []*FlagMeta                      // 获取所有标志元数据列表
	GetLongFlags() map[string]*FlagMeta            // 获取长标志映射
	GetShortFlags() map[string]*FlagMeta           // 获取短标志映射
	RegisterFlag(meta *FlagMeta) error             // 注册一个新的标志元数据到注册表中
	GetByLong(longName string) (*FlagMeta, bool)   // 通过长标志名称查找对应的标志元数据
	GetByShort(shortName string) (*FlagMeta, bool) // 通过短标志名称查找对应的标志元数据
	GetByName(name string) (*FlagMeta, bool)       // 通过标志名称查找标志元数据
}
```

FlagRegistryInterface 标志注册表接口, 定义了标志元数据的增删改查操作。

```go
type FlagType int
```

标志类型。

```go
const (
	FlagTypeInt      FlagType = iota + 1 // 整数类型
	FlagTypeInt64                        // 64位整数类型
	FlagTypeUint16                       // 16位无符号整数类型
	FlagTypeString                       // 字符串类型
	FlagTypeBool                         // 布尔类型
	FlagTypeFloat64                      // 64位浮点数类型
	FlagTypeEnum                         // 枚举类型
	FlagTypeDuration                     // 时间间隔类型
	FlagTypeSlice                        // 切片类型
	FlagTypeTime                         // 时间类型
	FlagTypeMap                          // 映射类型
	FlagTypePath                         // 路径类型
)
```

```go
type Float64Flag struct {
	BaseFlag[float64]
}
```

Float64Flag 浮点型标志结构体 继承BaseFlag[float64]泛型结构体,实现Flag接口。

```go
func (f *Float64Flag) Type() FlagType
```

Type 返回标志类型。

```go
type Int64Flag struct {
	BaseFlag[int64]
}
```

Int64Flag 64位整数类型标志结构体 继承BaseFlag[int64]泛型结构体,实现Flag接口。

```go
func (f *Int64Flag) SetRange(min, max int64)
```

SetRange 设置64位整数的有效范围。

min: 最小值 max: 最大值。

```go
func (f *Int64Flag) Type() FlagType
```

Type 返回标志类型。

```go
type IntFlag struct {
	BaseFlag[int]
}
```

IntFlag 整数类型标志结构体 继承BaseFlag[int]泛型结构体,实现Flag接口。

```go
func (f *IntFlag) SetRange(min, max int)
```

SetRange 设置整数的有效范围。

min: 最小值 max: 最大值。

```go
func (f *IntFlag) Type() FlagType
```

Type 返回标志类型。

```go
type MapFlag struct {
	BaseFlag[map[string]string]

	// Has unexported fields.
}
```

MapFlag 键值对类型标志结构体 继承BaseFlag[map[string]string]泛型结构体,实现Flag接口。

```go
func (f *MapFlag) Set(value string) error
```

Set 实现flag.Value接口,解析并设置键值对。

```go
func (f *MapFlag) SetDelimiters(keyDelimiter, valueDelimiter string)
```

SetDelimiters 设置键值对分隔符。

参数：

  - keyDelimiter 键值对分隔符
  - valueDelimiter 键值分隔符

```go
func (f *MapFlag) SetIgnoreCase(enable bool)
```

SetIgnoreCase 设置是否忽略键的大小写 enable为true时，所有键将转换为小写进行存储和比较。

```go
func (f *MapFlag) String() string
```

String 实现flag.Value接口,返回当前值的字符串表示。

```go
func (f *MapFlag) Type() FlagType
```

Type 返回标志类型。

```go
type PathFlag struct {
	BaseFlag[string]
}
```

PathFlag 路径类型标志结构体 继承BaseFlag[string]泛型结构体,实现Flag接口。

```go
func (f *PathFlag) Init(longName, shortName string, defValue string, usage string) error
```

Init 初始化路径标志。

```go
func (f *PathFlag) Set(value string) error
```

Set 实现flag.Value接口,解析并验证路径。

```go
func (f *PathFlag) String() string
```

String 实现flag.Value接口,返回当前值的字符串表示。

```go
func (f *PathFlag) Type() FlagType
```

Type 返回标志类型。

```go
type SliceFlag struct {
	BaseFlag[[]string] // 基类

	// Has unexported fields.
}
```

SliceFlag 切片类型标志结构体 继承BaseFlag[[]string]泛型结构体,实现Flag接口。

```go
func (f *SliceFlag) Clear() error
```

Clear 清空切片所有元素。

返回值: 操作成功返回nil, 否则返回错误信息。

注意：该方法会改变切片的指针。

```go
func (f *SliceFlag) Contains(element string) bool
```

Contains 检查切片是否包含指定元素 当切片未设置值时,将使用默认值进行检查。

```go
func (f *SliceFlag) GetDelimiters() []string
```

GetDelimiters 获取当前分隔符列表。

```go
func (f *SliceFlag) Init(longName, shortName string, defValue []string, usage string) error
```

Init 初始化切片类型标志。

参数:

  - longName: 长标志名称
  - shortName: 短标志字符
  - defValue: 默认值（切片类型）
  - usage: 帮助说明

返回值:

  - error: 初始化错误信息

```go
func (f *SliceFlag) Len() int
```

Len 获取切片长度。

返回: 获取切片长度。

```go
func (f *SliceFlag) Remove(element string) error
```

Remove 从切片中移除指定元素（支持移除空字符串元素）。

参数: element 待移除的元素（支持空字符串）。

返回值: 操作成功返回nil, 否则返回错误信息。

```go
func (f *SliceFlag) Set(value string) error
```

Set 实现flag.Value接口, 解析并设置切片值。

参数: value 待解析的切片值。

注意: 如果切片中包含分隔符,则根据分隔符进行分割, 否则将整个值作为单个元素 例如: "a,b,c" -> ["a", "b", "c"]。

```go
func (f *SliceFlag) SetDelimiters(delimiters []string)
```

SetDelimiters 设置切片解析的分隔符列表。

参数: delimiters 分隔符列表。

线程安全的分隔符更新。

```go
func (f *SliceFlag) SetSkipEmpty(skip bool)
```

SetSkipEmpty 设置是否跳过空元素。

参数: skip - 为true时跳过空元素, 为false时保留空元素。

```go
func (f *SliceFlag) Sort() error
```

Sort 对切片进行排序。

功能：对当前切片标志的值进行原地排序，修改原切片内容 排序规则:
采用Go标准库的sort.Strings()函数进行字典序排序(按Unicode代码点升序排列) 注意事项：
 1. 排序会直接修改当前标志的值，而非返回新切片
 2. 排序区分大小写, 遵循Unicode代码点比较规则(如'A' < 'a' < 'z')
 3. 若切片未设置值，将使用默认值进行排序。

返回值：

排序成功返回nil, 若排序过程中发生错误则返回错误信息。

```go
func (f *SliceFlag) String() string
```

String 实现flag.Value接口, 返回当前值的字符串表示。

```go
func (f *SliceFlag) Type() FlagType
```

Type 返回标志类型。

```go
type StringFlag struct {
	BaseFlag[string]
}
```

StringFlag 字符串类型标志结构体 继承BaseFlag[string]泛型结构体,实现Flag接口。

```go
func (f *StringFlag) Contains(substr string) bool
```

Contains 检查字符串是否包含指定子串。

参数: substr 子串。

返回值: true表示包含, false表示不包含。

```go
func (f *StringFlag) Len() int
```

Len 获取字符串标志的长度。

返回值：字符串的字符数(按UTF-8编码计算)。

```go
func (f *StringFlag) String() string
```

String 返回带引号的字符串值。

```go
func (f *StringFlag) ToLower() string
```

ToLower 将字符串标志值转换为小写。

```go
func (f *StringFlag) ToUpper() string
```

ToUpper 将字符串标志值转换为大写。

```go
func (f *StringFlag) Type() FlagType
```

Type 返回标志类型。

```go
type TimeFlag struct {
	BaseFlag[time.Time]

	// Has unexported fields.
}
```

TimeFlag 时间类型标志结构体 继承BaseFlag[time.Time]泛型结构体,实现Flag接口。

```go
func (f *TimeFlag) Set(value string) error
```

Set 实现flag.Value接口, 解析并设置时间值。

```go
func (f *TimeFlag) SetOutputFormat(format string)
```

SetOutputFormat 设置时间输出格式。

```go
func (f *TimeFlag) String() string
```

String 实现flag.Value接口, 返回当前时间的字符串表示 加锁保证outputFormat和value的并发安全访问。

```go
func (f *TimeFlag) Type() FlagType
```

Type 返回标志类型。

```go
type TypedFlag[T any] interface {
	Flag                    // 继承标志接口
	GetDefault() T          // 获取标志的具体类型默认值
	Get() T                 // 获取标志的具体类型值
	GetPointer() *T         // 获取标志值的指针
	Set(T) error            // 设置标志的具体类型值
	SetValidator(Validator) // 设置标志的验证器
}
```

TypedFlag 所有标志类型的通用接口,定义了标志的元数据访问方法和默认值访问方法。

```go
type Uint16Flag struct {
	BaseFlag[uint16] // 基类
	// Has unexported fields.
}
```

Uint16Flag 16位无符号整数类型标志结构体 继承BaseFlag[uint16]泛型结构体,实现Flag接口。

```go
func (f *Uint16Flag) Set(value string) error
```

Set 实现flag.Value接口, 解析并设置16位无符号整数值 验证值是否在uint16范围内(0-65535)。

参数:

    value: 待设置的值(0-65535)

返回值:

    error: 错误信息

```go
func (f *Uint16Flag) String() string
```

String 实现flag.Value接口, 返回当前值的字符串表示。

```go
func (f *Uint16Flag) Type() FlagType
```

Type 返回标志类型。

```go
type Validator interface {
	// Validate 验证参数值是否合法
	// value: 待验证的参数值
	// 返回值: 验证通过返回nil, 否则返回错误信息
	Validate(value any) error
}
```

Validator 验证器接口, 所有自定义验证器需实现此接口。