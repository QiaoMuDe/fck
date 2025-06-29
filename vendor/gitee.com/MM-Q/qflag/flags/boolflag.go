package flags

// BoolFlag 布尔类型标志结构体
// 继承BaseFlag[bool]泛型结构体,实现Flag接口
type BoolFlag struct {
	BaseFlag[bool]
}

// Type 返回标志类型
func (f *BoolFlag) Type() FlagType { return FlagTypeBool }
