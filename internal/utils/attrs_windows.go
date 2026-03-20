//go:build windows

// Package common 提供了 Windows 系统特定的文件属性检查功能。
// 该文件实现了 Windows 平台下的隐藏文件检测、只读属性检查等系统相关功能。
package common

import (
	"fmt"
	"path/filepath"
	"syscall"
)

// IsHidden 判断Windows文件或目录是否为隐藏
//
// 参数:
//   - path: 文件或目录路径
//
// 返回:
//   - bool: 是否为隐藏
func IsHidden(path string) bool {
	// 检查是否是盘符根目录(如 D: 或 D:\)
	if IsDriveRoot(path) {
		return false
	}

	// 优先检查Windows隐藏属性
	if IsHiddenWindows(path) {
		return true
	}

	// 其次检查Unix风格的点文件
	// 检查Unix风格的点文件(排除特殊目录)
	name := filepath.Base(path)
	return len(name) > 0 && name[0] == '.' && name != "." && name != ".."
}

// IsDriveRoot 检查路径是否是盘符根目录
//
// 参数:
//   - path: 路径
//
// 返回:
//   - bool: 是否是盘符根目录
func IsDriveRoot(path string) bool {
	// 检查路径是否是 D: 或 D:\
	if len(path) == 2 && path[1] == ':' {
		return true
	}
	if len(path) == 3 && path[1] == ':' && path[2] == '\\' {
		return true
	}
	return false
}

// IsHiddenWindows 检查Windows文件是否为隐藏
//
// 参数:
//   - path: 文件路径
//
// 返回:
//   - bool: 是否为隐藏文件
func IsHiddenWindows(path string) bool {
	// 获取文件属性
	utf16Path, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		fmt.Printf("Error converting path to UTF16: %v\n", err)
		return false
	}
	attributes, err := syscall.GetFileAttributes(utf16Path)
	if err != nil {
		fmt.Printf("Error getting file attributes for %s: %v\n", path, err)
		return false
	}

	// 检查隐藏属性
	return attributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0
}

// IsReadOnly 判断Windows文件或目录是否为只读
//
// 参数:
//   - path: 文件或目录的路径
//
// 返回:
//   - bool: 文件或目录是否为只读
func IsReadOnly(path string) bool {
	utf16Path, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return false
	}
	attrs, err := syscall.GetFileAttributes(utf16Path)
	if err != nil {
		return false
	}
	return (attrs & syscall.FILE_ATTRIBUTE_READONLY) != 0
}

// GetFileOwner 用于Windows环境下的占位函数
//
// 参数:
//   - filePath - 文件路径
//
// 返回:
//   - string - 文件所有者的用户名
//   - string - 文件所有者的组名
//
// 注意:
//   - 该函数在Windows环境下始终返回?问号。
func GetFileOwner(filePath string) (string, string) {
	_ = filePath
	// 在Windows环境下, 返回"?"占位符
	return "?", "?"
}
