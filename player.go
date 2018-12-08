package connect4

type Player int8

const (
	// RedPlayer ...
	RedPlayer Player = iota
	// YellowPlayer ...
	YellowPlayer
)

func (p Player) Next() Player {
	switch p {
	case RedPlayer:
		return YellowPlayer
	default:
		return RedPlayer
	}
}

func (p Player) Color() Color {
	switch p {
	case RedPlayer:
		return RedColor
	default:
		return YellowColor
	}
}

func (p Player) String() string {
	switch p {
	case RedPlayer:
		return "Red"
	default:
		return "Yellow"
	}
}
