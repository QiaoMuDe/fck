<div align="center">

# 🚀 QFlag

<a name="top"></a>

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8E6?style=flat&logo=go)](https://golang.org/) [![MIT License](https://img.shields.io/badge/License-MIT-green?style=flat)](https://opensource.org/licenses/MIT) [![Gitee](https://img.shields.io/badge/Gitee-qflag-red?style=flat)](https://gitee.com/MM-Q/qflag) [![GitHub](https://img.shields.io/badge/GitHub-qflag-black?style=flat)](https://github.com/QiaoMuDe/qflag) [![Ask DeepWiki](https://deepwiki.com/badge.svg?style=flat)](https://deepwiki.com/QiaoMuDe/qflag)

*泛型设计 • 自动路由 • 类型安全 • 并发安全 • 自动补全 • 子命令管理*

</div>

---

## 📖 项目简介

QFlag 是一个专为 Go 语言设计的命令行参数解析库, 提供了丰富的功能和优雅的 API, 帮助开发者快速构建专业的命令行工具。它支持多种标志类型、子命令、环境变量绑定、自动补全等高级特性, 同时保持简单易用的设计理念。

### ✨ 核心特性

#### 🎯 核心功能
- **类型安全** - 支持多种标志类型, 确保类型安全
- **高性能**   - 优化的解析算法, 快速高效
- **泛型设计** - 基于Go泛型的现代化设计, 类型安全且代码简洁
- **并发安全** - 使用读写锁保护, 支持并发访问标志和命令
- **自动路由** - 智能路由子命令, 自动匹配并执行

#### 🛡️ 错误和验证
- **错误处理** - 使用标准 Go 错误处理, 简洁直观
- **标志验证** - 丰富的内置验证器和自定义验证器支持
- **互斥标志** - 支持标志互斥组, 确保组内只有一个标志被设置
- **必需标志** - 支持标志必需组和条件性必需组, 灵活控制标志依赖关系

#### 🌍 国际化和帮助
- **国际化** - 支持中文和英文双语
- **帮助生成** - 自动生成专业的帮助文档
- **版本信息** - 内置版本信息显示功能
- **示例展示** - 支持在帮助中显示使用示例
- **注意事项** - 支持在帮助中显示注意事项

#### 🔄 配置和环境
- **环境变量** - 自动绑定环境变量, 支持环境变量前缀
- **配置选项** - 灵活的命令配置选项, 支持Logo、使用语法等自定义
- **配置优先级** - 命令行参数 > 环境变量 > 默认值
- **全局根命令** - 提供全局根命令实例, 简化使用方式

#### 📝 补全和文档
- **自动补全** - 生成 Bash 和 PowerShell 补全脚本
- **命令别名** - 支持命令长名和短名, 提供灵活的调用方式
- **标志别名** - 支持标志长名和短名, 提供灵活的调用方式

#### 🌳 命令和标志
- **子命令** - 完整的子命令支持, 支持嵌套子命令
- **便捷方法** - 简洁的API设计, 减少样板代码
- **命令选项** - 支持通过选项配置现有命令

#### 📊 标志类型
- **基础类型** - String, Bool, Int, Int64, Uint, Uint8, Uint16, Uint32, Uint64, Float64
- **特殊类型** - Enum (枚举类型)
- **时间类型** - Duration (持续时间), Time (时间点)
- **大小类型** - Size (支持KB、MB、GB等单位)
- **集合类型** - Map (键值对), StringSlice, IntSlice, Int64Slice

---

## 🔗 项目地址

该项目同时托管在 Gitee 和 GitHub 上，您可以选择合适的平台访问：

| 平台               | 地址                                                        | 描述                       |
| ------------------ | ----------------------------------------------------------- | -------------------------- |
| 🔴**Gitee**  | [gitee.com/MM-Q/qflag](https://gitee.com/MM-Q/qflag)           | 国内访问更快，主要开发仓库 |
| ⚫**GitHub** | [github.com/QiaoMuDe/qflag](https://github.com/QiaoMuDe/qflag) | 国际化平台，同步更新       |

---

## 📦 安装指南

使用 `go get` 命令安装：

```bash
go get -u gitee.com/MM-Q/qflag
```

然后在代码中导入：

```go
import "gitee.com/MM-Q/qflag"
```

---

## 💡 使用示例

### 🚀 全局根命令 (推荐) 

QFlag 提供了全局根命令 `qflag.Root`, 这是最简单、最直接的使用方式。**推荐优先使用**全局根命令作为命令行工具的入口点。

```go
package main

import (
    "fmt"
    "gitee.com/MM-Q/qflag"
)

func main() {
    // 直接使用全局根命令创建标志
    name := qflag.Root.String("name", "n", "用户名", "guest")
    age := qflag.Root.Int("age", "a", "年龄", 18)
    verbose := qflag.Root.Bool("verbose", "v", "详细模式", false)
    
    // 配置全局命令
    qflag.Root.SetDesc("示例应用程序")
    qflag.Root.SetVersion("1.0.0")
    
    // 解析命令行参数
    if err := qflag.Parse(); err != nil {
        fmt.Printf("解析错误: %v\n", err)
        return
    }
    
    // 使用参数
    fmt.Printf("用户名: %s\n", name.Get())
    fmt.Printf("年龄: %d\n", age.Get())
    fmt.Printf("详细模式: %t\n", verbose.Get())
    fmt.Printf("非标志参数: %v\n", qflag.Root.Args())
}
```

#### 全局根命令的优势

- 🎯 **简单易用**: 无需手动创建命令实例, 直接使用
- 🚀 **零配置**: 自动使用可执行文件名作为命令名
- 🔧 **功能完整**: 支持所有 QFlag 的高级功能
- 📦 **统一入口**: 所有操作都通过 `qflag.Root` 访问

#### 全局根命令支持的便捷函数

```go
// 解析函数
qflag.Parse()          // 解析命令行参数
qflag.ParseOnly()       // 仅解析当前命令
qflag.ParseAndRoute()   // 解析并路由到子命令

// 子命令管理
qflag.AddSubCmds(cmd1, cmd2)           // 添加子命令
qflag.AddSubCmdFrom([]Command{cmd1, cmd2}) // 从切片添加子命令

// 互斥组
qflag.AddMutexGroup("format", []string{"json", "xml"}, false)

// 必需组
qflag.AddRequiredGroup("connection", []string{"host", "port"}, false)  // 普通必需组
qflag.AddRequiredGroup("database", []string{"dbhost", "dbport"}, true) // 条件性必需组
```

### 基础用法

```go
package main

import (
    "fmt"
    "os"
    "gitee.com/MM-Q/qflag"
)

func main() {
    // 创建命令
    cmd := qflag.NewCmd("myapp", "m", qflag.ContinueOnError)
    cmd.SetDesc("我的应用程序")
    cmd.SetVersion("1.0.0")

    // 使用便捷方法创建标志
    nameFlag := cmd.String("name", "n", "用户名", "guest")
    verboseFlag := cmd.Bool("verbose", "v", "详细输出", false)

    // 解析参数
    if err := cmd.Parse(os.Args[1:]); err != nil {
        fmt.Printf("参数解析错误: %v\n", err)
        os.Exit(1)
    }

    // 使用参数
    fmt.Printf("用户名: %s\n", nameFlag.GetStr())
    fmt.Printf("详细模式: %v\n", verboseFlag.IsSet())
}
```

### 高级用法

#### 使用全局根命令的子命令支持

```go
package main

import (
    "fmt"
    "gitee.com/MM-Q/qflag"
)

func main() {
    // 配置全局根命令
    qflag.Root.SetDesc("命令行工具")
    qflag.Root.SetVersion("1.0.0")
    
    // 创建全局标志
    verbose := qflag.Root.Bool("verbose", "v", "详细输出", false)
    
    // 创建子命令
    listCmd := qflag.NewCmd("list", "ls", qflag.ContinueOnError)
    listCmd.SetDesc("列出所有项目")
    listCmd.Bool("all", "a", "显示所有项目", false)
    
    addCmd := qflag.NewCmd("add", "a", qflag.ContinueOnError)
    addCmd.SetDesc("添加新项目")
    addCmd.String("name", "n", "项目名称", "")
    
    // 添加子命令到全局根命令
    qflag.AddSubCmds(listCmd, addCmd)
    
    // 解析并路由
    if err := qflag.ParseAndRoute(); err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }
    
    // 如果是根命令执行
    if qflag.Root.NArg() == 0 {
        fmt.Printf("详细模式: %t\n", verbose.Get())
        fmt.Println("使用 'help' 查看可用命令")
    }
}
```

#### 传统子命令支持

```go
package main

import (
    "fmt"
    "os"
    "gitee.com/MM-Q/qflag"
)

func main() {
    // 创建根命令
    rootCmd := qflag.NewCmd("cli", "", qflag.ContinueOnError)
    rootCmd.SetDesc("命令行工具")

    // 使用便捷方法创建子命令
    listCmd := qflag.NewCmd("list", "ls", qflag.ContinueOnError)
    listCmd.SetDesc("列出所有项目")
    listCmd.Bool("all", "a", "显示所有项目", false)
    
    addCmd := qflag.NewCmd("add", "a", qflag.ContinueOnError)
    addCmd.SetDesc("添加新项目")
    addCmd.String("name", "n", "项目名称", "")
    
    // 添加子命令
    rootCmd.AddSubCmds(listCmd, addCmd)

    // 解析并路由
    if err := rootCmd.ParseAndRoute(os.Args[1:]); err != nil {
        fmt.Printf("错误: %v\n", err)
        os.Exit(1)
    }
}
```

#### 便捷方法创建标志

```go
package main

import (
    "fmt"
    "os"
    "time"
    "gitee.com/MM-Q/qflag"
)

func main() {
    cmd := qflag.NewCmd("server", "s", qflag.ContinueOnError)

    // 使用便捷方法创建多个标志
    cmd.String("host", "h", "主机地址", "localhost")
    cmd.Uint("port", "p", "端口号", 8080)
    cmd.Duration("timeout", "t", "超时时间", time.Second*30)
    cmd.Bool("debug", "d", "调试模式", false)

    fmt.Printf("成功添加 %d 个标志\n", len(cmd.Flags()))
}
```

#### 使用全局根命令的互斥标志组

```go
package main

import (
    "fmt"
    "gitee.com/MM-Q/qflag"
)

func main() {
    // 配置全局根命令
    qflag.Root.SetDesc("格式转换工具")
    
    // 使用全局根命令创建互斥标志
    jsonFlag := qflag.Root.Bool("json", "j", "JSON 格式", false)
    xmlFlag := qflag.Root.Bool("xml", "x", "XML 格式", false)
    yamlFlag := qflag.Root.Bool("yaml", "y", "YAML 格式", false)

    // 添加互斥组到全局根命令
    qflag.AddMutexGroup("format", []string{"json", "xml", "yaml"}, false)

    // 解析参数
    if err := qflag.Parse(); err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }

    // 使用参数
    if jsonFlag.Get() {
        fmt.Println("使用 JSON 格式")
    } else if xmlFlag.Get() {
        fmt.Println("使用 XML 格式")
    } else if yamlFlag.Get() {
        fmt.Println("使用 YAML 格式")
    }
}
```

#### 传统互斥标志组

```go
package main

import (
    "fmt"
    "os"
    "gitee.com/MM-Q/qflag"
)

func main() {
    cmd := qflag.NewCmd("converter", "", qflag.ContinueOnError)

    // 使用便捷方法创建互斥标志
    jsonFlag := cmd.Bool("json", "j", "JSON 格式", false)
    xmlFlag := cmd.Bool("xml", "x", "XML 格式", false)
    yamlFlag := cmd.Bool("yaml", "y", "YAML 格式", false)

    // 创建互斥标志组
    formatGroup := qflag.NewMutexGroup("format", "输出格式", true)
    formatGroup.AddFlags(jsonFlag, xmlFlag, yamlFlag)
    cmd.AddMutexGroup(formatGroup)

    // 解析参数
    if err := cmd.Parse(os.Args[1:]); err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }
}
```

#### 使用全局根命令的必需标志组

```go
package main

import (
    "fmt"
    "gitee.com/MM-Q/qflag"
)

func main() {
    // 配置全局根命令
    qflag.Root.SetDesc("数据库连接工具")
    
    // 使用全局根命令创建必需标志
    hostFlag := qflag.Root.String("host", "h", "主机地址", "")
    portFlag := qflag.Root.Uint("port", "p", "端口号", 0)

    // 添加普通必需组到全局根命令 - 所有标志都必须设置
    qflag.AddRequiredGroup("connection", []string{"host", "port"}, false)
    
    // 添加条件性必需组 - 如果使用其中一个则必须同时使用
    qflag.AddRequiredGroup("database", []string{"dbhost", "dbport"}, true)

    // 解析参数
    if err := qflag.Parse(); err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }

    // 使用参数
    fmt.Printf("连接到 %s:%d\n", hostFlag.Get(), portFlag.Get())
}
```

#### 传统必需标志组

```go
package main

import (
    "fmt"
    "os"
    "gitee.com/MM-Q/qflag"
)

func main() {
    cmd := qflag.NewCmd("db-connect", "", qflag.ContinueOnError)

    // 使用便捷方法创建必需标志
    hostFlag := cmd.String("host", "h", "主机地址", "")
    portFlag := cmd.Uint("port", "p", "端口号", 0)
    usernameFlag := cmd.String("username", "u", "用户名", "")
    passwordFlag := cmd.String("password", "P", "密码", "")

    // 添加普通必需组 - 所有标志都必须设置
    cmd.AddRequiredGroup("connection", []string{"host", "port"}, false)
    cmd.AddRequiredGroup("auth", []string{"username", "password"}, false)
    
    // 添加条件性必需组 - 如果使用其中一个则必须同时使用
    cmd.AddRequiredGroup("optional-features", []string{"feature-a", "feature-b"}, true)

    // 解析参数
    if err := cmd.Parse(os.Args[1:]); err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }

    // 使用参数
    fmt.Printf("连接到 %s:%d\n", hostFlag.Get(), portFlag.Get())
    fmt.Printf("用户: %s\n", usernameFlag.Get())
}
```

#### 条件性必需组

条件性必需组是QFlag提供的一种灵活的标志验证方式，它允许您定义一组标志，如果组中任何一个标志被使用，则所有标志都必须被设置。这与普通必需组不同，普通必需组要求所有标志都必须被设置，无论是否使用了其中任何一个。

```go
package main

import (
    "fmt"
    "gitee.com/MM-Q/qflag"
)

func main() {
    // 配置全局根命令
    qflag.Root.SetDesc("数据库连接工具")
    
    // 创建条件性必需组的标志
    dbHostFlag := qflag.Root.String("dbhost", "dh", "数据库主机", "")
    dbPortFlag := qflag.Root.Uint("dbport", "dp", "数据库端口", 0)
    dbUserFlag := qflag.Root.String("dbuser", "du", "数据库用户", "")
    dbPassFlag := qflag.Root.String("dbpass", "dp", "数据库密码", "")

    // 添加条件性必需组 - 如果使用其中一个则必须同时使用
    qflag.AddRequiredGroup("database", []string{"dbhost", "dbport", "dbuser", "dbpass"}, true)

    // 解析参数
    if err := qflag.Parse(); err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }

    // 使用参数
    if dbHostFlag.IsSet() {
        fmt.Printf("数据库连接: %s:%d, 用户: %s\n", 
            dbHostFlag.Get(), dbPortFlag.Get(), dbUserFlag.Get())
    } else {
        fmt.Println("未使用数据库连接")
    }
}
```

**条件性必需组的使用场景:**

1. **可选功能模块**: 当应用程序有多个可选功能模块，每个模块需要一组相关的配置参数时
2. **可选连接**: 当应用程序可以连接到多个可选的服务，每个服务需要自己的连接参数时
3. **可选输出格式**: 当应用程序支持多种输出格式，每种格式需要特定的参数时

**普通必需组 vs 条件性必需组:**

- **普通必需组**: 组中的所有标志都必须被设置，无论是否使用了其中任何一个
- **条件性必需组**: 只有当组中任何一个标志被设置时，所有标志才必须被设置

#### 环境变量绑定

QFlag 提供了三种环境变量绑定方式，满足不同场景需求。

##### 方式一：手动指定环境变量名

通过 `BindEnv()` 方法手动指定环境变量名称：

```go
package main

import (
    "fmt"
    "os"
    "gitee.com/MM-Q/qflag"
)

func main() {
    cmd := qflag.NewCmd("app", "", qflag.ContinueOnError)
    
    // 设置环境变量前缀
    cmd.SetEnvPrefix("MYAPP")

    // 手动指定环境变量名绑定
    dbFlag := cmd.String("database", "d", "数据库地址", "localhost")
    dbFlag.BindEnv("DATABASE_URL")  // 绑定到 MYAPP_DATABASE_URL
    
    if err := cmd.Parse(os.Args[1:]); err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }

    fmt.Printf("数据库地址: %s\n", dbFlag.GetStr())
}
```

##### 方式二：标志自动绑定

通过 `AutoBindEnv()` 方法自动使用标志长名称的大写形式作为环境变量名：

```go
package main

import (
    "fmt"
    "os"
    "gitee.com/MM-Q/qflag"
)

func main() {
    cmd := qflag.NewCmd("app", "", qflag.ContinueOnError)
    cmd.SetEnvPrefix("MYAPP")

    // 创建标志
    hostFlag := cmd.String("host", "H", "主机地址", "localhost")
    portFlag := cmd.Int("port", "p", "端口号", 8080)

    // 自动绑定：使用长名称大写作为环境变量名
    // host -> MYAPP_HOST, port -> MYAPP_PORT
    hostFlag.AutoBindEnv()
    portFlag.AutoBindEnv()
    
    if err := cmd.Parse(os.Args[1:]); err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }

    fmt.Printf("主机: %s, 端口: %d\n", hostFlag.GetStr(), portFlag.GetInt())
}
```

##### 方式三：命令批量自动绑定

通过 `AutoBindAllEnv()` 方法一次性为命令的所有标志自动绑定环境变量：

```go
package main

import (
    "fmt"
    "os"
    "gitee.com/MM-Q/qflag"
)

func main() {
    cmd := qflag.NewCmd("app", "", qflag.ContinueOnError)
    cmd.SetEnvPrefix("MYAPP")

    // 创建多个标志
    cmd.String("host", "H", "主机地址", "localhost")
    cmd.Int("port", "p", "端口号", 8080)
    cmd.String("user", "u", "用户名", "admin")
    cmd.String("password", "P", "密码", "")

    // 批量自动绑定所有标志的环境变量
    // host -> MYAPP_HOST
    // port -> MYAPP_PORT
    // user -> MYAPP_USER
    // password -> MYAPP_PASSWORD
    cmd.AutoBindAllEnv()
    
    if err := cmd.Parse(os.Args[1:]); err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }

    fmt.Printf("主机: %s, 端口: %d\n", cmd.GetString("host"), cmd.GetInt("port"))
}
```

##### 三种方式对比

| 方式 | 方法 | 适用场景 | 特点 |
|------|------|----------|------|
| 手动指定 | `BindEnv("NAME")` | 需要自定义环境变量名 | 灵活，可指定任意名称 |
| 标志自动绑定 | `AutoBindEnv()` | 单个标志自动绑定 | 使用长名称大写，简洁 |
| 命令批量绑定 | `AutoBindAllEnv()` | 批量绑定所有标志 | 一次性绑定，高效 |

---

## 📚 API 文档概述

QFlag 提供了简洁而强大的 API, 主要包含以下核心组件: 

### 🚀 全局根命令 (推荐使用方式) 

QFlag 提供了全局根命令 `qflag.Root`, 这是最简单、最直接的使用方式。**推荐优先使用**全局根命令作为命令行工具的入口点。

#### 全局根命令的优势

- 🎯 **简单易用**: 无需手动创建命令实例, 直接使用
- 🚀 **零配置**: 自动使用可执行文件名作为命令名
- 🔧 **功能完整**: 支持所有 QFlag 的高级功能
- 📦 **统一入口**: 所有操作都通过 `qflag.Root` 访问

#### 全局根命令支持的便捷函数

```go
// 解析函数
qflag.Parse()          // 解析命令行参数
qflag.ParseOnly()       // 仅解析当前命令
qflag.ParseAndRoute()   // 解析并路由到子命令

// 子命令管理
qflag.AddSubCmds(cmd1, cmd2)           // 添加子命令
qflag.AddSubCmdFrom([]Command{cmd1, cmd2}) // 从切片添加子命令

// 互斥组
qflag.AddMutexGroup("format", []string{"json", "xml"}, false)
```

#### 全局根命令的使用方式

```go
// 直接使用全局根命令创建标志
name := qflag.Root.String("name", "n", "用户名", "guest")
age := qflag.Root.Int("age", "a", "年龄", 18)
verbose := qflag.Root.Bool("verbose", "v", "详细模式", false)

// 配置全局命令
qflag.Root.SetDesc("示例应用程序")
qflag.Root.SetVersion("1.0.0")

// 解析命令行参数
if err := qflag.Parse(); err != nil {
    fmt.Printf("解析错误: %v\n", err)
    return
}
```

### 核心概念

- **Command (命令)** - 命令行工具的核心, 支持标志管理、参数解析、子命令等功能
- **Flag (标志)** - 命令行参数的抽象, 支持多种数据类型
- **MutexGroup (互斥组)** - 确保组内只有一个标志被设置
- **RequiredGroup (必需组)** - 确保组内所有标志都被设置, 支持普通必需组和条件性必需组两种模式

### 便捷方法

Command 类型提供了丰富的便捷方法来创建各种类型的标志, 无需手动创建和添加标志: 

```go
// 字符串和布尔标志
nameFlag := cmd.String("name", "n", "用户名", "guest")
verboseFlag := cmd.Bool("verbose", "v", "详细输出", false)

// 数值类型标志
portFlag := cmd.Uint("port", "p", "端口号", 8080)
timeoutFlag := cmd.Duration("timeout", "t", "超时时间", time.Second*30)

// 集合类型标志
filesFlag := cmd.StringSlice("files", "f", "文件列表", []string{})
tagsFlag := cmd.IntSlice("tags", "", "标签列表", []int{})
```

### 详细的 API 文档

完整的 API 文档和示例代码, 请参考项目中的 `examples/` 目录和源代码中的注释。

### 📚 标志使用语法

详细的标志使用语法和高级功能说明, 请参考 [FLAG_USAGE.md](FLAG_USAGE.md)。

---

## 🎯 支持的功能

### 标志功能

- ✅ 短标志名 (单字符, 如 `-v`) 
- ✅ 长标志名 (多字符, 如 `--verbose`) 
- ✅ 默认值设置
- ✅ 必需标志验证
- ✅ 条件性必需组验证
- ✅ 标志描述
- ✅ 环境变量绑定
- ✅ 标志值验证
- ✅ 枚举值限制
- ✅ 切片类型支持

### 命令功能

- ✅ 子命令支持
- ✅ 命令别名
- ✅ 命令描述
- ✅ 版本信息
- ✅ 帮助信息
- ✅ 示例展示
- ✅ 注意事项

### 高级功能

- ✅ 互斥标志组
- ✅ 必需标志组
- ✅ 条件性必需组
- ✅ 自动补全脚本生成
- ✅ 环境变量前缀
- ✅ 错误处理策略
- ✅ 中文/英文双语支持
- ✅ Logo 文本设置
- ✅ 使用语法自定义
- ✅ 禁用标志解析（DisableFlagParsing）
- ✅ 隐藏命令（Hidden）

---

## ⚙️ 配置选项说明

### 错误处理策略

QFlag 提供了三种错误处理策略: 

- **ContinueOnError** - 遇到错误继续解析
- **ExitOnError** - 遇到错误立即退出
- **ReturnOnError** - 遇到错误返回错误

### CmdOpts 配置选项

`CmdOpts` 提供了配置现有命令的方式，包含命令的所有可配置属性。

#### 字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| `Desc` | `string` | 命令描述 |
| `RunFunc` | `func(Command) error` | 命令执行函数 |
| `Version` | `string` | 版本号 |
| `UseChinese` | `bool` | 是否使用中文 |
| `EnvPrefix` | `string` | 环境变量前缀 |
| `UsageSyntax` | `string` | 命令使用语法 |
| `LogoText` | `string` | Logo文本 |
| `Completion` | `bool` | 是否启用自动补全标志 |
| `DisableFlagParsing` | `bool` | 是否禁用标志解析，所有参数作为位置参数处理 |
| `Hidden` | `bool` | 是否隐藏命令，不在帮助信息中显示 |
| `AutoBindEnv` | `bool` | 是否自动绑定所有标志的环境变量 |
| `Examples` | `map[string]string` | 示例使用 |
| `Notes` | `[]string` | 注意事项 |
| `SubCmds` | `[]Command` | 子命令列表 |
| `MutexGroups` | `[]MutexGroup` | 互斥组列表 |
| `RequiredGroups` | `[]RequiredGroup` | 必需组列表 |

#### 使用示例

```go
package main

import (
    "fmt"
    "os"
    "gitee.com/MM-Q/qflag"
)

func main() {
    // 创建命令
    cmd := qflag.NewCmd("myapp", "m", qflag.ExitOnError)
    
    // 添加标志
    cmd.String("host", "H", "主机地址", "localhost")
    cmd.Int("port", "p", "端口号", 8080)
    
    // 创建配置选项并应用
    opts := &qflag.CmdOpts{
        Desc:        "我的应用程序",
        Version:     "1.0.0",
        UseChinese:  true,
        EnvPrefix:   "MYAPP",
        AutoBindEnv: true,
        Examples: map[string]string{
            "启动服务":    "myapp --host 0.0.0.0 --port 8080",
            "使用环境变量": "MYAPP_HOST=0.0.0.0 MYAPP_PORT=8080 myapp",
        },
        Notes: []string{
            "端口号必须在 1-65535 范围内",
            "环境变量会自动加上 MYAPP_ 前缀",
        },
    }
    
    // 应用配置
    if err := cmd.ApplyOpts(opts); err != nil {
        fmt.Printf("应用配置失败: %v\n", err)
        return
    }
    
    // 解析参数
    if err := cmd.Parse(os.Args[1:]); err != nil {
        fmt.Printf("解析失败: %v\n", err)
        return
    }
    
    fmt.Printf("服务启动在 %s:%d\n", cmd.GetString("host"), cmd.GetInt("port"))
}
```

#### 特点

- **部分配置**: 未设置的属性不会被修改，保留原有值
- **批量设置**: 一次性设置多个命令属性
- **结构化管理**: 通过结构体集中管理配置

### 禁用标志解析

当需要将所有参数（包括 `--flag` 形式）作为位置参数处理时，可以禁用标志解析。

```go
package main

import (
    "fmt"
    "gitee.com/MM-Q/qflag"
)

func main() {
    // 创建命令并禁用标志解析
    cmd := qflag.NewCmd("exec", "e", qflag.ExitOnError)
    
    // 方式1: 通过 CmdOpts 配置
    opts := &qflag.CmdOpts{
        Desc:               "执行外部命令",
        DisableFlagParsing: true,  // 禁用标志解析
    }
    
    if err := cmd.ApplyOpts(opts); err != nil {
        fmt.Printf("应用配置失败: %v\n", err)
        return
    }
    
    cmd.SetRun(func(c qflag.Command) error {
        // 所有参数都作为位置参数
        args := c.Args()
        fmt.Printf("执行命令: %v\n", args)
        return nil
    })
    
    // 使用: exec -- kubectl get pods -n default
    // 所有参数（包括 --flag）都会作为位置参数
    if err := cmd.ParseAndRoute(); err != nil {
        fmt.Printf("错误: %v\n", err)
    }
}
```

**使用场景:**
- 透传参数给外部命令（如 `kubectl exec`）
- SSH 包装器
- 任何需要保留原始参数的场景

### 隐藏命令

隐藏命令不会显示在帮助信息的子命令列表中，但仍可通过命令行正常调用。

```go
package main

import (
    "fmt"
    "gitee.com/MM-Q/qflag"
)

func main() {
    // 创建隐藏命令
    debugCmd := qflag.NewCmd("debug", "d", qflag.ExitOnError)
    
    // 方式1: 通过 CmdOpts 配置
    opts := &qflag.CmdOpts{
        Desc:   "调试命令（内部使用）",
        Hidden: true,  // 隐藏命令
    }
    
    if err := debugCmd.ApplyOpts(opts); err != nil {
        fmt.Printf("应用配置失败: %v\n", err)
        return
    }
    
    debugCmd.SetRun(func(c qflag.Command) error {
        fmt.Println("调试模式已启用")
        // 执行调试逻辑
        return nil
    })
    
    // 方式2: 通过方法设置
    // debugCmd.SetHidden(true)
    
    // 添加到根命令
    qflag.Root.AddSubCmds(debugCmd)
    
    // 调试命令可以通过命令行调用，但不会显示在帮助中
    if err := qflag.ParseAndRoute(); err != nil {
        fmt.Printf("错误: %v\n", err)
    }
}
```

**使用场景:**
- 内部调试命令
- 已弃用但仍需兼容的命令
- 高级或实验性功能

---

## 📁 项目结构

```
qflag/
├── internal/              # 内部实现
│   ├── builtin/          # 内置标志
│   ├── cmd/              # 命令实现
│   ├── completion/       # 自动补全脚本生成
│   │   └── templates/    # 补全脚本模板
│   ├── flag/             # 标志类型实现
│   ├── help/             # 帮助文档生成
│   ├── mock/             # Mock 测试工具
│   ├── parser/           # 参数解析器
│   ├── registry/         # 注册表实现
│   ├── types/            # 类型定义
│   └── utils/            # 工具函数
├── validators/           # 验证器
├── examples/             # 使用示例
│   ├── builtin-flags/    # 内置标签示例
│   ├── cmdopts/          # 命令选项示例
│   ├── error-formatting/ # 错误格式化示例
│   ├── flag-constructors/# 标志构造器示例
│   ├── mutex-group/      # 互斥组示例
│   ├── nested-commands/  # 嵌套命令示例
│   └── required-groups/  # 必需组示例
├── docs/                 # 设计文档
├── qflag-cli/            # AI 技能包
│   └── references/       # 参考文档
├── exports.go            # 公共 API 导出
├── qflag.go              # 全局根命令和便捷函数
├── APIDOC.md             # API 文档
├── FLAG_USAGE.md         # 标志用法文档
├── qflag命令开发规范.md   # 命令开发规范
├── qflag命令行工具开发规范.md # 命令行工具开发规范
├── go.mod                # Go 模块文件
└── README.md             # 项目文档
```

---

## 🧪 测试说明

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/cmd

# 运行测试并显示覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 测试覆盖

项目使用全面的测试套件, 包括: 
- 单元测试
- 集成测试
- 边界条件测试
- 错误处理测试

---

## 📄 许可证和贡献指南

### 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

### 贡献指南

我们欢迎任何形式的贡献！

#### 贡献方式

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

#### 代码规范

- 遵循 Go 语言代码规范
- 添加适当的注释
- 编写测试用例
- 确保所有测试通过
- 更新相关文档

#### 问题反馈

如果您发现 bug 或有功能建议, 请: 
- 搜索现有的 issues
- 创建新的 issue, 详细描述问题
- 提供复现步骤和环境信息

---

## 📞 联系方式和相关链接

### 相关资源

- 📦 **仓库地址**: [https://gitee.com/MM-Q/qflag.git](https://gitee.com/MM-Q/qflag.git)
- 📖 **文档**: [项目文档](https://gitee.com/MM-Q/qflag)
- 🐛 **问题反馈**: [Issues](https://gitee.com/MM-Q/qflag/issues)
- 💬 **讨论**: [Discussions](https://gitee.com/MM-Q/qflag/discussions)

### 联系方式

- 📧 **邮箱**: [提交 Issue](https://gitee.com/MM-Q/qflag/issues)

---

## 🙏 致谢

感谢所有为本项目做出贡献的开发者！

---

<div align="center">

**[⬆ 返回顶部](#top)**

Made with ❤️ by QFlag Team

</div>