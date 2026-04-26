package term

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"
)

// Read 从指定的 reader 读取一行输入
//
// 参数:
//   - reader: 输入源（通常为 os.Stdin）
//   - prompt: 输入提示信息
//
// 返回:
//   - string: 输入的字符串（已去首尾空格）
func Read(reader io.Reader, prompt string) string {
	fmt.Print(prompt)

	bufReader := bufio.NewReader(reader)
	line, err := bufReader.ReadString('\n')
	if err != nil {
		return ""
	}

	return strings.TrimSpace(line)
}

// ReadLine 便捷函数，使用 os.Stdin 读取一行输入
//
// 参数:
//   - prompt: 输入提示信息
//
// 返回:
//   - string: 输入的字符串（已去首尾空格）
func ReadLine(prompt string) string {
	return Read(os.Stdin, prompt)
}

// ReadWithDef 带默认值的输入
//
// 参数:
//   - reader: 输入源（通常为 os.Stdin）
//   - prompt: 输入提示信息
//   - def: 默认值
//
// 返回:
//   - string: 输入的字符串（已去首尾空格），空输入时返回默认值
func ReadWithDef(reader io.Reader, prompt, def string) string {
	s := Read(reader, prompt)
	if s == "" {
		return def
	}
	return s
}

// ReadLineWithDef 便捷函数，使用 os.Stdin 读取带默认值的输入
//
// 参数:
//   - prompt: 输入提示信息
//   - def: 默认值
//
// 返回:
//   - string: 输入的字符串（已去首尾空格），空输入时返回默认值
func ReadLineWithDef(prompt, def string) string {
	return ReadWithDef(os.Stdin, prompt, def)
}

// ReadPassword 读取密码输入，不显示回显
//
// 参数:
//   - prompt: 输入提示信息
//
// 返回:
//   - string: 输入的密码（已去首尾空格）
//   - error: 错误信息
func ReadPassword(prompt string) (string, error) {
	fmt.Print(prompt)

	bytePwd, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(bytePwd)), nil
}

// Confirm 确认框，支持默认值
//
// 参数:
//   - reader: 输入源（通常为 os.Stdin）
//   - prompt: 确认提示信息
//   - defVal: 默认值 (true = 回车默认Yes)
//
// 返回:
//   - bool: 用户确认结果
func Confirm(reader io.Reader, prompt string, defVal bool) bool {
	suffix := " [Y/n] "
	if !defVal {
		suffix = " [y/N] "
	}

	input := Read(reader, prompt+suffix)
	input = strings.ToLower(strings.TrimSpace(input))

	switch input {
	case "y", "yes":
		return true
	case "n", "no":
		return false
	default:
		return defVal
	}
}

// ConfirmLine 便捷函数，使用 os.Stdin 进行确认
//
// 参数:
//   - prompt: 确认提示信息
//   - defVal: 默认值 (true = 回车默认Yes)
//
// 返回:
//   - bool: 用户确认结果
func ConfirmLine(prompt string, defVal bool) bool {
	return Confirm(os.Stdin, prompt, defVal)
}

// ReadInt 从指定的 reader 读取整数输入
//
// 参数:
//   - reader: 输入源（通常为 os.Stdin）
//   - prompt: 输入提示信息
//
// 返回:
//   - int: 输入的整数
//   - error: 输入为空或格式错误时返回错误
func ReadInt(reader io.Reader, prompt string) (int, error) {
	s := Read(reader, prompt)
	if s == "" {
		return 0, fmt.Errorf("输入不能为空")
	}

	result, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("无效的整数: %s", s)
	}

	return result, nil
}

// ReadIntLine 便捷函数，使用 os.Stdin 读取整数输入
//
// 参数:
//   - prompt: 输入提示信息
//
// 返回:
//   - int: 输入的整数
//   - error: 输入为空或格式错误时返回错误
func ReadIntLine(prompt string) (int, error) {
	return ReadInt(os.Stdin, prompt)
}

// ReadIntWithDef 从指定的 reader 读取整数输入，带默认值
//
// 参数:
//   - reader: 输入源（通常为 os.Stdin）
//   - prompt: 输入提示信息
//   - def: 默认值
//
// 返回:
//   - int: 输入的整数（空输入时返回默认值）
//   - error: 格式错误时返回错误
func ReadIntWithDef(reader io.Reader, prompt string, def int) (int, error) {
	s := Read(reader, prompt)
	if s == "" {
		return def, nil
	}

	result, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("无效的整数: %s", s)
	}

	return result, nil
}

// ReadIntLineWithDef 便捷函数，使用 os.Stdin 读取带默认值的整数输入
//
// 参数:
//   - prompt: 输入提示信息
//   - def: 默认值
//
// 返回:
//   - int: 输入的整数（空输入时返回默认值）
//   - error: 格式错误时返回错误
func ReadIntLineWithDef(prompt string, def int) (int, error) {
	return ReadIntWithDef(os.Stdin, prompt, def)
}

// ReadFloat 从指定的 reader 读取浮点数输入
//
// 参数:
//   - reader: 输入源（通常为 os.Stdin）
//   - prompt: 输入提示信息
//
// 返回:
//   - float64: 输入的浮点数
//   - error: 输入为空或格式错误时返回错误
func ReadFloat(reader io.Reader, prompt string) (float64, error) {
	s := Read(reader, prompt)
	if s == "" {
		return 0, fmt.Errorf("输入不能为空")
	}

	result, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("无效的数字: %s", s)
	}

	return result, nil
}

// ReadFloatLine 便捷函数，使用 os.Stdin 读取浮点数输入
//
// 参数:
//   - prompt: 输入提示信息
//
// 返回:
//   - float64: 输入的浮点数
//   - error: 输入为空或格式错误时返回错误
func ReadFloatLine(prompt string) (float64, error) {
	return ReadFloat(os.Stdin, prompt)
}

// ReadFloatWithDef 从指定的 reader 读取浮点数输入，带默认值
//
// 参数:
//   - reader: 输入源（通常为 os.Stdin）
//   - prompt: 输入提示信息
//   - def: 默认值
//
// 返回:
//   - float64: 输入的浮点数（空输入时返回默认值）
//   - error: 格式错误时返回错误
func ReadFloatWithDef(reader io.Reader, prompt string, def float64) (float64, error) {
	s := Read(reader, prompt)
	if s == "" {
		return def, nil
	}

	result, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("无效的数字: %s", s)
	}

	return result, nil
}

// ReadFloatLineWithDef 便捷函数，使用 os.Stdin 读取带默认值的浮点数输入
//
// 参数:
//   - prompt: 输入提示信息
//   - def: 默认值
//
// 返回:
//   - float64: 输入的浮点数（空输入时返回默认值）
//   - error: 格式错误时返回错误
func ReadFloatLineWithDef(prompt string, def float64) (float64, error) {
	return ReadFloatWithDef(os.Stdin, prompt, def)
}
