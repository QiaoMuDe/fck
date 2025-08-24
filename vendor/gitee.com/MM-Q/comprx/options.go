// Package comprx 提供压缩和解压缩操作的配置选项。
//
// 该文件定义了 Options 结构体和相关的配置方法，用于控制压缩和解压缩操作的行为。
// 支持压缩等级设置、进度条显示、文件过滤、路径验证等功能的配置。
//
// 主要类型：
//   - Options: 压缩/解压配置选项结构体
//
// 主要功能：
//   - 提供默认配置选项
//   - 支持链式配置方法
//   - 提供各种预设配置选项
package comprx

import (
	"gitee.com/MM-Q/comprx/types"
)

// Options 压缩/解压配置选项
type Options struct {
	CompressionLevel      types.CompressionLevel // 压缩等级
	OverwriteExisting     bool                   // 是否覆盖已存在的文件
	ProgressEnabled       bool                   // 是否启用进度显示
	ProgressStyle         types.ProgressStyle    // 进度条样式
	DisablePathValidation bool                   // 是否禁用路径验证
	Filter                types.FilterOptions    // 过滤选项
}

// DefaultOptions 返回默认配置选项
//
// 返回:
//   - Options: 默认配置选项
//
// 默认配置:
//   - CompressionLevel: 默认压缩等级
//   - OverwriteExisting: false (不覆盖已存在文件)
//   - ProgressEnabled: false (不显示进度)
//   - ProgressStyle: 文本样式
//   - DisablePathValidation: false (启用路径验证)
func DefaultOptions() Options {
	return Options{
		CompressionLevel:      types.CompressionLevelDefault,
		OverwriteExisting:     false,
		ProgressEnabled:       false,
		ProgressStyle:         types.ProgressStyleText,
		DisablePathValidation: false,
	}
}

// ProgressOptions 返回带进度显示的配置选项
//
// 参数:
//   - style: 进度条样式
//
// 返回:
//   - Options: 带进度显示的配置选项
func ProgressOptions(style types.ProgressStyle) Options {
	opts := DefaultOptions()
	opts.ProgressEnabled = true
	opts.ProgressStyle = style
	return opts
}

// TextProgressOptions 返回文本样式进度条配置选项
//
// 返回:
//   - Options: 文本样式进度条配置选项
//
// 使用示例:
//
//	err := PackOptions("output.zip", "input_dir", TextProgressOptions())
func TextProgressOptions() Options {
	opts := DefaultOptions()
	opts.ProgressEnabled = true
	opts.ProgressStyle = types.ProgressStyleText
	return opts
}

// UnicodeProgressOptions 返回Unicode样式进度条配置选项
//
// 返回:
//   - Options: Unicode样式进度条配置选项
//
// 使用示例:
//
//	err := PackOptions("output.zip", "input_dir", UnicodeProgressOptions())
func UnicodeProgressOptions() Options {
	opts := DefaultOptions()
	opts.ProgressEnabled = true
	opts.ProgressStyle = types.ProgressStyleUnicode
	return opts
}

// ASCIIProgressOptions 返回ASCII样式进度条配置选项
//
// 返回:
//   - Options: ASCII样式进度条配置选项
//
// 使用示例:
//
//	err := PackOptions("output.zip", "input_dir", ASCIIProgressOptions())
func ASCIIProgressOptions() Options {
	opts := DefaultOptions()
	opts.ProgressEnabled = true
	opts.ProgressStyle = types.ProgressStyleASCII
	return opts
}

// DefaultProgressOptions 返回默认样式进度条配置选项
//
// 返回:
//   - Options: 默认样式进度条配置选项
//
// 使用示例:
//
//	err := PackOptions("output.zip", "input_dir", DefaultProgressOptions())
func DefaultProgressOptions() Options {
	opts := DefaultOptions()
	opts.ProgressEnabled = true
	opts.ProgressStyle = types.ProgressStyleDefault
	return opts
}

// ForceOptions 返回强制模式配置选项
//
// 返回:
//   - Options: 强制模式配置选项
//
// 配置特点:
//   - OverwriteExisting: true (覆盖已存在文件)
//   - DisablePathValidation: true (禁用路径验证)
//   - ProgressEnabled: false (关闭进度条)
//
// 使用示例:
//
//	err := PackOptions("output.zip", "input_dir", ForceOptions())
func ForceOptions() Options {
	opts := DefaultOptions()
	opts.OverwriteExisting = true
	opts.DisablePathValidation = true
	opts.ProgressEnabled = false
	return opts
}

// NoCompressionOptions 返回禁用压缩且启用进度条的配置选项
//
// 返回:
//   - Options: 禁用压缩且启用进度条的配置选项
//
// 配置特点:
//   - CompressionLevel: 无压缩 (存储模式)
//   - ProgressEnabled: true (启用进度条)
//   - ProgressStyle: 文本样式
//
// 使用示例:
//
//	err := PackOptions("output.zip", "input_dir", NoCompressionOptions())
func NoCompressionOptions() Options {
	opts := DefaultOptions()
	opts.CompressionLevel = types.CompressionLevelNone
	opts.ProgressEnabled = true
	opts.ProgressStyle = types.ProgressStyleText
	return opts
}

// NoCompressionProgressOptions 返回禁用压缩且启用指定样式进度条的配置选项
//
// 参数:
//   - style: 进度条样式
//
// 返回:
//   - Options: 禁用压缩且启用指定样式进度条的配置选项
//
// 配置特点:
//   - CompressionLevel: 无压缩 (存储模式)
//   - ProgressEnabled: true (启用进度条)
//   - ProgressStyle: 指定样式
//
// 使用示例:
//
//	err := PackOptions("output.zip", "input_dir", NoCompressionProgressOptions(types.ProgressStyleUnicode))
func NoCompressionProgressOptions(style types.ProgressStyle) Options {
	opts := DefaultOptions()
	opts.CompressionLevel = types.CompressionLevelNone
	opts.ProgressEnabled = true
	opts.ProgressStyle = style
	return opts
}

// ==============================================
// Options Set 方法（直接设置，不返回对象）
// ==============================================

// SetCompressionLevel 设置压缩等级
//
// 参数:
//   - level: 压缩等级
//
// 使用示例:
//
//	opts := DefaultOptions()
//	opts.SetCompressionLevel(types.CompressionLevelBest)
func (o *Options) SetCompressionLevel(level types.CompressionLevel) {
	o.CompressionLevel = level
}

// SetOverwriteExisting 设置是否覆盖已存在的文件
//
// 参数:
//   - overwrite: 是否覆盖已存在文件
//
// 使用示例:
//
//	opts := DefaultOptions()
//	opts.SetOverwriteExisting(true)
func (o *Options) SetOverwriteExisting(overwrite bool) {
	o.OverwriteExisting = overwrite
}

// SetProgress 设置是否启用进度显示
//
// 参数:
//   - enabled: 是否启用进度显示
//
// 使用示例:
//
//	opts := DefaultOptions()
//	opts.SetProgress(true)
func (o *Options) SetProgress(enabled bool) {
	o.ProgressEnabled = enabled
}

// SetProgressStyle 设置进度条样式
//
// 参数:
//   - style: 进度条样式
//
// 使用示例:
//
//	opts := DefaultOptions()
//	opts.SetProgressStyle(types.ProgressStyleUnicode)
func (o *Options) SetProgressStyle(style types.ProgressStyle) {
	o.ProgressStyle = style
}

// SetProgressAndStyle 设置进度显示和样式
//
// 参数:
//   - enabled: 是否启用进度显示
//   - style: 进度条样式
//
// 使用示例:
//
//	opts := DefaultOptions()
//	opts.SetProgressAndStyle(true, types.ProgressStyleUnicode)
func (o *Options) SetProgressAndStyle(enabled bool, style types.ProgressStyle) {
	o.ProgressEnabled = enabled
	o.ProgressStyle = style
}

// SetDisablePathValidation 设置是否禁用路径验证
//
// 参数:
//   - disable: 是否禁用路径验证
//
// 使用示例:
//
//	opts := DefaultOptions()
//	opts.SetDisablePathValidation(true)
func (o *Options) SetDisablePathValidation(disable bool) {
	o.DisablePathValidation = disable
}

// SetFilter 设置过滤配置
//
// 参数:
//   - filter: 过滤选项
//
// 使用示例:
//
//	opts := DefaultOptions()
//	filter := types.FilterOptions{
//	    Include: []string{"*.go", "*.md"},
//	    Exclude: []string{"*_test.go"},
//	}
//	opts.SetFilter(filter)
func (o *Options) SetFilter(filter types.FilterOptions) {
	o.Filter = filter
}

// SetInclude 设置包含模式
//
// 参数:
//   - patterns: 包含模式列表
//
// 使用示例:
//
//	opts := DefaultOptions()
//	opts.SetInclude([]string{"*.go", "*.md"})
func (o *Options) SetInclude(patterns []string) {
	o.Filter.Include = patterns
}

// SetExclude 设置排除模式
//
// 参数:
//   - patterns: 排除模式列表
//
// 使用示例:
//
//	opts := DefaultOptions()
//	opts.SetExclude([]string{"*_test.go", "vendor/*"})
func (o *Options) SetExclude(patterns []string) {
	o.Filter.Exclude = patterns
}

// SetSizeFilter 设置文件大小过滤
//
// 参数:
//   - minSize: 最小文件大小（字节）
//   - maxSize: 最大文件大小（字节）
//
// 使用示例:
//
//	opts := DefaultOptions()
//	opts.SetSizeFilter(1024, 10*1024*1024) // 1KB - 10MB
func (o *Options) SetSizeFilter(minSize, maxSize int64) {
	o.Filter.MinSize = minSize
	o.Filter.MaxSize = maxSize
}

// SetMaxSize 设置最大文件大小
//
// 参数:
//   - maxSize: 最大文件大小（字节）
//
// 使用示例:
//
//	opts := DefaultOptions()
//	opts.SetMaxSize(10 * 1024 * 1024) // 10MB
func (o *Options) SetMaxSize(maxSize int64) {
	o.Filter.MaxSize = maxSize
}

// SetMinSize 设置最小文件大小
//
// 参数:
//   - minSize: 最小文件大小（字节）
//
// 使用示例:
//
//	opts := DefaultOptions()
//	opts.SetMinSize(1024) // 1KB
func (o *Options) SetMinSize(minSize int64) {
	o.Filter.MinSize = minSize
}

// ==============================================
// Options 链式配置方法（通过 Set 方法实现）
// ==============================================

// WithCompressionLevel 设置压缩等级
//
// 参数:
//   - level: 压缩等级
//
// 返回:
//   - Options: 配置选项（支持链式调用）
//
// 使用示例:
//
//	opts := DefaultOptions().WithCompressionLevel(types.CompressionLevelBest)
func (o Options) WithCompressionLevel(level types.CompressionLevel) Options {
	o.SetCompressionLevel(level)
	return o
}

// WithOverwriteExisting 设置是否覆盖已存在的文件
//
// 参数:
//   - overwrite: 是否覆盖已存在文件
//
// 返回:
//   - Options: 配置选项（支持链式调用）
//
// 使用示例:
//
//	opts := DefaultOptions().WithOverwriteExisting(true)
func (o Options) WithOverwriteExisting(overwrite bool) Options {
	o.SetOverwriteExisting(overwrite)
	return o
}

// WithProgress 设置是否启用进度显示
//
// 参数:
//   - enabled: 是否启用进度显示
//
// 返回:
//   - Options: 配置选项（支持链式调用）
//
// 使用示例:
//
//	opts := DefaultOptions().WithProgress(true)
func (o Options) WithProgress(enabled bool) Options {
	o.SetProgress(enabled)
	return o
}

// WithProgressStyle 设置进度条样式
//
// 参数:
//   - style: 进度条样式
//
// 返回:
//   - Options: 配置选项（支持链式调用）
//
// 使用示例:
//
//	opts := DefaultOptions().WithProgressStyle(types.ProgressStyleUnicode)
func (o Options) WithProgressStyle(style types.ProgressStyle) Options {
	o.SetProgressStyle(style)
	return o
}

// WithProgressAndStyle 设置进度显示和样式
//
// 参数:
//   - enabled: 是否启用进度显示
//   - style: 进度条样式
//
// 返回:
//   - Options: 配置选项（支持链式调用）
//
// 使用示例:
//
//	opts := DefaultOptions().WithProgressAndStyle(true, types.ProgressStyleUnicode)
func (o Options) WithProgressAndStyle(enabled bool, style types.ProgressStyle) Options {
	o.SetProgressAndStyle(enabled, style)
	return o
}

// WithDisablePathValidation 设置是否禁用路径验证
//
// 参数:
//   - disable: 是否禁用路径验证
//
// 返回:
//   - Options: 配置选项（支持链式调用）
//
// 使用示例:
//
//	opts := DefaultOptions().WithDisablePathValidation(true)
func (o Options) WithDisablePathValidation(disable bool) Options {
	o.SetDisablePathValidation(disable)
	return o
}

// WithFilter 设置过滤配置
//
// 参数:
//   - filter: 过滤选项
//
// 返回:
//   - Options: 配置选项（支持链式调用）
//
// 使用示例:
//
//	filter := types.FilterOptions{
//	    Include: []string{"*.go", "*.md"},
//	    Exclude: []string{"*_test.go"},
//	}
//	opts := DefaultOptions().WithFilter(filter)
func (o Options) WithFilter(filter types.FilterOptions) Options {
	o.SetFilter(filter)
	return o
}

// WithInclude 设置包含模式
//
// 参数:
//   - patterns: 包含模式列表
//
// 返回:
//   - Options: 配置选项（支持链式调用）
//
// 使用示例:
//
//	opts := DefaultOptions().WithInclude([]string{"*.go", "*.md"})
func (o Options) WithInclude(patterns []string) Options {
	o.SetInclude(patterns)
	return o
}

// WithExclude 设置排除模式
//
// 参数:
//   - patterns: 排除模式列表
//
// 返回:
//   - Options: 配置选项（支持链式调用）
//
// 使用示例:
//
//	opts := DefaultOptions().WithExclude([]string{"*_test.go", "vendor/*"})
func (o Options) WithExclude(patterns []string) Options {
	o.SetExclude(patterns)
	return o
}

// WithSizeFilter 设置文件大小过滤
//
// 参数:
//   - minSize: 最小文件大小（字节）
//   - maxSize: 最大文件大小（字节）
//
// 返回:
//   - Options: 配置选项（支持链式调用）
//
// 使用示例:
//
//	opts := DefaultOptions().WithSizeFilter(1024, 10*1024*1024) // 1KB - 10MB
func (o Options) WithSizeFilter(minSize, maxSize int64) Options {
	o.SetSizeFilter(minSize, maxSize)
	return o
}

// WithMaxSize 设置最大文件大小
//
// 参数:
//   - maxSize: 最大文件大小（字节）
//
// 返回:
//   - Options: 配置选项（支持链式调用）
//
// 使用示例:
//
//	opts := DefaultOptions().WithMaxSize(10 * 1024 * 1024) // 10MB
func (o Options) WithMaxSize(maxSize int64) Options {
	o.SetMaxSize(maxSize)
	return o
}

// WithMinSize 设置最小文件大小
//
// 参数:
//   - minSize: 最小文件大小（字节）
//
// 返回:
//   - Options: 配置选项（支持链式调用）
//
// 使用示例:
//
//	opts := DefaultOptions().WithMinSize(1024) // 1KB
func (o Options) WithMinSize(minSize int64) Options {
	o.SetMinSize(minSize)
	return o
}
