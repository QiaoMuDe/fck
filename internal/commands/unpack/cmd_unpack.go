package unpack

import (
	"errors"
	"fmt"

	"gitee.com/MM-Q/comprx"
	"gitee.com/MM-Q/fck/internal/types"
)

// UnpackConfig 解包配置结构体
type UnpackConfig struct {
	PackPath        string   // 压缩包路径
	DstPath         string   // 目标路径
	IncludePatterns []string // 包含的文件或目录
	ExcludePatterns []string // 排除的文件或目录
	MinSize         int64    // 最小文件大小
	MaxSize         int64    // 最大文件大小
	Overwrite       bool     // 覆盖已存在的文件
	Progress        bool     // 显示进度
	ProgressStyle   string   // 进度样式
	NoValidate      bool     // 是否禁用路径验证
}

// UnpackCmdMain 执行解包命令
//
// 参数:
//   - config: 解包配置
//
// 返回值:
//   - error: 解包过程中可能发生的错误
func UnpackCmdMain(config UnpackConfig) error {
	if config.PackPath == "" {
		return errors.New("压缩包名称不能为空")
	}

	dstPath := config.DstPath
	if dstPath == "" {
		dstPath = "."
	}

	filter := comprx.FilterOptions{
		Include: config.IncludePatterns,
		Exclude: config.ExcludePatterns,
		MinSize: config.MinSize,
		MaxSize: config.MaxSize,
	}

	progressStyleVal, isValid := types.GetProgressStyle(config.ProgressStyle)
	if !isValid {
		return fmt.Errorf("无效的进度条样式: %s", config.ProgressStyle)
	}

	opts := comprx.Options{
		CompressionLevel:      comprx.CompressionLevelDefault,
		OverwriteExisting:     config.Overwrite,
		ProgressEnabled:       config.Progress,
		ProgressStyle:         progressStyleVal,
		DisablePathValidation: config.NoValidate,
		Filter:                filter,
	}

	if unpackErr := comprx.UnpackOptions(config.PackPath, dstPath, opts); unpackErr != nil {
		return unpackErr
	}

	return nil
}
