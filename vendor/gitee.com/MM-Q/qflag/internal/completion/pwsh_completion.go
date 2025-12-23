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
// - escape: 字符串转义函数
func formatOptions(buf *bytes.Buffer, options []string, escape func(string) string) {
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
		buf.WriteByte('\'')
		buf.WriteString(escape(opt))
		buf.WriteByte('\'')
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

			builder.WriteByte('\'')
			builder.WriteString(escapePwshString(opt))
			builder.WriteByte('\'')
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
	formatOptions(rootOptsBuf, rootCmdOpts, escapePwshString)

	// 处理标志参数
	for i, param := range params {
		// 生成带枚举选项的标志参数条目
		enumOptions := ""
		if param.ValueType == "enum" && len(param.EnumOptions) > 0 {
			optionsBuf := bytes.NewBuffer(make([]byte, 0, len(param.EnumOptions)*15))
			formatOptions(optionsBuf, param.EnumOptions, escapePwshString)
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
	_, _ = buf.WriteString(completionReplacer.Replace(PwshFunctionHeader))
}

// pwshEscapeMap PowerShell特殊字符转义映射表
// 使用全局map提高转义性能, 避免重复的switch判断
var pwshEscapeMap = map[byte][]byte{
	'\'': {'\'', '\''}, // 单引号转义为两个单引号
	'\\': {'\\', '\\'}, // 反斜杠转义为两个反斜杠
	'$':  {'`', '$'},   // 美元符号转义
	'`':  {'`', '`'},   // 反引号转义
	'"':  {'`', '"'},   // 双引号转义
	'&':  {'`', '&'},   // 与符号转义
	'|':  {'`', '|'},   // 管道符转义
	';':  {'`', ';'},   // 分号转义
	'<':  {'`', '<'},   // 小于号转义
	'>':  {'`', '>'},   // 大于号转义
	'(':  {'`', '('},   // 左括号转义
	')':  {'`', ')'},   // 右括号转义
	'\r': {'`', 'r'},   // 回车符转义
	'\n': {'`', 'n'},   // 换行符转义
	'\t': {'`', 't'},   // 制表符转义
}

// escapePwshString 转义PowerShell字符串中的特殊字符
// 优化版本：使用全局map进行O(1)查找, 提升性能
//
// 参数:
// - s: 需要转义的字符串
//
// 返回:
// - 转义后的字符串
func escapePwshString(s string) string {
	// 预计算所需容量：最坏情况下每个字符都需要转义
	buf := make([]byte, 0, len(s)*2)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if escaped, exists := pwshEscapeMap[c]; exists {
			buf = append(buf, escaped...)
		} else {
			buf = append(buf, c)
		}
	}
	return string(buf)
}

const (
	// 标志参数条目(含枚举选项)
	PwshFlagParamItem = "	@{ Context = \"{{.Context}}\"; Parameter = \"{{.Parameter}}\"; ParamType = \"{{.ParamType}}\"; ValueType = \"{{.ValueType}}\"; Options = @({{.Options}}) }"
	// 命令树条目
	PwshCmdTreeItem = "	@{ Context = \"{{.Context}}\"; Options = @({{.Options}}) }"
)

const (
	// PowerShell自动补全脚本头部
	PwshFunctionHeader = `# -------------------------- Configuration Area (Need to be modified according to actual commands) --------------------------
# 命令名称
${{.SanitizedName}}_commandName = "{{.ProgramName}}"

# 1. 命令树结构
${{.SanitizedName}}_cmdTree = @(
{{.CmdTree}}
)

# 2. 标志参数定义
${{.SanitizedName}}_flagParams = @(
{{.FlagParams}}
)

# -----------------------------------------------------------------------------------

# ==================== 模糊补全配置参数 ====================
# 模糊补全功能开关 (设置为$false禁用, $true启用)
$script:{{.SanitizedName}}_FUZZY_COMPLETION_ENABLED = $true

# 启用模糊补全的最大候选项数量阈值
# 超过此数量将回退到传统前缀匹配以保证性能
$script:{{.SanitizedName}}_FUZZY_MAX_CANDIDATES = 120

# 模糊匹配的最小输入长度 (小于此长度不启用模糊匹配)
$script:{{.SanitizedName}}_FUZZY_MIN_PATTERN_LENGTH = 2

# 模糊匹配分数阈值 (0-100, 分数低于此值的匹配将被过滤)
$script:{{.SanitizedName}}_FUZZY_SCORE_THRESHOLD = 25

# 模糊匹配最大返回结果数
$script:{{.SanitizedName}}_FUZZY_MAX_RESULTS = 10

# 缓存大小控制参数
# 缓存条目数量超过此阈值时将清空缓存以防止内存无限增长
$script:{{.SanitizedName}}_FUZZY_CACHE_MAX_SIZE = 500

# 模糊匹配结果缓存 (格式: "pattern|candidate" -> score)
$script:{{.SanitizedName}}_fuzzyCache = @{}

# ==================== 模糊匹配核心算法 ====================

# 高性能模糊评分函数 - 使用优化的字符串操作
# 参数: $Pattern=输入模式, $Candidate=候选字符串
# 返回: 0-100的整数分数
function Get-{{.SanitizedName}}FuzzyScoreFast {
    param(
        [string]$Pattern,
        [string]$Candidate
    )
    
    $patternLen = $Pattern.Length
    $candidateLen = $Candidate.Length
    
    # 快速路径1: 空模式检查
    if ($patternLen -eq 0) {
        return 100
    }
    
    # 性能优化1: 长度预检查 - 候选项太短直接返回0
    if ($candidateLen -lt $patternLen) {
        return 0
    }
    
    # 性能优化2: 完全匹配检查 - 避免不必要的复杂计算
    if ($Candidate.StartsWith($Pattern, [System.StringComparison]::OrdinalIgnoreCase)) {
        return 100  # 前缀完全匹配给最高分
    }
    
    # 内存访问优化: 预转换字符数组，避免重复字符串索引访问
    $patternLower = $Pattern.ToLowerInvariant()
    $candidateLower = $Candidate.ToLowerInvariant()
    $patternChars = $patternLower.ToCharArray()
    $candidateChars = $candidateLower.ToCharArray()
    
    # 快速路径2: 单字符匹配优化
    if ($patternLen -eq 1) {
        for ($i = 0; $i -lt $candidateLen; $i++) {
            if ($candidateChars[$i] -eq $patternChars[0]) {
                return 100 - $i  # 位置越靠前分数越高
            }
        }
        return 0
    }
    
    # 性能优化3: 字符存在性预检查 - 快速排除不可能的匹配
    # 使用字符数组访问而非字符串索引，减少内存开销
    foreach ($char in $patternChars) {
        $found = $false
        foreach ($candidateChar in $candidateChars) {
            if ($candidateChar -eq $char) {
                $found = $true
                break
            }
        }
        if (-not $found) {
            return 0  # 必需字符不存在, 直接返回
        }
    }
    
    # 核心匹配算法 - 计算字符匹配度和连续性
    $matched = 0           # 匹配的字符数
    $consecutive = 0       # 当前连续匹配长度
    $maxConsecutive = 0    # 最大连续匹配长度
    $candidatePos = 0      # 候选字符串当前搜索位置
    $startBonus = 0        # 起始位置奖励
    
    # 检查是否从开头匹配 (大小写不敏感) - 使用字符数组更高效
    $startsWithPattern = $true
    if ($patternLen -le $candidateLen) {
        for ($i = 0; $i -lt $patternLen; $i++) {
            if ($patternChars[$i] -ne $candidateChars[$i]) {
                $startsWithPattern = $false
                break
            }
        }
        if ($startsWithPattern) {
            $startBonus = 20  # 起始匹配给20分奖励
        }
    }
    
    # 逐字符匹配算法 - 使用预转换的字符数组，减少内存访问开销
    for ($i = 0; $i -lt $patternLen; $i++) {
        $patternChar = $patternChars[$i]
        $found = $false
        
        # 在候选字符串中查找当前模式字符
        for ($j = $candidatePos; $j -lt $candidateLen; $j++) {
            if ($candidateChars[$j] -eq $patternChar) {
                $matched++
                $found = $true
                
                # 连续性检查 - 连续匹配的字符得分更高
                if ($j -eq $candidatePos) {
                    $consecutive++
                    if ($consecutive -gt $maxConsecutive) {
                        $maxConsecutive = $consecutive
                    }
                } else {
                    $consecutive = 1  # 重置连续计数
                }
                
                $candidatePos = $j + 1  # 更新搜索位置
                break
            }
        }
        
        # 如果某个字符未找到, 重置连续计数
        if (-not $found) {
            $consecutive = 0
        }
    }
    
    # 评分计算 - 使用整数运算
    # 基础分数: (匹配字符数 / 模式长度) * 60
    $baseScore = [Math]::Floor(($matched * 60) / $patternLen)
    
    # 连续性奖励: (最大连续长度 / 模式长度) * 20
    $consecutiveBonus = [Math]::Floor(($maxConsecutive * 20) / $patternLen)
    
    # 长度惩罚: 候选字符串越长, 分数略微降低
    $lengthPenalty = [Math]::Min(($candidateLen - $patternLen), 10)
    
    # 最终分数计算
    $finalScore = $baseScore + $consecutiveBonus + $startBonus - $lengthPenalty
    
    # 确保分数在0-100范围内
    return [Math]::Max(0, [Math]::Min(100, $finalScore))
}

# 带缓存的模糊评分函数 - 避免重复计算提高性能
# 参数: $Pattern=输入模式, $Candidate=候选字符串
function Get-{{.SanitizedName}}FuzzyScoreCached {
    param(
        [string]$Pattern,
        [string]$Candidate
    )
    
    $cacheKey = "$Pattern|$Candidate"
    
    # 缓存命中检查
    if ($script:{{.SanitizedName}}_fuzzyCache.ContainsKey($cacheKey)) {
        return $script:{{.SanitizedName}}_fuzzyCache[$cacheKey]
    }
    
    # 计算分数并缓存
    $score = Get-{{.SanitizedName}}FuzzyScoreFast -Pattern $Pattern -Candidate $Candidate
    
    # 缓存大小控制 - 防止内存无限增长
    if ($script:{{.SanitizedName}}_fuzzyCache.Count -gt $script:{{.SanitizedName}}_FUZZY_CACHE_MAX_SIZE) {
        $script:{{.SanitizedName}}_fuzzyCache.Clear()  # 清空缓存
    }
    
    $script:{{.SanitizedName}}_fuzzyCache[$cacheKey] = $score
    return $score
}

# 智能补全匹配函数 - 重构版匹配策略
# 参数: $Pattern=输入模式, $Options=候选选项数组  
function Get-{{.SanitizedName}}IntelligentMatches {
    param(
        [string]$Pattern,
        [array]$Options
    )
    
    $patternLen = $Pattern.Length
    $totalCandidates = $Options.Count
    
    # 空模式时返回所有选项 (用于Tab补全初始状态) 
    if ([string]::IsNullOrEmpty($Pattern)) {
        return $Options
    }
    
    # 🔥 新的智能匹配策略：多层级渐进式匹配
    
    # 第1级: 精确前缀匹配 (最高优先级) 
    $exactPrefixMatches = [System.Collections.ArrayList]::new()
    foreach ($option in $Options) {
        if ($option.StartsWith($Pattern, [System.StringComparison]::Ordinal)) {
            [void]$exactPrefixMatches.Add($option)
        }
    }
    
    # 精确前缀匹配如果有结果, 优先返回 (但不过度限制数量) 
    if ($exactPrefixMatches.Count -gt 0) {
        return $exactPrefixMatches.ToArray()
    }
    
    # 第2级: 大小写不敏感前缀匹配
    $caseInsensitiveMatches = [System.Collections.ArrayList]::new()
    foreach ($option in $Options) {
        if ($option.StartsWith($Pattern, [System.StringComparison]::OrdinalIgnoreCase)) {
            [void]$caseInsensitiveMatches.Add($option)
        }
    }
    
    # 大小写不敏感匹配如果有结果, 返回
    if ($caseInsensitiveMatches.Count -gt 0) {
        return $caseInsensitiveMatches.ToArray()
    }
    
    # 第3级: 子字符串匹配 (基本模糊匹配) 🔥重新加入
    $substringMatches = [System.Collections.ArrayList]::new()
    $patternLower = $Pattern.ToLowerInvariant()
    foreach ($option in $Options) {
        if ($option.ToLowerInvariant().Contains($patternLower)) {
            [void]$substringMatches.Add($option)
        }
    }
    
    # 子字符串匹配如果有结果, 返回
    if ($substringMatches.Count -gt 0) {
        return $substringMatches.ToArray()
    }
    
    # 第4级: 智能模糊匹配 (高级模糊匹配) 
    if ($script:{{.SanitizedName}}_FUZZY_COMPLETION_ENABLED -and $patternLen -ge $script:{{.SanitizedName}}_FUZZY_MIN_PATTERN_LENGTH -and $totalCandidates -le $script:{{.SanitizedName}}_FUZZY_MAX_CANDIDATES) {
        $scoredMatches = [System.Collections.ArrayList]::new()
        
        foreach ($option in $Options) {
            $score = Get-{{.SanitizedName}}FuzzyScoreCached -Pattern $Pattern -Candidate $option
            
            # 🔥降低阈值, 提高匹配率 (原阈值可能太高) 
            if ($score -ge ($script:{{.SanitizedName}}_FUZZY_SCORE_THRESHOLD * 0.7)) {
                [void]$scoredMatches.Add(@{
                    Option = $option
                    Score = $score
                })
            }
        }
        
        if ($scoredMatches.Count -gt 0) {
            # 按分数排序, 返回前N个最佳匹配
            $sortedMatches = $scoredMatches | Sort-Object Score -Descending
            
            $fuzzyResults = [System.Collections.ArrayList]::new()
            $count = 0
            foreach ($match in $sortedMatches) {
                if ($count -ge $script:{{.SanitizedName}}_FUZZY_MAX_RESULTS) { break }
                [void]$fuzzyResults.Add($match.Option)
                $count++
            }
            
            return $fuzzyResults.ToArray()
        }
    }
    
    # 🔥 最终 fallback：返回空数组 (让用户知道没有匹配到) 
    return @()
}

# ==================== 文件路径补全核心函数 ====================

# 专用文件路径补全函数 - 为{{.SanitizedName}}提供智能路径补全
# 参数: $WordToComplete=当前输入的单词
# 返回: 匹配的文件和目录路径数组
function Get-{{.SanitizedName}}PathCompletions {
    param(
        [string]$WordToComplete
    )
    
    $pathMatches = [System.Collections.ArrayList]::new()
    
    # 获取当前路径的目录部分
    $basePath = if ($WordToComplete -and (Split-Path $WordToComplete -Parent)) {
        Split-Path $WordToComplete -Parent
    } else {
        "."
    }
    
    # 获取文件名部分用于过滤
    $fileName = if ($WordToComplete) {
        Split-Path $WordToComplete -Leaf
    } else {
        ""
    }
    
    # 预编译文件名匹配模式
    $filePattern = "$fileName*"
    
    try {
        # 获取目录和文件
        $items = Get-ChildItem -Path $basePath -ErrorAction SilentlyContinue | Where-Object {
            $_.Name -like $filePattern
        }
        
        foreach ($item in $items) {
            $fullPath = if ($basePath -eq ".") {
                $item.Name
            } else {
                Join-Path $basePath $item.Name
            }
            
            # 目录添加路径分隔符
            if ($item.PSIsContainer) {
                [void]$pathMatches.Add("$fullPath/")
            } else {
                [void]$pathMatches.Add($fullPath)
            }
        }
    }
    catch {
        # 路径访问失败时返回空数组 - 静默处理错误
        Write-Debug "路径访问失败: $($_.Exception.Message)"
    }
    
    return $pathMatches.ToArray()
}

# -------------------------- Completion Logic Implementation ------------------------
$scriptBlock = {
    param(
        $wordToComplete,
        $commandAst,
        $cursorPosition
    )

    # 初始化缓存和索引 (仅在首次调用时创建) 
    if (-not $script:{{.SanitizedName}}_contextIndex) {
        $script:{{.SanitizedName}}_contextIndex = @{}
        $script:{{.SanitizedName}}_flagIndex = @{}
        
        # 预编译正则表达式以提高性能
        $script:{{.SanitizedName}}_flagRegex = [regex]::new('^-', [System.Text.RegularExpressions.RegexOptions]::Compiled)
        
        # 构建上下文索引以提高查找性能
        foreach ($item in ${{.SanitizedName}}_cmdTree) {
            if ($item.Context) {
                $script:{{.SanitizedName}}_contextIndex[$item.Context] = $item
            }
        }
        
        # 构建标志索引以提高查找性能
        foreach ($flag in ${{.SanitizedName}}_flagParams) {
            if ($flag.Context -and $flag.Parameter) {
                $key = "$($flag.Context)|$($flag.Parameter)"
                $script:{{.SanitizedName}}_flagIndex[$key] = $flag
            }
        }
    }

    try {
        # 1. 解析令牌
        $tokens = $commandAst.CommandElements | ForEach-Object { $_.Extent.Text }
        if (-not $tokens -or $tokens.Count -eq 0) {
            return @()
        }
        
        $currentIndex = $tokens.Count - 1
        $prevElement = if ($currentIndex -ge 1) { $tokens[$currentIndex - 1] } else { $null }

        # 快速路径：如果当前输入看起来像是路径，优先提供路径补全
        if ($wordToComplete -match '[/\~\.]' -or $wordToComplete -like './*' -or $wordToComplete -like '../*') {
            return Get-{{.SanitizedName}}PathCompletions -WordToComplete $wordToComplete
        }

        # 2. 计算当前命令上下文 (优化版本) 
        $context = "/"
        for ($i = 1; $i -le $currentIndex; $i++) {
            $elem = $tokens[$i]
            if ($script:{{.SanitizedName}}_flagRegex.IsMatch($elem)) { break }
            
            $nextContext = "$context$elem/"
            # 使用索引进行O(1)查找
            if ($script:{{.SanitizedName}}_contextIndex.ContainsKey($nextContext)) {
                $context = $nextContext
            } else {
                break
            }
        }

        # 3. 获取当前上下文的可用选项 (优化版本) 
        $currentContextItem = $script:{{.SanitizedName}}_contextIndex[$context]
        $currentOptions = if ($currentContextItem) { $currentContextItem.Options } else { @() }

        # 4. 优先补全当前级别的所有选项 (子命令 + 标志) - 使用智能匹配
        if ($currentOptions -and $currentOptions.Count -gt 0) {
            # 使用智能匹配获取最佳选项匹配 - 这是关键修复！
            $intelligentMatches = Get-{{.SanitizedName}}IntelligentMatches -Pattern $wordToComplete -Options $currentOptions
            
            if ($intelligentMatches.Count -gt 0) {
                # 使用ArrayList提高数组操作性能
                $matchingOptions = [System.Collections.ArrayList]::new()
                
                foreach ($option in $intelligentMatches) {
                    $result = if ($script:{{.SanitizedName}}_flagRegex.IsMatch($option)) { $option } else { "$option " }
                    [void]$matchingOptions.Add($result)
                }
                
                return $matchingOptions.ToArray()
            }
        }

        # 5. 枚举/预设值补全
        if ($prevElement -and $script:{{.SanitizedName}}_flagRegex.IsMatch($prevElement)) {
            $flagKey = "$context|$prevElement"
            $paramDef = $script:{{.SanitizedName}}_flagIndex[$flagKey]
            
            if ($paramDef) {
                switch ($paramDef.ValueType) {
                    'enum' {
                        # 统一使用智能匹配进行枚举值补全
                        # 空值时会智能返回所有枚举值, 有值时进行智能匹配
                        $enumMatches = Get-{{.SanitizedName}}IntelligentMatches -Pattern $wordToComplete -Options $paramDef.Options
                        return $enumMatches
                    }
                    'string' {
                        # 字符串类型 - 使用专用函数提供文件和目录路径补全
                        return Get-{{.SanitizedName}}PathCompletions -WordToComplete $wordToComplete
                    }
                    default {
                        # bool类型或其他非字符串类型标志后, 用户可能要输入新参数或路径, 使用专用函数提供文件路径补全
                        return Get-{{.SanitizedName}}PathCompletions -WordToComplete $wordToComplete
                    }
                }
            }
        }

         # 6. 补全标志本身 (如 --ty -> --type) - 使用智能匹配
        if ($script:{{.SanitizedName}}_flagRegex.IsMatch($wordToComplete)) {
            # 收集当前上下文的所有标志 - 使用ArrayList优化性能
            $contextFlags = [System.Collections.ArrayList]::new()
            foreach ($flag in ${{.SanitizedName}}_flagParams) {
                if ($flag.Context -eq $context) {
                    [void]$contextFlags.Add($flag.Parameter)
                }
            }
            
            if ($contextFlags.Count -gt 0) {
                # 使用智能匹配获取最佳标志匹配
                $flagMatches = Get-{{.SanitizedName}}IntelligentMatches -Pattern $wordToComplete -Options $contextFlags.ToArray()
                if ($flagMatches.Count -gt 0) {
                    return $flagMatches
                }
            }
        }

        # 7. 无匹配
        return @()
    }
    catch {
        # 错误处理：返回空数组而不是抛出异常
        Write-Debug "PowerShell补全错误: $($_.Exception.Message)"
        return @()
    }
}

# ==================== 调试和诊断功能 ====================

# 补全系统健康检查函数 (可选, 用于调试)
function Get-{{.SanitizedName}}CompletionDebug {
    Write-Host "=== {{.SanitizedName}} PowerShell补全系统诊断 ===" -ForegroundColor Cyan
    Write-Host "PowerShell版本: $($PSVersionTable.PSVersion)" -ForegroundColor Green
    Write-Host "补全函数状态: $(if (Get-Command Register-ArgumentCompleter -ErrorAction SilentlyContinue) { '已注册' } else { '未注册' })" -ForegroundColor Green
    Write-Host "命令树条目数: $(${{.SanitizedName}}_cmdTree.Count)" -ForegroundColor Green
    Write-Host "标志参数数: $(${{.SanitizedName}}_flagParams.Count)" -ForegroundColor Green
    Write-Host "模糊补全状态: $(if ($script:{{.SanitizedName}}_FUZZY_COMPLETION_ENABLED) { '启用' } else { '禁用' })" -ForegroundColor Green
    Write-Host "候选项阈值: $script:{{.SanitizedName}}_FUZZY_MAX_CANDIDATES" -ForegroundColor Green
    Write-Host "缓存条目数: $($script:{{.SanitizedName}}_fuzzyCache.Count)" -ForegroundColor Green
    Write-Host ""
    Write-Host "使用方法: 在PowerShell中输入 'Get-{{.SanitizedName}}CompletionDebug' 查看此信息" -ForegroundColor Yellow
}

# 注册补全函数-带原始名称 (可能包含扩展名) 
Register-ArgumentCompleter -CommandName ${{.SanitizedName}}_commandName -ScriptBlock $scriptBlock

# 注册补全函数-不带扩展名 (仅当与原始名称不同时才注册) 
${{.SanitizedName}}_withoutExt = [System.IO.Path]::GetFileNameWithoutExtension("{{.ProgramName}}")
if (${{.SanitizedName}}_withoutExt -ne ${{.SanitizedName}}_commandName) {
    Register-ArgumentCompleter -CommandName ${{.SanitizedName}}_withoutExt -ScriptBlock $scriptBlock
}
`
)
