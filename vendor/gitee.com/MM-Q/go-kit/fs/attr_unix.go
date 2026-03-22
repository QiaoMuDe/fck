//go:build !windows

package fs

import (
	"os"
	"path/filepath"
)

// isHidden 判断Unix文件或目录是否为隐藏
func isHidden(path string) bool {
	name := filepath.Base(path)
	// Unix系统中，以 "." 开头的文件为隐藏文件
	return len(name) > 1 && name[0] == '.'
}

// isReadOnly 判断Unix文件或目录是否为只读
func isReadOnly(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	// 检查是否没有写权限（所有者、组、其他用户都没有写权限）
	return info.Mode().Perm()&0222 == 0
}
