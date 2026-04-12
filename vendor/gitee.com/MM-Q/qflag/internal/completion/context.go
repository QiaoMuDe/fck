// context.go - 上下文计算指令实现
//
// 该文件实现了 __complete context 指令，用于动态计算命令行上下文路径
// 完全替代 Shell 脚本中的静态命令树查询逻辑

package completion

import (
	"fmt"
	"strings"

	"gitee.com/MM-Q/qflag/internal/types"
)

// ContextResult 上下文计算结果
// 包含上下文路径、当前命令信息以及可用选项等
type ContextResult struct {
	// 基础信息
	Context string // 上下文路径，如 "/server/start/"
	Command string // 当前命令名
	Depth   int    // 嵌套深度

	// 当前命令信息
	CurrentCmd  string // 当前命令名称
	CurrentDesc string // 当前命令描述

	// 可用选项
	SubCommands []string // 可用子命令列表
	Flags       []string // 可用标志列表（长短名称）

	// 上下文状态
	IsFlagContext   bool // 是否处于标志上下文
	FlagsStartIndex int  // 标志开始的位置（-1 表示无）

	// 父上下文
	ParentContext string // 父上下文路径
}

// HandleContext 处理 context 指令
// 这是 __complete 子命令的入口点
//
// 参数:
//   - root: 根命令实例，通过 cmdRegistry 查询子命令
//   - args: 命令行参数（子命令名称列表，不包含程序名）
//
// 返回值:
//   - error: 处理过程中的错误
//
// 示例:
//
//	HandleContext(root, []string{"server", "start"})  // 从 __complete 调用
func HandleContext(root types.Command, args []string) error {
	// 构建完整的 tokens 列表（程序名 + 子命令路径）
	tokens := append([]string{root.Name()}, args...)

	// 自动计算 cursorPos: 使用 tokens 的长度
	cursorPos := len(tokens)

	// 计算上下文
	result := CalculateContext(root, tokens, cursorPos)

	// 输出结果（只输出上下文路径）
	fmt.Println(result.Context)

	return nil
}

// CalculateContext 计算当前上下文
// 这是核心算法，完全替代 Shell 脚本中的静态命令树查询
//
// 算法逻辑:
//  1. 从索引 1 开始遍历 tokens（跳过程序名）
//  2. 遇到以 "-" 开头的 token，标记为标志上下文并停止
//  3. 在 cmdRegistry 中查找子命令
//  4. 找到则更新上下文，继续遍历
//  5. 未找到则停止遍历，保持当前上下文
//
// 参数:
//   - root: 根命令实例
//   - tokens: 完整的命令行参数列表（包括程序名）
//   - cursorPos: 当前光标位置
//
// 返回值:
//   - *ContextResult: 上下文计算结果
func CalculateContext(root types.Command, tokens []string, cursorPos int) *ContextResult {
	result := &ContextResult{
		Context:         "/",
		Command:         root.Name(),
		Depth:           0,
		CurrentCmd:      root.Name(),
		CurrentDesc:     root.Desc(),
		SubCommands:     []string{},
		Flags:           []string{},
		IsFlagContext:   false,
		FlagsStartIndex: -1,
		ParentContext:   "",
	}

	currentCmd := root

	// 从索引 1 开始遍历 (跳过程序名 arg0)
	for i := 1; i < cursorPos && i < len(tokens); i++ {
		token := tokens[i]

		// 规则 1: 遇到标志，停止上下文构建
		// 但 -- 是标志结束符，应该跳过继续解析后面的内容
		// 已完成的等号赋值 (如 --config=value) 也不是标志, 应该跳过
		if strings.HasPrefix(token, "-") {
			if token == "--" {
				// 标志结束符，继续解析后面的位置参数
				continue
			}
			if strings.Contains(token, "=") {
				// 已完成的等号赋值，跳过继续解析
				continue
			}
			result.IsFlagContext = true
			result.FlagsStartIndex = i
			break
		}

		// 规则 2: 在注册表中查找子命令
		subCmd, found := currentCmd.GetSubCmd(token)
		if !found {
			// 不是有效的子命令，停止遍历
			break
		}

		// 规则 3: 更新上下文
		result.ParentContext = result.Context
		result.Context += token + "/"
		result.Depth++
		result.CurrentCmd = token
		result.CurrentDesc = subCmd.Desc()
		currentCmd = subCmd
	}

	// 获取当前命令的可用选项
	result.SubCommands = getSubCommandNames(currentCmd)
	result.Flags = getFlagNames(currentCmd)

	return result
}

// getSubCommandNames 获取子命令名称列表
//
// 参数:
//   - cmd: 命令实例
//
// 返回值:
//   - []string: 子命令名称列表（已自动过滤隐藏命令）
func getSubCommandNames(cmd types.Command) []string {
	subCmds := cmd.SubCmds()
	names := make([]string, 0, len(subCmds))
	for _, subCmd := range subCmds {
		names = append(names, subCmd.Name())
	}
	return names
}

// getFlagNames 获取标志名称列表（包括长短名称）
//
// 参数:
//   - cmd: 命令实例
//
// 返回值:
//   - []string: 标志名称列表（长名称带 -- 前缀，短名称带 - 前缀）
func getFlagNames(cmd types.Command) []string {
	flags := cmd.Flags()
	names := make([]string, 0, len(flags)*2)

	for _, flag := range flags {
		// 添加长名称（带 -- 前缀）
		if flag.LongName() != "" {
			names = append(names, "--"+flag.LongName())
		}

		// 添加短名称（带 - 前缀，如果有）
		if flag.ShortName() != "" {
			names = append(names, "-"+flag.ShortName())
		}
	}

	return names
}

// findCommandByContext 根据上下文路径查找命令
//
// 参数:
//   - root: 根命令
//   - context: 上下文路径，如 "/server/start/"
//
// 返回值:
//   - types.Command: 找到的命令，如果未找到则返回 nil
func findCommandByContext(root types.Command, context string) types.Command {
	if context == "/" {
		return root
	}

	parts := strings.Split(strings.Trim(context, "/"), "/")
	current := root

	for _, part := range parts {
		subCmd, found := current.GetSubCmd(part)
		if !found {
			return nil
		}
		current = subCmd
	}

	return current
}
