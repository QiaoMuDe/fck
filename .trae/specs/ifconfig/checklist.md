# 检查清单

## 命令行参数
- [ ] `fck ifconfig` 默认显示物理网卡信息（表格）
- [ ] `fck ifconfig -a` 显示所有接口（含虚拟网卡）
- [ ] `fck ifconfig -s` 简略输出（名称、状态、IP）
- [ ] `fck ifconfig -j` JSON 格式输出
- [ ] `fck ifconfig --stats` 显示流量统计
- [ ] `fck ifconfig eth0` 只显示指定接口
- [ ] `fck ifconfig eth0 wlan0` 显示多个指定接口
- [ ] `fck ifconfig -ts r` 切换表格样式

## 业务逻辑
- [ ] `internal/commands/ifconfig/cmd_ifconfig.go` 编译通过
- [ ] `internal/cli/ifconfig.go` 编译通过
- [ ] `internal/cli/root.go` 注册 `IfconfigCmd` 后编译通过
- [ ] 接口信息显示正确：名称、MTU、MAC、状态、IP 地址
- [ ] 虚拟网卡默认隐藏，`-a` 时显示
- [ ] 速度字段取不到时显示 `-`
- [ ] JSON 输出格式正确可解析

## 跨平台
- [ ] Windows 上运行正常（接口名含中文/Unicode）
- [ ] gopsutil 流量统计在 Windows 上正常获取
