# Package types 

```go
import "gitee.com/MM-Q/qflag/internal/types"
```

Package types 定义了qflag项目的核心类型和接口

types 包提供了整个项目的基础类型定义, 包括:
  - 标志类型和接口定义
  - 命令接口定义
  - 注册表接口定义
  - 智能纠错错误类型（UnknownSubcommandError, UnknownFlagError）

这些类型和接口构成了整个框架的核心抽象层,  为具体的实现提供了统一的规范和契约。

---

## CONSTANTS

### 内置标志名称常量

```go
const (
    // HelpFlagName 帮助标志名称
    HelpFlagName = "help"

    // HelpFlagShortName 帮助标志短名称
    HelpFlagShortName = "h"

    // VersionFlagName 版本标志名称
    VersionFlagName = "version"

    // VersionFlagShortName 版本标志短名称
    VersionFlagShortName = "v"

    // CompletionFlagName 补全标志名称
    CompletionFlagName = "completion"
)
```

### Shell 类型常量

```go
const (
    // BashShell bash shell
    BashShell = "bash"

    // PwshShell pwsh shell
    PwshShell = "pwsh"

    // PowershellShell powershell shell
    PowershellShell = "powershell"
)
```

### 补全脚本生成相关常量定义

```go
const (
    // DefaultFlagParamsCapacity 预估的标志参数初始容量
    // 基于常见CLI工具分析, 大多数工具的标志数量在100-500之间
    DefaultFlagParamsCapacity = 256

    // NamesPerItem 每个标志/命令的名称数量(长名+短名)
    NamesPerItem = 2

    // MaxTraverseDepth 命令树遍历的最大深度限制
    // 防止循环引用导致的无限递归, 一般CLI工具很少超过20层
    MaxTraverseDepth = 50
)
```

### 帮助信息标题 - 中文

```go
const (
    HelpNameCN     = "名称:\n"
    HelpDescCN     = "\n描述:\n"
    HelpUsageCN    = "\n用法:\n"
    HelpOptionsCN  = "\n选项:\n"
    HelpSubCmdsCN  = "\n子命令:\n"
    HelpExamplesCN = "\n示例:\n"
    HelpNotesCN    = "\n注意:\n"
)
```

### 帮助信息标题 - 英文

```go
const (
    HelpNameEN     = "Name:\n"
    HelpDescEN     = "\nDesc:\n"
    HelpUsageEN    = "\nUsage:\n"
    HelpOptionsEN  = "\nOptions:\n"
    HelpSubCmdsEN  = "\nSubcmds:\n"
    HelpExamplesEN = "\nExamples:\n"
    HelpNotesEN    = "\nNotes:\n"
)
```

### 存储单位常量

```go
const (
    // 十进制单位 (KB, MB, GB等) 
    KB = 1000
    MB = 1000 * KB
    GB = 1000 * MB
    TB = 1000 * GB
    PB = 1000 * TB

    // 二进制单位 (KIB, MIB, GIB等) 
    KIB = 1024
    MIB = 1024 * KIB
    GIB = 1024 * MIB
    TIB = 1024 * GIB
    PIB = 1024 * TIB
)
```

### 帮助格式常量

```go
const HelpOptionSubCmdSpace = "      "
    选项和子命令中间空格, 缩进6个空格

const HelpPrefix = "  "
    统一的前缀, 缩进两个空格
```

### __complete 子命令指令常量

```go
const (
    // InstructionFuzzy 模糊匹配指令
    // 用法: __complete fuzzy <模式> <候选1> [候选2] ...
    // 输出: 每行一个匹配结果（按匹配质量降序）
    InstructionFuzzy = "fuzzy"

    // InstructionContext 上下文计算指令
    // 用法: __complete context <arg0> [arg1] ...
    // 输出: 上下文路径，如 "/server/start/"
    InstructionContext = "context"

    // InstructionCandidates 候选选项获取指令
    // 用法: __complete candidates <context>
    // 输出: 空格分隔的候选选项列表
    InstructionCandidates = "candidates"

    // InstructionEnum 枚举值获取指令
    // 用法: __complete enum <context> <flag-name>
    // 输出: 空格分隔的枚举值列表
    InstructionEnum = "enum"

    // InstructionAll 统一补全指令
    // 用法: __complete all <cur> <prev> [cmd_args...]
    // 输出: 多行格式，包含 CONTEXT, CUR, PREV, CANDIDATES, ENUM, MATCHES, IS_FLAG
    InstructionAll = "all"
)
```

### 补全命令名称常量

```go
const CompleteCmdName = "__complete"
    // 补全子命令名称，用于 Shell 自动补全脚本
```

### 补全安装相关常量

```go
const (
    // CompletionsDirName 补全脚本存放目录名
    CompletionsDirName = ".qflag_completions"

    // CompletionScriptComment 补全脚本注释模板
    CompletionScriptComment = "# qflag completion for %s\n"

    // PwshCompletionScriptExt PowerShell 补全脚本扩展名
    PwshCompletionScriptExt = ".ps1"

    // BashCompletionScriptExt Bash 补全脚本扩展名
    BashCompletionScriptExt = ".sh"

    // PwshProfileDirWindows Windows PowerShell 配置文件目录
    PwshProfileDirWindows = "Documents/PowerShell"

    // PwshProfileFileName PowerShell 配置文件名
    PwshProfileFileName = "Microsoft.PowerShell_profile.ps1"

    // PwshProfileDirUnix Unix PowerShell 配置文件目录
    PwshProfileDirUnix = ".config/powershell"

    // BashProfileFileNameDarwin macOS Bash 配置文件名
    BashProfileFileNameDarwin = ".bash_profile"

    // BashProfileFileNameLinux Linux Bash 配置文件名
    BashProfileFileNameLinux = ".bashrc"
)
```

### 补全加载命令模板

```go
const (
    // PwshLoadCommandTemplate PowerShell 加载命令模板
    // 参数: 程序名（3次）、脚本路径（1次）
    PwshLoadCommandTemplate = "$__qflag_comp_%s = '%s'; if (Test-Path $__qflag_comp_%s) { . $__qflag_comp_%s }"

    // BashLoadCommandTemplate Bash 加载命令模板
    // 参数: 脚本路径（2次）
    BashLoadCommandTemplate = "[ -f '%s' ] && source '%s'"
)
```

### 补全安装成功信息 - 中文

```go
const (
    // InstallSuccessScriptPathCN 脚本安装路径提示（中文）
    InstallSuccessScriptPathCN = "✓ 补全脚本已安装: %s"

    // InstallSuccessProfilePathCN 配置文件路径提示（中文）
    InstallSuccessProfilePathCN = "✓ 加载命令已添加到: %s"

    // InstallSuccessHintCN 重启提示（中文）
    InstallSuccessHintCN = "\n请重启终端或运行以下命令启用补全:"
)
```

### 补全安装成功信息 - 英文

```go
const (
    // InstallSuccessScriptPathEN 脚本安装路径提示（英文）
    InstallSuccessScriptPathEN = "✓ Completion script installed: %s"

    // InstallSuccessProfilePathEN 配置文件路径提示（英文）
    InstallSuccessProfilePathEN = "✓ Load command added to: %s"

    // InstallSuccessHintEN 重启提示（英文）
    InstallSuccessHintEN = "\nPlease restart your terminal or run the following command to enable completions:"
)
```

### 补全执行命令（Shell 命令本身不需要翻译）

```go
const (
    // InstallSuccessBashCmd Bash 执行命令
    InstallSuccessBashCmd = "  source %s"

    // InstallSuccessPwshCmd PowerShell 执行命令
    InstallSuccessPwshCmd = "  . %s"
)
```

---

## VARIABLES

### 常见时间格式常量

```go
var (
    // RFC3339 RFC3339 格式 (2006-01-02T15:04:05Z07:00)
    TimeFormatRFC3339 = time.RFC3339

    // RFC3339Nano RFC3339 纳秒格式 (2006-01-02T15:04:05.999999999Z07:00)
    TimeFormatRFC3339Nano = time.RFC3339Nano

    // RFC1123 RFC1123 格式 (Mon, 02 Jan 2006 15:04:05 MST)
    TimeFormatRFC1123 = time.RFC1123

    // RFC1123Z RFC1123 带时区格式 (Mon, 02 Jan 2006 15:04:05 -0700)
    TimeFormatRFC1123Z = time.RFC1123Z

    // RFC822 RFC822 格式 (02 Jan 06 15:04 MST)
    TimeFormatRFC822 = time.RFC822

    // RFC822Z RFC822 带时区格式 (02 Jan 06 15:04 -0700)
    TimeFormatRFC822Z = time.RFC822Z

    // Kitchen 厨房格式 (3:04PM)
    TimeFormatKitchen = time.Kitchen

    // Stamp 简单时间戳格式 (Jan _2 15:04:05)
    TimeFormatStamp = time.Stamp

    // StampMilli 毫秒时间戳格式 (Jan _2 15:04:05.000)
    TimeFormatStampMilli = time.StampMilli

    // StampMicro 微秒时间戳格式 (Jan _2 15:04:05.000000)
    TimeFormatStampMicro = time.StampMicro

    // StampNano 纳秒时间戳格式 (Jan _2 15:04:05.000000000)
    TimeFormatStampNano = time.StampNano

    // DateOnly 日期格式 (2006-01-02)
    TimeFormatDateOnly = "2006-01-02"

    // TimeOnly 时间格式 (15:04:05)
    TimeFormatTimeOnly = "15:04:05"

    // DateTime 日期时间格式 (2006-01-02 15:04:05)
    TimeFormatDateTime = "2006-01-02 15:04:05"

    // DateTimeSlash 斜杠分隔的日期时间格式 (01/02/2006 15:04:05)
    TimeFormatDateTimeSlash = "01/02/2006 15:04:05"

    // DateTimeCompact 紧凑日期时间格式 (20060102150405)
    TimeFormatDateTimeCompact = "20060102150405"

    // ISO8601 ISO8601 格式 (2006-01-02T15:04:05Z)
    TimeFormatISO8601 = "2006-01-02T15:04:05Z"

    // ISO8601Nano ISO8601 纳秒格式 (2006-01-02T15:04:05.999999999Z)
    TimeFormatISO8601Nano = "2006-01-02T15:04:05.999999999Z"
)
```

### CommonTimeFormats 常见时间格式列表, 按优先级排序

```go
var CommonTimeFormats = []string{
    TimeFormatRFC3339,
    TimeFormatRFC3339Nano,
    TimeFormatISO8601,
    TimeFormatISO8601Nano,
    TimeFormatDateTime,
    TimeFormatDateOnly,
    TimeFormatTimeOnly,
    TimeFormatRFC1123,
    TimeFormatRFC1123Z,
    TimeFormatDateTimeSlash,
    TimeFormatStamp,
    TimeFormatStampMilli,
    TimeFormatStampMicro,
    TimeFormatStampNano,
    TimeFormatRFC822,
    TimeFormatRFC822Z,
    TimeFormatKitchen,
    TimeFormatDateTimeCompact,
}
```

### 内置补全示例信息 - Windows 中文

```go
var HelpCompletionExampleWinCN = map[string]string{
    "临时启用": fmt.Sprintf("%s --completion pwsh | Out-String | Invoke-Expression", programName),
}
```

### 内置补全示例信息 - Windows 英文

```go
var HelpCompletionExampleWinEN = map[string]string{
    "Temporary": fmt.Sprintf("%s --completion pwsh | Out-String | Invoke-Expression", programName),
}
```

### 内置补全示例信息 - Unix-like 系统中文（Linux 和 macOS 共用）

```go
var HelpCompletionExampleUnixCN = map[string]string{
    "临时启用": fmt.Sprintf("source <(%s --completion bash)", programName),
}
```

### 内置补全示例信息 - Unix-like 系统英文（Linux 和 macOS 共用）

```go
var HelpCompletionExampleUnixEN = map[string]string{
    "Temporary": fmt.Sprintf("source <(%s --completion bash)", programName),
}
```

### SupportedShells Shell切片, 用于存储支持的Shell类型

```go
var SupportedShells = []string{
    BashShell,
    PwshShell,
    PowershellShell,
}
```

---

## FUNCTIONS

### func CurrentShell() string

```go
func CurrentShell() string
```

CurrentShell 根据当前平台返回对应的Shell类型

**返回值:**
  - string: 当前Shell类型

---

### func GetCompletionExample() map[string]string

```go
func GetCompletionExample() map[string]string
```

GetCompletionExample 获取当前平台的补全示例信息（中文）

**返回值:**
  - map[string]string: 包含补全示例信息的映射

**功能说明:**
  - 根据当前运行的操作系统返回对应的中文补全示例
  - 支持 Windows、Linux 和 macOS 平台
  - 提供临时启用和永久启用两种方式的示例

---

### func GetCompletionExampleEN() map[string]string

```go
func GetCompletionExampleEN() map[string]string
```

GetCompletionExampleEN 获取当前平台的补全示例信息（英文）

**返回值:**
  - map[string]string: 包含补全示例信息的映射

**功能说明:**
  - 根据当前运行的操作系统返回对应的英文补全示例
  - 支持 Windows、Linux 和 macOS 平台
  - 提供临时启用和永久启用两种方式的示例

---

### func GetInstallCompletionExample() map[string]string

```go
func GetInstallCompletionExample() map[string]string
```

GetInstallCompletionExample 获取当前平台的安装补全示例信息（中文）

**返回值:**
  - map[string]string: 包含安装补全示例信息的映射

**功能说明:**
  - 根据当前运行的操作系统返回对应的中文安装补全示例
  - Windows 使用 pwsh，其他平台使用 bash
  - 作为永久启用的推荐方式

---

### func GetInstallCompletionExampleEN() map[string]string

```go
func GetInstallCompletionExampleEN() map[string]string
```

GetInstallCompletionExampleEN 获取当前平台的安装补全示例信息（英文）

**返回值:**
  - map[string]string: 包含安装补全示例信息的映射

**功能说明:**
  - 根据当前运行的操作系统返回对应的英文安装补全示例
  - Windows 使用 pwsh，其他平台使用 bash
  - 作为永久启用的推荐方式

---

### func ParseTimeWithCommonFormats(value string) (time.Time, string, error)

```go
func ParseTimeWithCommonFormats(value string) (time.Time, string, error)
```

ParseTimeWithCommonFormats 尝试使用常见格式解析时间字符串

**函数功能: **
  - 尝试使用常见格式解析时间字符串

**参数:**
  - value: 要解析的时间字符串

**返回值:**
  - time.Time: 解析后的时间
  - string: 使用的时间格式
  - error: 如果解析失败返回错误

**说明: **
  - 使用常见时间格式列表进行解析
  - 返回第一个成功解析的时间和格式

---

### func ParseTimeWithFormats(value string, formats []string) (time.Time, string, error)

```go
func ParseTimeWithFormats(value string, formats []string) (time.Time, string, error)
```

ParseTimeWithFormats 尝试使用多种格式解析时间字符串

**参数:**
  - value: 要解析的时间字符串
  - formats: 要尝试的时间格式列表, 按优先级排序

**返回值:**
  - time.Time: 解析后的时间
  - string: 使用的时间格式
  - error: 如果解析失败返回错误

**功能说明: **
  - 按给定格式列表顺序尝试解析
  - 返回第一个成功解析的时间和格式
  - 如果所有格式都失败, 返回错误

---

## TYPES

### type UnknownFlagError struct

```go
type UnknownFlagError struct {
    Command     string   // 当前命令名
    Input       string   // 用户输入的错误标志
    Suggestions []string // 相似标志建议列表
}
```

UnknownFlagError 未知标志错误

当用户输入的标志不存在时返回此错误, 包含相似标志建议。

**字段说明:**
  - Command: 当前命令名
  - Input: 用户输入的错误标志
  - Suggestions: 相似标志建议列表

**错误格式示例:**
```
myapp: unknown flag: '--verboose'

The most similar flags are
	--verbose
	-v
```

#### func (e *UnknownFlagError) Error() string

```go
func (e *UnknownFlagError) Error() string
```

Error 实现 error 接口, 返回格式化的错误信息

**返回值:**
  - string: 格式化的错误信息, 包含建议列表

---

### type UnknownSubcommandError struct

```go
type UnknownSubcommandError struct {
    Command     string   // 当前命令名
    Input       string   // 用户输入的错误子命令
    Suggestions []string // 相似子命令建议列表
}
```

UnknownSubcommandError 未知子命令错误

当用户输入的子命令不存在时返回此错误, 包含相似子命令建议。

**字段说明:**
  - Command: 当前命令名
  - Input: 用户输入的错误子命令
  - Suggestions: 相似子命令建议列表

**错误格式示例:**
```
myapp: 'cnfig' is not a valid command. See 'myapp --help'.

The most similar commands are
	config
```

#### func (e *UnknownSubcommandError) Error() string

```go
func (e *UnknownSubcommandError) Error() string
```

Error 实现 error 接口, 返回格式化的错误信息

**返回值:**
  - string: 格式化的错误信息, 包含建议列表

---

### type BuiltinFlagHandler interface

```go
type BuiltinFlagHandler interface {
    // Handle 处理内置标志
    //
    // 参数:
    //   - cmd: 要处理的命令
    //
    // 返回值:
    //   - error: 处理失败时返回错误
    //
    // 功能说明: 
    //   - 执行内置标志的特定操作
    //   - 通常在执行后会调用 os.Exit 退出程序
    //   - 例如: 帮助标志会打印帮助信息并退出
    Handle(cmd Command) error

    // Type 返回标志类型
    //
    // 返回值:
    //   - BuiltinFlagType: 标志类型
    //
    // 功能说明: 
    //   - 用于标识处理器处理的标志类型
    //   - 在注册和管理时使用
    Type() BuiltinFlagType

    // ShouldRegister 判断是否应该注册此标志
    //
    // 参数:
    //   - cmd: 要检查的命令
    //
    // 返回值:
    //   - bool: 是否应该注册
    //
    // 功能说明: 
    //   - 根据命令的配置决定是否注册此标志
    //   - 例如: 版本标志只有在设置了版本信息时才注册
    //   - 帮助标志总是注册
    ShouldRegister(cmd Command) bool

    // ShouldSkipRegistration 判断是否应该跳过注册
    //
    // 参数:
    //   - cmd: 要检查的命令
    //
    // 返回值:
    //   - bool: 是否应该跳过注册
    //
    // 功能说明: 
    //   - 检查标志是否已经被注册（避免重复注册）
    //   - 支持重复解析场景
    //   - 由每个处理器自己实现检查逻辑
    //
    // 使用场景:
    //   - 重复调用 Parse() 方法时避免重复注册
    //   - 测试场景中多次解析同一命令
    //   - 支持幂等性操作
    ShouldSkipRegistration(cmd Command) bool
}
```

BuiltinFlagHandler 内置标志处理器接口

内置标志处理器负责处理特定类型的内置标志。 每种内置标志类型都应该有一个对应的处理器实现。

---

### type BuiltinFlagType int

```go
type BuiltinFlagType int
```

BuiltinFlagType 内置标志类型

内置标志是系统自动处理的特殊标志, 如帮助标志和版本标志。 这些标志在解析完成后会被自动检查, 如果被设置则执行相应的操作。

```go
const (
    // HelpFlag 帮助标志
    //
    // 帮助标志用于显示命令的帮助信息, 包括用法、选项、子命令等。
    // 总是会被注册, 因为所有命令都应该支持帮助功能。
    HelpFlag BuiltinFlagType = iota

    // VersionFlag 版本标志
    //
    // 版本标志用于显示命令的版本信息。
    // 只有在命令设置了版本信息时才会被注册。
    VersionFlag

    // CompletionFlag 补全标志
    //
    // 补全标志用于生成Shell自动补全脚本。
    // 总是会被注册, 支持bash和pwsh两种Shell类型。
    CompletionFlag
)
```

---

### type CmdConfig struct

```go
type CmdConfig struct {
    Version        string            // 版本号
    UseChinese     bool              // 是否使用中文
    EnvPrefix      string            // 环境变量前缀
    UsageSyntax    string            // 命令使用语法
    Example        map[string]string // 示例使用, key为描述, value为示例命令
    Notes          []string          // 注意事项
    LogoText       string            // 命令logo文本
    MutexGroups    []MutexGroup      // 互斥组列表
    RequiredGroups []RequiredGroup   // 必需组列表
    Completion     bool              // 是否启用自动补全标志
    DynamicCompletion bool           // 是否启用动态补全
}
```

CmdConfig 命令配置类型

#### func NewCmdConfig() *CmdConfig

```go
func NewCmdConfig() *CmdConfig
```

NewCmdConfig 创建新的命令配置

**返回值:**
  - *CmdConfig: 新创建的 CmdConfig 实例, 初始化为零值

#### func (c *CmdConfig) Clone() *CmdConfig

```go
func (c *CmdConfig) Clone() *CmdConfig
```

Clone 克隆命令配置

**返回值:**
  - *CmdConfig: 克隆后的新 CmdConfig 实例

**功能说明:**
  - 创建当前配置的深拷贝
  - 复制所有字段值
  - 复制切片和映射时创建新的底层数组/映射
  - 用于避免配置共享导致的副作用

---

### type CmdRegistry interface

```go
type CmdRegistry interface {
    // Register 注册新命令到注册表
    //
    // 参数:
    //   - cmd: 要注册的命令
    //
    // 返回值:
    //   - error: 注册失败时返回错误
    //
    // 错误情况: 
    //   - 命令为nil或名称为空
    //   - 命令名称已存在
    Register(cmd Command) error

    // Unregister 从注册表中移除指定命令
    //
    // 参数:
    //   - name: 要移除的命令名称
    //
    // 返回值:
    //   - error: 移除失败时返回错误
    //
    // 错误情况: 
    //   - 命令不存在
    Unregister(name string) error

    // Get 根据名称获取命令
    //
    // 参数:
    //   - name: 命令名称
    //
    // 返回值:
    //   - Command: 找到的命令
    //   - bool: 是否找到, true表示找到
    //
    // 功能说明: 
    //   - 支持长名称查找
    //   - 未找到时返回零值和false
    Get(name string) (Command, bool)

    // List 获取所有注册的命令列表
    //
    // 返回值:
    //   - []Command: 所有命令的切片
    //
    // 功能说明: 
    //   - 返回注册表中所有命令
    //   - 顺序不确定, 取决于实现
    List() []Command

    // Has 检查指定名称的命令是否存在
    //
    // 参数:
    //   - name: 要检查的命令名称
    //
    // 返回值:
    //   - bool: 是否存在, true表示存在
    //
    // 功能说明: 
    //   - 快速存在性检查
    //   - 不返回命令本身, 提高效率
    Has(name string) bool

    // Count 获取注册表中的命令数量
    //
    // 返回值:
    //   - int: 命令总数
    //
    // 功能说明: 
    //   - 返回当前注册的命令数量
    //   - 时间复杂度应为O(1)
    Count() int

    // Clear 清空注册表中的所有命令
    //
    // 功能说明: 
    //   - 移除所有命令
    //   - 重置注册表到初始状态
    //   - 释放相关内存
    Clear()
}
```

CmdRegistry 命令注册表接口

CmdRegistry 定义了命令注册和管理的标准接口, 提供了 命令的完整生命周期管理功能。

**核心功能: **
  - 命令的注册和注销
  - 基于名称的查找和检索
  - 批量操作和遍历支持
  - 存在性检查和计数

**设计特点: **
  - 支持长名称和短名称查找
  - 提供统一的错误处理
  - 支持别名管理 (通过具体实现) 
  - 线程安全由具体实现保证

---

### type Command interface

```go
type Command interface {
    // 基本属性
    Name() string      // 命令名称, 用于匹配和显示
    LongName() string  // 长名称, 用于显示和帮助
    ShortName() string // 短名称, 用于命令行输入
    Desc() string      // 命令描述, 用于帮助显示

    // 标志管理
    AddFlag(flag Flag) error          // 添加一个标志到命令
    AddFlags(flags ...Flag) error     // 添加多个标志到命令
    AddFlagsFrom(flags []Flag) error  // 从切片添加多个标志
    GetFlag(name string) (Flag, bool) // 根据名称获取标志
    Flags() []Flag                    // 获取所有标志
    FlagRegistry() FlagRegistry       // 获取标志注册器

    // 子命令管理
    AddSubCmds(cmds ...Command) error      // 添加多个子命令
    AddSubCmdFrom(cmds []Command) error    // 从切片添加子命令
    GetSubCmd(name string) (Command, bool) // 根据名称获取子命令
    SubCmds() []Command                    // 获取所有子命令
    HasSubCmd(name string) bool            // 是否有指定名称的子命令
    CmdRegistry() CmdRegistry              // 获取子命令注册器

    // 命令层次
    IsRootCmd() bool // 是否为根命令
    Path() string    // 命令的路径, 用于显示和帮助

    // 参数解析
    Parse(args []string) error         // 解析命令行参数
    ParseAndRoute(args []string) error // 解析并路由到子命令
    ParseOnly(args []string) error     // 仅解析参数, 不路由
    IsParsed() bool                    // 是否已解析参数
    SetParsed(parsed bool)             // 设置解析状态

    // 参数访问
    Args() []string        // 获取所有参数
    Arg(index int) string  // 获取指定索引的参数
    NArg() int             // 获取参数数量
    SetArgs(args []string) // 设置参数

    // 执行
    Run() error                    // 执行命令
    SetRun(fn func(Command) error) // 设置执行函数
    HasRunFunc() bool              // 是否有执行函数

    // 帮助信息
    Help() string // 获取命令帮助信息
    PrintHelp()   // 打印命令帮助信息

    // 配置
    SetParser(p Parser)                     // 设置解析器
    SetDesc(desc string)                    // 设置命令描述
    SetVersion(version string)              // 设置命令版本
    SetChinese(useChinese bool)             // 设置是否使用中文
    SetEnvPrefix(prefix string)             // 设置环境变量前缀
    SetUsageSyntax(syntax string)           // 设置命令行语法
    AddExample(title, cmd string)           // 添加一个示例
    AddExamples(examples map[string]string) // 添加多个示例
    AddNote(note string)                    // 添加一条注意事项
    AddNotes(notes []string)                // 添加多条注意事项
    SetLogoText(logo string)                // 设置命令logo文本
    Config() *CmdConfig                     // 获取命令配置

    // 禁用标志解析
    IsDisableFlagParsing() bool             // 检查是否禁用标志解析
    SetDisableFlagParsing(disable bool)     // 设置是否禁用标志解析

    // 隐藏命令
    IsHidden() bool                         // 检查命令是否隐藏
    SetHidden(hidden bool)                  // 设置命令是否隐藏

    // 环境变量绑定
    AutoBindAllEnv() // 为所有标志自动绑定环境变量
}
```

Command 接口定义了命令的核心行为

#### 禁用标志解析相关方法

##### func IsDisableFlagParsing() bool

```go
func IsDisableFlagParsing() bool
```

IsDisableFlagParsing 检查是否禁用标志解析

**返回值:**
  - bool: 如果禁用标志解析返回 true，否则返回 false

**功能说明:**
  - 获取命令的禁用标志解析状态
  - 当返回 true 时，解析器会跳过标志解析阶段
  - 所有参数（包括 `--flag` 形式）都作为位置参数处理
  - 不影响子命令的路由功能

**使用场景:**
  - 包装外部命令（如 kubectl exec、docker run）
  - 需要透传参数给子进程的场景
  - Shell 脚本包装器

**注意事项:**
  - 默认值为 false（不禁用）
  - 只影响当前命令的标志解析，不影响子命令
  - 禁用后，`--help` 和 `--version` 等内置标志也不会被特殊处理

---

##### func SetDisableFlagParsing(disable bool)

```go
func SetDisableFlagParsing(disable bool)
```

SetDisableFlagParsing 设置是否禁用标志解析

**参数:**
  - disable: 是否禁用标志解析，true 表示禁用，false 表示不禁用

**功能说明:**
  - 设置命令的禁用标志解析状态
  - 设置为 true 后，解析器会跳过标志解析阶段
  - 所有参数原样保留为位置参数
  - 子命令路由功能正常工作

**使用示例:**
```go
cmd := NewCmd("exec", "e", ExitOnError)
cmd.SetDisableFlagParsing(true)  // 禁用标志解析
cmd.SetRun(func(c Command) error {
    args := c.Args()  // 所有参数都作为位置参数
    // 透传给外部命令
    return nil
})
```

**注意事项:**
  - 应在解析前设置，通常在命令创建后立即设置
  - 每个命令可以独立设置，父命令的设置不影响子命令
  - 禁用后，环境变量绑定也会被跳过

---

#### 隐藏命令相关方法

##### func IsHidden() bool

```go
func IsHidden() bool
```

IsHidden 检查命令是否隐藏

**返回值:**
  - bool: 如果命令是隐藏的返回 true，否则返回 false

**功能说明:**
  - 获取命令的隐藏状态
  - 隐藏命令不会显示在帮助信息的子命令列表中
  - 但仍可以通过命令行正常调用
  - 默认值为 false（不隐藏）

**使用场景:**
  - 创建内部命令或调试命令
  - 隐藏已弃用但仍需兼容的命令
  - 隐藏高级或实验性功能

---

##### func SetHidden(hidden bool)

```go
func SetHidden(hidden bool)
```

SetHidden 设置命令是否隐藏

**参数:**
  - hidden: 是否隐藏命令，true 表示隐藏，false 表示不隐藏

**功能说明:**
  - 设置命令的隐藏状态
  - 隐藏后命令不会出现在帮助信息中
  - 不影响命令的正常执行和路由
  - 子命令可以独立设置隐藏状态

**使用示例:**
```go
cmd := NewCmd("debug", "d", ExitOnError)
cmd.SetHidden(true)  // 隐藏调试命令
cmd.SetRun(func(c Command) error {
    // 执行调试逻辑
    return nil
})
```

**注意事项:**
  - 隐藏命令仍可通过命令行正常调用
  - 只是不在帮助信息的子命令列表中显示
  - 适用于内部命令或高级功能

---

### type Error struct

```go
type Error struct {
    Code    string // 错误码, 用于错误分类
    Message string // 错误消息, 面向用户
    Cause   error  // 原始错误, 底层错误原因
}
```

Error 错误类型

Error 是qflag项目的标准错误类型, 提供了结构化的错误信息。 包含错误码、错误消息和原始错误, 便于错误分类和处理。

**字段说明: **
  - Code: 错误码, 用于错误分类和程序化处理
  - Message: 错误消息, 面向用户的描述信息
  - Cause: 原始错误, 包装的底层错误

**特性: **
  - 实现error接口
  - 支持错误链 (errors.Unwrap) 
  - 支持错误比较 (errors.Is) 
  - 提供错误码匹配

#### func NewError(code, message string, cause error) *Error

```go
func NewError(code, message string, cause error) *Error
```

NewError 创建新的错误

**参数:**
  - code: 错误码, 用于错误分类和识别
  - message: 错误消息, 面向用户的描述信息
  - cause: 原始错误, 可以为nil

**返回值:**
  - *Error: 新创建的错误实例

**功能说明: **

#### func (e *Error) Error() string

```go
func (e *Error) Error() string
```

Error 实现 error 接口

**返回值:**
  - string: 格式化的错误字符串

**功能说明: **
  - 返回用户友好的错误信息
  - 包含原始错误信息 (如果有) 
  - 格式: 消息 + ": " + 原始错误

#### func (e *Error) Is(target error) bool

```go
func (e *Error) Is(target error) bool
```

Is 判断错误是否相同

**参数:**
  - target: 要比较的目标错误

**返回值:**
  - bool: 是否相同, true表示相同

**功能说明: **
  - 基于错误码进行比较
  - 支持errors.Is函数
  - 忽略错误消息和原始错误

#### func (e *Error) Unwrap() error

```go
func (e *Error) Unwrap() error
```

Unwrap 实现 errors.Unwrap 接口

**返回值:**
  - error: 原始错误

**功能说明: **
  - 支持错误链操作
  - 允许使用errors.Unwrap获取底层错误
  - 支持errors.As和errors.Is

---

### type ErrorHandling = flag.ErrorHandling

```go
type ErrorHandling = flag.ErrorHandling
```

ErrorHandling 错误处理方式枚举

ErrorHandling 定义了解析错误时的处理策略, 直接使用标准库 flag包的错误处理方式, 保持兼容性。

**可选值: **
  - ContinueOnError: 解析错误时继续解析并返回错误
  - ExitOnError: 解析错误时退出程序
  - PanicOnError: 解析错误时触发panic

```go
var (
    // ContinueOnError 解析错误时继续解析并返回错误
    //
    // 使用场景: 
    //   - 需要收集所有错误
    //   - 自定义错误处理逻辑
    //   - 交互式应用
    ContinueOnError ErrorHandling = flag.ContinueOnError

    // ExitOnError 解析错误时退出程序
    //
    // 使用场景: 
    //   - 简单命令行工具
    //   - 错误即致命的应用
    //   - 脚本和自动化工具
    ExitOnError ErrorHandling = flag.ExitOnError

    // PanicOnError 解析错误时触发panic
    //
    // 使用场景: 
    //   - 开发和测试环境
    //   - 需要快速失败的场景
    //   - 调试和诊断
    PanicOnError ErrorHandling = flag.PanicOnError
)
```

---

### type Flag interface

```go
type Flag interface {
    // Name 获取标志名称
    //
    // 返回值:
    //   - string: 标志的完整名称
    //
    // 功能说明: 
    //   - 返回标志的主要标识符
    //   - 用于命令行参数和查找
    //   - 名称在注册表中必须唯一
    Name() string

    // LongName 获取标志长名称
    //
    // 返回值:
    //   - string: 标志的长名称
    //
    // 功能说明: 
    //   - 与Name方法功能相同
    //   - 提供语义明确的方法名
    //   - 保持接口一致性
    LongName() string

    // ShortName 获取标志短名称
    //
    // 返回值:
    //   - string: 标志的短名称
    //
    // 功能说明: 
    //   - 返回标志的简短形式
    //   - 通常为单个字符
    //   - 可能为空字符串
    ShortName() string

    // Desc 获取标志描述
    //
    // 返回值:
    //   - string: 标志的描述信息
    //
    // 功能说明: 
    //   - 返回标志的功能说明
    //   - 用于帮助信息生成
    //   - 应简洁明了地描述标志用途
    Desc() string

    // Type 获取标志类型
    //
    // 返回值:
    //   - FlagType: 标志的类型枚举
    //
    // 功能说明: 
    //   - 返回标志的数据类型
    //   - 用于类型检查和转换
    //   - 决定解析和验证逻辑
    Type() FlagType

    // Set 设置标志值
    //
    // 参数:
    //   - value: 要设置的字符串值
    //
    // 返回值:
    //   - error: 设置失败时返回错误
    //
    // 功能说明: 
    //   - 从字符串解析并设置值
    //   - 自动进行类型转换
    //   - 更新内部状态和标记
    Set(value string) error

    // GetDef 获取默认值
    //
    // 返回值:
    //   - any: 标志的默认值
    //
    // 功能说明: 
    //   - 返回初始化时设置的默认值
    //   - 用于帮助信息显示
    //   - 用户未设置值时使用此值
    GetDef() any

    // GetStr 获取标志当前值的字符串表示
    //
    // 返回值:
    //   - string: 标志当前值的字符串表示
    //
    // 功能说明: 
    //   - 获取标志当前值的字符串表示
    //   - 与String()方法不同, 此方法专注于值本身
    //   - 用于内置标志处理中获取标志值
    GetStr() string

    // IsSet 检查标志是否被用户设置
    //
    // 返回值:
    //   - bool: 是否被设置, true表示已设置
    //
    // 功能说明: 
    //   - 区分默认值和用户设置值
    //   - 用于必填标志检查
    //   - 影响某些标志的行为逻辑
    IsSet() bool

    // Reset 重置标志到默认状态
    //
    // 功能说明: 
    //   - 清除用户设置的值
    //   - 恢复到默认值
    //   - 重置设置状态标记
    Reset()

    // String 获取标志值的字符串表示
    //
    // 返回值:
    //   - string: 值的字符串表示
    //
    // 功能说明: 
    //   - 实现fmt.Stringer接口
    //   - 用于显示和日志输出
    //   - 格式应与输入格式兼容
    String() string

    // BindEnv 绑定环境变量
    //
    // 参数:
    //   - name: 环境变量名称
    //
    // 功能说明: 
    //   - 从环境变量读取默认值
    //   - 优先级低于命令行参数
    //   - 支持配置文件和环境变量
    BindEnv(name string)

    // AutoBindEnv 自动绑定环境变量
    //
    // 功能说明:
    //   - 自动使用标志的长名称作为环境变量名（转为大写）
    //   - 如果没有设置长名称，会触发 panic
    //   - 环境变量前缀（EnvPrefix）在解析时自动拼接，无需手动处理
    //
    // 注意事项:
    //   - 环境变量的优先级低于命令行参数
    //   - 必须设置长名称，否则会 panic
    //   - 短名称不会被使用，避免冲突
    //   - 自动转为大写，确保环境变量命名规范
    AutoBindEnv()

    // GetEnvVar 获取绑定的环境变量名
    //
    // 返回值:
    //   - string: 环境变量名称
    //
    // 功能说明: 
    //   - 返回当前绑定的环境变量
    //   - 未绑定时返回空字符串
    //   - 用于调试和配置管理
    GetEnvVar() string

    // EnumValues 获取枚举类型的可选值
    //
    // 返回值:
    //   - []string: 枚举类型的可选值列表
    //
    // 功能说明: 
    //   - 非枚举类型返回空切片
    //   - 枚举类型返回所有可选值
    //   - 用于补全脚本生成和验证
    EnumValues() []string
}
```

Flag 接口定义了标志的核心行为

Flag 是所有标志类型必须实现的基础接口, 定义了标志的 基本操作和属性。所有具体标志类型都应实现此接口。

**设计原则: **
  - 提供统一的标志操作接口
  - 支持多种数据类型
  - 支持验证和环境变量绑定
  - 提供完整的生命周期管理

---

### type FlagRegistry interface

```go
type FlagRegistry interface {
    // Register 注册新标志到注册表
    //
    // 参数:
    //   - flag: 要注册的标志
    //
    // 返回值:
    //   - error: 注册失败时返回错误
    //
    // 错误情况: 
    //   - 标志为nil或名称为空
    //   - 标志名称已存在
    Register(flag Flag) error

    // Unregister 从注册表中移除指定标志
    //
    // 参数:
    //   - name: 要移除的标志名称
    //
    // 返回值:
    //   - error: 移除失败时返回错误
    //
    // 错误情况: 
    //   - 标志不存在
    Unregister(name string) error

    // Get 根据名称获取标志
    //
    // 参数:
    //   - name: 标志名称
    //
    // 返回值:
    //   - Flag: 找到的标志
    //   - bool: 是否找到, true表示找到
    //
    // 功能说明: 
    //   - 支持长名称查找
    //   - 未找到时返回零值和false
    Get(name string) (Flag, bool)

    // List 获取所有注册的标志列表
    //
    // 返回值:
    //   - []Flag: 所有标志的切片
    //
    // 功能说明: 
    //   - 返回注册表中所有标志
    //   - 顺序不确定, 取决于实现
    List() []Flag

    // Has 检查指定名称的标志是否存在
    //
    // 参数:
    //   - name: 要检查的标志名称
    //
    // 返回值:
    //   - bool: 是否存在, true表示存在
    //
    // 功能说明: 
    //   - 快速存在性检查
    //   - 不返回标志本身, 提高效率
    Has(name string) bool

    // Count 获取注册表中的标志数量
    //
    // 返回值:
    //   - int: 标志总数
    //
    // 功能说明: 
    //   - 返回当前注册的标志数量
    //   - 时间复杂度应为O(1)
    Count() int

    // Clear 清空注册表中的所有标志
    //
    // 功能说明: 
    //   - 移除所有标志
    //   - 重置注册表到初始状态
    //   - 释放相关内存
    Clear()
}
```

FlagRegistry 标志注册表接口

FlagRegistry 定义了标志注册和管理的标准接口, 提供了 标志的完整生命周期管理功能。

**核心功能: **
  - 标志的注册和注销
  - 基于名称的查找和检索
  - 批量操作和遍历支持
  - 存在性检查和计数

**设计特点: **
  - 支持长名称和短名称查找
  - 提供统一的错误处理
  - 支持别名管理 (通过具体实现) 
  - 线程安全由具体实现保证

---

### type FlagType int

```go
type FlagType int
```

FlagType 标志类型枚举

FlagType 定义了所有支持的标志类型, 用于类型识别和 特定处理逻辑的实现。

**设计原则: **
  - 每种类型对应一种数据格式
  - 支持基础类型和复合类型
  - 便于类型检查和转换

```go
const (
    FlagTypeUnknown FlagType = iota // 未知标志类型, 用于错误处理

    // 基础类型
    FlagTypeString  // 字符串标志, 存储任意文本
    FlagTypeInt     // 整数标志, 平台相关int类型
    FlagTypeInt64   // 64位整数标志, 固定64位整数
    FlagTypeUint    // 无符号整数标志, 平台相关uint类型
    FlagTypeUint8   // 8位无符号整数标志, 0-255
    FlagTypeUint16  // 16位无符号整数标志, 0-65535
    FlagTypeUint32  // 32位无符号整数标志, 0-4294967295
    FlagTypeUint64  // 64位无符号整数标志, 0-18446744073709551615
    FlagTypeFloat64 // 64位浮点数标志, IEEE 754双精度
    FlagTypeBool    // 布尔标志, true/false值

    // 特殊类型
    FlagTypeEnum // 枚举标志, 限制为预定义值集合

    // 时间和大小类型
    FlagTypeDuration // 持续时间标志, 支持时间单位解析
    FlagTypeTime     // 时间标志, 支持多种时间格式
    FlagTypeSize     // 大小标志, 支持存储单位解析

    // 集合类型
    FlagTypeMap         // 映射标志, 键值对集合
    FlagTypeStringSlice // 字符串切片标志, 字符串数组
    FlagTypeIntSlice    // 整数切片标志, 整数数组
    FlagTypeInt64Slice  // 64位整数切片标志, 64位整数数组
)
```

#### func (t FlagType) IsNumericType() bool

```go
func (t FlagType) IsNumericType() bool
```

IsNumericType 检查是否为数值类型

**返回值:**
  - bool: 是否为数值类型, true表示是

**功能说明: **
  - 识别所有数值类型的标志
  - 包括整数、浮点数和大小类型
  - 用于数值范围验证

#### func (t FlagType) IsSliceType() bool

```go
func (t FlagType) IsSliceType() bool
```

IsSliceType 检查是否为切片类型

**返回值:**
  - bool: 是否为切片类型, true表示是

**功能说明: **
  - 识别所有切片类型的标志
  - 用于特殊处理逻辑
  - 支持多值输入的标志

#### func (t FlagType) IsValid() bool

```go
func (t FlagType) IsValid() bool
```

IsValid 检查标志类型是否有效

**返回值:**
  - bool: 是否有效, true表示有效

**功能说明: **
  - 排除未知类型
  - 用于类型验证
  - 确保类型在预定义范围内

#### func (t FlagType) String() string

```go
func (t FlagType) String() string
```

String 返回标志类型的字符串表示

**返回值:**
  - string: 类型的可读字符串表示

**功能说明: **
  - 提供人类可读的类型名称
  - 用于错误消息和日志
  - 未知类型返回格式化字符串

**示例: **
  - FlagTypeString -> "string"
  - FlagTypeIntSlice -> "[]int"
  - FlagType(999) -> "FlagType(999)"

---

### type MutexGroup struct

```go
type MutexGroup struct {
    Name      string   // 互斥组名称, 用于错误提示和标识
    Flags     []string // 互斥组中的标志名称列表
    AllowNone bool     // 是否允许一个都不设置
}
```

MutexGroup 互斥组定义

MutexGroup 定义了一组互斥的标志, 其中最多只能有一个被设置。 当用户设置了互斥组中的多个标志时, 解析器会返回错误。

**字段说明: **
  - Name: 互斥组名称, 用于错误提示和标识
  - Flags: 互斥组中的标志名称列表
  - AllowNone: 是否允许一个都不设置

**使用场景: **
  - 输出格式互斥 (如 --json 和 --xml 不能同时使用) 
  - 操作模式互斥 (如 --create 和 --update 不能同时使用) 
  - 必选选项 (如必须指定 --file 或 --url 中的一个) 

---

### type RequiredGroup struct

```go
type RequiredGroup struct {
    Name       string   // 必需组名称，用于错误提示和标识
    Flags      []string // 必需组中的标志名称列表
    Conditional bool     // 是否为条件性必需组
}
```

RequiredGroup 必需组定义

RequiredGroup 定义了一组标志的必需关系。支持两种模式:
1. 普通必需组: 组中的所有标志都必须被设置
2. 条件性必需组: 如果组中任何一个标志被设置, 则所有标志都必须被设置

**字段说明:**
  - Name: 必需组名称, 用于错误提示和标识
  - Flags: 必需组中的标志名称列表
  - Conditional: 是否为条件性必需组

**使用场景:**
  - 普通必需组: 连接参数必需 (如 --host 和 --port 必须同时设置)
  - 条件性必需组: 可选但相关的标志组合 (如 --host 和 --port, 如果使用其中一个则必须同时使用)
  - 认证参数: 普通必需组 (如 --username 和 --password 必须同时设置)
  - 配置文件路径: 条件性必需组 (如 --config 和 --env, 如果使用其中一个则必须同时使用)

**注意事项:**
  - 对于普通必需组, 组中的所有标志都必须被设置
  - 对于条件性必需组, 如果组中任何一个标志被设置, 则所有标志都必须被设置
  - 如果只设置了部分标志, 解析会失败
  - 必需组名称在命令中应该唯一

---

### type Validator[T any]

```go
type Validator[T any] func(value T) error
```

Validator 验证器函数类型

Validator 是一个泛型函数类型，用于验证标志值的有效性。验证器接收一个类型为 T 的值，返回错误信息。

**参数:**
  - value: 要验证的值

**返回值:**
  - error: 验证失败时返回错误，验证通过返回 nil

**功能说明:**
  - 验证器在标志的 Set 方法中被调用
  - 在解析完值后、设置值之前执行验证
  - 如果验证失败，Set 方法会返回错误，标志值不会被设置
  - 验证器是可选的，未设置时跳过验证
  - 重复设置验证器会覆盖之前的验证器

**使用示例:**

```go
// 端口号验证：1-65535
port.SetValidator(func(value int) error {
    if value < 1 || value > 65535 {
        return fmt.Errorf("端口 %d 超出范围 [1, 65535]", value)
    }
    return nil
})

// 字符串长度验证：3-20个字符
username.SetValidator(func(value string) error {
    if len(value) < 3 || len(value) > 20 {
        return fmt.Errorf("用户名长度 %d 超出范围 [3, 20]", len(value))
    }
    return nil
})

// 邮箱格式验证
email.SetValidator(func(value string) error {
    if !isValidEmail(value) {
        return fmt.Errorf("邮箱格式无效: %s", value)
    }
    return nil
})

// 自定义验证：检查端口是否被占用
port.SetValidator(func(value int) error {
    if isPortInUse(value) {
        return fmt.Errorf("端口 %d 已被占用", value)
    }
    return nil
})
```

**注意事项:**
  - 验证器应该快速执行，避免耗时操作
  - 验证器返回的错误应该清晰描述失败原因
  - 验证器执行时已经持有锁，验证器本身不需要处理并发
  - 空字符串（对于 BoolFlag 和集合类型）不经过验证
  - 验证器可以随时通过 ClearValidator 清除

---

### type OptionInfo struct

```go
type OptionInfo struct {
    NamePart string
    Desc     string
    DefValue string
}
```

用于存储选项的信息

---

### type Parser interface

```go
type Parser interface {
    // ParseOnly 解析当前命令的参数, 不递归解析子命令
    //
    // 参数:
    //   - cmd: 要解析的命令
    //   - args: 命令行参数列表
    //
    // 返回值:
    //   - error: 解析失败时返回错误
    //
    // 功能说明: 
    //   - 仅解析当前命令的标志和参数
    //   - 不处理子命令的解析
    //   - 适用于需要手动控制子命令处理的场景
    //
    // 使用场景: 
    //   - 多阶段解析
    //   - 自定义子命令处理逻辑
    //   - 参数预处理和验证
    ParseOnly(cmd Command, args []string) error

    // Parse 单纯解析, 递归解析子命令
    //
    // 参数:
    //   - cmd: 要解析的根命令
    //   - args: 命令行参数列表
    //
    // 返回值:
    //   - error: 解析失败时返回错误
    //
    // 功能说明: 
    //   - 解析根命令的标志和参数
    //   - 递归解析所有子命令
    //   - 构建完整的命令树结构
    //   - 不执行任何命令的运行函数
    //
    // 使用场景: 
    //   - 命令结构验证
    //   - 帮助信息生成
    //   - 配置预检查
    Parse(cmd Command, args []string) error

    // ParseAndRoute 解析并且路由执行
    //
    // 参数:
    //   - cmd: 要解析和执行的根命令
    //   - args: 命令行参数列表
    //
    // 返回值:
    //   - error: 解析或执行失败时返回错误
    //
    // 功能说明: 
    //   - 解析根命令和所有子命令
    //   - 根据参数路由到相应的命令
    //   - 执行最终目标命令的运行函数
    //   - 提供完整的命令行处理流程
    //
    // 执行流程: 
    //   1. 解析根命令的标志和参数
    //   2. 识别并解析子命令
    //   3. 递归处理直到找到最终命令
    //   4. 执行最终命令的Run函数
    //
    // 使用场景: 
    //   - 标准命令行应用
    //   - CLI工具的主入口
    //   - 自动化脚本执行
    ParseAndRoute(cmd Command, args []string) error
}
```

Parser 解析器接口

Parser 定义了命令行参数解析的标准接口, 提供了不同层次的 解析功能, 从简单的参数解析到完整的命令路由执行。

**设计理念: **
  - 分层设计: 提供不同层次的解析功能
  - 灵活性: 支持仅解析、解析+路由等多种使用模式
  - 可扩展性: 接口设计允许不同的解析策略实现

**使用场景: **
  - 命令行工具的参数解析
  - 子命令系统的路由管理
  - 配置管理和参数验证

---

### type SubCmdInfo struct

```go
type SubCmdInfo struct {
    Name string
    Desc string
}
```

用于存储子命令的信息