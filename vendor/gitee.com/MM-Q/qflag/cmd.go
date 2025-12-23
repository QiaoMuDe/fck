// package qflag 命令结构体和核心功能实现
// 本文件定义了Cmd结构体, 提供命令行解析、子命令管理、标志注册等核心功能。
// Cmd作为适配器连接内部函数式API和外部面向对象API。
package qflag

import (
	"fmt"
	"os"
	"sync"

	"gitee.com/MM-Q/qflag/flags"
	"gitee.com/MM-Q/qflag/internal/help"
	"gitee.com/MM-Q/qflag/internal/types"
	"gitee.com/MM-Q/qflag/internal/validator"
	"gitee.com/MM-Q/qflag/qerr"
)

// Cmd 命令结构体, 作为适配器连接内部函数式API和外部面向对象API
type Cmd struct {
	ctx       *types.CmdContext // 内部上下文，包含所有状态
	runFunc   func(*Cmd) error  // 存储Run函数, 用于执行命令逻辑
	subCmdMap map[string]*Cmd   // 存储子命令映射
	runMutex  sync.RWMutex      // 保护runFunc的读写锁
}

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
//   - qflag.ContinueOnError: 遇到错误时继续解析, 并将错误返回
//   - qflag.ExitOnError: 遇到错误时立即退出程序, 并将错误返回
//   - qflag.PanicOnError: 遇到错误时立即触发panic, 并将错误返回
func NewCmd(longName, shortName string, errorHandling ErrorHandling) *Cmd {
	// 创建内部上下文
	ctx := types.NewCmdContext(longName, shortName, errorHandling)

	// 创建命令实例
	cmd := &Cmd{
		ctx:       ctx,
		subCmdMap: make(map[string]*Cmd),
		runMutex:  sync.RWMutex{},
	}

	// 注册内置标志help
	cmd.BoolVar(cmd.ctx.BuiltinFlags.Help, flags.HelpFlagName, flags.HelpFlagShortName, false, flags.HelpFlagUsage)

	// 添加到内置标志名称映射
	cmd.ctx.BuiltinFlags.NameMap.Store(flags.HelpFlagName, true)
	cmd.ctx.BuiltinFlags.NameMap.Store(flags.HelpFlagShortName, true)

	return cmd
}

// Parse 完整解析命令行参数(递归解析子命令)
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

	// 如果返回了退出信号, 则需要主动退出程序, 否则通过返回的错误判断是否需要退出
	if shouldExit {
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

	// 如果返回了退出信号, 则需要主动退出程序, 否则通过返回的错误判断是否需要退出
	if shouldExit {
		os.Exit(0)
	}
	return err
}

// AddSubCmd 向当前命令添加一个或多个子命令
//
// 此方法会对所有子命令进行完整性验证，包括名称冲突检查、循环依赖检测等。
// 所有验证通过后，子命令将被注册到当前命令的子命令映射表和列表中。
// 操作过程中会自动设置子命令的父命令引用，确保命令树结构的完整性。
//
// 参数:
//   - subCmds: 要添加的子命令实例指针，支持传入多个子命令进行批量添加
//
// 返回值:
//   - error: 添加过程中的错误信息。如果任何子命令验证失败，将返回包含所有错误详情的聚合错误；
//     如果所有子命令成功添加，返回 nil
//
// 错误类型:
//   - ValidationError: 子命令为空、名称冲突、循环依赖等验证错误
//   - 其他错误: 内部状态异常或系统错误
//
// 使用示例:
//
//	cmd := qflag.NewCmd("parent", "p", "父命令")
//	subCmd1 := qflag.NewCmd("child1", "c1", "子命令1")
//	subCmd2 := qflag.NewCmd("child2", "c2", "子命令2")
//
//	if err := cmd.AddSubCmd(subCmd1, subCmd2); err != nil {
//	    log.Fatal(err)
//	}
func (c *Cmd) AddSubCmd(subCmds ...*Cmd) error {
	// 检查子命令是否为空
	if len(subCmds) == 0 {
		return qerr.NewValidationError("subCmds list cannot be empty")
	}

	// 提前获取锁，覆盖整个验证和添加过程
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()

	// 验证阶段 - 在获取锁之前进行，避免死锁
	var errors []error
	validCmds := make([]*Cmd, 0, len(subCmds)) // 预分配空间

	// 验证所有子命令(无锁操作)
	for cmdIndex, cmd := range subCmds {
		// 检查子命令是否为nil
		if cmd == nil {
			errors = append(errors, qerr.NewValidationErrorf("subCmd at index %d cannot be nil", cmdIndex))
			continue
		}

		// 执行子命令的验证方法(无锁操作)
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
			c.ctx.SubCmdMap[cmd.ctx.LongName] = cmd.ctx // 添加命令到ctx命令映射表
			c.subCmdMap[cmd.ctx.LongName] = cmd         // 添加命令到cmd命令映射表
		}

		// 将子命令的短名称和实例关联
		if cmd.ctx.ShortName != "" {
			c.ctx.SubCmdMap[cmd.ctx.ShortName] = cmd.ctx
			c.subCmdMap[cmd.ctx.ShortName] = cmd
		}

		// 先添加到临时切片
		tempList = append(tempList, cmd.ctx)
	}

	// 一次性合并到目标切片
	c.ctx.SubCmds = append(c.ctx.SubCmds, tempList...)

	return nil
}

// AddSubCmds 向当前命令添加子命令切片的便捷方法
//
// 此方法是 AddSubCmd 的便捷包装，专门用于处理子命令切片。
// 内部直接调用 AddSubCmd 方法，具有相同的验证逻辑和并发安全特性。
//
// 参数:
//   - subCmds: 子命令切片，包含要添加的所有子命令实例指针
//
// 返回值:
//   - error: 添加过程中的错误信息，与 AddSubCmd 返回的错误类型相同
//
// 使用示例:
//
//	cmd := qflag.NewCmd("parent", "p", "父命令")
//	subCmds := []*qflag.Cmd{
//	    qflag.NewCmd("child1", "c1", "子命令1"),
//	    qflag.NewCmd("child2", "c2", "子命令2"),
//	}
//
//	if err := cmd.AddSubCmds(subCmds); err != nil {
//	    log.Fatal(err)
//	}
func (c *Cmd) AddSubCmds(subCmds []*Cmd) error {
	return c.AddSubCmd(subCmds...)
}

// SubCmdMap 返回子命令映射表(长命令名+短命令名)
//
// 返回值:
//   - map[string]*Cmd: 子命令映射表
func (c *Cmd) SubCmdMap() map[string]*Cmd {
	c.ctx.Mutex.RLock()
	defer c.ctx.Mutex.RUnlock()

	// 检查子命令映射表是否为空
	if len(c.subCmdMap) == 0 {
		return nil
	}

	// 返回map副本避免外部修改
	subCmdMap := make(map[string]*Cmd, len(c.subCmdMap))

	// 遍历子命令映射表, 将每个子命令复制到新的map中
	for name, cmd := range c.subCmdMap {
		subCmdMap[name] = cmd
	}
	return subCmdMap
}

// GetSubCmd 根据名称获取子命令实例
//
// 参数:
//   - name: 子命令名称 (长名称或短名称)
//
// 返回值:
//   - *Cmd: 子命令实例，如果找不到则抛出恐慌
func (c *Cmd) GetSubCmd(name string) *Cmd {
	c.ctx.Mutex.RLock()
	defer c.ctx.Mutex.RUnlock()

	// 检查名称是否为空
	if name == "" {
		panic("subcommand name cannot be empty")
	}

	// 从子命令映射表中查找
	cmd, exists := c.subCmdMap[name]
	if !exists {
		panic(fmt.Sprintf("subcommand '%s' not found", name))
	}

	return cmd
}

// FlagRegistry 获取标志注册表的只读访问
//
// 返回值:
// - *FlagRegistry: 标志注册表的只读访问
func (c *Cmd) FlagRegistry() *FlagRegistry {
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
	fmt.Println(c.Help())
}

// HasSubCmd 检查子命令是否存在
//
// 参数:
//   - cmdName: 子命令名称
//
// 返回:
//   - bool: 子命令是否存在
func (c *Cmd) HasSubCmd(cmdName string) bool {
	c.ctx.Mutex.RLock()
	defer c.ctx.Mutex.RUnlock()

	// 检查子命令名称是否为空
	if cmdName == "" {
		return false
	}

	// 使用Cmd结构体的命令映射字段检查子命令是否存在
	_, ok := c.subCmdMap[cmdName]
	return ok
}

// IsParsed 检查命令是否已完成解析
//
// 返回值:
//   - bool: 解析状态,true表示已解析(无论成功失败), false表示未解析
func (c *Cmd) IsParsed() bool {
	return c.ctx.Parsed.Load()
}

// Version 获取版本信息
//
// 返回值:
// - string: 版本信息
func (c *Cmd) Version() string {
	c.ctx.Mutex.RLock()
	defer c.ctx.Mutex.RUnlock()
	return c.ctx.Config.Version
}

// Modules 获取自定义模块帮助信息
//
// 返回值:
//   - string: 自定义模块帮助信息
func (c *Cmd) Modules() string {
	c.ctx.Mutex.RLock()
	defer c.ctx.Mutex.RUnlock()
	return c.ctx.Config.ModuleHelps
}

// Logo 获取logo文本
//
// 返回值:
//   - string: logo文本字符串
func (c *Cmd) Logo() string {
	c.ctx.Mutex.RLock()
	defer c.ctx.Mutex.RUnlock()
	return c.ctx.Config.LogoText
}

// Chinese 获取是否使用中文帮助信息
//
// 返回值:
//   - bool: 是否使用中文帮助信息
func (c *Cmd) Chinese() bool {
	c.ctx.Mutex.RLock()
	defer c.ctx.Mutex.RUnlock()
	return c.ctx.Config.UseChinese
}

// Notes 获取所有备注信息
//
// 返回:
//   - 备注信息列表
func (c *Cmd) Notes() []string {
	c.ctx.Mutex.RLock()
	defer c.ctx.Mutex.RUnlock()
	// 返回切片副本而非原始引用
	notes := make([]string, len(c.ctx.Config.Notes))
	copy(notes, c.ctx.Config.Notes)
	return notes
}

// Desc 返回命令描述
//
// 返回值:
//   - string: 命令描述
func (c *Cmd) Desc() string {
	c.ctx.Mutex.RLock()
	defer c.ctx.Mutex.RUnlock()
	return c.ctx.Config.Desc
}

// Help 返回命令用法帮助信息
//
// 返回值:
//   - string: 命令用法帮助信息
func (c *Cmd) Help() string {
	// 获取读锁
	c.ctx.Mutex.RLock()
	defer c.ctx.Mutex.RUnlock()

	// 生成帮助信息或返回用户设置的帮助信息
	return help.GenerateHelp(c.ctx)
}

// Usage 获取自定义命令用法
//
// 返回值:
//   - string: 自定义命令用法
func (c *Cmd) Usage() string {
	c.ctx.Mutex.RLock()
	defer c.ctx.Mutex.RUnlock()
	return c.ctx.Config.UsageSyntax
}

// Examples 获取所有使用示例
//
// 返回:
//   - []ExampleInfo: 使用示例列表
func (c *Cmd) Examples() []ExampleInfo {
	c.ctx.Mutex.RLock()
	defer c.ctx.Mutex.RUnlock()
	examples := make([]ExampleInfo, len(c.ctx.Config.Examples))

	for i, e := range c.ctx.Config.Examples {
		examples[i] = ExampleInfo(e)
	}

	return examples
}

// SetNoFgExit 设置禁用内置标志自动退出
// 默认情况下为false, 当解析到内置参数时, QFlag将退出程序
//
// 参数:
//   - exit: 是否退出
func (c *Cmd) SetNoFgExit(exit bool) {
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()
	c.ctx.Config.NoFgExit = exit
}

// SetCompletion 设置是否启用自动补全, 只能在根命令上启用
//
// 参数:
//   - enable: true表示启用补全,false表示禁用
func (c *Cmd) SetCompletion(enable bool) {
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()

	// 只在根命令上启用自动补全
	if c.ctx.Parent != nil {
		return
	}

	// 设置启用状态
	c.ctx.Config.Completion = enable
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

// SetModules 设置自定义模块帮助信息
//
// 参数:
//   - moduleHelps: 自定义模块帮助信息
func (c *Cmd) SetModules(moduleHelps string) {
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()
	c.ctx.Config.ModuleHelps = moduleHelps
}

// SetLogo 设置logo文本
//
// 参数:
//   - logoText: logo文本字符串
func (c *Cmd) SetLogo(logoText string) {
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()
	c.ctx.Config.LogoText = logoText
}

// SetChinese 设置是否使用中文帮助信息
//
// 参数:
//   - useChinese: 是否使用中文帮助信息
func (c *Cmd) SetChinese(useChinese bool) {
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()
	c.ctx.Config.UseChinese = useChinese
}

// SetDesc 设置命令描述
//
// 参数:
//   - desc: 命令描述
func (c *Cmd) SetDesc(desc string) {
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()
	c.ctx.Config.Desc = desc
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

// SetUsage 设置自定义命令用法
//
// 参数:
//   - usageSyntax: 自定义命令用法
func (c *Cmd) SetUsage(usageSyntax string) {
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()
	c.ctx.Config.UsageSyntax = usageSyntax
}

// ApplyConfig 批量设置命令配置
// 通过传入一个CmdConfig结构体来一次性设置多个配置项
//
// 参数:
//   - config: 包含所有配置项的CmdConfig结构体
func (c *Cmd) ApplyConfig(config CmdConfig) {
	// 参数验证
	if c == nil || c.ctx == nil {
		return
	}

	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()

	// 设置版本信息
	if config.Version != "" {
		c.ctx.Config.Version = config.Version
	}

	// 设置命令描述
	if config.Desc != "" {
		c.ctx.Config.Desc = config.Desc
	}

	// 设置自定义帮助信息
	if config.Help != "" {
		c.ctx.Config.Help = config.Help
	}

	// 设置自定义用法格式
	if config.UsageSyntax != "" {
		c.ctx.Config.UsageSyntax = config.UsageSyntax
	}

	// 设置模块帮助信息
	if config.ModuleHelps != "" {
		c.ctx.Config.ModuleHelps = config.ModuleHelps
	}

	// 设置logo文本
	if config.LogoText != "" {
		c.ctx.Config.LogoText = config.LogoText
	}

	// 安全地设置备注信息 - 创建新切片避免内存泄漏
	if len(config.Notes) > 0 {
		newNotes := make([]string, len(c.ctx.Config.Notes)+len(config.Notes))
		copy(newNotes, c.ctx.Config.Notes)
		copy(newNotes[len(c.ctx.Config.Notes):], config.Notes)
		c.ctx.Config.Notes = newNotes
	}

	// 安全地设置示例信息 - 创建新切片避免内存泄漏
	if len(config.Examples) > 0 {
		newExamples := make([]types.ExampleInfo, len(c.ctx.Config.Examples)+len(config.Examples))
		copy(newExamples, c.ctx.Config.Examples)
		copy(newExamples[len(c.ctx.Config.Examples):], config.Examples)
		c.ctx.Config.Examples = newExamples
	}

	// 设置是否使用中文帮助信息
	c.ctx.Config.UseChinese = config.UseChinese

	// 设置内置标志是否自动退出
	c.ctx.Config.NoFgExit = config.NoFgExit

	// 设置是否启用自动补全功能 (只允许在根命令上设置)
	if c.ctx.Parent == nil {
		c.ctx.Config.Completion = config.Completion
	}
}

// SetRun 设置命令的执行函数
//
// 参数:
//   - run: 命令执行函数，接收*Cmd作为参数，返回error
func (c *Cmd) SetRun(run func(*Cmd) error) {
	c.runMutex.Lock()
	defer c.runMutex.Unlock()

	if run == nil {
		panic("run function cannot be nil")
	}
	c.runFunc = run
}

// Run 执行在命令设置的run函数, 如果未设置run函数, 则返回错误
//
// 返回值:
//   - error: 执行过程中的错误信息
func (c *Cmd) Run() error {
	c.runMutex.RLock()
	defer c.runMutex.RUnlock()

	// 检查命令是否已解析
	if !c.IsParsed() {
		return qerr.NewValidationError("command must be parsed before execution")
	}

	if c.runFunc == nil {
		return qerr.NewValidationError("no run function set for command")
	}
	return c.runFunc(c)
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
		Desc:  desc,
		Usage: usage,
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

// ParseAndRoute 解析参数并自动路由执行子命令
//
// 参数:
//   - args: 命令行参数列表(通常为 os.Args[1:])
//
// 返回值:
//   - error: 执行过程中遇到的错误
func (c *Cmd) ParseAndRoute(args []string) error {
	// 1. 只解析当前命令的标志参数（不递归子命令）
	if err := c.ParseFlagsOnly(args); err != nil {
		return err
	}

	// 2. 获取非标志参数
	nonFlagArgs := c.Args()

	// 3. 如果没有非标志参数，执行当前命令
	if len(nonFlagArgs) == 0 {
		if c.runFunc != nil {
			return c.Run()
		}
		c.PrintHelp()
		return nil
	}

	// 4. 第一个参数可能是子命令
	cmdName := nonFlagArgs[0]

	// 5. 查找子命令
	subCmdMap := c.SubCmdMap()

	if subCmd, exists := subCmdMap[cmdName]; exists {
		// 递归调用子命令，传递剩余参数
		return subCmd.ParseAndRoute(nonFlagArgs[1:])
	}

	// 6. 如果不是子命令, 则执行当前命令
	if c.runFunc != nil {
		return c.Run()
	}

	// 7. 如果不是子命令, 并且没有执行函数, 则显示帮助信息
	fmt.Printf("unknown command: %s\n", cmdName)
	c.PrintHelp()
	return nil
}
