/*
Package fuzzy 提供了模糊字符串匹配功能,
该功能针对文件名和代码符号进行了优化,
风格类似于 Sublime Text、VSCode、IntelliJ IDEA 等工具。
*/
package fuzzy

import (
	"sort"
	"unicode"
	"unicode/utf8"
)

// Find 在 data 中查找 pattern 并返回匹配结果,按匹配质量降序排列。
// 匹配质量由一组奖励和惩罚规则决定。
//
// 以下类型的匹配会获得奖励:
//   - 模式中的第一个字符与匹配字符串的第一个字符匹配
//   - 匹配的字符是驼峰命名
//   - 匹配的字符位于分隔符 (如下划线字符) 之后
//   - 匹配的字符与之前的匹配相邻
//
// 对于搜索字符串中每个未匹配的字符以及第一个匹配之前的所有前导字符,都会应用惩罚。
//
// 参数:
//   - pattern: 要查找的模式字符串
//   - data: 要搜索的字符串切片
//
// 返回值:
//   - Matches: 按匹配质量降序排列的匹配结果
func Find(pattern string, data []string) Matches {
	return FindFrom(pattern, stringSource(data))
}

// FindNoSort 在 data 中查找 pattern 并返回匹配结果,不排序。
// 该函数与 Find 功能相同,但返回的结果不会按匹配质量排序,性能略好。
//
// 参数:
//   - pattern: 要查找的模式字符串
//   - data: 要搜索的字符串切片
//
// 返回值:
//   - Matches: 匹配结果 (未排序)
func FindNoSort(pattern string, data []string) Matches {
	return FindFromNoSort(pattern, stringSource(data))
}

// FindFrom 使用 Source 接口在数据源中查找 pattern 并返回排序后的匹配结果。
// 该函数与 Find 功能相同,但接受 Source 接口而不是字符串切片,适用于自定义数据源。
// 结果按匹配质量降序排列。
//
// 参数:
//   - pattern: 要查找的模式字符串
//   - data: 实现 Source 接口的数据源
//
// 返回值:
//   - Matches: 按匹配质量降序排列的匹配结果
func FindFrom(pattern string, data Source) Matches {
	matches := FindFromNoSort(pattern, data)
	sort.Stable(matches)
	return matches
}

// FindFromNoSort 使用 Source 接口在数据源中查找 pattern,不排序。
// 该函数是模糊匹配的核心实现,负责计算匹配分数和匹配位置。
// 匹配质量由奖励和惩罚规则决定: 首字符匹配、驼峰命名匹配、分隔符后匹配、相邻匹配可获得奖励；
// 未匹配字符和前导字符会受到惩罚。
//
// 参数:
//   - pattern: 要查找的模式字符串
//   - data: 实现 Source 接口的数据源
//
// 返回值:
//   - Matches: 匹配结果 (未排序)
func FindFromNoSort(pattern string, data Source) Matches {
	// 空模式直接返回 nil,无需匹配
	if len(pattern) == 0 {
		return nil
	}

	// 空数据源直接返回 nil,无需匹配
	if data.Len() == 0 {
		return nil
	}

	// 将模式字符串转换为 rune 切片,以正确处理 Unicode 字符
	runes := []rune(pattern)
	var matches Matches
	var matchedIndexes []int

	// 遍历数据源中的每个字符串进行匹配
	for i := 0; i < data.Len(); i++ {
		// 创建当前字符串的匹配结果对象
		var match Match
		// 缓存字符串,避免多次调用 data.String(i) 带来的性能开销
		str := data.String(i)
		match.Str = str
		match.Index = i

		// 复用 matchedIndexes 切片以减少内存分配
		if matchedIndexes != nil {
			match.MatchedIndexes = matchedIndexes
		} else {
			match.MatchedIndexes = make([]int, 0, len(runes))
		}

		// 初始化匹配状态变量
		var score int               // 当前字符匹配的分数
		patternIndex := 0           // 当前正在匹配的模式字符索引
		bestScore := -1             // 当前模式字符的最佳匹配分数
		matchedIndex := -1          // 当前模式字符的最佳匹配位置
		currAdjacentMatchBonus := 0 // 相邻匹配奖励累计值
		var last rune               // 上一个遍历的字符
		var lastIndex int           // 上一个遍历的字符位置

		// 预读取字符串的第一个字符
		nextc, nextSize := utf8.DecodeRuneInString(str)
		var candidate rune
		var candidateSize int

		// 遍历字符串中的每个字符,寻找模式匹配
		for j := 0; j < len(str); j += candidateSize {
			candidate, candidateSize = nextc, nextSize

			// 检查当前字符是否匹配当前模式字符 (不区分大小写)
			if equalFold(candidate, runes[patternIndex]) {
				score = 0

				// 奖励1: 首字符匹配 (模式第一个字符匹配字符串开头)
				if j == 0 {
					score += firstCharMatchBonus
				}

				// 奖励2: 驼峰命名匹配 (小写后接大写)
				if unicode.IsLower(last) && unicode.IsUpper(candidate) {
					score += camelCaseMatchBonus
				}

				// 奖励3: 分隔符后匹配
				if j != 0 && isSeparator(last) {
					score += matchFollowingSeparatorBonus
				}

				// 奖励4: 相邻字符匹配 (连续匹配奖励递增)
				if len(match.MatchedIndexes) > 0 {
					lastMatch := match.MatchedIndexes[len(match.MatchedIndexes)-1]
					bonus := adjacentCharBonus(lastIndex, lastMatch, currAdjacentMatchBonus)
					score += bonus
					// 相邻匹配是递增的,并基于之前的相邻匹配不断增加
					// 因此我们需要保持当前的匹配奖励
					currAdjacentMatchBonus += bonus
				}

				// 记录当前模式字符的最佳匹配位置
				if score > bestScore {
					bestScore = score
					matchedIndex = j
				}
			}

			// 预读取下一个模式字符和字符串字符,用于决策
			var nextp rune
			if patternIndex < len(runes)-1 {
				nextp = runes[patternIndex+1]
			}

			// 读取字符串中的下一个字符 (ASCII 快速路径优化)
			if j+candidateSize < len(str) {
				if str[j+candidateSize] < utf8.RuneSelf { // Fast path for ASCII
					nextc, nextSize = rune(str[j+candidateSize]), 1
				} else {
					nextc, nextSize = utf8.DecodeRuneInString(str[j+candidateSize:])
				}
			} else {
				nextc, nextSize = 0, 0
			}

			// 当下一个匹配即将到来或搜索字符串结束时,我们应用最佳分数。
			// 跟踪下一个匹配何时到来使我们能够详尽地找到最佳匹配,而不一定是第一个匹配。
			// 例如给定模式 "tk" 和搜索字符串 "The Black Knight",详尽匹配使我们
			// 能够匹配第二个 k,从而给这个字符串更高的分数。
			if equalFold(nextp, nextc) || nextc == 0 {
				if matchedIndex > -1 {
					// 惩罚: 第一个匹配之前的未匹配前导字符
					if len(match.MatchedIndexes) == 0 {
						penalty := matchedIndex * unmatchedLeadingCharPenalty
						bestScore += max(penalty, maxUnmatchedLeadingCharPenalty)
					}

					// 将当前模式字符的匹配结果加入总分和匹配索引列表
					match.Score += bestScore
					match.MatchedIndexes = append(match.MatchedIndexes, matchedIndex)

					// 重置状态,准备匹配下一个模式字符
					score = 0
					bestScore = -1
					matchedIndex = -1
					patternIndex++
				}
			}

			// 保存当前字符状态供下一轮使用
			lastIndex = j
			last = candidate
		}

		// 惩罚: 对字符串中每个未匹配的字符应用惩罚
		// 公式: 匹配字符数 - 字符串总长度 (结果为负值,即惩罚)
		penalty := len(match.MatchedIndexes) - len(str)
		match.Score += penalty

		// 如果所有模式字符都匹配成功,则加入结果列表
		if len(match.MatchedIndexes) == len(runes) {
			matches = append(matches, match)
			matchedIndexes = nil
		} else {
			// 匹配失败,回收 matchedIndexes 切片用于下一个字符串 (减少GC压力)
			matchedIndexes = match.MatchedIndexes[:0]
		}
	}

	return matches
}

// equalFold 比较两个 rune 是否相等 (不区分大小写) 。
// 该函数实现了与 strings.EqualFold 类似的逻辑,支持 Unicode 大小写折叠。
// 先进行 ASCII 快速路径检查,再使用 unicode.SimpleFold 处理一般情况。
//
// 参数:
//   - tr: 第一个要比较的 rune
//   - sr: 第二个要比较的 rune
//
// 返回值:
//   - bool: 如果两个 rune 大小写不敏感相等则返回 true
func equalFold(tr, sr rune) bool {
	if tr == sr {
		return true
	}
	if tr < sr {
		tr, sr = sr, tr
	}
	// ASCII 快速检查。
	if tr < utf8.RuneSelf {
		// ASCII,且 sr 是大写。tr 必须是小写。
		if 'A' <= sr && sr <= 'Z' && tr == sr+'a'-'A' {
			return true
		}
		return false
	}

	// 一般情况。SimpleFold(x) 返回下一个等价的 rune > x
	// 或回绕到更小的值。
	r := unicode.SimpleFold(sr)
	for r != sr && r < tr {
		r = unicode.SimpleFold(r)
	}
	return r == tr
}

// adjacentCharBonus 计算相邻字符匹配的奖励分数。
// 当当前匹配位置与上一个匹配位置相邻时,给予递增的奖励,
// 以鼓励连续字符的匹配,使匹配结果更加紧凑。
//
// 参数:
//   - i: 当前字符位置
//   - lastMatch: 上一个匹配的字符位置
//   - currentBonus: 当前的相邻匹配奖励基数
//
// 返回值:
//   - int: 计算得到的奖励分数 (相邻时为递增奖励,否则为0)
func adjacentCharBonus(i int, lastMatch int, currentBonus int) int {
	if lastMatch == i {
		return currentBonus*2 + adjacentMatchBonus
	}
	return 0
}

// isSeparator 检查给定的 rune 是否为分隔符。
// 分隔符包括空格、下划线、连字符、点号等,用于识别单词边界。
// 匹配分隔符后的字符可获得额外奖励。
//
// 参数:
//   - s: 要检查的 rune
//
// 返回值:
//   - bool: 如果是分隔符则返回 true
func isSeparator(s rune) bool {
	for _, sep := range separators {
		if s == sep {
			return true
		}
	}
	return false
}
