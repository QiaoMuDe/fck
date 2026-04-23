// Package seq 实现了数字序列生成命令
// 类似于 Unix/Linux 的 seq 命令，支持指定起始值、结束值、步长
package seq

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// SeqConfig 配置结构体
type SeqConfig struct {
	Start      float64 // 起始值
	End        float64 // 结束值
	Step       float64 // 步长
	Separator  string  // 分隔符（默认换行）
	Terminator string  // 终止符（默认换行）
	EqualWidth bool    // 等宽输出（自动补零）
}

// SeqCmdMain 执行 seq 命令
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - error: 执行错误
func SeqCmdMain(config SeqConfig) error {
	// 验证步长
	if config.Step == 0 {
		return fmt.Errorf("step cannot be zero")
	}

	// 处理转义字符
	separator := unescape(config.Separator)
	terminator := unescape(config.Terminator)

	// 生成序列
	sequence := generateSequence(config.Start, config.End, config.Step)

	// 如果没有数据，直接返回
	if len(sequence) == 0 {
		return nil
	}

	// 确定格式
	format := ""
	if config.EqualWidth {
		format = calculateEqualWidthFormat(sequence)
	}

	// 输出序列
	for i, value := range sequence {
		if i > 0 {
			fmt.Print(separator)
		}
		if format != "" {
			fmt.Printf(format, value)
		} else {
			fmt.Print(formatNumber(value))
		}
	}
	fmt.Print(terminator)

	return nil
}

// generateSequence 生成数字序列
//
// 参数:
//   - start: 起始值
//   - end: 结束值
//   - step: 步长
//
// 返回值:
//   - []float64: 生成的序列
func generateSequence(start, end, step float64) []float64 {
	var sequence []float64

	// 判断序列方向
	isForward := step > 0

	// 使用容差处理浮点数比较
	epsilon := 1e-9

	if isForward {
		for current := start; current <= end+epsilon; current += step {
			sequence = append(sequence, current)
			// 防止无限循环（浮点精度问题）
			if current > end+epsilon {
				break
			}
		}
	} else {
		for current := start; current >= end-epsilon; current += step {
			sequence = append(sequence, current)
			// 防止无限循环（浮点精度问题）
			if current < end-epsilon {
				break
			}
		}
	}

	return sequence
}

// calculateEqualWidthFormat 计算等宽输出的格式字符串
//
// 参数:
//   - sequence: 数字序列
//
// 返回值:
//   - string: 格式字符串
func calculateEqualWidthFormat(sequence []float64) string {
	if len(sequence) == 0 {
		return "%g"
	}

	// 检查是否都是整数
	allIntegers := true
	for _, v := range sequence {
		if v != math.Trunc(v) {
			allIntegers = false
			break
		}
	}

	if allIntegers {
		// 计算最大宽度（考虑负号）
		maxWidth := 0
		for _, v := range sequence {
			width := len(fmt.Sprintf("%.0f", v))
			if width > maxWidth {
				maxWidth = width
			}
		}
		return fmt.Sprintf("%%%d.0f", maxWidth)
	}

	// 非整数使用默认格式
	return "%g"
}

// formatNumber 格式化数字，去掉不必要的尾随零
//
// 参数:
//   - num: 要格式化的数字
//
// 返回值:
//   - string: 格式化后的字符串
func formatNumber(num float64) string {
	// 如果是整数，按整数输出
	if num == math.Trunc(num) {
		return fmt.Sprintf("%.0f", num)
	}
	// 否则使用 %g 自动选择格式
	return fmt.Sprintf("%g", num)
}

// unescape 处理转义字符
//
// 参数:
//   - s: 输入字符串
//
// 返回值:
//   - string: 处理后的字符串
func unescape(s string) string {
	replacements := map[string]string{
		`\n`: "\n",
		`\t`: "\t",
		`\r`: "\r",
		`\\`: "\\",
	}

	for old, new := range replacements {
		s = strings.ReplaceAll(s, old, new)
	}

	return s
}

// ParseArgs 解析位置参数
//
// 参数:
//   - args: 位置参数列表
//
// 返回值:
//   - start: 起始值
//   - end: 结束值
//   - step: 步长
//   - error: 解析错误
func ParseArgs(args []string) (start, end, step float64, err error) {
	switch len(args) {
	case 1:
		// seq LAST
		end, err = strconv.ParseFloat(args[0], 64)
		if err != nil {
			return 0, 0, 0, fmt.Errorf("invalid number: %s", args[0])
		}
		start = 1
		step = 1
	case 2:
		// seq FIRST LAST
		start, err = strconv.ParseFloat(args[0], 64)
		if err != nil {
			return 0, 0, 0, fmt.Errorf("invalid number: %s", args[0])
		}
		end, err = strconv.ParseFloat(args[1], 64)
		if err != nil {
			return 0, 0, 0, fmt.Errorf("invalid number: %s", args[1])
		}
		step = 1
	case 3:
		// seq FIRST INCREMENT LAST
		start, err = strconv.ParseFloat(args[0], 64)
		if err != nil {
			return 0, 0, 0, fmt.Errorf("invalid number: %s", args[0])
		}
		step, err = strconv.ParseFloat(args[1], 64)
		if err != nil {
			return 0, 0, 0, fmt.Errorf("invalid number: %s", args[1])
		}
		end, err = strconv.ParseFloat(args[2], 64)
		if err != nil {
			return 0, 0, 0, fmt.Errorf("invalid number: %s", args[2])
		}
	default:
		return 0, 0, 0, fmt.Errorf("requires 1-3 arguments, got %d", len(args))
	}

	return start, end, step, nil
}
