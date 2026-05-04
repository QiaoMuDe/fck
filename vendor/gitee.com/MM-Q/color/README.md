<div align="center">

<a name="top"></a>

# 🎨 color

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org) [![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE) [![Gitee](https://img.shields.io/badge/Gitee-MM--Q/color-C71D23?style=flat&logo=gitee)](https://gitee.com/MM-Q/color.git)

> 🌈 一个优雅、高效的 Go 语言 ANSI 颜色输出库，让你的终端输出更加丰富多彩！

</div>

---

## 🎯 项目简介

`color` 是一个轻量级的 Go 语言库，用于在终端中输出彩色文本。它提供了简单易用的 API，支持 ANSI 颜色代码、256色、RGB真彩色，以及多种文本样式（加粗、斜体、下划线等）。

> 📌 **项目来源**: 本项目基于 [github.com/fatih/color](https://github.com/fatih/color) (v1.19.0) 进行二次开发和改造。感谢 [Fatih Arslan](https://github.com/fatih) 和原项目的所有贡献者！

### ✨ 为什么选择 color？

- 🚀 **简单易用** - 一行代码即可输出彩色文本
- 🎨 **功能丰富** - 支持标准色、高亮色、256色、RGB真彩色
- 🔧 **灵活配置** - 支持链式调用，自由组合颜色和样式
- 🖥️ **跨平台** - 原生支持 Windows、macOS、Linux
- 📦 **零依赖** - 仅依赖 Go 标准库和官方扩展库
- 🌏 **中文注释** - 完整的中文文档和代码注释
- 📚 **详细文档** - 完整的项目分析报告和使用指南

---

## 🔥 核心特性

| 特性 | 描述 | 状态 |
|------|------|------|
| 🎨 **标准颜色** | 8种标准前景色和背景色 | ✅ 支持 |
| ✨ **高亮颜色** | 8种高亮版本的颜色 | ✅ 支持 |
| 🌈 **256色** | 扩展的256色调色板 | ✅ 支持 |
| 🎆 **RGB真彩色** | 24位真彩色（1600万色） | ✅ 支持 |
| 💪 **文本样式** | 加粗、斜体、下划线、删除线等 | ✅ 支持 |
| 🔄 **链式调用** | 流畅的 API 设计 | ✅ 支持 |
| 🎯 **便捷函数** | 32个预定义颜色的快速调用 | ✅ 支持 |
| 🚫 **颜色禁用** | 支持 NO_COLOR 环境变量 | ✅ 支持 |
| 🖥️ **跨平台** | Windows、macOS、Linux | ✅ 支持 |
| 🔄 **管道检测** | 自动检测管道输出并禁用颜色 | ✅ 支持 |

---

## 📦 安装指南

### 环境要求

- Go 1.25 或更高版本

### 安装命令

```bash
go get gitee.com/MM-Q/color
```

### 导入包

```go
import "gitee.com/MM-Q/color"
```

---

## 🚀 使用示例

### 基础用法

#### 1️⃣ 使用便捷函数（最简单）

```go
package main

import (
    "gitee.com/MM-Q/color"
)

func main() {
    // 打印红色文本（自动追加换行符）
    color.Red("这是一条红色消息")
    
    // 打印绿色文本（自动追加换行符）
    color.Green("这是一条绿色消息")
    
    // 打印蓝色文本（自动追加换行符）
    color.Blue("这是一条蓝色消息")
    
    // 使用格式化字符串（不自动追加换行符）
    color.Yellowf("警告: %s\n", "磁盘空间不足")
    
    // 获取带颜色的字符串（不打印，不自动追加换行符）
    redText := color.SRed("红色文本")
    println(redText)
}
```

#### 2️⃣ 使用高亮颜色

```go
package main

import (
    "gitee.com/MM-Q/color"
)

func main() {
    // 高亮颜色比普通颜色更明亮（自动追加换行符）
    color.HiRed("高亮红色")
    color.HiGreen("高亮绿色")
    color.HiYellow("高亮黄色")
    color.HiBlue("高亮蓝色")
    color.HiCyan("高亮青色")
    color.HiMagenta("高亮洋红色")
    
    // 使用格式化版本（不自动追加换行符）
    color.HiRedf("高亮红色: %s\n", "信息")
}
```

### 高级用法

#### 3️⃣ 链式调用组合样式

```go
package main

import (
    "gitee.com/MM-Q/color"
)

func main() {
    // 创建颜色对象并添加多个属性
    boldRed := color.New(color.FgRed).Add(color.Bold)
    boldRed.Println("粗体红色文本")
    
    // 红色前景 + 白色背景
    redOnWhite := color.New(color.FgRed, color.BgWhite)
    redOnWhite.Println("红字白底")
    
    // 高亮绿色 + 加粗 + 下划线
    fancy := color.New(color.FgHiGreen, color.Bold, color.Underline)
    fancy.Println("高亮绿色粗体下划线")
    
    // 使用 Add 方法逐步添加属性
    c := color.New(color.FgBlue)
    c.Add(color.Bold)
    c.Add(color.Underline)
    c.Println("蓝色粗体下划线")
}
```

#### 4️⃣ 使用 256 色

```go
package main

import (
    "gitee.com/MM-Q/color"
)

func main() {
    // 使用 256 色模式
    color.New(color.FgHiWhite).AddBg(color.BgHiBlack).Println("高对比度")
    
    // 创建特定的 256 色
    myColor := color.New(color.Attribute(196)) // 鲜红色
    myColor.Println("256色鲜红色")
}
```

#### 5️⃣ 使用 RGB 真彩色

```go
package main

import (
    "gitee.com/MM-Q/color"
)

func main() {
    // RGB 前景色
    color.RGB(255, 128, 0).Println("橙色文本")
    
    // RGB 背景色
    color.BgRGB(0, 128, 255).Println("蓝色背景")
    
    // 组合 RGB 前景和背景
    color.RGB(255, 255, 255).AddBgRGB(255, 0, 0).Println("白字红底")
    
    // 链式调用
    color.New().
        AddRGB(255, 165, 0).      // 橙色前景
        AddBgRGB(0, 0, 128).      // 深蓝色背景
        Add(color.Bold).          // 加粗
        Println("RGB组合样式")
}
```

#### 6️⃣ 直接输出到指定位置

```go
package main

import (
    "os"
    "gitee.com/MM-Q/color"
)

func main() {
    // 输出到标准错误
    color.New(color.FgRed).Fprintln(os.Stderr, "错误信息")
    
    // 输出到文件
    file, _ := os.Create("output.txt")
    defer file.Close()
    
    color.New(color.FgGreen).Fprintln(file, "写入文件的内容")
}
```

#### 7️⃣ 使用函数生成器

```go
package main

import (
    "gitee.com/MM-Q/color"
)

func main() {
    // 创建可复用的打印函数
    danger := color.New(color.FgRed, color.Bold).PrintfFunc()
    success := color.New(color.FgGreen).PrintlnFunc()
    
    // 多次使用
    danger("危险: %s\n", "系统即将关机")
    danger("错误: %s\n", "连接超时")
    
    success("操作成功完成")
    success("数据已保存")
}
```

---

## 📚 API文档

### 便捷函数

**函数命名规则：**
- `Xxx(message)` - 打印彩色文本并自动追加换行符（仅接受字符串）
- `Xxxf(format, a...)` - 打印格式化彩色文本（不自动追加换行符）
- `SXxx(message)` - 返回彩色字符串（仅接受字符串）
- `SXxxf(format, a...)` - 返回格式化彩色字符串

| 函数 | 描述 | 示例 |
|------|------|------|
| `Red()` | 红色打印（自动换行） | `color.Red("text")` |
| `Redf()` | 红色格式化打印 | `color.Redf("text: %s", "value")` |
| `SRed()` | 返回红色字符串 | `s := color.SRed("text")` |
| `SRedf()` | 返回格式化红色字符串 | `s := color.SRedf("text: %s", "value")` |
| `HiRed()` | 高亮红色打印 | `color.HiRed("text")` |
| `HiRedf()` | 高亮红色格式化打印 | `color.HiRedf("text: %s", "value")` |
| `Gray()` | 灰色打印 | `color.Gray("text")` |
| `Grayf()` | 灰色格式化打印 | `color.Grayf("text: %s", "value")` |

> 完整列表：`Black/Red/Green/Yellow/Blue/Magenta/Cyan/White` 和对应的高亮版本 `HiBlack/HiRed/...`，以及灰色 `Gray/SGray`，每个都有 `f` 后缀的格式化版本和 `S` 前缀的字符串返回版本。

### Color 对象方法

| 方法 | 描述 | 返回值 |
|------|------|--------|
| `New(attributes...)` | 创建颜色对象 | `*Color` |
| `Add(attributes...)` | 添加属性 | `*Color` |
| `RGB(r, g, b)` | 设置RGB前景色 | `*Color` |
| `AddRGB(r, g, b)` | 添加RGB前景色 | `*Color` |
| `BgRGB(r, g, b)` | 设置RGB背景色 | `*Color` |
| `AddBgRGB(r, g, b)` | 添加RGB背景色 | `*Color` |
| `Print(a...)` | 打印 | `(int, error)` |
| `Printf(format, a...)` | 格式化打印 | `(int, error)` |
| `Println(a...)` | 打印并换行 | `(int, error)` |
| `Sprint(a...)` | 返回字符串 | `string` |
| `Sprintf(format, a...)` | 格式化返回字符串 | `string` |
| `PrintFunc()` | 返回打印函数 | `func(a...)` |
| `PrintfFunc()` | 返回格式化打印函数 | `func(format string, a...)` |
| `DisableColor()` | 禁用颜色 | - |
| `EnableColor()` | 启用颜色 | - |

---

## 🎨 支持功能

### 文本样式

| 样式 | 常量 | 说明 |
|------|------|------|
| 重置 | `color.Reset` | 重置所有属性 |
| 加粗 | `color.Bold` | 粗体文本 |
| 轻淡 | `color.Faint` | 降低亮度 |
| 斜体 | `color.Italic` | 斜体文本 |
| 下划线 | `color.Underline` | 下划线文本 |
| 慢闪烁 | `color.BlinkSlow` | 慢速闪烁 |
| 快闪烁 | `color.BlinkRapid` | 快速闪烁 |
| 反显 | `color.ReverseVideo` | 前景背景互换 |
| 隐蔽 | `color.Concealed` | 隐藏文本 |
| 删除线 | `color.CrossedOut` | 删除线文本 |

### 颜色常量

```go
// 标准前景色（30-37）
color.FgBlack, color.FgRed, color.FgGreen, color.FgYellow
color.FgBlue, color.FgMagenta, color.FgCyan, color.FgWhite

// 高亮前景色（90-97）
color.FgHiBlack, color.FgHiRed, color.FgHiGreen, color.FgHiYellow
color.FgHiBlue, color.FgHiMagenta, color.FgHiCyan, color.FgHiWhite

// 灰色（高亮黑色的别名）
color.FgGray  // 灰色前景

// 标准背景色（40-47）
color.BgBlack, color.BgRed, color.BgGreen, color.BgYellow
color.BgBlue, color.BgMagenta, color.BgCyan, color.BgWhite

// 高亮背景色（100-107）
color.BgHiBlack, color.BgHiRed, color.BgHiGreen, color.BgHiYellow
color.BgHiBlue, color.BgHiMagenta, color.BgHiCyan, color.BgHiWhite

// 灰色背景（高亮黑色背景的别名）
color.BgGray  // 灰色背景
```

---

## ⚙️ 配置选项

### 全局配置

```go
// 禁用所有颜色输出
color.NoColor = true

// 设置默认输出位置
color.Output = os.Stdout
color.Error = os.Stderr
```

### 环境变量

| 变量 | 说明 | 示例 |
|------|------|------|
| `NO_COLOR` | 禁用所有颜色输出 | `export NO_COLOR=1` |

### 管道输出自动检测

color 库会自动检测输出是否是终端，当检测到管道输出或重定向时，会自动禁用颜色输出，避免在日志文件中留下 ANSI 转义码。

```bash
# 直接运行（终端输出）- 显示颜色
go run main.go

# 管道输出 - 自动禁用颜色
go run main.go | cat

# 重定向到文件 - 自动禁用颜色
go run main.go > output.log
```

> 💡 **提示**: 这是默认行为，符合 Unix 哲学和大多数 CLI 工具（如 `ls`、`grep`、`git`）的惯例。

### 局部控制

```go
// 创建禁用颜色的对象
c := color.New(color.FgRed)
c.DisableColor()
c.Println("这行不会显示颜色")

// 重新启用
c.EnableColor()
c.Println("这行会显示红色")
```

---

## 📁 项目结构

```
color/
├── 📄 doc.go              # 包文档说明
├── 🎯 color.go            # 核心类型和方法
├── 🎨 attribute.go        # 颜色常量定义
├── ⚙️ output.go           # 输出控制
├── 🚀 helper.go           # 便捷函数（代码生成）
├── 🔧 utils.go            # 内部工具
├── 🌐 global.go           # 全局实例管理
├── � global_methods.go   # 全局实例方法（代码生成）
├── �🖥️ color_windows.go    # Windows平台适配
├── 🧪 color_test.go       # 单元测试
├── 📁 gen/                # 代码生成工具
│   ├── helper_gen.go      # helper.go 生成器
│   └── global_gen.go      # global_methods.go 生成器
├── 📦 go.mod              # 模块定义
├── 🔒 go.sum              # 依赖校验
├── 📜 LICENSE             # MIT许可证
└── 📝 README.md           # 项目说明
```

---

## 🧪 测试说明

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行测试并显示覆盖率
go test -v -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 测试内容

- ✅ 基础颜色输出测试
- ✅ 链式调用测试
- ✅ RGB颜色测试
- ✅ 便捷函数测试
- ✅ 颜色禁用测试
- ✅ 并发安全测试

---

## 📄 许可证

本项目采用 [MIT 许可证](LICENSE) 开源。

```
MIT License

Copyright (c) 2026 M乔木

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
```

---

## 📮 联系方式

### 项目链接

- 🏠 **项目仓库**: [https://gitee.com/MM-Q/color.git](https://gitee.com/MM-Q/color.git)
- 🐛 **问题反馈**: [https://gitee.com/MM-Q/color/issues](https://gitee.com/MM-Q/color/issues)
- 📖 **使用文档**: 参见 [doc.go](doc.go) 和 [color-analysis.md](color-analysis.md)

### 原项目致谢

本项目基于 [github.com/fatih/color](https://github.com/fatih/color) 进行二次开发，感谢原项目的所有贡献者：

- 👤 [Fatih Arslan](https://github.com/fatih) - 原项目作者
- 🏠 [github.com/fatih/color](https://github.com/fatih/color) - 原项目地址

### 从原项目迁移

如果你之前使用过 `github.com/fatih/color`，需要注意以下 API 变更：

#### 字符串返回方法命名变更

为了与 Go 标准库（如 `fmt.Sprint`）保持一致，返回字符串的方法名从 `XxxString` 格式改为 `SXxx` 格式：

| 原项目 (fatih/color) | 本项目 | 说明 |
|---------------------|--------|------|
| `color.RedString()` | `color.SRed()` | 返回红色字符串 |
| `color.GreenString()` | `color.SGreen()` | 返回绿色字符串 |
| `color.BlueString()` | `color.SBlue()` | 返回蓝色字符串 |
| `color.HiRedString()` | `color.SHiRed()` | 返回高亮红色字符串 |
| ... | ... | 其他颜色类似 |

**迁移示例：**

```go
// 原项目代码（支持格式化）
redStr := color.RedString("error: %s", errMsg)

// 迁移后的代码（使用带 f 后缀的方法支持格式化）
redStr := color.SRedf("error: %s", errMsg)

// 如果不需要格式化，直接使用字符串
redStr := color.SRed("error message")
```

#### 方法命名规则

本项目的方法命名遵循以下规则：

| 方法后缀 | 功能 | 示例 |
|---------|------|------|
| 无后缀 | 直接打印字符串 | `color.Red("text")` |
| `f` 后缀 | 支持格式化打印 | `color.Redf("value: %d", 42)` |
| `S` 前缀 | 返回字符串 | `color.SRed("text")` |
| `S` + `f` | 返回格式化字符串 | `color.SRedf("value: %d", 42)` |

**注意：** 只有带 `f` 后缀的方法支持格式化参数（如 `%s`, `%d` 等），不带后缀的方法只接受字符串参数。

### 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建你的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交你的修改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开一个 Pull Request

---

## 🙏 致谢

感谢所有为本项目做出贡献的开发者！

---

<div align="center">

**[⬆ 回到顶部](#top)**

</div>

---
 
> 🔗 **项目地址**: [https://gitee.com/MM-Q/color.git](https://gitee.com/MM-Q/color.git)
