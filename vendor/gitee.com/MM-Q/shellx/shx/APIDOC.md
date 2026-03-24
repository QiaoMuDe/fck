# shx API 文档

## 包信息

```go
package shx // import "gitee.com/MM-Q/shellx/shx"
```

## 概述

Package shx 提供了基于 mvdan.cc/sh/v3 的纯 Go shell 命令执行功能。

本包是 ShellX 的子包, 提供了与主包相似的 API 风格, 但具有更好的跨平台一致性。 它使用 mvdan.cc/sh/v3 进行命令执行, 不依赖系统 shell。

### 主要特性

- 纯 Go 实现, 不依赖系统 shell
- 更好的跨平台一致性 (Windows/Linux/macOS 行为一致)
- 链式调用 API, 支持流畅的方法链
- 支持超时控制和上下文取消
- 最小并发保护 (使用 atomic.Bool 防止重复执行)

### 基本用法

```go
import "gitee.com/MM-Q/shellx/shx"

// 简单执行
err := shx.Exec("echo hello world")

// 获取输出
output, err := shx.Output("ls -la")

// 链式配置
output, err := shx.New("echo hello").
    WithTimeout(5 * time.Second).
    WithDir("/tmp").
    ExecOutput()

// 使用上下文
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
err := shx.New("long-running-command").WithContext(ctx).Exec()
```

### 注意事项

- Shx 对象的配置方法 (WithXxx) 不是并发安全的, 不要在多个 goroutine 中并发配置
- 每个 Shx 对象只能执行一次, 重复执行会返回错误
- mvdan/sh 是同步执行的, 不提供异步 API, 如需异步请使用 goroutine 包装
- 不支持进程控制 (无 PID、Kill、Signal) , 只能通过 context 取消

---

## 变量

```go
var (
    // ErrAlreadyExecuted 表示命令已经执行过
    ErrAlreadyExecuted = errors.New("command has already been executed")

    // ErrNilContext 表示上下文为 nil
    ErrNilContext = errors.New("context cannot be nil")

    // ErrNilReader 表示 reader 为 nil
    ErrNilReader = errors.New("reader cannot be nil")

    // ErrNilWriter 表示 writer 为 nil
    ErrNilWriter = errors.New("writer cannot be nil")
)
```

---

## 函数

### IsExitStatus

```go
func IsExitStatus(err error) (uint8, bool)
```

检查错误是否是退出状态错误

**参数：**
- `err`: 错误对象

**返回：**
- `uint8`: 退出码
- `bool`: 是否是退出状态错误

---

### Out

```go
func Out(cmd string) ([]byte, error)
```

执行并获取输出

**参数：**
- `cmd`: 命令字符串

**返回：**
- `[]byte`: 命令输出
- `error`: 执行错误

**示例：**

```go
output, err := shx.Out("ls -la")
```

---

### OutCtx

```go
func OutCtx(ctx context.Context, cmd string) ([]byte, error)
```

使用上下文执行并获取输出

**参数：**
- `ctx`: 上下文
- `cmd`: 命令字符串

**返回：**
- `[]byte`: 命令输出
- `error`: 执行错误

**示例：**

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
output, err := shx.OutCtx(ctx, "ls -la")
```

---

### OutWith

```go
func OutWith(cmd string, timeout time.Duration) ([]byte, error)
```

超时执行并获取输出

**参数：**
- `cmd`: 命令字符串
- `timeout`: 超时时间

**返回：**
- `[]byte`: 命令输出
- `error`: 执行错误

**示例：**

```go
output, err := shx.OutWith("sleep 5", 10*time.Second)
```

---

### OutWithIO

```go
func OutWithIO(cmd string, stdin io.Reader, stdout, stderr io.Writer) ([]byte, error)
```

使用自定义输入输出执行并获取输出

**参数：**
- `cmd`: 命令字符串
- `stdin`: 标准输入
- `stdout`: 标准输出
- `stderr`: 标准错误

**返回：**
- `[]byte`: 命令输出
- `error`: 执行错误

**示例：**

```go
var buf bytes.Buffer
output, err := shx.OutWithIO("cat", strings.NewReader("hello"), &buf, os.Stderr)
```

---

### Run

```go
func Run(cmd string) error
```

执行命令

**参数：**
- `cmd`: 命令字符串

**返回：**
- `error`: 执行错误

**示例：**

```go
err := shx.Run("echo hello")
```

---

### RunCtx

```go
func RunCtx(ctx context.Context, cmd string) error
```

使用上下文执行

**参数：**
- `ctx`: 上下文
- `cmd`: 命令字符串

**返回：**
- `error`: 执行错误

**示例：**

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
err := shx.RunCtx(ctx, "sleep 10")
```

---

### RunToTerminal

```go
func RunToTerminal(cmd string) error
```

执行命令并输出到终端

**参数：**
- `cmd`: 命令字符串

**返回：**
- `error`: 执行错误

**示例：**

```go
err := shx.RunToTerminal("echo hello")
```

---

### RunWith

```go
func RunWith(cmd string, timeout time.Duration) error
```

超时执行

**参数：**
- `cmd`: 命令字符串
- `timeout`: 超时时间

**返回：**
- `error`: 执行错误

**示例：**

```go
err := shx.RunWith("sleep 10", 5*time.Second)
```

---

### RunWithIO

```go
func RunWithIO(cmd string, stdin io.Reader, stdout, stderr io.Writer) error
```

使用自定义输入输出执行

**参数：**
- `cmd`: 命令字符串
- `stdin`: 标准输入
- `stdout`: 标准输出
- `stderr`: 标准错误

**返回：**
- `error`: 执行错误

**示例：**

```go
var buf bytes.Buffer
err := shx.RunWithIO("cat", strings.NewReader("hello"), &buf, os.Stderr)
```

---

## 类型

### ExitStatus

```go
type ExitStatus struct {
    Code uint8
}
```

ExitStatus 包装退出状态错误

---

#### Error

```go
func (e ExitStatus) Error() string
```

Error 实现 error 接口

---

### Shx

```go
type Shx struct {
    // Has unexported fields.
}
```

Shx 表示一个待执行的 shell 命令

---

#### New

```go
func New(cmdStr string) *Shx
```

从字符串创建命令

**参数:**
- `cmdStr`: 命令字符串

**返回:**
- `*Shx`: 命令对象

**示例:**

```go
cmd := shx.New("echo hello world")
cmd := shx.New("ls -la | grep .go")
```

---

#### NewArgs

```go
func NewArgs(cmd string, args ...string) *Shx
```

从命令名和可变参数创建命令

**参数:**
- `cmd`: 命令名
- `args`: 可变参数列表

**返回:**
- `*Shx`: 命令对象

**示例:**

```go
cmd := shx.NewArgs("ls", "-la", "/tmp")
cmd := shx.NewArgs("git", "commit", "-m", "message")
```

---

#### NewCmds

```go
func NewCmds(cmds []string) *Shx
```

从命令切片创建命令

**参数:**
- `cmds`: 命令切片，每个元素是一个完整的命令部分

**返回:**
- `*Shx`: 命令对象

**示例:**

```go
cmd := shx.NewCmds([]string{"ls", "-la", "|", "grep", ".go"})
cmd := shx.NewCmds([]string{"echo", "hello", ">", "output.txt"})
```

---

#### NewWithParser

```go
func NewWithParser(cmdStr string, parser *syntax.Parser) *Shx
```

使用自定义解析器创建命令

**参数:**
- `cmdStr`: 命令字符串
- `parser`: 自定义解析器

**返回:**
- `*Shx`: 命令对象

---

#### Context

```go
func (s *Shx) Context() context.Context
```

获取上下文

**返回:**
- `context.Context`: 上下文

---

#### Dir

```go
func (s *Shx) Dir() string
```

获取工作目录

**返回:**
- `string`: 工作目录

---

#### Env

```go
func (s *Shx) Env() expand.Environ
```

获取环境变量

**返回:**
- `expand.Environ`: 环境变量

---

#### Exec

```go
func (s *Shx) Exec() error
```

执行命令 (阻塞)

**返回:**
- `error`: 执行过程中的错误, 不包含退出码错误

**线程安全:**
- 使用 atomic.Bool 确保重复执行检测的线程安全

**示例:**

```go
err := shx.New("echo hello").Exec()
if err != nil {
    log.Fatal(err)
}
```

---

#### ExecContext

```go
func (s *Shx) ExecContext(ctx context.Context) error
```

在指定上下文中执行命令

**参数:**
- `ctx`: 上下文 (用于取消执行)

**返回:**
- `error`: 执行过程中的错误

**注意:**
- 此方法会覆盖之前通过 WithContext 设置的上下文
- 此方法不受 WithTimeout 影响 (上下文的超时优先)

---

#### ExecContextOutput

```go
func (s *Shx) ExecContextOutput(ctx context.Context) ([]byte, error)
```

在指定上下文中执行并返回输出

**参数:**
- `ctx`: 上下文

**返回:**
- `[]byte`: 命令输出
- `error`: 执行过程中的错误

---

#### ExecOutput

```go
func (s *Shx) ExecOutput() ([]byte, error)
```

执行命令并返回输出

**返回:**
- `[]byte`: 命令输出 (stdout 和 stderr 合并)
- `error`: 执行过程中的错误

**注意:**
- 内部会自动捕获 stdout 和 stderr
- 如果需要区分 stdout 和 stderr, 请使用 WithStdout 和 WithStderr 自定义

---

#### IsExecuted

```go
func (s *Shx) IsExecuted() bool
```

检查命令是否已经执行过

**返回:**
- `bool`: 是否已执行

---

#### Raw

```go
func (s *Shx) Raw() string
```

获取原始命令字符串

**返回:**
- `string`: 原始命令字符串

---

#### Timeout

```go
func (s *Shx) Timeout() time.Duration
```

获取超时时间

**返回:**
- `time.Duration`: 超时时间

---

#### WithContext

```go
func (s *Shx) WithContext(ctx context.Context) *Shx
```

设置上下文

**参数:**
- `ctx`: 上下文

**返回:**
- `*Shx`: 命令对象 (支持链式调用)

**注意:**
- 如果命令已经执行过, 会 panic
- 如果 ctx 为 nil, 会 panic
- 设置的上下文会完全覆盖 WithTimeout 设置的超时

---

#### WithDir

```go
func (s *Shx) WithDir(dir string) *Shx
```

设置工作目录

**参数:**
- `dir`: 工作目录路径

**返回:**
- `*Shx`: 命令对象 (支持链式调用)

**注意:**
- 如果命令已经执行过, 会 panic
- 如果目录不存在或不是目录, 会 panic

---

#### WithEnv

```go
func (s *Shx) WithEnv(key, value string) *Shx
```

设置环境变量

**参数:**
- `key`: 环境变量名
- `value`: 环境变量值

**返回:**
- `*Shx`: 命令对象 (支持链式调用)

**注意:**
- 如果命令已经执行过, 会 panic
- 如果 key 为空, 则忽略

---

#### WithEnvMap

```go
func (s *Shx) WithEnvMap(envs map[string]string) *Shx
```

批量设置环境变量 (map 方式)

**参数:**
- `envs`: 环境变量映射 (key-value)

**返回:**
- `*Shx`: 命令对象 (支持链式调用)

**注意:**
- 如果命令已经执行过, 会 panic
- 同名的变量, 后遍历的会覆盖先遍历的 (map 遍历顺序不确定)

**示例:**

```go
cmd := shx.New("echo $FOO $BAR").
    WithEnvMap(map[string]string{
        "FOO": "hello",
        "BAR": "world",
    })
```

---

#### WithEnvs

```go
func (s *Shx) WithEnvs(envs []string) *Shx
```

批量设置环境变量 (切片方式)

**参数:**
- `envs`: 环境变量切片, 每个元素格式为 `"key=value"`

**返回:**
- `*Shx`: 命令对象 (支持链式调用)

**注意:**
- 如果命令已经执行过, 会 panic
- 格式错误的项会被忽略 (缺少 `=` 或 key 为空)
- 同名的变量, 后出现的会覆盖先出现的
- 新变量会覆盖旧环境中的同名变量

**示例:**

```go
// 使用字符串切片
cmd := shx.New("echo $PATH").
    WithEnvs([]string{
        "PATH=/usr/local/bin:/usr/bin",
        "HOME=/home/user",
    })

// 配合 os.Environ() 使用
envs := os.Environ()
envs = append(envs, "CUSTOM_VAR=value")
cmd := shx.New("env").WithEnvs(envs)
```

---

#### WithStderr

```go
func (s *Shx) WithStderr(w io.Writer) *Shx
```

设置标准错误

**参数:**
- `w`: 错误输出写入器

**返回:**
- `*Shx`: 命令对象 (支持链式调用)

**注意:**
- 如果命令已经执行过, 会 panic
- 如果 w 为 nil, 会 panic

---

#### WithStdin

```go
func (s *Shx) WithStdin(r io.Reader) *Shx
```

设置标准输入

**参数:**
- `r`: 输入读取器

**返回:**
- `*Shx`: 命令对象 (支持链式调用)

**注意:**
- 如果命令已经执行过, 会 panic
- 如果 r 为 nil, 会 panic

---

#### WithStdout

```go
func (s *Shx) WithStdout(w io.Writer) *Shx
```

设置标准输出

**参数:**
- `w`: 输出写入器

**返回:**
- `*Shx`: 命令对象 (支持链式调用)

**注意:**
- 如果命令已经执行过, 会 panic
- 如果 w 为 nil, 会 panic

---

#### WithTimeout

```go
func (s *Shx) WithTimeout(d time.Duration) *Shx
```

设置超时时间

**参数:**
- `d`: 超时时间

**返回:**
- `*Shx`: 命令对象 (支持链式调用)

**注意:**
- 如果命令已经执行过, 会 panic
- 如果 d <= 0, 则忽略 (不设置超时)
