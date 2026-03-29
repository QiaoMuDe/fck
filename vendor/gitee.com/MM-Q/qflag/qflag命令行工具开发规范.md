# qflag 命令行工具开发规范

本文档定义了使用 qflag 库开发命令行工具的标准规范，强调**简洁、统一、易维护**的目录结构。

## 目录结构规范

### 核心原则

1. **入口统一**: 程序入口固定在 `cmd/` 目录
2. **命令集中**: 所有命令定义统一放在 `internal/cli/` 目录
3. **单文件原则**: 每个命令 = 单个文件（包含：初始化 + 标志 + run 函数）
4. **按需创建**: 只有两级子命令才创建同名子目录

### 标准目录结构

```
your-project/
├── cmd/                          # 程序入口（固定）
│   └── yourapp/                  # 你的工具名
│       └── main.go               # 唯一入口：只调用 cli.Run()
├── internal/
│   ├── cli/                      # 所有命令定义（核心目录）
│   │   ├── root.go               # 根命令：help/usage
│   │   ├── run.go                # 一级子命令：yourapp run
│   │   ├── build.go              # 一级子命令：yourapp build
│   │   ├── version.go            # 一级子命令：yourapp version
│   │   ├── config.go             # 一级命令：yourapp config
│   │   └── config/               # 二级子命令目录（同名文件夹）
│   │       ├── get.go            # yourapp config get
│   │       └── set.go            # yourapp config set
│   ├── utils/                    # 工具函数
│   ├── config/                   # 配置相关
│   └── service/                  # 业务服务
├── go.mod
└── README.md
```

### 目录结构说明

| 目录/文件 | 说明 | 必需 |
|----------|------|------|
| `cmd/yourapp/main.go` | 程序入口，调用 `cli.Run()` | ✅ |
| `internal/cli/root.go` | 根命令定义和配置 | ✅ |
| `internal/cli/*.go` | 一级子命令定义（单文件） | 按需 |
| `internal/cli/*/` | 二级子命令目录 | 按需 |
| `internal/utils/` | 工具函数 | 按需 |
| `internal/config/` | 配置相关 | 按需 |
| `internal/service/` | 业务服务 | 按需 |

---

## 命令文件规范

### 统一文件结构

每个命令文件必须遵循以下结构：

```go
package cli

import (
    "fmt"
    "gitee.com/MM-Q/qflag"
)

// ============================================
// 1. 全局命令变量（供注册到根命令）
// ============================================
var (
    // Cmd 命令对象（必须导出，供 root.go 注册）
    Cmd *qflag.Cmd
)

// ============================================
// 2. 全局标志变量（用于在 run 函数中传递）
// ============================================
var (
    // 示例标志
    flagName    *qflag.StringFlag
    flagCount   *qflag.IntFlag
    flagVerbose *qflag.BoolFlag
)

// ============================================
// 3. init() 初始化命令、定义标志
// ============================================
func init() {
    // 初始化命令
    Cmd = qflag.NewCmd("command-name", "c", qflag.ExitOnError)
    
    // 定义标志
    flagName = Cmd.String("name", "n", "名称描述", "default")
    flagCount = Cmd.Int("count", "c", "计数", 1)
    flagVerbose = Cmd.Bool("verbose", "v", "详细输出", false)
    
    // 应用命令配置
    cmdOpts := &qflag.CmdOpts{
        Desc:        "命令描述",
        UsageSyntax: fmt.Sprintf("%s command-name [option] [args]", qflag.Root.Name()),
        UseChinese:  true,
    }
    
    if err := Cmd.ApplyOpts(cmdOpts); err != nil {
        panic(fmt.Errorf("apply opts err: %w", err))
    }
    
    // 设置运行函数
    Cmd.SetRun(run)
}

// ============================================
// 4. run() 函数：业务逻辑，运行逻辑
// ============================================
func run(cmd qflag.Command) error {
    // 获取标志值
    name := flagName.Get()
    count := flagCount.Get()
    verbose := flagVerbose.Get()
    
    // 执行业务逻辑
    // ...
    
    return nil
}
```

### 文件结构要点

1. **包声明**: 所有命令文件统一使用 `package cli`
2. **命令变量**: 必须导出（首字母大写），命名为 `Cmd`
3. **标志变量**: 使用 `flag` 前缀命名，私有变量
4. **init() 函数**: 负责初始化命令、定义标志、设置配置
5. **run() 函数**: 负责业务逻辑执行，返回 error

---

## 根命令定义规范

### 核心区别

**根命令 ≠ 子命令**：
- 子命令：可以在 `init()` 中创建，导出命令对象供上级注册
- 根命令：没有上级，需要使用 `InitAndRun()` 模式直接操作 `qflag.Root`

### 文件位置

`internal/cli/root.go`

### 完整示例

```go
package cli

import (
    "fmt"

    "gitee.com/MM-Q/qflag"
    "gitee.com/MM-Q/verman"
    "your-project/internal/types"
)

// ============================================
// 1. 全局标志变量（定义在根命令上）
// ============================================
var (
    listFlag     *qflag.BoolFlag
    runFlag      *qflag.StringFlag
    forceFlag    *qflag.BoolFlag
    filePathFlag *qflag.StringFlag
)

// ============================================
// 2. InitAndRun 初始化并运行根命令
// ============================================
// InitAndRun 初始化并运行根命令
//
// 返回值:
//   - err: 初始化或运行命令时可能发生的错误
func InitAndRun() (err error) {
    // defer 捕获 panic
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic recovered: %v", r)
        }
    }()

    // 注册根命令的标志（直接在 qflag.Root 上定义）
    listFlag = qflag.Root.Bool("list", "l", "列出所有可用任务", false)
    runFlag = qflag.Root.String("run", "r", "运行指定任务", "")
    forceFlag = qflag.Root.Bool("force", "f", "执行强制类操作", false)
    filePathFlag = qflag.Root.String("path", "p", "指定任务配置文件路径", "")

    // 配置根命令
    rootCmdOpts := &qflag.CmdOpts{
        Version:    verman.V.Version(),
        Desc:       "yourapp - 你的应用描述",
        UseChinese: true,
        Completion: true,
        LogoText:   types.LogoText,
        RunFunc:    run,
        Examples: map[string]string{
            "列出所有任务":     fmt.Sprintf("%s -l", qflag.Root.Name()),
            "运行指定任务":     fmt.Sprintf("%s -r deploy", qflag.Root.Name()),
            "指定配置文件运行": fmt.Sprintf("%s -r deploy -p custom.toml", qflag.Root.Name()),
            "强制运行任务":     fmt.Sprintf("%s -r deploy -f", qflag.Root.Name()),
        },
        Notes: []string{
            "默认查找的配置文件: config.toml",
            "未指定配置文件时，按优先级查找当前目录下的配置文件",
        },
        MutexGroups: []qflag.MutexGroup{
            {
                Name:      "run-or-list",
                Flags:     []string{"run", "list"},
                AllowNone: true,
            },
        },
        SubCmds: []qflag.Command{
            RunCmd,     // yourapp run
            BuildCmd,   // yourapp build
            VersionCmd, // yourapp version
            ConfigCmd,  // yourapp config
        },
    }

    // 应用根命令配置
    if err = qflag.ApplyOpts(rootCmdOpts); err != nil {
        err = fmt.Errorf("apply opts failed: %w", err)
        return err
    }

    // 解析并自动路由到子命令
    if err = qflag.ParseAndRoute(); err != nil {
        err = fmt.Errorf("parse and route failed: %w", err)
        return err
    }

    return nil
}

// ============================================
// 3. run() 函数：根命令的业务逻辑
// ============================================
// run 是根命令的执行函数
//
// 参数:
//   - cmd: 根命令接口
//
// 返回值:
//   - error: 执行时可能遇到的错误
func run(cmd qflag.Command) error {
    // 获取标志值
    list := listFlag.Get()
    runTask := runFlag.Get()
    force := forceFlag.Get()
    filePath := filePathFlag.Get()

    // 根据标志执行不同逻辑
    if list {
        // 执行列表逻辑
        return listTasks()
    }

    if runTask != "" {
        // 执行运行任务逻辑
        return runTask(runTask, filePath, force)
    }

    // 默认显示帮助
    cmd.PrintHelp()
    return nil
}

// listTasks 列出所有任务
func listTasks() error {
    // 实现列出任务逻辑
    return nil
}

// runTask 运行指定任务
func runTask(taskName, filePath string, force bool) error {
    // 实现运行任务逻辑
    return nil
}
```

### 根命令要点

1. **InitAndRun() 函数**: 统一入口函数，由 `main.go` 调用
2. **直接操作 qflag.Root**: 在 `InitAndRun()` 中直接定义标志
3. **defer 捕获 panic**: 保证 panic 转换为 error 返回
4. **命名返回值**: 使用 `err` 作为命名返回值，方便 defer 和错误处理
5. **SubCmds 注册**: 所有一级子命令在这里注册
6. **自动路由**: 使用 `qflag.ParseAndRoute()` 自动路由到子命令

---

## 一级子命令规范

### 文件位置

`internal/cli/command.go`

### 完整示例

```go
package cli

import (
    "fmt"

    "gitee.com/MM-Q/qflag"
)

// ============================================
// 1. 全局命令变量
// ============================================
var (
    // RunCmd yourapp run 命令
    RunCmd *qflag.Cmd
)

// ============================================
// 2. 全局标志变量
// ============================================
var (
    runInput    *qflag.StringFlag
    runOutput   *qflag.StringFlag
    runParallel *qflag.BoolFlag
)

// ============================================
// 3. init() 初始化
// ============================================
func init() {
    // 初始化命令
    RunCmd = qflag.NewCmd("run", "r", qflag.ExitOnError)
    
    // 定义标志
    runInput = RunCmd.String("input", "i", "输入文件路径", "")
    runOutput = RunCmd.String("output", "o", "输出文件路径", "")
    runParallel = RunCmd.Bool("parallel", "p", "并行执行", false)
    
    // 应用命令配置
    cmdOpts := &qflag.CmdOpts{
        Desc:        "运行任务",
        UsageSyntax: fmt.Sprintf("%s run [option] [args]", qflag.Root.Name()),
        UseChinese:  true,
    }
    
    if err := RunCmd.ApplyOpts(cmdOpts); err != nil {
        panic(fmt.Errorf("apply opts err: %w", err))
    }
    
    // 设置运行函数
    RunCmd.SetRun(runRun)
}

// ============================================
// 4. run() 函数
// ============================================
func runRun(cmd qflag.Command) error {
    // 获取标志值
    input := runInput.Get()
    output := runOutput.Get()
    parallel := runParallel.Get()
    
    // 执行业务逻辑
    fmt.Printf("运行任务: input=%s, output=%s, parallel=%v\n", input, output, parallel)
    
    return nil
}
```

### 命名规范

- **命令变量**: `XxxCmd` 格式（如 `RunCmd`、`BuildCmd`）
- **标志变量**: `命令前缀 + 标志名` 格式（如 `runInput`、`runOutput`）
- **run 函数**: `run + 命令名` 格式（如 `runRun`、`runBuild`）

---

## 二级子命令规范

### 目录结构

```
internal/cli/
├── config.go           # 一级命令: yourapp config
└── config/             # 二级命令目录
    ├── get.go          # yourapp config get
    └── set.go          # yourapp config set
```

### 一级命令示例（无业务逻辑）

`internal/cli/config.go`:

```go
package cli

import (
    "fmt"

    "gitee.com/MM-Q/qflag"
)

var (
    // ConfigCmd yourapp config 命令
    ConfigCmd *qflag.Cmd
)

func init() {
    ConfigCmd = qflag.NewCmd("config", "c", qflag.ExitOnError)
    
    // 注册二级子命令
    cmdOpts := &qflag.CmdOpts{
        Desc:        "配置管理",
        UsageSyntax: fmt.Sprintf("%s config [command]", qflag.Root.Name()),
        UseChinese:  true,
        SubCmds: []qflag.Command{
            ConfigGetCmd, // yourapp config get
            ConfigSetCmd, // yourapp config set
        },
    }
    
    if err := ConfigCmd.ApplyOpts(cmdOpts); err != nil {
        panic(fmt.Errorf("apply opts err: %w", err))
    }
    
    ConfigCmd.SetRun(runConfig)
}

func runConfig(cmd qflag.Command) error {
    cmd.PrintHelp()
    return nil
}
```

### 二级命令示例

`internal/cli/config/get.go`:

```go
package config

import (
    "fmt"

    "gitee.com/MM-Q/qflag"
)

var (
    // ConfigGetCmd yourapp config get 命令
    ConfigGetCmd *qflag.Cmd
)

var (
    getKey    *qflag.StringFlag
    getGlobal *qflag.BoolFlag
)

func init() {
    ConfigGetCmd = qflag.NewCmd("get", "g", qflag.ExitOnError)
    
    // 定义标志
    getKey = ConfigGetCmd.String("key", "k", "配置键名", "")
    getGlobal = ConfigGetCmd.Bool("global", "G", "全局配置", false)
    
    // 应用命令配置
    cmdOpts := &qflag.CmdOpts{
        Desc:        "获取配置项",
        UsageSyntax: fmt.Sprintf("%s config get [选项]", qflag.Root.Name()),
        UseChinese:  true,
    }
    
    if err := ConfigGetCmd.ApplyOpts(cmdOpts); err != nil {
        panic(fmt.Errorf("apply opts err: %w", err))
    }
    
    ConfigGetCmd.SetRun(runConfigGet)
}

func runConfigGet(cmd qflag.Command) error {
    key := getKey.Get()
    global := getGlobal.Get()
    
    fmt.Printf("获取配置: key=%s, global=%v\n", key, global)
    
    return nil
}
```

`internal/cli/config/set.go`:

```go
package config

import (
    "fmt"

    "gitee.com/MM-Q/qflag"
)

var (
    // ConfigSetCmd yourapp config set 命令
    ConfigSetCmd *qflag.Cmd
)

var (
    setKey    *qflag.StringFlag
    setValue  *qflag.StringFlag
    setGlobal *qflag.BoolFlag
)

func init() {
    ConfigSetCmd = qflag.NewCmd("set", "s", qflag.ExitOnError)
    
    setKey = ConfigSetCmd.String("key", "k", "配置键名", "")
    setValue = ConfigSetCmd.String("value", "v", "配置值", "")
    setGlobal = ConfigSetCmd.Bool("global", "G", "全局配置", false)
    
    cmdOpts := &qflag.CmdOpts{
        Desc:        "设置配置项",
        UsageSyntax: fmt.Sprintf("%s config set [选项]", qflag.Root.Name()),
        UseChinese:  true,
    }
    
    if err := ConfigSetCmd.ApplyOpts(cmdOpts); err != nil {
        panic(fmt.Errorf("apply opts err: %w", err))
    }
    
    ConfigSetCmd.SetRun(runConfigSet)
}

func runConfigSet(cmd qflag.Command) error {
    key := setKey.Get()
    value := setValue.Get()
    global := setGlobal.Get()
    
    fmt.Printf("设置配置: key=%s, value=%s, global=%v\n", key, value, global)
    
    return nil
}
```

### 二级命令要点

1. **包名**: 使用命令名作为包名（如 `package config`）
2. **UsageSyntax**: 使用 `qflag.Root.Name()` 拼接完整命令路径（如 `yourapp config get [选项]`）
3. **注册位置**: 在一级命令的 `SubCmds` 中注册
4. **命名规范**: `一级命令 + 二级命令 + Cmd`（如 `ConfigGetCmd`）
5. **无需导入父包**: 直接使用 `qflag.Root.Name()` 即可

---

## 程序入口规范

### 文件位置

`cmd/yourapp/main.go`

### 完整示例

```go
package main

import (
    "os"

    "your-project/internal/cli"
    "your-project/internal/utils"
)

func main() {
    // 调用 cli.InitAndRun() 初始化并运行
    if err := cli.InitAndRun(); err != nil {
        utils.LogErr(err.Error())
        os.Exit(1)
    }
}
```

### 入口文件要点

1. **极简原则**: 只负责调用 `cli.InitAndRun()`
2. **错误处理**: 统一使用 `utils.LogErr()` 输出错误
3. **退出码**: 错误时使用 `os.Exit(1)`

---

## 标志类型规范

### 支持的标志类型

| 类型 | 方法 | 示例 |
|------|------|------|
| 字符串 | `Cmd.String()` | `Cmd.String("name", "n", "描述", "default")` |
| 布尔 | `Cmd.Bool()` | `Cmd.Bool("verbose", "v", "描述", false)` |
| 整数 | `Cmd.Int()` | `Cmd.Int("count", "c", "描述", 1)` |
| 枚举 | `Cmd.Enum()` | `Cmd.Enum("type", "t", "描述", "opt1", []string{"opt1", "opt2"})` |
| 时长 | `Cmd.Duration()` | `Cmd.Duration("timeout", "t", "描述", time.Second*10)` |
| 时间 | `Cmd.Time()` | `Cmd.Time("time", "t", "描述", time.Now())` |
| 大小 | `Cmd.Size()` | `Cmd.Size("size", "s", "描述", 1024)` |
| 字符串切片 | `Cmd.StringSlice()` | `Cmd.StringSlice("tags", "t", "描述", []string{"tag1"})` |
| 整数切片 | `Cmd.IntSlice()` | `Cmd.IntSlice("ports", "p", "描述", []int{8080})` |
| 映射 | `Cmd.Map()` | `Cmd.Map("labels", "l", "描述", map[string]string{"key": "val"})` |

### 标志命名规范

1. **长名称**: 使用小写字母，多个单词用连字符连接（如 `--output-file`）
2. **短名称**: 使用单个小写字母（如 `-o`）
3. **变量名**: 使用驼峰命名，加命令前缀（如 `runOutput`、`buildOutputFile`）

---

## 互斥组和必需组

### 互斥组（MutexGroup）

确保组内**最多只有一个**标志被设置。

```go
cmdOpts := &qflag.CmdOpts{
    Desc: "部署命令",
    MutexGroups: []qflag.MutexGroup{
        {
            Name:      "deploy-target",      // 组名
            Flags:     []string{"dev", "staging", "prod"},  // 互斥标志
            AllowNone: true,                 // 是否允许一个都不设置
        },
    },
}
```

**使用场景**：
```bash
# ✅ 正确：设置一个标志
yourapp deploy --dev

# ✅ 正确：都不设置（AllowNone=true）
yourapp deploy

# ❌ 错误：设置多个标志
yourapp deploy --dev --staging  # 报错：互斥标志冲突
```

### 必需组（RequiredGroup）

确保组内**所有**标志都被设置。

```go
cmdOpts := &qflag.CmdOpts{
    Desc: "数据库配置",
    RequiredGroups: []qflag.RequiredGroup{
        {
            Name:        "database",         // 组名
            Flags:       []string{"db-host", "db-port", "db-name"},  // 必需标志
            Conditional: false,              // 是否为条件性必需
        },
    },
}
```

**条件性必需组**：
```go
RequiredGroups: []qflag.RequiredGroup{
    {
        Name:        "auth",
        Flags:       []string{"username", "password"},
        Conditional: true,  // 只有设置了其中一个，才要求设置全部
    },
}
```

**使用场景**：
```bash
# 普通必需组（Conditional=false）
# ❌ 错误：缺少 db-name
yourapp migrate --db-host=localhost --db-port=3306

# ✅ 正确：设置所有标志
yourapp migrate --db-host=localhost --db-port=3306 --db-name=mydb

# 条件性必需组（Conditional=true）
# ✅ 正确：都不设置
yourapp run

# ❌ 错误：设置了一个但未设置全部
yourapp run --username=admin  # 缺少 password

# ✅ 正确：设置全部
yourapp run --username=admin --password=secret
```

---

## 环境变量绑定

QFlag 提供了三种环境变量绑定方式，可根据实际需求选择。

### 方式一：手动指定环境变量名

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

### 方式二：标志自动绑定

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

### 方式三：命令批量自动绑定

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

### 方式四：通过 CmdOpts 配置

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

### 四种方式对比

| 方式 | 方法 | 适用场景 | 特点 |
|------|------|----------|------|
| 手动指定 | `BindEnv("NAME")` | 需要自定义环境变量名 | 灵活，可指定任意名称 |
| 标志自动绑定 | `AutoBindEnv()` | 单个标志自动绑定 | 使用长名称大写，简洁 |
| 命令批量绑定 | `AutoBindAllEnv()` | 批量绑定所有标志 | 一次性绑定，高效 |
| CmdOpts 配置 | `AutoBindEnv: true` | 配置化管理 | 与其他配置一起设置 |

### 环境变量绑定注意事项

1. **前缀设置**：使用 `SetEnvPrefix()` 或 `CmdOpts.EnvPrefix` 设置环境变量前缀
2. **命名规则**：环境变量名 = 前缀 + _ + 标志名（大写）
3. **优先级**：命令行参数 > 环境变量 > 默认值
4. **长名称要求**：`AutoBindEnv()` 和 `AutoBindAllEnv()` 要求标志必须有长名称，否则会 panic

---

## 自动补全功能

### 启用自动补全

在根命令配置中启用：

```go
rootCmdOpts := &qflag.CmdOpts{
    Completion: true,  // 启用自动补全
    // ...
}
```

### 生成补全脚本

qflag 会自动注册 `--completion` 标志：

```bash
# 生成 Bash 补全脚本
yourapp --completion bash > /etc/bash_completion.d/yourapp

# 生成 PowerShell 补全脚本
yourapp --completion pwsh > yourapp-completion.ps1

# 生成 Fish 补全脚本
yourapp --completion fish > ~/.config/fish/completions/yourapp.fish
```

### 安装补全脚本

**Bash**:
```bash
# 方式1：系统级安装
yourapp --completion bash | sudo tee /etc/bash_completion.d/yourapp
source ~/.bashrc

# 方式2：用户级安装
yourapp --completion bash >> ~/.bash_completion
source ~/.bashrc
```

**PowerShell**:
```powershell
# 生成脚本
yourapp --completion pwsh > $HOME\yourapp-completion.ps1

# 在 $PROFILE 中添加
echo ". $HOME\yourapp-completion.ps1" >> $PROFILE
. $PROFILE
```

### 补全效果

```bash
# 自动补全命令
yourapp [TAB]
build    config   init     run      version

# 自动补全标志
yourapp run --[TAB]
--config    --help      --output    --parallel  --verbose

# 自动补全标志值（枚举类型）
yourapp build --target [TAB]
dev      staging  prod
```

---

## 测试规范

### 测试文件结构

```
your-project/
├── internal/cli/
│   ├── run.go
│   └── run_test.go      # 测试文件与源文件同目录
```

### 单元测试示例

```go
package cli

import (
    "testing"
    
    "gitee.com/MM-Q/qflag"
)

func TestRunCommand(t *testing.T) {
    // 创建测试命令
    testCmd := qflag.NewCmd("test", "t", qflag.ContinueOnError)
    
    // 定义标志
    input := testCmd.String("input", "i", "输入文件", "")
    output := testCmd.String("output", "o", "输出文件", "")
    
    // 应用配置
    cmdOpts := &qflag.CmdOpts{
        Desc:       "测试命令",
        UseChinese: true,
    }
    
    if err := testCmd.ApplyOpts(cmdOpts); err != nil {
        t.Fatalf("应用配置失败: %v", err)
    }
    
    // 测试参数解析
    args := []string{"--input=test.txt", "--output=out.txt"}
    if err := testCmd.Parse(args); err != nil {
        t.Fatalf("解析失败: %v", err)
    }
    
    // 验证标志值
    if input.Get() != "test.txt" {
        t.Errorf("期望 input=test.txt, 实际=%s", input.Get())
    }
    
    if output.Get() != "out.txt" {
        t.Errorf("期望 output=out.txt, 实际=%s", output.Get())
    }
}
```

### 测试错误场景

```go
func TestMutexGroupError(t *testing.T) {
    testCmd := qflag.NewCmd("test", "t", qflag.ContinueOnError)
    
    devFlag := testCmd.Bool("dev", "d", "开发环境", false)
    prodFlag := testCmd.Bool("prod", "p", "生产环境", false)
    
    cmdOpts := &qflag.CmdOpts{
        Desc: "测试互斥组",
        MutexGroups: []qflag.MutexGroup{
            {
                Name:      "env",
                Flags:     []string{"dev", "prod"},
                AllowNone: true,
            },
        },
    }
    
    if err := testCmd.ApplyOpts(cmdOpts); err != nil {
        t.Fatalf("应用配置失败: %v", err)
    }
    
    // 测试互斥标志冲突
    args := []string{"--dev", "--prod"}
    err := testCmd.Parse(args)
    
    if err == nil {
        t.Error("期望返回互斥标志冲突错误，但解析成功")
    }
}
```

### 测试覆盖率

```bash
# 运行测试并生成覆盖率报告
go test -cover ./internal/cli/...

# 生成详细覆盖率报告
go test -coverprofile=coverage.out ./internal/cli/...
go tool cover -html=coverage.out
```

---

## 最佳实践

### 1. 全局配置管理

```go
// internal/config/global.go
package config

import "sync"

var (
    once     sync.Once
    instance *GlobalConfig
)

// GlobalConfig 全局配置
type GlobalConfig struct {
    Verbose bool
    Debug   bool
    ConfigPath string
}

// GetGlobalConfig 获取全局配置单例
func GetGlobalConfig() *GlobalConfig {
    once.Do(func() {
        instance = &GlobalConfig{}
    })
    return instance
}

// internal/cli/root.go
var (
    verboseFlag    *qflag.BoolFlag
    debugFlag      *qflag.BoolFlag
    configPathFlag *qflag.StringFlag
)

func InitAndRun() (err error) {
    // ... 定义标志
    
    // 在 run 函数中初始化全局配置
    // run 函数参考前面的示例
}

func initGlobalConfig() {
    cfg := config.GetGlobalConfig()
    cfg.Verbose = verboseFlag.Get()
    cfg.Debug = debugFlag.Get()
    cfg.ConfigPath = configPathFlag.Get()
}
```

### 2. 配置文件处理

```go
// internal/config/loader.go
package config

import (
    "fmt"
    "os"
    
    "github.com/pelletier/go-toml/v2"
)

// LoadFromFile 从文件加载配置
func LoadFromFile(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("读取配置文件失败: %w", err)
    }
    
    var cfg Config
    if err := toml.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("解析配置文件失败: %w", err)
    }
    
    return &cfg, nil
}

// LoadWithPriority 按优先级加载配置
// 优先级：指定路径 > 当前目录默认文件 > 默认配置
func LoadWithPriority(specifiedPath string, defaultNames []string) (*Config, error) {
    // 1. 使用指定路径
    if specifiedPath != "" {
        return LoadFromFile(specifiedPath)
    }
    
    // 2. 查找默认文件
    for _, name := range defaultNames {
        if _, err := os.Stat(name); err == nil {
            return LoadFromFile(name)
        }
    }
    
    // 3. 返回默认配置
    return &Config{}, nil
}
```

### 3. 标志验证器

```go
import (
    "gitee.com/MM-Q/qflag"
    "gitee.com/MM-Q/qflag/validators"
)

// 使用内置验证器
portFlag := Cmd.Int("port", "p", "端口", 8080).
    SetValidator(validators.Range(1, 65535))

emailFlag := Cmd.String("email", "e", "邮箱", "").
    SetValidator(validators.Email())

urlFlag := Cmd.String("url", "u", "URL", "").
    SetValidator(validators.URL())

// 自定义验证器
pathFlag := Cmd.String("path", "p", "文件路径", "").
    SetValidator(func(value string) error {
        if value == "" {
            return nil  // 空值不验证
        }
        if _, err := os.Stat(value); os.IsNotExist(err) {
            return fmt.Errorf("文件不存在: %s", value)
        }
        return nil
    })
```

### 4. 优雅的错误处理

```go
// internal/cli/root.go
func run(cmd qflag.Command) error {
    // 检查必需条件
    if runTask := runFlag.Get(); runTask != "" {
        // 加载配置文件
        cfg, err := loadConfig(filePathFlag.Get())
        if err != nil {
            return fmt.Errorf("加载配置失败: %w", err)
        }
        
        // 执行任务
        if err := executeTask(runTask, cfg); err != nil {
            // 根据错误类型处理
            if errors.Is(err, ErrTaskNotFound) {
                fmt.Printf("任务不存在: %s\n", runTask)
                fmt.Println("使用 --list 查看所有可用任务")
                return nil
            }
            return fmt.Errorf("执行任务失败: %w", err)
        }
    }
    
    return nil
}
```

---

## 常见问题

### Q1: 如何处理可选参数？

```go
// 使用指针判断是否设置
var optionalFlag *qflag.StringFlag

// 在 run 函数中
if optionalFlag.IsSet() {
    // 标志被设置了
    value := optionalFlag.Get()
} else {
    // 使用默认值或其他逻辑
}
```

### Q2: 如何实现全局标志？

全局标志定义在根命令上，所有子命令都可以访问：

```go
// internal/cli/root.go
var (
    verboseFlag *qflag.BoolFlag
    debugFlag   *qflag.BoolFlag
)

func InitAndRun() (err error) {
    // 在 qflag.Root 上定义全局标志
    verboseFlag = qflag.Root.Bool("verbose", "v", "详细输出", false)
    debugFlag = qflag.Root.Bool("debug", "d", "调试模式", false)
    // ...
}

// internal/cli/run.go
func runRun(cmd qflag.Command) error {
    // 访问全局标志
    verbose := verboseFlag.Get()
    debug := debugFlag.Get()
    // ...
}
```

### Q3: 如何实现命令别名？

目前 qflag 不直接支持命令别名，可以通过创建多个命令变量实现：

```go
// internal/cli/build.go
var (
    BuildCmd *qflag.Cmd
    BCmd     *qflag.Cmd  // 别名
)

func init() {
    // 主命令
    BuildCmd = qflag.NewCmd("build", "b", qflag.ExitOnError)
    // ... 配置 BuildCmd
    
    // 别名命令（指向同一个 run 函数）
    BCmd = qflag.NewCmd("b", "", qflag.ExitOnError)
    BCmd.SetRun(runBuild)
}

// 在 root.go 中注册
SubCmds: []qflag.Command{
    BuildCmd,
    BCmd,  // 注册别名
}
```

### Q4: 如何自定义帮助信息？

```go
cmdOpts := &qflag.CmdOpts{
    Desc: "构建项目",
    UsageSyntax: fmt.Sprintf("%s build [选项]", qflag.Root.Name()),
    Examples: map[string]string{
        "构建当前项目":     fmt.Sprintf("%s build", qflag.Root.Name()),
        "构建指定平台":     fmt.Sprintf("%s build --platform linux/amd64", qflag.Root.Name()),
        "构建并压缩":       fmt.Sprintf("%s build --compress", qflag.Root.Name()),
    },
    Notes: []string{
        "默认构建当前平台的二进制文件",
        "使用 --platform 可以指定目标平台",
        "支持的平台: linux/amd64, windows/amd64, darwin/amd64",
    },
}
```

### Q5: 如何处理子命令的参数？

```go
func runRun(cmd qflag.Command) error {
    // 获取位置参数
    args := cmd.Args()
    
    if len(args) == 0 {
        return fmt.Errorf("缺少必需参数")
    }
    
    // 获取第一个参数
    taskName := cmd.Arg(0)
    
    // 获取所有参数
    for i, arg := range args {
        fmt.Printf("参数 %d: %s\n", i, arg)
    }
    
    return nil
}
```

---

## 内部包组织规范

### 1. 工具包 (internal/utils)

```go
package utils

import "fmt"

const (
    logErrPrefix = "err:"
    logPrefix    = "===>>"
    logTitle     = "======>>"
)

// LogTitle 输出标题日志
func LogTitle(msg string) {
    fmt.Printf("%s %s\n", logTitle, msg)
}

// Log 打印普通日志
func Log(msg string) {
    fmt.Printf("%s %s\n", logPrefix, msg)
}

// LogErr 打印错误日志
func LogErr(msg string) {
    fmt.Printf("%s %s\n", logErrPrefix, msg)
}
```

### 2. 配置包 (internal/config)

```go
package config

// Config 应用配置结构
type Config struct {
    Debug   bool
    LogPath string
    Port    int
}

// Load 加载配置
func Load(path string) (*Config, error) {
    // 实现配置加载逻辑
    return &Config{}, nil
}
```

### 3. 服务包 (internal/service)

```go
package service

// Service 业务服务接口
type Service interface {
    Process(data string) (string, error)
}

// NewService 创建服务实例
func NewService() Service {
    return &serviceImpl{}
}

type serviceImpl struct{}

func (s *serviceImpl) Process(data string) (string, error) {
    return data, nil
}
```

---

## 代码注释规范

所有公共函数必须添加函数级注释，遵循 Go 语言标准注释格式：

```go
// FunctionName 函数功能简述
// 函数功能详细描述（可选）
//
// 参数:
//   - param1: 参数1描述
//   - param2: 参数2描述
//
// 返回值:
//   - returnType1: 返回值1描述
//   - error: 错误信息描述
func FunctionName(param1 string, param2 int) (returnType1, error) {
    // 函数实现
}
```

---

## 错误处理规范

### 错误创建

```go
// 使用 fmt.Errorf 创建错误信息
if err != nil {
    return fmt.Errorf("操作失败: %w", err)
}
```

### 错误处理策略

| 策略 | 使用场景 | 说明 |
|------|---------|------|
| `qflag.ExitOnError` | 生产环境 | 遇到错误立即退出程序 |
| `qflag.ContinueOnError` | 测试环境 | 返回错误，不退出程序 |
| `qflag.PanicOnError` | 开发环境 | 遇到错误触发 panic |

---

## 完整示例项目

### 项目结构

```
gob/
├── cmd/
│   └── gob/
│       └── main.go
├── internal/
│   ├── cli/
│   │   ├── root.go
│   │   ├── init.go
│   │   ├── run.go
│   │   ├── build.go
│   │   ├── version.go
│   │   └── config/
│   │       ├── get.go
│   │       └── set.go
│   ├── utils/
│   │   └── log.go
│   └── config/
│       └── config.go
├── go.mod
└── README.md
```

### main.go

```go
package main

import (
    "os"

    "gitee.com/MM-Q/gob/internal/cli"
    "gitee.com/MM-Q/gob/internal/utils"
)

func main() {
    if err := cli.InitAndRun(); err != nil {
        utils.LogErr(err.Error())
        os.Exit(1)
    }
}
```

### internal/cli/root.go

```go
package cli

import (
    "fmt"

    "gitee.com/MM-Q/gob/internal/config"
    "gitee.com/MM-Q/qflag"
    "gitee.com/MM-Q/verman"
)

const logoText = `
  ____ _           
 / ___| |__   ___  
| |   | '_ \ / _ \ 
| |___| | | |  __/ 
 \____|_| |_|\___| 
`

var (
    verboseFlag *qflag.BoolFlag
    configFlag  *qflag.StringFlag
)

// InitAndRun 初始化并运行根命令
//
// 返回值:
//   - err: 初始化或运行命令时可能发生的错误
func InitAndRun() (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic recovered: %v", r)
        }
    }()

    // 注册根命令的标志
    verboseFlag = qflag.Root.Bool("verbose", "v", "详细输出", false)
    configFlag = qflag.Root.String("config", "c", "配置文件路径", "")

    // 配置根命令
    rootCmdOpts := &qflag.CmdOpts{
        Version:    verman.V.Version(),
        Desc:       "Go 项目构建工具",
        LogoText:   logoText,
        UseChinese: true,
        Completion: true,
        RunFunc:    run,
        Examples: map[string]string{
            "初始化项目": fmt.Sprintf("%s init", qflag.Root.Name()),
            "运行任务":   fmt.Sprintf("%s run -i main.go", qflag.Root.Name()),
            "构建项目":   fmt.Sprintf("%s build", qflag.Root.Name()),
        },
        SubCmds: []qflag.Command{
            InitCmd,
            RunCmd,
            BuildCmd,
            VersionCmd,
        },
    }

    if err = qflag.ApplyOpts(rootCmdOpts); err != nil {
        err = fmt.Errorf("apply opts failed: %w", err)
        return err
    }

    if err = qflag.ParseAndRoute(); err != nil {
        err = fmt.Errorf("parse and route failed: %w", err)
        return err
    }

    return nil
}

func run(cmd qflag.Command) error {
    cmd.PrintHelp()
    return nil
}
```

---

## 规范总结

### 核心原则

1. ✅ **入口统一**: `cmd/yourapp/main.go`
2. ✅ **命令集中**: `internal/cli/`
3. ✅ **单文件原则**: 每个命令 = 单个文件
4. ✅ **按需创建**: 只有两级子命令才创建目录
5. ✅ **根命令特殊**: 使用 `InitAndRun()` 模式，直接操作 `qflag.Root`

### 根命令 vs 子命令

| 特性 | 根命令 | 子命令 |
|------|--------|--------|
| **初始化方式** | `InitAndRun()` 函数 | `init()` 函数 |
| **命令对象** | 直接使用 `qflag.Root` | 创建新命令对象 `qflag.NewCmd()` |
| **标志定义** | 在 `qflag.Root` 上定义 | 在命令对象上定义 |
| **命令变量** | 不需要导出 | 必须导出供上级注册 |
| **调用方式** | `cli.InitAndRun()` | 由根命令的 `SubCmds` 注册 |

### 文件结构四要素

#### 根命令 (root.go)

1. **全局标志变量** - 定义在 `qflag.Root` 上
2. **InitAndRun() 函数** - 初始化和运行入口
3. **defer 捕获 panic** - 保证错误统一处理
4. **run() 函数** - 业务逻辑实现

#### 子命令 (*.go)

1. **全局命令变量** - 供注册使用
2. **全局标志变量** - 用于 run 函数
3. **init() 函数** - 初始化命令和标志
4. **run() 函数** - 业务逻辑实现

### 核心功能清单

| 功能 | 状态 | 说明 |
|------|------|------|
| 互斥组 | ✅ | `MutexGroups` 确保组内最多一个标志 |
| 必需组 | ✅ | `RequiredGroups` 确保组内所有标志 |
| 自动补全 | ✅ | `--completion bash/pwsh/fish` |
| 标志验证 | ✅ | 使用 `SetValidator()` 设置验证器 |
| 全局标志 | ✅ | 在根命令上定义，所有子命令可访问 |
| 多级子命令 | ✅ | 支持二级及更多级子命令 |
| 配置文件 | ✅ | 结合 `internal/config` 使用 |

### 命名规范

- 命令变量: `XxxCmd`
- 标志变量: `命令前缀 + 标志名`
- run 函数: `run + 命令名`

### 导入顺序

```go
import (
    // 标准库
    "fmt"
    "os"
    
    // 第三方库
    "gitee.com/MM-Q/qflag"
    
    // 内部包
    "your-project/internal/utils"
    "your-project/internal/config"
)
```

### 开发流程

1. **创建项目结构**: 按照目录规范创建项目
2. **定义根命令**: 在 `internal/cli/root.go` 中使用 `InitAndRun()`
3. **定义子命令**: 在 `internal/cli/*.go` 中使用 `init()` + `Cmd` 变量
4. **注册子命令**: 在父命令的 `SubCmds` 中注册
5. **实现业务逻辑**: 在 `run()` 函数中实现
6. **编写测试**: 为每个命令编写单元测试
7. **生成补全**: 使用 `--completion` 生成补全脚本

### 检查清单

在完成命令行工具开发后，检查以下项目：

- [ ] 入口文件 `main.go` 是否极简（只调用 `cli.InitAndRun()`）
- [ ] 根命令是否使用 `InitAndRun()` 模式
- [ ] 子命令是否导出 `Cmd` 变量
- [ ] 标志变量是否使用正确的前缀命名
- [ ] 是否正确注册子命令到父命令
- [ ] 是否添加了必要的验证器
- [ ] 是否配置了互斥组和必需组（如需要）
- [ ] 是否编写了单元测试
- [ ] 是否生成了自动补全脚本
- [ ] 是否添加了代码注释
- [ ] 是否更新了 README 文档

遵循此规范可以保证项目的**可维护性、可扩展性和一致性**。
