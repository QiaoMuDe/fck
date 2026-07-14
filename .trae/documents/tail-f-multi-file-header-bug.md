# tail -f 多文件模式缺少文件名标题 Bug 修复

## 问题描述

`fck tail -f *.log` 在有新内容追加时，输出的行没有文件名标题（`==> filename <==`），无法区分新内容是来自哪个文件。

**用户观察到的现象：**
```
==> INFO.log <==        ← 初始内容有标题
2026-07-14 11:09:43 | INFO | ...
==> app.log <==         ← 初始内容有标题
2026-07-14 11:09:43 | INFO | ...
11                      ← 动态追加内容没有标题，不知道来自哪个文件
121
121
```

## 根因分析

[`followFile`](file:///d:%5C%E8%B5%84%E6%BA%90%E6%B1%A0%5C%E4%B8%8B%E6%B0%B4%E9%81%93%5CDev%5C%E6%9C%AC%E5%9C%B0%E9%A1%B9%E7%9B%AE%5Cfck%5Cinternal%5Ccommands%5Ctail%5Ccmd_tail.go) 函数中第 329 行的标题打印条件为 `showHeader && tf.Size == 0`：

```go
if showHeader && tf.Size == 0 {
    fmt.Printf("\n==> %s <==\n", tf.Path)
}
```

`tf.Size == 0` 只在文件被截断/替换后才会成立（第 310 行 `tf.Size = 0`）。对于正常追加的场景，`tf.Size` 始终大于 0，因此每次轮询检测到新内容时不会打印标题。

**关键点：** 第 314-316 行的 `newSize == tf.Size` 提前返回已确保 `followFile` 只在有增量内容时才会执行到读取循环，因此这里的 `tf.Size == 0` 条件是多余的——能走到这里必然有新内容，多文件模式就应该打标题。

## 修改方案

修改 [`internal/commands/tail/cmd_tail.go`](file:///d:%5C%E8%B5%84%E6%BA%90%E6%B1%A0%5C%E4%B8%8B%E6%B0%B4%E9%81%93%5CDev%5C%E6%9C%AC%E5%9C%B0%E9%A1%B9%E7%9B%AE%5Cfck%5Cinternal%5Ccommands%5Ctail%5Ccmd_tail.go) 的 `followFile` 函数：

1. **移除** 第 329-331 行的 `if showHeader && tf.Size == 0 { ... }`（读取循环内的标题打印）
2. **新增** 在新内容读取前（第 317 行 return nil 之后），统一打印标题的条件 `showHeader`

优化后的逻辑流程：
```
新内容检测
  ├── 文件被截断 → 打印标题 (保留)
  ├── 无新内容 → return nil (保留)
  └── 有新内容 →
        ├── 多文件模式 → 打印标题 (新增, 原先缺失)
        └── 读取并输出新内容
```

## 验证

1. `go build ./internal/commands/tail/` 编译通过
2. 手动测试：`echo "line1" > a.log && echo "line1" > b.log && tail -f *.log`，追加内容时验证标题正确显示
