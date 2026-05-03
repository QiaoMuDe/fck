# md 命令管道输入支持设计方案

## 一、需求背景

当前 `md` 命令只能从文件读取 Markdown 内容进行渲染预览：

```bash
fck md README.md
```

需要支持从管道/标准输入读取内容：

```bash
echo "# Hello" | fck md
cat README.md | fck md
git log --oneline | head -20 | fck md
```

## 二、设计目标

1. **支持管道输入**：从 stdin 读取 Markdown 内容
2. **与文件参数互斥**：管道输入时不能指定文件参数
3. **统一体验**：管道和文件模式的渲染效果一致
4. **分页器支持**：管道输入也可使用 `-l` 分页器查看

## 三、技术方案

### 3.1 管道检测

使用项目已有的工具函数：

```go
// internal/utils/utils.go
func IsStdinPipe() bool
```

检测逻辑：
- `true`：stdin 是管道或文件重定向
- `false`：stdin 是终端输入

### 3.2 参数处理逻辑

```
if IsStdinPipe() {
    // 管道模式
    if len(args) > 0 {
        return error("cannot specify file when reading from pipe")
    }
    source = os.Stdin
    filename = "stdin"
} else {
    // 文件模式
    if len(args) == 0 {
        return error("no markdown file specified")
    }
    if len(args) > 1 {
        return error("only one markdown file can be specified at a time")
    }
    source = os.Open(args[0])
    filename = args[0]
}
```

### 3.3 oviewer 集成

oviewer 的 `ControlReader` 方法支持从 `io.Reader` 加载内容：

```go
// 创建 Document
doc, err := oviewer.NewDocument()
doc.FileName = filename  // 管道时显示 "stdin"

// 从 Reader 加载（支持文件或 stdin）
err = doc.ControlReader(source, nil)
```

### 3.4 架构调整

#### 3.4.1 CLI 层 (internal/cli/md.go)

```go
func runMd(cmd qflag.Command) error {
    args := cmd.Args()
    
    // 检测管道输入
    isPipe := term.IsStdinPipe()
    
    if isPipe {
        if len(args) > 0 {
            return fmt.Errorf("cannot specify file when reading from pipe")
        }
    } else {
        if len(args) == 0 {
            return fmt.Errorf("no markdown file specified")
        }
        if len(args) > 1 {
            return fmt.Errorf("only one markdown file can be specified at a time")
        }
    }
    
    config := md.MdConfig{
        File:      args[0],  // 管道时为空字符串
        UsePipe:   isPipe,   // 新增字段
        UsePager:  mdLess.Get(),
        ShowRaw:   mdRaw.Get(),
        Style:     mdStyle.Get(),
        WordWidth: mdWidth.Get(),
        MaxSize:   mdMaxSize.Get(),
    }
    
    return md.MdCmdMain(config)
}
```

#### 3.4.2 配置层 (internal/commands/md/cmd_md.go)

```go
type MdConfig struct {
    File      string    // 文件路径（管道模式为空）
    UsePipe   bool      // 是否从管道读取
    UsePager  bool      // 使用分页器
    ShowRaw   bool      // 显示原始文件
    Style     string    // 渲染样式
    WordWidth int       // 换行宽度
    MaxSize   int64     // 最大文件大小
}
```

#### 3.4.3 核心逻辑层 (internal/commands/md/viewer.go)

```go
func (v *MdViewer) Run() error {
    if v.config.UsePipe {
        return v.runFromPipe()
    }
    if v.config.UsePager {
        return v.runWithPager()
    }
    return v.runDirect()
}

// 管道输入处理
func (v *MdViewer) runFromPipe() error {
    // 从 stdin 读取全部内容
    source, err := io.ReadAll(os.Stdin)
    if err != nil {
        return fmt.Errorf("failed to read from pipe: %w", err)
    }
    
    if v.config.UsePager {
        return v.runPagerWithContent(source, "stdin")
    }
    
    // 直接输出模式
    rendered, err := v.renderer.RenderBytes(source)
    if err != nil {
        return fmt.Errorf("failed to render: %w", err)
    }
    fmt.Print(string(rendered))
    return nil
}

// 统一的分页器处理（支持文件和管道）
func (v *MdViewer) runPagerWithContent(content []byte, filename string) error {
    // 创建渲染版文档
    renderDoc, err := oviewer.NewDocument()
    renderDoc.FileName = filename + " (rendered)"
    
    rendered, err := v.renderer.RenderBytes(content)
    if err != nil {
        return fmt.Errorf("failed to render: %w", err)
    }
    
    if err := renderDoc.ControlReader(bytes.NewBuffer(rendered), nil); err != nil {
        return fmt.Errorf("failed to load rendered content: %w", err)
    }
    
    // 根据配置决定是否加载原始版
    if v.config.ShowRaw {
        originalDoc, err := oviewer.NewDocument()
        originalDoc.FileName = filename
        if err := originalDoc.ControlReader(bytes.NewBuffer(content), nil); err != nil {
            return fmt.Errorf("failed to load original content: %w", err)
        }
        
        ov, err := oviewer.NewOviewer(renderDoc, originalDoc)
        if err != nil {
            return fmt.Errorf("failed to create oviewer: %w", err)
        }
        return ov.Run()
    }
    
    // 仅渲染版
    ov, err := oviewer.NewOviewer(renderDoc)
    if err != nil {
        return fmt.Errorf("failed to create oviewer: %w", err)
    }
    return ov.Run()
}
```

## 四、使用示例

### 4.1 管道输入

```bash
# 基本管道输入
echo "# Hello World" | fck md

# 从其他命令输出
cat README.md | fck md
git log --oneline | head -20 | fck md

# 管道 + 分页器
cat README.md | fck md -l

# 管道 + 分页器 + 原始视图
cat README.md | fck md -l -r
```

### 4.2 错误场景

```bash
# 错误：管道输入时指定文件
cat README.md | fck md other.md
# 输出: err: cannot specify file when reading from pipe

# 错误：终端模式无参数
fck md
# 输出: err: no markdown file specified
```

## 五、边缘情况处理

| 场景 | 处理方式 |
|------|----------|
| 空管道 | 正常渲染，输出空内容 |
| 二进制数据 | glamour 渲染失败，返回错误 |
| 超大管道数据 | 受 `-S` 参数限制，超过时提示使用分页器 |
| 管道 + 分页器 | 正常支持，显示 "stdin (rendered)" |
| 管道 + 原始视图 | 支持 `-r` 标志，可在渲染/原始间切换 |

## 六、实现步骤

1. **修改 `internal/commands/md/cmd_md.go`**
   - 添加 `UsePipe` 字段到 `MdConfig`

2. **修改 `internal/commands/md/viewer.go`**
   - 重构 `Run()` 方法，支持管道分支
   - 添加 `runFromPipe()` 方法
   - 提取 `runPagerWithContent()` 统一处理分页逻辑

3. **修改 `internal/cli/md.go`**
   - 更新 `runMd()` 函数，添加管道检测和参数校验

4. **更新帮助文档**
   - 添加管道使用示例
   - 更新注意事项说明

## 七、兼容性

- **向后兼容**：现有文件模式完全不受影响
- **Unix 惯例**：符合 `cat`, `less`, `grep` 等工具的管道处理习惯
- **Windows 支持**：`term.IsStdinPipe()` 已考虑跨平台兼容性

## 八、预估代码量

- `cmd_md.go`: +1 行（添加字段）
- `viewer.go`: +40 行（新增方法和重构）
- `md.go`: +10 行（管道检测和参数校验）

**总计约 50 行代码**
