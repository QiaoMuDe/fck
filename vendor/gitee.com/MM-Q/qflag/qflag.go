// package qflag 提供对标准库flag的封装, 自动实现长短标志, 并默认绑定-h/--help标志打印帮助信息。
// 用户可通过Cmd.Help字段自定义帮助内容, 支持直接赋值字符串或从文件加载。
// 该包是一个功能强大的命令行参数解析库, 支持子命令、多种数据类型标志、参数验证等高级特性。
package qflag

import (
	"os"
	"path/filepath"

	"gitee.com/MM-Q/qflag/internal/types"
)

/*
项目地址: https://gitee.com/MM-Q/qflag
*/

// Root 全局根命令实例, 提供对全局标志和子命令的访问
// 用户可以通过 qflag.Root.String() 这样的方式直接创建全局标志
// 这是访问命令行功能的主要入口点, 推荐优先使用
var Root *Cmd

// init 包初始化函数，直接创建全局根命令实例
func init() {
	// 使用一致的命令名生成逻辑
	cmdName := "myapp"
	if len(os.Args) > 0 {
		cmdName = filepath.Base(os.Args[0])
	}

	// 直接创建全局根命令实例
	Root = NewCmd(cmdName, "", ExitOnError)
}

// Parse 解析所有命令行参数, 包括根命令和所有子命令的标志参数
//
// 返回：
//   - error: 解析过程中遇到的错误, 若成功则为 nil
func Parse() error {
	return Root.Parse(os.Args[1:])
}

// ParseFlagsOnly 解析根命令的所有标志参数, 不包括子命令
//
// 返回：
//   - error: 解析过程中遇到的错误, 若成功则为 nil
func ParseFlagsOnly() error {
	return Root.ParseFlagsOnly(os.Args[1:])
}

// ParseAndRoute 解析参数并自动路由执行子命令
// 这是推荐使用的命令行参数处理方式，会自动处理子命令路由
//
// 返回：
//   - error: 执行过程中遇到的错误, 若成功则为 nil
func ParseAndRoute() error {
	return Root.ParseAndRoute(os.Args[1:])
}

// ApplyConfig 批量设置根命令配置
// 通过传入一个CmdConfig结构体来一次性设置多个配置项
// 这是全局函数，直接操作全局根命令实例
//
// 参数:
//   - config: 包含所有配置项的CmdConfig结构体
func ApplyConfig(config types.CmdConfig) {
	Root.ApplyConfig(config)
}

// AddSubCmd 向根命令添加一个或多个子命令
// 这是全局函数，直接操作全局根命令实例
//
// 此方法会对所有子命令进行完整性验证，包括名称冲突检查、循环依赖检测等。
// 所有验证通过后，子命令将被注册到根命令的子命令映射表和列表中。
// 操作过程中会自动设置子命令的父命令引用，确保命令树结构的完整性。
//
// 参数:
//   - subCmds: 要添加的子命令实例指针，支持传入多个子命令进行批量添加
//
// 返回值:
//   - error: 添加过程中的错误信息。如果任何子命令验证失败，将返回包含所有错误详情的聚合错误；
//     如果所有子命令成功添加，返回 nil
func AddSubCmd(subCmds ...*Cmd) error {
	return Root.AddSubCmd(subCmds...)
}

// AddSubCmds 向根命令添加子命令切片的便捷方法
// 这是全局函数，直接操作全局根命令实例
//
// 此方法是 AddSubCmd 的便捷包装，专门用于处理子命令切片。
// 内部直接调用 AddSubCmd 方法，具有相同的验证逻辑和并发安全特性。
//
// 参数:
//   - subCmds: 子命令切片，包含要添加的所有子命令实例指针
//
// 返回值:
//   - error: 添加过程中的错误信息，与 AddSubCmd 返回的错误类型相同
func AddSubCmds(subCmds []*Cmd) error {
	return Root.AddSubCmds(subCmds)
}
