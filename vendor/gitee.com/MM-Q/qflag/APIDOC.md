# Package qflag

Package qflag 根包统一导出入口 本文件用于将各子包的核心功能导出到根包，简化外部使用

Package qflag 提供对标准库 flag 的封装，自动实现长短标志，并默认绑定 -h/--help 标志打印帮助信息。用户可通过 Cmd.Help 字段自定义帮助内容，支持直接赋值字符串或从文件加载。

## VARIABLES

### NewCmd

```go
var NewCmd = cmd.NewCmd
```

NewCmd 导出 cmd 包中的 NewCmd 函数

### QCommandLine

```go
var QCommandLine = cmd.QCommandLine
```

QCommandLine 导出 cmd 包的全局默认 Command 实例

## FUNCTIONS

### AddExample

```go
func AddExample(e cmd.ExampleInfo)
```

AddExample 添加示例 该函数用于添加命令行标志的示例，这些示例将在命令行帮助信息中显示

**参数:**

  * `e` : 示例信息，ExampleInfo 类型。

### AddNote

```go
func AddNote(note string)
```

AddNote 添加注意事项 该函数用于添加命令行标志的注意事项，这些注意事项将在命令行帮助信息中显示

**参数:**

  * `note` : 注意事项内容，字符串类型。

### AddSubCmd

```go
func AddSubCmd(subCmds ...*cmd.Cmd) error
```

AddSubCmd 向全局默认命令实例 `QCommandLine` 添加一个或多个子命令 该函数会调用全局默认命令实例的 `AddSubCmd` 方法，支持批量添加子命令 在添加过程中，会检查子命令是否为 `nil` 以及是否存在循环引用，若有异常则返回错误信息

**参数:**

  * `subCmds` : 可变参数，接收一个或多个 `*Cmd` 类型的子命令实例

**返回值:**

  * `error` : 若添加子命令过程中出现错误（如子命令为 `nil` 或存在循环引用），则返回错误信息；否则返回 `nil` 。

### Arg

```go
func Arg(i int) string
```

Arg 获取全局默认命令实例 `QCommandLine` 解析后的指定索引位置的非标志参数 索引从 0 开始，若索引超出非标志参数切片的范围，将返回空字符串

**参数:**

  * `i` : 非标志参数的索引位置，从 0 开始计数

**返回值:**

  * `string` : 指定索引位置的非标志参数；若索引越界，则返回空字符串

### Args

```go
func Args() []string
```

Args 获取全局默认命令实例 `QCommandLine` 解析后的非标志参数切片。非标志参数是指命令行中未被识别为标志的参数

**返回值:**

  * `[]string` : 包含所有非标志参数的字符串切片。

### Bool

```go
func Bool(longName, shortName string, defValue bool, usage string) *flags.BoolFlag
```

Bool 为全局默认命令创建一个布尔类型的命令行标志。该函数会调用全局默认命令实例的 Bool 方法，为命令行添加一个支持长短标志的布尔参数。

**参数说明：**

  * `name` : 标志的长名称，在命令行中以 --name 的形式使用。
  * `shortName` : 标志的短名称，在命令行中以 -shortName 的形式使用。
  * `defValue` : 标志的默认值，当命令行未指定该标志时使用。
  * `usage` : 标志的帮助说明信息，用于在显示帮助信息时展示。

**返回值：**

  * `*flags.BoolFlag` : 指向新创建的布尔标志对象的指针。

### BoolVar

```go
func BoolVar(f *flags.BoolFlag, longName, shortName string, defValue bool, usage string)
```

BoolVar 函数的作用是将布尔类型的命令行标志绑定到全局默认命令实例 `QCommandLine` 中。它会调用全局默认命令实例的 `BoolVar` 方法，为命令行添加一个支持长短和短标志的布尔参数，并将该参数与传入的 `BoolFlag` 指针建立关联，后续可以通过该指针获取和使用该标志的值。

**参数说明：**

  * `f` : 指向 `BoolFlag` 类型的指针，用于存储和管理布尔类型命令行标志的相关信息，如当前值、默认值等。
  * `name` : 标志的长名称，在命令行中以 `--name` 的形式使用。
  * `shortName` : 标志的短名称，在命令行中以 `-shortName` 的形式使用。
  * `defValue` : 标志的默认值，当命令行未指定该标志时，会使用此默认值。
  * `usage` : 标志的帮助说明信息，用于在显示帮助信息时展示给用户，解释该标志的用途。

### CmdExists

```go
func CmdExists(cmdName string) bool
```

CmdExists 检查子命令是否存在

**参数:**

  * `cmdName` : 子命令名称

**返回:**

  * `bool` : 子命令是否存在

### Duration

```go
func Duration(longName, shortName string, defValue time.Duration, usage string) *flags.DurationFlag
```

Duration 为全局默认命令定义一个时间间隔类型的命令行标志。该函数会调用全局默认命令实例 `QCommandLine` 的 `Duration` 方法，为命令行添加支持长短标志的时间间隔类型参数，

**参数说明：**

  * `longName` : 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式。
  * `shortName` : 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式。
  * `defValue` : 该命令行标志的默认值。当用户在命令行中未指定该标志时，会采用此默认值。该值会被复制一份，避免外部修改影响内部状态。
  * `usage` : 该命令行标志的帮助说明信息，在显示帮助信息时会呈现给用户，用以解释该标志的具体用途。

**返回值：**

  * `*flags.DurationFlag` : 指向新创建的时间间隔类型标志对象的指针。

### DurationVar

```go
func DurationVar(f *flags.DurationFlag, longName, shortName string, defValue time.Duration, usage string)
```

DurationVar 为全局默认命令将一个时间间隔类型的命令行标志绑定到指定的 `DurationFlag` 指针。该函数会调用全局默认命令实例 `QCommandLine` 的 `DurationVar` 方法，为命令行添加支持长短标志的时间间隔类型参数，

**参数说明：**

  * `f` : 指向 `DurationFlag` 类型的指针，此指针用于存储和管理时间间隔类型命令行标志的各类信息，如当前标志的值、默认值等。
  * `longName` : 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式。
  * `shortName` : 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式。
  * `defValue` : 该命令行标志的默认值。当用户在命令行中未指定该标志时，会采用此默认值。该值会被复制一份，避免外部修改影响内部状态。
  * `usage` : 该命令行标志的帮助说明信息，在显示帮助信息时会呈现给用户，用以解释该标志的具体用途。

### Enum

```go
func Enum(longName, shortName string, defValue string, usage string, enumValues []string) *flags.EnumFlag
```

Enum 为全局默认命令定义一个枚举类型的命令行标志。该函数会调用全局默认命令实例 `QCommandLine` 的 `Enum` 方法，为命令行添加支持长短标志的枚举类型参数，

**参数说明：**

  * `name` : 标志的长名称，在命令行中以 `--name` 的形式使用。
  * `shortName` : 标志的短名称，在命令行中以 `-shortName` 的形式使用。
  * `defValue` : 标志的默认值，当命令行未指定该标志时使用。
  * `usage` : 标志的帮助说明信息，用于在显示帮助信息时展示。
  * `enumValues` : 枚举值的集合，用于指定标志可接受的取值范围。

**返回值：**

  * `*flags.EnumFlag` : 指向新创建的枚举类型标志对象的指针。

### EnumVar

```go
func EnumVar(f *flags.EnumFlag, longName, shortName string, defValue string, usage string, enumValues []string)
```

EnumVar 为全局默认命令将一个枚举类型的命令行标志绑定到指定的 `EnumFlag` 指针。该函数会调用全局默认命令实例 `QCommandLine` 的 `EnumVar` 方法，为命令行添加支持长短标志的枚举类型参数，

**参数说明：**

  * `f` : 指向 `EnumFlag` 类型的指针，此指针用于存储和管理枚举类型命令行标志的各类信息，如当前标志的值、默认值等。
  * `name` : 命令行标志的长名称，在命令行中使用时需遵循 `--name` 的格式。
  * `shortName` : 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式。
  * `defValue` : 该命令行标志的默认值。当用户在命令行中未指定该标志时，会采用此默认值。该值会被复制一份，避免外部修改影响内部状态。
  * `usage` : 该命令行标志的帮助说明信息，在显示帮助信息时会呈现给用户，用于解释该标志的具体用途。
  * `enumValues` : 枚举值的集合，用于指定标志可接受的取值范围。

### FlagExists

```go
func FlagExists(name string) bool
```

FlagExists 检查全局默认命令实例 `QCommandLine` 中是否存在指定名称的标志 该函数会调用全局默认命令实例的 `FlagExists` 方法，用于检查命令行中是否存在指定名称的标志

**参数:**

  * `name` : 要检查的标志名称，可以是长名称或短名称。

**返回值:**

  * `bool` : 若存在指定名称的标志，则返回 `true` ；否则返回 `false` 。

### FlagRegistry

```go
func FlagRegistry() *flags.FlagRegistry
```

FlagRegistry 获取标志注册表

**返回值:**

  * `*flags.FlagRegistry` : 标志注册表

### Float64

```go
func Float64(longName, shortName string, defValue float64, usage string) *flags.Float64Flag
```

Float64 为全局默认命令创建一个浮点数类型的命令行标志。该函数会调用全局默认命令实例的 Float64 方法，为命令行添加一个支持长短标志的浮点数参数。

**参数说明：**

  * `name` : 标志的长名称，在命令行中以 --name 的形式使用。
  * `shortName` : 标志的短名称，在命令行中以 -shortName 的形式使用。
  * `defValue` : 标志的默认值，当命令行未指定该标志时使用。
  * `usage` : 标志的帮助说明信息，用于在显示帮助信息时展示。

**返回值：**

  * `*flags.FloatFlag` : 指向新创建的浮点数标志对象的指针。

### Float64Var

```go
func Float64Var(f *flags.Float64Flag, longName, shortName string, defValue float64, usage string)
```

Float64Var 为全局默认命令绑定一个浮点数类型的命令行标志到指定的 `FloatFlag` 指针。该函数会调用全局默认命令实例 `QCommandLine` 的 `Float64Var` 方法，为命令行添加支持长短标志的浮点数参数，并将该参数与传入的 `FloatFlag` 指针关联，以便后续获取和使用该标志的值。

**参数说明：**

  * `f` : 指向 `FloatFlag` 的指针，用于存储和管理该浮点数类型命令行标志的相关信息，包括当前值、默认值等。
  * `name` : 命令行标志的长名称，在命令行中需以 `--name` 的格式使用。
  * `shortName` : 命令行标志的短名称，在命令行中需以 `-shortName` 的格式使用。
  * `defValue` : 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值。
  * `usage` : 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户，用于解释该标志的用途。

### GetDescription

```go
func GetDescription() string
```

GetDescription 获取命令描述信息

### GetExamples

```go
func GetExamples() []cmd.ExampleInfo
```

GetExamples 获取示例信息 该函数用于获取命令行标志的示例信息列表

**返回值:**

  * `[]ExampleInfo` : 示例信息列表，每个元素为 ExampleInfo 类型。

### GetHelp

```go
func GetHelp() string
```

GetHelp 返回全局默认命令实例 `QCommandLine` 的帮助信息

**返回值:**

  * `string` : 命令行帮助信息。

### GetLogoText

```go
func GetLogoText() string
```

GetLogoText 获取全局默认命令实例 `QCommandLine` 的 logo 文本

**返回值:**

  * `string` : 配置的 logo 文本。

### GetModuleHelps

```go
func GetModuleHelps() string
```

GetModuleHelps 获取模块帮助信息

**返回值:**

  * `string` : 模块帮助信息。

### GetNotes

```go
func GetNotes() []string
```

GetNotes 获取所有备注信息

### GetUsageSyntax

```go
func GetUsageSyntax() string
```

GetUsageSyntax 获取全局默认命令实例 `QCommandLine` 的用法信息

**返回值:**

  * `string` : 命令行用法信息。

### GetUseChinese

```go
func GetUseChinese() bool
```

GetUseChinese 获取是否使用中文 该函数用于获取当前命令行标志是否使用中文

**返回值:**

  * `bool` : 如果使用中文, 则返回 true; 否则返回 false。

### GetVersion

```go
func GetVersion() string
```

GetVersion 获取全局默认命令的版本信息

**返回值：**

  * `string` : 版本信息字符串。

### IP4

```go
func IP4(longName, shortName string, defValue string, usage string) *flags.IP4Flag
```

IP4 为全局默认命令创建一个 IPv4 地址类型的命令行标志。该函数会调用全局默认命令实例的 IP4 方法，为命令行添加一个支持长短标志的 IPv4 地址类型参数。

**参数说明：**

  * `longName` : 标志的长名称，在命令行中以 --longName 的形式使用。
  * `shortName` : 标志的短名称，在命令行中以 -shortName 的形式使用。
  * `defValue` : 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值
  * `usage` : 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户。

**返回值：**

  * `*flags.IP4Flag` : 指向新创建的 IPv4 地址标志对象的指针。

### IP4Var

```go
func IP4Var(f *flags.IP4Flag, longName, shortName string, defValue string, usage string)
```

IP4Var 为全局默认命令将一个 IPv4 地址类型的命令行标志绑定到指定的 IP4Flag 指针。该函数会调用全局默认命令实例的 IP4Var 方法，为命令行添加支持长短标志的 IPv4 地址类型参数，并将参数值绑定到指定的 IP4Flag 指针变量中。

**参数说明：**

  * `f` : 指向 IP4Flag 的指针，用于存储和管理该 IPv4 地址类型命令行标志的相关信息。
  * `longName` : 命令行标志的长名称，在命令行中需以 --longName 的格式使用。
  * `shortName` : 命令行标志的短名称，在命令行中需以 -shortName 的格式使用。
  * `defValue` : 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值。
  * `usage` : 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户。

### IP6

```go
func IP6(longName, shortName string, defValue string, usage string) *flags.IP6Flag
```

IP6 为全局默认命令创建一个 IPv6 地址类型的命令行标志。该函数会调用全局默认命令实例的 IP6 方法，为命令行添加一个支持长短标志的 IPv6 地址类型参数。

**参数说明：**

  * `longName` : 标志的长名称，在命令行中以 --longName 的形式使用。
  * `shortName` : 标志的短名称，在命令行中以 -shortName 的形式使用。
  * `defValue` : 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值
  * `usage` : 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户。

**返回值：**

  * `*flags.IP6Flag` : 指向新创建的 IPv6 地址标志对象的指针。

### IP6Var

```go
func IP6Var(f *flags.IP6Flag, longName, shortName string, defValue string, usage string)
```

IP6Var 为全局默认命令将一个 IPv6 地址类型的命令行标志绑定到指定的 IP6Flag 指针。该函数会调用全局默认命令实例的 IP6Var 方法，为命令行添加支持长短标志的 IPv6 地址类型参数，并将参数值绑定到指定的 IP6Flag 指针变量中。

**参数说明：**

  * `f` : 指向 IP6Flag 的指针，用于存储和管理该 IPv6 地址类型命令行标志的相关信息。
  * `longName` : 命令行标志的长名称，在命令行中需以 --longName 的格式使用。
  * `shortName` : 命令行标志的短名称，在命令行中需以 -shortName 的格式使用。
  * `defValue` : 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值。
  * `usage` : 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户。

### Int

```go
func Int(longName, shortName string, defValue int, usage string) *flags.IntFlag
```

Int 为全局默认命令创建一个整数类型的命令行标志。该函数会调用全局默认命令实例的 Int 方法，为命令行添加一个支持长短标志的整数参数。

**参数说明：**

  * `name` : 标志的长名称，在命令行中以 --name 的形式使用。
  * `shortName` : 标志的短名称，在命令行中以 -shortName 的形式使用。
  * `defValue` : 标志的默认值，当命令行未指定该标志时使用。
  * `usage` : 标志的帮助说明信息，用于在显示帮助信息时展示。

**返回值：**

  * `*flags.IntFlag` : 指向新创建的整数标志对象的指针。

### Int64

```go
func Int64(longName, shortName string, defValue int64, usage string) *flags.Int64Flag
```

Int64 为全局默认命令定义一个 64 位整数类型的命令行标志。该函数会调用全局默认命令实例 `QCommandLine` 的 `Int64` 方法，为命令行添加支持长短标志的 64 位整数类型参数，

**参数说明：**

  * `longName` : 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式。
  * `shortName` : 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式。
  * `defValue` : 命令行标志的默认值。
  * `usage` : 命令行标志的用法说明。

**返回值：**

  * `*flags.Int64Flag` : 指向新创建的 64 位整数类型标志对象的指针。

### Int64Var

```go
func Int64Var(f *flags.Int64Flag, longName, shortName string, defValue int64, usage string)
```

Int64Var 函数创建一个 64 位整数类型标志，并将其绑定到指定的 `Int64Flag` 指针 该函数会调用全局默认命令实例 `QCommandLine` 的 `Int64Var` 方法，为命令行添加支持长短标志的 64 位整数类型参数，

**参数说明：**

  * `f` : 指向要绑定的 `Int64Flag` 对象的指针。
  * `longName` : 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式。
  * `shortName` : 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式。
  * `defValue` : 命令行标志的默认值。
  * `usage` : 命令行标志的用法说明。

### IntVar

```go
func IntVar(f *flags.IntFlag, longName, shortName string, defValue int, usage string)
```

IntVar 函数的作用是将整数类型的命令行标志绑定到全局默认命令的 `IntFlag` 指针上。它借助全局默认命令实例 `QCommandLine` 的 `IntVar` 方法，为命令行添加支持长短标志的整数参数，并将该参数与传入的 `IntFlag` 指针建立关联，方便后续对该标志的值进行获取和使用。

**参数说明：**

  * `f` : 指向 `IntFlag` 类型的指针，此指针用于存储和管理整数类型命令行标志的各类信息，如当前标志的值、默认值等。
  * `name` : 命令行标志的长名称，在命令行中使用时需遵循 `--name` 的格式。
  * `shortName` : 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式。
  * `defValue` : 该命令行标志的默认值。当用户在命令行中未指定该标志时，会采用此默认值。
  * `usage` : 该命令行标志的帮助说明信息，在显示帮助信息时会呈现给用户，用以解释该标志的具体用途。

### IsParsed

```go
func IsParsed() bool
```

IsParsed 检查命令行参数是否已解析

**返回:**

  * `bool` : 是否已解析

### LoadHelp

```go
func LoadHelp(filepath string) error
```

LoadHelp 从文件中加载帮助信息

**参数:**

  * `filepath` : 文件路径，字符串类型。

**返回值:**

  * `error` : 如果加载失败，则返回错误信息；否则返回 nil。

**示例:**

  ```go
  qflag.LoadHelp("help.txt")
  ```

### LongName

```go
func LongName() string
```

LongName 获取命令长名称

### Map

```go
func Map(longName, shortName string, defValue map[string]string, usage string) *flags.MapFlag
```

Map 为全局默认命令创建一个键值对类型的命令行标志。该函数会调用全局默认命令实例的 Map 方法，为命令行添加一个支持长短标志的键值对参数。

**参数说明：**

  * `longName` : 标志的长名称，在命令行中以 --longName 的形式使用。
  * `shortName` : 标志的短名称，在命令行中以 -shortName 的形式使用。
  * `defValue` : 标志的默认值，当命令行未指定该标志时使用。
  * `usage` : 标志的帮助说明信息，用于在显示帮助信息时展示。

**返回值：**

  * `*flags.MapFlag` : 指向新创建的键值对标志对象的指针。

### MapVar

```go
func MapVar(f *flags.MapFlag, longName, shortName string, defValue map[string]string, usage string)
```

MapVar 为全局默认命令将一个键值对类型的命令行标志绑定到指定的 MapFlag 指针。该函数会调用全局默认命令实例的 MapVar 方法，为命令行添加支持长短标志的键值对参数，并将该参数与传入的 MapFlag 指针关联，以便后续获取和使用该标志的值。

**参数说明：**

  * `f` : 指向 MapFlag 的指针，用于存储和管理该键值对类型命令行标志的相关信息。
  * `longName` : 命令行标志的长名称，在命令行中需以 --longName 的格式使用。
  * `shortName` : 命令行标志的短名称，在命令行中需以 -shortName 的格式使用。
  * `defValue` : 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值。
  * `usage` : 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户。

### NArg

```go
func NArg() int
```

NArg 获取全局默认命令实例 `QCommandLine` 解析后的非标志参数的数量

**返回值:**

  * `int` : 非标志参数的数量。

### NFlag

```go
func NFlag() int
```

NFlag 获取全局默认命令实例 `QCommandLine` 解析后已定义和使用的标志的数量

**返回值:**

  * `int` : 标志的数量。

### Name

```go
func Name() string
```

Name 获取全局默认命令实例 `QCommandLine` 的名称

**返回值:**

  * 优先返回长名称, 如果长名称不存在则返回短名称

### Parse

```go
func Parse() error
```

Parse 完整解析命令行参数（含子命令处理） 主要功能：

  1. 解析当前命令的长短标志及内置标志
  2. 自动检测并解析子命令及其参数（若存在）
  3. 验证枚举类型标志的有效性

**参数：**

  * `args` : 原始命令行参数切片（包含可能的子命令及参数）

**返回值：**

  * 解析过程中遇到的错误（如标志格式错误、子命令解析失败等）

**注意事项：**

  * 每个 Cmd 实例仅会被解析一次（线程安全）
  * 若检测到子命令，会将剩余参数传递给子命令的 Parse 方法
  * 处理内置标志执行逻辑

### ParseFlagsOnly

```go
func ParseFlagsOnly() error
```

ParseFlagsOnly 仅解析当前命令的标志参数（忽略子命令） 主要功能：

  1. 解析当前命令的长短标志及内置标志
  2. 验证枚举类型标志的有效性
  3. 明确忽略所有子命令及后续参数

**参数：**

  * `args` : 原始命令行参数切片（子命令及后续参数会被忽略）

**返回值：**

  * 解析过程中遇到的错误（如标志格式错误等）

**注意事项：**

  * 每个 Cmd 实例仅会被解析一次（线程安全）
  * 不会处理任何子命令，所有参数均视为当前命令的标志或位置参数
  * 处理内置标志逻辑

### Path

```go
func Path(longName, shortName string, defValue string, usage string) *flags.PathFlag
```

Path 为全局默认命令创建一个路径类型的命令行标志。该函数会调用全局默认命令实例的 Path 方法，为命令行添加一个支持长短标志的路径参数。

**参数说明：**

  * `longName` : 标志的长名称，在命令行中以 --longName 的形式使用。
  * `shortName` : 标志的短名称，在命令行中以 -shortName 的形式使用。
  * `defValue` : 标志的默认值，当命令行未指定该标志时使用。
  * `usage` : 标志的帮助说明信息，用于在显示帮助信息时展示。

**返回值：**

  * `*flags.PathFlag` : 指向新创建的路径标志对象的指针。

### PathVar

```go
func PathVar(f *flags.PathFlag, longName, shortName string, defValue string, usage string)
```

PathVar 为全局默认命令将一个路径类型的命令行标志绑定到指定的 PathFlag 指针。该函数会调用全局默认命令实例的 PathVar 方法，为命令行添加支持长短标志的路径参数，并将该参数与传入的 PathFlag 指针关联，以便后续获取和使用该标志的值。

**参数说明：**

  * `f` : 指向 PathFlag 的指针，用于存储和管理该路径类型命令行标志的相关信息。
  * `longName` : 命令行标志的长名称，在命令行中需以 --longName 的格式使用。
  * `shortName` : 命令行标志的短名称，在命令行中需以 -shortName 的格式使用。
  * `defValue` : 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值。
  * `usage` : 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户。

### PrintHelp

```go
func PrintHelp()
```

PrintHelp 输出全局默认命令实例 `QCommandLine` 的帮助信息。帮助信息通常包含命令的名称、可用的标志及其描述等内容。

### SetDescription

```go
func SetDescription(desc string)
```

SetDescription 设置命令描述信息

### SetEnableCompletion

```go
func SetEnableCompletion(enable bool)
```

SetEnableCompletion 设置是否启用自动完成功能

**参数:**

  * `enable` : 是否启用自动完成功能

**注意:**

  * 该函数只能在根命令实例上调用

### SetExitOnBuiltinFlags

```go
func SetExitOnBuiltinFlags(exit bool)
```

SetExitOnBuiltinFlags 设置是否在解析内置参数时退出 默认情况下为 true，当解析到内置参数时，QFlag 将退出程序

**参数:**

  * `exit` : 是否退出

### SetHelp

```go
func SetHelp(help string)
```

SetHelp 配置全局默认命令实例 `QCommandLine` 的帮助信息

**参数:**

  * `help` : 新的帮助信息，字符串类型。

### SetLogoText

```go
func SetLogoText(logoText string)
```

SetLogoText 配置全局默认命令实例 `QCommandLine` 的 logo 文本

**参数:**

  * `logoText` : 配置的 logo 文本，字符串类型。

### SetModuleHelps

```go
func SetModuleHelps(moduleHelps string)
```

SetModuleHelps 配置模块帮助信息

**参数:**

  * `moduleHelps` : 模块帮助信息，字符串类型。

### SetUsageSyntax

```go
func SetUsageSyntax(usageSyntax string)
```

SetUsageSyntax 配置全局默认命令实例 `QCommandLine` 的用法信息

**参数:**

  * `usage` : 新的用法信息，字符串类型。

**示例:**

  ```go
  qflag.SetUsageSyntax("Usage: qflag [options]")
  ```

### SetUseChinese

```go
func SetUseChinese(useChinese bool)
```

SetUseChinese 设置是否使用中文 该函数用于设置当前命令行标志是否使用中文

**参数:**

  * `useChinese` : 如果使用中文, 则传入 true; 否则传入 false。

### SetVersion

```go
func SetVersion(version string)
```

SetVersion 为全局默认命令设置版本信息

**参数说明：**

  * `version` : 版本信息字符串，用于标识命令的版本。

### ShortName

```go
func ShortName() string
```

ShortName 获取命令短名称

### Slice

```go
func Slice(longName, shortName string, defValue []string, usage string) *flags.SliceFlag
```

Slice 为全局默认命令定义一个字符串切片类型的命令行标志。该函数会调用全局默认命令实例 `QCommandLine` 的 `Slice` 方法，为命令行添加支持长短标志的字符串切片类型参数，

**参数说明：**

  * `longName` : 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式。
  * `shortName` : 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式。
  * `defValue` : 该命令行标志的默认值。当用户在命令行中未指定该标志时，会采用此默认值。该值会被复制一份，避免外部修改影响内部状态。
  * `usage` : 该命令行标志的帮助说明信息，在显示帮助信息时会呈现给用户，用以解释该标志的具体用途。

**返回值：**

  * `*flags.SliceFlag` : 指向新创建的字符串切片类型标志对象的指针。

### SliceVar

```go
func SliceVar(f *flags.SliceFlag, longName, shortName string, defValue []string, usage string)
```

SliceVar 为全局默认命令将一个字符串切片类型的命令行标志绑定到指定的 `SliceFlag` 指针。该函数会调用全局默认命令实例 `QCommandLine` 的 `SliceVar` 方法，为命令行添加支持长短标志的字符串切片类型参数，

**参数说明：**

  * `f` : 指向要绑定的 `SliceFlag` 对象的指针。
  * `longName` : 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式。
  * `shortName` : 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式。
  * `defValue` : 该命令行标志的默认值。当用户在命令行中未指定该标志时，会采用此默认值。该值会被复制一份，避免外部修改影响内部状态。
  * `usage` : 该命令行标志的帮助说明信息，在显示帮助信息时会呈现给用户，用以解释该标志的具体用途。

### String

```go
func String(longName, shortName, defValue, usage string) *flags.StringFlag
```

String 为全局默认命令创建一个字符串类型的命令行标志。该函数会调用全局默认命令实例的 String 方法，为命令行添加一个支持长短标志的字符串参数。

**参数说明：**

  * `name` : 标志的长名称，在命令行中以 --name 的形式使用。
  * `shortName` : 标志的短名称，在命令行中以 -shortName 的形式使用。
  * `defValue` : 标志的默认值，当命令行未指定该标志时使用。
  * `usage` : 标志的帮助说明信息，用于在显示帮助信息时展示。

**返回值：**

  * `*flags.StringFlag` : 指向新创建的字符串标志对象的指针。

### StringVar

```go
func StringVar(f *flags.StringFlag, longName, shortName, defValue, usage string)
```

StringVar 函数的作用是将一个字符串类型的命令行标志绑定到全局默认命令的 `StringFlag` 指针上。借助全局默认命令实例 `QCommandLine` 的 `StringVar` 方法，为命令行添加支持长短标志的字符串参数，并将该参数与传入的 `StringFlag` 指针关联，以便后续获取和使用该标志的值。

**参数说明：**

  * `f` : 指向 `StringFlag` 的指针，用于存储和管理该字符串类型命令行标志的相关信息，包括当前值、默认值等。
  * `name` : 命令行标志的长名称，在命令行中需以 `--name` 的格式使用。
  * `shortName` : 命令行标志的短名称，在命令行中需以 `-shortName` 的格式使用。
  * `defValue` : 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值。
  * `usage` : 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户，用于解释该标志的用途。

### SubCmdMap

```go
func SubCmdMap() map[string]*cmd.Cmd
```

SubCmdMap 获取所有已注册的子命令映射

### SubCmds

```go
func SubCmds() []*cmd.Cmd
```

SubCmds 获取所有已注册的子命令列表

### Time

```go
func Time(longName, shortName string, defValue time.Time, usage string) *flags.TimeFlag
```

Time 为全局默认命令定义一个时间类型的命令行标志。该函数会调用全局默认命令实例 `QCommandLine` 的 `Time` 方法，为命令行添加支持长短标志的时间类型参数，

**参数说明：**

  * `longName` : 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式。
  * `shortName` : 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式。
  * `defValue` : 命令行标志的默认值。
  * `usage` : 命令行标志的用法说明。

**返回值：**

  * `*flags.TimeFlag` : 指向新创建的时间类型标志对象的指针。

### TimeVar

```go
func TimeVar(f *flags.TimeFlag, longName, shortName string, defValue time.Time, usage string)
```

TimeVar 为全局默认命令定义一个时间类型的命令行标志，并将其绑定到指定的 `TimeFlag` 指针。该函数会调用全局默认命令实例 `QCommandLine` 的 `TimeVar` 方法，为命令行添加支持长短标志的时间类型参数，

**参数说明：**

  * `f` : 指向要绑定的 `TimeFlag` 对象的指针。
  * `longName` : 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式。
  * `shortName` : 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式。
  * `defValue` : 命令行标志的默认值。
  * `usage` : 命令行标志的用法说明。

### URL

```go
func URL(longName, shortName string, defValue string, usage string) *flags.URLFlag
```

URL 为全局默认命令创建一个 URL 地址类型的命令行标志。该函数会调用全局默认命令实例的 URL 方法，为命令行添加一个支持长短标志的 URL 地址类型参数。

**参数说明：**

  * `longName` : 标志的长名称，在命令行中以 --longName 的形式使用。
  * `shortName` : 标志的短名称，在命令行中以 -shortName 的形式使用。
  * `defValue` : 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值
  * `usage` : 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户。

**返回值：**

  * `*flags.URLFlag` : 指向新创建的 URL 地址标志对象的指针。

### URLVar

```go
func URLVar(f *flags.URLFlag, longName, shortName string, defValue string, usage string)
```

URLVar 为全局默认命令将一个 URL 地址类型的命令行标志绑定到指定的 URLFlag 指针。该函数会调用全局默认命令实例的 URLVar 方法，为命令行添加支持长短标志的 URL 地址类型参数，并将参数值绑定到指定的 URLFlag 指针变量中。

**参数说明：**

  * `f` : 指向 URLFlag 的指针，用于存储和管理该 URL 地址类型命令行标志的相关信息。
  * `longName` : 命令行标志的长名称，在命令行中需以 --longName 的格式使用
  * `shortName` : 命令行标志的短名称，在命令行中需以 -shortName 的格式使用
  * `defValue` : 命令行标志的默认值
  * `usage` : 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户

### Uint16

```go
func Uint16(longName, shortName string, defValue uint16, usage string) *flags.Uint16Flag
```

Uint16 为全局默认命令定义一个无符号 16 位整数类型的命令行标志。该函数会调用全局默认命令实例 `QCommandLine` 的 `Uint16` 方法，为命令行添加支持长短标志的无符号 16 位整数类型参数，

**参数说明：**

  * `longName` : 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式。
  * `shortName` : 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式。
  * `defValue` : 命令行标志的默认值。
  * `usage` : 命令行标志的用法说明。

**返回值：**

  * `*flags.Uint16Flag` : 指向新创建的无符号 16 位整数类型标志对象的指针。

### Uint16Var

```go
func Uint16Var(f *flags.Uint16Flag, longName, shortName string, defValue uint16, usage string)
```

Uint16Var 函数创建一个无符号 16 位整数类型标志，并将其绑定到指定的 `Uint16Flag` 指针 该函数会调用全局默认命令实例 `QCommandLine` 的 `Uint16Var` 方法，为命令行添加支持长短标志的无符号 16 位整数类型参数，

**参数说明：**

  * `f` : 指向要绑定的 `Uint16Flag` 对象的指针。
  * `longName` : 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式。
  * `shortName` : 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式。
  * `defValue` : 命令行标志的默认值。
  * `usage` : 命令行标志的用法说明。

### Uint32

```go
func Uint32(longName, shortName string, defValue uint32, usage string) *flags.Uint32Flag
```

Uint32 为全局默认命令创建一个无符号 32 位整数类型的命令行标志。该函数会调用全局默认命令实例的 Uint32 方法，为命令行添加一个支持长短标志的无符号 32 位整数类型参数。

**参数说明：**

  * `longName` : 标志的长名称，在命令行中以 --longName 的形式使用。
  * `shortName` : 标志的短名称，在命令行中以 -shortName 的形式使用。
  * `defValue` : 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值。
  * `usage` : 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户。

**返回值：**

  * `*flags.Uint32Flag` : 指向新创建的无符号 32 位整数标志对象的指针。

### Uint32Var

```go
func Uint32Var(f *flags.Uint32Flag, longName, shortName string, defValue uint32, usage string)
```

Uint32Var 创建并绑定一个无符号 32 位整数标志。

**参数：**

  * `f` : 指向要绑定的标志对象的指针。
  * `longName` : 标志的完整名称，在命令行中以 --longName 的形式使用。
  * `shortName` : 标志的短名称，在命令行中以 -shortName 的形式使用。
  * `defValue` : 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值
  * `usage` : 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户。

### Uint64

```go
func Uint64(longName, shortName string, defValue uint64, usage string) *flags.Uint64Flag
```

Uint64 为全局默认命令创建一个无符号 64 位整数类型的命令行标志。该函数会调用全局默认命令实例的 Uint64 方法，为命令行添加一个支持长短标志的无符号 64 位整数类型参数。

**参数说明：**

  * `longName` : 标志的长名称，在命令行中以 --longName 的形式使用。
  * `shortName` : 标志的短名称，在命令行中以 -s 的形式使用。
  * `defValue` : 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值
  * `usage` : 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户。

**返回值：**

  * `*flags.Uint64Flag` : 指向新创建的无符号 64 位整数标志对象的指针。

### Uint64Var

```go
func Uint64Var(f *flags.Uint64Flag, longName, shortName string, defValue uint64, usage string)
```

Uint64Var 为全局默认命令将一个无符号 64 位整数类型的命令行标志绑定到指定的 Uint64Flag 指针。该函数会调用全局默认命令实例的 Uint64Var 方法，为命令行添加支持长短标志的无符号 64 位整数类型参数，并将参数值绑定到指定的 Uint64Flag 指针变量中。

**参数说明：**

  * `f` : 指向 Uint64Flag 的指针，用于存储和管理该无符号 64 位整数类型命令行标志的相关信息。
  * `longName` : 命令行标志的长名称，在命令行中需以 --longName 的格式使用。
  * `shortName` : 命令行标志的短名称，在命令行中需以 -shortName 的格式使用。
  * `defValue` : 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值。
  * `usage` : 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户。

## TYPES

### BoolFlag

```go
type BoolFlag = flags.BoolFlag
```

BoolFlag 导出 flag 包中的 BoolFlag 结构体

### Cmd

```go
type Cmd = cmd.Cmd
```

Cmd 导出 cmd 包中的 Cmd 结构体

### DurationFlag

```go
type DurationFlag = flags.DurationFlag
```

DurationFlag 导出 flag 包中的 DurationFlag 结构体

### EnumFlag

```go
type EnumFlag = flags.EnumFlag
```

EnumFlag 导出 flag 包中的 EnumFlag 结构体

### ExampleInfo

```go
type ExampleInfo = cmd.ExampleInfo
```

ExampleInfo 导出 cmd 包中的 ExampleInfo 结构体

### Flag

```go
type Flag = flags.Flag
```

Flag 导出 flag 包中的 Flag 结构体

### Float64Flag

```go
type Float64Flag = flags.Float64Flag
```

Float64Flag 导出 flag 包中的 Float64Flag 结构体

### IP4Flag

```go
type IP4Flag = flags.IP4Flag
```

IP4Flag 导出 flag 包中的 Ip4Flag 结构体

### IP6Flag

```go
type IP6Flag = flags.IP6Flag
```

IP6Flag 导出 flag 包中的 Ip6Flag 结构体

### Int64Flag

```go
type Int64Flag = flags.Int64Flag
```

Int64Flag 导出 flag 包中的 Int64Flag 结构体

### IntFlag

```go
type IntFlag = flags.IntFlag
```

IntFlag 导出 flag 包中的 IntFlag 结构体

### MapFlag

```go
type MapFlag = flags.MapFlag
```

MapFlag 导出 flag 包中的 MapFlag 结构体

### PathFlag

```go
type PathFlag = flags.PathFlag
```

PathFlag 导出 flag 包中的 PathFlag 结构体

### SliceFlag

```go
type SliceFlag = flags.SliceFlag
```

SliceFlag 导出 flag 包中的 SliceFlag 结构体

### StringFlag

```go
type StringFlag = flags.StringFlag
```

StringFlag 导出 flag 包中的 StringFlag 结构体

### TimeFlag

```go
type TimeFlag = flags.TimeFlag
```

TimeFlag 导出 flag 包中的 TimeFlag 结构体

### URLFlag

```go
type URLFlag = flags.URLFlag
```

URLFlag 导出 flag 包中的 URLFlag 结构体

### Uint16Flag

```go
type Uint16Flag = flags.Uint16Flag
```

Uint16Flag 导出 flag 包中的 UintFlag 结构体

### Uint32Flag

```go
type Uint32Flag = flags.Uint32Flag
```

Uint32Flag 导出 flag 包中的 Uint32Flag 结构体

### Uint64Flag

```go
type Uint64Flag = flags.Uint64Flag
```

Uint64Flag 导出 flag 包中的 Uint64Flag 结构体