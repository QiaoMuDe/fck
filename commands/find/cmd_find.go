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
	"gitee.com/MM-Q/fck/commands/internal/common"
	"gitee.com/MM-Q/fck/commands/internal/types"
)

// FindCmdMain 是 find 子命令的主函数 - 重构后作为协调器
func FindCmdMain(cl *colorlib.ColorLib) error {
	// 获取第一个参数作为查找路径
	findPath := findCmd.Arg(0)
	if findPath == "" {
		findPath = "."
	}
	findPath = filepath.Clean(findPath)

	// 创建验证器并验证参数
	validator := NewConfigValidator()
	if err := validator.ValidateArgs(findPath); err != nil {
		return err
	}

	// 设置颜色
	cl.SetColor(findCmdColor.Get())

	// 创建配置
	config, err := createFindConfig(cl)
	if err != nil {
		return err
	}

	// 创建模式匹配器（使用缓存）
	matcher := NewPatternMatcher(100) // 最多缓存100个正则表达式

	// 创建文件操作器
	operator := NewFileOperator(cl)

	// 创建搜索器
	searcher := NewFileSearcher(config, matcher, operator)

	// 单线程搜索
	if err := searcher.Search(findPath); err != nil {
		return err
	}

	// 如果启用了count标志, 只输出匹配数量
	if findCmdCount.Get() {
		fmt.Println(config.MatchCount.Load())
	}

	return nil
}

// createFindConfig 创建查找配置
//
// 参数:
//   - cl: 颜色库
//
// 返回:
//   - *types.FindConfig: 查找配置
//   - error: 错误信息
func createFindConfig(cl *colorlib.ColorLib) (*types.FindConfig, error) {
	// 准备正则表达式模式
	isRegex := findCmdRegex.Get()       // 是否启用正则模式
	wholeWord := findCmdWholeWord.Get() // 是否匹配完整关键字
	caseSensitive := findCmdCase.Get()  // 是否区分大小写

	// 定义正则表达式
	var nameRegex, exNameRegex, pathRegex, exPathRegex *regexp.Regexp
	var err error

	// 如果启用正则模式，编译正则表达式
	if isRegex {
		if nameRegex, err = compileRegexPattern(findCmdName.Get(), isRegex, wholeWord, caseSensitive); err != nil {
			return nil, fmt.Errorf("文件名正则表达式编译错误: %v", err)
		}
		if exNameRegex, err = compileRegexPattern(findCmdExcludeName.Get(), isRegex, wholeWord, caseSensitive); err != nil {
			return nil, fmt.Errorf("排除文件名正则表达式编译错误: %v", err)
		}
		if pathRegex, err = compileRegexPattern(findCmdPath.Get(), isRegex, wholeWord, caseSensitive); err != nil {
			return nil, fmt.Errorf("路径正则表达式编译错误: %v", err)
		}
		if exPathRegex, err = compileRegexPattern(findCmdExcludePath.Get(), isRegex, wholeWord, caseSensitive); err != nil {
			return nil, fmt.Errorf("排除路径正则表达式编译错误: %v", err)
		}
	}

	// 创建匹配计数器
	matchCount := atomic.Int64{}
	matchCount.Store(0)

	// 创建配置实例
	config := &types.FindConfig{
		Cl:            cl,                       // 颜色库
		NameRegex:     nameRegex,                // 文件名正则
		ExNameRegex:   exNameRegex,              // 排除文件名正则
		PathRegex:     pathRegex,                // 路径正则
		ExPathRegex:   exPathRegex,              // 排除路径正则
		IsRegex:       isRegex,                  // 是否启用正则模式
		WholeWord:     wholeWord,                // 是否匹配完整关键字
		CaseSensitive: caseSensitive,            // 是否区分大小写
		MatchCount:    &matchCount,              // 匹配计数器
		NamePattern:   findCmdName.Get(),        // 文件名模式
		PathPattern:   findCmdPath.Get(),        // 路径模式
		ExNamePattern: findCmdExcludeName.Get(), // 排除文件名模式
		ExPathPattern: findCmdExcludePath.Get(), // 排除路径模式
	}

	// 处理扩展名参数
	if err := processExtensions(config); err != nil {
		return nil, err
	}

	return config, nil
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

	// 构建正则表达式
	escapedPattern := common.RegexBuilder(pattern, isRegex, wholeWord, caseSensitive)
	return common.CompileRegex(escapedPattern)
}

// processExtensions 处理扩展名参数
//
// 参数:
//   - config: 配置参数
//
// 返回:
//   - error: 错误信息
func processExtensions(config *types.FindConfig) error {
	// 如果扩展名切片不为空，则处理扩展名参数
	if findCmdExt.Len() > 0 {
		for _, ext := range findCmdExt.Get() {
			// 如果扩展名不包含"."则添加"."
			if !strings.HasPrefix(ext, ".") {
				ext = "." + ext
			}
			config.FindExtSliceMap.Store(ext, true)
		}
	}
	return nil
}
