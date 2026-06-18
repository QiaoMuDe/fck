// Package utils 提供了跨模块共享的通用工具函数和实用程序。
// 该文件包含文件哈希计算、路径处理、错误处理、颜色输出等常用功能。
package utils

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"gitee.com/MM-Q/fck/internal/types"
)

// GetLast8Chars 获取输入字符串的最后8个字符
//
// 参数:
//   - s: 输入字符串
//
// 返回:
//   - string: 字符串的最后8个字符，如果字符串长度不足8个字符则返回原字符串
func GetLast8Chars(s string) string {
	if s == "" {
		return ""
	}
	if len(s) <= 8 {
		return s
	}
	return s[len(s)-8:]
}

// IsSystemFileOrDir 检查文件或目录是否是系统文件或特殊目录
//
// 参数:
//   - name: 文件或目录名称
//
// 返回:
//   - bool: 是否为系统文件或特殊目录
//
// 注意:
//   - 检查预定义的系统文件和目录列表
//   - 包括Windows系统文件如pagefile.sys、$RECYCLE.BIN等
func IsSystemFileOrDir(name string) bool {
	// 检查文件或目录是否在列表中
	if types.SystemFilesAndDirs[name] {
		return true
	}

	return false
}

// RegexBuilder 构建正则表达式模式字符串
//
// 参数:
//   - pattern: 原始匹配模式
//   - isRegex: 是否启用正则表达式模式（true时不转义特殊字符）
//   - wholeWord: 是否启用全字匹配（在模式前后添加^和$）
//   - caseSensitive: 是否区分大小写（false时添加(?i)标志）
//
// 返回:
//   - string: 构建后的正则表达式字符串
//
// 注意:
//   - 非正则模式下会转义特殊字符
//   - 全字匹配会在模式前后添加^和$
//   - 不区分大小写时会添加(?i)标志
func RegexBuilder(pattern string, isRegex, wholeWord, caseSensitive bool) string {
	if pattern == "" {
		return ""
	}

	// 非正则模式下转义特殊字符
	// if !isRegex {
	// 	pattern = regexp.QuoteMeta(pattern)
	// }

	// 全字匹配处理
	if wholeWord {
		pattern = "^" + pattern + "$"
	}

	// 大小写处理
	if !caseSensitive {
		pattern = "(?i)" + pattern
	}

	return pattern
}

// CompileRegex 编译正则表达式
//
// 参数:
//   - pattern: 正则表达式模式字符串
//
// 返回:
//   - *regexp.Regexp: 编译后的正则表达式对象，如果模式为空则返回nil
//   - error: 编译错误信息，如果编译失败
//
// 注意:
//   - 仅在模式不为空时进行编译
//   - 空模式会返回nil而不是错误
func CompileRegex(pattern string) (*regexp.Regexp, error) {
	if pattern == "" {
		return nil, nil
	}
	return regexp.Compile(pattern)
}

// HandleError 处理路径检查时出现的错误
//
// 参数:
//   - path: 当前正在检查的路径（文件路径或目录路径）
//   - err: 在检查路径过程中产生的错误对象
//
// 返回:
//   - error: 包含更具描述性错误信息的新错误对象
//
// 注意:
//   - 根据不同的错误类型生成对应的错误提示信息
//   - 支持的错误类型：无效字符、权限错误、路径不存在等
//   - 便于调用者定位和处理问题
func HandleError(path string, err error) error {
	// 检查路径是否包含无效字符
	if errors.Is(err, os.ErrInvalid) {
		return fmt.Errorf("path %s contains invalid characters: %v", path, err)
	}

	// 检查是否为权限错误
	if errors.Is(err, os.ErrPermission) {
		return fmt.Errorf("permission denied when checking path %s: %v", path, err)
	}

	// 检查路径是否不存在
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("directory %s does not exist", path)
	}

	// 其他未知错误的通用处理
	return fmt.Errorf("error when checking path %s: %v", path, err)
}
