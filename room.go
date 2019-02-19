package tarantula

import (
	"github.com/avalchev94/tarantula/games"
	"github.com/gorilla/websocket"
)

type Message struct {
	Move   games.MoveData
	Player games.PlayerID
	State  games.GameState
}

type Room struct {
	Name    string
	Players Players
	games.Game
}

func (r Room) Run() error {
	if err := r.Players.StartGame(r.CurrentPlayer()); err != nil {
		return err
	}

	for r.State() == games.Running {
		currentPlayer := r.CurrentPlayer()

		// get message from the current player
		var msg Message
		if err := r.Players[currentPlayer].ReadJSON(&msg); err != nil {
			return err
		}

		// update game logic with the message data
		if err := r.Move(msg.Player, msg.Move); err != nil {
			return err
		}

		// send them message to the rest of the players
		msg.State = games.Running
		if err := r.Players.Send(msg, currentPlayer); err != nil {
			return err
		}
	}

	// end game
	return r.Players.EndGame(r.State(), r.CurrentPlayer())
}

// todo: check if room is full?
func (r Room) AddPlayer(conn *websocket.Conn) error {
	id, err := r.Game.AddPlayer()
	if err != nil {
		return err
	}

	r.Players[id] = conn

	return nil
}
