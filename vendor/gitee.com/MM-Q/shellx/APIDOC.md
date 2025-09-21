# ShellX API 文档

```go
import "gitee.com/MM-Q/shellx"
```

## 📚 API 概述

### 核心类型

| 类型 | 描述 |
|------|------|
| `Command` | 命令对象，集配置、构建、执行于一体 |
| `Result` | 命令执行结果，包含输出、错误、时间等信息 |
| `ShellType` | Shell 类型枚举，支持多种 shell |

### 主要方法

#### 创建命令

```go
// 可变参数方式
func NewCmd(name string, args ...string) *Command

// 切片方式
func NewCmds(cmdArgs []string) *Command

// 字符串解析方式
func NewCmdStr(cmdStr string) *Command
```

#### 链式配置

```go
func (c *Command) WithWorkDir(dir string) *Command
func (c *Command) WithEnv(key, value string) *Command
func (c *Command) WithEnvs(envs []string) *Command
func (c *Command) WithTimeout(timeout time.Duration) *Command
func (c *Command) WithContext(ctx context.Context) *Command
func (c *Command) WithStdin(stdin io.Reader) *Command
func (c *Command) WithStdout(stdout io.Writer) *Command
func (c *Command) WithStderr(stderr io.Writer) *Command
func (c *Command) WithShell(shell ShellType) *Command
```

#### 信息获取

```go
func (c *Command) CmdStr() string  // 获取命令字符串
```

#### 便捷函数

```go
// 基础执行函数
func Exec(name string, args ...string) error
func ExecStr(cmdStr string) error
func ExecOut(name string, args ...string) ([]byte, error)
func ExecOutStr(cmdStr string) ([]byte, error)

// 带超时的执行函数
func ExecT(timeout time.Duration, name string, args ...string) error
func ExecStrT(timeout time.Duration, cmdStr string) error
func ExecOutT(timeout time.Duration, name string, args ...string) ([]byte, error)
func ExecOutStrT(timeout time.Duration, cmdStr string) ([]byte, error)
```

#### 命令执行

```go
// 同步执行
func (c *Command) Exec() error
func (c *Command) ExecOutput() ([]byte, error)
func (c *Command) ExecStdout() ([]byte, error)
func (c *Command) ExecResult() (*Result, error)

// 异步执行
func (c *Command) ExecAsync() error
func (c *Command) Wait() error

// 进程控制
func (c *Command) Kill() error
func (c *Command) Signal(sig os.Signal) error
func (c *Command) IsRunning() bool
func (c *Command) GetPID() int
func (c *Command) IsExecuted() bool

// 信息获取
func (c *Command) CmdStr() string
```

### Shell 类型支持

| Shell 类型 | 常量 | 平台支持 | 描述 |
|------------|------|----------|------|
| **sh** | `ShellSh` | Unix/Linux/macOS | 标准 Unix shell |
| **bash** | `ShellBash` | Unix/Linux/macOS | Bash shell |
| **cmd** | `ShellCmd` | Windows | Windows 命令提示符 |
| **powershell** | `ShellPowerShell` | Windows | Windows PowerShell |
| **pwsh** | `ShellPwsh` | 跨平台 | PowerShell Core |
| **none** | `ShellNone` | 跨平台 | 直接执行，不使用 shell |
| **default** | `ShellDefault` | 跨平台 | 根据操作系统自动选择 |

---

## 📖 详细文档

Package shellx 定义了shell命令执行库的核心数据类型。本文件定义了Command结构体，集配置、构建、执行于一体的一体化设计。

Command是命令对象的核心实现，支持：
- 配置方法：WithWorkDir、WithEnv、WithTimeout、WithContext等链式调用
- 同步执行：Exec、ExecOutput、ExecStdout、ExecResult
- 异步执行：ExecAsync、Wait
- 进程控制：Kill、Signal、IsRunning、GetPID
- 执行状态管理：IsExecuted（确保命令只执行一次）
- 完整的执行结果：Result对象包含输出、错误、时间、退出码等信息
- 延迟构建：真正的exec.Cmd对象在执行时才创建，确保超时控制精确

Package shellx 提供了一个功能完善、易于使用的Go语言shell命令执行库。

本库基于Go标准库的os/exec包进行封装，提供了更加友好的API和丰富的功能，支持同步和异步命令执行、输入输出重定向、超时控制、上下文管理、多种shell类型支持等功能，并提供类型安全的API和友好的链式调用接口。

## 主要特性

- 支持三种命令创建方式：NewCmd(可变参数)、NewCmds(切片)、NewCmdStr(字符串解析)
- 链式调用API，支持流畅的方法链
- 完整的错误处理和类型安全
- 支持多种shell类型（sh、bash、cmd、powershell、pwsh等）
- 同步和异步执行支持
- 命令执行状态管理和进程控制
- 输入输出重定向和环境变量设置
- 精确的超时控制和上下文取消
- 并发安全的设计
- 跨平台兼容（Windows、Linux、macOS）
- 一体化设计：无需Build()方法，直接执行

## 核心组件

- Command: 命令对象，集配置、构建、执行于一体
- Result: 命令执行结果，包含输出、错误、时间等信息
- ShellType: Shell类型枚举，支持多种shell

## 基本用法

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

// 方式3：使用切片创建命令
result, err := shellx.NewCmds([]string{"git", "status"}).
	WithTimeout(10 * time.Second).
	ExecResult()

if err != nil {
	log.Fatal(err)
}

fmt.Printf("Exit Code: %d\n", result.Code())
fmt.Printf("Success: %t\n", result.Success())
fmt.Printf("Duration: %v\n", result.Duration())
fmt.Printf("Output: %s\n", result.Output())
```

## 便捷函数用法

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

## 高级用法

```go
// 设置标准输入输出
var stdout, stderr bytes.Buffer
stdin := strings.NewReader("input data")

err := shellx.NewCmd("cat").
	WithStdin(stdin).
	WithStdout(&stdout).
	WithStderr(&stderr).
	Exec()

// 使用上下文控制
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

err := shellx.NewCmd("long-running-command").
	WithContext(ctx).
	Exec()

// 进程控制
cmd := shellx.NewCmd("sleep", "100")
cmd.ExecAsync()
pid := cmd.GetPID()
isRunning := cmd.IsRunning()
cmd.Kill() // 或 cmd.Signal(syscall.SIGTERM)
```

## 超时控制

```go
// 方式1：使用WithTimeout方法
err := shellx.NewCmd("sleep", "10").
	WithTimeout(3*time.Second).  // 3秒后超时
	Exec()

// 方式2：使用便捷函数
err := shellx.ExecT(3*time.Second, "sleep", "10")

// 方式3：用户上下文优先（会忽略WithTimeout）
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

err := shellx.NewCmd("sleep", "10").
	WithContext(ctx).
	WithTimeout(3*time.Second).  // 这个会被忽略
	Exec()
```

## 命令解析

```go
// 支持复杂的命令字符串解析，包括引号处理
cmd := shellx.NewCmdStr(`git commit -m "Initial commit with 'quotes'"`)
// 解析结果：["git", "commit", "-m", "Initial commit with 'quotes'"]
```

## Shell类型

```go
// 支持多种shell类型
shellx.ShellSh         // sh shell
shellx.ShellBash       // bash shell
shellx.ShellCmd        // Windows cmd
shellx.ShellPowerShell // Windows PowerShell
shellx.ShellPwsh       // PowerShell Core
shellx.ShellNone       // 直接执行，不使用shell
shellx.ShellDefault    // 根据操作系统自动选择
```

## 注意事项

- 每个Command对象只能执行一次，重复执行会返回错误
- Command是并发安全的，可以在多个goroutine中安全使用
- 命令执行会继承父进程的环境变量，可通过WithEnv添加额外变量
- 超时控制精确：exec.Cmd在执行时才创建，避免配置到执行之间的时间损耗
- 用户上下文优先级高于WithTimeout设置的超时时间
- 异步执行需要调用Wait()等待完成或使用Kill()终止

---

## FUNCTIONS

### func FindCmd

```go
func FindCmd(name string) (string, error)
```

FindCmd 查找命令

**参数:**
- name: 命令名称

**返回:**
- string: 命令路径
- error: 错误信息

### func Exec

```go
func Exec(name string, args ...string) error
```

Exec 执行命令(阻塞)

**参数:**
- name: 命令名
- args: 命令参数

**返回:**
- error: 错误信息

### func ExecStr

```go
func ExecStr(cmdStr string) error
```

ExecStr 执行命令(阻塞)

**参数:**
- cmdStr: 命令字符串

**返回:**
- error: 错误信息

### func ExecOut

```go
func ExecOut(name string, args ...string) ([]byte, error)
```

ExecOut 执行命令并返回合并后的输出(阻塞)

**参数:**
- name: 命令名
- args: 命令参数

**返回:**
- []byte: 输出
- error: 错误信息

### func ExecOutStr

```go
func ExecOutStr(cmdStr string) ([]byte, error)
```

ExecOutStr 执行命令并返回合并后的输出(阻塞)

**参数:**
- cmdStr: 命令字符串

**返回:**
- []byte: 输出
- error: 错误信息

### func ExecT

```go
func ExecT(timeout time.Duration, name string, args ...string) error
```

ExecT 执行命令(阻塞，带超时)

**参数:**
- timeout: 超时时间，如果为0则不设置超时
- name: 命令名
- args: 命令参数

**返回:**
- error: 错误信息

### func ExecStrT

```go
func ExecStrT(timeout time.Duration, cmdStr string) error
```

ExecStrT 执行命令(阻塞，带超时)

**参数:**
- timeout: 超时时间，如果为0则不设置超时
- cmdStr: 命令字符串

**返回:**
- error: 错误信息

### func ExecOutT

```go
func ExecOutT(timeout time.Duration, name string, args ...string) ([]byte, error)
```

ExecOutT 执行命令并返回合并后的输出(阻塞，带超时)

**参数:**
- timeout: 超时时间，如果为0则不设置超时
- name: 命令名
- args: 命令参数

**返回:**
- []byte: 合并后的输出
- error: 错误信息

### func ExecOutStrT

```go
func ExecOutStrT(timeout time.Duration, cmdStr string) ([]byte, error)
```

ExecOutStrT 执行命令并返回合并后的输出(阻塞，带超时)

**参数:**
- timeout: 超时时间，如果为0则不设置超时
- cmdStr: 命令字符串

**返回:**
- []byte: 合并后的输出
- error: 错误信息

### func ParseCmd

```go
func ParseCmd(cmdStr string) []string
```

ParseCmd 将命令字符串解析为命令切片，支持引号处理(单引号、双引号、反引号)，出错时返回空切片

**实现原理：**
1. 去除首尾空白
2. 遍历每个字符
3. 处理引号状态切换
4. 在非引号状态下遇到空格时分割
5. 检查引号是否闭合

**参数:**
- cmdStr: 要解析的命令字符串

**返回值:**
- []string: 解析后的命令切片

## TYPES

### type Command

```go
type Command struct {
	// Has unexported fields.
}
```

Command 命令对象 - 集配置、构建、执行于一体

#### func NewCmd

```go
func NewCmd(name string, args ...string) *Command
```

NewCmd 创建新的命令对象 (数组方式 - 可变参数)

**参数：**
- name: 命令名
- args: 命令参数列表

**返回：**
- *Command: 命令对象

#### func NewCmdStr

```go
func NewCmdStr(cmdStr string) *Command
```

NewCmdStr 创建新的命令对象 (字符串方式)

**参数：**
- cmdStr: 命令字符串

**返回：**
- *Command: 命令对象

#### func NewCmds

```go
func NewCmds(cmdArgs []string) *Command
```

NewCmds 创建新的命令对象 (数组方式 - 切片参数)

**参数：**
- cmdArgs: 命令参数列表，第一个元素为命令名，后续元素为参数

**返回：**
- *Command: 命令对象

#### func (*Command) Args

```go
func (c *Command) Args() []string
```

Args 获取命令参数列表

**返回:**
- []string: 命令参数列表

#### func (*Command) CmdStr

```go
func (c *Command) CmdStr() string
```

CmdStr 获取命令字符串

**返回:**
- string: 命令字符串

**说明:**
- 如果 exec.Cmd 对象已构建，返回其 String() 方法的结果
- 如果 exec.Cmd 对象未构建，返回原始命令字符串
- 该方法可用于调试和日志记录

#### func (*Command) Env

```go
func (c *Command) Env() []string
```

Env 获取命令环境变量列表

**返回:**
- []string: 命令环境变量列表

#### func (*Command) Name

```go
func (c *Command) Name() string
```

Name 获取命令名称

**返回:**
- string: 命令名称

#### func (*Command) Raw

```go
func (c *Command) Raw() string
```

Raw 获取原始命令字符串

**返回:**
- string: 原始命令字符串

#### func (*Command) ShellType

```go
func (c *Command) ShellType() ShellType
```

ShellType 获取shell类型

**返回:**
- ShellType: shell类型

#### func (*Command) Timeout

```go
func (c *Command) Timeout() time.Duration
```

Timeout 获取命令执行超时时间

**返回:**
- time.Duration: 命令执行超时时间

#### func (*Command) WorkDir

```go
func (c *Command) WorkDir() string
```

WorkDir 获取命令执行的工作目录

**返回:**
- string: 命令执行目录

#### func (*Command) WithContext

```go
func (c *Command) WithContext(ctx context.Context) *Command
```

WithContext 设置命令的上下文

**参数：**
- ctx: context.Context类型，用于取消命令执行和超时控制

**返回：**
- *Command: 命令对象

**注意:**
- 该方法会验证上下文是否为空，如果为空则panic.
- 该上下文会覆盖之前设置的超时时间.

#### func (*Command) WithEnv

```go
func (c *Command) WithEnv(key, value string) *Command
```

WithEnv 设置命令的环境变量

**参数：**
- key: 环境变量的键
- value: 环境变量的值

**返回：**
- *Command: 命令对象

**注意:**
- 该方法会验证key是否为空, 如果为空则忽略。
- 无需添加系统环境变量os.Environ(), 系统环境变量会自动继承.

#### func (*Command) WithEnvs

```go
func (c *Command) WithEnvs(envs []string) *Command
```

WithEnvs 批量设置命令的环境变量

**参数：**
- envs: []string类型，环境变量列表，每个元素为"key=value"格式

**返回：**
- *Command: 命令对象

**注意:**
- 该方法会验证环境变量格式，只添加验证通过的环境变量。
- 无需添加系统环境变量os.Environ(), 系统环境变量会自动继承.

#### func (*Command) WithShell

```go
func (c *Command) WithShell(shell ShellType) *Command
```

WithShell 设置命令的shell类型

**参数：**
- shell: ShellType类型，表示要使用的shell类型

**返回：**
- *Command: 命令对象

#### func (*Command) WithStderr

```go
func (c *Command) WithStderr(stderr io.Writer) *Command
```

WithStderr 设置命令的标准错误输出

**参数：**
- stderr: io.Writer类型，用于接收命令的标准错误输出

**返回：**
- *Command: 命令对象

#### func (*Command) WithStdin

```go
func (c *Command) WithStdin(stdin io.Reader) *Command
```

WithStdin 设置命令的标准输入

**参数：**
- stdin: io.Reader类型，用于提供命令的标准输入

**返回：**
- *Command: 命令对象

#### func (*Command) WithStdout

```go
func (c *Command) WithStdout(stdout io.Writer) *Command
```

WithStdout 设置命令的标准输出

**参数：**
- stdout: io.Writer类型，用于接收命令的标准输出

**返回：**
- *Command: 命令对象

#### func (*Command) WithTimeout

```go
func (c *Command) WithTimeout(timeout time.Duration) *Command
```

WithTimeout 设置命令的超时时间(便捷方式)

**参数：**
- timeout: time.Duration类型，命令执行的超时时间

**返回：**
- *Command: 命令对象

**注意:**
- 该方法会验证超时时间是否小于等于0, 如果小于等于0则忽略。
- 该超时时间优先级低于上下文设置的超时时间.

#### func (*Command) WithWorkDir

```go
func (c *Command) WithWorkDir(dir string) *Command
```

WithWorkDir 设置命令的工作目录

**参数：**
- dir: 命令的工作目录

**返回：**
- *Command: 命令对象

#### func (*Command) Cmd

```go
func (c *Command) Cmd() *exec.Cmd
```

Cmd 获取底层的 exec.Cmd 对象

**返回:**
- *exec.Cmd: 底层的 exec.Cmd 对象

#### func (*Command) Exec

```go
func (c *Command) Exec() error
```

Exec 执行命令(阻塞)

**返回:**
- error: 错误信息

#### func (*Command) ExecAsync

```go
func (c *Command) ExecAsync() error
```

ExecAsync 异步执行命令(非阻塞)

**返回:**
- error: 错误信息

#### func (*Command) ExecOutput

```go
func (c *Command) ExecOutput() ([]byte, error)
```

ExecOutput 执行命令并返回合并后的输出(阻塞)

**返回:**
- []byte: 命令输出
- error: 错误信息

#### func (*Command) ExecResult

```go
func (c *Command) ExecResult() (*Result, error)
```

ExecResult 执行命令并返回完整的执行结果(阻塞)

**使用示例:**

```go
result, err := cmd.ExecResult()
if err != nil {
    // 处理错误情况
    log.Printf("Command failed: %v", err)
    return
}
// 处理成功情况
fmt.Println(string(result.Output()))
```

**返回:**
- *Result: 执行结果对象，包含输出、时间、退出码等信息
- error: 执行过程中的错误信息

#### func (*Command) ExecStdout

```go
func (c *Command) ExecStdout() ([]byte, error)
```

ExecStdout 执行命令并返回标准输出(阻塞)

**返回:**
- []byte: 标准输出
- error: 错误信息

#### func (*Command) GetPID

```go
func (c *Command) GetPID() int
```

GetPID 获取进程ID

**返回:**
- int: 进程ID，如果进程不存在返回0

#### func (*Command) IsExecuted

```go
func (c *Command) IsExecuted() bool
```

IsExecuted 检查命令是否已经执行过

**返回:**
- bool: 是否已执行

#### func (*Command) IsRunning

```go
func (c *Command) IsRunning() bool
```

IsRunning 检查进程是否还在运行

**返回:**
- bool: 是否在运行

#### func (*Command) Kill

```go
func (c *Command) Kill() error
```

Kill 杀死当前命令的进程

**返回:**
- error: 错误信息

#### func (*Command) Signal

```go
func (c *Command) Signal(sig os.Signal) error
```

Signal 向当前进程发送信号

**参数:**
- sig: 信号类型

**返回:**
- error: 错误信息

#### func (*Command) Wait

```go
func (c *Command) Wait() error
```

Wait 等待命令执行完成(仅在异步执行时有效)

**返回:**
- error: 错误信息

### type Result

```go
type Result struct {
	// Has unexported fields.
}
```

Result 表示命令执行的结果

#### func (*Result) Code

```go
func (r *Result) Code() int
```

Code 获取命令退出码

#### func (*Result) Duration

```go
func (r *Result) Duration() time.Duration
```

Duration 获取命令执行时长

#### func (*Result) End

```go
func (r *Result) End() time.Time
```

End 获取命令结束时间

#### func (*Result) Output

```go
func (r *Result) Output() []byte
```

Output 获取命令输出

#### func (*Result) Start

```go
func (r *Result) Start() time.Time
```

Start 获取命令开始时间

#### func (*Result) Success

```go
func (r *Result) Success() bool
```

Success 获取命令是否执行成功

### type ShellType

```go
type ShellType int
```

ShellType 定义shell类型

```go
const (
	ShellSh         ShellType = iota // sh shell
	ShellBash                        // bash shell
	ShellPwsh                        // pwsh (PowerShell Core)
	ShellPowerShell                  // powershell (Windows PowerShell)
	ShellCmd                         // cmd (Windows Command Prompt)
	ShellNone                        // 无shell, 直接原生的执行命令
	ShellDefault                     // 默认shell, 根据操作系统自动选择(Windows系统默认为cmd, 其他系统默认为sh)
)
```

#### func (ShellType) String

```go
func (s ShellType) String() string
```

String 返回shell类型的字符串表示