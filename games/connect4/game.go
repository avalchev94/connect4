package connect4

import (
	"fmt"

	"github.com/avalchev94/tarantula/games"
)

const (
	maxPlayers = 2
)

type Game struct {
	field      Field
	players    []Color
	currPlayer Color
	state      games.GameState
}

func NewGame(cols, rows int) *Game {
	return &Game{
		field:      NewField(cols, rows),
		players:    []Color{},
		currPlayer: RedColor,
		state:      games.Running,
	}
}

func (g *Game) Move(player games.PlayerID, move games.MoveData) error {
	data, ok := move.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unsuccessful cast MoveData to map[string]interface{}")
	}

	color := Color(player)
	if color != RedColor && color != YellowColor {
		return fmt.Errorf("player with id %v does not exist", player)
	}

	column := int(data["col"].(float64))
	return g.moveInternal(color, column)
}

func (g *Game) moveInternal(player Color, column int) error {
	switch {
	case g.state != games.Running:
		return fmt.Errorf("game is not running")
	case g.currPlayer != player:
		return fmt.Errorf("player is not consistent?")
	case !g.field.InRange(Cell{column, 0}):
		return fmt.Errorf("out of range column %d", column)
	}

	cell, err := g.field.Update(column, player)
	if err != nil {
		return err
	}

	four := g.field.FindFour(cell)
	switch {
	case four != nil:
		g.state = games.EndWin
	case g.field.Full():
		g.state = games.EndDraw
	default:
		g.currPlayer = player.Next()
	}

	return nil
}

func (g *Game) State() games.GameState {
	return g.state
}

func (g *Game) AddPlayer() (games.PlayerID, error) {
	if len(g.players) == maxPlayers {
		return -1, fmt.Errorf("game has reached maximum players")
	}

	if len(g.players) == 0 {
		g.players = append(g.players, RedColor)
		return games.PlayerID(RedColor), nil
	}

	g.players = append(g.players, YellowColor)
	return games.PlayerID(YellowColor), nil
}

func (g *Game) CurrentPlayer() games.PlayerID {
	return games.PlayerID(g.currPlayer)
}

func (g *Game) Name() string {
	return "Connect4"
}

func (g *Game) Field() Field {
	return g.field
}
