package unpack

import (
	"errors"
	"fmt"

	"gitee.com/MM-Q/comprx"
	cxtypes "gitee.com/MM-Q/comprx/types"
	"gitee.com/MM-Q/fck/commands/internal/types"
)

// UnpackCmdMain unpack命令的主函数
func UnpackCmdMain() error {
	// 获取压缩包路径
	packPath := unpackCmd.Arg(0)
	if packPath == "" {
		return errors.New("压缩包名称不能为空")
	}

	// 获取目标路径
	dstPath := unpackCmd.Arg(1)
	if dstPath == "" {
		dstPath = "."
	}

	// 过滤器配置
	filter := cxtypes.FilterOptions{
		Include: includePatterns.Get(), // 包含的文件或目录
		Exclude: excludePatterns.Get(), // 排除的文件或目录
		MinSize: minSize.Get(),         // 最小文件大小
		MaxSize: maxSize.Get(),         // 最大文件大小
	}

	// 获取进度条样式并检查有效性
	progressStyleVal, isValid := types.GetProgressStyle(progressStyle.Get())
	if !isValid {
		return fmt.Errorf("无效的进度条样式: %s", progressStyle.Get())
	}

	// 压缩配置
	opts := comprx.Options{
		CompressionLevel:      cxtypes.CompressionLevelDefault, // 压缩级别（解压时使用默认值）
		OverwriteExisting:     overwrite.Get(),                 // 覆盖已存在的文件
		ProgressEnabled:       progress.Get(),                  // 显示进度
		ProgressStyle:         progressStyleVal,                // 进度样式
		DisablePathValidation: noValidate.Get(),                // 是否禁用路径验证
		Filter:                filter,                          // 过滤器
	}

	// 解压
	if unpackErr := comprx.UnpackOptions(packPath, dstPath, opts); unpackErr != nil {
		return unpackErr
	}

	return nil
}
