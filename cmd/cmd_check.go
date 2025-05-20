package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/globals"
)

func checkCmdMain(cl *colorlib.ColorLib) error {
	// 校验逻辑
	// 先获取指定目录的所有文件的哈希值并存储到map中，键是文件名，值是哈希值
	// 然后逐行读取checksum.hash文件，获取文件的哈希值和文件名并与map中的哈希值进行比较
	// 如果不一致，则打印错误信息
	// 并检查两边的是否有校验文件中有的文件，目标目录中没有的文件
	// 又或者是目标目录中有的文件，校验文件中没有的文件

	if *checkCmdFile != "" {
		if err := fileCheck(*checkCmdFile, cl); err != nil {
			return err
		}
	}

	return nil
}

func fileCheck(checkFile string, cl *colorlib.ColorLib) error {
	// 存储校验值文件的map
	checkFileHashes := make(map[string]string)

	// 存储目标目录文件的map
	targetDirHashes := make(map[string]string)

	// 检查校验文件是否为空
	if *checkCmdFile == "" {
		return fmt.Errorf("在校验文件时，必须指定一个校验文件 checksum.hash")
	}

	// 检查校验文件是否存在
	if _, err := os.Stat(*checkCmdFile); err != nil {
		return fmt.Errorf("校验文件不存在: %s", *checkCmdFile)
	}

	// 读取校验文件
	checkFileRead, openErr := os.OpenFile(*checkCmdFile, os.O_RDONLY, 0644)
	if openErr != nil {
		return fmt.Errorf("无法打开校验文件: %v", openErr)
	}
	defer checkFileRead.Close()

	// 获取哈希算法
	headerInfo := make([]byte, 1024)
	n, readErr := checkFileRead.ReadAt(headerInfo, 0)
	if readErr != nil && readErr != io.EOF {
		return fmt.Errorf("读取校验文件时出错: %v", readErr)
	}
	headerInfo = headerInfo[:n] // 调整切片长度以匹配实际读取的字节数

	// 检查是否是#开头的文件
	// 原代码中 headerInfo 是 nil 切片，直接访问索引会报错，这里先检查切片长度
	if len(headerInfo) == 0 || headerInfo[0] != '#' {
		return fmt.Errorf("校验文件头格式错误, 必须以#开头")
	}
	// 解析哈希算法，以井号为分隔符
	// 解析哈希算法和时间戳，以井号为分隔符
	parts := strings.Split(string(headerInfo), "#")
	if len(parts) < 3 {
		fmt.Println("校验文件头格式错误, 格式应为 #hashType#timestamp")
		return fmt.Errorf("校验文件头格式错误, 格式应为 #hashType#timestamp")
	}
	hashType := parts[1] // 哈希算法
	if hashType == "" {
		fmt.Println("校验文件头格式错误, 必须指定哈希算法")
		return fmt.Errorf("校验文件头格式错误, 必须指定哈希算法")
	}

	// 检查哈希算法是否支持
	hashFunc, ok := globals.SupportedAlgorithms[string(hashType)]
	if !ok {
		return fmt.Errorf("不支持的哈希算法: %s", string(hashType))
	}

	// 重置文件指针到开头
	_, seekErr := checkFileRead.Seek(0, io.SeekStart)
	if seekErr != nil {
		return fmt.Errorf("重置文件指针时出错: %v", seekErr)
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
		if len(parts) != 2 {          // 如果分割后的长度不为 2，则跳过
			cl.PrintErrf("error: 校验文件格式错误, 文件 %s 的第 %d 行, %s", checkFile, lineCount, line)
			continue
		}
		expectedHash := parts[0]                 // 哈希值
		filePath := strings.Join(parts[1:], " ") // 文件路径

		// 去除路径中的引号
		filePath = strings.Trim(filePath, `"`)

		// 存储到 map 中
		checkFileHashes[filePath] = expectedHash
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取校验文件时出错: %v", err)
	}

	// 计算目标目录文件哈希值
	for file := range checkFileHashes {
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
			cl.PrintErrf("文件 %s 不一致, 预期Hash值: %s, 实际Hash值: %s", filePath, checkHash[len(checkHash)-8:], targetHash[len(targetHash)-8:])
			checkCount++
		}
	}

	// 检查checkCount是否为0
	if checkCount == 0 {
		cl.PrintOk("校验成功，无文件差异")
	}

	return nil
}
