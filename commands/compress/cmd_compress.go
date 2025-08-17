package compress

import (
	"gitee.com/MM-Q/colorlib"
)

func CompressCmdMain(cl *colorlib.ColorLib) error {
	// 创建压缩配置
	config, err := newCompressConfig()
	if err != nil {
		return err
	}

	// 根据格式选择压缩器
	switch config.compressType {
	}

	return nil
}
