// Package shellx 定义了shell命令执行库的核心数据类型。
// 本文件定义了工具函数，提供命令字符串处理和解析功能。
//
// 主要功能：
//   - getCmdStr: 从Builder对象获取完整的命令字符串
//   - ParseCmd: 智能解析命令字符串，支持复杂的引号处理
//   - FindCmd: 查找系统中的命令路径
//
// ParseCmd函数特性：
//   - 支持单引号、双引号、反引号三种引号类型
//   - 正确处理引号内的空格和特殊字符
//   - 支持引号嵌套（不同类型的引号可以嵌套）
//   - 自动检测未闭合的引号并返回空结果
//   - 处理多个连续空格和制表符
//   - 支持复杂的命令行参数解析
//
// 解析示例：
//   - `ls -la` → ["ls", "-la"]
//   - `echo "hello world"` → ["echo", "hello world"]
//   - `git commit -m "fix: update 'config' file"` → ["git", "commit", "-m", "fix: update 'config' file"]
//   - `find . -name "*.go" -exec grep "pattern" {} \;` → ["find", ".", "-name", "*.go", "-exec", "grep", "pattern", "{}", "\\;"]
package shellx

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// ########################################
// 便捷函数
// ########################################

// ExecStr 执行命令(阻塞)
//
// 参数:
//   - cmdStr: 命令字符串
//
// 返回:
//   - error: 错误信息
func ExecStr(cmdStr string) error {
	return NewCmdStr(cmdStr).WithStdout(os.Stdout).WithStderr(os.Stderr).Exec()
}

// Exec 执行命令(阻塞)
//
// 函数:
//   - name: 命令名
//   - args: 命令参数
//
// 返回:
//   - error: 错误信息
func Exec(name string, args ...string) error {
	return NewCmd(name, args...).WithStdout(os.Stdout).WithStderr(os.Stderr).Exec()
}

// ExecOutStr 执行命令并返回合并后的输出(阻塞)
//
// 参数:
//   - cmdStr: 命令字符串
//
// 返回:
//   - []byte: 输出
//   - error: 错误信息
//
// 注意:
//   - 由于需要捕获默认的stdout和stderr合并输出, 内部已经设置了WithStdout(os.Stdout)和WithStderr(os.Stderr)
func ExecOutStr(cmdStr string) ([]byte, error) {
	return NewCmdStr(cmdStr).ExecOutput()
}

// ExecOut 执行命令并返回合并后的输出(阻塞)
//
// 函数:
//   - name: 命令名
//   - args: 命令参数
//
// 返回:
//   - []byte: 输出
//   - error: 错误信息
//
// 注意:
//   - 由于需要捕获默认的stdout和stderr合并输出, 内部已经设置了WithStdout(os.Stdout)和WithStderr(os.Stderr)
func ExecOut(name string, args ...string) ([]byte, error) {
	return NewCmd(name, args...).ExecOutput()
}

// ExecStrT 执行命令(阻塞，带超时)
//
// 参数:
//   - timeout: 超时时间，如果为0则不设置超时
//   - cmdStr: 命令字符串
//
// 返回:
//   - error: 错误信息
func ExecStrT(timeout time.Duration, cmdStr string) error {
	cmd := NewCmdStr(cmdStr).WithStdout(os.Stdout).WithStderr(os.Stderr)

	// 设置超时
	if timeout > 0 {
		cmd = cmd.WithTimeout(timeout)
	}

	return cmd.Exec()
}

// ExecT 执行命令(阻塞，带超时)
//
// 参数:
//   - timeout: 超时时间，如果为0则不设置超时
//   - name: 命令名
//   - args: 命令参数
//
// 返回:
//   - error: 错误信息
func ExecT(timeout time.Duration, name string, args ...string) error {
	cmd := NewCmd(name, args...).WithStdout(os.Stdout).WithStderr(os.Stderr)

	// 设置超时
	if timeout > 0 {
		cmd = cmd.WithTimeout(timeout)
	}

	return cmd.Exec()
}

// ExecOutStrT 执行命令并返回合并后的输出(阻塞，带超时)
//
// 参数:
//   - timeout: 超时时间，如果为0则不设置超时
//   - cmdStr: 命令字符串
//
// 返回:
//   - []byte: 合并后的输出
//   - error: 错误信息
func ExecOutStrT(timeout time.Duration, cmdStr string) ([]byte, error) {
	cmd := NewCmdStr(cmdStr)

	// 设置超时
	if timeout > 0 {
		cmd = cmd.WithTimeout(timeout)
	}

	return cmd.ExecOutput()
}

// ExecOutT 执行命令并返回合并后的输出(阻塞，带超时)
//
// 参数:
//   - timeout: 超时时间，如果为0则不设置超时
//   - name: 命令名
//   - args: 命令参数
//
// 返回:
//   - []byte: 合并后的输出
//   - error: 错误信息
func ExecOutT(timeout time.Duration, name string, args ...string) ([]byte, error) {
	cmd := NewCmd(name, args...)

	// 设置超时
	if timeout > 0 {
		cmd = cmd.WithTimeout(timeout)
	}

	return cmd.ExecOutput()
}

// ParseCmd 将命令字符串解析为命令切片，支持引号处理(单引号、双引号、反引号)，出错时返回空切片
//
// 实现原理：
//  1. 去除首尾空白
//  2. 遍历每个字符
//  3. 处理引号状态切换
//  4. 在非引号状态下遇到空格时分割
//  5. 检查引号是否闭合
//
// 参数:
//   - cmdStr: 要解析的命令字符串
//
// 返回值:
//   - []string: 解析后的命令切片
func ParseCmd(cmdStr string) []string {
	// 去除首尾空白
	cmdStr = strings.TrimSpace(cmdStr)
	if cmdStr == "" {
		return []string{}
	}

	var (
		result    []string // 解析结果
		current   []rune   // 当前命令片段
		inQuotes  bool     // 是否在引号中
		quote     rune     // 当前引号类型
		hadQuotes bool     // 当前片段是否包含过引号
	)

	// 遍历每个字符
	for _, r := range cmdStr {
		if r == '"' || r == '\'' || r == '`' {
			if !inQuotes {
				inQuotes = true // 开始引号
				quote = r
				hadQuotes = true // 标记当前片段包含引号

			} else if r == quote {
				inQuotes = false // 引号闭合

			} else {
				current = append(current, r) // 引号内的字符直接添加
			}

		} else if (r == ' ' || r == '\t') && !inQuotes {
			if len(current) > 0 || hadQuotes {
				result = append(result, string(current)) // 非引号状态下遇到空格或制表符，添加当前命令片段
				current = current[:0]
				hadQuotes = false
			}

		} else {
			current = append(current, r)
		}
	}

	// 添加最后一个命令片段
	if len(current) > 0 || hadQuotes {
		result = append(result, string(current))
	}

	// 检查引号是否闭合
	if inQuotes {
		return []string{}
	}

	return result
}

// FindCmd 查找命令
//
// 参数:
//   - name: 命令名称
//
// 返回:
//   - string: 命令路径
//   - error: 错误信息
func FindCmd(name string) (string, error) {
	return exec.LookPath(name)
}

// validateEnvVar 验证环境变量格式
//
// 参数:
//   - env: 环境变量字符串，格式为 "key=value"
//
// 返回:
//   - error: 错误信息
func validateEnvVar(env string) error {
	if strings.TrimSpace(env) == "" {
		return fmt.Errorf("environment variable cannot be empty")
	}

	parts := strings.SplitN(env, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid environment variable format, expected 'key=value': %s", env)
	}

	key := strings.TrimSpace(parts[0])
	if key == "" {
		return fmt.Errorf("environment variable key cannot be empty: %s", env)
	}

	return nil
}
