package grep

import (
	"bufio"
	"bytes"
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

	// 二进制文件处理
	Text         bool // -a, --text 强制将二进制文件视为文本处理
	IgnoreBinary bool // -I 忽略二进制文件（不输出提示，静默跳过）

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
	// 情况1: 未指定目标文件（从 stdin 读取）
	if config.Target == "" {
		// 递归标志与管道输入冲突时，优先处理管道
		if (config.Recursive || config.FollowSymlink) && isStdinPipe() {
			if !config.NoMessages {
				fmt.Fprintln(os.Stderr, "-r/-R flags ignored when reading from stdin")
			}
			config.Recursive = false
			config.FollowSymlink = false
		}

		// 从标准输入读取
		if !config.Recursive && !config.FollowSymlink {
			config.filename = "(standard input)"
			return processFile(os.Stdin, &config)
		}

		// 递归模式默认使用当前目录
		config.Target = "."
	}

	// 情况2: 检查目标类型（文件或目录）
	info, err := os.Stat(config.Target)
	if err != nil {
		if config.NoMessages {
			return nil
		}
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", config.Target)
		}
		return fmt.Errorf("cannot access file: %w", err)
	}

	// 根据类型处理
	if info.IsDir() {
		// 目录：需要 -r 或 -R 标志
		if !config.Recursive && !config.FollowSymlink {
			if config.NoMessages {
				return nil
			}
			return fmt.Errorf("is a directory, use -r or -R for recursive search: %s", config.Target)
		}
		return processPathRecursive(config.Target, &config)
	}

	// 情况3: 普通文件：直接处理
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

	skip, err := handleBinaryFile(file, path, config)
	if err != nil {
		return err
	}
	if skip {
		return nil
	}

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

// handleBinaryFile 处理二进制文件检测和配置响应
//
// 根据 config.Text 和 config.IgnoreBinary 决定如何处理二进制文件
//
// 参数:
//   - file: 已打开的文件句柄
//   - path: 文件路径（用于输出提示）
//   - config: 命令配置
//
// 返回:
//   - skip: true 表示应跳过此文件, false 表示继续处理
//   - err: 检测错误 (非静默模式下)
func handleBinaryFile(file *os.File, path string, config *GrepConfig) (skip bool, err error) {
	isBinary, err := isBinaryFile(file)
	if err != nil {
		if config.NoMessages {
			return true, nil // 静默模式下跳过错误文件
		}
		return true, fmt.Errorf("failed to detect file type: %w", err)
	}

	if !isBinary {
		return false, nil // 文本文件，继续处理
	}

	// 二进制文件，根据配置处理
	if config.Text {
		// -a/--text 模式：强制作为文本处理
		return false, nil
	}
	if config.IgnoreBinary {
		// -I 模式：静默跳过
		return true, nil
	}

	// 默认行为：输出提示并跳过
	if config.WithFilename {
		fmt.Printf("Binary file %s matches\n", path)
	} else {
		fmt.Println("Binary file matches")
	}
	return true, nil
}

// isStdinPipe 检查标准输入是否为管道/重定向（而非终端）
//
// 返回:
//   - bool: true 表示 stdin 是管道或文件重定向, false 表示是终端输入
func isStdinPipe() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	// 检查是否为命名管道或常规文件（重定向）
	return info.Mode()&os.ModeNamedPipe != 0 || info.Mode().IsRegular()
}

// isBinaryFile 检测文件是否为二进制文件
//
// 原理：读取文件前 8000 字节，检查是否包含空字符(\0)
//
// 注意：
//   - 只支持普通文件, stdin/pipe 默认返回 false (视为文本)
//   - 检测后会重置文件指针到开头
//
// 参数:
//   - file: 已打开的文件句柄
//
// 返回:
//   - bool: true 表示二进制文件, false 表示文本文件或无法检测
//   - error: 读取或重置指针错误
func isBinaryFile(file *os.File) (bool, error) {
	// 获取文件信息，检查是否为普通文件
	info, err := file.Stat()
	if err != nil {
		return false, err
	}

	// 非普通文件 (stdin、pipe、设备文件等) 默认视为文本
	if !info.Mode().IsRegular() {
		return false, nil
	}

	// 读取文件前 8000 字节
	buf := make([]byte, 8000)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return false, err
	}

	// 检查是否包含空字符（二进制文件特征）
	isBinary := bytes.Contains(buf[:n], []byte{0})

	// 重置文件指针到开头，供后续读取使用
	if _, seekErr := file.Seek(0, io.SeekStart); seekErr != nil {
		return false, seekErr
	}

	return isBinary, nil
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
