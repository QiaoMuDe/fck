package parser

import (
	"fmt"

	"gitee.com/MM-Q/qflag/internal/types"
)

// buildFlagDisplayName 构建标志的显示名称
//
// 参数:
//   - flag: 标志对象
//
// 返回值:
//   - string: 显示名称
//   - 有长短标志："--long/-s"
//   - 只有长标志："--long"
//   - 只有短标志："-s"
//
// 功能说明:
//   - 根据标志的长名和短名生成友好的显示名称
//   - 优先显示长名，短名作为补充
//   - 用于错误信息中显示用户使用的标志
func (p *DefaultParser) buildFlagDisplayName(flag types.Flag) string {
	longName := flag.LongName()
	shortName := flag.ShortName()

	// 长标志和短标志都存在
	if longName != "" && shortName != "" {
		return fmt.Sprintf("--%s/-%s", longName, shortName)
	}

	// 只有长标志
	if longName != "" {
		return fmt.Sprintf("--%s", longName)
	}

	// 只有短标志
	if shortName != "" {
		return fmt.Sprintf("-%s", shortName)
	}

	// 都没有（理论上不应该发生），返回 Name()
	return flag.Name()
}

// buildSetFlagsMap 构建已设置标志的映射和显示名称映射
//
// 参数:
//   - cmd: 要验证的命令
//
// 返回值:
//   - map[string]bool: 已设置标志名称的映射，key为标志名，value为true
//
// 功能说明:
//   - 遍历所有标志，收集已设置的标志
//   - 同时构建所有标志的显示名称映射
//   - 支持混合使用长短名：将长名和短名都存储到 map 中
//   - 将结果缓存到解析器的 setFlagsMap 和 flagDisplayNames 字段中
//   - 根据标志数量预分配map空间，减少扩容开销
func (p *DefaultParser) buildSetFlagsMap(cmd types.Command) map[string]bool {
	// 根据标志数量预分配map空间，减少扩容开销
	flags := cmd.Flags()
	p.setFlagsMap = make(map[string]bool, len(flags)*2)        // 预留空间存储长名和短名
	p.flagDisplayNames = make(map[string]string, len(flags)*2) // 预分配空间存储显示名称

	for _, flag := range flags {
		// 生成显示名称
		displayName := p.buildFlagDisplayName(flag)

		// 存储长名和缓存显示名称
		if flag.LongName() != "" {
			p.flagDisplayNames[flag.LongName()] = displayName
			if flag.IsSet() {
				p.setFlagsMap[flag.LongName()] = true
			}
		}

		// 存储短名和缓存显示名称
		if flag.ShortName() != "" {
			p.flagDisplayNames[flag.ShortName()] = displayName
			if flag.IsSet() {
				p.setFlagsMap[flag.ShortName()] = true
			}
		}
	}

	return p.setFlagsMap
}

// validateMutexGroups 验证命令的互斥组规则
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - error: 如果互斥组验证失败返回错误
//
// 功能说明:
//   - 检查每个互斥组中是否有多个标志被设置
//   - 检查不允许为空的互斥组中是否有至少一个标志被设置
//   - 提供清晰的错误信息, 指出冲突的标志和组名
//
// 验证规则:
//   - 互斥组中最多只能有一个标志被设置
//   - 如果 AllowNone 为 false, 则必须至少有一个标志被设置
//
// 错误处理:
//   - 使用 fmt.Errorf 创建错误
//   - 错误信息包含互斥组名称和冲突的标志列表
//
// 性能优化:
//   - 使用缓存的已设置标志映射，避免重复的 GetFlag() 和 IsSet() 调用
func (p *DefaultParser) validateMutexGroups(config *types.CmdConfig) error {
	// 检查互斥组是否为空
	if len(config.MutexGroups) == 0 {
		return nil
	}

	// 使用缓存的已设置标志映射
	setFlags := p.setFlagsMap

	// 遍历所有互斥组
	for _, group := range config.MutexGroups {
		var setFlagsList []string
		seenDisplayNames := make(map[string]bool, len(group.Flags)) // 去重map，防止重复显示相同的标志

		// 检查互斥组中的每个标志是否被设置
		for _, flagName := range group.Flags {
			// 如果拿组里的标志没有获取到显示名称, 则表示为不是一个有效标志, 返回错误
			displayName, ok := p.flagDisplayNames[flagName]
			if !ok {
				return fmt.Errorf("invalid flag '%s' in mutex group '%s'", flagName, group.Name)
			}

			// 如果标志被设置, 添加到已设置列表
			if setFlags[flagName] {
				if !seenDisplayNames[displayName] {
					// 添加去重检查，避免同一个标志的多个名称重复显示
					seenDisplayNames[displayName] = true
					setFlagsList = append(setFlagsList, displayName)
				}
			}
		}

		// 直接计算实际设置的标志数量（用于验证逻辑）
		setCount := len(setFlagsList)

		// 如果互斥组中设置了多个标志, 返回错误
		if setCount > 1 {
			return fmt.Errorf("mutually exclusive flags %v in group '%s' cannot be used together", setFlagsList, group.Name)
		}

		// 如果不允许为空, 且互斥组中没有设置任何标志, 返回错误
		if !group.AllowNone && setCount == 0 {
			return fmt.Errorf("one of flags %v in mutex group '%s' must be set", group.Flags, group.Name)
		}
	}

	return nil
}

// validateRequiredGroups 验证命令的必需组规则
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - error: 如果必需组验证失败返回错误
//
// 功能说明:
//   - 检查每个必需组中是否有标志未被设置
//   - 对于条件性必需组，只有当组中任何一个标志被设置时，才检查所有标志
//   - 提供清晰的错误信息，指出未设置的标志和组名
//
// 验证规则:
//   - 普通必需组：所有标志都必须被设置
//   - 条件性必需组：如果任何一个标志被设置，则所有标志都必须被设置
//
// 错误处理:
//   - 使用 fmt.Errorf 创建错误
//   - 错误信息包含必需组名称和未设置的标志列表
//   - 条件性必需组的错误信息会明确指出是因为使用了某个标志而要求其他标志
//
// 性能优化:
//   - 使用缓存的已设置标志映射，避免重复的 GetFlag() 和 IsSet() 调用
func (p *DefaultParser) validateRequiredGroups(config *types.CmdConfig) error {
	if len(config.RequiredGroups) == 0 {
		return nil
	}

	// 使用缓存的已设置标志映射
	setFlags := p.setFlagsMap

	// 遍历所有必需组
	for _, group := range config.RequiredGroups {
		var unsetFlags []string
		var setFlagsCount int

		// 检查组中是否有任何标志被设置
		for _, flagName := range group.Flags {
			if setFlags[flagName] {
				setFlagsCount++
			}
		}

		// 如果是条件性必需组，且没有任何标志被设置，则跳过验证
		if group.Conditional && setFlagsCount == 0 {
			continue
		}

		seenUnsetDisplayNames := make(map[string]bool, len(group.Flags)) // 去重map，防止重复显示相同的标志

		// 遍历组中的每个标志
		for _, flagName := range group.Flags {
			// 如果拿组里的标志没有获取到显示名称, 则表示为不是一个有效标志, 返回错误
			displayName, ok := p.flagDisplayNames[flagName]
			if !ok {
				return fmt.Errorf("invalid flag '%s' in required group '%s'", flagName, group.Name)
			}

			// 如果标志未被设置, 添加到未设置列表
			if !setFlags[flagName] {
				if !seenUnsetDisplayNames[displayName] {
					// 添加去重检查，避免同一个标志的多个名称重复显示
					seenUnsetDisplayNames[displayName] = true
					unsetFlags = append(unsetFlags, displayName)
				}
			}
		}

		// 如果组中有未设置的标志, 返回错误
		if len(unsetFlags) > 0 {
			if group.Conditional {
				return fmt.Errorf("flags %v in group '%s' must all be set", unsetFlags, group.Name)
			}
			return fmt.Errorf("required flags %v in group '%s' must be set", unsetFlags, group.Name)
		}
	}

	return nil
}
