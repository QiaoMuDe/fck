# Test 命令设计方案

## 功能概述
实现一个用于检测指定路径是否存在的命令 `test`，支持检测文件、目录、符号链接是否存在，提供便捷的路径检测功能。

## 核心功能

### 1. 路径存在性检测
- 支持检测路径是否存在
- 存在则正常退出（退出码 0）
- 不存在则报错退出（退出码 1）

### 2. 文件类型检测
- 支持检测文件是否存在（`-f`）
- 支持检测目录是否存在（`-d`）
- 支持检测符号链接是否存在（`-l`）

### 3. 错误处理
- 提供友好的错误提示
- 明确指出路径不存在或类型不匹配

## 命令选项设计

| 选项 | 短选项 | 说明 | 默认值 |
|------|--------|------|--------|
| `--file` | `-f` | 检查指定文件是否存在 | false |
| `--dir` | `-d` | 检查指定目录是否存在 | false |
| `--link` | `-l` | 检查指定符号链接是否存在 | false |

## 使用示例

### 检测路径是否存在
```bash
fck test /path/to/file.txt
# 如果存在：正常退出（退出码 0）
# 如果不存在：报错退出（退出码 1）
```

### 检测文件是否存在
```bash
fck test -f /path/to/file.txt
# 如果是文件：正常退出
# 如果不是文件或不存在：报错退出
```

### 检测目录是否存在
```bash
fck test -d /path/to/directory
# 如果是目录：正常退出
# 如果不是目录或不存在：报错退出
```

### 检测符号链接是否存在
```bash
fck test -l /path/to/symlink
# 如果是符号链接：正常退出
# 如果不是符号链接或不存在：报错退出
```

### 在脚本中使用
```bash
# 检查文件是否存在
if fck test -f /path/to/file.txt; then
    echo "文件存在"
else
    echo "文件不存在"
fi

# 检查目录是否存在
if fck test -d /path/to/directory; then
    echo "目录存在"
else
    echo "目录不存在"
fi
```

## 文件结构

```
internal/
├── commands/
│   └── test/
│       └── cmd_test.go            # 业务逻辑实现
└── cli/
    └── test.go                    # CLI定义和配置
```

## 实现细节

### 1. 配置结构体

```go
type TestConfig struct {
    CheckFile    bool
    CheckDir    bool
    CheckLink   bool
    Path        string
}
```

### 2. 核心处理流程

```
1. 解析命令行参数，构建 TestConfig
2. 验证参数：
   - 检查是否指定了路径
   - 检查是否指定了多个类型标志（互斥）
3. 执行检测：
   - 如果指定了类型标志，检测指定类型
   - 如果未指定类型标志，检测路径是否存在
4. 输出结果：
   - 如果检测成功，正常退出（退出码 0）
   - 如果检测失败，报错退出（退出码 1）
```

### 3. 检测路径是否存在

```go
import (
    "os"
)

// 检测路径是否存在
func testPath(path string) error {
    _, err := os.Stat(path)
    if err != nil {
        if os.IsNotExist(err) {
            return fmt.Errorf("路径不存在: %s", path)
        }
        return fmt.Errorf("检测路径失败: %w", err)
    }
    return nil
}
```

### 4. 检测文件是否存在

```go
import (
    "os"
)

// 检测文件是否存在
func testFile(path string) error {
    info, err := os.Stat(path)
    if err != nil {
        if os.IsNotExist(err) {
            return fmt.Errorf("文件不存在: %s", path)
        }
        return fmt.Errorf("检测文件失败: %w", err)
    }

    if info.IsDir() {
        return fmt.Errorf("不是文件: %s", path)
    }

    return nil
}
```

### 5. 检测目录是否存在

```go
import (
    "os"
)

// 检测目录是否存在
func testDir(path string) error {
    info, err := os.Stat(path)
    if err != nil {
        if os.IsNotExist(err) {
            return fmt.Errorf("目录不存在: %s", path)
        }
        return fmt.Errorf("检测目录失败: %w", err)
    }

    if !info.IsDir() {
        return fmt.Errorf("不是目录: %s", path)
    }

    return nil
}
```

### 6. 检测符号链接是否存在

```go
import (
    "os"
)

// 检测符号链接是否存在
func testLink(path string) error {
    info, err := os.Lstat(path)
    if err != nil {
        if os.IsNotExist(err) {
            return fmt.Errorf("符号链接不存在: %s", path)
        }
        return fmt.Errorf("检测符号链接失败: %w", err)
    }

    if info.Mode()&os.ModeSymlink == 0 {
        return fmt.Errorf("不是符号链接: %s", path)
    }

    return nil
}
```

## 边界情况处理

1. **未指定路径**：报错，提示需要指定路径
2. **指定多个类型标志**：报错，提示只能指定一个类型标志
3. **路径不存在**：报错，提示路径不存在
4. **类型不匹配**：报错，提示不是指定类型
5. **权限不足**：报错，提示权限问题
6. **符号链接指向不存在的目标**：符号链接本身存在则返回成功

## 安全考虑

### 1. 路径验证
- 检查路径是否合法
- 验证路径长度限制
- 防止路径遍历攻击

### 2. 权限验证
- 验证用户是否有读取路径的权限
- 提供友好的权限错误提示

## 测试用例

### 单元测试
- 测试检测路径是否存在
- 测试检测文件是否存在
- 测试检测目录是否存在
- 测试检测符号链接是否存在
- 测试边界情况

### 集成测试
- 测试命令行参数解析
- 测试选项组合
- 测试错误处理
- 测试退出码

### 安全测试
- 测试路径验证
- 测试权限验证

## 注意事项

1. 与系统 `test` 命令的兼容性：保持基本用法一致
2. 跨平台支持：Windows 和 Linux/macOS 的文件操作
3. 错误处理：提供友好的错误提示
4. 退出码：成功返回 0，失败返回 1
5. 符号链接：使用 `os.Lstat()` 而不是 `os.Stat()` 来检测符号链接本身

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
- `os.Stat()`：获取文件信息（跟随符号链接）
- `os.Lstat()`：获取文件信息（不跟随符号链接）
- `os.IsNotExist()`：判断错误是否为"不存在"
- `os.ModeSymlink`：符号链接模式
- `fmt.Println()`：输出信息

## 后续扩展（可选）

1. 支持检测文件权限（如 -r、-w、-x）
2. 支持检测文件大小（如 -s）
3. 支持检测文件类型（如 -b 块设备、-c 字符设备）
4. 支持检测文件是否为空（如 -e）
5. 支持检测文件修改时间（如 -newer、-older）
6. 支持组合条件（如 -a、-o）
