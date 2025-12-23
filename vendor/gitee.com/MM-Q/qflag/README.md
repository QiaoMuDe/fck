<div align="center">

# 🚀 qflag

**功能强大的 Go 语言命令行参数解析库**

[![Go Version](https://img.shields.io/badge/Go-1.24.0-blue.svg)](https://golang.org/)
[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![Gitee](https://img.shields.io/badge/Gitee-qflag-red.svg)](https://gitee.com/MM-Q/qflag)
[![GitHub](https://img.shields.io/badge/GitHub-qflag-black.svg)](https://github.com/QiaoMuDe/qflag)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/QiaoMuDe/qflag)

*泛型设计 • 自动路由 • 类型安全 • 并发安全 • 自动补全 • 子命令管理*

[📖 快速开始](#快速开始) • [🔧 安装指南](#安装) • [📚 API 文档](#api-文档) • [🤝 贡献指南](#贡献指南)

</div>

---

## ✨ 项目简介

qflag 是一个基于 Go 泛型的现代化命令行参数解析库，对标准库 flag 进行了全面增强。它采用模块化架构设计，提供了 16+ 种标志类型（包括基础类型、切片类型、复杂类型如枚举、时间、映射、大小等）、完整的子命令系统、智能的自动路由机制、强大的参数验证框架、智能的 Shell 自动补全（支持 Bash/PowerShell）、环境变量绑定等企业级特性。通过泛型设计确保类型安全，内置并发保护机制，支持中英文帮助信息，为构建专业的 CLI 应用提供了完整的解决方案。

## 🔗 项目地址

该项目同时托管在 Gitee 和 GitHub 上，您可以选择合适的平台访问：

| 平台               | 地址                                                        | 描述                       |
| ------------------ | ----------------------------------------------------------- | -------------------------- |
| 🔴**Gitee**  | [gitee.com/MM-Q/qflag](https://gitee.com/MM-Q/qflag)           | 国内访问更快，主要开发仓库 |
| ⚫**GitHub** | [github.com/QiaoMuDe/qflag](https://github.com/QiaoMuDe/qflag) | 国际化平台，同步更新       |

## 安装

使用 `go get` 命令安装：

```bash
go get -u gitee.com/MM-Q/qflag
```

然后在代码中导入：

```go
import "gitee.com/MM-Q/qflag"
```

## 核心特性

### 🚀 丰富的数据类型支持

- **基础类型**：字符串、整数（int/int64/uint16/uint32/uint64）、布尔值、浮点数
- **高级类型**：枚举、时间间隔、时间、切片([]string, []int64, []int)、映射、大小
- **泛型设计**：基于 Go 泛型的类型安全标志系统

### 🎯 强大的命令管理

- **子命令支持**：构建复杂的命令树结构
- **命令嵌套**：支持多层级子命令
- **命令别名**：长短名称支持，提升用户体验

### ✅ 完善的参数验证

- **内置验证器**：字符串长度、数值范围、正则表达式、路径验证等
- **自定义验证器**：实现 `Validator` 接口，支持复杂业务逻辑验证
- **类型安全**：编译时类型检查，运行时验证保障

### 🔧 便捷的开发体验

- **自动补全**：支持 Bash 和 PowerShell 的自动补全脚本生成
- **环境变量绑定**：标志可自动从环境变量加载默认值
- **帮助信息生成**：自动生成格式化的帮助文档，支持中英文
- **执行函数接口**：通过 `SetRun` 和 `Run` 方法提供灵活的命令执行逻辑定义，支持并发安全
- **自动路由**：`ParseAndRoute` 方法支持自动解析参数并路由到对应的子命令，简化命令行应用开发
- **错误处理**：详细的错误类型和信息，便于调试

### 🛡️ 企业级特性

- **并发安全**：使用 `sync.RWMutex` 保证线程安全
- **内存优化**：高效的内存使用和垃圾回收友好设计
- **扩展性**：模块化架构，易于扩展和定制

## 支持的标志类型

| 标志类型            | 创建函数          | 绑定函数             | 描述           | 示例                               |
| ------------------- | ----------------- | -------------------- | -------------- | ---------------------------------- |
| **基础类型**  |                   |                      |                |                                    |
| `StringFlag`      | `String()`      | `StringVar()`      | 字符串类型     | `--name "example"`               |
| `IntFlag`         | `Int()`         | `IntVar()`         | 32位整数       | `--port 8080`                    |
| `Int64Flag`       | `Int64()`       | `Int64Var()`       | 64位整数       | `--size 1073741824`              |
| `Uint16Flag`      | `Uint16()`      | `Uint16Var()`      | 16位无符号整数 | `--timeout 300`                  |
| `Uint32Flag`      | `Uint32()`      | `Uint32Var()`      | 32位无符号整数 | `--max-conn 1000`                |
| `Uint64Flag`      | `Uint64()`      | `Uint64Var()`      | 64位无符号整数 | `--max-size 9223372036854775807` |
| `BoolFlag`        | `Bool()`        | `BoolVar()`        | 布尔类型       | `--debug`                        |
| `Float64Flag`     | `Float64()`     | `Float64Var()`     | 64位浮点数     | `--threshold 0.95`               |
| **高级类型**  |                   |                      |                |                                    |
| `EnumFlag`        | `Enum()`        | `EnumVar()`        | 枚举类型       | `--mode "debug"`                 |
| `StringSliceFlag` | `StringSlice()` | `StringSliceVar()` | 字符串切片     | `--files file1,file2`            |
| `IntSliceFlag`    | `IntSlice()`    | `IntSliceVar()`    | 整数切片       | `--ports 8080,9000,3000`         |
| `Int64SliceFlag`  | `Int64Slice()`  | `Int64SliceVar()`  | 64位整数切片   | `--sizes 1024,2048,4096`         |
| `DurationFlag`    | `Duration()`    | `DurationVar()`    | 时间间隔       | `--timeout 30s`                  |
| `TimeFlag`        | `Time()`        | `TimeVar()`        | 时间类型       | `--start "2024-01-01T00:00:00"`  |
| `MapFlag`         | `Map()`         | `MapVar()`         | 键值对映射     | `--config key=value,key2=value2` |
| `SizeFlag`        | `Size()`        | `SizeVar()`        | 大小类型       | `--max-size 1024MB`              |

## 快速开始

### 基本使用示例

```go
package main

import (
    "fmt"
    "os"
    "gitee.com/MM-Q/qflag"
)

func main() {
    // 定义标志
    name := qflag.Root.String("name", "n", "world", "要问候的名称")
    count := qflag.Root.Int("count", "c", 1, "问候次数")
    verbose := qflag.Root.Bool("verbose", "v", false, "详细输出")
  
    // 解析命令行参数
    if err := qflag.Parse(); err != nil {
        fmt.Printf("解析参数错误: %v\n", err)
        os.Exit(1)
    }
  
    // 使用参数值
    for i := 0; i < count.Get(); i++ {
        if verbose.Get() {
            fmt.Printf("第 %d 次问候: ", i+1)
        }
        fmt.Printf("Hello, %s!\n", name.Get())
    }
}
```

使用方式：

```bash
./app --name "Alice" --count 3 --verbose
./app -n "Bob" -c 2 -v
```

### 子命令示例（自动路由方式）

```go
package main

import (
    "fmt"
    "os"
    "gitee.com/MM-Q/qflag"
)

func main() {
    // 创建根命令
    rootCmd := qflag.NewCmd("myapp", "", qflag.ExitOnError)
    rootCmd.SetDesc("示例应用程序")
    
    // 全局标志
    verbose := rootCmd.Bool("verbose", "v", false, "详细输出")
  
    // 创建启动服务子命令
    startCmd := qflag.NewCmd("start", "s", qflag.ExitOnError)
    startCmd.SetDesc("启动服务")
  
    // 为子命令添加标志
    port := startCmd.Int("port", "p", 8080, "服务端口")
    host := startCmd.String("host", "h", "localhost", "服务主机")
    
    // 设置启动服务的执行函数
    startCmd.SetRun(func(cmd *qflag.Cmd) error {
        if verbose.Get() {
            fmt.Printf("启动服务在 %s:%d\n", host.Get(), port.Get())
        }
        // 启动服务逻辑...
        fmt.Printf("服务启动成功！监听地址: %s:%d\n", host.Get(), port.Get())
        return nil
    })
  
    // 创建停止服务子命令
    stopCmd := qflag.NewCmd("stop", "st", qflag.ExitOnError)
    stopCmd.SetDesc("停止服务")
  
    pidFile := stopCmd.String("pid-file", "f", "/var/run/app.pid", "PID文件路径")
    
    // 设置停止服务的执行函数
    stopCmd.SetRun(func(cmd *qflag.Cmd) error {
        if verbose.Get() {
            fmt.Printf("从 %s 读取PID并停止服务\n", pidFile.Get())
        }
        // 停止服务逻辑...
        fmt.Printf("服务停止成功！PID文件: %s\n", pidFile.Get())
        return nil
    })
  
    // 注册子命令到根命令
    rootCmd.AddSubCmd(startCmd, stopCmd)
  
    // 使用ParseAndRoute自动解析并路由到对应子命令
    if err := rootCmd.ParseAndRoute(os.Args[1:]); err != nil {
        fmt.Printf("错误: %v\n", err)
        os.Exit(1)
    }
}
```

使用方式：

```bash
# 启动服务
./myapp start --port 9000 --host 0.0.0.0 --verbose
# 输出: 启动服务在 0.0.0.0:9000
# 输出: 服务启动成功！监听地址: 0.0.0.0:9000

# 停止服务
./myapp stop --pid-file /tmp/app.pid -v
# 输出: 从 /tmp/app.pid 读取PID并停止服务
# 输出: 服务停止成功！PID文件: /tmp/app.pid

# 查看帮助
./myapp --help
./myapp start --help
./myapp stop --help
```

### Run函数执行示例

```go
package main

import (
    "fmt"
    "os"
    "gitee.com/MM-Q/qflag"
)

func main() {
    // 创建命令
    serverCmd := qflag.NewCmd("server", "s", qflag.ExitOnError)
    port := serverCmd.Int("port", "p", 8080, "服务器端口")
    debug := serverCmd.Bool("debug", "d", false, "调试模式")
    
    // 设置执行函数 - 使用SetRun方法
    serverCmd.SetRun(func(cmd *qflag.Cmd) error {
        fmt.Printf("启动服务器: localhost:%d (调试模式: %v)\n", port.Get(), debug.Get())
        // 这里放置实际的服务器启动逻辑
        return nil
    })
    
    // 使用全局函数添加到根命令
    qflag.AddSubCmd(serverCmd)
    
    // 使用全局函数解析参数
    if err := qflag.Parse(); err != nil {
        fmt.Printf("解析错误: %v\n", err)
        os.Exit(1)
    }
    
    // 直接执行Run方法 - 内部会自动检查是否已解析
    if err := serverCmd.Run(); err != nil {
        fmt.Printf("执行错误: %v\n", err)
        os.Exit(1)
    }
}
```

使用方式：

```bash
./app server --port 3000 --debug
# 输出: 启动服务器: localhost:3000 (调试模式: true)
```

## 高级功能示例

### 1. 枚举类型标志

```go
package main

import (
    "fmt"
    "os"
    "gitee.com/MM-Q/qflag"
)

func main() {
    // 创建枚举标志
    logLevel := qflag.Root.Enum("log-level", "l", "info", 
        "日志级别", []string{"debug", "info", "warn", "error"})
  
    // 设置大小写敏感（可选）
    logLevel.SetCaseSensitive(false)
  
    if err := qflag.Parse(); err != nil {
        fmt.Printf("解析参数错误: %v\n", err)
        os.Exit(1)
    }
  
    fmt.Printf("当前日志级别: %s\n", logLevel.Get())
}
```

### 2. 切片类型标志

```go
package main

import (
    "fmt"
    "os"
    "gitee.com/MM-Q/qflag"
)

func main() {
    // 创建各种切片标志
    files := qflag.Root.StringSlice("files", "f", []string{}, "要处理的文件列表")
    ports := qflag.Root.IntSlice("ports", "p", []int{8080}, "服务端口列表") 
    sizes := qflag.Root.Int64Slice("sizes", "s", []int64{}, "文件大小列表")
    
    // 自定义分隔符（可选）
    files.SetDelimiters([]string{";"})
    
    if err := qflag.Parse(); err != nil {
        fmt.Printf("解析参数错误: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Printf("文件: %v, 端口: %v, 大小: %v\n", files.Get(), ports.Get(), sizes.Get())
}
```

使用方式：

```bash
./app --files file1.txt;file2.txt;file3.txt --ports 8080,9000,3000 --sizes 1024,2048,4096
```

### 3. 映射类型标志

```go
package main

import (
    "fmt"
    "os"
    "gitee.com/MM-Q/qflag"
)

func main() {
    // 创建映射标志
    config := qflag.Root.Map("config", "c", map[string]string{}, "配置键值对")
  
    // 设置分隔符（键值对分隔符，键值分隔符）
    config.SetDelimiters(",", ":")
  
    if err := qflag.Parse(); err != nil {
        fmt.Printf("解析参数错误: %v\n", err)
        os.Exit(1)
    }
  
    fmt.Printf("配置: %v\n", config.Get())
}
```

使用方式：`./app --config server:localhost,port:8080,debug:true`

### 4. 参数验证

```go
package main

import (
    "fmt"
    "os"
    "gitee.com/MM-Q/qflag"
    "gitee.com/MM-Q/qflag/validator"
)

func main() {
    // 端口范围验证（1024-65535）
    port := qflag.Root.Int("port", "p", 8080, "服务端口（1024-65535）")
    port.SetValidator(&validator.IntRangeValidator{Min: 1024, Max: 65535})
    
    // 字符串长度验证（3-20字符）
    name := qflag.Root.String("name", "n", "", "服务名称（3-20字符）")
    name.SetValidator(&validator.StringLengthValidator{Min: 3, Max: 20})
    
    if err := qflag.Parse(); err != nil {
        fmt.Printf("解析参数错误: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Printf("服务 %s 将在端口 %d 启动\n", name.Get(), port.Get())
}
```

### 5. 环境变量绑定

```go
package main

import (
    "fmt"
    "os"
    "gitee.com/MM-Q/qflag"
)

func main() {
    // 绑定环境变量（DATABASE_HOST, DATABASE_PORT）
    dbHost := qflag.Root.String("db-host", "", "localhost", "数据库主机")
    dbPort := qflag.Root.Int("db-port", "", 5432, "数据库端口")
    dbHost.BindEnv("DATABASE_HOST")
    dbPort.BindEnv("DATABASE_PORT")
    
    if err := qflag.Parse(); err != nil {
        fmt.Printf("解析参数错误: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Printf("连接数据库: %s:%d\n", dbHost.Get(), dbPort.Get())
}
```

使用方式：

```bash
export DATABASE_HOST=prod-db.example.com DATABASE_PORT=3306
./app  # 使用环境变量
./app --db-host localhost --db-port 5432  # 命令行参数优先级更高
```

### 6. 自定义验证器

```go
package main

import (
    "errors"
    "fmt"
    "os"
    "strings"
    "gitee.com/MM-Q/qflag"
)

// 超简化的自定义验证器
type EmailValidator struct{}
func (v *EmailValidator) Validate(value any) error {
    email, _ := value.(string)
    if !strings.Contains(email, "@") {
        return errors.New("邮箱必须包含@符号")
    }
    return nil
}

func main() {
    email := qflag.Root.String("email", "e", "", "用户邮箱")
    email.SetValidator(&EmailValidator{})
    
    if err := qflag.Parse(); err != nil {
        fmt.Printf("解析参数错误: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Printf("用户邮箱: %s\n", email.Get())
}
```

## 自动补全

qflag 支持为 Bash 和 PowerShell 生成自动补全脚本：

### Bash 补全

```bash
# 生成 Bash 补全脚本
./your-app --completion bash > your-app-completion.sh

# 安装补全脚本
sudo cp your-app-completion.sh /etc/profile.d/
source /etc/profile.d/your-app-completion.sh
```

### PowerShell 补全

```powershell
# 生成 PowerShell 补全脚本
./your-app.exe --completion pwsh > your-app-completion.ps1

# 安装补全脚本
. ./your-app-completion.ps1
```

## 帮助信息定制

```go
package main

import (
    "gitee.com/MM-Q/qflag"
    "gitee.com/MM-Q/qflag/internal/types"
)

func main() {
    // 使用全局函数配置根命令
    qflag.ApplyConfig(types.CmdConfig{
        Version: "1.0.0",
        Desc: "这是一个示例应用程序",
        UseChinese: true,
        Examples: []types.ExampleInfo{
            {Desc: "启动服务", Usage: "myapp start --port 8080"},
            {Desc: "查看状态", Usage: "myapp status --verbose"},
        },
    })
    
    // 定义标志并使用...
    name := qflag.Root.String("name", "n", "world", "要问候的名称")
    
    if err := qflag.Parse(); err != nil {
        return
    }
    
    // 应用逻辑...
}
```

## 项目架构

qflag 采用模块化设计，主要包含以下包：

- **`qflag`** - 主包，提供全局 API 和便捷函数
- **`flags`** - 标志类型定义，包含所有标志类型的实现
- **`validator`** - 参数验证器，提供常用验证器和验证接口
- **`qerr`** - 错误处理，定义错误类型和错误处理机制
- **`utils`** - 工具函数，提供通用的辅助功能
- **`internal`** - 内部实现包，包含核心解析逻辑

## API 文档

完整的 API 文档按模块组织：

- **[qflag 包文档](./APIDOC.md)** - 全局 API 和便捷函数
- **[flags 包文档](./flags/APIDOC.md)** - 标志类型定义和使用方法
- **[validator 包文档](./validator/APIDOC.md)** - 参数验证器接口和实现
- **[qerr 包文档](./qerr/APIDOC.md)** - 错误处理相关 API

## 性能特性

- **内存效率**：优化的内存分配策略，减少 GC 压力
- **并发安全**：全面的线程安全保护，支持并发访问
- **解析速度**：高效的参数解析算法，适合大型应用
- **类型安全**：编译时类型检查，避免运行时类型错误

## 兼容性

- **Go 版本**：要求 Go 1.24+ （支持泛型）
- **操作系统**：支持 Windows、Linux、macOS
- **Shell 支持**：Bash、PowerShell

## 🧪 测试说明

qflag 提供了完整的测试套件，确保代码质量和功能稳定性。

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行测试并显示覆盖率
go test -cover ./...

# 生成详细的覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# 运行基准测试
go test -bench=. ./...

# 运行特定包的测试
go test ./flags
go test ./cmd
go test ./validator
```

### 测试覆盖率目标

- **整体覆盖率**：> 90%
- **核心包覆盖率**：> 95%
- **关键功能**：100% 覆盖

### 持续集成

项目配置了自动化测试流程：

- **代码质量检查**：使用 `golangci-lint` 进行静态分析
- **多版本测试**：在 Go 1.24+ 版本上测试
- **跨平台测试**：Windows、Linux、macOS 环境验证

## 🤝 贡献指南

我们欢迎社区贡献！请遵循以下步骤：

1. Fork 项目到您的 GitHub/Gitee 账户
2. 创建特性分支：`git checkout -b feature/amazing-feature`
3. 提交更改：`git commit -m 'Add amazing feature'`
4. 推送分支：`git push origin feature/amazing-feature`
5. 创建 Pull Request

### 开发规范

- 遵循 Go 官方代码规范
- 添加适当的单元测试
- 更新相关文档
- 确保所有测试通过

## 许可证

本项目采用 MIT 许可证。详情请参阅 [LICENSE](LICENSE) 文件。

## 支持与反馈

- **问题报告**：[Gitee Issues](https://gitee.com/MM-Q/qflag/issues)
- **功能请求**：[GitHub Issues](https://github.com/QiaoMuDe/qflag/issues)
- **讨论交流**：欢迎在 Issues 中讨论使用问题和改进建议

---

<div align="center">

**qflag** - 让命令行参数解析变得简单而强大！ 🚀

</div>
