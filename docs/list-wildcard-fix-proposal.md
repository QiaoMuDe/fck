# List 命令通配符问题修复方案

## 问题描述

当使用通配符查看路径时（如 `fck list .\vendor\git*`），当前实现会直接进入匹配到的目录内部并分组显示内容，而期望的行为是只显示匹配到的目录本身（作为条目列出）。

### 当前行为
```bash
$ fck list .\vendor\git*

vendor\gitee.com:
  MM-Q/  ...

vendor\github.com:
  jedib0t/  mattn/  mitchellh/  rivo/  schollz/
```

### 期望行为
```bash
$ fck list .\vendor\git*

  d  rwxrwxrwx  0  B  2026-03-24 13:49:22   gitee.com
  d  rwxrwxrwx  0  B  2026-03-24 13:49:22   github.com
```

## 问题根源分析

问题位于 `internal/commands/list/scanner.go` 文件的 `ScanWithOriginalPaths` 函数（第 68-86 行）：

```go
// 检查是否为通配符展开的目录
isWildcardDir := strings.ContainsAny(originalPath, "*?[]")
if isWildcardDir {
    if info, err := os.Stat(expandedPath); err == nil && info.IsDir() {
        // 通配符展开的目录：扫描目录内容，但保持原始路径为目录路径
        files, err := s.scanDirectoryWithOriginal(expandedPath, expandedPath, expandedPath, opts)
        if err != nil {
            return nil, fmt.Errorf("扫描目录 %s 失败: %v", expandedPath, err)
        }
        allFiles = append(allFiles, files...)
        continue
    }
}
```

**问题**：当检测到通配符匹配的目录时，代码直接调用 `scanDirectoryWithOriginal` 扫描目录**内部内容**，而不是将目录本身作为条目返回。

## 修复方案

### 方案一：移除通配符特殊处理（推荐）

**核心思路**：通配符匹配的目录应该与普通目录一样处理，只显示目录条目本身，不自动进入内部。

**修改文件**：`internal/commands/list/scanner.go`

**修改内容**：
删除 `ScanWithOriginalPaths` 函数中的通配符特殊处理逻辑（第 68-86 行），让通配符匹配的目录走正常的 `scanSinglePathWithOriginal` 流程。

```go
func (s *FileScanner) ScanWithOriginalPaths(originalPaths, expandedPaths []string, opts ScanOptions) (FileInfoList, error) {
    var allFiles FileInfoList

    // 创建原始路径到展开路径的映射
    pathMapping := s.createPathMapping(originalPaths, expandedPaths)

    for _, expandedPath := range expandedPaths {
        // 找到对应的原始路径
        originalPath := s.findOriginalPath(expandedPath, pathMapping)

        // 删除以下通配符特殊处理代码块：
        // isWildcardDir := strings.ContainsAny(originalPath, "*?[]")
        // if isWildcardDir {
        //     if info, err := os.Stat(expandedPath); err == nil && info.IsDir() {
        //         files, err := s.scanDirectoryWithOriginal(expandedPath, expandedPath, expandedPath, opts)
        //         ...
        //     }
        // }

        // 所有路径统一使用 scanSinglePathWithOriginal 处理
        files, err := s.scanSinglePathWithOriginal(expandedPath, originalPath, opts)
        if err != nil {
            return nil, fmt.Errorf("扫描路径 %s 失败: %v", expandedPath, err)
        }
        allFiles = append(allFiles, files...)
    }

    return allFiles, nil
}
```

**效果**：
- 通配符匹配到目录时，只显示目录条目
- 通配符匹配到文件时，显示文件条目
- 使用 `-r` 递归选项时，正常递归显示目录内容

### 方案二：添加条件判断

**核心思路**：只在递归模式下才扫描通配符目录的内部内容。

**修改内容**：
```go
// 检查是否为通配符展开的目录
isWildcardDir := strings.ContainsAny(originalPath, "*?[]")
if isWildcardDir && opts.Recursive {  // 添加递归条件判断
    if info, err := os.Stat(expandedPath); err == nil && info.IsDir() {
        files, err := s.scanDirectoryWithOriginal(expandedPath, expandedPath, expandedPath, opts)
        ...
    }
}
```

**对比**：
| 方案 | 优点 | 缺点 |
|------|------|------|
| 方案一 | 逻辑简单，行为一致 | 改变现有行为 |
| 方案二 | 保留递归时进入目录的能力 | 逻辑稍复杂 |

## 推荐方案

**推荐方案一**，原因：
1. 符合 Unix/Linux `ls` 命令的行为习惯
2. 逻辑更简单，易于维护
3. 用户可以通过 `-r` 选项明确指定是否需要递归查看

## 边缘案例考虑

1. **通配符匹配到混合类型**（文件+目录）：应该分别显示文件条目和目录条目
2. **通配符未匹配到任何内容**：保持现有警告提示
3. **使用 `-d` 选项列出目录本身**：与方案一行为一致
4. **使用 `-r` 选项递归**：正常递归显示所有匹配目录的内容

## 测试用例建议

```bash
# 测试1：通配符匹配目录，非递归模式
fck list .\vendor\git*
# 期望：显示 gitee.com 和 github.com 两个目录条目

# 测试2：通配符匹配目录，递归模式
fck list -r .\vendor\git*
# 期望：显示两个目录及其内部内容

# 测试3：通配符匹配文件
fck list *.go
# 期望：显示匹配到的 Go 文件条目

# 测试4：混合路径（通配符+具体路径）
fck list .\vendor\git* .\docs\
# 期望：分别显示匹配的目录和 docs 目录的内容（如果 docs 不是通配符匹配的话）

# 测试5：多层通配符
fck list .\vendor\*\*
# 期望：显示匹配到的所有条目
```
