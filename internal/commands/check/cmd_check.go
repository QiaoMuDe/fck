// Package check 实现了文件完整性校验命令的主要逻辑。
// 该文件包含 check 子命令的入口函数，负责解析校验文件并执行文件完整性验证。
package check

import (
	"fmt"
	"os"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/internal/types"
)

// CheckConfig 校验配置结构体
type CheckConfig struct {
	CheckFile string // 校验文件路径
	BaseDir   string // 基准目录
	Quiet     bool   // 静默模式
	Color     bool   // 颜色输出
}

// CheckCmdMain 执行校验命令
//
// 参数:
//   - cl: 颜色库实例
//   - config: 校验配置
//
// 返回值:
//   - error: 校验过程中可能发生的错误
func CheckCmdMain(cl *colorlib.ColorLib, config CheckConfig) error {
	checkFile := config.CheckFile
	if checkFile == "" {
		checkFile = types.OutputFileName
	}

	cl.SetColor(config.Color)

	if _, err := os.Stat(checkFile); err != nil {
		return fmt.Errorf("指定的校验文件不存在: %s, 请确认文件路径是否正确", checkFile)
	}

	cl.Blue("正在校验完整性...")

	parser := newHashFileParser(cl)

	userBaseDir := config.BaseDir

	hashMap, hashFunc, err := parser.parseFile(checkFile, userBaseDir)
	if err != nil {
		return fmt.Errorf("解析校验文件失败: %v", err)
	}

	checker := newFileChecker(cl, hashFunc, config.Quiet)

	if err := checker.checkFiles(hashMap); err != nil {
		return fmt.Errorf("文件校验失败: %v", err)
	}

	return nil
}
