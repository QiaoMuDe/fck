package df

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"gitee.com/MM-Q/fck/internal/types"
	"gitee.com/MM-Q/go-kit/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/shirou/gopsutil/v3/disk"
)

// DFConfig df 命令配置
type DFConfig struct {
	LocalOnly  bool   // 只显示本地文件系统
	FSFilter   string // 按文件系统类型过滤
	ShowTotal  bool   // 显示总计行
	ListMode   bool   // 简洁模式
	TableStyle string // 表格样式
}

// DFInfo 磁盘分区信息
type DFInfo struct {
	Filesystem string  // 文件系统设备名
	Size       uint64  // 总大小（字节）
	Used       uint64  // 已使用（字节）
	Avail      uint64  // 可用（字节）
	UsePercent float64 // 使用百分比
	MountedOn  string  // 挂载点
	FSType     string  // 文件系统类型
}

// DFStats 统计信息
type DFStats struct {
	TotalSize  uint64
	TotalUsed  uint64
	TotalAvail uint64
	Count      int
}

// 常见的网络文件系统类型
var networkFSTypes = []string{
	"nfs", "nfs4", "smbfs", "cifs", "afs",
	"fuse.sshfs", "fuse", "glusterfs", "ceph",
}

// DFCmdMain df 命令主入口
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - error: 执行错误
func DFCmdMain(config DFConfig) error {
	// 获取所有分区信息
	partitions, err := disk.Partitions(true)
	if err != nil {
		return fmt.Errorf("failed to get partitions: %w", err)
	}

	// 收集各分区使用情况
	var dfInfos []DFInfo
	stats := &DFStats{}

	for _, p := range partitions {
		// 过滤本地文件系统
		if config.LocalOnly && !isLocalFS(p.Fstype) {
			continue
		}

		// 按类型过滤
		if config.FSFilter != "" && !strings.EqualFold(p.Fstype, config.FSFilter) {
			continue
		}

		// 获取使用情况
		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue // 跳过无法访问的分区
		}

		info := DFInfo{
			Filesystem: p.Device,
			Size:       usage.Total,
			Used:       usage.Used,
			Avail:      usage.Free,
			UsePercent: usage.UsedPercent,
			MountedOn:  normalizeMountpoint(p.Mountpoint),
			FSType:     p.Fstype,
		}

		dfInfos = append(dfInfos, info)
		stats.TotalSize += usage.Total
		stats.TotalUsed += usage.Used
		stats.TotalAvail += usage.Free
		stats.Count++
	}

	if len(dfInfos) == 0 {
		fmt.Println("No filesystems found")
		return nil
	}

	// 渲染输出
	if config.ListMode {
		return renderListMode(dfInfos)
	}
	return renderTableMode(dfInfos, stats, config)
}

// isLocalFS 判断是否为本地文件系统
//
// 参数:
//   - fsType: 文件系统类型
//
// 返回值:
//   - bool: 是否为本地文件系统
func isLocalFS(fsType string) bool {
	for _, nfs := range networkFSTypes {
		if strings.EqualFold(fsType, nfs) {
			return false
		}
	}
	return true
}

// normalizeMountpoint 规范化挂载点显示
//
// 参数:
//   - mount: 原始挂载点
//
// 返回值:
//   - string: 规范化后的挂载点
func normalizeMountpoint(mount string) string {
	// Windows 盘符统一显示格式
	if runtime.GOOS == "windows" {
		mount = strings.TrimSuffix(mount, "\\")
		mount = strings.TrimSuffix(mount, "/")
	}
	return mount
}

// formatSize 格式化大小为人类可读格式
//
// 参数:
//   - bytes: 字节数
//
// 返回值:
//   - string: 格式化后的字符串
func formatSize(bytes uint64) string {
	return utils.FormatBytes(int64(bytes))
}

// renderTableMode 表格模式渲染
//
// 参数:
//   - infos: 分区信息列表
//   - stats: 统计信息
//   - config: 命令配置
//
// 返回值:
//   - error: 错误
func renderTableMode(infos []DFInfo, stats *DFStats, config DFConfig) error {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// 表头
	t.AppendHeader(table.Row{"Filesystem", "Size", "Used", "Avail", "Use%", "Mounted on"})

	// 数据行
	for _, info := range infos {
		t.AppendRow(table.Row{
			info.Filesystem,
			formatSize(info.Size),
			formatSize(info.Used),
			formatSize(info.Avail),
			fmt.Sprintf("%.0f%%", info.UsePercent),
			info.MountedOn,
		})
	}

	// 总计行
	if config.ShowTotal && stats.Count > 1 {
		t.AppendSeparator()
		usePercent := float64(0)
		if stats.TotalSize > 0 {
			usePercent = float64(stats.TotalUsed) / float64(stats.TotalSize) * 100
		}
		t.AppendRow(table.Row{
			"total",
			formatSize(stats.TotalSize),
			formatSize(stats.TotalUsed),
			formatSize(stats.TotalAvail),
			fmt.Sprintf("%.0f%%", usePercent),
			"",
		})
	}

	// 设置样式
	if config.TableStyle != "" {
		if style, ok := types.TableStyleMap[config.TableStyle]; ok {
			t.SetStyle(style)
		}
	}

	t.Render()
	return nil
}

// renderListMode 简洁模式渲染
//
// 参数:
//   - infos: 分区信息列表
//
// 返回值:
//   - error: 错误
func renderListMode(infos []DFInfo) error {
	for _, info := range infos {
		fmt.Printf("%-20s %8s %4.0f%%\n", info.MountedOn, formatSize(info.Size), info.UsePercent)
	}
	return nil
}
