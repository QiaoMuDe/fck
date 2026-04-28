# TCP 命令重构设计方案

> **文档版本**: v1.0  
> **设计日期**: 2026-04-28  
> **遵循规范**: qflag 命令行工具开发规范 / qflag 命令开发规范

---

## 一、需求概述

重构 tcp 子命令，支持三种运行模式：

1. **扫描模式 (scan)**: 端口扫描功能，测试目标主机上哪些 TCP 端口处于开放状态
2. **客户端模式 (client)**: TCP 客户端功能，支持发送字符串、文件、目录及交互式模式
3. **服务端模式 (server)**: TCP 服务端功能，监听指定端口并接收客户端数据，必须返回响应包

---

## 二、命令行参数设计

### 2.1 命令结构

```
fck tcp [子命令] [选项] [参数]

子命令:
  scan    端口扫描模式
  client  TCP 客户端模式
  server  TCP 服务端模式
```

### 2.2 扫描模式 (tcp scan)

```bash
fck tcp scan [选项] <目标主机>
```

| 标志 | 短标志 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| --ports | -p | string | "1-1024" | 端口范围，支持格式: "80", "1-1024", "80,443,8080" |
| --timeout | -t | duration | 2s | 单个端口扫描超时时间 |
| --concurrent | -c | int | CPU核心数*2 | 并发扫描数量 |
| --show-closed | -s | bool | false | 显示关闭的端口 |
| --output | -o | string | "" | 输出结果到文件 |
| --format | -f | enum | "table" | 输出格式: table, json, csv |
| --verbose | -v | bool | false | 详细输出模式 |

**使用示例**:
```bash
# 扫描常用端口
fck tcp scan 192.168.1.1

# 扫描指定端口范围
fck tcp scan -p 1-65535 192.168.1.1

# 扫描多个指定端口
fck tcp scan -p 22,80,443,3306,8080 example.com

# 高并发快速扫描
fck tcp scan -c 500 -t 1s 192.168.1.1

# 导出 JSON 结果
fck tcp scan -f json -o result.json 192.168.1.1
```

### 2.3 客户端模式 (tcp client)

```bash
fck tcp client [选项] <服务器地址:端口>
```

| 标志 | 短标志 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| --message | -m | string | "" | 要发送的字符串消息 |
| --path | -p | string | "" | 要发送的文件/目录路径，支持通配符（如 *.txt） |
| --interactive | -i | bool | false | 交互式模式，持续输入发送 |
| --timeout | -t | duration | 10s | 连接和传输超时时间 |
| --buffer-size | -b | size | 4KB | 发送缓冲区大小 |
| --no-response | -n | bool | false | 不等待服务器响应 |
| --delimiter | -D | string | "\n" | 消息分隔符（交互式模式使用） |

**互斥组**:
- `message`, `path`, `interactive` 三者互斥，必须指定其中一个

**path 参数处理逻辑**:
- 如果是文件：直接发送该文件
- 如果是目录：发送该目录下所有文件（不递归子目录）
- 如果是通配符：匹配并发送所有匹配的文件

**使用示例**:
```bash
# 发送字符串
fck tcp client -m "Hello Server" 192.168.1.1:8080

# 发送单个文件
fck tcp client -p /path/to/file.txt 192.168.1.1:8080

# 发送目录下所有文件
fck tcp client -p /path/to/dir 192.168.1.1:8080

# 发送通配符匹配的文件
fck tcp client -p "/path/to/*.txt" 192.168.1.1:8080

# 交互式模式
fck tcp client -i 192.168.1.1:8080

# 使用自定义分隔符的交互式模式
fck tcp client -i -D "EOF" 192.168.1.1:8080
```

**响应处理说明**:
- 三种发送模式（字符串、路径、交互式）都会处理服务端返回的数据包
- 默认将响应内容输出到 stdout
- 使用 `-n/--no-response` 可禁用响应等待（仅发送，不接收）
- 使用 `-t/--timeout` 设置响应等待超时时间

### 2.4 服务端模式 (tcp server)

```bash
fck tcp server [选项]
```

| 标志 | 短标志 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| --port | -p | int | 8080 | 监听端口 |
| --address | -a | string | "0.0.0.0" | 监听地址 |
| --timeout | -t | duration | 30s | 连接超时时间 |
| --max-conn | -m | int | CPU核心数*2 | 最大并发连接数 |
| --buffer-size | -b | size | 4KB | 接收缓冲区大小 |
| --response | -r | string | "ACK" | 响应消息内容 |
| --output | -o | string | "" | 接收数据保存目录 |
| --echo | -e | bool | false | 回声模式（将接收数据原样返回） |
| --verbose | -v | bool | false | 详细输出模式 |

**使用示例**:
```bash
# 启动默认服务端
fck tcp server

# 指定端口和地址
fck tcp server -p 9090 -a 127.0.0.1

# 回声模式（用于测试）
fck tcp server -e -p 8080

# 保存接收的数据到文件
fck tcp server -o /tmp/received -p 8080

# 自定义响应消息
fck tcp server -r "Received OK" -p 8080
```

---

## 三、模块划分与代码组织结构

### 3.1 目录结构

```
internal/
├── cli/
│   ├── tcp.go              # 一级命令: fck tcp
│   └── tcp/                # 二级子命令目录
│       ├── scan.go         # fck tcp scan
│       ├── client.go       # fck tcp client
│       └── server.go       # fck tcp server
└── commands/
    └── tcp/                # tcp 业务逻辑目录
        ├── cmd_tcp.go      # 公共类型和工具函数
        ├── scanner.go      # 端口扫描实现
        ├── client.go       # TCP 客户端实现
        └── server.go       # TCP 服务端实现
```

### 3.2 文件职责说明

| 文件 | 职责 |
|------|------|
| `internal/cli/tcp.go` | 定义 tcp 一级命令，注册 scan/client/server 三个子命令 |
| `internal/cli/tcp/scan.go` | scan 子命令的 flag 定义和配置组装 |
| `internal/cli/tcp/client.go` | client 子命令的 flag 定义和配置组装 |
| `internal/cli/tcp/server.go` | server 子命令的 flag 定义和配置组装 |
| `internal/commands/tcp/cmd_tcp.go` | 公共类型定义（Config 结构体、常量、工具函数） |
| `internal/commands/tcp/scanner.go` | 端口扫描业务逻辑实现 |
| `internal/commands/tcp/client.go` | TCP 客户端业务逻辑实现 |
| `internal/commands/tcp/server.go` | TCP 服务端业务逻辑实现 |

---

## 四、核心功能实现逻辑

### 4.1 公共类型定义 (cmd_tcp.go)

```go
// Package tcp 实现 TCP 网络工具功能，包括端口扫描、客户端通信和服务端监听
package tcp

import (
    "time"
)

// 输出格式常量
const (
    FormatTable = "table"
    FormatJSON  = "json"
    FormatCSV   = "csv"
)

// 端口状态常量
const (
    PortOpen     = "open"
    PortClosed   = "closed"
    PortFiltered = "filtered"
    PortTimeout  = "timeout"
)

// ScanConfig 端口扫描配置
type ScanConfig struct {
    Target      string        // 目标主机
    Ports       string        // 端口范围字符串
    Timeout     time.Duration // 单个端口超时
    Concurrent  int           // 并发数
    ShowClosed  bool          // 显示关闭端口
    Output      string        // 输出文件
    Format      string        // 输出格式
    Verbose     bool          // 详细模式
}

// PortResult 端口扫描结果
type PortResult struct {
    Port     int       `json:"port"`
    Status   string    `json:"status"`
    Service  string    `json:"service,omitempty"`
    Response time.Duration `json:"response_time_ms"`
}

// ScanResult 扫描结果集合
type ScanResult struct {
    Target    string       `json:"target"`
    StartTime time.Time    `json:"start_time"`
    EndTime   time.Time    `json:"end_time"`
    Duration  time.Duration `json:"duration_ms"`
    Results   []PortResult `json:"results"`
    Stats     ScanStats    `json:"stats"`
}

// ScanStats 扫描统计
type ScanStats struct {
    Total    int `json:"total"`
    Open     int `json:"open"`
    Closed   int `json:"closed"`
    Filtered int `json:"filtered"`
    Timeout  int `json:"timeout"`
}

// ClientConfig TCP 客户端配置
type ClientConfig struct {
    Address      string        // 服务器地址
    Message      string        // 发送消息
    File         string        // 发送文件路径
    Dir          string        // 发送目录路径
    Interactive  bool          // 交互式模式
    Timeout      time.Duration // 超时时间
    BufferSize   int           // 缓冲区大小
    NoResponse   bool          // 不等待响应
    Delimiter    string        // 消息分隔符
}

// ServerConfig TCP 服务端配置
type ServerConfig struct {
    Address     string        // 监听地址
    Port        int           // 监听端口
    Timeout     time.Duration // 连接超时
    MaxConn     int           // 最大连接数
    BufferSize  int           // 缓冲区大小
    Response    string        // 响应消息
    OutputDir   string        // 输出目录
    Echo        bool          // 回声模式
    Verbose     bool          // 详细模式
}

// TransferStats 传输统计
type TransferStats struct {
    BytesSent     int64         `json:"bytes_sent"`
    BytesReceived int64         `json:"bytes_received"`
    Duration      time.Duration `json:"duration_ms"`
    FilesSent     int           `json:"files_sent,omitempty"`
}
```

### 4.2 端口扫描实现 (scanner.go)

```go
// ScanCmdMain 端口扫描主入口
//
// 参数:
//   - config: 扫描配置
//
// 返回值:
//   - error: 执行错误
func ScanCmdMain(config ScanConfig) error {
    // 1. 解析端口范围
    ports, err := parsePortRange(config.Ports)
    if err != nil {
        return fmt.Errorf("invalid port range: %w", err)
    }

    // 2. 解析目标地址
    target, err := resolveTarget(config.Target)
    if err != nil {
        return fmt.Errorf("failed to resolve target: %w", err)
    }

    // 3. 执行扫描
    result := &ScanResult{
        Target:    target,
        StartTime: time.Now(),
        Results:   make([]PortResult, 0, len(ports)),
    }

    // 使用工作池控制并发
    results := scanPorts(target, ports, config)
    result.Results = results
    result.EndTime = time.Now()
    result.Duration = result.EndTime.Sub(result.StartTime)

    // 4. 统计结果
    result.Stats = calculateStats(result.Results, config.ShowClosed)

    // 5. 输出结果
    return outputScanResult(result, config)
}

// scanPorts 并发扫描端口
//
// 参数:
//   - target: 目标地址
//   - ports: 端口列表
//   - config: 扫描配置
//
// 返回值:
//   - []PortResult: 扫描结果列表
func scanPorts(target string, ports []int, config ScanConfig) []PortResult {
    // 使用有缓冲通道作为工作队列
    portChan := make(chan int, config.Concurrent)
    resultChan := make(chan PortResult, len(ports))
    
    var wg sync.WaitGroup
    
    // 启动工作协程
    for i := 0; i < config.Concurrent; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for port := range portChan {
                result := scanSinglePort(target, port, config.Timeout)
                resultChan <- result
            }
        }()
    }
    
    // 分发任务
    go func() {
        for _, port := range ports {
            portChan <- port
        }
        close(portChan)
    }()
    
    // 等待所有扫描完成
    go func() {
        wg.Wait()
        close(resultChan)
    }()
    
    // 收集结果
    var results []PortResult
    for result := range resultChan {
        results = append(results, result)
    }
    
    return results
}

// scanSinglePort 扫描单个端口
//
// 参数:
//   - target: 目标地址
//   - port: 端口号
//   - timeout: 超时时间
//
// 返回值:
//   - PortResult: 扫描结果
func scanSinglePort(target string, port int, timeout time.Duration) PortResult {
    address := fmt.Sprintf("%s:%d", target, port)
    start := time.Now()
    
    conn, err := net.DialTimeout("tcp", address, timeout)
    duration := time.Since(start)
    
    result := PortResult{
        Port:     port,
        Response: duration,
    }
    
    if err != nil {
        if os.IsTimeout(err) || strings.Contains(err.Error(), "timeout") {
            result.Status = PortTimeout
        } else if strings.Contains(err.Error(), "refused") {
            result.Status = PortClosed
        } else {
            result.Status = PortFiltered
        }
        return result
    }
    
    defer conn.Close()
    result.Status = PortOpen
    result.Service = guessService(port)
    return result
}
```

### 4.3 TCP 客户端实现 (client.go)

```go
// ClientCmdMain TCP 客户端主入口
//
// 参数:
//   - config: 客户端配置
//
// 返回值:
//   - error: 执行错误
func ClientCmdMain(config ClientConfig) error {
    // 1. 建立连接
    conn, err := net.DialTimeout("tcp", config.Address, config.Timeout)
    if err != nil {
        return fmt.Errorf("failed to connect to %s: %w", config.Address, err)
    }
    defer conn.Close()

    // 2. 根据模式执行发送
    var stats *TransferStats
    switch {
    case config.Message != "":
        stats, err = sendString(conn, config)
    case config.File != "":
        stats, err = sendFile(conn, config)
    case config.Dir != "":
        stats, err = sendDir(conn, config)
    case config.Interactive:
        stats, err = sendInteractive(conn, config)
    }

    if err != nil {
        return err
    }

    // 3. 输出统计
    printClientStats(stats)
    return nil
}

// sendString 发送字符串消息
//
// 参数:
//   - conn: TCP 连接
//   - config: 客户端配置
//
// 返回值:
//   - *TransferStats: 传输统计
//   - error: 错误
func sendString(conn net.Conn, config ClientConfig) (*TransferStats, error) {
    stats := &TransferStats{StartTime: time.Now()}
    
    // 发送数据
    data := []byte(config.Message)
    n, err := conn.Write(data)
    if err != nil {
        return nil, fmt.Errorf("failed to send message: %w", err)
    }
    stats.BytesSent = int64(n)

    // 接收响应（如果不禁用）
    if !config.NoResponse {
        response, err := readResponse(conn, config.Timeout, config.BufferSize)
        if err != nil {
            return nil, fmt.Errorf("failed to read response: %w", err)
        }
        stats.BytesReceived = int64(len(response))
        fmt.Printf("Response: %s\n", string(response))
    }

    stats.Duration = time.Since(stats.StartTime)
    return stats, nil
}

// sendFile 发送单个文件
//
// 参数:
//   - conn: TCP 连接
//   - config: 客户端配置
//
// 返回值:
//   - *TransferStats: 传输统计
//   - error: 错误
func sendFile(conn net.Conn, config ClientConfig) (*TransferStats, error) {
    stats := &TransferStats{StartTime: time.Now(), FilesSent: 1}
    
    // 打开文件
    file, err := os.Open(config.File)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    // 发送文件头（文件名和大小）
    fileInfo, err := file.Stat()
    if err != nil {
        return nil, fmt.Errorf("failed to get file info: %w", err)
    }

    header := fmt.Sprintf("FILE:%s:%d\n", filepath.Base(config.File), fileInfo.Size())
    _, err = conn.Write([]byte(header))
    if err != nil {
        return nil, fmt.Errorf("failed to send file header: %w", err)
    }

    // 发送文件内容
    sent, err := io.Copy(conn, file)
    if err != nil {
        return nil, fmt.Errorf("failed to send file content: %w", err)
    }
    stats.BytesSent = sent + int64(len(header))

    // 接收响应
    if !config.NoResponse {
        response, err := readResponse(conn, config.Timeout, config.BufferSize)
        if err != nil {
            return nil, fmt.Errorf("failed to read response: %w", err)
        }
        stats.BytesReceived = int64(len(response))
    }

    stats.Duration = time.Since(stats.StartTime)
    return stats, nil
}

// sendDir 发送目录下所有文件（不递归）
//
// 参数:
//   - conn: TCP 连接
//   - config: 客户端配置
//
// 返回值:
//   - *TransferStats: 传输统计
//   - error: 错误
func sendDir(conn net.Conn, config ClientConfig) (*TransferStats, error) {
    stats := &TransferStats{StartTime: time.Now()}
    
    // 读取目录内容（不递归）
    entries, err := os.ReadDir(config.Dir)
    if err != nil {
        return nil, fmt.Errorf("failed to read directory: %w", err)
    }

    // 发送目录头
    dirHeader := fmt.Sprintf("DIR:%s:%d\n", filepath.Base(config.Dir), len(entries))
    _, err = conn.Write([]byte(dirHeader))
    if err != nil {
        return nil, fmt.Errorf("failed to send dir header: %w", err)
    }
    stats.BytesSent += int64(len(dirHeader))

    // 发送每个文件
    for _, entry := range entries {
        if entry.IsDir() {
            continue // 跳过子目录
        }

        filePath := filepath.Join(config.Dir, entry.Name())
        fileConfig := config
        fileConfig.File = filePath

        fileStats, err := sendFileInSession(conn, fileConfig)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Warning: failed to send %s: %v\n", entry.Name(), err)
            continue
        }

        stats.BytesSent += fileStats.BytesSent
        stats.BytesReceived += fileStats.BytesReceived
        stats.FilesSent++
    }

    stats.Duration = time.Since(stats.StartTime)
    return stats, nil
}

// sendInteractive 交互式发送模式
//
// 参数:
//   - conn: TCP 连接
//   - config: 客户端配置
//
// 返回值:
//   - *TransferStats: 传输统计
//   - error: 错误
func sendInteractive(conn net.Conn, config ClientConfig) (*TransferStats, error) {
    stats := &TransferStats{StartTime: time.Now()}
    reader := bufio.NewReader(os.Stdin)

    fmt.Println("Interactive mode started. Type your message and press Enter to send.")
    fmt.Printf("Use delimiter '%s' on a separate line to exit.\n", config.Delimiter)
    fmt.Println("----------------------------------------")

    for {
        fmt.Print("> ")
        line, err := reader.ReadString('\n')
        if err != nil {
            if err == io.EOF {
                break
            }
            return nil, fmt.Errorf("failed to read input: %w", err)
        }

        // 去除换行符
        line = strings.TrimRight(line, "\r\n")

        // 检查分隔符
        if line == config.Delimiter {
            fmt.Println("Exiting interactive mode...")
            break
        }

        // 发送消息
        data := []byte(line + "\n")
        n, err := conn.Write(data)
        if err != nil {
            return nil, fmt.Errorf("failed to send message: %w", err)
        }
        stats.BytesSent += int64(n)

        // 接收响应
        if !config.NoResponse {
            response, err := readResponse(conn, config.Timeout, config.BufferSize)
            if err != nil {
                fmt.Fprintf(os.Stderr, "Warning: failed to read response: %v\n", err)
            } else {
                stats.BytesReceived += int64(len(response))
                fmt.Printf("< %s\n", string(response))
            }
        }
    }

    stats.Duration = time.Since(stats.StartTime)
    return stats, nil
}
```

### 4.4 TCP 服务端实现 (server.go)

```go
// ServerCmdMain TCP 服务端主入口
//
// 参数:
//   - config: 服务端配置
//
// 返回值:
//   - error: 执行错误
func ServerCmdMain(config ServerConfig) error {
    // 1. 创建监听器
    address := fmt.Sprintf("%s:%d", config.Address, config.Port)
    listener, err := net.Listen("tcp", address)
    if err != nil {
        return fmt.Errorf("failed to start server on %s: %w", address, err)
    }
    defer listener.Close()

    fmt.Printf("TCP Server listening on %s\n", address)
    fmt.Printf("Press Ctrl+C to stop\n\n")

    // 2. 创建连接限制信号量
    sem := make(chan struct{}, config.MaxConn)

    // 3. 接受连接循环
    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Fprintf(os.Stderr, "Failed to accept connection: %v\n", err)
            continue
        }

        // 获取信号量（限制并发）
        sem <- struct{}{}

        // 处理连接
        go func(c net.Conn) {
            defer func() { <-sem }()
            defer c.Close()

            handleConnection(c, config)
        }(conn)
    }
}

// handleConnection 处理单个客户端连接
//
// 参数:
//   - conn: 客户端连接
//   - config: 服务端配置
func handleConnection(conn net.Conn, config ServerConfig) {
    clientAddr := conn.RemoteAddr().String()
    startTime := time.Now()

    if config.Verbose {
        fmt.Printf("[%s] Client connected\n", clientAddr)
    }

    // 设置超时
    if config.Timeout > 0 {
        conn.SetDeadline(time.Now().Add(config.Timeout))
    }

    // 读取数据
    buffer := make([]byte, config.BufferSize)
    var receivedData []byte
    var totalBytes int

    for {
        n, err := conn.Read(buffer)
        if n > 0 {
            receivedData = append(receivedData, buffer[:n]...)
            totalBytes += n

            // 检查是否是文件传输
            if isFileTransfer(receivedData) {
                handleFileTransfer(conn, receivedData, config)
                return
            }
        }

        if err != nil {
            if err != io.EOF {
                if config.Verbose {
                    fmt.Fprintf(os.Stderr, "[%s] Read error: %v\n", clientAddr, err)
                }
            }
            break
        }

        // 如果数据量较小，可能是单次消息，尝试处理
        if n < config.BufferSize {
            break
        }
    }

    // 处理接收到的数据
    if len(receivedData) > 0 {
        // 保存到文件（如果配置了输出目录）
        if config.OutputDir != "" {
            saveReceivedData(config.OutputDir, clientAddr, receivedData)
        }

        // 发送响应
        var response []byte
        if config.Echo {
            response = receivedData
        } else {
            response = []byte(config.Response)
        }

        _, err := conn.Write(response)
        if err != nil && config.Verbose {
            fmt.Fprintf(os.Stderr, "[%s] Failed to send response: %v\n", clientAddr, err)
        }

        duration := time.Since(startTime)
        if config.Verbose {
            fmt.Printf("[%s] Received %d bytes, sent %d bytes in %v\n",
                clientAddr, totalBytes, len(response), duration)
        }
    }

    if config.Verbose {
        fmt.Printf("[%s] Client disconnected\n", clientAddr)
    }
}

// isFileTransfer 检查是否是文件传输请求
//
// 参数:
//   - data: 接收到的数据
//
// 返回值:
//   - bool: 是否是文件传输
func isFileTransfer(data []byte) bool {
    return bytes.HasPrefix(data, []byte("FILE:")) || bytes.HasPrefix(data, []byte("DIR:"))
}

// handleFileTransfer 处理文件传输
//
// 参数:
//   - conn: 客户端连接
//   - data: 接收到的数据
//   - config: 服务端配置
func handleFileTransfer(conn net.Conn, data []byte, config ServerConfig) {
    // 解析文件头
    header, fileData, err := parseFileHeader(data)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to parse file header: %v\n", err)
        conn.Write([]byte("ERROR: Invalid file header"))
        return
    }

    // 如果还有更多数据需要读取
    if len(fileData) < header.Size {
        remaining := header.Size - len(fileData)
        extraData := make([]byte, remaining)
        _, err := io.ReadFull(conn, extraData)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Failed to read file content: %v\n", err)
            conn.Write([]byte("ERROR: Failed to receive file"))
            return
        }
        fileData = append(fileData, extraData...)
    }

    // 保存文件
    if config.OutputDir != "" {
        filePath := filepath.Join(config.OutputDir, header.Filename)
        err = os.WriteFile(filePath, fileData, 0644)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Failed to save file: %v\n", err)
            conn.Write([]byte("ERROR: Failed to save file"))
            return
        }
        fmt.Printf("Saved file: %s (%d bytes)\n", filePath, len(fileData))
    }

    // 发送确认响应
    response := fmt.Sprintf("RECEIVED:%s:%d", header.Filename, len(fileData))
    conn.Write([]byte(response))
}

// FileHeader 文件传输头部信息
type FileHeader struct {
    Filename string
    Size     int
}

// parseFileHeader 解析文件传输头部
//
// 参数:
//   - data: 接收到的数据
//
// 返回值:
//   - *FileHeader: 文件头信息
//   - []byte: 已接收的文件内容
//   - error: 错误
func parseFileHeader(data []byte) (*FileHeader, []byte, error) {
    // 格式: FILE:filename:size\ncontent
    parts := bytes.SplitN(data, []byte("\n"), 2)
    if len(parts) < 1 {
        return nil, nil, fmt.Errorf("invalid file header format")
    }

    headerStr := string(parts[0])
    var content []byte
    if len(parts) > 1 {
        content = parts[1]
    }

    // 解析 FILE:filename:size
    headerParts := strings.Split(headerStr, ":")
    if len(headerParts) != 3 {
        return nil, nil, fmt.Errorf("invalid file header format")
    }

    size, err := strconv.Atoi(headerParts[2])
    if err != nil {
        return nil, nil, fmt.Errorf("invalid file size: %w", err)
    }

    return &FileHeader{
        Filename: headerParts[1],
        Size:     size,
    }, content, nil
}
```

---

## 五、错误处理机制

### 5.1 错误分类

| 错误类型 | 说明 | 处理方式 |
|----------|------|----------|
| 参数错误 | 命令行参数无效 | 打印帮助信息并返回错误 |
| 网络错误 | 连接失败、超时等 | 包装错误并返回，支持重试 |
| IO 错误 | 文件读写失败 | 包装错误，继续处理其他文件 |
| 协议错误 | 数据传输格式错误 | 记录错误，发送错误响应 |
| 系统错误 | 资源不足、权限问题 | 直接返回错误 |

### 5.2 错误处理原则

1. **错误包装**: 使用 `fmt.Errorf("...: %w", err)` 包装底层错误
2. **友好提示**: 错误信息使用中文，便于用户理解
3. **分级处理**: 致命错误直接返回，非致命错误记录并继续
4. **超时控制**: 所有网络操作都设置超时

### 5.3 错误处理示例

```go
// 参数验证
if config.Target == "" {
    return fmt.Errorf("未指定目标主机")
}

// 网络错误处理
conn, err := net.DialTimeout("tcp", address, timeout)
if err != nil {
    if os.IsTimeout(err) {
        return fmt.Errorf("连接超时: %w", err)
    }
    return fmt.Errorf("连接失败: %w", err)
}

// IO 错误处理（非致命）
for _, file := range files {
    err := processFile(file)
    if err != nil {
        fmt.Fprintf(os.Stderr, "处理文件 %s 失败: %v\n", file, err)
        continue // 继续处理下一个
    }
}
```

---

## 六、日志记录策略

### 6.1 日志级别

| 级别 | 使用场景 | 输出目标 |
|------|----------|----------|
| ERROR | 错误信息 | stderr |
| WARN | 警告信息 | stderr |
| INFO | 一般信息 | stdout |
| DEBUG | 调试信息（verbose 模式） | stdout |

### 6.2 日志内容

```go
// 服务端连接日志
fmt.Printf("[%s] Client connected\n", clientAddr)
fmt.Printf("[%s] Received %d bytes\n", clientAddr, bytesReceived)
fmt.Printf("[%s] Client disconnected\n", clientAddr)

// 扫描进度日志（verbose 模式）
if config.Verbose {
    fmt.Printf("Scanning port %d/%d...\n", current, total)
}

// 传输统计
fmt.Printf("Sent: %d bytes, Received: %d bytes, Duration: %v\n",
    stats.BytesSent, stats.BytesReceived, stats.Duration)
```

### 6.3 日志格式

- 时间戳: 使用 `[HH:MM:SS]` 格式
- 客户端标识: `[IP:Port]`
- 结构化输出: JSON 格式用于文件导出

---

## 七、测试用例设计

### 7.1 单元测试

```go
// scanner_test.go

// TestParsePortRange 测试端口范围解析
func TestParsePortRange(t *testing.T) {
    tests := []struct {
        input    string
        expected []int
        wantErr  bool
    }{
        {"80", []int{80}, false},
        {"1-5", []int{1, 2, 3, 4, 5}, false},
        {"80,443,8080", []int{80, 443, 8080}, false},
        {"1-3,80,100-102", []int{1, 2, 3, 80, 100, 101, 102}, false},
        {"0", nil, true},      // 端口 0 无效
        {"65536", nil, true},  // 超出范围
        {"abc", nil, true},    // 无效格式
    }

    for _, tt := range tests {
        t.Run(tt.input, func(t *testing.T) {
            got, err := parsePortRange(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("parsePortRange() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.expected) {
                t.Errorf("parsePortRange() = %v, want %v", got, tt.expected)
            }
        })
    }
}

// TestScanSinglePort 测试单个端口扫描
func TestScanSinglePort(t *testing.T) {
    // 启动测试服务器
    listener, err := net.Listen("tcp", "127.0.0.1:0")
    if err != nil {
        t.Fatal(err)
    }
    defer listener.Close()

    port := listener.Addr().(*net.TCPAddr).Port

    // 在后台接受连接
    go func() {
        conn, _ := listener.Accept()
        if conn != nil {
            conn.Close()
        }
    }()

    // 测试开放端口
    result := scanSinglePort("127.0.0.1", port, time.Second)
    if result.Status != PortOpen {
        t.Errorf("Expected port %d to be open, got %s", port, result.Status)
    }

    // 测试关闭端口
    result = scanSinglePort("127.0.0.1", 1, time.Second)
    if result.Status != PortClosed && result.Status != PortFiltered {
        t.Errorf("Expected port 1 to be closed/filtered, got %s", result.Status)
    }
}
```

### 7.2 集成测试

```go
// tcp_integration_test.go

// TestClientServerCommunication 测试客户端服务端通信
func TestClientServerCommunication(t *testing.T) {
    // 启动服务端
    serverConfig := ServerConfig{
        Address:    "127.0.0.1",
        Port:       0, // 自动分配端口
        Response:   "ACK",
        BufferSize: 4096,
    }

    listener, err := net.Listen("tcp", "127.0.0.1:0")
    if err != nil {
        t.Fatal(err)
    }
    port := listener.Addr().(*net.TCPAddr).Port
    listener.Close()

    serverConfig.Port = port

    // 在后台启动服务端
    go func() {
        ServerCmdMain(serverConfig)
    }()

    time.Sleep(100 * time.Millisecond) // 等待服务端启动

    // 测试客户端发送字符串
    clientConfig := ClientConfig{
        Address:    fmt.Sprintf("127.0.0.1:%d", port),
        Message:    "Hello",
        Timeout:    time.Second * 5,
        BufferSize: 4096,
    }

    err = ClientCmdMain(clientConfig)
    if err != nil {
        t.Errorf("ClientCmdMain failed: %v", err)
    }
}

// TestFileTransfer 测试文件传输
func TestFileTransfer(t *testing.T) {
    // 创建临时文件
    tmpFile, err := os.CreateTemp("", "test-*.txt")
    if err != nil {
        t.Fatal(err)
    }
    defer os.Remove(tmpFile.Name())

    testData := []byte("Hello, World!")
    tmpFile.Write(testData)
    tmpFile.Close()

    // 创建临时目录接收文件
    tmpDir, err := os.MkdirTemp("", "tcp-test-*")
    if err != nil {
        t.Fatal(err)
    }
    defer os.RemoveAll(tmpDir)

    // 测试文件传输逻辑...
}
```

### 7.3 边界测试

| 测试场景 | 输入 | 预期结果 |
|----------|------|----------|
| 空消息 | `-m ""` | 发送空字符串 |
| 大文件 | 1GB 文件 | 成功传输，内存不溢出 |
| 超长端口范围 | `-p 1-65535` | 正常扫描 |
| 无效地址 | `256.256.256.256` | 返回解析错误 |
| 端口占用 | 启动两个服务端 | 第二个返回绑定错误 |
| 连接中断 | 服务端中途关闭 | 客户端返回连接错误 |

---

## 八、CLI 层实现

### 8.1 一级命令 (tcp.go)

```go
package cli

import (
    "fmt"

    "gitee.com/MM-Q/fck/internal/cli/tcp"
    "gitee.com/MM-Q/qflag"
)

var TcpCmd *qflag.Cmd

func init() {
    TcpCmd = qflag.NewCmd("tcp", "", qflag.ExitOnError)

    cmdOpts := &qflag.CmdOpts{
        Desc:        "TCP 网络工具，支持端口扫描、客户端通信和服务端监听",
        UsageSyntax: fmt.Sprintf("%s tcp [command] [options]", qflag.Root.Name()),
        UseChinese:  true,
        Notes: []string{
            "支持三种模式: scan(扫描), client(客户端), server(服务端)",
            "使用 'fck tcp [command] -h' 查看子命令详细帮助",
        },
        Examples: map[string]string{
            "端口扫描":     fmt.Sprintf("%s tcp scan 192.168.1.1 -p 1-1024", qflag.Root.Name()),
            "发送消息":     fmt.Sprintf("%s tcp client -m 'hello' 192.168.1.1:8080", qflag.Root.Name()),
            "启动服务端":   fmt.Sprintf("%s tcp server -p 8080", qflag.Root.Name()),
        },
        SubCmds: []qflag.Command{
            tcp.ScanCmd,
            tcp.ClientCmd,
            tcp.ServerCmd,
        },
    }

    if err := TcpCmd.ApplyOpts(cmdOpts); err != nil {
        panic(fmt.Errorf("apply opts err: %w", err))
    }

    TcpCmd.SetRun(runTcp)
}

func runTcp(cmd qflag.Command) error {
    cmd.PrintHelp()
    return nil
}
```

### 8.2 扫描子命令 (tcp/scan.go)

```go
package tcp

import (
    "fmt"

    "gitee.com/MM-Q/fck/internal/commands/tcp"
    "gitee.com/MM-Q/fck/internal/types"
    "gitee.com/MM-Q/qflag"
)

var ScanCmd *qflag.Cmd

var (
    scanPorts      *qflag.StringFlag
    scanTimeout    *qflag.DurationFlag
    scanConcurrent *qflag.IntFlag
    scanShowClosed *qflag.BoolFlag
    scanOutput     *qflag.StringFlag
    scanFormat     *qflag.EnumFlag
    scanVerbose    *qflag.BoolFlag
)

var scanFormatOptions = []string{"table", "json", "csv"}

func init() {
    ScanCmd = qflag.NewCmd("scan", "sc", qflag.ExitOnError)

    scanPorts = ScanCmd.String("ports", "p", "端口范围，支持格式: 80, 1-1024, 80,443,8080", "1-1024")
    scanTimeout = ScanCmd.Duration("timeout", "t", "单个端口扫描超时时间", types.DefaultTimeout)
    scanConcurrent = ScanCmd.Int("concurrent", "c", "并发扫描数量", 100)
    scanShowClosed = ScanCmd.Bool("show-closed", "s", "显示关闭的端口", false)
    scanOutput = ScanCmd.String("output", "o", "输出结果到文件", "")
    scanFormat = ScanCmd.Enum("format", "f", "输出格式", "table", scanFormatOptions)
    scanVerbose = ScanCmd.Bool("verbose", "v", "详细输出模式", false)

    cmdOpts := &qflag.CmdOpts{
        Desc:        "TCP 端口扫描工具",
        UsageSyntax: fmt.Sprintf("%s tcp scan [options] <target>", qflag.Root.Name()),
        UseChinese:  true,
        Notes: []string{
            "默认扫描 1-1024 端口",
            "支持域名和 IP 地址",
            "高并发扫描可能会被防火墙拦截",
        },
        Examples: map[string]string{
            "扫描常用端口":       fmt.Sprintf("%s tcp scan 192.168.1.1", qflag.Root.Name()),
            "扫描指定范围":       fmt.Sprintf("%s tcp scan -p 1-65535 192.168.1.1", qflag.Root.Name()),
            "扫描多个端口":       fmt.Sprintf("%s tcp scan -p 22,80,443 example.com", qflag.Root.Name()),
            "高并发快速扫描":     fmt.Sprintf("%s tcp scan -c 500 -t 1s 192.168.1.1", qflag.Root.Name()),
            "导出 JSON 结果":    fmt.Sprintf("%s tcp scan -f json -o result.json 192.168.1.1", qflag.Root.Name()),
        },
    }

    if err := ScanCmd.ApplyOpts(cmdOpts); err != nil {
        panic(fmt.Errorf("apply opts err: %w", err))
    }

    ScanCmd.SetRun(runScan)
}

func runScan(cmd qflag.Command) error {
    args := cmd.Args()
    if len(args) == 0 {
        return fmt.Errorf("未指定目标主机")
    }

    config := tcp.ScanConfig{
        Target:     args[0],
        Ports:      scanPorts.Get(),
        Timeout:    scanTimeout.Get(),
        Concurrent: scanConcurrent.Get(),
        ShowClosed: scanShowClosed.Get(),
        Output:     scanOutput.Get(),
        Format:     scanFormat.Get(),
        Verbose:    scanVerbose.Get(),
    }

    return tcp.ScanCmdMain(config)
}
```

### 8.3 客户端子命令 (tcp/client.go)

```go
package tcp

import (
    "fmt"

    "gitee.com/MM-Q/fck/internal/commands/tcp"
    "gitee.com/MM-Q/qflag"
)

var ClientCmd *qflag.Cmd

var (
    clientMessage     *qflag.StringFlag
    clientFile        *qflag.StringFlag
    clientDir         *qflag.StringFlag
    clientInteractive *qflag.BoolFlag
    clientTimeout     *qflag.DurationFlag
    clientBufferSize  *qflag.SizeFlag
    clientNoResponse  *qflag.BoolFlag
    clientDelimiter   *qflag.StringFlag
)

func init() {
    ClientCmd = qflag.NewCmd("client", "c", qflag.ExitOnError)

    clientMessage = ClientCmd.String("message", "m", "要发送的字符串消息", "")
    clientFile = ClientCmd.String("file", "f", "要发送的文件路径", "")
    clientDir = ClientCmd.String("dir", "d", "要发送的目录路径", "")
    clientInteractive = ClientCmd.Bool("interactive", "i", "交互式模式", false)
    clientTimeout = ClientCmd.Duration("timeout", "t", "连接和传输超时时间", tcp.DefaultTimeout)
    clientBufferSize = ClientCmd.Size("buffer-size", "b", "发送缓冲区大小", tcp.DefaultBufferSize)
    clientNoResponse = ClientCmd.Bool("no-response", "n", "不等待服务器响应", false)
    clientDelimiter = ClientCmd.String("delimiter", "D", "消息分隔符", "EOF")

    cmdOpts := &qflag.CmdOpts{
        Desc:        "TCP 客户端工具",
        UsageSyntax: fmt.Sprintf("%s tcp client [options] <address:port>", qflag.Root.Name()),
        UseChinese:  true,
        Notes: []string{
            "-m, -f, -d, -i 四个选项互斥，必须指定其中一个",
            "交互式模式下使用分隔符退出（默认 EOF）",
            "目录发送不会递归子目录",
        },
        Examples: map[string]string{
            "发送字符串":     fmt.Sprintf("%s tcp client -m 'Hello Server' 192.168.1.1:8080", qflag.Root.Name()),
            "发送文件":       fmt.Sprintf("%s tcp client -f /path/to/file.txt 192.168.1.1:8080", qflag.Root.Name()),
            "发送目录":       fmt.Sprintf("%s tcp client -d /path/to/dir 192.168.1.1:8080", qflag.Root.Name()),
            "交互式模式":     fmt.Sprintf("%s tcp client -i 192.168.1.1:8080", qflag.Root.Name()),
        },
        MutexGroups: []qflag.MutexGroup{
            {
                Name:      "send-mode",
                Flags:     []string{"message", "file", "dir", "interactive"},
                AllowNone: false,
            },
        },
    }

    if err := ClientCmd.ApplyOpts(cmdOpts); err != nil {
        panic(fmt.Errorf("apply opts err: %w", err))
    }

    ClientCmd.SetRun(runClient)
}

func runClient(cmd qflag.Command) error {
    args := cmd.Args()
    if len(args) == 0 {
        return fmt.Errorf("未指定服务器地址")
    }

    config := tcp.ClientConfig{
        Address:     args[0],
        Message:     clientMessage.Get(),
        File:        clientFile.Get(),
        Dir:         clientDir.Get(),
        Interactive: clientInteractive.Get(),
        Timeout:     clientTimeout.Get(),
        BufferSize:  int(clientBufferSize.Get()),
        NoResponse:  clientNoResponse.Get(),
        Delimiter:   clientDelimiter.Get(),
    }

    return tcp.ClientCmdMain(config)
}
```

### 8.4 服务端子命令 (tcp/server.go)

```go
package tcp

import (
    "fmt"

    "gitee.com/MM-Q/fck/internal/commands/tcp"
    "gitee.com/MM-Q/qflag"
)

var ServerCmd *qflag.Cmd

var (
    serverPort       *qflag.IntFlag
    serverAddress    *qflag.StringFlag
    serverTimeout    *qflag.DurationFlag
    serverMaxConn    *qflag.IntFlag
    serverBufferSize *qflag.SizeFlag
    serverResponse   *qflag.StringFlag
    serverOutput     *qflag.StringFlag
    serverEcho       *qflag.BoolFlag
    serverVerbose    *qflag.BoolFlag
)

func init() {
    ServerCmd = qflag.NewCmd("server", "s", qflag.ExitOnError)

    serverPort = ServerCmd.Int("port", "p", "监听端口", 8080)
    serverAddress = ServerCmd.String("address", "a", "监听地址", "0.0.0.0")
    serverTimeout = ServerCmd.Duration("timeout", "t", "连接超时时间", tcp.DefaultServerTimeout)
    serverMaxConn = ServerCmd.Int("max-conn", "m", "最大并发连接数", 100)
    serverBufferSize = ServerCmd.Size("buffer-size", "b", "接收缓冲区大小", tcp.DefaultBufferSize)
    serverResponse = ServerCmd.String("response", "r", "响应消息内容", "ACK")
    serverOutput = ServerCmd.String("output", "o", "接收数据保存目录", "")
    serverEcho = ServerCmd.Bool("echo", "e", "回声模式", false)
    serverVerbose = ServerCmd.Bool("verbose", "v", "详细输出模式", false)

    cmdOpts := &qflag.CmdOpts{
        Desc:        "TCP 服务端工具",
        UsageSyntax: fmt.Sprintf("%s tcp server [options]", qflag.Root.Name()),
        UseChinese:  true,
        Notes: []string{
            "默认监听 0.0.0.0:8080",
            "回声模式会将接收的数据原样返回",
            "使用 Ctrl+C 停止服务端",
        },
        Examples: map[string]string{
            "启动默认服务端":   fmt.Sprintf("%s tcp server", qflag.Root.Name()),
            "指定端口":        fmt.Sprintf("%s tcp server -p 9090", qflag.Root.Name()),
            "回声模式":        fmt.Sprintf("%s tcp server -e -p 8080", qflag.Root.Name()),
            "保存接收数据":    fmt.Sprintf("%s tcp server -o /tmp/received -p 8080", qflag.Root.Name()),
        },
    }

    if err := ServerCmd.ApplyOpts(cmdOpts); err != nil {
        panic(fmt.Errorf("apply opts err: %w", err))
    }

    ServerCmd.SetRun(runServer)
}

func runServer(cmd qflag.Command) error {
    config := tcp.ServerConfig{
        Address:    serverAddress.Get(),
        Port:       serverPort.Get(),
        Timeout:    serverTimeout.Get(),
        MaxConn:    serverMaxConn.Get(),
        BufferSize: int(serverBufferSize.Get()),
        Response:   serverResponse.Get(),
        OutputDir:  serverOutput.Get(),
        Echo:       serverEcho.Get(),
        Verbose:    serverVerbose.Get(),
    }

    return tcp.ServerCmdMain(config)
}
```

---

## 九、注册到根命令

在 `internal/cli/root.go` 的 `SubCmds` 列表中添加 `TcpCmd`:

```go
SubCmds: []qflag.Command{
    // ... 现有命令
    DnsCmd,
    TcpCmd,  // 添加 TCP 命令
    // ... 其他命令
},
```

---

## 十、实现检查清单

### 10.1 文件创建清单

- [ ] `internal/cli/tcp.go` - 一级命令定义
- [ ] `internal/cli/tcp/scan.go` - 扫描子命令
- [ ] `internal/cli/tcp/client.go` - 客户端子命令
- [ ] `internal/cli/tcp/server.go` - 服务端子命令
- [ ] `internal/commands/tcp/cmd_tcp.go` - 公共类型和常量
- [ ] `internal/commands/tcp/scanner.go` - 扫描业务逻辑
- [ ] `internal/commands/tcp/client.go` - 客户端业务逻辑
- [ ] `internal/commands/tcp/server.go` - 服务端业务逻辑

### 10.2 功能实现清单

- [ ] 端口范围解析（支持多种格式）
- [ ] 并发端口扫描
- [ ] 扫描结果格式化输出（table/json/csv）
- [ ] TCP 客户端连接管理
- [ ] 字符串发送功能
- [ ] 文件发送功能
- [ ] 目录发送功能（非递归）
- [ ] 交互式发送模式
- [ ] TCP 服务端监听
- [ ] 并发连接处理
- [ ] 文件接收功能
- [ ] 响应消息发送
- [ ] 回声模式

### 10.3 测试清单

- [ ] 端口范围解析单元测试
- [ ] 端口扫描单元测试
- [ ] 客户端-服务端集成测试
- [ ] 文件传输测试
- [ ] 边界条件测试

### 10.4 文档清单

- [ ] 命令帮助文档
- [ ] 使用示例
- [ ] 错误处理说明

---

## 十一、依赖说明

### 11.1 标准库依赖

```go
import (
    "bufio"
    "bytes"
    "context"
    "encoding/csv"
    "encoding/json"
    "fmt"
    "io"
    "net"
    "os"
    "path/filepath"
    "reflect"
    "sort"
    "strconv"
    "strings"
    "sync"
    "testing"
    "time"
)
```

### 11.2 项目内部依赖

```go
import (
    "gitee.com/MM-Q/fck/internal/types"
    "gitee.com/MM-Q/qflag"
)
```

### 11.3 第三方依赖

```go
import (
    "github.com/jedib0t/go-pretty/v6/table"
)
```

---

## 十二、风险评估与注意事项

### 12.1 安全风险

| 风险 | 说明 | 缓解措施 |
|------|------|----------|
| 端口扫描被检测 | 高频率扫描可能触发防火墙 | 提供可配置的并发数和超时时间 |
| 文件传输安全 | 接收任意文件可能存在风险 | 服务端限制接收目录，验证文件名 |
| 资源耗尽 | 大量并发连接消耗资源 | 限制最大并发连接数 |

### 12.2 性能考虑

- 端口扫描使用连接池控制并发
- 文件传输使用流式读写，避免大文件占用过多内存
- 服务端使用信号量限制并发连接数

### 12.3 兼容性

- 支持 Windows、Linux、macOS 跨平台运行
- 处理不同系统的换行符差异
- 支持 IPv4 和 IPv6 地址

---

**文档结束**
