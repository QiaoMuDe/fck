package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gitee.com/MM-Q/color"
	"gitee.com/MM-Q/fck/internal/types"
)

// GetFileColorByExtension 根据文件扩展名返回相应颜色
//
// 参数:
//   - ext: 文件扩展名
//   - path: 文件路径
//   - cl: color.GlobalColor 实例
//
// 返回值:
//   - string: 着色后的文件路径
func GetFileColorByExtension(ext, path string, cl *color.GlobalColor) string {
	// 非Windows系统，直接使用白色
	if runtime.GOOS != "windows" {
		return cl.SWhite(path)
	}

	// 处理特殊的macOS系统文件
	base := filepath.Base(path)
	if base == ".DS_Store" || base == ".localized" || strings.HasPrefix(base, "._") {
		return cl.SGray(path) // macOS系统文件使用灰色
	}

	// 处理特殊的无扩展名配置文件
	if ext == "" && isSpecialConfigFile(base) {
		return cl.SYellow(path) // 特殊配置文件使用黄色
	}

	// 统一转换为小写进行匹配
	lowerExt := strings.ToLower(ext)

	// 根据扩展名分类着色
	switch {
	case greenExtensions[lowerExt]:
		return cl.SGreen(path) // 绿色系文件
	case yellowExtensions[lowerExt]:
		return cl.SYellow(path) // 黄色系文件
	case redExtensions[lowerExt]:
		return cl.SRed(path) // 红色系文件
	case magentaExtensions[lowerExt]:
		return cl.SMagenta(path) // 紫色系文件
	default:
		return cl.SWhite(path) // 其他文件使用白色
	}
}

// isSpecialConfigFile 检查是否为特殊的配置文件(无扩展名)
func isSpecialConfigFile(filename string) bool {
	// 转换为小写进行匹配
	lower := strings.ToLower(filename)
	return specialConfigFiles[lower]
}

// SprintStringColor 根据路径类型以不同颜色输出字符串
//
// 参数:
//   - p: 要检查的路径，用于获取文件类型信息
//   - s: 要着色的字符串内容
//   - cl: color.GlobalColor实例，用于彩色输出
//
// 返回:
//   - string: 根据路径类型以不同颜色返回的字符串
func SprintStringColor(p string, s string, cl *color.GlobalColor) string {
	// 获取路径信息
	pathInfo, statErr := os.Lstat(p)
	if statErr != nil {
		return cl.SRed(s) // 如果获取路径信息失败, 返回红色输出
	}

	// 根据路径类型设置颜色
	switch mode := pathInfo.Mode(); {
	case mode&os.ModeSymlink != 0:
		// 符号链接 - 使用青色输出
		return cl.SCyan(s)

	case runtime.GOOS == "windows" && mode.IsRegular() && types.WindowsSymlinkExts[filepath.Ext(p)]:
		// Windows下的快捷方式文件 - 使用青色输出
		return cl.SCyan(s)

	case mode.IsDir():
		// 目录 - 使用蓝色输出
		return cl.SBlue(s)

	case mode&os.ModeDevice != 0:
		// 设备文件 - 使用黄色输出
		return cl.SYellow(s)

	case mode&os.ModeNamedPipe != 0:
		// 命名管道 - 使用黄色输出
		return cl.SYellow(s)

	case mode&os.ModeSocket != 0:
		// 套接字文件 - 使用黄色输出
		return cl.SYellow(s)

	case mode&os.ModeType == 0 && mode&os.ModeCharDevice != 0:
		// 字符设备文件 - 使用黄色输出
		return cl.SYellow(s)

	case mode.IsRegular() && pathInfo.Size() == 0:
		// 空文件 - 使用灰色输出
		return cl.SGray(s)

	case mode.IsRegular() && mode&0111 != 0:
		// 可执行文件 - 使用绿色输出
		return cl.SGreen(s)

	case runtime.GOOS == "windows" && mode.IsRegular() && types.WindowsExecutableExts[filepath.Ext(p)]:
		// Windows下的可执行文件 - 使用绿色输出
		return cl.SGreen(s)

	case mode.IsRegular():
		// 普通文件 - 使用白色输出
		return cl.SWhite(s)

	default:
		// 其他类型文件 - 使用白色输出
		return cl.SWhite(s)
	}
}
