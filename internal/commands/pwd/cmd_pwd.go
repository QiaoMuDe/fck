package pwd

import (
	"fmt"
	"os"
	"path/filepath"
)

// PwdConfig 配置结构体
type PwdConfig struct {
	Physical bool
	Logical  bool
}

// PwdCmdMain 主函数
func PwdCmdMain(config PwdConfig) error {
	path, err := processPath(config)
	if err != nil {
		return err
	}

	fmt.Println(path)
	return nil
}

// getWorkingDirectory 获取当前工作目录
func getWorkingDirectory() (string, error) {
	return os.Getwd()
}

// resolveSymlinks 解析符号链接，获取物理路径
func resolveSymlinks(path string) (string, error) {
	return filepath.EvalSymlinks(path)
}

// processPath 处理路径
func processPath(config PwdConfig) (string, error) {
	cwd, err := getWorkingDirectory()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	if config.Physical {
		resolved, err := resolveSymlinks(cwd)
		if err != nil {
			return "", fmt.Errorf("failed to resolve symlinks: %w", err)
		}
		return resolved, nil
	}

	return cwd, nil
}
