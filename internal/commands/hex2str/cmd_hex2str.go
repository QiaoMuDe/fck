// Package hex2str 实现了十六进制与字符串的相互转换功能
// 支持编码（字符串→十六进制）和解码（十六进制→字符串）两种模式
package hex2str

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"gitee.com/MM-Q/go-kit/term"
)

// Hex2strConfig 配置结构体
type Hex2strConfig struct {
	Encode bool   // 编码模式：字符串→十六进制
	Upper  bool   // 编码输出使用大写字母
	Input  string // 位置参数输入
}

// Hex2strCmdMain 执行 hex2str 命令
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - error: 执行错误
func Hex2strCmdMain(config Hex2strConfig) error {
	// 获取输入数据
	data, err := readInput(config)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return fmt.Errorf("no input provided")
	}

	// 根据模式执行转换
	if config.Encode {
		return encode(data, config.Upper)
	}
	return decode(data)
}

// readInput 读取输入数据
//
// 优先级:
// 1. 管道/stdin 输入（优先处理管道，符合 Unix 工具惯例）
// 2. 位置参数
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - []byte: 输入数据
//   - error: 读取错误
func readInput(config Hex2strConfig) ([]byte, error) {
	// 1. 优先检测管道/stdin 输入
	if term.IsStdinPipe() {
		reader := bufio.NewReader(os.Stdin)
		data, err := io.ReadAll(reader)
		if err != nil {
			return nil, fmt.Errorf("failed to read from stdin: %w", err)
		}
		// 去除末尾换行符
		return bytes.TrimSpace(data), nil
	}

	// 2. 使用位置参数
	if config.Input != "" {
		return []byte(config.Input), nil
	}

	return nil, fmt.Errorf("no input provided: specify a string argument or use pipe input")
}

// encode 将字符串编码为十六进制
//
// 参数:
//   - data: 输入数据
//   - upper: 是否使用大写字母
//
// 返回值:
//   - error: 执行错误
func encode(data []byte, upper bool) error {
	result := hex.EncodeToString(data)
	if upper {
		result = strings.ToUpper(result)
	}
	fmt.Println(result)
	return nil
}

// decode 将十六进制解码为字符串
//
// 参数:
//   - data: 十六进制输入数据
//
// 返回值:
//   - error: 执行错误
func decode(data []byte) error {
	// 直接解码，不对数据进行任何修改
	result, err := hex.DecodeString(string(data))
	if err != nil {
		return fmt.Errorf("invalid hex data: %w", err)
	}

	fmt.Println(string(result))
	return nil
}
