package flags

import (
	"fmt"
	"strconv"
	"sync"

	"gitee.com/MM-Q/qflag/qerr"
)

// Uint16Flag 16位无符号整数类型标志结构体
// 继承BaseFlag[uint16]泛型结构体,实现Flag接口
type Uint16Flag struct {
	BaseFlag[uint16]            // 基类
	mu               sync.Mutex // 互斥锁
}

// Type 返回标志类型
//
// 返回值:
//   - FlagType: 标志类型枚举值
func (f *Uint16Flag) Type() FlagType { return FlagTypeUint16 }

// String 实现flag.Value接口, 返回当前值的字符串表示
//
// 返回值:
//   - string: 当前值的字符串表示
func (f *Uint16Flag) String() string {
	return fmt.Sprint(f.Get())
}

// Set 实现flag.Value接口, 解析并设置16位无符号整数值
// 验证值是否在uint16范围内(0-65535)
//
// 参数:
//   - value: 待设置的值(0-65535)
//
// 返回值:
//   - error: 解析或验证失败时返回错误信息
func (f *Uint16Flag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 检查是否为空
	if value == "" {
		return qerr.NewValidationError("empty value")
	}

	// 解析字符串为uint64
	num, err := strconv.ParseUint(value, 10, 16)
	if err != nil {
		return qerr.NewValidationErrorf("invalid uint16 value: %v", err)
	}

	// 转换为uint16
	val := uint16(num)

	// 调用基类方法设置值
	return f.BaseFlag.Set(val)
}

// Uint32Flag 32位无符号整数类型标志结构体
// 继承BaseFlag[uint32]泛型结构体,实现Flag接口
type Uint32Flag struct {
	BaseFlag[uint32]            // 基类
	mu               sync.Mutex // 互斥锁
}

// Type 返回标志类型
func (f *Uint32Flag) Type() FlagType { return FlagTypeUint32 }

// String 实现flag.Value接口, 返回当前值的字符串表示
func (f *Uint32Flag) String() string {
	return fmt.Sprint(f.Get())
}

// Set 实现flag.Value接口, 解析并设置32位无符号整数值
// 验证值是否在uint32范围内(0-4294967295)
//
// 参数:
//   - value: 待设置的值(0-4294967295)
//
// 返回值:
//   - error: 解析或验证失败时返回错误信息
func (f *Uint32Flag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 检查是否为空
	if value == "" {
		return qerr.NewValidationError("empty value")
	}

	// 将字符串解析为无符号整型
	num, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return qerr.NewValidationErrorf("invalid uint32 value: %v", err)
	}

	val := uint32(num)
	return f.BaseFlag.Set(val)
}

// Uint64Flag 64位无符号整数类型标志结构体
// 继承BaseFlag[uint64]泛型结构体,实现Flag接口
type Uint64Flag struct {
	BaseFlag[uint64]            // 基类
	mu               sync.Mutex // 互斥锁
}

// Type 返回标志类型
func (f *Uint64Flag) Type() FlagType { return FlagTypeUint64 }

// String 实现flag.Value接口, 返回当前值的字符串表示
func (f *Uint64Flag) String() string {
	return fmt.Sprint(f.Get())
}

// Set 实现flag.Value接口, 解析并设置64位无符号整数值
// 验证值是否在uint64范围内(0-18446744073709551615)
//
// 参数:
//   - value: 待设置的值(0-18446744073709551615)
//
// 返回值:
//   - error: 解析或验证失败时返回错误信息
func (f *Uint64Flag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 检查是否为空
	if value == "" {
		return qerr.NewValidationError("empty value")
	}

	// 将字符串解析为无符号整型
	num, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return qerr.NewValidationErrorf("invalid uint64 value: %v", err)
	}

	val := uint64(num)
	return f.BaseFlag.Set(val)
}
