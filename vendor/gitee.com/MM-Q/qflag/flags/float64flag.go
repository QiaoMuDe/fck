package flags

import (
	"strconv"
	"sync"

	"gitee.com/MM-Q/qflag/qerr"
)

// Float64Flag 浮点型标志结构体
// 继承BaseFlag[float64]泛型结构体,实现Flag接口
type Float64Flag struct {
	BaseFlag[float64]
	mu sync.Mutex
}

// Type 返回标志类型
func (f *Float64Flag) Type() FlagType { return FlagTypeFloat64 }

// Set 实现flag.Value接口,解析并设置浮点值
func (f *Float64Flag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	floatVal, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return qerr.NewValidationErrorf("failed to parse float64 value: %v", err)
	}
	return f.BaseFlag.Set(floatVal)
}
