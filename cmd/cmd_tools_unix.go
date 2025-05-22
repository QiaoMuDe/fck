//go:build !windows

package cmd

import (
	"os"
	"path/filepath"
)

// isHidden 判断Unix文件或目录是否为隐藏
func isHidden(path string) bool {
	name := filepath.Base(path)
	return len(name) > 2 && name[0] == '.'
}

// isReadOnly 判断Unix文件或目录是否为只读
func isReadOnly(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().Perm()&0222 == 0
}
