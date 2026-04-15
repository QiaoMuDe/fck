package cat

import (
	"fmt"

	"github.com/noborus/ov/oviewer"
)

// runOVMode 使用 ov 库进行分页查看
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

	// 使用 ov 打开文件
	root, err := oviewer.Open(config.Targets[0])
	if err != nil {
		return fmt.Errorf("failed to open file with ov: %w", err)
	}

	// 运行分页器
	return root.Run()
}
