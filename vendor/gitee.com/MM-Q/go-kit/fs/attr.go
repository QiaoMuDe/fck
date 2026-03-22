package fs

// IsHidden 判断文件或目录是否为隐藏
// 用于跨平台检查文件或目录的隐藏属性
//
// 参数:
//   - path: 文件或目录路径
//
// 返回:
//   - bool: 文件为隐藏返回true，否则返回false
func IsHidden(path string) bool {
	return isHidden(path)
}

// IsReadOnly 判断文件或目录是否为只读
// 用于跨平台检查文件或目录的只读属性
//
// 参数:
//   - path: 文件或目录路径
//
// 返回:
//   - bool: 文件为只读返回true，否则返回false
func IsReadOnly(path string) bool {
	return isReadOnly(path)
}
