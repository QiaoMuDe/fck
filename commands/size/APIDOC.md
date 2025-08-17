# Package size

Package size 实现了文件和目录大小计算功能。该文件提供了计算文件、目录大小的核心功能，支持通配符路径展开、隐藏文件处理、进度显示和表格输出。

Package size 定义了 size 子命令的命令行标志和参数配置。该文件包含 size 命令支持的所有选项，如颜色输出、表格样式、隐藏文件处理等。

## FUNCTIONS

### InitSizeCmd

初始化

```go
func InitSizeCmd() *cmd.Cmd
```

### SizeCmdMain

SizeCmdMain 是 size 子命令的主函数

```go
func SizeCmdMain(cl *colorlib.ColorLib) error
```

参数：
- `cl`: 用于打印输出的 ColorLib 对象

返回：
- `error`: 如果发生错误，返回错误信息，否则返回 nil