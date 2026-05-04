package color

import (
	"os"

	"golang.org/x/sys/windows"
)

// init 函数在包被导入时自动执行
// 用于在 Windows 系统上启用 ANSI 颜色支持
func init() {
	// 为当前进程启用 ANSI 颜色支持
	// 参考文档: https://learn.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences#output-sequences

	// 检查标准输出是否为空（例如以 Windows 服务运行时）
	if os.Stdout == nil {
		return
	}

	// 定义变量存储控制台输出模式
	var outMode uint32

	// 获取标准输出的 Windows 句柄
	out := windows.Handle(os.Stdout.Fd())

	// 获取当前控制台模式
	// 如果失败（例如输出被重定向到文件），则直接返回
	if err := windows.GetConsoleMode(out, &outMode); err != nil {
		return
	}

	// 启用以下控制台模式：
	// - ENABLE_PROCESSED_OUTPUT: 启用处理输出控制序列
	// - ENABLE_VIRTUAL_TERMINAL_PROCESSING: 启用虚拟终端处理（支持 ANSI 转义序列）
	outMode |= windows.ENABLE_PROCESSED_OUTPUT | windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING

	// 设置新的控制台模式
	// 错误被忽略，因为如果设置失败，程序仍可以正常运行（只是没有颜色）
	_ = windows.SetConsoleMode(out, outMode)
}
