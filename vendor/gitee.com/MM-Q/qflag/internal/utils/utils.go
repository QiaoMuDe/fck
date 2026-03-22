package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"gitee.com/MM-Q/qflag/internal/types"
)

// GetCmdName 返回命令名称
//
// 参数:
//   - cmd: 要获取名称的命令
//
// 返回值:
//   - string: 命令名称
func GetCmdName(cmd types.Command) string {
	// 如果命令有长名和短名, 返回长名和短名
	if cmd.LongName() != "" && cmd.ShortName() != "" {
		return fmt.Sprintf("%s, %s\n", cmd.LongName(), cmd.ShortName())
	}

	// 如果命令有名称, 返回名称
	if cmd.Name() != "" {
		return cmd.Name() + "\n"
	}

	// 如果命令没有名称, 返回程序文件名
	return filepath.Base(os.Args[0]) + "\n"
}

// FormatDefaultValue 返回格式化后的默认值
//
// 参数:
//   - flagType: 标志类型
//   - defValue: 默认值
//
// 返回值:
//   - string: 格式化后的默认值
func FormatDefaultValue(flagType types.FlagType, defValue any) string {
	if defValue == nil {
		return ""
	}

	// 根据标志类型进行专门处理
	switch flagType {
	case types.FlagTypeString, types.FlagTypeEnum:
		if v, ok := defValue.(string); ok {
			if v == "" {
				return `""`
			}
			return v
		}
		return fmt.Sprintf("%v", defValue)

	case types.FlagTypeInt, types.FlagTypeInt64, types.FlagTypeUint, types.FlagTypeUint8, types.FlagTypeUint16, types.FlagTypeUint32, types.FlagTypeUint64:
		return fmt.Sprintf("%d", defValue)

	case types.FlagTypeFloat64:
		if v, ok := defValue.(float64); ok {
			return fmt.Sprintf("%.2f", v)
		}
		return fmt.Sprintf("%v", defValue)

	case types.FlagTypeBool:
		if v, ok := defValue.(bool); ok {
			return strconv.FormatBool(v)
		}
		return fmt.Sprintf("%v", defValue)

	case types.FlagTypeDuration:
		if v, ok := defValue.(time.Duration); ok {
			return v.String()
		}
		return fmt.Sprintf("%v", defValue)

	case types.FlagTypeTime:
		if v, ok := defValue.(time.Time); ok {
			if v.IsZero() {
				return `""`
			}
			return v.Format("2006-01-02 15:04:05")
		}
		return fmt.Sprintf("%v", defValue)

	case types.FlagTypeSize:
		if v, ok := defValue.(int64); ok {
			return FormatSize(v)
		}
		return fmt.Sprintf("%d bytes", defValue)

	case types.FlagTypeMap:
		if v, ok := defValue.(map[string]string); ok {
			if len(v) == 0 {
				return "{}"
			}
			// 将map格式化为更易读的字符串
			pairs := make([]string, 0, len(v))
			for key, value := range v {
				pairs = append(pairs, fmt.Sprintf("%s=%s", key, value))
			}
			return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
		}
		return fmt.Sprintf("%v", defValue)

	case types.FlagTypeStringSlice:
		return fmt.Sprintf("%v", defValue)

	case types.FlagTypeIntSlice:
		return fmt.Sprintf("%v", defValue)

	case types.FlagTypeInt64Slice:
		return fmt.Sprintf("%v", defValue)

	default:
		// 未知类型, 使用通用格式化
		formatted := fmt.Sprintf("%v", defValue)

		// 对于空切片或空映射, 提供更友好的显示
		if formatted == "[]" || formatted == "map[]" {
			return formatted
		}

		// 对于空字符串, 确保显示为 ""
		if formatted == "" {
			return `""`
		}

		return formatted
	}
}

// CalcOptionMaxWidth 计算选项名称最大宽度
//
// 参数:
//   - options: 选项信息列表
//
// 返回值:
//   - int: 选项名称最大宽度
func CalcOptionMaxWidth(options []types.OptionInfo) int {
	maxWidth := 0
	for _, opt := range options {
		if len(opt.NamePart) > maxWidth {
			maxWidth = len(opt.NamePart)
		}
	}
	return maxWidth
}

// CalcSubCmdMaxLen 计算子命令名称最大宽度
//
// 参数:
//   - subCmds: 子命令信息列表
//
// 返回值:
//   - int: 子命令名称最大宽度
func CalcSubCmdMaxLen(subCmds []types.SubCmdInfo) int {
	maxLen := 0
	for _, info := range subCmds {
		if len(info.Name) > maxLen {
			maxLen = len(info.Name)
		}
	}
	return maxLen
}

// getOptionGroup 获取选项的分组优先级
//
// 参数:
//   - namePart: 选项名称部分
//
// 返回值:
//   - int: 分组优先级 (0: 长短都有, 1: 仅长选项, 2: 仅短选项)
func getOptionGroup(namePart string) int {
	if strings.Contains(namePart, ", --") {
		return 0
	}
	if strings.HasPrefix(namePart, "--") {
		return 1
	}
	return 2
}

// getSubCmdGroup 获取子命令的分组优先级
//
// 参数:
//   - name: 子命令名称
//
// 返回值:
//   - int: 分组优先级 (0: 长短都有, 1: 仅单一名字)
func getSubCmdGroup(name string) int {
	if strings.Contains(name, ", ") {
		return 0
	}
	return 1
}

// SortOptions 对选项列表进行排序
// 排序规则: 长短都有 > 仅长选项 > 仅短选项, 每组内按首字母排序 (忽略大小写)
//
// 参数:
//   - options: 选项信息列表
func SortOptions(options []types.OptionInfo) {
	sort.Slice(options, func(i, j int) bool {
		iGroup := getOptionGroup(options[i].NamePart)
		jGroup := getOptionGroup(options[j].NamePart)
		if iGroup != jGroup {
			return iGroup < jGroup
		}
		// 使用 strings.ToLower 进行大小写不敏感的字符串比较
		return strings.ToLower(options[i].NamePart) < strings.ToLower(options[j].NamePart)
	})
}

// SortSubCmds 对子命令列表进行排序
// 排序规则: 长短名都有 > 仅单一名字, 每组内按首字母排序 (忽略大小写)
//
// 参数:
//   - subCmds: 子命令信息列表
func SortSubCmds(subCmds []types.SubCmdInfo) {
	sort.Slice(subCmds, func(i, j int) bool {
		iGroup := getSubCmdGroup(subCmds[i].Name)
		jGroup := getSubCmdGroup(subCmds[j].Name)
		if iGroup != jGroup {
			return iGroup < jGroup
		}
		// 使用 strings.ToLower 进行大小写不敏感的字符串比较
		return strings.ToLower(subCmds[i].Name) < strings.ToLower(subCmds[j].Name)
	})
}

// ValidateFlagName 验证标志名称是否可用
//
// 参数:
//   - cmd: 命令对象, 用于检查标志注册表
//   - longName: 长标志名称
//   - shortName: 短标志名称
//
// 返回值:
//   - error: 如果验证失败返回错误, 成功返回 nil
func ValidateFlagName(cmd types.Command, longName, shortName string) error {
	if cmd == nil {
		return types.NewError("INVALID_COMMAND", "cmd cannot be nil", nil)
	}

	cmdName := cmd.Name()

	if longName == "" && shortName == "" {
		return types.NewError("INVALID_FLAG_NAME", fmt.Sprintf("cmd %q: flag name cannot be empty", cmdName), nil)
	}

	// 分别检查长名称和短名称
	if longName != "" {
		if _, ok := cmd.GetFlag(longName); ok {
			return types.NewError("FLAG_ALREADY_EXISTS", fmt.Sprintf("cmd %q: flag '--%s' already exists", cmdName, longName), nil)
		}
	}

	if shortName != "" {
		if _, ok := cmd.GetFlag(shortName); ok {
			return types.NewError("FLAG_ALREADY_EXISTS", fmt.Sprintf("cmd %q: flag '-%s' already exists", cmdName, shortName), nil)
		}
	}

	return nil
}

// FormatFlagName 格式化标志名称
//
// 参数:
//   - longName: 长标志名称
//   - shortName: 短标志名称
//
// 返回值:
//   - string: 格式化后的标志名称
func FormatFlagName(longName, shortName string) string {
	if longName != "" && shortName != "" {
		return fmt.Sprintf("-%s, --%s", shortName, longName)
	}
	if longName != "" {
		return fmt.Sprintf("--%s", longName)
	}
	return fmt.Sprintf("-%s", shortName)
}

// ToStrSlice 将任何类型的切片转换为字符串切片
// 主要用于将数字类型切片转换为字符串切片, 以便用于 EnumFlag
//
// 参数:
//   - slice: 任意类型的切片
//
// 返回值:
//   - []string: 转换后的字符串切片
//   - error: 如果转换失败返回错误
func ToStrSlice(slice any) ([]string, error) {
	if slice == nil {
		return nil, types.NewError("NIL_SLICE", "slice cannot be nil", nil)
	}

	// 使用反射获取切片的值和类型
	val := reflect.ValueOf(slice)
	if val.Kind() != reflect.Slice {
		return nil, types.NewError("NOT_SLICE", "input must be a slice", nil)
	}

	result := make([]string, 0, val.Len())
	for i := 0; i < val.Len(); i++ {
		elem := val.Index(i)

		// 根据元素类型进行转换
		switch elem.Kind() {
		case reflect.String:
			result = append(result, elem.String())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			result = append(result, fmt.Sprintf("%d", elem.Int()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			result = append(result, fmt.Sprintf("%d", elem.Uint()))
		case reflect.Float32, reflect.Float64:
			result = append(result, fmt.Sprintf("%g", elem.Float()))
		case reflect.Bool:
			result = append(result, fmt.Sprintf("%t", elem.Bool()))
		default:
			// 对于其他类型, 使用 fmt.Sprintf 尝试转换
			result = append(result, fmt.Sprintf("%v", elem.Interface()))
		}
	}

	return result, nil
}

// FormatSize 格式化大小值为带单位的字符串
//
// 参数:
//   - size: 大小值 (字节)
//
// 返回值:
//   - string: 格式化后的大小字符串, 如 "10KB", "4MB"
//
// 功能说明:
//   - 自动选择合适的单位 (B, KB, MB, GB, TB)
//   - 保留两位小数
//   - 对于小于1KB的值, 直接显示字节数
func FormatSize(size int64) string {
	if size < 0 {
		return "0B"
	}

	// 根据大小选择合适的单位 (KB, MB, GB 等使用 1000 为基数)
	switch {
	case size < types.KB:
		return fmt.Sprintf("%dB", size)
	case size < types.MB:
		return fmt.Sprintf("%.2fKB", float64(size)/float64(types.KB))
	case size < types.GB:
		return fmt.Sprintf("%.2fMB", float64(size)/float64(types.MB))
	case size < types.TB:
		return fmt.Sprintf("%.2fGB", float64(size)/float64(types.GB))
	case size < types.PB:
		return fmt.Sprintf("%.2fTB", float64(size)/float64(types.TB))
	default:
		return fmt.Sprintf("%.2fPB", float64(size)/float64(types.PB))
	}
}
