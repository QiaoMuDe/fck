package flags

import (
	"fmt"
	"strconv"
)

// Uint16Flag 16位无符号整数类型标志结构体
// 继承BaseFlag[uint16]泛型结构体,实现Flag接口
type Uint16Flag struct {
	BaseFlag[uint16]
}

// Type 返回标志类型
func (f *Uint16Flag) Type() FlagType { return FlagTypeUint16 }

// String 实现flag.Value接口, 返回当前值的字符串表示
func (f *Uint16Flag) String() string {
	return fmt.Sprint(f.Get())
}

// Set 实现flag.Value接口, 解析并设置16位无符号整数值
// 验证值是否在uint16范围内(0-65535)
//
// 参数:
//
//	value: 待设置的值(0-65535)
//
// 返回值:
//
//	error: 错误信息
func (f *Uint16Flag) Set(value string) error {
	// 解析字符串为uint64
	num, err := strconv.ParseUint(value, 10, 16)
	if err != nil {
		return fmt.Errorf("invalid uint16 value: %v", err)
	}
	// 转换为uint16
	val := uint16(num)
	// 调用基类方法设置值
	return f.BaseFlag.Set(val)
}
