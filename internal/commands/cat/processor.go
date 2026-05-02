package cat

import (
	"bytes"
	"fmt"
)

// Processor 内容处理器
type Processor struct {
	config CatConfig
}

// NewProcessor 创建内容处理器
//
// 参数:
//   - config: 命令配置
//
// 返回:
//   - *Processor: 内容处理器实例
func NewProcessor(config CatConfig) *Processor {
	return &Processor{config: config}
}

// Process 处理内容源
//
// 参数:
//   - source: 内容源
//
// 返回:
//   - []byte: 处理后的内容
//   - error: 错误
func (p *Processor) Process(source ContentSource) ([]byte, error) {
	// 1. 读取内容 (FileSource 和 StdinSource 已在读取时进行大小限制)
	content, err := source.Read()
	if err != nil {
		return nil, err
	}

	// 2. 二进制检测（文件源）
	if !p.config.Text {
		isBinary, err := source.IsBinary()
		if err != nil {
			return nil, fmt.Errorf("cannot detect content type for %s: %w", source.Name(), err)
		}
		if isBinary {
			// -I 模式：静默跳过
			if p.config.IgnoreBinary {
				return nil, nil
			}
			// 默认行为：输出提示并跳过
			if !p.config.Quiet {
				fmt.Printf("bin file %s matches\n", source.Name())
			}
			return nil, nil
		}
	}

	// 3. 统一换行符
	content = normalizeLineEndings(content)

	return content, nil
}

// normalizeLineEndings 统一换行符为 \n
//
// 参数:
//   - content: 原始内容
//
// 返回:
//   - []byte: 换行符统一后的内容
func normalizeLineEndings(content []byte) []byte {
	// 将 \r\n 替换为 \n
	content = bytes.ReplaceAll(content, []byte("\r\n"), []byte("\n"))
	// 将剩余的 \r 替换为 \n
	content = bytes.ReplaceAll(content, []byte("\r"), []byte("\n"))
	return content
}

// SplitLines 将内容按行分割
//
// 参数:
//   - content: 内容
//
// 返回:
//   - [][]byte: 按行分割后的内容切片
func SplitLines(content []byte) [][]byte {
	if len(content) == 0 {
		return nil
	}
	return bytes.Split(content, []byte("\n"))
}
