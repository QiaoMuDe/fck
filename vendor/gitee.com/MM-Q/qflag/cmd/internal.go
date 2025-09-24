// cmd_internal 包含 Cmd 的内部实现细节，不对外暴露
// Package cmd 内部实现和辅助功能
// 本文件包含了Cmd结构体的内部实现方法和辅助功能，提供命令行解析的核心逻辑。
// 变更需同步更新 cmd.go 中的公共接口文档。
package cmd

import (
	"fmt"
	"os"
	"runtime/debug"

	"gitee.com/MM-Q/qflag/flags"
	"gitee.com/MM-Q/qflag/internal/completion"
	"gitee.com/MM-Q/qflag/internal/parser"
	"gitee.com/MM-Q/qflag/internal/types"
	"gitee.com/MM-Q/qflag/qerr"
)

// parseCommon 解析命令行参数的公共逻辑
//
// 参数:
//   - args: 命令行参数切片
//   - parseSubcommands: 是否解析子命令
//
// 返回值:
//   - shouldExit: 是否需要退出程序
//   - err: 解析过程中遇到的错误
func (c *Cmd) parseCommon(args []string, parseSubcommands bool) (shouldExit bool, err error) {
	defer func() {
		// 添加panic捕获逻辑
		if r := recover(); r != nil {
			err = fmt.Errorf("%w: %v\nStack: %s", qerr.ErrPanicRecovered, r, debug.Stack())
		}
	}()

	// 检查命令是否为nil
	if c == nil {
		return false, fmt.Errorf("cmd: nil command")
	}

	// 调用提取的组件校验方法
	if err := c.validateComponents(); err != nil {
		return false, err
	}

	c.ctx.ParseOnce.Do(func() {
		defer c.ctx.Parsed.Store(true) // 在返回时, 无论成功失败均标记为已解析

		// 调用内置标志注册方法
		c.registerBuiltinFlags()

		// 设置底层flag库的Usage函数
		c.ctx.FlagSet.Usage = func() {
			c.PrintHelp()
		}

		// 解析当前命令的参数
		if parseErr := parser.ParseCommand(c.ctx, args); parseErr != nil {
			err = parseErr
			return
		}

		// 处理内置标志
		exit, handleErr := c.handleBuiltinFlags()
		if handleErr != nil {
			err = handleErr
			return
		}

		// 内置标志处理是否需要退出程序
		if exit {
			shouldExit = true
			return
		}

		// 解析子命令参数
		if parseSubcommands {
			exit, parseErr := c.parseSubCommands()
			if parseErr != nil {
				err = parseErr
				return
			}
			if exit {
				shouldExit = true
				return
			}
		}

		// 执行解析钩子
		if c.ctx.ParseHook != nil {
			hookErr, hookExit := c.ctx.ParseHook(c.ctx)
			if hookErr != nil {
				err = hookErr
				return
			}
			if hookExit {
				shouldExit = true
				return
			}
		}
	})

	return shouldExit, err
}

// parseSubCommands 解析子命令
//
// 返回值:
//   - shouldExit: 是否需要退出程序
//   - err: 解析过程中遇到的错误
func (c *Cmd) parseSubCommands() (bool, error) {
	// 如果没有非标志参数或者没有注册子命令，则无需解析子命令
	if len(c.ctx.Args) == 0 || len(c.ctx.SubCmdMap) == 0 {
		return false, nil
	}

	// 获取非标志参数的第一个参数(子命令名称)
	subCmdName := c.ctx.Args[0]

	// 检查子命令是否存在
	subCmd, exists := c.ctx.SubCmdMap[subCmdName]
	if !exists {
		return false, nil
	}

	// 获取除子命令名称外的剩余参数
	argsToProcess := make([]string, len(c.ctx.Args)-1)
	copy(argsToProcess, c.ctx.Args[1:])

	// 解析子命令的参数(这里创建了临时Cmd实例包装器)
	exit, parseErr := tempCmd(subCmd).parseCommon(argsToProcess, true)
	if parseErr != nil {
		return false, fmt.Errorf("%w for '%s': %v", qerr.ErrSubCommandParseFailed, subCmdName, parseErr)
	}

	return exit, nil
}

// validateComponents 校验核心组件和内置标志的初始化状态
//
// 返回值：
//   - 组件校验错误信息
func (c *Cmd) validateComponents() error {
	// 核心功能组件校验 (必须初始化)
	if c.ctx.FlagSet == nil {
		return fmt.Errorf("flag.FlagSet instance is not initialized")
	}
	if c.ctx.FlagRegistry == nil {
		return fmt.Errorf("FlagRegistry instance is not initialized")
	}
	if c.ctx.SubCmds == nil {
		return fmt.Errorf("subCmdMap cannot be nil")
	}

	// 内置标志校验
	if c.ctx.BuiltinFlags.Help == nil {
		return fmt.Errorf("help flag is not initialized")
	}
	if c.ctx.BuiltinFlags.Version == nil {
		return fmt.Errorf("version flag is not initialized")
	}

	return nil
}

// registerBuiltinFlags 注册内置标志
// 仅在顶级命令中注册
func (c *Cmd) registerBuiltinFlags() {
	// 仅在顶级命令中注册内置标志
	if c.ctx.Parent != nil {
		return
	}

	// 如果启用了自动补全功能
	if c.ctx.Config.EnableCompletion {
		// 语言配置结构体：集中管理所有语言相关资源
		type languageConfig struct {
			notes     []string            // 注意事项
			shellDesc string              // 用法描述
			examples  []types.ExampleInfo // 示例
		}

		// 一次性判断语言并初始化所有相关资源
		useChinese := c.ctx.Config.UseChinese
		var langConfig languageConfig

		// 根据语言选择对应的资源
		if useChinese {
			langConfig = languageConfig{
				notes:     completion.CompletionNotesCN,                               // 注意事项
				shellDesc: fmt.Sprintf(flags.CompletionShellDescCN, flags.ShellSlice), // 使用帮助
				examples:  completion.CompletionExamplesCN,                            // 示例
			}
		} else {
			langConfig = languageConfig{
				notes:     completion.CompletionNotesEN,                               // 注意事项
				shellDesc: fmt.Sprintf(flags.CompletionShellDescEN, flags.ShellSlice), // 使用帮助
				examples:  completion.CompletionExamplesEN,                            // 示例
			}
		}

		// 注册补全标志
		c.EnumVar(c.ctx.BuiltinFlags.Completion, flags.CompletionShellFlagLongName, "", flags.ShellNone, langConfig.shellDesc, flags.ShellSlice)

		// 注册到内置的标志名称映射
		c.ctx.BuiltinFlags.NameMap.Store(flags.CompletionShellFlagLongName, true)

		// 添加自动补全子命令的注意事项
		c.ctx.Config.Notes = append(c.ctx.Config.Notes, langConfig.notes...)

		// 获取运行的程序名
		cmdName := os.Args[0]

		// 添加自动补全子命令的示例
		for _, ex := range langConfig.examples {
			// 直接添加到底层切片中
			c.ctx.Config.Examples = append(c.ctx.Config.Examples, types.ExampleInfo{
				Description: ex.Description,
				Usage:       fmt.Sprintf(ex.Usage, cmdName),
			})
		}
	}

	// 仅在版本信息不为空时注册(-v/--version)
	if c.ctx.Config.Version != "" {
		// 注册版本信息标志
		c.BoolVar(c.ctx.BuiltinFlags.Version, flags.VersionFlagLongName, flags.VersionFlagShortName, false, flags.VersionFlagUsage)

		// 添加到内置标志名称映射
		c.ctx.BuiltinFlags.NameMap.Store(flags.VersionFlagLongName, true)
		c.ctx.BuiltinFlags.NameMap.Store(flags.VersionFlagShortName, true)
	}
}

// handleBuiltinFlags 处理内置标志(-h/--help, -v/--version, --completion等)的逻辑
//
// 返回值:
//   - 是否需要退出程序
//   - 处理过程中遇到的错误
func (c *Cmd) handleBuiltinFlags() (bool, error) {
	// 检查是否使用-h/--help标志
	if c.ctx.BuiltinFlags.Help.Get() {
		c.PrintHelp()
		if c.ctx.Config.ExitOnBuiltinFlags {
			return true, nil // 标记需要退出
		}
		return false, nil
	}

	// 只有在顶级命令中处理-v/--version标志
	if c.ctx.Parent == nil {
		// 检查是否使用-v/--version标志
		if c.ctx.BuiltinFlags.Version.Get() && c.ctx.Config.Version != "" {
			fmt.Println(c.ctx.Config.Version)
			if c.ctx.Config.ExitOnBuiltinFlags {
				return true, nil // 标记需要退出
			}
			return false, nil
		}
	}

	// 检查是否启用补全功能
	if c.ctx.Config.EnableCompletion {
		// 获取shell类型
		shell := c.ctx.BuiltinFlags.Completion.Get()

		// 只有不是默认值时才生成补全脚本
		if shell != flags.ShellNone {
			completion, err := completion.GenerateShellCompletion(c.ctx, shell)
			if err != nil {
				return false, err
			}
			fmt.Println(completion)
			return c.ctx.Config.ExitOnBuiltinFlags, nil
		}
	}

	// 检查枚举类型标志是否有效
	for _, meta := range c.ctx.FlagRegistry.GetFlagMetaList() {
		// 不是枚举类型标志，跳过
		if meta.GetFlagType() != flags.FlagTypeEnum {
			continue
		}

		// 获取枚举类型标志失败，跳过
		enumFlag, ok := meta.GetFlag().(*flags.EnumFlag)
		if !ok {
			continue
		}

		// 调用IsCheck方法进行验证
		if checkErr := enumFlag.IsCheck(enumFlag.Get()); checkErr != nil {
			// 添加标志名称到错误信息, 便于定位问题
			return false, fmt.Errorf("flag %s: %w", meta.GetName(), checkErr)
		}
	}

	return false, nil
}

// tempCmd 创建临时的 Cmd 包装器用于递归解析子命令
//
// 这个函数解决了 *types.CmdContext 无法直接调用 parseCommon 方法的问题
// 创建的临时实例仅用于方法调用，不会被其他地方引用
//
// 参数:
//   - ctx: 子命令的上下文
//
// 返回值:
//   - *Cmd: 临时的 Cmd 包装器
func tempCmd(ctx *types.CmdContext) *Cmd {
	return &Cmd{ctx: ctx}
}
