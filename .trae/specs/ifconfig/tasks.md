# Tasks

- [ ] Task 1: 创建 `internal/commands/ifconfig/cmd_ifconfig.go` — 业务逻辑实现
  - [ ] 1.1 定义 `IfconfigConfig` 结构体（All, Short, JSON, Stats, TableStyle, Interfaces 字段）
  - [ ] 1.2 实现 `IfconfigCmdMain()` 入口函数
  - [ ] 1.3 实现接口信息收集：使用 `net.Interfaces()` 获取基本信息
  - [ ] 1.4 实现 IP 地址收集：使用 `interface.Addrs()` 获取 IPv4/IPv6
  - [ ] 1.5 实现虚拟网卡过滤逻辑（名称含 Hyper-V / VirtualBox / VMware / WSL / Bluetooth / Loopback / Tailscale 等关键词判定为虚拟）
  - [ ] 1.6 实现扩展信息收集：使用 gopsutil `net.IOCountersByFile(true)` 获取流量统计
  - [ ] 1.7 实现 gopsutil 接口速度获取，取不到时显示 `-`
  - [ ] 1.8 实现默认表格输出：使用 go-pretty/table 渲染名称、MTU、MAC、状态、IP、IPv6、速度
  - [ ] 1.9 实现简略输出（`-s`）：名称、状态、主 IP
  - [ ] 1.10 实现 JSON 输出（`-j`）：`encoding/json` 序列化
  - [ ] 1.11 实现流量统计输出（`--stats`）：额外显示 TX/RX 字节、包数、错误数
  - [ ] 1.12 实现指定接口过滤：从 config.Interfaces 筛选匹配的接口名

- [ ] Task 2: 创建 `internal/cli/ifconfig.go` — CLI 接口
  - [ ] 2.1 定义 `IfconfigCmd *qflag.Cmd` 和标志变量
  - [ ] 2.2 使用 qflag-cli 注册标志：`-a`/`--all`、`-s`/`--short`、`-j`/`--json`、`--stats`、`-ts`/`--table-style`
  - [ ] 2.3 设置命令描述、使用示例（中英文双语）
  - [ ] 2.4 实现 `runIfconfig()` 构建 `IfconfigConfig` 并调用 `IfconfigCmdMain()`
  - [ ] 2.5 将 `cmd.Args()` 传入 config.Interfaces（指定接口过滤）

- [ ] Task 3: 注册子命令到根命令
  - [ ] 3.1 在 `internal/cli/root.go` 的 `SubCmds` 中添加 `IfconfigCmd`
  - [ ] 3.2 编译验证命令可正常注册

# Task Dependencies

- Task 3 依赖 Task 2
- Task 1 和 Task 2 可并行实现
