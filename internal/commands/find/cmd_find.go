// Package find 实现了文件查找命令的主要逻辑和配置管理。
// 该文件包含 find 子命令的入口函数，负责参数验证、配置创建和搜索执行。
package find

import (
	"fmt"

	"gitee.com/MM-Q/color"
)

// FindCmdMain 执行查找命令
//
// 参数:
//   - cl: 颜色库实例
//   - config: 查找配置指针
//
// 返回值:
//   - error: 查找过程中可能发生的错误
func FindCmdMain(cl *color.GlobalColor, config *FindConfig) error {
	// 创建验证器
	validator := NewConfigValidator(config)

	// 设置颜色
	cl.SetNoColor(!config.Color)

	// 初始化配置（包含路径处理、正则编译、扩展名映射等）
	if err := config.Init(cl); err != nil {
		return err
	}

	// 验证参数
	if err := validator.ValidateArgs(config.FindPath); err != nil {
		return err
	}

	matcher := NewPatternMatcher(100)

	operator := NewFileOperator(cl, config.PrintActions, config.MovePath)

	// 传递 config 指针
	searcher := NewFileSearcher(config, matcher, operator)

	if err := searcher.Search(config.FindPath); err != nil {
		return err
	}

	if config.Count {
		fmt.Println(config.MatchCount.Load())
	}

	return nil
}
