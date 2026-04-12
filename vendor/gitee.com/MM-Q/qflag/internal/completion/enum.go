// enum.go - 枚举值补全指令实现
//
// 该文件实现了 __complete enum 指令, 用于根据上下文路径和标志名
// 获取枚举类型标志的所有可选值

package completion

import (
	"fmt"
	"strings"

	"gitee.com/MM-Q/qflag/internal/types"
)

// HandleEnum 处理 enum 指令
//
// 参数:
//   - root: 根命令实例
//   - args: [context, flag-name]
//
// 返回值:
//   - error: 处理错误
//
// 示例:
//
//	HandleEnum(root, []string{"/server/start/", "--log-level"})
//	HandleEnum(root, []string{"/server/start/", "-o"})
func HandleEnum(root types.Command, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: __complete enum <context> <flag-name>")
	}

	context := args[0]  // 获取上下文路径
	flagName := args[1] // 获取标志名称

	// 根据上下文查找命令
	cmd := findCommandByContext(root, context)
	if cmd == nil {
		// 无效的上下文, 返回空
		return nil
	}

	// 查找标志
	targetFlag := findFlagByName(cmd, flagName)
	if targetFlag == nil {
		// 标志不存在, 返回空
		return nil
	}

	// 获取枚举值
	enumValues := getEnumValues(targetFlag)
	if len(enumValues) == 0 {
		// 不是枚举类型或没有枚举值, 返回空
		return nil
	}

	// 输出 (空格分隔)
	fmt.Println(strings.Join(enumValues, " "))

	return nil
}

// GetEnumValues 获取枚举值 (供程序内部使用)
//
// 参数:
//   - root: 根命令实例
//   - context: 上下文路径
//   - flagName: 标志名称
//
// 返回值:
//   - []string: 枚举值列表
//   - error: 处理错误
func GetEnumValues(root types.Command, context string, flagName string) ([]string, error) {
	// 根据上下文查找命令
	cmd := findCommandByContext(root, context)
	if cmd == nil {
		// 无效的上下文, 返回空列表
		return []string{}, nil
	}

	// 查找标志
	targetFlag := findFlagByName(cmd, flagName)
	if targetFlag == nil {
		// 标志不存在, 返回空列表
		return []string{}, nil
	}

	// 获取枚举值
	return getEnumValues(targetFlag), nil
}

// findFlagByName 根据名称查找标志
//
// 参数:
//   - cmd: 命令实例
//   - flagName: 标志名称 (支持长名称带 "--"、短名称带 "-", 以及带 "=" 后缀的形式)
//
// 返回值:
//   - types.Flag: 找到的标志, 如果未找到则返回 nil
func findFlagByName(cmd types.Command, flagName string) types.Flag {
	// 移除可能的 "=" 后缀
	flagName = strings.TrimSuffix(flagName, "=")

	for _, flag := range cmd.Flags() {
		// 匹配长名称 (flagName 带 "--" 前缀, flag.LongName() 不带)
		if strings.HasPrefix(flagName, "--") {
			if flag.LongName() == flagName[2:] { // 去掉 "--" 前缀再比较
				return flag
			}
		}
		// 匹配短名称 (flagName 带 "-" 前缀, flag.ShortName() 不带)
		if strings.HasPrefix(flagName, "-") && flag.ShortName() != "" {
			if flag.ShortName() == flagName[1:] { // 去掉 "-" 前缀再比较
				return flag
			}
		}
	}

	return nil
}

// getEnumValues 获取标志的枚举值列表
//
// 参数:
//   - flag: 标志实例
//
// 返回值:
//   - []string: 枚举值列表, 如果不是枚举类型则返回空切片
func getEnumValues(flag types.Flag) []string {
	// 检查是否是枚举类型
	if flag.Type() != types.FlagTypeEnum {
		return []string{}
	}

	// 调用 EnumValues 方法获取枚举值
	return flag.EnumValues()
}
