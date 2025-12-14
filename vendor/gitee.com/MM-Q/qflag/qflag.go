// Package qflag 提供对标准库flag的封装, 自动实现长短标志, 并默认绑定-h/--help标志打印帮助信息。
// 用户可通过Cmd.Help字段自定义帮助内容, 支持直接赋值字符串或从文件加载。
// 该包是一个功能强大的命令行参数解析库, 支持子命令、多种数据类型标志、参数验证等高级特性。
package qflag

import (
	"flag"
	"os"
	"path/filepath"

	"gitee.com/MM-Q/qflag/cmd"
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
	Root = cmd.NewCmd(cmdName, "", flag.ExitOnError)
}

// ================================================================================
// 为了保持向后兼容性, 保留最常用的全局函数
// 用户应该优先使用 Root 全局命令实例访问所有方法
// ================================================================================

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
