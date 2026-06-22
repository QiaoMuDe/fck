// Package ifconfig 实现了网络接口信息查看功能。
// 该文件提供了查看系统网络接口配置信息的功能，支持默认表格输出、简略输出和 JSON 输出模式。
package ifconfig

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"gitee.com/MM-Q/fck/internal/types"
	"gitee.com/MM-Q/go-kit/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	gopsutilNet "github.com/shirou/gopsutil/v3/net"
)

// IfconfigConfig ifconfig 命令配置结构体
type IfconfigConfig struct {
	All        bool     // -a: 显示所有接口（含虚拟网卡）
	Short      bool     // -s: 简略输出
	JSON       bool     // -j: JSON 输出
	Stats      bool     // --stats: 显示流量统计
	TableStyle string   // -ts: 表格样式
	Interfaces []string // 位置参数：指定接口名（为空表示全部）
}

// interfaceInfo 网络接口信息结构体
type interfaceInfo struct {
	Name      string // 接口名称
	Index     int    // 接口索引
	MTU       int    // MTU 大小
	MAC       string // MAC 地址
	Status    string // "UP" 或 "DOWN"
	IPv4      string // 第一个 IPv4 地址（含 CIDR）
	IPv6      string // 第一个 IPv6 地址（含 CIDR）
	Speed     string // 链路速度（Mbps），取不到显示 "-"
	TXBytes   uint64 // 发送字节数
	RXBytes   uint64 // 接收字节数
	TXPackets uint64 // 发送包数
	RXPackets uint64 // 接收包数
	TXErrs    uint64 // 发送错误数
	RXErrs    uint64 // 接收错误数
}

// IfconfigCmdMain ifconfig 命令主入口
//
// 参数:
//   - config: 命令配置
//
// 返回:
//   - error: 执行错误
func IfconfigCmdMain(config IfconfigConfig) error {
	interfaces, err := collectInterfaces(config)
	if err != nil {
		return fmt.Errorf("failed to collect interface information: %w", err)
	}

	if len(interfaces) == 0 {
		fmt.Println("No network interfaces found")
		return nil
	}

	switch {
	case config.JSON:
		return renderJSON(interfaces)
	case config.Short:
		renderShort(interfaces, config)
	default:
		renderTable(interfaces, config)
	}

	return nil
}

// virtualKeywords 虚拟网卡关键词
var virtualKeywords = []string{
	"Hyper-V", "VirtualBox", "VMware", "WSL",
	"Bluetooth", "Loopback", "Tailscale", "Docker",
}

// isVirtualAdapter 判断是否为虚拟网卡
//
// 参数:
//   - name: 接口名称
//
// 返回:
//   - bool: 是否为虚拟网卡
func isVirtualAdapter(name string) bool {
	nameLower := strings.ToLower(name)
	for _, keyword := range virtualKeywords {
		if strings.Contains(nameLower, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

// collectInterfaces 收集网络接口信息
//
// 参数:
//   - config: 命令配置
//
// 返回:
//   - []interfaceInfo: 接口信息列表
//   - error: 错误
func collectInterfaces(config IfconfigConfig) ([]interfaceInfo, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get network interfaces: %w", err)
	}

	// 构建 IO 计数器映射（按接口名索引）
	ioCounters, ioErr := gopsutilNet.IOCounters(true)
	var ioMap map[string]gopsutilNet.IOCountersStat
	if ioErr == nil {
		ioMap = make(map[string]gopsutilNet.IOCountersStat, len(ioCounters))
		for _, counter := range ioCounters {
			ioMap[counter.Name] = counter
		}
	}

	var result []interfaceInfo

	for _, iface := range interfaces {
		// 过滤虚拟网卡
		if !config.All && isVirtualAdapter(iface.Name) {
			continue
		}

		// 按指定接口名过滤
		if len(config.Interfaces) > 0 {
			matched := false
			for _, name := range config.Interfaces {
				if strings.EqualFold(iface.Name, name) {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}

		info := interfaceInfo{
			Name:  iface.Name,
			Index: iface.Index,
			MTU:   iface.MTU,
			MAC:   "N/A",
			Speed: "-",
		}

		// 判断接口状态
		if iface.Flags&net.FlagUp != 0 {
			info.Status = "UP"
		} else {
			info.Status = "DOWN"
		}

		// 获取 MAC 地址
		if len(iface.HardwareAddr) > 0 {
			info.MAC = iface.HardwareAddr.String()
		}

		// 获取 IP 地址
		addrs, addrErr := iface.Addrs()
		if addrErr == nil {
			for _, addr := range addrs {
				ipnet, ok := addr.(*net.IPNet)
				if !ok {
					continue
				}
				if ipnet.IP.To4() != nil {
					if info.IPv4 == "" {
						info.IPv4 = ipnet.String()
					}
				} else if info.IPv6 == "" {
					info.IPv6 = ipnet.String()
				}
			}
		}

		// 获取流量统计
		if ioMap != nil {
			if counter, ok := ioMap[iface.Name]; ok {
				info.TXBytes = counter.BytesSent
				info.RXBytes = counter.BytesRecv
				info.TXPackets = counter.PacketsSent
				info.RXPackets = counter.PacketsRecv
				info.TXErrs = counter.Errout
				info.RXErrs = counter.Errin
			}
		}

		result = append(result, info)
	}

	return result, nil
}

// renderTable 默认表格模式渲染
//
// 参数:
//   - interfaces: 接口信息列表
//   - config: 命令配置
func renderTable(interfaces []interfaceInfo, config IfconfigConfig) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// 设置表头
	headers := table.Row{"Name", "MTU", "MAC", "Status", "IPv4", "IPv6", "Speed"}
	if config.Stats {
		headers = append(headers,
			"TX Bytes", "RX Bytes",
			"TX Packets", "RX Packets",
			"TX Errs", "RX Errs",
		)
	}
	t.AppendHeader(headers)

	// 添加数据行
	for _, iface := range interfaces {
		row := table.Row{
			iface.Name,
			iface.MTU,
			iface.MAC,
			iface.Status,
			iface.IPv4,
			iface.IPv6,
			iface.Speed,
		}
		if config.Stats {
			row = append(row,
				utils.FormatBytes(int64(iface.TXBytes)),
				utils.FormatBytes(int64(iface.RXBytes)),
				iface.TXPackets,
				iface.RXPackets,
				iface.TXErrs,
				iface.RXErrs,
			)
		}
		t.AppendRow(row)
	}

	// 应用表格样式
	style := table.StyleDefault
	if config.TableStyle != "" {
		if s, ok := types.TableStyleMap[config.TableStyle]; ok {
			style = s
		}
	}
	t.SetStyle(style)

	t.Render()
}

// renderShort 简略模式渲染
//
// 参数:
//   - interfaces: 接口信息列表
//   - config: 命令配置
func renderShort(interfaces []interfaceInfo, config IfconfigConfig) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.AppendHeader(table.Row{"Name", "Status", "IPv4"})

	for _, iface := range interfaces {
		t.AppendRow(table.Row{iface.Name, iface.Status, iface.IPv4})
	}

	// 简略模式：仅当明确指定了非 "none" 样式时才显示边框
	if config.TableStyle != "" && config.TableStyle != "none" {
		if s, ok := types.TableStyleMap[config.TableStyle]; ok {
			t.SetStyle(s)
		} else {
			t.SetStyle(table.StyleDefault)
		}
	} else {
		// 默认无边框
		t.SetStyle(table.Style{
			Box: table.BoxStyle{
				PaddingLeft:  " ",
				PaddingRight: " ",
			},
		})
	}

	t.Render()
}

// renderJSON JSON 模式渲染
//
// 参数:
//   - interfaces: 接口信息列表
//
// 返回:
//   - error: 错误
func renderJSON(interfaces []interfaceInfo) error {
	data, err := json.MarshalIndent(interfaces, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(data))
	return nil
}
