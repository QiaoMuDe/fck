package check

import (
	"bufio"
	"fmt"
	"hash"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/commands/internal/common"
	"gitee.com/MM-Q/fck/commands/internal/types"
)

// CheckCmdMain 是 check 命令的主函数
func CheckCmdMain(cl *colorlib.ColorLib) error {
	// 获取校验文件路径
	checkFile := checkCmd.Arg(0)
	if checkFile == "" {
		checkFile = types.OutputFileName
	}

	// 检查checkFile是否存在
	if _, err := os.Stat(checkFile); err != nil {
		return fmt.Errorf("校验文件不存在: %s", checkFile)
	}

	cl.PrintOk("正在校验目录完整性...")
	// 执行校验文件
	if err := fileCheck(checkFile, cl); err != nil {
		return err
	}
	return nil
}

// readHashFileToMap 读取校验文件并加载到 map 中
func readHashFileToMap(checkFile string, cl *colorlib.ColorLib, isRelPath bool) (types.VirtualHashMap, func() hash.Hash, error) {
	// 创建一个新的映射，用于存储替换后的路径
	replaceMap := make(types.VirtualHashMap)

	// 检查校验文件是否存在
	if _, err := os.Stat(checkFile); err != nil {
		return nil, nil, fmt.Errorf("校验文件不存在: %s", checkFile)
	}

	// 打开校验文件
	checkFileRead, openErr := os.OpenFile(checkFile, os.O_RDONLY, 0600)
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
				cl.PrintErrorf("error: 校验文件格式错误, 文件 %s 的第 %d 行, %s\n", checkFile, lineCount, line)
				continue
			}

			// 构建虚拟路径
			virtualPath := filepath.Join(types.VirtualRootDir, relPath)

			// 存储到 map 中
			replaceMap[virtualPath] = types.VirtualHashEntry{
				RealPath: filePath,     // 真实路径
				Hash:     expectedHash, // 哈希值
			}
		} else {
			// 如果哈希值或文件路径为空，则跳过
			if expectedHash == "" || filePath == "" {
				cl.PrintErrorf("error: 校验文件格式错误, 文件 %s 的第 %d 行, %s\n", checkFile, lineCount, line)
				continue
			}

			// 存储到 map 中
			replaceMap[filePath] = types.VirtualHashEntry{
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
	hashFunc, ok := types.SupportedAlgorithms[hashType]
	if !ok {
		return nil, fmt.Errorf("不支持的哈希算法: %s", hashType)
	}

	return hashFunc, nil
}

// fileCheck 根据单校验文件校验目录完整性的逻辑
func fileCheck(checkFile string, cl *colorlib.ColorLib) error {
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
		hash, checksumErr = common.Checksum(entry.RealPath, hashFunc)
		if checksumErr != nil {
			cl.PrintErrorf("计算文件哈希失败: %v\n", checksumErr)
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
			cl.PrintErrorf("文件 %s 不一致, 预期Hash值: %s, 实际Hash值: %s\n", filePath, common.GetLast8Chars(checkEntry.Hash), common.GetLast8Chars(targetHash))
			checkCount++
		}
	}

	// 检查 checkCount 是否为 0
	if checkCount == 0 {
		cl.PrintOk("校验成功，无文件差异")
	}

	return nil
}
