// Package completion Bash Shell 自动补全实现
// 本文件实现了Bash Shell环境下的命令行自动补全功能，
// 生成Bash补全脚本，支持标志和子命令的智能补全。
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
// - 高性能模糊匹配算法 (纯整数运算，无bc依赖)
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
		// 使用对象池构建键字符串，避免fmt.Sprintf的内存分配
		if param.CommandPath == "" {
			key = param.Name
		} else {
			key = buildString(func(builder *strings.Builder) {
				builder.WriteString(param.CommandPath)
				builder.WriteByte('|')
				builder.WriteString(param.Name)
			})
		}

		// 使用对象池构建参数值字符串，避免字符串拼接的内存分配
		paramValue := buildString(func(builder *strings.Builder) {
			builder.WriteString(param.Type)
			builder.WriteByte('|')
			builder.WriteString(param.ValueType)
		})

		// 使用对象池构建完整的标志参数项，避免fmt.Fprintf的格式化开销
		flagParamItem := buildString(func(builder *strings.Builder) {
			builder.WriteString(programName)
			builder.WriteString("_flag_params[\"")
			builder.WriteString(key)
			builder.WriteString("\"]=\"")
			builder.WriteString(paramValue)
			builder.WriteString("\"\n")
		})
		flagParamsBuf.WriteString(flagParamItem)

		// 如果参数类型为枚举，则生成枚举选项
		if param.ValueType == "enum" && len(param.EnumOptions) > 0 {
			// 预分配切片容量以提高性能
			escapedOpts := make([]string, len(param.EnumOptions))
			for i, opt := range param.EnumOptions {
				escapedOpts[i] = escapeSpecialChars(opt)
			}

			// 将枚举选项转换为字符串，使用|分隔符与其他选项保持一致
			options := strings.Join(escapedOpts, "|")

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
	_, _ = tmpl.WriteString(buf, BashFunctionHeader)
}

// bashEscapeMap Bash特殊字符转义映射表
// 使用全局map提高转义性能，避免重复的switch判断
var bashEscapeMap = map[rune]string{
	'\\': "\\\\", // 反斜杠
	'"':  "\\\"", // 双引号
	' ':  "\\ ",  // 空格
	'$':  "\\$",  // 美元符号
	'`':  "\\`",  // 反引号
	'|':  "\\|",  // 管道符
	'&':  "\\&",  // 与符号
	';':  "\\;",  // 分号
	'(':  "\\(",  // 左括号
	')':  "\\)",  // 右括号
	'<':  "\\<",  // 小于号
	'>':  "\\>",  // 大于号
	'*':  "\\*",  // 星号
	'?':  "\\?",  // 问号
	'[':  "\\[",  // 左方括号
	']':  "\\]",  // 右方括号
	'{':  "\\{",  // 左花括号
	'}':  "\\}",  // 右花括号
	'~':  "\\~",  // 波浪号
	'#':  "\\#",  // 井号
}

// escapeSpecialChars 处理字符串中的特殊字符转义
// 优化版本：使用全局map进行O(1)查找，提升性能
// 用于确保生成的bash脚本中的字符串安全性
//
// 参数:
//   - s: 需要处理的字符串
//
// 返回值:
//   - 转义后的字符串，可安全用于bash脚本中
func escapeSpecialChars(s string) string {
	var builder strings.Builder
	builder.Grow(len(s) * 2) // 预分配容量以减少重新分配

	for _, r := range s {
		if escaped, exists := bashEscapeMap[r]; exists {
			builder.WriteString(escaped)
		} else {
			builder.WriteRune(r)
		}
	}

	return builder.String()
}

const (
	BashFlagParamItem = "{{.ProgramName}}_flag_params[%q]=%q\n"  // 标志参数项格式
	BashEnumOptions   = "{{.ProgramName}}_enum_options[%q]=%q\n" // 枚举选项格式
)

const (
	// 优化的Bash补全模板 - 集成高性能模糊匹配功能
	BashFunctionHeader = `#!/usr/bin/env bash

# ==================== 模糊补全配置参数 ====================
# 模糊补全功能开关 (设置为0禁用，1启用)
readonly {{.ProgramName}}_FUZZY_COMPLETION_ENABLED=1

# 启用模糊补全的最大候选项数量阈值
# 超过此数量将回退到传统前缀匹配以保证性能
readonly {{.ProgramName}}_FUZZY_MAX_CANDIDATES=150

# 模糊匹配的最小输入长度 (小于此长度不启用模糊匹配)
readonly {{.ProgramName}}_FUZZY_MIN_PATTERN_LENGTH=2

# 模糊匹配分数阈值 (0-100，分数低于此值的匹配将被过滤)
readonly {{.ProgramName}}_FUZZY_SCORE_THRESHOLD=30

# 模糊匹配最大返回结果数
readonly {{.ProgramName}}_FUZZY_MAX_RESULTS=8

# 缓存大小控制参数
# 缓存条目数量超过此阈值时将清空缓存以防止内存无限增长
readonly {{.ProgramName}}_FUZZY_CACHE_MAX_SIZE=500

# ==================== 静态数据定义 ====================
# 静态命令树定义 - 在函数外预初始化
declare -A {{.ProgramName}}_cmd_tree
{{.ProgramName}}_cmd_tree["/"]="{{.RootCmdOpts}}"
{{.CmdTreeEntries}}

# 标志参数定义 - 存储类型和值类型 (type|valueType)
declare -A {{.ProgramName}}_flag_params
{{.FlagParams}}

# 枚举选项定义 - 存储枚举标志的允许值
declare -A {{.ProgramName}}_enum_options
{{.EnumOptions}}

# 模糊匹配结果缓存 (格式: "pattern|candidate" -> score)
declare -A {{.ProgramName}}_fuzzy_cache

# ==================== 模糊匹配核心算法 ====================
# 高性能模糊评分函数 - 使用纯整数运算避免bc开销
# 参数: $1=输入模式, $2=候选字符串
# 返回: 0-100的整数分数 (通过echo输出)
_{{.ProgramName}}_fuzzy_score_fast() {
    local pattern="$1"
    local candidate="$2"
    local pattern_len=${#pattern}
    local candidate_len=${#candidate}
    
    # 性能优化1: 长度预检查 - 候选项太短直接返回0
    if [[ $candidate_len -lt $pattern_len ]]; then
        echo "0"
        return
    fi
    
    # 性能优化2: 完全匹配检查 - 避免不必要的复杂计算
    if [[ "$candidate" == "$pattern"* ]]; then
        echo "100"  # 前缀完全匹配给最高分
        return
    fi
    
    # 性能优化3: 字符存在性预检查 - 快速排除不可能的匹配
    local pattern_lower="${pattern,,}"  # 转小写用于大小写不敏感匹配
    local candidate_lower="${candidate,,}"
    local i
    for ((i=0; i<pattern_len; i++)); do
        local char="${pattern_lower:$i:1}"
        if [[ "$candidate_lower" != *"$char"* ]]; then
            echo "0"  # 必需字符不存在，直接返回
            return
        fi
    done
    
    # 核心匹配算法 - 计算字符匹配度和连续性
    local matched=0           # 匹配的字符数
    local consecutive=0       # 当前连续匹配长度
    local max_consecutive=0   # 最大连续匹配长度
    local candidate_pos=0     # 候选字符串当前搜索位置
    local start_bonus=0       # 起始位置奖励
    
    # 检查是否从开头匹配 (大小写不敏感)
    if [[ "$candidate_lower" == "$pattern_lower"* ]]; then
        start_bonus=20  # 起始匹配给20分奖励
    fi
    
    # 逐字符匹配算法
    for ((i=0; i<pattern_len; i++)); do
        local pattern_char="${pattern_lower:$i:1}"
        local found=0
        
        # 在候选字符串中查找当前模式字符
        for ((j=candidate_pos; j<candidate_len; j++)); do
            local candidate_char="${candidate_lower:$j:1}"
            if [[ "$pattern_char" == "$candidate_char" ]]; then
                ((matched++))
                found=1
                
                # 连续性检查 - 连续匹配的字符得分更高
                if [[ $j -eq $candidate_pos ]]; then
                    ((consecutive++))
                    if [[ $consecutive -gt $max_consecutive ]]; then
                        max_consecutive=$consecutive
                    fi
                else
                    consecutive=1  # 重置连续计数
                fi
                
                candidate_pos=$((j+1))  # 更新搜索位置
                break
            fi
        done
        
        # 如果某个字符未找到，重置连续计数
        if [[ $found -eq 0 ]]; then
            consecutive=0
        fi
    done
    
    # 评分计算 - 使用整数运算避免浮点数和bc
    # 基础分数: (匹配字符数 / 模式长度) * 60
    local base_score=$((matched * 60 / pattern_len))
    
    # 连续性奖励: (最大连续长度 / 模式长度) * 20
    local consecutive_bonus=$((max_consecutive * 20 / pattern_len))
    
    # 长度惩罚: 候选字符串越长，分数略微降低
    local length_penalty=$((candidate_len - pattern_len))
    if [[ $length_penalty -gt 10 ]]; then
        length_penalty=10  # 最大惩罚10分
    fi
    
    # 最终分数计算
    local final_score=$((base_score + consecutive_bonus + start_bonus - length_penalty))
    
    # 确保分数在0-100范围内
    if [[ $final_score -lt 0 ]]; then
        final_score=0
    elif [[ $final_score -gt 100 ]]; then
        final_score=100
    fi
    
    echo "$final_score"
}

# 带缓存的模糊评分函数 - 避免重复计算提高性能
# 参数: $1=输入模式, $2=候选字符串
_{{.ProgramName}}_fuzzy_score_cached() {
    local pattern="$1"
    local candidate="$2"
    local cache_key="${pattern}|${candidate}"
    
    # 缓存命中检查
    if [[ -n "${{{.ProgramName}}_fuzzy_cache[$cache_key]}" ]]; then
        echo "${{{.ProgramName}}_fuzzy_cache[$cache_key]}"
        return
    fi
    
    # 计算分数并缓存
    local score
    score=$(_{{.ProgramName}}_fuzzy_score_fast "$pattern" "$candidate")
    
    # 缓存大小控制 - 防止内存无限增长
    if [[ ${#{{.ProgramName}}_fuzzy_cache[@]} -gt ${{.ProgramName}}_FUZZY_CACHE_MAX_SIZE ]]; then
        {{.ProgramName}}_fuzzy_cache=()  # 清空缓存
    fi
    
    {{.ProgramName}}_fuzzy_cache["$cache_key"]="$score"
    echo "$score"
}

# 智能补全匹配函数 - 分级匹配策略
# 参数: $1=输入模式, $2=候选选项字符串(用|分隔)
_{{.ProgramName}}_intelligent_match() {
    local pattern="$1"
    local options_str="$2"
    local pattern_len=${#pattern}
    
    # 解析候选选项到数组
    local opts_arr
    IFS='|' read -ra opts_arr <<< "$options_str"
    local total_candidates=${#opts_arr[@]}
    
    # 性能保护: 候选项过多时禁用模糊匹配
    if [[ $total_candidates -gt ${{.ProgramName}}_FUZZY_MAX_CANDIDATES ]]; then
        # 回退到传统compgen前缀匹配
        local opts
        printf -v opts '%s ' "${opts_arr[@]}"
        opts="${opts% }"
        COMPREPLY=($(compgen -W "$opts" -- "$pattern"))
        return 0
    fi
    
    local matches=()
    
    # 第1级: 精确前缀匹配 (最快，优先级最高)
    local exact_matches=()
    local opt
    for opt in "${opts_arr[@]}"; do
        if [[ "$opt" == "$pattern"* ]]; then
            exact_matches+=("$opt")
        fi
    done
    
    # 如果有精确匹配且数量合理，直接返回
    if [[ ${#exact_matches[@]} -gt 0 && ${#exact_matches[@]} -le 15 ]]; then
        COMPREPLY=("${exact_matches[@]}")
        return 0
    fi
    
    # 第2级: 大小写不敏感前缀匹配
    if [[ ${#exact_matches[@]} -eq 0 ]]; then
        local pattern_lower="${pattern,,}"
        local case_insensitive_matches=()
        
        for opt in "${opts_arr[@]}"; do
            local opt_lower="${opt,,}"
            if [[ "$opt_lower" == "$pattern_lower"* ]]; then
                case_insensitive_matches+=("$opt")
            fi
        done
        
        # 如果有大小写不敏感匹配，返回
        if [[ ${#case_insensitive_matches[@]} -gt 0 ]]; then
            COMPREPLY=("${case_insensitive_matches[@]}")
            return 0
        fi
    fi
    
    # 第3级: 模糊匹配 (最慢，仅在必要时使用)
    if [[ ${{.ProgramName}}_FUZZY_COMPLETION_ENABLED -eq 1 && $pattern_len -ge ${{.ProgramName}}_FUZZY_MIN_PATTERN_LENGTH ]]; then
        local fuzzy_results=()
        local scored_matches=()
        
        # 对所有候选项进行模糊评分
        for opt in "${opts_arr[@]}"; do
            local score
            score=$(_{{.ProgramName}}_fuzzy_score_cached "$pattern" "$opt")
            
            # 只保留分数达到阈值的匹配
            if [[ $score -ge ${{.ProgramName}}_FUZZY_SCORE_THRESHOLD ]]; then
                scored_matches+=("$score:$opt")
            fi
        done
        
        # 如果有模糊匹配结果，按分数排序并返回前N个
        if [[ ${#scored_matches[@]} -gt 0 ]]; then
            # 按分数降序排序
            local sorted_matches
            mapfile -t sorted_matches < <(printf '%s\n' "${scored_matches[@]}" | sort -t: -k1 -nr)
            
            # 提取选项名称，限制返回数量
            local count=0
            local match
            for match in "${sorted_matches[@]}"; do
                if [[ $count -ge ${{.ProgramName}}_FUZZY_MAX_RESULTS ]]; then
                    break
                fi
                fuzzy_results+=("${match#*:}")  # 移除分数前缀
                ((count++))
            done
            
            COMPREPLY=("${fuzzy_results[@]}")
            return 0
        fi
    fi
    
    # 第4级: 子字符串匹配 (最后的备选方案)
    local substring_matches=()
    local pattern_lower="${pattern,,}"
    
    for opt in "${opts_arr[@]}"; do
        local opt_lower="${opt,,}"
        if [[ "$opt_lower" == *"$pattern_lower"* ]]; then
            substring_matches+=("$opt")
        fi
    done
    
    if [[ ${#substring_matches[@]} -gt 0 ]]; then
        COMPREPLY=("${substring_matches[@]}")
        return 0
    fi
    
    # 如果所有匹配策略都失败，返回空结果
    COMPREPLY=()
    return 1
}

# ==================== 主补全函数 ====================
_{{.ProgramName}}() {
	local cur prev words cword context opts i arg
	COMPREPLY=()

	# 使用 _get_comp_words_by_ref 获取补全参数以提高健壮性
	if declare -F _get_comp_words_by_ref >/dev/null 2>&1; then
		_get_comp_words_by_ref -n =: cur prev words cword
	else
		# 与旧版本 Bash 补全环境的兼容性
		words=("${COMP_WORDS[@]}")
		cword=$COMP_CWORD
		cur="${words[cword]}"
		prev="${words[cword-1]}"
	fi

	# 输入验证
	if [[ $cword -lt 0 || ${#words[@]} -eq 0 ]]; then
		return 1
	fi

	# 查找当前命令上下文
	local context="/"
	local i
	for ((i=1; i < cword; i++)); do
		local arg="${words[i]}"
		# 如果遇到标志，停止上下文构建
		if [[ "$arg" =~ ^- ]]; then
			break
		fi
		
		local next_context="$context$arg/"
		# 验证上下文是否存在
		if [[ -n "${{{.ProgramName}}_cmd_tree[$next_context]}" ]]; then
			context="$next_context"
		else
			break
		fi
	done

	# 获取当前上下文的可用选项
	local current_context_opts="${{{.ProgramName}}_cmd_tree[$context]}"
	if [[ -z "$current_context_opts" ]]; then
		return 1
	fi
	
	# 检查前一个参数是否需要值并获取其类型
	local prev_param_type=""
	local prev_value_type=""
	if [[ $cword -gt 1 ]]; then
		local prev_arg="${words[cword-1]}"
		local key="${context}|${prev_arg}"
		local prev_param_info="${{{.ProgramName}}_flag_params[$key]:-}"
		
		if [[ -n "$prev_param_info" ]]; then
			IFS='|' read -r prev_param_type prev_value_type <<< "$prev_param_info"
		fi
	fi

	# 根据参数类型动态生成补全
	if [[ -n "$prev_param_type" && "$prev_param_type" == "required" ]]; then
		case "$prev_value_type" in
			enum)
				local enum_key="${context}|${words[cword-1]}"
				local enum_opts="${{{.ProgramName}}_enum_options[$enum_key]:-}"
				
				if [[ -n "$enum_opts" ]]; then
					# 对枚举选项也使用智能匹配
					_{{.ProgramName}}_intelligent_match "$cur" "$enum_opts"
					return 0
				fi
				;;
			string)
				# 字符串类型 - 提供文件和目录路径补全
				COMPREPLY=($(compgen -f -d -- "$cur"))
				return 0
				;;
			*)
				# 默认值补全 - 使用智能匹配
				_{{.ProgramName}}_intelligent_match "$cur" "$current_context_opts"
				return 0
				;;
		esac
	fi

	# 主要的标志和命令补全 - 使用智能匹配算法
	_{{.ProgramName}}_intelligent_match "$cur" "$current_context_opts"
	return 0
}

# ==================== 调试和诊断函数 ====================
# 补全系统健康检查函数 (可选，用于调试)
_{{.ProgramName}}_completion_debug() {
    echo "=== {{.ProgramName}} 补全系统诊断 ==="
    echo "Bash版本: $BASH_VERSION"
    echo "补全函数状态: $(type -t _{{.ProgramName}})"
    echo "命令树条目数: ${#{{.ProgramName}}_cmd_tree[@]}"
    echo "标志参数数: ${#{{.ProgramName}}_flag_params[@]}"
    echo "枚举选项数: ${#{{.ProgramName}}_enum_options[@]}"
    echo "模糊补全状态: $([ ${{.ProgramName}}_FUZZY_COMPLETION_ENABLED -eq 1 ] && echo "启用" || echo "禁用")"
    echo "候选项阈值: ${{.ProgramName}}_FUZZY_MAX_CANDIDATES"
    echo "缓存条目数: ${#{{.ProgramName}}_fuzzy_cache[@]}"
    echo ""
    echo "使用方法: 在命令行输入 '_{{.ProgramName}}_completion_debug' 查看此信息"
}

# 注册补全函数
complete -F _{{.ProgramName}} {{.ProgramName}}
`
)
