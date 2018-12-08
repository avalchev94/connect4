package connect4

import "fmt"

type Game struct {
	field  Field
	player Player
	state  State
	four   []Cell
}

func New(cols, rows int, starting Player) *Game {
	return &Game{
		field:  NewField(cols, rows),
		player: starting,
		state:  Running,
		four:   nil,
	}
}

func (g *Game) Turn(col int) (Cell, error) {
	switch {
	case !g.Running():
		return Cell{}, fmt.Errorf("game is not running(over)")
	case !g.field.InRange(Cell{col, 0}):
		return Cell{}, fmt.Errorf("out of range: col %d", col)
	}

	cell, err := g.field.Update(col, g.player.Color())
	if err != nil {
		return Cell{}, err
	}

	four := g.findFour(cell)
	switch {
	case four != nil:
		switch g.player {
		case RedPlayer:
			g.state = RedWin
		case YellowPlayer:
			g.state = YellowWin
		}
		g.four = four
	case g.field.Full():
		g.state = EndDraw
	default:
		g.player = g.player.Next()
	}
	return cell, nil
}

func (g *Game) Running() bool {
	return g.state == Running
}

func (g *Game) State() State {
	return g.state
}

func (g *Game) Player() Player {
	return g.player
}

func (g *Game) Field() Field {
	return g.field
}

func (g *Game) Four() []Cell {
	return g.four
}
