package testcmd

import (
	"fmt"
	"os"
)

// TestConfig 配置结构体
type TestConfig struct {
	CheckFile bool
	CheckDir  bool
	CheckLink bool
	Path      string
}

// TestCmdMain 主函数
func TestCmdMain(config TestConfig) error {
	if config.Path == "" {
		return fmt.Errorf("未指定路径")
	}

	count := 0
	if config.CheckFile {
		count++
	}
	if config.CheckDir {
		count++
	}
	if config.CheckLink {
		count++
	}

	if count > 1 {
		return fmt.Errorf("只能指定一个类型标志")
	}

	if config.CheckFile {
		return testFile(config.Path)
	}

	if config.CheckDir {
		return testDir(config.Path)
	}

	if config.CheckLink {
		return testLink(config.Path)
	}

	return testPath(config.Path)
}

// testPath 检测路径是否存在
func testPath(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("路径不存在: %s", path)
		}
		return fmt.Errorf("检测路径失败: %w", err)
	}
	return nil
}

// testFile 检测文件是否存在
func testFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("文件不存在: %s", path)
		}
		return fmt.Errorf("检测文件失败: %w", err)
	}

	if info.IsDir() {
		return fmt.Errorf("不是文件: %s", path)
	}

	return nil
}

// testDir 检测目录是否存在
func testDir(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("目录不存在: %s", path)
		}
		return fmt.Errorf("检测目录失败: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("不是目录: %s", path)
	}

	return nil
}

// testLink 检测符号链接是否存在
func testLink(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("符号链接不存在: %s", path)
		}
		return fmt.Errorf("检测符号链接失败: %w", err)
	}

	if info.Mode()&os.ModeSymlink == 0 {
		return fmt.Errorf("不是符号链接: %s", path)
	}

	return nil
}
