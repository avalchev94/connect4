package poker

import (
	"net/http"
	"time"

	"github.com/avalchev94/tarantula/games"
)

type Game struct {
	host   string
	room   string
	client http.Client
}

func NewGame(host, room string) *Game {
	return &Game{
		host: host,
		room: room,
		client: http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (g *Game) Start() error {
	return nil
}

func (g *Game) Pause() error {
	return nil
}

func (g *Game) Move(player games.PlayerID, move games.MoveData) error {
	return nil
}

func (g *Game) State() games.GameState {
	return 0
}

func (g *Game) AddPlayer() (games.PlayerID, error) {
	return 0, nil
}

func (g *Game) CurrentPlayer() games.PlayerID {
	return 0
}

func (g *Game) Settings() games.Settings {
	return Settings{}
}
