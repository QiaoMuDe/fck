// Package completion 自动补全内部实现
// 本文件包含了自动补全功能的内部实现逻辑, 提供补全算法、
// 匹配策略等核心功能的底层支持。
package completion

import (
	"bytes"
	"container/list"
	_ "embed"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"gitee.com/MM-Q/qflag/internal/types"
)

// FlagParam 表示标志参数及其需求类型和值类型
type FlagParam struct {
	CommandPath string   // 命令路径, 如 "/cmd/subcmd"
	Name        string   // 标志名称(保留原始大小写)
	Type        string   // 参数需求类型: "required"|"optional"|"none"
	ValueType   string   // 参数值类型: "path"|"string"|"number"|"enum"|"bool"等
	EnumOptions []string // 枚举类型的可选值列表
}

//go:embed templates/bash.tmpl
var bashTemplate string

//go:embed templates/pwsh.tmpl
var pwshTemplate string

// GenAndPrint 生成并打印补全脚本
//
// 参数:
//   - cmd: 要生成补全脚本的命令
//   - shellType: Shell类型 (bash, pwsh, powershell)
func GenAndPrint(cmd types.Command, shellType string) {
	st, err := Generate(cmd, shellType)
	if err != nil {
		fmt.Printf("Error generating completion script: %v\n", err)
	}
	fmt.Println(st)
}

// Generate 生成补全脚本
//
// 参数:
//   - cmd: 要生成补全脚本的命令
//   - shellType: Shell类型 (bash, pwsh, powershell)
//
// 返回值:
//   - string: 生成的补全脚本
//   - error: 生成失败时返回错误
func Generate(cmd types.Command, shellType string) (string, error) {
	// 缓冲区
	var buf bytes.Buffer

	// 程序名称
	programName := filepath.Base(os.Args[0])

	// 获取根命令的补全选项
	rootCmdOpts := collectCompletionOptions(cmd)

	// 构建命令树条目缓冲区
	var cmdTreeEntries bytes.Buffer

	// 从根命令的子命令开始添加条目
	traverseCommandTree(&cmdTreeEntries, "", cmd.SubCmds(), shellType)

	// 缓存标志参数, 避免重复计算
	params := collectFlagParameters(cmd)

	// 根据shell类型调用对应的处理函数
	switch shellType {
	case types.BashShell: // Bash特定处理
		generateBashCompletion(&buf, params, rootCmdOpts, cmdTreeEntries.String(), programName)

	case types.PwshShell, types.PowershellShell: // PowerShell特定处理
		generatePwshCompletion(&buf, params, rootCmdOpts, cmdTreeEntries.String(), programName)

	default:
		return "", types.NewError("UNSUPPORTED_SHELL", "unsupported shell", nil)
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
func traverseCommandTree(buf *bytes.Buffer, rootPath string, cmds []types.Command, shellType string) {
	// node 代表一个即将写入的"树条目"
	type node struct {
		name       string        // 真正写入的名字 (长名或短名)
		parentPath string        // 父路径
		cmd        types.Command // 命令接口
	}

	if len(cmds) == 0 {
		return
	}

	// 初始化队列: 把第一层命令按长短名拆开入队
	q := list.New()
	for _, c := range cmds {
		if c == nil {
			continue
		}

		// 长名上下文入队
		if n := c.LongName(); n != "" {
			q.PushBack(node{name: n, parentPath: rootPath, cmd: c})
		}

		// 短名上下文入队
		if n := c.ShortName(); n != "" {
			q.PushBack(node{name: n, parentPath: rootPath, cmd: c})
		}
	}

	// 记录已经写了多少条, 最后一条不补逗号
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
		opts := collectCompletionOptions(cur.cmd)
		switch shellType {
		case types.BashShell: // Bash特定处理
			programName := filepath.Base(os.Args[0])
			generateBashCommandTreeEntry(buf, fullPath, opts, programName)

		case types.PwshShell, types.PowershellShell: // Powershell特定处理
			generatePwshCommandTreeEntry(buf, fullPath, opts)
		}

		// 子节点入队 (同样按长短名拆分)
		for _, sub := range cur.cmd.SubCmds() {
			if sub == nil {
				continue
			}

			// 长名入队
			if n := sub.LongName(); n != "" {
				q.PushBack(node{name: n, parentPath: fullPath, cmd: sub})
			}

			// 短名入队
			if n := sub.ShortName(); n != "" {
				q.PushBack(node{name: n, parentPath: fullPath, cmd: sub})
			}
		}

		// 计数器累加
		written++

		// 只要不是最后一条就补逗号
		if q.Len() > 0 {
			// 如果是pwsh或者powershell, 才处理
			if shellType == types.PwshShell || shellType == types.PowershellShell {
				buf.WriteString(",\n")
			}
		}
	}
}

// collectFlagParameters 使用广度优先搜索收集整个命令树的所有标志参数信息
//
// 算法说明:
// 1. 使用BFS遍历整个命令树, 确保所有层级的标志都被收集
// 2. 为每个标志生成唯一键(路径+标志名), 避免重复收集
// 3. 同时收集长短名标志, 因为shell补全需要支持两种形式
// 4. 特殊处理枚举类型标志, 提取可选值列表用于补全
//
// 参数:
//   - cmd: 根命令上下文, 作为遍历起点
//
// 返回值:
//   - []FlagParam: 包含所有标志参数信息的切片
func collectFlagParameters(cmd types.Command) []FlagParam {
	// 预分配切片容量, 基于常见CLI工具的标志数量估算
	params := make([]FlagParam, 0, types.DefaultFlagParamsCapacity)

	// 使用map进行去重, 键为"路径+标志名"的组合
	seen := make(map[string]struct{}, types.DefaultFlagParamsCapacity)

	// 队列 BFS
	type node struct {
		cmd  types.Command
		path string // 已标准化好的路径
	}
	q := []node{{cmd: cmd, path: "/"}}

	// add 添加一个标志到结果集
	add := func(prefix, name string, cur node, flag types.Flag) {
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
		ft := flag.Type()
		param := FlagParam{
			CommandPath: cur.path,
			Name:        prefix + name,
			Type:        getParamTypeByFlagType(ft),
			ValueType:   getValueTypeByFlagType(ft),
		}

		// 如果是枚举标志, 则获取枚举选项
		if ft == types.FlagTypeEnum {
			param.EnumOptions = flag.EnumValues()
		}
		params = append(params, param)
	}

	// 遍历命令节点
	for len(q) > 0 {
		cur := q[0] // 获取当前节点
		q = q[1:]   // 移除当前节点

		// 遍历当前命令的所有标志
		for _, flag := range cur.cmd.Flags() {
			// 如果短标志不为空, 则添加短标志
			if flag.ShortName() != "" {
				add("-", flag.ShortName(), cur, flag)
			}

			// 如果长标志不为空, 则添加长标志
			if flag.LongName() != "" {
				add("--", flag.LongName(), cur, flag)
			}
		}

		// 子命令入队
		for _, sub := range cur.cmd.SubCmds() {
			if sub == nil {
				continue
			}

			// 子命令路径
			subPath := path.Join(cur.path, sub.Name()) + "/"

			// 入队
			q = append(q, node{cmd: sub, path: subPath})
		}
	}
	return params
}

// collectCompletionOptions 收集命令的补全选项, 包括标志和子命令
//
// 参数:
//   - cmd: 命令接口
//
// 返回值:
//   - 包含所有标志选项和子命令名称的字符串切片
func collectCompletionOptions(cmd types.Command) []string {
	if cmd == nil {
		return nil
	}

	// 获取所有标志数量(长标志+短标志)
	flags := cmd.Flags()
	flagCnt := len(flags)

	// 计算总容量 (标志数量+子命令数量×每项名称数)
	// 每个子命令都有长名和短名, 所以乘以NamesPerItem
	capacity := flagCnt + len(cmd.SubCmds())*types.NamesPerItem

	// 创建一个用于存储标志选项和子命令名称的切片
	seen := make(map[string]struct{}, capacity)

	// 定义一个添加选项的函数
	add := func(s string) {
		if s != "" {
			seen[s] = struct{}{}
		}
	}

	// 1. flags  (同时展开长短名)
	for _, flag := range flags {
		if flag == nil {
			continue
		}

		if flag.LongName() != "" {
			add("--" + flag.LongName())
		}

		if flag.ShortName() != "" {
			add("-" + flag.ShortName())
		}
	}

	// 2. sub-commands (同时展开长短名)
	for _, sub := range cmd.SubCmds() {
		if sub == nil {
			continue
		}
		add(sub.LongName())
		add(sub.ShortName())
	}

	// 3. 转回切片
	opts := make([]string, 0, len(seen))
	for k := range seen {
		opts = append(opts, k)
	}
	return opts
}

// 通用字符串构建器对象池 - 供bash和pwsh共同使用
var stringBuilderPool = sync.Pool{
	New: func() interface{} {
		builder := &strings.Builder{}
		builder.Grow(512) // 预分配512字节容量
		return builder
	},
}

// buildString 使用对象池构建字符串的通用辅助函数
// 供bash_completion.go和pwsh_completion.go共同使用
func buildString(fn func(*strings.Builder)) string {
	builder := stringBuilderPool.Get().(*strings.Builder)
	defer func() {
		// 如果容量过大则不回收, 避免内存浪费
		if builder.Cap() <= 8192 {
			builder.Reset()
			stringBuilderPool.Put(builder)
		}
	}()

	fn(builder)
	return builder.String()
}

// getValueTypeByFlagType 根据标志类型获取值类型
//
// 参数:
//   - flagType - 标志类型
//
// 返回值:
//   - string: 值类型
func getValueTypeByFlagType(flagType types.FlagType) string {
	switch flagType {
	case types.FlagTypeBool:
		return "bool"
	case types.FlagTypeEnum:
		return "enum"
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
func getParamTypeByFlagType(flagType types.FlagType) string {
	if flagType == types.FlagTypeBool {
		return "none"
	}
	return "required"
}
