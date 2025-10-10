// Package find 实现了文件查找结果的彩色输出功能。
// 该文件提供了根据文件类型(目录、可执行文件、符号链接等)进行彩色显示的功能。
package find

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/commands/internal/common"
	"gitee.com/MM-Q/fck/commands/internal/types"
)

const (
	// 目录颜色常量
	dirColor = ColorBlue
)

// ColorType 定义了颜色类型
type ColorType uint8

const (
	ColorRed     ColorType = iota // 红色
	ColorGreen                    // 绿色
	ColorYellow                   // 黄色
	ColorBlue                     // 蓝色
	ColorCyan                     // 青色
	ColorWhite                    // 白色
	ColorGray                     // 灰色
	ColorDefault                  // 默认颜色
)

// isWindowsSymlink 检查是否为Windows快捷方式
func isWindowsSymlink(mode os.FileMode, ext string) bool {
	return runtime.GOOS == "windows" && mode.IsRegular() && types.WindowsSymlinkExts[ext]
}

// isWindowsExecutable 检查是否为Windows可执行文件
func isWindowsExecutable(mode os.FileMode, ext string) bool {
	return runtime.GOOS == "windows" && mode.IsRegular() && types.WindowsExecutableExts[ext]
}

// isUnixExecutable 检查是否为Unix可执行文件
func isUnixExecutable(mode os.FileMode) bool {
	return mode&os.ModeType == 0 && mode&0111 != 0
}

// isSpecialDevice 检查是否为特殊设备文件
func isSpecialDevice(mode os.FileMode) bool {
	return mode&(os.ModeDevice|os.ModeCharDevice|os.ModeNamedPipe|os.ModeSocket) != 0
}

// printPathColor 根据路径类型以不同颜色输出路径字符串
//
// 参数:
//   - path: 要检查的路径，用于获取文件类型信息
//   - cl: colorlib.ColorLib实例, 用于彩色输出
//   - d: 匹配到的DirEntry对象
//
// 注意:
//   - 该函数直接输出到标准输出，不返回值
func printPathColor(path string, cl *colorlib.ColorLib, d os.DirEntry) {
	// 路径为空或DirEntry为空或ColorLib为空时直接返回
	if path == "" || d == nil || cl == nil {
		return
	}

	// 禁用颜色输出时直接输出路径
	if !findCmdColor.Get() {
		fmt.Println(path)
		return
	}

	mode := d.Type()          // 获取文件类型
	ext := filepath.Ext(path) // 缓存扩展名，避免重复计算

	// 确定文件颜色
	var fileColor ColorType
	switch {
	case mode&os.ModeSymlink != 0, isWindowsSymlink(mode, ext):
		fileColor = ColorCyan // 符号链接和Windows快捷方式

	case mode.IsDir():
		cl.Blue(path) // 目录直接输出(蓝色)
		return

	case isSpecialDevice(mode):
		fileColor = ColorYellow // 各种设备文件

	case isUnixExecutable(mode), isWindowsExecutable(mode, ext):
		fileColor = ColorGreen // 可执行文件

	case mode.IsRegular():
		// 普通文件：分段渲染（目录蓝 + 文件名按扩展名着色）并直接打印
		printColorByExtension(path, ext, cl)
		return

	default:
		fileColor = ColorWhite // 其他类型文件
	}

	// 输出路径
	printColor(path, cl, dirColor, fileColor)
}

// printColor 打印彩色路径
//
// 参数:
//   - p: 文件路径
//   - cl: colorlib.ColorLib实例, 用于彩色输出
//   - dirColor: 目录部分的颜色类型
//   - fileColor: 文件名部分的颜色类型
func printColor(p string, cl *colorlib.ColorLib, dirColor ColorType, fileColor ColorType) {
	// 获取路径的目录和文件名
	dir, file := filepath.Split(p)

	// 如果目录为空，则返回文件名
	if dir == "" {
		fmt.Println(getColoredString(cl, fileColor, file))
		return
	}

	// 打印渐进式颜色
	fmt.Printf("%s%s\n", getColoredString(cl, dirColor, dir), getColoredString(cl, fileColor, file))
}

// printColorByExtension 分段渲染：目录部分使用既有 dirColor（蓝色），文件名部分使用按扩展名配色
func printColorByExtension(p string, ext string, cl *colorlib.ColorLib) {
	if p == "" || cl == nil {
		return
	}
	dir, file := filepath.Split(p)
	// 仅对文件名部分着色，目录部分保持蓝色
	coloredFile := common.GetFileColorByExtension(ext, file, cl)
	if dir == "" {
		fmt.Println(coloredFile)
		return
	}
	fmt.Printf("%s%s\n", getColoredString(cl, dirColor, dir), coloredFile)
}

// getColoredString 根据颜色字符串调用对应的颜色方法
//
// 参数:
//   - cl: colorlib.ColorLib实例, 用于彩色输出
//   - color: 颜色类型, 用于指定输出颜色
//   - text: 要输出的文本
//
// 返回:
//   - string: 彩色文本
func getColoredString(cl *colorlib.ColorLib, color ColorType, text string) string {
	switch color {
	case ColorRed:
		return cl.Sred(text)
	case ColorGreen:
		return cl.Sgreen(text)
	case ColorYellow:
		return cl.Syellow(text)
	case ColorBlue:
		return cl.Sblue(text)
	case ColorCyan:
		return cl.Scyan(text)
	case ColorWhite:
		return cl.Swhite(text)
	case ColorGray:
		return cl.Sgray(text)
	default:
		return text
	}
}
