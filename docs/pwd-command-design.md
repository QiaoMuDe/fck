# Pwd 命令设计方案

## 功能概述
实现一个用于打印当前工作目录的命令 `pwd`，支持显示物理路径和逻辑路径，提供便捷的目录查看功能。

## 核心功能

### 1. 显示当前工作目录
- 显示当前工作目录的绝对路径
- 支持显示逻辑路径（不解析符号链接）
- 支持显示物理路径（解析符号链接）

### 2. 路径类型
- **逻辑路径**：不解析符号链接，显示当前工作目录的实际路径（默认）
- **物理路径**：解析所有符号链接，显示最终的物理路径

### 3. 错误处理
- 检测目录是否存在
- 检测权限是否足够
- 提供友好的错误提示

## 命令选项设计

| 选项 | 短选项 | 说明 | 默认值 |
|------|--------|------|--------|
| `--physical` | `-P` | 显示物理路径（解析符号链接） | false |
| `--logical` | `-L` | 显示逻辑路径（不解析符号链接） | true |

## 使用示例

### 显示当前工作目录（逻辑路径）
```bash
fck pwd
# 输出: /home/user/project
```

### 显示物理路径（解析符号链接）
```bash
fck pwd -P
# 输出: /mnt/data/projects/project
```

### 显示逻辑路径（显式指定）
```bash
fck pwd -L
# 输出: /home/user/project
```

### 符号链接示例
```bash
# 假设 /home/user/project 是一个符号链接，指向 /mnt/data/projects/project

# 显示逻辑路径（不解析符号链接）
fck pwd
# 输出: /home/user/project

# 显示物理路径（解析符号链接）
fck pwd -P
# 输出: /mnt/data/projects/project
```

## 文件结构

```
internal/
├── commands/
│   └── pwd/
│       └── cmd_pwd.go           # 业务逻辑实现
└── cli/
    └── pwd.go                    # CLI定义和配置
```

## 实现细节

### 1. 配置结构体

```go
type PwdConfig struct {
    Physical bool
    Logical  bool
}
```

### 2. 核心处理流程

```
1. 解析命令行参数，构建 PwdConfig
2. 获取当前工作目录：
   - 使用 os.Getwd() 获取当前工作目录
3. 处理路径类型：
   - 如果 Physical=true，解析符号链接
   - 如果 Logical=true（默认），不解析符号链接
4. 输出路径：
   - 打印路径到标准输出
5. 错误处理：
   - 处理获取目录失败的情况
   - 处理解析符号链接失败的情况
```

### 3. 获取当前工作目录

```go
import "os"

// 获取当前工作目录
func getWorkingDirectory() (string, error) {
    return os.Getwd()
}
```

### 4. 解析符号链接

```go
import "path/filepath"

// 解析符号链接，获取物理路径
func resolveSymlinks(path string) (string, error) {
    return filepath.EvalSymlinks(path)
}
```

### 5. 路径处理逻辑

```go
// 处理路径
func processPath(config PwdConfig) (string, error) {
    // 获取当前工作目录
    cwd, err := getWorkingDirectory()
    if err != nil {
        return "", fmt.Errorf("获取当前工作目录失败: %w", err)
    }

    // 如果需要物理路径，解析符号链接
    if config.Physical {
        resolved, err := resolveSymlinks(cwd)
        if err != nil {
            return "", fmt.Errorf("解析符号链接失败: %w", err)
        }
        return resolved, nil
    }

    // 默认返回逻辑路径
    return cwd, nil
}
```

### 6. 主函数

```go
// PwdCmdMain 主函数
func PwdCmdMain(config PwdConfig) error {
    path, err := processPath(config)
    if err != nil {
        return err
    }

    fmt.Println(path)
    return nil
}
```

## 边界情况处理

1. **目录不存在**：报错，提示目录不存在
2. **权限不足**：报错，提示权限问题
3. **符号链接循环**：报错，提示符号链接循环
4. **符号链接断裂**：报错，提示符号链接断裂
5. **路径过长**：报错，提示路径过长

## 安全考虑

### 1. 路径验证
- 验证路径是否合法
- 验证路径长度限制
- 防止路径遍历攻击

### 2. 权限检查
- 检查是否有读取目录的权限
- 检查是否有解析符号链接的权限

### 3. 符号链接安全
- 检测符号链接循环
- 限制符号链接解析深度
- 处理符号链接断裂情况

## 测试用例

### 单元测试
- 测试获取当前工作目录
- 测试显示逻辑路径
- 测试显示物理路径
- 测试解析符号链接
- 测试错误处理

### 集成测试
- 测试命令行参数解析
- 测试选项组合
- 测试边界情况

### 安全测试
- 测试路径验证
- 测试权限检查
- 测试符号链接安全

## 注意事项

1. 与系统 `pwd` 命令的兼容性：保持基本用法一致
2. 跨平台支持：Windows 和 Linux/macOS 的路径处理
3. 路径格式：使用平台特定的路径分隔符
4. 错误处理：提供友好的错误提示
5. 符号链接：正确处理符号链接
6. 性能考虑：快速响应，无性能问题

## 平台差异说明

### Windows 平台
- 路径分隔符：`\`（反斜杠）
- 驱动器盘符：如 `C:\`
- 符号链接：支持符号链接（需要管理员权限）

### Unix-like 平台
- 路径分隔符：`/`（正斜杠）
- 根目录：`/`
- 符号链接：广泛使用

## 标准库函数

使用以下标准库函数：
- `os.Getwd()`：获取当前工作目录
- `filepath.EvalSymlinks()`：解析符号链接
- `filepath.Clean()`：清理路径
- `fmt.Println()`：输出路径

## 后续扩展（可选）

1. 支持显示相对路径（如 --relative）
2. 支持显示路径的详细信息（如 --verbose）
3. 支持自定义输出格式（如 --format）
4. 支持颜色输出（如 --color）
5. 支持路径规范化（如 --normalize）
