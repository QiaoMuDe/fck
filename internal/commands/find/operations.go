// Package find 实现了文件查找的操作功能。
// 该文件提供了文件操作器，支持删除、移动文件以及执行自定义命令等操作。
package find

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/color"
	"gitee.com/MM-Q/go-kit/fs"
	"gitee.com/MM-Q/shx"
)

// FileOperator 负责所有文件操作：删除、移动、执行命令
type FileOperator struct {
	cl           *color.GlobalColor
	printActions bool
	movePath     string
}

// NewFileOperator 创建新的文件操作器
func NewFileOperator(cl *color.GlobalColor, printActions bool, movePath string) *FileOperator {
	return &FileOperator{
		cl:           cl,
		printActions: printActions,
		movePath:     movePath,
	}
}

// Delete 删除匹配的文件或目录
//
// 参数:
//   - path: 要删除的文件/目录路径
//
// 返回:
//   - error: 错误信息
func (o *FileOperator) Delete(path string) error {
	if o.printActions {
		o.cl.Redf("del: %s\n", path)
	}

	if rmErr := os.RemoveAll(path); rmErr != nil {
		return fmt.Errorf("failed to delete: %s: %v", path, rmErr)
	}

	return nil
}

// Move 移动匹配的文件或目录到指定位置
//
// 参数:
//   - srcPath: 源文件/目录路径
//   - targetPath: 目标路径
//
// 返回:
//   - error: 错误信息
func (o *FileOperator) Move(srcPath, targetPath string) error {
	// 检查源路径是否为空
	if srcPath == "" {
		return fmt.Errorf("source path is empty")
	}

	// 检查目标路径是否为空
	if targetPath == "" {
		return fmt.Errorf("target path is empty")
	}

	// 检查源路径是否存在
	if _, err := os.Lstat(srcPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("source file or directory does not exist: %s", srcPath)
		}
		return fmt.Errorf("failed to check source file or directory: %s: %v", srcPath, err)
	}

	// fs.MoveEx 会自动处理：
	// 1. 空路径检查
	// 2. 获取绝对路径
	// 3. 循环移动检查（防止移动到子目录）
	// 4. 如果目标是已存在的目录，自动追加源文件名
	// 5. 如果目标不存在或不是目录，使用精确路径
	// 6. 自动创建父目录

	// 执行移动操作
	if err := fs.MoveEx(srcPath, targetPath, false, o.printActions); err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("target file already exists: %s -> %s", srcPath, targetPath)
		}
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied, cannot move file: %v", err)
		}
		return fmt.Errorf("failed to move file: %s -> %s: %v", srcPath, targetPath, err)
	}

	return nil
}

// Execute 执行指定的命令，支持直接执行和shell执行两种模式
//
// 参数:
//   - cmdStr: 要执行的命令字符串
//   - path: 文件路径，用于替换命令中的{}占位符
//
// 返回:
//   - error: 错误信息
//
// 执行模式:
//   - 默认直接执行: "cat {}" (更安全，性能更好)
func (o *FileOperator) Execute(cmdStr, path string) error {
	// 检查cmdStr是否为空
	if cmdStr == "" {
		return fmt.Errorf("command is empty")
	}

	// 检查路径是否为空
	if path == "" {
		return fmt.Errorf("path is empty")
	}

	// 检查是否包含{}
	if !strings.Contains(cmdStr, "{}") {
		return fmt.Errorf("command must contain {} as path placeholder when using -exec flag")
	}

	// 检查路径是否存在
	if _, err := os.Lstat(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file or directory does not exist: %s", path)
		}
		return fmt.Errorf("cannot access file or directory: %s", path)
	}

	// 安全地替换{}为实际的文件路径
	// 使用引号包裹路径，防止 Windows 下反斜杠被当作转义字符处理
	cleanPath := filepath.Clean(path)
	cmdStr = strings.ReplaceAll(cmdStr, "{}", fmt.Sprintf("%q", cleanPath))

	// 如果启用了print-actions输出, 打印执行的命令
	if o.printActions {
		o.cl.Redf("exec: %v\n", cmdStr)
	}

	// 执行命令
	if err := shx.RunToTerminal(cmdStr); err != nil {
		return fmt.Errorf("failed to execute command: %v", err)
	}

	return nil
}
