// Package completion Bash Shell 自动补全实现
// 本文件实现了Bash Shell环境下的命令行自动补全功能,
// 生成Bash补全脚本, 支持标志和子命令的智能补全。
package completion

import (
	"bytes"
	"fmt"
	"strings"
)

// generateBashCommandTreeEntry 生成Bash命令树条目
//
// 参数:
// - cmdTreeEntries: 命令树条目缓冲区
// - cmdPath: 命令路径
// - cmdOpts: 命令选项
// - programName: 程序名称
func generateBashCommandTreeEntry(cmdTreeEntries *bytes.Buffer, cmdPath string, cmdOpts []string, programName string) {
	fmt.Fprintf(cmdTreeEntries, "%s_cmd_tree[%q]=%q\n", programName, cmdPath, strings.Join(cmdOpts, "|"))
}

// generateBashCompletion 生成优化的Bash自动补全脚本
//
// 新增功能:
// - 高性能模糊匹配算法 (纯整数运算, 无bc依赖)
// - 智能分级匹配策略 (精确->大小写不敏感->模糊->子字符串)
// - 性能保护机制 (候选项过多时自动回退到传统匹配)
// - 结果缓存系统 (避免重复计算)
// - 可配置的性能参数
//
// 参数:
// - buf: 输出缓冲区
// - params: 标志参数列表
// - rootCmdOpts: 根命令选项
// - cmdTreeEntries: 命令树条目
// - programName: 程序名称
func generateBashCompletion(buf *bytes.Buffer, params []FlagParam, rootCmdOpts []string, cmdTreeEntries string, programName string) {
	// 构建标志参数映射
	var flagParamsBuf bytes.Buffer
	var enumOptionsBuf bytes.Buffer

	// 遍历标志参数并生成相应的Bash自动补全脚本
	for _, param := range params {
		var key string
		// 使用对象池构建键字符串, 避免fmt.Sprintf的内存分配
		if param.CommandPath == "" {
			key = param.Name
		} else {
			key = buildString(func(builder *strings.Builder) {
				builder.WriteString(param.CommandPath)
				builder.WriteByte('|')
				builder.WriteString(param.Name)
			})
		}

		// 使用对象池构建参数值字符串, 避免字符串拼接的内存分配
		paramValue := buildString(func(builder *strings.Builder) {
			builder.WriteString(param.Type)
			builder.WriteByte('|')
			builder.WriteString(param.ValueType)
		})

		// 使用对象池构建完整的标志参数项, 避免fmt.Fprintf的格式化开销
		flagParamItem := buildString(func(builder *strings.Builder) {
			builder.WriteString(programName)
			builder.WriteString("_flag_params[\"")
			builder.WriteString(key)
			builder.WriteString("\"]=\"")
			builder.WriteString(paramValue)
			builder.WriteString("\"\n")
		})
		flagParamsBuf.WriteString(flagParamItem)

		// 如果参数类型为枚举, 则生成枚举选项
		if param.ValueType == "enum" && len(param.EnumOptions) > 0 {
			// 将枚举选项转换为字符串, 使用|分隔符与其他选项保持一致
			options := strings.Join(param.EnumOptions, "|")

			// 写入枚举选项
			enumOptionItem := buildString(func(builder *strings.Builder) {
				builder.WriteString(programName)
				builder.WriteString("_enum_options[\"")
				builder.WriteString(key)
				builder.WriteString("\"]=\"")
				builder.WriteString(options)
				builder.WriteString("\"\n")
			})
			enumOptionsBuf.WriteString(enumOptionItem)
		}
	}

	// 使用命名模板生成Bash自动补全脚本
	tmpl := strings.NewReplacer(
		"{{.RootCmdOpts}}", strings.Join(rootCmdOpts, "|"), // 根命令选项
		"{{.CmdTreeEntries}}", cmdTreeEntries, // 命令树条目
		"{{.FlagParams}}", flagParamsBuf.String(), // 标志参数
		"{{.EnumOptions}}", enumOptionsBuf.String(), // 枚举选项
		"{{.ProgramName}}", programName, // 程序名称
	)

	// 写入Bash函数头部
	_, _ = tmpl.WriteString(buf, bashTemplate)
}

const (
	BashFlagParamItem = "{{.ProgramName}}_flag_params[%q]=%q\n"  // 标志参数项格式
	BashEnumOptions   = "{{.ProgramName}}_enum_options[%q]=%q\n" // 枚举选项格式
)
