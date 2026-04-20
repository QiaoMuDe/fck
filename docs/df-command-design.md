# df 命令实现方案

## 1. 设计目标

实现 `fck df` 命令，用于查看磁盘空间使用情况，功能对标 Linux `df -h` 命令。

## 2. 使用示例

```bash
# 查看所有分区（人类可读格式）
fck df

# 只显示本地文件系统
fck df -l

# 按文件系统类型过滤
fck df -t ext4
fck df -t ntfs

# 显示总计
fck df --total

# 简洁模式
fck df -s

# 指定表格样式
fck df --table-style cb
```

## 3. 输出格式

### 标准模式

```
Filesystem      Size   Used   Avail  Use%  Mounted on
/dev/sda1       20GB   15GB   4.5GB  77%   /
/dev/sdb1       100GB  30GB   70GB   30%   /data
C:              100GB  45GB   55GB   45%   C:\
D:              500GB  200GB  300GB  40%   D:\
total           620GB  245GB  375GB  40%
```

### 简洁模式 (-s)

```
/       20GB   77%
/data   100GB  30%
C:\      100GB  45%
```

## 4. 数据结构设计

### 4.1 配置结构体

```go
// DFConfig df 命令配置
type DFConfig struct {
    LocalOnly   bool   // 只显示本地文件系统
    FSFilter    string // 按文件系统类型过滤
    ShowTotal   bool   // 显示总计行
    ListMode    bool   // 简洁模式
    TableStyle  string // 表格样式
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
```

## 5. CLI 参数设计

### 5.1 参数定义

```go
var (
    dfLocal      *qflag.BoolFlag   // -l, --local      只显示本地文件系统
    dfType       *qflag.StringFlag // -t, --type       按文件系统类型过滤
    dfTotal      *qflag.BoolFlag   // --total          显示总计
    dfList       *qflag.BoolFlag   // -s, --simple     简洁模式
    dfTableStyle *qflag.EnumFlag   // --table-style    表格样式
)
```

### 5.2 命令配置

```go
cmdOpts := &qflag.CmdOpts{
    Desc:       "查看磁盘空间使用情况",
    UseChinese: true,
    Examples: map[string]string{
        "查看所有分区":       "fck df",
        "只显示本地分区":     "fck df -l",
        "按类型过滤":        "fck df -t ext4",
        "显示总计":          "fck df --total",
        "简洁模式":          "fck df -s",
    },
    Notes: []string{
        "默认显示所有分区（包括网络文件系统）",
        "大小自动转换为人类可读格式（GB/MB/KB）",
        "Windows 上显示为 C:, D: 等盘符",
    },
}
```

## 6. 核心实现逻辑

### 6.1 主函数流程

```go
func DFCmdMain(config DFConfig) error {
    // 1. 获取所有分区信息
    partitions, err := disk.Partitions(true)
    if err != nil {
        return fmt.Errorf("failed to get partitions: %w", err)
    }

    // 2. 收集各分区使用情况
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
            MountedOn:  p.Mountpoint,
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
    
    // 3. 渲染输出
    if config.ListMode {
        return renderListMode(dfInfos, config)
    }
    return renderTableMode(dfInfos, stats, config)
}
```

### 6.2 本地文件系统判断

```go
// 常见的网络文件系统类型
var networkFSTypes = []string{
    "nfs", "nfs4", "smbfs", "cifs", "afs", 
    "fuse.sshfs", "fuse", "glusterfs", "ceph",
}

func isLocalFS(fsType string) bool {
    for _, nfs := range networkFSTypes {
        if strings.EqualFold(fsType, nfs) {
            return false
        }
    }
    return true
}
```

### 6.3 格式化大小

```go
func formatSize(bytes uint64) string {
    const (
        KB = 1024
        MB = 1024 * KB
        GB = 1024 * MB
        TB = 1024 * GB
    )
    
    switch {
    case bytes >= TB:
        return fmt.Sprintf("%.1fTB", float64(bytes)/TB)
    case bytes >= GB:
        return fmt.Sprintf("%.1fGB", float64(bytes)/GB)
    case bytes >= MB:
        return fmt.Sprintf("%.1fMB", float64(bytes)/MB)
    case bytes >= KB:
        return fmt.Sprintf("%.1fKB", float64(bytes)/KB)
    default:
        return fmt.Sprintf("%dB", bytes)
    }
}
```

### 6.4 表格渲染

```go
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
        t.AppendRow(table.Row{
            "total",
            formatSize(stats.TotalSize),
            formatSize(stats.TotalUsed),
            formatSize(stats.TotalAvail),
            fmt.Sprintf("%.0f%%", float64(stats.TotalUsed)/float64(stats.TotalSize)*100),
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
```

## 7. 跨平台处理

### 7.1 Windows 特殊处理

```go
func normalizeMountpoint(mount string) string {
    // Windows 盘符统一显示为 C:, D: 格式
    if runtime.GOOS == "windows" {
        mount = strings.TrimSuffix(mount, "\\")
        mount = strings.TrimSuffix(mount, "/")
    }
    return mount
}
```

### 7.2 不同平台的分区过滤

| 平台 | 特殊处理 |
|------|----------|
| Windows | 显示盘符如 C:, D: |
| Linux | 过滤 /proc, /sys, /dev 等虚拟分区 |
| macOS | 过滤 /dev, /System/Volumes 等 |

## 8. 文件结构

```
internal/commands/df/
└── cmd_df.go      # 业务逻辑

internal/cli/
└── df.go          # CLI 定义
```

## 9. 实现步骤

1. 创建 `internal/commands/df/cmd_df.go`
   - 定义配置结构体
   - 实现主函数
   - 实现辅助函数（格式化、过滤、渲染）

2. 创建 `internal/cli/df.go`
   - 定义 CLI 参数
   - 配置帮助信息
   - 实现运行函数

3. 修改 `internal/cli/root.go`
   - 注册 df 命令

4. 编译验证
   - `go build ./...`
   - 测试各平台功能

## 10. 与 Linux df 对比

| 功能 | Linux df | fck df |
|------|----------|--------|
| 基本显示 | `df -h` | `fck df` |
| 本地分区 | `df -l` | `fck df -l` |
| 按类型过滤 | `df -t ext4` | `fck df -t ext4` |
| 显示总计 | `df --total` | `fck df --total` |
| 简洁输出 | 无 | `fck df -s` |
| 表格样式 | 无 | `fck df --table-style cb` |

## 11. 依赖

```go
import (
    "github.com/shirou/gopsutil/v3/disk"
    "github.com/jedib0t/go-pretty/v6/table"
)
```

## 12. 注意事项

1. **权限问题**：某些分区可能需要管理员权限才能获取使用情况
2. **网络分区**：网络文件系统（NFS/SMB）可能响应较慢，添加超时处理
3. **虚拟分区**：过滤掉 /proc, /sys, /dev 等虚拟文件系统
4. **Unicode**：挂载点路径可能包含中文，确保正确显示
