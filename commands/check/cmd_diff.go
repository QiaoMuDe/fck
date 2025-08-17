package check

import (
	"fmt"
	"os"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/commands/internal/types"
)

// CheckCmdMain 是 check 命令的主函数
func CheckCmdMain(cl *colorlib.ColorLib) error {
	// 获取校验文件路径
	checkFile := checkCmd.Arg(0)
	if checkFile == "" {
		checkFile = types.OutputFileName
	}

	// 检查校验文件是否存在
	if _, err := os.Stat(checkFile); err != nil {
		return fmt.Errorf("校验文件不存在: %s", checkFile)
	}

	cl.PrintOk("正在校验目录完整性...")

	// 创建解析器
	parser := newHashFileParser(cl)

	// 解析校验文件
	hashMap, hashFunc, err := parser.parseFile(checkFile, false)
	if err != nil {
		return fmt.Errorf("解析校验文件失败: %v", err)
	}

	// 创建校验器
	checker := newFileChecker(cl, hashFunc)

	// 执行文件校验
	if err := checker.checkFiles(hashMap); err != nil {
		return fmt.Errorf("文件校验失败: %v", err)
	}

	return nil
}
