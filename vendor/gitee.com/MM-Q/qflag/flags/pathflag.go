package flags

import (
	"fmt"
	"path/filepath"

	"gitee.com/MM-Q/qflag/validator"
)

// PathFlag 路径类型标志结构体
// 继承BaseFlag[string]泛型结构体,实现Flag接口
type PathFlag struct {
	BaseFlag[string]
}

// Type 返回标志类型
func (f *PathFlag) Type() FlagType { return FlagTypePath }

// String 实现flag.Value接口,返回当前值的字符串表示
func (f *PathFlag) String() string { return f.Get() }

// Set 实现flag.Value接口,解析并验证路径
func (f *PathFlag) Set(value string) error {
	if value == "" {
		return fmt.Errorf("path cannot be empty")
	}

	// 规范化路径为绝对路径
	absPath, err := filepath.Abs(value)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}

	// 调用基类方法设置值(会触发验证器验证)
	return f.BaseFlag.Set(absPath)
}

// Init 初始化路径标志
func (f *PathFlag) Init(longName, shortName string, defValue string, usage string) error {
	// 初始化路径标志值指针
	valuePtr := new(string)

	// 规范化默认路径为绝对路径
	absDefValue, err := filepath.Abs(defValue)
	if err != nil {
		return fmt.Errorf("failed to normalize default path: %v", err)
	}

	// 设置默认值
	*valuePtr = absDefValue

	// 调用基类方法初始化
	if err := f.BaseFlag.Init(longName, shortName, defValue, usage, valuePtr); err != nil {
		return err
	}

	// 设置路径验证器
	f.SetValidator(&validator.PathValidator{})
	return nil
}
