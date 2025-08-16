package find

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/commands/internal/common"
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
func (o *FileOperator) Delete(path string, isDir bool) error {
	// 检查是否为空路径
	if path == "" {
		return fmt.Errorf("没有可删除的路径")
	}

	// 先检查文件/目录是否存在
	if _, err := os.Lstat(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("文件/目录不存在: %s", path)
		}
		return fmt.Errorf("检查文件/目录时出错: %s: %v", path, err)
	}

	// 打印删除信息
	if findCmdPrintDelete.Get() {
		o.cl.Redf("del: %s\n", path)
	}

	// 根据类型选择删除方法
	var rmErr error
	if isDir {
		// 检查目录是否为空
		dirEntries, readDirErr := os.ReadDir(path)
		if readDirErr != nil {
			return common.HandleError(path, readDirErr)
		}

		// 根据目录是否为空选择删除方法
		if len(dirEntries) > 0 {
			// 目录不为空, 递归删除
			rmErr = os.RemoveAll(path)
		} else {
			// 删除空目录
			rmErr = os.Remove(path)
		}
	} else {
		// 删除文件
		rmErr = os.Remove(path)
	}

	if rmErr != nil {
		return fmt.Errorf("删除失败: %s: %v", path, rmErr)
	}

	return nil
}

// Move 移动匹配的文件或目录到指定位置
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

	// 尝试移动前先检查权限
	if err := o.checkWritePermission(filepath.Dir(absTargetPath)); err != nil {
		return fmt.Errorf("目标目录无写入权限: %v", err)
	}

	// 检查目标文件是否已存在
	if _, err := os.Stat(absTargetPath); err == nil {
		// 如果是移动操作, 直接跳过而不是报错
		if findCmdMove.Get() != "" {
			return nil
		}
		return fmt.Errorf("目标文件已存在: %s", absTargetPath)
	}

	// 打印移动信息
	if findCmdPrintMove.Get() {
		o.cl.Redf("%s -> %s\n", absSearchPath, absTargetPath)
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

// Execute 执行指定的命令，支持跨平台并检查shell是否存在
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
	safeCmdStr, err := o.sanitizeCommand(cmdStr, path)
	if err != nil {
		return fmt.Errorf("命令安全检查失败: %v", err)
	}

	// 根据操作系统选择shell和参数
	shell, args := o.getShellCommand(safeCmdStr)

	// 检查shell是否存在
	if _, err := exec.LookPath(shell); err != nil {
		return fmt.Errorf("找不到 %s 解释器: %v", shell, err)
	}

	// 如果启用了print-cmd输出, 打印执行的命令
	if findCmdPrintCmd.Get() {
		o.cl.Redf("%s %s\n", shell, strings.Join(args, " "))
	}

	// 构建命令并设置输出
	cmd := exec.Command(shell, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 执行命令并捕获错误
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("命令执行失败: %v", err)
	}

	return nil
}

// sanitizeCommand 安全地处理命令字符串，防止命令注入
func (o *FileOperator) sanitizeCommand(cmdStr, path string) (string, error) {
	// 清理路径，移除潜在的危险字符
	cleanPath := filepath.Clean(path)

	// 检查路径是否包含危险字符
	dangerousChars := []string{";", "&", "|", "`", "$", "(", ")", "<", ">", "\"", "'"}
	for _, char := range dangerousChars {
		if strings.Contains(cleanPath, char) {
			// 对危险字符进行转义而不是拒绝
			cleanPath = strings.ReplaceAll(cleanPath, char, "\\"+char)
		}
	}

	// 根据系统类型选择引用方式
	var quotedPath string
	if runtime.GOOS == "windows" {
		quotedPath = fmt.Sprintf("\"%s\"", cleanPath) // Windows使用双引号
	} else {
		quotedPath = fmt.Sprintf("'%s'", cleanPath) // Linux使用单引号
	}

	// 替换{}为安全的路径
	safeCmdStr := strings.ReplaceAll(cmdStr, "{}", quotedPath)

	return safeCmdStr, nil
}

// getShellCommand 根据操作系统获取shell命令和参数
func (o *FileOperator) getShellCommand(cmdStr string) (string, []string) {
	if runtime.GOOS == "windows" {
		// 先尝试使用 PowerShell
		if _, err := exec.LookPath("powershell"); err == nil {
			return "powershell", []string{"-Command", cmdStr}
		}
		// 如果 PowerShell 不存在, 使用 cmd
		return "cmd", []string{"/C", cmdStr}
	}
	// Unix-like系统使用bash
	return "bash", []string{"-c", cmdStr}
}

// checkWritePermission 检查目录的写入权限
func (o *FileOperator) checkWritePermission(dir string) error {
	tmpFile := filepath.Join(dir, ".fck_tmp")
	if err := os.WriteFile(tmpFile, []byte{}, 0600); err != nil {
		return err
	}
	if err := os.Remove(tmpFile); err != nil {
		o.cl.PrintErrorf("删除临时文件失败: %v\n", err)
	}
	return nil
}
