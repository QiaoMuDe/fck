package md

import (
	"fmt"
)

// MdConfig md 命令配置
type MdConfig struct {
	File      string // 要预览的文件（管道模式为空）
	UsePipe   bool   // 是否从管道读取
	UsePager  bool   // 使用分页器
	ShowRaw   bool   // 分页器中显示原始文件
	Style     string // 渲染样式
	WordWidth int    // 换行宽度
	MaxSize   int64  // 最大文件大小（字节）
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
