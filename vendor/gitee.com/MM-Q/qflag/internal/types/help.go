package types

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// programName 缓存程序名，避免重复计算
// 在包初始化时立即获取
var programName = filepath.Base(os.Args[0])

// 用于存储选项的信息
type OptionInfo struct {
	NamePart string // 选项名称部分
	Desc     string // 选项描述
	DefValue string // 选项默认值
}

// 用于存储子命令的信息
type SubCmdInfo struct {
	Name string // 子命令名称
	Desc string // 子命令描述
}

// 帮助信息标题 - 中文
const (
	HelpNameCN     = "名称:\n"
	HelpDescCN     = "\n描述:\n"
	HelpUsageCN    = "\n用法:\n"
	HelpOptionsCN  = "\n选项:\n"
	HelpSubCmdsCN  = "\n子命令:\n"
	HelpExamplesCN = "\n示例:\n"
	HelpNotesCN    = "\n注意:\n"
)

// 帮助信息标题 - 英文
const (
	HelpNameEN     = "Name:\n"
	HelpDescEN     = "\nDesc:\n"
	HelpUsageEN    = "\nUsage:\n"
	HelpOptionsEN  = "\nOptions:\n"
	HelpSubCmdsEN  = "\nSubcmds:\n"
	HelpExamplesEN = "\nExamples:\n"
	HelpNotesEN    = "\nNotes:\n"
)

// 统一的前缀, 缩进两个空格
const HelpPrefix = "  "

// 选项和子命令中间空格, 缩进6个空格
const HelpOptionSubCmdSpace = "      "

// 补全示例信息 - Windows 中文
var HelpCompletionExampleWinCN = map[string]string{
	"临时启用": fmt.Sprintf("%s --completion pwsh | Out-String | Invoke-Expression", programName),
}

// 补全示例信息 - Windows 英文
var HelpCompletionExampleWinEN = map[string]string{
	"Temporary": fmt.Sprintf("%s --completion pwsh | Out-String | Invoke-Expression", programName),
}

// 补全示例信息 - Unix-like 系统中文（Linux 和 macOS 共用）
var HelpCompletionExampleUnixCN = map[string]string{
	"临时启用": fmt.Sprintf("source <(%s --completion bash)", programName),
}

// 补全示例信息 - Unix-like 系统英文（Linux 和 macOS 共用）
var HelpCompletionExampleUnixEN = map[string]string{
	"Temporary": fmt.Sprintf("source <(%s --completion bash)", programName),
}

// GetCompletionExample 获取当前平台的补全示例信息（中文）
//
// 返回值:
//   - map[string]string: 包含补全示例信息的映射
//
// 功能说明:
//   - 根据当前运行的操作系统返回对应的中文补全示例
//   - 支持 Windows、Linux 和 macOS 平台
//   - 提供临时启用和永久启用两种方式的示例
func GetCompletionExample() map[string]string {
	switch runtime.GOOS {
	case "windows":
		return HelpCompletionExampleWinCN
	default:
		// Linux 和 macOS 使用相同的示例
		return HelpCompletionExampleUnixCN
	}
}

// GetCompletionExampleEN 获取当前平台的补全示例信息（英文）
//
// 返回值:
//   - map[string]string: 包含补全示例信息的映射
//
// 功能说明:
//   - 根据当前运行的操作系统返回对应的英文补全示例
//   - 支持 Windows、Linux 和 macOS 平台
//   - 提供临时启用和永久启用两种方式的示例
func GetCompletionExampleEN() map[string]string {
	switch runtime.GOOS {
	case "windows":
		return HelpCompletionExampleWinEN
	default:
		// Linux 和 macOS 使用相同的示例
		return HelpCompletionExampleUnixEN
	}
}

// GetInstallCompletionExample 获取当前平台的安装补全示例信息（中文）
//
// 返回值:
//   - map[string]string: 包含安装补全示例信息的映射
//
// 功能说明:
//   - 根据当前运行的操作系统返回对应的中文安装补全示例
//   - Windows 使用 pwsh，其他平台使用 bash
//   - 作为永久启用的推荐方式
func GetInstallCompletionExample() map[string]string {
	shell := CurrentShell()
	return map[string]string{
		"永久启用": fmt.Sprintf("%s --install-completion %s", programName, shell),
	}
}

// GetInstallCompletionExampleEN 获取当前平台的安装补全示例信息（英文）
//
// 返回值:
//   - map[string]string: 包含安装补全示例信息的映射
//
// 功能说明:
//   - 根据当前运行的操作系统返回对应的英文安装补全示例
//   - Windows 使用 pwsh，其他平台使用 bash
//   - 作为永久启用的推荐方式
func GetInstallCompletionExampleEN() map[string]string {
	shell := CurrentShell()
	return map[string]string{
		"Permanent": fmt.Sprintf("%s --install-completion %s", programName, shell),
	}
}
