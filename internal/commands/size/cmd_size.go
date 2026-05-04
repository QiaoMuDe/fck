// Package size 实现了文件和目录大小计算功能。
// 该文件提供了计算文件、目录大小的核心功能，支持通配符路径展开、隐藏文件处理、进度显示和表格输出。
package size

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/color"
	"gitee.com/MM-Q/fck/internal/types"
	"gitee.com/MM-Q/fck/internal/utils"
	gfs "gitee.com/MM-Q/go-kit/fs"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/schollz/progressbar/v3"
)

type item struct {
	Name string
	Size string
}

type items []item

type SizeConfig struct {
	Args          []string
	NoColor       bool
	TableStyle    string
	Hidden        bool
	Human         bool
	FollowSymlink bool
}

func SizeCmdMain(cl *color.GlobalColor, config SizeConfig) error {
	targetPaths := config.Args

	if len(targetPaths) == 0 {
		targetPaths = []string{"*"}
	}

	cl.SetNoColor(config.NoColor)

	var itemList items

	for _, targetPath := range targetPaths {
		targetPath = filepath.Clean(targetPath)

		pathsToProcess, err := expandPath(targetPath)
		if err != nil {
			cl.Redf("failed to expand path: %v\n", err)
			continue
		}

		for _, path := range pathsToProcess {
			addPathToList(path, &itemList, cl, config)
		}
	}

	if len(itemList) > 0 {
		printSizeTable(itemList, cl, config)
	}

	return nil
}

func expandPath(path string) ([]string, error) {
	if !strings.Contains(path, "*") {
		return []string{path}, nil
	}

	filePaths, err := filepath.Glob(path)
	if err != nil {
		return nil, err
	}

	if len(filePaths) == 0 {
		return nil, fmt.Errorf("no matching files found: %s", path)
	}

	return filePaths, nil
}

func addPathToList(path string, itemList *items, cl *color.GlobalColor, config SizeConfig) {
	if !config.Hidden && gfs.IsHidden(path) {
		return
	}

	size, err := getPathSize(path, config)
	if err != nil {
		cl.Redf("failed to calculate size of: %s - %v\n", path, err)
		return
	}

	var sizeStr string
	if config.Human {
		sizeStr = humanReadableSize(size, 2)
	} else {
		sizeStr = fmt.Sprintf("%d", size)
	}

	*itemList = append(*itemList, item{
		Name: path,
		Size: sizeStr,
	})
}

func getPathSize(path string, config SizeConfig) (int64, error) {
	includeHidden := config.Hidden
	if !includeHidden && gfs.IsHidden(path) {
		return 0, nil
	}

	// 根据配置选择是否跟随符号链接
	var info os.FileInfo
	var err error
	if config.FollowSymlink {
		info, err = os.Stat(path)
	} else {
		info, err = os.Lstat(path)
	}
	if err != nil {
		switch {
		case os.IsPermission(err):
			return 0, fmt.Errorf("permission denied for: path %s", path)
		case os.IsNotExist(err):
			return 0, fmt.Errorf("file not found: path %s", path)
		default:
			return 0, fmt.Errorf("failed to get file info: path %s, error: %v", path, err)
		}
	}

	if !info.IsDir() {
		return info.Size(), nil
	}

	var totalSize int64
	var skippedFiles int

	bar := progressbar.NewOptions64(
		-1,
		progressbar.OptionClearOnFinish(),
		progressbar.OptionSetDescription(filepath.Base(path)+"++"),
	)
	defer func() {
		_ = bar.Finish()
		_ = bar.Close()
	}()

	walkErr := filepath.WalkDir(path, func(filePath string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			if !os.IsNotExist(err) {
				skippedFiles++
			}
			return nil
		}

		if filePath == path {
			return nil
		}

		if !includeHidden && gfs.IsHidden(filePath) {
			return nil
		}

		// 处理符号链接
		if dirEntry.Type()&os.ModeSymlink != 0 {
			if config.FollowSymlink {
				// 跟随符号链接，获取目标信息
				targetInfo, statErr := os.Stat(filePath)
				if statErr != nil {
					skippedFiles++
					return nil
				}
				if !targetInfo.IsDir() {
					totalSize += targetInfo.Size()
					_ = bar.Add64(targetInfo.Size())
				}
				// 如果是软链接的是目录则跳过（不计入大小）
			}
			// 不跟随符号链接时，跳过（不计入大小）
			return nil
		}

		if !dirEntry.IsDir() {
			fileInfo, infoErr := dirEntry.Info()
			if infoErr != nil {
				skippedFiles++
				return nil
			}

			fileSize := fileInfo.Size()
			totalSize += fileSize
			_ = bar.Add64(fileSize)
		}

		return nil
	})

	if walkErr != nil {
		switch {
		case os.IsPermission(walkErr):
			return 0, fmt.Errorf("permission denied for: path %s", path)
		case os.IsNotExist(walkErr):
			return 0, fmt.Errorf("file not found: path %s", path)
		default:
			return 0, fmt.Errorf("failed to walk directory path: %v", walkErr)
		}
	}

	return totalSize, nil
}

var (
	units      = []string{"B", "KB", "MB", "GB", "TB", "PB"}
	thresholds = []float64{
		1,
		1024,
		1024 * 1024,
		1024 * 1024 * 1024,
		1024 * 1024 * 1024 * 1024,
		1024 * 1024 * 1024 * 1024 * 1024,
	}
)

func humanReadableSize(size int64, fn int) string {
	if size == 0 {
		return "0 B"
	}

	sizeFloat := float64(size)

	unitIndex := 0
	for i := len(thresholds) - 1; i > 0; i-- {
		if sizeFloat >= thresholds[i] {
			unitIndex = i
			sizeFloat /= thresholds[i]
			break
		}
	}

	var formatted string
	if sizeFloat < 10 && unitIndex > 0 {
		decimals := fn
		if decimals < 1 {
			decimals = 1
		}
		formatted = fmt.Sprintf("%."+fmt.Sprintf("%d", decimals)+"f", sizeFloat)
		formatted = strings.TrimSuffix(formatted, ".0")
	} else {
		formatted = fmt.Sprintf("%.0f", sizeFloat)
	}

	return fmt.Sprintf("%s %s", formatted, units[unitIndex])
}

func printSizeTable(its items, cl *color.GlobalColor, config SizeConfig) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	if config.TableStyle != "none" {
		t.AppendHeader(table.Row{"Size", "Name"})
	}

	for i := range its {
		colorSize := cl.SWhite(its[i].Size)
		colorName := utils.SprintStringColor(its[i].Name, its[i].Name, cl)
		t.AppendRow(table.Row{colorSize, colorName})
	}

	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "Size", Align: text.AlignRight},
		{Name: "Name", Align: text.AlignLeft},
	})

	if style, ok := types.TableStyleMap[config.TableStyle]; ok {
		t.SetStyle(style)
	}

	t.Render()
}
