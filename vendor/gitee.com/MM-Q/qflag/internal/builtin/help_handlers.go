package builtin

import (
	"os"

	"gitee.com/MM-Q/qflag/internal/types"
)

// HelpHandler 帮助标志处理器
//
// HelpHandler 负责处理帮助标志 (-h/--help) 。
// 当用户指定帮助标志时, 会打印命令的帮助信息并退出程序。
type HelpHandler struct{}

// Handle 处理帮助标志
//
// 参数:
//   - cmd: 要处理的命令
//
// 返回值:
//   - error: 处理失败时返回错误
//
// 功能说明:
//   - 打印命令的帮助信息
//   - 使用状态码0退出程序
func (h *HelpHandler) Handle(cmd types.Command) error {
	cmd.PrintHelp()
	os.Exit(0)
	return nil
}

// Type 返回标志类型
//
// 返回值:
//   - types.BuiltinFlagType: HelpFlag
func (h *HelpHandler) Type() types.BuiltinFlagType {
	return types.HelpFlag
}

// ShouldRegister 判断是否应该注册帮助标志
//
// 参数:
//   - cmd: 要检查的命令
//
// 返回值:
//   - bool: 总是返回true
//
// 功能说明:
//   - 帮助标志总是注册, 因为所有命令都应该支持帮助功能
func (h *HelpHandler) ShouldRegister(cmd types.Command) bool {
	return true
}

// ShouldSkipRegistration 判断是否应该跳过注册
//
// 参数:
//   - cmd: 要检查的命令
//
// 返回值:
//   - bool: 如果标志已存在则返回true
//
// 功能说明:
//   - 检查帮助标志是否已经被注册
//   - 避免重复注册，支持重复解析场景
func (h *HelpHandler) ShouldSkipRegistration(cmd types.Command) bool {
	// 检查长名称和短名称是否已存在
	_, longExists := cmd.GetFlag(types.HelpFlagName)
	_, shortExists := cmd.GetFlag(types.HelpFlagShortName)
	return longExists || shortExists
}
