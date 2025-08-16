<div align="center">

# 🚀 FCK - 文件系统命令行工具集

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/QiaoMuDe/fck)
[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-GPL%20v3.0-green.svg)](LICENSE)

**一个强大、高效、现代化的文件系统操作工具集**

[功能特性](#-功能特性) • [快速开始](#-快速开始) • [命令详解](#-命令详解) • [使用示例](#-使用示例) • [贡献指南](#-贡献指南)

</div>

---

## ✨ 功能特性

### 🔐 文件校验 (hash)
- **多算法支持**: MD5、SHA1、SHA256、SHA512
- **批量处理**: 支持通配符和递归扫描
- **完整性验证**: 生成和验证校验文件
- **高性能**: 并发计算，显著提升处理速度

### 📊 智能统计 (size)
- **精确计算**: 文件和目录大小统计
- **人性化显示**: 自动选择最佳单位 (B/KB/MB/GB/TB)
- **进度显示**: 大目录扫描时显示实时进度
- **隐藏文件**: 可选择包含或排除隐藏文件

### 🔍 高级查找 (find)
- **多条件筛选**: 按名称、大小、时间、类型等组合查找
- **正则支持**: 强大的正则表达式匹配
- **并发搜索**: 多线程并行处理，提升搜索效率
- **批量操作**: 支持删除、移动、执行命令等批量操作

### 📋 目录列表 (list)
- **多种排序**: 按名称、大小、时间排序
- **彩色显示**: 根据文件类型智能着色
- **表格样式**: 20+种表格样式可选
- **详细信息**: 显示权限、用户组、修改时间等

### 🔄 差异对比 (diff)
- **目录比较**: 基于哈希值的精确比较
- **完整性检查**: 根据校验文件验证目录完整性
- **详细报告**: 生成完整的差异分析报告
- **批量验证**: 支持多目录批量对比

---

## 🚀 快速开始

### 环境要求
- Go 1.24+ 
- 支持 Windows、Linux、macOS

### 安装方式

#### 方式一：源码编译
```bash
# 克隆仓库
git clone https://gitee.com/MM-Q/fck.git
cd fck

# 开发版本（生成到 output 目录）
python3 build.py

# 正式版本（安装到 $GOPATH/bin）
python3 build.py -s -ai -f

# 发布版本（压缩包）
python3 build.py -batch -z
```

#### 方式二：直接运行
```bash
go run main.go [命令] [选项]
```

---

## 📖 命令详解

### 🔐 hash - 文件哈希计算

计算文件或目录的哈希值，支持完整性验证。

```bash
# 基本用法
fck hash [文件/目录...] [选项]

# 常用选项
-t, --type string     哈希算法 (md5|sha1|sha256|sha512) [默认: sha256]
-r, --recursive       递归处理目录
-w, --write          将结果写入 checksum.hash 文件
```

**使用示例：**
```bash
# 计算单个文件的SHA256值
fck hash document.pdf

# 递归计算目录所有文件的MD5值
fck hash ./photos -t md5 -r

# 生成校验文件
fck hash ./important_files -r -w
```

### 📊 size - 文件大小统计

统计文件或目录的磁盘占用空间。

```bash
# 基本用法
fck size [路径...] [选项]

# 常用选项
-H, --hidden         包含隐藏文件
-c, --color          启用彩色输出
-ts, --table-style   表格样式 [默认: default]
```

**使用示例：**
```bash
# 统计当前目录各项大小
fck size

# 统计指定目录，包含隐藏文件
fck size /home/user -H

# 彩色表格显示
fck size ./projects -c -ts cb
```

### 🔍 find - 高级文件查找

强大的文件搜索工具，支持多条件组合查找。

```bash
# 基本用法
fck find [路径] [选项]

# 搜索选项
-n, --name string        按文件名搜索
-p, --path string        按路径搜索
-type string            文件类型 (f|d|l|h|e|x等)
-size string            文件大小 (+100M|-50K)
-mtime string           修改时间 (+7|-30)

# 高级选项
-r, --regex             启用正则表达式
-i, --case              区分大小写
-w, --whole-word        全词匹配
-x, --concurrent        并发搜索
-H, --hidden            包含隐藏文件

# 操作选项
-delete                 删除匹配的文件
-mv string              移动到指定目录
-exec string            对匹配文件执行命令
-count                  只显示匹配数量
```

**使用示例：**
```bash
# 查找所有PDF文件
fck find . -n "*.pdf"

# 查找大于100MB的文件
fck find /home -size +100M -type f

# 查找7天内修改的图片文件
fck find ./photos -mtime -7 -n "*.jpg|*.png" -r

# 删除空文件
fck find . -type e -delete

# 并发搜索日志文件
fck find /var/log -n "*.log" -x -H
```

### 📋 list - 目录内容列表

增强版的目录列表工具，提供丰富的显示选项。

```bash
# 基本用法
fck list [路径...] [选项]

# 显示选项
-l, --long              详细信息显示
-a, --all               显示隐藏文件
-r, --recursive         递归显示
-c, --color             彩色输出

# 排序选项
-s, --sort-size         按大小排序
-t, --sort-time         按时间排序
-n, --sort-name         按名称排序
-R, --reverse           反向排序

# 样式选项
-ts, --table-style      表格样式
-q, --quote             文件名加引号
-u, --show-user-group   显示用户和组
```

**使用示例：**
```bash
# 详细列表显示
fck list -l -c

# 按大小排序，包含隐藏文件
fck list -s -a

# 递归显示，按时间倒序
fck list -r -t -R

# 使用圆角表格样式
fck list -l -ts r -c
```

### 🔄 diff - 目录差异对比

比较目录或验证文件完整性。

```bash
# 基本用法
fck diff [选项]

# 比较模式
-a, --dir-a string      目录A路径
-b, --dir-b string      目录B路径
-f, --file string       校验文件路径
-d, --dirs string       要验证的目录

# 输出选项
-w, --write             将结果写入文件
-t, --type string       哈希算法 [默认: sha256]
```

**使用示例：**
```bash
# 比较两个目录
fck diff -a ./backup -b ./current

# 根据校验文件验证目录完整性
fck diff -f checksum.hash

# 比较校验文件与指定目录
fck diff -f backup.hash -d ./restore

# 生成比较报告
fck diff -a ./old -b ./new -w
```

---

## 💡 使用示例

### 场景一：备份验证
```bash
# 1. 为重要文件生成校验文件
fck hash ./important_docs -r -w -t sha256

# 2. 备份文件到其他位置
cp -r ./important_docs ./backup/

# 3. 验证备份完整性
fck diff -f checksum.hash -d ./backup/important_docs
```

### 场景二：磁盘清理
```bash
# 1. 找出大文件
fck find /home -size +1G -type f

# 2. 查找重复文件（通过大小）
fck list /home -r -s | grep "同样大小"

# 3. 清理临时文件
fck find /tmp -name "*.tmp" -mtime +7 -delete
```

### 场景三：项目管理
```bash
# 1. 统计项目各模块大小
fck size ./src/* -c -ts cb

# 2. 查找未使用的文件
fck find ./assets -mtime +30 -type f

# 3. 生成项目文件清单
fck list ./src -r -l > project_files.txt
```

---

## 🎨 表格样式

FCK 支持 20+ 种表格样式，让输出更美观：

| 样式代码 | 说明 | 样式代码 | 说明 |
|---------|------|---------|------|
| `default` | 默认样式 | `l` | 浅色样式 |
| `r` | 圆角样式 | `bd` | 粗体样式 |
| `cb` | 彩色亮色 | `cd` | 彩色暗色 |
| `db` | 双线样式 | `none` | 无边框 |

使用示例：
```bash
fck list -c -ts r    # 圆角彩色表格
fck size -c -ts cb   # 彩色亮色表格
```

---

## 🔧 高级配置

### 性能优化
- 使用 `-x` 选项启用并发处理
- 大目录操作时合理设置缓冲区大小
- 根据系统资源调整并发数量

### 安全考虑
- 路径遍历保护：自动检测和阻止危险路径
- 权限检查：操作前验证文件访问权限
- 输入验证：严格验证所有用户输入

---

## 🤝 贡献指南

我们欢迎所有形式的贡献！

### 如何贡献
1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

### 开发环境
```bash
# 克隆仓库
git clone https://gitee.com/MM-Q/fck.git
cd fck

# 安装依赖
go mod tidy

# 运行测试
go test ./...

# 构建开发版本
python3 build.py
```

### 代码规范
- 遵循 Go 官方代码规范
- 添加必要的单元测试
- 更新相关文档

---

## 📄 许可证

本项目采用 [GNU General Public License v3.0](LICENSE) 许可证。

---

## 🙏 致谢

感谢所有贡献者和使用者的支持！

---

<div align="center">

**如果这个项目对你有帮助，请给个 ⭐ Star！**

[报告问题](https://gitee.com/MM-Q/fck/issues) • [功能建议](https://gitee.com/MM-Q/fck/issues) • [讨论交流](https://gitee.com/MM-Q/fck/discussions)

</div>