package fs

import "os"

// Exists 检查指定路径的文件或目录是否存在
// 用于验证文件系统中指定路径是否存在，权限错误等异常情况视为不存在
//
// 参数:
//   - path: 要检查的路径
//
// 返回:
//   - bool: 文件或目录存在返回true，否则返回false
func Exists(path string) bool {
	// 使用os.Stat尝试获取文件信息
	_, err := os.Stat(path)

	// 如果没有错误，说明文件/目录存在
	if err == nil {
		return true
	}

	// 如果错误是文件不存在，则返回false
	if os.IsNotExist(err) {
		return false
	}

	// 其他错误情况（如权限问题等）也视为不存在
	// 根据实际需求，也可以选择返回错误
	return false
}

// IsFile 检查指定路径是否为文件
// 用于验证指定路径是否为普通文件
//
// 参数:
//   - path: 要检查的路径
//
// 返回:
//   - bool: 是文件返回true，否则返回false
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}

// IsDir 检查指定路径是否为目录
// 用于验证指定路径是否为目录
//
// 参数:
//   - path: 要检查的路径
//
// 返回:
//   - bool: 是目录返回true，否则返回false
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
