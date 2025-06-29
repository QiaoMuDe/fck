package flags

import "gitee.com/MM-Q/qflag/validator"

// Int64Flag 64位整数类型标志结构体
// 继承BaseFlag[int64]泛型结构体,实现Flag接口
type Int64Flag struct {
	BaseFlag[int64]
}

// Type 返回标志类型
func (f *Int64Flag) Type() FlagType { return FlagTypeInt64 }

// SetRange 设置64位整数的有效范围
//
// min: 最小值
// max: 最大值
func (f *Int64Flag) SetRange(min, max int64) {
	validator := &validator.IntRangeValidator64{Min: min, Max: max}
	f.SetValidator(validator)
}
