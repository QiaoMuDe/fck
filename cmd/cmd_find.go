package cmd

import (
	"errors"
	"flag"
	"fmt"
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
func findCmdMain(cl *colorlib.ColorLib, cmd *flag.FlagSet) error {
	// 获取第一个参数作为查找路径
	findPath := cmd.Arg(0)

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

	// 根据*findCmdColor选项设置颜色
	if *findCmdColor {
		cl.NoColor.Store(false)
	} else {
		cl.NoColor.Store(true)
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

	// 新增完整关键字匹配逻辑
	if *findCmdWholeWord {
		// 名字在关键字前后添加开头和结尾匹配
		if escapedName != "" {
			escapedName = "^" + escapedName + "$"
		}

		// 排除的名字在关键字前后添加开头和结尾匹配
		if exCludedName != "" {
			exCludedName = "^" + exCludedName + "$"
		}

		// 路径在关键字前后添加开头和结尾匹配
		if escapedPath != "" {
			escapedPath = "^" + escapedPath + "$"
		}

		// 排除的路径在关键字前后添加开头和结尾匹配
		if exCludedPath != "" {
			exCludedPath = "^" + exCludedPath + "$"
		}
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

	// 定义一个统计匹配项的
	matchCount := atomic.Int64{}
	matchCount.Store(0)

	// 检查是否启用并发模式
	if *findCmdX {
		// 启用并发模式
		if err := processWalkDirConcurrent(cl, nameRegex, exNameRegex, pathRegex, exPathRegex, findPath, &matchCount); err != nil {
			return err
		}
	} else {
		// 运行processWalkDir单线程遍历模式
		if walkDirErr := processWalkDir(cl, nameRegex, exNameRegex, pathRegex, exPathRegex, findPath, &matchCount); walkDirErr != nil {
			return walkDirErr
		}
	}

	// 如果启用了count标志, 只输出匹配数量
	if *findCmdCount {
		fmt.Println(matchCount.Load())
		return nil
	}

	return nil
}

// processWalkDir 用于处理 filepath.WalkDir 的逻辑
// 参数:
// - cl: 颜色库
// - nameRegex: 文件名正则表达式
// - exNameRegex: 排除的文件名正则表达式
// - pathRegex: 路径正则表达式
// - exPathRegex: 排除的路径正则表达式
// - findPath: 要查找的路径
// 返回值:
// - error: 错误信息
func processWalkDir(cl *colorlib.ColorLib, nameRegex, exNameRegex, pathRegex, exPathRegex *regexp.Regexp, findPath string, matchCount *atomic.Int64) error {
	// 使用 filepath.WalkDir 遍历目录
	walkDirErr := filepath.WalkDir(findPath, func(path string, entry os.DirEntry, err error) error {
		// 检查遍历过程中是否遇到错误
		if err != nil {
			// 忽略不存在的报错
			if os.IsNotExist(err) {
				// 路径不存在，跳过
				return nil
			}

			// 检查是否为权限不足的报错
			if os.IsPermission(err) {
				cl.PrintErrf("权限不足，无法访问某些目录: %s\n", path)
				return nil
			}

			return fmt.Errorf("访问时出错：%s", err)
		}

		// 跳过*findCmdPath目录本身
		if path == findPath {
			return nil
		}

		// 检查当前路径的深度是否超过最大深度(先将路径统一转为/分隔符)
		depth := strings.Count(path[len(findPath):], string(filepath.Separator))
		if *findCmdMaxDepth >= 0 && depth > *findCmdMaxDepth {
			return filepath.SkipDir
		}

		// 检查是否为符号链接循环, 然后做相应处理
		if entry.Type()&os.ModeSymlink != 0 {
			if isSymlinkLoop(path) {
				return filepath.SkipDir
			}
		}

		// 处理文件或目录
		if processErr := processFindCmd(cl, nameRegex, exNameRegex, pathRegex, exPathRegex, entry, path, matchCount); processErr != nil {
			return processErr
		}

		return nil
	})

	// 检查遍历过程中是否遇到错误
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

// processFindCmd 用于处理 find 命令的逻辑
// 参数:
// - cl: 颜色库
// - nameRegex: 文件名正则表达式
// - exNameRegex: 排除的文件名正则表达式
// - pathRegex: 路径正则表达式
// - exPathRegex: 排除的路径正则表达式
// - entry: 文件或目录的DirEntry对象
// - path: 文件或目录的完整路径
// - matchCount: 匹配数量
// 返回值:
// - error: 错误信息
func processFindCmd(cl *colorlib.ColorLib, nameRegex, exNameRegex, pathRegex, exPathRegex *regexp.Regexp, entry os.DirEntry, path string, matchCount *atomic.Int64) error {

	// 如果指定了-n和-p参数, 并且指定了-or参数, 则只检查文件名或路径是否匹配(默认为或操作)
	if *findCmdOr && *findCmdName != "" && *findCmdPath != "" {
		// 执行或操作
		if nameRegex.MatchString(entry.Name()) || pathRegex.MatchString(path) {
			if err := filterConditions(entry, path, cl, exNameRegex, exPathRegex, matchCount); err != nil {
				return err
			}
		}
		return nil
	}

	// 如果指定了-n和-p参数, 则同时检查文件名和路径是否匹配(默认为-and操作)
	if *findCmdAnd && *findCmdName != "" && *findCmdPath != "" {
		if nameRegex.MatchString(entry.Name()) && pathRegex.MatchString(path) {
			// 如果同时匹配, 则执行筛选条件
			if err := filterConditions(entry, path, cl, exNameRegex, exPathRegex, matchCount); err != nil {
				return err
			}
		}

		return nil
	}

	// 如果指定了-n参数, 则检查文件名是否匹配
	if *findCmdName != "" {
		if nameRegex.MatchString(entry.Name()) {
			if err := filterConditions(entry, path, cl, exNameRegex, exPathRegex, matchCount); err != nil {
				return err
			}
		}
		return nil
	}

	// 如果指定了-p参数, 则检查路径是否匹配
	if *findCmdPath != "" {
		if pathRegex.MatchString(path) {
			if err := filterConditions(entry, path, cl, exNameRegex, exPathRegex, matchCount); err != nil {
				return err
			}
		}
		return nil
	}

	// 默认情况下
	if err := filterConditions(entry, path, cl, exNameRegex, exPathRegex, matchCount); err != nil {
		return err
	}

	return nil
}

// filterConditions 用于在循环中筛选条件的函数
// 参数:
// - entry: 文件或目录的DirEntry对象
// - path: 文件或目录的完整路径
// - cl: 颜色库
// - exNameRegex: 排除的文件名正则表达式
// - exPathRegex: 排除的路径正则表达式
// - matchCount: 匹配数量
// 返回值:
// - error: 错误信息
func filterConditions(entry os.DirEntry, path string, cl *colorlib.ColorLib, exNameRegex, exPathRegex *regexp.Regexp, matchCount *atomic.Int64) error {
	// 默认隐藏文件或隐藏目录不参与匹配
	// 如果没有启用隐藏标志且是隐藏目录或文件, 则跳过
	if !*findCmdHidden {
		if isHidden(path) {
			// 如果是隐藏目录，跳过整个目录
			if entry.IsDir() {
				return filepath.SkipDir
			}

			// 如果是隐藏文件，跳过单个文件
			return nil
		}
	}

	// 如果指定了排除文件或目录名, 跳过匹配的文件或目录
	if *findCmdExcludeName != "" {
		if exNameRegex.MatchString(entry.Name()) {
			return nil
		}
	}

	// 如果指定了排除路径, 跳过匹配的路径
	if *findCmdExcludePath != "" {
		if exPathRegex.MatchString(path) {
			return nil
		}
	}

	// 根据findCmdType的值，跳过不符合要求的文件或目录
	switch *findCmdType {
	case globals.FindTypeFile, globals.FindTypeFileShort:
		// 如果只查找文件，跳过目录
		if entry.IsDir() {
			return nil
		}
	case globals.FindTypeDir, globals.FindTypeDirShort:
		// 如果只查找目录，跳过文件
		if !entry.IsDir() {
			return nil
		}
	case globals.FindTypeSymlink, globals.FindTypeSymlinkShort:
		// 如果只查找软链接，跳过非软链接
		// Windows系统特殊处理
		if runtime.GOOS == "windows" {
			// 如果不是不是.lnk文件，跳过
			if !strings.HasSuffix(entry.Name(), ".lnk") || !strings.HasSuffix(entry.Name(), ".url") {
				return nil
			}
		}

		// 统一检查符号链接标志
		if entry.Type()&os.ModeSymlink == 0 {
			return nil // 不是符号链接, 跳过
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
			fileInfo, sizeErr := entry.Info()
			if sizeErr != nil || fileInfo.Size() > 0 {
				return nil
			}
		}
	case globals.FindTypeExecutable, globals.FindTypeExecutableShort:
		// 如果只查找可执行文件或目录
		if runtime.GOOS == "windows" {
			// Windows系统检查文件扩展名
			ext := strings.ToLower(filepath.Ext(entry.Name()))
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
		// Windows不支持Unix domain socket，直接返回nil
		if runtime.GOOS == "windows" {
			return nil
		}

		// 如果只查找socket文件, 跳过非socket文件
		if entry.Type()&os.ModeSocket == 0 {
			return nil
		}
	case globals.FindTypePipe, globals.FindTypePipeShort:
		// Windows不支持Unix命名管道，直接返回nil
		if runtime.GOOS == "windows" {
			return nil
		}

		// 如果只查找管道文件, 跳过非管道文件
		if entry.Type()&os.ModeNamedPipe == 0 {
			return nil
		}
	case globals.FindTypeBlock, globals.FindTypeBlockShort:
		// Windows不支持Unix块设备文件，直接返回nil
		if runtime.GOOS == "windows" {
			return nil
		}

		// 如果只查找块设备文件, 跳过非块设备文件
		if entry.Type()&os.ModeDevice == 0 {
			return nil
		}
	case globals.FindTypeChar, globals.FindTypeCharShort:
		// Windows不支持Unix字符设备文件，直接返回nil
		if runtime.GOOS == "windows" {
			return nil
		}

		// 如果只查找字符设备文件, 跳过非字符设备文件
		if entry.Type()&os.ModeCharDevice == 0 {
			return nil
		}
	case globals.FindTypeAppend, globals.FindTypeAppendShort:
		// Windows文件系统不支持Unix追加模式标志，直接返回nil
		if runtime.GOOS == "windows" {
			return nil
		}

		// 如果只查找追加模式文件, 跳过非追加模式文件
		if entry.Type()&os.ModeAppend == 0 {
			return nil
		}
	case globals.FindTypeNonAppend, globals.FindTypeNonAppendShort:
		// Windows文件系统不支持Unix追加模式标志，直接返回nil
		if runtime.GOOS == "windows" {
			return nil
		}

		// 如果只查找非追加模式文件, 跳过追加模式文件
		if entry.Type()&os.ModeAppend != 0 {
			return nil
		}
	case globals.FindTypeExclusive, globals.FindTypeExclusiveShort:
		// Windows文件系统不支持Unix独占模式标志，直接返回nil
		if runtime.GOOS == "windows" {
			return nil
		}

		// 如果只查找独占模式文件, 跳过非独占模式文件
		if entry.Type()&os.ModeExclusive == 0 {
			return nil
		}
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

	// 扩展名检查
	if *findCmdExt != "" {
		ext := filepath.Ext(entry.Name())       // 获取文件扩展名
		exts := strings.Split(*findCmdExt, ",") // 分割多个扩展名
		match := false                          // 标记是否匹配
		for _, e := range exts {                // 遍历所有扩展名
			if ext == strings.TrimSpace(e) { // 检查当前扩展名是否匹配
				match = true // 匹配成功
				break
			}
		}
		// 如果没有匹配到任何扩展名, 跳过
		if !match {
			return nil
		}
	}

	// 如果启用了count标志, 则不执行任何操作
	if !*findCmdCount {
		// 如果启用了delete标志, 删除匹配的文件或目录
		if *findCmdDelete {
			if err := deleteMatchedItem(path, entry.IsDir(), cl); err != nil {
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
			// 执行命令
			if err := runCommand(*findCmdExec, cl, path); err != nil {
				return fmt.Errorf("执行-exec命令时发生了错误: %v", err)
			}

			return nil
		}
	}

	// 根据标志, 输出完整路径还是匹配到的路径
	if *findCmdFullPath {
		// 增加匹配计数
		matchCount.Add(1)

		// 如果启用了count标志, 则不输出路径
		if !*findCmdCount {
			// 获取完整路径
			fullPath, pathErr := filepath.Abs(path)
			if pathErr != nil {
				fullPath = path // 如果获取完整路径失败, 则使用相对路径
			}

			// 输出完整路径
			printPathColor(fullPath, cl)
		}

		return nil
	}

	// 输出匹配的路径
	matchCount.Add(1) // 增加匹配计数

	// 如果没有启用count标志, 才输出匹配路径
	if !*findCmdCount {
		printPathColor(path, cl)
	}

	return nil
}

// 用于检查find命令的相关参数是否正确
// 参数:
// - findPath: 要查找的路径
func checkFindCmdArgs(findPath string) error {
	// 检查要查找的最大深度是否小于 -1
	if *findCmdMaxDepth < -1 {
		return fmt.Errorf("查找最大深度不能小于 -1")
	}

	// 将限制查找类型转为小写
	*findCmdType = strings.ToLower(*findCmdType)

	// 检查是否为受支持限制查找类型
	if !globals.IsValidFindType(*findCmdType) {
		return fmt.Errorf("无效的类型: %s, 请使用%s", *findCmdType, globals.GetSupportedFindTypes()[:])
	}

	// 如果只显示隐藏文件或目录, 则必须指定 -H 标志
	if !*findCmdHidden && (*findCmdType == globals.FindTypeHidden || *findCmdType == globals.FindTypeHiddenShort) {
		return fmt.Errorf("必须指定 -H 标志才能使用 -type hidden 或 -type h 选项")
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

	// 检查-exec标志是否与-delete或-mv一起使用
	if *findCmdExec != "" && (*findCmdDelete || *findCmdMove != "") {
		return fmt.Errorf("使用-exec标志时不能同时指定-delete或-mv标志")
	}

	// 检查-delete标志是否与-exec或-mv一起使用
	if *findCmdDelete && (*findCmdExec != "" || *findCmdMove != "") {
		return fmt.Errorf("使用-delete标志时不能同时指定-exec或-mv标志")
	}

	// 检查-mv标志是否与-exec或-delete一起使用
	if *findCmdMove != "" && (*findCmdExec != "" || *findCmdDelete) {
		return fmt.Errorf("使用-mv标志时不能同时指定-exec或-delete标志")
	}

	// 检查-mv标志指定的路径是否为文件
	if *findCmdMove != "" {
		if info, err := os.Stat(*findCmdMove); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("-mv 标志指定的路径不存在: %s", *findCmdMove)
			}

			return fmt.Errorf("获取文件信息失败: %v", err)

		} else if !info.IsDir() {
			// 如果指定的路径是文件, 则返回错误
			return fmt.Errorf("-mv标志指定的路径必须为目录")
		}
	}

	// 检查软连接最大解析深度是否小于 1
	if *findCmdMaxDepthLimit < 1 {
		return fmt.Errorf("软连接最大解析深度不能小于 1")
	}

	// 检查如果指定了-count则不能同时指定 -exec、-mv、-delete
	if *findCmdCount && (*findCmdExec != "" || *findCmdMove != "" || *findCmdDelete) {
		return fmt.Errorf("使用-count标志时不能同时指定-exec、-mv、-delete标志")
	}

	// 检查如果指定了-count则不能使用 -exec或-mv或-delete
	if *findCmdCount && (*findCmdExec != "" || *findCmdMove != "" || *findCmdDelete) {
		return fmt.Errorf("使用-count标志时不能同时指定-exec、-mv、-delete标志")
	}

	// 检查要查找的路径是否存在
	if _, err := os.Lstat(findPath); err != nil {
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

	// 添加扩展名参数验证
	if *findCmdExt != "" {
		exts := strings.Split(*findCmdExt, ",")
		for _, ext := range exts {
			ext = strings.TrimSpace(ext)
			if strings.ContainsAny(ext, "\\/:") {
				return fmt.Errorf("扩展名包含非法字符: %s", ext)
			}
			if !strings.HasPrefix(ext, ".") {
				return fmt.Errorf("扩展名应以点开头: %s (建议使用 .%s)", ext, ext)
			}
		}
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

// runCommand 执行单个命令， 支持跨平台并检查shell是否存在
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

	// 替换{}为实际的文件路径，并根据系统类型选择引用方式
	if runtime.GOOS == "windows" {
		cmdStr = strings.Replace(*findCmdExec, "{}", fmt.Sprintf("\"%s\"", p), -1) // Windows使用双引号
	} else {
		cmdStr = strings.Replace(*findCmdExec, "{}", fmt.Sprintf("'%s'", p), -1) // Linux使用单引号
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

	var rmErr error

	// 打印删除信息
	if *findCmdPrintDelete {
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
		cl.Redf("%s -> %s\n", absSearchPath, absTargetPath)
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

// isSymlinkLoop 检查符号链接是否存在循环
// 参数:
//
//	path: 要检查的路径
//
// 返回值:
//
//	bool: true表示存在循环，false表示无循环
func isSymlinkLoop(path string) bool {
	maxDepth := *findCmdMaxDepthLimit // 最大解析深度限制
	visited := make(map[string]bool)
	currentPath := filepath.Clean(path)

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
func processWalkDirConcurrent(cl *colorlib.ColorLib, nameRegex, exNameRegex, pathRegex, exPathRegex *regexp.Regexp, findPath string, matchCount *atomic.Int64) error {
	var wg sync.WaitGroup
	pathChan := make(chan string, 30000) // 通道缓冲区
	errorChan := make(chan error, 100)   // 错误通道

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
				processErr := processFindCmd(cl, nameRegex, exNameRegex, pathRegex, exPathRegex, dirEntry, path, matchCount)
				if processErr != nil {
					if processErr == filepath.SkipDir {
						// 如果是 SkipDir，跳过该目录即可
						continue
					}

					if errors.Is(processErr, os.ErrPermission) {
						cl.PrintErrf("路径 %s 权限不足，已跳过\n", path)
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
				cl.PrintErrf("路径 %s 因权限不足已跳过\n", p)
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
		if *findCmdMaxDepth >= 0 && depth > *findCmdMaxDepth {
			return filepath.SkipDir
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

	// 遍历完成，关闭路径通道
	close(pathChan)

	// 等待所有 worker 完成
	wg.Wait()

	// 关闭错误通道
	close(errorChan)

	// 优先检查遍历过程中是否遇到错误
	if walkDirErr != nil {
		if os.IsPermission(walkDirErr) {
			return fmt.Errorf("权限不足，无法访问某些目录:  %v", walkDirErr)
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
