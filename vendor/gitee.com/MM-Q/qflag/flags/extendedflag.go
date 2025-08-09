// Package flags 扩展数据类型标志实现
// 本文件实现了枚举、时间间隔、切片、时间、映射等扩展数据类型的标志结构体，
// 提供了复杂数据类型的解析、验证和格式化功能。
package flags

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"gitee.com/MM-Q/qflag/qerr"
)

// =============================================================================
// 时间类型标志
// =============================================================================

// 支持的时间格式列表
var supportedTimeFormats = []string{
	time.RFC3339,          // 2006-01-02T15:04:05Z07:00
	"2006-01-02",          // 2006-01-02
	"2006-01-02 15:04",    // 2006-01-02 15:04
	"2006-01-02 15:04:05", // 2006-01-02 15:04:05
	time.ANSIC,            // Mon Jan _2 15:04:05 MST 2006
	time.UnixDate,         // Mon Jan _2 15:04:05 MST 2006
	time.RubyDate,         // Mon Jan 02 15:04:05 -0700 2006
	time.RFC822,           // 01 Jan 06 15:04 MST
	time.RFC822Z,          // 02 Jan 06 15:04 -0700
	time.RFC850,           // Monday, 02-Jan-06 15:04:05 MST
	time.RFC1123,          // Mon, 02 Jan 2006 15:04:05 MST
	time.RFC1123Z,         // Mon, 02 Jan 2006 15:04:05 -0700
	time.Kitchen,          // 3:04PM
	time.Stamp,            // Mon Jan _2 15:04:05
	time.StampMilli,       // Mon Jan _2 15:04:05.000
	time.StampMicro,       // Mon Jan _2 15:04:05.000000
	time.StampNano,        // Mon Jan _2 15:04:05.000000000
}

// TimeFlag 时间类型标志结构体
// 继承BaseFlag[time.Time]泛型结构体,实现Flag接口
type TimeFlag struct {
	BaseFlag[time.Time]
	outputFormat string     // 自定义输出格式
	mu           sync.Mutex // 保护outputFormat和value的并发访问
}

// Type 返回标志类型
//
// 返回值:
//   - FlagType: 标志类型枚举值
func (f *TimeFlag) Type() FlagType { return FlagTypeTime }

// Set 实现flag.Value接口, 解析并设置时间值
//
// 参数:
//   - value: 待解析的时间字符串
//
// 返回值:
//   - error: 解析或验证失败时返回错误信息
func (f *TimeFlag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	var t time.Time
	var err error

	// 尝试解析时间字符串
	for _, format := range supportedTimeFormats {
		t, err = time.Parse(format, value)
		if err == nil {
			break
		}
	}
	if err != nil {
		return qerr.NewValidationErrorf("invalid time format: %v (supported formats include %v)", err, supportedTimeFormats)
	}

	// 调用基类设置值
	return f.BaseFlag.Set(t)
}

// String 实现flag.Value接口, 返回当前时间的字符串表示
//
// 返回值:
//   - string: 格式化后的时间字符串
//
// 注意: 加锁保证outputFormat和value的并发安全访问
func (f *TimeFlag) String() string {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 获取当前值和输出格式
	value := f.Get()
	format := f.outputFormat

	// 如果设置了输出格式, 则使用该格式
	if format != "" {
		return value.Format(format)
	}
	return value.Format(time.RFC3339) // 默认格式
}

// SetOutputFormat 设置时间输出格式
//
// 参数:
//   - format: 时间格式化字符串
//
// 注意: 此方法线程安全
func (f *TimeFlag) SetOutputFormat(format string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.outputFormat = format
}

// Init 初始化时间类型标志，支持字符串类型默认值
//
// 参数:
//   - longName: 长标志名称
//   - shortName: 短标志字符
//   - defValue: 默认值（字符串格式，支持多种时间表达）
//   - usage: 帮助说明
//
// 返回值:
//   - error: 初始化错误信息
//
// 支持的默认值格式:
//   - "now" 或 "" : 当前时间
//   - "zero" : 零时间 (time.Time{})
//   - "1h", "30m", "-2h" : 相对时间（基于当前时间的偏移）
//   - "2006-01-02", "2006-01-02 15:04:05" : 绝对时间格式
//   - RFC3339等标准格式
func (f *TimeFlag) Init(longName, shortName string, defValue string, usage string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 解析字符串默认值为 time.Time
	parsedTime, err := f.parseTimeString(defValue)
	if err != nil {
		return qerr.NewValidationErrorf("invalid default time value '%s': %v", defValue, err)
	}

	// 创建时间值指针
	timePtr := new(time.Time)
	*timePtr = parsedTime

	// 调用基类的 Init 方法
	return f.BaseFlag.Init(longName, shortName, usage, timePtr)
}

// parseTimeString 解析时间字符串为 time.Time
//
// 参数:
//   - s: 时间字符串
//
// 返回值:
//   - time.Time: 解析后的时间
//   - error: 解析错误
func (f *TimeFlag) parseTimeString(s string) (time.Time, error) {
	// 处理特殊值
	switch strings.ToLower(strings.TrimSpace(s)) {
	// now 或 空字符串 表示当前时间
	case "now", "":
		return time.Now(), nil
	// zero 表示零时间
	case "zero":
		return time.Time{}, nil
	}

	// 尝试解析相对时间（如 "1h", "-30m", "2h30m"）
	if d, err := time.ParseDuration(s); err == nil {
		return time.Now().Add(d), nil
	}

	// 尝试解析绝对时间格式
	for _, format := range supportedTimeFormats {
		if t, err := time.Parse(format, s); err == nil {
			return t, nil
		}
	}

	// 如果都解析失败，返回错误
	return time.Time{}, qerr.NewValidationError("unsupported time format, supported formats: 'now', 'zero', duration (1h, 30m), or standard time formats")
}

// =============================================================================
// 时间间隔类型标志
// =============================================================================

// DurationFlag 时间间隔类型标志结构体
// 继承BaseFlag[time.Duration]泛型结构体,实现Flag接口
type DurationFlag struct {
	BaseFlag[time.Duration]
	mu sync.Mutex // 互斥锁, 用于保护并发访问
}

// Type 返回标志类型
//
// 返回值:
//   - FlagType: 标志类型枚举值
func (f *DurationFlag) Type() FlagType { return FlagTypeDuration }

// Set 实现flag.Value接口, 解析并设置时间间隔值
//
// 参数:
//   - value: 待设置的值
//
// 返回值:
//   - error: 解析或验证失败时返回错误信息
func (f *DurationFlag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 检查空值
	if value == "" {
		return qerr.NewValidationError("duration cannot be empty")
	}

	// 将单位转换为小写, 确保解析的准确性
	lowercaseValue := strings.ToLower(value)

	// 解析时间间隔字符串
	duration, err := time.ParseDuration(lowercaseValue)
	if err != nil {
		return qerr.NewValidationErrorf("invalid duration format: %v (valid units: ns/us/ms/s/m/h)", err)
	}

	// 检查负值（可选）
	if duration < 0 {
		return qerr.NewValidationError("negative duration not allowed")
	}

	// 调用基类方法设置值
	return f.BaseFlag.Set(duration)
}

// String 实现flag.Value接口, 返回当前值的字符串表示
//
// 返回值:
//   - string: 当前值的字符串表示
func (f *DurationFlag) String() string {
	return f.Get().String()
}

// =============================================================================
// 键值对类型标志
// =============================================================================

// MapFlag 键值对类型标志结构体
// 继承BaseFlag[map[string]string]泛型结构体,实现Flag接口
type MapFlag struct {
	BaseFlag[map[string]string]
	keyDelimiter   string       // 键值对之间的分隔符
	valueDelimiter string       // 键和值之间的分隔符
	mu             sync.RWMutex // 读写锁,保护并发访问
	ignoreCase     bool         // 是否忽略键的大小写
}

// SetIgnoreCase 设置是否忽略键的大小写
//
// 参数:
//   - enable: 是否忽略键的大小写
//
// 注意:
//   - 当enable为true时,所有键将转换为小写进行存储和比较
func (f *MapFlag) SetIgnoreCase(enable bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.ignoreCase = enable
}

// Type 返回标志类型
//
// 返回值:
//   - FlagType: 标志类型枚举值
func (f *MapFlag) Type() FlagType { return FlagTypeMap }

// String 实现flag.Value接口,返回当前值的字符串表示
//
// 返回值:
//   - string: 当前值的字符串表示
func (f *MapFlag) String() string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	m := f.Get()
	if m == nil {
		return ""
	}
	var parts []string
	for k, v := range m {
		parts = append(parts, fmt.Sprintf("%s%s%s", k, f.valueDelimiter, v))
	}
	return strings.Join(parts, f.keyDelimiter)
}

// Set 实现flag.Value接口,解析并设置键值对
//
// 参数:
//   - value: 待设置的值
//
// 返回值:
//   - error: 解析或验证失败时返回错误信息
func (f *MapFlag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if value == "" {
		return qerr.NewValidationError("map value cannot be empty")
	}

	// 确保分隔符已设置，如果没有设置则使用默认值
	keyDelim := f.keyDelimiter
	valueDelim := f.valueDelimiter
	if keyDelim == "" {
		keyDelim = FlagSplitComma // 默认使用逗号
	}
	if valueDelim == "" {
		valueDelim = FlagKVEqual // 默认使用等号
	}

	// 获取当前值
	current := f.Get()
	if current == nil {
		// 如果当前值为空, 初始化一个空map
		current = make(map[string]string)
	}

	// 简化的解析逻辑: 先按键分隔符分割，再处理每个键值对
	pairs := strings.Split(value, keyDelim)

	// 处理每个键值对
	for _, pair := range pairs {
		// 使用SplitN限制分割次数为2，这样值中可以包含值分隔符, 例如: "key=value,key2=value2"
		kv := strings.SplitN(pair, valueDelim, 2)

		// 检查键值对是否包含两个部分
		if len(kv) != 2 {
			return qerr.NewValidationErrorf("validation failed: invalid key-value pair format: %s", pair)
		}

		// 去除键和值的前后空格
		key := strings.TrimSpace(kv[0])
		val := strings.TrimSpace(kv[1])

		// 如果需要忽略大小写,则将键转换为小写
		if f.ignoreCase {
			key = strings.ToLower(key)
		}

		// 检查键是否为空
		if key == "" {
			return qerr.NewValidationErrorf("validation failed: empty key in key-value pair: %s", pair)
		}
		// 注意：空值是允许的，不需要检查 val == ""

		// 更新当前值
		current[key] = val
	}

	return f.BaseFlag.Set(current)
}

// SetDelimiters 设置键值对分隔符
//
// 参数：
//   - keyDelimiter 键值对分隔符
//   - valueDelimiter 键值分隔符
func (f *MapFlag) SetDelimiters(keyDelimiter, valueDelimiter string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if keyDelimiter == "" {
		keyDelimiter = FlagSplitComma // 默认使用逗号
	}
	if valueDelimiter == "" {
		valueDelimiter = FlagKVEqual // 默认使用等号
	}

	// 检查键分隔符和值分隔符是否相同
	if keyDelimiter == valueDelimiter {
		// 如果相同，使用默认值组合
		keyDelimiter = FlagSplitComma
		valueDelimiter = FlagKVEqual
	}

	// 设置分隔符
	f.keyDelimiter = keyDelimiter
	f.valueDelimiter = valueDelimiter
}

// =============================================================================
// 切片类型标志
// =============================================================================

// SliceFlag 切片类型标志结构体
// 继承BaseFlag[[]string]泛型结构体,实现Flag接口
type SliceFlag struct {
	BaseFlag[[]string]              // 基类
	delimiters         []string     // 分隔符
	mu                 sync.RWMutex // 读写锁
	skipEmpty          bool         // 是否跳过空元素
}

// Type 返回标志类型
func (f *SliceFlag) Type() FlagType { return FlagTypeSlice }

// String 实现flag.Value接口, 返回当前值的字符串表示
func (f *SliceFlag) String() string {
	return strings.Join(f.Get(), ",")
}

// Set 实现flag.Value接口, 解析并设置切片值
//
// 参数:
//   - value 待解析的切片值
//
// 注意:
//   - 如果切片中包含分隔符,则根据分隔符进行分割, 否则将整个值作为单个元素
//   - 例如: "a,b,c" -> ["a", "b", "c"]
func (f *SliceFlag) Set(value string) error {
	// 加读锁保护分隔符切片访问
	f.mu.Lock()
	defer f.mu.Unlock()

	// 检查空值
	if value == "" {
		return fmt.Errorf("slice cannot be empty")
	}

	// 存储分割后的元素
	var elements []string

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
	if f.skipEmpty {
		filtered := make([]string, 0, len(elements))
		for _, e := range elements {
			if e != "" {
				filtered = append(filtered, e)
			}
		}
		elements = filtered
	}

	// 调用基类方法设置值
	return f.BaseFlag.Set(elements)
}

// SetDelimiters 设置切片解析的分隔符列表
//
// 参数:
//   - delimiters 分隔符列表
func (f *SliceFlag) SetDelimiters(delimiters []string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 检查分隔符是否为空
	if len(delimiters) == 0 {
		// 使用默认分隔符（与Init保持一致）
		delimiters = FlagSplitSlice
	}

	// 更新分隔符
	f.delimiters = delimiters
}

// GetDelimiters 获取当前分隔符列表
func (f *SliceFlag) GetDelimiters() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	// 返回拷贝避免外部修改内部切片
	res := make([]string, len(f.delimiters))
	copy(res, f.delimiters)
	return res
}

// SetSkipEmpty 设置是否跳过空元素
//
// 参数:
//   - skip - 为true时跳过空元素, 为false时保留空元素
//
// 线程安全的空元素跳过更新
func (f *SliceFlag) SetSkipEmpty(skip bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.skipEmpty = skip
}

// Len 获取切片长度
//
// 返回:
//   - 获取切片长度
func (f *SliceFlag) Len() int {
	// 返回切片长度
	return len(f.Get())
}

// Contains 检查切片是否包含指定元素
//
// 参数:
//   - element 待检查的元素
//
// 返回:
//   - 若切片包含指定元素, 返回true, 否则返回false
//
// 注意:
//   - 当切片未设置值时,将使用默认值进行检查
func (f *SliceFlag) Contains(element string) bool {
	// 通过Get()获取当前值(已处理nil情况和线程安全)
	current := f.Get()

	// 加读锁保护分隔符切片访问
	f.mu.RLock()
	defer f.mu.RUnlock()

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
// 返回值:
//   - 操作成功返回nil, 否则返回错误信息
//
// 注意：
//   - 该方法会改变切片的指针
func (f *SliceFlag) Clear() error {
	// 使用BaseFlag的Set方法确保线程安全
	return f.BaseFlag.Set([]string{})
}

// Remove 从切片中移除指定元素（支持移除空字符串元素）
//
// 参数:
//   - element 待移除的元素（支持空字符串）
//
// 返回值:
//   - 操作成功返回nil, 否则返回错误信息
func (f *SliceFlag) Remove(element string) error {
	// 获取当前切片
	current := f.Get()

	// 加写锁保护切片访问
	f.mu.Lock()
	defer f.mu.Unlock()

	// 遍历当前切片，移除指定元素
	newSlice := []string{}
	for _, item := range current {
		if item != element {
			newSlice = append(newSlice, item)
		}
	}

	return f.BaseFlag.Set(newSlice)
}

// Sort 对切片进行排序
// 对当前切片标志的值进行原地排序，修改原切片内容
// 采用Go标准库的sort.Strings()函数进行字典序排序(按Unicode代码点升序排列)
//
// 注意：
//   - 排序会直接修改当前标志的值，而非返回新切片
//   - 排序区分大小写, 遵循Unicode代码点比较规则(如'A' < 'a' < 'z')
//   - 若切片未设置值，将使用默认值进行排序
//
// 返回值：
//   - 排序成功返回nil, 若排序过程中发生错误则返回错误信息
func (f *SliceFlag) Sort() error {
	current := f.Get()
	sort.Strings(current)
	return f.BaseFlag.Set(current)
}

// Init 初始化切片类型标志
//
// 参数:
//   - longName: 长标志名称
//   - shortName: 短标志字符
//   - defValue: 默认值（切片类型）
//   - usage: 帮助说明
//
// 返回值:
//   - error: 初始化错误信息
func (f *SliceFlag) Init(longName, shortName string, defValue []string, usage string) error {
	// 确保默认值不为nil
	if defValue == nil {
		defValue = []string{}
	}

	// 1. 初始化值指针（切片需创建副本避免外部修改影响）
	valueCopy := make([]string, len(defValue))
	copy(valueCopy, defValue)
	valuePtr := &valueCopy

	// 2. 调用基类初始化通用字段
	if err := f.BaseFlag.Init(longName, shortName, usage, valuePtr); err != nil {
		return err
	}

	// 3. 初始化切片特有字段(通过SetDelimiters保证线程安全)
	f.SetDelimiters(FlagSplitSlice)

	return nil
}

// =============================================================================
// 无符号整数类型标志
// =============================================================================

// Uint16Flag 16位无符号整数类型标志结构体
// 继承BaseFlag[uint16]泛型结构体,实现Flag接口
type Uint16Flag struct {
	BaseFlag[uint16]            // 基类
	mu               sync.Mutex // 互斥锁
}

// Type 返回标志类型
//
// 返回值:
//   - FlagType: 标志类型枚举值
func (f *Uint16Flag) Type() FlagType { return FlagTypeUint16 }

// String 实现flag.Value接口, 返回当前值的字符串表示
//
// 返回值:
//   - string: 当前值的字符串表示
func (f *Uint16Flag) String() string {
	return fmt.Sprint(f.Get())
}

// Set 实现flag.Value接口, 解析并设置16位无符号整数值
// 验证值是否在uint16范围内(0-65535)
//
// 参数:
//   - value: 待设置的值(0-65535)
//
// 返回值:
//   - error: 解析或验证失败时返回错误信息
func (f *Uint16Flag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 检查是否为空
	if value == "" {
		return qerr.NewValidationError("empty value")
	}

	// 解析字符串为uint64
	num, err := strconv.ParseUint(value, 10, 16)
	if err != nil {
		return qerr.NewValidationErrorf("invalid uint16 value: %v", err)
	}

	// 转换为uint16
	val := uint16(num)

	// 调用基类方法设置值
	return f.BaseFlag.Set(val)
}

// Uint32Flag 32位无符号整数类型标志结构体
// 继承BaseFlag[uint32]泛型结构体,实现Flag接口
type Uint32Flag struct {
	BaseFlag[uint32]            // 基类
	mu               sync.Mutex // 互斥锁
}

// Type 返回标志类型
func (f *Uint32Flag) Type() FlagType { return FlagTypeUint32 }

// String 实现flag.Value接口, 返回当前值的字符串表示
func (f *Uint32Flag) String() string {
	return fmt.Sprint(f.Get())
}

// Set 实现flag.Value接口, 解析并设置32位无符号整数值
// 验证值是否在uint32范围内(0-4294967295)
//
// 参数:
//   - value: 待设置的值(0-4294967295)
//
// 返回值:
//   - error: 解析或验证失败时返回错误信息
func (f *Uint32Flag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 检查是否为空
	if value == "" {
		return qerr.NewValidationError("empty value")
	}

	// 将字符串解析为无符号整型
	num, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return qerr.NewValidationErrorf("invalid uint32 value: %v", err)
	}

	val := uint32(num)
	return f.BaseFlag.Set(val)
}

// Uint64Flag 64位无符号整数类型标志结构体
// 继承BaseFlag[uint64]泛型结构体,实现Flag接口
type Uint64Flag struct {
	BaseFlag[uint64]            // 基类
	mu               sync.Mutex // 互斥锁
}

// Type 返回标志类型
func (f *Uint64Flag) Type() FlagType { return FlagTypeUint64 }

// String 实现flag.Value接口, 返回当前值的字符串表示
func (f *Uint64Flag) String() string {
	return fmt.Sprint(f.Get())
}

// Set 实现flag.Value接口, 解析并设置64位无符号整数值
// 验证值是否在uint64范围内(0-18446744073709551615)
//
// 参数:
//   - value: 待设置的值(0-18446744073709551615)
//
// 返回值:
//   - error: 解析或验证失败时返回错误信息
func (f *Uint64Flag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 检查是否为空
	if value == "" {
		return qerr.NewValidationError("empty value")
	}

	// 将字符串解析为无符号整型
	num, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return qerr.NewValidationErrorf("invalid uint64 value: %v", err)
	}

	val := uint64(num)
	return f.BaseFlag.Set(val)
}

// Float64Flag 浮点型标志结构体
// 继承BaseFlag[float64]泛型结构体,实现Flag接口
type Float64Flag struct {
	BaseFlag[float64]
	mu sync.Mutex
}

// =============================================================================
// 枚举类型标志
// =============================================================================

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
