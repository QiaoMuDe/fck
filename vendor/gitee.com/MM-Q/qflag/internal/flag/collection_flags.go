package flag

import (
	"fmt"
	"strconv"
	"strings"

	"gitee.com/MM-Q/qflag/internal/types"
)

// StringSliceFlag 字符串切片标志
type StringSliceFlag struct {
	*BaseFlag[[]string]
}

// NewStringSliceFlag 创建新的字符串切片标志
func NewStringSliceFlag(longName, shortName, desc string, default_ []string) *StringSliceFlag {
	return &StringSliceFlag{
		BaseFlag: NewBaseFlag(types.FlagTypeStringSlice, longName, shortName, desc, default_),
	}
}

// Set 设置字符串切片标志的值
func (f *StringSliceFlag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 处理空字符串，设置为空切片（不验证）
	if value == "" {
		*f.value = []string{}
		f.isSet = true
		return nil
	}

	// 使用逗号分割字符串
	parts := strings.Split(value, ",")

	// 过滤掉空字符串元素
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	// 验证（如果设置了验证器）
	if f.validator != nil {
		if err := f.validator(result); err != nil {
			return err
		}
	}

	// 设置值并标记为已设置
	*f.value = result
	f.isSet = true

	return nil
}

// Length 获取切片长度
func (f *StringSliceFlag) Length() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(*f.value)
}

// IsEmpty 检查切片是否为空
func (f *StringSliceFlag) IsEmpty() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(*f.value) == 0
}

// IntSliceFlag 整数切片标志
type IntSliceFlag struct {
	*BaseFlag[[]int]
}

// NewIntSliceFlag 创建新的整数切片标志
func NewIntSliceFlag(longName, shortName, desc string, default_ []int) *IntSliceFlag {
	return &IntSliceFlag{
		BaseFlag: NewBaseFlag(types.FlagTypeIntSlice, longName, shortName, desc, default_),
	}
}

// Set 设置整数切片标志的值
func (f *IntSliceFlag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 处理空字符串，设置为空切片（不验证）
	if value == "" {
		*f.value = []int{}
		f.isSet = true
		return nil
	}

	// 使用逗号分割字符串
	parts := strings.Split(value, ",")

	// 转换为整数
	result := make([]int, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)

		// 跳过空字符串
		if part == "" {
			continue
		}

		n, err := strconv.Atoi(part)
		if err != nil {
			return types.WrapParseError(err, "int slice", part)
		}
		result = append(result, n)
	}

	// 验证（如果设置了验证器）
	if f.validator != nil {
		if err := f.validator(result); err != nil {
			return err
		}
	}

	// 设置值并标记为已设置
	*f.value = result
	f.isSet = true

	return nil
}

// Length 获取切片长度
func (f *IntSliceFlag) Length() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(*f.value)
}

// IsEmpty 检查切片是否为空
func (f *IntSliceFlag) IsEmpty() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(*f.value) == 0
}

// Int64SliceFlag 64位整数切片标志
type Int64SliceFlag struct {
	*BaseFlag[[]int64]
}

// NewInt64SliceFlag 创建新的64位整数切片标志
func NewInt64SliceFlag(longName, shortName, desc string, default_ []int64) *Int64SliceFlag {
	return &Int64SliceFlag{
		BaseFlag: NewBaseFlag(types.FlagTypeInt64Slice, longName, shortName, desc, default_),
	}
}

// Set 设置64位整数切片标志的值
func (f *Int64SliceFlag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 处理空字符串，设置为空切片（不验证）
	if value == "" {
		*f.value = []int64{}
		f.isSet = true
		return nil
	}

	// 使用逗号分割字符串
	parts := strings.Split(value, ",")

	// 转换为64位整数
	result := make([]int64, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)

		// 跳过空字符串
		if part == "" {
			continue
		}

		n, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return types.WrapParseError(err, "int64 slice", part)
		}
		result = append(result, n)
	}

	// 验证（如果设置了验证器）
	if f.validator != nil {
		if err := f.validator(result); err != nil {
			return err
		}
	}

	// 设置值并标记为已设置
	*f.value = result
	f.isSet = true

	return nil
}

// Length 获取切片长度
func (f *Int64SliceFlag) Length() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(*f.value)
}

// IsEmpty 检查切片是否为空
func (f *Int64SliceFlag) IsEmpty() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(*f.value) == 0
}

// MapFlag 用于处理键值对映射类型的命令行参数。
// 支持的格式: key1=value1,key2=value2
//
// 空值处理:
//   - 空字符串 "" 表示创建空映射
//   - ",,," 中的空对会被跳过
//   - 使用 SetKV 方法设置键值对时, 键不能为空
//   - 使用 Clear 方法可以清空映射
type MapFlag struct {
	*BaseFlag[map[string]string]
}

// NewMapFlag 创建新的映射标志
//
// 参数:
//   - longName: 长选项名
//   - shortName: 短选项名
//   - desc: 标志描述
//   - default_: 默认值, 如果为nil则创建空映射
//
// 返回值:
//   - *MapFlag: 映射标志实例
func NewMapFlag(longName, shortName, desc string, default_ map[string]string) *MapFlag {
	// 确保默认值不是nil, 如果是nil则创建空map
	if default_ == nil {
		default_ = make(map[string]string)
	}

	return &MapFlag{
		BaseFlag: NewBaseFlag(types.FlagTypeMap, longName, shortName, desc, default_),
	}
}

// Set 设置映射标志的值
//
// 支持格式: key1=value1,key2=value2
//
// 空值处理:
//   - 空字符串 "" 表示创建空映射
//   - ",,," 中的空对会被跳过
//   - 键不能为空, 否则返回错误
//
// 参数:
//   - value: 映射字符串
//
// 返回值:
//   - error: 如果解析失败返回错误
func (f *MapFlag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 检查空字符串输入，设置为空映射（不验证）
	if value == "" {
		*f.value = make(map[string]string)
		f.isSet = true
		return nil
	}

	// 使用逗号分割字符串, 然后对每个键值对用等号分割
	result := make(map[string]string)
	pairs := strings.Split(value, ",")

	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			// 跳过空对, 但继续处理其他对
			continue
		}

		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return types.NewError("INVALID_MAP_FORMAT", fmt.Sprintf("invalid map format: %s", pair), nil)
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		if key == "" {
			return types.NewError("INVALID_MAP_KEY", fmt.Sprintf("empty key in map format: %s", pair), nil)
		}

		result[key] = val
	}

	// 验证（如果设置了验证器）
	if f.validator != nil {
		if err := f.validator(result); err != nil {
			return err
		}
	}

	// 设置值并标记为已设置
	*f.value = result
	f.isSet = true

	return nil
}

// Length 获取映射长度
func (f *MapFlag) Length() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(*f.value)
}

// IsEmpty 检查映射是否为空
func (f *MapFlag) IsEmpty() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(*f.value) == 0
}

// Clear 清空映射
// 将映射设置为空映射, 并标记为已设置
func (f *MapFlag) Clear() {
	f.mu.Lock()
	defer f.mu.Unlock()

	*f.value = make(map[string]string)
	f.isSet = true
}

// GetKey 获取映射中指定键的值
func (f *MapFlag) GetKey(key string) (string, bool) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if key == "" {
		return "", false
	}

	val, exists := (*f.value)[key]
	return val, exists
}

// HasKey 检查映射中是否包含指定键
func (f *MapFlag) HasKey(key string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if key == "" {
		return false
	}

	_, exists := (*f.value)[key]
	return exists
}

// Keys 获取映射的所有键
func (f *MapFlag) Keys() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	keys := make([]string, 0, len(*f.value))
	for k := range *f.value {
		keys = append(keys, k)
	}
	return keys
}
