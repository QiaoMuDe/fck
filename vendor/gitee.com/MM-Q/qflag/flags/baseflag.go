// Package flags 基础标志结构体定义
// 本文件定义了BaseFlag泛型结构体，为所有标志类型提供通用的字段和方法实现，
// 包括标志的初始化、值获取设置、验证器支持等核心功能。
package flags

import (
	"fmt"
	"sync"

	"gitee.com/MM-Q/qflag/qerr"
)

// BaseFlag 泛型基础标志结构体,封装所有标志的通用字段和方法
type BaseFlag[T any] struct {
	longName     string       // 长标志名称
	shortName    string       // 短标志字符
	initialValue T            // 初始默认值
	usage        string       // 帮助说明
	value        *T           // 标志值指针
	baseMu       sync.RWMutex // 基类读写锁
	validator    Validator    // 验证器接口
	initialized  bool         // 标志是否已初始化
	isSet        bool         // 标志是否已被设置值
	envVar       string       // 存储标志关联的环境变量名称
}

// BindEnv 绑定环境变量到标志
//
// 参数:
//   - envName 环境变量名
//
// 返回:
//   - 标志对象本身,支持链式调用
func (f *BaseFlag[T]) BindEnv(envName string) *BaseFlag[T] {
	f.baseMu.Lock()
	defer f.baseMu.Unlock()
	f.envVar = envName
	return f
}

// GetEnvVar 获取绑定的环境变量名
//
// 返回值:
//   - string: 环境变量名
func (f *BaseFlag[T]) GetEnvVar() string {
	f.baseMu.RLock()
	defer f.baseMu.RUnlock()
	return f.envVar
}

// Init 初始化标志的元数据和值指针, 无需显式调用, 仅在创建标志对象时自动调用
//
// 参数:
//   - longName: 长标志名称
//   - shortName: 短标志字符
//   - usage: 帮助说明
//   - value: 标志值指针
//
// 返回值:
//   - error: 初始化错误信息
func (f *BaseFlag[T]) Init(longName, shortName string, usage string, value *T) error {
	f.baseMu.Lock()
	defer f.baseMu.Unlock()

	// 检查是否已初始化
	if f.initialized {
		return qerr.NewValidationErrorf("flag %s/%s already initialized", f.shortName, f.longName)
	}

	// 检查长短标志是否同时为空
	if longName == "" && shortName == "" {
		return qerr.NewValidationError("longName and shortName cannot both be empty")
	}

	// 验证值指针（避免后续空指针异常）
	if value == nil {
		return qerr.NewValidationError("value pointer cannot be nil")
	}

	f.longName = longName   // 初始化长标志名
	f.shortName = shortName // 初始化短标志名
	f.initialValue = *value // 保存初始默认值
	f.usage = usage         // 初始化标志用途
	f.value = value         // 初始化值指针
	f.initialized = true    // 设置初始化完成标志

	return nil
}

// Name 获取标志的名称
//
// 返回值:
//   - string: 标志名称, 优先返回长名称, 如果长名称为空则返回短名称
func (f *BaseFlag[T]) Name() string {
	if f.longName != "" {
		return f.longName
	}
	return f.shortName
}

// LongName 获取标志的长名称
//
// 返回值:
//   - string: 长标志名称
func (f *BaseFlag[T]) LongName() string { return f.longName }

// ShortName 获取标志的短名称
//
// 返回值:
//   - string: 短标志字符
func (f *BaseFlag[T]) ShortName() string { return f.shortName }

// Usage 获取标志的用法说明
//
// 返回值:
//   - string: 用法说明
func (f *BaseFlag[T]) Usage() string { return f.usage }

// GetDefault 获取标志的初始默认值(泛型类型)
//
// 返回值:
//   - T: 初始默认值
func (f *BaseFlag[T]) GetDefault() T { return f.initialValue }

// GetDefaultAny 获取标志的初始默认值(any类型)
//
// 返回值:
//   - any: 初始默认值
func (f *BaseFlag[T]) GetDefaultAny() any { return f.initialValue }

// IsSet 判断标志是否已被设置值
//
// 返回值:
//   - bool: true表示已设置值, false表示未设置
func (f *BaseFlag[T]) IsSet() bool {
	f.baseMu.RLock()
	defer f.baseMu.RUnlock()
	return f.isSet
}

// Get 获取标志的实际值(泛型类型)
//
// 返回值:
//   - T: 标志值
func (f *BaseFlag[T]) Get() T {
	f.baseMu.RLock()
	defer f.baseMu.RUnlock()

	// 如果标志未设置值, 返回初始默认值
	if !f.isSet {
		return f.initialValue
	}

	// 返回标志值
	return *f.value
}

// Set 设置标志的值(泛型类型)
//
// 参数:
//   - value T: 标志值
//
// 返回:
//   - error: 错误信息
func (f *BaseFlag[T]) Set(value T) error {
	f.baseMu.Lock()
	defer f.baseMu.Unlock()

	// 创建一个副本
	v := value

	// 设置标志值前先进行验证
	if f.validator != nil {
		if err := f.validator.Validate(v); err != nil {
			return qerr.NewValidationErrorf("invalid value for %s: %v", f.longName, err)
		}
	}

	// 设置标志值
	f.value = &v

	// 标志已设置
	f.isSet = true

	return nil
}

// SetValidator 设置标志的验证器(泛型类型)
//
// 参数:
//   - validator Validator: 验证器接口
func (f *BaseFlag[T]) SetValidator(validator Validator) {
	f.baseMu.Lock()
	defer f.baseMu.Unlock()
	f.validator = validator
}

// Reset 将标志重置为初始默认值
func (f *BaseFlag[T]) Reset() {
	f.baseMu.Lock()
	defer f.baseMu.Unlock()
	v := f.initialValue // 获取初始默认值
	f.value = &v        // 重置标志值指针
	f.isSet = false     // 标志未设置
}

// String 返回标志的字符串表示
func (f *BaseFlag[T]) String() string {
	return fmt.Sprint(f.Get())
}

// Type 返回标志类型, 默认实现返回0, 需要子类重写
//
// 注意：
//   - 具体标志类型需要重写此方法返回正确的FlagType
func (f *BaseFlag[T]) Type() FlagType {
	return 0 // 默认实现，需要被子类重写
}
