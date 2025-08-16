// Package buffer 提供了高性能的字符串构建器对象池实现。
// 该文件实现了 BuilderPool 结构体，用于管理 strings.Builder 对象的复用，
// 支持智能缩容、统计信息收集等功能，有效减少内存分配和垃圾回收压力。
package buffer

import (
	"strings"
	"sync"
)

// BuilderPool 字符串构建器对象池实现
type BuilderPool struct {
	pool    *sync.Pool // 对象池
	maxSize int        // 最大构建器大小限制
}

// NewBuilderPool 创建新的字符串构建器池
//
// 参数:
//   - maxSize: 最大构建器大小限制，超过此大小的构建器不会被放回池中
//
// 返回值:
//   - *BuilderPool: 字符串构建器池实例
func NewBuilderPool(maxSize int) *BuilderPool {
	// 创建新的字符串构建器池实例
	bp := &BuilderPool{
		maxSize: maxSize, // 最大构建器大小限制
	}

	// 创建对象池实例
	bp.pool = &sync.Pool{
		New: func() any {
			return &strings.Builder{} // 返回新的字符串构建器实例
		},
	}

	return bp
}

// 常量定义默认大小
const (
	DefaultMaxSize = 64 * 1024 // 64KB 默认最大大小
)

// NewDefaultBuilderPool 创建默认配置的字符串构建器池
//
// 返回值:
//   - *BuilderPool: 使用默认配置的字符串构建器池实例
func NewDefaultBuilderPool() *BuilderPool {
	return NewBuilderPool(DefaultMaxSize)
}

// Get 从池中获取字符串构建器
//
// 返回值:
//   - *strings.Builder: 干净的字符串构建器实例
func (bp *BuilderPool) Get() *strings.Builder {
	obj := bp.pool.Get()                  // 从对象池获取对象
	builder, ok := obj.(*strings.Builder) // 类型断言
	if !ok {
		// 断言失败，创建新的构建器
		builder = &strings.Builder{}
	}

	// 重置构建器
	builder.Reset()
	return builder
}

// Put 将字符串构建器归还到池中（智能缩容优化版）
//
// 参数:
//   - builder: 要归还的字符串构建器
func (bp *BuilderPool) Put(builder *strings.Builder) {
	if builder == nil {
		return
	}

	// 正常大小直接归还
	if bp.maxSize <= 0 {
		builder.Reset()
		bp.pool.Put(builder)
		return
	}

	// 小容量对象直接归还
	cap := builder.Cap()
	if cap <= bp.maxSize {
		builder.Reset()
		bp.pool.Put(builder)
		return
	}

	// 智能缩容: 超大对象处理
	if cap <= bp.maxSize<<1 { // 2倍以内，使用位运算优化
		// 创建小容量替代品，避免池变空
		newBuilder := &strings.Builder{}
		// 使用更小的替代品容量，减少内存占用
		newBuilder.Grow(bp.maxSize >> 3) // maxSize/8
		newBuilder.Reset()               // 重置新构建器
		bp.pool.Put(newBuilder)          // 归还小容量替代品
	}
	// 超大对象直接丢弃，让GC处理
}

// WithBuilder 使用字符串构建器执行函数
//
// 参数:
//   - fn: 要执行的函数，接收字符串构建器作为参数
func (bp *BuilderPool) WithBuilder(fn func(*strings.Builder)) {
	builder := bp.Get()   // 从对象池获取对象
	defer bp.Put(builder) // 函数结束后归还对象
	fn(builder)           // 执行函数
}

// BuildString 构建字符串
//
// 参数:
//   - fn: 构建函数，接收字符串构建器作为参数
//
// 返回值:
//   - string: 构建的字符串
func (bp *BuilderPool) BuildString(fn func(*strings.Builder)) string {
	builder := bp.Get()     // 从对象池获取对象
	defer bp.Put(builder)   // 函数结束后归还对象
	fn(builder)             // 执行函数
	return builder.String() // 返回构建的字符串

}
