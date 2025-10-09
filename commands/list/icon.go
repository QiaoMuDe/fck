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
		SymlinkType:     "",       // 软链接
		SocketType:      "",       // 套接字
		PipeType:        "",       // 管道
		BlockDeviceType: "",       // 块设备
		CharDeviceType:  "",       // 字符设备
		ExecutableType:  "",       // 可执行文件
		EmptyType:       "",       // 空文件
		FileType:        "\uF15B", // 普通文件
		UnknownType:     "\uEBC3", // 未知类型
	},

	ByExt: map[string]string{
		// 压缩/归档
		".zip": "\uf410", "zip": "\uf410", // 󰐐 nf-oct-file_zip
	},

	Default: "\uF413", // 默认图标
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
	if info.EntryType == FileType {
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
