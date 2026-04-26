// 结构化菜单功能
//
// 本文件提供结构化菜单功能，适用于需要自定义菜单键、退出选项、默认值等复杂场景。
//
// 主要功能:
//   - 支持自定义菜单键 (如 "1", "a", "add")
//   - 支持默认值
//   - 支持退出选项
//   - 支持自定义样式
//   - 支持循环显示
//
// 使用示例:
//
// 示例1: 基本使用
//
//	menu := &term.Menu{
//	    Title: "主菜单",
//	    Items: []term.MenuItem{
//	        {Key: "1", Value: "查看列表"},
//	        {Key: "2", Value: "添加项目"},
//	        {Key: "3", Value: "删除项目"},
//	    },
//	    Prompt:   "请选择 (1-3): ",
//	    Default:  "1",
//	    AllowExit: true,
//	    ExitKey:  "q",
//	    ExitText: "退出",
//	}
//
//	choice, err := term.ShowMenuLine(menu)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("你选择了: %s\n", choice)
//
// 示例2: 字母键菜单
//
//	menu := &term.Menu{
//	    Title: "操作菜单",
//	    Items: []term.MenuItem{
//	        {Key: "a", Value: "添加项目"},
//	        {Key: "d", Value: "删除项目"},
//	        {Key: "v", Value: "查看列表"},
//	    },
//	    Prompt:   "请选择: ",
//	    AllowExit: true,
//	    ExitKey:  "q",
//	    ExitText: "退出",
//	}
//
//	choice, _ := term.ShowMenuLine(menu)
//	fmt.Printf("你选择了: %s\n", choice)
//
// 示例3: 循环菜单
//
//	menu := &term.Menu{
//	    Title: "主菜单",
//	    Items: []term.MenuItem{
//	        {Key: "1", Value: "查看列表"},
//	        {Key: "2", Value: "添加项目"},
//	    },
//	    Prompt:   "请选择: ",
//	    AllowExit: true,
//	    ExitKey:  "q",
//	    ExitText: "退出",
//	}
//
//	err := term.ShowMenuLoopLine(menu, func(key string) bool {
//	    switch key {
//	    case "1":
//	        fmt.Println("查看列表...")
//	    case "2":
//	        fmt.Println("添加项目...")
//	    case "q":
//	        return false  // 退出循环
//	    }
//	    return true  // 继续循环
//	})
//	if err != nil {
//	    fmt.Println("菜单循环错误:", err)
//	}
//
// 示例4: 自定义样式
//
//	menu := &term.Menu{
//	    Title: "主菜单",
//	    Items: []term.MenuItem{
//	        {Key: "1", Value: "查看列表"},
//	    },
//	    Prompt: "请选择: ",
//	}
//
//	style := &term.MenuStyle{
//	    Prefix:    ">> ",
//	    Separator: " | ",
//	    Indent:    2,
//	    ShowTitle: true,
//	}
//
//	menuText := term.RenderMenu(menu, style)
//	fmt.Println(menuText)
//
// 示例5: 嵌套菜单
//
//	mainMenu := &term.Menu{
//	    Title: "主菜单",
//	    Items: []term.MenuItem{
//	        {Key: "1", Value: "用户管理"},
//	        {Key: "2", Value: "系统设置"},
//	    },
//	    Prompt:   "请选择: ",
//	    AllowExit: true,
//	    ExitKey:  "q",
//	    ExitText: "退出",
//	}
//
//	term.ShowMenuLoopLine(mainMenu, func(key string) bool {
//	    switch key {
//	    case "1":
//	        showUserMenu()
//	    case "2":
//	        showSettingsMenu()
//	    case "q":
//	        return false
//	    }
//	    return true
//	})
package term

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// MenuItem 菜单项
type MenuItem struct {
	Key   string // 选项键 (如 "1", "a", "add")
	Value string // 选项描述 (如 "查看列表")
}

// Menu 菜单配置 (结构化菜单)
type Menu struct {
	Title     string     // 菜单标题, 如 "主菜单"
	Items     []MenuItem // 菜单项列表 (如 ["1. 查看", "2. 添加"])
	Prompt    string     // 提示信息, 如 "请输入选择"
	Default   string     // 默认值 (可选), 如 "1"
	AllowExit bool       // 是否允许退出, 如 true
	ExitKey   string     // 退出键 (如 "q", "0"), 如 "q"
	ExitText  string     // 退出选项描述 (如 "退出"), 如 "退出"
}

// MenuStyle 菜单样式
type MenuStyle struct {
	Prefix    string // 选项前缀 (如 "1. ", "[a] ")
	Separator string // 分隔符 (如 " - ", ": ")
	Indent    int    // 缩进空格数
	ShowTitle bool   // 是否显示标题
}

// MenuError 菜单错误类型
type MenuError struct {
	Code    string // 错误代码
	Message string // 错误消息
	Detail  string // 详细信息
}

// Error 实现 error 接口
func (e *MenuError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Detail)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// 预定义错误
var (
	ErrEmptyMenu      = &MenuError{Code: "EMPTY_MENU", Message: "菜单不能为空"}
	ErrInvalidChoice  = &MenuError{Code: "INVALID_CHOICE", Message: "无效的选择"}
	ErrDuplicateKey   = &MenuError{Code: "DUPLICATE_KEY", Message: "菜单项键重复"}
	ErrInvalidDefault = &MenuError{Code: "INVALID_DEFAULT", Message: "默认值无效"}
)

// NewMenuError 创建菜单错误
func NewMenuError(code, message, detail string) *MenuError {
	return &MenuError{
		Code:    code,
		Message: message,
		Detail:  detail,
	}
}

// GetDefaultMenuStyle 获取默认菜单样式
//
// 返回:
//   - *MenuStyle: 默认菜单样式
func GetDefaultMenuStyle() *MenuStyle {
	return &MenuStyle{
		Prefix:    "",
		Separator: ". ",
		Indent:    0,
		ShowTitle: true,
	}
}

// ValidateMenu 验证菜单配置
//
// 参数:
//   - menu: 菜单配置
//
// 返回:
//   - error: 验证错误
func ValidateMenu(menu *Menu) error {
	if menu == nil {
		return NewMenuError("NIL_MENU", "菜单不能为空", "")
	}

	if len(menu.Items) == 0 {
		return ErrEmptyMenu
	}

	keyMap := make(map[string]bool)
	for _, item := range menu.Items {
		if item.Key == "" {
			return NewMenuError("EMPTY_KEY", "菜单项键不能为空", "")
		}
		if keyMap[item.Key] {
			return NewMenuError("DUPLICATE_KEY", "菜单项键重复", item.Key)
		}
		keyMap[item.Key] = true
	}

	if menu.Default != "" {
		if !keyMap[menu.Default] && menu.Default != menu.ExitKey {
			return NewMenuError("INVALID_DEFAULT", "默认值无效", menu.Default)
		}
	}

	if menu.AllowExit && menu.ExitKey == "" {
		return NewMenuError("EMPTY_EXIT_KEY", "退出键不能为空", "")
	}

	return nil
}

// RenderMenu 渲染菜单为字符串
//
// 参数:
//   - menu: 菜单配置
//   - style: 菜单样式 (可选, 使用默认样式)
//
// 返回:
//   - string: 渲染后的菜单字符串
func RenderMenu(menu *Menu, style *MenuStyle) string {
	if style == nil {
		style = GetDefaultMenuStyle()
	}

	var builder strings.Builder

	if style.ShowTitle && menu.Title != "" {
		builder.WriteString(menu.Title)
		builder.WriteString("\n")
		builder.WriteString(strings.Repeat("=", len(menu.Title)))
		builder.WriteString("\n")
	}

	indent := strings.Repeat(" ", style.Indent)
	for _, item := range menu.Items {
		builder.WriteString(indent)
		builder.WriteString(style.Prefix)
		builder.WriteString(item.Key)
		builder.WriteString(style.Separator)
		builder.WriteString(item.Value)
		builder.WriteString("\n")
	}

	if menu.AllowExit {
		builder.WriteString(indent)
		builder.WriteString(style.Prefix)
		builder.WriteString(menu.ExitKey)
		builder.WriteString(style.Separator)
		builder.WriteString(menu.ExitText)
		builder.WriteString("\n")
	}

	return builder.String()
}

// ShowMenu 显示结构化菜单并获取用户选择
//
// 参数:
//   - menu: 菜单配置
//   - reader: 输入源
//
// 返回:
//   - string: 用户选择的键
//   - error: 错误信息
func ShowMenu(menu *Menu, reader io.Reader) (string, error) {
	if err := ValidateMenu(menu); err != nil {
		return "", err
	}

	fmt.Print(RenderMenu(menu, nil))

	prompt := menu.Prompt
	if menu.Default != "" {
		prompt = fmt.Sprintf("%s [默认: %s]", menu.Prompt, menu.Default)
	}
	prompt += " "

	choice := Read(reader, prompt)
	if choice == "" && menu.Default != "" {
		return menu.Default, nil
	}

	if choice == "" {
		return "", ErrInvalidChoice
	}

	validKeys := make(map[string]bool)
	for _, item := range menu.Items {
		validKeys[item.Key] = true
	}
	if menu.AllowExit {
		validKeys[menu.ExitKey] = true
	}

	if !validKeys[choice] {
		return "", NewMenuError("INVALID_CHOICE", "无效的选择", choice)
	}

	return choice, nil
}

// ShowMenuLine 便捷函数, 使用 os.Stdin 显示菜单
//
// 参数:
//   - menu: 菜单配置
//
// 返回:
//   - string: 用户选择的键
//   - error: 错误信息
func ShowMenuLine(menu *Menu) (string, error) {
	return ShowMenu(menu, os.Stdin)
}

// ShowMenuLoop 循环显示菜单, 直到用户选择退出
//
// 参数:
//   - menu: 菜单配置
//   - reader: 输入源
//   - handler: 处理函数, 返回 true 继续循环, false 退出
//
// 返回:
//   - error: 错误信息
func ShowMenuLoop(menu *Menu, reader io.Reader, handler func(key string) bool) error {
	if err := ValidateMenu(menu); err != nil {
		return err
	}

	for {
		choice, err := ShowMenu(menu, reader)
		if err != nil {
			return err
		}

		if menu.AllowExit && choice == menu.ExitKey {
			return nil
		}

		if !handler(choice) {
			return nil
		}
	}
}

// ShowMenuLoopLine 便捷函数, 使用 os.Stdin 循环显示菜单
//
// 参数:
//   - menu: 菜单配置
//   - handler: 处理函数, 返回 true 继续循环, false 退出
//
// 返回:
//   - error: 错误信息
func ShowMenuLoopLine(menu *Menu, handler func(key string) bool) error {
	return ShowMenuLoop(menu, os.Stdin, handler)
}
