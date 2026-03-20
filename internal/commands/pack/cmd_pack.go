package pack

import (
	"errors"
	"fmt"

	"gitee.com/MM-Q/comprx"
	"gitee.com/MM-Q/fck/internal/types"
)

// PackConfig 打包配置结构体
type PackConfig struct {
	PackPath         string   // 压缩包路径
	SrcPath          string   // 源路径
	IncludePatterns  []string // 包含的文件或目录
	ExcludePatterns  []string // 排除的文件或目录
	MinSize          int64    // 最小文件大小
	MaxSize          int64    // 最大文件大小
	CompressionLevel string   // 压缩级别
	Overwrite        bool     // 覆盖已存在的文件
	Progress         bool     // 显示进度
	ProgressStyle    string   // 进度样式
	NoValidate       bool     // 是否禁用路径验证
}

// PackCmdMain 执行打包命令
//
// 参数:
//   - config: 打包配置
//
// 返回值:
//   - error: 打包过程中可能发生的错误
func PackCmdMain(config PackConfig) error {
	if config.PackPath == "" {
		return errors.New("压缩包名称不能为空")
	}

	if config.SrcPath == "" {
		return errors.New("源路径不能为空")
	}

	filter := comprx.FilterOptions{
		Include: config.IncludePatterns, // 包含的文件或目录
		Exclude: config.ExcludePatterns, // 排除的文件或目录
		MinSize: config.MinSize,         // 最小文件大小
		MaxSize: config.MaxSize,         // 最大文件大小
	}

	compressionLevelVal, isValid := types.GetCompressionLevel(config.CompressionLevel)
	if !isValid {
		return fmt.Errorf("无效的压缩级别: %s", config.CompressionLevel)
	}

	progressStyleVal, isValid := types.GetProgressStyle(config.ProgressStyle)
	if !isValid {
		return fmt.Errorf("无效的进度条样式: %s", config.ProgressStyle)
	}

	opts := comprx.Options{
		CompressionLevel:      compressionLevelVal, // 压缩级别
		OverwriteExisting:     config.Overwrite,    // 是否覆盖已存在的文件
		ProgressEnabled:       config.Progress,     // 显示进度
		ProgressStyle:         progressStyleVal,    //  进度条样式
		DisablePathValidation: config.NoValidate,   // 是否禁用路径验证
		Filter:                filter,              //  过滤选项
	}

	if packErr := comprx.PackOptions(config.PackPath, config.SrcPath, opts); packErr != nil {
		return packErr
	}

	return nil
}
