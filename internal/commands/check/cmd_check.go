// Package check 实现了文件完整性校验命令的主要逻辑。
// 该文件包含 check 子命令的入口函数，负责解析校验文件并执行文件完整性验证。
package check

import (
	"fmt"
	"os"

	"gitee.com/MM-Q/color"
	"gitee.com/MM-Q/fck/internal/types"
)

// CheckConfig 校验配置结构体
type CheckConfig struct {
	CheckFile string // 校验文件路径
	BaseDir   string // 基准目录
	Verbose   bool   // 详细模式（显示校验通过的文件）
}

// CheckCmdMain 执行校验命令
//
// 参数:
//   - cl: 颜色库实例
//   - config: 校验配置
//
// 返回值:
//   - error: 校验过程中可能发生的错误
func CheckCmdMain(cl *color.GlobalColor, config CheckConfig) error {
	checkFile := config.CheckFile
	if checkFile == "" {
		checkFile = types.OutputFileName
	}

	if _, err := os.Stat(checkFile); err != nil {
		return fmt.Errorf("checksum file not found: %s", checkFile)
	}

	cl.Blue("checking file integrity...")

	parser := newHashFileParser(cl)

	userBaseDir := config.BaseDir

	hashMap, hashFunc, err := parser.parseFile(checkFile, userBaseDir)
	if err != nil {
		return fmt.Errorf("failed to parse checksum file: %v", err)
	}

	checker := newFileChecker(cl, hashFunc, config.Verbose)

	if err := checker.checkFiles(hashMap); err != nil {
		return fmt.Errorf("checksum verification failed: %v", err)
	}

	return nil
}
