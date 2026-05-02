# types.go 拆分方案

## 目标
将 `internal/types/types.go` (505行) 拆分为3个文件，提高代码可维护性。

## 拆分方案

### 文件结构

```
internal/types/
├── types.go          # 基础类型 + 数据结构 + Windows扩展名
├── format.go         # 表格样式 + 语法高亮常量
└── command.go        # 查找 + DNS + TCP 命令相关类型
```

### 内容分配

#### 1. types.go
- 基础常量：
  - `OutputFileName`
  - `OutputCheckFileName`
  - `TimestampFormat`
  - `VirtualRootDir`
  - `ChecksumModePortable` / `ChecksumModeLocal`
  - `InitialBufferSize` / `DefaultMaxBufferSize`
  - `SedBackupSuffix` / `SedTempFilePattern`
- 数据结构：
  - `VirtualHashEntry`
  - `VirtualHashMap`
  - `DirEntryWrapper`
- Windows扩展名：
  - `WindowsExecutableExts`
  - `WindowsSymlinkExts`

#### 2. format.go
- 语法高亮常量：
  - `HighlightFormatter256`
  - `HighlightFormatter16m`
  - `HighlightFormatter16`
  - `HighlightStyleDefault`
- 表格样式：
  - `TableStyleMap`
  - `TableStyles`
  - `StyleNone`

#### 3. command.go
- 查找类型常量：
  - `FindTypeAll`, `FindTypeFile`, `FindTypeDir`, ...
  - `FindTypeFileShort`, `FindTypeDirShort`, ...
- 查找类型变量：
  - `FindTypeLimits`
  - `ListTypeLimits`
  - `FindLimits`
- 查找类型函数：
  - `IsValidFindType()`
  - `GetSupportedFindTypes()`
- DNS类型：
  - `DNSTypeA`, `DNSTypeAAAA`, ...
  - `DNSQueryTypes`
  - `GetDNSAnyTypes()`
- TCP扫描类型：
  - `TCPScanFormatDefault`, `TCPScanFormatTable`, ...
  - `TCPScanFormatOptions`

## 实施步骤

1. 创建 `format.go`，迁移高亮和表格相关内容
2. 创建 `command.go`，迁移查找、DNS、TCP相关内容
3. 修改 `types.go`，删除已迁移的内容
4. 编译测试，确保所有引用正常

## 注意事项

- 所有文件保持 `package types` 包名
- 导入语句按需保留
- 函数注释一并迁移
