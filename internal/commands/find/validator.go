// Package find 实现了文件查找命令参数的验证功能。
// 该文件提供了配置验证器，用于验证所有命令行参数的合法性和安全性。
package find

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gitee.com/MM-Q/fck/internal/types"
)

// ConfigValidator 负责验证find命令的所有参数
type ConfigValidator struct {
	config *FindConfig
}

// NewConfigValidator 创建新的配置验证器
func NewConfigValidator(config *FindConfig) *ConfigValidator {
	return &ConfigValidator{
		config: config,
	}
}

// ValidateArgs 验证find命令的所有参数
func (v *ConfigValidator) ValidateArgs(findPath string) error {
	if err := v.ValidatePath(findPath); err != nil {
		return err
	}

	if err := v.ValidateFlags(); err != nil {
		return err
	}

	return nil
}

// ValidatePath 验证路径的安全性和有效性
func (v *ConfigValidator) ValidatePath(findPath string) error {
	// 检查路径是否存在
	if _, err := os.Lstat(findPath); err != nil {
		// 检查是否是权限不足的错误
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied, cannot access some directories: %s", findPath)
		}

		// 如果是不存在错误, 则返回路径不存在
		if os.IsNotExist(err) {
			return fmt.Errorf("path does not exist: %s", findPath)
		}

		// 其他错误, 返回错误信息
		return fmt.Errorf("error checking path: %s: %v", findPath, err)
	}

	return nil
}

// ValidateFlags 验证所有标志参数的合法性
func (v *ConfigValidator) ValidateFlags() error {
	if v.config.MaxDepth < -1 {
		return fmt.Errorf("max depth cannot be less than -1")
	}

	v.config.Type = strings.ToLower(v.config.Type)

	if !types.IsValidFindType(v.config.Type) {
		return fmt.Errorf("invalid type flag: %s, please use %s", v.config.Type, types.GetSupportedFindTypes()[:])
	}

	// 验证隐藏文件标志
	if !v.config.Hidden && (v.config.Type == types.FindTypeHidden || v.config.Type == types.FindTypeHiddenShort) {
		return fmt.Errorf("must specify -H flag to use -type hidden or -type h")
	}

	// 验证文件大小格式
	if err := v.validateSizeFormat(); err != nil {
		return err
	}

	// 验证修改时间格式
	if err := v.validateTimeFormat(); err != nil {
		return err
	}

	// 验证exec相关参数
	if err := v.validateExecFlags(); err != nil {
		return err
	}

	// 验证mv标志参数
	if v.config.MovePath != "" {
		if info, err := os.Stat(v.config.MovePath); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("path does not exist: %s", v.config.MovePath)
			}
			return fmt.Errorf("error checking path: %v", err)
		} else if !info.IsDir() {
			return fmt.Errorf("-mv flag specified path must be a directory")
		}
	}

	// 验证扩展名参数
	if err := v.validateExtensions(); err != nil {
		return err
	}

	if v.config.MaxDepthLimit < 1 {
		return fmt.Errorf("max depth limit cannot be less than 1")
	}

	return nil
}

// validateSizeFormat 验证文件大小格式
func (v *ConfigValidator) validateSizeFormat() error {
	if v.config.SizePattern != "" {
		sizeRegex := regexp.MustCompile(`^([+-])(\d+)([BKMGbkmg])$`)
		match := sizeRegex.FindStringSubmatch(v.config.SizePattern)
		if match == nil {
			return fmt.Errorf("invalid size pattern, format like +5M (> 5M), (< -5M), support units B/K/M/G (uppercase)")
		}
		_, err := strconv.Atoi(match[2])
		if err != nil {
			return fmt.Errorf("invalid size pattern")
		}
	}
	return nil
}

// validateTimeFormat 验证修改时间格式
func (v *ConfigValidator) validateTimeFormat() error {
	if v.config.ModTimePattern != "" {
		timeRegex := regexp.MustCompile(`^([+-])(\d+)$`)
		match := timeRegex.FindStringSubmatch(v.config.ModTimePattern)
		if match == nil {
			return fmt.Errorf("invalid time pattern, format like +5 (5 days ago) or -5 days (5 days ago)")
		}
		_, err := strconv.Atoi(match[2])
		if err != nil {
			return fmt.Errorf("invalid time pattern")
		}
	}
	return nil
}

// validateExecFlags 验证exec相关标志
func (v *ConfigValidator) validateExecFlags() error {
	if v.config.ExecCmd != "" && !strings.Contains(v.config.ExecCmd, "{}") {
		return fmt.Errorf("invalid exec command, must contain {} as path placeholder when using -exec flag")
	}

	return nil
}

// validateExtensions 验证扩展名参数
func (v *ConfigValidator) validateExtensions() error {
	for _, ext := range v.config.ExtSlice {
		if strings.ContainsAny(ext, " \t\n\r\\/:*?\"<>|") {
			return fmt.Errorf("invalid extension: %s", ext)
		}
	}
	return nil
}
