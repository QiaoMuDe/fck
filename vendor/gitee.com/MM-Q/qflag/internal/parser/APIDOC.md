# Package parser

**Import Path:** `gitee.com/MM-Q/qflag/internal/parser`

Package parser 命令行参数解析器。本包实现了命令行参数的解析逻辑和环境变量处理功能，为命令行应用程序提供完整的参数解析支持。

## 功能模块

- **环境变量解析和处理** - 实现了环境变量的解析和处理逻辑，支持从环境变量中读取标志值，为命令行参数提供环境变量绑定和默认值设置功能
- **命令行参数解析器** - 实现了命令行参数的解析逻辑，包括标志解析、参数分离、子命令识别等核心解析功能，为命令行参数处理提供基础支持

## 目录

- [函数](#函数)
  - [LoadEnvVars](#loadenvvars)
  - [ParseCommand](#parsecommand)

## 函数

### LoadEnvVars

```go
func LoadEnvVars(ctx *types.CmdContext) error
```

LoadEnvVars 从环境变量加载参数值。纯函数设计，不依赖结构体状态。

**参数:**
- `ctx`: 命令上下文，包含需要从环境变量加载值的标志信息

**返回值:**
- `error`: 如果加载过程中出现错误，返回错误信息；成功时返回 nil

**功能特点:**
- 纯函数设计，无副作用
- 支持从环境变量中读取标志值
- 提供环境变量绑定功能
- 支持默认值设置

**使用场景:**
- 在命令行解析之前预加载环境变量值
- 为标志提供环境变量作为默认值来源
- 支持配置文件和环境变量的混合使用

### ParseCommand

```go
func ParseCommand(ctx *types.CmdContext, args []string) (err error)
```

ParseCommand 解析单个命令的标志和参数。

**参数:**
- `ctx`: 命令上下文，用于存储解析结果和配置信息
- `args`: 命令行参数切片，包含需要解析的所有参数

**返回值:**
- `error`: 如果解析失败，返回错误信息；成功时返回 nil

**功能特点:**
- 支持标志解析（长选项和短选项）
- 实现参数分离逻辑
- 支持子命令识别
- 提供核心解析功能

**解析能力:**
- **标志解析**: 识别和处理 `--flag` 和 `-f` 格式的选项
- **参数分离**: 区分标志参数和位置参数
- **子命令识别**: 识别和路由到相应的子命令
- **值绑定**: 将解析的值绑定到相应的标志变量

## 使用示例

### 基本解析流程

```go
package main

import (
    "os"
    "gitee.com/MM-Q/qflag/internal/parser"
    "gitee.com/MM-Q/qflag/internal/types"
)

func main() {
    // 创建命令上下文
    ctx := &types.CmdContext{
        // ... 初始化命令信息
    }
    
    // 1. 首先加载环境变量
    if err := parser.LoadEnvVars(ctx); err != nil {
        panic(err)
    }
    
    // 2. 然后解析命令行参数
    args := os.Args[1:] // 排除程序名
    if err := parser.ParseCommand(ctx, args); err != nil {
        panic(err)
    }
    
    // 3. 使用解析结果
    // ctx 现在包含了所有解析的标志和参数值
}
```

### 环境变量支持示例

```go
// 假设有以下环境变量设置
// export MY_APP_DEBUG=true
// export MY_APP_PORT=8080

func setupFlags(ctx *types.CmdContext) {
    // 设置支持环境变量的标志
    // 解析器会自动从 MY_APP_DEBUG 环境变量读取值
    debugFlag := &types.Flag{
        LongName:  "debug",
        ShortName: "d",
        EnvVar:    "MY_APP_DEBUG",
        // ... 其他配置
    }
    
    portFlag := &types.Flag{
        LongName:  "port",
        ShortName: "p", 
        EnvVar:    "MY_APP_PORT",
        // ... 其他配置
    }
    
    // 添加到上下文
    ctx.Flags = append(ctx.Flags, debugFlag, portFlag)
}

func main() {
    ctx := &types.CmdContext{}
    setupFlags(ctx)
    
    // LoadEnvVars 会自动读取环境变量值
    parser.LoadEnvVars(ctx)
    
    // ParseCommand 会处理命令行参数，命令行参数优先级更高
    parser.ParseCommand(ctx, os.Args[1:])
}
```

### 子命令解析示例

```go
func main() {
    ctx := &types.CmdContext{
        SubCommands: []*types.Command{
            {
                LongName: "serve",
                ShortName: "s",
                // ... 子命令配置
            },
            {
                LongName: "build", 
                ShortName: "b",
                // ... 子命令配置
            },
        },
    }
    
    // 解析命令: myapp serve --port 8080
    args := []string{"serve", "--port", "8080"}
    
    if err := parser.ParseCommand(ctx, args); err != nil {
        panic(err)
    }
    
    // 检查解析结果
    if ctx.CurrentSubCommand != nil {
        fmt.Printf("执行子命令: %s\n", ctx.CurrentSubCommand.LongName)
    }
}
```

## 解析规则

### 标志格式支持

- **长选项**: `--flag`, `--flag=value`, `--flag value`
- **短选项**: `-f`, `-f=value`, `-f value`
- **组合短选项**: `-abc` (等同于 `-a -b -c`)

### 优先级规则

1. **命令行参数** - 最高优先级
2. **环境变量** - 中等优先级  
3. **默认值** - 最低优先级

### 特殊处理

- **布尔标志**: 支持 `--flag` 和 `--no-flag` 格式
- **数组标志**: 支持多次指定同一标志
- **分隔符**: 支持 `--` 分隔符，后续参数不作为标志处理

## 错误处理

解析器会返回以下类型的错误：

- **未知标志错误**: 遇到未定义的标志时
- **缺少参数错误**: 标志需要值但未提供时
- **类型转换错误**: 参数值无法转换为期望类型时
- **子命令错误**: 子命令不存在或配置错误时

## 设计特点

1. **纯函数设计** - 所有函数都是纯函数，不依赖全局状态
2. **上下文驱动** - 通过 `CmdContext` 传递所有必要信息
3. **分离关注点** - 环境变量加载和命令行解析分别处理
4. **错误友好** - 提供详细的错误信息和位置
5. **扩展性强** - 支持自定义标志类型和解析规则

## 性能考虑

- 解析过程为 O(n) 时间复杂度，其中 n 为参数数量
- 内存使用与标志数量和参数长度成正比
- 环境变量查找使用系统调用，建议在程序启动时一次性加载

## 注意事项

- 环境变量名通常使用大写字母和下划线
- 命令行参数会覆盖环境变量设置的值
- 子命令解析是递归的，支持多级子命令
- 解析失败时，上下文状态可能处于部分解析状态