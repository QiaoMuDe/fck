package curl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Formatter 响应格式化器
type Formatter struct {
	color  bool
	pretty bool
}

// NewFormatter 创建格式化器
//
// 参数:
//   - color: 是否启用彩色输出
//   - pretty: 是否格式化 JSON
//
// 返回值:
//   - *Formatter: 格式化器实例
func NewFormatter(color, pretty bool) *Formatter {
	return &Formatter{
		color:  color,
		pretty: pretty,
	}
}

// PrintVerbose 打印详细输出
//
// 参数:
//   - resp: 响应对象
func (f *Formatter) PrintVerbose(resp *Response) {
	// 状态行
	if f.color {
		fmt.Printf("\033[36m< %s\033[0m\n", resp.Status)
	} else {
		fmt.Printf("< %s\n", resp.Status)
	}

	// 响应头
	for key, values := range resp.Headers {
		for _, value := range values {
			if f.color {
				fmt.Printf("\033[33m< %s:\033[0m %s\n", key, value)
			} else {
				fmt.Printf("< %s: %s\n", key, value)
			}
		}
	}
	fmt.Println()

	// 响应体
	f.printBodyContent(resp.Body)

	// 统计信息
	fmt.Printf("\nTime: %v\n", resp.Time)
	fmt.Printf("Size: %d bytes\n", len(resp.Body))
}

// PrintHeaders 打印响应头
//
// 参数:
//   - resp: 响应对象
func (f *Formatter) PrintHeaders(resp *Response) {
	if f.color {
		fmt.Printf("\033[36m%s\033[0m\n", resp.Status)
	} else {
		fmt.Println(resp.Status)
	}

	for key, values := range resp.Headers {
		for _, value := range values {
			if f.color {
				fmt.Printf("\033[33m%s:\033[0m %s\n", key, value)
			} else {
				fmt.Printf("%s: %s\n", key, value)
			}
		}
	}
	fmt.Println()
}

// PrintBody 打印响应体
//
// 参数:
//   - resp: 响应对象
func (f *Formatter) PrintBody(resp *Response) {
	f.printBodyContent(resp.Body)
}

// printBodyContent 打印响应体内容
//
// 参数:
//   - body: 响应体字节
func (f *Formatter) printBodyContent(body []byte) {
	if len(body) == 0 {
		return
	}

	// 尝试格式化 JSON
	if f.pretty && f.isJSON(body) {
		var buf bytes.Buffer
		if err := json.Indent(&buf, body, "", "  "); err == nil {
			if f.color {
				f.printColoredJSON(buf.Bytes())
			} else {
				fmt.Println(buf.String())
			}
			return
		}
	}

	// 普通文本输出
	fmt.Println(string(body))
}

// isJSON 检查是否为 JSON
//
// 参数:
//   - data: 数据
//
// 返回值:
//   - bool: 是否为 JSON
func (f *Formatter) isJSON(data []byte) bool {
	data = bytes.TrimSpace(data)
	return len(data) > 0 && (data[0] == '{' || data[0] == '[')
}

// printColoredJSON 打印彩色 JSON
//
// 参数:
//   - data: JSON 数据
func (f *Formatter) printColoredJSON(data []byte) {
	// 简单的 JSON 语法高亮
	str := string(data)
	var result strings.Builder

	inString := false
	escape := false

	for i := 0; i < len(str); i++ {
		ch := str[i]

		if escape {
			result.WriteByte(ch)
			escape = false
			continue
		}

		if ch == '\\' && inString {
			escape = true
			result.WriteByte(ch)
			continue
		}

		if ch == '"' {
			inString = !inString
			if inString {
				result.WriteString("\033[32m\"")
			} else {
				result.WriteString("\"\033[0m")
			}
			continue
		}

		if inString {
			result.WriteByte(ch)
			continue
		}

		// 关键字和符号着色
		switch ch {
		case '{', '}', '[', ']':
			result.WriteString(fmt.Sprintf("\033[33m%c\033[0m", ch))
		case ':':
			result.WriteString("\033[37m:\033[0m")
		case ',':
			result.WriteString("\033[37m,\033[0m")
		case 't', 'f', 'n': // true, false, null
			if i+4 <= len(str) && str[i:i+4] == "true" {
				result.WriteString("\033[36mtrue\033[0m")
				i += 3
			} else if i+5 <= len(str) && str[i:i+5] == "false" {
				result.WriteString("\033[36mfalse\033[0m")
				i += 4
			} else if i+4 <= len(str) && str[i:i+4] == "null" {
				result.WriteString("\033[36mnull\033[0m")
				i += 3
			} else {
				result.WriteByte(ch)
			}
		default:
			result.WriteByte(ch)
		}
	}

	fmt.Println(result.String())
}

// PrintError 打印错误
//
// 参数:
//   - err: 错误
func (f *Formatter) PrintError(err error) {
	if f.color {
		fmt.Fprintf(os.Stderr, "\033[31mError: %v\033[0m\n", err)
	} else {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}
