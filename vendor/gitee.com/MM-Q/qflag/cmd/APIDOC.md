# Package cmd

`cmd` 包提供了命令行标志管理和解析的功能，支持多种标志类型和自定义帮助信息。

## VARIABLES

### 中文模板实例

```markdown
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

### 英文模板实例

```markdown
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

### 保持兼容 API

```markdown
var NewCmd = NewCommand
```

## FUNCTIONS

### GetExecutablePath

获取程序的绝对安装路径，如果无法通过 `os.Executable` 获取路径，则使用 `os.Args[0]` 作为替代。

```markdown
func GetExecutablePath() string
```

## TYPES

### Cmd

命令行标志管理结构体，封装参数解析、长短标志互斥及帮助系统。

```markdown
type Cmd struct {
    fs                  *flag.FlagSet
    flagRegistry        *flags.FlagRegistry
    initFlagBound       bool
    initFlagOnce        sync.Once
    parseOnce           sync.Once
    subCmds             []*Cmd
    parentCmd           *Cmd
    args                []string
    addMu               sync.Mutex
    setMu               sync.Mutex
    builtinFlagNameMap  sync.Map
    userInfo            UserInfo
    helpFlag            *flags.BoolFlag
    showInstallPathFlag *flags.BoolFlag
    versionFlag         *flags.BoolFlag
}
```

### CmdInterface

命令接口定义，封装命令行程序的核心功能。

```markdown
type CmdInterface interface {
    LongName() string
    ShortName() string
    GetDescription() string
    SetDescription(desc string)
    GetHelp() string
    SetHelp(help string)
    LoadHelp(filePath string) error
    SetUsageSyntax(usageSyntax string)
    GetUsageSyntax() string
    GetUseChinese() bool
    SetUseChinese(useChinese bool)
    AddSubCmd(subCmd *Cmd)
    SubCmds() []*Cmd
    Parse(args []string) error
    ParseFlagsOnly(args []string) error
    Args() []string
    Arg(i int) string
    NArg() int
    NFlag() int
    FlagExists(name string) bool
    PrintHelp()
    AddNote(note string)
    GetNotes() []string
    AddExample(e ExampleInfo)
    GetExamples() []ExampleInfo
    SetVersion(version string)
    GetVersion() string
    SetLogoText(logoText string)
    GetLogoText() string
    SetModuleHelps(moduleHelps string)
    GetModuleHelps() string
    String(longName, shortName, usage, defValue string) *flags.StringFlag
    Int(longName, shortName, usage string, defValue int) *flags.IntFlag
    Int64(longName, shortName, usage string, defValue int64) *flags.Int64Flag
    Bool(longName, shortName, usage string, defValue bool) *flags.BoolFlag
    Float64(longName, shortName, usage string, defValue float64) *flags.Float64Flag
    Duration(longName, shortName, usage string, defValue time.Duration) *flags.DurationFlag
    Enum(longName, shortName string, defValue string, usage string, options []string) *flags.EnumFlag
    Slice(longName, shortName string, defValue []string, usage string) *flags.SliceFlag
    Uint16(longName, shortName string, defValue uint16, usage string) *flags.Uint16Flag
    Time(longName, shortName string, defValue time.Time, usage string) *flags.TimeFlag
    Map(longName, shortName string, defValue map[string]string, usage string) *flags.MapFlag
    Path(longName, shortName string, defValue string, usage string) *flags.PathFlag
    StringVar(f *flags.StringFlag, longName, shortName, defValue, usage string)
    IntVar(f *flags.IntFlag, longName, shortName string, defValue int, usage string)
    Int64Var(f *flags.Int64Flag, longName, shortName string, defValue int64, usage string)
    BoolVar(f *flags.BoolFlag, longName, shortName string, defValue bool, usage string)
    Float64Var(f *flags.Float64Flag, longName, shortName string, defValue float64, usage string)
    DurationVar(f *flags.DurationFlag, longName, shortName string, defValue time.Duration, usage string)
    EnumVar(f *flags.EnumFlag, longName, shortName string, defValue string, usage string, options []string)
    SliceVar(f *flags.SliceFlag, longName, shortName string, defValue []string, usage string)
    Uint16Var(f *flags.Uint16Flag, longName, shortName string, defValue uint16, usage string)
    TimeVar(f *flags.TimeFlag, longName, shortName string, defValue time.Time, usage string)
    MapVar(f *flags.MapFlag, longName, shortName string, defValue map[string]string, usage string)
    PathVar(f *flags.PathFlag, longName, shortName string, defValue string, usage string)
}
```

### ExampleInfo

示例信息结构体，用于存储命令的使用示例。

```markdown
type ExampleInfo struct {
    Description string
    Usage       string
}
```

### HelpTemplate

帮助信息模板结构体。

```markdown
type HelpTemplate struct {
    CmdName              string
    CmdNameWithShort     string
    CmdDescription       string
    UsagePrefix          string
    UsageSubCmd          string
    UsageInfoWithOptions string
    UsageGlobalOptions   string
    OptionsHeader        string
    Option1              string
    Option2              string
    Option3              string
    OptionDefault        string
    SubCmdsHeader        string
    SubCmd               string
    SubCmdWithShort      string
    NotesHeader          string
    NoteItem             string
    DefaultNote          string
    ExamplesHeader       string
    ExampleItem          string
}
```

### UserInfo

存储用户自定义信息的嵌套结构体。

```markdown
type UserInfo struct {
    longName       string
    shortName      string
    version        string
    description    string
    help           string
    usageSyntax    string
    moduleHelps    string
    logoText       string
    notes          []string
    examples       []ExampleInfo
    useChinese     bool
}
```