package check

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// hashLineValidator 校验文件行验证器
type hashLineValidator struct {
	hashRegex *regexp.Regexp
}

// newHashLineValidator 创建新的行验证器
func newHashLineValidator() *hashLineValidator {
	// 匹配常见哈希格式：32位(MD5)、40位(SHA1)、64位(SHA256)、128位(SHA512)
	hashRegex := regexp.MustCompile(`^[a-fA-F0-9]{32}$|^[a-fA-F0-9]{40}$|^[a-fA-F0-9]{64}$|^[a-fA-F0-9]{128}$`)
	return &hashLineValidator{
		hashRegex: hashRegex,
	}
}

// validateLine 验证校验文件中的单行内容
func (v *hashLineValidator) validateLine(line string, lineNum int) (hash, filePath string, err error) {
	// 跳过空行和注释行
	if line == "" || strings.HasPrefix(line, "#") {
		return "", "", nil
	}

	parts := strings.Fields(line)
	if len(parts) < 2 {
		return "", "", fmt.Errorf("第%d行格式错误: 缺少哈希值或文件路径", lineNum)
	}

	hash = parts[0]
	filePath = strings.Join(parts[1:], " ")

	// 验证哈希值格式
	if !v.hashRegex.MatchString(hash) {
		return "", "", fmt.Errorf("第%d行哈希值格式无效: %s", lineNum, hash)
	}

	// 清理文件路径
	filePath = v.cleanFilePath(filePath)

	// 验证文件路径安全性
	if err := v.validateFilePath(filePath, lineNum); err != nil {
		return "", "", err
	}

	return hash, filePath, nil
}

// cleanFilePath 清理文件路径
func (v *hashLineValidator) cleanFilePath(filePath string) string {
	// 去除引号
	filePath = strings.Trim(filePath, `"`)
	// 替换双反斜杠
	filePath = strings.ReplaceAll(filePath, `\\`, `\`)
	// 清理路径
	return filepath.Clean(filePath)
}

// validateFilePath 验证文件路径安全性
func (v *hashLineValidator) validateFilePath(filePath string, lineNum int) error {
	// 检查路径遍历攻击
	if strings.Contains(filePath, "..") {
		return fmt.Errorf("第%d行路径包含非法字符 '..': %s", lineNum, filePath)
	}

	// 检查空路径
	if filePath == "" {
		return fmt.Errorf("第%d行文件路径为空", lineNum)
	}

	// 限制路径长度（防止过长路径攻击）
	if len(filePath) > 4096 {
		return fmt.Errorf("第%d行文件路径过长: %s", lineNum, filePath)
	}

	return nil
}
