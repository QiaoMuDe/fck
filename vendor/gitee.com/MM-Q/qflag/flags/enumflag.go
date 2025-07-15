package flags

import (
	"strings"
	"sync"

	"gitee.com/MM-Q/qflag/qerr"
)

// EnumFlag 枚举类型标志结构体
// 继承BaseFlag[string]泛型结构体,增加枚举特有的选项验证
type EnumFlag struct {
	BaseFlag[string]
	optionMap       map[string]bool // 枚举值映射
	originalOptions []string        // 原始选项(未处理)
	caseSensitive   bool            // 是否区分大小写
	mu              sync.RWMutex    // 读写锁
}

// Type 返回标志类型
//
// 返回值:
//   - FlagType: 标志类型枚举值
func (f *EnumFlag) Type() FlagType { return FlagTypeEnum }

// SetCaseSensitive 设置枚举值是否区分大小写
//
// 参数:
//   - sensitive - true表示区分大小写，false表示不区分（默认）
//
// 返回值:
//   - *EnumFlag - 返回自身以支持链式调用
func (f *EnumFlag) SetCaseSensitive(sensitive bool) *EnumFlag {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 仅当设置值改变时才重建映射
	if f.caseSensitive == sensitive {
		return f
	}

	// 更新大小写敏感设置
	f.caseSensitive = sensitive

	// 根据新的大小写设置重建选项映射
	f.optionMap = make(map[string]bool, len(f.originalOptions))

	// 添加枚举值到选项映射
	for _, opt := range f.originalOptions {
		if opt == "" {
			continue
		}

		key := opt
		// 如果不区分大小写，则将枚举值转换为小写
		if !f.caseSensitive {
			key = strings.ToLower(opt)
		}
		f.optionMap[key] = true
	}

	return f
}

// IsCheck 检查枚举值是否有效
//
// 参数:
//   - value: 待检查的枚举值
//
// 返回值:
//   - error: 为nil, 说明值有效,否则返回错误信息
func (f *EnumFlag) IsCheck(value string) error {
	f.mu.RLock()
	defer f.mu.RUnlock()

	// 如果枚举map为空,则不需要检查
	if len(f.optionMap) == 0 {
		return nil
	}

	// 根据大小写敏感设置处理值
	checkValue := value
	if !f.caseSensitive {
		checkValue = strings.ToLower(checkValue)
	}

	// 检查值是否在枚举map中
	if _, valid := f.optionMap[checkValue]; !valid {
		return qerr.NewValidationErrorf("invalid enum value '%s', options are %v", value, f.originalOptions)
	}
	return nil
}

// Set 实现flag.Value接口, 解析并设置枚举值
//
// 参数:
//   - value: 待设置的值
//
// 返回值:
//   - error: 解析或验证失败时返回错误信息
func (f *EnumFlag) Set(value string) error {
	// 先验证值是否有效
	if err := f.IsCheck(value); err != nil {
		return qerr.NewValidationErrorf("failed to set enum value: %v", err)
	}
	// 调用基类方法设置值
	return f.BaseFlag.Set(value)
}

// String 实现flag.Value接口, 返回当前值的字符串表示
//
// 返回值:
//   - string: 当前值的字符串表示
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
	f.mu.Lock()
	defer f.mu.Unlock()

	// 初始化枚举值
	if options == nil {
		options = make([]string, 0)
	}

	// 1. 初始化基类字段
	valuePtr := new(string)

	// 根据大小写敏感设置处理默认值
	if !f.caseSensitive {
		*valuePtr = strings.ToLower(defValue)
	} else {
		*valuePtr = defValue
	}

	// 调用基类方法初始化字段
	if err := f.BaseFlag.Init(longName, shortName, usage, valuePtr); err != nil {
		return err
	}

	// 2. 初始化枚举optionMap（仅在Init阶段修改，无需额外锁）
	// 注意：无需额外锁，因BaseFlag.Init已保证单例初始化
	f.optionMap = make(map[string]bool)                 // 枚举值映射
	f.originalOptions = make([]string, 0, len(options)) // 原始选项切片
	for _, opt := range options {
		if opt == "" {
			return qerr.NewValidationError("enum option cannot be empty")
		}
		f.originalOptions = append(f.originalOptions, opt) // 保存原始选项

		// 如果不区分大小写，则将枚举值转换为小写
		key := opt
		if !f.caseSensitive {
			key = strings.ToLower(opt)
		}
		f.optionMap[key] = true
	}

	// 3. 验证默认值有效性
	checkValue := defValue // 根据大小写敏感设置处理默认值
	if !f.caseSensitive {
		checkValue = strings.ToLower(checkValue)
	}
	if len(options) > 0 && !f.optionMap[checkValue] {
		return qerr.NewValidationErrorf("default value '%s' not in enum options %v", defValue, options)
	}

	return nil
}

// GetOptions 返回枚举的所有可选值
//
// 返回值:
//   - []string: 枚举的所有可选值
func (f *EnumFlag) GetOptions() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	// 返回副本以避免外部修改
	options := make([]string, len(f.originalOptions))
	copy(options, f.originalOptions)
	return options
}
