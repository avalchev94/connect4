package connect4

import (
	"fmt"

	"github.com/avalchev94/tarantula/games"
	"github.com/pkg/errors"
)

const (
	maxPlayers = 2
)

type Game struct {
	field       Field
	players     map[Color]bool
	currPlayer  Color
	state       games.GameState
	stateUpdate chan games.GameState
}

func NewGame(cols, rows int) *Game {
	return &Game{
		field:       NewField(cols, rows),
		players:     map[Color]bool{},
		currPlayer:  RedColor,
		state:       games.Starting,
		stateUpdate: make(chan games.GameState, 1),
	}
}

func (g *Game) setState(state games.GameState) {
	if g.state != state {
		g.state = state
		g.stateUpdate <- state
	}
}

func (g *Game) Start() error {
	switch g.state {
	case games.Running:
		return errors.New("game is already running")
	case games.Starting, games.Paused:
		if len(g.players) < maxPlayers {
			return errors.Errorf("Players are %d, but needed %d", len(g.players), maxPlayers)
		}

		for player, connected := range g.players {
			if !connected {
				return errors.Errorf("Player %q is not connected", player)
			}
		}

		g.setState(games.Running)
	default:
		return errors.New("not handeled for end game")
	}

	return nil
}

func (g *Game) Pause() error {
	switch g.state {
	case games.Paused:
		return errors.New("game is already paused")
	case games.Running:
		g.setState(games.Paused)
	default:
		return errors.New("not handled pause state")
	}

	return nil
}

func (g *Game) Move(player games.PlayerID, move games.MoveData) error {
	if move.TimeExpired() {
		g.currPlayer = g.currPlayer.Next()
		g.setState(games.EndWin)
		return nil
	}

	data := struct {
		Column int `json:"col"`
	}{}
	if err := move.Decode(&data); err != nil {
		return fmt.Errorf("failed to Decode MoveData: %v", err)
	}

	color := Color(player)
	if color != RedColor && color != YellowColor {
		return fmt.Errorf("player with id %v does not exist", player)
	}

	return g.moveInternal(color, data.Column)
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
		g.setState(games.EndWin)
	case g.field.Full():
		g.setState(games.EndDraw)
	default:
		g.currPlayer = player.Next()
	}

	return nil
}

func (g *Game) State() games.GameState {
	return g.state
}

func (g *Game) StateUpdated() <-chan games.GameState {
	return g.stateUpdate
}

func (g *Game) AddPlayer(connected bool) (games.PlayerID, error) {
	if len(g.players) == maxPlayers {
		return -1, fmt.Errorf("game has reached maximum players")
	}

	// try start the game after player is added
	defer g.Start()

	player := RedColor
	if _, ok := g.players[RedColor]; ok {
		player = YellowColor
	}

	g.players[player] = connected
	return games.PlayerID(player), nil
}

func (g *Game) DelPlayer(player games.PlayerID) error {
	color := Color(player)
	if _, ok := g.players[color]; !ok {
		return errors.Errorf("player %q doesn't exist", color)
	}

	delete(g.players, color)
	g.setState(games.Starting)

	return nil
}

func (g *Game) SetPlayerStatus(player games.PlayerID, connected bool) error {
	color := Color(player)
	if _, ok := g.players[color]; !ok {
		return errors.Errorf("player %q doesn't exist", color)
	}

	g.players[color] = connected
	if connected {
		g.Start()
	} else {
		g.Pause()
	}

	return nil
}

func (g *Game) CurrentPlayer() games.PlayerID {
	return games.PlayerID(g.currPlayer)
}

func (g *Game) Field() Field {
	return g.field
}

func (g *Game) Settings() games.Settings {
	return Settings{
		Rows: len(g.field),
		Cols: len(g.field[0]),
		GameProgress: GameProgress{
			Player: g.currPlayer,
			Field:  g.field,
			State:  g.state,
		},
	}
}
