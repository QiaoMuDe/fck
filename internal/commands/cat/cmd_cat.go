package cat

import (
	"fmt"

	"gitee.com/MM-Q/go-kit/term"
)

// CatConfig cat 命令配置
type CatConfig struct {
	// CLI 参数
	Targets      []string // 目标文件列表
	ShowLineNum  bool     // -n 显示所有行号
	ShowNonBlank bool     // -b 显示非空行行号
	ShowEnd      bool     // -E 显示行尾$
	ShowTabs     bool     // -T 将制表符显示为^I
	ShowAll      bool     // -A 等价于 -ET
	Quiet        bool     // -q 静默模式 (不显示错误信息)
	Text         bool     // -a, --text 强制将二进制文件视为文本处理
	IgnoreBinary bool     // -I, --ignore-binary 完全忽略二进制文件
	UseLess      bool     // -l, --less 使用分页器查看文件内容
	Highlight    bool     // -H, --highlight 启用语法高亮
	MaxSize      int64    // -S, --max-size 最大文件大小 (字节)
}

// CatCmdMain 执行 cat 命令
//
// 参数:
//   - config: 命令配置
//
// 返回:
//   - error: 执行错误 (如果有)
func CatCmdMain(config CatConfig) error {
	// 1. 处理标志冲突 (-b 优先级高于 -n)
	if config.ShowNonBlank {
		config.ShowLineNum = false
	}

	// 2. 处理 --show-all
	if config.ShowAll {
		config.ShowEnd = true
		config.ShowTabs = true
	}

	// 3. 检测是否为管道输入
	isPipe := term.IsStdinPipe()

	// 4. 如果没有文件参数且不是管道输入，报错
	if len(config.Targets) == 0 && !isPipe {
		return fmt.Errorf("no file specified")
	}

	// 如果为管道输入或重定向则禁用高亮模式
	if isPipe {
		config.Highlight = false
	}

	// 5. 获取内容源
	var source ContentSource
	if isPipe {
		source = NewStdinSource("", config.MaxSize)
	} else {
		// 只处理第一个文件
		source = NewFileSource(config.Targets[0], config.MaxSize)
	}

	// 6. 处理内容
	processor := NewProcessor(config)
	content, err := processor.Process(source)
	if err != nil {
		return err
	}

	// 7. 如果内容为空（如二进制文件被跳过），直接返回
	if len(content) == 0 {
		return nil
	}

	// 8. 根据模式输出
	if config.UseLess {
		return OutputWithPager(content, source.Name(), config)
	}

	return OutputDirectly(content, source.Name(), config)
}
