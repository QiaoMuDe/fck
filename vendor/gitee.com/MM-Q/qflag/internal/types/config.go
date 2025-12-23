// Package types 配置结构体和选项定义
// 本文件定义了命令配置相关的结构体和选项，包括命令的各种配置参数、
// 帮助信息设置、版本信息等配置数据的定义和管理。
package types

// CmdConfig 命令行配置
type CmdConfig struct {
	// 版本信息
	Version string

	// 自定义描述
	Desc string

	// 自定义的完整命令行帮助信息
	Help string

	// 自定义用法格式说明
	UsageSyntax string

	// 模块帮助信息
	ModuleHelps string

	// logo文本
	LogoText string

	// 备注内容切片
	Notes []string

	// 示例信息切片
	Examples []ExampleInfo

	// 是否使用中文帮助信息
	UseChinese bool

	// 禁用内置标志自动退出
	NoFgExit bool

	// 控制是否启用自动补全功能
	Completion bool
}

// ExampleInfo 示例信息结构体
// 用于存储命令的使用示例，包括描述和示例内容
//
// 字段:
//   - Desc: 示例描述
//   - Usage: 示例使用方式
type ExampleInfo struct {
	Desc  string // 示例描述
	Usage string // 示例使用方式
}

// NewCmdConfig 创建一个新的CmdConfig实例
func NewCmdConfig() *CmdConfig {
	return &CmdConfig{
		Notes:      []string{},      // 备注内容切片
		Examples:   []ExampleInfo{}, // 示例信息切片
		UseChinese: false,           // 是否使用中文帮助信息
		NoFgExit:   false,           // 禁用内置标志自动退出
		Completion: false,           // 控制是否启用自动补全功能
	}
}
