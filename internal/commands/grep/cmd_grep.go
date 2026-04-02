package grep

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// 颜色常量
const (
	ColorRed   = "\033[1;31m" // 红色加粗
	ColorBlue  = "\033[1;34m" // 蓝色加粗
	ColorCyan  = "\033[1;36m" // 青色加粗
	ColorWhite = "\033[1;37m" // 白色加粗
	ColorReset = "\033[0m"    // 重置
)

// GrepConfig grep 命令配置
type GrepConfig struct {
	// CLI 参数
	Pattern       string // 搜索模式
	Target        string // 目标文件（空表示标准输入）
	IgnoreCase    bool   // -i 忽略大小写
	InvertMatch   bool   // -v 反向匹配
	LineNumber    bool   // -n 显示行号
	Count         bool   // -c 只显示计数
	WithFilename  bool   // -H 显示文件名
	Regexp        bool   // -E 使用正则表达式
	MaxCount      int    // -m 最大匹配数
	BeforeContext int    // -B 前文行数
	AfterContext  int    // -A 后文行数
	Context       int    // -C 上下文行数
	NoColor       bool   // --no-color 禁用颜色

	// 递归搜索相关
	Recursive     bool     // -r 递归搜索
	FollowSymlink bool     // -R 递归并跟随符号链接
	Include       []string // --include 包含文件模式
	Exclude       []string // --exclude 排除文件模式
	IncludeDir    []string // --include-dir 包含目录模式
	ExcludeDir    []string // --exclude-dir 排除目录模式
	NoMessages    bool     // -s 静默模式，不显示错误信息

	// 隐藏文件支持
	Hidden bool // --hidden 处理隐藏文件/目录（默认跳过）

	// 运行时
	compiledPattern *regexp.Regexp // 编译后的正则
	matchCount      int            // 当前匹配计数
	filename        string         // 实际文件名（用于输出）
	lastPrintedFile string         // 最后打印的文件名（用于递归模式分组显示）
}

// lineInfo 存储行信息
type lineInfo struct {
	content string // 行内容
	lineNum int    // 行号
}

// GrepCmdMain 执行 grep 命令
//
// 参数:
//   - config: 命令配置
//
// 返回:
//   - error: 执行错误 (如果有)
func GrepCmdMain(config GrepConfig) error {
	// 1. 验证参数
	if config.Pattern == "" {
		return fmt.Errorf("no pattern specified")
	}

	// 2. 处理 -C 选项（同时设置 -B 和 -A）
	if config.Context > 0 {
		if config.BeforeContext == 0 {
			config.BeforeContext = config.Context
		}
		if config.AfterContext == 0 {
			config.AfterContext = config.Context
		}
	}

	// 3. 编译模式
	if err := compilePattern(&config); err != nil {
		return fmt.Errorf("invalid pattern: %w", err)
	}

	// 4. 确定输入源并处理
	switch {
	case config.Target == "" && (config.Recursive || config.FollowSymlink):
		// 递归模式未指定目标，默认使用当前目录
		config.Target = "."

	case config.Target == "":
		// 非递归模式未指定目标，从标准输入读取
		config.filename = "(standard input)"
		return processFile(os.Stdin, &config)
	}

	// 检查是否为递归模式
	if config.Recursive || config.FollowSymlink {
		return processPathRecursive(config.Target, &config)
	}

	// 单文件模式
	return processSingleFileWithError(config.Target, &config)
}

// processSingleFileWithError 处理单个文件（单文件模式）
// 根据 NoMessages 配置决定是否返回错误
//
// 参数:
//   - path: 文件路径
//   - config: 命令配置
//
// 返回:
//   - error: 处理错误 (如果 NoMessages 为 false)
func processSingleFileWithError(path string, config *GrepConfig) error {
	// 检查是否为隐藏文件（默认跳过）
	if !config.Hidden && isHidden(filepath.Base(path)) {
		if config.NoMessages {
			return nil // 静默模式，跳过错误
		}
		return fmt.Errorf("hidden file skipped (use --hidden/-a to include): %s", path)
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		if config.NoMessages {
			return nil // 静默模式，跳过错误
		}
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", path)
		}
		return fmt.Errorf("cannot access file: %w", err)
	}
	if fileInfo.IsDir() {
		if config.NoMessages {
			return nil // 静默模式，跳过错误
		}
		return fmt.Errorf("is a directory, use -r or -R for recursive search: %s", path)
	}

	file, err := os.Open(path)
	if err != nil {
		if config.NoMessages {
			return nil // 静默模式，跳过错误
		}
		return fmt.Errorf("cannot open file: %w", err)
	}
	defer func() { _ = file.Close() }()

	config.filename = path
	return processFile(file, config)
}

// compilePattern 编译搜索模式
//
// 参数:
//   - config: 命令配置指针
//
// 返回:
//   - error: 编译错误 (如果有)
func compilePattern(config *GrepConfig) error {
	// 非正则模式：不需要编译，直接返回
	if !config.Regexp {
		return nil
	}

	// 正则模式：编译正则表达式
	pattern := config.Pattern
	flags := ""
	if config.IgnoreCase {
		flags = "(?i)"
	}

	re, err := regexp.Compile(flags + pattern)
	if err != nil {
		return err
	}
	config.compiledPattern = re
	return nil
}

// processFile 处理文件内容
//
// 参数:
//   - reader: 输入读取器
//   - config: 命令配置
//
// 返回:
//   - error: 处理错误 (如果有)
func processFile(reader io.Reader, config *GrepConfig) error {
	scanner := bufio.NewScanner(reader)

	// 用于上下文显示的行缓冲区
	var beforeLines []lineInfo
	var afterLinesToPrint int
	var lastMatchLine int

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// 匹配检查
		matched := matchLine(line, config)

		if matched {
			config.matchCount++

			// 达到最大匹配数则停止
			if config.MaxCount > 0 && config.matchCount > config.MaxCount {
				break
			}

			// 打印前文
			if config.BeforeContext > 0 {
				for _, info := range beforeLines {
					if info.lineNum > lastMatchLine {
						printLine(info.content, info.lineNum, "-", config)
					}
				}
			}

			// 打印匹配行
			printLine(line, lineNum, ":", config)
			lastMatchLine = lineNum

			// 设置需要打印的后文行数
			afterLinesToPrint = config.AfterContext
		} else {
			// 打印后文（如果还有剩余）
			if afterLinesToPrint > 0 {
				printLine(line, lineNum, "-", config)
				afterLinesToPrint--
				if afterLinesToPrint == 0 {
					// 后文结束，打印分隔行（如果有多个匹配）
					if config.matchCount > 0 && config.matchCount < config.MaxCount {
						fmt.Println("--")
					}
				}
			}
		}

		// 维护前文缓冲区
		if config.BeforeContext > 0 {
			beforeLines = append(beforeLines, lineInfo{content: line, lineNum: lineNum})
			if len(beforeLines) > config.BeforeContext {
				beforeLines = beforeLines[1:]
			}
		}
	}

	// 6. 计数模式输出
	switch {
	case config.Count && (config.Recursive || config.FollowSymlink):
		// 递归模式：先输出文件名，再输出计数，最后空行分隔
		fmt.Println(formatPathWithColor(config.filename, config.NoColor))
		fmt.Println(config.matchCount)
		fmt.Println()

	case config.Count:
		// 单文件模式：只输出计数
		fmt.Println(config.matchCount)
	}

	return scanner.Err()
}

// matchLine 检查行是否匹配
//
// 参数:
//   - line: 行内容
//   - config: 命令配置
//
// 返回:
//   - bool: 是否匹配
func matchLine(line string, config *GrepConfig) bool {
	var matched bool

	if config.Regexp {
		// 正则模式：使用编译后的正则
		matched = config.compiledPattern.MatchString(line)
	} else {
		// 固定字符串模式：使用 strings.Contains
		pattern := config.Pattern
		if config.IgnoreCase {
			// 忽略大小写：都转小写再比较
			matched = strings.Contains(strings.ToLower(line), strings.ToLower(pattern))
		} else {
			matched = strings.Contains(line, pattern)
		}
	}

	if config.InvertMatch {
		matched = !matched
	}
	return matched
}

// highlightLine 对匹配的关键字进行高亮
//
// 参数:
//   - line: 原始行内容
//   - config: 命令配置
//
// 返回:
//   - string: 高亮后的行内容
func highlightLine(line string, config *GrepConfig) string {
	// 禁用颜色或计数模式，直接返回
	if config.NoColor || config.Count {
		return line
	}

	if config.Regexp {
		return highlightRegex(line, config)
	}
	return highlightString(line, config)
}

// highlightRegex 正则模式高亮
//
// 参数:
//   - line: 原始行内容
//   - config: 命令配置
//
// 返回:
//   - string: 高亮后的行内容
func highlightRegex(line string, config *GrepConfig) string {
	// 使用 FindAllStringIndex 找到所有匹配位置
	matches := config.compiledPattern.FindAllStringIndex(line, -1)
	if matches == nil {
		return line
	}

	// 从后向前插入颜色码，避免位置偏移
	result := line
	for i := len(matches) - 1; i >= 0; i-- {
		start, end := matches[i][0], matches[i][1]
		matched := result[start:end]
		result = result[:start] + ColorRed + matched + ColorReset + result[end:]
	}
	return result
}

// highlightString 固定字符串模式高亮
//
// 参数:
//   - line: 原始行内容
//   - config: 命令配置
//
// 返回:
//   - string: 高亮后的行内容
func highlightString(line string, config *GrepConfig) string {
	pattern := config.Pattern
	if pattern == "" {
		return line
	}

	// 根据忽略大小写选项处理
	searchLine := line
	searchPattern := pattern
	if config.IgnoreCase {
		searchLine = strings.ToLower(line)
		searchPattern = strings.ToLower(pattern)
	}

	// 找到所有匹配位置
	var positions [][2]int
	start := 0
	for {
		idx := strings.Index(searchLine[start:], searchPattern)
		if idx == -1 {
			break
		}
		actualStart := start + idx
		actualEnd := actualStart + len(pattern)
		positions = append(positions, [2]int{actualStart, actualEnd})
		start = actualEnd
	}

	if len(positions) == 0 {
		return line
	}

	// 从后向前插入颜色码，避免位置偏移
	// 必须从后向前处理，因为插入颜色码会增加字符串长度，影响后续位置
	result := line
	for i := len(positions) - 1; i >= 0; i-- {
		start, end := positions[i][0], positions[i][1]
		// 从原始 line 取匹配内容，避免 result 中已插入的颜色码干扰
		matched := line[start:end]
		result = result[:start] + ColorRed + matched + ColorReset + result[end:]
	}
	return result
}

// printLine 输出行
// 根据模式选择单文件格式或递归格式
//
// 参数:
//   - line: 行内容
//   - lineNum: 行号
//   - separator: 分隔符（: 表示匹配行，- 表示上下文行）
//   - config: 命令配置
func printLine(line string, lineNum int, separator string, config *GrepConfig) {
	// 计数模式不输出具体内容
	if config.Count {
		return
	}

	// 高亮匹配内容（仅匹配行）
	highlightedLine := line
	if separator == ":" {
		highlightedLine = highlightLine(line, config)
	}

	// 根据模式选择输出格式
	if config.Recursive || config.FollowSymlink {
		printLineRecursive(highlightedLine, lineNum, separator, config)
	} else {
		printLineSingle(highlightedLine, lineNum, separator, config)
	}
}

// formatPathWithColor 格式化路径（带颜色）
// 整个路径统一使用蓝色
//
// 参数:
//   - path: 文件路径
//   - noColor: 是否禁用颜色
//
// 返回:
//   - string: 格式化后的路径
func formatPathWithColor(path string, noColor bool) string {
	if noColor {
		return path
	}

	// 整个路径统一蓝色
	return ColorBlue + path + ColorReset
}

// printLineSingle 单文件模式输出行
// 格式: 文件名:行号:内容 或 行号:内容 或 内容
//
// 参数:
//   - highlightedLine: 高亮后的行内容
//   - lineNum: 行号
//   - separator: 分隔符
//   - config: 命令配置
func printLineSingle(highlightedLine string, lineNum int, separator string, config *GrepConfig) {
	var result string

	// 文件名（蓝色）
	if config.WithFilename {
		result += formatPathWithColor(config.filename, config.NoColor)
	}

	// 行号（青色）
	if config.LineNumber {
		if result != "" {
			result += ":"
		}
		if !config.NoColor {
			result += ColorCyan + strconv.Itoa(lineNum) + ColorReset
		} else {
			result += strconv.Itoa(lineNum)
		}
	}

	// 内容（红色高亮匹配）
	if result != "" {
		result += separator + highlightedLine
	} else {
		result = highlightedLine
	}

	fmt.Println(result)
}

// printLineRecursive 递归模式输出行
// 格式: 文件名单独一行，然后行号:内容（文件之间空行分隔）
//
// 参数:
//   - highlightedLine: 高亮后的行内容
//   - lineNum: 行号
//   - separator: 分隔符
//   - config: 命令配置
func printLineRecursive(highlightedLine string, lineNum int, separator string, config *GrepConfig) {
	// 检查是否需要打印文件名（文件变化时）
	if config.filename != config.lastPrintedFile {
		// 如果不是第一个文件，先打印空行分隔
		if config.lastPrintedFile != "" {
			fmt.Println()
		}
		// 打印文件名（蓝色）
		fmt.Println(formatPathWithColor(config.filename, config.NoColor))
		config.lastPrintedFile = config.filename
	}

	// 打印行号和内容
	var result string
	if config.LineNumber {
		if !config.NoColor {
			result += ColorCyan + strconv.Itoa(lineNum) + ColorReset
		} else {
			result += strconv.Itoa(lineNum)
		}
		result += separator + highlightedLine
	} else {
		result = highlightedLine
	}
	fmt.Println(result)
}

// isHidden 判断文件或目录是否为隐藏项
// 在 Unix/Linux 系统中，以 . 开头的文件/目录为隐藏项
// 在 Windows 系统中，同样以 . 开头判断（大多数工具行为）
//
// 参数:
//   - name: 文件或目录名称
//
// 返回:
//   - bool: 是否为隐藏项
func isHidden(name string) bool {
	return strings.HasPrefix(name, ".")
}
