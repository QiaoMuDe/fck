# fuzzy API 文档

```go
package fuzzy // import "gitee.com/MM-Q/go-kit/fuzzy"
```

Package fuzzy 提供了模糊字符串匹配功能，该功能针对文件名和代码符号进行了优化，风格类似于 Sublime Text、VSCode、IntelliJ IDEA 等工具。

fuzzy 包定义了模糊匹配的类型。

## 类型

### Match

```go
type Match struct {
	Str            string // 匹配到的字符串
	Index          int    // 匹配字符串在提供的切片中的索引
	MatchedIndexes []int  // 匹配字符的索引。用于高亮显示匹配项
	Score          int    // 用于对匹配结果进行排序的分数
}
```

Match 表示一个匹配到的字符串。

### Matches

```go
type Matches []Match
```

Matches 是 Match 结构体的切片

### Source

```go
type Source interface {
	String(i int) string // 位置 i 处要匹配的字符串
	Len() int            // 源的长度。通常是你想要匹配的事物切片的长度
}
```

Source 表示字符串列表的抽象源。Source 必须是可迭代类型，如切片。源将被迭代直到 Len()，对每个元素调用 String(i)，其中 i 是元素的索引。你可以在 README 中找到工作示例。

## 函数

### Find

```go
func Find(pattern string, data []string) Matches
```

Find 在 data 中查找 pattern 并返回匹配结果，按匹配质量降序排列。匹配质量由一组奖励和惩罚规则决定。

以下类型的匹配会获得奖励：
- 模式中的第一个字符与匹配字符串的第一个字符匹配
- 匹配的字符是驼峰命名
- 匹配的字符位于分隔符（如下划线字符）之后
- 匹配的字符与之前的匹配相邻

对于搜索字符串中每个未匹配的字符以及第一个匹配之前的所有前导字符，都会应用惩罚。

**参数:**
- `pattern`: 要查找的模式字符串
- `data`: 要搜索的字符串切片

**返回值:**
- `Matches`: 按匹配质量降序排列的匹配结果

### FindFrom

```go
func FindFrom(pattern string, data Source) Matches
```

FindFrom 使用 Source 接口在数据源中查找 pattern 并返回排序后的匹配结果。该函数与 Find 功能相同，但接受 Source 接口而不是字符串切片，适用于自定义数据源。结果按匹配质量降序排列。

**参数:**
- `pattern`: 要查找的模式字符串
- `data`: 实现 Source 接口的数据源

**返回值:**
- `Matches`: 按匹配质量降序排列的匹配结果

### FindFromNoSort

```go
func FindFromNoSort(pattern string, data Source) Matches
```

FindFromNoSort 使用 Source 接口在数据源中查找 pattern，不排序。该函数是模糊匹配的核心实现，负责计算匹配分数和匹配位置。匹配质量由奖励和惩罚规则决定：首字符匹配、驼峰命名匹配、分隔符后匹配、相邻匹配可获得奖励；未匹配字符和前导字符会受到惩罚。

**参数:**
- `pattern`: 要查找的模式字符串
- `data`: 实现 Source 接口的数据源

**返回值:**
- `Matches`: 匹配结果（未排序）

### FindNoSort

```go
func FindNoSort(pattern string, data []string) Matches
```

FindNoSort 在 data 中查找 pattern 并返回匹配结果，不排序。该函数与 Find 功能相同，但返回的结果不会按匹配质量排序，性能略好。

**参数:**
- `pattern`: 要查找的模式字符串
- `data`: 要搜索的字符串切片

**返回值:**
- `Matches`: 匹配结果（未排序）

### Complete

```go
func Complete(pattern string, candidates []string) Matches
```

Complete 在 candidates 中搜索 pattern，优先前缀匹配，其次模糊匹配。返回按匹配质量排序的结果，分数越高表示匹配度越好。

适用于命令行标志补全场景，匹配优先级如下：
1. 精确匹配 (1000分)
2. 前缀匹配 (100-200分，越短的候选分数越高)
3. 模糊匹配 (0-99分，复用 fuzzy 算法但降低权重)

**参数:**
- `pattern`: 要查找的模式字符串
- `candidates`: 候选字符串切片

**返回值:**
- `Matches`: 按匹配质量降序排列的匹配结果

**示例:**

```go
flags := []string{"--verbose", "--version", "-v"}
matches := fuzzy.Complete("--v", flags)
// matches[0].Str = "--verbose"
// matches[1].Str = "--version"
// matches[2].Str = "-v"
```

### CompletePrefix

```go
func CompletePrefix(pattern string, candidates []string) Matches
```

CompletePrefix 仅返回以 pattern 为前缀的候选。不区分大小写，返回结果按候选长度升序排列（短的优先）。

**参数:**
- `pattern`: 前缀模式字符串
- `candidates`: 候选字符串切片

**返回值:**
- `Matches`: 前缀匹配的结果列表

**示例:**

```go
flags := []string{"--verbose", "--version", "-v"}
matches := fuzzy.CompletePrefix("--v", flags)
// 结果: [--verbose, --version] (不包含 -v)
```

### CompleteExact

```go
func CompleteExact(pattern string, candidates []string) Matches
```

CompleteExact 仅返回与 pattern 完全相等的候选。区分大小写。

**参数:**
- `pattern`: 精确匹配字符串
- `candidates`: 候选字符串切片

**返回值:**
- `Matches`: 精确匹配的结果列表（最多一个）

**示例:**

```go
flags := []string{"--verbose", "--version", "-v"}
matches := fuzzy.CompleteExact("-v", flags)
// 结果: [-v]
```

## Matches 方法

### Len

```go
func (a Matches) Len() int
```

Len 返回匹配结果的长度。

**返回值:**
- `int`: 匹配结果的长度

### Less

```go
func (a Matches) Less(i, j int) bool
```

Less 返回第一个匹配结果是否应该排在第二个匹配结果之前。按分数降序排列，分数相同时按原始索引升序排列（稳定排序）。

**参数:**
- `i`: 第一个匹配结果的索引
- `j`: 第二个匹配结果的索引

**返回值:**
- `bool`: 如果第一个匹配结果应该排在第二个之前，则返回 true

### Swap

```go
func (a Matches) Swap(i, j int)
```

Swap 交换两个匹配结果的顺序。

**参数:**
- `i`: 第一个匹配结果的索引
- `j`: 第二个匹配结果的索引
