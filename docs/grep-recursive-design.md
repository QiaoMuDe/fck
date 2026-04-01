# grep 递归搜索功能设计方案

## 功能概述
为 `fck grep` 命令添加递归目录搜索功能，支持文件/目录过滤和符号链接控制。

## 新增标志

| 标志 | 类型 | 说明 |
|------|------|------|
| `-r, --recursive` | bool | 递归搜索目录（不跟随符号链接） |
| `-R` | bool | 递归搜索目录（跟随符号链接） |
| `--include` | string | 只搜索匹配模式的文件，多个模式用逗号分隔（如 `*.go,*.js`） |
| `--exclude` | string | 排除匹配模式的文件，多个模式用逗号分隔 |
| `--include-dir` | string | 只进入匹配模式的目录，多个模式用逗号分隔 |
| `--exclude-dir` | string | 排除匹配模式的目录，多个模式用逗号分隔 |

## 使用示例

```bash
# 递归搜索当前目录
fck grep -r "pattern" ./

# 递归并跟随符号链接
fck grep -R "pattern" ./

# 只搜索 .go 和 .js 文件
fck grep -r --include="*.go,*.js" "pattern" ./

# 排除 test 文件
fck grep -r --exclude="*_test.go" "pattern" ./

# 排除 node_modules 目录
fck grep -r --exclude-dir="node_modules" "pattern" ./

# 组合使用
fck grep -r --include="*.go" --exclude-dir="vendor,node_modules" "func main" ./
```

## 技术实现

### 1. 配置结构扩展

```go
type GrepConfig struct {
    // 现有字段...
    
    // 递归搜索相关
    Recursive     bool     // -r 递归搜索
    FollowSymlink bool     // -R 跟随符号链接
    Include       []string // --include 包含文件模式（逗号分隔解析为切片）
    Exclude       []string // --exclude 排除文件模式（逗号分隔解析为切片）
    IncludeDir    []string // --include-dir 包含目录模式（逗号分隔解析为切片）
    ExcludeDir    []string // --exclude-dir 排除目录模式（逗号分隔解析为切片）
}
```

### 2. 核心遍历逻辑

```go
func processPath(path string, config *GrepConfig) error {
    fileInfo, err := os.Stat(path)
    if err != nil {
        return err
    }
    
    // 单文件模式
    if !config.Recursive && !config.FollowSymlink {
        if fileInfo.IsDir() {
            return fmt.Errorf("is a directory, use -r or -R for recursive search: %s", path)
        }
        return processSingleFile(path, config)
    }
    
    // 递归模式
    if !fileInfo.IsDir() {
        return processSingleFile(path, config)
    }
    
    return filepath.WalkDir(path, func(filePath string, d fs.DirEntry, err error) error {
        if err != nil {
            // 权限错误等，跳过
            return nil
        }
        
        // 处理符号链接
        if d.Type()&fs.ModeSymlink != 0 {
            if !config.FollowSymlink {
                return nil // 跳过符号链接
            }
            // 解析符号链接...
        }
        
        // 目录过滤
        if d.IsDir() {
            dirName := d.Name()
            if !matchDirFilters(dirName, config) {
                return filepath.SkipDir
            }
            return nil
        }
        
        // 文件过滤
        fileName := d.Name()
        if !matchFileFilters(fileName, config) {
            return nil
        }
        
        // 处理文件
        return processSingleFile(filePath, config)
    })
}
```

### 3. 过滤函数

```go
// matchFileFilters 检查文件是否匹配过滤条件
func matchFileFilters(fileName string, config *GrepConfig) bool {
    // 检查 exclude
    for _, pattern := range config.Exclude {
        if matched, _ := filepath.Match(pattern, fileName); matched {
            return false
        }
    }
    
    // 检查 include（如果有指定，必须匹配其中一个）
    if len(config.Include) > 0 {
        for _, pattern := range config.Include {
            if matched, _ := filepath.Match(pattern, fileName); matched {
                return true
            }
        }
        return false
    }
    
    return true
}

// matchDirFilters 检查目录是否匹配过滤条件
func matchDirFilters(dirName string, config *GrepConfig) bool {
    // 检查 exclude-dir
    for _, pattern := range config.ExcludeDir {
        if matched, _ := filepath.Match(pattern, dirName); matched {
            return false
        }
    }
    
    // 检查 include-dir（如果有指定，必须匹配其中一个）
    if len(config.IncludeDir) > 0 {
        for _, pattern := range config.IncludeDir {
            if matched, _ := filepath.Match(pattern, dirName); matched {
                return true
            }
        }
        return false
    }
    
    return true
}
```

### 4. 输出格式调整

递归模式下，文件名默认显示（相当于自动启用 `-H`），格式：
```
path/to/file.go:10:matched content
```

配合 `--group` 标志（后续可添加）可实现：
```
path/to/file.go
10:matched content
20:another match

another/file.go
5:matched content
```

## 边界情况处理

1. **符号链接循环**：使用 `map[string]bool` 记录已访问的真实路径
2. **权限错误**：跳过无权限访问的文件/目录，继续遍历
3. **二进制文件**：默认跳过（检查文件内容是否包含 `\x00`）
4. **-r 和 -R 同时指定**：`-R` 优先级更高（或报错）
5. **无匹配文件**：正常退出，返回码 0（无匹配）或 1（有匹配）

## 改动范围

- `internal/cli/grep.go`：新增 6 个标志定义
- `internal/commands/grep/cmd_grep.go`：
  - 扩展 `GrepConfig` 结构
  - 新增 `processPath` 函数
  - 新增过滤相关函数
  - 修改 `GrepCmdMain` 入口逻辑

预估改动：约 150-200 行代码
