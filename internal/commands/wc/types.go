package wc

// 扩展 scanner 缓冲区以支持长行（1MB）
const maxCapacity = 1024 * 1024

// WcConfig 命令配置
type WcConfig struct {
	ShowLines   bool     // 显示行数
	ShowWords   bool     // 显示单词数
	ShowBytes   bool     // 显示字节数
	ShowChars   bool     // 显示字符数
	ShowMaxLine bool     // 显示最大行长度
	ShowAll     bool     // 显示所有统计项
	TableStyle  string   // 表格样式
	Files       []string // 输入文件列表
}

// WcStats 统计结果
type WcStats struct {
	Lines      int64  // 行数
	Words      int64  // 单词数
	Bytes      int64  // 字节数
	Chars      int64  // 字符数
	MaxLineLen int64  // 最大行长度
	Filename   string // 文件名（管道输入时为 "-"）
}
