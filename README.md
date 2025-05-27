# FCK 工具

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/QiaoMuDe/fck)

一个多功能命令行工具，提供文件哈希计算、大小统计、校验和查找等功能。

## 功能特性

- 计算文件哈希值（支持MD5、SHA1、SHA256、SHA512算法）
- 递归处理目录
- 多线程并发处理
- 文件大小统计
- 文件查找功能
- 哈希校验功能

## 安装

1. 确保已安装Go环境（1.24+）
2. 克隆仓库：
   ```
   git clone https://gitee.com/MM-Q/fck.git
   ```
3. 构建项目：
   ```
   cd fck
   python3 build.py
   ```

## 使用说明

### 基本命令

```
./fck [命令] [选项]
```

### 可用命令

- `hash`: 计算文件哈希值（支持MD5/SHA1/SHA256/SHA512算法）

- `size`: 统计文件大小（支持递归计算目录大小）

- `diff`: 目录差异比较（支持文件哈希值校验）

- `find`: 查找文件（支持按名称、大小、时间等条件筛选）

## 许可证

本项目采用 GNU General Public License v3.0 许可证。

## 贡献

欢迎提交Issue和Pull Request。