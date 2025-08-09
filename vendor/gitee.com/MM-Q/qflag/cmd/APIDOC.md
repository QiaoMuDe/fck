# Package cmd

**Import Path:** `gitee.com/MM-Q/qflag/cmd`

Package cmd 提供基础标志创建和管理功能，包含命令结构体和核心功能实现。本包定义了Cmd结构体，提供命令行解析、子命令管理、标志注册等核心功能，作为适配器连接内部函数式API和外部面向对象API。

## 功能模块

- **基础标志创建和管理功能** - 提供字符串、整数、布尔、浮点数等基本类型标志的创建和绑定功能
- **命令结构体和核心功能实现** - 定义Cmd结构体，提供命令行解析、子命令管理、标志注册等核心功能
- **扩展标志类型支持** - 提供枚举、时间间隔、切片、时间、映射等高级类型标志的创建和绑定功能
- **内部实现和辅助功能** - 包含Cmd结构体的内部实现方法和辅助功能，提供命令行解析的核心逻辑

## Variables

### New

```go
var New = NewCmd
```

New 创建新的命令实例(NewCmd的简写)

## Types

### Cmd

```go
type Cmd struct {
    // Has unexported fields.
}
```

Cmd 简化的命令结构体，作为适配器连接内部函数式API和外部面向对象API

#### NewCmd

```go
func NewCmd(longName, shortName string, errorHandling flag.ErrorHandling) *Cmd
```

NewCmd 创建新的命令实例

**参数:**
- `longName`: 命令的全称(如: ls, rm, mkdir 等)
- `shortName`: 命令的简称(如: l, r, m 等)
- `errorHandling`: 标志解析错误处理策略

**返回值:**
- `*Cmd`: 新创建的命令实例

**errorHandling可选值:**
- `flag.ContinueOnError`: 遇到错误时继续解析, 并将错误返回
- `flag.ExitOnError`: 遇到错误时立即退出程序, 并将错误返回
- `flag.PanicOnError`: 遇到错误时立即触发panic, 并将错误返回

#### 示例管理方法

##### AddExample

```go
func (c *Cmd) AddExample(desc, usage string)
```

AddExample 为命令添加使用示例

**参数:**
- `desc`: 示例描述
- `usage`: 示例用法

##### AddExamples

```go
func (c *Cmd) AddExamples(examples []ExampleInfo)
```

AddExamples 为命令添加使用示例列表

**参数:**
- `examples`: 示例信息列表

##### GetExamples

```go
func (c *Cmd) GetExamples() []ExampleInfo
```

GetExamples 获取所有使用示例

**返回:**
- `[]ExampleInfo`: 使用示例列表

#### 备注管理方法

##### AddNote

```go
func (c *Cmd) AddNote(note string)
```

AddNote 添加备注信息到命令

**参数:**
- `note`: 备注信息

##### AddNotes

```go
func (c *Cmd) AddNotes(notes []string)
```

AddNotes 添加备注信息切片到命令

**参数:**
- `notes`: 备注信息列表

##### GetNotes

```go
func (c *Cmd) GetNotes() []string
```

GetNotes 获取所有备注信息

**返回:**
- 备注信息列表

#### 子命令管理方法

##### AddSubCmd

```go
func (c *Cmd) AddSubCmd(subCmds ...*Cmd) error
```

AddSubCmd 添加外部子命令到当前命令。支持批量添加多个子命令, 遇到错误时收集所有错误并返回

**参数:**
- `subCmds`: 一个或多个子命令实例指针

**返回值:**
- 错误信息, 如果所有子命令添加成功则返回nil

##### CmdExists

```go
func (c *Cmd) CmdExists(cmdName string) bool
```

CmdExists 检查子命令是否存在

**参数:**
- `cmdName`: 子命令名称

**返回:**
- `bool`: 子命令是否存在

##### SubCmdMap

```go
func (c *Cmd) SubCmdMap() map[string]*Cmd
```

SubCmdMap 返回子命令映射表(长命令名+短命令名)

**返回值:**
- `map[string]*Cmd`: 子命令映射表

##### SubCmds

```go
func (c *Cmd) SubCmds() []*Cmd
```

SubCmds 返回子命令切片

**返回值:**
- `[]*Cmd`: 子命令切片

#### 参数访问方法

##### Arg

```go
func (c *Cmd) Arg(i int) string
```

Arg 获取指定索引的非标志参数

**参数:**
- `i`: 参数索引

**返回值:**
- `string`: 指定索引位置的非标志参数；若索引越界，则返回空字符串

##### Args

```go
func (c *Cmd) Args() []string
```

Args 获取非标志参数切片

**返回值:**
- `[]string`: 参数切片

##### NArg

```go
func (c *Cmd) NArg() int
```

NArg 获取非标志参数的数量

**返回值:**
- `int`: 参数数量

#### 基础类型标志方法

##### Bool

```go
func (c *Cmd) Bool(longName, shortName string, defValue bool, usage string) *flags.BoolFlag
```

Bool 添加布尔类型标志, 返回标志对象指针

**参数值:**
- `longName`: string - 长标志名
- `shortName`: string - 短标志
- `defValue`: bool - 默认值
- `usage`: string - 帮助说明

**返回值:**
- `*flags.BoolFlag` - 布尔标志对象指针

##### BoolVar

```go
func (c *Cmd) BoolVar(f *flags.BoolFlag, longName, shortName string, defValue bool, usage string)
```

BoolVar 绑定布尔类型标志到指针并内部注册Flag对象

**参数值:**
- `f`: *flags.BoolFlag - 布尔标志对象指针
- `longName`: string - 长标志名
- `shortName`: string - 短标志
- `defValue`: bool - 默认值
- `usage`: string - 帮助说明

##### String

```go
func (c *Cmd) String(longName, shortName, defValue, usage string) *flags.StringFlag
```

String 添加字符串类型标志, 返回标志对象指针

**参数值:**
- `longName`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明

**返回值:**
- `*flags.StringFlag`: 字符串标志对象指针

##### StringVar

```go
func (c *Cmd) StringVar(f *flags.StringFlag, longName, shortName, defValue, usage string)
```

StringVar 绑定字符串类型标志到指针并内部注册Flag对象

**参数值:**
- `f`: 字符串标志指针
- `longName`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明

##### Int

```go
func (c *Cmd) Int(longName, shortName string, defValue int, usage string) *flags.IntFlag
```

Int 添加整数类型标志, 返回标志对象指针

**参数值:**
- `longName`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明

**返回值:**
- `*flags.IntFlag`: 整数标志对象指针

##### IntVar

```go
func (c *Cmd) IntVar(f *flags.IntFlag, longName, shortName string, defValue int, usage string)
```

IntVar 绑定整数类型标志到指针并内部注册Flag对象

**参数值:**
- `f`: 整数标志指针
- `longName`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明

##### Int64

```go
func (c *Cmd) Int64(longName, shortName string, defValue int64, usage string) *flags.Int64Flag
```

Int64 添加64位整数类型标志, 返回标志对象指针

**参数值:**
- `longName`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明

**返回值:**
- `*flags.Int64Flag`: 64位整数标志对象指针

##### Int64Var

```go
func (c *Cmd) Int64Var(f *flags.Int64Flag, longName, shortName string, defValue int64, usage string)
```

Int64Var 绑定64位整数类型标志到指针并内部注册Flag对象

**参数值:**
- `f`: 64位整数标志指针
- `longName`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明

##### Float64

```go
func (c *Cmd) Float64(longName, shortName string, defValue float64, usage string) *flags.Float64Flag
```

Float64 添加浮点型标志, 返回标志对象指针

**参数值:**
- `longName` - 长标志名
- `shortName` - 短标志
- `defValue` - 默认值
- `usage` - 帮助说明

**返回值:**
- `*flags.Float64Flag` - 浮点型标志对象指针

##### Float64Var

```go
func (c *Cmd) Float64Var(f *flags.Float64Flag, longName, shortName string, defValue float64, usage string)
```

Float64Var 绑定浮点型标志到指针并内部注册Flag对象

**参数值:**
- `f`: *flags.Float64Flag - 浮点型标志对象指针
- `longName`: string - 长标志名
- `shortName`: string - 短标志
- `defValue`: float64 - 默认值
- `usage`: string - 帮助说明

#### 无符号整数类型标志方法

##### Uint16

```go
func (c *Cmd) Uint16(longName, shortName string, defValue uint16, usage string) *flags.Uint16Flag
```

Uint16 添加16位无符号整数类型标志, 返回标志对象指针

**参数值:**
- `longName`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明

**返回值:**
- `*flags.Uint16Flag`: 16位无符号整数标志对象指针

##### Uint16Var

```go
func (c *Cmd) Uint16Var(f *flags.Uint16Flag, longName, shortName string, defValue uint16, usage string)
```

Uint16Var 绑定16位无符号整数类型标志到指针并内部注册Flag对象

**参数值:**
- `f`: 16位无符号整数标志指针
- `longName`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明

##### Uint32

```go
func (c *Cmd) Uint32(longName, shortName string, defValue uint32, usage string) *flags.Uint32Flag
```

Uint32 添加32位无符号整数类型标志, 返回标志对象指针

**参数值:**
- `longName`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明

**返回值:**
- `*flags.Uint32Flag`: 32位无符号整数标志对象指针

##### Uint32Var

```go
func (c *Cmd) Uint32Var(f *flags.Uint32Flag, longName, shortName string, defValue uint32, usage string)
```

Uint32Var 绑定32位无符号整数类型标志到指针并内部注册Flag对象

**参数值:**
- `f`: 32位无符号整数标志指针
- `longName`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明

##### Uint64

```go
func (c *Cmd) Uint64(longName, shortName string, defValue uint64, usage string) *flags.Uint64Flag
```

Uint64 添加64位无符号整数类型标志, 返回标志对象指针

**参数值:**
- `longName`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明

**返回值:**
- `*flags.Uint64Flag`: 64位无符号整数标志对象指针

##### Uint64Var

```go
func (c *Cmd) Uint64Var(f *flags.Uint64Flag, longName, shortName string, defValue uint64, usage string)
```

Uint64Var 绑定64位无符号整数类型标志到指针并内部注册Flag对象

**参数值:**
- `f`: 64位无符号整数标志指针
- `longName`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明

#### 扩展类型标志方法

##### Duration

```go
func (c *Cmd) Duration(longName, shortName string, defValue time.Duration, usage string) *flags.DurationFlag
```

Duration 添加时间间隔类型标志, 返回标志对象指针

**参数值:**
- `longName`: string - 长标志名
- `shortName`: string - 短标志
- `defValue`: time.Duration - 默认值
- `usage`: string - 帮助说明

**返回值:**
- `*flags.DurationFlag` - 时间间隔标志对象指针

##### DurationVar

```go
func (c *Cmd) DurationVar(f *flags.DurationFlag, longName, shortName string, defValue time.Duration, usage string)
```

DurationVar 绑定时间间隔类型标志到指针并内部注册Flag对象

**参数值:**
- `f`: *flags.DurationFlag - 时间间隔标志对象指针
- `longName`: string - 长标志名
- `shortName`: string - 短标志
- `defValue`: time.Duration - 默认值
- `usage`: string - 帮助说明

##### Enum

```go
func (c *Cmd) Enum(longName, shortName string, defValue string, usage string, options []string) *flags.EnumFlag
```

Enum 添加枚举类型标志, 返回标志对象指针

**参数值:**
- `longName`: string - 长标志名
- `shortName`: string - 短标志
- `defValue`: string - 默认值
- `usage`: string - 帮助说明
- `options`: []string - 限制该标志取值的枚举值切片

**返回值:**
- `*flags.EnumFlag` - 枚举标志对象指针

##### EnumVar

```go
func (c *Cmd) EnumVar(f *flags.EnumFlag, longName, shortName string, defValue string, usage string, options []string)
```

EnumVar 绑定枚举类型标志到指针并内部注册Flag对象

**参数值:**
- `f`: *flags.EnumFlag - 枚举标志对象指针
- `longName`: string - 长标志名
- `shortName`: string - 短标志
- `defValue`: string - 默认值
- `usage`: string - 帮助说明
- `options`: []string - 限制该标志取值的枚举值切片

##### Slice

```go
func (c *Cmd) Slice(longName, shortName string, defValue []string, usage string) *flags.SliceFlag
```

Slice 绑定字符串切片类型标志并内部注册Flag对象

**参数值:**
- `longName`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明

**返回值:**
- `*flags.SliceFlag`: 字符串切片标志对象指针

##### SliceVar

```go
func (c *Cmd) SliceVar(f *flags.SliceFlag, longName, shortName string, defValue []string, usage string)
```

SliceVar 绑定字符串切片类型标志到指针并内部注册Flag对象

**参数值:**
- `f`: 字符串切片标志指针
- `longName`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明

##### Time

```go
func (c *Cmd) Time(longName, shortName string, defValue string, usage string) *flags.TimeFlag
```

Time 添加时间类型标志, 返回标志对象指针

**参数值:**
- `longName`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值(时间表达式, 如"now", "zero", "1h", "2006-01-02")
- `usage`: 帮助说明

**返回值:**
- `*flags.TimeFlag`: 时间标志对象指针

**支持的默认值格式:**
- `"now"` 或 `""` : 当前时间
- `"zero"` : 零时间 (time.Time{})
- `"1h"`, `"30m"`, `"-2h"` : 相对时间（基于当前时间的偏移）
- `"2006-01-02"`, `"2006-01-02 15:04:05"` : 绝对时间格式
- RFC3339等标准格式

##### TimeVar

```go
func (c *Cmd) TimeVar(f *flags.TimeFlag, longName, shortName string, defValue string, usage string)
```

TimeVar 绑定时间类型标志到指针并内部注册Flag对象

**参数值:**
- `f`: 时间标志指针
- `longName`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值(时间表达式, 如"now", "zero", "1h", "2006-01-02")
- `usage`: 帮助说明

**支持的默认值格式:**
- `"now"` 或 `""` : 当前时间
- `"zero"` : 零时间 (time.Time{})
- `"1h"`, `"30m"`, `"-2h"` : 相对时间（基于当前时间的偏移）
- `"2006-01-02"`, `"2006-01-02 15:04:05"` : 绝对时间格式
- RFC3339等标准格式

##### Map

```go
func (c *Cmd) Map(longName, shortName string, defValue map[string]string, usage string) *flags.MapFlag
```

Map 添加键值对类型标志, 返回标志对象指针

**参数值:**
- `longName`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明

**返回值:**
- `*flags.MapFlag`: 键值对标志对象指针

##### MapVar

```go
func (c *Cmd) MapVar(f *flags.MapFlag, longName, shortName string, defValue map[string]string, usage string)
```

MapVar 绑定键值对类型标志到指针并内部注册Flag对象

**参数值:**
- `f`: 键值对标志指针
- `longName`: 长标志名
- `shortName`: 短标志名
- `defValue`: 默认值
- `usage`: 帮助说明

#### 标志管理方法

##### FlagExists

```go
func (c *Cmd) FlagExists(name string) bool
```

FlagExists 检查指定名称的标志是否存在

**参数:**
- `name`: 标志名称

**返回值:**
- `bool`: 标志是否存在

##### FlagRegistry

```go
func (c *Cmd) FlagRegistry() *flags.FlagRegistry
```

FlagRegistry 获取标志注册表的只读访问

**返回值:**
- `*flags.FlagRegistry`: 标志注册表的只读访问

##### NFlag

```go
func (c *Cmd) NFlag() int
```

NFlag 获取标志的数量

**返回值:**
- `int`: 标志数量

#### 命令信息方法

##### GetDescription

```go
func (c *Cmd) GetDescription() string
```

GetDescription 返回命令描述

**返回值:**
- `string`: 命令描述

##### GetHelp

```go
func (c *Cmd) GetHelp() string
```

GetHelp 返回命令用法帮助信息

**返回值:**
- `string`: 命令用法帮助信息

##### GetLogoText

```go
func (c *Cmd) GetLogoText() string
```

GetLogoText 获取logo文本

**返回值:**
- `string`: logo文本字符串

##### GetModuleHelps

```go
func (c *Cmd) GetModuleHelps() string
```

GetModuleHelps 获取自定义模块帮助信息

**返回值:**
- `string`: 自定义模块帮助信息

##### GetUsageSyntax

```go
func (c *Cmd) GetUsageSyntax() string
```

GetUsageSyntax 获取自定义命令用法

**返回值:**
- `string`: 自定义命令用法

##### GetUseChinese

```go
func (c *Cmd) GetUseChinese() bool
```

GetUseChinese 获取是否使用中文帮助信息

**返回值:**
- `bool`: 是否使用中文帮助信息

##### GetVersion

```go
func (c *Cmd) GetVersion() string
```

GetVersion 获取版本信息

**返回值:**
- `string`: 版本信息

##### LongName

```go
func (c *Cmd) LongName() string
```

LongName 返回命令长名称

**返回值:**
- `string`: 命令长名称

##### Name

```go
func (c *Cmd) Name() string
```

Name 获取命令名称

**返回值:**
- `string`: 命令名称

**说明:**
- 优先返回长名称, 如果长名称不存在则返回短名称

##### ShortName

```go
func (c *Cmd) ShortName() string
```

ShortName 返回命令短名称

**返回值:**
- `string`: 命令短名称

#### 解析方法

##### Parse

```go
func (c *Cmd) Parse(args []string) (err error)
```

Parse 完整解析命令行参数(含子命令处理)

**主要功能：**
1. 解析当前命令的长短标志及内置标志
2. 自动检测并解析子命令及其参数(若存在)
3. 验证枚举类型标志的有效性

**参数：**
- `args`: 原始命令行参数切片(包含可能的子命令及参数)

**返回值：**
- `error`: 解析过程中遇到的错误(如标志格式错误、子命令解析失败等)

**注意事项：**
- 每个Cmd实例仅会被解析一次(线程安全)
- 若检测到子命令, 会将剩余参数传递给子命令的Parse方法
- 处理内置标志执行逻辑

##### ParseFlagsOnly

```go
func (c *Cmd) ParseFlagsOnly(args []string) (err error)
```

ParseFlagsOnly 仅解析当前命令的标志参数(忽略子命令)

**主要功能：**
1. 解析当前命令的长短标志及内置标志
2. 验证枚举类型标志的有效性
3. 明确忽略所有子命令及后续参数

**参数：**
- `args`: 原始命令行参数切片(子命令及后续参数会被忽略)

**返回值：**
- `error`: 解析过程中遇到的错误(如标志格式错误等)

**注意事项：**
- 每个Cmd实例仅会被解析一次(线程安全)
- 不会处理任何子命令, 所有参数均视为当前命令的标志或位置参数
- 处理内置标志逻辑

##### IsParsed

```go
func (c *Cmd) IsParsed() bool
```

IsParsed 检查命令是否已完成解析

**返回值:**
- `bool`: 解析状态,true表示已解析(无论成功失败), false表示未解析

#### 帮助和配置方法

##### PrintHelp

```go
func (c *Cmd) PrintHelp()
```

PrintHelp 打印命令的帮助信息, 优先打印用户的帮助信息, 否则自动生成帮助信息

**注意:**
- 打印帮助信息时, 不会自动退出程序

##### SetDescription

```go
func (c *Cmd) SetDescription(desc string)
```

SetDescription 设置命令描述

**参数:**
- `desc`: 命令描述

##### SetEnableCompletion

```go
func (c *Cmd) SetEnableCompletion(enable bool)
```

SetEnableCompletion 设置是否启用自动补全, 只能在根命令上启用

**参数:**
- `enable`: true表示启用补全,false表示禁用

##### SetExitOnBuiltinFlags

```go
func (c *Cmd) SetExitOnBuiltinFlags(exit bool)
```

SetExitOnBuiltinFlags 设置是否在解析内置参数时退出。默认情况下为true, 当解析到内置参数时, QFlag将退出程序

**参数:**
- `exit`: 是否退出

##### SetHelp

```go
func (c *Cmd) SetHelp(help string)
```

SetHelp 设置用户自定义命令帮助信息

**参数:**
- `help`: 用户自定义命令帮助信息

##### SetLogoText

```go
func (c *Cmd) SetLogoText(logoText string)
```

SetLogoText 设置logo文本

**参数:**
- `logoText`: logo文本字符串

##### SetModuleHelps

```go
func (c *Cmd) SetModuleHelps(moduleHelps string)
```

SetModuleHelps 设置自定义模块帮助信息

**参数:**
- `moduleHelps`: 自定义模块帮助信息

##### SetUsageSyntax

```go
func (c *Cmd) SetUsageSyntax(usageSyntax string)
```

SetUsageSyntax 设置自定义命令用法

**参数:**
- `usageSyntax`: 自定义命令用法

##### SetUseChinese

```go
func (c *Cmd) SetUseChinese(useChinese bool)
```

SetUseChinese 设置是否使用中文帮助信息

**参数:**
- `useChinese`: 是否使用中文帮助信息

##### SetVersion

```go
func (c *Cmd) SetVersion(version string)
```

SetVersion 设置版本信息

**参数:**
- `version`: 版本信息

##### SetVersionf

```go
func (c *Cmd) SetVersionf(format string, args ...any)
```

SetVersionf 设置版本信息

**参数:**
- `format`: 版本信息格式字符串
- `args`: 格式化参数

#### 链式调用方法

##### WithDescription

```go
func (c *Cmd) WithDescription(desc string) *Cmd
```

WithDescription 设置命令描述（链式调用）

**参数:**
- `desc`: 命令描述

**返回值:**
- `*Cmd`: 返回命令实例，支持链式调用

##### WithEnableCompletion

```go
func (c *Cmd) WithEnableCompletion(enable bool) *Cmd
```

WithEnableCompletion 设置是否启用自动补全（链式调用）

**参数:**
- `enable`: true表示启用补全,false表示禁用

**返回值:**
- `*Cmd`: 返回命令实例，支持链式调用

##### WithExample

```go
func (c *Cmd) WithExample(desc, usage string) *Cmd
```

WithExample 为命令添加使用示例（链式调用）

**参数:**
- `desc`: 示例描述
- `usage`: 示例用法

**返回值:**
- `*Cmd`: 返回命令实例，支持链式调用

##### WithExamples

```go
func (c *Cmd) WithExamples(examples []ExampleInfo) *Cmd
```

WithExamples 添加使用示例列表到命令（链式调用）

**参数:**
- `examples`: 示例信息列表

**返回值:**
- `*Cmd`: 返回命令实例，支持链式调用

##### WithExitOnBuiltinFlags

```go
func (c *Cmd) WithExitOnBuiltinFlags(exit bool) *Cmd
```

WithExitOnBuiltinFlags 设置是否在解析内置参数时退出（链式调用）

**参数:**
- `exit`: 是否退出

**返回值:**
- `*Cmd`: 返回命令实例，支持链式调用

##### WithHelp

```go
func (c *Cmd) WithHelp(help string) *Cmd
```

WithHelp 设置用户自定义命令帮助信息（链式调用）

**参数:**
- `help`: 用户自定义命令帮助信息

**返回值:**
- `*Cmd`: 返回命令实例，支持链式调用

##### WithLogoText

```go
func (c *Cmd) WithLogoText(logoText string) *Cmd
```

WithLogoText 设置logo文本（链式调用）

**参数:**
- `logoText`: logo文本字符串

**返回值:**
- `*Cmd`: 返回命令实例，支持链式调用

##### WithModuleHelps

```go
func (c *Cmd) WithModuleHelps(moduleHelps string) *Cmd
```

WithModuleHelps 设置自定义模块帮助信息（链式调用）

**参数:**
- `moduleHelps`: 自定义模块帮助信息

**返回值:**
- `*Cmd`: 返回命令实例，支持链式调用

##### WithNote

```go
func (c *Cmd) WithNote(note string) *Cmd
```

WithNote 添加备注信息到命令（链式调用）

**参数:**
- `note`: 备注信息

**返回值:**
- `*Cmd`: 返回命令实例，支持链式调用

##### WithNotes

```go
func (c *Cmd) WithNotes(notes []string) *Cmd
```

WithNotes 添加备注信息切片到命令（链式调用）

**参数:**
- `notes`: 备注信息列表

**返回值:**
- `*Cmd`: 返回命令实例，支持链式调用

##### WithUsageSyntax

```go
func (c *Cmd) WithUsageSyntax(usageSyntax string) *Cmd
```

WithUsageSyntax 设置自定义命令用法（链式调用）

**参数:**
- `usageSyntax`: 自定义命令用法

**返回值:**
- `*Cmd`: 返回命令实例，支持链式调用

##### WithUseChinese

```go
func (c *Cmd) WithUseChinese(useChinese bool) *Cmd
```

WithUseChinese 设置是否使用中文帮助信息（链式调用）

**参数:**
- `useChinese`: 是否使用中文帮助信息

**返回值:**
- `*Cmd`: 返回命令实例，支持链式调用

##### WithVersion

```go
func (c *Cmd) WithVersion(version string) *Cmd
```

WithVersion 设置版本信息（链式调用）

**参数:**
- `version`: 版本信息

**返回值:**
- `*Cmd`: 返回命令实例，支持链式调用

##### WithVersionf

```go
func (c *Cmd) WithVersionf(format string, args ...any) *Cmd
```

WithVersionf 设置版本信息（链式调用，支持格式化）

**参数:**
- `format`: 版本信息格式字符串
- `args`: 格式化参数

**返回值:**
- `*Cmd`: 返回命令实例，支持链式调用

### ExampleInfo

```go
type ExampleInfo = types.ExampleInfo
```

ExampleInfo 导出示例信息类型

## 使用示例

### 基本用法

```go
// 创建命令
cmd := NewCmd("myapp", "app", flag.ExitOnError)

// 添加标志
verbose := cmd.Bool("verbose", "v", false, "启用详细输出")
output := cmd.String("output", "o", "", "输出文件路径")

// 解析参数
err := cmd.Parse(os.Args[1:])
if err != nil {
    log.Fatal(err)
}

// 使用标志值
if *verbose {
    fmt.Println("详细模式已启用")
}
```

### 链式调用

```go
cmd := NewCmd("myapp", "app", flag.ExitOnError).
    WithDescription("我的应用程序").
    WithVersion("1.0.0").
    WithUseChinese(true).
    WithExample("基本用法", "myapp --verbose --output result.txt")
```

### 子命令

```go
// 创建主命令
rootCmd := NewCmd("git", "g", flag.ExitOnError)

// 创建子命令
addCmd := NewCmd("add", "a", flag.ExitOnError)
commitCmd := NewCmd("commit", "c", flag.ExitOnError)

// 添加子命令
err := rootCmd.AddSubCmd(addCmd, commitCmd)
if err != nil {
    log.Fatal(err)
}
```

### 扩展类型标志

```go
// 枚举类型
logLevel := cmd.Enum("log-level", "l", "info", "设置日志级别", 
    []string{"debug", "info", "warn", "error"})

// 时间类型
startTime := cmd.Time("start-time", "s", "now", "设置开始时间")

// 切片类型
tags := cmd.Slice("tags", "t", []string{}, "设置标签列表")

// 映射类型
env := cmd.Map("env", "e", map[string]string{}, "设置环境变量")
```

## 注意事项

1. **解析状态**: 每个Cmd实例只能被解析一次，多次调用Parse方法会被忽略
2. **线程安全**: 解析过程是线程安全的
3. **错误处理**: 根据创建时指定的errorHandling策略处理解析错误
4. **内置标志**: 支持--help、--version等内置标志的自动处理
5. **自动补全**: 可在根命令上启用Shell自动补全功能
6. **枚举验证**: 枚举类型标志会自动验证输入值是否在允许的选项中
7. **时间格式**: 时间类型标志支持多种时间格式，包括相对时间和绝对时间

## 相关包

- `gitee.com/MM-Q/qflag/flags` - 标志类型定义
- `gitee.com/MM-Q/qflag/types` - 通用类型定义
- `flag` - Go标准库标志包（用于错误处理策略）
