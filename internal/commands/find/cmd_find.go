// Package find 实现了文件查找命令的主要逻辑和配置管理。
// 该文件包含 find 子命令的入口函数，负责参数验证、配置创建和搜索执行。
package find

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"sync/atomic"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/internal/types"
	common "gitee.com/MM-Q/fck/internal/utils"
)

// FindConfig 查找配置结构体
type FindConfig struct {
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
}

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

	findConfig, err := createFindConfig(cl, config)
	if err != nil {
		return err
	}

	matcher := NewPatternMatcher(100)

	operator := NewFileOperator(cl, config.PrintActions, config.UseShell, config.MovePath)

	searcher := NewFileSearcher(findConfig, matcher, operator)

	if err := searcher.Search(findPath); err != nil {
		return err
	}

	if config.Count {
		fmt.Println(findConfig.MatchCount.Load())
	}

	return nil
}

// createFindConfig 创建查找配置
//
// 参数:
//   - cl: 颜色库
//   - config: 查找配置
//
// 返回:
//   - *types.FindConfig: 查找配置
//   - error: 错误信息
func createFindConfig(cl *colorlib.ColorLib, config FindConfig) (*types.FindConfig, error) {
	isRegex := config.Regex
	wholeWord := config.WholeWord
	caseSensitive := config.CaseSensitive

	var nameRegex, exNameRegex, pathRegex, exPathRegex *regexp.Regexp
	var err error

	if isRegex {
		if nameRegex, err = compileRegexPattern(config.NamePattern, isRegex, wholeWord, caseSensitive); err != nil {
			return nil, fmt.Errorf("文件名正则表达式编译错误: %v", err)
		}
		if exNameRegex, err = compileRegexPattern(config.ExcludeName, isRegex, wholeWord, caseSensitive); err != nil {
			return nil, fmt.Errorf("排除文件名正则表达式编译错误: %v", err)
		}
		if pathRegex, err = compileRegexPattern(config.PathPattern, isRegex, wholeWord, caseSensitive); err != nil {
			return nil, fmt.Errorf("路径正则表达式编译错误: %v", err)
		}
		if exPathRegex, err = compileRegexPattern(config.ExcludePath, isRegex, wholeWord, caseSensitive); err != nil {
			return nil, fmt.Errorf("排除路径正则表达式编译错误: %v", err)
		}
	}

	matchCount := atomic.Int64{}
	matchCount.Store(0)

	findConfig := &types.FindConfig{
		Cl:             cl,
		NameRegex:      nameRegex,
		ExNameRegex:    exNameRegex,
		PathRegex:      pathRegex,
		ExPathRegex:    exPathRegex,
		IsRegex:        isRegex,
		WholeWord:      wholeWord,
		CaseSensitive:  caseSensitive,
		MatchCount:     &matchCount,
		NamePattern:    config.NamePattern,
		PathPattern:    config.PathPattern,
		ExNamePattern:  config.ExcludeName,
		ExPathPattern:  config.ExcludePath,
		Quiet:          config.Quiet,
		MaxDepth:       config.MaxDepth,
		Hidden:         config.Hidden,
		SizePattern:    config.SizePattern,
		ModTimePattern: config.ModTimePattern,
		ExtSlice:       config.ExtSlice,
		Type:           config.Type,
		Count:          config.Count,
		FullPath:       config.FullPath,
		Delete:         config.Delete,
		MovePath:       config.MovePath,
		ExecCmd:        config.ExecCmd,
		And:            config.And,
		Or:             config.Or,
		MaxDepthLimit:  config.MaxDepthLimit,
	}

	if err := processExtensions(findConfig, config.ExtSlice); err != nil {
		return nil, err
	}

	return findConfig, nil
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

// processExtensions 处理扩展名参数
//
// 参数:
//   - findConfig: 查找配置
//   - extSlice: 扩展名切片
//
// 返回:
//   - error: 错误信息
func processExtensions(findConfig *types.FindConfig, extSlice []string) error {
	for _, ext := range extSlice {
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		findConfig.FindExtSliceMap.Store(ext, true)
	}
	return nil
}
