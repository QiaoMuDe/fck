# qflag

[![Go Version](https://img.shields.io/badge/Go-1.24.4-blue.svg)](https://golang.org/)
[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/QiaoMuDe/qflag)

qflag 是一个功能强大的 Go 语言命令行参数解析库，支持多种数据类型、子命令管理、参数验证、自动补全等高级特性。该库提供了简洁易用的 API，同时保持了高度的灵活性和扩展性。

## 项目地址

该项目托管在 Gitee 和 GitHub 上，您可以访问以下链接查看源代码和最新动态：

- [Gitee](https://gitee.com/MM-Q/qflag.git)
- [GitHub](https://github.com/QiaoMuDe/qflag.git)

## 安装

使用 `go get` 命令安装：

```bash
go get -u gitee.com/MM-Q/qflag
```

然后在代码中导入：

```go
import "gitee.com/MM-Q/qflag"
import "gitee.com/MM-Q/qflag/cmd"
import "gitee.com/MM-Q/qflag/flags"
import "gitee.com/MM-Q/qflag/validator"
```

## 核心特性

### 🚀 丰富的数据类型支持
- **基础类型**：字符串、整数（int/int64/uint16/uint32/uint64）、布尔值、浮点数
- **高级类型**：枚举、时间间隔、时间、切片、映射
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
- **错误处理**：详细的错误类型和信息，便于调试

### 🛡️ 企业级特性
- **并发安全**：使用 `sync.RWMutex` 保证线程安全
- **内存优化**：高效的内存使用和垃圾回收友好设计
- **扩展性**：模块化架构，易于扩展和定制

## 支持的标志类型

| 标志类型 | 创建函数 | 绑定函数 | 描述 | 示例 |
|----------|----------|----------|------|------|
| **基础类型** |
| `StringFlag` | `String()` | `StringVar()` | 字符串类型 | `--name "example"` |
| `IntFlag` | `Int()` | `IntVar()` | 32位整数 | `--port 8080` |
| `Int64Flag` | `Int64()` | `Int64Var()` | 64位整数 | `--size 1073741824` |
| `Uint16Flag` | `Uint16()` | `Uint16Var()` | 16位无符号整数 | `--timeout 300` |
| `Uint32Flag` | `Uint32()` | `Uint32Var()` | 32位无符号整数 | `--max-conn 1000` |
| `Uint64Flag` | `Uint64()` | `Uint64Var()` | 64位无符号整数 | `--max-size 9223372036854775807` |
| `BoolFlag` | `Bool()` | `BoolVar()` | 布尔类型 | `--debug` |
| `Float64Flag` | `Float64()` | `Float64Var()` | 64位浮点数 | `--threshold 0.95` |
| **高级类型** |
| `EnumFlag` | `Enum()` | `EnumVar()` | 枚举类型 | `--mode "debug"` |
| `SliceFlag` | `Slice()` | `SliceVar()` | 字符串切片 | `--files file1,file2` |
| `DurationFlag` | `Duration()` | `DurationVar()` | 时间间隔 | `--timeout 30s` |
| `TimeFlag` | `Time()` | `TimeVar()` | 时间类型 | `--start "2024-01-01T00:00:00"` |
| `MapFlag` | `Map()` | `MapVar()` | 键值对映射 | `--config key=value,key2=value2` |

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
    name := qflag.String("name", "n", "world", "要问候的名称")
    count := qflag.Int("count", "c", 1, "问候次数")
    verbose := qflag.Bool("verbose", "v", false, "详细输出")
    
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

### 子命令示例

```go
package main

import (
    "flag"
    "fmt"
    "os"
    "gitee.com/MM-Q/qflag"
)

func main() {
    // 全局标志
    verbose := qflag.Bool("verbose", "v", false, "详细输出")
    
    // 创建子命令
    startCmd := qflag.NewCmd("start", "s", flag.ExitOnError)
    startCmd.SetDescription("启动服务")
    
    // 为子命令添加标志
    port := startCmd.Int("port", "p", 8080, "服务端口")
    host := startCmd.String("host", "h", "localhost", "服务主机")
    
    // 创建另一个子命令
    stopCmd := qflag.NewCmd("stop", "st", flag.ExitOnError)
    stopCmd.SetDescription("停止服务")
    
    pidFile := stopCmd.String("pid-file", "f", "/var/run/app.pid", "PID文件路径")
    
    // 注册子命令
    qflag.AddSubCmd(startCmd, stopCmd)
    
    // 解析参数
    if err := qflag.Parse(); err != nil {
        fmt.Printf("解析参数错误: %v\n", err)
        os.Exit(1)
    }
    
    // 处理命令逻辑
    if startCmd.IsParsed() {
        if verbose.Get() {
            fmt.Printf("启动服务在 %s:%d\n", host.Get(), port.Get())
        }
        // 启动服务逻辑...
    } else if stopCmd.IsParsed() {
        if verbose.Get() {
            fmt.Printf("从 %s 读取PID并停止服务\n", pidFile.Get())
        }
        // 停止服务逻辑...
    }
}
```

使用方式：
```bash
./app start --port 9000 --host 0.0.0.0 --verbose
./app stop --pid-file /tmp/app.pid -v
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
    logLevel := qflag.Enum("log-level", "l", "info", 
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
    // 创建切片标志
    files := qflag.Slice("files", "f", []string{}, "要处理的文件列表")
    
    // 自定义分隔符（默认为逗号）
    files.SetSeparator(";")
    
    if err := qflag.Parse(); err != nil {
        fmt.Printf("解析参数错误: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Printf("要处理的文件: %v\n", files.Get())
}
```

使用方式：`./app --files file1.txt;file2.txt;file3.txt`

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
    config := qflag.Map("config", "c", map[string]string{}, "配置键值对")
    
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
    // 创建带验证的标志
    port := qflag.Int("port", "p", 8080, "服务端口")
    
    // 设置端口范围验证器
    port.SetValidator(&validator.IntRangeValidator{
        Min: 1024,
        Max: 65535,
    })
    
    // 字符串长度验证
    name := qflag.String("name", "n", "", "服务名称")
    name.SetValidator(&validator.StringLengthValidator{
        Min: 3,
        Max: 20,
    })
    
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
    // 创建标志并绑定环境变量
    dbHost := qflag.String("db-host", "", "localhost", "数据库主机")
    dbHost.BindEnv("DATABASE_HOST")
    
    dbPort := qflag.Int("db-port", "", 5432, "数据库端口")
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
export DATABASE_HOST=prod-db.example.com
export DATABASE_PORT=3306
./app  # 将使用环境变量的值
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
    "gitee.com/MM-Q/qflag/flags"
)

// 自定义邮箱验证器
type EmailValidator struct{}

func (v *EmailValidator) Validate(value any) error {
    email, ok := value.(string)
    if !ok {
        return errors.New("value is not a string")
    }
    
    if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
        return errors.New("invalid email format")
    }
    
    return nil
}

func main() {
    email := qflag.String("email", "e", "", "用户邮箱")
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
./your-app --generate-shell-completion bash > your-app-completion.sh

# 安装补全脚本
sudo cp your-app-completion.sh /etc/profile.d/
source /etc/profile.d/your-app-completion.sh
```

### PowerShell 补全

```powershell
# 生成 PowerShell 补全脚本
./your-app.exe --generate-shell-completion pwsh > your-app-completion.ps1

# 安装补全脚本
. ./your-app-completion.ps1
```

## 帮助信息定制

```go
package main

import (
    "gitee.com/MM-Q/qflag"
)

func main() {
    // 设置应用信息
    qflag.SetVersion("1.0.0")
    qflag.SetDescription("这是一个示例应用程序")
    qflag.SetUsageSyntax("myapp [选项] <命令> [参数...]")
    
    // 添加使用示例
    qflag.AddExample("启动服务", "myapp start --port 8080")
    qflag.AddExample("查看状态", "myapp status --verbose")
    
    // 添加注意事项
    qflag.AddNote("配置文件默认位置: ~/.myapp/config.yaml")
    qflag.AddNote("日志文件位置: /var/log/myapp.log")
    
    // 设置中文帮助信息
    qflag.SetUseChinese(true)
    
    // 定义标志...
    name := qflag.String("name", "n", "world", "要问候的名称")
    
    if err := qflag.Parse(); err != nil {
        return
    }
    
    // 应用逻辑...
}
```

## 项目架构

qflag 采用模块化设计，主要包含以下包：

- **`qflag`** - 主包，提供全局 API 和便捷函数
- **`cmd`** - 命令管理，处理子命令和命令树结构
- **`flags`** - 标志类型定义，包含所有标志类型的实现
- **`validator`** - 参数验证器，提供常用验证器和验证接口
- **`qerr`** - 错误处理，定义错误类型和错误处理机制
- **`utils`** - 工具函数，提供通用的辅助功能
- **`internal`** - 内部实现包，包含核心解析逻辑

## API 文档

完整的 API 文档按模块组织：

- **[qflag 包文档](./APIDOC.md)** - 全局 API 和便捷函数
- **[cmd 包文档](./cmd/APIDOC.md)** - 命令管理相关 API
- **[flags 包文档](./flags/APIDOC.md)** - 标志类型定义和使用方法
- **[validator 包文档](./validator/APIDOC.md)** - 参数验证器接口和实现
- **[qerr 包文档](./qerr/APIDOC.md)** - 错误处理相关 API
- **[utils 包文档](./utils/APIDOC.md)** - 工具函数文档

## 性能特性

- **内存效率**：优化的内存分配策略，减少 GC 压力
- **并发安全**：全面的线程安全保护，支持并发访问
- **解析速度**：高效的参数解析算法，适合大型应用
- **类型安全**：编译时类型检查，避免运行时类型错误

## 兼容性

- **Go 版本**：要求 Go 1.24+ （支持泛型）
- **操作系统**：支持 Windows、Linux、macOS
- **Shell 支持**：Bash、PowerShell

## 贡献指南

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

**qflag** - 让命令行参数解析变得简单而强大！ 🚀