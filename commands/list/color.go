// Package list 实现了文件列表显示的颜色输出功能。
// 该文件提供了统一的跨平台颜色方案，根据文件类型和扩展名进行彩色显示。
package list

import (
	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/commands/internal/common"
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
		return common.GetFileColorByExtension(info.FileExt, path, cl)
	default:
		// 未知类型使用白色
		return cl.Swhite(path)
	}
}
