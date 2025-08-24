// Package utils 提供缓冲区管理功能，通过对象池优化内存使用。
//
// 该文件实现了基于 sync.Pool 的缓冲区对象池，用于减少频繁的内存分配和回收。
// 通过复用缓冲区，可以显著提升文件读写操作的性能，特别是在处理大量小文件时。
//
// 主要功能：
//   - 缓冲区对象池管理
//   - 动态大小缓冲区获取
//   - 自动内存回收控制
//   - 防止内存泄漏的大小限制
//
// 性能优化：
//   - 使用 sync.Pool 减少 GC 压力
//   - 支持不同大小的缓冲区需求
//   - 自动限制大缓冲区回收
//
// 使用示例：
//
//	// 获取缓冲区
//	buffer := utils.GetBuffer(64 * 1024)
//
//	// 使用缓冲区进行文件操作
//	_, err := io.CopyBuffer(dst, src, buffer)
//
//	// 归还缓冲区到对象池
//	utils.PutBuffer(buffer)
package utils

import "sync"

// 缓冲区对象池，复用缓冲区减少内存分配
var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 32*1024) // 默认32KB缓冲区
	},
}

// GetBuffer 从对象池获取缓冲区
//
// 参数:
//   - size: 缓冲区大小
//
// 返回值:
//   - []byte: 获取到的缓冲区
func GetBuffer(size int) []byte {
	buffer, ok := bufferPool.Get().([]byte)
	if !ok || len(buffer) < size {
		// 如果类型断言失败或池中的缓冲区太小，创建新的
		return make([]byte, size)
	}
	return buffer[:size]
}

// PutBuffer 将缓冲区归还到对象池
//
// 参数:
//   - buffer: 要归还的缓冲区
//
// 说明:
//   - 该函数将缓冲区归还到对象池，以便后续复用。
//   - 只有容量不超过1MB的缓冲区才会被归还，以避免对象池占用过多内存。
func PutBuffer(buffer []byte) {
	if cap(buffer) <= 1024*1024 { // 只回收不超过1MB的缓冲区
		//nolint:staticcheck // SA6002: 忽略装箱警告，对象池的性能收益远大于装箱开销
		bufferPool.Put(buffer)
	}
}
