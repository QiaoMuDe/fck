# Package flag 

```go
import "gitee.com/MM-Q/qflag/internal/flag"
```

---

## CONSTANTS

### const (IntSize, UintSize)

```go
const (
    // IntSize 当前平台上int类型的位数
    // 在32位系统上为32, 在64位系统上为64
    IntSize = strconv.IntSize
    // UintSize 当前平台上uint类型的位数, 与int相同
    // 在32位系统上为32, 在64位系统上为64
    UintSize = strconv.IntSize
)
```

平台相关的整数位数

---

## TYPES

### type BaseFlag[T any] struct

```go
type BaseFlag[T any] struct {
    // Has unexported fields.
}
```

BaseFlag 泛型基础标志结构体

BaseFlag 是所有标志类型的基础结构, 使用泛型支持多种数据类型。 它提供了标志的基本功能, 包括名称管理、值存储、默认值处理和环境变量绑定等。

**线程安全: **
  - 所有公共方法都使用读写锁保护, 确保并发安全
  - 读操作使用读锁, 写操作使用写锁

**字段说明: **
  - mu: 读写锁, 保护并发访问
  - longName: 长选项名称, 如 "--help"
  - shortName: 短选项名称, 如 "-h"
  - desc: 标志描述信息
  - flagType: 标志类型枚举值
  - value: 指向当前值的指针
  - default_: 默认值
  - isSet: 标志是否已被设置
  - envVar: 关联的环境变量名

#### func NewBaseFlag[T any](flagType types.FlagType, longName, shortName, desc string, default_ T) *BaseFlag[T]

```go
func NewBaseFlag[T any](flagType types.FlagType, longName, shortName, desc string, default_ T) *BaseFlag[T]
```

NewBaseFlag 创建新的基础标志实例

**参数:**
  - flagType: 标志类型枚举值
  - longName: 长选项名, 如 "help"
  - shortName: 短选项名, 如 "h"
  - default_: 默认值

**返回值:**
  - *BaseFlag[T]: 基础标志实例

**注意事项: **
  - 此函数会初始化内部值指针, 并将默认值复制到值中
  - 创建后的标志初始状态为未设置(isSet=false)

#### func (f *BaseFlag[T]) BindEnv(name string)

```go
func (f *BaseFlag[T]) BindEnv(name string)
```

BindEnv 绑定环境变量

**参数:**
  - name: 环境变量名

**注意事项: **
  - 绑定后, 标志可以从指定的环境变量读取值
  - 环境变量的优先级低于命令行参数

#### func (f *BaseFlag[T]) Desc() string

```go
func (f *BaseFlag[T]) Desc() string
```

Desc 获取标志的描述信息

**返回值:**
  - string: 标志描述文本

#### func (f *BaseFlag[T]) EnumValues() []string

```go
func (f *BaseFlag[T]) EnumValues() []string
```

EnumValues 获取枚举类型的可选值

**返回值:**
  - []string: 枚举类型的可选值列表

**功能说明: **
  - 实现 Flag 接口的 EnumValues 方法
  - 非枚举类型返回空切片
  - 此方法是线程安全的

#### func (f *BaseFlag[T]) Get() T

```go
func (f *BaseFlag[T]) Get() T
```

Get 获取标志的当前值

**返回值:**
  - T: 标志的当前值

#### func (f *BaseFlag[T]) GetDef() any

```go
func (f *BaseFlag[T]) GetDef() any
```

GetDef 获取标志的默认值

**返回值:**
  - any: 标志的默认值, 使用any类型以支持泛型

#### func (f *BaseFlag[T]) GetEnvVar() string

```go
func (f *BaseFlag[T]) GetEnvVar() string
```

GetEnvVar 获取关联的环境变量名

**返回值:**
  - string: 环境变量名, 如果未绑定则返回空字符串

#### func (f *BaseFlag[T]) GetStr() string

```go
func (f *BaseFlag[T]) GetStr() string
```

GetStr 获取标志当前值的字符串表示

**返回值:**
  - string: 标志当前值的字符串表示

**功能说明: **
  - 获取标志当前值的字符串表示
  - 与String()方法不同, 此方法专注于值本身
  - 用于内置标志处理中获取标志值

#### func (f *BaseFlag[T]) GetValuePtr() *T

```go
func (f *BaseFlag[T]) GetValuePtr() *T
```

GetValuePtr 返回值指针, 用于注册到标准库 flag 包

**返回值:**
  - *T: 指向标志值的指针

**注意事项: ** 1. 此方法主要用于与标准库 flag 包集成, 不推荐在常规代码中使用 2. 返回的指针指向内部状态, 直接修改可能破坏线程安全 3. 仅应在程序初始化阶段 (标志注册时) 使用, 避免并发访问 4. 如需在多线程环境中访问标志值, 请使用 Get() 方法

#### func (f *BaseFlag[T]) IsSet() bool

```go
func (f *BaseFlag[T]) IsSet() bool
```

IsSet 检查标志是否已被设置

**返回值:**
  - bool: 如果标志已被设置返回true, 否则返回false

#### func (f *BaseFlag[T]) LongName() string

```go
func (f *BaseFlag[T]) LongName() string
```

LongName 获取标志的长名称

**返回值:**
  - string: 长选项名称, 如 "help"

#### func (f *BaseFlag[T]) Name() string

```go
func (f *BaseFlag[T]) Name() string
```

Name 获取标志名称

**返回值:**
  - string: 优先返回长名称, 如果长名称为空则返回短名称

#### func (f *BaseFlag[T]) Reset()

```go
func (f *BaseFlag[T]) Reset()
```

Reset 重置标志为默认值

将标志的值重置为默认值, 并将isSet状态设置为false

#### func (f *BaseFlag[T]) Set(value string) error

```go
func (f *BaseFlag[T]) Set(value string) error
```

Set 设置标志的值

**参数:**
  - value: 要设置的字符串值

**返回值:**
  - error: 如果设置失败返回错误

**注意事项: **
  - 这是基础实现, 具体子类应该重写此方法实现自己的解析逻辑
  - 基础实现仅返回nil, 不进行任何实际设置操作

#### func (f *BaseFlag[T]) ShortName() string

```go
func (f *BaseFlag[T]) ShortName() string
```

ShortName 获取标志的短名称

**返回值:**
  - string: 短选项名称, 如 "h"

#### func (f *BaseFlag[T]) String() string

```go
func (f *BaseFlag[T]) String() string
```

String 返回标志的格式化名称

**返回值:**
  - string: 格式化的标志名称, 用于显示

**注意事项: **
  - 使用utils.FormatFlagName进行格式化
  - 通常用于帮助信息显示

#### func (f *BaseFlag[T]) Type() types.FlagType

```go
func (f *BaseFlag[T]) Type() types.FlagType
```

Type 获取标志的类型

**返回值:**
  - types.FlagType: 标志类型枚举值

#### func (f *BaseFlag[T]) SetValidator(validator types.Validator[T])

```go
func (f *BaseFlag[T]) SetValidator(validator types.Validator[T])
```

SetValidator 设置验证器

**参数:**
  - validator: 验证器函数

**功能说明:**
  - 设置标志的验证器
  - 如果之前已设置验证器，会被覆盖
  - 验证器会在 Set 方法中解析完值后被调用
  - 如果验证失败，Set 方法会返回错误，标志值不会被设置

**使用示例:**
```go
// 端口号验证：1-65535
port.SetValidator(func(value int) error {
    if value < 1 || value > 65535 {
        return fmt.Errorf("端口 %d 超出范围 [1, 65535]", value)
    }
    return nil
})
```

#### func (f *BaseFlag[T]) ClearValidator()

```go
func (f *BaseFlag[T]) ClearValidator()
```

ClearValidator 清除验证器

**功能说明:**
  - 移除标志的验证器
  - 之后调用 Set 方法将不会进行验证

#### func (f *BaseFlag[T]) HasValidator() bool

```go
func (f *BaseFlag[T]) HasValidator() bool
```

HasValidator 检查是否设置了验证器

**返回值:**
  - bool: 是否设置了验证器

**功能说明:**
  - 用于判断标志是否配置了验证逻辑
  - 此方法是线程安全的

---

### type BoolFlag struct

```go
type BoolFlag struct {
    *BaseFlag[bool]
}
```

BoolFlag 布尔标志

BoolFlag 用于处理布尔类型的命令行参数。 它接受多种布尔值表示形式, 包括 "true", "false", "1", "0", "t", "f", "TRUE", "FALSE" 等。

#### func NewBoolFlag(longName, shortName, desc string, default_ bool) *BoolFlag

```go
func NewBoolFlag(longName, shortName, desc string, default_ bool) *BoolFlag
```

NewBoolFlag 创建布尔标志

**参数:**
  - longName: 长选项名, 如 "verbose"
  - shortName: 短选项名, 如 "v"
  - desc: 标志描述
  - default_: 默认值

**返回值:**
  - *BoolFlag: 布尔标志实例

#### func (f *BoolFlag) Set(value string) error

```go
func (f *BoolFlag) Set(value string) error
```

Set 设置布尔标志的值

**参数:**
  - value: 要设置的字符串值

**返回值:**
  - error: 如果解析失败返回错误

**注意事项: **
  - 支持标准库的布尔值解析格式
  - 空字符串会被解析为true (这是Go flag包的标准行为) 
  - 支持的值: 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False

---

### type DurationFlag struct

```go
type DurationFlag struct {
    *BaseFlag[time.Duration]
}
```

DurationFlag 持续时间标志

DurationFlag 用于处理时间间隔类型的命令行参数。 支持Go标准库time.ParseDuration所支持的所有格式, 如 "300ms", "-1.5h", "2h45m" 等。

**支持的格式: **
  - "ns": 纳秒
  - "us" (或 "µs"): 微秒
  - "ms": 毫秒
  - "s": 秒
  - "m": 分钟
  - "h": 小时

**注意事项: **
  - 支持负数表示负时间间隔
  - 支持小数表示部分时间单位
  - 可以组合多个单位, 如 "1h30m"

#### func NewDurationFlag(longName, shortName, desc string, default_ time.Duration) *DurationFlag

```go
func NewDurationFlag(longName, shortName, desc string, default_ time.Duration) *DurationFlag
```

NewDurationFlag 创建新的持续时间标志

**参数:**
  - longName: 长选项名, 如 "timeout"
  - shortName: 短选项名, 如 "t"
  - desc: 标志描述
  - default_: 默认值

**返回值:**
  - *DurationFlag: 持续时间标志实例

#### func (f *DurationFlag) Set(value string) error

```go
func (f *DurationFlag) Set(value string) error
```

Set 设置持续时间标志的值

**参数:**
  - value: 要设置的时间间隔字符串

**返回值:**
  - error: 如果解析失败返回错误

**注意事项: **
  - 使用 time.ParseDuration 解析字符串
  - 支持所有Go标准库支持的时间格式
  - 如果值无法解析为时间间隔, 返回解析错误

---

### type EnumFlag struct

```go
type EnumFlag struct {
    *BaseFlag[string]

    // Has unexported fields.
}
```

EnumFlag 枚举标志

EnumFlag 用于处理枚举类型的命令行参数, 限制输入值必须在预定义的允许值列表中。 使用映射表(map)实现O(1)时间复杂度的值查找, 提高性能。

**特性: **
  - 使用映射表进行快速值验证
  - 不允许空字符串作为枚举值
  - 默认值必须在允许值列表中
  - 不允许设置空值

#### func NewEnumFlag(longName, shortName, desc, default_ string, allowedValues []string) *EnumFlag

```go
func NewEnumFlag(longName, shortName, desc, default_ string, allowedValues []string) *EnumFlag
```

NewEnumFlag 创建枚举标志

**参数:**
  - longName: 长选项名, 如 "mode"
  - shortName: 短选项名, 如 "m"
  - desc: 标志描述
  - default_: 默认值, 必须在允许值列表中
  - allowedValues: 允许的枚举值列表, 不能为空且不能包含空字符串

**返回值:**
  - *EnumFlag: 枚举标志实例

**注意事项: **
  - 允许值列表不能为空
  - 允许值列表中不能包含空字符串
  - 默认值必须在允许值列表中
  - 如果验证失败, 会panic

#### func (f *EnumFlag) EnumValues() []string

```go
func (f *EnumFlag) EnumValues() []string
```

EnumValues 获取枚举类型的可选值

**返回值:**
  - []string: 枚举类型的可选值列表

**功能说明: **
  - 实现 Flag 接口的 EnumValues 方法
  - 返回所有允许的枚举值
  - 此方法是线程安全的

#### func (f *EnumFlag) GetAllowedValues() []string

```go
func (f *EnumFlag) GetAllowedValues() []string
```

GetAllowedValues 获取允许的枚举值

**返回值:**
  - []string: 允许的枚举值列表

**注意事项: **
  - 返回的切片顺序可能不一致, 因为基于map的key生成
  - 此方法是线程安全的

#### func (f *EnumFlag) IsAllowed(value string) bool

```go
func (f *EnumFlag) IsAllowed(value string) bool
```

IsAllowed 检查值是否在允许的枚举值中

**参数:**
  - value: 要检查的值

**返回值:**
  - bool: 如果值在允许列表中返回true, 否则返回false

**注意事项: **
  - 此方法是线程安全的
  - 使用映射表进行O(1)时间复杂度的查找

#### func (f *EnumFlag) Set(value string) error

```go
func (f *EnumFlag) Set(value string) error
```

Set 设置枚举标志的值

**参数:**
  - value: 要设置的字符串值

**返回值:**
  - error: 如果值不在允许列表中或为空, 返回错误

**注意事项: **
  - 不允许设置空值
  - 使用映射表进行O(1)时间复杂度的值验证
  - 错误消息会列出所有允许的值

---

### type Float64Flag struct

```go
type Float64Flag struct {
    *BaseFlag[float64]
}
```

Float64Flag 64位浮点数标志

Float64Flag 用于处理64位浮点数类型的命令行参数。 支持整数、小数和科学计数法表示的数值。

**注意事项: **
  - 支持正数和负数
  - 支持十进制格式和科学计数法
  - 支持特殊值: NaN、+Inf、-Inf
  - 精度遵循IEEE 754双精度浮点数标准

#### func NewFloat64Flag(longName, shortName, desc string, default_ float64) *Float64Flag

```go
func NewFloat64Flag(longName, shortName, desc string, default_ float64) *Float64Flag
```

NewFloat64Flag 创建新的64位浮点数标志

**参数:**
  - longName: 长选项名, 如 "ratio"
  - shortName: 短选项名, 如 "r"
  - desc: 标志描述
  - default_: 默认值

**返回值:**
  - *Float64Flag: 64位浮点数标志实例

#### func (f *Float64Flag) Set(value string) error

```go
func (f *Float64Flag) Set(value string) error
```

Set 设置64位浮点数标志的值

**参数:**
  - value: 要设置的64位浮点数字符串

**返回值:**
  - error: 如果解析失败返回错误

**注意事项: **
  - 使用 strconv.ParseFloat 解析字符串
  - 固定使用64位精度
  - 支持十进制格式和科学计数法
  - 如果值格式无效, 返回解析错误

---

### type Int64Flag struct

```go
type Int64Flag struct {
    *BaseFlag[int64]
}
```

Int64Flag 64位整数标志

Int64Flag 用于处理64位整数类型的命令行参数。 在所有平台上都使用固定的64位整数, 提供一致的行为。

**注意事项: **
  - 支持正数和负数
  - 支持十进制格式
  - 范围: -9,223,372,036,854,775,808 到 9,223,372,036,854,775,807

#### func NewInt64Flag(longName, shortName, desc string, default_ int64) *Int64Flag

```go
func NewInt64Flag(longName, shortName, desc string, default_ int64) *Int64Flag
```

NewInt64Flag 创建64位整数标志

**参数:**
  - longName: 长选项名, 如 "timestamp"
  - shortName: 短选项名, 如 "ts"
  - desc: 标志描述
  - default_: 默认值

**返回值:**
  - *Int64Flag: 64位整数标志实例

#### func (f *Int64Flag) Set(value string) error

```go
func (f *Int64Flag) Set(value string) error
```

Set 设置64位整数标志的值

**参数:**
  - value: 要设置的64位整数字符串

**返回值:**
  - error: 如果解析失败返回错误

**注意事项: **
  - 使用 strconv.ParseInt 解析字符串
  - 固定使用64位精度
  - 如果值超出64位整数范围, 返回解析错误

---

### type Int64SliceFlag struct

```go
type Int64SliceFlag struct {
    *BaseFlag[[]int64]
}
```

Int64SliceFlag 64位整数切片标志

#### func NewInt64SliceFlag(longName, shortName, desc string, default_ []int64) *Int64SliceFlag

```go
func NewInt64SliceFlag(longName, shortName, desc string, default_ []int64) *Int64SliceFlag
```

NewInt64SliceFlag 创建新的64位整数切片标志

#### func (f *Int64SliceFlag) IsEmpty() bool

```go
func (f *Int64SliceFlag) IsEmpty() bool
```

IsEmpty 检查切片是否为空

#### func (f *Int64SliceFlag) Length() int

```go
func (f *Int64SliceFlag) Length() int
```

Length 获取切片长度

#### func (f *Int64SliceFlag) Set(value string) error

```go
func (f *Int64SliceFlag) Set(value string) error
```

Set 设置64位整数切片标志的值

---

### type IntFlag struct

```go
type IntFlag struct {
    *BaseFlag[int]
}
```

IntFlag 整数标志

IntFlag 用于处理整数类型的命令行参数。 使用平台相关的int类型, 在32位系统上为32位整数, 在64位系统上为64位整数。

**注意事项: **
  - 支持正数和负数
  - 支持十进制格式
  - 超出平台int范围会返回错误

#### func NewIntFlag(longName, shortName, desc string, default_ int) *IntFlag

```go
func NewIntFlag(longName, shortName, desc string, default_ int) *IntFlag
```

NewIntFlag 创建整数标志

**参数:**
  - longName: 长选项名, 如 "count"
  - shortName: 短选项名, 如 "c"
  - desc: 标志描述
  - default_: 默认值

**返回值:**
  - *IntFlag: 整数标志实例

#### func (f *IntFlag) Set(value string) error

```go
func (f *IntFlag) Set(value string) error
```

Set 设置整数标志的值

**参数:**
  - value: 要设置的整数字符串

**返回值:**
  - error: 如果解析失败返回错误

**注意事项: **
  - 使用 strconv.ParseInt 解析字符串
  - 使用平台相关的位数(IntSize)
  - 如果值超出平台int范围, 返回解析错误

---

### type IntSliceFlag struct

```go
type IntSliceFlag struct {
    *BaseFlag[[]int]
}
```

IntSliceFlag 整数切片标志

#### func NewIntSliceFlag(longName, shortName, desc string, default_ []int) *IntSliceFlag

```go
func NewIntSliceFlag(longName, shortName, desc string, default_ []int) *IntSliceFlag
```

NewIntSliceFlag 创建新的整数切片标志

#### func (f *IntSliceFlag) IsEmpty() bool

```go
func (f *IntSliceFlag) IsEmpty() bool
```

IsEmpty 检查切片是否为空

#### func (f *IntSliceFlag) Length() int

```go
func (f *IntSliceFlag) Length() int
```

Length 获取切片长度

#### func (f *IntSliceFlag) Set(value string) error

```go
func (f *IntSliceFlag) Set(value string) error
```

Set 设置整数切片标志的值

---

### type MapFlag struct

```go
type MapFlag struct {
    *BaseFlag[map[string]string]
}
```

MapFlag 用于处理键值对映射类型的命令行参数。 支持的格式: key1=value1,key2=value2

**空值处理:**
  - 空字符串 "" 表示创建空映射
  - ",,," 中的空对会被跳过
  - 使用 Clear 方法可以清空映射

#### func NewMapFlag(longName, shortName, desc string, default_ map[string]string) *MapFlag

```go
func NewMapFlag(longName, shortName, desc string, default_ map[string]string) *MapFlag
```

NewMapFlag 创建新的映射标志

**参数:**
  - longName: 长选项名
  - shortName: 短选项名
  - desc: 标志描述
  - default_: 默认值, 如果为nil则创建空映射

**返回值:**
  - *MapFlag: 映射标志实例

#### func (f *MapFlag) Clear()

```go
func (f *MapFlag) Clear()
```

Clear 清空映射 将映射设置为空映射, 并标记为已设置

#### func (f *MapFlag) GetKey(key string) (string, bool)

```go
func (f *MapFlag) GetKey(key string) (string, bool)
```

GetKey 获取映射中指定键的值

#### func (f *MapFlag) HasKey(key string) bool

```go
func (f *MapFlag) HasKey(key string) bool
```

HasKey 检查映射中是否包含指定键

#### func (f *MapFlag) IsEmpty() bool

```go
func (f *MapFlag) IsEmpty() bool
```

IsEmpty 检查映射是否为空

#### func (f *MapFlag) Keys() []string

```go
func (f *MapFlag) Keys() []string
```

Keys 获取映射的所有键

#### func (f *MapFlag) Length() int

```go
func (f *MapFlag) Length() int
```

Length 获取映射长度

#### func (f *MapFlag) Set(value string) error

```go
func (f *MapFlag) Set(value string) error
```

Set 设置映射标志的值

支持格式: key1=value1,key2=value2

**空值处理:**
  - 空字符串 "" 表示创建空映射
  - ",,," 中的空对会被跳过
  - 键不能为空, 否则返回错误

**参数:**
  - value: 映射字符串

**返回值:**
  - error: 如果解析失败返回错误

---

### type SizeFlag struct

```go
type SizeFlag struct {
    *BaseFlag[int64]
}
```

SizeFlag 大小标志 (支持KB、MB、GB等单位)

SizeFlag 用于处理大小类型的命令行参数, 支持多种大小单位。 可以解析带有单位的大小值, 并将其转换为字节数。

**支持的单位: **
  - B/b: 字节
  - KB/kb/K/k: 千字节 (1024字节)
  - MB/mb/M/m: 兆字节 (1024^2字节)
  - GB/gb/G/g: 吉字节 (1024^3字节)
  - TB/tb/T/t: 太字节 (1024^4字节)
  - PB/pb/P/p: 拍字节 (1024^5字节)
  - KiB/kib: 二进制千字节 (1024字节)
  - MiB/mib: 二进制兆字节 (1024^2字节)
  - GiB/gib: 二进制吉字节 (1024^3字节)
  - TiB/tib: 二进制太字节 (1024^4字节)
  - PiB/pib: 二进制拍字节 (1024^5字节)

**注意事项: **
  - 支持小数, 如 "1.5MB"
  - 不支持负数
  - 默认单位为字节(B)
  - 大小写不敏感

#### func NewSizeFlag(longName, shortName, desc string, default_ int64) *SizeFlag

```go
func NewSizeFlag(longName, shortName, desc string, default_ int64) *SizeFlag
```

NewSizeFlag 创建新的大小标志

**参数:**
  - longName: 长选项名, 如 "max-size"
  - shortName: 短选项名, 如 "s"
  - desc: 标志描述
  - default_: 默认值(以字节为单位)

**返回值:**
  - *SizeFlag: 大小标志实例

#### func (f *SizeFlag) Set(value string) error

```go
func (f *SizeFlag) Set(value string) error
```

Set 设置大小标志的值

**参数:**
  - value: 要设置的大小字符串, 可包含单位

**返回值:**
  - error: 如果解析失败返回错误

**注意事项: **
  - 支持多种大小单位, 大小写不敏感
  - 支持小数值, 如 "1.5MB"
  - 不支持负数
  - 如果未指定单位, 默认为字节(B)
  - 如果值超出int64范围, 返回错误

#### func (f *SizeFlag) String() string

```go
func (f *SizeFlag) String() string
```

String 返回格式化的大小字符串

**返回值:**
  - string: 格式化的大小字符串, 如 "1024B"、"2.5MB" 等

**注意事项: **
  - 此方法是线程安全的
  - 实现了fmt.Stringer接口

---

### type StringFlag struct

```go
type StringFlag struct {
    *BaseFlag[string]
}
```

StringFlag 字符串标志

StringFlag 用于处理字符串类型的命令行参数。 它接受任何字符串值, 包括空字符串。

#### func NewStringFlag(longName, shortName, desc string, default_ string) *StringFlag

```go
func NewStringFlag(longName, shortName, desc string, default_ string) *StringFlag
```

NewStringFlag 创建字符串标志

**参数:**
  - longName: 长选项名, 如 "output"
  - shortName: 短选项名, 如 "o"
  - desc: 标志描述
  - default_: 默认值

**返回值:**
  - *StringFlag: 字符串标志实例

#### func (f *StringFlag) Set(value string) error

```go
func (f *StringFlag) Set(value string) error
```

Set 设置字符串标志的值

**参数:**
  - value: 要设置的字符串值

**返回值:**
  - error: 总是返回nil

**注意事项: **
  - 字符串标志接受任何字符串值, 包括空字符串
  - 不会对输入值进行任何验证或转换

---

### type StringSliceFlag struct

```go
type StringSliceFlag struct {
    *BaseFlag[[]string]
}
```

StringSliceFlag 字符串切片标志

#### func NewStringSliceFlag(longName, shortName, desc string, default_ []string) *StringSliceFlag

```go
func NewStringSliceFlag(longName, shortName, desc string, default_ []string) *StringSliceFlag
```

NewStringSliceFlag 创建新的字符串切片标志

#### func (f *StringSliceFlag) IsEmpty() bool

```go
func (f *StringSliceFlag) IsEmpty() bool
```

IsEmpty 检查切片是否为空

#### func (f *StringSliceFlag) Length() int

```go
func (f *StringSliceFlag) Length() int
```

Length 获取切片长度

#### func (f *StringSliceFlag) Set(value string) error

```go
func (f *StringSliceFlag) Set(value string) error
```

Set 设置字符串切片标志的值

---

### type TimeFlag struct

```go
type TimeFlag struct {
    *BaseFlag[time.Time]

    // Has unexported fields.
}
```

TimeFlag 时间标志

TimeFlag 用于处理时间类型的命令行参数。 支持自动检测多种常见时间格式, 也支持指定特定格式进行解析。

**特性: **
  - 自动检测常见时间格式
  - 支持自定义格式解析
  - 记录当前使用的格式
  - 线程安全的格式存储

**常见支持格式: **
  - RFC3339: "2006-01-02T15:04:05Z07:00"
  - RFC1123: "Mon, 02 Jan 2006 15:04:05 MST"
  - 日期格式: "2006-01-02", "2006/01/02"
  - 时间格式: "15:04:05", "15:04"
  - 其他常见格式

#### func NewTimeFlag(longName, shortName, desc string, default_ time.Time) *TimeFlag

```go
func NewTimeFlag(longName, shortName, desc string, default_ time.Time) *TimeFlag
```

NewTimeFlag 创建新的时间标志

**参数:**
  - longName: 长选项名, 如 "start-time"
  - shortName: 短选项名, 如 "s"
  - desc: 标志描述
  - default_: 默认值

**返回值:**
  - *TimeFlag: 时间标志实例

#### func (f *TimeFlag) FormatTime(t time.Time) string

```go
func (f *TimeFlag) FormatTime(t time.Time) string
```

FormatTime 使用当前解析格式格式化时间

**参数:**
  - t: 要格式化的时间值

**返回值:**
  - string: 格式化后的时间字符串

**注意事项: **
  - 如果当前格式为空, 使用RFC3339格式
  - 此方法是线程安全的

#### func (f *TimeFlag) GetFormat() string

```go
func (f *TimeFlag) GetFormat() string
```

GetFormat 获取当前使用的时间格式

**返回值:**
  - string: 当前使用的时间格式, 如果未设置则返回空字符串

**注意事项: **
  - 此方法是线程安全的
  - 返回的格式可用于格式化其他时间值

#### func (f *TimeFlag) Set(value string) error

```go
func (f *TimeFlag) Set(value string) error
```

Set 使用常见格式自动解析时间

**参数:**
  - value: 要设置的时间字符串

**返回值:**
  - error: 如果解析失败返回错误

**注意事项: **
  - 使用 types.ParseTimeWithCommonFormats 自动检测格式
  - 成功解析后会记录使用的格式
  - 如果无法匹配任何格式, 返回错误

#### func (f *TimeFlag) SetWithFormat(value, format string) error

```go
func (f *TimeFlag) SetWithFormat(value, format string) error
```

SetWithFormat 使用指定格式解析时间

**参数:**
  - value: 要设置的时间字符串
  - format: 时间格式字符串, 遵循Go的time.Format布局

**返回值:**
  - error: 如果解析失败返回错误

**注意事项: **
  - 使用 time.Parse 按指定格式解析
  - 成功解析后会更新当前使用的格式
  - 格式字符串必须遵循Go的time.Format布局规则

#### func (f *TimeFlag) String() string

```go
func (f *TimeFlag) String() string
```

String 返回格式化的时间字符串

**返回值:**
  - string: 格式化的时间字符串

**注意事项: **
  - 如果当前格式为空, 使用RFC3339格式
  - 此方法是线程安全的
  - 实现了fmt.Stringer接口

---

### type Uint16Flag struct

```go
type Uint16Flag struct {
    *BaseFlag[uint16]
}
```

Uint16Flag 16位无符号整数标志

Uint16Flag 用于处理16位无符号整数类型的命令行参数。 适用于处理端口号、短范围计数器等场景。

**注意事项: **
  - 只支持非负数
  - 支持十进制格式
  - 范围: 0 到 65,535

#### func NewUint16Flag(longName, shortName, desc string, default_ uint16) *Uint16Flag

```go
func NewUint16Flag(longName, shortName, desc string, default_ uint16) *Uint16Flag
```

NewUint16Flag 创建新的16位无符号整数标志

**参数:**
  - longName: 长选项名, 如 "port"
  - shortName: 短选项名, 如 "p"
  - desc: 标志描述
  - default_: 默认值

**返回值:**
  - *Uint16Flag: 16位无符号整数标志实例

#### func (f *Uint16Flag) Set(value string) error

```go
func (f *Uint16Flag) Set(value string) error
```

Set 设置16位无符号整数标志的值

**参数:**
  - value: 要设置的16位无符号整数字符串

**返回值:**
  - error: 如果解析失败返回错误

**注意事项: **
  - 使用 strconv.ParseUint 解析字符串
  - 固定使用16位精度
  - 如果值超出0-65,535范围或为负数, 返回解析错误

---

### type Uint32Flag struct

```go
type Uint32Flag struct {
    *BaseFlag[uint32]
}
```

Uint32Flag 32位无符号整数标志

Uint32Flag 用于处理32位无符号整数类型的命令行参数。 适用于处理IP地址、大范围计数器等场景。

**注意事项: **
  - 只支持非负数
  - 支持十进制格式
  - 范围: 0 到 4,294,967,295

#### func NewUint32Flag(longName, shortName, desc string, default_ uint32) *Uint32Flag

```go
func NewUint32Flag(longName, shortName, desc string, default_ uint32) *Uint32Flag
```

NewUint32Flag 创建新的32位无符号整数标志

**参数:**
  - longName: 长选项名, 如 "ip"
  - shortName: 短选项名, 如 "i"
  - desc: 标志描述
  - default_: 默认值

**返回值:**
  - *Uint32Flag: 32位无符号整数标志实例

#### func (f *Uint32Flag) Set(value string) error

```go
func (f *Uint32Flag) Set(value string) error
```

Set 设置32位无符号整数标志的值

**参数:**
  - value: 要设置的32位无符号整数字符串

**返回值:**
  - error: 如果解析失败返回错误

**注意事项: **
  - 使用 strconv.ParseUint 解析字符串
  - 固定使用32位精度
  - 如果值超出0-4,294,967,295范围或为负数, 返回解析错误

---

### type Uint64Flag struct

```go
type Uint64Flag struct {
    *BaseFlag[uint64]
}
```

Uint64Flag 64位无符号整数标志

Uint64Flag 用于处理64位无符号整数类型的命令行参数。 在所有平台上都使用固定的64位无符号整数, 提供一致的行为。

**注意事项: **
  - 只支持非负数
  - 支持十进制格式
  - 范围: 0 到 18,446,744,073,709,551,615

#### func NewUint64Flag(longName, shortName, desc string, default_ uint64) *Uint64Flag

```go
func NewUint64Flag(longName, shortName, desc string, default_ uint64) *Uint64Flag
```

NewUint64Flag 创建新的64位无符号整数标志

**参数:**
  - longName: 长选项名, 如 "id"
  - shortName: 短选项名, 如 "i"
  - desc: 标志描述
  - default_: 默认值

**返回值:**
  - *Uint64Flag: 64位无符号整数标志实例

#### func (f *Uint64Flag) Set(value string) error

```go
func (f *Uint64Flag) Set(value string) error
```

Set 设置64位无符号整数标志的值

**参数:**
  - value: 要设置的64位无符号整数字符串

**返回值:**
  - error: 如果解析失败返回错误

**注意事项: **
  - 使用 strconv.ParseUint 解析字符串
  - 固定使用64位精度
  - 如果值超出0-18,446,744,073,709,551,615范围或为负数, 返回解析错误

---

### type Uint8Flag struct

```go
type Uint8Flag struct {
    *BaseFlag[uint8]
}
```

Uint8Flag 8位无符号整数标志

Uint8Flag 用于处理8位无符号整数类型的命令行参数。 适用于处理字节值、小范围计数器等场景。

**注意事项: **
  - 只支持非负数
  - 支持十进制格式
  - 范围: 0 到 255

#### func NewUint8Flag(longName, shortName, desc string, default_ uint8) *Uint8Flag

```go
func NewUint8Flag(longName, shortName, desc string, default_ uint8) *Uint8Flag
```

NewUint8Flag 创建新的8位无符号整数标志

**参数:**
  - longName: 长选项名, 如 "byte"
  - shortName: 短选项名, 如 "b"
  - desc: 标志描述
  - default_: 默认值

**返回值:**
  - *Uint8Flag: 8位无符号整数标志实例

#### func (f *Uint8Flag) Set(value string) error

```go
func (f *Uint8Flag) Set(value string) error
```

Set 设置8位无符号整数标志的值

**参数:**
  - value: 要设置的8位无符号整数字符串

**返回值:**
  - error: 如果解析失败返回错误

**注意事项: **
  - 使用 strconv.ParseUint 解析字符串
  - 固定使用8位精度
  - 如果值超出0-255范围或为负数, 返回解析错误

---

### type UintFlag struct

```go
type UintFlag struct {
    *BaseFlag[uint]
}
```

UintFlag 无符号整数标志

UintFlag 用于处理无符号整数类型的命令行参数。 使用平台相关的uint类型, 在32位系统上为32位无符号整数, 在64位系统上为64位无符号整数。

**注意事项: **
  - 只支持非负数
  - 支持十进制格式
  - 超出平台uint范围会返回错误

#### func NewUintFlag(longName, shortName, desc string, default_ uint) *UintFlag

```go
func NewUintFlag(longName, shortName, desc string, default_ uint) *UintFlag
```

NewUintFlag 创建新的无符号整数标志

**参数:**
  - longName: 长选项名, 如 "port"
  - shortName: 短选项名, 如 "p"
  - desc: 标志描述
  - default_: 默认值

**返回值:**
  - *UintFlag: 无符号整数标志实例

#### func (f *UintFlag) Set(value string) error

```go
func (f *UintFlag) Set(value string) error
```

Set 设置无符号整数标志的值

**参数:**
  - value: 要设置的无符号整数字符串

**返回值:**
  - error: 如果解析失败返回错误

**注意事项: **
  - 使用 strconv.ParseUint 解析字符串
  - 使用平台相关的位数(UintSize)
  - 如果值超出平台uint范围或为负数, 返回解析错误