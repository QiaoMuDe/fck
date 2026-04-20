// Package port 实现了端口查看功能。
// 该文件提供了查看系统网络连接和端口占用情况的功能，支持 TCP/UDP 协议和进程信息查询。
package port

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"gitee.com/MM-Q/fck/internal/types"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// PortConfig 端口命令配置结构体
type PortConfig struct {
	// 目标端口列表
	Ports []int
	// 只显示 TCP
	TCPOnly bool
	// 只显示 UDP
	UDPOnly bool
	// 连接状态过滤
	State string
	// 指定 PID
	PID int32
	// 指定进程名
	ProcessName string
	// 简洁模式
	ListMode bool
	// 表格样式
	TableStyle string
	// 只显示监听端口
	Listening bool
}

// PortInfo 端口信息结构体
type PortInfo struct {
	Protocol    string
	LocalAddr   string
	LocalPort   uint32
	RemoteAddr  string
	RemotePort  uint32
	State       string
	PID         int32
	ProcessName string
}

// PortStats 统计信息结构体
type PortStats struct {
	TotalCount    int
	TCPCount      int
	UDPCount      int
	ListenCount   int
	ProcessErrors int
}

// ProcessNameCache 进程名缓存结构体
type ProcessNameCache struct {
	cache map[int32]string
	stats *PortStats
}

// NewProcessNameCache 创建进程名缓存
//
// 参数:
//   - stats: 统计信息指针，用于记录错误
//
// 返回:
//   - *ProcessNameCache: 进程名缓存实例
func NewProcessNameCache(stats *PortStats) *ProcessNameCache {
	return &ProcessNameCache{
		cache: make(map[int32]string),
		stats: stats,
	}
}

// GetProcessName 获取进程名，带缓存机制
//
// 参数:
//   - pid: 进程 ID
//
// 返回:
//   - string: 进程名，如果获取失败返回 "-"
func (p *ProcessNameCache) GetProcessName(pid int32) string {
	// PID 无效
	if pid <= 0 {
		return "-"
	}

	// 先查缓存
	if name, ok := p.cache[pid]; ok {
		return name
	}

	// 缓存未命中，查询系统
	processName := "-"
	proc, err := process.NewProcess(pid)
	if err == nil {
		name, err := proc.Name()
		if err == nil {
			processName = name
		}
	} else {
		// 记录错误
		if p.stats != nil {
			p.stats.ProcessErrors++
		}
	}

	// 写入缓存
	p.cache[pid] = processName
	return processName
}

// formatAddr 格式化地址，将 0.0.0.0 和 :: 显示为 *
//
// 参数:
//   - addr: 原始地址
//   - port: 端口号
//
// 返回:
//   - string: 格式化后的地址
func formatAddr(addr string, port uint32) string {
	// 将 0.0.0.0 (IPv4) 和 :: (IPv6) 显示为 *
	if addr == "0.0.0.0" || addr == "::" {
		return fmt.Sprintf("*:%d", port)
	}
	return fmt.Sprintf("%s:%d", addr, port)
}

// PortCmdMain 端口命令主入口
//
// 参数:
//   - config: 端口命令配置
//
// 返回:
//   - error: 执行错误
func PortCmdMain(config *PortConfig) error {
	// 收集端口信息
	ports, stats, warnings, err := collectPorts(config)
	if err != nil {
		return fmt.Errorf("failed to collect port information: %w", err)
	}

	// 打印警告信息
	for _, warning := range warnings {
		fmt.Fprintf(os.Stderr, "warn: %s\n", warning)
	}

	if len(ports) == 0 {
		fmt.Println("No matching ports found")
		return nil
	}

	// 根据模式渲染输出
	var renderErr error
	if config.ListMode {
		renderErr = renderListMode(ports)
	} else {
		renderErr = renderTableMode(ports, stats, config)
	}

	if renderErr != nil {
		return renderErr
	}

	// 如果有进程信息获取失败，提示用户可能需要提升权限
	if stats.ProcessErrors > 0 {
		fmt.Fprintf(os.Stderr, "\nwarn: Failed to get process info for %d connections (may need elevated privileges)\n", stats.ProcessErrors)
	}

	return nil
}

// collectPorts 收集端口信息
//
// 参数:
//   - config: 端口命令配置
//
// 返回:
//   - []PortInfo: 端口信息列表
//   - *PortStats: 统计信息
//   - []string: 警告信息列表
//   - error: 错误
func collectPorts(config *PortConfig) ([]PortInfo, *PortStats, []string, error) {
	var allPorts []PortInfo
	stats := &PortStats{}
	var warnings []string

	// 创建进程名缓存
	processCache := NewProcessNameCache(stats)

	// 确定要查询的协议
	protocols := []string{"tcp", "udp"}
	if config.TCPOnly {
		protocols = []string{"tcp"}
	} else if config.UDPOnly {
		protocols = []string{"udp"}
	}

	// 查询每种协议的连接
	for _, proto := range protocols {
		connections, err := net.Connections(proto)
		if err != nil {
			// 收集错误信息但继续处理其他协议
			warnings = append(warnings, fmt.Sprintf("failed to get %s connections: %v", proto, err))
			continue
		}

		for _, conn := range connections {
			// 过滤监听端口
			// UDP 是无连接协议，没有 LISTEN 状态，只对 TCP 进行过滤
			if config.Listening && proto == "tcp" && conn.Status != "LISTEN" {
				continue
			}

			// 过滤指定端口
			if len(config.Ports) > 0 {
				found := false
				for _, port := range config.Ports {
					if int(conn.Laddr.Port) == port {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}

			// 过滤指定 PID
			if config.PID > 0 && conn.Pid != config.PID {
				continue
			}

			// 过滤指定状态
			if config.State != "" && !strings.EqualFold(conn.Status, config.State) {
				continue
			}

			// 获取进程名（使用缓存）
			processName := processCache.GetProcessName(conn.Pid)

			// 过滤指定进程名
			if config.ProcessName != "" && !strings.Contains(strings.ToLower(processName), strings.ToLower(config.ProcessName)) {
				continue
			}

			// 构建端口信息
			// 根据地址判断是 IPv4 还是 IPv6
			protocol := strings.ToUpper(proto)
			if conn.Laddr.IP == "::" || strings.Contains(conn.Laddr.IP, ":") {
				protocol += "6" // TCP6 / UDP6
			} else {
				protocol += "4" // TCP4 / UDP4
			}

			portInfo := PortInfo{
				Protocol:    protocol,
				LocalAddr:   conn.Laddr.IP,
				LocalPort:   conn.Laddr.Port,
				RemoteAddr:  conn.Raddr.IP,
				RemotePort:  conn.Raddr.Port,
				State:       conn.Status,
				PID:         conn.Pid,
				ProcessName: processName,
			}

			allPorts = append(allPorts, portInfo)

			// 更新统计
			stats.TotalCount++
			if proto == "tcp" {
				stats.TCPCount++
			} else {
				stats.UDPCount++
			}
			if conn.Status == "LISTEN" {
				stats.ListenCount++
			}
		}
	}

	// 按端口排序
	sort.Slice(allPorts, func(i, j int) bool {
		if allPorts[i].LocalPort != allPorts[j].LocalPort {
			return allPorts[i].LocalPort < allPorts[j].LocalPort
		}
		return allPorts[i].Protocol < allPorts[j].Protocol
	})

	return allPorts, stats, warnings, nil
}

// renderTableMode 表格模式渲染
//
// 参数:
//   - ports: 端口信息列表
//   - stats: 统计信息
//   - config: 端口命令配置
//
// 返回:
//   - error: 错误
func renderTableMode(ports []PortInfo, stats *PortStats, config *PortConfig) error {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// 设置表头
	t.AppendHeader(table.Row{"Protocol", "Local Address", "Remote Address", "State", "PID", "Process"})

	// 添加数据行
	for _, port := range ports {
		localAddr := formatAddr(port.LocalAddr, port.LocalPort)
		remoteAddr := "-"
		if port.RemoteAddr != "" {
			remoteAddr = formatAddr(port.RemoteAddr, port.RemotePort)
		}

		pidStr := "-"
		if port.PID > 0 {
			pidStr = strconv.Itoa(int(port.PID))
		}

		t.AppendRow(table.Row{port.Protocol, localAddr, remoteAddr, port.State, pidStr, port.ProcessName})
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
		fmt.Sprintf("TCP: %d", stats.TCPCount),
		fmt.Sprintf("UDP: %d", stats.UDPCount),
		fmt.Sprintf("Listen: %d", stats.ListenCount),
		"",
		"",
	})

	t.Render()
	return nil
}

// renderListMode 简洁模式渲染
//
// 参数:
//   - ports: 端口信息列表
//
// 返回:
//   - error: 错误
func renderListMode(ports []PortInfo) error {
	for _, port := range ports {
		localAddr := formatAddr(port.LocalAddr, port.LocalPort)

		var output string
		if port.PID > 0 && port.ProcessName != "-" {
			output = fmt.Sprintf("%s (%s, PID:%d)", localAddr, port.ProcessName, port.PID)
		} else {
			output = localAddr
		}

		output = fmt.Sprintf("[%s] %s", port.Protocol, output)
		fmt.Println(output)
	}
	return nil
}
