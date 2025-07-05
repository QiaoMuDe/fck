package flags

import (
	"strconv"
	"sync"

	"gitee.com/MM-Q/qflag/validator"
)

// IntFlag 整数类型标志结构体
// 继承BaseFlag[int]泛型结构体,实现Flag接口
type IntFlag struct {
	BaseFlag[int]
	mu sync.RWMutex
}

// Type 返回标志类型
func (f *IntFlag) Type() FlagType { return FlagTypeInt }

// SetRange 设置整数的有效范围
//
// min: 最小值
// max: 最大值
func (f *IntFlag) SetRange(min, max int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	validator := &validator.IntRangeValidator{Min: min, Max: max}
	f.SetValidator(validator)
}

// Set 实现flag.Value接口,解析并设置整数值
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
func (f *IntFlag) String() string {
	return f.BaseFlag.String()
}
