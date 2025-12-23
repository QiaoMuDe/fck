# qflag API文档

## 包概述

```go
package qflag // import "gitee.com/MM-Q/qflag"
```

### 基础标志创建和管理功能
本文件提供了Cmd结构体的基础标志创建方法，包括字符串、整数、布尔、浮点数等基本类型标志的创建和绑定功能。

### 命令结构体和核心功能实现
本文件定义了Cmd结构体，提供命令行解析、子命令管理、标志注册等核心功能。
Cmd作为适配器连接内部函数式API和外部面向对象API。

### 根包统一导出入口
本文件用于将各子包的核心功能导出到根包，简化外部使用。通过类型别名和变量导出的方式，为用户提供统一的API接口。

### 扩展标志类型支持
本文件提供了Cmd结构体的扩展标志创建方法，包括枚举、时间间隔、切片、时间、映射等高级类型标志的创建和绑定功能。

### 内部实现
cmd_internal 包含 Cmd 的内部实现细节，不对外暴露


## 函数

### func Parse() error
```go
func Parse() error
```
Parse 解析所有命令行参数, 包括根命令和所有子命令的标志参数

**返回：**
  - error: 解析过程中遇到的错误, 若成功则为 nil

### func ParseFlagsOnly() error
```go
func ParseFlagsOnly() error
```
ParseFlagsOnly 解析根命令的所有标志参数, 不包括子命令

**返回：**
  - error: 解析过程中遇到的错误, 若成功则为 nil

### func ParseAndRoute() error
```go
func ParseAndRoute() error
```
ParseAndRoute 解析参数并自动路由执行子命令
这是推荐使用的命令行参数处理方式，会自动处理子命令路由

**返回：**
  - error: 执行过程中遇到的错误, 若成功则为 nil

### func ApplyConfig(config types.CmdConfig)
```go
func ApplyConfig(config types.CmdConfig)
```
ApplyConfig 批量设置根命令配置
通过传入一个CmdConfig结构体来一次性设置多个配置项
这是全局函数，直接操作全局根命令实例

**参数:**
  - config: 包含所有配置项的CmdConfig结构体

### func AddSubCmd(subCmds ...*Cmd) error
```go
func AddSubCmd(subCmds ...*Cmd) error
```
AddSubCmd 向根命令添加一个或多个子命令
这是全局函数，直接操作全局根命令实例

此方法会对所有子命令进行完整性验证，包括名称冲突检查、循环依赖检测等。所有验证通过后，子命令将被注册到根命令的子命令映射表和列表中。
操作过程中会自动设置子命令的父命令引用，确保命令树结构的完整性。

**参数:**
  - subCmds: 要添加的子命令实例指针，支持传入多个子命令进行批量添加

**返回值:**
  - error: 添加过程中的错误信息。如果任何子命令验证失败，将返回包含所有错误详情的聚合错误；
    如果所有子命令成功添加，返回 nil

### func AddSubCmds(subCmds []*Cmd) error
```go
func AddSubCmds(subCmds []*Cmd) error
```
AddSubCmds 向根命令添加子命令切片的便捷方法
这是全局函数，直接操作全局根命令实例

此方法是 AddSubCmd 的便捷包装，专门用于处理子命令切片。内部直接调用 AddSubCmd 方法，具有相同的验证逻辑和并发安全特性。

**参数:**
  - subCmds: 子命令切片，包含要添加的所有子命令实例指针

**返回值:**
  - error: 添加过程中的错误信息，与 AddSubCmd 返回的错误类型相同

## 类型

### type BoolFlag = flags.BoolFlag
```go
type BoolFlag = flags.BoolFlag
```
BoolFlag 导出flag包中的BoolFlag结构体

### type Cmd struct
```go
type Cmd struct {
  // Has unexported fields.
}
```
Cmd 命令结构体，作为适配器连接内部函数式API和外部面向对象API

**字段说明:**
- Run: 执行函数接口，用于定义命令的执行逻辑。当命令需要执行时，会调用此函数并传入命令实例本身

### var Root *Cmd
```go
var Root *Cmd
```
Root 全局根命令实例, 提供对全局标志和子命令的访问 用户可以通过 qflag.Root.String() 这样的方式直接创建全局标志
这是访问命令行功能的主要入口点, 推荐优先使用

### func NewCmd(longName, shortName string, errorHandling ErrorHandling) *Cmd
```go
func NewCmd(longName, shortName string, errorHandling ErrorHandling) *Cmd
```
NewCmd 创建新的命令实例

**参数:**
  - longName: 命令的全称(如: ls, rm, mkdir 等)
  - shortName: 命令的简称(如: l, r, m 等)
  - errorHandling: 标志解析错误处理策略

**返回值:**
  - *Cmd: 新创建的命令实例

**errorHandling可选值:**
  - qflag.ContinueOnError: 遇到错误时继续解析, 并将错误返回
  - qflag.ExitOnError: 遇到错误时立即退出程序, 并将错误返回
  - qflag.PanicOnError: 遇到错误时立即触发panic, 并将错误返回

### func (c *Cmd) AddExample(desc, usage string)
```go
func (c *Cmd) AddExample(desc, usage string)
```
AddExample 为命令添加使用示例

**参数:**
  - desc: 示例描述
  - usage: 示例用法

### func (c *Cmd) AddExamples(examples []ExampleInfo)
```go
func (c *Cmd) AddExamples(examples []ExampleInfo)
```
AddExamples 为命令添加使用示例列表

**参数:**
  - examples: 示例信息列表

### func (c *Cmd) AddNote(note string)
```go
func (c *Cmd) AddNote(note string)
```
AddNote 添加备注信息到命令

**参数:**
  - note: 备注信息

### func (c *Cmd) AddNotes(notes []string)
```go
func (c *Cmd) AddNotes(notes []string)
```
AddNotes 添加备注信息切片到命令

**参数:**
  - notes: 备注信息列表

### func (c *Cmd) AddSubCmd(subCmds ...*Cmd) error
```go
func (c *Cmd) AddSubCmd(subCmds ...*Cmd) error
```
AddSubCmd 向当前命令添加一个或多个子命令

此方法会对所有子命令进行完整性验证，包括名称冲突检查、循环依赖检测等。所有验证通过后，子命令将被注册到当前命令的子命令映射表和列表中。
操作过程中会自动设置子命令的父命令引用，确保命令树结构的完整性。

**并发安全:** 此方法使用互斥锁保护，可安全地在多个 goroutine 中并发调用。

**参数:**
  - subCmds: 要添加的子命令实例指针，支持传入多个子命令进行批量添加

**返回值:**
  - error: 添加过程中的错误信息。如果任何子命令验证失败，将返回包含所有错误详情的聚合错误；如果所有子命令成功添加，返回 nil

**错误类型:**
  - ValidationError: 子命令为空、名称冲突、循环依赖等验证错误
  - 其他错误: 内部状态异常或系统错误

**使用示例:**

```go
cmd := qflag.NewCmd("parent", "p", "父命令")
subCmd1 := qflag.NewCmd("child1", "c1", "子命令1")
subCmd2 := qflag.NewCmd("child2", "c2", "子命令2")

if err := cmd.AddSubCmd(subCmd1, subCmd2); err != nil {
    log.Fatal(err)
}
```

### func (c *Cmd) AddSubCmds(subCmds []*Cmd) error
```go
func (c *Cmd) AddSubCmds(subCmds []*Cmd) error
```
AddSubCmds 向当前命令添加子命令切片的便捷方法

此方法是 AddSubCmd 的便捷包装，专门用于处理子命令切片。内部直接调用 AddSubCmd 方法，具有相同的验证逻辑和并发安全特性。

**并发安全:** 此方法通过调用 AddSubCmd 实现，继承其互斥锁保护特性。

**参数:**
  - subCmds: 子命令切片，包含要添加的所有子命令实例指针

**返回值:**
  - error: 添加过程中的错误信息，与 AddSubCmd 返回的错误类型相同

**使用示例:**

```go
cmd := qflag.NewCmd("parent", "p", "父命令")
subCmds := []*qflag.Cmd{
    qflag.NewCmd("child1", "c1", "子命令1"),
    qflag.NewCmd("child2", "c2", "子命令2"),
}

if err := cmd.AddSubCmds(subCmds); err != nil {
    log.Fatal(err)
}
```

### func (c *Cmd) ApplyConfig(config CmdConfig)
```go
func (c *Cmd) ApplyConfig(config CmdConfig)
```
ApplyConfig 批量设置命令配置 通过传入一个CmdConfig结构体来一次性设置多个配置项

**参数:**
  - config: 包含所有配置项的CmdConfig结构体

### func (c *Cmd) Arg(i int) string
```go
func (c *Cmd) Arg(i int) string
```
Arg 获取指定索引的非标志参数

**参数:**
  - i: 参数索引

**返回值:**
  - string: 指定索引位置的非标志参数；若索引越界，则返回空字符串

### func (c *Cmd) Args() []string
```go
func (c *Cmd) Args() []string
```
Args 获取非标志参数切片

**返回值:**
  - []string: 参数切片

### func (c *Cmd) Bool(longName, shortName string, defValue bool, usage string) *BoolFlag
```go
func (c *Cmd) Bool(longName, shortName string, defValue bool, usage string) *BoolFlag
```
Bool 添加布尔类型标志, 返回标志对象指针

**参数值:**
  - longName: string - 长标志名
  - shortName: string - 短标志
  - defValue: bool - 默认值
  - usage: string - 帮助说明

**返回值:**
  - *BoolFlag - 布尔标志对象指针

### func (c *Cmd) BoolVar(f *BoolFlag, longName, shortName string, defValue bool, usage string)
```go
func (c *Cmd) BoolVar(f *BoolFlag, longName, shortName string, defValue bool, usage string)
```
BoolVar 绑定布尔类型标志到指针并内部注册Flag对象

**参数值:**
  - f: *BoolFlag - 布尔标志对象指针
  - longName: string - 长标志名
  - shortName: string - 短标志
  - defValue: bool - 默认值
  - usage: string - 帮助说明

### func (c *Cmd) Chinese() bool
```go
func (c *Cmd) Chinese() bool
```
Chinese 获取是否使用中文帮助信息

**返回值:**
  - bool: 是否使用中文帮助信息

### func (c *Cmd) HasSubCmd(cmdName string) bool
```go
func (c *Cmd) HasSubCmd(cmdName string) bool
```
HasSubCmd 检查子命令是否存在

**参数:**
  - cmdName: 子命令名称

**返回:**
  - bool: 子命令是否存在

### func (c *Cmd) Desc() string
```go
func (c *Cmd) Desc() string
```
Desc 返回命令描述

**返回值:**
  - string: 命令描述

### func (c *Cmd) Duration(longName, shortName string, defValue time.Duration, usage string) *DurationFlag
```go
func (c *Cmd) Duration(longName, shortName string, defValue time.Duration, usage string) *DurationFlag
```
Duration 添加时间间隔类型标志, 返回标志对象指针

**参数值:**
  - longName: string - 长标志名
  - shortName: string - 短标志
  - defValue: time.Duration - 默认值
  - usage: string - 帮助说明

**返回值:**
  - *DurationFlag - 时间间隔标志对象指针

### func (c *Cmd) DurationVar(f *DurationFlag, longName, shortName string, defValue time.Duration, usage string)
```go
func (c *Cmd) DurationVar(f *DurationFlag, longName, shortName string, defValue time.Duration, usage string)
```
DurationVar 绑定时间间隔类型标志到指针并内部注册Flag对象

**参数值:**
  - f: *DurationFlag - 时间间隔标志对象指针
  - longName: string - 长标志名
  - shortName: string - 短标志
  - defValue: time.Duration - 默认值
  - usage: string - 帮助说明

### func (c *Cmd) Enum(longName, shortName string, defValue string, usage string, options []string) *EnumFlag
```go
func (c *Cmd) Enum(longName, shortName string, defValue string, usage string, options []string) *EnumFlag
```
Enum 添加枚举类型标志, 返回标志对象指针

**参数值:**
  - longName: string - 长标志名
  - shortName: string - 短标志
  - defValue: string - 默认值
  - usage: string - 帮助说明
  - options: []string - 限制该标志取值的枚举值切片

**返回值:**
  - *EnumFlag - 枚举标志对象指针

### func (c *Cmd) EnumVar(f *EnumFlag, longName, shortName string, defValue string, usage string, options []string)
```go
func (c *Cmd) EnumVar(f *EnumFlag, longName, shortName string, defValue string, usage string, options []string)
```
EnumVar 绑定枚举类型标志到指针并内部注册Flag对象

**参数值:**
  - f: *EnumFlag - 枚举标志对象指针
  - longName: string - 长标志名
  - shortName: string - 短标志
  - defValue: string - 默认值
  - usage: string - 帮助说明
  - options: []string - 限制该标志取值的枚举值切片

### func (c *Cmd) Examples() []ExampleInfo
```go
func (c *Cmd) Examples() []ExampleInfo
```
Examples 获取所有使用示例

**返回:**
  - []ExampleInfo: 使用示例列表

### func (c *Cmd) FlagExists(name string) bool
```go
func (c *Cmd) FlagExists(name string) bool
```
FlagExists 检查指定名称的标志是否存在

**参数:**
  - name: 标志名称

**返回值:**
  - bool: 标志是否存在

### func (c *Cmd) FlagRegistry() *FlagRegistry
```go
func (c *Cmd) FlagRegistry() *FlagRegistry
```
FlagRegistry 获取标志注册表的只读访问

**返回值:** - *FlagRegistry: 标志注册表的只读访问

### func (c *Cmd) GetSubCmd(name string) *Cmd
```go
func (c *Cmd) GetSubCmd(name string) *Cmd
```
GetSubCmd 根据名称获取子命令实例

**参数:**
  - name: 子命令名称 (长名称或短名称)

**返回值:**
  - *Cmd: 子命令实例

**恐慌:**
  - 当name为空字符串时，会抛出"subcommand name cannot be empty"的恐慌
  - 当找不到指定名称的子命令时，会抛出"fmt.Sprintf("subcommand '%s' not found", name)"的恐慌

**并发安全:** 此方法使用读锁保护，可安全地在多个 goroutine 中并发调用。

### func (c *Cmd) Float64(longName, shortName string, defValue float64, usage string) *Float64Flag
```go
func (c *Cmd) Float64(longName, shortName string, defValue float64, usage string) *Float64Flag
```
Float64 添加浮点型标志, 返回标志对象指针

**参数值:**
  - longName - 长标志名
  - shortName - 短标志
  - defValue - 默认值
  - usage - 帮助说明

**返回值:**
  - *Float64Flag - 浮点型标志对象指针

### func (c *Cmd) Float64Var(f *Float64Flag, longName, shortName string, defValue float64, usage string)
```go
func (c *Cmd) Float64Var(f *Float64Flag, longName, shortName string, defValue float64, usage string)
```
Float64Var 绑定浮点型标志到指针并内部注册Flag对象

**参数值:**
  - f: *Float64Flag - 浮点型标志对象指针
  - longName: string - 长标志名
  - shortName: string - 短标志
  - defValue: float64 - 默认值
  - usage: string - 帮助说明

### func (c *Cmd) Help() string
```go
func (c *Cmd) Help() string
```
Help 返回命令用法帮助信息

**返回值:**
  - string: 命令用法帮助信息

### func (c *Cmd) Int(longName, shortName string, defValue int, usage string) *IntFlag
```go
func (c *Cmd) Int(longName, shortName string, defValue int, usage string) *IntFlag
```
Int 添加整数类型标志, 返回标志对象指针

**参数值:**
  - longName: 长标志名
  - shortName: 短标志名
  - defValue: 默认值
  - usage: 帮助说明

**返回值:**
  - *IntFlag: 整数标志对象指针

### func (c *Cmd) Int64(longName, shortName string, defValue int64, usage string) *Int64Flag
```go
func (c *Cmd) Int64(longName, shortName string, defValue int64, usage string) *Int64Flag
```
Int64 添加64位整数类型标志, 返回标志对象指针

**参数值:**
  - longName: 长标志名
  - shortName: 短标志名
  - defValue: 默认值
  - usage: 帮助说明

**返回值:**
  - *Int64Flag: 64位整数标志对象指针

### func (c *Cmd) Int64Slice(longName, shortName string, defValue []int64, usage string) *Int64SliceFlag
```go
func (c *Cmd) Int64Slice(longName, shortName string, defValue []int64, usage string) *Int64SliceFlag
```
Int64Slice 绑定64位整数切片类型标志并内部注册Flag对象

**参数值:**
  - longName: 长标志名
  - shortName: 短标志名
  - defValue: 默认值
  - usage: 帮助说明

**返回值:**
  - *Int64SliceFlag: 64位整数切片标志对象指针

### func (c *Cmd) Int64SliceVar(f *Int64SliceFlag, longName, shortName string, defValue []int64, usage string)
```go
func (c *Cmd) Int64SliceVar(f *Int64SliceFlag, longName, shortName string, defValue []int64, usage string)
```
Int64SliceVar 绑定64位整数切片类型标志到指针并内部注册Flag对象

**参数值:**
  - f: 64位整数切片标志指针
  - longName: 长标志名
  - shortName: 短标志名
  - defValue: 默认值
  - usage: 帮助说明

### func (c *Cmd) Int64Var(f *Int64Flag, longName, shortName string, defValue int64, usage string)
```go
func (c *Cmd) Int64Var(f *Int64Flag, longName, shortName string, defValue int64, usage string)
```
Int64Var 绑定64位整数类型标志到指针并内部注册Flag对象

**参数值:**
  - f: 64位整数标志指针
  - longName: 长标志名
  - shortName: 短标志名
  - defValue: 默认值
  - usage: 帮助说明

### func (c *Cmd) IntSlice(longName, shortName string, defValue []int, usage string) *IntSliceFlag
```go
func (c *Cmd) IntSlice(longName, shortName string, defValue []int, usage string) *IntSliceFlag
```
IntSlice 绑定整数切片类型标志并内部注册Flag对象

**参数值:**
  - longName: 长标志名
  - shortName: 短标志名
  - defValue: 默认值
  - usage: 帮助说明

**返回值:**
  - *IntSliceFlag: 整数切片标志对象指针

### func (c *Cmd) IntSliceVar(f *IntSliceFlag, longName, shortName string, defValue []int, usage string)
```go
func (c *Cmd) IntSliceVar(f *IntSliceFlag, longName, shortName string, defValue []int, usage string)
```
IntSliceVar 绑定整数切片类型标志到指针并内部注册Flag对象

**参数值:**
  - f: 整数切片标志指针
  - longName: 长标志名
  - shortName: 短标志名
  - defValue: 默认值
  - usage: 帮助说明

### func (c *Cmd) IntVar(f *IntFlag, longName, shortName string, defValue int, usage string)
```go
func (c *Cmd) IntVar(f *IntFlag, longName, shortName string, defValue int, usage string)
```
IntVar 绑定整数类型标志到指针并内部注册Flag对象

**参数值:**
  - f: 整数标志指针
  - longName: 长标志名
  - shortName: 短标志名
  - defValue: 默认值
  - usage: 帮助说明

### func (c *Cmd) IsParsed() bool
```go
func (c *Cmd) IsParsed() bool
```
IsParsed 检查命令是否已完成解析

**返回值:**
  - bool: 解析状态,true表示已解析(无论成功失败), false表示未解析

### func (c *Cmd) Logo() string
```go
func (c *Cmd) Logo() string
```
Logo 获取logo文本

**返回值:**
  - string: logo文本字符串

### func (c *Cmd) LongName() string
```go
func (c *Cmd) LongName() string
```
LongName 返回命令长名称

**返回值:**
  - string: 命令长名称

### func (c *Cmd) Map(longName, shortName string, defValue map[string]string, usage string) *MapFlag
```go
func (c *Cmd) Map(longName, shortName string, defValue map[string]string, usage string) *MapFlag
```
Map 添加键值对类型标志, 返回标志对象指针

**参数值:**
  - longName: 长标志名
  - shortName: 短标志名
  - defValue: 默认值
  - usage: 帮助说明

**返回值:**
  - *MapFlag: 键值对标志对象指针

### func (c *Cmd) MapVar(f *MapFlag, longName, shortName string, defValue map[string]string, usage string)
```go
func (c *Cmd) MapVar(f *MapFlag, longName, shortName string, defValue map[string]string, usage string)
```
MapVar 绑定键值对类型标志到指针并内部注册Flag对象

**参数值:**
  - f: 键值对标志指针
  - longName: 长标志名
  - shortName: 短标志名
  - defValue: 默认值
  - usage: 帮助说明

### func (c *Cmd) Modules() string
```go
func (c *Cmd) Modules() string
```
Modules 获取自定义模块帮助信息

**返回值:**
  - string: 自定义模块帮助信息

### func (c *Cmd) NArg() int
```go
func (c *Cmd) NArg() int
```
NArg 获取非标志参数的数量

**返回值:**
  - int: 参数数量

### func (c *Cmd) NFlag() int
```go
func (c *Cmd) NFlag() int
```
NFlag 获取标志的数量

**返回值:**
  - int: 标志数量

### func (c *Cmd) Name() string
```go
func (c *Cmd) Name() string
```
Name 获取命令名称

**返回值:**
  - string: 命令名称

**说明:**
  - 优先返回长名称, 如果长名称不存在则返回短名称

### func (c *Cmd) Notes() []string
```go
func (c *Cmd) Notes() []string
```
Notes 获取所有备注信息

**返回:**
  - 备注信息列表

### func (c *Cmd) Parse(args []string) (err error)
```go
func (c *Cmd) Parse(args []string) (err error)
```
Parse 完整解析命令行参数(含子命令处理)

**主要功能：**
 1. 解析当前命令的长短标志及内置标志
 2. 自动检测并解析子命令及其参数(若存在)
 3. 验证枚举类型标志的有效性

**参数：**
  - args: 原始命令行参数切片(包含可能的子命令及参数)

**返回值：**
  - error: 解析过程中遇到的错误(如标志格式错误、子命令解析失败等)

**注意事项：**
  - 每个Cmd实例仅会被解析一次(线程安全)
  - 若检测到子命令, 会将剩余参数传递给子命令的Parse方法
  - 处理内置标志执行逻辑
  - **重要**: Parse方法只负责解析参数，不会自动执行Run函数。执行需要手动调用

### func (c *Cmd) 手动执行Run函数
```go
// 示例：手动执行Run函数的标准流程
if cmd.Run != nil {
    if err := cmd.Run(cmd); err != nil {
        // 处理执行错误
        fmt.Printf("执行失败: %v\n", err)
        os.Exit(1)
    }
}
```

**执行模式说明：**
- qflag采用分离式架构：解析(Parse)与执行(Run)完全分离
- Parse方法只解析命令行参数，不会自动调用Run函数
- 用户需要手动检查并执行Run函数，提供完全的控制权
- Run函数接收命令实例本身作为参数，可以访问所有标志值

**最佳实践：**
```go
cmd := NewCmd("server", "s", ExitOnError)
port := cmd.Int("port", "p", 8080, "服务端口")

// 设置执行函数
cmd.Run = func(c *Cmd) error {
    fmt.Printf("启动服务器，端口: %d\n", port.Get())
    // 实际的服务器启动逻辑
    return nil
}

// 解析参数
if err := cmd.Parse(os.Args[1:]); err != nil {
    log.Fatal(err)
}

// 手动执行（重要：不会自动执行！）
if cmd.Run != nil {
    if err := cmd.Run(cmd); err != nil {
        log.Fatal(err)
    }
}
```

### func (c *Cmd) ParseFlagsOnly(args []string) (err error)
```go
func (c *Cmd) ParseFlagsOnly(args []string) (err error)
```
ParseFlagsOnly 仅解析当前命令的标志参数(忽略子命令)

**主要功能：**
 1. 解析当前命令的长短标志及内置标志
 2. 验证枚举类型标志的有效性
 3. 明确忽略所有子命令及后续参数

**参数：**
  - args: 原始命令行参数切片(子命令及后续参数会被忽略)

**返回值：**
  - error: 解析过程中遇到的错误(如标志格式错误等)

**注意事项：**
  - 每个Cmd实例仅会被解析一次(线程安全)
  - 不会处理任何子命令, 所有参数均视为当前命令的标志或位置参数
  - 处理内置标志逻辑

### func (c *Cmd) ParseAndRoute(args []string) error
```go
func (c *Cmd) ParseAndRoute(args []string) error
```
ParseAndRoute 解析参数并自动路由执行子命令

**主要功能：**
 1. 解析当前命令的参数
 2. 自动路由并执行匹配的子命令
 3. 如果没有匹配的子命令，执行当前命令或显示帮助

**参数：**
  - args: 命令行参数列表(通常为 os.Args[1:])

**返回值：**
  - error: 执行过程中遇到的错误

**使用示例：**
```go
cmd := qflag.NewCmd("myapp", "", qflag.ExitOnError)
// 添加子命令...
subCmd := qflag.NewCmd("sub", "s", qflag.ExitOnError)
cmd.AddSubCmd(subCmd)

// 一行代码完成解析和路由执行
if err := cmd.ParseAndRoute(os.Args[1:]); err != nil {
    log.Fatal(err)
}
```

### func (c *Cmd) PrintHelp()
```go
func (c *Cmd) PrintHelp()
```
PrintHelp 打印命令的帮助信息, 优先打印用户的帮助信息, 否则自动生成帮助信息

**注意:**
  - 打印帮助信息时, 不会自动退出程序

### func (c *Cmd) SetChinese(useChinese bool)
```go
func (c *Cmd) SetChinese(useChinese bool)
```
SetChinese 设置是否使用中文帮助信息

**参数:**
  - useChinese: 是否使用中文帮助信息

### func (c *Cmd) SetCompletion(enable bool)
```go
func (c *Cmd) SetCompletion(enable bool)
```
SetCompletion 设置是否启用自动补全, 只能在根命令上启用

**参数:**
  - enable: true表示启用补全,false表示禁用

### func (c *Cmd) SetRun(run func(*Cmd) error)
```go
func (c *Cmd) SetRun(run func(*Cmd) error)
```
SetRun 设置命令的执行函数

**参数:**
  - run: 命令执行函数，接收*Cmd作为参数，返回error

**恐慌:**
  - 当run为nil时，会抛出"run function cannot be nil"的恐慌

### func (c *Cmd) SetDesc(desc string)
```go
func (c *Cmd) SetDesc(desc string)
```
SetDesc 设置命令描述

**参数:**
  - desc: 命令描述

### func (c *Cmd) SetHelp(help string)
```go
func (c *Cmd) SetHelp(help string)
```
SetHelp 设置用户自定义命令帮助信息

**参数:**
  - help: 用户自定义命令帮助信息

### func (c *Cmd) SetLogo(logoText string)
```go
func (c *Cmd) SetLogo(logoText string)
```
SetLogo 设置logo文本

**参数:**
  - logoText: logo文本字符串

### func (c *Cmd) SetModules(moduleHelps string)
```go
func (c *Cmd) SetModules(moduleHelps string)
```
SetModules 设置自定义模块帮助信息

**参数:**
  - moduleHelps: 自定义模块帮助信息

### func (c *Cmd) SetNoFgExit(exit bool)
```go
func (c *Cmd) SetNoFgExit(exit bool)
```
SetNoFgExit 设置禁用内置标志自动退出 默认情况下为false, 当解析到内置参数时, QFlag将退出程序

**参数:**
  - exit: 是否退出

### func (c *Cmd) SetUsage(usageSyntax string)
```go
func (c *Cmd) SetUsage(usageSyntax string)
```
SetUsage 设置自定义命令用法

**参数:**
  - usageSyntax: 自定义命令用法

### func (c *Cmd) SetVersion(version string)
```go
func (c *Cmd) SetVersion(version string)
```
SetVersion 设置版本信息

**参数:**
  - version: 版本信息

### func (c *Cmd) SetVersionf(format string, args ...any)
```go
func (c *Cmd) SetVersionf(format string, args ...any)
```
SetVersionf 设置版本信息

**参数:**
  - format: 版本信息格式字符串
  - args: 格式化参数

### func (c *Cmd) ShortName() string
```go
func (c *Cmd) ShortName() string
```
ShortName 返回命令短名称

**返回值:**
  - string: 命令短名称

### func (c *Cmd) Size(longName, shortName string, defValue int64, usage string) *SizeFlag
```go
func (c *Cmd) Size(longName, shortName string, defValue int64, usage string) *SizeFlag
```
Size 添加大小类型标志, 返回标志对象指针

**参数值:**
  - longName: string - 长标志名
  - shortName: string - 短标志名
  - defValue: int64 - 默认值(单位为字节)
  - usage: string - 帮助说明

**返回值:**
  - *SizeFlag - 大小标志对象指针

**支持的单位格式:**
  - 字节: "B", "b", "byte", "bytes"
  - 十进制: "KB", "MB", "GB", "TB", "PB" 或简写 "K", "M", "G", "T", "P"
  - 二进制: "KiB", "MiB", "GiB", "TiB", "PiB"
  - 支持小数: "1.5GB", "2.5MB"
  - 支持负数: "-1GB", "-500MB"
  - 特殊值: "0" (零值特例)

### func (c *Cmd) SizeVar(f *SizeFlag, longName, shortName string, defValue int64, usage string)
```go
func (c *Cmd) SizeVar(f *SizeFlag, longName, shortName string, defValue int64, usage string)
```
SizeVar 绑定大小类型标志到指针并内部注册Flag对象

**参数值:**
  - f: *SizeFlag - 大小标志对象指针
  - longName: string - 长标志名
  - shortName: string - 短标志名
  - defValue: int64 - 默认值(单位为字节)
  - usage: string - 帮助说明

**支持的单位格式:**
  - 字节: "B", "b", "byte", "bytes"
  - 十进制: "KB", "MB", "GB", "TB", "PB" 或简写 "K", "M", "G", "T", "P"
  - 二进制: "KiB", "MiB", "GiB", "TiB", "PiB"
  - 支持小数: "1.5GB", "2.5MB"
  - 支持负数: "-1GB", "-500MB"
  - 特殊值: "0" (零值特例)

### func (c *Cmd) String(longName, shortName, defValue, usage string) *StringFlag
```go
func (c *Cmd) String(longName, shortName, defValue, usage string) *StringFlag
```
String 添加字符串类型标志, 返回标志对象指针

**参数值:**
  - longName: 长标志名
  - shortName: 短标志名
  - defValue: 默认值
  - usage: 帮助说明

**返回值:**
  - *StringFlag: 字符串标志对象指针

### func (c *Cmd) StringSlice(longName, shortName string, defValue []string, usage string) *StringSliceFlag
```go
func (c *Cmd) StringSlice(longName, shortName string, defValue []string, usage string) *StringSliceFlag
```
StringSlice 绑定字符串切片类型标志并内部注册Flag对象

**参数值:**
  - longName: 长标志名
  - shortName: 短标志名
  - defValue: 默认值
  - usage: 帮助说明

**返回值:**
  - *StringSliceFlag: 字符串切片标志对象指针

### func (c *Cmd) StringSliceVar(f *StringSliceFlag, longName, shortName string, defValue []string, usage string)
```go
func (c *Cmd) StringSliceVar(f *StringSliceFlag, longName, shortName string, defValue []string, usage string)
```
StringSliceVar 绑定字符串切片类型标志到指针并内部注册Flag对象

**参数值:**
  - f: 字符串切片标志指针
  - longName: 长标志名
  - shortName: 短标志名
  - defValue: 默认值
  - usage: 帮助说明

### func (c *Cmd) StringVar(f *StringFlag, longName, shortName, defValue, usage string)
```go
func (c *Cmd) StringVar(f *StringFlag, longName, shortName, defValue, usage string)
```
StringVar 绑定字符串类型标志到指针并内部注册Flag对象

**参数值:**
  - f: 字符串标志指针
  - longName: 长标志名
  - shortName: 短标志名
  - defValue: 默认值
  - usage: 帮助说明

### func (c *Cmd) SubCmdMap() map[string]*Cmd
```go
func (c *Cmd) SubCmdMap() map[string]*Cmd
```
SubCmdMap 返回子命令映射表(长命令名+短命令名)

**返回值:**
  - map[string]*Cmd: 子命令映射表

### func (c *Cmd) Time(longName, shortName string, defValue string, usage string) *TimeFlag
```go
func (c *Cmd) Time(longName, shortName string, defValue string, usage string) *TimeFlag
```
Time 添加时间类型标志, 返回标志对象指针

**参数值:**
  - longName: 长标志名
  - shortName: 短标志名
  - defValue: 默认值(时间表达式, 如"now", "zero", "1h", "2006-01-02")
  - usage: 帮助说明

**返回值:**
  - *TimeFlag: 时间标志对象指针

**支持的默认值格式:**
  - "now" 或 "" : 当前时间
  - "zero" : 零时间 (time.Time{})
  - "1h", "30m", "-2h" : 相对时间（基于当前时间的偏移）
  - "2006-01-02", "2006-01-02 15:04:05" : 绝对时间格式
  - RFC3339等标准格式

### func (c *Cmd) TimeVar(f *TimeFlag, longName, shortName string, defValue string, usage string)
```go
func (c *Cmd) TimeVar(f *TimeFlag, longName, shortName string, defValue string, usage string)
```
TimeVar 绑定时间类型标志到指针并内部注册Flag对象

**参数值:**
  - f: 时间标志指针
  - longName: 长标志名
  - shortName: 短标志名
  - defValue: 默认值(时间表达式, 如"now", "zero", "1h", "2006-01-02")
  - usage: 帮助说明

**支持的默认值格式:**
  - "now" 或 "" : 当前时间
  - "zero" : 零时间 (time.Time{})
  - "1h", "30m", "-2h" : 相对时间（基于当前时间的偏移）
  - "2006-01-02", "2006-01-02 15:04:05" : 绝对时间格式
  - RFC3339等标准格式

### func (c *Cmd) Uint16(longName, shortName string, defValue uint16, usage string) *Uint16Flag
```go
func (c *Cmd) Uint16(longName, shortName string, defValue uint16, usage string) *Uint16Flag
```
Uint16 添加16位无符号整数类型标志, 返回标志对象指针

**参数值:**
  - longName: 长标志名
  - shortName: 短标志名
  - defValue: 默认值
  - usage: 帮助说明

**返回值:**
  - *Uint16Flag: 16位无符号整数标志对象指针

### func (c *Cmd) Uint16Var(f *Uint16Flag, longName, shortName string, defValue uint16, usage string)
```go
func (c *Cmd) Uint16Var(f *Uint16Flag, longName, shortName string, defValue uint16, usage string)
```
Uint16Var 绑定16位无符号整数类型标志到指针并内部注册Flag对象

**参数值:**
  - f: 16位无符号整数标志指针
  - longName: 长标志名
  - shortName: 短标志名
  - defValue: 默认值
  - usage: 帮助说明

### func (c *Cmd) Uint32(longName, shortName string, defValue uint32, usage string) *Uint32Flag
```go
func (c *Cmd) Uint32(longName, shortName string, defValue uint32, usage string) *Uint32Flag
```
Uint32 添加32位无符号整数类型标志, 返回标志对象指针

**参数值:**
  - longName: 长标志名
  - shortName: 短标志名
  - defValue: 默认值
  - usage: 帮助说明

**返回值:**
  - *Uint32Flag: 32位无符号整数标志对象指针

### func (c *Cmd) Uint32Var(f *Uint32Flag, longName, shortName string, defValue uint32, usage string)
```go
func (c *Cmd) Uint32Var(f *Uint32Flag, longName, shortName string, defValue uint32, usage string)
```
Uint32Var 绑定32位无符号整数类型标志到指针并内部注册Flag对象

**参数值:**
  - f: 32位无符号整数标志指针
  - longName: 长标志名
  - shortName: 短标志名
  - defValue: 默认值
  - usage: 帮助说明

### func (c *Cmd) Uint64(longName, shortName string, defValue uint64, usage string) *Uint64Flag
```go
func (c *Cmd) Uint64(longName, shortName string, defValue uint64, usage string) *Uint64Flag
```
Uint64 添加64位无符号整数类型标志, 返回标志对象指针

**参数值:**
  - longName: 长标志名
  - shortName: 短标志名
  - defValue: 默认值
  - usage: 帮助说明

**返回值:**
  - *Uint64Flag: 64位无符号整数标志对象指针

### func (c *Cmd) Uint64Var(f *Uint64Flag, longName, shortName string, defValue uint64, usage string)
```go
func (c *Cmd) Uint64Var(f *Uint64Flag, longName, shortName string, defValue uint64, usage string)
```
Uint64Var 绑定64位无符号整数类型标志到指针并内部注册Flag对象

**参数值:**
  - f: 64位无符号整数标志指针
  - longName: 长标志名
  - shortName: 短标志名
  - defValue: 默认值
  - usage: 帮助说明

### func (c *Cmd) Usage() string
```go
func (c *Cmd) Usage() string
```
Usage 获取自定义命令用法

**返回值:**
  - string: 自定义命令用法

### func (c *Cmd) Version() string
```go
func (c *Cmd) Version() string
```
Version 获取版本信息

**返回值:** - string: 版本信息

## 类型别名

### type CmdConfig = types.CmdConfig
```go
type CmdConfig = types.CmdConfig
```
CmdConfig 导出cmd包中的CmdConfig结构体

### type DurationFlag = flags.DurationFlag
```go
type DurationFlag = flags.DurationFlag
```
DurationFlag 导出flag包中的DurationFlag结构体

### type EnumFlag = flags.EnumFlag
```go
type EnumFlag = flags.EnumFlag
```
EnumFlag 导出flag包中的EnumFlag结构体

### type ErrorHandling = flags.ErrorHandling
```go
type ErrorHandling = flags.ErrorHandling
```
ErrorHandling 错误处理策略

## 变量

```go
var (
	// ContinueOnError 解析错误时继续解析并返回错误
	ContinueOnError ErrorHandling = flags.ContinueOnError
	// ExitOnError 解析错误时退出程序
	ExitOnError ErrorHandling = flags.ExitOnError
	// PanicOnError 解析错误时触发panic
	PanicOnError ErrorHandling = flags.PanicOnError
)
```
ErrorHandling 错误处理策略常量

## 类型别名

### type ExampleInfo = types.ExampleInfo
```go
type ExampleInfo = types.ExampleInfo
```
ExampleInfo 导出示例信息类型

### type Flag = flags.Flag
```go
type Flag = flags.Flag
```
Flag 导出flag包中的Flag结构体

### type FlagRegistry = flags.FlagRegistry
```go
type FlagRegistry = flags.FlagRegistry
```
FlagRegistry 导出flag包中的FlagRegistry结构体

### type Float64Flag = flags.Float64Flag
```go
type Float64Flag = flags.Float64Flag
```
Float64Flag 导出flag包中的Float64Flag结构体

### type Int64Flag = flags.Int64Flag
```go
type Int64Flag = flags.Int64Flag
```
Int64Flag 导出flag包中的Int64Flag结构体

### type Int64SliceFlag = flags.Int64SliceFlag
```go
type Int64SliceFlag = flags.Int64SliceFlag
```
Int64SliceFlag 导出flag包中的Int64SliceFlag结构体

### type IntFlag = flags.IntFlag
```go
type IntFlag = flags.IntFlag
```
IntFlag 导出flag包中的IntFlag结构体

### type IntSliceFlag = flags.IntSliceFlag
```go
type IntSliceFlag = flags.IntSliceFlag
```
IntSliceFlag 导出flag包中的IntSliceFlag结构体

### type MapFlag = flags.MapFlag
```go
type MapFlag = flags.MapFlag
```
MapFlag 导出flag包中的MapFlag结构体

### type SizeFlag = flags.SizeFlag
```go
type SizeFlag = flags.SizeFlag
```
SizeFlag 导出flag包中的SizeFlag结构体

### type StringFlag = flags.StringFlag
```go
type StringFlag = flags.StringFlag
```
StringFlag 导出flag包中的StringFlag结构体

### type StringSliceFlag = flags.StringSliceFlag
```go
type StringSliceFlag = flags.StringSliceFlag
```
StringSliceFlag 导出flag包中的StringSliceFlag结构体

### type TimeFlag = flags.TimeFlag
```go
type TimeFlag = flags.TimeFlag
```
TimeFlag 导出flag包中的TimeFlag结构体

### type Uint16Flag = flags.Uint16Flag
```go
type Uint16Flag = flags.Uint16Flag
```
Uint16Flag 导出flag包中的UintFlag结构体

### type Uint32Flag = flags.Uint32Flag
```go
type Uint32Flag = flags.Uint32Flag
```
Uint32Flag 导出flag包中的Uint32Flag结构体

### type Uint64Flag = flags.Uint64Flag
```go
type Uint64Flag = flags.Uint64Flag
```
Uint64Flag 导出flag包中的Uint64Flag结构体