// Package base64 实现了 base64 编解码命令
// 支持编码/解码、URL 安全变体、无填充模式
package base64

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"

	"gitee.com/MM-Q/go-kit/term"
)

// Base64Config 配置结构体
type Base64Config struct {
	Decode    bool     // 解码模式
	Strings   []string // 位置参数字符串
	FilePath  string   // 输入文件路径
	Output    string   // 输出文件路径
	Wrap      int      // 每行字符数（0表示不换行）
	URLSafe   bool     // URL 安全变体
	NoPadding bool     // 禁用填充
}

// Base64CmdMain 执行 base64 命令
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - error: 执行错误
func Base64CmdMain(config Base64Config) error {
	if config.Decode {
		return decode(config)
	}
	return encode(config)
}

// encode 执行编码操作
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - error: 执行错误
func encode(config Base64Config) error {
	// 读取输入
	data, err := readInput(config)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return fmt.Errorf("no input provided")
	}

	// 创建 base64 编码器
	enc := base64.StdEncoding
	if config.URLSafe {
		enc = base64.URLEncoding
	}
	if config.NoPadding {
		enc = enc.WithPadding(base64.NoPadding)
	}

	// 编码
	encoded := enc.EncodeToString(data)

	// 换行处理
	if config.Wrap > 0 {
		encoded = wrapLines(encoded, config.Wrap)
	}

	// 输出（编码结果添加换行符）
	encoded += "\n"
	return writeOutput([]byte(encoded), config.Output)
}

// decode 执行解码操作
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - error: 执行错误
func decode(config Base64Config) error {
	// 读取输入
	data, err := readInput(config)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return fmt.Errorf("no input provided")
	}

	// 清理换行符（解码时忽略）
	data = bytes.ReplaceAll(data, []byte("\n"), nil)
	data = bytes.ReplaceAll(data, []byte("\r"), nil)

	// 创建解码器
	enc := base64.StdEncoding
	if config.URLSafe {
		enc = base64.URLEncoding
	}

	// 解码
	decoded, err := enc.DecodeString(string(data))
	if err != nil {
		return fmt.Errorf("decode failed: %w", err)
	}

	// 输出（二进制数据不加换行）
	return writeOutput(decoded, config.Output)
}

// readInput 读取输入数据
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - []byte: 输入数据
//   - error: 读取错误
func readInput(config Base64Config) ([]byte, error) {
	// 1. 优先检测管道/重定向输入
	if term.IsStdinPipe() {
		return io.ReadAll(os.Stdin)
	}

	// 2. 从文件读取
	if config.FilePath != "" {
		return os.ReadFile(config.FilePath)
	}

	// 3. 从位置参数读取（空格连接）
	if len(config.Strings) > 0 {
		return []byte(strings.Join(config.Strings, " ")), nil
	}

	return nil, fmt.Errorf("no input provided")
}

// wrapLines 按指定宽度换行
//
// 参数:
//   - s: 输入字符串
//   - width: 每行宽度
//
// 返回值:
//   - string: 换行后的字符串
func wrapLines(s string, width int) string {
	if width <= 0 || len(s) <= width {
		return s
	}

	var result strings.Builder
	for i := 0; i < len(s); i += width {
		end := i + width
		if end > len(s) {
			end = len(s)
		}
		result.WriteString(s[i:end])
		if end < len(s) {
			result.WriteString("\n")
		}
	}
	return result.String()
}

// writeOutput 写入输出
//
// 参数:
//   - data: 输出数据
//   - filePath: 输出文件路径（空表示标准输出）
//
// 返回值:
//   - error: 写入错误
func writeOutput(data []byte, filePath string) error {
	if filePath == "" {
		_, err := os.Stdout.Write(data)
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}
