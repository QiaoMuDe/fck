// Package buffer 定义了字符串构建器池的接口和数据结构。
// 该文件包含了 Pool 接口定义和 PoolStats 统计信息结构体，
// 为字符串构建器池的实现提供了标准的接口规范。
package buffer

import "strings"

// Pool 字符串构建器池接口
type Pool interface {
	Get() *strings.Builder                     // 获取字符串构建器
	Put(*strings.Builder)                      // 归还字符串构建器
	WithBuilder(func(*strings.Builder))        // 使用字符串构建器
	BuildString(func(*strings.Builder)) string // 构建字符串
}
