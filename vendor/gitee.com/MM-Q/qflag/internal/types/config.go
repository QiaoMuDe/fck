package types

// DepType 依赖关系类型
type DepType int

const (
	// DepMutex 互斥依赖
	// 当触发标志被设置时，目标标志不能被设置
	DepMutex DepType = iota

	// DepRequired 必需依赖
	// 当触发标志被设置时，所有目标标志必须被设置
	DepRequired
)

// String 返回依赖类型的字符串表示
func (d DepType) String() string {
	switch d {
	case DepMutex:
		return "mutex"
	case DepRequired:
		return "required"
	default:
		return "unknown"
	}
}

// FlagDependency 标志依赖关系定义
//
// FlagDependency 定义了当某个标志（触发标志）被设置时，
// 对其他标志（目标标志）的约束条件。
//
// 字段说明:
//   - Name: 依赖关系名称，用于错误提示和标识
//   - Trigger: 触发标志的名称，当此标志被设置时触发依赖检查
//   - Targets: 目标标志名称列表，这些标志会受到约束
//   - Type: 依赖关系类型（互斥或必需）
//
// 使用场景:
//   - 远程模式与本地路径互斥 (trigger="remote", targets=["local-path"], type=DepMutex)
//   - SSL模式需要证书和密钥 (trigger="ssl", targets=["cert","key"], type=DepRequired)
//   - 配置文件模式与其他配置互斥 (trigger="config", targets=["port","host"], type=DepMutex)
type FlagDependency struct {
	Name    string   // 依赖关系名称，用于错误提示和标识
	Trigger string   // 触发标志名称
	Targets []string // 目标标志名称列表
	Type    DepType  // 依赖关系类型
}

// MutexGroup 互斥组定义
//
// MutexGroup 定义了一组互斥的标志, 其中最多只能有一个被设置。
// 当用户设置了互斥组中的多个标志时, 解析器会返回错误。
//
// 字段说明:
//   - Name: 互斥组名称, 用于错误提示和标识
//   - Flags: 互斥组中的标志名称列表
//   - AllowNone: 是否允许一个都不设置
//
// 使用场景:
//   - 输出格式互斥 (如 --json 和 --xml 不能同时使用)
//   - 操作模式互斥 (如 --create 和 --update 不能同时使用)
//   - 必选选项 (如必须指定 --file 或 --url 中的一个)
type MutexGroup struct {
	Name      string   // 互斥组名称, 用于错误提示和标识
	Flags     []string // 互斥组中的标志名称列表
	AllowNone bool     // 是否允许一个都不设置
}

// RequiredGroup 必需组定义
//
// RequiredGroup 定义了一组必需的标志, 其中所有标志都必须被设置。
// 当用户没有设置必需组中的某些标志时, 解析器会返回错误。
//
// 字段说明:
//   - Name: 必需组名称, 用于错误提示和标识
//   - Flags: 必需组中的标志名称列表
//   - Conditional: 是否为条件性必需组, 如果为true, 则只有当组中任何一个标志被设置时, 才要求所有标志都被设置
//
// 使用场景:
//   - 数据库连接配置 (host, port, user, pass)
//   - API 认证配置 (api-key, api-secret)
//   - 文件上传配置 (file-path, upload-url)
//   - 条件性配置 (如果使用了任何一个标志, 则必须使用所有标志)
type RequiredGroup struct {
	Name        string   // 必需组名称, 用于错误提示和标识
	Flags       []string // 必需组中的标志名称列表
	Conditional bool     // 是否为条件性必需组
}

// CmdConfig 命令配置类型
type CmdConfig struct {
	Version           string            // 版本号
	UseChinese        bool              // 是否使用中文
	EnvPrefix         string            // 环境变量前缀
	UsageSyntax       string            // 命令使用语法
	Example           map[string]string // 示例使用, key为描述, value为示例命令
	Notes             []string          // 注意事项
	LogoText          string            // 命令logo文本
	MutexGroups       []MutexGroup      // 互斥组列表
	RequiredGroups    []RequiredGroup   // 必需组列表
	FlagDependencies  []FlagDependency  // 标志依赖关系列表
	Completion        bool              // 是否启用自动补全标志
	DynamicCompletion bool              // 是否启用动态补全
}

// NewCmdConfig 创建新的命令配置
//
// 返回值:
//   - *CmdConfig: 新创建的 CmdConfig 实例, 初始化为零值
func NewCmdConfig() *CmdConfig {
	return &CmdConfig{
		Version:           "",
		UseChinese:        false,
		EnvPrefix:         "",
		UsageSyntax:       "",
		Example:           map[string]string{},
		Notes:             []string{},
		LogoText:          "",
		MutexGroups:       []MutexGroup{},
		RequiredGroups:    []RequiredGroup{},
		FlagDependencies:  []FlagDependency{},
		Completion:        false,
		DynamicCompletion: false,
	}
}

// Clone 克隆命令配置
//
// 返回值:
//   - *CmdConfig: 克隆后的新 CmdConfig 实例
//
// 功能说明:
//   - 创建当前配置的深拷贝
//   - 复制所有字段值
//   - 复制切片和映射时创建新的底层数组/映射
//   - 用于避免配置共享导致的副作用
func (c *CmdConfig) Clone() *CmdConfig {
	if c == nil {
		return nil
	}

	clone := &CmdConfig{
		Version:           c.Version,
		UseChinese:        c.UseChinese,
		EnvPrefix:         c.EnvPrefix,
		UsageSyntax:       c.UsageSyntax,
		LogoText:          c.LogoText,
		Completion:        c.Completion,
		DynamicCompletion: c.DynamicCompletion,
	}

	// 深拷贝 Example 映射
	if len(c.Example) > 0 {
		clone.Example = make(map[string]string, len(c.Example))
		for k, v := range c.Example {
			clone.Example[k] = v
		}
	}

	// 深拷贝 Notes 切片
	if len(c.Notes) > 0 {
		clone.Notes = make([]string, len(c.Notes))
		copy(clone.Notes, c.Notes)
	}

	// 深拷贝 MutexGroups 切片
	if len(c.MutexGroups) > 0 {
		clone.MutexGroups = make([]MutexGroup, len(c.MutexGroups))
		copy(clone.MutexGroups, c.MutexGroups)
	}

	// 深拷贝 RequiredGroups 切片
	if len(c.RequiredGroups) > 0 {
		clone.RequiredGroups = make([]RequiredGroup, len(c.RequiredGroups))
		copy(clone.RequiredGroups, c.RequiredGroups)
	}

	// 深拷贝 FlagDependencies 切片
	if len(c.FlagDependencies) > 0 {
		clone.FlagDependencies = make([]FlagDependency, len(c.FlagDependencies))
		copy(clone.FlagDependencies, c.FlagDependencies)
	}

	return clone
}
