package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/globals"
)

// 定义全局常量的颜色映射
var PermissionColorMap = map[int]string{
	1: "green",  // 所有者-读-绿色
	2: "yellow", // 所有者-写-黄色
	3: "red",    // 所有者-执行-红色
	4: "green",  // 组-读-绿色
	5: "yellow", // 组-写-黄色
	6: "red",    // 组-执行-红色
	7: "green",  // 其他-读-绿色
	8: "yellow", // 其他-写-黄色
	9: "red",    // 其他-执行-红色
}

// 颜色映射表，用于开发环境配色方案
var devColorMap = map[string]map[string]bool{
	// 代码文件或脚本文件
	"green": {
		".go":          true, // go 源文件
		".py":          true, // python 源文件
		".pyw":         true, // python 源文件
		".sh":          true, // shell 脚本
		".bash":        true, // shell 脚本
		".zsh":         true, // shell 脚本
		".js":          true, // javascript 源文件
		".ts":          true, // typescript 源文件
		".jsx":         true, // react js
		".tsx":         true, // react ts
		".html":        true, // html 文件
		".css":         true, // css 文件
		".scss":        true, // sass 样式文件
		".sass":        true, // sass 样式文件
		".less":        true, // less 样式文件
		".rs":          true, // rust 源文件
		".rb":          true, // ruby 源文件
		".php":         true, // php 文件
		".java":        true, // java 文件
		".c":           true, // c 源文件
		".cpp":         true, // c++ 源文件
		".h":           true, // c++ 头文件
		".hpp":         true, // c++ 头文件
		".m":           true, // objective-c 源文件
		".mm":          true, // objective-c++ 源文件
		".swift":       true, // swift 源文件
		".kt":          true, // kotlin 源文件
		".kts":         true, // kotlin 脚本文件
		".scala":       true, // scala 源文件
		".lua":         true, // lua 源文件
		".pl":          true, // perl 源文件
		".pm":          true, // perl 模块文件
		".r":           true, // R 语言源文件
		".sql":         true, // sql 文件
		".mssql":       true, // mssql 脚本文件
		".ps1":         true, // powershell 脚本文件
		".psm1":        true, // powershell 模块文件
		".bat":         true, // windows 批处理文件
		".cmd":         true, // windows 批处理文件
		".vbs":         true, // vbscript 脚本文件
		".vba":         true, // vba 脚本文件
		".asm":         true, // 汇编语言源文件
		".s":           true, // 汇编语言源文件
		".ada":         true, // ada 源文件
		".fs":          true, // f# 源文件
		".hs":          true, // haskell 源文件
		".lisp":        true, // lisp 源文件
		".clj":         true, // clojure 源文件
		".erl":         true, // erlang 源文件
		".hrl":         true, // erlang 头文件
		".ex":          true, // elixir 源文件
		".exs":         true, // elixir 脚本文件
		".groovy":      true, // groovy 源文件
		".dart":        true, // dart 源文件
		".vue":         true, // vue 组件文件
		".svelte":      true, // svelte 组件文件
		".astro":       true, // astro 组件文件
		".pug":         true, // pug 模板文件
		".jade":        true, // jade 模板文件
		".haml":        true, // haml 模板文件
		".mdx":         true, // markdown 带 jsx 文件
		".aspx":        true, // asp.net 页面文件
		".ascx":        true, // asp.net 用户控件文件
		".cs":          true, // c# 源文件
		".vb":          true, // visual basic 源文件
		".fsx":         true, // f# 脚本文件
		".csproj":      true, // c# 项目文件
		".vbproj":      true, // visual basic 项目文件
		".fsproj":      true, // f# 项目文件
		".sln":         true, // visual studio 解决方案文件
		".jsp":         true, // java server pages
		".jspx":        true, // java server pages xml
		".php3":        true, // php 3 文件
		".php4":        true, // php 4 文件
		".php5":        true, // php 5 文件
		".phtml":       true, // php html 混合文件
		".twig":        true, // twig 模板文件
		".smarty":      true, // smarty 模板文件
		".erb":         true, // ruby erb 模板文件
		".tpl":         true, // 通用模板文件
		".vm":          true, // velocity 模板文件
		".ftl":         true, // freemarker 模板文件
		".njk":         true, // nunjucks 模板文件
		".hbs":         true, // handlebars 模板文件
		".mustache":    true, // mustache 模板文件
		".ejs":         true, // embedded javascript 模板文件
		".tsbuildinfo": true, // typescript 增量编译信息文件
	},
	// 配置文件或环境文件
	"yellow": {
		".ini":                     true, // ini 配置文件
		".conf":                    true, // conf 配置文件
		".cfg":                     true, // 通用配置文件
		".json":                    true, // JSON 格式配置文件
		".yaml":                    true, // YAML 格式配置文件
		".yml":                     true, // YAML 格式配置文件（简写）
		".xml":                     true, // XML 格式配置文件
		".toml":                    true, // TOML 格式配置文件
		".env":                     true, // 环境变量配置文件
		".properties":              true, // Java 属性配置文件
		".config":                  true, // 通用配置文件
		".settings":                true, // 设置配置文件
		".md":                      true, // Markdown 文档文件
		".markdown":                true, // Markdown 文档文件（完整写法）
		".mod":                     true, // Go 模块文件
		".sum":                     true, // Go 模块校验和文件
		".lock":                    true, // 依赖锁定文件
		".npmrc":                   true, // npm 配置文件
		".yarnrc":                  true, // yarn 配置文件
		".pnpmrc":                  true, // pnpm 配置文件
		".editorconfig":            true, // 编辑器配置文件
		".gitconfig":               true, // git 全局配置文件
		".gitignore":               true, // git 忽略配置文件
		".gitattributes":           true, // git 属性配置文件
		".dockerignore":            true, // docker 忽略配置文件
		".dockerfile":              true, // docker 构建文件
		".docker-compose.yml":      true, // docker 组合配置文件
		".docker-compose.yaml":     true, // docker 组合配置文件
		".kubeconfig":              true, // kubernetes 配置文件
		".helm":                    true, // helm 配置文件
		".travis.yml":              true, // travis ci 配置文件
		".gitlab-ci.yml":           true, // gitlab ci 配置文件
		".github/workflows/*.yml":  true, // github actions 配置文件
		".github/workflows/*.yaml": true, // github actions 配置文件
		".eslintrc":                true, // eslint 配置文件
		".eslintrc.json":           true, // eslint 配置文件
		".eslintrc.js":             true, // eslint 配置文件
		".eslintrc.yml":            true, // eslint 配置文件
		".prettierrc":              true, // prettier 配置文件
		".prettierrc.json":         true, // prettier 配置文件
		".prettierrc.js":           true, // prettier 配置文件
		".prettierrc.yml":          true, // prettier 配置文件
		".prettierignore":          true, // prettier 忽略配置文件
		".stylelintrc":             true, // stylelint 配置文件
		".stylelintrc.json":        true, // stylelint 配置文件
		".stylelintrc.js":          true, // stylelint 配置文件
		".stylelintrc.yml":         true, // stylelint 配置文件
		".stylelintignore":         true, // stylelint 忽略配置文件
		".babelrc":                 true, // babel 配置文件
		".babelrc.json":            true, // babel 配置文件
		".babelrc.js":              true, // babel 配置文件
		".webpack.config.js":       true, // webpack 配置文件
		".rollup.config.js":        true, // rollup 配置文件
		".vite.config.js":          true, // vite 配置文件
		".vite.config.ts":          true, // vite 配置文件
		".tsconfig.json":           true, // typescript 配置文件
		".jsconfig.json":           true, // javascript 配置文件
		".vue.config.js":           true, // vue 项目配置文件
		".nuxt.config.js":          true, // nuxt 项目配置文件
		".nuxt.config.ts":          true, // nuxt 项目配置文件
		".next.config.js":          true, // next.js 项目配置文件
		".gatsby-config.js":        true, // gatsby 项目配置文件
		".postcss.config.js":       true, // postcss 配置文件
		".tailwind.config.js":      true, // tailwind 配置文件
		".jest.config.js":          true, // jest 配置文件
		".cypress.json":            true, // cypress 配置文件
		".env.development":         true, // 开发环境变量配置文件
		".env.production":          true, // 生产环境变量配置文件
		".env.test":                true, // 测试环境变量配置文件
		".env.local":               true, // 本地环境变量配置文件
		".env.development.local":   true, // 本地开发环境变量配置文件
		".env.production.local":    true, // 本地生产环境变量配置文件
		".env.test.local":          true, // 本地测试环境变量配置文件
		".htaccess":                true, // apache 配置文件
		".nginx.conf":              true, // nginx 配置文件
		".vhost":                   true, // 虚拟主机配置文件
		".htpasswd":                true, // apache 认证文件
		".env.example":             true, // 环境变量示例文件
		".env.sample":              true, // 环境变量示例文件
		".env.dist":                true, // 环境变量分发文件
		".env.template":            true, // 环境变量模板文件
		".circleci/config.yml":     true, // circleci 配置文件
		".jenkinsfile":             true, // jenkins 配置文件
		".makefile":                true, // makefile 文件
		".mk":                      true, // makefile 文件
		".cmake":                   true, // cmake 配置文件
		".cscope.out":              true, // cscope 配置文件
		".tags":                    true, // ctags 配置文件
		".project":                 true, // eclipse 项目文件
		".classpath":               true, // eclipse 类路径文件
		".iml":                     true, // intellij idea 模块文件
		".ipr":                     true, // intellij idea 项目文件
		".iws":                     true, // intellij idea 工作区文件
		".eslintignore":            true, // eslint 忽略文件
		".npmignore":               true, // npm 忽略文件
		".yarnignore":              true, // yarn 忽略文件
		".gitmodules":              true, // git 子模块配置文件
		".d.ts":                    true, // typescript 声明文件
		".backup":                  true, // 备份文件
		".bakup":                   true, // 备份文件
		".metadata":                true, // 元数据文件
		".swp":                     true, // Vim 交换文件
		".swo":                     true, // Vim 交换文件
		".swn":                     true, // Vim 交换文件
		".lck":                     true, // 锁文件
		".pid":                     true, // 进程 ID 文件
		".err":                     true, // 错误日志文件
		".workspace":               true, // Eclipse 工作区文件
		".nvmrc":                   true, // Node Version Manager 配置文件
		".gdbinit":                 true, // GDB 初始化文件
		".lldbinit":                true, // LLDB 初始化文件
		".valgrindrc":              true, // Valgrind 配置文件
		".clang-format":            true, // Clang 格式化配置文件
		".clang-tidy":              true, // Clang 静态分析配置文件
	},
	// 数据文件或压缩文件或数据库文件
	"red": {
		".zip":       true, // ZIP 压缩文件
		".tar":       true, // TAR 归档文件
		".gz":        true, // GZIP 压缩文件
		".bz2":       true, // BZIP2 压缩文件
		".rar":       true, // RAR 压缩文件
		".7z":        true, // 7-Zip 压缩文件
		".tar.gz":    true, // TAR 归档并 GZIP 压缩文件
		".tar.bz2":   true, // TAR 归档并 BZIP2 压缩文件
		".tgz":       true, // TAR 归档并 GZIP 压缩文件（简写）
		".xz":        true, // XZ 压缩文件
		".lzma":      true, // LZMA 压缩文件
		".tar.xz":    true, // TAR 归档并 XZ 压缩文件
		".tbz2":      true, // TAR 归档并 BZIP2 压缩文件（简写）
		".tbz":       true, // TAR 归档并 BZIP2 压缩文件（简写）
		".txz":       true, // TAR 归档并 XZ 压缩文件（简写）
		".lz":        true, // LZ 压缩文件
		".lz4":       true, // LZ4 压缩文件
		".jar":       true, // Java 归档文件
		".war":       true, // Java Web 应用归档文件
		".ear":       true, // Java 企业应用归档文件
		".apk":       true, // Android 应用程序包
		".ipa":       true, // iOS 应用程序包
		".db":        true, // SQLite 数据库文件
		".sqlite":    true, // SQLite 数据库文件（简写）
		".db3":       true, // SQLite 3 数据库文件
		".sqlite3":   true, // SQLite 3 数据库文件（简写）
		".mdb":       true, // Microsoft Access 数据库文件
		".pdb":       true, // Microsoft Access 数据库文件
		".accdb":     true, // Microsoft Access 2007+ 数据库文件
		".sqlitedb":  true, // SQLite 数据库文件
		".sqlite3db": true, // SQLite 3 数据库文件
		".dbf":       true, // dBase 数据库文件
		".mdf":       true, // Microsoft SQL Server 数据库文件
		".ldf":       true, // Microsoft SQL Server 日志文件
		".ndf":       true, // Microsoft SQL Server 辅助数据文件
		".tar.bz":    true, // tar.bz 压缩文件
		".tar.lz":    true, // tar.lz 压缩文件
		".tar.lz4":   true, // tar.lz4 压缩文件
		".lzop":      true, // lzop 压缩文件
		".z":         true, // compress 压缩文件
		".zst":       true, // zstandard 压缩文件
		".zstd":      true, // zstandard 压缩文件
		".br":        true, // brotli 压缩文件
		".cab":       true, // Windows  Cabinet 压缩文件
		".arj":       true, // ARJ 压缩文件
		".lzh":       true, // LZH 压缩文件
		".lha":       true, // LHA 压缩文件
		".alz":       true, // ALZip 压缩文件
		".ace":       true, // ACE 压缩文件
		".bz":        true, // BZ 压缩文件
		".zoo":       true, // ZOO 压缩文件
		".arc":       true, // ARC 压缩文件
		".sit":       true, // StuffIt 压缩文件
		".sitx":      true, // StuffIt 扩展压缩文件
		".dmg":       true, // macOS 磁盘映像文件
		".iso":       true, // 光盘映像文件
		".img":       true, // 磁盘映像文件
		".vhd":       true, // 虚拟硬盘文件
		".vmdk":      true, // VMware 虚拟磁盘文件
		".qcow2":     true, // QEMU 磁盘映像文件
		".raw":       true, // 原始磁盘映像文件
		".bak":       true, // 备份文件
		".bkf":       true, // Windows 备份文件
		".fdb":       true, // Firebird 数据库文件
		".gdb":       true, // InterBase 数据库文件
		".ibd":       true, // MySQL 索引和数据文件
		".frm":       true, // MySQL 表结构文件
		".myd":       true, // MySQL 数据文件
		".myi":       true, // MySQL 索引文件
		".ndb":       true, // MySQL Cluster 数据文件
		".pgdata":    true, // PostgreSQL 数据目录
		".whl":       true, // Python wheel 包文件
		".pex":       true, // Python 可执行包文件
		".rpm":       true, // Red Hat 包管理器文件
		".deb":       true, // Debian 包文件
		".aab":       true, // Android App Bundle
		".xap":       true, // Windows Phone 应用包
		".appx":      true, // Windows 应用包
		".pkg":       true, // macOS 安装包
		".msi":       true, // Windows 安装包
	},
	// 库文件或静态库文件或编译后的产物
	"purple": {
		".so":            true, // Linux 共享对象库文件
		".dll":           true, // Windows 动态链接库文件
		".lib":           true, // Windows 静态链接库文件
		".pyc":           true, // Python 字节码文件
		".pyd":           true, // Python 扩展模块文件
		".o":             true, // 目标文件
		".a":             true, // Linux 静态库文件
		".dylib":         true, // macOS 动态链接库文件
		".class":         true, // Java 字节码文件
		".exe":           true, // 可执行文件
		".com":           true, // 命令文件
		".bin":           true, // 二进制文件
		".elf":           true, // Linux 可执行文件
		".out":           true, // 编译输出文件
		".obj":           true, // 目标文件
		".exp":           true, // Windows 导出文件
		".pdb":           true, // Windows 调试符号文件
		".ilk":           true, // Windows 增量链接文件
		".pch":           true, // 预编译头文件
		".sbr":           true, // 源代码浏览器文件
		".idb":           true, // Visual Studio 增量调试文件
		".idc":           true, // IDA Pro 脚本文件
		".id0":           true, // IDA Pro 数据库文件
		".id1":           true, // IDA Pro 数据库文件
		".id2":           true, // IDA Pro 数据库文件
		".nam":           true, // IDA Pro 名称文件
		".til":           true, // IDA Pro 类型库文件
		".pyo":           true, // Python 优化字节码文件
		".pyi":           true, // Python 类型提示文件
		".egg":           true, // Python 包文件
		".jmod":          true, // Java 模块文件
		".kotlin_module": true, // Kotlin 模块文件
		".ktx":           true, // Kotlin 扩展文件
		".swiftmodule":   true, // Swift 模块文件
		".node":          true, // Node.js 原生插件
		".wasm":          true, // WebAssembly 模块
		".bc":            true, // LLVM 位码文件
		".ll":            true, // LLVM 中间表示文件
		".map":           true, // 映射文件
		".lst":           true, // 列表文件
		".d":             true, // 依赖文件
		".rlib":          true, // Rust 库文件
	},
}

// splitPathColor 函数用于根据路径类型以不同颜色返回字符串
func splitPathColor(p string, cl *colorlib.ColorLib, dirCode int, fileCode int) string {
	// 获取路径的目录和文件名
	dir, file := filepath.Split(p)

	// 如果目录为空，则返回文件名
	if dir == "" {
		return fmt.Sprint(cl.SColor(fileCode, file))
	}

	// 设置渐进式颜色
	return fmt.Sprint(cl.SColor(dirCode, dir), cl.SColor(fileCode, file))
}

// SPrintStringColor 根据路径类型以不同颜色输出字符串
// 参数:
//
//	p: 要检查的路径(用于获取文件类型信息)
//	s: 返回字符串内容
//	cl: colorlib.ColorLib实例, 用于彩色输出
func SprintStringColor(p string, s string, cl *colorlib.ColorLib) string {
	// 获取路径信息
	pathInfo, statErr := os.Lstat(p)
	if statErr != nil {
		return cl.Sred(s) // 如果获取路径信息失败, 返回红色输出
	}

	// 根据路径类型设置颜色
	switch mode := pathInfo.Mode(); {
	case mode&os.ModeSymlink != 0:
		// 符号链接 - 使用青色输出
		return cl.Scyan(s)
	case runtime.GOOS == "windows" && mode.IsRegular() && globals.WindowsSymlinkExts[filepath.Ext(p)]:
		// Windows下的快捷方式文件 - 使用青色输出
		return cl.Scyan(s)
	case mode.IsDir():
		// 目录 - 使用蓝色输出
		return cl.Sblue(s)
	case mode&os.ModeDevice != 0:
		// 设备文件 - 使用黄色输出
		return cl.Syellow(s)
	case mode&os.ModeNamedPipe != 0:
		// 命名管道 - 使用黄色输出
		return cl.Syellow(s)
	case mode&os.ModeSocket != 0:
		// 套接字文件 - 使用黄色输出
		return cl.Syellow(s)
	case mode&os.ModeType == 0 && mode&os.ModeCharDevice != 0:
		// 字符设备文件 - 使用黄色输出
		return cl.Syellow(s)
	case mode.IsRegular() && pathInfo.Size() == 0:
		// 空文件 - 使用灰色输出
		return cl.Sgray(s)
	case mode.IsRegular() && mode&0111 != 0:
		// 可执行文件 - 使用绿色输出
		return cl.Sgreen(s)
	case runtime.GOOS == "windows" && mode.IsRegular() && globals.WindowsExecutableExts[filepath.Ext(p)]:
		// Windows下的可执行文件 - 使用绿色输出
		return cl.Sgreen(s)
	case mode.IsRegular():
		// 普通文件 - 使用白色输出
		return cl.Swhite(s)
	default:
		// 其他类型文件 - 使用白色输出
		return cl.Swhite(s)
	}
}

// printPathColor 根据路径类型以不同颜色输出路径字符串
// 参数:
//
//	path: 要检查的路径(用于获取文件类型信息)
//	cl: colorlib.ColorLib实例, 用于彩色输出
func printPathColor(path string, cl *colorlib.ColorLib) {
	// 获取路径信息
	pathInfo, statErr := os.Lstat(path)
	if statErr != nil {
		// 如果获取路径信息失败, 输出红色的路径
		fmt.Println(splitPathColor(path, cl, colorlib.Cyan, colorlib.Red))
		return
	}

	// 根据路径类型设置颜色
	switch mode := pathInfo.Mode(); {
	case mode&os.ModeSymlink != 0:
		// 符号链接 - 使用青色输出
		fmt.Println(splitPathColor(path, cl, colorlib.Cyan, colorlib.Cyan))
	case runtime.GOOS == "windows" && mode.IsRegular() && globals.WindowsSymlinkExts[filepath.Ext(path)]:
		// Windows快捷方式 - 使用青色输出
		fmt.Println(splitPathColor(path, cl, colorlib.Cyan, colorlib.Cyan))
	case mode.IsDir():
		// 目录 - 使用蓝色输出
		fmt.Println(splitPathColor(path, cl, colorlib.Cyan, colorlib.Blue))
	case mode&os.ModeDevice != 0:
		// 设备文件 - 使用黄色输出
		fmt.Println(splitPathColor(path, cl, colorlib.Cyan, colorlib.Yellow))
	case mode&os.ModeCharDevice != 0:
		// 字符设备文件 - 使用黄色输出
		fmt.Println(splitPathColor(path, cl, colorlib.Cyan, colorlib.Yellow))
	case mode&os.ModeNamedPipe != 0:
		// 命名管道 - 使用黄色输出
		fmt.Println(splitPathColor(path, cl, colorlib.Cyan, colorlib.Yellow))
	case mode&os.ModeSocket != 0:
		// 套接字文件 - 使用黄色输出
		fmt.Println(splitPathColor(path, cl, colorlib.Cyan, colorlib.Yellow))
	case mode&os.ModeType == 0 && mode&0111 != 0:
		// 可执行文件 - 使用绿色输出
		fmt.Println(splitPathColor(path, cl, colorlib.Cyan, colorlib.Green))
	case runtime.GOOS == "windows" && mode.IsRegular() && globals.WindowsExecutableExts[filepath.Ext(path)]:
		// Windows可执行文件 - 使用绿色输出
		fmt.Println(splitPathColor(path, cl, colorlib.Cyan, colorlib.Green))
	case pathInfo.Size() == 0:
		// 空文件 - 使用灰色输出
		fmt.Println(splitPathColor(path, cl, colorlib.Cyan, colorlib.Gray))
	case mode.IsRegular():
		// 普通文件 - 使用白色输出
		fmt.Println(splitPathColor(path, cl, colorlib.Cyan, colorlib.White))
	default:
		// 其他类型文件 - 使用白色输出
		fmt.Println(splitPathColor(path, cl, colorlib.Cyan, colorlib.White))
	}
}

// getColorString 函数的作用是根据传入的文件信息、路径字符串以及颜色库实例，返回带有相应颜色的路径字符串。
// 参数:
// info: 包含文件类型和文件后缀名等信息的 globals.ListInfo 结构体实例。
// pF: 要处理的路径字符串。
// cl: 用于彩色输出的 colorlib.ColorLib 实例。
// 返回值:
// colorString: 经过颜色处理后的路径字符串。
func getColorString(info globals.ListInfo, pF string, cl *colorlib.ColorLib) string {
	// 如果启用了开发环境模式, 则返回开发模式下的颜色处理结果
	if listCmdDevColor.Get() {
		return getDevColorString(info, pF, cl)
	}

	// 依据文件的类型来确定输出的颜色
	switch info.EntryType {
	case globals.SymlinkType:
		// 若文件类型为符号链接，则使用青色来渲染字符串
		return cl.Scyan(pF)
	case globals.DirType:
		// 若文件类型为目录，则使用蓝色来渲染字符串
		return cl.Sblue(pF)
	case globals.ExecutableType:
		// 若文件类型为可执行文件，则使用绿色来渲染字符串
		return cl.Sgreen(pF)
	case globals.SocketType, globals.PipeType, globals.BlockDeviceType, globals.CharDeviceType:
		// 若文件类型为套接字、管道、块设备、字符设备，则使用黄色来渲染字符串
		return cl.Syellow(pF)
	case globals.EmptyType:
		// 若文件类型为空文件, 则使用灰色来渲染字符串
		return cl.Sgray(pF)
	case globals.FileType:
		// 若文件类型为普通文件，则根据平台差异来设置颜色
		if runtime.GOOS == "windows" {
			switch {
			case globals.WindowsExecutableExts[info.FileExt]:
				// 对于 Windows 系统下的可执行文件，使用绿色来渲染字符串
				return cl.Sgreen(pF)
			case globals.WindowsSymlinkExts[info.FileExt]:
				// 对于 Windows 系统下的符号链接，使用青色来渲染字符串
				return cl.Scyan(pF)
			default:
				// 对于其他文件类型，使用白色来渲染字符串
				return cl.Swhite(pF)
			}
		}

		// 添加MacOS特殊文件处理
		if runtime.GOOS == "darwin" {
			base := filepath.Base(pF)
			switch {
			case base == ".DS_Store" || base == ".localized" || strings.HasPrefix(base, "._"):
				return cl.Sgray(pF) // MacOS系统文件使用灰色
			case filepath.Ext(pF) == ".app":
				return cl.Sgreen(pF) // MacOS应用程序包使用绿色
			}
		}

		// 对于 Linux 系统下的普通文件，使用白色来渲染字符串
		return cl.Swhite(pF)
	default:
		// 对于未匹配的类型，使用白色来渲染字符串
		return cl.Swhite(pF)
	}
}

// getDeColorString 函数的作用是根据传入的文件信息、路径字符串以及颜色库实例，返回不带颜色的路径字符串(开发环境配色方案)
// 函数参数:
// info: 包含文件类型和文件后缀名等信息的 globals.ListInfo 结构体实例。
// pF: 要处理的路径字符串。
// cl: 用于处理颜色的 colorlib.ColorLib 实例。
// 返回值:
// decolorString: 不带颜色的路径字符串。
func getDevColorString(info globals.ListInfo, pF string, cl *colorlib.ColorLib) string {
	// 依据文件的类型来确定输出的颜色
	switch info.EntryType {
	case globals.SymlinkType:
		// 若文件类型为符号链接，则使用青色来渲染字符串
		return cl.Scyan(pF)
	case globals.DirType:
		// 若文件类型为目录，则使用蓝色来渲染字符串
		return cl.Sblue(pF)
	case globals.ExecutableType:
		// 若文件类型为可执行文件，则使用绿色来渲染字符串
		return cl.Sgreen(pF)
	case globals.SocketType, globals.PipeType, globals.BlockDeviceType, globals.CharDeviceType, globals.EmptyType:
		// 若文件类型为套接字, 管道, 块设备, 字符设备, 空文件, 则使用灰色来渲染字符串
		return cl.Sgray(pF)
	case globals.FileType:
		// 若文件类型为普通文件，则根据开发环境配色方案来设置颜色
		if runtime.GOOS == "windows" {
			switch {
			case globals.WindowsExecutableExts[info.FileExt]:
				// 对于 Windows 系统下的可执行文件，使用绿色来渲染字符串
				return cl.Sgreen(pF)
			case globals.WindowsSymlinkExts[info.FileExt]:
				// 对于 Windows 系统下的符号链接，使用青色来渲染字符串
				return cl.Scyan(pF)
			}
		}

		// 常规配色方案
		for color, extMap := range devColorMap {
			if extMap[info.FileExt] {
				switch color {
				case "yellow":
					return cl.Syellow(pF)
				case "green":
					return cl.Sgreen(pF)
				case "red":
					return cl.Sred(pF)
				case "purple":
					return cl.Spurple(pF)
				}
			}
		}

		// 普通文件，使用白色来渲染字符串
		return cl.Swhite(pF)
	default:
		// 对于未匹配的类型, 使用灰色来渲染字符串
		return cl.Sgray(pF)
	}
}
