package fs

import (
	"fmt"
	"os"
	"strconv"
)

// OctStrToMode 将4位八进制字符串权限转换为 os.FileMode 类型
// 用于将八进制权限字符串转换为Go标准库的文件权限类型
//
// 参数:
//   - octalStr: 4位八进制字符串，例如 "0755" 或 "0644"
//
// 返回:
//   - os.FileMode: 转换后的文件权限
//   - error: 输入不合法时返回错误
//
// 示例:
//
//	mode, err := OctStrToMode("0755")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(mode) // 输出: -rwxr-xr-x
func OctStrToMode(octalStr string) (os.FileMode, error) {
	// 检查输入是否为4位
	if len(octalStr) != 4 {
		return 0, fmt.Errorf("输入的权限字符串长度必须为4位")
	}

	// 尝试将字符串解析为八进制数
	mode, err := strconv.ParseUint(octalStr, 8, 32)
	if err != nil {
		return 0, fmt.Errorf("输入的权限字符串不是合法的八进制数: %w", err)
	}

	// 将 uint 转换为 os.FileMode
	return os.FileMode(mode), nil
}

// IsOctPerm 检查输入的权限是否是合法的4位八进制数
// 用于验证权限字符串格式是否符合八进制权限规范
//
// 参数:
//   - permission: 输入的权限字符串，例如 "0755" 或 "0644"
//
// 返回:
//   - bool: 权限格式合法返回true，否则返回false
func IsOctPerm(permission string) bool {
	// 检查长度是否为4位
	if len(permission) != 4 {
		return false
	}

	// 尝试将字符串解析为八进制数
	_, err := strconv.ParseInt(permission, 8, 32)
	return err == nil
}
