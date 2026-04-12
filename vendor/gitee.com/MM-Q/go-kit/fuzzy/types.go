/*
fuzzy 包定义了模糊匹配的类型。
*/
package fuzzy

// Match 表示一个匹配到的字符串。
type Match struct {
	Str            string // 匹配到的字符串
	Index          int    // 匹配字符串在提供的切片中的索引
	MatchedIndexes []int  // 匹配字符的索引。用于高亮显示匹配项
	Score          int    // 用于对匹配结果进行排序的分数
}

const (
	firstCharMatchBonus            = 10  // 第一个字符匹配奖励
	matchFollowingSeparatorBonus   = 20  // 匹配分隔符后的奖励
	camelCaseMatchBonus            = 20  // 驼峰命名匹配奖励
	adjacentMatchBonus             = 5   // 相邻匹配奖励
	unmatchedLeadingCharPenalty    = -5  // 未匹配的前导字符惩罚
	maxUnmatchedLeadingCharPenalty = -15 // 最大未匹配的前导字符惩罚
)

// 分隔符字符。这些字符在匹配时会被忽略，但会获得奖励。
// 用于将单词分隔开。
// 例如，"hello world" 中的 " " 就是一个分隔符。
var separators = []rune("/-_ .\\")

// Matches 是 Match 结构体的切片
type Matches []Match

// Len 返回匹配结果的长度。
//
// 返回值:
//   - int: 匹配结果的长度
func (a Matches) Len() int { return len(a) }

// Swap 交换两个匹配结果的顺序。
//
// 参数:
//   - i: 第一个匹配结果的索引
//   - j: 第二个匹配结果的索引
func (a Matches) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Less 返回第一个匹配结果是否应该排在第二个匹配结果之前。
// 按分数降序排列，分数相同时按原始索引升序排列（稳定排序）。
//
// 参数:
//   - i: 第一个匹配结果的索引
//   - j: 第二个匹配结果的索引
//
// 返回值:
//   - bool: 如果第一个匹配结果应该排在第二个之前，则返回 true
func (a Matches) Less(i, j int) bool {
	if a[i].Score == a[j].Score {
		return a[i].Index < a[j].Index
	}
	return a[i].Score > a[j].Score
}

// Source 表示字符串列表的抽象源。Source 必须是可迭代类型，如切片。
// 源将被迭代直到 Len()，对每个元素调用 String(i)，其中 i 是元素的索引。
// 你可以在 README 中找到工作示例。
type Source interface {
	String(i int) string // 位置 i 处要匹配的字符串
	Len() int            // 源的长度。通常是你想要匹配的事物切片的长度
}

/*
stringSource 是 Source 的一个实现，用于字符串切片。
*/
type stringSource []string

// String 返回字符串切片中位置 i 处的字符串。
//
// 参数:
//   - i: 要返回的字符串的索引
//
// 返回值:
//   - string: 位置 i 处的字符串
func (ss stringSource) String(i int) string {
	return ss[i]
}

// Len 返回字符串切片的长度。
//
// 返回值:
//   - int: 字符串切片的长度
func (ss stringSource) Len() int { return len(ss) }
