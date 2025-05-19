package cmd

import (
	"flag"
	"fmt"
	"hash"
	"path/filepath"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/globals"
)

func hashCmdMain(cmd *flag.FlagSet, cl *colorlib.ColorLib) error {
	// 获取指定的路径
	targetPath := cmd.Arg(0)

	// 如果没有指定路径，则打印帮助信息并退出
	if targetPath == "" {
		return fmt.Errorf("在校验哈希值时，必须指定一个路径")
	}

	// 清理路径
	targetPath = filepath.Clean(targetPath)

	// 定义一个切片来存储文件列表
	var files []string

	// 通过glob获取文件列表
	files, globRrr := filepath.Glob(targetPath)
	if globRrr != nil {
		return fmt.Errorf("在校验哈希值时，路径 %s 无效: %v", targetPath, globRrr)
	}

	// 检查文件列表是否为空
	if len(files) == 0 {
		return fmt.Errorf("在校验哈希值时，路径 %s 没有找到任何文件", targetPath)
	}

	// 检查指定的哈希算法是否有效
	hashType, ok := globals.SupportedAlgorithms[*hashCmdType]
	if !ok {
		return fmt.Errorf("在校验哈希值时，哈希算法 %s 无效", *hashCmdType)
	}

	// 执行哈希值校验
	if runErr := hashCmdRun(hashType, targetPath); runErr != nil {
		return fmt.Errorf("在校验哈希值时，发生错误: %v", runErr)
	}

	return nil
}

func hashCmdRun(hashType func() hash.Hash, targetPath string) error {
	//
	return nil
}
