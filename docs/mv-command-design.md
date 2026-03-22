# Mv 命令设计方案

## 功能概述
实现一个用于移动文件和目录的命令 `mv`，支持重命名文件、移动文件到目录、移动目录等功能，提供便捷的文件移动操作。

## 核心功能

### 1. 移动文件
- 支持移动单个文件
- 支持移动多个文件到目录
- 支持文件重命名

### 2. 移动目录
- 支持移动单个目录
- 支持移动多个目录

### 3. 覆盖控制
- 支持强制覆盖
- 支持交互式覆盖

### 4. 详细输出
- 支持显示移动的文件/目录
- 支持显示操作统计信息

### 5. 错误处理
- 检测源文件是否存在
- 检测目标目录是否存在
- 检测权限是否足够
- 提供友好的错误提示

## 命令选项设计

| 选项 | 短选项 | 说明 | 默认值 |
|------|--------|------|--------|
| `--force` | `-f` | 强制覆盖已存在的文件 | false |
| `--interactive` | `-i` | 交互式覆盖（覆盖前提示） | false |
| `--verbose` | `-v` | 显示移动的文件/目录 | false |

## 使用示例

### 移动单个文件（重命名）
```bash
fck mv file1.txt file2.txt
# 将 file1.txt 重命名为 file2.txt
```

### 移动多个文件到目录
```bash
fck mv file1.txt file2.txt file3.txt directory/
# 将多个文件移动到 directory/ 目录
```

### 移动目录
```bash
fck mv source_dir/ target_dir/
# 移动 source_dir/ 到 target_dir/
```

### 强制覆盖
```bash
fck mv -f file.txt existing.txt
# 强制覆盖已存在的 existing.txt
```

### 交互式覆盖
```bash
fck mv -i file.txt existing.txt
# 覆盖前提示确认
```

### 详细输出
```bash
fck mv -v file1.txt file2.txt
# 输出: 移动: file1.txt -> file2.txt
```

### 组合使用
```bash
fck mv -v source_dir/ target_dir/
# 移动目录并显示详情
```

## 文件结构

```
internal/
├── commands/
│   └── mv/
│       └── cmd_mv.go              # 业务逻辑实现
└── cli/
    └── mv.go                       # CLI定义和配置
```

## 实现细节

### 1. 配置结构体

```go
type MvConfig struct {
    Force       bool
    Interactive bool
    Verbose     bool
    Sources     []string
    Target      string
}
```

### 2. 核心处理流程

```
1. 解析命令行参数，构建 MvConfig
2. 验证参数：
   - 检查源文件是否存在
   - 检查目标路径是否有效
3. 确定移动类型：
   - 单文件到单文件（重命名）
   - 多文件到目录
   - 目录到目录
4. 执行移动：
   - 使用 fs.MoveEx 进行移动
   - 处理覆盖（根据选项）
5. 输出结果：
   - 如果 Verbose=true，显示移动的文件
   - 显示操作统计信息
```

### 3. 移动文件或目录

使用 `gitee.com/MM-Q/go-kit/fs` 包的 `MoveEx` 函数进行移动，该函数已经实现了所有逻辑，包括：
- 文件移动
- 目录移动
- 覆盖控制

```go
import (
    "fmt"
    "gitee.com/MM-Q/go-kit/fs"
)

// 移动文件或目录
func moveItem(src, dst string, force, verbose bool) error {
    // 使用 fs.MoveEx 进行移动
    if err := fs.MoveEx(src, dst, force); err != nil {
        return fmt.Errorf("移动 '%s' 到 '%s' 失败: %v", src, dst, err)
    }

    if verbose {
        fmt.Printf("移动: %s -> %s\n", src, dst)
    }

    return nil
}
```

### 4. 交互式覆盖

```go
import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

// 交互式覆盖确认
func confirmOverwrite(dst string) bool {
    reader := bufio.NewReader(os.Stdin)
    fmt.Printf("覆盖 %s? (y/n): ", dst)

    response, err := reader.ReadString('\n')
    if err != nil {
        return false
    }

    response = strings.TrimSpace(strings.ToLower(response))
    return response == "y" || response == "yes"
}
```

## 边界情况处理

1. **空源列表**：提示错误，不执行任何操作
2. **源文件不存在**：报错，提示源文件不存在
3. **目标目录不存在**：报错，提示目标目录不存在
4. **权限不足**：报错，提示权限问题
5. **同名文件**：根据选项处理覆盖
6. **跨设备移动**：由 fs.MoveEx 自动处理

## 安全考虑

### 1. 权限验证
- 验证用户是否有读取源文件的权限
- 验证用户是否有写入目标目录的权限

### 2. 路径验证
- 检查路径是否合法
- 验证路径长度限制
- 防止路径遍历攻击

### 3. 覆盖安全
- 防止意外覆盖重要文件
- 提供交互式确认
- 支持不覆盖选项

## 测试用例

### 单元测试
- 测试移动单个文件
- 测试移动多个文件
- 测试移动目录
- 测试文件重命名
- 测试覆盖控制
- 测试详细输出

### 集成测试
- 测试命令行参数解析
- 测试选项组合
- 测试错误处理
- 测试边界情况

### 安全测试
- 测试权限验证
- 测试路径验证
- 测试覆盖安全

## 注意事项

1. 与系统 `mv` 命令的兼容性：保持基本用法一致
2. 跨平台支持：Windows 和 Linux/macOS 的文件操作
3. 错误处理：提供友好的错误提示
4. 性能考虑：大文件操作时的性能优化
5. 原子性：移动操作应该是原子的（由 fs.MoveEx 处理）

## 平台差异说明

### Windows 平台
- 路径分隔符：`\`（反斜杠）
- 驱动器盘符：如 `C:\`
- 跨驱动器移动：需要特殊处理（由 fs.MoveEx 处理）

### Unix-like 平台
- 路径分隔符：`/`（正斜杠）
- 根目录：`/`
- 文件系统：支持跨文件系统移动（由 fs.MoveEx 处理）

## 标准库函数

使用以下标准库和第三方库函数：
- `fs.MoveEx()`：移动文件或目录（来自 `gitee.com/MM-Q/go-kit/fs`）
- `os.Stat()`：获取文件信息
- `os.PathSeparator`：路径分隔符
- `fmt.Println()`：输出信息

## 后续扩展（可选）

1. 支持通配符模式（如 *.txt）
2. 支持排除模式（如 --exclude）
3. 支持进度显示（如 --progress）
4. 支持备份原文件（如 --backup）
5. 支持从文件读取移动列表
