<a name="top"></a>
<div align="center">

# 🚀 FCK - 全功能命令行工具集

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/QiaoMuDe/fck) [![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org) [![License](https://img.shields.io/badge/License-GPL%20v3.0-green.svg)](LICENSE) [![Platform](https://img.shields.io/badge/Platform-Windows%2CLinux%2CmacOS-green.svg)]()


**一站式文件操作、文本处理、系统管理与网络诊断工具集**

[功能特性](#-功能特性) • [快速开始](#-快速开始) • [使用示例](#-使用示例) • [贡献指南](#-贡献指南)

</div>

---

## 简介

FCK 是一个功能丰富的跨平台命令行工具集，集成了 **46+** 实用工具，涵盖文件操作、文本处理、系统监控、网络诊断、压缩解压等多个领域。无论是日常文件管理、开发辅助工作，还是系统运维、网络排障，FCK 都能提供高效便捷的解决方案。

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

### 🛠️ 开发辅助
JSON 处理、Base64 编解码、哈希计算、编码转换、序列生成等开发常用工具。

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

#### 方式二：通过Go安装
```bash
go install gitee.com/MM-Q/fck/cmd/fck@latest
```

---

## 📚 命令帮助

执行 `fck --help` 可查看所有可用命令及全局选项：

```bash
👉  fck --help


    ________      ________          ___  __       
   |\  _____\    |\   ____\        |\  \|\  \     
   \ \  \__/     \ \  \___|        \ \  \/  /|_   
    \ \   __\     \ \  \            \ \   ___  \  
     \ \  \_|      \ \  \____        \ \  \\ \  \ 
      \ \__\        \ \_______\       \ \__\\ \__\
       \|__|         \|_______|        \|__| \|__|
                   FCK CLI

名称:
  fck.exe

描述:
  一站式文件操作、文本处理与系统管理工具集, 集成46+实用命令, 覆盖文件管理、文本处理、系统监控、开发辅助等多个场景

用法:
  fck.exe [options] [args...]

选项:
  -h, --help <bool>                显示帮助信息 (default: false)
  -v, --version <bool>             显示版本信息 (default: false)
  --completion <enum>              生成Shell自动补全脚本, 支持的Shell: [bash pwsh powershell] (default: pwsh)
  --install-completion <enum>      安装Shell自动补全脚本到系统, 支持的Shell: [bash pwsh powershell] (default: pwsh)

子命令:
  base64, b64       Base64 编解码工具
  check, ck         校验指定文件的哈希值是否与校验文件中的记录一致
  curl, c           HTTP 客户端工具
  find, f           文件目录查找工具
  grep, g           文本搜索工具
  hex2str, h2s      十六进制与字符串相互转换工具
  iconv, icv        文件编码转换工具
  json, j           JSON 数据处理工具
  list, ls          文件目录列表工具
  newline, nl       文件换行符检测与转换工具
  pack, pk          智能打包压缩工具
  preview, pv       压缩包预览工具
  proc, ps          查看系统进程信息
  unpack, upk       智能解压缩工具
  xargs, x          从标准输入或文件读取参数，批量执行指定命令
  awk               文本字段处理工具
  cat               显示文件内容
  cp                文件目录复制工具
  date              时间获取和格式化工具
  df                查看磁盘空间使用情况
  dns               DNS 查询工具
  echo              文本输出工具
  gm                获取Git仓库的各种元数据信息
  hash              文件哈希计算工具
  head              显示文件开头内容
  home              显示用户主目录
  md                预览 Markdown 文件
  mkdir             目录创建工具
  mv                文件移动工具
  ping              测试网络连通性
  port              查看系统端口占用情况
  pwd               打印当前工作目录
  rm                文件删除工具
  sed               文本替换工具
  seq               生成数字序列
  shck              Shell 脚本语法检查工具
  shfmt             Shell 脚本格式化工具
  shx               Shell 命令/脚本执行工具
  size              文件目录大小计算工具
  tail              显示文件结尾内容
  tcp               TCP 网络工具，支持端口扫描、客户端通信和服务端监听
  tee               从标准输入读取并输出到多个目标
  test              路径检测工具
  touch             文件创建和时间戳更新工具
  tr                字符转换工具
  truncate          文件截断工具
  watch             命令监控工具 (周期性执行命令并显示输出)
  wc                统计文件的行数、单词数、字节数和字符数
  which             在环境变量 PATH 中查找命令的可执行文件路径

示例:
  1. 永久启用
     fck.exe --install-completion pwsh

  2. 临时启用
     fck.exe --completion pwsh | Out-String | Invoke-Expression

注意:
  1. 各子命令有独立帮助文档，可通过 --help/-h 参数查看, 例如 'fck.exe <子命令> -h' 查看各子命令详细帮助

```

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
