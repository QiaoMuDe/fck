package extract

import (
	"fmt"

	"gitee.com/MM-Q/fck/commands/internal/types"
	"gitee.com/MM-Q/qflag"
)

// 解压配置
type extractConfig struct {
	archiveFilePath string             // 压缩文件路径
	destPath        string             // 解压目标路径
	extractType     types.CompressType // 压缩格式
}

// newExtractConfig 创建一个解压配置
//
// 返回:
//   - *extractConfig: 解压配置
//   - error: 错误信息
func newExtractConfig() (*extractConfig, error) {
	// 获取命令行参数
	extractCmdArgs := extractCmd.Args()

	// 判断参数是否至少有1个
	if len(extractCmdArgs) < 1 {
		return nil, fmt.Errorf("参数错误, 示例: %s %s <压缩文件> [目标目录]", qflag.LongName(), extractCmd.LongName())
	}

	// 获取压缩文件路径
	archiveFilePath := extractCmdArgs[0]

	// 获取目标路径，如果未指定则使用当前目录
	destPath := "."
	if len(extractCmdArgs) > 1 {
		destPath = extractCmdArgs[1]
	}

	// 智能检测压缩文件格式
	extractType, err := types.DetectCompressFormat(archiveFilePath)
	if err != nil {
		return nil, fmt.Errorf("无法识别压缩文件格式: %s, 错误: %v", archiveFilePath, err)
	}

	// 创建解压配置
	config := &extractConfig{
		archiveFilePath: archiveFilePath, // 压缩文件路径
		destPath:        destPath,        // 解压目标路径
		extractType:     extractType,     // 压缩格式
	}

	return config, nil
}
