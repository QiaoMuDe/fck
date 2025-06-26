# qflag

## 包说明

`package qflag // import "gitee.com/MM-Q/qflag"`

Package qflag 提供对标准库flag的封装, 自动实现长短标志, 并默认绑定-h/--help标志打印帮助信息。
用户可通过Cmd.Help字段自定义帮助内容, 支持直接赋值字符串或从文件加载。

## 目录

- [常量](#常量)
- [变量](#变量)
- [函数](#函数)
- [类型](#类型)

## 常量

```go
const (
	ErrFlagParseFailed       = "Parameter parsing error"  // 全局实例标志解析错误
	ErrSubCommandParseFailed = "Subcommand parsing error" // 子命令标志解析错误
	ErrPanicRecovered        = "panic recovered"          // 恐慌捕获错误
	ErrValidationFailed      = "Validation failed"        // 参数验证失败错误
)
```

命令行解析相关错误常量，用于标识不同类型的解析和验证错误。

## 变量

### ChineseTemplate

```go
var ChineseTemplate = HelpTemplate{
	CmdName:               "名称: %s\n\n",
	UsagePrefix:           "用法: ",
	UsageSubCmd:           " [子命令]",
	UseageInfoWithOptions: " [选项]\n\n",
	UseageGlobalOptions:   " [全局选项]",
	CmdNameWithShort:      "名称: %s(%s)\n\n",
	CmdDescription:        "描述: %s\n\n",
	OptionsHeader:         "选项:\n",
	Option1:               "  --%s, -%s",
	Option2:               "  --%s",
	OptionDefault:         "%s%*s%s (默认值: %s)\n",
	SubCmdsHeader:         "\n子命令:\n",
	SubCmd:                "  %s\t%s\n",
	SubCmdWithShort:       "  %s, %s\t%s\n",
	NotesHeader:           "\n注意事项:\n",
	NoteItem:              "  %d、%s\n",
	DefaultNote:           "当长选项和短选项同时使用时，最后指定的选项将优先生效。",
	ExamplesHeader:        "\n示例:\n",
	ExampleItem:           "  %d、%s\n    %s\n",
}

```

中文帮助信息模板实例，定义了中文环境下帮助信息的格式化方式。

### EnglishTemplate

```go
var EnglishTemplate = HelpTemplate{
	CmdName:               "Name: %s\n\n",
	UsagePrefix:           "Usage: ",
	UsageSubCmd:           " [subcmd]",
	UseageInfoWithOptions: " [options]\n\n",
	UseageGlobalOptions:   " [global options]",
	CmdNameWithShort:      "Name: %s(%s)\n\n",
	CmdDescription:        "Desc: %s\n\n",
	OptionsHeader:         "Options:\n",
	Option1:               "  --%s, -%s",
	Option2:               "  --%s",
	OptionDefault:         "%s%*s%s (default: %s)\n",
	SubCmdsHeader:         "\nSubCmds:\n",
	SubCmd:                "  %s\t%s\n",
	SubCmdWithShort:       "  %s, %s\t%s\n",
	NotesHeader:           "\nNotes:\n",
	NoteItem:              "  %d. %s\n",
	DefaultNote:           "In the case where both long options and short options are used at the same time,\n the option specified last shall take precedence.",
	ExamplesHeader:        "\nExamples:\n",
	ExampleItem:           "  %d. %s\n    %s\n",
}
```

英文帮助信息模板实例，定义了英文环境下帮助信息的格式化方式。

## 函数

### AddExample

```go
func AddExample(e ExampleInfo)
```

添加示例信息到命令行帮助文档。参数 `e`为示例信息结构体，包含描述和使用方式。
    AddExample 添加示例 该函数用于添加命令行标志的示例，这些示例将在命令行帮助信息中显示。 参数:
      - e: 示例信息，ExampleInfo 类型。

### AddNote

```go
func AddNote(note string)
```

添加注意事项到命令行帮助文档。参数 `note`为注意事项内容字符串，将在帮助信息的注意事项部分显示。

### AddSubCmd

```go
func AddSubCmd(subCmds ...*Cmd) error
```

向全局默认命令实例 `QCommandLine` 添加一个或多个子命令。

**参数:**

- `subCmds`: 可变参数，接收一个或多个 `*Cmd` 类型的子命令实例。

**返回值:**

- `error`: 若添加过程中出现错误（如子命令为 `nil` 或存在循环引用），则返回错误信息；否则返回 `nil`。

### Arg

```go
func Arg(i int) string
```

获取全局默认命令实例 `QCommandLine` 解析后的指定索引位置的非标志参数。

**参数:**

- `i`: 非标志参数的索引位置，从 0 开始计数。

**返回值:**

- `string`: 指定索引位置的非标志参数；若索引越界，则返回空字符串。

### Args

```go
func Args() []string
```

获取全局默认命令实例 `QCommandLine` 解析后的非标志参数切片。

**返回值:**

- `[]string`: 包含所有非标志参数的字符串切片。 非标志参数是指命令行中未被识别为标志的参数。

### BoolVar

```go
func BoolVar(f *BoolFlag, longName, shortName string, defValue bool, usage string)
```

将布尔类型的命令行标志绑定到全局默认命令实例 `QCommandLine` 中。

**参数:**

- `f`: 指向 `BoolFlag` 类型的指针，用于存储和管理布尔类型命令行标志的相关信息。
- `longName`: 标志的长名称，在命令行中以 `--longName` 的形式使用。
- `shortName`: 标志的短名称，在命令行中以 `-shortName` 的形式使用。
- `defValue`: 标志的默认值，当命令行未指定该标志时使用。
- `usage`: 标志的帮助说明信息，用于在显示帮助信息时展示给用户。

### GetDescription

```go
func GetDescription() string
```

获取全局默认命令实例 `QCommandLine` 的描述信息。

**返回值:**

- `string`: 命令的描述信息字符串。

### DurationVar

```go
func DurationVar(f *DurationFlag, longName, shortName string, defValue time.Duration, usage string)
```

将时间 duration 类型的命令行标志绑定到全局默认命令实例 `QCommandLine` 中。

**参数:**

- `f`: 指向 `DurationFlag` 类型的指针，用于存储和管理 duration 类型命令行标志的相关信息。
- `longName`: 标志的长名称，在命令行中以 `--longName` 的形式使用。
- `shortName`: 标志的短名称，在命令行中以 `-shortName` 的形式使用。
- `defValue`: 标志的默认值，当命令行未指定该标志时使用。
- `usage`: 标志的帮助说明信息，用于在显示帮助信息时展示给用户。

### EnumVar

```go
func EnumVar(f *EnumFlag, longName, shortName string, defValue string, usage string, enumValues []string)
```

将枚举类型的命令行标志绑定到全局默认命令实例 `QCommandLine` 中。

**参数:**

- `f`: 指向 `EnumFlag` 类型的指针，用于存储和管理枚举类型命令行标志的相关信息。
- `longName`: 标志的长名称，在命令行中以 `--longName` 的形式使用。
- `shortName`: 标志的短名称，在命令行中以 `-shortName` 的形式使用。
- `defValue`: 标志的默认值，当命令行未指定该标志时使用。
- `usage`: 标志的帮助说明信息，用于在显示帮助信息时展示给用户。
- `enumValues`: 枚举值的字符串切片，指定该标志允许的取值范围。

### FlagExists

```go
func FlagExists(name string) bool
```

检查全局默认命令实例 `QCommandLine` 中是否存在指定名称的标志。

**参数:**

- `name`: 要检查的标志名称，可以是长名称或短名称。

**返回值:**

- `bool`: 若标志存在，则返回 `true`；否则返回 `false`。

### FloatVar

```go
func FloatVar(f *FloatFlag, longName, shortName string, defValue float64, usage string)
```

将浮点类型的命令行标志绑定到全局默认命令实例 `QCommandLine` 中。

**参数:**

- `f`: 指向 `FloatFlag` 类型的指针，用于存储和管理浮点类型命令行标志的相关信息。
- `longName`: 标志的长名称，在命令行中以 `--longName` 的形式使用。
- `shortName`: 标志的短名称，在命令行中以 `-shortName` 的形式使用。
- `defValue`: 标志的默认值，当命令行未指定该标志时使用。
- `usage`: 标志的帮助说明信息，用于在显示帮助信息时展示给用户。

### GetExecutablePath

```go
func GetExecutablePath() string
```

获取当前程序的可执行文件路径。

**返回值:**

- `string`: 当前程序的可执行文件绝对路径字符串。

### GetLogoText

```go
func GetLogoText() string
```

获取全局默认命令实例 `QCommandLine` 的 Logo 文本。

**返回值:**

- `string`: Logo 文本字符串。

### GetModuleHelps

```go
func GetModuleHelps() string
```

获取全局默认命令实例 `QCommandLine` 的模块帮助信息。

**返回值:**

- `string`: 模块帮助信息字符串。

### GetNotes

```go
func GetNotes() []string
```

获取全局默认命令实例 `QCommandLine` 的注意事项列表。

**返回值:**

- `[]string`: 包含所有注意事项的字符串切片。

### GetUseChinese

```go
func GetUseChinese() bool
```

获取全局默认命令实例 `QCommandLine` 是否使用中文的配置状态。

**返回值:**

- `bool`: 若启用中文显示，则返回 `true`；否则返回 `false`。

### GetHelp

```go
func GetHelp() string
```

生成全局默认命令实例 `QCommandLine` 的帮助文档。

**返回值:**

- `string`: 包含命令行帮助信息的格式化字符串。

### IntVar

```go
func IntVar(f *IntFlag, longName, shortName string, defValue int, usage string)
```

将整数类型的命令行标志绑定到全局默认命令实例 `QCommandLine` 中。

**参数:**

- `f`: 指向 `IntFlag` 类型的指针，用于存储和管理整数类型命令行标志的相关信息。
- `longName`: 标志的长名称，在命令行中以 `--longName` 的形式使用。
- `shortName`: 标志的短名称，在命令行中以 `-shortName` 的形式使用。
- `defValue`: 标志的默认值，当命令行未指定该标志时使用。
- `usage`: 标志的帮助说明信息，用于在显示帮助信息时展示给用户。

### LongName

```go
func LongName() string
```

获取全局默认命令实例 `QCommandLine` 的长名称。

**返回值:**

- `string`: 命令的长名称字符串。

### NArg

```go
func NArg() int
```

获取全局默认命令实例 `QCommandLine` 的非标志参数数量。

**返回值:**

- `int`: 非标志参数的数量。

### NFlag

```go
func NFlag() int
```

获取全局默认命令实例 `QCommandLine` 解析后的标志参数数量。

**返回值:**

- `int`: 标志参数的数量。

### NewValidationError

```go
func NewValidationError(message string) error
```

创建一个新的验证错误实例。通常在命令行参数验证失败时使用。

**参数:**

- `message`: 错误消息字符串，描述验证失败的原因。

**返回值:**

- `error`: 包含指定消息的验证错误实例。

### NewValidationErrorf

```go
func NewValidationErrorf(format string, v ...interface{}) error
```

创建一个格式化的验证错误实例。支持使用格式化字符串和可变参数生成错误消息。

**参数:**

- `format`: 格式化字符串，用于指定错误消息的格式。
- `v`: 可变参数列表，用于填充格式化字符串中的占位符。

**返回值:**

- `error`: 包含格式化消息的验证错误实例。

### Parse

```go
func Parse() error
```

解析命令行参数并填充到相应的标志和参数中。该函数会调用全局默认命令实例 `QCommandLine` 的 `Parse` 方法。

**返回值:**

- `error`: 若解析过程中发生错误（如无效的标志、缺少必填参数等），则返回错误信息；否则返回 `nil`。

### ParseFlagsOnly

```go
func ParseFlagsOnly() error
```

解析命令行参数并填充到相应的标志中，但不处理子命令的参数解析。该函数会调用全局默认命令实例 `QCommandLine` 的 `ParseFlagsOnly` 方法。

**返回值:**

- `error`: 若解析过程中发生错误（如无效的标志、缺少必填参数等），则返回错误信息；否则返回 `nil`。

### PrintHelp

```go
func PrintHelp()
```

打印全局默认命令实例 `QCommandLine` 的帮助文档到标准输出。

**说明:**
该函数会调用全局默认命令实例的 `PrintHelp` 方法，将格式化的帮助信息（包括标志说明、使用示例等）输出到控制台。

### SetDescription

```go
func SetDescription(desc string)
```

设置全局默认命令实例 `QCommandLine` 的描述信息。

**参数:**

- `desc`: 命令的描述信息字符串，将显示在帮助文档中。

### SetHelp

```go
func SetHelp(help string)
```

设置全局默认命令实例 `QCommandLine` 的帮助信息。

**参数:**

- `help`: 命令的帮助信息字符串，将显示在帮助文档中。

### SetLogoText

```go
func SetLogoText(logoText string)
```

设置全局默认命令实例 `QCommandLine` 的 Logo 文本。

**参数:**

- `logoText`: Logo 文本字符串，将显示在帮助文档顶部。

### SetModuleHelps

```go
func SetModuleHelps(moduleHelps string)
```

设置全局默认命令实例 `QCommandLine` 的模块帮助信息。

**参数:**

- `moduleHelps`: 模块帮助信息字符串，用于描述模块功能和使用方法。

### SetUsageSyntax

```go
func SetUsageSyntax(usageSyntax string)
```

设置全局默认命令实例 `QCommandLine` 的使用说明。

**参数:**

- `usage`: 使用说明字符串，描述命令的基本用法格式。

### SetUseChinese

```go
func SetUseChinese(useChinese bool)
```

设置全局默认命令实例 `QCommandLine` 是否使用中文显示。

**参数:**

- `useChinese`: 布尔值，`true` 启用中文显示，`false` 使用英文显示。

### ShortName

```go
func ShortName() string
```

获取全局默认命令实例 `QCommandLine` 的短名称。

**返回值:**

- `string`: 命令的短名称字符串。

### StringVar

```go
func StringVar(f *StringFlag, longName, shortName, defValue, usage string)
```

将字符串类型的命令行标志绑定到全局默认命令实例 `QCommandLine` 中。

**参数:**

- `f`: 指向 `StringFlag` 类型的指针，用于存储和管理字符串类型命令行标志的相关信息。
- `longName`: 标志的长名称，在命令行中以 `--longName` 的形式使用。
- `shortName`: 标志的短名称，在命令行中以 `-shortName` 的形式使用。
- `defValue`: 标志的默认值，当命令行未指定该标志时使用。
- `usage`: 标志的帮助说明信息，用于在显示帮助信息时展示给用户。

## 类型

### BaseFlag 结构体

```go
type BaseFlag[T any] struct {
    // Has unexported fields.
}
```

泛型基础标志结构体，封装所有标志的通用字段和方法。

字段：

* `Name` ：标志的长名称。
* `Short` ：标志的短名称。
* `Usage` ：标志的帮助说明信息。
* `DefValue` ：标志的默认值。
* `Value` ：标志的当前值。
* `Changed` ：标志的值是否被修改过的标志位。
* `Validator` ：标志的验证器接口，用于自定义参数验证逻辑。

```go
func (f *BaseFlag[T]) Get() T
```

Get 获取标志的实际值。

```go
func (f *BaseFlag[T]) GetDefault() T
```

GetDefault 获取标志的默认值。

```go
func (f *BaseFlag[T]) GetDefaultAny() any
```

GetDefaultAny 获取标志的默认值(any类型)。

```go
func (f *BaseFlag[T]) LongName() string
```

LongName 获取标志的长名称。

```go
func (f *BaseFlag[T]) Set(value T) error
```

Set 设置标志的值。

```go
func (f *BaseFlag[T]) SetValidator(validator Validator)
```

SetValidator 设置标志的验证器 参数：validator 验证器接口。

```go
func (f *BaseFlag[T]) ShortName() string
```

ShortName 获取标志的短名称。

```go
func (f *BaseFlag[T]) Usage() string
```

Usage 获取标志的用法说明。

### BoolFlag 结构体

```go
type BoolFlag struct {
    BaseFlag[bool]
}
```

布尔类型标志的结构体，继承自 BaseFlag[bool] 泛型结构体，实现 Flag 接口。

字段：

* `BaseFlag[bool]` ：基础标志泛型结构体，包含通用属性和方法（如名称、短名称、默认值等）。

```go
func Bool(longName, shortName string, defValue bool, usage string) *BoolFlag
```

Bool 为全局默认命令创建一个布尔类型的命令行标志。该函数会调用全局默认命令实例的 Bool 方法，为命令行添加一个支持长短标志的布尔参数。

参数说明：

* `name` ：标志的长名称，在命令行中以 `--name` 的形式使用。
* `shortName` ：标志的短名称，在命令行中以 `-shortName` 的形式使用。
* `defValue` ：标志的默认值，当命令行未指定该标志时使用。
* `usage` ：标志的帮助说明信息，用于在显示帮助信息时展示。

返回值：

* `*BoolFlag` ：指向新创建的布尔标志对象的指针。

```go
func (f *BoolFlag) SetValidator(validator Validator)
```

SetValidator 设置标志的验证器 参数：validator 验证器接口。

```go
func (f *BoolFlag) Type() FlagType
```

Type 返回标志类型。

### Cmd 结构体

```go
type Cmd struct {
    // Has unexported fields.
}
```

Cmd 命令行标志管理结构体，封装参数解析、长短标志互斥及帮助系统。

```go
var QCommandLine *Cmd
```

QCommandLine 全局默认 Cmd 实例。

```go
func NewCmd(longName string, shortName string, errorHandling flag.ErrorHandling) *Cmd
```

NewCmd 创建新的命令实例 参数：longName：命令长名称 shortName：命令短名称 errorHandling：错误处理方式 返回值：*Cmd 命令实例指针 errorHandling 可选值：flag.ContinueOnError、flag.ExitOnError、flag.PanicOnError。

```go
func SubCmds() []*Cmd
```

SubCmds 获取所有已注册的子命令列表。

```go
func (c *Cmd) AddExample(e ExampleInfo)
```

AddExample 为命令添加使用示例 description：示例描述 usage：示例使用方式。

```go
func (c *Cmd) AddNote(note string)
```

AddNote 添加备注信息到命令。

```go
func (c *Cmd) AddSubCmd(subCmds ...*Cmd) error
```

AddSubCmd 关联一个或多个子命令到当前命令 支持批量添加多个子命令，遇到错误时收集所有错误并返回 参数：subCmds：一个或多个子命令实例指针 返回值：错误信息列表，如果所有子命令添加成功则返回 nil。

```go
func (c *Cmd) Arg(i int) string
```

Arg 获取指定索引的非标志参数。

```go
func (c *Cmd) Args() []string
```

Args 获取非标志参数切片。

```go
func (c *Cmd) Bool(longName, shortName string, defValue bool, usage string) *BoolFlag
```

Bool 添加布尔类型标志，返回标志对象指针 参数依次为：长标志名、短标志、默认值、帮助说明。

```go
func (c *Cmd) BoolVar(f *BoolFlag, longName, shortName string, defValue bool, usage string)
```

BoolVar 绑定布尔类型标志到指针并内部注册 Flag 对象 参数依次为：布尔标志指针、长标志名、短标志、默认值、帮助说明。

```go
func (c *Cmd) Description() string
```

Description 返回命令描述。

```go
func (c *Cmd) Duration(longName, shortName string, defValue time.Duration, usage string) *DurationFlag
```

Duration 添加时间间隔类型标志，返回标志对象指针 参数依次为：长标志名、短标志、默认值、帮助说明。

```go
func (c *Cmd) DurationVar(f *DurationFlag, longName, shortName string, defValue time.Duration, usage string)
```

DurationVar 绑定时间间隔类型标志到指针并内部注册 Flag 对象 参数依次为：时间间隔标志指针、长标志名、短标志、默认值、帮助说明。

```go
func (c *Cmd) Enum(longName, shortName string, defValue string, usage string, options []string) *EnumFlag
```

Enum 添加枚举类型标志，返回标志对象指针 参数依次为：长标志名、短标志、默认值、帮助说明、限制该标志取值的枚举值切片。

```go
func (c *Cmd) EnumVar(f *EnumFlag, longName, shortName string, defValue string, usage string, options []string)
```

EnumVar 绑定枚举类型标志到指针并内部注册 Flag 对象 参数依次为：枚举标志指针、长标志名、短标志、默认值、帮助说明、限制该标志取值的枚举值切片。

```go
func (c *Cmd) FlagExists(name string) bool
```

FlagExists 检查指定名称的标志是否存在。

```go
func (c *Cmd) Float(longName, shortName string, defValue float64, usage string) *FloatFlag
```

Float 添加浮点型标志，返回标志对象指针 参数依次为：长标志名、短标志、默认值、帮助说明。

```go
func (c *Cmd) FloatVar(f *FloatFlag, longName, shortName string, defValue float64, usage string)
```

FloatVar 绑定浮点型标志到指针并内部注册 Flag 对象 参数依次为：浮点数标志指针、长标志名、短标志、默认值、帮助说明。

```go
func (c *Cmd) GetExamples() []ExampleInfo
```

GetExamples 获取所有使用示例 返回示例切片的副本，防止外部修改。

```go
func (c *Cmd) GetLogoText() string
```

GetLogoText 获取 logo 文本。

```go
func (c *Cmd) GetModuleHelps() string
```

GetModuleHelps 获取自定义模块帮助信息。

```go
func (c *Cmd) GetNotes() []string
```

GetNotes 获取所有备注信息。

```go
func (c *Cmd) GetUseChinese() bool
```

GetUseChinese 获取是否使用中文帮助信息。

```go
func (c *Cmd) Help() string
```

Help 返回命令用法。

```go
func (c *Cmd) Int(longName, shortName string, defValue int, usage string) *IntFlag
```

Int 添加整数类型标志，返回标志对象指针 参数依次为：长标志名、短标志、默认值、帮助说明。

```go
func (c *Cmd) IntVar(f *IntFlag, longName, shortName string, defValue int, usage string)
```

IntVar 绑定整数类型标志到指针并内部注册 Flag 对象 参数依次为：整数标志指针、长标志名、短标志、默认值、帮助说明。

```go
func (c *Cmd) LongName() string
```

LongName 返回命令长名称。

```go
func (c *Cmd) NArg() int
```

NArg 获取非标志参数的数量。

```go
func (c *Cmd) NFlag() int
```

NFlag 获取标志的数量。

```go
func (c *Cmd) Parse(args []string) (err error)
```

Parse 解析命令行参数，自动检查长短标志，并处理内置标志 如果有子命令则会自动解析子命令的参数 参数：args：命令行参数切片 注意：该方法保证每个 Cmd 实例只会解析一次。

```go
func (c *Cmd) ParseFlagsOnly(args []string) (err error)
```

ParseFlagsOnly 仅解析当前命令的标志参数（忽略子命令, 不会自动解析子命令） 参数：args：命令行参数切片 注意：该方法保证每个 Cmd 实例只会解析一次。

```go
func (c *Cmd) PrintHelp()
```

PrintHelp 打印命令的帮助信息，优先打印用户的帮助信息，否则自动生成帮助信息。

```go
func (c *Cmd) SetDescription(desc string)
```

SetDescription 设置命令描述。

```go
func (c *Cmd) SetHelp(help string)
```

SetHelp 设置用户自定义命令帮助信息。

```go
func (c *Cmd) SetLogoText(logoText string)
```

SetLogoText 设置 logo 文本。

```go
func (c *Cmd) SetModuleHelps(moduleHelps string)
```

SetModuleHelps 设置自定义模块帮助信息。

```go
func (c *Cmd) SetUsage(usage string)
```

SetUsage 设置自定义命令用法。

```go
func (c *Cmd) SetUseChinese(useChinese bool)
```

SetUseChinese 设置是否使用中文帮助信息。

```go
func (c *Cmd) ShortName() string
```

ShortName 返回命令短名称。

```go
func (c *Cmd) String(longName, shortName, defValue, usage string) *StringFlag
```

String 添加字符串类型标志，返回标志对象指针 参数依次为：长标志名、短标志、默认值、帮助说明。

```go
func (c *Cmd) StringVar(f *StringFlag, longName, shortName, defValue, usage string)
```

StringVar 绑定字符串类型标志到指针并内部注册 Flag 对象 参数依次为：字符串标志指针、长标志名、短标志、默认值、帮助说明。

```go
func (c *Cmd) SubCmds() []*Cmd
```

SubCmds 返回子命令列表。

### CmdInterface 接口

```go
type CmdInterface interface {
    LongName() string // 获取命令名称(长名称)，如"app"
    ShortName() string // 获取命令短名称，如"a"
    GetDescription() string // 获取命令描述信息
    SetDescription(desc string) // 设置命令描述信息，用于帮助输出
    GetHelp() string // 获取自定义帮助信息
    SetHelp(help string) // 设置用户自定义命令帮助信息，覆盖自动生成内容
    SetUsageSyntax(usageSyntax string) // 设置自定义命令用法，覆盖自动生成内容
    GetUseChinese() bool // 获取是否使用中文帮助信息
    SetUseChinese(useChinese bool) // 设置是否使用中文帮助信息
    AddSubCmd(subCmd *Cmd) // 添加子命令，子命令会继承父命令的上下文
    SubCmds() []*Cmd // 获取所有已注册的子命令列表
    Parse(args []string) error // 解析命令行参数，自动处理标志和子命令
    ParseFlagsOnly(args []string) (err error) // 仅解析标志参数，不处理子命令
    Args() []string // 获取所有非标志参数(未绑定到任何标志的参数)
    Arg(i int) string // 获取指定索引的非标志参数，索引越界返回空字符串
    NArg() int // 获取非标志参数的数量
    NFlag() int // 获取已解析的标志数量
    FlagExists(name string) bool // 检查指定名称的标志是否存在(支持长/短名称)
    PrintHelp() // 打印命令帮助信息
    AddNote(note string) // 添加备注信息
    GetNotes() []string // 获取所有备注信息
    AddExample(e ExampleInfo) // 添加示例信息
    GetExamples() []ExampleInfo // 获取所有示例信息
    SetVersion(version string)  // 设置版本信息
    GetVersion() string     // 获取版本信息
    String(longName, shortName, usage, defValue string) *StringFlag // 添加字符串类型标志
    Int(longName, shortName, usage string, defValue int) *IntFlag // 添加整数类型标志
    Bool(longName, shortName, usage string, defValue bool) *BoolFlag // 添加布尔类型标志
    Float(longName, shortName, usage string, defValue float64) *FloatFlag // 添加浮点数类型标志
    Duration(longName, shortName, usage string, defValue time.Duration) *DurationFlag // 添加时间间隔类型标志
    Enum(longName, shortName string, defValue string, usage string, options []string) *EnumFlag // 添加枚举类型标志
    Slice(longName, shortName string, defValue []string, usage string) *SliceFlag                  // 添加字符串切片类型标志  
    StringVar(f *StringFlag, longName, shortName, defValue, usage string) // 绑定字符串标志到指定变量
    IntVar(f *IntFlag, longName, shortName string, defValue int, usage string) // 绑定整数标志到指定变量
    BoolVar(f *BoolFlag, longName, shortName string, defValue bool, usage string) // 绑定布尔标志到指定变量
    FloatVar(f *FloatFlag, longName, shortName string, defValue float64, usage string) // 绑定浮点数标志到指定变量
    DurationVar(f *DurationFlag, longName, shortName string, defValue time.Duration, usage string) // 绑定时间间隔类型标志到指定变量
    EnumVar(f *EnumFlag, longName, shortName string, defValue string, usage string, options []string) // 绑定枚举标志到指定变量
    SliceVar(f *SliceFlag, longName, shortName string, defValue []string, usage string)                  // 绑定字符串切片标志到指定变量  
    SetLogoText(logoText string) // 设置logo文本
    GetLogoText() string // 获取logo文本
    SetModuleHelps(moduleHelps string) // 设置自定义模块帮助信息
    GetModuleHelps() string // 获取自定义模块帮助信息
}
```

CmdInterface 命令接口定义，封装命令行程序的核心功能 提供统一的命令管理、参数解析和帮助系统 实现类需保证线程安全，所有方法应支持并发调用。

示例用法：`cmd := NewCmd("app", "a", flag.ContinueOnError)` `cmd.SetDescription("示例应用程序")` `cmd.String("config", "c", "配置文件路径", "/etc/app.conf")`

### DurationFlag 结构体

```go
type DurationFlag struct {
    BaseFlag[time.Duration]
}
```

DurationFlag 时间间隔类型标志结构体 继承 BaseFlag[time.Duration] 泛型结构体，实现 Flag 接口。

```go
func Duration(longName, shortName string, defValue time.Duration, usage string) *DurationFlag
```

Duration 为全局默认命令定义一个时间间隔类型的命令行标志。该函数会调用全局默认命令实例 `QCommandLine` 的 `Duration` 方法，为命令行添加支持长短标志的时间间隔类型参数 参数说明： - longName：命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式。 - shortName：命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式。 - defValue：该命令行标志的默认值。当用户在命令行中未指定该标志时，会采用此默认值。该值会被复制一份，避免外部修改影响内部状态。 - usage：该命令行标志的帮助说明信息，在显示帮助信息时会呈现给用户，用以解释该标志的具体用途。 返回值：。

```go
func (f *DurationFlag) Set(value string) error
```

Set 实现 flag.Value 接口，解析并设置时间间隔值。

```go
func (f *DurationFlag) SetValidator(validator Validator)
```

SetValidator 设置标志的验证器 参数：validator 验证器接口。

```go
func (f *DurationFlag) String() string
```

String 实现 flag.Value 接口，返回当前值的字符串表示。

```go
func (f *DurationFlag) Type() FlagType
```

Type 返回标志类型。

### EnumFlag 结构体

```go
type EnumFlag struct {
    BaseFlag[string]
    // Has unexported fields.
}
```

EnumFlag 枚举类型标志结构体 继承 BaseFlag[string] 泛型结构体，增加枚举特有的选项验证。

```go
func Enum(longName, shortName string, defValue string, usage string, enumValues []string) *EnumFlag
```

Enum 为全局默认命令定义一个枚举类型的命令行标志。该函数会调用全局默认命令实例 `QCommandLine` 的 `Enum` 方法，为命令行添加支持长短标志的枚举类型参数 参数说明： - name：标志的长名称，在命令行中以 `--name` 的形式使用。 - shortName：标志的短名称，在命令行中以 `-shortName` 的形式使用。 - defValue：标志的默认值，当命令行未指定该标志时使用。 - usage：标志的帮助说明信息，用于在显示帮助信息时展示。 - enumValues：枚举值的集合，用于指定标志可接受的取值范围。 返回值： - *EnumFlag：指向新创建的枚举类型标志对象的指针。

```go
func (f *EnumFlag) IsCheck(value string) error
```

IsCheck 检查枚举值是否有效 返回值：为 nil，说明值有效，否则返回错误信息。

```go
func (f *EnumFlag) Set(value string) error
```

Set 实现 flag.Value 接口，解析并设置枚举值。

```go
func (f *EnumFlag) SetValidator(validator Validator)
```

SetValidator 设置标志的验证器 参数：validator 验证器接口。

```go
func (f *EnumFlag) String() string
```

String 实现 flag.Value 接口，返回当前值的字符串表示。

```go
func (f *EnumFlag) Type() FlagType
```

实现 Flag 接口。

### ExampleInfo 结构体

```go
type ExampleInfo struct {
    Description string // 示例描述
    Usage string // 示例使用方式
}
```

ExampleInfo 示例信息结构体 用于存储命令的使用示例，包括描述和示例内容。

```go
func GetExamples() []ExampleInfo
```

GetExamples 获取示例信息 该函数用于获取命令行标志的示例信息列表。 返回值： - []ExampleInfo：示例信息列表，每个元素为 ExampleInfo 类型。

### Flag 接口

```go
type Flag interface {
    LongName() string // 获取标志的长名称
    ShortName() string // 获取标志的短名称
    Usage() string // 获取标志的用法
    Type() FlagType // 获取标志类型
    GetDefaultAny() any // 获取标志的默认值
    String() string     // 获取标志的字符串表示 
    IsSet() bool        // 判断标志是否已设置值 
    Reset()             // 重置标志值为默认值
}
```

Flag 所有标志类型的通用接口，定义了标志的元数据访问方法。

### FlagInfo 结构体

```go
type FlagInfo struct {
    // Has unexported fields.
}
```

FlagInfo 标志信息结构体 用于存储命令行标志的元数据，包括长名称、短名称、用法说明和默认值。

### FlagMeta 结构体

```go
type FlagMeta struct {
    // Has unexported fields.
}
```

FlagMeta 统一存储标志的完整元数据。

```go
func (m *FlagMeta) GetDefault() any
```

GetDefault 获取标志的默认值。

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

### FlagMetaInterface 接口

```go
type FlagMetaInterface interface {
    GetFlagType() FlagType // 获取标志类型
    GetFlag() Flag // 获取标志对象
    GetLongName() string // 获取标志的长名称
    GetShortName() string // 获取标志的短名称
    GetUsage() string // 获取标志的用法描述
    GetDefault() any // 获取标志的默认值
    GetValue() any // 获取标志的当前值
}
```

FlagMetaInterface 标志元数据接口，定义了标志元数据的获取方法。

### FlagRegistry 结构体

```go
type FlagRegistry struct {
    // Has unexported fields.
}
```

FlagRegistry 集中管理所有标志元数据及索引。

```go
func (r *FlagRegistry) GetAllFlags() []*FlagMeta
```

GetAllFlags 获取所有标志元数据列表 返回值： - []*FlagMeta：所有标志元数据的切片。

```go
func (r *FlagRegistry) GetByLong(longName string) (*FlagMeta, bool)
```

GetByLong 通过长标志名称查找对应的标志元数据 参数： - longName：标志的长名称(如 "help") 返回值： - *FlagMeta：找到的标志元数据指针，未找到时为 nil - bool：是否找到标志，true 表示找到。

```go
func (r *FlagRegistry) GetByName(name string) (*FlagMeta, bool)
```

GetByName 通过标志名称查找标志元数据 参数 name 可以是长名称(如 "help")或短名称(如 "h") 返回值： - *FlagMeta：找到的标志元数据指针，未找到时为 nil - bool：是否找到标志，true 表示找到。

```go
func (r *FlagRegistry) GetByShort(shortName string) (*FlagMeta, bool)
```

GetByShort 通过短标志名称查找对应的标志元数据 参数： - shortName：标志的短名称(如 "h" 对应 "help") 返回值： - *FlagMeta：找到的标志元数据指针，未找到时为 nil - bool：是否找到标志，true 表示找到。

```go
func (r *FlagRegistry) GetLongFlags() map[string]*FlagMeta
```

GetLongFlags 获取长标志映射 返回值： - map[string]*FlagMeta：长标志名称到标志元数据的映射。

```go
func (r *FlagRegistry) GetShortFlags() map[string]*FlagMeta
```

GetShortFlags 获取短标志映射 返回值： - map[string]*FlagMeta：短标志名称到标志元数据的映射。

```go
func (r *FlagRegistry) RegisterFlag(meta *FlagMeta) error
```

RegisterFlag 注册一个新的标志元数据到注册表中 该方法会执行以下操作： 1. 检查长名称和短名称是否已存在 2. 将标志添加到长名称索引 3. 将标志添加到短名称索引 4. 将标志添加到所有标志列表 注意：该方法线程安全，但发现重复标志时会 panic。

### FlagRegistryInterface 接口

```go
type FlagRegistryInterface interface {
    GetAllFlags() []*FlagMeta // 获取所有标志元数据列表
    GetLongFlags() map[string]*FlagMeta // 获取长标志映射
    GetShortFlags() map[string]*FlagMeta // 获取短标志映射
    RegisterFlag(meta *FlagMeta) error // 注册一个新的标志元数据到注册表中
    GetByLong(longName string) (*FlagMeta, bool) // 通过长标志名称查找对应的标志元数据
    GetByShort(shortName string) (*FlagMeta, bool) // 通过短标志名称查找对应的标志元数据
    GetByName(name string) (*FlagMeta, bool) // 通过标志名称查找标志元数据
}
```

FlagRegistryInterface 标志注册表接口，定义了标志元数据的增删改查操作。

### FlagType 类型

```go
type FlagType int
```

标志类型。

```go
const (
    FlagTypeInt      FlagType = iota + 1 // 整数类型
    FlagTypeString                       // 字符串类型
    FlagTypeBool                         // 布尔类型
    FlagTypeFloat                        // 浮点数类型
    FlagTypeEnum                         // 枚举类型
    FlagTypeDuration                     // 时间间隔类型
)
```

### FloatFlag 结构体

```go
type FloatFlag struct {
    BaseFlag[float64]
}
```

FloatFlag 浮点型标志结构体 继承 BaseFlag[float64] 泛型结构体，实现 Flag 接口。

```go
func Float(longName, shortName string, defValue float64, usage string) *FloatFlag
```

Float 为全局默认命令创建一个浮点数类型的命令行标志。该函数会调用全局默认命令实例的 Float 方法，为命令行添加一个支持长短标志的浮点数参数。 参数说明： - name：标志的长名称，在命令行中以 `--name` 的形式使用。 - shortName：标志的短名称，在命令行中以 `-shortName` 的形式使用。 - defValue：标志的默认值，当命令行未指定该标志时使用。 - usage：标志的帮助说明信息，用于在显示帮助信息时展示。 返回值： - *FloatFlag：指向新创建的浮点数标志对象的指针。

```go
func (f *FloatFlag) SetValidator(validator Validator)
```

SetValidator 设置标志的验证器 参数：validator 验证器接口。

```go
func (f *FloatFlag) Type() FlagType
```

Type 返回标志类型。

### HelpTemplate 结构体

```go
type HelpTemplate struct {
    CmdName               string // 命令名称模板
    CmdNameWithShort      string // 命令名称带短名称模板
    CmdDescription        string // 命令描述模板
    UsagePrefix           string // 用法说明前缀模板
    UsageSubCmd           string // 用法说明子命令模板
    UseageInfoWithOptions string // 带选项的用法说明信息模板
    UseageGlobalOptions   string // 全局选项部分
    OptionsHeader         string // 选项头部模板
    Option1               string // 选项模板(带短选项)
    Option2               string // 选项模板(无短选项)
    OptionDefault         string // 选项模板的默认值
    SubCmdsHeader         string // 子命令头部模板
    SubCmd                string // 子命令模板
    SubCmdWithShort       string // 子命令带短名称模板
    NotesHeader           string // 注意事项头部模板
    NoteItem              string // 注意事项项模板
    DefaultNote           string // 默认注意事项
    ExamplesHeader        string // 示例信息头部模板
    ExampleItem           string // 示例信息项模板
}
```

HelpTemplate 帮助信息模板结构体。

### IntFlag 结构体

```go
type IntFlag struct {
    BaseFlag[int]
}
```

IntFlag 整数类型标志结构体 继承 BaseFlag[int] 泛型结构体，实现 Flag 接口。

```go
func Int(longName, shortName string, defValue int, usage string) *IntFlag
```

Int 为全局默认命令创建一个整数类型的命令行标志。该函数会调用全局默认命令实例的 Int 方法，为命令行添加一个支持长短标志的整数参数。 参数说明： - name：标志的长名称，在命令行中以 `--name` 的形式使用。 - shortName：标志的短名称，在命令行中以 `-shortName` 的形式使用。 - defValue：标志的默认值，当命令行未指定该标志时使用。 - usage：标志的帮助说明信息，用于在显示帮助信息时展示。 返回值： - *IntFlag：指向新创建的整数标志对象的指针。

```go
func (f *IntFlag) SetValidator(validator Validator)
```

SetValidator 设置标志的验证器 参数：validator 验证器接口。

```go
func (f *IntFlag) Type() FlagType
```

Type 返回标志类型。

### QCommandLineInterface 接口

```go
type QCommandLineInterface interface {
    LongName() string // 获取命令长名称
    ShortName() string // 获取命令短名称
    GetDescription() string // 获取命令描述信息
    SetDescription(desc string) // 设置命令描述信息
    GetHelp() string // 获取命令帮助信息
    SetHelp(help string) // 设置命令帮助信息
    SetUsageSyntax(usage string) // 设置命令用法格式
    GetUseChinese() bool // 获取是否使用中文帮助信息
    SetUseChinese(useChinese bool) // 设置是否使用中文帮助信息
    AddSubCmd(subCmd *Cmd) // 添加子命令，子命令会继承父命令的上下文
    SubCmds() []*Cmd // 获取所有已注册的子命令列表
    Parse() error // 解析命令行参数，自动处理标志和子命令
    ParseFlagsOnly() error // 解析命令行参数，仅处理标志，不处理子命令
    Args() []string // 获取所有非标志参数(未绑定到任何标志的参数)
    Arg(i int) string // 获取指定索引的非标志参数，索引越界返回空字符串
    NArg() int // 获取非标志参数的数量
    NFlag() int // 获取已解析的标志数量
    PrintHelp() // 打印命令帮助信息
    FlagExists(name string) bool // 检查指定名称的标志是否存在(支持长/短名称)
    AddNote(note string) // 添加一个注意事项
    GetNotes() []string // 获取所有备注信息
    AddExample(e ExampleInfo) // 添加一个示例信息
    GetExamples() []ExampleInfo // 获取示例信息列表
    SetVersion(version string) // 设置版本信息
    GetVersion() string // 获取版本信息
    String(longName, shortName, defValue, usage string) *StringFlag // 添加字符串类型标志
    Int(longName, shortName string, defValue int, usage string) *IntFlag // 添加整数类型标志
    Bool(longName, shortName string, defValue bool, usage string) *BoolFlag // 添加布尔类型标志
    Float(longName, shortName string, defValue float64, usage string) *FloatFlag // 添加浮点数类型标志
    Duration(longName, shortName string, defValue time.Duration, usage string) *DurationFlag // 添加时间间隔类型标志
    Enum(longName, shortName string, defValue string, usage string, enumValues []string) *EnumFlag // 添加枚举类型标志
    Slice(longName, shortName string, defValue []string, usage string) *SliceFlag                  // 添加字符串切片类型标志  
    StringVar(f *StringFlag, longName, shortName, defValue, usage string) // 绑定字符串标志到指定变量
    IntVar(f *IntFlag, longName, shortName string, defValue int, usage string) // 绑定整数标志到指定变量
    BoolVar(f *BoolFlag, longName, shortName string, defValue bool, usage string) // 绑定布尔标志到指定变量
    FloatVar(f *FloatFlag, longName, shortName string, defValue float64, usage string) // 绑定浮点数标志到指定变量
    DurationVar(f *DurationFlag, longName, shortName string, defValue time.Duration, usage string) // 绑定时间间隔类型标志到指定变量
    EnumVar(f *EnumFlag, longName, shortName string, defValue string, usage string, enumValues []string) // 绑定枚举标志到指定变量
    SliceVar(f *SliceFlag, longName, shortName string, defValue []string, usage string)                  // 绑定字符串切片标志到指定变量  
    SetLogoText(logoText string) // 设置logo文本
    GetLogoText() string // 获取logo文本
    SetModuleHelps(moduleHelps string) // 设置自定义模块帮助信息
    GetModuleHelps() string // 获取自定义模块帮助信息
}
```

QCommandLineInterface 定义了全局默认命令行接口，提供统一的命令行参数管理功能 该接口封装了命令行程序的常用操作，包括标志添加、参数解析和帮助信息展示。

### StringFlag 结构体

```go
type StringFlag struct {
    BaseFlag[string]
}
```

StringFlag 字符串类型标志结构体 继承 BaseFlag[string] 泛型结构体，实现 Flag 接口。

```go
func String(longName, shortName, defValue, usage string) *StringFlag
```

String 为全局默认命令创建一个字符串类型的命令行标志。该函数会调用全局默认命令实例的 String 方法，为命令行添加一个支持长短标志的字符串参数。 参数说明： - name：标志的长名称，在命令行中以 `--name` 的形式使用。 - shortName：标志的短名称，在命令行中以 `-shortName` 的形式使用。 - defValue：标志的默认值，当命令行未指定该标志时使用。 - usage：标志的帮助说明信息，用于在显示帮助信息时展示。 返回值： - *StringFlag：指向新创建的字符串标志对象的指针。

```go
func (f *StringFlag) SetValidator(validator Validator)
```

SetValidator 设置标志的验证器 参数：validator 验证器接口。

```go
func (f *StringFlag) Type() FlagType
```

Type 返回标志类型。

### TypedFlag 接口

```go
type TypedFlag[T any] interface {
    Flag                     // 继承标志接口
    GetDefault() T           // 获取标志的具体类型默认值
    Get() T                  // 获取标志的具体类型值
    GetPointer() *T          // 获取标志的具体类型值指针
    Set(T) error             // 设置标志的具体类型值
    SetValidator(Validator)  // 设置标志的验证器
}
```

TypedFlag 所有标志类型的通用接口，定义了标志的元数据访问方法和默认值访问方法。

### UserInfo 结构体

```go
type UserInfo struct {
    // Has unexported fields.
}
```

UserInfo 存储用户自定义信息的嵌套结构体。

### Validator 接口

```go
type Validator interface {
    // Validate 验证参数值是否合法
    // value: 待验证的参数值
    // 返回值: 验证通过返回 nil, 否则返回错误信息
    Validate(value any) error
}
```

Validator 验证器接口，所有自定义验证器需实现此接口。

```

```
