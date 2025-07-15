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

	// (1) 执行根据单校验文件校验目录完整性的逻辑 -f 参数
	if diffCmdFile.Get() != "" && diffCmdDirs.Get() == "" && diffCmdDirA.Get() == "" && diffCmdDirB.Get() == "" {
		cl.PrintOk("正在校验目录完整性...")
		// 执行校验文件
		if err := fileCheck(diffCmdFile.Get(), cl); err != nil {
			return err
		}
		return nil
	}

	// (2) 如果指定校验文件不为空，同时也通过diffCmdDirs.Get()指定了目录 -f 参数和 -d 参数
	if diffCmdFile.Get() != "" && diffCmdDirs.Get() != "" && diffCmdDirA.Get() == "" && diffCmdDirB.Get() == "" {
		// 执行校验文件和目录的逻辑
		if err := checkWithFileAndDir(diffCmdFile.Get(), diffCmdDirs.Get(), cl); err != nil {
			return err
		}
		return nil
	}

	// (3) 执行对比目录A和目录B的逻辑 -a 参数和 -b 参数
	if diffCmdDirA.Get() == "" || diffCmdDirB.Get() == "" {
		return fmt.Errorf("对比目录时，必须同时指定 -a 和 -b 参数。或 -h 参数查看帮助信息。")
	}
	if diffCmdDirA.Get() != "" && diffCmdDirB.Get() != "" {
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
	if _, err := os.Stat(diffCmdDirA.Get()); err != nil {
		// 检查是否是权限错误
		if os.IsPermission(err) {
			// 如果是权限错误，返回错误信息
			return fmt.Errorf("权限不足: %s", diffCmdDirA.Get())
		}

		// 检查是否为文件不存在错误
		if os.IsNotExist(err) {
			// 如果是文件不存在错误，返回错误信息
			return fmt.Errorf("目录不存在: %s", diffCmdDirA.Get())
		}

		// 如果是其他错误，返回错误信息
		return fmt.Errorf("无法访问目录: %s", diffCmdDirA.Get())
	}
	if _, err := os.Stat(diffCmdDirB.Get()); err != nil {
		// 检查是否是权限错误
		if os.IsPermission(err) {
			// 如果是权限错误，返回错误信息
			return fmt.Errorf("权限不足: %s", diffCmdDirB.Get())
		}

		// 检查是否为文件不存在错误
		if os.IsNotExist(err) {
			// 如果是文件不存在错误，返回错误信息
			return fmt.Errorf("目录不存在: %s", diffCmdDirB.Get())
		}

		// 如果是其他错误，返回错误信息
		return fmt.Errorf("无法访问目录: %s", diffCmdDirB.Get())
	}

	// 检查目录A是否为绝对路径，如果不是，则转换为绝对路径
	if !filepath.IsAbs(diffCmdDirA.Get()) {
		absDirA, err := filepath.Abs(diffCmdDirA.Get())
		if err != nil {
			return fmt.Errorf("无法获取目录A的绝对路径: %v", err)
		}
		if setErr := diffCmdDirA.Set(absDirA); setErr != nil {
			return fmt.Errorf("无法设置目录A的绝对路径: %v", setErr)
		}
	}

	// 检查目录B是否为绝对路径，如果不是，则转换为绝对路径
	if !filepath.IsAbs(diffCmdDirB.Get()) {
		absDirB, err := filepath.Abs(diffCmdDirB.Get())
		if err != nil {
			return fmt.Errorf("无法获取目录B的绝对路径: %v", err)
		}
		if setErr := diffCmdDirB.Set(absDirB); setErr != nil {
			return fmt.Errorf("无法设置目录B的绝对路径: %v", setErr)
		}
	}

	// 校验目录A 和 目录B //
	// 检查指定的哈希算法是否有效
	hashType, ok := globals.SupportedAlgorithms[diffCmdType.Get()]
	if !ok {
		return fmt.Errorf("在校验哈希值时，哈希算法 %s 无效", diffCmdType.Get())
	}

	// 获取两个目录下的文件列表
	filesA, filesB, err := getFilesFromDirs(diffCmdDirA.Get(), diffCmdDirB.Get())
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	// 检查是否需要写入文件
	var fileWrite *os.File
	if diffCmdWrite.Get() {
		var err error
		// 打开文件以写入
		fileWrite, err = os.OpenFile(globals.OutputCheckFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("打开文件 %s 失败: %v", globals.OutputCheckFileName, err)
		}
		defer func() {
			if err := fileWrite.Close(); err != nil {
				cl.PrintErrf("colse file failed: %v\n", err)
			}
		}()

		// 写入文件头
		if err := writeFileHeader(fileWrite, hashCmdType.Get(), globals.TimestampFormat); err != nil {
			return fmt.Errorf("写入文件头失败: %v", err)
		}
	}

	// 比较两个目录的文件
	if err := compareFiles(filesA, filesB, hashType, cl, fileWrite); err != nil {
		return fmt.Errorf("比较文件失败: %v", err)
	}

	// 如果是写入文件模式，则打印文件路径
	if diffCmdWrite.Get() {
		cl.PrintOkf("比较结果已写入文件: %s\n", globals.OutputCheckFileName)
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
	defer func() {
		if err := checkFileRead.Close(); err != nil {
			fmt.Printf("close file failed: %v\n", err)
		}
	}()

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
				cl.PrintErrf("error: 校验文件格式错误, 文件 %s 的第 %d 行, %s\n", checkFile, lineCount, line)
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
				cl.PrintErrf("error: 校验文件格式错误, 文件 %s 的第 %d 行, %s\n", checkFile, lineCount, line)
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
	if diffCmdFile.Get() == "" {
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
			cl.PrintWarnf("在进行校验时，发现文件 %s 不存在，将跳过该文件的校验\n", entry.RealPath)
			continue
		}

		// 存储哈希值
		var hash string

		// 存储错误信息
		var checksumErr error

		// 计算哈希值
		hash, checksumErr = checksum(entry.RealPath, hashFunc)
		if checksumErr != nil {
			cl.PrintErrf("计算文件哈希失败: %v\n", checksumErr)
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
			cl.PrintErrf("文件 %s 不一致, 预期Hash值: %s, 实际Hash值: %s\n", filePath, getLast8Chars(checkEntry.Hash), getLast8Chars(targetHash))
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
	walkDirErr := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
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

	if walkDirErr != nil {
		// 如果是权限错误，则返回错误信息
		if os.IsPermission(walkDirErr) {
			return nil, fmt.Errorf("权限不足: %s", dir)
		}

		// 如果是文件不存在错误，则返回错误信息
		if os.IsNotExist(walkDirErr) {
			return nil, fmt.Errorf("目录不存在: %s", dir)
		}

		// 其他错误，返回错误信息
		return nil, fmt.Errorf("遍历目录时出错: %v", walkDirErr)
	}

	return files, nil
}

// compareFiles 比较两个目录的文件
func compareFiles(filesA, filesB map[string]string, hashType func() hash.Hash, cl *colorlib.ColorLib, fileWrite *os.File) error {
	// 初始化统计计数器
	sameCount := 0  // 相同文件计数
	diffCount := 0  // 不同文件计数
	onlyACount := 0 // 目录 A 中独有的文件计数
	onlyBCount := 0 // 目录 B 中独有的文件计数

	var hasCompared bool // 新增标志位，表示是否有进行过比较

	// 比较相同路径的文件
	matchFileNameFiles := "=== 比较具有相同名称的文件 ==="
	if diffCmdWrite.Get() {
		if _, writeErr := fileWrite.WriteString(matchFileNameFiles + "\n"); writeErr != nil {
			return fmt.Errorf("写入文件时出错: %v", writeErr)
		}
	} else {
		cl.Green(matchFileNameFiles)
	}

	// 初始化同名文件计数器
	var sameNameCount int

	// 遍历目录 A 中的文件
	for relPath, pathA := range filesA {
		if pathB, ok := filesB[relPath]; ok {
			// 标记已进行过比较
			hasCompared = true

			// 获取文件大小
			fileInfoA, err := os.Lstat(pathA)
			if err != nil {
				cl.PrintErrf("获取文件 %s 的大小时出错: %v\n", pathA, err)
				continue
			}
			fileInfoB, err := os.Lstat(pathB)
			if err != nil {
				cl.PrintErrf("获取文件 %s 的大小时出错: %v\n", pathB, err)
				continue
			}

			// 比较文件大小
			if fileInfoA.Size() != fileInfoB.Size() {
				diffCount++     // 增加不同文件计数
				sameNameCount++ // 增加同名文件计数
				// 根据 -w 参数决定是否将结果写入文件
				result := fmt.Sprintf("%d. 文件 %s 的大小不同:\n  目录 A: %d 字节  路径: %s\n  目录 B: %d 字节  路径: %s", sameNameCount, relPath, fileInfoA.Size(), pathA, fileInfoB.Size(), pathB)
				if diffCmdWrite.Get() {
					if _, writeErr := fileWrite.WriteString(result + "\n"); writeErr != nil {
						return fmt.Errorf("写入文件时出错: %v", writeErr)
					}
				} else {
					fmt.Println(result)
				}
			} else {
				// 如果文件大小一致，使用 errgroup 并行计算校验值
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
						cl.PrintErrf("计算文件 %s 的 %s 值时出错: %v\n", diffCmdType.Get(), pathA, errA)
					}
					if errB != nil {
						cl.PrintErrf("计算文件 %s 的 %s 值时出错: %v\n", diffCmdType.Get(), pathB, errB)
					}
					continue
				}

				// 比较两个文件的校验值
				if hashValueA != hashValueB {
					diffCount++     // 增加不同文件计数
					sameNameCount++ // 增加同名文件计数
					// 根据 -w 参数决定是否将结果写入文件
					result := fmt.Sprintf("%d. 文件 %s 的 %s 值不同:\n  目录 A: %s  路径: %s\n  目录 B: %s  路径: %s", sameNameCount, relPath, diffCmdType.Get(), getLast8Chars(hashValueA), pathA, getLast8Chars(hashValueB), pathB)
					if diffCmdWrite.Get() {
						if _, writeErr := fileWrite.WriteString(result + "\n"); writeErr != nil {
							return fmt.Errorf("写入文件时出错: %v", writeErr)
						}
					} else {
						fmt.Println(result)
					}
				} else {
					sameCount++ // 增加相同文件计数
				}
			}

			// 从 filesB 中移除已比较的文件
			delete(filesB, relPath)
			// 从 filesA 中移除已比较的文件
			delete(filesA, relPath)
		}
	}

	// 根据比较结果输出提示
	if hasCompared {
		if sameCount == 0 && diffCount == 0 {
			if diffCmdWrite.Get() {
				if _, writeErr := fileWrite.WriteString("暂无匹配文件\n"); writeErr != nil {
					return fmt.Errorf("写入文件时出错: %v", writeErr)
				}
			} else {
				fmt.Println("暂无匹配文件")
			}

		} else {
			if sameCount > 0 {
				if diffCmdWrite.Get() {
					if _, writeErr := fmt.Fprintf(fileWrite, "找到 %d 个相同文件\n", sameCount); writeErr != nil {
						return fmt.Errorf("写入文件时出错: %v", writeErr)
					}
				} else {
					fmt.Printf("找到 %d 个相同文件\n", sameCount)
				}
			}
			if diffCount > 0 {
				if diffCmdWrite.Get() {
					if _, writeErr := fmt.Fprintf(fileWrite, "找到 %d 个不同文件\n", diffCount); writeErr != nil {
						return fmt.Errorf("写入文件时出错: %v", writeErr)
					}
				} else {
					fmt.Printf("找到 %d 个不同文件\n", diffCount)
				}
			}
		}
	} else {
		if diffCmdWrite.Get() {
			if _, writeErr := fileWrite.WriteString("无匹配文件\n"); writeErr != nil {
				return fmt.Errorf("写入文件时出错: %v", writeErr)
			}
		} else {
			fmt.Println("无匹配文件")
		}
	}

	// 检查仅存在于目录 A 的文件
	fileOnlyInDirectoryA := "\n=== 仅存在于目录 A 的文件 ==="
	if diffCmdWrite.Get() {
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
		if diffCmdWrite.Get() {
			if _, writeErr := fileWrite.WriteString(result + "\n"); writeErr != nil {
				return fmt.Errorf("写入文件时出错: %v", writeErr)
			}
		} else {
			fmt.Println(result)
		}
	}
	if onlyACountDisplay == 0 {
		if diffCmdWrite.Get() {
			if _, writeErr := fileWrite.WriteString("无匹配文件\n"); writeErr != nil {
				return fmt.Errorf("写入文件时出错: %v", writeErr)
			}
		} else {
			fmt.Println("无匹配文件")
		}
	}

	// 检查仅存在于目录 B 的文件
	fileOnlyInDirectoryB := "\n=== 仅存在于目录 B 的文件 ==="
	if diffCmdWrite.Get() {
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
		if diffCmdWrite.Get() {
			if _, writeErr := fileWrite.WriteString(result + "\n"); writeErr != nil {
				return fmt.Errorf("写入文件时出错: %v", writeErr)
			}
		} else {
			fmt.Println(result)
		}
	}
	if onlyBCountDisplay == 0 {
		if diffCmdWrite.Get() {
			if _, writeErr := fileWrite.WriteString("无匹配文件\n"); writeErr != nil {
				return fmt.Errorf("写入文件时出错: %v", writeErr)
			}
		} else {
			fmt.Println("无匹配文件")
		}
	}

	// 输出统计结果
	result := fmt.Sprintf("\n=== 统计结果 ===\n相同文件: %d\n不同文件: %d\n仅A目录文件: %d\n仅B目录文件: %d", sameCount, diffCount, onlyACount, onlyBCount)
	if diffCmdWrite.Get() {
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

	var hasCompared bool // 新增标志位，表示是否有进行过比较

	// 比较相同文件名的文件
	matchFileNameFiles := "=== 比较具有相同名称的文件 ==="
	if diffCmdWrite.Get() {
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
			// 标记已进行过比较
			hasCompared = true

			// 如果目录中存在同名文件，计算其哈希值并比较
			hashValue, err := checksum(targetPath, hashFunc)
			if err != nil {
				cl.PrintErrf("计算文件 %s 的 %s 值时出错: %v\n", diffCmdType.Get(), targetPath, err)
				continue
			}

			// 比较两个文件的校验值
			if hashValue != checkEntry.Hash {
				diffCount++
				sameNameCount++
				// 根据 -w 参数决定是否将结果写入文件
				result := fmt.Sprintf("%d. 文件 %s 的 %s 值不同:\n  校验文件: %s  路径: %s\n  目录文件: %s  路径: %s", sameNameCount, filepath.Base(virtualPath), diffCmdType.Get(), getLast8Chars(checkEntry.Hash), checkEntry.RealPath, getLast8Chars(hashValue), targetPath)
				if diffCmdWrite.Get() {
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

	// 根据比较结果输出提示
	if hasCompared {
		if sameCount == 0 && diffCount == 0 {
			if diffCmdWrite.Get() {
				if _, writeErr := fileWrite.WriteString("暂无匹配文件\n"); writeErr != nil {
					return fmt.Errorf("写入文件时出错: %v", writeErr)
				}
			} else {
				fmt.Println("暂无匹配文件")
			}

		} else {
			if sameCount > 0 {
				if diffCmdWrite.Get() {
					if _, writeErr := fmt.Fprintf(fileWrite, "找到 %d 个相同文件\n", sameCount); writeErr != nil {
						return fmt.Errorf("写入文件时出错: %v", writeErr)
					}
				} else {
					fmt.Printf("找到 %d 个相同文件\n", sameCount)
				}
			}
			if diffCount > 0 {
				if diffCmdWrite.Get() {
					if _, writeErr := fmt.Fprintf(fileWrite, "找到 %d 个不同文件\n", diffCount); writeErr != nil {
						return fmt.Errorf("写入文件时出错: %v", writeErr)
					}
				} else {
					fmt.Printf("找到 %d 个不同文件\n", diffCount)
				}
			}
		}
	} else {
		if diffCmdWrite.Get() {
			if _, writeErr := fileWrite.WriteString("无匹配文件\n"); writeErr != nil {
				return fmt.Errorf("写入文件时出错: %v", writeErr)
			}
		} else {
			fmt.Println("无匹配文件")
		}
	}

	// 检查仅存在于校验文件中的文件
	validationFileOnly := "\n=== 仅存在于校验文件中的文件 ==="
	if diffCmdWrite.Get() {
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
		if diffCmdWrite.Get() {
			if _, writeErr := fileWrite.WriteString(result + "\n"); writeErr != nil {
				return fmt.Errorf("写入文件 %s 失败: %v", globals.OutputFileName, writeErr)
			}
		} else {
			fmt.Println(result)
		}
	}
	if onlyCheckFileCountDisplay == 0 {
		if diffCmdWrite.Get() {
			if _, writeErr := fileWrite.WriteString("无匹配文件\n"); writeErr != nil {
				return fmt.Errorf("写入文件 %s 失败: %v", globals.OutputFileName, writeErr)
			}
		} else {
			fmt.Println("无匹配文件")
		}
	}

	// 检查仅存在于目录中的文件
	directoryOnlyFile := "\n=== 仅存在于目录中的文件 ==="
	if diffCmdWrite.Get() {
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
		if diffCmdWrite.Get() {
			if _, writeErr := fileWrite.WriteString(result + "\n"); writeErr != nil {
				return fmt.Errorf("写入文件 %s 失败: %v", globals.OutputFileName, writeErr)
			}
		} else {
			fmt.Println(result)
		}
	}
	if onlyDirFileCountDisplay == 0 {
		if diffCmdWrite.Get() {
			if _, writeErr := fileWrite.WriteString("无匹配文件\n"); writeErr != nil {
				return fmt.Errorf("写入文件 %s 失败: %v", globals.OutputFileName, writeErr)
			}
		} else {
			fmt.Println("无匹配文件")
		}
	}

	// 输出统计结果
	result := fmt.Sprintf("\n=== 统计结果 ===\n相同文件: %d\n不同文件: %d\n仅校验文件: %d\n仅目录文件: %d", sameCount, diffCount, onlyCheckFileCount, onlyDirFileCount)
	if diffCmdWrite.Get() {
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
	if setErr := diffCmdDirs.Set(filepath.Clean(diffCmdDirs.Get())); setErr != nil {
		return fmt.Errorf("清理路径 %s 失败: %v", diffCmdDirs.Get(), setErr)
	}

	// 检查指定的目录是否包含禁止输入的路径
	if _, ok := globals.ForbiddenPaths[diffCmdDirs.Get()]; ok {
		return fmt.Errorf("指定的目录包含禁止输入的路径: %s", diffCmdDirs.Get())
	}

	// 检查checkDir是否包含超过1个点, 如果包含则报错
	if strings.Contains(diffCmdDirs.Get(), "..") {
		return fmt.Errorf("指定的目录包含禁止输入的路径: %s", diffCmdDirs.Get())
	}

	// 检查checkDir是否以分隔符结尾, 如果是则去掉
	if strings.HasSuffix(checkDir, string(filepath.Separator)) {
		checkDir = strings.TrimSuffix(checkDir, string(filepath.Separator))
	}

	// 检查checkDir的目录层级是否超过1
	if strings.Count(checkDir, string(filepath.Separator)) > 1 {
		return fmt.Errorf("目录 %s 的层级不能超过1", checkDir)
	}

	// 检查目录是否存在
	if _, statErr := os.Stat(checkDir); os.IsNotExist(statErr) {
		return fmt.Errorf("目录 %s 不存在", checkDir)
	}

	// 检查目录是否为绝对路径，如果是则退出
	if filepath.IsAbs(checkDir) {
		return fmt.Errorf("目录 %s 不能为绝对路径", checkDir)
	}

	// 读取校验文件并加载到 map 中
	checkFileHashes, hashFunc, readErr := readHashFileToMap(checkFile, cl, true)
	if readErr != nil {
		return readErr
	}

	// 获取指定目录下的文件列表
	targetFiles, err := getFiles(checkDir)
	if err != nil {
		return fmt.Errorf("读取目录 %s 时出错: %v", checkDir, err)
	}

	// 检查是否需要写入文件
	var fileWrite *os.File
	if diffCmdWrite.Get() {
		var err error
		// 打开文件以写入
		fileWrite, err = os.OpenFile(globals.OutputCheckFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("打开文件 %s 失败: %v", globals.OutputCheckFileName, err)
		}
		defer func() {
			if err := fileWrite.Close(); err != nil {
				cl.PrintErrf("close file failed: %v\n", err)
			}
		}()

		// 写入文件头
		if err := writeFileHeader(fileWrite, hashCmdType.Get(), globals.TimestampFormat); err != nil {
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
	if diffCmdWrite.Get() {
		cl.PrintOkf("比较结果已写入文件: %s\n", globals.OutputCheckFileName)
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
