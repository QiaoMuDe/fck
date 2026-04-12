# fuzzy - 模糊字符串匹配库

> 该子包源自开源项目：https://github.com/sahilm/fuzzy.git

## 📖 简介

`fuzzy` 是一个高性能的模糊字符串匹配库，专为文件名和代码符号搜索优化，匹配风格类似于 Sublime Text、VSCode、IntelliJ IDEA 等主流编辑器的智能搜索功能。

## ✨ 核心特性

- **智能评分系统** - 基于奖励和惩罚机制计算匹配质量
- **多种匹配奖励** - 首字符匹配、驼峰命名匹配、分隔符后匹配、相邻字符匹配
- **Unicode 支持** - 完整支持多语言字符匹配（大小写不敏感）
- **高性能** - 零外部依赖，ASCII 快速路径优化
- **灵活接口** - 支持自定义数据源（Source 接口）

## 🔧 安装

```bash
go get gitee.com/MM-Q/go-kit/fuzzy
```

## 🚀 快速开始

### 基础用法

```go
package main

import (
    "fmt"
    "gitee.com/MM-Q/go-kit/fuzzy"
)

func main() {
    // 文件名列表
    files := []string{
        "main.go",
        "Makefile",
        "my_app.go",
        "application.go",
        "user_controller.go",
    }
    
    // 搜索 "mg" 模式
    matches := fuzzy.Find("mg", files)
    
    // 输出匹配结果（按匹配质量排序）
    for _, match := range matches {
        fmt.Printf("Score: %d, Match: %s\n", match.Score, match.Str)
    }
}
```

**输出示例：**
```
Score: 10, Match: main.go
Score: 5, Match: Makefile
```

### 高亮匹配字符

```go
matches := fuzzy.Find("tk", []string{"The Black Knight", "test_key"})

for _, match := range matches {
    fmt.Printf("匹配: %s (分数: %d)\n", match.Str, match.Score)
    fmt.Print("高亮: ")
    for i := 0; i < len(match.Str); i++ {
        if contains(i, match.MatchedIndexes) {
            fmt.Printf("[%c]", match.Str[i]) // 高亮显示
        } else {
            fmt.Printf("%c", match.Str[i])
        }
    }
    fmt.Println()
}
```

### 使用自定义数据源

```go
// 实现 Source 接口
type MySource struct {
    items []MyItem
}

func (s MySource) String(i int) string {
    return s.items[i].Name
}

func (s MySource) Len() int {
    return len(s.items)
}

// 使用自定义数据源
source := MySource{items: myItems}
matches := fuzzy.FindFrom("pattern", source)
```

### 不排序的搜索（性能更好）

```go
// 如果不需要排序结果，使用 FindNoSort 性能更好
matches := fuzzy.FindNoSort("pattern", data)
```

## 📊 评分规则

### 奖励机制

| 匹配类型 | 奖励分数 | 说明 |
|---------|---------|------|
| 首字符匹配 | +10 | 模式第一个字符匹配字符串开头 |
| 驼峰命名匹配 | +20 | 小写字母后接大写字母（如 `aB`） |
| 分隔符后匹配 | +20 | 匹配 `/` `-` `_` ` ` `.` `\` 后的字符 |
| 相邻字符匹配 | +5（递增） | 连续匹配的字符获得递增奖励 |

### 惩罚机制

| 惩罚类型 | 惩罚分数 | 说明 |
|---------|---------|------|
| 前导未匹配字符 | -5/字符（最大-15） | 第一个匹配前的字符 |
| 未匹配字符 | -1/字符 | 字符串中未匹配的字符 |

### 匹配示例

| 模式 | 字符串 | 分数 | 说明 |
|-----|--------|-----|------|
| `tk` | The Black K**n**ight | 较高 | 匹配第二个 `k`，相邻奖励 |
| `uc` | **u**ser_**c**ontroller | 高 | 分隔符后匹配 + 首字符匹配 |
| `ma` | **M**akefile | 高 | 首字符匹配 |
| `go` | main.**go** | 较低 | 扩展名匹配，无特殊奖励 |

## 📚 API 文档

### 函数

| 函数 | 签名 | 说明 |
|-----|------|------|
| `Find` | `Find(pattern string, data []string) Matches` | 在字符串切片中搜索，返回排序结果 |
| `FindNoSort` | `FindNoSort(pattern string, data []string) Matches` | 在字符串切片中搜索，返回未排序结果（更快） |
| `FindFrom` | `FindFrom(pattern string, data Source) Matches` | 在自定义数据源中搜索，返回排序结果 |
| `FindFromNoSort` | `FindFromNoSort(pattern string, data Source) Matches` | 在自定义数据源中搜索，返回未排序结果 |

### 类型

```go
// Match 表示一个匹配结果
type Match struct {
    Str            string // 匹配到的字符串
    Index          int    // 在原切片中的索引
    MatchedIndexes []int  // 匹配字符的位置索引（用于高亮）
    Score          int    // 匹配分数（越高越匹配）
}

// Matches 是 Match 的切片，实现 sort.Interface
type Matches []Match

// Source 数据源接口
type Source interface {
    String(i int) string // 获取第 i 个字符串
    Len() int            // 数据源长度
}
```

## 💡 使用场景

- **命令行工具** - 文件名/命令快速搜索（如 fzf、peco）
- **IDE 插件** - 代码符号、文件快速定位
- **Web 应用** - 搜索框智能提示
- **文件管理器** - 快速文件查找

## 🏗️ 项目信息

- **许可证**: MIT
- **Go 版本要求**: 1.16+
- **外部依赖**: 无（标准库实现）

## 📝 源码说明

本包基于 [sahilm/fuzzy](https://github.com/sahilm/fuzzy) 开源实现，针对中文注释和代码文档进行了完善。

---

**使用建议**: 对于需要高性能搜索的场景（如 10K+ 文件），建议使用 `FindNoSort` 并自行缓存结果。
