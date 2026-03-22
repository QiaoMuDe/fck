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

	longName  string           // 长命令名
	shortName string           // 短命令名
	desc      string           // 命令描述
	config    *types.CmdConfig // 命令配置

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
//   - 标志为nil: 返回INVALID_FLAG错误
//   - 标志名称冲突: 返回FLAG_ALREADY_EXISTS错误
func (c *Cmd) AddFlag(f types.Flag) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if f == nil {
		return types.NewError("INVALID_FLAG", "flag cannot be nil", nil)
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
//   - 标志为nil: 返回INVALID_FLAG错误
//   - 标志名称冲突: 返回FLAG_ALREADY_EXISTS错误
func (c *Cmd) AddFlags(flags ...types.Flag) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, flag := range flags {
		if flag == nil {
			return types.NewError("INVALID_FLAG", "flag cannot be nil", nil)
		}

		if err := c.flagRegistry.Register(flag); err != nil {
			return err
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
//   - 子命令为nil: 返回INVALID_COMMAND错误
//   - 子命令类型错误: 返回INVALID_COMMAND_TYPE错误
//   - 子命令名称冲突: 返回COMMAND_ALREADY_EXISTS错误
func (c *Cmd) AddSubCmds(cmds ...types.Command) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, cmd := range cmds {
		if cmd == nil {
			return types.NewError("INVALID_COMMAND", "cmd cannot be nil", nil)
		}

		// 检查子命令是否为 *Cmd 类型
		if subCmd, ok := cmd.(*Cmd); ok {
			subCmd.parent = c // 设置子命令的父命令为当前命令
		} else {
			return types.NewError("INVALID_COMMAND_TYPE", "cmd must be *Cmd type", nil)
		}

		// 注册子命令
		if err := c.cmdRegistry.Register(cmd); err != nil {
			return err
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

// Parse 解析命令行参数（可重复解析）
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

// ParseOnce 解析命令行参数（只解析一次）
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

// ParseAndRoute 解析并路由执行命令（可重复解析）
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

// ParseAndRouteOnce 解析并路由执行命令（只解析一次）
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

// ParseOnly 仅解析当前命令, 不递归解析子命令（可重复解析）
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

// ParseOnlyOnce 仅解析当前命令, 不递归解析子命令（只解析一次）
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
		return fmt.Errorf("cmd must be parsed before execution")
	}

	if c.runFunc == nil {
		return fmt.Errorf("no run function set")
	}

	return c.runFunc(c)
}

// SetRun 设置命令的运行函数
//
// 参数:
//   - fn: 运行函数, 接收命令并返回错误
//
// 功能说明:
//   - 实现types.Command接口
//   - 设置命令的执行逻辑
//   - 支持并发安全的设置
//   - 可多次设置, 最后一次生效
func (c *Cmd) SetRun(fn func(types.Command) error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.runFunc = fn
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

// SetDesc 设置命令描述
//
// 参数:
//   - desc: 命令描述信息
//
// 功能说明:
//   - 实现types.Command接口
//   - 设置命令的功能描述
//   - 用于帮助信息显示
//   - 支持并发安全的设置
func (c *Cmd) SetDesc(desc string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.desc = desc
}

// SetVersion 设置命令版本
//
// 参数:
//   - version: 版本字符串
//
// 功能说明:
//   - 设置命令的版本信息
//   - 存储在配置中
//   - 用于版本显示和帮助信息
//   - 支持并发安全的设置
//   - 只有根命令才能设置版本信息
func (c *Cmd) SetVersion(version string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 只有根命令才能设置版本信息
	if c.parent != nil {
		return
	}

	c.config.Version = version
}

// SetChinese 设置是否使用中文
//
// 参数:
//   - useChinese: 是否使用中文
//
// 功能说明:
//   - 设置帮助信息的语言
//   - 影响错误消息和提示
//   - 存储在配置中
//   - 支持并发安全的设置
func (c *Cmd) SetChinese(useChinese bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.config.UseChinese = useChinese
}

// SetCompletion 设置是否启用自动补全标志
//
// 参数:
//   - enable: 是否启用自动补全标志
//
// 功能说明:
//   - 控制是否注册 --completion 标志
//   - 默认为 false，不启用自动补全
//   - 存储在配置中
//   - 支持并发安全的设置
func (c *Cmd) SetCompletion(enable bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.config.Completion = enable
}

// SetEnvPrefix 设置环境变量前缀
//
// 参数:
//   - prefix: 环境变量前缀
//
// 功能说明:
//   - 设置环境变量的前缀
//   - 自动添加下划线后缀
//   - 空字符串表示不使用前缀
//   - 支持并发安全的设置
func (c *Cmd) SetEnvPrefix(prefix string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if prefix != "" {
		c.config.EnvPrefix = prefix + "_"
	} else {
		c.config.EnvPrefix = ""
	}
}

// SetUsageSyntax 设置使用语法
//
// 参数:
//   - syntax: 使用语法字符串
//
// 功能说明:
//   - 设置命令的使用语法
//   - 用于帮助信息显示
//   - 存储在配置中
//   - 支持并发安全的设置
func (c *Cmd) SetUsageSyntax(syntax string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.config.UsageSyntax = syntax
}

// AddExample 添加单个示例
//
// 参数:
//   - title: 示例标题
//   - cmd: 示例命令
//
// 功能说明:
//   - 添加命令使用示例
//   - 用于帮助信息显示
//   - 存储在配置中
//   - 支持并发安全的添加
func (c *Cmd) AddExample(title, cmd string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.config.Example[title] = cmd
}

// AddExamples 批量添加示例
//
// 参数:
//   - examples: 示例映射, 标题为键, 命令为值
//
// 功能说明:
//   - 批量添加多个示例
//   - 空映射直接返回
//   - 覆盖同名的示例
//   - 支持并发安全的添加
func (c *Cmd) AddExamples(examples map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(examples) == 0 {
		return
	}

	for title, cmd := range examples {
		c.config.Example[title] = cmd
	}
}

// AddNote 添加单个注释
//
// 参数:
//   - note: 注释内容
//
// 功能说明:
//   - 添加命令的额外说明
//   - 用于帮助信息显示
//   - 空字符串被忽略
//   - 支持并发安全的添加
func (c *Cmd) AddNote(note string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if note == "" {
		return
	}

	c.config.Notes = append(c.config.Notes, note)
}

// AddNotes 批量添加注释
//
// 参数:
//   - notes: 注释切片
//
// 功能说明:
//   - 批量添加多个注释
//   - 空切片直接返回
//   - 空字符串被忽略
//   - 支持并发安全的添加
func (c *Cmd) AddNotes(notes []string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(notes) == 0 {
		return
	}

	c.config.Notes = append(c.config.Notes, notes...)
}

// SetLogoText 设置Logo文本
//
// 参数:
//   - logo: Logo文本内容
//
// 功能说明:
//   - 设置命令的Logo
//   - 用于帮助信息显示
//   - 存储在配置中
//   - 支持并发安全的设置
func (c *Cmd) SetLogoText(logo string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.config.LogoText = logo
}

// Config 获取命令配置
//
// 返回值:
//   - *types.CmdConfig: 命令配置的副本
//
// 功能说明:
//   - 实现types.Command接口
//   - 返回命令的配置对象
//   - 注意: 返回的是副本, 修改不会影响命令
//   - 支持并发安全的访问
func (c *Cmd) Config() *types.CmdConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var mutexGroups []types.MutexGroup
	if len(c.config.MutexGroups) > 0 {
		mutexGroups = make([]types.MutexGroup, len(c.config.MutexGroups))
		copy(mutexGroups, c.config.MutexGroups)
	}

	var requiredGroups []types.RequiredGroup
	if len(c.config.RequiredGroups) > 0 {
		requiredGroups = make([]types.RequiredGroup, len(c.config.RequiredGroups))
		copy(requiredGroups, c.config.RequiredGroups)
	}

	return &types.CmdConfig{
		Version:        c.config.Version,
		UseChinese:     c.config.UseChinese,
		EnvPrefix:      c.config.EnvPrefix,
		UsageSyntax:    c.config.UsageSyntax,
		Example:        c.config.Example,
		Notes:          c.config.Notes,
		LogoText:       c.config.LogoText,
		MutexGroups:    mutexGroups,
		RequiredGroups: requiredGroups,
		Completion:     c.config.Completion,
	}
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

// SetParser 设置命令的解析器
//
// 参数:
//   - p: 解析器接口实现
//
// 功能说明:
//   - 替换默认的解析器
//   - 允许自定义解析逻辑
//   - nil值会触发panic
//   - 支持并发安全的设置
func (c *Cmd) SetParser(p types.Parser) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if p == nil {
		panic("parser cannot be nil")
	}

	c.parser = p
}

// SetArgs 设置命令行参数
//
// 参数:
//   - args: 命令行参数列表
//
// 功能说明:
//   - 手动设置命令行参数
//   - 通常由解析器调用
//   - 空切片被忽略
//   - 支持并发安全的设置
func (c *Cmd) SetArgs(args []string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(args) == 0 {
		return
	}

	c.args = args
}

// SetParsed 设置解析状态
//
// 参数:
//   - parsed: 解析状态
//
// 功能说明:
//   - 手动设置解析状态
//   - 通常由解析器调用
//   - 影响后续操作的行为
//   - 支持并发安全的设置
func (c *Cmd) SetParsed(parsed bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.parsed = parsed
}

// ApplyOpts 应用选项到命令
//
// 参数:
//   - opts: 命令选项
//
// 返回值:
//   - error: 应用选项失败时返回错误
//
// 功能说明:
//   - 将选项结构体的所有属性应用到当前命令
//   - 支持部分配置（未设置的属性不会被修改）
//   - 使用defer捕获panic, 转换为错误返回
//
// 应用顺序:
//  1. 基本属性（Desc、RunFunc）
//  2. 配置选项（Version、UseChinese、EnvPrefix、UsageSyntax、LogoText）
//  3. 示例和说明（Examples、Notes）
//  4. 互斥组（MutexGroups）
//  5. 子命令（SubCmds）
//
// 错误处理:
//   - 选项为 nil: 返回 INVALID_CMDOPTS 错误
//   - 添加子命令失败: 返回 FAILED_TO_ADD_SUBCMDS 错误
//   - panic: 转换为 PANIC 错误
//
// 线程安全:
//   - 方法内部使用读写锁保护并发访问
//   - 可以安全地在多个 goroutine 中调用
//
// 设计说明:
//   - 调用现有的 SetDesc、SetVersion、AddExamples 等方法
//   - 避免重复代码，降低维护成本
//   - 保持行为一致性，与用户手动调用方法完全一致
//   - 保留方法中的验证、通知等逻辑
func (c *Cmd) ApplyOpts(opts *CmdOpts) error {
	var err error

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
		}
	}()

	// 验证选项
	if opts == nil {
		return types.NewError("INVALID_CMDOPTS", "cmd opts cannot be nil", nil)
	}

	// 1. 设置基本属性 - 调用现有方法
	if opts.Desc != "" {
		c.SetDesc(opts.Desc)
	}
	if opts.RunFunc != nil {
		c.SetRun(opts.RunFunc)
	}

	// 2. 设置配置选项 - 调用现有方法
	if opts.Version != "" {
		c.SetVersion(opts.Version)
	}
	if opts.EnvPrefix != "" {
		c.SetEnvPrefix(opts.EnvPrefix)
	}
	if opts.UsageSyntax != "" {
		c.SetUsageSyntax(opts.UsageSyntax)
	}
	if opts.LogoText != "" {
		c.SetLogoText(opts.LogoText)
	}
	c.SetChinese(opts.UseChinese)
	c.SetCompletion(opts.Completion)

	// 3. 添加示例和说明 - 调用现有方法
	if len(opts.Examples) > 0 {
		c.AddExamples(opts.Examples)
	}
	if len(opts.Notes) > 0 {
		c.AddNotes(opts.Notes)
	}

	// 4. 添加互斥组 - 调用现有方法
	if len(opts.MutexGroups) > 0 {
		for _, group := range opts.MutexGroups {
			if err := c.AddMutexGroup(group.Name, group.Flags, group.AllowNone); err != nil {
				return types.WrapError(err, "FAILED_TO_ADD_MUTEX_GROUP", "failed to add mutex group")
			}
		}
	}

	// 5. 添加必需组 - 调用现有方法
	if len(opts.RequiredGroups) > 0 {
		for _, group := range opts.RequiredGroups {
			if err := c.AddRequiredGroup(group.Name, group.Flags, group.Conditional); err != nil {
				return types.WrapError(err, "FAILED_TO_ADD_REQUIRED_GROUP", "failed to add required group")
			}
		}
	}

	// 6. 添加子命令 - 调用现有方法
	if len(opts.SubCmds) > 0 {
		if err := c.AddSubCmds(opts.SubCmds...); err != nil {
			return types.WrapError(err, "FAILED_TO_ADD_SUBCMDS", "failed to add subcommands")
		}
	}

	return err
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

// AddMutexGroup 添加互斥组到命令
//
// 参数:
//   - name: 互斥组名称, 用于错误提示和标识
//   - flags: 互斥组中的标志名称列表
//   - allowNone: 是否允许一个都不设置
//
// 返回值:
//   - error: 添加失败时返回错误
//
// 功能说明:
//   - 创建新的互斥组并添加到命令配置中
//   - 互斥组中的标志最多只能有一个被设置
//   - 如果 allowNone 为 false, 则必须至少有一个标志被设置
//   - 使用写锁保护并发安全
//
// 注意事项:
//   - 标志名称必须是已注册的标志
//   - 互斥组名称在命令中应该唯一
//   - 如果组名已存在，返回错误
//
// 错误码:
//   - MUTEX_GROUP_ALREADY_EXISTS: 互斥组已存在
//   - EMPTY_MUTEX_GROUP: 互斥组标志列表为空
//   - FLAG_NOT_FOUND: 标志不存在
func (c *Cmd) AddMutexGroup(name string, flags []string, allowNone bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 互斥组名称不能为空
	if name == "" {
		return types.NewError("EMPTY_MUTEX_GROUP_NAME",
			"mutex group name cannot be empty", nil)
	}

	// 检查互斥组名称是否已存在
	for _, group := range c.config.MutexGroups {
		if group.Name == name {
			return types.NewError("MUTEX_GROUP_ALREADY_EXISTS",
				fmt.Sprintf("mutex group '%s' already exists", name), nil)
		}
	}

	// 检查标志是否为空
	if len(flags) == 0 {
		return types.NewError("EMPTY_MUTEX_GROUP",
			"mutex group cannot be empty", nil)
	}

	// 检查标志名称是否存在
	for _, flagName := range flags {
		if _, exists := c.flagRegistry.Get(flagName); !exists {
			return types.NewError("FLAG_NOT_FOUND",
				fmt.Sprintf("flag '%s' not found", flagName), nil)
		}
	}

	group := types.MutexGroup{
		Name:      name,
		Flags:     flags,
		AllowNone: allowNone,
	}

	c.config.MutexGroups = append(c.config.MutexGroups, group)
	return nil
}

// GetMutexGroups 获取命令的所有互斥组
//
// 返回值:
//   - []types.MutexGroup: 互斥组列表的副本
//
// 功能说明:
//   - 返回命令中定义的所有互斥组
//   - 返回副本以防止外部修改内部状态
//   - 使用读锁保护并发安全
func (c *Cmd) GetMutexGroups() []types.MutexGroup {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// 返回副本以防止外部修改
	groups := make([]types.MutexGroup, len(c.config.MutexGroups))
	copy(groups, c.config.MutexGroups)
	return groups
}

// RemoveMutexGroup 移除指定名称的互斥组
//
// 参数:
//   - name: 要移除的互斥组名称
//
// 返回值:
//   - error: 移除失败时返回错误
//
// 功能说明:
//   - 根据名称查找并移除互斥组
//   - 使用写锁保护并发安全
//   - 如果找不到对应名称的互斥组, 返回错误
//
// 错误码:
//   - MUTEX_GROUP_NOT_FOUND: 互斥组不存在
func (c *Cmd) RemoveMutexGroup(name string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, group := range c.config.MutexGroups {
		if group.Name == name {
			c.config.MutexGroups = append(c.config.MutexGroups[:i], c.config.MutexGroups[i+1:]...)
			return nil
		}
	}

	return types.NewError("MUTEX_GROUP_NOT_FOUND",
		fmt.Sprintf("mutex group '%s' not found", name), nil)
}

// GetMutexGroup 获取指定名称的互斥组
//
// 参数:
//   - name: 要获取的互斥组名称
//
// 返回值:
//   - *types.MutexGroup: 互斥组指针, 如果找到则返回对应的互斥组
//   - bool: 是否找到, true表示找到
//
// 功能说明:
//   - 根据名称查找互斥组
//   - 使用读锁保护并发安全
//   - 如果找不到对应名称的互斥组, 返回nil和false
func (c *Cmd) GetMutexGroup(name string) (*types.MutexGroup, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for i := range c.config.MutexGroups {
		if c.config.MutexGroups[i].Name == name {
			return &c.config.MutexGroups[i], true
		}
	}
	return nil, false
}

// MutexGroups 获取所有互斥组
//
// 返回值:
//   - []types.MutexGroup: 互斥组列表的副本
//
// 功能说明:
//   - 返回命令中定义的所有互斥组
//   - 返回副本以防止外部修改内部状态
//   - 使用读锁保护并发安全
func (c *Cmd) MutexGroups() []types.MutexGroup {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.config.MutexGroups) == 0 {
		return []types.MutexGroup{}
	}
	groups := make([]types.MutexGroup, len(c.config.MutexGroups))
	copy(groups, c.config.MutexGroups)
	return groups
}

// AddRequiredGroup 添加必需组
//
// 参数:
//   - name: 必需组名称
//   - flags: 必需组中的标志名称列表
//   - conditional: 是否为条件性必需组，如果为true，则只有当组中任何一个标志被设置时，才要求所有标志都被设置
//
// 返回值:
//   - error: 添加失败时返回错误
//
// 功能说明:
//   - 添加一个必需组到命令配置
//   - 如果组名已存在，返回错误
//   - 如果标志列表为空，返回错误
//   - 如果标志不存在，返回错误
//   - 支持条件性必需组，当conditional为true时，只有当组中任何一个标志被设置时，才要求所有标志都被设置
//
// 错误码:
//   - REQUIRED_GROUP_ALREADY_EXISTS: 必需组已存在
//   - EMPTY_REQUIRED_GROUP: 必需组标志列表为空
//   - FLAG_NOT_FOUND: 标志不存在
func (c *Cmd) AddRequiredGroup(name string, flags []string, conditional bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 必需组名称不能为空
	if name == "" {
		return types.NewError("EMPTY_REQUIRED_GROUP_NAME",
			"required group name cannot be empty", nil)
	}

	// 检查必需组名称是否已存在
	for _, group := range c.config.RequiredGroups {
		if group.Name == name {
			return types.NewError("REQUIRED_GROUP_ALREADY_EXISTS",
				fmt.Sprintf("required group '%s' already exists", name), nil)
		}
	}

	// 必需组标志列表不能为空
	if len(flags) == 0 {
		return types.NewError("EMPTY_REQUIRED_GROUP",
			"required group cannot be empty", nil)
	}

	// 检查必需组标志是否存在
	for _, flagName := range flags {
		if _, exists := c.flagRegistry.Get(flagName); !exists {
			return types.NewError("FLAG_NOT_FOUND",
				fmt.Sprintf("flag '%s' not found", flagName), nil)
		}
	}

	// 添加必需组
	group := types.RequiredGroup{
		Name:        name,
		Flags:       flags,
		Conditional: conditional,
	}

	c.config.RequiredGroups = append(c.config.RequiredGroups, group)
	return nil
}

// RemoveRequiredGroup 移除必需组
//
// 参数:
//   - name: 必需组名称
//
// 返回值:
//   - error: 移除失败时返回错误
//
// 功能说明:
//   - 从命令配置中移除指定的必需组
//   - 如果组不存在，返回错误
//
// 错误码:
//   - REQUIRED_GROUP_NOT_FOUND: 必需组不存在
func (c *Cmd) RemoveRequiredGroup(name string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, group := range c.config.RequiredGroups {
		if group.Name == name {
			c.config.RequiredGroups = append(c.config.RequiredGroups[:i], c.config.RequiredGroups[i+1:]...)
			return nil
		}
	}

	return types.NewError("REQUIRED_GROUP_NOT_FOUND",
		fmt.Sprintf("required group '%s' not found", name), nil)
}

// GetRequiredGroup 获取必需组
//
// 参数:
//   - name: 必需组名称
//
// 返回值:
//   - *types.RequiredGroup: 必需组指针
//   - bool: 是否找到
//
// 功能说明:
//   - 根据名称获取必需组
//   - 如果组不存在，返回 nil 和 false
func (c *Cmd) GetRequiredGroup(name string) (*types.RequiredGroup, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for i := range c.config.RequiredGroups {
		if c.config.RequiredGroups[i].Name == name {
			return &c.config.RequiredGroups[i], true
		}
	}

	return nil, false
}

// RequiredGroups 获取所有必需组
//
// 返回值:
//   - []types.RequiredGroup: 所有必需组列表
//
// 功能说明:
//   - 返回命令配置中的所有必需组
//   - 返回的是副本，修改不会影响原配置
func (c *Cmd) RequiredGroups() []types.RequiredGroup {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.config.RequiredGroups) == 0 {
		return []types.RequiredGroup{}
	}
	result := make([]types.RequiredGroup, len(c.config.RequiredGroups))
	copy(result, c.config.RequiredGroups)
	return result
}
