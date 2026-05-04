package color

import (
	"fmt"
	"io"
	"os"
	"sync"
)

//go:generate go run gen/global_gen.go

// GlobalColor 是全局颜色实例的类型
// 包含独立的配置和输出设置，与普通的 Color 实例完全分离
type GlobalColor struct {
	config *StyleConfig // 样式配置
	color  *Color       // 复用的 Color 对象
	mu     sync.RWMutex // 读写锁，保证线程安全
}

// StyleConfig 定义颜色样式配置
type StyleConfig struct {
	NoColor    bool      // 是否禁用颜色
	Bold       bool      // 加粗
	Underline  bool      // 下划线
	Italic     bool      // 斜体
	Blink      bool      // 闪烁
	Faint      bool      // 暗淡
	CrossedOut bool      // 删除线
	Output     io.Writer // 输出目标
}

// Clone 创建样式配置的深拷贝
// 返回一个新的 StyleConfig 实例，复制当前配置的所有字段
// 注意: Output 字段是接口类型，只会复制引用，不会深拷贝底层的写入器
//
// 返回值:
//   - *StyleConfig: 新的样式配置实例
func (s *StyleConfig) Clone() *StyleConfig {
	if s == nil {
		return defaultStyleConfig()
	}
	return &StyleConfig{
		NoColor:    s.NoColor,
		Bold:       s.Bold,
		Underline:  s.Underline,
		Italic:     s.Italic,
		Blink:      s.Blink,
		Faint:      s.Faint,
		CrossedOut: s.CrossedOut,
		Output:     s.Output,
	}
}

// defaultStyleConfig 返回默认样式配置
// 默认启用颜色输出和加粗样式
// NoColor 会根据全局 NoColor 变量自动判断（考虑环境变量和终端检测）
//
// 返回值:
//   - *StyleConfig: 默认样式配置实例
func defaultStyleConfig() *StyleConfig {
	return &StyleConfig{
		NoColor:    NoColor, // 使用全局 NoColor 判断（包含 NO_COLOR 环境变量和终端检测）
		Bold:       true,    // 默认启用加粗
		Underline:  false,
		Italic:     false,
		Blink:      false,
		Faint:      false,
		CrossedOut: false,
		Output:     os.Stdout,
	}
}

// 全局实例
var (
	globalOnce sync.Once
	globalInst *GlobalColor
)

// initGlobal 初始化全局实例
// 使用 sync.Once 确保线程安全的单例模式
func initGlobal() {
	globalOnce.Do(func() {
		globalInst = &GlobalColor{
			config: defaultStyleConfig(),
			color:  New(),
		}
	})
}

// GetGlobal 返回全局颜色实例
// 如果实例未初始化，会自动创建默认实例
//
// 返回值:
//   - *GlobalColor: 全局颜色实例
//
// 示例:
//
//	g := color.GetGlobal()
//	g.Red("红色文字")
func GetGlobal() *GlobalColor {
	initGlobal()
	return globalInst
}

// G 是 GetGlobal 的快捷方式，返回全局颜色实例
// 使用更短的函数名，方便频繁调用
//
// 返回值:
//   - *GlobalColor: 全局颜色实例
//
// 示例:
//
//	c := color.G()
//	c.Red("红色文字")
func G() *GlobalColor {
	return GetGlobal()
}

// ResetGlobal 重置全局实例到默认状态
// 会重新创建实例，丢弃之前的配置
func ResetGlobal() {
	globalOnce = sync.Once{}
	globalInst = nil
	initGlobal()
}

// ===========================================================
// 配置方法
// ===========================================================

// SetConfig 设置样式配置
// 会自动克隆传入的配置，防止外部修改影响全局实例
// 如果传入 nil，则使用默认配置
//
// 参数:
//   - config: 要设置的样式配置
//
// 返回:
//   - *GlobalColor: 当前 GlobalColor 对象，支持链式调用
//
// 示例:
//
//	config := &color.StyleConfig{Bold: true, NoColor: false}
//	color.G().SetConfig(config).Red("粗体红色文字")
func (g *GlobalColor) SetConfig(config *StyleConfig) *GlobalColor {
	g.mu.Lock()
	defer g.mu.Unlock()
	if config == nil {
		g.config = defaultStyleConfig()
	} else {
		g.config = config.Clone()
	}
	return g
}

// GetConfig 获取当前样式配置
// 返回的是配置对象的指针，直接修改会影响全局实例
// 如果需要修改配置，建议使用 GetConfigClone() 或 SetConfig() 方法
//
// 返回值:
//   - *StyleConfig: 当前样式配置
//
// 示例:
//
//	config := color.G().GetConfig()
//	fmt.Println(config.Bold) // 查看是否启用加粗
func (g *GlobalColor) GetConfig() *StyleConfig {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.config
}

// GetConfigClone 获取当前样式配置的克隆
// 返回一个新的配置对象，修改它不会影响全局实例
// 适合用于基于当前配置创建新配置的场景
//
// 返回值:
//   - *StyleConfig: 当前样式配置的克隆
//
// 示例:
//
//	newConfig := color.G().GetConfigClone()
//	newConfig.Bold = false // 不影响全局实例
func (g *GlobalColor) GetConfigClone() *StyleConfig {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.config.Clone()
}

// SetOutput 设置输出目标
//
// 参数:
//   - w: 输出目标，如果为 nil 则忽略
//
// 返回:
//   - *GlobalColor: 当前 GlobalColor 对象，支持链式调用
//
// 示例:
//
//	var buf bytes.Buffer
//	color.G().SetOutput(&buf).Red("输出到缓冲区")
func (g *GlobalColor) SetOutput(w io.Writer) *GlobalColor {
	g.mu.Lock()
	defer g.mu.Unlock()
	if w != nil {
		g.config.Output = w
	}
	return g
}

// SetNoColor 设置是否禁用颜色
//
// 参数:
//   - noColor: true 表示禁用颜色，false 表示启用颜色
//
// 返回:
//   - *GlobalColor: 当前 GlobalColor 对象，支持链式调用
//
// 示例:
//
//	color.G().SetNoColor(true).Red("这段文字不会显示颜色")
func (g *GlobalColor) SetNoColor(noColor bool) *GlobalColor {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.config.NoColor = noColor
	return g
}

// SetBold 设置是否启用加粗
//
// 参数:
//   - bold: true 表示启用加粗，false 表示禁用
//
// 返回:
//   - *GlobalColor: 当前 GlobalColor 对象，支持链式调用
//
// 示例:
//
//	color.G().SetBold(true).Red("粗体红色文字")
func (g *GlobalColor) SetBold(bold bool) *GlobalColor {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.config.Bold = bold
	return g
}

// SetUnderline 设置是否启用下划线
//
// 参数:
//   - underline: true 表示启用下划线，false 表示禁用
//
// 返回:
//   - *GlobalColor: 当前 GlobalColor 对象，支持链式调用
//
// 示例:
//
//	color.G().SetUnderline(true).Blue("带下划线的蓝色文字")
func (g *GlobalColor) SetUnderline(underline bool) *GlobalColor {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.config.Underline = underline
	return g
}

// SetItalic 设置是否启用斜体
//
// 参数:
//   - italic: true 表示启用斜体，false 表示禁用
//
// 返回:
//   - *GlobalColor: 当前 GlobalColor 对象，支持链式调用
//
// 示例:
//
//	color.G().SetItalic(true).Green("斜体绿色文字")
func (g *GlobalColor) SetItalic(italic bool) *GlobalColor {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.config.Italic = italic
	return g
}

// SetBlink 设置是否启用闪烁
//
// 参数:
//   - blink: true 表示启用闪烁，false 表示禁用
//
// 返回:
//   - *GlobalColor: 当前 GlobalColor 对象，支持链式调用
//
// 示例:
//
//	color.G().SetBlink(true).Yellow("闪烁的黄色文字")
func (g *GlobalColor) SetBlink(blink bool) *GlobalColor {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.config.Blink = blink
	return g
}

// SetFaint 设置是否启用暗淡效果
//
// 参数:
//   - faint: true 表示启用暗淡，false 表示禁用
//
// 返回:
//   - *GlobalColor: 当前 GlobalColor 对象，支持链式调用
//
// 示例:
//
//	color.G().SetFaint(true).Cyan("暗淡的青色文字")
func (g *GlobalColor) SetFaint(faint bool) *GlobalColor {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.config.Faint = faint
	return g
}

// SetCrossedOut 设置是否启用删除线
//
// 参数:
//   - crossedOut: true 表示启用删除线，false 表示禁用
//
// 返回:
//   - *GlobalColor: 当前 GlobalColor 对象，支持链式调用
//
// 示例:
//
//	color.G().SetCrossedOut(true).Red("带删除线的红色文字")
func (g *GlobalColor) SetCrossedOut(crossedOut bool) *GlobalColor {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.config.CrossedOut = crossedOut
	return g
}

// ===========================================================
// 内部方法
// ===========================================================

// buildParams 根据前景色和配置构建 SGR 参数列表
// 如果禁用颜色，返回空参数（所有样式效果都被禁用）
//
// 参数:
//   - fgColor: 前景色属性
//
// 返回值:
//   - []Attribute: SGR 参数列表
func (g *GlobalColor) buildParams(fgColor Attribute) []Attribute {
	config := g.config
	params := make([]Attribute, 0, 7) // 预分配容量：1个颜色 + 最多6个样式

	// 如果禁用颜色，返回空参数（所有样式效果都被禁用）
	if config.NoColor {
		return params
	}

	// 添加前景色
	params = append(params, fgColor)

	// 添加样式属性
	if config.Bold {
		params = append(params, Bold)
	}
	if config.Faint {
		params = append(params, Faint)
	}
	if config.Italic {
		params = append(params, Italic)
	}
	if config.Underline {
		params = append(params, Underline)
	}
	if config.Blink {
		params = append(params, BlinkSlow)
	}
	if config.CrossedOut {
		params = append(params, CrossedOut)
	}

	return params
}

// output 获取当前输出目标
// 如果未设置输出目标，返回 os.Stdout
//
// 返回值:
//   - io.Writer: 当前输出目标
func (g *GlobalColor) output() io.Writer {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if g.config.Output != nil {
		return g.config.Output
	}
	return os.Stdout
}

// printColor 设置颜色并打印（内部方法）
//
// 参数:
//   - fgColor: 前景色属性
//   - format: 格式字符串
//   - a: 格式化参数
func (g *GlobalColor) printColor(fgColor Attribute, format string, a ...interface{}) {
	g.mu.Lock()
	g.color.params = g.buildParams(fgColor)
	c := g.color
	noColor := g.config.NoColor
	g.mu.Unlock()

	if c == nil || noColor {
		_, _ = fmt.Fprintf(g.output(), format, a...)
		return
	}
	_, _ = fmt.Fprint(g.output(), c.Sprintf(format, a...))
}

// sprintColor 设置颜色并返回字符串（不换行）
// 如果禁用颜色，直接返回原始字符串
//
// 参数:
//   - fgColor: 前景色属性
//   - format: 格式字符串
//   - a: 格式化参数
//
// 返回值:
//   - string: 带颜色的字符串
func (g *GlobalColor) sprintColor(fgColor Attribute, format string, a ...interface{}) string {
	g.mu.Lock()
	g.color.params = g.buildParams(fgColor)
	c := g.color
	noColor := g.config.NoColor
	g.mu.Unlock()

	if c == nil || noColor {
		if len(a) == 0 {
			return format
		}
		return fmt.Sprintf(format, a...)
	}
	if len(a) == 0 {
		return c.Sprint(format)
	}
	return c.Sprintf(format, a...)
}
