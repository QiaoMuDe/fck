package qflag

import (
	"os"
	"path/filepath"
)

// Root 全局根命令实例, 提供对全局标志和子命令的访问
// 用户可以通过 qflag.Root.String() 这样的方式直接创建全局标志
// 这是访问命令行功能的主要入口点, 推荐优先使用
var Root *Cmd

// init 包初始化函数, 直接创建全局根命令实例
func init() {
	// 使用一致的命令名生成逻辑
	cmdName := "myapp"
	if len(os.Args) > 0 {
		cmdName = filepath.Base(os.Args[0])
	}

	// 直接创建全局根命令实例
	Root = NewCmd(cmdName, "", ExitOnError)
}

// Parse 解析命令行参数
//
// 返回值:
//   - error: 解析失败时返回错误
//
// 功能说明:
//   - 使用全局根命令解析命令行参数
//   - 可以重复调用, 会覆盖之前的解析结果
//   - 递归解析所有子命令
//
// 注意事项:
//   - 如果需要确保只解析一次, 请使用 ParseOnce
func Parse() error {
	return Root.Parse(os.Args[1:])
}

// ParseOnce 解析命令行参数（只解析一次）
//
// 返回值:
//   - error: 解析失败时返回错误
//
// 功能说明:
//   - 使用全局根命令解析命令行参数
//   - 使用ParseOnce确保只解析一次
//   - 重复执行无错误、仅首次执行解析
//   - 递归解析所有子命令
//
// 注意事项:
//   - 建议在普通场景使用此方法, 避免误用
//   - 如果需要重复解析, 请使用 Parse
func ParseOnce() error {
	return Root.ParseOnce(os.Args[1:])
}

// ParseOnly 仅解析当前命令, 不递归解析子命令
//
// 返回值:
//   - error: 解析失败时返回错误
//
// 功能说明:
//   - 使用全局根命令解析命令行参数
//   - 可以重复调用, 会覆盖之前的解析结果
//   - 不处理子命令解析
//
// 注意事项:
//   - 如果需要确保只解析一次, 请使用 ParseOnlyOnce
func ParseOnly() error {
	return Root.ParseOnly(os.Args[1:])
}

// ParseOnlyOnce 仅解析当前命令, 不递归解析子命令（只解析一次）
//
// 返回值:
//   - error: 解析失败时返回错误
//
// 功能说明:
//   - 使用全局根命令解析命令行参数
//   - 使用ParseOnlyOnce确保只解析一次
//   - 重复执行无错误、仅首次执行解析
//   - 不处理子命令解析
//
// 注意事项:
//   - 建议在普通场景使用此方法, 避免误用
//   - 如果需要重复解析, 请使用 ParseOnly
func ParseOnlyOnce() error {
	return Root.ParseOnlyOnce(os.Args[1:])
}

// ParseAndRoute 解析并路由执行命令
//
// 返回值:
//   - error: 解析或执行失败时返回错误
//
// 功能说明:
//   - 使用全局根命令解析命令行参数
//   - 可以重复调用, 会覆盖之前的解析结果
//   - 完整的解析和执行流程
//
// 注意事项:
//   - 如果需要确保只解析一次, 请使用 ParseAndRouteOnce
func ParseAndRoute() error {
	return Root.ParseAndRoute(os.Args[1:])
}

// ParseAndRouteOnce 解析并路由执行命令（只解析一次）
//
// 返回值:
//   - error: 解析或执行失败时返回错误
//
// 功能说明:
//   - 使用全局根命令解析命令行参数
//   - 使用ParseAndRouteOnce确保只解析一次
//   - 重复执行无错误、仅首次执行解析
//   - 完整的解析和执行流程
//
// 注意事项:
//   - 建议在普通场景使用此方法, 避免误用
//   - 如果需要重复解析, 请使用 ParseAndRoute
func ParseAndRouteOnce() error {
	return Root.ParseAndRouteOnce(os.Args[1:])
}

// AddSubCmds 添加子命令到全局根命令
//
// 参数:
//   - cmd: 要添加的子命令实例
//
// 返回值:
//   - error: 添加子命令过程中遇到的错误, 如果没有错误则返回 nil
func AddSubCmds(cmds ...Command) error {
	return Root.AddSubCmds(cmds...)
}

// AddSubCmdFrom 从切片添加子命令
//
// 参数:
//   - cmds: 要添加的子命令实例切片
//
// 返回值:
//   - error: 添加子命令过程中遇到的错误, 如果没有错误则返回 nil
func AddSubCmdFrom(cmds []Command) error {
	return Root.AddSubCmdFrom(cmds)
}

// ApplyOpts 应用选项到全局根命令
//
// 参数:
//   - opts: 要应用的选项结构体实例
//
// 返回值:
//   - error: 应用选项过程中遇到的错误, 如果没有错误则返回 nil
//
// 功能说明:
//   - 将选项结构体的所有属性应用到全局根命令实例
//   - 支持部分配置（未设置的属性不会被修改）
//   - 使用写锁保护并发安全
func ApplyOpts(opts *CmdOpts) error {
	return Root.ApplyOpts(opts)
}
