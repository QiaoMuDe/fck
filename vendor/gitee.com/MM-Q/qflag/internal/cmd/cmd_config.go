// Package cmd 提供命令实现和命令管理功能
//
// cmd_config.go 包含命令配置相关的功能实现
//
// 本文件提供了以下主要功能:
//   - 命令配置管理 (Config)
//   - 命令属性设置方法 (SetDesc, SetHidden, SetDisableFlagParsing, SetVersion 等)
//   - 解析器和运行函数管理 (SetParser, SetRun)
//   - 参数和示例管理 (AddExample, AddNotes)
//   - 选项批量应用 (ApplyOpts)
//
// 主要方法列表:
//   - Config: 获取命令配置
//   - SetDesc/SetHidden/SetDisableFlagParsing: 设置基本属性
//   - SetVersion/SetChinese/SetCompletion: 设置配置选项
//   - SetParser/SetArgs/SetParsed/SetRun: 设置解析器和运行函数
//   - AddExample/AddExamples/AddNote/AddNotes: 添加示例和注释
//   - ApplyOpts: 批量应用选项到命令
//
// 线程安全:
//   - 所有公共方法都使用读写锁保护
//   - 支持并发安全的访问和修改
package cmd

import (
	"fmt"

	"gitee.com/MM-Q/qflag/internal/types"
)

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

	// 使用克隆方法创建配置的深拷贝
	return c.config.Clone()
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

// SetHidden 设置命令是否隐藏
//
// 参数:
//   - hidden: 是否隐藏命令
//
// 功能说明:
//   - 设置命令的隐藏状态
//   - 隐藏的命令不会显示在帮助信息中
//   - 支持并发安全的设置
func (c *Cmd) SetHidden(hidden bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.hidden = hidden
}

// SetDisableFlagParsing 设置是否禁用标志解析
//
// 参数:
//   - disable: 是否禁用标志解析
//
// 功能说明:
//   - 设置命令的标志解析状态
//   - 禁用后所有参数（包括--开头的）都作为位置参数
//   - 支持并发安全的设置
func (c *Cmd) SetDisableFlagParsing(disable bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.disableFlagParsing = disable
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
//   - 默认为 false, 不启用自动补全
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
		panic(fmt.Sprintf("nil parser in '%s'", c.Name()))
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
//   - 支持部分配置 (未设置的属性不会被修改)
//   - 使用defer捕获panic, 转换为错误返回
//
// 应用顺序:
//  1. 基本属性 (Desc、RunFunc)
//  2. 配置选项 (Version、UseChinese、EnvPrefix、UsageSyntax、LogoText)
//  3. 示例和说明 (Examples、Notes)
//  4. 互斥组 (MutexGroups)
//  5. 子命令 (SubCmds)
//
// 错误处理:
//   - 选项为 nil: 返回错误
//   - 添加子命令失败: 返回错误
//   - panic: 转换为错误
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
				err = fmt.Errorf("panic in ApplyOpts '%s': %s", c.Name(), x)
			case error:
				err = fmt.Errorf("panic in ApplyOpts '%s': %w", c.Name(), x)
			default:
				err = fmt.Errorf("panic in ApplyOpts '%s': %v", c.Name(), x)
			}
		}
	}()

	// 验证选项
	if opts == nil {
		return fmt.Errorf("nil opts in '%s'", c.Name())
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
	c.SetHidden(opts.Hidden)
	c.SetDisableFlagParsing(opts.DisableFlagParsing)

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
				return fmt.Errorf("add mutex group '%s' failed in '%s': %w", group.Name, c.Name(), err)
			}
		}
	}

	// 5. 添加必需组 - 调用现有方法
	if len(opts.RequiredGroups) > 0 {
		for _, group := range opts.RequiredGroups {
			if err := c.AddRequiredGroup(group.Name, group.Flags, group.Conditional); err != nil {
				return fmt.Errorf("add required group '%s' failed in '%s': %w", group.Name, c.Name(), err)
			}
		}
	}

	// 6. 添加子命令 - 调用现有方法
	if len(opts.SubCmds) > 0 {
		if err := c.AddSubCmds(opts.SubCmds...); err != nil {
			return fmt.Errorf("add subcommands failed in '%s': %w", c.Name(), err)
		}
	}

	// 7. 自动绑定环境变量
	if opts.AutoBindEnv {
		c.AutoBindAllEnv()
	}

	return err
}

// AutoBindAllEnv 为所有标志自动绑定环境变量
//
// 功能说明:
//   - 遍历命令的所有标志
//   - 为每个标志调用 AutoBindEnv() 方法
//   - 批量设置环境变量绑定
//
// 使用示例:
//
//	cmd.String("host", "h", "主机地址", "localhost")
//	cmd.Uint("port", "p", "端口号", 8080)
//	cmd.AutoBindAllEnv()  // 自动绑定 HOST 和 PORT
//
// 注意事项:
//   - 如果标志没有长名称，会触发 panic
//   - 环境变量名为标志长名称的大写形式
//   - 环境变量前缀在解析时自动拼接
func (c *Cmd) AutoBindAllEnv() {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, f := range c.flagRegistry.List() {
		f.AutoBindEnv()
	}
}
