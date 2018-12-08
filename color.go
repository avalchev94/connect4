package connect4

// Color describes the state of a single cell
type Color int8

const (
	// NullColor - cell is empty
	NullColor Color = iota
	// RedColor - red player has the cell
	RedColor
	// YellowColor - yellow player has the cell
	YellowColor
)

func ParseColor(b []byte) Color {
	switch string(b) {
	case "R":
		return RedColor
	case "Y":
		return YellowColor
	default:
		return NullColor
	}
}

func (c Color) String() string {
	switch c {
	case RedColor:
		return "R"
	case YellowColor:
		return "Y"
	default:
		return "_"
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
