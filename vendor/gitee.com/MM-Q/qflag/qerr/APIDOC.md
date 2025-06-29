# qerr API 文档

## 常量

### 命令行解析相关错误常量

```go
const (
	ErrFlagParseFailed       = "Parameter parsing error"  // 全局实例标志解析错误
	ErrSubCommandParseFailed = "Subcommand parsing error" // 子命令标志解析错误
	ErrPanicRecovered        = "panic recovered"          // 恐慌捕获错误
	ErrValidationFailed      = "Validation failed"        // 参数验证失败错误
)
```

## 函数

### JoinErrors

将错误切片合并为单个错误，并去除重复错误。

```go
func JoinErrors(errors []error) error
```

- **参数**：
  - `errors`：错误切片。
- **返回值**：
  - 合并后的错误信息。

### NewValidationError

创建一个新的验证错误。

```go
func NewValidationError(message string) error
```

- **参数**：
  - `message`：错误信息。
- **返回值**：
  - 验证错误。

### NewValidationErrorf

创建一个格式化的验证错误。

```go
func NewValidationErrorf(format string, v ...interface{}) error
```

- **参数**：
  - `format`：格式化字符串。
  - `v`：格式化参数。
- **返回值**：
  - 格式化的验证错误。