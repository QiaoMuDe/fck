# Package comprx

Package comprx 提供了一个统一的压缩和解压缩库，支持多种压缩格式，包括 ZIP、TAR、GZIP、BZIP2、ZLIB 和 TGZ。它还支持进度条显示、文件过滤、并发安全操作等高级功能。

## 主要功能

- 压缩和解压缩多种格式的文件
- 支持进度条显示
- 文件过滤功能
- 线程安全操作
- 灵活的配置选项

## 基本使用示例

```go
// 简单压缩
err := comprx.Pack("output.zip", "input_dir")

// 简单解压
err := comprx.Unpack("archive.zip", "output_dir")

// 带进度条的压缩
err := comprx.PackProgress("output.zip", "input_dir")
```

## 文件过滤功能

Package comprx 提供文件过滤功能，支持从忽略文件加载排除模式。

### 主要功能

- 从忽略文件加载排除模式
- 支持注释行和空行处理
- 自动去重排除模式
- 支持 glob 模式匹配
- 提供文件不存在时的容错处理

### 支持的忽略文件格式

- 每行一个模式
- `#` 开头的注释行
- 空行自动忽略
- 支持标准 glob 通配符

### 使用示例

```go
// 加载忽略文件，文件不存在会报错
patterns, err := comprx.LoadExcludeFromFile(".gitignore")

// 加载忽略文件，文件不存在返回空列表
patterns, err := comprx.LoadExcludeFromFileOrEmpty(".comprxignore")
```

## 压缩包内容列表和信息查看功能

Package comprx 提供查看压缩包内容的各种方法，包括列出文件信息、打印压缩包信息等。支持多种压缩格式，提供简洁和详细两种显示样式，支持文件过滤和数量限制。

### 主要功能

- 列出压缩包内的文件信息
- 打印压缩包基本信息
- 支持文件名模式匹配
- 支持限制显示文件数量
- 提供简洁和详细两种显示样式

### 使用示例

```go
// 列出压缩包内所有文件
info, err := comprx.List("archive.zip")

// 打印压缩包信息（简洁样式）
err := comprx.PrintLs("archive.zip")

// 打印匹配模式的文件（详细样式）
err := comprx.PrintLlMatch("archive.zip", "*.go")
```

## 内存中的压缩和解压缩功能

Package comprx 提供 GZIP 和 ZLIB 格式的内存压缩和流式压缩功能。支持字节数组、字符串和流式数据的压缩与解压缩操作。

### 主要功能

- GZIP 内存压缩：字节数组和字符串的压缩解压
- GZIP 流式压缩：支持 `io.Reader` 和 `io.Writer` 接口
- ZLIB 内存压缩：字节数组和字符串的压缩解压
- ZLIB 流式压缩：支持 `io.Reader` 和 `io.Writer` 接口
- 支持自定义压缩等级

### 使用示例

```go
// GZIP 压缩字符串
compressed, err := comprx.GzipString("hello world")

// ZLIB 解压字节数据
decompressed, err := comprx.UnzlibBytes(compressedData)
```

## 压缩和解压缩操作的配置选项

Package comprx 定义了 `Options` 结构体和相关的配置方法，用于控制压缩和解压缩操作的行为。支持压缩等级设置、进度条显示、文件过滤、路径验证等功能的配置。

### 主要类型

- `Options`: 压缩/解压配置选项结构体

### 主要功能

- 提供默认配置选项
- 支持链式配置方法
- 提供各种预设配置选项

### 使用示例

```go
opts := comprx.Options{
    CompressionLevel: config.CompressionLevelBest,
    OverwriteExisting: true,
    ProgressEnabled: true,
    ProgressStyle: types.ProgressStyleUnicode,
}
err := comprx.PackOptions("output.zip", "input_dir", opts)
```

## 文件和目录大小计算功能

Package comprx 提供计算文件或目录大小的实用函数。支持单个文件大小获取和目录递归大小计算，提供安全和详细两种版本。

### 主要功能

- 获取单个文件的大小
- 递归计算目录的总大小
- 提供安全版本（出错返回 0）和详细版本（返回错误信息）
- 自动忽略符号链接等特殊文件

### 使用示例

```go
// 安全版本，出错时返回 0
size := comprx.GetSizeOrZero("./mydir")

// 详细版本，返回错误信息
size, err := comprx.GetSize("./myfile.txt")
```

## FUNCTIONS

### GetSize

```go
func GetSize(path string) (int64, error)
```

- **描述**: 获取文件或目录的大小（字节）
- **参数**:
  - `path`: 文件或目录路径
- **返回**:
  - `int64`: 文件或目录的总大小（字节）
  - `error`: 错误信息
- **注意**:
  - 如果是文件，返回文件大小
  - 如果是目录，返回目录中所有文件的总大小
  - 如果路径不存在，返回错误
  - 只计算普通文件的大小，忽略符号链接等特殊文件

### GetSizeOrZero

```go
func GetSizeOrZero(path string) int64
```

- **描述**: 获取文件或目录的大小，出错时返回 0
- **参数**:
  - `path`: 文件或目录路径
- **返回**:
  - `int64`: 文件或目录的总大小（字节），出错时返回 0
- **功能**:
  - 如果是文件，返回文件大小
  - 如果是目录，返回目录中所有普通文件的总大小
  - 忽略符号链接等特殊文件
  - 发生任何错误时返回 0，不抛出异常
- **注意**:
  - 此函数为 `GetSize` 的安全版本，适用于不需要错误处理的场景
  - 如需详细错误信息，请使用 `GetSize` 函数

### GzipBytes

```go
func GzipBytes(data []byte) ([]byte, error)
```

- **描述**: 压缩字节数据（使用默认压缩等级）
- **参数**:
  - `data`: 要压缩的字节数据
- **返回**:
  - `[]byte`: 压缩后的数据
  - `error`: 错误信息
- **使用示例**:

```go
compressed, err := GzipBytes([]byte("hello world"))
```

### GzipBytesWithLevel

```go
func GzipBytesWithLevel(data []byte, level types.CompressionLevel) ([]byte, error)
```

- **描述**: 压缩字节数据（指定压缩等级）
- **参数**:
  - `data`: 要压缩的字节数据
  - `level`: 压缩级别
- **返回**:
  - `[]byte`: 压缩后的数据
  - `error`: 错误信息
- **使用示例**:

```go
compressed, err := GzipBytesWithLevel([]byte("hello world"), types.CompressionLevelBest)
```

### GzipStream

```go
func GzipStream(dst io.Writer, src io.Reader) error
```

- **描述**: 流式压缩数据（使用默认压缩等级）
- **参数**:
  - `dst`: 目标写入器
  - `src`: 源读取器
- **返回**:
  - `error`: 错误信息
- **使用示例**:

```go
file, _ := os.Open("input.txt")
defer file.Close()

var buf bytes.Buffer
err := GzipStream(&buf, file)
```

### GzipStreamWithLevel

```go
func GzipStreamWithLevel(dst io.Writer, src io.Reader, level types.CompressionLevel) error
```

- **描述**: 流式压缩数据（指定压缩等级）
- **参数**:
  - `dst`: 目标写入器
  - `src`: 源读取器
  - `level`: 压缩级别
- **返回**:
  - `error`: 错误信息
- **使用示例**:

```go
file, _ := os.Open("input.txt")
defer file.Close()

output, _ := os.Create("output.gz")
defer output.Close()

err := GzipStreamWithLevel(output, file, types.CompressionLevelBest)
```

### GzipString

```go
func GzipString(text string) ([]byte, error)
```

- **描述**: 压缩字符串（使用默认压缩等级）
- **参数**:
  - `text`: 要压缩的字符串
- **返回**:
  - `[]byte`: 压缩后的数据
  - `error`: 错误信息
- **使用示例**:

```go
compressed, err := GzipString("hello world")
```

### GzipStringWithLevel

```go
func GzipStringWithLevel(text string, level types.CompressionLevel) ([]byte, error)
```

- **描述**: 压缩字符串（指定压缩等级）
- **参数**:
  - `text`: 要压缩的字符串
  - `level`: 压缩级别
- **返回**:
  - `[]byte`: 压缩后的数据
  - `error`: 错误信息
- **使用示例**:

```go
compressed, err := GzipStringWithLevel("hello world", types.CompressionLevelBest)
```

### List

```go
func List(archivePath string) (*types.ArchiveInfo, error)
```

- **描述**: 列出压缩包的所有文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
- **返回**:
  - `*types.ArchiveInfo`: 压缩包信息
  - `error`: 错误信息

### ListLimit

```go
func ListLimit(archivePath string, limit int) (*types.ArchiveInfo, error)
```

- **描述**: 列出指定数量的文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `limit`: 限制返回的文件数量
- **返回**:
  - `*types.ArchiveInfo`: 压缩包信息
  - `error`: 错误信息

### ListMatch

```go
func ListMatch(archivePath string, pattern string) (*types.ArchiveInfo, error)
```

- **描述**: 列出匹配指定模式的文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `pattern`: 文件名匹配模式 (支持通配符 `*` 和 `?`)
- **返回**:
  - `*types.ArchiveInfo`: 压缩包信息
  - `error`: 错误信息

### LoadExcludeFromFile

```go
func LoadExcludeFromFile(ignoreFilePath string) ([]string, error)
```

- **描述**: 从忽略文件加载排除模式
- **参数**:
  - `ignoreFilePath`: 忽略文件路径（如 `.comprxignore`, `.gitignore`）
- **返回**:
  - `[]string`: 排除模式列表（已去重）
  - `error`: 错误信息
- **支持的文件格式**:
  - 每行一个模式
  - 支持 `#` 开头的注释行
  - 自动忽略空行
  - 支持 glob 模式匹配
  - 自动去除重复模式
- **使用示例**:

```go
patterns, err := comprx.LoadExcludeFromFile(".comprxignore")
```

### LoadExcludeFromFileOrEmpty

```go
func LoadExcludeFromFileOrEmpty(ignoreFilePath string) ([]string, error)
```

- **描述**: 从忽略文件加载排除模式，文件不存在时返回空列表
- **参数**:
  - `ignoreFilePath`: 忽略文件路径
- **返回**:
  - `[]string`: 排除模式列表，文件不存在时返回空列表
  - `error`: 错误信息（文件不存在不算错误）
- **使用示例**:

```go
patterns, err := comprx.LoadExcludeFromFileOrEmpty(".comprxignore")
```

### Pack

```go
func Pack(dst string, src string) error
```

- **描述**: 压缩文件或目录（禁用进度条） - 线程安全
- **参数**:
  - `dst`: 目标文件路径
  - `src`: 源文件路径
- **返回**:
  - `error`: 错误信息
- **使用示例**:

```go
err := Pack("output.zip", "input_dir")
```

### PackOptions

```go
func PackOptions(dst string, src string, opts Options) error
```

- **描述**: 使用指定配置压缩文件或目录 - 线程安全
- **参数**:
  - `dst`: 目标文件路径
  - `src`: 源文件路径
  - `opts`: 配置选项
- **返回**:
  - `error`: 错误信息
- **使用示例**:

```go
opts := Options{
    CompressionLevel: config.CompressionLevelBest,
    OverwriteExisting: true,
    ProgressEnabled: true,
    ProgressStyle: types.ProgressStyleUnicode,
}
err := PackOptions("output.zip", "input_dir", opts)
```

### PackProgress

```go
func PackProgress(dst string, src string) error
```

- **描述**: 压缩文件或目录（启用进度条） - 线程安全
- **参数**:
  - `dst`: 目标文件路径
  - `src`: 源文件路径
- **返回**:
  - `error`: 错误信息
- **使用示例**:

```go
err := PackProgress("output.zip", "input_dir")
```

### PrintArchiveAndFiles

```go
func PrintArchiveAndFiles(archivePath string, detailed bool) error
```

- **描述**: 打印压缩包信息和所有文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `detailed`: `true`=详细样式, `false`=简洁样式(默认)
- **返回**:
  - `error`: 错误信息

### PrintArchiveAndFilesLimit

```go
func PrintArchiveAndFilesLimit(archivePath string, limit int, detailed bool) error
```

- **描述**: 打印压缩包信息和指定数量的文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `limit`: 限制打印的文件数量
  - `detailed`: `true`=详细样式, `false`=简洁样式(默认)
- **返回**:
  - `error`: 错误信息

### PrintArchiveAndFilesMatch

```go
func PrintArchiveAndFilesMatch(archivePath string, pattern string, detailed bool) error
```

- **描述**: 打印压缩包信息和匹配指定模式的文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `pattern`: 文件名匹配模式 (支持通配符 `*` 和 `?`)
  - `detailed`: `true`=详细样式, `false`=简洁样式(默认)
- **返回**:
  - `error`: 错误信息

### PrintArchiveInfo

```go
func PrintArchiveInfo(archivePath string) error
```

- **描述**: 打印压缩包本身的基本信息
- **参数**:
  - `archivePath`: 压缩包文件路径
- **返回**:
  - `error`: 错误信息

### PrintFiles

```go
func PrintFiles(archivePath string, detailed bool) error
```

- **描述**: 打印压缩包内所有文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `detailed`: `true`=详细样式, `false`=简洁样式(默认)
- **返回**:
  - `error`: 错误信息

### PrintFilesLimit

```go
func PrintFilesLimit(archivePath string, limit int, detailed bool) error
```

- **描述**: 打印压缩包内指定数量的文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `limit`: 限制打印的文件数量
  - `detailed`: `true`=详细样式, `false`=简洁样式(默认)
- **返回**:
  - `error`: 错误信息

### PrintFilesMatch

```go
func PrintFilesMatch(archivePath string, pattern string, detailed bool) error
```

- **描述**: 打印压缩包内匹配指定模式的文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `pattern`: 文件名匹配模式 (支持通配符 `*` 和 `?`)
  - `detailed`: `true`=详细样式, `false`=简洁样式(默认)
- **返回**:
  - `error`: 错误信息

### PrintInfo

```go
func PrintInfo(archivePath string) error
```

- **描述**: 打印压缩包信息和所有文件信息（简洁样式）
- **参数**:
  - `archivePath`: 压缩包文件路径
- **返回**:
  - `error`: 错误信息

### PrintInfoDetailed

```go
func PrintInfoDetailed(archivePath string) error
```

- **描述**: 打印压缩包信息和所有文件信息（详细样式）
- **参数**:
  - `archivePath`: 压缩包文件路径
- **返回**:
  - `error`: 错误信息

### PrintInfoDetailedLimit

```go
func PrintInfoDetailedLimit(archivePath string, limit int) error
```

- **描述**: 打印压缩包信息和指定数量的文件信息（详细样式）
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `limit`: 限制打印的文件数量
- **返回**:
  - `error`: 错误信息

### PrintInfoDetailedMatch

```go
func PrintInfoDetailedMatch(archivePath string, pattern string) error
```

- **描述**: 打印压缩包信息和匹配指定模式的文件信息（详细样式）
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `pattern`: 文件名匹配模式 (支持通配符 `*` 和 `?`)
- **返回**:
  - `error`: 错误信息

### PrintInfoLimit

```go
func PrintInfoLimit(archivePath string, limit int) error
```

- **描述**: 打印压缩包信息和指定数量的文件信息（简洁样式）
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `limit`: 限制打印的文件数量
- **返回**:
  - `error`: 错误信息

### PrintInfoMatch

```go
func PrintInfoMatch(archivePath string, pattern string) error
```

- **描述**: 打印压缩包信息和匹配指定模式的文件信息（简洁样式）
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `pattern`: 文件名匹配模式 (支持通配符 `*` 和 `?`)
- **返回**:
  - `error`: 错误信息

### PrintLl

```go
func PrintLl(archivePath string) error
```

- **描述**: 打印压缩包内所有文件信息（详细样式）
- **参数**:
  - `archivePath`: 压缩包文件路径
- **返回**:
  - `error`: 错误信息

### PrintLlLimit

```go
func PrintLlLimit(archivePath string, limit int) error
```

- **描述**: 打印压缩包内指定数量的文件信息（详细样式）
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `limit`: 限制打印的文件数量
- **返回**:
  - `error`: 错误信息

### PrintLlMatch

```go
func PrintLlMatch(archivePath string, pattern string) error
```

- **描述**: 打印压缩包内匹配指定模式的文件信息（详细样式）
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `pattern`: 文件名匹配模式 (支持通配符 `*` 和 `?`)
- **返回**:
  - `error`: 错误信息

### PrintLs

```go
func PrintLs(archivePath string) error
```

- **描述**: 打印压缩包内所有文件信息（简洁样式）
- **参数**:
  - `archivePath`: 压缩包文件路径
- **返回**:
  - `error`: 错误信息

### PrintLsLimit

```go
func PrintLsLimit(archivePath string, limit int) error
```

- **描述**: 打印压缩包内指定数量的文件信息（简洁样式）
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `limit`: 限制打印的文件数量
- **返回**:
  - `error`: 错误信息

### PrintLsMatch

```go
func PrintLsMatch(archivePath string, pattern string) error
```

- **描述**: 打印压缩包内匹配指定模式的文件信息（简洁样式）
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `pattern`: 文件名匹配模式 (支持通配符 `*` 和 `?`)
- **返回**:
  - `error`: 错误信息

### UngzipBytes

```go
func UngzipBytes(compressedData []byte) ([]byte, error)
```

- **描述**: 解压字节数据
- **参数**:
  - `compressedData`: 压缩的字节数据
- **返回**:
  - `[]byte`: 解压后的数据
  - `error`: 错误信息
- **使用示例**:

```go
decompressed, err := UngzipBytes(compressedData)
```

### UngzipStream

```go
func UngzipStream(dst io.Writer, src io.Reader) error
```

- **描述**: 流式解压数据
- **参数**:
  - `dst`: 目标写入器
  - `src`: 源读取器（压缩数据）
- **返回**:
  - `error`: 错误信息
- **使用示例**:

```go
compressedFile, _ := os.Open("input.gz")
defer compressedFile.Close()

output, _ := os.Create("output.txt")
defer output.Close()

err := UngzipStream(output, compressedFile)
```

### UngzipString

```go
func UngzipString(compressedData []byte) (string, error)
```

- **描述**: 解压为字符串
- **参数**:
  - `compressedData`: 压缩的字节数据
- **返回**:
  - `string`: 解压后的字符串
  - `error`: 错误信息
- **使用示例**:

```go
text, err := UngzipString(compressedData)
```

### Unpack

```go
func Unpack(src string, dst string) error
```

- **描述**: 解压文件（禁用进度条） - 线程安全
- **参数**:
  - `src`: 源文件路径
  - `dst`: 目标目录路径
- **返回**:
  - `error`: 错误信息
- **使用示例**:

```go
err := Unpack("archive.zip", "output_dir")
```

### UnpackDir

```go
func UnpackDir(archivePath string, dirName string, outputDir string) error
```

- **描述**: 解压指定目录 - 线程安全
- **参数**:
  - `archivePath`: 压缩包路径
  - `dirName`: 要解压的目录名
  - `outputDir`: 输出目录路径
- **返回**:
  - `error`: 错误信息
- **使用示例**:

```go
err := UnpackDir("archive.zip", "src", "output/")
```

### UnpackFile

```go
func UnpackFile(archivePath string, fileName string, outputDir string) error
```

- **描述**: 解压指定文件名 - 线程安全
- **参数**:
  - `archivePath`: 压缩包路径
  - `fileName`: 要解压的文件名
  - `outputDir`: 输出目录路径
- **返回**:
  - `error`: 错误信息
- **使用示例**:

```go
err := UnpackFile("archive.zip", "config.json", "output/")
```

### UnpackMatch

```go
func UnpackMatch(archivePath string, keyword string, outputDir string) error
```

- **描述**: 解压匹配关键字的文件 - 线程安全
- **参数**:
  - `archivePath`: 压缩包路径
  - `keyword`: 匹配关键字
  - `outputDir`: 输出目录路径
- **返回**:
  - `error`: 错误信息
- **使用示例**:

```go
err := UnpackMatch("archive.zip", "test", "output/")
```

### UnpackOptions

```go
func UnpackOptions(src string, dst string, opts Options) error
```

- **描述**: 使用指定配置解压文件 - 线程安全
- **参数**:
  - `src`: 源文件路径
  - `dst`: 目标目录路径
  - `opts`: 配置选项
- **返回**:
  - `error`: 错误信息
- **使用示例**:

```go
opts := Options{
    OverwriteExisting: true,
    ProgressEnabled: true,
    ProgressStyle: types.ProgressStyleASCII,
}
err := UnpackOptions("archive.zip", "output_dir", opts)
```

### UnpackProgress

```go
func UnpackProgress(src string, dst string) error
```

- **描述**: 解压文件（启用进度条） - 线程安全
- **参数**:
  - `src`: 源文件路径
  - `dst`: 目标目录路径
- **返回**:
  - `error`: 错误信息
- **使用示例**:

```go
err := UnpackProgress("archive.zip", "output_dir")
```

### UnzlibBytes

```go
func UnzlibBytes(compressedData []byte) ([]byte, error)
```

- **描述**: 解压字节数据
- **参数**:
  - `compressedData`: 压缩的字节数据
- **返回**:
  - `[]byte`: 解压后的数据
  - `error`: 错误信息
- **使用示例**:

```go
decompressed, err := UnzlibBytes(compressedData)
```

### UnzlibStream

```go
func UnzlibStream(dst io.Writer, src io.Reader) error
```

- **描述**: 流式解压数据
- **参数**:
  - `dst`: 目标写入器
  - `src`: 源读取器（压缩数据）
- **返回**:
  - `error`: 错误信息
- **使用示例**:

```go
compressedFile, _ := os.Open("input.zlib")
defer compressedFile.Close()

output, _ := os.Create("output.txt")
defer output.Close()

err := UnzlibStream(output, compressedFile)
```

### UnzlibString

```go
func UnzlibString(compressedData []byte) (string, error)
```

- **描述**: 解压为字符串
- **参数**:
  - `compressedData`: 压缩的字节数据
- **返回**:
  - `string`: 解压后的字符串
  - `error`: 错误信息
- **使用示例**:

```go
text, err := UnzlibString(compressedData)
```

### ZlibBytes

```go
func ZlibBytes(data []byte) ([]byte, error)
```

- **描述**: 压缩字节数据（使用默认压缩等级）
- **参数**:
  - `data`: 要压缩的字节数据
- **返回**:
  - `[]byte`: 压缩后的数据
  - `error`: 错误信息
- **使用示例**:

```go
compressed, err := ZlibBytes([]byte("hello world"))
```

### ZlibBytesWithLevel

```go
func ZlibBytesWithLevel(data []byte, level types.CompressionLevel) ([]byte, error)
```

- **描述**: 压缩字节数据（指定压缩等级）
- **参数**:
  - `data`: 要压缩的字节数据
  - `level`: 压缩级别
- **返回**:
  - `[]byte`: 压缩后的数据
  - `error`: 错误信息
- **使用示例**:

```go
compressed, err := ZlibBytesWithLevel([]byte("hello world"), types.CompressionLevelBest)
```

### ZlibStream

```go
func ZlibStream(dst io.Writer, src io.Reader) error
```

- **描述**: 流式压缩数据（使用默认压缩等级）
- **参数**:
  - `dst`: 目标写入器
  - `src`: 源读取器
- **返回**:
  - `error`: 错误信息
- **使用示例**:

```go
file, _ := os.Open("input.txt")
defer file.Close()

var buf bytes.Buffer
err := ZlibStream(&buf, file)
```

### ZlibStreamWithLevel

```go
func ZlibStreamWithLevel(dst io.Writer, src io.Reader, level types.CompressionLevel) error
```

- **描述**: 流式压缩数据（指定压缩等级）
- **参数**:
  - `dst`: 目标写入器
  - `src`: 源读取器
  - `level`: 压缩级别
- **返回**:
  - `error`: 错误信息
- **使用示例**:

```go
file, _ := os.Open("input.txt")
defer file.Close()

output, _ := os.Create("output.zlib")
defer output.Close()

err := ZlibStreamWithLevel(output, file, types.CompressionLevelBest)
```

### ZlibString

```go
func ZlibString(text string) ([]byte, error)
```

- **描述**: 压缩字符串（使用默认压缩等级）
- **参数**:
  - `text`: 要压缩的字符串
- **返回**:
  - `[]byte`: 压缩后的数据
  - `error`: 错误信息
- **使用示例**:

```go
compressed, err := ZlibString("hello world")
```

### ZlibStringWithLevel

```go
func ZlibStringWithLevel(text string, level types.CompressionLevel) ([]byte, error)
```

- **描述**: 压缩字符串（指定压缩等级）
- **参数**:
  - `text`: 要压缩的字符串
  - `level`: 压缩级别
- **返回**:
  - `[]byte`: 压缩后的数据
  - `error`: 错误信息
- **使用示例**:

```go
compressed, err := ZlibStringWithLevel("hello world", types.CompressionLevelBest)
```

## TYPES

### Options

```go
type Options struct {
    CompressionLevel      types.CompressionLevel // 压缩等级
    OverwriteExisting     bool                   // 是否覆盖已存在的文件
    ProgressEnabled       bool                   // 是否启用进度显示
    ProgressStyle         types.ProgressStyle    // 进度条样式
    DisablePathValidation bool                   // 是否禁用路径验证
    Filter                types.FilterOptions    // 过滤选项
}
```

- **描述**: 压缩/解压配置选项

### ASCIIProgressOptions

```go
func ASCIIProgressOptions() Options
```

- **描述**: 返回ASCII样式进度条配置选项
- **返回**:
  - `Options`: ASCII样式进度条配置选项
- **使用示例**:

```go
err := PackOptions("output.zip", "input_dir", ASCIIProgressOptions())
```

### DefaultOptions

```go
func DefaultOptions() Options
```

- **描述**: 返回默认配置选项
- **返回**:
  - `Options`: 默认配置选项
- **默认配置**:
  - `CompressionLevel`: 默认压缩等级
  - `OverwriteExisting`: `false` (不覆盖已存在文件)
  - `ProgressEnabled`: `false` (不显示进度)
  - `ProgressStyle`: 文本样式
  - `DisablePathValidation`: `false` (启用路径验证)

### DefaultProgressOptions

```go
func DefaultProgressOptions() Options
```

- **描述**: 返回默认样式进度条配置选项
- **返回**:
  - `Options`: 默认样式进度条配置选项
- **使用示例**:

```go
err := PackOptions("output.zip", "input_dir", DefaultProgressOptions())
```

### ForceOptions

```go
func ForceOptions() Options
```

- **描述**: 返回强制模式配置选项
- **返回**:
  - `Options`: 强制模式配置选项
- **配置特点**:
  - `OverwriteExisting`: `true` (覆盖已存在文件)
  - `DisablePathValidation`: `true` (禁用路径验证)
  - `ProgressEnabled`: `false` (关闭进度条)
- **使用示例**:

```go
err := PackOptions("output.zip", "input_dir", ForceOptions())
```

### NoCompressionOptions

```go
func NoCompressionOptions() Options
```

- **描述**: 返回禁用压缩且启用进度条的配置选项
- **返回**:
  - `Options`: 禁用压缩且启用进度条的配置选项
- **配置特点**:
  - `CompressionLevel`: 无压缩 (存储模式)
  - `ProgressEnabled`: `true` (启用进度条)
  - `ProgressStyle`: 文本样式
- **使用示例**:

```go
err := PackOptions("output.zip", "input_dir", NoCompressionOptions())
```

### NoCompressionProgressOptions

```go
func NoCompressionProgressOptions(style types.ProgressStyle) Options
```

- **描述**: 返回禁用压缩且启用指定样式进度条的配置选项
- **参数**:
  - `style`: 进度条样式
- **返回**:
  - `Options`: 禁用压缩且启用指定样式进度条的配置选项
- **配置特点**:
  - `CompressionLevel`: 无压缩 (存储模式)
  - `ProgressEnabled`: `true` (启用进度条)
  - `ProgressStyle`: 指定样式
- **使用示例**:

```go
err := PackOptions("output.zip", "input_dir", NoCompressionProgressOptions(types.ProgressStyleUnicode))
```

### ProgressOptions

```go
func ProgressOptions(style types.ProgressStyle) Options
```

- **描述**: 返回带进度显示的配置选项
- **参数**:
  - `style`: 进度条样式
- **返回**:
  - `Options`: 带进度显示的配置选项

### TextProgressOptions

```go
func TextProgressOptions() Options
```

- **描述**: 返回文本样式进度条配置选项
- **返回**:
  - `Options`: 文本样式进度条配置选项
- **使用示例**:

```go
err := PackOptions("output.zip", "input_dir", TextProgressOptions())
```

### UnicodeProgressOptions

```go
func UnicodeProgressOptions() Options
```

- **描述**: 返回Unicode样式进度条配置选项
- **返回**:
  - `Options`: Unicode样式进度条配置选项
- **使用示例**:

```go
err := PackOptions("output.zip", "input_dir", UnicodeProgressOptions())
```

### Options 方法

#### SetCompressionLevel

```go
func (o *Options) SetCompressionLevel(level types.CompressionLevel)
```

- **描述**: 设置压缩等级
- **参数**:
  - `level`: 压缩等级
- **使用示例**:

```go
opts := DefaultOptions()
opts.SetCompressionLevel(types.CompressionLevelBest)
```

#### SetDisablePathValidation

```go
func (o *Options) SetDisablePathValidation(disable bool)
```

- **描述**: 设置是否禁用路径验证
- **参数**:
  - `disable`: 是否禁用路径验证
- **使用示例**:

```go
opts := DefaultOptions()
opts.SetDisablePathValidation(true)
```

#### SetExclude

```go
func (o *Options) SetExclude(patterns []string)
```

- **描述**: 设置排除模式
- **参数**:
  - `patterns`: 排除模式列表
- **使用示例**:

```go
opts := DefaultOptions()
opts.SetExclude([]string{"*_test.go", "vendor/*"})
```

#### SetFilter

```go
func (o *Options) SetFilter(filter types.FilterOptions)
```

- **描述**: 设置过滤配置
- **参数**:
  - `filter`: 过滤选项
- **使用示例**:

```go
opts := DefaultOptions()
filter := types.FilterOptions{
    Include: []string{"*.go", "*.md"},
    Exclude: []string{"*_test.go"},
}
opts.SetFilter(filter)
```

#### SetInclude

```go
func (o *Options) SetInclude(patterns []string)
```

- **描述**: 设置包含模式
- **参数**:
  - `patterns`: 包含模式列表
- **使用示例**:

```go
opts := DefaultOptions()
opts.SetInclude([]string{"*.go", "*.md"})
```

#### SetMaxSize

```go
func (o *Options) SetMaxSize(maxSize int64)
```

- **描述**: 设置最大文件大小
- **参数**:
  - `maxSize`: 最大文件大小（字节）
- **使用示例**:

```go
opts := DefaultOptions()
opts.SetMaxSize(10 * 1024 * 1024) // 10MB
```

#### SetMinSize

```go
func (o *Options) SetMinSize(minSize int64)
```

- **描述**: 设置最小文件大小
- **参数**:
  - `minSize`: 最小文件大小（字节）
- **使用示例**:

```go
opts := DefaultOptions()
opts.SetMinSize(1024) // 1KB
```

#### SetOverwriteExisting

```go
func (o *Options) SetOverwriteExisting(overwrite bool)
```

- **描述**: 设置是否覆盖已存在的文件
- **参数**:
  - `overwrite`: 是否覆盖已存在文件
- **使用示例**:

```go
opts := DefaultOptions()
opts.SetOverwriteExisting(true)
```

#### SetProgress

```go
func (o *Options) SetProgress(enabled bool)
```

- **描述**: 设置是否启用进度显示
- **参数**:
  - `enabled`: 是否启用进度显示
- **使用示例**:

```go
opts := DefaultOptions()
opts.SetProgress(true)
```

#### SetProgressAndStyle

```go
func (o *Options) SetProgressAndStyle(enabled bool, style types.ProgressStyle)
```

- **描述**: 设置进度显示和样式
- **参数**:
  - `enabled`: 是否启用进度显示
  - `style`: 进度条样式
- **使用示例**:

```go
opts := DefaultOptions()
opts.SetProgressAndStyle(true, types.ProgressStyleUnicode)
```

#### SetProgressStyle

```go
func (o *Options) SetProgressStyle(style types.ProgressStyle)
```

- **描述**: 设置进度条样式
- **参数**:
  - `style`: 进度条样式
- **使用示例**:

```go
opts := DefaultOptions()
opts.SetProgressStyle(types.ProgressStyleUnicode)
```

#### SetSizeFilter

```go
func (o *Options) SetSizeFilter(minSize, maxSize int64)
```

- **描述**: 设置文件大小过滤
- **参数**:
  - `minSize`: 最小文件大小（字节）
  - `maxSize`: 最大文件大小（字节）
- **使用示例**:

```go
opts := DefaultOptions()
opts.SetSizeFilter(1024, 10*1024*1024) // 1KB - 10MB
```

### Options 链式调用方法

#### WithCompressionLevel

```go
func (o Options) WithCompressionLevel(level types.CompressionLevel) Options
```

- **描述**: 设置压缩等级
- **参数**:
  - `level`: 压缩等级
- **返回**:
  - `Options`: 配置选项（支持链式调用）
- **使用示例**:

```go
opts := DefaultOptions().WithCompressionLevel(types.CompressionLevelBest)
```

#### WithDisablePathValidation

```go
func (o Options) WithDisablePathValidation(disable bool) Options
```

- **描述**: 设置是否禁用路径验证
- **参数**:
  - `disable`: 是否禁用路径验证
- **返回**:
  - `Options`: 配置选项（支持链式调用）
- **使用示例**:

```go
opts := DefaultOptions().WithDisablePathValidation(true)
```

#### WithExclude

```go
func (o Options) WithExclude(patterns []string) Options
```

- **描述**: 设置排除模式
- **参数**:
  - `patterns`: 排除模式列表
- **返回**:
  - `Options`: 配置选项（支持链式调用）
- **使用示例**:

```go
opts := DefaultOptions().WithExclude([]string{"*_test.go", "vendor/*"})
```

#### WithFilter

```go
func (o Options) WithFilter(filter types.FilterOptions) Options
```

- **描述**: 设置过滤配置
- **参数**:
  - `filter`: 过滤选项
- **返回**:
  - `Options`: 配置选项（支持链式调用）
- **使用示例**:

```go
filter := types.FilterOptions{
    Include: []string{"*.go", "*.md"},
    Exclude: []string{"*_test.go"},
}
opts := DefaultOptions().WithFilter(filter)
```

#### WithInclude

```go
func (o Options) WithInclude(patterns []string) Options
```

- **描述**: 设置包含模式
- **参数**:
  - `patterns`: 包含模式列表
- **返回**:
  - `Options`: 配置选项（支持链式调用）
- **使用示例**:

```go
opts := DefaultOptions().WithInclude([]string{"*.go", "*.md"})
```

#### WithMaxSize

```go
func (o Options) WithMaxSize(maxSize int64) Options
```

- **描述**: 设置最大文件大小
- **参数**:
  - `maxSize`: 最大文件大小（字节）
- **返回**:
  - `Options`: 配置选项（支持链式调用）
- **使用示例**:

```go
opts := DefaultOptions().WithMaxSize(10 * 1024 * 1024) // 10MB
```

#### WithMinSize

```go
func (o Options) WithMinSize(minSize int64) Options
```

- **描述**: 设置最小文件大小
- **参数**:
  - `minSize`: 最小文件大小（字节）
- **返回**:
  - `Options`: 配置选项（支持链式调用）
- **使用示例**:

```go
opts := DefaultOptions().WithMinSize(1024) // 1KB
```

#### WithOverwriteExisting

```go
func (o Options) WithOverwriteExisting(overwrite bool) Options
```

- **描述**: 设置是否覆盖已存在的文件
- **参数**:
  - `overwrite`: 是否覆盖已存在文件
- **返回**:
  - `Options`: 配置选项（支持链式调用）
- **使用示例**:

```go
opts := DefaultOptions().WithOverwriteExisting(true)
```

#### WithProgress

```go
func (o Options) WithProgress(enabled bool) Options
```

- **描述**: 设置是否启用进度显示
- **参数**:
  - `enabled`: 是否启用进度显示
- **返回**:
  - `Options`: 配置选项（支持链式调用）
- **使用示例**:

```go
opts := DefaultOptions().WithProgress(true)
```

#### WithProgressAndStyle

```go
func (o Options) WithProgressAndStyle(enabled bool, style types.ProgressStyle) Options
```

- **描述**: 设置进度显示和样式
- **参数**:
  - `enabled`: 是否启用进度显示
  - `style`: 进度条样式
- **返回**:
  - `Options`: 配置选项（支持链式调用）
- **使用示例**:

```go
opts := DefaultOptions().WithProgressAndStyle(true, types.ProgressStyleUnicode)
```

#### WithProgressStyle

```go
func (o Options) WithProgressStyle(style types.ProgressStyle) Options
```

- **描述**: 设置进度条样式
- **参数**:
  - `style`: 进度条样式
- **返回**:
  - `Options`: 配置选项（支持链式调用）
- **使用示例**:

```go
opts := DefaultOptions().WithProgressStyle(types.ProgressStyleUnicode)
```

#### WithSizeFilter

```go
func (o Options) WithSizeFilter(minSize, maxSize int64) Options
```

- **描述**: 设置文件大小过滤
- **参数**:
  - `minSize`: 最小文件大小（字节）
  - `maxSize`: 最大文件大小（字节）
- **返回**:
  - `Options`: 配置选项（支持链式调用）
- **使用示例**:

```go
opts := DefaultOptions().WithSizeFilter(1024, 10*1024*1024) // 1KB - 10MB
```