// Package list 扩展名映射表定义
// 该文件定义了各种文件类型的扩展名映射表，用于颜色分类
package list

// 特殊配置文件名映射表（无扩展名）
var specialConfigFiles = map[string]bool{
	"makefile":       true,
	"dockerfile":     true,
	"vagrantfile":    true,
	"rakefile":       true,
	"gemfile":        true,
	"podfile":        true,
	"brewfile":       true,
	"procfile":       true,
	"readme":         true,
	"license":        true,
	"changelog":      true,
	"authors":        true,
	"contributors":   true,
	"copying":        true,
	"install":        true,
	"news":           true,
	"todo":           true,
	"version":        true,
	"cmakelists.txt": true,

	// 现代项目文件
	"justfile":    true, // Just命令运行器
	"taskfile":    true, // Task运行器
	"flake.nix":   true, // Nix Flakes
	"default.nix": true, // Nix表达式
	"shell.nix":   true, // Nix Shell
	"deno.json":   true, // Deno配置
	"bun.lockb":   true, // Bun锁文件

	// CI/CD文件
	"buildkite.yml":           true, // Buildkite
	"azure-pipelines.yml":     true, // Azure DevOps
	"bitbucket-pipelines.yml": true, // Bitbucket
}

// greenExtensions 绿色系文件扩展名映射表
// 这些文件使用绿色显示
var greenExtensions = map[string]bool{
	// Go语言
	".go": true,

	// Python
	".py":  true,
	".pyw": true,

	// Shell脚本
	".sh":   true,
	".bash": true,
	".zsh":  true,

	// JavaScript/TypeScript
	".js":  true,
	".ts":  true,
	".jsx": true,
	".tsx": true,

	// Web前端
	".html": true,
	".css":  true,
	".scss": true,
	".sass": true,
	".less": true,
	".vue":  true,

	// 其他编程语言
	".rs":     true, // Rust
	".rb":     true, // Ruby
	".php":    true, // PHP
	".java":   true, // Java
	".c":      true, // C
	".cpp":    true, // C++
	".h":      true, // C/C++ 头文件
	".hpp":    true, // C++ 头文件
	".m":      true, // Objective-C
	".mm":     true, // Objective-C++
	".swift":  true, // Swift
	".kt":     true, // Kotlin
	".kts":    true, // Kotlin 脚本
	".scala":  true, // Scala
	".lua":    true, // Lua
	".pl":     true, // Perl
	".pm":     true, // Perl 模块
	".r":      true, // R语言
	".sql":    true, // SQL
	".ps1":    true, // PowerShell
	".psm1":   true, // PowerShell 模块
	".bat":    true, // Windows 批处理
	".cmd":    true, // Windows 命令
	".vbs":    true, // VBScript
	".asm":    true, // 汇编
	".s":      true, // 汇编
	".cs":     true, // C#
	".vb":     true, // Visual Basic
	".fs":     true, // F#
	".hs":     true, // Haskell
	".dart":   true, // Dart
	".ex":     true, // Elixir
	".exs":    true, // Elixir 脚本
	".groovy": true, // Groovy

	// 模板文件
	".pug":      true, // Pug
	".jade":     true, // Jade
	".haml":     true, // Haml
	".erb":      true, // ERB
	".tpl":      true, // 模板
	".hbs":      true, // Handlebars
	".mustache": true, // Mustache
	".ejs":      true, // EJS

	// 其他代码相关
	".mdx":    true, // MDX
	".svelte": true, // Svelte
	".astro":  true, // Astro

	// 新兴语言和框架
	".zig":     true, // Zig语言
	".nim":     true, // Nim语言
	".crystal": true, // Crystal语言
	".v":       true, // V语言
	".odin":    true, // Odin语言
	".gleam":   true, // Gleam语言
	".roc":     true, // Roc语言

	// Web开发相关
	".mjs":        true, // ES模块
	".cjs":        true, // CommonJS模块
	".coffee":     true, // CoffeeScript
	".livescript": true, // LiveScript
	".elm":        true, // Elm语言
	".purescript": true, // PureScript
	".reason":     true, // ReasonML
	".rescript":   true, // ReScript

	// 移动开发
	".xaml":       true, // XAML标记
	".storyboard": true, // iOS Storyboard
	".xib":        true, // iOS XIB文件

	// 游戏开发
	".gd":   true, // Godot脚本
	".hlsl": true, // HLSL着色器
	".glsl": true, // GLSL着色器
}

// yellowExtensions 黄色系文件扩展名映射表
// 这些文件使用黄色显示
var yellowExtensions = map[string]bool{
	// 配置文件格式
	".ini":        true,
	".conf":       true,
	".cfg":        true,
	".json":       true,
	".yaml":       true,
	".yml":        true,
	".xml":        true,
	".toml":       true,
	".env":        true,
	".properties": true,
	".config":     true,
	".settings":   true,

	// 文档文件
	".md":        true, // Markdown
	".markdown":  true, // Markdown
	".txt":       true, // 文本文件
	".readme":    true, // README文件
	".license":   true, // LICENSE文件
	".changelog": true, // CHANGELOG文件

	// 项目配置文件
	".mod":                 true, // Go模块
	".sum":                 true, // Go校验和
	".lock":                true, // 依赖锁定
	".npmrc":               true, // npm配置
	".yarnrc":              true, // yarn配置
	".pnpmrc":              true, // pnpm配置
	".nvmrc":               true, // nvm配置
	".editorconfig":        true, // 编辑器配置
	".gitconfig":           true, // Git配置
	".gitignore":           true, // Git忽略
	".gitattributes":       true, // Git属性
	".dockerignore":        true, // Docker忽略
	".dockerfile":          true, // Dockerfile
	".docker-compose.yml":  true, // Docker Compose
	".docker-compose.yaml": true, // Docker Compose

	// 开发工具配置
	".eslintrc":         true, // ESLint
	".eslintrc.json":    true,
	".eslintrc.js":      true,
	".eslintrc.yml":     true,
	".eslintignore":     true,
	".prettierrc":       true, // Prettier
	".prettierrc.json":  true,
	".prettierrc.js":    true,
	".prettierrc.yml":   true,
	".prettierignore":   true,
	".stylelintrc":      true, // Stylelint
	".stylelintrc.json": true,
	".stylelintrc.js":   true,
	".stylelintrc.yml":  true,
	".stylelintignore":  true,
	".babelrc":          true, // Babel
	".babelrc.json":     true,
	".babelrc.js":       true,

	// 构建工具配置
	".webpack.config.js":  true, // Webpack
	".rollup.config.js":   true, // Rollup
	".vite.config.js":     true, // Vite
	".vite.config.ts":     true,
	".tsconfig.json":      true, // TypeScript
	".jsconfig.json":      true, // JavaScript
	".vue.config.js":      true, // Vue
	".nuxt.config.js":     true, // Nuxt
	".nuxt.config.ts":     true,
	".next.config.js":     true, // Next.js
	".gatsby-config.js":   true, // Gatsby
	".postcss.config.js":  true, // PostCSS
	".tailwind.config.js": true, // Tailwind
	".jest.config.js":     true, // Jest
	".cypress.json":       true, // Cypress

	// 环境配置
	".env.development":       true,
	".env.production":        true,
	".env.test":              true,
	".env.local":             true,
	".env.development.local": true,
	".env.production.local":  true,
	".env.test.local":        true,
	".env.example":           true,
	".env.sample":            true,
	".env.dist":              true,
	".env.template":          true,

	// 服务器配置
	".htaccess":   true, // Apache
	".nginx.conf": true, // Nginx
	".vhost":      true, // 虚拟主机
	".htpasswd":   true, // Apache认证

	// CI/CD配置
	".travis.yml":    true, // Travis CI
	".gitlab-ci.yml": true, // GitLab CI
	".jenkinsfile":   true, // Jenkins

	// 构建配置
	".makefile": true, // Makefile
	".mk":       true,
	".cmake":    true, // CMake

	// IDE配置
	".project":   true, // Eclipse
	".classpath": true,
	".iml":       true, // IntelliJ IDEA
	".ipr":       true,
	".iws":       true,
	".workspace": true, // Eclipse工作区

	// 其他配置
	".d.ts":         true, // TypeScript声明
	".swp":          true, // Vim交换文件
	".swo":          true,
	".swn":          true,
	".lck":          true, // 锁文件
	".pid":          true, // 进程ID
	".gdbinit":      true, // GDB
	".lldbinit":     true, // LLDB
	".clang-format": true, // Clang格式化
	".clang-tidy":   true, // Clang静态分析

	// 容器和虚拟化
	".k8s.yml":     true, // Kubernetes配置
	".helm.yml":    true, // Helm配置
	".compose.yml": true, // Docker Compose简写

	// 现代构建工具
	".swcrc":      true, // SWC配置
	".turbo.json": true, // Turbo配置
	".nx.json":    true, // Nx配置
	".rush.json":  true, // Rush配置

	// 云服务配置
	".serverless.yml": true, // Serverless配置
	".tf":             true, // Terraform文件
	".tfvars":         true, // Terraform变量

	// 包管理器
	".cargo.toml":   true, // Rust Cargo
	".pubspec.yaml": true, // Dart/Flutter
	".mix.exs":      true, // Elixir Mix
	".shard.yml":    true, // Crystal Shards

	// 文档格式
	".rst":  true, // reStructuredText
	".adoc": true, // AsciiDoc
}

// redExtensions 红色系文件扩展名映射表
// 这些文件使用红色显示
var redExtensions = map[string]bool{
	// 压缩文件
	".zip":     true, // ZIP
	".tar":     true, // TAR
	".gz":      true, // GZIP
	".bz2":     true, // BZIP2
	".rar":     true, // RAR
	".7z":      true, // 7-Zip
	".tar.gz":  true, // TAR.GZ
	".tar.bz2": true, // TAR.BZ2
	".tgz":     true, // TGZ
	".xz":      true, // XZ
	".lzma":    true, // LZMA
	".tar.xz":  true, // TAR.XZ
	".tbz2":    true, // TBZ2
	".tbz":     true, // TBZ
	".txz":     true, // TXZ
	".lz":      true, // LZ
	".lz4":     true, // LZ4
	".zst":     true, // Zstandard
	".zstd":    true, // Zstandard
	".br":      true, // Brotli
	".zlib":    true, // Zlib
	".wal":     true, // python wal

	// 应用程序包
	".jar":  true, // Java归档
	".war":  true, // Java Web应用
	".ear":  true, // Java企业应用
	".apk":  true, // Android应用
	".ipa":  true, // iOS应用
	".whl":  true, // Python wheel
	".rpm":  true, // Red Hat包
	".deb":  true, // Debian包
	".msi":  true, // Windows安装包
	".pkg":  true, // macOS安装包
	".dmg":  true, // macOS磁盘映像
	".appx": true, // Windows应用包

	// 数据库文件
	".db":      true, // SQLite
	".sqlite":  true, // SQLite
	".db3":     true, // SQLite 3
	".sqlite3": true, // SQLite 3
	".mdb":     true, // Access
	".accdb":   true, // Access 2007+
	".dbf":     true, // dBase
	".mdf":     true, // SQL Server数据
	".ldf":     true, // SQL Server日志
	".ndf":     true, // SQL Server辅助
	".fdb":     true, // Firebird
	".gdb":     true, // InterBase
	".ibd":     true, // MySQL索引数据
	".frm":     true, // MySQL表结构
	".myd":     true, // MySQL数据
	".myi":     true, // MySQL索引

	// 磁盘映像
	".iso":   true, // 光盘映像
	".img":   true, // 磁盘映像
	".vhd":   true, // 虚拟硬盘
	".vmdk":  true, // VMware虚拟磁盘
	".qcow2": true, // QEMU磁盘映像
	".raw":   true, // 原始磁盘映像

	// 备份文件
	".bak":    true, // 备份文件
	".bkf":    true, // Windows备份
	".backup": true, // 备份文件

	// 其他数据文件
	".csv":  true, // CSV数据
	".tsv":  true, // TSV数据
	".xlsx": true, // Excel
	".xls":  true, // Excel旧版
	".ods":  true, // OpenDocument表格
	".pdf":  true, // PDF文档
	".doc":  true, // Word文档
	".docx": true, // Word文档
	".ppt":  true, // PowerPoint
	".pptx": true, // PowerPoint
	".odt":  true, // OpenDocument文本
	".odp":  true, // OpenDocument演示
	// 图片文件
	".jpg":  true, // JPEG图片
	".jpeg": true, // JPEG图片
	".png":  true, // PNG图片
	".gif":  true, // GIF图片
	".bmp":  true, // BMP图片
	".ico":  true, // 图标
	".svg":  true, // SVG图片
	".webp": true, // WebP图片
	".tif":  true, // TIFF图片
	".tiff": true, // TIFF图片
	".psd":  true, // Photoshop文件
	".eps":  true, // EPS文件
	".ai":   true, // Adobe Illustrator文件
	".ps":   true, // PostScript文件
	".rtf":  true, // Rich Text Format

	// 视频文件
	".mp4":  true, // MP4视频
	".avi":  true, // AVI视频
	".mkv":  true, // MKV视频
	".mov":  true, // QuickTime视频
	".wmv":  true, // Windows Media视频
	".webm": true, // WebM视频
	".ogv":  true, // Ogg视频
	".flv":  true, // Flash视频
	".m4v":  true, // MPEG-4视频
	".3gp":  true, // 3GP视频

	// 音频文件
	".mp3":  true, // MP3音频
	".wav":  true, // WAV音频
	".flac": true, // FLAC音频
	".ogg":  true, // Ogg音频
	".m4a":  true, // M4A音频
	".aac":  true, // AAC音频
	".wma":  true, // Windows Media音频
	".opus": true, // Opus音频
	".aiff": true, // AIFF音频
	".au":   true, // AU音频

	// 字体文件
	".ttf":   true, // TrueType字体
	".otf":   true, // OpenType字体
	".woff":  true, // Web字体
	".woff2": true, // Web字体2
	".eot":   true, // Embedded OpenType字体

	// 3D和设计文件
	".blend":  true, // Blender文件
	".fbx":    true, // FBX 3D模型
	".obj":    true, // Wavefront OBJ
	".dae":    true, // COLLADA 3D模型
	".sketch": true, // Sketch设计文件
	".fig":    true, // Figma文件
	".xd":     true, // Adobe XD文件
	".vsd":    true, // Visio文件
	".vsdx":   true, // Visio文件
	".xmind":  true, // XMind文件"

	// 科学数据
	".hdf5":   true, // HDF5数据
	".h5":     true, // HDF5数据
	".nc":     true, // NetCDF数据
	".netcdf": true, // NetCDF数据
	".mat":    true, // MATLAB数据
}

// magentaExtensions 品红系文件扩展名映射表
// 这些文件使用紫色显示
var magentaExtensions = map[string]bool{
	// 动态链接库
	".so":    true, // Linux共享对象
	".dll":   true, // Windows动态链接库
	".dylib": true, // macOS动态链接库

	// 静态库
	".a":   true, // Linux静态库
	".lib": true, // Windows静态库

	// 可执行文件
	".exe": true, // Windows可执行文件
	".com": true, // DOS命令文件
	".bin": true, // 二进制文件
	".elf": true, // Linux可执行文件
	".out": true, // 编译输出文件

	// 目标文件
	".o":   true, // 目标文件
	".obj": true, // Windows目标文件

	// 调试和符号文件
	".pdb": true, // Windows调试符号
	".exp": true, // Windows导出文件
	".ilk": true, // Windows增量链接
	".pch": true, // 预编译头文件
	".sbr": true, // 源代码浏览器
	".idb": true, // Visual Studio增量调试

	// 字节码文件
	".pyc":   true, // Python字节码
	".pyo":   true, // Python优化字节码
	".pyd":   true, // Python扩展模块
	".class": true, // Java字节码

	// 包文件
	".egg":  true, // Python包
	".jmod": true, // Java模块

	// 现代编译产物
	".wasm":        true, // WebAssembly
	".bc":          true, // LLVM位码
	".ll":          true, // LLVM中间表示
	".node":        true, // Node.js原生插件
	".rlib":        true, // Rust库
	".swiftmodule": true, // Swift模块

	// 映射和列表文件
	".map": true, // 映射文件
	".lst": true, // 列表文件
	".d":   true, // 依赖文件

	// 现代编译产物
	".dSYM":     true, // macOS调试符号
	".vsix":     true, // VS Code扩展
	".nupkg":    true, // NuGet包
	".gem":      true, // Ruby Gem
	".crate":    true, // Rust Crate
	".wheel":    true, // Python Wheel
	".snap":     true, // Snap包
	".flatpak":  true, // Flatpak包
	".appimage": true, // AppImage包
}
