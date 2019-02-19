package connect4

import (
	"github.com/avalchev94/tarantula/games"
)

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

func (c Color) Next() Color {
	switch c {
	case RedColor:
		return YellowColor
	default:
		return RedColor
	}
}

func (c Color) PlayerID() games.PlayerID {
	return games.PlayerID(c)
}
