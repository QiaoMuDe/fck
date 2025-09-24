# Package qflag

**Import Path:** `gitee.com/MM-Q/qflag`

Package qflag 根包统一导出入口。本文件用于将各子包的核心功能导出到根包，简化外部使用。通过类型别名和变量导出的方式，为用户提供统一的API接口。

Package qflag 提供对标准库flag的封装，自动实现长短标志，并默认绑定-h/--help标志打印帮助信息。用户可通过Cmd.Help字段自定义帮助内容，支持直接赋值字符串或从文件加载。该包是一个功能强大的命令行参数解析库，支持子命令、多种数据类型标志、参数验证等高级特性。

Package qflag 全局标志函数定义文件。本文件提供了全局默认命令实例的各种标志创建和绑定函数，包括字符串、整数、布尔、浮点数、枚举、时间间隔、切片、时间、映射等类型的标志支持。

## Variables

```go
var (
    // qCommandLine 全局默认Command实例（保持原名，与标准库flag对齐）
    qCommandLine *Cmd
)
```

## Functions

### AddExample

```go
func AddExample(desc, usage string)
```

AddExample 添加示例。该函数用于添加命令行标志的示例，这些示例将在命令行帮助信息中显示。

**参数:**
- `desc`: 示例描述，字符串类型
- `usage`: 示例用法，字符串类型

### AddExamples

```go
func AddExamples(examples []cmd.ExampleInfo)
```

AddExamples 添加示例。该函数用于添加命令行标志的示例，这些示例将在命令行帮助信息中显示。

**参数:**
- `examples`: 示例列表，每个元素为 ExampleInfo 类型

### AddNote

```go
func AddNote(note string)
```

AddNote 添加注意事项。该函数用于添加命令行标志的注意事项，这些注意事项将在命令行帮助信息中显示。

**参数:**
- `note`: 注意事项内容，字符串类型

### AddNotes

```go
func AddNotes(notes []string)
```

AddNotes 添加注意事项。该函数用于添加命令行标志的注意事项，这些注意事项将在命令行帮助信息中显示。

**参数:**
- `notes`: 注意事项内容，字符串切片，每个元素为一个注意事项

### AddSubCmd

```go
func AddSubCmd(subCmds ...*Cmd) error
```

AddSubCmd 向全局默认命令实例 `QCommandLine` 添加一个或多个子命令。该函数会调用全局默认命令实例的 `AddSubCmd` 方法，支持批量添加子命令。在添加过程中，会检查子命令是否为 `nil` 以及是否存在循环引用，若有异常则返回错误信息。

**参数:**
- `subCmds`: 可变参数，接收一个或多个 `*Cmd` 类型的子命令实例

**返回值:**
- `error`: 若添加子命令过程中出现错误（如子命令为 `nil` 或存在循环引用），则返回错误信息；否则返回 `nil`

### AddSubCmds

```go
func AddSubCmds(subCmds []*Cmd) error
```

AddSubCmds 向全局默认命令实例 `QCommandLine` 添加子命令切片的便捷函数。该函数是 AddSubCmd 的便捷包装，专门用于处理子命令切片，内部直接调用全局默认命令实例的 `AddSubCmds` 方法，具有相同的验证逻辑和并发安全特性。

**参数:**
- `subCmds`: 子命令切片，包含要添加的所有子命令实例指针

**返回值:**
- `error`: 添加过程中的错误信息，与 AddSubCmd 返回的错误类型相同

**使用场景对比:**
```go
// 使用 AddSubCmd - 适合已知数量的子命令
qflag.AddSubCmd(subCmd1, subCmd2, subCmd3)

// 使用 AddSubCmds - 适合动态生成的子命令切片
subCmds := []*qflag.Cmd{subCmd1, subCmd2, subCmd3}
qflag.AddSubCmds(subCmds)
```

### Arg

```go
func Arg(i int) string
```

Arg 获取全局默认命令实例 `QCommandLine` 解析后的指定索引位置的非标志参数。索引从 0 开始，若索引超出非标志参数切片的范围，将返回空字符串。

**参数:**
- `i`: 非标志参数的索引位置，从 0 开始计数

**返回值:**
- `string`: 指定索引位置的非标志参数；若索引越界，则返回空字符串

### Args

```go
func Args() []string
```

Args 获取全局默认命令实例 `QCommandLine` 解析后的非标志参数切片。非标志参数是指命令行中未被识别为标志的参数。

**返回值:**
- `[]string`: 包含所有非标志参数的字符串切片

### Bool

```go
func Bool(longName, shortName string, defValue bool, usage string) *flags.BoolFlag
```

Bool 为全局默认命令创建一个布尔类型的命令行标志。该函数会调用全局默认命令实例的 Bool 方法，为命令行添加一个支持长短标志的布尔参数。

**参数:**
- `longName`: 标志的长名称，在命令行中以 --name 的形式使用
- `shortName`: 标志的短名称，在命令行中以 -shortName 的形式使用
- `defValue`: 标志的默认值，当命令行未指定该标志时使用
- `usage`: 标志的帮助说明信息，用于在显示帮助信息时展示

**返回值:**
- `*flags.BoolFlag`: 指向新创建的布尔标志对象的指针

### BoolVar

```go
func BoolVar(f *flags.BoolFlag, longName, shortName string, defValue bool, usage string)
```

BoolVar 函数的作用是将布尔类型的命令行标志绑定到全局默认命令实例 `QCommandLine` 中。它会调用全局默认命令实例的 `BoolVar` 方法，为命令行添加一个支持长短和短标志的布尔参数，并将该参数与传入的 `BoolFlag` 指针建立关联，后续可以通过该指针获取和使用该标志的值。

**参数:**
- `f`: 指向 `BoolFlag` 类型的指针，用于存储和管理布尔类型命令行标志的相关信息，如当前值、默认值等
- `longName`: 标志的长名称，在命令行中以 `--name` 的形式使用
- `shortName`: 标志的短名称，在命令行中以 `-shortName` 的形式使用
- `defValue`: 标志的默认值，当命令行未指定该标志时，会使用此默认值
- `usage`: 标志的帮助说明信息，用于在显示帮助信息时展示给用户，解释该标志的用途

### CmdExists

```go
func CmdExists(cmdName string) bool
```

CmdExists 检查子命令是否存在。

**参数:**
- `cmdName`: 子命令名称

**返回值:**
- `bool`: 子命令是否存在

### Duration

```go
func Duration(longName, shortName string, defValue time.Duration, usage string) *flags.DurationFlag
```

Duration 为全局默认命令定义一个时间间隔类型的命令行标志。该函数会调用全局默认命令实例 `QCommandLine` 的 `Duration` 方法，为命令行添加支持长短标志的时间间隔类型参数。

**参数:**
- `longName`: 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式
- `shortName`: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式
- `defValue`: 该命令行标志的默认值。当用户在命令行中未指定该标志时，会采用此默认值。该值会被复制一份，避免外部修改影响内部状态
- `usage`: 该命令行标志的帮助说明信息，在显示帮助信息时会呈现给用户，用以解释该标志的具体用途

**返回值:**
- `*flags.DurationFlag`: 指向新创建的时间间隔类型标志对象的指针

### DurationVar

```go
func DurationVar(f *flags.DurationFlag, longName, shortName string, defValue time.Duration, usage string)
```

DurationVar 为全局默认命令将一个时间间隔类型的命令行标志绑定到指定的 `DurationFlag` 指针。该函数会调用全局默认命令实例 `QCommandLine` 的 `DurationVar` 方法，为命令行添加支持长短标志的时间间隔类型参数。

**参数:**
- `f`: 指向 `DurationFlag` 类型的指针，此指针用于存储和管理时间间隔类型命令行标志的各类信息，如当前标志的值、默认值等
- `longName`: 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式
- `shortName`: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式
- `defValue`: 该命令行标志的默认值。当用户在命令行中未指定该标志时，会采用此默认值。该值会被复制一份，避免外部修改影响内部状态
- `usage`: 该命令行标志的帮助说明信息，在显示帮助信息时会呈现给用户，用以解释该标志的具体用途

### Enum

```go
func Enum(longName, shortName string, defValue string, usage string, enumValues []string) *flags.EnumFlag
```

Enum 为全局默认命令定义一个枚举类型的命令行标志。该函数会调用全局默认命令实例 `QCommandLine` 的 `Enum` 方法，为命令行添加支持长短标志的枚举类型参数。

**参数:**
- `longName`: 标志的长名称，在命令行中以 `--name` 的形式使用
- `shortName`: 标志的短名称，在命令行中以 `-shortName` 的形式使用
- `defValue`: 标志的默认值，当命令行未指定该标志时使用
- `usage`: 标志的帮助说明信息，用于在显示帮助信息时展示
- `enumValues`: 枚举值的集合，用于指定标志可接受的取值范围

**返回值:**
- `*flags.EnumFlag`: 指向新创建的枚举类型标志对象的指针

### EnumVar

```go
func EnumVar(f *flags.EnumFlag, longName, shortName string, defValue string, usage string, enumValues []string)
```

EnumVar 为全局默认命令将一个枚举类型的命令行标志绑定到指定的 `EnumFlag` 指针。该函数会调用全局默认命令实例 `QCommandLine` 的 `EnumVar` 方法，为命令行添加支持长短标志的枚举类型参数。

**参数:**
- `f`: 指向 `EnumFlag` 类型的指针，此指针用于存储和管理枚举类型命令行标志的各类信息，如当前标志的值、默认值等
- `longName`: 命令行标志的长名称，在命令行中使用时需遵循 `--name` 的格式
- `shortName`: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式
- `defValue`: 该命令行标志的默认值。当用户在命令行中未指定该标志时，会采用此默认值。该值会被复制一份，避免外部修改影响内部状态
- `usage`: 该命令行标志的帮助说明信息，在显示帮助信息时会呈现给用户，用以解释该标志的具体用途
- `enumValues`: 枚举值的集合，用于指定标志可接受的取值范围

### FlagExists

```go
func FlagExists(name string) bool
```

FlagExists 检查全局默认命令实例 `QCommandLine` 中是否存在指定名称的标志。该函数会调用全局默认命令实例的 `FlagExists` 方法，用于检查命令行中是否存在指定名称的标志。

**参数:**
- `name`: 要检查的标志名称，可以是长名称或短名称

**返回值:**
- `bool`: 若存在指定名称的标志，则返回 `true`；否则返回 `false`

### FlagRegistry

```go
func FlagRegistry() *flags.FlagRegistry
```

FlagRegistry 获取标志注册表。

**返回值:**
- `*flags.FlagRegistry`: 标志注册表

### Float64

```go
func Float64(longName, shortName string, defValue float64, usage string) *flags.Float64Flag
```

Float64 为全局默认命令创建一个浮点数类型的命令行标志。该函数会调用全局默认命令实例的 Float64 方法，为命令行添加一个支持长短标志的浮点数参数。

**参数:**
- `longName`: 标志的长名称，在命令行中以 --name 的形式使用
- `shortName`: 标志的短名称，在命令行中以 -shortName 的形式使用
- `defValue`: 标志的默认值，当命令行未指定该标志时使用
- `usage`: 标志的帮助说明信息，用于在显示帮助信息时展示

**返回值:**
- `*flags.Float64Flag`: 指向新创建的浮点数标志对象的指针

### Float64Var

```go
func Float64Var(f *flags.Float64Flag, longName, shortName string, defValue float64, usage string)
```

Float64Var 为全局默认命令绑定一个浮点数类型的命令行标志到指定的 `Float64Flag` 指针。该函数会调用全局默认命令实例 `QCommandLine` 的 `Float64Var` 方法，为命令行添加支持长短标志的浮点数参数，并将该参数与传入的 `Float64Flag` 指针关联，以便后续获取和使用该标志的值。

**参数:**
- `f`: 指向 `Float64Flag` 的指针，用于存储和管理该浮点数类型命令行标志的相关信息，包括当前值、默认值等
- `longName`: 命令行标志的长名称，在命令行中需以 `--name` 的格式使用
- `shortName`: 命令行标志的短名称，在命令行中需以 `-shortName` 的格式使用
- `defValue`: 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值
- `usage`: 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户，用于解释该标志的用途

### Desc

```go
func Desc() string
```

Desc 获取命令描述信息。

**返回值:**
- `string`: 命令描述信息

### Examples

```go
func Examples() []cmd.ExampleInfo
```

Examples 获取示例信息。该函数用于获取命令行标志的示例信息列表。

**返回值:**
- `[]cmd.ExampleInfo`: 示例信息列表，每个元素为 ExampleInfo 类型

### Help

```go
func Help() string
```

Help 返回全局默认命令实例 `QCommandLine` 的帮助信息。

**返回值:**
- `string`: 命令行帮助信息

### Logo

```go
func Logo() string
```

Logo 获取全局默认命令实例 `QCommandLine` 的 logo 文本。

**返回值:**
- `string`: 配置的 logo 文本

### Modules

```go
func Modules() string
```

Modules 获取模块帮助信息。

**返回值:**
- `string`: 模块帮助信息

### Notes

```go
func Notes() []string
```

Notes 获取所有备注信息。

**返回值:**
- `[]string`: 备注信息列表

### Usage

```go
func Usage() string
```

Usage 获取全局默认命令实例 `QCommandLine` 的用法信息。

**返回值:**
- `string`: 命令行用法信息

### Chinese

```go
func Chinese() bool
```

Chinese 获取是否使用中文。该函数用于获取当前命令行标志是否使用中文。

**返回值:**
- `bool`: 如果使用中文，则返回true；否则返回false

### Version

```go
func Version() string
```

Version 获取全局默认命令的版本信息。

**返回值:**
- `string`: 版本信息字符串

### Int

```go
func Int(longName, shortName string, defValue int, usage string) *flags.IntFlag
```

Int 为全局默认命令创建一个整数类型的命令行标志。该函数会调用全局默认命令实例的 Int 方法，为命令行添加一个支持长短标志的整数参数。

**参数:**
- `longName`: 标志的长名称，在命令行中以 --name 的形式使用
- `shortName`: 标志的短名称，在命令行中以 -shortName 的形式使用
- `defValue`: 标志的默认值，当命令行未指定该标志时使用
- `usage`: 标志的帮助说明信息，用于在显示帮助信息时展示

**返回值:**
- `*flags.IntFlag`: 指向新创建的整数标志对象的指针

### Int64

```go
func Int64(longName, shortName string, defValue int64, usage string) *flags.Int64Flag
```

Int64 为全局默认命令定义一个64位整数类型的命令行标志。该函数会调用全局默认命令实例 `QCommandLine` 的 `Int64` 方法，为命令行添加支持长短标志的64位整数类型参数。

**参数:**
- `longName`: 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式
- `shortName`: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式
- `defValue`: 命令行标志的默认值
- `usage`: 命令行标志的用法说明

**返回值:**
- `*flags.Int64Flag`: 指向新创建的64位整数类型标志对象的指针

### Int64Var

```go
func Int64Var(f *flags.Int64Flag, longName, shortName string, defValue int64, usage string)
```

Int64Var 函数创建一个64位整数类型标志，并将其绑定到指定的 `Int64Flag` 指针。该函数会调用全局默认命令实例 `QCommandLine` 的 `Int64Var` 方法，为命令行添加支持长短标志的64位整数类型参数。

**参数:**
- `f`: 指向要绑定的 `Int64Flag` 对象的指针
- `longName`: 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式
- `shortName`: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式
- `defValue`: 命令行标志的默认值
- `usage`: 命令行标志的用法说明

### IntVar

```go
func IntVar(f *flags.IntFlag, longName, shortName string, defValue int, usage string)
```

IntVar 函数的作用是将整数类型的命令行标志绑定到全局默认命令的 `IntFlag` 指针上。它借助全局默认命令实例 `QCommandLine` 的 `IntVar` 方法，为命令行添加支持长短标志的整数参数，并将该参数与传入的 `IntFlag` 指针建立关联，方便后续对该标志的值进行获取和使用。

**参数:**
- `f`: 指向 `IntFlag` 类型的指针，此指针用于存储和管理整数类型命令行标志的各类信息
- `longName`: 命令行标志的长名称，在命令行中使用时需遵循 `--name` 的格式
- `shortName`: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式
- `defValue`: 该命令行标志的默认值。当用户在命令行中未指定该标志时，会采用此默认值
- `usage`: 该命令行标志的帮助说明信息，在显示帮助信息时会呈现给用户，用以解释该标志的具体用途

### IsParsed

```go
func IsParsed() bool
```

IsParsed 检查命令行参数是否已解析。

**返回值:**
- `bool`: 是否已解析

### LongName

```go
func LongName() string
```

LongName 获取命令长名称。

### Map

```go
func Map(longName, shortName string, defValue map[string]string, usage string) *flags.MapFlag
```

Map 为全局默认命令创建一个键值对类型的命令行标志。该函数会调用全局默认命令实例的 Map 方法，为命令行添加一个支持长短标志的键值对参数。

**参数:**
- `longName`: 标志的长名称，在命令行中以 --longName 的形式使用
- `shortName`: 标志的短名称，在命令行中以 -shortName 的形式使用
- `defValue`: 标志的默认值，当命令行未指定该标志时使用
- `usage`: 标志的帮助说明信息，用于在显示帮助信息时展示

**返回值:**
- `*flags.MapFlag`: 指向新创建的键值对标志对象的指针

### MapVar

```go
func MapVar(f *flags.MapFlag, longName, shortName string, defValue map[string]string, usage string)
```

MapVar 为全局默认命令将一个键值对类型的命令行标志绑定到指定的 MapFlag 指针。该函数会调用全局默认命令实例的 MapVar 方法，为命令行添加支持长短标志的键值对参数，并将该参数与传入的 MapFlag 指针关联，以便后续获取和使用该标志的值。

**参数:**
- `f`: 指向 MapFlag 的指针，用于存储和管理该键值对类型命令行标志的相关信息
- `longName`: 命令行标志的长名称，在命令行中需以 --longName 的格式使用
- `shortName`: 命令行标志的短名称，在命令行中需以 -shortName 的格式使用
- `defValue`: 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值
- `usage`: 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户

### NArg

```go
func NArg() int
```

NArg 获取全局默认命令实例 `QCommandLine` 解析后的非标志参数的数量。

**返回值:**
- `int`: 非标志参数的数量

### NFlag

```go
func NFlag() int
```

NFlag 获取全局默认命令实例 `QCommandLine` 解析后已定义和使用的标志的数量。

**返回值:**
- `int`: 标志的数量

### Name

```go
func Name() string
```

Name 获取全局默认命令实例 `QCommandLine` 的名称。

**返回值:**
- 优先返回长名称，如果长名称不存在则返回短名称

### Parse

```go
func Parse() error
```

Parse 完整解析命令行参数（含子命令处理）。

**主要功能:**
1. 解析当前命令的长短标志及内置标志
2. 自动检测并解析子命令及其参数（若存在）
3. 验证枚举类型标志的有效性

**参数:**
- `args`: 原始命令行参数切片（包含可能的子命令及参数）

**返回值:**
- 解析过程中遇到的错误（如标志格式错误、子命令解析失败等）

**注意事项:**
- 每个Cmd实例仅会被解析一次（线程安全）
- 若检测到子命令，会将剩余参数传递给子命令的Parse方法
- 处理内置标志执行逻辑

### ParseFlagsOnly

```go
func ParseFlagsOnly() error
```

ParseFlagsOnly 仅解析当前命令的标志参数（忽略子命令）。

**主要功能:**
1. 解析当前命令的长短标志及内置标志
2. 验证枚举类型标志的有效性
3. 明确忽略所有子命令及后续参数

**参数:**
- `args`: 原始命令行参数切片（子命令及后续参数会被忽略）

**返回值:**
- 解析过程中遇到的错误（如标志格式错误等）

**注意事项:**
- 每个Cmd实例仅会被解析一次（线程安全）
- 不会处理任何子命令，所有参数均视为当前命令的标志或位置参数
- 处理内置标志逻辑

### PrintHelp

```go
func PrintHelp()
```

PrintHelp 输出全局默认命令实例 `QCommandLine` 的帮助信息。帮助信息通常包含命令的名称、可用的标志及其描述等内容。

### SetDesc

```go
func SetDesc(desc string)
```

SetDesc 设置命令描述信息。

**参数:**
- `desc`: 命令描述信息

### SetCompletion

```go
func SetCompletion(enable bool)
```

SetCompletion 设置是否启用自动完成功能。

**参数:**
- `enable`: 是否启用自动完成功能

### SetAutoExit

```go
func SetAutoExit(exit bool)
```

SetAutoExit 设置是否在解析内置参数时退出。默认情况下为true，当解析到内置参数时，QFlag将退出程序。

**参数:**
- `exit`: 是否退出

### SetHelp

```go
func SetHelp(help string)
```

SetHelp 配置全局默认命令实例 `QCommandLine` 的帮助信息。

**参数:**
- `help`: 新的帮助信息，字符串类型

### SetLogo

```go
func SetLogo(logoText string)
```

SetLogo 配置全局默认命令实例 `QCommandLine` 的 logo 文本。

**参数:**
- `logoText`: 配置的 logo 文本，字符串类型

### SetModules

```go
func SetModules(moduleHelps string)
```

SetModules 配置模块帮助信息。

**参数:**
- `moduleHelps`: 模块帮助信息，字符串类型

### SetUsage

```go
func SetUsage(usageSyntax string)
```

SetUsage 配置全局默认命令实例 `QCommandLine` 的用法信息。

**参数:**
- `usage`: 新的用法信息，字符串类型

**示例:**
```go
qflag.SetUsage("Usage: qflag [options]")
```

### SetChinese

```go
func SetChinese(useChinese bool)
```

SetChinese 设置是否使用中文。该函数用于设置当前命令行标志是否使用中文。

**参数:**
- `useChinese`: 如果使用中文，则传入true；否则传入false

### SetVersion

```go
func SetVersion(version string)
```

SetVersion 为全局默认命令设置版本信息。

**参数:**
- `version`: 版本信息字符串，用于标识命令的版本

### SetVersionf

```go
func SetVersionf(format string, args ...any)
```

SetVersionf 为全局默认命令设置版本信息。

**参数:**
- `format`: 格式化字符串，用于标识命令的版本
- `args`: 可变参数列表，用于替换格式化字符串中的占位符

### ShortName

```go
func ShortName() string
```

ShortName 获取命令短名称。

### Size

```go
func Size(longName, shortName string, defValue int64, usage string) *flags.SizeFlag
```

Size 为全局默认命令定义一个大小类型的命令行标志。
该函数会调用全局默认命令实例 `QCommandLine` 的 `Size` 方法，为命令行添加支持长短标志的大小类型参数。
用户可以输入 "1KB", "5MB", "2.5GiB" 等带单位的字符串，该标志会自动解析为字节数。

**参数:**
- `longName`: 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式。
- `shortName`: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式。
- `defValue`: 该命令行标志的默认值，单位为字节 (int64)。
- `usage`: 该命令行标志的帮助说明信息，在显示帮助信息时会呈现给用户，用以解释该标志的具体用途。

**返回值:**
- `*flags.SizeFlag`: 指向新创建的大小类型标志对象的指针。

### SizeVar

```go
func SizeVar(f *flags.SizeFlag, longName, shortName string, defValue int64, usage string)
```

SizeVar 为全局默认命令将一个大小类型的命令行标志绑定到指定的 `SizeFlag` 指针。
该函数会调用全局默认命令实例 `QCommandLine` 的 `SizeVar` 方法，为命令行添加支持长短标志的大小类型参数。

**参数:**
- `f`: 指向要绑定的 `SizeFlag` 对象的指针。
- `longName`: 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式。
- `shortName`: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式。
- `defValue`: 该命令行标志的默认值，单位为字节 (int64)。
- `usage`: 该命令行标志的帮助说明信息，在显示帮助信息时会呈现给用户，用以解释该标志的具体用途。

### StringSlice

```go
func StringSlice(longName, shortName string, defValue []string, usage string) *flags.StringSliceFlag
```

StringSlice 为全局默认命令定义一个字符串切片类型的命令行标志。该函数会调用全局默认命令实例 `QCommandLine` 的 `StringSlice` 方法，为命令行添加支持长短标志的字符串切片类型参数。

**参数:**
- `longName`: 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式
- `shortName`: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式
- `defValue`: 该命令行标志的默认值。当用户在命令行中未指定该标志时，会采用此默认值。该值会被复制一份，避免外部修改影响内部状态
- `usage`: 该命令行标志的帮助说明信息，在显示帮助信息时会呈现给用户，用以解释该标志的具体用途

**返回值:**
- `*flags.StringSliceFlag`: 指向新创建的字符串切片类型标志对象的指针

### StringSliceVar

```go
func StringSliceVar(f *flags.StringSliceFlag, longName, shortName string, defValue []string, usage string)
```

StringSliceVar 为全局默认命令将一个字符串切片类型的命令行标志绑定到指定的 `StringSliceFlag` 指针。该函数会调用全局默认命令实例 `QCommandLine` 的 `StringSliceVar` 方法，为命令行添加支持长短标志的字符串切片类型参数。

**参数:**
- `f`: 指向要绑定的 `StringSliceFlag` 对象的指针
- `longName`: 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式
- `shortName`: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式
- `defValue`: 该命令行标志的默认值。当用户在命令行中未指定该标志时，会采用此默认值。该值会被复制一份，避免外部修改影响内部状态
- `usage`: 该命令行标志的帮助说明信息，在显示帮助信息时会呈现给用户，用以解释该标志的具体用途

### IntSlice

```go
func IntSlice(longName, shortName string, defValue []int, usage string) *flags.IntSliceFlag
```

IntSlice 为全局默认命令定义一个整数切片类型的命令行标志。该函数会调用全局默认命令实例 `QCommandLine` 的 `IntSlice` 方法，为命令行添加支持长短标志的整数切片类型参数。

**参数:**
- `longName`: 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式
- `shortName`: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式
- `defValue`: 该命令行标志的默认值。当用户在命令行中未指定该标志时，会采用此默认值。该值会被复制一份，避免外部修改影响内部状态
- `usage`: 该命令行标志的帮助说明信息，在显示帮助信息时会呈现给用户，用以解释该标志的具体用途

**返回值:**
- `*flags.IntSliceFlag`: 指向新创建的整数切片类型标志对象的指针

### IntSliceVar

```go
func IntSliceVar(f *flags.IntSliceFlag, longName, shortName string, defValue []int, usage string)
```

IntSliceVar 为全局默认命令将一个整数切片类型的命令行标志绑定到指定的 `IntSliceFlag` 指针。该函数会调用全局默认命令实例 `QCommandLine` 的 `IntSliceVar` 方法，为命令行添加支持长短标志的整数切片类型参数。

**参数:**
- `f`: 指向要绑定的 `IntSliceFlag` 对象的指针
- `longName`: 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式
- `shortName`: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式
- `defValue`: 该命令行标志的默认值。当用户在命令行中未指定该标志时，会采用此默认值。该值会被复制一份，避免外部修改影响内部状态
- `usage`: 该命令行标志的帮助说明信息，在显示帮助信息时会呈现给用户，用以解释该标志的具体用途

### Int64Slice

```go
func Int64Slice(longName, shortName string, defValue []int64, usage string) *flags.Int64SliceFlag
```

Int64Slice 为全局默认命令定义一个64位整数切片类型的命令行标志。该函数会调用全局默认命令实例 `QCommandLine` 的 `Int64Slice` 方法，为命令行添加支持长短标志的64位整数切片类型参数。

**参数:**
- `longName`: 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式
- `shortName`: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式
- `defValue`: 该命令行标志的默认值。当用户在命令行中未指定该标志时，会采用此默认值。该值会被复制一份，避免外部修改影响内部状态
- `usage`: 该命令行标志的帮助说明信息，在显示帮助信息时会呈现给用户，用以解释该标志的具体用途

**返回值:**
- `*flags.Int64SliceFlag`: 指向新创建的64位整数切片类型标志对象的指针

### Int64SliceVar

```go
func Int64SliceVar(f *flags.Int64SliceFlag, longName, shortName string, defValue []int64, usage string)
```

Int64SliceVar 为全局默认命令将一个64位整数切片类型的命令行标志绑定到指定的 `Int64SliceFlag` 指针。该函数会调用全局默认命令实例 `QCommandLine` 的 `Int64SliceVar` 方法，为命令行添加支持长短标志的64位整数切片类型参数。

**参数:**
- `f`: 指向要绑定的 `Int64SliceFlag` 对象的指针
- `longName`: 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式
- `shortName`: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式
- `defValue`: 该命令行标志的默认值。当用户在命令行中未指定该标志时，会采用此默认值。该值会被复制一份，避免外部修改影响内部状态
- `usage`: 该命令行标志的帮助说明信息，在显示帮助信息时会呈现给用户，用以解释该标志的具体用途

### String

```go
func String(longName, shortName, defValue, usage string) *flags.StringFlag
```

String 为全局默认命令创建一个字符串类型的命令行标志。该函数会调用全局默认命令实例的 String 方法，为命令行添加一个支持长短标志的字符串参数。

**参数:**
- `longName`: 标志的长名称，在命令行中以 --name 的形式使用
- `shortName`: 标志的短名称，在命令行中以 -shortName 的形式使用
- `defValue`: 标志的默认值，当命令行未指定该标志时使用
- `usage`: 标志的帮助说明信息，用于在显示帮助信息时展示

**返回值:**
- `*flags.StringFlag`: 指向新创建的字符串标志对象的指针

### StringVar

```go
func StringVar(f *flags.StringFlag, longName, shortName, defValue, usage string)
```

StringVar 函数的作用是将一个字符串类型的命令行标志绑定到全局默认命令的 `StringFlag` 指针上。借助全局默认命令实例 `QCommandLine` 的 `StringVar` 方法，为命令行添加支持长短标志的字符串参数，并将该参数与传入的 `StringFlag` 指针关联，以便后续获取和使用该标志的值。

**参数:**
- `f`: 指向 `StringFlag` 的指针，用于存储和管理该字符串类型命令行标志的相关信息，包括当前值、默认值等
- `longName`: 命令行标志的长名称，在命令行中需以 `--name` 的格式使用
- `shortName`: 命令行标志的短名称，在命令行中需以 `-shortName` 的格式使用
- `defValue`: 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值
- `usage`: 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户，用于解释该标志的用途

### SubCmdMap

```go
func SubCmdMap() map[string]*Cmd
```

SubCmdMap 获取所有已注册的子命令映射。

### SubCmds

```go
func SubCmds() []*Cmd
```

SubCmds 获取所有已注册的子命令列表。

### Time

```go
func Time(longName, shortName string, defValue string, usage string) *flags.TimeFlag
```

Time 为全局默认命令定义一个时间类型的命令行标志。该函数会调用全局默认命令实例 `QCommandLine` 的 `Time` 方法，为命令行添加支持长短标志的时间类型参数。

**参数:**
- `longName`: 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式
- `shortName`: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式
- `defValue`: 命令行标志的默认值(时间表达式，如"now", "zero", "1h", "2006-01-02")
- `usage`: 命令行标志的用法说明

**返回值:**
- `*flags.TimeFlag`: 指向新创建的时间类型标志对象的指针

**支持的默认值格式:**
- `"now"` 或 `""`: 当前时间
- `"zero"`: 零时间 (time.Time{})
- `"1h"`, `"30m"`, `"-2h"`: 相对时间（基于当前时间的偏移）
- `"2006-01-02"`, `"2006-01-02 15:04:05"`: 绝对时间格式
- RFC3339等标准格式

### TimeVar

```go
func TimeVar(f *flags.TimeFlag, longName, shortName string, defValue string, usage string)
```

TimeVar 为全局默认命令定义一个时间类型的命令行标志，并将其绑定到指定的 `TimeFlag` 指针。该函数会调用全局默认命令实例 `QCommandLine` 的 `TimeVar` 方法，为命令行添加支持长短标志的时间类型参数。

**参数:**
- `f`: 指向要绑定的 `TimeFlag` 对象的指针
- `longName`: 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式
- `shortName`: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式
- `defValue`: 命令行标志的默认值(时间表达式，如"now", "zero", "1h", "2006-01-02")
- `usage`: 命令行标志的用法说明

**支持的默认值格式:**
- `"now"` 或 `""`: 当前时间
- `"zero"`: 零时间 (time.Time{})
- `"1h"`, `"30m"`, `"-2h"`: 相对时间（基于当前时间的偏移）
- `"2006-01-02"`, `"2006-01-02 15:04:05"`: 绝对时间格式
- RFC3339等标准格式

### Uint16

```go
func Uint16(longName, shortName string, defValue uint16, usage string) *flags.Uint16Flag
```

Uint16 为全局默认命令定义一个无符号16位整数类型的命令行标志。该函数会调用全局默认命令实例 `QCommandLine` 的 `Uint16` 方法，为命令行添加支持长短标志的无符号16位整数类型参数。

**参数:**
- `longName`: 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式
- `shortName`: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式
- `defValue`: 命令行标志的默认值
- `usage`: 命令行标志的用法说明

**返回值:**
- `*flags.Uint16Flag`: 指向新创建的无符号16位整数类型标志对象的指针

### Uint16Var

```go
func Uint16Var(f *flags.Uint16Flag, longName, shortName string, defValue uint16, usage string)
```

Uint16Var 函数创建一个无符号16位整数类型标志，并将其绑定到指定的 `Uint16Flag` 指针。该函数会调用全局默认命令实例 `QCommandLine` 的 `Uint16Var` 方法，为命令行添加支持长短标志的无符号16位整数类型参数。

**参数:**
- `f`: 指向要绑定的 `Uint16Flag` 对象的指针
- `longName`: 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式
- `shortName`: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式
- `defValue`: 命令行标志的默认值
- `usage`: 命令行标志的用法说明

### Uint32

```go
func Uint32(longName, shortName string, defValue uint32, usage string) *flags.Uint32Flag
```

Uint32 为全局默认命令创建一个无符号32位整数类型的命令行标志。该函数会调用全局默认命令实例的 Uint32 方法，为命令行添加一个支持长短标志的无符号32位整数类型参数。

**参数:**
- `longName`: 标志的长名称，在命令行中以 --longName 的形式使用
- `shortName`: 标志的短名称，在命令行中以 -shortName 的形式使用
- `defValue`: 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值
- `usage`: 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户

**返回值:**
- `*flags.Uint32Flag`: 指向新创建的无符号32位整数标志对象的指针

### Uint32Var

```go
func Uint32Var(f *flags.Uint32Flag, longName, shortName string, defValue uint32, usage string)
```

Uint32Var 创建并绑定一个无符号32位整数标志。

**参数:**
- `f`: 指向要绑定的标志对象的指针
- `longName`: 标志的完整名称，在命令行中以 --longName 的形式使用
- `shortName`: 标志的短名称，在命令行中以 -shortName 的形式使用
- `defValue`: 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值
- `usage`: 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户

### Uint64

```go
func Uint64(longName, shortName string, defValue uint64, usage string) *flags.Uint64Flag
```

Uint64 为全局默认命令创建一个无符号64位整数类型的命令行标志。该函数会调用全局默认命令实例的 Uint64 方法，为命令行添加一个支持长短标志的无符号64位整数类型参数。

**参数:**
- `longName`: 标志的长名称，在命令行中以 --longName 的形式使用
- `shortName`: 标志的短名称，在命令行中以 -s 的形式使用
- `defValue`: 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值
- `usage`: 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户

**返回值:**
- `*flags.Uint64Flag`: 指向新创建的无符号64位整数标志对象的指针

### Uint64Var

```go
func Uint64Var(f *flags.Uint64Flag, longName, shortName string, defValue uint64, usage string)
```

Uint64Var 为全局默认命令将一个无符号64位整数类型的命令行标志绑定到指定的 Uint64Flag 指针。该函数会调用全局默认命令实例的 Uint64Var 方法，为命令行添加支持长短标志的无符号64位整数类型参数，并将参数值绑定到指定的 Uint64Flag 指针变量中。

**参数:**
- `f`: 指向 Uint64Flag 的指针，用于存储和管理该无符号64位整数类型命令行标志的相关信息
- `longName`: 命令行标志的长名称，在命令行中需以 --longName 的格式使用
- `shortName`: 命令行标志的短名称，在命令行中需以 -shortName 的格式使用
- `defValue`: 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值
- `usage`: 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户

### WithDesc

```go
func WithDesc(desc string) *Cmd
```

WithDesc 设置命令描述（链式调用）。

**参数:**
- `desc`: 命令描述

**返回值:**
- `*Cmd`: 返回命令实例，支持链式调用

### WithCompletion

```go
func WithCompletion(enable bool) *Cmd
```

WithCompletion 设置是否启用自动补全（链式调用）。

**参数:**
- `enable`: true表示启用补全，false表示禁用

**返回值:**
- `*Cmd`: 返回命令实例，支持链式调用

### WithExample

```go
func WithExample(desc, usage string) *Cmd
```

WithExample 为命令添加使用示例（链式调用）。

**参数:**
- `desc`: 示例描述
- `usage`: 示例用法

**返回值:**
- `*Cmd`: 返回命令实例，支持链式调用

### WithExamples

```go
func WithExamples(examples []cmd.ExampleInfo) *Cmd
```

WithExamples 添加使用示例列表到命令（链式调用）。

**参数:**
- `examples`: 示例信息列表，每个元素为 ExampleInfo 类型

**返回值:**
- `*Cmd`: 返回命令实例，支持链式调用

### WithAutoExit

```go
func WithAutoExit(exit bool) *Cmd
```

WithAutoExit 设置是否在解析内置参数时退出（链式调用）。

**参数:**
- `exit`: 是否退出

**返回值:**
- `*Cmd`: 返回命令实例，支持链式调用

### WithHelp

```go
func WithHelp(help string) *Cmd
```

WithHelp 设置用户自定义命令帮助信息（链式调用）。

**参数:**
- `help`: 用户自定义命令帮助信息

**返回值:**
- `*Cmd`: 返回命令实例，支持链式调用

### WithLogo

```go
func WithLogo(logoText string) *Cmd
```

WithLogo 设置logo文本（链式调用）。

**参数:**
- `logoText`: logo文本字符串

**返回值:**
- `*Cmd`: 返回命令实例，支持链式调用

### WithModules

```go
func WithModules(moduleHelps string) *Cmd
```

WithModules 设置自定义模块帮助信息（链式调用）。

**参数:**
- `moduleHelps`: 自定义模块帮助信息

**返回值:**
- `*Cmd`: 返回命令实例，支持链式调用

### WithNote

```go
func WithNote(note string) *Cmd
```

WithNote 添加备注信息到命令（链式调用）。

**参数:**
- `note`: 备注信息

**返回值:**
- `*Cmd`: 返回命令实例，支持链式调用

### WithNotes

```go
func WithNotes(notes []string) *Cmd
```

WithNotes 添加备注信息切片到命令（链式调用）。

**参数:**
- `notes`: 备注信息列表

**返回值:**
- `*Cmd`: 返回命令实例，支持链式调用

### WithUsage

```go
func WithUsage(usageSyntax string) *Cmd
```

WithUsage 设置自定义命令用法（链式调用）。

**参数:**
- `usageSyntax`: 自定义命令用法

**返回值:**
- `*Cmd`: 返回命令实例，支持链式调用

### WithChinese

```go
func WithChinese(useChinese bool) *Cmd
```

WithChinese 设置是否使用中文帮助信息（链式调用）。

**参数:**
- `useChinese`: 是否使用中文帮助信息

**返回值:**
- `*Cmd`: 返回命令实例，支持链式调用

### WithVersion

```go
func WithVersion(version string) *Cmd
```

WithVersion 设置版本信息（链式调用）。

**参数:**
- `version`: 版本信息

**返回值:**
- `*Cmd`: 返回命令实例，支持链式调用

### WithVersionf

```go
func WithVersionf(format string, args ...any) *Cmd
```

WithVersionf 设置版本信息（链式调用，支持格式化）。

**参数:**
- `format`: 版本信息格式字符串
- `args`: 格式化参数

**返回值:**
- `*Cmd`: 返回命令实例，支持链式调用

## Types

### BoolFlag

```go
type BoolFlag = flags.BoolFlag
```

BoolFlag 导出flag包中的BoolFlag结构体。

### DurationFlag

```go
type DurationFlag = flags.DurationFlag
```

DurationFlag 导出flag包中的DurationFlag结构体。

### EnumFlag

```go
type EnumFlag = flags.EnumFlag
```

EnumFlag 导出flag包中的EnumFlag结构体。

### Flag

```go
type Flag = flags.Flag
```

Flag 导出flag包中的Flag结构体。

### Float64Flag

```go
type Float64Flag = flags.Float64Flag
```

Float64Flag 导出flag包中的Float64Flag结构体。

### Int64Flag

```go
type Int64Flag = flags.Int64Flag
```

Int64Flag 导出flag包中的Int64Flag结构体。

### IntFlag

```go
type IntFlag = flags.IntFlag
```

IntFlag 导出flag包中的IntFlag结构体。

### MapFlag

```go
type MapFlag = flags.MapFlag
```

MapFlag 导出flag包中的MapFlag结构体。

### StringSliceFlag

```go
type StringSliceFlag = flags.StringSliceFlag
```

StringSliceFlag 导出flag包中的StringSliceFlag结构体。

### IntSliceFlag

```go
type IntSliceFlag = flags.IntSliceFlag
```

IntSliceFlag 导出flag包中的IntSliceFlag结构体。

### Int64SliceFlag

```go
type Int64SliceFlag = flags.Int64SliceFlag
```

Int64SliceFlag 导出flag包中的Int64SliceFlag结构体。

### SizeFlag

```go
type SizeFlag = flags.SizeFlag
```

SizeFlag 导出flag包中的SizeFlag结构体。

### StringFlag

```go
type StringFlag = flags.StringFlag
```

StringFlag 导出flag包中的StringFlag结构体。

### TimeFlag

```go
type TimeFlag = flags.TimeFlag
```

TimeFlag 导出flag包中的TimeFlag结构体。

### Uint16Flag

```go
type Uint16Flag = flags.Uint16Flag
```

Uint16Flag 导出flag包中的UintFlag结构体。

### Uint32Flag

```go
type Uint32Flag = flags.Uint32Flag
```

Uint32Flag 导出flag包中的Uint32Flag结构体。

### Uint64Flag

```go
type Uint64Flag = flags.Uint64Flag
```

Uint64Flag 导出flag包中的Uint64Flag结构体。

### Cmd

```go
type Cmd = cmd.Cmd
```

Cmd 导出cmd包中的Cmd结构体。