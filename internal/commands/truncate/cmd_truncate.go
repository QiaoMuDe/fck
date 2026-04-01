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
		return fmt.Errorf("no target files specified")
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
		return fmt.Errorf("file size cannot be negative")
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
		fmt.Printf("operation completed: %d created files, %d modified files", stats.Created, stats.Modified)
		if stats.Errors > 0 {
			fmt.Printf(", %d errors", stats.Errors)
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
			return fmt.Errorf("file does not exist: %s", path)
		}

		file, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("create file failed, error: %w", err)
		}
		_ = file.Close()

		if verbose {
			fmt.Printf("create file: %s (%d)\n", path, size)
		}

		stats.Created++
		return nil
	}

	if err != nil {
		return fmt.Errorf("error checking file: %w", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("error checking file: %w", err)
	}
	currentSize := info.Size()

	if currentSize == size {
		if verbose {
			fmt.Printf("skip: %s (already size %d)\n", path, size)
		}
		return nil
	}

	if err := os.Truncate(path, size); err != nil {
		return fmt.Errorf("truncate file failed, error: %w", err)
	}

	if verbose {
		fmt.Printf("truncate file: %s (%d -> %d)\n", path, currentSize, size)
	}

	stats.Modified++
	return nil
}

// getReferenceSize 获取参考文件大小
func getReferenceSize(refPath string) (int64, error) {
	info, err := os.Stat(refPath)
	if err != nil {
		return 0, fmt.Errorf("error checking reference file: %w", err)
	}
	return info.Size(), nil
}
