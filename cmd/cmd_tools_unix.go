//go:build darwin

package cmd

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/sys/unix"
)

// isHidden 判断Unix文件或目录是否为隐藏
func isHidden(path string) bool {
	// 检查Unix风格的点文件(排除特殊目录)
	name := filepath.Base(path)
	return len(name) > 0 && name[0] == '.' && name != "." && name != ".."
}

// isReadOnly 判断Unix文件或目录是否为只读
func isReadOnly(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().Perm()&0222 == 0
}

// getFileOwner 获取文件的所属用户和组
// 在 Linux 和 macOS 上返回用户和组名称
// 在 Windows 上返回问号 (?)
func getFileOwner(filePath string) (string, string) {
	// 使用 unix.Stat 获取文件状态
	var stat unix.Stat_t
	if err := unix.Stat(filePath, &stat); err != nil {
		return "?", "?"
	}

	// 获取 UID 和 GID
	uid := stat.Uid
	gid := stat.Gid

	// 获取用户信息
	userInfo, err := user.LookupId(fmt.Sprintf("%d", uid))
	if err != nil {
		return "?", "?"
	}

	// 获取组信息
	groupInfo, err := user.LookupGroupId(fmt.Sprintf("%d", gid))
	if err != nil {
		return "?", "?"
	}

	return userInfo.Username, groupInfo.Name
}
