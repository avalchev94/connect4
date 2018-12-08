package server

import (
	"fmt"

	"github.com/avalchev94/connect4"
)

type Room struct {
	Players players
	Game    *connect4.Game
}

func (r Room) Run() error {
	if err := r.Players.StartGame(r.Game.Player()); err != nil {
		return err
	}

	for r.Game.Running() {
		player := r.Players[r.Game.Player()]

		var msg Message
		if err := player.Read(&msg); err != nil {
			return err
		}

		cell, err := r.Game.Turn(msg.Cell.Col)
		if err != nil {
			err := player.WriteError(err)
			if err != nil {
				return err
			}
		}

		if cell != msg.Cell {
			return fmt.Errorf("client cell is different from the server cell")
		}

		switch r.Game.State() {
		case connect4.Running:
			if err := r.Players[r.Game.Player()].Write(Message{
				GameState: r.Game.State(),
				Cell:      msg.Cell,
			}); err != nil {
				return err
			}
		default:
			if err := r.Players.EndGame(r.Game.State()); err != nil {
				return err
			}
		}
	}

	return nil
}
