// Package shellx 定义了shell命令执行库的核心数据类型。
// 本文件定义了ShellType枚举，提供了shell类型管理。
//
// 主要类型：
//   - ShellType: Shell类型枚举，支持sh、bash、cmd、powershell等多种shell
//
// ShellType支持的shell类型：
//   - ShellSh: Unix/Linux sh shell
//   - ShellBash: Bash shell
//   - ShellCmd: Windows Command Prompt
//   - ShellPowerShell: Windows PowerShell
//   - ShellPwsh: PowerShell Core (跨平台)
//   - ShellNone: 直接执行命令，不使用shell
//   - ShellDef1: 根据操作系统自动选择默认shell(Windows系统默认为cmd, 其他系统默认为sh)
//   - ShellDef2: 根据操作系统自动选择默认shell(Windows系统默认为powershell, 其他系统默认为sh)
package shellx

import (
	"runtime"
)

// ShellType 定义shell类型
type ShellType int

const (
	ShellSh         ShellType = iota // sh shell
	ShellBash                        // bash shell
	ShellPwsh                        // pwsh (PowerShell Core)
	ShellPowerShell                  // powershell (Windows PowerShell)
	ShellCmd                         // cmd (Windows Command Prompt)
	ShellNone                        // 无shell, 直接原生的执行命令
	ShellDef1                        // 默认shell, 根据操作系统自动选择(Windows系统默认为cmd, 其他系统默认为sh)
	ShellDef2                        // 默认shell, 根据操作系统自动选择(Windows系统默认为powershell, 其他系统默认为sh)
)

// String 返回shell类型的字符串表示
func (s ShellType) String() string {
	switch s {
	case ShellSh:
		return "sh"

	case ShellBash:
		return "bash"

	case ShellPwsh:
		return "pwsh"

	case ShellPowerShell:
		return "powershell"

	case ShellCmd:
		return "cmd"

	case ShellNone:
		return "none"

	case ShellDef1:
		if runtime.GOOS == "windows" {
			return "cmd"
		}
		return "sh"

	case ShellDef2:
		if runtime.GOOS == "windows" {
			return "powershell"
		}
		return "sh"

	default:
		return "unknown"
	}
}

// shellFlags 用于返回shell类型的标志
func (s ShellType) shellFlags() string {
	switch s {
	case ShellSh:
		return "-c"

	case ShellBash:
		return "-c"

	case ShellPwsh:
		return "-Command"

	case ShellPowerShell:
		return "-Command"

	case ShellCmd:
		return "/c"

	case ShellNone:
		return ""

	case ShellDef1:
		if runtime.GOOS == "windows" {
			return "/c"
		}
		return "-c"

	case ShellDef2:
		if runtime.GOOS == "windows" {
			return "-Command"
		}
		return "-c"

	default:
		return ""
	}
}
