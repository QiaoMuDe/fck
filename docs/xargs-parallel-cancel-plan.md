# xargs 并行执行取消机制修复方案

## 问题描述

当前 `runParallel` 函数在 `-e`（出错停止）模式下，当某个批次执行失败时：
1. 已经启动的 goroutine 会继续执行完成（浪费资源）
2. 阻塞在信号量的批次仍会被启动

## 修复方案

### 方案一：使用 context.Context（推荐）

通过 `context.WithCancel` 创建可取消的上下文，出错时取消后续执行。

```go
// runParallel 并行执行批次（支持取消）
func runParallel(batches [][]string, config XargsConfig, stats *XargsStats) error {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    var wg sync.WaitGroup
    semaphore := make(chan struct{}, config.MaxProcs)
    var mu sync.Mutex
    var firstErr error

    for _, batch := range batches {
        // 检查是否已取消
        select {
        case <-ctx.Done():
            break // 已取消，停止启动新任务
        default:
        }

        wg.Add(1)
        
        select {
        case semaphore <- struct{}{}: // 获取信号量
        case <-ctx.Done(): // 等待期间被取消
            wg.Done()
            break
        }

        go func(b []string) {
            defer wg.Done()
            defer func() { <-semaphore }()

            // 检查是否已取消
            select {
            case <-ctx.Done():
                return
            default:
            }

            if err := executeBatch(b, config, stats); err != nil {
                mu.Lock()
                if config.ExitOnError && firstErr == nil {
                    firstErr = err
                    cancel() // 取消其他任务
                }
                mu.Unlock()
            }
        }(batch)
    }

    wg.Wait()
    return firstErr
}
```

**优点**：
- 标准做法，符合 Go 最佳实践
- 可以优雅地取消已启动的任务
- 新任务不会启动

**缺点**：
- 需要修改 `executeBatch` 签名，支持 `context.Context`
- 代码复杂度增加

---

### 方案二：使用原子标志位（简单）

使用 `atomic.Bool` 作为取消标志，轻量级实现。

```go
// runParallel 并行执行批次（使用原子标志）
func runParallel(batches [][]string, config XargsConfig, stats *XargsStats) error {
    var wg sync.WaitGroup
    semaphore := make(chan struct{}, config.MaxProcs)
    var mu sync.Mutex
    var firstErr error
    var cancelled atomic.Bool

    for _, batch := range batches {
        // 检查是否已取消
        if cancelled.Load() {
            break
        }

        wg.Add(1)
        semaphore <- struct{}{} // 获取信号量

        go func(b []string) {
            defer wg.Done()
            defer func() { <-semaphore }()

            // 检查是否已取消
            if cancelled.Load() {
                return
            }

            if err := executeBatch(b, config, stats); err != nil {
                mu.Lock()
                if config.ExitOnError && firstErr == nil {
                    firstErr = err
                    cancelled.Store(true) // 设置取消标志
                }
                mu.Unlock()
            }
        }(batch)
    }

    wg.Wait()
    return firstErr
}
```

**优点**：
- 实现简单，不需要修改其他函数
- 性能开销小

**缺点**：
- 已启动的 goroutine 需要主动检查标志
- 不如 context 标准

---

### 方案三：提前退出（最小改动）

只在 `executeBatch` 层面优化，不等待所有任务完成。

```go
// runParallel 并行执行批次（提前退出）
func runParallel(batches [][]string, config XargsConfig, stats *XargsStats) error {
    var wg sync.WaitGroup
    semaphore := make(chan struct{}, config.MaxProcs)
    var mu sync.Mutex
    var firstErr error
    done := make(chan struct{}) // 完成信号

    for _, batch := range batches {
        select {
        case <-done: // 已有错误，停止启动新任务
            break
        default:
        }

        wg.Add(1)
        semaphore <- struct{}{}

        go func(b []string) {
            defer wg.Done()
            defer func() { <-semaphore }()

            if err := executeBatch(b, config, stats); err != nil {
                mu.Lock()
                if config.ExitOnError && firstErr == nil {
                    firstErr = err
                    close(done) // 通知其他 goroutine
                }
                mu.Unlock()
            }
        }(batch)
    }

    wg.Wait()
    return firstErr
}
```

**优点**：
- 改动最小
- 新任务不会启动

**缺点**：
- 已启动的任务仍会继续执行
- 使用 `close(done)` 只能通知一次

---

## 推荐方案

**方案二（原子标志位）**，原因：
1. 实现简单，风险低
2. 不需要大规模修改代码
3. 能满足基本需求（阻止新任务启动）
4. 已启动的任务可以通过检查标志提前退出

## 实施步骤

1. 添加 `sync/atomic` 导入
2. 修改 `runParallel` 函数
3. 可选：在 `executeBatch` 中添加取消检查点
4. 测试验证

## 兼容性

- 不影响现有功能
- 只在 `-e` 模式下生效
- 默认行为不变
