package pack

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

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
	Timestamp        bool     // 在压缩包名中添加时间戳
}

// PackCmdMain 执行打包命令
//
// 参数:
//   - config: 打包配置
//
// 返回值:
//   - error: 打包过程中可能发生的错误
func PackCmdMain(config PackConfig) error {
	if config.SrcPath == "" {
		return errors.New("no source path specified")
	}

	// 确定压缩包路径
	packPath := config.PackPath
	if packPath == "" {
		packPath = generateDefaultPackPath(config.SrcPath)
	}

	// 如果启用了时间戳，添加时间戳到压缩包名
	if config.Timestamp {
		packPath = addTimestampToPackPath(packPath)
	}

	filter := comprx.FilterOptions{
		Include: config.IncludePatterns, // 包含的文件或目录
		Exclude: config.ExcludePatterns, // 排除的文件或目录
		MinSize: config.MinSize,         // 最小文件大小
		MaxSize: config.MaxSize,         // 最大文件大小
	}

	compressionLevelVal, isValid := types.GetCompressionLevel(config.CompressionLevel)
	if !isValid {
		return fmt.Errorf("invalid compression level: %s", config.CompressionLevel)
	}

	progressStyleVal, isValid := types.GetProgressStyle(config.ProgressStyle)
	if !isValid {
		return fmt.Errorf("invalid progress style: %s", config.ProgressStyle)
	}

	opts := comprx.Options{
		CompressionLevel:      compressionLevelVal, // 压缩级别
		OverwriteExisting:     config.Overwrite,    // 是否覆盖已存在的文件
		ProgressEnabled:       config.Progress,     // 显示进度
		ProgressStyle:         progressStyleVal,    //  进度条样式
		DisablePathValidation: config.NoValidate,   // 是否禁用路径验证
		Filter:                filter,              //  过滤选项
	}

	if packErr := comprx.PackOptions(packPath, config.SrcPath, opts); packErr != nil {
		return packErr
	}

	return nil
}

// generateDefaultPackPath 根据源路径生成默认压缩包路径
//
// 参数:
//   - srcPath: 源路径
//
// 返回值:
//   - string: 默认压缩包路径
func generateDefaultPackPath(srcPath string) string {
	// 获取文件名（不含目录）
	base := filepath.Base(srcPath)

	// 去掉现有扩展名
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	// 添加 .zip 后缀
	return name + ".zip"
}

// addTimestampToPackPath 在压缩包路径中添加时间戳
// 格式: 名称_YYYYMMDDHHMMSS.扩展名
//
// 参数:
//   - packPath: 压缩包路径
//
// 返回值:
//   - string: 添加时间戳后的压缩包路径
func addTimestampToPackPath(packPath string) string {
	// 获取目录和文件名
	dir := filepath.Dir(packPath)
	base := filepath.Base(packPath)

	// 分离文件名和扩展名
	// 处理多段扩展名如 .tar.gz
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	// 检查是否还有扩展名（如 .tar.gz）
	if ext != "" {
		// 检查是否是 .tar 前缀的扩展名
		if strings.HasSuffix(name, ".tar") {
			ext = ".tar" + ext
			name = strings.TrimSuffix(name, ".tar")
		}
	}

	// 生成时间戳
	timestamp := time.Now().Format("20060102150405")

	// 组装新文件名
	newBase := fmt.Sprintf("%s_%s%s", name, timestamp, ext)

	// 如果原路径包含目录，拼接回去
	if dir == "." || dir == "" {
		return newBase
	}
	return filepath.Join(dir, newBase)
}
