// Package proc 实现了进程查看功能。
// 该文件提供了查看系统进程信息的功能，支持过滤、排序和多种输出格式。
package proc

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"gitee.com/MM-Q/fck/internal/types"
	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/shirou/gopsutil/v3/process"
)

// ProcInfo 进程信息结构体
type ProcInfo struct {
	// 基础信息
	PID     int32  // 进程 ID
	PPID    int32  // 父进程 ID
	Name    string // 进程名
	CmdLine string // 命令行

	// 用户信息
	Username string // 用户名

	// 资源使用
	CPUPercent float64 // CPU 使用率
	MemPercent float32 // 内存使用率
	MemRSS     uint64  // 实际使用内存(RSS) MB

	// 状态信息
	Status     string // 进程状态
	NumThreads int32  // 线程数

	// 时间信息
	CreateTime int64  // 创建时间(时间戳)
	RunTime    string // 运行时长
}

// ProcConfig 进程命令配置结构体
type ProcConfig struct {
	// 过滤条件（核心）
	Name string // 进程名
	PID  int32  // 单个 PID
	PIDs []int  // 多个 PID

	// 排序
	SortBy string
	Ascend bool

	// 显示模式
	ListMode   bool
	TreeMode   bool
	JSONMode   bool
	TableStyle string // 表格样式
}

// ProcStats 统计信息结构体
type ProcStats struct {
	TotalCount    int
	RunningCount  int
	SleepingCount int
	StoppedCount  int
	ZombieCount   int
}

// ProcCmdMain 进程命令主入口
//
// 参数:
//   - config: 进程命令配置
//
// 返回:
//   - error: 执行错误
func ProcCmdMain(config *ProcConfig) error {
	// 收集进程信息
	procs, stats, err := collectProcs(config)
	if err != nil {
		return fmt.Errorf("failed to collect process information: %w", err)
	}

	if len(procs) == 0 {
		fmt.Println("No matching processes found")
		return nil
	}

	// 根据模式渲染输出
	switch {
	case config.JSONMode:
		return renderJSONMode(procs, stats)
	case config.TreeMode:
		return renderTreeMode(procs, config)
	case config.ListMode:
		return renderListMode(procs)
	default:
		return renderTableMode(procs, stats, config)
	}
}

// collectProcs 收集进程信息
//
// 参数:
//   - config: 进程命令配置
//
// 返回:
//   - []ProcInfo: 进程信息列表
//   - *ProcStats: 统计信息
//   - error: 错误
func collectProcs(config *ProcConfig) ([]ProcInfo, *ProcStats, error) {
	var allProcs []ProcInfo
	stats := &ProcStats{}

	// 获取所有进程
	pids, err := process.Pids()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get process list: %w", err)
	}

	// 如果指定了单个 PID，转换为 PIDs 列表
	if config.PID > 0 {
		config.PIDs = []int{int(config.PID)}
	}

	// 遍历所有进程
	for _, pid := range pids {
		// 如果指定了 PID 列表，检查是否匹配
		if len(config.PIDs) > 0 {
			found := false
			for _, targetPID := range config.PIDs {
				if int(pid) == targetPID {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// 获取进程信息
		procInfo, err := getProcInfo(pid)
		if err != nil {
			// 某些进程可能无法访问，跳过
			continue
		}

		// 按进程名过滤
		if config.Name != "" && !strings.Contains(strings.ToLower(procInfo.Name), strings.ToLower(config.Name)) {
			continue
		}

		allProcs = append(allProcs, procInfo)

		// 更新统计
		stats.TotalCount++
		switch procInfo.Status {
		case "R":
			stats.RunningCount++
		case "S":
			stats.SleepingCount++
		case "T":
			stats.StoppedCount++
		case "Z":
			stats.ZombieCount++
		}
	}

	// 排序
	sortProcs(allProcs, config.SortBy, config.Ascend)

	return allProcs, stats, nil
}

// getProcInfo 获取单个进程信息
//
// 参数:
//   - pid: 进程 ID
//
// 返回:
//   - ProcInfo: 进程信息
//   - error: 错误
func getProcInfo(pid int32) (ProcInfo, error) {
	var info ProcInfo
	info.PID = pid

	proc, err := process.NewProcess(pid)
	if err != nil {
		return info, err
	}

	// 进程名
	name, err := proc.Name()
	if err == nil {
		info.Name = name
	}

	// 父进程 ID
	ppid, err := proc.Ppid()
	if err == nil {
		info.PPID = ppid
	}

	// 命令行
	cmdline, err := proc.Cmdline()
	if err == nil {
		info.CmdLine = cmdline
	}

	// 用户名
	username, err := proc.Username()
	if err == nil {
		info.Username = username
	}

	// CPU 使用率
	cpuPercent, err := proc.CPUPercent()
	if err == nil {
		info.CPUPercent = cpuPercent
	}

	// 内存信息
	memInfo, err := proc.MemoryInfo()
	if err == nil {
		info.MemRSS = memInfo.RSS / 1024 / 1024 // 转换为 MB
	}

	// 内存使用率
	memPercent, err := proc.MemoryPercent()
	if err == nil {
		info.MemPercent = memPercent
	}

	// 状态
	status, err := proc.Status()
	if err == nil && len(status) > 0 {
		info.Status = status[0]
	}

	// 线程数
	numThreads, err := proc.NumThreads()
	if err == nil {
		info.NumThreads = numThreads
	}

	// 创建时间
	createTime, err := proc.CreateTime()
	if err == nil {
		info.CreateTime = createTime
		info.RunTime = formatDuration(time.Since(time.Unix(createTime/1000, 0)))
	}

	return info, nil
}

// formatDuration 格式化时长
//
// 参数:
//   - d: 时长
//
// 返回:
//   - string: 格式化后的字符串
func formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd%dh", days, hours)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh%dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

// sortProcs 排序进程列表
//
// 参数:
//   - procs: 进程列表
//   - sortBy: 排序字段
//   - ascend: 是否升序
func sortProcs(procs []ProcInfo, sortBy string, ascend bool) {
	less := func(i, j int) bool {
		var result bool
		switch sortBy {
		case "name":
			result = procs[i].Name < procs[j].Name
		case "cpu":
			result = procs[i].CPUPercent < procs[j].CPUPercent
		case "mem":
			result = procs[i].MemRSS < procs[j].MemRSS
		case "time":
			result = procs[i].CreateTime < procs[j].CreateTime
		default: // pid
			result = procs[i].PID < procs[j].PID
		}
		if ascend {
			return result
		}
		return !result
	}
	sort.Slice(procs, less)
}

// renderTableMode 表格模式渲染
//
// 参数:
//   - procs: 进程信息列表
//   - stats: 统计信息
//   - config: 进程命令配置
//
// 返回:
//   - error: 错误
func renderTableMode(procs []ProcInfo, stats *ProcStats, config *ProcConfig) error {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// 设置表头
	t.AppendHeader(table.Row{"PID", "Name", "User", "Status", "CPU(%)", "Mem(MB)", "Threads", "Runtime"})

	// 添加数据行
	for _, proc := range procs {
		status := formatStatus(proc.Status)
		name := proc.Name
		if name == "" {
			name = "-"
		}
		t.AppendRow(table.Row{
			proc.PID,
			truncateString(name, 30),
			truncateString(proc.Username, 20),
			status,
			fmt.Sprintf("%.1f", proc.CPUPercent),
			proc.MemRSS,
			proc.NumThreads,
			proc.RunTime,
		})
	}

	// 设置表格样式
	if config.TableStyle != "" {
		if style, ok := types.TableStyleMap[config.TableStyle]; ok {
			t.SetStyle(style)
		}
	}

	// 添加统计信息
	t.AppendSeparator()
	t.AppendFooter(table.Row{
		fmt.Sprintf("Total: %d", stats.TotalCount),
		fmt.Sprintf("R: %d", stats.RunningCount),
		fmt.Sprintf("S: %d", stats.SleepingCount),
		fmt.Sprintf("T: %d", stats.StoppedCount),
		fmt.Sprintf("Z: %d", stats.ZombieCount),
		"",
		"",
		"",
	})

	t.Render()
	return nil
}

// formatStatus 格式化状态
//
// 参数:
//   - status: 原始状态
//
// 返回:
//   - string: 格式化后的状态
func formatStatus(status string) string {
	switch status {
	case "R":
		return "Running"
	case "S":
		return "Sleeping"
	case "T":
		return "Stopped"
	case "Z":
		return "Zombie"
	default:
		return status
	}
}

// truncateString 截断字符串
//
// 参数:
//   - s: 原始字符串
//   - maxLen: 最大长度
//
// 返回:
//   - string: 截断后的字符串
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// renderListMode 简洁模式渲染
//
// 参数:
//   - procs: 进程信息列表
//
// 返回:
//   - error: 错误
func renderListMode(procs []ProcInfo) error {
	for _, proc := range procs {
		status := formatStatus(proc.Status)
		name := proc.Name
		if name == "" {
			name = "-"
		}
		output := fmt.Sprintf("[%d] %s (%s, CPU: %.1f%%, MEM: %dMB, %s)",
			proc.PID, name, proc.Username, proc.CPUPercent, proc.MemRSS, status)
		fmt.Println(output)
	}
	return nil
}

// renderTreeMode 树形模式渲染
//
// 参数:
//   - procs: 进程信息列表
//   - config: 进程命令配置
//
// 返回:
//   - error: 错误
func renderTreeMode(procs []ProcInfo, config *ProcConfig) error {
	// 构建进程树
	procMap := make(map[int32]ProcInfo)
	childrenMap := make(map[int32][]int32)

	for _, proc := range procs {
		procMap[proc.PID] = proc
		childrenMap[proc.PPID] = append(childrenMap[proc.PPID], proc.PID)
	}

	// 找到根进程（PPID 不在当前列表中，或者 PPID 为 0/1）
	var roots []int32
	for _, proc := range procs {
		if proc.PPID == 0 || proc.PPID == 1 {
			roots = append(roots, proc.PID)
		} else if _, exists := procMap[proc.PPID]; !exists {
			roots = append(roots, proc.PID)
		}
	}

	// 如果没有找到根，使用 PID 最小的作为根
	if len(roots) == 0 && len(procs) > 0 {
		minPID := procs[0].PID
		for _, proc := range procs {
			if proc.PID < minPID {
				minPID = proc.PID
			}
		}
		roots = append(roots, minPID)
	}

	// 使用 go-pretty/list 渲染树
	l := list.NewWriter()
	l.SetStyle(list.StyleConnectedRounded)

	// 递归构建树结构
	for _, root := range roots {
		visited := make(map[int32]bool)
		buildTreeWithList(root, procMap, childrenMap, l, visited)
	}

	fmt.Println(l.Render())
	return nil
}

// buildTreeWithList 使用 list.Writer 递归构建树
//
// 参数:
//   - pid: 当前进程 ID
//   - procMap: 进程映射
//   - childrenMap: 子进程映射
//   - l: 列表写入器
//   - visited: 已访问的 PID 集合（防止循环引用）
func buildTreeWithList(pid int32, procMap map[int32]ProcInfo, childrenMap map[int32][]int32, l list.Writer, visited map[int32]bool) {
	// 检查循环引用
	if visited[pid] {
		return
	}
	visited[pid] = true

	proc, exists := procMap[pid]
	if !exists {
		return
	}

	// 处理进程名为空的情况
	name := proc.Name
	if name == "" {
		name = "-"
	}

	// 添加当前节点
	l.AppendItem(fmt.Sprintf("%d %s", proc.PID, name))

	// 递归添加子节点
	children := childrenMap[pid]
	if len(children) > 0 {
		l.Indent()
		for _, childPID := range children {
			buildTreeWithList(childPID, procMap, childrenMap, l, visited)
		}
		l.UnIndent()
	}
}

// renderJSONMode JSON 模式渲染
//
// 参数:
//   - procs: 进程信息列表
//   - stats: 统计信息
//
// 返回:
//   - error: 错误
func renderJSONMode(procs []ProcInfo, stats *ProcStats) error {
	output := struct {
		Stats     ProcStats  `json:"stats"`
		Processes []ProcInfo `json:"processes"`
	}{
		Stats:     *stats,
		Processes: procs,
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}
