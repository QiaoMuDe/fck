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
			fmt.Fprintf(os.Stderr, "warn: syntax highlight failed: %v\n", hlErr)
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
