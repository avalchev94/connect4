package tarantula

import (
	"context"
	"sync"
	"time"

	"github.com/avalchev94/tarantula/games"
	"github.com/pkg/errors"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type Player struct {
	id     games.PlayerID
	socket *websocket.Conn
	mutex  *sync.Mutex
	read   chan Message
	send   chan Message
}

func NewPlayer(id games.PlayerID) *Player {
	return &Player{
		id:     id,
		socket: nil,
		mutex:  &sync.Mutex{},
		read:   make(chan Message, 100),
		send:   make(chan Message, 100),
	}
}

func (p *Player) Send(msg Message) {
	p.send <- msg
}

func (p *Player) Read() <-chan Message {
	return p.read
}

func (p *Player) Disconnected(timeout time.Duration) error {
	p.setConnection(nil)

	if timeout <= 0 {
		return nil
	}

	timer := time.NewTimer(timeout)
	for {
		select {
		case <-timer.C:
			if p.socket == nil {
				return errors.New("wait for player connection timeouted")
			}
			return nil
		default:
			if p.socket != nil {
				return nil
			}
		}
	}
}

func (p *Player) Connected(conn *websocket.Conn) {
	p.setConnection(conn)
}

func (p *Player) Connection() *websocket.Conn {
	return p.socket
}

func (p *Player) setConnection(conn *websocket.Conn) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.socket = conn
}

func (p *Player) ProcessMessages(parentCtx context.Context) error {
	errorChan := make(chan error, 2)
	ctx, cancelFunc := context.WithCancel(parentCtx)

	wg := sync.WaitGroup{}
	wg.Add(2)

	// read messages from the socket
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				var msg Message
				if err := wsjson.Read(ctx, p.socket, &msg); err != nil {
					errorChan <- errors.WithMessage(err, "Failed to read message")
					return
				}
				p.read <- msg
			}
		}
	}()

	// send messages to the socket
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-p.send:
				if err := wsjson.Write(ctx, p.socket, msg); err != nil {
					// failed to send the message, add it back
					p.send <- msg

					errorChan <- errors.WithMessage(err, "Failed to send message")
					return
				}
			}
		}
	}()

	socketErr := <-errorChan
	cancelFunc()
	wg.Wait()

	return socketErr
}

type Players map[UUID]*Player

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
		if player.Connection() != nil {
			player.Send(msg)
		}
	}
}
