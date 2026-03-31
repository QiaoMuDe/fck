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
			return fmt.Errorf("权限不足, 无法访问某些目录: %s", findPath)
		}

		// 如果是不存在错误, 则返回路径不存在
		if os.IsNotExist(err) {
			return fmt.Errorf("路径不存在: %s", findPath)
		}

		// 其他错误, 返回错误信息
		return fmt.Errorf("检查路径时出错: %s: %v", findPath, err)
	}

	return nil
}

// ValidateFlags 验证所有标志参数的合法性
func (v *ConfigValidator) ValidateFlags() error {
	if v.config.MaxDepth < -1 {
		return fmt.Errorf("查找最大深度不能小于 -1")
	}

	v.config.Type = strings.ToLower(v.config.Type)

	if !types.IsValidFindType(v.config.Type) {
		return fmt.Errorf("无效的类型: %s, 请使用%s", v.config.Type, types.GetSupportedFindTypes()[:])
	}

	if !v.config.Hidden && (v.config.Type == types.FindTypeHidden || v.config.Type == types.FindTypeHiddenShort) {
		return fmt.Errorf("必须指定 -H 标志才能使用 -type hidden 或 -type h 选项")
	}

	if err := v.validateSizeFormat(); err != nil {
		return err
	}

	if err := v.validateTimeFormat(); err != nil {
		return err
	}

	if err := v.validateExecFlags(); err != nil {
		return err
	}

	if err := v.validateOperationFlags(); err != nil {
		return err
	}

	if err := v.validateExtensions(); err != nil {
		return err
	}

	if v.config.MaxDepthLimit < 1 {
		return fmt.Errorf("软连接最大解析深度不能小于1")
	}

	return nil
}

// validateSizeFormat 验证文件大小格式
func (v *ConfigValidator) validateSizeFormat() error {
	if v.config.SizePattern != "" {
		sizeRegex := regexp.MustCompile(`^([+-])(\d+)([BKMGbkmg])$`)
		match := sizeRegex.FindStringSubmatch(v.config.SizePattern)
		if match == nil {
			return fmt.Errorf("文件大小格式错误, 格式如+5M(大于5M)或-5M(小于5M), 支持单位B/K/M/G(大写)")
		}
		_, err := strconv.Atoi(match[2])
		if err != nil {
			return fmt.Errorf("文件大小格式错误")
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
			return fmt.Errorf("文件时间格式错误, 格式如+5(5天前)或-5(5天内)")
		}
		_, err := strconv.Atoi(match[2])
		if err != nil {
			return fmt.Errorf("文件时间格式错误")
		}
	}
	return nil
}

// validateExecFlags 验证exec相关标志
func (v *ConfigValidator) validateExecFlags() error {
	if v.config.ExecCmd != "" && !strings.Contains(v.config.ExecCmd, "{}") {
		return fmt.Errorf("使用-exec标志时必须包含{}作为路径占位符")
	}

	return nil
}

// validateOperationFlags 验证操作标志之间的冲突
func (v *ConfigValidator) validateOperationFlags() error {
	if v.config.ExecCmd != "" && (v.config.Delete || v.config.MovePath != "") {
		return fmt.Errorf("使用-exec标志时不能同时指定-delete或-mv标志")
	}

	if v.config.Delete && (v.config.ExecCmd != "" || v.config.MovePath != "") {
		return fmt.Errorf("使用-delete标志时不能同时指定-exec或-mv标志")
	}

	if v.config.MovePath != "" && (v.config.ExecCmd != "" || v.config.Delete) {
		return fmt.Errorf("使用-mv标志时不能同时指定-exec或-delete标志")
	}

	if v.config.MovePath != "" {
		if info, err := os.Stat(v.config.MovePath); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("-mv 标志指定的路径不存在: %s", v.config.MovePath)
			}
			return fmt.Errorf("获取文件信息失败: %v", err)
		} else if !info.IsDir() {
			return fmt.Errorf("-mv标志指定的路径必须为目录")
		}
	}

	if v.config.Count && (v.config.ExecCmd != "" || v.config.MovePath != "" || v.config.Delete) {
		return fmt.Errorf("使用-count标志时不能同时指定-exec、-mv、-delete标志")
	}

	return nil
}

// validateExtensions 验证扩展名参数
func (v *ConfigValidator) validateExtensions() error {
	for _, ext := range v.config.ExtSlice {
		if strings.ContainsAny(ext, " \t\n\r\\/:*?\"<>|") {
			return fmt.Errorf("扩展名包含非法字符: %s", ext)
		}
	}
	return nil
}
