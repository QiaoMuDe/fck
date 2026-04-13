// Package sed 实现文本替换功能
package sed

import (
	"regexp"
	"strings"
)

// replaceWithString 使用字符串替换
//
// 参数:
//   - line: 原始行内容
//   - config: 命令配置指针
//
// 返回:
//   - string: 处理后的行内容
//   - int: 本次替换次数
func replaceWithString(line string, config *SedConfig) (string, int) {
	pattern := config.Pattern
	replacement := config.Replacement

	// 统计匹配次数
	count := strings.Count(line, pattern)
	if count == 0 {
		return line, 0
	}

	// 计算实际可替换次数（考虑全局限制）
	actual := count
	if config.MaxCount > 0 {
		remaining := config.MaxCount - config.replaceCount
		actual = calcActualReplaceCount(count, remaining)
	}

	if actual <= 0 {
		return line, 0
	}

	result := replaceN(line, pattern, replacement, actual)
	return result, actual
}

// replaceStringIgnoreCase 大小写不敏感的字符串替换
//
// 参数:
//   - line: 原始行内容
//   - config: 命令配置指针
//
// 返回:
//   - string: 处理后的行内容
//   - int: 本次替换次数
func replaceStringIgnoreCase(line string, config *SedConfig) (string, int) {
	pattern := config.Pattern
	replacement := config.Replacement

	lowerLine := strings.ToLower(line)
	lowerPattern := strings.ToLower(pattern)

	// 统计匹配次数
	count := 0
	start := 0
	for {
		idx := strings.Index(lowerLine[start:], lowerPattern)
		if idx == -1 {
			break
		}
		count++
		start += idx + len(pattern)
	}

	if count == 0 {
		return line, 0
	}

	// 计算实际可替换次数
	actual := count
	if config.MaxCount > 0 {
		remaining := config.MaxCount - config.replaceCount
		actual = calcActualReplaceCount(count, remaining)
	}

	if actual <= 0 {
		return line, 0
	}

	// 执行替换
	var result strings.Builder
	start = 0
	replaced := 0
	for replaced < actual {
		idx := strings.Index(lowerLine[start:], lowerPattern)
		if idx == -1 {
			break
		}
		actualIdx := start + idx
		result.WriteString(line[start:actualIdx])
		result.WriteString(replacement)
		start = actualIdx + len(pattern)
		replaced++
	}
	result.WriteString(line[start:])

	return result.String(), actual
}

// replaceN 替换前 N 个匹配项
//
// 参数:
//   - s: 原始字符串
//   - old: 匹配模式
//   - new: 替换内容
//   - n: 最大替换次数
//
// 返回:
//   - string: 处理后的字符串
func replaceN(s, old, new string, n int) string {
	if n <= 0 {
		return s
	}

	var result strings.Builder
	start := 0
	count := 0

	for count < n {
		idx := strings.Index(s[start:], old)
		if idx == -1 {
			break
		}

		actualIdx := start + idx
		result.WriteString(s[start:actualIdx])
		result.WriteString(new)
		start = actualIdx + len(old)
		count++
	}

	result.WriteString(s[start:])
	return result.String()
}

// replaceWithRegex 使用正则表达式替换
//
// 参数:
//   - line: 原始行内容
//   - config: 命令配置指针
//
// 返回:
//   - string: 处理后的行内容
//   - int: 本次替换次数
func replaceWithRegex(line string, config *SedConfig) (string, int) {
	if config.compiledPattern == nil {
		return line, 0
	}

	// 检查是否有匹配
	matches := config.compiledPattern.FindAllStringIndex(line, -1)
	if len(matches) == 0 {
		return line, 0
	}

	// 计算实际可替换次数
	actual := len(matches)
	if config.MaxCount > 0 {
		remaining := config.MaxCount - config.replaceCount
		actual = calcActualReplaceCount(len(matches), remaining)
	}

	if actual <= 0 {
		return line, 0
	}

	// 执行替换
	if actual == len(matches) {
		// 全部替换
		result := config.compiledPattern.ReplaceAllString(line, config.Replacement)
		return result, actual
	}

	// 部分替换
	result := replaceRegexN(line, config.compiledPattern, config.Replacement, actual)
	return result, actual
}

// replaceRegexN 替换正则前 N 个匹配项
//
// 参数:
//   - line: 原始行内容
//   - re: 正则表达式
//   - replacement: 替换内容
//   - n: 最大替换次数
//
// 返回:
//   - string: 处理后的行内容
func replaceRegexN(line string, re *regexp.Regexp, replacement string, n int) string {
	if n <= 0 {
		return line
	}

	matches := re.FindAllStringIndex(line, n)
	if len(matches) == 0 {
		return line
	}

	var result strings.Builder
	lastEnd := 0

	for _, match := range matches {
		start, end := match[0], match[1]
		result.WriteString(line[lastEnd:start])
		result.WriteString(replacement)
		lastEnd = end
	}

	result.WriteString(line[lastEnd:])
	return result.String()
}

// inLineRange 检查行号是否在指定范围内
//
// 参数:
//   - config: 命令配置指针
//   - lineNum: 当前行号
//
// 返回:
//   - bool: true 表示在范围内, false 表示不在
func inLineRange(config *SedConfig, lineNum int) bool {
	if config.lineStart == 0 {
		return true // 未指定范围
	}
	if lineNum < config.lineStart {
		return false
	}
	if config.lineEnd > 0 && lineNum > config.lineEnd {
		return false
	}
	return true
}

// calcActualReplaceCount 计算实际可替换次数
//
// 参数:
//   - desired: 期望替换次数
//   - remaining: 剩余可替换次数
//
// 返回:
//   - int: 实际可替换次数
func calcActualReplaceCount(desired, remaining int) int {
	if desired > remaining {
		return remaining
	}
	return desired
}

// processLine 处理单行替换
//
// 参数:
//   - line: 原始行内容
//   - config: 命令配置指针
//   - lineNum: 当前行号
//
// 返回:
//   - string: 处理后的行内容
//   - bool: 是否发生了替换
func processLine(line string, config *SedConfig, lineNum int) (string, bool) {
	// 1. 行范围检查
	if !inLineRange(config, lineNum) {
		return line, false
	}

	// 2. 全局替换次数检查（快速路径）
	if config.MaxCount > 0 && config.replaceCount >= config.MaxCount {
		return line, false
	}

	// 3. 执行替换
	var result string
	var count int

	if config.Regexp {
		// 正则模式 (-E):忽略大小写通过 (?i) 前缀在编译时处理
		result, count = replaceWithRegex(line, config)

	} else if config.IgnoreCase {
		// 字符串模式 + 忽略大小写 (-I 但无 -E)
		result, count = replaceStringIgnoreCase(line, config)

	} else {
		// 普通字符串模式 (无 -E 无 -I)
		result, count = replaceWithString(line, config)
	}

	// 4. 更新全局计数
	if count > 0 {
		config.replaceCount += count
		return result, true
	}

	return line, false
}
