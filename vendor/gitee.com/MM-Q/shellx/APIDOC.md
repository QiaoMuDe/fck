# shellx API 文档

## 包信息

```go
package shellx // import "gitee.com/MM-Q/shellx"
```

## 概述

Package shellx 提供了一个功能完善、易于使用的Go语言shell命令执行库。

本库基于Go标准库的os/exec包进行封装，提供了更加友好的API和丰富的功能，支持同步和异步命令执行、输入输出重定向、精确超时控制、上下文管理、多种shell类型支持等功能，并提供类型安全的API和友好的链式调用接口。

### 主要特性

- 一体化设计：Command集配置、构建、执行于一体，无需Build()方法
- 支持三种命令创建方式：NewCmd(可变参数)、NewCmds(切片)、NewCmdStr(字符串解析)
- 丰富的便捷函数：Exec、ExecStr、ExecOut、ExecOutStr及其带超时版本
- 链式调用API，支持流畅的方法链
- 精确超时控制：延迟构建exec.Cmd，确保超时计时精确
- 完整的错误处理和类型安全
- 支持多种shell类型（sh、bash、cmd、powershell、pwsh等）
- 同步和异步执行支持
- 命令执行状态管理和进程控制
- 输入输出重定向和环境变量设置
- 上下文取消和优先级控制
- 无锁设计，高性能
- 跨平台兼容（Windows、Linux、macOS）

### 并发安全说明

- Command 对象的配置方法 (WithXxx) 不是并发安全的，不要在多个 goroutine 中并发配置
- 每个 Command 对象只能执行一次，重复执行会返回错误
- 执行方法是并发安全的，使用 atomic.Bool 防止重复执行
- 属性获取方法不是并发安全的，不要在多个 goroutine 中并发调用

### 核心组件

- **Command**: 命令对象，集配置、构建、执行于一体
- **ShellType**: Shell类型枚举，支持多种shell

### 基本用法

```go
import "gitee.com/MM-Q/shellx"

// 方式1：使用可变参数创建命令（无需Build）
err := shellx.NewCmd("ls", "-la").
    WithWorkDir("/tmp").
    WithTimeout(30 * time.Second).
    WithShell(shellx.ShellBash).
    Exec()

// 方式2：使用字符串创建命令
output, err := shellx.NewCmdStr(`echo "hello world"`).
    WithEnv("MY_VAR", "value").
    ExecOutput()
```

### 便捷函数用法

```go
// 基础执行函数
err := shellx.Exec("ls", "-la")                    // 执行命令，输出到控制台
err := shellx.ExecStr("echo hello")                // 字符串方式执行
output, err := shellx.ExecOut("ls", "-la")         // 执行并返回输出
output, err := shellx.ExecOutStr("echo hello")     // 字符串方式执行并返回输出

// 带超时的执行函数
err := shellx.ExecT(5*time.Second, "sleep", "10")                    // 5秒超时
err := shellx.ExecStrT(3*time.Second, "ping google.com")       // 字符串方式，3秒超时
output, err := shellx.ExecOutT(2*time.Second, "curl", "example.com") // 返回输出，2秒超时
output, err := shellx.ExecOutStrT(1*time.Second, "date")             // 字符串方式，返回输出，1秒超时
```

### 高级用法

```go
// 设置标准输入输出
var stdout, stderr bytes.Buffer
stdin := strings.NewReader("input data")

err := shellx.NewCmd("cat").
    WithStdin(stdin).
    WithStdout(&stdout).
    WithStderr(&stderr).
    Exec()

// 使用上下文控制（优先级高于WithTimeout）
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

err := shellx.NewCmd("long-running-command").
    WithContext(ctx).
    WithTimeout(5*time.Second).  // 这个会被忽略
    Exec()

// 异步执行和进程控制
cmd := shellx.NewCmd("sleep", "100")
err := cmd.ExecAsync()
pid := cmd.GetPID()
isRunning := cmd.IsRunning()
cmd.Kill() // 或 cmd.Signal(syscall.SIGTERM)
err = cmd.Wait()
```

### 命令解析

```go
// 支持复杂的命令字符串解析，包括引号处理
cmd := shellx.NewCmdStr(`git commit -m "Initial commit with 'quotes'"`)
// 解析结果：["git", "commit", "-m", "Initial commit with 'quotes'"]
```

### Shell类型

```go
// 支持多种shell类型
shellx.ShellSh         // sh shell
shellx.ShellBash       // bash shell
shellx.ShellCmd        // Windows cmd
shellx.ShellPowerShell // Windows PowerShell
shellx.ShellPwsh       // PowerShell Core
shellx.ShellNone       // 直接执行，不使用shell
shellx.ShellDef1       // 默认shell，根据操作系统自动选择 (Windows: cmd, 其他: sh)
shellx.ShellDef2       // 默认shell，根据操作系统自动选择 (Windows: PowerShell, 其他: sh)
```

### 注意事项

- 每个 Command 对象只能执行一次，重复执行会返回错误
- Command 对象的配置方法 (WithXxx) 不是并发安全的，不要在多个 goroutine 中并发配置
- 执行方法是并发安全的，使用 atomic.Bool 防止重复执行
- 属性获取方法不是并发安全的，不要在多个 goroutine 中并发调用
- 命令执行会继承父进程的环境变量，可通过WithEnv添加额外变量
- 超时控制在执行时创建上下文，确保计时精确
- 异步执行需要调用Wait()等待完成或使用Kill()终止

---

## 变量

```go
var (
    // ErrAlreadyExecuted 表示命令已经执行过
    ErrAlreadyExecuted = errors.New("command has already been executed")
    // ErrNotStarted 表示命令尚未启动
    ErrNotStarted = errors.New("command has not been started")
    // ErrNoProcess 表示没有进程可操作
    ErrNoProcess = errors.New("no process to operate")
)
```

---

## 函数

### Exec

```go
func Exec(name string, args ...string) error
```

执行命令(阻塞)

**参数:**
- `name`: 命令名
- `args`: 命令参数

**返回:**
- `error`: 错误信息

---

### ExecCode

```go
func ExecCode(name string, args ...string) (int, error)
```

执行命令并返回退出码(阻塞)

**参数:**
- `name`: 命令名
- `args`: 命令参数

**返回:**
- `int`: 退出码
- `error`: 错误信息

---

### ExecCodeStr

```go
func ExecCodeStr(cmdStr string) (int, error)
```

字符串方式执行命令并返回退出码(阻塞)

**参数:**
- `cmdStr`: 命令字符串

**返回:**
- `int`: 退出码
- `error`: 错误信息

---

### ExecOut

```go
func ExecOut(name string, args ...string) ([]byte, error)
```

执行命令并返回合并后的输出(阻塞)

**参数:**
- `name`: 命令名
- `args`: 命令参数

**返回:**
- `[]byte`: 输出
- `error`: 错误信息

**注意:**
- 由于需要捕获默认的stdout和stderr合并输出, 内部已经设置了WithStdout(os.Stdout)和WithStderr(os.Stderr)

---

### ExecOutStr

```go
func ExecOutStr(cmdStr string) ([]byte, error)
```

执行命令并返回合并后的输出(阻塞)

**参数:**
- `cmdStr`: 命令字符串

**返回:**
- `[]byte`: 输出
- `error`: 错误信息

**注意:**
- 由于需要捕获默认的stdout和stderr合并输出, 内部已经设置了WithStdout(os.Stdout)和WithStderr(os.Stderr)

---

### ExecOutStrT

```go
func ExecOutStrT(timeout time.Duration, cmdStr string) ([]byte, error)
```

执行命令并返回合并后的输出(阻塞，带超时)

**参数:**
- `timeout`: 超时时间，如果为0则不设置超时
- `cmdStr`: 命令字符串

**返回:**
- `[]byte`: 合并后的输出
- `error`: 错误信息

---

### ExecOutT

```go
func ExecOutT(timeout time.Duration, name string, args ...string) ([]byte, error)
```

执行命令并返回合并后的输出(阻塞，带超时)

**参数:**
- `timeout`: 超时时间，如果为0则不设置超时
- `name`: 命令名
- `args`: 命令参数

**返回:**
- `[]byte`: 合并后的输出
- `error`: 错误信息

---

### ExecStr

```go
func ExecStr(cmdStr string) error
```

执行命令(阻塞)

**参数:**
- `cmdStr`: 命令字符串

**返回:**
- `error`: 错误信息

---

### ExecStrT

```go
func ExecStrT(timeout time.Duration, cmdStr string) error
```

执行命令(阻塞，带超时)

**参数:**
- `timeout`: 超时时间，如果为0则不设置超时
- `cmdStr`: 命令字符串

**返回:**
- `error`: 错误信息

---

### ExecT

```go
func ExecT(timeout time.Duration, name string, args ...string) error
```

执行命令(阻塞，带超时)

**参数:**
- `timeout`: 超时时间，如果为0则不设置超时
- `name`: 命令名
- `args`: 命令参数

**返回:**
- `error`: 错误信息

---

### FindCmd

```go
func FindCmd(name string) (string, error)
```

查找命令

**参数:**
- `name`: 命令名称

**返回:**
- `string`: 命令路径
- `error`: 错误信息

---

### Split

```go
func Split(cmdStr string) []string
```

将命令字符串拆分为命令切片，支持引号处理(单引号、双引号、反引号)

**功能：**
- 智能拆分 Shell 命令字符串为参数数组
- 支持单引号、双引号、反引号包裹的内容
- 正确处理转义字符和特殊字符
- 自动处理命令分隔符(;|&&||)

**参数:**
- `cmdStr`: 要拆分的命令字符串

**返回值:**
- `[]string`: 拆分后的命令切片 (最佳结果)

**注意：**
- 此函数忽略拆分错误，返回最佳拆分结果。如需错误信息，请使用 SplitE 函数。
- 转义字符保持原样，不进行解释处理
- 支持多字符操作符如 &&、||、>>、<< 等

**特性：**
- 支持单引号、双引号、反引号
- 正确处理空引号和连续引号
- 支持常见空白字符处理（空格、制表符、换行符等）
- 支持跨平台换行符（\n, \r, \r\n）
- 使用strings.Builder优化性能
- 完整的 Unicode 字符支持（中文、emoji 等多字节字符）

**示例:**
```go
// 基本拆分
result := shellx.Split("echo hello world")
// 结果: ["echo", "hello", "world"]

// 引号处理
result := shellx.Split(`echo "hello world" 'test'`)
// 结果: ["echo", "hello world", "test"]

// 空引号处理
result := shellx.Split(`echo ""hello""`)
// 结果: ["echo", "hello"]

// 未闭合引号（忽略错误）
result := shellx.Split(`echo "hello world`)
// 结果: ["echo", "hello world"]  // 忽略未闭合引号错误

// Unicode 字符支持
result := shellx.Split("echo 你好世界 🎉")
// 结果: ["echo", "你好世界", "🎉"]

// 跨平台换行符
result := shellx.Split("echo \"line1\nline2\"")
// 结果: ["echo", "line1\nline2"]
```

---

### SplitE

```go
func SplitE(cmdStr string) ([]string, error)
```

将命令字符串拆分为命令切片（带错误信息），支持引号处理(单引号、双引号、反引号)

**功能：**
- 智能拆分 Shell 命令字符串为参数数组
- 支持单引号、双引号、反引号包裹的内容
- 正确处理转义字符和特殊字符
- 自动处理命令分隔符(;|&&||)
- 检测并返回拆分过程中的错误

**参数:**
- `cmdStr`: 要拆分的命令字符串

**返回值:**
- `[]string`: 拆分后的命令切片
- `error`: 拆分错误，成功时为 nil

**错误类型：**
- `*UnclosedQuoteError`: 未闭合的引号错误
- 其他可能的语法错误

**注意：**
- 转义字符保持原样，不进行解释处理
- 支持多字符操作符如 &&、||、>>、<< 等

**特性：**
- 支持单引号、双引号、反引号
- 正确处理空引号和连续引号
- 支持常见空白字符处理（空格、制表符、换行符等）
- 支持跨平台换行符（\n, \r, \r\n）
- 使用strings.Builder优化性能
- 返回详细的错误信息

**示例:**
```go
// 基本拆分
result, err := shellx.SplitE("echo hello world")
// 结果: ["echo", "hello", "world"], err: nil

// 引号处理
result, err := shellx.SplitE(`echo "hello world" 'test'`)
// 结果: ["echo", "hello world", "test"], err: nil

// 未闭合引号错误
result, err := shellx.SplitE(`echo "hello world`)
// 结果: ["echo", "hello world"], err: *UnclosedQuoteError{QuoteType: '"'}

// 检查错误类型
if err != nil {
    if unclosedErr, ok := err.(*shellx.UnclosedQuoteError); ok {
        fmt.Printf("未闭合的引号类型: %q\n", unclosedErr.QuoteType)
    }
}

// Unicode 字符支持
result, err := shellx.SplitE("echo 你好世界 🎉")
// 结果: ["echo", "你好世界", "🎉"], err: nil
```

---

## 类型

### UnclosedQuoteError

```go
type UnclosedQuoteError struct {
    QuoteType rune // 未闭合的引号类型 ('"', '\'', 或 '`')
}
```

UnclosedQuoteError 表示命令字符串中存在未闭合的引号

**方法:**

#### Error

```go
func (e *UnclosedQuoteError) Error() string
```

实现 error 接口，返回错误信息

**返回:**
- `string`: 错误信息

#### GetQuoteType

```go
func (e *UnclosedQuoteError) GetQuoteType() rune
```

返回未闭合的引号类型

**返回:**
- `rune`: 未闭合的引号字符

---

### Command

```go
type Command struct {
    // Has unexported fields.
}
```

Command 命令对象 - 集配置、构建、执行于一体

**注意事项:**
- Command 对象的配置方法 (WithXxx) 不是并发安全的，不要在多个 goroutine 中并发配置
- 每个 Command 对象只能执行一次，重复执行会返回错误
- 执行方法是并发安全的，使用 atomic.Bool 防止重复执行
- 属性获取方法不是并发安全的，不要在多个 goroutine 中并发调用

---

#### NewCmd

```go
func NewCmd(name string, args ...string) *Command
```

创建新的命令对象 (数组方式 - 可变参数)

**参数：**
- `name`: 命令名
- `args`: 命令参数列表

**返回：**
- `*Command`: 命令对象

**注意:**
- 默认通过shell执行, 可以通过WithShell方法指定shell类型
- 默认为ShellDef1, 根据操作系统自动选择shell(Windows系统默认为cmd, 其他系统默认为sh)
- 默认继承父进程的环境变量, 可以通过WithEnv方法设置环境变量

---

#### NewCmdStr

```go
func NewCmdStr(cmdStr string) *Command
```

创建新的命令对象 (字符串方式)

**参数：**
- `cmdStr`: 命令字符串

**返回：**
- `*Command`: 命令对象

**注意:**
- 默认通过shell执行, 可以通过WithShell方法指定shell类型
- 默认为ShellDef1, 根据操作系统自动选择shell(Windows系统默认为cmd, 其他系统默认为sh)
- 默认继承父进程的环境变量, 可以通过WithEnv方法设置环境变量

---

#### NewCmds

```go
func NewCmds(cmdArgs []string) *Command
```

创建新的命令对象 (数组方式 - 切片参数)

**参数：**
- `cmdArgs`: 命令参数列表，第一个元素为命令名，后续元素为参数

**返回：**
- `*Command`: 命令对象

**注意:**
- 默认通过shell执行, 可以通过WithShell方法指定shell类型
- 默认为ShellDef1, 根据操作系统自动选择shell(Windows系统默认为cmd, 其他系统默认为sh)
- 默认继承父进程的环境变量, 可以通过WithEnv方法设置环境变量

---

#### Args

```go
func (c *Command) Args() []string
```

获取命令参数列表

**返回:**
- `[]string`: 命令参数列表

**注意:**
- 此方法不是并发安全的，不要在多个goroutine中并发调用
- 返回的是参数的副本，修改返回值不会影响原始对象

---

#### Cmd

```go
func (c *Command) Cmd() *exec.Cmd
```

获取底层的 exec.Cmd 对象

**返回:**
- `*exec.Cmd`: 底层的 exec.Cmd 对象

---

#### CmdStr

```go
func (c *Command) CmdStr() string
```

获取命令字符串

**返回:**
- `string`: 命令字符串

---

#### Env

```go
func (c *Command) Env() []string
```

获取命令环境变量列表

**返回:**
- `[]string`: 命令环境变量列表

**注意:**
- 此方法不是并发安全的，不要在多个goroutine中并发调用
- 返回的是环境变量的副本，修改返回值不会影响原始对象

---

#### Exec

```go
func (c *Command) Exec() error
```

执行命令(阻塞)

**返回:**
- `error`: 错误信息，可通过 IsTimeoutError() 和 IsCanceledError() 判断错误类型

---

#### ExecAsync

```go
func (c *Command) ExecAsync() error
```

异步执行命令(非阻塞)

**返回:**
- `error`: 错误信息，可通过 IsTimeoutError() 和 IsCanceledError() 判断错误类型

---

#### ExecOutput

```go
func (c *Command) ExecOutput() ([]byte, error)
```

执行命令并返回合并后的输出(阻塞)

**返回:**
- `[]byte`: 命令输出
- `error`: 错误信息，可通过 IsTimeoutError() 和 IsCanceledError() 判断错误类型

**注意:**
- 由于需要捕获默认的stdout和stderr合并输出, 内部已经设置了WithStdout(os.Stdout)和WithStderr(os.Stderr)

---

#### ExecStdout

```go
func (c *Command) ExecStdout() ([]byte, error)
```

执行命令并返回标准输出(阻塞)

**返回:**
- `[]byte`: 标准输出
- `error`: 错误信息，可通过 IsTimeoutError() 和 IsCanceledError() 判断错误类型

---

#### GetPID

```go
func (c *Command) GetPID() int
```

获取进程ID

**返回:**
- `int`: 进程ID, 如果进程不存在返回0

---

#### IsExecuted

```go
func (c *Command) IsExecuted() bool
```

检查命令是否已经执行过

**返回:**
- `bool`: 是否已执行

---

#### IsRunning

```go
func (c *Command) IsRunning() bool
```

检查进程是否还在运行

**注意:** 此方法提供基本的进程状态检查，可能不是100%准确，特别是在Windows系统上可能存在权限问题。对于精确的进程状态管理，建议使用Wait或WaitWithCode方法。

**返回:**
- `bool`: 是否在运行

---

#### Kill

```go
func (c *Command) Kill() error
```

杀死当前命令的进程

**返回:**
- `error`: 错误信息

---

#### Name

```go
func (c *Command) Name() string
```

获取命令名称

**返回:**
- `string`: 命令名称

**注意:**
- 此方法不是并发安全的，不要在多个goroutine中并发调用

---

#### Raw

```go
func (c *Command) Raw() string
```

获取原始命令字符串

**返回:**
- `string`: 原始命令字符串

**注意:**
- 此方法不是并发安全的，不要在多个goroutine中并发调用

---

#### ShellType

```go
func (c *Command) ShellType() ShellType
```

获取shell类型

**返回:**
- `ShellType`: shell类型

**注意:**
- 此方法不是并发安全的，不要在多个goroutine中并发调用

---

#### Signal

```go
func (c *Command) Signal(sig os.Signal) error
```

向当前进程发送信号

**参数:**
- `sig`: 信号类型

**返回:**
- `error`: 错误信息

---

#### Timeout

```go
func (c *Command) Timeout() time.Duration
```

获取命令执行超时时间

**返回:**
- `time.Duration`: 命令执行超时时间

**注意:**
- 此方法不是并发安全的，不要在多个goroutine中并发调用

---

#### Wait

```go
func (c *Command) Wait() error
```

等待命令执行完成(仅在异步执行时有效)

**返回:**
- `error`: 错误信息，可通过 IsTimeoutError() 和 IsCanceledError() 判断错误类型

---

#### WaitWithCode

```go
func (c *Command) WaitWithCode() (int, error)
```

等待命令执行完成并返回退出码(仅在异步执行时有效)

**返回:**
- `int`: 命令退出码(0表示成功，-1表示无法提取的执行错误，其他值表示命令返回的退出码)
- `error`: 错误信息，可通过 IsTimeoutError() 和 IsCanceledError() 判断错误类型

---

#### WithContext

```go
func (c *Command) WithContext(ctx context.Context) *Command
```

设置命令的上下文

**参数：**
- `ctx`: context.Context类型，用于取消命令执行和超时控制

**返回：**
- `*Command`: 命令对象

**注意:**
- 该方法会验证上下文是否为空，如果为空则panic.
- 该上下文会覆盖之前设置的超时时间.
- 此方法不是并发安全的，不要在多个goroutine中并发配置

---

#### WithEnv

```go
func (c *Command) WithEnv(key, value string) *Command
```

设置命令的环境变量

**参数：**
- `key`: 环境变量的键
- `value`: 环境变量的值

**返回：**
- `*Command`: 命令对象

**注意:**
- 该方法会验证key是否为空, 如果为空则忽略。
- 无需添加系统环境变量os.Environ(), 系统环境变量会自动继承.
- 此方法不是并发安全的，不要在多个goroutine中并发配置

---

#### WithEnvs

```go
func (c *Command) WithEnvs(envs []string) *Command
```

批量设置命令的环境变量

**参数：**
- `envs`: []string类型，环境变量列表，每个元素为"key=value"格式

**返回：**
- `*Command`: 命令对象

**注意:**
- 该方法会验证环境变量格式，只添加验证通过的环境变量。
- 无需添加系统环境变量os.Environ(), 系统环境变量会自动继承.
- 此方法不是并发安全的，不要在多个goroutine中并发配置

---

#### WithShell

```go
func (c *Command) WithShell(shell ShellType) *Command
```

设置命令的shell类型

**参数：**
- `shell`: ShellType类型，表示要使用的shell类型

**返回：**
- `*Command`: 命令对象

**注意:**
- 此方法不是并发安全的，不要在多个goroutine中并发配置

---

#### WithStderr

```go
func (c *Command) WithStderr(stderr io.Writer) *Command
```

设置命令的标准错误输出

**参数：**
- `stderr`: io.Writer类型，用于接收命令的标准错误输出

**返回：**
- `*Command`: 命令对象

**注意:**
- 此方法不是并发安全的，不要在多个goroutine中并发配置

---

#### WithStdin

```go
func (c *Command) WithStdin(stdin io.Reader) *Command
```

设置命令的标准输入

**参数：**
- `stdin`: io.Reader类型，用于提供命令的标准输入

**返回：**
- `*Command`: 命令对象

**注意:**
- 此方法不是并发安全的，不要在多个goroutine中并发配置

---

#### WithStdout

```go
func (c *Command) WithStdout(stdout io.Writer) *Command
```

设置命令的标准输出

**参数：**
- `stdout`: io.Writer类型，用于接收命令的标准输出

**返回：**
- `*Command`: 命令对象

**注意:**
- 此方法不是并发安全的，不要在多个goroutine中并发配置

---

#### WithTimeout

```go
func (c *Command) WithTimeout(timeout time.Duration) *Command
```

设置命令的超时时间(便捷方式)

**参数：**
- `timeout`: time.Duration类型，命令执行的超时时间

**返回：**
- `*Command`: 命令对象

**注意:**
- 该方法会验证超时时间是否小于等于0, 如果小于等于0则忽略。
- 该超时时间优先级低于上下文设置的超时时间.
- 此方法不是并发安全的，不要在多个goroutine中并发配置

---

#### WithWorkDir

```go
func (c *Command) WithWorkDir(dir string) *Command
```

设置命令的工作目录

**参数：**
- `dir`: 命令的工作目录

**返回：**
- `*Command`: 命令对象

**注意:**
- 此方法不是并发安全的，不要在多个goroutine中并发配置

---

#### WorkDir

```go
func (c *Command) WorkDir() string
```

获取命令执行的工作目录

**返回:**
- `string`: 命令执行目录

**注意:**
- 此方法不是并发安全的，不要在多个goroutine中并发调用

---

### ShellType

```go
type ShellType int
```

ShellType 定义shell类型

---

#### 常量

```go
const (
    ShellSh         ShellType = iota // sh shell
    ShellBash                        // bash shell
    ShellPwsh                        // pwsh (PowerShell Core)
    ShellPowerShell                  // powershell (Windows PowerShell)
    ShellCmd                         // cmd (Windows Command Prompt)
    ShellNone                        // 无shell, 直接原生的执行命令
    ShellDef1                        // 默认shell, 根据操作系统自动选择(Windows系统默认为cmd, 其他系统默认为sh)
    ShellDef2                        // 默认shell, 根据操作系统自动选择(Windows系统默认为powershell, 其他系统默认为sh)
)
```

---

#### String

```go
func (s ShellType) String() string
```

String 返回shell类型的字符串表示
