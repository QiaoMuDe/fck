# Ping 命令设计文档

## 1. 功能概述

Ping 命令用于测试网络连通性，向指定主机发送 ICMP Echo Request 包并等待响应。

## 2. CLI 接口设计

### 2.1 命令定义

```go
// internal/cli/ping.go
package cli

import (
	"fmt"
	"gitee.com/MM-Q/fck/internal/commands/ping"
	"gitee.com/MM-Q/qflag"
)

var PingCmd *qflag.Cmd

var (
	pingCount    *qflag.IntFlag      // -c, --count      发送包数量
	pingInterval *qflag.DurationFlag // -i, --interval   发送间隔
	pingTimeout  *qflag.DurationFlag // -W, --timeout    超时时间
	pingSize     *qflag.IntFlag      // -s, --size       数据包大小
	pingTTL      *qflag.IntFlag      // -t, --ttl        TTL 值
	pingQuiet    *qflag.BoolFlag     // -q, --quiet      静默模式
	pingNumeric  *qflag.BoolFlag     // -n, --numeric    不进行 DNS 解析
	pingDeadline *qflag.DurationFlag // -w, --deadline   总超时时间
)

func init() {
	PingCmd = qflag.NewCmd("ping", "", qflag.ExitOnError)

	pingCount = PingCmd.Int("count", "c", "发送包数量，0 表示无限", 4)
	pingInterval = PingCmd.Duration("interval", "i", "发送间隔，如: 1s, 500ms", "1s")
	pingTimeout = PingCmd.Duration("timeout", "W", "单个包超时时间", "5s")
	pingSize = PingCmd.Int("size", "s", "数据包大小(字节)", 56)
	pingTTL = PingCmd.Int("ttl", "t", "TTL 生存时间", 64)
	pingQuiet = PingCmd.Bool("quiet", "q", "静默模式，只显示统计结果", false)
	pingNumeric = PingCmd.Bool("numeric", "n", "不进行 DNS 反向解析", false)
	pingDeadline = PingCmd.Duration("deadline", "w", "总超时时间，到达后自动停止", "")

	cmdOpts := &qflag.CmdOpts{
		Desc:        "测试网络连通性",
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s ping [options] <host>", qflag.Root.Name()),
		Notes: []string{
			"host 可以是 IP 地址或域名",
			"需要管理员/root 权限发送 ICMP 包",
			"Windows 上可能需要关闭防火墙",
			"按 Ctrl+C 可以中断 ping",
		},
		Examples: map[string]string{
			"ping 百度":         fmt.Sprintf("%s ping baidu.com", qflag.Root.Name()),
			"ping 指定次数":      fmt.Sprintf("%s ping -c 10 baidu.com", qflag.Root.Name()),
			"ping 指定间隔":      fmt.Sprintf("%s ping -i 2s baidu.com", qflag.Root.Name()),
			"ping 大包测试":      fmt.Sprintf("%s ping -s 1400 baidu.com", qflag.Root.Name()),
			"ping IP 地址":       fmt.Sprintf("%s ping 8.8.8.8", qflag.Root.Name()),
			"静默模式":          fmt.Sprintf("%s ping -q baidu.com", qflag.Root.Name()),
			"快速 ping":         fmt.Sprintf("%s ping -i 100ms -c 100 baidu.com", qflag.Root.Name()),
		},
		MutexGroups: []qflag.MutexGroup{
			{
				Name:      "count-or-deadline",
				Flags:     []string{"count", "deadline"},
				AllowNone: true,
			},
		},
	}

	if err := PingCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts error: %w", err))
	}

	PingCmd.SetRun(runPing)
}

func runPing(cmd qflag.Command) error {
	targets := cmd.Args()
	if len(targets) == 0 {
		return fmt.Errorf("please specify a host to ping")
	}
	if len(targets) > 1 {
		return fmt.Errorf("only one host can be pinged at a time")
	}

	config := ping.PingConfig{
		Host:     targets[0],
		Count:    pingCount.Get(),
		Interval: pingInterval.Get(),
		Timeout:  pingTimeout.Get(),
		Size:     pingSize.Get(),
		TTL:      pingTTL.Get(),
		Quiet:    pingQuiet.Get(),
		Numeric:  pingNumeric.Get(),
		Deadline: pingDeadline.Get(),
	}

	return ping.PingCmdMain(config)
}
```

### 2.2 标志说明

| 标志 | 短标志 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| count | -c | int | 4 | 发送包数量，0 表示无限 |
| interval | -i | duration | 1s | 发送间隔 |
| timeout | -W | duration | 5s | 单个包超时时间 |
| size | -s | int | 56 | 数据包大小(字节) |
| ttl | -t | int | 64 | TTL 生存时间 |
| quiet | -q | bool | false | 静默模式 |
| numeric | -n | bool | false | 不进行 DNS 解析 |
| deadline | -w | duration | - | 总超时时间 |

## 3. 业务逻辑设计

### 3.1 配置结构体

```go
// internal/commands/ping/cmd_ping.go
package ping

import (
	"fmt"
	"net"
	"time"
)

// PingConfig 配置结构体
type PingConfig struct {
	Host     string        // 目标主机
	Count    int           // 发送包数量
	Interval time.Duration // 发送间隔
	Timeout  time.Duration // 单个包超时
	Size     int           // 数据包大小
	TTL      int           // TTL 值
	Quiet    bool          // 静默模式
	Numeric  bool          // 不进行 DNS 解析
	Deadline time.Duration // 总超时时间
}

// PingStats 统计信息
type PingStats struct {
	Transmitted int           // 发送包数
	Received    int           // 接收包数
	Lost        int           // 丢失包数
	LossRate    float64       // 丢包率
	MinTime     time.Duration // 最小延迟
	MaxTime     time.Duration // 最大延迟
	AvgTime     time.Duration // 平均延迟
	StdDev      time.Duration // 标准差
}
```

### 3.2 主函数流程

```go
// PingCmdMain 执行 ping 命令
func PingCmdMain(config PingConfig) error {
	// 1. 解析目标地址
	ipAddr, err := resolveHost(config.Host, config.Numeric)
	if err != nil {
		return fmt.Errorf("failed to resolve host: %w", err)
	}

	// 2. 创建 ICMP 连接
	pinger, err := createPinger(ipAddr, config)
	if err != nil {
		return fmt.Errorf("failed to create pinger: %w", err)
	}
	defer pinger.Close()

	// 3. 打印开始信息
	if !config.Quiet {
		printStartInfo(config.Host, ipAddr, config.Size)
	}

	// 4. 执行 ping 循环
	stats, err := runPingLoop(pinger, config)
	if err != nil {
		return err
	}

	// 5. 打印统计结果
	printStats(config.Host, stats)

	return nil
}

	// 3. 打印开始信息
	if !config.Quiet {
		printStartInfo(config.Host, ipAddr, config.Size)
	}

	// 4. 执行 ping 循环
	stats, err := runPingLoop(pinger, config)
	if err != nil {
		return err
	}

	// 5. 打印统计结果
	printStats(config.Host, stats)

	return nil
}
```

### 3.3 核心功能实现

```go
// resolveHost 解析主机地址
func resolveHost(host string, numeric bool) (*net.IPAddr, error) {
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

	return nil, fmt.Errorf("no IP address found for host: %s", host)
}

// createPinger 创建 ICMP pinger
func createPinger(addr *net.IPAddr, config PingConfig) (*icmp.Pinger, error) {
	// 使用 golang.org/x/net/icmp 或第三方库
	// 或平台特定的实现
}

// runPingLoop 执行 ping 循环
func runPingLoop(pinger *icmp.Pinger, config PingConfig) (*PingStats, error) {
	stats := &PingStats{
		MinTime: time.Duration(1<<63 - 1),
	}

	deadline := time.Now().Add(config.Deadline)
	interval := time.NewTicker(config.Interval)
	defer interval.Stop()

	seq := 0
	var times []time.Duration

	for {
		// 检查计数
		if config.Count > 0 && seq >= config.Count {
			break
		}

		// 检查截止时间
		if config.Deadline > 0 && time.Now().After(deadline) {
			break
		}

		// 发送 ping
		start := time.Now()
		err := pinger.Send(seq)
		if err != nil {
			return nil, err
		}
		stats.Transmitted++

		// 等待响应
		reply, err := pinger.Receive(config.Timeout)
		duration := time.Since(start)

		if err == nil && reply != nil {
			stats.Received++
			times = append(times, duration)

			if duration < stats.MinTime {
				stats.MinTime = duration
			}
			if duration > stats.MaxTime {
				stats.MaxTime = duration
			}

			if !config.Quiet {
				printReply(seq, reply, duration)
			}
		} else {
			stats.Lost++
			if !config.Quiet {
				fmt.Printf("Request timeout for icmp_seq %d\n", seq)
			}
		}

		seq++

		// 等待下一个间隔
		if config.Count == 0 || seq < config.Count {
			<-interval.C
		}
	}

	// 计算统计
	if len(times) > 0 {
		stats.AvgTime = calculateAvg(times)
		stats.StdDev = calculateStdDev(times, stats.AvgTime)
	}

	if stats.Transmitted > 0 {
		stats.LossRate = float64(stats.Lost) / float64(stats.Transmitted) * 100
	}

	return stats, nil
}
```

## 4. 输出格式

### 4.1 开始信息

```
PING baidu.com (110.242.68.66): 56 data bytes
```

### 4.2 正常响应

```
64 bytes from 110.242.68.66: icmp_seq=0 ttl=52 time=24.5 ms
64 bytes from 110.242.68.66: icmp_seq=1 ttl=52 time=23.8 ms
64 bytes from 110.242.68.66: icmp_seq=2 ttl=52 time=25.1 ms
```

### 4.3 超时

```
Request timeout for icmp_seq=3
```

### 4.4 统计结果

```
--- baidu.com ping statistics ---
4 packets transmitted, 4 received, 0% packet loss, time 3005ms
rtt min/avg/max/mdev = 23.8/24.5/25.1/0.5 ms
```

## 5. 技术选型

### 5.1 ICMP 实现方案

方案 1: 使用 `golang.org/x/net/icmp`
- 优点：官方库，跨平台
- 缺点：需要 root 权限，Windows 支持有限

方案 2: 使用 `github.com/go-ping/ping`
- 优点：功能完善，支持非特权 ping
- 缺点：第三方依赖

方案 3: 使用系统 ping 命令
- 优点：简单，无需处理权限
- 缺点：依赖外部命令，输出解析复杂

**推荐**：方案 2 `go-ping/ping`，功能完善且支持非特权模式。

### 5.2 依赖

```go
import (
	"github.com/go-ping/ping"
)
```

## 6. 文件结构

```
internal/
├── commands/
│   └── ping/
│       └── cmd_ping.go       # 业务逻辑
└── cli/
    └── ping.go               # CLI 定义
```

## 7. 错误处理

| 错误场景 | 错误信息 |
|---------|---------|
| 主机名解析失败 | `failed to resolve host: %v` |
| 权限不足 | `permission denied, try running with sudo/administrator` |
| 网络不可达 | `network is unreachable` |
| 主机不可达 | `host unreachable` |
| 超时 | `Request timeout` |

## 8. 测试用例

### 8.1 单元测试

```go
func TestResolveHost(t *testing.T) {
	// 测试 IP 地址解析
	// 测试域名解析
	// 测试无效主机名
}

func TestPingStats(t *testing.T) {
	// 测试统计计算
}
```

### 8.2 集成测试

```bash
# 测试本地回环
fck ping 127.0.0.1 -c 4

# 测试公网地址
fck ping 8.8.8.8 -c 4

# 测试域名
fck ping baidu.com -c 4

# 测试大包
fck ping baidu.com -s 1400 -c 4

# 测试快速 ping
fck ping baidu.com -i 100ms -c 100
```

## 9. 注意事项

1. **权限问题**：ICMP 需要 root/admin 权限，考虑使用非特权模式
2. **信号处理**：支持 Ctrl+C 中断
3. **并发安全**：避免多个 ping 实例冲突
4. **平台兼容**：Windows/Linux/macOS 差异处理
5. **防火墙**：提示用户可能需要关闭防火墙
