package pack

import (
	"errors"
	"fmt"

	"gitee.com/MM-Q/comprx"
	"gitee.com/MM-Q/fck/commands/internal/types"
)

func PackCmdMain() error {
	// 获取压缩包路径
	packPath := packCmd.Arg(0)
	if packPath == "" {
		return errors.New("压缩包名称不能为空")
	}

	// 获取源路径
	srcPath := packCmd.Arg(1)
	if srcPath == "" {
		return errors.New("源路径不能为空")
	}

	// 过滤器配置
	filter := comprx.FilterOptions{
		Include: includePatterns.Get(), // 包含的文件或目录
		Exclude: excludePatterns.Get(), // 排除的文件或目录
		MinSize: minSize.Get(),         // 最小文件大小
		MaxSize: maxSize.Get(),         // 最大文件大小
	}

	// 获取压缩级别并检查有效性
	compressionLevelVal, isValid := types.GetCompressionLevel(compressionLevel.Get())
	if !isValid {
		return fmt.Errorf("无效的压缩级别: %s", compressionLevel.Get())
	}

	// 获取进度条样式并检查有效性
	progressStyleVal, isValid := types.GetProgressStyle(progressStyle.Get())
	if !isValid {
		return fmt.Errorf("无效的进度条样式: %s", progressStyle.Get())
	}

	// 压缩配置
	opts := comprx.Options{
		CompressionLevel:      compressionLevelVal, // 压缩级别
		OverwriteExisting:     overwrite.Get(),     // 覆盖已存在的文件
		ProgressEnabled:       progress.Get(),      // 显示进度
		ProgressStyle:         progressStyleVal,    // 进度样式
		DisablePathValidation: noValidate.Get(),    // 是否禁用路径验证
		Filter:                filter,              // 过滤器
	}

	// 打包
	if packErr := comprx.PackOptions(packPath, srcPath, opts); packErr != nil {
		return packErr
	}

	return nil
}
