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
	"gitee.com/MM-Q/fck/commands/internal/types"
)

// hashFileParser 校验文件解析器
type hashFileParser struct {
	validator *hashLineValidator
	cl        *colorlib.ColorLib
}

// newHashFileParser 创建新的文件解析器
func newHashFileParser(cl *colorlib.ColorLib) *hashFileParser {
	return &hashFileParser{
		validator: newHashLineValidator(),
		cl:        cl,
	}
}

// parseFile 解析校验文件
func (p *hashFileParser) parseFile(checkFile string, isRelPath bool) (types.VirtualHashMap, func() hash.Hash, error) {
	// 检查文件是否存在
	if _, err := os.Stat(checkFile); err != nil {
		return nil, nil, fmt.Errorf("校验文件不存在: %s", checkFile)
	}

	// 打开文件
	file, err := os.Open(checkFile)
	if err != nil {
		return nil, nil, fmt.Errorf("无法打开校验文件: %v", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			p.cl.PrintWarnf("关闭文件失败: %v\n", closeErr)
		}
	}()

	scanner := bufio.NewScanner(file)

	// 解析文件头
	hashFunc, err := p.parseHeader(scanner)
	if err != nil {
		return nil, nil, err
	}

	// 解析文件内容
	hashMap, err := p.parseContent(scanner, isRelPath)
	if err != nil {
		return nil, nil, err
	}

	if len(hashMap) == 0 {
		return nil, nil, fmt.Errorf("没有找到有效的校验文件内容")
	}

	return hashMap, hashFunc, nil
}

// parseHeader 解析校验文件头
func (p *hashFileParser) parseHeader(scanner *bufio.Scanner) (func() hash.Hash, error) {
	if !scanner.Scan() {
		return nil, fmt.Errorf("校验文件为空")
	}

	headerLine := scanner.Text()
	headerRegex := regexp.MustCompile(`^#(\w+)#(.+)$`)
	matches := headerRegex.FindStringSubmatch(headerLine)

	if matches == nil {
		return nil, fmt.Errorf("校验文件头格式错误, 格式应为 #hashType#timestamp")
	}

	hashType := matches[1]
	if hashType == "" {
		return nil, fmt.Errorf("校验文件头格式错误, 必须指定哈希算法")
	}

	hashFunc, ok := types.SupportedAlgorithms[hashType]
	if !ok {
		return nil, fmt.Errorf("不支持的哈希算法: %s", hashType)
	}

	return hashFunc, nil
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
