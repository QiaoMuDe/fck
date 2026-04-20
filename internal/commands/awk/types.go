package awk

// CharRange 字符范围定义
type CharRange struct {
	Start int // 起始位置（1-based，包含）
	End   int // 结束位置（1-based，包含），-1表示到行尾
}

// IsValid 检查范围是否有效
func (cr CharRange) IsValid() bool {
	if cr.Start < 1 {
		return false
	}
	if cr.End != -1 && cr.End < cr.Start {
		return false
	}
	return true
}

// Extract 从字符串中提取字符范围
func (cr CharRange) Extract(s string) string {
	runes := []rune(s) // 支持Unicode
	start := cr.Start - 1

	if start >= len(runes) {
		return ""
	}

	end := cr.End
	if end == -1 || end > len(runes) {
		end = len(runes)
	}

	return string(runes[start:end])
}

// AwkConfig 命令配置
type AwkConfig struct {
	Pattern     string      // 正则匹配模式
	Fields      []int       // 字段索引列表（0=整行, -1=NF, 1+=正常字段）
	FieldSep    string      // 输入分隔符
	OutputSep   string      // 输出分隔符
	ShowLineNum bool        // 显示行号
	Files       []string    // 输入文件列表
	CharRanges  []CharRange // 字符范围列表
}

// IsFieldMode 是否为字段提取模式
func (c AwkConfig) IsFieldMode() bool {
	return len(c.Fields) > 0
}

// IsCharMode 是否为字符提取模式
func (c AwkConfig) IsCharMode() bool {
	return len(c.CharRanges) > 0
}
