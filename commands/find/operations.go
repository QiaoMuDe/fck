// Package find 实现了文件查找的操作功能。
// 该文件提供了文件操作器，支持删除、移动文件以及执行自定义命令等操作。
package find

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/shellx"
)

// FileOperator 负责所有文件操作：删除、移动、执行命令
type FileOperator struct {
	cl *colorlib.ColorLib
}

// NewFileOperator 创建新的文件操作器
func NewFileOperator(cl *colorlib.ColorLib) *FileOperator {
	return &FileOperator{
		cl: cl,
	}
}

// Delete 删除匹配的文件或目录
//
// 参数:
//   - path: 要删除的文件/目录
//
// 返回:
//   - error: 错误信息
func (o *FileOperator) Delete(path string) error {
	// 打印删除信息
	if findCmdPrintActions.Get() {
		o.cl.Redf("del: %s\n", path)
	}

	// 调用删除函数统一处理
	if rmErr := os.RemoveAll(path); rmErr != nil {
		return fmt.Errorf("删除失败: %s: %v", path, rmErr)
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
		return fmt.Errorf("源路径为空")
	}

	// 检查目标路径是否为空
	if targetPath == "" {
		return fmt.Errorf("没有指定目标路径")
	}

	// 检查源路径是否存在
	if _, err := os.Lstat(srcPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("源文件/目录不存在: %s", srcPath)
		}
		return fmt.Errorf("检查源文件/目录时出错: %s: %v", srcPath, err)
	}

	// 获取目标路径的绝对路径
	absTargetPath, err := filepath.Abs(targetPath)
	if err != nil {
		return fmt.Errorf("获取目标路径绝对路径失败: %v", err)
	}

	// 获取源路径的绝对路径
	absSearchPath, err := filepath.Abs(srcPath)
	if err != nil {
		return fmt.Errorf("获取源路径绝对路径失败: %v", err)
	}

	// 检查目标路径是否是源路径的子目录(防止循环移动)
	if strings.HasPrefix(absTargetPath, absSearchPath) {
		return fmt.Errorf("不能将目录移动到自身或其子目录中")
	}

	// 确保目标目录存在
	if err := os.MkdirAll(filepath.Dir(absTargetPath), 0755); err != nil {
		return fmt.Errorf("创建目标目录失败: %v", err)
	}

	// 组装完整的目标路径: 目标路径 + 源路径的文件名
	if filepath.Base(absSearchPath) != "" {
		absTargetPath = filepath.Join(absTargetPath, filepath.Base(absSearchPath))
	}

	// 检查目标文件是否已存在
	if _, err := os.Stat(absTargetPath); err == nil {
		// 如果是移动操作, 直接跳过而不是报错
		if findCmdMove.Get() != "" {
			return nil
		}
		return fmt.Errorf("目标已存在: %s", absTargetPath)
	}

	// 打印移动信息
	if findCmdPrintActions.Get() {
		o.cl.Redf("mv: %s -> %s\n", absSearchPath, absTargetPath)
	}

	// 执行移动操作
	if err := os.Rename(absSearchPath, absTargetPath); err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("目标文件已存在: %s -> %s", absSearchPath, absTargetPath)
		}
		if os.IsPermission(err) {
			return fmt.Errorf("权限不足, 无法移动文件: %v", err)
		}
		return fmt.Errorf("移动失败: %s -> %s: %v", absSearchPath, absTargetPath, err)
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
//   - 使用--use-shell/-us启用shell执行: 支持管道、重定向等shell功能
func (o *FileOperator) Execute(cmdStr, path string) error {
	// 检查cmdStr是否为空
	if cmdStr == "" {
		return fmt.Errorf("命令为空")
	}

	// 检查路径是否为空
	if path == "" {
		return fmt.Errorf("路径为空")
	}

	// 检查是否包含{}
	if !strings.Contains(cmdStr, "{}") {
		return fmt.Errorf("使用-exec标志时必须包含{}作为路径占位符")
	}

	// 检查路径是否存在
	if _, err := os.Lstat(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("文件/目录不存在: %s", path)
		}
		return fmt.Errorf("无法访问文件/目录: %s", path)
	}

	// 安全地替换{}为实际的文件路径
	cmdStr = strings.ReplaceAll(cmdStr, "{}", filepath.Clean(path))

	// 解析命令
	cmds := shellx.ParseCmd(cmdStr)

	// 检查命令是否存在, 仅对非shell执行有效
	if !findCmdUseShell.Get() {
		if _, err := shellx.FindCmd(cmds[0]); err != nil {
			return fmt.Errorf("找不到命令 %s: (提示: 对于内置命令如echo, 请使用--use-shell/-us标志)", cmds[0])
		}
	}

	// 如果启用了print-actions输出, 打印执行的命令
	if findCmdPrintActions.Get() {
		o.cl.Redf("exec: %v\n", cmdStr)
	}

	// 根据--use-shell/-us标志选择执行方式
	if findCmdUseShell.Get() {
		// 使用shell执行
		if err := shellx.NewCmdStr(cmdStr).WithStdout(os.Stdout).WithStderr(os.Stderr).Exec(); err != nil {
			return fmt.Errorf("命令执行失败: %v", err)
		}

	} else {
		// 原生直接执行
		if err := shellx.NewCmds(cmds).WithStdout(os.Stdout).WithStderr(os.Stderr).WithShell(shellx.ShellNone).Exec(); err != nil {
			return fmt.Errorf("命令执行失败: %v", err)
		}
	}

	return nil

}
