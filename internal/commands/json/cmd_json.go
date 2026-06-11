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
	"gitee.com/MM-Q/go-kit/fs"
	"gitee.com/MM-Q/go-kit/term"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// JsonConfig 配置结构体
type JsonConfig struct {
	Pretty     bool     // 美化输出
	Compact    bool     // 压缩输出
	Validate   bool     // 验证模式
	Query      string   // 查询路径
	SetValue   []string // 设置值表达式: ["path1=val1", "path2=val2"]
	DeletePath []string // 删除路径: ["path1", "path2"]
	ValueType  string   // 值类型: auto/string/number/bool
	Highlight  bool     // 语法高亮
	Raw        bool     // 原始字符串输出
	Write      bool     // 原地写入
	Backup     bool     // 写入前备份
	Files      []string // 位置参数（文件路径）
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

	// 3. 设置模式 — 使用 sjson 设置字段值
	if len(config.SetValue) > 0 {
		result, err := setJSON(data, config.SetValue, config.ValueType)
		if err != nil {
			return err
		}
		// 如果请求美化输出，重新格式化
		if config.Pretty {
			return writeFormatted(result, config)
		}
		return writeOutput(result, config)
	}

	// 4. 删除模式 — 使用 sjson 删除字段
	if len(config.DeletePath) > 0 {
		result, err := deleteJSON(data, config.DeletePath)
		if err != nil {
			return err
		}
		// 如果请求美化输出，重新格式化
		if config.Pretty {
			return writeFormatted(result, config)
		}
		return writeOutput(result, config)
	}

	// 5. 查询处理 (优先于解析，直接操作原始数据)
	if config.Query != "" {
		result, err := queryJSON(data, config.Query)
		if err != nil {
			return err
		}
		// 查询结果直接输出
		return writeOutput([]byte(result.Raw), config)
	}

	// 6. 解析 JSON
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("json parse failed: %w", err)
	}

	// 7. 格式化输出
	output, err := formatOutput(jsonData, config)
	if err != nil {
		return err
	}

	// 8. 输出结果
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
		// 管道输入时禁用 -w
		if config.Write {
			return nil, fmt.Errorf("cannot use -w flag with pipe input")
		}

		// 管道输入时禁用高亮模式
		if config.Highlight {
			config.Highlight = false
		}

		// 读取标准输入
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

// setJSON 使用路径表达式设置 JSON 字段值
//
// 使用 tidwall/sjson 库实现高性能设置
// 支持多个操作: ["path1=val1", "path2=val2"]
// 路径语法与 gjson 兼容，支持:
//   - 对象属性: users.name
//   - 数组索引: users.0
//   - 负数索引: users.-1 (追加)
//   - 嵌套路径: address.city
//
// 参数:
//   - data: JSON 原始数据
//   - pairs: 设置表达式列表，元素格式为 "path=value"
//   - valueType: 值类型 (auto/string/number/bool)
//
// 返回值:
//   - []byte: 修改后的 JSON
//   - error: 设置错误
func setJSON(data []byte, pairs []string, valueType string) ([]byte, error) {
	if len(pairs) == 0 {
		return data, nil
	}

	result := string(data)

	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		idx := strings.Index(pair, "=")
		if idx < 0 {
			return nil, fmt.Errorf("invalid set expression (missing '='): %s", pair)
		}

		path := strings.TrimSpace(pair[:idx])
		value := strings.TrimSpace(pair[idx+1:])

		if path == "" {
			return nil, fmt.Errorf("invalid set expression (empty path): %s", pair)
		}

		var sjsonResult string
		var err error

		switch valueType {
		case "number":
			sjsonResult, err = sjson.Set(result, path, json.Number(value))
		case "bool":
			boolVal := value == "true"
			sjsonResult, err = sjson.Set(result, path, boolVal)
		case "string":
			sjsonResult, err = sjson.Set(result, path, value)
		default: // "auto"
			// 自动推断：sjson 会根据 Go 类型自动推断
			// 空值处理
			if value == "null" {
				sjsonResult, err = sjson.Set(result, path, nil)
			} else if value == "true" {
				sjsonResult, err = sjson.Set(result, path, true)
			} else if value == "false" {
				sjsonResult, err = sjson.Set(result, path, false)
			} else if isNumeric(value) {
				sjsonResult, err = sjson.Set(result, path, json.Number(value))
			} else {
				sjsonResult, err = sjson.Set(result, path, value)
			}
		}

		if err != nil {
			return nil, fmt.Errorf("set failed for path '%s': %w", path, err)
		}
		result = sjsonResult
	}

	return []byte(result), nil
}

// isNumeric 检查字符串是否为有效的数字
//
// 参数:
//   - s: 待检查的字符串
//
// 返回值:
//   - bool: 是否为数字
func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			if c != '.' && c != '-' && c != '+' {
				return false
			}
		}
	}
	return true
}

// deleteJSON 使用路径表达式删除 JSON 字段
//
// 使用 tidwall/sjson 库实现
// 支持多个路径: ["path1", "path2"]
//
// 参数:
//   - data: JSON 原始数据
//   - paths: 删除路径列表
//
// 返回值:
//   - []byte: 修改后的 JSON
//   - error: 删除错误
func deleteJSON(data []byte, paths []string) ([]byte, error) {
	if len(paths) == 0 {
		return data, nil
	}

	result := string(data)

	for _, path := range paths {
		path = strings.TrimSpace(path)
		if path == "" {
			continue
		}

		var sjsonResult string
		var err error
		sjsonResult, err = sjson.Delete(result, path)
		if err != nil {
			return nil, fmt.Errorf("delete failed for path '%s': %w", path, err)
		}
		result = sjsonResult
	}

	return []byte(result), nil
}

// writeFormatted 格式化 JSON 数据并输出
//
// 对 set/delete 的结果进行重新格式化（美化或压缩），然后输出
//
// 参数:
//   - data: JSON 原始数据
//   - config: 命令配置
//
// 返回值:
//   - error: 输出错误
func writeFormatted(data []byte, config JsonConfig) error {
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("json parse failed: %w", err)
	}
	output, err := formatOutput(jsonData, config)
	if err != nil {
		return err
	}
	return writeOutput(output, config)
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

	// 原地写入模式
	if config.Write {
		if len(config.Files) == 0 {
			return fmt.Errorf("-w flag requires a file path argument")
		}
		targetFile := config.Files[0]

		// 创建备份（如果启用）
		if config.Backup {
			backupFile := targetFile + ".bak"
			if err := fs.CopyEx(targetFile, backupFile, true); err != nil {
				return fmt.Errorf("failed to create backup: %w", err)
			}
		}

		// 原子写入：临时文件+移动
		tmpFile := targetFile + ".tmp"
		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write temp file: %w", err)
		}
		if err := fs.MoveEx(tmpFile, targetFile, true); err != nil {
			return fmt.Errorf("failed to move temp file: %w", err)
		}

		return nil
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
