package rm

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// RMConfig 配置结构体
type RMConfig struct {
	Targets   []string
	Recursive bool
	Force     bool
	Verbose   bool
}

// DeleteStats 删除统计
type DeleteStats struct {
	Files       int
	Directories int
	Errors      int
}

var confirmAll = false

// RMCmdMain 主函数
func RMCmdMain(config RMConfig) error {
	if len(config.Targets) == 0 {
		return fmt.Errorf("no targets specified")
	}

	stats := &DeleteStats{}

	for _, target := range config.Targets {
		err := deletePath(target, config, stats)
		if err != nil {
			stats.Errors++
			if !config.Force {
				return err
			}
		}
	}

	if config.Verbose {
		fmt.Printf("delete completed: %d files, %d directories", stats.Files, stats.Directories)
		if stats.Errors > 0 {
			fmt.Printf(", %d errors", stats.Errors)
		}
		fmt.Println()
	}

	return nil
}

// deletePath 删除路径
func deletePath(path string, config RMConfig, stats *DeleteStats) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	if err := validatePath(absPath); err != nil {
		return err
	}

	info, err := os.Lstat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			if config.Verbose {
				fmt.Printf("skip: %s (not exists)\n", path)
			}
			return nil
		}
		return fmt.Errorf("failed to access path: %w", err)
	}

	if !config.Force && !confirmDelete(path, getFileType(info)) {
		if config.Verbose {
			fmt.Printf("skip: %s\n", path)
		}
		return nil
	}

	if err := os.RemoveAll(absPath); err != nil {
		return fmt.Errorf("failed to delete path: %w", err)
	}

	if config.Verbose {
		fmt.Printf("delete: %s\n", path)
	}

	if info.IsDir() {
		stats.Directories++
	} else {
		stats.Files++
	}

	return nil
}

// getFileType 获取文件类型
func getFileType(info os.FileInfo) string {
	if info.IsDir() {
		return "directory"
	}
	return "file"
}

// validatePath 验证路径
func validatePath(path string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	cwdAbs, err := filepath.Abs(cwd)
	if err != nil {
		return fmt.Errorf("failed to get absolute path of current working directory: %w", err)
	}

	if path == cwdAbs {
		return fmt.Errorf("refuse to delete current working directory")
	}

	if filepath.VolumeName(path) == path {
		return fmt.Errorf("refuse to delete root directory")
	}

	return nil
}

// confirmDelete 确认删除
func confirmDelete(path, fileType string) bool {
	if confirmAll {
		return true
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("delete %s %s? (y/n/a/q): ", fileType, path)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		switch input {
		case "y":
			return true
		case "n":
			return false
		case "a":
			confirmAll = true
			return true
		case "q":
			os.Exit(0)
		default:
			fmt.Println("invalid option, please input y/n/a/q")
		}
	}
}
