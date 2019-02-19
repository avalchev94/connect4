package tarantula

import (
	"fmt"
	"log"

	"github.com/avalchev94/tarantula/games"
	"github.com/gorilla/websocket"
)

type Players map[games.PlayerID]*websocket.Conn

func (p Players) StartGame(firstPlayer games.PlayerID) error {
	log.Println("starting game")
	for id, conn := range p {
		if err := conn.WriteJSON(Message{
			Player: id,
			State:  games.Starting,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (p Players) EndGame(state games.GameState, player games.PlayerID) error {
	if state == games.Starting || state == games.Running {
		return fmt.Errorf("game has not ended")
	}

	msg := Message{
		Player: player,
		State:  state,
	}

	for _, conn := range p {
		if err := conn.WriteJSON(msg); err != nil {
			return err
		}
	}
	return nil
}

func (p Players) Send(msg Message, sender games.PlayerID) error {
	for player, conn := range p {
		if player == sender {
			continue
		}

		if err := conn.WriteJSON(msg); err != nil {
			return err
		}
	}
	return nil
}
