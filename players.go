package tarantula

import (
	"github.com/avalchev94/tarantula/games"
	"github.com/gorilla/websocket"
)

type Player struct {
	id     games.PlayerID
	socket *websocket.Conn
	read   chan Message
	send   chan Message
}

func NewPlayer(id games.PlayerID) *Player {
	return &Player{
		id:     id,
		socket: nil,
		read:   make(chan Message, 100),
		send:   make(chan Message, 100),
	}
}

func (p *Player) Send(msg Message) {
	p.send <- msg
}

func (p *Player) Read() Message {
	return <-p.read
}

// func (p *Player) Ping() error {
// 	return p.socket.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second))
// }

func (p *Player) ProcessMessages() error {
	errorChan := make(chan error, 2)

	// read messages from the socket
	go func() {
		for {
			var msg Message
			if err := p.socket.ReadJSON(&msg); err != nil {
				errorChan <- err
				break
			}
			p.read <- msg
		}
	}()

	// send messages to the socket
	go func() {
		for msg := range p.send {
			if err := p.socket.WriteJSON(msg); err != nil {
				// failed to send the message, add it back
				p.send <- msg

				errorChan <- err
				break
			}
		}
	}()

	return <-errorChan
}

type Players map[string]*Player

func (p Players) SendBut(msg Message, sender *Player) {
	for _, player := range p {
		if player.id == sender.id {
			continue
		}

		player.Send(msg)
	}
}

func (p Players) SendAll(msg Message) {
	for _, player := range p {
		player.Send(msg)
	}
}
