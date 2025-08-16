package find

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/commands/internal/common"
)

// shell 检测缓存
var (
	detectedShell     string    // 检测到的shell
	detectedShellArgs []string  // shell参数
	shellMutex        sync.Once // shell检测互斥锁
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
//   - isDir: 文件/目录类型
//
// 返回:
//   - error: 错误信息
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
	if findCmdPrintMove.Get() {
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

	// 根据--use-shell标志选择执行方式
	if findCmdUseShell.Get() {
		return o.executeWithShell(cmdStr, path)
	} else {
		return o.executeDirect(cmdStr, path)
	}
}

// quotePath 安全地引用文件路径，防止路径中的特殊字符导致问题
//
// 参数:
//   - path: 文件路径
//
// 返回:
//   - string: 引用后的安全路径
func (o *FileOperator) quotePath(path string) string {
	// 清理路径
	cleanPath := filepath.Clean(path)

	// 根据系统类型选择引用方式
	if runtime.GOOS == "windows" {
		// Windows使用双引号，并转义内部的双引号
		escapedPath := strings.ReplaceAll(cleanPath, "\"", "\\\"")
		return fmt.Sprintf("\"%s\"", escapedPath)
	} else {
		// Unix系统使用单引号，并处理内部的单引号
		escapedPath := strings.ReplaceAll(cleanPath, "'", "'\"'\"'")
		return fmt.Sprintf("'%s'", escapedPath)
	}
}

// executeDirect 直接执行命令（默认模式）
//
// 参数:
//   - cmdStr: 要执行的命令字符串
//   - path: 文件路径
//
// 返回:
//   - error: 错误信息
func (o *FileOperator) executeDirect(cmdStr, path string) error {
	// 安全地替换{}为实际的文件路径
	safePath := o.quotePath(path)
	finalCmd := strings.ReplaceAll(cmdStr, "{}", safePath)

	// 解析命令参数
	args, err := o.parseCommand(finalCmd)
	if err != nil {
		return fmt.Errorf("解析命令失败: %v", err)
	}

	if len(args) == 0 {
		return fmt.Errorf("解析后的命令为空")
	}

	// 检查命令是否存在
	if _, err := exec.LookPath(args[0]); err != nil {
		return fmt.Errorf("找不到命令 %s (提示: 对于内置命令如echo, 请使用--use-shell/-us标志)", args[0])
	}

	// 如果启用了print-cmd输出, 打印执行的命令
	if findCmdPrintCmd.Get() {
		o.cl.Redf("exec(direct): %v\n", args)
	}

	// 构建命令并设置输出
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = filepath.Dir(path) // 设置工作目录
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 执行命令并捕获错误
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("命令执行失败: %v", err)
	}

	return nil
}

// executeWithShell 通过shell执行命令（--use-shell模式）
//
// 参数:
//   - cmdStr: 要执行的命令字符串
//   - path: 文件路径
//
// 返回:
//   - error: 错误信息
func (o *FileOperator) executeWithShell(cmdStr, path string) error {
	// 安全地替换{}为实际的文件路径
	safePath := o.quotePath(path)
	finalCmd := strings.ReplaceAll(cmdStr, "{}", safePath)

	// 根据操作系统选择shell和参数
	shell, args := o.getShellCommand(finalCmd)

	// 检查shell是否存在
	if _, err := exec.LookPath(shell); err != nil {
		return fmt.Errorf("找不到 %s 解释器: %v", shell, err)
	}

	// 如果启用了print-cmd输出, 打印执行的命令
	if findCmdPrintCmd.Get() {
		o.cl.Redf("exec(shell): [%s %s]\n", shell, strings.Join(args, " "))
	}

	// 构建命令并设置输出
	cmd := exec.Command(shell, args...)
	cmd.Dir = filepath.Dir(path) // 设置工作目录
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 执行命令并捕获错误
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("命令执行失败: %v", err)
	}

	return nil
}

// initShellDetection 初始化shell检测，只执行一次
func initShellDetection() {
	if runtime.GOOS == "windows" {
		// Windows 平台按优先级检测: pwsh -> powershell -> cmd
		windowsShells := []struct {
			name string
			flag string
		}{
			{"pwsh", "-Command"},       // PowerShell Core (跨平台)
			{"powershell", "-Command"}, // Windows PowerShell
			{"cmd", "/C"},              // Command Prompt
		}

		for _, shell := range windowsShells {
			if _, err := exec.LookPath(shell.name); err == nil {
				detectedShell = shell.name
				detectedShellArgs = []string{shell.flag}
				return
			}
		}

		// 如果都找不到，使用默认的 cmd（理论上不会发生）
		detectedShell = "cmd"
		detectedShellArgs = []string{"/C"}
	} else {
		// Unix/Linux 平台按优先级检测: bash -> sh -> zsh
		unixShells := []string{"bash", "sh", "zsh"}

		for _, shell := range unixShells {
			if _, err := exec.LookPath(shell); err == nil {
				detectedShell = shell
				detectedShellArgs = []string{"-c"}
				return
			}
		}

		// 如果都找不到，使用默认的 sh（理论上不会发生）
		detectedShell = "sh"
		detectedShellArgs = []string{"-c"}
	}
}

// getShellCommand 根据操作系统获取shell命令和参数（带缓存优化）
//
// 参数:
//   - cmdStr: 要执行的命令字符串
//
// 返回:
//   - string: shell命令
//   - []string: shell参数
//
// 支持的shell:
//   - Windows: pwsh -> powershell -> cmd (按优先级)
//   - Unix/Linux: bash -> sh -> zsh (按优先级)
func (o *FileOperator) getShellCommand(cmdStr string) (string, []string) {
	// 使用 sync.Once 确保只检测一次
	shellMutex.Do(initShellDetection)

	// 构建完整的参数列表
	args := make([]string, len(detectedShellArgs)+1) // +1 是为了包含命令字符串
	copy(args, detectedShellArgs)                    // 复制shell参数
	args[len(args)-1] = cmdStr                       // 将执行的命令添加到最后

	return detectedShell, args
}

// parseCommand 解析命令字符串为参数数组
//
// 参数:
//   - cmdStr: 命令字符串
//
// 返回:
//   - []string: 解析后的参数数组
//   - error: 解析错误
func (o *FileOperator) parseCommand(cmdStr string) ([]string, error) {
	args := strings.Fields(strings.TrimSpace(cmdStr))
	if len(args) == 0 {
		return nil, fmt.Errorf("命令为空")
	}
	return args, nil
}
