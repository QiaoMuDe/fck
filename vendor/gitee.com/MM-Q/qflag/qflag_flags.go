// flags 定义了所有标志类型的通用接口和基础标志结构体
package qflag

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// Validator 验证器接口, 所有自定义验证器需实现此接口
type Validator interface {
	// Validate 验证参数值是否合法
	// value: 待验证的参数值
	// 返回值: 验证通过返回nil, 否则返回错误信息
	Validate(value any) error
}

// Flag 所有标志类型的通用接口,定义了标志的元数据访问方法
type Flag interface {
	LongName() string   // 获取标志的长名称
	ShortName() string  // 获取标志的短名称
	Usage() string      // 获取标志的用法
	Type() FlagType     // 获取标志类型
	GetDefaultAny() any // 获取标志的默认值(any类型)
	String() string     // 获取标志的字符串表示
	IsSet() bool        // 判断标志是否已设置值
	Reset()             // 重置标志值为默认值
}

// TypedFlag 所有标志类型的通用接口,定义了标志的元数据访问方法和默认值访问方法
type TypedFlag[T any] interface {
	Flag                    // 继承标志接口
	GetDefault() T          // 获取标志的具体类型默认值
	Get() T                 // 获取标志的具体类型值
	GetPointer() *T         // 获取标志值的指针
	Set(T) error            // 设置标志的具体类型值
	SetValidator(Validator) // 设置标志的验证器
}

// BaseFlag 泛型基础标志结构体,封装所有标志的通用字段和方法
type BaseFlag[T any] struct {
	cmd       *Cmd       // 所属的命令实例
	longName  string     // 长标志名称
	shortName string     // 短标志字符
	defValue  T          // 默认值
	usage     string     // 帮助说明
	value     *T         // 标志值指针
	mu        sync.Mutex // 并发访问锁
	validator Validator  // 验证器接口
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
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.value != nil
}

// Get 获取标志的实际值
// 优先级：已设置的值 > 默认值
// 线程安全：使用互斥锁保证并发访问安全
func (f *BaseFlag[T]) Get() T {
	f.mu.Lock()
	defer f.mu.Unlock()

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
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.value
}

// Set 设置标志的值
//
// 参数: value 标志值
//
// 返回: 错误信息
func (f *BaseFlag[T]) Set(value T) error {
	f.mu.Lock()
	defer f.mu.Unlock()

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

	return nil
}

// SetValidator 设置标志的验证器
//
// 参数: validator 验证器接口
func (f *BaseFlag[T]) SetValidator(validator Validator) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.validator = validator
}

// String 返回标志的字符串表示
func (f *BaseFlag[T]) String() string {
	return fmt.Sprint(f.Get())
}

// Reset 将标志值重置为默认值
// 线程安全：使用互斥锁保证并发安全
func (f *BaseFlag[T]) Reset() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.value = nil // 重置为未设置状态,下次Get()会返回默认值
}

// IntFlag 整数类型标志结构体
// 继承BaseFlag[int]泛型结构体,实现Flag接口
type IntFlag struct {
	BaseFlag[int]
}

// Type 返回标志类型
func (f *IntFlag) Type() FlagType { return FlagTypeInt }

// StringFlag 字符串类型标志结构体
// 继承BaseFlag[string]泛型结构体,实现Flag接口
type StringFlag struct {
	BaseFlag[string]
}

// Type 返回标志类型
func (f *StringFlag) Type() FlagType { return FlagTypeString }

// String 返回带引号的字符串值
func (f *StringFlag) String() string {
	return fmt.Sprintf("%q", f.Get())
}

// Len 获取字符串标志的长度
//
// 返回值：字符串的字符数(按UTF-8编码计算)
func (f *StringFlag) Len() int {
	return len(f.Get())
}

// BoolFlag 布尔类型标志结构体
// 继承BaseFlag[bool]泛型结构体,实现Flag接口
type BoolFlag struct {
	BaseFlag[bool]
}

// Type 返回标志类型
func (f *BoolFlag) Type() FlagType { return FlagTypeBool }

// FloatFlag 浮点型标志结构体
// 继承BaseFlag[float64]泛型结构体,实现Flag接口
type FloatFlag struct {
	BaseFlag[float64]
}

// Type 返回标志类型
func (f *FloatFlag) Type() FlagType { return FlagTypeFloat }

// DurationFlag 时间间隔类型标志结构体
// 继承BaseFlag[time.Duration]泛型结构体,实现Flag接口
type DurationFlag struct {
	BaseFlag[time.Duration]
}

// Type 返回标志类型
func (f *DurationFlag) Type() FlagType { return FlagTypeDuration }

// Set 实现flag.Value接口, 解析并设置时间间隔值
func (f *DurationFlag) Set(value string) error {
	// 检查空值
	if value == "" {
		return fmt.Errorf("duration cannot be empty")
	}

	// 将单位转换为小写, 确保解析的准确性
	lowercaseValue := strings.ToLower(value)

	// 解析时间间隔字符串
	duration, err := time.ParseDuration(lowercaseValue)
	if err != nil {
		return fmt.Errorf("invalid duration format: %v (valid units: ns/us/ms/s/m/h)", err)
	}

	// 检查负值（可选）
	if duration < 0 {
		return fmt.Errorf("negative duration not allowed")
	}

	// 调用基类方法设置值
	return f.BaseFlag.Set(duration)
}

// String 实现flag.Value接口, 返回当前值的字符串表示
func (f *DurationFlag) String() string {
	return f.Get().String()
}

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

// SliceFlag 切片类型标志结构体
// 继承BaseFlag[[]string]泛型结构体,实现Flag接口
type SliceFlag struct {
	BaseFlag[[]string]            // 基类
	delimiters         []string   // 分隔符
	mu                 sync.Mutex // 锁
	SkipEmpty          bool       // 是否跳过空元素
}

// Type 返回标志类型
func (f *SliceFlag) Type() FlagType { return FlagTypeSlice }

// String 实现flag.Value接口, 返回当前值的字符串表示
func (f *SliceFlag) String() string {
	return strings.Join(f.Get(), ",")
}

// Set 实现flag.Value接口, 解析并设置切片值
//
// 参数: value 待解析的切片值
//
// 注意: 如果切片中包含分隔符,则根据分隔符进行分割, 否则将整个值作为单个元素
// 例如: "a,b,c" -> ["a", "b", "c"]
func (f *SliceFlag) Set(value string) error {
	// 检查空值
	if value == "" {
		return fmt.Errorf("slice cannot be empty")
	}

	// 获取当前切片值
	current := f.Get()
	var elements []string // 存储分割后的元素

	// 加锁保护分隔符切片访问
	f.mu.Lock()
	defer f.mu.Unlock()

	// 检查是否包含分隔符切片中的任何分隔符
	found := false
	for _, delimiter := range f.delimiters {
		if strings.Contains(value, delimiter) {
			// 根据分隔符分割字符串
			elements = strings.Split(value, delimiter)
			// 去除每个元素的首尾空白字符
			for i, e := range elements {
				elements[i] = strings.TrimSpace(e)
			}
			found = true
			break // 找到第一个匹配的分隔符后停止
		}
	}

	// 如果没有找到分隔符,将整个值作为单个元素
	if !found {
		elements = []string{strings.TrimSpace(value)}
	}

	// 过滤空元素（如果启用）
	if f.SkipEmpty {
		filtered := make([]string, 0, len(elements))
		for _, e := range elements {
			if e != "" {
				filtered = append(filtered, e)
			}
		}
		elements = filtered
	}

	// 预分配切片容量以减少内存分配
	newValues := make([]string, 0, len(current)+len(elements))

	// 将当前值和新增的值添加到新的切片中
	newValues = append(newValues, current...)
	newValues = append(newValues, elements...)

	// 调用基类方法设置值
	return f.BaseFlag.Set(newValues)
}

// SetDelimiters 设置切片解析的分隔符列表
//
// 参数: delimiters 分隔符列表
//
// 线程安全的分隔符更新
func (f *SliceFlag) SetDelimiters(delimiters []string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.delimiters = delimiters
}

// GetDelimiters 获取当前分隔符列表
func (f *SliceFlag) GetDelimiters() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	// 返回拷贝避免外部修改内部切片
	res := make([]string, len(f.delimiters))
	copy(res, f.delimiters)
	return res
}

// SetSkipEmpty 设置是否跳过空元素
//
// 参数: skip - 为true时跳过空元素, 为false时保留空元素
func (f *SliceFlag) SetSkipEmpty(skip bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.SkipEmpty = skip
}

// Len 获取切片长度
//
// 返回: 获取切片长度
func (f *SliceFlag) Len() int {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 如果切片为nil, 则返回0
	if f.value == nil {
		return 0
	}

	// 返回切片长度
	return len(*f.value)
}

// Contains 检查切片是否包含指定元素
// 当切片未设置值时,将使用默认值进行检查
func (f *SliceFlag) Contains(element string) bool {
	// 通过Get()获取当前值(已处理nil情况和线程安全)
	current := f.Get()

	// 直接遍历当前值(已确保非nil)
	for _, item := range current {
		if item == element {
			return true
		}
	}
	return false
}

// Clear 清空切片所有元素
//
// 注意：该方法会改变切片的指针
func (f *SliceFlag) Clear() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.value = &[]string{}
}
