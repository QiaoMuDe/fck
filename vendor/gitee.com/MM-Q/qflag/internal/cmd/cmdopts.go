package cmd

import (
	"gitee.com/MM-Q/qflag/internal/types"
)

// CmdOpts 命令选项
//
// CmdOpts 提供了配置现有命令的方式，包含命令的所有可配置属性。
//
// 使用场景:
//   - 已有命令实例，需要批量设置属性
//   - 需要结构化的配置管理
//   - 需要部分配置（未设置的属性不会被修改）
//
// 示例:
//
//	cmd := NewCmd("myapp", "m", types.ExitOnError)
//	opts := &CmdOpts{
//	    Desc: "我的应用程序",
//	    Version: "1.0.0",
//	    UseChinese: true,
//	}
//	cmd.ApplyOpts(opts)
type CmdOpts struct {
	// 基本属性
	Desc               string // 命令描述
	Hidden             bool   // 是否隐藏命令, 不在帮助信息中显示
	DisableFlagParsing bool   // 是否禁用标志解析, 所有参数都作为位置参数

	// 运行函数
	RunFunc func(types.Command) error // 命令执行函数

	// 配置选项
	Version     string // 版本号
	UseChinese  bool   // 是否使用中文
	EnvPrefix   string // 环境变量前缀
	UsageSyntax string // 命令使用语法
	LogoText    string // Logo文本
	Completion  bool   // 是否启用自动补全标志

	// 环境变量绑定
	AutoBindEnv bool // 是否自动绑定所有标志的环境变量

	// 示例和说明
	Examples map[string]string // 示例使用, key为描述, value为示例命令
	Notes    []string          // 注意事项

	// 子命令和互斥组
	SubCmds        []types.Command       // 子命令列表, 用于添加到命令中
	MutexGroups    []types.MutexGroup    // 互斥组列表
	RequiredGroups []types.RequiredGroup // 必需组列表（支持条件性）
}

// NewCmdOpts 创建新的命令选项
//
// 返回值:
//   - *CmdOpts: 初始化的命令选项
//
// 功能说明:
//   - 创建基本命令选项
//   - 初始化所有字段为零值
//   - 初始化 map 和 slice 避免空指针
func NewCmdOpts() *CmdOpts {
	return &CmdOpts{
		Examples:       make(map[string]string),
		Notes:          []string{},
		SubCmds:        []types.Command{},
		MutexGroups:    []types.MutexGroup{},
		RequiredGroups: []types.RequiredGroup{},
	}
}
