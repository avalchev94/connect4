package connect4

type State int8

const (
	Running State = iota
	RedWin
	YellowWin
	EndDraw
)

func (s State) String() string {
	switch s {
	case EndDraw:
		return "Draw"
	case RedWin:
		return "Red is Winner"
	case YellowWin:
		return "Yellow is Winner"
	default:
		return "Running"
	}
}
