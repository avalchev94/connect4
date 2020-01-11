package connect4

// Color describes the state of a single cell
type Color int

const (
	// NullColor - cell is empty
	NullColor Color = iota
	// RedColor - red player has the cell
	RedColor
	// YellowColor - yellow player has the cell
	YellowColor
)

func (c Color) String() string {
	switch c {
	case RedColor:
		return "red"
	case YellowColor:
		return "yellow"
	default:
		return "undefined"
	}
}

func (c Color) Next() Color {
	switch c {
	case RedColor:
		return YellowColor
	default:
		return RedColor
	}
}
