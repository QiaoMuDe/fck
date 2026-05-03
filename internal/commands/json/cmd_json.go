// Package json 实现了 JSON 数据处理命令
// 提供格式化、验证、查询等功能
package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gitee.com/MM-Q/fck/internal/types"
	"gitee.com/MM-Q/go-kit/term"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/tidwall/gjson"
)

// JsonConfig 配置结构体
type JsonConfig struct {
	Pretty    bool     // 美化输出
	Compact   bool     // 压缩输出
	Validate  bool     // 验证模式
	Query     string   // 查询路径
	Highlight bool     // 语法高亮
	Raw       bool     // 原始字符串输出
	Files     []string // 位置参数（文件路径）
}

// JsonStats 操作统计
type JsonStats struct {
	IsValid    bool   // 是否有效
	ParseTime  string // 解析耗时
	OutputSize int    // 输出大小
}

// JsonCmdMain 执行 json 命令
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - error: 执行错误
func JsonCmdMain(config JsonConfig) error {
	// 1. 读取输入
	data, err := readInput(config)
	if err != nil {
		return err
	}

	// 2. 验证模式
	if config.Validate {
		return validateJSON(data)
	}

	// 3. 查询处理 (优先于解析，直接操作原始数据)
	if config.Query != "" {
		result, err := queryJSON(data, config.Query)
		if err != nil {
			return err
		}
		// 查询结果直接输出
		return writeOutput([]byte(result.Raw), config)
	}

	// 4. 解析 JSON
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("json parse failed: %w", err)
	}

	// 5. 格式化输出
	output, err := formatOutput(jsonData, config)
	if err != nil {
		return err
	}

	// 6. 输出结果
	return writeOutput(output, config)
}

// readInput 读取输入数据
//
// 输入方式:
//  1. 管道/重定向输入 (使用 term.IsStdinPipe() 检测) - 传递JSON字符串
//  2. 位置参数 - 指定文件路径，读取文件内容
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - []byte: 输入数据
//   - error: 读取错误
func readInput(config JsonConfig) ([]byte, error) {
	// 1. 优先检测管道/重定向输入
	if term.IsStdinPipe() {
		return readAllStdin()
	}

	// 2. 从文件读取
	if len(config.Files) > 0 {
		// 检查是否传入多个文件，如果是则报错
		if len(config.Files) > 1 {
			return nil, fmt.Errorf("only one file path can be specified, got %d", len(config.Files))
		}
		// 读取单个文件
		return os.ReadFile(config.Files[0])
	}

	// 3. 没有管道输入也没有文件参数，报错提示
	return nil, fmt.Errorf("no input data provided, please pass JSON string via pipe or specify file path")
}

// readAllStdin 读取标准输入所有数据
//
// 返回值:
//   - []byte: 输入数据
//   - error: 读取错误
func readAllStdin() ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("failed to read stdin: %w", err)
	}
	return buf.Bytes(), nil
}

// validateJSON 验证 JSON 有效性
//
// 参数:
//   - data: JSON 数据
//
// 返回值:
//   - error: 验证错误 (nil 表示有效)
func validateJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("json validation failed: %w", err)
	}
	fmt.Println("✓ json is valid")
	return nil
}

// queryJSON 使用路径表达式查询 JSON
//
// 使用 tidwall/gjson 库实现高性能查询
// 路径语法:
//   - 对象属性: users.name
//   - 数组索引: users.0
//   - 负数索引: users.-1 (最后一个)
//   - 通配符:   users.*.name (所有元素, 内部转换为 users.#.name)
//
// 参数:
//   - data: JSON 原始数据 ([]byte)
//   - path: 查询路径
//
// 返回值:
//   - gjson.Result: 查询结果
//   - error: 查询错误
func queryJSON(data []byte, path string) (gjson.Result, error) {
	// 将 * 通配符转换为 gjson 的 # 语法
	path = convertWildcardToGjson(path)

	result := gjson.GetBytes(data, path)
	if !result.Exists() {
		return result, fmt.Errorf("query path invalid or result does not exist: %s", path)
	}

	return result, nil
}

// convertWildcardToGjson 将 * 通配符转换为 gjson 的 # 语法
//
// 参数:
//   - path: 原始路径
//
// 返回值:
//   - string: 转换后的路径
func convertWildcardToGjson(path string) string {
	// 将 .* 替换为 .#
	// 将 [*] 替换为 [#]
	path = strings.ReplaceAll(path, ".*", ".#")
	path = strings.ReplaceAll(path, "[*]", "[#]")
	return path
}

// formatOutput 格式化输出
//
// 参数:
//   - data: JSON 数据
//   - config: 命令配置
//
// 返回值:
//   - []byte: 格式化后的数据
//   - error: 格式化错误
func formatOutput(data interface{}, config JsonConfig) ([]byte, error) {
	var output []byte
	var err error

	if config.Compact {
		// 压缩输出
		output, err = json.Marshal(data)
	} else {
		// 美化输出
		output, err = json.MarshalIndent(data, "", "  ")
	}

	if err != nil {
		return nil, fmt.Errorf("json format failed: %w", err)
	}

	// 添加换行符
	output = append(output, '\n')
	return output, nil
}

// writeOutput 输出结果
//
// 参数:
//   - data: 输出数据
//   - config: 命令配置
//
// 返回值:
//   - error: 输出错误
func writeOutput(data []byte, config JsonConfig) error {
	content := string(data)

	// 处理原始字符串输出 (-r)
	if config.Raw {
		content = processRawOutput(content)
	}

	// 确保内容以换行符结尾
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}

	// 如果需要语法高亮
	if config.Highlight {
		return highlightAndPrint(content)
	}

	// 直接输出
	_, err := os.Stdout.WriteString(content)
	return err
}

// processRawOutput 处理原始字符串输出
// 如果内容是带引号的字符串，去除引号
//
// 参数:
//   - content: 原始内容
//
// 返回值:
//   - string: 处理后的内容
func processRawOutput(content string) string {
	content = strings.TrimSpace(content)

	// 检查是否是带引号的字符串
	if len(content) >= 2 && content[0] == '"' && content[len(content)-1] == '"' {
		// 去除首尾引号并处理转义字符
		var result strings.Builder
		i := 1
		for i < len(content)-1 {
			if content[i] == '\\' && i+1 < len(content)-1 {
				// 处理转义字符
				switch content[i+1] {
				case '"':
					result.WriteByte('"')
				case '\\':
					result.WriteByte('\\')
				case '/':
					result.WriteByte('/')
				case 'b':
					result.WriteByte('\b')
				case 'f':
					result.WriteByte('\f')
				case 'n':
					result.WriteByte('\n')
				case 'r':
					result.WriteByte('\r')
				case 't':
					result.WriteByte('\t')
				default:
					result.WriteByte(content[i+1])
				}
				i += 2
			} else {
				result.WriteByte(content[i])
				i++
			}
		}
		return result.String()
	}

	return content
}

// highlightAndPrint 语法高亮并输出
//
// 参数:
//   - content: JSON 内容
//
// 返回值:
//   - error: 输出错误
func highlightAndPrint(content string) error {
	lexer := lexers.Get("json")
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	style := styles.Get(types.HighlightStyleDefault)
	if style == nil {
		style = styles.Fallback
	}

	formatter := formatters.Get(types.HighlightFormatter256)
	if formatter == nil {
		formatter = formatters.Fallback
	}

	iterator, err := lexer.Tokenise(nil, content)
	if err != nil {
		return fmt.Errorf("syntax highlighting failed: %w", err)
	}

	return formatter.Format(os.Stdout, style, iterator)
}
