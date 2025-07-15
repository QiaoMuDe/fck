package flags

import (
	"fmt"
	"strings"
)

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
