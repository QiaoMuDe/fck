// Package find 实现了文件查找的模式匹配功能。
// 该文件提供了模式匹配器，支持文件名、路径、大小、时间等多种匹配条件，并包含正则表达式缓存机制。
package find

import (
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"gitee.com/MM-Q/fck/commands/internal/types"
)

// PatternMatcher 负责所有模式匹配逻辑，包含正则表达式缓存
type PatternMatcher struct {
	regexCache   map[string]*regexp.Regexp // 正则表达式缓存
	cacheMutex   sync.RWMutex              // 缓存读写锁
	maxCacheSize int                       // 最大缓存大小
}

// NewPatternMatcher 创建新的模式匹配器
func NewPatternMatcher(maxCacheSize int) *PatternMatcher {
	return &PatternMatcher{
		regexCache:   make(map[string]*regexp.Regexp),
		maxCacheSize: maxCacheSize,
	}
}

// GetRegex 获取编译好的正则表达式，使用缓存机制
//
// 参数:
//   - pattern: 正则表达式模式
//
// 返回:
//   - *regexp.Regexp: 编译后的正则表达式
//   - error: 编译错误（如果有）
func (m *PatternMatcher) GetRegex(pattern string) (*regexp.Regexp, error) {
	if pattern == "" {
		return nil, nil
	}

	// 读锁检查缓存
	m.cacheMutex.RLock()
	if compiled, exists := m.regexCache[pattern]; exists {
		m.cacheMutex.RUnlock()
		return compiled, nil
	}
	m.cacheMutex.RUnlock()

	// 写锁编译和缓存
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	// 双重检查
	if compiled, exists := m.regexCache[pattern]; exists {
		return compiled, nil
	}

	// 缓存大小限制
	if len(m.regexCache) >= m.maxCacheSize {
		// 清理最旧的条目（简单策略：删除第一个）
		for k := range m.regexCache {
			delete(m.regexCache, k)
			break
		}
	}

	// 编译正则表达式
	compiled, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	m.regexCache[pattern] = compiled
	return compiled, nil
}

// MatchName 匹配文件名
//
// 该函数负责匹配文件名与给定模式的逻辑
//
// 参数:
//   - name: 文件名
//   - pattern: 匹配模式
//   - config: 查找配置
//
// 返回:
//   - bool: 是否匹配成功
func (m *PatternMatcher) MatchName(name, pattern string, config *types.FindConfig) bool {
	return m.matchPattern(name, pattern, config.NameRegex, config)
}

// MatchPath 匹配路径
//
// 该函数负责匹配路径与给定模式的逻辑
//
// 参数:
//   - path: 路径
//   - pattern: 匹配模式
//   - config: 查找配置
//
// 返回:
//   - bool: 是否匹配成功
func (m *PatternMatcher) MatchPath(path, pattern string, config *types.FindConfig) bool {
	return m.matchPattern(path, pattern, config.PathRegex, config)
}

// matchPattern 通用匹配函数
// 该函数负责匹配输入字符串与给定模式的逻辑, 默认使用字符串匹配模式, 如果启用正则匹配模式, 则使用正则表达式匹配
//
// 参数:
//   - input: 输入字符串
//   - pattern: 匹配模式
//   - regex: 编译后的正则表达式
//   - config: 查找配置
//
// 返回:
//   - bool: 是否匹配成功
func (m *PatternMatcher) matchPattern(input, pattern string, regex *regexp.Regexp, config *types.FindConfig) bool {
	// 如果模式为空, 则不匹配
	if pattern == "" {
		return false
	}

	// 如果启用正则匹配, 使用正则表达式匹配
	if config.IsRegex {
		if regex == nil {
			return false
		}
		return regex.MatchString(input)
	}

	// 根据大小写敏感性处理字符串
	var s, p string
	if config.CaseSensitive {
		// 区分大小写
		s = input
		p = pattern
	} else {
		// 默认不区分大小写
		s = strings.ToLower(input)
		p = strings.ToLower(pattern)
	}

	// 全字匹配处理
	if config.WholeWord {
		return s == p
	}

	// 匹配模式处理
	return strings.Contains(s, p)
}

// MatchSize 检查文件大小是否符合指定的条件
//
// 该函数负责检查文件大小是否符合指定的条件, 支持字节、KB、MB、GB 单位
//
// 参数:
//   - fileSize: 文件大小
//   - sizeCondition: 大小条件, 格式如"+100"表示大于100, "-100"表示小于100
//
// 返回:
//   - bool: 是否匹配成功
func (m *PatternMatcher) MatchSize(fileSize int64, sizeCondition string) bool {
	// 检查大小条件是否为空
	if len(sizeCondition) < 2 {
		return false
	}

	// 获取比较符号和数值部分
	comparator := sizeCondition[0] // 比较符号
	sizeStr := sizeCondition[1:]   // 数值部分

	// 获取单位和数值部分
	unit := sizeStr[len(sizeStr)-1]          // 单位
	sizeValueStr := sizeStr[:len(sizeStr)-1] // 数值部分

	// 转换数值部分
	sizeValue, err := strconv.ParseFloat(sizeValueStr, 64)
	if err != nil {
		return false
	}

	// 根据单位转换为字节
	var sizeInBytes float64
	switch unit {
	case 'B', 'b':
		sizeInBytes = sizeValue
	case 'K', 'k':
		sizeInBytes = sizeValue * 1024
	case 'M', 'm':
		sizeInBytes = sizeValue * 1024 * 1024
	case 'G', 'g':
		sizeInBytes = sizeValue * 1024 * 1024 * 1024
	default:
		return false
	}

	// 根据比较符号进行比较
	switch comparator {
	case '+': // 大于
		return float64(fileSize) > sizeInBytes
	case '-': // 小于
		return float64(fileSize) < sizeInBytes
	default:
		return false
	}
}

// MatchTime 检查文件时间是否符合指定的条件
//
// 该函数负责检查文件时间是否符合指定的条件, 支持天单位
//
// 参数:
//   - fileTime: 文件时间
//   - timeCondition: 时间条件, 格式如"+10"表示10天前, "-10"表示10天后
//
// 返回:
//   - bool: 是否匹配成功
func (m *PatternMatcher) MatchTime(fileTime time.Time, timeCondition string) bool {
	// 检查时间条件是否为空
	if len(timeCondition) < 2 {
		return false
	}

	// 获取比较符号和数值部分
	comparator := timeCondition[0] // 比较符号
	daysStr := timeCondition[1:]   // 数值部分

	// 转换天数
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		return false
	}

	// 计算时间阈值
	threshold := time.Now().AddDate(0, 0, -days)

	// 根据比较符号进行比较
	switch comparator {
	case '+': // 大于
		return fileTime.After(threshold) // 检查文件时间是否在阈值之后
	case '-': // 小于
		return fileTime.Before(threshold) // 检查文件时间是否在阈值之前
	default:
		return false
	}
}

// ClearCache 清空正则表达式缓存
//
// 该函数负责清空正则表达式缓存, 释放内存
func (m *PatternMatcher) ClearCache() {
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()
	m.regexCache = make(map[string]*regexp.Regexp)
}

// GetCacheSize 获取当前缓存大小
func (m *PatternMatcher) GetCacheSize() int {
	m.cacheMutex.RLock()
	defer m.cacheMutex.RUnlock()
	return len(m.regexCache)
}
