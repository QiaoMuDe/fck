# grep 命令改造设计方案

> **需求**: 新增支持默认跳过隐藏文件/目录，通过标志启用才处理；默认读取 .gitignore 文件加载排除规则  
> **设计日期**: 2026-04-02  
> **实施日期**: 2026-04-02  
> **状态**: ✅ 已完成

---

## 一、需求变更摘要

| 需求项 | 当前行为 | 期望行为 |
|--------|----------|----------|
| **隐藏文件/目录** | 无特殊处理 | 默认跳过，通过标志启用才处理 |
| **.gitignore 支持** | 无 | 默认读取 `.gitignore` 加载排除规则 |

---

## 二、新增 CLI 标志设计

在 `internal/cli/grep.go` 中添加以下标志：

```go
// 新增标志变量
var (
    // 隐藏文件控制
    grepHidden *qflag.BoolFlag  // -H, --hidden 处理隐藏文件/目录（默认跳过）
    
    // gitignore 控制
    grepNoGitignore *qflag.BoolFlag     // --no-gitignore 禁用 .gitignore 读取
    grepGitignorePath *qflag.StringFlag // --gitignore-path 指定自定义 .gitignore 路径
)
```

### 标志详细设计

| 标志 | 短选项 | 默认值 | 说明 |
|------|--------|--------|------|
| `--hidden` | `-H` | `false` | 启用处理隐藏文件/目录（默认跳过） |
| `--no-gitignore` | - | `false` | 禁用自动读取 `.gitignore` 文件 |
| `--gitignore-path` | - | `""` | 指定自定义 `.gitignore` 文件路径 |

---

## 三、配置结构体扩展

在 `internal/commands/grep/cmd_grep.go` 的 `GrepConfig` 中添加：

```go
type GrepConfig struct {
    // ... 现有字段 ...

    // 新增：隐藏文件控制
    Hidden bool // -H 处理隐藏文件/目录（默认跳过）
    
    // 新增：gitignore 控制
    NoGitignore     bool     // --no-gitignore 禁用 .gitignore
    GitignorePath   string   // --gitignore-path 自定义 .gitignore 路径
    
    // 内部使用：解析后的 gitignore 规则
    gitignoreRules  []gitignoreRule // 解析后的排除规则
}
```

---

## 四、.gitignore 解析规则设计

### 4.1 支持的标准规则

```go
// gitignoreRule 表示一条 gitignore 规则
type gitignoreRule struct {
    pattern    string      // 原始模式
    regex      *regexp.Regexp // 编译后的正则（用于复杂模式）
    isNegation bool        // 是否否定规则（以!开头）
    isDirOnly  bool        // 是否仅目录（以/结尾）
    matchPath  bool        // 是否匹配完整路径（包含/）
}
```

### 4.2 支持的语法

| 模式 | 含义 | 示例 |
|------|------|------|
| `*.log` | 忽略所有 `.log` 文件 | `*.log` |
| `temp/` | 忽略 `temp` 目录 | `node_modules/` |
| `/build` | 仅忽略根目录的 `build` | `/build` |
| `doc/*.txt` | 忽略 `doc` 下的 `.txt` | `docs/*.md` |
| `!important.log` | 不忽略 `important.log` | `!.gitkeep` |
| `# comment` | 注释，忽略 | `# 这是注释` |
| ` ` | 空行，忽略 |  |

### 4.3 查找策略

```
1. 如果 --no-gitignore，跳过 gitignore 加载
2. 如果指定了 --gitignore-path，使用该路径
3. 否则，使用当前目录下的 .gitignore 文件（如果存在）
4. 如果 .gitignore 不存在，静默跳过，不报错
```

---

## 五、核心逻辑改造点

### 5.1 隐藏文件判断

```go
// isHidden 判断文件/目录是否为隐藏
func isHidden(path string) bool {
    base := filepath.Base(path)
    // Unix/Linux: 以 . 开头
    // Windows: 使用文件属性检查
    if runtime.GOOS == "windows" {
        // 使用 Windows API 检查 FILE_ATTRIBUTE_HIDDEN
        // 或简单判断以 . 开头（大多数工具行为）
    }
    return strings.HasPrefix(base, ".")
}
```

### 5.2 递归遍历改造（recursive.go）

在 `processPathRecursive` 函数中增加过滤逻辑：

```go
func processPathRecursive(path string, config *GrepConfig) error {
    // ... 现有代码 ...
    
    // 加载 gitignore 规则（如果启用）
    if !config.NoGitignore {
        if err := loadGitignoreRules(config); err != nil && !config.NoMessages {
            // 可选：输出警告
        }
    }
    
    return filepath.WalkDir(path, func(filePath string, d fs.DirEntry, err error) error {
        // ... 现有错误处理 ...
        
        // 检查隐藏文件/目录
        if !config.Hidden && isHidden(d.Name()) {
            if d.IsDir() {
                return filepath.SkipDir
            }
            return nil
        }
        
        // 检查 gitignore 规则
        if shouldIgnoreByGitignore(filePath, d.IsDir(), config) {
            if d.IsDir() {
                return filepath.SkipDir
            }
            return nil
        }
        
        // ... 后续处理 ...
    })
}
```

### 5.3 单文件模式改造

在 `processSingleFileWithError` 中增加检查：

```go
func processSingleFileWithError(path string, config *GrepConfig) error {
    // 检查是否为隐藏文件
    if !config.Hidden && isHidden(path) {
        if config.NoMessages {
            return nil
        }
        return fmt.Errorf("hidden file skipped (use -H to include): %s", path)
    }
    
    // ... 现有代码 ...
}
```

---

## 六、文件变更清单

| 文件 | 变更类型 | 变更内容 |
|------|----------|----------|
| `internal/cli/grep.go` | 修改 | 新增 3 个标志定义和配置传递 |
| `internal/commands/grep/cmd_grep.go` | 修改 | GrepConfig 新增 3 个字段 |
| `internal/commands/grep/recursive.go` | 修改 | 遍历逻辑增加隐藏文件和 gitignore 过滤 |
| `internal/commands/grep/gitignore.go` | **新增** | gitignore 解析和匹配逻辑 |

---

## 七、新增文件设计：gitignore.go

```go
// internal/commands/grep/gitignore.go
package grep

import (
    "bufio"
    "os"
    "path/filepath"
    "regexp"
    "strings"
)

// gitignoreRule 表示一条解析后的 gitignore 规则
type gitignoreRule struct {
    pattern    string
    regex      *regexp.Regexp
    isNegation bool
    isDirOnly  bool
    matchPath  bool
}

// loadGitignoreRules 加载 .gitignore 规则
func loadGitignoreRules(config *GrepConfig) error {
    gitignorePath := config.GitignorePath
    if gitignorePath == "" {
        // 从起始目录向上查找
        var err error
        gitignorePath, err = findGitignore(config.Target)
        if err != nil {
            return err // 未找到，不报错
        }
    }
    
    rules, err := parseGitignore(gitignorePath)
    if err != nil {
        return err
    }
    
    config.gitignoreRules = rules
    return nil
}

// findGitignore 从指定路径向上查找 .gitignore
func findGitignore(startPath string) (string, error) {
    // 实现向上查找逻辑
}

// parseGitignore 解析 .gitignore 文件
func parseGitignore(path string) ([]gitignoreRule, error) {
    // 实现解析逻辑
}

// shouldIgnoreByGitignore 检查路径是否应被忽略
func shouldIgnoreByGitignore(path string, isDir bool, config *GrepConfig) bool {
    // 实现匹配逻辑（支持否定规则）
}
```

---

## 八、使用示例

```bash
# 默认行为：跳过隐藏文件和 .gitignore 中指定的文件
fck grep "pattern" ./myproject

# 包含隐藏文件
fck grep -H "pattern" ./myproject

# 禁用 .gitignore
fck grep --no-gitignore "pattern" ./myproject

# 使用自定义 .gitignore
fck grep --gitignore-path ./custom-ignore "pattern" ./myproject

# 组合使用
fck grep -H -n --no-gitignore "TODO" ./src
```

---

## 九、边缘案例考虑

| 场景 | 处理方案 |
|------|----------|
| .gitignore 不存在 | 静默跳过，不报错 |
| .gitignore 格式错误 | 跳过错误行，继续解析 |
| 否定规则冲突 | 按 git 标准：后匹配的规则优先 |
| 符号链接指向隐藏文件 | 跟随 -R 标志行为，再判断隐藏属性 |
| 多层 .gitignore | 暂不支持，只读取最近的一层 |

---

## 十、与现有命令的一致性

参考 `hash` 命令的隐藏文件处理方式：
- `hash` 使用 `-H, --hidden` 启用处理隐藏文件
- `grep` 将采用相同的标志命名和行为模式

---

## 十一、实施检查清单

- [x] `internal/cli/grep.go` - 添加 3 个新标志
- [x] `internal/commands/grep/cmd_grep.go` - 扩展 GrepConfig 结构体
- [x] `internal/commands/grep/gitignore.go` - 新建 gitignore 解析模块
- [x] `internal/commands/grep/recursive.go` - 修改递归遍历逻辑
- [x] 编译验证通过
- [ ] 更新命令帮助文档中的示例（可选）
- [ ] 测试各种边缘场景（建议补充单元测试）

---

> **✅ 方案实施完成**
