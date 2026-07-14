package preview

import (
	"errors"
	"fmt"

	"gitee.com/MM-Q/comprx"
)

// PreviewConfig 预览配置结构体
type PreviewConfig struct {
	PackPaths []string // 压缩包路径列表
	Info      bool     // 打印压缩包信息
	Ls        bool     // 简洁文件列表
	Ll        bool     // 详细文件列表
	Limit     int      // 限制文件数量
}

// previewSingle 预览单个压缩包
//
// 参数:
//   - path: 压缩包路径
//   - config: 预览配置
//
// 返回值:
//   - error: 预览过程中可能发生的错误
func previewSingle(path string, config PreviewConfig) error {
	if config.Info {
		return comprx.PrintArchiveInfo(path)
	}

	if config.Ls {
		if config.Limit != 0 {
			return comprx.PrintFilesLimit(path, config.Limit, false)
		}
		return comprx.PrintLs(path)
	}

	if config.Ll {
		if config.Limit != 0 {
			return comprx.PrintFilesLimit(path, config.Limit, true)
		}
		return comprx.PrintLl(path)
	}

	return comprx.PrintArchiveAndFiles(path, true)
}

// PreviewCmdMain 执行预览命令
//
// 参数:
//   - config: 预览配置
//
// 返回值:
//   - error: 预览过程中可能发生的错误
func PreviewCmdMain(config PreviewConfig) error {
	if len(config.PackPaths) == 0 {
		return errors.New("no pack path specified")
	}

	for i, path := range config.PackPaths {
		if i > 0 {
			fmt.Println()
		}
		if len(config.PackPaths) > 1 {
			fmt.Printf("==> %s <==\n", path)
		}

		if err := previewSingle(path, config); err != nil {
			return err
		}
	}

	return nil
}
