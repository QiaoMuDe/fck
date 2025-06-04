//go:build windows

package cmd

import (
	"fmt"
	"path/filepath"
	"syscall"
)

// isHidden 判断Windows文件或目录是否为隐藏
func isHidden(path string) bool {
	// 检查是否是盘符根目录(如 D: 或 D:\)
	if isDriveRoot(path) {
		return false
	}

	name := filepath.Base(path)
	if len(name) > 2 && name[0] == '.' {
		return true
	}

	return isHiddenWindows(path)
}

// isDriveRoot 检查路径是否是盘符根目录
func isDriveRoot(path string) bool {
	// 检查路径是否是 D: 或 D:\
	if len(path) == 2 && path[1] == ':' {
		return true
	}
	if len(path) == 3 && path[1] == ':' && path[2] == '\\' {
		return true
	}
	return false
}

// isHiddenWindows 检查Windows文件是否为隐藏
func isHiddenWindows(path string) bool {
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

// isReadOnly 判断Windows文件或目录是否为只读
func isReadOnly(path string) bool {
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

// getFileOwner 用于Windows环境下的占位函数
// 参数: filePath - 文件路径
// 返回: 用户和组的信息, Windows环境下返回"?"占位符
func getFileOwner(filePath string) (string, string) {
	_ = filePath
	// 在Windows环境下, 返回"?"占位符
	return "?", "?"
}
