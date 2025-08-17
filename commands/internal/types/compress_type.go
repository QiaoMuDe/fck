package types

import (
	"fmt"
	"path/filepath"
	"strings"
)

// 支持的压缩格式
type CompressType string

const (
	CompressTypeZip   CompressType = ".zip"    // zip 压缩格式
	CompressTypeTar   CompressType = ".tar"    // tar 压缩格式
	CompressTypeTgz   CompressType = ".tgz"    // tgz 压缩格式
	CompressTypeTarGz CompressType = ".tar.gz" // tar.gz 压缩格式
	CompressTypeGz    CompressType = ".gz"     // gz 压缩格式
)

// supportedCompressTypes 受支持的压缩格式map, key是压缩格式类型，value是空结构体
var supportedCompressTypes = map[CompressType]struct{}{
	CompressTypeZip:   {}, // zip 压缩格式
	CompressTypeTar:   {}, // tar 压缩格式
	CompressTypeTgz:   {}, // tgz 压缩格式
	CompressTypeTarGz: {}, // tar.gz 压缩格式
	CompressTypeGz:    {}, // gz 压缩格式
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
	// 处理.tar.gz特殊情况
	if strings.HasSuffix(strings.ToLower(filename), ".tar.gz") {
		return CompressTypeTarGz, nil
	}

	// 获取文件扩展名
	ext := filepath.Ext(filename)
	if !IsSupportedCompressType(ext) {
		return "", fmt.Errorf("不支持的压缩文件格式: %s, 支持的格式: %v", ext, SupportedCompressTypes())
	}

	return CompressType(ext), nil
}
