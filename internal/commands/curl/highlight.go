package curl

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// ANSI 颜色码
const (
	colorReset   = "\033[0m"
	colorGray    = "\033[90m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorWhite   = "\033[37m"
)

// HighlightConfig 高亮配置
type HighlightConfig struct {
	Enabled   bool   // 是否启用高亮
	Style     string // 主题样式
	Formatter string // 格式化器类型
}

// DefaultHighlightConfig 返回默认高亮配置
func DefaultHighlightConfig() HighlightConfig {
	return HighlightConfig{
		Enabled:   true,
		Style:     "monokai",
		Formatter: "terminal256",
	}
}

// highlightHeaders 高亮响应头
//
// 参数:
//   - status: 状态行 (如 "HTTP/1.1 200 OK")
//   - headers: 响应头
//   - color: 是否启用颜色
//
// 返回值:
//   - string: 高亮后的字符串
func highlightHeaders(status string, headers http.Header, color bool) string {
	if !color {
		var buf bytes.Buffer
		buf.WriteString(status)
		buf.WriteString("\n")
		for key, values := range headers {
			for _, value := range values {
				buf.WriteString(fmt.Sprintf("%s: %s\n", key, value))
			}
		}
		return buf.String()
	}

	var buf bytes.Buffer

	// 高亮状态行
	buf.WriteString(highlightStatusLine(status))
	buf.WriteString("\n")

	// 高亮响应头
	for key, values := range headers {
		for _, value := range values {
			buf.WriteString(highlightHeader(key, value))
			buf.WriteString("\n")
		}
	}

	return buf.String()
}

// highlightStatusLine 高亮状态行
//
// 参数:
//   - status: 状态行 (如 "HTTP/1.1 200 OK")
//
// 返回值:
//   - string: 高亮后的状态行
func highlightStatusLine(status string) string {
	parts := strings.SplitN(status, " ", 3)
	if len(parts) < 2 {
		return status
	}

	var buf bytes.Buffer

	// HTTP 版本 - 灰色
	buf.WriteString(colorGray)
	buf.WriteString(parts[0])
	buf.WriteString(colorReset)
	buf.WriteString(" ")

	// 状态码 - 根据状态码着色
	statusCode := parts[1]
	statusColor := getStatusColor(statusCode)
	buf.WriteString(statusColor)
	buf.WriteString(statusCode)
	buf.WriteString(colorReset)

	// 状态描述
	if len(parts) > 2 {
		buf.WriteString(" ")
		buf.WriteString(parts[2])
	}

	return buf.String()
}

// getStatusColor 根据状态码返回颜色
//
// 参数:
//   - statusCode: 状态码字符串
//
// 返回值:
//   - string: ANSI 颜色码
func getStatusColor(statusCode string) string {
	if len(statusCode) < 1 {
		return colorWhite
	}

	switch statusCode[0] {
	case '2':
		return colorGreen // 2xx 成功
	case '3':
		return colorYellow // 3xx 重定向
	case '4':
		return colorRed // 4xx 客户端错误
	case '5':
		return colorMagenta // 5xx 服务器错误
	default:
		return colorWhite
	}
}

// highlightHeader 高亮单个响应头
//
// 参数:
//   - key: Header 键
//   - value: Header 值
//
// 返回值:
//   - string: 高亮后的 Header
func highlightHeader(key, value string) string {
	var buf bytes.Buffer

	// Header Key - 蓝色
	buf.WriteString(colorBlue)
	buf.WriteString(key)
	buf.WriteString(colorReset)

	// 冒号 - 灰色
	buf.WriteString(colorGray)
	buf.WriteString(":")
	buf.WriteString(colorReset)

	// 空格和 Value
	buf.WriteString(" ")
	buf.WriteString(highlightHeaderValue(key, value))

	return buf.String()
}

// highlightHeaderValue 根据 Header 类型高亮值
//
// 参数:
//   - key: Header 键
//   - value: Header 值
//
// 返回值:
//   - string: 高亮后的值
func highlightHeaderValue(key, value string) string {
	lowerKey := strings.ToLower(key)

	// Content-Type 特殊处理
	if lowerKey == "content-type" {
		return highlightContentTypeValue(value)
	}

	// Set-Cookie 特殊处理
	if lowerKey == "set-cookie" || lowerKey == "cookie" {
		return highlightCookieValue(value)
	}

	return value
}

// highlightContentTypeValue 高亮 Content-Type 值
//
// 参数:
//   - value: Content-Type 值
//
// 返回值:
//   - string: 高亮后的值
func highlightContentTypeValue(value string) string {
	// 分离 MIME 类型和参数
	parts := strings.SplitN(value, ";", 2)
	mimeType := strings.TrimSpace(parts[0])

	var buf bytes.Buffer

	// MIME 类型 - 青色
	buf.WriteString(colorCyan)
	buf.WriteString(mimeType)
	buf.WriteString(colorReset)

	// 参数（如 charset=utf-8）
	if len(parts) > 1 {
		buf.WriteString(colorGray)
		buf.WriteString(";")
		buf.WriteString(parts[1])
		buf.WriteString(colorReset)
	}

	return buf.String()
}

// highlightCookieValue 高亮 Cookie 值
//
// 参数:
//   - value: Cookie 值
//
// 返回值:
//   - string: 高亮后的值
func highlightCookieValue(value string) string {
	// 分离 name=value 和属性
	parts := strings.Split(value, ";")
	if len(parts) == 0 {
		return value
	}

	var buf bytes.Buffer

	// 第一个部分是 name=value
	nameValue := strings.TrimSpace(parts[0])
	nvParts := strings.SplitN(nameValue, "=", 2)
	if len(nvParts) == 2 {
		// name - 黄色
		buf.WriteString(colorYellow)
		buf.WriteString(nvParts[0])
		buf.WriteString(colorReset)

		// = - 灰色
		buf.WriteString(colorGray)
		buf.WriteString("=")
		buf.WriteString(colorReset)

		// value - 绿色
		buf.WriteString(colorGreen)
		buf.WriteString(nvParts[1])
		buf.WriteString(colorReset)
	} else {
		buf.WriteString(nameValue)
	}

	// 其他属性
	for i := 1; i < len(parts); i++ {
		buf.WriteString(colorGray)
		buf.WriteString("; ")
		buf.WriteString(strings.TrimSpace(parts[i]))
		buf.WriteString(colorReset)
	}

	return buf.String()
}

// highlightBody 高亮响应体
//
// 参数:
//   - body: 响应体内容
//   - contentType: Content-Type
//   - color: 是否启用颜色
//
// 返回值:
//   - string: 高亮后的内容
func highlightBody(body []byte, contentType string, color bool) string {
	if !color || len(body) == 0 {
		return string(body)
	}

	// 检测语言
	lang := detectLanguage(contentType)
	if lang == "" {
		return string(body)
	}

	// 使用 chroma 高亮
	highlighted, ok, err := highlightWithChroma(body, lang, DefaultHighlightConfig())
	if err != nil || !ok {
		return string(body)
	}

	return string(highlighted)
}

// detectLanguage 根据 Content-Type 检测语言
//
// 参数:
//   - contentType: Content-Type 值
//
// 返回值:
//   - string: chroma 语言名称，空字符串表示不支持
func detectLanguage(contentType string) string {
	if contentType == "" {
		return ""
	}

	// 提取 MIME 类型（去掉 charset 等参数）
	mimeType := strings.SplitN(contentType, ";", 2)[0]
	mimeType = strings.TrimSpace(strings.ToLower(mimeType))

	// 从 MIME 类型中提取语言（如 text/html -> html, application/json -> json）
	parts := strings.Split(mimeType, "/")
	if len(parts) == 2 {
		// 提取后半部分作为语言名称
		lang := parts[1]

		// 移除可能的 +xml 等后缀，保留主类型
		if idx := strings.Index(lang, "+"); idx > 0 {
			lang = lang[:idx]
		}

		// 检查 chroma 是否支持该语言
		if lexer := lexers.Get(lang); lexer != nil {
			return lang
		}
	}

	return ""
}

// highlightWithChroma 使用 chroma 进行语法高亮
//
// 参数:
//   - content: 内容
//   - lang: 语言名称
//   - config: 高亮配置
//
// 返回值:
//   - []byte: 高亮后的内容
//   - bool: 是否成功
//   - error: 错误
func highlightWithChroma(content []byte, lang string, config HighlightConfig) ([]byte, bool, error) {
	if !config.Enabled {
		return content, false, nil
	}

	// 获取 lexer
	lexer := lexers.Get(lang)
	if lexer == nil {
		return content, false, nil
	}

	// 获取 formatter
	formatter := formatters.Get(config.Formatter)
	if formatter == nil {
		formatter = formatters.Fallback
	}

	// 获取样式
	style := styles.Get(config.Style)
	if style == nil {
		style = styles.Fallback
	}

	// 执行高亮
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
