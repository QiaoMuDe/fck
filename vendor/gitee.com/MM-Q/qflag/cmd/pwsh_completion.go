// pwsh_completion.go - PowerShell 补全脚本生成器
package cmd

import (
	"bytes"
	"fmt"
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
//
// 参数:
// - cmdTreeEntries: 命令树条目缓冲区
// - cmdPath: 命令路径
// - cmdOpts: 命令选项
func generatePwshCommandTreeEntry(cmdTreeEntries *bytes.Buffer, cmdPath string, cmdOpts []string) {
	// 预分配缓冲区容量，减少内存分配
	optsBuf := bytes.NewBuffer(make([]byte, 0, len(cmdOpts)*20)) // 假设平均每个选项20字节

	// 格式化命令选项
	formatOptions(optsBuf, cmdOpts, escapePwshString)

	// 添加命令树条目
	fmt.Fprintf(cmdTreeEntries, PwshCmdTreeItem, cmdPath, optsBuf.String())
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
		// 生成标志参数条目
		// 生成带枚举选项的标志参数条目
		enumOptions := ""
		if param.ValueType == "enum" && len(param.EnumOptions) > 0 {
			optionsBuf := bytes.NewBuffer(make([]byte, 0, len(param.EnumOptions)*15))
			formatOptions(optionsBuf, param.EnumOptions, escapePwshString)
			enumOptions = optionsBuf.String()
		}
		fmt.Fprintf(flagParamsBuf, PwshFlagParamItem,
			param.CommandPath,
			param.Name,
			param.Type,
			param.ValueType,
			enumOptions,
		)

		// 条目之间添加逗号，非最后一个条目
		if i < len(params)-1 {
			flagParamsBuf.WriteString(",\n")
		}
	}

	// 生成根命令条目
	rootCmdEntry := fmt.Sprintf(PwshCmdTreeItem, "/", rootOptsBuf.String())
	if cmdTreeEntries != "" {
		rootCmdEntry += ",\n" + cmdTreeEntries
	}

	// 写入PowerShell自动补全脚本
	fmt.Fprintf(buf, PwshFunctionHeader,
		programName,
		rootCmdEntry,
		flagParamsBuf.String(),
	) // 移除独立的枚举选项数组
}

// escapePwshString 转义PowerShell字符串中的特殊字符
// 优化：单次循环处理所有转义，减少字符串分配
func escapePwshString(s string) string {
	// 预计算所需容量：最坏情况下每个字符都需要转义
	buf := make([]byte, 0, len(s)*2)
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		case '\'':
			buf = append(buf, '\'', '\'') // 单引号转义为两个单引号
		case '\\':
			buf = append(buf, '\\', '\\') // 反斜杠转义为两个反斜杠
		default:
			buf = append(buf, c)
		}
	}
	return string(buf)
}

const (
	// 标志参数条目(含枚举选项)
	PwshFlagParamItem = "	@{ Context = \"%s\"; Parameter = \"%s\"; ParamType = \"%s\"; ValueType = \"%s\"; Options = @(%s) }"
	// 命令树条目
	PwshCmdTreeItem = "	@{ Context = \"%s\"; Options = @(%s) }"
)

const (
	// PowerShell自动补全脚本头部
	PwshFunctionHeader = `# -------------------------- Configuration Area (Need to be modified according to actual commands) --------------------------
# Command Name
$commandName = "%s"

# 1. Command Tree
$cmdTree = @(
%s
)

# 2. Flag Parameter Definitions
$flagParams = @(
%s
)

# -----------------------------------------------------------------------------------

# -------------------------- Completion Logic Implementation ------------------------
$scriptBlock = {
    param(
        $wordToComplete,
        $commandAst,
        $cursorPosition
    )

    # 1. Parse tokens
    $tokens = $commandAst.CommandElements | ForEach-Object { $_.Extent.Text }
    $currentIndex = $tokens.Count - 1
    $prevElement = if ($currentIndex -ge 1) { $tokens[$currentIndex - 1] } else { $null }

    # 2. Calculate the current command context
    $context = "/"
    for ($i = 1; $i -le $currentIndex; $i++) {
        $elem = $tokens[$i]
        if ($elem -match '^-') { break }
        $nextContext = "$context$elem/"
        $contextMatch = $cmdTree | Where-Object { $_.Context -eq $nextContext }
        if ($contextMatch) {
            $context = $nextContext
        } else {
            break
        }
    }

    # 3. Available options in the current context
    $currentOptions = ($cmdTree | Where-Object { $_.Context -eq $context }).Options

    # 4. First complete all options (subcommands + flags) at the current level
    if ($currentOptions) {
        $matchingOptions = $currentOptions | Where-Object {
            $_ -like "$wordToComplete*"
        }
        if ($matchingOptions) {
            return $matchingOptions | ForEach-Object {
                if ($_ -match '^-') { $_ } else { "$_ " }
            }
        }
    }

    # 5. Complete flags themselves (like --ty -> --type)
    if ($wordToComplete -match '^-') {
        $flagDefs = $flagParams | Where-Object { $_.Context -eq $context }
        $flagMatches = $flagDefs | Where-Object {
            $_.Parameter -like "$wordToComplete*"
        } | ForEach-Object { $_.Parameter }
        return $flagMatches
    }

    # 6. Enum/Preset value completion
    # 6a Current token is empty → Complete all enum values of the previous flag
    if (-not $wordToComplete -and $prevElement -match '^-') {
        $paramDef = $flagParams | Where-Object {
            $_.Context -eq $context -and $_.Parameter -eq $prevElement
        }
        if ($paramDef) {
            switch ($paramDef.ValueType) {
                'enum'   { return $paramDef.Options }
                'path'   { return Get-ChildItem  -Name}
                'number' { return 1..10 | ForEach-Object { "$_" } }
                'ip'     { return @('192.168.','10.0.','172.16.','127.0.0.') }
                'url'    { return @('http://','https://','ftp://') }
                default  { return @() }
            }
        }
    }

    # 6b The current token is not empty, and the previous token is a flag that requires a value → Filter with prefix
    $flagForValue = $tokens[$currentIndex - 1]
    if ($flagForValue -match '^-' -and $currentIndex -ge 1) {
        $paramDef = $flagParams | Where-Object {
            $_.Context -eq $context -and $_.Parameter -eq $flagForValue
        }
        if ($paramDef) {
            switch ($paramDef.ValueType) {
                'path' {
                    $pattern = if ($wordToComplete) { "$wordToComplete*" } else { '*' }
                    Get-ChildItem -Name $pattern -ErrorAction SilentlyContinue
                }
                'number' { return 1..100 | Where-Object { "$_" -like "$wordToComplete*" } }
                'ip'     { return @('192.168.','10.0.','172.16.','127.0.0.') | Where-Object { $_ -like "$wordToComplete*" } }
                'enum'   { return $paramDef.Options | Where-Object { $_ -like "$wordToComplete*" } }
                'url'    { return @('http://','https://','ftp://') | Where-Object { $_ -like "$wordToComplete*" } }
                default  { return @() }
            }
        }
    }

    # 7. No match
    return @()
}

Register-ArgumentCompleter -CommandName $commandName -ScriptBlock $scriptBlock
`
)
