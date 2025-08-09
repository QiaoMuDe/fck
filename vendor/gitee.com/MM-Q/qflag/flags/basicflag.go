// Package flags 基本数据类型标志实现
// 本文件实现了整数、浮点数、布尔、字符串等基本数据类型的标志结构体，
// 提供了相应的解析、验证和类型转换功能。
package flags

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"gitee.com/MM-Q/qflag/qerr"
	"gitee.com/MM-Q/qflag/validator"
)

// =============================================================================
// 整数类型标志
// =============================================================================

// IntFlag 整数类型标志结构体
// 继承BaseFlag[int]泛型结构体,实现Flag接口
type IntFlag struct {
	BaseFlag[int]
	mu sync.RWMutex
}

// Type 返回标志类型
//
// 返回值:
//   - FlagType: 标志类型枚举值
func (f *IntFlag) Type() FlagType { return FlagTypeInt }

// SetRange 设置整数的有效范围
//
// 参数:
//   - min: 最小值
//   - max: 最大值
func (f *IntFlag) SetRange(min, max int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	validator := &validator.IntRangeValidator{Min: min, Max: max}
	f.SetValidator(validator)
}

// Set 实现flag.Value接口,解析并验证整数值
//
// 参数:
//   - value: 待解析的整数值
//
// 返回值:
//   - error: 解析错误或验证错误
func (f *IntFlag) Set(value string) error {
	f.mu.RLock()
	defer f.mu.RUnlock()
	intVal, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	return f.BaseFlag.Set(intVal)
}

// String 实现flag.Value接口,返回当前整数值的字符串表示
//
// 返回值:
//   - string: 当前整数值的字符串表示
func (f *IntFlag) String() string {
	return f.BaseFlag.String()
}

// Int64Flag 64位整数类型标志结构体
// 继承BaseFlag[int64]泛型结构体,实现Flag接口
type Int64Flag struct {
	BaseFlag[int64]
	mu sync.Mutex // 互斥锁
}

// Type 返回标志类型
//
// 返回值:
//   - FlagType: 标志类型枚举值
func (f *Int64Flag) Type() FlagType { return FlagTypeInt64 }

// SetRange 设置64位整数的有效范围
//
// 参数:
//   - min: 最小值
//   - max: 最大值
func (f *Int64Flag) SetRange(min, max int64) {
	f.mu.Lock()
	defer f.mu.Unlock()
	validator := &validator.IntRangeValidator64{Min: min, Max: max}
	f.SetValidator(validator)
}

// Set 实现flag.Value接口,解析并设置64位整数值
//
// 参数:
//   - value: 待解析的64位整数值
//
// 返回值:
//   - error: 解析错误或验证错误
func (f *Int64Flag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	int64Val, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return qerr.NewValidationErrorf("failed to parse int64 value: %v", err)
	}
	return f.BaseFlag.Set(int64Val)
}

// =============================================================================
// 64位浮点数类型标志
// =============================================================================

// Type 返回标志类型
//
// 返回值:
//   - FlagType: 标志类型枚举值
func (f *Float64Flag) Type() FlagType { return FlagTypeFloat64 }

// Set 实现flag.Value接口,解析并设置浮点值
//
// 参数:
//   - value: 待解析的浮点值
//
// 返回值:
//   - error: 解析错误或验证错误
func (f *Float64Flag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	floatVal, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return qerr.NewValidationErrorf("failed to parse float64 value: %v", err)
	}
	return f.BaseFlag.Set(floatVal)
}

// =============================================================================
// 布尔类型标志
// =============================================================================

// BoolFlag 布尔类型标志结构体
// 继承BaseFlag[bool]泛型结构体,实现Flag接口
type BoolFlag struct {
	BaseFlag[bool]
	mu sync.Mutex
}

// Type 返回标志类型
//
// 返回值:
//   - FlagType: 标志类型枚举值
func (f *BoolFlag) Type() FlagType { return FlagTypeBool }

// Set 实现flag.Value接口,解析并设置布尔值
//
// 支持以下布尔值格式（大小写不敏感）:
//   - 真值: "true", "1", "t", "T", "TRUE", "True"
//   - 假值: "false", "0", "f", "F", "FALSE", "False"
//
// 参数:
//   - value: 待设置的布尔值字符串
//
// 返回值:
//   - error: 解析或验证失败时返回错误信息
//
// 示例:
//   - flag.Set("true")   // ✅ 成功，值为 true
//   - flag.Set("1")      // ✅ 成功，值为 true
//   - flag.Set("FALSE")  // ✅ 成功，值为 false
//   - flag.Set("0")      // ✅ 成功，值为 false
//   - flag.Set("yes")    // ❌ 失败，返回解析错误
func (f *BoolFlag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 解析字符串为布尔值
	boolVal, err := strconv.ParseBool(strings.ToLower(value))
	if err != nil {
		return err
	}
	return f.BaseFlag.Set(boolVal)
}

// String 实现flag.Value接口,返回布尔值字符串
//
// 返回值:
//   - string: 布尔值字符串
func (f *BoolFlag) String() string {
	return f.BaseFlag.String()
}

// IsBoolFlag 实现flag.boolFlag接口,返回true
func (f *BoolFlag) IsBoolFlag() bool { return true }

// =============================================================================
// 字符串类型标志
// =============================================================================

// StringFlag 字符串类型标志结构体
// 继承BaseFlag[string]泛型结构体,实现Flag接口
type StringFlag struct {
	BaseFlag[string]
}

// Type 返回标志类型
//
// 返回值:
//   - FlagType: 标志类型枚举值
func (f *StringFlag) Type() FlagType { return FlagTypeString }

// String 返回带引号的字符串值
//
// 返回值:
//   - string: 带引号的字符串值
func (f *StringFlag) String() string {
	return fmt.Sprintf("%q", f.Get())
}

// Len 获取字符串标志的长度
//
// 返回值：
//   - 字符串的字符数(按UTF-8编码计算)
func (f *StringFlag) Len() int {
	return len(f.Get())
}

// ToUpper 将字符串标志值转换为大写
func (f *StringFlag) ToUpper() string {
	return strings.ToUpper(f.Get())
}

// ToLower 将字符串标志值转换为小写
func (f *StringFlag) ToLower() string {
	return strings.ToLower(f.Get())
}

// Contains 检查字符串是否包含指定子串
//
// 参数:
//   - substr 子串
//
// 返回值:
//   - bool: 如果包含子串则返回true,否则返回false
func (f *StringFlag) Contains(substr string) bool {
	return strings.Contains(f.Get(), substr)
}

// Set 实现flag.Value接口的Set方法
// 将字符串值解析并设置到标志中
//
// 参数:
//   - value: 待设置的字符串值
//
// 返回值:
//   - error: 设置失败时返回错误信息
func (f *StringFlag) Set(value string) error {
	return f.BaseFlag.Set(value)
}
