# Package hash

Package hash 实现了文件哈希计算命令的主要逻辑。该文件包含 hash 子命令的入口函数，负责参数验证、文件收集和哈希计算任务的执行。

Package hash 实现了文件收集功能，用于哈希计算前的文件路径收集。该文件提供了文件收集器，支持单个文件、目录遍历和通配符匹配等多种文件收集方式。

Package hash 定义了文件哈希计算命令的标志和参数配置。该文件负责初始化 hash 子命令的命令行参数解析和帮助信息设置。

Package hash 实现了文件哈希计算的并发任务管理功能。该文件提供了哈希任务管理器，支持多线程并发计算文件哈希值，并可选择将结果写入文件。

## FUNCTIONS

### HashCmdMain

HashCmdMain 是 hash 子命令的主函数

```go
func HashCmdMain(cl *colorlib.ColorLib) error
```

### InitHashCmd

初始化 hash 命令

```go
func InitHashCmd() *cmd.Cmd
```

## TYPES

### FileWriterWrapper

文件写入器包装

```go
type FileWriterWrapper struct {
	// Has unexported fields.
}
```

### HashResult

哈希计算结果

```go
type HashResult struct {
	FilePath  string // 文件路径
	HashValue string // 哈希值
	Error     error  // 错误信息
}
```

### HashTaskManager

哈希任务管理器

```go
type HashTaskManager struct {
	// Has unexported fields.
}
```

### NewHashTaskManager

创建哈希任务管理器

```go
func NewHashTaskManager(files []string, hashType func() hash.Hash) *HashTaskManager
```

参数：
- `files`: 文件列表
- `hashType`: 哈希类型

返回值：
- `*HashTaskManager`: 哈希任务管理器

### GetStats

获取统计信息

```go
func (m *HashTaskManager) GetStats() (processed, errors int64)
```

返回值：
- `processed`: 已处理的文件数
- `errors`: 错误数

### Run

执行所有哈希任务

```go
func (m *HashTaskManager) Run() []error
```

返回值：
- `[]error`: 错误列表

### WriteRequest

写入请求

```go
type WriteRequest struct {
	Content string     // 要写入的内容
	Done    chan error // 完成通知通道
}
```