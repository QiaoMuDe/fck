// Package cmd 提供命令实现和命令管理功能
//
// cmd_group.go 包含互斥组和必需组相关的功能实现
//
// 本文件提供了以下主要功能:
//   - 互斥组管理: 限制组内标志最多只能设置一个
//   - 必需组管理: 要求组内标志全部设置或条件性必需
//
// 主要方法列表:
//   - AddMutexGroup: 添加互斥组
//   - GetMutexGroup/GetMutexGroups: 获取互斥组
//   - RemoveMutexGroup: 移除互斥组
//   - MutexGroups: 获取所有互斥组
//   - AddRequiredGroup: 添加必需组
//   - GetRequiredGroup/RequiredGroups: 获取必需组
//   - RemoveRequiredGroup: 移除必需组
//
// 互斥组特性:
//   - 组内标志互斥, 只能设置其中一个
//   - 支持 allowNone 参数控制是否允许都不设置
//   - 适用于排他性的选项组合
//
// 必需组特性:
//   - 支持条件性必需 (Conditional)
//   - 当 conditional=true 时, 只有组中任一标志被设置才要求全部设置
//   - 适用于依赖性选项
//
// 线程安全:
//   - 所有公共方法都使用读写锁保护
//   - 支持并发安全的访问和修改
package cmd

import (
	"fmt"

	"gitee.com/MM-Q/qflag/internal/types"
)

// AddMutexGroup 添加互斥组到命令
//
// 参数:
//   - name: 互斥组名称, 用于错误提示和标识
//   - flags: 互斥组中的标志名称列表
//   - allowNone: 是否允许一个都不设置
//
// 返回值:
//   - error: 添加失败时返回错误
//
// 功能说明:
//   - 创建新的互斥组并添加到命令配置中
//   - 互斥组中的标志最多只能有一个被设置
//   - 如果 allowNone 为 false, 则必须至少有一个标志被设置
//   - 使用写锁保护并发安全
//
// 注意事项:
//   - 标志名称必须是已注册的标志
//   - 互斥组名称在命令中应该唯一
//   - 如果组名已存在，返回错误
func (c *Cmd) AddMutexGroup(name string, flags []string, allowNone bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 互斥组名称不能为空
	if name == "" {
		return fmt.Errorf("empty mutex group name in '%s'", c.Name())
	}

	// 检查互斥组名称是否已存在
	for _, group := range c.config.MutexGroups {
		if group.Name == name {
			return fmt.Errorf("duplicate mutex group '%s' in '%s'", name, c.Name())
		}
	}

	// 检查标志是否为空
	if len(flags) == 0 {
		return fmt.Errorf("empty mutex group '%s' in '%s'", name, c.Name())
	}

	// 检查标志名称是否存在
	for _, flagName := range flags {
		if _, exists := c.flagRegistry.Get(flagName); !exists {
			return fmt.Errorf("flag '%s' not found in '%s'", flagName, c.Name())
		}
	}

	group := types.MutexGroup{
		Name:      name,
		Flags:     flags,
		AllowNone: allowNone,
	}

	c.config.MutexGroups = append(c.config.MutexGroups, group)
	return nil
}

// GetMutexGroups 获取命令的所有互斥组
//
// 返回值:
//   - []types.MutexGroup: 互斥组列表的副本
//
// 功能说明:
//   - 返回命令中定义的所有互斥组
//   - 返回副本以防止外部修改内部状态
//   - 使用读锁保护并发安全
func (c *Cmd) GetMutexGroups() []types.MutexGroup {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// 返回副本以防止外部修改
	groups := make([]types.MutexGroup, len(c.config.MutexGroups))
	copy(groups, c.config.MutexGroups)
	return groups
}

// RemoveMutexGroup 移除指定名称的互斥组
//
// 参数:
//   - name: 要移除的互斥组名称
//
// 返回值:
//   - error: 移除失败时返回错误
//
// 功能说明:
//   - 根据名称查找并移除互斥组
//   - 使用写锁保护并发安全
//   - 如果找不到对应名称的互斥组, 返回错误
func (c *Cmd) RemoveMutexGroup(name string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, group := range c.config.MutexGroups {
		if group.Name == name {
			c.config.MutexGroups = append(c.config.MutexGroups[:i], c.config.MutexGroups[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("mutex group '%s' not found in '%s'", name, c.Name())
}

// GetMutexGroup 获取指定名称的互斥组
//
// 参数:
//   - name: 要获取的互斥组名称
//
// 返回值:
//   - *types.MutexGroup: 互斥组指针, 如果找到则返回对应的互斥组
//   - bool: 是否找到, true表示找到
//
// 功能说明:
//   - 根据名称查找互斥组
//   - 使用读锁保护并发安全
//   - 如果找不到对应名称的互斥组, 返回nil和false
func (c *Cmd) GetMutexGroup(name string) (*types.MutexGroup, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for i := range c.config.MutexGroups {
		if c.config.MutexGroups[i].Name == name {
			return &c.config.MutexGroups[i], true
		}
	}
	return nil, false
}

// MutexGroups 获取所有互斥组
//
// 返回值:
//   - []types.MutexGroup: 互斥组列表的副本
//
// 功能说明:
//   - 返回命令中定义的所有互斥组
//   - 返回副本以防止外部修改内部状态
//   - 使用读锁保护并发安全
func (c *Cmd) MutexGroups() []types.MutexGroup {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.config.MutexGroups) == 0 {
		return []types.MutexGroup{}
	}
	groups := make([]types.MutexGroup, len(c.config.MutexGroups))
	copy(groups, c.config.MutexGroups)
	return groups
}

// AddRequiredGroup 添加必需组
//
// 参数:
//   - name: 必需组名称
//   - flags: 必需组中的标志名称列表
//   - conditional: 是否为条件性必需组，如果为true，则只有当组中任何一个标志被设置时，才要求所有标志都被设置
//
// 返回值:
//   - error: 添加失败时返回错误
//
// 功能说明:
//   - 添加一个必需组到命令配置
//   - 如果组名已存在，返回错误
//   - 如果标志列表为空，返回错误
//   - 如果标志不存在，返回错误
//   - 支持条件性必需组，当conditional为true时，只有当组中任何一个标志被设置时，才要求所有标志都被设置
func (c *Cmd) AddRequiredGroup(name string, flags []string, conditional bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 必需组名称不能为空
	if name == "" {
		return fmt.Errorf("empty required group name in '%s'", c.Name())
	}

	// 检查必需组名称是否已存在
	for _, group := range c.config.RequiredGroups {
		if group.Name == name {
			return fmt.Errorf("duplicate required group '%s' in '%s'", name, c.Name())
		}
	}

	// 必需组标志列表不能为空
	if len(flags) == 0 {
		return fmt.Errorf("empty required group '%s' in '%s'", name, c.Name())
	}

	// 检查必需组标志是否存在
	for _, flagName := range flags {
		if _, exists := c.flagRegistry.Get(flagName); !exists {
			return fmt.Errorf("flag '%s' not found in '%s'", flagName, c.Name())
		}
	}

	// 添加必需组
	group := types.RequiredGroup{
		Name:        name,
		Flags:       flags,
		Conditional: conditional,
	}

	c.config.RequiredGroups = append(c.config.RequiredGroups, group)
	return nil
}

// RemoveRequiredGroup 移除必需组
//
// 参数:
//   - name: 必需组名称
//
// 返回值:
//   - error: 移除失败时返回错误
//
// 功能说明:
//   - 从命令配置中移除指定的必需组
//   - 如果组不存在，返回错误
func (c *Cmd) RemoveRequiredGroup(name string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, group := range c.config.RequiredGroups {
		if group.Name == name {
			c.config.RequiredGroups = append(c.config.RequiredGroups[:i], c.config.RequiredGroups[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("required group '%s' not found in '%s'", name, c.Name())
}

// GetRequiredGroup 获取必需组
//
// 参数:
//   - name: 必需组名称
//
// 返回值:
//   - *types.RequiredGroup: 必需组指针
//   - bool: 是否找到
//
// 功能说明:
//   - 根据名称获取必需组
//   - 如果组不存在，返回 nil 和 false
func (c *Cmd) GetRequiredGroup(name string) (*types.RequiredGroup, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for i := range c.config.RequiredGroups {
		if c.config.RequiredGroups[i].Name == name {
			return &c.config.RequiredGroups[i], true
		}
	}

	return nil, false
}

// RequiredGroups 获取所有必需组
//
// 返回值:
//   - []types.RequiredGroup: 所有必需组列表
//
// 功能说明:
//   - 返回命令配置中的所有必需组
//   - 返回的是副本，修改不会影响原配置
func (c *Cmd) RequiredGroups() []types.RequiredGroup {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.config.RequiredGroups) == 0 {
		return []types.RequiredGroup{}
	}
	result := make([]types.RequiredGroup, len(c.config.RequiredGroups))
	copy(result, c.config.RequiredGroups)
	return result
}
