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

// printPathColor 根据路径类型以不同颜色输出路径字符串
//
// 参数:
//   - path: 要检查的路径，用于获取文件类型信息
//   - cl: colorlib.ColorLib实例，用于彩色输出
//
// 注意:
//   - 该函数直接输出到标准输出，不返回值
func printPathColor(path string, cl *colorlib.ColorLib) {
	// 获取路径信息
	pathInfo, statErr := os.Lstat(path)
	if statErr != nil {
		// 如果获取路径信息失败, 直接输出红色的路径
		printColor(path, cl, "cyan", "red")
		return
	}

	// 根据路径类型设置颜色并直接输出
	switch mode := pathInfo.Mode(); {
	case mode&os.ModeSymlink != 0:
		// 符号链接 - 使用青色输出
		printColor(path, cl, "cyan", "cyan")
	case runtime.GOOS == "windows" && mode.IsRegular() && types.WindowsSymlinkExts[filepath.Ext(path)]:
		// Windows快捷方式 - 使用青色输出
		printColor(path, cl, "cyan", "cyan")
	case mode.IsDir():
		// 目录 - 使用蓝色输出
		printColor(path, cl, "cyan", "blue")
	case mode&os.ModeDevice != 0:
		// 设备文件 - 使用黄色输出
		printColor(path, cl, "cyan", "yellow")
	case mode&os.ModeCharDevice != 0:
		// 字符设备文件 - 使用黄色输出
		printColor(path, cl, "cyan", "yellow")
	case mode&os.ModeNamedPipe != 0:
		// 命名管道 - 使用黄色输出
		printColor(path, cl, "cyan", "yellow")
	case mode&os.ModeSocket != 0:
		// 套接字文件 - 使用黄色输出
		printColor(path, cl, "cyan", "yellow")
	case mode&os.ModeType == 0 && mode&0111 != 0:
		// 可执行文件 - 使用绿色输出
		printColor(path, cl, "cyan", "green")
	case runtime.GOOS == "windows" && mode.IsRegular() && types.WindowsExecutableExts[filepath.Ext(path)]:
		// Windows可执行文件 - 使用绿色输出
		printColor(path, cl, "cyan", "green")
	case pathInfo.Size() == 0:
		// 空文件 - 使用灰色输出
		printColor(path, cl, "cyan", "gray")
	case mode.IsRegular():
		// 普通文件 - 使用白色输出
		printColor(path, cl, "cyan", "white")
	default:
		// 其他类型文件 - 使用白色输出
		printColor(path, cl, "cyan", "white")
	}
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
