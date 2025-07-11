// cmd_completion.go - 自动补全命令的实现
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/qflag/flags"
)

// FlagParam 表示标志参数及其需求类型
type FlagParam struct {
	Name string // 标志名称(保留原始大小写)
	Type string // 参数需求类型: "required"|"optional"|"none"
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

// 补全脚本模板常量
const (
	// Bash补全模板
	//BashFunctionHeader   = "#!/usr/bin/env bash\n\n_%s() {\n\tlocal cur prev words cword context opts i arg\n\tCOMPREPLY=()\n\n\t# 使用_get_comp_words_by_ref获取补全参数, 提高健壮性\n\tif [[ -z \"${_get_comp_words_by_ref}\" ]]; then\n\t\t# 兼容旧版本Bash补全环境\n\t\twords=(\"${COMP_WORDS[@]}\")\n\t\tcword=$COMP_CWORD\n\telse\n\t\t_get_comp_words_by_ref -n =: cur prev words cword\n\tfi\n\n\tcur=\"${words[cword]}\"\n\tprev=\"${words[cword-1]}\"\n\n\t# 构建命令树结构\n\tdeclare -A cmd_tree\n\tcmd_tree[/]=\"%s\"\n%s\n\n\t# 查找当前命令上下文\n\tlocal context=\"/\"\n\tlocal i\n\tfor ((i=1; i < cword; i++)); do\n\t\tlocal arg=\"${words[i]}\"\n\t\tif [[ -n \"${cmd_tree[$context$arg/]}\" ]]; then\n\t\t\tcontext=\"$context$arg/\"\n\t\tfi\n\tdone\n\n\t# 获取当前上下文可用选项\n\topts=\"${cmd_tree[$context]}\"\n\tCOMPREPLY=($(compgen -W \"${opts}\" -- ${cur}))\n\treturn 0\n\t}\n\ncomplete -F _%s %s\n" // Bash补全函数头部
	BashFunctionHeader = `#!/usr/bin/env bash

_%s() {
	local cur prev words cword context opts i arg
	COMPREPLY=()

	# 使用_get_comp_words_by_ref获取补全参数, 提高健壮性
	if [[ -z "${_get_comp_words_by_ref}" ]]; then
		# 兼容旧版本Bash补全环境
		words=("${COMP_WORDS[@]}")
		cword=$COMP_CWORD
	else
		_get_comp_words_by_ref -n =: cur prev words cword
	fi

	cur="${words[cword]}"
	prev="${words[cword-1]}"

	# 构建命令树结构
	declare -A cmd_tree
	cmd_tree[/]="%s"
%s

	# 查找当前命令上下文
	local context="/"
	local i
	for ((i=1; i < cword; i++)); do
		local arg="${words[i]}"
		if [[ -n "${cmd_tree[$context$arg/]}" ]]; then
			context="$context$arg/"
		fi
	done

	# 获取当前上下文可用选项
	opts="${cmd_tree[$context]}"
	# 添加-o filenames选项处理特殊字符和空格
	COMPREPLY=($(compgen -o filenames -W "${opts}" -- ${cur}))

	# 模糊匹配与纠错提示：当无精确匹配时，从所有选项中查找包含关键词的项
	if [[ ${#COMPREPLY[@]} -eq 0 ]]; then
		local all_opts=()
		# 使用循环安全收集所有选项，避免空格分割问题
		for path in "${!cmd_tree[@]}"; do
			for opt in ${cmd_tree[$path]}; do
				all_opts+=($opt)
			done
		done
		# 去重并生成补全结果
		local unique_opts=($(printf "%%s\n" "${all_opts[@]}" | sort -u))
		COMPREPLY=($(compgen -o filenames -W "${unique_opts[*]}" -- ${cur}))
	fi

	return 0
	}

complete -F _%s %s
`
	BashCommandTreeEntry = "\tcmd_tree[/%s/]=\"%s\"\n" // 命令树条目格式

	// PowerShell补全模板
	//PwshFunctionHeader   = "Register-ArgumentCompleter -CommandName %s -ScriptBlock {\n    param($wordToComplete, $commandAst, $cursorPosition, $commandName, $parameterName)\n\n    # 标志参数需求映射\n    $flagParams = @{\n%s\n    }\n\n    # 构建命令树结构\n    $cmdTree = @{\n        '' = '%s'\n%s\n    }\n\n    # 解析命令行参数获取当前上下文\n    $context = ''\n    $args = $commandAst.CommandElements | Select-Object -Skip 1 | ForEach-Object { $_.ToString() }\n    $index = 0\n    $count = $args.Count\n\n    while ($index -lt $count) {\n        $arg = $args[$index]\n        # 处理选项参数及其值\n        if ($arg -like '-*' -and $flagParams.ContainsKey($arg)) {\n            $paramType = $flagParams[$arg]\n            $index++\n            \n            # 根据参数类型决定是否跳过下一个参数\n            if ($paramType -eq 'required' -or ($paramType -eq 'optional' -and $index -lt $count -and $args[$index] -notlike '-*')) {\n                $index++\n            }\n            continue\n        }\n\n        $nextContext = if ($context) { \"$context.$arg\" } else { $arg }\n        if ($cmdTree.ContainsKey($nextContext)) {\n            $context = $nextContext\n            $index++\n        } else {\n            break\n        }\n    }\n\n    # 获取当前上下文可用选项并过滤\n    $options = @()\n    if ($cmdTree.ContainsKey($context)) {\n        $options = $cmdTree[$context] -split ' ' | Where-Object { $_ -like \"$wordToComplete*\" }\n    }\n\n    $options | ForEach-Object { [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterName', $_) }\n}\n" // PowerShell补全函数头部
	PwshFunctionHeader = `Register-ArgumentCompleter -CommandName %s -ScriptBlock {
		param($wordToComplete, $commandAst, $cursorPosition, $commandName, $parameterName)
	
		# 标志参数需求数组(保留原始大小写)
		$flagParams = @(
%s      )
	
		# 构建命令树结构
		$cmdTree = @{
%s      }
	
		# 解析命令行参数获取当前上下文
		$context = ''
		$args = $commandAst.CommandElements | Select-Object -Skip 1 | ForEach-Object { $_.ToString() }
		$index = 0
		$count = $args.Count
	
		while ($index -lt $count) {
			$arg = $args[$index]
			# 使用大小写敏感匹配查找标志
			$paramInfo = $flagParams | Where-Object { $_.Name -ceq $arg } | Select-Object -First 1
			if ($paramInfo) {
				$paramType = $paramInfo.Type
				$index++
				
				# 根据参数类型决定是否跳过下一个参数
				if ($paramType -eq 'required' -or ($paramType -eq 'optional' -and $index -lt $count -and $args[$index] -notlike '-*')) {
					$index++
				}
				continue
			}
	
			$nextContext = if ($context) { "$context.$arg" } else { $arg }
			if ($cmdTree.ContainsKey($nextContext)) {
				$context = $nextContext
				$index++
			} else {
				break
			}
		}
	
		# 获取当前上下文可用选项并过滤
		$options = @()
		if ($cmdTree.ContainsKey($context)) {
			$options = $cmdTree[$context] -split ' ' | Where-Object { $_ -like "$wordToComplete*" }
		}
	
		# 模糊匹配与纠错提示：当无精确匹配时，从所有选项中查找包含关键词的项
		if (-not $options) {
			# 递归收集所有层级的选项
			$allOptions = @()
			$cmdTree.Values | ForEach-Object { $allOptions += $_ -split ' ' }
			$options = $allOptions | Select-Object -Unique | Where-Object { $_ -like "*$wordToComplete*" }
		}
	
		$options | ForEach-Object { [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterName', $_) }
	}`
	PwshCommandTreeEntryRoot = "\t\t'' = '%s'\n"                    // 根命令树条目格式
	PwshCommandTreeEntry     = "\t\t'%s' = '%s'\n"                  // 命令树条目格式
	PwshCommandTreeOption    = "\t\t@{ Name = '%s'; Type = '%s'}\n" // 选项参数需求条目格式
)

// addSubCommandsBash 迭代方式添加子命令到命令树，替代递归实现
//
// 参数:
//   - cmdTreeEntries - 用于存储命令树条目的缓冲区
//   - parentPath - 父命令路径
//   - cmds - 子命令列表
func addSubCommandsBash(cmdTreeEntries *bytes.Buffer, parentPath string, cmds []*Cmd) {
	// 使用队列实现广度优先遍历
	type cmdNode struct {
		cmd        *Cmd
		parentPath string
	}
	queue := make([]cmdNode, 0, len(cmds))

	// 初始化队列
	for _, cmd := range cmds {
		queue = append(queue, cmdNode{cmd: cmd, parentPath: parentPath})
	}

	// 定义处理命令名称的匿名函数 (直接传入cmd对象)
	processCmdName := func(name string, currentParentPath string, cmd *Cmd, queue *[]cmdNode) {
		if name == "" || cmd == nil {
			return
		}
		// 使用strings.Builder优化路径拼接
		var cmdPathBuilder strings.Builder
		cmdPathBuilder.WriteString(currentParentPath)
		cmdPathBuilder.WriteString(name)
		cmdPathBuilder.WriteString("/")
		cmdPath := cmdPathBuilder.String()
		cmdPathBuilder.Reset()

		// 去除末尾的斜杠 (通过切片操作避免额外字符串分配)
		var trimmedPath string
		if len(cmdPath) > 0 && cmdPath[len(cmdPath)-1] == '/' {
			// 显式检查斜杠，避免极端情况下的切片越界
			trimmedPath = cmdPath[:len(cmdPath)-1]
		} else {
			trimmedPath = cmdPath // 保留原始路径，避免切片越界
		}

		// 获取子命令补全选项 (通过cmd对象直接调用方法)
		cmdOpts := cmd.collectCompletionOptions()

		// 写入命令树条目
		fmt.Fprintf(cmdTreeEntries, BashCommandTreeEntry, trimmedPath, strings.Join(cmdOpts, " "))

		// 将子命令加入队列 (通过cmd对象访问子命令)
		for _, subCmd := range cmd.subCmds {
			*queue = append(*queue, cmdNode{cmd: subCmd, parentPath: cmdPath})
		}
	}

	// 处理队列中的所有命令
	for len(queue) > 0 {
		// 出队
		node := queue[0]                     // 获取当前命令节点
		queue = queue[1:]                    // 移除已处理的元素
		cmd := node.cmd                      // 获取当前命令
		currentParentPath := node.parentPath // 获取当前命令的父路径

		// 直接传递cmd对象给匿名函数
		processCmdName(cmd.LongName(), currentParentPath, cmd, &queue)
		processCmdName(cmd.ShortName(), currentParentPath, cmd, &queue)
	}
}

// generateBashCompletion 生成Bash自动补全脚本
//
// 返回值：
//   - string: Bash自动补全脚本
func (c *Cmd) generateBashCompletion() (string, error) {
	// 缓冲区
	var buf bytes.Buffer

	// 检查命令树注册表是否为空
	if c == nil {
		return "", fmt.Errorf("command instance is nil")
	}

	// 如果不是父命令则返回错误
	if c.parentCmd != nil {
		return "", fmt.Errorf("invalid command state: not a root command")
	}

	// 检查标志注册表是否为空
	if c.flagRegistry == nil {
		return "", fmt.Errorf("invalid command state: flag registry is nil")
	}

	// 程序名称
	programName := filepath.Base(os.Args[0])

	// 获取根命令的补全选项
	rootCmdOpts := c.collectCompletionOptions()

	// 构建命令树条目缓冲区
	var cmdTreeEntries bytes.Buffer

	// 添加根命令选项
	rootOpts := strings.Join(rootCmdOpts, " ")

	// 从根命令的子命令开始添加条目
	addSubCommandsBash(&cmdTreeEntries, "", c.subCmds)

	// 写入补全函数头部和命令树
	fmt.Fprintf(&buf, BashFunctionHeader, programName, rootOpts, cmdTreeEntries.String(), programName, programName)

	return buf.String(), nil
}

// generatePwshCompletion 生成PowerShell自动补全脚本
//
// 返回值：
//   - string: PowerShell自动补全脚本
func (c *Cmd) generatePwshCompletion() (string, error) {
	// 缓冲区
	var buf bytes.Buffer

	// 检查命令树注册表是否为空
	if c == nil {
		return "", fmt.Errorf("command instance is nil")
	}

	// 如果不是父命令则返回错误
	if c.parentCmd != nil {
		return "", fmt.Errorf("invalid command state: not a root command")
	}

	// 检查标志注册表是否为空
	if c.flagRegistry == nil {
		return "", fmt.Errorf("invalid command state: flag registry is nil")
	}

	// 程序名称
	programName := filepath.Base(os.Args[0])

	// 构建命令树条目缓冲区
	var cmdTreeEntries bytes.Buffer

	// 获取根命令的补全选项
	rootCmdOpts := c.collectCompletionOptions()

	// 添加根命令选项
	rootOpts := strings.Join(rootCmdOpts, " ")

	// 从根命令的子命令开始添加条目
	addSubCommandsPwsh(&cmdTreeEntries, "", c.subCmds, rootOpts)

	// 构建标志参数需求映射
	var flagParamsBuf bytes.Buffer

	// 收集当前命令的标志(从根命令开始)
	flagParams := c.collectFlagParameters() // 现在返回[]FlagParam

	// 写入标志参数需求条目 - 使用数组而非哈希表
	for _, param := range flagParams {
		fmt.Fprintf(&flagParamsBuf, PwshCommandTreeOption, param.Name, param.Type)
	}

	// 写入补全函数头部和命令树
	fmt.Fprintf(&buf, PwshFunctionHeader, programName, flagParamsBuf.String(), cmdTreeEntries.String())

	return buf.String(), nil
}

// addSubCommandsPwsh 迭代方式添加子命令到PowerShell命令树，替代递归实现
//
// 参数:
//   - cmdTreeEntries - 用于存储命令树条目的缓冲区
//   - parentPath - 父命令路径(使用.作为分隔符)
//   - cmds - 子命令列表
func addSubCommandsPwsh(cmdTreeEntries *bytes.Buffer, parentPath string, cmds []*Cmd, rootOpts string) {
	// 使用队列实现广度优先遍历
	type cmdNode struct {
		cmd        *Cmd
		parentPath string
	}
	queue := make([]cmdNode, 0, len(cmds))

	// 初始化队列
	for _, cmd := range cmds {
		queue = append(queue, cmdNode{cmd: cmd, parentPath: parentPath})
	}

	// 写入根命令条目
	fmt.Fprintf(cmdTreeEntries, PwshCommandTreeEntryRoot, rootOpts)

	// 定义处理命令名称的匿名函数 (抽取重复逻辑)
	processCmdName := func(name string, currentParentPath string, cmd *Cmd, cmdOpts []string) {
		if name == "" {
			return
		}
		// 构建命令路径(使用.作为分隔符)
		cmdPath := currentParentPath
		if cmdPath != "" {
			cmdPath += fmt.Sprintf(".%s", name)
		} else {
			cmdPath = name
		}

		// 写入命令树条目
		fmt.Fprintf(cmdTreeEntries, PwshCommandTreeEntry, cmdPath, strings.Join(cmdOpts, " "))

		// 将子命令加入队列
		for _, subCmd := range cmd.subCmds {
			queue = append(queue, cmdNode{cmd: subCmd, parentPath: cmdPath})
		}
	}

	// 处理队列中的所有命令
	for len(queue) > 0 {
		// 出队
		node := queue[0]                     // 获取当前命令节点
		queue = queue[1:]                    // 移除已处理的元素
		cmd := node.cmd                      // 获取当前命令
		currentParentPath := node.parentPath // 获取当前命令的父命令路径

		if cmd == nil {
			continue
		}

		// 获取子命令补全选项
		cmdOpts := cmd.collectCompletionOptions()

		// 处理长命令和短命令 (通过匿名函数消除重复)
		processCmdName(cmd.LongName(), currentParentPath, cmd, cmdOpts)  // 处理长命令
		processCmdName(cmd.ShortName(), currentParentPath, cmd, cmdOpts) // 处理短命令
	}
}

// collectFlagParameters 收集所有命令标志参数需求，返回标志名称到参数需求类型的映射
// 参数需求类型: "required"|"optional"|"none"
func (c *Cmd) collectFlagParameters() []FlagParam {
	params := make([]FlagParam, 0) // 使用切片存储标志参数需求
	seen := make(map[string]bool)  // 使用原始标志名称作为键，区分大小写

	// 定义匿名函数处理标志添加逻辑，包含参数类型判断
	addFlagParam := func(flag *flags.FlagMeta, prefix, opt string) {
		if opt == "" {
			return
		}

		// 拼接标志名称
		flagName := prefix + opt

		// 只有在标志名称未被添加过时才添加
		if !seen[flagName] {
			seen[flagName] = true // 标记为已添加

			// 根据标志类型设置参数类型
			paramType := "required"
			if flag.GetFlagType() == flags.FlagTypeBool {
				paramType = "none"
			}

			// 添加标志参数需求
			params = append(params, FlagParam{Name: flagName, Type: paramType})
		}
	}

	// 使用队列实现广度优先遍历
	queue := make([]*Cmd, 0, 10)
	queue = append(queue, c)

	// 遍历队列中的所有命令
	for len(queue) > 0 {
		cmd := queue[0]
		queue = queue[1:]

		// 收集当前命令的标志 - 同时处理长短选项
		for _, flag := range cmd.flagRegistry.GetAllFlagMetas() {
			// 处理短选项
			addFlagParam(flag, "-", flag.GetShortName())
			// 处理长选项
			addFlagParam(flag, "--", flag.GetLongName())
		}

		// 将子命令加入队列
		queue = append(queue, cmd.subCmds...)
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
			opts = append(opts, fmt.Sprint("--", m.GetLongName()))
		}
		if m.GetShortName() != "" {
			opts = append(opts, fmt.Sprint("-", m.GetShortName()))
		}
	}

	// 获取所有子命令(为空时不会循环)
	for subCmd := range subCmdMap {
		opts = append(opts, subCmd)
	}

	// 返回所有选项
	return opts
}
