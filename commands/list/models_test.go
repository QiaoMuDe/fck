package list

import (
	"testing"
	"time"

	"gitee.com/MM-Q/fck/commands/internal/types"
)

func TestScanOptions(t *testing.T) {
	opts := ScanOptions{
		Recursive:  true,
		ShowHidden: false,
		FileTypes:  []string{types.FindTypeFile},
		DirItself:  false,
	}

	if !opts.Recursive {
		t.Error("ScanOptions.Recursive 设置失败")
	}

	if opts.ShowHidden {
		t.Error("ScanOptions.ShowHidden 应该为 false")
	}

	if len(opts.FileTypes) != 1 {
		t.Errorf("ScanOptions.FileTypes 长度 = %v, 期望 1", len(opts.FileTypes))
	}

	if opts.FileTypes[0] != types.FindTypeFile {
		t.Errorf("ScanOptions.FileTypes[0] = %v, 期望 %v", opts.FileTypes[0], types.FindTypeFile)
	}
}

func TestProcessOptions(t *testing.T) {
	opts := ProcessOptions{
		SortBy:     "name",
		Reverse:    true,
		GroupByDir: false,
	}

	if opts.SortBy != "name" {
		t.Errorf("ProcessOptions.SortBy = %v, 期望 name", opts.SortBy)
	}

	if !opts.Reverse {
		t.Error("ProcessOptions.Reverse 应该为 true")
	}

	if opts.GroupByDir {
		t.Error("ProcessOptions.GroupByDir 应该为 false")
	}
}

func TestFormatOptions(t *testing.T) {
	opts := FormatOptions{
		LongFormat:    true,
		UseColor:      false,
		TableStyle:    "default",
		QuoteNames:    true,
		ShowUserGroup: false,
	}

	if !opts.LongFormat {
		t.Error("FormatOptions.LongFormat 应该为 true")
	}

	if opts.UseColor {
		t.Error("FormatOptions.UseColor 应该为 false")
	}

	if opts.TableStyle != "default" {
		t.Errorf("FormatOptions.TableStyle = %v, 期望 default", opts.TableStyle)
	}

	if !opts.QuoteNames {
		t.Error("FormatOptions.QuoteNames 应该为 true")
	}
}

func TestFileInfo(t *testing.T) {
	now := time.Now()
	fileInfo := FileInfo{
		Name:           "test.txt",
		Path:           "/path/to/test.txt",
		EntryType:      FileType,
		Size:           1024,
		ModTime:        now,
		Perm:           "-rw-r--r--",
		Owner:          "user",
		Group:          "group",
		FileExt:        ".txt",
		LinkTargetPath: "",
	}

	if fileInfo.Name != "test.txt" {
		t.Errorf("FileInfo.Name = %v, 期望 test.txt", fileInfo.Name)
	}

	if fileInfo.Path != "/path/to/test.txt" {
		t.Errorf("FileInfo.Path = %v, 期望 /path/to/test.txt", fileInfo.Path)
	}

	if fileInfo.EntryType != FileType {
		t.Errorf("FileInfo.EntryType = %v, 期望 %v", fileInfo.EntryType, FileType)
	}

	if fileInfo.Size != 1024 {
		t.Errorf("FileInfo.Size = %v, 期望 1024", fileInfo.Size)
	}

	if !fileInfo.ModTime.Equal(now) {
		t.Errorf("FileInfo.ModTime = %v, 期望 %v", fileInfo.ModTime, now)
	}

	if fileInfo.Perm != "-rw-r--r--" {
		t.Errorf("FileInfo.Perm = %v, 期望 -rw-r--r--", fileInfo.Perm)
	}

	if fileInfo.Owner != "user" {
		t.Errorf("FileInfo.Owner = %v, 期望 user", fileInfo.Owner)
	}

	if fileInfo.Group != "group" {
		t.Errorf("FileInfo.Group = %v, 期望 group", fileInfo.Group)
	}

	if fileInfo.FileExt != ".txt" {
		t.Errorf("FileInfo.FileExt = %v, 期望 .txt", fileInfo.FileExt)
	}

	if fileInfo.LinkTargetPath != "" {
		t.Errorf("FileInfo.LinkTargetPath = %v, 期望空字符串", fileInfo.LinkTargetPath)
	}
}

func TestFileInfoList(t *testing.T) {
	now := time.Now()
	files := FileInfoList{
		{
			Name:      "file1.txt",
			EntryType: FileType,
			Size:      100,
			ModTime:   now,
		},
		{
			Name:      "file2.txt",
			EntryType: FileType,
			Size:      200,
			ModTime:   now.Add(time.Hour),
		},
	}

	if len(files) != 2 {
		t.Errorf("FileInfoList 长度 = %v, 期望 2", len(files))
	}

	if files[0].Name != "file1.txt" {
		t.Errorf("FileInfoList[0].Name = %v, 期望 file1.txt", files[0].Name)
	}

	if files[1].Size != 200 {
		t.Errorf("FileInfoList[1].Size = %v, 期望 200", files[1].Size)
	}
}

func TestPermissionColorMap(t *testing.T) {
	// 测试权限颜色映射是否正确初始化
	expectedColors := map[int]colorType{
		1: colorTypeGreen,  // 所有者-读-绿色
		2: colorTypeYellow, // 所有者-写-黄色
		3: colorTypeRed,    // 所有者-执行-红色
		4: colorTypeGreen,  // 组-读-绿色
		5: colorTypeYellow, // 组-写-黄色
		6: colorTypeRed,    // 组-执行-红色
		7: colorTypeGreen,  // 其他-读-绿色
		8: colorTypeYellow, // 其他-写-黄色
		9: colorTypeRed,    // 其他-执行-红色
	}

	for pos, expectedColor := range expectedColors {
		if color, exists := permissionColorMap[pos]; !exists {
			t.Errorf("permissionColorMap[%d] 不存在", pos)
		} else if color != expectedColor {
			t.Errorf("permissionColorMap[%d] = %v, 期望 %v", pos, color, expectedColor)
		}
	}

	// 验证映射表的完整性
	if len(permissionColorMap) != 9 {
		t.Errorf("permissionColorMap 长度 = %v, 期望 9", len(permissionColorMap))
	}
}

// 测试结构体的零值
func TestStructZeroValues(t *testing.T) {
	var scanOpts ScanOptions
	if scanOpts.Recursive != false {
		t.Error("ScanOptions 零值 Recursive 应该为 false")
	}
	if scanOpts.ShowHidden != false {
		t.Error("ScanOptions 零值 ShowHidden 应该为 false")
	}
	if scanOpts.FileTypes != nil {
		t.Error("ScanOptions 零值 FileTypes 应该为 nil")
	}

	var processOpts ProcessOptions
	if processOpts.SortBy != "" {
		t.Error("ProcessOptions 零值 SortBy 应该为空字符串")
	}
	if processOpts.Reverse != false {
		t.Error("ProcessOptions 零值 Reverse 应该为 false")
	}

	var formatOpts FormatOptions
	if formatOpts.LongFormat != false {
		t.Error("FormatOptions 零值 LongFormat 应该为 false")
	}
	if formatOpts.UseColor != false {
		t.Error("FormatOptions 零值 UseColor 应该为 false")
	}

	var fileInfo FileInfo
	if fileInfo.Name != "" {
		t.Error("FileInfo 零值 Name 应该为空字符串")
	}
	if fileInfo.Size != 0 {
		t.Error("FileInfo 零值 Size 应该为 0")
	}
}

// 测试结构体的复制
func TestStructCopy(t *testing.T) {
	original := FileInfo{
		Name:      "original.txt",
		Size:      1024,
		EntryType: FileType,
	}

	// 值复制
	copied := original
	copied.Name = "copied.txt"

	if original.Name != "original.txt" {
		t.Error("原始结构体不应该被修改")
	}

	if copied.Name != "copied.txt" {
		t.Error("复制的结构体应该被修改")
	}

	// 验证其他字段保持一致
	if copied.Size != original.Size {
		t.Error("复制的结构体其他字段应该保持一致")
	}

	if copied.EntryType != original.EntryType {
		t.Error("复制的结构体其他字段应该保持一致")
	}
}

// 测试切片操作
func TestFileInfoListOperations(t *testing.T) {
	files := FileInfoList{}

	// 测试追加
	file1 := FileInfo{Name: "file1.txt"}
	files = append(files, file1)

	if len(files) != 1 {
		t.Errorf("追加后长度 = %v, 期望 1", len(files))
	}

	// 测试批量追加
	file2 := FileInfo{Name: "file2.txt"}
	file3 := FileInfo{Name: "file3.txt"}
	files = append(files, file2, file3)

	if len(files) != 3 {
		t.Errorf("批量追加后长度 = %v, 期望 3", len(files))
	}

	// 测试索引访问
	if files[0].Name != "file1.txt" {
		t.Errorf("files[0].Name = %v, 期望 file1.txt", files[0].Name)
	}

	if files[2].Name != "file3.txt" {
		t.Errorf("files[2].Name = %v, 期望 file3.txt", files[2].Name)
	}

	// 测试切片
	subset := files[1:3]
	if len(subset) != 2 {
		t.Errorf("切片长度 = %v, 期望 2", len(subset))
	}

	if subset[0].Name != "file2.txt" {
		t.Errorf("subset[0].Name = %v, 期望 file2.txt", subset[0].Name)
	}
}
