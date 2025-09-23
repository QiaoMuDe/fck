# Find 命令并发模式设计文档

## 概述

本文档描述了为 `fck find` 命令新增并发遍历模式的设计方案。该方案基于工作队列模式，实现类似于 `fd` 工具的高性能并行文件系统遍历。

## 设计目标

1. **性能提升**: 在大型目录结构中实现 2-4倍 的性能提升
2. **向下兼容**: 保持现有 API 和命令行接口不变
3. **简单易用**: 通过布尔标志控制执行模式
4. **资源控制**: 合理控制并发度，避免系统资源耗尽
5. **错误处理**: 在并发环境下正确处理和报告错误

## 核心架构

### 1. 执行模式控制

新增 `-X/--parallel` 布尔标志，控制执行模式：

```bash
# 单线程模式（默认，保持兼容性）
fck find <path>
fck find -X=false <path>

# 并发模式
fck find -X <path>
fck find -X=true <path>
```

### 2. 并发架构设计

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Main Thread   │───▶│   Path Queue     │◀───│  Worker Pool    │
│                 │    │ (Buffer: 10000)  │    │ (CPU Cores)     │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                │                        │
                                ▼                        ▼
                       ┌──────────────────┐    ┌─────────────────┐
                       │   Result Chan    │    │  Error Chan     │
                       │  (Match Results) │    │ (Error Collect) │
                       └──────────────────┘    └─────────────────┘
```

### 3. 关键组件

#### ParallelSearcher 结构体
```go
type ParallelSearcher struct {
    config       *types.FindConfig
    matcher      *PatternMatcher
    operator     *FileOperator
    
    // 并发控制
    pathQueue    chan string
    resultQueue  chan SearchResult
    errorQueue   chan error
    
    // 同步控制
    wg           sync.WaitGroup
    ctx          context.Context
    cancel       context.CancelFunc
}
```

#### SearchResult 结构体
```go
type SearchResult struct {
    Path     string
    Entry    os.DirEntry
    IsMatch  bool
    Error    error
}
```



## 实现细节

### 1. 标志定义更新

在 `flags.go` 中新增：

```go
var (
    // 新增并发模式标志
    findCmdParallel *qflag.BoolFlag // parallel 标志
)

// 在 InitFindCmd() 中添加
findCmdParallel = findCmd.Bool("parallel", "X", false, 
    "启用并发模式进行文件搜索，默认使用单线程模式")
```

### 2. 核心算法流程

#### 主搜索流程
```go
func (ps *ParallelSearcher) Search(rootPath string) error {
    // 1. 初始化
    ps.initializeChannels()
    ps.startWorkers()
    ps.startResultProcessor()
    
    // 2. 添加根路径
    ps.pathQueue <- rootPath
    
    // 3. 等待完成
    ps.waitForCompletion()
    
    // 4. 清理资源
    ps.cleanup()
    
    return ps.collectErrors()
}
```

#### 工作协程逻辑
```go
func (ps *ParallelSearcher) worker(workerID int) {
    defer ps.wg.Done()
    
    for {
        select {
        case path := <-ps.pathQueue:
            ps.processPath(path, workerID)
        case <-ps.ctx.Done():
            return
        }
    }
}

func (ps *ParallelSearcher) processPath(path string, workerID int) {
    entries, err := os.ReadDir(path)
    if err != nil {
        ps.handleError(err, path)
        return
    }
    
    for _, entry := range entries {
        fullPath := filepath.Join(path, entry.Name())
        
        // 应用过滤条件
        if ps.shouldSkip(entry, fullPath) {
            continue
        }
        
        if entry.IsDir() {
            // 目录：添加到队列继续遍历
            select {
            case ps.pathQueue <- fullPath:
            case <-ps.ctx.Done():
                return
            default:
                // 队列满时的处理策略
                ps.handleQueueFull(fullPath)
            }
        } else {
            // 文件：检查匹配并发送结果
            result := ps.checkMatch(entry, fullPath)
            ps.sendResult(result)
        }
    }
}
```

### 3. 性能优化策略

#### 内存池优化
```go
var (
    // 复用 SearchResult 对象
    resultPool = sync.Pool{
        New: func() interface{} {
            return &SearchResult{}
        },
    }
    
    // 复用路径字符串切片
    pathSlicePool = sync.Pool{
        New: func() interface{} {
            return make([]string, 0, 100)
        },
    }
)
```

### 4. 错误处理机制

#### 错误分类和处理
```go
type ErrorType int

const (
    ErrorTypePermission ErrorType = iota
    ErrorTypeNotFound
    ErrorTypeIO
    ErrorTypeTimeout
)

type SearchError struct {
    Type    ErrorType
    Path    string
    Err     error
    Worker  int
}

func (ps *ParallelSearcher) handleError(err error, path string) {
    searchErr := &SearchError{
        Type: ps.classifyError(err),
        Path: path,
        Err:  err,
    }
    
    select {
    case ps.errorQueue <- searchErr:
    default:
        // 错误队列满时的处理
    }
}
```

## 配置参数

### 默认配置
```go
type ParallelConfig struct {
    Enabled   bool          // 是否启用并发模式
    Timeout   time.Duration // 超时时间
    MaxMemory int64         // 最大内存使用
}

var DefaultParallelConfig = ParallelConfig{
    Enabled:   false, // 默认关闭，保持兼容性
    Timeout:   30 * time.Second,
    MaxMemory: 100 * 1024 * 1024, // 100MB
}

// 固定的并发参数
const (
    DefaultWorkerCount = runtime.NumCPU() // 工作协程数 = CPU核心数
    DefaultQueueSize   = 10000            // 队列缓冲大小 = 10000
)
```

## 兼容性保证

### 1. API 兼容性
- 保持现有 `FileSearcher` 接口不变
- 通过工厂模式选择实现：
```go
func NewSearcher(config *types.FindConfig) Searcher {
    if config.ParallelMode {
        return NewParallelSearcher(config)
    }
    return NewSequentialSearcher(config)
}
```

### 2. 输出兼容性
- 并发模式下保持输出格式一致
- 可选择是否保持输出顺序（性能 vs 顺序权衡）

## 测试策略

### 1. 单元测试
- 工作协程逻辑测试
- 队列管理测试
- 错误处理测试

### 2. 集成测试
- 大型目录结构测试
- 并发安全性测试
- 内存泄漏测试
- 性能基准测试

### 3. 压力测试
```go
func BenchmarkParallelSearch(b *testing.B) {
    // 创建测试目录结构
    testDir := createLargeTestDir(10000, 5) // 10000个文件，5层深度
    defer os.RemoveAll(testDir)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        searcher := NewParallelSearcher(testConfig)
        searcher.Search(testDir)
    }
}
```

## 监控和调试

### 1. 调试模式
```bash
# 启用详细日志
fck find -X --debug <path>

# 显示性能统计
fck find -X --stats <path>
```

## 部署和发布

### 1. 渐进式发布
- Phase 1: 实现基础并发功能，默认关闭
- Phase 2: 完善错误处理和性能优化
- Phase 3: 根据用户反馈优化并推广

### 2. 配置管理
```yaml
# ~/.fck/config.yaml
find:
  parallel: false      # 默认单线程模式
```

## 风险评估

### 1. 技术风险
- **内存使用**: 并发可能增加内存消耗
- **文件句柄**: 大量并发可能耗尽文件句柄
- **竞态条件**: 并发访问共享资源的风险

### 2. 缓解措施
- 实现内存监控和限制
- 控制并发度和队列大小
- 使用原子操作和适当的同步机制

## 总结

本设计方案通过引入 `-X` 布尔标志实现了简单直观的执行模式控制，在保持向下兼容的同时，为用户提供了高性能的并发文件搜索能力。通过工作队列模式和合理的资源控制，预期可以在大型目录结构中实现显著的性能提升。

该方案的核心优势：
1. **高性能**: 并发处理提升搜索速度
2. **简单易用**: 布尔标志控制，用户友好
3. **兼容性**: 保持现有接口和行为不变
4. **可靠性**: 完善的错误处理和资源控制
5. **零配置**: 自动使用最优的并发参数

## 使用示例

```bash
# 基本用法
fck find .                    # 单线程模式（默认）
fck find -X .                 # 并发模式（自动使用CPU核心数个协程，队列大小10000）

# 配合其他参数
fck find -X -n "*.go" .       # 并发查找Go文件
fck find -X -t f .            # 并发查找文件
fck find -X -s "+1M" .        # 并发查找大于1M的文件
```

## 实现要点

- **工作协程数**: 固定为CPU核心数，无需用户配置
- **队列大小**: 固定为10000，提供足够的缓冲能力
- **零配置**: 用户只需要 `-X` 标志即可启用最优并发模式
- **自动优化**: 系统自动选择最适合的并发参数