package flags

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

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
		defValue = make([]string, 0)
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
