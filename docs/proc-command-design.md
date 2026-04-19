# fck proc 进程查看工具设计方案

## 1. 命令概述

### 命令名称
`proc` - 进程查看和管理工具

### 命令用途
查看系统进程信息，支持按多种条件过滤、排序和格式化输出，类似于 Linux 的 `ps` 和 Windows 的 `tasklist` 命令的增强版。

### 使用示例
```bash
fck proc                    # 查看所有进程
fck proc -n nginx           # 查看名为 nginx 的进程
fck proc -p 1234            # 查看指定 PID 的进程
fck proc -u username        # 查看指定用户的进程
fck proc -c "chrome"        # 按命令行过滤
fck proc --sort cpu         # 按 CPU 使用率排序
fck proc --tree             # 树形显示进程关系
fck proc -l                 # 简洁模式
```

---

## 2. 功能特性

### 2.1 基础功能
- [ ] 列出所有进程
- [ ] 显示进程详细信息（PID、名称、用户、CPU、内存等）
- [ ] 支持多种输出格式（表格、简洁、JSON）
- [ ] 按多种条件过滤
- [ ] 支持排序

### 2.2 过滤功能
- [ ] 按进程名过滤 (`-n, --name`) - 支持部分匹配
- [ ] 按 PID 过滤 (`-p, --pid`) - 单个 PID
- [ ] 按 PID 列表过滤 (`-P, --pids`) - 多个 PID，如: 1234,5678

> 其他过滤需求可通过管道使用 `grep` 命令实现

### 2.3 排序功能
- [ ] 按 PID 排序
- [ ] 按进程名排序
- [ ] 按 CPU 使用率排序
- [ ] 按内存使用率排序
- [ ] 按启动时间排序
- [ ] 支持升序/降序 (`--asc`, `--desc`)

### 2.4 显示功能
- [ ] 表格模式（默认）
- [ ] 简洁模式 (`-l, --list`)
- [ ] 树形模式 (`--tree`) - 显示进程父子关系
- [ ] JSON 输出 (`--json`)
- [ ] 自定义显示列 (`--columns`)

### 2.5 高级功能（可选）
- [ ] 实时刷新模式 (`--watch`)
- [ ] 查找进程打开的端口（与 port 命令联动）
- [ ] 显示进程环境变量 (`--env`)
- [ ] 显示进程启动命令完整路径 (`--full-cmd`)

---

## 3. 命令行接口设计

### 3.1 标志定义

| 短标志 | 长标志 | 类型 | 默认值 | 说明 |
|--------|--------|------|--------|------|
| -n | --name | string | "" | 按进程名过滤，支持部分匹配 |
| -p | --pid | int | 0 | 指定单个 PID |
| -P | --pids | int[] | [] | 指定多个 PID，如: 1234,5678 |
| -s | --sort | string | "pid" | 排序字段: pid, name, cpu, mem, time |
| | --asc | bool | false | 升序排列（默认降序） |
| -l | --list | bool | false | 简洁模式 |
| | --tree | bool | false | 树形显示进程关系 |
| | --json | bool | false | JSON 格式输出 |
| -h | --help | bool | false | 显示帮助信息 |

**过滤建议**：其他过滤（用户、CPU、内存、状态等）通过管道使用 `grep` 实现

### 3.2 互斥组
- `--list`, `--tree`, `--json` 互斥

### 3.3 示例命令

```bash
# 基础使用
fck proc                           # 显示所有进程
fck proc -n chrome                 # 查找 chrome 进程
fck proc -p 1234                   # 查看 PID 1234 的进程
fck proc -P 1234,5678,9012         # 查看多个指定进程

# 排序
fck proc -s cpu                    # 按 CPU 排序（高到低）
fck proc -s mem --asc              # 按内存排序（低到高）

# 显示模式
fck proc -l                        # 简洁模式
fck proc --tree                    # 树形显示
fck proc --json                    # JSON 输出

# 配合 grep 进行高级过滤
fck proc -l | grep admin           # 过滤用户为 admin 的进程
fck proc -l | grep "CPU:.*[5-9]"   # 过滤 CPU > 5% 的进程
fck proc -l | grep -i running      # 过滤运行中的进程
fck proc | grep 192.168            # 过滤包含特定 IP 的进程
```

---

## 4. 数据结构

### 4.1 进程信息结构体

```go
// ProcInfo 进程信息结构体
type ProcInfo struct {
    // 基础信息
    PID        int32   // 进程 ID
    PPID       int32   // 父进程 ID
    Name       string  // 进程名
    CmdLine    string  // 命令行
    FullCmd    string  // 完整命令行
    
    // 用户信息
    Username   string  // 用户名
    UID        int32   // 用户 ID
    GID        int32   // 组 ID
    
    // 资源使用
    CPUPercent float64 // CPU 使用率
    MemPercent float64 // 内存使用率
    MemRSS     uint64  // 实际使用内存(RSS)
    MemVMS     uint64  // 虚拟内存(VMS)
    
    // 状态信息
    Status     string  // 进程状态
    NumThreads int32   // 线程数
    NumFDs     int32   // 文件描述符数
    
    // 时间信息
    CreateTime int64   // 创建时间(时间戳)
    RunTime    string  // 运行时长
    
    // 网络信息（可选，与 port 命令联动）
    Connections int    // 网络连接数
    ListenPorts []int   // 监听端口列表
}

// ProcConfig 进程命令配置结构体
type ProcConfig struct {
    // 过滤条件（核心）
    Name     string  // 进程名
    PID      int32   // 单个 PID
    PIDs     []int   // 多个 PID
    
    // 排序
    SortBy string
    Ascend bool
    
    // 显示模式
    ListMode bool
    TreeMode bool
    JSONMode bool
}

// ProcStats 统计信息结构体
type ProcStats struct {
    TotalCount   int
    RunningCount int
    SleepingCount int
    StoppedCount int
    ZombieCount  int
    TotalCPU     float64
    TotalMem     uint64
}
```

---

## 5. 实现方案

### 5.1 文件结构

```
internal/
├── cli/
│   └── proc.go              # CLI 定义和参数解析
├── commands/
│   └── proc/
│       └── cmd_proc.go      # 业务逻辑实现
```

### 5.2 核心函数

```go
// ProcCmdMain 进程命令主入口
func ProcCmdMain(config *ProcConfig) error

// collectProcs 收集进程信息
func collectProcs(config *ProcConfig) ([]ProcInfo, *ProcStats, error)

// filterProcs 过滤进程列表
func filterProcs(procs []ProcInfo, config *ProcConfig) []ProcInfo

// sortProcs 排序进程列表
func sortProcs(procs []ProcInfo, sortBy string, ascend bool)

// renderTableMode 表格模式渲染
func renderTableMode(procs []ProcInfo, stats *ProcStats, config *ProcConfig) error

// renderListMode 简洁模式渲染
func renderListMode(procs []ProcInfo) error

// renderTreeMode 树形模式渲染
func renderTreeMode(procs []ProcInfo, config *ProcConfig) error

// renderJSONMode JSON 模式渲染
func renderJSONMode(procs []ProcInfo, stats *ProcStats) error

// buildProcTree 构建进程树
func buildProcTree(procs []ProcInfo) map[int32][]ProcInfo
```

### 5.3 依赖库

- `github.com/shirou/gopsutil/v3/process` - 获取进程信息
- `github.com/jedib0t/go-pretty/v6/table` - 表格格式化
- `gitee.com/MM-Q/qflag` - CLI 框架

---

## 6. 输出格式示例

### 6.1 表格模式（默认）

```
┌──────┬─────────────────┬─────────┬────────┬──────────┬──────────┬──────────┐
│ PID  │ Name            │ User    │ Status │ CPU(%)   │ Mem(MB)  │ Runtime  │
├──────┼─────────────────┼─────────┼────────┼──────────┼──────────┼──────────┤
│ 1234 │ nginx           │ www     │ R      │ 2.5      │ 128.5    │ 2d 3h    │
│ 5678 │ python          │ admin   │ S      │ 15.2     │ 512.0    │ 5h 30m   │
│ 9012 │ chrome          │ admin   │ R      │ 8.7      │ 1024.0   │ 1d 12h   │
├──────┴─────────────────┴─────────┴────────┴──────────┴──────────┴──────────┤
│ Total: 156  Running: 45  Sleeping: 100  Stopped: 1  Zombie: 0               │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 6.2 简洁模式

```
[1234] nginx (www, CPU: 2.5%, MEM: 128.5MB)
[5678] python (admin, CPU: 15.2%, MEM: 512.0MB)
[9012] chrome (admin, CPU: 8.7%, MEM: 1024.0MB)
```

### 6.3 树形模式

```
1 systemd
  ├─ 234 nginx
  │   ├─ 235 nginx-worker
  │   └─ 236 nginx-worker
  ├─ 456 sshd
  │   └─ 789 sshd-session
  │       └─ 790 bash
  └─ 901 python
      └─ 902 python-thread
```

### 6.4 JSON 模式

```json
{
  "total": 156,
  "running": 45,
  "sleeping": 100,
  "processes": [
    {
      "pid": 1234,
      "ppid": 1,
      "name": "nginx",
      "user": "www",
      "cpu_percent": 2.5,
      "mem_percent": 1.2,
      "mem_rss": 134217728,
      "status": "running",
      "create_time": 1704067200
    }
  ]
}
```

---

## 7. 与现有命令的联动

### 7.1 与 port 命令联动

```bash
# 查看进程打开的端口
fck proc -n nginx --ports

# 查看占用特定端口的进程
fck port -P 8080 --proc
```

### 7.2 与 find 命令联动

```bash
# 查找进程相关的文件
fck proc -n nginx --files
```

---

## 8. 性能考虑

1. **缓存机制**：进程名、用户名等信息缓存，避免重复查询
2. **增量更新**：`--watch` 模式下只更新变化的数据
3. **分页显示**：进程数量过多时支持分页
4. **并发收集**：使用 goroutine 并发收集进程信息

---

## 9. 跨平台支持

| 功能 | Linux | Windows | macOS |
|------|-------|---------|-------|
| 基础进程信息 | ✅ | ✅ | ✅ |
| CPU/内存使用率 | ✅ | ✅ | ✅ |
| 进程状态 | ✅ | ✅ | ✅ |
| 用户/组信息 | ✅ | ✅ | ✅ |
| 环境变量 | ✅ | ⚠️ | ✅ |
| 完整命令行 | ✅ | ⚠️ | ✅ |
| 进程树 | ✅ | ✅ | ✅ |

> ⚠️ Windows 部分功能可能受限或需要管理员权限

---

## 10. 实现优先级

### P0 - 核心功能（必须实现）
- [ ] 基础进程列表显示（PID、名称、用户、CPU、内存、状态）
- [ ] 按 PID、名称过滤
- [ ] 表格和简洁模式
- [ ] 按 PID、CPU、内存排序

### P1 - 重要功能（建议实现）
- [ ] 树形显示模式
- [ ] JSON 输出

### P2 - 增强功能（可选实现）
- [ ] 与 port 命令联动

---

## 11. 风险评估

| 风险 | 可能性 | 影响 | 缓解措施 |
|------|--------|------|----------|
| 获取进程信息需要 root 权限 | 中 | 中 | 优雅处理权限错误，提供警告信息 |
| 进程数量过多导致性能问题 | 中 | 中 | 实现分页和并发收集 |
| 跨平台兼容性差异 | 高 | 低 | 使用 gopsutil 库，处理平台差异 |
| 实时刷新模式资源占用高 | 低 | 中 | 限制刷新频率，提供退出机制 |

---

## 12. 总结

`fck proc` 命令将提供一个功能丰富、跨平台的进程查看工具，支持多种过滤、排序和显示模式，与现有的 `port` 命令形成良好的工具链。实现时采用模块化设计，便于后续扩展和维护。
