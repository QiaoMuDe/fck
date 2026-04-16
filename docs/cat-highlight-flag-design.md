# Cat 命令语法高亮控制标志设计方案

> **设计日期**: 2026-04-16  
> **版本**: v1.0  
> **状态**: 待评审

---

## 一、需求背景

当前 cat 命令的语法高亮功能仅在分页器模式（`-l/--less`）下可用，且默认启用。用户希望：
1. 为日常查看模式也添加语法高亮支持
2. 新增一个独立的标志来控制语法高亮的启用/禁用
3. 默认禁用语法高亮，避免性能开销和颜色干扰

---

## 二、设计方案

### 2.1 标志设计

| 标志 | 短标志 | 说明 | 默认值 |
|------|--------|------|--------|
| `--highlight` | `-H` | 启用语法高亮 | `false` |

**命名理由**:
- `--highlight`: 语义明确，表示"高亮显示"
- `-H`: 大写 H，与 `-h`（help）区分，且 `H` 代表 Highlight
- 参考其他工具：`bat` 使用 `--color` 和 `--decorations`，`eza` 使用 `--color`

### 2.2 使用场景

```bash
# 1. 日常查看 + 语法高亮
fck cat -H main.go
fck cat --highlight script.py

# 2. 分页查看 + 语法高亮（默认已支持，但新增标志可显式控制）
fck cat -l -H main.go

# 3. 语法高亮 + 行号
fck cat -H -n main.go

# 4. 仅语法高亮，无其他装饰
fck cat -H file.js

# 5. 不启用高亮（默认行为）
fck cat file.go        # 纯文本输出
fck cat -l file.go     # 分页查看，无高亮
```

### 2.3 配置结构变更

```go
// CatConfig cat 命令配置
type CatConfig struct {
    // CLI 参数
    Targets      []string // 目标文件列表
    ShowLineNum  bool     // -n 显示所有行号
    ShowNonBlank bool     // -b 显示非空行行号
    ShowEnd      bool     // -E 显示行尾$
    ShowTabs     bool     // -T 显示制表符为^I
    ShowAll      bool     // -A 等价于 -ET
    ShowNewline  bool     // -N 显示换行符类型
    HeadLines    int      // --head 显示前N行 (0表示全部)
    TailLines    int      // --tail 显示后N行 (0表示全部)
    Quiet        bool     // -q 静默模式 (不显示错误信息)
    Text         bool     // -a, --text 强制将二进制文件视为文本处理
    IgnoreBinary bool     // -I, --ignore-binary 完全忽略二进制文件
    UseLess      bool     // -l, --less 使用分页器查看文件内容
    Highlight    bool     // -H, --highlight 启用语法高亮（新增）

    // 运行时
    LineCounter int // 行号计数器
}
```

### 2.4 实现逻辑

#### 2.4.1 分页器模式（UseLess = true）

```go
if config.UseLess {
    // 传递 Highlight 配置给分页器
    return runOVMode(config)
}
```

**变更点**:
- `ov_pager.go` 中的 `runOVMode` 需要根据 `config.Highlight` 决定是否启用高亮
- 当前分页器默认启用高亮，改为根据标志控制

#### 2.4.2 日常查看模式（UseLess = false）

```go
// 在 processFile 中增加高亮处理分支
if config.Highlight && shouldHighlight(path) {
    return processFileWithHighlight(file, path, config)
}
```

**新增函数**: `processFileWithHighlight`

```go
// processFileWithHighlight 处理带语法高亮的文件
//
// 实现逻辑:
// 1. 读取整个文件内容
// 2. 调用 highlightContent 进行语法高亮
// 3. 将高亮后的内容按行分割
// 4. 逐行输出（保留行号、行尾标记等功能）
```

### 2.5 高亮与现有功能的兼容性

| 功能组合 | 支持情况 | 说明 |
|----------|----------|------|
| `-H` + `-n` | ✅ 支持 | 高亮内容 + 行号 |
| `-H` + `-b` | ✅ 支持 | 高亮内容 + 非空行行号 |
| `-H` + `-E` | ⚠️ 冲突 | 高亮已包含行尾颜色，`-E` 的 `$` 可能显示异常 |
| `-H` + `-T` | ⚠️ 冲突 | 高亮已处理制表符，`-T` 的 `^I` 可能显示异常 |
| `-H` + `-N` | ⚠️ 冲突 | 高亮内容中显示换行符标记会干扰颜色 |
| `-H` + `-u` | ✅ 支持 | 高亮前 N 行 |
| `-H` + `-d` | ✅ 支持 | 高亮后 N 行 |
| `-H` + `-l` | ✅ 支持 | 分页器 + 高亮 |

**处理策略**:
- 高亮模式启用时，`-E`、`-T`、`-N` 自动禁用（输出警告）
- 或者：高亮优先级高于这些标志，直接忽略

### 2.6 性能考虑

| 场景 | 内存占用 | 处理时间 | 优化建议 |
|------|----------|----------|----------|
| 小文件 (<1MB) | 低 | 快 | 无 |
| 中等文件 (1-10MB) | 中 | 可接受 | 无 |
| 大文件 (>10MB) | 高 | 慢 | 跳过语法高亮，输出警告 |

**建议**: 添加文件大小限制，超过 10MB 的文件自动跳过语法高亮

---

## 三、文件修改清单

### 3.1 CLI 层

**文件**: `internal/cli/cat.go`

```go
// 新增标志变量
var (
    // ... 现有标志 ...
    catHighlight *qflag.BoolFlag // -H, --highlight 启用语法高亮
)

// init 函数中注册标志
catHighlight = CatCmd.Bool("highlight", "H", "启用语法高亮", false)

// runCat 函数中传递配置
config := cat.CatConfig{
    // ... 现有字段 ...
    Highlight: catHighlight.Get(),
}
```

### 3.2 命令逻辑层

**文件**: `internal/commands/cat/cmd_cat.go`

```go
type CatConfig struct {
    // ... 现有字段 ...
    Highlight bool // -H, --highlight 启用语法高亮
}

// processFile 函数中增加分支
if config.Highlight {
    return processFileWithHighlight(file, path, config)
}
```

**新增文件**: `internal/commands/cat/highlight_processor.go`

```go
// processFileWithHighlight 实现
// - 读取文件
// - 调用 highlightContent
// - 逐行输出（处理行号等）
```

### 3.3 分页器层

**文件**: `internal/commands/cat/ov_pager.go`

```go
// runOVMode 函数修改
// 根据 config.Highlight 决定是否启用高亮
if config.Highlight && shouldHighlight(filename) {
    // 启用高亮
} else {
    // 不启用高亮
}
```

---

## 四、测试用例

### 4.1 基本功能测试

```bash
# 测试 1: 启用高亮
fck cat -H main.go

# 测试 2: 禁用高亮（默认）
fck cat main.go

# 测试 3: 分页器 + 高亮
fck cat -l -H main.go

# 测试 4: 分页器 + 无高亮
fck cat -l main.go
```

### 4.2 组合功能测试

```bash
# 测试 5: 高亮 + 行号
fck cat -H -n main.go

# 测试 6: 高亮 + 前 N 行
fck cat -H -u 10 main.go

# 测试 7: 多文件 + 高亮
fck cat -H file1.go file2.go
```

### 4.3 边界情况测试

```bash
# 测试 8: 不支持的文件类型
fck cat -H file.txt

# 测试 9: 二进制文件 + 高亮
fck cat -H -a binary.bin

# 测试 10: 大文件 + 高亮
fck cat -H largefile.log
```

---

## 五、实现优先级

1. **P0**: 添加 `--highlight/-H` 标志和配置字段
2. **P1**: 实现日常查看模式的语法高亮 (`processFileWithHighlight`)
3. **P2**: 修改分页器模式，根据标志控制高亮
4. **P3**: 添加文件大小限制和性能优化
5. **P4**: 更新文档和示例

---

## 六、备选方案

### 方案 B: 使用 `--color` 标志

```bash
fck cat --color=always file.go   # 总是高亮
fck cat --color=auto file.go     # 自动检测（TTY 时高亮）
fck cat --color=never file.go    # 从不高亮
```

**优点**: 与 `ls`、`grep` 等工具一致  
**缺点**: 需要处理 `auto` 模式的 TTY 检测，复杂度较高

### 方案 C: 环境变量控制

```bash
export FCK_HIGHLIGHT=1
fck cat file.go   # 自动启用高亮
```

**优点**: 全局配置，无需每次输入标志  
**缺点**: 不够直观，与现有标志体系不一致

---

**推荐**: 采用方案 A（`-H/--highlight`），简单直观，符合项目现有设计。
