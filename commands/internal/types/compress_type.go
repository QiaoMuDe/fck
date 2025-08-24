package types

import (
	"gitee.com/MM-Q/comprx/types"
)

// 进度条样式
const (
	ProgressStyleText    = "text"    // 文本样式
	ProgressStyleDefault = "default" // 默认样式
	ProgressStyleUnicode = "unicode" // unicode样式
	ProgressStyleAscii   = "ascii"   // ascii样式
)

// 受支持的进度条样式
var SupportedProgressStyles = []string{
	ProgressStyleText,
	ProgressStyleDefault,
	ProgressStyleUnicode,
	ProgressStyleAscii,
}

// 进度条样式映射
var ProgressStyleMap = map[string]types.ProgressStyle{
	ProgressStyleText:    types.ProgressStyleText,    // 文本
	ProgressStyleDefault: types.ProgressStyleDefault, // 默认
	ProgressStyleUnicode: types.ProgressStyleUnicode, // unicode
	ProgressStyleAscii:   types.ProgressStyleASCII,   // ascii
}

// 压缩级别
const (
	CompressionLevelDefault = "default"      // 默认压缩级别
	CompressionLevelNone    = "none"         // 不压缩
	CompressionLevelFast    = "fast"         // 快速压缩
	CompressionLevelBest    = "best"         // 最佳压缩
	CompressionLevelHuffman = "huffman-only" // 仅使用霍夫曼编码
)

// 受支持的压缩级别
var SupportedCompressionLevels = []string{
	CompressionLevelDefault,
	CompressionLevelNone,
	CompressionLevelFast,
	CompressionLevelBest,
	CompressionLevelHuffman,
}

// 压缩级别映射
var CompressionLevelMap = map[string]types.CompressionLevel{
	CompressionLevelDefault: types.CompressionLevelDefault,     // 默认
	CompressionLevelNone:    types.CompressionLevelNone,        // 不压缩
	CompressionLevelFast:    types.CompressionLevelFast,        // 快速
	CompressionLevelBest:    types.CompressionLevelBest,        // 最佳
	CompressionLevelHuffman: types.CompressionLevelHuffmanOnly, //  Huffman
}

// GetCompressionLevel 获取压缩级别，如果无效则返回默认级别
//
// 参数:
//   - level: 压缩级别字符串
//
// 返回值:
//   - types.CompressionLevel: 压缩级别枚举值
//   - bool: 是否成功获取到压缩级别
func GetCompressionLevel(level string) (types.CompressionLevel, bool) {
	compressionLevel, ok := CompressionLevelMap[level]
	if !ok {
		return types.CompressionLevelDefault, false
	}
	return compressionLevel, true
}

// GetProgressStyle 获取进度条样式，如果无效则返回默认样式
//
// 参数:
//   - style: 进度条样式字符串
//
// 返回值:
//   - types.ProgressStyle: 进度条样式枚举值
//   - bool: 是否成功获取到进度条样式
func GetProgressStyle(style string) (types.ProgressStyle, bool) {
	progressStyle, ok := ProgressStyleMap[style]
	if !ok {
		return types.ProgressStyleDefault, false
	}
	return progressStyle, true
}
