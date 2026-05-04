package color

// Attribute 定义单个 SGR（Select Graphic Rendition）代码，
// 用于控制终端文本的显示属性，如颜色、样式等。
type Attribute int

// escape 定义 ANSI 转义序列的起始字符。
const escape = "\x1b"

// 基础文本样式属性。
// 这些属性控制文本的显示样式，如加粗、斜体、下划线等。
const (
	Reset        Attribute = iota // 重置所有属性为默认值
	Bold                          // 加粗文本
	Faint                         // 轻淡/暗淡文本（降低亮度）
	Italic                        // 斜体文本
	Underline                     // 下划线文本
	BlinkSlow                     // 慢速闪烁（每秒少于 150 次）
	BlinkRapid                    // 快速闪烁（每分钟 150 次以上）
	ReverseVideo                  // 反显（前景色和背景色互换）
	Concealed                     // 隐蔽/隐藏文本（与背景色相同）
	CrossedOut                    // 删除线文本
)

// 重置属性。
// 这些属性用于取消对应的文本样式。
const (
	ResetBold       Attribute = iota + 22 // 重置加粗
	ResetItalic                           // 重置斜体
	ResetUnderline                        // 重置下划线
	ResetBlinking                         // 重置闪烁
	_                                     // 保留（无对应重置属性）
	ResetReversed                         // 重置反显
	ResetConcealed                        // 重置蔽蔽
	ResetCrossedOut                       // 重置删除线
)

// mapResetAttributes 映射每个样式属性到其对应的重置属性。
// 用于在取消颜色时正确重置文本样式。
var mapResetAttributes map[Attribute]Attribute = map[Attribute]Attribute{
	Bold:         ResetBold,       // 重置加粗
	Faint:        ResetBold,       // 重置轻淡/暗淡文本
	Italic:       ResetItalic,     // 重置斜体
	Underline:    ResetUnderline,  // 重置下划线
	BlinkSlow:    ResetBlinking,   // 重置慢速闪烁
	BlinkRapid:   ResetBlinking,   // 重置快速闪烁
	ReverseVideo: ResetReversed,   // 重置反显
	Concealed:    ResetConcealed,  // 重置蔽蔽
	CrossedOut:   ResetCrossedOut, // 重置删除线
}

// 前景文本颜色（标准 8 色）。
// 这些颜色在大多数终端中都受支持。
const (
	FgBlack   Attribute = iota + 30 // 黑色前景
	FgRed                           // 红色前景
	FgGreen                         // 绿色前景
	FgYellow                        // 黄色前景
	FgBlue                          // 蓝色前景
	FgMagenta                       // 洋红色/品红色前景
	FgCyan                          // 青色前景
	FgWhite                         // 白色前景

	// FgGray 是 FgHiBlack 的别名，在终端中显示为灰色。
	// 提供这个别名是为了让代码更具可读性。
	FgGray = FgHiBlack

	// foreground 是内部常量，用于 256 色和 24 位真彩色模式。
	// 不应在代码中直接使用。
	foreground
)

// 前景高亮文本颜色（高亮 8 色）。
// 这些是高亮版本的标准前景色，比标准色更明亮。
const (
	FgHiBlack   Attribute = iota + 90 // 高亮黑色前景（灰色）
	FgHiRed                           // 高亮红色前景
	FgHiGreen                         // 高亮绿色前景
	FgHiYellow                        // 高亮黄色前景
	FgHiBlue                          // 高亮蓝色前景
	FgHiMagenta                       // 高亮洋红色前景
	FgHiCyan                          // 高亮青色前景
	FgHiWhite                         // 高亮白色前景
)

// 背景文本颜色（标准 8 色）。
// 这些颜色用于设置文本的背景色。
const (
	BgBlack   Attribute = iota + 40 // 黑色背景
	BgRed                           // 红色背景
	BgGreen                         // 绿色背景
	BgYellow                        // 黄色背景
	BgBlue                          // 蓝色背景
	BgMagenta                       // 洋红色/品红色背景
	BgCyan                          // 青色背景
	BgWhite                         // 白色背景

	// BgGray 是 BgHiBlack 的别名，在终端中显示为灰色背景。
	// 提供这个别名是为了让代码更具可读性。
	BgGray = BgHiBlack

	// background 是内部常量，用于 256 色和 24 位真彩色模式。
	// 不应在代码中直接使用。
	background
)

// 背景高亮文本颜色（高亮 8 色）。
// 这些是高亮版本的背景色，比标准色更明亮。
const (
	BgHiBlack   Attribute = iota + 100 // 高亮黑色背景（灰色）
	BgHiRed                            // 高亮红色背景
	BgHiGreen                          // 高亮绿色背景
	BgHiYellow                         // 高亮黄色背景
	BgHiBlue                           // 高亮蓝色背景
	BgHiMagenta                        // 高亮洋红色背景
	BgHiCyan                           // 高亮青色背景
	BgHiWhite                          // 高亮白色背景
)
