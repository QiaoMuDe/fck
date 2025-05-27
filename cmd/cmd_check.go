package cmd

import (
	"bufio"
	"fmt"
	"hash"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/globals"

	"golang.org/x/sync/errgroup"
)

// diffCmdMain 是 check 命令的主函数
func diffCmdMain(cl *colorlib.ColorLib) error {
	// 启动一个 goroutine，在用户按下 Ctrl+C 时取消操作
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		// 使用自定义错误作为取消原因
		cl.PrintWarn("用户中断操作")
		os.Exit(1) // 退出程序
	}()

	// 检查三个参数是否都为空
	if *diffCmdFile == "" && *diffCmdDirA == "" && *diffCmdDirB == "" {
		return fmt.Errorf("必须指定一个校验文件或两个目录。或 -h 参数查看帮助信息。")
	}

	// (1) 执行根据单校验文件校验目录完整性的逻辑 -f 参数
	if *diffCmdFile != "" && *diffCmdDirs == "" {
		if err := fileCheck(*diffCmdFile, cl); err != nil {
			return err
		}
		return nil
	}

	// (2) 如果指定校验文件不为空，同时也通过*diffCmdDirs指定了目录 -f 参数和 -d 参数
	if *diffCmdFile != "" && *diffCmdDirs != "" {
		// 执行校验文件和目录的逻辑
		if err := checkWithFileAndDir(*diffCmdFile, *diffCmdDirs, cl); err != nil {
			return err
		}
		return nil
	}

	// (3) 执行对比目录A和目录B的逻辑 -a 参数和 -b 参数
	if *diffCmdDirA == "" || *diffCmdDirB == "" {
		return fmt.Errorf("对比目录时，必须同时指定 -a 和 -b 参数。或 -h 参数查看帮助信息。")
	}
	if *diffCmdDirA != "" && *diffCmdDirB != "" {
		if err := checkWithDirAndDir(cl); err != nil {
			return err
		}
		return nil
	}

	return nil
}

// 对比dirA和dirB两个目录的文件内容是否一致
func checkWithDirAndDir(cl *colorlib.ColorLib) error {
	// 检查目录A 和 目录B 是否存在
	if _, err := os.Stat(*diffCmdDirA); err != nil {
		return fmt.Errorf("目录A不存在: %s", *diffCmdDirA)
	}
	if _, err := os.Stat(*diffCmdDirB); err != nil {
		return fmt.Errorf("目录B不存在: %s", *diffCmdDirB)
	}

	// 检查目录A是否为绝对路径，如果不是，则转换为绝对路径
	if !filepath.IsAbs(*diffCmdDirA) {
		absDirA, err := filepath.Abs(*diffCmdDirA)
		if err != nil {
			return fmt.Errorf("无法获取目录A的绝对路径: %v", err)
		}
		*diffCmdDirA = absDirA
	}

	// 检查目录B是否为绝对路径，如果不是，则转换为绝对路径
	if !filepath.IsAbs(*diffCmdDirB) {
		absDirB, err := filepath.Abs(*diffCmdDirB)
		if err != nil {
			return fmt.Errorf("无法获取目录B的绝对路径: %v", err)
		}
		*diffCmdDirB = absDirB
	}

	// 校验目录A 和 目录B //
	// 检查指定的哈希算法是否有效
	hashType, ok := globals.SupportedAlgorithms[*diffCmdType]
	if !ok {
		return fmt.Errorf("在校验哈希值时，哈希算法 %s 无效", *diffCmdType)
	}

	// 获取两个目录下的文件列表
	filesA, filesB, err := getFilesFromDirs(*diffCmdDirA, *diffCmdDirB)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	// 检查是否需要写入文件
	var fileWrite *os.File
	if *diffCmdWrite {
		var err error
		// 打开文件以写入
		fileWrite, err = os.OpenFile(globals.OutputCheckFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("打开文件 %s 失败: %v", globals.OutputCheckFileName, err)
		}
		defer fileWrite.Close()

		// 写入文件头
		if err := writeFileHeader(fileWrite, *hashCmdType, globals.TimestampFormat); err != nil {
			return fmt.Errorf("写入文件头失败: %v", err)
		}
	}

	// 比较两个目录的文件
	if err := compareFiles(filesA, filesB, hashType, cl, fileWrite); err != nil {
		return fmt.Errorf("比较文件失败: %v", err)
	}

	// 如果是写入文件模式，则打印文件路径
	if *diffCmdWrite {
		cl.PrintOkf("比较结果已写入文件: %s", globals.OutputCheckFileName)
	}

	return nil
}

// readHashFileToMap 读取校验文件并加载到 map 中
func readHashFileToMap(checkFile string, cl *colorlib.ColorLib, isRelPath bool) (globals.VirtualHashMap, func() hash.Hash, error) {
	// 创建一个新的映射，用于存储替换后的路径
	replaceMap := make(globals.VirtualHashMap)

	// 检查校验文件是否存在
	if _, err := os.Stat(checkFile); err != nil {
		return nil, nil, fmt.Errorf("校验文件不存在: %s", checkFile)
	}

	// 打开校验文件
	checkFileRead, openErr := os.OpenFile(checkFile, os.O_RDONLY, 0644)
	if openErr != nil {
		return nil, nil, fmt.Errorf("无法打开校验文件: %v", openErr)
	}
	defer checkFileRead.Close()

	// 解析校验文件内容
	scanner := bufio.NewScanner(checkFileRead)

	// 解析文件头
	hashFunc, err := parseHeader(scanner)
	if err != nil {
		return nil, nil, err
	}

	// 定义一个计数器
	var lineCount int

	// 逐行读取校验文件
	for scanner.Scan() { // 逐行读取
		line := scanner.Text() // 获取当前行的文本

		lineCount++ // 计数器加 1

		// 如果当前行是空行或以#开头，则跳过
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析校验文件中的哈希值和文件路径
		parts := strings.Fields(line) // 按空格分割

		// 获取哈希值
		expectedHash := parts[0]

		// 获取文件路径, 从第二个元素开始到结尾
		filePath := strings.Join(parts[1:], " ")

		// 去除路径中的引号
		filePath = strings.Trim(filePath, `"`)

		// 将路径中的双\\替换为单\
		filePath = strings.ReplaceAll(filePath, `\\`, `\`)

		// 根据参数判断是否使用相对路径
		if isRelPath {
			// 手动解析路径，找到根目录部分
			rootDir := strings.Split(filePath, string(filepath.Separator))[0]

			// 获取相对路径
			relPath, err := filepath.Rel(rootDir, filePath)
			if err != nil {
				return nil, nil, fmt.Errorf("获取相对路径时出错: %v", err)
			}

			// 如果哈希值或文件路径为空，则跳过
			if expectedHash == "" || relPath == "" {
				cl.PrintErrf("error: 校验文件格式错误, 文件 %s 的第 %d 行, %s", checkFile, lineCount, line)
				continue
			}

			// 构建虚拟路径
			virtualPath := filepath.Join(globals.VirtualRootDir, relPath)

			// 存储到 map 中
			replaceMap[virtualPath] = globals.VirtualHashEntry{
				RealPath: filePath,     // 真实路径
				Hash:     expectedHash, // 哈希值
			}
		} else {
			// 如果哈希值或文件路径为空，则跳过
			if expectedHash == "" || filePath == "" {
				cl.PrintErrf("error: 校验文件格式错误, 文件 %s 的第 %d 行, %s", checkFile, lineCount, line)
				continue
			}

			// 存储到 map 中
			replaceMap[filePath] = globals.VirtualHashEntry{
				RealPath: filePath,     // 路径
				Hash:     expectedHash, // 哈希值
			}
		}
	}

	// 检查是否有错误发生
	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("读取校验文件时出错: %v", err)
	}

	// 检查map是否为空
	if len(replaceMap) == 0 {
		return nil, nil, fmt.Errorf("没有找到有效的校验文件内容")
	}

	return replaceMap, hashFunc, nil
}

// parseHeader 解析校验文件的文件头，提取哈希算法
func parseHeader(scanner *bufio.Scanner) (func() hash.Hash, error) {
	// 读取文件头
	if !scanner.Scan() {
		return nil, fmt.Errorf("校验文件为空")
	}
	headerLine := scanner.Text() // 获取第一行文件头

	// 使用正则表达式解析文件头
	headerRegex := regexp.MustCompile(`^#(\w+)#(.+)$`)
	matches := headerRegex.FindStringSubmatch(headerLine)
	if matches == nil {
		return nil, fmt.Errorf("校验文件头格式错误, 格式应为 #hashType#timestamp")
	}

	hashType := matches[1] // 哈希算法
	if hashType == "" {
		return nil, fmt.Errorf("校验文件头格式错误, 必须指定哈希算法")
	}

	// 检查哈希算法是否支持
	hashFunc, ok := globals.SupportedAlgorithms[hashType]
	if !ok {
		return nil, fmt.Errorf("不支持的哈希算法: %s", hashType)
	}

	return hashFunc, nil
}

// fileCheck 根据单校验文件校验目录完整性的逻辑 -f 参数
func fileCheck(checkFile string, cl *colorlib.ColorLib) error {
	// 检查校验文件是否为空
	if *diffCmdFile == "" {
		return fmt.Errorf("在校验文件时，必须指定一个校验文件 checksum.hash")
	}

	// 读取校验文件并加载到 map 中
	checkFileHashes, hashFunc, err := readHashFileToMap(checkFile, cl, false)
	if err != nil {
		return err
	}

	// 存储目标目录文件的 map
	targetDirHashes := make(map[string]string)

	// 计算目标目录文件哈希值
	for _, entry := range checkFileHashes {
		// 检查文件是否存在，若不存在则输出更详细的提示信息并跳过当前文件处理
		if _, err := os.Stat(entry.RealPath); err != nil {
			cl.PrintWarnf("在进行校验时，发现文件 %s 不存在，将跳过该文件的校验", entry.RealPath)
			continue
		}

		// 存储哈希值
		var hash string

		// 存储错误信息
		var checksumErr error

		// 计算哈希值
		hash, checksumErr = checksum(entry.RealPath, hashFunc)
		if checksumErr != nil {
			cl.PrintErrf("计算文件哈希失败: %v", checksumErr)
		}

		// 将哈希值存储在目标目录哈希值映射中
		targetDirHashes[entry.RealPath] = hash
	}

	// 对比哈希值
	var checkCount int
	for filePath, checkEntry := range checkFileHashes {
		// 获取实际的哈希值
		targetHash, ok := targetDirHashes[filePath]
		if !ok {
			// 如果实际的哈希值不存在，则跳过
			continue
		}

		// 比较哈希值
		if targetHash != checkEntry.Hash {
			cl.PrintErrf("文件 %s 不一致, 预期Hash值: %s, 实际Hash值: %s", filePath, getLast8Chars(checkEntry.Hash), getLast8Chars(targetHash))
			checkCount++
		}
	}

	// 检查 checkCount 是否为 0
	if checkCount == 0 {
		cl.PrintOk("校验成功，无文件差异")
	}

	return nil
}

// getFilesFromDirs 获取两个目录下的文件列表
func getFilesFromDirs(dirA, dirB string) (map[string]string, map[string]string, error) {
	var eg errgroup.Group

	// 获取目录 A 的文件列表
	var filesA map[string]string
	var getFilesAErr error
	eg.Go(func() error {
		filesA, getFilesAErr = getFiles(dirA)
		return getFilesAErr
	})

	// 获取目录 B 的文件列表
	var filesB map[string]string
	var getFilesBErr error
	eg.Go(func() error {
		filesB, getFilesBErr = getFiles(dirB)
		return getFilesBErr
	})

	// 等待两个目录的文件列表获取完成
	if err := eg.Wait(); err != nil {
		return nil, nil, fmt.Errorf("获取文件列表时出错: %v", err)
	}

	// 如果获取文件列表时出错，则返回错误
	if getFilesAErr != nil {
		return nil, nil, fmt.Errorf("读取目录 A 时出错: %v", getFilesAErr)
	}
	if getFilesBErr != nil {
		return nil, nil, fmt.Errorf("读取目录 B 时出错: %v", getFilesBErr)
	}

	return filesA, filesB, nil
}

// getFiles 遍历指定目录，返回目录下所有文件的名称到路径的映射
//
// 参数：
//   - dir: 要遍历的目录路径
//
// 返回值：
//   - map[string]string: 返回一个映射, key 为相对于 dir 的文件路径, value 为文件的完整路径
func getFiles(dir string) (map[string]string, error) {
	files := make(map[string]string)
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// 如果是文件，则加入到 map 中
		if !d.IsDir() {
			relPath, err := filepath.Rel(dir, path)
			if err != nil {
				return fmt.Errorf("获取相对路径失败: %v", err)
			}
			if relPath == "" {
				return fmt.Errorf("相对路径为空")
			}

			files[relPath] = path // 使用相对路径作为键
		}
		return nil
	})
	return files, err
}

// compareFiles 比较两个目录的文件
func compareFiles(filesA, filesB map[string]string, hashType func() hash.Hash, cl *colorlib.ColorLib, fileWrite *os.File) error {
	// 初始化统计计数器
	sameCount := 0
	diffCount := 0
	onlyACount := 0
	onlyBCount := 0

	// 比较相同路径的文件
	matchFileNameFiles := "=== 比较具有相同名称的文件 ==="
	if *diffCmdWrite {
		if _, writeErr := fileWrite.WriteString(matchFileNameFiles + "\n"); writeErr != nil {
			return fmt.Errorf("写入文件时出错: %v", writeErr)
		}
	} else {
		cl.Green(matchFileNameFiles)
	}
	var sameNameCount int
	// 遍历目录 A 中的文件
	for relPath, pathA := range filesA {
		if pathB, ok := filesB[relPath]; ok {
			// 如果目录 B 中存在同名文件，使用errgroup并行计算校验值
			var eg errgroup.Group
			var hashValueA, hashValueB string
			var errA, errB error

			// 使用errgroup并行计算校验值
			eg.Go(func() error {
				hashValueA, errA = checksum(pathA, hashType)
				return errA
			})

			eg.Go(func() error {
				hashValueB, errB = checksum(pathB, hashType)
				return errB
			})

			// 等待两个校验值计算完成
			if err := eg.Wait(); err != nil {
				if errA != nil {
					cl.PrintErrf("计算文件 %s 的 %s 值时出错: %v\n", *diffCmdType, pathA, errA)
				}
				if errB != nil {
					cl.PrintErrf("计算文件 %s 的 %s 值时出错: %v\n", *diffCmdType, pathB, errB)
				}
				continue
			}

			// 比较两个文件的校验值
			if hashValueA != hashValueB {
				diffCount++
				sameNameCount++
				// 根据 -w 参数决定是否将结果写入文件
				result := fmt.Sprintf("%d. 文件 %s 的 %s 值不同:\n  目录 A: %s  路径: %s\n  目录 B: %s  路径: %s", sameNameCount, relPath, *diffCmdType, getLast8Chars(hashValueA), pathA, getLast8Chars(hashValueB), pathB)
				if *diffCmdWrite {
					if _, writeErr := fileWrite.WriteString(result + "\n"); writeErr != nil {
						return fmt.Errorf("写入文件时出错: %v", writeErr)
					}
				} else {
					fmt.Println(result)
				}
			} else {
				sameCount++
			}

			// 从 filesB 中移除已比较的文件
			delete(filesB, relPath)
			// 从 filesA 中移除已比较的文件
			delete(filesA, relPath)
		}
	}

	// 如果没有相同文件则输出提示
	if sameCount == 0 {
		if *diffCmdWrite {
			if _, writeErr := fileWrite.WriteString("暂无相同文件\n"); writeErr != nil {
				return fmt.Errorf("写入文件时出错: %v", writeErr)
			}
		} else {
			fmt.Println("暂无相同文件")
		}
	}

	// 如果没有不同文件则输出提示
	if diffCount == 0 {
		if *diffCmdWrite {
			if _, writeErr := fileWrite.WriteString("暂无不同文件\n"); writeErr != nil {
				return fmt.Errorf("写入文件时出错: %v", writeErr)
			}
		} else {
			fmt.Println("暂无不同文件")
		}
	}

	// 检查仅存在于目录 A 的文件
	fileOnlyInDirectoryA := "\n=== 仅存在于目录 A 的文件 ==="
	if *diffCmdWrite {
		if _, writeErr := fileWrite.WriteString(fileOnlyInDirectoryA + "\n"); writeErr != nil {
			return fmt.Errorf("写入文件时出错: %v", writeErr)
		}
	} else {
		cl.Green(fileOnlyInDirectoryA)
	}
	onlyACount = len(filesA)
	var onlyACountDisplay int
	for relPath, pathA := range filesA {
		onlyACountDisplay++
		// 根据 -w 参数决定是否将结果写入文件
		result := fmt.Sprintf("%d. 文件 %s 仅存在于目录 A: %s", onlyACountDisplay, filepath.Base(relPath), pathA)
		if *diffCmdWrite {
			if _, writeErr := fileWrite.WriteString(result + "\n"); writeErr != nil {
				return fmt.Errorf("写入文件时出错: %v", writeErr)
			}
		} else {
			fmt.Println(result)
		}
	}
	if onlyACountDisplay == 0 {
		if *diffCmdWrite {
			if _, writeErr := fileWrite.WriteString("无匹配文件\n"); writeErr != nil {
				return fmt.Errorf("写入文件时出错: %v", writeErr)
			}
		} else {
			fmt.Println("无匹配文件")
		}
	}

	// 检查仅存在于目录 B 的文件
	fileOnlyInDirectoryB := "\n=== 仅存在于目录 B 的文件 ==="
	if *diffCmdWrite {
		if _, writeErr := fileWrite.WriteString(fileOnlyInDirectoryB + "\n"); writeErr != nil {
			return fmt.Errorf("写入文件时出错: %v", writeErr)
		}
	} else {
		cl.Green(fileOnlyInDirectoryB)
	}
	onlyBCount = len(filesB)
	var onlyBCountDisplay int
	for relPath, pathB := range filesB {
		onlyBCountDisplay++
		// 根据 -w 参数决定是否将结果写入文件
		result := fmt.Sprintf("%d. 文件 %s 仅存在于目录 B: %s", onlyBCountDisplay, filepath.Base(relPath), pathB)
		if *diffCmdWrite {
			if _, writeErr := fileWrite.WriteString(result + "\n"); writeErr != nil {
				return fmt.Errorf("写入文件时出错: %v", writeErr)
			}
		} else {
			fmt.Println(result)
		}
	}
	if onlyBCountDisplay == 0 {
		if *diffCmdWrite {
			if _, writeErr := fileWrite.WriteString("无匹配文件\n"); writeErr != nil {
				return fmt.Errorf("写入文件时出错: %v", writeErr)
			}
		} else {
			fmt.Println("无匹配文件")
		}
	}

	// 输出统计结果
	result := fmt.Sprintf("\n=== 统计结果 ===\n相同文件: %d\n不同文件: %d\n仅A目录文件: %d\n仅B目录文件: %d", sameCount, diffCount, onlyACount, onlyBCount)
	if *diffCmdWrite {
		if _, writeErr := fileWrite.WriteString(result + "\n"); writeErr != nil {
			return fmt.Errorf("写入文件时出错: %v", writeErr)
		}
	} else {
		cl.Green(result)
	}

	return nil
}

// compareDirWithCheckFile 对比校验文件与目录文件
func compareDirWithCheckFile(checkFileHashes globals.VirtualHashMap, targetFiles map[string]string, hashFunc func() hash.Hash, cl *colorlib.ColorLib, fileWrite *os.File) error {
	// 初始化统计计数器
	sameCount := 0
	diffCount := 0
	onlyCheckFileCount := 0
	onlyDirFileCount := 0

	// 比较相同文件名的文件
	matchFileNameFiles := "=== 比较具有相同名称的文件 ==="
	if *diffCmdWrite {
		if _, writeErr := fileWrite.WriteString(matchFileNameFiles + "\n"); writeErr != nil {
			return fmt.Errorf("写入文件 %s 失败: %v", globals.OutputFileName, writeErr)
		}
	} else {
		cl.Green(matchFileNameFiles)
	}
	var sameNameCount int
	// 遍历校验文件中的文件
	for virtualPath, checkEntry := range checkFileHashes {
		if targetPath, ok := targetFiles[virtualPath]; ok {
			// 如果目录中存在同名文件，计算其哈希值并比较
			hashValue, err := checksum(targetPath, hashFunc)
			if err != nil {
				cl.PrintErrf("计算文件 %s 的 %s 值时出错: %v\n", *diffCmdType, targetPath, err)
				continue
			}

			// 比较两个文件的校验值
			if hashValue != checkEntry.Hash {
				diffCount++
				sameNameCount++
				// 根据 -w 参数决定是否将结果写入文件
				result := fmt.Sprintf("%d. 文件 %s 的 %s 值不同:\n  校验文件: %s  路径: %s\n  目录文件: %s  路径: %s", sameNameCount, filepath.Base(virtualPath), *diffCmdType, getLast8Chars(checkEntry.Hash), checkEntry.RealPath, getLast8Chars(hashValue), targetPath)
				if *diffCmdWrite {
					if _, writeErr := fileWrite.WriteString(result + "\n"); writeErr != nil {
						return fmt.Errorf("写入文件 %s 失败: %v", globals.OutputFileName, writeErr)
					}
				} else {
					fmt.Println(result)
				}
			} else {
				sameCount++
			}

			delete(targetFiles, virtualPath)     // 从 targetFiles 中移除已比较的文件
			delete(checkFileHashes, virtualPath) // 从 checkFileHashes 中移除已比较的文件
		}
	}

	// 如果没有相同文件则输出提示
	if sameCount == 0 {
		if *diffCmdWrite {
			if _, writeErr := fileWrite.WriteString("暂无相同文件\n"); writeErr != nil {
				return fmt.Errorf("写入文件 %s 失败: %v", globals.OutputFileName, writeErr)
			}
		} else {
			fmt.Println("暂无相同文件")
		}
	}

	// 如果没有不同文件则输出提示
	if diffCount == 0 {
		if *diffCmdWrite {
			if _, writeErr := fileWrite.WriteString("暂无不同文件\n"); writeErr != nil {
				return fmt.Errorf("写入文件 %s 失败: %v", globals.OutputFileName, writeErr)
			}
		} else {
			fmt.Println("暂无不同文件")
		}
	}

	// 检查仅存在于校验文件中的文件
	validationFileOnly := "\n=== 仅存在于校验文件中的文件 ==="
	if *diffCmdWrite {
		if _, writeErr := fileWrite.WriteString(validationFileOnly + "\n"); writeErr != nil {
			return fmt.Errorf("写入文件 %s 失败: %v", globals.OutputFileName, writeErr)
		}
	} else {
		cl.Green(validationFileOnly)
	}
	onlyCheckFileCount = len(checkFileHashes)
	var onlyCheckFileCountDisplay int
	for virtualPath, checkEntry := range checkFileHashes {
		onlyCheckFileCountDisplay++
		// 根据 -w 参数决定是否将结果写入文件
		result := fmt.Sprintf("%d. 文件 %s 仅存在于校验文件: %s", onlyCheckFileCountDisplay, filepath.Base(virtualPath), checkEntry.RealPath)
		if *diffCmdWrite {
			if _, writeErr := fileWrite.WriteString(result + "\n"); writeErr != nil {
				return fmt.Errorf("写入文件 %s 失败: %v", globals.OutputFileName, writeErr)
			}
		} else {
			fmt.Println(result)
		}
	}
	if onlyCheckFileCountDisplay == 0 {
		if *diffCmdWrite {
			if _, writeErr := fileWrite.WriteString("无匹配文件\n"); writeErr != nil {
				return fmt.Errorf("写入文件 %s 失败: %v", globals.OutputFileName, writeErr)
			}
		} else {
			fmt.Println("无匹配文件")
		}
	}

	// 检查仅存在于目录中的文件
	directoryOnlyFile := "\n=== 仅存在于目录中的文件 ==="
	if *diffCmdWrite {
		if _, writeErr := fileWrite.WriteString(directoryOnlyFile + "\n"); writeErr != nil {
			return fmt.Errorf("写入文件 %s 失败: %v", globals.OutputFileName, writeErr)
		}
	} else {
		cl.Green(directoryOnlyFile)
	}
	onlyDirFileCount = len(targetFiles)
	var onlyDirFileCountDisplay int
	for virtualPath, targetPath := range targetFiles {
		onlyDirFileCountDisplay++
		// 根据 -w 参数决定是否将结果写入文件
		result := fmt.Sprintf("%d. 文件 %s 仅存在于目录: %s", onlyDirFileCountDisplay, virtualPath, targetPath)
		if *diffCmdWrite {
			if _, writeErr := fileWrite.WriteString(result + "\n"); writeErr != nil {
				return fmt.Errorf("写入文件 %s 失败: %v", globals.OutputFileName, writeErr)
			}
		} else {
			fmt.Println(result)
		}
	}
	if onlyDirFileCountDisplay == 0 {
		if *diffCmdWrite {
			if _, writeErr := fileWrite.WriteString("无匹配文件\n"); writeErr != nil {
				return fmt.Errorf("写入文件 %s 失败: %v", globals.OutputFileName, writeErr)
			}
		} else {
			fmt.Println("无匹配文件")
		}
	}

	// 输出统计结果
	result := fmt.Sprintf("\n=== 统计结果 ===\n相同文件: %d\n不同文件: %d\n仅校验文件: %d\n仅目录文件: %d", sameCount, diffCount, onlyCheckFileCount, onlyDirFileCount)
	if *diffCmdWrite {
		if _, writeErr := fileWrite.WriteString(result + "\n"); writeErr != nil {
			return fmt.Errorf("写入文件 %s 失败: %v", globals.OutputFileName, writeErr)
		}
	} else {
		cl.Green(result)
	}

	return nil
}

// checkWithFileAndDir 根据校验文件和目录进行校验的逻辑
func checkWithFileAndDir(checkFile, checkDir string, cl *colorlib.ColorLib) error {
	// 清理路径
	*diffCmdDirs = filepath.Clean(*diffCmdDirs)

	// 检查指定的目录是否包含禁止输入的路径
	if _, ok := globals.ForbiddenPaths[*diffCmdDirs]; ok {
		return fmt.Errorf("指定的目录包含禁止输入的路径: %s", *diffCmdDirs)
	}

	// 检查checkDir是否包含超过1个点, 如果包含则报错
	if strings.Contains(*diffCmdDirs, "..") {
		return fmt.Errorf("指定的目录包含禁止输入的路径: %s", *diffCmdDirs)
	}

	// 检查checkDir是否以分隔符结尾, 如果是则去掉
	if strings.HasSuffix(checkDir, string(filepath.Separator)) {
		checkDir = strings.TrimSuffix(checkDir, string(filepath.Separator))
	}

	// 检查checkDir的目录层级是否超过1
	if strings.Count(checkDir, string(filepath.Separator)) > 1 {
		return fmt.Errorf("目录 %s 的层级不能超过1", checkDir)
	}

	// 读取校验文件并加载到 map 中
	checkFileHashes, hashFunc, readErr := readHashFileToMap(checkFile, cl, true)
	if readErr != nil {
		return readErr
	}

	// 检查目录是否存在
	if _, statErr := os.Stat(checkDir); os.IsNotExist(statErr) {
		return fmt.Errorf("目录 %s 不存在", checkDir)
	}

	// 检查目录是否为绝对路径，如果不是，则转换为绝对路径
	if !filepath.IsAbs(checkDir) {
		var absErr error
		// 获取目录的绝对路径
		if checkDir, absErr = filepath.Abs(checkDir); absErr != nil {
			return fmt.Errorf("获取目录 %s 的绝对路径失败: %v", checkDir, absErr)
		}
	}

	// 获取指定目录下的文件列表
	targetFiles, err := getFiles(checkDir)
	if err != nil {
		return fmt.Errorf("读取目录 %s 时出错: %v", checkDir, err)
	}

	// 检查是否需要写入文件
	var fileWrite *os.File
	if *diffCmdWrite {
		var err error
		// 打开文件以写入
		fileWrite, err = os.OpenFile(globals.OutputCheckFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("打开文件 %s 失败: %v", globals.OutputCheckFileName, err)
		}
		defer fileWrite.Close()

		// 写入文件头
		if err := writeFileHeader(fileWrite, *hashCmdType, globals.TimestampFormat); err != nil {
			return fmt.Errorf("写入文件头失败: %v", err)
		}
	}

	// 替换目录map中的路径为基于虚拟父目录路径的绝对路径
	replaceDirFiles := replacePath(globals.VirtualRootDir, targetFiles)

	// 对比校验文件与目录文件
	if err := compareDirWithCheckFile(checkFileHashes, replaceDirFiles, hashFunc, cl, fileWrite); err != nil {
		return fmt.Errorf("校验文件与目录文件对比失败: %v", err)
	}

	// 如果是写入文件模式，则打印文件路径
	if *diffCmdWrite {
		cl.PrintOkf("比较结果已写入文件: %s", globals.OutputCheckFileName)
	}

	return nil
}

// replacePath 将相对路径替换为基于虚拟父目录路径的绝对路径
func replacePath(virtualBaseDir string, relativePaths map[string]string) map[string]string {
	// 创建一个新的映射，用于存储替换后的路径
	virtualPaths := make(map[string]string)

	// 遍历原始路径映射，将相对路径替换为基于虚拟父目录路径的绝对路径
	for relativePath, realPath := range relativePaths {
		// 构建虚拟路径 = 虚拟父目录路径 + 相对路径
		virtualPath := filepath.Join(virtualBaseDir, relativePath)
		// 将虚拟路径和真实路径添加到新的映射中
		virtualPaths[virtualPath] = realPath
	}

	return virtualPaths
}
