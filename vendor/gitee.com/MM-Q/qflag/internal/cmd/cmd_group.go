// Package cmd 提供命令实现和命令管理功能
//
// cmd_group.go 包含互斥组、必需组和标志依赖相关的功能实现
//
// 本文件提供了以下主要功能:
//   - 互斥组管理: 限制组内标志最多只能设置一个
//   - 必需组管理: 要求组内标志全部设置或条件性必需
//   - 标志依赖管理: 定义标志之间的依赖关系
//
// 主要方法列表:
//   - AddMutexGroup: 添加互斥组
//   - MutexGroups: 获取所有互斥组
//   - AddRequiredGroup: 添加必需组
//   - RequiredGroups: 获取所有必需组
//   - AddFlagDependency: 添加标志依赖关系
//   - FlagDependencies: 获取所有标志依赖关系
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
// 标志依赖特性:
//   - 支持互斥依赖: 触发标志设置时, 目标标志不能设置
//   - 支持必需依赖: 触发标志设置时, 目标标志必须设置
//   - 适用于条件性约束
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

// RequiredGroups 获取所有必需组
//
// 返回值:
//   - []types.RequiredGroup: 必需组列表的副本
//
// 功能说明:
//   - 返回命令中定义的所有必需组
//   - 返回副本以防止外部修改内部状态
//   - 使用读锁保护并发安全
func (c *Cmd) RequiredGroups() []types.RequiredGroup {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.config.RequiredGroups) == 0 {
		return []types.RequiredGroup{}
	}
	groups := make([]types.RequiredGroup, len(c.config.RequiredGroups))
	copy(groups, c.config.RequiredGroups)
	return groups
}

// AddFlagDependency 添加标志依赖关系
//
// 参数:
//   - name: 依赖关系名称
//   - trigger: 触发标志名称
//   - targets: 目标标志名称列表
//   - depType: 依赖关系类型
//
// 返回值:
//   - error: 添加失败时返回错误
func (c *Cmd) AddFlagDependency(name, trigger string, targets []string, depType types.DepType) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 验证名称
	if name == "" {
		return fmt.Errorf("empty flag dependency name in '%s'", c.Name())
	}

	// 检查是否已存在
	for _, dep := range c.config.FlagDependencies {
		if dep.Name == name {
			return fmt.Errorf("duplicate flag dependency '%s' in '%s'", name, c.Name())
		}
	}

	// 验证触发标志
	if trigger == "" {
		return fmt.Errorf("empty trigger flag in '%s'", c.Name())
	}

	// 验证目标标志列表
	if len(targets) == 0 {
		return fmt.Errorf("empty target flags in '%s'", c.Name())
	}

	// 检查自依赖
	for _, target := range targets {
		if target == trigger {
			return fmt.Errorf("trigger flag '%s' cannot be in targets in '%s'", trigger, c.Name())
		}
	}

	// 验证触发标志是否存在
	if _, exists := c.flagRegistry.Get(trigger); !exists {
		return fmt.Errorf("trigger flag '%s' not found in '%s'", trigger, c.Name())
	}

	// 验证目标标志是否存在
	for _, target := range targets {
		if _, exists := c.flagRegistry.Get(target); !exists {
			return fmt.Errorf("target flag '%s' not found in '%s'", target, c.Name())
		}
	}

	// 创建依赖关系
	dep := types.FlagDependency{
		Name:    name,
		Trigger: trigger,
		Targets: targets,
		Type:    depType,
	}

	// 添加到配置
	c.config.FlagDependencies = append(c.config.FlagDependencies, dep)
	return nil
}

// FlagDependencies 获取所有标志依赖关系
//
// 返回值:
//   - []types.FlagDependency: 标志依赖关系列表的副本
//
// 功能说明:
//   - 返回命令中定义的所有标志依赖关系
//   - 返回副本以防止外部修改内部状态
//   - 使用读锁保护并发安全
func (c *Cmd) FlagDependencies() []types.FlagDependency {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.config.FlagDependencies) == 0 {
		return []types.FlagDependency{}
	}
	deps := make([]types.FlagDependency, len(c.config.FlagDependencies))
	copy(deps, c.config.FlagDependencies)
	return deps
}
