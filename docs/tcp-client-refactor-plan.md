# TCP 客户端重构方案

> **日期**: 2026-05-01  
> **目标**: 简化 TCP 客户端发送模式，保留交互式、字符串、管道三种模式，移除文件路径模式

---

## 一、当前实现分析

### 现有发送模式

| 标志 | 模式 | 状态 |
|------|------|------|
| `-m, --message` | 发送字符串消息 | 保留 |
| `-p, --path` | 发送文件/目录/通配符 | **移除** |
| `-i, --interactive` | 交互式模式 | 保留 |

### 当前代码结构

```
internal/commands/tcp/
├── client.go      # 包含 sendString, sendPath, sendInteractive 三个发送函数
├── cmd_tcp.go     # 包含 ClientConfig 结构体定义

internal/cli/tcp/
├── client.go      # 包含 flag 定义和命令配置
```

---

## 二、重构目标

### 新模式设计

| 优先级 | 模式 | 触发条件 | 说明 |
|--------|------|----------|------|
| 1 | 管道模式 | `IsStdinPipe() == true` | 自动检测，无需标志 |
| 2 | 字符串模式 | `-m "消息内容"` | 显式指定消息 |
| 3 | 交互模式 | `-i` 或无标志 | 默认进入交互式输入 |

### 核心变更点

1. **移除 `-p, --path` 标志** 及所有相关逻辑
2. **新增管道自动检测**：利用 `utils.IsStdinPipe()` 优先处理管道输入
3. **简化模式判断逻辑**：管道 > 字符串 > 交互式
4. **更新文档和示例**：反映新的使用方式

---

## 三、详细改动方案

### 3.1 internal/commands/tcp/cmd_tcp.go

#### 修改内容

**移除字段：**
```go
// 移除 Path 字段
Path string
```

**ClientConfig 新结构：**
```go
type ClientConfig struct {
    Address     string
    Message     string        // -m 指定消息
    Interactive bool          // -i 交互模式
    Timeout     time.Duration
    BufferSize  int
    NoResponse  bool
    Delimiter   string        // 交互模式分隔符
}
```

---

### 3.2 internal/commands/tcp/client.go

#### 3.2.1 ClientCmdMain 函数重构

**新逻辑流程：**

```go
func ClientCmdMain(config ClientConfig) error {
    // 1. 建立连接
    conn, err := net.DialTimeout("tcp", config.Address, config.Timeout)
    ...

    // 2. 优先级判断发送模式
    var stats *TransferStats
    switch {
    case utils.IsStdinPipe():           // 优先级1: 管道输入
        stats, err = sendStdin(conn, config)
    case config.Message != "":          // 优先级2: 指定消息
        stats, err = sendString(conn, config)
    default:                            // 优先级3: 交互模式（默认）
        stats, err = sendInteractive(conn, config)
    }

    // 3. 输出统计
    printClientStats(stats)
    return nil
}
```

#### 3.2.2 新增 sendStdin 函数

```go
// sendStdin 从标准输入（管道）读取并发送数据
//
// 参数:
//   - conn: TCP 连接
//   - config: 客户端配置
//
// 返回值:
//   - *TransferStats: 传输统计
//   - error: 错误
func sendStdin(conn net.Conn, config ClientConfig) (*TransferStats, error) {
    startTime := time.Now()
    stats := &TransferStats{}

    // 读取全部 stdin 内容
    data, err := io.ReadAll(os.Stdin)
    if err != nil {
        return nil, fmt.Errorf("failed to read from stdin: %w", err)
    }

    // 发送数据
    n, err := conn.Write(data)
    if err != nil {
        return nil, fmt.Errorf("failed to send data: %w", err)
    }
    stats.BytesSent = int64(n)

    // 关闭写入端
    if tcpConn, ok := conn.(*net.TCPConn); ok {
        if err := tcpConn.CloseWrite(); err != nil {
            return nil, fmt.Errorf("failed to close write: %w", err)
        }
    }

    // 接收响应
    if !config.NoResponse {
        response, err := readResponse(conn, config.Timeout, config.BufferSize)
        if err != nil {
            return nil, fmt.Errorf("failed to read response: %w", err)
        }
        stats.BytesReceived = int64(len(response))
        fmt.Printf("Response: %s\n", string(response))
    }

    stats.Duration = time.Since(startTime)
    return stats, nil
}
```

#### 3.2.3 移除的函数

需要完全删除以下函数：

- `sendPath()` - 文件路径发送主函数
- `resolvePath()` - 路径解析函数（支持通配符）
- `sendRawFileData()` - 原始文件数据发送函数

#### 3.2.4 保留的函数

以下函数保持不变：

- `sendString()` - 字符串发送（逻辑不变）
- `sendInteractive()` - 交互式发送（逻辑不变）
- `readResponse()` - 读取响应（逻辑不变）
- `printClientStats()` - 打印统计（逻辑不变）

---

### 3.3 internal/cli/tcp/client.go

#### 3.3.1 移除 flag 定义

**删除：**
```go
clientPath = ClientCmd.String("path", "p", "要发送的文件/目录路径, 支持通配符 (如 *.txt) ", "")
```

#### 3.3.2 更新互斥组配置

**原配置：**
```go
MutexGroups: []qflag.MutexGroup{
    {
        Name:      "send-mode",
        Flags:     []string{"message", "path", "interactive"},
        AllowNone: true,  // 允许不指定，但代码里会检查报错
    },
},
```

**新配置：**
```go
// 移除互斥组，改为代码逻辑控制优先级
// 或者保留但只包含 message 和 interactive
MutexGroups: []qflag.MutexGroup{
    {
        Name:      "send-mode",
        Flags:     []string{"message", "interactive"},
        AllowNone: true,  // 允许都不指定，此时默认进入交互模式
    },
},
```

#### 3.3.3 更新 runClient 函数

**删除路径检查：**
```go
// 删除以下检查
if clientMessage.Get() == "" && clientPath.Get() == "" && !clientInteractive.Get() {
    return fmt.Errorf("must specify send mode: -m (message), -p (path) or -i (interactive)")
}
```

**新配置构建：**
```go
config := tcp.ClientConfig{
    Address:     args[0],
    Message:     clientMessage.Get(),
    // Path:     clientPath.Get(),  // 删除
    Interactive: clientInteractive.Get(),
    Timeout:     clientTimeout.Get(),
    BufferSize:  int(clientBufferSize.Get()),
    NoResponse:  clientNoResponse.Get(),
    Delimiter:   clientDelimiter.Get(),
}
```

#### 3.3.4 更新文档说明

**原 Notes：**
```go
Notes: []string{
    "-m, -p, -i 三个选项互斥, 必须指定其中一个",
    "交互式模式下使用分隔符退出 (默认 EOF) ",
    "路径参数支持文件、目录 (不递归子目录) 和通配符",
    "三种发送模式都会处理服务端返回的数据包",
},
```

**新 Notes：**
```go
Notes: []string{
    "支持三种发送模式（按优先级）:",
    "  1. 管道输入: echo 'data' | fck tcp client <address>",
    "  2. 字符串模式: -m '消息内容'（与管道互斥）",
    "  3. 交互模式: -i 或无标志（默认）",
    "交互式模式下使用分隔符退出 (默认 EOF)",
    "管道模式自动检测，无需额外标志",
},
```

#### 3.3.5 更新使用示例

**原 Examples：**
```go
Examples: map[string]string{
    "发送字符串":     fmt.Sprintf("%s tcp client -m 'Hello Server' 192.168.1.1:8080", qflag.Root.Name()),
    "发送单个文件":    fmt.Sprintf("%s tcp client -p /path/to/file.txt 192.168.1.1:8080", qflag.Root.Name()),
    "发送目录下所有文件": fmt.Sprintf("%s tcp client -p /path/to/dir 192.168.1.1:8080", qflag.Root.Name()),
    "发送通配符匹配文件": fmt.Sprintf("%s tcp client -p '/path/to/*.txt' 192.168.1.1:8080", qflag.Root.Name()),
    "交互式模式":     fmt.Sprintf("%s tcp client -i 192.168.1.1:8080", qflag.Root.Name()),
    "不等待响应":     fmt.Sprintf("%s tcp client -m 'hello' -n 192.168.1.1:8080", qflag.Root.Name()),
},
```

**新 Examples：**
```go
Examples: map[string]string{
    "管道发送":       fmt.Sprintf("echo 'Hello Server' | %s tcp client 192.168.1.1:8080", qflag.Root.Name()),
    "文件内容管道发送": fmt.Sprintf("cat data.txt | %s tcp client 192.168.1.1:8080", qflag.Root.Name()),
    "发送字符串":     fmt.Sprintf("%s tcp client -m 'Hello Server' 192.168.1.1:8080", qflag.Root.Name()),
    "交互式模式":     fmt.Sprintf("%s tcp client -i 192.168.1.1:8080", qflag.Root.Name()),
    "默认交互模式":   fmt.Sprintf("%s tcp client 192.168.1.1:8080", qflag.Root.Name()),
    "不等待响应":     fmt.Sprintf("%s tcp client -m 'hello' -n 192.168.1.1:8080", qflag.Root.Name()),
},
```

---

## 四、改动文件清单

| 序号 | 文件路径 | 改动类型 | 预估行数 |
|------|----------|----------|----------|
| 1 | `internal/commands/tcp/cmd_tcp.go` | 修改 | ~3 行（删除 Path 字段） |
| 2 | `internal/commands/tcp/client.go` | 修改 | ~100 行（删除 sendPath 相关，添加 sendStdin） |
| 3 | `internal/cli/tcp/client.go` | 修改 | ~40 行（删除 flag，更新文档） |

**总计**：约 143 行改动，涉及 3 个文件

---

## 五、使用方式对比

### 重构前

```bash
# 发送字符串
fck tcp client -m "hello" 192.168.1.1:8080

# 发送文件（移除）
fck tcp client -p /path/to/file.txt 192.168.1.1:8080

# 发送目录（移除）
fck tcp client -p /path/to/dir 192.168.1.1:8080

# 交互模式
fck tcp client -i 192.168.1.1:8080
```

### 重构后

```bash
# 管道发送（新增）
echo "hello" | fck tcp client 192.168.1.1:8080
cat file.txt | fck tcp client 192.168.1.1:8080

# 发送字符串
fck tcp client -m "hello" 192.168.1.1:8080

# 交互模式
fck tcp client -i 192.168.1.1:8080

# 默认交互模式（无标志）
fck tcp client 192.168.1.1:8080
```

---

## 六、边缘案例与测试建议

### 6.1 边缘案例

| 场景 | 预期行为 |
|------|----------|
| 管道 + `-m` 同时存在 | 管道优先，忽略 `-m` |
| 管道 + `-i` 同时存在 | 管道优先，忽略 `-i` |
| 空管道输入 | 发送空数据，正常关闭连接 |
| 大文件管道 | 使用 io.ReadAll，需注意内存占用 |
| 无标志无管道 | 进入默认交互模式 |

### 6.2 测试用例建议

```bash
# 测试1: 管道发送
echo "test message" | fck tcp client 127.0.0.1:8080

# 测试2: 文件管道
cat test.txt | fck tcp client 127.0.0.1:8080

# 测试3: 字符串模式
fck tcp client -m "test" 127.0.0.1:8080

# 测试4: 交互模式
fck tcp client -i 127.0.0.1:8080

# 测试5: 默认交互模式
fck tcp client 127.0.0.1:8080

# 测试6: 管道优先于 -m
echo "pipe" | fck tcp client -m "msg" 127.0.0.1:8080  # 应发送 "pipe"

# 测试7: 空管道
echo -n "" | fck tcp client 127.0.0.1:8080
```

---

## 七、实现注意事项

1. **导入包**：`sendStdin` 需要导入 `io` 和 `utils` 包
2. **错误处理**：保持与现有代码一致的错误处理风格
3. **统计信息**：`sendStdin` 应正确填充 `TransferStats`
4. **向后兼容**：此变更为破坏性变更，移除 `-p` 标志

---

**方案确认后实施代码修改**
