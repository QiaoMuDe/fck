# Find 命令配置结构体合并方案

## 一、背景

当前 find 命令存在两个配置结构体：
- `find.FindConfig`：CLI 层传入的原始配置（`cmd_find.go`）
- `types.FindConfig`：内部使用的运行时配置（`types/types.go`）

两个结构体有大量重复字段，且需要额外的转换函数 `createFindConfig`，增加了代码复杂度。

---

## 二、合并方案

### 2.1 新结构体定义

创建新文件 `internal/commands/find/config.go`：

```go
// Package find 实现了文件查找命令的配置管理。
// 该文件包含统一的查找配置结构体，整合了 CLI 参数和运行时状态。
package find

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"

	"gitee.com/MM-Q/colorlib"
	common "gitee.com/MM-Q/fck/internal/utils"
)

// FindConfig 统一的查找配置结构体
// 包含 CLI 传入的原始参数和运行时生成的状态
type FindConfig struct {
	// ==================== CLI 传入的原始配置 ====================
	FindPath       string   // 查找路径
	NamePattern    string   // 文件名模式
	PathPattern    string   // 路径模式
	ExtSlice       []string // 扩展名切片
	MaxDepth       int      // 最大深度
	SizePattern    string   // 大小模式
	ModTimePattern string   // 修改时间模式
	CaseSensitive  bool     // 区分大小写
	FullPath       bool     // 完整路径
	Hidden         bool     // 隐藏文件
	Color          bool     // 颜色输出
	Regex          bool     // 正则表达式
	ExcludeName    string   // 排除文件名
	ExcludePath    string   // 排除路径
	ExecCmd        string   // 执行命令
	Delete         bool     // 删除
	MovePath       string   // 移动路径
	PrintActions   bool     // 打印操作
	And            bool     // AND条件
	Or             bool     // OR条件
	MaxDepthLimit  int      // 最大深度限制
	Count          bool     // 统计数量
	Type           string   // 类型
	WholeWord      bool     // 完整关键字
	UseShell       bool     // 使用shell
	Quiet          bool     // 静默模式

	// ==================== 运行时生成的配置 ====================
	// 以下字段由 Init() 方法初始化，不应直接从 CLI 传入

	Cl              *colorlib.ColorLib // 颜色库实例
	NameRegex       *regexp.Regexp     // 编译后的文件名正则
	ExNameRegex     *regexp.Regexp     // 编译后的排除文件名正则
	PathRegex       *regexp.Regexp     // 编译后的路径正则
	ExPathRegex     *regexp.Regexp     // 编译后的排除路径正则
	MatchCount      *atomic.Int64      // 匹配计数原子变量
	FindExtSliceMap sync.Map           // 扩展名映射（线程安全）
}

// Init 初始化运行时字段
//
// 参数:
//   - cl: 颜色库实例
//
// 返回:
//   - error: 初始化过程中的错误
func (c *FindConfig) Init(cl *colorlib.ColorLib) error {
	c.Cl = cl

	// 初始化匹配计数器
	c.MatchCount = &atomic.Int64{}
	c.MatchCount.Store(0)

	// 编译正则表达式（仅在 Regex 模式下）
	if c.Regex {
		if err := c.compileRegexes(); err != nil {
			return err
		}
	}

	// 初始化扩展名映射
	if err := c.initExtSliceMap(); err != nil {
		return err
	}

	return nil
}

// compileRegexes 编译所有正则表达式
//
// 返回:
//   - error: 编译错误
func (c *FindConfig) compileRegexes() error {
	var err error

	if c.NameRegex, err = compileRegexPattern(c.NamePattern, c.Regex, c.WholeWord, c.CaseSensitive); err != nil {
		return fmt.Errorf("文件名正则表达式编译错误: %v", err)
	}

	if c.ExNameRegex, err = compileRegexPattern(c.ExcludeName, c.Regex, c.WholeWord, c.CaseSensitive); err != nil {
		return fmt.Errorf("排除文件名正则表达式编译错误: %v", err)
	}

	if c.PathRegex, err = compileRegexPattern(c.PathPattern, c.Regex, c.WholeWord, c.CaseSensitive); err != nil {
		return fmt.Errorf("路径正则表达式编译错误: %v", err)
	}

	if c.ExPathRegex, err = compileRegexPattern(c.ExcludePath, c.Regex, c.WholeWord, c.CaseSensitive); err != nil {
		return fmt.Errorf("排除路径正则表达式编译错误: %v", err)
	}

	return nil
}

// initExtSliceMap 初始化扩展名映射
//
// 返回:
//   - error: 处理错误
func (c *FindConfig) initExtSliceMap() error {
	for _, ext := range c.ExtSlice {
		// 确保扩展名以点开头
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		c.FindExtSliceMap.Store(ext, true)
	}
	return nil
}

// compileRegexPattern 编译正则表达式模式
//
// 参数:
//   - pattern: 正则表达式模式
//   - isRegex: 是否启用正则模式
//   - wholeWord: 是否匹配完整关键字
//   - caseSensitive: 是否区分大小写
//
// 返回:
//   - *regexp.Regexp: 编译后的正则表达式对象
//   - error: 错误信息
func compileRegexPattern(pattern string, isRegex, wholeWord, caseSensitive bool) (*regexp.Regexp, error) {
	if pattern == "" {
		return nil, nil
	}

	escapedPattern := common.RegexBuilder(pattern, isRegex, wholeWord, caseSensitive)
	return common.CompileRegex(escapedPattern)
}
```

---

## 三、需要修改的文件清单

### 3.1 新增文件

| 文件 | 说明 |
|------|------|
| `internal/commands/find/config.go` | 新的统一配置结构体 |

### 3.2 修改文件

| 文件 | 修改内容 | 预计行数 |
|------|----------|----------|
| `cmd_find.go` | 删除旧 FindConfig，使用新结构体；删除 createFindConfig 函数 | -80 行 |
| `searcher.go` | 修改参数类型 `*types.FindConfig` → `*FindConfig` | +5 行 |
| `matcher.go` | 修改参数类型 `*types.FindConfig` → `*FindConfig` | +10 行 |
| `types/types.go` | 删除 `types.FindConfig` 结构体 | -35 行 |

---

## 四、各文件修改详情

### 4.1 cmd_find.go

```go
// 删除旧的 FindConfig 结构体定义（第17-45行）

// FindCmdMain 执行查找命令
//
// 参数:
//   - cl: 颜色库实例
//   - config: 查找配置
//
// 返回值:
//   - error: 查找过程中可能发生的错误
func FindCmdMain(cl *colorlib.ColorLib, config FindConfig) error {
	findPath := config.FindPath
	if findPath == "" {
		findPath = "."
	}
	findPath = filepath.Clean(findPath)

	// 创建验证器并验证参数
	validator := NewConfigValidator(config)
	if err := validator.ValidateArgs(findPath); err != nil {
		return err
	}

	cl.SetColor(config.Color)

	// 初始化配置（替代原来的 createFindConfig）
	if err := config.Init(cl); err != nil {
		return err
	}

	matcher := NewPatternMatcher(100)

	operator := NewFileOperator(cl, config.PrintActions, config.UseShell, config.MovePath)

	// 传递 config 的指针
	searcher := NewFileSearcher(&config, matcher, operator)

	if err := searcher.Search(findPath); err != nil {
		return err
	}

	if config.Count {
		fmt.Println(config.MatchCount.Load())
	}

	return nil
}

// 删除 createFindConfig 函数（第99-175行）
// 删除 processExtensions 函数（第193-208行）
// 删除 compileRegexPattern 函数（第177-191行，移到 config.go）
```

### 4.2 searcher.go

```go
// FileSearcher 文件搜索器
type FileSearcher struct {
	config   *FindConfig // 修改：*types.FindConfig → *FindConfig
	matcher  *PatternMatcher
	operator *FileOperator
}

// NewFileSearcher 创建文件搜索器
//
// 参数:
//   - config: 查找配置
//   - matcher: 模式匹配器
//   - operator: 文件操作器
//
// 返回:
//   - *FileSearcher: 文件搜索器
func NewFileSearcher(config *FindConfig, matcher *PatternMatcher, operator *FileOperator) *FileSearcher {
	// ... 原有逻辑
}
```

### 4.3 matcher.go

```go
// MatchName 匹配文件名
//
// 参数:
//   - name: 文件名
//   - pattern: 匹配模式
//   - config: 查找配置
//
// 返回:
//   - bool: 是否匹配
func (m *PatternMatcher) MatchName(name, pattern string, config *FindConfig) bool {
	// ... 原有逻辑
}

// MatchPath 匹配路径
//
// 参数:
//   - path: 路径
//   - pattern: 匹配模式
//   - config: 查找配置
//
// 返回:
//   - bool: 是否匹配
func (m *PatternMatcher) MatchPath(path, pattern string, config *FindConfig) bool {
	// ... 原有逻辑
}

// matchPattern 统一匹配逻辑
//
// 参数:
//   - input: 输入字符串
//   - pattern: 匹配模式
//   - regex: 正则表达式
//   - config: 查找配置
//
// 返回:
//   - bool: 是否匹配
func (m *PatternMatcher) matchPattern(input, pattern string, regex *regexp.Regexp, config *FindConfig) bool {
	// ... 原有逻辑
}
```

### 4.4 types/types.go

```go
// 删除 types.FindConfig 结构体（第419-452行）
// 注意：需要确认是否有其他包使用 types.FindConfig
```

---

## 五、依赖检查

需要确认 `types.FindConfig` 是否被其他包使用：

```bash
# 在项目根目录执行
grep -r "types\.FindConfig" --include="*.go" .
```

如果只有 `find` 包使用，可以安全删除；否则需要评估影响。

---

## 六、迁移步骤

### 步骤 1：创建新文件
1. 创建 `internal/commands/find/config.go`
2. 添加新的 `FindConfig` 结构体和 `Init` 方法

### 步骤 2：修改 searcher.go 和 matcher.go
1. 修改导入语句（删除 `types` 导入）
2. 修改函数参数类型

### 步骤 3：修改 cmd_find.go
1. 删除旧的 `FindConfig` 定义
2. 删除 `createFindConfig` 函数
3. 修改 `FindCmdMain` 使用 `config.Init()`

### 步骤 4：删除 types.FindConfig
1. 在 `types/types.go` 中删除旧的结构体

### 步骤 5：验证
1. 编译项目：`go build ./...`
2. 运行测试：`go test ./internal/commands/find/...`
3. 手动测试 find 命令功能

---

## 七、优点

| 优点 | 说明 |
|------|------|
| 减少重复代码 | 消除两个结构体的字段重复 |
| 简化调用链 | 无需 `createFindConfig` 转换函数 |
| 职责清晰 | 一个结构体管理所有配置状态 |
| 易于维护 | 修改配置只需修改一处 |

## 八、风险

| 风险 | 缓解措施 |
|------|----------|
| 其他包依赖 types.FindConfig | 提前检查依赖关系 |
| 运行时字段误用 | 通过注释明确标记 |
| 初始化顺序错误 | Init() 方法必须在使用前调用 |

---

## 九、代码统计

| 项目 | 数量 |
|------|------|
| 新增文件 | 1 个 |
| 修改文件 | 4 个 |
| 预计减少代码行数 | 约 100 行 |
| 预计新增代码行数 | 约 120 行 |
| 净变化 | 约 +20 行（但结构更清晰） |
