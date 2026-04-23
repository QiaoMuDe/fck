// Package completion 自动补全内部实现
// 本文件包含动态补全逻辑，用于处理 __complete 子命令
package completion

import (
	"fmt"
	"strings"

	"gitee.com/MM-Q/go-kit/fuzzy"
	"gitee.com/MM-Q/qflag/internal/types"
)

// HandleDynamicComplete 处理 __complete 子命令的路由
//
// 参数:
//   - root: 根命令实例，用于查询命令树
//   - instruction: 指令名称
//   - params: 指令参数列表
//
// 返回值:
//   - error: 处理错误
func HandleDynamicComplete(root types.Command, instruction string, params []string) error {
	switch instruction {
	case types.InstructionFuzzy:
		return handleFuzzy(params)

	case types.InstructionContext:
		return HandleContext(root, params)

	case types.InstructionCandidates:
		return HandleCandidates(root, params)

	case types.InstructionEnum:
		return HandleEnum(root, params)

	case types.InstructionAll:
		return handleAll(root, params)

	default:
		return fmt.Errorf("unknown instruction: %s", instruction)
	}
}

// handleFuzzy 处理 fuzzy 指令
//
// 参数:
//   - args: 参数列表，第一个是模式，后面是候选列表
//
// 返回值:
//   - error: 处理错误
//
// 输出格式: 每行一个匹配结果（按匹配质量降序）
func handleFuzzy(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: __complete fuzzy <pattern> <candidates...>")
	}

	pattern := args[0]     // 模式参数
	candidates := args[1:] // 候选参数

	// 使用 fuzzy 包执行前缀补全
	matches := fuzzy.CompletePrefix(pattern, candidates)

	// 输出匹配结果（只输出匹配的字符串）
	for _, match := range matches {
		fmt.Println(match.Str)
	}

	return nil
}

// handleAll 处理 all 指令，一次性返回所有补全信息
//
// 参数:
//   - root: 根命令实例
//   - args: [cur, prev, cmd_args...]
//     args[0]: cur - 当前输入的词
//     args[1]: prev - 前一个词
//     args[2:]: cmd_args - 已输入的子命令参数
//
// 返回值:
//   - error: 处理错误
//
// 输出格式:
//
//	CONTEXT:<上下文路径>
//	CUR:<当前输入>
//	PREV:<前一个词>
//	CANDIDATES:<候选项列表>
//	ENUM:<枚举值列表>
//	MATCHES:<匹配结果>
//	IS_FLAG:<true|false>
func handleAll(root types.Command, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: __complete all <cur> <prev> [cmd_args...]")
	}

	// 解析参数
	// 去除可能的引号 (PowerShell 等 shell 会保留引号在输入中)
	cur := strings.Trim(args[0], `"'`)  // 当前输入
	prev := strings.Trim(args[1], `"'`) // 前一个词
	cmdArgs := []string{}
	if len(args) > 2 {
		cmdArgs = args[2:] // 子命令参数
	}

	// 1. 计算上下文
	// CalculateContext 期望 tokens 包含程序名作为第一个元素 (从索引1开始遍历)
	// 当 cmdArgs 为空时 (如刚输入命令名后按 Tab) ,tokens = [""] 表示只有程序名
	// 当 cmdArgs = ["config"] 时, tokens = ["", "config"], CalculateContext 能正确识别子命令
	// 所以我们需要在 cmdArgs 前面添加一个空字符串作为占位符
	tokens := append([]string{""}, cmdArgs...)
	contextResult := CalculateContext(root, tokens, len(tokens))
	context := "/"
	if contextResult != nil {
		context = contextResult.Context
	}

	// 2. 获取候选项
	candidates, _ := GetCandidates(root, context)

	// 3. 执行补全逻辑
	var matchStrings []string
	var enumValues []string

	// 3. 判断补全类型并执行相应逻辑
	isFlagValueCompletion := isFlagValueContext(cur, prev)

	if isFlagValueCompletion {
		// ========== 标志值补全 ==========
		flagType, found := getFlagType(root, context, prev)

		if !found {
			// 标志不存在，按普通候选项补全
			matchStrings = fuzzyMatch(candidates, cur)
		} else {
			switch flagType {
			case types.FlagTypeBool:
				// 布尔标志：不需要值，补全其他标志/子命令
				matchStrings = fuzzyMatch(candidates, cur)

			case types.FlagTypeEnum:
				// 枚举标志：获取枚举值并模糊匹配
				enumValues, _ = GetEnumValues(root, context, prev)
				matchStrings = fuzzyMatch(enumValues, cur)

			default:
				// 其他类型（String/Int/Duration/Size等）：需要值
				// matchStrings 保持为空，由 Shell 回退到路径补全
			}
		}
	} else {
		// ========== 普通候选项补全 ==========
		matchStrings = fuzzyMatch(candidates, cur)
	}

	// 5. 输出结果（带前缀的多行格式）
	fmt.Printf("CONTEXT:%s\n", context)
	fmt.Printf("CUR:%s\n", cur)
	fmt.Printf("PREV:%s\n", prev)
	fmt.Printf("CANDIDATES:%s\n", strings.Join(candidates, " "))
	fmt.Printf("ENUM:%s\n", strings.Join(enumValues, " "))
	fmt.Printf("MATCHES:%s\n", strings.Join(matchStrings, " "))
	fmt.Printf("IS_FLAG:%v\n", isFlagValueCompletion && len(enumValues) > 0)

	return nil
}

// getFlagType 获取指定上下文中标志的类型
//
// 参数:
//   - root: 根命令实例
//   - context: 上下文路径
//   - flagName: 标志名称
//
// 返回值:
//   - FlagType: 标志类型
//   - bool: 是否找到标志
//
// 说明:
//   - 内置标志（如 --help, --version）在解析时动态注册
//   - 这里通过名称匹配来识别内置标志的类型
func getFlagType(root types.Command, context string, flagName string) (types.FlagType, bool) {
	cmd := findCommandByContext(root, context)
	if cmd == nil {
		return types.FlagTypeUnknown, false
	}

	// 首先尝试从命令的 Flags() 中查找
	flag := findFlagByName(cmd, flagName)
	if flag != nil {
		return flag.Type(), true
	}

	// 如果没有找到，检查是否是内置标志
	// 内置标志在解析时动态注册，但我们可以根据名称识别其类型
	return getBuiltinFlagType(flagName, context, cmd)
}

// getBuiltinFlagType 根据标志名称识别内置标志的类型
//
// 参数:
//   - flagName: 标志名称
//   - context: 上下文路径
//   - cmd: 命令实例
//
// 返回值:
//   - FlagType: 标志类型
//   - bool: 是否为内置标志
func getBuiltinFlagType(flagName string, context string, cmd types.Command) (types.FlagType, bool) {
	// 移除 "-" 或 "--" 前缀
	name := strings.TrimPrefix(flagName, "--")
	if name == flagName {
		name = strings.TrimPrefix(flagName, "-")
	}

	// 帮助标志：所有命令都有
	if name == types.HelpFlagName || name == types.HelpFlagShortName {
		return types.FlagTypeBool, true
	}

	// 根命令特有的内置标志
	if context == "/" {
		config := cmd.Config()

		// 版本标志
		if config.Version != "" {
			if name == types.VersionFlagName || name == types.VersionFlagShortName {
				return types.FlagTypeBool, true
			}
		}

		// 补全标志
		if config.Completion {
			if name == types.CompletionFlagName {
				return types.FlagTypeEnum, true
			}
		}
	}

	return types.FlagTypeUnknown, false
}

// fuzzyMatch 对候选列表进行模糊匹配
//
// 参数:
//   - candidates: 候选列表
//   - cur: 当前输入
//
// 返回值:
//   - []string: 匹配结果
func fuzzyMatch(candidates []string, cur string) []string {
	// 如果当前输入为空，返回所有候选项
	if cur == "" {
		return candidates
	}

	// 执行模糊匹配
	matches := fuzzy.CompletePrefix(cur, candidates)
	result := make([]string, len(matches))
	for i, match := range matches {
		result[i] = match.Str
	}
	return result
}

// isFlagValueContext 判断当前是否处于标志值补全上下文
//
// 标志值补全 vs 普通候选项补全的区别：
//
//  1. 标志值补全场景（返回 true）：
//     用户已输入标志名，正在输入该标志的值
//     例如：--format <cur> 或 -f <cur>
//     此时 cur 是标志的值，不是标志名
//
//  2. 普通候选项补全场景（返回 false）：
//     用户正在输入标志名、子命令名或其他候选项
//     例如：<cur> 或 --<cur> 或 -<cur>
//     此时 cur 可能是标志名、子命令名等
//
// 参数:
//   - cur: 当前输入的词（用户正在输入的内容）
//   - prev: 前一个词（可能是标志名）
//
// 返回值:
//   - bool: true 表示处于标志值补全上下文，false 表示普通候选项补全
//
// 判断逻辑:
//  1. cur 不以 "-" 开头 → 用户在输入值，不是标志名
//  2. prev 以 "-" 开头 → 前一个词是标志
//  3. prev 不是 "--" → 排除单独的 "--" 分隔符
//  4. prev 不包含 "=" → 排除 --flag=value 形式（值已经在 prev 中）
//
// 示例:
//
//	isFlagValueContext("json", "--format") → true  （--format 的值）
//	isFlagValueContext("--ty", "list")     → false （正在输入标志名 --type）
//	isFlagValueContext("-t", "list")       → false （正在输入标志名 -t）
//	isFlagValueContext("", "--format")     → true  （--format 的值，空输入）
func isFlagValueContext(cur, prev string) bool {
	// 快速路径：如果 cur 以 - 开头，说明用户在补全标志名，不是标志值
	if strings.HasPrefix(cur, "-") {
		return false
	}

	// 检查 prev 是否是需要值的标志
	// 条件：以 - 开头、不是 --、不包含 =（避免 --flag=value 形式）
	return strings.HasPrefix(prev, "-") && prev != "--" && !strings.Contains(prev, "=")
}
