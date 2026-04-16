# Cat 命令语法高亮功能设计方案

## 1. 需求概述

基于 `chroma` 语法高亮库和 `ov` 分页器，为 cat 命令的 `--ov` 模式添加根据文件后缀自动语法高亮的功能。

## 2. 库 API 分析

### 2.1 Chroma 库 API

```
github.com/alecthomas/chroma/v2
├── lexers/           # 词法分析器（200+ 语言支持）
│   └── lexers.go     # GlobalLexerRegistry
├── formatters/       # 格式化输出器
│   ├── api.go        # Registry, Get(), Register()
│   ├── tty_indexed.go # 256色ANSI格式
│   └── tty_truecolour.go # 真彩色ANSI格式
├── styles/           # 主题样式（40+ 主题）
│   └── api.go        # Registry, Get()
└── chroma.go         # 核心接口定义
```

**关键 API**:
```go
// lexers/lexers.go
func Match(filename string) chroma.Lexer    // 根据文件名匹配 lexer
func Get(name string) chroma.Lexer          // 根据语言名获取 lexer

// formatters/api.go  
func Get(name string) chroma.Formatter      // "terminal256" 或 "terminal16m"
var Fallback chroma.Formatter               // 默认 formatter

// styles/api.go
func Get(name string) *chroma.Style         // 获取样式，如 "monokai"
var Fallback *chroma.Style                 // 默认样式

// 核心高亮流程
iterator, err := lexer.Tokenise(nil, code)  // 词法分析
formatter.Format(writer, style, iterator)   // 格式化输出
```

### 2.2 OV 库 API

```
github.com/noborus/ov/oviewer
├── oviewer.go        # Root 结构（主入口）
│   ├── Open(files...) -> *Root             // 打开文件（当前使用）
│   ├── NewOviewer(docs...) -> *Root        // 从 Document 创建
│   └── Run() -> error                       // 运行分页器
├── document.go       # Document 结构
│   ├── NewDocument() -> *Document          // 创建空文档
│   ├── OpenDocument(file) -> *Document     // 打开文件
│   └── ControlReader(reader, reload)       // 从 Reader 加载
└── control.go        # 控制逻辑
    └── ControlReader(r io.Reader, reload)  // 关键：从 io.Reader 读取
```

**关键发现 - ControlReader API**:
```go
// Document.ControlReader 是核心 API
func (m *Document) ControlReader(r io.Reader, reload func() *bufio.Reader) error

// 使用流程：
doc, _ := oviewer.NewDocument()          // 1. 创建空文档
doc.FileName = filename                   // 2. 设置文件名
doc.ControlReader(reader, nil)           // 3. 从 Reader 加载内容
root, _ := oviewer.NewOviewer(doc)       // 4. 创建 Root
root.Run()                               // 5. 运行
```

## 3. 实现方案

### 3.1 方案对比

| 方案 | 描述 | 优点 | 缺点 | 推荐度 |
|------|------|------|------|--------|
| A | 直接修改文件后传给 ov | 实现简单 | 需要临时文件或内存缓冲 | ⭐⭐⭐ |
| B | 使用 ControlReader API | 内存直接传递，无临时文件 | 需要了解 ov 内部 API | ⭐⭐⭐⭐⭐ |
| C | 通过 stdin 管道 | 无需修改 ov 调用方式 | 可能丢失文件名信息 | ⭐⭐ |

### 3.2 推荐方案 B：ControlReader API

**流程图**:
```
┌─────────────────┐
│  读取源文件内容  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐     否      ┌────────────────────────┐
│ 是否支持高亮?   │────────────▶│ doc.ControlReader(     │
│ (根据后缀判断)  │             │   fileReader, nil)     │
└────────┬────────┘             └────────────────────────┘
         │ 是
         ▼
┌─────────────────┐
│ 使用 Chroma 高亮 │
│ - lexers.Match  │
│ - Tokenise      │
│ - ANSI 格式化   │
└────────┬────────┘
         │
         ▼
┌────────────────────────┐
│ doc.ControlReader(     │
│   highlightedReader,   │
│   nil)                 │
└────────┬───────────────┘
         │
         ▼
┌─────────────────┐
│ NewOviewer(doc) │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ root.Run()      │
└─────────────────┘
```

## 4. 详细设计

### 4.1 新增文件

**`internal/commands/cat/syntax_highlight.go`**

```go
package cat

import (
	"bytes"
	"fmt"
	"io"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// HighlightConfig 语法高亮配置
type HighlightConfig struct {
	Enabled     bool   // 是否启用高亮
	Style       string // 主题样式
	Formatter   string // 格式化器类型: "terminal256", "terminal16m", "terminal"
}

// DefaultHighlightConfig 返回默认高亮配置
func DefaultHighlightConfig() HighlightConfig {
	return HighlightConfig{
		Enabled:   true,
		Style:     "monokai",
		Formatter: "terminal256", // 使用 256 色，兼容性最好
	}
}

// highlightContent 对内容进行语法高亮
//
// 参数:
//   - content: 文件内容
//   - filename: 文件名（用于检测语言）
//   - config: 高亮配置
//
// 返回:
//   - []byte: 高亮后的 ANSI 格式内容
//   - bool: 是否成功高亮
//   - error: 错误信息
func highlightContent(content []byte, filename string, config HighlightConfig) ([]byte, bool, error) {
	if !config.Enabled {
		return content, false, nil
	}

	// 1. 根据文件名获取 lexer
	// lexers.Match 会根据文件名模式匹配最合适的 lexer
	// 如果返回 nil，表示 Chroma 不支持该文件类型
	lexer := lexers.Match(filename)
	if lexer == nil {
		return content, false, nil // 不支持的语言，返回原内容
	}

	// 2. 获取 formatter（ANSI 终端格式）
	formatter := formatters.Get(config.Formatter)
	if formatter == nil {
		formatter = formatters.Fallback
	}

	// 3. 获取样式
	style := styles.Get(config.Style)
	if style == nil {
		style = styles.Fallback
	}

	// 4. 执行高亮
	iterator, err := lexer.Tokenise(nil, string(content))
	if err != nil {
		return content, false, fmt.Errorf("tokenise failed: %w", err)
	}

	var buf bytes.Buffer
	err = formatter.Format(&buf, style, iterator)
	if err != nil {
		return content, false, fmt.Errorf("format failed: %w", err)
	}

	return buf.Bytes(), true, nil
}

// shouldHighlight 判断文件是否应该高亮
//
// 参数:
//   - filename: 文件名
//
// 返回:
//   - bool: 是否应该高亮
//
// 说明:
//   使用 Chroma 的 lexers.Match 函数自动检测文件类型。
//   如果返回 nil，表示 Chroma 不支持该文件类型的语法高亮。
func shouldHighlight(filename string) bool {
	return lexers.Match(filename) != nil
}
```

### 4.2 修改文件

**`internal/commands/cat/ov_pager.go`**

```go
package cat

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/noborus/ov/oviewer"
)

// runOVMode 使用 ov 库进行分页查看（支持语法高亮）
//
// 参数:
//   - config: 命令配置
//
// 返回:
//   - error: 执行错误
func runOVMode(config CatConfig) error {
	// 验证参数: ov 模式只支持单个文件
	if len(config.Targets) != 1 {
		return fmt.Errorf("ov mode only supports single file, got %d files", len(config.Targets))
	}

	filename := config.Targets[0]

	// 读取文件内容
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// 创建 Document
	doc, err := oviewer.NewDocument()
	if err != nil {
		return fmt.Errorf("failed to create document: %w", err)
	}

	doc.FileName = filename

	// 判断是否进行语法高亮
	var reader io.Reader
	hlConfig := DefaultHighlightConfig()

	if shouldHighlight(filename) {
		hlContent, ok, hlErr := highlightContent(content, filename, hlConfig)
		if hlErr != nil {
			// 高亮失败，输出警告但继续使用原内容
			fmt.Fprintf(os.Stderr, "Warning: syntax highlight failed: %v\n", hlErr)
			reader = bytes.NewReader(content)
		} else if ok {
			reader = bytes.NewReader(hlContent)
			doc.Caption = "[syntax highlighted]"
		} else {
			reader = bytes.NewReader(content)
		}
	} else {
		reader = bytes.NewReader(content)
	}

	// 使用 ControlReader 加载内容
	if err := doc.ControlReader(reader, nil); err != nil {
		return fmt.Errorf("failed to load content: %w", err)
	}

	// 创建 oviewer Root
	root, err := oviewer.NewOviewer(doc)
	if err != nil {
		return fmt.Errorf("failed to create oviewer: %w", err)
	}

	// 运行分页器
	return root.Run()
}
```

### 4.3 依赖更新

需要确保以下包在 vendor 中：

```bash
# 添加 Chroma 子包
go get github.com/alecthomas/chroma/v2/formatters
go get github.com/alecthomas/chroma/v2/lexers
go get github.com/alecthomas/chroma/v2/styles

go mod vendor
```

## 5. 支持的语言

基于 Chroma 支持的语言，常见文件类型：

| 类别 | 扩展名 |
|------|--------|
| Go | .go |
| Python | .py, .pyw |
| JavaScript/TypeScript | .js, .jsx, .ts, .tsx |
| Java | .java |
| C/C++ | .c, .cpp, .cc, .h, .hpp |
| Rust | .rs |
| Ruby | .rb |
| PHP | .php |
| Swift | .swift |
| Kotlin | .kt |
| Scala | .scala |
| HTML/XML | .html, .htm, .xml |
| CSS/SCSS | .css, .scss, .sass, .less |
| JSON/YAML/TOML | .json, .yaml, .yml, .toml |
| Markdown | .md, .markdown |
| Shell | .sh, .bash, .zsh, .fish |
| SQL | .sql |
| Lua | .lua |
| Vim | .vim |
| PowerShell | .ps1 |
| Perl | .pl |
| Haskell | .hs |
| Elm | .elm |
| Elixir | .ex |
| Erlang | .erl |
| R | .r |
| Makefile | Makefile |
| Dockerfile | Dockerfile |

## 6. 主题样式

Chroma 内置 40+ 主题，常用：

- `monokai` - 默认，深色背景，高对比度
- `github` - GitHub 风格，浅色
- `solarized-dark` / `solarized-light` - Solarized 配色
- `dracula` - 流行的深色主题
- `vs` - Visual Studio 风格
- `xcode` - Xcode 风格
- `vim` - Vim 风格

## 7. 实现步骤

1. **添加依赖**
   ```bash
   go get github.com/alecthomas/chroma/v2/formatters
   go get github.com/alecthomas/chroma/v2/lexers
   go get github.com/alecthomas/chroma/v2/styles
   go mod vendor
   ```

2. **创建 `syntax_highlight.go`**
   - 实现 `highlightContent()` 函数
   - 实现 `shouldHighlight()` 函数
   - 定义 `HighlightConfig` 结构体

3. **修改 `ov_pager.go`**
   - 集成高亮逻辑
   - 使用 `ControlReader()` API 替代 `Open()`

4. **测试**
   ```bash
   # 测试 Go 文件高亮
   fck cat --ov main.go
   
   # 测试 Python 文件高亮
   fck cat --ov script.py
   
   # 测试不支持高亮的文件（应正常显示）
   fck cat --ov file.txt
   
   # 测试不同主题（后续扩展）
   fck cat --ov --highlight-style=dracula main.go
   ```

## 8. 注意事项

1. **性能考虑**
   - 大文件高亮可能较慢（>10MB）
   - 考虑添加文件大小限制（如 >50MB 跳过高亮）

2. **内存使用**
   - 高亮后的内容会全部加载到内存
   - 超大文件建议跳过高亮

3. **ANSI 代码**
   - Chroma 的 terminal formatter 输出 ANSI 颜色代码
   - ov 原生支持 ANSI 颜色显示

4. **错误处理**
   - 高亮失败应优雅降级，显示原内容
   - 输出警告信息到 stderr

5. **ControlReader 行为**
   - `ControlReader` 是异步读取的
   - 需要确保 reader 在调用前已准备好所有数据
   - 对于大文件，考虑使用 `bufio.Reader` 包装

## 9. 扩展建议

1. **添加更多配置选项**
   - `--highlight-style` 选择主题
   - `--no-highlight` 禁用高亮
   - `--true-color` 使用 24 位真彩色

2. **支持从 stdin 高亮**
   - 通过 `--lang` 参数指定语言
   - 例如：`cat file.go | fck cat --ov --lang=go`

3. **性能优化**
   - 对大文件使用流式高亮
   - 缓存高亮结果

4. **行号显示**
   - ov 支持 `LineNumber` 设置
   - 可在配置中开启

## 10. 代码示例

### 完整使用示例

```go
package main

import (
    "fmt"
    "os"
    
    "github.com/alecthomas/chroma/v2/formatters"
    "github.com/alecthomas/chroma/v2/lexers"
    "github.com/alecthomas/chroma/v2/styles"
    "github.com/noborus/ov/oviewer"
)

func main() {
    filename := "main.go"
    content, _ := os.ReadFile(filename)
    
    // 1. 高亮处理
    lexer := lexers.Match(filename)
    style := styles.Get("monokai")
    formatter := formatters.Get("terminal256")
    
    import "bytes"
    var buf bytes.Buffer
    iterator, _ := lexer.Tokenise(nil, string(content))
    formatter.Format(&buf, style, iterator)
    
    // 2. 创建 Document 并加载
    doc, _ := oviewer.NewDocument()
    doc.FileName = filename
    doc.ControlReader(&buf, nil)
    
    // 3. 运行 ov
    root, _ := oviewer.NewOviewer(doc)
    root.Run()
}
```
