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
		return fmt.Errorf("path not specified")
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
		return fmt.Errorf("only one type flag can be specified")
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
			return fmt.Errorf("path does not exist: %s", path)
		}
		return fmt.Errorf("error checking path: %w", err)
	}
	return nil
}

// testFile 检测文件是否存在
func testFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", path)
		}
		return fmt.Errorf("error checking file: %w", err)
	}

	if info.IsDir() {
		return fmt.Errorf("not a file: %s", path)
	}

	return nil
}

// testDir 检测目录是否存在
func testDir(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %s", path)
		}
		return fmt.Errorf("error checking directory: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("not a directory: %s", path)
	}

	return nil
}

// testLink 检测符号链接是否存在
func testLink(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("symbolic link does not exist: %s", path)
		}
		return fmt.Errorf("error checking symbolic link: %w", err)
	}

	if info.Mode()&os.ModeSymlink == 0 {
		return fmt.Errorf("not a symbolic link: %s", path)
	}

	return nil
}
