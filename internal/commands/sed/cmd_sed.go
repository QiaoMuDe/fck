// Package sed 实现文本替换功能
// 提供对文件进行文本替换的能力，支持正则表达式和行号范围
package sed

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// SedConfig sed 命令配置
type SedConfig struct {
	// CLI 参数
	Target      string // 目标文件路径
	Pattern     string // -p 匹配模式
	Replacement string // -r 替换内容
	Regexp      bool   // -E 使用正则
	LineRange   string // -n 行号范围
	InPlace     bool   // -i 原地修改
	Backup      bool   // -b 创建备份
	MaxCount    int    // -c 最大替换次数
	IgnoreCase  bool   // -I 忽略大小写

	// 运行时
	compiledPattern *regexp.Regexp // 编译后的正则
	lineStart       int            // 起始行号
	lineEnd         int            // 结束行号（0表示到末尾）
	replaceCount    int            // 当前已替换次数
}

// SedCmdMain 执行 sed 命令
//
// 参数:
//   - config: 命令配置
//
// 返回:
//   - error: 执行错误
func SedCmdMain(config SedConfig) error {
	// 1. 验证参数
	if err := validateConfig(&config); err != nil {
		return err
	}

	// 2. 解析行号范围
	if config.LineRange != "" {
		if err := parseLineRange(&config); err != nil {
			return err
		}
	}

	// 3. 编译正则（如果需要）
	if config.Regexp {
		if err := compilePattern(&config); err != nil {
			return err
		}
	}

	// 4. 处理文件
	return processFile(&config)
}

// validateConfig 验证配置
//
// 参数:
//   - config: 命令配置指针
//
// 返回:
//   - error: 验证错误
func validateConfig(config *SedConfig) error {
	if config.Target == "" {
		return fmt.Errorf("no target file specified")
	}

	if config.Pattern == "" {
		return fmt.Errorf("please use -p to specify the pattern")
	}

	if config.Replacement == "" {
		return fmt.Errorf("please use -r to specify the replacement")
	}

	// 检查文件是否存在
	if _, err := os.Stat(config.Target); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", config.Target)
		}
		return fmt.Errorf("cannot access file: %w", err)
	}

	return nil
}

// parseLineRange 解析行号范围
// 支持格式: "5" (单行), "1,10" (范围), "5," (从第5行到末尾)
//
// 参数:
//   - config: 命令配置指针
//
// 返回:
//   - error: 解析错误
func parseLineRange(config *SedConfig) error {
	rangeStr := config.LineRange

	// 处理 "5," 格式（从第5行到末尾）
	if strings.HasSuffix(rangeStr, ",") {
		startStr := strings.TrimSuffix(rangeStr, ",")
		start, err := strconv.Atoi(startStr)
		if err != nil || start < 1 {
			return fmt.Errorf("invalid line range format: %q, examples: 5 or 1,10", rangeStr)
		}
		config.lineStart = start
		config.lineEnd = 0 // 0 means to the end
		return nil
	}

	// 处理 "1,10" 格式
	if strings.Contains(rangeStr, ",") {
		parts := strings.Split(rangeStr, ",")
		if len(parts) != 2 {
			return fmt.Errorf("invalid line range format: %q, examples: 5 or 1,10", rangeStr)
		}

		start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil || start < 1 {
			return fmt.Errorf("invalid line range format: %q, examples: 5 or 1,10", rangeStr)
		}

		end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil || end < 1 {
			return fmt.Errorf("invalid line range format: %q, examples: 5 or 1,10", rangeStr)
		}

		if start > end {
			return fmt.Errorf("start line %d cannot be greater than end line %d", start, end)
		}

		config.lineStart = start
		config.lineEnd = end
		return nil
	}

	// 处理 "5" 格式（单行）
	lineNum, err := strconv.Atoi(rangeStr)
	if err != nil || lineNum < 1 {
		return fmt.Errorf("invalid line range format: %q, examples: 5 or 1,10", rangeStr)
	}

	config.lineStart = lineNum
	config.lineEnd = lineNum
	return nil
}

// compilePattern 编译正则表达式
//
// 参数:
//   - config: 命令配置指针
//
// 返回:
//   - error: 编译错误
func compilePattern(config *SedConfig) error {
	pattern := config.Pattern

	// 处理忽略大小写
	if config.IgnoreCase {
		pattern = "(?i)" + pattern
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid regular expression: %w", err)
	}

	config.compiledPattern = re
	return nil
}

// processFile 处理文件
//
// 参数:
//   - config: 命令配置指针
//
// 返回:
//   - error: 处理错误
func processFile(config *SedConfig) error {
	// 打开文件
	file, err := os.Open(config.Target)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}
	defer func() { _ = file.Close() }()

	// 读取并处理每一行
	var resultLines []string
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		processedLine, _ := processLine(line, config, lineNum)
		resultLines = append(resultLines, processedLine)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// 输出或写入文件
	if config.InPlace {
		return writeFile(config, resultLines)
	}

	// 输出到标准输出
	for _, line := range resultLines {
		fmt.Println(line)
	}

	return nil
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
	// 检查行号范围
	if config.lineStart > 0 {
		if lineNum < config.lineStart {
			return line, false
		}
		if config.lineEnd > 0 && lineNum > config.lineEnd {
			return line, false
		}
	}

	// 检查最大替换次数
	if config.MaxCount > 0 && config.replaceCount >= config.MaxCount {
		return line, false
	}

	// 执行替换
	if config.Regexp {
		return replaceWithRegex(line, config)
	}
	return replaceWithString(line, config)
}

// replaceWithString 使用字符串替换
//
// 参数:
//   - line: 原始行内容
//   - config: 命令配置指针
//
// 返回:
//   - string: 处理后的行内容
//   - bool: 是否发生了替换
func replaceWithString(line string, config *SedConfig) (string, bool) {
	pattern := config.Pattern
	replacement := config.Replacement

	// 处理忽略大小写
	if config.IgnoreCase {
		// 使用大小写不敏感的方式替换
		return replaceStringIgnoreCase(line, pattern, replacement, config)
	}

	// 检查最大替换次数
	if config.MaxCount > 0 {
		count := strings.Count(line, pattern)
		if count == 0 {
			return line, false
		}

		remaining := config.MaxCount - config.replaceCount
		if count > remaining {
			// 只替换前 remaining 个
			result := replaceN(line, pattern, replacement, remaining)
			config.replaceCount += remaining
			return result, true
		}

		result := strings.ReplaceAll(line, pattern, replacement)
		config.replaceCount += count
		return result, true
	}

	// 无限制替换
	if !strings.Contains(line, pattern) {
		return line, false
	}

	result := strings.ReplaceAll(line, pattern, replacement)
	return result, true
}

// replaceStringIgnoreCase 大小写不敏感的字符串替换
//
// 参数:
//   - line: 原始行内容
//   - pattern: 匹配模式
//   - replacement: 替换内容
//   - config: 命令配置指针
//
// 返回:
//   - string: 处理后的行内容
//   - bool: 是否发生了替换
func replaceStringIgnoreCase(line, pattern, replacement string, config *SedConfig) (string, bool) {
	lowerLine := strings.ToLower(line)
	lowerPattern := strings.ToLower(pattern)

	var result strings.Builder
	start := 0
	replaced := false

	for {
		idx := strings.Index(lowerLine[start:], lowerPattern)
		if idx == -1 {
			break
		}

		// 检查最大替换次数
		if config.MaxCount > 0 && config.replaceCount >= config.MaxCount {
			break
		}

		actualIdx := start + idx
		result.WriteString(line[start:actualIdx])
		result.WriteString(replacement)
		start = actualIdx + len(pattern)
		config.replaceCount++
		replaced = true
	}

	result.WriteString(line[start:])
	return result.String(), replaced
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

	for {
		if count >= n {
			break
		}

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
//   - bool: 是否发生了替换
func replaceWithRegex(line string, config *SedConfig) (string, bool) {
	if config.compiledPattern == nil {
		return line, false
	}

	// 检查是否有匹配
	if !config.compiledPattern.MatchString(line) {
		return line, false
	}

	// 处理最大替换次数
	if config.MaxCount > 0 {
		// 计算匹配次数
		matches := config.compiledPattern.FindAllStringIndex(line, -1)
		if len(matches) == 0 {
			return line, false
		}

		remaining := config.MaxCount - config.replaceCount
		if len(matches) > remaining {
			// 只替换前 remaining 个
			result := replaceRegexN(line, config.compiledPattern, config.Replacement, remaining)
			config.replaceCount += remaining
			return result, true
		}
	}

	// 无限制替换
	result := config.compiledPattern.ReplaceAllString(line, config.Replacement)
	if config.MaxCount > 0 {
		config.replaceCount += len(config.compiledPattern.FindAllString(line, -1))
	}
	return result, true
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

// writeFile 写入文件（原地修改）
//
// 参数:
//   - config: 命令配置指针
//   - lines: 要写入的行列表
//
// 返回:
//   - error: 写入错误
func writeFile(config *SedConfig, lines []string) error {
	// 如果需要备份
	if config.Backup {
		backupPath := config.Target + ".bak"
		if err := copyFile(config.Target, backupPath); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
	}

	// 创建临时文件
	dir := filepath.Dir(config.Target)
	tempFile, err := os.CreateTemp(dir, ".sed-*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tempPath := tempFile.Name()

	// 写入内容
	writer := bufio.NewWriter(tempFile)
	for i, line := range lines {
		if i > 0 {
			writer.WriteByte('\n')
		}
		writer.WriteString(line)
	}

	if err := writer.Flush(); err != nil {
		tempFile.Close()
		os.Remove(tempPath)
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	if err := tempFile.Close(); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	// 原子替换
	if err := os.Rename(tempPath, config.Target); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// copyFile 复制文件
//
// 参数:
//   - src: 源文件路径
//   - dst: 目标文件路径
//
// 返回:
//   - error: 复制错误
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = destFile.ReadFrom(sourceFile)
	return err
}
