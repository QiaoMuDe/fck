package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/globals"
)

// 定义全局常量的颜色映射
var PermissionColorMap = map[int]string{
	1: "green",  // 所有者-读-绿色
	2: "yellow", // 所有者-写-黄色
	3: "red",    // 所有者-执行-红色
	4: "green",  // 组-读-绿色
	5: "yellow", // 组-写-黄色
	6: "red",    // 组-执行-红色
	7: "green",  // 其他-读-绿色
	8: "yellow", // 其他-写-黄色
	9: "red",    // 其他-执行-红色
}

// splitPathColor 函数用于根据路径类型以不同颜色返回字符串
func splitPathColor(p string, cl *colorlib.ColorLib, dirCode int, fileCode int) string {
	// 获取路径的目录和文件名
	dir, file := filepath.Split(p)

	// 如果目录为空，则返回文件名
	if dir == "" {
		return fmt.Sprint(cl.SColor(fileCode, file))
	}

	// 设置渐进式颜色
	return fmt.Sprint(cl.SColor(dirCode, dir), cl.SColor(fileCode, file))
}

// printSizeColor 根据路径类型以不同颜色输出字符串
// 参数:
//
//	path: 路径或文件名
//	s: 文件大小
//	cl: colorlib.ColorLib实例, 用于彩色输出
func printSizeColor(p string, s int64, cl *colorlib.ColorLib) {
	// 第一列固定为白色显示大小
	sizeStr := cl.Swhite(humanReadableSize(s, 2))

	// 获取路径信息
	pathInfo, statErr := os.Lstat(p)
	if statErr != nil {
		if *sizeCmdColor {
			fmt.Printf("%-25s\t%s\n", sizeStr, cl.Sred(p))
		} else {
			fmt.Printf("%-15s\t%s\n", sizeStr, p)
		}
		return
	}

	// 根据路径类型设置第二列颜色
	var pathStr string
	switch mode := pathInfo.Mode(); {
	case mode&os.ModeSymlink != 0: // 符号链接 - 青色
		pathStr = cl.Scyan(p)
	case runtime.GOOS == "windows" && mode.IsRegular() && (filepath.Ext(p) == ".lnk" || filepath.Ext(p) == ".url"): // Windows快捷方式 - 青色
		pathStr = cl.Scyan(p)
	case mode.IsDir(): // 目录 - 蓝色
		pathStr = cl.Sblue(p)
	case mode&os.ModeDevice != 0: // 设备文件 - 黄色
		pathStr = cl.Syellow(p)
	case mode&os.ModeNamedPipe != 0: // 命名管道 - 黄色
		pathStr = cl.Syellow(p)
	case mode&os.ModeSocket != 0: // 套接字 - 黄色
		pathStr = cl.Syellow(p)
	case mode&os.ModeType == 0 && mode&os.ModeCharDevice != 0: // 字符设备 - 黄色
		pathStr = cl.Syellow(p)
	case mode.IsRegular() && pathInfo.Size() == 0: // 空文件 - 灰色
		pathStr = cl.Sgray(p)
	case mode.IsRegular() && mode&0111 != 0: // 可执行文件 - 绿色
		pathStr = cl.Sgreen(p)
	case runtime.GOOS == "windows" && mode.IsRegular() && (filepath.Ext(p) == ".exe" || filepath.Ext(p) == ".bat" || filepath.Ext(p) == ".cmd" || filepath.Ext(p) == ".msi" || filepath.Ext(p) == ".ps1" || filepath.Ext(p) == ".psm1"): // Windows可执行文件 - 绿色
		pathStr = cl.Sgreen(p)
	default: // 其他文件 - 白色
		pathStr = cl.Swhite(p)
	}

	// 格式化输出
	if *sizeCmdColor {
		fmt.Printf("%-25s\t%s\n", sizeStr, pathStr)
	} else {
		fmt.Printf("%-15s\t%s\n", sizeStr, p)
	}
}

// SPrintStringColor 根据路径类型以不同颜色输出字符串
// 参数:
//
//	p: 要检查的路径(用于获取文件类型信息)
//	s: 返回字符串内容
//	cl: colorlib.ColorLib实例, 用于彩色输出
func SprintStringColor(p string, s string, cl *colorlib.ColorLib) string {
	// 获取路径信息
	pathInfo, statErr := os.Lstat(p)
	if statErr != nil {
		return cl.Sred(s) // 如果获取路径信息失败, 返回红色输出
	}

	// 根据路径类型设置颜色
	switch mode := pathInfo.Mode(); {
	case mode&os.ModeSymlink != 0:
		// 符号链接 - 使用青色输出
		return cl.Scyan(s)
	case runtime.GOOS == "windows" && mode.IsRegular() && (filepath.Ext(p) == ".lnk" || filepath.Ext(p) == ".url"):
		// Windows下的快捷方式文件 - 使用青色输出
		return cl.Scyan(s)
	case mode.IsDir():
		// 目录 - 使用蓝色输出
		return cl.Sblue(s)
	case mode&os.ModeDevice != 0:
		// 设备文件 - 使用黄色输出
		return cl.Syellow(s)
	case mode&os.ModeNamedPipe != 0:
		// 命名管道 - 使用黄色输出
		return cl.Syellow(s)
	case mode&os.ModeSocket != 0:
		// 套接字文件 - 使用黄色输出
		return cl.Syellow(s)
	case mode&os.ModeType == 0 && mode&os.ModeCharDevice != 0:
		// 字符设备文件 - 使用黄色输出
		return cl.Syellow(s)
	case mode.IsRegular() && pathInfo.Size() == 0:
		// 空文件 - 使用灰色输出
		return cl.Sgray(s)
	case mode.IsRegular() && mode&0111 != 0:
		// 可执行文件 - 使用绿色输出
		return cl.Sgreen(s)
	case runtime.GOOS == "windows" && mode.IsRegular() && (filepath.Ext(p) == ".exe" || filepath.Ext(p) == ".bat" || filepath.Ext(p) == ".cmd" || filepath.Ext(p) == ".msi" || filepath.Ext(p) == ".ps1" || filepath.Ext(p) == ".psm1"):
		// Windows下的可执行文件 - 使用绿色输出
		return cl.Sgreen(s)
	case mode.IsRegular():
		// 普通文件 - 使用白色输出
		return cl.Swhite(s)
	default:
		// 其他类型文件 - 使用白色输出
		return cl.Swhite(s)
	}
}

// printPathColor 根据路径类型以不同颜色输出路径字符串
// 参数:
//
//	path: 要检查的路径(用于获取文件类型信息)
//	cl: colorlib.ColorLib实例, 用于彩色输出
func printPathColor(path string, cl *colorlib.ColorLib) {
	// 获取路径信息
	pathInfo, statErr := os.Lstat(path)
	if statErr != nil {
		// 如果获取路径信息失败, 输出红色的路径
		fmt.Println(splitPathColor(path, cl, colorlib.Cyan, colorlib.Red))
		return
	}

	// 根据路径类型设置颜色
	switch mode := pathInfo.Mode(); {
	case mode&os.ModeSymlink != 0:
		// 符号链接 - 使用青色输出
		fmt.Println(splitPathColor(path, cl, colorlib.Cyan, colorlib.Cyan))
	case runtime.GOOS == "windows" && mode.IsRegular() && (filepath.Ext(path) == ".lnk" || filepath.Ext(path) == ".url"):
		// Windows快捷方式 - 使用青色输出
		fmt.Println(splitPathColor(path, cl, colorlib.Cyan, colorlib.Cyan))
	case mode.IsDir():
		// 目录 - 使用蓝色输出
		fmt.Println(splitPathColor(path, cl, colorlib.Cyan, colorlib.Blue))
	case mode&os.ModeDevice != 0:
		// 设备文件 - 使用黄色输出
		fmt.Println(splitPathColor(path, cl, colorlib.Cyan, colorlib.Yellow))
	case mode&os.ModeCharDevice != 0:
		// 字符设备文件 - 使用黄色输出
		fmt.Println(splitPathColor(path, cl, colorlib.Cyan, colorlib.Yellow))
	case mode&os.ModeNamedPipe != 0:
		// 命名管道 - 使用黄色输出
		fmt.Println(splitPathColor(path, cl, colorlib.Cyan, colorlib.Yellow))
	case mode&os.ModeSocket != 0:
		// 套接字文件 - 使用黄色输出
		fmt.Println(splitPathColor(path, cl, colorlib.Cyan, colorlib.Yellow))
	case mode&os.ModeType == 0 && mode&0111 != 0:
		// 可执行文件 - 使用绿色输出
		fmt.Println(splitPathColor(path, cl, colorlib.Cyan, colorlib.Green))
	case runtime.GOOS == "windows" && mode.IsRegular() && (filepath.Ext(path) == ".exe" || filepath.Ext(path) == ".bat" || filepath.Ext(path) == ".cmd" || filepath.Ext(path) == ".msi" || filepath.Ext(path) == ".ps1" || filepath.Ext(path) == ".psm1"):
		// Windows可执行文件 - 使用绿色输出
		fmt.Println(splitPathColor(path, cl, colorlib.Cyan, colorlib.Green))
	case pathInfo.Size() == 0:
		// 空文件 - 使用灰色输出
		fmt.Println(splitPathColor(path, cl, colorlib.Cyan, colorlib.Gray))
	case mode.IsRegular():
		// 普通文件 - 使用白色输出
		fmt.Println(splitPathColor(path, cl, colorlib.Cyan, colorlib.White))
	default:
		// 其他类型文件 - 使用白色输出
		fmt.Println(splitPathColor(path, cl, colorlib.Cyan, colorlib.White))
	}
}

// getColorString 函数的作用是根据传入的文件信息、路径字符串以及颜色库实例，返回带有相应颜色的路径字符串。
// 参数:
// info: 包含文件类型和文件后缀名等信息的 globals.ListInfo 结构体实例。
// pF: 要处理的路径字符串。
// cl: 用于彩色输出的 colorlib.ColorLib 实例。
// 返回值:
// colorString: 经过颜色处理后的路径字符串。
func getColorString(info globals.ListInfo, pF string, cl *colorlib.ColorLib) string {
	// 依据文件的类型来确定输出的颜色
	switch info.EntryType {
	case globals.SymlinkType:
		// 若文件类型为符号链接，则使用青色来渲染字符串
		return cl.Scyan(pF)
	case globals.DirType:
		// 若文件类型为目录，则使用蓝色来渲染字符串
		return cl.Sblue(pF)
	case globals.ExecutableType:
		// 若文件类型为可执行文件，则使用绿色来渲染字符串
		return cl.Sgreen(pF)
	case globals.SocketType, globals.PipeType, globals.BlockDeviceType, globals.CharDeviceType:
		// 若文件类型为套接字、管道、块设备、字符设备，则使用黄色来渲染字符串
		return cl.Syellow(pF)
	case globals.EmptyType:
		// 若文件类型为空文件, 则使用灰色来渲染字符串
		return cl.Sgray(pF)
	case globals.FileType:
		// 若文件类型为普通文件，则根据平台差异来设置颜色
		if runtime.GOOS == "windows" {
			switch info.FileExt {
			case "exe", "bat", "cmd", "ps1", "psm1", "msi":
				// 对于 Windows 系统下的可执行文件，使用绿色来渲染字符串
				return cl.Sgreen(pF)
			case "lnk", "url":
				// 对于 Windows 系统下的符号链接，使用青色来渲染字符串
				return cl.Scyan(pF)
			default:
				// 对于其他文件类型，使用白色来渲染字符串
				return cl.Swhite(pF)
			}
		}

		// 添加MacOS特殊文件处理
		if runtime.GOOS == "darwin" {
			base := filepath.Base(pF)
			switch {
			case base == ".DS_Store" || base == ".localized" || strings.HasPrefix(base, "._"):
				return cl.Sgray(pF) // MacOS系统文件使用灰色
			case filepath.Ext(pF) == ".app":
				return cl.Sgreen(pF) // MacOS应用程序包使用绿色
			}
		}

		// 对于 Linux 系统下的普通文件，使用白色来渲染字符串
		return cl.Swhite(pF)
	default:
		// 对于未匹配的类型，使用白色来渲染字符串
		return cl.Swhite(pF)
	}
}
