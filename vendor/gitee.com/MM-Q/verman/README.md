<div align="center">

# verman 🚀

[![Go Version](https://img.shields.io/badge/Go-%3E%3D1.18-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/gitee.com/MM-Q/verman)](https://goreportcard.com/report/gitee.com/MM-Q/verman)
[![Release](https://img.shields.io/badge/release-v0.0.17-brightgreen)](https://gitee.com/MM-Q/verman/releases)

</div>

> 🎯 **轻量级 Go 版本信息管理库** - 专为简化应用程序版本管理而设计的现代化解决方案

verman 是一个功能强大且易于使用的 Go 语言版本信息管理库，支持编译时注入版本信息，提供 **12 种不同格式** 的版本输出，帮助开发者轻松管理和展示应用版本信息。

---

## ✨ 核心特性

- 🔧 **编译时注入** - 支持通过 `-ldflags` 在编译时注入版本信息
- 📊 **丰富输出格式** - 提供 **12 种不同** 的版本信息输出格式
- 🌐 **运行时信息** - 自动获取 Go 版本和平台信息
- 🎨 **简洁 API** - 通过全局实例 `V` 调用，无需构造函数
- 🚀 **零依赖** - 仅使用 Go 标准库，轻量级设计
- 🔒 **数据安全** - 私有变量设计，防止运行时意外修改
- ⚡ **高性能** - 优化的字符串格式化，支持并发访问
- 📋 **多行格式** - 支持横幅、表格、JSON 等多行显示格式

## 📦 安装指南

### 从 Gitee 仓库安装

```bash
go get gitee.com/MM-Q/verman
```

### 验证安装

```bash
go mod tidy
```

## 🚀 使用示例

### 基础用法

```go
package main

import (
    "fmt"
    "gitee.com/MM-Q/verman"
)

func main() {
    // 使用全局实例 V 调用方法
    fmt.Println("版本信息:", verman.V.Version())
    fmt.Println("短格式:", verman.V.Short())
    fmt.Println("完整信息:", verman.V.Full())
    
    // 多行格式展示
    fmt.Println("横幅格式:")
    fmt.Println(verman.V.Banner())
    
    fmt.Println("表格格式:")
    fmt.Println(verman.V.Table())
}
```

### 高级用法 - 编译时注入

```bash
# 使用 ldflags 注入版本信息（注意变量名为小写）
go build -ldflags "
-X 'gitee.com/MM-Q/verman.appName=myapp' 
-X 'gitee.com/MM-Q/verman.gitVersion=v1.2.3' 
-X 'gitee.com/MM-Q/verman.gitCommit=abc1234' 
-X 'gitee.com/MM-Q/verman.gitTreeState=clean' 
-X 'gitee.com/MM-Q/verman.gitCommitTime=2024-01-01T12:00:00Z' 
-X 'gitee.com/MM-Q/verman.buildTime=2024-01-01T12:30:00Z'
" main.go
```

### CLI 应用示例

```go
package main

import (
    "flag"
    "fmt"
    "os"
    "gitee.com/MM-Q/verman"
)

func main() {
    var showVersion = flag.Bool("version", false, "显示版本信息")
    var format = flag.String("format", "simple", "版本信息格式")
    flag.Parse()

    if *showVersion {
        switch *format {
        case "short":
            fmt.Println(verman.V.Short())
        case "long":
            fmt.Println(verman.V.Long())
        case "standard":
            fmt.Println(verman.V.Standard())
        case "full":
            fmt.Println(verman.V.Full())
        case "lines":
            fmt.Println(verman.V.Lines())
        case "table":
            fmt.Println(verman.V.Table())
        case "json":
            fmt.Println(verman.V.JSON())
        case "kv":
            fmt.Println(verman.V.KV())
        case "banner":
            fmt.Println(verman.V.Banner())
        case "csv":
            fmt.Println(verman.V.CSV())
        case "uri":
            fmt.Println(verman.V.URI())
        default:
            fmt.Println(verman.V.Version())
        }
        os.Exit(0)
    }

    // 应用程序主逻辑...
}
```

## 📚 API 文档概述

### Info 结构体

```go
type Info struct {
    AppName       string // 应用程序名称
    GitVersion    string // Git 语义化版本号
    GitCommit     string // Git 提交哈希值
    GitTreeState  string // Git 仓库状态
    GitCommitTime string // Git 提交时间
    BuildTime     string // 构建时间
    GoVersion     string // Go 运行时版本
    Platform      string // 平台信息
}
```

### 全局实例

| 变量名 | 类型 | 描述 |
|--------|------|------|
| `V` | `*Info` | 全局版本信息实例，供外部使用 |

### 结构体方法

| 方法名 | 返回值 | 描述 |
|--------|--------|------|
| `Version()` | `string` | 标准版本格式 |
| `Short()` | `string` | 短格式 |
| `Long()` | `string` | 长格式 |
| `Standard()` | `string` | 标准格式 |
| `Full()` | `string` | 完整版本格式 |
| `Lines()` | `string` | 分行格式 |
| `Table()` | `string` | 表格格式（多行） |
| `JSON()` | `string` | JSON格式 |
| `KV()` | `string` | 键值对格式 |
| `Banner()` | `string` | 横幅格式（多行） |
| `CSV()` | `string` | 逗号分隔格式 |
| `URI()` | `string` | URI格式 |

## 🎨 支持的输出格式

### 单行格式

| 格式 | 示例输出 | 使用场景 |
|------|----------|----------|
| `Version()` | `MyApp version v1.0.0 linux/amd64` | 标准版本显示 |
| `Short()` | `MyApp v1.0.0` | 短格式显示 |
| `Long()` | `MyApp version v1.0.0 linux/amd64 with go1.22.1` | 长格式显示 |
| `Standard()` | `MyApp v1.0.0 (abc1234) [linux/amd64]` | 标准格式显示 |
| `Full()` | `MyApp v1.0.0 (abc1234) linux/amd64 built 2024-01-01 with go1.22.1` | 包含完整信息 |
| `JSON()` | `{"name":"MyApp","version":"v1.0.0","commit":"abc1234","platform":"linux/amd64","buildTime":"2024-01-01T12:00:00Z","goVersion":"go1.22.1"}` | API返回或配置 |
| `KV()` | `app=MyApp version=v1.0.0 commit=abc1234 platform=linux/amd64 build=2024-01-01 go=go1.22.1` | 键值对格式 |
| `CSV()` | `MyApp,v1.0.0,abc1234,linux/amd64,2024-01-01,go1.22.1` | 逗号分隔格式 |
| `URI()` | `verman://MyApp/v1.0.0?commit=abc1234&platform=linux/amd64&build=2024-01-01&go=go1.22.1` | URI格式 |

### 多行格式

| 格式 | 示例输出 | 使用场景 |
|------|----------|----------|
| `Banner()` | `╔═══════════════════════════════════════╗`<br>`║           MyApp v1.0.0                ║`<br>`╠═══════════════════════════════════════╣`<br>`║ Commit:  abc1234 (clean)              ║`<br>`║ Platform: linux/amd64                 ║`<br>`║ Built:   2024-01-01 with go1.22.1     ║`<br>`╚═══════════════════════════════════════╝` | 程序启动横幅 |
| `Lines()` | `MyApp v1.0.0`<br>`Commit: abc1234 (clean)`<br>`Platform: linux/amd64`<br>`Built: 2024-01-01 with go1.22.1` | 分行显示 |
| `Table()` | `┌─────────────┬──────────────────┐`<br>`│ Application │ MyApp           │`<br>`│ Version     │ v1.0.0          │`<br>`│ Commit      │ abc1234         │`<br>`│ Platform    │ linux/amd64     │`<br>`│ Build Time  │ 2024-01-01      │`<br>`│ Go Version  │ go1.22.1        │`<br>`└─────────────┴──────────────────┘` | 详细信息表格 |

## ⚙️ 配置选项

### 编译时注入变量

通过 `-ldflags -X` 可以注入以下私有变量（**注意变量名为小写**）：

```bash
-X 'gitee.com/MM-Q/verman.appName=应用名称'
-X 'gitee.com/MM-Q/verman.gitVersion=版本号'
-X 'gitee.com/MM-Q/verman.gitCommit=提交哈希'
-X 'gitee.com/MM-Q/verman.gitTreeState=仓库状态'
-X 'gitee.com/MM-Q/verman.gitCommitTime=提交时间'
-X 'gitee.com/MM-Q/verman.buildTime=构建时间'
```

### 默认值

如果未注入相应值，将使用以下默认值：
- 所有字符串变量：`"unknown"`
- 运行时变量：自动获取当前环境信息

## 📁 项目结构

```
verman/
├── 📄 verman.go          # 主要库文件
├── 🧪 verman_test.go     # 测试文件
├── 📖 README.md          # 项目文档
├── 📋 APIDOC.md          # API 文档
├── 📜 LICENSE            # 许可证文件
├── 📦 go.mod             # Go 模块文件
├── 📁 script/            # 构建脚本目录
│   ├── build.bat         # Windows 构建脚本
│   ├── build.sh          # Linux/macOS 构建脚本
│   └── build.py          # 跨平台 Python 构建脚本
└── 📁 test/              # 测试项目
    ├── go.mod
    └── main.go
```

## 📄 许可证

本项目采用 [MIT 许可证](LICENSE)。

## 🤝 贡献指南

我们欢迎所有形式的贡献！

### 如何贡献

1. 🍴 Fork 本仓库
2. 🌿 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 💾 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 📤 推送到分支 (`git push origin feature/AmazingFeature`)
5. 🔄 创建 Pull Request

### 贡献类型

- 🐛 Bug 修复
- ✨ 新功能开发
- 📚 文档改进
- 🧪 测试用例添加
- 🎨 代码优化

## 📞 联系方式和相关链接

### 🔗 仓库地址
- **主仓库**: [https://gitee.com/MM-Q/verman](https://gitee.com/MM-Q/verman)

### 📋 相关资源
- 📖 [API 文档](APIDOC.md)
- 🐛 [问题反馈](https://gitee.com/MM-Q/verman/issues)
- 💡 [功能建议](https://gitee.com/MM-Q/verman/issues)
- 📦 [发布版本](https://gitee.com/MM-Q/verman/releases)

### 👨‍💻 维护者
- **MM-Q** - *项目创建者和主要维护者*

---

<div align="center">

**⭐ 如果这个项目对你有帮助，请给它一个 Star！**

[🏠 返回仓库首页](https://gitee.com/MM-Q/verman) | [📖 查看文档](APIDOC.md) | [🐛 报告问题](https://gitee.com/MM-Q/verman/issues)

</div>