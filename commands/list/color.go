// Package list 实现了文件列表显示的颜色输出功能。
// 该文件提供了统一的跨平台颜色方案，根据文件类型和扩展名进行彩色显示。
package list

import (
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/colorlib"
)

// GetColorString 根据文件信息返回带有相应颜色的路径字符串
//
// 参数:
//   - info: 文件信息，包含文件类型和扩展名等信息
//   - path: 要处理的路径字符串
//   - cl: 用于彩色输出的colorlib.ColorLib实例
//
// 返回:
//   - string: 经过颜色处理后的路径字符串
//
// 颜色方案:
//   - 蓝色: 目录
//   - 青色: 符号链接
//   - 绿色: 可执行文件和代码文件
//   - 黄色: 设备文件和配置文件
//   - 红色: 数据文件和压缩包
//   - 紫色: 库文件和编译产物
//   - 灰色: 空文件
//   - 白色: 其他文件
func GetColorString(info FileInfo, path string, cl *colorlib.ColorLib) string {
	// 1. 基础文件类型优先处理（跨平台统一）
	switch info.EntryType {
	case DirType:
		// 目录使用蓝色
		return cl.Sblue(path)
	case SymlinkType:
		// 符号链接使用青色
		return cl.Scyan(path)
	case ExecutableType:
		// 可执行文件使用绿色
		return cl.Sgreen(path)
	case SocketType, PipeType, BlockDeviceType, CharDeviceType:
		// 设备文件使用黄色
		return cl.Syellow(path)
	case EmptyType:
		// 空文件使用灰色
		return cl.Sgray(path)
	case FileType:
		// 2. 普通文件按扩展名分类
		return getFileColorByExtension(info.FileExt, path, cl)
	default:
		// 未知类型使用白色
		return cl.Swhite(path)
	}
}

// getFileColorByExtension 根据文件扩展名返回相应颜色
func getFileColorByExtension(ext, path string, cl *colorlib.ColorLib) string {
	// 处理特殊的macOS系统文件
	base := filepath.Base(path)
	if base == ".DS_Store" || base == ".localized" || strings.HasPrefix(base, "._") {
		return cl.Sgray(path) // macOS系统文件使用灰色
	}

	// 处理特殊的无扩展名配置文件
	if ext == "" && isSpecialConfigFile(base) {
		return cl.Syellow(path)
	}

	// 统一转换为小写进行匹配
	lowerExt := strings.ToLower(ext)

	// 根据扩展名分类着色
	switch {
	case greenExtensions[lowerExt]:
		return cl.Sgreen(path) // 绿色系文件
	case yellowExtensions[lowerExt]:
		return cl.Syellow(path) // 黄色系文件
	case redExtensions[lowerExt]:
		return cl.Sred(path) // 红色系文件
	case magentaExtensions[lowerExt]:
		return cl.Smagenta(path) // 紫色系文件
	default:
		return cl.Swhite(path) // 其他文件使用白色
	}
}

// isSpecialConfigFile 检查是否为特殊的配置文件(无扩展名)
func isSpecialConfigFile(filename string) bool {
	// 转换为小写进行匹配
	lower := strings.ToLower(filename)
	return specialConfigFiles[lower]
}
