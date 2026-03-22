package types

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// 用于存储选项的信息
type OptionInfo struct {
	NamePart string
	Desc     string
	DefValue string
}

// 用于存储子命令的信息
type SubCmdInfo struct {
	Name string
	Desc string
}

// 帮助信息标题
const (
	// 中文帮助信息标题
	HelpNameCN     = "名称:\n"
	HelpDescCN     = "\n描述:\n"
	HelpUsageCN    = "\n用法:\n"
	HelpOptionsCN  = "\n选项:\n"
	HelpSubCmdsCN  = "\n子命令:\n"
	HelpExamplesCN = "\n示例:\n"
	HelpNotesCN    = "\n注意:\n"

	// 英文帮助信息标题
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

// 内置补全示例信息 - Windows
var HelpCompletionExampleWin = map[string]string{
	"Windows 临时启用": fmt.Sprintf("%s --completion pwsh | Out-String | Invoke-Expression", filepath.Base(os.Args[0])),
	"Windows 永久启用": fmt.Sprintf("echo '%s --completion pwsh | Out-String | Invoke-Expression' >> $PROFILE", filepath.Base(os.Args[0])),
}

// 内置补全示例信息 - Linux
var HelpCompletionExampleLinux = map[string]string{
	"Linux 临时启用": fmt.Sprintf("source <(%s --completion bash)", filepath.Base(os.Args[0])),
	"Linux 永久启用": fmt.Sprintf("echo 'source <(%s --completion bash)' >> ~/.bashrc", filepath.Base(os.Args[0])),
}

// 内置补全示例信息 - macOS
var HelpCompletionExampleMac = map[string]string{
	"macOS 临时启用": fmt.Sprintf("source <(%s --completion bash)", filepath.Base(os.Args[0])),
	"macOS 永久启用": fmt.Sprintf("echo 'source <(%s --completion bash)' >> ~/.bash_profile", filepath.Base(os.Args[0])),
}

// GetCompletionExample 获取当前平台的补全示例信息
//
// 返回值:
//   - map[string]string: 包含补全示例信息的映射
//
// 功能说明:
//   - 根据当前运行的操作系统返回对应的补全示例
//   - 支持 Windows、Linux 和 macOS 平台
//   - 提供临时启用和永久启用两种方式的示例
func GetCompletionExample() map[string]string {
	switch runtime.GOOS {
	case "windows":
		return HelpCompletionExampleWin
	case "linux":
		return HelpCompletionExampleLinux
	case "darwin":
		return HelpCompletionExampleMac
	default:
		// 默认返回 Linux 示例, 适用于大多数 Unix-like 系统
		return HelpCompletionExampleLinux
	}
}
