# pack 命令多文件/多目录打包设计方案

## 一、需求概述

当前 `pack` 命令内部只支持打包单个文件或目录，需要扩展支持：
- 同时打包多个文件
- 同时打包多个目录
- 混合打包文件和目录
- 支持通配符匹配

## 二、核心思路

当指定多个源路径时：
1. 创建临时目录
2. 将所有源文件/目录复制到临时目录
3. 打包临时目录
4. 清理临时目录

## 三、设计方案

### 3.1 命令语法

```bash
# 单文件/目录（现有功能）
fck pack source.txt
fck pack source/

# 多文件
fck pack file1.txt file2.txt file3.txt

# 多目录
fck pack dir1/ dir2/ dir3/

# 混合模式
fck pack file1.txt dir1/ file2.txt

# 使用通配符
fck pack *.txt
fck pack src/*

# 指定输出路径
fck pack -o archive.zip file1.txt file2.txt
fck pack -o backup.tar.gz dir1/ dir2/
```

### 3.2 参数解析规则

- 所有位置参数 → 源路径列表（支持通配符）
- `-o, --output` 标志 → 压缩包输出路径（可选）

**默认压缩包名生成规则**（多源时）：
1. 如果指定了 `-o`，使用指定值
2. 如果未指定 `-o`：
   - 使用第一个源路径的名称（去掉扩展名）
   - 添加 `.zip` 后缀
   - 示例：`fck pack file1.txt file2.txt` → `file1.zip`

### 3.3 配置结构更新

```go
// PackConfig 打包配置结构体
type PackConfig struct {
    SrcPaths         []string // 源路径列表（支持通配符）
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

### 4.1 源路径展开（支持通配符）

```go
import "gitee.com/MM-Q/go-kit/fs"

// expandSourcePaths 展开源路径列表（支持通配符）
func expandSourcePaths(patterns []string) ([]string, error) {
    // 使用 fs.Expand 展开通配符（包含文件和目录）
    paths, err := fs.Expand(patterns)
    if err != nil {
        return nil, err
    }

    if len(paths) == 0 {
        return nil, errors.New("no sources to pack")
    }

    return paths, nil
}
```

### 4.2 多源打包主逻辑

```go
// PackCmdMain 执行打包命令
func PackCmdMain(config PackConfig) error {
    if len(config.SrcPaths) == 0 {
        return errors.New("no source paths specified")
    }

    // 展开源路径（支持通配符）
    srcPaths, err := expandSourcePaths(config.SrcPaths)
    if err != nil {
        return err
    }

    // 确定压缩包路径
    packPath := config.PackPath
    if packPath == "" {
        packPath = generateDefaultPackPath(srcPaths[0])
    }

    // 单源直接打包，多源需要合并
    var srcPath string
    if len(srcPaths) == 1 {
        srcPath = srcPaths[0]
    } else {
        // 多源：创建临时目录合并
        tempDir, err := mergeSourcesToTemp(srcPaths)
        if err != nil {
            return fmt.Errorf("failed to merge sources: %w", err)
        }
        defer cleanupTempDir(tempDir) // 清理临时目录
        srcPath = tempDir
    }

    // 执行打包
    return packSingle(packPath, srcPath, config)
}

// packSingle 打包单个源
func packSingle(packPath, srcPath string, config PackConfig) error {
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

    return comprx.PackOptions(packPath, srcPath, config)
}
```

### 4.3 临时目录合并逻辑

```go
import (
    "os"
    "path/filepath"
    "gitee.com/MM-Q/go-kit/fs"
)

// mergeSourcesToTemp 将多个源合并到临时目录
func mergeSourcesToTemp(srcPaths []string) (string, error) {
    // 创建临时目录
    tempDir, err := os.MkdirTemp("", "fck-pack-*")
    if err != nil {
        return "", fmt.Errorf("failed to create temp dir: %w", err)
    }

    // 复制每个源到临时目录
    for _, src := range srcPaths {
        if err := copyToTemp(tempDir, src); err != nil {
            cleanupTempDir(tempDir)
            return "", fmt.Errorf("failed to copy %s: %w", src, err)
        }
    }

    return tempDir, nil
}

// copyToTemp 将源复制到临时目录
func copyToTemp(tempDir, src string) error {
    // 获取源文件名
    base := filepath.Base(src)
    dst := filepath.Join(tempDir, base)

    // 检查目标是否已存在（处理同名文件/目录）
    if _, err := os.Stat(dst); err == nil {
        // 同名存在，添加序号后缀
        dst = generateUniqueName(tempDir, base)
    }

    // 复制文件或目录
    info, err := os.Stat(src)
    if err != nil {
        return err
    }

    if info.IsDir() {
        return fs.CopyDir(src, dst)
    }
    return fs.CopyFile(src, dst)
}

// generateUniqueName 生成唯一名称（处理同名冲突）
func generateUniqueName(dir, base string) string {
    ext := filepath.Ext(base)
    name := strings.TrimSuffix(base, ext)

    for i := 1; i < 1000; i++ {
        newName := fmt.Sprintf("%s_%d%s", name, i, ext)
        dst := filepath.Join(dir, newName)
        if _, err := os.Stat(dst); os.IsNotExist(err) {
            return dst
        }
    }

    // 兜底方案使用时间戳
    return filepath.Join(dir, fmt.Sprintf("%s_%d%s", name, time.Now().Unix(), ext))
}

// cleanupTempDir 清理临时目录
func cleanupTempDir(tempDir string) {
    if tempDir != "" {
        _ = os.RemoveAll(tempDir)
    }
}
```

## 五、CLI 层更新

### 5.1 参数解析更新

```go
func runPack(cmd qflag.Command) error {
    args := cmd.Args()
    if len(args) < 1 {
        return fmt.Errorf("at least one source path is required")
    }

    // 所有位置参数都是源路径（支持通配符）
    srcPaths := args

    config := pack.PackConfig{
        SrcPaths:         srcPaths,
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

### 5.2 帮助文档更新

```go
cmdOpts := &qflag.CmdOpts{
    Desc: "智能打包压缩工具",
    Notes: []string{
        "支持的格式有: .zip, .tar, .tar.gz, .tgz, .gz, .bz2, .bzip2, .zlib",
        "通过文件后缀指定压缩格式",
        "不指定 -o 时，自动生成源文件名+.zip",
        "支持多文件/多目录同时打包",
        "支持通配符匹配 (如 *.txt)",
    },
    UseChinese:  true,
    UsageSyntax: fmt.Sprintf("%s pack [options] <source...>", qflag.Root.Name()),
    Examples: map[string]string{
        "打包单文件":        fmt.Sprintf("%s pack source.txt", qflag.Root.Name()),
        "打包单目录":        fmt.Sprintf("%s pack source/", qflag.Root.Name()),
        "打包多文件":        fmt.Sprintf("%s pack file1.txt file2.txt", qflag.Root.Name()),
        "打包多目录":        fmt.Sprintf("%s pack dir1/ dir2/", qflag.Root.Name()),
        "使用通配符":        fmt.Sprintf("%s pack *.txt", qflag.Root.Name()),
        "指定输出路径":      fmt.Sprintf("%s pack -o backup.zip file1.txt file2.txt", qflag.Root.Name()),
    },
}
```

## 六、关键问题处理

### 6.1 同名文件/目录冲突

当多个源有相同名称时：
- 第一个使用原名称
- 后续添加序号后缀（`_1`, `_2`...）
- 示例：`dir1/file.txt` 和 `dir2/file.txt` → `file.txt`, `file_1.txt`

### 6.2 临时目录清理

- 使用 `defer` 确保清理
- 即使打包失败也要清理
- 使用 `os.RemoveAll` 递归删除

### 6.3 错误处理

- 复制失败时立即清理已创建的临时目录
- 返回详细的错误信息（包含具体哪个源失败）

## 七、使用示例

### 多文件打包
```bash
$ fck pack file1.txt file2.txt file3.txt
# 生成: file1.zip (包含 file1.txt, file2.txt, file3.txt)
```

### 多目录打包
```bash
$ fck pack dir1/ dir2/ dir3/
# 生成: dir1.zip (包含 dir1/, dir2/, dir3/)
```

### 混合打包
```bash
$ fck pack file1.txt dir1/ file2.txt
# 生成: file1.zip (包含 file1.txt, dir1/, file2.txt)
```

### 通配符打包
```bash
$ fck pack *.txt
$ fck pack -o docs.zip *.md *.txt
```

## 八、实现步骤

1. 更新 `PackConfig` 结构体，`SrcPath string` → `SrcPaths []string`
2. 创建 `expandSourcePaths` 函数支持通配符
3. 创建 `mergeSourcesToTemp` 函数合并多源到临时目录
4. 创建 `copyToTemp` 和 `generateUniqueName` 辅助函数
5. 更新 `PackCmdMain` 支持多源处理
6. 更新 CLI 层参数解析（支持多位置参数）
7. 更新帮助文档和示例
8. 测试验证

## 九、参考实现

参考 `unpack` 命令的实现方式：
- `internal/commands/unpack/cmd_unpack.go`
- `internal/cli/unpack.go`

参考 `go-kit/fs` 包的文件操作：
- `fs.Expand` - 展开通配符
- `fs.CopyFile` - 复制文件
- `fs.CopyDir` - 复制目录

---

**设计日期**: 2026-05-14  
**预计工作量**: 中（约 150-200 行代码修改）
