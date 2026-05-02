package cat

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"gitee.com/MM-Q/fck/internal/types"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/noborus/ov/oviewer"
)

// ansiRegex 匹配 ANSI 转义序列
var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

// HighlightConfig 语法高亮配置
type HighlightConfig struct {
	Enabled   bool   // 是否启用高亮
	Style     string // 主题样式
	Formatter string // 格式化器类型: "terminal256", "terminal16m", "terminal"
}

// DefaultHighlightConfig 返回默认高亮配置
//
// 返回:
//   - HighlightConfig: 默认高亮配置
func DefaultHighlightConfig() HighlightConfig {
	return HighlightConfig{
		Enabled:   true,
		Style:     types.HighlightStyleDefault,
		Formatter: types.HighlightFormatter256,
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

// OutputDirectly 直接输出内容
//
// 参数:
//   - content: 内容
//   - name: 名称（用于高亮检测）
//   - config: 配置
//
// 返回:
//   - error: 错误
func OutputDirectly(content []byte, name string, config CatConfig) error {
	if len(content) == 0 {
		return nil
	}

	lines := SplitLines(content)

	// 高亮模式
	if config.Highlight {
		return outputHighlighted(lines, name, config)
	}

	// 普通模式
	return outputPlain(lines, config)
}

// outputPlain 普通模式输出
//
// 参数:
//   - lines: 按行分割的内容
//   - config: 配置
//
// 返回:
//   - error: 错误
func outputPlain(lines [][]byte, config CatConfig) error {
	lineNum := 0

	for _, line := range lines {
		content := string(line)
		isBlank := len(content) == 0

		// 处理行号 (-n 或 -b)
		if config.ShowLineNum || (config.ShowNonBlank && !isBlank) {
			lineNum++
			fmt.Printf("%6d\t", lineNum)
		}

		// 处理 -T: 制表符显示为 ^I
		if config.ShowTabs {
			content = strings.ReplaceAll(content, "\t", "^I")
		}

		// 输出内容
		fmt.Print(content)

		// 处理 -E: 显示行尾 $
		if config.ShowEnd {
			fmt.Print("$")
		}

		fmt.Println()
	}

	return nil
}

// outputHighlighted 高亮模式输出
//
// 参数:
//   - lines: 按行分割的内容
//   - name: 名称（用于高亮检测）
//   - config: 配置
//
// 返回:
//   - error: 错误
func outputHighlighted(lines [][]byte, name string, config CatConfig) error {
	// 合并内容
	content := bytes.Join(lines, []byte("\n"))

	// 语法高亮
	hlConfig := DefaultHighlightConfig()
	hlContent, ok, err := highlightContent(content, name, hlConfig)
	if err != nil || !ok {
		// 高亮失败，降级为普通输出
		return outputPlain(lines, config)
	}

	// 按行分割高亮后的内容
	hlLines := bytes.Split(hlContent, []byte("\n"))

	// 输出行（带行号）
	lineNum := 0

	for _, line := range hlLines {
		content := string(line)
		// 移除 ANSI 码后判断是否为空行
		isBlank := len(strings.TrimSpace(stripANSI(content))) == 0

		// 处理行号 (-n 或 -b)
		if config.ShowLineNum || (config.ShowNonBlank && !isBlank) {
			lineNum++
			fmt.Printf("%6d\t", lineNum)
		}

		// 输出高亮内容
		fmt.Println(content)
	}

	return nil
}

// OutputWithPager 使用分页器输出
//
// 参数:
//   - content: 内容
//   - name: 名称（用于高亮检测）
//   - config: 配置
//
// 返回:
//   - error: 错误
func OutputWithPager(content []byte, name string, config CatConfig) error {
	if len(content) == 0 {
		return nil
	}

	// 创建 Document
	doc, err := oviewer.NewDocument()
	if err != nil {
		return fmt.Errorf("failed to create document: %w", err)
	}

	doc.FileName = name

	// 判断是否进行语法高亮
	var reader io.Reader

	if config.Highlight && shouldHighlight(name) {
		hlConfig := DefaultHighlightConfig()
		hlContent, ok, hlErr := highlightContent(content, name, hlConfig)
		if hlErr != nil {
			// 高亮失败，输出警告但继续使用原内容
			fmt.Fprintf(os.Stderr, "warn: syntax highlight failed: %v\n", hlErr)
			reader = bytes.NewReader(content)
		} else if ok {
			reader = bytes.NewReader(hlContent)
			doc.Caption = "[syntax highlighted]"
		} else {
			reader = bytes.NewReader(content)
		}
	} else {
		reader = bytes.NewReader(content)
	}

	// 使用 ControlReader 加载内容
	if err := doc.ControlReader(reader, nil); err != nil {
		return fmt.Errorf("failed to load content: %w", err)
	}

	// 创建 oviewer Root
	root, err := oviewer.NewOviewer(doc)
	if err != nil {
		return fmt.Errorf("failed to create oviewer: %w", err)
	}

	// 运行分页器
	return root.Run()
}
