# TCP 交互式模式设计方案

## 概述
为 TCP 命令新增交互式发送数据模式，支持实时输入、发送和接收响应。

**技术选型**: 使用 `github.com/chzyer/readline` 库实现，提供 Tab 补全、历史记录、行编辑等完整交互体验。

## 依赖

```bash
go get github.com/chzyer/readline
```

## 新增标志

| 标志 | 简写 | 说明 | 默认值 |
|------|------|------|--------|
| `--interactive` | `-I` | 启用交互式模式（大写I避免与 `-i` 间隔冲突） | false |
| `--hex` | - | 以十六进制格式显示接收的数据 | false |

## 功能特性

### 基本功能
1. 连接成功后进入交互式 shell
2. 用户输入一行，回车立即发送
3. 发送后自动等待并显示响应（可配置等待时间 `-w`）
4. 支持特殊命令控制交互会话

### readline 提供的增强功能

| 功能 | 快捷键 | 说明 |
|------|--------|------|
| **Tab 补全** | `Tab` | 命令和参数自动补全 |
| **历史记录** | `↑` / `↓` | 浏览历史命令 |
| **行编辑** | `←` / `→` | 光标移动 |
| **行首/行尾** | `Ctrl+A` / `Ctrl+E` | 快速跳转 |
| **删除字符** | `Backspace` / `Delete` | 删除字符 |
| **清空行** | `Ctrl+U` / `Ctrl+K` | 删除光标前/后内容 |
| **搜索历史** | `Ctrl+R` | 反向搜索历史命令 |
| **中断** | `Ctrl+C` | 优雅退出 |

### 特殊命令
在交互模式下，输入以下命令执行特殊操作：

| 命令 | 说明 |
|------|------|
| `quit` / `exit` | 退出交互模式，关闭连接 |
| `close` | 关闭当前连接（可重新连接）|
| `hex` | 切换十六进制显示模式 |
| `file <path>` | 发送文件内容 |
| `help` | 显示帮助信息 |
| `status` | 显示连接状态信息 |

### 交互界面示例

```
$ fck tcp -I target.com 8080
Connecting to target.com (192.168.1.1):8080...
Connected to target.com:8080
Type 'help' for available commands

tcp> hello
[←] Received (5 bytes): world
tcp> test
[←] Received (4 bytes): ok
tcp> quit
Connection closed.
```

### Tab 补全示例

```
tcp> he<Tab>     # 补全为 help
tcp> fi<Tab>     # 补全为 file
tcp> qu<Tab>     # 补全为 quit
```

### 十六进制显示模式

```
tcp> hello
[←] Received (5 bytes):
00000000  77 6f 72 6c 64                                    |world|
tcp> hex         # 切换回文本模式
Hex mode: OFF
tcp>
```

## 实现要点

### 核心逻辑
1. 使用 `readline.Instance` 创建交互式读取器
2. 配置 Tab 补全器（PrefixCompleter）
3. 使用 goroutine 分离发送和接收
4. 历史记录自动保存到文件

### 并发模型
```
主 goroutine: readline.Readline() → 处理输入 → 发送数据
子 goroutine: 持续读取响应 → 打印到控制台
```

### Tab 补全配置
```go
var completer = readline.NewPrefixCompleter(
    readline.PcItem("help"),
    readline.PcItem("quit"),
    readline.PcItem("exit"),
    readline.PcItem("close"),
    readline.PcItem("hex"),
    readline.PcItem("status"),
    readline.PcItem("file",
        readline.PcItemDynamic(listFiles),
    ),
)
```

### 核心函数实现
```go
// runInteractiveMode 使用 readline 运行交互式模式
func runInteractiveMode(config TcpConfig, ipAddr net.IP) error {
    // 建立连接
    conn, err := net.DialTimeout("tcp", target, config.Timeout)
    if err != nil {
        return err
    }
    defer conn.Close()

    // 配置 readline
    rl, err := readline.NewEx(&readline.Config{
        Prompt:          "tcp> ",
        HistoryFile:     "/tmp/fck-tcp-history",
        AutoComplete:    completer,
        InterruptPrompt: "^C",
        EOFPrompt:       "exit",
    })
    if err != nil {
        return err
    }
    defer rl.Close()

    // 启动接收 goroutine
    go receiveLoop(conn, config)

    // 主循环：读取用户输入
    for {
        line, err := rl.Readline()
        if err != nil { // io.EOF, readline.ErrInterrupt
            break
        }

        // 处理命令或发送数据
        if err := handleInput(conn, line); err != nil {
            fmt.Printf("Error: %v\n", err)
        }
    }

    return nil
}
```

### 数据流
```
用户输入 → readline → [tcp>] 提示 → 发送 → 等待响应 → [←] 标记 → 显示响应
                                              ↑
                                         子 goroutine 监听
```

### 特殊处理
- 空行：发送换行符（可配置行为）
- 超时：提示超时但仍保持连接，可继续发送
- 连接断开：提示并自动尝试重连（可选）
- Ctrl+C：优雅退出，关闭连接
- 历史记录：自动保存到 `~/.fck/tcp_history`

## 使用示例

```bash
# 基础交互模式
fck tcp -I target.com 8080

# 交互模式 + 2秒响应等待
fck tcp -I -w 2s target.com 8080

# 交互模式 + 十六进制显示
fck tcp -I --hex target.com 8080

# 交互模式 + 静默（仅显示收发数据，无提示信息）
fck tcp -I -q target.com 8080
```

## 与现有标志的互斥关系

交互式模式应与以下标志互斥（通过 MutexGroup 实现）：
- `-d, --data`（单次发送数据）
- `-f, --file`（从文件发送）
- `-c, --count`（多次连接）
- `-l, --listen`（监听模式）

## 技术实现细节

### 新增结构体
```go
// InteractiveSession 交互式会话状态
type InteractiveSession struct {
    Conn       net.Conn
    Config     TcpConfig
    HexMode    bool
    StartTime  time.Time
    BytesSent  int64
    BytesRecv  int64
    Readline   *readline.Instance
}
```

### 核心函数
```go
// runInteractiveMode 运行交互式模式
func runInteractiveMode(config TcpConfig, ipAddr net.IP) error

// handleInteractiveInput 处理用户输入
func handleInteractiveInput(session *InteractiveSession, input string) error

// receiveLoop 接收数据循环（goroutine）
func receiveLoop(session *InteractiveSession, done chan<- bool)

// printReceivedData 打印接收到的数据
func printReceivedData(data []byte, hexMode bool)

// setupReadline 配置 readline 实例
func setupReadline(config TcpConfig) (*readline.Instance, error)
```

### 输入处理流程
1. readline 读取一行输入（支持编辑、历史、补全）
2. 检查是否为特殊命令
3. 如果是命令，执行对应操作
4. 如果是普通数据，发送到服务器
5. 等待响应（根据 `-w` 超时设置）
6. 显示响应数据

## 错误处理

| 场景 | 处理方式 |
|------|----------|
| 连接失败 | 显示错误，退出程序 |
| 发送失败 | 显示错误，保持连接 |
| 接收超时 | 提示超时，保持连接 |
| 连接断开 | 提示断开，询问是否重连 |
| 文件读取失败 | 显示错误，保持连接 |
| readline 初始化失败 | 回退到基础模式 |

## 扩展性考虑

1. **脚本模式**：支持从脚本文件读取命令序列
2. **自动重连**：断开后自动重连选项
3. **多连接**：同时管理多个连接（类似 tab）
4. **宏录制**：录制和回放命令序列
5. **插件系统**：支持自定义命令扩展

## 方案对比

| 功能 | 基础实现 | readline 实现 |
|------|----------|---------------|
| 历史记录 | ❌ | ✅ 自动保存到文件 |
| Tab 补全 | ❌ | ✅ 命令补全 |
| 行编辑 | ❌ | ✅ 光标移动、删除 |
| 快捷键 | ❌ | ✅ Ctrl+A/E/C 等 |
| 提示定制 | 简单 | 丰富（颜色、前缀）|
| 依赖 | 无 | 1 个外部库 |
| 实现复杂度 | 低 | 中等 |
