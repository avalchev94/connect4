package connect4

// Color describes the state of a single cell
type Color string

const (
	NullColor   Color = ""
	RedColor    Color = "red"
	YellowColor Color = "yellow"
)

func (c Color) Next() Color {
	switch c {
	case RedColor:
		return YellowColor
	default:
		return RedColor
	}
}
