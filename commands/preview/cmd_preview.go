package preview

import (
	"errors"

	"gitee.com/MM-Q/comprx"
)

func PreviewCmdMain() error {
	// 获取压缩包路径
	packPath := previewCmd.Arg(0)
	if packPath == "" {
		return errors.New("压缩包路径不能为空")
	}

	// 打印压缩包信息
	if infoFlag.Get() {
		if err := comprx.PrintArchiveInfo(packPath); err != nil {
			return err
		}
		return nil
	}

	// 以简洁方式打印文件列表
	if lsFlag.Get() {
		// 限制打印文件数量
		if limitFlag.Get() != 0 {
			if err := comprx.PrintFilesLimit(packPath, limitFlag.Get(), false); err != nil {
				return err
			}
			return nil
		}

		// 打印所有文件
		if err := comprx.PrintLs(packPath); err != nil {
			return err
		}
		return nil
	}

	// 以详细方式打印文件列表
	if llFlag.Get() {
		// 限制打印文件数量
		if limitFlag.Get() != 0 {
			if err := comprx.PrintFilesLimit(packPath, limitFlag.Get(), true); err != nil {
				return err
			}
			return nil
		}

		// 打印所有文件详细信息
		if err := comprx.PrintLl(packPath); err != nil {
			return err
		}
		return nil
	}

	// 默认打印压缩包的信息和所有文件
	if err := comprx.PrintArchiveAndFiles(packPath, true); err != nil {
		return err
	}

	return nil
}
