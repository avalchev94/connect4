package tarantula

import (
	"fmt"

	"github.com/avalchev94/tarantula/games"
	"github.com/gorilla/websocket"
)

type Players map[games.PlayerUUID]*websocket.Conn

func (p Players) StartGame(firstPlayer games.PlayerUUID) error {
	for uuid, conn := range p {
		if err := conn.WriteJSON(Message{
			Player: uuid,
			State:  games.Starting,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (p Players) EndGame(state games.GameState, player games.PlayerUUID) error {
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

func (p Players) Send(msg Message, sender games.PlayerUUID) error {
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
