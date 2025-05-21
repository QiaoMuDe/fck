package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
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

	// 检查如果指定了文件大小, 格式是否正确(格式为 +5M 或 -5M), 单位必须为 B/K/M/G 同时为大写
	if *findCmdSize != "" {
		// 使用正则表达式匹配文件大小条件
		sizeRegex := regexp.MustCompile(`^([+-])(\d+)([BKMGbkmg])$`) // 正确分组：符号、数字、单位
		match := sizeRegex.FindStringSubmatch(*findCmdSize)          // 查找匹配项
		if match == nil {
			return fmt.Errorf("文件大小格式错误, 格式如+5M(大于5M)或-5M(小于5M), 支持单位B/K/M/G(大写)")
		}
		_, err := strconv.Atoi(match[2]) // 转换数字部分(match[2])
		if err != nil {
			return fmt.Errorf("文件大小格式错误")
		}
	}

	// 检查是否同时指定了文件和目录和软链接
	if *findCmdFile && *findCmdDir && *findCmdSymlink {
		return fmt.Errorf("不能同时指定 -f、-d 和 -l 标志")
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
			if *findCmdSymlink {
				// 如果只查找软链接，检查文件类型
				fileInfo, err := entry.Info() // 获取文件信息
				if err != nil {
					return nil
				}
				if fileInfo.Mode()&os.ModeSymlink == 0 { // 检查文件是否为软链接
					return nil
				}
			}
			if *findCmdSize != "" {
				// 检查文件大小是否符合要求
				fileInfo, err := entry.Info()
				if err != nil {
					return nil
				}
				if !matchFileSize(fileInfo.Size(), *findCmdSize) {
					return nil
				}
			}
			if *findCmdModTime != "" {
				// 检查修改时间是否符合要求
				fileInfo, err := entry.Info()
				if err != nil {
					return nil
				}
				// 检查文件时间是否符合要求
				if !matchFileTime(fileInfo.ModTime(), *findCmdModTime) {
					return nil
				}
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

// matchFileSize 检查文件大小是否符合指定的条件
func matchFileSize(fileSize int64, sizeCondition string) bool {
	if len(sizeCondition) < 2 {
		return false
	}

	// 获取比较符号和数值部分
	comparator := sizeCondition[0]
	sizeStr := sizeCondition[1:]

	// 获取单位
	unit := sizeStr[len(sizeStr)-1]
	sizeValueStr := sizeStr[:len(sizeStr)-1]

	// 转换数值部分
	sizeValue, err := strconv.ParseFloat(sizeValueStr, 64)
	if err != nil {
		return false
	}

	// 根据单位转换为字节
	var sizeInBytes float64
	switch unit {
	case 'B':
		sizeInBytes = sizeValue
	case 'b':
		sizeInBytes = sizeValue
	case 'K':
		sizeInBytes = sizeValue * 1024
	case 'k':
		sizeInBytes = sizeValue * 1024
	case 'M':
		sizeInBytes = sizeValue * 1024 * 1024
	case 'm':
		sizeInBytes = sizeValue * 1024 * 1024
	case 'G':
		sizeInBytes = sizeValue * 1024 * 1024 * 1024
	case 'g':
		sizeInBytes = sizeValue * 1024 * 1024 * 1024
	default:
		return false
	}

	// 根据比较符号进行比较
	switch comparator {
	case '+':
		return float64(fileSize) > sizeInBytes
	case '-':
		return float64(fileSize) < sizeInBytes
	default:
		return false
	}
}

// matchFileTime 检查文件时间是否符合指定的条件
func matchFileTime(fileTime time.Time, timeCondition string) bool {
	// 检查时间条件是否为空
	if len(timeCondition) < 2 {
		return false
	}

	// 获取比较符号和数值部分
	comparator := timeCondition[0]
	daysStr := timeCondition[1:]

	// 转换天数
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		return false
	}

	// 计算时间阈值
	threshold := time.Now().AddDate(0, 0, -days)

	// 根据比较符号进行比较
	switch comparator {
	case '+':
		return fileTime.After(threshold) // 检查文件时间是否在阈值之后
	case '-':
		return fileTime.Before(threshold) // 检查文件时间是否在阈值之前
	default:
		return false
	}
}
