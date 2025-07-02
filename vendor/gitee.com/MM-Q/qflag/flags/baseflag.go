package flags

import (
	"fmt"
	"sync"
)

// BaseFlag 泛型基础标志结构体,封装所有标志的通用字段和方法
type BaseFlag[T any] struct {
	longName    string       // 长标志名称
	shortName   string       // 短标志字符
	defValue    T            // 默认值
	usage       string       // 帮助说明
	value       *T           // 标志值指针
	baseMu      sync.RWMutex // 基类读写锁
	validator   Validator    // 验证器接口
	initialized bool         // 标志是否已初始化
	isSet       bool         // 标志是否已被设置值
}

// Init 初始化标志的元数据和值指针, 无需显式调用, 仅在创建标志对象时自动调用
//
// 参数:
//   - longName: 长标志名称
//   - shortName: 短标志字符
//   - defValue: 默认值
//   - usage: 帮助说明
//   - value: 标志值指针
//
// 返回值:
//   - error: 初始化错误信息
func (f *BaseFlag[T]) Init(longName, shortName string, defValue T, usage string, value *T) error {
	f.baseMu.Lock()
	defer f.baseMu.Unlock()

	// 检查是否已初始化
	if f.initialized {
		return fmt.Errorf("flag %s/%s already initialized", f.shortName, f.longName)
	}

	// 检查长短标志是否同时为空
	if longName == "" && shortName == "" {
		return fmt.Errorf("longName and shortName cannot both be empty")
	}

	// 验证值指针（避免后续空指针异常）
	if value == nil {
		return fmt.Errorf("value pointer cannot be nil")
	}

	f.longName = longName   // 初始化长标志名
	f.shortName = shortName // 初始化短标志名
	f.defValue = defValue   // 初始化默认值
	f.usage = usage         // 初始化标志用途
	f.value = value         // 初始化值指针
	f.initialized = true    // 设置初始化完成标志

	return nil
}

// LongName 获取标志的长名称
func (f *BaseFlag[T]) LongName() string { return f.longName }

// ShortName 获取标志的短名称
func (f *BaseFlag[T]) ShortName() string { return f.shortName }

// Usage 获取标志的用法说明
func (f *BaseFlag[T]) Usage() string { return f.usage }

// GetDefault 获取标志的默认值
func (f *BaseFlag[T]) GetDefault() T { return f.defValue }

// GetDefaultAny 获取标志的默认值(any类型)
func (f *BaseFlag[T]) GetDefaultAny() any { return f.defValue }

// IsSet 判断标志是否已被设置值
//
// 返回值: true表示已设置值, false表示未设置
func (f *BaseFlag[T]) IsSet() bool {
	f.baseMu.RLock()
	defer f.baseMu.RUnlock()
	return f.isSet
}

// Get 获取标志的实际值
// 优先级：已设置的值 > 默认值
// 线程安全：使用互斥锁保证并发访问安全
func (f *BaseFlag[T]) Get() T {
	f.baseMu.RLock()
	defer f.baseMu.RUnlock()

	// 如果标志值不为空,则返回标志值
	if f.value != nil {
		return *f.value
	}

	// 否则返回默认值
	return f.defValue
}

// GetPointer 返回标志值的指针
//
// 注意:
//  1. 获取指针过程受锁保护, 但直接修改指针指向的值仍会绕过验证机制
//  2. 多线程环境下修改时需额外同步措施, 建议优先使用Set()方法
func (f *BaseFlag[T]) GetPointer() *T {
	f.baseMu.RLock()
	defer f.baseMu.RUnlock()
	return f.value
}

// Set 设置标志的值
//
// 参数: value 标志值
//
// 返回: 错误信息
func (f *BaseFlag[T]) Set(value T) error {
	f.baseMu.Lock()
	defer f.baseMu.Unlock()

	// 创建一个副本
	v := value

	// 设置标志值前先进行验证
	if f.validator != nil {
		if err := f.validator.Validate(v); err != nil {
			return fmt.Errorf("invalid value for %s: %w", f.longName, err)
		}
	}

	// 设置标志值
	f.value = &v

	// 标志已设置
	f.isSet = true

	return nil
}

// SetValidator 设置标志的验证器
//
// 参数: validator 验证器接口
func (f *BaseFlag[T]) SetValidator(validator Validator) {
	f.baseMu.Lock()
	defer f.baseMu.Unlock()
	f.validator = validator
}

// String 返回标志的字符串表示
func (f *BaseFlag[T]) String() string {
	return fmt.Sprint(f.Get())
}

// Reset 将标志值重置为默认值
// 线程安全：使用互斥锁保证并发安全
func (f *BaseFlag[T]) Reset() {
	f.baseMu.Lock()
	defer f.baseMu.Unlock()
	f.value = nil   // 重置为未设置状态,下次Get()会返回默认值
	f.isSet = false // 重置设置状态
}
