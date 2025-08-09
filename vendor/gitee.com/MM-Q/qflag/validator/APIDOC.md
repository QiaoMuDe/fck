# Package validator

**导入路径:** `gitee.com/MM-Q/qflag/validator`

Package validator 参数验证器实现。本文件提供了常用的参数验证器实现，包括字符串长度验证、正则表达式验证、数值范围验证、枚举值验证、路径验证等功能，为各种标志类型提供值的有效性验证支持。

## 类型

### BoolValidator

```go
type BoolValidator struct{}
```

BoolValidator 验证布尔值。

#### Validate

```go
func (v *BoolValidator) Validate(value any) error
```

Validate 实现Validator接口，检查值是否为布尔类型。

**参数:**
- `value any`: 待验证的值

**返回值:**
- `error`: 验证错误，如果验证通过则返回nil

**示例:**
```go
validator := &BoolValidator{}
err := validator.Validate(true)  // 返回nil
err := validator.Validate("yes") // 返回错误
```

---

### DurationValidator

```go
type DurationValidator struct {
    Min time.Duration // 最小时间间隔，包含在内
    Max time.Duration // 最大时间间隔，包含在内
}
```

DurationValidator 验证时间间隔是否有效。

#### Validate

```go
func (v *DurationValidator) Validate(value any) error
```

Validate 实现Validator接口，检查时间间隔是否有效且在指定范围内。

**参数:**
- `value any`: 待验证的值

**返回值:**
- `error`: 验证错误，如果验证通过则返回nil

**示例:**
```go
validator := &DurationValidator{
    Min: time.Second,
    Max: time.Hour,
}
err := validator.Validate(time.Minute * 30) // 返回nil
err := validator.Validate(time.Millisecond) // 返回错误（小于最小值）
```

---

### EnumValidator

```go
type EnumValidator struct {
    AllowedValues []any // 允许的值列表
}
```

EnumValidator 验证值是否在枚举列表中。

#### Validate

```go
func (v *EnumValidator) Validate(value any) error
```

Validate 实现Validator接口，检查值是否在允许的枚举列表中。

**参数:**
- `value any`: 待验证的值

**返回值:**
- `error`: 验证错误，如果验证通过则返回nil

**示例:**
```go
validator := &EnumValidator{
    AllowedValues: []any{"debug", "info", "warn", "error"},
}
err := validator.Validate("info")  // 返回nil
err := validator.Validate("trace") // 返回错误
```

---

### FloatRangeValidator

```go
type FloatRangeValidator struct {
    Min float64 // 最小值，包含在内
    Max float64 // 最大值，包含在内
}
```

FloatRangeValidator 验证浮点数是否在指定范围内。

#### Validate

```go
func (v *FloatRangeValidator) Validate(value any) error
```

Validate 实现Validator接口，检查浮点数是否在[Min, Max]范围内。

**参数:**
- `value any`: 待验证的值

**返回值:**
- `error`: 验证错误，如果验证通过则返回nil

**示例:**
```go
validator := &FloatRangeValidator{
    Min: 0.0,
    Max: 100.0,
}
err := validator.Validate(50.5)  // 返回nil
err := validator.Validate(150.0) // 返回错误
```

---

### IntRangeValidator

```go
type IntRangeValidator struct {
    Min int // 最小值，包含在内
    Max int // 最大值，包含在内
}
```

IntRangeValidator 验证整数是否在指定范围内。

**验证逻辑:**
- 检查整数是否在[Min, Max]闭区间范围内
- 支持所有整数类型（int、int8、int16、int32、uint等）的验证
- 实现了Validator接口

**注意:** 此版本使用int类型而非int64，适用于32位整数场景。如需64位整数验证，请使用IntRangeValidator64。

#### Validate

```go
func (v *IntRangeValidator) Validate(value any) error
```

Validate 实现Validator接口，检查整数是否在指定的int范围内。

**参数:**
- `value any`: 待验证的值

**返回值:**
- `error`: 验证错误，如果验证通过则返回nil

**示例:**
```go
validator := &IntRangeValidator{
    Min: 1,
    Max: 100,
}
err := validator.Validate(50)  // 返回nil
err := validator.Validate(150) // 返回错误
```

---

### IntRangeValidator64

```go
type IntRangeValidator64 struct {
    Min int64 // 最小值，包含在内
    Max int64 // 最大值，包含在内
}
```

IntRangeValidator64 验证整数是否在指定范围内（64位版本）。

#### Validate

```go
func (v *IntRangeValidator64) Validate(value any) error
```

Validate 实现Validator接口，检查整数是否在[Min, Max]范围内。

**参数:**
- `value any`: 待验证的值

**返回值:**
- `error`: 验证错误，如果验证通过则返回nil

**示例:**
```go
validator := &IntRangeValidator64{
    Min: 1,
    Max: 9223372036854775807, // int64最大值
}
err := validator.Validate(int64(1000000000000)) // 返回nil
```

---

### IntValueValidator

```go
type IntValueValidator struct {
    AllowedValues []int // 允许的整数值列表
}
```

IntValueValidator 验证整数是否为指定值之一。

**功能特点:**
- 支持验证整数是否匹配预定义的允许值列表中的任何一个值
- 适用于需要严格限制输入为特定离散值的场景
- 支持的整数类型：int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64
- 空的允许值列表将导致所有值都验证失败

#### Validate

```go
func (v *IntValueValidator) Validate(value any) error
```

Validate 实现Validator接口，验证值是否为允许的整数之一。

**参数:**
- `value any`: 待验证的值

**返回值:**
- `error`: 验证错误，如果验证通过则返回nil

**示例:**
```go
validator := &IntValueValidator{
    AllowedValues: []int{1, 3, 5},
}
err := validator.Validate(3) // 返回nil
err := validator.Validate(2) // 返回错误
```

---

### PathValidator

```go
type PathValidator struct {
    MustExist   bool // 是否必须存在，默认为true
    IsDirectory bool // 是否必须是目录，默认为false
}
```

PathValidator 路径验证器，实现Validator接口，用于验证路径是否存在且规范化。

#### Validate

```go
func (v *PathValidator) Validate(value any) error
```

Validate 验证路径是否符合指定规则。

**参数:**
- `value any`: 待验证的值

**返回值:**
- `error`: 验证错误，如果验证通过则返回nil

**示例:**
```go
// 验证文件必须存在
validator := &PathValidator{
    MustExist:   true,
    IsDirectory: false,
}
err := validator.Validate("/path/to/file.txt")

// 验证目录必须存在
dirValidator := &PathValidator{
    MustExist:   true,
    IsDirectory: true,
}
err := dirValidator.Validate("/path/to/directory")
```

---

### SliceLengthValidator

```go
type SliceLengthValidator struct {
    Min int // 最小长度，包含在内
    Max int // 最大长度，包含在内，0表示不限制
}
```

SliceLengthValidator 验证切片长度是否在指定范围内。

#### Validate

```go
func (v *SliceLengthValidator) Validate(value any) error
```

Validate 实现Validator接口，检查切片长度是否在[Min, Max]范围内。

**参数:**
- `value any`: 待验证的值

**返回值:**
- `error`: 验证错误，如果验证通过则返回nil

**示例:**
```go
validator := &SliceLengthValidator{
    Min: 1,
    Max: 10,
}
err := validator.Validate([]string{"a", "b", "c"}) // 返回nil
err := validator.Validate([]string{})              // 返回错误（长度为0，小于最小值1）
```

---

### StringLengthValidator

```go
type StringLengthValidator struct {
    Min int // 最小长度，包含在内
    Max int // 最大长度，包含在内，0表示不限制
}
```

StringLengthValidator 验证字符串长度是否在指定范围内。

#### Validate

```go
func (v *StringLengthValidator) Validate(value any) error
```

Validate 实现Validator接口，检查字符串长度是否在[Min, Max]范围内。

**参数:**
- `value any`: 待验证的值

**返回值:**
- `error`: 验证错误，如果验证通过则返回nil

**示例:**
```go
validator := &StringLengthValidator{
    Min: 3,
    Max: 20,
}
err := validator.Validate("hello")     // 返回nil
err := validator.Validate("hi")        // 返回错误（长度为2，小于最小值3）
err := validator.Validate("very long string that exceeds limit") // 返回错误
```

---

### StringRegexValidator

```go
type StringRegexValidator struct {
    Pattern string         // 正则表达式模式
    Regex   *regexp.Regexp // 编译后的正则表达式
}
```

StringRegexValidator 验证字符串是否匹配正则表达式。

#### Validate

```go
func (v *StringRegexValidator) Validate(value any) error
```

Validate 实现Validator接口，检查字符串是否匹配正则表达式。

**参数:**
- `value any`: 待验证的值

**返回值:**
- `error`: 验证错误，如果验证通过则返回nil

**示例:**
```go
validator := &StringRegexValidator{
    Pattern: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
}
// 需要先编译正则表达式
validator.Regex = regexp.MustCompile(validator.Pattern)

err := validator.Validate("user@example.com") // 返回nil
err := validator.Validate("invalid-email")    // 返回错误