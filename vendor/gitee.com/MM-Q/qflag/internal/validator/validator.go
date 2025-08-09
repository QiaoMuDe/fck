// internal/subcmd/validator.go
// Package validator 内部验证器实现
// 本文件实现了内部使用的验证器功能，提供命令和标志的验证逻辑，
// 包括循环引用检测、命名冲突检查等内部验证机制。
package validator

import (
	"fmt"

	"gitee.com/MM-Q/qflag/internal/types"
)

const (
	// 父子命令依赖检测最大深度
	maxDepth = 10
)

// ValidateSubCommand 验证单个子命令的有效性
//
// 参数：
//   - parent: 当前上下文实例
//   - child: 待添加的上下文实例
//
// 返回值：
//   - error: 验证失败时返回的错误信息, 否则返回nil
func ValidateSubCommand(parent, child *types.CmdContext) error {
	if child == nil {
		return fmt.Errorf("subcmd %s is nil", GetCmdIdentifier(child))
	}

	if parent == nil {
		return nil // 如果父命令为nil，跳过验证
	}

	// 检查父命令的SubCmdMap是否已初始化
	if parent.SubCmdMap == nil {
		panic("父命令的SubCmdMap未初始化")
	}

	// 检测循环引用 - 在获取锁之前进行，避免死锁
	if HasCycleFast(parent, child) {
		return fmt.Errorf("cyclic reference detected: Command %s already exists in the command chain", GetCmdIdentifier(child))
	}

	// 检查名称冲突（外部API层已提供锁保护）
	// 检查长名称冲突
	if child.LongName != "" {
		if _, exists := parent.SubCmdMap[child.LongName]; exists {
			return fmt.Errorf("long name '%s' already exists", child.LongName)
		}
	}

	// 检查短名称冲突
	if child.ShortName != "" {
		if _, exists := parent.SubCmdMap[child.ShortName]; exists {
			return fmt.Errorf("short name '%s' already exists", child.ShortName)
		}
	}

	return nil
}

// HasCycleFast 快速检测父命令和子命令之间是否存在循环依赖
//
// 核心原理：
// 1. 只检查child的父链向上遍历，避免复杂的子树遍历
// 2. 利用CLI工具命令层级浅的特点（通常<10层）
// 3. 时间复杂度从O(n²)优化到O(d)，其中d是命令深度
//
// 参数:
//   - parent: 待添加的父命令上下文
//   - child: 待添加的子命令上下文
//
// 返回值:
//   - bool: true表示存在循环依赖，false表示安全
//
// 使用场景：
//   - 在AddSubCmd函数中调用，防止添加会造成循环依赖的子命令
func HasCycleFast(parent, child *types.CmdContext) bool {
	// 基础安全检查：空指针保护
	// 如果任一参数为空，不可能形成循环依赖
	if parent == nil || child == nil {
		return false // 任一为空，不可能形成循环
	}

	// 直接循环检查：自引用情况
	// 如果parent和child是同一个对象，直接形成循环
	if parent == child {
		return true // 自己指向自己，直接循环
	}

	// 核心算法：向上遍历child的父命令链
	// 从child的直接父命令开始，沿着父链向上查找
	current := child.Parent

	// 循环条件说明：
	// - depth < maxDepth: 限制最大深度，防止异常情况下的无限循环
	// - current != nil: 当到达根命令时自然终止（根命令的Parent为nil）
	for depth := 0; depth < maxDepth && current != nil; depth++ {
		// 循环检测：如果在child的祖先链中找到了parent
		// 说明添加parent->child的边会形成循环
		if current == parent {
			return true // 发现循环依赖
		}

		// 向上移动：继续检查上一级父命令
		// 当current.Parent为nil时（到达根命令），循环自然结束
		current = current.Parent
	}

	// 遍历完成：没有在child的祖先链中找到parent
	// 说明添加parent->child的边不会形成循环
	return false
}

// GetCmdIdentifier 获取命令的标识字符串，用于错误信息
//
// 参数：
//   - cmd: 命令对象
//
// 返回：
//   - 命令标识字符串, 如果为空则返回 <nil>
func GetCmdIdentifier(cmd *types.CmdContext) string {
	if cmd == nil {
		return "<nil>"
	}
	return cmd.GetName()
}
