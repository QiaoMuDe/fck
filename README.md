<a name="top"></a>
<div align="center">

# 🚀 FCK - 全功能命令行工具集

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/QiaoMuDe/fck)
[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-GPL%20v3.0-green.svg)](LICENSE)

**一站式文件操作、文本处理、系统管理与网络诊断工具集**

[功能特性](#-功能特性) • [快速开始](#-快速开始) • [使用示例](#-使用示例) • [贡献指南](#-贡献指南)

</div>

---

## 简介

FCK 是一个功能丰富的跨平台命令行工具集，集成了 **40+** 实用工具，涵盖文件操作、文本处理、系统监控、网络诊断、压缩解压等多个领域。无论是日常文件管理、开发辅助工作，还是系统运维、网络排障，FCK 都能提供高效便捷的解决方案。

### 为什么选择 FCK？

- **一站式解决方案**：覆盖文件、文本、系统、网络等多个场景，一个工具搞定多种需求
- **跨平台支持**：Windows、Linux、macOS 全平台兼容，统一的使用体验
- **现代化设计**：彩色输出、表格样式、进度显示、交互式体验
- **高性能**：并发处理、流式操作、内存优化，轻松应对大文件和批量任务
- **管道友好**：所有命令都支持标准输入输出，便于组合使用和脚本集成
- **持续更新**：活跃的开发和维护，不断添加新功能和优化体验

---

## ✨ 功能特性

### � 文件操作
文件校验、打包解包、查找搜索、复制移动、目录列表、压缩预览等全套文件管理工具。

### � 文本处理
类 Unix 文本处理工具集：grep 搜索、sed 替换、awk 字段提取、wc 统计、tr 字符转换、head/tail 查看等。

### �️ 系统监控
进程查看、端口监控、磁盘空间、文件大小统计、命令监控等系统管理工具。

### 🌐 网络工具
TCP 连接测试（支持交互模式）、DNS 查询（多记录类型）、Ping 测试、路由追踪等网络诊断工具。

### �️ 开发辅助
JSON 处理、Base64 编解码、哈希计算、序列生成等开发常用工具。

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

### 查看帮助
```bash
# 查看所有可用命令
fck --help

# 查看具体命令帮助
fck [命令] --help
```

---

## � 使用示例

### 文件校验
```bash
# 计算文件 MD5
fck hash file.txt

# 递归计算目录所有文件的 SHA256
fck hash -r ./mydir

# 验证校验文件
fck check checksum.hash
```

### 文本处理
```bash
# 搜索包含 "error" 的行
fck grep "error" log.txt

# 替换文本
fck sed -s "old/new/g" file.txt

# 提取第 1、3 列
fck awk -f 1,3 data.csv

# 查看文件前 20 行
fck head -n 20 large.log
```

### 系统监控
```bash
# 查看进程列表
fck proc

# 查看端口占用
fck port

# 监控磁盘空间
fck df

# 周期性执行命令（每 2 秒）
fck watch -n 2 "fck proc"
```

### 网络工具
```bash
# DNS 查询
fck dns www.example.com

# 查询所有记录类型
fck dns -a www.example.com

# TCP 连接测试（交互模式）
fck tcp -i host:port

# Ping 测试
fck ping 8.8.8.8

# 路由追踪
fck tracepath www.example.com
```

### 文件打包
```bash
# 打包目录
fck pack -o backup.zip ./mydir

# 解压文件
fck unpack archive.zip -d ./output

# 预览压缩包内容
fck preview archive.zip
```

---

## 🎨 表格样式

FCK 支持 20+ 种表格样式，让输出更美观：

| 样式代码 | 说明 | 样式代码 | 说明 |
|---------|------|---------|------|
| `def` | 默认样式 | `l` | 浅色样式 |
| `r` | 圆角样式 | `bd` | 粗体样式 |
| `cb` | 彩色亮色样式 | `cd` | 彩色暗色样式 |
| `db` | 双线样式 | `none` | 禁用样式 |
| `cbb` | 黑色背景蓝色字体 | `cbc` | 青色背景蓝色字体 |
| `cbg` | 绿色背景蓝色字体 | `cbm` | 紫色背景蓝色字体 |
| `cby` | 黄色背景蓝色字体 | `cbr` | 红色背景蓝色字体 |
| `cwb` | 蓝色背景白色字体 | `ccw` | 青色背景白色字体 |
| `cgw` | 绿色背景白色字体 | `cmw` | 紫色背景白色字体 |
| `crw` | 红色背景白色字体 | `cyw` | 黄色背景白色字体 |

使用示例：
```bash
fck list -S r          # 使用圆角样式
fck proc -S cbg        # 使用绿色背景样式
```

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

[报告问题](https://gitee.com/MM-Q/fck/issues) • [功能建议](https://gitee.com/MM-Q/fck/issues) • [讨论交流](https://gitee.com/MM-Q/fck/discussions) • <a href="#top">返回顶部</a>

</div>
