package preview

import (
	"errors"

	"gitee.com/MM-Q/comprx"
)

// PreviewConfig 预览配置结构体
type PreviewConfig struct {
	PackPath string // 压缩包路径
	Info     bool   // 打印压缩包信息
	Ls       bool   // 简洁文件列表
	Ll       bool   // 详细文件列表
	Limit    int    // 限制文件数量
}

// PreviewCmdMain 执行预览命令
//
// 参数:
//   - config: 预览配置
//
// 返回值:
//   - error: 预览过程中可能发生的错误
func PreviewCmdMain(config PreviewConfig) error {
	if config.PackPath == "" {
		return errors.New("no pack path specified")
	}

	if config.Info {
		if err := comprx.PrintArchiveInfo(config.PackPath); err != nil {
			return err
		}
		return nil
	}

	if config.Ls {
		if config.Limit != 0 {
			if err := comprx.PrintFilesLimit(config.PackPath, config.Limit, false); err != nil {
				return err
			}
			return nil
		}

		if err := comprx.PrintLs(config.PackPath); err != nil {
			return err
		}
		return nil
	}

	if config.Ll {
		if config.Limit != 0 {
			if err := comprx.PrintFilesLimit(config.PackPath, config.Limit, true); err != nil {
				return err
			}
			return nil
		}

		if err := comprx.PrintLl(config.PackPath); err != nil {
			return err
		}
		return nil
	}

	if err := comprx.PrintArchiveAndFiles(config.PackPath, true); err != nil {
		return err
	}

	return nil
}
