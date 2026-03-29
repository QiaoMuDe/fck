# qflag 命令开发规范

本文档定义了使用 qflag 库开发 CLI 命令的标准规范。

## 适用范围

- 基于 qflag 库的 CLI 工具
- 命令式架构的 Go 应用程序

## 目录结构规范

推荐目录结构（可根据项目调整）：

```
<project>/
├── internal/                    # 或 pkg/，视项目而定
│   ├── commands/                # 业务逻辑层
│   │   └── <command_name>/      # 命令目录（小写）
│   │       └── cmd_<command_name>.go
│   └── cli/                     # CLI 定义层
│       ├── <command_name>.go
│       └── root.go              # 命令注册
└── cmd/
    └── main.go                  # 程序入口
```

**说明**：`internal/` 和 `pkg/` 选择取决于项目需求：
- `internal/`：私有代码，仅本项目可用（推荐）
- `pkg/`：公共代码，可被外部导入

## 命令命名规范

### 1. 目录和文件命名
- **命令目录**：小写，如 `mkdir`, `touch`, `rm`
- **业务逻辑文件**：`cmd_<command_name>.go`，如 `cmd_mkdir.go`
- **CLI 定义文件**：`<command_name>.go`，如 `mkdir.go`

### 2. 代码命名
- **包名**：与命令目录同名，小写，如 `package mkdir`
- **配置结构体**：`<Command>Config`，如 `MkdirConfig`
- **主函数**：`<Command>CmdMain`，如 `MkdirCmdMain`
- **命令变量**：`<Command>Cmd`，如 `MkdirCmd`
- **选项变量**：`<command><Option>`，如 `mkdirParents`, `mkdirMode`
- **运行函数**：`run<Command>`，如 `runMkdir`

## 业务逻辑文件规范

### 1. 文件结构模板

```go
package <command_name>

import (
	"fmt"
	"os"
	"path/filepath"
	// 其他必要的导入
)

// <Command>Config 配置结构体
type <Command>Config struct {
	Targets   []string
	Option1   bool
	Option2   int
	Option3   string
	Verbose   bool
}

// <Command>Stats 操作统计（可选）
type <Command>Stats struct {
	Processed int
	Errors    int
}

// <Command>CmdMain 主函数
func <Command>CmdMain(config <Command>Config) error {
	if len(config.Targets) == 0 {
		return fmt.Errorf("未指定要操作的目标")
	}

	stats := &<Command>Stats{}

	for _, target := range config.Targets {
		err := processTarget(target, config, stats)
		if err != nil {
			stats.Errors++
			return err
		}
	}

	if config.Verbose {
		fmt.Printf("操作完成: %d 个处理", stats.Processed)
		fmt.Println()
	}

	return nil
}

// processTarget 处理单个目标
func processTarget(path string, config <Command>Config, stats *<Command>Stats) error {
	// 处理逻辑
	return nil
}
```

### 2. 编码规范
- 使用函数级注释

**函数级注释示例：**
```go
// processTarget 处理单个目标
//
// 参数:
//   - path: 目标路径
//   - config: 命令配置
//   - stats: 统计信息
//
// 返回值:
//   - error: 处理过程中的错误
func processTarget(path string, config <Command>Config, stats *<Command>Stats) error {
	// 处理逻辑
	return nil
}
```

- 错误处理使用 `fmt.Errorf` 包装错误
- 提供友好的错误信息
- 使用 `%w` 包装底层错误（Go 1.13+）

## CLI 定义文件规范

### 1. 文件结构模板

```go
package cli

import (
	"fmt"

	"<module_name>/internal/commands/<command_name>"  // 替换为实际模块路径
	"gitee.com/MM-Q/qflag"
)

var <Command>Cmd *qflag.Cmd

var (
	<command>Option1 *qflag.BoolFlag   // 选项说明
	<command>Option2 *qflag.StringFlag // 选项说明
	<command>Verbose *qflag.BoolFlag   // 显示详细信息
)

func init() {
	<Command>Cmd = qflag.NewCmd("<command>", "<short>", qflag.ExitOnError)

	<command>Option1 = <Command>Cmd.Bool("option1", "o1", "选项说明", false)
	<command>Option2 = <Command>Cmd.String("option2", "o2", "选项说明", "default")
	<command>Verbose = <Command>Cmd.Bool("verbose", "v", "显示详细信息", false)

	cmdOpts := &qflag.CmdOpts{
		Desc: "命令描述",
		Notes: []string{
			"说明1",
			"说明2",
		},
		UseChinese: true,
	}

	if err := <Command>Cmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	<Command>Cmd.SetRun(run<Command>)
}

func run<Command>(cmd qflag.Command) error {
	config := <command_name>.<Command>Config{
		Targets:   cmd.Args(),
		Option1:  <command>Option1.Get(),
		Option2:  <command>Option2.Get(),
		Verbose:  <command>Verbose.Get(),
	}

	return <command_name>.<Command>CmdMain(config)
}
```

### 2. 标志定义规范

| 类型 | 方法 | 示例 |
|------|------|------|
| Bool | `Bool("long", "short", "说明", 默认值)` | `cmd.Bool("force", "f", "强制", false)` |
| String | `String("long", "short", "说明", 默认值)` | `cmd.String("output", "o", "输出", "")` |
| Int | `Int("long", "short", "说明", 默认值)` | `cmd.Int("count", "c", "数量", 0)` |
| Int64 | `Int64("long", "short", "说明", 默认值)` | `cmd.Int64("size", "s", "大小", 0)` |
| Uint | `Uint("long", "short", "说明", 默认值)` | `cmd.Uint("port", "p", "端口", 8080)` |
| Uint8 | `Uint8("long", "short", "说明", 默认值)` | `cmd.Uint8("level", "l", "级别", 1)` |
| Uint16 | `Uint16("long", "short", "说明", 默认值)` | `cmd.Uint16("code", "c", "代码", 100)` |
| Uint32 | `Uint32("long", "short", "说明", 默认值)` | `cmd.Uint32("id", "i", "ID", 0)` |
| Uint64 | `Uint64("long", "short", "说明", 默认值)` | `cmd.Uint64("total", "t", "总数", 0)` |
| Float64 | `Float64("long", "short", "说明", 默认值)` | `cmd.Float64("rate", "r", "比率", 0.0)` |
| Enum | `Enum("long", "short", "说明", 默认值, 可选值)` | `cmd.Enum("type", "t", "类型", "md5", []string{"md5", "sha1"})` |
| Duration | `Duration("long", "short", "说明", 默认值)` | `cmd.Duration("timeout", "t", "超时", 0)` |
| Time | `Time("long", "short", "说明", 默认值)` | `cmd.Time("start", "s", "开始时间", time.Time{})` |
| Size | `Size("long", "short", "说明", 默认值)` | `cmd.Size("limit", "l", "限制", 0)` |
| StringSlice | `StringSlice("long", "short", "说明", 默认值)` | `cmd.StringSlice("tags", "t", "标签", nil)` |
| IntSlice | `IntSlice("long", "short", "说明", 默认值)` | `cmd.IntSlice("ids", "i", "ID列表", nil)` |
| Int64Slice | `Int64Slice("long", "short", "说明", 默认值)` | `cmd.Int64Slice("values", "v", "值列表", nil)` |
| Map | `Map("long", "short", "说明", 默认值)` | `cmd.Map("env", "e", "环境变量", map[string]string{})` |

### 3. 选项配置规范

`CmdOpts` 支持以下配置项：

| 配置项 | 类型 | 说明 | 示例 |
|--------|------|------|------|
| `Desc` | string | 命令描述 | `"创建目录"` |
| `RunFunc` | `func(Command) error` | 命令执行函数 | `run` |
| `Version` | string | 版本号（仅在根命令生效） | `"1.0.0"` |
| `UseChinese` | bool | 使用中文帮助 | `true` |
| `EnvPrefix` | string | 环境变量前缀 | `"FCK"` |
| `UsageSyntax` | string | 命令使用语法（使用 `fmt.Sprintf` 替换 `%s`） | `fmt.Sprintf("%s 当前子命令名 [选项] [位置参数...]", qflag.Root.Name())` |
| `LogoText` | string | Logo文本 | `"FCK Tools"` |
| `Completion` | bool | 启用自动补全（仅在根命令生效） | `true` |
| `AutoBindEnv` | bool | 自动绑定所有标志的环境变量 | `true` |
| `Examples` | map[string]string | 使用示例 | `map[string]string{"创建单个目录": "mkdir test"}` |
| `Notes` | []string | 注意事项列表 | `[]string{"说明1", "说明2"}` |
| `SubCmds` | []Command | 子命令列表 | `[]qflag.Command{RunCmd}` |
| `MutexGroups` | []MutexGroup | 互斥组 | 定义互斥的标志 |
| `RequiredGroups` | []RequiredGroup | 必需组 | 定义必需的标志 |

**基础配置示例：**
```go
cmdOpts := &qflag.CmdOpts{
    Desc:       "命令描述",
    Version:    "1.0.0",
    UseChinese: true,
    Notes: []string{
        "说明1",
        "说明2",
    },
}
```

**完整配置示例：**
```go
cmdOpts := &qflag.CmdOpts{
    Desc:        "创建目录",
    RunFunc:     run,                    // 命令执行函数
    Version:     "1.0.0",                // 版本号（仅在根命令生效）
    UseChinese:  true,
    EnvPrefix:   "FCK",                  // 环境变量前缀
    UsageSyntax: fmt.Sprintf("%s 当前子命令名 [选项] [位置参数...]", qflag.Root.Name()),
    LogoText:    "FCK Tools",
    Completion:  true,                   // 启用自动补全（仅在根命令生效）
    AutoBindEnv: true,                   // 自动绑定所有标志的环境变量
    Examples: map[string]string{
        "创建单个目录":   "mkdir test",
        "递归创建目录":   "mkdir -p a/b/c",
    },
    Notes: []string{
        "支持递归创建",
        "支持设置权限",
    },
    SubCmds: []qflag.Command{
        BuildCmd,
        ConfigCmd,
    },
    MutexGroups: []types.MutexGroup{
        {
            Name:      "format",
            Flags:     []string{"json", "xml"},
            AllowNone: true,
        },
    },
    RequiredGroups: []types.RequiredGroup{
        {
            Name:        "auth",
            Flags:       []string{"user", "pass"},
            Conditional: true,
        },
    },
}
```

## 环境变量绑定规范

QFlag 提供了三种环境变量绑定方式，可根据实际需求选择。

### 1. 手动指定环境变量名

通过 `BindEnv()` 方法手动指定环境变量名称：

```go
func init() {
    Cmd = qflag.NewCmd("run", "r", qflag.ExitOnError)
    Cmd.SetEnvPrefix("MYAPP")  // 设置环境变量前缀
    
    // 手动绑定：绑定到 MYAPP_DATABASE_URL
    dbFlag := Cmd.String("database", "d", "数据库地址", "localhost")
    dbFlag.BindEnv("DATABASE_URL")
    
    // ...
}
```

### 2. 标志自动绑定

通过 `AutoBindEnv()` 方法自动使用标志长名称的大写形式作为环境变量名：

```go
func init() {
    Cmd = qflag.NewCmd("run", "r", qflag.ExitOnError)
    Cmd.SetEnvPrefix("MYAPP")
    
    // 自动绑定：host -> MYAPP_HOST, port -> MYAPP_PORT
    hostFlag := Cmd.String("host", "H", "主机地址", "localhost")
    portFlag := Cmd.Int("port", "p", "端口号", 8080)
    
    hostFlag.AutoBindEnv()
    portFlag.AutoBindEnv()
    
    // ...
}
```

### 3. 命令批量自动绑定

通过 `AutoBindAllEnv()` 方法一次性为命令的所有标志自动绑定环境变量：

```go
func init() {
    Cmd = qflag.NewCmd("run", "r", qflag.ExitOnError)
    Cmd.SetEnvPrefix("MYAPP")
    
    // 创建多个标志
    Cmd.String("host", "H", "主机地址", "localhost")
    Cmd.Int("port", "p", "端口号", 8080)
    Cmd.String("user", "u", "用户名", "admin")
    
    // 批量自动绑定所有标志
    Cmd.AutoBindAllEnv()
    
    // ...
}
```

### 4. 通过 CmdOpts 配置自动绑定

在 `CmdOpts` 中设置 `AutoBindEnv` 字段：

```go
func init() {
    Cmd = qflag.NewCmd("run", "r", qflag.ExitOnError)
    
    // 创建标志
    Cmd.String("host", "H", "主机地址", "localhost")
    Cmd.Int("port", "p", "端口号", 8080)
    
    cmdOpts := &qflag.CmdOpts{
        Desc:        "运行服务",
        EnvPrefix:   "MYAPP",
        AutoBindEnv: true,  // 自动绑定所有标志的环境变量
        UseChinese:  true,
    }
    
    if err := Cmd.ApplyOpts(cmdOpts); err != nil {
        panic(fmt.Errorf("apply opts err: %w", err))
    }
    
    Cmd.SetRun(run)
}
```

### 5. 三种方式对比

| 方式 | 方法 | 适用场景 | 特点 |
|------|------|----------|------|
| 手动指定 | `BindEnv("NAME")` | 需要自定义环境变量名 | 灵活，可指定任意名称 |
| 标志自动绑定 | `AutoBindEnv()` | 单个标志自动绑定 | 使用长名称大写，简洁 |
| 命令批量绑定 | `AutoBindAllEnv()` | 批量绑定所有标志 | 一次性绑定，高效 |
| CmdOpts 配置 | `AutoBindEnv: true` | 配置化管理 | 与其他配置一起设置 |

### 6. 环境变量绑定注意事项

1. **前缀设置**：使用 `SetEnvPrefix()` 或 `CmdOpts.EnvPrefix` 设置环境变量前缀
2. **命名规则**：环境变量名 = 前缀 + _ + 标志名（大写）
3. **优先级**：命令行参数 > 环境变量 > 默认值
4. **长名称要求**：`AutoBindEnv()` 和 `AutoBindAllEnv()` 要求标志必须有长名称，否则会 panic

## 命令注册规范

### 1. 注册位置
在 `internal/cli/root.go` 的 `SubCmds` 列表中添加命令：

```go
SubCmds: []qflag.Command{
	MkdirCmd,
	TouchCmd,
	<Command>Cmd,  // 新增命令
},
```

### 2. 排序规则

按功能逻辑排序，建议：
1. 核心/常用命令在前
2. 相关功能命令放在一起
3. 新命令插入到合适的功能组

## 开发流程

### 1. 创建业务逻辑文件
```bash
# 创建命令目录
mkdir internal/commands/<command_name>

# 创建业务逻辑文件
touch internal/commands/<command_name>/cmd_<command_name>.go
```

### 2. 编写业务逻辑
- 定义配置结构体
- 实现主函数
- 实现辅助函数
- 添加错误处理

### 3. 创建 CLI 定义文件
```bash
# 创建 CLI 定义文件
touch internal/cli/<command_name>.go
```

### 4. 编写 CLI 定义
- 定义命令变量
- 定义选项变量
- 实现 init() 函数
- 实现运行函数

### 5. 注册命令
在 `internal/cli/root.go` 中注册命令

### 6. 编译验证
```bash
# 编译整个项目
go build ./...

# 或编译特定包
go build ./internal/cli/...
```

## 代码风格规范

### 1. 注释规范
- 使用函数级注释
- 注释使用中文
- 注释简洁明了

**函数级注释示例：**
```go
// <Command>CmdMain 执行<命令名称>命令
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - error: 执行错误
func <Command>CmdMain(config <Command>Config) error {
	// 实现逻辑
	return nil
}

// processTarget 处理单个目标
//
// 参数:
//   - path: 目标路径
//   - config: 命令配置
//   - stats: 统计信息
//
// 返回值:
//   - error: 处理过程中的错误
func processTarget(path string, config <Command>Config, stats *<Command>Stats) error {
	// 处理逻辑
	return nil
}
```

### 2. 错误处理

- 使用 `fmt.Errorf` 包装错误
- 提供友好的错误信息
- 使用 `%w` 包装底层错误（Go 1.13+）

**常见错误处理模式：**

```go
// 参数验证错误
if len(config.Targets) == 0 {
    return fmt.Errorf("未指定目标")
}

// 文件操作错误
if err != nil {
    return fmt.Errorf("读取文件失败: %w", err)
}

// 命令执行错误
if exitCode != 0 {
    return fmt.Errorf("命令执行失败，退出码: %d", exitCode)
}
```

### 3. 变量命名
- 导出变量：首字母大写（如 `MkdirCmd`）
- 私有变量：首字母小写（如 `mkdirParents`）
- 常量：全大写（如 `MAX_DEPTH`）

### 4. 函数命名
- 导出函数：首字母大写（如 `MkdirCmdMain`）
- 私有函数：首字母小写（如 `createDirectory`）

## 完整示例

参考示例实现：
- `internal/commands/example/cmd_example.go` - 业务逻辑示例
- `internal/cli/example.go` - CLI 定义示例

**提示**：可在本文档所在项目的 `internal/commands/` 目录中找到实际示例。

## 注意事项

1. **包导入**：使用绝对路径导入，如 `"<module_name>/internal/commands/<command_name>"`
2. **错误处理**：所有可能的错误都要处理，提供友好的错误信息
3. **参数验证**：验证所有输入参数的有效性
4. **跨平台**：确保代码在 Windows、Linux、macOS 上都能运行
5. **性能考虑**：注意性能优化，避免不必要的资源消耗
6. **测试验证**：开发完成后进行编译验证和功能测试

## 快速检查清单

在提交新命令前，检查以下项目：

- [ ] 目录结构符合规范
- [ ] 命名符合规范
- [ ] 业务逻辑文件已创建
- [ ] CLI 定义文件已创建
- [ ] 命令已注册到 root.go
- [ ] 编译通过
- [ ] 帮助文档完整
- [ ] 错误处理完善
- [ ] 代码风格一致
- [ ] 功能测试通过

## 附录

### qflag 库信息

- **仓库**: `gitee.com/MM-Q/qflag`
- **文档**: 参考 qflag 官方文档
- **版本要求**: v0.5.10+

### 相关文档

- [qflag 使用指南](https://gitee.com/MM-Q/qflag)
- [Go 命令行工具最佳实践](https://golang.org/doc/code.html)
