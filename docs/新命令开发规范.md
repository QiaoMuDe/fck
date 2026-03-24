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

### 2. 选项定义规范

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

### 3. 帮助文档规范

`CmdOpts` 支持以下配置项：

| 配置项 | 类型 | 说明 | 示例 |
|--------|------|------|------|
| `Desc` | string | 命令描述 | `"创建目录"` |
| `Version` | string | 版本号（仅在根命令生效） | `"1.0.0"` |
| `UseChinese` | bool | 使用中文帮助 | `true` |
| `EnvPrefix` | string | 环境变量前缀 | `"FCK"` |
| `UsageSyntax` | string | 命令使用语法（使用 `fmt.Sprintf` 替换 `%s`） | `fmt.Sprintf("%s 当前子命令名 [选项] [位置参数...]", qflag.Root.Name())` |
| `LogoText` | string | Logo文本 | `"FCK Tools"` |
| `Completion` | bool | 启用自动补全（仅在根命令生效） | `true` |
| `Notes` | []string | 注意事项列表 | `[]string{"说明1", "说明2"}` |
| `Examples` | map[string]string | 使用示例 | `map[string]string{"创建单个目录": "mkdir test"}` |
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
    Version:     "1.0.0", // 版本号（仅在根命令生效）
    UseChinese:  true, 
    EnvPrefix:   "FCK",
    UsageSyntax: fmt.Sprintf("%s 当前子命令名 [选项] [位置参数...]", qflag.Root.Name()),
    LogoText:    "FCK Tools",
    Completion:  true, // 启用自动补全（仅在根命令生效）
    Notes: []string{
        "支持递归创建",
        "支持设置权限",
    },
    Examples: map[string]string{
        "创建单个目录":   "mkdir test",
        "递归创建目录":   "mkdir -p a/b/c",
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
- **版本要求**: v0.5.9+

### 相关文档

- [qflag 使用指南](https://gitee.com/MM-Q/qflag)
- [Go 命令行工具最佳实践](https://golang.org/doc/code.html)
