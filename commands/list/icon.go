package list

// 图标映射与选择逻辑
// 本文件提供：
//  1) IconMap 数据结构：集中维护“类型/扩展名 → 图标编码”的映射
//  2) getIcon(info FileInfo)：选择策略为：
//     - 普通文件优先按扩展名匹配（如 ".zip"/"zip"）；
//     - 非普通文件直接按 EntryType 匹配；

//     - 若未命中扩展名或类型映射，将返回对应类型在 ByType 中的值；若为空则回退到 Default；
//
// 说明：
//  - 图标编码默认使用 Nerd Font 私有区字符；若终端未安装对应字体，可能显示为方块，可在后续增加 emoji/纯文本降级。
//  - 本文件仅返回“图标编码字符串”，不负责着色或与名称拼接，由 formatter 决定是否显示、如何着色。

// IconMap 定义图标映射集合
type IconMap struct {
	// 按扩展名映射，键为小写扩展名（支持包含"."或不包含"."的两种）
	ByExt map[string]string
	// 按文件类型映射
	ByType map[EntryType]string
	// 默认图标
	Default string
}

// DefaultIcons 默认图标映射表
// 注：编码均为字符串，可直接拼接到名称前；是否加空格/着色由调用方决定。
var DefaultIcons = IconMap{
	ByType: map[EntryType]string{
		DirType:         "\uf4d4", // 目录
		SymlinkType:     "\uF482", // 软链接
		SocketType:      "\uf4d6", // 套接字
		PipeType:        "\uf4d6", // 管道
		BlockDeviceType: "\uf4d6", // 块设备
		CharDeviceType:  "\uf4d6", // 字符设备
		ExecutableType:  "\uf4d3", // 可执行文件
		EmptyType:       "\uf4d3", // 空文件
		FileType:        "\uf4d3", // 普通文件
		UnknownType:     "\uf4d3", // 未知类型
	},

	ByExt: map[string]string{
		// Go语言
		".go": "\uE627",

		// Python
		".py":  "\uED1B",
		".pyw": "\uED1B",

		// Shell脚本
		".sh":   "\uE691",
		".bash": "\uE691",
		".zsh":  "\uE691",

		// JavaScript/TypeScript
		".js":  "\uf2ee",
		".ts":  "\uE69D",
		".jsx": "\uf2ee",
		".tsx": "\uE69D",

		// Web前端
		".html": "\uE736",
		".css":  "\uE736",
		".scss": "\uE736",
		".sass": "\uE736",
		".less": "\uE736",
		".vue":  "\uE736",

		// 其他编程语言
		".rs":     "\uE7A8", // Rust
		".rb":     "\ueae9", // Ruby
		".php":    "\uE73D", // PHP
		".java":   "\ue66d", // Java
		".c":      "\ue649", // C
		".cpp":    "\ue649", // C++
		".h":      "\ue649", // C/C++ 头文件
		".hpp":    "\ue649", // C++ 头文件
		".m":      "\ueae9", // Objective-C
		".mm":     "\ueae9", // Objective-C++
		".swift":  "\ueae9", // Swift
		".kt":     "\ueae9", // Kotlin
		".kts":    "\ueae9", // Kotlin 脚本
		".scala":  "\ueae9", // Scala
		".lua":    "\ueae9", // Lua
		".pl":     "\ueae9", // Perl
		".pm":     "\ueae9", // Perl 模块
		".r":      "\ueae9", // R语言
		".sql":    "\uf472", // SQL
		".ps1":    "\ue86c", // PowerShell
		".psm1":   "\ue86c", // PowerShell 模块
		".bat":    "\uebc4", // Windows 批处理
		".cmd":    "\uebc4", // Windows 命令
		".vbs":    "\ueae9", // VBScript
		".asm":    "\ueae9", // 汇编
		".s":      "\ueae9", // 汇编
		".cs":     "\ue649", // C#
		".vb":     "\ueae9", // Visual Basic
		".fs":     "\ueae9", // F#
		".hs":     "\ueae9", // Haskell
		".dart":   "\ueae9", // Dart
		".ex":     "\ueae9", // Elixir
		".exs":    "\ueae9", // Elixir 脚本
		".groovy": "\ueae9", // Groovy

		// 模板文件
		".pug":      "\ueae9", // Pug
		".jade":     "\ueae9", // Jade
		".haml":     "\ueae9", // Haml
		".erb":      "\ueae9", // ERB
		".tpl":      "\ueae9", // 模板
		".hbs":      "\ueae9", // Handlebars
		".mustache": "\ueae9", // Mustache
		".ejs":      "\ueae9", // EJS

		// 其他代码相关
		".mdx":    "\ueae9", // MDX
		".svelte": "\ueae9", // Svelte
		".astro":  "\ueae9", // Astro

		// 新兴语言和框架
		".zig":     "\ueae9", // Zig语言
		".nim":     "\ueae9", // Nim语言
		".crystal": "\ueae9", // Crystal语言
		".v":       "\ue6ac", // V语言
		".odin":    "\ueae9", // Odin语言
		".gleam":   "\ueae9", // Gleam语言
		".roc":     "\ueae9", // Roc语言

		// Web开发相关
		".mjs":        "\ueae9", // ES模块
		".cjs":        "\ueae9", // CommonJS模块
		".coffee":     "\ueae9", // CoffeeScript
		".livescript": "\ueae9", // LiveScript
		".elm":        "\ueae9", // Elm语言
		".purescript": "\ueae9", // PureScript
		".reason":     "\ueae9", // ReasonML
		".rescript":   "\ueae9", // ReScript

		// 移动开发
		".xaml":       "\ueae9", // XAML标记
		".storyboard": "\ueae9", // iOS Storyboard
		".xib":        "\ueae9", // iOS XIB文件

		// 游戏开发
		".gd":   "\ueae9", // Godot脚本
		".hlsl": "\ueae9", // HLSL着色器
		".glsl": "\ueae9", // GLSL着色器

		// 配置文件格式
		".ini":        "\ue5fc",
		".conf":       "\ue5fc",
		".cfg":        "\ue5fc",
		".json":       "\ueb0f",
		".yaml":       "\ue8eb",
		".yml":        "\ue8eb",
		".xml":        "\ue5fc",
		".toml":       "\ue6b2",
		".env":        "\ue5fc",
		".properties": "\ue5fc",
		".config":     "\ue5fc",
		".settings":   "\ue5fc",

		// 文档文件
		".md":        "\uf48a", // Markdown
		".markdown":  "\uf48a", // Markdown
		".txt":       "\ue5fc", // 文本文件
		".readme":    "\ueda4", // README文件
		".license":   "\uf4d1", // LICENSE文件
		".changelog": "\ue5fc", // CHANGELOG文件

		// 项目配置文件
		".mod":                 "\ue5fc", // Go模块
		".sum":                 "\ue5fc", // Go校验和
		".lock":                "\uf023", // 依赖锁定
		".npmrc":               "\ue5fc", // npm配置
		".yarnrc":              "\ue5fc", // yarn配置
		".pnpmrc":              "\ue5fc", // pnpm配置
		".nvmrc":               "\ue5fc", // nvm配置
		".editorconfig":        "\ue5fc", // 编辑器配置
		".gitconfig":           "\ue5fc", // Git配置
		".gitignore":           "\ue5fc", // Git忽略
		".gitattributes":       "\ue5fc", // Git属性
		".dockerignore":        "\uf21f", // Docker忽略
		".dockerfile":          "\uf21f", // Dockerfile
		".docker-compose.yml":  "\uf21f", // Docker Compose
		".docker-compose.yaml": "\uf21f", // Docker Compose

		// 开发工具配置
		".eslintrc":         "\ue5fc", // ESLint
		".eslintrc.json":    "\ue5fc",
		".eslintrc.js":      "\ue5fc",
		".eslintrc.yml":     "\ue5fc",
		".eslintignore":     "\ue5fc",
		".prettierrc":       "\ue5fc", // Prettier
		".prettierrc.json":  "\ue5fc",
		".prettierrc.js":    "\ue5fc",
		".prettierrc.yml":   "\ue5fc",
		".prettierignore":   "\ue5fc",
		".stylelintrc":      "\ue5fc", // Stylelint
		".stylelintrc.json": "\ue5fc",
		".stylelintrc.js":   "\ue5fc",
		".stylelintrc.yml":  "\ue5fc",
		".stylelintignore":  "\ue5fc",
		".babelrc":          "\ue5fc", // Babel
		".babelrc.json":     "\ue5fc",
		".babelrc.js":       "\ue5fc",

		// 构建工具配置
		".webpack.config.js":  "\ue5fc", // Webpack
		".rollup.config.js":   "\ue5fc", // Rollup
		".vite.config.js":     "\ue5fc", // Vite
		".vite.config.ts":     "\ue5fc",
		".tsconfig.json":      "\ue5fc", // TypeScript
		".jsconfig.json":      "\ue5fc", // JavaScript
		".vue.config.js":      "\ue5fc", // Vue
		".nuxt.config.js":     "\ue5fc", // Nuxt
		".nuxt.config.ts":     "\ue5fc",
		".next.config.js":     "\ue5fc", // Next.js
		".gatsby-config.js":   "\ue5fc", // Gatsby
		".postcss.config.js":  "\ue5fc", // PostCSS
		".tailwind.config.js": "\ue5fc", // Tailwind
		".jest.config.js":     "\ue5fc", // Jest
		".cypress.json":       "\ue5fc", // Cypress

		// 环境配置
		".env.development":       "\ue5fc",
		".env.production":        "\ue5fc",
		".env.test":              "\ue5fc",
		".env.local":             "\ue5fc",
		".env.development.local": "\ue5fc",
		".env.production.local":  "\ue5fc",
		".env.test.local":        "\ue5fc",
		".env.example":           "\ue5fc",
		".env.sample":            "\ue5fc",
		".env.dist":              "\ue5fc",
		".env.template":          "\ue5fc",

		// 服务器配置
		".htaccess":   "\ue5fc", // Apache
		".nginx.conf": "\ue5fc", // Nginx
		".vhost":      "\ue5fc", // 虚拟主机
		".htpasswd":   "\ue5fc", // Apache认证

		// CI/CD配置
		".travis.yml":    "\ue5fc", // Travis CI
		".gitlab-ci.yml": "\ue5fc", // GitLab CI
		".jenkinsfile":   "\ue5fc", // Jenkins

		// 构建配置
		".makefile": "\ue5fc", // Makefile
		".mk":       "\ue5fc",
		".cmake":    "\ue5fc", // CMake

		// IDE配置
		".project":   "\ue5fc", // Eclipse
		".classpath": "\ue5fc",
		".iml":       "\ue5fc", // IntelliJ IDEA
		".ipr":       "\ue5fc",
		".iws":       "\ue5fc",
		".workspace": "\ue5fc", // Eclipse工作区

		// 其他配置
		".d.ts":         "\ue5fc", // TypeScript声明
		".swp":          "\ue5fc", // Vim交换文件
		".swo":          "\ue5fc",
		".swn":          "\ue5fc",
		".lck":          "\uf023", // 锁文件
		".pid":          "\ue5fc", // 进程ID
		".gdbinit":      "\ue5fc", // GDB
		".lldbinit":     "\ue5fc", // LLDB
		".clang-format": "\ue5fc", // Clang格式化
		".clang-tidy":   "\ue5fc", // Clang静态分析

		// 容器和虚拟化
		".k8s.yml":     "\ue81d", // Kubernetes配置
		".helm.yml":    "\ue7fb", // Helm配置
		".compose.yml": "\uf21f", // Docker Compose简写

		// 现代构建工具
		".swcrc":      "\ue5fc", // SWC配置
		".turbo.json": "\ue5fc", // Turbo配置
		".nx.json":    "\ue5fc", // Nx配置
		".rush.json":  "\ue5fc", // Rush配置

		// 云服务配置
		".serverless.yml": "\ue5fc", // Serverless配置
		".tf":             "\ue5fc", // Terraform文件
		".tfvars":         "\ue5fc", // Terraform变量

		// 包管理器
		".cargo.toml":   "\ue5fc", // Rust Cargo
		".pubspec.yaml": "\ue5fc", // Dart/Flutter
		".mix.exs":      "\ue5fc", // Elixir Mix
		".shard.yml":    "\ue5fc", // Crystal Shards

		// 文档格式
		".rst":  "\ue5fc", // reStructuredText
		".adoc": "\ue5fc", // AsciiDoc

		// 压缩文件
		".zip":     "\uf410", // ZIP
		".tar":     "\uf410", // TAR
		".gz":      "\uf410", // GZIP
		".bz2":     "\uf410", // BZIP2
		".rar":     "\uf410", // RAR
		".7z":      "\uf410", // 7-Zip
		".tar.gz":  "\uf410", // TAR.GZ
		".tar.bz2": "\uf410", // TAR.BZ2
		".tgz":     "\uf410", // TGZ
		".xz":      "\uf410", // XZ
		".lzma":    "\uf410", // LZMA
		".tar.xz":  "\uf410", // TAR.XZ
		".tbz2":    "\uf410", // TBZ2
		".tbz":     "\uf410", // TBZ
		".txz":     "\uf410", // TXZ
		".lz":      "\uf410", // LZ
		".lz4":     "\uf410", // LZ4
		".zst":     "\uf410", // Zstandard
		".zstd":    "\uf410", // Zstandard
		".br":      "\uf410", // Brotli
		".zlib":    "\uf410", // Zlib
		".wal":     "\ued1b", // python wal

		// 应用程序包
		".jar":  "\ue66d", // Java归档
		".war":  "\ue66d", // Java Web应用
		".ear":  "\ue66d", // Java企业应用
		".apk":  "\uf410", // Android应用
		".ipa":  "\uf410", // iOS应用
		".whl":  "\uf410", // Python wheel
		".rpm":  "\uf410", // Red Hat包
		".deb":  "\uf410", // Debian包
		".msi":  "\uf410", // Windows安装包
		".pkg":  "\uf410", // macOS安装包
		".dmg":  "\uf410", // macOS磁盘映像
		".appx": "\uf410", // Windows应用包

		// 数据库文件
		".db":      "\uf472", // SQLite
		".sqlite":  "\uf472", // SQLite
		".db3":     "\uf472", // SQLite 3
		".sqlite3": "\uf472", // SQLite 3
		".mdb":     "\uf472", // Access
		".accdb":   "\uf472", // Access 2007+
		".dbf":     "\uf472", // dBase
		".mdf":     "\uf472", // SQL Server数据
		".ldf":     "\uf472", // SQL Server日志
		".ndf":     "\uf472", // SQL Server辅助
		".fdb":     "\uf472", // Firebird
		".gdb":     "\uf472", // InterBase
		".ibd":     "\uf472", // MySQL索引数据
		".frm":     "\uf472", // MySQL表结构
		".myd":     "\uf472", // MySQL数据
		".myi":     "\uf472", // MySQL索引

		// 其他数据文件
		".csv":  "\ueefc", // CSV数据
		".xlsx": "\uf1c3", // Excel
		".xls":  "\uf1c3", // Excel旧版
		".pdf":  "\uf1c1", // PDF文档
		".doc":  "\uf1c2", // Word文档
		".docx": "\uf1c2", // Word文档
		".ppt":  "\uf1c4", // PowerPoint
		".pptx": "\uf1c4", // PowerPoint

		// 图片文件
		".jpg":  "\uf03e", // JPEG图片
		".jpeg": "\uf03e", // JPEG图片
		".png":  "\uf03e", // PNG图片
		".gif":  "\uf03e", // GIF图片
		".bmp":  "\uf03e", // BMP图片
		".ico":  "\uf03e", // 图标
		".svg":  "\uf03e", // SVG图片
		".webp": "\uf03e", // WebP图片
		".tif":  "\uf03e", // TIFF图片
		".tiff": "\uf03e", // TIFF图片
		".psd":  "\uf03e", // Photoshop文件
		".eps":  "\uf03e", // EPS文件
		".ai":   "\uf03e", // Adobe Illustrator文件
		".ps":   "\uf03e", // PostScript文件
		".rtf":  "\uf03e", // Rich Text Format

		// 视频文件
		".mp4":  "\uf52c", // MP4视频
		".avi":  "\uf52c", // AVI视频
		".mkv":  "\uf52c", // MKV视频
		".mov":  "\uf52c", // QuickTime视频
		".wmv":  "\uf52c", // Windows Media视频
		".webm": "\uf52c", // WebM视频
		".ogv":  "\uf52c", // Ogg视频
		".flv":  "\uf52c", // Flash视频
		".m4v":  "\uf52c", // MPEG-4视频
		".3gp":  "\uf52c", // 3GP视频

		// 音频文件
		".mp3":  "\uec1b", // MP3音频
		".wav":  "\uec1b", // WAV音频
		".flac": "\uec1b", // FLAC音频
		".ogg":  "\uec1b", // Ogg音频
		".m4a":  "\uec1b", // M4A音频
		".aac":  "\uec1b", // AAC音频
		".wma":  "\uec1b", // Windows Media音频
		".opus": "\uec1b", // Opus音频
		".aiff": "\uec1b", // AIFF音频
		".au":   "\uec1b", // AU音频

		// 字体文件
		".ttf":   "\ue659", // TrueType字体
		".otf":   "\ue659", // OpenType字体
		".woff":  "\ue659", // Web字体
		".woff2": "\ue659", // Web字体2
		".eot":   "\ue659", // Embedded OpenType字体

		// 动态链接库
		".so":    "\ueb9c", // Linux共享对象
		".dll":   "\ueb9c", // Windows动态链接库
		".dylib": "\ueb9c", // macOS动态链接库

		// 静态库
		".a":   "\ueb9c", // Linux静态库
		".lib": "\ueb9c", // Windows静态库

		// 可执行文件
		".exe": "\uebc4", // Windows可执行文件
		".com": "\uebc4", // DOS命令文件
		".bin": "\ueb9c", // 二进制文件
		".elf": "\ueb9c", // Linux可执行文件
		".out": "\ueb9c", // 编译输出文件

		// 目标文件
		".o":   "\ueb9c", // 目标文件
		".obj": "\ueb9c", // Windows目标文件

		// 调试和符号文件
		".pdb": "\ueb9c", // Windows调试符号
		".exp": "\ueb9c", // Windows导出文件
		".ilk": "\ueb9c", // Windows增量链接
		".pch": "\ueb9c", // 预编译头文件
		".sbr": "\ueb9c", // 源代码浏览器
		".idb": "\ueb9c", // Visual Studio增量调试

		// 字节码文件
		".pyc":   "\ued1b", // Python字节码
		".pyo":   "\ued1b", // Python优化字节码
		".pyd":   "\ued1b", // Python扩展模块
		".class": "\ue66d", // Java字节码

		// 包文件
		".egg":  "\ued1b", // Python包
		".jmod": "\ue66d", // Java模块

		// 现代编译产物
		".wasm":        "\ueb9c", // WebAssembly
		".bc":          "\ueb9c", // LLVM位码
		".ll":          "\ueb9c", // LLVM中间表示
		".node":        "\ueb9c", // Node.js原生插件
		".rlib":        "\ueb9c", // Rust库
		".swiftmodule": "\ueb9c", // Swift模块

		// 映射和列表文件
		".map": "\ueb9c", // 映射文件
		".lst": "\ueb9c", // 列表文件
		".d":   "\ueb9c", // 依赖文件

		// 现代编译产物
		".dSYM":     "\ueb9c", // macOS调试符号
		".vsix":     "\ueb9c", // VS Code扩展
		".nupkg":    "\ueb9c", // NuGet包
		".gem":      "\ueb9c", // Ruby Gem
		".crate":    "\ueb9c", // Rust Crate
		".wheel":    "\ueb9c", // Python Wheel
		".snap":     "\ueb9c", // Snap包
		".flatpak":  "\ueb9c", // Flatpak包
		".appimage": "\ueb9c", // AppImage包
	},

	Default: "\uf4d3", // 默认图标
}

// getIcon 根据文件信息返回图标编码。
//
// 规则：
//   - 优先按普通文件类型映射，若未命中则按文件扩展名映射，最后使用默认图标。
//
// 参数：
//   - info: 文件信息结构体，包含文件类型、扩展名等。
//
// 返回值：
//   - string: 图标编码字符串，若不存在图标则为空字符串。
//
// 注意：
//   - 仅在存在图标时在末尾追加一个空格，避免无图标时产生多余空格。
func getIcon(info FileInfo) string {
	// 先确定图标
	var icon string
	if info.EntryType == FileType || info.EntryType == ExecutableType {
		if v, ok := DefaultIcons.ByExt[info.FileExt]; ok {
			icon = v
		} else {
			icon = DefaultIcons.ByType[info.EntryType]
		}
	} else {
		icon = DefaultIcons.ByType[info.EntryType]
	}

	// 若未匹配到或映射值为空，回退到默认图标
	if icon == "" {
		icon = DefaultIcons.Default
	}

	// 有图标则追加空格返回；没有则返回空字符串
	if icon != "" {
		return icon + " "
	}

	return ""
}
