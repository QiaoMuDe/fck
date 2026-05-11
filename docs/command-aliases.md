# FCK 子命令别名设计方案

## 设计原则
1. 优先使用常见 Unix 命令的缩写习惯
2. 保持简洁，通常 1-3 个字符
3. 避免冲突，确保别名唯一
4. 易于记忆，符合直觉

## 别名映射表

| 命令 | 建议别名 | 说明 |
|------|----------|------|
| list, ls | ls | 已有别名，保持不变 |
| awk | awk | 保持原名，Unix 标准工具 |
| base64 | b64 | Base64 常用缩写 |
| cat | cat | 保持原名，Unix 标准工具 |
| check | chk | Check 的常用缩写 |
| cp | cp | 保持原名，Unix 标准工具 |
| curl | curl | 保持原名，知名工具 |
| date | date | 保持原名，Unix 标准工具 |
| df | df | 保持原名，Unix 标准工具 |
| dns | dns | 保持原名，清晰明了 |
| echo | echo | 保持原名，Unix 标准工具 |
| find | find | 保持原名，Unix 标准工具 |
| gm | gm | 保持原名，Git Manager 缩写 |
| grep | grep | 保持原名，Unix 标准工具 |
| hash | hash | 保持原名，清晰明了 |
| head | head | 保持原名，Unix 标准工具 |
| home | ~ | 主目录符号，直观 |
| iconv | icv | IConv 的缩写 |
| json | json | 保持原名，清晰明了 |
| md | md | 保持原名，Markdown 缩写 |
| mkdir | mkdir | 保持原名，Unix 标准工具 |
| mv | mv | 保持原名，Unix 标准工具 |
| newline | nl | NewLine 的缩写，与 Unix nl 命令一致 |
| pack | pk | Pack 的缩写 |
| ping | ping | 保持原名，Unix 标准工具 |
| port | pt | Port 的缩写 |
| preview | pv | PreView 的缩写 |
| proc | ps | Process 的标准缩写，与 Unix ps 一致 |
| pwd | pwd | 保持原名，Unix 标准工具 |
| rm | rm | 保持原名，Unix 标准工具 |
| sed | sed | 保持原名，Unix 标准工具 |
| seq | seq | 保持原名，Unix 标准工具 |
| size | sz | Size 的缩写 |
| tail | tail | 保持原名，Unix 标准工具 |
| tcp | tcp | 保持原名，协议名 |
| tee | tee | 保持原名，Unix 标准工具 |
| test | test | 保持原名，Unix 标准工具 |
| touch | touch | 保持原名，Unix 标准工具 |
| tr | tr | 保持原名，Unix 标准工具 |
| truncate | trunc | Truncate 的缩写 |
| unpack | upk | UnPack 的缩写 |
| watch | wch | Watch 的缩写（避免与 wc 冲突） |
| wc | wc | 保持原名，Unix 标准工具 |
| which | wh | Which 的缩写 |
| xargs | xargs | 保持原名，Unix 标准工具 |

## 别名分类统计

### 保持原名的命令（22个）
awk, cat, cp, curl, date, df, dns, echo, find, grep, hash, head, mkdir, mv, ping, pwd, rm, sed, seq, tail, tcp, tee, test, touch, tr, wc, xargs

### 新增别名（11个）
- base64 → b64
- check → chk
- home → ~
- iconv → icv
- newline → nl
- pack → pk
- port → pt
- preview → pv
- proc → ps
- size → sz
- truncate → trunc
- unpack → upk
- watch → wch
- which → wh

## 使用示例

```bash
# 使用别名
fck nl file.txt          # 检测换行符
fck icv -t UTF-8 file.txt # 编码转换
fck pk -o backup.zip dir  # 打包
fck upk archive.zip       # 解压
fck ps                    # 查看进程
fck sz dir/               # 计算大小
fck chk checksum.hash     # 校验文件
fck ~                     # 显示主目录
```
