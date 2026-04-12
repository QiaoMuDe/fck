// suggestion.go - 智能纠错查找器
//
// 该文件实现智能纠错功能，用于在子命令或标志输入错误时提供相似建议

package parser

import (
	"strings"

	"gitee.com/MM-Q/go-kit/fuzzy"
	"gitee.com/MM-Q/qflag/internal/types"
)

// SuggestionFinder 智能纠错查找器
//
// 封装建议查找逻辑，保持解析方法简洁
type SuggestionFinder struct {
	maxSuggestions int // 最大建议数量
}

// NewSuggestionFinder 创建查找器
//
// 参数:
//   - maxSuggestions: 最大建议数量
//
// 返回值:
//   - *SuggestionFinder: 查找器实例
func NewSuggestionFinder(maxSuggestions int) *SuggestionFinder {
	return &SuggestionFinder{maxSuggestions: maxSuggestions}
}

// FindForSubcommand 查找子命令建议
//
// 根据输入字符串，在命令的所有子命令中查找相似的命令名
//
// 参数:
//   - input: 用户输入的错误子命令
//   - cmd: 当前命令
//
// 返回值:
//   - []string: 相似子命令列表
func (f *SuggestionFinder) FindForSubcommand(input string, cmd types.Command) []string {
	subCmds := cmd.SubCmds()
	if len(subCmds) == 0 {
		return nil
	}

	names := make([]string, 0, len(subCmds)*2)
	for _, sc := range subCmds {
		if sc.LongName() != "" {
			names = append(names, sc.LongName())
		}
		if sc.ShortName() != "" {
			names = append(names, sc.ShortName())
		}
	}

	return f.findSimilar(input, names)
}

// FindForFlag 查找标志建议
//
// 根据输入字符串，在命令的所有标志中查找相似的标志名
//
// 参数:
//   - input: 用户输入的错误标志
//   - cmd: 当前命令
//
// 返回值:
//   - []string: 相似标志列表
func (f *SuggestionFinder) FindForFlag(input string, cmd types.Command) []string {
	flags := cmd.FlagRegistry().List()
	if len(flags) == 0 {
		return nil
	}

	names := make([]string, 0, len(flags)*2)
	for _, fl := range flags {
		if fl.LongName() != "" {
			names = append(names, "--"+fl.LongName())
		}
		if fl.ShortName() != "" {
			names = append(names, "-"+fl.ShortName())
		}
	}

	cleanInput := strings.TrimLeft(input, "-")
	return f.findSimilar(cleanInput, names)
}

// findSimilar 内部模糊匹配
//
// 参数:
//   - input: 输入字符串
//   - candidates: 候选列表
//
// 返回值:
//   - []string: 匹配结果列表
func (f *SuggestionFinder) findSimilar(input string, candidates []string) []string {
	// 使用前缀匹配, 其次模糊匹配
	matches := fuzzy.Complete(input, candidates)

	// 限制返回结果数量
	if len(matches) > f.maxSuggestions {
		matches = matches[:f.maxSuggestions]
	}
	result := make([]string, len(matches))
	for i, m := range matches {
		result[i] = m.Str
	}
	return result
}

// newUnknownSubcommandError 创建未知子命令错误（带建议）
//
// 参数:
//   - cmd: 当前命令
//   - input: 用户输入的错误子命令
//
// 返回值:
//   - error: 如果找到建议返回错误，否则返回 nil (不拦截)
func newUnknownSubcommandError(cmd types.Command, input string) error {
	// 创建模糊匹配查找器, 最多返回3个建议
	finder := NewSuggestionFinder(3)
	suggestions := finder.FindForSubcommand(input, cmd)

	// 没有找到建议，不拦截，返回 nil
	if len(suggestions) == 0 {
		return nil
	}

	return &types.UnknownSubcommandError{
		Command:     cmd.Name(),
		Input:       input,
		Suggestions: suggestions,
	}
}

// newUnknownFlagError 创建未知标志错误（带建议）
//
// 参数:
//   - cmd: 当前命令
//   - input: 用户输入的错误标志
//
// 返回值:
//   - error: 带建议的错误
func newUnknownFlagError(cmd types.Command, input string) error {
	finder := NewSuggestionFinder(3)
	suggestions := finder.FindForFlag(input, cmd)

	return &types.UnknownFlagError{
		Command:     cmd.Name(),
		Input:       input,
		Suggestions: suggestions,
	}
}

// checkUnknownFlags 预扫描参数，检查未知标志
//
// 在调用标准库 flag.FlagSet.Parse() 之前，先扫描参数列表，
// 提前发现未知标志并返回带建议的错误。
//
// 判断逻辑:
//   - 以 - 或 -- 开头的参数一定是标志
//   - 如果不在已注册标志列表中，就是错误的标志
//   - 遇到 -- 停止扫描，后面的都视为位置参数
//   - 遇到子命令名时停止扫描（后续标志由子命令处理）
//
// 参数:
//   - cmd: 当前命令
//   - args: 命令行参数列表
//
// 返回值:
//   - error: 如果发现未知标志返回错误，否则返回 nil
func checkUnknownFlags(cmd types.Command, args []string) error {
	// 获取所有已注册的标志名（长短名称都包括）
	registeredFlags := make(map[string]bool)
	for _, f := range cmd.FlagRegistry().List() {
		if f.LongName() != "" {
			registeredFlags["--"+f.LongName()] = true
			registeredFlags["-"+f.LongName()] = true // 支持单横杠长名称
		}
		if f.ShortName() != "" {
			registeredFlags["-"+f.ShortName()] = true
		}
	}

	// 扫描参数
	for i := 0; i < len(args); i++ {
		arg := args[i]

		// 遇到 -- 停止扫描，后面的都视为位置参数
		if arg == "--" {
			break
		}

		// 不是标志格式，检查是否为子命令
		if !strings.HasPrefix(arg, "-") {
			// 如果是子命令名，停止扫描（后续标志由子命令处理）
			if _, isSubCmd := cmd.CmdRegistry().Get(arg); isSubCmd {
				break
			}
			continue
		}

		// 处理 --flag=value 格式
		flagName := arg
		if idx := strings.Index(arg, "="); idx != -1 {
			flagName = arg[:idx]
		}

		// 不是已注册的标志 → 纠错
		if !registeredFlags[flagName] {
			return newUnknownFlagError(cmd, flagName)
		}
	}

	return nil
}
