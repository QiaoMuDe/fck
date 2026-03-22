package builtin

import (
	"os"

	"gitee.com/MM-Q/qflag/internal/completion"
	"gitee.com/MM-Q/qflag/internal/types"
)

// CompletionHandler 补全标志处理器
//
// CompletionHandler 负责处理补全标志 (--completion)
// 当用户指定补全标志时, 会生成对应的Shell自动补全脚本。
type CompletionHandler struct{}

// Handle 处理补全标志
//
// 参数:
//   - cmd: 要处理的命令
//
// 返回值:
//   - error: 处理失败时返回错误
//
// 功能说明:
//   - 从命令行参数获取Shell类型
//   - 生成对应的补全脚本
//   - 输出脚本并退出程序
func (h *CompletionHandler) Handle(cmd types.Command) error {
	// 获取shell类型参数
	shellType := getShellTypeFromArgs(cmd)

	// 生成补全脚本
	completion.GenAndPrint(cmd, shellType)
	os.Exit(0)
	return nil
}

// Type 返回标志类型
//
// 返回值:
//   - types.BuiltinFlagType: CompletionFlag
func (h *CompletionHandler) Type() types.BuiltinFlagType {
	return types.CompletionFlag
}

// ShouldRegister 判断是否应该注册补全标志
//
// 参数:
//   - cmd: 要检查的命令
//
// 返回值:
//   - bool: 根据配置决定是否注册
//
// 功能说明:
//   - 补全标志只在根命令中注册
//   - 只有当命令配置中 Completion 为 true 时才注册
//   - 默认不启用自动补全功能
func (h *CompletionHandler) ShouldRegister(cmd types.Command) bool {
	return cmd.IsRootCmd() && cmd.Config().Completion
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
//   - 检查补全标志是否已经被注册
//   - 避免重复注册，支持重复解析场景
func (h *CompletionHandler) ShouldSkipRegistration(cmd types.Command) bool {
	// 检查长名称是否已存在（补全标志没有短名称）
	_, exists := cmd.GetFlag(types.CompletionFlagName)
	return exists
}

// getShellTypeFromArgs 从命令行参数获取Shell类型
//
// 参数:
//   - cmd: 命令实例
//
// 返回值:
//   - string: Shell类型 (bash, pwsh)
//
// 功能说明:
//   - 检查命令行参数中的--completion参数值
//   - 默认返回bash
func getShellTypeFromArgs(cmd types.Command) string {
	f, ok := cmd.GetFlag(types.CompletionFlagName)
	if ok {
		return f.GetStr()
	}

	// 默认返回当前平台的Shell类型
	return types.CurrentShell()
}
