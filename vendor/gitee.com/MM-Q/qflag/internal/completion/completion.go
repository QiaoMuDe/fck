// Package completion 命令行自动补全功能
// 本文件实现了命令行自动补全的核心功能，包括标志补全、子命令补全、
// 参数值补全等，为用户提供便捷的命令行交互体验。
package completion

import (
	"bytes"
	"container/list"
	"os"
	"path"
	"path/filepath"

	"gitee.com/MM-Q/qflag/flags"
	"gitee.com/MM-Q/qflag/internal/types"
)

// 补全脚本生成相关常量定义
const (
	// DefaultFlagParamsCapacity 预估的标志参数初始容量
	// 基于常见CLI工具分析，大多数工具的标志数量在100-500之间
	DefaultFlagParamsCapacity = 256

	// NamesPerItem 每个标志/命令的名称数量(长名+短名)
	NamesPerItem = 2

	// MaxTraverseDepth 命令树遍历的最大深度限制
	// 防止循环引用导致的无限递归，一般CLI工具很少超过20层
	MaxTraverseDepth = 50
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
	// CompletionNotesCN 中文版本注意事项
	CompletionNotesCN = []string{
		"Windows环境: 需要PowerShell 5.1或更高版本以支持Register-ArgumentCompleter",
		"Linux环境: 需要bash 4.0或更高版本以支持关联数组特性",
		"请确保您的环境满足上述版本要求，否则自动补全功能可能无法正常工作",
	}

	// CompletionNotesEN 英文版本注意事项
	CompletionNotesEN = []string{
		"Windows environment: Requires PowerShell 5.1 or higher to support Register-ArgumentCompleter",
		"Linux environment: Requires bash 4.0 or higher to support associative array features",
		"Please ensure your environment meets the above version requirements, otherwise the auto-completion feature may not work properly",
	}
)

// 内置自动补全命令的示例使用（中文）
var CompletionExamplesCN = []types.ExampleInfo{
	{Description: "Linux环境 临时启用", Usage: "source <(%s --generate-shell-completion bash)"},
	{Description: "Linux环境 永久启用(添加到~/.bashrc)", Usage: "echo \"source <(%s --generate-shell-completion bash)\" >> ~/.bashrc"},
	//{Description: "Linux环境 系统级安装(至/etc/profile.d)", Usage: "sudo %s --generate-shell-completion bash > /etc/profile.d/qflag_completion.bash"},
	{Description: "Windows环境 临时启用", Usage: "%s --generate-shell-completion powershell | Out-String | Invoke-Expression"},
	{Description: "Windows环境 永久启用(添加到PowerShell配置文件)", Usage: "echo \"%s --generate-shell-completion powershell | Out-String | Invoke-Expression\" >> $PROFILE"},
	//{Description: "Windows环境 系统级安装(至ProgramFiles)", Usage: "%s --generate-shell-completion powershell > $env:ProgramFiles\\qflag\\completion.ps1"},
}

// 内置自动补全命令的示例使用（英文）
var CompletionExamplesEN = []types.ExampleInfo{
	{Description: "Linux environment temporary activation", Usage: "source <(%s --generate-shell-completion bash)"},
	{Description: "Linux environment permanent activation (add to ~/.bashrc)", Usage: "echo \"source <(%s --generate-shell-completion bash)\" >> ~/.bashrc"},
	//{Description: "Linux system-wide installation (to /etc/profile.d)", Usage: "sudo %s --generate-shell-completion bash > /etc/profile.d/qflag_completion.bash"},
	{Description: "Windows environment temporary activation", Usage: "%s --generate-shell-completion powershell | Out-String | Invoke-Expression"},
	{Description: "Windows environment permanent activation (add to PowerShell profile)", Usage: "echo \"%s --generate-shell-completion powershell | Out-String | Invoke-Expression\" >> $PROFILE"},
	//{Description: "Windows system-wide installation (to ProgramFiles)", Usage: "%s --generate-shell-completion powershell > $env:ProgramFiles\\qflag\\completion.ps1"},
}

// GenerateShellCompletion 生成shell自动补全脚本
//
// 参数:
//   - ctx: 命令上下文
//   - shellType: shell类型 ("bash", "pwsh", "powershell")
//
// 返回值：
//   - string: 自动补全脚本
//   - error: 错误信息
func GenerateShellCompletion(ctx *types.CmdContext, shellType string) (string, error) {
	// 缓冲区
	var buf bytes.Buffer

	// 检查生成自动补全脚本的必要条件
	if checkErr := validateCompletionGeneration(ctx); checkErr != nil {
		return "", checkErr
	}

	// 程序名称
	programName := filepath.Base(os.Args[0])

	// 获取根命令的补全选项
	rootCmdOpts := collectCompletionOptions(ctx)

	// 构建命令树条目缓冲区
	var cmdTreeEntries bytes.Buffer

	// 从根命令的子命令开始添加条目
	traverseCommandTree(&cmdTreeEntries, "", ctx.SubCmds, shellType)

	// 缓存标志参数，避免重复计算
	params := collectFlagParameters(ctx)

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

// traverseCommandTree 遍历命令树
//
// 参数:
//   - buf: 缓冲区
//   - rootPath: 根路径
//   - cmds: 命令列表
//   - shellType: 目标 shell 类型
func traverseCommandTree(buf *bytes.Buffer, rootPath string, cmdContexts []*types.CmdContext, shellType string) {
	// node 代表一个即将写入的“树条目”
	type node struct {
		name       string            // 真正写入的名字（长名或短名）
		parentPath string            // 父路径
		ctx        *types.CmdContext // 上下文
	}

	if len(cmdContexts) == 0 {
		return
	}

	// 初始化队列：把第一层命令按长短名拆开入队
	q := list.New()
	for _, c := range cmdContexts {
		if c == nil {
			continue
		}

		// 长名上下文入队
		if n := c.LongName; n != "" {
			q.PushBack(node{name: n, parentPath: rootPath, ctx: c})
		}

		// 短名上下文入队
		if n := c.ShortName; n != "" {
			q.PushBack(node{name: n, parentPath: rootPath, ctx: c})
		}
	}

	// 记录已经写了多少条，最后一条不补逗号
	var written int

	for q.Len() > 0 {
		cur := q.Remove(q.Front()).(node)

		// 计算当前节点完整路径
		var fullPath string
		if cur.parentPath == "/" {
			// 根节点
			fullPath = path.Join("/", cur.name) + "/"
		} else {
			// 子节点
			fullPath = path.Join("/", cur.parentPath, cur.name) + "/"
		}

		// 根据 shell 写树条目
		opts := collectCompletionOptions(cur.ctx)
		switch shellType {
		case flags.ShellBash: // Bash特定处理
			// 需要传递程序名称参数
			programName := filepath.Base(os.Args[0])
			generateBashCommandTreeEntry(buf, fullPath, opts, programName)
		case flags.ShellPwsh, flags.ShellPowershell: // Powershell特定处理
			generatePwshCommandTreeEntry(buf, fullPath, opts)
		}

		// 子节点入队（同样按长短名拆分）
		for _, sub := range cur.ctx.SubCmds {
			if sub == nil {
				continue
			}

			// 长名入队
			if n := sub.LongName; n != "" {
				q.PushBack(node{name: n, parentPath: fullPath, ctx: sub})
			}

			// 短名入队
			if n := sub.ShortName; n != "" {
				q.PushBack(node{name: n, parentPath: fullPath, ctx: sub})
			}
		}

		// 计数器累加
		written++

		// 只要不是最后一条就补逗号
		if q.Len() > 0 {
			// 如果是pwsh或者powershell，才处理
			if shellType == flags.ShellPwsh || shellType == flags.ShellPowershell {
				buf.WriteString(",\n")
			}
		}
	}
}

// collectFlagParameters 使用广度优先搜索收集整个命令树的所有标志参数信息
//
// 算法说明:
// 1. 使用BFS遍历整个命令树，确保所有层级的标志都被收集
// 2. 为每个标志生成唯一键(路径+标志名)，避免重复收集
// 3. 同时收集长短名标志，因为shell补全需要支持两种形式
// 4. 特殊处理枚举类型标志，提取可选值列表用于补全
//
// 参数:
//   - ctx: 根命令上下文，作为遍历起点
//
// 返回值:
//   - []FlagParam: 包含所有标志参数信息的切片，每个元素包含:
//   - CommandPath: 标志所属命令的路径
//   - Name: 标志名称(包含前缀，如--verbose或-v)
//   - Type: 参数需求类型("required"|"optional"|"none")
//   - ValueType: 参数值类型("path"|"string"|"number"|"enum"|"bool")
//   - EnumOptions: 枚举类型的可选值列表
func collectFlagParameters(ctx *types.CmdContext) []FlagParam {
	// 预分配切片容量，基于常见CLI工具的标志数量估算
	params := make([]FlagParam, 0, DefaultFlagParamsCapacity)

	// 使用map进行去重，键为"路径+标志名"的组合
	seen := make(map[string]struct{}, DefaultFlagParamsCapacity)

	// 队列 BFS
	type node struct {
		ctx  *types.CmdContext
		path string // 已标准化好的路径
	}
	q := []node{{ctx: ctx, path: "/"}}

	// add 添加一个标志到结果集
	add := func(prefix, name string, cur node, meta *flags.FlagMeta) {
		if name == "" {
			return
		}

		// 用“路径+flag”做唯一键
		key := cur.path + prefix + name
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}

		// 添加标志参数
		ft := meta.GetFlagType()
		param := FlagParam{
			CommandPath: cur.path,
			Name:        prefix + name,
			Type:        getParamTypeByFlagType(ft),
			ValueType:   getValueTypeByFlagType(ft),
		}

		// 如果是枚举标志，则获取枚举选项
		if ft == flags.FlagTypeEnum {
			if ef, ok := meta.GetFlag().(*flags.EnumFlag); ok {
				param.EnumOptions = ef.GetOptions()
			}
		}
		params = append(params, param)
	}

	// 遍历命令节点
	for len(q) > 0 {
		cur := q[0] // 获取当前节点
		q = q[1:]   // 移除当前节点

		// 遍历当前命令的所有标志
		for _, meta := range cur.ctx.FlagRegistry.GetFlagMetaList() {
			// 如果短标志不为空，则添加短标志
			if meta.GetShortName() != "" {
				add("-", meta.GetShortName(), cur, meta)
			}

			// 如果长标志不为空，则添加长标志
			if meta.GetLongName() != "" {
				add("--", meta.GetLongName(), cur, meta)
			}
		}

		// 子命令入队
		for _, sub := range cur.ctx.SubCmds {
			if sub == nil {
				continue
			}

			// 子命令路径
			subPath := path.Join(cur.path, sub.GetName()) + "/"

			// 入队
			q = append(q, node{ctx: sub, path: subPath})
		}
	}
	return params
}

// collectCompletionOptions 收集命令的补全选项，包括标志和子命令
//
// 参数：
//   - ctx: 命令上下文
//
// 返回值：
//   - 包含所有标志选项和子命令名称的字符串切片
func collectCompletionOptions(ctx *types.CmdContext) []string {
	if ctx == nil || ctx.FlagRegistry == nil {
		return nil
	}

	// 获取所有标志数量(长标志+短标志)
	flagCnt := ctx.FlagRegistry.GetAllFlagsCount()

	// 计算总容量（标志数量+子命令数量×每项名称数）
	// 每个子命令都有长名和短名，所以乘以NamesPerItem
	capacity := flagCnt + len(ctx.SubCmds)*NamesPerItem

	// 创建一个用于存储标志选项和子命令名称的切片
	seen := make(map[string]struct{}, capacity)

	// 定义一个添加选项的函数
	add := func(s string) {
		if s != "" {
			seen[s] = struct{}{}
		}
	}

	// 1. flags （同时展开长短名）
	for _, m := range ctx.FlagRegistry.GetFlagMetaList() {
		if m == nil {
			continue
		}

		if m.GetLongName() != "" {
			add("--" + m.GetLongName())
		}

		if m.GetShortName() != "" {
			add("-" + m.GetShortName())
		}
	}

	// 2. sub-commands（同时展开长短名）
	for _, sub := range ctx.SubCmds {
		if sub == nil {
			continue
		}
		add(sub.LongName)
		add(sub.ShortName)
	}

	// 3. 转回切片
	opts := make([]string, 0, len(seen))
	for k := range seen {
		opts = append(opts, k)
	}
	return opts
}
