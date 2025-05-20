package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// findCmdMain 是 find 子命令的主函数
func findCmdMain() error {
	// 检查要查找的路径是否为空
	if *findCmdPath == "" {
		return fmt.Errorf("查找路径不能为空")
	}

	// 检查要查找的路径是否存在
	if _, err := os.Stat(*findCmdPath); err != nil {
		return fmt.Errorf("查找路径不存在: %s", *findCmdPath)
	}

	// 检查要查找的最大深度是否小于 -1
	if *findCmdMaxDepth < -1 {
		return fmt.Errorf("查找最大深度不能小于 -1")
	}

	// 检查是否同时指定了文件和目录
	if *findCmdFile && *findCmdDir {
		return fmt.Errorf("不能同时指定 -f 和 -d 标志")
	}

	// 编译关键字为正则表达式
	keywordRegex, err := regexp.Compile(*findCmdKeyword)
	if err != nil {
		return fmt.Errorf("关键字格式错误: %s", err)
	}

	// 使用 filepath.WalkDir 遍历目录
	walkDirErr := filepath.WalkDir(*findCmdPath, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("访问文件时出错：%s", err)
		}

		// 检查当前路径的深度是否超过最大深度
		depth := strings.Count(filepath.ToSlash(path[len(*findCmdPath):]), "/")
		if *findCmdMaxDepth >= 0 && depth > *findCmdMaxDepth {
			return filepath.SkipDir
		}

		// 检查文件名是否匹配关键字
		if keywordRegex.MatchString(entry.Name()) {
			// 根据用户选择过滤文件或目录
			if *findCmdFile && entry.IsDir() {
				// 如果只查找文件，跳过目录
				return nil
			}
			if *findCmdDir && !entry.IsDir() {
				// 如果只查找目录，跳过文件
				return nil
			}
			// 输出匹配的路径
			fmt.Println(path)
		}

		return nil
	})

	if walkDirErr != nil {
		return fmt.Errorf("遍历目录时出错: %s", walkDirErr)
	}

	return nil
}
