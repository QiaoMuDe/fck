# pack 命令使用方式重构设计方案

## 一、需求概述

当前 `pack` 命令的使用方式不够直观，需要调整参数传递方式，使其更符合用户直觉：
- 压缩包路径通过标志指定（而非位置参数）
- 位置参数只传递源路径
- 支持自动生成默认压缩包名称

## 二、现有实现分析

### 当前命令语法
```bash
fck pack archive.zip source.txt    # 位置参数1=压缩包路径, 位置参数2=源路径
fck pack archive.tar source/       # 位置参数1=压缩包路径, 位置参数2=源路径
```

### 当前限制
- 必须显式指定压缩包路径
- 参数顺序固定，不够灵活
- 无法快速打包（需要想压缩包名）

## 三、新设计方案

### 3.1 命令语法

```bash
# 快速打包（自动生成压缩包名）
fck pack source.txt              # 生成 source.zip
fck pack source/                 # 生成 source.zip（目录）
fck pack app.exe                 # 生成 app.zip

# 指定压缩包路径
fck pack -o backup.zip source.txt           # 指定名称
fck pack -o backup.tar.gz source/           # 指定名称和格式
fck pack -o archive.zip file1.txt           # 单文件打包
```

### 3.2 参数解析规则

**方案**：
- 位置参数 → 源路径（文件或目录）
- `-o, --output` 标志 → 压缩包输出路径（可选）

**默认压缩包名生成规则**：
1. 如果指定了 `-o`，使用指定值
2. 如果未指定 `-o`：
   - 去掉源路径的扩展名
   - 添加 `.zip` 后缀
   - 示例：`source.txt` → `source.zip`，`app.exe` → `app.zip`

### 3.3 标志调整

新增/修改标志：
- `-o, --output` - 压缩包输出路径（新增，默认空，自动生成）

保持现有标志不变：
- `-i, --include` - 包含的文件模式
- `-e, --exclude` - 排除的文件模式
- `-c, --compression` - 压缩级别
- `-f, --overwrite` - 覆盖已存在文件
- `-p, --progress` - 显示进度
- 等等...

### 3.4 配置结构更新

```go
// PackConfig 打包配置结构体
type PackConfig struct {
    SrcPath          string   // 源路径（文件或目录）
    PackPath         string   // 压缩包输出路径
    IncludePatterns  []string // 包含的文件或目录
    ExcludePatterns  []string // 排除的文件或目录
    MinSize          int64    // 最小文件大小
    MaxSize          int64    // 最大文件大小
    CompressionLevel string   // 压缩级别
    Overwrite        bool     // 覆盖已存在的文件
    Progress         bool     // 显示进度
    ProgressStyle    string   // 进度样式
    NoValidate       bool     // 是否禁用路径验证
}
```

## 四、核心实现逻辑

### 4.1 默认压缩包名生成

```go
import (
    "path/filepath"
    "strings"
)

// generateDefaultPackPath 根据源路径生成默认压缩包路径
func generateDefaultPackPath(srcPath string) string {
    // 获取文件名（不含目录）
    base := filepath.Base(srcPath)
    
    // 去掉现有扩展名
    ext := filepath.Ext(base)
    name := strings.TrimSuffix(base, ext)
    
    // 添加 .zip 后缀
    return name + ".zip"
}
```

### 4.2 打包主逻辑

```go
// PackCmdMain 执行打包命令
func PackCmdMain(config PackConfig) error {
    if config.SrcPath == "" {
        return errors.New("no source path specified")
    }

    // 确定压缩包路径
    packPath := config.PackPath
    if packPath == "" {
        packPath = generateDefaultPackPath(config.SrcPath)
    }

    filter := comprx.FilterOptions{
        Include: config.IncludePatterns,
        Exclude: config.ExcludePatterns,
        MinSize: config.MinSize,
        MaxSize: config.MaxSize,
    }

    compressionLevelVal, isValid := types.GetCompressionLevel(config.CompressionLevel)
    if !isValid {
        return fmt.Errorf("invalid compression level: %s", config.CompressionLevel)
    }

    progressStyleVal, isValid := types.GetProgressStyle(config.ProgressStyle)
    if !isValid {
        return fmt.Errorf("invalid progress style: %s", config.ProgressStyle)
    }

    opts := comprx.Options{
        CompressionLevel:      compressionLevelVal,
        OverwriteExisting:     config.Overwrite,
        ProgressEnabled:       config.Progress,
        ProgressStyle:         progressStyleVal,
        DisablePathValidation: config.NoValidate,
        Filter:                filter,
    }

    return comprx.PackOptions(packPath, config.SrcPath, opts)
}
```

## 五、CLI 层更新

### 5.1 参数解析更新

```go
var (
    // ... 其他标志不变
    packOutput *qflag.StringFlag // 压缩包输出路径
)

func init() {
    // ... 其他标志初始化不变
    packOutput = PackCmd.String("output", "o", "压缩包输出路径，默认为源文件名+.zip", "")

    cmdOpts := &qflag.CmdOpts{
        Desc: "智能打包压缩工具",
        Notes: []string{
            "支持的格式有: .zip, .tar, .tar.gz, .tgz, .gz, .bz2, .bzip2, .zlib",
            "通过文件后缀指定压缩格式",
            "不指定 -o 时，自动生成源文件名+.zip",
        },
        UseChinese:  true,
        UsageSyntax: fmt.Sprintf("%s pack [options] <source>", qflag.Root.Name()),
        Examples: map[string]string{
            "快速打包文件":      fmt.Sprintf("%s pack source.txt", qflag.Root.Name()),
            "快速打包目录":      fmt.Sprintf("%s pack source/", qflag.Root.Name()),
            "指定压缩包名称":    fmt.Sprintf("%s pack -o backup.zip source.txt", qflag.Root.Name()),
            "指定压缩格式":      fmt.Sprintf("%s pack -o backup.tar.gz source/", qflag.Root.Name()),
        },
    }
}

func runPack(cmd qflag.Command) error {
    args := cmd.Args()
    if len(args) < 1 {
        return fmt.Errorf("source path is required")
    }

    srcPath := args[0] // 源路径

    config := pack.PackConfig{
        SrcPath:          srcPath,
        PackPath:         packOutput.Get(),
        IncludePatterns:  packIncludePatterns.Get(),
        ExcludePatterns:  packExcludePatterns.Get(),
        MinSize:          packMinSize.Get(),
        MaxSize:          packMaxSize.Get(),
        CompressionLevel: packCompressionLevel.Get(),
        Overwrite:        packOverwrite.Get(),
        Progress:         packProgress.Get(),
        ProgressStyle:    packProgressStyle.Get(),
        NoValidate:       packNoValidate.Get(),
    }

    return pack.PackCmdMain(config)
}
```

## 六、使用示例

### 快速打包
```bash
# 文件打包
$ fck pack app.exe
# 生成: app.zip

# 目录打包
$ fck pack myproject/
# 生成: myproject.zip

# 带扩展名的文件
$ fck pack document.txt
# 生成: document.zip
```

### 指定输出路径
```bash
# 指定名称
$ fck pack -o backup.zip source.txt

# 指定格式
$ fck pack -o archive.tar.gz source/

# 指定目录
$ fck pack -o /path/to/backup.zip file.txt
```

## 七、向后兼容性

**此方案为破坏性变更**，旧用法不再兼容：
- 旧：`fck pack archive.zip source.txt`
- 新：`fck pack source.txt -o archive.zip`

需要在版本说明中标注此变更。

## 八、实现步骤

1. 添加 `generateDefaultPackPath` 函数生成默认压缩包名
2. 更新 `PackConfig` 结构体字段顺序（SrcPath 放前面）
3. 更新 `PackCmdMain` 添加默认路径生成逻辑
4. 更新 CLI 层：
   - 添加 `packOutput` 标志
   - 修改参数解析逻辑
   - 更新帮助文档和示例
5. 测试验证

## 九、参考实现

参考 `unpack` 命令的实现方式：
- `internal/commands/unpack/cmd_unpack.go`
- `internal/cli/unpack.go`

---

**设计日期**: 2026-05-14  
**预计工作量**: 小（约 50-80 行代码修改）
