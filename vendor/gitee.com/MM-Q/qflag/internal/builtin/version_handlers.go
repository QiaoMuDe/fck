package builtin

import (
	"fmt"
	"os"

	"gitee.com/MM-Q/qflag/internal/types"
)

// VersionHandler 版本标志处理器
//
// VersionHandler 负责处理版本标志 (-v/--version) 。
// 当用户指定版本标志时, 会打印命令的版本信息并退出程序。
type VersionHandler struct{}

// Handle 处理版本标志
//
// 参数:
//   - cmd: 要处理的命令
//
// 返回值:
//   - error: 处理失败时返回错误
//
// 功能说明:
//   - 打印命令的版本信息
//   - 使用状态码0退出程序
func (h *VersionHandler) Handle(cmd types.Command) error {
	fmt.Println(cmd.Config().Version)
	os.Exit(0)
	return nil
}

// Type 返回标志类型
//
// 返回值:
//   - types.BuiltinFlagType: VersionFlag
func (h *VersionHandler) Type() types.BuiltinFlagType {
	return types.VersionFlag
}

// ShouldRegister 判断是否应该注册版本标志
//
// 参数:
//   - cmd: 要检查的命令
//
// 返回值:
//   - bool: 如果设置了版本信息返回true, 否则返回false
//
// 功能说明:
//   - 只有在命令设置了版本信息时才注册版本标志
//   - 版本标志只在根命令中注册
func (h *VersionHandler) ShouldRegister(cmd types.Command) bool {
	return cmd.IsRootCmd() && cmd.Config().Version != ""
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
//   - 检查版本标志是否已经被注册
//   - 避免重复注册，支持重复解析场景
func (h *VersionHandler) ShouldSkipRegistration(cmd types.Command) bool {
	// 检查长名称和短名称是否已存在
	_, longExists := cmd.GetFlag(types.VersionFlagName)
	_, shortExists := cmd.GetFlag(types.VersionFlagShortName)
	return longExists || shortExists
}
