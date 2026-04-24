# TCP 命令设计文档

## 1. 命令概述

### 1.1 命令名称
- **中文名**: TCP 连接测试工具
- **英文名**: tcp
- **功能**: 测试 TCP 端口连通性、执行 TCP 端口扫描、发送/接收 TCP 数据

### 1.2 使用场景
- 测试服务器端口是否开放
- 检测网络连通性
- 简单的 TCP 端口扫描
- 发送 TCP 探测包
- 测试服务响应时间

## 2. 命令设计

### 2.1 命令语法

```bash
fck tcp [options] <host> <port>
```

### 2.2 参数说明

| 参数 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `host` | string | 是 | 目标主机（IP 或域名） |
| `port` | int | 是 | 目标端口号（1-65535） |

### 2.3 选项设计

| 长选项 | 短选项 | 类型 | 默认值 | 说明 |
|--------|--------|------|--------|------|
| `--timeout` | `-t` | Duration | 5s | 连接超时时间 |
| `--count` | `-c` | int | 1 | 连接次数（0 表示无限） |
| `--interval` | `-i` | Duration | 1s | 连接间隔时间 |
| `--scan` | `-s` | bool | false | 端口扫描模式 |
| `--range` | `-r` | string | "" | 端口范围（如: 80-100,443） |
| `--open` | `-o` | bool | false | 仅显示开放的端口 |
| `--banner` | `-b` | bool | false | 尝试获取服务 banner |
| `--data` | `-d` | string | "" | 发送的数据内容 |
| `--file` | `-f` | string | "" | 从文件读取发送数据 |
| `--wait` | `-w` | Duration | 0 | 等待响应时间 |
| `--quiet` | `-q` | bool | false | 静默模式 |
| `--json` | `-j` | bool | false | JSON 格式输出 |

## 3. 功能设计

### 3.1 基础连接测试

默认模式，测试单个端口的连通性：

```bash
# 测试 80 端口
fck tcp baidu.com 80

# 指定超时时间
fck tcp -t 10s baidu.com 443

# 多次连接测试
fck tcp -c 5 -i 2s baidu.com 80
```

**输出示例**:
```
Connecting to baidu.com (111.63.65.104):80...
Connected to baidu.com:80 (time=23.456ms)
Connection closed.

--- baidu.com:80 tcp statistics ---
1 connections attempted, 1 succeeded, 0 failed
rtt min/avg/max = 23.456/23.456/23.456 ms
```

### 3.2 端口扫描模式

使用 `--scan` 或 `--range` 进行端口扫描：

```bash
# 扫描连续端口范围
fck tcp -s -r 1-1000 baidu.com

# 扫描指定端口列表
fck tcp -s -r 80,443,8080,3306 baidu.com

# 仅显示开放端口
fck tcp -s -r 1-100 -o baidu.com

# JSON 格式输出
fck tcp -s -r 1-100 -j baidu.com
```

**输出示例**:
```
Scanning baidu.com (111.63.65.104) ports 1-100...
PORT     STATE    SERVICE
22/tcp   filtered ssh
80/tcp   open     http
443/tcp  open     https

Scan completed: 100 ports scanned, 2 open, 98 closed/filtered
Time taken: 12.345s
```

### 3.3 Banner 抓取

使用 `--banner` 获取服务 banner：

```bash
# 获取服务 banner
fck tcp -b baidu.com 80

# 结合等待时间
fck tcp -b -w 2s baidu.com 22
```

**输出示例**:
```
Connecting to baidu.com (111.63.65.104):80...
Connected (time=23.456ms)

--- Received Banner ---
HTTP/1.1 400 Bad Request
Server: bfe/1.0.8.18
...
```

### 3.4 数据发送模式

使用 `--data` 或 `--file` 发送自定义数据：

```bash
# 发送字符串数据
fck tcp -d "GET / HTTP/1.1\r\nHost: baidu.com\r\n\r\n" baidu.com 80

# 从文件读取数据
fck tcp -f request.txt baidu.com 80

# 发送并等待响应
fck tcp -d "ping" -w 5s baidu.com 8080
```

## 4. 互斥组设计

### 4.1 模式互斥

以下标志代表不同的工作模式，不能同时使用：

```go
MutexGroups: []qflag.MutexGroup{
    {
        Name:      "mode",
        Flags:     []string{"scan", "banner", "data", "file"},
        AllowNone: true,  // 允许都不设置（默认连接模式）
    },
}
```

### 4.2 数据源互斥

`--data` 和 `--file` 不能同时使用：

```go
MutexGroups: []qflag.MutexGroup{
    {
        Name:      "data-source",
        Flags:     []string{"data", "file"},
        AllowNone: true,
    },
}
```

## 5. 必需组设计

### 5.1 扫描模式必需参数

使用 `--scan` 时，`--range` 必须设置：

```go
RequiredGroups: []qflag.RequiredGroup{
    {
        Name:        "scan-params",
        Flags:       []string{"scan", "range"},
        Conditional: true,  // 设置了 scan 就必须设置 range
    },
}
```

## 6. 代码结构设计

### 6.1 文件结构

```
internal/
├── commands/tcp/
│   └── cmd_tcp.go          # 业务逻辑
└── cli/
    ├── tcp.go              # CLI 定义
    └── root.go             # 注册 TcpCmd
```

### 6.2 配置结构体

```go
// TcpConfig 配置结构体
type TcpConfig struct {
    Host     string        // 目标主机
    Port     int           // 目标端口
    Timeout  time.Duration // 连接超时时间
    Count    int           // 连接次数
    Interval time.Duration // 连接间隔
    Scan     bool          // 扫描模式
    Range    string        // 端口范围
    OpenOnly bool          // 仅显示开放端口
    Banner   bool          // 获取 banner
    Data     string        // 发送数据
    File     string        // 数据文件路径
    Wait     time.Duration // 等待响应时间
    Quiet    bool          // 静默模式
    Json     bool          // JSON 输出
}
```

### 6.3 统计结构体

```go
// TcpStats 连接统计
type TcpStats struct {
    Attempted   int           // 尝试连接数
    Succeeded   int           // 成功连接数
    Failed      int           // 失败连接数
    MinTime     time.Duration // 最小耗时
    MaxTime     time.Duration // 最大耗时
    TotalTime   time.Duration // 总耗时
    OpenPorts   []int         // 开放端口列表（扫描模式）
    ClosedPorts []int         // 关闭端口列表（扫描模式）
    Banner      string        // 获取的 banner
}

// PortResult 端口扫描结果
type PortResult struct {
    Port    int           // 端口号
    State   string        // 状态: open/closed/filtered
    Service string        // 服务名称（推测）
    Time    time.Duration // 响应时间
}
```

### 6.4 主函数设计

```go
// TcpCmdMain 执行 TCP 命令
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - error: 执行错误
func TcpCmdMain(config TcpConfig) error {
    // 1. 解析目标地址
    // 2. 根据模式执行不同逻辑
    //    - 扫描模式: 执行端口扫描
    //    - Banner 模式: 获取服务 banner
    //    - 数据模式: 发送并接收数据
    //    - 默认模式: 连接测试
    // 3. 输出结果
    return nil
}
```

## 7. 核心功能实现

### 7.1 端口解析

```go
// parsePortRange 解析端口范围字符串
// 支持格式: "80", "80-100", "80,443,8080", "1-100,443,8080-8090"
func parsePortRange(rangeStr string) ([]int, error) {
    // 实现端口范围解析
}
```

### 7.2 连接测试

```go
// testConnection 测试单个端口连接
func testConnection(host string, port int, timeout time.Duration) (time.Duration, error) {
    // 实现 TCP 连接测试
    // 返回连接耗时和错误
}
```

### 7.3 端口扫描

```go
// scanPorts 扫描端口范围
func scanPorts(host string, ports []int, config TcpConfig) ([]PortResult, error) {
    // 实现端口扫描逻辑
    // 支持并发扫描提高效率
}
```

### 7.4 Banner 获取

```go
// grabBanner 尝试获取服务 banner
func grabBanner(host string, port int, timeout time.Duration) (string, error) {
    // 实现 banner 抓取
    // 连接后读取响应数据
}
```

### 7.5 数据发送

```go
// sendData 发送数据并等待响应
func sendData(host string, port int, data []byte, wait time.Duration) ([]byte, error) {
    // 实现数据发送和接收
}
```

## 8. 输出格式

### 8.1 普通模式输出

```
Connecting to <host> (<ip>):<port>...
Connected to <host>:<port> (time=<duration>)
Connection closed.

--- <host>:<port> tcp statistics ---
<count> connections attempted, <succeeded> succeeded, <failed> failed
rtt min/avg/max = <min>/<avg>/<max> ms
```

### 8.2 扫描模式输出

```
Scanning <host> (<ip>) ports <range>...
PORT     STATE    SERVICE    TIME
<port>   <state>  <service>  <time>
...

Scan completed: <total> ports scanned, <open> open, <closed> closed/filtered
Time taken: <duration>
```

### 8.3 JSON 模式输出

```json
{
  "host": "baidu.com",
  "ip": "111.63.65.104",
  "ports": [
    {
      "port": 80,
      "state": "open",
      "service": "http",
      "time_ms": 23.456
    }
  ],
  "statistics": {
    "attempted": 100,
    "succeeded": 2,
    "failed": 98,
    "time_taken_ms": 12345
  }
}
```

## 9. 错误处理

### 9.1 常见错误

| 错误类型 | 错误信息 | 处理方式 |
|----------|----------|----------|
| 无效主机 | "invalid host: <host>" | 返回错误 |
| 无效端口 | "invalid port: <port>, must be 1-65535" | 返回错误 |
| 连接超时 | "connection timeout" | 记录失败 |
| 连接拒绝 | "connection refused" | 记录失败 |
| DNS 解析失败 | "failed to resolve host: <host>" | 返回错误 |
| 文件读取失败 | "failed to read file: <path>" | 返回错误 |

### 9.2 错误处理原则

- 参数错误：立即返回错误
- 连接错误：记录到统计，继续执行
- 扫描错误：单端口失败不影响其他端口

## 10. 使用示例

### 10.1 基础示例

```bash
# 测试单个端口
fck tcp baidu.com 80

# 测试多个端口（扫描）
fck tcp -s -r 80,443 baidu.com

# 测试端口范围
fck tcp -s -r 1-1000 baidu.com

# 仅显示开放端口
fck tcp -s -r 1-100 -o baidu.com
```

### 10.2 高级示例

```bash
# 获取服务 banner
fck tcp -b -w 2s baidu.com 80

# 发送 HTTP 请求
fck tcp -d "GET / HTTP/1.1\r\nHost: baidu.com\r\n\r\n" -w 5s baidu.com 80

# 从文件发送数据
fck tcp -f payload.bin -w 10s target.com 8080

# JSON 格式扫描
fck tcp -s -r 22,80,443 -j target.com | jq '.ports[] | select(.state == "open")'

# 多次连接测试（压力测试）
fck tcp -c 100 -i 100ms target.com 80
```

## 11. 注意事项

1. **权限要求**: 扫描低端口（<1024）可能需要管理员权限
2. **防火墙**: 某些端口可能被防火墙拦截，显示 filtered
3. **扫描频率**: 避免高频扫描，可能被目标视为攻击
4. **并发控制**: 扫描模式应限制并发数，避免资源耗尽
5. **超时设置**: 合理设置超时时间，避免长时间等待

## 12. 依赖库

```go
// 标准库
import (
    "encoding/json"
    "fmt"
    "net"
    "os"
    "strconv"
    "strings"
    "sync"
    "time"
)

// 第三方库（如有需要）
// 无特殊依赖，使用标准库 net 包即可
```

## 13. 实现优先级

| 优先级 | 功能 | 说明 |
|--------|------|------|
| P0 | 基础连接测试 | 核心功能，必须实现 |
| P0 | 端口扫描 | 核心功能，必须实现 |
| P1 | JSON 输出 | 便于脚本处理 |
| P1 | Banner 获取 | 常用功能 |
| P2 | 数据发送 | 高级功能 |
| P2 | 并发扫描优化 | 性能优化 |

---

**文档版本**: v1.0  
**创建日期**: 2026-04-24  
**作者**: AI Assistant
