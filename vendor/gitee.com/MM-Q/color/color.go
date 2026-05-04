package color

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Color 定义一个由 SGR 参数定义的自定义颜色对象。
type Color struct {
	params  []Attribute // SGR 参数列表
	noColor *bool       // 是否禁用颜色
}

// New 返回一个新创建的颜色对象。
//
// 参数:
//   - value: 任意数量的 SGR 参数。
//
// 返回值:
//   - *Color: 新创建的颜色对象。
func New(value ...Attribute) *Color {
	c := &Color{
		params: make([]Attribute, 0),
	}

	if noColorIsSet() {
		c.noColor = boolPtr(true)
	}

	c.Add(value...)
	return c
}

// RGB 返回一个新的 24 位 RGB 前景色。
// 参数值会被自动截断到 0-255 范围内。
//
// 参数:
//   - r: 红色分量 (0-255)
//   - g: 绿色分量 (0-255)
//   - b: 蓝色分量 (0-255)
//
// 返回值:
//   - *Color: 配置好的颜色对象
func RGB(r, g, b int) *Color {
	r = clamp255(r)
	g = clamp255(g)
	b = clamp255(b)
	return New(foreground, 2, Attribute(r), Attribute(g), Attribute(b))
}

// BgRGB 返回一个新的 24 位 RGB 背景色。
// 参数值会被自动截断到 0-255 范围内。
//
// 参数:
//   - r: 红色分量 (0-255)
//   - g: 绿色分量 (0-255)
//   - b: 蓝色分量 (0-255)
//
// 返回值:
//   - *Color: 配置好的颜色对象
func BgRGB(r, g, b int) *Color {
	r = clamp255(r)
	g = clamp255(g)
	b = clamp255(b)
	return New(background, 2, Attribute(r), Attribute(g), Attribute(b))
}

// AddRGB 用于链式添加前景 RGB SGR 参数。
// 可以使用任意数量的参数进行组合并创建自定义颜色对象。
//
// 参数:
//   - r: 红色分量 (0-255)
//   - g: 绿色分量 (0-255)
//   - b: 蓝色分量 (0-255)
//
// 返回值:
//   - *Color: 当前颜色对象, 支持链式调用
//
// 示例:
//
//	c.AddRGB(255, 128, 0).Println("橙色文本")
func (c *Color) AddRGB(r, g, b int) *Color {
	c.params = append(c.params, foreground, 2, Attribute(r), Attribute(g), Attribute(b))
	return c
}

// AddBgRGB 用于链式添加背景 RGB SGR 参数。
// 可以使用任意数量的参数进行组合并创建自定义颜色对象。
//
// 参数:
//   - r: 红色分量 (0-255)
//   - g: 绿色分量 (0-255)
//   - b: 蓝色分量 (0-255)
//
// 返回值:
//   - *Color: 当前颜色对象, 支持链式调用
//
// 示例:
//
//	c.AddBgRGB(255, 128, 0).Println("橙色背景")
func (c *Color) AddBgRGB(r, g, b int) *Color {
	c.params = append(c.params, background, 2, Attribute(r), Attribute(g), Attribute(b))
	return c
}

// Set 设置 SGR 序列。
// 如果颜色被禁用, 则不执行任何操作。
//
// 返回值:
//   - *Color: 当前颜色对象, 支持链式调用
func (c *Color) Set() *Color {
	if c.isNoColorSet() {
		return c
	}

	_, _ = fmt.Fprint(Output, c.format())
	return c
}

// unset 重置 SGR 序列。
// 如果颜色被禁用, 则不执行任何操作。
func (c *Color) unset() {
	if c.isNoColorSet() {
		return
	}

	Unset()
}

// SetWriter 用于使用给定的 io.Writer 设置 SGR 序列。
// 这是一个底层函数, 用户应该使用更高级的函数, 如 color.Fprint、color.Print 等。
//
// 参数:
//   - w: 目标写入器
//
// 返回值:
//   - *Color: 当前颜色对象, 支持链式调用
func (c *Color) SetWriter(w io.Writer) *Color {
	_, _ = c.setWriter(w)
	return c
}

// setWriter 向指定的写入器写入 SGR 序列。
//
// 参数:
//   - w: 目标写入器
//
// 返回值:
//   - int: 写入的字节数
//   - error: 写入过程中的错误 (如果有)
func (c *Color) setWriter(w io.Writer) (int, error) {
	if c.isNoColorSet() {
		return 0, nil
	}

	return fmt.Fprint(w, c.format())
}

// UnsetWriter 使用给定的 io.Writer 重置所有转义属性并清除输出。
// 通常在 SetWriter() 之后调用。
//
// 参数:
//   - w: 目标写入器
func (c *Color) UnsetWriter(w io.Writer) {
	_, _ = c.unsetWriter(w)
}

// unsetWriter 向指定的写入器写入重置序列。
//
// 参数:
//   - w: 目标写入器
//
// 返回值:
//   - int: 写入的字节数
//   - error: 写入过程中的错误 (如果有)
func (c *Color) unsetWriter(w io.Writer) (int, error) {
	if c.isNoColorSet() {
		return 0, nil
	}

	return fmt.Fprintf(w, "%s[%dm", escape, Reset)
}

// Add 用于链式添加 SGR 参数。
// 可以使用任意数量的参数进行组合并创建自定义颜色对象。
//
// 参数:
//   - value: 任意数量的 SGR 参数
//
// 返回值:
//   - *Color: 当前颜色对象, 支持链式调用
//
// 示例:
//
//	c.Add(FgRed, Underline).Println("红色下划线文本")
func (c *Color) Add(value ...Attribute) *Color {
	c.params = append(c.params, value...)
	return c
}

// Fprint 使用其操作数的默认格式进行格式化并写入 w。
// 当操作数都不是字符串时, 在它们之间添加空格。
//
// 参数:
//   - w: 目标写入器
//   - a: 要格式化的操作数
//
// 返回值:
//   - int: 写入的字节数
//   - error: 写入过程中的错误 (如果有)
//
// 注意:
//
//	在 Windows 上, 如果 w 是 *os.File 类型, 用户应该用 colorable.NewColorable() 包装 w。
func (c *Color) Fprint(w io.Writer, a ...interface{}) (n int, err error) {
	n, err = c.setWriter(w)
	if err != nil {
		return n, err
	}

	nn, err := fmt.Fprint(w, a...)
	n += nn
	if err != nil {
		return
	}

	nn, err = c.unsetWriter(w)
	n += nn
	return n, err
}

// Print 使用其操作数的默认格式进行格式化并写入标准输出。
// 当操作数都不是字符串时, 在它们之间添加空格。
// 这是用给定颜色包装的标准 fmt.Print() 方法。
//
// 参数:
//   - a: 要格式化的操作数
//
// 返回值:
//   - int: 写入的字节数
//   - error: 写入过程中的错误 (如果有)
func (c *Color) Print(a ...interface{}) (n int, err error) {
	c.Set()
	defer c.unset()

	return fmt.Fprint(Output, a...)
}

// Fprintf 根据格式说明符进行格式化并写入 w。
//
// 参数:
//   - w: 目标写入器
//   - format: 格式字符串
//   - a: 要格式化的操作数
//
// 返回值:
//   - int: 写入的字节数
//   - error: 写入过程中的错误 (如果有)
//
// 注意:
//
//	在 Windows 上, 如果 w 是 *os.File 类型, 用户应该用 colorable.NewColorable() 包装 w。
func (c *Color) Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	n, err = c.setWriter(w)
	if err != nil {
		return n, err
	}

	nn, err := fmt.Fprintf(w, format, a...)
	n += nn
	if err != nil {
		return
	}

	nn, err = c.unsetWriter(w)
	n += nn
	return n, err
}

// Printf 根据格式说明符进行格式化并写入标准输出。
// 这是用给定颜色包装的标准 fmt.Printf() 方法。
//
// 参数:
//   - format: 格式字符串
//   - a: 要格式化的操作数
//
// 返回值:
//   - int: 写入的字节数
//   - error: 写入过程中的错误 (如果有)
func (c *Color) Printf(format string, a ...interface{}) (n int, err error) {
	c.Set()
	defer c.unset()

	return fmt.Fprintf(Output, format, a...)
}

// Fprintln 使用其操作数的默认格式进行格式化并写入 w。
// 操作数之间始终添加空格, 并追加换行符。
//
// 参数:
//   - w: 目标写入器
//   - a: 要格式化的操作数
//
// 返回值:
//   - int: 写入的字节数
//   - error: 写入过程中的错误 (如果有)
//
// 注意:
//
//	在 Windows 上, 如果 w 是 *os.File 类型, 用户应该用 colorable.NewColorable() 包装 w。
func (c *Color) Fprintln(w io.Writer, a ...interface{}) (n int, err error) {
	return fmt.Fprintln(w, c.wrap(sprintln(a...)))
}

// Println 使用其操作数的默认格式进行格式化并写入标准输出。
// 操作数之间始终添加空格, 并追加换行符。
// 这是用给定颜色包装的标准 fmt.Println() 方法。
//
// 参数:
//   - a: 要格式化的操作数
//
// 返回值:
//   - int: 写入的字节数
//   - error: 写入过程中的错误 (如果有)
func (c *Color) Println(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(Output, c.wrap(sprintln(a...)))
}

// Sprint 类似于 Print, 但返回字符串而不是打印它。
//
// 参数:
//   - a: 要格式化的操作数
//
// 返回值:
//   - string: 格式化后的彩色字符串
func (c *Color) Sprint(a ...interface{}) string {
	return c.wrap(fmt.Sprint(a...))
}

// Sprintln 类似于 Println, 但返回字符串而不是打印它。
//
// 参数:
//   - a: 要格式化的操作数
//
// 返回值:
//   - string: 格式化后的彩色字符串 (包含换行符)
func (c *Color) Sprintln(a ...interface{}) string {
	return c.wrap(sprintln(a...)) + "\n"
}

// Sprintf 类似于 Printf, 但返回字符串而不是打印它。
//
// 参数:
//   - format: 格式字符串
//   - a: 要格式化的操作数
//
// 返回值:
//   - string: 格式化后的彩色字符串
func (c *Color) Sprintf(format string, a ...interface{}) string {
	return c.wrap(fmt.Sprintf(format, a...))
}

// FprintFunc 返回一个新函数, 该函数使用 color.Fprint() 将传入的参数打印为彩色。
//
// 返回值:
//   - func(w io.Writer, a ...interface{}): 打印函数
func (c *Color) FprintFunc() func(w io.Writer, a ...interface{}) {
	return func(w io.Writer, a ...interface{}) {
		_, _ = c.Fprint(w, a...)
	}
}

// PrintFunc 返回一个新函数, 该函数使用 color.Print() 将传入的参数打印为彩色。
//
// 返回值:
//   - func(a ...interface{}): 打印函数
func (c *Color) PrintFunc() func(a ...interface{}) {
	return func(a ...interface{}) {
		_, _ = c.Print(a...)
	}
}

// FprintfFunc 返回一个新函数, 该函数使用 color.Fprintf() 将传入的参数打印为彩色。
//
// 返回值:
//   - func(w io.Writer, format string, a ...interface{}): 格式化打印函数
func (c *Color) FprintfFunc() func(w io.Writer, format string, a ...interface{}) {
	return func(w io.Writer, format string, a ...interface{}) {
		_, _ = c.Fprintf(w, format, a...)
	}
}

// PrintfFunc 返回一个新函数, 该函数使用 color.Printf() 将传入的参数打印为彩色。
//
// 返回值:
//   - func(format string, a ...interface{}): 格式化打印函数
func (c *Color) PrintfFunc() func(format string, a ...interface{}) {
	return func(format string, a ...interface{}) {
		_, _ = c.Printf(format, a...)
	}
}

// FprintlnFunc 返回一个新函数, 该函数使用 color.Fprintln() 将传入的参数打印为彩色。
//
// 返回值:
//   - func(w io.Writer, a ...interface{}): 打印函数 (带换行)
func (c *Color) FprintlnFunc() func(w io.Writer, a ...interface{}) {
	return func(w io.Writer, a ...interface{}) {
		_, _ = c.Fprintln(w, a...)
	}
}

// PrintlnFunc 返回一个新函数, 该函数使用 color.Println() 将传入的参数打印为彩色。
//
// 返回值:
//   - func(a ...interface{}): 打印函数 (带换行)
func (c *Color) PrintlnFunc() func(a ...interface{}) {
	return func(a ...interface{}) {
		_, _ = c.Println(a...)
	}
}

// SprintFunc 返回一个新函数, 该函数使用 fmt.Sprint() 为给定参数返回彩色字符串。
// 可用于放入或混合到其他字符串中。
//
// 返回值:
//   - func(a ...interface{}) string: 字符串生成函数
//
// 示例:
//
//	put := New(FgYellow).SprintFunc()
//	fmt.Fprintf(color.Output, "This is a %s", put("warning"))
//
// 注意:
//
//	Windows 用户应将其与 color.Output 一起使用。
func (c *Color) SprintFunc() func(a ...interface{}) string {
	return func(a ...interface{}) string {
		return c.wrap(fmt.Sprint(a...))
	}
}

// SprintfFunc 返回一个新函数, 该函数使用 fmt.Sprintf() 为给定参数返回彩色字符串。
// 可用于放入或混合到其他字符串中。
//
// 返回值:
//   - func(format string, a ...interface{}) string: 格式化字符串生成函数
//
// 注意:
//
//	Windows 用户应将其与 color.Output 一起使用。
func (c *Color) SprintfFunc() func(format string, a ...interface{}) string {
	return func(format string, a ...interface{}) string {
		return c.wrap(fmt.Sprintf(format, a...))
	}
}

// SprintlnFunc 返回一个新函数, 该函数使用 fmt.Sprintln() 为给定参数返回彩色字符串。
// 可用于放入或混合到其他字符串中。
//
// 返回值:
//   - func(a ...interface{}) string: 字符串生成函数 (带换行)
//
// 注意:
//
//	Windows 用户应将其与 color.Output 一起使用。
func (c *Color) SprintlnFunc() func(a ...interface{}) string {
	return func(a ...interface{}) string {
		return c.wrap(sprintln(a...)) + "\n"
	}
}

// sequence 返回格式化的 SGR 序列, 用于插入 "\x1b[...m"。
//
// 返回值:
//   - string: SGR 序列字符串, 例如 "1;36" 表示粗体青色
func (c *Color) sequence() string {
	format := make([]string, len(c.params))
	for i, v := range c.params {
		format[i] = strconv.Itoa(int(v))
	}

	return strings.Join(format, ";")
}

// wrap 用颜色属性包装字符串 s。
//
// 参数:
//   - s: 要包装的原始字符串
//
// 返回值:
//   - string: 包装后的彩色字符串, 可以直接打印
func (c *Color) wrap(s string) string {
	if c.isNoColorSet() {
		return s
	}

	return c.format() + s + c.unformat()
}

// format 返回 SGR 转义序列的开头部分。
//
// 返回值:
//   - string: 转义序列, 例如 "\x1b[31m"
func (c *Color) format() string {
	return fmt.Sprintf("%s[%sm", escape, c.sequence())
}

// unformat 返回 SGR 转义序列的结尾部分 (重置序列) 。
// 对于序列中的每个元素, 使用特定的重置转义, 如果未找到则使用通用重置。
//
// 返回值:
//   - string: 重置转义序列
func (c *Color) unformat() string {
	format := make([]string, len(c.params))
	for i, v := range c.params {
		format[i] = strconv.Itoa(int(Reset))
		ra, ok := mapResetAttributes[v]
		if ok {
			format[i] = strconv.Itoa(int(ra))
		}
	}

	return fmt.Sprintf("%s[%sm", escape, strings.Join(format, ";"))
}

// DisableColor 禁用颜色输出。
// 可用于在不更改任何现有代码的情况下禁用颜色输出, 例如配合 "--no-color" 标志使用。
// 要重新启用, 请使用 EnableColor() 方法。
func (c *Color) DisableColor() {
	c.noColor = boolPtr(true)
}

// EnableColor 启用颜色输出。
// 与 DisableColor() 一起使用。如果颜色未被禁用, 此方法没有副作用。
func (c *Color) EnableColor() {
	c.noColor = boolPtr(false)
}

// isNoColorSet 检查颜色是否被禁用。
// 首先检查用户设置的颜色选项, 如果没有则返回全局选项。
//
// 返回值:
//   - bool: 如果颜色被禁用则返回 true
func (c *Color) isNoColorSet() bool {
	// 首先检查是否有用户设置的选项
	if c.noColor != nil {
		return *c.noColor
	}

	// 如果没有, 则返回全局选项, 默认情况下是禁用的
	return NoColor
}

// Equals 比较两种颜色是否相等。
//
// 参数:
//   - c2: 要比较的另一个颜色对象
//
// 返回值:
//   - bool: 如果两种颜色相等则返回 true
func (c *Color) Equals(c2 *Color) bool {
	if c == nil && c2 == nil {
		return true
	}
	if c == nil || c2 == nil {
		return false
	}

	if len(c.params) != len(c2.params) {
		return false
	}

	counts := make(map[Attribute]int, len(c.params))
	for _, attr := range c.params {
		counts[attr]++
	}

	for _, attr := range c2.params {
		if counts[attr] == 0 {
			return false
		}
		counts[attr]--
	}

	return true
}
