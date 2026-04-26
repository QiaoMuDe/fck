package utils

import (
	"gitee.com/MM-Q/go-kit/fs"
)

// IsHidden 判断文件或目录是否为隐藏
//
// 参数:
//   - path: 文件或目录路径
//
// 返回:
//   - bool: 是否为隐藏
func IsHidden(path string) bool {
	return fs.IsHidden(path)
}

// IsReadOnly 判断文件或目录是否为只读
//
// 参数:
//   - path: 文件或目录的路径
//
// 返回:
//   - bool: 文件或目录是否为只读
func IsReadOnly(path string) bool {
	return fs.IsReadOnly(path)
}

// GetFileOwner 获取文件的所属用户和组
//
// 参数:
//   - filePath: 文件路径
//
// 返回:
//   - string: 文件所有者的用户名
//   - string: 文件所有者的组名
func GetFileOwner(filePath string) (string, string) {
	return fs.GetFileOwner(filePath)
}
