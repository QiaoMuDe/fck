package shellx

import (
	"os"
	"os/exec"
	"time"
)

// Split 将命令字符串拆分为命令切片，支持引号处理(单引号、双引号、反引号)
//
// 功能：
//   - 智能拆分 Shell 命令字符串为参数数组
//   - 支持单引号、双引号、反引号包裹的内容
//   - 正确处理转义字符和特殊字符
//   - 自动处理命令分隔符(;|&&||)
//
// 参数:
//   - cmdStr: 要拆分的命令字符串
//
// 返回值:
//   - []string: 拆分后的命令切片 (最佳结果)
//
// 注意：
//   - 此函数忽略拆分错误，返回最佳拆分结果。如需错误信息，请使用 SplitE 函数。
//   - 转义字符保持原样，不进行解释处理
//   - 支持多字符操作符如 &&、||、>>、<< 等
func Split(cmdStr string) []string {
	result, _ := splitInternal(cmdStr)
	return result
}

// SplitE 将命令字符串拆分为命令切片（带错误信息），支持引号处理(单引号、双引号、反引号)
//
// 功能：
//   - 智能拆分 Shell 命令字符串为参数数组
//   - 支持单引号、双引号、反引号包裹的内容
//   - 正确处理转义字符和特殊字符
//   - 自动处理命令分隔符(;|&&||)
//   - 检测并返回拆分过程中的错误
//
// 参数:
//   - cmdStr: 要拆分的命令字符串
//
// 返回值:
//   - []string: 拆分后的命令切片
//   - error: 拆分错误，成功时为 nil
//
// 错误类型：
//   - UnclosedQuoteError: 未闭合的引号错误
//   - 其他可能的语法错误
//
// 注意：
//   - 转义字符保持原样，不进行解释处理
//   - 支持多字符操作符如 &&、||、>>、<< 等
func SplitE(cmdStr string) ([]string, error) {
	return splitInternal(cmdStr)
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

// ExecCode 执行命令并返回退出码(阻塞)
//
// 参数:
//   - name: 命令名
//   - args: 命令参数
//
// 返回:
//   - int: 退出码
//   - error: 错误信息
func ExecCode(name string, args ...string) (int, error) {
	err := NewCmd(name, args...).Exec()
	if err != nil {
		return -1, err
	}
	return extractExitCode(err), nil
}

// ExecCodeStr 字符串方式执行命令并返回退出码(阻塞)
//
// 参数:
//   - cmdStr: 命令字符串
//
// 返回:
//   - int: 退出码
//   - error: 错误信息
func ExecCodeStr(cmdStr string) (int, error) {
	err := NewCmdStr(cmdStr).Exec()
	if err != nil {
		return -1, err
	}
	return extractExitCode(err), nil
}
