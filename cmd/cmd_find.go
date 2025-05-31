package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"gitee.com/MM-Q/colorlib"
)

// findCmdMain 是 find 子命令的主函数
func findCmdMain(cl *colorlib.ColorLib, cmd *flag.FlagSet) error {
	// 获取第一个参数作为查找路径
	findPath := cmd.Arg(0)

	// 如果没有指定查找路径, 则使用当前工作目录
	if findPath == "" {
		findPath = "."
	}

	// 检查参数
	if err := checkFindCmdArgs(findPath); err != nil {
		return err
	}

	// 根据标志决定是否启用正则表达式匹配
	var escapedName, escapedPath, exCludedName, exCludedPath string
	if *findCmdRegex {
		// 不转义关键字中的特殊字符, 即启用正则表达式匹配
		escapedName = *findCmdName

		// 不转义要排除的关键字中的特殊字符, 即启用正则表达式匹配
		exCludedName = *findCmdExcludeName

		// 如果指定了排除路径, 则转义路径中的特殊字符
		exCludedPath = *findCmdExcludePath

		// 如果指定了查找路径, 则转义路径中的特殊字符
		escapedPath = *findCmdPath
	} else {
		// 转义关键字中的特殊字符, 即不启用正则表达式匹配
		escapedName = regexp.QuoteMeta(*findCmdName)

		// 转义要排除的关键字中的特殊字符, 即不启用正则表达式匹配
		exCludedName = regexp.QuoteMeta(*findCmdExcludeName)

		// 如果指定了排除路径, 则转义路径中的特殊字符
		exCludedPath = regexp.QuoteMeta(*findCmdExcludePath)

		// 如果指定了查找路径, 则转义路径中的特殊字符
		escapedPath = regexp.QuoteMeta(*findCmdPath)
	}

	// 根据用户选择是否区分大小写
	var nameRegex, exNameRegex, pathRegex, exPathRegex *regexp.Regexp
	var nameRegexErr, exNameRegexErr, pathRegexErr, exPathRegexErr error
	// 默认不区分大小写
	if *findCmdCase {
		// 编译查找文件或目录名的正则表达式
		nameRegex, nameRegexErr = regexp.Compile(escapedName)

		// 编译排除的关键字的正则表达式
		exNameRegex, exNameRegexErr = regexp.Compile(exCludedName)

		// 编译查找路径的正则表达式
		pathRegex, pathRegexErr = regexp.Compile(escapedPath)

		// 编译排除路径的正则表达式
		exPathRegex, exPathRegexErr = regexp.Compile(exCludedPath)
	} else {
		// 不区分大小写的正则表达式
		nameRegex, nameRegexErr = regexp.Compile("(?i)" + escapedName)

		// 编译排除的关键字的正则表达式
		exNameRegex, exNameRegexErr = regexp.Compile("(?i)" + exCludedName)

		// 编译查找路径的正则表达式
		pathRegex, pathRegexErr = regexp.Compile("(?i)" + escapedPath)

		// 编译排除路径的正则表达式
		exPathRegex, exPathRegexErr = regexp.Compile("(?i)" + exCludedPath)
	}
	if nameRegexErr != nil || exNameRegexErr != nil || pathRegexErr != nil || exPathRegexErr != nil {
		return fmt.Errorf("表达式编译错误: %v, %v, %v, %v", nameRegexErr, exNameRegexErr, pathRegexErr, exPathRegexErr)
	}

	// 使用 filepath.WalkDir 遍历目录
	walkDirErr := filepath.WalkDir(findPath, func(path string, entry os.DirEntry, err error) error {
		// 检查路径是否存在
		if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
			// 忽略不存在的路径
			return nil
		}

		// 检查遍历过程中是否遇到错误
		if err != nil {
			return fmt.Errorf("访问文件时出错：%s", err)
		}

		// 跳过*findCmdPath目录本身
		if path == findPath {
			return nil
		}

		// 检查当前路径的深度是否超过最大深度
		depth := strings.Count(filepath.ToSlash(path[len(findPath):]), "/")
		if *findCmdMaxDepth >= 0 && depth > *findCmdMaxDepth {
			return filepath.SkipDir
		}

		// 检查是否为符号链接循环, 然后做相应处理
		if entry.Type()&os.ModeSymlink != 0 {
			if isSymlinkLoop(path) {
				return filepath.SkipDir
			}
		}

		// 如果指定了-n和-p参数, 并且指定了-or参数, 则只检查文件名或路径是否匹配(默认为或操作)
		if *findCmdOr && *findCmdName != "" && *findCmdPath != "" {
			// 执行或操作
			if nameRegex.MatchString(entry.Name()) || pathRegex.MatchString(path) {
				if err := filterConditions(entry, path, cl, exNameRegex, exPathRegex); err != nil {
					return err
				}
			}
			return nil
		}

		// 如果指定了-n和-p参数, 则同时检查文件名和路径是否匹配(默认为-and操作)
		if *findCmdAnd && *findCmdName != "" && *findCmdPath != "" {
			if nameRegex.MatchString(entry.Name()) && pathRegex.MatchString(path) {
				// 如果同时匹配, 则执行筛选条件
				if err := filterConditions(entry, path, cl, exNameRegex, exPathRegex); err != nil {
					return err
				}
			}

			return nil
		}

		// 如果指定了-n参数, 则检查文件名是否匹配
		if *findCmdName != "" {
			if nameRegex.MatchString(entry.Name()) {
				if err := filterConditions(entry, path, cl, exNameRegex, exPathRegex); err != nil {
					return err
				}
			}
			return nil
		}

		// 如果指定了-p参数, 则检查路径是否匹配
		if *findCmdPath != "" {
			if pathRegex.MatchString(path) {
				if err := filterConditions(entry, path, cl, exNameRegex, exPathRegex); err != nil {
					return err
				}
			}
			return nil
		}

		// 默认情况下
		if err := filterConditions(entry, path, cl, exNameRegex, exPathRegex); err != nil {
			return err
		}

		return nil
	})

	if walkDirErr != nil {
		if os.IsPermission(walkDirErr) {
			return fmt.Errorf("权限不足，无法访问某些目录: %v", walkDirErr)
		} else if os.IsNotExist(walkDirErr) {
			return fmt.Errorf("路径不存在: %v", walkDirErr)
		}
		return fmt.Errorf("遍历目录时出错: %v", walkDirErr)
	}

	return nil
}

// 用于在循环中筛选条件的函数
func filterConditions(entry os.DirEntry, path string, cl *colorlib.ColorLib, exNameRegex, exPathRegex *regexp.Regexp) error {
	// 如果只查找文件，跳过目录
	if *findCmdFile && entry.IsDir() {
		return nil
	}

	// 如果只查找目录，跳过文件
	if *findCmdDir && !entry.IsDir() {
		return nil
	}

	// 如果只查找软链接，跳过非软链接
	if *findCmdSymlink {
		fileInfo, linkErr := os.Lstat(path)
		if linkErr != nil {
			return nil
		}

		// Windows系统特殊处理
		if runtime.GOOS == "windows" {
			// 检查.lnk后缀
			if strings.HasSuffix(entry.Name(), ".lnk") {
				// 是快捷方式，继续后续处理
			} else if fileInfo.Mode()&os.ModeSymlink != 0 {
				// 是mklink创建的符号链接
			} else {
				// 既不是.lnk也不是符号链接, 跳过
				return nil
			}
		} else {
			// 非Windows系统只检查ModeSymlink标志
			if fileInfo.Mode()&os.ModeSymlink == 0 {
				return nil // 不是符号链接, 跳过
			}
		}
	}

	// 如果只显示隐藏文件或目录, 跳过非隐藏文件或目录
	if *findCmdHidden && !isHidden(path) && *findCmdHiddenOnly {
		return nil
	}

	// 如果只显示只读文件或目录 跳过非只读文件或目录
	if *findCmdReadOnly && !isReadOnly(path) {
		return nil
	}

	// 如果指定了文件大小, 跳过不符合条件的文件
	if *findCmdSize != "" {
		fileInfo, sizeErr := entry.Info()
		if sizeErr != nil {
			return nil
		}
		if !matchFileSize(fileInfo.Size(), *findCmdSize) {
			return nil
		}
	}

	// 如果指定了修改时间, 跳过不符合条件的文件
	if *findCmdModTime != "" {
		fileInfo, mtimeErr := entry.Info()
		if mtimeErr != nil {
			return nil
		}
		// 检查文件时间是否符合要求
		if !matchFileTime(fileInfo.ModTime(), *findCmdModTime) {
			return nil
		}
	}

	// 如果没有启用隐藏标志且是隐藏目录或文件, 则跳过
	if !*findCmdHidden && isHidden(path) {
		// 如果是隐藏目录，跳过整个目录
		if entry.IsDir() {
			return filepath.SkipDir
		}

		// 如果是隐藏文件，跳过单个文件
		return nil
	}

	// 如果指定了排除文件或目录名, 跳过匹配的文件或目录
	if *findCmdExcludeName != "" && exNameRegex.MatchString(entry.Name()) {
		return nil
	}

	// 如果指定了排除路径, 跳过匹配的路径
	if *findCmdExcludePath != "" && exPathRegex.MatchString(path) {
		return nil
	}

	// 如果启用了delete标志, 删除匹配的文件或目录
	if *findCmdDelete {
		if err := deleteMatchedItem(path, entry.IsDir()); err != nil {
			return err
		}

		// 如果是目录, 跳过整个目录
		if entry.IsDir() {
			return filepath.SkipDir
		}

		// 如果是文件, 跳过单个文件
		return nil
	}

	// 如果启用了-mv标志, 将匹配的文件或目录移动到指定位置
	if *findCmdMove != "" {
		if err := moveMatchedItem(path, *findCmdMove, cl); err != nil {
			return err
		}

		// 如果是目录, 跳过整个目录
		if entry.IsDir() {
			return filepath.SkipDir
		}

		// 如果是文件, 跳过单个文件
		return nil
	}

	// 如果启用了-exec标志, 执行指定的命令
	if *findCmdExec != "" {
		// 替换{}为实际的文件路径，并根据系统类型选择引用方式
		var cmdStr string
		if runtime.GOOS == "windows" {
			cmdStr = strings.Replace(*findCmdExec, "{}", fmt.Sprintf("\"%s\"", path), -1) // Windows使用双引号
		} else {
			cmdStr = strings.Replace(*findCmdExec, "{}", fmt.Sprintf("'%s'", path), -1) // Linux使用单引号
		}

		// 执行命令
		if err := runCommand(cmdStr, cl); err != nil {
			return fmt.Errorf("执行-exec命令时发生了错误: %v", err)
		}

		return nil
	}

	// 根据标志, 输出输出完整路径还是匹配到的路径
	if *findCmdFullPath {
		// 获取完整路径
		fullPath, pathErr := filepath.Abs(path)
		if pathErr != nil {
			return fmt.Errorf("获取完整路径时出错: %s", pathErr)
		}

		// 输出完整路径
		if *findCmdColor {
			if err := printPathColor(fullPath, cl); err != nil {
				return fmt.Errorf("输出路径时出错: %s", err)
			}
		} else {
			// 如果没有启用颜色输出, 直接打印路径
			fmt.Println(fullPath)
		}
	} else {
		// 输出匹配到的路径
		if *findCmdColor {
			if err := printPathColor(path, cl); err != nil {
				return fmt.Errorf("输出路径时出错: %s", err)
			}
		} else {
			// 如果没有启用颜色输出, 直接打印路径
			fmt.Println(path)
		}
	}

	return nil
}

// 用于检查find命令的相关参数是否正确
func checkFindCmdArgs(findPath string) error {
	// 检查要查找的路径是否存在
	if _, err := os.Stat(findPath); err != nil {
		// 检查是否是权限不足的错误
		if os.IsPermission(err) {
			return fmt.Errorf("权限不足，无法访问某些目录: %s", findPath)
		}

		// 如果是不存在错误, 则返回路径不存在
		if os.IsNotExist(err) {
			return fmt.Errorf("路径不存在: %s", findPath)
		}

		// 其他错误, 返回错误信息
		return fmt.Errorf("检查路径时出错: %s: %v", findPath, err)
	}

	// 检查要查找的最大深度是否小于 -1
	if *findCmdMaxDepth < -1 {
		return fmt.Errorf("查找最大深度不能小于 -1")
	}

	// 检查是否同时指定了文件和目录和软链接
	if *findCmdFile && *findCmdDir && *findCmdSymlink {
		return fmt.Errorf("不能同时指定 -f、-d 和 -l 标志")
	}

	// 检查是否同时指定了文件和目录
	if *findCmdFile && *findCmdDir {
		return fmt.Errorf("不能同时指定 -f 和 -d 标志")
	}

	// 检查是否同时指定了文件和软链接
	if *findCmdFile && *findCmdSymlink {
		return fmt.Errorf("不能同时指定 -f 和 -l 标志")
	}

	// 检查是否同时指定了目录和软链接
	if *findCmdDir && *findCmdSymlink {
		return fmt.Errorf("不能同时指定 -d 和 -l 标志")
	}

	// 检查是否同时指定了文件和只读文件
	if *findCmdFile && *findCmdReadOnly {
		return fmt.Errorf("不能同时指定 -f 和 -ro 标志")
	}

	// 检查是否同时指定了目录和只读文件
	if *findCmdDir && *findCmdReadOnly {
		return fmt.Errorf("不能同时指定 -d 和 -ro 标志")
	}

	// 如果只显示隐藏文件或目录, 则必须指定 -H 标志
	if *findCmdHiddenOnly && !*findCmdHidden {
		return fmt.Errorf("必须指定 -H 标志才能使用 -ho 标志")
	}

	// 检查是否同时指定了-or
	if *findCmdOr {
		*findCmdAnd = false // 如果使用-or, 则不能同时使用-and
	}

	// 检查如果指定了文件大小, 格式是否正确(格式为 +5M 或 -5M), 单位必须为 B/K/M/G 同时为大写
	if *findCmdSize != "" {
		// 使用正则表达式匹配文件大小条件
		sizeRegex := regexp.MustCompile(`^([+-])(\d+)([BKMGbkmg])$`) // 正确分组：符号、数字、单位
		match := sizeRegex.FindStringSubmatch(*findCmdSize)          // 查找匹配项
		if match == nil {
			return fmt.Errorf("文件大小格式错误, 格式如+5M(大于5M)或-5M(小于5M), 支持单位B/K/M/G(大写)")
		}
		_, err := strconv.Atoi(match[2]) // 转换数字部分(match[2])
		if err != nil {
			return fmt.Errorf("文件大小格式错误")
		}
	}

	// 检查如果指定了修改时间, 格式是否正确(格式为 +5 或 -5), 单位必须为 天
	if *findCmdModTime != "" {
		// 使用正则表达式匹配文件时间条件
		timeRegex := regexp.MustCompile(`^([+-])(\d+)$`) // 正确分组：符号、数字
		match := timeRegex.FindStringSubmatch(*findCmdModTime)
		if match == nil {
			return fmt.Errorf("文件时间格式错误, 格式如+5(5天前)或-5(5天内)")
		}
		_, err := strconv.Atoi(match[2]) // 转换数字部分(match[2])
		if err != nil {
			return fmt.Errorf("文件时间格式错误")
		}
	}

	// 检查-exec标志是否包含{}
	if *findCmdExec != "" && !strings.Contains(*findCmdExec, "{}") {
		return fmt.Errorf("使用-exec标志时必须包含{}作为路径占位符")
	}

	// 检查-print-cmd标志是否与-exec一起使用
	if *findCmdPrintCmd && *findCmdExec == "" {
		return fmt.Errorf("使用-print-cmd标志时必须同时指定-exec标志")
	}

	// 检查-delete标志是否与-exec一起使用
	if *findCmdDelete && *findCmdExec != "" {
		return fmt.Errorf("使用-delete标志时不能同时指定-exec标志")
	}

	// 检查-mv标志是否与-exec一起使用
	if *findCmdMove != "" && *findCmdExec != "" {
		return fmt.Errorf("使用-mv标志时不能同时指定-exec标志")
	}

	// 检查-mv标志是否与-delete一起使用
	if *findCmdMove != "" && *findCmdDelete {
		return fmt.Errorf("使用-mv标志时不能同时指定-delete标志")
	}

	// 检查-mv标志指定的路径是否为文件
	if *findCmdMove != "" {
		if info, err := os.Stat(*findCmdMove); err == nil {
			// 如果指定的路径是文件, 则返回错误
			if !info.IsDir() {
				return fmt.Errorf("-mv标志指定的路径必须为目录")
			}
		}
	}

	return nil
}

// matchFileSize 检查文件大小是否符合指定的条件
func matchFileSize(fileSize int64, sizeCondition string) bool {
	if len(sizeCondition) < 2 {
		return false
	}

	// 获取比较符号和数值部分
	comparator := sizeCondition[0]
	sizeStr := sizeCondition[1:]

	// 获取单位
	unit := sizeStr[len(sizeStr)-1]
	sizeValueStr := sizeStr[:len(sizeStr)-1]

	// 转换数值部分
	sizeValue, err := strconv.ParseFloat(sizeValueStr, 64)
	if err != nil {
		return false
	}

	// 根据单位转换为字节
	var sizeInBytes float64
	switch unit {
	case 'B':
		sizeInBytes = sizeValue
	case 'b':
		sizeInBytes = sizeValue
	case 'K':
		sizeInBytes = sizeValue * 1024
	case 'k':
		sizeInBytes = sizeValue * 1024
	case 'M':
		sizeInBytes = sizeValue * 1024 * 1024
	case 'm':
		sizeInBytes = sizeValue * 1024 * 1024
	case 'G':
		sizeInBytes = sizeValue * 1024 * 1024 * 1024
	case 'g':
		sizeInBytes = sizeValue * 1024 * 1024 * 1024
	default:
		return false
	}

	// 根据比较符号进行比较
	switch comparator {
	case '+':
		return float64(fileSize) > sizeInBytes
	case '-':
		return float64(fileSize) < sizeInBytes
	default:
		return false
	}
}

// matchFileTime 检查文件时间是否符合指定的条件
func matchFileTime(fileTime time.Time, timeCondition string) bool {
	// 检查时间条件是否为空
	if len(timeCondition) < 2 {
		return false
	}

	// 获取比较符号和数值部分
	comparator := timeCondition[0]
	daysStr := timeCondition[1:]

	// 转换天数
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		return false
	}

	// 计算时间阈值
	threshold := time.Now().AddDate(0, 0, -days)

	// 根据比较符号进行比较
	switch comparator {
	case '+':
		return fileTime.After(threshold) // 检查文件时间是否在阈值之后
	case '-':
		return fileTime.Before(threshold) // 检查文件时间是否在阈值之前
	default:
		return false
	}
}

// runCommand 执行单个命令， 支持跨平台并检查shell是否存在
// 参数cmdStr是要执行的命令字符串
// 参数cl是颜色库实例
// 返回执行过程中的错误
func runCommand(cmdStr string, cl *colorlib.ColorLib) error {
	var shell string
	var args []string

	// 根据系统选择默认shell
	if runtime.GOOS == "windows" {
		// 先尝试使用 PowerShell
		if _, err := exec.LookPath("powershell"); err == nil {
			shell = "powershell"
			args = []string{"-Command", cmdStr}
		} else {
			// 如果 PowerShell 不存在，使用 cmd
			shell = "cmd"
			args = []string{"/C", cmdStr}
		}
	} else {
		shell = "bash"
		args = []string{"-c", cmdStr}
	}

	// 检查shell是否存在
	if _, err := exec.LookPath(shell); err != nil {
		return fmt.Errorf("找不到 %s 解释器: %v", shell, err)
	}

	// 如果启用了print-cmd输出, 打印执行的命令
	if *findCmdPrintCmd {
		cl.Redf("%s %s", shell, strings.Join(args, " "))
	}

	// 构建命令并设置输出
	cmd := exec.Command(shell, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 执行命令并捕获错误
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("命令执行失败: %v", err)
	}
	return nil
}

// deleteMatchedItem 删除匹配的文件或目录
// 参数path是要删除的路径
// 参数isDir指示是否是目录
// 参数cl是颜色库实例
// 返回删除过程中的错误
func deleteMatchedItem(path string, isDir bool) error {
	// 检查是否为空路径
	if path == "" {
		return fmt.Errorf("没有可删除的路径")
	}

	// 先检查文件/目录是否存在
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("文件/目录不存在: %s", path)
		}
		return fmt.Errorf("检查文件/目录时出错: %s: %v", path, err)
	}

	var rmErr error

	// 根据类型选择删除方法
	if isDir {
		// 检查目录是否为空
		dirEntries, readDirErr := os.ReadDir(path)
		if readDirErr != nil {
			return fmt.Errorf("读取目录失败: %s: %v", path, readDirErr)
		}

		// 检查目录是否为空
		if len(dirEntries) > 0 {
			// 目录不为空, 递归删除
			rmErr = os.RemoveAll(path)
		} else {
			// 删除空目录
			rmErr = os.Remove(path)
		}
	} else {
		rmErr = os.Remove(path)
	}

	if rmErr != nil {
		return fmt.Errorf("删除失败: %s: %v", path, rmErr)
	}

	return nil
}

// moveMatchedItem 移动匹配的文件或目录到指定位置
// 参数path是要移动的源路径
// 参数destPath是目标路径
// 返回移动过程中的错误
func moveMatchedItem(path string, targetPath string, cl *colorlib.ColorLib) error {
	// 检查源路径是否为空
	if path == "" {
		return fmt.Errorf("源路径为空")
	}

	// 检查目标路径是否为空
	if targetPath == "" {
		return fmt.Errorf("没有指定目标路径")
	}

	// 检查源路径是否存在
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("源文件/目录不存在: %s", path)
		}
		return fmt.Errorf("检查源文件/目录时出错: %s: %v", path, err)
	}

	// 获取目标路径的绝对路径
	absTargetPath, err := filepath.Abs(targetPath)
	if err != nil {
		return fmt.Errorf("获取目标路径绝对路径失败: %v", err)
	}

	// 获取源路径的绝对路径
	absSearchPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("获取源路径绝对路径失败: %v", err)
	}

	// 检查目标路径是否是源路径的子目录(防止循环移动)
	if strings.HasPrefix(absTargetPath, absSearchPath) {
		return fmt.Errorf("不能将目录移动到自身或其子目录中")
	}

	// 确保目标目录存在
	if err := os.MkdirAll(filepath.Dir(absTargetPath), 0755); err != nil {
		return fmt.Errorf("创建目标目录失败: %v", err)
	}

	// 组装完整的目标路径: 目标路径 + 源路径的文件名
	if filepath.Base(absSearchPath) != "" {
		absTargetPath = filepath.Join(absTargetPath, filepath.Base(absSearchPath))
	}

	// 尝试移动前先检查权限
	if err := os.WriteFile(filepath.Join(filepath.Dir(absTargetPath), ".fck_tmp"), []byte{}, 0644); err != nil {
		return fmt.Errorf("目标目录无写入权限: %v", err)
	}
	os.Remove(filepath.Join(filepath.Dir(absTargetPath), ".fck_tmp"))

	// 检查目标文件是否已存在
	if _, err := os.Stat(absTargetPath); err == nil {
		// 如果是移动操作，直接跳过而不是报错
		if *findCmdMove != "" {
			return nil
		}
		return fmt.Errorf("目标文件已存在: %s", absTargetPath)
	}

	// 打印移动信息
	if *findCmdPrintMove {
		cl.Redf("%s -> %s", absSearchPath, absTargetPath)
	}

	// 执行移动操作
	if err := os.Rename(absSearchPath, absTargetPath); err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("目标文件已存在: %s -> %s", absSearchPath, absTargetPath)
		}
		if os.IsPermission(err) {
			return fmt.Errorf("权限不足，无法移动文件: %v", err)
		}

		return fmt.Errorf("移动失败: %s -> %s: %v", absSearchPath, absTargetPath, err)
	}

	return nil
}

// isSymlinkLoop 检查路径是否是符号链接循环
func isSymlinkLoop(path string) bool {
	visited := make(map[string]bool)
	for {
		if visited[path] {
			return true // 检测到循环
		}
		visited[path] = true

		info, err := os.Lstat(path)
		if err != nil || info.Mode()&os.ModeSymlink == 0 {
			return false // 不是符号链接或出错
		}

		newPath, err := os.Readlink(path)
		if err != nil {
			return false
		}

		if !filepath.IsAbs(newPath) {
			newPath = filepath.Join(filepath.Dir(path), newPath)
		}
		path = newPath
	}
}
