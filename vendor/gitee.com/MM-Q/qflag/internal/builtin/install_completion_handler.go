package builtin

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gitee.com/MM-Q/qflag/internal/completion"
	"gitee.com/MM-Q/qflag/internal/types"
)

// InstallCompletionHandler 安装补全标志处理器
//
// 负责处理 --install-completion 标志，自动完成：
// 1. 创建 ~/.qflag_completions/ 目录
// 2. 生成补全脚本到该目录
// 3. 将加载命令添加到 Shell 配置文件
type InstallCompletionHandler struct{}

// Handle 处理安装补全标志
//
// 参数:
//   - cmd: 要处理的命令
//
// 返回值:
//   - error: 处理失败时返回错误
//
// 功能说明:
//   - 获取用户家目录
//   - 创建补全脚本存放目录
//   - 生成补全脚本到文件
//   - 添加加载命令到 Shell 配置文件
//   - 输出成功信息并退出程序
func (h *InstallCompletionHandler) Handle(cmd types.Command) error {
	shellType := h.getShellTypeFromArgs(cmd)
	// 获取程序名并去掉可执行文件扩展名（如 .exe），将特殊字符替换为下划线
	programName := sanitizeProgramName(strings.TrimSuffix(filepath.Base(os.Args[0]), filepath.Ext(os.Args[0])))

	// 1. 获取用户家目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// 2. 创建补全目录
	completionsDir := filepath.Join(homeDir, types.CompletionsDirName)
	if err := os.MkdirAll(completionsDir, 0755); err != nil {
		return fmt.Errorf("failed to create completions directory: %w", err)
	}

	// 3. 确定脚本文件名和路径
	scriptName := programName
	if shellType == types.PwshShell || shellType == types.PowershellShell {
		scriptName += types.PwshCompletionScriptExt
	} else {
		scriptName += types.BashCompletionScriptExt
	}
	scriptPath := filepath.Join(completionsDir, scriptName)

	// 4. 生成补全脚本到文件
	scriptContent, err := completion.Generate(cmd, shellType)
	if err != nil {
		return fmt.Errorf("failed to generate completion script: %w", err)
	}

	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0644); err != nil {
		return fmt.Errorf("failed to write completion script: %w", err)
	}

	// 5. 添加加载命令到配置文件
	profilePath := h.getProfilePath(homeDir, shellType)
	if err := h.addLoadCommandToProfile(profilePath, scriptPath, shellType, programName); err != nil {
		return fmt.Errorf("failed to add load command to profile: %w", err)
	}

	// 6. 输出成功信息并退出程序
	h.printSuccessMessages(scriptPath, profilePath, shellType, cmd)

	os.Exit(0)
	return nil
}

// addLoadCommandToProfile 添加加载命令到配置文件
//
// 参数:
//   - profilePath: 配置文件路径
//   - scriptPath: 补全脚本路径
//   - shellType: Shell 类型
//   - programName: 程序名称
//
// 返回值:
//   - error: 添加失败时返回错误
func (h *InstallCompletionHandler) addLoadCommandToProfile(profilePath, scriptPath, shellType, programName string) error {
	loadCommand := h.generateLoadCommand(scriptPath, shellType, programName)

	// 读取现有内容
	content, err := os.ReadFile(profilePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// 检查是否已存在
	if err == nil && strings.Contains(string(content), scriptPath) {
		// 已存在，不需要重复添加
		return nil
	}

	// 追加到配置文件
	f, err := os.OpenFile(profilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	// 添加换行和注释
	if len(content) > 0 && !strings.HasSuffix(string(content), "\n") {
		_, _ = f.WriteString("\n")
	}
	_, _ = fmt.Fprintf(f, "\n"+types.CompletionScriptComment, programName)
	_, _ = f.WriteString(loadCommand + "\n")

	return nil
}

// getProfilePath 获取配置文件路径
//
// 参数:
//   - homeDir: 用户家目录
//   - shellType: Shell 类型
//
// 返回值:
//   - string: 配置文件路径
func (h *InstallCompletionHandler) getProfilePath(homeDir, shellType string) string {
	switch shellType {
	case types.PwshShell, types.PowershellShell:
		// PowerShell 配置文件
		if runtime.GOOS == "windows" {
			return filepath.Join(homeDir, types.PwshProfileDirWindows, types.PwshProfileFileName)
		}
		return filepath.Join(homeDir, types.PwshProfileDirUnix, types.PwshProfileFileName)
	default:
		// Bash 配置文件
		if runtime.GOOS == "darwin" {
			return filepath.Join(homeDir, types.BashProfileFileNameDarwin)
		}
		return filepath.Join(homeDir, types.BashProfileFileNameLinux)
	}
}

// generateLoadCommand 生成加载命令
//
// 参数:
//   - scriptPath: 补全脚本路径
//   - shellType: Shell 类型
//   - programName: 程序名称（用于生成唯一变量名）
//
// 返回值:
//   - string: 加载命令
func (h *InstallCompletionHandler) generateLoadCommand(scriptPath, shellType, programName string) string {
	switch shellType {
	case types.PwshShell, types.PowershellShell:
		// PowerShell: 使用程序名生成唯一变量名，避免多个程序冲突
		return fmt.Sprintf(types.PwshLoadCommandTemplate, programName, scriptPath, programName, programName)
	default:
		// Bash: 使用 -f 检查文件存在
		return fmt.Sprintf(types.BashLoadCommandTemplate, scriptPath, scriptPath)
	}
}

// getShellTypeFromArgs 从命令行参数获取 Shell 类型
//
// 参数:
//   - cmd: 命令实例
//
// 返回值:
//   - string: Shell 类型
func (h *InstallCompletionHandler) getShellTypeFromArgs(cmd types.Command) string {
	f, ok := cmd.GetFlag(types.InstallCompletionFlagName)
	if ok {
		return f.GetStr()
	}

	// 默认返回当前平台的 Shell 类型
	return types.CurrentShell()
}

// Type 返回标志类型
//
// 返回值:
//   - types.BuiltinFlagType: InstallCompletionFlag
func (h *InstallCompletionHandler) Type() types.BuiltinFlagType {
	return types.InstallCompletionFlag
}

// ShouldRegister 判断是否应该注册此标志
//
// 参数:
//   - cmd: 要检查的命令
//
// 返回值:
//   - bool: 是否应该注册
//
// 功能说明:
//   - 只在根命令中注册
//   - 只有当命令配置中 Completion 为 true 时才注册
func (h *InstallCompletionHandler) ShouldRegister(cmd types.Command) bool {
	return cmd.IsRootCmd() && cmd.Config().Completion
}

// ShouldSkipRegistration 判断是否应该跳过注册
//
// 参数:
//   - cmd: 要检查的命令
//
// 返回值:
//   - bool: 如果标志已存在则返回 true
//
// 功能说明:
//   - 检查安装补全标志是否已经被注册
//   - 避免重复注册，支持重复解析场景
func (h *InstallCompletionHandler) ShouldSkipRegistration(cmd types.Command) bool {
	_, exists := cmd.GetFlag(types.InstallCompletionFlagName)
	return exists
}

// printSuccessMessages 打印安装成功信息
//
// 参数:
//   - scriptPath: 补全脚本路径
//   - profilePath: 配置文件路径
//   - shellType: Shell 类型
//   - cmd: 命令实例（用于获取语言配置）
//
// 功能说明:
//   - 根据命令配置的语言选择中文或英文输出
//   - 根据 Shell 类型选择正确的执行命令（source 或 .）
//   - 使用预定义的常量格式化输出信息
func (h *InstallCompletionHandler) printSuccessMessages(scriptPath, profilePath, shellType string, cmd types.Command) {
	// 根据语言配置选择输出内容
	if cmd.Config().UseChinese {
		fmt.Printf(types.InstallSuccessScriptPathCN+"\n", scriptPath)
		fmt.Printf(types.InstallSuccessProfilePathCN+"\n", profilePath)
		fmt.Println(types.InstallSuccessHintCN)
	} else {
		fmt.Printf(types.InstallSuccessScriptPathEN+"\n", scriptPath)
		fmt.Printf(types.InstallSuccessProfilePathEN+"\n", profilePath)
		fmt.Println(types.InstallSuccessHintEN)
	}

	// 根据 Shell 类型选择执行命令（不区分中英文）
	if shellType == types.PwshShell || shellType == types.PowershellShell {
		fmt.Printf(types.InstallSuccessPwshCmd+"\n", profilePath)
	} else {
		fmt.Printf(types.InstallSuccessBashCmd+"\n", profilePath)
	}
}

// sanitizeProgramName 清理程序名，将特殊字符替换为下划线
//
// 参数:
//   - name: 原始程序名
//
// 返回值:
//   - string: 清理后的程序名
//
// 功能说明:
//   - 将横杠、空格等特殊字符替换为下划线
//   - 确保生成的变量名在 Shell 中合法
func sanitizeProgramName(name string) string {
	// 定义需要替换的特殊字符
	replacer := strings.NewReplacer(
		"-", "_",
		" ", "_",
		".", "_",
		":", "_",
		"/", "_",
		"\\", "_",
	)
	return replacer.Replace(name)
}
