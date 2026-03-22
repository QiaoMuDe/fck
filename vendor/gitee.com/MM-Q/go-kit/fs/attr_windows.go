//go:build windows

package fs

import (
	"path/filepath"
	"syscall"
)

// isHidden 判断Windows文件或目录是否为隐藏
func isHidden(path string) bool {
	name := filepath.Base(path)

	// 检查文件名是否以 "." 开头
	if len(name) > 1 && name[0] == '.' {
		return true
	}

	// 检查Windows隐藏属性
	utf16Path, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return false
	}

	attributes, err := syscall.GetFileAttributes(utf16Path)
	if err != nil {
		return false
	}

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
