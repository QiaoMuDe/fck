<div align="center">

# 🚀 QFlag

<a name="top"></a>

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8E6?style=flat&logo=go)](https://golang.org/) [![MIT License](https://img.shields.io/badge/License-MIT-green?style=flat)](https://opensource.org/licenses/MIT) [![Gitee](https://img.shields.io/badge/Gitee-qflag-red?style=flat)](https://gitee.com/MM-Q/qflag) [![GitHub](https://img.shields.io/badge/GitHub-qflag-black?style=flat)](https://github.com/QiaoMuDe/qflag) [![Ask DeepWiki](https://deepwiki.com/badge.svg?style=flat)](https://deepwiki.com/QiaoMuDe/qflag)

*泛型设计 • 自动路由 • 类型安全 • 并发安全 • 自动补全 • 智能纠错*

</div>

---

## 📖 项目简介

QFlag 是一个专为 Go 语言设计的命令行参数解析库, 提供了丰富的功能和优雅的 API, 帮助开发者快速构建专业的命令行工具。它支持多种标志类型、子命令、环境变量绑定、自动补全等高级特性, 同时保持简单易用的设计理念。

### ✨ 核心特性

| 特性分类 | 功能描述 |
|---------|---------|
| **类型系统** | 泛型设计、类型安全、17种标志类型（基础类型、枚举、时间、大小、集合） |
| **验证机制** | 标志验证、互斥组、必需组、条件性必需组、标志依赖关系 |
| **命令系统** | 子命令支持、命令/标志别名、嵌套命令、自动路由 |
| **配置环境** | 环境变量绑定、配置优先级（命令行 > 环境变量 > 默认值）、全局根命令 |
| **开发体验** | 智能纠错、Shell补全（Bash/Pwsh）、自动生成帮助文档、双语支持 |
| **工程特性** | 并发安全（读写锁保护）、零外部依赖、标准Go错误处理 |

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
```

### 高级用法

#### 智能纠错功能

QFlag 内置智能纠错功能，当用户输入错误的子命令或标志时，会自动推荐相似的选项。

```bash
# 子命令纠错示例
$ myapp cnfig
myapp: 'cnfig' is not a valid command. See 'myapp --help'.

The most similar commands are
        config

# 标志纠错示例
$ myapp config --verb
myapp: unknown flag: '--verb'

The most similar flags are
        --verbose
        -v
```

智能纠错功能无需额外配置，自动集成在解析流程中。

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
    qflag.Root.AddMutexGroup("format", []string{"json", "xml", "yaml"}, false)

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
    qflag.Root.AddRequiredGroup("connection", []string{"host", "port"}, false)
    
    // 添加条件性必需组 - 如果使用其中一个则必须同时使用
    qflag.Root.AddRequiredGroup("database", []string{"dbhost", "dbport"}, true)

    // 解析参数
    if err := qflag.Parse(); err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }

    // 使用参数
    fmt.Printf("连接到 %s:%d\n", hostFlag.Get(), portFlag.Get())
}
```

#### 条件性必需组

条件性必需组允许定义一组标志，如果组中任何一个标志被使用，则所有标志都必须被设置。

```go
// 创建条件性必需组的标志
dbHostFlag := qflag.Root.String("dbhost", "dh", "数据库主机", "")
dbPortFlag := qflag.Root.Uint("dbport", "dp", "数据库端口", 0)

// 添加条件性必需组 - 第三个参数为 true
qflag.Root.AddRequiredGroup("database", []string{"dbhost", "dbport"}, true)
```

#### 标志依赖关系

定义标志之间的依赖约束，支持互斥依赖（`DepMutex`）和必需依赖（`DepRequired`）。

```go
// 创建标志
remoteFlag := qflag.Root.Bool("remote", "r", "使用远程模式", false)
sslFlag := qflag.Root.Bool("ssl", "s", "启用SSL", false)
certFlag := qflag.Root.String("cert", "c", "证书路径", "")
keyFlag := qflag.Root.String("key", "k", "密钥路径", "")

// 互斥依赖：远程模式与本地路径互斥
qflag.Root.AddFlagDependency("remote_mutex_local", "remote", []string{"local-path"}, qflag.DepMutex)

// 必需依赖：SSL模式需要证书和密钥
qflag.Root.AddFlagDependency("ssl_requires_cert_key", "ssl", []string{"cert", "key"}, qflag.DepRequired)
```

#### 环境变量绑定

QFlag 支持通过环境变量设置标志值，优先级：命令行参数 > 环境变量 > 默认值。

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

    // 创建标志
    hostFlag := cmd.String("host", "H", "主机地址", "localhost")
    portFlag := cmd.Int("port", "p", "端口号", 8080)

    // 方式1: 手动指定环境变量名
    hostFlag.BindEnv("HOST_ADDR")  // 绑定到 MYAPP_HOST_ADDR
    
    // 方式2: 自动绑定（使用标志长名称大写）
    portFlag.AutoBindEnv()  // 绑定到 MYAPP_PORT
    
    // 方式3: 批量自动绑定所有标志
    // cmd.AutoBindAllEnv()
    
    if err := cmd.Parse(os.Args[1:]); err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }

    fmt.Printf("主机: %s, 端口: %d\n", hostFlag.GetStr(), portFlag.GetInt())
}
```

**三种绑定方式对比：**

| 方式 | 方法 | 适用场景 | 环境变量名 |
|------|------|----------|------------|
| 手动指定 | `BindEnv("NAME")` | 需要自定义名称 | `前缀_NAME` |
| 自动绑定 | `AutoBindEnv()` | 单个标志绑定 | `前缀_标志长名大写` |
| 批量绑定 | `AutoBindAllEnv()` | 批量绑定所有标志 | `前缀_标志长名大写` |

---

## 📚 API 文档概述

QFlag 提供了简洁而强大的 API, 主要包含以下核心组件: 

### 核心概念

- **Command (命令)** - 命令行工具的核心, 支持标志管理、参数解析、子命令等功能
- **Flag (标志)** - 命令行参数的抽象, 支持多种数据类型
- **MutexGroup (互斥组)** - 确保组内只有一个标志被设置
- **RequiredGroup (必需组)** - 确保组内所有标志都被设置, 支持普通必需组和条件性必需组两种模式
- **FlagDependency (标志依赖)** - 定义触发标志与目标标志之间的约束关系（互斥或必需）
- **智能纠错** - 输入错误的子命令或标志时，自动推荐相似的选项

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

完整的 API 文档和示例代码, 请参考项目中的 [_examples/](_examples/) 目录和源代码中的注释。

### 📚 更多文档

- [API 文档](APIDOC.md) - 完整的 API 参考
- [qflag-cli](qflag-cli) - CLI 开发技能包，包含快速开始指南、核心概念速查、完整开发规范和示例代码（适用于 AI 技能调用和人类学习） 

---

## ⚙️ 配置选项说明

### 错误处理策略

QFlag 提供了三种错误处理策略: 

- **ContinueOnError** - 遇到错误继续解析
- **ExitOnError** - 遇到错误立即退出
- **ReturnOnError** - 遇到错误返回错误

### CmdOpts 配置选项

`CmdOpts` 提供了配置现有命令的方式，包含命令的所有可配置属性。完整字段说明请参考 [APIDOC.md](APIDOC.md)。

#### 使用示例

```go
// 创建命令并添加标志
cmd := qflag.NewCmd("myapp", "m", qflag.ExitOnError)
cmd.String("host", "H", "主机地址", "localhost")
cmd.Int("port", "p", "端口号", 8080)

// 创建配置选项
opts := &qflag.CmdOpts{
    Desc:        "我的应用程序",
    Version:     "1.0.0",
    UseChinese:  true,
    EnvPrefix:   "MYAPP",
    AutoBindEnv: true,
}

// 应用配置
cmd.ApplyOpts(opts)
```

#### 特点

- **部分配置**: 未设置的属性不会被修改，保留原有值
- **批量设置**: 一次性设置多个命令属性
- **结构化管理**: 通过结构体集中管理配置

### 禁用标志解析

通过设置 `DisableFlagParsing: true` 可将所有参数（包括 `--flag` 形式）作为位置参数处理。

**使用场景:**
- 透传参数给外部命令（如 `kubectl exec`）
- SSH 包装器
- 需要保留原始参数的场景

**配置方式:**
```go
opts := &qflag.CmdOpts{
    DisableFlagParsing: true,
}
cmd.ApplyOpts(opts)
```

完整示例请参考 [_examples/cmdopts](_examples/cmdopts) 目录。

### 隐藏命令

设置 `Hidden: true` 后，命令不会显示在帮助信息中，但仍可通过命令行正常调用。

**使用场景:**
- 内部调试命令
- 已弃用但仍需兼容的命令
- 高级或实验性功能

**配置方式:**
```go
opts := &qflag.CmdOpts{
    Hidden: true,
}
cmd.ApplyOpts(opts)
// 或
cmd.SetHidden(true)
```

完整示例请参考 [_examples/cmdopts](_examples/cmdopts) 目录。

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