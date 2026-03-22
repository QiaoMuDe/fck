package builtin

import (
	"fmt"

	"gitee.com/MM-Q/qflag/internal/flag"
	"gitee.com/MM-Q/qflag/internal/types"
)

// BuiltinFlagManager 内置标志管理器
//
// BuiltinFlagManager 负责管理所有内置标志的注册和处理。
// 它维护一个处理器映射表, 根据标志类型找到对应的处理器。
type BuiltinFlagManager struct {
	handlers map[types.BuiltinFlagType]types.BuiltinFlagHandler // 处理器映射表
}

// NewBuiltinFlagManager 创建内置标志管理器
//
// 返回值:
//   - *BuiltinFlagManager: 内置标志管理器实例
//
// 功能说明:
//   - 初始化处理器映射表
//   - 注册默认的内置标志处理器
func NewBuiltinFlagManager() *BuiltinFlagManager {
	m := &BuiltinFlagManager{
		handlers: make(map[types.BuiltinFlagType]types.BuiltinFlagHandler),
	}

	// 注册默认处理器
	m.RegisterHandler(&HelpHandler{})
	m.RegisterHandler(&VersionHandler{})
	m.RegisterHandler(&CompletionHandler{})

	return m
}

// RegisterHandler 注册内置标志处理器
//
// 参数:
//   - handler: 要注册的处理器
//
// 功能说明:
//   - 将处理器添加到处理器映射表
//   - 不再预注册标志名映射
func (m *BuiltinFlagManager) RegisterHandler(handler types.BuiltinFlagHandler) {
	flagType := handler.Type()
	m.handlers[flagType] = handler
}

// RegisterBuiltinFlags 注册内置标志
//
// 参数:
//   - cmd: 要注册标志的命令
//
// 返回值:
//   - error: 注册失败时返回错误
//
// 功能说明:
//   - 遍历所有处理器, 检查是否应该注册对应的标志
//   - 根据命令的语言设置使用相应的描述信息
//   - 创建并注册标志到命令中
//   - 如果标志已存在, 则跳过注册, 支持重复解析
func (m *BuiltinFlagManager) RegisterBuiltinFlags(cmd types.Command) error {
	for _, handler := range m.handlers {
		// 检查是否应该注册标志, 如果不应该, 则跳过
		if !handler.ShouldRegister(cmd) {
			continue
		}

		// 检查是否应该跳过注册（避免重复注册）
		if handler.ShouldSkipRegistration(cmd) {
			continue
		}

		// 根据标志类型创建并注册标志
		switch handler.Type() {
		case types.HelpFlag: // 注册帮助标志
			// 根据命令的语言设置使用相应的描述信息
			var desc string
			if cmd.Config().UseChinese {
				desc = "显示帮助信息"
			} else {
				desc = "Show help information"
			}
			helpFlag := flag.NewBoolFlag(types.HelpFlagName, types.HelpFlagShortName, desc, false)
			if err := cmd.AddFlag(helpFlag); err != nil {
				return err
			}

		case types.VersionFlag: // 注册版本标志
			// 根据命令的语言设置使用相应的描述信息
			var desc string
			if cmd.Config().UseChinese {
				desc = "显示版本信息"
			} else {
				desc = "Show version information"
			}
			versionFlag := flag.NewBoolFlag(types.VersionFlagName, types.VersionFlagShortName, desc, false)
			if err := cmd.AddFlag(versionFlag); err != nil {
				return err
			}

		case types.CompletionFlag: // 注册自动完成标志
			// 根据命令的语言设置使用相应的描述信息
			var desc string
			if cmd.Config().UseChinese {
				desc = fmt.Sprintf("生成Shell自动补全脚本, 支持的Shell: %v", types.SupportedShells)
			} else {
				desc = fmt.Sprintf("Generate shell completion script. Supported shells: %v", types.SupportedShells)
			}
			completionFlag := flag.NewEnumFlag(types.CompletionFlagName, "", desc, types.CurrentShell(), types.SupportedShells)
			if err := cmd.AddFlag(completionFlag); err != nil {
				return err
			}

			// 注册补全标志后注册内置的示例信息
			cmd.AddExamples(types.GetCompletionExample())
		}
	}

	return nil
}

// HandleBuiltinFlags 处理内置标志
//
// 参数:
//   - cmd: 要处理标志的命令
//
// 返回值:
//   - error: 处理失败时返回错误
//
// 功能说明:
//   - 遍历命令的所有标志, 检查是否是内置标志
//   - 如果是内置标志且被设置, 则执行对应的处理器
//   - 传入当前命令进行动态检查
func (m *BuiltinFlagManager) HandleBuiltinFlags(cmd types.Command) error {
	flags := cmd.Flags()

	for _, f := range flags {
		// 传入当前命令进行动态检查
		if flagType, isBuiltin := m.isBuiltinFlag(f, cmd); isBuiltin {
			// 检查是否被设置
			if f.IsSet() {
				// 执行处理器
				if handler, exists := m.handlers[flagType]; exists {
					return handler.Handle(cmd)
				}
			}
		}
	}

	return nil
}

// isBuiltinFlag 检查是否是内置标志
//
// 参数:
//   - f: 要检查的标志
//   - cmd: 当前命令实例
//
// 返回值:
//   - types.BuiltinFlagType: 标志类型
//   - bool: 是否是内置标志
//
// 功能说明:
//   - 动态检查标志是否为内置标志
//   - 基于当前命令的配置和 ShouldRegister 判断
//   - 解决名称冲突问题
func (m *BuiltinFlagManager) isBuiltinFlag(f types.Flag, cmd types.Command) (types.BuiltinFlagType, bool) {
	// 遍历所有处理器，检查是否应该注册且名称匹配
	for _, handler := range m.handlers {
		// 检查该标志类型是否真的会在当前命令中注册
		if !handler.ShouldRegister(cmd) {
			continue
		}

		// 检查名称是否匹配
		switch handler.Type() {
		case types.HelpFlag:
			// 检查是否为帮助标志
			if f.LongName() == types.HelpFlagName || f.ShortName() == types.HelpFlagShortName {
				return types.HelpFlag, true
			}

		case types.VersionFlag:
			// 检查是否为版本标志
			if f.LongName() == types.VersionFlagName || f.ShortName() == types.VersionFlagShortName {
				return types.VersionFlag, true
			}

		case types.CompletionFlag:
			// 检查是否为自动补全标志
			if f.LongName() == types.CompletionFlagName {
				return types.CompletionFlag, true
			}
		}
	}

	return 0, false
}
