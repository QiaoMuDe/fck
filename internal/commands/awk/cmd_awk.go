package awk

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gitee.com/MM-Q/go-kit/fs"
	"gitee.com/MM-Q/go-kit/term"
)

// AwkCmdMain 执行 awk 命令
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - error: 执行错误
func AwkCmdMain(config AwkConfig) error {
	// 编译正则表达式（如果提供了模式）
	var re *regexp.Regexp
	if config.Pattern != "" {
		var err error
		re, err = regexp.Compile(config.Pattern)
		if err != nil {
			return fmt.Errorf("invalid pattern: %w", err)
		}
	}

	// 优先判断管道输入（管道输入覆盖文件参数）
	if term.IsStdinPipe() {
		return processInput(os.Stdin, config, re)
	}

	// 无管道输入时，处理文件
	if len(config.Files) == 0 {
		return fmt.Errorf("no input file specified and no pipe input detected")
	}

	// 展开通配符并收集所有文件
	allFiles, err := fs.Expand(config.Files)
	if err != nil {
		return err
	}

	for _, file := range allFiles {
		if err := processFile(file, config, re); err != nil {
			return err
		}
	}

	return nil
}

// processFile 处理单个文件
//
// 参数:
//   - filename: 文件名
//   - config: 命令配置
//   - re: 正则表达式（可为 nil）
//
// 返回值:
//   - error: 处理错误
func processFile(filename string, config AwkConfig, re *regexp.Regexp) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("cannot open file %s: %w", filename, err)
	}
	defer func() {
		_ = file.Close()
	}()

	return processInput(file, config, re)
}

// processInput 处理输入流（支持大文件和管道）
//
// 使用 bufio.Reader 替代 Scanner，支持任意行长度
//
// 参数:
//   - input: 输入文件
//   - config: 命令配置
//   - re: 正则表达式（可为 nil）
//
// 返回值:
//   - error: 处理错误
func processInput(input *os.File, config AwkConfig, re *regexp.Regexp) error {
	reader := bufio.NewReader(input)
	lineNum := 0

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// 处理最后一行（可能没有换行符）
				if len(line) > 0 {
					lineNum++
					if err := processLine(lineNum, line, config, re); err != nil {
						return err
					}
				}
				break
			}
			return fmt.Errorf("read error: %w", err)
		}

		// 去掉末尾换行符
		line = strings.TrimSuffix(line, "\n")
		line = strings.TrimSuffix(line, "\r") // Windows 兼容

		lineNum++
		if err := processLine(lineNum, line, config, re); err != nil {
			return err
		}
	}

	return nil
}

// processLine 处理单行
//
// 参数:
//   - lineNum: 行号
//   - line: 行内容
//   - config: 命令配置
//   - re: 正则表达式（可为 nil）
//
// 返回值:
//   - error: 处理错误
func processLine(lineNum int, line string, config AwkConfig, re *regexp.Regexp) error {
	// 模式匹配检查
	if re != nil && !re.MatchString(line) {
		return nil
	}

	// 字符提取模式不需要分割字段
	if config.IsCharMode() {
		return outputLine(lineNum, line, nil, 0, config)
	}

	// 分割字段（未指定分隔符时使用任意空白）
	var fields []string
	if config.FieldSep == "" {
		fields = strings.Fields(line)
	} else {
		fields = strings.Split(line, config.FieldSep)
	}
	nf := len(fields)

	// 输出处理
	if err := outputLine(lineNum, line, fields, nf, config); err != nil {
		return err
	}

	return nil
}

// outputLine 输出行
//
// 参数:
//   - nr: 行号
//   - line: 整行内容
//   - fields: 分割后的字段
//   - nf: 字段数量
//   - config: 命令配置
//
// 返回值:
//   - error: 输出错误
func outputLine(nr int, line string, fields []string, nf int, config AwkConfig) error {
	var parts []string

	// 添加行号
	if config.ShowLineNum {
		parts = append(parts, fmt.Sprintf("%d", nr))
	}

	// 字符提取模式
	if config.IsCharMode() {
		for _, cr := range config.CharRanges {
			value := cr.Extract(line)
			parts = append(parts, value)
		}
	} else {
		// 字段提取模式（原有逻辑）
		for _, idx := range config.Fields {
			var value string
			switch idx {
			case 0: // 整行
				value = line
			case -1: // 最后一个字段
				if nf > 0 {
					value = fields[nf-1]
				}
			default: // 正常字段（1-based）
				if idx > 0 && idx <= nf {
					value = fields[idx-1]
				}
			}
			parts = append(parts, value)
		}
	}

	// 输出
	fmt.Println(strings.Join(parts, config.OutputSep))
	return nil
}

// ParseCharRanges 解析字符范围字符串
// 支持格式: "1-10,15,20-,5"
func ParseCharRanges(s string) ([]CharRange, error) {
	var ranges []CharRange
	parts := strings.Split(s, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		cr, err := parseSingleRange(part)
		if err != nil {
			return nil, fmt.Errorf("invalid range %q: %w", part, err)
		}
		ranges = append(ranges, cr)
	}

	return ranges, nil
}

// parseSingleRange 解析单个范围
// 支持: "5" -> {5,5}, "1-10" -> {1,10}, "5-" -> {5,-1}, "-10" -> {1,10}
func parseSingleRange(s string) (CharRange, error) {
	// 处理 "-N" 格式（前N个字符）
	if strings.HasPrefix(s, "-") {
		n, err := strconv.Atoi(s[1:])
		if err != nil || n < 1 {
			return CharRange{}, fmt.Errorf("invalid range")
		}
		return CharRange{Start: 1, End: n}, nil
	}

	// 处理 "N-" 格式（从N到结尾）
	if strings.HasSuffix(s, "-") {
		n, err := strconv.Atoi(s[:len(s)-1])
		if err != nil || n < 1 {
			return CharRange{}, fmt.Errorf("invalid range")
		}
		return CharRange{Start: n, End: -1}, nil
	}

	// 处理 "N-M" 格式
	if strings.Contains(s, "-") {
		parts := strings.SplitN(s, "-", 2)
		start, err1 := strconv.Atoi(parts[0])
		end, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil || start < 1 || end < 1 || end < start {
			return CharRange{}, fmt.Errorf("invalid range")
		}
		return CharRange{Start: start, End: end}, nil
	}

	// 处理单个数字 "N"
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 {
		return CharRange{}, fmt.Errorf("invalid range")
	}
	return CharRange{Start: n, End: n}, nil
}
