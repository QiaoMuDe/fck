// Package list 实现了文件列表的数据处理功能。
// 该文件提供了文件列表的排序处理，支持按名称、时间、大小等多种方式排序。
package list

import (
	"sort"
	"strings"
)

// FileProcessor 文件数据处理器
type FileProcessor struct{}

// NewFileProcessor 创建新的文件处理器
func NewFileProcessor() *FileProcessor {
	return &FileProcessor{}
}

// Process 处理文件列表: 过滤 -> 排序 -> 分组
//
// 参数:
//   - files: 文件列表
//   - opts: 处理选项
//
// 返回:
//   - FileInfoList: 处理后的文件列表
func (p *FileProcessor) Process(files FileInfoList, opts ProcessOptions) FileInfoList {
	// 当前版本主要处理排序，过滤已在扫描阶段完成
	return p.sort(files, opts)
}

// sort 对文件列表进行排序
//
// 参数:
//   - files: 文件列表
//   - opts: 排序选项
//
// 返回:
//   - FileInfoList: 排序后的文件列表
func (p *FileProcessor) sort(files FileInfoList, opts ProcessOptions) FileInfoList {
	if len(files) <= 1 {
		return files
	}

	// 创建副本避免修改原始数据
	sorted := make(FileInfoList, len(files))
	copy(sorted, files)

	// 根据排序类型进行排序
	switch opts.SortBy {
	case "time": // 按修改时间排序
		p.sortByTime(sorted, opts.Reverse)
	case "size": // 按文件大小排序
		p.sortBySize(sorted, opts.Reverse)
	case "name": // 按名称排序
		p.sortByName(sorted, opts.Reverse)
	default:
		// 默认按名称排序
		p.sortByName(sorted, opts.Reverse)
	}

	return sorted
}

// sortByTime 按修改时间排序
//
// 参数:
//   - files: 文件列表
//   - reverse: 是否逆序
func (p *FileProcessor) sortByTime(files FileInfoList, reverse bool) {
	sort.Slice(files, func(i, j int) bool {
		result := files[i].ModTime.After(files[j].ModTime)
		if reverse {
			return !result
		}
		return result
	})
}

// sortBySize 按文件大小排序
//
// 参数:
//   - files: 文件列表
//   - reverse: 是否逆序
func (p *FileProcessor) sortBySize(files FileInfoList, reverse bool) {
	sort.Slice(files, func(i, j int) bool {
		result := files[i].Size > files[j].Size
		if reverse {
			return !result
		}
		return result
	})
}

// sortByName 按文件名排序
//
// 参数:
//   - files: 文件列表
//   - reverse: 是否逆序
func (p *FileProcessor) sortByName(files FileInfoList, reverse bool) {
	sort.Slice(files, func(i, j int) bool {
		name1 := strings.ToLower(files[i].Name) // 转换为小写
		name2 := strings.ToLower(files[j].Name) // 转换为小写
		result := name1 < name2                 // 判断是否小于
		if reverse {
			return !result
		}
		return result
	})
}
