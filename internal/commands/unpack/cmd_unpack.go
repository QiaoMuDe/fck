package unpack

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gitee.com/MM-Q/comprx"
	"gitee.com/MM-Q/fck/internal/types"
	"gitee.com/MM-Q/go-kit/fs"
)

// UnpackConfig 解包配置结构体
type UnpackConfig struct {
	Archives        []string // 压缩包路径列表（支持通配符）
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
	if len(config.Archives) == 0 {
		return errors.New("no archives specified")
	}

	dstPath := config.DstPath
	if dstPath == "" {
		dstPath = "."
	}

	// 展开压缩包列表（支持通配符）
	archives, err := fs.ExpandFiles(config.Archives)
	if err != nil {
		return err
	}

	// 检查是否有压缩包可处理
	if len(archives) == 0 {
		return errors.New("no archives to process")
	}

	// 清理所有压缩包路径
	for i, archive := range archives {
		archives[i] = filepath.Clean(archive)
	}

	// 统计信息
	stats := struct {
		total   int
		success int
		failed  int
	}{}

	// 处理每个压缩包
	for _, archive := range archives {
		stats.total++

		// 构建单个压缩包的配置
		singleConfig := UnpackConfig{
			Archives:        []string{archive},
			DstPath:         dstPath,
			IncludePatterns: config.IncludePatterns,
			ExcludePatterns: config.ExcludePatterns,
			MinSize:         config.MinSize,
			MaxSize:         config.MaxSize,
			Overwrite:       config.Overwrite,
			Progress:        config.Progress,
			ProgressStyle:   config.ProgressStyle,
			NoValidate:      config.NoValidate,
		}

		if err := unpackSingle(singleConfig); err != nil {
			stats.failed++
			fmt.Fprintf(os.Stderr, "%s: %v\n", archive, err)
			continue
		}

		stats.success++
	}

	// 输出统计信息
	if stats.failed > 0 {
		return fmt.Errorf("unpack completed with %d errors", stats.failed)
	}

	return nil
}

// unpackSingle 解压单个压缩包
func unpackSingle(config UnpackConfig) error {
	if len(config.Archives) != 1 {
		return errors.New("unpackSingle requires exactly one archive")
	}

	// 获取压缩包路径
	archive := config.Archives[0]

	// 构建过滤选项
	filter := comprx.FilterOptions{
		Include: config.IncludePatterns,
		Exclude: config.ExcludePatterns,
		MinSize: config.MinSize,
		MaxSize: config.MaxSize,
	}

	// 构建解包选项
	progressStyleVal, isValid := types.GetProgressStyle(config.ProgressStyle)
	if !isValid {
		return fmt.Errorf("invalid progress style: %s", config.ProgressStyle)
	}
	opts := comprx.Options{
		CompressionLevel:      comprx.CompressionLevelDefault,
		OverwriteExisting:     config.Overwrite,
		ProgressEnabled:       config.Progress,
		ProgressStyle:         progressStyleVal,
		DisablePathValidation: config.NoValidate,
		Filter:                filter,
	}

	// 执行解包
	return comprx.UnpackOptions(archive, config.DstPath, opts)
}
