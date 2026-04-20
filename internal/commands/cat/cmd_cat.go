package cat

import (
	"fmt"
	"os"

	"gitee.com/MM-Q/fck/internal/utils"
)

// CatConfig cat 命令配置
type CatConfig struct {
	// CLI 参数
	Targets      []string // 目标文件列表
	ShowLineNum  bool     // -n 显示所有行号
	ShowNonBlank bool     // -b 显示非空行行号
	ShowEnd      bool     // -E 显示行尾$
	ShowTabs     bool     // -T 显示制表符为^I
	ShowAll      bool     // -A 等价于 -ET
	Quiet        bool     // -q 静默模式 (不显示错误信息)
	Text         bool     // -a, --text 强制将二进制文件视为文本处理
	IgnoreBinary bool     // -I, --ignore-binary 完全忽略二进制文件
	UseLess      bool     // -l, --less 使用分页器查看文件内容
	Highlight    bool     // -H, --highlight 启用语法高亮
	MaxSize      int64    // -S, --max-size 最大文件大小 (字节)

	// 运行时
	LineCounter int // 行号计数器
}

// CatCmdMain 执行 cat 命令
//
// 参数:
//   - config: 命令配置
//
// 返回:
//   - error: 执行错误 (如果有)
func CatCmdMain(config CatConfig) error {
	// 1. 验证参数
	if len(config.Targets) == 0 {
		return fmt.Errorf("no file specified")
	}

	// 2. 如果使用分页器模式，直接调用分页器查看
	if config.UseLess {
		return runOVMode(config)
	}

	// 3. 处理标志冲突 (-b 优先级高于 -n)
	if config.ShowNonBlank {
		config.ShowLineNum = false
	}

	// 4. 处理 --show-all
	if config.ShowAll {
		config.ShowEnd = true
		config.ShowTabs = true
	}

	// 5. 处理文件
	return processFile(config.Targets[0], &config)
}

// processFile 处理单个文件
//
// 参数:
//   - path: 文件路径
//   - config: 命令配置
//
// 返回:
//   - error: 处理错误 (如果有)
func processFile(path string, config *CatConfig) error {
	// 1. 基础检查（目录、二进制）
	if err := checkFile(path, config); err != nil {
		return err
	}

	// 2. 创建统一查看器
	viewer := NewFileViewer(config)

	// 3. 统一查看（内部处理 head/tail/高亮等逻辑）
	return viewer.View(path)
}

// checkFile 检查文件类型
//
// 参数:
//   - path: 文件路径
//   - config: 命令配置
//
// 返回:
//   - error: 检查错误，如果是二进制文件且需要跳过，返回 nil
func checkFile(path string, config *CatConfig) error {
	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer func() {
		_ = file.Close()
	}()

	// 获取文件信息 (用于判断是否是目录)
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info %s: %w", path, err)
	}
	if info.IsDir() {
		return fmt.Errorf("%s is a directory", path)
	}

	// 二进制文件检测 (除非强制文本模式)
	if !config.Text {
		isBinary, err := utils.IsBinaryFile(file)
		if err != nil {
			return fmt.Errorf("cannot detect file type for %s: %w", path, err)
		}

		// 处理二进制文件
		if isBinary {
			// -I 模式：静默跳过
			if config.IgnoreBinary {
				return nil
			}

			// 默认行为：输出提示并跳过
			if !config.Quiet {
				fmt.Printf("bin file %s matches\n", path)
			}
			return nil
		}
	}

	return nil
}
