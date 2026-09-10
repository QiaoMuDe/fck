package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/doc2md"
	"gitee.com/MM-Q/qflag"
)

// Doc2mdCmd doc2md 命令
var Doc2mdCmd *qflag.Cmd

var (
	doc2mdOutput   *qflag.StringFlag // -o, --output 输出 Markdown 文件路径
	doc2mdExt      *qflag.EnumFlag   // -x, --extension stdin 管道输入时的扩展名提示（枚举限制）
	doc2mdMime     *qflag.StringFlag // -m, --mime-type MIME 类型提示
	doc2mdCharset  *qflag.StringFlag // -c, --charset 字符集提示
	doc2mdKeepData *qflag.BoolFlag   // --keep-data-uris 保留完整 base64 图片数据 URI
)

// supportedDocExtensions 是 -x/--extension 支持的文档扩展名枚举列表。
// 首项 DocExtNone 为哨兵默认值，用于校验管道输入时未显式指定扩展名的情况。
var supportedDocExtensions = []string{doc2md.DocExtNone, ".csv", ".docx", ".epub", ".html", ".ipynb", ".md", ".pdf", ".pptx", ".rss", ".txt", ".xls", ".xlsx", ".zip"}

func init() {
	Doc2mdCmd = qflag.NewCmd("doc2md", "", qflag.ExitOnError)

	doc2mdOutput = Doc2mdCmd.String("output", "o", "输出 Markdown 文件路径（默认输出到终端）", "")
	doc2mdExt = Doc2mdCmd.Enum("extension", "x", qflag.EnumHelp("输入文档扩展名（stdin 管道输入时必填，none 表示未指定），支持:", supportedDocExtensions, "\t\t\t\t\t"), doc2md.DocExtNone, supportedDocExtensions)
	doc2mdMime = Doc2mdCmd.String("mime-type", "m", "输入文档 MIME 类型提示", "")
	doc2mdCharset = Doc2mdCmd.String("charset", "c", "输入文档字符集提示（如: gbk）", "")
	doc2mdKeepData = Doc2mdCmd.Bool("keep-data-uris", "", "保留完整 base64 图片数据 URI（默认截断）", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "将办公文档（Word/Excel/PPT/PDF 等）转换为 Markdown",
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s doc2md [options] [files]...", qflag.Root.Name()),
		Examples: map[string]string{
			"转换 Word 文档":  fmt.Sprintf("%s doc2md input.docx", qflag.Root.Name()),
			"输出到文件":       fmt.Sprintf("%s doc2md -o out.md input.docx", qflag.Root.Name()),
			"管道输入":        fmt.Sprintf("cat input.docx | %s doc2md -x .docx", qflag.Root.Name()),
			"保留完整图片数据URI": fmt.Sprintf("%s doc2md --keep-data-uris input.docx", qflag.Root.Name()),
		},
		Notes: []string{
			"stdin 管道输入时必须使用 -x/--extension 提示文档扩展名",
			"支持文件路径通配符展开",
			"默认输出到终端",
			"指定 -o/--output 输出文件时一次只能转换一个输入文件",
		},
	}

	if err := Doc2mdCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	Doc2mdCmd.SetRun(runDoc2md)
}

// runDoc2md 运行 doc2md 命令，将参数映射为配置并调用业务实现
//
// 参数:
//   - cmd: 命令接口
//
// 返回值:
//   - error: 错误
func runDoc2md(cmd qflag.Command) error {
	config := doc2md.Doc2mdConfig{
		Files:        cmd.Args(),
		Extension:    doc2mdExt.Get(),
		MIMEType:     doc2mdMime.Get(),
		Charset:      doc2mdCharset.Get(),
		Output:       doc2mdOutput.Get(),
		KeepDataURIs: doc2mdKeepData.Get(),
	}

	return doc2md.Doc2mdCmdMain(config)
}
