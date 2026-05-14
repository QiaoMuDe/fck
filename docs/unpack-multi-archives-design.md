# unpack 命令多压缩包处理设计方案

## 一、需求概述

当前 `unpack` 命令仅支持解压单个压缩包，需要扩展支持：
- 同时处理多个压缩包
- 支持通配符匹配（如 `*.zip`）
- 批量解压统计信息
- 保持向后兼容（现有单文件解压语法）

## 二、现有实现分析

### 当前命令语法
```bash
fck unpack archive.zip              # 解压到当前目录
fck unpack archive.zip /path/to/dst # 解压到指定目录
```

### 当前限制
- 只接受一个压缩包路径参数
- 不支持通配符展开
- 不支持批量处理统计

## 三、新设计方案

### 3.1 命令语法

```bash
# 单文件（向后兼容）
fck unpack archive.zip
fck unpack -o /path/to/dst archive.zip 

# 多文件（新功能）
fck unpack archive1.zip archive2.tar.gz
fck unpack *.zip *.tar.gz
fck unpack -o /path/to/dst *.zip
fck unpack -o /path/to/dst *.tar.gz

# 混合模式
fck unpack  -o /path/to/dst archive1.zip archive2.zip
```

### 3.2 参数解析规则

**方案**：使用 `-o, --output` 标志指定输出目录，默认为当前目录 `.`

- 所有位置参数 → 视为压缩包路径（支持通配符）
- `-o` 标志 → 指定解压目标目录（默认为 `.`）

```go
// 解析逻辑
args := cmd.Args()
if len(args) == 0 {
    return error("no archive specified")
}

// 所有位置参数都是压缩包路径
archives = args
// 目标路径从 -o 标志获取，默认为 "."
dstPath = unpackOutput.Get()
```

### 3.3 标志调整

新增/修改标志：
- `-o, --output` - 解压目标目录（新增，默认 `.`）

注意：此 `-o` 标志与现有设计一致，但之前是位置参数，现在改为显式标志。

保持现有标志不变：
- `-i, --include` - 包含的文件模式
- `-e, --exclude` - 排除的文件模式
- `-f, --overwrite` - 覆盖已存在文件
- `-p, --progress` - 显示进度
- 等等...

### 3.4 配置结构更新

```go
// UnpackConfig 解包配置结构体
type UnpackConfig struct {
    Archives        []string // 压缩包路径列表（支持通配符）
    DstPath         string   // 目标路径
    IncludePatterns []string // 包含的文件或目录
    ExcludePatterns []string // 排除的文件或目录
    MinSize         int64    // 最小文件大小
    MaxSize         int64    // 最大文件大小
    Overwrite       bool     // 覆盖已存在的文件
    Progress        bool     // 显示进度
    ProgressStyle   string   // 进度样式
    NoValidate      bool     // 是否禁用路径验证
}
```

## 四、核心实现逻辑

### 4.1 文件展开（支持通配符）

使用 `go-kit/fs` 包的 `ExpandFiles` 函数直接展开通配符（自动过滤目录，只返回文件）：

```go
import "gitee.com/MM-Q/go-kit/fs"

// 在 UnpackCmdMain 中直接调用
archives, err := fs.ExpandFiles(config.Archives)
if err != nil {
    return err
}
```

### 4.2 批量解压主逻辑

```go
// UnpackCmdMain 执行解包命令
func UnpackCmdMain(config UnpackConfig) error {
    // 1. 展开压缩包列表（支持通配符）
    archives, err := fs.ExpandFiles(config.Archives)
    if err != nil {
        return err
    }

    if len(archives) == 0 {
        return errors.New("no archives to process")
    }

    // 2. 统计信息
    stats := struct {
        total   int
        success int
        failed  int
    }{}

    // 3. 处理每个压缩包
    for _, archive := range archives {
        stats.total++

        // 构建单个压缩包的配置
        singleConfig := UnpackConfig{
            Archives:        []string{archive},
            DstPath:         config.DstPath,
            IncludePatterns: config.IncludePatterns,
            ExcludePatterns: config.ExcludePatterns,
            MinSize:         config.MinSize,
            MaxSize:         config.MaxSize,
            Overwrite:       config.Overwrite,
            Progress:        config.Progress,
            ProgressStyle:   config.ProgressStyle,
            NoValidate:      config.NoValidate,
        }

        if err := unpackSingle(singleConfig); err != nil {
            stats.failed++
            fmt.Fprintf(os.Stderr, "%s: %v\n", archive, err)
            continue
        }

        stats.success++
    }

    // 4. 输出统计信息（多个文件时）
    if stats.total > 1 {
        fmt.Printf("\nTotal: %d archives\n", stats.total)
        fmt.Printf("Success: %d\n", stats.success)
        if stats.failed > 0 {
            fmt.Printf("Failed: %d\n", stats.failed)
        }
    }

    if stats.failed > 0 {
        return fmt.Errorf("unpack completed with %d errors", stats.failed)
    }

    return nil
}

// unpackSingle 解压单个压缩包
func unpackSingle(config UnpackConfig) error {
    if len(config.Archives) != 1 {
        return errors.New("unpackSingle requires exactly one archive")
    }

    archive := config.Archives[0]

    // 使用 comprx 解压
    filter := comprx.FilterOptions{
        Include: config.IncludePatterns,
        Exclude: config.ExcludePatterns,
        MinSize: config.MinSize,
        MaxSize: config.MaxSize,
    }

    progressStyleVal, isValid := types.GetProgressStyle(config.ProgressStyle)
    if !isValid {
        return fmt.Errorf("invalid progress style: %s", config.ProgressStyle)
    }

    opts := comprx.Options{
        CompressionLevel:      comprx.CompressionLevelDefault,
        OverwriteExisting:     config.Overwrite,
        ProgressEnabled:       config.Progress,
        ProgressStyle:         progressStyleVal,
        DisablePathValidation: config.NoValidate,
        Filter:                filter,
    }

    return comprx.UnpackOptions(archive, config.DstPath, opts)
}
```

## 五、CLI 层更新

### 5.1 参数解析更新

```go
func runUnpack(cmd qflag.Command) error {
    args := cmd.Args()
    if len(args) < 1 {
        return fmt.Errorf("at least one archive path is required")
    }

    // 所有位置参数都是压缩包路径（支持通配符）
    archives := args
    // 目标路径从 -o 标志获取
    dstPath := unpackOutput.Get()

    config := unpack.UnpackConfig{
        Archives:        archives,
        DstPath:         dstPath,
        IncludePatterns: unpackIncludePatterns.Get(),
        ExcludePatterns: unpackExcludePatterns.Get(),
        MinSize:         unpackMinSize.Get(),
        MaxSize:         unpackMaxSize.Get(),
        Overwrite:       unpackOverwrite.Get(),
        Progress:        unpackProgress.Get(),
        ProgressStyle:   unpackProgressStyle.Get(),
        NoValidate:      unpackNoValidate.Get(),
    }

    return unpack.UnpackCmdMain(config)
}
```

### 5.2 帮助文档更新

```go
cmdOpts := &qflag.CmdOpts{
    Desc:        "智能解压缩工具",
    Notes: []string{
        "支持的格式有: .zip, .tar, .tar.gz, .tgz, .gz, .bz2, .bzip2, .zlib",
        "支持同时解压多个压缩包",
        "支持通配符匹配 (如 *.zip)",
        "使用 -o 指定解压目标目录，默认为当前目录",
    },
    UseChinese:  true,
    UsageSyntax: fmt.Sprintf("%s unpack [options] <archive...>", qflag.Root.Name()),
    Examples: map[string]string{
        "解压单个压缩包":      fmt.Sprintf("%s unpack archive.zip", qflag.Root.Name()),
        "解压到指定目录":      fmt.Sprintf("%s unpack archive.zip -o /path/to/dst", qflag.Root.Name()),
        "解压多个压缩包":      fmt.Sprintf("%s unpack archive1.zip archive2.tar.gz", qflag.Root.Name()),
        "使用通配符":         fmt.Sprintf("%s unpack *.zip", qflag.Root.Name()),
        "批量解压到目录":      fmt.Sprintf("%s unpack *.tar.gz -o /path/to/dst", qflag.Root.Name()),
    },
}
```

## 六、输出格式

### 单文件解压（向后兼容）
```
$ fck unpack archive.zip
# 正常输出，无统计信息
```

### 多文件解压
```
$ fck unpack *.zip
Processing: archive1.zip
Processing: archive2.zip
Processing: archive3.zip

Total: 3 archives
Success: 3
```

### 有错误时
```
$ fck unpack *.zip
Processing: archive1.zip
Processing: archive2.zip
archive2.zip: file is corrupted
Processing: archive3.zip

Total: 3 archives
Success: 2
Failed: 1
```

## 七、错误处理

1. **通配符匹配失败**：保留原路径，后续验证文件是否存在
2. **文件不存在**：报错并继续处理其他文件
3. **文件是目录**：跳过该路径
4. **解压失败**：记录错误，继续处理其他压缩包
5. **所有文件都失败**：返回错误码

## 八、向后兼容性

- 现有单文件解压语法完全兼容
- 现有标志行为不变
- 新增功能对旧用法无影响

## 九、实现步骤

1. 更新 `UnpackConfig` 结构体，`PackPath` 改为 `Archives []string`
2. 更新 `UnpackCmdMain` 使用 `fs.ExpandFiles` 支持通配符和批量处理
3. 创建 `unpackSingle` 函数处理单个压缩包
4. 更新 CLI 层：
   - 添加 `unpackOutput` 标志（`-o, --output`）
   - 简化参数解析逻辑
   - 更新帮助文档和示例
5. 测试验证

## 十、参考实现

参考 `iconv` 和 `newline` 命令的批量处理逻辑：
- `internal/commands/iconv/cmd_iconv.go:expandFileList`
- `internal/commands/newline/cmd_newline.go:expandArchiveList`

通配符展开使用 `go-kit/fs` 包的 `Expand` 函数：
- `vendor/gitee.com/MM-Q/go-kit/fs/expand.go`

---

**设计日期**: 2026-05-11  
**预计工作量**: 小（约 100-150 行代码修改）
