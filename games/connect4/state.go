package connect4

type State int8

const (
	Starting State = iota
	Running
	EndWin
	EndDraw
)

func (s State) Starting() bool {
	return s == Starting
}

func (s State) Running() bool {
	return s == Running
}

func (s State) EndWin() bool {
	return s == EndWin
}

func (s State) EndDraw() bool {
	return s == EndDraw
}
