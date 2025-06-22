package qflag

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// Cmd 命令行标志管理结构体,封装参数解析、长短标志互斥及帮助系统。
type Cmd struct {
	/* 内部使用属性*/
	fs                  *flag.FlagSet // 底层flag集合, 处理参数解析
	flagRegistry        *FlagRegistry // 标志注册表, 统一管理标志的元数据
	helpFlag            *BoolFlag     // 帮助标志指针,用于绑定和检查
	initFlagBound       bool          // 标记内置标志是否已绑定
	initFlagOnce        sync.Once     // 用于确保内置标志只被绑定一次
	showInstallPathFlag *BoolFlag     // 安装路径标志指针,用于绑定和检查
	subCmds             []*Cmd        // 子命令列表, 用于关联子命令
	parentCmd           *Cmd          // 父命令指针,用于递归调用, 根命令的父命令为nil
	help                string        // 自定义帮助内容
	usage               string        // 用户自定义命令用法信息
	description         string        // 自定义描述,用于帮助信息中显示
	longName            string        // 命令长名称,用于帮助信息中显示
	shortName           string        // 命令短名称,用于帮助信息中显示
	args                []string      // 命令行参数切片
	addMu               sync.Mutex    // 互斥锁,确保并发安全操作
	parseOnce           sync.Once     // 用于确保命令只被解析一次
	setMu               sync.Mutex    // 互斥锁,确保并发安全操作
	builtinFlagNameMap  sync.Map      // 用于存储内置标志名称的映射
	useChinese          bool          // 控制是否使用中文帮助信息
	notes               []string      // 存储备注内容
	examples            []ExampleInfo // 存储示例信息
}

// CmdInterface 命令接口定义，封装命令行程序的核心功能
// 提供统一的命令管理、参数解析和帮助系统
// 实现类需保证线程安全，所有方法应支持并发调用
//
// 示例用法:
// cmd := NewCmd("app", "a", flag.ContinueOnError)
// cmd.SetDescription("示例应用程序")
// cmd.String("config", "c", "配置文件路径", "/etc/app.conf")
type CmdInterface interface {
	LongName() string                  // 获取命令名称(长名称)，如"app"
	ShortName() string                 // 获取命令短名称，如"a"
	Description() string               // 获取命令描述信息
	SetDescription(desc string)        // 设置命令描述信息，用于帮助输出
	Help() string                      // 获取自定义帮助信息
	SetHelp(help string)               // 设置自定义帮助信息，覆盖自动生成内容
	SetUsage(usage string)             // 设置自定义用法说明，覆盖自动生成内容
	GetUseChinese() bool               // 获取是否使用中文帮助信息
	SetUseChinese(useChinese bool)     // 设置是否使用中文帮助信息
	AddSubCmd(subCmd *Cmd)             // 添加子命令，子命令会继承父命令的上下文
	SubCmds() []*Cmd                   // 获取所有已注册的子命令列表
	AddMutexGroup(flags ...Flag) error // 添加互斥组，互斥组内的标志不能同时使用
	Parse(args []string) error         // 解析命令行参数，自动处理标志和子命令
	Args() []string                    // 获取所有非标志参数(未绑定到任何标志的参数)
	Arg(i int) string                  // 获取指定索引的非标志参数，索引越界返回空字符串
	NArg() int                         // 获取非标志参数的数量
	NFlag() int                        // 获取已解析的标志数量
	FlagExists(name string) bool       // 检查指定名称的标志是否存在(支持长/短名称)
	PrintHelp()                        // 打印命令帮助信息
	AddNote(note string)               // 添加备注信息
	GetNotes() []string                // 获取所有备注信息
	AddExample(e ExampleInfo)          // 添加示例信息
	GetExamples() []ExampleInfo        // 获取所有示例信息

	String(longName, shortName, usage, defValue string) *StringFlag                                   // 添加字符串类型标志
	Int(longName, shortName, usage string, defValue int) *IntFlag                                     // 添加整数类型标志
	Bool(longName, shortName, usage string, defValue bool) *BoolFlag                                  // 添加布尔类型标志
	Float(longName, shortName, usage string, defValue float64) *FloatFlag                             // 添加浮点数类型标志
	Duration(longName, shortName, usage string, defValue time.Duration) *DurationFlag                 // 添加时间间隔类型标志
	Enum(longName, shortName string, defValue string, usage string, options []string) *EnumFlag       // 添加枚举类型标志
	StringVar(f *StringFlag, longName, shortName, defValue, usage string)                             // 绑定字符串标志到指定变量
	IntVar(f *IntFlag, longName, shortName string, defValue int, usage string)                        // 绑定整数标志到指定变量
	BoolVar(f *BoolFlag, longName, shortName string, defValue bool, usage string)                     // 绑定布尔标志到指定变量
	FloatVar(f *FloatFlag, longName, shortName string, defValue float64, usage string)                // 绑定浮点数标志到指定变量
	DurationVar(f *DurationFlag, longName, shortName string, defValue time.Duration, usage string)    // 绑定时间间隔类型标志到指定变量
	EnumVar(f *EnumFlag, longName, shortName string, defValue string, usage string, options []string) // 绑定枚举标志到指定变量
}

// GetUseChinese 获取是否使用中文帮助信息
func (c *Cmd) GetUseChinese() bool {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.useChinese
}

// SetUseChinese 设置是否使用中文帮助信息
func (c *Cmd) SetUseChinese(useChinese bool) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.useChinese = useChinese
}

// GetNotes 获取所有备注信息
func (c *Cmd) GetNotes() []string {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	// 返回切片副本而非原始引用
	notes := make([]string, len(c.notes))
	copy(notes, c.notes)
	return notes
}

// LongName 返回命令长名称
func (c *Cmd) LongName() string { return c.longName }

// ShortName 返回命令短名称
func (c *Cmd) ShortName() string { return c.shortName }

// Description 返回命令描述
func (c *Cmd) Description() string { return c.description }

// SetDescription 设置命令描述
func (c *Cmd) SetDescription(desc string) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.description = desc
}

// Help 返回命令用法
func (c *Cmd) Help() string {
	if c.help != "" {
		return c.help
	}

	// 自动生成帮助信息
	return generateHelpInfo(c)
}

// SetUsage 设置自定义命令用法
func (c *Cmd) SetUsage(usage string) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.usage = usage
}

// SetHelp 设置用户自定义命令帮助信息
func (c *Cmd) SetHelp(help string) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.help = help
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
	c.notes = append(c.notes, note)
}

// AddExample 为命令添加使用示例
// description: 示例描述
// usage: 示例使用方式
func (c *Cmd) AddExample(e ExampleInfo) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	// 添加到使用示例列表中
	c.examples = append(c.examples, e)
}

// GetExamples 获取所有使用示例
// 返回示例切片的副本，防止外部修改
func (c *Cmd) GetExamples() []ExampleInfo {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	examples := make([]ExampleInfo, len(c.examples))
	copy(examples, c.examples)
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
			c.flagRegistry = &FlagRegistry{
				mu:       sync.RWMutex{},
				byLong:   make(map[string]*FlagMeta),
				byShort:  make(map[string]*FlagMeta),
				allFlags: []*FlagMeta{},
			}
		}

		if c.helpFlag == nil {
			// 为空时自动初始化
			c.helpFlag = &BoolFlag{}
		}

		// 绑定帮助标志
		helpUsage := "Show help information"
		if c.useChinese {
			helpUsage = "显示帮助信息"
		}
		c.BoolVar(c.helpFlag, helpFlagName, helpFlagShortName, false, helpUsage)

		// 绑定显示安装路径标志
		if c.showInstallPathFlag == nil {
			c.showInstallPathFlag = &BoolFlag{}
		}

		// 绑定显示安装路径标志
		installPathUsage := "Show install path"
		if c.useChinese {
			installPathUsage = "显示安装路径"
		}
		c.BoolVar(c.showInstallPathFlag, showInstallPathFlagName, showInstallPathFlagShortName, false, installPathUsage)

		// 添加内置标志到检测映射
		c.builtinFlagNameMap.Store(helpFlagName, true)
		c.builtinFlagNameMap.Store(helpFlagShortName, true)
		c.builtinFlagNameMap.Store(showInstallPathFlagName, true)
		c.builtinFlagNameMap.Store(showInstallPathFlagShortName, true)

		// 添加默认的注意事项
		if c.useChinese {
			c.AddNote(ChineseTemplate.DefaultNote)
		} else {
			c.AddNote(EnglishTemplate.DefaultNote)
		}

		// 设置内置标志已绑定
		c.initFlagBound = true
	})
}

// // printHelp 打印帮助内容, 优先显示用户自定义的Usage, 否则自动生成
// func (c *Cmd) printHelp() {
// 	// 确保内置标志已初始化
// 	c.initBuiltinFlags()

// 	// 如果用户自定义了Usage，则直接打印
// 	if c.help != "" {
// 		fmt.Println(c.help)
// 	} else {
// 		// 自动生成帮助信息
// 		fmt.Println(generateHelpInfo(c))
// 	}
// }

// validateFlag 通用标志验证逻辑
// 参数:
// longName: 长名称
// shortName: 短名称
// 返回值:
// error: 如果验证失败则返回错误信息,否则返回nil
func (c *Cmd) validateFlag(longName, shortName string) error {
	// 新增格式校验
	if strings.ContainsAny(longName, invalidFlagChars) {
		return fmt.Errorf("The flag name '%s' contains illegal characters", longName)
	}

	// 检查标志名称和短名称是否为空
	if longName == "" {
		return fmt.Errorf("Flag name cannot be empty")
	}
	// if shortName == "" {
	// 	return fmt.Errorf("Flag short name cannot be empty")
	// }

	// 检查标志是否已存在
	if _, exists := c.flagRegistry.GetByName(longName); exists {
		return fmt.Errorf("Flag long name %s already exists", longName)
	}

	if _, exists := c.flagRegistry.GetByName(shortName); exists {
		return fmt.Errorf("Flag short name %s already exists", shortName)
	}

	// 检查标志是否为内置标志
	if _, ok := c.builtinFlagNameMap.Load(longName); ok {
		return fmt.Errorf("Flag long name %s is reserved", longName)
	}
	if _, ok := c.builtinFlagNameMap.Load(shortName); ok {
		return fmt.Errorf("Flag short name %s is reserved", shortName)
	}

	return nil
}

// NewCmd 创建新的命令实例
// 参数:
// longName: 命令长名称
// shortName: 命令短名称
// errorHandling: 错误处理方式
// 返回值: *Cmd命令实例指针
// errorHandling可选值: flag.ContinueOnError、flag.ExitOnError、flag.PanicOnError
func NewCmd(longName string, shortName string, errorHandling flag.ErrorHandling) *Cmd {
	// 检查命令名称是否为空
	if longName == "" {
		panic("cmd long name cannot be empty")
	}

	// 设置默认的错误处理方式为ContinueOnError, 避免测试时意外退出
	if errorHandling == 0 {
		errorHandling = flag.ContinueOnError
	}

	// 创建标志注册表
	flagRegistry := &FlagRegistry{
		mu:       sync.RWMutex{},             // 并发读写锁
		byLong:   make(map[string]*FlagMeta), // 存储长标志的映射
		byShort:  make(map[string]*FlagMeta), // 存储短标志的映射
		allFlags: []*FlagMeta{},              // 存储所有标志的切片
	}

	// 创建新的Cmd实例
	cmd := &Cmd{
		fs:                  flag.NewFlagSet(longName, errorHandling), // 创建新的flag集
		longName:            longName,                                 // 命令长名称, 用于帮助信息中显示
		shortName:           shortName,                                // 命令短名称, 用于帮助信息中显示
		args:                []string{},                               // 命令行参数
		flagRegistry:        flagRegistry,                             // 初始化标志注册表
		helpFlag:            &BoolFlag{},                              // 初始化帮助标志
		showInstallPathFlag: &BoolFlag{},                              // 初始化显示安装路径标志
	}

	return cmd
}

// AddSubCmd 关联一个或多个子命令到当前命令
// 支持批量添加多个子命令，遇到错误时收集所有错误并返回
// 参数:
//
//	subCmds: 一个或多个子命令实例指针
//
// 返回值:
//
//	错误信息列表，如果所有子命令添加成功则返回nil
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
		subCmdNames.Store(strings.ToLower(cmd.longName), true)
		subCmdNames.Store(strings.ToLower(cmd.shortName), true)
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
		if hasCycle(c, cmd) {
			errors = append(errors, fmt.Errorf("Cyclic reference detected: Command %s already exists in the command chain", cmd.longName))
			continue
		}

		// 检测子命令名称是否已存在（大小写不敏感）
		if _, loaded := subCmdNames.LoadOrStore(strings.ToLower(cmd.longName), true); loaded {
			errors = append(errors, fmt.Errorf("Subcommand %s already exists", cmd.longName))
			continue
		}
		if _, loaded := subCmdNames.LoadOrStore(strings.ToLower(cmd.shortName), true); loaded {
			errors = append(errors, fmt.Errorf("Subcommand %s already exists", cmd.shortName))
			continue
		}

		// 如果没有错误，则将子命令添加到切片中
		addedCmds = append(addedCmds, cmd)
	}

	// 如果有验证错误，返回所有错误信息
	if len(errors) > 0 {
		return fmt.Errorf("Failed to add subcommands: %w", joinErrors(errors))
	}

	// 第二阶段：批量添加子命令
	for _, cmd := range addedCmds {
		cmd.parentCmd = c                  // 设置父命令指针
		c.subCmds = append(c.subCmds, cmd) // 添加到子命令列表
	}

	return nil
}

// Parse 解析命令行参数, 自动检查长短标志, 并处理帮助标志
// 参数:
//
//	args: 命令行参数切片
//
// 注意: 该方法保证每个Cmd实例只会解析一次
func (c *Cmd) Parse(args []string) (err error) {
	defer func() {
		// 添加panic捕获
		if r := recover(); r != nil {
			// 使用预定义的恐慌错误常量
			err = fmt.Errorf("%s: %v", ErrPanicRecovered, r)
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
			err = fmt.Errorf("%s: %w", ErrFlagParseFailed, parseErr)
			return
		}

		// 检查是否使用-h/--help标志
		if c.helpFlag.Get() {
			if c.fs.ErrorHandling() != flag.ContinueOnError {
				c.PrintHelp() // 只有在ExitOnError或PanicOnError时才打印使用说明
				os.Exit(0)
			}
			return
		}

		// 检查是否使用-sip/--show-install-path标志
		if c.showInstallPathFlag.Get() {
			if c.fs.ErrorHandling() != flag.ContinueOnError {
				// 只有在ExitOnError或PanicOnError时才打印安装路径
				fmt.Println(GetExecutablePath())
				os.Exit(0)
			}
			return
		}

		// 设置非标志参数
		c.args = append(c.args, c.fs.Args()...)

		// 检查是否有子命令
		if len(c.args) > 0 {
			for _, subCmd := range c.subCmds {
				if c.args[0] == subCmd.longName || c.args[0] == subCmd.shortName {
					// 将剩余参数传递给子命令解析
					if parseErr := subCmd.Parse(c.args[1:]); parseErr != nil {
						err = fmt.Errorf("%s: %w", ErrSubCommandParseFailed, parseErr)
					}
					return
				}
			}
		}

		// 检查枚举类型标志是否有效
		for _, meta := range c.flagRegistry.GetAllFlags() {
			if meta.GetFlagType() == FlagTypeEnum {
				if enumFlag, ok := meta.flag.(*EnumFlag); ok {
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

// String 添加字符串类型标志, 返回标志对象指针
// 参数依次为: 长标志名、短标志、默认值、帮助说明
func (c *Cmd) String(longName, shortName, defValue, usage string) *StringFlag {
	f := &StringFlag{}
	c.StringVar(f, longName, shortName, defValue, usage)
	return f
}

// StringVar 绑定字符串类型标志到指针并内部注册Flag对象
// 参数依次为: 字符串标志指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) StringVar(f *StringFlag, longName, shortName, defValue, usage string) {
	// 检查指针是否为nil
	if f == nil {
		panic("StringFlag pointer cannot be nil")
	}

	// 参数校验（复用公共函数）
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	// 显式初始化当前值的默认值
	currentStr := defValue

	// 修改传入的标志对象
	f.cmd = c               // 修改标志对象 - 命令对象
	f.longName = longName   // 修改标志对象 - 长标志名
	f.shortName = shortName // 修改标志对象 - 短标志名
	f.defValue = defValue   // 修改标志对象 - 默认值
	f.usage = usage         // 修改标志对象 - 帮助说明
	f.value = &currentStr   // 修改标志对象 - 当前值

	// 创建FlagMeta对象
	meta := &FlagMeta{
		flag: f, // 添加标志对象 - Flag对象
	}

	// 绑定短标志
	if shortName != "" {
		c.fs.StringVar(&currentStr, shortName, defValue, usage)
	}

	// 绑定长标志
	c.fs.StringVar(&currentStr, longName, defValue, usage)

	// 注册Flag对象
	if registerErr := c.flagRegistry.RegisterFlag(meta); registerErr != nil {
		panic(registerErr)
	}
}

// IntVar 绑定整数类型标志到指针并内部注册Flag对象
// 参数依次为: 整数标志指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) IntVar(f *IntFlag, longName, shortName string, defValue int, usage string) {
	// 检查指针是否为nil
	if f == nil {
		panic("IntFlag pointer cannot be nil")
	}

	// 参数校验（复用公共函数）
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	// 初始化默认值
	currentInt := defValue

	// 修改传入的标志对象
	f.cmd = c               // 修改标志对象 - 命令对象
	f.longName = longName   // 修改标志对象 - 长标志名
	f.shortName = shortName // 修改标志对象 - 短标志名
	f.defValue = defValue   // 修改标志对象 - 默认值
	f.usage = usage         // 修改标志对象 - 帮助说明
	f.value = &currentInt   // 修改标志对象 - 当前值

	// 创建FlagMeta对象
	meta := &FlagMeta{
		flag: f, // 添加标志对象 - Flag对象
	}

	// 绑定短标志
	if shortName != "" {
		c.fs.IntVar(&currentInt, shortName, defValue, usage)
	}

	// 绑定长标志
	c.fs.IntVar(&currentInt, longName, defValue, usage)

	// 注册Flag对象
	if registerErr := c.flagRegistry.RegisterFlag(meta); registerErr != nil {
		panic(registerErr)
	}
}

// Int 添加整数类型标志, 返回标志对象指针
// 参数依次为: 长标志名、短标志、默认值、帮助说明
func (c *Cmd) Int(longName, shortName string, defValue int, usage string) *IntFlag {
	f := &IntFlag{}
	c.IntVar(f, longName, shortName, defValue, usage)
	return f
}

// BoolVar 绑定布尔类型标志到指针并内部注册Flag对象
// 参数依次为: 布尔标志指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) BoolVar(f *BoolFlag, longName, shortName string, defValue bool, usage string) {
	// 检查指针是否为nil
	if f == nil {
		panic("BoolFlag pointer cannot be nil")
	}

	// 参数校验（复用公共函数）
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	// 修改传入的标志对象
	f.cmd = c               // 修改标志对象 - 命令对象
	f.longName = longName   // 修改标志对象 - 长标志名
	f.shortName = shortName // 修改标志对象 - 短标志名
	f.defValue = defValue   // 修改标志对象 - 默认值
	f.usage = usage         // 修改标志对象 - 帮助说明
	f.value = new(bool)     // 创建当前值指针
	*f.value = defValue

	// 创建FlagMeta对象
	meta := &FlagMeta{
		flag: f, // 添加标志对象 - Flag对象
	}

	// 绑定短标志
	if shortName != "" {
		c.fs.BoolVar(f.value, shortName, defValue, usage)
	}

	// 绑定长标志
	c.fs.BoolVar(f.value, longName, defValue, usage)

	// 注册Flag对象
	if registerErr := c.flagRegistry.RegisterFlag(meta); registerErr != nil {
		panic(registerErr)
	}
}

// Bool 添加布尔类型标志, 返回标志对象指针
// 参数依次为: 长标志名、短标志、默认值、帮助说明
func (c *Cmd) Bool(longName, shortName string, defValue bool, usage string) *BoolFlag {
	f := &BoolFlag{}
	c.BoolVar(f, longName, shortName, defValue, usage)
	return f
}

// Float 添加浮点型标志, 返回标志对象指针
// 参数依次为: 长标志名、短标志、默认值、帮助说明
func (c *Cmd) Float(longName, shortName string, defValue float64, usage string) *FloatFlag {
	f := &FloatFlag{}
	c.FloatVar(f, longName, shortName, defValue, usage)
	return f
}

// FloatVar 绑定浮点型标志到指针并内部注册Flag对象
// 参数依次为: 浮点数标志指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) FloatVar(f *FloatFlag, longName, shortName string, defValue float64, usage string) {
	// 检查指针是否为空
	if f == nil {
		panic("FloatFlag pointer cannot be nil")
	}

	// 参数校验（复用公共函数）
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	// 显式初始化默认值
	currentFloat := new(float64) // 显式堆分配
	*currentFloat = defValue

	// 修改传入的标志对象
	f.cmd = c               // 修改标志对象 - 命令对象
	f.longName = longName   // 修改标志对象 - 长标志名
	f.shortName = shortName // 修改标志对象 - 短标志名
	f.defValue = defValue   // 修改标志对象 - 默认值
	f.usage = usage         // 修改标志对象 - 帮助说明
	f.value = currentFloat  // 修改标志对象 - 当前值

	// 创建FlagMeta对象
	meta := &FlagMeta{
		flag: f, // 添加标志对象 - Flag对象
	}

	// 绑定短标志
	if shortName != "" {
		c.fs.Float64Var(currentFloat, shortName, defValue, usage)
	}

	// 绑定长标志
	c.fs.Float64Var(currentFloat, longName, defValue, usage)

	// 注册Flag对象
	if registerErr := c.flagRegistry.RegisterFlag(meta); registerErr != nil {
		panic(registerErr)
	}
}

// DurationVar 绑定时间间隔类型标志到指针并内部注册Flag对象
// 参数依次为: 时间间隔标志指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) DurationVar(f *DurationFlag, longName, shortName string, defValue time.Duration, usage string) {
	// 检查指针是否为空
	if f == nil {
		panic("DurationFlag pointer cannot be nil")
	}

	// 参数校验
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	// 初始化默认值
	currentDuration := new(time.Duration)
	*currentDuration = defValue

	// 设置标志属性
	f.cmd = c
	f.longName = longName
	f.shortName = shortName
	f.defValue = defValue
	f.usage = usage
	f.value = currentDuration

	// 绑定短标志
	if shortName != "" {
		c.fs.DurationVar(currentDuration, shortName, defValue, usage)
	}

	// 绑定长标志
	c.fs.DurationVar(currentDuration, longName, defValue, usage)

	// 创建并注册标志元数据
	meta := &FlagMeta{
		flag: f, // 添加标志对象 - Flag对象
	}

	// 注册标志元数据
	if registerErr := c.flagRegistry.RegisterFlag(meta); registerErr != nil {
		panic(registerErr)
	}
}

// Duration 添加时间间隔类型标志, 返回标志对象指针
// 参数依次为: 长标志名、短标志、默认值、帮助说明
func (c *Cmd) Duration(longName, shortName string, defValue time.Duration, usage string) *DurationFlag {
	f := &DurationFlag{}
	c.DurationVar(f, longName, shortName, defValue, usage)
	return f
}

// Enum 添加枚举类型标志, 返回标志对象指针
// 参数依次为: 长标志名、短标志、默认值、帮助说明、限制该标志取值的枚举值切片
func (c *Cmd) Enum(longName, shortName string, defValue string, usage string, options []string) *EnumFlag {
	f := &EnumFlag{}
	c.EnumVar(f, longName, shortName, defValue, usage, options)
	return f
}

// EnumVar 绑定枚举类型标志到指针并内部注册Flag对象
// 参数依次为: 枚举标志指针、长标志名、短标志、默认值、帮助说明、限制该标志取值的枚举值切片
func (c *Cmd) EnumVar(f *EnumFlag, longName, shortName string, defValue string, usage string, options []string) {
	// 检查指针是否为空
	if f == nil {
		panic("EnumFlag pointer cannot be nil")
	}

	// 参数校验（复用公共函数）
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	// 初始化枚举值
	if options == nil {
		options = make([]string, 0)
	}

	// 显式初始化当前值的默认值
	currentStr := defValue

	// 创建枚举map
	optionMap := make(map[string]bool)
	if len(options) > 0 {
		for _, v := range options {
			// 转换为小写
			v = strings.ToLower(v)
			optionMap[v] = true
		}
	}

	// 修改传入的标志对象
	f.cmd = c               // 修改标志对象 - 命令对象
	f.longName = longName   // 修改标志对象 - 长标志名
	f.shortName = shortName // 修改标志对象 - 短标志名
	f.defValue = defValue   // 修改标志对象 - 默认值
	f.usage = usage         // 修改标志对象 - 帮助说明
	f.value = &currentStr   // 修改标志对象 - 当前值
	f.optionMap = optionMap // 修改标志对象 - 枚举值map

	// 创建FlagMeta对象
	meta := &FlagMeta{
		flag: f, // 添加标志对象 - Flag对象
	}

	// 绑定短标志
	if shortName != "" {
		c.fs.StringVar(&currentStr, shortName, defValue, usage)
	}

	// 绑定长标志
	c.fs.StringVar(&currentStr, longName, defValue, usage)

	// 注册Flag对象
	if registerErr := c.flagRegistry.RegisterFlag(meta); registerErr != nil {
		panic(registerErr)
	}
}
