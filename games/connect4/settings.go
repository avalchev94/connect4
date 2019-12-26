package connect4

import (
	"github.com/avalchev94/tarantula/games"
)

type GameProgress struct {
	Player Color           `json:"player"`
	Field  [][]Color       `json:"field"`
	State  games.GameState `json:"state"`
}

type Settings struct {
	Cols         int          `json:"cols"`
	Rows         int          `json:"rows"`
	GameProgress GameProgress `json:"gameProgress"`
}

func (s Settings) Name() string {
	return "connect4"
}
