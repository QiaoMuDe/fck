# TCP 命令重构设计 - 对话历史记录

> **导出日期**: 2026-04-28  
> **项目**: FCK - 一站式文件与系统管理工具集  
> **主题**: TCP 子命令重构设计方案讨论  
> **状态**: 设计方案已确定，待实现

---

## 一、项目背景

### 1.1 项目概述
- **项目名称**: FCK
- **定位**: 一站式文件与系统管理工具集
- **技术栈**: Go 1.25.0
- **CLI 框架**: 自研 qflag 库
- **版本**: v1.8

### 1.2 项目结构
```
fck/
├── cmd/main.go              # 程序入口
├── internal/
│   ├── cli/                 # 命令 flag 定义
│   ├── commands/            # 业务逻辑
│   ├── types/               # 全局类型和常量
│   └── utils/               # 通用工具
├── docs/                    # 设计文档
└── vendor/                  # 自研依赖库
```

### 1.3 开发规范
项目遵循以下开发规范文档：
- `docs/qflag命令行工具开发规范.md` - 目录结构和命令定义规范
- `docs/qflag命令开发规范.md` - 业务逻辑和 CLI 定义规范

---

## 二、需求分析

### 2.1 原始需求
重构 tcp 子命令，支持三种运行模式：

1. **扫描模式 (scan)**: 端口扫描功能，测试目标主机上哪些 TCP 端口处于开放状态
2. **客户端模式 (client)**: TCP 客户端功能，支持：
   - 字符串发送
   - 文件发送
   - 目录发送（不递归子目录）
   - 交互式发送模式
3. **服务端模式 (server)**: TCP 服务端功能，监听指定端口并接收客户端数据，必须返回响应包

### 2.2 方案设计文档
完整设计方案保存在: `docs/tcp-command-refactor-design.md`

---

## 三、对话历史与技术决策

### 3.1 初始方案评审

**用户初始要求**:
- 基于 qflag 开发规范设计 TCP 命令
- 三种模式：扫描、客户端、服务端
- 需要包含：命令行参数设计、模块划分、核心逻辑、错误处理、日志策略、测试用例

**AI 完成的工作**:
1. 分析了现有项目结构（port、dns 等命令的实现方式）
2. 阅读了 qflag 开发规范文档
3. 设计了完整的 TCP 命令重构方案

---

### 3.2 技术决策讨论

#### 决策 1: 并发扫描数量默认值

**原始设计**:
```
--concurrent | -c | int | 100 | 并发扫描数量
```

**用户反馈**:
> 默认应该是 cpu 的核心*2

**修改结果**:
```
--concurrent | -c | int | CPU核心数*2 | 并发扫描数量
```

**实现提示**:
代码中使用 `runtime.NumCPU() * 2` 获取默认值

---

#### 决策 2: 客户端路径参数合并

**原始设计**:
```
--file | -f | string | "" | 要发送的文件路径
--dir  | -d | string | "" | 要发送的目录路径（不递归子目录）
```

**用户反馈**:
> 这个完全可以只有一个标志 --path/-p。然后如果是文件就直接处理，如果是目录就读取他一层内容里的所有文件，如果是通配符则匹配处理就行

**修改结果**:
```
--path | -p | string | "" | 要发送的文件/目录路径，支持通配符（如 *.txt）
```

**path 参数处理逻辑**:
- 如果是文件：直接发送该文件
- 如果是目录：发送该目录下所有文件（不递归子目录）
- 如果是通配符：匹配并发送所有匹配的文件

**互斥组调整**:
- 原始: `message`, `file`, `dir`, `interactive` 四者互斥
- 修改后: `message`, `path`, `interactive` 三者互斥

---

#### 决策 3: 服务端最大并发连接数默认值

**原始设计**:
```
--max-conn | -m | int | 100 | 最大并发连接数
```

**用户反馈**:
> 这个也默认 cpu 核心*2 吧

**修改结果**:
```
--max-conn | -m | int | CPU核心数*2 | 最大并发连接数
```

---

#### 决策 4: 客户端响应处理说明

**用户反馈**:
> 你这客户端模式，发送字符串，发送路径，交互式发送，这是三个模式都要能处理服务端返回的数据包，你加个这个说明

**添加的文档内容**:
```markdown
**响应处理说明**:
- 三种发送模式（字符串、路径、交互式）都会处理服务端返回的数据包
- 默认将响应内容输出到 stdout
- 使用 `-n/--no-response` 可禁用响应等待（仅发送，不接收）
- 使用 `-t/--timeout` 设置响应等待超时时间
```

---

## 四、最终设计方案概要

### 4.1 命令结构
```
fck tcp [子命令] [选项] [参数]

子命令:
  scan    端口扫描模式
  client  TCP 客户端模式
  server  TCP 服务端模式
```

### 4.2 扫描模式 (tcp scan)

| 标志 | 短标志 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| --ports | -p | string | "1-1024" | 端口范围 |
| --timeout | -t | duration | 2s | 单个端口扫描超时时间 |
| --concurrent | -c | int | CPU核心数*2 | 并发扫描数量 |
| --show-closed | -s | bool | false | 显示关闭的端口 |
| --output | -o | string | "" | 输出结果到文件 |
| --format | -f | enum | "table" | 输出格式: table, json, csv |
| --verbose | -v | bool | false | 详细输出模式 |

### 4.3 客户端模式 (tcp client)

| 标志 | 短标志 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| --message | -m | string | "" | 要发送的字符串消息 |
| --path | -p | string | "" | 要发送的文件/目录路径，支持通配符 |
| --interactive | -i | bool | false | 交互式模式 |
| --timeout | -t | duration | 10s | 连接和传输超时时间 |
| --buffer-size | -b | size | 4KB | 发送缓冲区大小 |
| --no-response | -n | bool | false | 不等待服务器响应 |
| --delimiter | -D | string | "\n" | 消息分隔符 |

**互斥组**: `message`, `path`, `interactive` 三者互斥

**path 处理逻辑**:
- 文件：直接发送
- 目录：发送目录下所有文件（不递归）
- 通配符：匹配并发送所有匹配文件

**响应处理**:
- 三种模式都会处理服务端返回的数据包
- 默认输出响应到 stdout
- 支持禁用响应等待和设置超时

### 4.4 服务端模式 (tcp server)

| 标志 | 短标志 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| --port | -p | int | 8080 | 监听端口 |
| --address | -a | string | "0.0.0.0" | 监听地址 |
| --timeout | -t | duration | 30s | 连接超时时间 |
| --max-conn | -m | int | CPU核心数*2 | 最大并发连接数 |
| --buffer-size | -b | size | 4KB | 接收缓冲区大小 |
| --response | -r | string | "ACK" | 响应消息内容 |
| --output | -o | string | "" | 接收数据保存目录 |
| --echo | -e | bool | false | 回声模式 |
| --verbose | -v | bool | false | 详细输出模式 |

---

## 五、代码组织结构

### 5.1 目录结构
```
internal/
├── cli/
│   ├── tcp.go              # 一级命令: fck tcp
│   └── tcp/
│       ├── scan.go         # fck tcp scan
│       ├── client.go       # fck tcp client
│       └── server.go       # fck tcp server
└── commands/
    └── tcp/
        ├── cmd_tcp.go      # 公共类型和工具函数
        ├── scanner.go      # 端口扫描实现
        ├── client.go       # TCP 客户端实现
        └── server.go       # TCP 服务端实现
```

### 5.2 文件职责

| 文件 | 职责 |
|------|------|
| `internal/cli/tcp.go` | 定义 tcp 一级命令，注册三个子命令 |
| `internal/cli/tcp/scan.go` | scan 子命令的 flag 定义和配置组装 |
| `internal/cli/tcp/client.go` | client 子命令的 flag 定义和配置组装 |
| `internal/cli/tcp/server.go` | server 子命令的 flag 定义和配置组装 |
| `internal/commands/tcp/cmd_tcp.go` | 公共类型定义（Config 结构体、常量、工具函数） |
| `internal/commands/tcp/scanner.go` | 端口扫描业务逻辑实现 |
| `internal/commands/tcp/client.go` | TCP 客户端业务逻辑实现 |
| `internal/commands/tcp/server.go` | TCP 服务端业务逻辑实现 |

---

## 六、关键实现要点

### 6.1 并发控制
- 扫描模式：使用工作池控制并发，默认 `runtime.NumCPU() * 2`
- 服务端：使用信号量限制并发连接数，默认 `runtime.NumCPU() * 2`

### 6.2 路径处理逻辑
```go
func processPath(path string) ([]string, error) {
    // 1. 检查是否是通配符
    if strings.Contains(path, "*") || strings.Contains(path, "?") {
        return filepath.Glob(path)
    }
    
    // 2. 检查是否是目录
    info, err := os.Stat(path)
    if err != nil {
        return nil, err
    }
    
    if info.IsDir() {
        // 读取目录下所有文件（不递归）
        entries, err := os.ReadDir(path)
        // ... 过滤文件
        return files, nil
    }
    
    // 3. 普通文件
    return []string{path}, nil
}
```

### 6.3 响应处理
- 所有发送模式默认等待并处理服务端响应
- 响应内容输出到 stdout
- 支持 `--no-response` 禁用响应等待
- 支持 `--timeout` 设置超时

### 6.4 文件传输协议
```
# 文件头格式
FILE:filename:size\ncontent

# 目录头格式
DIR:dirname:filecount\n
# 响应格式
RECEIVED:filename:size
```

---

## 七、测试用例设计

### 7.1 单元测试
- 端口范围解析测试
- 单个端口扫描测试
- 路径处理逻辑测试（文件/目录/通配符）
- 文件头解析测试

### 7.2 集成测试
- 客户端-服务端通信测试
- 文件传输测试
- 并发扫描测试

### 7.3 边界测试
- 空消息发送
- 大文件传输（1GB）
- 全端口范围扫描（1-65535）
- 无效地址处理
- 端口占用处理
- 连接中断处理

---

## 八、待办事项

### 8.1 文件创建清单
- [ ] `internal/cli/tcp.go` - 一级命令定义
- [ ] `internal/cli/tcp/scan.go` - 扫描子命令
- [ ] `internal/cli/tcp/client.go` - 客户端子命令
- [ ] `internal/cli/tcp/server.go` - 服务端子命令
- [ ] `internal/commands/tcp/cmd_tcp.go` - 公共类型和常量
- [ ] `internal/commands/tcp/scanner.go` - 扫描业务逻辑
- [ ] `internal/commands/tcp/client.go` - 客户端业务逻辑
- [ ] `internal/commands/tcp/server.go` - 服务端业务逻辑

### 8.2 功能实现清单
- [ ] 端口范围解析（支持多种格式）
- [ ] 并发端口扫描
- [ ] 扫描结果格式化输出（table/json/csv）
- [ ] TCP 客户端连接管理
- [ ] 字符串发送功能
- [ ] 文件发送功能
- [ ] 目录发送功能（非递归）
- [ ] 通配符匹配功能
- [ ] 交互式发送模式
- [ ] TCP 服务端监听
- [ ] 并发连接处理
- [ ] 文件接收功能
- [ ] 响应消息发送
- [ ] 回声模式

### 8.3 测试清单
- [ ] 端口范围解析单元测试
- [ ] 端口扫描单元测试
- [ ] 客户端-服务端集成测试
- [ ] 文件传输测试
- [ ] 边界条件测试

### 8.4 注册命令
在 `internal/cli/root.go` 的 `SubCmds` 列表中添加 `TcpCmd`

---

## 九、参考文档

1. `docs/fck-analysis.md` - 项目架构分析
2. `docs/qflag命令行工具开发规范.md` - 目录结构和命令定义规范
3. `docs/qflag命令开发规范.md` - 业务逻辑和 CLI 定义规范
4. `docs/tcp-command-refactor-design.md` - 完整设计方案（本文档的详细版本）

---

## 十、关键代码片段

### 10.1 获取 CPU 核心数默认值
```go
import "runtime"

func getDefaultConcurrent() int {
    return runtime.NumCPU() * 2
}
```

### 10.2 路径处理函数
```go
func resolvePath(path string) ([]string, error) {
    // 通配符匹配
    if strings.Contains(path, "*") || strings.Contains(path, "?") {
        matches, err := filepath.Glob(path)
        if err != nil {
            return nil, fmt.Errorf("invalid glob pattern: %w", err)
        }
        return matches, nil
    }
    
    info, err := os.Stat(path)
    if err != nil {
        return nil, fmt.Errorf("failed to stat path: %w", err)
    }
    
    // 目录处理
    if info.IsDir() {
        entries, err := os.ReadDir(path)
        if err != nil {
            return nil, fmt.Errorf("failed to read directory: %w", err)
        }
        
        var files []string
        for _, entry := range entries {
            if !entry.IsDir() {
                files = append(files, filepath.Join(path, entry.Name()))
            }
        }
        return files, nil
    }
    
    // 普通文件
    return []string{path}, nil
}
```

---

## 十一、注意事项

### 11.1 安全风险
- 端口扫描可能被防火墙检测，提供可配置的并发数和超时
- 文件传输需验证文件名，防止目录遍历攻击
- 限制最大并发连接数，防止资源耗尽

### 11.2 性能考虑
- 使用连接池控制并发
- 文件传输使用流式读写
- 服务端使用信号量限制并发

### 11.3 兼容性
- 支持 Windows、Linux、macOS
- 处理不同系统的换行符差异
- 支持 IPv4 和 IPv6

---

**文档结束**

---

## 附录：对话时间线

| 时间 | 事件 |
|------|------|
| 初始 | 用户要求基于 qflag 规范设计 TCP 命令重构方案 |
| | AI 分析项目结构和现有命令实现 |
| | AI 完成完整设计方案文档 |
| 修改 1 | 用户反馈：并发扫描默认值改为 CPU核心数*2 |
| | AI 修改文档 |
| 修改 2 | 用户反馈：合并 --file 和 --dir 为 --path，支持通配符 |
| | AI 修改文档，调整互斥组 |
| 修改 3 | 用户反馈：服务端最大并发连接数也改为 CPU核心数*2 |
| | AI 修改文档 |
| 修改 4 | 用户反馈：添加客户端响应处理说明 |
| | AI 添加响应处理说明文档 |
| 结束 | 用户要求导出对话历史 |
| | AI 整理并导出本对话历史文档 |
