package cmd

import (
	"fmt"
	"os"
	"regexp"
	"time"
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

// getLast8Chars 函数用于获取输入字符串的最后8个字符
func getLast8Chars(s string) string {
	if s == "" {
		return ""
	}
	if len(s) <= 8 {
		return s
	}
	return s[len(s)-8:]
}

// writeFileHeader 写入文件头信息
func writeFileHeader(file *os.File, hashType string, timestampFormat string) error {
	// 获取当前时间
	now := time.Now()

	// 构造文件头内容
	header := fmt.Sprintf("#%s#%s\n", hashType, now.Format(timestampFormat))

	// 写入文件头
	if _, err := file.WriteString(header); err != nil {
		return fmt.Errorf("写入文件头失败: %v", err)
	}
	return nil
}

// isSystemFileOrDir 检查文件或目录是否是系统文件或特殊目录
func isSystemFileOrDir(name string) bool {
	// 检查文件或目录是否在列表中
	if systemFilesAndDirs[name] {
		return true
	}

	return false
}

// RegexBuilder 构建正则表达式模式字符串
// pattern: 原始匹配模式
// isRegex: 是否启用正则表达式模式（true时不转义特殊字符）
// wholeWord: 是否启用全字匹配（在模式前后添加^和$）
// caseSensitive: 是否区分大小写（false时添加(?i)标志）
// 返回构建后的正则表达式字符串
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

// CompileRegex 编译正则表达式（仅在模式不为空时）
// pattern: 正则表达式模式字符串
// 返回编译后的正则表达式对象和可能的错误
func CompileRegex(pattern string) (*regexp.Regexp, error) {
	if pattern == "" {
		return nil, nil
	}
	return regexp.Compile(pattern)
}
