// Package common 提供了跨模块共享的通用工具函数和实用程序。
// 该文件包含文件哈希计算、路径处理、错误处理、颜色输出等常用功能。
package common

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/commands/internal/types"
)

// 定义需要跳过的系统文件和特殊目录
var systemFilesAndDirs = map[string]bool{
	"pagefile.sys":              true,
	"$RECYCLE.BIN":              true,
	"System Volume Information": true,
	"hiberfil.sys":              true,
	"swapfile.sys":              true,
	"DumpStack.log.tmp":         true,
	"Thumbs.db":                 true,
	"Desktop.ini":               true,
	"Autorun.inf":               true,
	"bootmgr":                   true,
	"BOOTNXT":                   true,
	"ntldr":                     true,
	"ntdetect.com":              true,
	"ntbootdd.sys":              true,
}

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
	if systemFilesAndDirs[name] {
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
		return fmt.Errorf("路径 %s 包含无效字符: %v", path, err)
	}

	// 检查是否为权限错误
	if errors.Is(err, os.ErrPermission) {
		return fmt.Errorf("检查路径 %s 时发生了权限错误: %v", path, err)
	}

	// 检查路径是否不存在
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("目录 %s 不存在", path)
	}

	// 其他未知错误的通用处理
	return fmt.Errorf("检查路径 %s 时发生了错误: %v", path, err)
}

// SprintStringColor 根据路径类型以不同颜色输出字符串
//
// 参数:
//   - p: 要检查的路径，用于获取文件类型信息
//   - s: 要着色的字符串内容
//   - cl: colorlib.ColorLib实例，用于彩色输出
//
// 返回:
//   - string: 根据路径类型以不同颜色返回的字符串
func SprintStringColor(p string, s string, cl *colorlib.ColorLib) string {
	// 获取路径信息
	pathInfo, statErr := os.Lstat(p)
	if statErr != nil {
		return cl.Sred(s) // 如果获取路径信息失败, 返回红色输出
	}

	// 根据路径类型设置颜色
	switch mode := pathInfo.Mode(); {
	case mode&os.ModeSymlink != 0:
		// 符号链接 - 使用青色输出
		return cl.Scyan(s)
	case runtime.GOOS == "windows" && mode.IsRegular() && types.WindowsSymlinkExts[filepath.Ext(p)]:
		// Windows下的快捷方式文件 - 使用青色输出
		return cl.Scyan(s)
	case mode.IsDir():
		// 目录 - 使用蓝色输出
		return cl.Sblue(s)
	case mode&os.ModeDevice != 0:
		// 设备文件 - 使用黄色输出
		return cl.Syellow(s)
	case mode&os.ModeNamedPipe != 0:
		// 命名管道 - 使用黄色输出
		return cl.Syellow(s)
	case mode&os.ModeSocket != 0:
		// 套接字文件 - 使用黄色输出
		return cl.Syellow(s)
	case mode&os.ModeType == 0 && mode&os.ModeCharDevice != 0:
		// 字符设备文件 - 使用黄色输出
		return cl.Syellow(s)
	case mode.IsRegular() && pathInfo.Size() == 0:
		// 空文件 - 使用灰色输出
		return cl.Sgray(s)
	case mode.IsRegular() && mode&0111 != 0:
		// 可执行文件 - 使用绿色输出
		return cl.Sgreen(s)
	case runtime.GOOS == "windows" && mode.IsRegular() && types.WindowsExecutableExts[filepath.Ext(p)]:
		// Windows下的可执行文件 - 使用绿色输出
		return cl.Sgreen(s)
	case mode.IsRegular():
		// 普通文件 - 使用白色输出
		return cl.Swhite(s)
	default:
		// 其他类型文件 - 使用白色输出
		return cl.Swhite(s)
	}
}
