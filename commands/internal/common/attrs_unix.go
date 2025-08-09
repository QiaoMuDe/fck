//go:build darwin

package common

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/sys/unix"
)

// IsHidden 判断Unix文件或目录是否为隐藏
//
// 参数:
//   - path: 文件或目录路径
//
// 返回:
//   - bool: 是否为隐藏
func IsHidden(path string) bool {
	// 检查Unix风格的点文件(排除特殊目录)
	name := filepath.Base(path)
	return len(name) > 0 && name[0] == '.' && name != "." && name != ".."
}

// IsReadOnly 判断Unix文件或目录是否为只读
//
// 参数:
//   - path: 文件或目录的路径
//
// 返回:
//   - bool: 文件或目录是否为只读
func IsReadOnly(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().Perm()&0222 == 0
}

// GetFileOwner 获取文件的所属用户和组
//
// 参数:
//   - filePath: 文件路径
//
// 返回:
//   - string: 文件所有者的用户名
//   - string: 文件所有者的组名
//
// 注意:
//   - 在 Linux 和 macOS 上返回用户和组名称
//   - 在 Windows 上返回问号 (?)
func GetFileOwner(filePath string) (string, string) {
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
