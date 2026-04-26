package term

import (
	"os"
	"strconv"

	"golang.org/x/term"
)

// IsStdinPipe 检测标准输入是否为管道或文件重定向
//
// 该函数用于区分程序是从终端直接接收输入，还是通过管道(|)或文件重定向(<)接收输入
// 常用于判断是否需要显示交互式提示
//
// 注意: 如果获取 stdin 状态失败，该函数返回 false（按终端处理）
// 如需获取错误详情，请使用 IsStdinPipeWithError
//
// 返回:
//   - bool: true 表示 stdin 是管道或文件重定向; false 表示是终端输入或出错
//
// 示例:
//
//	if term.IsStdinPipe() {
//	    // 从管道读取数据
//	    data, _ := io.ReadAll(os.Stdin)
//	} else {
//	    // 交互式终端，显示提示
//	    fmt.Print("请输入: ")
//	}
func IsStdinPipe() bool {
	isPipe, _ := IsStdinPipeWithError()
	return isPipe
}

// IsStdinPipeWithError 检测标准输入是否为管道或文件重定向
//
// 与 IsStdinPipe 功能相同，但会返回详细的错误信息
// 适用于需要精确错误处理的场景
//
// 返回:
//   - bool: true 表示 stdin 是管道或文件重定向; false 表示是终端输入
//   - error: 获取 stdin 状态时发生的错误，nil 表示成功
//
// 示例:
//
//	isPipe, err := term.IsStdinPipeWithError()
//	if err != nil {
//	    log.Printf("无法检测输入类型: %v", err)
//	    return
//	}
//	if isPipe {
//	    // 处理管道输入
//	} else {
//	    // 处理交互式输入
//	}
func IsStdinPipeWithError() (bool, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false, err
	}
	// 检查是否为命名管道或常规文件（重定向）
	isPipe := info.Mode()&os.ModeNamedPipe != 0 || info.Mode().IsRegular()
	return isPipe, nil
}

// GetSafeTerminalWidth 获取当前终端的宽度（字符列数）
//
// 该函数会尝试多种方式获取终端宽度，并确保返回的值在合理范围内:
// 1. 首先检查 COLUMNS 环境变量
// 2. 其次尝试通过系统调用获取终端尺寸
// 3. 如果都失败，返回默认值
//
// 返回:
//   - int: 终端宽度（字符列数），范围在 [40, 1200] 之间
//
// 示例:
//
//	width := term.GetSafeTerminalWidth()
//	fmt.Printf("当前终端宽度: %d 列\n", width)
func GetSafeTerminalWidth() int {
	defaultWidth := 80 // 默认宽度
	minWidth := 40     // 最小宽度
	maxWidth := 1200   // 最大宽度

	// 检查环境变量
	if cols := os.Getenv("COLUMNS"); cols != "" {
		if width, err := strconv.Atoi(cols); err == nil && width >= minWidth && width <= maxWidth {
			return width
		}
	}

	// 检查是否为终端
	fd := os.Stdout.Fd()
	if fd > 1024 || !term.IsTerminal(int(fd)) {
		return defaultWidth
	}

	// 安全的类型转换和获取尺寸
	if fd <= uintptr(^uint(0)>>1) { // 确保不会溢出
		if width, _, err := term.GetSize(int(fd)); err == nil {
			if width >= minWidth && width <= maxWidth {
				return width
			}
		}
	}

	return defaultWidth
}
