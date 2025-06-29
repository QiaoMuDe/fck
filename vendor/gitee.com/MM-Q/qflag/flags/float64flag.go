package flags

// Float64Flag 浮点型标志结构体
// 继承BaseFlag[float64]泛型结构体,实现Flag接口
type Float64Flag struct {
	BaseFlag[float64]
}

// Type 返回标志类型
func (f *Float64Flag) Type() FlagType { return FlagTypeFloat64 }
