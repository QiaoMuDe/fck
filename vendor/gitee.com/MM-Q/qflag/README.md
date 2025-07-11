# qflag

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/QiaoMuDe/qflag)

qflag 是一个用于解析和管理命令行参数的 Go 语言库，支持多种参数类型和高级功能如自动补全、环境变量加载、参数验证等。

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
import "gitee.com/MM-Q/qflag/"
import "gitee.com/MM-Q/qflag/cmd"
import "gitee.com/MM-Q/qflag/flags"
import "gitee.com/MM-Q/qflag/qerr"
```

## 特性

- **支持多种参数类型**：包括布尔值、整数、浮点数、字符串、时间、IP 地址、URL、枚举、路径等。
- **子命令支持**：提供对子命令的管理，可以构建复杂的命令树。
- **自动补全**：支持 Bash 和 PowerShell 的自动补全脚本生成。
- **参数验证**：提供参数验证接口，确保输入值符合预期。
- **帮助文档生成**：支持生成中英文帮助文档，可自定义描述、示例、用法等。
- **环境变量支持**：可从环境变量加载参数值。
- **错误处理**：提供详细的错误类型和错误信息，便于调试和处理异常情况。


## 标志类型

以下是qflag支持的所有标志类型及其说明：

| 标志类型 | 描述 | 示例 | 注意事项 |
|----------|------|------|--------|
| `StringFlag` / `StringVar` | 字符串类型标志 | `--name "example"` | 无 |
| `IntFlag` / `IntVar` | 整数类型标志 | `--port 8080` | 无 |
| `Int64Flag` / `Int64Var` | 64位整数类型标志 | `--size 1073741824` | 无 |
| `Uint16Flag` / `Uint16Var` | 无符号16位整数类型标志 | `--timeout 300` | 无 |
| `Uint32Flag` / `Uint32Var` | 无符号32位整数类型标志 | `--max-connections 1000` | 无 |
| `Uint64Flag` / `Uint64Var` | 无符号64位整数类型标志 | `--max-memory 9223372036854775807` | 无 |
| `BoolFlag` / `BoolVar` | 布尔类型标志 | `--debug` | 无 |
| `Float64Flag` / `Float64Var` | 64位浮点数类型标志 | `--threshold 0.95` | 无 |
| `EnumFlag` / `EnumVar` | 枚举类型标志 | `--mode "debug"` | 支持大小写敏感设置，通过`SetCaseSensitive(true)`启用 |
| `SliceFlag` / `SliceVar` | 切片类型标志，支持自定义分隔符 | `--files file1.txt,file2.txt` | 默认使用逗号分隔，可通过`SetSeparator`修改 |
| `DurationFlag` / `DurationVar` | 时间间隔类型标志 | `--timeout 30s` | 无 |
| `TimeFlag` / `TimeVar` | 时间类型标志 | `--start-time "2024-01-01T00:00:00"` | 无 |
| `MapFlag` / `MapVar` | 映射类型标志 | `--config key=value,key2=value2` | 支持`=`和`:`作为键值分隔符，可通过`SetDelimiters`修改 |
| `PathFlag` / `PathVar` | 路径类型标志，支持路径验证 | `--log-path "./logs"` | 可通过`IsDirectory(true)`和`MustExist(true)`设置验证规则 |
| `IP4Flag` / `IP4Var` | IPv4地址类型标志 | `--server-ip 192.168.1.1` | 无 |
| `IP6Flag` / `IP6Var` | IPv6地址类型标志 | `--server-ipv6 ::1` | 无 |
| `URLFlag` / `URLVar` | URL类型标志 | `--api-url https://api.example.com` |	无 |

## 使用示例

### SliceFlag 示例

展示如何使用切片标志并自定义分隔符：

```go
package main

import (
	"fmt"
	"os"

	"gitee.com/MM-Q/qflag"
)

func main() {
	// 创建切片标志，默认使用逗号分隔
	filesF := qflag.Slice("files", "f", []string{}, "要处理的文件列表")
	// 设置自定义分隔符为分号
	filesF.SetSeparator(";")

	if err := qflag.Parse(); err != nil {
		fmt.Printf("解析参数错误: %v\n", err)
		os.Exit(1)
	}

	// 获取切片值
	files := filesF.Get()
	fmt.Printf("要处理的文件: %v\n", files)
}
```

使用方式：`./app --files file1.txt;file2.txt;file3.txt`

### MapFlag 示例

展示如何使用映射标志并自定义分隔符：

```go
package main

import (
	"fmt"
	"os"

	"gitee.com/MM-Q/qflag"
)

func main() {
	// 创建映射标志
	configF := qflag.Map("config", "c", map[string]string{}, "配置键值对")
	// 设置键值对分隔符为逗号，键值分隔符为冒号
	configF.SetDelimiters(",", ":")

	if err := qflag.Parse(); err != nil {
		fmt.Printf("解析参数错误: %v\n", err)
		os.Exit(1)
	}

	// 获取映射值
	config := configF.Get()
	fmt.Printf("配置: %v\n", config)
}
```

使用方式：`./app --config server:localhost,port:8080,timeout:30s`

### PathFlag 示例

展示路径标志的路径验证功能：

```go
package main

import (
	"fmt"
	"os"

	"gitee.com/MM-Q/qflag"
)

func main() {
	// 创建路径标志并设置验证规则
	logPathF := qflag.Path("log-path", "l", "/var/log/app", "日志目录")
		.IsDirectory(true)  // 必须是目录
		.MustExist(false)   // 路径不存在时自动创建

	if err := qflag.Parse(); err != nil {
		fmt.Printf("解析参数错误: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("日志目录: %s\n", logPathF.Get())
}
```

### EnumFlag 大小写敏感示例

展示枚举标志的大小写敏感设置：

```go
package main

import (
	"fmt"
	"os"

	"gitee.com/MM-Q/qflag"
)

func main() {
	// 创建枚举标志并启用大小写敏感
	modeF := qflag.Enum("mode", "m", "debug", []string{"Debug", "Release", "Test"}, "运行模式")
	modeF.SetCaseSensitive(true)

	if err := qflag.Parse(); err != nil {
		fmt.Printf("解析参数错误: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("运行模式: %s\n", modeF.Get())
}
```

注意：启用大小写敏感后，`--mode debug`会报错，必须使用`--mode Debug`

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

## 贡献

欢迎贡献代码或提出问题。请提交 PR 或 Issue 到 [Gitee 仓库](https://gitee.com/MM-Q/qflag)。
