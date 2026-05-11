// Package iconv 提供文件编码检测和转换功能
package iconv

import (
	"unicode/utf8"

	"gitee.com/MM-Q/fck/internal/types"
)

// DetectEncoding 自动检测文件编码
//
// 参数:
//   - data: 文件内容样本（前 8KB）
//
// 返回值:
//   - string: 检测到的编码名称
//   - float64: 置信度（0-1）
//   - error: 检测错误（如果有）
func DetectEncoding(data []byte) (string, float64, error) {
	// 1. 检查 UTF-8 BOM
	if hasBOM(data, types.EncodingUTF8) {
		return types.EncodingUTF8BOM, 1.0, nil
	}

	// 2. 检查 UTF-16 BOM
	if hasBOM(data, types.EncodingUTF16LE) {
		return types.EncodingUTF16LE, 1.0, nil
	}
	if hasBOM(data, types.EncodingUTF16BE) {
		return types.EncodingUTF16BE, 1.0, nil
	}

	// 3. 检查是否为纯 ASCII（也是有效的 UTF-8）
	if isASCII(data) {
		return types.EncodingUTF8, 1.0, nil
	}

	// 4. 检查是否为有效的 UTF-8
	if utf8.Valid(data) {
		return types.EncodingUTF8, 0.95, nil
	}

	// 5. 使用启发式算法检测中文编码（GBK/Big5）
	// 统计字节分布特征
	encoding, confidence := detectChineseEncoding(data)
	if confidence > 0.8 {
		return encoding, confidence, nil
	}

	// 6. 默认返回 GBK（Windows 中文环境最常见）
	return types.EncodingGBK, 0.5, nil
}

// hasBOM 检查数据是否以指定编码的 BOM 开头
func hasBOM(data []byte, encoding string) bool {
	switch encoding {
	case types.EncodingUTF8:
		return len(data) >= 3 &&
			data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF
	case types.EncodingUTF16LE:
		return len(data) >= 2 &&
			data[0] == 0xFF && data[1] == 0xFE
	case types.EncodingUTF16BE:
		return len(data) >= 2 &&
			data[0] == 0xFE && data[1] == 0xFF
	default:
		return false
	}
}

// isASCII 检查数据是否为纯 ASCII
func isASCII(data []byte) bool {
	for _, b := range data {
		if b > 0x7F {
			return false
		}
	}
	return true
}

// detectChineseEncoding 检测中文编码（GBK/Big5）
// 使用启发式算法统计字节分布特征
func detectChineseEncoding(data []byte) (string, float64) {
	// GBK 和 Big5 的常用字符范围统计
	gbkCount := 0
	big5Count := 0
	totalMultiByte := 0

	for i := 0; i < len(data); {
		if data[i] < 0x80 {
			// ASCII 字符
			i++
			continue
		}

		// 多字节字符
		if i+1 >= len(data) {
			break
		}

		b1, b2 := data[i], data[i+1]
		totalMultiByte++

		// GBK 范围: 0x81-0xFE 开头，第二个字节 0x40-0xFE
		if b1 >= 0x81 && b1 <= 0xFE {
			if b2 >= 0x40 && b2 <= 0xFE {
				gbkCount++
			}
		}

		// Big5 范围: 0x81-0xFE 开头，第二个字节 0x40-0x7E 或 0xA1-0xFE
		if b1 >= 0x81 && b1 <= 0xFE {
			if (b2 >= 0x40 && b2 <= 0x7E) || (b2 >= 0xA1 && b2 <= 0xFE) {
				big5Count++
			}
		}

		i += 2
	}

	if totalMultiByte == 0 {
		// 没有多字节字符，可能是纯 ASCII
		return types.EncodingUTF8, 1.0
	}

	// 计算置信度
	gbkRatio := float64(gbkCount) / float64(totalMultiByte)
	big5Ratio := float64(big5Count) / float64(totalMultiByte)

	if gbkRatio > 0.9 && gbkRatio > big5Ratio {
		return types.EncodingGBK, gbkRatio
	}
	if big5Ratio > 0.9 && big5Ratio > gbkRatio {
		return types.EncodingBig5, big5Ratio
	}

	// 无法确定，默认 GBK（Windows 中文环境最常见）
	return types.EncodingGBK, 0.5
}
