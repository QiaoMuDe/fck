package flags

import (
	"fmt"
	"strconv"
	"sync"
)

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
func (f *Uint32Flag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 检查是否为空
	if value == "" {
		return fmt.Errorf("empty value")
	}

	// 将字符串解析为无符号整型
	num, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid uint32 value: %v", err)
	}

	val := uint32(num)
	return f.BaseFlag.Set(val)
}
