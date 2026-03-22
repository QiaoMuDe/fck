# Home 命令设计方案

## 功能概述
实现一个用于显示用户主目录的命令 `home`，支持显示不同平台的主目录路径，提供便捷的主目录查看功能。

## 核心功能

### 1. 显示用户主目录
- 显示当前用户的主目录路径
- 支持显示绝对路径
- 支持跨平台（Windows/Linux/macOS）

### 2. 跨平台支持
- **Windows 平台**：使用 `USERPROFILE` 环境变量
- **Unix-like 平台（Linux/macOS）**：使用 `HOME` 环境变量
- **Go 标准库**：使用 `os.UserHomeDir()` 函数

### 3. 错误处理
- 检测主目录是否存在
- 检测权限是否足够
- 提供友好的错误提示

## 命令选项设计

该命令不需要任何选项，直接显示用户主目录路径。

## 使用示例

### 显示用户主目录
```bash
fck home
# Windows 输出: C:\Users\username
# Linux/macOS 输出: /home/username
```

### 在脚本中使用
```bash
# 获取主目录并切换到该目录
cd $(fck home)

# 在主目录中创建文件
touch $(fck home)/test.txt
```

## 文件结构

```
internal/
├── commands/
│   └── home/
│       └── cmd_home.go          # 业务逻辑实现
└── cli/
    └── home.go                   # CLI定义和配置
```

## 实现细节

### 1. 配置结构体

```go
type HomeConfig struct {
    // 该命令不需要任何配置
}
```

### 2. 核心处理流程

```
1. 获取用户主目录：
   - 使用 os.UserHomeDir() 函数
2. 输出路径：
   - 打印路径到标准输出
3. 错误处理：
   - 处理获取主目录失败的情况
```

### 3. 获取用户主目录

```go
import "os"

// 获取用户主目录
func getUserHomeDir() (string, error) {
    return os.UserHomeDir()
}
```

### 4. 主函数

```go
// HomeCmdMain 主函数
func HomeCmdMain(config HomeConfig) error {
    homeDir, err := getUserHomeDir()
    if err != nil {
        return fmt.Errorf("获取用户主目录失败: %w", err)
    }

    fmt.Println(homeDir)
    return nil
}
```

## 边界情况处理

1. **主目录不存在**：报错，提示主目录不存在
2. **权限不足**：报错，提示权限问题
3. **环境变量未设置**：报错，提示环境变量未设置
4. **路径过长**：报错，提示路径过长

## 安全考虑

### 1. 路径验证
- 验证路径是否合法
- 验证路径长度限制
- 防止路径遍历攻击

### 2. 权限检查
- 检查是否有读取主目录的权限
- 检查是否有访问环境变量的权限

## 测试用例

### 单元测试
- 测试获取用户主目录
- 测试错误处理
- 测试跨平台支持

### 集成测试
- 测试命令行参数解析
- 测试边界情况

### 跨平台测试
- 测试 Windows 平台
- 测试 Linux 平台
- 测试 macOS 平台

### 安全测试
- 测试路径验证
- 测试权限检查

## 注意事项

1. 与系统 `echo $HOME` 或 `echo ~` 的兼容性：保持基本用法一致
2. 跨平台支持：Windows 和 Linux/macOS 的路径处理
3. 路径格式：使用平台特定的路径分隔符
4. 错误处理：提供友好的错误提示
5. 性能考虑：快速响应，无性能问题

## 平台差异说明

### Windows 平台
- 路径分隔符：`\`（反斜杠）
- 驱动器盘符：如 `C:\`
- 主目录：`C:\Users\username`
- 环境变量：`USERPROFILE`

### Unix-like 平台
- 路径分隔符：`/`（正斜杠）
- 根目录：`/`
- 主目录：`/home/username`（Linux）或 `/Users/username`（macOS）
- 环境变量：`HOME`

## 标准库函数

使用以下标准库函数：
- `os.UserHomeDir()`：获取用户主目录
- `fmt.Println()`：输出路径

## 与其他命令的对比

| 命令 | 功能 | 区别 |
|------|------|------|
| `fck home` | 显示用户主目录 | 专门用于显示主目录 |
| `fck pwd` | 显示当前工作目录 | 显示当前所在目录 |
| `echo $HOME` | 显示主目录（Unix） | 需要环境变量支持 |
| `echo ~` | 显示主目录（Unix） | Shell 扩展功能 |

## 使用场景

### 1. 脚本中使用
```bash
# 在脚本中获取主目录
HOME_DIR=$(fck home)
cd $HOME_DIR
```

### 2. 配置文件中使用
```bash
# 在配置文件中引用主目录
config_path=$(fck home)/.config/app/config.json
```

### 3. 跨平台脚本
```bash
# 跨平台脚本中使用主目录
HOME_DIR=$(fck home)
echo "主目录: $HOME_DIR"
```

### 4. 快速导航
```bash
# 快速切换到主目录
cd $(fck home)
```

## 后续扩展（可选）

1. 支持显示其他用户的主目录（如 `fck home username`）
2. 支持显示相对路径（如 --relative）
3. 支持显示路径的详细信息（如 --verbose）
4. 支持自定义输出格式（如 --format）
5. 支持颜色输出（如 --color）
6. 支持路径规范化（如 --normalize）
7. 支持检查主目录是否存在（如 --check）
