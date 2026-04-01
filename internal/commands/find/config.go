// Package find 实现了文件查找命令的配置管理。
// 该文件包含统一的查找配置结构体，整合了 CLI 参数和运行时状态。
package find

import (
	"fmt"
	"path/filepath"
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
	MaxDepthLimit  int      // 最大深度限制
	Count          bool     // 统计数量
	Type           string   // 类型
	WholeWord      bool     // 完整关键字
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

	// 处理查找路径：默认为当前目录，并清理路径格式
	if c.FindPath == "" {
		c.FindPath = "."
	}
	c.FindPath = filepath.Clean(c.FindPath)

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
		return fmt.Errorf("failed to compile filename regex: %v", err)
	}

	if c.ExNameRegex, err = compileRegexPattern(c.ExcludeName, c.Regex, c.WholeWord, c.CaseSensitive); err != nil {
		return fmt.Errorf("failed to compile exclude filename regex: %v", err)
	}

	if c.PathRegex, err = compileRegexPattern(c.PathPattern, c.Regex, c.WholeWord, c.CaseSensitive); err != nil {
		return fmt.Errorf("failed to compile path regex: %v", err)
	}

	if c.ExPathRegex, err = compileRegexPattern(c.ExcludePath, c.Regex, c.WholeWord, c.CaseSensitive); err != nil {
		return fmt.Errorf("failed to compile exclude path regex: %v", err)
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
