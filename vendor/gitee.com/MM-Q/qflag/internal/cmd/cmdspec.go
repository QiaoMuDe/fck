// Package cmd 提供命令规格结构体, 用于通过配置创建命令
//
// cmdspec 包实现了通过规格结构体创建命令的功能, 提供了更直观、集中的命令配置方式。
// 主要组件:
//   - CmdSpec: 命令规格结构体
//   - NewCmdSpec: 便捷构造函数
//   - NewCmdFromSpec: 从规格创建命令的函数
//
// 特性:
//   - 支持所有命令属性的集中配置
//   - 支持嵌套子命令
//   - 提供默认值处理
//   - 完全兼容现有API
package cmd

import (
	"fmt"

	"gitee.com/MM-Q/qflag/internal/types"
)

// CmdSpec 命令规格结构体
//
// CmdSpec 提供了通过规格创建命令的方式, 包含命令的所有属性。
// 这种方式比函数式配置更加直观和集中。
type CmdSpec struct {
	// 基本属性
	LongName      string              // 命令长名称
	ShortName     string              // 命令短名称
	Desc          string              // 命令描述
	ErrorHandling types.ErrorHandling // 错误处理策略

	// 运行函数
	RunFunc func(types.Command) error // 命令执行函数

	// 配置选项
	Version     string // 版本号
	UseChinese  bool   // 是否使用中文
	EnvPrefix   string // 环境变量前缀
	UsageSyntax string // 命令使用语法
	LogoText    string // Logo文本
	Completion  bool   // 是否启用自动补全标志

	// 示例和说明
	Examples map[string]string // 示例使用, key为描述, value为示例命令
	Notes    []string          // 注意事项

	// 子命令和互斥组
	SubCmds        []types.Command       // 子命令列表, 用于添加到命令中
	MutexGroups    []types.MutexGroup    // 互斥组列表
	RequiredGroups []types.RequiredGroup // 必需组列表（支持条件性）
}

// NewCmdSpec 创建新的命令规格
//
// 参数:
//   - longName: 命令长名称
//   - shortName: 命令短名称
//
// 返回值:
//   - *CmdSpec: 初始化的命令规格
//
// 功能说明:
//   - 创建基本命令规格
//   - 设置默认值
//   - 初始化所有字段
func NewCmdSpec(longName, shortName string) *CmdSpec {
	return &CmdSpec{
		LongName:       longName,
		ShortName:      shortName,
		ErrorHandling:  types.ExitOnError, // 默认错误处理策略
		UseChinese:     false,             // 默认不使用中文
		Completion:     false,             // 默认不启用自动补全
		Examples:       make(map[string]string),
		Notes:          []string{},
		SubCmds:        []types.Command{},
		MutexGroups:    []types.MutexGroup{},
		RequiredGroups: []types.RequiredGroup{},
	}
}

// NewCmdFromSpec 从规格创建命令
//
// 参数:
//   - spec: 命令规格结构体
//
// 返回值:
//   - *Cmd: 创建的命令实例
//   - error: 创建失败时返回错误
//
// 功能说明:
//   - 根据规格结构体创建命令
//   - 自动设置所有属性和配置
//   - 递归创建子命令
//   - 支持默认值处理
//   - 使用defer捕获panic, 转换为错误返回
func NewCmdFromSpec(spec *CmdSpec) (cmd *Cmd, err error) {
	// 使用defer捕获panic, 转换为错误返回
	defer func() {
		if r := recover(); r != nil {
			// 将panic转换为错误
			switch x := r.(type) {
			case string:
				err = types.NewError("PANIC", x, nil)
			case error:
				err = types.NewError("PANIC", x.Error(), x)
			default:
				err = types.NewError("PANIC", fmt.Sprintf("%v", x), nil)
			}
			cmd = nil
		}
	}()

	// 验证规格结构体
	if spec == nil {
		return nil, types.NewError("INVALID_CMD_SPEC", "command spec cannot be nil", nil)
	}

	// 创建基本命令
	cmd = NewCmd(spec.LongName, spec.ShortName, spec.ErrorHandling)

	// 设置基本属性
	cmd.SetDesc(spec.Desc)   // 命令描述
	cmd.SetRun(spec.RunFunc) // 命令执行函数

	// 设置配置选项
	cmd.SetVersion(spec.Version)         // 版本号
	cmd.SetChinese(spec.UseChinese)      // 是否使用中文
	cmd.SetEnvPrefix(spec.EnvPrefix)     // 环境变量前缀
	cmd.SetUsageSyntax(spec.UsageSyntax) // 命令使用语法
	cmd.SetLogoText(spec.LogoText)       // Logo文本
	cmd.SetCompletion(spec.Completion)   // 是否启用自动补全

	// 添加示例和说明
	if len(spec.Examples) > 0 {
		cmd.AddExamples(spec.Examples)
	}
	if len(spec.Notes) > 0 {
		cmd.AddNotes(spec.Notes)
	}

	// 添加互斥组
	if len(spec.MutexGroups) > 0 {
		for _, group := range spec.MutexGroups {
			if err := cmd.AddMutexGroup(group.Name, group.Flags, group.AllowNone); err != nil {
				return nil, types.WrapError(err, "FAILED_TO_ADD_MUTEX_GROUP", "failed to add mutex group")
			}
		}
	}

	// 添加必需组
	if len(spec.RequiredGroups) > 0 {
		for _, group := range spec.RequiredGroups {
			if err := cmd.AddRequiredGroup(group.Name, group.Flags, group.Conditional); err != nil {
				return nil, types.WrapError(err, "FAILED_TO_ADD_REQUIRED_GROUP", "failed to add required group")
			}
		}
	}

	// 添加子命令
	if len(spec.SubCmds) > 0 {
		if err := cmd.AddSubCmds(spec.SubCmds...); err != nil {
			return nil, types.WrapError(err, "FAILED_TO_ADD_SUBCMDS", "failed to add subcommands")
		}
	}

	return cmd, nil
}
