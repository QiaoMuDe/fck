# Checklist

## shfmt 命令
- [x] `fck shfmt script.sh` 格式化脚本并输出到 stdout
- [x] `cat script.sh | fck shfmt` 从管道读取并格式化
- [x] `fck shfmt -w script.sh` 原地写入格式化内容
- [x] `fck shfmt -w -b script.sh` 原地写入并创建备份
- [x] `fck shfmt -w *.sh` 通配符匹配批量格式化多个文件
- [x] `fck shfmt -w **/*.sh` 递归匹配子目录中的 .sh 文件

## shck 命令
- [x] `fck shck valid.sh` 语法正确时输出 `✓ syntax is valid`
- [x] `fck shck invalid.sh` 语法错误时显示行号、列号和错误描述
- [x] `cat script.sh | fck shck` 从管道读取并检查语法
- [x] `fck shck src/*.sh` 通配符匹配批量检查多个文件

## shx 命令
- [x] `fck shx "echo hello"` 执行命令字符串
- [x] `fck shx script.sh` 执行 .sh 脚本文件
- [x] `fck shx -t 5s "sleep 10"` 超时控制
- [x] `fck shx -d /tmp "pwd"` 指定工作目录
- [x] `fck shx -e FOO=bar "echo $FOO"` 设置环境变量
- [x] `echo "hello" | fck shx "cat"` 管道 stdin 传递

## 注册与编译
- [x] 3 个新命令在 `fck --help` 列表中可见
- [x] `fck shfmt --help`、`fck shck --help`、`fck shx --help` 显示正确帮助信息
- [x] 项目编译通过，无错误
