package awk

// AwkConfig 命令配置
type AwkConfig struct {
	Pattern     string   // 正则匹配模式
	Fields      []int    // 字段索引列表（0=整行, -1=NF, 1+=正常字段）
	FieldSep    string   // 输入分隔符
	OutputSep   string   // 输出分隔符
	ShowLineNum bool     // 显示行号
	Files       []string // 输入文件列表
}
