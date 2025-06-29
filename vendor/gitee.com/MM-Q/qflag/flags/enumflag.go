package flags

import (
	"fmt"
	"strings"
)

// EnumFlag 枚举类型标志结构体
// 继承BaseFlag[string]泛型结构体,增加枚举特有的选项验证
type EnumFlag struct {
	BaseFlag[string]
	optionMap map[string]bool // 枚举值映射
}

// 实现Flag接口
func (f *EnumFlag) Type() FlagType { return FlagTypeEnum }

// IsCheck 检查枚举值是否有效
// 返回值: 为nil, 说明值有效,否则返回错误信息
func (f *EnumFlag) IsCheck(value string) error {
	// 如果枚举map为空,则不需要检查
	if len(f.optionMap) == 0 {
		return nil
	}

	// 转换为小写
	value = strings.ToLower(value)

	// 检查值是否在枚举map中
	if _, valid := f.optionMap[value]; !valid {
		var options []string
		for k := range f.optionMap {
			// 添加枚举值
			options = append(options, k)
		}
		return fmt.Errorf("invalid enum value '%s', options are %v", value, options)
	}
	return nil
}

// Set 实现flag.Value接口, 解析并设置枚举值
func (f *EnumFlag) Set(value string) error {
	// 先验证值是否有效
	if err := f.IsCheck(value); err != nil {
		return err
	}
	// 调用基类方法设置值
	return f.BaseFlag.Set(value)
}

// String 实现flag.Value接口, 返回当前值的字符串表示
func (f *EnumFlag) String() string { return f.Get() }

// Init 初始化枚举类型标志, 无需显式调用, 仅在创建标志对象时自动调用
//
// 参数:
//   - longName: 长标志名称
//   - shortName: 短标志字符
//   - defValue: 默认值
//   - usage: 帮助说明
//   - options: 枚举可选值列表
//
// 返回值:
//   - error: 初始化错误信息
func (f *EnumFlag) Init(longName, shortName string, defValue string, usage string, options []string) error {
	// 初始化枚举值
	if options == nil {
		options = make([]string, 0)
	}

	// 1. 初始化基类字段
	valuePtr := new(string)

	// 默认值小写处理
	*valuePtr = strings.ToLower(defValue)

	// 调用基类方法初始化字段
	if err := f.BaseFlag.Init(longName, shortName, defValue, usage, valuePtr); err != nil {
		return err
	}

	// 2. 初始化枚举optionMap（仅在Init阶段修改，无需额外锁）
	// 注意：无需额外锁，因BaseFlag.Init已保证单例初始化
	f.optionMap = make(map[string]bool)
	for _, opt := range options {
		if opt == "" {
			return fmt.Errorf("enum option cannot be empty")
		}
		f.optionMap[strings.ToLower(opt)] = true
	}

	// 3. 验证默认值有效性
	if len(options) > 0 && !f.optionMap[strings.ToLower(defValue)] {
		return fmt.Errorf("default value '%s' not in enum options %v", defValue, options)
	}

	return nil
}
