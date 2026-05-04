package color

import (
	"fmt"
	"io"
	"os"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
)

var (
	// NoColor 定义输出是否着色。它根据 stdout 的文件描述符是否指向终端
	// 动态设置为 false 或 true。如果设置了 NO_COLOR 环境变量（无论其值是什么），
	// 它也会被设置为 true。这是一个全局选项, 影响所有颜色。
	// 如需对每个颜色块进行更多控制，请单独使用 DisableColor() 方法。
	NoColor = noColorIsSet() || os.Getenv("TERM") == "dumb" || !stdoutIsTerminal()

	// Output 定义打印函数的标准输出。默认使用 stdOut()。
	Output = stdOut()

	// Error 定义打印函数的标准错误输出。默认使用 stdErr()。
	Error = stdErr()
)

// noColorIsSet 检查环境变量 NO_COLOR 是否设置。
//
// 返回值:
//   - bool: 如果 NO_COLOR 设置为非空字符串则返回 true
func noColorIsSet() bool {
	return os.Getenv("NO_COLOR") != ""
}

// stdoutIsTerminal 检查 os.Stdout 是否是终端。
//
// 返回值:
//   - bool: 如果 os.Stdout 是终端则返回 true, 如果为 nil (例如作为 Windows 服务运行时) 则返回 false
func stdoutIsTerminal() bool {
	if os.Stdout == nil {
		return false
	}
	return isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
}

// stdOut 返回用于颜色输出的写入器。
//
// 返回值:
//   - io.Writer: 颜色输出写入器，如果 os.Stdout 为 nil 则返回 io.Discard
func stdOut() io.Writer {
	if os.Stdout == nil {
		return io.Discard
	}
	return colorable.NewColorableStdout()
}

// stdErr 返回用于颜色错误输出的写入器。
//
// 返回值:
//   - io.Writer: 颜色错误输出写入器，如果 os.Stderr 为 nil 则返回 io.Discard
func stdErr() io.Writer {
	if os.Stderr == nil {
		return io.Discard
	}
	return colorable.NewColorableStderr()
}

// Set 立即设置给定的 SGR 参数。
// 将使用给定的 SGR 参数更改输出颜色，直到调用 color.Unset() 为止。
//
// 参数:
//   - p: 任意数量的 SGR 参数
//
// 返回值:
//   - *Color: 配置好的颜色对象
func Set(p ...Attribute) *Color {
	c := New(p...)
	c.Set()
	return c
}

// Unset 重置所有转义属性并清除输出。
// 通常在 Set() 之后调用。
func Unset() {
	if NoColor {
		return
	}

	_, _ = fmt.Fprintf(Output, "%s[%dm", escape, Reset)
}
