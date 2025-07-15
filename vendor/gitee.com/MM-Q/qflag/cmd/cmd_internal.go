// cmd_internal 包含 Cmd 的内部实现细节，不对外暴露
// 变更需同步更新 cmd.go 中的公共接口文档
package cmd

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"gitee.com/MM-Q/qflag/flags"
	"gitee.com/MM-Q/qflag/qerr"
)

// parseCommon 命令行参数解析公共逻辑
//
// 主要功能：
//  1. 通用参数解析流程(标志解析、内置标志处理、错误处理)
//  2. 枚举类型标志验证
//  3. 可选的子命令解析支持
//
// 参数：
//
//	args: 原始命令行参数切片
//	parseSubcommands: 是否解析子命令(true: 解析子命令, false: 忽略子命令)
//
// 返回值：
//
//   - 解析过程中遇到的错误(如标志格式错误、子命令解析失败等)
//   - 是否需要退出程序, 用于处理内部选项标志的解析处理情况(true: 需要退出, false: 不需要退出)
//
// 注意事项：
//   - 每个Cmd实例仅会被解析一次(线程安全)
//   - 内置标志(-h/--help, -v/--version等)处理逻辑在此实现
//   - 子命令解析仅在parseSubcommands=true时执行
func (c *Cmd) parseCommon(args []string, parseSubcommands bool) (err error, shouldExit bool) {
	defer func() {
		// 添加panic捕获逻辑
		if r := recover(); r != nil {
			// 使用预定义的恐慌错误变量，并添加详细的堆栈信息
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			err = fmt.Errorf("%w: %v\nStack: %s", qerr.ErrPanicRecovered, r, string(buf[:n]))
		}
	}()

	// 如果命令为空, 则返回错误
	if c == nil {
		return fmt.Errorf("cmd cannot be nil"), false
	}

	// 调用提取的组件校验方法
	if err = c.validateComponents(); err != nil {
		return err, false
	}

	// 确保只解析一次
	c.parseOnce.Do(func() {
		defer c.parsed.Store(true) // 在返回时, 无论成功失败均标记为已解析

		// 调用内置标志注册方法
		c.registerBuiltinFlags()

		// 添加默认的注意事项
		if c.GetUseChinese() {
			c.AddNote(ChineseTemplate.DefaultNote)
		} else {
			c.AddNote(EnglishTemplate.DefaultNote)
		}

		// 设置底层flag库的Usage函数
		c.fs.Usage = func() {
			c.PrintHelp()
		}

		// 解析前加载环境变量的参数值
		if err = c.loadEnvVars(); err != nil {
			err = fmt.Errorf("%w: %v", qerr.ErrEnvLoadFailed, err)
			return
		}

		// 调用底层flag库解析参数
		if parseErr := c.fs.Parse(args); parseErr != nil {
			err = fmt.Errorf("%w: %w", qerr.ErrFlagParseFailed, parseErr)
			return
		}

		// 调用内置标志处理方法
		exit, handleErr := c.handleBuiltinFlags()
		if handleErr != nil {
			// 处理内置标志错误
			err = handleErr
			return
		}

		// 内置标志处理是否需要退出程序
		if exit {
			shouldExit = true
			return
		}

		// 设置当前命令的非标志参数
		c.args = append(c.args, c.fs.Args()...)

		// 如果允许解析子命令, 则进入子命令解析阶段, 否则跳过子命令解析
		if parseSubcommands {
			// 如果存在子命令并且非标志参数不为0
			if len(c.args) > 0 && (len(c.subCmdMap) > 0 && len(c.subCmds) > 0) {
				// 获取参数的第一个值(子命令名称: 长名或短名)
				arg := c.args[0]

				// 保存剩余参数
				remainingArgs := make([]string, len(c.args)-1)
				copy(remainingArgs, c.args[1:])

				// 直接通过参数0查找子命令, 如果存在则解析子命令
				if subCmd, ok := c.subCmdMap[arg]; ok {
					// 将剩余参数传递给子命令解析
					if parseErr := subCmd.Parse(remainingArgs); parseErr != nil {
						err = fmt.Errorf("%w for '%s': %v", qerr.ErrSubCommandParseFailed, arg, parseErr)
					}
					return
				}
			}
		}

		// 调用自定义解析阶段钩子函数
		if c.ParseHook != nil {
			// 执行钩子函数
			hookErr, hookExit := c.ParseHook(c)
			// 处理钩子函数错误
			if hookErr != nil {
				err = hookErr
				return
			}
			// 钩子函数是否需要退出程序
			if hookExit {
				shouldExit = true
				return
			}
		}
	})

	// 检查是否报错
	if err != nil {
		return err, false
	}

	// 根据内置标志处理结果决定是否退出
	return nil, shouldExit
}

// validateComponents 校验核心组件和内置标志的初始化状态
//
// 返回值：
//   - 组件校验错误信息
func (c *Cmd) validateComponents() error {
	// 核心功能组件校验 (必须初始化)
	if c.fs == nil {
		return fmt.Errorf("flag.FlagSet instance is not initialized")
	}
	if c.flagRegistry == nil {
		return fmt.Errorf("FlagRegistry instance is not initialized")
	}
	if c.subCmdMap == nil {
		return fmt.Errorf("subCmdMap cannot be nil")
	}

	// 内置标志校验
	if c.helpFlag == nil {
		return fmt.Errorf("help flag is not initialized")
	}
	if c.versionFlag == nil {
		return fmt.Errorf("version flag is not initialized")
	}

	return nil
}

// registerBuiltinFlags 注册内置标志
// 仅在顶级命令中注册
func (c *Cmd) registerBuiltinFlags() {
	// 仅在顶级命令中注册内置标志
	if c.parentCmd != nil {
		return
	}

	// 如果启用了自动补全功能
	if c.enableCompletion {
		// 语言配置结构体：集中管理所有语言相关资源
		type languageConfig struct {
			notes     []string      // 注意事项
			shellDesc string        // 用法描述
			examples  []ExampleInfo // 示例
		}

		// 一次性判断语言并初始化所有相关资源
		useChinese := c.GetUseChinese()
		var langConfig languageConfig

		// 根据语言选择对应的资源
		if useChinese {
			langConfig = languageConfig{
				notes:     completionNotesCN,
				shellDesc: fmt.Sprintf(flags.CompletionShellDescCN, flags.ShellSlice),
				examples:  completionExamplesCN,
			}
		} else {
			langConfig = languageConfig{
				notes:     completionNotesEN,
				shellDesc: fmt.Sprintf(flags.CompletionShellDescEN, flags.ShellSlice),
				examples:  completionExamplesEN,
			}
		}

		// 注册补全标志
		c.EnumVar(c.completionShell, flags.CompletionShellFlagLongName, flags.CompletionShellFlagShortName, flags.ShellNone, langConfig.shellDesc, flags.ShellSlice)

		// 注册到内置的标志名称映射
		c.builtinFlagNameMap.Store(flags.CompletionShellFlagLongName, true)
		c.builtinFlagNameMap.Store(flags.CompletionShellFlagShortName, true)

		// 添加自动补全子命令的注意事项
		c.userInfo.notes = append(c.userInfo.notes, langConfig.notes...)

		// 获取运行的程序名
		cmdName := os.Args[0]

		// 添加自动补全子命令的示例
		for _, ex := range langConfig.examples {
			// 直接添加到底层切片中
			c.userInfo.examples = append(c.userInfo.examples, ExampleInfo{
				Description: ex.Description,
				Usage:       fmt.Sprintf(ex.Usage, cmdName),
			})
		}
	}

	// 仅在版本信息不为空时注册(-v/--version)
	if c.GetVersion() != "" {
		// 定义版本信息标志提示信息
		versionUsage := flags.VersionFlagUsageEn
		if c.GetUseChinese() {
			versionUsage = flags.VersionFlagUsageZh
		}

		// 注册版本信息标志
		c.BoolVar(c.versionFlag, flags.VersionFlagLongName, flags.VersionFlagShortName, false, versionUsage)

		// 添加到内置标志名称映射
		c.builtinFlagNameMap.Store(flags.VersionFlagLongName, true)
		c.builtinFlagNameMap.Store(flags.VersionFlagShortName, true)
	}
}

// handleBuiltinFlags 处理内置标志(-h/--help, -v/--version, --generate-shell-completion/-gsc等)的逻辑
//
// 返回值:
//   - 是否需要退出程序
//   - 处理过程中遇到的错误
func (c *Cmd) handleBuiltinFlags() (bool, error) {
	// 检查是否使用-h/--help标志
	if c.helpFlag.Get() {
		c.PrintHelp()
		if c.exitOnBuiltinFlags {
			return true, nil // 标记需要退出
		}
		return false, nil
	}

	// 只有在顶级命令中处理-v/--version标志
	if c.parentCmd == nil {
		// 检查是否使用-v/--version标志
		if c.versionFlag.Get() && c.GetVersion() != "" {
			fmt.Println(c.GetVersion())
			if c.exitOnBuiltinFlags {
				return true, nil // 标记需要退出
			}
			return false, nil
		}
	}

	// 检查是否启用补全功能
	if c.enableCompletion {
		// 获取shell类型
		shell := c.completionShell.Get()

		// 只有不是默认值时才生成补全脚本
		if shell != flags.ShellNone {
			// 生成对应shell的补全脚本
			switch shell {
			case flags.ShellBash, flags.ShellPowershell, flags.ShellPwsh: // 受支持的shell类型
				shellCompletion, err := c.generateShellCompletion(shell)
				if err != nil {
					return false, err
				}
				fmt.Println(shellCompletion)
			default:
				return false, fmt.Errorf("unsupported shell: %s. Supported shells are: %v", shell, flags.ShellSlice)
			}

			// 仅在生成补全脚本后检查退出标志
			if c.exitOnBuiltinFlags {
				return true, nil // 标记需要退出
			}
			return false, nil
		}
	}

	// 检查枚举类型标志是否有效
	for _, meta := range c.flagRegistry.GetAllFlagMetas() {
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

// loadEnvVars 从环境变量加载参数值
//
// 优先级：命令行参数 > 环境变量 > 默认值
//
// 参数：无
//
// 返回值：
//
//	error - 加载过程中的错误（如有）
func (c *Cmd) loadEnvVars() error {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()

	// 存储读取错误
	var errors []error

	// 预分配map容量以提高性能,初始容量为已注册标志数量
	// 使用所有标志总数作为容量最大基准, 确保独立长/短标志场景下容量充足
	processedEnvs := make(map[string]bool, c.flagRegistry.GetALLFlagsCount()) // 跟踪已处理的环境变量，避免重复处理

	// 遍历所有已注册的标志
	c.fs.VisitAll(func(f *flag.Flag) {
		// 获取标志实例
		flagInstance, ok := f.Value.(flags.Flag)
		if !ok {
			return
		}

		// 获取环境变量名称
		envVar := flagInstance.GetEnvVar()
		if envVar == "" {
			// 环境变量未设置，提前返回
			return
		}

		// 检查是否已处理过该环境变量（避免长短标志重复处理）
		if processedEnvs[envVar] {
			return
		}

		// 读取环境变量值
		envValue := os.Getenv(envVar)
		if envValue == "" {
			return // 环境变量未设置，提前返回
		}

		// 标记该环境变量为已处理
		processedEnvs[envVar] = true

		// 设置标志值(使用现有Set方法进行类型转换)
		if err := f.Value.Set(envValue); err != nil {
			errors = append(errors, qerr.NewValidationErrorf("Failed to parse environment variable %s for flag %s: %v", envVar, f.Name, err))
		}
	})

	// 函数末尾返回聚合错误
	if len(errors) > 0 {
		return qerr.NewValidationErrorf("Failed to load environment variables: %v", qerr.JoinErrors(errors))
	}

	return nil
}

// hasCycle 检测当前命令与待添加子命令间是否存在循环引用
// 循环引用场景包括：
// 1. 子命令直接或间接引用当前命令
// 2. 子命令的父命令链中包含当前命令
// 参数:
//
//	child: 待添加的子命令实例
//
// 返回值:
//
//	存在循环引用返回true，否则返回false
func (c *Cmd) hasCycle(child *Cmd) bool {
	if c == nil || child == nil {
		return false
	}

	visited := make(map[*Cmd]bool)

	// 添加初始深度参数0
	return c.dfs(child, visited, 0)
}

// dfs 深度优先搜索检测循环引用
// 递归检查当前节点及其子命令、父命令链中是否包含目标节点(q)
// 参数:
//
//		current: 当前遍历的命令节点
//		visited: 已访问节点集合，防止重复遍历
//	  depth: 当前递归深度，用于防止无限递归
//
// 返回值:
//
//	找到目标节点返回true, 否则返回false
func (c *Cmd) dfs(current *Cmd, visited map[*Cmd]bool, depth int) bool {
	// 添加递归深度限制(100层)
	if depth > 100 {
		panic(fmt.Sprintf("Potential circular reference detected (recursion depth exceeds %d), there may be circular dependencies between commands", depth))
		//return true // 视为存在循环以中断递归
	}

	// 已访问过当前节点，直接返回避免无限循环
	if visited[current] {
		return false
	}
	visited[current] = true

	// 找到目标节点，存在循环引用
	if current == c {
		return true
	}

	// 递归检查所有子命令
	for _, subCmd := range current.subCmds {
		if c.dfs(subCmd, visited, depth+1) {
			return true
		}
	}

	// 检查父命令链
	if current.parentCmd != nil {
		return c.dfs(current.parentCmd, visited, depth+1)
	}

	return false
}

// getCmdIdentifier 获取命令的标识字符串，用于错误信息
//
// 参数：
//   - cmd: 命令对象
//
// 返回：
//   - 命令标识字符串, 如果为空则返回 <nil>
func getCmdIdentifier(cmd *Cmd) string {
	if cmd == nil {
		return "<nil>"
	}
	return cmd.Name()
}

// validateFlag 通用标志验证逻辑
//
// 参数:
//   - longName: 长名称
//   - shortName: 短名称
//
// 返回值:
//   - error: 如果验证失败则返回错误信息,否则返回nil
func (c *Cmd) validateFlag(longName, shortName string) error {
	// 检查标志名称和短名称是否同时为空
	if longName == "" && shortName == "" {
		return fmt.Errorf("flag long name and short name cannot both be empty")
	}

	// 检查长标志相关逻辑
	if longName != "" {
		// 检查长名称是否包含非法字符
		if strings.ContainsAny(longName, flags.InvalidFlagChars) {
			return fmt.Errorf("the flag long name '%s' contains illegal characters", longName)
		}

		// 检查长标志是否已存在
		if _, exists := c.flagRegistry.GetByName(longName); exists {
			return fmt.Errorf("flag long name %s already exists", longName)
		}

		// 检查长标志是否为内置标志
		if _, ok := c.builtinFlagNameMap.Load(longName); ok {
			return fmt.Errorf("flag long name %s is reserved", longName)
		}
	}

	// 检查短标志相关逻辑
	if shortName != "" {
		// 检查短名称是否包含非法字符
		if strings.ContainsAny(shortName, flags.InvalidFlagChars) {
			return fmt.Errorf("the flag short name '%s' contains illegal characters", shortName)
		}

		// 检查短标志是否已存在
		if _, exists := c.flagRegistry.GetByName(shortName); exists {
			return fmt.Errorf("flag short name %s already exists", shortName)
		}

		// 检查短标志是否为内置标志
		if _, ok := c.builtinFlagNameMap.Load(shortName); ok {
			return fmt.Errorf("flag short name %s is reserved", shortName)
		}
	}

	return nil
}

// validateSubCmd 验证单个子命令的有效性
//
// 参数：
//   - cmd: 要验证的子命令实例
//
// 返回值：
//   - error: 验证失败时返回的错误信息, 否则返回nil
func (c *Cmd) validateSubCmd(cmd *Cmd) error {
	if cmd == nil {
		return fmt.Errorf("subcmd %s is nil", getCmdIdentifier(cmd))
	}

	// 检测循环引用
	if c.hasCycle(cmd) {
		return fmt.Errorf("cyclic reference detected: Command %s already exists in the command chain", getCmdIdentifier(cmd))
	}

	// 检查子命令的长名称是否已存在
	if _, exists := c.subCmdMap[cmd.LongName()]; exists {
		return fmt.Errorf("long name '%s' already exists", cmd.LongName())
	}

	// 检查子命令的短名称是否已存在
	if _, exists := c.subCmdMap[cmd.ShortName()]; exists {
		return fmt.Errorf("short name '%s' already exists", cmd.ShortName())
	}

	return nil
}
