package list // import "gitee.com/MM-Q/fck/commands/list"

Package list 实现了文件列表显示命令的主要逻辑。 该文件包含 list 子命令的入口函数，负责参数验证、路径处理、文件扫描和格式化输出。

Package list 实现了文件列表显示的颜色输出功能。 该文件提供了统一的跨平台颜色方案，根据文件类型和扩展名进行彩色显示。

Package list 扩展名映射表定义 该文件定义了各种文件类型的扩展名映射表，用于颜色分类

Package list 定义了 list 子命令的命令行标志和参数配置。 该文件包含所有 list 命令支持的选项，如排序方式、显示格式、过滤条件等。

Package list 实现了文件列表的格式化输出功能。 该文件提供了表格格式和网格格式两种显示方式，支持颜色输出、权限显示、文件大小格式化等功能。

Package list 定义了 list 子命令使用的数据模型和结构体。 该文件包含文件信息结构体、扫描选项、处理选项、格式化选项等核心数据类型定义。

Package list 实现了文件列表的数据处理功能。 该文件提供了文件列表的排序处理，支持按名称、时间、大小等多种方式排序。

Package list 实现了文件系统扫描功能。 该文件提供了文件和目录的扫描、过滤、类型识别等核心功能，支持递归扫描和多种文件类型过滤。

VARIABLES

var DefaultIcons = IconMap{
	ByType: map[EntryType]string{
		DirType:     "\uf07b",
		FileType:    "\uf016",
		SymlinkType: "\uf0c1",
		PipeType:    "\uf781",
		SocketType:  "\uf1e6",
	},
	ByExt: map[string]string{

		".zip": "\uf410", "zip": "\uf410",
		".tar": "\uf410", "tar": "\uf410",
		".gz": "\uf410", "gz": "\uf410",
		".rar": "\uf410", "rar": "\uf410",
		".7z": "\uf410", "7z": "\uf410",

		".txt": "\uf15c", "txt": "\uf15c",
		".md": "\uf48a", "md": "\uf48a",
		".rst": "\uf15c", "rst": "\uf15c",

		".go": "\ue627", "go": "\ue627",
		".py": "\ue73c", "py": "\ue73c",
		".js": "\ue74e", "js": "\ue74e",
		".ts": "\ue628", "ts": "\ue628",
		".tsx": "\ue7ba", "tsx": "\ue7ba",
		".jsx": "\ue7ba", "jsx": "\ue7ba",
		".json": "\ue60b", "json": "\ue60b",
		".yml": "\uf481", "yml": "\uf481",
		".yaml": "\uf481", "yaml": "\uf481",
		".sh": "\uf489", "sh": "\uf489",
		".bat": "\uf489", "bat": "\uf489",
		".ps1": "\uf489", "ps1": "\uf489",

		".pdf": "\uf1c1", "pdf": "\uf1c1",
		".doc": "\uf1c2", "doc": "\uf1c2",
		".docx": "\uf1c2", "docx": "\uf1c2",
		".xls": "\uf1c3", "xls": "\uf1c3",
		".xlsx": "\uf1c3", "xlsx": "\uf1c3",
		".ppt": "\uf1c4", "ppt": "\uf1c4",
		".pptx": "\uf1c4", "pptx": "\uf1c4",

		".png": "\uf1c5", "png": "\uf1c5",
		".jpg": "\uf1c5", "jpg": "\uf1c5",
		".jpeg": "\uf1c5", "jpeg": "\uf1c5",
		".gif": "\uf1c5", "gif": "\uf1c5",
		".svg": "\uf1c5", "svg": "\uf1c5",
		".mp3": "\uf1c7", "mp3": "\uf1c7",
		".wav": "\uf1c7", "wav": "\uf1c7",
		".mp4": "\uf1c8", "mp4": "\uf1c8",
		".mkv": "\uf1c8", "mkv": "\uf1c8",

		".exe": "\uf489", "exe": "\uf489",
		".bin": "\uf489", "bin": "\uf489",
	},

	Default: "\uf016",
}
    DefaultIcons 默认图标映射表 注：编码均为字符串，可直接拼接到名称前；是否加空格/着色由调用方决定。


FUNCTIONS

func GetColorString(info FileInfo, path string, cl *colorlib.ColorLib) string
    GetColorString 根据文件信息返回带有相应颜色的路径字符串

    参数:
      - info: 文件信息，包含文件类型和扩展名等信息
      - path: 要处理的路径字符串
      - cl: 用于彩色输出的colorlib.ColorLib实例

    返回:
      - string: 经过颜色处理后的路径字符串

    颜色方案:
      - 蓝色: 目录
      - 青色: 符号链接
      - 绿色: 可执行文件和代码文件
      - 黄色: 设备文件和配置文件
      - 红色: 数据文件和压缩包
      - 紫色: 库文件和编译产物
      - 灰色: 空文件
      - 白色: 其他文件

func GetIcon(info FileInfo) string
    GetIcon 根据文件信息返回图标编码。 规则：优先扩展名，后按类型，最后默认。

func InitListCmd() *cmd.Cmd
func ListCmdMain(cl *colorlib.ColorLib) error
    ListCmdMain list 命令主函数

    参数:
      - cl: 颜色库

    返回:
      - error: 错误信息


TYPES

type EntryType string
    EntryType 定义文件类型

const (
	DirType         EntryType = "d" // 目录类型
	SymlinkType     EntryType = "l" // 符号链接类型
	SocketType      EntryType = "s" // 套接字类型
	PipeType        EntryType = "p" // 管道类型
	BlockDeviceType EntryType = "b" // 块设备类型
	CharDeviceType  EntryType = "c" // 字符设备类型
	ExecutableType  EntryType = "x" // 可执行文件类型
	EmptyType       EntryType = "e" // 空文件类型
	FileType        EntryType = "f" // 普通文件类型
	UnknownType     EntryType = "?" // 未知类型
)
    定义文件类型标识符常量

func (e EntryType) String() string
    String 返回文件类型对应的字符串

type FileFormatter struct {
	// Has unexported fields.
}
    FileFormatter 文件格式化器

func NewFileFormatter(cl *colorlib.ColorLib) *FileFormatter
    NewFileFormatter 创建新的文件格式化器

func (f *FileFormatter) Render(files FileInfoList, opts FormatOptions) error
    Render 渲染文件列表

    参数:
      - files: 文件列表
      - opts: 格式选项

    返回:
      - error: 错误

type FileInfo struct {
	Name           string    // 文件名 - BaseName
	Path           string    // 文件路径 - 绝对路径
	OriginalPath   string    // 原始路径 - 用户指定的路径（用于分组显示）
	EntryType      EntryType // 类型 - 文件/目录/软链接
	Size           int64     // 大小 - 字节数
	ModTime        time.Time // 修改时间 - time.Time
	Perm           string    // 权限 - 类型-所有者-组-其他用户
	Owner          string    // 所属用户 - windows环境为?
	Group          string    // 所属组 - windows环境为?
	FileExt        string    // 扩展名 - 扩展名
	LinkTargetPath string    // 如果是软链接，则是指向的文件路径，否则为空字符串
}
    list子命令用于存储文件信息的结构体

type FileInfoList []FileInfo
    FileInfoList 文件信息列表类型

type FileProcessor struct{}
    FileProcessor 文件数据处理器

func NewFileProcessor() *FileProcessor
    NewFileProcessor 创建新的文件处理器

func (p *FileProcessor) Process(files FileInfoList, opts ProcessOptions) FileInfoList
    Process 处理文件列表: 过滤 -> 排序 -> 分组

    参数:
      - files: 文件列表
      - opts: 处理选项

    返回:
      - FileInfoList: 处理后的文件列表

type FileScanner struct {
	// Has unexported fields.
}
    FileScanner 文件扫描器

func NewFileScanner() *FileScanner
    NewFileScanner 创建新的文件扫描器

func (s *FileScanner) Scan(paths []string, opts ScanOptions) (FileInfoList, error)
    Scan 扫描指定路径的文件

    参数:
      - paths: 要扫描的路径列表
      - opts: 扫描选项

    返回:
      - FileInfoList: 扫描到的文件信息列表
      - error: 扫描过程中的错误

func (s *FileScanner) ScanWithOriginalPaths(originalPaths, expandedPaths []string, opts ScanOptions) (FileInfoList, error)
    ScanWithOriginalPaths 扫描指定路径的文件（保持原始路径和展开路径的对应关系）

    参数:
      - originalPaths: 用户输入的原始路径列表
      - expandedPaths: 展开后的路径列表
      - opts: 扫描选项

    返回:
      - FileInfoList: 扫描到的文件信息列表
      - error: 扫描过程中的错误

type FormatOptions struct {
	LongFormat    bool   // 是否长格式显示
	UseColor      bool   // 是否使用颜色
	TableStyle    string // 表格样式
	QuoteNames    bool   // 是否引用文件名
	ShowUserGroup bool   // 是否显示用户组
	ShouldGroup   bool   // 是否应该分组显示 (新增：避免重复判断)
}
    FormatOptions 格式化选项

type IconMap struct {
	// 按扩展名映射，键为小写扩展名（支持包含"."或不包含"."的两种）
	ByExt map[string]string
	// 按文件类型映射
	ByType map[EntryType]string
	// 默认图标
	Default string
}
    IconMap 定义图标映射集合

type ProcessOptions struct {
	SortBy      string // 排序方式: "name", "time", "size"
	Reverse     bool   // 是否反向排序
	GroupByDir  bool   // 是否按目录分组 (原有的递归分组)
	GroupByPath bool   // 是否按路径分组 (新增：用于多路径/通配符场景)
	IsMultiPath bool   // 是否为多路径场景 (新增：标识符)
}
    ProcessOptions 处理选项

type ScanOptions struct {
	Recursive  bool     // 是否递归扫描
	ShowHidden bool     // 是否显示隐藏文件
	FileTypes  []string // 文件类型过滤
	DirItself  bool     // 是否只显示目录本身
}
    ScanOptions 扫描选项

