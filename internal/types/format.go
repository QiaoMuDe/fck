package types

import "github.com/jedib0t/go-pretty/v6/table"

// 语法高亮相关常量
const (
	// HighlightFormatter256 256色格式化器 (兼容性最好)
	HighlightFormatter256 = "terminal256"
	// HighlightFormatter16m 真彩色格式化器 (最佳效果)
	HighlightFormatter16m = "terminal16m"
	// HighlightFormatter16 16色格式化器 (最低兼容)
	HighlightFormatter16 = "terminal"
	// HighlightStyleDefault 默认高亮主题
	HighlightStyleDefault = "monokai"
)

// Table样式映射表
var TableStyleMap = map[string]table.Style{
	"def":  table.StyleDefault,                    // 默认样式
	"l":    table.StyleLight,                      // 浅色样式
	"r":    table.StyleRounded,                    // 圆角样式
	"bd":   table.StyleBold,                       // 粗体样式
	"cb":   table.StyleColoredBright,              // 彩色亮色样式
	"cd":   table.StyleColoredDark,                // 彩色暗色样式
	"db":   table.StyleDouble,                     // 双线样式
	"cbb":  table.StyleColoredBlackOnBlueWhite,    // 黑色背景蓝色字体样式
	"cbc":  table.StyleColoredBlackOnCyanWhite,    // 青色背景蓝色字体样式
	"cbg":  table.StyleColoredBlackOnGreenWhite,   // 绿色背景蓝色字体样式
	"cbm":  table.StyleColoredBlackOnMagentaWhite, // 紫色背景蓝色字体样式
	"cby":  table.StyleColoredBlackOnYellowWhite,  // 黄色背景蓝色字体样式
	"cbr":  table.StyleColoredBlackOnRedWhite,     // 红色背景蓝色字体样式
	"cwb":  table.StyleColoredBlueWhiteOnBlack,    // 蓝色背景白色字体样式
	"ccw":  table.StyleColoredCyanWhiteOnBlack,    // 青色背景白色字体样式
	"cgw":  table.StyleColoredGreenWhiteOnBlack,   // 绿色背景白色字体样式
	"cmw":  table.StyleColoredMagentaWhiteOnBlack, // 紫色背景白色字体样式
	"crw":  table.StyleColoredRedWhiteOnBlack,     // 红色背景白色字体样式
	"cyw":  table.StyleColoredYellowWhiteOnBlack,  // 黄色背景白色字体样式
	"none": StyleNone,                             // 禁用样式
}

// Table样式切片
var TableStyles = []string{
	"def",  // 默认样式
	"l",    // 浅色样式
	"r",    // 圆角样式
	"bd",   // 粗体样式
	"cb",   // 彩色亮色样式
	"cd",   // 彩色暗色样式
	"db",   // 双线样式
	"cbb",  // 黑色背景蓝色字体样式
	"cbc",  // 青色背景蓝色字体样式
	"cbg",  // 绿色背景蓝色字体样式
	"cbm",  // 紫色背景蓝色字体样式
	"cby",  // 黄色背景蓝色字体样式
	"cbr",  // 红色背景蓝色字体样式
	"cwb",  // 蓝色背景白色字体样式
	"ccw",  // 青色背景白色字体样式
	"cgw",  // 绿色背景白色字体样式
	"cmw",  // 紫色背景白色字体样式
	"crw",  // 红色背景白色字体样式
	"cyw",  // 黄色背景白色字体样式
	"none", // 禁用样式
}

// 定义禁用样式
var StyleNone = table.Style{
	Box: table.BoxStyle{
		PaddingLeft:      " ", // 左边框
		PaddingRight:     " ", // 右边框
		MiddleHorizontal: " ", // 水平线
		MiddleVertical:   " ", // 垂直线
		TopLeft:          " ", // 左上角
		TopRight:         " ", // 右上角
		BottomLeft:       " ", // 左下角
		BottomRight:      " ", // 右下角
	},
}
