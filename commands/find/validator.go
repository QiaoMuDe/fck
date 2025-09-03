// Package find 实现了文件查找命令参数的验证功能。
// 该文件提供了配置验证器，用于验证所有命令行参数的合法性和安全性。
package find

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gitee.com/MM-Q/fck/commands/internal/types"
)

// ConfigValidator 负责验证find命令的所有参数
type ConfigValidator struct{}

// NewConfigValidator 创建新的配置验证器
func NewConfigValidator() *ConfigValidator {
	return &ConfigValidator{}
}

// ValidateArgs 验证find命令的所有参数
func (v *ConfigValidator) ValidateArgs(findPath string) error {
	// 验证路径安全性
	if err := v.ValidatePath(findPath); err != nil {
		return err
	}

	// 验证标志参数
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
	// 检查要查找的最大深度是否小于 -1
	if findCmdMaxDepth.Get() < -1 {
		return fmt.Errorf("查找最大深度不能小于 -1")
	}

	// 将限制查找类型转为小写
	if setErr := findCmdType.Set(strings.ToLower(findCmdType.Get())); setErr != nil {
		return fmt.Errorf("转换查找类型失败: %v", setErr)
	}

	// 检查是否为受支持限制查找类型
	if !types.IsValidFindType(findCmdType.Get()) {
		return fmt.Errorf("无效的类型: %s, 请使用%s", findCmdType.Get(), types.GetSupportedFindTypes()[:])
	}

	// 如果只显示隐藏文件或目录, 则必须指定 -H 标志
	if !findCmdHidden.Get() && (findCmdType.Get() == types.FindTypeHidden || findCmdType.Get() == types.FindTypeHiddenShort) {
		return fmt.Errorf("必须指定 -H 标志才能使用 -type hidden 或 -type h 选项")
	}

	// 检查是否同时指定了-or
	if findCmdOr.Get() {
		// 如果使用-or, 则不能同时使用-and
		if setErr := findCmdAnd.Set("false"); setErr != nil {
			return fmt.Errorf("设置 -and 失败: %v", setErr)
		}
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

	// 验证操作标志冲突
	if err := v.validateOperationFlags(); err != nil {
		return err
	}

	// 验证扩展名参数
	if err := v.validateExtensions(); err != nil {
		return err
	}

	// 检查软连接最大解析深度是否小于 1
	if findCmdMaxDepthLimit.Get() < 1 {
		return fmt.Errorf("软连接最大解析深度不能小于1")
	}

	return nil
}

// validateSizeFormat 验证文件大小格式
func (v *ConfigValidator) validateSizeFormat() error {
	if findCmdSize.Get() != "" {
		// 使用正则表达式匹配文件大小条件
		sizeRegex := regexp.MustCompile(`^([+-])(\d+)([BKMGbkmg])$`)
		match := sizeRegex.FindStringSubmatch(findCmdSize.Get())
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
	if findCmdModTime.Get() != "" {
		// 使用正则表达式匹配文件时间条件
		timeRegex := regexp.MustCompile(`^([+-])(\d+)$`)
		match := timeRegex.FindStringSubmatch(findCmdModTime.Get())
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
	// 检查-exec标志是否包含{}
	if findCmdExec.Get() != "" && !strings.Contains(findCmdExec.Get(), "{}") {
		return fmt.Errorf("使用-exec标志时必须包含{}作为路径占位符")
	}

	return nil
}

// validateOperationFlags 验证操作标志之间的冲突
func (v *ConfigValidator) validateOperationFlags() error {
	// 检查-exec标志是否与-delete或-mv一起使用
	if findCmdExec.Get() != "" && (findCmdDelete.Get() || findCmdMove.Get() != "") {
		return fmt.Errorf("使用-exec标志时不能同时指定-delete或-mv标志")
	}

	// 检查-delete标志是否与-exec或-mv一起使用
	if findCmdDelete.Get() && (findCmdExec.Get() != "" || findCmdMove.Get() != "") {
		return fmt.Errorf("使用-delete标志时不能同时指定-exec或-mv标志")
	}

	// 检查-mv标志是否与-exec或-delete一起使用
	if findCmdMove.Get() != "" && (findCmdExec.Get() != "" || findCmdDelete.Get()) {
		return fmt.Errorf("使用-mv标志时不能同时指定-exec或-delete标志")
	}

	// 检查-mv标志指定的路径是否为目录
	if findCmdMove.Get() != "" {
		if info, err := os.Stat(findCmdMove.Get()); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("-mv 标志指定的路径不存在: %s", findCmdMove.Get())
			}
			return fmt.Errorf("获取文件信息失败: %v", err)
		} else if !info.IsDir() {
			return fmt.Errorf("-mv标志指定的路径必须为目录")
		}
	}

	// 检查如果指定了-count则不能同时指定 -exec、-mv、-delete
	if findCmdCount.Get() && (findCmdExec.Get() != "" || findCmdMove.Get() != "" || findCmdDelete.Get()) {
		return fmt.Errorf("使用-count标志时不能同时指定-exec、-mv、-delete标志")
	}

	return nil
}

// validateExtensions 验证扩展名参数
func (v *ConfigValidator) validateExtensions() error {
	if findCmdExt.Len() > 0 {
		for _, ext := range findCmdExt.Get() {
			// 检查扩展名是否包含常见危险字符
			if strings.ContainsAny(ext, " \t\n\r\\/:*?\"<>|") {
				return fmt.Errorf("扩展名包含非法字符: %s", ext)
			}
		}
	}
	return nil
}
