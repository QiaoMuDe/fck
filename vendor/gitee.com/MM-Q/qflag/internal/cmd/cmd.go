// Package cmd 提供命令实现和命令管理功能
//
// cmd 包实现了 types.Command 接口, 提供了完整的命令行命令功能。
// 主要组件:
//   - Cmd: 命令结构体, 实现了所有命令相关接口
//   - 命令生命周期管理
//   - 标志和子命令管理
//   - 解析和执行功能
//
// 特性:
//   - 线程安全的命令结构
//   - 支持嵌套子命令
//   - 灵活的配置选项
//   - 完整的帮助系统
package cmd

import (
	"fmt"
	"sync"

	"gitee.com/MM-Q/qflag/internal/help"
	"gitee.com/MM-Q/qflag/internal/parser"
	"gitee.com/MM-Q/qflag/internal/registry"
	"gitee.com/MM-Q/qflag/internal/types"
)

// Cmd 是一个命令结构体, 实现了 types.Command 接口
//
// Cmd 提供了完整的命令行命令实现, 支持标志管理、子命令、
// 参数解析和执行等功能。使用读写锁保证并发安全。
//
// 字段说明:
//   - mu: 读写锁, 保护所有字段的并发访问
//   - longName/shortName: 命令的长名称和短名称
//   - desc: 命令的描述信息
//   - config: 命令的配置选项
//   - hidden: 是否隐藏命令, 隐藏的命令不会显示在帮助信息中
//   - disableFlagParsing: 是否禁用标志解析, 禁用后所有参数都作为位置参数处理
//   - flagRegistry: 标志注册器, 管理命令的所有标志
//   - cmdRegistry: 子命令注册器, 管理所有子命令
//   - args: 命令行参数列表
//   - parsed: 标记是否已解析
//   - parseOnce: 确保解析只执行一次
//   - runFunc: 命令的执行函数
//   - parser: 命令的解析器
//   - parent: 父命令引用, 用于构建命令树
//
// 线程安全:
//   - 所有公共方法都使用读写锁保护
//   - 支持并发读取和独占写入
//   - 解析操作使用sync.Once确保只执行一次
type Cmd struct {
	mu sync.RWMutex // 读写锁, 用于保护命令的并发访问

	longName           string           // 长命令名
	shortName          string           // 短命令名
	desc               string           // 命令描述
	config             *types.CmdConfig // 命令配置
	hidden             bool             // 是否隐藏命令
	disableFlagParsing bool             // 是否禁用标志解析

	flagRegistry types.FlagRegistry // 标志注册器
	cmdRegistry  types.CmdRegistry  // 子命令注册器

	args      []string  // 命令行参数
	parsed    bool      // 是否已解析
	parseOnce sync.Once // 解析一次标志

	runFunc func(types.Command) error // 运行函数

	parser types.Parser // 解析器

	parent *Cmd // 父命令, 默认为 nil
}

// NewCmd 创建新的命令实例
//
// 参数:
//   - longName: 命令的长名称
//   - shortName: 命令的短名称
//   - errorHandling: 错误处理策略
//
// 返回值:
//   - *Cmd: 初始化完成的命令实例
//
// 功能说明:
//   - 创建命令并初始化基本字段
//   - 创建标志和子命令注册器
//   - 设置默认解析器
//   - 初始化配置选项
func NewCmd(longName, shortName string, errorHandling types.ErrorHandling) *Cmd {
	return &Cmd{
		longName:     longName,
		shortName:    shortName,
		config:       types.NewCmdConfig(),
		flagRegistry: registry.NewFlagRegistry(),
		cmdRegistry:  registry.NewCmdRegistry(),
		args:         []string{},
		parsed:       false,
		parser:       parser.NewDefaultParser(errorHandling),
	}
}

// Name 获取命令名称
//
// 返回值:
//   - string: 命令的名称, 优先返回长名称
//
// 功能说明:
//   - 实现types.Command接口
//   - 优先返回长名称, 为空时返回短名称
//   - 用作命令的主要标识符
func (c *Cmd) Name() string {
	if c.longName != "" {
		return c.longName
	}
	return c.shortName
}

// LongName 获取命令长名称
//
// 返回值:
//   - string: 命令的长名称
//
// 功能说明:
//   - 实现types.Command接口
//   - 线程安全地访问长名称
//   - 用于命令的完整标识
func (c *Cmd) LongName() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.longName
}

// ShortName 获取命令短名称
//
// 返回值:
//   - string: 命令的短名称
//
// 功能说明:
//   - 实现types.Command接口
//   - 线程安全地访问短名称
//   - 用于命令的简短标识
func (c *Cmd) ShortName() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.shortName
}

// Desc 获取命令描述
//
// 返回值:
//   - string: 命令的描述信息
//
// 功能说明:
//   - 实现types.Command接口
//   - 线程安全地访问描述信息
//   - 用于帮助信息显示
func (c *Cmd) Desc() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.desc
}

// AddFlag 添加标志到命令
//
// 参数:
//   - f: 要添加的标志
//
// 返回值:
//   - error: 添加失败时返回错误
//
// 功能说明:
//   - 实现types.Command接口
//   - 将标志注册到命令的标志注册器
//   - 支持并发安全的添加操作
//
// 错误情况:
//   - 标志为nil: 返回错误
//   - 标志名称冲突: 返回错误
func (c *Cmd) AddFlag(f types.Flag) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if f == nil {
		return fmt.Errorf("nil flag in '%s'", c.Name())
	}

	return c.flagRegistry.Register(f)
}

// AddFlags 添加多个标志到命令
//
// 参数:
//   - flags: 要添加的标志列表
//
// 返回值:
//   - error: 添加失败时返回错误
//
// 功能说明:
//   - 实现types.Command接口
//   - 批量添加多个标志
//   - 支持并发安全的添加操作
//
// 错误情况:
//   - 标志为nil: 返回错误
//   - 标志名称冲突: 返回错误
func (c *Cmd) AddFlags(flags ...types.Flag) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, flag := range flags {
		if flag == nil {
			return fmt.Errorf("nil flag in '%s'", c.Name())
		}

		if err := c.flagRegistry.Register(flag); err != nil {
			return fmt.Errorf("register flag failed in '%s': %w", c.Name(), err)
		}
	}
	return nil
}

// AddFlagsFrom 从切片添加多个标志
//
// 参数:
//   - flags: 标志切片
//
// 返回值:
//   - error: 添加失败时返回错误
//
// 功能说明:
//   - 实现types.Command接口
//   - 从切片中添加多个标志
//   - 空切片直接返回成功
//   - 内部调用AddFlags实现
func (c *Cmd) AddFlagsFrom(flags []types.Flag) error {
	if len(flags) == 0 {
		return nil
	}
	return c.AddFlags(flags...)
}

// GetFlag 根据名称获取标志
//
// 参数:
//   - name: 标志名称
//
// 返回值:
//   - types.Flag: 找到的标志
//   - bool: 是否找到, true表示找到
//
// 功能说明:
//   - 实现types.Command接口
//   - 从标志注册器中查找标志
//   - 支持并发安全的查找操作
func (c *Cmd) GetFlag(name string) (types.Flag, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.flagRegistry.Get(name)
}

// Flags 获取所有标志
//
// 返回值:
//   - []types.Flag: 所有标志的切片副本
//
// 功能说明:
//   - 实现types.Command接口
//   - 返回所有注册的标志
//   - 创建副本避免外部修改
//   - 支持并发安全的访问
func (c *Cmd) Flags() []types.Flag {
	c.mu.RLock()
	defer c.mu.RUnlock()

	flags := c.flagRegistry.List()
	if len(flags) == 0 {
		return []types.Flag{}
	}
	result := make([]types.Flag, len(flags))
	copy(result, flags)
	return result
}

// AddSubCmds 添加子命令到命令
//
// 参数:
//   - cmds: 要添加的子命令列表
//
// 返回值:
//   - error: 添加失败时返回错误
//
// 功能说明:
//   - 实现types.Command接口
//   - 批量添加多个子命令
//   - 自动设置父子关系
//   - 支持并发安全的添加操作
//
// 错误情况:
//   - 子命令为nil: 返回错误
//   - 子命令类型错误: 返回错误
//   - 子命令名称冲突: 返回错误
func (c *Cmd) AddSubCmds(cmds ...types.Command) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, cmd := range cmds {
		if cmd == nil {
			return fmt.Errorf("nil subcommand in '%s'", c.Name())
		}

		// 检查子命令是否为 *Cmd 类型
		if subCmd, ok := cmd.(*Cmd); ok {
			subCmd.parent = c // 设置子命令的父命令为当前命令
		} else {
			return fmt.Errorf("invalid subcommand type in '%s'", c.Name())
		}

		// 注册子命令
		if err := c.cmdRegistry.Register(cmd); err != nil {
			return fmt.Errorf("register subcommand '%s' failed in '%s': %w", cmd.Name(), c.Name(), err)
		}
	}
	return nil
}

// AddSubCmdFrom 从切片添加子命令
//
// 参数:
//   - cmds: 子命令切片
//
// 返回值:
//   - error: 添加失败时返回错误
//
// 功能说明:
//   - 实现types.Command接口
//   - 从切片中添加子命令
//   - 空切片直接返回成功
//   - 内部调用AddSubCmds实现
func (c *Cmd) AddSubCmdFrom(cmds []types.Command) error {
	if len(cmds) == 0 {
		return nil
	}
	return c.AddSubCmds(cmds...)
}

// GetSubCmd 根据名称获取子命令
//
// 参数:
//   - name: 子命令名称
//
// 返回值:
//   - types.Command: 找到的子命令
//   - bool: 是否找到, true表示找到
//
// 功能说明:
//   - 实现types.Command接口
//   - 从子命令注册器中查找
//   - 支持并发安全的查找操作
func (c *Cmd) GetSubCmd(name string) (types.Command, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.cmdRegistry.Get(name)
}

// SubCmds 获取所有子命令
//
// 返回值:
//   - []types.Command: 所有子命令的切片副本
//
// 功能说明:
//   - 实现types.Command接口
//   - 返回所有注册的子命令
//   - 创建副本避免外部修改
//   - 支持并发安全的访问
func (c *Cmd) SubCmds() []types.Command {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cmds := c.cmdRegistry.List()
	if len(cmds) == 0 {
		return []types.Command{}
	}
	result := make([]types.Command, len(cmds))
	copy(result, cmds)
	return result
}

// HasSubCmd 检查是否存在指定名称的子命令
//
// 参数:
//   - name: 子命令名称
//
// 返回值:
//   - bool: 是否存在, true表示存在
//
// 功能说明:
//   - 实现types.Command接口
//   - 快速检查子命令存在性
//   - 支持并发安全的检查
func (c *Cmd) HasSubCmd(name string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.cmdRegistry.Has(name)
}

// Parse 解析命令行参数 (可重复解析)
//
// 参数:
//   - args: 命令行参数列表
//
// 返回值:
//   - error: 解析失败时返回错误
//
// 功能说明:
//   - 实现types.Command接口
//   - 可以重复调用，会覆盖之前的解析结果
//   - 调用解析器的Parse方法
//   - 递归解析所有子命令
//
// 注意事项:
//   - 重复解析会覆盖之前的解析结果
//   - 如果需要确保只解析一次，请使用 ParseOnce
func (c *Cmd) Parse(args []string) error {
	return c.parser.Parse(c, args)
}

// ParseOnce 解析命令行参数 (只解析一次)
//
// 参数:
//   - args: 命令行参数列表
//
// 返回值:
//   - error: 解析失败时返回错误
//
// 功能说明:
//   - 使用sync.Once确保只解析一次
//   - 重复执行无错误、仅首次执行解析
//   - 调用解析器的Parse方法
//   - 递归解析所有子命令
//
// 注意事项:
//   - 如果需要重复解析，请使用 Parse 方法
//   - 建议在普通场景使用此方法，避免误用
func (c *Cmd) ParseOnce(args []string) error {
	var err error
	c.parseOnce.Do(func() {
		err = c.parser.Parse(c, args)
	})
	return err
}

// ParseAndRoute 解析并路由执行命令 (可重复解析)
//
// 参数:
//   - args: 命令行参数列表
//
// 返回值:
//   - error: 解析或执行失败时返回错误
//
// 功能说明:
//   - 实现types.Command接口
//   - 可以重复调用，会覆盖之前的解析结果
//   - 调用解析器的ParseAndRoute方法
//   - 完整的解析和执行流程
//
// 注意事项:
//   - 重复解析会覆盖之前的解析结果
//   - 如果需要确保只解析一次，请使用 ParseAndRouteOnce
func (c *Cmd) ParseAndRoute(args []string) error {
	return c.parser.ParseAndRoute(c, args)
}

// ParseAndRouteOnce 解析并路由执行命令 (只解析一次)
//
// 参数:
//   - args: 命令行参数列表
//
// 返回值:
//   - error: 解析或执行失败时返回错误
//
// 功能说明:
//   - 使用sync.Once确保只执行一次
//   - 重复执行无错误、仅首次执行解析
//   - 调用解析器的ParseAndRoute方法
//   - 完整的解析和执行流程
//
// 注意事项:
//   - 如果需要重复解析，请使用 ParseAndRoute 方法
//   - 建议在普通场景使用此方法，避免误用
func (c *Cmd) ParseAndRouteOnce(args []string) error {
	var err error
	c.parseOnce.Do(func() {
		err = c.parser.ParseAndRoute(c, args)
	})
	return err
}

// ParseOnly 仅解析当前命令, 不递归解析子命令 (可重复解析)
//
// 参数:
//   - args: 命令行参数列表
//
// 返回值:
//   - error: 解析失败时返回错误
//
// 功能说明:
//   - 实现types.Command接口
//   - 可以重复调用，会覆盖之前的解析结果
//   - 调用解析器的ParseOnly方法
//   - 不处理子命令解析
//
// 注意事项:
//   - 重复解析会覆盖之前的解析结果
//   - 如果需要确保只解析一次，请使用 ParseOnlyOnce
func (c *Cmd) ParseOnly(args []string) error {
	return c.parser.ParseOnly(c, args)
}

// ParseOnlyOnce 仅解析当前命令, 不递归解析子命令 (只解析一次)
//
// 参数:
//   - args: 命令行参数列表
//
// 返回值:
//   - error: 解析失败时返回错误
//
// 功能说明:
//   - 使用sync.Once确保只执行一次
//   - 重复执行无错误、仅首次执行解析
//   - 调用解析器的ParseOnly方法
//   - 不处理子命令解析
//
// 注意事项:
//   - 如果需要重复解析，请使用 ParseOnly 方法
//   - 建议在普通场景使用此方法，避免误用
func (c *Cmd) ParseOnlyOnce(args []string) error {
	var err error
	c.parseOnce.Do(func() {
		err = c.parser.ParseOnly(c, args)
	})
	return err
}

// IsParsed 检查命令是否已解析
//
// 返回值:
//   - bool: 是否已解析, true表示已解析
//
// 功能说明:
//   - 实现types.Command接口
//   - 线程安全地检查解析状态
//   - 用于避免重复解析
func (c *Cmd) IsParsed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.parsed
}

// Args 获取命令行参数
//
// 返回值:
//   - []string: 命令行参数的副本
//
// 功能说明:
//   - 实现types.Command接口
//   - 返回解析后的参数列表
//   - 创建副本避免外部修改
//   - 支持并发安全的访问
func (c *Cmd) Args() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if len(c.args) == 0 {
		return []string{}
	}
	result := make([]string, len(c.args))
	copy(result, c.args)
	return result
}

// Arg 获取指定索引的命令行参数
//
// 参数:
//   - index: 命令行参数的索引
//
// 返回值:
//   - string: 对应索引的命令行参数值
//
// 注意:
//   - 索引从 0 开始计数
//   - 如果索引超出范围, 返回空字符串
func (c *Cmd) Arg(index int) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if index >= 0 && index < len(c.args) {
		return c.args[index]
	}
	return ""
}

// NArg 获取命令行参数数量
//
// 返回值:
//   - int: 参数数量
//
// 功能说明:
//   - 实现types.Command接口
//   - 线程安全地获取参数数量
//   - 用于参数范围检查
func (c *Cmd) NArg() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.args)
}

// Run 执行命令
//
// 返回值:
//   - error: 执行失败时返回错误
//
// 功能说明:
//   - 实现types.Command接口
//   - 检查解析状态和运行函数
//   - 调用设置的运行函数
//   - 支持并发安全的执行
//
// 错误情况:
//   - 未解析: 返回解析错误
//   - 无运行函数: 返回运行函数错误
func (c *Cmd) Run() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.parsed {
		return fmt.Errorf("'%s' not parsed", c.Name())
	}

	if c.runFunc == nil {
		return fmt.Errorf("no run function in '%s'", c.Name())
	}

	return c.runFunc(c)
}

// HasRunFunc 检查是否设置了运行函数
//
// 返回值:
//   - bool: 是否设置了运行函数, true表示已设置
//
// 功能说明:
//   - 实现types.Command接口
//   - 线程安全地检查运行函数
//   - 用于执行前的状态检查
func (c *Cmd) HasRunFunc() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.runFunc != nil
}

// Help 获取帮助信息
//
// 返回值:
//   - string: 格式化的帮助信息
//
// 功能说明:
//   - 实现types.Command接口
//   - 使用help包生成帮助信息
//   - 包含标志、子命令和示例
//   - 支持并发安全的访问
func (c *Cmd) Help() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return help.GenHelp(c)
}

// PrintHelp 打印帮助信息
//
// 功能说明:
//   - 实现types.Command接口
//   - 直接输出帮助信息到控制台
//   - 使用标准fmt包输出
//   - 支持并发安全的访问
func (c *Cmd) PrintHelp() {
	fmt.Println(help.GenHelp(c))
}

// IsHidden 检查命令是否隐藏
//
// 返回值:
//   - bool: 是否隐藏, true表示隐藏
//
// 功能说明:
//   - 实现types.Command接口
//   - 线程安全地检查隐藏状态
//   - 用于帮助信息过滤
func (c *Cmd) IsHidden() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.hidden
}

// IsDisableFlagParsing 检查是否禁用标志解析
//
// 返回值:
//   - bool: 是否禁用标志解析, true表示禁用
//
// 功能说明:
//   - 实现types.Command接口
//   - 线程安全地检查标志解析状态
//   - 用于解析器决定是否解析标志
func (c *Cmd) IsDisableFlagParsing() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.disableFlagParsing
}

// FlagRegistry 获取标志注册器
//
// 返回值:
//   - types.FlagRegistry: 标志注册器接口
//
// 功能说明:
//   - 实现types.Command接口
//   - 返回命令的标志注册器
//   - 用于直接操作标志注册
//   - 支持并发安全的访问
func (c *Cmd) FlagRegistry() types.FlagRegistry {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.flagRegistry
}

// CmdRegistry 获取子命令注册器
//
// 返回值:
//   - types.CmdRegistry: 子命令注册器接口
//
// 功能说明:
//   - 实现types.Command接口
//   - 返回命令的子命令注册器
//   - 用于直接操作子命令注册
//   - 支持并发安全的访问
func (c *Cmd) CmdRegistry() types.CmdRegistry {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cmdRegistry
}

// IsRootCmd 检查是否为根命令
//
// 返回值:
//   - bool: 是否为根命令, true表示是根命令
//
// 功能说明:
//   - 实现types.Command接口
//   - 通过检查父命令判断
//   - 根命令没有父命令
//   - 支持并发安全的检查
func (c *Cmd) IsRootCmd() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.parent == nil
}

// Path 获取命令路径
//
// 返回值:
//   - string: 完整的命令路径
//
// 功能说明:
//   - 实现types.Command接口
//   - 递归构建完整路径
//   - 格式: 父路径 + 空格 + 命令名
//   - 根命令直接返回名称
//   - 用于帮助信息和错误显示
//   - 用于帮助信息和错误显示
func (c *Cmd) Path() string {
	// 根命令直接返回名称
	if c.parent == nil {
		return c.Name()
	}

	// 子命令路径为: 父命令路径 + 空格 + 子命令名称
	return c.parent.Path() + " " + c.Name()
}
