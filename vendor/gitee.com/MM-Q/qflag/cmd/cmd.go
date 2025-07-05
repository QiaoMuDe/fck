package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"gitee.com/MM-Q/qflag/flags"
	"gitee.com/MM-Q/qflag/qerr"
)

// QCommandLine 全局默认Command实例
var QCommandLine *Cmd

// 在包初始化时创建全局默认Cmd实例
func init() {
	// 处理可能的空os.Args情况
	if len(os.Args) == 0 {
		// 如果os.Args为空,则创建一个新的Cmd对象,命令行参数为"myapp",短名字为"",错误处理方式为ExitOnError
		QCommandLine = NewCmd("myapp", "", flag.ExitOnError)
	} else {
		// 如果os.Args不为空,则创建一个新的Cmd对象,命令行参数为filepath.Base(os.Args[0]),错误处理方式为ExitOnError
		QCommandLine = NewCmd(filepath.Base(os.Args[0]), "", flag.ExitOnError)
	}
}

// UserInfo 存储用户自定义信息的嵌套结构体
type UserInfo struct {
	// 命令长名称
	longName string

	// 命令短名称
	shortName string

	// 版本信息
	version string

	// 自定义描述
	description string

	// 自定义的完整命令行帮助信息
	help string

	// 自定义用法格式说明
	usageSyntax string

	// 模块帮助信息
	moduleHelps string

	// logo文本
	logoText string

	// 备注内容切片
	notes []string

	// 示例信息切片
	examples []ExampleInfo

	// 是否使用中文帮助信息
	useChinese bool
}

// Cmd 命令行标志管理结构体,封装参数解析、长短标志互斥及帮助系统。
type Cmd struct {
	// 底层flag集合, 处理参数解析
	fs *flag.FlagSet

	// 标志注册表, 统一管理标志的元数据
	flagRegistry *flags.FlagRegistry

	// 用于确保命令只被解析一次
	parseOnce sync.Once
	parsed    atomic.Bool // 标记是否已完成解析

	// 子命令映射表, 用于关联和查找子命令
	subCmdMap map[string]*Cmd

	// 子命令切片, 用于存储唯一实例(无重复)
	subCmds []*Cmd

	// 父命令指针,用于递归调用, 根命令的父命令为nil
	parentCmd *Cmd

	// 命令行参数切片
	args []string

	// 读写锁, 确保并发安全操作, 读操作使用RLock/RUnlock, 写操作使用Lock/Unlock
	rwMu sync.RWMutex

	// 用于存储内置标志名称的映射
	builtinFlagNameMap sync.Map

	// 用户自定义信息
	userInfo UserInfo

	// 帮助标志指针,用于绑定和检查
	helpFlag *flags.BoolFlag

	// 安装路径标志指针,用于绑定和检查
	showInstallPathFlag *flags.BoolFlag

	// 版本标志指针,用于绑定和检查
	versionFlag *flags.BoolFlag

	// 控制内置标志是否自动退出
	exitOnBuiltinFlags bool

	// 控制是否禁用内置标志注册
	disableBuiltinFlags bool
}

// CmdInterface 命令接口定义, 封装命令行程序的核心功能
// 提供统一的命令管理、参数解析和帮助系统
// 实现类需保证线程安全, 所有方法应支持并发调用
//
// 示例用法:
// cmd := NewCmd("app", "a", flag.ContinueOnError)
// cmd.SetDescription("示例应用程序")
// cmd.String("config", "c", "配置文件路径", "/etc/app.conf")
type CmdInterface interface {
	// 元数据操作方法
	Name() string                             // 获取命令名称
	LongName() string                         // 获取命令名称(长名称), 如"app"
	ShortName() string                        // 获取命令短名称, 如"a"
	GetDescription() string                   // 获取命令描述信息
	SetDescription(desc string)               // 设置命令描述信息, 用于帮助输出
	GetHelp() string                          // 获取自定义帮助信息
	SetHelp(help string)                      // 设置自定义帮助信息, 覆盖自动生成内容
	LoadHelp(filePath string) error           // 加载自定义帮助信息, 从文件中读取
	SetUsageSyntax(usageSyntax string)        // 设置自定义用法说明, 覆盖自动生成内容
	GetUsageSyntax() string                   // 获取自定义用法说明
	GetUseChinese() bool                      // 获取是否使用中文帮助信息
	SetUseChinese(useChinese bool)            // 设置是否使用中文帮助信息
	AddSubCmd(subCmd *Cmd)                    // 添加子命令, 子命令会继承父命令的上下文
	SubCmds() []*Cmd                          // 获取所有已注册的子命令列表
	SubCmdMap() map[string]*Cmd               // 获取所有已注册的子命令映射表
	Args() []string                           // 获取所有非标志参数(未绑定到任何标志的参数)
	Arg(i int) string                         // 获取指定索引的非标志参数, 索引越界返回空字符串
	NArg() int                                // 获取非标志参数的数量
	NFlag() int                               // 获取已解析的标志数量
	FlagExists(name string) bool              // 检查指定名称的标志是否存在(支持长/短名称)
	PrintHelp()                               // 打印命令帮助信息
	AddNote(note string)                      // 添加备注信息
	GetNotes() []string                       // 获取所有备注信息
	AddExample(e ExampleInfo)                 // 添加示例信息
	GetExamples() []ExampleInfo               // 获取所有示例信息
	SetVersion(version string)                // 设置版本信息
	GetVersion() string                       // 获取版本信息
	SetLogoText(logoText string)              // 设置logo文本
	GetLogoText() string                      // 获取logo文本
	SetModuleHelps(moduleHelps string)        // 设置自定义模块帮助信息
	GetModuleHelps() string                   // 获取自定义模块帮助信息
	SetExitOnBuiltinFlags(exit bool) *Cmd     // 设置是否在添加内置标志时退出
	SetDisableBuiltinFlags(disable bool) *Cmd // 设置是否禁用内置标志注册
	CmdExists(cmdName string) bool            // 判断命令行参数中是否存在指定标志

	// 标志解析方法
	Parse(args []string) error                // 解析命令行参数, 自动处理标志和子命令
	ParseFlagsOnly(args []string) (err error) // 仅解析标志参数, 不处理子命令
	IsParsed() bool                           // 检查是否已解析命令行参数

	// 添加标志方法
	String(longName, shortName, usage, defValue string) *flags.StringFlag                             // 添加字符串类型标志
	Int(longName, shortName, usage string, defValue int) *flags.IntFlag                               // 添加整数类型标志
	Int64(longName, shortName, usage string, defValue int64) *flags.Int64Flag                         // 添加64位整数类型标志
	Bool(longName, shortName, usage string, defValue bool) *flags.BoolFlag                            // 添加布尔类型标志
	Float64(longName, shortName, usage string, defValue float64) *flags.Float64Flag                   // 添加浮点数类型标志
	Duration(longName, shortName, usage string, defValue time.Duration) *flags.DurationFlag           // 添加时间间隔类型标志
	Enum(longName, shortName string, defValue string, usage string, options []string) *flags.EnumFlag // 添加枚举类型标志
	Slice(longName, shortName string, defValue []string, usage string) *flags.SliceFlag               // 添加字符串切片类型标志
	uint16(longName, shortName string, defValue uint16, usage string) *flags.Uint16Flag               // 添加无符号16位整型标志
	Time(longName, shortName string, defValue time.Time, usage string) *flags.TimeFlag                // 添加时间类型标志
	Map(longName, shortName string, defValue map[string]string, usage string) *flags.MapFlag          // 添加Map标志
	Path(longName, shortName string, defValue string, usage string) *flags.PathFlag                   // 添加路径标志

	// 绑定标志方法
	StringVar(f *flags.StringFlag, longName, shortName, defValue, usage string)                             // 绑定字符串标志到指定变量
	IntVar(f *flags.IntFlag, longName, shortName string, defValue int, usage string)                        // 绑定整数标志到指定变量
	Int64Var(f *flags.Int64Flag, longName, shortName string, defValue int64, usage string)                  // 绑定64位整数标志到指定变量
	BoolVar(f *flags.BoolFlag, longName, shortName string, defValue bool, usage string)                     // 绑定布尔标志到指定变量
	Float64Var(f *flags.Float64Flag, longName, shortName string, defValue float64, usage string)            // 绑定浮点数标志到指定变量
	DurationVar(f *flags.DurationFlag, longName, shortName string, defValue time.Duration, usage string)    // 绑定时间间隔类型标志到指定变量
	EnumVar(f *flags.EnumFlag, longName, shortName string, defValue string, usage string, options []string) // 绑定枚举标志到指定变量
	SliceVar(f *flags.SliceFlag, longName, shortName string, defValue []string, usage string)               // 绑定字符串切片标志到指定变量
	Uint16Var(f *flags.Uint16Flag, longName, shortName string, defValue uint16, usage string)               // 绑定无符号16位整型标志到指定变量
	TimeVar(f *flags.TimeFlag, longName, shortName string, defValue time.Time, usage string)                // 绑定时间类型标志到指定变量
	MapVar(f *flags.MapFlag, longName, shortName string, defValue map[string]string, usage string)          // 绑定字符串映射标志到指定变量
	PathVar(f *flags.PathFlag, longName, shortName string, defValue string, usage string)                   // 绑定路径标志到指定变量
}

// NewCmd 创建新的命令实例
//
// 参数:
//
//   - longName: 命令长名称
//   - shortName: 命令短名称
//   - errorHandling: 错误处理方式
//
// 返回值:
//   - *Cmd: 新的命令实例指针
//
// errorHandling可选值:
//
//   - flag.ContinueOnError: 解析标志时遇到错误继续解析, 并返回错误信息
//   - flag.ExitOnError: 解析标志时遇到错误立即退出程序, 并返回错误信息
//   - flag.PanicOnError: 解析标志时遇到错误立即触发panic
func NewCmd(longName string, shortName string, errorHandling flag.ErrorHandling) *Cmd {
	// 检查命令名称是否同时为空
	if longName == "" && shortName == "" {
		panic("cmd long name and short name cannot both be empty")
	}

	// 设置默认的错误处理方式为ContinueOnError, 避免测试时意外退出
	if errorHandling == 0 {
		errorHandling = flag.ContinueOnError
	}

	// 创建标志注册表
	flagRegistry := flags.NewFlagRegistry()

	// 确定命令名称：优先使用长名称, 如果长名称为空则使用短名称
	cmdName := longName
	if cmdName == "" {
		cmdName = shortName
	}

	// 创建新的Cmd实例
	cmd := &Cmd{
		fs:                  flag.NewFlagSet(cmdName, errorHandling), // 创建新的flag集
		args:                []string{},                              // 命令行参数
		subCmdMap:           map[string]*Cmd{},                       // 子命令映射
		subCmds:             []*Cmd{},                                // 子命令切片
		flagRegistry:        flagRegistry,                            // 初始化标志注册表
		helpFlag:            &flags.BoolFlag{},                       // 初始化帮助标志
		showInstallPathFlag: &flags.BoolFlag{},                       // 初始化显示安装路径标志
		versionFlag:         &flags.BoolFlag{},                       // 初始化版本信息标志
		userInfo: UserInfo{
			longName:  longName,        // 命令长名称
			shortName: shortName,       // 命令短名称
			notes:     []string{},      // 命令备注
			examples:  []ExampleInfo{}, // 命令示例
		},
		exitOnBuiltinFlags: true,       // 默认保持原有行为, 在解析内置标志后退出
		builtinFlagNameMap: sync.Map{}, // 内置标志名称映射
	}

	return cmd
}

// SetExitOnBuiltinFlags 设置是否在解析内置参数时退出
// 默认情况下为true, 当解析到内置参数时, QFlag将退出程序
//
// 参数:
//   - exit: 是否退出
//
// 返回值:
//   - *cmd.Cmd: 当前命令对象
func (c *Cmd) SetExitOnBuiltinFlags(exit bool) *Cmd {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()
	c.exitOnBuiltinFlags = exit

	// 返回当前Cmd实例
	return c
}

// SetDisableBuiltinFlags 设置是否禁用内置标志注册
//
// 参数: disable - true表示禁用内置标志注册, false表示启用(默认)
//
// 返回值: 当前Cmd实例, 支持链式调用
func (c *Cmd) SetDisableBuiltinFlags(disable bool) *Cmd {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()
	c.disableBuiltinFlags = disable
	return c
}

// SubCmdMap 返回子命令映射表
func (c *Cmd) SubCmdMap() map[string]*Cmd {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()

	// 返回map副本避免外部修改
	subCmdMap := make(map[string]*Cmd, len(c.subCmdMap))

	// 遍历子命令映射表, 将每个子命令复制到新的map中
	for name, cmd := range c.subCmdMap {
		subCmdMap[name] = cmd
	}
	return subCmdMap
}

// SubCmds 返回子命令切片
func (c *Cmd) SubCmds() []*Cmd {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()

	// 创建一个切片副本
	result := make([]*Cmd, len(c.subCmds))

	// 拷贝子命令切片
	copy(result, c.subCmds)

	return result
}

// AddSubCmd 关联一个或多个子命令到当前命令
// 支持批量添加多个子命令, 遇到错误时收集所有错误并返回
//
// 参数:
//
//   - subCmds: 一个或多个子命令实例指针
//
// 返回值:
//
//   - 错误信息, 如果所有子命令添加成功则返回nil
func (c *Cmd) AddSubCmd(subCmds ...*Cmd) error {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

	// 检查子命令是否为空
	if len(subCmds) == 0 {
		return fmt.Errorf("subcommand list cannot be empty")
	}

	// 检查子命令map是否为nil
	if c.subCmdMap == nil {
		c.subCmdMap = make(map[string]*Cmd)
	}

	// 验证阶段 - 收集所有错误
	var errors []error
	validCmds := make([]*Cmd, 0, len(subCmds)) // 预分配空间

	// 验证所有子命令
	for _, cmd := range subCmds {
		if err := c.validateSubCmd(cmd); err != nil {
			errors = append(errors, fmt.Errorf("invalid subcommand %s: %w", getCmdIdentifier(cmd), err))
			continue
		}
		validCmds = append(validCmds, cmd)
	}

	// 如果有验证错误, 返回所有错误信息
	if len(errors) > 0 {
		return fmt.Errorf("failed to add subcommands: %w", qerr.JoinErrors(errors))
	}

	// 预分配临时切片(容量=validCmds长度, 避免多次扩容)
	tempList := make([]*Cmd, 0, len(validCmds))

	// 添加阶段 - 仅处理通过验证的命令
	for _, cmd := range validCmds {
		cmd.parentCmd = c                  // 设置子命令的父命令指针
		c.subCmdMap[cmd.ShortName()] = cmd // 将子命令的短名称和实例关联
		c.subCmdMap[cmd.LongName()] = cmd  // 将子命令的长名称和实例关联
		tempList = append(tempList, cmd)   // 先添加到临时切片
	}

	// 一次性合并到目标切片
	c.subCmds = append(c.subCmds, tempList...)

	return nil
}

// Parse 完整解析命令行参数(含子命令处理)
// 主要功能：
//  1. 解析当前命令的长短标志及内置标志
//  2. 自动检测并解析子命令及其参数(若存在)
//  3. 验证枚举类型标志的有效性
//
// 参数：
//
//	args: 原始命令行参数切片(包含可能的子命令及参数)
//
// 返回值：
//
//	解析过程中遇到的错误(如标志格式错误、子命令解析失败等)
//
// 注意事项：
//   - 每个Cmd实例仅会被解析一次(线程安全)
//   - 若检测到子命令, 会将剩余参数传递给子命令的Parse方法
//   - 处理内置标志执行逻辑
func (c *Cmd) Parse(args []string) (err error) {
	err, shouldExit := c.parseCommon(args, true)
	if shouldExit {
		// 延迟处理内置标志的退出
		os.Exit(0)
	}
	return err
}

// ParseFlagsOnly 仅解析当前命令的标志参数(忽略子命令)
// 主要功能：
//  1. 解析当前命令的长短标志及内置标志
//  2. 验证枚举类型标志的有效性
//  3. 明确忽略所有子命令及后续参数
//
// 参数：
//
//	args: 原始命令行参数切片(子命令及后续参数会被忽略)
//
// 返回值：
//
//	解析过程中遇到的错误(如标志格式错误等)
//
// 注意事项：
//   - 每个Cmd实例仅会被解析一次(线程安全)
//   - 不会处理任何子命令, 所有参数均视为当前命令的标志或位置参数
//   - 处理内置标志逻辑
func (c *Cmd) ParseFlagsOnly(args []string) (err error) {
	err, shouldExit := c.parseCommon(args, false)
	if shouldExit {
		os.Exit(0)
	}
	return err
}

// parseCommon 命令行参数解析公共逻辑
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
		// 添加panic捕获
		if r := recover(); r != nil {
			// 使用预定义的恐慌错误变量
			err = fmt.Errorf("%w: %v", qerr.ErrPanicRecovered, r)
		}
	}()

	// 如果命令为空, 则返回错误
	if c == nil {
		return fmt.Errorf("cmd cannot be nil"), false
	}

	// 核心功能组件校验 (必须初始化)
	if c.fs == nil {
		return fmt.Errorf("flag.FlagSet instance is not initialized"), false
	}
	if c.flagRegistry == nil {
		return fmt.Errorf("FlagRegistry instance is not initialized"), false
	}
	if c.subCmdMap == nil {
		return fmt.Errorf("subCmdMap cannot be nil"), false
	}

	// 内置标志校验 (根据启用状态决定是否需要校验)
	if !c.disableBuiltinFlags {
		if c.helpFlag == nil {
			return fmt.Errorf("help flag is not initialized"), false
		}
		if c.versionFlag == nil {
			return fmt.Errorf("version flag is not initialized"), false
		}
		if c.showInstallPathFlag == nil {
			return fmt.Errorf("showInstallPath flag is not initialized"), false
		}
	}

	// 确保只解析一次
	c.parseOnce.Do(func() {
		defer c.parsed.Store(true) // 无论成功失败均标记为已解析
		// 只有在没有禁用内置标志时才注册内置标志
		if !c.disableBuiltinFlags {
			// 定义帮助标志提示信息
			helpUsage := flags.HelpFlagUsageEn
			if c.GetUseChinese() {
				helpUsage = flags.HelpFlagUsageZh
			}

			// 注册帮助标志
			c.BoolVar(c.helpFlag, flags.HelpFlagName, flags.HelpFlagShortName, false, helpUsage)

			// 添加到内置标志名称映射
			c.builtinFlagNameMap.Store(flags.HelpFlagName, true)
			c.builtinFlagNameMap.Store(flags.HelpFlagShortName, true)

			// 只有在根命令上注册显示程序安装路径标志和版本信息标志
			if c.parentCmd == nil {
				// 定义显示安装路径标志提示信息
				installPathUsage := flags.ShowInstallPathFlagUsageEn
				if c.GetUseChinese() {
					installPathUsage = flags.ShowInstallPathFlagUsageZh
				}

				// 绑定显示安装路径标志
				c.BoolVar(c.showInstallPathFlag, "", flags.ShowInstallPathFlagName, false, installPathUsage)

				// 添加到内置标志名称映射
				c.builtinFlagNameMap.Store(flags.ShowInstallPathFlagName, true)

				// 只有在版本信息不为空时才注册版本信息标志
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
		}

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

		// 调用flag库解析参数
		if parseErr := c.fs.Parse(args); parseErr != nil {
			err = fmt.Errorf("%w: %w", qerr.ErrFlagParseFailed, parseErr)
			return
		}

		// 只有在没有禁用内置标志时才处理内置标志
		if !c.disableBuiltinFlags {
			// 检查是否使用-h/--help标志
			if c.helpFlag.Get() {
				c.PrintHelp()
				if c.exitOnBuiltinFlags {
					// 仅在ExitOnBuiltinFlags时才退出
					shouldExit = true // 标记需要退出
					return            // 退出当前函数, 执行defer清理
				}
				return
			}

			// 只有在顶级命令中处理-sip/--show-install-path和-v/--version标志
			if c.parentCmd == nil {
				// 检查是否使用-sip/--show-install-path标志
				if c.showInstallPathFlag.Get() {
					fmt.Println(GetExecutablePath())
					if c.exitOnBuiltinFlags {
						// 仅在ExitOnBuiltinFlags时才退出
						shouldExit = true // 标记需要退出
						return            // 退出当前函数, 执行defer清理
					}
					return
				}

				// 检查是否使用-v/--version标志
				if c.versionFlag.Get() && c.GetVersion() != "" {
					fmt.Println(c.GetVersion())
					if c.exitOnBuiltinFlags {
						// 仅在ExitOnBuiltinFlags时才退出
						shouldExit = true // 标记需要退出
						return            // 退出当前函数, 执行defer清理
					}
					return
				}
			}
		}

		// 检查枚举类型标志是否有效
		for _, meta := range c.flagRegistry.GetAllFlags() {
			if meta.GetFlagType() == flags.FlagTypeEnum {
				if enumFlag, ok := meta.GetFlag().(*flags.EnumFlag); ok {
					// 调用IsCheck方法进行验证
					if checkErr := enumFlag.IsCheck(enumFlag.Get()); checkErr != nil {
						// 添加标志名称到错误信息, 便于定位问题
						err = fmt.Errorf("flag %s: %w", meta.GetName(), checkErr)
						break // 发现第一个错误后立即退出循环
					}
				}
			}
		}

		// 枚举标志检查失败立即返回
		if err != nil {
			return
		}

		// 设置非标志参数
		c.args = append(c.args, c.fs.Args()...)

		// 如果允许解析子命令, 则进入子命令解析阶段, 否则跳过
		if parseSubcommands {
			// 如果存在子命令并且非标志参数不为0
			if len(c.args) > 0 && len(c.subCmdMap) > 0 {
				// 获取参数
				arg := c.args[0]

				// 直接通过参数查找(map键已包含长名称和短名称)
				if subCmd, ok := c.subCmdMap[arg]; ok {
					// 将剩余参数传递给子命令解析
					if parseErr := subCmd.Parse(c.args[1:]); parseErr != nil {
						err = fmt.Errorf("%w for '%s': %v", qerr.ErrSubCommandParseFailed, arg, parseErr)
					}
					return
				}
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
