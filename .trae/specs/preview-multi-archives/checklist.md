# Checklist

- [x] CLI 层调用 `fs.ExpandFiles` 展开通配符
- [x] 命令层 `PreviewConfig.PackPaths` 为 `[]string` 类型
- [x] 单个压缩包时行为与之前完全一致（无标题）
- [x] 多个压缩包时在输出间打印 `==> filename <==` 分隔
- [x] 无参数时返回 `"missing archive path"` 错误
- [x] 展开后列表为空时返回 `"no archive files found"` 错误
- [x] `UsageSyntax` 和 `Examples` 更新反映多文件支持
- [x] 编译通过
- [x] 通配符无匹配项时原样传递（`fs.ExpandFiles` 行为）
