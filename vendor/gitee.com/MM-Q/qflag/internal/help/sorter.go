// Package help 帮助信息排序和组织
// 本文件实现了帮助信息的排序和组织功能，包括标志排序、子命令排序等，
// 确保帮助信息以合理的顺序展示给用户。
package help

import (
	"sort"

	"gitee.com/MM-Q/qflag/internal/types"
)

// NamedItem 表示具有长名称和短名称的项目接口
type NamedItem interface {
	GetLongName() string
	GetShortName() string
}

// flagInfoItem 为 flagInfo 实现 NamedItem 接口
type flagInfoItem struct {
	flagInfo
}

func (f flagInfoItem) GetLongName() string {
	return f.longFlag
}

func (f flagInfoItem) GetShortName() string {
	return f.shortFlag
}

// subCmdItem 为子命令实现 NamedItem 接口
type subCmdItem struct {
	*types.CmdContext
}

func (s subCmdItem) GetLongName() string {
	return s.LongName
}

func (s subCmdItem) GetShortName() string {
	return s.ShortName
}

// sortByNamePriority 通用排序函数，按短名称优先级排序
//
// 排序优先级: 1.有短名称的优先 2.按长名称字母序 3.短名称字母序
//
// 参数：
//   - items: 实现了 NamedItem 接口的项目切片
func sortByNamePriority(items []NamedItem) {
	sort.Slice(items, func(i, j int) bool {
		a, b := items[i], items[j]
		return sortWithShortNamePriority(
			a.GetShortName() != "",
			b.GetShortName() != "",
			a.GetLongName(),
			b.GetLongName(),
			a.GetShortName(),
			b.GetShortName(),
		)
	})
}

// sortWithShortNamePriority 通用排序比较函数
//
// 排序优先级: 1.有短名称的优先 2.按长名称字母序 3.短名称字母序
//
// 参数：
//   - aHasShort: a是否有短名称
//   - bHasShort: b是否有短名称
//   - aName: a的长名称
//   - bName: b的长名称
//   - aShort: a的短名称
//   - bShort: b的短名称
//
// 返回：
//   - bool: a是否应该排在b之前
func sortWithShortNamePriority(aHasShort, bHasShort bool, aName, bName, aShort, bShort string) bool {
	// 1. 有短名称的优先
	if aHasShort != bHasShort {
		return aHasShort
	}

	// 2. 按长名称字母顺序排序
	if aName != bName {
		return aName < bName
	}

	// 3. 都有短名称则按短名称字母顺序排序
	return aShort < bShort
}

// sortSubCommands 对子命令进行排序
//
// 参数：
//   - subCmds - 需要排序的子命令列表
func sortSubCommands(subCmds []*types.CmdContext) {
	// 将子命令转换为 NamedItem 接口并排序
	items := make([]NamedItem, len(subCmds))
	for i, subCmd := range subCmds {
		items[i] = subCmdItem{subCmd}
	}

	// 使用通用排序函数
	sortByNamePriority(items)

	// 将排序结果写回原切片
	for i, item := range items {
		subCmds[i] = item.(subCmdItem).CmdContext
	}
}

// sortFlags 按短标志字母顺序排序，有短标志的选项优先
//
// 参数：
//   - flags - 需要排序的标志列表
func sortFlags(flags []flagInfo) {
	// 将 flagInfo 转换为 NamedItem 接口
	items := make([]NamedItem, len(flags))
	for i, flag := range flags {
		items[i] = flagInfoItem{flag}
	}

	// 使用通用排序函数
	sortByNamePriority(items)

	// 将排序结果写回原切片
	for i, item := range items {
		flags[i] = item.(flagInfoItem).flagInfo
	}
}
