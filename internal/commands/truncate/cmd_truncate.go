package truncate

import (
	"fmt"
	"os"
)

// TruncateConfig 配置结构体
type TruncateConfig struct {
	Size      int64
	Reference string
	Targets   []string
	Create    bool
	Verbose   bool
}

// TruncateStats 操作统计
type TruncateStats struct {
	Modified int
	Created  int
	Errors   int
}

// TruncateCmdMain 主函数
func TruncateCmdMain(config TruncateConfig) error {
	if len(config.Targets) == 0 {
		return fmt.Errorf("未指定要截断的文件")
	}

	size := config.Size
	var err error

	if config.Reference != "" {
		size, err = getReferenceSize(config.Reference)
		if err != nil {
			return err
		}
	}

	if size < 0 {
		return fmt.Errorf("文件大小不能为负数")
	}

	stats := &TruncateStats{}

	for _, target := range config.Targets {
		err := truncateFile(target, size, config.Create, config.Verbose, stats)
		if err != nil {
			stats.Errors++
			return err
		}
	}

	if config.Verbose {
		fmt.Printf("操作完成: %d 个创建, %d 个修改", stats.Created, stats.Modified)
		if stats.Errors > 0 {
			fmt.Printf(", %d 个错误", stats.Errors)
		}
		fmt.Println()
	}

	return nil
}

// truncateFile 截断文件
func truncateFile(path string, size int64, create bool, verbose bool, stats *TruncateStats) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if !create {
			return fmt.Errorf("文件不存在: %s", path)
		}

		file, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("创建文件失败: %w", err)
		}
		_ = file.Close()

		if verbose {
			fmt.Printf("创建文件: %s (%d)\n", path, size)
		}

		stats.Created++
		return nil
	}

	if err != nil {
		return fmt.Errorf("获取文件信息失败: %w", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("获取文件信息失败: %w", err)
	}
	currentSize := info.Size()

	if currentSize == size {
		if verbose {
			fmt.Printf("跳过: %s (大小已为 %d)\n", path, size)
		}
		return nil
	}

	if err := os.Truncate(path, size); err != nil {
		return fmt.Errorf("截断文件失败: %w", err)
	}

	if verbose {
		fmt.Printf("截断文件: %s (%d -> %d)\n", path, currentSize, size)
	}

	stats.Modified++
	return nil
}

// getReferenceSize 获取参考文件大小
func getReferenceSize(refPath string) (int64, error) {
	info, err := os.Stat(refPath)
	if err != nil {
		return 0, fmt.Errorf("获取参考文件信息失败: %w", err)
	}
	return info.Size(), nil
}
