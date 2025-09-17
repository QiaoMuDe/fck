// Package check 实现了校验文件的解析功能。
// 该文件提供了校验文件解析器，用于解析包含文件哈希信息的校验文件格式。
package check

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/commands/internal/types"
	"gitee.com/MM-Q/go-kit/hash"
)

// hashFileParser 校验文件解析器
type hashFileParser struct {
	validator *hashLineValidator
	cl        *colorlib.ColorLib
}

// HeaderInfo 文件头信息
type HeaderInfo struct {
	HashType  string
	Timestamp string
	Mode      string // PORTABLE, LOCAL, 或空（兼容旧格式）
	BasePath  string // LOCAL模式的基准路径
}

// newHashFileParser 创建新的文件解析器
func newHashFileParser(cl *colorlib.ColorLib) *hashFileParser {
	return &hashFileParser{
		validator: newHashLineValidator(),
		cl:        cl,
	}
}

// parseFile 解析校验文件
func (p *hashFileParser) parseFile(checkFile string, isRelPath bool) (types.VirtualHashMap, string, error) {
	// 检查文件是否存在
	if _, err := os.Stat(checkFile); err != nil {
		return nil, "", fmt.Errorf("校验文件不存在: %s", checkFile)
	}

	// 打开文件
	file, err := os.Open(checkFile)
	if err != nil {
		return nil, "", fmt.Errorf("无法打开校验文件: %v", err)
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)

	// 解析文件头
	hashType, err := p.parseHeader(scanner)
	if err != nil {
		return nil, "", err
	}

	// 解析文件内容
	hashMap, err := p.parseContent(scanner, isRelPath)
	if err != nil {
		return nil, "", err
	}

	if len(hashMap) == 0 {
		return nil, "", fmt.Errorf("没有找到有效的校验文件内容")
	}

	return hashMap, hashType, nil
}

// parseFileEnhanced 解析校验文件（增强版）
func (p *hashFileParser) parseFileEnhanced(checkFile string, userBaseDir string) (types.VirtualHashMap, string, error) {
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

	// 解析文件头（增强版）
	headerInfo, err := p.parseHeaderEnhanced(scanner)
	if err != nil {
		return nil, "", err
	}

	// 解析文件内容（增强版）
	hashMap, err := p.parseContentEnhanced(scanner, headerInfo, userBaseDir)
	if err != nil {
		return nil, "", err
	}

	if len(hashMap) == 0 {
		return nil, "", fmt.Errorf("没有找到有效的校验文件内容")
	}

	return hashMap, headerInfo.HashType, nil
}

// parseHeader 解析校验文件头
func (p *hashFileParser) parseHeader(scanner *bufio.Scanner) (string, error) {
	if !scanner.Scan() {
		return "", fmt.Errorf("校验文件为空")
	}

	headerLine := scanner.Text()
	headerRegex := regexp.MustCompile(`^#(\w+)#(.+)$`)
	matches := headerRegex.FindStringSubmatch(headerLine)

	if matches == nil {
		return "", fmt.Errorf("校验文件头格式错误, 格式应为 #hashType#timestamp")
	}

	hashType := matches[1]
	if hashType == "" {
		return "", fmt.Errorf("校验文件头格式错误, 必须指定哈希算法")
	}

	if ok := hash.IsAlgorithmSupported(hashType); !ok {
		return "", fmt.Errorf("不支持的哈希算法: %s", hashType)
	}

	return hashType, nil
}

// parseHeaderEnhanced 解析校验文件头（增强版）
func (p *hashFileParser) parseHeaderEnhanced(scanner *bufio.Scanner) (*HeaderInfo, error) {
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

	headerInfo := &HeaderInfo{
		HashType:  matches[1],
		Timestamp: matches[2],
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
func (p *hashFileParser) parseContent(scanner *bufio.Scanner, isRelPath bool) (types.VirtualHashMap, error) {
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

		// 处理相对路径逻辑
		if isRelPath {
			virtualPath, realPath, err := p.processRelativePath(filePath)
			if err != nil {
				p.cl.PrintErrorf("处理相对路径失败: %v\n", err)
				continue
			}
			hashMap[virtualPath] = types.VirtualHashEntry{
				RealPath: realPath,
				Hash:     hash,
			}
		} else {
			hashMap[filePath] = types.VirtualHashEntry{
				RealPath: filePath,
				Hash:     hash,
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取校验文件时出错: %v", err)
	}

	return hashMap, nil
}

// parseContentEnhanced 解析文件内容（增强版）
func (p *hashFileParser) parseContentEnhanced(scanner *bufio.Scanner, headerInfo *HeaderInfo, userBaseDir string) (types.VirtualHashMap, error) {
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

// resolveFilePath 智能路径解析
func (p *hashFileParser) resolveFilePath(filePath string, headerInfo *HeaderInfo, userBaseDir string) (string, error) {
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

// processRelativePath 处理相对路径逻辑
func (p *hashFileParser) processRelativePath(filePath string) (virtualPath, realPath string, err error) {
	// 检查空路径
	if filePath == "" {
		return "", "", fmt.Errorf("无效的文件路径: %s", filePath)
	}

	// 使用 / 作为分隔符来解析路径，因为校验文件中通常使用 /
	parts := strings.Split(strings.ReplaceAll(filePath, "\\", "/"), "/")
	if len(parts) == 0 {
		return "", "", fmt.Errorf("无效的文件路径: %s", filePath)
	}

	// 检查是否只有根目录
	if len(parts) == 1 {
		return "", "", fmt.Errorf("获取相对路径时出错: Rel: can't make %s relative to C", filePath)
	}

	// 跳过第一个部分（根目录），获取相对路径
	relParts := parts[1:]
	if len(relParts) == 0 {
		return "", "", fmt.Errorf("相对路径为空")
	}

	relPath := strings.Join(relParts, string(filepath.Separator))
	virtualPath = filepath.Join(types.VirtualRootDir, relPath)
	return virtualPath, filePath, nil
}
