//go:build windows

// Package utils 提供了 Windows 系统特定的文件属性检查功能。
// 该文件实现了 Windows 平台下的隐藏文件检测、只读属性检查等系统相关功能。
package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

	// 检查Unix风格的点文件(排除特殊目录)
	name := filepath.Base(path)

	// 条件: 文件名以点开头，且不是特殊目录 (.) 或 (..)
	return len(name) > 0 && name[0] == '.' && name != "." && name != ".."
}

// IsDriveRoot 检查路径是否是盘符根目录
//
// 支持格式:
//   - D:
//   - D:\
//   - D:/
//   - d:
//   - d:\
//   - d:/
//
// 参数:
//   - path: 路径
//
// 返回:
//   - bool: 是否是盘符根目录
func IsDriveRoot(path string) bool {
	// 统一转换为小写并去除首尾空格
	path = strings.TrimSpace(strings.ToLower(path))

	// 检查路径长度
	if len(path) < 2 || len(path) > 3 {
		return false
	}

	// 检查是否是盘符格式 (X:)
	if path[1] != ':' {
		return false
	}

	// 检查第一个字符是否是字母
	if (path[0] < 'a' || path[0] > 'z') && (path[0] < 'A' || path[0] > 'Z') {
		return false
	}

	// 检查是否有后缀 (\ 或 /)
	if len(path) == 2 {
		// 只有盘符，如 "D:"
		return true
	}

	// 长度为3，检查后缀
	return path[2] == '\\' || path[2] == '/'
}

// IsHiddenWindows 检查Windows文件是否为隐藏
//
// 参数:
//   - path: 文件路径
//
// 返回:
//   - bool: 是否为隐藏文件
//
// 注意:
//   - 如果路径无效或获取属性失败，返回 false
//   - 错误信息不会输出，调用者如需错误处理请使用其他方式
func IsHiddenWindows(path string) bool {
	// 获取文件属性
	utf16Path, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		// 路径转换失败，返回 false
		return false
	}
	attributes, err := syscall.GetFileAttributes(utf16Path)
	if err != nil {
		// 获取属性失败，返回 false
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

// GetFileOwner 获取Windows文件的所有者和组信息
//
// 参数:
//   - filePath: 文件路径
//
// 返回:
//   - string: 文件所有者的用户名
//   - string: 文件所有者的组名
//
// 注意:
//   - 当前使用简化实现，通过 os.Lstat 获取文件信息
//   - 如果获取失败，返回 "?" 作为占位符
//   - Windows 完整的安全描述符需要使用 golang.org/x/sys/windows 包
func GetFileOwner(filePath string) (string, string) {
	// 使用 os.Lstat 获取文件信息
	info, err := os.Lstat(filePath)
	if err != nil {
		return "?", "?"
	}

	// 获取系统特定的文件信息
	if sysInfo, ok := info.Sys().(*syscall.Win32FileAttributeData); ok && sysInfo != nil {
		// 尝试获取文件所有者（通过文件句柄）
		owner, err := getFileOwnerByHandle(filePath)
		if err == nil && owner != "" {
			return owner, "?"
		}
	}

	// 如果无法获取，返回 "?"
	return "?", "?"
}

// getFileOwnerByHandle 通过文件句柄获取所有者
func getFileOwnerByHandle(filePath string) (string, error) {
	// 打开文件获取句柄
	utf16Path, err := syscall.UTF16PtrFromString(filePath)
	if err != nil {
		return "", err
	}

	// 打开文件（只读模式）
	handle, err := syscall.CreateFile(
		utf16Path,
		syscall.GENERIC_READ,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE,
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_ATTRIBUTE_NORMAL,
		0,
	)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = syscall.CloseHandle(handle)
	}()

	// 获取文件信息
	var fileInfo syscall.ByHandleFileInformation
	err = syscall.GetFileInformationByHandle(handle, &fileInfo)
	if err != nil {
		return "", err
	}

	// 返回文件索引作为标识（简化实现）
	_ = fileInfo
	return "", fmt.Errorf("not implemented")
}
