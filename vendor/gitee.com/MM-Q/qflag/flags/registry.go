// Package flags 标志注册表和元数据管理
// 本文件实现了FlagRegistry标志注册表，提供标志的注册、查找、索引管理等功能，
// 支持按长名称、短名称进行标志查找和管理。
package flags

import (
	"strings"
	"sync"

	"gitee.com/MM-Q/qflag/qerr"
)

// FlagMeta 统一存储标志的完整元数据
type FlagMeta struct {
	Flag Flag // 标志对象
}

// GetLongName 获取标志的长名称
func (m *FlagMeta) GetLongName() string { return m.Flag.LongName() }

// GetShortName 获取标志的短名称
func (m *FlagMeta) GetShortName() string { return m.Flag.ShortName() }

// GetName 获取标志的名称
//
// 优先返回长名称, 如果长名称为空, 则返回短名称
func (m *FlagMeta) GetName() string {
	if m.GetLongName() != "" {
		return m.GetLongName()
	}
	return m.GetShortName()
}

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
	mu           sync.RWMutex         // 并发访问锁（读写锁）
	byLong       map[string]*FlagMeta // 按长名称索引
	byShort      map[string]*FlagMeta // 按短名称索引
	allFlagMetas []*FlagMeta          // 所有标志元数据切片
}

// 创建一个空的标志注册表
//
// 返回值:
//   - *FlagRegistry: 创建的标志注册表指针
func NewFlagRegistry() *FlagRegistry {
	return &FlagRegistry{
		mu:           sync.RWMutex{},
		byLong:       map[string]*FlagMeta{},
		byShort:      map[string]*FlagMeta{},
		allFlagMetas: []*FlagMeta{},
	}
}

// RegisterFlag 注册一个新的标志元数据到注册表中
//
// 参数:
//   - meta: 要注册的标志元数据
//
// 该方法会执行以下操作:
//   - 1.检查长名称和短名称是否已存在
//   - 2.将标志添加到长名称索引
//   - 3.将标志添加到短名称索引
//   - 4.将标志添加到所有标志列表
//
// 返回值:
//   - error: 错误信息, 无错误时为nil
func (r *FlagRegistry) RegisterFlag(meta *FlagMeta) error {
	r.mu.Lock()         // 获取写锁, 保证并发安全
	defer r.mu.Unlock() // 函数返回时释放写锁

	// 获取长标志名称
	longName := meta.GetLongName()
	// 获取短标志名称
	shortName := meta.GetShortName()

	// 检查长短标志是否都为空
	if longName == "" && shortName == "" {
		return qerr.NewValidationError("flag must have at least one name")
	}

	// 如果长标志名称不为空, 则进行检查和添加索引
	if longName != "" {
		// 验证长标志
		if err := r.validateFlagName(meta.GetLongName(), "long", r.byLong); err != nil {
			return err
		}

		// 添加长标志索引
		r.byLong[meta.GetLongName()] = meta
	}

	// 如果短标志名称不为空, 则进行检查和添加索引
	if shortName != "" {
		// 验证短标志
		if err := r.validateFlagName(meta.GetShortName(), "short", r.byShort); err != nil {
			return err
		}

		// 添加短标志索引
		r.byShort[meta.GetShortName()] = meta
	}

	// 添加到标志元数据列表
	r.allFlagMetas = append(r.allFlagMetas, meta)

	return nil
}

// GetByLong 通过长标志名称查找对应的标志元数据
//
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
//
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
//
// 参数:
//   - name可以是长名称(如"help")或短名称(如"h")
//
// 返回值:
//   - *FlagMeta: 找到的标志元数据指针, 未找到时为nil
//   - bool: 是否找到标志, true表示找到
func (r *FlagRegistry) GetByName(name string) (*FlagMeta, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 先尝试按长名称查找
	if meta, exists := r.byLong[name]; exists {
		return meta, exists
	}

	// 再尝试按短名称查找
	if meta, exists := r.byShort[name]; exists {
		return meta, exists
	}

	// 未找到
	return nil, false
}

// GetFlagMetaList 获取所有标志元数据列表
//
// 返回值:
//   - []*FlagMeta: 所有标志元数据的切片
func (r *FlagRegistry) GetFlagMetaList() []*FlagMeta {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.allFlagMetas
}

// GetFlagNameMap 获取所有标志映射(长标志+短标志)
//
// 返回值:
//   - map[string]*FlagMeta: 长短标志名称到标志元数据的映射
func (r *FlagRegistry) GetFlagNameMap() map[string]*FlagMeta {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 创建一个空的标志映射
	allFlags := make(map[string]*FlagMeta, len(r.byLong)+len(r.byShort))

	// 拷贝一份长标志映射
	for k, v := range r.byLong {
		allFlags[k] = v
	}

	// 拷贝一份短标志映射
	for k, v := range r.byShort {
		allFlags[k] = v
	}

	// 返回拷贝后的所有标志映射
	return allFlags
}

// GetLongFlagMap 获取长标志映射
//
// 返回值:
//   - map[string]*FlagMeta: 长标志名称到标志元数据的映射
func (r *FlagRegistry) GetLongFlagMap() map[string]*FlagMeta {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 拷贝一份长标志映射
	byLong := make(map[string]*FlagMeta, len(r.byLong))
	for k, v := range r.byLong {
		byLong[k] = v
	}

	// 返回拷贝后的长标志映射
	return byLong
}

// GetShortFlagMap 获取短标志映射
//
// 返回值:
//   - map[string]*FlagMeta: 短标志名称到标志元数据的映射
func (r *FlagRegistry) GetShortFlagMap() map[string]*FlagMeta {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 拷贝一份短标志映射
	byShort := make(map[string]*FlagMeta, len(r.byShort))
	for k, v := range r.byShort {
		byShort[k] = v
	}

	// 返回拷贝后的短标志映射
	return byShort
}

// GetLongFlagsCount 获取长标志数量
//
// 返回值:
//   - int: 长标志的数量
func (r *FlagRegistry) GetLongFlagsCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.byLong)
}

// GetShortFlagsCount 获取短标志数量
//
// 返回值:
//   - int: 短标志的数量
func (r *FlagRegistry) GetShortFlagsCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.byShort)
}

// GetAllFlagsCount 获取所有标志数量(长标志+短标志)
//
// 返回值:
//   - int: 所有标志的数量
func (r *FlagRegistry) GetAllFlagsCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.byLong) + len(r.byShort)
}

// GetFlagMetaCount 获取标志元数据数量
//
// 返回值:
//   - int: 标志元数据的数量
func (r *FlagRegistry) GetFlagMetaCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.allFlagMetas)
}

// validateFlagName 验证标志名称的通用函数
//
// 参数:
//   - name: 标志名称
//   - nameType: 标志类型(如"long"或"short")
//   - existingMap: 已存在的标志映射
//
// 返回值:
//   - error: 验证错误, 验证通过返回nil
func (r *FlagRegistry) validateFlagName(name, nameType string, existingMap map[string]*FlagMeta) error {
	if name == "" {
		return nil // 空名称直接返回，由调用方处理
	}

	// 检查名称是否包含非法字符
	if strings.ContainsAny(name, InvalidFlagChars) {
		return qerr.NewValidationErrorf("%s flag name '%s' contains illegal characters", nameType, name)
	}

	// 检查标志是否已存在
	if _, exists := existingMap[name]; exists {
		return qerr.NewValidationErrorf("%s flag '%s' already exists", nameType, name)
	}

	return nil
}
