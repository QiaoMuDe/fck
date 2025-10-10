package common

import (
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/colorlib"
)

// GetFileColorByExtension 根据文件扩展名返回相应颜色
//
// 参数:
//   - ext: 文件扩展名
//   - path: 文件路径
//   - cl: colorlib.ColorLib 实例
//
// 返回值:
//   - string: 着色后的文件路径
func GetFileColorByExtension(ext, path string, cl *colorlib.ColorLib) string {
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
