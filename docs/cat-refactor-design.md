# cat 命令重构设计方案

## 一、问题分析

当前实现的问题：
1. 文件处理和管道输入是两套独立逻辑
2. 管道输入不支持分页器 (`-l` 标志)
3. 代码重复：行分割、换行符统一、标志处理等逻辑分散
4. 分页器模式 (`runOVMode`) 只支持文件，不支持管道

## 二、重构目标

1. **统一内容处理**：文件和管道输入使用同一套处理逻辑
2. **分页器支持管道**：`-l` 标志对管道输入也生效
3. **代码复用**：提取公共逻辑，减少重复代码
4. **保持兼容性**：现有功能不受影响

## 三、核心设计

### 1. 类型定义 (types.go)

```go
// ContentSource 内容源接口
type ContentSource interface {
    Name() string                    // 名称（文件名或 "stdin"）
    Read() ([]byte, error)           // 读取内容
    Size() (int64, error)            // 获取大小（用于检查限制）
    IsBinary() (bool, error)         // 检测是否为二进制
}

// FileSource 文件内容源
type FileSource struct {
    path string
}

func (f *FileSource) Name() string { ... }
func (f *FileSource) Read() ([]byte, error) { ... }
func (f *FileSource) Size() (int64, error) { ... }
func (f *FileSource) IsBinary() (bool, error) { ... }

// StdinSource 管道内容源
type StdinSource struct {
    name string // 默认 "stdin" 或用户指定
}

func (s *StdinSource) Name() string { ... }
func (s *StdinSource) Read() ([]byte, error) { ... }
func (s *StdinSource) Size() (int64, error) { ... }
func (s *StdinSource) IsBinary() (bool, error) { ... }

// ContentProcessor 内容处理器
type ContentProcessor struct {
    config CatConfig
}

// ProcessorOption 处理器选项
type ProcessorOption func(*ContentProcessor)
```

### 2. 统一处理流程

```
+------------------+
|   ContentSource  |
+--------+---------+
         |
         v
+--------+---------+
|     Read()       |  读取原始内容
+--------+---------+
         |
         v
+--------+---------+
|  normalizeLines  |  统一换行符
+--------+---------+
         |
         v
+--------+---------+
|  processContent  |  统一处理（高亮/标志等）
+--------+---------+
         |
    +----+----+
    |         |
    v         v
+---+----+  +--+---+
| 普通输出 |  | 分页器 |
+--------+  +------+
```

### 3. 改造后的函数结构

```go
// CatCmdMain 主入口
func CatCmdMain(config CatConfig) error {
    // 1. 确定内容源（文件或管道）
    source, err := getContentSource(config)
    if err != nil {
        return err
    }
    
    // 2. 统一处理
    return processContent(source, config)
}

// processContent 统一处理内容
func processContent(source ContentSource, config CatConfig) error {
    // 读取内容
    content, err := source.Read()
    if err != nil { ... }
    
    // 检查大小限制
    if config.MaxSize > 0 && len(content) > int(config.MaxSize) { ... }
    
    // 统一换行符
    content = normalizeLineEndings(content)
    
    // 根据模式处理
    if config.UseLess {
        return outputWithPager(content, source.Name(), config)
    }
    return outputDirectly(content, source.Name(), config)
}

// outputWithPager 使用分页器输出
func outputWithPager(content []byte, name string, config CatConfig) error {
    // 支持语法高亮
    // 支持所有 ov 功能
    // 文件和管道统一处理
}

// outputDirectly 直接输出
func outputDirectly(content []byte, name string, config CatConfig) error {
    // 处理 -n, -b, -E, -T 等标志
    // 支持语法高亮
}
```

## 四、代码结构改造

### 文件结构

```
internal/commands/cat/
├── cmd_cat.go       # 主入口和配置
├── types.go         # 类型定义（新增）：ContentSource 接口和实现
├── processor.go     # 统一处理逻辑（新增）
├── output.go        # 输出相关（普通输出+分页器）
├── viewer.go        # 保留但简化（复用 processor）
└── ov_pager.go      # 合并到 output.go 或保留简化
```

### 关键改动

1. **cmd_cat.go**
   - 简化 `CatCmdMain`，只负责确定内容源
   - 移除 `processFile`、`processStdin` 等重复逻辑

2. **types.go** (新增)
   - 定义 `ContentSource` 接口
   - 定义 `FileSource` 和 `StdinSource` 结构体
   - 实现接口方法
   - 定义 `ContentProcessor` 和相关类型

3. **processor.go** (新增)
   - `NewProcessor()` 创建处理器
   - `Process()` 统一处理入口
   - `normalizeLineEndings()` 换行符统一

4. **output.go** (新增或改造 viewer.go)
   - `outputDirectly()` 直接输出（支持所有标志）
   - `outputWithPager()` 分页器输出（支持高亮）
   - 复用现有的高亮逻辑

## 五、处理流程对比

### 当前流程

```
文件输入:
  CatCmdMain -> processFile -> FileViewer.View -> 输出

管道输入:
  CatCmdMain -> processStdin -> outputLines -> 输出

分页器:
  CatCmdMain -> runOVMode -> 直接读取文件 -> 分页显示
  （不支持管道！）
```

### 重构后流程

```
文件输入:
  CatCmdMain -> getContentSource(文件) -> processContent -> 输出/分页器

管道输入:
  CatCmdMain -> getContentSource(管道) -> processContent -> 输出/分页器

统一处理:
  - 大小检查
  - 换行符统一
  - 二进制检测
  - 语法高亮
  - 标志处理 (-n, -b, -E, -T)
  - 分页器支持
```

## 六、边缘案例处理

1. **空管道输入**：直接返回，不报错
2. **二进制管道输入**：根据 `-a`、`-I` 标志处理
3. **大文件管道输入**：检查大小限制
4. **分页器+管道**：正常支持，文件名显示为 "stdin"
5. **高亮+管道**：尝试根据内容检测语言，或默认不高亮

## 七、实现优先级

1. **P0** - 创建 `ContentSource` 接口和实现
2. **P0** - 实现 `processContent()` 统一处理
3. **P0** - 改造分页器支持 `ContentSource`
4. **P1** - 整合并移除重复代码
5. **P1** - 测试所有场景（文件/管道/分页器/高亮）

## 八、兼容性保证

- 所有现有 CLI 标志行为不变
- 文件处理逻辑不变
- 新增管道输入支持 `-l` 分页器
- 性能不下降（避免重复读取）
