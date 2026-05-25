package which

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gitee.com/MM-Q/shellx"
)

// windowsExts 定义 Windows 可执行文件扩展名集合
var windowsExts = map[string]bool{
	".exe": true,
	".bat": true,
	".cmd": true,
	".com": true,
}

// WhichConfig 配置结构体
type WhichConfig struct {
	Commands []string // 要查找的命令列表（支持多个）
	All      bool     // 显示所有匹配的路径（而不仅是第一个）
	Silent   bool     // 静默模式，不输出，只返回退出码
}

// WhichCmdMain 执行 which 命令
//
// 参数:
//   - config: 命令配置，包含要查找的命令列表、是否显示所有匹配、是否静默模式
//
// 返回值:
//   - error: 查找过程中的错误，nil 表示所有命令都找到
func WhichCmdMain(config WhichConfig) error {
	if len(config.Commands) == 0 {
		return fmt.Errorf("no command specified")
	}

	var notFound []string
	var results []string

	for _, cmd := range config.Commands {
		// 查找命令路径
		paths, err := findCommand(cmd, config.All)
		if err != nil {
			return err
		}

		// 处理查找结果
		if len(paths) == 0 {
			notFound = append(notFound, cmd)
		} else {
			results = append(results, paths...)
		}
	}

	// 非静默模式下，输出找到的路径或错误信息
	if !config.Silent {
		for _, path := range results {
			fmt.Println(path)
		}
	}

	// 如果有命令未找到，返回错误信息
	if len(notFound) > 0 {
		return fmt.Errorf("%s: command not found", notFound[0])
	}

	return nil
}

// findCommand 查找单个命令的可执行路径
// 优先使用 FindCommandPath（基于 exec.LookPath），需要 -a 时遍历所有 PATH 目录
//
// 参数:
//   - name: 命令名称
//   - all: 是否返回所有匹配路径
//
// 返回值:
//   - []string: 找到的路径列表
//   - error: 查找过程中的系统错误
func findCommand(name string, all bool) ([]string, error) {
	// 如果不需要 -a（所有匹配），直接使用 FindCommandPath
	if !all {
		path := shellx.FindCommandPath(name)
		if path == "" {
			// 没找到
			return nil, nil
		}
		return []string{path}, nil
	}

	// 需要 -a：遍历所有 PATH 目录返回所有匹配
	var matches []string

	// 如果命令包含路径分隔符，直接检查该路径
	if strings.ContainsAny(name, `/\`) {
		// 直接检查该路径是否可执行
		if isExecutable(name) {
			abs, err := filepath.Abs(name)
			if err != nil {
				return nil, err
			}
			return []string{abs}, nil
		}
		return nil, nil
	}

	// 从 PATH 环境变量中查找
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return nil, nil
	}

	// 分隔 PATH 环境变量中的目录，并遍历每个目录
	separator := string(os.PathListSeparator)
	paths := strings.Split(pathEnv, separator)

	for _, dir := range paths {
		dir = strings.TrimSpace(dir)
		if dir == "" {
			continue
		}

		// 在当前目录中查找可执行文件
		found := findExecutableInDir(name, dir)

		// 添加到结果列表中
		for _, f := range found {
			abs, err := filepath.Abs(f)
			if err != nil {
				continue
			}
			matches = append(matches, abs)
		}
	}

	return matches, nil
}

// findExecutableInDir 在指定目录中查找可执行文件
//
// 参数:
//   - name: 命令名称
//   - dir: 目录路径
//
// 返回值:
//   - []string: 找到的可执行文件完整路径列表
func findExecutableInDir(name, dir string) []string {
	// Unix/Linux/macOS: 直接检查文件
	if runtime.GOOS != "windows" {
		fullPath := filepath.Join(dir, name)
		if isExecutable(fullPath) {
			return []string{fullPath}
		}
		return nil
	}

	// Windows: 需要处理扩展名
	return findWindowsExecutables(name, dir)
}

// findWindowsExecutables 在 Windows 目录中查找可执行文件
//
// 参数:
//   - name: 命令名称
//   - dir: 目录路径
//
// 返回值:
//   - []string: 找到的可执行文件完整路径列表
func findWindowsExecutables(name, dir string) []string {
	var results []string

	// 获取文件扩展名（小写）
	ext := strings.ToLower(filepath.Ext(name))

	// 情况1: 命令已包含可执行扩展名（如 go.exe）
	if windowsExts[ext] {
		fullPath := filepath.Join(dir, name)
		if isExecutable(fullPath) {
			results = append(results, fullPath)
		}
		return results
	}

	// 情况2: 命令无扩展名（如 go），尝试添加各种扩展名
	for ext := range windowsExts {
		fullPath := filepath.Join(dir, name+ext)
		if isExecutable(fullPath) {
			results = append(results, fullPath)
		}
	}

	return results
}

// isExecutable 检查文件是否可执行
//
// 参数:
//   - path: 文件路径
//
// 返回值:
//   - bool: 是否可执行
func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	if info.IsDir() {
		return false
	}

	// Windows: 文件存在且有可执行扩展名即可
	if runtime.GOOS == "windows" {
		return true
	}

	// Unix/Linux/macOS: 检查可执行权限
	mode := info.Mode()
	return mode&0111 != 0
}
