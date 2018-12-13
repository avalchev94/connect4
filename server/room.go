package server

import (
	"fmt"

	"github.com/gorilla/websocket"

	"github.com/avalchev94/connect4"
)

type room struct {
	Name    string
	Players players
	Game    *connect4.Game
}

func newRoom(name string /*, gameOptions*/) room {
	return room{
		Name:    name,
		Players: players{},
		Game:    connect4.NewGame(7, 6, connect4.RedPlayer),
	}
}

func (r room) Run() error {
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

func (r room) AddPlayer(conn *websocket.Conn) {
	switch len(r.Players) {
	case 0:
		r.Players[connect4.RedPlayer] = player{conn, connect4.RedColor}
	case 1:
		r.Players[connect4.YellowPlayer] = player{conn, connect4.YellowColor}
	default:
		// todo: add watchers
	}
}
