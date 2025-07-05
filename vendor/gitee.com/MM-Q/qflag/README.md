# qflag

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/QiaoMuDe/qflag)

qflag 是一个用于解析命令行参数的 Go 语言库。它提供了丰富的功能，包括多种类型的标志（flag）、子命令支持、帮助信息生成等。

## 项目地址

该项目托管在 Gitee 和 GitHub 上，您可以访问以下链接查看源代码和最新动态：

- [Gitee](https://gitee.com/MM-Q/qflag.git)
- [GitHub](https://github.com/QiaoMuDe/qflag.git)

## 安装

要使用 qflag，您需要先安装 Go 环境。然后可以通过以下命令安装：

```bash
go get -u gitee.com/MM-Q/qflag
```

## 特性

- 支持多种类型的标志：字符串、整数、布尔值、64位浮点数、枚举、切片、时间间隔、64位整数、无符号16位整数、映射、路径、时间、IPv4地址、IPv6地址和URL。
- 支持子命令。
- 提供内置的帮助信息和安装路径显示功能。
- 支持自定义帮助信息。
- 线程安全设计。
- 循环引用检测。
- 动态帮助信息生成。
- 标志命名规则。
- 内置参数验证器，支持路径存在性、数值范围等验证。

## 标志类型

以下是qflag支持的所有标志类型及其说明：

| 标志类型 | 描述 | 示例 |
|----------|------|------|
| String | 字符串类型标志 | `--name "example"` |
| Int | 整数类型标志 | `--port 8080` |
| Int64 | 64位整数类型标志 | `--size 1073741824` |
| Uint16 | 无符号16位整数类型标志 | `--timeout 300` |
| Bool | 布尔类型标志 | `--debug` |
| Float64 | 64位浮点数类型标志 | `--threshold 0.95` |
| Enum | 枚举类型标志 | `--mode "debug"` |
| Slice | 切片类型标志 | `--files file1.txt,file2.txt` |
| Duration | 时间间隔类型标志 | `--timeout 30s` |
| Time | 时间类型标志 | `--start-time "2024-01-01T00:00:00"` |
| Map | 映射类型标志 | `--config key=value` |
| Path | 路径类型标志 | `--log-path "./logs"` |
| IP4 | IPv4地址类型标志 | `--server-ip 192.168.1.1` |
| IP6 | IPv6地址类型标志 | `--server-ipv6 ::1` |
| URL | URL类型标志 | `--api-url https://api.example.com` |

## 使用示例

### 基本使用示例

以下是一个简单的使用示例，展示如何创建一个命令并添加一些标志：

```go
package main

import (
	"fmt"
	"os"

	"gitee.com/MM-Q/qflag"
)

func main() {
	// 定义一个启动app的Bool型标志
	runF := qflag.Bool("run", "r", false, "run app")
	nameF := qflag.String("name", "n", "", "app name")
	pathF := qflag.String("path", "p", qflag.GetExecutablePath(), "app path")

	// 解析参数
	if err := qflag.Parse(); err != nil {
		fmt.Printf("解析参数错误: %v\n", err)
		os.Exit(1)
	}

	// 获取参数值
	if runF.Get() {
		fmt.Printf("启动app: %s\n", nameF.Get())
		fmt.Printf("app路径: %s\n", pathF.Get())
	}
}

```

### 子命令示例

以下是一个使用子命令的示例：

```go
package main

import (
	"flag"
	"fmt"
	"os"

	"gitee.com/MM-Q/qflag"
)

func main() {
	// 定义一个启动app的Bool型标志
	runF := qflag.Bool("run", "r", false, "运行app")
	nameF := qflag.String("name", "n", "", "指定app名称")
	pathF := qflag.String("path", "p", qflag.GetExecutablePath(), "指定app路径")

	// 定义子命令
	stopCmd := qflag.NewCmd("stop", "st", flag.ExitOnError)
	stopF := stopCmd.Bool("stop", "s", false, "停止app")
	stopCmdName := stopCmd.String("name", "n", "", "指定app名称")

	// 添加子命令描述
	stopCmd.SetDescription("停止app")

	// 添加子命令到全局QCommandLine
	qflag.AddSubCmd(stopCmd)

	// 解析命令
	if err := qflag.Parse(); err != nil {
		fmt.Printf("解析参数错误: %v\n", err)
		os.Exit(1)
	}

	// 获取参数值
	if runF.Get() {
		fmt.Printf("启动app: %s\n", nameF.GetValue())
		fmt.Printf("app路径: %s\n", pathF.GetValue())
	}

	// 判断是否执行了stop子命令
	if stopF.Get() {
		fmt.Printf("停止app: %s\n", stopCmdName.Get())
	}
}
```

## API文档

qflag提供了完善的API文档，按模块组织如下：

- **全局命令处理**: [qflag包文档](./APIDOC.md) - 包含全局命令创建、参数解析和子命令管理相关API
- **核心命令处理**: [cmd包文档](./cmd/APIDOC.md) - 包含命令创建、参数解析和子命令管理相关API
- **标志类型定义**: [flags包文档](./flags/APIDOC.md) - 包含所有标志类型的详细定义和使用方法
- **错误处理**: [qerr包文档](./qerr/APIDOC.md) - 包含错误类型和处理相关API
- **参数验证**: [validator包文档](./validator/APIDOC.md) - 包含参数验证器接口和实现

完整的API文档可通过访问各模块对应的APIDOC.md文件查看。

## 许可证

本项目采用 MIT 许可证。详情请参阅 [LICENSE](LICENSE) 文件。
