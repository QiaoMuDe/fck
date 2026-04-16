package cat

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// ansiRegex 匹配 ANSI 转义序列
var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

// FileViewer 统一文件查看器
type FileViewer struct {
	config   *CatConfig
	hlConfig HighlightConfig
}

// NewFileViewer 创建文件查看器
func NewFileViewer(config *CatConfig) *FileViewer {
	return &FileViewer{
		config:   config,
		hlConfig: DefaultHighlightConfig(),
	}
}

// View 查看文件（统一入口）
//
// 参数:
//   - path: 文件路径
//
// 返回:
//   - error: 处理错误
func (v *FileViewer) View(path string) error {
	// 1. 读取文件
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}

	// 2. 统一换行符
	content = bytes.ReplaceAll(content, []byte("\r\n"), []byte("\n"))

	// 3. 按行分割
	allLines := bytes.Split(content, []byte("\n"))

	// 4. 根据 head/tail 切片
	lines := v.sliceLines(allLines)

	// 5. 处理特殊显示标志（仅非高亮模式）
	if !v.config.Highlight {
		lines = v.applySpecialFlags(lines)
	}

	// 6. 输出（根据是否高亮决定）
	if v.config.Highlight {
		return v.outputHighlighted(lines, path)
	}
	return v.outputPlain(lines)
}

// sliceLines 根据 head/tail 标志切片
//
// 参数:
//   - lines: 所有行
//
// 返回:
//   - [][]byte: 切片后的行
func (v *FileViewer) sliceLines(lines [][]byte) [][]byte {
	// head 模式
	if v.config.HeadLines > 0 {
		if v.config.HeadLines < len(lines) {
			return lines[:v.config.HeadLines]
		}
		return lines
	}

	// tail 模式
	if v.config.TailLines > 0 {
		if v.config.TailLines < len(lines) {
			start := len(lines) - v.config.TailLines
			return lines[start:]
		}
		return lines
	}

	// 全部
	return lines
}

// applySpecialFlags 应用特殊显示标志
//
// 参数:
//   - lines: 行数据
//
// 返回:
//   - [][]byte: 处理后的行
//
// 说明:
//   - 仅在高亮模式为 false 时调用
//   - 处理 -E(行尾$), -T(制表符^I)
func (v *FileViewer) applySpecialFlags(lines [][]byte) [][]byte {
	// 仅处理 -E(行尾$), -T(制表符^I)
	if !v.config.ShowEnd && !v.config.ShowTabs {
		return lines
	}

	result := make([][]byte, len(lines))
	for i, line := range lines {
		var parts []string

		// 处理 -T: 制表符显示为 ^I
		content := string(line)
		if v.config.ShowTabs {
			content = strings.ReplaceAll(content, "\t", "^I")
		}
		parts = append(parts, content)

		// 处理 -E: 显示行尾 $
		if v.config.ShowEnd {
			parts = append(parts, "$")
		}

		result[i] = []byte(strings.Join(parts, ""))
	}

	return result
}

// outputPlain 普通模式输出
//
// 参数:
//   - lines: 行数据
//
// 返回:
//   - error: 输出错误
func (v *FileViewer) outputPlain(lines [][]byte) error {
	lineNum := 0

	for _, line := range lines {
		content := string(line)
		isBlank := len(strings.TrimSpace(content)) == 0

		// 处理行号 (-n 或 -b)
		if v.config.ShowLineNum || (v.config.ShowNonBlank && !isBlank) {
			lineNum++
			fmt.Printf("%6d\t", lineNum)
		}

		// 输出内容
		fmt.Println(content)
	}

	return nil
}

// outputHighlighted 高亮模式输出
//
// 参数:
//   - lines: 行数据
//   - path: 文件路径（用于检测语言）
//
// 返回:
//   - error: 输出错误
func (v *FileViewer) outputHighlighted(lines [][]byte, path string) error {
	// 1. 合并内容
	content := bytes.Join(lines, []byte("\n"))

	// 2. 语法高亮
	hlContent, ok, err := highlightContent(content, path, v.hlConfig)
	if err != nil || !ok {
		// 高亮失败，降级为普通输出
		return v.outputPlain(lines)
	}

	// 3. 按行分割高亮后的内容
	hlLines := bytes.Split(hlContent, []byte("\n"))

	// 4. 输出行（带行号）
	lineNum := 0

	for _, line := range hlLines {
		content := string(line)
		// 移除 ANSI 码后判断是否为空行
		isBlank := len(strings.TrimSpace(stripANSI(content))) == 0

		// 处理行号 (-n 或 -b)
		if v.config.ShowLineNum || (v.config.ShowNonBlank && !isBlank) {
			lineNum++
			fmt.Printf("%6d\t", lineNum)
		}

		// 输出高亮内容
		fmt.Println(content)
	}

	return nil
}

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

// stripANSI 移除字符串中的 ANSI 转义序列
//
// 参数:
//   - s: 输入字符串
//
// 返回:
//   - string: 移除 ANSI 码后的字符串
func stripANSI(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}
