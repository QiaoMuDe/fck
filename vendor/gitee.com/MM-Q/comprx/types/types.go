// Package types 定义了 comprx 库使用的核心数据类型和常量。
//
// 该包提供了压缩库所需的基础类型定义，包括进度条样式、压缩格式、
// 压缩等级等核心类型，以及相关的验证和转换方法。
//
// 主要类型：
//   - ProgressStyle: 进度条样式类型
//   - CompressType: 压缩格式类型
//   - CompressionLevel: 压缩等级类型
//
// 主要功能：
//   - 定义支持的压缩格式和进度条样式
//   - 提供压缩格式自动检测功能
//   - 提供类型验证和转换方法
//   - 定义压缩等级常量和验证
package types

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ProgressStyle 进度条样式类型
//
// 进度条样式类型定义:
//   - ProgressStyleText: 文本样式进度条 - 使用文字描述进度
//   - ProgressStyleUnicode: Unicode样式进度条 - 使用Unicode字符绘制精美进度条
//   - ProgressStyleASCII: ASCII样式进度条 - 使用基础ASCII字符绘制兼容性最好的进度条
type ProgressStyle string

// 进度条样式常量
const (
	// ProgressStyleText 文本样式进度条 - 使用文字描述进度
	ProgressStyleText ProgressStyle = "text"

	// ProgressStyleDefault 默认进度条样式 - progress库的默认进度条样式
	ProgressStyleDefault ProgressStyle = "default"

	// ProgressStyleUnicode Unicode样式进度条 - 使用Unicode字符绘制精美进度条
	// 示例: ████████████░░░░░░░░ 60%
	ProgressStyleUnicode ProgressStyle = "unicode"

	// ProgressStyleASCII ASCII样式进度条 - 使用基础ASCII字符绘制兼容性最好的进度条
	// 示例: [##########          ] 50%
	ProgressStyleASCII ProgressStyle = "ascii"
)

// String 返回进度条样式的字符串表示
func (ps ProgressStyle) String() string {
	return string(ps)
}

// IsValid 检查进度条样式是否有效
func (ps ProgressStyle) IsValid() bool {
	switch ps {
	case ProgressStyleText, ProgressStyleDefault, ProgressStyleUnicode, ProgressStyleASCII:
		return true
	default:
		return false
	}
}

// SupportedProgressStyles 返回所有支持的进度条样式
func SupportedProgressStyles() []ProgressStyle {
	return []ProgressStyle{
		ProgressStyleText,
		ProgressStyleDefault,
		ProgressStyleUnicode,
		ProgressStyleASCII,
	}
}

// 支持的压缩格式
//
// 压缩格式类型定义:
//   - CompressTypeZip: zip 压缩格式
//   - CompressTypeTar: tar 压缩格式
//   - CompressTypeTgz: tgz 压缩格式
//   - CompressTypeTarGz: tar.gz 压缩格式
//   - CompressTypeGz: gz 压缩格式
//   - CompressTypeBz2: bz2 压缩格式
//   - CompressTypeBzip2: bzip2 压缩格式
type CompressType string

const (
	CompressTypeZip   CompressType = ".zip"    // zip 压缩格式
	CompressTypeTar   CompressType = ".tar"    // tar 压缩格式
	CompressTypeTgz   CompressType = ".tgz"    // tgz 压缩格式
	CompressTypeTarGz CompressType = ".tar.gz" // tar.gz 压缩格式
	CompressTypeGz    CompressType = ".gz"     // gz 压缩格式
	CompressTypeBz2   CompressType = ".bz2"    // bz2 压缩格式
	CompressTypeBzip2 CompressType = ".bzip2"  // bzip2 压缩格式
	CompressTypeZlib  CompressType = ".zlib"   // zlib 压缩格式
)

// supportedCompressTypes 受支持的压缩格式map, key是压缩格式类型，value是空结构体
var supportedCompressTypes = map[CompressType]struct{}{
	CompressTypeZip:   {}, // zip 压缩格式
	CompressTypeTar:   {}, // tar 压缩格式
	CompressTypeTgz:   {}, // tgz 压缩格式
	CompressTypeTarGz: {}, // tar.gz 压缩格式
	CompressTypeGz:    {}, // gz 压缩格式
	CompressTypeBz2:   {}, // bz2 压缩格式
	CompressTypeBzip2: {}, // bzip2 压缩格式
	CompressTypeZlib:  {}, // zlib 压缩格式
}

// String 压缩格式的字符串表示
//
// 返回:
//   - string: 压缩格式的字符串表示
func (c CompressType) String() string {
	return string(c)
}

// IsSupportedCompressType 判断是否受支持的压缩格式
//
// 参数:
//   - ct: 压缩格式字符串
//
// 返回:
//   - bool: 如果是受支持的压缩格式, 返回 true, 否则返回 false
func IsSupportedCompressType(ct string) bool {
	_, ok := supportedCompressTypes[CompressType(ct)]
	return ok
}

// SupportedCompressTypes 返回受支持的压缩格式字符串列表
//
// 返回:
//   - []string: 受支持的压缩格式字符串列表
func SupportedCompressTypes() []string {
	var compressTypes []string
	for ct := range supportedCompressTypes {
		compressTypes = append(compressTypes, ct.String())
	}
	return compressTypes
}

// DetectCompressFormat 智能检测压缩文件格式
//
// 参数:
//   - filename: 文件名
//
// 返回:
//   - types.CompressType: 检测到的压缩格式
//   - error: 错误信息
func DetectCompressFormat(filename string) (CompressType, error) {
	// 转换为小写进行处理
	lowerFilename := strings.ToLower(filename)

	// 处理.tar.gz特殊情况
	if strings.HasSuffix(lowerFilename, ".tar.gz") {
		return CompressTypeTarGz, nil
	}

	// 获取文件扩展名并转换为小写
	ext := strings.ToLower(filepath.Ext(filename))
	if !IsSupportedCompressType(ext) {
		return "", fmt.Errorf("不支持的压缩文件格式: %s, 支持的格式: %v", ext, SupportedCompressTypes())
	}

	return CompressType(ext), nil
}

// CompressionLevel 压缩等级类型
//
// 支持的压缩等级：
//   - CompressionLevelDefault: 默认压缩等级(zip仅支持该等级)
//   - CompressionLevelNone: 不压缩(zip仅支持该等级)
//   - CompressionLevelFast: 快速压缩
//   - CompressionLevelBest: 最佳压缩
//   - CompressionLevelHuffmanOnly: 仅使用Huffman编码
type CompressionLevel int

const (
	// 压缩等级常量
	CompressionLevelDefault     CompressionLevel = -1 // 默认压缩等级(zip仅支持该等级)
	CompressionLevelNone        CompressionLevel = 0  // 不压缩(zip仅支持该等级)
	CompressionLevelFast        CompressionLevel = 1  // 快速压缩
	CompressionLevelBest        CompressionLevel = 9  // 最佳压缩
	CompressionLevelHuffmanOnly CompressionLevel = -2 // 仅使用Huffman编码
)

// IsValid 检查压缩等级是否有效
// 有效范围：-2 到 9
func (cl CompressionLevel) IsValid() bool {
	return cl >= CompressionLevelHuffmanOnly && cl <= CompressionLevelBest
}

// String 返回压缩等级的字符串表示
func (cl CompressionLevel) String() string {
	switch cl {
	case CompressionLevelHuffmanOnly:
		return "huffman-only"
	case CompressionLevelDefault:
		return "default"
	case CompressionLevelNone:
		return "none"
	case CompressionLevelFast:
		return "fast"
	case CompressionLevelBest:
		return "best"
	default:
		if cl >= 0 && cl <= 9 {
			return fmt.Sprintf("level-%d", int(cl))
		}
		return fmt.Sprintf("invalid-%d", int(cl))
	}
}

// SupportedCompressionLevels 返回所有预定义的压缩等级
func SupportedCompressionLevels() []CompressionLevel {
	return []CompressionLevel{
		CompressionLevelHuffmanOnly,
		CompressionLevelDefault,
		CompressionLevelNone,
		CompressionLevelFast,
		CompressionLevelBest,
	}
}
