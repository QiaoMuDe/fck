// bash_completion.go - bash 补全脚本生成器
package cmd

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
func generateBashCommandTreeEntry(cmdTreeEntries *bytes.Buffer, cmdPath string, cmdOpts []string) {
	fmt.Fprintf(cmdTreeEntries, BashCommandTreeEntry, cmdPath, strings.Join(cmdOpts, "|"))
}

// generateBashCompletion 生成Bash自动补全脚本
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
	for _, param := range params {
		var key string
		if param.CommandPath == "" {
			key = param.Name
		} else {
			key = fmt.Sprintf("%s|%s", param.CommandPath, param.Name)
		}
		fmt.Fprintf(&flagParamsBuf, BashFlagParamItem, key, param.Type+"|"+param.ValueType)
		if param.ValueType == "enum" && len(param.EnumOptions) > 0 {
			// 使用bytes.Buffer减少内存分配
			var optionsBuf bytes.Buffer
			for i, opt := range param.EnumOptions {
				if i > 0 {
					optionsBuf.WriteString(" ")
				}
				// 增强特殊字符转义处理: 引号、反斜杠和空格
				escapedOpt := strings.ReplaceAll(opt, "\\", "\\\\")
				escapedOpt = strings.ReplaceAll(escapedOpt, "\"", "\\\"")
				escapedOpt = strings.ReplaceAll(escapedOpt, " ", "\\ ")
				optionsBuf.WriteString(escapedOpt)
			}
			options := optionsBuf.String()
			fmt.Fprintf(&enumOptionsBuf, BashEnumOptions, key, options)
		}
	}

	// 写入Bash自动补全脚本头
	fmt.Fprintf(buf, BashFunctionHeader, strings.Join(rootCmdOpts, "|"), cmdTreeEntries, flagParamsBuf.String(), enumOptionsBuf.String(), programName, programName, programName)
}

const (
	BashCommandTreeEntry = "cmd_tree[%s]=\"%s\"\n" // 命令树条目格式
	BashFlagParamItem    = "flag_params[%q]=%q\n"  // 标志参数项格式
	BashEnumOptions      = "enum_options[%q]=%q\n" // 枚举选项格式
)

const (
	// Bash补全模板
	BashFunctionHeader = `#!/usr/bin/env bash

# Static command tree definition - Pre-initialized outside the function
declare -A cmd_tree
cmd_tree[/]="%s"
%s

# Flag parameters definition - stores type and value type (type|valueType)
declare -A flag_params
%s

# Enum options definition - stores allowed values for enum flags
declare -A enum_options
%s

_%s() {
	local cur prev words cword context opts i arg
	COMPREPLY=()

	# Use _get_comp_words_by_ref to get completion parameters for better robustness
	if [[ -z "${_get_comp_words_by_ref}" ]]; then
		# Compatibility with older versions of Bash completion environment
		words=("${COMP_WORDS[@]}")
		cword=$COMP_CWORD
	else
		_get_comp_words_by_ref -n =: cur prev words cword
	fi

	cur="${words[cword]}"
	prev="${words[cword-1]}"

	# Find the current command context
	local context="/"
	local i
	for ((i=1; i < cword; i++)); do
		local arg="${words[i]}"
		if [[ -n "${cmd_tree[$context$arg/]}" ]]; then
			context="$context$arg/"
		fi
	done

	# Get the available options for the current context
	IFS='|' read -ra opts_arr <<< "${cmd_tree[$context]}"
	opts=$(IFS=' '; echo "${opts_arr[*]}")
	
	# Check if the previous parameter needs a value and get its type
	prev_param_type=""
	prev_value_type=""
	if [[ $cword -gt 1 ]]; then
		prev_arg="${words[cword-1]}"
		key="${context}|${prev_arg}"
		prev_param_info=${flag_params[$key]}
		IFS='|' read -r prev_param_type prev_value_type <<< "${prev_param_info}"
	fi

	# Dynamically generate completion based on parameter type
	if [[ -n "$prev_param_type" && $prev_param_type == "required" ]]; then
		case "$prev_value_type" in
			path)
				# Path type parameter, use file and directory completion
				COMPREPLY=($(compgen -f -d -- "${cur}"))
				return 0
				;;
			number)
				# Number type parameter, provide basic number completion
				COMPREPLY=($(compgen -W "$(seq 1 100)" -- "${cur}"))
				return 0
				;;
			ip)
				# IP address type parameter, provide basic IP completion
				COMPREPLY=($(compgen -W "192.168. 10.0. 172.16." -- "${cur}"))
				return 0
				;;
			enum)
			    # 当前单词为空且前一个参数是枚举标志 → 直接列出所有枚举值
                if [[ -z "$cur" && "$prev_value_type" == "enum" ]]; then
                    COMPREPLY=($(compgen -W "${enum_options[$key]}" -- ""))
                    return 0
                fi
				
				# 前缀过滤（大小写不敏感）
                COMPREPLY=($(compgen -W "${enum_options[$key]}" -- "${cur}"))
                # 只保留以 $cur(忽略大小写)开头的
                COMPREPLY=($(echo "${COMPREPLY[@]}" | grep -i "^${cur}"))
				return 0
				;;

			url)
				# URL type parameter, provide common URL prefix completion
				COMPREPLY=($(compgen -W "http:// https:// ftp://" -- "${cur}"))
				return 0
				;;
			*)
                # Default value completion
				COMPREPLY=($(compgen -W "${opts}" -- "${cur}"))
				return 0
				;;
			esac
	fi

	# Flag parameter completion
	COMPREPLY=($(compgen -W "${opts}" -- "${cur}"))

	return 0
	}

complete -F _%s %s
`
)
