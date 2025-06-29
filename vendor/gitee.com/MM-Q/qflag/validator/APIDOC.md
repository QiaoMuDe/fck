# Validator API 文档

## BoolValidator

### 类型定义

```go
type BoolValidator struct{}
```

### 方法

#### Validate

验证值是否为布尔类型。

```go
func (v *BoolValidator) Validate(value any) error
```

- **参数**：
  - `value`：待验证的值。
- **返回值**：
  - 如果值为布尔类型，返回 `nil`；否则返回错误信息。

## DurationValidator

### 类型定义

```go
type DurationValidator struct {
	Min time.Duration // 最小时间间隔, 包含在内
	Max time.Duration // 最大时间间隔, 包含在内
}
```

### 方法

#### Validate

验证时间间隔是否有效且在指定范围内。

```go
func (v *DurationValidator) Validate(value any) error
```

- **参数**：
  - `value`：待验证的时间间隔，支持字符串类型（如 `"5m"`）和 `time.Duration` 类型。
- **返回值**：
  - 如果时间间隔有效且在指定范围内，返回 `nil`；否则返回错误信息。

## EnumValidator

### 类型定义

```go
type EnumValidator struct {
	AllowedValues []any // 允许的值列表
}
```

### 方法

#### Validate

验证值是否在枚举列表中。

```go
func (v *EnumValidator) Validate(value any) error
```

- **参数**：
  - `value`：待验证的值。
- **返回值**：
  - 如果值在枚举列表中，返回 `nil`；否则返回错误信息。

## FloatRangeValidator

### 类型定义

```go
type FloatRangeValidator struct {
	Min float64 // 最小值, 包含在内
	Max float64 // 最大值, 包含在内
}
```

### 方法

#### Validate

验证浮点数是否在指定范围内。

```go
func (v *FloatRangeValidator) Validate(value any) error
```

- **参数**：
  - `value`：待验证的浮点数。
- **返回值**：
  - 如果浮点数在指定范围内，返回 `nil`；否则返回错误信息。

## IntRangeValidator

### 类型定义

```go
type IntRangeValidator struct {
	Min int // 最小值, 包含在内
	Max int // 最大值, 包含在内
}
```

### 方法

#### Validate

验证整数是否在指定范围内。

```go
func (v *IntRangeValidator) Validate(value any) error
```

- **参数**：
  - `value`：待验证的整数。
- **返回值**：
  - 如果整数在指定范围内，返回 `nil`；否则返回错误信息。

## IntRangeValidator64

### 类型定义

```go
type IntRangeValidator64 struct {
	Min int64 // 最小值, 包含在内
	Max int64 // 最大值, 包含在内
}
```

### 方法

#### Validate

验证整数是否在指定范围内。

```go
func (v *IntRangeValidator64) Validate(value any) error
```

- **参数**：
  - `value`：待验证的整数。
- **返回值**：
  - 如果整数在指定范围内，返回 `nil`；否则返回错误信息。

## IntValueValidator

### 类型定义

```go
type IntValueValidator struct {
	AllowedValues []int // 允许的整数值列表
}
```

### 方法

#### Validate

验证整数是否为指定值之一。

```go
func (v *IntValueValidator) Validate(value any) error
```

- **参数**：
  - `value`：待验证的整数。
- **返回值**：
  - 如果整数在允许值列表中，返回 `nil`；否则返回错误信息。

## PathValidator

### 类型定义

```go
type PathValidator struct{}
```

### 方法

#### Validate

验证路径是否存在且规范化。

```go
func (v *PathValidator) Validate(value any) error
```

- **参数**：
  - `value`：待验证的路径。
- **返回值**：
  - 如果路径存在且规范化，返回 `nil`；否则返回错误信息。

## SliceLengthValidator

### 类型定义

```go
type SliceLengthValidator struct {
	Min int // 最小长度, 包含在内
	Max int // 最大长度, 包含在内, 0表示不限制
}
```

### 方法

#### Validate

验证切片长度是否在指定范围内。

```go
func (v *SliceLengthValidator) Validate(value any) error
```

- **参数**：
  - `value`：待验证的切片。
- **返回值**：
  - 如果切片长度在指定范围内，返回 `nil`；否则返回错误信息。

## StringLengthValidator

### 类型定义

```go
type StringLengthValidator struct {
	Min int // 最小长度, 包含在内
	Max int // 最大长度, 包含在内, 0表示不限制
}
```

### 方法

#### Validate

验证字符串长度是否在指定范围内。

```go
func (v *StringLengthValidator) Validate(value any) error
```

- **参数**：
  - `value`：待验证的字符串。
- **返回值**：
  - 如果字符串长度在指定范围内，返回 `nil`；否则返回错误信息。

## StringRegexValidator

### 类型定义

```go
type StringRegexValidator struct {
	Pattern string         // 正则表达式模式
	Regex   *regexp.Regexp // 编译后的正则表达式
}
```

### 方法

#### Validate

验证字符串是否匹配正则表达式。

```go
func (v *StringRegexValidator) Validate(value any) error
```

- **参数**：
  - `value`：待验证的字符串。
- **返回值**：
  - 如果字符串匹配正则表达式，返回 `nil`；否则返回错误信息。