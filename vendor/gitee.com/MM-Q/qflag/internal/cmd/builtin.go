package cmd

import (
	"fmt"
	"os"

	"gitee.com/MM-Q/qflag/internal/completion"
	"gitee.com/MM-Q/qflag/internal/types"
)

// createCompleteCmd 创建动态补全子命令
//
// 参数:
//   - root: 根命令实例
//
// 返回值:
//   - error: 注册子命令时失败的错误信息
//
// 功能说明:
//   - 创建一个隐藏的子命令，用于执行动态补全
//   - 接收当前命令行作为参数
//   - 计算并返回补全结果
func createCompleteCmd(root *Cmd) error {
	// 检查是否已存在相同名称的子命令
	if root.cmdRegistry.Has(types.CompleteCmdName) {
		return nil
	}

	// 创建子命令
	cmd := NewCmd(types.CompleteCmdName, "", types.ExitOnError)
	opts := &CmdOpts{
		Desc:               "内部动态补全指令，为 Shell 补全脚本提供实时命令树查询服务",
		Hidden:             true, // 隐藏在命令列表中
		DisableFlagParsing: true, // 禁用标志解析，只处理指令参数
		Examples: map[string]string{
			"执行模糊匹配补全": fmt.Sprintf("%s %s %s <模式> <候选1> [候选2] ...", root.Name(), types.CompleteCmdName, types.InstructionFuzzy),
			"计算上下文路径":  fmt.Sprintf("%s %s %s <arg0> [arg1] ...", root.Name(), types.CompleteCmdName, types.InstructionContext),
			"获取候选选项":   fmt.Sprintf("%s %s %s <上下文路径>", root.Name(), types.CompleteCmdName, types.InstructionCandidates),
			"获取枚举值":    fmt.Sprintf("%s %s %s <上下文路径> <标志名>", root.Name(), types.CompleteCmdName, types.InstructionEnum),
			"统一获取补全信息": fmt.Sprintf("%s %s %s <当前输入> <前一个输入> [子命令参数...]", root.Name(), types.CompleteCmdName, types.InstructionAll),
		},
		Notes: []string{
			"本命令为内部命令，用于 Shell 自动补全脚本动态获取补全信息",
			"",
			"调试模式：",
			"  设置环境变量 QFLAG_COMPLETE_DEBUG=1 可启用调试输出，",
			"  此时所有错误信息将显示在终端，便于排查问题",
			"",
			"生产模式：",
			"  默认情况下所有错误都被静默处理，避免干扰补全脚本解析",
		},
		RunFunc: func(c types.Command) error {
			args := c.Args()

			// 校验参数数量：至少需要指定指令
			if len(args) < 1 {
				// 检查是否启用调试模式
				if os.Getenv("QFLAG_COMPLETE_DEBUG") == "1" {
					return fmt.Errorf("usage: %s %s <instruction> [args...]", root.Name(), types.CompleteCmdName)
				}
				// 静默处理参数错误，不影响补全脚本
				return nil
			}

			// 拆分参数：第一个是指令，后续是参数列表
			instruction := args[0]
			params := []string{}
			if len(args) > 1 {
				params = args[1:]
			}

			// 传递给 completion 包处理，传入 root 命令
			err := completion.HandleDynamicComplete(root, instruction, params)

			// 检查是否启用调试模式
			if err != nil && os.Getenv("QFLAG_COMPLETE_DEBUG") == "1" {
				// 调试模式下返回错误
				return err
			}

			// 静默处理错误，避免错误信息干扰补全脚本
			return nil
		},
	}

	// 设置选项
	if err := cmd.ApplyOpts(opts); err != nil {
		return err
	}

	// 设置补全子命令的父命令
	// 直接设置防止锁问题
	cmd.parent = root

	// 注册子命令
	if err := root.cmdRegistry.Register(cmd); err != nil {
		return fmt.Errorf("register subcommand '%s' failed in '%s': %w", cmd.Name(), root.Name(), err)
	}

	return nil
}
