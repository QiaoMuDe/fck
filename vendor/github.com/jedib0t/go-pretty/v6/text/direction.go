package text

// Direction defines the overall flow of text. Similar to bidi.Direction, but
// simplified and specific to this package.
type Direction int

// Available Directions.
const (
	Default Direction = iota
	LeftToRight
	RightToLeft
)

const (
	RuneL2R = '\u202a'
	RuneR2L = '\u202b'
)

// Modifier returns a character to force the given direction for the text that
// follows the modifier.
func (d Direction) Modifier() string {
	switch d {
	case LeftToRight:
		return string(RuneL2R)
	case RightToLeft:
		return string(RuneR2L)
	}
	return ""
}
