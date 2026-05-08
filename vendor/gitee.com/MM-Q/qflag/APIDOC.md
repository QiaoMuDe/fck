# qflag

```go
package qflag // import "gitee.com/MM-Q/qflag"
```

---

## 目录

- [变量](#变量)
- [函数](#函数)
- [类型](#类型)

---

## 变量

### 错误处理策略常量

```go
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
```

### GenAndPrintCompletion

```go
var GenAndPrintCompletion = completion.GenAndPrint
```

生成并打印补全脚本

**参数:**

- `cmd`: 要生成补全脚本的命令
- `shellType`: Shell类型 (bash, pwsh, powershell)

**功能说明:**

- 生成自动补全脚本
- 直接输出到标准输出
- 便于在命令行中直接使用

### GenerateCompletion

```go
var GenerateCompletion = completion.Generate
```

生成补全脚本

**参数:**

- `cmd`: 要生成补全脚本的命令
- `shellType`: Shell类型 (bash, pwsh, powershell)

**返回值:**

- `string`: 生成的补全脚本
- `error`: 生成失败时返回错误

**功能说明:**

- 为指定命令生成自动补全脚本
- 支持多种Shell类型
- 包含完整的命令树和标志信息

### NewCmd

```go
var NewCmd = cmd.NewCmd
```

创建新的命令实例

**参数:**

- `longName`: 命令的长名称
- `shortName`: 命令的短名称
- `errorHandling`: 错误处理策略

**返回值:**

- `*Cmd`: 初始化完成的命令实例

**功能说明:**

- 创建命令并初始化基本字段
- 创建标志和子命令注册器
- 设置默认解析器
- 初始化配置选项

### NewCmdOpts

```go
var NewCmdOpts = cmd.NewCmdOpts
```

创建新的命令选项

**返回值:**

- `*CmdOpts`: 初始化的命令选项

**功能说明:**

- 创建基本命令选项
- 初始化所有字段为零值
- 初始化 map 和 slice 避免空指针

---

## 函数

### AddSubCmdFrom

```go
func AddSubCmdFrom(cmds []Command) error
```

从切片添加子命令

**参数:**

- `cmds`: 要添加的子命令实例切片

**返回值:**

- `error`: 添加子命令过程中遇到的错误, 如果没有错误则返回 nil

### AddSubCmds

```go
func AddSubCmds(cmds ...Command) error
```

添加子命令到全局根命令

**参数:**

- `cmd`: 要添加的子命令实例

**返回值:**

- `error`: 添加子命令过程中遇到的错误, 如果没有错误则返回 nil

### ApplyOpts

```go
func ApplyOpts(opts *CmdOpts) error
```

应用选项到全局根命令

**参数:**

- `opts`: 要应用的选项结构体实例

**返回值:**

- `error`: 应用选项过程中遇到的错误, 如果没有错误则返回 nil

**功能说明:**

- 将选项结构体的所有属性应用到全局根命令实例
- 支持部分配置（未设置的属性不会被修改）
- 使用写锁保护并发安全

### Parse

```go
func Parse() error
```

解析命令行参数

**返回值:**

- `error`: 解析失败时返回错误

**功能说明:**

- 使用全局根命令解析命令行参数
- 可以重复调用，会覆盖之前的解析结果
- 递归解析所有子命令

**注意事项:**

- 如果需要确保只解析一次，请使用 `ParseOnce`

### ParseAndRoute

```go
func ParseAndRoute() error
```

解析并路由执行命令

**返回值:**

- `error`: 解析或执行失败时返回错误

**功能说明:**

- 使用全局根命令解析命令行参数
- 可以重复调用，会覆盖之前的解析结果
- 完整的解析和执行流程

**注意事项:**

- 如果需要确保只解析一次，请使用 `ParseAndRouteOnce`

### ParseAndRouteOnce

```go
func ParseAndRouteOnce() error
```

解析并路由执行命令（只解析一次）

**返回值:**

- `error`: 解析或执行失败时返回错误

**功能说明:**

- 使用全局根命令解析命令行参数
- 使用`ParseAndRouteOnce`确保只解析一次
- 重复执行无错误、仅首次执行解析
- 完整的解析和执行流程

**注意事项:**

- 建议在普通场景使用此方法，避免误用
- 如果需要重复解析，请使用 `ParseAndRoute`

### ParseOnce

```go
func ParseOnce() error
```

解析命令行参数（只解析一次）

**返回值:**

- `error`: 解析失败时返回错误

**功能说明:**

- 使用全局根命令解析命令行参数
- 使用`ParseOnce`确保只解析一次
- 重复执行无错误、仅首次执行解析
- 递归解析所有子命令

**注意事项:**

- 建议在普通场景使用此方法，避免误用
- 如果需要重复解析，请使用 `Parse`

### ParseOnly

```go
func ParseOnly() error
```

仅解析当前命令, 不递归解析子命令

**返回值:**

- `error`: 解析失败时返回错误

**功能说明:**

- 使用全局根命令解析命令行参数
- 可以重复调用，会覆盖之前的解析结果
- 不处理子命令解析

**注意事项:**

- 如果需要确保只解析一次，请使用 `ParseOnlyOnce`

### ParseOnlyOnce

```go
func ParseOnlyOnce() error
```

仅解析当前命令, 不递归解析子命令（只解析一次）

**返回值:**

- `error`: 解析失败时返回错误

**功能说明:**

- 使用全局根命令解析命令行参数
- 使用`ParseOnlyOnce`确保只解析一次
- 重复执行无错误、仅首次执行解析
- 不处理子命令解析

**注意事项:**

- 建议在普通场景使用此方法，避免误用
- 如果需要重复解析，请使用 `ParseOnly`

---

## 类型

### BoolFlag

```go
type BoolFlag = flag.BoolFlag
```

**BoolFlag 布尔标志**

`BoolFlag` 用于处理布尔类型的命令行参数。它接受多种布尔值表示形式, 包括 `"true"`, `"false"`, `"1"`, `"0"`, `"t"`, `"f"`, `"TRUE"`, `"FALSE"` 等。

### Cmd

```go
type Cmd = cmd.Cmd
```

**Cmd 命令结构体类型**

`Cmd` 是一个命令结构体, 实现了 `types.Command` 接口，提供了完整的命令行命令实现, 支持标志管理、子命令、参数解析和执行等功能。使用读写锁保证并发安全。

### Root

```go
var Root *Cmd
```

**Root 全局根命令实例**

提供对全局标志和子命令的访问。用户可以通过 `qflag.Root.String()` 这样的方式直接创建全局标志。这是访问命令行功能的主要入口点, 推荐优先使用。

### CmdConfig

```go
type CmdConfig = types.CmdConfig
```

**CmdConfig**

包含了命令的各种配置选项, 用于自定义命令的行为和外观。这些配置会影响命令的帮助信息显示、环境变量处理、错误提示等。

### CmdOpts

```go
type CmdOpts = cmd.CmdOpts
```

**CmdOpts 命令选项结构体**

`CmdOpts` 提供了配置现有命令的方式，包含命令的所有可配置属性。

### CmdRegistry

```go
type CmdRegistry = types.CmdRegistry
```

**CmdRegistry 命令注册表接口**

`CmdRegistry` 定义了命令注册和管理的标准接口, 提供了命令的完整生命周期管理功能。

**核心功能:**

- 命令的注册和注销
- 基于名称的查找和检索
- 批量操作和遍历支持
- 存在性检查和计数

**设计特点:**

- 支持长名称和短名称查找
- 提供统一的错误处理
- 支持别名管理 (通过具体实现)
- 线程安全由具体实现保证

### Command

```go
type Command = types.Command
```

**Command**

定义了命令行工具中命令的基本接口, 包括标志管理、参数解析、子命令管理等功能。实现此接口的类型可以作为命令行工具的命令使用。

### DurationFlag

```go
type DurationFlag = flag.DurationFlag
```

**DurationFlag 持续时间标志**

`DurationFlag` 用于处理时间间隔类型的命令行参数。支持Go标准库`time.ParseDuration`所支持的所有格式, 如 `"300ms"`, `"-1.5h"`, `"2h45m"` 等。

**支持的格式:**

- `"ns"`: 纳秒
- `"us"` (或 `"µs"`): 微秒
- `"ms"`: 毫秒
- `"s"`: 秒
- `"m"`: 分钟
- `"h"`: 小时

**注意事项:**

- 支持负数表示负时间间隔
- 支持小数表示部分时间单位
- 可以组合多个单位, 如 `"1h30m"`

### EnumFlag

```go
type EnumFlag = flag.EnumFlag
```

**EnumFlag 枚举标志**

`EnumFlag` 用于处理枚举类型的命令行参数, 限制输入值必须在预定义的允许值列表中。使用映射表(map)实现O(1)时间复杂度的值查找, 提高性能。

**特性:**

- 使用映射表进行快速值验证
- 不允许空字符串作为枚举值
- 默认值必须在允许值列表中
- 不允许设置空值

### ErrorHandling

```go
type ErrorHandling = types.ErrorHandling
```

**ErrorHandling 错误处理方式枚举**

`ErrorHandling` 定义了解析错误时的处理策略, 直接使用标准库 flag包的错误处理方式, 保持兼容性。

**可选值:**

- `ContinueOnError`: 解析错误时继续解析并返回错误
- `ExitOnError`: 解析错误时退出程序
- `PanicOnError`: 解析错误时触发panic

### Flag

```go
type Flag = types.Flag
```

**Flag 接口**

定义了标志的核心行为。`Flag` 是所有标志类型必须实现的基础接口, 定义了标志的基本操作和属性。所有具体标志类型都应实现此接口。

**设计原则:**

- 提供统一的标志操作接口
- 支持多种数据类型
- 支持验证和环境变量绑定
- 提供完整的生命周期管理

### FlagRegistry

```go
type FlagRegistry = types.FlagRegistry
```

**FlagRegistry 标志注册表接口**

`FlagRegistry` 定义了标志注册和管理的标准接口, 提供了标志的完整生命周期管理功能。

**核心功能:**

- 标志的注册和注销
- 基于名称的查找和检索
- 批量操作和遍历支持
- 存在性检查和计数

**设计特点:**

- 支持长名称和短名称查找
- 提供统一的错误处理
- 支持别名管理 (通过具体实现)
- 线程安全由具体实现保证

### FlagType

```go
type FlagType = types.FlagType
```

**FlagType 标志类型枚举**

`FlagType` 定义了所有支持的标志类型, 用于类型识别和特定处理逻辑的实现。

**设计原则:**

- 每种类型对应一种数据格式
- 支持基础类型和复合类型
- 便于类型检查和转换

#### 标志类型常量

```go
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
```

### Float64Flag

```go
type Float64Flag = flag.Float64Flag
```

**Float64Flag 64位浮点数标志**

`Float64Flag` 用于处理64位浮点数类型的命令行参数。支持整数、小数和科学计数法表示的数值。

**注意事项:**

- 支持正数和负数
- 支持十进制格式和科学计数法
- 支持特殊值: NaN、+Inf、-Inf
- 精度遵循IEEE 754双精度浮点数标准

### Int64Flag

```go
type Int64Flag = flag.Int64Flag
```

**Int64Flag 64位整数标志**

`Int64Flag` 用于处理64位整数类型的命令行参数。在所有平台上都使用固定的64位整数, 提供一致的行为。

**注意事项:**

- 支持正数和负数
- 支持十进制格式
- 范围: -9,223,372,036,854,775,808 到 9,223,372,036,854,775,807

### Int64SliceFlag

```go
type Int64SliceFlag = flag.Int64SliceFlag
```

**Int64SliceFlag 64位整数切片标志**

### IntFlag

```go
type IntFlag = flag.IntFlag
```

**IntFlag 整数标志**

`IntFlag` 用于处理整数类型的命令行参数。使用平台相关的int类型, 在32位系统上为32位整数, 在64位系统上为64位整数。

**注意事项:**

- 支持正数和负数
- 支持十进制格式
- 超出平台int范围会返回错误

### IntSliceFlag

```go
type IntSliceFlag = flag.IntSliceFlag
```

**IntSliceFlag 整数切片标志**

### MapFlag

```go
type MapFlag = flag.MapFlag
```

**MapFlag**

用于处理键值对映射类型的命令行参数。支持的格式: `key1=value1,key2=value2`

**空值处理:**

- 空字符串 `""` 表示创建空映射
- `",,,"` 中的空对会被跳过
- 使用 `SetKV` 方法设置键值对时, 键不能为空
- 使用 `Clear` 方法可以清空映射

### MutexGroup

```go
type MutexGroup = types.MutexGroup
```

**MutexGroup**

定义了一组互斥的标志, 其中最多只能有一个被设置。当用户设置了互斥组中的多个标志时, 解析器会返回错误。

### Parser

```go
type Parser = types.Parser
```

**Parser 解析器接口**

`Parser` 定义了命令行参数解析的标准接口, 提供了不同层次的解析功能, 从简单的参数解析到完整的命令路由执行。

**设计理念:**

- 分层设计: 提供不同层次的解析功能
- 灵活性: 支持仅解析、解析+路由等多种使用模式
- 可扩展性: 接口设计允许不同的解析策略实现

**使用场景:**

- 命令行工具的参数解析
- 子命令系统的路由管理
- 配置管理和参数验证

### RequiredGroup

```go
type RequiredGroup = types.RequiredGroup
```

**RequiredGroup**

定义了一组必需的标志，其中所有标志都必须被设置。当用户没有设置必需组中的某些标志时，解析器会返回错误。

### FlagDependency

```go
type FlagDependency = types.FlagDependency
```

**FlagDependency 标志依赖关系**

定义了标志之间的依赖关系。当触发标志被设置时，目标标志会受到约束（互斥或必需）。

**字段说明:**
  - Name: 依赖关系名称，用于错误提示和标识
  - Trigger: 触发标志名称，当此标志被设置时触发依赖检查
  - Targets: 目标标志名称列表，这些标志会受到约束
  - Type: 依赖关系类型（DepMutex 或 DepRequired）

**使用场景:**
  - 远程模式与本地路径互斥 (trigger="remote", targets=["local-path"], type=DepMutex)
  - SSL模式需要证书和密钥 (trigger="ssl", targets=["cert","key"], type=DepRequired)
  - 配置文件模式与其他配置互斥 (trigger="config", targets=["port","host"], type=DepMutex)

### DepType

```go
type DepType = types.DepType
```

**DepType 依赖关系类型**

定义了标志依赖关系的类型，用于区分互斥依赖和必需依赖。

**常量值:**
  - DepMutex: 互斥依赖，触发标志被设置时，目标标志不能被设置
  - DepRequired: 必需依赖，触发标志被设置时，目标标志必须被设置

**使用示例:**
```go
// 互斥依赖：远程模式与本地路径互斥
dep := qflag.FlagDependency{
    Name:    "remote_mutex_local",
    Trigger: "remote",
    Targets: []string{"local-path"},
    Type:    qflag.DepMutex,
}

// 必需依赖：SSL模式需要证书和密钥
dep := qflag.FlagDependency{
    Name:    "ssl_requires_cert_key",
    Trigger: "ssl",
    Targets: []string{"cert", "key"},
    Type:    qflag.DepRequired,
}
```

### SizeFlag

```go
type SizeFlag = flag.SizeFlag
```

**SizeFlag 大小标志 (支持KB、MB、GB等单位)**

`SizeFlag` 用于处理大小类型的命令行参数, 支持多种大小单位。可以解析带有单位的大小值, 并将其转换为字节数。

**支持的单位:**

- `B/b`: 字节
- `KB/kb/K/k`: 千字节 (1024字节)
- `MB/mb/M/m`: 兆字节 (1024^2字节)
- `GB/gb/G/g`: 吉字节 (1024^3字节)
- `TB/tb/T/t`: 太字节 (1024^4字节)
- `PB/pb/P/p`: 拍字节 (1024^5字节)
- `KiB/kib`: 二进制千字节 (1024字节)
- `MiB/mib`: 二进制兆字节 (1024^2字节)
- `GiB/gib`: 二进制吉字节 (1024^3字节)
- `TiB/tib`: 二进制太字节 (1024^4字节)
- `PiB/pib`: 二进制拍字节 (1024^5字节)

**注意事项:**

- 支持小数, 如 `"1.5MB"`
- 不支持负数
- 默认单位为字节(B)
- 大小写不敏感

### StringFlag

```go
type StringFlag = flag.StringFlag
```

**StringFlag 字符串标志**

`StringFlag` 用于处理字符串类型的命令行参数。它接受任何字符串值, 包括空字符串。

### StringSliceFlag

```go
type StringSliceFlag = flag.StringSliceFlag
```

**StringSliceFlag 字符串切片标志**

### TimeFlag

```go
type TimeFlag = flag.TimeFlag
```

**TimeFlag 时间标志**

`TimeFlag` 用于处理时间类型的命令行参数。支持自动检测多种常见时间格式, 也支持指定特定格式进行解析。

**特性:**

- 自动检测常见时间格式
- 支持自定义格式解析
- 记录当前使用的格式
- 线程安全的格式存储

**常见支持格式:**

- RFC3339: `"2006-01-02T15:04:05Z07:00"`
- RFC1123: `"Mon, 02 Jan 2006 15:04:05 MST"`
- 日期格式: `"2006-01-02"`, `"2006/01/02"`
- 时间格式: `"15:04:05"`, `"15:04"`
- 其他常见格式

### Uint16Flag

```go
type Uint16Flag = flag.Uint16Flag
```

**Uint16Flag 16位无符号整数标志**

`Uint16Flag` 用于处理16位无符号整数类型的命令行参数。适用于处理端口号、短范围计数器等场景。

**注意事项:**

- 只支持非负数
- 支持十进制格式
- 范围: 0 到 65,535

### Uint32Flag

```go
type Uint32Flag = flag.Uint32Flag
```

**Uint32Flag 32位无符号整数标志**

`Uint32Flag` 用于处理32位无符号整数类型的命令行参数。适用于处理IP地址、大范围计数器等场景。

**注意事项:**

- 只支持非负数
- 支持十进制格式
- 范围: 0 到 4,294,967,295

### Uint64Flag

```go
type Uint64Flag = flag.Uint64Flag
```

**Uint64Flag 64位无符号整数标志**

`Uint64Flag` 用于处理64位无符号整数类型的命令行参数。在所有平台上都使用固定的64位无符号整数, 提供一致的行为。

**注意事项:**

- 只支持非负数
- 支持十进制格式
- 范围: 0 到 18,446,744,073,709,551,615

### Uint8Flag

```go
type Uint8Flag = flag.Uint8Flag
```

**Uint8Flag 8位无符号整数标志**

`Uint8Flag` 用于处理8位无符号整数类型的命令行参数。适用于处理字节值、小范围计数器等场景。

**注意事项:**

- 只支持非负数
- 支持十进制格式
- 范围: 0 到 255

### UintFlag

```go
type UintFlag = flag.UintFlag
```

**UintFlag 无符号整数标志**

`UintFlag` 用于处理无符号整数类型的命令行参数。使用平台相关的uint类型, 在32位系统上为32位无符号整数, 在64位系统上为64位无符号整数。

**注意事项:**

- 只支持非负数
- 支持十进制格式
- 超出平台uint范围会返回错误

### Validator

```go
type Validator[T any] = types.Validator[T]
```

**Validator 验证器函数类型**

`Validator` 是一个泛型函数类型，用于验证标志值的有效性。验证器接收一个类型为 T 的值，返回错误信息。

**参数:**

- `value`: 要验证的值

**返回值:**

- `error`: 验证失败时返回错误，验证通过返回 nil

**功能说明:**

- 验证器在标志的 Set 方法中被调用
- 在解析完值后、设置值之前执行验证
- 如果验证失败，Set 方法会返回错误，标志值不会被设置
- 验证器是可选的，未设置时跳过验证
- 重复设置验证器会覆盖之前的验证器

**空值处理:**

```go
Validator: qflag.Validator[string](func(value string) error {
    if value == "" {
        return fmt.Errorf("值不能为空")
    }
    return nil
})
```

**注意事项:**

- 验证器应该快速执行，避免耗时操作
- 验证器返回的错误应该清晰描述失败原因
- 验证器执行时已经持有锁，验证器本身不需要处理并发
- 验证器可以随时通过 `ClearValidator` 清除
