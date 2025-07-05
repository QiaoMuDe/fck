# Package cmd

Package cmd 提供了命令行参数解析和管理的功能，定义了命令行标志管理结构体 Cmd，封装了参数解析、长短标志互斥及帮助系统。

## 示例用法
```go
cmd := NewCmd("app", "a", flag.ContinueOnError)
cmd.SetDescription("示例应用程序")
cmd.String("config", "c", "配置文件路径", "/etc/app.conf")
```

## VARIABLES

### ChineseTemplate
中文模板实例，用于生成命令行工具的中文帮助信息。
```go
var ChineseTemplate = HelpTemplate{
    CmdName:              "名称: %s\n\n",
    UsagePrefix:          "用法: ",
    UsageSubCmd:          " [子命令]",
    UsageInfoWithOptions: " [选项]\n\n",
    UsageGlobalOptions:   " [全局选项]",
    CmdNameWithShort:     "名称: %s, %s\n\n",
    CmdDescription:       "描述: %s\n\n",
    OptionsHeader:        "选项:\n",
    Option1:              "  --%s, -%s %s",
    Option2:              "  --%s %s",
    Option3:              "  -%s %s",
    OptionDefault:        "%s%*s%s (默认值: %s)\n",
    SubCmdsHeader:        "\n子命令:\n",
    SubCmd:               "  %s\t%s\n",
    SubCmdWithShort:      "  %s, %s\t%s\n",
    NotesHeader:          "\n注意事项:\n",
    NoteItem:             "  %d、%s\n",
    DefaultNote:          "当长选项和短选项同时使用时，最后指定的选项将优先生效。",
    ExamplesHeader:       "\n示例:\n",
    ExampleItem:          "  %d、%s\n     %s\n",
}
```

### EnglishTemplate
英文模板实例，用于生成命令行工具的英文帮助信息。
```go
var EnglishTemplate = HelpTemplate{
    CmdName:              "Name: %s\n\n",
    UsagePrefix:          "Usage: ",
    UsageSubCmd:          " [subcmd]",
    UsageInfoWithOptions: " [options]\n\n",
    UsageGlobalOptions:   " [global options]",
    CmdNameWithShort:     "Name: %s, %s\n\n",
    CmdDescription:       "Desc: %s\n\n",
    OptionsHeader:        "Options:\n",
    Option1:              "  --%s, -%s %s",
    Option2:              "  --%s %s",
    Option3:              "  -%s %s",
    OptionDefault:        "%s%*s%s (default: %s)\n",
    SubCmdsHeader:        "\nSubCmds:\n",
    SubCmd:               "  %s\t%s\n",
    SubCmdWithShort:      "  %s, %s\t%s\n",
    NotesHeader:          "\nNotes:\n",
    NoteItem:             "  %d. %s\n",
    DefaultNote:          "In the case where both long options and short options are used at the same time,\n the option specified last shall take precedence.",
    ExamplesHeader:       "\nExamples:\n",
    ExampleItem:          "  %d. %s\n     %s\n",
}
```

## FUNCTIONS

### GetExecutablePath
获取程序的绝对安装路径，如果无法通过 os.Executable 获取路径，则使用 os.Args[0] 作为替代。
```go
func GetExecutablePath() string
```
返回：程序的绝对路径字符串。

## TYPES

### Cmd
命令行标志管理结构体，封装参数解析、长短标志互斥及帮助系统。
```go
type Cmd struct {
    // Has unexported fields.
}
```

#### QCommandLine
全局默认 Command 实例。
```go
var QCommandLine *Cmd
```

#### NewCmd
创建新的命令实例。
```go
func NewCmd(longName string, shortName string, errorHandling flag.ErrorHandling) *Cmd
```
参数：
- `longName`: 命令长名称
- `shortName`: 命令短名称
- `errorHandling`: 错误处理方式

返回值：
- `*Cmd`: 新的命令实例指针

`errorHandling`可选值：
- `flag.ContinueOnError`: 解析标志时遇到错误继续解析, 并返回错误信息
- `flag.ExitOnError`: 解析标志时遇到错误立即退出程序, 并返回错误信息
- `flag.PanicOnError`: 解析标志时遇到错误立即触发 panic

#### AddExample
为命令添加使用示例。
```go
func (c *Cmd) AddExample(e ExampleInfo)
```
参数：
- `e`: 示例信息，包含 description（示例描述）和 usage（示例使用方式）

#### AddNote
添加备注信息到命令。
```go
func (c *Cmd) AddNote(note string)
```

#### AddSubCmd
关联一个或多个子命令到当前命令，支持批量添加多个子命令，遇到错误时收集所有错误并返回。
```go
func (c *Cmd) AddSubCmd(subCmds ...*Cmd) error
```
参数：
- `subCmds`: 一个或多个子命令实例指针

返回值：
- 错误信息，如果所有子命令添加成功则返回 nil

#### Arg
获取指定索引的非标志参数。
```go
func (c *Cmd) Arg(i int) string
```

#### Args
获取非标志参数切片。
```go
func (c *Cmd) Args() []string
```

#### Bool
添加布尔类型标志，返回标志对象指针。
```go
func (c *Cmd) Bool(longName, shortName string, defValue bool, usage string) *flags.BoolFlag
```
参数依次为：长标志名、短标志、默认值、帮助说明。
返回值：布尔标志对象指针。

#### BoolVar
绑定布尔类型标志到指针并内部注册 Flag 对象。
```go
func (c *Cmd) BoolVar(f *flags.BoolFlag, longName, shortName string, defValue bool, usage string)
```
参数依次为：布尔标志指针、长标志名、短标志、默认值、帮助说明。

#### CmdExists
检查子命令是否存在。
```go
func (c *Cmd) CmdExists(cmdName string) bool
```
参数：
- `cmdName`: 子命令名称

返回：
- `bool`: 子命令是否存在

#### Duration
添加时间间隔类型标志，返回标志对象指针。
```go
func (c *Cmd) Duration(longName, shortName string, defValue time.Duration, usage string) *flags.DurationFlag
```
参数依次为：长标志名、短标志、默认值、帮助说明。
返回值：时间间隔标志对象指针。

#### DurationVar
绑定时间间隔类型标志到指针并内部注册 Flag 对象。
```go
func (c *Cmd) DurationVar(f *flags.DurationFlag, longName, shortName string, defValue time.Duration, usage string)
```
参数依次为：时间间隔标志指针、长标志名、短标志、默认值、帮助说明。

#### Enum
添加枚举类型标志，返回标志对象指针。
```go
func (c *Cmd) Enum(longName, shortName string, defValue string, usage string, options []string) *flags.EnumFlag
```
参数依次为：长标志名、短标志、默认值、帮助说明、限制该标志取值的枚举值切片。
返回值：枚举标志对象指针。

#### EnumVar
绑定枚举类型标志到指针并内部注册 Flag 对象。
```go
func (c *Cmd) EnumVar(f *flags.EnumFlag, longName, shortName string, defValue string, usage string, options []string)
```
参数依次为：枚举标志指针、长标志名、短标志、默认值、帮助说明、限制该标志取值的枚举值切片。

#### FlagExists
检查指定名称的标志是否存在。
```go
func (c *Cmd) FlagExists(name string) bool
```

#### Float64
添加浮点型标志，返回标志对象指针。
```go
func (c *Cmd) Float64(longName, shortName string, defValue float64, usage string) *flags.Float64Flag
```
参数依次为：长标志名、短标志、默认值、帮助说明。
返回值：浮点型标志对象指针。

#### Float64Var
绑定浮点型标志到指针并内部注册 Flag 对象。
```go
func (c *Cmd) Float64Var(f *flags.Float64Flag, longName, shortName string, defValue float64, usage string)
```
参数依次为：浮点数标志指针、长标志名、短标志、默认值、帮助说明。

#### GetDescription
返回命令描述。
```go
func (c *Cmd) GetDescription() string
```

#### GetExamples
获取所有使用示例，返回示例切片的副本，防止外部修改。
```go
func (c *Cmd) GetExamples() []ExampleInfo
```

#### GetHelp
返回命令用法帮助信息。
```go
func (c *Cmd) GetHelp() string
```

#### GetLogoText
获取 logo 文本。
```go
func (c *Cmd) GetLogoText() string
```

#### GetModuleHelps
获取自定义模块帮助信息。
```go
func (c *Cmd) GetModuleHelps() string
```

#### GetNotes
获取所有备注信息。
```go
func (c *Cmd) GetNotes() []string
```

#### GetUsageSyntax
获取自定义命令用法。
```go
func (c *Cmd) GetUsageSyntax() string
```

#### GetUseChinese
获取是否使用中文帮助信息。
```go
func (c *Cmd) GetUseChinese() bool
```

#### GetVersion
获取版本信息。
```go
func (c *Cmd) GetVersion() string
```

#### IP4
添加 IPv4 地址类型标志，返回标志对象指针。
```go
func (c *Cmd) IP4(longName, shortName string, defValue string, usage string) *flags.IP4Flag
```
参数依次为：长标志名、短标志、默认值、帮助说明。
返回值：IPv4 地址标志对象指针。

#### IP4Var
绑定 IPv4 地址类型标志到指针并内部注册 Flag 对象。
```go
func (c *Cmd) IP4Var(f *flags.IP4Flag, longName, shortName string, defValue string, usage string)
```
参数依次为：IPv4 标志指针、长标志名、短标志、默认值、帮助说明。

#### IP6
添加 IPv6 地址类型标志，返回标志对象指针。
```go
func (c *Cmd) IP6(longName, shortName string, defValue string, usage string) *flags.IP6Flag
```
参数依次为：长标志名、短标志、默认值、帮助说明。
返回值：IPv6 地址标志对象指针。

#### IP6Var
绑定 IPv6 地址类型标志到指针并内部注册 Flag 对象。
```go
func (c *Cmd) IP6Var(f *flags.IP6Flag, longName, shortName string, defValue string, usage string)
```
参数依次为：IPv6 标志指针、长标志名、短标志、默认值、帮助说明。

#### Int
添加整数类型标志，返回标志对象指针。
```go
func (c *Cmd) Int(longName, shortName string, defValue int, usage string) *flags.IntFlag
```
参数依次为：长标志名、短标志、默认值、帮助说明。
返回值：整数标志对象指针。

#### Int64
添加 64 位整数类型标志，返回标志对象指针。
```go
func (c *Cmd) Int64(longName, shortName string, defValue int64, usage string) *flags.Int64Flag
```
参数依次为：长标志名、短标志、默认值、帮助说明。
返回值：64 位整数标志对象指针。

#### Int64Var
绑定 64 位整数类型标志到指针并内部注册 Flag 对象。
```go
func (c *Cmd) Int64Var(f *flags.Int64Flag, longName, shortName string, defValue int64, usage string)
```
参数依次为：64 位整数标志指针、长标志名、短标志、默认值、帮助说明。

#### IntVar
绑定整数类型标志到指针并内部注册 Flag 对象。
```go
func (c *Cmd) IntVar(f *flags.IntFlag, longName, shortName string, defValue int, usage string)
```
参数依次为：整数标志指针、长标志名、短标志、默认值、帮助说明。

#### LoadHelp
从指定文件加载帮助信息。
```go
func (c *Cmd) LoadHelp(filePath string) error
```
参数：
- `filePath`: 帮助信息文件路径

返回值：
- `error`: 如果文件不存在或读取文件失败，则返回错误信息。

#### LongName
返回命令长名称。
```go
func (c *Cmd) LongName() string
```

#### Map
添加键值对类型标志，返回标志对象指针。
```go
func (c *Cmd) Map(longName, shortName string, defValue map[string]string, usage string) *flags.MapFlag
```
参数依次为：长标志名、短标志、默认值、帮助说明。
返回值：键值对标志对象指针。

#### MapVar
绑定键值对类型标志到指针并内部注册 Flag 对象。
```go
func (c *Cmd) MapVar(f *flags.MapFlag, longName, shortName string, defValue map[string]string, usage string)
```
参数依次为：键值对标志指针、长标志名、短标志、默认值、帮助说明。

#### NArg
获取非标志参数的数量。
```go
func (c *Cmd) NArg() int
```

#### NFlag
获取标志的数量。
```go
func (c *Cmd) NFlag() int
```

#### Name
获取命令名称，优先返回长名称，如果长名称不存在则返回短名称。
```go
func (c *Cmd) Name() string
```

#### Parse
完整解析命令行参数（含子命令处理）。
```go
func (c *Cmd) Parse(args []string) (err error)
```
主要功能：
1. 解析当前命令的长短标志及内置标志
2. 自动检测并解析子命令及其参数（若存在）
3. 验证枚举类型标志的有效性

参数：
- `args`: 原始命令行参数切片（包含可能的子命令及参数）

返回值：
- 解析过程中遇到的错误（如标志格式错误、子命令解析失败等）

注意事项：
- 每个 Cmd 实例仅会被解析一次（线程安全）
- 若检测到子命令，会将剩余参数传递给子命令的 Parse 方法
- 处理内置标志执行逻辑

#### ParseFlagsOnly
仅解析当前命令的标志参数，忽略子命令。
```go
func (c *Cmd) ParseFlagsOnly(args []string) (err error)
```
主要功能：
1. 解析当前命令的长短标志及内置标志
2. 验证枚举类型标志的有效性
3. 明确忽略所有子命令及后续参数

参数：
- `args`: 原始命令行参数切片（子命令及后续参数会被忽略）

返回值：
- 解析过程中遇到的错误（如标志格式错误等）

注意事项：
- 每个 Cmd 实例仅会被解析一次（线程安全）
- 不会处理任何子命令，所有参数均视为当前命令的标志或位置参数
- 处理内置标志逻辑

#### IsParsed
检查命令是否已完成解析。
```go
func (c *Cmd) IsParsed() bool
```
返回值：
- `bool`: 解析状态，true表示已解析（无论成功失败），false表示未解析

#### Path
添加路径类型标志，返回标志对象指针。
```go
func (c *Cmd) Path(longName, shortName string, defValue string, usage string) *flags.PathFlag
```
参数依次为：长标志名、短标志、默认值、帮助说明。
返回值：路径标志对象指针。

#### PathVar
绑定路径类型标志到指针并内部注册 Flag 对象。
```go
func (c *Cmd) PathVar(f *flags.PathFlag, longName, shortName string, defValue string, usage string)
```
参数依次为：路径标志指针、长标志名、短标志、默认值、帮助说明。

#### PrintHelp
打印命令的帮助信息，优先打印用户的帮助信息，否则自动生成帮助信息。
```go
func (c *Cmd) PrintHelp()
```
注意：
- 打印帮助信息时，不会自动退出程序。

#### SetDescription
设置命令描述。
```go
func (c *Cmd) SetDescription(desc string)
```

#### SetDisableBuiltinFlags
设置是否禁用内置标志注册。
```go
func (c *Cmd) SetDisableBuiltinFlags(disable bool) *Cmd
```
参数：
- `disable`: true 表示禁用内置标志注册，false 表示启用（默认）

返回值：
- 当前 Cmd 实例，支持链式调用。

#### SetExitOnBuiltinFlags
设置是否在解析内置参数时退出，默认情况下为 true，当解析到内置参数时，QFlag 将退出程序。
```go
func (c *Cmd) SetExitOnBuiltinFlags(exit bool) *Cmd
```
参数：
- `exit`: 是否退出。

返回值：
- *cmd.Cmd: 当前命令对象。

#### SetHelp
设置用户自定义命令帮助信息。
```go
func (c *Cmd) SetHelp(help string)
```

#### SetLogoText
设置 logo 文本。
```go
func (c *Cmd) SetLogoText(logoText string)
```

#### SetModuleHelps
设置自定义模块帮助信息。
```go
func (c *Cmd) SetModuleHelps(moduleHelps string)
```

#### SetUsageSyntax
设置自定义命令用法。
```go
func (c *Cmd) SetUsageSyntax(usageSyntax string)
```

#### SetUseChinese
设置是否使用中文帮助信息。
```go
func (c *Cmd) SetUseChinese(useChinese bool)
```

#### SetVersion
设置版本信息。
```go
func (c *Cmd) SetVersion(version string)
```

#### ShortName
返回命令短名称。
```go
func (c *Cmd) ShortName() string
```

#### Slice
绑定字符串切片类型标志并内部注册 Flag 对象。
```go
func (c *Cmd) Slice(longName, shortName string, defValue []string, usage string) *flags.SliceFlag
```
参数依次为：长标志名、短标志、默认值、帮助说明。
返回值：字符串切片标志对象指针。

#### SliceVar
绑定字符串切片类型标志到指针并内部注册 Flag 对象。
```go
func (c *Cmd) SliceVar(f *flags.SliceFlag, longName, shortName string, defValue []string, usage string)
```
参数依次为：字符串切片标志指针、长标志名、短标志、默认值、帮助说明。

#### String
添加字符串类型标志，返回标志对象指针。
```go
func (c *Cmd) String(longName, shortName, defValue, usage string) *flags.StringFlag
```
参数依次为：长标志名、短标志、默认值、帮助说明。
返回值：字符串标志对象指针。

#### StringVar
绑定字符串类型标志到指针并内部注册 Flag 对象。
```go
func (c *Cmd) StringVar(f *flags.StringFlag, longName, shortName, defValue, usage string)
```
参数依次为：字符串标志指针、长标志名、短标志、默认值、帮助说明。

#### SubCmdMap
返回子命令映射表。
```go
func (c *Cmd) SubCmdMap() map[string]*Cmd
```

#### SubCmds
返回子命令切片。
```go
func (c *Cmd) SubCmds() []*Cmd
```

#### Time
添加时间类型标志，返回标志对象指针。
```go
func (c *Cmd) Time(longName, shortName string, defValue time.Time, usage string) *flags.TimeFlag
```
参数依次为：长标志名、短标志、默认值、帮助说明。
返回值：时间标志对象指针。

#### TimeVar
绑定时间类型标志到指针并内部注册 Flag 对象。
```go
func (c *Cmd) TimeVar(f *flags.TimeFlag, longName, shortName string, defValue time.Time, usage string)
```
参数依次为：时间标志指针、长标志名、短标志、默认值、帮助说明。

#### URL
添加 URL 类型标志，返回标志对象指针。
```go
func (c *Cmd) URL(longName, shortName string, defValue string, usage string) *flags.URLFlag
```
参数依次为：长标志名、短标志、默认值、帮助说明。
返回值：URL 标志对象指针。

#### URLVar
绑定 URL 类型标志到指针并内部注册 Flag 对象。
```go
func (c *Cmd) URLVar(f *flags.URLFlag, longName, shortName string, defValue string, usage string)
```
参数依次为：URL 标志指针、长标志名、短标志、默认值、帮助说明。

#### Uint16
添加 16 位无符号整数类型标志，返回标志对象指针。
```go
func (c *Cmd) Uint16(longName, shortName string, defValue uint16, usage string) *flags.Uint16Flag
```
参数依次为：长标志名、短标志、默认值、帮助说明。
返回值：16 位无符号整数标志对象指针。

#### Uint16Var
绑定 16 位无符号整数类型标志到指针并内部注册 Flag 对象。
```go
func (c *Cmd) Uint16Var(f *flags.Uint16Flag, longName, shortName string, defValue uint16, usage string)
```
参数依次为：16 位无符号整数标志指针、长标志名、短标志、默认值、帮助说明。

#### Uint32
添加 32 位无符号整数类型标志，返回标志对象指针。
```go
func (c *Cmd) Uint32(longName, shortName string, defValue uint32, usage string) *flags.Uint32Flag
```
参数依次为：长标志名、短标志、默认值、帮助说明。
返回值：32 位无符号整数标志对象指针。

#### Uint32Var
绑定 32 位无符号整数类型标志到指针并内部注册 Flag 对象。
```go
func (c *Cmd) Uint32Var(f *flags.Uint32Flag, longName, shortName string, defValue uint32, usage string)
```
参数依次为：32 位无符号整数标志指针、长标志名、短标志、默认值、帮助说明。

#### Uint64
添加 64 位无符号整数类型标志，返回标志对象指针。
```go
func (c *Cmd) Uint64(longName, shortName string, defValue uint64, usage string) *flags.Uint64Flag
```
参数依次为：长标志名、短标志、默认值、帮助说明。
返回值：64 位无符号整数标志对象指针。

#### Uint64Var
绑定 64 位无符号整数类型标志到指针并内部注册 Flag 对象。
```go
func (c *Cmd) Uint64Var(f *flags.Uint64Flag, longName, shortName string, defValue uint64, usage string)
```
参数依次为：64 位无符号整数标志指针、长标志名、短标志、默认值、帮助说明。

### CmdInterface
命令接口定义，封装命令行程序的核心功能，提供统一的命令管理、参数解析和帮助系统。实现类需保证线程安全，所有方法应支持并发调用。
```go
// CmdInterface 命令接口定义, 封装命令行程序的核心功能
// 提供统一的命令管理、参数解析和帮助系统
// 实现类需保证线程安全, 所有方法应支持并发调用
//
// 示例用法:
// cmd := NewCmd("app", "a", flag.ContinueOnError)
// cmd.SetDescription("示例应用程序")
// cmd.String("config", "c", "配置文件路径", "/etc/app.conf")
type CmdInterface interface {
	// 元数据操作方法
	Name() string                             // 获取命令名称
	LongName() string                         // 获取命令名称(长名称), 如"app"
	ShortName() string                        // 获取命令短名称, 如"a"
	GetDescription() string                   // 获取命令描述信息
	SetDescription(desc string)               // 设置命令描述信息, 用于帮助输出
	GetHelp() string                          // 获取自定义帮助信息
	SetHelp(help string)                      // 设置自定义帮助信息, 覆盖自动生成内容
	LoadHelp(filePath string) error           // 加载自定义帮助信息, 从文件中读取
	SetUsageSyntax(usageSyntax string)        // 设置自定义用法说明, 覆盖自动生成内容
	GetUsageSyntax() string                   // 获取自定义用法说明
	GetUseChinese() bool                      // 获取是否使用中文帮助信息
	SetUseChinese(useChinese bool)            // 设置是否使用中文帮助信息
	AddSubCmd(subCmd *Cmd)                    // 添加子命令, 子命令会继承父命令的上下文
	SubCmds() []*Cmd                          // 获取所有已注册的子命令列表
	SubCmdMap() map[string]*Cmd               // 获取所有已注册的子命令映射表
	Args() []string                           // 获取所有非标志参数(未绑定到任何标志的参数)
	Arg(i int) string                         // 获取指定索引的非标志参数, 索引越界返回空字符串
	NArg() int                                // 获取非标志参数的数量
	NFlag() int                               // 获取已解析的标志数量
	FlagExists(name string) bool              // 检查指定名称的标志是否存在(支持长/短名称)
	PrintHelp()                               // 打印命令帮助信息
	AddNote(note string)                      // 添加备注信息
	GetNotes() []string                       // 获取所有备注信息
	AddExample(e ExampleInfo)                 // 添加示例信息
	GetExamples() []ExampleInfo               // 获取所有示例信息
	SetVersion(version string)                // 设置版本信息
	GetVersion() string                       // 获取版本信息
	SetLogoText(logoText string)              // 设置logo文本
	GetLogoText() string                      // 获取logo文本
	SetModuleHelps(moduleHelps string)        // 设置自定义模块帮助信息
	GetModuleHelps() string                   // 获取自定义模块帮助信息
	SetExitOnBuiltinFlags(exit bool) *Cmd     // 设置是否在添加内置标志时退出
	SetDisableBuiltinFlags(disable bool) *Cmd // 设置是否禁用内置标志注册
	CmdExists(cmdName string) bool            // 判断命令行参数中是否存在指定标志

	// 标志解析方法
	Parse(args []string) error                // 解析命令行参数, 自动处理标志和子命令
	ParseFlagsOnly(args []string) (err error) // 仅解析标志参数, 不处理子命令
	IsParsed() bool                           // 检查是否已解析命令行参数

	// 添加标志方法
	String(longName, shortName, usage, defValue string) *flags.StringFlag                             // 添加字符串类型标志
	Int(longName, shortName, usage string, defValue int) *flags.IntFlag                               // 添加整数类型标志
	Int64(longName, shortName, usage string, defValue int64) *flags.Int64Flag                         // 添加64位整数类型标志
	Bool(longName, shortName, usage string, defValue bool) *flags.BoolFlag                            // 添加布尔类型标志
	Float64(longName, shortName, usage string, defValue float64) *flags.Float64Flag                   // 添加浮点数类型标志
	Duration(longName, shortName, usage string, defValue time.Duration) *flags.DurationFlag           // 添加时间间隔类型标志
	Enum(longName, shortName string, defValue string, usage string, options []string) *flags.EnumFlag // 添加枚举类型标志
	Slice(longName, shortName string, defValue []string, usage string) *flags.SliceFlag               // 添加字符串切片类型标志
	uint16(longName, shortName string, defValue uint16, usage string) *flags.Uint16Flag               // 添加无符号16位整型标志
	Time(longName, shortName string, defValue time.Time, usage string) *flags.TimeFlag                // 添加时间类型标志
	Map(longName, shortName string, defValue map[string]string, usage string) *flags.MapFlag          // 添加Map标志
	Path(longName, shortName string, defValue string, usage string) *flags.PathFlag                   // 添加路径标志

	// 绑定标志方法
	StringVar(f *flags.StringFlag, longName, shortName, defValue, usage string)                             // 绑定字符串标志到指定变量
	IntVar(f *flags.IntFlag, longName, shortName string, defValue int, usage string)                        // 绑定整数标志到指定变量
	Int64Var(f *flags.Int64Flag, longName, shortName string, defValue int64, usage string)                  // 绑定64位整数标志到指定变量
	BoolVar(f *flags.BoolFlag, longName, shortName string, defValue bool, usage string)                     // 绑定布尔标志到指定变量
	Float64Var(f *flags.Float64Flag, longName, shortName string, defValue float64, usage string)            // 绑定浮点数标志到指定变量
	DurationVar(f *flags.DurationFlag, longName, shortName string, defValue time.Duration, usage string)    // 绑定时间间隔类型标志到指定变量
	EnumVar(f *flags.EnumFlag, longName, shortName string, defValue string, usage string, options []string) // 绑定枚举标志到指定变量
	SliceVar(f *flags.SliceFlag, longName, shortName string, defValue []string, usage string)               // 绑定字符串切片标志到指定变量
	Uint16Var(f *flags.Uint16Flag, longName, shortName string, defValue uint16, usage string)               // 绑定无符号16位整型标志到指定变量
	TimeVar(f *flags.TimeFlag, longName, shortName string, defValue time.Time, usage string)                // 绑定时间类型标志到指定变量
	MapVar(f *flags.MapFlag, longName, shortName string, defValue map[string]string, usage string)          // 绑定字符串映射标志到指定变量
	PathVar(f *flags.PathFlag, longName, shortName string, defValue string, usage string)                   // 绑定路径标志到指定变量
}
```

### ExampleInfo
示例信息结构体，用于存储命令的使用示例，包括描述和示例内容。
```go
type ExampleInfo struct {
    Description string // 示例描述
    Usage       string // 示例使用方式
}
```

### HelpTemplate
帮助信息模板结构体。
```go
type HelpTemplate struct {
    CmdName              string // 命令名称模板
    CmdNameWithShort     string // 命令名称带短名称模板
    CmdDescription       string // 命令描述模板
    UsagePrefix          string // 用法说明前缀模板
    UsageSubCmd          string // 用法说明子命令模板
    UsageInfoWithOptions string // 带选项的用法说明信息模板
    UsageGlobalOptions   string // 全局选项部分
    OptionsHeader        string // 选项头部模板
    Option1              string // 选项模板(带短选项)
    Option2              string // 选项模板(无短选项)
    Option3              string // 选项模板(无长选项)
    OptionDefault        string // 选项模板的默认值
    SubCmdsHeader        string // 子命令头部模板
    SubCmd               string // 子命令模板
    SubCmdWithShort      string // 子命令带短名称模板
    NotesHeader          string // 注意事项头部模板
    NoteItem             string // 注意事项项模板
    DefaultNote          string // 默认注意事项
    ExamplesHeader       string // 示例信息头部模板
    ExampleItem          string // 示例信息项模板
}
```

### UserInfo
存储用户自定义信息的嵌套结构体。
```go
type UserInfo struct {
    // Has unexported fields.
}
```