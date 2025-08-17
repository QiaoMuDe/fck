# FCK工具 - COMPRESS/EXTRACT多格式压缩解压子命令设计方案

## 📋 概述

本文档详细描述了为FCK工具新增compress和extract子命令的完整设计方案。采用智能格式识别的统一命令设计，支持ZIP、TAR、TAR.GZ、TAR.BZ2、7Z等多种压缩格式，提供简洁易用的用户体验。

## 🗜️ COMPRESS子命令设计

### 核心功能
- 智能识别输出格式（基于文件扩展名）
- 支持多种压缩格式：ZIP、TAR、TAR.GZ、TAR.BZ2、TAR.XZ、7Z
- 支持多种压缩级别
- 支持密码保护（ZIP、7Z格式）
- 支持排除/包含模式
- 支持进度显示和详细输出

### 支持的压缩格式

| 格式 | 扩展名 | 压缩 | 密码 | 分卷 | 说明 |
|------|--------|------|------|------|------|
| ZIP | .zip | ✅ | ✅ | ✅ | 最常用格式，跨平台兼容性好 |
| TAR | .tar | ✅ | ❌ | ❌ | Unix传统归档格式，无压缩 |
| TAR.GZ | .tar.gz, .tgz | ✅ | ❌ | ❌ | TAR+GZIP，Linux常用 |
| TAR.BZ2 | .tar.bz2, .tbz2 | ✅ | ❌ | ❌ | TAR+BZIP2，压缩比更高 |
| TAR.XZ | .tar.xz, .txz | ✅ | ❌ | ❌ | TAR+XZ，最高压缩比 |
| 7Z | .7z | ✅ | ✅ | ✅ | 高压缩比，功能丰富 |

### 命令语法
```bash
fck compress [options] <archive> <files/dirs...>
```

### 参数说明

#### 必需参数
- `<archive>` - 输出的压缩文件名（格式由扩展名自动识别）
- `<files/dirs...>` - 要压缩的文件或目录

#### 可选标志
| 标志 | 长标志 | 参数 | 描述 |
|------|--------|------|------|
| `-f` | `--format` | `<format>` | 强制指定格式 (zip, tar, tgz, tbz2, txz, 7z) |
| `-l` | `--level` | `<0-9>` | 压缩级别 (0=无压缩, 9=最大压缩, 默认6) |
| `-p` | `--password` | `<pwd>` | 设置密码保护 (仅ZIP、7Z格式) |
| `-r` | `--recursive` | - | 递归压缩目录 (默认启用) |
| `-x` | `--exclude` | `<pattern>` | 排除匹配模式的文件 (支持多个) |
| `-i` | `--include` | `<pattern>` | 仅包含匹配模式的文件 |
| `--force` | `--force` | - | 强制覆盖已存在的压缩文件 |
| `-v` | `--verbose` | - | 显示详细信息 |
| `-q` | `--quiet` | - | 静默模式 |
| `-t` | `--test` | - | 压缩后测试文件完整性 |
| `-s` | `--split` | `<size>` | 分卷压缩 (如: 100MB, 1GB，仅ZIP、7Z) |
| | `--progress` | - | 显示进度条 |
| | `--preserve-permissions` | - | 保留文件权限 (Unix系统) |
| | `--preserve-timestamps` | - | 保留时间戳 |
| | `--comment` | `<text>` | 添加压缩文件注释 |
| | `--threads` | `<num>` | 并行压缩线程数 (默认CPU核心数) |

### 使用示例

#### 自动格式识别压缩
```bash
# ZIP格式
fck compress backup.zip ./documents ./photos

# TAR.GZ格式
fck compress backup.tar.gz ./documents ./photos

# 7Z格式
fck compress backup.7z ./documents ./photos
```

#### 高压缩比 + 密码保护
```bash
fck compress -l 9 -p mypassword secure.zip ./sensitive_data
fck compress -l 9 -p mypassword secure.7z ./sensitive_data
```

#### 排除特定文件类型
```bash
fck compress -x "*.tmp" -x "*.log" clean.tar.gz ./project
```

#### 分卷压缩
```bash
fck compress -s 100MB large.zip ./big_directory
fck compress -s 500MB large.7z ./big_directory
```

#### 强制指定格式
```bash
fck compress -f zip backup ./documents  # 输出backup.zip
fck compress -f tgz backup ./documents  # 输出backup.tar.gz
```

## 📦 EXTRACT子命令设计

### 核心功能
- 智能识别压缩格式（基于文件扩展名和文件头）
- 支持解压多种格式：ZIP、TAR、TAR.GZ、TAR.BZ2、TAR.XZ、7Z、RAR
- 支持密码保护的压缩文件
- 支持选择性解压
- 支持解压到指定目录
- 支持测试压缩文件完整性

### 支持的解压格式

| 格式 | 扩展名 | 解压 | 密码 | 说明 |
|------|--------|------|------|------|
| ZIP | .zip | ✅ | ✅ | 完全支持 |
| TAR | .tar | ✅ | ❌ | 完全支持 |
| TAR.GZ | .tar.gz, .tgz | ✅ | ❌ | 完全支持 |
| TAR.BZ2 | .tar.bz2, .tbz2 | ✅ | ❌ | 完全支持 |
| TAR.XZ | .tar.xz, .txz | ✅ | ❌ | 完全支持 |
| 7Z | .7z | ✅ | ✅ | 完全支持 |
| RAR | .rar | ✅ | ✅ | 仅解压支持 |

### 命令语法
```bash
fck extract [options] <archive> [destination]
```

### 参数说明

#### 必需参数
- `<archive>` - 要解压的压缩文件

#### 可选参数
- `[destination]` - 解压目标目录 (默认当前目录)

#### 可选标志
| 标志 | 长标志 | 参数 | 描述 |
|------|--------|------|------|
| `-f` | `--format` | `<format>` | 强制指定格式 (auto, zip, tar, tgz, tbz2, txz, 7z, rar) |
| `-p` | `--password` | `<pwd>` | 压缩文件密码 |
| `-l` | `--list` | - | 仅列出压缩文件内容，不解压 |
| `-t` | `--test` | - | 测试压缩文件完整性 |
| `-o` | `--overwrite` | - | 覆盖已存在的文件 |
| `-n` | `--never-overwrite` | - | 从不覆盖已存在的文件 |
| `-u` | `--update` | - | 仅解压更新的文件 |
| `-j` | `--junk-paths` | - | 忽略目录结构，解压到同一目录 |
| `-x` | `--exclude` | `<pattern>` | 排除匹配模式的文件 |
| `-i` | `--include` | `<pattern>` | 仅解压匹配模式的文件 |
| `-v` | `--verbose` | - | 显示详细信息 |
| `-q` | `--quiet` | - | 静默模式 |
| | `--progress` | - | 显示进度条 |
| | `--preserve-permissions` | - | 保留文件权限 |
| | `--preserve-timestamps` | - | 保留时间戳 |
| | `--threads` | `<num>` | 并行解压线程数 |

### 使用示例

#### 自动格式识别解压
```bash
# 自动识别ZIP格式
fck extract backup.zip

# 自动识别TAR.GZ格式
fck extract backup.tar.gz ./restore

# 自动识别7Z格式
fck extract backup.7z
```

#### 带密码解压
```bash
fck extract -p mypassword secure.zip
fck extract -p mypassword secure.7z
```

#### 仅列出内容
```bash
fck extract -l archive.zip
fck extract -l backup.tar.gz
fck extract -l data.7z
```

#### 选择性解压
```bash
fck extract -i "*.txt" -i "*.md" docs.zip
fck extract -x "*.tmp" -x "*.log" backup.tar.gz
```

#### 测试完整性
```bash
fck extract -t backup.zip
fck extract -t backup.7z
```

## 📋 LIST-ARCHIVE子命令设计

### 核心功能
- 列出压缩文件内容
- 支持多种显示格式
- 显示文件详细信息
- 支持过滤和排序

### 命令语法
```bash
fck list-archive [options] <archive>
```

### 参数说明
| 标志 | 长标志 | 参数 | 描述 |
|------|--------|------|------|
| `-f` | `--format` | `<format>` | 输出格式 (table, json, csv, tree) |
| `-s` | `--sort` | `<field>` | 排序字段 (name, size, time, ratio) |
| `-r` | `--reverse` | - | 反向排序 |
| `-H` | `--human-readable` | - | 人类可读的文件大小 |
| `-p` | `--password` | `<pwd>` | 压缩文件密码 |
| | `--filter` | `<pattern>` | 文件名过滤模式 |

### 使用示例
```bash
# 表格格式显示
fck list-archive backup.zip

# JSON格式输出
fck list-archive -f json backup.tar.gz

# 树形结构显示
fck list-archive -f tree backup.7z

# 按大小排序
fck list-archive -s size -r backup.zip
```

## 🏗️ 技术实现方案

### 目录结构
```
commands/
├── compress/
│   ├── cmd_compress.go     # compress子命令主逻辑
│   ├── flags.go           # 标志定义
│   ├── detector.go        # 格式检测器
│   ├── validator.go       # 参数验证
│   ├── progress.go        # 进度跟踪
│   ├── formats/           # 格式处理器
│   │   ├── zip.go         # ZIP格式处理
│   │   ├── tar.go         # TAR格式处理
│   │   ├── targz.go       # TAR.GZ格式处理
│   │   ├── tarbz2.go      # TAR.BZ2格式处理
│   │   ├── tarxz.go       # TAR.XZ格式处理
│   │   ├── sevenzip.go    # 7Z格式处理
│   │   └── interface.go   # 格式处理器接口
│   ├── APIDOC.md          # API文档
│   └── cmd_compress_test.go # 测试文件
├── extract/
│   ├── cmd_extract.go     # extract子命令主逻辑
│   ├── flags.go           # 标志定义
│   ├── detector.go        # 格式检测器
│   ├── validator.go       # 参数验证
│   ├── progress.go        # 进度跟踪
│   ├── formats/           # 格式处理器
│   │   ├── zip.go         # ZIP格式处理
│   │   ├── tar.go         # TAR格式处理
│   │   ├── targz.go       # TAR.GZ格式处理
│   │   ├── tarbz2.go      # TAR.BZ2格式处理
│   │   ├── tarxz.go       # TAR.XZ格式处理
│   │   ├── sevenzip.go    # 7Z格式处理
│   │   ├── rar.go         # RAR格式处理
│   │   ├── auto.go        # 自动格式检测
│   │   └── interface.go   # 格式处理器接口
│   ├── APIDOC.md          # API文档
│   └── cmd_extract_test.go # 测试文件
└── list-archive/
    ├── cmd_list_archive.go # list-archive子命令主逻辑
    ├── flags.go           # 标志定义
    ├── formatter.go       # 输出格式化器
    ├── APIDOC.md          # API文档
    └── cmd_list_archive_test.go # 测试文件
```

### 核心接口设计

#### 1. 格式检测器接口
```go
type FormatDetector interface {
    // 根据文件扩展名检测格式
    DetectByExtension(filename string) (Format, error)
    
    // 根据文件头检测格式
    DetectByHeader(reader io.Reader) (Format, error)
    
    // 自动检测格式（优先文件头，后备扩展名）
    AutoDetect(filename string) (Format, error)
}

// 支持的格式枚举
type Format int

const (
    FormatUnknown Format = iota
    FormatZIP
    FormatTAR
    FormatTARGZ
    FormatTARBZ2
    FormatTARXZ
    Format7Z
    FormatRAR
)
```

#### 2. 压缩器接口
```go
type Compressor interface {
    // 设置压缩选项
    SetOptions(options CompressOptions) error
    
    // 添加文件到压缩包
    AddFile(srcPath, archivePath string) error
    
    // 添加目录到压缩包
    AddDirectory(srcPath, archivePath string) error
    
    // 完成压缩并关闭
    Close() error
    
    // 获取支持的扩展名
    SupportedExtensions() []string
}

// 压缩选项
type CompressOptions struct {
    Level               int      // 压缩级别
    Password            string   // 密码保护
    Comment             string   // 注释
    PreservePermissions bool     // 保留权限
    PreserveTimestamps  bool     // 保留时间戳
    ExcludePatterns     []string // 排除模式
    IncludePatterns     []string // 包含模式
    SplitSize           int64    // 分卷大小
    Threads             int      // 线程数
}
```

#### 3. 解压器接口
```go
type Extractor interface {
    // 设置解压选项
    SetOptions(options ExtractOptions) error
    
    // 列出压缩文件内容
    List() ([]ArchiveFileInfo, error)
    
    // 解压所有文件
    ExtractAll(destPath string) error
    
    // 解压指定文件
    ExtractFiles(files []string, destPath string) error
    
    // 测试压缩文件完整性
    Test() error
    
    // 关闭解压器
    Close() error
}

// 解压选项
type ExtractOptions struct {
    Password            string   // 密码
    Overwrite           bool     // 覆盖已存在文件
    NeverOverwrite      bool     // 从不覆盖
    UpdateOnly          bool     // 仅更新
    JunkPaths           bool     // 忽略路径
    PreservePermissions bool     // 保留权限
    PreserveTimestamps  bool     // 保留时间戳
    ExcludePatterns     []string // 排除模式
    IncludePatterns     []string // 包含模式
    Threads             int      // 线程数
}
```

#### 4. 进度跟踪器
```go
type ProgressTracker struct {
    TotalFiles     int64     // 总文件数
    ProcessedFiles int64     // 已处理文件数
    TotalBytes     int64     // 总字节数
    ProcessedBytes int64     // 已处理字节数
    StartTime      time.Time // 开始时间
    CurrentFile    string    // 当前处理文件
    Operation      string    // 操作类型（压缩/解压）
}

func (p *ProgressTracker) Update(filename string, bytes int64)
func (p *ProgressTracker) GetProgress() float64
func (p *ProgressTracker) GetETA() time.Duration
func (p *ProgressTracker) GetSpeed() int64
func (p *ProgressTracker) Display() // 显示进度条
```

#### 5. 压缩文件信息结构
```go
type ArchiveFileInfo struct {
    Name           string    // 文件名
    Path           string    // 完整路径
    Size           int64     // 原始大小
    CompressedSize int64     // 压缩后大小
    ModTime        time.Time // 修改时间
    IsDir          bool      // 是否为目录
    Mode           os.FileMode // 文件权限
    CRC32          uint32    // CRC32校验值
    CompressionRatio float64 // 压缩比
}
```

### 格式检测实现

#### 智能格式检测逻辑
```go
func (d *FormatDetector) AutoDetect(filename string) (Format, error) {
    // 1. 首先尝试通过文件头检测
    if file, err := os.Open(filename); err == nil {
        defer file.Close()
        if format, err := d.DetectByHeader(file); err == nil && format != FormatUnknown {
            return format, nil
        }
    }
    
    // 2. 后备方案：通过扩展名检测
    return d.DetectByExtension(filename)
}

func (d *FormatDetector) DetectByExtension(filename string) (Format, error) {
    filename = strings.ToLower(filename)
    
    switch {
    case strings.HasSuffix(filename, ".zip"):
        return FormatZIP, nil
    case strings.HasSuffix(filename, ".tar.gz") || strings.HasSuffix(filename, ".tgz"):
        return FormatTARGZ, nil
    case strings.HasSuffix(filename, ".tar.bz2") || strings.HasSuffix(filename, ".tbz2"):
        return FormatTARBZ2, nil
    case strings.HasSuffix(filename, ".tar.xz") || strings.HasSuffix(filename, ".txz"):
        return FormatTARXZ, nil
    case strings.HasSuffix(filename, ".tar"):
        return FormatTAR, nil
    case strings.HasSuffix(filename, ".7z"):
        return Format7Z, nil
    case strings.HasSuffix(filename, ".rar"):
        return FormatRAR, nil
    default:
        return FormatUnknown, fmt.Errorf("不支持的文件格式: %s", filename)
    }
}

func (d *FormatDetector) DetectByHeader(reader io.Reader) (Format, error) {
    // 读取文件头部字节
    header := make([]byte, 16)
    n, err := reader.Read(header)
    if err != nil || n < 4 {
        return FormatUnknown, err
    }
    
    // 检查文件签名
    switch {
    case bytes.HasPrefix(header, []byte("PK\x03\x04")) || bytes.HasPrefix(header, []byte("PK\x05\x06")):
        return FormatZIP, nil
    case bytes.HasPrefix(header, []byte("7z\xBC\xAF\x27\x1C")):
        return Format7Z, nil
    case bytes.HasPrefix(header, []byte("Rar!\x1A\x07\x00")) || bytes.HasPrefix(header, []byte("Rar!\x1A\x07\x01\x00")):
        return FormatRAR, nil
    case bytes.HasPrefix(header, []byte("\x1F\x8B")):
        return FormatTARGZ, nil // 可能是GZIP压缩的TAR
    case bytes.HasPrefix(header, []byte("BZh")):
        return FormatTARBZ2, nil // 可能是BZIP2压缩的TAR
    case bytes.HasPrefix(header, []byte("\xFD7zXZ\x00")):
        return FormatTARXZ, nil // 可能是XZ压缩的TAR
    default:
        // 检查是否为TAR格式（通过TAR头部结构）
        if d.isTarHeader(header) {
            return FormatTAR, nil
        }
        return FormatUnknown, fmt.Errorf("无法识别的文件格式")
    }
}
```

### 依赖库选择

#### Go标准库
- `archive/zip` - ZIP格式支持
- `archive/tar` - TAR格式支持
- `compress/gzip` - GZIP压缩支持
- `compress/bzip2` - BZIP2解压支持
- `compress/lzw` - LZW压缩支持

#### 第三方库
- `github.com/alexmullins/zip` - 支持密码保护的ZIP
- `github.com/ulikunitz/xz` - XZ压缩支持
- `github.com/bodgit/sevenzip` - 7Z格式支持
- `github.com/nwaples/rardecode` - RAR解压支持
- `github.com/klauspost/compress` - 高性能压缩库

### 错误处理策略

#### 错误类型定义
```go
var (
    ErrUnsupportedFormat    = errors.New("不支持的压缩格式")
    ErrInvalidArchive       = errors.New("无效的压缩文件")
    ErrPasswordRequired     = errors.New("需要密码")
    ErrWrongPassword        = errors.New("密码错误")
    ErrFileExists           = errors.New("文件已存在")
    ErrInsufficientSpace    = errors.New("磁盘空间不足")
    ErrPermissionDenied     = errors.New("权限不足")
    ErrCorruptedArchive     = errors.New("压缩文件已损坏")
    ErrFormatMismatch       = errors.New("文件格式与扩展名不匹配")
)
```

#### 错误处理原则
1. 使用现有的colorlib进行错误输出
2. 提供详细的错误信息和建议
3. 支持错误恢复和重试机制
4. 记录详细的错误日志
5. 格式检测失败时提供建议

## 🔧 架构集成方案

### 1. 主命令调度器集成

#### 修改 `commands/cmd.go`
```go
// 在 Run() 函数中添加compress、extract和list-archive子命令的初始化
compressCmd := compress.InitCompressCmd()
extractCmd := extract.InitExtractCmd()
listArchiveCmd := listarchive.InitListArchiveCmd()

// 添加到子命令列表
if addCmdErr := qflag.AddSubCmd(sizeCmd, listCmd, checkCmd, hashCmd, findCmd, compressCmd, extractCmd, listArchiveCmd); addCmdErr != nil {
    // 错误处理
}

// 在switch语句中添加处理逻辑
case compressCmd.LongName(), compressCmd.ShortName():
    if err := compress.CompressCmdMain(cmdCL); err != nil {
        fmt.Printf("err: %v\n", err)
        os.Exit(1)
    }
case extractCmd.LongName(), extractCmd.ShortName():
    if err := extract.ExtractCmdMain(cmdCL); err != nil {
        fmt.Printf("err: %v\n", err)
        os.Exit(1)
    }
case listArchiveCmd.LongName(), listArchiveCmd.ShortName():
    if err := listarchive.ListArchiveCmdMain(cmdCL); err != nil {
        fmt.Printf("err: %v\n", err)
        os.Exit(1)
    }
```

### 2. 配置管理集成

#### 复用现有配置模式
- 使用qflag库进行参数解析
- 遵循现有的标志命名约定
- 集成到现有的帮助系统
- 支持配置文件（未来扩展）

### 3. 颜色输出集成

#### 使用现有colorlib
```go
// 成功信息
cl.PrintOkf("压缩完成: %s (%s -> %s, 压缩比: %.1f%%)\n", 
    archiveName, originalSize, compressedSize, ratio)

// 警告信息
cl.PrintWarnf("跳过文件: %s (权限不足)\n", filename)

// 错误信息
cl.PrintErrorf("压缩失败: %s\n", err.Error())

// 进度信息
cl.PrintInfof("正在处理: %s [%d/%d]\n", filename, current, total)

// 格式检测信息
cl.PrintInfof("检测到格式: %s\n", formatName)
```

### 4. 通用工具复用

#### 使用现有common包功能
- 文件路径处理和验证
- 错误处理工具函数
- 进度条显示（复用progressbar）
- 文件大小格式化
- 权限检查工具

## 🎯 高级功能扩展

### 未来可考虑的功能

#### 1. 格式转换功能
```bash
# 格式转换子命令
fck convert backup.zip backup.7z    # ZIP转7Z
fck convert backup.tar.gz backup.zip # TAR.GZ转ZIP
```

#### 2. 压缩分析功能
```bash
# 分析压缩效果
fck analyze backup.zip
# 输出：文件类型分布、压缩比统计、重复文件检测等
```

#### 3. 批量操作功能
```bash
# 批量压缩
fck compress-batch -f zip *.txt     # 将每个txt文件单独压缩
fck compress-batch -f tgz ./*/      # 将每个子目录单独压缩

# 批量解压
fck extract-batch *.zip             # 解压所有zip文件
```

#### 4. 云存储集成
```bash
# 直接压缩到云存储
fck compress --upload s3://bucket/backup.zip ./data

# 从云存储解压
fck extract --download s3://bucket/backup.zip
```

#### 5. 增量压缩功能
```bash
# 增量压缩（基于时间戳）
fck compress --incremental --since "2024-01-01" backup.zip ./data

# 增量压缩（基于基准压缩包）
fck compress --incremental --base old_backup.zip new_backup.zip ./data
```

#### 6. 压缩包管理功能
```bash
# 压缩包信息
fck info backup.zip                 # 显示详细信息

# 压缩包修复
fck repair corrupted.zip            # 尝试修复损坏的压缩包

# 压缩包合并
fck merge output.zip part1.zip part2.zip  # 合并多个压缩包
```

## 📊 与现有子命令协同

### 1. 与find子命令结合
```bash
# 查找并压缩特定文件
fck find -n "*.log" -exec "fck compress logs.zip {}"

# 查找大文件并分别压缩
fck find -size +100MB -exec "fck compress {}.7z {}"

# 查找并排除压缩
fck find -type f | grep -v "\.tmp$" | xargs fck compress clean.tar.gz
```

### 2. 与size子命令结合
```bash
# 压缩前后大小对比
fck size ./data
fck compress data.7z ./data
fck size data.7z

# 分析压缩效果
echo "原始大小: $(fck size --total ./data)"
echo "压缩后大小: $(fck size data.7z)"
```

### 3. 与hash子命令结合
```bash
# 压缩后验证完整性
fck compress backup.7z ./important_data
fck hash backup.7z

# 解压后验证完整性
fck extract backup.7z ./restore
fck hash ./restore

# 压缩包内容哈希验证
fck extract -l backup.zip | fck hash --stdin
```

### 4. 与list子命令结合
```bash
# 比较目录和压缩包内容
fck list ./source --format csv > source.csv
fck extract -l backup.zip --format csv > backup.csv
diff source.csv backup.csv

# 压缩包内容详细列表
fck list-archive backup.zip --format table
```

### 5. 与check子命令结合
```bash
# 检查压缩包完整性
fck check backup.zip backup.7z

# 批量检查压缩包
fck find -n "*.zip" -exec "fck check {}"
```

## 🧪 测试策略

### 单元测试覆盖
- **格式检测测试**: 测试各种文件格式的正确识别
- **压缩功能测试**: 测试各种格式的压缩功能
- **解压功能测试**: 测试各种格式的解压功能
- **参数验证测试**: 测试所有参数的有效性验证
- **错误处理测试**: 测试各种错误情况的处理
- **边界条件测试**: 测试极限情况（空文件、大文件等）

### 集成测试覆盖
- **跨格式测试**: 测试不同格式间的兼容性
- **与其他子命令协同测试**: 测试命令组合使用
- **并发操作测试**: 测试多线程压缩/解压
- **大文件处理测试**: 测试GB级文件的处理
- **跨平台兼容性测试**: Windows/Linux/macOS测试

### 性能测试
- **压缩速度基准测试**: 不同格式的压缩速度对比
- **解压速度基准测试**: 不同格式的解压速度对比
- **内存使用量测试**: 监控内存占用情况
- **并发性能测试**: 多线程处理性能测试
- **大文件性能测试**: GB级文件处理性能

### 测试用例示例
```go
func TestFormatDetection(t *testing.T) {
    tests := []struct {
        filename string
        expected Format
    }{
        {"test.zip", FormatZIP},
        {"test.tar.gz", FormatTARGZ},
        {"test.tgz", FormatTARGZ},
        {"test.7z", Format7Z},
        {"test.rar", FormatRAR},
    }
    
    detector := NewFormatDetector()
    for _, tt := range tests {
        format, err := detector.DetectByExtension(tt.filename)
        assert.NoError(t, err)
        assert.Equal(t, tt.expected, format)
    }
}

func TestCompressExtract(t *testing.T) {
    formats := []Format{FormatZIP, FormatTARGZ, Format7Z}
    
    for _, format := range formats {
        t.Run(format.String(), func(t *testing.T) {
            // 创建测试数据
            testDir := createTestData(t)
            defer os.RemoveAll(testDir)
            
            // 压缩
            archivePath := filepath.Join(t.TempDir(), "test"+format.Extension())
            err := CompressDirectory(testDir, archivePath, format, CompressOptions{})
            assert.NoError(t, err)
            
            // 解压
            extractDir := t.TempDir()
            err = ExtractArchive(archivePath, extractDir, ExtractOptions{})
            assert.NoError(t, err)
            
            // 验证内容
            assert.True(t, compareDirectories(testDir, extractDir))
        })
    }
}
```

## 📈 性能优化策略

### 压缩优化
1. **并行压缩**: 支持多线程并行处理多个文件
2. **内存缓冲区优化**: 根据文件大小动态调整缓冲区
3. **压缩算法选择**: 根据文件类型选择最优算法
4. **预处理优化**: 文件预分析和排序优化

### 解压优化
1. **并行解压**: 支持多线程并行解压
2. **流式处理**: 大文件流式解压，减少内存占用
3. **索引缓存**: 压缩包索引缓存，加速重复操作
4. **磁盘I/O优化**: 批量写入，减少磁盘碎片

### 内存优化
1. **分块处理**: 大文件分块处理，控制内存使用
2. **对象池**: 复用缓冲区对象，减少GC压力
3. **延迟加载**: 按需加载压缩包内容
4. **内存映射**: 大文件使用内存映射技术

## 🔒 安全考虑

### 安全措施
1. **路径遍历防护**: 防止../等路径遍历攻击
2. **ZIP炸弹检测**: 检测恶意构造的压缩包
3. **文件大小限制**: 限制单个文件和总解压大小
4. **权限验证**: 验证文件读写权限
5. **符号链接检查**: 防止符号链接攻击

### 密码安全
1. **安全密码输入**: 隐藏密码输入，防止肩窥
2. **密码强度验证**: 提供密码强度建议
3. **内存清理**: 使用后立即清理内存中的密码
4. **加密算法**: 使用强加密算法（AES-256等）

### 安全实现示例
```go
func ValidateExtractPath(basePath, targetPath string) error {
    // 清理路径
    cleanTarget := filepath.Clean(targetPath)
    cleanBase := filepath.Clean(basePath)
    
    // 检查是否在基础路径内
    if !strings.HasPrefix(cleanTarget, cleanBase) {
        return fmt.Errorf("路径遍历攻击检测: %s", targetPath)
    }
    
    // 检查路径长度
    if len(cleanTarget) > 4096 {
        return fmt.Errorf("路径过长: %s", targetPath)
    }
    
    return nil
}

func DetectZipBomb(archive *zip.Reader) error {
    var totalUncompressed int64
    var totalCompressed int64
    
    for _, file := range archive.File {
        totalUncompressed += int64(file.UncompressedSize64)
        totalCompressed += int64(file.CompressedSize64)
        
        // 检查单个文件大小
        if file.UncompressedSize64 > MaxSingleFileSize {
            return fmt.Errorf("文件过大: %s (%d bytes)", file.Name, file.UncompressedSize64)
        }
    }
    
    // 检查总大小
    if totalUncompressed > MaxTotalUncompressedSize {
        return fmt.Errorf("解压后总大小过大: %d bytes", totalUncompressed)
    }
    
    // 检查压缩比
    if totalCompressed > 0 && totalUncompressed/totalCompressed > MaxCompressionRatio {
        return fmt.Errorf("疑似ZIP炸弹，压缩比异常: %d", totalUncompressed/totalCompressed)
    }
    
    return nil
}
```

## 📋 开发计划

### 第一阶段：基础功能实现（2-3周）
- [ ] 格式检测器实现
- [ ] ZIP格式压缩/解压支持
- [ ] TAR格式压缩/解压支持
- [ ] TAR.GZ格式压缩/解压支持
- [ ] 基础参数解析和验证
- [ ] 错误处理机制
- [ ] 基础测试用例

### 第二阶段：功能增强（2-3周）
- [ ] 进度条显示
- [ ] 密码保护支持（ZIP格式）
- [ ] 文件过滤功能（包含/排除模式）
- [ ] 详细输出模式
- [ ] 并行处理优化
- [ ] TAR.BZ2和TAR.XZ格式支持

### 第三阶段：高级功能（3-4周）
- [ ] 7Z格式支持
- [ ] RAR格式解压支持
- [ ] 分卷压缩支持
- [ ] list-archive子命令
- [ ] 与其他子命令集成
- [ ] 性能优化和基准测试

### 第四阶段：扩展功能（4-6周）
- [ ] 格式转换功能
- [ ] 批量操作支持
- [ ] 增量压缩功能
- [ ] 云存储集成
- [ ] Web界面（可选）
- [ ] 插件系统（可选）

### 第五阶段：完善和发布（1-2周）
- [ ] 全面测试和bug修复
- [ ] 文档完善
- [ ] 性能调优
- [ ] 发布准备

## 📚 参考资料和标准

### 文件格式规范
- [ZIP文件格式规范](https://pkware.cachefly.net/webdocs/casestudies/APPNOTE.TXT)
- [TAR文件格式规范](https://www.gnu.org/software/tar/manual/html_node/Standard.html)
- [7Z文件格式规范](https://www.7-zip.org/7z.html)
- [RAR文件格式规范](https://www.rarlab.com/technote.htm)

### Go语言相关
- [Go archive/zip 包文档](https://pkg.go.dev/archive/zip)
- [Go archive/tar 包文档](https://pkg.go.dev/archive/tar)
- [Go compress 包文档](https://pkg.go.dev/compress)

### 项目相关
- [FCK工具现有架构](./commands/)
- [qflag库使用指南](https://gitee.com/MM-Q/qflag)
- [colorlib库使用指南](https://gitee.com/MM-Q/colorlib)

### 安全参考
- [OWASP文件上传安全指南](https://owasp.org/www-community/vulnerabilities/Unrestricted_File_Upload)
- [ZIP炸弹防护指南](https://blog.ostorlab.co/zip-packages-exploitation.html)

## 🎯 成功指标

### 功能指标
- ✅ 支持6种以上压缩格式
- ✅ 压缩/解压成功率 > 99.9%
- ✅ 与现有子命令100%兼容
- ✅ 测试覆盖率 > 90%

### 性能指标
- ✅ 压缩速度不低于系统原生工具
- ✅ 内存使用量 < 100MB（处理1GB文件）
- ✅ 并发处理提升效率 > 50%
- ✅ 启动时间 < 100ms

### 用户体验指标
- ✅ 命令学习成本 < 5分钟
- ✅ 错误信息清晰易懂
- ✅ 进度显示准确及时
- ✅ 帮助文档完整

---

**文档版本**: v2.0  
**创建日期**: 2025-08-17  
**最后更新**: 2025-08-17  
**作者**: CodeBuddy  
**状态**: 设计阶段  
**变更**: 从ZIP/UNZIP单一格式设计升级为多格式统一设计
