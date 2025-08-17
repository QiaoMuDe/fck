package compress

import (
	"fmt"
	"path/filepath"

	"gitee.com/MM-Q/fck/commands/internal/types"
	"gitee.com/MM-Q/qflag"
)

// 压缩配置
type compressConfig struct {
	compressFilePath string             // 压缩文件路径
	sourceFilePaths  []string           // 压缩源文件路径列表
	compressType     types.CompressType // 压缩格式
}

// newCompressConfig 创建一个压缩配置
//
// 返回:
//   - *compressConfig: 压缩配置
//   - error: 错误信息
func newCompressConfig() (*compressConfig, error) {
	// 获取命令行参数
	compressCmdArgs := compressCmd.Args()

	// 判断参数是否大于2
	if len(compressCmdArgs) < 2 {
		return nil, fmt.Errorf("参数错误, 示例: %s %s <压缩文件> <源文件...>", qflag.LongName(), compressCmd.LongName())
	}

	// 获取压缩文件路径
	compressFilePath := compressCmdArgs[0]

	// 获取源文件路径列表
	sourceFilePaths := compressCmdArgs[1:]

	// 获取压缩文件后缀
	compressFileExt := filepath.Ext(compressFilePath)

	// 检查是否为受支持的压缩文件类型
	if !types.IsSupportedCompressType(compressFileExt) {
		return nil, fmt.Errorf("不支持的压缩文件类型: %s, 支持的类型: %v", compressFileExt, types.SupportedCompressTypes())
	}

	// 创建压缩配置
	config := &compressConfig{
		compressFilePath: compressFilePath,                    // 压缩文件路径
		sourceFilePaths:  sourceFilePaths,                     // 压缩源文件路径列表
		compressType:     types.CompressType(compressFileExt), // 压缩格式
	}

	return config, nil
}
