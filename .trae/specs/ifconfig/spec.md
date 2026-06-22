# ifconfig 网络接口信息命令 Spec

## Why

Windows 缺乏原生美观的网络接口信息查看工具（ipconfig 输出杂乱），标准 `ifconfig` 命令在 Windows 上不可用。新增 `ifconfig` 命令提供跨平台统一的网络接口信息查看能力，包括接口状态、IP 地址、MAC 地址、MTU、链路速度、流量统计等。

## What Changes

- 新增 `internal/commands/ifconfig/` 业务逻辑目录
- 新增 `internal/cli/ifconfig.go` CLI 接口文件
- 在 `internal/cli/root.go` 注册 `ifconfig` 子命令
- 无新增外部依赖（使用 stdlib `net` + 已存在的 `gopsutil` + `go-pretty`）

## Impact

- Affected specs: 网络工具类新增 ifconfig 命令
- Affected code:
  - `internal/commands/ifconfig/cmd_ifconfig.go` — 主业务逻辑
  - `internal/cli/ifconfig.go` — CLI 标志定义与命令注册
  - `internal/cli/root.go` — 添加 `IfconfigCmd` 到 SubCmds

## ADDED Requirements

### Requirement 1: 接口列表显示

系统 SHALL 显示本机所有物理网络接口的信息，默认隐藏虚拟网卡。

#### Scenario: 默认显示
- **WHEN** 用户执行 `fck ifconfig`
- **THEN** 显示所有物理网络接口的表格信息（名称、MTU、MAC、状态、IP 地址）
- **AND** 默认隐藏虚拟网卡（Hyper-V、VMware、WSL、Bluetooth 等）
- **AND** 使用 go-pretty/table 输出彩色表格

#### Scenario: 显示全部接口
- **WHEN** 用户执行 `fck ifconfig -a`
- **THEN** 显示所有接口（含虚拟网卡）

#### Scenario: 指定接口
- **WHEN** 用户执行 `fck ifconfig eth0`
- **THEN** 只显示指定名称的接口信息

#### Scenario: 简略模式
- **WHEN** 用户执行 `fck ifconfig -s`
- **THEN** 以简略格式输出（名称、状态、IP 地址，无表格边框）

### Requirement 2: 流量统计

系统 SHALL 显示每个接口的收发流量统计。

#### Scenario: 显示流量统计
- **WHEN** 用户执行 `fck ifconfig --stats`
- **THEN** 额外显示 TX/RX 字节数、包数、错误数等统计信息

### Requirement 3: JSON 输出

系统 SHALL 支持 JSON 格式输出，便于脚本消费。

#### Scenario: JSON 输出
- **WHEN** 用户执行 `fck ifconfig -j`
- **THEN** 以 JSON 数组格式输出所有接口信息

### Requirement 4: 表格样式

系统 SHALL 支持表格样式切换，复用现有的 `--style` / `-S` 标志。

#### Scenario: 表格样式
- **WHEN** 用户执行 `fck ifconfig -S r`
- **THEN** 使用圆角表格样式输出

### Requirement 5: `--table-style` 标志

系统 SHALL 复用项目已有的 `--table-style` / `-ts` 标志来选择 go-pretty 表格样式，默认禁用边框样式，与其他表格命令（size、list、port、proc、wc、df）保持一致。

#### Scenario: 切换表格样式
- **WHEN** 用户执行 `fck ifconfig -ts r`
- **THEN** 使用圆角样式渲染表格
- **WHEN** 用户不指定 `-ts`
- **THEN** 使用 `"none"` 默认样式（禁用边框）

## MODIFIED Requirements

（无）

## REMOVED Requirements

（无）
