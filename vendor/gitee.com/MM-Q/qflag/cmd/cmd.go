package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
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
		QCommandLine = NewCommand("myapp", "", flag.ExitOnError)
	} else {
		// 如果os.Args不为空,则创建一个新的Cmd对象,命令行参数为filepath.Base(os.Args[0]),错误处理方式为ExitOnError
		QCommandLine = NewCommand(filepath.Base(os.Args[0]), "", flag.ExitOnError)
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
	// 标记内置标志是否已绑定
	initFlagBound bool
	// 用于确保内置标志只被绑定一次
	initFlagOnce sync.Once
	// 用于确保命令只被解析一次
	parseOnce sync.Once
	// 子命令列表, 用于关联子命令
	subCmds []*Cmd
	// 父命令指针,用于递归调用, 根命令的父命令为nil
	parentCmd *Cmd
	// 命令行参数切片
	args []string
	// 互斥锁,确保并发安全操作
	addMu sync.Mutex
	// 互斥锁,确保并发安全操作
	setMu sync.Mutex
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
}

// CmdInterface 命令接口定义，封装命令行程序的核心功能
// 提供统一的命令管理、参数解析和帮助系统
// 实现类需保证线程安全，所有方法应支持并发调用
//
// 示例用法:
// cmd := NewCommand("app", "a", flag.ContinueOnError)
// cmd.SetDescription("示例应用程序")
// cmd.String("config", "c", "配置文件路径", "/etc/app.conf")
type CmdInterface interface {
	// 元数据操作方法
	LongName() string                         // 获取命令名称(长名称)，如"app"
	ShortName() string                        // 获取命令短名称，如"a"
	GetDescription() string                   // 获取命令描述信息
	SetDescription(desc string)               // 设置命令描述信息，用于帮助输出
	GetHelp() string                          // 获取自定义帮助信息
	SetHelp(help string)                      // 设置自定义帮助信息，覆盖自动生成内容
	LoadHelp(filePath string) error           // 加载自定义帮助信息，从文件中读取
	SetUsageSyntax(usageSyntax string)        // 设置自定义用法说明，覆盖自动生成内容
	GetUsageSyntax() string                   // 获取自定义用法说明
	GetUseChinese() bool                      // 获取是否使用中文帮助信息
	SetUseChinese(useChinese bool)            // 设置是否使用中文帮助信息
	AddSubCmd(subCmd *Cmd)                    // 添加子命令，子命令会继承父命令的上下文
	SubCmds() []*Cmd                          // 获取所有已注册的子命令列表
	Parse(args []string) error                // 解析命令行参数，自动处理标志和子命令
	ParseFlagsOnly(args []string) (err error) // 仅解析标志参数，不处理子命令
	Args() []string                           // 获取所有非标志参数(未绑定到任何标志的参数)
	Arg(i int) string                         // 获取指定索引的非标志参数，索引越界返回空字符串
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

// 保持兼容API
// 支持 NewCmd 别名
var NewCmd = NewCommand

// NewCommand 创建新的命令实例
// 参数:
// longName: 命令长名称
// shortName: 命令短名称
// errorHandling: 错误处理方式
// 返回值: *Cmd命令实例指针
// errorHandling可选值: flag.ContinueOnError、flag.ExitOnError、flag.PanicOnError
func NewCommand(longName string, shortName string, errorHandling flag.ErrorHandling) *Cmd {
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

	// 确定命令名称：优先使用长名称，如果长名称为空则使用短名称
	cmdName := longName
	if cmdName == "" {
		cmdName = shortName
	}

	// 创建新的Cmd实例
	cmd := &Cmd{
		fs:                  flag.NewFlagSet(cmdName, errorHandling), // 创建新的flag集
		args:                []string{},                              // 命令行参数
		flagRegistry:        flagRegistry,                            // 初始化标志注册表
		helpFlag:            &flags.BoolFlag{},                       // 初始化帮助标志
		showInstallPathFlag: &flags.BoolFlag{},                       // 初始化显示安装路径标志
		versionFlag:         &flags.BoolFlag{},                       // 初始化版本信息标志
		userInfo: UserInfo{
			longName:  longName,  // 命令长名称
			shortName: shortName, // 命令短名称
		},
	}

	return cmd
}

// SetVersion 设置版本信息
func (c *Cmd) SetVersion(version string) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.userInfo.version = version
}

// GetVersion 获取版本信息
func (c *Cmd) GetVersion() string {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.userInfo.version
}

// SetModuleHelps 设置自定义模块帮助信息
func (c *Cmd) SetModuleHelps(moduleHelps string) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.userInfo.moduleHelps = moduleHelps
}

// GetModuleHelps 获取自定义模块帮助信息
func (c *Cmd) GetModuleHelps() string {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.userInfo.moduleHelps
}

// SetLogoText 设置logo文本
func (c *Cmd) SetLogoText(logoText string) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.userInfo.logoText = logoText
}

// GetLogoText 获取logo文本
func (c *Cmd) GetLogoText() string {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.userInfo.logoText
}

// GetUseChinese 获取是否使用中文帮助信息
func (c *Cmd) GetUseChinese() bool {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.userInfo.useChinese
}

// SetUseChinese 设置是否使用中文帮助信息
func (c *Cmd) SetUseChinese(useChinese bool) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.userInfo.useChinese = useChinese
}

// GetNotes 获取所有备注信息
func (c *Cmd) GetNotes() []string {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	// 返回切片副本而非原始引用
	notes := make([]string, len(c.userInfo.notes))
	copy(notes, c.userInfo.notes)
	return notes
}

// LongName 返回命令长名称
func (c *Cmd) LongName() string { return c.userInfo.longName }

// ShortName 返回命令短名称
func (c *Cmd) ShortName() string { return c.userInfo.shortName }

// GetDescription 返回命令描述
func (c *Cmd) GetDescription() string { return c.userInfo.description }

// SetDescription 设置命令描述
func (c *Cmd) SetDescription(desc string) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.userInfo.description = desc
}

// GetHelp 返回命令用法帮助信息
func (c *Cmd) GetHelp() string {
	if c.userInfo.help != "" {
		return c.userInfo.help
	}

	// 自动生成帮助信息
	return generateHelpInfo(c)
}

// SetUsageSyntax 设置自定义命令用法
func (c *Cmd) SetUsageSyntax(usageSyntax string) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.userInfo.usageSyntax = usageSyntax
}

// GetUsageSyntax 获取自定义命令用法
func (c *Cmd) GetUsageSyntax() string {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.userInfo.usageSyntax
}

// SetHelp 设置用户自定义命令帮助信息
func (c *Cmd) SetHelp(help string) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.userInfo.help = help
}

// LoadHelp 从指定文件加载帮助信息
//
// 参数:
// filePath: 帮助信息文件路径
//
// 返回值:
// error: 如果文件不存在或读取文件失败，则返回错误信息
func (c *Cmd) LoadHelp(filePath string) error {
	// 检查是否为空
	if filePath == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	// 直接读取文件内容并处理可能的错误（包括文件不存在的情况）
	// 移除了单独的os.Stat检查以避免竞态条件
	content, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("File %s does not exist", filePath)
		}
		return fmt.Errorf("Failed to read file %s: %w", filePath, err)
	}

	// 设置帮助信息
	c.SetHelp(string(content))
	return nil
}

// SubCmds 返回子命令列表
func (c *Cmd) SubCmds() []*Cmd {
	c.addMu.Lock()
	defer c.addMu.Unlock()
	// 返回子命令切片副本
	subCmds := make([]*Cmd, len(c.subCmds))
	copy(subCmds, c.subCmds)
	return subCmds
}

// Args 获取非标志参数切片
func (c *Cmd) Args() []string {
	// 返回参数切片副本
	args := make([]string, len(c.args))
	copy(args, c.args)
	return args
}

// Arg 获取指定索引的非标志参数
func (c *Cmd) Arg(i int) string {
	if i >= 0 && i < len(c.args) {
		return c.args[i]
	}
	return ""
}

// NArg 获取非标志参数的数量
func (c *Cmd) NArg() int { return len(c.args) }

// NFlag 获取标志的数量
func (c *Cmd) NFlag() int { return c.fs.NFlag() }

// PrintHelp 打印命令的帮助信息, 优先打印用户的帮助信息, 否则自动生成帮助信息
func (c *Cmd) PrintHelp() {
	fmt.Println(generateHelpInfo(c))
}

// FlagExists 检查指定名称的标志是否存在
func (c *Cmd) FlagExists(name string) bool {
	if _, exists := c.flagRegistry.GetByName(name); exists {
		return true
	}

	return false
}

// AddNote 添加备注信息到命令
func (c *Cmd) AddNote(note string) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.userInfo.notes = append(c.userInfo.notes, note)
}

// AddExample 为命令添加使用示例
// description: 示例描述
// usage: 示例使用方式
func (c *Cmd) AddExample(e ExampleInfo) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	// 添加到使用示例列表中
	c.userInfo.examples = append(c.userInfo.examples, e)
}

// GetExamples 获取所有使用示例
// 返回示例切片的副本，防止外部修改
func (c *Cmd) GetExamples() []ExampleInfo {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	examples := make([]ExampleInfo, len(c.userInfo.examples))
	copy(examples, c.userInfo.examples)
	return examples
}

// initBuiltinFlags 初始化内置标志
func (c *Cmd) initBuiltinFlags() {
	// 检查是否已绑定
	if c.initFlagBound {
		return // 避免重复绑定
	}

	// 初始化内置标志
	c.initFlagOnce.Do(func() {
		// 新增: 确保flagRegistry已初始化
		if c.flagRegistry == nil {
			// 为空时自动初始化
			c.flagRegistry = flags.NewFlagRegistry()
		}

		if c.helpFlag == nil {
			// 为空时自动初始化
			c.helpFlag = &flags.BoolFlag{}
		}

		// 绑定帮助标志
		helpUsage := "Show help information"
		if c.GetUseChinese() {
			helpUsage = "显示帮助信息"
		}
		c.BoolVar(c.helpFlag, flags.HelpFlagName, flags.HelpFlagShortName, false, helpUsage)

		// 添加内置标志到检测映射
		c.builtinFlagNameMap.Store(flags.HelpFlagName, true)
		c.builtinFlagNameMap.Store(flags.HelpFlagShortName, true)

		// 只有在根命令上绑定显示程序安装路径标志
		if c.parentCmd == nil {
			// 初始化显示安装路径标志
			if c.showInstallPathFlag == nil {
				c.showInstallPathFlag = &flags.BoolFlag{}
			}

			// 定义显示安装路径标志提示
			installPathUsage := "Show the installation path of the program"
			if c.GetUseChinese() {
				installPathUsage = "显示程序的安装路径"
			}

			// 绑定显示安装路径标志
			c.BoolVar(c.showInstallPathFlag, "", flags.ShowInstallPathFlagName, false, installPathUsage)

			// 添加内置标志到检测映射
			c.builtinFlagNameMap.Store(flags.ShowInstallPathFlagName, true)

			// 绑定版本信息标志
			if c.versionFlag == nil {
				c.versionFlag = &flags.BoolFlag{}
			}

			// 只有在设置了版本信息时才绑定版本信息标志
			if c.GetVersion() != "" {
				// 定义版本标志提示
				versionUsage := "Show the version of the program"
				if c.GetUseChinese() {
					versionUsage = "显示程序的版本信息"
				}

				// 绑定版本信息标志
				c.BoolVar(c.versionFlag, flags.VersionFlagLongName, flags.VersionFlagShortName, false, versionUsage)

				// 添加内置标志到检测映射
				c.builtinFlagNameMap.Store(flags.VersionFlagLongName, true)
				c.builtinFlagNameMap.Store(flags.VersionFlagShortName, true)
			}
		}

		// 添加默认的注意事项
		if c.GetUseChinese() {
			c.AddNote(ChineseTemplate.DefaultNote)
		} else {
			c.AddNote(EnglishTemplate.DefaultNote)
		}

		// 设置内置标志已绑定
		c.initFlagBound = true
	})
}

// validateFlag 通用标志验证逻辑
// 参数:
// longName: 长名称
// shortName: 短名称
// 返回值:
// error: 如果验证失败则返回错误信息,否则返回nil
func (c *Cmd) validateFlag(longName, shortName string) error {
	// 检查标志名称和短名称是否同时为空
	if longName == "" && shortName == "" {
		return fmt.Errorf("Flag long name and short name cannot both be empty")
	}

	// 检查长标志相关逻辑
	if longName != "" {
		// 检查长名称是否包含非法字符
		if strings.ContainsAny(longName, flags.InvalidFlagChars) {
			return fmt.Errorf("The flag long name '%s' contains illegal characters", longName)
		}

		// 检查长标志是否已存在
		if _, exists := c.flagRegistry.GetByName(longName); exists {
			return fmt.Errorf("Flag long name %s already exists", longName)
		}

		// 检查长标志是否为内置标志
		if _, ok := c.builtinFlagNameMap.Load(longName); ok {
			return fmt.Errorf("Flag long name %s is reserved", longName)
		}
	}

	// 检查短标志相关逻辑
	if shortName != "" {
		// 检查短名称是否包含非法字符
		if strings.ContainsAny(shortName, flags.InvalidFlagChars) {
			return fmt.Errorf("The flag short name '%s' contains illegal characters", shortName)
		}

		// 检查短标志是否已存在
		if _, exists := c.flagRegistry.GetByName(shortName); exists {
			return fmt.Errorf("Flag short name %s already exists", shortName)
		}

		// 检查短标志是否为内置标志
		if _, ok := c.builtinFlagNameMap.Load(shortName); ok {
			return fmt.Errorf("Flag short name %s is reserved", shortName)
		}
	}

	return nil
}

// AddSubCmd 关联一个或多个子命令到当前命令
// 支持批量添加多个子命令，遇到错误时收集所有错误并返回
// 参数:
//
//	subCmds: 一个或多个子命令实例指针
//
// 返回值:
//
//	错误信息列表, 如果所有子命令添加成功则返回nil
func (c *Cmd) AddSubCmd(subCmds ...*Cmd) error {
	c.addMu.Lock()
	defer c.addMu.Unlock()

	// 检查子命令是否为空
	if len(subCmds) == 0 {
		return fmt.Errorf("subcommand list cannot be empty")
	}

	// 创建错误切片
	var errors []error

	// 使用sync.Map来存储子命令名称, 解决并发安全问题
	var subCmdNames sync.Map
	for _, cmd := range c.subCmds {
		// 存储子命令的长名称（如果存在）
		if cmd.LongName() != "" {
			subCmdNames.Store(strings.ToLower(cmd.LongName()), true)
		}

		// 存储子命令的短名称（如果存在）
		if cmd.ShortName() != "" {
			subCmdNames.Store(strings.ToLower(cmd.ShortName()), true)
		}
	}

	// 创建一个空的切片，用于存储已添加的子命令
	addedCmds := make([]*Cmd, 0, len(subCmds))

	// 第一阶段：验证所有子命令
	for _, cmd := range subCmds {
		if cmd == nil {
			errors = append(errors, fmt.Errorf("Subcommand cannot be nil"))
			continue
		}

		// 检测循环引用
		if c.hasCycle(cmd) {
			errors = append(errors, fmt.Errorf("Cyclic reference detected: Command %s already exists in the command chain", cmd.LongName()))
			continue
		}

		// 如果设置了长名称，则检查长名称是否已存在（大小写不敏感）
		if cmd.LongName() != "" {
			if _, loaded := subCmdNames.LoadOrStore(strings.ToLower(cmd.LongName()), true); loaded {
				errors = append(errors, fmt.Errorf("Subcommand %s already exists", cmd.LongName()))
				continue
			}
		}

		// 如果设置了短名称，则检查短名称是否已存在（大小写不敏感）
		if cmd.ShortName() != "" {
			if _, loaded := subCmdNames.LoadOrStore(strings.ToLower(cmd.ShortName()), true); loaded {
				errors = append(errors, fmt.Errorf("Subcommand %s already exists", cmd.ShortName()))
				continue
			}
		}

		// 如果没有错误，则将子命令添加到切片中
		addedCmds = append(addedCmds, cmd)
	}

	// 如果有验证错误，返回所有错误信息
	if len(errors) > 0 {
		return fmt.Errorf("Failed to add subcommands: %w", qerr.JoinErrors(errors))
	}

	// 第二阶段：批量添加子命令
	for _, cmd := range addedCmds {
		cmd.parentCmd = c                  // 设置父命令指针
		c.subCmds = append(c.subCmds, cmd) // 添加到子命令列表
	}

	return nil
}

// Parse 完整解析命令行参数（含子命令处理）
// 主要功能：
//  1. 解析当前命令的长短标志及内置标志
//  2. 自动检测并解析子命令及其参数（若存在）
//  3. 验证枚举类型标志的有效性
//
// 参数：
//
//	args: 原始命令行参数切片（包含可能的子命令及参数）
//
// 返回值：
//
//	解析过程中遇到的错误（如标志格式错误、子命令解析失败等）
//
// 注意事项：
//   - 每个Cmd实例仅会被解析一次（线程安全）
//   - 若检测到子命令，会将剩余参数传递给子命令的Parse方法
//   - 处理内置标志执行逻辑
func (c *Cmd) Parse(args []string) (err error) {
	defer func() {
		// 添加panic捕获
		if r := recover(); r != nil {
			// 使用预定义的恐慌错误常量
			err = fmt.Errorf("%s: %v", qerr.ErrPanicRecovered, r)
		}
	}()

	// 如果命令为空，则返回错误
	if c == nil {
		return fmt.Errorf("cmd cannot be nil")
	}

	// 确保只解析一次
	c.parseOnce.Do(func() {
		// 初始化内置标志
		c.initBuiltinFlags()

		// 设置使用说明
		c.fs.Usage = func() {
			c.PrintHelp()
		}

		// 调用flag库解析参数
		if parseErr := c.fs.Parse(args); parseErr != nil {
			err = fmt.Errorf("%s: %w", qerr.ErrFlagParseFailed, parseErr)
			return
		}

		// 检查是否使用-h/--help标志
		if c.helpFlag.Get() {
			c.PrintHelp()
			if c.fs.ErrorHandling() != flag.ContinueOnError {
				// 只有在ExitOnError或PanicOnError时才退出
				os.Exit(0)
			}
			return
		}

		// 只有在顶级命令中处理-sip/--show-install-path和-v/--version标志
		if c.parentCmd == nil {
			// 检查是否使用-sip/--show-install-path标志
			if c.showInstallPathFlag.Get() {
				fmt.Println(GetExecutablePath())
				if c.fs.ErrorHandling() != flag.ContinueOnError {
					// 只有在ExitOnError或PanicOnError时才退出
					os.Exit(0)
				}
				return
			}

			// 检查是否使用-v/--version标志
			if c.versionFlag.Get() {
				fmt.Println(c.GetVersion())
				if c.fs.ErrorHandling() != flag.ContinueOnError {
					// 只有在ExitOnError或PanicOnError时才退出
					os.Exit(0)
				}
				return
			}
		}

		// 设置非标志参数
		c.args = append(c.args, c.fs.Args()...)

		// 检查是否有子命令
		if len(c.args) > 0 {
			for _, subCmd := range c.subCmds {
				// 第一个非标志参数如果匹配到子命令，则解析子命令
				if c.args[0] == subCmd.LongName() || c.args[0] == subCmd.ShortName() {
					// 将剩余参数传递给子命令解析
					if parseErr := subCmd.Parse(c.args[1:]); parseErr != nil {
						err = fmt.Errorf("%s: %w", qerr.ErrSubCommandParseFailed, parseErr)
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
						// 如果验证失败，则返回错误信息，错误信息： 无效的枚举值, 可选值: [a, b, c]
						err = checkErr
					}
				}
			}
		}
	})

	// 检查是否报错
	if err != nil {
		return err
	}

	return nil
}

// ParseFlagsOnly 仅解析当前命令的标志参数（忽略子命令）
// 主要功能：
//  1. 解析当前命令的长短标志及内置标志
//  2. 验证枚举类型标志的有效性
//  3. 明确忽略所有子命令及后续参数
//
// 参数：
//
//	args: 原始命令行参数切片（子命令及后续参数会被忽略）
//
// 返回值：
//
//	解析过程中遇到的错误（如标志格式错误等）
//
// 注意事项：
//   - 每个Cmd实例仅会被解析一次（线程安全）
//   - 不会处理任何子命令，所有参数均视为当前命令的标志或位置参数
//   - 处理内置标志逻辑
func (c *Cmd) ParseFlagsOnly(args []string) (err error) {
	defer func() {
		// 添加panic捕获
		if r := recover(); r != nil {
			// 使用预定义的恐慌错误常量
			err = fmt.Errorf("%s: %v", qerr.ErrPanicRecovered, r)
		}
	}()

	// 确保只解析一次
	c.parseOnce.Do(func() {
		// 初始化内置标志
		c.initBuiltinFlags()

		// 设置使用说明
		c.fs.Usage = func() {
			c.PrintHelp()
		}

		// 调用flag库解析参数
		if parseErr := c.fs.Parse(args); parseErr != nil {
			err = fmt.Errorf("%s: %w", qerr.ErrFlagParseFailed, parseErr)
			return
		}

		// 检查是否使用-h/--help标志
		if c.helpFlag.Get() {
			c.PrintHelp()
			if c.fs.ErrorHandling() != flag.ContinueOnError {
				// 只有在ExitOnError或PanicOnError时才退出
				os.Exit(0)
			}
			return
		}

		// 只有在顶级命令中处理-sip/--show-install-path和-v/--version标志
		if c.parentCmd == nil {
			// 检查是否使用-sip/--show-install-path标志
			if c.showInstallPathFlag.Get() {
				fmt.Println(GetExecutablePath())
				if c.fs.ErrorHandling() != flag.ContinueOnError {
					// 只有在ExitOnError或PanicOnError时才退出
					os.Exit(0)
				}
				return
			}

			// 检查是否使用-v/--version标志
			if c.versionFlag.Get() {
				fmt.Println(c.GetVersion())
				if c.fs.ErrorHandling() != flag.ContinueOnError {
					// 只有在ExitOnError或PanicOnError时才退出
					os.Exit(0)
				}
				return
			}
		}

		// 设置非标志参数
		c.args = append(c.args, c.fs.Args()...)

		// 检查枚举类型标志是否有效
		for _, meta := range c.flagRegistry.GetAllFlags() {
			if meta.GetFlagType() == flags.FlagTypeEnum {
				if enumFlag, ok := meta.GetFlag().(*flags.EnumFlag); ok {
					// 调用IsCheck方法进行验证
					if checkErr := enumFlag.IsCheck(enumFlag.Get()); checkErr != nil {
						// 如果验证失败，则返回错误信息，错误信息： 无效的枚举值, 可选值: [a, b, c]
						err = checkErr
					}
				}
			}
		}
	})

	// 检查是否报错
	if err != nil {
		return err
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
	return c.dfs(child, visited)
}

// dfs 深度优先搜索检测循环引用
// 递归检查当前节点及其子命令、父命令链中是否包含目标节点(q)
// 参数:
//
//	current: 当前遍历的命令节点
//	visited: 已访问节点集合，防止重复遍历
//
// 返回值:
//
//	找到目标节点返回true，否则返回false
func (c *Cmd) dfs(current *Cmd, visited map[*Cmd]bool) bool {
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
		if c.dfs(subCmd, visited) {
			return true
		}
	}

	// 检查父命令链
	if current.parentCmd != nil {
		return c.dfs(current.parentCmd, visited)
	}

	return false
}

// GetExecutablePath 获取程序的绝对安装路径
// 如果无法通过 os.Executable 获取路径,则使用 os.Args[0] 作为替代
// 返回：程序的绝对路径字符串
func GetExecutablePath() string {
	// 尝试使用 os.Executable 获取可执行文件的绝对路径
	exePath, err := os.Executable()
	if err != nil {
		// 如果 os.Executable 报错,使用 os.Args[0] 作为替代
		exePath = os.Args[0]
	}
	// 使用 filepath.Abs 确保路径是绝对路径
	absPath, err := filepath.Abs(exePath)
	if err != nil {
		// 如果 filepath.Abs 报错,直接返回原始路径
		return exePath
	}
	return absPath
}
