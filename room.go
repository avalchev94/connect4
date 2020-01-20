package tarantula

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/avalchev94/tarantula/games"
	"github.com/avalchev94/tarantula/pkg/timerx"
	"github.com/pkg/errors"
	"nhooyr.io/websocket"
)

type Room struct {
	game       games.Game
	mutex      *sync.Mutex
	players    Players
	join       chan chanMessage
	connect    chan chanMessage
	disconnect chan chanMessage
	leave      chan chanMessage
}

type chanMessage struct {
	uuid    UUID
	conn    *websocket.Conn
	errChan chan<- error
}

func NewRoom(game games.Game) *Room {
	return &Room{
		game:       game,
		mutex:      &sync.Mutex{},
		players:    Players{},
		join:       make(chan chanMessage),
		connect:    make(chan chanMessage),
		disconnect: make(chan chanMessage),
		leave:      make(chan chanMessage),
	}
}

func (r *Room) GameSettings() games.Settings {
	return r.game.Settings()
}

func (r *Room) PlayersCount() int {
	return len(r.players)
}

func (r *Room) PlayerExist(uuid UUID) (games.PlayerID, bool) {
	if player, ok := r.players[uuid]; ok {
		return player.id, true
	}

	return "", false
}

func (r *Room) Join(uuid UUID) error {
	errChan := make(chan error)
	r.join <- chanMessage{uuid, nil, errChan}

	return <-errChan
}

func (r *Room) Leave(uuid UUID) error {
	errChan := make(chan error)
	r.leave <- chanMessage{uuid, nil, errChan}

	return <-errChan
}

func (r *Room) Connect(uuid UUID, conn *websocket.Conn) error {
	errChan := make(chan error)
	r.connect <- chanMessage{uuid, conn, errChan}

	return <-errChan
}

func (r *Room) Run() {
	var (
		moveTimer     = timerx.NewTimer(30 * time.Second)
		currentPlayer = NewPlayer("")
	)
	moveTimer.Pause()

	for {
		select {
		case msg := <-r.join:
			if err := r.handleJoin(msg.uuid); msg.errChan != nil {
				msg.errChan <- err
			}
		case msg := <-r.connect:
			if err := r.handleConnect(msg.uuid, msg.conn); msg.errChan != nil {
				msg.errChan <- err
			}
		case msg := <-r.disconnect:
			if err := r.handleDisconnect(msg.uuid); msg.errChan != nil {
				msg.errChan <- err
			}
		case msg := <-r.leave:
			if err := r.handleLeave(msg.uuid); msg.errChan != nil {
				msg.errChan <- err
			}
		case state := <-r.game.StateUpdated():
			switch state {
			case games.Running:
				r.players.SendAll(Message{
					Type: GameStarting,
					Payload: payloadGameStarting{
						Staring:       r.game.CurrentPlayer(),
						MoveRemaining: int(moveTimer.Remaining().Seconds()),
					},
				})

				_, currentPlayer = r.findPlayer(r.game.CurrentPlayer())
				moveTimer.Start()
			case games.Paused:
				r.players.SendAll(Message{
					Type: GamePaused,
				})
				moveTimer.Pause()
			case games.EndDraw, games.EndWin:
				r.handleEnd()
				moveTimer.Stop()
			}

		case <-moveTimer.C:
			r.handleTimeExpired(currentPlayer)

		case msg := <-currentPlayer.Read():
			r.handleMove(msg, currentPlayer)

			_, currentPlayer = r.findPlayer(r.game.CurrentPlayer())
			moveTimer.Reset(30 * time.Second)
		}
	}
}

func (r *Room) handleJoin(uuid UUID) error {
	if _, ok := r.players[uuid]; ok {
		return errors.Errorf("player with uuid %q already exist", uuid)
	}

	playerID, err := r.game.AddPlayer(false)
	if err != nil {
		return err
	}

	player := NewPlayer(playerID)
	r.players[uuid] = player

	msg := Message{
		Type: PlayerJoined,
		Payload: payloadPlayer{
			Player: player.id,
		},
	}
	r.players.SendBut(msg, player)

	go func() {
		if err := player.Disconnected(30 * time.Second); err != nil {
			// log
			r.leave <- chanMessage{uuid: uuid}
		}
	}()

	return nil
}

func (r *Room) handleConnect(uuid UUID, conn *websocket.Conn) error {
	player, ok := r.players[uuid]
	if !ok {
		return errors.Errorf("player with uuid %q does not exist", uuid)
	} else if player.Connection() != nil {
		return errors.Errorf("player with uuid %q already has active connection", uuid)
	}

	if err := r.game.SetPlayerStatus(player.id, true); err != nil {
		return err
	}
	player.Connected(conn)

	msg := Message{
		Type: PlayerConnected,
		Payload: payloadPlayer{
			Player: player.id,
		},
	}
	r.players.SendBut(msg, player)

	go func() {
		if err := player.ProcessMessages(context.TODO()); err != nil {
			log.Printf("[Room %q] player %q disconected: %v", "name", uuid, err)
		}
		r.disconnect <- chanMessage{uuid: uuid}
	}()

	return nil
}

func (r *Room) handleDisconnect(uuid UUID) error {
	player, ok := r.players[uuid]
	if !ok {
		return errors.Errorf("player with uuid %q does not exist", uuid)
	} else if player.Connection() == nil {
		return errors.Errorf("player with uuid %q is already disconnected", uuid)
	}

	if err := r.game.SetPlayerStatus(player.id, false); err != nil {
		return err
	}

	msg := Message{
		Type: PlayerDisconnected,
		Payload: payloadPlayer{
			Player: player.id,
		},
	}
	r.players.SendBut(msg, player)

	go func() {
		if err := player.Disconnected(30 * time.Second); err != nil {
			// log
			r.leave <- chanMessage{uuid: uuid}
		}
	}()

	return nil
}

func (r *Room) handleLeave(uuid UUID) error {
	player, ok := r.players[uuid]
	if !ok {
		return errors.Errorf("player with uuid %q does not exist", uuid)
	}

	if err := r.game.DelPlayer(player.id); err != nil {
		return nil
	}
	delete(r.players, uuid)

	msg := Message{
		Type: PlayerLeft,
		Payload: payloadPlayer{
			Player: player.id,
		},
	}
	r.players.SendAll(msg)

	return nil
}

func (r *Room) handleMove(msg Message, player *Player) error {
	movePayload, ok := msg.Payload.(payloadPlayerMove)
	if !ok {
		return errors.Errorf("failed to parse move payload: %v", msg)
	}

	// update game logic with the message data
	if err := r.game.Move(movePayload.Player, movePayload.Move); err != nil {
		return err
	}

	// send them message to the rest of the players
	r.players.SendBut(msg, player)

	return nil
}

func (r *Room) handleEnd() error {
	msg := Message{
		Type: GameEnded,
		Payload: payloadGameEnded{
			State:  r.game.State(),
			Winner: r.game.CurrentPlayer(),
		},
	}

	r.players.SendAll(msg)
	return nil
}

func (r *Room) handleTimeExpired(player *Player) error {
	msg := Message{
		Type: PlayerMoveExpired,
		Payload: payloadPlayer{
			Player: player.id,
		},
	}
	r.players.SendBut(msg, player)

	return r.game.Move(player.id, moveData{Expired: true})
}

func (r *Room) findPlayer(playerID games.PlayerID) (UUID, *Player) {
	for uuid, player := range r.players {
		if player.id == playerID {
			return uuid, player
		}
	}
	return UUID(""), nil
}
