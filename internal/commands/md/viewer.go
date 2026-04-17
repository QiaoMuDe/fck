package md

import (
	"bytes"
	"fmt"
	"os"

	"gitee.com/MM-Q/go-kit/utils"
	"github.com/charmbracelet/glamour"
	"github.com/noborus/ov/oviewer"
)

// MdViewer Markdown 查看器
type MdViewer struct {
	config   MdConfig
	renderer *glamour.TermRenderer
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

	return viewer, nil
}

// Run 启动查看器
//
// 返回值:
//   - error: 运行错误
func (v *MdViewer) Run() error {
	if v.config.UsePager {
		return v.runWithPager()
	}
	return v.runDirect()
}

// runDirect 直接输出到终端
//
// 返回值:
//   - error: 运行错误
func (v *MdViewer) runDirect() error {
	file := v.config.File

	// 检查文件大小
	fileInfo, err := os.Stat(file)
	if err != nil {
		return fmt.Errorf("failed to stat %s: %w", file, err)
	}

	if fileInfo.Size() > v.config.MaxSize {
		return fmt.Errorf("file %s size (%s) exceeds max size limit (%s), use -l flag to view with pager",
			file, utils.FormatBytes(fileInfo.Size()), utils.FormatBytes(v.config.MaxSize))
	}

	source, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", file, err)
	}

	rendered, err := v.renderer.RenderBytes(source)
	if err != nil {
		return fmt.Errorf("failed to render %s: %w", file, err)
	}

	fmt.Print(string(rendered))
	return nil
}

// runWithPager 使用分页器查看
//
// 返回值:
//   - error: 运行错误
func (v *MdViewer) runWithPager() error {
	file := v.config.File

	source, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", file, err)
	}

	// 创建渲染版文档
	renderDoc, err := oviewer.NewDocument()
	if err != nil {
		return fmt.Errorf("failed to create render document: %w", err)
	}
	renderDoc.FileName = file + " (rendered)"

	rendered, err := v.renderer.RenderBytes(source)
	if err != nil {
		return fmt.Errorf("failed to render markdown: %w", err)
	}

	if err := renderDoc.ControlReader(bytes.NewBuffer(rendered), nil); err != nil {
		return fmt.Errorf("failed to load rendered content: %w", err)
	}

	// 创建原始版文档
	originalDoc, err := oviewer.NewDocument()
	if err != nil {
		return fmt.Errorf("failed to create original document: %w", err)
	}
	originalDoc.FileName = file

	if err := originalDoc.ControlReader(bytes.NewBuffer(source), nil); err != nil {
		return fmt.Errorf("failed to load original content: %w", err)
	}

	ov, err := oviewer.NewOviewer(renderDoc, originalDoc)
	if err != nil {
		return fmt.Errorf("failed to create oviewer: %w", err)
	}

	return ov.Run()
}
