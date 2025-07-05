package flags

import (
	"fmt"
	"strconv"
	"sync"
)

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
func (f *Uint64Flag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 检查是否为空
	if value == "" {
		return fmt.Errorf("empty value")
	}

	// 将字符串解析为无符号整型
	num, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid uint64 value: %v", err)
	}

	val := uint64(num)
	return f.BaseFlag.Set(val)
}
