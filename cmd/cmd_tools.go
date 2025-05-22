package cmd

import (
	"os"
	"path/filepath"
	"runtime"
	"syscall"
)

// getLast8Chars 函数用于获取输入字符串的最后 8 个字符。
// 如果输入字符串为空，则返回空字符串；
// 如果输入字符串的长度小于等于 8，则返回该字符串本身。
func getLast8Chars(s string) string {
	// 检查输入字符串是否为空，若为空则直接返回空字符串
	if s == "" {
		return ""
	}

	// 检查输入字符串的长度是否小于等于 8，若是则直接返回该字符串本身
	if len(s) <= 8 {
		return s
	}

	// 若输入字符串长度大于 8，则截取并返回其最后 8 个字符
	return s[len(s)-8:]
}

// isHidden 判断文件或目录是否为隐藏
func isHidden(path string) bool {
	// 获取文件名
	name := filepath.Base(path)

	// 检查文件名是否以 "." 开头（适用于所有系统，包括.git目录）
	if len(name) > 2 && name[0] == '.' {
		return true
	}

	// 在 Linux 和 macOS 上，隐藏文件以 "." 开头
	if runtime.GOOS != "windows" {
		return false
	}

	// 在 Windows 上，需要通过文件属性判断
	// 使用 syscall 库来检查文件属性
	return isHiddenWindows(path)
}

// isHiddenWindows 检查 Windows 文件是否为隐藏
func isHiddenWindows(path string) bool {
	// 调用 Windows API 获取文件属性
	// 转换路径为 UTF-16 编码
	utf16Path, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return false
	}

	// 调用 GetFileAttributes 函数
	attributes, err := syscall.GetFileAttributes(utf16Path)
	if err != nil {
		return false
	}

	// 检查隐藏属性
	return attributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0
}

// isReadOnly 判断文件或目录是否为只读
func isReadOnly(path string) bool {
	// 获取文件信息
	info, err := os.Stat(path)
	if err != nil {
		// 如果文件不存在或无法访问，返回 false
		return false
	}

	// 在 Windows 系统上，检查文件属性中的只读标志
	if runtime.GOOS == "windows" {
		// 转换路径为 UTF-16 格式
		utf16Path, err := syscall.UTF16PtrFromString(path)
		if err != nil {
			return false
		}

		// 获取文件属性
		attrs, err := syscall.GetFileAttributes(utf16Path)
		if err != nil {
			return false
		}

		// 检查只读属性
		return (attrs & syscall.FILE_ATTRIBUTE_READONLY) != 0
	}

	// 在 Linux/Unix 系统上，检查文件权限
	// 如果没有写权限（用户、组、其他都没有写权限），则认为是只读
	return info.Mode().Perm()&0222 == 0
}
