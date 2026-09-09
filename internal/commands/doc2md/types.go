// Package doc2md 提供办公文档转 Markdown 的转换功能。
// 基于 gitee.com/MM-Q/doc2md 库封装，支持本地文件与 stdin 管道输入。
package doc2md

// DocExtNone 是 -x/--extension 未指定时的哨兵默认值。
// 管道输入时若仍为该值，说明用户未显式指定扩展名，应报错。
const DocExtNone = "none"

// Doc2mdConfig 是 doc2md 命令的配置结构
type Doc2mdConfig struct {
	// Files 输入文档路径列表（支持通配符展开）
	Files []string
	// Extension 扩展名提示（stdin 管道输入时必须显式指定，未指定时为 DocExtNone）
	Extension string
	// MIMEType MIME 类型提示
	MIMEType string
	// Charset 字符集提示（如: gbk）
	Charset string
	// Output 输出 Markdown 文件路径（为空时输出到 stdout）
	Output string
	// KeepDataURIs 是否保留完整 base64 图片数据 URI（默认截断）
	KeepDataURIs bool
}
