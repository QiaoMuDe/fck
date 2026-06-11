package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/json"
	"gitee.com/MM-Q/qflag"
)

var JsonCmd *qflag.Cmd

var (
	jsonPretty    *qflag.BoolFlag        // -p, --pretty      美化输出
	jsonCompact   *qflag.BoolFlag        // -c, --compact     压缩输出
	jsonValidate  *qflag.BoolFlag        // -v, --validate    验证模式
	jsonQuery     *qflag.StringFlag      // -q, --query       查询路径
	jsonHighlight *qflag.BoolFlag        // -H, --highlight   语法高亮
	jsonRaw       *qflag.BoolFlag        // -r, --raw         原始字符串输出
	jsonWrite     *qflag.BoolFlag        // -w, --write       原地写入
	jsonBackup    *qflag.BoolFlag        // -b, --backup      写入前备份
	jsonSet       *qflag.StringSliceFlag // -s, --set         设置指定路径的值
	jsonDelete    *qflag.StringSliceFlag // -D, --delete      删除指定路径
	jsonType      *qflag.EnumFlag        // -t, --type        值类型 (auto/string/number/bool)
)

func init() {
	JsonCmd = qflag.NewCmd("json", "j", qflag.ExitOnError)

	jsonPretty = JsonCmd.Bool("pretty", "p", "格式化JSON", false)
	jsonCompact = JsonCmd.Bool("compact", "c", "压缩JSON", false)
	jsonValidate = JsonCmd.Bool("validate", "v", "验证 JSON 语法有效性", false)
	jsonQuery = JsonCmd.String("query", "q", "使用路径表达式查询数据，如: users.0.name", "")
	jsonHighlight = JsonCmd.Bool("highlight", "H", "语法高亮显示", false)
	jsonRaw = JsonCmd.Bool("raw", "r", "原始字符串输出 (不带引号)", false)
	jsonWrite = JsonCmd.Bool("write", "w", "原地写入处理后的 JSON 到原文件", false)
	jsonBackup = JsonCmd.Bool("backup", "b", "写入前创建备份文件（.bak后缀）", false)
	jsonSet = JsonCmd.StringSlice("set", "s", "设置指定路径的值（多个值以逗号分隔）", nil)
	jsonDelete = JsonCmd.StringSlice("delete", "D", "删除指定路径（多个路径以逗号分隔）", nil)
	jsonType = JsonCmd.Enum("type", "t", "值类型 (auto/string/number/bool)", "auto", []string{"auto", "string", "number", "bool"})

	cmdOpts := &qflag.CmdOpts{
		Desc:        "JSON 数据处理工具",
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s json [options] [file]", qflag.Root.Name()),
		Notes: []string{
			"输入方式: 管道传递JSON字符串 或 位置参数指定文件路径",
			"查询路径使用点号分隔，如: users.0.name",
			"数组索引支持负数，-1 表示最后一个元素",
			"支持通配符 * 匹配数组所有元素，如: users.*.name",
			"使用 -w 原地写入时，必须指定文件路径参数（不支持管道）",
			"使用 -b 备份时，会创建 .bak 后缀的备份文件",
			"-s 设置值时，路径与值用 = 分隔，如: name=Tom",
			"-t 指定值类型: auto 自动推断、string 字符串、number 数值、bool 布尔",
			"-D 删除路径后，源文件内容将被修改（配合 -w 原地写入）",
		},
		Examples: map[string]string{
			"格式化 JSON (管道)": fmt.Sprintf("echo '{\"a\":1}' | %s json -p", qflag.Root.Name()),
			"格式化 JSON (文件)": fmt.Sprintf("%s json -p data.json", qflag.Root.Name()),
			"压缩 JSON":       fmt.Sprintf("%s json -c file.json", qflag.Root.Name()),
			"验证 JSON":       fmt.Sprintf("%s json -v invalid.json", qflag.Root.Name()),
			"查询数据 (管道)":     fmt.Sprintf("echo '{\"users\":[{\"name\":\"Tom\"}]}' | %s json -q users.0.name", qflag.Root.Name()),
			"查询数据 (文件)":     fmt.Sprintf("%s json -q users.0.name data.json", qflag.Root.Name()),
			"高亮显示":          fmt.Sprintf("%s json -pH large.json", qflag.Root.Name()),
			"提取数组所有元素":      fmt.Sprintf("echo '{\"items\":[1,2,3]}' | %s json -q items.*", qflag.Root.Name()),
			"原地格式化文件":       fmt.Sprintf("%s json -p -w data.json", qflag.Root.Name()),
			"原地压缩并备份":       fmt.Sprintf("%s json -c -w -b config.json", qflag.Root.Name()),
			"设置字段值 (字符串)":   fmt.Sprintf("echo '{\"name\":\"old\"}' | %s json -s name=new", qflag.Root.Name()),
			"设置字段值 (数字)":    fmt.Sprintf("echo '{\"count\":0}' | %s json -s count=42 -t number", qflag.Root.Name()),
			"设置嵌套字段值":       fmt.Sprintf("echo '{\"user\":{\"name\":\"Tom\"}}' | %s json -s user.name=Jerry", qflag.Root.Name()),
			"删除字段":          fmt.Sprintf("echo '{\"name\":\"Tom\",\"age\":20}' | %s json -D age", qflag.Root.Name()),
			"删除数组元素":        fmt.Sprintf("echo '{\"items\":[1,2,3]}' | %s json -D items.1", qflag.Root.Name()),
		},
		MutexGroups: []qflag.MutexGroup{
			{
				// 操作模式互斥
				Name:      "operation-mode",
				Flags:     []string{"validate", "query", "set", "delete"},
				AllowNone: true, // 允许都不选，默认使用 pretty 模式
			},
			{
				// 写入模式与查询和验证模式互斥
				Name:      "write-query-validate-mutex",
				Flags:     []string{"write", "query", "validate"},
				AllowNone: true,
			},
			{
				// 备份模式与查询验证模式互斥
				Name:      "backup-query-validate-mutex",
				Flags:     []string{"backup", "query", "validate"},
				AllowNone: true,
			},
		},
		FlagDependencies: []qflag.FlagDependency{
			{
				Name:    "backup-write",
				Trigger: "backup",
				Targets: []string{"write"},
				Type:    qflag.DepRequired,
			},
		},
	}

	if err := JsonCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	JsonCmd.SetRun(runJson)
}

func runJson(cmd qflag.Command) error {
	// 检查写入模式必须与格式化、压缩、设置或删除一起使用
	if jsonWrite.Get() && !jsonPretty.Get() && !jsonCompact.Get() && len(jsonSet.Get()) == 0 && len(jsonDelete.Get()) == 0 {
		return fmt.Errorf("-w flag must be used with -p (pretty), -c (compact), -s (set), or -D (delete)")
	}

	// 如果启用原地写入，禁用高亮模式
	highlight := jsonHighlight.Get()
	if jsonWrite.Get() {
		highlight = false
	}

	config := json.JsonConfig{
		Pretty:     jsonPretty.Get(),
		Compact:    jsonCompact.Get(),
		Validate:   jsonValidate.Get(),
		Query:      jsonQuery.Get(),
		Highlight:  highlight, // 高亮模式
		Raw:        jsonRaw.Get(),
		Write:      jsonWrite.Get(),
		Backup:     jsonBackup.Get(),
		SetValue:   jsonSet.Get(),
		DeletePath: jsonDelete.Get(),
		ValueType:  jsonType.Get(),
		Files:      cmd.Args(),
	}

	return json.JsonCmdMain(config)
}
