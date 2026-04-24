package tcp

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

// outputListenJSON 以JSON格式输出监听接收的数据
//
// 参数:
//   - connID: 连接ID
//   - remoteAddr: 远程地址
//   - startTime: 连接开始时间
//   - duration: 连接持续时间
//   - totalBytes: 接收的总字节数
//   - data: 接收的数据列表
//   - hasError: 是否有错误
//
// 返回值:
//   - error: JSON编码错误
func outputListenJSON(connID int64, remoteAddr string, startTime time.Time, duration time.Duration, totalBytes int, data []string, hasError bool) error {
	result := map[string]interface{}{
		"conn_id":     connID,
		"remote_addr": remoteAddr,
		"start_time":  startTime.Format(time.RFC3339),
		"duration_ms": float64(duration.Microseconds()) / 1000.0,
		"total_bytes": totalBytes,
		"data":        data,
		"has_error":   hasError,
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

// printConnectStats 打印连接测试统计信息
// 显示连接次数、成功/失败数量、响应时间统计
//
// 参数:
//   - stats: 连接统计信息
//
// 输出格式:
//
//	--- Summary ---
//	N attempts, M success, K failed
//	min/avg/max = x.xxx/y.yyy/z.zzz ms
//
// 说明:
//   - 仅在多次连接时调用
//   - 仅统计成功的连接计算平均时间
func printConnectStats(stats *ConnectionStats) {
	fmt.Printf("\n--- Summary ---\n")
	fmt.Printf("%d attempts, %d success, %d failed\n",
		stats.Attempted, stats.Succeeded, stats.Failed)

	if stats.Succeeded > 0 {
		avgTime := stats.TotalTime / time.Duration(stats.Succeeded)
		fmt.Printf("min/avg/max = %.3f/%.3f/%.3f ms\n",
			float64(stats.MinTime.Microseconds())/1000.0,
			float64(avgTime.Microseconds())/1000.0,
			float64(stats.MaxTime.Microseconds())/1000.0)
	}
}

// printScanResults 打印端口扫描结果
// 以表格形式显示每个端口的状态、服务猜测和响应时间
//
// 参数:
//   - config: TCP配置,控制输出选项
//   - results: 端口扫描结果列表
//   - duration: 总扫描耗时
//
// 输出格式:
//
//	PORT     STATE    SERVICE    TIME(ms)
//	80       open     http       15.234
//	443      open     https      18.567
//
// 功能特性:
//   - 支持仅显示开放端口(-o选项)
//   - 自动统计开放/关闭端口数量
//   - 显示总扫描耗时
func printScanResults(config TcpConfig, results []PortResult, duration time.Duration) {
	if !config.Quiet {
		fmt.Printf("\nPORT     STATE    SERVICE    TIME(ms)\n")
	}

	open := 0
	closed := 0

	for _, result := range results {
		if result.State == "open" {
			open++
			if !config.OpenOnly {
				fmt.Printf("%-8d %-8s %-10s %.3f\n",
					result.Port, result.State, result.Service,
					result.TimeMs)
			}
		} else {
			closed++
			if !config.OpenOnly {
				fmt.Printf("%-8d %-8s %-10s %.3f\n",
					result.Port, result.State, result.Service,
					result.TimeMs)
			}
		}
	}

	if config.OpenOnly {
		for _, result := range results {
			if result.State == "open" {
				fmt.Printf("%-8d %-8s %-10s %.3f\n",
					result.Port, result.State, result.Service,
					result.TimeMs)
			}
		}
	}

	fmt.Printf("\nScan completed: %d ports scanned, %d open, %d closed/filtered\n",
		len(results), open, closed)
	fmt.Printf("Time taken: %.3fs\n", duration.Seconds())
}

// outputConnectJSON 以JSON格式输出连接测试结果
// 包含每次连接的详细结果和汇总统计
//
// 参数:
//   - config: TCP配置,包含目标主机和端口
//   - stats: 连接统计信息
//   - results: 每次连接的详细结果
//
// 返回值:
//   - error: JSON编码错误
//
// 输出示例:
//
//	{
//	  "host": "baidu.com",
//	  "port": 80,
//	  "results": [...],
//	  "summary": {...}
//	}
func outputConnectJSON(config TcpConfig, stats *ConnectionStats, results []ConnectResult) error {
	avgTime := time.Duration(0)
	if stats.Succeeded > 0 {
		avgTime = stats.TotalTime / time.Duration(stats.Succeeded)
	}

	result := map[string]interface{}{
		"host":    config.Host,
		"port":    config.Port,
		"results": results,
		"summary": map[string]interface{}{
			"attempted": stats.Attempted,
			"success":   stats.Succeeded,
			"failed":    stats.Failed,
			"min_ms":    float64(stats.MinTime.Microseconds()) / 1000.0,
			"avg_ms":    float64(avgTime.Microseconds()) / 1000.0,
			"max_ms":    float64(stats.MaxTime.Microseconds()) / 1000.0,
		},
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

// outputBannerJSON 以JSON格式输出Banner获取结果
//
// 参数:
//   - config: TCP配置,包含目标主机和端口
//   - ipAddr: 目标IP地址
//   - banner: 获取到的Banner内容
//   - state: 状态(success/timeout/refused等)
//   - duration: 连接耗时
//
// 返回值:
//   - error: JSON编码错误
func outputBannerJSON(config TcpConfig, ipAddr net.IP, banner string, state string, duration time.Duration) error {
	result := map[string]interface{}{
		"host":      config.Host,
		"ip":        ipAddr.String(),
		"port":      config.Port,
		"state":     state,
		"time_ms":   float64(duration.Microseconds()) / 1000.0,
		"banner":    banner,
		"wait_time": config.Wait.String(),
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

// outputDataJSON 以JSON格式输出数据发送结果
//
// 参数:
//   - config: TCP配置,包含目标主机和端口
//   - ipAddr: 目标IP地址
//   - data: 发送的数据内容
//   - dataSource: 数据来源(file路径或"inline")
//   - response: 接收到的响应
//   - state: 状态(success/send_failed/timeout等)
//   - dialDuration: 连接耗时
//   - sendDuration: 发送耗时
//
// 返回值:
//   - error: JSON编码错误
func outputDataJSON(config TcpConfig, ipAddr net.IP, data []byte, dataSource string, response string, state string, dialDuration time.Duration, sendDuration time.Duration) error {
	result := map[string]interface{}{
		"host":        config.Host,
		"ip":          ipAddr.String(),
		"port":        config.Port,
		"state":       state,
		"data_source": dataSource,
		"data_size":   len(data),
		"dial_ms":     float64(dialDuration.Microseconds()) / 1000.0,
		"send_ms":     float64(sendDuration.Microseconds()) / 1000.0,
		"wait_time":   config.Wait.String(),
		"response":    response,
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

// outputScanJSON 以JSON格式输出端口扫描结果
// 包含每个端口的详细结果和汇总统计
//
// 参数:
//   - config: TCP配置,包含目标主机
//   - ipAddr: 目标IP地址
//   - results: 端口扫描结果列表
//   - duration: 总扫描耗时
//
// 返回值:
//   - error: JSON编码错误
//
// 输出示例:
//
//	{
//	  "host": "target.com",
//	  "ip": "192.168.1.1",
//	  "ports": [...],
//	  "total": 100,
//	  "open": 3,
//	  ...
//	}
func outputScanJSON(config TcpConfig, ipAddr net.IP, results []PortResult, duration time.Duration) error {
	open := 0
	for _, r := range results {
		if r.State == "open" {
			open++
		}
	}

	result := ScanResult{
		Host:      config.Host,
		IP:        ipAddr.String(),
		Ports:     results,
		Total:     len(results),
		Open:      open,
		Closed:    len(results) - open,
		TimeTaken: fmt.Sprintf("%.3fs", duration.Seconds()),
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}
