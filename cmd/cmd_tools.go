package cmd

import (
	"fmt"
	"os"
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
