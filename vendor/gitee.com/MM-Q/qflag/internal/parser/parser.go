package parser

import (
	"flag"
	"fmt"
	"strings"

	"gitee.com/MM-Q/qflag/internal/builtin"
	"gitee.com/MM-Q/qflag/internal/types"
)

// DefaultParser 默认解析器实现
//
// DefaultParser 是types.Parser接口的默认实现, 基于Go标准库的flag包。
// 它负责解析命令行参数、处理环境变量和路由子命令。
//
// 特性:
//   - 支持所有标准标志类型
//   - 支持环境变量绑定
//   - 支持子命令解析和路由
//   - 支持内置标志自动处理
type DefaultParser struct {
	flagSet          *flag.FlagSet               // 标准库flag.FlagSet实例
	errorHandling    types.ErrorHandling         // 错误处理策略
	builtinMgr       *builtin.BuiltinFlagManager // 内置标志管理器
	setFlagsMap      map[string]bool             // 已设置标志映射（缓存）
	flagDisplayNames map[string]string           // 所有标志的显示名称映射（缓存）
}

// NewDefaultParser 创建默认解析器实例
//
// 参数:
//   - errorHandling: 错误处理策略, 决定解析错误时的行为
//
// 返回值:
//   - types.Parser: 解析器接口实例
func NewDefaultParser(errorHandling types.ErrorHandling) types.Parser {
	return &DefaultParser{
		errorHandling: errorHandling,
		builtinMgr:    builtin.NewBuiltinFlagManager(), // 初始化内置标志管理器
	}
}

// ParseOnly 仅解析命令行参数, 不执行子命令路由
//
// 参数:
//   - cmd: 要解析的命令
//   - args: 命令行参数列表
//
// 返回值:
//   - error: 如果解析失败返回错误
//
// 注意事项:
//   - 重置所有标志到默认状态（避免重复解析时的遗留值）
//   - 注册内置标志
//   - 创建新的FlagSet实例进行解析
//   - 注册命令的所有标志到FlagSet
//   - 先解析命令行参数
//   - 再加载环境变量 (仅在标志未被命令行参数设置时)
//   - 处理内置标志
//   - 不处理子命令路由
//   - 使用defer确保命令状态和参数在函数返回时被设置
func (p *DefaultParser) ParseOnly(cmd types.Command, args []string) error {
	// 如果禁用标志解析，直接设置参数并返回
	if cmd.IsDisableFlagParsing() {
		cmd.SetParsed(true)
		cmd.SetArgs(args)
		return nil
	}

	// 创建新的 FlagSet 实例
	p.flagSet = flag.NewFlagSet("", p.errorHandling)

	// 自定义 Usage 函数, 避免打印默认的使用说明
	p.flagSet.Usage = func() {
		cmd.PrintHelp()
	}

	// 重置所有标志到默认状态
	// 这对于重复解析场景至关重要：
	// 1. 清除上次解析的遗留值，恢复到默认值
	// 2. 重置 isSet 状态，确保环境变量能正确加载
	// 3. 确保互斥组和必需组验证基于正确的状态
	flagRegistry := cmd.FlagRegistry()
	for _, f := range flagRegistry.List() {
		f.Reset()
	}

	// 注册内置标志
	if err := p.builtinMgr.RegisterBuiltinFlags(cmd); err != nil {
		return err
	}

	// 使用defer确保命令状态和参数在函数返回时被设置
	defer func() {
		cmd.SetParsed(true)
		cmd.SetArgs(p.flagSet.Args())
	}()

	// 注册命令行标志
	for _, f := range flagRegistry.List() {
		p.registerFlag(f)
	}

	// 预检查：扫描未知标志
	if err := checkUnknownFlags(cmd, args); err != nil {
		return err
	}

	// 先解析命令行参数
	if err := p.flagSet.Parse(args); err != nil {
		return err
	}

	// 获取命令配置, 检查是否为nil
	config := cmd.Config()
	if config == nil {
		return fmt.Errorf("nil config in '%s'", cmd.Name())
	}

	// 加载环境变量 (仅在标志未被命令行参数设置时)
	if err := p.loadEnvVars(cmd, config.EnvPrefix); err != nil {
		return err
	}

	// 如果有互斥组或必需组，需要验证, 则构建已设置标志映射
	if len(config.MutexGroups) > 0 || len(config.RequiredGroups) > 0 {
		// 构建已设置标志映射（在验证前构建，确保标志状态已确定）
		p.buildSetFlagsMap(cmd)

		// 验证互斥组规则
		if err := p.validateMutexGroups(config); err != nil {
			return err
		}

		// 验证必需组规则
		if err := p.validateRequiredGroups(config); err != nil {
			return err
		}
	}

	// 处理内置标志
	if err := p.builtinMgr.HandleBuiltinFlags(cmd); err != nil {
		return err
	}

	return nil
}

// Parse 解析命令行参数并处理子命令
//
// 参数:
//   - cmd: 要解析的命令
//   - args: 命令行参数列表
//
// 返回值:
//   - error: 如果解析失败返回错误
//
// 注意事项:
//   - 首先调用ParseOnly解析参数
//   - 检查剩余参数是否为子命令
//   - 如果是子命令, 递归解析子命令
//   - 不执行子命令的运行函数
func (p *DefaultParser) Parse(cmd types.Command, args []string) error {
	// 先解析参数 (ParseOnly 会处理禁用标志解析的情况)
	if err := p.ParseOnly(cmd, args); err != nil {
		return err
	}

	// 检查剩余参数是否为子命令
	cmdRegistry := cmd.CmdRegistry()
	remainingArgs := cmd.Args()

	// 如果有剩余参数, 检查是否为子命令
	if len(remainingArgs) > 0 {
		// 获取第一个参数
		firstArg := remainingArgs[0]

		// 检查是否为子命令, 如果是, 递归解析并执行子命令
		if subCmd, ok := cmdRegistry.Get(firstArg); ok {
			return subCmd.Parse(remainingArgs[1:])
		}

		// 有子命令但没匹配上 → 纠错（但参数以 - 开头的不是子命令）
		if len(cmd.SubCmds()) > 0 && !strings.HasPrefix(firstArg, "-") {
			// 尝试纠错，如果有建议则返回错误，否则继续处理
			if err := newUnknownSubcommandError(cmd, firstArg); err != nil {
				return err
			}
			// 没有找到建议，不拦截，作为普通参数继续处理
		}

		// 没有子命令或参数以 - 开头 → 是普通参数，正常处理
	}

	return nil
}

// ParseAndRoute 解析命令行参数、处理子命令并执行
//
// 参数:
//   - cmd: 要解析的命令
//   - args: 命令行参数列表
//
// 返回值:
//   - error: 如果解析或执行失败返回错误
//
// 注意事项:
//   - 首先调用ParseOnly解析参数
//   - 检查剩余参数是否为子命令
//   - 如果是子命令, 递归解析并执行子命令
//   - 如果不是子命令, 执行当前命令的运行函数
//   - 如果命令没有设置运行函数, 返回错误
func (p *DefaultParser) ParseAndRoute(cmd types.Command, args []string) error {
	// 先解析参数 (ParseOnly 会处理禁用标志解析的情况)
	if err := p.ParseOnly(cmd, args); err != nil {
		return err
	}

	// 检查剩余参数是否为子命令
	cmdRegistry := cmd.CmdRegistry()
	remainingArgs := cmd.Args()

	// 如果是子命令, 递归解析并执行子命令
	if len(remainingArgs) > 0 {
		firstArg := remainingArgs[0] // 获取第一个参数

		// 检查是否为子命令, 如果是, 递归解析并执行子命令
		if subCmd, ok := cmdRegistry.Get(firstArg); ok {
			return subCmd.ParseAndRoute(remainingArgs[1:])
		}

		// 有子命令但没匹配上 → 纠错（但参数以 - 开头的不是子命令）
		if len(cmd.SubCmds()) > 0 && !strings.HasPrefix(firstArg, "-") {
			// 尝试纠错，如果有建议则返回错误，否则继续处理
			if err := newUnknownSubcommandError(cmd, firstArg); err != nil {
				return err
			}
			// 没有找到建议，不拦截，作为普通参数继续处理
		}

		// 没有子命令或参数以 - 开头 → 是普通参数，正常处理
	}

	// 如果不是子命令, 执行当前命令的运行函数
	if cmd.HasRunFunc() {
		return cmd.Run()
	}

	return fmt.Errorf("cmd %q has no run function set", cmd.Name())
}
