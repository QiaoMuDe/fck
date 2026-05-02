package curl

import (
	"fmt"
	"net/http"
)

// Formatter 响应格式化器
type Formatter struct {
	color bool
}

// NewFormatter 创建格式化器
//
// 参数:
//   - color: 是否启用彩色输出
//
// 返回值:
//   - *Formatter: 格式化器实例
func NewFormatter(color bool) *Formatter {
	return &Formatter{
		color: color,
	}
}

// PrintVerbose 打印详细输出
//
// 参数:
//   - resp: 响应对象
func (f *Formatter) PrintVerbose(resp *Response) {
	// 状态行和响应头
	fmt.Print(highlightHeaders(resp.Status, resp.Headers, f.color))
	fmt.Println()

	// 响应体
	f.PrintBody(resp)

	// 统计信息
	fmt.Printf("Time: %v\n", resp.Time)
	fmt.Printf("Size: %d bytes\n", len(resp.Body))
}

// PrintHeaders 打印响应头
//
// 参数:
//   - resp: 响应对象
func (f *Formatter) PrintHeaders(resp *Response) {
	fmt.Print(highlightHeaders(resp.Status, resp.Headers, f.color))
	fmt.Println()
}

// PrintBody 打印响应体
//
// 参数:
//   - resp: 响应对象
func (f *Formatter) PrintBody(resp *Response) {
	if len(resp.Body) == 0 {
		return
	}

	// 获取 Content-Type
	contentType := resp.Headers.Get("Content-Type")

	// 使用 chroma 高亮
	highlighted := highlightBody(resp.Body, contentType, f.color)
	fmt.Println(highlighted)
}

// PrintError 打印错误
//
// 参数:
//   - err: 错误
func (f *Formatter) PrintError(err error) {
	if f.color {
		fmt.Printf("%sError: %v%s\n", colorRed, err, colorReset)
	} else {
		fmt.Printf("Error: %v\n", err)
	}
}

// PrintResponse 打印完整响应（带高亮）
//
// 参数:
//   - resp: 响应对象
//   - includeHeaders: 是否包含响应头
//   - headOnly: 是否仅显示响应头
func (f *Formatter) PrintResponse(resp *Response, includeHeaders, headOnly bool) {
	// 打印响应头
	if includeHeaders || headOnly {
		fmt.Print(highlightHeaders(resp.Status, resp.Headers, f.color))
		if headOnly {
			return
		}
		fmt.Println()
	}

	// 打印响应体
	if len(resp.Body) > 0 {
		contentType := resp.Headers.Get("Content-Type")
		highlighted := highlightBody(resp.Body, contentType, f.color)
		fmt.Println(highlighted)
	}
}

// PrintHighlightedHeaders 打印高亮响应头（兼容旧代码）
//
// 参数:
//   - status: 状态行
//   - headers: 响应头
func (f *Formatter) PrintHighlightedHeaders(status string, headers http.Header) {
	fmt.Print(highlightHeaders(status, headers, f.color))
	fmt.Println()
}

// PrintHighlightedBody 打印高亮响应体（兼容旧代码）
//
// 参数:
//   - body: 响应体
//   - contentType: Content-Type
func (f *Formatter) PrintHighlightedBody(body []byte, contentType string) {
	if len(body) == 0 {
		return
	}
	highlighted := highlightBody(body, contentType, f.color)
	fmt.Println(highlighted)
}
