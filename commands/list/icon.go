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
	// 按特殊文件名映射
	BySpecial map[string]string
	// 默认图标
	Default string
}

// defaultIcons 默认图标映射表
// 注：编码均为字符串，可直接拼接到名称前；是否加空格/着色由调用方决定。
var defaultIcons = IconMap{
	ByType: map[EntryType]string{
		DirType:         "\uf115", // 目录
		SymlinkType:     "\uF08E", // 软链接
		SocketType:      "\uf4d6", // 套接字
		PipeType:        "\uf4d6", // 管道
		BlockDeviceType: "\uf4d6", // 块设备
		CharDeviceType:  "\uf4d6", // 字符设备
		ExecutableType:  "\uf489", // 可执行文件
		EmptyType:       "\uf4a5", // 空文件
		FileType:        "\uf4a5", // 普通文件
		UnknownType:     "\uf4a5", // 未知类型
	},

	ByExt: map[string]string{
		// Go语言
		".go": "\uE627",

		// Python
		".py":  "\uED1B",
		".pyw": "\uED1B",

		// 脚本
		".sh":   "\uf489",
		".bash": "\uf489",
		".zsh":  "\uf489",
		".fsh":  "\uf489",
		".fish": "\uf489",
		".run":  "\uf489",

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
		".rb":     "\uf1c9", // Ruby
		".php":    "\uE73D", // PHP
		".java":   "\uec15", // Java
		".c":      "\ue649", // C
		".cpp":    "\ue649", // C++
		".h":      "\ue649", // C/C++ 头文件
		".hpp":    "\ue649", // C++ 头文件
		".m":      "\uf1c9", // Objective-C
		".mm":     "\uf1c9", // Objective-C++
		".swift":  "\uf1c9", // Swift
		".kt":     "\uf1c9", // Kotlin
		".kts":    "\uf1c9", // Kotlin 脚本
		".scala":  "\uf1c9", // Scala
		".lua":    "\uf1c9", // Lua
		".pl":     "\uf1c9", // Perl
		".pm":     "\uf1c9", // Perl 模块
		".r":      "\uf1c9", // R语言
		".sql":    "\uf472", // SQL
		".ps1":    "\uf489", // PowerShell
		".psm1":   "\uf489", // PowerShell 模块
		".bat":    "\uf489", // Windows 批处理
		".cmd":    "\uf489", // Windows 命令
		".vbs":    "\uf1c9", // VBScript
		".asm":    "\uf1c9", // 汇编
		".s":      "\uf1c9", // 汇编
		".cs":     "\ue649", // C#
		".vb":     "\uf1c9", // Visual Basic
		".fs":     "\uf1c9", // F#
		".hs":     "\uf1c9", // Haskell
		".dart":   "\uf1c9", // Dart
		".ex":     "\uf1c9", // Elixir
		".exs":    "\uf1c9", // Elixir 脚本
		".groovy": "\uf1c9", // Groovy

		// 模板文件
		".pug":      "\uf1c9", // Pug
		".jade":     "\uf1c9", // Jade
		".haml":     "\uf1c9", // Haml
		".erb":      "\uf1c9", // ERB
		".tpl":      "\uf1c9", // 模板
		".hbs":      "\uf1c9", // Handlebars
		".mustache": "\uf1c9", // Mustache
		".ejs":      "\uf1c9", // EJS

		// 其他代码相关
		".mdx":    "\uf1c9", // MDX
		".svelte": "\uf1c9", // Svelte
		".astro":  "\uf1c9", // Astro

		// 新兴语言和框架
		".zig":     "\uf1c9", // Zig语言
		".nim":     "\uf1c9", // Nim语言
		".crystal": "\uf1c9", // Crystal语言
		".v":       "\ue6ac", // V语言
		".odin":    "\uf1c9", // Odin语言
		".gleam":   "\uf1c9", // Gleam语言
		".roc":     "\uf1c9", // Roc语言

		// Web开发相关
		".mjs":        "\uf1c9", // ES模块
		".cjs":        "\uf1c9", // CommonJS模块
		".coffee":     "\uf1c9", // CoffeeScript
		".livescript": "\uf1c9", // LiveScript
		".elm":        "\uf1c9", // Elm语言
		".purescript": "\uf1c9", // PureScript
		".reason":     "\uf1c9", // ReasonML
		".rescript":   "\uf1c9", // ReScript

		// 移动开发
		".xaml":       "\uf1c9", // XAML标记
		".storyboard": "\uf1c9", // iOS Storyboard
		".xib":        "\uf1c9", // iOS XIB文件

		// 游戏开发
		".gd":   "\uf1c9", // Godot脚本
		".hlsl": "\uf1c9", // HLSL着色器
		".glsl": "\uf1c9", // GLSL着色器

		// 配置文件格式
		".ini":        "\ue5fc",
		".conf":       "\ue5fc",
		".cfg":        "\ue5fc",
		".json":       "\ue60b",
		".yaml":       "\ue6a8",
		".yml":        "\ue6a8",
		".xml":        "\ue5fc",
		".toml":       "\ue6b2",
		".env":        "\ue5fc",
		".properties": "\ue5fc",
		".config":     "\ue5fc",
		".settings":   "\ue5fc",
		".service":    "\ue5fc",

		// 文档文件
		".md":        "\uf48a", // Markdown
		".markdown":  "\uf48a", // Markdown
		".txt":       "\ue5fc", // 文本文件
		".readme":    "\ueda4", // README文件
		".license":   "\uf4d1", // LICENSE文件
		".lic":       "\uf4d1", // LIC文件
		".crt":       "\uf4d1", // 证书文件
		".cer":       "\uf4d1", // 证书文件
		".pem":       "\uf4d1", // 证书文件
		".changelog": "\ue5fc", // CHANGELOG文件

		// 项目配置文件
		".mod":           "\ue5fc", // Go模块
		".sum":           "\ue5fc", // Go校验和
		".lock":          "\uf023", // 依赖锁定
		".npmrc":         "\ue5fc", // npm配置
		".yarnrc":        "\ue5fc", // yarn配置
		".pnpmrc":        "\ue5fc", // pnpm配置
		".nvmrc":         "\ue5fc", // nvm配置
		".editorconfig":  "\ue5fc", // 编辑器配置
		".gitconfig":     "\ue5fc", // Git配置
		".gitignore":     "\ue5fc", // Git忽略
		".gitattributes": "\ue5fc", // Git属性
		".dockerignore":  "\uf21f", // Docker忽略
		".dockerfile":    "\uf21f", // Dockerfile

		// 开发工具配置
		".eslintrc":        "\ue5fc", // ESLint
		".eslintignore":    "\ue5fc",
		".prettierrc":      "\ue5fc", // Prettier
		".prettierignore":  "\ue5fc",
		".stylelintrc":     "\ue5fc", // Stylelint
		".stylelintignore": "\ue5fc",
		".babelrc":         "\ue5fc", // Babel

		// 服务器配置
		".htaccess": "\ue5fc", // Apache
		".vhost":    "\ue5fc", // 虚拟主机
		".htpasswd": "\ue5fc", // Apache认证

		// CI/CD配置
		".jenkinsfile": "\ue5fc", // Jenkins

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

		// 现代构建工具
		".swcrc": "\ue5fc", // SWC配置

		// 云服务配置
		".tf":     "\ue5fc", // Terraform文件
		".tfvars": "\ue5fc", // Terraform变量

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
		".jar":  "\uec15", // Java归档
		".war":  "\uec15", // Java Web应用
		".ear":  "\uec15", // Java企业应用
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
		".dump":    "\uf472", // sql备份"

		// 其他数据文件
		".csv":  "\ueefc", // CSV数据
		".xlsx": "\uf1c3", // Excel
		".xls":  "\uf1c3", // Excel旧版
		".pdf":  "\uf1c1", // PDF文档
		".doc":  "\uf1c2", // Word文档
		".docx": "\uf1c2", // Word文档
		".ppt":  "\uf1c4", // PowerPoint
		".pptx": "\uf1c4", // PowerPoint
		".tsv":  "\ueefc", // TSV数据 (Tab Separated Values)
		".odt":  "\uf1c2", // OpenDocument Text (与Word文档图标一致)
		".ods":  "\uf1c3", // OpenDocument Spreadsheet (与Excel图标一致)
		".odp":  "\uf1c4", // OpenDocument Presentation (与PowerPoint图标一致)
		".rtf":  "\uf1c2", // rtf文档

		// 图片文件
		".jpg":  "\uf03e", // JPEG图片
		".jpeg": "\uf03e", // JPEG图片
		".png":  "\uf03e", // PNG图片
		".gif":  "\uf03e", // GIF图片
		".bmp":  "\uf03e", // BMP图片
		".ico":  "\uf03e", // 图标
		".svg":  "\uf03e", // SVG图片
		".webp": "\uf03e", // WebP图片
		".heic": "\uf03e", // HEIC图片
		".avif": "\uf03e", // AVIF图片
		".raw":  "\uf03e", // 相机原始图片
		".cr2":  "\uf03e", // Canon Raw
		".nef":  "\uf03e", // Nikon Raw
		".tif":  "\uf03e", // TIFF图片
		".tiff": "\uf03e", // TIFF图片
		".psd":  "\uf03e", // Photoshop文件
		".eps":  "\uf03e", // EPS文件
		".ai":   "\uf03e", // Adobe Illustrator文件
		".ps":   "\uf03e", // PostScript文件
		".helf": "\uf03e", // HEIF图片

		// 视频文件
		".mp4":  "\uf03d", // MP4视频
		".avi":  "\uf03d", // AVI视频
		".mkv":  "\uf03d", // MKV视频
		".mov":  "\uf03d", // QuickTime视频
		".wmv":  "\uf03d", // Windows Media视频
		".webm": "\uf03d", // WebM视频
		".ogv":  "\uf03d", // Ogg视频
		".flv":  "\uf03d", // Flash视频
		".m4v":  "\uf03d", // MPEG-4视频
		".3gp":  "\uf03d", // 3GP视频
		".vob":  "\uf03d", // DVD Video Object
		".mpg":  "\uf03d", // MPEG
		".mpeg": "\uf03d", // MPEG
		".rm":   "\uf03d", // RealMedia
		".rmvb": "\uf03d", // RealMedia Variable Bitrate
		".asf":  "\uf03d", // Advanced Systems Format
		".m3u8": "\uf03d", // M3U8视频

		// 音频文件
		".mp3":  "\ue638", // MP3音频
		".wav":  "\ue638", // WAV音频
		".flac": "\ue638", // FLAC音频
		".ogg":  "\ue638", // Ogg音频
		".m4a":  "\ue638", // M4A音频
		".aac":  "\ue638", // AAC音频
		".wma":  "\ue638", // Windows Media音频
		".opus": "\ue638", // Opus音频
		".aiff": "\ue638", // AIFF音频
		".au":   "\ue638", // AU音频
		".mid":  "\ue638", // MIDI文件
		".midi": "\ue638", // MIDI文件
		".ape":  "\ue638", // Monkey's Audio (无损)
		".wv":   "\ue638", // WavPack (无损)
		".dsf":  "\ue638", // DSD Stream File
		".dff":  "\ue638", // DSD Interchange File Format

		// 字体文件
		".ttf":   "\ue659", // TrueType字体
		".otf":   "\ue659", // OpenType字体
		".woff":  "\ue659", // Web字体
		".woff2": "\ue659", // Web字体2
		".eot":   "\ue659", // Embedded OpenType字体

		// 动态链接库
		".so":    "\uf471", // Linux共享对象
		".dll":   "\uf471", // Windows动态链接库
		".dylib": "\uf471", // macOS动态链接库

		// 静态库
		".a":   "\uf471", // Linux静态库
		".lib": "\uf471", // Windows静态库

		// 可执行文件
		".exe": "\uf489", // Windows可执行文件
		".com": "\uf489", // DOS命令文件
		".bin": "\uf471", // 二进制文件
		".elf": "\uf471", // Linux可执行文件
		".out": "\uf471", // 编译输出文件

		// 目标文件
		".o":   "\uf471", // 目标文件
		".obj": "\uf471", // Windows目标文件

		// 调试和符号文件
		".pdb": "\uf471", // Windows调试符号
		".exp": "\uf471", // Windows导出文件
		".ilk": "\uf471", // Windows增量链接
		".pch": "\uf471", // 预编译头文件
		".sbr": "\uf471", // 源代码浏览器
		".idb": "\uf471", // Visual Studio增量调试

		// 字节码文件
		".pyc":   "\ued1b", // Python字节码
		".pyo":   "\ued1b", // Python优化字节码
		".pyd":   "\ued1b", // Python扩展模块
		".class": "\uec15", // Java字节码

		// 包文件
		".egg":  "\ued1b", // Python包
		".jmod": "\uec15", // Java模块

		// 现代编译产物
		".wasm":        "\ue8e0", // WebAssembly
		".bc":          "\uf471", // LLVM位码
		".ll":          "\uf471", // LLVM中间表示
		".node":        "\uf471", // Node.js原生插件
		".rlib":        "\uf471", // Rust库
		".swiftmodule": "\uf471", // Swift模块

		// 映射和列表文件
		".map": "\uf471", // 映射文件
		".lst": "\uf471", // 列表文件
		".d":   "\uf471", // 依赖文件

		// 现代编译产物
		".dSYM":     "\uf471", // macOS调试符号
		".vsix":     "\uf471", // VS Code扩展
		".nupkg":    "\uf471", // NuGet包
		".gem":      "\uf471", // Ruby Gem
		".crate":    "\uf471", // Rust Crate
		".wheel":    "\uf471", // Python Wheel
		".snap":     "\uf471", // Snap包
		".flatpak":  "\uf471", // Flatpak包
		".appimage": "\uf471", // AppImage包

		// 日志文件
		".log": "\uf4ed", // 日志文件
	},

	Default: "\uf4a5", // 默认图标

	BySpecial: map[string]string{
		"Makefile":            "\ue673", // Makefile文件
		"makefile":            "\ue673", // Makefile文件
		"Dockerfile":          "\uf21f", // Dockerfile文件
		"docker-compose.yml":  "\uf21f", // Docker Compose文件
		"docker-compose.yaml": "\uf21f", // Docker Compose文件
		"compose.yml":         "\uf21f", // Docker Compose v2 推荐简写
		"compose.yaml":        "\uf21f", // Docker Compose v2 推荐简写
		"LICENSE":             "\uf4d1", // 许可证文件
		"LICENCE":             "\uf4d1", // 许可证文件
		"Jenkinsfile":         "\uf2ec", // Jenkinsfile文件
		"README.md":           "\uEDA4", // README文件
		"API.md":              "\uEDA4", // API文档文件
		"APIDOC.md":           "\uEDA4", // API文档文件
		"APIDoc.md":           "\uEDA4", // API文档文件
		".env":                "\ue5fc", // 环境变量文件
		".gitignore":          "\uf1d3", // Git忽略文件
		".gitattributes":      "\uf1d3", // Git属性配置
		"access.log":          "\uE776", // 访问日志
		"error.log":           "\uE776", // 错误日志
		"nginx.conf":          "\uE776", // Nginx配置文件
		"Taskfile.yml":        "\uf0ae", // 任务文件
		"taskfile.yml":        "\uf0ae", // 任务文件
		"Taskfile.yaml":       "\uf0ae", // 任务文件
		"taskfile.yaml":       "\uf0ae", // 任务文件
		"Taskfile.dist.yml":   "\uf0ae", // 默认任务文件
		"taskfile.dist.yml":   "\uf0ae", // 默认任务文件
		"Taskfile.dist.yaml":  "\uf0ae", // 默认任务文件
		"taskfile.dist.yaml":  "\uf0ae", // 默认任务文件
		"gob.toml":            "\uf0ae", // gob配置文件
		"go.mod":              "\uE627", // go模块文件
		"go.sum":              "\uE627", // go模块文件
		".bashrc":             "\ue760", // bash配置文件
		".bash_profile":       "\ue760", // bash配置文件
		".bash_history":       "\ue760", // bash历史文件
	},
}

// getIcon 根据文件信息返回图标编码。
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
	var icon string

	// 1. 最优先：按特殊文件名匹配
	//    例如 "Makefile", "Dockerfile", "README", "LICENSE" 等
	if v, ok := defaultIcons.BySpecial[info.Name]; ok {
		return addSpace(v) // 找到即返回
	}

	// 2. 其次：按文件类型匹配
	//    这包括目录、软链接等非普通文件，以及普通文件的通用类型图标
	if v, ok := defaultIcons.ByType[info.EntryType]; ok {
		// 对于非普通文件，ByType 的图标通常是最终的
		if info.EntryType != FileType {
			return addSpace(v) // 找到即返回
		}

		// 如果是普通文件，先暂存ByType的通用图标，继续尝试ByExt
		icon = v
	}

	// 3. 再次：按文件扩展名匹配 (仅对普通文件有效)
	//    如果文件是普通文件，并且之前没有找到更具体的图标（例如特殊文件名），
	//    则尝试通过扩展名来获取更具体的图标。
	if info.EntryType == FileType {
		if v, ok := defaultIcons.ByExt[info.FileExt]; ok {
			return addSpace(v) // 找到即返回
		}
	}

	// 4. 如果前面都没有匹配到，但ByType为普通文件提供了通用图标，则使用它
	if icon != "" {
		return addSpace(icon)
	}

	// 5. 最后：回退到默认图标
	return addSpace(defaultIcons.Default)
}

// addSpace 用于给图标添加空格的函数
//
// 参数：
//   - icon: 图标编码字符串，若为空则不添加空格。
//
// 返回值：
//   - string: 图标编码字符串，若为空则返回空字符串。
//
// 注意：
//   - 仅在图标非空时在末尾追加一个空格，避免无图标时产生多余空格。
func addSpace(icon string) string {
	if icon != "" {
		return icon + " "
	}

	return ""
}
