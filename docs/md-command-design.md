# md 子命令设计方案

## 一、功能概述

新增 `md` 子命令，用于预览 Markdown 文件。结合 glamour 库进行 Markdown 渲染，oviewer 库提供分页查看功能。

## 二、参考实现分析

参考 `mdviewer` 项目的核心设计：

1. **双文档模式**：同时加载渲染版和原始版，支持按键切换
2. **glamour 渲染**：使用 `glamour.NewTermRenderer` 创建渲染器
3. **oviewer 集成**：通过 `oviewer.NewDocument()` 创建文档，`oviewer.NewOviewer()` 启动查看器

## 三、目录结构

```
internal/
├── commands/md/
│   ├── cmd_md.go       # 业务逻辑主文件
│   └── viewer.go       # glamour + oviewer 封装
└── cli/
    └── md.go           # CLI 定义
```

## 四、CLI 定义 (internal/cli/md.go)

```go
package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/md"
	"gitee.com/MM-Q/qflag"
)

var MdCmd *qflag.Cmd

var (
	mdOnly    *qflag.BoolFlag   // -o, --only 仅渲染模式，不显示原始文件
	mdStyle   *qflag.StringFlag // -s, --style 渲染样式
	mdWidth   *qflag.IntFlag    // -w, --width 换行宽度
)

func init() {
	MdCmd = qflag.NewCmd("md", "mdv", qflag.ExitOnError)

	mdOnly = MdCmd.Bool("only", "o", "仅显示渲染后的内容，不显示原始文件", false)
	mdStyle = MdCmd.String("style", "s", "渲染样式 (auto/dark/light/dracula/pink/notty)", "auto")
	mdWidth = MdCmd.Int("width", "w", "换行宽度 (0 表示自动)", 0)

	cmdOpts := &qflag.CmdOpts{
		Desc:       "预览 Markdown 文件",
		UseChinese: true,
		Examples: map[string]string{
			"预览单个文件":   "fck md README.md",
			"预览多个文件":   "fck md doc1.md doc2.md",
			"仅渲染模式":     "fck md -o README.md",
			"指定样式":       "fck md -s dark README.md",
			"指定宽度":       "fck md -w 100 README.md",
		},
		Notes: []string{
			"支持 glamour 的所有内置样式: auto, dark, light, dracula, pink, notty",
			"默认使用 auto 样式，根据终端背景自动选择",
			"在分页器中按 ']' 查看原始文件，按 '[' 返回渲染视图",
			"支持同时打开多个文件，使用 n/p 切换",
		},
	}

	if err := MdCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	MdCmd.SetRun(runMd)
}

func runMd(cmd qflag.Command) error {
	args := cmd.Args()
	if len(args) == 0 {
		return fmt.Errorf("请指定要预览的 Markdown 文件")
	}

	config := md.MdConfig{
		Files:     args,
		Only:      mdOnly.Get(),
		Style:     mdStyle.Get(),
		WordWidth: mdWidth.Get(),
	}

	return md.MdCmdMain(config)
}
```

## 五、业务逻辑 (internal/commands/md/cmd_md.go)

```go
package md

import (
	"fmt"
)

// MdConfig md 命令配置
type MdConfig struct {
	Files     []string // 要预览的文件列表
	Only      bool     // 仅渲染模式
	Style     string   // 渲染样式
	WordWidth int      // 换行宽度
}

// MdCmdMain md 命令主入口
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - error: 执行错误
func MdCmdMain(config MdConfig) error {
	viewer, err := NewMdViewer(config)
	if err != nil {
		return fmt.Errorf("failed to create viewer: %w", err)
	}

	return viewer.Run()
}
```

## 六、Viewer 封装 (internal/commands/md/viewer.go)

```go
package md

import (
	"bytes"
	"fmt"
	"os"

	"github.com/charmbracelet/glamour"
	"github.com/noborus/ov/oviewer"
)

// MdViewer Markdown 查看器
type MdViewer struct {
	config    MdConfig
	renderer  *glamour.TermRenderer
	documents []*oviewer.Document
}

// NewMdViewer 创建 Markdown 查看器
//
// 参数:
//   - config: 查看器配置
//
// 返回值:
//   - *MdViewer: 查看器实例
//   - error: 创建错误
func NewMdViewer(config MdConfig) (*MdViewer, error) {
	// 创建 glamour 渲染器
	opts := []glamour.TermRendererOption{
		glamour.WithStylePath(config.Style),
	}
	
	if config.WordWidth > 0 {
		opts = append(opts, glamour.WithWordWrap(config.WordWidth))
	}

	renderer, err := glamour.NewTermRenderer(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create renderer: %w", err)
	}

	viewer := &MdViewer{
		config:   config,
		renderer: renderer,
	}

	// 加载所有文档
	for _, file := range config.Files {
		docs, err := viewer.loadDocuments(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to load %s: %v\n", file, err)
			continue
		}
		viewer.documents = append(viewer.documents, docs...)
	}

	if len(viewer.documents) == 0 {
		return nil, fmt.Errorf("no valid files to display")
	}

	return viewer, nil
}

// loadDocuments 加载单个文件的文档
//
// 参数:
//   - fileName: 文件名
//
// 返回值:
//   - []*oviewer.Document: 文档列表（渲染版和原始版）
//   - error: 加载错误
func (v *MdViewer) loadDocuments(fileName string) ([]*oviewer.Document, error) {
	source, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// 创建渲染版文档
	renderDoc, err := oviewer.NewDocument()
	if err != nil {
		return nil, fmt.Errorf("failed to create render document: %w", err)
	}
	renderDoc.FileName = fileName + " (rendered)"

	rendered, err := v.renderer.RenderBytes(source)
	if err != nil {
		return nil, fmt.Errorf("failed to render markdown: %w", err)
	}

	if err := renderDoc.ControlReader(bytes.NewBuffer(rendered), nil); err != nil {
		return nil, fmt.Errorf("failed to load rendered content: %w", err)
	}

	// 仅渲染模式，不加载原始文件
	if v.config.Only {
		return []*oviewer.Document{renderDoc}, nil
	}

	// 创建原始版文档
	originalDoc, err := oviewer.NewDocument()
	if err != nil {
		return nil, fmt.Errorf("failed to create original document: %w", err)
	}
	originalDoc.FileName = fileName

	if err := originalDoc.ControlReader(bytes.NewBuffer(source), nil); err != nil {
		return nil, fmt.Errorf("failed to load original content: %w", err)
	}

	return []*oviewer.Document{renderDoc, originalDoc}, nil
}

// Run 启动查看器
//
// 返回值:
//   - error: 运行错误
func (v *MdViewer) Run() error {
	ov, err := oviewer.NewOviewer(v.documents...)
	if err != nil {
		return fmt.Errorf("failed to create oviewer: %w", err)
	}

	return ov.Run()
}
```

## 七、依赖添加

在 `go.mod` 中添加：

```go
require (
	// ... 现有依赖
	github.com/charmbracelet/glamour v0.8.0
)
```

**注意**：oviewer 已在 vendor 中，glamour 需要新增。

## 八、支持的样式

| 样式 | 说明 |
|------|------|
| `auto` | 自动检测终端背景色 |
| `dark` | 暗色主题 |
| `light` | 亮色主题 |
| `dracula` | Dracula 主题 |
| `pink` | 粉色主题 |
| `notty` | 无颜色输出 |

## 九、使用示例

```bash
# 基本使用
fck md README.md

# 预览多个文件
fck md README.md CHANGELOG.md docs/guide.md

# 仅渲染模式（不显示原始文件）
fck md -o README.md

# 指定暗色主题
fck md -s dark README.md

# 指定换行宽度
fck md -w 120 README.md
```

## 十、分页器快捷键

继承 oviewer 的所有快捷键：

| 按键 | 功能 |
|------|------|
| `]` | 切换到原始文件视图 |
| `[` | 切换到渲染视图 |
| `n` | 下一个文件 |
| `p` | 上一个文件 |
| `q` | 退出 |
| `h` | 显示帮助 |

## 十一、边缘情况处理

1. **文件不存在**：输出警告，继续处理其他文件
2. **空文件**：正常显示，内容为空
3. **非 Markdown 文件**：尝试渲染，可能显示异常
4. **大文件**：oviewer 自动处理分页
5. **二进制文件**：渲染失败，显示警告

## 十二、后续优化方向

1. 支持从管道读取 Markdown 内容
2. 支持自定义 glamour 样式文件
3. 支持导出渲染后的内容为 HTML/ANSI
4. 集成到 cat 命令作为自动检测（根据文件扩展名）
