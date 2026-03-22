# Truncate 命令设计方案

## 功能概述
实现一个用于截断文件到指定大小的命令 `truncate`，支持创建新文件、扩展文件、缩小文件等功能，提供便捷的文件大小管理操作。

## 核心功能

### 1. 截断文件大小
- 支持将文件截断到指定大小
- 支持扩展文件（增大文件）
- 支持缩小文件（减小文件）
- 支持创建新文件

### 2. 大小单位支持
- 支持字节（B）
- 支持千字节（K、KB）
- 支持兆字节（M、MB）
- 支持吉字节（G、GB）
- 支持太字节（T、TB）

### 3. 文件操作
- 支持单个文件操作
- 支持多个文件操作
- 支持文件不存在时创建
- 支持参考其他文件大小

### 4. 错误处理
- 检测文件是否存在
- 检测权限是否足够
- 检测大小格式是否有效
- 提供友好的错误提示

## 命令选项设计

| 选项 | 短选项 | 说明 | 默认值 |
|------|--------|------|--------|
| `--size` | `-s` | 设置文件大小（支持单位：B、K、M、G、T） | 必需 |
| `--create` | `-c` | 如果文件不存在则创建 | false |
| `--reference` | `-r` | 参考文件大小 | "" |
| `--verbose` | `-v` | 显示操作的文件 | false |

## 使用示例

### 截断文件到指定大小
```bash
fck truncate -s 1M file.txt
# 将 file.txt 截断到 1MB 大小
```

### 创建指定大小的文件
```bash
fck truncate -s 100M newfile.bin
# 创建一个 100MB 的新文件
```

### 缩小文件
```bash
fck truncate -s 10K largefile.txt
# 将 largefile.txt 缩小到 10KB
```

### 扩展文件
```bash
fck truncate -s 1G smallfile.txt
# 将 smallfile.txt 扩展到 1GB
```

### 参考其他文件大小
```bash
fck truncate -r reference.txt target.txt
# 将 target.txt 设置为与 reference.txt 相同的大小
```

### 创建多个文件
```bash
fck truncate -s 1M file1.txt file2.txt file3.txt
# 创建三个 1MB 的文件
```

### 详细输出
```bash
fck truncate -v -s 1M file.txt
# 输出: 截断文件: file.txt (0 -> 1048576)
```

### 组合使用
```bash
fck truncate -cv -s 1M newfile.txt
# 创建文件并显示详情
```

## 文件结构

```
internal/
├── commands/
│   └── truncate/
│       └── cmd_truncate.go      # 业务逻辑实现
└── cli/
    └── truncate.go              # CLI定义和配置
```

## 实现细节

### 1. 配置结构体

```go
type TruncateConfig struct {
    Size      string
    Reference string
    Targets   []string
    Create    bool
    Verbose   bool
}
```

### 2. 核心处理流程

```
1. 解析命令行参数，构建 TruncateConfig
2. 解析大小：
   - 如果指定了 --reference，读取参考文件大小
   - 否则解析 --size 参数
3. 处理每个目标文件：
   - 检查文件是否存在
   - 如果不存在且 Create=true，创建文件
   - 截断文件到指定大小
4. 输出结果：
   - 如果 Verbose=true，显示操作的文件
   - 显示操作统计信息
```

### 3. 大小解析

使用 qflag 库的 SizeFlag，支持的大小格式：
- 数字：如 `100`（字节）
- 带单位：如 `100K`、`100M`、`1G`、`1T`
- 大小写单位：如 `KB`、`MB`、`GB`、`TB`

SizeFlag 自动解析为 int64，无需手动解析。

### 4. 截断文件

```go
import "os"

// 截断文件
func truncateFile(path string, size int64, create bool, verbose bool) error {
    // 检查文件是否存在
    if _, err := os.Stat(path); os.IsNotExist(err) {
        if !create {
            return fmt.Errorf("文件不存在: %s", path)
        }
        // 创建文件
        file, err := os.Create(path)
        if err != nil {
            return fmt.Errorf("创建文件失败: %w", err)
        }
        file.Close()
    }

    // 获取当前文件大小
    info, err := os.Stat(path)
    if err != nil {
        return fmt.Errorf("获取文件信息失败: %w", err)
    }
    currentSize := info.Size()

    // 截断文件
    if err := os.Truncate(path, size); err != nil {
        return fmt.Errorf("截断文件失败: %w", err)
    }

    if verbose {
        fmt.Printf("截断文件: %s (%d -> %d)\n", path, currentSize, size)
    }

    return nil
}
```

### 5. 参考文件大小

```go
// 获取参考文件大小
func getReferenceSize(refPath string) (int64, error) {
    info, err := os.Stat(refPath)
    if err != nil {
        return 0, fmt.Errorf("获取参考文件信息失败: %w", err)
    }
    return info.Size(), nil
}
```

## 边界情况处理

1. **空目标列表**：提示错误，不执行任何操作
2. **大小格式无效**：提示错误，不执行操作
3. **文件不存在且未指定创建**：报错，提示文件不存在
4. **权限不足**：报错，提示权限问题
5. **参考文件不存在**：报错，提示参考文件不存在
6. **大小为负数**：报错，提示大小不能为负数
7. **磁盘空间不足**：报错，提示磁盘空间不足

## 安全考虑

### 1. 权限验证
- 验证用户是否有读取/写入文件的权限
- 验证用户是否有创建文件的权限
- 验证用户是否有读取参考文件的权限

### 2. 路径验证
- 检查路径是否合法
- 验证路径长度限制
- 防止路径遍历攻击

### 3. 大小验证
- 验证大小格式是否有效
- 验证大小是否在合理范围内
- 防止设置过大的文件大小

### 4. 磁盘空间检查
- 检查磁盘空间是否足够
- 防止磁盘空间耗尽

## 测试用例

### 单元测试
- 测试截断文件
- 测试创建文件
- 测试扩展文件
- 测试缩小文件
- 测试大小解析
- 测试参考文件大小
- 测试详细输出

### 集成测试
- 测试命令行参数解析
- 测试选项组合
- 测试错误处理
- 测试边界情况

### 安全测试
- 测试权限验证
- 测试路径验证
- 测试大小验证
- 测试磁盘空间检查

## 注意事项

1. 与系统 `truncate` 命令的兼容性：保持基本用法一致
2. 跨平台支持：Windows 和 Linux/macOS 的文件操作
3. 大小单位：支持多种单位（B、K、M、G、T）
4. 错误处理：提供友好的错误提示
5. 性能考虑：大文件操作时的性能优化
6. 磁盘空间：检查磁盘空间是否足够

## 大小单位说明

| 单位 | 说明 | 字节数 |
|------|------|--------|
| B | 字节 | 1 |
| K/KB | 千字节 | 1024 |
| M/MB | 兆字节 | 1024*1024 |
| G/GB | 吉字节 | 1024*1024*1024 |
| T/TB | 太字节 | 1024*1024*1024*1024 |

## 标准库函数

使用以下标准库函数：
- `os.Truncate()`：截断文件
- `os.Create()`：创建文件
- `os.Stat()`：获取文件信息
- `os.OpenFile()`：打开文件
- `fmt.Println()`：输出信息

## 后续扩展（可选）

1. 支持百分比大小（如 --size 50%）
2. 支持相对大小（如 +10M、-10M）
3. 支持排除模式（如 --exclude）
4. 支持批量操作
5. 支持文件权限设置（如 --mode）
6. 支持稀疏文件（如 --sparse）
7. 支持从文件读取大小列表
