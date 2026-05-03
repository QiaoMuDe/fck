package wc

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"gitee.com/MM-Q/fck/internal/types"
	"gitee.com/MM-Q/go-kit/term"
	"github.com/jedib0t/go-pretty/v6/table"
)

// WcCmdMain 执行 wc 命令
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - error: 执行错误
func WcCmdMain(config WcConfig) error {
	// 如果启用 --all, 显示所有统计项
	if config.ShowAll {
		config.ShowLines = true
		config.ShowWords = true
		config.ShowBytes = true
		config.ShowChars = true
		config.ShowMaxLine = true
	}

	// 如果没有指定任何显示选项，默认显示行数、单词数、字节数
	if !config.ShowLines && !config.ShowWords && !config.ShowBytes &&
		!config.ShowChars && !config.ShowMaxLine {
		config.ShowLines = true
		config.ShowWords = true
		config.ShowBytes = true
	}

	// 优先判断管道输入
	if term.IsStdinPipe() {
		stats, err := countStats(os.Stdin, "-")
		if err != nil {
			return fmt.Errorf("failed to count stdin: %w", err)
		}
		printStatsTable([]*WcStats{stats}, config, false)
		return nil
	}

	// 无管道输入时，处理文件
	if len(config.Files) == 0 {
		return fmt.Errorf("no input file specified and no pipe input detected")
	}

	// 展开通配符并收集所有文件
	var allFiles []string
	for _, pattern := range config.Files {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return fmt.Errorf("invalid pattern %s: %w", pattern, err)
		}
		if len(matches) == 0 {
			// 检查是否包含通配符字符
			if strings.ContainsAny(pattern, "*?[") {
				return fmt.Errorf("no files match pattern: %s", pattern)
			}
			allFiles = append(allFiles, pattern)
		} else {
			allFiles = append(allFiles, matches...)
		}
	}

	// 处理多个文件
	var allStats []*WcStats
	var total WcStats
	for _, file := range allFiles {
		stats, err := processFile(file)
		if err != nil {
			return err
		}
		allStats = append(allStats, stats)
		total.Lines += stats.Lines
		total.Words += stats.Words
		total.Bytes += stats.Bytes
		total.Chars += stats.Chars
		if stats.MaxLineLen > total.MaxLineLen {
			total.MaxLineLen = stats.MaxLineLen
		}
	}

	// 多文件时添加总计
	if len(allFiles) > 1 {
		total.Filename = "total"
		allStats = append(allStats, &total)
	}

	// 使用表格输出
	printStatsTable(allStats, config, true)

	return nil
}

// processFile 处理单个文件
//
// 参数:
//   - filename: 文件名
//
// 返回值:
//   - *WcStats: 统计结果
//   - error: 处理错误
func processFile(filename string) (*WcStats, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot open file %s: %w", filename, err)
	}
	defer func() {
		_ = file.Close()
	}()

	return countStats(file, filename)
}

// countStats 统计 reader 中的数据
//
// 参数:
//   - reader: 输入流
//   - filename: 文件名（用于标识）
//
// 返回值:
//   - *WcStats: 统计结果
//   - error: 读取错误
func countStats(reader io.Reader, filename string) (*WcStats, error) {
	stats := &WcStats{Filename: filename}
	scanner := bufio.NewScanner(reader)

	// 创建一个缓冲区，用于处理大文件和长行
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		line := scanner.Text()
		stats.Lines++

		// 字节数（包含换行符，与 Linux wc -c 保持一致）
		// scanner.Bytes() 返回不含换行符的内容，需要 +1 计算换行符
		lineBytes := len(scanner.Bytes()) + 1
		stats.Bytes += int64(lineBytes)

		// 字符数（Unicode）
		stats.Chars += int64(utf8.RuneCountInString(line))

		// 最大行长度（字符数）
		lineLen := int64(len([]rune(line)))
		if lineLen > stats.MaxLineLen {
			stats.MaxLineLen = lineLen
		}

		// 单词数（以空白分隔）
		stats.Words += int64(len(splitFields(line)))
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read error: %w", err)
	}

	return stats, nil
}

// splitFields 分割字段（以任意空白分隔）
func splitFields(s string) []string {
	var fields []string
	fieldStart := -1

	for i, r := range s {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' || r == '\f' || r == '\v' {
			if fieldStart >= 0 {
				fields = append(fields, s[fieldStart:i])
				fieldStart = -1
			}
		} else {
			if fieldStart < 0 {
				fieldStart = i
			}
		}
	}

	if fieldStart >= 0 {
		fields = append(fields, s[fieldStart:])
	}

	return fields
}

// printStatsTable 使用表格输出统计结果
//
// 参数:
//   - allStats: 统计结果列表
//   - config: 显示配置
//   - showFilename: 是否显示文件名
func printStatsTable(allStats []*WcStats, config WcConfig, showFilename bool) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// 构建表头
	var header table.Row
	if config.ShowLines {
		header = append(header, "Lines")
	}
	if config.ShowWords {
		header = append(header, "Words")
	}
	if config.ShowChars {
		header = append(header, "Chars")
	}
	if config.ShowBytes {
		header = append(header, "Bytes")
	}
	if config.ShowMaxLine {
		header = append(header, "MaxLine")
	}
	if showFilename {
		header = append(header, "File")
	}
	t.AppendHeader(header)

	// 构建数据行
	for _, stats := range allStats {
		var row table.Row
		if config.ShowLines {
			row = append(row, stats.Lines)
		}
		if config.ShowWords {
			row = append(row, stats.Words)
		}
		if config.ShowChars {
			row = append(row, stats.Chars)
		}
		if config.ShowBytes {
			row = append(row, stats.Bytes)
		}
		if config.ShowMaxLine {
			row = append(row, stats.MaxLineLen)
		}
		if showFilename {
			row = append(row, stats.Filename)
		}
		t.AppendRow(row)
	}

	// 设置表格样式
	if style, ok := types.TableStyleMap[config.TableStyle]; ok {
		t.SetStyle(style)
	}

	t.Render()
}
