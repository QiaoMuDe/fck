//go:build windows

package cmd

import (
	"path/filepath"
	"syscall"
)

// isHidden 判断Windows文件或目录是否为隐藏
func isHidden(path string) bool {
	name := filepath.Base(path)
	if len(name) > 2 && name[0] == '.' {
		return true
	}
	return isHiddenWindows(path)
}

// isHiddenWindows 检查Windows文件是否为隐藏
func isHiddenWindows(path string) bool {
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
