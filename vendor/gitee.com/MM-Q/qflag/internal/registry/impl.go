// Package registry 提供泛型注册表实现
//
// registry 包基于 Go 泛型提供了通用的注册表实现, 支持任意类型的注册、查找和管理。
// 核心功能包括:
//   - 基于名称的注册和查找
//   - 支持短名称/别名的快速查找
//   - 线程安全的存储结构
//   - 遍历和批量操作支持
//
// 设计特点:
//   - 使用泛型实现类型安全
//   - 双索引结构 (名称索引和短名称索引)
//   - 统一的错误处理机制
//   - 支持类型特定的名称提取
package registry

import (
	"gitee.com/MM-Q/qflag/internal/types"
)

// registry 泛型注册表实现
//
// registry 是基于泛型的通用注册表实现, 支持任意类型的存储和管理。
// 使用数字索引结构提高查找效率, 支持长名称和短名称的快速定位。
//
// 字段说明:
//   - items: 主存储, 键为数字索引, 值为注册项 (只存储一次)
//   - nameIndex: 名称到索引的映射 (支持长名称和短名称, 都指向同一个索引)
//   - nextID: 下一个可用的ID (自增)
//
// 并发安全:
//   - 本身不提供并发安全, 使用时需外部同步
//   - 适用于单线程或已同步的场景
type registry[T any] struct {
	items     map[int]T      // 主存储, 使用数字索引存储实际对象 (只存储一次)
	nameIndex map[string]int // 名称到索引的映射 (支持长名称和短名称, 都指向同一个索引)
	nextID    int            // 下一个可用的ID (自增)
}

// NewRegistry 创建新的泛型注册表实例
//
// 返回值:
//   - *registry[T]: 初始化完成的注册表实例
//
// 功能说明:
//   - 初始化空的存储映射
//   - 创建名称索引映射
//   - 初始化ID计数器
//   - 准备接收注册项
func NewRegistry[T any]() *registry[T] {
	return &registry[T]{
		items:     make(map[int]T),
		nameIndex: make(map[string]int),
		nextID:    1, // 从1开始, 0表示无效ID
	}
}

// Register 注册项到注册表 (支持长名称和短名称)
//
// 参数:
//   - item: 要注册的项
//   - longName: 项的长名称
//   - shortName: 项的短名称
//
// 返回值:
//   - error: 注册失败时返回错误
//
// 错误情况:
//   - 长名称和短名称都为空: 返回 INVALID_NAME 错误
//   - 长名称已存在: 返回 ErrAlreadyExists 错误
//   - 短名称已存在: 返回 ErrAlreadyExists 错误
//
// 功能说明:
//   - 验证参数有效性
//   - 检查名称冲突
//   - 获取索引并存储到主存储
//   - 建立名称索引映射 (长名称和短名称都指向同一个索引)
func (r *registry[T]) Register(item T, longName, shortName string) error {
	// 第一步: 验证参数
	if longName == "" && shortName == "" {
		return types.NewError("INVALID_NAME", "long name and short name cannot both be empty", nil)
	}

	// 第二步: 检查名称冲突
	if longName != "" {
		if _, exists := r.nameIndex[longName]; exists {
			return types.ErrFlagAlreadyExists
		}
	}

	if shortName != "" {
		if _, exists := r.nameIndex[shortName]; exists {
			return types.ErrFlagAlreadyExists
		}
	}

	// 第三步: 获取索引并存储到主存储
	id := r.nextID
	r.nextID++
	r.items[id] = item // 直接获取索引, 把实际对象存储进去

	// 第四步: 建立名称索引映射
	// 通过判断长短名称如果不为空就注册一个, value就是实际的索引
	if longName != "" {
		r.nameIndex[longName] = id // 长名称映射到索引
	}

	if shortName != "" {
		r.nameIndex[shortName] = id // 短名称映射到索引
	}

	return nil
}

// Unregister 从注册表中移除指定项
//
// 参数:
//   - name: 要移除的项名称 (长名称或短名称)
//
// 返回值:
//   - error: 移除失败时返回错误
//
// 错误情况:
//   - 项不存在: 返回 ErrFlagNotFound 错误
//
// 功能说明:
//   - 通过名称找到对应的ID
//   - 找出所有指向该ID的名称并删除
//   - 删除主存储中的项
func (r *registry[T]) Unregister(name string) error {
	id, exists := r.nameIndex[name]
	if !exists {
		return types.ErrFlagNotFound
	}

	// 找出所有指向该ID的名称并删除
	for name, nameID := range r.nameIndex {
		if nameID == id {
			delete(r.nameIndex, name)
		}
	}

	// 删除项
	delete(r.items, id)

	return nil
}

// Get 通过名称获取项 (支持长名称或短名称)
//
// 参数:
//   - name: 项的名称 (长名称或短名称)
//
// 返回值:
//   - T: 找到的项
//   - bool: 是否找到, true表示找到
//
// 功能说明:
//   - 通过名称映射找到对应的ID
//   - 通过ID获取项
//   - 统一处理长名称和短名称
func (r *registry[T]) Get(name string) (T, bool) {
	id, exists := r.nameIndex[name]
	if !exists {
		var zero T
		return zero, false
	}

	item, exists := r.items[id]
	return item, exists
}

// List 获取所有注册项的列表
//
// 返回值:
//   - []T: 所有注册项的切片
//
// 功能说明:
//   - 遍历主映射中的所有项
//   - 返回新的切片, 不影响原数据
//   - 顺序不确定, 取决于map遍历顺序
func (r *registry[T]) List() []T {
	result := make([]T, 0, len(r.items))
	for _, item := range r.items {
		result = append(result, item)
	}
	return result
}

// Has 检查指定名称的项是否存在
//
// 参数:
//   - name: 要检查的名称 (长名称或短名称)
//
// 返回值:
//   - bool: 是否存在, true表示存在
//
// 功能说明:
//   - 快速检查名称索引
//   - 不返回项本身, 提高效率
func (r *registry[T]) Has(name string) bool {
	_, exists := r.nameIndex[name]
	return exists
}

// Count 获取注册表中的项数量
//
// 返回值:
//   - int: 注册项总数
//
// 功能说明:
//   - 直接返回主映射长度
//   - 时间复杂度O(1)
func (r *registry[T]) Count() int {
	return len(r.items)
}

// Clear 清空注册表中的所有项
//
// 功能说明:
//   - 重新初始化主映射
//   - 重新初始化名称索引
//   - 重置ID计数器
//   - 释放原有内存
func (r *registry[T]) Clear() {
	r.items = make(map[int]T)
	r.nameIndex = make(map[string]int)
	r.nextID = 1
}

// Range 遍历注册表中的所有项
//
// 参数:
//   - f: 遍历函数, 接收名称和项, 返回是否继续遍历
//
// 功能说明:
//   - 按名称索引顺序遍历
//   - 支持提前终止遍历
//   - 遍历过程中修改注册表可能导致不确定行为
func (r *registry[T]) Range(f func(name string, item T) bool) {
	for name, id := range r.nameIndex {
		if item, exists := r.items[id]; exists {
			if !f(name, item) {
				break
			}
		}
	}
}
