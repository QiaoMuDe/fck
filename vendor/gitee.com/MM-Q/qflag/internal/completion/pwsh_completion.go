// Package completion PowerShell 自动补全实现
// 本文件实现了PowerShell环境下的命令行自动补全功能,
// 生成PowerShell补全脚本, 支持标志和子命令的智能补全。
package completion

import (
	"bytes"
	"path/filepath"
	"strings"
)

// formatOptions 将选项列表格式化为PowerShell数组字符串
//
// 参数:
// - buf: 输出缓冲区
// - options: 选项列表
func formatOptions(buf *bytes.Buffer, options []string) {
	for i, opt := range options {
		// 只有不为空的选项才添加到缓冲区
		if opt == "" {
			continue
		}

		// 如果不是第一个选项, 则添加逗号
		if i > 0 {
			buf.WriteString(", ")
		}

		// 添加选项
		buf.WriteString("'" + opt + "'")
	}
}

// generatePwshCommandTreeEntry 生成PowerShell命令树条目
// 使用对象池优化内存分配, 避免创建临时缓冲区和Replacer
//
// 参数:
// - cmdTreeEntries: 命令树条目缓冲区
// - cmdPath: 命令路径
// - cmdOpts: 命令选项
func generatePwshCommandTreeEntry(cmdTreeEntries *bytes.Buffer, cmdPath string, cmdOpts []string) {
	// 使用对象池构建命令树条目, 避免创建临时缓冲区和strings.NewReplacer的开销
	cmdTreeItem := buildString(func(builder *strings.Builder) {
		builder.WriteString("\t@{ Context = \"")
		builder.WriteString(cmdPath)
		builder.WriteString("\"; Options = @(")

		// 直接在builder中格式化选项, 避免额外的字符串分配
		first := true
		for _, opt := range cmdOpts {
			if opt == "" {
				continue
			}

			if !first {
				builder.WriteString(", ")
			}
			first = false

			builder.WriteString("'" + opt + "'")
		}

		builder.WriteString(") }")
	})

	cmdTreeEntries.WriteString(cmdTreeItem)
}

// generatePwshCompletion 生成PowerShell自动补全脚本
//
// 参数:
// - buf: 输出缓冲区
// - params: 标志参数列表
// - rootCmdOpts: 根命令选项
// - cmdTreeEntries: 命令树条目
// - programName: 程序名称
func generatePwshCompletion(buf *bytes.Buffer, params []FlagParam, rootCmdOpts []string, cmdTreeEntries string, programName string) {
	// 构建标志参数和枚举选项
	flagParamsBuf := bytes.NewBuffer(make([]byte, 0, len(params)*100)) // 预分配容量

	// 处理根命令选项
	rootOptsBuf := bytes.NewBuffer(make([]byte, 0, len(rootCmdOpts)*20))
	formatOptions(rootOptsBuf, rootCmdOpts)

	// 处理标志参数
	for i, param := range params {
		// 生成带枚举选项的标志参数条目
		enumOptions := ""
		if param.ValueType == "enum" && len(param.EnumOptions) > 0 {
			optionsBuf := bytes.NewBuffer(make([]byte, 0, len(param.EnumOptions)*15))
			formatOptions(optionsBuf, param.EnumOptions)
			enumOptions = optionsBuf.String()
		}

		// 使用命名占位符替换位置参数
		flagReplacer := strings.NewReplacer(
			"{{.Context}}", param.CommandPath,
			"{{.Parameter}}", param.Name,
			"{{.ParamType}}", param.Type,
			"{{.ValueType}}", param.ValueType,
			"{{.Options}}", enumOptions,
		)
		flagParamsBuf.WriteString(flagReplacer.Replace(PwshFlagParamItem))

		// 条目之间添加逗号, 非最后一个条目
		if i < len(params)-1 {
			flagParamsBuf.WriteString(",\n")
		}
	}

	// 清理程序名, 去除可能的后缀
	sanitizedProgramName := strings.TrimSuffix(programName, filepath.Ext(programName))

	// 生成根命令条目
	rootReplacer := strings.NewReplacer(
		"{{.Context}}", "/",
		"{{.Options}}", rootOptsBuf.String(),
	)
	// 生成根命令条目
	rootCmdEntry := rootReplacer.Replace(PwshCmdTreeItem)

	// 如果命令树条目不为空, 则添加逗号
	if cmdTreeEntries != "" {
		rootCmdEntry += ",\n" + cmdTreeEntries
	}

	// 使用命名占位符替换位置参数
	completionReplacer := strings.NewReplacer(
		"{{.SanitizedName}}", sanitizedProgramName, // 替换程序名称
		"{{.ProgramName}}", programName, // 替换程序名称
		"{{.CmdTree}}", rootCmdEntry, // 替换命令树条目
		"{{.FlagParams}}", flagParamsBuf.String(), // 替换标志参数
	)

	// 写入PowerShell自动补全脚本
	_, _ = buf.WriteString(completionReplacer.Replace(pwshTemplate))
}

const (
	// 标志参数条目(含枚举选项)
	PwshFlagParamItem = "	@{ Context = \"{{.Context}}\"; Parameter = \"{{.Parameter}}\"; ParamType = \"{{.ParamType}}\"; ValueType = \"{{.ValueType}}\"; Options = @({{.Options}}) }"
	// 命令树条目
	PwshCmdTreeItem = "	@{ Context = \"{{.Context}}\"; Options = @({{.Options}}) }"
)
