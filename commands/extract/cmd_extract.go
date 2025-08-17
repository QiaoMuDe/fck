package extract

import (
	"gitee.com/MM-Q/colorlib"
)

func ExtractCmdMain(cl *colorlib.ColorLib) error {
	// 创建解压配置
	config, err := newExtractConfig()
	if err != nil {
		return err
	}

	// 根据格式选择解压器
	switch config.extractType {
	}

	return nil
}
