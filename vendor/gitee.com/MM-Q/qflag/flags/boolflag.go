package flags

import (
	"strconv"
	"strings"
	"sync"
)

// BoolFlag 布尔类型标志结构体
// 继承BaseFlag[bool]泛型结构体,实现Flag接口
type BoolFlag struct {
	BaseFlag[bool]
	mu sync.Mutex
}

// Type 返回标志类型
func (f *BoolFlag) Type() FlagType { return FlagTypeBool }

// Set 实现flag.Value接口,解析并设置布尔值
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
func (f *BoolFlag) String() string {
	return f.BaseFlag.String()
}

// IsBoolFlag 实现flag.Value接口,返回布尔值
func (f *BoolFlag) IsBoolFlag() bool { return true }
