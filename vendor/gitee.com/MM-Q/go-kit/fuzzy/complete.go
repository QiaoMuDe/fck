/*
Package fuzzy 提供了命令行补全搜索功能。

Complete 函数优先匹配前缀, 其次进行模糊匹配, 适用于命令行标志补全场景。

基本用法:

	flags := []string{"--verbose", "--version", "-v", "-h"}

	// 前缀优先的补全
	matches := fuzzy.Complete("--v", flags)
	// 结果: [--verbose, --version, -v]

	// 仅前缀匹配
	matches := fuzzy.CompletePrefix("--v", flags)
	// 结果: [--verbose, --version]

	// 仅精确匹配
	matches := fuzzy.CompleteExact("-v", flags)
	// 结果: [-v]
*/
package fuzzy

import (
	"sort"
	"strings"
)

// Complete 在 candidates 中搜索 pattern, 优先前缀匹配, 其次模糊匹配
// 返回按匹配质量排序的结果, 分数越高表示匹配度越好
//
// 匹配优先级:
//  1. 精确匹配 (1000分)
//  2. 前缀匹配 (100+分, 越短的候选分数越高)
//  3. 模糊匹配 (0-99分, 复用fuzzy算法但降低权重)
//
// 参数:
//   - pattern: 要查找的模式字符串
//   - candidates: 候选字符串切片
//
// 返回值:
//   - Matches: 按匹配质量降序排列的匹配结果
//
// 示例:
//
//	flags := []string{"--verbose", "--version", "-v"}
//	matches := fuzzy.Complete("--v", flags)
//	// matches[0].Str = "--verbose"
//	// matches[1].Str = "--version"
//	// matches[2].Str = "-v"
func Complete(pattern string, candidates []string) Matches {
	if len(pattern) == 0 || len(candidates) == 0 {
		return nil
	}

	var matches Matches

	// 遍历所有候选, 为每个打分并记录匹配位置
	for i, candidate := range candidates {
		score, matchedIndexes := matchComplete(pattern, candidate)
		if score > 0 {
			matches = append(matches, Match{
				Str:            candidate,
				Index:          i,
				Score:          score,
				MatchedIndexes: matchedIndexes,
			})
		}
	}

	// 按分数降序排序
	sort.Stable(matches)
	return matches
}

// CompletePrefix 仅返回以 pattern 为前缀的候选
// 不区分大小写, 返回结果按候选长度升序排列 (短的优先)
//
// 参数:
//   - pattern: 前缀模式字符串
//   - candidates: 候选字符串切片
//
// 返回值:
//   - Matches: 前缀匹配的结果列表
//
// 示例:
//
//	flags := []string{"--verbose", "--version", "-v"}
//	matches := fuzzy.CompletePrefix("--v", flags)
//	// 结果: [--verbose, --version] (不包含 -v)
func CompletePrefix(pattern string, candidates []string) Matches {
	if len(pattern) == 0 || len(candidates) == 0 {
		return nil
	}

	var matches Matches
	patternLower := strings.ToLower(pattern)

	for i, candidate := range candidates {
		if strings.HasPrefix(strings.ToLower(candidate), patternLower) {
			// 越短的候选分数越高
			score := 1000 - len(candidate)
			if score < 0 {
				score = 0
			}

			// 记录匹配位置
			matchedIndexes := make([]int, 0, len(pattern))
			for j := 0; j < len(pattern) && j < len(candidate); j++ {
				matchedIndexes = append(matchedIndexes, j)
			}

			matches = append(matches, Match{
				Str:            candidate,
				Index:          i,
				Score:          score,
				MatchedIndexes: matchedIndexes,
			})
		}
	}

	// 按分数降序排序 (短的在前)
	sort.Stable(matches)
	return matches
}

// CompleteExact 仅返回与 pattern 完全相等的候选
// 区分大小写
//
// 参数:
//   - pattern: 精确匹配字符串
//   - candidates: 候选字符串切片
//
// 返回值:
//   - Matches: 精确匹配的结果列表 (最多一个)
//
// 示例:
//
//	flags := []string{"--verbose", "--version", "-v"}
//	matches := fuzzy.CompleteExact("-v", flags)
//	// 结果: [-v]
func CompleteExact(pattern string, candidates []string) Matches {
	if len(pattern) == 0 || len(candidates) == 0 {
		return nil
	}

	for i, candidate := range candidates {
		if pattern == candidate {
			return Matches{{
				Str:            candidate,
				Index:          i,
				Score:          10000,
				MatchedIndexes: getAllIndexes(candidate),
			}}
		}
	}

	return nil
}

// matchComplete 判断匹配类型并返回分数和匹配位置
// 返回 (0, nil) 表示不匹配
func matchComplete(pattern, candidate string) (int, []int) {
	// 1. 精确匹配 (区分大小写)
	if pattern == candidate {
		return 1000, getAllIndexes(candidate)
	}

	// 2. 前缀匹配 (不区分大小写)
	patternLower := strings.ToLower(pattern)
	candidateLower := strings.ToLower(candidate)

	if strings.HasPrefix(candidateLower, patternLower) {
		// 越短的候选分数越高 (优先短候选)
		// 基础分 200, 减去长度惩罚, 确保分数在 100-200 之间
		score := 200 - len(candidate)
		if score < 100 {
			score = 100
		}

		// 记录前缀匹配位置
		matchedIndexes := make([]int, 0, len(pattern))
		for i := 0; i < len(pattern) && i < len(candidate); i++ {
			matchedIndexes = append(matchedIndexes, i)
		}

		return score, matchedIndexes
	}

	// 3. 模糊匹配 (复用现有算法)
	matches := FindFromNoSort(pattern, stringSource([]string{candidate}))
	if len(matches) > 0 && matches[0].Score > 0 {
		// 降低模糊匹配权重, 确保前缀匹配优先
		score := matches[0].Score / 10
		if score > 99 {
			score = 99
		}
		return score, matches[0].MatchedIndexes
	}

	return 0, nil
}

// getAllIndexes 获取字符串的所有字节索引
func getAllIndexes(str string) []int {
	indexes := make([]int, 0, len(str))
	for i := range str {
		indexes = append(indexes, i)
	}
	return indexes
}
