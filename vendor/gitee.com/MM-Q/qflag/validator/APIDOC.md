# Package validator

validator 提供常用的参数验证器实现。

## Types

### BoolValidator

验证布尔值。

```go
type BoolValidator struct{}
```

**方法**

```go
func (v *BoolValidator) Validate(value any) error
```

实现 Validator 接口，检查值是否为布尔类型。

### DurationValidator

验证时间间隔是否有效。

```go
type DurationValidator struct {
    Min time.Duration // 最小时间间隔，包含在内
    Max time.Duration // 最大时间间隔，包含在内
}
```

**方法**

```go
func (v *DurationValidator) Validate(value any) error
```

实现 Validator 接口，检查时间间隔是否有效且在指定范围内。

### EnumValidator

验证值是否在枚举列表中。

```go
type EnumValidator struct {
    AllowedValues []any // 允许的值列表
}
```

**方法**

```go
func (v *EnumValidator) Validate(value any) error
```

实现 Validator 接口，检查值是否在允许的枚举列表中。

### FloatRangeValidator

验证浮点数是否在指定范围内。

```go
type FloatRangeValidator struct {
    Min float64 // 最小值，包含在内
    Max float64 // 最大值，包含在内
}
```

**方法**

```go
func (v *FloatRangeValidator) Validate(value any) error
```

实现 Validator 接口，检查浮点数是否在 [Min, Max] 范围内。

### IntRangeValidator

验证整数是否在指定范围内。

```go
type IntRangeValidator struct {
    Min int // 最小值，包含在内
    Max int // 最大值，包含在内
}
```

验证逻辑：
- 检查整数是否在 [Min, Max] 闭区间范围内。
- 支持所有整数类型（int、int8、int16、int32、uint 等）的验证。
- 注意：此版本使用 int 类型而非 int64，适用于 32 位整数场景。如需 64 位整数验证，请使用 Int64RangeValidator。

**方法**

```go
func (v *IntRangeValidator) Validate(value any) error
```

IntRangeValidator 验证整数是否在指定的 int 范围内。

注意：
1. 支持多种整数类型转换（int/int8/uint 等），但最终会转换为 int 处理。
2. 从宽类型（如 uint64）转换为 int 可能导致溢出。
3. 如需严格类型检查，请使用自定义验证器。

### IntRangeValidator64

验证整数是否在指定范围内。

```go
type IntRangeValidator64 struct {
    Min int64 // 最小值，包含在内
    Max int64 // 最大值，包含在内
}
```

**方法**

```go
func (v *IntRangeValidator64) Validate(value any) error
```

实现 Validator 接口，检查整数是否在 [Min, Max] 范围内。

### IntValueValidator

验证整数是否为指定值之一。

```go
type IntValueValidator struct {
    AllowedValues []int // 允许的整数值列表
}
```

支持验证整数是否匹配预定义的允许值列表中的任何一个值。适用于需要严格限制输入为特定离散值的场景。

使用示例：
```go
validator := &IntValueValidator{AllowedValues: []int{1, 3, 5}}
flag.SetValidator(validator)
```
这将只允许值为 1、3 或 5 的整数通过验证。

注意：空的允许值列表将导致所有值都验证失败。

**方法**

```go
func (v *IntValueValidator) Validate(value any) error
```

实现 Validator 接口，验证值是否为允许的整数之一。

验证逻辑：检查输入整数是否在允许值列表中。

参数：
- `value`：待验证的整数。

返回值：
- 验证通过返回 nil，否则返回错误信息。

支持的整数类型：int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64。

示例：
```go
validator := &IntValueValidator{AllowedValues: []int{1, 3, 5}}
err := validator.Validate(3) // 返回 nil
err := validator.Validate(2) // 返回错误
```

注意：允许值列表为空时，所有值都将验证失败。

### PathValidator

路径验证器。

```go
type PathValidator struct {
    MustExist   bool // 是否必须存在，默认为 true
    IsDirectory bool // 是否必须是目录，默认为 false
}
```

实现 Validator 接口，用于验证路径是否存在且规范化。

**方法**

```go
func (v *PathValidator) Validate(value any) error
```

验证路径是否符合指定规则。

### SliceLengthValidator

验证切片长度是否在指定范围内。

```go
type SliceLengthValidator struct {
    Min int // 最小长度，包含在内
    Max int // 最大长度，包含在内，0 表示不限制
}
```

**方法**

```go
func (v *SliceLengthValidator) Validate(value any) error
```

实现 Validator 接口，检查切片长度是否在 [Min, Max] 范围内。

### StringLengthValidator

验证字符串长度是否在指定范围内。

```go
type StringLengthValidator struct {
    Min int // 最小长度，包含在内
    Max int // 最大长度，包含在内，0 表示不限制
}
```

**方法**

```go
func (v *StringLengthValidator) Validate(value any) error
```

实现 Validator 接口，检查字符串长度是否在 [Min, Max] 范围内。

### StringRegexValidator

验证字符串是否匹配正则表达式。

```go
type StringRegexValidator struct {
    Pattern string         // 正则表达式模式
    Regex   *regexp.Regexp // 编译后的正则表达式
}
```

**方法**

```go
func (v *StringRegexValidator) Validate(value any) error
```

实现 Validator 接口，检查字符串是否匹配正则表达式。
