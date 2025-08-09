# Package qerr

**Import Path:** `gitee.com/MM-Q/qflag/qerr`

Package qerr 错误处理和自定义错误类型定义。本文件定义了qflag包使用的各种自定义错误类型，提供统一的错误处理机制，包括验证错误、解析错误等不同类型的错误定义和处理方法。

## Variables

### 命令行解析相关错误变量

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

预定义的命令行解析相关错误变量，用于标识不同类型的解析和处理错误。

## Functions

### JoinErrors

```go
func JoinErrors(errors []error) error
```

JoinErrors 将错误切片合并为单个错误，并去除重复错误

**参数:**
- `errors []error`: 错误切片

**返回值:**
- `error`: 合并后的错误

**功能说明:**
- 将多个错误合并为一个错误
- 自动去除重复的错误
- 返回合并后的错误信息

**示例:**
```go
err1 := errors.New("first error")
err2 := errors.New("second error")
err3 := errors.New("first error") // 重复错误

combinedErr := qerr.JoinErrors([]error{err1, err2, err3})
// 结果只包含 "first error" 和 "second error"
```

### NewValidationError

```go
func NewValidationError(message string) error
```

NewValidationError 创建一个新的验证错误

**参数:**
- `message string`: 错误消息

**返回值:**
- `error`: 验证错误

**功能说明:**
- 创建标准的验证错误
- 用于参数验证失败的场景
- 返回包含指定消息的错误

**示例:**
```go
err := qerr.NewValidationError("invalid email format")
if err != nil {
    fmt.Printf("Validation failed: %s\n", err.Error())
}
```

### NewValidationErrorf

```go
func NewValidationErrorf(format string, v ...interface{}) error
```

NewValidationErrorf 创建一个格式化的验证错误

**参数:**
- `format string`: 格式化字符串
- `v ...interface{}`: 格式化参数

**返回值:**
- `error`: 验证错误

**功能说明:**
- 创建格式化的验证错误
- 支持printf风格的格式化
- 便于创建包含动态内容的错误消息

**示例:**
```go
fieldName := "email"
value := "invalid-email"
err := qerr.NewValidationErrorf("field '%s' has invalid value: %s", fieldName, value)
// 结果: "field 'email' has invalid value: invalid-email"
```

## 使用场景

### 错误合并

当需要收集多个验证错误并统一处理时：

```go
var errors []error

// 收集各种验证错误
if name == "" {
    errors = append(errors, qerr.NewValidationError("name is required"))
}
if email == "" {
    errors = append(errors, qerr.NewValidationError("email is required"))
}
if age < 0 {
    errors = append(errors, qerr.NewValidationErrorf("age must be positive, got %d", age))
}

// 合并所有错误
if len(errors) > 0 {
    return qerr.JoinErrors(errors)
}
```

### 验证错误处理

在命令行参数验证中：

```go
func validateFlags(config *Config) error {
    if config.Port < 1 || config.Port > 65535 {
        return qerr.NewValidationErrorf("port must be between 1 and 65535, got %d", config.Port)
    }
    
    if config.Host == "" {
        return qerr.NewValidationError("host cannot be empty")
    }
    
    return nil
}
```

### 错误类型检查

检查特定类型的错误：

```go
func handleError(err error) {
    if err == nil {
        return
    }
    
    // 检查是否为预定义的错误类型
    switch err {
    case qerr.ErrFlagParseFailed:
        fmt.Println("Failed to parse command line flags")
    case qerr.ErrValidationFailed:
        fmt.Println("Parameter validation failed")
    case qerr.ErrSubCommandParseFailed:
        fmt.Println("Failed to parse subcommand")
    default:
        fmt.Printf("Unknown error: %s\n", err.Error())
    }
}
```

## 最佳实践

1. **使用合适的错误类型** - 根据错误场景选择合适的预定义错误或创建新的验证错误
2. **错误消息清晰** - 提供清晰、具体的错误消息，帮助用户理解问题
3. **错误合并** - 在需要收集多个错误时使用JoinErrors进行合并
4. **格式化错误** - 使用NewValidationErrorf创建包含动态信息的错误消息
5. **错误检查** - 在适当的地方检查特定的错误类型并进行相应处理

## 注意事项

- 所有函数都返回标准的error接口，与Go标准错误处理兼容
- JoinErrors会自动去重，避免重复的错误消息
- 验证错误函数适用于参数验证场景
- 预定义的错误变量可用于错误类型比较和处理
