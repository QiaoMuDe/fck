package flags

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// =============================================================================
// 切片类型标志
// =============================================================================

// StringSliceFlag 切片类型标志结构体
// 继承BaseFlag[[]string]泛型结构体,实现Flag接口
type StringSliceFlag struct {
	BaseFlag[[]string]              // 基类
	delimiters         []string     // 分隔符
	mu                 sync.RWMutex // 读写锁
	skipEmpty          bool         // 是否跳过空元素
	initOnce           sync.Once    // 确保只初始化一次
}

// Type 返回标志类型
func (f *StringSliceFlag) Type() FlagType { return FlagTypeStringSlice }

// String 实现flag.Value接口, 返回当前值的字符串表示
func (f *StringSliceFlag) String() string {
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
func (f *StringSliceFlag) Set(value string) error {
	// 加读锁保护分隔符切片访问
	f.mu.RLock()
	defer f.mu.RUnlock()

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
func (f *StringSliceFlag) SetDelimiters(delimiters []string) {
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
func (f *StringSliceFlag) GetDelimiters() []string {
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
func (f *StringSliceFlag) SetSkipEmpty(skip bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.skipEmpty = skip
}

// Len 获取切片长度
//
// 返回:
//   - 获取切片长度
func (f *StringSliceFlag) Len() int {
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
func (f *StringSliceFlag) Contains(element string) bool {
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
// 返回值:
//   - 操作成功返回nil, 否则返回错误信息
//
// 注意：
//   - 该方法会改变切片的指针
func (f *StringSliceFlag) Clear() error {
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
func (f *StringSliceFlag) Remove(element string) error {
	// 获取当前切片
	current := f.Get()

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
func (f *StringSliceFlag) Sort() error {
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
//
// 注意: 重复调用此方法是安全的，后续调用将被忽略
func (f *StringSliceFlag) Init(longName, shortName string, defValue []string, usage string) error {
	var initErr error
	f.initOnce.Do(func() {
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
			initErr = err
			return
		}

		// 3. 初始化切片特有字段(通过SetDelimiters保证线程安全)
		f.SetDelimiters(FlagSplitSlice)
	})
	return initErr
}

// =============================================================================
// 64位整数切片类型标志
// =============================================================================

// Int64SliceFlag 64位整数切片类型标志结构体
// 继承BaseFlag[[]int64]泛型结构体,实现Flag接口
type Int64SliceFlag struct {
	BaseFlag[[]int64]              // 基类
	delimiters        []string     // 分隔符
	mu                sync.RWMutex // 读写锁
	skipEmpty         bool         // 是否跳过空元素
	initOnce          sync.Once    // 确保只初始化一次
}

// Type 返回标志类型
func (f *Int64SliceFlag) Type() FlagType { return FlagTypeInt64Slice }

// String 实现flag.Value接口, 返回当前值的字符串表示
func (f *Int64SliceFlag) String() string {
	values := f.Get()
	strValues := make([]string, len(values))
	for i, v := range values {
		strValues[i] = strconv.FormatInt(v, 10)
	}
	return strings.Join(strValues, ",")
}

// Set 实现flag.Value接口, 解析并设置64位整数切片值
//
// 参数:
//   - value 待解析的切片值
//
// 注意:
//   - 如果切片中包含分隔符,则根据分隔符进行分割, 否则将整个值作为单个元素
//   - 例如: "1,2,3" -> [1, 2, 3]
func (f *Int64SliceFlag) Set(value string) error {
	// 加读锁保护分隔符切片访问
	f.mu.RLock()
	defer f.mu.RUnlock()

	// 检查空值
	if value == "" {
		return fmt.Errorf("int64 slice cannot be empty")
	}

	// 存储分割后的字符串元素
	var strElements []string

	// 检查是否包含分隔符切片中的任何分隔符
	found := false
	for _, delimiter := range f.delimiters {
		if strings.Contains(value, delimiter) {
			// 根据分隔符分割字符串
			strElements = strings.Split(value, delimiter)
			// 去除每个元素的首尾空白字符
			for i, e := range strElements {
				strElements[i] = strings.TrimSpace(e)
			}
			found = true
			break // 找到第一个匹配的分隔符后停止
		}
	}

	// 如果没有找到分隔符,将整个值作为单个元素
	if !found {
		strElements = []string{strings.TrimSpace(value)}
	}

	// 过滤空元素（如果启用）
	if f.skipEmpty {
		filtered := make([]string, 0, len(strElements))
		for _, e := range strElements {
			if e != "" {
				filtered = append(filtered, e)
			}
		}
		strElements = filtered
	}

	// 转换字符串元素为64位整数
	elements := make([]int64, 0, len(strElements))
	for _, strElement := range strElements {
		if strElement == "" && f.skipEmpty {
			continue // 跳过空元素
		}

		int64Val, err := strconv.ParseInt(strElement, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid int64 value '%s': %v", strElement, err)
		}
		elements = append(elements, int64Val)
	}

	// 调用基类方法设置值
	return f.BaseFlag.Set(elements)
}

// SetDelimiters 设置切片解析的分隔符列表
//
// 参数:
//   - delimiters 分隔符列表
func (f *Int64SliceFlag) SetDelimiters(delimiters []string) {
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
func (f *Int64SliceFlag) GetDelimiters() []string {
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
func (f *Int64SliceFlag) SetSkipEmpty(skip bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.skipEmpty = skip
}

// Len 获取切片长度
//
// 返回:
//   - 获取切片长度
func (f *Int64SliceFlag) Len() int {
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
func (f *Int64SliceFlag) Contains(element int64) bool {
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
// 返回值:
//   - 操作成功返回nil, 否则返回错误信息
//
// 注意：
//   - 该方法会改变切片的指针
func (f *Int64SliceFlag) Clear() error {
	// 使用BaseFlag的Set方法确保线程安全
	return f.BaseFlag.Set([]int64{})
}

// Remove 从切片中移除指定元素
//
// 参数:
//   - element 待移除的元素
//
// 返回值:
//   - 操作成功返回nil, 否则返回错误信息
func (f *Int64SliceFlag) Remove(element int64) error {
	// 获取当前切片
	current := f.Get()

	// 遍历当前切片，移除指定元素
	newSlice := []int64{}
	for _, item := range current {
		if item != element {
			newSlice = append(newSlice, item)
		}
	}

	return f.BaseFlag.Set(newSlice)
}

// Sort 对切片进行排序
// 对当前切片标志的值进行原地排序，修改原切片内容
// 采用Go标准库的sort.Slice()函数进行数值升序排序
//
// 注意：
//   - 排序会直接修改当前标志的值，而非返回新切片
//   - 若切片未设置值，将使用默认值进行排序
//
// 返回值：
//   - 排序成功返回nil, 若排序过程中发生错误则返回错误信息
func (f *Int64SliceFlag) Sort() error {
	current := f.Get()
	sort.Slice(current, func(i, j int) bool {
		return current[i] < current[j]
	})
	return f.BaseFlag.Set(current)
}

// Init 初始化64位整数切片类型标志
//
// 参数:
//   - longName: 长标志名称
//   - shortName: 短标志字符
//   - defValue: 默认值（64位整数切片类型）
//   - usage: 帮助说明
//
// 返回值:
//   - error: 初始化错误信息
//
// 注意: 重复调用此方法是安全的，后续调用将被忽略
func (f *Int64SliceFlag) Init(longName, shortName string, defValue []int64, usage string) error {
	var initErr error
	f.initOnce.Do(func() {
		// 确保默认值不为nil
		if defValue == nil {
			defValue = []int64{}
		}

		// 1. 初始化值指针（切片需创建副本避免外部修改影响）
		valueCopy := make([]int64, len(defValue))
		copy(valueCopy, defValue)
		valuePtr := &valueCopy

		// 2. 调用基类初始化通用字段
		if err := f.BaseFlag.Init(longName, shortName, usage, valuePtr); err != nil {
			initErr = err
			return
		}

		// 3. 初始化切片特有字段(通过SetDelimiters保证线程安全)
		f.SetDelimiters(FlagSplitSlice)
	})
	return initErr
}

// =============================================================================
// 整数切片类型标志
// =============================================================================

// IntSliceFlag 整数切片类型标志结构体
// 继承BaseFlag[[]int]泛型结构体,实现Flag接口
type IntSliceFlag struct {
	BaseFlag[[]int]              // 基类
	delimiters      []string     // 分隔符
	mu              sync.RWMutex // 读写锁
	skipEmpty       bool         // 是否跳过空元素
	initOnce        sync.Once    // 确保只初始化一次
}

// Type 返回标志类型
func (f *IntSliceFlag) Type() FlagType { return FlagTypeIntSlice }

// String 实现flag.Value接口, 返回当前值的字符串表示
func (f *IntSliceFlag) String() string {
	values := f.Get()
	strValues := make([]string, len(values))
	for i, v := range values {
		strValues[i] = strconv.Itoa(v)
	}
	return strings.Join(strValues, ",")
}

// Set 实现flag.Value接口, 解析并设置整数切片值
//
// 参数:
//   - value 待解析的切片值
//
// 注意:
//   - 如果切片中包含分隔符,则根据分隔符进行分割, 否则将整个值作为单个元素
//   - 例如: "1,2,3" -> [1, 2, 3]
func (f *IntSliceFlag) Set(value string) error {
	// 加读锁保护分隔符切片访问
	f.mu.RLock()
	defer f.mu.RUnlock()

	// 检查空值
	if value == "" {
		return fmt.Errorf("int slice cannot be empty")
	}

	// 存储分割后的字符串元素
	var strElements []string

	// 检查是否包含分隔符切片中的任何分隔符
	found := false
	for _, delimiter := range f.delimiters {
		if strings.Contains(value, delimiter) {
			// 根据分隔符分割字符串
			strElements = strings.Split(value, delimiter)
			// 去除每个元素的首尾空白字符
			for i, e := range strElements {
				strElements[i] = strings.TrimSpace(e)
			}
			found = true
			break // 找到第一个匹配的分隔符后停止
		}
	}

	// 如果没有找到分隔符,将整个值作为单个元素
	if !found {
		strElements = []string{strings.TrimSpace(value)}
	}

	// 过滤空元素（如果启用）
	if f.skipEmpty {
		filtered := make([]string, 0, len(strElements))
		for _, e := range strElements {
			if e != "" {
				filtered = append(filtered, e)
			}
		}
		strElements = filtered
	}

	// 转换字符串元素为整数
	elements := make([]int, 0, len(strElements))
	for _, strElement := range strElements {
		if strElement == "" && f.skipEmpty {
			continue // 跳过空元素
		}

		intVal, err := strconv.Atoi(strElement)
		if err != nil {
			return fmt.Errorf("invalid integer value '%s': %v", strElement, err)
		}
		elements = append(elements, intVal)
	}

	// 调用基类方法设置值
	return f.BaseFlag.Set(elements)
}

// SetDelimiters 设置切片解析的分隔符列表
//
// 参数:
//   - delimiters 分隔符列表
func (f *IntSliceFlag) SetDelimiters(delimiters []string) {
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
func (f *IntSliceFlag) GetDelimiters() []string {
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
func (f *IntSliceFlag) SetSkipEmpty(skip bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.skipEmpty = skip
}

// Len 获取切片长度
//
// 返回:
//   - 获取切片长度
func (f *IntSliceFlag) Len() int {
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
func (f *IntSliceFlag) Contains(element int) bool {
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
// 返回值:
//   - 操作成功返回nil, 否则返回错误信息
//
// 注意：
//   - 该方法会改变切片的指针
func (f *IntSliceFlag) Clear() error {
	// 使用BaseFlag的Set方法确保线程安全
	return f.BaseFlag.Set([]int{})
}

// Remove 从切片中移除指定元素
//
// 参数:
//   - element 待移除的元素
//
// 返回值:
//   - 操作成功返回nil, 否则返回错误信息
func (f *IntSliceFlag) Remove(element int) error {
	// 获取当前切片
	current := f.Get()

	// 遍历当前切片，移除指定元素
	newSlice := []int{}
	for _, item := range current {
		if item != element {
			newSlice = append(newSlice, item)
		}
	}

	return f.BaseFlag.Set(newSlice)
}

// Sort 对切片进行排序
// 对当前切片标志的值进行原地排序，修改原切片内容
// 采用Go标准库的sort.Ints()函数进行数值升序排序
//
// 注意：
//   - 排序会直接修改当前标志的值，而非返回新切片
//   - 若切片未设置值，将使用默认值进行排序
//
// 返回值：
//   - 排序成功返回nil, 若排序过程中发生错误则返回错误信息
func (f *IntSliceFlag) Sort() error {
	current := f.Get()
	sort.Ints(current)
	return f.BaseFlag.Set(current)
}

// Init 初始化整数切片类型标志
//
// 参数:
//   - longName: 长标志名称
//   - shortName: 短标志字符
//   - defValue: 默认值（整数切片类型）
//   - usage: 帮助说明
//
// 返回值:
//   - error: 初始化错误信息
//
// 注意: 重复调用此方法是安全的，后续调用将被忽略
func (f *IntSliceFlag) Init(longName, shortName string, defValue []int, usage string) error {
	var initErr error
	f.initOnce.Do(func() {
		// 确保默认值不为nil
		if defValue == nil {
			defValue = []int{}
		}

		// 1. 初始化值指针（切片需创建副本避免外部修改影响）
		valueCopy := make([]int, len(defValue))
		copy(valueCopy, defValue)
		valuePtr := &valueCopy

		// 2. 调用基类初始化通用字段
		if err := f.BaseFlag.Init(longName, shortName, usage, valuePtr); err != nil {
			initErr = err
			return
		}

		// 3. 初始化切片特有字段(通过SetDelimiters保证线程安全)
		f.SetDelimiters(FlagSplitSlice)
	})

	return initErr
}
