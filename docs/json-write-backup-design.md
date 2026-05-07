# JSON 命令原地写入与备份功能设计方案

## 1. 功能概述

为 `fck json` 命令增加原地写入（-w）和备份（-b）功能，允许用户直接将处理后的 JSON 写回原文件。

## 2. 命令行参数设计

### 新增 Flags

| 短标志 | 长标志 | 类型 | 默认值 | 说明 |
|--------|--------|------|--------|------|
| -w | --write | bool | false | 原地写入处理后的 JSON 到原文件 |
| -b | --backup | bool | false | 写入前创建备份文件（仅在 -w 生效时有效） |

### 互斥组调整

```go
MutexGroups: []qflag.MutexGroup{
    {
        Name:      "operation-mode",
        Flags:     []string{"pretty", "compact", "validate", "query"},
        AllowNone: true,
    },
},
```

## 3. 核心处理逻辑

```
处理流程：
┌─────────────────────────────────────┐
│ 1. 解析输入（管道或文件）              │
└─────────────┬───────────────────────┘
              ▼
┌─────────────────────────────────────┐
│ 2. 处理 JSON（格式化/压缩/验证/查询）  │
│    - 验证 JSON 有效性                │
│    - 如果无效，直接返回错误           │
└─────────────┬───────────────────────┘
              ▼
┌─────────────────────────────────────┐
│ 3. 检查 -w 标志                       │
│    └─ 否 → 输出到 stdout            │
│    └─ 是 → 继续下一步                │
└─────────────┬───────────────────────┘
              ▼
┌─────────────────────────────────────┐
│ 4. 验证写入条件                       │
│    - 检查是否有文件参数               │
│    - 无文件参数 → 报错                │
└─────────────┬───────────────────────┘
              ▼
┌─────────────────────────────────────┐
│ 5. 检查 -b 标志（备份）               │
│    └─ 是 → 创建 .bak 备份文件        │
│    └─ 否 → 跳过                      │
└─────────────┬───────────────────────┘
              ▼
┌─────────────────────────────────────┐
│ 6. 写入原文件                         │
│    - 使用临时文件方式确保原子性       │
│    - 写入失败时恢复原文件             │
└─────────────────────────────────────┘
```

## 4. 备份文件策略

### 命名规则

- 原文件：`data.json`
- 备份文件：`data.json.bak`

### 备份已存在时的处理

- 策略：直接覆盖旧备份
- 原因：保持简单，避免产生多个备份文件

## 5. 安全机制

### 5.1 JSON 验证

- 在写入前必须验证 JSON 有效性
- 无效 JSON 直接返回错误，不写入文件

### 5.2 原子写入

```go
// 使用临时文件实现原子写入
1. 写入到临时文件：data.json.tmp
2. 重命名临时文件覆盖原文件：data.json.tmp → data.json
3. 如果失败，保留原文件
```

### 5.3 错误恢复

- 备份创建失败：报错，不继续写入
- 写入失败：如果已创建备份，尝试从备份恢复（可选）

## 6. 代码结构修改

### 6.1 cli/json.go

```go
// 新增 flags
jsonWrite   *qflag.BoolFlag   // -w, --write   原地写入
jsonBackup  *qflag.BoolFlag   // -b, --backup  写入前备份

// init 函数中初始化
jsonWrite = JsonCmd.Bool("write", "w", "原地写入处理后的 JSON 到原文件", false)
jsonBackup = JsonCmd.Bool("backup", "b", "写入前创建备份文件", false)

// Notes 中新增说明
"使用 -w 原地写入时，必须指定文件路径参数",
"使用 -b 备份时，会创建 .bak 后缀的备份文件",

// Examples 中新增示例
"原地格式化文件":    fmt.Sprintf("%s json -p -w data.json", qflag.Root.Name()),
"原地压缩并备份":    fmt.Sprintf("%s json -c -w -b config.json", qflag.Root.Name()),
"查询结果写入文件":  fmt.Sprintf("%s json -q -w users.0 data.json", qflag.Root.Name()),

// runJson 中传递配置
config := json.JsonConfig{
    // ... 原有字段
    Write:  jsonWrite.Get(),
    Backup: jsonBackup.Get(),
}
```

### 6.2 commands/json/cmd_json.go

```go
// JsonConfig 结构体新增字段
type JsonConfig struct {
    // ... 原有字段
    Write   bool     // 原地写入
    Backup  bool     // 创建备份
    Files   []string // 文件路径列表
}

// readInput 读取输入数据（修改后）
func readInput(config JsonConfig) ([]byte, error) {
    // 1. 优先检测管道/重定向输入
    if term.IsStdinPipe() {
        // 管道输入时禁用 -w
        if config.Write {
            return nil, fmt.Errorf("cannot use -w flag with pipe input")
        }
        return readAllStdin()
    }

    // 2. 从文件读取
    if len(config.Files) > 0 {
        if len(config.Files) > 1 {
            return nil, fmt.Errorf("only one file path can be specified")
        }
        return os.ReadFile(config.Files[0])
    }

    return nil, fmt.Errorf("no input data provided")
}

// JsonCmdMain 修改逻辑
func JsonCmdMain(config JsonConfig) error {
    // 1. 获取输入数据（管道或文件）
    input, err := readInput(config)
    if err != nil {
        return err
    }

    // 2. 处理 JSON
    result, err := processJSON(input, config)
    if err != nil {
        return err
    }

    // 3. 检查是否需要原地写入
    if config.Write {
        targetFile := config.Files[0]

        // 创建备份（如果启用）
        if config.Backup {
            backupFile := targetFile + ".bak"
            if err := fs.CopyEx(targetFile, backupFile, true); err != nil {
                return fmt.Errorf("failed to create backup: %w", err)
            }
        }

        // 写入文件（原子写入：临时文件+重命名）
        tmpFile := targetFile + ".tmp"
        if err := os.WriteFile(tmpFile, []byte(result), 0644); err != nil {
            return fmt.Errorf("failed to write temp file: %w", err)
        }
        if err := fs.MoveEx(tmpFile, targetFile, true); err != nil {
            return fmt.Errorf("failed to move temp file: %w", err)
        }

        fmt.Printf("Written to: %s\n", targetFile)
        return nil
    }

    // 4. 输出到 stdout
    fmt.Println(result)
    return nil
}

// 使用 fs.CopyEx 创建备份
// fs.CopyEx(src, dst, overwrite bool) error

// 原子写入使用 os.WriteFile 或临时文件+重命名
```

## 7. 使用示例

```bash
# 原地格式化文件（美化输出）
fck json -pw data.json

# 原地压缩 JSON 并创建备份
fck json -cwb config.json

# 查询结果写入原文件
fck json -qw "users.0.name" data.json

# 验证并格式化（验证失败则不写入）
fck json -vpw data.json

# 高亮显示并写入（高亮只影响终端输出，写入的是纯文本）
fck json -pHw data.json
```

## 8. 错误场景

```bash
# 错误：管道输入时启用 -w
$ echo '{"a":1}' | fck json -pw
Error: -w flag requires a file path argument

# 错误：备份创建失败
$ fck json -pwb readonly.json
Error: failed to create backup: permission denied

# 错误：JSON 无效时不写入
$ fck json -pw invalid.json
Error: invalid JSON syntax
```

## 9. 注意事项

1. **管道输入冲突**：`-w` 不能与管道输入同时使用，因为没有文件可写入
2. **多文件处理**：如果指定多个文件，只处理第一个文件（或报错）
3. **备份覆盖**：备份文件已存在时会直接覆盖
4. **原子性**：使用临时文件+重命名确保写入的原子性
5. **权限问题**：写入前检查文件权限，避免写入失败
