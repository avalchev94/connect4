package games

import "errors"

// Game abstracts turn(move)-based game.
// Implement the following methods to use the game server.
// - Move(MoveData): notifies game engine about player move.
// - State() GameState: current state of the game.
// - CurrentPlayer() PlayerID: player that's on move.
// - AddPlayer() PlayerID: notifies game enging about new player, new id is returned.
type Game interface {
	Start() error
	Pause() error

	Move(PlayerID, MoveData) error
	State() GameState
	StateUpdated() <-chan GameState

	AddPlayer() (PlayerID, error)
	DeletePlayer(PlayerID) error
	CurrentPlayer() PlayerID

	Settings() Settings
}

// MoveData - data descriping a player move. All games need the "MoveData" in it's
// specific way. So, every game's frontend and backend should agree on a structure.
type MoveData interface {
	TimeExpired() bool
	Decode(interface{}) error
}

// Settings - basic game settings, that all games have. For more specific settings,
// object should be casted to the game's original settings type.
type Settings interface {
	Name() string
}

// PlayerID - the connection between game's logic and tarantula's logic.
type PlayerID int

// GameState - basic game state,
type GameState int8

// Self - descriptive
const (
	Starting GameState = iota
	Running
	Paused
	EndDraw
	EndWin
)

var (
	PlayersNotEnough = errors.New("players are not enough to start the game")
)
