package home

import (
	"fmt"
	"os"
)

// HomeConfig 配置结构体
type HomeConfig struct {
}

// HomeCmdMain 主函数
func HomeCmdMain(config HomeConfig) error {
	homeDir, err := getUserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home dir: %w", err)
	}

	fmt.Println(homeDir)
	return nil
}

// getUserHomeDir 获取用户主目录
func getUserHomeDir() (string, error) {
	return os.UserHomeDir()
}
