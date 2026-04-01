// Package size 实现了文件和目录大小计算功能。
// 该文件提供了计算文件、目录大小的核心功能，支持通配符路径展开、隐藏文件处理、进度显示和表格输出。
package size

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/internal/types"
	common "gitee.com/MM-Q/fck/internal/utils"
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
	Args       []string
	Color      bool
	TableStyle string
	Hidden     bool
	Human      bool
}

func SizeCmdMain(cl *colorlib.ColorLib, config SizeConfig) error {
	targetPaths := config.Args

	if len(targetPaths) == 0 {
		targetPaths = []string{"*"}
	}

	cl.SetColor(config.Color)

	var itemList items

	for _, targetPath := range targetPaths {
		targetPath = filepath.Clean(targetPath)

		pathsToProcess, err := expandPath(targetPath)
		if err != nil {
			cl.PrintErrorf("failed to expand path: %v\n", err)
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

func addPathToList(path string, itemList *items, cl *colorlib.ColorLib, config SizeConfig) {
	if !config.Hidden && common.IsHidden(path) {
		return
	}

	size, err := getPathSize(path, config)
	if err != nil {
		cl.PrintErrorf("failed to calculate size of: %s - %v\n", path, err)
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
	if !includeHidden && common.IsHidden(path) {
		return 0, nil
	}

	info, err := os.Lstat(path)
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

		if !includeHidden && common.IsHidden(filePath) {
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

func printSizeTable(its items, cl *colorlib.ColorLib, config SizeConfig) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	if config.TableStyle != "none" {
		t.AppendHeader(table.Row{"Size", "Name"})
	}

	for i := range its {
		colorSize := cl.Swhite(its[i].Size)
		colorName := common.SprintStringColor(its[i].Name, its[i].Name, cl)
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
