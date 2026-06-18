# Checklist

## shck -q 静默模式
- [x] `fck shck -q valid.sh` 不输出内容，退出码 0
- [x] `fck shck -q invalid.sh` 仅输出错误信息，退出码非 0
- [x] `fck shck -q src/*.sh`（部分错误）仅输出失败文件错误信息
- [x] `cat valid.sh | fck shck -q` 管道输入正确，不输出内容
- [x] `fck shck --quiet valid.sh` 长选项同样生效
- [x] 未指定 `-q` 时行为不变，仍输出 ✓ 成功消息

## shfmt -l 列表模式
- [x] `fck shfmt -l script.sh`（需要格式化）输出 `script.sh`
- [x] `fck shfmt -l script.sh`（已格式化）不输出内容
- [x] `fck shfmt -l *.sh` 仅输出需要格式化的文件路径（每行一个）
- [x] `fck shfmt -l -w script.sh` -l 生效，不写入文件，仅列表
- [x] `fck shfmt --list script.sh` 长选项同样生效
- [x] 未指定 `-l` 时行为不变，正常输出格式化结果

## shx 退出码透传
- [x] `fck shx "exit 0"` 退出码为 0
- [x] `fck shx "exit 1"` 退出码为 1，无额外错误信息输出
- [x] `fck shx "exit 42"` 退出码为 42
- [x] `fck shx -t 1s "sleep 10"` 超时输出错误信息，退出码非 0
- [x] `fck shx "echo hello"` 正常命令退出码为 0

## 编译
- [x] go build 编译通过，无错误
