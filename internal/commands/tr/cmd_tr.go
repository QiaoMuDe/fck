// Package tr 实现了字符转换功能
// 支持字符映射、删除、压缩等功能
package tr

import (
	"bufio"
	"fmt"
	"os"
	"unicode/utf8"
)

// TrConfig tr 命令配置
type TrConfig struct {
	Set1           string // 源字符集
	Set2           string // 目标字符集
	Delete         bool   // 删除模式
	SqueezeRepeats bool   // 压缩重复字符
	Complement     bool   // 使用补集
	TruncateSet1   bool   // 截断 set1
}

// TrCmdMain 主函数
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - error: 执行错误
func TrCmdMain(config TrConfig) error {
	// 1. 解析字符集
	set1, err := parseCharSet(config.Set1)
	if err != nil {
		return err
	}

	// 2. 根据模式处理
	if config.Delete {
		return deleteMode(set1, config.Complement)
	}

	if config.SqueezeRepeats {
		return squeezeMode(set1)
	}

	// 3. 转换模式（需要 set2）
	if config.Set2 == "" {
		return fmt.Errorf("set2 is required for translate mode")
	}

	set2, err := parseCharSet(config.Set2)
	if err != nil {
		return err
	}

	mapping := buildMapping(set1, set2, config.TruncateSet1)
	return translateMode(mapping)
}

// parseCharSet 解析字符集字符串
//
// 参数:
//   - set: 字符集字符串
//
// 返回值:
//   - []rune: 解析后的字符列表
//   - error: 解析错误
func parseCharSet(set string) ([]rune, error) {
	if set == "" {
		return nil, fmt.Errorf("character set cannot be empty")
	}

	// 先将字符串解码为 rune 数组，正确处理 UTF-8
	runes := []rune(set)
	var result []rune
	i := 0
	for i < len(runes) {
		// 处理转义字符
		if runes[i] == '\\' && i+1 < len(runes) {
			switch runes[i+1] {
			case 'n':
				result = append(result, '\n')
				i += 2
				continue
			case 't':
				result = append(result, '\t')
				i += 2
				continue
			case '\\':
				result = append(result, '\\')
				i += 2
				continue
			}
		}

		// 处理范围（如 a-z, 一-三）
		if i+2 < len(runes) && runes[i+1] == '-' && runes[i+2] > runes[i] {
			start := runes[i]
			end := runes[i+2]
			for r := start; r <= end; r++ {
				result = append(result, r)
			}
			i += 3
			continue
		}

		// 普通字符
		result = append(result, runes[i])
		i++
	}

	return result, nil
}

// buildMapping 构建字符映射表
//
// 参数:
//   - set1: 源字符集
//   - set2: 目标字符集
//   - truncate: 是否截断
//
// 返回值:
//   - map[rune]rune: 字符映射表
func buildMapping(set1, set2 []rune, truncate bool) map[rune]rune {
	mapping := make(map[rune]rune)

	if truncate && len(set1) > len(set2) {
		// 截断 set1 到 set2 的长度
		set1 = set1[:len(set2)]
	}

	for i, r1 := range set1 {
		var r2 rune
		if i < len(set2) {
			r2 = set2[i]
		} else {
			// set1 比 set2 长，使用 set2 最后一个字符
			r2 = set2[len(set2)-1]
		}
		mapping[r1] = r2
	}

	return mapping
}

// deleteMode 删除模式
//
// 参数:
//   - set1: 要删除的字符集
//   - complement: 是否使用补集
//
// 返回值:
//   - error: 错误
func deleteMode(set1 []rune, complement bool) error {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanRunes)

	set1Map := make(map[rune]bool, len(set1))
	for _, r := range set1 {
		set1Map[r] = true
	}

	writer := bufio.NewWriter(os.Stdout)
	defer func() {
		_ = writer.Flush()
	}()

	for scanner.Scan() {
		r, ok := getRuneFromScanner(scanner)
		if !ok {
			continue
		}
		shouldDelete := set1Map[r]
		if complement {
			shouldDelete = !shouldDelete
		}

		if !shouldDelete {
			if _, err := writer.WriteRune(r); err != nil {
				return err
			}
		}
	}

	return scanner.Err()
}

// squeezeMode 压缩模式
//
// 参数:
//   - set1: 要压缩的字符集
//
// 返回值:
//   - error: 错误
func squeezeMode(set1 []rune) error {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanRunes)

	set1Map := make(map[rune]bool, len(set1))
	for _, r := range set1 {
		set1Map[r] = true
	}

	writer := bufio.NewWriter(os.Stdout)
	defer func() {
		_ = writer.Flush()
	}()

	var lastRune rune
	var hasLast bool

	for scanner.Scan() {
		r, ok := getRuneFromScanner(scanner)
		if !ok {
			continue
		}

		if set1Map[r] {
			if hasLast && lastRune == r {
				continue // 跳过重复
			}
		}

		if _, err := writer.WriteRune(r); err != nil {
			return err
		}
		lastRune = r
		hasLast = true
	}

	return scanner.Err()
}

// translateMode 转换模式
//
// 参数:
//   - mapping: 字符映射表
//
// 返回值:
//   - error: 错误
func translateMode(mapping map[rune]rune) error {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanRunes)

	writer := bufio.NewWriter(os.Stdout)
	defer func() {
		_ = writer.Flush()
	}()

	for scanner.Scan() {
		r, ok := getRuneFromScanner(scanner)
		if !ok {
			continue
		}

		if newR, ok := mapping[r]; ok {
			if _, err := writer.WriteRune(newR); err != nil {
				return err
			}
		} else {
			if _, err := writer.WriteRune(r); err != nil {
				return err
			}
		}
	}

	return scanner.Err()
}

// getRuneFromScanner 从 scanner 获取一个 rune
//
// 参数:
//   - scanner: bufio.Scanner 实例
//
// 返回值:
//   - rune: 获取到的字符
//   - bool: 是否成功获取
func getRuneFromScanner(scanner *bufio.Scanner) (rune, bool) {
	text := scanner.Text()
	if len(text) == 0 {
		return 0, false
	}
	r, _ := utf8.DecodeRuneInString(text)
	return r, true
}
