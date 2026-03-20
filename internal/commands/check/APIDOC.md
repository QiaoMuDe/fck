# Package check

Package check 实现了文件哈希校验功能。该文件提供了并发文件校验器，用于验证文件的完整性和一致性。

## 功能介绍

### 文件完整性校验命令的主要逻辑

Package check 实现了文件完整性校验命令的主要逻辑。该文件包含 check 子命令的入口函数，负责解析校验文件并执行文件完整性验证。

### 命令标志和参数配置

Package check 定义了文件校验命令的标志和参数配置。该文件负责初始化 check 子命令的命令行参数解析和帮助信息设置。

### 校验文件的解析功能

Package check 实现了校验文件的解析功能。该文件提供了校验文件解析器，用于解析包含文件哈希信息的校验文件格式。

### 校验文件行内容的验证功能

Package check 实现了校验文件行内容的验证功能。该文件提供了校验文件行验证器，用于验证哈希值格式和文件路径的安全性。

## Functions

### CheckCmdMain

CheckCmdMain 是 check 命令的主函数。

```go
func CheckCmdMain(cl *colorlib.ColorLib) error
```

### InitCheckCmd

```go
func InitCheckCmd() *cmd.Cmd
```