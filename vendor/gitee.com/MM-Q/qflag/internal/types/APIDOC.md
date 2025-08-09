# Package types

**Import Path:** `gitee.com/MM-Q/qflag/internal/types`

Package types 内置类型和数据结构定义。本包定义了qflag包内部使用的内置类型和数据结构，包括内置标志、配置选项等核心数据类型的定义。

## 功能模块

- **内置类型和数据结构定义** - 定义了qflag包内部使用的内置类型和数据结构，包括内置标志、配置选项等核心数据类型的定义
- **配置结构体和选项定义** - 定义了命令配置相关的结构体和选项，包括命令的各种配置参数、帮助信息设置、版本信息等配置数据的定义和管理
- **命令上下文和状态管理** - 定义了命令上下文结构体，用于管理命令的状态、子命令、标志注册表等信息，提供命令执行过程中的状态维护和数据共享功能

## 目录

- [类型](#类型)
  - [BuiltinFlags](#builtinflags)
  - [CmdConfig](#cmdconfig)
  - [CmdContext](#cmdcontext)
  - [ExampleInfo](#exampleinfo)

## 类型

### BuiltinFlags

```go
type BuiltinFlags struct {
    Help       *flags.BoolFlag // 标志-帮助
    Version    *flags.BoolFlag // 标志-版本
    Completion *flags.EnumFlag // 标志-自动完成
    NameMap    sync.Map        // 内置标志名称映射
}
```

BuiltinFlags 内置标志结构体，管理系统预定义的标志。

**字段说明:**
- `Help`: 帮助标志，用于显示帮助信息
- `Version`: 版本标志，用于显示版本信息
- `Completion`: 自动完成标志，用于生成shell自动完成脚本
- `NameMap`: 内置标志名称映射表，用于快速查找

#### 构造函数

##### NewBuiltinFlags

```go
func NewBuiltinFlags() *BuiltinFlags
```

NewBuiltinFlags 创建内置标志实例。

**返回值:**
- `*BuiltinFlags`: 新创建的内置标志实例

**功能:**
- 初始化所有内置标志
- 设置默认的标志配置
- 建立名称映射关系

#### 方法

##### IsBuiltinFlag

```go
func (bf *BuiltinFlags) IsBuiltinFlag(name string) bool
```

IsBuiltinFlag 检查是否为内置标志。

**参数:**
- `name`: 标志名称

**返回值:**
- `bool`: 是否为内置标志

**用途:**
- 区分用户定义标志和系统内置标志
- 在解析过程中进行特殊处理
- 避免名称冲突

##### MarkAsBuiltin

```go
func (bf *BuiltinFlags) MarkAsBuiltin(names ...string)
```

MarkAsBuiltin 标记为内置标志。

**参数:**
- `names`: 标志名称列表（可变参数）

**功能:**
- 将指定名称标记为内置标志
- 更新内部映射表
- 支持批量标记

### CmdConfig

```go
type CmdConfig struct {
    // 版本信息
    Version string

    // 自定义描述
    Description string

    // 自定义的完整命令行帮助信息
    Help string

    // 自定义用法格式说明
    UsageSyntax string

    // 模块帮助信息
    ModuleHelps string

    // logo文本
    LogoText string

    // 备注内容切片
    Notes []string

    // 示例信息切片
    Examples []ExampleInfo

    // 是否使用中文帮助信息
    UseChinese bool

    // 控制内置标志是否自动退出
    ExitOnBuiltinFlags bool

    // 控制是否启用自动补全功能
    EnableCompletion bool
}
```

CmdConfig 命令行配置结构体，包含命令的所有配置选项。

**配置分类:**

**基本信息:**
- `Version`: 应用程序版本号
- `Description`: 命令描述信息
- `Help`: 自定义帮助信息
- `UsageSyntax`: 用法语法说明

**扩展信息:**
- `ModuleHelps`: 模块帮助信息
- `LogoText`: 应用程序Logo文本
- `Notes`: 注意事项列表
- `Examples`: 使用示例列表

**行为控制:**
- `UseChinese`: 是否使用中文帮助信息
- `ExitOnBuiltinFlags`: 内置标志是否自动退出程序
- `EnableCompletion`: 是否启用自动补全功能

#### 构造函数

##### NewCmdConfig

```go
func NewCmdConfig() *CmdConfig
```

NewCmdConfig 创建一个新的CmdConfig实例。

**返回值:**
- `*CmdConfig`: 新创建的配置实例

**默认配置:**
- 启用中文帮助信息
- 启用内置标志自动退出
- 启用自动补全功能

### CmdContext

```go
type CmdContext struct {
    // 长命令名称
    LongName string
    // 短命令名称
    ShortName string

    // 标志注册表, 统一管理标志的元数据
    FlagRegistry *flags.FlagRegistry
    // 底层flag集合, 处理参数解析
    FlagSet *flag.FlagSet

    // 命令行参数(非标志参数)
    Args []string
    // 是否已经解析过参数
    Parsed atomic.Bool
    // 用于确保参数解析只执行一次
    ParseOnce sync.Once
    // 读写锁
    Mutex sync.RWMutex

    // 子命令上下文切片
    SubCmds []*CmdContext
    // 子命令映射表
    SubCmdMap map[string]*CmdContext
    // 父命令上下文
    Parent *CmdContext

    // 配置信息
    Config *CmdConfig

    // 内置标志结构体
    BuiltinFlags *BuiltinFlags

    // ParseHook 解析阶段钩子函数
    // 在标志解析完成后、子命令参数处理后调用
    //
    // 参数:
    //   - 当前命令上下文
    //
    // 返回值:
    //   - error: 错误信息, 非nil时会中断解析流程
    //   - bool: 是否需要退出程序
    ParseHook func(*CmdContext) (error, bool)
}
```

CmdContext 命令上下文，包含所有必要的状态信息。这是所有函数操作的核心数据结构。

**结构分类:**

**基本标识:**
- `LongName`: 命令的长名称
- `ShortName`: 命令的短名称

**标志管理:**
- `FlagRegistry`: 标志注册表，管理标志元数据
- `FlagSet`: 底层标志集合，处理参数解析

**解析状态:**
- `Args`: 解析后的非标志参数
- `Parsed`: 原子布尔值，标记是否已解析
- `ParseOnce`: 确保解析只执行一次
- `Mutex`: 读写锁，保证并发安全

**层级结构:**
- `SubCmds`: 子命令上下文列表
- `SubCmdMap`: 子命令映射表，用于快速查找
- `Parent`: 父命令上下文引用

**配置和扩展:**
- `Config`: 命令配置信息
- `BuiltinFlags`: 内置标志管理
- `ParseHook`: 解析钩子函数

#### 构造函数

##### NewCmdContext

```go
func NewCmdContext(longName, shortName string, errorHandling flag.ErrorHandling) *CmdContext
```

NewCmdContext 创建新的命令上下文。

**参数:**
- `longName`: 长命令名称
- `shortName`: 短命令名称
- `errorHandling`: 错误处理方式

**返回值:**
- `*CmdContext`: 新创建的命令上下文

**错误处理选项:**
- `flag.ContinueOnError`: 解析标志时遇到错误继续解析，并返回错误信息
- `flag.ExitOnError`: 解析标志时遇到错误立即退出程序，并返回错误信息
- `flag.PanicOnError`: 解析标志时遇到错误立即触发panic

#### 方法

##### GetName

```go
func (ctx *CmdContext) GetName() string
```

GetName 获取命令名称。如果长命令名称不为空则返回长命令名称，否则返回短命令名称。

**返回值:**
- `string`: 命令名称

**逻辑:**
- 优先返回长名称
- 长名称为空时返回短名称
- 用于显示和日志记录

### ExampleInfo

```go
type ExampleInfo struct {
    Description string // 示例描述
    Usage       string // 示例使用方式
}
```

ExampleInfo 示例信息结构体，用于存储命令的使用示例，包括描述和示例内容。

**字段:**
- `Description`: 示例描述，说明示例的用途和场景
- `Usage`: 示例使用方式，具体的命令行用法

**用途:**
- 在帮助信息中展示使用示例
- 提供用户参考和学习材料
- 支持多个示例的组织和管理

## 使用示例

### 创建基本命令上下文

```go
package main

import (
    "flag"
    "gitee.com/MM-Q/qflag/internal/types"
)

func main() {
    // 创建命令上下文
    ctx := types.NewCmdContext("myapp", "app", flag.ExitOnError)
    
    // 设置配置
    ctx.Config = types.NewCmdConfig()
    ctx.Config.Version = "1.0.0"
    ctx.Config.Description = "我的应用程序"
    
    // 使用上下文...
}
```

### 配置命令信息

```go
func setupCommand() *types.CmdContext {
    ctx := types.NewCmdContext("deploy", "d", flag.ContinueOnError)
    
    // 配置基本信息
    config := types.NewCmdConfig()
    config.Version = "2.1.0"
    config.Description = "部署应用程序到服务器"
    config.UsageSyntax = "deploy [选项] <目标环境>"
    
    // 添加注意事项
    config.Notes = []string{
        "部署前请确保目标环境可访问",
        "建议在部署前进行备份",
    }
    
    // 添加使用示例
    config.Examples = []types.ExampleInfo{
        {
            Description: "部署到生产环境",
            Usage:       "deploy --env production --config prod.yaml",
        },
        {
            Description: "部署到测试环境",
            Usage:       "deploy --env staging --dry-run",
        },
    }
    
    ctx.Config = config
    return ctx
}
```

### 管理子命令

```go
func setupSubCommands(parent *types.CmdContext) {
    // 创建子命令
    startCmd := types.NewCmdContext("start", "s", flag.ExitOnError)
    startCmd.Config = types.NewCmdConfig()
    startCmd.Config.Description = "启动服务"
    startCmd.Parent = parent
    
    stopCmd := types.NewCmdContext("stop", "st", flag.ExitOnError)
    stopCmd.Config = types.NewCmdConfig()
    stopCmd.Config.Description = "停止服务"
    stopCmd.Parent = parent
    
    // 添加到父命令
    parent.SubCmds = append(parent.SubCmds, startCmd, stopCmd)
    
    // 建立映射关系
    if parent.SubCmdMap == nil {
        parent.SubCmdMap = make(map[string]*types.CmdContext)
    }
    parent.SubCmdMap["start"] = startCmd
    parent.SubCmdMap["s"] = startCmd
    parent.SubCmdMap["stop"] = stopCmd
    parent.SubCmdMap["st"] = stopCmd
}
```

### 使用内置标志

```go
func setupBuiltinFlags(ctx *types.CmdContext) {
    // 创建内置标志
    ctx.BuiltinFlags = types.NewBuiltinFlags()
    
    // 检查是否为内置标志
    if ctx.BuiltinFlags.IsBuiltinFlag("help") {
        fmt.Println("help 是内置标志")
    }
    
    // 标记自定义标志为内置
    ctx.BuiltinFlags.MarkAsBuiltin("debug", "verbose")
}
```

### 设置解析钩子

```go
func setupParseHook(ctx *types.CmdContext) {
    ctx.ParseHook = func(context *types.CmdContext) (error, bool) {
        // 在解析完成后执行自定义逻辑
        fmt.Printf("解析完成，命令: %s\n", context.GetName())
        
        // 检查必需的参数
        if len(context.Args) == 0 {
            return fmt.Errorf("缺少必需的参数"), false
        }
        
        // 执行验证逻辑
        if err := validateArgs(context.Args); err != nil {
            return err, false
        }
        
        // 正常继续执行
        return nil, false
    }
}

func validateArgs(args []string) error {
    // 自定义验证逻辑
    for _, arg := range args {
        if arg == "" {
            return fmt.Errorf("参数不能为空")
        }
    }
    return nil
}
```

## 设计特点

1. **类型安全** - 使用强类型定义，避免运行时错误
2. **并发安全** - 使用原子操作和锁机制保证并发安全
3. **层级结构** - 支持命令和子命令的层级组织
4. **配置驱动** - 通过配置结构体灵活控制行为
5. **扩展性强** - 提供钩子函数支持自定义逻辑

## 内存管理

- 上下文结构体使用指针传递，避免大结构体复制
- 子命令映射使用map提供O(1)查找性能
- 原子操作避免锁竞争，提高并发性能
- 合理的字段布局减少内存对齐开销

## 最佳实践

1. **初始化顺序**: 先创建上下文，再设置配置，最后注册标志
2. **错误处理**: 根据应用场景选择合适的错误处理策略
3. **并发使用**: 在多goroutine环境中注意使用读写锁
4. **资源清理**: 及时清理不再使用的上下文引用
5. **配置管理**: 使用配置结构体统一管理命令行为

## 注意事项

- 命令上下文创建后不应修改基本标识字段
- 解析状态字段由系统管理，不应手动修改
- 子命令映射表需要与子命令列表保持同步
- 钩子函数应避免长时间阻塞操作
- 配置修改应在解析之前完成