# Package commands

Package commands 实现了 fck 命令行工具的主要入口和子命令调度功能。该包负责初始化各个子命令（size、list、check、hash、find），解析命令行参数，并根据用户输入调度到相应的子命令执行器。Package commands 提供了主命令的初始化和配置功能。该文件负责设置 fck 工具的版本信息、帮助文档、Logo 显示等全局配置。

## FUNCTIONS

### InitMainCmd

InitMainCmd 初始化主命令

```go
func InitMainCmd()
```

### Run

Run 运行命令行工具

```go
func Run()
```