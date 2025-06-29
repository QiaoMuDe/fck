package flags

import (
	"fmt"
	"sync"
)

// FlagMeta 统一存储标志的完整元数据
type FlagMeta struct {
	Flag Flag // 标志对象
}

// FlagMetaInterface 标志元数据接口, 定义了标志元数据的获取方法
type FlagMetaInterface interface {
	GetFlagType() FlagType // 获取标志类型
	GetFlag() Flag         // 获取标志对象
	GetLongName() string   // 获取标志的长名称
	GetShortName() string  // 获取标志的短名称
	GetUsage() string      // 获取标志的用法描述
	GetDefault() any       // 获取标志的默认值
	GetValue() any         // 获取标志的当前值
}

// GetLongName 获取标志的长名称
func (m *FlagMeta) GetLongName() string { return m.Flag.LongName() }

// GetShortName 获取标志的短名称
func (m *FlagMeta) GetShortName() string { return m.Flag.ShortName() }

// GetUsage 获取标志的用法描述
func (m *FlagMeta) GetUsage() string { return m.Flag.Usage() }

// GetFlagType 获取标志的类型
func (m *FlagMeta) GetFlagType() FlagType { return m.Flag.Type() }

// GetDefault 获取标志的默认值
func (m *FlagMeta) GetDefault() any { return m.Flag.GetDefaultAny() }

// GetFlag 获取标志对象
func (m *FlagMeta) GetFlag() Flag { return m.Flag }

// FlagRegistry 集中管理所有标志元数据及索引
type FlagRegistry struct {
	mu       sync.RWMutex         // 并发访问锁
	byLong   map[string]*FlagMeta // 按长名称索引
	byShort  map[string]*FlagMeta // 按短名称索引
	allFlags []*FlagMeta          // 所有标志元数据列表
}

// FlagRegistryInterface 标志注册表接口, 定义了标志元数据的增删改查操作
type FlagRegistryInterface interface {
	GetAllFlags() []*FlagMeta                      // 获取所有标志元数据列表
	GetLongFlags() map[string]*FlagMeta            // 获取长标志映射
	GetShortFlags() map[string]*FlagMeta           // 获取短标志映射
	RegisterFlag(meta *FlagMeta) error             // 注册一个新的标志元数据到注册表中
	GetByLong(longName string) (*FlagMeta, bool)   // 通过长标志名称查找对应的标志元数据
	GetByShort(shortName string) (*FlagMeta, bool) // 通过短标志名称查找对应的标志元数据
	GetByName(name string) (*FlagMeta, bool)       // 通过标志名称查找标志元数据
}

// 创建一个空的标志注册表
func NewFlagRegistry() *FlagRegistry {
	return &FlagRegistry{
		mu:       sync.RWMutex{},
		byLong:   make(map[string]*FlagMeta),
		byShort:  make(map[string]*FlagMeta),
		allFlags: make([]*FlagMeta, 0),
	}
}

// RegisterFlag 注册一个新的标志元数据到注册表中
// 该方法会执行以下操作:
// 1. 检查长名称和短名称是否已存在
// 2. 将标志添加到长名称索引
// 3. 将标志添加到短名称索引
// 4. 将标志添加到所有标志列表
// 注意: 该方法线程安全, 但发现重复标志时会panic
func (r *FlagRegistry) RegisterFlag(meta *FlagMeta) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 检查长短标志是否都为空
	if meta.GetLongName() == "" && meta.GetShortName() == "" {
		return fmt.Errorf("flag must have at least one name")
	}

	// 检查长标志是否已存在
	if meta.GetLongName() != "" {
		if _, exists := r.byLong[meta.GetLongName()]; exists {
			return fmt.Errorf("long flag %s already exists", meta.GetLongName())
		}
	}

	// 检查短标志是否已存在
	// 只在短标志不为空时检查重复
	if meta.GetShortName() != "" {
		if _, exists := r.byShort[meta.GetShortName()]; exists {
			return fmt.Errorf("short flag %s already exists", meta.GetShortName())
		}
	}

	// 添加长标志索引
	if meta.GetLongName() != "" {
		r.byLong[meta.GetLongName()] = meta
	}

	// 只在短标志不为空时添加短标志索引
	if meta.GetShortName() != "" {
		r.byShort[meta.GetShortName()] = meta
	}

	// 添加到所有标志列表
	r.allFlags = append(r.allFlags, meta)

	return nil
}

// GetByLong 通过长标志名称查找对应的标志元数据
// 参数:
//   - longName: 标志的长名称(如"help")
//
// 返回值:
//   - *FlagMeta: 找到的标志元数据指针, 未找到时为nil
//   - bool: 是否找到标志, true表示找到
func (r *FlagRegistry) GetByLong(longName string) (*FlagMeta, bool) {
	r.mu.RLock()                       // 获取读锁, 保证并发安全
	defer r.mu.RUnlock()               // 函数返回时释放读锁
	meta, exists := r.byLong[longName] // 从长名称索引中查找
	return meta, exists                // 返回查找结果
}

// GetByShort 通过短标志名称查找对应的标志元数据
// 参数:
//   - shortName: 标志的短名称(如"h"对应"help")
//
// 返回值:
//   - *FlagMeta: 找到的标志元数据指针, 未找到时为nil
//   - bool: 是否找到标志, true表示找到
func (r *FlagRegistry) GetByShort(shortName string) (*FlagMeta, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	meta, exists := r.byShort[shortName]
	return meta, exists
}

// GetByName 通过标志名称查找标志元数据
// 参数name可以是长名称(如"help")或短名称(如"h")
// 返回值:
//   - *FlagMeta: 找到的标志元数据指针, 未找到时为nil
//   - bool: 是否找到标志, true表示找到
func (r *FlagRegistry) GetByName(name string) (*FlagMeta, bool) {
	// 先尝试按长名称查找
	if meta, exists := r.GetByLong(name); exists {
		return meta, exists
	}

	// 再尝试按短名称查找
	if meta, exists := r.GetByShort(name); exists {
		return meta, exists
	}

	// 未找到
	return nil, false
}

// GetAllFlags 获取所有标志元数据列表
// 返回值:
//   - []*FlagMeta: 所有标志元数据的切片
func (r *FlagRegistry) GetAllFlags() []*FlagMeta {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.allFlags
}

// GetLongFlags 获取长标志映射
// 返回值:
//   - map[string]*FlagMeta: 长标志名称到标志元数据的映射
func (r *FlagRegistry) GetLongFlags() map[string]*FlagMeta {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.byLong
}

// GetShortFlags 获取短标志映射
// 返回值:
//   - map[string]*FlagMeta: 短标志名称到标志元数据的映射
func (r *FlagRegistry) GetShortFlags() map[string]*FlagMeta {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.byShort
}
