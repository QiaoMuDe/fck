# Chown 命令设计方案

## 功能概述
实现一个用于修改文件和目录所有者的命令 `chown`，支持递归修改、用户/组设置、详细输出等功能，提供便捷的所有权管理操作。

## 核心功能

### 1. 修改文件所有者
- 支持修改单个文件所有者
- 支持修改多个文件所有者
- 支持设置用户和组

### 2. 修改目录所有者
- 支持修改单个目录所有者
- 支持修改多个目录所有者
- 支持递归修改目录及其内容的所有者

### 3. 递归修改
- 支持递归修改目录及其所有子目录和文件
- 自动遍历目录树
- 支持跳过符号链接

### 4. 用户/组设置
- 支持设置用户（如 `user`）
- 支持设置用户和组（如 `user:group`）
- 支持只设置组（如 `:group`）
- 支持使用 UID/GID（如 `1000:1000`）

### 5. 详细输出
- 支持显示修改的文件/目录
- 支持显示所有者变化
- 支持显示操作统计信息

### 6. 错误处理
- 检测权限是否足够
- 检测路径是否存在
- 检测用户/组是否存在
- 提供友好的错误提示

## 命令选项设计

| 选项 | 短选项 | 说明 | 默认值 |
|------|--------|------|--------|
| `--recursive` | `-R` | 递归修改目录及其内容 | false |
| `--verbose` | `-v` | 显示修改的文件/目录 | false |
| `--changes` | `-c` | 只显示有变化的文件/目录 | false |
| `--quiet` | `-f` | 抑制错误信息 | false |

## 使用示例

### 修改单个文件所有者
```bash
fck chown user file.txt
# 将 file.txt 所有者设置为 user
```

### 修改文件所有者和组
```bash
fck chown user:group file.txt
# 将 file.txt 所有者设置为 user，组设置为 group
```

### 只修改组
```bash
fck chown :group file.txt
# 将 file.txt 组设置为 group（不修改所有者）
```

### 使用 UID/GID
```bash
fck chown 1000:1000 file.txt
# 将 file.txt 所有者 UID 设置为 1000，组 GID 设置为 1000
```

### 修改多个文件所有者
```bash
fck chown user file1.txt file2.txt file3.txt
# 将多个文件所有者设置为 user
```

### 修改目录所有者
```bash
fck chown user directory/
# 将 directory/ 所有者设置为 user
```

### 递归修改所有者
```bash
fck chown -R user directory/
# 递归修改 directory/ 及其所有内容的所有者
```

### 详细输出
```bash
fck chown -v user file.txt
# 输出: 修改所有者: file.txt (olduser:oldgroup -> user:group)
```

### 只显示有变化的文件
```bash
fck chown -c user file.txt
# 只在所有者实际改变时显示
```

### 抑制错误信息
```bash
fck chown -f user nonexistent.txt
# 不显示错误信息
```

### 组合使用
```bash
fck chown -Rv user:group directory/
# 递归修改并显示详情
```

## 文件结构

```
internal/
├── commands/
│   └── chown/
│       └── cmd_chown.go          # 业务逻辑实现
└── cli/
    └── chown.go                   # CLI定义和配置
```

## 实现细节

### 1. 配置结构体

```go
type ChownConfig struct {
    Owner     string
    Targets   []string
    Recursive bool
    Verbose   bool
    Changes   bool
    Quiet     bool
}
```

### 2. 所有者解析

支持的所有者格式：
- `user`：只设置用户
- `user:group`：设置用户和组
- `:group`：只设置组
- `UID:GID`：使用数字 ID

解析逻辑：
```go
// 解析所有者字符串
parts := strings.Split(ownerStr, ":")
if len(parts) == 1 {
    // 只设置用户
    user = parts[0]
    group = ""
} else if len(parts) == 2 {
    // 设置用户和组
    user = parts[0]
    group = parts[1]
}
```

### 3. 核心处理流程

```
1. 解析命令行参数，构建 ChownConfig
2. 解析所有者字符串：
   - 验证所有者格式
   - 解析用户和组
3. 处理每个目标文件/目录：
   - 检查路径是否存在
   - 如果是目录且 Recursive=true，递归处理
   - 如果是文件或目录且 Recursive=false，直接修改
4. 输出结果：
   - 如果 Verbose=true，显示修改的文件
   - 如果 Changes=true，只显示有变化的文件
   - 显示操作统计信息
```

### 4. 所有者修改逻辑

```
1. 获取当前所有者
2. 解析目标所有者
3. 比较所有者是否相同
4. 如果不同：
   - 修改所有者
   - 如果 Verbose=true，输出修改信息
   - 如果 Changes=true，输出修改信息
5. 如果相同：
   - 如果 Changes=true，跳过
   - 如果 Verbose=true，输出跳过信息
```

### 5. 递归处理逻辑

```
1. 遍历目录及其所有内容
2. 对每个文件/目录：
   - 跳过符号链接（可选）
   - 修改所有者
   - 如果是目录，继续递归
3. 处理错误：
   - 如果 Quiet=true，抑制错误
   - 如果 Quiet=false，显示错误
```

### 6. 跨平台考虑

使用标准库 `os.Chown` 函数，该函数会自动处理平台差异：
- Windows 平台：可能有限制，但会提供友好的错误提示
- Unix-like 平台（Linux、macOS）：支持完整的 UID/GID 功能

### 7. 用户/组解析

```go
import (
    "os/user"
    "strconv"
)

// 解析用户
func parseUser(userStr string) (int, error) {
    // 尝试解析为 UID
    if uid, err := strconv.Atoi(userStr); err == nil {
        return uid, nil
    }

    // 尝试查找用户名
    u, err := user.Lookup(userStr)
    if err != nil {
        return -1, err
    }

    return strconv.Atoi(u.Uid)
}

// 解析组
func parseGroup(groupStr string) (int, error) {
    // 尝试解析为 GID
    if gid, err := strconv.Atoi(groupStr); err == nil {
        return gid, nil
    }

    // 尝试查找组名
    g, err := user.LookupGroup(groupStr)
    if err != nil {
        return -1, err
    }

    return strconv.Atoi(g.Gid)
}
```

## 边界情况处理

1. **空目标列表**：提示错误，不执行任何操作
2. **所有者字符串无效**：提示错误，不执行操作
3. **文件不存在**：报错，提示文件不存在
4. **权限不足**：报错，提示权限问题
5. **用户/组不存在**：报错，提示用户/组不存在
6. **符号链接**：可选择跳过或跟随链接
7. **递归循环**：检测并避免循环
8. **所有者相同**：跳过修改，或显示提示

## 安全考虑

### 1. 权限验证
- 验证用户是否有修改所有者的权限
- 防止未授权的所有权修改
- 验证目标所有者的有效性

### 2. 路径验证
- 检查路径是否合法
- 验证路径长度限制
- 防止路径遍历攻击

### 3. 递归安全
- 检测循环链接
- 限制递归深度
- 跟随符号链接（可选）

### 4. 跨平台安全
- Windows 平台的特殊处理
- Unix-like 平台的权限检查
- 提供平台特定的错误提示

## 测试用例

### 单元测试
- 测试修改单个文件所有者
- 测试修改多个文件所有者
- 测试修改目录所有者
- 测试递归修改所有者
- 测试所有者解析
- 测试用户/组解析
- 测试详细输出
- 测试只显示变化

### 集成测试
- 测试命令行参数解析
- 测试选项组合
- 测试错误处理
- 测试边界情况

### 跨平台测试
- 测试 Windows 平台
- 测试 Linux 平台
- 测试 macOS 平台

### 安全测试
- 测试权限验证
- 测试路径验证
- 测试递归安全

## 注意事项

1. 与系统 `chown` 命令的兼容性：保持基本用法一致
2. 跨平台支持：Windows 和 Linux/macOS 的所有权处理
3. 所有者格式：支持用户名、组名、UID、GID
4. 错误处理：提供友好的错误提示
5. 性能考虑：大量文件操作时的性能优化
6. 符号链接：默认跳过符号链接
7. Windows 限制：Windows 平台可能有功能限制

## 平台差异说明

### Windows 平台
- 使用 SID 而非 UID/GID
- 需要 `golang.org/x/sys/windows` 包
- 可能需要管理员权限
- 功能可能受限

### Unix-like 平台
- 使用 UID/GID
- 支持用户名和组名
- 需要 root 权限修改其他用户的文件

## 后续扩展（可选）

1. 支持参考文件所有者（如 --reference file.txt）
2. 支持通配符模式（如 *.txt）
3. 支持排除模式（如 --exclude）
4. 支持批量操作
5. 支持所有权继承（如 --preserve-root）
6. 支持从文件读取所有者列表
