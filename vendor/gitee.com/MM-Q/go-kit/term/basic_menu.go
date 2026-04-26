// 基础菜单功能
//
// 本文件提供基础菜单功能，适用于简单的数字选项场景。
// 依赖 menu.go 中的类型和错误定义。
//
// 主要功能:
//   - 简单的数字选项 (1, 2, 3...)
//   - 自动编号选项
//   - 支持自定义样式
//   - 支持循环显示
//   - 严格的输入验证 (使用 strconv.Atoi)
//
// 使用示例:
//
// 示例1: 基本使用
//
//	options := []string{
//	    "查看列表",
//	    "添加项目",
//	    "删除项目",
//	    "退出程序",
//	}
//
//	choice, err := term.ShowBasicMenuLine("主菜单", options, "请选择 (1-4): ")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("你选择了: %s\n", options[choice])
//
// 示例2: 循环菜单
//
//	options := []string{
//	    "查看列表",
//	    "添加项目",
//	    "删除项目",
//	    "退出程序",
//	}
//
//	err := term.ShowBasicMenuLoopLine("主菜单", options, "请选择 (1-4): ", func(index int) bool {
//	    switch index {
//	    case 0:
//	        fmt.Println("查看列表...")
//	        showList()
//	    case 1:
//	        fmt.Println("添加项目...")
//	        addItem()
//	    case 2:
//	        fmt.Println("删除项目...")
//	        deleteItem()
//	    case 3:
//	        fmt.Println("退出程序")
//	        return false  // 退出循环
//	    }
//	    return true  // 继续循环
//	})
//
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// 示例3: 自定义样式
//
//	options := []string{
//	    "查看列表",
//	    "添加项目",
//	}
//
//	style := &term.MenuStyle{
//	    Prefix:    ">> ",
//	    Separator: " | ",
//	    Indent:    2,
//	    ShowTitle: true,
//	}
//
//	menuText := term.RenderBasicMenu("主菜单", options, style)
//	fmt.Println(menuText)
//
// 示例4: 嵌套菜单
//
//	mainOptions := []string{
//	    "用户管理",
//	    "系统设置",
//	    "退出",
//	}
//
//	err := term.ShowBasicMenuLoopLine("主菜单", mainOptions, "请选择: ", func(index int) bool {
//	    switch index {
//	    case 0:
//	        showUserMenu()
//	    case 1:
//	        showSettingsMenu()
//	    case 2:
//	        return false  // 退出
//	    }
//	    return true
//	})
//
// 示例5: 错误处理
//
//	options := []string{
//	    "查看列表",
//	    "添加项目",
//	}
//
//	choice, err := term.ShowBasicMenuLine("主菜单", options, "请选择: ")
//	if err != nil {
//	    if menuErr, ok := err.(*term.MenuError); ok {
//	        fmt.Printf("菜单错误: %s\n", menuErr.Message)
//	    } else {
//	        fmt.Printf("未知错误: %v\n", err)
//	    }
//	    return
//	}
//	fmt.Printf("你选择了: %s\n", options[choice])
package term

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// # RenderBasicMenu 渲染基础菜单为字符串
//
// 参数:
//   - title: 菜单标题
//   - options: 选项列表
//   - style: 菜单样式 (可选, 使用默认样式)
//
// 返回:
//   - string: 渲染后的菜单字符串
func RenderBasicMenu(title string, options []string, style *MenuStyle) string {
	if style == nil {
		style = GetDefaultMenuStyle()
	}

	var builder strings.Builder

	if style.ShowTitle && title != "" {
		builder.WriteString(title)
		builder.WriteString("\n")
		builder.WriteString(strings.Repeat("=", len(title)))
		builder.WriteString("\n")
	}

	indent := strings.Repeat(" ", style.Indent)
	for i, option := range options {
		builder.WriteString(indent)
		builder.WriteString(style.Prefix)
		fmt.Fprintf(&builder, "%d", i+1)
		builder.WriteString(style.Separator)
		builder.WriteString(option)
		builder.WriteString("\n")
	}

	return builder.String()
}

// ShowBasicMenu 显示基础菜单并获取用户选择
//
// 参数:
//   - title: 菜单标题
//   - options: 选项列表
//   - reader: 输入源
//   - prompt: 提示信息
//
// 返回:
//   - int: 用户选择的索引 (从 0 开始)
//   - error: 错误信息
func ShowBasicMenu(title string, options []string, reader io.Reader, prompt string) (int, error) {
	if len(options) == 0 {
		return 0, ErrEmptyMenu
	}

	fmt.Print(RenderBasicMenu(title, options, nil))

	choiceStr := Read(reader, prompt+" ")
	if choiceStr == "" {
		return 0, ErrInvalidChoice
	}

	choice, err := strconv.Atoi(choiceStr)
	if err != nil {
		return 0, NewMenuError("INVALID_CHOICE", "无效的选择", choiceStr)
	}

	if choice < 1 || choice > len(options) {
		return 0, NewMenuError("INVALID_CHOICE", "无效的选择", fmt.Sprintf("%d (范围: 1-%d)", choice, len(options)))
	}

	return choice - 1, nil
}

// ShowBasicMenuLine 便捷函数, 使用 os.Stdin 显示菜单
//
// 参数:
//   - title: 菜单标题
//   - options: 选项列表
//   - prompt: 提示信息
//
// 返回:
//   - int: 用户选择的索引 (从 0 开始)
//   - error: 错误信息
func ShowBasicMenuLine(title string, options []string, prompt string) (int, error) {
	return ShowBasicMenu(title, options, os.Stdin, prompt)
}

// ShowBasicMenuLoop 循环显示基础菜单, 直到用户选择退出
//
// 参数:
//   - title: 菜单标题
//   - options: 选项列表
//   - reader: 输入源
//   - prompt: 提示信息
//   - handler: 处理函数, 返回 true 继续循环, false 退出
//
// 返回:
//   - error: 错误信息
func ShowBasicMenuLoop(title string, options []string, reader io.Reader, prompt string, handler func(index int) bool) error {
	if len(options) == 0 {
		return ErrEmptyMenu
	}

	for {
		choice, err := ShowBasicMenu(title, options, reader, prompt)
		if err != nil {
			return err
		}

		if !handler(choice) {
			return nil
		}
	}
}

// ShowBasicMenuLoopLine 便捷函数, 使用 os.Stdin 循环显示菜单
//
// 参数:
//   - title: 菜单标题
//   - options: 选项列表
//   - prompt: 提示信息
//   - handler: 处理函数, 返回 true 继续循环, false 退出
//
// 返回:
//   - error: 错误信息
func ShowBasicMenuLoopLine(title string, options []string, prompt string, handler func(index int) bool) error {
	return ShowBasicMenuLoop(title, options, os.Stdin, prompt, handler)
}
