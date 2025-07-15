# Package cmd

Package cmd 提供了命令行标志管理、子命令支持和帮助系统等功能，用于简化命令行工具的开发。它封装了参数解析、长短标志互斥及帮助文档生成等复杂逻辑，并提供了丰富的类型安全标志支持。

## Package 文件结构

### 文件描述

- `bash_completion.go`: Bash 补全脚本生成器
- `cmd_completion.go`: 自动补全命令的实现
- `cmd_genhelp`: 命令行帮助信息生成器
- `cmd_internal`: 包含 Cmd 的内部实现细节（不对外暴露）
- `pwsh_completion.go`: PowerShell 补全脚本生成器

## 常量定义

### Bash 补全脚本常量

```go
const (
    BashCommandTreeEntry = "cmd_tree[%s]=\"%s\"\n" // 命令树条目格式
    BashFlagParamItem    = "flag_params[%q]=%q\n"  // 标志参数项格式
    BashEnumOptions      = "enum_options[%q]=%q\n" // 枚举选项格式
)
```

### PowerShell 补全脚本常量

```go
const (
    PwshFlagParamItem = "	@{ Context = \"%s\"; Parameter = \"%s\"; ParamType = \"%s\"; ValueType = \"%s\"; Options = @(%s) }"
    PwshCmdTreeItem   = "	@{ Context = \"%s\"; Options = @(%s) }"
)
```

### Bash 补全模板

```go
const (
    BashFunctionHeader = `#!/usr/bin/env bash

# Static command tree definition - Pre-initialized outside the function
declare -A cmd_tree
cmd_tree[/]="%s"
%s

# Flag parameters definition - stores type and value type (type|valueType)
declare -A flag_params
%s

# Enum options definition - stores allowed values for enum flags
declare -A enum_options
%s

_%s() {
    local cur prev words cword context opts i arg
    COMPREPLY=()

    # Use _get_comp_words_by_ref to get completion parameters for better robustness
    if [[ -z "${_get_comp_words_by_ref}" ]]; then
        # Compatibility with older versions of Bash completion environment
        words=("${COMP_WORDS[@]}")
        cword=$COMP_CWORD
    else
        _get_comp_words_by_ref -n =: cur prev words cword
    fi

    cur="${words[cword]}"
    prev="${words[cword-1]}"

    # Find the current command context
    local context="/"
    local i
    for ((i=1; i < cword; i++)); do
        local arg="${words[i]}"
        if [[ -n "${cmd_tree[$context$arg/]}" ]]; then
            context="$context$arg/"
        fi
    done

    # Get the available options for the current context
    IFS='|' read -ra opts_arr <<< "${cmd_tree[$context]}"
    opts=$(IFS=' '; echo "${opts_arr[*]}")
    
    # Check if the previous parameter needs a value and get its type
    prev_param_type=""
    prev_value_type=""
    if [[ $cword -gt 1 ]]; then
        prev_arg="${words[cword-1]}"
        key="${context}|${prev_arg}"
        prev_param_info=${flag_params[$key]}
        IFS='|' read -r prev_param_type prev_value_type <<< "${prev_param_info}"
    fi

    # Dynamically generate completion based on parameter type
    if [[ -n "$prev_param_type" && $prev_param_type == "required" ]]; then
        case "$prev_value_type" in
            path)
                COMPREPLY=($(compgen -f -d -- "${cur}"))
                return 0
                ;;
            number)
                COMPREPLY=($(compgen -W "$(seq 1 100)" -- "${cur}"))
                return 0
                ;;
            ip)
                COMPREPLY=($(compgen -W "192.168. 10.0. 172.16." -- "${cur}"))
                return 0
                ;;
            enum)
                if [[ -z "$cur" && "$prev_value_type" == "enum" ]]; then
                    COMPREPLY=($(compgen -W "${enum_options[$key]}" -- ""))
                    return 0
                fi
                COMPREPLY=($(compgen -W "${enum_options[$key]}" -- "${cur}"))
                COMPREPLY=($(echo "${COMPREPLY[@]}" | grep -i "^${cur}"))
                return 0
                ;;
            url)
                COMPREPLY=($(compgen -W "http:// https:// ftp://" -- "${cur}"))
                return 0
                ;;
            *)
                COMPREPLY=($(compgen -W "${opts}" -- "${cur}"))
                return 0
                ;;
        esac
    fi

    COMPREPLY=($(compgen -W "${opts}" -- "${cur}"))

    return $?
}

complete -F _%s %s
`
)
```

### PowerShell 补全脚本头部模板

```go
const (
    PwshFunctionHeader = `# -------------------------- Configuration Area (Need to be modified according to actual commands) --------------------------
# Command Name
$commandName = "%s"

# 1. Command Tree
$cmdTree = @(
%s
)

# 2. Flag Parameter Definitions
$flagParams = @(
%s
)

# -----------------------------------------------------------------------------------

# -------------------------- Completion Logic Implementation ------------------------
$scriptBlock = {
    param(
        $wordToComplete,
        $commandAst,
        $cursorPosition
    )

    # 1. Parse tokens
    $tokens = $commandAst.CommandElements | ForEach-Object { $_.Extent.Text }
    $currentIndex = $tokens.Count - 1
    $prevElement = if ($currentIndex -ge 1) { $tokens[$currentIndex - 1] } else { $null }

    # 2. Calculate the current command context
    $context = "/"
    for ($i = 1; $i -le $currentIndex; $i++) {
        $elem = $tokens[$i]
        if ($elem -match '^-') { break }
        $nextContext = "$context$elem/"
        $contextMatch = $cmdTree | Where-Object { $_.Context -eq $nextContext }
        if ($contextMatch) {
            $context = $nextContext
        } else {
            break
        }
    }

    # 3. Available options in the current context
    $currentOptions = ($cmdTree | Where-Object { $_.Context -eq $context }).Options

    # 4. First complete all options (subcommands + flags) at the current level
    if ($currentOptions) {
        $matchingOptions = $currentOptions | Where-Object {
            $_ -like "$wordToComplete*"
        }
        if ($matchingOptions) {
            return $matchingOptions | ForEach-Object {
                if ($_ -match '^-') { $_ } else { "$_ " }
            }
        }
    }

    # 5. Complete flags themselves (like --ty -> --type)
    if ($wordToComplete -match '^-') {
        $flagDefs = $flagParams | Where-Object { $_.Context -eq $context }
        $flagMatches = $flagDefs | Where-Object {
            $_.Parameter -like "$wordToComplete*"
        } | ForEach-Object { $_.Parameter }
        return $flagMatches
    }

    # 6. Enum/Preset value completion
    # 6a Current token is empty → Complete all enum values of the previous flag
    if (-not $wordToComplete -and $prevElement -match '^-') {
        $paramDef = $flagParams | Where-Object {
            $_.Context -eq $context -and $_.Parameter -eq $prevElement
        }
        if ($paramDef) {
            switch ($paramDef.ValueType) {
                'enum'   { return $paramDef.Options }
                'path'   { return Get-ChildItem  -Name}
                'number' { return 1..10 | ForEach-Object { "$_" } }
                'ip'     { return @('192.168.','10.0.','172.16.','127.0.0.') }
                'url'    { return @('http://','https://','ftp://') }
                default  { return @() }
            }
        }
    }

    # 6b The current token is not empty, and the previous token is a flag that requires a value → Filter with prefix
    $flagForValue = $tokens[$currentIndex - 1]
    if ($flagForValue -match '^-' -and $currentIndex -ge 1) {
        $paramDef = $flagParams | Where-Object {
            $_.Context -eq $context -and $_.Parameter -eq $flagForValue
        }
        if ($paramDef) {
            switch ($paramDef.ValueType) {
                'path' {
                    $pattern = if ($wordToComplete) { "$wordToComplete*" } else { '*' }
                    Get-ChildItem -Name $pattern -ErrorAction SilentlyContinue
                }
                'number' { return 1..100 | Where-Object { "$_" -like "$wordToComplete*" } }
                'ip'     { return @('192.168.','10.0.','172.16.','127.0.0.') | Where-Object { $_ -like "$wordToComplete*" } }
                'enum'   { return $paramDef.Options | Where-Object { $_ -like "$wordToComplete*" } }
                'url'    { return @('http://','https://','ftp://') | Where-Object { $_ -like "$wordToComplete*" } }
                default  { return @() }
            }
        }
    }

    # 7. No match
    return @()
}

Register-ArgumentCompleter -CommandName $commandName -ScriptBlock $scriptBlock
`
)
```

## 变量定义

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

## 类型定义

### Cmd 结构体

```go
type Cmd struct {
	// 解析阶段钩子函数
	// 在标志解析完成后、子命令参数处理后调用
	//
	// 参数:
	//   - 当前命令实例
	//
	// 返回值:
	//   - error: 错误信息, 非nil时会中断解析流程
	//   - bool: 是否需要退出程序
	ParseHook func(*Cmd) (error, bool)
	// Has unexported fields.
}
```

### Cmd 命令行标志管理结构体

- 封装参数解析、长短标志互斥及帮助系统。
- 提供丰富的类型安全标志支持。
- 支持子命令和自动补全。

### QCommandLine 全局默认 Command 实例

```go
var QCommandLine *Cmd
```

## 函数定义

### NewCmd 函数

```go
func NewCmd(longName string, shortName string, errorHandling flag.ErrorHandling) *Cmd
```

- 创建新的命令实例。
- 参数：
  - `longName`: 命令长名称
  - `shortName`: 命令短名称
  - `errorHandling`: 错误处理方式
- 返回值：
  - `*Cmd`: 新的命令实例指针
- `errorHandling`可选值：
  - `flag.ContinueOnError`: 解析标志时遇到错误继续解析, 并返回错误信息
  - `flag.ExitOnError`: 解析标志时遇到错误立即退出程序, 并返回错误信息
  - `flag.PanicOnError`: 解析标志时遇到错误立即触发panic

### AddExample 方法

```go
func (c *Cmd) AddExample(e ExampleInfo)
```

- 为命令添加使用示例。
- 参数：
  - `e`: 使用示例信息
    - `Description`: 示例描述
    - `Usage`: 示例用法

### AddNote 方法

```go
func (c *Cmd) AddNote(note string)
```

- 添加备注信息到命令。
- 参数：
  - `note`: 备注信息

### AddSubCmd 方法

```go
func (c *Cmd) AddSubCmd(subCmds ...*Cmd) error
```

- 添加外部子命令到当前命令。
- 支持批量添加多个子命令。
- 遇到错误时收集所有错误并返回。
- 参数：
  - `subCmds`: 一个或多个子命令实例指针
- 返回值：
  - 错误信息, 如果所有子命令添加成功则返回nil

### Arg 方法

```go
func (c *Cmd) Arg(i int) string
```

- 获取指定索引的非标志参数。
- 参数：
  - `i`: 参数索引
- 返回值：
  - `string`: 指定索引位置的非标志参数；若索引越界，则返回空字符串

### Args 方法

```go
func (c *Cmd) Args() []string
```

- 获取非标志参数切片。
- 返回值：
  - `[]string`: 参数切片

### Bool 方法

```go
func (c *Cmd) Bool(longName, shortName string, defValue bool, usage string) *flags.BoolFlag
```

- 添加布尔类型标志。
- 返回标志对象指针。
- 参数值：
  - `longName`: string - 长标志名
  - `shortName`: string - 短标志
  - `defValue`: bool - 默认值
  - `usage`: string - 帮助说明
- 返回值：
  - `*flags.BoolFlag` - 布尔标志对象指针

### BoolVar 方法

```go
func (c *Cmd) BoolVar(f *flags.BoolFlag, longName, shortName string, defValue bool, usage string)
```

- 绑定布尔类型标志到指针并内部注册Flag对象。
- 参数值：
  - `f`: *flags.BoolFlag - 布尔标志对象指针
  - `longName`: string - 长标志名
  - `shortName`: string - 短标志
  - `defValue`: bool - 默认值
  - `usage`: string - 帮助说明

### CmdExists 方法

```go
func (c *Cmd) CmdExists(cmdName string) bool
```

- 检查子命令是否存在。
- 参数：
  - `cmdName`: 子命令名称
- 返回值：
  - `bool`: 子命令是否存在

### Duration 方法

```go
func (c *Cmd) Duration(longName, shortName string, defValue time.Duration, usage string) *flags.DurationFlag
```

- 添加时间间隔类型标志。
- 返回标志对象指针。
- 参数值：
  - `longName`: string - 长标志名
  - `shortName`: string - 短标志
  - `defValue`: time.Duration - 默认值
  - `usage`: string - 帮助说明
- 返回值：
  - `*flags.DurationFlag` - 时间间隔标志对象指针

### DurationVar 方法

```go
func (c *Cmd) DurationVar(f *flags.DurationFlag, longName, shortName string, defValue time.Duration, usage string)
```

- 绑定时间间隔类型标志到指针并内部注册Flag对象。
- 参数值：
  - `f`: *flags.DurationFlag - 时间间隔标志对象指针
  - `longName`: string - 长标志名
  - `shortName`: string - 短标志
  - `defValue`: time.Duration - 默认值
  - `usage`: string - 帮助说明

### Enum 方法

```go
func (c *Cmd) Enum(longName, shortName string, defValue string, usage string, options []string) *flags.EnumFlag
```

- 添加枚举类型标志。
- 返回标志对象指针。
- 参数值：
  - `longName`: string - 长标志名
  - `shortName`: string - 短标志
  - `defValue`: string - 默认值
  - `usage`: string - 帮助说明
  - `options`: []string - 限制该标志取值的枚举值切片
- 返回值：
  - `*flags.EnumFlag` - 枚举标志对象指针

### EnumVar 方法

```go
func (c *Cmd) EnumVar(f *flags.EnumFlag, longName, shortName string, defValue string, usage string, options []string)
```

- 绑定枚举类型标志到指针并内部注册Flag对象。
- 参数值：
  - `f`: *flags.EnumFlag - 枚举标志对象指针
  - `longName`: string - 长标志名
  - `shortName`: string - 短标志
  - `defValue`: string - 默认值
  - `usage`: string - 帮助说明
  - `options`: []string - 限制该标志取值的枚举值切片

### FlagExists 方法

```go
func (c *Cmd) FlagExists(name string) bool
```

- 检查指定名称的标志是否存在。
- 参数：
  - `name`: 标志名称
- 返回值：
  - `bool`: 标志是否存在

### FlagRegistry 方法

```go
func (c *Cmd) FlagRegistry() *flags.FlagRegistry
```

- 获取标志注册表的只读访问。
- 返回值：
  - `*flags.FlagRegistry`: 标志注册表的只读访问

### Float64 方法

```go
func (c *Cmd) Float64(longName, shortName string, defValue float64, usage string) *flags.Float64Flag
```

- 添加浮点型标志。
- 返回标志对象指针。
- 参数值：
  - `longName`: 长标志名
  - `shortName`: 短标志
  - `defValue`: 默认值
  - `usage`: 帮助说明
- 返回值：
  - `*flags.Float64Flag`: 浮点型标志对象指针

### Float64Var 方法

```go
func (c *Cmd) Float64Var(f *flags.Float64Flag, longName, shortName string, defValue float64, usage string)
```

- 绑定浮点型标志到指针并内部注册Flag对象。
- 参数值：
  - `f`: *flags.Float64Flag - 浮点型标志对象指针
  - `longName`: string - 长标志名
  - `shortName`: string - 短标志
  - `defValue`: float64 - 默认值
  - `usage`: string - 帮助说明

### GetDescription 方法

```go
func (c *Cmd) GetDescription() string
```

- 返回命令描述。
- 返回值：
  - `string`: 命令描述

### GetExamples 方法

```go
func (c *Cmd) GetExamples() []ExampleInfo
```

- 获取所有使用示例。
- 返回值：
  - 示例信息列表

### GetHelp 方法

```go
func (c *Cmd) GetHelp() string
```

- 返回命令用法帮助信息。
- 返回值：
  - `string`: 命令用法帮助信息

### GetLogoText 方法

```go
func (c *Cmd) GetLogoText() string
```

- 获取logo文本。
- 返回值：
  - `string`: logo文本

### GetModuleHelps 方法

```go
func (c *Cmd) GetModuleHelps() string
```

- 获取自定义模块帮助信息。
- 返回值：
  - `string`: 自定义模块帮助信息

### GetNotes 方法

```go
func (c *Cmd) GetNotes() []string
```

- 获取所有备注信息。
- 返回值：
  - 备注信息列表

### GetUsageSyntax 方法

```go
func (c *Cmd) GetUsageSyntax() string
```

- 获取自定义命令用法。
- 返回值：
  - `string`: 自定义命令用法

### GetUseChinese 方法

```go
func (c *Cmd) GetUseChinese() bool
```

- 获取是否使用中文帮助信息。
- 返回值：
  - `bool`: 是否使用中文帮助信息

### GetVersion 方法

```go
func (c *Cmd) GetVersion() string
```

- 获取版本信息。
- 返回值：
  - `string`: 版本信息

### IP4 方法

```go
func (c *Cmd) IP4(longName, shortName string, defValue string, usage string) *flags.IP4Flag
```

- 添加IPv4地址类型标志。
- 返回标志对象指针。
- 参数值：
  - `longName`: 长标志名
  - `shortName`: 短标志名
  - `defValue`: 默认值
  - `usage`: 帮助说明
- 返回值：
  - `*flags.IP4Flag`: IPv4地址标志对象指针

### IP4Var 方法

```go
func (c *Cmd) IP4Var(f *flags.IP4Flag, longName, shortName string, defValue string, usage string)
```

- 绑定IPv4地址类型标志到指针并内部注册Flag对象。
- 参数值：
  - `f`: IPv4标志指针
  - `longName`: 长标志名
  - `shortName`: 短标志名
  - `defValue`: 默认值
  - `usage`: 帮助说明

### IP6 方法

```go
func (c *Cmd) IP6(longName, shortName string, defValue string, usage string) *flags.IP6Flag
```

- 添加IPv6地址类型标志。
- 返回标志对象指针。
- 参数值：
  - `longName`: 长标志名
  - `shortName`: 短标志名
  - `defValue`: 默认值
  - `usage`: 帮助说明
- 返回值：
  - `*flags.IP6Flag`: IPv6地址标志对象指针

### IP6Var 方法

```go
func (c *Cmd) IP6Var(f *flags.IP6Flag, longName, shortName string, defValue string, usage string)
```

- 绑定IPv6地址类型标志到指针并内部注册Flag对象。
- 参数值：
  - `f`: IPv6标志指针
  - `longName`: 长标志名
  - `shortName`: 短标志名
  - `defValue`: 默认值
  - `usage`: 帮助说明

### Int 方法

```go
func (c *Cmd) Int(longName, shortName string, defValue int, usage string) *flags.IntFlag
```

- 添加整数类型标志。
- 返回标志对象指针。
- 参数值：
  - `longName`: 长标志名
  - `shortName`: 短标志名
  - `defValue`: 默认值
  - `usage`: 帮助说明
- 返回值：
  - `*flags.IntFlag`: 整数标志对象指针

### Int64 方法

```go
func (c *Cmd) Int64(longName, shortName string, defValue int64, usage string) *flags.Int64Flag
```

- 添加64位整数类型标志。
- 返回标志对象指针。
- 参数值：
  - `longName`: 长标志名
  - `shortName`: 短标志名
  - `defValue`: 默认值
  - `usage`: 帮助说明
- 返回值：
  - `*flags.Int64Flag`: 64位整数标志对象指针

### Int64Var 方法

```go
func (c *Cmd) Int64Var(f *flags.Int64Flag, longName, shortName string, defValue int64, usage string)
```

- 绑定64位整数类型标志到指针并内部注册Flag对象。
- 参数值：
  - `f`: 64位整数标志指针
  - `longName`: 长标志名
  - `shortName`: 短标志名
  - `defValue`: 默认值
  - `usage`: 帮助说明

### IntVar 方法

```go
func (c *Cmd) IntVar(f *flags.IntFlag, longName, shortName string, defValue int, usage string)
```

- 绑定整数类型标志到指针并内部注册Flag对象。
- 参数值：
  - `f`: 整数标志指针
  - `longName`: 长标志名
  - `shortName`: 短标志名
  - `defValue`: 默认值
  - `usage`: 帮助说明

### IsParsed 方法

```go
func (c *Cmd) IsParsed() bool
```

- 检查命令是否已完成解析。
- 返回值：
  - `bool`: 解析状态, true表示已解析(无论成功失败), false表示未解析

### LoadHelp 方法

```go
func (c *Cmd) LoadHelp(filePath string) error
```

- 从指定文件加载帮助信息。
- 参数：
  - `filePath`: 帮助信息文件路径
- 返回值：
  - `error`: 如果文件不存在或读取文件失败，则返回错误信息

### LongName 方法

```go
func (c *Cmd) LongName() string
```

- 返回命令长名称。
- 返回值：
  - `string`: 命令长名称

### Map 方法

```go
func (c *Cmd) Map(longName, shortName string, defValue map[string]string, usage string) *flags.MapFlag
```

- 添加键值对类型标志。
- 返回标志对象指针。
- 参数值：
  - `longName`: 长标志名
  - `shortName`: 短标志名
  - `defValue`: 默认值
  - `usage`: 帮助说明
- 返回值：
  - `*flags.MapFlag`: 键值对标志对象指针

### MapVar 方法

```go
func (c *Cmd) MapVar(f *flags.MapFlag, longName, shortName string, defValue map[string]string, usage string)
```

- 绑定键值对类型标志到指针并内部注册Flag对象。
- 参数值：
  - `f`: 键值对标志指针
  - `longName`: 长标志名
  - `shortName`: 短标志名
  - `defValue`: 默认值
  - `usage`: 帮助说明

### NArg 方法

```go
func (c *Cmd) NArg() int
```

- 获取非标志参数的数量。
- 返回值：
  - `int`: 参数数量

### NFlag 方法

```go
func (c *Cmd) NFlag() int
```

- 获取标志的数量。
- 返回值：
  - `int`: 标志数量

### Name 方法

```go
func (c *Cmd) Name() string
```

- 获取命令名称。
- 返回值：
  - 优先返回长名称, 如果长名称不存在则返回短名称

### Parse 方法

```go
func (c *Cmd) Parse(args []string) (err error)
```

- 完整解析命令行参数(含子命令处理)。
- 主要功能：
  1. 解析当前命令的长短标志及内置标志
  2. 自动检测并解析子命令及其参数(若存在)
  3. 验证枚举类型标志的有效性
- 参数：
  - `args`: 原始命令行参数切片(包含可能的子命令及参数)
- 返回值：
  - 解析过程中遇到的错误(如标志格式错误、子命令解析失败等)
- 注意事项：
  - 每个Cmd实例仅会被解析一次(线程安全)
  - 若检测到子命令, 会将剩余参数传递给子命令的Parse方法
  - 处理内置标志执行逻辑

### ParseFlagsOnly 方法

```go
func (c *Cmd) ParseFlagsOnly(args []string) (err error)
```

- 仅解析当前命令的标志参数(忽略子命令)。
- 主要功能：
  1. 解析当前命令的长短标志及内置标志
  2. 验证枚举类型标志的有效性
  3. 明确忽略所有子命令及后续参数
- 参数：
  - `args`: 原始命令行参数切片(子命令及后续参数会被忽略)
- 返回值：
  - 解析过程中遇到的错误(如标志格式错误等)
- 注意事项：
  - 每个Cmd实例仅会被解析一次(线程安全)
  - 不会处理任何子命令, 所有参数均视为当前命令的标志或位置参数
  - 处理内置标志逻辑

### Path 方法

```go
func (c *Cmd) Path(longName, shortName string, defValue string, usage string) *flags.PathFlag
```

- 添加路径类型标志。
- 返回标志对象指针。
- 参数值：
  - `longName`: 长标志名
  - `shortName`: 短标志名
  - `defValue`: 默认值
  - `usage`: 帮助说明
- 返回值：
  - `*flags.PathFlag`: 路径标志对象指针

### PathVar 方法

```go
func (c *Cmd) PathVar(f *flags.PathFlag, longName, shortName string, defValue string, usage string)
```

- 绑定路径类型标志到指针并内部注册Flag对象。
- 参数值：
  - `f`: 路径标志指针
  - `longName`: 长标志名
  - `shortName`: 短标志名
  - `defValue`: 默认值
  - `usage`: 帮助说明

### PrintHelp 方法

```go
func (c *Cmd) PrintHelp()
```

- 打印命令的帮助信息。
- 优先打印用户的帮助信息，否则自动生成帮助信息。
- 注意事项：
  - 打印帮助信息时, 不会自动退出程序

### SetDescription 方法

```go
func (c *Cmd) SetDescription(desc string)
```

- 设置命令描述。
- 参数：
  - `desc`: 命令描述

### SetEnableCompletion 方法

```go
func (c *Cmd) SetEnableCompletion(enable bool)
```

- 设置是否启用自动补全。
- 只能在根命令上启用。
- 参数：
  - `enable`: true表示启用补全, false表示禁用

### SetExitOnBuiltinFlags 方法

```go
func (c *Cmd) SetExitOnBuiltinFlags(exit bool)
```

- 设置是否在解析内置参数时退出。
- 默认情况下为true, 当解析到内置参数时, QFlag将退出程序。
- 参数：
  - `exit`: 是否退出

### SetHelp 方法

```go
func (c *Cmd) SetHelp(help string)
```

- 设置用户自定义命令帮助信息。
- 参数：
  - `help`: 用户自定义命令帮助信息

### SetLogoText 方法

```go
func (c *Cmd) SetLogoText(logoText string)
```

- 设置logo文本。
- 参数：
  - `logoText`: logo文本

### SetModuleHelps 方法

```go
func (c *Cmd) SetModuleHelps(moduleHelps string)
```

- 设置自定义模块帮助信息。
- 参数：
  - `moduleHelps`: 自定义模块帮助信息

### SetUsageSyntax 方法

```go
func (c *Cmd) SetUsageSyntax(usageSyntax string)
```

- 设置自定义命令用法。
- 参数：
  - `usageSyntax`: 自定义命令用法

### SetUseChinese 方法

```go
func (c *Cmd) SetUseChinese(useChinese bool)
```

- 设置是否使用中文帮助信息。
- 参数：
  - `useChinese`: 是否使用中文帮助信息

### SetVersion 方法

```go
func (c *Cmd) SetVersion(version string)
```

- 设置版本信息。
- 参数：
  - `version`: 版本信息

### ShortName 方法

```go
func (c *Cmd) ShortName() string
```

- 返回命令短名称。
- 返回值：
  - `string`: 命令短名称

### Slice 方法

```go
func (c *Cmd) Slice(longName, shortName string, defValue []string, usage string) *flags.SliceFlag
```

- 绑定字符串切片类型标志并内部注册Flag对象。
- 参数值：
  - `longName`: 长标志名
  - `shortName`: 短标志名
  - `defValue`: 默认值
  - `usage`: 帮助说明
- 返回值：
  - `*flags.SliceFlag`: 字符串切片标志对象指针

### SliceVar 方法

```go
func (c *Cmd) SliceVar(f *flags.SliceFlag, longName, shortName string, defValue []string, usage string)
```

- 绑定字符串切片类型标志到指针并内部注册Flag对象。
- 参数值：
  - `f`: 字符串切片标志指针
  - `longName`: 长标志名
  - `shortName`: 短标志名
  - `defValue`: 默认值
  - `usage`: 帮助说明

### String 方法

```go
func (c *Cmd) String(longName, shortName, defValue, usage string) *flags.StringFlag
```

- 添加字符串类型标志。
- 返回标志对象指针。
- 参数值：
  - `longName`: 长标志名
  - `shortName`: 短标志名
  - `defValue`: 默认值
  - `usage`: 帮助说明
- 返回值：
  - `*flags.StringFlag`: 字符串标志对象指针

### StringVar 方法

```go
func (c *Cmd) StringVar(f *flags.StringFlag, longName, shortName, defValue, usage string)
```

- 绑定字符串类型标志到指针并内部注册Flag对象。
- 参数值：
  - `f`: 字符串标志指针
  - `longName`: 长标志名
  - `shortName`: 短标志名
  - `defValue`: 默认值
  - `usage`: 帮助说明

### SubCmdMap 方法

```go
func (c *Cmd) SubCmdMap() map[string]*Cmd
```

- 返回子命令映射表。
- 返回值：
  - `map[string]*Cmd`: 子命令映射表

### SubCmds 方法

```go
func (c *Cmd) SubCmds() []*Cmd
```

- 返回子命令切片。
- 返回值：
  - `[]*Cmd`: 子命令切片

### Time 方法

```go
func (c *Cmd) Time(longName, shortName string, defValue time.Time, usage string) *flags.TimeFlag
```

- 添加时间类型标志。
- 返回标志对象指针。
- 参数值：
  - `longName`: 长标志名
  - `shortName`: 短标志名
  - `defValue`: 默认值
  - `usage`: 帮助说明
- 返回值：
  - `*flags.TimeFlag`: 时间标志对象指针

### TimeVar 方法

```go
func (c *Cmd) TimeVar(f *flags.TimeFlag, longName, shortName string, defValue time.Time, usage string)
```

- 绑定时间类型标志到指针并内部注册Flag对象。
- 参数值：
  - `f`: 时间标志指针
  - `longName`: 长标志名
  - `shortName`: 短标志名
  - `defValue`: 默认值
  - `usage`: 帮助说明

### URL 方法

```go
func (c *Cmd) URL(longName, shortName string, defValue string, usage string) *flags.URLFlag
```

- 添加URL类型标志。
- 返回标志对象指针。
- 参数值：
  - `longName`: 长标志名
  - `shortName`: 短标志名
  - `defValue`: 默认值
  - `usage`: 帮助说明
- 返回值：
  - `*flags.URLFlag`: URL标志对象指针

### URLVar 方法

```go
func (c *Cmd) URLVar(f *flags.URLFlag, longName, shortName string, defValue string, usage string)
```

- 绑定URL类型标志到指针并内部注册Flag对象。
- 参数值：
  - `f`: URL标志对象指针
  - `longName`: 长标志名
  - `shortName`: 短标志名
  - `defValue`: 默认值
  - `usage`: 帮助说明

### Uint16 方法

```go
func (c *Cmd) Uint16(longName, shortName string, defValue uint16, usage string) *flags.Uint16Flag
```

- 添加16位无符号整数类型标志。
- 返回标志对象指针。
- 参数值：
  - `longName`: 长标志名
  - `shortName`: 短标志名
  - `defValue`: 默认值
  - `usage`: 帮助说明
- 返回值：
  - `*flags.Uint16Flag`: 16位无符号整数标志对象指针

### Uint16Var 方法

```go
func (c *Cmd) Uint16Var(f *flags.Uint16Flag, longName, shortName string, defValue uint16, usage string)
```

- 绑定16位无符号整数类型标志到指针并内部注册Flag对象。
- 参数值：
  - `f`: 16位无符号整数标志指针
  - `longName`: 长标志名
  - `shortName`: 短标志名
  - `defValue`: 默认值
  - `usage`: 帮助说明

### Uint32 方法

```go
func (c *Cmd) Uint32(longName, shortName string, defValue uint32, usage string) *flags.Uint32Flag
```

- 添加32位无符号整数类型标志。
- 返回标志对象指针。
- 参数值：
  - `longName`: 长标志名
  - `shortName`: 短标志名
  - `defValue`: 默认值
  - `usage`: 帮助说明
- 返回值：
  - `*flags.Uint32Flag`: 32位无符号整数标志对象指针

### Uint32Var 方法

```go
func (c *Cmd) Uint32Var(f *flags.Uint32Flag, longName, shortName string, defValue uint32, usage string)
```

- 绑定32位无符号整数类型标志到指针并内部注册Flag对象。
- 参数值：
  - `f`: 32位无符号整数标志指针
  - `longName`: 长标志名
  - `shortName`: 短标志名
  - `defValue`: 默认值
  - `usage`: 帮助说明

### Uint64 方法

```go
func (c *Cmd) Uint64(longName, shortName string, defValue uint64, usage string) *flags.Uint64Flag
```

- 添加64位无符号整数类型标志。
- 返回标志对象指针。
- 参数值：
  - `longName`: 长标志名
  - `shortName`: 短标志名
  - `defValue`: 默认值
  - `usage`: 帮助说明
- 返回值：
  - `*flags.Uint64Flag`: 64位无符号整数标志对象指针

### Uint64Var 方法

```go
func (c *Cmd) Uint64Var(f *flags.Uint64Flag, longName, shortName string, defValue uint64, usage string)
```

- 绑定64位无符号整数类型标志到指针并内部注册Flag对象。
- 参数值：
  - `f`: 64位无符号整数标志指针
  - `longName`: 长标志名
  - `shortName`: 短标志名
  - `defValue`: 默认值
  - `usage`: 帮助说明

## 结构体定义

### ExampleInfo 结构体

```go
type ExampleInfo struct {
	Description string // 示例描述
	Usage       string // 示例使用方式
}
```

- 用于存储命令的使用示例，包括描述和示例内容。

### FlagParam 结构体

```go
type FlagParam struct {
	CommandPath string   // 命令路径，如 "/cmd/subcmd"
	Name        string   // 标志名称(保留原始大小写)
	Type        string   // 参数需求类型: "required"|"optional"|"none"
	ValueType   string   // 参数值类型: "path"|"string"|"number"|"enum"|"bool"等
	EnumOptions []string // 枚举类型的可选值列表
}
```

- 表示标志参数及其需求类型和值类型。

### HelpTemplate 结构体

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

- 帮助信息模板结构体，用于定义命令帮助信息的格式和内容。