# sed 命令设计方案

> **需求**: 实现文本替换功能，支持匹配和替换，操作指定文件
> **设计日期**: 2026-04-02
> **状态**: 待实施

---

## 一、功能概述

实现一个简化版的 sed 命令，用于对文件进行文本替换操作。不支持流编辑的复杂功能，专注于最常用的替换场景。

## 二、命令规格

### 命令格式

```bash
fck sed [options] <file>
```

### 核心功能

1. **文本替换**: 将匹配的内容替换为指定文本
2. **正则支持**: 可选使用正则表达式匹配
3. **全局替换**: 默认替换所有匹配项
4. **行号范围**: 可选限定替换的行号范围
5. **原地修改**: 默认输出到标准输出，可选直接修改文件

## 三、标志设计

| 长标志 | 短标志 | 类型 | 必填 | 默认值 | 说明 |
|--------|--------|------|------|--------|------|
| `--pattern` | `-p` | string | ✅ | - | 匹配模式（固定字符串或正则） |
| `--replacement` | `-r` | string | ✅ | - | 替换内容 |
| `--regexp` | `-E` | bool | ❌ | false | 使用正则表达式匹配 |
| `--line` | `-n` | string | ❌ | - | 指定行号范围（如: 1,5 或 10） |
| `--in-place` | `-i` | bool | ❌ | false | 原地修改文件 |
| `--backup` | `-b` | bool | ❌ | false | 修改前创建 .bak 备份 |
| `--count` | `-c` | int | ❌ | 0 | 最大替换次数（0表示无限制） |
| `--ignore-case` | `-I` | bool | ❌ | false | 忽略大小写 |

## 四、使用示例

### 基础替换

```bash
# 将所有 "old" 替换为 "new"
fck sed -p "old" -r "new" file.txt

# 使用正则表达式
fck sed -E -p "foo\d+" -r "bar" file.txt

# 忽略大小写
fck sed -I -p "hello" -r "world" file.txt
```

### 行号范围

```bash
# 只替换第 5 行
fck sed -p "old" -r "new" -n 5 file.txt

# 替换 1-10 行
fck sed -p "old" -r "new" -n 1,10 file.txt

# 替换第 5 行及之后
fck sed -p "old" -r "new" -n 5, file.txt
```

### 原地修改

```bash
# 直接修改文件（输出到文件）
fck sed -p "old" -r "new" -i file.txt

# 修改前创建备份 file.txt.bak
fck sed -p "old" -r "new" -i -b file.txt
```

### 限制替换次数

```bash
# 只替换前 3 处匹配
fck sed -p "old" -r "new" -c 3 file.txt
```

## 五、目录结构

```
internal/
├── cli/
│   └── sed.go              # CLI 定义
└── commands/
    └── sed/
        └── cmd_sed.go      # 业务逻辑
```

## 六、配置结构体

```go
// SedConfig sed 命令配置
type SedConfig struct {
    // CLI 参数
    Target      string // 目标文件路径
    Pattern     string // -p 匹配模式
    Replacement string // -r 替换内容
    Regexp      bool   // -E 使用正则
    LineRange   string // -n 行号范围
    InPlace     bool   // -i 原地修改
    Backup      bool   // -b 创建备份
    MaxCount    int    // -c 最大替换次数
    IgnoreCase  bool   // -I 忽略大小写

    // 运行时
    compiledPattern *regexp.Regexp // 编译后的正则
    lineStart       int            // 起始行号
    lineEnd         int            // 结束行号（0表示到末尾）
}
```

## 七、核心算法

### 行号范围解析

```go
// 解析行号范围字符串
// "5" -> start=5, end=5
// "1,10" -> start=1, end=10
// "5," -> start=5, end=0(末尾)
func parseLineRange(rangeStr string) (start, end int, err error)
```

### 替换逻辑

```go
// 处理单行替换
func processLine(line string, config *SedConfig, lineNum int) (string, bool) {
    // 1. 检查行号范围
    if !inRange(lineNum, config.lineStart, config.lineEnd) {
        return line, false
    }

    // 2. 执行替换
    if config.Regexp {
        return replaceWithRegex(line, config)
    }
    return replaceWithString(line, config)
}
```

### 文件处理流程

```go
func processFile(config *SedConfig) error {
    // 1. 读取文件
    // 2. 逐行处理
    // 3. 如果 -i，写入原文件（可选备份）
    // 4. 否则输出到 stdout
}
```

## 八、错误处理

| 错误场景 | 错误信息 |
|----------|----------|
| 未指定文件 | "未指定目标文件" |
| 文件不存在 | "文件不存在: <path>" |
| 未指定 pattern | "请使用 -p 指定匹配模式" |
| 未指定 replacement | "请使用 -r 指定替换内容" |
| 行号格式错误 | "行号范围格式错误，示例: 5 或 1,10" |
| 正则编译失败 | "正则表达式无效: <error>" |
| 权限不足 | "无法写入文件: <path>" |

## 九、安全考虑

1. **备份机制**: `-i` 模式下建议配合 `-b` 使用，防止误操作
2. **原子写入**: 先写入临时文件，再重命名，避免中途失败导致文件损坏
3. **权限检查**: 修改前检查文件写权限
4. **二进制文件检测**: 简单检测是否为文本文件，避免破坏二进制文件

## 十、实现检查清单

- [ ] `internal/cli/sed.go` - CLI 定义
- [ ] `internal/commands/sed/cmd_sed.go` - 业务逻辑
- [ ] 行号范围解析功能
- [ ] 正则/字符串替换功能
- [ ] 原地修改功能（含备份）
- [ ] 错误处理
- [ ] 帮助文档示例
- [ ] 编译验证

---

> **方案设计完成，待批准后实施**
