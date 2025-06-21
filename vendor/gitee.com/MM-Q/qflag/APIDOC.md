# qflag API文档

项目地址: [https://gitee.com/MM-Q/qflag](https://gitee.com/MM-Q/qflag)

qflag是一个Go语言命令行参数解析库，提供了比标准库flag更丰富的功能，包括长短标志绑定、子命令支持、自动帮助信息生成等特性。

## 概述

qflag是一个Go语言命令行参数解析库，提供了比标准库flag更丰富的功能，包括长短标志绑定、子命令支持、自动帮助信息生成等特性。

## 核心类型

### FlagType枚举

定义标志的类型常量，用于标识不同种类的命令行标志。

```go
type FlagType int

const (
  FlagTypeInt      FlagType = iota + 1 // 整数类型
  FlagTypeString                       // 字符串类型
  FlagTypeBool                         // 布尔类型
  FlagTypeFloat                        // 浮点数类型

  FlagTypeEnum                         // 枚举类型
  FlagTypeDuration                     // 时间间隔类型
)
```

### Flag接口

所有标志类型的通用接口，定义了标志的基本属性访问方法。

```go
type Flag interface {
  LongName() string   // 获取标志的长名称
  ShortName() string  // 获取标志的短名称
  Usage() string      // 获取标志的用法
  Type() FlagType     // 获取标志类型
  GetDefaultAny() any // 获取标志的默认值
}

### TypedFlag接口
带类型的标志接口，扩展了Flag接口以支持类型安全的值访问。

type TypedFlag[T any] interface {
  Flag
  GetDefault() T // 获取类型化的默认值
  Get() T        // 获取类型化的当前值
  Set(T) error   // 设置类型化的值
}

### BaseFlag结构体
泛型基础标志结构体，所有具体标志类型的基类，封装了通用字段和方法。

#### 字段说明
- `cmd *Cmd` 所属命令实例的引用
- `longName string` 长标志名称
- `shortName string` 短标志字符
- `defValue T` 默认值
- `usage string` 帮助说明
- `value *T` 标志值指针
- `mu sync.Mutex` 并发访问锁

```go
type BaseFlag[T any] struct {
  cmd       *Cmd       // 所属的命令实例
  longName  string     // 长标志名称
  shortName string     // 短标志字符
  defValue  T          // 默认值
  usage     string     // 帮助说明
  value     *T         // 标志值指针
  mu        sync.Mutex // 并发访问锁
}

// 通用方法
func (f *BaseFlag[T]) Get() T              // 获取标志值
func (f *BaseFlag[T]) Set(value T) error   // 设置标志值
func (f *BaseFlag[T]) GetDefault() T       // 获取默认值
```

### FlagMetaInterface接口

定义标志元数据的标准访问方法。

#### 方法说明

- `GetDefault() any` 获取标志的默认值

```go
type FlagMetaInterface interface {
  GetFlagType() FlagType // 获取标志类型
  GetFlag() Flag         // 获取标志对象
  GetLongName() string   // 获取长名称
  GetShortName() string  // 获取短名称
  GetUsage() string      // 获取用法描述
  GetDefault() any       // 获取默认值
  GetValue() any         // 获取当前值
}

### FlagMeta结构体
实现FlagMetaInterface接口，统一存储标志的完整元数据。
统一存储标志的完整元数据，包括标志对象及其类型信息。

```go
type FlagMeta struct {
  flag Flag // 标志对象
}

// 元数据访问方法
func (m *FlagMeta) GetLongName() string   // 获取长名称
func (m *FlagMeta) GetShortName() string  // 获取短名称
func (m *FlagMeta) GetUsage() string      // 获取用法描述
func (m *FlagMeta) GetFlagType() FlagType // 获取标志类型
func (m *FlagMeta) GetDefault() any       // 获取默认值
```

### FlagRegistryInterface接口

定义标志注册表的标准操作方法。

```go
type FlagRegistryInterface interface {
  GetAllFlags() []*FlagMeta                      // 获取所有标志元数据
  GetLongFlags() map[string]*FlagMeta            // 获取长标志映射
  GetShortFlags() map[string]*FlagMeta           // 获取短标志映射
  RegisterFlag(meta *FlagMeta) error             // 注册标志元数据
  GetByLong(longName string) (*FlagMeta, bool)   // 按长名称查找
  GetByShort(shortName string) (*FlagMeta, bool) // 按短名称查找
  GetByName(name string) (*FlagMeta, bool)       // 按名称查找
}

### FlagRegistry结构体
实现FlagRegistryInterface接口，集中管理所有标志元数据及索引。
集中管理所有标志元数据及索引，提供线程安全的标志注册和查询功能。

#### 主要方法
- `GetLongFlags() map[string]*FlagMeta` 获取所有长标志的元数据映射
- `GetShortFlags() map[string]*FlagMeta` 获取所有短标志的元数据映射

```go
type FlagRegistry struct {
  mu       sync.RWMutex         // 并发访问锁
  byLong   map[string]*FlagMeta // 按长名称索引
  byShort  map[string]*FlagMeta // 按短名称索引
  allFlags []*FlagMeta          // 所有标志元数据列表
}

// 注册表方法
func (r *FlagRegistry) RegisterFlag(meta *FlagMeta) error             // 注册标志
func (r *FlagRegistry) GetByLong(longName string) (*FlagMeta, bool)   // 按长名称查找
func (r *FlagRegistry) GetByShort(shortName string) (*FlagMeta, bool) // 按短名称查找
func (r *FlagRegistry) GetByName(name string) (*FlagMeta, bool)       // 按名称查找
func (r *FlagRegistry) GetAllFlags() []*FlagMeta                      // 获取所有标志


### Flag接口
`Flag`是所有标志类型的通用接口，定义了标志的元数据访问方法。

```go
type Flag interface {
    LongName() string   // 获取标志的长名称
    ShortName() string  // 获取标志的短名称
    Usage() string      // 获取标志的用法
    Type() FlagType     // 获取标志类型
    getDefaultAny() any // 获取默认值(内部使用)
}
```

### TypedFlag接口

`TypedFlag`是带类型的标志接口，继承自 `Flag`并提供类型化的默认值和值访问方法。

```go
type TypedFlag[T any] interface {
    Flag
    GetDefault() T // 获取标志的默认值
    Get() T        // 获取标志的实际值
    Set(T)         // 设置标志的值
}
```

### 具体标志类型

qflag提供以下具体标志类型，均实现了 `TypedFlag`接口：

- `StringFlag`: 字符串类型标志
- `IntFlag`: 整数类型标志
- `BoolFlag`: 布尔类型标志
- `FloatFlag`: 浮点数类型标志
- `EnumFlag`: 枚举类型标志，限制输入值为预定义选项集合

每个标志类型都有对应的 `GetDefault()`、`GetValue()`和 `SetValue()`方法。

### Cmd结构体

`Cmd`是qflag库的核心结构体，实现了 `Command`接口，用于管理命令行标志和子命令。

#### 主要字段

- `fs *flag.FlagSet`: 底层flag集合，处理参数解析
- `name string`: 命令名称
- `shortName string`: 命令短名称
- `description string`: 命令描述
- `usage string`: 自定义帮助内容
- `subCmds []*Cmd`: 子命令列表
- `useChinese bool`: 控制是否使用中文帮助信息
- `notes []string`: 存储备注内容

#### 主要方法

##### GetUseChinese

```go
func (c *Cmd) GetUseChinese() bool
```

获取是否使用中文帮助信息的状态。

返回值:

- `bool`: 当前是否启用中文帮助信息

##### SetUseChinese

```go
func (c *Cmd) SetUseChinese(useChinese bool)
```

设置是否使用中文帮助信息。

参数:

- `useChinese`: 为true时启用中文帮助信息

##### Args

```go
func (c *Cmd) Args() []string
```

获取非标志参数切片。

返回值:

- `[]string`: 非标志参数切片

##### Arg

```go
func (c *Cmd) Arg(i int) string
```

获取指定索引的非标志参数。

参数:

- `i`: 参数索引

返回值:

- `string`: 指定索引的参数值，若索引无效则返回空字符串

##### NArg

```go
func (c *Cmd) NArg() int
```

##### AddNote

```go
func (c *Cmd) AddNote(note string)
```

为命令添加备注信息，备注将显示在帮助信息的注意事项部分。
获取非标志参数的数量。

返回值:

- `int`: 非标志参数数量

##### NFlag

```go
func (c *Cmd) NFlag() int
```

获取已设置的标志数量。

返回值:

- `int`: 已设置的标志数量

##### FlagExists

```go
func (c *Cmd) FlagExists(name string) bool
```

检查指定名称的标志是否存在。

参数:

- `name`: 标志名称（长标志或短标志）

返回值:

- `bool`: 标志是否存在

##### AddExample

```go
func (c *Cmd) AddExample(e ExampleInfo)
```

为命令添加示例信息，示例将显示在帮助信息的示例部分。

参数:

- `e`: 示例信息对象

##### GetExamples

```go
func (c *Cmd) GetExamples() []ExampleInfo
```

获取命令的示例信息列表，返回所有添加的示例信息。

返回值:

- `[]ExampleInfo`: 示例信息列表副本

##### NewCmd

```go
func NewCmd(name string, shortName string, errorHandling flag.ErrorHandling) *Cmd
```

创建新的命令实例。

参数:

- `name`: 命令名称
- `shortName`: 命令短名称
- `errorHandling`: 错误处理方式(flag.ContinueOnError, flag.ExitOnError, flag.PanicOnError)

返回值:

- `*Cmd`: 新创建的命令实例

##### AddSubCmd

```go
func (c *Cmd) AddSubCmd(subCmds ...*Cmd) error
```

为当前命令添加子命令。

参数:

- `subCmds`: 一个或多个子命令

返回值:

- `error`: 如果检测到循环引用或nil子命令则返回错误

##### AddMutexGroup

```go

```

为当前命令添加标志互斥组，互斥组内的标志不能同时被设置。

参数:

- `flags`: 构成互斥组的一个或多个标志实例

返回值:

- `error`: 如果标志为nil或不属于当前命令则返回错误

##### Parse

```go
func (c *Cmd) Parse(args []string) error
```

解析命令行参数。

参数:

- `args`: 命令行参数切片

返回值:

- `error`: 解析过程中遇到的错误

##### PrintUsage

```go
func (c *Cmd) PrintUsage()
```

打印命令的帮助信息，优先使用自定义帮助内容，否则自动生成。

## 标志操作函数

### 字符串标志

#### String

```go
func String(name, shortName, defValue, usage string) *StringFlag
```

创建字符串类型标志(全局默认命令)。

参数:

- `name`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明

返回值:

- `*StringFlag`: 字符串标志实例

#### StringVar

```go
func StringVar(p *string, name, shortName, defValue, usage string)
```

绑定字符串类型标志到指针(全局默认命令)。

参数:

- `p`: 指向字符串变量的指针
- `name`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明

### 整数标志

#### Int

```go
func Int(name, shortName string, defValue int, usage string) *IntFlag
```

创建整数类型标志(全局默认命令)。

参数:

- `name`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明

返回值:

- `*IntFlag`: 整数标志实例

#### IntVar

```go
func IntVar(p *int, name, shortName string, defValue int, usage string)
```

绑定整数类型标志到指针(全局默认命令)。

参数:

- `p`: 指向整数变量的指针
- `name`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明

### 布尔标志

#### Bool

```go
func Bool(name, shortName string, defValue bool, usage string) *BoolFlag
```

创建布尔类型标志(全局默认命令)。

参数:

- `name`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明

返回值:

- `*BoolFlag`: 布尔标志实例

#### BoolVar

```go
func BoolVar(p *bool, name, shortName string, defValue bool, usage string)
```

绑定布尔类型标志到指针(全局默认命令)。

参数:

- `p`: 指向布尔变量的指针
- `name`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明

### 浮点数标志

#### Float

```go
func Float(name, shortName string, defValue float64, usage string) *FloatFlag
```

创建浮点数类型标志(全局默认命令)。

参数:

- `name`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明

返回值:

- `*FloatFlag`: 浮点数标志实例

#### FloatVar

```go
func FloatVar(p *float64, name, shortName string, defValue float64, usage string)
```

绑定浮点数类型标志到指针(全局默认命令)。

参数:

- `p`: 指向浮点数变量的指针
- `name`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明

### 枚举标志

#### EnumFlag特有方法

```go
func (f *EnumFlag) IsCheck(value string) error
```

验证值是否在枚举选项范围内。

参数:

- `value`: 待验证的值

返回值:

- `error`: 验证失败时返回错误信息

### 时间间隔标志

#### DurationFlag特有方法

```go
func (f *DurationFlag) Set(value string) error
```

解析并设置时间间隔值，支持ns/us/ms/s/m/h等单位。

参数:

- `value`: 时间间隔字符串

返回值:

- `error`: 解析失败时返回错误信息

#### Enum

```go
func Enum(name, shortName string, defValue string, usage string, allowedValues ...string) *EnumFlag
```

创建枚举类型标志(全局默认命令)，限制输入值为预定义选项集合。

参数:

- `name`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明
- `allowedValues`: 允许的枚举值列表

返回值:

- `*EnumFlag`: 枚举标志实例

#### EnumVar

```go
func EnumVar(p *string, name, shortName string, defValue string, usage string, allowedValues ...string)
```

绑定枚举类型标志到指针(全局默认命令)，限制输入值为预定义选项集合。

参数:

- `p`: 指向字符串变量的指针
- `name`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明
- `allowedValues`: 允许的枚举值列表

### 时间间隔标志

#### Duration

```go
func Duration(name, shortName string, defValue time.Duration, usage string) *DurationFlag
```

创建时间间隔类型标志(全局默认命令)，支持解析时间单位如"s"(秒), "m"(分钟), "h"(小时)等。

参数:

- `name`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值，类型为time.Duration
- `usage`: 帮助说明

返回值:

- `*DurationFlag`: 时间间隔标志实例

#### DurationVar

```go
func DurationVar(p *time.Duration, name, shortName string, defValue time.Duration, usage string)
```

绑定时间间隔类型标志到指针(全局默认命令)，支持解析时间单位如"s"(秒), "m"(分钟), "h"(小时)等。

参数:

- `p`: 指向time.Duration变量的指针
- `name`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值，类型为time.Duration
- `usage`: 帮助说明

## 内置标志

qflag自动绑定了以下内置标志：

### 帮助标志

- 长标志: `--help`
- 短标志: `-h`
- 功能: 显示命令的帮助信息
- 使用示例: `myapp --help` 或 `myapp -h`

### 显示安装路径标志

- 长标志: `--show-install-path`
- 短标志: `-sip`
- 功能: 显示应用程序的安装路径
- 使用示例: `myapp --show-install-path` 或 `myapp -sip`

## 高级特性

### 线程安全

qflag库使用 `sync.Mutex`和 `sync.Once`确保所有标志操作和解析过程是线程安全的，可以在并发环境中安全使用。

### 循环引用检测

添加子命令时，qflag会自动检测命令间的循环引用，避免出现无限递归的命令结构。

### 动态帮助信息生成

qflag会根据命令和标志的定义自动生成格式化的帮助信息，包括命令描述、标志说明、子命令列表等。支持中英文双语切换，通过SetUseChinese方法控制。### 标志注册表
FlagRegistry提供集中式标志管理，支持通过名称快速查找标志元数据，包括标志类型、默认值、当前值等信息。

### 标志命名规则

qflag对标志名称有严格的字符限制，禁止使用以下字符：

```go
const invalidFlagChars = " !@#$%^&*(){}[]|\\;:'\"<>,.?/"
```

命名建议：

- 使用小写字母和连字符(-)组合，如 `--config-path`
- 短标志建议使用单个字母，如 `-c`对应 `--config`
- 避免使用保留标志名称：`help`、`h`、`show-install-path`、`sip`
  FlagRegistry提供集中式标志管理，支持通过名称快速查找标志元数据，包括标志类型、默认值、当前值等信息。
  qflag会根据命令和标志的定义自动生成格式化的帮助信息，包括命令描述、标志说明、子命令列表等。

## 使用示例

### 基本用法

```go
package main

import (
  "fmt"
  "gitee.com/MM-Q/qflag"
)

func main() {
  // 创建字符串标志
  nameFlag := qflag.String("name", "n", "", "Your name")
  // 创建整数标志
  ageFlag := qflag.Int("age", "a", 0, "Your age")
  // 创建布尔标志
  verboseFlag := qflag.Bool("verbose", "v", false, "Verbose output")

  // 解析命令行参数
  if err := qflag.Parse(); err != nil {
    fmt.Println("Error parsing flags:", err)
    return
  }

  // 使用标志值
  fmt.Printf("Hello, %s! You are %d years old.\n", nameFlag.Get(), ageFlag.Get())
  if verboseFlag.GetValue() {
    fmt.Println("Verbose mode enabled")
    // 访问默认值示例
    fmt.Printf("Default verbose value: %v\n", verboseFlag.GetDefault())
    // 获取当前值示例
    fmt.Printf("Current verbose value: %v\n", verboseFlag.Get())
  }
}
```

### 自定义帮助信息

```go
package main

import (
  "fmt"
  "gitee.com/MM-Q/qflag"
  "flag"
)

func main() {
  // 创建自定义命令
  cmd := qflag.NewCmd("greet", "g", flag.ExitOnError)
  cmd.SetDescription("A simple greeting program")
  
  // 设置自定义帮助信息
  cmd.SetUsage(`Usage: greet [options] <name>

A simple program to greet people.

Options:
  -n, --name <name>   Your name (required)
  -a, --age <age>     Your age
  -h, --help          Show this help message`)
  
  // 创建标志
  name := cmd.String("name", "n", "", "Your name")
  age := cmd.Int("age", "a", 0, "Your age")
  
  // 解析参数
  if err := cmd.Parse(os.Args[1:]); err != nil {
    fmt.Println("Error parsing flags:", err)
    return
  }
  
  // 显示帮助信息（当--help/-h被使用时自动调用）
  if *cmd.helpFlag {
    cmd.PrintUsage()
    return
  }
  
  // 检查必填参数
  if *name == "" {
    fmt.Println("Error: name is required")
    cmd.PrintUsage()
    return
  }
  
  // 执行逻辑
  fmt.Printf("Hello, %s!", *name)
  if *age > 0 {
    fmt.Printf(" You are %d years old.", *age)
  }
  fmt.Println()
}
```

### 内置标志说明

```go
package main

import (
  "fmt"
  "gitee.com/MM-Q/qflag"
  "os"
)

func main() {
  // 创建标志
  name := qflag.String("name", "n", "", "Your name")
  age := qflag.Int("age", "a", 0, "Your age")

  // 解析命令行参数
  if err := qflag.Parse(); err != nil {
    fmt.Println("Error parsing flags:", err)
    os.Exit(1)
  }

  // 检查是否请求帮助
  if *qflag.QCommandLine.helpFlag {
    qflag.PrintUsage()
    return
  }

  // 检查是否请求显示安装路径
  if *qflag.QCommandLine.showInstallPathFlag {
    fmt.Println("Installation path: C:\Program Files\MyApp")
    return
  }

  // 使用标志值
  fmt.Printf("Hello, %s! You are %d years old.\n", *name, *age)
}
```

使用示例:

```bash
# 显示帮助信息
myapp -h
myapp --help

# 显示安装路径
myapp -sip
myapp --show-install-path

# 正常运行
myapp -n Alice -a 30
```

### 错误处理示例

```go
package main

import (
  "fmt"
  "gitee.com/MM-Q/qflag"
  "flag"
)

func main() {
  // 创建自定义命令并设置错误处理方式
  cmd := qflag.NewCmd("myapp", "ma", flag.ContinueOnError)
  cmd.Int("port", "p", 8080, "Server port")

  // 解析参数
  err := cmd.Parse(os.Args[1:])
  if err != nil {
    // 处理解析错误
    if err == flag.ErrHelp {
      // 帮助信息已显示，无需额外处理
      return
    }
    fmt.Printf("Error parsing command line arguments: %v\n", err)
    cmd.PrintUsage()
    os.Exit(1)
  }

  // 程序逻辑...
}
```

## 错误处理

### 错误常量

qflag定义了以下错误常量，用于标识不同类型的解析错误：

| 常量名称                     | 描述                 |
| ---------------------------- | -------------------- |
| `ErrFlagParseFailed`       | 全局实例标志解析错误 |
| `ErrSubCommandParseFailed` | 子命令标志解析错误   |
| `ErrPanicRecovered`        | 恐慌捕获错误         |

## API摘要

| 类别       | 名称                                 | 描述               |
| ---------- | ------------------------------------ | ------------------ |
| 命令创建   | NewCmd                               | 创建新的命令实例   |
| 标志操作   | String, Int, Bool, Float, Enum       | 创建指定类型的标志 |
| 标志绑定   | StringVar, IntVar, BoolVar, FloatVar | 绑定标志到变量指针 |
| 参数解析   | Parse                                | 解析命令行参数     |
| 帮助信息   | PrintUsage                           | 打印命令帮助信息   |
| 子命令管理 | AddSubCmd                            | 为命令添加子命令   |
| 参数访问   | Args, Arg, NArg                      | 获取非标志参数     |
| 标志查询   | NFlag, FlagExists                    | 查询标志状态       |

## 参数优先级规则

当长标志和短标志同时使用时，后指定的标志将覆盖先指定的标志。例如：

```bash
myapp --name Alice -n Bob
```

上述命令中，最终的name值为"Bob"，因为短标志 `-n`在长标志 `--name`之后指定。

## 使用注意事项

1. 标志名称和短名称不能为空，且不能与内置标志冲突
2. 添加子命令时要避免循环引用
3. Parse方法应只调用一次
4. 布尔类型标志不需要显式赋值， presence即为true
5. 使用自定义错误处理方式时，需手动处理帮助标志和错误信息

### 子命令用法

```go
package main

import (
  "fmt"
  "os"
  "qflag"
  "flag"
)

func main() {
  // 创建主命令
  rootCmd := qflag.NewCmd("myapp", "ma", flag.ExitOnError)
  rootCmd.SetDescription("My application description")
  rootCmd.String("log-level", "l", "info", "Set log level (debug, info, warn, error)")

  // 创建子命令
  serverCmd := qflag.NewCmd("server", "s", flag.ExitOnError)
  serverCmd.SetDescription("Start the server")
  serverCmd.Int("port", "p", 8080, "Server port")
  serverCmd.Bool("tls", "t", false, "Enable TLS")

  // 创建嵌套子命令
  configCmd := qflag.NewCmd("config", "c", flag.ExitOnError)
  configCmd.SetDescription("Configure the server")
  configCmd.String("file", "f", "config.json", "Config file path")

  // 构建命令层级
  serverCmd.AddSubCmd(configCmd)
  rootCmd.AddSubCmd(serverCmd)

  // 解析命令行参数
  if err := rootCmd.Parse(os.Args[1:]); err != nil {
    fmt.Println("Error parsing flags:", err)
    return
  }

  // 处理命令
  if len(rootCmd.Args()) > 0 {
    switch rootCmd.Args()[0] {
    case "server":
      fmt.Printf("Starting server on port %d\n", serverCmd.Arg(0))
      if len(serverCmd.Args()) > 0 && serverCmd.Args()[0] == "config" {
        fmt.Printf("Using config file: %s\n", configCmd.Arg(0))
      }
    }
  }
}
```

### 内置标志说明

```go
package main

import (
  "fmt"
  "gitee.com/MM-Q/qflag"
  "flag"
)

func main() {
  // 创建主命令
  rootCmd := qflag.NewCmd("myapp", "ma", flag.ExitOnError)
  rootCmd.SetDescription("My application description")

  // 创建子命令
  subCmd := qflag.NewCmd("sub", "s", flag.ExitOnError)
  subCmd.SetDescription("Subcommand description")
  subCmd.String("config", "c", "config.json", "Config file path")

  // 添加子命令
  rootCmd.AddSubCmd(subCmd)

  // 解析命令行参数
  if err := rootCmd.Parse(os.Args[1:]); err != nil {
    fmt.Println("Error parsing flags:", err)
    return
  }

  // 处理子命令
  if len(rootCmd.Args()) > 0 && rootCmd.Args()[0] == "sub" {
    fmt.Println("Running subcommand")
    // 可以通过subCmd访问子命令的标志
  }
}
```
