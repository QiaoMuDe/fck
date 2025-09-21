# 并发目录搜索重新设计方案

## 当前实现的问题分析

### 现有架构问题
1. **单线程遍历瓶颈**: 使用 `filepath.WalkDir` 进行单线程目录遍历，成为性能瓶颈
2. **通道传输开销**: 路径通过 channel 传输，增加了内存分配和上下文切换开销
3. **工作负载不均**: 某些 worker 可能处理简单路径，而其他 worker 处理复杂路径
4. **内存使用效率低**: 大量路径字符串在 channel 中缓存，占用内存

### 参考 Rust fd 的优势
- **并发目录读取**: 每个协程独立处理目录分支
- **就地处理**: 遇到匹配项立即处理，无需通道传输
- **动态负载均衡**: 工作窃取机制实现负载均衡
- **内存效率**: 减少不必要的字符串复制和缓存

## 新架构设计

### 核心设计理念
```
传统模式: 单线程遍历 -> Channel -> 多Worker处理
新模式:   多Worker并发遍历 + 就地处理 + 工作窃取
```

### 1. 工作单元设计

```go
// WorkItem 表示一个工作单元（目录或文件）
type WorkItem struct {
    Path     string    // 路径
    Depth    int       // 当前深度
    IsDir    bool      // 是否为目录
    Priority int       // 优先级（深度越浅优先级越高）
}

// WorkQueue 工作队列（支持优先级和工作窃取）
type WorkQueue struct {
    items    []WorkItem
    mutex    sync.RWMutex
    cond     *sync.Cond
    closed   bool
}
```

### 2. Worker 架构设计

```go
// ConcurrentWorker 并发工作器
type ConcurrentWorker struct {
    id          int
    searcher    *FileSearcher
    workQueue   *WorkQueue
    localQueue  []WorkItem        // 本地工作队列
    stats       *WorkerStats      // 统计信息
    ctx         context.Context
    cancel      context.CancelFunc
}

// WorkerPool 工作器池
type WorkerPool struct {
    workers     []*ConcurrentWorker
    workQueue   *WorkQueue
    maxWorkers  int
    stats       *PoolStats
}
```

### 3. 并发处理流程

#### 阶段1: 初始化
```
1. 创建全局工作队列
2. 启动 N 个 Worker 协程
3. 将根目录作为初始工作项加入队列
```

#### 阶段2: 并发处理
```
每个 Worker 的处理循环:
1. 从全局队列获取工作项（优先本地队列）
2. 如果是目录：
   - 使用 os.ReadDir() 读取目录内容
   - 对每个子项进行匹配检查
   - 匹配的文件立即输出结果
   - 子目录加入本地队列继续处理
3. 如果是文件：
   - 直接进行匹配检查和输出
4. 本地队列为空时，尝试从其他 Worker 窃取工作
```

#### 阶段3: 工作窃取
```
当 Worker 空闲时：
1. 检查全局队列是否有新工作
2. 随机选择其他 Worker 尝试窃取工作
3. 窃取优先级低的工作项（深度较深的目录）
4. 所有队列都空时，Worker 进入等待状态
```

### 4. 关键算法实现

#### 4.1 目录读取算法
```go
func (w *ConcurrentWorker) processDirectory(path string, depth int) error {
    // 检查深度限制
    if w.exceedsMaxDepth(depth) {
        return nil
    }
    
    // 读取目录内容
    entries, err := os.ReadDir(path)
    if err != nil {
        return w.handleError(path, err)
    }
    
    // 分类处理
    var subdirs []WorkItem
    for _, entry := range entries {
        fullPath := filepath.Join(path, entry.Name())
        
        // 跳过隐藏文件检查
        if w.shouldSkip(fullPath, entry) {
            continue
        }
        
        if entry.IsDir() {
            // 子目录加入本地队列
            subdirs = append(subdirs, WorkItem{
                Path:     fullPath,
                Depth:    depth + 1,
                IsDir:    true,
                Priority: depth + 1,
            })
        } else {
            // 文件立即处理
            w.processFile(fullPath, entry)
        }
    }
    
    // 将子目录加入本地队列
    w.addToLocalQueue(subdirs)
    return nil
}
```

#### 4.2 工作窃取算法
```go
func (w *ConcurrentWorker) stealWork() *WorkItem {
    // 随机选择目标 Worker
    targetID := rand.Intn(len(w.pool.workers))
    if targetID == w.id {
        return nil
    }
    
    target := w.pool.workers[targetID]
    return target.stealFromLocal()
}

func (w *ConcurrentWorker) stealFromLocal() *WorkItem {
    w.localMutex.Lock()
    defer w.localMutex.Unlock()
    
    if len(w.localQueue) == 0 {
        return nil
    }
    
    // 窃取优先级最低的工作项（队列末尾）
    item := w.localQueue[len(w.localQueue)-1]
    w.localQueue = w.localQueue[:len(w.localQueue)-1]
    return &item
}
```

#### 4.3 负载均衡策略
```go
// 动态调整策略
type LoadBalancer struct {
    workerLoads    []int64        // 每个 Worker 的负载
    lastBalance    time.Time      // 上次均衡时间
    balanceInterval time.Duration // 均衡间隔
}

func (lb *LoadBalancer) shouldRebalance() bool {
    return time.Since(lb.lastBalance) > lb.balanceInterval
}

func (lb *LoadBalancer) rebalance(workers []*ConcurrentWorker) {
    // 计算平均负载
    totalLoad := int64(0)
    for _, load := range lb.workerLoads {
        totalLoad += load
    }
    avgLoad := totalLoad / int64(len(workers))
    
    // 重新分配过载 Worker 的工作
    for i, worker := range workers {
        if lb.workerLoads[i] > avgLoad*2 {
            worker.redistributeWork()
        }
    }
}
```

### 5. 性能优化策略

#### 5.1 内存优化
- **对象池**: 复用 WorkItem 对象，减少 GC 压力
- **字符串优化**: 使用 string interning 减少重复路径字符串
- **批量处理**: 批量读取目录内容，减少系统调用

#### 5.2 I/O 优化
- **预读策略**: 根据目录大小调整读取缓冲区
- **异步 I/O**: 对于大目录，使用异步读取
- **缓存策略**: 缓存最近访问的目录信息

#### 5.3 并发优化
- **自适应 Worker 数**: 根据系统负载动态调整 Worker 数量
- **优先级队列**: 优先处理浅层目录，提高响应速度
- **背压控制**: 防止内存使用过多

### 6. 错误处理和监控

#### 6.1 错误处理策略
```go
type ErrorHandler struct {
    maxErrors     int
    errorCount    int64
    errorTypes    map[string]int64
    mutex         sync.RWMutex
}

func (eh *ErrorHandler) handleError(err error, path string) bool {
    eh.mutex.Lock()
    defer eh.mutex.Unlock()
    
    // 分类统计错误
    errorType := classifyError(err)
    eh.errorTypes[errorType]++
    
    // 检查是否超过最大错误数
    if atomic.AddInt64(&eh.errorCount, 1) > int64(eh.maxErrors) {
        return false // 停止处理
    }
    
    return true // 继续处理
}
```

#### 6.2 性能监控
```go
type PerformanceMonitor struct {
    startTime      time.Time
    processedFiles int64
    processedDirs  int64
    workerStats    []WorkerStats
}

type WorkerStats struct {
    ProcessedItems int64
    IdleTime       time.Duration
    WorkStolen     int64
    WorkGiven      int64
}
```

### 7. 配置参数

```go
type ConcurrentConfig struct {
    MaxWorkers      int           // 最大 Worker 数
    QueueSize       int           // 队列大小
    StealThreshold  int           // 工作窃取阈值
    BalanceInterval time.Duration // 负载均衡间隔
    BatchSize       int           // 批处理大小
    EnableProfiling bool          // 启用性能分析
}
```

## 预期性能提升

### 理论分析
1. **并发度提升**: 从单线程遍历提升到多线程并发遍历
2. **减少通道开销**: 消除路径字符串的通道传输
3. **负载均衡**: 工作窃取机制实现动态负载均衡
4. **内存效率**: 减少不必要的内存分配和复制

### 预期指标
- **吞吐量**: 提升 2-4 倍（取决于目录结构）
- **内存使用**: 减少 30-50%
- **响应时间**: 减少 40-60%
- **CPU 利用率**: 提升到 80-90%

## 实现计划

### 阶段1: 核心架构 (1-2天)
- 实现 WorkItem 和 WorkQueue
- 实现基础的 ConcurrentWorker
- 实现简单的工作分配机制

### 阶段2: 工作窃取 (1天)
- 实现工作窃取算法
- 实现负载均衡机制
- 添加性能监控

### 阶段3: 优化和测试 (1-2天)
- 内存和 I/O 优化
- 错误处理完善
- 性能测试和调优

### 阶段4: 集成和验证 (1天)
- 与现有代码集成
- 功能验证和回归测试
- 文档更新

## 风险评估

### 技术风险
- **复杂性增加**: 并发控制和工作窃取增加了代码复杂性
- **调试困难**: 多线程问题难以重现和调试
- **内存竞争**: 多个 Worker 访问共享数据结构

### 缓解措施
- **充分测试**: 包括单元测试、集成测试和压力测试
- **渐进式实现**: 分阶段实现，每个阶段都进行验证
- **性能基准**: 建立性能基准，确保改进效果
- **回退机制**: 保留原有实现作为备选方案

## 总结

这个新设计将显著提升并发搜索的性能，特别是在处理大型目录结构时。通过消除单线程遍历瓶颈、减少通道开销、实现动态负载均衡，预期能够实现 2-4 倍的性能提升。

关键创新点：
1. **多线程并发目录遍历**替代单线程 WalkDir
2. **就地处理**替代通道传输
3. **工作窃取**实现动态负载均衡
4. **优先级队列**优化处理顺序

这个设计借鉴了 Rust fd 的成功经验，同时结合 Go 语言的特性进行了适配和优化。