# Package qerr

Package qerr 提供命令行解析相关错误处理。

## VARIABLES

```go
var (
    ErrFlagParseFailed       = errors.New("parameter parsing error")             // 全局实例标志解析错误
    ErrSubCommandParseFailed = errors.New("subcommand parsing error")            // 子命令标志解析错误
    ErrPanicRecovered        = errors.New("panic recovered")                     // 恐慌捕获错误
    ErrValidationFailed      = errors.New("validation failed")                   // 参数验证失败错误
    ErrEnvLoadFailed         = errors.New("environment variable loading failed") // 环境变量加载失败错误
    ErrAddSubCommandFailed   = errors.New("add subcommand failed")               // 添加子命令失败错误
)
```

命令行解析相关错误变量。

## FUNCTIONS

### JoinErrors

```go
func JoinErrors(errors []error) error
```

JoinErrors 将错误切片合并为单个错误，并去除重复错误。

- 参数值：
  - errors []error: 错误切片。

- 返回值：
  - error: 合并后的错误。

### NewValidationError

```go
func NewValidationError(message string) error
```

NewValidationError 创建一个新的验证错误。

- 参数值：
  - message string: 错误消息。

- 返回值：
  - error: 验证错误。

### NewValidationErrorf

```go
func NewValidationErrorf(format string, v ...interface{}) error
```

NewValidationErrorf 创建一个格式化的验证错误。

- 参数值：
  - format string: 格式化字符串。
  - v ...interface{}: 格式化参数。

- 返回值：
  - error: 验证错误。