// Package validator 参数验证器实现
// 本文件提供了常用的参数验证器实现，包括字符串长度验证、正则表达式验证、
// 数值范围验证、枚举值验证、路径验证等功能，为各种标志类型提供值的有效性验证支持。
package validator

import (
	"errors"
	"fmt"
	"math"
	"os"
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
//
// 参数值:
//   - value any: 待验证的值
//
// 返回值:
//   - error: 验证错误, 如果验证通过则返回nil
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
//
// 参数值:
//   - value any: 待验证的值
//
// 返回值:
//   - error: 验证错误, 如果验证通过则返回nil
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
//
// 验证逻辑：检查整数是否在[Min, Max]闭区间范围内
// 支持所有整数类型（int、int8、int16、int32、uint等）的验证
// 实现了Validator接口
// IntRangeValidator 验证整数是否在指定范围内
//
// 注意：此版本使用int类型而非int64，适用于32位整数场景
// 如需64位整数验证，请使用Int64RangeValidator
//
// 实现了Validator接口
type IntRangeValidator struct {
	Min int // 最小值, 包含在内
	Max int // 最大值, 包含在内
}

// Validate 实现Validator接口, 检查整数是否在指定的 int 范围内
//
// 参数值:
//   - value any: 待验证的值
//
// 返回值:
//   - error: 验证错误, 如果验证通过则返回nil
func (v *IntRangeValidator) Validate(value any) error {
	var num int

	// 处理不同整数类型的转换
	switch val := value.(type) {
	case int:
		num = val
	case int8:
		num = int(val)
	case int16:
		num = int(val)
	case int32:
		num = int(val)
	case int64:
		if val > math.MaxInt32 || val < math.MinInt32 {
			return fmt.Errorf("int64 value %d exceeds int range [%d, %d]", val, math.MinInt32, math.MaxInt32)
		}
		num = int(val)
	case uint:
		num = int(val)
	case uint8:
		num = int(val)
	case uint16:
		num = int(val)
	case uint32:
		num = int(val)
	case uint64:
		if val > uint64(math.MaxInt32) {
			return fmt.Errorf("uint64 value %d exceeds int max value %d", val, math.MaxInt32)
		}
		num = int(val)
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

// IntRangeValidator64 验证整数是否在指定范围内
type IntRangeValidator64 struct {
	Min int64 // 最小值, 包含在内
	Max int64 // 最大值, 包含在内
}

// Validate 实现Validator接口, 检查整数是否在[Min, Max]范围内
//
// 参数值:
//   - value any: 待验证的值
//
// 返回值:
//   - error: 验证错误, 如果验证通过则返回nil
func (v *IntRangeValidator64) Validate(value any) error {
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
		// 增加uint64转换的显式溢出检查
		if val > math.MaxInt64 {
			return fmt.Errorf("uint64 value %d exceeds int64 max value %d", val, math.MaxInt64)
		}
		num = int64(val)
	default:
		return errors.New("value is not an int64-compatible integer type")
	}

	if num < v.Min {
		return fmt.Errorf("value must be at least %d", v.Min)
	}

	if num > v.Max {
		return fmt.Errorf("value must be at most %d", v.Max)
	}

	return nil
}

// IntValueValidator 验证整数是否为指定值之一
//
// 支持验证整数是否匹配预定义的允许值列表中的任何一个值
// 适用于需要严格限制输入为特定离散值的场景
//
// 使用示例:
// validator := &IntValueValidator{AllowedValues: []int{1, 3, 5}}
// flag.SetValidator(validator)
// 这将只允许值为1、3或5的整数通过验证
//
// 注意: 空的允许值列表将导致所有值都验证失败
//
// 实现了Validator接口
// IntValueValidator 验证整数是否为指定值之一
//
// 验证逻辑: 检查输入整数是否在允许值列表中
//
// 参数:
//
//	value: 待验证的整数
//
// 返回值:
//
//	验证通过返回nil，否则返回错误信息
//
// 支持的整数类型: int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64
//
// 示例:
//
//	validator := &IntValueValidator{AllowedValues: []int{1, 3, 5}}
//	err := validator.Validate(3) // 返回nil
//	err := validator.Validate(2) // 返回错误
//
// 注意: 允许值列表为空时，所有值都将验证失败
type IntValueValidator struct {
	AllowedValues []int // 允许的整数值列表
}

// Validate 实现Validator接口, 验证值是否为允许的整数之一
//
// 参数值:
//   - value any: 待验证的值
//
// 返回值:
//   - error: 验证错误, 如果验证通过则返回nil
func (v *IntValueValidator) Validate(value any) error {
	// 检查允许值列表是否为空
	if len(v.AllowedValues) == 0 {
		return errors.New("no allowed values specified")
	}

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

	// 检查值是否在允许值列表中
	found := false
	for _, allowed := range v.AllowedValues {
		if int64(allowed) == num {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("value %d is not in allowed values %v", num, v.AllowedValues)
	}

	return nil
}

// FloatRangeValidator 验证浮点数是否在指定范围内
type FloatRangeValidator struct {
	Min float64 // 最小值, 包含在内
	Max float64 // 最大值, 包含在内
}

// Validate 实现Validator接口, 检查浮点数是否在[Min, Max]范围内
//
// 参数值:
//   - value any: 待验证的值
//
// 返回值:
//   - error: 验证错误, 如果验证通过则返回nil
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
//
// 参数值:
//   - value any: 待验证的值
//
// 返回值:
//   - error: 验证错误, 如果验证通过则返回nil
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
//
// 参数值:
//   - value any: 待验证的值
//
// 返回值:
//   - error: 验证错误, 如果验证通过则返回nil
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
//
// 参数值:
//   - value any: 待验证的值
//
// 返回值:
//   - error: 验证错误, 如果验证通过则返回nil
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
//
// 参数值:
//   - value any: 待验证的值
//
// 返回值:
//   - error: 验证错误, 如果验证通过则返回nil
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

// PathValidator 路径验证器
// 实现Validator接口,用于验证路径是否存在且规范化
type PathValidator struct {
	MustExist   bool // 是否必须存在，默认为true
	IsDirectory bool // 是否必须是目录，默认为false
}

// Validate 验证路径是否符合指定规则
//
// 参数值:
//   - value any: 待验证的值
//
// 返回值:
//   - error: 验证错误, 如果验证通过则返回nil
func (v *PathValidator) Validate(value any) error {
	path, ok := value.(string)
	if !ok {
		return fmt.Errorf("path must be a string")
	}

	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	// 检查路径是否存在（如果需要）
	var fi os.FileInfo
	var err error
	if v.MustExist || v.IsDirectory {
		fi, err = os.Stat(path)
		if os.IsNotExist(err) {
			return fmt.Errorf("path does not exist: %s", path)
		} else if err != nil {
			return fmt.Errorf("failed to check path: %v", err)
		}
	}

	// 检查是否为目录（独立于存在性检查）
	if v.IsDirectory {
		if fi == nil {
			return fmt.Errorf("cannot verify directory type for non-existent path")
		} else if !fi.IsDir() {
			return fmt.Errorf("path is not a directory: %s", path)
		}
	}

	return nil
}
