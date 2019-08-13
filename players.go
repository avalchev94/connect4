package tarantula

import (
	"fmt"

	"github.com/avalchev94/tarantula/games"
	"github.com/gorilla/websocket"
)

type Message struct {
	Move   games.MoveData
	Player games.PlayerID
	State  games.GameState
}

type Player struct {
	id     games.PlayerID
	socket *websocket.Conn
}

func (p *Player) Send(msg Message) error {
	if err := p.socket.WriteJSON(msg); err != nil {
		// Connection failed
		p.socket.Close()
		p.socket = nil

		return err
	}
	return nil
}

func (p *Player) Read() (Message, error) {
	var msg Message
	if err := p.socket.ReadJSON(&msg); err != nil {
		// Connection failed
		p.socket.Close()
		p.socket = nil
		return msg, err
	}
	return msg, nil
}

type Players map[string]*Player

func (p Players) StartGame() error {
	for _, player := range p {
		msg := Message{
			Player: player.id,
			State:  games.Starting,
		}
		if err := player.Send(msg); err != nil {
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
	return p.SendAll(msg)
}

func (p Players) Send(msg Message, sender *Player) error {
	for _, player := range p {
		if player.id == sender.id {
			continue
		}

		if err := player.Send(msg); err != nil {
			return err
		}
	}
	return nil
}

func (p Players) SendAll(msg Message) error {
	for _, player := range p {
		if err := player.Send(msg); err != nil {
			return err
		}
	}
	return nil
}
