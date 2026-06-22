# Package shx

`import "gitee.com/MM-Q/shx"`

Package shx 提供了基于 mvdan.cc/sh/v3 的纯 Go shell 命令执行功能。

它使用 mvdan.cc/sh/v3 进行命令解析和执行, 不依赖系统 shell。

## 主要特性

- 纯 Go 实现, 不依赖系统 shell
- 更好的跨平台一致性 (Windows/Linux/macOS 行为一致)
- 链式调用 API, 支持流畅的方法链
- 支持超时控制和上下文取消
- 支持执行 .sh 脚本文件

## 基本用法

```go
import "gitee.com/MM-Q/shx"

// 简单执行
err := shx.Run("echo hello world")

// 获取输出
output, err := shx.Out("ls -la")

// 链式配置
output, err := shx.New("echo hello").
	WithTimeout(5 * time.Second).
	WithDir("/tmp").
	ExecOutput()

// 使用上下文
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
err := shx.New("long-running-command").WithContext(ctx).Exec()

// 执行脚本文件
err = shx.RunScript("deploy.sh")
```

## 注意事项

- Shx 对象的配置方法 (WithXxx) 不是并发安全的, 不要在多个 goroutine 中并发配置
- mvdan/sh 是同步执行的, 不提供异步 API, 如需异步请使用 goroutine 包装
- 不支持进程控制 (无 PID、Kill、Signal) , 只能通过 context 取消

---

## Variables

### ErrNilContext

```go
var (
	// ErrNilContext 表示上下文为 nil
	ErrNilContext = errors.New("context cannot be nil")
)
```

预定义错误

---

## Functions

### CheckScriptSyntax

```go
func CheckScriptSyntax(filePath string) error
```

CheckScriptSyntax 检查 shell 脚本文件的语法

**参数：**

- `filePath`: 脚本文件路径

**返回：**

- `error`: 语法正确返回 nil；语法错误返回 `*SyntaxError`；文件不存在等系统错误返回原始错误

**示例：**

```go
err := shx.CheckScriptSyntax("deploy.sh")
if err != nil {
    fmt.Println(err)
}
```

---

### CheckSyntax

```go
func CheckSyntax(script string) error
```

CheckSyntax 检查 shell 命令字符串的语法

**参数：**

- `script`: shell 命令字符串

**返回：**

- `error`: 语法正确返回 nil；语法错误返回 `*SyntaxError`

**示例：**

```go
err := shx.CheckSyntax("echo hello")
if err != nil {
    fmt.Println(err)
}
```

---

### FindCmd

```go
func FindCmd(name string) (string, error)
```

FindCmd 查找命令

增强版，在标准库 exec.LookPath 基础上增加了以下能力：
- 处理 Go 1.19+ 的 ErrDot 安全限制（当前目录程序）
- 返回绝对路径
- Windows 上检查可执行文件扩展名

**参数：**

- `name`: 命令名称

**返回：**

- `string`: 命令的绝对路径
- `error`: 错误信息

---

### FindCommandPath

```go
func FindCommandPath(name string) string
```

FindCommandPath 查找单个命令的绝对路径

供其他包复用，只返回第一个匹配的路径 优先使用标准库 exec.LookPath, 处理 ErrDot 情况，找不到则返回空字符串

**参数：**

- `name`: 命令名称

**返回：**

- `string`: 命令的绝对路径，如果找不到则返回空字符串

---

### IsExitStatus

```go
func IsExitStatus(err error) (uint8, bool)
```

IsExitStatus 检查错误是否是退出状态错误

**参数：**

- `err`: 错误对象

**返回：**

- `uint8`: 退出码
- `bool`: 是否是退出状态错误

---

### Format

```go
func Format(script string) (string, error)
```

Format 使用默认选项格式化 shell 命令字符串。

默认启用：缩进 4 空格、case 语句缩进、注释保留。

**参数：**

- `script`: shell 命令字符串

**返回：**

- `string`: 格式化后的字符串
- `error`: 解析错误或系统错误

**示例：**

```go
formatted, err := shx.Format("for i in 1 2 3;do echo $i;done")
if err != nil {
    log.Fatal(err)
}
fmt.Println(formatted)
```

---

### FormatScript

```go
func FormatScript(filePath string) (string, error)
```

FormatScript 格式化 shell 脚本文件

**参数：**

- `filePath`: 脚本文件路径

**返回：**

- `string`: 格式化后的脚本内容
- `error`: 解析错误或系统错误

**示例：**

```go
formatted, err := shx.FormatScript("deploy.sh")
if err != nil {
    log.Fatal(err)
}
fmt.Println(formatted)
```

---

### Out

```go
func Out(cmd string) ([]byte, error)
```

Out 执行并获取输出

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

OutCtx 使用上下文执行并获取输出

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

### OutCtxScript

```go
func OutCtxScript(ctx context.Context, filePath string) ([]byte, error)
```

OutCtxScript 使用上下文执行 bash 脚本文件并获取输出

**参数：**

- `ctx`: 上下文
- `filePath`: 脚本文件路径

**返回：**

- `[]byte`: 命令输出
- `error`: 执行错误

**示例：**

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
output, err := shx.OutCtxScript(ctx, "build.sh")
```

---

### OutScript

```go
func OutScript(filePath string) ([]byte, error)
```

OutScript 执行 bash 脚本文件并获取输出

**参数：**

- `filePath`: 脚本文件路径

**返回：**

- `[]byte`: 命令输出
- `error`: 执行错误

**示例：**

```go
output, err := shx.OutScript("deploy.sh")
```

---

### OutScriptWith

```go
func OutScriptWith(filePath string, timeout time.Duration) ([]byte, error)
```

OutScriptWith 超时执行 bash 脚本文件并获取输出

**参数：**

- `filePath`: 脚本文件路径
- `timeout`: 超时时间

**返回：**

- `[]byte`: 命令输出
- `error`: 执行错误

**示例：**

```go
output, err := shx.OutScriptWith("build.sh", 60*time.Second)
```

---

### OutScriptWithIO

```go
func OutScriptWithIO(filePath string, stdin io.Reader, stdout, stderr io.Writer) ([]byte, error)
```

OutScriptWithIO 使用自定义输入输出执行 bash 脚本文件并获取输出

**参数：**

- `filePath`: 脚本文件路径
- `stdin`: 标准输入
- `stdout`: 标准输出
- `stderr`: 标准错误

**返回：**

- `[]byte`: 命令输出
- `error`: 执行错误

**示例：**

```go
var buf bytes.Buffer
output, err := shx.OutScriptWithIO("script.sh", strings.NewReader("input"), &buf, os.Stderr)
```

---

### OutWith

```go
func OutWith(cmd string, timeout time.Duration) ([]byte, error)
```

OutWith 超时执行并获取输出

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

OutWithIO 使用自定义输入输出执行并获取输出

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

Run 执行命令

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

RunCtx 使用上下文执行

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

### RunCtxScript

```go
func RunCtxScript(ctx context.Context, filePath string) error
```

RunCtxScript 使用上下文执行 bash 脚本文件

**参数：**

- `ctx`: 上下文
- `filePath`: 脚本文件路径

**返回：**

- `error`: 执行错误

**示例：**

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
err := shx.RunCtxScript(ctx, "long_task.sh")
```

---

### RunScript

```go
func RunScript(filePath string) error
```

RunScript 执行 bash 脚本文件

**参数：**

- `filePath`: 脚本文件路径

**返回：**

- `error`: 执行错误

**示例：**

```go
err := shx.RunScript("deploy.sh")
```

---

### RunScriptToTerminal

```go
func RunScriptToTerminal(filePath string) error
```

RunScriptToTerminal 执行 bash 脚本文件并输出到终端

**参数：**

- `filePath`: 脚本文件路径

**返回：**

- `error`: 执行错误

**示例：**

```go
err := shx.RunScriptToTerminal("deploy.sh")
```

---

### RunScriptWith

```go
func RunScriptWith(filePath string, timeout time.Duration) error
```

RunScriptWith 超时执行 bash 脚本文件

**参数：**

- `filePath`: 脚本文件路径
- `timeout`: 超时时间

**返回：**

- `error`: 执行错误

**示例：**

```go
err := shx.RunScriptWith("long_task.sh", 30*time.Second)
```

---

### RunScriptWithIO

```go
func RunScriptWithIO(filePath string, stdin io.Reader, stdout, stderr io.Writer) error
```

RunScriptWithIO 使用自定义输入输出执行 bash 脚本文件

**参数：**

- `filePath`: 脚本文件路径
- `stdin`: 标准输入
- `stdout`: 标准输出
- `stderr`: 标准错误

**返回：**

- `error`: 执行错误

**示例：**

```go
var buf bytes.Buffer
err := shx.RunScriptWithIO("script.sh", strings.NewReader("input"), &buf, os.Stderr)
```

---

### RunToTerminal

```go
func RunToTerminal(cmd string) error
```

RunToTerminal 执行命令并输出到终端

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

RunWith 超时执行

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

RunWithIO 使用自定义输入输出执行

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

### Split

```go
func Split(cmdStr string) []string
```

Split 将命令字符串拆分为命令切片，支持引号处理(单引号、双引号、反引号)

**功能：**

- 智能拆分 Shell 命令字符串为参数数组
- 支持单引号、双引号、反引号包裹的内容
- 正确处理转义字符和特殊字符
- 自动处理命令分隔符(;|&&||)

**参数：**

- `cmdStr`: 要拆分的命令字符串

**返回值：**

- `[]string`: 拆分后的命令切片 (最佳结果)

**注意：**

- 此函数忽略拆分错误，返回最佳拆分结果。如需错误信息，请使用 SplitE 函数。
- 转义字符保持原样，不进行解释处理
- 支持多字符操作符如 &&、||、>>、<< 等

---

### SplitE

```go
func SplitE(cmdStr string) ([]string, error)
```

SplitE 将命令字符串拆分为命令切片（带错误信息），支持引号处理(单引号、双引号、反引号)

**功能：**

- 智能拆分 Shell 命令字符串为参数数组
- 支持单引号、双引号、反引号包裹的内容
- 正确处理转义字符和特殊字符
- 自动处理命令分隔符(;|&&||)
- 检测并返回拆分过程中的错误

**参数：**

- `cmdStr`: 要拆分的命令字符串

**返回值：**

- `[]string`: 拆分后的命令切片
- `error`: 拆分错误，成功时为 nil

**错误类型：**

- `UnclosedQuoteError`: 未闭合的引号错误
- 其他可能的语法错误

**注意：**

- 转义字符保持原样，不进行解释处理
- 支持多字符操作符如 &&、||、>>、<< 等

---

## Types

### type ExitStatus

```go
type ExitStatus struct {
	Code uint8
	// 内含未导出字段（err error，保存原始错误用于错误链）
}
```

ExitStatus 包装退出状态错误。
支持错误链穿透：可通过 `errors.Is(err, interp.ExitStatus(N))` 检测退出码。

#### func (ExitStatus) Error

```go
func (e ExitStatus) Error() string
```

Error 实现 error 接口

#### func (ExitStatus) Unwrap

```go
func (e ExitStatus) Unwrap() error
```

Unwrap 返回原始错误，支持 `errors.Is`/`errors.As` 遍历错误链。
当 `ExitStatus` 由 `handleError` 构造时，返回原始的 `interp.ExitStatus` 错误；
当通过结构体字面量 `ExitStatus{Code: N}` 直接创建时，返回 nil。

---

### type Shx

```go
type Shx struct {
	// 导出字段（可通过结构体字面量或 WithXxx 方法配置）
	Env     expand.Environ  // 环境变量
	Dir     string          // 工作目录
	Timeout time.Duration   // 超时时间
	Ctx     context.Context // 上下文
	Stdin   io.Reader       // 标准输入
	Stdout  io.Writer       // 标准输出
	Stderr  io.Writer       // 标准错误

	// Has unexported fields.
}
```

Shx 表示一个待执行的 shell 命令或脚本文件

支持两种输入方式：
- 命令字符串（通过 New/NewArgs/NewCmds 创建）
- bash 脚本文件（通过 NewScript 创建）

字段可通过结构体字面量或 WithXxx 方法配置：

```go
cmd := &Shx{
	Dir:     "/tmp",
	Timeout: 5 * time.Second,
}
```

#### func New

```go
func New(cmdStr string) *Shx
```

New 从字符串创建命令

**参数：**

- `cmdStr`: 命令字符串

**返回：**

- `*Shx`: 命令对象

**示例：**

```go
cmd := shx.New("echo hello world")
cmd := shx.New("ls -la | grep .go")
```

#### func NewArgs

```go
func NewArgs(cmd string, args ...string) *Shx
```

NewArgs 从命令名和可变参数创建命令

**参数：**

- `cmd`: 命令名
- `args`: 可变参数列表

**返回：**

- `*Shx`: 命令对象

**示例：**

```go
cmd := shx.NewArgs("ls", "-la", "/tmp")
cmd := shx.NewArgs("git", "commit", "-m", "message")
```

#### func NewCmds

```go
func NewCmds(cmds []string) *Shx
```

NewCmds 从命令切片创建命令

**参数：**

- `cmds`: 命令切片，每个元素是一个完整的命令部分

**返回：**

- `*Shx`: 命令对象

**示例：**

```go
cmd := shx.NewCmds([]string{"ls", "-la", "|", "grep", ".go"})
cmd := shx.NewCmds([]string{"echo", "hello", ">", "output.txt"})
```

#### func NewScript

```go
func NewScript(filePath string) *Shx
```

NewScript 从 bash 脚本文件创建命令

**参数：**

- `filePath`: 脚本文件路径，必须以 .sh 结尾

**返回：**

- `*Shx`: 命令对象

**示例：**

```go
cmd := shx.NewScript("deploy.sh")
cmd := shx.NewScript("/path/to/script.sh")
```

**注意：**

- 如果 filePath 为空或不是 .sh 后缀, 会 panic

#### func (s *Shx) Exec

```go
func (s *Shx) Exec() error
```

Exec 执行命令

**返回：**

- `error`: 执行过程中的错误, 不包含退出码错误

**示例：**

```go
err := shx.New("echo hello").Exec()
if err != nil {
	log.Fatal(err)
}
```

#### func (s *Shx) ExecContext

```go
func (s *Shx) ExecContext(ctx context.Context) error
```

ExecContext 在指定上下文中执行命令

**参数：**

- `ctx`: 上下文 (用于取消执行)

**返回：**

- `error`: 执行过程中的错误

**注意：**

- 此方法会覆盖之前通过 WithContext 设置的上下文
- 此方法不受 WithTimeout 影响 (上下文的超时优先)

#### func (s *Shx) ExecContextOutput

```go
func (s *Shx) ExecContextOutput(ctx context.Context) ([]byte, error)
```

ExecContextOutput 在指定上下文中执行并返回输出

**参数：**

- `ctx`: 上下文

**返回：**

- `[]byte`: 命令输出
- `error`: 执行过程中的错误

#### func (s *Shx) ExecOutput

```go
func (s *Shx) ExecOutput() ([]byte, error)
```

ExecOutput 执行命令并返回输出

**返回：**

- `[]byte`: 命令输出 (stdout 和 stderr 合并)
- `error`: 执行过程中的错误

**注意：**

- 内部会自动捕获 stdout 和 stderr
- 如果需要区分 stdout 和 stderr, 请使用 WithStdout 和 WithStderr 自定义

#### func (s *Shx) Raw

```go
func (s *Shx) Raw() string
```

Raw 获取原始命令字符串

**返回：**

- `string`: 原始命令字符串

#### func (s *Shx) WithContext

```go
func (s *Shx) WithContext(ctx context.Context) *Shx
```

WithContext 设置上下文

**参数：**

- `ctx`: 上下文

**返回：**

- `*Shx`: 命令对象 (支持链式调用)

**注意：**

- 设置的上下文会完全覆盖 WithTimeout 设置的超时

#### func (s *Shx) WithDir

```go
func (s *Shx) WithDir(dir string) *Shx
```

WithDir 设置工作目录

**参数：**

- `dir`: 工作目录路径

**返回：**

- `*Shx`: 命令对象 (支持链式调用)

**注意：**

- 如果目录不存在或不是目录, 会 panic

#### func (s *Shx) WithEnv

```go
func (s *Shx) WithEnv(key, value string) *Shx
```

WithEnv 设置环境变量

**参数：**

- `key`: 环境变量名
- `value`: 环境变量值

**返回：**

- `*Shx`: 命令对象 (支持链式调用)

**注意：**

- 如果 key 为空, 会 panic

#### func (s *Shx) WithEnvs

```go
func (s *Shx) WithEnvs(envs []string) *Shx
```

WithEnvs 批量设置环境变量

**参数：**

- `envs`: 环境变量切片, 每个元素格式为 "key=value"

**返回：**

- `*Shx`: 命令对象 (支持链式调用)

**注意：**

- 格式错误的项会被忽略
- 同名的变量, 后出现的会覆盖先出现的

#### func (s *Shx) WithStderr

```go
func (s *Shx) WithStderr(w io.Writer) *Shx
```

WithStderr 设置标准错误

**参数：**

- `w`: 错误输出写入器

**返回：**

- `*Shx`: 命令对象 (支持链式调用)

#### func (s *Shx) WithStdin

```go
func (s *Shx) WithStdin(r io.Reader) *Shx
```

WithStdin 设置标准输入

**参数：**

- `r`: 输入读取器

**返回：**

- `*Shx`: 命令对象 (支持链式调用)

#### func (s *Shx) WithStdout

```go
func (s *Shx) WithStdout(w io.Writer) *Shx
```

WithStdout 设置标准输出

**参数：**

- `w`: 输出写入器

**返回：**

- `*Shx`: 命令对象 (支持链式调用)

#### func (s *Shx) WithTimeout

```go
func (s *Shx) WithTimeout(d time.Duration) *Shx
```

WithTimeout 设置超时时间

**参数：**

- `d`: 超时时间

**返回：**

- `*Shx`: 命令对象 (支持链式调用)

**注意：**

- 如果 d <= 0, 则忽略 (不设置超时)

---

### type FormatOptions

```go
type FormatOptions struct {
    Indent            uint // 缩进空格数（0 表示使用 tab）
    SwitchCaseIndent  bool // case 语句体是否缩进
    KeepComments      bool // 是否保留注释
    BinaryNextLine    bool // &&、|| 等二元操作符是否换行显示
    FunctionNextLine  bool // 函数体 { 是否换行
    SpaceRedirects    bool // 重定向符前后是否加空格
    SingleLine        bool // 是否单行输出
    Minify            bool // 是否最小化输出（压缩模式）
}
```

FormatOptions 控制 shell 脚本格式化的行为选项，传递给 `FormatWithOptions`/`FormatScriptWithOptions`。

#### func DefaultFormatOptions

```go
func DefaultFormatOptions() FormatOptions
```

DefaultFormatOptions 返回默认格式化选项：
- 缩进: 4 空格
- case 语句缩进: 启用
- 注释保留: 启用

#### func FormatWithOptions

```go
func FormatWithOptions(script string, opts FormatOptions) (string, error)
```

FormatWithOptions 使用指定选项格式化 shell 命令字符串。

**参数：**

- `script`: shell 命令字符串
- `opts`: 格式化选项

**返回：**

- `string`: 格式化后的字符串
- `error`: 解析错误或系统错误

**示例：**

```go
formatted, err := shx.FormatWithOptions("for i in 1 2 3;do echo $i;done", shx.DefaultFormatOptions())
if err != nil {
    log.Fatal(err)
}
fmt.Println(formatted)
```

#### func FormatScriptWithOptions

```go
func FormatScriptWithOptions(filePath string, opts FormatOptions) (string, error)
```

FormatScriptWithOptions 使用指定选项格式化 shell 脚本文件。

**参数：**

- `filePath`: 脚本文件路径
- `opts`: 格式化选项

**返回：**

- `string`: 格式化后的脚本内容
- `error`: 解析错误或系统错误

**示例：**

```go
formatted, err := shx.FormatScriptWithOptions("deploy.sh", shx.DefaultFormatOptions())
if err != nil {
    log.Fatal(err)
}
fmt.Println(formatted)
```

---

### type SyntaxError

```go
type SyntaxError struct {
    File    string // 文件名（如果来自文件则为文件路径，否则为空字符串）
    Line    int    // 错误行号（从 1 开始）
    Column  int    // 错误列号（从 1 开始）
    Message string // 错误描述
}
```

SyntaxError 表示 shell 语法错误，实现了 error 接口

当语法检查发现错误时，`CheckSyntax` 和 `CheckScriptSyntax` 返回此类型的值。

#### func (*SyntaxError) Error

```go
func (e *SyntaxError) Error() string
```

返回格式化的错误信息，格式为：

- 有文件名：`file:line:col: message`
- 无文件名：`line:col: message`
- 仅有 Message：`message`

**示例输出：**

```
1:6: reached EOF without closing quote "
```
```
deploy.sh:1:6: unexpected EOF while looking for matching `'
```

---

### type UnclosedQuoteError

```go
type UnclosedQuoteError struct {
	QuoteType rune // 未闭合的引号类型 (', ", 或 `)
}
```

UnclosedQuoteError 表示命令字符串中存在未闭合的引号

#### func (*UnclosedQuoteError) Error

```go
func (e *UnclosedQuoteError) Error() string
```

#### func (*UnclosedQuoteError) GetQuoteType

```go
func (e *UnclosedQuoteError) GetQuoteType() rune
```

QuoteType 返回未闭合的引号类型