package cmd

import (
	"bufio"
	"fmt"
	"hash"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/globals"

	"golang.org/x/sync/errgroup"
)

// checkCmdMain 是 check 命令的主函数
func checkCmdMain(cl *colorlib.ColorLib) error {
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
	if *checkCmdFile == "" && *checkCmdDirA == "" && *checkCmdDirB == "" {
		return fmt.Errorf("必须指定一个校验文件或两个目录。或 -h 参数查看帮助信息。")
	}

	// (1) 执行根据单校验文件校验目录完整性的逻辑 -f 参数
	if *checkCmdFile != "" && *checkCmdDirs == "" {
		if err := fileCheck(*checkCmdFile, cl); err != nil {
			return err
		}
		return nil
	}

	// (2) 如果指定校验文件不为空，同时也通过*checkCmdDirs指定了目录 -f 参数和 -d 参数
	if *checkCmdFile != "" && *checkCmdDirs != "" {
		// 检查指定的目录是否为.或者..
		if *checkCmdDirs == "." || *checkCmdDirs == ".." || *checkCmdDirs == "./" || *checkCmdDirs == "../" {
			return fmt.Errorf("不能指定目录为.或者.., 请指定校验目录的上级目录或 -h 参数查看帮助信息")
		}

		// 执行校验文件和目录的逻辑
		if err := checkWithFileAndDir(*checkCmdFile, *checkCmdDirs, cl); err != nil {
			return err
		}
		return nil
	}

	// (3) 执行对比目录A和目录B的逻辑 -a 参数和 -b 参数
	if *checkCmdDirA == "" || *checkCmdDirB == "" {
		return fmt.Errorf("对比目录时，必须同时指定 -a 和 -b 参数。或 -h 参数查看帮助信息。")
	}
	if *checkCmdDirA != "" && *checkCmdDirB != "" {
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
	if _, err := os.Stat(*checkCmdDirA); err != nil {
		return fmt.Errorf("目录A不存在: %s", *checkCmdDirA)
	}
	if _, err := os.Stat(*checkCmdDirB); err != nil {
		return fmt.Errorf("目录B不存在: %s", *checkCmdDirB)
	}

	// 校验目录A 和 目录B //
	// 检查指定的哈希算法是否有效
	hashType, ok := globals.SupportedAlgorithms[*checkCmdType]
	if !ok {
		return fmt.Errorf("在校验哈希值时，哈希算法 %s 无效", *checkCmdType)
	}

	// 获取两个目录下的文件列表
	filesA, filesB, err := getFilesFromDirs(*checkCmdDirA, *checkCmdDirB)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	// 检查是否需要写入文件
	var fileWrite *os.File
	if *checkCmdWrite {
		var err error
		// 打开文件以写入
		fileWrite, err = os.OpenFile(globals.OutputCheckFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("打开文件 %s 失败: %v", globals.OutputCheckFileName, err)
		}
		defer fileWrite.Close()

		// 获取时间
		now := time.Now()

		// 写入文件头
		if _, err := fileWrite.WriteString(fmt.Sprintf("#%s#%s\n\n", *checkCmdType, now.Format("2006-01-02 15:04:05"))); err != nil {
			return fmt.Errorf("写入文件头失败: %v", err)
		}
	}

	// 比较两个目录的文件
	compareFiles(filesA, filesB, hashType, cl, fileWrite)

	// 如果是写入文件模式，则打印文件路径
	if *checkCmdWrite {
		cl.PrintOkf("比较结果已写入文件: %s", globals.OutputCheckFileName)
	}

	return nil
}

// readHashFileToMap 读取校验文件并加载到 map 中
func readHashFileToMap(checkFile string, cl *colorlib.ColorLib, isRelPath bool) (map[string]string, func() hash.Hash, error) {
	// 存储校验值文件的 map
	checkFileHashes := make(map[string]string)

	// 检查校验文件是否存在
	if _, err := os.Stat(checkFile); err != nil {
		return nil, nil, fmt.Errorf("校验文件不存在: %s", checkFile)
	}

	// 读取校验文件
	checkFileRead, openErr := os.OpenFile(checkFile, os.O_RDONLY, 0644)
	if openErr != nil {
		return nil, nil, fmt.Errorf("无法打开校验文件: %v", openErr)
	}
	defer checkFileRead.Close()

	// 获取哈希算法
	headerInfo := make([]byte, 1024)
	n, readErr := checkFileRead.ReadAt(headerInfo, 0)
	if readErr != nil && readErr != io.EOF {
		return nil, nil, fmt.Errorf("读取校验文件时出错: %v", readErr)
	}
	headerInfo = headerInfo[:n] // 调整切片长度以匹配实际读取的字节数

	// 检查是否是#开头的文件
	if len(headerInfo) == 0 || headerInfo[0] != '#' {
		return nil, nil, fmt.Errorf("校验文件头格式错误, 必须以#开头")
	}

	// 解析哈希算法和时间戳，以井号为分隔符
	parts := strings.Split(string(headerInfo), "#")
	if len(parts) < 3 {
		return nil, nil, fmt.Errorf("校验文件头格式错误, 格式应为 #hashType#timestamp")
	}
	hashType := parts[1] // 哈希算法
	if hashType == "" {
		return nil, nil, fmt.Errorf("校验文件头格式错误, 必须指定哈希算法")
	}

	// 检查哈希算法是否支持
	hashFunc, ok := globals.SupportedAlgorithms[string(hashType)]
	if !ok {
		return nil, nil, fmt.Errorf("不支持的哈希算法: %s", string(hashType))
	}

	// 重置文件指针到开头
	_, seekErr := checkFileRead.Seek(0, io.SeekStart)
	if seekErr != nil {
		return nil, nil, fmt.Errorf("重置文件指针时出错: %v", seekErr)
	}

	// 解析校验文件内容
	scanner := bufio.NewScanner(checkFileRead)

	// 定义一个计数器
	var lineCount int

	// 逐行读取校验文件
	for scanner.Scan() { // 逐行读取
		line := scanner.Text() // 获取当前行的文本

		lineCount++ // 计数器加 1

		// 如果当前行是空行，则跳过
		if line == "" {
			continue
		}

		// 如果当前行以#开头，则跳过
		if strings.HasPrefix(line, "#") {
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

			// 存储到 map 中
			checkFileHashes[relPath] = expectedHash
		} else {
			// 如果哈希值或文件路径为空，则跳过
			if expectedHash == "" || filePath == "" {
				cl.PrintErrf("error: 校验文件格式错误, 文件 %s 的第 %d 行, %s", checkFile, lineCount, line)
				continue
			}

			// 存储到 map 中
			checkFileHashes[filePath] = expectedHash
		}
	}

	// 检查是否有错误发生
	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("读取校验文件时出错: %v", err)
	}

	// 检查map是否为空
	if len(checkFileHashes) == 0 {
		return nil, nil, fmt.Errorf("没有找到有效的校验文件内容")
	}

	return checkFileHashes, hashFunc, nil
}

// fileCheck 根据单校验文件校验目录完整性的逻辑 -f 参数
func fileCheck(checkFile string, cl *colorlib.ColorLib) error {
	// 检查校验文件是否为空
	if *checkCmdFile == "" {
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
	for file := range checkFileHashes {
		// 检查文件是否存在，若不存在则输出更详细的提示信息并跳过当前文件处理
		if _, err := os.Stat(file); err != nil {
			cl.PrintWarnf("在进行校验时，发现文件 %s 不存在，将跳过该文件的校验", file)
			continue
		}

		var hash string       // 存储哈希值
		var checksumErr error // 存储错误信息
		// 计算哈希值
		hash, checksumErr = checksum(file, hashFunc)
		if checksumErr != nil {
			cl.PrintErrf("计算文件哈希失败: %v", checksumErr)
		}
		// 将哈希值存储在目标目录哈希值映射中
		targetDirHashes[file] = hash
	}

	// 对比哈希值
	var checkCount int
	for filePath, checkHash := range checkFileHashes {
		// 获取实际的哈希值
		targetHash, ok := targetDirHashes[filePath]
		if !ok {
			// 如果实际的哈希值不存在，则跳过
			continue
		}

		// 比较哈希值
		if targetHash != checkHash {
			cl.PrintErrf("文件 %s 不一致, 预期Hash值: %s, 实际Hash值: %s", filePath, getLast8Chars(checkHash), getLast8Chars(targetHash))
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
func compareFiles(filesA, filesB map[string]string, hashType func() hash.Hash, cl *colorlib.ColorLib, fileWrite *os.File) {
	// 初始化统计计数器
	sameCount := 0
	diffCount := 0
	onlyACount := 0
	onlyBCount := 0

	// 比较相同路径的文件
	if *checkCmdWrite {
		fileWrite.WriteString("=== 比较具有相同名称的文件 ===\n")
	} else {
		cl.Green("=== 比较具有相同名称的文件 ===")
	}
	var sameNameCount int
	// 遍历目录 A 中的文件
	for path, pathA := range filesA {
		if pathB, ok := filesB[path]; ok {
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
					cl.PrintErrf("计算文件 %s 的 %s 值时出错: %v\n", *checkCmdType, pathA, errA)
				}
				if errB != nil {
					cl.PrintErrf("计算文件 %s 的 %s 值时出错: %v\n", *checkCmdType, pathB, errB)
				}
				continue
			}

			// 比较两个文件的校验值
			if hashValueA != hashValueB {
				diffCount++
				sameNameCount++
				// 根据 -w 参数决定是否将结果写入文件
				if *checkCmdWrite {
					// 根据 -full 参数决定是否显示完整路径
					if *checkCmdFullPath {
						fileWrite.WriteString(fmt.Sprintf("%d. 文件 %s 的 %s 值不同:\n  目录 A: %s\n  目录 B: %s\n", sameNameCount, path, *checkCmdType, getLast8Chars(hashValueA), getLast8Chars(hashValueB)))
					} else {
						fileWrite.WriteString(fmt.Sprintf("%d. 文件 %s 的 %s 值不同:\n  目录 A: %s\n  目录 B: %s\n", sameNameCount, filepath.Base(path), *checkCmdType, getLast8Chars(hashValueA), getLast8Chars(hashValueB)))
					}
				} else {
					// 根据 -full 参数决定是否显示完整路径
					if *checkCmdFullPath {
						fmt.Printf("%d. 文件 %s 的 %s 值不同:\n  目录 A: %s\n  目录 B: %s\n", sameNameCount, path, *checkCmdType, getLast8Chars(hashValueA), getLast8Chars(hashValueB))
					} else {
						fmt.Printf("%d. 文件 %s 的 %s 值不同:\n  目录 A: %s\n  目录 B: %s\n", sameNameCount, filepath.Base(path), *checkCmdType, getLast8Chars(hashValueA), getLast8Chars(hashValueB))
					}
				}
			} else {
				sameCount++
			}

			// 从 filesB 中移除已比较的文件
			delete(filesB, path)
			// 从 filesA 中移除已比较的文件
			delete(filesA, path)
		}
	}

	// 如果没有相同文件则输出提示
	if sameCount == 0 {
		if *checkCmdWrite {
			fileWrite.WriteString("暂无相同文件\n")
		} else {
			fmt.Println("暂无相同文件")
		}
	}

	// 如果没有不同文件则输出提示
	if diffCount == 0 {
		if *checkCmdWrite {
			fileWrite.WriteString("暂无不同文件\n")
		} else {
			fmt.Println("暂无不同文件")
		}
	}

	// 检查仅存在于目录 A 的文件
	if *checkCmdWrite {
		fileWrite.WriteString("\n=== 仅存在于目录 A 的文件 ===\n")
	} else {
		cl.Green("\n=== 仅存在于目录 A 的文件 ===")
	}
	onlyACount = len(filesA)
	var onlyACountDisplay int
	for path, pathA := range filesA {
		onlyACountDisplay++
		// 根据 -w 参数决定是否将结果写入文件
		if *checkCmdWrite {
			// 根据 -full 参数决定是否显示完整路径
			if *checkCmdFullPath {
				fileWrite.WriteString(fmt.Sprintf("%d. 文件 %s 仅存在于目录 A: %s\n", onlyACountDisplay, path, pathA))
			} else {
				fileWrite.WriteString(fmt.Sprintf("%d. 文件 %s 仅存在于目录 A: %s\n", onlyACountDisplay, filepath.Base(path), pathA))
			}
		} else {
			// 根据 -full 参数决定是否显示完整路径
			if *checkCmdFullPath {
				fmt.Printf("%d. 文件 %s 仅存在于目录 A: %s\n", onlyACountDisplay, path, pathA)
			} else {
				fmt.Printf("%d. 文件 %s 仅存在于目录 A: %s\n", onlyACountDisplay, filepath.Base(path), pathA)
			}
		}
	}
	if onlyACountDisplay == 0 {
		if *checkCmdWrite {
			fileWrite.WriteString("无匹配文件\n")
		} else {
			fmt.Println("无匹配文件")
		}
	}

	// 检查仅存在于目录 B 的文件
	if *checkCmdWrite {
		fileWrite.WriteString("\n=== 仅存在于目录 B 的文件 ===\n")
	} else {
		cl.Green("\n=== 仅存在于目录 B 的文件 ===")
	}
	onlyBCount = len(filesB)
	var onlyBCountDisplay int
	for path, pathB := range filesB {
		onlyBCountDisplay++
		// 根据 -w 参数决定是否将结果写入文件
		if *checkCmdWrite {
			// 根据 -full 参数决定是否显示完整路径
			if *checkCmdFullPath {
				fileWrite.WriteString(fmt.Sprintf("%d. 文件 %s 仅存在于目录 B: %s\n", onlyBCountDisplay, path, pathB))
			} else {
				fileWrite.WriteString(fmt.Sprintf("%d. 文件 %s 仅存在于目录 B: %s\n", onlyBCountDisplay, filepath.Base(path), pathB))
			}
		} else {
			// 根据 -full 参数决定是否显示完整路径
			if *checkCmdFullPath {
				fmt.Printf("%d. 文件 %s 仅存在于目录 B: %s\n", onlyBCountDisplay, path, pathB)
			} else {
				fmt.Printf("%d. 文件 %s 仅存在于目录 B: %s\n", onlyBCountDisplay, filepath.Base(path), pathB)
			}
		}
	}
	if onlyBCountDisplay == 0 {
		if *checkCmdWrite {
			fileWrite.WriteString("无匹配文件\n")
		} else {
			fmt.Println("无匹配文件")
		}
	}

	// 输出统计结果
	if *checkCmdWrite {
		fileWrite.WriteString(fmt.Sprintf("\n=== 统计结果 ===\n相同文件: %d\n不同文件: %d\n仅A目录文件: %d\n仅B目录文件: %d\n", sameCount, diffCount, onlyACount, onlyBCount))
	} else {
		cl.Green(fmt.Sprintf("\n=== 统计结果 ===\n相同文件: %d\n不同文件: %d\n仅A目录文件: %d\n仅B目录文件: %d", sameCount, diffCount, onlyACount, onlyBCount))
	}
}

// compareDirWithCheckFile 对比校验文件与目录文件
func compareDirWithCheckFile(checkFileHashes map[string]string, dirFiles map[string]string, hashFunc func() hash.Hash, cl *colorlib.ColorLib, fileWrite *os.File) {
	// 初始化统计计数器
	sameCount := 0
	diffCount := 0
	onlyCheckFileCount := 0
	onlyDirFileCount := 0

	// 比较相同文件名的文件
	if *checkCmdWrite {
		fileWrite.WriteString("=== 比较具有相同名称的文件 ===\n")
	} else {
		cl.Green("=== 比较具有相同名称的文件 ===")
	}
	var sameNameCount int
	// 遍历校验文件中的文件
	for path, checkHash := range checkFileHashes {
		if dirPath, ok := dirFiles[path]; ok {
			// 如果目录中存在同名文件，计算其哈希值并比较
			hashValue, err := checksum(dirPath, hashFunc)
			if err != nil {
				cl.PrintErrf("计算文件 %s 的 %s 值时出错: %v\n", *checkCmdType, dirPath, err)
				continue
			}

			// 比较两个文件的校验值
			if hashValue != checkHash {
				diffCount++
				sameNameCount++
				// 根据 -w 参数决定是否将结果写入文件
				if *checkCmdWrite {
					// 根据 -full 参数决定是否显示完整路径
					if *checkCmdFullPath {
						fileWrite.WriteString(fmt.Sprintf("%d. 文件 %s 的 %s 值不同:\n  校验文件: %s\n  目录文件: %s\n", sameNameCount, path, *checkCmdType, getLast8Chars(checkHash), getLast8Chars(hashValue)))
					} else {
						fileWrite.WriteString(fmt.Sprintf("%d. 文件 %s 的 %s 值不同:\n  校验文件: %s\n  目录文件: %s\n", sameNameCount, filepath.Base(path), *checkCmdType, getLast8Chars(checkHash), getLast8Chars(hashValue)))
					}
				} else {
					// 根据 -full 参数决定是否显示完整路径
					if *checkCmdFullPath {
						fmt.Printf("%d. 文件 %s 的 %s 值不同:\n  校验文件: %s\n  目录文件: %s\n", sameNameCount, path, *checkCmdType, getLast8Chars(checkHash), getLast8Chars(hashValue))
					} else {
						fmt.Printf("%d. 文件 %s 的 %s 值不同:\n  校验文件: %s\n  目录文件: %s\n", sameNameCount, filepath.Base(path), *checkCmdType, getLast8Chars(checkHash), getLast8Chars(hashValue))
					}
				}
			} else {
				sameCount++
			}

			delete(dirFiles, path)        // 从 dirFiles 中移除已比较的文件
			delete(checkFileHashes, path) // 从 checkFileHashes 中移除已比较的文件
		}
	}

	// 如果没有相同文件则输出提示
	if sameCount == 0 {
		if *checkCmdWrite {
			fileWrite.WriteString("暂无相同文件\n")
		} else {
			fmt.Println("暂无相同文件")
		}
	}

	// 如果没有不同文件则输出提示
	if diffCount == 0 {
		if *checkCmdWrite {
			fileWrite.WriteString("暂无不同文件\n")
		} else {
			fmt.Println("暂无不同文件")
		}
	}

	// 检查仅存在于校验文件中的文件
	if *checkCmdWrite {
		fileWrite.WriteString("\n=== 仅存在于校验文件中的文件 ===\n")
	} else {
		cl.Green("\n=== 仅存在于校验文件中的文件 ===")
	}
	onlyCheckFileCount = len(checkFileHashes)
	var onlyCheckFileCountDisplay int
	for path, checkHash := range checkFileHashes {
		onlyCheckFileCountDisplay++
		// 根据 -w 参数决定是否将结果写入文件
		if *checkCmdWrite {
			// 根据 -full 参数决定是否显示完整路径
			if *checkCmdFullPath {
				fileWrite.WriteString(fmt.Sprintf("%d. 文件 %s 仅存在于校验文件: %s\n", onlyCheckFileCountDisplay, path, checkHash))
			} else {
				fileWrite.WriteString(fmt.Sprintf("%d. 文件 %s 仅存在于校验文件: %s\n", onlyCheckFileCountDisplay, filepath.Base(path), checkHash))
			}
		} else {
			// 根据 -full 参数决定是否显示完整路径
			if *checkCmdFullPath {
				fmt.Printf("%d. 文件 %s 仅存在于校验文件: %s\n", onlyCheckFileCountDisplay, path, checkHash)
			} else {
				fmt.Printf("%d. 文件 %s 仅存在于校验文件: %s\n", onlyCheckFileCountDisplay, filepath.Base(path), checkHash)
			}
		}
	}
	if onlyCheckFileCountDisplay == 0 {
		if *checkCmdWrite {
			fileWrite.WriteString("无匹配文件\n")
		} else {
			fmt.Println("无匹配文件")
		}
	}

	// 检查仅存在于目录中的文件
	if *checkCmdWrite {
		fileWrite.WriteString("\n=== 仅存在于目录中的文件 ===\n")
	} else {
		cl.Green("\n=== 仅存在于目录中的文件 ===")
	}
	onlyDirFileCount = len(dirFiles)
	var onlyDirFileCountDisplay int
	for path, dirPath := range dirFiles {
		onlyDirFileCountDisplay++
		// 根据 -w 参数决定是否将结果写入文件
		if *checkCmdWrite {
			// 根据 -full 参数决定是否显示完整路径
			if *checkCmdFullPath {
				fileWrite.WriteString(fmt.Sprintf("%d. 文件 %s 仅存在于目录: %s\n", onlyDirFileCountDisplay, path, dirPath))
			} else {
				fileWrite.WriteString(fmt.Sprintf("%d. 文件 %s 仅存在于目录: %s\n", onlyDirFileCountDisplay, filepath.Base(path), dirPath))
			}
		} else {
			// 根据 -full 参数决定是否显示完整路径
			if *checkCmdFullPath {
				fmt.Printf("%d. 文件 %s 仅存在于目录: %s\n", onlyDirFileCountDisplay, path, dirPath)
			} else {
				fmt.Printf("%d. 文件 %s 仅存在于目录: %s\n", onlyDirFileCountDisplay, filepath.Base(path), dirPath)
			}
		}
	}
	if onlyDirFileCountDisplay == 0 {
		if *checkCmdWrite {
			fileWrite.WriteString("无匹配文件\n")
		} else {
			fmt.Println("无匹配文件")
		}
	}

	// 输出统计结果
	if *checkCmdWrite {
		fileWrite.WriteString(fmt.Sprintf("\n=== 统计结果 ===\n相同文件: %d\n不同文件: %d\n仅校验文件: %d\n仅目录文件: %d\n", sameCount, diffCount, onlyCheckFileCount, onlyDirFileCount))
	} else {
		cl.Green(fmt.Sprintf("\n=== 统计结果 ===\n相同文件: %d\n不同文件: %d\n仅校验文件: %d\n仅目录文件: %d", sameCount, diffCount, onlyCheckFileCount, onlyDirFileCount))
	}
}

// checkWithFileAndDir 根据校验文件和目录进行校验的逻辑
func checkWithFileAndDir(checkFile, checkDir string, cl *colorlib.ColorLib) error {
	// 读取校验文件并加载到 map 中
	checkFileHashes, hashFunc, err := readHashFileToMap(checkFile, cl, true)
	if err != nil {
		return err
	}

	// 获取指定目录下的文件列表
	dirFiles, err := getFiles(checkDir)
	if err != nil {
		return fmt.Errorf("读取目录 %s 时出错: %v", checkDir, err)
	}

	// 检查是否需要写入文件
	var fileWrite *os.File
	if *checkCmdWrite {
		var err error
		// 打开文件以写入
		fileWrite, err = os.OpenFile(globals.OutputCheckFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("打开文件 %s 失败: %v", globals.OutputCheckFileName, err)
		}
		defer fileWrite.Close()

		// 获取时间
		now := time.Now()

		// 写入文件头
		if _, err := fileWrite.WriteString(fmt.Sprintf("#%s#%s\n\n", *checkCmdType, now.Format("2006-01-02 15:04:05"))); err != nil {
			return fmt.Errorf("写入文件头失败: %v", err)
		}
	}

	// 对比校验文件与目录文件
	compareDirWithCheckFile(checkFileHashes, dirFiles, hashFunc, cl, fileWrite)

	// 如果是写入文件模式，则打印文件路径
	if *checkCmdWrite {
		cl.PrintOkf("比较结果已写入文件: %s", globals.OutputCheckFileName)
	}

	return nil
}
