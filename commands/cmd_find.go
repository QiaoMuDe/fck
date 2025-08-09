package commands

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/globals"
)

// findCmdMain 是 find 子命令的主函数
func findCmdMain(cl *colorlib.ColorLib) error {
	// 获取第一个参数作为查找路径
	findPath := findCmd.Arg(0)

	// 如果没有指定查找路径, 则使用当前工作目录
	if findPath == "" {
		findPath = "."
	}

	// 清理路径
	findPath = filepath.Clean(findPath)

	// 检查参数
	if err := checkFindCmdArgs(findPath); err != nil {
		return err
	}

	// 根据findCmdColor.Get()选项设置颜色
	if findCmdColor.Get() {
		cl.NoColor.Store(false)
	} else {
		cl.NoColor.Store(true)
	}

	// 准备正则表达式模式
	isRegex := findCmdRegex.Get()       // 是否启用正则表达式
	wholeWord := findCmdWholeWord.Get() // 是否匹配完整关键字
	caseSensitive := findCmdCase.Get()  // 是否区分大小写

	// 定义正则表达式和错误处理
	var (
		nameRegex, exNameRegex, pathRegex, exPathRegex *regexp.Regexp
		nameErr, exNameErr, pathErr, exPathErr         error
	)

	// 启用正则模式
	if isRegex {
		// 构建检索文件名的正则表达式模式
		escapedName := RegexBuilder(findCmdName.Get(), isRegex, wholeWord, caseSensitive)

		// 构建排除文件名的正则表达式模式
		excludedName := RegexBuilder(findCmdExcludeName.Get(), isRegex, wholeWord, caseSensitive)

		// 构建检索路径的正则表达式模式
		escapedPath := RegexBuilder(findCmdPath.Get(), isRegex, wholeWord, caseSensitive)

		// 构建排除路径的正则表达式模式
		excludedPath := RegexBuilder(findCmdExcludePath.Get(), isRegex, wholeWord, caseSensitive)

		// 构建文件名的正则表达式
		if nameRegex, nameErr = CompileRegex(escapedName); nameErr != nil {
			return fmt.Errorf("文件名正则表达式编译错误: %v", nameErr)
		}

		// 构建排除文件名的正则表达式
		if exNameRegex, exNameErr = CompileRegex(excludedName); exNameErr != nil {
			return fmt.Errorf("排除文件名正则表达式编译错误: %v", exNameErr)
		}

		// 构建路径的正则表达式
		if pathRegex, pathErr = CompileRegex(escapedPath); pathErr != nil {
			return fmt.Errorf("路径正则表达式编译错误: %v", pathErr)
		}

		// 构建排除路径的正则表达式
		if exPathRegex, exPathErr = CompileRegex(excludedPath); exPathErr != nil {
			return fmt.Errorf("排除路径正则表达式编译错误: %v", exPathErr)
		}
	} else {
		// 禁用正则匹配
		nameRegex = nil
		pathRegex = nil
		exNameRegex = nil
		exPathRegex = nil
	}

	// 定义一个统计匹配项的
	matchCount := atomic.Int64{}
	matchCount.Store(0)

	// 创建FindConfig实例
	config := &globals.FindConfig{
		Cl:            cl,                       // 颜色库实例
		NameRegex:     nameRegex,                // 文件名匹配正则
		ExNameRegex:   exNameRegex,              // 排除文件名正则
		PathRegex:     pathRegex,                // 路径匹配正则
		ExPathRegex:   exPathRegex,              // 排除路径正则
		IsRegex:       isRegex,                  // 是否启用正则匹配
		WholeWord:     wholeWord,                // 是否全词匹配
		CaseSensitive: caseSensitive,            // 是否区分大小写
		MatchCount:    &matchCount,              // 匹配计数
		NamePattern:   findCmdName.Get(),        // 文件名匹配模式
		PathPattern:   findCmdPath.Get(),        // 路径匹配模式
		ExNamePattern: findCmdExcludeName.Get(), // 排除文件名模式
		ExPathPattern: findCmdExcludePath.Get(), // 排除路径模式
	}

	// 检查是否指定了ext参数, 如果指定则存储到config.FindExtSliceMap中
	if findCmdExt.Len() > 0 {
		// 遍历ext切片
		for _, ext := range findCmdExt.Get() {
			// 检查扩展名是否包含常见危险字符（空格、换行符、制表符、特殊路径分隔符等）
			if strings.ContainsAny(ext, " \t\n\r\\/:*?\"<>|") {
				return fmt.Errorf("扩展名包含非法字符: %s", ext)
			}
			// 如果扩展名不包含"."则添加"."
			if !strings.HasPrefix(ext, ".") {
				ext = fmt.Sprint(".", ext)
			}
			// 存储ext标志
			config.FindExtSliceMap.Store(ext, true)
		}
	}

	// 检查是否启用并发模式
	if findCmdX.Get() {
		// 启用并发模式
		if err := processWalkDirConcurrent(config, findPath); err != nil {
			return err
		}
	} else {
		// 运行processWalkDir单线程遍历模式
		if walkDirErr := processWalkDir(config, findPath); walkDirErr != nil {
			return walkDirErr
		}
	}

	// 如果启用了count标志, 只输出匹配数量
	if findCmdCount.Get() {
		fmt.Println(matchCount.Load())
		return nil
	}

	return nil
}

// processWalkDir 用于处理 filepath.WalkDir 的逻辑
// 参数:
// - config: 查找配置结构体
// - findPath: 要查找的路径
// 返回值:
// - error: 错误信息
func processWalkDir(config *globals.FindConfig, findPath string) error {
	// 使用 filepath.WalkDir 遍历目录
	walkDirErr := filepath.WalkDir(findPath, func(path string, entry os.DirEntry, err error) error {
		// 检查遍历过程中是否遇到错误
		if err != nil {
			// 忽略不存在的报错
			if os.IsNotExist(err) {
				// 路径不存在, 跳过
				return nil
			}

			// 检查是否为权限不足的报错
			if os.IsPermission(err) {
				config.Cl.PrintErrf("权限不足, 无法访问某些目录: %s\n", path)
				return nil
			}

			return fmt.Errorf("访问时出错：%s", err)
		}

		// 跳过findCmdPath.Get()目录本身
		if path == findPath {
			return nil
		}

		// 检查当前路径的深度是否超过最大深度(先将路径统一转为/分隔符)
		depth := strings.Count(path[len(findPath):], string(filepath.Separator))
		if findCmdMaxDepth.Get() >= 0 && depth > findCmdMaxDepth.Get() {
			return filepath.SkipDir
		}

		// 检查是否为符号链接循环, 然后做相应处理
		if entry.Type()&os.ModeSymlink != 0 {
			if isSymlinkLoop(path) {
				return filepath.SkipDir
			}
		}

		// 处理文件或目录
		if processErr := processFindCmd(config, entry, path); processErr != nil {
			return processErr
		}

		return nil
	})

	// 检查遍历过程中是否遇到错误
	if walkDirErr != nil {
		if os.IsPermission(walkDirErr) {
			return fmt.Errorf("权限不足, 无法访问某些目录: %v", walkDirErr)
		} else if os.IsNotExist(walkDirErr) {
			return fmt.Errorf("路径不存在: %v", walkDirErr)
		}
		return fmt.Errorf("遍历目录时出错: %v", walkDirErr)
	}

	return nil
}

// processFindCmd 用于处理 find 命令的逻辑
// 参数:
// - config: 查找配置结构体
// - entry: 文件或目录的DirEntry对象
// - path: 文件或目录的完整路径
// 返回值:
// - error: 错误信息
func processFindCmd(config *globals.FindConfig, entry os.DirEntry, path string) error {
	// 如果指定了-n和-p参数, 并且指定了-or参数, 则只检查文件名或路径是否匹配(默认为或操作)
	if findCmdOr.Get() && config.NamePattern != "" && config.PathPattern != "" {
		// 执行或操作
		if matchPattern(entry.Name(), config.NamePattern, config.NameRegex, config) || matchPattern(path, config.PathPattern, config.PathRegex, config) {
			if err := filterConditions(config, entry, path); err != nil {
				return err
			}
		}
		return nil
	}

	// 如果指定了-n和-p参数, 则同时检查文件名和路径是否匹配(默认为-and操作)
	if findCmdAnd.Get() && config.NamePattern != "" && config.PathPattern != "" {
		if matchPattern(entry.Name(), config.NamePattern, config.NameRegex, config) && matchPattern(path, config.PathPattern, config.PathRegex, config) {
			// 如果同时匹配, 则执行筛选条件
			if err := filterConditions(config, entry, path); err != nil {
				return err
			}
		}

		return nil
	}

	// 如果指定了-n参数, 则检查文件名是否匹配
	if config.NamePattern != "" {
		if matchPattern(entry.Name(), config.NamePattern, config.NameRegex, config) {
			if err := filterConditions(config, entry, path); err != nil {
				return err
			}
		}
		return nil
	}

	// 如果指定了-p参数, 则检查路径是否匹配
	if config.PathPattern != "" {
		if matchPattern(path, config.PathPattern, config.PathRegex, config) {
			if err := filterConditions(config, entry, path); err != nil {
				return err
			}
		}
		return nil
	}

	// 默认情况下
	if err := filterConditions(config, entry, path); err != nil {
		return err
	}

	return nil
}

// filterConditions 用于在循环中筛选条件的函数
// 参数:
// - config: 查找配置结构体
// - entry: 文件或目录的DirEntry对象
// - path: 文件或目录的完整路径
// 返回值:
// - error: 错误信息
func filterConditions(config *globals.FindConfig, entry os.DirEntry, path string) error {
	// 默认隐藏文件或隐藏目录不参与匹配
	// 如果没有启用隐藏标志且是隐藏目录或文件, 则跳过
	if !findCmdHidden.Get() {
		if isHidden(path) {
			// 如果是隐藏目录, 跳过整个目录
			if entry.IsDir() {
				return filepath.SkipDir
			}

			// 如果是隐藏文件, 跳过单个文件
			return nil
		}
	}

	// 如果指定了排除文件或目录名, 跳过匹配的文件或目录
	if config.ExNamePattern != "" {
		if matchPattern(entry.Name(), config.ExNamePattern, config.ExNameRegex, config) {
			if entry.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}
	}

	// 如果指定了排除路径, 跳过匹配的路径
	if config.ExPathPattern != "" {
		if matchPattern(path, config.ExPathPattern, config.ExPathRegex, config) {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
	}

	// 仅在需要文件元信息时才获取（空文件检查/修改时间/文件大小）
	var cacheInfo fs.FileInfo
	var cacheErr error

	// 判断是否需要获取文件信息：空文件检查、修改时间筛选或大小筛选
	if (findCmdType.Get() == globals.FindTypeEmpty || findCmdType.Get() == globals.FindTypeEmptyShort) || findCmdModTime.Get() != "" || findCmdSize.Get() != "" {
		cacheInfo, cacheErr = entry.Info()
		if cacheErr != nil {
			return nil
		}
	}

	// 预计算文件扩展名并缓存, 避免多次调用filepath.Ext()
	entryExt := filepath.Ext(entry.Name())

	// 根据findCmdType的值, 跳过不符合要求的文件或目录
	switch findCmdType.Get() {
	case globals.FindTypeFile, globals.FindTypeFileShort:
		// 如果只查找文件, 跳过目录
		if entry.IsDir() {
			return nil
		}
	case globals.FindTypeDir, globals.FindTypeDirShort:
		// 如果只查找目录, 跳过文件
		if !entry.IsDir() {
			return nil
		}
	case globals.FindTypeSymlink, globals.FindTypeSymlinkShort:
		// 如果只查找软链接, 跳过非软链接

		// Windows系统特殊处理
		if runtime.GOOS == "windows" {
			// 如果不是.lnk或.url文件, 跳过
			if !globals.WindowsSymlinkExts[entryExt] {
				return nil
			}
		} else {
			// 非Windows系统检查符号链接
			if entry.Type()&os.ModeSymlink == 0 {
				return nil // 不是符号链接, 跳过
			}
		}

	case globals.FindTypeHidden, globals.FindTypeHiddenShort:
		// 如果只显示隐藏文件或目录, 跳过非隐藏文件或目录
		if !isHidden(path) {
			return nil
		}
	case globals.FindTypeReadonly, globals.FindTypeReadonlyShort:
		// 如果只显示只读文件或目录 跳过非只读文件或目录
		if !isReadOnly(path) {
			return nil
		}
	case globals.FindTypeEmpty, globals.FindTypeEmptyShort:
		// 如果只查找空文件或目录
		if entry.IsDir() {
			// 检查目录是否为空
			dirEntries, err := os.ReadDir(path)
			if err != nil || len(dirEntries) > 0 {
				return nil
			}
		} else {
			// 如果是文件, 检查文件是否为空
			if cacheInfo.Size() > 0 {
				return nil
			}
		}
	case globals.FindTypeExecutable, globals.FindTypeExecutableShort:
		// 如果只查找可执行文件或目录
		if runtime.GOOS == "windows" {
			// Windows系统检查文件扩展名
			ext := strings.ToLower(entryExt)
			if !globals.WindowsExecutableExts[ext] {
				return nil
			}
		} else {
			// 非Windows系统检查可执行权限
			if entry.Type().IsRegular() && entry.Type()&0111 == 0 {
				return nil
			}
		}
	case globals.FindTypeSocket, globals.FindTypeSocketShort:
		// Windows不支持Unix domain socket, 直接返回nil
		if runtime.GOOS == "windows" {
			return nil
		}

		// 如果只查找socket文件, 跳过非socket文件
		if entry.Type()&os.ModeSocket == 0 {
			return nil
		}
	case globals.FindTypePipe, globals.FindTypePipeShort:
		// Windows不支持Unix命名管道, 直接返回nil
		if runtime.GOOS == "windows" {
			return nil
		}

		// 如果只查找管道文件, 跳过非管道文件
		if entry.Type()&os.ModeNamedPipe == 0 {
			return nil
		}
	case globals.FindTypeBlock, globals.FindTypeBlockShort:
		// Windows不支持Unix块设备文件, 直接返回nil
		if runtime.GOOS == "windows" {
			return nil
		}

		// 如果只查找块设备文件, 跳过非块设备文件
		if entry.Type()&os.ModeDevice == 0 {
			return nil
		}
	case globals.FindTypeChar, globals.FindTypeCharShort:
		// Windows不支持Unix字符设备文件, 直接返回nil
		if runtime.GOOS == "windows" {
			return nil
		}

		// 如果只查找字符设备文件, 跳过非字符设备文件
		if entry.Type()&os.ModeCharDevice == 0 {
			return nil
		}
	case globals.FindTypeAppend, globals.FindTypeAppendShort:
		// Windows文件系统不支持Unix追加模式标志, 直接返回nil
		if runtime.GOOS == "windows" {
			return nil
		}

		// 如果只查找追加模式文件, 跳过非追加模式文件
		if entry.Type()&os.ModeAppend == 0 {
			return nil
		}
	case globals.FindTypeNonAppend, globals.FindTypeNonAppendShort:
		// Windows文件系统不支持Unix追加模式标志, 直接返回nil
		if runtime.GOOS == "windows" {
			return nil
		}

		// 如果只查找非追加模式文件, 跳过追加模式文件
		if entry.Type()&os.ModeAppend != 0 {
			return nil
		}
	case globals.FindTypeExclusive, globals.FindTypeExclusiveShort:
		// Windows文件系统不支持Unix独占模式标志, 直接返回nil
		if runtime.GOOS == "windows" {
			return nil
		}

		// 如果只查找独占模式文件, 跳过非独占模式文件
		if entry.Type()&os.ModeExclusive == 0 {
			return nil
		}
	}

	// 如果指定了文件大小, 跳过不符合条件的文件
	if findCmdSize.Get() != "" {
		if !matchFileSize(cacheInfo.Size(), findCmdSize.Get()) {
			return nil
		}
	}

	// 如果指定了修改时间, 跳过不符合条件的文件
	if findCmdModTime.Get() != "" {
		// 检查文件时间是否符合要求
		if !matchFileTime(cacheInfo.ModTime(), findCmdModTime.Get()) {
			return nil
		}
	}

	// 如果指定了文件扩展名, 跳过不符合条件的文件
	if findCmdExt.Len() > 0 {
		if _, ok := config.FindExtSliceMap.Load(entryExt); !ok {
			return nil
		}
	}

	// 如果启用了count标志, 则不执行任何操作
	if !findCmdCount.Get() {
		// 如果启用了delete标志, 删除匹配的文件或目录
		if findCmdDelete.Get() {
			if err := deleteMatchedItem(path, entry.IsDir(), config.Cl); err != nil {
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
		if findCmdMove.Get() != "" {
			if err := moveMatchedItem(path, findCmdMove.Get(), config.Cl); err != nil {
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
		if findCmdExec.Get() != "" {
			// 执行命令
			if err := runCommand(findCmdExec.Get(), config.Cl, path); err != nil {
				return fmt.Errorf("执行-exec命令时发生了错误: %v", err)
			}

			return nil
		}
	}

	// 根据标志, 输出完整路径还是匹配到的路径
	if findCmdFullPath.Get() {
		// 增加匹配计数
		config.MatchCount.Add(1)

		// 如果启用了count标志, 则不输出路径
		if !findCmdCount.Get() {
			// 获取完整路径
			fullPath, pathErr := filepath.Abs(path)
			if pathErr != nil {
				fullPath = path // 如果获取完整路径失败, 则使用相对路径
			}

			// 输出完整路径
			printPathColor(fullPath, config.Cl)
		}

		return nil
	}

	// 输出匹配的路径
	config.MatchCount.Add(1) // 增加匹配计数

	// 如果没有启用count标志, 才输出匹配路径
	if !findCmdCount.Get() {
		printPathColor(path, config.Cl)
	}

	return nil
}

// 用于检查find命令的相关参数是否正确
// 参数:
// - findPath: 要查找的路径
func checkFindCmdArgs(findPath string) error {
	// 检查要查找的最大深度是否小于 -1
	if findCmdMaxDepth.Get() < -1 {
		return fmt.Errorf("查找最大深度不能小于 -1")
	}

	// 将限制查找类型转为小写
	if setErr := findCmdType.Set(strings.ToLower(findCmdType.Get())); setErr != nil {
		return fmt.Errorf("转换查找类型失败: %v", setErr)
	}

	// 检查是否为受支持限制查找类型
	if !globals.IsValidFindType(findCmdType.Get()) {
		return fmt.Errorf("无效的类型: %s, 请使用%s", findCmdType.Get(), globals.GetSupportedFindTypes()[:])
	}

	// 如果只显示隐藏文件或目录, 则必须指定 -H 标志
	if !findCmdHidden.Get() && (findCmdType.Get() == globals.FindTypeHidden || findCmdType.Get() == globals.FindTypeHiddenShort) {
		return fmt.Errorf("必须指定 -H 标志才能使用 -type hidden 或 -type h 选项")
	}

	// 检查是否同时指定了-or
	if findCmdOr.Get() {
		// 如果使用-or, 则不能同时使用-and
		if setErr := findCmdAnd.Set("false"); setErr != nil {
			return fmt.Errorf("设置 -and 失败: %v", setErr)
		}
	}

	// 检查如果指定了文件大小, 格式是否正确(格式为 +5M 或 -5M), 单位必须为 B/K/M/G 同时为大写
	if findCmdSize.Get() != "" {
		// 使用正则表达式匹配文件大小条件
		sizeRegex := regexp.MustCompile(`^([+-])(\d+)([BKMGbkmg])$`) // 正确分组：符号、数字、单位
		match := sizeRegex.FindStringSubmatch(findCmdSize.Get())     // 查找匹配项
		if match == nil {
			return fmt.Errorf("文件大小格式错误, 格式如+5M(大于5M)或-5M(小于5M), 支持单位B/K/M/G(大写)")
		}
		_, err := strconv.Atoi(match[2]) // 转换数字部分(match[2])
		if err != nil {
			return fmt.Errorf("文件大小格式错误")
		}
	}

	// 检查如果指定了修改时间, 格式是否正确(格式为 +5 或 -5), 单位必须为 天
	if findCmdModTime.Get() != "" {
		// 使用正则表达式匹配文件时间条件
		timeRegex := regexp.MustCompile(`^([+-])(\d+)$`) // 正确分组：符号、数字
		match := timeRegex.FindStringSubmatch(findCmdModTime.Get())
		if match == nil {
			return fmt.Errorf("文件时间格式错误, 格式如+5(5天前)或-5(5天内)")
		}
		_, err := strconv.Atoi(match[2]) // 转换数字部分(match[2])
		if err != nil {
			return fmt.Errorf("文件时间格式错误")
		}
	}

	// 检查-exec标志是否包含{}
	if findCmdExec.Get() != "" && !strings.Contains(findCmdExec.Get(), "{}") {
		return fmt.Errorf("使用-exec标志时必须包含{}作为路径占位符")
	}

	// 检查-print-cmd标志是否与-exec一起使用
	if findCmdPrintCmd.Get() && findCmdExec.Get() == "" {
		return fmt.Errorf("使用-print-cmd标志时必须同时指定-exec标志")
	}

	// 检查-exec标志是否与-delete或-mv一起使用
	if findCmdExec.Get() != "" && (findCmdDelete.Get() || findCmdMove.Get() != "") {
		return fmt.Errorf("使用-exec标志时不能同时指定-delete或-mv标志")
	}

	// 检查-delete标志是否与-exec或-mv一起使用
	if findCmdDelete.Get() && (findCmdExec.Get() != "" || findCmdMove.Get() != "") {
		return fmt.Errorf("使用-delete标志时不能同时指定-exec或-mv标志")
	}

	// 检查-mv标志是否与-exec或-delete一起使用
	if findCmdMove.Get() != "" && (findCmdExec.Get() != "" || findCmdDelete.Get()) {
		return fmt.Errorf("使用-mv标志时不能同时指定-exec或-delete标志")
	}

	// 检查-mv标志指定的路径是否为文件
	if findCmdMove.Get() != "" {
		if info, err := os.Stat(findCmdMove.Get()); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("-mv 标志指定的路径不存在: %s", findCmdMove.Get())
			}

			return fmt.Errorf("获取文件信息失败: %v", err)

		} else if !info.IsDir() {
			// 如果指定的路径是文件, 则返回错误
			return fmt.Errorf("-mv标志指定的路径必须为目录")
		}
	}

	// 检查软连接最大解析深度是否小于 1
	if findCmdMaxDepthLimit.Get() < 1 {
		return fmt.Errorf("软连接最大解析深度不能小于 1")
	}

	// 检查如果指定了-count则不能同时指定 -exec、-mv、-delete
	if findCmdCount.Get() && (findCmdExec.Get() != "" || findCmdMove.Get() != "" || findCmdDelete.Get()) {
		return fmt.Errorf("使用-count标志时不能同时指定-exec、-mv、-delete标志")
	}

	// 检查如果指定了-count则不能使用 -exec或-mv或-delete
	if findCmdCount.Get() && (findCmdExec.Get() != "" || findCmdMove.Get() != "" || findCmdDelete.Get()) {
		return fmt.Errorf("使用-count标志时不能同时指定-exec、-mv、-delete标志")
	}

	// 检查要查找的路径是否存在
	if _, err := os.Lstat(findPath); err != nil {
		// 检查是否是权限不足的错误
		if os.IsPermission(err) {
			return fmt.Errorf("权限不足, 无法访问某些目录: %s", findPath)
		}

		// 如果是不存在错误, 则返回路径不存在
		if os.IsNotExist(err) {
			return fmt.Errorf("路径不存在: %s", findPath)
		}

		// 其他错误, 返回错误信息
		return fmt.Errorf("检查路径时出错: %s: %v", findPath, err)
	}

	// 检查软连接最大解析深度是否小于 1
	if findCmdMaxDepthLimit.Get() < 1 {
		return fmt.Errorf("软连接最大解析深度不能小于1")
	}

	return nil
}

// matchFileSize 检查文件大小是否符合指定的条件
// 参数:
//
//	fileSize: 文件大小
//	sizeCondition: 文件大小条件
//
// 返回值:
//
//	bool: 如果文件大小符合条件, 返回true, 否则返回false
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

// runCommand 执行单个命令,  支持跨平台并检查shell是否存在
// 参数:
//
//	cmdStr: 要执行的命令字符串
//	cl: 颜色库实例
//	p: 用于替换{}的路径
//
// 返回值:
//
//	error: 如果执行过程中发生错误, 返回错误信息, 否则返回nil
func runCommand(cmdStr string, cl *colorlib.ColorLib, p string) error {
	// 检查cmdStr是否为空
	if cmdStr == "" {
		return fmt.Errorf("cmd is nil")
	}

	// 检查路径是否为空
	if p == "" {
		return fmt.Errorf("path is nil")
	}

	// 检查是否包含{}
	if !strings.Contains(cmdStr, "{}") {
		return fmt.Errorf("使用-exec标志时必须包含{}作为路径占位符")
	}

	// 检查路径是否存在
	if _, err := os.Lstat(p); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("文件/目录不存在: %s", p)
		}

		return fmt.Errorf("无法访问文件/目录: %s", p)
	}

	// 替换{}为实际的文件路径, 并根据系统类型选择引用方式
	if runtime.GOOS == "windows" {
		cmdStr = strings.ReplaceAll(findCmdExec.Get(), "{}", fmt.Sprintf("\"%s\"", p)) // Windows使用双引号
	} else {
		cmdStr = strings.ReplaceAll(findCmdExec.Get(), "{}", fmt.Sprintf("'%s'", p)) // Linux使用单引号
	}

	// 根据操作系统选择shell
	var shell string

	// 定义参数数组
	var args []string

	// 根据系统选择默认shell
	if runtime.GOOS == "windows" {
		// 先尝试使用 PowerShell
		if _, err := exec.LookPath("powershell"); err == nil {
			shell = "powershell"
			args = []string{"-Command", cmdStr}
		} else {
			// 如果 PowerShell 不存在, 使用 cmd
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
	if findCmdPrintCmd.Get() {
		cl.Redf("%s %s\n", shell, strings.Join(args, " "))
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
func deleteMatchedItem(path string, isDir bool, cl *colorlib.ColorLib) error {
	// 检查是否为空路径
	if path == "" {
		return fmt.Errorf("没有可删除的路径")
	}

	// 先检查文件/目录是否存在
	if _, err := os.Lstat(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("文件/目录不存在: %s", path)
		}
		return fmt.Errorf("检查文件/目录时出错: %s: %v", path, err)
	}

	// 删除错误
	var rmErr error

	// 打印删除信息
	if findCmdPrintDelete.Get() {
		cl.Redf("del: %s\n", path)
	}

	// 根据类型选择删除方法
	if isDir {
		// 检查目录是否为空
		dirEntries, readDirErr := os.ReadDir(path)
		if readDirErr != nil {
			return handleError(path, readDirErr)
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
		// 删除文件
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
	if _, err := os.Lstat(path); err != nil {
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
	if err := os.Remove(filepath.Join(filepath.Dir(absTargetPath), ".fck_tmp")); err != nil {
		cl.PrintErrf("delete tmp file error: %v\n", err)
	}

	// 检查目标文件是否已存在
	if _, err := os.Stat(absTargetPath); err == nil {
		// 如果是移动操作, 直接跳过而不是报错
		if findCmdMove.Get() != "" {
			return nil
		}
		return fmt.Errorf("目标文件已存在: %s", absTargetPath)
	}

	// 打印移动信息
	if findCmdPrintMove.Get() {
		cl.Redf("%s -> %s\n", absSearchPath, absTargetPath)
	}

	// 执行移动操作
	if err := os.Rename(absSearchPath, absTargetPath); err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("目标文件已存在: %s -> %s", absSearchPath, absTargetPath)
		}
		if os.IsPermission(err) {
			return fmt.Errorf("权限不足, 无法移动文件: %v", err)
		}

		return fmt.Errorf("移动失败: %s -> %s: %v", absSearchPath, absTargetPath, err)
	}

	return nil
}

// isSymlinkLoop 检查符号链接是否存在循环
// 参数:
//
//	path: 要检查的路径
//
// 返回值:
//
//	bool: true表示存在循环, false表示无循环
func isSymlinkLoop(path string) bool {
	maxDepth := findCmdMaxDepthLimit.Get() // 最大解析深度限制
	visited := make(map[string]bool)       // 已访问路径记录
	currentPath := filepath.Clean(path)    // 清理当前路径

	for depth := 0; depth < maxDepth; depth++ {
		// 检查是否已访问过当前路径
		if visited[currentPath] {
			return true
		}
		visited[currentPath] = true

		// 获取文件信息
		info, err := os.Lstat(currentPath)
		if err != nil || info.Mode()&os.ModeSymlink == 0 {
			return false
		}

		// 解析符号链接
		newPath, err := os.Readlink(currentPath)
		if err != nil {
			return false
		}

		// 处理相对路径
		if !filepath.IsAbs(newPath) {
			newPath = filepath.Join(filepath.Dir(currentPath), newPath)
		}
		currentPath = filepath.Clean(newPath)
	}

	return false // 达到最大深度仍未发现循环
}

// 并发版本的目录遍历函数
func processWalkDirConcurrent(config *globals.FindConfig, findPath string) error {
	var wg sync.WaitGroup
	pathChan := make(chan string, runtime.NumCPU()*2000) // 通道缓冲区
	errorChan := make(chan error, 100)                   // 错误通道

	// 设置最大并发数量为 CPU 核心数的两倍
	maxWorkers := runtime.NumCPU() * 2

	// 启动多个 worker goroutine 处理路径
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer func() {
				// 捕获 panic
				if r := recover(); r != nil {
					errorChan <- fmt.Errorf("worker goroutine 发生 panic: %v", r)
				}

				wg.Done()
			}()

			// 处理路径
			for path := range pathChan {
				// 检查是否存在
				entry, lstatErr := os.Lstat(path)
				if lstatErr != nil {
					// 忽略不存在的路径错误
					if os.IsNotExist(lstatErr) {
						continue
					}
					if len(errorChan) < cap(errorChan) {
						errorChan <- fmt.Errorf("无法访问 %s: %v", path, lstatErr)
					}
					continue
				}

				// 构建 DirEntryWrapper
				dirEntry := &globals.DirEntryWrapper{
					NameVal:  entry.Name(),
					IsDirVal: entry.IsDir(),
					ModeVal:  entry.Mode(),
				}

				// 通过 processFindCmd 处理路径
				processErr := processFindCmd(config, dirEntry, path)
				if processErr != nil {
					if processErr == filepath.SkipDir {
						// 如果是 SkipDir, 跳过该目录即可
						continue
					}

					if errors.Is(processErr, os.ErrPermission) {
						config.Cl.PrintErrf("路径 %s 权限不足, 已跳过\n", path)
						// 跳过该路径
						continue
					}

					if len(errorChan) < cap(errorChan) {
						errorChan <- fmt.Errorf("处理路径失败 %s: %v", path, processErr)
					}
				}
			}
		}()
	}

	// 主逻辑通过 WalkDir 将路径发送到 channel
	walkDirErr := filepath.WalkDir(findPath, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			// 忽略不存在的路径错误
			if os.IsNotExist(err) {
				return nil
			}

			// 忽略权限错误
			if os.IsPermission(err) {
				config.Cl.PrintErrf("路径 %s 因权限不足已跳过\n", p)
				return nil
			}

			return fmt.Errorf("访问 %s 时出错: %s", p, err)
		}

		// 跳过findPath本身
		if p == findPath {
			return nil
		}

		// 检查当前路径的深度是否超过最大深度(先将路径统一转为/分隔符)
		depth := strings.Count(p[len(findPath):], string(filepath.Separator))
		if findCmdMaxDepth.Get() >= 0 && depth > findCmdMaxDepth.Get() {
			return filepath.SkipDir
		}

		// 默认隐藏文件或隐藏目录不参与匹配
		// 如果没有启用隐藏标志且是隐藏目录或文件, 则跳过
		if !findCmdHidden.Get() {
			if isHidden(p) {
				// 如果是隐藏目录, 跳过整个目录
				if d.IsDir() {
					return filepath.SkipDir
				}

				// 如果是隐藏文件, 跳过单个文件
				return nil
			}
		}

		// 检查是否为符号链接循环, 然后做相应处理
		if d.Type()&os.ModeSymlink != 0 {
			if isSymlinkLoop(p) {
				return filepath.SkipDir
			}
		}

		// 发送路径到 channel
		pathChan <- p

		return nil
	})

	// 遍历完成, 关闭路径通道
	close(pathChan)

	// 等待所有 worker 完成
	wg.Wait()

	// 关闭错误通道
	close(errorChan)

	// 优先检查遍历过程中是否遇到错误
	if walkDirErr != nil {
		if os.IsPermission(walkDirErr) {
			return fmt.Errorf("权限不足, 无法访问某些目录:  %v", walkDirErr)
		} else if os.IsNotExist(walkDirErr) {
			return fmt.Errorf("路径不存在: %v", walkDirErr)
		}
		return fmt.Errorf("遍历目录时出错: %v", walkDirErr)
	}

	// 收集并分类错误（限制最多显示5个不同错误）
	errorMap := make(map[string]error)
	errorCount := 0
	for err := range errorChan {
		if err != nil && errorCount < 5 {
			if _, exists := errorMap[err.Error()]; !exists {
				errorMap[err.Error()] = err
				errorCount++
			}
		}
	}

	// 合并错误并添加统计信息
	if len(errorMap) > 0 {
		var combinedErr error
		var errorList []string
		for _, err := range errorMap {
			combinedErr = errors.Join(combinedErr, err)
			errorList = append(errorList, fmt.Sprintf("\t- %s", err))
		}
		return fmt.Errorf("共发现%d类错误:\n%s", len(errorMap), strings.Join(errorList, "\n"))
	}

	return nil
}

// matchPattern 通用匹配函数, 检查输入字符串是否匹配模式
// 参数:
//
//	input: 输入字符串(文件名或路径)
//	pattern: 匹配模式字符串
//	regex: 编译好的正则表达式(如果启用正则)
//	config: 配置参数
//
// 返回值:
//
//	bool: 是否匹配
func matchPattern(input, pattern string, regex *regexp.Regexp, config *globals.FindConfig) bool {
	// 如果模式为空, 则不匹配
	if pattern == "" {
		return false
	}

	// 如果启用正则匹配, 使用正则表达式匹配
	if config.IsRegex {
		if regex == nil {
			return false
		}
		return regex.MatchString(input)
	}

	// 根据大小写敏感性处理字符串
	var s, p string
	if config.CaseSensitive {
		// 区分大小写
		s = input
		p = pattern
	} else {
		// 默认不区分大小写
		s = strings.ToLower(input)
		p = strings.ToLower(pattern)
	}

	// 全字匹配处理
	if config.WholeWord {
		return s == p
	}

	// 匹配模式处理
	return strings.Contains(s, p)
}
