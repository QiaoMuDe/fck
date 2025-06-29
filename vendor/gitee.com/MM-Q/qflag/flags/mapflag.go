package flags

import (
	"fmt"
	"strings"
	"sync"
)

// MapFlag 键值对类型标志结构体
// 继承BaseFlag[map[string]string]泛型结构体,实现Flag接口
type MapFlag struct {
	BaseFlag[map[string]string]
	keyDelimiter   string     // 键值对之间的分隔符
	valueDelimiter string     // 键和值之间的分隔符
	mu             sync.Mutex // 互斥锁
	IgnoreCase     bool       // 是否忽略键的大小写
}

// SetIgnoreCase 设置是否忽略键的大小写
// enable为true时，所有键将转换为小写进行存储和比较
func (f *MapFlag) SetIgnoreCase(enable bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.IgnoreCase = enable
}

// Type 返回标志类型
func (f *MapFlag) Type() FlagType { return FlagTypeMap }

// String 实现flag.Value接口,返回当前值的字符串表示
func (f *MapFlag) String() string {
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
func (f *MapFlag) Set(value string) error {
	if value == "" {
		return fmt.Errorf("map value cannot be empty")
	}

	// 获取当前值
	current := f.Get()
	if current == nil {
		current = make(map[string]string)
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	// 使用键分隔符分割多个键值对
	pairs := strings.Split(value, f.keyDelimiter)
	for _, pair := range pairs {
		// 使用值分隔符分割键和值
		kv := strings.SplitN(pair, f.valueDelimiter, 2)

		// 检查键值对是否包含两个部分
		if len(kv) != 2 {
			return fmt.Errorf("invalid key-value pair format: %s", pair)
		}

		// 去除键和值的前后空格
		key := strings.TrimSpace(kv[0])
		val := strings.TrimSpace(kv[1])

		// 如果需要忽略大小写,则将键转换为小写
		if f.IgnoreCase {
			key = strings.ToLower(key)
		}

		// 检查键和值是否为空
		if key == "" {
			return fmt.Errorf("empty key in key-value pair: %s", pair)
		}
		if val == "" {
			return fmt.Errorf("empty value in key-value pair: %s", pair)
		}

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

	// 设置分隔符
	f.keyDelimiter = keyDelimiter
	f.valueDelimiter = valueDelimiter
}
