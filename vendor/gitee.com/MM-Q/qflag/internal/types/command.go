package types

// Command 接口定义了命令的核心行为
type Command interface {
	// 基本属性
	Name() string      // 命令名称, 用于匹配和显示
	LongName() string  // 长名称, 用于显示和帮助
	ShortName() string // 短名称, 用于命令行输入
	Desc() string      // 命令描述, 用于帮助显示

	// 标志管理
	AddFlag(flag Flag) error          // 添加一个标志到命令
	AddFlags(flags ...Flag) error     // 添加多个标志到命令
	AddFlagsFrom(flags []Flag) error  // 从切片添加多个标志
	GetFlag(name string) (Flag, bool) // 根据名称获取标志
	Flags() []Flag                    // 获取所有标志
	FlagRegistry() FlagRegistry       // 获取标志注册器

	// 子命令管理
	AddSubCmds(cmds ...Command) error      // 添加多个子命令
	AddSubCmdFrom(cmds []Command) error    // 从切片添加子命令
	GetSubCmd(name string) (Command, bool) // 根据名称获取子命令
	SubCmds() []Command                    // 获取所有子命令
	HasSubCmd(name string) bool            // 是否有指定名称的子命令
	CmdRegistry() CmdRegistry              // 获取子命令注册器

	// 命令层次
	IsRootCmd() bool // 是否为根命令
	Path() string    // 命令的路径, 用于显示和帮助

	// 参数解析
	Parse(args []string) error         // 解析命令行参数
	ParseAndRoute(args []string) error // 解析并路由到子命令
	ParseOnly(args []string) error     // 仅解析参数, 不路由
	IsParsed() bool                    // 是否已解析参数
	SetParsed(parsed bool)             // 设置解析状态

	// 参数访问
	Args() []string        // 获取所有参数
	Arg(index int) string  // 获取指定索引的参数
	NArg() int             // 获取参数数量
	SetArgs(args []string) // 设置参数

	// 执行
	Run() error                    // 执行命令
	SetRun(fn func(Command) error) // 设置执行函数
	HasRunFunc() bool              // 是否有执行函数

	// 帮助信息
	Help() string // 获取命令帮助信息
	PrintHelp()   // 打印命令帮助信息

	// 配置
	SetParser(p Parser)                     // 设置解析器
	SetDesc(desc string)                    // 设置命令描述
	SetVersion(version string)              // 设置命令版本
	SetChinese(useChinese bool)             // 设置是否使用中文
	SetEnvPrefix(prefix string)             // 设置环境变量前缀
	SetUsageSyntax(syntax string)           // 设置命令行语法
	AddExample(title, cmd string)           // 添加一个示例
	AddExamples(examples map[string]string) // 添加多个示例
	AddNote(note string)                    // 添加一条注意事项
	AddNotes(notes []string)                // 添加多条注意事项
	SetLogoText(logo string)                // 设置命令logo文本
	Config() *CmdConfig                     // 获取命令配置
}
