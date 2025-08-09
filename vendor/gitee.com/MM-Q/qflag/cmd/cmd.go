// Package cmd 命令结构体和核心功能实现
// 本文件定义了Cmd结构体，提供命令行解析、子命令管理、标志注册等核心功能。
// Cmd作为适配器连接内部函数式API和外部面向对象API。
package cmd

import (
	"flag"
	"fmt"
	"os"

	"gitee.com/MM-Q/qflag/flags"
	"gitee.com/MM-Q/qflag/internal/help"
	"gitee.com/MM-Q/qflag/internal/types"
	"gitee.com/MM-Q/qflag/internal/validator"
	"gitee.com/MM-Q/qflag/qerr"
)

// ExampleInfo 导出示例信息类型
type ExampleInfo = types.ExampleInfo

// Cmd 简化的命令结构体，作为适配器连接内部函数式API和外部面向对象API
type Cmd struct {
	ctx *types.CmdContext // 内部上下文，包含所有状态
}

// New 创建新的命令实例(NewCmd的简写)
var New = NewCmd

// ================================================================================
// 操作方法 - 解析与管理 (17个)
// ================================================================================

// NewCmd 创建新的命令实例
//
// 参数:
//   - longName: 命令的全称(如: ls, rm, mkdir 等)
//   - shortName: 命令的简称(如: l, r, m 等)
//   - errorHandling: 标志解析错误处理策略
//
// 返回值:
//   - *Cmd: 新创建的命令实例
//
// errorHandling可选值:
//   - flag.ContinueOnError: 遇到错误时继续解析, 并将错误返回
//   - flag.ExitOnError: 遇到错误时立即退出程序, 并将错误返回
//   - flag.PanicOnError: 遇到错误时立即触发panic, 并将错误返回
func NewCmd(longName, shortName string, errorHandling flag.ErrorHandling) *Cmd {
	// 创建内部上下文
	ctx := types.NewCmdContext(longName, shortName, errorHandling)

	// 创建命令实例
	cmd := &Cmd{ctx: ctx}

	// 注册内置标志help
	cmd.BoolVar(cmd.ctx.BuiltinFlags.Help, flags.HelpFlagName, flags.HelpFlagShortName, false, flags.HelpFlagUsage)

	// 添加到内置标志名称映射
	cmd.ctx.BuiltinFlags.NameMap.Store(flags.HelpFlagName, true)
	cmd.ctx.BuiltinFlags.NameMap.Store(flags.HelpFlagShortName, true)

	return cmd
}

// Parse 完整解析命令行参数(含子命令处理)
//
// 主要功能：
//  1. 解析当前命令的长短标志及内置标志
//  2. 自动检测并解析子命令及其参数(若存在)
//  3. 验证枚举类型标志的有效性
//
// 参数：
//   - args: 原始命令行参数切片(包含可能的子命令及参数)
//
// 返回值：
//   - error: 解析过程中遇到的错误(如标志格式错误、子命令解析失败等)
//
// 注意事项：
//   - 每个Cmd实例仅会被解析一次(线程安全)
//   - 若检测到子命令, 会将剩余参数传递给子命令的Parse方法
//   - 处理内置标志执行逻辑
func (c *Cmd) Parse(args []string) (err error) {
	shouldExit, err := c.parseCommon(args, true)
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
//   - args: 原始命令行参数切片(子命令及后续参数会被忽略)
//
// 返回值：
//   - error: 解析过程中遇到的错误(如标志格式错误等)
//
// 注意事项：
//   - 每个Cmd实例仅会被解析一次(线程安全)
//   - 不会处理任何子命令, 所有参数均视为当前命令的标志或位置参数
//   - 处理内置标志逻辑
func (c *Cmd) ParseFlagsOnly(args []string) (err error) {
	shouldExit, err := c.parseCommon(args, false)
	if shouldExit {
		os.Exit(0)
	}
	return err
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
	// 检查子命令是否为空
	if len(subCmds) == 0 {
		return qerr.NewValidationError("subCmds list cannot be empty")
	}

	// 验证阶段 - 在获取锁之前进行，避免死锁
	var errors []error
	validCmds := make([]*Cmd, 0, len(subCmds)) // 预分配空间

	// 验证所有子命令（无锁操作）
	for cmdIndex, cmd := range subCmds {
		// 检查子命令是否为nil
		if cmd == nil {
			errors = append(errors, qerr.NewValidationErrorf("subCmd at index %d cannot be nil", cmdIndex))
			continue
		}

		// 执行子命令的验证方法（无锁操作）
		if err := validator.ValidateSubCommand(c.ctx, cmd.ctx); err != nil {
			errors = append(errors, fmt.Errorf("invalid subcommand %s: %w", validator.GetCmdIdentifier(cmd.ctx), err))
			continue
		}
		validCmds = append(validCmds, cmd)
	}

	// 如果有验证错误, 返回所有错误信息
	if len(errors) > 0 {
		return qerr.NewValidationErrorf("%s: %v", qerr.ErrAddSubCommandFailed, qerr.JoinErrors(errors))
	}

	// 获取锁进行实际的添加操作
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()

	// 检查子命令map是否为nil
	if c.ctx.SubCmdMap == nil {
		return qerr.NewValidationError("subCmdMap cannot be nil")
	}

	// 检查子命令数组是否为nil
	if c.ctx.SubCmds == nil {
		return qerr.NewValidationError("subCmds cannot be nil")
	}

	// 预分配临时切片(容量=validCmds长度, 避免多次扩容)
	tempList := make([]*types.CmdContext, 0, len(validCmds))

	// 添加阶段 - 仅处理通过验证的命令
	for _, cmd := range validCmds {
		// 设置子命令的父命令指针
		cmd.ctx.Parent = c.ctx

		// 将子命令的长名称和实例关联
		if cmd.ctx.LongName != "" {
			c.ctx.SubCmdMap[cmd.ctx.LongName] = cmd.ctx
		}

		// 将子命令的短名称和实例关联
		if cmd.ctx.ShortName != "" {
			c.ctx.SubCmdMap[cmd.ctx.ShortName] = cmd.ctx
		}

		// 先添加到临时切片
		tempList = append(tempList, cmd.ctx)
	}

	// 一次性合并到目标切片
	c.ctx.SubCmds = append(c.ctx.SubCmds, tempList...)

	return nil
}

// SubCmdMap 返回子命令映射表(长命令名+短命令名)
//
// 返回值:
//   - map[string]*Cmd: 子命令映射表
func (c *Cmd) SubCmdMap() map[string]*Cmd {
	c.ctx.Mutex.RLock()
	defer c.ctx.Mutex.RUnlock()

	// 检查子命令映射表是否为空
	if len(c.ctx.SubCmdMap) == 0 {
		return nil
	}

	// 返回map副本避免外部修改
	subCmdMap := make(map[string]*Cmd, len(c.ctx.SubCmdMap))

	// 遍历子命令映射表, 将每个子命令复制到新的map中
	for name, ctx := range c.ctx.SubCmdMap {
		subCmdMap[name] = &Cmd{ctx: ctx}
	}
	return subCmdMap
}

// SubCmds 返回子命令切片
//
// 返回值:
//   - []*Cmd: 子命令切片
func (c *Cmd) SubCmds() []*Cmd {
	c.ctx.Mutex.RLock()
	defer c.ctx.Mutex.RUnlock()

	// 检查子命令是否为空
	if len(c.ctx.SubCmds) == 0 {
		return nil
	}

	// 创建一个切片副本
	result := make([]*Cmd, len(c.ctx.SubCmds))

	// 拷贝子命令切片
	for i, ctx := range c.ctx.SubCmds {
		result[i] = &Cmd{ctx: ctx}
	}

	return result
}

// FlagRegistry 获取标志注册表的只读访问
//
// 返回值:
// - *flags.FlagRegistry: 标志注册表的只读访问
func (c *Cmd) FlagRegistry() *flags.FlagRegistry {
	c.ctx.Mutex.RLock()
	defer c.ctx.Mutex.RUnlock()
	return c.ctx.FlagRegistry
}

// Name 获取命令名称
//
// 返回值:
//   - string: 命令名称
//
// 说明:
//   - 优先返回长名称, 如果长名称不存在则返回短名称
func (c *Cmd) Name() string {
	if c.ctx.LongName != "" {
		return c.ctx.LongName
	}

	return c.ctx.ShortName
}

// LongName 返回命令长名称
//
// 返回值:
//   - string: 命令长名称
func (c *Cmd) LongName() string { return c.ctx.LongName }

// ShortName 返回命令短名称
//
// 返回值:
//   - string: 命令短名称
func (c *Cmd) ShortName() string { return c.ctx.ShortName }

// Args 获取非标志参数切片
//
// 返回值:
//   - []string: 参数切片
func (c *Cmd) Args() []string {
	c.ctx.Mutex.RLock()
	defer c.ctx.Mutex.RUnlock()

	// 检查参数是否为空
	if len(c.ctx.Args) == 0 {
		return nil
	}

	// 返回参数切片副本
	args := make([]string, len(c.ctx.Args))
	copy(args, c.ctx.Args)
	return args
}

// Arg 获取指定索引的非标志参数
//
// 参数:
//   - i: 参数索引
//
// 返回值:
//   - string: 指定索引位置的非标志参数；若索引越界，则返回空字符串
func (c *Cmd) Arg(i int) string {
	c.ctx.Mutex.RLock()
	defer c.ctx.Mutex.RUnlock()

	// 返回参数
	if i >= 0 && i < len(c.ctx.Args) {
		return c.ctx.Args[i]
	}
	return ""
}

// NArg 获取非标志参数的数量
//
// 返回值:
//   - int: 参数数量
func (c *Cmd) NArg() int {
	c.ctx.Mutex.RLock()
	defer c.ctx.Mutex.RUnlock()
	return len(c.ctx.Args)
}

// NFlag 获取标志的数量
//
// 返回值:
//   - int: 标志数量
func (c *Cmd) NFlag() int {
	c.ctx.Mutex.RLock()
	defer c.ctx.Mutex.RUnlock()

	// 返回独立的标志数量
	return c.ctx.FlagRegistry.GetFlagMetaCount()
}

// FlagExists 检查指定名称的标志是否存在
//
// 参数:
//   - name: 标志名称
//
// 返回值:
//   - bool: 标志是否存在
func (c *Cmd) FlagExists(name string) bool {
	c.ctx.Mutex.RLock()
	defer c.ctx.Mutex.RUnlock()

	// 检查标志是否存在
	if _, exists := c.ctx.FlagRegistry.GetByName(name); exists {
		return true
	}

	return false
}

// PrintHelp 打印命令的帮助信息, 优先打印用户的帮助信息, 否则自动生成帮助信息
//
// 注意:
//   - 打印帮助信息时, 不会自动退出程序
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
	c.ctx.Mutex.RLock()
	defer c.ctx.Mutex.RUnlock()

	// 检查子命令名称是否为空
	if cmdName == "" {
		return false
	}

	// 检查子命令是否存在
	_, ok := c.ctx.SubCmdMap[cmdName]
	return ok
}

// IsParsed 检查命令是否已完成解析
//
// 返回值:
//   - bool: 解析状态,true表示已解析(无论成功失败), false表示未解析
func (c *Cmd) IsParsed() bool {
	return c.ctx.Parsed.Load()
}

// ================================================================================
// Get 方法 - 获取配置信息(9个)
// ================================================================================

// GetVersion 获取版本信息
//
// 返回值:
// - string: 版本信息
func (c *Cmd) GetVersion() string {
	c.ctx.Mutex.RLock()
	defer c.ctx.Mutex.RUnlock()
	return c.ctx.Config.Version
}

// GetModuleHelps 获取自定义模块帮助信息
//
// 返回值:
//   - string: 自定义模块帮助信息
func (c *Cmd) GetModuleHelps() string {
	c.ctx.Mutex.RLock()
	defer c.ctx.Mutex.RUnlock()
	return c.ctx.Config.ModuleHelps
}

// GetLogoText 获取logo文本
//
// 返回值:
//   - string: logo文本字符串
func (c *Cmd) GetLogoText() string {
	c.ctx.Mutex.RLock()
	defer c.ctx.Mutex.RUnlock()
	return c.ctx.Config.LogoText
}

// GetUseChinese 获取是否使用中文帮助信息
//
// 返回值:
//   - bool: 是否使用中文帮助信息
func (c *Cmd) GetUseChinese() bool {
	c.ctx.Mutex.RLock()
	defer c.ctx.Mutex.RUnlock()
	return c.ctx.Config.UseChinese
}

// GetNotes 获取所有备注信息
//
// 返回:
//   - 备注信息列表
func (c *Cmd) GetNotes() []string {
	c.ctx.Mutex.RLock()
	defer c.ctx.Mutex.RUnlock()
	// 返回切片副本而非原始引用
	notes := make([]string, len(c.ctx.Config.Notes))
	copy(notes, c.ctx.Config.Notes)
	return notes
}

// GetDescription 返回命令描述
//
// 返回值:
//   - string: 命令描述
func (c *Cmd) GetDescription() string {
	c.ctx.Mutex.RLock()
	defer c.ctx.Mutex.RUnlock()
	return c.ctx.Config.Description
}

// GetHelp 返回命令用法帮助信息
//
// 返回值:
//   - string: 命令用法帮助信息
func (c *Cmd) GetHelp() string {
	// 获取读锁
	c.ctx.Mutex.RLock()
	defer c.ctx.Mutex.RUnlock()

	// 生成帮助信息或返回用户设置的帮助信息
	return help.GenerateHelp(c.ctx)
}

// GetUsageSyntax 获取自定义命令用法
//
// 返回值:
//   - string: 自定义命令用法
func (c *Cmd) GetUsageSyntax() string {
	c.ctx.Mutex.RLock()
	defer c.ctx.Mutex.RUnlock()
	return c.ctx.Config.UsageSyntax
}

// GetExamples 获取所有使用示例
//
// 返回:
//   - []ExampleInfo: 使用示例列表
func (c *Cmd) GetExamples() []ExampleInfo {
	c.ctx.Mutex.RLock()
	defer c.ctx.Mutex.RUnlock()
	examples := make([]ExampleInfo, len(c.ctx.Config.Examples))

	for i, e := range c.ctx.Config.Examples {
		examples[i] = ExampleInfo(e)
	}

	return examples
}

// ================================================================================
// Set 方法 - 设置配置信息(15个)
// ================================================================================

// SetExitOnBuiltinFlags 设置是否在解析内置参数时退出
// 默认情况下为true, 当解析到内置参数时, QFlag将退出程序
//
// 参数:
//   - exit: 是否退出
func (c *Cmd) SetExitOnBuiltinFlags(exit bool) {
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()
	c.ctx.Config.ExitOnBuiltinFlags = exit
}

// SetEnableCompletion 设置是否启用自动补全, 只能在根命令上启用
//
// 参数:
//   - enable: true表示启用补全,false表示禁用
func (c *Cmd) SetEnableCompletion(enable bool) {
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()

	// 只在根命令上启用自动补全
	if c.ctx.Parent != nil {
		return
	}

	// 设置启用状态
	c.ctx.Config.EnableCompletion = enable
}

// SetVersion 设置版本信息
//
// 参数:
//   - version: 版本信息
func (c *Cmd) SetVersion(version string) {
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()

	// 检查版本信息是否为空
	if version == "" {
		return
	}

	// 设置版本信息
	c.ctx.Config.Version = version
}

// SetVersionf 设置版本信息
//
// 参数:
//   - format: 版本信息格式字符串
//   - args: 格式化参数
func (c *Cmd) SetVersionf(format string, args ...any) {
	c.SetVersion(fmt.Sprintf(format, args...))
}

// SetModuleHelps 设置自定义模块帮助信息
//
// 参数:
//   - moduleHelps: 自定义模块帮助信息
func (c *Cmd) SetModuleHelps(moduleHelps string) {
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()
	c.ctx.Config.ModuleHelps = moduleHelps
}

// SetLogoText 设置logo文本
//
// 参数:
//   - logoText: logo文本字符串
func (c *Cmd) SetLogoText(logoText string) {
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()
	c.ctx.Config.LogoText = logoText
}

// SetUseChinese 设置是否使用中文帮助信息
//
// 参数:
//   - useChinese: 是否使用中文帮助信息
func (c *Cmd) SetUseChinese(useChinese bool) {
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()
	c.ctx.Config.UseChinese = useChinese
}

// SetDescription 设置命令描述
//
// 参数:
//   - desc: 命令描述
func (c *Cmd) SetDescription(desc string) {
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()
	c.ctx.Config.Description = desc
}

// SetHelp 设置用户自定义命令帮助信息
//
// 参数:
//   - help: 用户自定义命令帮助信息
func (c *Cmd) SetHelp(help string) {
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()
	c.ctx.Config.Help = help
}

// SetUsageSyntax 设置自定义命令用法
//
// 参数:
//   - usageSyntax: 自定义命令用法
func (c *Cmd) SetUsageSyntax(usageSyntax string) {
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()
	c.ctx.Config.UsageSyntax = usageSyntax
}

// AddNote 添加备注信息到命令
//
// 参数:
//   - note: 备注信息
func (c *Cmd) AddNote(note string) {
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()
	c.ctx.Config.Notes = append(c.ctx.Config.Notes, note)
}

// AddNotes 添加备注信息切片到命令
//
// 参数:
//   - notes: 备注信息列表
func (c *Cmd) AddNotes(notes []string) {
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()
	c.ctx.Config.Notes = append(c.ctx.Config.Notes, notes...)
}

// AddExample 为命令添加使用示例
//
// 参数:
//   - desc: 示例描述
//   - usage: 示例用法
func (c *Cmd) AddExample(desc, usage string) {
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()

	// 检查描述和用法是否为空
	if desc == "" || usage == "" {
		return
	}

	// 新建示例信息
	e := ExampleInfo{
		Description: desc,
		Usage:       usage,
	}

	// 添加到使用示例列表中
	c.ctx.Config.Examples = append(c.ctx.Config.Examples, e)
}

// AddExamples 为命令添加使用示例列表
//
// 参数:
//   - examples: 示例信息列表
func (c *Cmd) AddExamples(examples []ExampleInfo) {
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()

	// 检查示例信息列表是否为空
	if len(examples) == 0 {
		return
	}

	// 添加到使用示例列表中
	c.ctx.Config.Examples = append(c.ctx.Config.Examples, examples...)
}

// ================================================================================
// 链式调用方法 - 用于构建器模式，提供更流畅的API体验(14个)
// ================================================================================

// WithDescription 设置命令描述（链式调用）
//
// 参数:
//   - desc: 命令描述
//
// 返回值:
//   - *Cmd: 返回命令实例，支持链式调用
func (c *Cmd) WithDescription(desc string) *Cmd {
	c.SetDescription(desc)
	return c
}

// WithVersion 设置版本信息（链式调用）
//
// 参数:
//   - version: 版本信息
//
// 返回值:
//   - *Cmd: 返回命令实例，支持链式调用
func (c *Cmd) WithVersion(version string) *Cmd {
	c.SetVersion(version)
	return c
}

// WithVersionf 设置版本信息（链式调用，支持格式化）
//
// 参数:
//   - format: 版本信息格式字符串
//   - args: 格式化参数
//
// 返回值:
//   - *Cmd: 返回命令实例，支持链式调用
func (c *Cmd) WithVersionf(format string, args ...any) *Cmd {
	c.SetVersionf(format, args...)
	return c
}

// WithUseChinese 设置是否使用中文帮助信息（链式调用）
//
// 参数:
//   - useChinese: 是否使用中文帮助信息
//
// 返回值:
//   - *Cmd: 返回命令实例，支持链式调用
func (c *Cmd) WithUseChinese(useChinese bool) *Cmd {
	c.SetUseChinese(useChinese)
	return c
}

// WithUsageSyntax 设置自定义命令用法（链式调用）
//
// 参数:
//   - usageSyntax: 自定义命令用法
//
// 返回值:
//   - *Cmd: 返回命令实例，支持链式调用
func (c *Cmd) WithUsageSyntax(usageSyntax string) *Cmd {
	c.SetUsageSyntax(usageSyntax)
	return c
}

// WithLogoText 设置logo文本（链式调用）
//
// 参数:
//   - logoText: logo文本字符串
//
// 返回值:
//   - *Cmd: 返回命令实例，支持链式调用
func (c *Cmd) WithLogoText(logoText string) *Cmd {
	c.SetLogoText(logoText)
	return c
}

// WithHelp 设置用户自定义命令帮助信息（链式调用）
//
// 参数:
//   - help: 用户自定义命令帮助信息
//
// 返回值:
//   - *Cmd: 返回命令实例，支持链式调用
func (c *Cmd) WithHelp(help string) *Cmd {
	c.SetHelp(help)
	return c
}

// WithNote 添加备注信息到命令（链式调用）
//
// 参数:
//   - note: 备注信息
//
// 返回值:
//   - *Cmd: 返回命令实例，支持链式调用
func (c *Cmd) WithNote(note string) *Cmd {
	c.AddNote(note)
	return c
}

// WithNotes 添加备注信息切片到命令（链式调用）
//
// 参数:
//   - notes: 备注信息列表
//
// 返回值:
//   - *Cmd: 返回命令实例，支持链式调用
func (c *Cmd) WithNotes(notes []string) *Cmd {
	c.AddNotes(notes)
	return c
}

// WithExample 为命令添加使用示例（链式调用）
//
// 参数:
//   - desc: 示例描述
//   - usage: 示例用法
//
// 返回值:
//   - *Cmd: 返回命令实例，支持链式调用
func (c *Cmd) WithExample(desc, usage string) *Cmd {
	c.AddExample(desc, usage)
	return c
}

// WithExamples 添加使用示例列表到命令（链式调用）
//
// 参数:
//   - examples: 示例信息列表
//
// 返回值:
//   - *Cmd: 返回命令实例，支持链式调用
func (c *Cmd) WithExamples(examples []ExampleInfo) *Cmd {
	c.AddExamples(examples)
	return c
}

// WithExitOnBuiltinFlags 设置是否在解析内置参数时退出（链式调用）
//
// 参数:
//   - exit: 是否退出
//
// 返回值:
//   - *Cmd: 返回命令实例，支持链式调用
func (c *Cmd) WithExitOnBuiltinFlags(exit bool) *Cmd {
	c.SetExitOnBuiltinFlags(exit)
	return c
}

// WithEnableCompletion 设置是否启用自动补全（链式调用）
//
// 参数:
//   - enable: true表示启用补全,false表示禁用
//
// 返回值:
//   - *Cmd: 返回命令实例，支持链式调用
func (c *Cmd) WithEnableCompletion(enable bool) *Cmd {
	c.SetEnableCompletion(enable)
	return c
}

// WithModuleHelps 设置自定义模块帮助信息（链式调用）
//
// 参数:
//   - moduleHelps: 自定义模块帮助信息
//
// 返回值:
//   - *Cmd: 返回命令实例，支持链式调用
func (c *Cmd) WithModuleHelps(moduleHelps string) *Cmd {
	c.SetModuleHelps(moduleHelps)
	return c
}
