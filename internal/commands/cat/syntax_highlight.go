package cat

import (
	"bytes"
	"fmt"

	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// HighlightConfig 语法高亮配置
type HighlightConfig struct {
	Enabled   bool   // 是否启用高亮
	Style     string // 主题样式
	Formatter string // 格式化器类型: "terminal256", "terminal16m", "terminal"
}

// DefaultHighlightConfig 返回默认高亮配置
func DefaultHighlightConfig() HighlightConfig {
	return HighlightConfig{
		Enabled:   true,
		Style:     "monokai",
		Formatter: "terminal256", // 使用 256 色，兼容性最好
	}
}

// highlightContent 对内容进行语法高亮
//
// 参数:
//   - content: 文件内容
//   - filename: 文件名（用于检测语言）
//   - config: 高亮配置
//
// 返回:
//   - []byte: 高亮后的 ANSI 格式内容
//   - bool: 是否成功高亮
//   - error: 错误信息
func highlightContent(content []byte, filename string, config HighlightConfig) ([]byte, bool, error) {
	if !config.Enabled {
		return content, false, nil
	}

	// 1. 根据文件名获取 lexer
	// lexers.Match 会根据文件名模式匹配最合适的 lexer
	// 如果返回 nil，表示 Chroma 不支持该文件类型
	lexer := lexers.Match(filename)
	if lexer == nil {
		return content, false, nil // 不支持的语言，返回原内容
	}

	// 2. 获取 formatter（ANSI 终端格式）
	formatter := formatters.Get(config.Formatter)
	if formatter == nil {
		formatter = formatters.Fallback
	}

	// 3. 获取样式
	style := styles.Get(config.Style)
	if style == nil {
		style = styles.Fallback
	}

	// 4. 执行高亮
	iterator, err := lexer.Tokenise(nil, string(content))
	if err != nil {
		return content, false, fmt.Errorf("tokenise failed: %w", err)
	}

	var buf bytes.Buffer
	err = formatter.Format(&buf, style, iterator)
	if err != nil {
		return content, false, fmt.Errorf("format failed: %w", err)
	}

	return buf.Bytes(), true, nil
}

// shouldHighlight 判断文件是否应该高亮
//
// 参数:
//   - filename: 文件名
//
// 返回:
//   - bool: 是否应该高亮
//
// 说明:
//
//	使用 Chroma 的 lexers.Match 函数自动检测文件类型。
//	如果返回 nil，表示 Chroma 不支持该文件类型的语法高亮。
func shouldHighlight(filename string) bool {
	return lexers.Match(filename) != nil
}
