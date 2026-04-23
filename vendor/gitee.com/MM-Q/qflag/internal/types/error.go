// error.go - 错误类型定义
//
// 该文件包含 qflag 使用的自定义错误类型，支持智能纠错提示功能

package types

import (
	"fmt"
	"strings"
)

// UnknownSubcommandError 未知子命令错误
//
// 当用户输入的子命令不存在时返回此错误，包含相似子命令建议
type UnknownSubcommandError struct {
	Command     string   // 当前命令名
	Input       string   // 用户输入的错误子命令
	Suggestions []string // 相似子命令建议列表
}

// Error 实现 error 接口，返回格式化的错误信息
//
// 格式示例：
//
//	myapp: 'cnfig' is not a valid command. See 'myapp --help'.
//
//	The most similar commands are
//	        config
func (e *UnknownSubcommandError) Error() string {
	var sb strings.Builder
	_, _ = fmt.Fprintf(&sb, "%s: '%s' is not a valid command. See '%s --help'.\n",
		e.Command, e.Input, e.Command)

	if len(e.Suggestions) > 0 {
		sb.WriteString("\nThe most similar commands are\n")
		for _, sug := range e.Suggestions {
			_, _ = fmt.Fprintf(&sb, "\t%s\n", sug)
		}
	}

	return sb.String()
}

// UnknownFlagError 未知标志错误
//
// 当用户输入的标志不存在时返回此错误，包含相似标志建议
type UnknownFlagError struct {
	Command     string   // 当前命令名
	Input       string   // 用户输入的错误标志
	Suggestions []string // 相似标志建议列表
}

// Error 实现 error 接口，返回格式化的错误信息
//
// 格式示例：
//
//	myapp: unknown flag: '--verboose'
//
//	The most similar flags are
//	        --verbose
//	        -v
func (e *UnknownFlagError) Error() string {
	var sb strings.Builder
	_, _ = fmt.Fprintf(&sb, "%s: unknown flag: '%s'\n",
		e.Command, e.Input)

	if len(e.Suggestions) > 0 {
		sb.WriteString("\nThe most similar flags are\n")
		for _, sug := range e.Suggestions {
			_, _ = fmt.Fprintf(&sb, "\t%s\n", sug)
		}
	}

	return sb.String()
}
