# Cp 命令设计方案

## 功能概述
实现一个用于复制文件和目录的命令 `cp`，支持递归复制、保留文件属性、覆盖控制等功能，提供便捷的文件复制操作。

## 核心功能

### 1. 复制文件
- 支持复制单个文件
- 支持复制多个文件到目录
- 支持文件重命名

### 2. 复制目录
- 支持复制单个目录
- 支持递归复制目录及其内容（自动处理）
- 支持复制多个目录

### 3. 保留文件属性
- 自动保留文件权限（由 fs.CopyEx 处理）
- 自动保留文件时间戳（由 fs.CopyEx 处理）
- 自动递归复制目录（由 fs.CopyEx 处理）

### 4. 覆盖控制
- 支持强制覆盖
- 支持交互式覆盖

### 5. 详细输出
- 支持显示复制的文件/目录
- 支持显示操作统计信息

### 6. 错误处理
- 检测源文件是否存在
- 检测目标目录是否存在
- 检测权限是否足够
- 提供友好的错误提示

## 命令选项设计

| 选项 | 短选项 | 说明 | 默认值 |
|------|--------|------|--------|
| `--force` | `-f` | 强制覆盖已存在的文件 | false |
| `--interactive` | `-i` | 交互式覆盖（覆盖前提示） | false |
| `--verbose` | `-v` | 显示复制的文件/目录 | false |

## 使用示例

### 复制单个文件
```bash
fck cp file1.txt file2.txt
# 将 file1.txt 复制为 file2.txt
```

### 复制多个文件到目录
```bash
fck cp file1.txt file2.txt file3.txt directory/
# 将多个文件复制到 directory/ 目录
```

### 递归复制目录
```bash
fck cp source_dir/ target_dir/
# 递归复制 source_dir/ 到 target_dir/（自动处理）
```

### 保留文件属性
```bash
fck cp file.txt copy.txt
# 复制文件并自动保留权限和时间戳
```

### 强制覆盖
```bash
fck cp -f file.txt existing.txt
# 强制覆盖已存在的 existing.txt
```

### 交互式覆盖
```bash
fck cp -i file.txt existing.txt
# 覆盖前提示确认
```

### 详细输出
```bash
fck cp -v file1.txt file2.txt
# 输出: 复制: file1.txt -> file2.txt
```

### 组合使用
```bash
fck cp -v source_dir/ target_dir/
# 递归复制并显示详情
```

## 文件结构

```
internal/
├── commands/
│   └── cp/
│       └── cmd_cp.go              # 业务逻辑实现
└── cli/
    └── cp.go                       # CLI定义和配置
```

## 实现细节

### 1. 配置结构体

```go
type CpConfig struct {
    Force       bool
    Interactive bool
    Verbose     bool
    Sources     []string
    Target      string
}
```

### 2. 核心处理流程

```
1. 解析命令行参数，构建 CpConfig
2. 验证参数：
   - 检查源文件是否存在
   - 检查目标路径是否有效
3. 确定复制类型：
   - 单文件到单文件
   - 多文件到目录
   - 目录到目录
4. 执行复制：
   - 使用 fs.CopyEx 进行复制（自动处理递归和属性保留）
   - 处理覆盖（根据选项）
5. 输出结果：
   - 如果 Verbose=true，显示复制的文件
   - 显示操作统计信息
```

### 3. 复制文件

使用 `gitee.com/MM-Q/go-kit/fs` 包的 `CopyEx` 函数进行复制，该函数已经实现了所有逻辑，包括：
- 文件复制
- 目录复制（递归）
- 覆盖控制
- 属性保留

```go
import (
    "fmt"
    "gitee.com/MM-Q/go-kit/fs"
)

// 复制文件或目录
func copyItem(src, dst string, force, verbose bool) error {
    // 使用 fs.CopyEx 进行复制（自动处理递归和属性保留）
    if err := fs.CopyEx(src, dst, force); err != nil {
        return fmt.Errorf("复制 '%s' 到 '%s' 失败: %v", src, dst, err)
    }

    if verbose {
        fmt.Printf("复制: %s -> %s\n", src, dst)
    }

    return nil
}
```

### 4. 递归复制目录

使用 `fs.CopyEx` 函数自动处理递归复制，无需手动实现。

### 5. 交互式覆盖

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
5. **磁盘空间不足**：报错，提示磁盘空间不足
6. **符号链接**：可选择跳过或跟随链接
7. **递归循环**：检测并避免循环
8. **同名文件**：根据选项处理覆盖

## 安全考虑

### 1. 权限验证
- 验证用户是否有读取源文件的权限
- 验证用户是否有写入目标目录的权限
- 验证用户是否有修改文件属性的权限

### 2. 路径验证
- 检查路径是否合法
- 验证路径长度限制
- 防止路径遍历攻击

### 3. 递归安全
- 检测循环链接
- 限制递归深度
- 跟随符号链接（可选）

### 4. 覆盖安全
- 防止意外覆盖重要文件
- 提供交互式确认
- 支持不覆盖选项

## 测试用例

### 单元测试
- 测试复制单个文件
- 测试复制多个文件
- 测试递归复制目录
- 测试保留文件属性
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
- 测试递归安全
- 测试覆盖安全

## 注意事项

1. 与系统 `cp` 命令的兼容性：保持基本用法一致
2. 跨平台支持：Windows 和 Linux/macOS 的文件操作
3. 文件属性：不同平台的文件属性支持
4. 错误处理：提供友好的错误提示
5. 性能考虑：大文件操作时的性能优化
6. 磁盘空间：检查磁盘空间是否足够
7. 符号链接：默认跳过符号链接

## 平台差异说明

### Windows 平台
- 路径分隔符：`\`（反斜杠）
- 驱动器盘符：如 `C:\`
- 文件属性：权限支持有限
- 符号链接：支持符号链接（需要管理员权限）

### Unix-like 平台
- 路径分隔符：`/`（正斜杠）
- 根目录：`/`
- 文件属性：完整的权限和时间戳支持
- 符号链接：广泛使用

## 标准库函数

使用以下标准库函数：
- `io.Copy()`：复制文件内容
- `os.Open()`：打开文件
- `os.Create()`：创建文件
- `os.Stat()`：获取文件信息
- `os.MkdirAll()`：创建目录
- `os.ReadDir()`：读取目录
- `os.Chmod()`：修改文件权限
- `os.Chtimes()`：修改文件时间戳
- `filepath.Join()`：拼接路径
- `filepath.Walk()`：遍历目录树

## 后续扩展（可选）

1. 支持通配符模式（如 *.txt）
2. 支持排除模式（如 --exclude）
3. 支持进度显示（如 --progress）
4. 支持校验和验证（如 --checksum）
5. 支持压缩传输（如 --compress）
6. 支持增量复制（如 --incremental）
7. 支持从文件读取复制列表
