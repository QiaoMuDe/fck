// cmd 命令行标志管理结构体,封装参数解析、长短标志互斥及帮助系统。
package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"gitee.com/MM-Q/qflag/flags"
	"gitee.com/MM-Q/qflag/qerr"
)

// QCommandLine 全局默认Command实例
var QCommandLine *Cmd

// 在包初始化时创建全局默认Cmd实例
func init() {
	// 使用一致的命令名生成逻辑
	cmdName := "myapp"
	if len(os.Args) > 0 {
		cmdName = filepath.Base(os.Args[0])
	}

	// 创建全局默认Cmd实例
	QCommandLine = NewCmd(cmdName, "", flag.ExitOnError)
}

// userInfo 存储用户自定义信息的嵌套结构体
type userInfo struct {
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
	userInfo userInfo

	// 帮助标志指针,用于绑定和检查
	helpFlag *flags.BoolFlag

	// 版本标志指针,用于绑定和检查
	versionFlag *flags.BoolFlag

	// 自动补全标志指针,用于绑定和检查
	completionShell *flags.EnumFlag

	// 控制内置标志是否自动退出
	exitOnBuiltinFlags bool

	// 控制是否启用自动补全功能
	enableCompletion bool

	// 解析阶段钩子函数
	// 在标志解析完成后、子命令参数处理后调用
	//
	// 参数:
	//   - 当前命令实例
	//
	// 返回值:
	//   - error: 错误信息, 非nil时会中断解析流程
	//   - bool: 是否需要退出程序
	ParseHook func(*Cmd) (error, bool)
}

// NewCmd 创建新的命令实例
//
// 参数:
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

	// 创建标志注册表
	flagRegistry := flags.NewFlagRegistry()

	// 确定命令名称：优先使用长名称, 如果长名称为空则使用短名称
	cmdName := longName
	if cmdName == "" {
		cmdName = shortName
	}

	// 创建新的Cmd实例
	cmd := &Cmd{
		fs:              flag.NewFlagSet(cmdName, errorHandling), // 创建新的flag集
		args:            []string{},                              // 命令行参数
		subCmdMap:       map[string]*Cmd{},                       // 子命令映射
		subCmds:         []*Cmd{},                                // 子命令切片
		flagRegistry:    flagRegistry,                            // 初始化标志注册表
		helpFlag:        &flags.BoolFlag{},                       // 初始化帮助标志
		versionFlag:     &flags.BoolFlag{},                       // 初始化版本信息标志
		completionShell: &flags.EnumFlag{},                       // 初始化自动完成标志
		userInfo: userInfo{
			longName:  longName,        // 命令长名称
			shortName: shortName,       // 命令短名称
			notes:     []string{},      // 命令备注
			examples:  []ExampleInfo{}, // 命令示例
		},
		exitOnBuiltinFlags: true,  // 默认保持原有行为, 在解析内置标志后退出
		enableCompletion:   false, // 默认关闭自动补全
	}

	// 注册帮助标志
	cmd.BoolVar(cmd.helpFlag, flags.HelpFlagName, flags.HelpFlagShortName, false, flags.HelpFlagUsageEn)

	// 添加到内置标志名称映射
	cmd.builtinFlagNameMap.Store(flags.HelpFlagName, true)
	cmd.builtinFlagNameMap.Store(flags.HelpFlagShortName, true)

	return cmd
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

// AddSubCmd 添加外部子命令到当前命令
// 支持批量添加多个子命令, 遇到错误时收集所有错误并返回
//
// 参数:
//   - subCmds: 一个或多个子命令实例指针
//
// 返回值:
//   - 错误信息, 如果所有子命令添加成功则返回nil
func (c *Cmd) AddSubCmd(subCmds ...*Cmd) error {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

	// 检查子命令是否为空
	if len(subCmds) == 0 {
		return qerr.NewValidationError("subCmds list cannot be empty")
	}

	// 检查子命令map是否为nil
	if c.subCmdMap == nil {
		return qerr.NewValidationError("subCmdMap cannot be nil")
	}

	// 检查子命令数组是否为nil
	if c.subCmds == nil {
		return qerr.NewValidationError("subCmds cannot be nil")
	}

	// 验证阶段 - 收集所有错误
	var errors []error
	validCmds := make([]*Cmd, 0, len(subCmds)) // 预分配空间

	// 验证所有子命令
	for cmdIndex, cmd := range subCmds {
		// 检查子命令是否为nil
		if cmd == nil {
			errors = append(errors, qerr.NewValidationErrorf("subCmd at index %d cannot be nil", cmdIndex))
			continue
		}

		// 执行子命令的验证方法
		if err := c.validateSubCmd(cmd); err != nil {
			errors = append(errors, fmt.Errorf("invalid subcommand %s: %w", getCmdIdentifier(cmd), err))
			continue
		}
		validCmds = append(validCmds, cmd)
	}

	// 如果有验证错误, 返回所有错误信息
	if len(errors) > 0 {
		return qerr.NewValidationErrorf("%s: %v", qerr.ErrAddSubCommandFailed, qerr.JoinErrors(errors))
	}

	// 预分配临时切片(容量=validCmds长度, 避免多次扩容)
	tempList := make([]*Cmd, 0, len(validCmds))

	// 添加阶段 - 仅处理通过验证的命令
	for _, cmd := range validCmds {
		// 设置子命令的父命令指针
		cmd.parentCmd = c

		// 将子命令的长名称和实例关联
		if cmd.LongName() != "" {
			c.subCmdMap[cmd.LongName()] = cmd
		}

		// 将子命令的短名称和实例关联
		if cmd.ShortName() != "" {
			c.subCmdMap[cmd.ShortName()] = cmd
		}

		// 先添加到临时切片
		tempList = append(tempList, cmd)
	}

	// 一次性合并到目标切片
	c.subCmds = append(c.subCmds, tempList...)

	return nil
}

// Parse 完整解析命令行参数(含子命令处理)
//
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
//
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

// SetEnableCompletion 设置是否启用自动补全, 只能在根命令上启用
//
// 参数:
//   - enable: true表示启用补全,false表示禁用
func (c *Cmd) SetEnableCompletion(enable bool) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

	// 只在根命令上启用自动补全
	if c.parentCmd != nil {
		return
	}

	// 设置启用状态
	c.enableCompletion = enable
}

// FlagRegistry 获取标志注册表的只读访问
//
// 返回值:
// - *flags.FlagRegistry: 标志注册表的只读访问
func (c *Cmd) FlagRegistry() *flags.FlagRegistry {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()
	return c.flagRegistry
}

// SetVersion 设置版本信息
func (c *Cmd) SetVersion(version string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()
	c.userInfo.version = version
}

// GetVersion 获取版本信息
func (c *Cmd) GetVersion() string {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()
	return c.userInfo.version
}

// SetModuleHelps 设置自定义模块帮助信息
func (c *Cmd) SetModuleHelps(moduleHelps string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()
	c.userInfo.moduleHelps = moduleHelps
}

// GetModuleHelps 获取自定义模块帮助信息
func (c *Cmd) GetModuleHelps() string {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()
	return c.userInfo.moduleHelps
}

// SetLogoText 设置logo文本
func (c *Cmd) SetLogoText(logoText string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()
	c.userInfo.logoText = logoText
}

// GetLogoText 获取logo文本
func (c *Cmd) GetLogoText() string {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()
	return c.userInfo.logoText
}

// GetUseChinese 获取是否使用中文帮助信息
func (c *Cmd) GetUseChinese() bool {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()
	return c.userInfo.useChinese
}

// SetUseChinese 设置是否使用中文帮助信息
func (c *Cmd) SetUseChinese(useChinese bool) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()
	c.userInfo.useChinese = useChinese
}

// GetNotes 获取所有备注信息
func (c *Cmd) GetNotes() []string {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()
	// 返回切片副本而非原始引用
	notes := make([]string, len(c.userInfo.notes))
	copy(notes, c.userInfo.notes)
	return notes
}

// Name 获取命令名称
//
// 返回值:
// - 优先返回长名称, 如果长名称不存在则返回短名称
func (c *Cmd) Name() string {
	if c.LongName() != "" {
		return c.LongName()
	}

	return c.ShortName()
}

// LongName 返回命令长名称
func (c *Cmd) LongName() string { return c.userInfo.longName }

// ShortName 返回命令短名称
func (c *Cmd) ShortName() string { return c.userInfo.shortName }

// GetDescription 返回命令描述
func (c *Cmd) GetDescription() string {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()
	return c.userInfo.description
}

// SetDescription 设置命令描述
func (c *Cmd) SetDescription(desc string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()
	c.userInfo.description = desc
}

// GetHelp 返回命令用法帮助信息
func (c *Cmd) GetHelp() string {
	// 获取读锁
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()

	// 生成帮助信息或返回用户设置的帮助信息
	return generateHelpInfo(c)
}

// SetUsageSyntax 设置自定义命令用法
func (c *Cmd) SetUsageSyntax(usageSyntax string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()
	c.userInfo.usageSyntax = usageSyntax
}

// GetUsageSyntax 获取自定义命令用法
func (c *Cmd) GetUsageSyntax() string {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()
	return c.userInfo.usageSyntax
}

// SetHelp 设置用户自定义命令帮助信息
func (c *Cmd) SetHelp(help string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()
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

	// 清理路径并检查有效性
	cleanPath := filepath.Clean(filePath)
	if cleanPath == "" || strings.TrimSpace(cleanPath) == "" {
		return fmt.Errorf("file path cannot be empty or contain only whitespace")
	}

	// 直接读取文件内容并处理可能的错误（包括文件不存在的情况）
	content, err := os.ReadFile(cleanPath)
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

// AddNote 添加备注信息到命令
func (c *Cmd) AddNote(note string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()
	c.userInfo.notes = append(c.userInfo.notes, note)
}

// AddExample 为命令添加使用示例
// description: 示例描述
// usage: 示例使用方式
func (c *Cmd) AddExample(e ExampleInfo) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()
	// 添加到使用示例列表中
	c.userInfo.examples = append(c.userInfo.examples, e)
}

// GetExamples 获取所有使用示例
// 返回示例切片的副本，防止外部修改
func (c *Cmd) GetExamples() []ExampleInfo {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()
	examples := make([]ExampleInfo, len(c.userInfo.examples))
	copy(examples, c.userInfo.examples)
	return examples
}

// Args 获取非标志参数切片
func (c *Cmd) Args() []string {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()
	// 返回参数切片副本
	args := make([]string, len(c.args))
	copy(args, c.args)
	return args
}

// Arg 获取指定索引的非标志参数
func (c *Cmd) Arg(i int) string {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()
	// 返回参数
	if i >= 0 && i < len(c.args) {
		return c.args[i]
	}
	return ""
}

// NArg 获取非标志参数的数量
func (c *Cmd) NArg() int {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()
	return len(c.args)
}

// NFlag 获取标志的数量
func (c *Cmd) NFlag() int {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()
	return c.fs.NFlag()
}

// FlagExists 检查指定名称的标志是否存在
func (c *Cmd) FlagExists(name string) bool {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()

	// 检查标志是否存在
	if _, exists := c.flagRegistry.GetByName(name); exists {
		return true
	}

	return false
}

// PrintHelp 打印命令的帮助信息, 优先打印用户的帮助信息, 否则自动生成帮助信息
//
// 注意:
//
//	打印帮助信息时, 不会自动退出程序
func (c *Cmd) PrintHelp() {
	// 打印帮助信息
	fmt.Println(c.GetHelp())
}

// CmdExists 检查子命令是否存在
//
// 参数:
//   - cmdName: 子命令名称
//
// 返回:
//   - bool: 子命令是否存在
func (c *Cmd) CmdExists(cmdName string) bool {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()
	// 检查子命令是否存在
	_, ok := c.subCmdMap[cmdName]
	return ok
}

// IsParsed 检查命令是否已完成解析
//
// 返回值:
//
//   - bool: 解析状态,true表示已解析(无论成功失败), false表示未解析
func (c *Cmd) IsParsed() bool {
	return c.parsed.Load()
}

// SetExitOnBuiltinFlags 设置是否在解析内置参数时退出
// 默认情况下为true, 当解析到内置参数时, QFlag将退出程序
//
// 参数:
//   - exit: 是否退出
func (c *Cmd) SetExitOnBuiltinFlags(exit bool) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()
	c.exitOnBuiltinFlags = exit
}
