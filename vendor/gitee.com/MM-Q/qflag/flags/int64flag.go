package flags

import (
	"strconv"
	"sync"

	"gitee.com/MM-Q/qflag/validator"
)

// Int64Flag 64位整数类型标志结构体
// 继承BaseFlag[int64]泛型结构体,实现Flag接口
type Int64Flag struct {
	BaseFlag[int64]
	mu sync.Mutex // 互斥锁
}

// Type 返回标志类型
func (f *Int64Flag) Type() FlagType { return FlagTypeInt64 }

// SetRange 设置64位整数的有效范围
//
// min: 最小值
// max: 最大值
func (f *Int64Flag) SetRange(min, max int64) {
	f.mu.Lock()
	defer f.mu.Unlock()
	validator := &validator.IntRangeValidator64{Min: min, Max: max}
	f.SetValidator(validator)
}

// Set 实现flag.Value接口,解析并设置64位整数值
func (f *Int64Flag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	int64Val, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	return f.BaseFlag.Set(int64Val)
}
