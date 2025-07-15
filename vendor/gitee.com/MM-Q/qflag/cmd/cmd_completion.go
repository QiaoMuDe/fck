// cmd_completion.go - 自动补全命令的实现
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"gitee.com/MM-Q/qflag/flags"
)

// FlagParam 表示标志参数及其需求类型和值类型
type FlagParam struct {
	CommandPath string   // 命令路径，如 "/cmd/subcmd"
	Name        string   // 标志名称(保留原始大小写)
	Type        string   // 参数需求类型: "required"|"optional"|"none"
	ValueType   string   // 参数值类型: "path"|"string"|"number"|"enum"|"bool"等
	EnumOptions []string // 枚举类型的可选值列表
}

// 生成标志的注意事项
var (
	// completionNotesCN 中文版本注意事项
	completionNotesCN = []string{
		"Windows环境: 需要PowerShell 5.1或更高版本以支持Register-ArgumentCompleter",
		"Linux环境: 需要bash 4.0或更高版本以支持关联数组特性",
		"请确保您的环境满足上述版本要求，否则自动补全功能可能无法正常工作",
	}

	// completionNotesEN 英文版本注意事项
	completionNotesEN = []string{
		"Windows environment: Requires PowerShell 5.1 or higher to support Register-ArgumentCompleter",
		"Linux environment: Requires bash 4.0 or higher to support associative array features",
		"Please ensure your environment meets the above version requirements, otherwise the auto-completion feature may not work properly",
	}
)

// 内置自动补全命令的示例使用（中文）
var completionExamplesCN = []ExampleInfo{
	{Description: "Linux环境 临时启用", Usage: "source <(%s --generate-shell-completion bash)"},
	{Description: "Linux环境 永久启用(添加到~/.bashrc)", Usage: "echo \"source <(%s --generate-shell-completion bash)\" >> ~/.bashrc"},
	{Description: "Linux环境 系统级安装(至/etc/profile.d)", Usage: "sudo %s --generate-shell-completion bash > /etc/profile.d/qflag_completion.bash"},
	{Description: "Windows环境 临时启用", Usage: "%s --generate-shell-completion powershell | Out-String | Invoke-Expression"},
	{Description: "Windows环境 永久启用(添加到PowerShell配置文件)", Usage: "echo \"%s --generate-shell-completion powershell | Out-String | Invoke-Expression\" >> $PROFILE"},
	{Description: "Windows环境 系统级安装(至ProgramFiles)", Usage: "%s --generate-shell-completion powershell > $env:ProgramFiles\\qflag\\completion.ps1"},
}

// 内置自动补全命令的示例使用（英文）
var completionExamplesEN = []ExampleInfo{
	{Description: "Linux environment temporary activation", Usage: "source <(%s --generate-shell-completion bash)"},
	{Description: "Linux environment permanent activation (add to ~/.bashrc)", Usage: "echo \"source <(%s --generate-shell-completion bash)\" >> ~/.bashrc"},
	{Description: "Linux system-wide installation (to /etc/profile.d)", Usage: "sudo %s --generate-shell-completion bash > /etc/profile.d/qflag_completion.bash"},
	{Description: "Windows environment temporary activation", Usage: "%s --generate-shell-completion powershell | Out-String | Invoke-Expression"},
	{Description: "Windows environment permanent activation (add to PowerShell profile)", Usage: "echo \"%s --generate-shell-completion powershell | Out-String | Invoke-Expression\" >> $PROFILE"},
	{Description: "Windows system-wide installation (to ProgramFiles)", Usage: "%s --generate-shell-completion powershell > $env:ProgramFiles\\qflag\\completion.ps1"},
}

// getValueTypeByFlagType 根据标志类型获取值类型
//
// 参数:
//   - flagType - 标志类型
//
// 返回值:
//   - string: 值类型
func getValueTypeByFlagType(flagType flags.FlagType) string {
	switch flagType {
	case flags.FlagTypeBool:
		return "bool"
	case flags.FlagTypeInt, flags.FlagTypeUint16, flags.FlagTypeUint32, flags.FlagTypeUint64, flags.FlagTypeInt64, flags.FlagTypeFloat64:
		return "number"
	case flags.FlagTypePath:
		return "path"
	case flags.FlagTypeEnum:
		return "enum"
	case flags.FlagTypeIP4, flags.FlagTypeIP6:
		return "ip"
	case flags.FlagTypeURL:
		return "url"
	default:
		return "string"
	}
}

// getParamTypeByFlagType 根据标志类型获取参数需求类型
//
// 参数:
//   - flagType - 标志类型
//
// 返回值:
//   - string: 参数需求类型
func getParamTypeByFlagType(flagType flags.FlagType) string {
	if flagType == flags.FlagTypeBool {
		return "none"
	}
	return "required"
}

// traverseCommandTree 通用的命令树遍历函数
//
// 参数:
//   - cmdTreeEntries - 用于存储命令树条目的缓冲区
//   - parentPath - 父命令路径
//   - cmds - 子命令列表
//   - shellType - shell类型 ("bash", "pwsh", "powershell")
func traverseCommandTree(cmdTreeEntries *bytes.Buffer, parentPath string, cmds []*Cmd, shellType string) {
	// 使用队列实现广度优先遍历
	type cmdNode struct {
		cmd        *Cmd
		parentPath string
	}
	queue := make([]cmdNode, 0, len(cmds))

	// 初始化队列并计算总节点数
	totalNodes := new(int)
	*totalNodes = 0
	for _, cmd := range cmds {
		if cmd != nil {
			// 为长名称和短名称分别计数
			if cmd.LongName() != "" {
				*totalNodes++
			}
			if cmd.ShortName() != "" {
				*totalNodes++
			}
			queue = append(queue, cmdNode{cmd: cmd, parentPath: parentPath})
		}
	}

	// 初始化已处理节点数
	processedNodes := new(int)
	*processedNodes = 0

	// 定义处理命令名称的匿名函数
	processCmdName := func(name string, currentParentPath string, cmd *Cmd, queue *[]cmdNode) {
		// 检查命令名称和命令是否有效以及计数器是否有效
		if name == "" || cmd == nil {
			return
		}

		// 获取子命令补全选项
		cmdOpts := cmd.collectCompletionOptions()

		// 构建命令路径（改进版）
		var cmdPath string
		if currentParentPath != "/" && currentParentPath != "" {
			cmdPath = path.Join("/", currentParentPath, name) + "/"
		} else if currentParentPath == "/" {
			cmdPath = path.Join(currentParentPath, name) + "/"
		} else {
			cmdPath = path.Join("/", name) + "/"
		}

		// 根据shell类型调用对应的处理函数
		switch shellType {
		case flags.ShellBash: // Bash
			// 调用generateBashCommandTreeEntry函数生成Bash自动补全条目
			generateBashCommandTreeEntry(cmdTreeEntries, cmdPath, cmdOpts)

		case flags.ShellPwsh, flags.ShellPowershell: // Powershell和Pwsh
			// 调用generatePwshCommandTreeEntry函数生成Powershell自动补全条目
			generatePwshCommandTreeEntry(cmdTreeEntries, cmdPath, cmdOpts)

			// 判断是否为最后一个节点, 如果不是最后一个条目，添加逗号
			if *processedNodes != *totalNodes {
				cmdTreeEntries.WriteString(",\n")
			}
		}

		// 将子命令加入队列
		for _, subCmd := range cmd.subCmds {
			if subCmd != nil {
				// 添加命令
				*queue = append(*queue, cmdNode{cmd: subCmd, parentPath: cmdPath})

				// 添加长名称
				if subCmd.LongName() != "" {
					*totalNodes++ // 总节点数加1
				}

				// 添加短名称
				if subCmd.ShortName() != "" {
					*totalNodes++ // 总节点数加1
				}
			}
		}
	}
	// 处理队列中的所有命令
	for len(queue) > 0 {
		// 出队
		node := queue[0]                     // 获取当前命令节点
		queue = queue[1:]                    // 移除已处理的元素
		cmd := node.cmd                      // 获取当前命令
		currentParentPath := node.parentPath // 获取当前命令的父路径

		// 处理长命名
		if cmd.LongName() != "" {
			*processedNodes++ // 已处理节点数加1
			processCmdName(cmd.LongName(), currentParentPath, cmd, &queue)
		}

		// 处理短命名
		if cmd.ShortName() != "" {
			*processedNodes++ // 已处理节点数加1
			processCmdName(cmd.ShortName(), currentParentPath, cmd, &queue)
		}
	}
}

// generateShellCompletion 生成shell自动补全脚本
//
// 参数:
//   - shellType: shell类型 ("bash", "pwsh", "powershell")
//
// 返回值：
//   - string: 自动补全脚本
//   - error: 错误信息
func (c *Cmd) generateShellCompletion(shellType string) (string, error) {
	// 缓冲区
	var buf bytes.Buffer

	// 检查生成自动补全脚本的必要条件
	if checkErr := validateCompletionGeneration(c); checkErr != nil {
		return "", checkErr
	}

	// 程序名称
	programName := filepath.Base(os.Args[0])

	// 获取根命令的补全选项
	rootCmdOpts := c.collectCompletionOptions()

	// 构建命令树条目缓冲区
	var cmdTreeEntries bytes.Buffer

	// 从根命令的子命令开始添加条目
	traverseCommandTree(&cmdTreeEntries, "", c.subCmds, shellType)

	// 缓存标志参数，避免重复计算
	params := c.collectFlagParameters()

	// 根据shell类型调用对应的处理函数
	switch shellType {
	case flags.ShellBash: // Bash特定处理
		generateBashCompletion(&buf, params, rootCmdOpts, cmdTreeEntries.String(), programName)
	case flags.ShellPwsh, flags.ShellPowershell: // PowerShell特定处理
		generatePwshCompletion(&buf, params, rootCmdOpts, cmdTreeEntries.String(), programName)
	}

	// 返回自动补全脚本
	return buf.String(), nil
}

// collectFlagParameters 收集所有命令标志参数需求，返回标志名称到参数需求类型的映射
// 参数需求类型: "required"|"optional"|"none"
func (c *Cmd) collectFlagParameters() []FlagParam {
	// 基于命令数量和平均标志密度预分配容量 (每个命令约8个标志)
	initialCapacity := len(c.subCmds)*8 + 16        // 额外+16应对根命令标志
	params := make([]FlagParam, 0, initialCapacity) // 使用切片存储标志参数需求
	seen := make(map[string]bool)                   // 使用原始标志名称作为键，区分大小写

	// 定义匿名函数处理标志添加逻辑，包含参数类型判断和命令路径
	addFlagParam := func(flag *flags.FlagMeta, prefix, opt string, cmdPath string) {
		if opt == "" {
			return
		}

		// 拼接标志名称
		flagName := prefix + opt

		// 只有在标志名称未被添加过时才添加
		if !seen[flagName] {
			seen[flagName] = true // 标记为已添加

			// 根据标志类型获取参数类型和值类型
			flagType := flag.GetFlagType()
			paramType := getParamTypeByFlagType(flagType)
			valueType := getValueTypeByFlagType(flagType)
			var enumOptions []string

			if flagType == flags.FlagTypeEnum {
				if enumFlag, ok := flag.GetFlag().(*flags.EnumFlag); ok {
					enumOptions = enumFlag.GetOptions()
				}
			}

			// 添加标志参数需求，包含命令路径
			params = append(params, FlagParam{CommandPath: cmdPath, Name: flagName, Type: paramType, ValueType: valueType, EnumOptions: enumOptions})
		}
	}

	// 使用队列实现广度优先遍历，记录命令路径
	type cmdNode struct {
		cmd        *Cmd
		parentPath string
	}
	queue := []cmdNode{{cmd: c, parentPath: ""}}

	// 循环遍历命令树
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		cmd := node.cmd
		currentParentPath := node.parentPath

		// 构建当前命令路径，拼接父路径和当前命令名
		cmdPath := "/"
		if currentParentPath != "" {
			cmdPath = currentParentPath + cmd.Name() + "/"
		}

		// 收集当前命令的标志 - 同时处理长短选项
		for _, flag := range cmd.flagRegistry.GetAllFlagMetas() {
			// 处理短选项
			addFlagParam(flag, "-", flag.GetShortName(), cmdPath)
			// 处理长选项
			addFlagParam(flag, "--", flag.GetLongName(), cmdPath)
		}

		// 将子命令加入队列
		for _, subCmd := range cmd.subCmds {
			queue = append(queue, cmdNode{cmd: subCmd, parentPath: cmdPath})
		}
	}

	return params
}

// collectCompletionOptions 收集命令的补全选项，包括标志和子命令
//
// 参数：
//   - c: 当前命令实例
//
// 返回值：
//   - 包含所有标志选项和子命令名称的字符串切片
func (c *Cmd) collectCompletionOptions() []string {

	// 防御性检查
	if c == nil || c.flagRegistry == nil {
		return []string{}
	}

	// 获取所有标志
	flags := c.flagRegistry.GetAllFlagMetas()

	// 获取所有子命令
	subCmdMap := c.subCmdMap

	// 预分配切片容量，减少动态扩容的开销
	capacity := len(flags)*2 + len(subCmdMap) // 每个标志和子命令(子命令map已包含长短名)最多2个选项
	opts := make([]string, 0, capacity)

	// 获取所有长选项和短选项(为空时不会循环)
	for _, m := range flags {
		if m.GetLongName() != "" {
			opts = append(opts, "--"+m.GetLongName())
		}
		if m.GetShortName() != "" {
			opts = append(opts, "-"+m.GetShortName())
		}
	}

	// 获取所有子命令(为空时不会循环)
	for subCmd := range subCmdMap {
		opts = append(opts, subCmd)
	}

	// 返回所有选项
	return opts
}

// validateCompletionGeneration 验证补全脚本生成所需的命令状态
//
// 参数:
//   - c: 命令实例
//
// 返回值:
//   - error: 如果验证失败, 返回相应的错误信息; 否则返回nil
func validateCompletionGeneration(c *Cmd) error {
	if c == nil {
		return fmt.Errorf("command instance is nil")
	}
	if c.parentCmd != nil {
		return fmt.Errorf("invalid command state: not a root command")
	}
	if c.flagRegistry == nil {
		return fmt.Errorf("invalid command state: flag registry is nil")
	}
	return nil
}
