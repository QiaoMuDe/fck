# Package cmd 

```go
import "gitee.com/MM-Q/qflag/internal/cmd"
```

---

## 包介绍

Package cmd 提供命令实现和命令管理功能

cmd 包实现了 types.Command 接口, 提供了完整的命令行命令功能。 主要组件: 
  - Cmd: 命令结构体, 实现了所有命令相关接口
  - 命令生命周期管理
  - 标志和子命令管理
  - 解析和执行功能

特性: 
  - 线程安全的命令结构
  - 支持嵌套子命令
  - 灵活的配置选项
  - 完整的帮助系统
  - 内置动态补全支持

---

## 内置子命令

### __complete 子命令

当启用动态补全功能时（必须先设置 `Completion: true`，再设置 `DynamicCompletion: true`），cmd 包会自动注册 `__complete` 内部子命令。

**注意**: 如果未启用 `Completion` 就直接启用 `DynamicCompletion`，会触发 panic。

**用途**: 为 Shell 自动补全脚本提供实时命令树查询服务

**支持指令**:
- `fuzzy` - 执行模糊匹配补全
- `context` - 计算上下文路径
- `candidates` - 获取候选选项
- `enum` - 获取枚举值
- `all` - 统一获取所有补全信息

**示例**:
```bash
# 计算上下文路径
myapp __complete context server start

# 获取候选选项
myapp __complete candidates /server/start/

# 统一获取补全信息
myapp __complete all "" "" server start
```

**调试模式**:
设置环境变量 `QFLAG_COMPLETE_DEBUG=1` 可启用调试输出

---

## TYPES

### type Cmd struct

```go
type Cmd struct {
    // Has unexported fields.
}
```

Cmd 是一个命令结构体, 实现了 types.Command 接口

Cmd 提供了完整的命令行命令实现, 支持标志管理、子命令、 参数解析和执行等功能。使用读写锁保证并发安全。

线程安全: 
  - 所有公共方法都使用读写锁保护
  - 支持并发读取和独占写入
  - 解析操作使用sync.Once确保只执行一次

#### func NewCmd(longName, shortName string, errorHandling types.ErrorHandling) *Cmd

```go
func NewCmd(longName, shortName string, errorHandling types.ErrorHandling) *Cmd
```

NewCmd 创建新的命令实例

**参数:**
  - longName: 命令的长名称
  - shortName: 命令的短名称
  - errorHandling: 错误处理策略

**返回值:**
  - *Cmd: 初始化完成的命令实例

**功能说明:**
  - 创建命令并初始化基本字段
  - 创建标志和子命令注册器
  - 设置默认解析器
  - 初始化配置选项

### type CmdOpts struct

```go
type CmdOpts struct {
    // 基本属性
    Desc string // 命令描述

    // 运行函数
    RunFunc func(types.Command) error // 命令执行函数

    // 配置选项
    Version     string // 版本号
    UseChinese  bool   // 是否使用中文
    EnvPrefix   string // 环境变量前缀
    UsageSyntax string // 命令使用语法
    LogoText    string // Logo文本
    Completion  bool   // 是否启用自动补全标志
    DynamicCompletion bool // 是否启用动态补全标志，默认 false。启用后会自动注册 __complete 子命令

    // 禁用标志解析
    DisableFlagParsing bool // 是否禁用标志解析，默认 false

    // 隐藏命令
    Hidden bool // 是否隐藏命令，默认 false。隐藏命令不会显示在帮助信息中

    // 环境变量绑定
    AutoBindEnv bool // 是否自动绑定所有标志的环境变量

    // 示例和说明
    Examples map[string]string // 示例使用, key为描述, value为示例命令
    Notes    []string          // 注意事项

    // 子命令和互斥组
    SubCmds          []types.Command        // 子命令列表, 用于添加到命令中
    MutexGroups      []types.MutexGroup     // 互斥组列表
    RequiredGroups   []types.RequiredGroup  // 必需组列表
    FlagDependencies []types.FlagDependency // 标志依赖关系列表
}
```

CmdOpts 是命令选项结构体, 提供了配置现有命令的方式

CmdOpts 包含了命令的所有可配置属性, 用于配置已存在的命令。

字段说明:
  - Desc: 命令描述
  - RunFunc: 命令执行函数
  - Version: 版本号
  - UseChinese: 是否使用中文
  - EnvPrefix: 环境变量前缀
  - UsageSyntax: 命令使用语法
  - LogoText: Logo文本
  - Completion: 是否启用自动补全标志，默认 false
  - DynamicCompletion: 是否启用动态补全标志，默认 false。启用后会自动注册 __complete 子命令
  - DisableFlagParsing: 是否禁用标志解析，默认 false。当设置为 true 时，所有参数（包括 `--flag` 形式）都作为位置参数处理
  - Hidden: 是否隐藏命令，默认 false。隐藏命令不会显示在帮助信息中，但仍可正常调用
  - AutoBindEnv: 是否自动绑定所有标志的环境变量
  - Examples: 示例使用, key为描述, value为示例命令
  - Notes: 注意事项
  - SubCmds: 子命令列表, 用于添加到命令中
  - MutexGroups: 互斥组列表
  - RequiredGroups: 必需组列表
  - FlagDependencies: 标志依赖关系列表

使用场景: 
  - 已有命令实例, 需要批量设置属性
  - 需要结构化的配置管理
  - 需要部分配置（未设置的属性不会被修改）

#### func NewCmdOpts() *CmdOpts

```go
func NewCmdOpts() *CmdOpts
```

NewCmdOpts 创建新的命令选项

**返回值:**
  - *CmdOpts: 初始化的命令选项

**功能说明:**
  - 创建基本命令选项
  - 初始化所有字段为零值
  - 初始化 map 和 slice 避免空指针

**默认值:**
  - Examples: 空映射
  - Notes: 空切片
  - SubCmds: 空切片
  - MutexGroups: 空切片
  - RequiredGroups: 空切片
  - FlagDependencies: 空切片

#### func (c *Cmd) ApplyOpts(opts *CmdOpts) error

```go
func (c *Cmd) ApplyOpts(opts *CmdOpts) error
```

ApplyOpts 应用选项到命令

**参数:**
  - opts: 命令选项

**返回值:**
  - error: 应用选项失败时返回错误

**功能说明:**
  - 将选项结构体的所有属性应用到当前命令
  - 支持部分配置（未设置的属性不会被修改）
  - 使用defer捕获panic, 转换为错误返回

**应用顺序:**
  1. 基本属性（Desc、RunFunc）
  2. 配置选项（Version、UseChinese、EnvPrefix、UsageSyntax、LogoText、Completion）
  3. 示例和说明（Examples、Notes）
  4. 互斥组（MutexGroups）
  5. 必需组（RequiredGroups）
  6. 标志依赖关系（FlagDependencies）
  7. 子命令（SubCmds）

**错误情况:**
  - 选项为 nil: 返回 INVALID_CMDOPTS 错误
  - 添加子命令失败: 返回 FAILED_TO_ADD_SUBCMDS 错误
  - panic: 转换为 PANIC 错误

**线程安全:**
  - 方法内部使用读写锁保护并发访问
  - 可以安全地在多个 goroutine 中调用

**设计说明:**
  - 调用现有的 SetDesc、SetVersion、AddExamples 等方法
  - 避免重复代码，降低维护成本
  - 保持行为一致性，与用户手动调用方法完全一致
  - 保留方法中的验证、通知等逻辑

#### func (c *Cmd) AddExample(title, cmd string)

```go
func (c *Cmd) AddExample(title, cmd string)
```

AddExample 添加单个示例

**参数:**
  - title: 示例标题
  - cmd: 示例命令

**功能说明:**
  - 添加命令使用示例
  - 用于帮助信息显示
  - 存储在配置中
  - 支持并发安全的添加

#### func (c *Cmd) AddExamples(examples map[string]string)

```go
func (c *Cmd) AddExamples(examples map[string]string)
```

AddExamples 批量添加示例

**参数:**
  - examples: 示例映射, 标题为键, 命令为值

**功能说明:**
  - 批量添加多个示例
  - 空映射直接返回
  - 覆盖同名的示例
  - 支持并发安全的添加

#### func (c *Cmd) AddFlag(f types.Flag) error

```go
func (c *Cmd) AddFlag(f types.Flag) error
```

AddFlag 添加标志到命令

**参数:**
  - f: 要添加的标志

**返回值:**
  - error: 添加失败时返回错误

**功能说明:**
  - 实现types.Command接口
  - 将标志注册到命令的标志注册器
  - 支持并发安全的添加操作

**错误情况:**
  - 标志为nil: 返回INVALID_FLAG错误
  - 标志名称冲突: 返回FLAG_ALREADY_EXISTS错误

#### func (c *Cmd) AddFlags(flags ...types.Flag) error

```go
func (c *Cmd) AddFlags(flags ...types.Flag) error
```

AddFlags 添加多个标志到命令

**参数:**
  - flags: 要添加的标志列表

**返回值:**
  - error: 添加失败时返回错误

**功能说明:**
  - 实现types.Command接口
  - 批量添加多个标志
  - 支持并发安全的添加操作

**错误情况:**
  - 标志为nil: 返回INVALID_FLAG错误
  - 标志名称冲突: 返回FLAG_ALREADY_EXISTS错误

#### func (c *Cmd) AddFlagsFrom(flags []types.Flag) error

```go
func (c *Cmd) AddFlagsFrom(flags []types.Flag) error
```

AddFlagsFrom 从切片添加多个标志

**参数:**
  - flags: 标志切片

**返回值:**
  - error: 添加失败时返回错误

**功能说明:**
  - 实现types.Command接口
  - 从切片中添加多个标志
  - 空切片直接返回成功
  - 内部调用AddFlags实现

#### func (c *Cmd) AddMutexGroup(name string, flags []string, allowNone bool) error

```go
func (c *Cmd) AddMutexGroup(name string, flags []string, allowNone bool) error
```

AddMutexGroup 添加互斥组到命令

**参数:**
  - name: 互斥组名称, 用于错误提示和标识
  - flags: 互斥组中的标志名称列表
  - allowNone: 是否允许一个都不设置

**返回值:**
  - error: 添加失败时返回错误

**功能说明:**
  - 创建新的互斥组并添加到命令配置中
  - 互斥组中的标志最多只能有一个被设置
  - 如果 allowNone 为 false, 则必须至少有一个标志被设置
  - 使用写锁保护并发安全

**错误情况:**
  - 组名已存在: 返回 MUTEX_GROUP_ALREADY_EXISTS 错误
  - 标志列表为空: 返回 EMPTY_MUTEX_GROUP 错误
  - 标志不存在: 返回 FLAG_NOT_FOUND 错误

**注意事项:**
  - 标志名称必须是已注册的标志
  - 互斥组名称在命令中应该唯一
  - 重复添加同名互斥组会返回错误

#### func (c *Cmd) AddNote(note string)

```go
func (c *Cmd) AddNote(note string)
```

AddNote 添加单个注释

**参数:**
  - note: 注释内容

**功能说明:**
  - 添加命令的额外说明
  - 用于帮助信息显示
  - 空字符串被忽略
  - 支持并发安全的添加

#### func (c *Cmd) AddNotes(notes []string)

```go
func (c *Cmd) AddNotes(notes []string)
```

AddNotes 批量添加注释

**参数:**
  - notes: 注释切片

**功能说明:**
  - 批量添加多个注释
  - 空切片直接返回
  - 空字符串被忽略
  - 支持并发安全的添加

#### func (c *Cmd) AddSubCmdFrom(cmds []types.Command) error

```go
func (c *Cmd) AddSubCmdFrom(cmds []types.Command) error
```

AddSubCmdFrom 从切片添加子命令

**参数:**
  - cmds: 子命令切片

**返回值:**
  - error: 添加失败时返回错误

**功能说明:**
  - 实现types.Command接口
  - 从切片中添加子命令
  - 空切片直接返回成功
  - 内部调用AddSubCmds实现

#### func (c *Cmd) AddSubCmds(cmds ...types.Command) error

```go
func (c *Cmd) AddSubCmds(cmds ...types.Command) error
```

AddSubCmds 添加子命令到命令

**参数:**
  - cmds: 要添加的子命令列表

**返回值:**
  - error: 添加失败时返回错误

**功能说明:**
  - 实现types.Command接口
  - 批量添加多个子命令
  - 自动设置父子关系
  - 支持并发安全的添加操作

**错误情况:**
  - 子命令为nil: 返回INVALID_COMMAND错误
  - 子命令类型错误: 返回INVALID_COMMAND_TYPE错误
  - 子命令名称冲突: 返回COMMAND_ALREADY_EXISTS错误

#### func (c *Cmd) Arg(index int) string

```go
func (c *Cmd) Arg(index int) string
```

Arg 获取指定索引的命令行参数

**参数:**
  - index: 命令行参数的索引

**返回值:**
  - string: 对应索引的命令行参数值

**注意:**
  - 索引从 0 开始计数
  - 如果索引超出范围, 返回空字符串

#### func (c *Cmd) Args() []string

```go
func (c *Cmd) Args() []string
```

Args 获取命令行参数

**返回值:**
  - []string: 命令行参数的副本

**功能说明:**
  - 实现types.Command接口
  - 返回解析后的参数列表
  - 创建副本避免外部修改
  - 支持并发安全的访问

#### func (c *Cmd) Bool(longName, shortName, description string, default_ bool) *flag.BoolFlag

```go
func (c *Cmd) Bool(longName, shortName, description string, default_ bool) *flag.BoolFlag
```

Bool 创建布尔标志

**参数:**
  - longName: 长标志名 (如 --long-name)
  - shortName: 短标志名 (如 -s)
  - description: 标志的描述信息
  - default_: 标志的默认值

**返回值:**
  - *flag.BoolFlag: 新创建的布尔标志

#### func (c *Cmd) CmdRegistry() types.CmdRegistry

```go
func (c *Cmd) CmdRegistry() types.CmdRegistry
```

CmdRegistry 获取子命令注册器

**返回值:**
  - types.CmdRegistry: 子命令注册器接口

**功能说明:**
  - 实现types.Command接口
  - 返回命令的子命令注册器
  - 用于直接操作子命令注册
  - 支持并发安全的访问

#### func (c *Cmd) Config() *types.CmdConfig

```go
func (c *Cmd) Config() *types.CmdConfig
```

Config 获取命令配置

**返回值:**
  - *types.CmdConfig: 命令配置的副本

**功能说明:**
  - 实现types.Command接口
  - 返回命令的配置对象
  - 注意: 返回的是副本, 修改不会影响命令
  - 支持并发安全的访问

#### func (c *Cmd) Desc() string

```go
func (c *Cmd) Desc() string
```

Desc 获取命令描述

**返回值:**
  - string: 命令的描述信息

**功能说明:**
  - 实现types.Command接口
  - 线程安全地访问描述信息
  - 用于帮助信息显示

#### func (c *Cmd) Duration(longName, shortName, description string, default_ time.Duration) *flag.DurationFlag

```go
func (c *Cmd) Duration(longName, shortName, description string, default_ time.Duration) *flag.DurationFlag
```

Duration 创建持续时间标志

**参数:**
  - longName: 长标志名 (如 --long-name)
  - shortName: 短标志名 (如 -s)
  - description: 标志的描述信息
  - default_: 标志的默认值

**返回值:**
  - *flag.DurationFlag: 新创建的持续时间标志

#### func (c *Cmd) Enum(longName, shortName, description, default_ string, allowedValues []string) *flag.EnumFlag

```go
func (c *Cmd) Enum(longName, shortName, description, default_ string, allowedValues []string) *flag.EnumFlag
```

Enum 创建枚举标志

**参数:**
  - longName: 长标志名 (如 --long-name)
  - shortName: 短标志名 (如 -s)
  - description: 标志的描述信息
  - default_: 标志的默认值
  - allowedValues: 允许的枚举值列表

**返回值:**
  - *flag.EnumFlag: 新创建的枚举标志

#### func (c *Cmd) FlagRegistry() types.FlagRegistry

```go
func (c *Cmd) FlagRegistry() types.FlagRegistry
```

FlagRegistry 获取标志注册器

**返回值:**
  - types.FlagRegistry: 标志注册器接口

**功能说明:**
  - 实现types.Command接口
  - 返回命令的标志注册器
  - 用于直接操作标志注册
  - 支持并发安全的访问

#### func (c *Cmd) Flags() []types.Flag

```go
func (c *Cmd) Flags() []types.Flag
```

Flags 获取所有标志

**返回值:**
  - []types.Flag: 所有标志的切片副本

**功能说明:**
  - 实现types.Command接口
  - 返回所有注册的标志
  - 创建副本避免外部修改
  - 支持并发安全的访问

#### func (c *Cmd) Float64(longName, shortName, description string, default_ float64) *flag.Float64Flag

```go
func (c *Cmd) Float64(longName, shortName, description string, default_ float64) *flag.Float64Flag
```

Float64 创建64位浮点数标志

**参数:**
  - longName: 长标志名 (如 --long-name)
  - shortName: 短标志名 (如 -s)
  - description: 标志的描述信息
  - default_: 标志的默认值

**返回值:**
  - *flag.Float64Flag: 新创建的64位浮点数标志

#### func (c *Cmd) GetFlag(name string) (types.Flag, bool)

```go
func (c *Cmd) GetFlag(name string) (types.Flag, bool)
```

GetFlag 根据名称获取标志

**参数:**
  - name: 标志名称

**返回值:**
  - types.Flag: 找到的标志
  - bool: 是否找到, true表示找到

**功能说明:**
  - 实现types.Command接口
  - 从标志注册器中查找标志
  - 支持并发安全的查找操作

#### func (c *Cmd) MutexGroups() []types.MutexGroup

```go
func (c *Cmd) MutexGroups() []types.MutexGroup
```

MutexGroups 获取命令的所有互斥组

**返回值:**
  - []types.MutexGroup: 互斥组列表的副本

**功能说明:**
  - 返回命令中定义的所有互斥组
  - 返回副本以防止外部修改内部状态
  - 使用读锁保护并发安全

#### func (c *Cmd) GetSubCmd(name string) (types.Command, bool)

```go
func (c *Cmd) GetSubCmd(name string) (types.Command, bool)
```

GetSubCmd 根据名称获取子命令

**参数:**
  - name: 子命令名称

**返回值:**
  - types.Command: 找到的子命令
  - bool: 是否找到, true表示找到

**功能说明:**
  - 实现types.Command接口
  - 从子命令注册器中查找
  - 支持并发安全的查找操作

#### func (c *Cmd) HasRunFunc() bool

```go
func (c *Cmd) HasRunFunc() bool
```

HasRunFunc 检查是否设置了运行函数

**返回值:**
  - bool: 是否设置了运行函数, true表示已设置

**功能说明:**
  - 实现types.Command接口
  - 线程安全地检查运行函数
  - 用于执行前的状态检查

#### func (c *Cmd) HasSubCmd(name string) bool

```go
func (c *Cmd) HasSubCmd(name string) bool
```

HasSubCmd 检查是否存在指定名称的子命令

**参数:**
  - name: 子命令名称

**返回值:**
  - bool: 是否存在, true表示存在

**功能说明:**
  - 实现types.Command接口
  - 快速检查子命令存在性
  - 支持并发安全的检查

#### func (c *Cmd) IsDisableFlagParsing() bool

```go
func (c *Cmd) IsDisableFlagParsing() bool
```

IsDisableFlagParsing 检查是否禁用标志解析

**返回值:**
  - bool: 如果禁用标志解析返回 true，否则返回 false

**功能说明:**
  - 实现types.Command接口
  - 获取命令的禁用标志解析状态
  - 线程安全地访问状态字段
  - 默认值为 false

**使用场景:**
  - 检查命令是否配置了禁用标志解析
  - 在解析前进行状态确认
  - 调试和日志输出

---

#### func (c *Cmd) SetDisableFlagParsing(disable bool)

```go
func (c *Cmd) SetDisableFlagParsing(disable bool)
```

SetDisableFlagParsing 设置是否禁用标志解析

**参数:**
  - disable: 是否禁用标志解析，true 表示禁用，false 表示不禁用

**功能说明:**
  - 实现types.Command接口
  - 设置命令的禁用标志解析状态
  - 线程安全地修改状态字段
  - 应在解析前设置

**使用示例:**
```go
cmd := NewCmd("exec", "e", ExitOnError)
cmd.SetDisableFlagParsing(true)  // 禁用标志解析
cmd.SetRun(func(c types.Command) error {
    args := c.Args()  // 所有参数都作为位置参数
    // 处理透传的参数
    return nil
})
```

**注意事项:**
  - 应在解析前设置，通常在命令创建后立即设置
  - 每个命令可以独立设置，父命令的设置不影响子命令
  - 禁用后，环境变量绑定也会被跳过
  - 内置标志（如 --help）也不会被特殊处理

#### func (c *Cmd) IsHidden() bool

```go
func (c *Cmd) IsHidden() bool
```

IsHidden 检查命令是否隐藏

**返回值:**
  - bool: 如果命令是隐藏的返回 true，否则返回 false

**功能说明:**
  - 实现types.Command接口
  - 获取命令的隐藏状态
  - 线程安全地访问状态字段
  - 默认值为 false

**使用场景:**
  - 检查命令是否配置了隐藏
  - 在生成帮助信息时过滤隐藏命令
  - 调试和日志输出

---

#### func (c *Cmd) SetHidden(hidden bool)

```go
func (c *Cmd) SetHidden(hidden bool)
```

SetHidden 设置命令是否隐藏

**参数:**
  - hidden: 是否隐藏命令，true 表示隐藏，false 表示不隐藏

**功能说明:**
  - 实现types.Command接口
  - 设置命令的隐藏状态
  - 线程安全地修改状态字段
  - 隐藏命令不会显示在帮助信息中

**使用示例:**
```go
cmd := NewCmd("debug", "d", ExitOnError)
cmd.SetHidden(true)  // 隐藏调试命令
cmd.SetRun(func(c types.Command) error {
    // 执行调试逻辑
    fmt.Println("调试模式")
    return nil
})
```

**注意事项:**
  - 隐藏命令仍可通过命令行正常调用
  - 只是不在帮助信息的子命令列表中显示
  - 适用于内部命令、调试命令或高级功能
  - 子命令可以独立设置隐藏状态

#### func (c *Cmd) Help() string

```go
func (c *Cmd) Help() string
```

Help 获取帮助信息

**返回值:**
  - string: 格式化的帮助信息

**功能说明:**
  - 实现types.Command接口
  - 使用help包生成帮助信息
  - 包含标志、子命令和示例
  - 支持并发安全的访问

#### func (c *Cmd) Int(longName, shortName, description string, default_ int) *flag.IntFlag

```go
func (c *Cmd) Int(longName, shortName, description string, default_ int) *flag.IntFlag
```

Int 创建整数标志

**参数:**
  - longName: 长标志名 (如 --long-name)
  - shortName: 短标志名 (如 -s)
  - description: 标志的描述信息
  - default_: 标志的默认值

**返回值:**
  - *flag.IntFlag: 新创建的整数标志

#### func (c *Cmd) Int64(longName, shortName, description string, default_ int64) *flag.Int64Flag

```go
func (c *Cmd) Int64(longName, shortName, description string, default_ int64) *flag.Int64Flag
```

Int64 创建64位整数标志

**参数:**
  - longName: 长标志名 (如 --long-name)
  - shortName: 短标志名 (如 -s)
  - description: 标志的描述信息
  - default_: 标志的默认值

**返回值:**
  - *flag.Int64Flag: 新创建的64位整数标志

#### func (c *Cmd) Int64Slice(longName, shortName, description string, default_ []int64) *flag.Int64SliceFlag

```go
func (c *Cmd) Int64Slice(longName, shortName, description string, default_ []int64) *flag.Int64SliceFlag
```

Int64Slice 创建64位整数切片标志

**参数:**
  - longName: 长标志名 (如 --long-name)
  - shortName: 短标志名 (如 -s)
  - description: 标志的描述信息
  - default_: 标志的默认值

**返回值:**
  - *flag.Int64SliceFlag: 新创建的64位整数切片标志

#### func (c *Cmd) IntSlice(longName, shortName, description string, default_ []int) *flag.IntSliceFlag

```go
func (c *Cmd) IntSlice(longName, shortName, description string, default_ []int) *flag.IntSliceFlag
```

IntSlice 创建整数切片标志

**参数:**
  - longName: 长标志名 (如 --long-name)
  - shortName: 短标志名 (如 -s)
  - description: 标志的描述信息
  - default_: 标志的默认值

**返回值:**
  - *flag.IntSliceFlag: 新创建的整数切片标志

#### func (c *Cmd) IsParsed() bool

```go
func (c *Cmd) IsParsed() bool
```

IsParsed 检查命令是否已解析

**返回值:**
  - bool: 是否已解析, true表示已解析

**功能说明:**
  - 实现types.Command接口
  - 线程安全地检查解析状态
  - 用于避免重复解析

#### func (c *Cmd) IsRootCmd() bool

```go
func (c *Cmd) IsRootCmd() bool
```

IsRootCmd 检查是否为根命令

**返回值:**
  - bool: 是否为根命令, true表示是根命令

**功能说明:**
  - 实现types.Command接口
  - 通过检查父命令判断
  - 根命令没有父命令
  - 支持并发安全的检查

#### func (c *Cmd) LongName() string

```go
func (c *Cmd) LongName() string
```

LongName 获取命令长名称

**返回值:**
  - string: 命令的长名称

**功能说明:**
  - 实现types.Command接口
  - 线程安全地访问长名称
  - 用于命令的完整标识

#### func (c *Cmd) Map(longName, shortName, description string, default_ map[string]string) *flag.MapFlag

```go
func (c *Cmd) Map(longName, shortName, description string, default_ map[string]string) *flag.MapFlag
```

Map 创建映射标志

**参数:**
  - longName: 长标志名 (如 --long-name)
  - shortName: 短标志名 (如 -s)
  - description: 标志的描述信息
  - default_: 标志的默认值

**返回值:**
  - *flag.MapFlag: 新创建的映射标志

#### func (c *Cmd) NArg() int

```go
func (c *Cmd) NArg() int
```

NArg 获取命令行参数数量

**返回值:**
  - int: 参数数量

**功能说明:**
  - 实现types.Command接口
  - 线程安全地获取参数数量
  - 用于参数范围检查

#### func (c *Cmd) Name() string

```go
func (c *Cmd) Name() string
```

Name 获取命令名称

**返回值:**
  - string: 命令的名称, 优先返回长名称

**功能说明:**
  - 实现types.Command接口
  - 优先返回长名称, 为空时返回短名称
  - 用作命令的主要标识符

#### func (c *Cmd) Parse(args []string) error

```go
func (c *Cmd) Parse(args []string) error
```

Parse 解析命令行参数（可重复解析）

**参数:**
  - args: 命令行参数列表

**返回值:**
  - error: 解析失败时返回错误

**功能说明:**
  - 实现types.Command接口
  - 可以重复调用，会覆盖之前的解析结果
  - 调用解析器的Parse方法
  - 递归解析所有子命令

**注意事项:**
  - 重复解析会覆盖之前的解析结果
  - 如果需要确保只解析一次，请使用 ParseOnce

#### func (c *Cmd) ParseOnce(args []string) error

```go
func (c *Cmd) ParseOnce(args []string) error
```

ParseOnce 解析命令行参数（只解析一次）

**参数:**
  - args: 命令行参数列表

**返回值:**
  - error: 解析失败时返回错误

**功能说明:**
  - 使用sync.Once确保只解析一次
  - 重复执行无错误、仅首次执行解析
  - 调用解析器的Parse方法
  - 递归解析所有子命令

**注意事项:**
  - 如果需要重复解析，请使用 Parse 方法
  - 建议在普通场景使用此方法，避免误用

#### func (c *Cmd) ParseAndRoute(args []string) error

```go
func (c *Cmd) ParseAndRoute(args []string) error
```

ParseAndRoute 解析并路由执行命令（可重复解析）

**参数:**
  - args: 命令行参数列表

**返回值:**
  - error: 解析或执行失败时返回错误

**功能说明:**
  - 实现types.Command接口
  - 可以重复调用，会覆盖之前的解析结果
  - 调用解析器的ParseAndRoute方法
  - 完整的解析和执行流程

**注意事项:**
  - 重复解析会覆盖之前的解析结果
  - 如果需要确保只解析一次，请使用 ParseAndRouteOnce

#### func (c *Cmd) ParseAndRouteOnce(args []string) error

```go
func (c *Cmd) ParseAndRouteOnce(args []string) error
```

ParseAndRouteOnce 解析并路由执行命令（只解析一次）

**参数:**
  - args: 命令行参数列表

**返回值:**
  - error: 解析或执行失败时返回错误

**功能说明:**
  - 使用sync.Once确保只执行一次
  - 重复执行无错误、仅首次执行解析
  - 调用解析器的ParseAndRoute方法
  - 完整的解析和执行流程

**注意事项:**
  - 如果需要重复解析，请使用 ParseAndRoute 方法
  - 建议在普通场景使用此方法，避免误用

#### func (c *Cmd) ParseOnly(args []string) error

```go
func (c *Cmd) ParseOnly(args []string) error
```

ParseOnly 仅解析当前命令, 不递归解析子命令（可重复解析）

**参数:**
  - args: 命令行参数列表

**返回值:**
  - error: 解析失败时返回错误

**功能说明:**
  - 实现types.Command接口
  - 可以重复调用，会覆盖之前的解析结果
  - 调用解析器的ParseOnly方法
  - 不处理子命令解析

**注意事项:**
  - 重复解析会覆盖之前的解析结果
  - 如果需要确保只解析一次，请使用 ParseOnlyOnce

#### func (c *Cmd) ParseOnlyOnce(args []string) error

```go
func (c *Cmd) ParseOnlyOnce(args []string) error
```

ParseOnlyOnce 仅解析当前命令, 不递归解析子命令（只解析一次）

**参数:**
  - args: 命令行参数列表

**返回值:**
  - error: 解析失败时返回错误

**功能说明:**
  - 使用sync.Once确保只执行一次
  - 重复执行无错误、仅首次执行解析
  - 调用解析器的ParseOnly方法
  - 不处理子命令解析

**注意事项:**
  - 如果需要重复解析，请使用 ParseOnly 方法
  - 建议在普通场景使用此方法，避免误用

#### func (c *Cmd) Path() string

```go
func (c *Cmd) Path() string
```

Path 获取命令路径

**返回值:**
  - string: 完整的命令路径

**功能说明:**
  - 实现types.Command接口
  - 递归构建完整路径
  - 格式: 父路径 + 空格 + 命令名
  - 根命令直接返回名称
  - 用于帮助信息和错误显示
  - 用于帮助信息和错误显示

#### func (c *Cmd) PrintHelp()

```go
func (c *Cmd) PrintHelp()
```

PrintHelp 打印帮助信息

**功能说明:**
  - 实现types.Command接口
  - 直接输出帮助信息到控制台
  - 使用标准fmt包输出
  - 支持并发安全的访问

#### func (c *Cmd) AddRequiredGroup(name string, flags []string, conditional bool) error

```go
func (c *Cmd) AddRequiredGroup(name string, flags []string, conditional bool) error
```

AddRequiredGroup 添加必需组到命令

**参数:**
  - name: 必需组名称, 用于错误提示和标识
  - flags: 必需组中的标志名称列表
  - conditional: 是否为条件性必需组

**返回值:**
  - error: 添加失败时返回错误

**功能说明:**
  - 创建新的必需组并添加到命令配置中
  - 对于普通必需组, 组中的所有标志都必须被设置
  - 对于条件性必需组, 如果组中任何一个标志被设置, 则所有标志都必须被设置
  - 使用写锁保护并发安全

**错误情况:**
  - 组名已存在: 返回 REQUIRED_GROUP_ALREADY_EXISTS 错误
  - 标志列表为空: 返回 EMPTY_REQUIRED_GROUP 错误
  - 标志不存在: 返回 FLAG_NOT_FOUND 错误

**注意事项:**
  - 标志名称必须是已注册的标志
  - 必需组名称在命令中应该唯一
  - 重复添加同名必需组会返回错误
  - 条件性必需组提供了更灵活的标志验证方式, 适用于某些可选但相关的标志组合

#### func (c *Cmd) RequiredGroups() []types.RequiredGroup

```go
func (c *Cmd) RequiredGroups() []types.RequiredGroup
```

RequiredGroups 获取命令的所有必需组

**返回值:**
  - []types.RequiredGroup: 必需组列表的副本

**功能说明:**
  - 返回命令中定义的所有必需组
  - 返回副本以防止外部修改内部状态
  - 使用读锁保护并发安全

#### func (c *Cmd) AddFlagDependency(name, trigger string, targets []string, depType types.DepType) error

```go
func (c *Cmd) AddFlagDependency(name, trigger string, targets []string, depType types.DepType) error
```

AddFlagDependency 添加标志依赖关系到命令

**参数:**
  - name: 依赖关系名称, 用于错误提示和标识
  - trigger: 触发标志名称, 当此标志被设置时触发依赖检查
  - targets: 目标标志名称列表, 这些标志会受到约束
  - depType: 依赖关系类型 (types.DepMutex 或 types.DepRequired)

**返回值:**
  - error: 添加失败时返回错误

**功能说明:**
  - 创建新的标志依赖关系并添加到命令配置中
  - DepMutex: 当触发标志被设置时, 目标标志不能被设置
  - DepRequired: 当触发标志被设置时, 所有目标标志必须被设置
  - 使用写锁保护并发安全

**错误情况:**
  - 名称为空: 返回错误
  - 名称已存在: 返回错误
  - 触发标志为空: 返回错误
  - 目标标志列表为空: 返回错误
  - 触发标志自依赖: 返回错误 (触发标志不能在目标列表中)
  - 触发标志不存在: 返回错误
  - 目标标志不存在: 返回错误

**使用场景:**
  - 远程模式与本地路径互斥 (trigger="remote", targets=["local-path"], type=DepMutex)
  - SSL模式需要证书和密钥 (trigger="ssl", targets=["cert","key"], type=DepRequired)
  - 配置文件模式与其他配置互斥 (trigger="config", targets=["port","host"], type=DepMutex)

**注意事项:**
  - 依赖关系是单向的, 从触发标志指向目标标志
  - 触发标志和目标标志必须都已注册到命令
  - 依赖关系名称在命令中应该唯一

#### func (c *Cmd) FlagDependencies() []types.FlagDependency

```go
func (c *Cmd) FlagDependencies() []types.FlagDependency
```

FlagDependencies 获取命令的所有标志依赖关系

**返回值:**
  - []types.FlagDependency: 标志依赖关系列表的副本

**功能说明:**
  - 返回命令中定义的所有标志依赖关系
  - 返回副本以防止外部修改内部状态
  - 使用读锁保护并发安全

#### func (c *Cmd) Run() error

```go
func (c *Cmd) Run() error
```

Run 执行命令

**返回值:**
  - error: 执行失败时返回错误

**功能说明:**
  - 实现types.Command接口
  - 检查解析状态和运行函数
  - 调用设置的运行函数
  - 支持并发安全的执行

**错误情况:**
  - 未解析: 返回解析错误
  - 无运行函数: 返回运行函数错误

#### func (c *Cmd) SetArgs(args []string)

```go
func (c *Cmd) SetArgs(args []string)
```

SetArgs 设置命令行参数

**参数:**
  - args: 命令行参数列表

**功能说明:**
  - 手动设置命令行参数
  - 通常由解析器调用
  - 空切片被忽略
  - 支持并发安全的设置

#### func (c *Cmd) SetChinese(useChinese bool)

```go
func (c *Cmd) SetChinese(useChinese bool)
```

SetChinese 设置是否使用中文

**参数:**
  - useChinese: 是否使用中文

**功能说明:**
  - 设置帮助信息的语言
  - 影响错误消息和提示
  - 存储在配置中
  - 支持并发安全的设置

#### func (c *Cmd) SetCompletion(enable bool)

```go
func (c *Cmd) SetCompletion(enable bool)
```

SetCompletion 设置是否启用自动补全标志

**参数:**
  - enable: 是否启用自动补全标志

**功能说明:**
  - 控制是否注册 --completion 标志
  - 默认为 false，不启用自动补全
  - 存储在配置中
  - 支持并发安全的设置

#### func (c *Cmd) SetDynamicCompletion(enable bool)

```go
func (c *Cmd) SetDynamicCompletion(enable bool)
```

SetDynamicCompletion 设置是否启用动态补全

**参数:**
  - enable: 是否启用动态补全

**功能说明:**
  - 启用后会自动注册 __complete 子命令
  - 用于 Shell 自动补全脚本动态获取补全信息
  - 支持并发安全的设置

**前置条件:**
  - 必须先启用 `Completion`（自动补全标志），否则会 panic

**注意事项:**
  - 如果未启用 `Completion` 就直接调用，会 panic: "dynamic completion cannot be enabled when completion is disabled"
  - 如果启用时创建 __complete 子命令失败也会 panic
  - 只能在根命令上启用动态补全, 子命令上无效
  - 建议在解析前设置

**使用示例:**
```go
cmd := NewCmd("myapp", "", ExitOnError)
cmd.SetCompletion(true)              // 先启用自动补全
cmd.SetDynamicCompletion(true) // 再启用动态补全
```

#### func (c *Cmd) SubCmds() []types.Command

```go
func (c *Cmd) SubCmds() []types.Command
```

SubCmds 获取所有非隐藏子命令

**返回值:**
  - []types.Command: 非隐藏子命令的切片

**功能说明:**
  - 实现types.Command接口
  - 返回所有非隐藏的子命令列表
  - 自动过滤隐藏的命令（如 __complete）
  - 支持并发安全的访问

**使用场景:**
  - 获取命令的子命令列表
  - 帮助信息生成
  - 补全候选计算
  - 子命令遍历操作

**注意事项:**
  - 隐藏命令（通过 SetHidden(true) 设置）不会出现在返回结果中
  - 如需访问隐藏命令，请使用 GetSubCmd(name) 按名称获取
  - 返回的切片是副本，修改不影响原命令状态

#### func (c *Cmd) SetDesc(desc string)

```go
func (c *Cmd) SetDesc(desc string)
```

SetDesc 设置命令描述

**参数:**
  - desc: 命令描述信息

**功能说明:**
  - 实现types.Command接口
  - 设置命令的功能描述
  - 用于帮助信息显示
  - 支持并发安全的设置

#### func (c *Cmd) SetEnvPrefix(prefix string)

```go
func (c *Cmd) SetEnvPrefix(prefix string)
```

SetEnvPrefix 设置环境变量前缀

**参数:**
  - prefix: 环境变量前缀

**功能说明:**
  - 设置环境变量的前缀
  - 自动添加下划线后缀
  - 空字符串表示不使用前缀
  - 支持并发安全的设置

#### func (c *Cmd) AutoBindAllEnv()

```go
func (c *Cmd) AutoBindAllEnv()
```

AutoBindAllEnv 为所有标志自动绑定环境变量

**功能说明:**
  - 遍历命令的所有标志
  - 为每个标志调用 AutoBindEnv() 方法
  - 批量设置环境变量绑定

**使用示例:**
```go
cmd.String("host", "h", "主机地址", "localhost")
cmd.Uint("port", "p", "端口号", 8080)
cmd.AutoBindAllEnv()  // 自动绑定 HOST 和 PORT
```

**注意事项:**
  - 如果标志没有长名称，会触发 panic
  - 环境变量名为标志长名称的大写形式
  - 环境变量前缀在解析时自动拼接

#### func (c *Cmd) SetLogoText(logo string)

```go
func (c *Cmd) SetLogoText(logo string)
```

SetLogoText 设置Logo文本

**参数:**
  - logo: Logo文本内容

**功能说明:**
  - 设置命令的Logo
  - 用于帮助信息显示
  - 存储在配置中
  - 支持并发安全的设置

#### func (c *Cmd) SetParsed(parsed bool)

```go
func (c *Cmd) SetParsed(parsed bool)
```

SetParsed 设置解析状态

**参数:**
  - parsed: 解析状态

**功能说明:**
  - 手动设置解析状态
  - 通常由解析器调用
  - 影响后续操作的行为
  - 支持并发安全的设置

#### func (c *Cmd) SetParser(p types.Parser)

```go
func (c *Cmd) SetParser(p types.Parser)
```

SetParser 设置命令的解析器

**参数:**
  - p: 解析器接口实现

**功能说明:**
  - 替换默认的解析器
  - 允许自定义解析逻辑
  - nil值会触发panic