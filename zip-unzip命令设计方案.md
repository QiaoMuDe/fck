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
├── zip/
│   ├── cmd_zip.go          # zip子命令主逻辑
│   ├── flags.go            # 标志定义
│   ├── compressor.go       # 压缩核心逻辑
│   ├── validator.go        # 参数验证
│   ├── progress.go         # 进度跟踪
│   ├── APIDOC.md          # API文档
│   └── cmd_zip_test.go     # 测试文件
└── unzip/
    ├── cmd_unzip.go        # unzip子命令主逻辑
    ├── flags.go            # 标志定义
    ├── extractor.go        # 解压核心逻辑
    ├── validator.go        # 参数验证
    ├── progress.go         # 进度跟踪
    ├── APIDOC.md          # API文档
    └── cmd_unzip_test.go   # 测试文件
```

### 核心接口设计

#### 1. 压缩器接口
```go
type Compressor interface {
    // 添加文件到压缩包
    AddFile(path string, info os.FileInfo) error
    
    // 添加目录到压缩包
    AddDirectory(path string) error
    
    // 设置密码保护
    SetPassword(password string)
    
    // 设置压缩级别
    SetCompressionLevel(level int)
    
    // 设置压缩方法
    SetCompressionMethod(method string)
    
    // 添加注释
    SetComment(comment string)
    
    // 关闭压缩器
    Close() error
}
```

#### 2. 解压器接口
```go
type Extractor interface {
    // 列出ZIP文件内容
    List() ([]FileInfo, error)
    
    // 解压所有文件
    Extract(destination string) error
    
    // 解压指定文件
    ExtractFile(filename, destination string) error
    
    // 测试ZIP文件完整性
    Test() error
    
    // 设置密码
    SetPassword(password string)
    
    // 关闭解压器
    Close() error
}
```

#### 3. 进度跟踪器
```go
type ProgressTracker struct {
    TotalFiles     int64     // 总文件数
    ProcessedFiles int64     // 已处理文件数
    TotalBytes     int64     // 总字节数
    ProcessedBytes int64     // 已处理字节数
    StartTime      time.Time // 开始时间
    CurrentFile    string    // 当前处理文件
}

func (p *ProgressTracker) Update(filename string, bytes int64)
func (p *ProgressTracker) GetProgress() float64
func (p *ProgressTracker) GetETA() time.Duration
func (p *ProgressTracker) GetSpeed() int64
```

#### 4. 文件信息结构
```go
type FileInfo struct {
    Name         string    // 文件名
    Size         int64     // 文件大小
    CompressedSize int64   // 压缩后大小
    ModTime      time.Time // 修改时间
    IsDir        bool      // 是否为目录
    Mode         os.FileMode // 文件权限
    CRC32        uint32    // CRC32校验值
}
```

### 依赖库选择

#### 标准库
- `archive/zip` - Go标准ZIP库
- `compress/flate` - 压缩算法
- `path/filepath` - 路径处理
- `os` - 文件系统操作

#### 第三方库（可选）
- `github.com/alexmullins/zip` - 支持密码保护的ZIP库
- `github.com/schollz/progressbar/v3` - 进度条（已在项目中使用）

### 错误处理策略

#### 错误类型定义
```go
var (
    ErrInvalidArchive    = errors.New("无效的ZIP文件")
    ErrPasswordRequired  = errors.New("需要密码")
    ErrWrongPassword     = errors.New("密码错误")
    ErrFileExists        = errors.New("文件已存在")
    ErrInsufficientSpace = errors.New("磁盘空间不足")
    ErrPermissionDenied  = errors.New("权限不足")
)
```

#### 错误处理原则
1. 使用现有的colorlib进行错误输出
2. 提供详细的错误信息和建议
3. 支持错误恢复和重试机制
4. 记录详细的错误日志

## 🔧 架构集成方案

### 1. 主命令调度器集成

#### 修改 `commands/cmd.go`
```go
// 在 Run() 函数中添加zip和unzip子命令的初始化
zipCmd := zip.InitZipCmd()
unzipCmd := unzip.InitUnzipCmd()

// 添加到子命令列表
if addCmdErr := qflag.AddSubCmd(sizeCmd, listCmd, checkCmd, hashCmd, findCmd, zipCmd, unzipCmd); addCmdErr != nil {
    // 错误处理
}

// 在switch语句中添加处理逻辑
case zipCmd.LongName(), zipCmd.ShortName():
    if err := zip.ZipCmdMain(cmdCL); err != nil {
        fmt.Printf("err: %v\n", err)
        os.Exit(1)
    }
case unzipCmd.LongName(), unzipCmd.ShortName():
    if err := unzip.UnzipCmdMain(cmdCL); err != nil {
        fmt.Printf("err: %v\n", err)
        os.Exit(1)
    }
```

### 2. 配置管理集成

#### 复用现有配置模式
- 使用qflag库进行参数解析
- 遵循现有的标志命名约定
- 集成到现有的帮助系统

### 3. 颜色输出集成

#### 使用现有colorlib
```go
// 成功信息
cl.PrintOkf("压缩完成: %s (%s -> %s)\n", archiveName, originalSize, compressedSize)

// 警告信息
cl.PrintWarnf("跳过文件: %s (权限不足)\n", filename)

// 错误信息
cl.PrintErrorf("压缩失败: %s\n", err.Error())

// 进度信息
cl.PrintInfof("正在压缩: %s\n", filename)
```

### 4. 通用工具复用

#### 使用现有common包功能
- 文件路径处理
- 错误处理工具
- 进度条显示
- 文件大小格式化

## 🎯 高级功能扩展

### 未来可考虑的功能

#### 1. 多格式支持
- tar.gz格式支持
- 7z格式支持
- rar格式支持（仅解压）

#### 2. 云存储集成
- 直接压缩到云存储
- 从云存储解压文件
- 支持AWS S3、阿里云OSS等

#### 3. 增量压缩
- 基于时间戳的增量压缩
- 基于文件哈希的变更检测
- 压缩包版本管理

#### 4. 压缩分析
- 压缩比统计
- 文件类型分析
- 压缩效果报告

#### 5. 批量操作
- 批量压缩多个目录
- 批量解压多个ZIP文件
- 压缩任务队列管理

## 📊 与现有子命令协同

### 1. 与find子命令结合
```bash
# 查找并压缩特定文件
fck find -n "*.log" -exec "fck zip logs.zip {}"

# 查找大文件并分别压缩
fck find -size +100MB -exec "fck zip {}.zip {}"
```

### 2. 与size子命令结合
```bash
# 压缩前后大小对比
fck size ./data
fck zip data.zip ./data
fck size data.zip

# 分析压缩效果
fck size --compare-with data.zip ./data
```

### 3. 与hash子命令结合
```bash
# 压缩后验证完整性
fck zip backup.zip ./important_data
fck hash backup.zip

# 解压后验证完整性
fck unzip backup.zip ./restore
fck hash ./restore
```

### 4. 与list子命令结合
```bash
# 列出ZIP内容（详细格式）
fck unzip -l archive.zip | fck list --format table

# 比较目录和ZIP内容
fck list ./source > source.list
fck unzip -l backup.zip > backup.list
diff source.list backup.list
```

## 🧪 测试策略

### 单元测试
- 压缩/解压核心功能测试
- 参数验证测试
- 错误处理测试
- 边界条件测试

### 集成测试
- 与其他子命令的协同测试
- 大文件处理测试
- 并发操作测试
- 跨平台兼容性测试

### 性能测试
- 压缩速度基准测试
- 内存使用量测试
- 大文件处理性能测试
- 并发压缩性能测试

## 📈 性能优化

### 压缩优化
- 并行压缩支持
- 内存缓冲区优化
- 压缩算法选择
- 文件预处理优化

### 解压优化
- 并行解压支持
- 流式解压处理
- 内存映射文件
- 磁盘I/O优化

## 🔒 安全考虑

### 安全措施
- 路径遍历攻击防护
- ZIP炸弹检测
- 文件大小限制
- 权限验证

### 密码安全
- 安全的密码输入
- 密码强度验证
- 内存中密码清理
- 加密算法选择

## 📋 开发计划

### 第一阶段（基础功能）
- [ ] 基本压缩/解压功能
- [ ] 参数解析和验证
- [ ] 错误处理机制
- [ ] 基础测试用例

### 第二阶段（功能增强）
- [ ] 进度条显示
- [ ] 密码保护支持
- [ ] 文件过滤功能
- [ ] 详细输出模式

### 第三阶段（高级功能）
- [ ] 分卷压缩支持
- [ ] 并行处理优化
- [ ] 与其他子命令集成
- [ ] 性能优化

### 第四阶段（扩展功能）
- [ ] 多格式支持
- [ ] 云存储集成
- [ ] 增量压缩
- [ ] 批量操作

## 📚 参考资料

- [Go archive/zip 包文档](https://pkg.go.dev/archive/zip)
- [ZIP文件格式规范](https://pkware.cachefly.net/webdocs/casestudies/APPNOTE.TXT)
- [FCK工具现有架构分析](./commands/)
- [qflag库使用指南](https://gitee.com/MM-Q/qflag)
- [colorlib库使用指南](https://gitee.com/MM-Q/colorlib)

---

**文档版本**: v1.0  
**创建日期**: 2025-08-17  
**最后更新**: 2025-08-17  
**作者**: CodeBuddy  
**状态**: 设计阶段