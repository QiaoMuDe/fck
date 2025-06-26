// qflag_validator.go
// 提供常用的参数验证器实现
package qflag

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"time"
)

// StringLengthValidator 验证字符串长度是否在指定范围内
type StringLengthValidator struct {
	Min int // 最小长度, 包含在内
	Max int // 最大长度, 包含在内, 0表示不限制
}

// Validate 实现Validator接口, 检查字符串长度是否在[Min, Max]范围内
func (v *StringLengthValidator) Validate(value any) error {
	s, ok := value.(string)
	if !ok {
		return errors.New("value is not a string")
	}

	if len(s) < v.Min {
		return fmt.Errorf("string length must be at least %d", v.Min)
	}

	if v.Max > 0 && len(s) > v.Max {
		return fmt.Errorf("string length must be at most %d", v.Max)
	}

	return nil
}

// StringRegexValidator 验证字符串是否匹配正则表达式
type StringRegexValidator struct {
	Pattern string         // 正则表达式模式
	Regex   *regexp.Regexp // 编译后的正则表达式
}

// Validate 实现Validator接口, 检查字符串是否匹配正则表达式
func (v *StringRegexValidator) Validate(value any) error {
	s, ok := value.(string)
	if !ok {
		return errors.New("value is not a string")
	}

	if v.Regex == nil {
		if v.Pattern == "" {
			return errors.New("regex pattern is empty")
		}
		var err error
		v.Regex, err = regexp.Compile(v.Pattern)
		if err != nil {
			return fmt.Errorf("invalid regex pattern: %v", err)
		}
	}

	if !v.Regex.MatchString(s) {
		return fmt.Errorf("string does not match pattern: %s", v.Pattern)
	}

	return nil
}

// IntRangeValidator 验证整数是否在指定范围内
type IntRangeValidator struct {
	Min int64 // 最小值, 包含在内
	Max int64 // 最大值, 包含在内
}

// Validate 实现Validator接口, 检查整数是否在[Min, Max]范围内
func (v *IntRangeValidator) Validate(value any) error {
	var num int64

	// 处理不同整数类型的转换
	switch val := value.(type) {
	case int:
		num = int64(val)
	case int8:
		num = int64(val)
	case int16:
		num = int64(val)
	case int32:
		num = int64(val)
	case int64:
		num = val
	case uint:
		num = int64(val)
	case uint8:
		num = int64(val)
	case uint16:
		num = int64(val)
	case uint32:
		num = int64(val)
	case uint64:
		num = int64(val)
	default:
		return errors.New("value is not an integer type")
	}

	if num < v.Min {
		return fmt.Errorf("value must be at least %d", v.Min)
	}

	if num > v.Max {
		return fmt.Errorf("value must be at most %d", v.Max)
	}

	return nil
}

// FloatRangeValidator 验证浮点数是否在指定范围内
type FloatRangeValidator struct {
	Min float64 // 最小值, 包含在内
	Max float64 // 最大值, 包含在内
}

// Validate 实现Validator接口, 检查浮点数是否在[Min, Max]范围内
func (v *FloatRangeValidator) Validate(value any) error {
	var num float64

	// 处理不同浮点数类型的转换
	switch val := value.(type) {
	case float32:
		num = float64(val)
	case float64:
		num = val
	default:
		return errors.New("value is not a float type")
	}

	if num < v.Min {
		return fmt.Errorf("value must be at least %f", v.Min)
	}

	if num > v.Max {
		return fmt.Errorf("value must be at most %f", v.Max)
	}

	return nil
}

// BoolValidator 验证布尔值
type BoolValidator struct{}

// Validate 实现Validator接口, 检查值是否为布尔类型
func (v *BoolValidator) Validate(value any) error {
	_, ok := value.(bool)
	if !ok {
		return errors.New("value is not a boolean")
	}
	return nil
}

// DurationValidator 验证时间间隔是否有效
type DurationValidator struct {
	Min time.Duration // 最小时间间隔, 包含在内
	Max time.Duration // 最大时间间隔, 包含在内
}

// Validate 实现Validator接口, 检查时间间隔是否有效且在指定范围内
func (v *DurationValidator) Validate(value any) error {
	// 支持字符串类型的时间间隔（如"5m"）和time.Duration类型
	var d time.Duration
	var err error

	switch val := value.(type) {
	case string:
		d, err = time.ParseDuration(val)
		if err != nil {
			return fmt.Errorf("invalid duration string: %v", err)
		}
	case time.Duration:
		d = val
	default:
		return errors.New("value is not a duration string or time.Duration")
	}

	if d < v.Min {
		return fmt.Errorf("duration must be at least %v", v.Min)
	}

	if v.Max > 0 && d > v.Max {
		return fmt.Errorf("duration must be at most %v", v.Max)
	}

	return nil
}

// SliceLengthValidator 验证切片长度是否在指定范围内
type SliceLengthValidator struct {
	Min int // 最小长度, 包含在内
	Max int // 最大长度, 包含在内, 0表示不限制
}

// Validate 实现Validator接口, 检查切片长度是否在[Min, Max]范围内
func (v *SliceLengthValidator) Validate(value any) error {
	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		return errors.New("value is not a slice or array")
	}

	length := val.Len()
	if length < v.Min {
		return fmt.Errorf("slice length must be at least %d", v.Min)
	}

	if v.Max > 0 && length > v.Max {
		return fmt.Errorf("slice length must be at most %d", v.Max)
	}

	return nil
}

// EnumValidator 验证值是否在枚举列表中
type EnumValidator struct {
	AllowedValues []any // 允许的值列表
}

// Validate 实现Validator接口, 检查值是否在允许的枚举列表中
func (v *EnumValidator) Validate(value any) error {
	if len(v.AllowedValues) == 0 {
		return errors.New("no allowed values specified")
	}

	for _, allowed := range v.AllowedValues {
		if reflect.DeepEqual(value, allowed) {
			return nil
		}
	}

	return fmt.Errorf("value %v is not in allowed values list", value)
}
