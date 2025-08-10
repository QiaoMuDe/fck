// Package completion PowerShell 自动补全实现
// 本文件实现了PowerShell环境下的命令行自动补全功能，
// 生成PowerShell补全脚本，支持标志和子命令的智能补全。
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

		// 如果不是第一个选项，则添加逗号
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
// 使用对象池优化内存分配，避免创建临时缓冲区和Replacer
//
// 参数:
// - cmdTreeEntries: 命令树条目缓冲区
// - cmdPath: 命令路径
// - cmdOpts: 命令选项
func generatePwshCommandTreeEntry(cmdTreeEntries *bytes.Buffer, cmdPath string, cmdOpts []string) {
	// 使用对象池构建命令树条目，避免创建临时缓冲区和strings.NewReplacer的开销
	cmdTreeItem := buildString(func(builder *strings.Builder) {
		builder.WriteString("\t@{ Context = \"")
		builder.WriteString(cmdPath)
		builder.WriteString("\"; Options = @(")

		// 直接在builder中格式化选项，避免额外的字符串分配
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

		// 条目之间添加逗号，非最后一个条目
		if i < len(params)-1 {
			flagParamsBuf.WriteString(",\n")
		}
	}

	// 清理程序名，去除可能的后缀
	sanitizedProgramName := strings.TrimSuffix(programName, filepath.Ext(programName))

	// 生成根命令条目
	rootReplacer := strings.NewReplacer(
		"{{.Context}}", "/",
		"{{.Options}}", rootOptsBuf.String(),
	)
	// 生成根命令条目
	rootCmdEntry := rootReplacer.Replace(PwshCmdTreeItem)

	// 如果命令树条目不为空，则添加逗号
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
// 使用全局map提高转义性能，避免重复的switch判断
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
// 优化版本：使用全局map进行O(1)查找，提升性能
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
# 模糊补全功能开关 (设置为$false禁用，$true启用)
$script:{{.SanitizedName}}_FUZZY_COMPLETION_ENABLED = $true

# 启用模糊补全的最大候选项数量阈值
# 超过此数量将回退到传统前缀匹配以保证性能
$script:{{.SanitizedName}}_FUZZY_MAX_CANDIDATES = 120

# 模糊匹配的最小输入长度 (小于此长度不启用模糊匹配)
$script:{{.SanitizedName}}_FUZZY_MIN_PATTERN_LENGTH = 2

# 模糊匹配分数阈值 (0-100，分数低于此值的匹配将被过滤)
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
    
    # 性能优化1: 长度预检查 - 候选项太短直接返回0
    if ($candidateLen -lt $patternLen) {
        return 0
    }
    
    # 性能优化2: 完全匹配检查 - 避免不必要的复杂计算
    if ($Candidate.StartsWith($Pattern, [System.StringComparison]::OrdinalIgnoreCase)) {
        return 100  # 前缀完全匹配给最高分
    }
    
    # 性能优化3: 字符存在性预检查 - 快速排除不可能的匹配
    $patternLower = $Pattern.ToLowerInvariant()
    $candidateLower = $Candidate.ToLowerInvariant()
    
    foreach ($char in $patternLower.ToCharArray()) {
        if ($candidateLower.IndexOf($char) -eq -1) {
            return 0  # 必需字符不存在，直接返回
        }
    }
    
    # 核心匹配算法 - 计算字符匹配度和连续性
    $matched = 0           # 匹配的字符数
    $consecutive = 0       # 当前连续匹配长度
    $maxConsecutive = 0    # 最大连续匹配长度
    $candidatePos = 0      # 候选字符串当前搜索位置
    $startBonus = 0        # 起始位置奖励
    
    # 检查是否从开头匹配 (大小写不敏感)
    if ($candidateLower.StartsWith($patternLower)) {
        $startBonus = 20  # 起始匹配给20分奖励
    }
    
    # 逐字符匹配算法
    for ($i = 0; $i -lt $patternLen; $i++) {
        $patternChar = $patternLower[$i]
        $found = $false
        
        # 在候选字符串中查找当前模式字符
        for ($j = $candidatePos; $j -lt $candidateLen; $j++) {
            if ($candidateLower[$j] -eq $patternChar) {
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
        
        # 如果某个字符未找到，重置连续计数
        if (-not $found) {
            $consecutive = 0
        }
    }
    
    # 评分计算 - 使用整数运算
    # 基础分数: (匹配字符数 / 模式长度) * 60
    $baseScore = [Math]::Floor(($matched * 60) / $patternLen)
    
    # 连续性奖励: (最大连续长度 / 模式长度) * 20
    $consecutiveBonus = [Math]::Floor(($maxConsecutive * 20) / $patternLen)
    
    # 长度惩罚: 候选字符串越长，分数略微降低
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

# 智能补全匹配函数 - 分级匹配策略
# 参数: $Pattern=输入模式, $Options=候选选项数组
function Get-{{.SanitizedName}}IntelligentMatches {
    param(
        [string]$Pattern,
        [array]$Options
    )
    
    $patternLen = $Pattern.Length
    $totalCandidates = $Options.Count
    
    # 性能保护: 候选项过多时禁用模糊匹配
    if ($totalCandidates -gt $script:{{.SanitizedName}}_FUZZY_MAX_CANDIDATES) {
        # 回退到传统前缀匹配
        $prefixMatches = @()
        foreach ($option in $Options) {
            if ($option -like "$Pattern*") {
                $prefixMatches += $option
            }
        }
        return $prefixMatches
    }
    
    # 第1级: 精确前缀匹配 (最快，优先级最高) - 使用ArrayList优化性能
    $exactMatches = [System.Collections.ArrayList]::new()
    foreach ($option in $Options) {
        if ($option.StartsWith($Pattern, [System.StringComparison]::Ordinal)) {
            [void]$exactMatches.Add($option)
        }
    }
    
    # 如果有精确匹配且数量合理，直接返回
    if ($exactMatches.Count -gt 0 -and $exactMatches.Count -le 12) {
        return $exactMatches.ToArray()
    }
    
    # 第2级: 大小写不敏感前缀匹配
    if ($exactMatches.Count -eq 0) {
        $caseInsensitiveMatches = [System.Collections.ArrayList]::new()
        foreach ($option in $Options) {
            if ($option.StartsWith($Pattern, [System.StringComparison]::OrdinalIgnoreCase)) {
                [void]$caseInsensitiveMatches.Add($option)
            }
        }
        
        # 如果有大小写不敏感匹配，返回
        if ($caseInsensitiveMatches.Count -gt 0) {
            return $caseInsensitiveMatches.ToArray()
        }
    }
    
    # 第3级: 模糊匹配 (最慢，仅在必要时使用) - 使用ArrayList优化性能
    if ($script:{{.SanitizedName}}_FUZZY_COMPLETION_ENABLED -and $patternLen -ge $script:{{.SanitizedName}}_FUZZY_MIN_PATTERN_LENGTH) {
        $scoredMatches = [System.Collections.ArrayList]::new()
        
        # 对所有候选项进行模糊评分
        foreach ($option in $Options) {
            $score = Get-{{.SanitizedName}}FuzzyScoreCached -Pattern $Pattern -Candidate $option
            
            # 只保留分数达到阈值的匹配
            if ($score -ge $script:{{.SanitizedName}}_FUZZY_SCORE_THRESHOLD) {
                [void]$scoredMatches.Add(@{
                    Option = $option
                    Score = $score
                })
            }
        }
        
        # 如果有模糊匹配结果，按分数排序并返回前N个
        if ($scoredMatches.Count -gt 0) {
            # 按分数降序排序
            $sortedMatches = $scoredMatches | Sort-Object Score -Descending
            
            # 提取选项名称，限制返回数量 - 使用ArrayList优化
            $fuzzyResults = [System.Collections.ArrayList]::new()
            $count = 0
            foreach ($match in $sortedMatches) {
                if ($count -ge $script:{{.SanitizedName}}_FUZZY_MAX_RESULTS) {
                    break
                }
                [void]$fuzzyResults.Add($match.Option)
                $count++
            }
            
            return $fuzzyResults.ToArray()
        }
    }
    
    # 第4级: 子字符串匹配 (最后的备选方案) - 使用ArrayList优化性能
    $substringMatches = [System.Collections.ArrayList]::new()
    $patternLower = $Pattern.ToLowerInvariant()
    
    foreach ($option in $Options) {
        $optionLower = $option.ToLowerInvariant()
        if ($optionLower.Contains($patternLower)) {
            [void]$substringMatches.Add($option)
        }
    }
    
    if ($substringMatches.Count -gt 0) {
        return $substringMatches.ToArray()
    }
    
    # 如果所有匹配策略都失败，返回空数组
    return @()
}

# -------------------------- Completion Logic Implementation ------------------------
$scriptBlock = {
    param(
        $wordToComplete,
        $commandAst,
        $cursorPosition
    )

    # 初始化缓存和索引（仅在首次调用时创建）
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

        # 2. 计算当前命令上下文（优化版本）
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

        # 3. 获取当前上下文的可用选项（优化版本）
        $currentContextItem = $script:{{.SanitizedName}}_contextIndex[$context]
        $currentOptions = if ($currentContextItem) { $currentContextItem.Options } else { @() }

        # 4. 优先补全当前级别的所有选项（子命令 + 标志）- 使用智能匹配
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

        # 5. 补全标志本身（如 --ty -> --type）- 使用智能匹配
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

        # 6. 枚举/预设值补全
        if ($prevElement -and $script:{{.SanitizedName}}_flagRegex.IsMatch($prevElement)) {
            $flagKey = "$context|$prevElement"
            $paramDef = $script:{{.SanitizedName}}_flagIndex[$flagKey]
            
            if ($paramDef) {
                switch ($paramDef.ValueType) {
                    'enum' {
                        if (-not $wordToComplete) {
                            # 当前单词为空 → 返回所有枚举值
                            return $paramDef.Options
                        } else {
                            # 使用智能匹配进行枚举值补全
                            $enumMatches = Get-{{.SanitizedName}}IntelligentMatches -Pattern $wordToComplete -Options $paramDef.Options
                            return $enumMatches
                        }
                    }
                    'string' {
                        # 字符串类型 - 提供文件和目录路径补全 - 使用ArrayList优化性能
                        $pathMatches = [System.Collections.ArrayList]::new()
                        
                        # 获取当前路径的目录部分
                        $basePath = if ($wordToComplete -and (Split-Path $wordToComplete -Parent)) {
                            Split-Path $wordToComplete -Parent
                        } else {
                            "."
                        }
                        
                        # 获取文件名部分用于过滤
                        $fileName = if ($wordToComplete) {
                            Split-Path $wordToComplete -Leaf
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
                            # 路径访问失败时返回空数组
                        }
                        
                        return $pathMatches.ToArray()
                    }
                    default {
                        return @()
                    }
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

# 补全系统健康检查函数 (可选，用于调试)
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

# 模糊匹配测试函数 (用于调试和验证)
function Test-{{.SanitizedName}}FuzzyMatch {
    param(
        [Parameter(Mandatory=$true)]
        [string]$Pattern,
        [Parameter(Mandatory=$true)]
        [string]$Candidate
    )
    
    $score = Get-{{.SanitizedName}}FuzzyScoreFast -Pattern $Pattern -Candidate $Candidate
    Write-Host "模式: '$Pattern' 匹配候选: '$Candidate' 得分: $score" -ForegroundColor Cyan
    
    # 详细分析
    if ($score -ge 80) {
        Write-Host "匹配质量: 优秀" -ForegroundColor Green
    } elseif ($score -ge 50) {
        Write-Host "匹配质量: 良好" -ForegroundColor Yellow
    } elseif ($score -ge $script:{{.SanitizedName}}_FUZZY_SCORE_THRESHOLD) {
        Write-Host "匹配质量: 可接受" -ForegroundColor DarkYellow
    } else {
        Write-Host "匹配质量: 不匹配" -ForegroundColor Red
    }
    
    return $score
}

Register-ArgumentCompleter -CommandName ${{.SanitizedName}}_commandName -ScriptBlock $scriptBlock
`
)
