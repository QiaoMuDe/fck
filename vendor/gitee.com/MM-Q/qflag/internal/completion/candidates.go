// candidates.go - 候选选项获取指令实现
//
// 该文件实现了 __complete candidates 指令，用于根据上下文路径
// 获取该上下文下所有可用的补全选项（子命令和标志）

package completion

import (
	"fmt"
	"strings"

	"gitee.com/MM-Q/qflag/internal/types"
)

// HandleCandidates 处理 candidates 指令
//
// 参数:
//   - root: 根命令实例
//   - args: [context]
//
// 返回值:
//   - error: 处理错误
//
// 示例:
//
//	HandleCandidates(root, []string{"/server/start/"})
func HandleCandidates(root types.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: __complete candidates <context>")
	}

	context := args[0]

	// 根据上下文查找命令
	cmd := findCommandByContext(root, context)
	if cmd == nil {
		// 无效的上下文，返回空
		return nil
	}

	// 收集所有候选选项
	var candidates []string

	// 添加子命令
	candidates = append(candidates, getSubCommandNames(cmd)...)

	// 添加标志
	candidates = append(candidates, getFlagNames(cmd)...)

	// 添加内置标志
	candidates = append(candidates, getBuiltinFlagNames(cmd, context)...)

	// 输出（空格分隔）
	fmt.Println(strings.Join(candidates, " "))

	return nil
}

// GetCandidates 获取候选选项（供程序内部使用）
//
// 参数:
//   - root: 根命令实例
//   - context: 上下文路径
//
// 返回值:
//   - []string: 候选选项列表
//   - error: 处理错误
func GetCandidates(root types.Command, context string) ([]string, error) {
	// 根据上下文查找命令
	cmd := findCommandByContext(root, context)
	if cmd == nil {
		// 无效的上下文，返回空列表
		return []string{}, nil
	}

	// 收集所有候选选项
	var candidates []string

	// 添加子命令
	candidates = append(candidates, getSubCommandNames(cmd)...)

	// 添加标志
	candidates = append(candidates, getFlagNames(cmd)...)

	// 添加内置标志
	candidates = append(candidates, getBuiltinFlagNames(cmd, context)...)

	return candidates, nil
}

// getBuiltinFlagNames 获取应该注册的内置标志名称列表
//
// 规则:
//   - 所有命令（包括根命令和子命令）默认都添加帮助标志
//   - 只有根命令根据配置添加版本标志和补全标志
//
// 参数:
//   - cmd: 命令实例
//   - context: 上下文路径，用于判断是否是根命令
//
// 返回值:
//   - []string: 内置标志名称列表（长名称和短名称）
func getBuiltinFlagNames(cmd types.Command, context string) []string {
	var names []string

	// 帮助标志 - 所有命令都添加
	names = append(names, "--"+types.HelpFlagName)
	names = append(names, "-"+types.HelpFlagShortName)

	// 根命令特殊处理：根据配置添加版本和补全标志
	if context == "/" {
		config := cmd.Config()

		// 版本标志 - 根据是否有版本信息决定
		if config.Version != "" {
			names = append(names, "--"+types.VersionFlagName)
			names = append(names, "-"+types.VersionFlagShortName)
		}

		// 补全标志 - 根据配置决定
		if config.DynamicCompletion {
			names = append(names, "--"+types.CompletionFlagName)
		}
	}

	return names
}
