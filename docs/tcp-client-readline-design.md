# TCP 客户端交互模式增强设计方案

> **目标**: 使用 readline 库实现交互式发送模式，支持历史记录和内置命令

---

## 一、功能需求

### 1.1 核心功能

| 功能 | 说明 |
|------|------|
| **历史记录** | 支持上下箭头浏览历史输入 |
| **行编辑** | 支持左右移动光标、删除、退格等 |
| **内置命令** | 以 `/` 开头的特殊命令 |
| **自动补全** | Tab 键补全内置命令 |
| **历史记录** | 内存中保存，退出后清空 |

### 1.2 内置命令设计

| 命令 | 功能 | 示例 |
|------|------|------|
| `/quit` 或 `/q` | 退出交互模式 | `/quit` |
| `/history` 或 `/h` | 显示历史记录 | `/history` |
| `/clear` 或 `/c` | 清屏 | `/clear` |
| `/info` 或 `/i` | 显示连接信息和统计 | `/info` |
| `/ping` | 发送测试消息并计算延迟 | `/ping` |
| `/help` 或 `/?` | 显示帮助信息 | `/help` |
| `/send <file>` | 发送文件内容 | `/send data.txt` |
| `/hex <data>` | 发送十六进制数据 | `/hex 48656c6c6f` |
| `/interval <ms>` | 设置发送间隔(毫秒) | `/interval 1000` |
| `/repeat <n>` | 重复发送上一次内容 | `/repeat 5` |

---

## 二、技术方案

### 2.1 依赖库选择

**推荐库**: `github.com/chzyer/readline`

**理由**:
- 纯 Go 实现，无 CGO 依赖
- 支持 Windows/Linux/Mac
- 功能完善：历史、编辑、补全、提示
- 使用广泛，文档完善

**替代方案**:
- `github.com/c-bata/go-prompt` - 更现代化，但体积大
- `github.com/peterh/liner` - 轻量，但功能较少

### 2.2 目录结构

```
internal/commands/tcp/
├── client.go           # 现有客户端逻辑
├── client_interactive.go # 新增：交互式模式（readline 实现）
├── client_builtin.go   # 新增：内置命令处理
├── cmd_tcp.go          # 配置定义
├── server.go           # 服务端
└── scanner.go          # 端口扫描
```

### 2.3 关键数据结构

```go
// InteractiveSession 交互式会话状态
type InteractiveSession struct {
    conn       net.Conn           // TCP 连接
    config     ClientConfig       // 客户端配置
    stats      *TransferStats     // 传输统计
    rl         *readline.Instance // readline 实例
    history    []string           // 内存中的历史记录
    lastSend   string             // 上一次发送的内容
    interval   time.Duration      // 发送间隔
}

// BuiltinCommand 内置命令定义
type BuiltinCommand struct {
    Name        string                           // 命令名
    Aliases     []string                         // 别名
    Description string                           // 描述
    Handler     func(*InteractiveSession, []string) error // 处理函数
}
```

---

## 三、详细设计

### 3.1 初始化流程

```go
func runInteractiveMode(conn net.Conn, config ClientConfig) error {
    // 1. 创建 readline 配置（不设置 HistoryFile，历史仅内存存储）
    rlConfig := &readline.Config{
        Prompt:          "tcp> ",
        AutoComplete:    buildCompleter(),
        InterruptPrompt: "^C",
        EOFPrompt:       "exit",
    }
    
    // 2. 创建 readline 实例
    rl, err := readline.NewEx(rlConfig)
    if err != nil {
        return err
    }
    defer rl.Close()
    
    // 3. 创建会话
    session := &InteractiveSession{
        conn:     conn,
        config:   config,
        stats:    &TransferStats{},
        rl:       rl,
        interval: 0,
    }
    
    // 4. 显示欢迎信息
    printWelcome()
    
    // 5. 主循环
    for {
        line, err := rl.Readline()
        if err != nil { // io.EOF, readline.ErrInterrupt
            break
        }
        
        // 处理输入
        if err := session.handleInput(line); err != nil {
            fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        }
    }
    
    return nil
}
```

### 3.2 输入处理逻辑

```go
func (s *InteractiveSession) handleInput(line string) error {
    line = strings.TrimSpace(line)
    if line == "" {
        return nil
    }
    
    // 检查是否是内置命令（以 / 开头）
    if strings.HasPrefix(line, "/") {
        return s.handleBuiltinCommand(line)
    }
    
    // 普通消息发送
    return s.sendMessage(line)
}
```

### 3.3 内置命令实现

```go
// 内置命令注册表
var builtinCommands = []*BuiltinCommand{
    {
        Name:        "quit",
        Aliases:     []string{"q", "exit"},
        Description: "退出交互模式",
        Handler:     cmdQuit,
    },
    {
        Name:        "history",
        Aliases:     []string{"h"},
        Description: "显示历史记录",
        Handler:     cmdHistory,
    },
    {
        Name:        "clear",
        Aliases:     []string{"c"},
        Description: "清屏",
        Handler:     cmdClear,
    },
    {
        Name:        "info",
        Aliases:     []string{"i"},
        Description: "显示连接信息和统计",
        Handler:     cmdInfo,
    },
    {
        Name:        "ping",
        Aliases:     []string{},
        Description: "发送测试消息并计算延迟",
        Handler:     cmdPing,
    },
    {
        Name:        "help",
        Aliases:     []string{"?"},
        Description: "显示帮助信息",
        Handler:     cmdHelp,
    },
    {
        Name:        "send",
        Aliases:     []string{},
        Description: "发送文件内容",
        Handler:     cmdSendFile,
    },
    {
        Name:        "hex",
        Aliases:     []string{},
        Description: "发送十六进制数据",
        Handler:     cmdSendHex,
    },
    {
        Name:        "interval",
        Aliases:     []string{"i"},
        Description: "设置发送间隔(毫秒)",
        Handler:     cmdInterval,
    },
    {
        Name:        "repeat",
        Aliases:     []string{"r"},
        Description: "重复发送上一次内容",
        Handler:     cmdRepeat,
    },
}

// 命令处理示例：quit
func cmdQuit(s *InteractiveSession, args []string) error {
    return io.EOF // 触发退出
}

// 命令处理示例：info（显示连接信息和统计）
func cmdInfo(s *InteractiveSession, args []string) error {
    fmt.Printf("Local Address:  %s\n", s.conn.LocalAddr())
    fmt.Printf("Remote Address: %s\n", s.conn.RemoteAddr())
    fmt.Printf("Duration:       %v\n", time.Since(s.stats.StartTime))
    fmt.Printf("Bytes Sent:     %d\n", s.stats.BytesSent)
    fmt.Printf("Bytes Received: %d\n", s.stats.BytesReceived)
    return nil
}

// 命令处理示例：ping（发送测试消息并计算延迟）
func cmdPing(s *InteractiveSession, args []string) error {
    start := time.Now()
    msg := "PING"
    if err := s.sendMessage(msg); err != nil {
        return err
    }
    latency := time.Since(start)
    fmt.Printf("Latency: %v\n", latency)
    return nil
}
```

### 3.4 自动补全实现

```go
func buildCompleter() *readline.PrefixCompleter {
    return readline.NewPrefixCompleter(
        readline.PcItem("/quit",
            readline.PcItem("-h"),
        ),
        readline.PcItem("/history"),
        readline.PcItem("/clear"),
        readline.PcItem("/info"),
        readline.PcItem("/ping"),
        readline.PcItem("/help"),
        readline.PcItem("/send",
            readline.PcItemDynamic(listFiles),
        ),
        readline.PcItem("/hex"),
        readline.PcItem("/interval"),
        readline.PcItem("/repeat"),
    )
}

// 动态补全：列出当前目录文件
func listFiles(prefix string) []string {
    files, _ := filepath.Glob(prefix + "*")
    return files
}
```

### 3.5 历史记录处理

历史记录仅保存在内存中，不持久化到文件。退出后历史记录清空。

```go
// 在 InteractiveSession 中维护历史
history []string

// 添加历史记录
func (s *InteractiveSession) addHistory(line string) {
    // 避免重复添加相同的连续记录
    if len(s.history) > 0 && s.history[len(s.history)-1] == line {
        return
    }
    s.history = append(s.history, line)
    // 限制历史记录数量（如最多 1000 条）
    if len(s.history) > 1000 {
        s.history = s.history[1:]
    }
}

// 显示历史记录（/history 命令）
func cmdHistory(s *InteractiveSession, args []string) error {
    for i, h := range s.history {
        fmt.Printf("%3d: %s\n", i+1, h)
    }
    return nil
}
```

---

## 四、改动范围

### 4.1 新增文件

| 文件 | 功能 | 预估行数 |
|------|------|----------|
| `client_interactive.go` | readline 交互模式主逻辑 | ~200 |
| `client_builtin.go` | 内置命令实现 | ~150 |

### 4.2 修改文件

| 文件 | 改动 | 预估行数 |
|------|------|----------|
| `client.go` | 移除旧的 sendInteractive，调用新的交互模式 | ~20 |
| `go.mod` | 添加 readline 依赖 | ~1 |

### 4.3 依赖添加

```bash
go get github.com/chzyer/readline
```

---

## 五、使用示例

```bash
# 启动交互模式
$ fck tcp client 127.0.0.1:8888
TCP Interactive Client
Connected to 127.0.0.1:8888
Type /help for available commands

tcp> hello
└ [2026-05-01 12:00:00] [OK] [ACK] Received 5 bytes

tcp> world
└ [2026-05-01 12:00:01] [OK] [ACK] Received 5 bytes

tcp> /info
Local Address:  192.168.1.5:12345
Remote Address: 127.0.0.1:8888
Duration:       1m30s
Bytes Sent:     10
Bytes Received: 106

tcp> /ping
PING
└ [2026-05-01 12:00:30] [OK] [ACK] Received 4 bytes
Latency: 2.5ms

tcp> /history
1: hello
2: world
3: /info
4: /ping

tcp> /send data.txt
Sent file data.txt (1024 bytes)
└ [2026-05-01 12:00:30] [OK] [ACK] Received 1024 bytes

tcp> /quit
Connection closed.
```

---

## 六、边缘案例处理

| 场景 | 处理方案 |
|------|----------|
| 内置命令不存在 | 提示错误，建议 /help |
| 发送文件不存在 | 提示错误，列出可用文件 |
| 十六进制格式错误 | 提示正确格式示例 |
| 连接断开 | 提示并退出交互模式 |
| 重复发送无历史 | 提示先发送一条消息 |
| 历史记录过多 | 自动清理旧记录（保留最近 1000 条） |

---

## 七、备选方案

### 方案B：不引入 readline，自己实现简单历史

如果担心依赖问题，可以用 `bufio.Scanner` + 自己维护历史切片，但功能有限。

### 方案C：使用更现代的 go-prompt

界面更美观，支持语法高亮，但依赖更多，体积更大。

---

**推荐采用方案A（chzyer/readline）**，功能完善且稳定。

**方案确认后实施代码修改**
