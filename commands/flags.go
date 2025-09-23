// Package commands 提供了主命令的初始化和配置功能。
// 该文件负责设置 fck 工具的版本信息、帮助文档、Logo 显示等全局配置。
package commands

import (
	"gitee.com/MM-Q/fck/commands/internal/types"
	"gitee.com/MM-Q/qflag"
	"gitee.com/MM-Q/verman"
)

// InitMainCmd 初始化主命令
func InItMainCmd() { // 获取版本信息
	qflag.SetVersion(verman.V.Version())                                 // 设置版本信息
	qflag.SetUseChinese(true)                                            // 启用中文帮助信息
	qflag.SetDescription("多功能文件处理工具集, 提供文件哈希计算、大小统计、查找和校验等实用功能")         // 设置命令行描述
	qflag.AddNote("各子命令有独立帮助文档，可通过-h参数查看, 例如 'fck <子命令> -h' 查看各子命令详细帮助") // 设置命令行note
	qflag.AddNote("所有路径参数支持Windows和Unix风格")                              // 添加命令行note
	qflag.SetLogoText(types.FckHelpLogo)                                 // 设置命令行logo
	qflag.SetEnableCompletion(true)                                      // 启用自动补全
}
