# Package validator

Package validator 提供常用的参数验证器实现。

## Types

### BoolValidator

```go
type BoolValidator struct{}
```

BoolValidator 验证布尔值。

#### Validate

```go
func (v *BoolValidator) Validate(value any) error
```

Validate 实现 Validator 接口，检查值是否为布尔类型。

- 参数：
  - value any: 待验证的值。

- 返回值：
  - error: 验证错误，如果验证通过则返回 nil。

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

Validate 实现 Validator 接口，检查时间间隔是否有效且在指定范围内。

- 参数：
  - value any: 待验证的值。

- 返回值：
  - error: 验证错误，如果验证通过则返回 nil。

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

Validate 实现 Validator 接口，检查值是否在允许的枚举列表中。

- 参数：
  - value any: 待验证的值。

- 返回值：
  - error: 验证错误，如果验证通过则返回 nil。

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

Validate 实现 Validator 接口，检查浮点数是否在 [Min, Max] 范围内。

- 参数：
  - value any: 待验证的值。

- 返回值：
  - error: 验证错误，如果验证通过则返回 nil。

### IntRangeValidator

```go
type IntRangeValidator struct {
	Min int // 最小值，包含在内
	Max int // 最大值，包含在内
}
```

IntRangeValidator 验证整数是否在指定范围内。

- 验证逻辑：检查整数是否在 [Min, Max] 闭区间范围内。
- 支持所有整数类型（int、int8、int16、int32、uint 等）的验证。
- 注意：此版本使用 int 类型而非 int64，适用于 32 位整数场景。如需 64 位整数验证，请使用 Int64RangeValidator。

#### Validate

```go
func (v *IntRangeValidator) Validate(value any) error
```

Validate 实现 Validator 接口，检查整数是否在指定的 int 范围内。

- 参数：
  - value any: 待验证的值。

- 返回值：
  - error: 验证错误，如果验证通过则返回 nil。

### IntRangeValidator64

```go
type IntRangeValidator64 struct {
	Min int64 // 最小值，包含在内
	Max int64 // 最大值，包含在内
}
```

IntRangeValidator64 验证整数是否在指定范围内。

#### Validate

```go
func (v *IntRangeValidator64) Validate(value any) error
```

Validate 实现 Validator 接口，检查整数是否在 [Min, Max] 范围内。

- 参数：
  - value any: 待验证的值。

- 返回值：
  - error: 验证错误，如果验证通过则返回 nil。

### IntValueValidator

```go
type IntValueValidator struct {
	AllowedValues []int // 允许的整数值列表
}
```

IntValueValidator 验证整数是否为指定值之一。

- 支持验证整数是否匹配预定义的允许值列表中的任何一个值。
- 适用于需要严格限制输入为特定离散值的场景。
- 使用示例：
  ```go
  validator := &IntValueValidator{AllowedValues: []int{1, 3, 5}}
  flag.SetValidator(validator)
  ```
  这将只允许值为 1、3 或 5 的整数通过验证。
- 注意：空的允许值列表将导致所有值都验证失败。

#### Validate

```go
func (v *IntValueValidator) Validate(value any) error
```

Validate 实现 Validator 接口，验证值是否为允许的整数之一。

- 参数：
  - value any: 待验证的值。

- 返回值：
  - error: 验证错误，如果验证通过则返回 nil。

### PathValidator

```go
type PathValidator struct {
	MustExist   bool // 是否必须存在，默认为 true
	IsDirectory bool // 是否必须是目录，默认为 false
}
```

PathValidator 路径验证器，实现 Validator 接口，用于验证路径是否存在且规范化。

#### Validate

```go
func (v *PathValidator) Validate(value any) error
```

Validate 验证路径是否符合指定规则。

- 参数：
  - value any: 待验证的值。

- 返回值：
  - error: 验证错误，如果验证通过则返回 nil。

### SliceLengthValidator

```go
type SliceLengthValidator struct {
	Min int // 最小长度，包含在内
	Max int // 最大长度，包含在内，0 表示不限制
}
```

SliceLengthValidator 验证切片长度是否在指定范围内。

#### Validate

```go
func (v *SliceLengthValidator) Validate(value any) error
```

Validate 实现 Validator 接口，检查切片长度是否在 [Min, Max] 范围内。

- 参数：
  - value any: 待验证的值。

- 返回值：
  - error: 验证错误，如果验证通过则返回 nil。

### StringLengthValidator

```go
type StringLengthValidator struct {
	Min int // 最小长度，包含在内
	Max int // 最大长度，包含在内，0 表示不限制
}
```

StringLengthValidator 验证字符串长度是否在指定范围内。

#### Validate

```go
func (v *StringLengthValidator) Validate(value any) error
```

Validate 实现 Validator 接口，检查字符串长度是否在 [Min, Max] 范围内。

- 参数：
  - value any: 待验证的值。

- 返回值：
  - error: 验证错误，如果验证通过则返回 nil。

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

Validate 实现 Validator 接口，检查字符串是否匹配正则表达式。

- 参数：
  - value any: 待验证的值。

- 返回值：
  - error: 验证错误，如果验证通过则返回 nil。