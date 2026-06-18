# Checklist

## 业务逻辑
- [x] `ShfmtConfig` 包含 `FormatOptions shx.FormatOptions` 字段
- [x] 默认使用 `shx.DefaultFormatOptions()` 保持一致行为（CLI 层标志默认值）
- [x] `processStdin` 和 `processFile` 使用 `shx.FormatWithOptions()`

## CLI 标志
- [x] `--indent` / `-i` 正确解析 uint 值并传递
- [x] 所有 bool 标志正确解析并传递

## 功能验证
- [ ] `fck shfmt -i 0 script.sh` 使用 tab 缩进格式化
- [ ] `fck shfmt -i 4 script.sh` 使用 4 空格缩进
- [ ] `fck shfmt --switch-case-indent=false script.sh` 禁用 case 缩进
- [ ] `fck shfmt --keep-comments=false script.sh` 移除注释
- [ ] `fck shfmt --binary-next-line script.sh` 操作符换行
- [ ] `fck shfmt --function-next-line script.sh` 函数体换行
- [ ] `fck shfmt --space-redirects script.sh` 重定向符加空格
- [ ] `fck shfmt --single-line script.sh` 单行输出
- [ ] `fck shfmt -m script.sh` 最小化输出
- [ ] `fck shfmt -i 4 -m -w script.sh` 多标志组合正常

## 回归
- [ ] 不指定任何格式化标志时，行为与旧版 `shx.Format()` 一致
- [ ] `fck shfmt -w script.sh` 原地写入仍正常工作
- [ ] `fck shfmt -w -b script.sh` 备份功能仍正常工作
- [ ] `cat script.sh | fck shfmt` 管道输入仍正常工作
- [ ] `fck shfmt -w *.sh` 通配符批量仍正常工作

## 编译与规范
- [x] `go build ./...` 编译通过
- [x] `go vet ./...` 无错误
- [x] `golangci-lint` 无新增错误
