// Package cache 提供了ANSI序列缓存管理功能。
// 该文件实现了ANSICache结构体，用于缓存预构建的ANSI颜色序列，
// 支持预热、统计、清理等功能，显著提升颜色输出的性能。
package cache

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
)

// ANSICache ANSI序列缓存管理器
type ANSICache struct {
	cache    sync.Map           // 线程安全的缓存存储 map[CacheKey]string
	hitCount map[CacheKey]int64 // 命中次数统计
	mu       sync.RWMutex       // 保护hitCount的读写锁
	stats    CacheStats         // 缓存统计信息
}

// CacheStats 缓存统计信息
type CacheStats struct {
	TotalRequests atomic.Int64 // 总请求次数
	CacheHits     atomic.Int64 // 缓存命中次数
	CacheMisses   atomic.Int64 // 缓存未命中次数
}

// NewANSICache 创建新的ANSI序列缓存管理器
//
// 返回值:
//   - *ANSICache: ANSI缓存管理器实例
func NewANSICache() *ANSICache {
	cache := &ANSICache{
		hitCount: make(map[CacheKey]int64),
	}

	// 预热常用组合
	cache.preheat()
	return cache
}

// preheat 预热常用的ANSI序列组合
func (ac *ANSICache) preheat() {
	// 基于ColorLib使用模式分析的热点组合
	hotCombinations := []struct {
		color                  int
		bold, underline, blink bool
		description            string
	}{
		// 日志级别常用组合
		{31, true, false, false, "错误信息 - 红色粗体"},
		{32, true, false, false, "成功信息 - 绿色粗体"},
		{33, true, false, false, "警告信息 - 黄色粗体"},
		{34, false, false, false, "普通信息 - 蓝色"},
		{35, false, false, false, "调试信息 - 品红色"},

		// 基础颜色无样式组合
		{30, false, false, false, "黑色"},
		{31, false, false, false, "红色"},
		{32, false, false, false, "绿色"},
		{33, false, false, false, "黄色"},
		{34, false, false, false, "蓝色"},
		{35, false, false, false, "品红色"},
		{36, false, false, false, "青色"},
		{37, false, false, false, "白色"},
		{90, false, false, false, "灰色"},

		// 亮色系列
		{91, false, false, false, "亮红色"},
		{92, false, false, false, "亮绿色"},
		{93, false, false, false, "亮黄色"},
		{94, false, false, false, "亮蓝色"},
		{95, false, false, false, "亮品红色"},
		{96, false, false, false, "亮青色"},
		{97, false, false, false, "亮白色"},

		// 常用样式组合
		{31, false, true, false, "红色下划线"},
		{32, false, true, false, "绿色下划线"},
		{33, false, false, true, "黄色闪烁"},
		{31, true, true, false, "红色粗体下划线"},
	}

	for _, combo := range hotCombinations {
		key := BuildCacheKey(combo.color, combo.bold, combo.underline, combo.blink)
		ansi := ac.buildANSI(key)
		ac.cache.Store(key, ansi)
	}
}

// GetANSI 获取ANSI序列，优先从缓存获取
//
// 参数:
//   - colorCode: 颜色代码
//   - bold: 是否粗体
//   - underline: 是否下划线
//   - blink: 是否闪烁
//
// 返回值:
//   - string: ANSI序列字符串
func (ac *ANSICache) GetANSI(colorCode int, bold, underline, blink bool) string {
	key := BuildCacheKey(colorCode, bold, underline, blink)
	ac.stats.TotalRequests.Add(1)

	// 尝试从缓存获取
	if cached, ok := ac.cache.Load(key); ok {
		ac.recordHit(key)
		ac.stats.CacheHits.Add(1)
		return cached.(string)
	}

	// 缓存未命中，构建并缓存
	ac.stats.CacheMisses.Add(1)
	ansi := ac.buildANSI(key)
	ac.cache.Store(key, ansi)
	return ansi
}

// buildANSI 构建ANSI序列字符串
//
// 参数:
//   - key: 缓存键
//
// 返回值:
//   - string: 构建的ANSI序列
func (ac *ANSICache) buildANSI(key CacheKey) string {
	colorCode, bold, underline, blink := key.Parse()

	var builder strings.Builder
	builder.WriteString("\033[")

	// 收集所有需要的样式代码
	var codes []string

	// 添加样式代码
	if bold {
		codes = append(codes, "1")
	}
	if underline {
		codes = append(codes, "4")
	}
	if blink {
		codes = append(codes, "5")
	}

	// 添加颜色代码
	codes = append(codes, fmt.Sprintf("%d", colorCode))

	// 用分号连接所有代码
	builder.WriteString(strings.Join(codes, ";"))
	builder.WriteString("m")

	return builder.String()
}

// recordHit 记录缓存命中
//
// 参数:
//   - key: 缓存键
func (ac *ANSICache) recordHit(key CacheKey) {
	ac.mu.Lock()
	ac.hitCount[key]++
	ac.mu.Unlock()
}

// GetStats 获取缓存统计信息
//
// 返回值:
//   - *CacheStats: 缓存统计信息指针
func (ac *ANSICache) GetStats() *CacheStats {
	return &ac.stats
}

// GetHitRate 获取缓存命中率
//
// 返回值:
//   - float64: 命中率百分比 (0-100)
func (ac *ANSICache) GetHitRate() float64 {
	total := ac.stats.TotalRequests.Load()
	if total == 0 {
		return 0.0
	}
	hits := ac.stats.CacheHits.Load()
	return float64(hits) / float64(total) * 100.0
}

// GetCacheSize 获取当前缓存大小
//
// 返回值:
//   - int: 缓存中的条目数量
func (ac *ANSICache) GetCacheSize() int {
	count := 0
	ac.cache.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

// GetTopHits 获取命中次数最多的缓存键
//
// 参数:
//   - limit: 返回的最大数量
//
// 返回值:
//   - []CacheKey: 按命中次数排序的缓存键列表
func (ac *ANSICache) GetTopHits(limit int) []CacheKey {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	type keyHit struct {
		key CacheKey
		hit int64
	}

	var hits []keyHit
	for key, count := range ac.hitCount {
		hits = append(hits, keyHit{key, count})
	}

	// 简单的冒泡排序（因为数量不大）
	for i := 0; i < len(hits)-1; i++ {
		for j := 0; j < len(hits)-i-1; j++ {
			if hits[j].hit < hits[j+1].hit {
				hits[j], hits[j+1] = hits[j+1], hits[j]
			}
		}
	}

	// 返回前limit个
	result := make([]CacheKey, 0, limit)
	for i := 0; i < len(hits) && i < limit; i++ {
		result = append(result, hits[i].key)
	}

	return result
}

// Clear 清空缓存
func (ac *ANSICache) Clear() {
	ac.cache.Range(func(key, value interface{}) bool {
		ac.cache.Delete(key)
		return true
	})

	ac.mu.Lock()
	ac.hitCount = make(map[CacheKey]int64)
	ac.mu.Unlock()

	// 重置统计信息
	ac.stats.TotalRequests.Store(0)
	ac.stats.CacheHits.Store(0)
	ac.stats.CacheMisses.Store(0)
}

// String 返回缓存状态的字符串表示
//
// 返回值:
//   - string: 缓存状态信息
func (ac *ANSICache) String() string {
	return fmt.Sprintf("ANSICache{size:%d, hitRate:%.2f%%, requests:%d}",
		ac.GetCacheSize(),
		ac.GetHitRate(),
		ac.stats.TotalRequests.Load())
}
