# Package find

Package find 实现了文件查找命令的主要逻辑和配置管理。该文件包含 find 子命令的入口函数，负责参数验证、配置创建和搜索执行。

## 功能介绍

### 文件查找结果的彩色输出功能

Package find 实现了文件查找结果的彩色输出功能。该文件提供了根据文件类型（目录、可执行文件、符号链接等）进行彩色显示的功能。

### 并发搜索功能

Package find 实现了文件查找的并发搜索功能。该文件提供了多线程并发搜索器，用于提高大目录结构的搜索性能。

### 命令参数配置

Package find 定义了文件查找命令的标志和参数配置。该文件负责初始化 find 子命令的所有命令行参数、选项和帮助信息设置。

### 模式匹配功能

Package find 实现了文件查找的模式匹配功能。该文件提供了模式匹配器，支持文件名、路径、大小、时间等多种匹配条件，并包含正则表达式缓存机制。

### 文件操作功能

Package find 实现了文件查找的操作功能。该文件提供了文件操作器，支持删除、移动文件以及执行自定义命令等操作。

### 核心搜索逻辑

Package find 实现了文件查找的核心搜索逻辑。该文件提供了文件搜索器，负责遍历目录、应用过滤条件、执行操作和输出结果。

### 参数验证功能

Package find 实现了文件查找命令参数的验证功能。该文件提供了配置验证器，用于验证所有命令行参数的合法性和安全性。

## Functions

### FindCmdMain

FindCmdMain 是 find 子命令的主函数 - 重构后作为协调器。

```go
func FindCmdMain(cl *colorlib.ColorLib) error
```

### InitFindCmd

```go
func InitFindCmd() *cmd.Cmd
```

### TestMain

TestMain 在所有测试运行前执行初始化。

```go
func TestMain(m *testing.M)
```

## TYPES

### ConcurrentSearcher

ConcurrentSearcher 负责并发搜索协调。

```go
type ConcurrentSearcher struct {
	// Has unexported fields.
}
```

#### NewConcurrentSearcher

NewConcurrentSearcher 创建新的并发搜索器。

```go
func NewConcurrentSearcher(searcher *FileSearcher, maxWorkers int) *ConcurrentSearcher
```

- 参数：
  - `searcher`: 基础搜索器
  - `maxWorkers`: 最大并发worker数量
- 返回：
  - `ConcurrentSearcher`: 并发搜索器对象

#### SearchConcurrent

SearchConcurrent 执行并发搜索。

```go
func (cs *ConcurrentSearcher) SearchConcurrent(findPath string) error
```

- 参数：
  - `findPath`: 搜索路径
- 返回：
  - `error`: 错误信息

### ConfigValidator

ConfigValidator 负责验证find命令的所有参数。

```go
type ConfigValidator struct{}
```

#### NewConfigValidator

NewConfigValidator 创建新的配置验证器。

```go
func NewConfigValidator() *ConfigValidator
```

#### ValidateArgs

ValidateArgs 验证find命令的所有参数。

```go
func (v *ConfigValidator) ValidateArgs(findPath string) error
```

#### ValidateFlags

ValidateFlags 验证所有标志参数的合法性。

```go
func (v *ConfigValidator) ValidateFlags() error
```

#### ValidatePath

ValidatePath 验证路径的安全性和有效性。

```go
func (v *ConfigValidator) ValidatePath(findPath string) error
```

### FileOperator

FileOperator 负责所有文件操作：删除、移动、执行命令。

```go
type FileOperator struct {
	// Has unexported fields.
}
```

#### NewFileOperator

NewFileOperator 创建新的文件操作器。

```go
func NewFileOperator(cl *colorlib.ColorLib) *FileOperator
```

#### Delete

Delete 删除匹配的文件或目录。

```go
func (o *FileOperator) Delete(path string, isDir bool) error
```

- 参数：
  - `path`: 要删除的文件/目录
  - `isDir`: 文件/目录类型
- 返回：
  - `error`: 错误信息

#### Execute

Execute 执行指定的命令，支持直接执行和shell执行两种模式。

```go
func (o *FileOperator) Execute(cmdStr, path string) error
```

- 参数：
  - `cmdStr`: 要执行的命令字符串
  - `path`: 文件路径，用于替换命令中的{}占位符
- 返回：
  - `error`: 错误信息
- 执行模式：
  - 默认直接执行: "cat {}" (更安全，性能更好)
  - 使用--use-shell/-us启用shell执行: 支持管道、重定向等shell功能

#### Move

Move 移动匹配的文件或目录到指定位置。

```go
func (o *FileOperator) Move(srcPath, targetPath string) error
```

- 参数：
  - `srcPath`: 源文件/目录路径
  - `targetPath`: 目标路径
- 返回：
  - `error`: 错误信息

### FileSearcher

FileSearcher 负责核心搜索逻辑。

```go
type FileSearcher struct {
	// Has unexported fields.
}
```

#### NewFileSearcher

NewFileSearcher 创建新的文件搜索器。

```go
func NewFileSearcher(config *types.FindConfig, matcher *PatternMatcher, operator *FileOperator) *FileSearcher
```

#### Search

Search 执行文件搜索。

```go
func (s *FileSearcher) Search(findPath string) error
```

- 参数：
  - `findPath`: 查找路径
- 返回：
  - `error`: 搜索错误（如果有）

### PatternMatcher

PatternMatcher 负责所有模式匹配逻辑，包含正则表达式缓存。

```go
type PatternMatcher struct {
	// Has unexported fields.
}
```

#### NewPatternMatcher

NewPatternMatcher 创建新的模式匹配器。

```go
func NewPatternMatcher(maxCacheSize int) *PatternMatcher
```

#### ClearCache

ClearCache 清空正则表达式缓存。

```go
func (m *PatternMatcher) ClearCache()
```

- 该函数负责清空正则表达式缓存, 释放内存。

#### GetCacheSize

GetCacheSize 获取当前缓存大小。

```go
func (m *PatternMatcher) GetCacheSize() int
```

#### GetRegex

GetRegex 获取编译好的正则表达式，使用缓存机制。

```go
func (m *PatternMatcher) GetRegex(pattern string) (*regexp.Regexp, error)
```

- 参数：
  - `pattern`: 正则表达式模式
- 返回：
  - `*regexp.Regexp`: 编译后的正则表达式
  - `error`: 编译错误（如果有）

#### MatchName

MatchName 匹配文件名。

```go
func (m *PatternMatcher) MatchName(name, pattern string, config *types.FindConfig) bool
```

- 该函数负责匹配文件名与给定模式的逻辑。
- 参数：
  - `name`: 文件名
  - `pattern`: 匹配模式
  - `config`: 查找配置
- 返回：
  - `bool`: 是否匹配成功

#### MatchPath

MatchPath 匹配路径。

```go
func (m *PatternMatcher) MatchPath(path, pattern string, config *types.FindConfig) bool
```

- 该函数负责匹配路径与给定模式的逻辑。
- 参数：
  - `path`: 路径
  - `pattern`: 匹配模式
  - `config`: 查找配置
- 返回：
  - `bool`: 是否匹配成功

#### MatchSize

MatchSize 检查文件大小是否符合指定的条件。

```go
func (m *PatternMatcher) MatchSize(fileSize int64, sizeCondition string) bool
```

- 该函数负责检查文件大小是否符合指定的条件, 支持字节、KB、MB、GB 单位。
- 参数：
  - `fileSize`: 文件大小
  - `sizeCondition`: 大小条件, 格式如"+100"表示大于100, "-100"表示小于100
- 返回：
  - `bool`: 是否匹配成功

#### MatchTime

MatchTime 检查文件时间是否符合指定的条件。

```go
func (m *PatternMatcher) MatchTime(fileTime time.Time, timeCondition string) bool
```

- 该函数负责检查文件时间是否符合指定的条件, 支持天单位。
- 参数：
  - `fileTime`: 文件时间
  - `timeCondition`: 时间条件, 格式如"+10"表示10天前, "-10"表示10天后
- 返回：
  - `bool`: 是否匹配成功