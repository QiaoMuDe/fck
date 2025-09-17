// Package check 实现校验文件解析功能
package check

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/commands/internal/types"
	"gitee.com/MM-Q/go-kit/hash"
)

// hashFileParser 校验文件解析器
type hashFileParser struct {
	validator *hashLineValidator
	cl        *colorlib.ColorLib
}

// newHashFileParser 创建校验文件解析器
func newHashFileParser(cl *colorlib.ColorLib) *hashFileParser {
	return &hashFileParser{
		validator: newHashLineValidator(),
		cl:        cl,
	}
}

// parseFile 解析校验文件
//
// 参数:
//   - checkFile: 校验文件路径
//   - userBaseDir: 用户指定的基准目录
//
// 返回值:
//   - types.VirtualHashMap: 虚拟哈希映射表
//   - string: 哈希算法类型
//   - error: 错误信息
func (p *hashFileParser) parseFile(checkFile string, userBaseDir string) (types.VirtualHashMap, string, error) {
	// 检查文件是否存在
	if _, err := os.Stat(checkFile); err != nil {
		return nil, "", fmt.Errorf("校验文件不存在: %s", checkFile)
	}

	// 打开文件
	file, err := os.Open(checkFile)
	if err != nil {
		return nil, "", fmt.Errorf("无法打开校验文件: %v", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			p.cl.PrintWarnf("关闭文件失败: %v\n", closeErr)
		}
	}()

	scanner := bufio.NewScanner(file)

	// 解析文件头
	headerInfo, err := p.parseHeader(scanner)
	if err != nil {
		return nil, "", err
	}

	// 解析文件内容
	hashMap, err := p.parseContent(scanner, headerInfo, userBaseDir)
	if err != nil {
		return nil, "", err
	}

	if len(hashMap) == 0 {
		return nil, "", fmt.Errorf("没有找到有效的校验文件内容")
	}

	return hashMap, headerInfo.HashType, nil
}

// parseHeader 解析校验文件头
//
// 参数:
//   - scanner: 文件扫描器
//
// 返回值:
//   - *types.ChecksumHeader: 校验文件头信息
//   - error: 错误信息
func (p *hashFileParser) parseHeader(scanner *bufio.Scanner) (*types.ChecksumHeader, error) {
	if !scanner.Scan() {
		return nil, fmt.Errorf("校验文件为空")
	}

	headerLine := scanner.Text()

	// 尝试解析新格式：#hashType#timestamp#mode#basePath 或 #hashType#timestamp#mode
	headerRegex := regexp.MustCompile(`^#(\w+)#([^#]+)(?:#([^#]+)(?:#(.+))?)?$`)
	matches := headerRegex.FindStringSubmatch(headerLine)

	if matches == nil {
		return nil, fmt.Errorf("校验文件头格式错误")
	}

	headerInfo := &types.ChecksumHeader{
		HashType:  matches[1], // hashType
		Timestamp: matches[2], // timestamp
	}

	// 检查哈希算法是否支持
	if headerInfo.HashType == "" {
		return nil, fmt.Errorf("校验文件头格式错误, 必须指定哈希算法")
	}

	if ok := hash.IsAlgorithmSupported(headerInfo.HashType); !ok {
		return nil, fmt.Errorf("不支持的哈希算法: %s", headerInfo.HashType)
	}

	// 解析模式和基准路径
	if len(matches) > 3 && matches[3] != "" {
		headerInfo.Mode = matches[3]
		if len(matches) > 4 && matches[4] != "" {
			headerInfo.BasePath = matches[4]
		}
	} else {
		// 兼容旧格式，默认为便携模式
		headerInfo.Mode = types.ChecksumModePortable
	}

	return headerInfo, nil
}

// parseContent 解析文件内容
//
// 参数:
//   - scanner: 文件扫描器
//   - headerInfo: 文件头信息
//   - userBaseDir: 用户指定基准目录
//
// 返回值:
//   - types.VirtualHashMap: 虚拟哈希映射表
//   - error: 错误信息
func (p *hashFileParser) parseContent(scanner *bufio.Scanner, headerInfo *types.ChecksumHeader, userBaseDir string) (types.VirtualHashMap, error) {
	hashMap := make(types.VirtualHashMap)
	lineNum := 1 // 从第二行开始计数（第一行是头部）

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// 验证并解析行内容
		hash, filePath, err := p.validator.validateLine(line, lineNum)
		if err != nil {
			p.cl.PrintErrorf("解析错误: %v\n", err)
			continue
		}

		// 跳过空行和注释行
		if hash == "" || filePath == "" {
			continue
		}

		// 解析文件路径
		resolvedPath, err := p.resolveFilePath(filePath, headerInfo, userBaseDir)
		if err != nil {
			p.cl.PrintErrorf("路径解析失败: %v\n", err)
			continue
		}

		hashMap[filePath] = types.VirtualHashEntry{
			RealPath: resolvedPath,
			Hash:     hash,
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取校验文件时出错: %v", err)
	}

	return hashMap, nil
}

// resolveFilePath 解析文件路径
//
// 参数:
//   - filePath: 原始文件路径
//   - headerInfo: 文件头信息
//   - userBaseDir: 用户指定基准目录
//
// 返回值:
//   - string: 解析后的绝对路径
//   - error: 错误信息
func (p *hashFileParser) resolveFilePath(filePath string, headerInfo *types.ChecksumHeader, userBaseDir string) (string, error) {
	// 1. 绝对路径直接使用
	if filepath.IsAbs(filePath) {
		return filePath, nil
	}

	// 2. 用户手动指定基准目录优先
	if userBaseDir != "" {
		return filepath.Join(userBaseDir, filePath), nil
	}

	// 3. 根据文件头模式自动处理
	switch headerInfo.Mode {
	case types.ChecksumModeLocal:
		// LOCAL模式：使用文件头中的基准路径
		if headerInfo.BasePath != "" {
			return filepath.Join(headerInfo.BasePath, filePath), nil
		}
		// 如果没有基准路径，降级为便携模式处理
		fallthrough
	case types.ChecksumModePortable, "": // 便携模式或旧格式（默认便携模式）
		// 直接使用当前目录作为基准目录
		return filepath.Join(".", filePath), nil
	default:
		return "", fmt.Errorf("未知的校验文件模式: %s", headerInfo.Mode)
	}
}
