// Package sed 实现文本替换功能
// 提供对文件进行文本替换的能力, 支持正则表达式和行号范围
package sed

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"gitee.com/MM-Q/fck/internal/types"
	"gitee.com/MM-Q/go-kit/fs"
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
	MaxBuffer   int64  // --buffer-size 最大行缓冲区大小(字节)

	// 运行时
	compiledPattern *regexp.Regexp // 编译后的正则
	lineStart       int            // 起始行号
	lineEnd         int            // 结束行号 (0表示到末尾)
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
// 根据配置选择预览模式或原地修改模式, 采用流式处理避免大文件OOM
//
// 参数:
//   - config: 命令配置指针
//
// 返回:
//   - error: 处理错误
func processFile(config *SedConfig) error {
	// 原地修改模式
	if config.InPlace {
		return processFileInPlace(config)
	}

	// 预览模式
	return processFilePreview(config)
}

// processFilePreview 预览模式处理文件
// 边读取边处理边输出到标准输出, 不缓存任何行, 内存占用O(1)
//
// 参数:
//   - config: 命令配置指针
//
// 返回:
//   - error: 处理错误
func processFilePreview(config *SedConfig) error {
	file, err := os.Open(config.Target)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}
	defer func() { _ = file.Close() }()

	// 设置动态缓冲区: 初始64KB, 最大可扩容到配置值
	scanner := bufio.NewScanner(file)
	buf := make([]byte, types.InitialBufferSize)
	scanner.Buffer(buf, int(config.MaxBuffer))

	// 遍历文件行, 边读边处理边输出
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		processedLine, _ := processLine(line, config, lineNum)
		fmt.Println(processedLine)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	return nil
}

// processFileInPlace 原地修改模式处理文件
// 边读取边处理边写入临时文件, 最后原子替换, 内存占用O(1)
//
// 参数:
//   - config: 命令配置指针
//
// 返回:
//   - error: 处理错误
func processFileInPlace(config *SedConfig) (err error) {
	// 如果需要备份, 先创建备份（此时源文件还未打开）
	// 使用 CopyEx 允许覆盖已存在的备份文件, 与 Linux sed 行为一致
	if config.Backup {
		backupPath := config.Target + types.SedBackupSuffix
		if err := fs.CopyEx(config.Target, backupPath, true); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
	}

	// 打开源文件
	sourceFile, err := os.Open(config.Target)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}

	// 创建临时文件（在同一目录, 确保原子重命名有效）
	dir := filepath.Dir(config.Target)
	tempFile, err := os.CreateTemp(dir, types.SedTempFilePattern)
	if err != nil {
		_ = sourceFile.Close()
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tempPath := tempFile.Name()

	// 延迟清理: 只在出错时删除临时文件
	defer func() {
		if err != nil {
			_ = tempFile.Close()
			_ = os.Remove(tempPath)
		}
	}()

	// 设置动态缓冲区: 初始64KB, 最大可扩容到配置值
	scanner := bufio.NewScanner(sourceFile)
	buf := make([]byte, types.InitialBufferSize)
	scanner.Buffer(buf, int(config.MaxBuffer))

	// 使用缓冲写入器提高性能
	writer := bufio.NewWriter(tempFile)

	lineNum := 0
	firstLine := true
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		processedLine, _ := processLine(line, config, lineNum)

		// 除第一行外, 每行前面添加换行符
		if !firstLine {
			if err := writer.WriteByte('\n'); err != nil {
				_ = sourceFile.Close()
				return fmt.Errorf("failed to write temp file: %w", err)
			}
		}
		firstLine = false

		if _, err := writer.WriteString(processedLine); err != nil {
			_ = sourceFile.Close()
			return fmt.Errorf("failed to write temp file: %w", err)
		}
	}

	// 关闭源文件（Windows 上必须先关闭才能重命名）
	if err := sourceFile.Close(); err != nil {
		return fmt.Errorf("failed to close source file: %w", err)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// 刷新缓冲区
	if err = writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush temp file: %w", err)
	}

	// 关闭临时文件
	if err = tempFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	// 原子替换
	if err = os.Rename(tempPath, config.Target); err != nil {
		return fmt.Errorf("failed to replace file: %w", err)
	}

	return nil
}
