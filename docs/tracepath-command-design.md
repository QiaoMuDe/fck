# Tracepath 命令设计方案

## 命令概述

类似于 Linux `tracepath` 的网络路由追踪工具，用于追踪数据包从本地到目标主机的网络路径，显示每一跳的延迟和主机信息。

## 与现有命令的区别

- **ping**: 测试单个主机的连通性和延迟
- **tcp**: TCP 端口扫描和连接测试
- **tracepath**: 追踪网络路由路径，显示每一跳信息

## 命令结构

```
fck tracepath [选项] <目标主机>
```

## 功能特性

### 核心功能
1. **路由追踪** - 使用递增 TTL 发现路径上的每一跳路由器
2. **延迟测量** - 测量到每一跳的往返时间 (RTT)
3. **主机名解析** - 自动解析每一跳的域名
4. **MTU 探测** - 发现路径 MTU (最大传输单元)
5. **多协议支持** - 支持 ICMP、UDP、TCP 探测

### 探测模式
- **ICMP 模式** (默认) - 使用 ICMP Echo Request
- **UDP 模式** - 使用 UDP 数据包，某些防火墙允许通过
- **TCP 模式** - 使用 TCP SYN，可穿透部分防火墙

### 输出信息
- 跳数 (Hop)
- 主机 IP 地址
- 主机名 (反向解析)
- 往返时间 (RTT) - 发送 3 个探测包，显示最小/平均/最大
- MTU 限制信息
- 超时标记

## CLI 定义

### 文件位置
`internal/cli/tracepath.go`

### 标志定义

| 长名称 | 短名称 | 类型 | 默认值 | 说明 |
|--------|--------|------|--------|------|
| icmp | I | bool | true | 使用 ICMP 探测 |
| udp | U | bool | false | 使用 UDP 探测 |
| tcp | T | bool | false | 使用 TCP 探测 |
| port | p | int | 0 | TCP/UDP 目标端口 (0=随机) |
| max-hops | m | int | 30 | 最大跳数 |
| first-ttl | f | int | 1 | 起始 TTL |
| timeout | t | duration | 5s | 每跳超时时间 |
| probes | q | int | 3 | 每跳探测次数 |
| resolve | n | bool | true | 解析主机名 |
| mtu | M | bool | false | 探测路径 MTU |
| min-mtu | | int | 576 | MTU 探测最小值 |
| max-mtu | | int | 65535 | MTU 探测最大值 |
| source | s | string | "" | 源 IP 地址 |
| interface | i | string | "" | 网络接口 |
| no-dns | | bool | false | 禁用 DNS 解析 |
| numeric | | bool | false | 仅显示 IP，不解析主机名 |
| backup | b | bool | false | 使用备用路径探测 |
| asym | a | bool | false | 显示非对称跳数 |

### 代码实现

```go
package cli

import (
	"fmt"
	"time"

	"gitee.com/MM-Q/fck/internal/commands/tracepath"
	"gitee.com/MM-Q/qflag"
)

var TracepathCmd *qflag.Cmd

var (
	tpICMP      *qflag.BoolFlag
	tpUDP       *qflag.BoolFlag
	tpTCP       *qflag.BoolFlag
	tpPort      *qflag.IntFlag
	tpMaxHops   *qflag.IntFlag
	tpFirstTTL  *qflag.IntFlag
	tpTimeout   *qflag.DurationFlag
	tpProbes    *qflag.IntFlag
	tpResolve   *qflag.BoolFlag
	tpMTU       *qflag.BoolFlag
	tpMinMTU    *qflag.IntFlag
	tpMaxMTU    *qflag.IntFlag
	tpSource    *qflag.StringFlag
	tpInterface *qflag.StringFlag
	tpNoDNS     *qflag.BoolFlag
	tpNumeric   *qflag.BoolFlag
	tpBackup    *qflag.BoolFlag
	tpAsym      *qflag.BoolFlag
)

func init() {
	TracepathCmd = qflag.NewCmd("tracepath", "tp", qflag.ExitOnError)

	tpICMP = TracepathCmd.Bool("icmp", "I", "使用 ICMP 探测", true)
	tpUDP = TracepathCmd.Bool("udp", "U", "使用 UDP 探测", false)
	tpTCP = TracepathCmd.Bool("tcp", "T", "使用 TCP 探测", false)
	tpPort = TracepathCmd.Int("port", "p", "TCP/UDP 目标端口 (0=随机)", 0)
	tpMaxHops = TracepathCmd.Int("max-hops", "m", "最大跳数", 30)
	tpFirstTTL = TracepathCmd.Int("first-ttl", "f", "起始 TTL", 1)
	tpTimeout = TracepathCmd.Duration("timeout", "t", "每跳超时时间", time.Second*5)
	tpProbes = TracepathCmd.Int("probes", "q", "每跳探测次数", 3)
	tpResolve = TracepathCmd.Bool("resolve", "n", "解析主机名", true)
	tpMTU = TracepathCmd.Bool("mtu", "M", "探测路径 MTU", false)
	tpMinMTU = TracepathCmd.Int("min-mtu", "", "MTU 探测最小值", 576)
	tpMaxMTU = TracepathCmd.Int("max-mtu", "", "MTU 探测最大值", 65535)
	tpSource = TracepathCmd.String("source", "s", "源 IP 地址", "")
	tpInterface = TracepathCmd.String("interface", "i", "网络接口", "")
	tpNoDNS = TracepathCmd.Bool("no-dns", "", "禁用 DNS 解析", false)
	tpNumeric = TracepathCmd.Bool("numeric", "", "仅显示 IP，不解析主机名", false)
	tpBackup = TracepathCmd.Bool("backup", "b", "使用备用路径探测", false)
	tpAsym = TracepathCmd.Bool("asym", "a", "显示非对称跳数", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "追踪网络路由路径",
		UsageSyntax: fmt.Sprintf("%s tracepath [选项] <目标主机>", qflag.Root.Name()),
		UseChinese:  true,
		Notes: []string{
			"目标主机可以是 IP 地址或域名",
			"默认使用 ICMP 探测，某些网络可能需要使用 UDP 或 TCP",
			"需要管理员/root 权限发送原始套接字",
			"Windows 系统可能需要关闭防火墙",
		},
		Examples: map[string]string{
			"基础路由追踪":     fmt.Sprintf("%s tracepath www.example.com", qflag.Root.Name()),
			"使用 UDP 探测":   fmt.Sprintf("%s tracepath -U 8.8.8.8", qflag.Root.Name()),
			"使用 TCP 探测":   fmt.Sprintf("%s tracepath -T -p 80 www.example.com", qflag.Root.Name()),
			"最大 20 跳":      fmt.Sprintf("%s tracepath -m 20 www.example.com", qflag.Root.Name()),
			"探测路径 MTU":    fmt.Sprintf("%s tracepath -M www.example.com", qflag.Root.Name()),
			"仅显示 IP":       fmt.Sprintf("%s tracepath --numeric www.example.com", qflag.Root.Name()),
			"指定网络接口":    fmt.Sprintf("%s tracepath -i eth0 www.example.com", qflag.Root.Name()),
		},
		MutexGroups: []qflag.MutexGroup{
			{
				Name:      "probe-mode",
				Flags:     []string{"icmp", "udp", "tcp"},
				AllowNone: true,
			},
		},
	}

	if err := TracepathCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	TracepathCmd.SetRun(runTracepath)
}

func runTracepath(cmd qflag.Command) error {
	args := cmd.Args()
	if len(args) == 0 {
		return fmt.Errorf("缺少目标主机参数")
	}

	// 确定探测模式
	probeMode := "icmp"
	switch {
	case tpUDP.Get():
		probeMode = "udp"
	case tpTCP.Get():
		probeMode = "tcp"
	}

	config := tracepath.TracepathConfig{
		Target:      args[0],
		ProbeMode:   probeMode,
		Port:        tpPort.Get(),
		MaxHops:     tpMaxHops.Get(),
		FirstTTL:    tpFirstTTL.Get(),
		Timeout:     tpTimeout.Get(),
		Probes:      tpProbes.Get(),
		Resolve:     tpResolve.Get() && !tpNoDNS.Get() && !tpNumeric.Get(),
		MTU:         tpMTU.Get(),
		MinMTU:      tpMinMTU.Get(),
		MaxMTU:      tpMaxMTU.Get(),
		Source:      tpSource.Get(),
		Interface:   tpInterface.Get(),
		Numeric:     tpNumeric.Get() || tpNoDNS.Get(),
		Backup:      tpBackup.Get(),
		Asymmetric:  tpAsym.Get(),
	}

	return tracepath.TracepathCmdMain(config)
}
```

## 业务逻辑

### 文件位置
`internal/commands/tracepath/cmd_tracepath.go`

### 配置结构体

```go
package tracepath

import "time"

// TracepathConfig 路由追踪配置
type TracepathConfig struct {
	Target     string
	ProbeMode  string        // icmp/udp/tcp
	Port       int
	MaxHops    int
	FirstTTL   int
	Timeout    time.Duration
	Probes     int
	Resolve    bool
	MTU        bool
	MinMTU     int
	MaxMTU     int
	Source     string
	Interface  string
	Numeric    bool
	Backup     bool
	Asymmetric bool
}

// HopResult 单跳结果
type HopResult struct {
	TTL       int
	IP        string
	Hostname  string
	RTT       []time.Duration // 多次探测的 RTT
	Lost      int             // 丢包数
	MTU       int             // 该跳的 MTU (如果有)
	Asymmetric int            // 非对称跳数
	Error     string
}

// TraceResult 完整追踪结果
type TraceResult struct {
	Target      string
	TargetIP    string
	Hops        []HopResult
	TotalHops   int
	Success     bool
	PathMTU     int
	Asymmetric  bool
	TotalTime   time.Duration
}
```

### 主函数

```go
// TracepathCmdMain 路由追踪主函数
func TracepathCmdMain(config TracepathConfig) error {
	// 解析目标地址
	targetIP, err := resolveTarget(config.Target)
	if err != nil {
		return fmt.Errorf("解析目标失败: %w", err)
	}

	// 创建探测器
	prober, err := createProber(config)
	if err != nil {
		return fmt.Errorf("创建探测器失败: %w", err)
	}
	defer prober.Close()

	// 执行路由追踪
	result, err := traceRoute(prober, targetIP, config)
	if err != nil {
		return fmt.Errorf("路由追踪失败: %w", err)
	}

	// MTU 探测
	if config.MTU {
		result.PathMTU = probeMTU(prober, targetIP, config)
	}

	// 输出结果
	return outputResult(result, config)
}
```

## 核心功能实现

### 1. 探测器接口

```go
// Prober 探测器接口
type Prober interface {
	// Probe 发送探测包
	// ttl: TTL 值
	// returns: 是否收到响应, 响应地址, RTT, 错误
	Probe(ttl int) (bool, string, time.Duration, error)
	
	// ProbeMTU 探测 MTU
	ProbeMTU(size int) (bool, error)
	
	// Close 关闭探测器
	Close() error
}
```

### 2. ICMP 探测器

使用原始套接字发送 ICMP Echo Request，通过递增 TTL 触发每一跳返回 ICMP Time Exceeded。

### 3. UDP 探测器

发送 UDP 到高端口，触发 ICMP Port Unreachable 或 Time Exceeded。

### 4. TCP 探测器

发送 TCP SYN，某些防火墙配置允许通过。

### 5. 输出格式

**标准输出示例:**
```
追踪到 www.example.com (93.184.216.34)，最大 30 跳

 1:  192.168.1.1 (gateway.local)        0.5ms   0.4ms   0.6ms
 2:  10.0.0.1                           1.2ms   1.1ms   1.3ms
 3:  172.16.0.1 (core-router.isp.net)   2.5ms   2.3ms   2.8ms
 4:  * * * (超时)
 5:  93.184.216.34 (www.example.com)   15.2ms  15.0ms  15.5ms

到达目标! 总跳数: 5
```

**带 MTU 探测:**
```
追踪到 www.example.com (93.184.216.34)，最大 30 跳

 1:  192.168.1.1                       0.5ms   0.4ms   0.6ms  MTU=1500
 2:  10.0.0.1                          1.2ms   1.1ms   1.3ms  MTU=1500
 ...
路径 MTU: 1500 字节
```

## 平台适配

### Windows
- 使用 `iphlpapi.dll` 或原始套接字
- 可能需要管理员权限
- 防火墙可能拦截 ICMP

### Linux
- 使用原始套接字
- 需要 root 权限或使用 capabilities

### macOS
- 类似 Linux 实现

## 使用示例

```bash
# 基础路由追踪
fck tracepath www.example.com

# 使用 UDP 探测
fck tracepath -U 8.8.8.8

# 使用 TCP 探测特定端口
fck tracepath -T -p 443 www.example.com

# 最大 15 跳
fck tracepath -m 15 www.example.com

# 探测路径 MTU
fck tracepath -M www.example.com

# 仅显示 IP，不解析主机名
fck tracepath --numeric www.example.com

# 指定网络接口
fck tracepath -i eth0 www.example.com

# 更快的探测 (减少超时和探测次数)
fck tracepath -t 2s -q 1 www.example.com
```

## 注册到根命令

在 `internal/cli/root.go` 的 `SubCmds` 中添加：

```go
SubCmds: []qflag.Command{
	// ... 其他命令
	TracepathCmd,
},
```

## 实现步骤

1. 创建业务逻辑目录 `internal/commands/tracepath/`
2. 创建 `internal/commands/tracepath/cmd_tracepath.go` - 主逻辑
3. 创建 `internal/commands/tracepath/prober.go` - 探测器接口和实现
4. 创建 `internal/commands/tracepath/output.go` - 输出格式化
5. 创建 `internal/cli/tracepath.go` - CLI 定义
6. 在 `internal/cli/root.go` 中注册命令
7. 测试编译和功能验证

## 注意事项

1. **权限要求**: 需要管理员/root 权限发送原始套接字
2. **防火墙**: 某些网络可能过滤 ICMP，需要使用 UDP/TCP 模式
3. **负载均衡**: 某些路径可能有多个路由，每次探测结果可能不同
4. **NAT**: 经过 NAT 后的跳数可能显示为同一 IP
