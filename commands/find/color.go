// Package find 实现了文件查找结果的彩色输出功能。
// 该文件提供了根据文件类型（目录、可执行文件、符号链接等）进行彩色显示的功能。
package find

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/commands/internal/types"
)

const (
	// 目录颜色常量
	dirColor = "blue"
)

// isEmptyFile 检查DirEntry是否为空文件
//
// 参数:
//   - d: 要检查的DirEntry对象
//
// 返回:
//   - bool: 如果是空文件返回true，否则返回false
func isEmptyFile(d os.DirEntry) bool {
	// 空DirEntry对象直接返回false
	if d == nil {
		return false
	}

	info, err := d.Info()
	// 如果获取文件信息失败或info为nil，返回false
	if err != nil || info == nil {
		return false
	}

	return info.Size() == 0
}

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
//   - cl: colorlib.ColorLib实例，用于彩色输出
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
	var fileColor string
	switch {
	case mode&os.ModeSymlink != 0, isWindowsSymlink(mode, ext):
		fileColor = "cyan" // 符号链接和Windows快捷方式

	case mode.IsDir():
		cl.Blue(path) // 目录直接输出(蓝色)
		return

	case isSpecialDevice(mode):
		fileColor = "yellow" // 各种设备文件

	case isUnixExecutable(mode), isWindowsExecutable(mode, ext):
		fileColor = "green" // 可执行文件

	case isEmptyFile(d):
		fileColor = "gray" // 空文件

	case mode.IsRegular():
		fileColor = "white" // 普通文件

	default:
		fileColor = "white" // 其他类型文件
	}

	// 输出路径
	printColor(path, cl, dirColor, fileColor)
}

// printColor 打印彩色路径
//
// 参数:
//   - p: 文件路径
//   - cl: colorlib.ColorLib实例，用于彩色输出
//   - dirColor: 目录部分的颜色字符串
//   - fileColor: 文件名部分的颜色字符串
func printColor(p string, cl *colorlib.ColorLib, dirColor string, fileColor string) {
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

// getColoredString 根据颜色字符串调用对应的颜色方法
//
// 参数:
//   - cl: colorlib.ColorLib实例，用于彩色输出
//   - color: 颜色字符串，用于指定输出颜色
//   - text: 要输出的文本
//
// 返回:
//   - string: 彩色文本
func getColoredString(cl *colorlib.ColorLib, color string, text string) string {
	switch color {
	case "red":
		return cl.Sred(text)
	case "green":
		return cl.Sgreen(text)
	case "yellow":
		return cl.Syellow(text)
	case "blue":
		return cl.Sblue(text)
	case "cyan":
		return cl.Scyan(text)
	case "white":
		return cl.Swhite(text)
	case "gray":
		return cl.Sgray(text)
	default:
		return text
	}
}
