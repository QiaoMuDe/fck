// Package ping 实现了网络连通性测试命令
package ping

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

// PingConfig 配置结构体
type PingConfig struct {
	Host     string        // 目标主机
	Count    int           // 发送包数量
	Interval time.Duration // 发送间隔
	Timeout  time.Duration // 总超时时间
	Size     int           // 数据包大小
	TTL      int           // TTL 值
	Quiet    bool          // 静默模式
}

// PingCmdMain 执行 ping 命令
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - error: 执行错误
func PingCmdMain(config PingConfig) error {
	// 1. 解析目标地址
	ipAddr, err := resolveHost(config.Host)
	if err != nil {
		return fmt.Errorf("failed to resolve host: %w", err)
	}

	// 2. 创建 pinger
	pinger, err := probing.NewPinger(ipAddr.String())
	if err != nil {
		return fmt.Errorf("failed to create pinger: %w", err)
	}

	// 3. 配置 pinger
	pinger.Count = config.Count
	pinger.Interval = config.Interval
	pinger.Timeout = config.Timeout
	pinger.Size = config.Size
	pinger.TTL = config.TTL
	pinger.SetPrivileged(true) // 使用特权模式

	// 4. 设置用户中断标志
	var interrupted bool

	// 5. 设置回调
	if !config.Quiet {
		// 打印开始信息
		fmt.Printf("PING %s (%s): %d data bytes\n", config.Host, ipAddr.String(), config.Size)

		pinger.OnRecv = func(pkt *probing.Packet) {
			fmt.Printf("%d bytes from %s: icmp_seq=%d ttl=%d time=%.3f ms\n",
				pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.TTL, float64(pkt.Rtt.Microseconds())/1000.0)
		}

		pinger.OnDuplicateRecv = func(pkt *probing.Packet) {
			fmt.Printf("%d bytes from %s: icmp_seq=%d ttl=%d time=%.3f ms (DUP!)\n",
				pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.TTL, float64(pkt.Rtt.Microseconds())/1000.0)
		}
	}

	// 6. 设置信号处理 (Ctrl+C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		interrupted = true
		if pinger != nil {
			pinger.Stop()
		}
	}()

	// 7. 执行 ping
	err = pinger.Run()
	if err != nil {
		// 区分用户中断和真正的错误
		if interrupted {
			// 用户主动中断，不返回错误，继续打印统计结果
		} else {
			return fmt.Errorf("ping failed: %w", err)
		}
	}

	// 8. 获取统计结果并打印
	results := pinger.Statistics()
	printStats(config.Host, results)

	return nil
}

// resolveHost 解析主机地址
//
// 参数:
//   - host: 主机名或 IP 地址
//
// 返回值:
//   - *net.IPAddr: 解析后的 IP 地址
//   - error: 解析错误
func resolveHost(host string) (*net.IPAddr, error) {
	// 尝试解析为 IP
	ip := net.ParseIP(host)
	if ip != nil {
		return &net.IPAddr{IP: ip}, nil
	}

	// DNS 解析
	addrs, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}

	// 优先使用 IPv4
	for _, addr := range addrs {
		if addr.To4() != nil {
			return &net.IPAddr{IP: addr}, nil
		}
	}

	// 如果没有 IPv4，使用第一个
	if len(addrs) > 0 {
		return &net.IPAddr{IP: addrs[0]}, nil
	}

	return nil, errors.New("no IP address found for host: " + host)
}

// printStats 打印统计结果
//
// 参数:
//   - host: 目标主机
//   - stats: 统计信息
func printStats(host string, stats *probing.Statistics) {
	// 计算丢包率
	lossRate := 0.0
	if stats.PacketsSent > 0 {
		lossRate = float64(stats.PacketsSent-stats.PacketsRecv) / float64(stats.PacketsSent) * 100
	}

	fmt.Printf("\n--- %s ping statistics ---\n", host)
	fmt.Printf("%d packets transmitted, %d received, %.1f%% packet loss\n",
		stats.PacketsSent, stats.PacketsRecv, lossRate)

	if stats.PacketsRecv > 0 {
		fmt.Printf("rtt min/avg/max/mdev = %.3f/%.3f/%.3f/%.3f ms\n",
			float64(stats.MinRtt.Microseconds())/1000.0,
			float64(stats.AvgRtt.Microseconds())/1000.0,
			float64(stats.MaxRtt.Microseconds())/1000.0,
			float64(stats.StdDevRtt.Microseconds())/1000.0)
	}
}
