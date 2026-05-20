# curl 下载进度条实现方案分析

## 一、需求概述

为 curl 命令添加类似 wget 的下载进度条功能，当使用 `-o` 参数保存文件到本地时，显示：
- 下载进度百分比
- 已下载大小 / 总大小
- 下载速度
- 预计剩余时间

## 二、progressbar/v3 库核心能力分析

### 2.1 关键 API

```go
// 创建进度条（支持 int64 最大值）
func NewOptions64(max int64, options ...Option) *ProgressBar

// 创建包装了进度条的 Reader
func NewReader(r io.Reader, bar *ProgressBar) Reader

// ProgressBar 本身实现了 io.Writer 接口
func (p *ProgressBar) Write(b []byte) (n int, err error)
```

### 2.2 未知总大小的处理

progressbar 库对 `max = -1` 的处理（源码第 424 行）：
```go
elapsedTime: max == -1,  // 当 max 为 -1 时，启用已用时间显示
```

当总大小未知时，进度条会显示为**旋转动画（spinner）**模式，而不是百分比进度条。

### 2.3 关键配置选项

| 选项 | 作用 |
|------|------|
| `OptionSetDescription(desc)` | 设置描述文字 |
| `OptionShowBytes(true)` | 显示字节数 |
| `OptionShowTotalBytes(true)` | 显示总大小 |
| `OptionSetWidth(50)` | 设置进度条宽度 |
| `OptionSetTheme(ThemeASCII)` | 设置主题样式 |
| `OptionEnableColorCodes(true)` | 启用颜色 |
| `OptionSetPredictTime(true)` | 显示预计剩余时间 |
| `OptionThrottle(duration)` | 限制更新频率 |

## 三、实现方案

### 3.1 方案对比

#### 方案 A：使用 NewReader 包装（推荐）

```go
// 1. 根据 Content-Length 确定总大小
var totalSize int64 = -1
if resp.ContentLength > 0 {
    totalSize = resp.ContentLength
}

// 2. 创建进度条
bar := progressbar.NewOptions64(
    totalSize,
    progressbar.OptionSetDescription("Downloading..."),
    progressbar.OptionShowBytes(true),
    progressbar.OptionShowTotalBytes(true),
    progressbar.OptionSetWidth(50),
    progressbar.OptionSetPredictTime(true),
)

// 3. 创建带进度追踪的 Reader
reader := progressbar.NewReader(resp.Body, bar)

// 4. 流式复制到文件
file, _ := os.Create(outputPath)
defer file.Close()

_, err := io.Copy(file, &reader)
```

**优点**：
- 代码简洁，只需几行
- 自动处理进度更新
- 自动调用 `bar.Finish()` 当 Reader 关闭时

**缺点**：
- 需要处理 Reader 的 Close 时机

#### 方案 B：使用 progressbar.Write 接口

```go
// 1. 创建进度条（同上）
bar := progressbar.NewOptions64(totalSize, ...)

// 2. 创建 MultiWriter，同时写入文件和进度条
file, _ := os.Create(outputPath)
defer file.Close()

// progressbar 实现了 io.Writer
multiWriter := io.MultiWriter(file, bar)

// 3. 流式复制
_, err := io.Copy(multiWriter, resp.Body)

// 4. 手动完成
bar.Finish()
```

**优点**：
- 更直观，符合 io.MultiWriter 模式
- 控制更灵活

**缺点**：
- 需要手动调用 Finish()
- 会实际写入数据到 progressbar（虽然它只计数不存储）

#### 方案 C：分块读取手动更新

```go
buffer := make([]byte, 32*1024)
for {
    n, err := resp.Body.Read(buffer)
    if n > 0 {
        file.Write(buffer[:n])
        bar.Add(n)  // 手动更新进度
    }
    if err == io.EOF {
        break
    }
}
bar.Finish()
```

**优点**：
- 完全控制
- 可以在读取过程中做其他处理

**缺点**：
- 代码冗余
- 需要自行处理缓冲区

### 3.2 推荐方案：方案 B（MultiWriter 模式）

理由：
1. 代码清晰易读
2. 符合 Go 的 io 接口设计哲学
3. 不需要处理 Reader 的包装和关闭
4. 与现有代码结构最匹配

## 四、具体实现步骤

### 4.1 修改点分析

当前代码位置：`internal/commands/curl/cmd_curl.go`

**现有保存逻辑**（第 227-229 行）：
```go
// 保存到文件
if config.Output != "" {
    return os.WriteFile(config.Output, resp.Body, 0644)
}
```

**问题**：
1. 使用 `io.ReadAll` 在 Execute 函数第 58 行已经读取了整个 body 到内存
2. 需要先改为流式处理

### 4.2 完整改造流程

#### 步骤 1：修改 Execute 函数，支持流式处理

```go
func Execute(config Config) error {
    // ... 前面的代码保持不变 ...
    
    // 执行请求
    resp, err := client.Do(req)
    // ... 错误处理 ...
    defer resp.Body.Close()
    
    // 不再使用 io.ReadAll，改为流式处理
    // 根据配置决定输出方式
    if config.Output != "" {
        // 保存到文件，带进度条
        return downloadWithProgress(resp, config)
    }
    
    // 其他模式（输出到控制台）仍然读取全部内容
    body, err := io.ReadAll(resp.Body)
    // ... 后续处理 ...
}
```

#### 步骤 2：实现 downloadWithProgress 函数

```go
// downloadWithProgress 带进度条的文件下载
func downloadWithProgress(resp *http.Response, config Config) error {
    // 创建输出文件
    file, err := os.Create(config.Output)
    if err != nil {
        return fmt.Errorf("failed to create file: %w", err)
    }
    defer file.Close()
    
    // 获取文件总大小
    var totalSize int64 = -1
    if resp.ContentLength > 0 {
        totalSize = resp.ContentLength
    }
    
    // 非静默模式显示进度条
    if !config.Silent {
        // 打印下载信息（类似 wget）
        fmt.Printf("正在下载: %s\n", config.URL)
        if totalSize > 0 {
            fmt.Printf("大小: %s\n", humanizeBytes(totalSize))
        } else {
            fmt.Println("大小: 未知")
        }
        
        // 创建进度条
        bar := progressbar.NewOptions64(
            totalSize,
            progressbar.OptionSetDescription("进度"),
            progressbar.OptionShowBytes(true),
            progressbar.OptionShowTotalBytes(true),
            progressbar.OptionSetWidth(50),
            progressbar.OptionSetPredictTime(true),
            progressbar.OptionEnableColorCodes(true),
        )
        
        // 使用 MultiWriter 同时写入文件和进度条
        multiWriter := io.MultiWriter(file, bar)
        _, err = io.Copy(multiWriter, resp.Body)
        if err != nil {
            return fmt.Errorf("download failed: %w", err)
        }
        
        // 完成进度条
        bar.Finish()
        fmt.Println() // 换行
        
        // 打印完成信息
        fmt.Printf("已保存到: %s\n", config.Output)
    } else {
        // 静默模式直接复制
        _, err = io.Copy(file, resp.Body)
    }
    
    return err
}
```

### 4.3 边界情况处理

| 场景 | 处理方式 |
|------|----------|
| Content-Length = -1（未知） | 使用 `max = -1`，进度条显示为 spinner 动画 |
| Content-Length = 0 | 正常显示，进度条立即完成 |
| 静默模式 (-s) | 不显示进度条，直接下载 |
| 非终端环境 | 检测 `isatty`，非终端不显示进度条 |
| 小文件 (< 1KB) | 可以跳过进度条或快速完成 |
| 下载中断 | 需要处理信号，清理临时文件 |

### 4.4 终端检测

```go
import "golang.org/x/term"

// 检测是否是终端
func isTerminal() bool {
    return term.IsTerminal(int(os.Stdout.Fd()))
}

// 使用
if !config.Silent && isTerminal() {
    // 显示进度条
}
```

## 五、预期输出效果

### 5.1 已知文件大小

```
正在下载: https://example.com/file.zip
大小: 10.50 MB
进度: [████████████████████░░░░░░░░░░░░░░░░░░░░] 50% | 5.25 MB/10.50 MB | 2.1 MB/s | 剩余 2s
```

### 5.2 未知文件大小

```
正在下载: https://example.com/stream
大小: 未知
进度: [░░░░░░░░░░░░░░░░░░░░] ?% | 1.25 MB/? | 500 kB/s | 已用 3s
```

显示为旋转动画：`- \ | /`

### 5.3 静默模式

```bash
$ fck curl -s -o file.zip https://example.com/file.zip
# 无输出，直接完成
```

## 六、代码修改清单

| 文件 | 修改类型 | 说明 |
|------|----------|------|
| `internal/commands/curl/cmd_curl.go` | 修改 | 添加 `downloadWithProgress` 函数，修改 `Execute` 函数支持流式下载 |
| `internal/commands/curl/types.go` | 可选 | 如需要可添加 `ShowProgress` 配置项 |
| `internal/cli/curl.go` | 可选 | 如需要可添加 `--progress` 或 `-#` 参数 |

## 七、注意事项

1. **流式 vs 全量读取**：只有保存到文件时才使用流式下载，其他模式（如 verbose、include 等）仍需要读取完整 body 进行处理

2. **进度条库已存在**：项目已在 `go.mod` 中依赖 `github.com/schollz/progressbar/v3`，无需额外安装

3. **线程安全**：progressbar 是线程安全的，可以在并发场景下使用

4. **性能影响**：进度条更新会有轻微性能开销，可通过 `OptionThrottle` 限制更新频率

5. **错误处理**：下载中断时需要确保临时文件被清理，可以考虑使用 `.part` 临时文件模式

## 八、参考代码片段

### 8.1 字节人性化显示

```go
func humanizeBytes(bytes int64) string {
    const unit = 1024
    if bytes < unit {
        return fmt.Sprintf("%d B", bytes)
    }
    div, exp := int64(unit), 0
    for n := bytes / unit; n >= unit; n /= unit {
        div *= unit
        exp++
    }
    return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
```

### 8.2 临时文件模式（可选增强）

```go
// 下载到临时文件，完成后重命名
tempFile := config.Output + ".part"
file, err := os.Create(tempFile)
// ... 下载 ...
os.Rename(tempFile, config.Output)
```

这样可以避免下载中断时留下不完整的文件。

---

## 总结

使用 **progressbar/v3** 库的 `MultiWriter` 模式是最简洁的实现方案：

1. 创建进度条，根据 `resp.ContentLength` 设置总大小（-1 表示未知）
2. 使用 `io.MultiWriter(file, bar)` 同时写入文件和更新进度
3. 使用 `io.Copy` 完成流式下载
4. 调用 `bar.Finish()` 完成进度条

整个改动大约需要 50-80 行代码，对现有代码结构影响较小。
