package games

// Game abstracts turn(move)-based game.
// Implement the following methods to use the game server.
// - Move(MoveData): notifies game engine about player move.
// - State() GameState: current state of the game.
// - CurrentPlayer() PlayerID: player that's on move.
// - AddPlayer() PlayerID: notifies game enging about new player, new id is returned.
// - Name(): name of the game.
type Game interface {
	Move(PlayerID, MoveData) error
	State() GameState

	CurrentPlayer() PlayerID
	AddPlayer() (PlayerID, error)

	Name() string
}

// MoveData is an interface keeping data for player's move.
type MoveData interface {
}

// PlayerID - need to connect the players in the game enging and game server.
type PlayerID int

// GameState is an enum for the basic game states.
type GameState int8

const (
	// Starting - first state sent just before starting
	Starting GameState = iota
	// Running - game is in progress
	Running
	// EndDraw - game ended as draw.
	EndDraw
	// EndWin - game ended as win. Use Game.CurrentPlayer to see who.
	EndWin
)
