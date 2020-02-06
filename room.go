package tarantula

import (
	"context"
	"time"

	"github.com/avalchev94/tarantula/games"
	"github.com/avalchev94/tarantula/pkg/timerx"
	"github.com/pkg/errors"
	"nhooyr.io/websocket"
)

type Room struct {
	game       games.Game
	players    Players
	logger     Logger
	join       chan chanMessage
	leave      chan chanMessage
	connect    chan chanMessage
	disconnect chan chanMessage
}

type chanMessage struct {
	uuid    UUID
	conn    *websocket.Conn
	errChan chan<- error
}

func NewRoom(game games.Game) *Room {
	return &Room{
		game:       game,
		players:    Players{},
		logger:     &dummyLogger{},
		join:       make(chan chanMessage),
		leave:      make(chan chanMessage),
		connect:    make(chan chanMessage),
		disconnect: make(chan chanMessage),
	}
}

func (r *Room) SetLogger(logger Logger) {
	if logger != nil {
		r.logger = logger
	} else {
		r.logger = dummyLogger{}
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

func (r *Room) Run(ctx context.Context) {
	var (
		moveTimer     = timerx.NewTimer(30 * time.Second)
		currentPlayer = NewPlayer("")
	)
	moveTimer.Pause()

	r.logger.Info("event loop running...")

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-r.join:
			if err := r.handleJoin(ctx, msg.uuid); msg.errChan != nil {
				msg.errChan <- err
			}
		case msg := <-r.connect:
			if err := r.handleConnect(ctx, msg.uuid, msg.conn); msg.errChan != nil {
				msg.errChan <- err
			}
		case msg := <-r.disconnect:
			if err := r.handleDisconnect(ctx, msg.uuid); msg.errChan != nil {
				msg.errChan <- err
			}
		case msg := <-r.leave:
			if err := r.handleLeave(msg.uuid); msg.errChan != nil {
				msg.errChan <- err
			}
		case state := <-r.game.StateUpdated():
			switch state {
			case games.Running:
				_, currentPlayer = r.findPlayer(r.game.CurrentPlayer())
				moveTimer.Start()

				r.players.SendAll(Message{
					Type: GameStarting,
					Payload: payloadGameStarting{
						Starting:      currentPlayer.id,
						MoveRemaining: int(moveTimer.Remaining().Seconds()),
					},
				})
				r.logger.Debugf("game running, on move: %q", currentPlayer.id)
			case games.Paused:
				moveTimer.Pause()

				r.players.SendAll(Message{
					Type: GamePaused,
				})
				r.logger.Debugf("game paused")
			case games.EndDraw, games.EndWin:
				moveTimer.Stop()

				r.players.SendAll(Message{
					Type: GameEnded,
					Payload: payloadGameEnded{
						State:  r.game.State(),
						Winner: r.game.CurrentPlayer(),
					},
				})
				r.logger.Debugf("game ended, state: %q, player: %q", r.game.State(), r.game.CurrentPlayer())
			}

		case <-moveTimer.C:
			if err := r.handleTimeExpired(currentPlayer); err != nil {
				r.players.SendAll(Message{
					Type:    GameError,
					Payload: payloadGameError{err},
				})
			}

		case msg := <-currentPlayer.Read():
			if err := r.handleMove(msg, currentPlayer); err != nil {
				currentPlayer.Send(Message{
					Type:    GameError,
					Payload: payloadGameError{err},
				})
			} else {
				_, currentPlayer = r.findPlayer(r.game.CurrentPlayer())
				moveTimer.Reset(30 * time.Second)
			}
		}
	}
}

func (r *Room) handleJoin(ctx context.Context, uuid UUID) error {
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
		if err := player.Disconnected(ctx, 30*time.Second); err != nil {
			r.leave <- chanMessage{uuid: uuid}
		}
	}()
	r.logger.Debugf("player %q joined", player.id)

	return nil
}

func (r *Room) handleConnect(ctx context.Context, uuid UUID, conn *websocket.Conn) error {
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
		if err := player.ProcessMessages(ctx); err != nil {
			r.disconnect <- chanMessage{uuid: uuid}
		}
	}()
	r.logger.Debugf("player %q connected", player.id)

	return nil
}

func (r *Room) handleDisconnect(ctx context.Context, uuid UUID) error {
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
		if err := player.Disconnected(ctx, 30*time.Second); err != nil {
			r.leave <- chanMessage{uuid: uuid}
		}
	}()
	r.logger.Debugf("player %q disconnected", player.id)

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

	r.logger.Debugf("player %q left", player.id)

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

	r.logger.Debugf("player %q made move: %+v", player.id, msg.Payload)

	// send them message to the rest of the players
	r.players.SendBut(msg, player)
	return nil
}

func (r *Room) handleTimeExpired(player *Player) error {
	if err := r.game.Move(player.id, moveData{Expired: true}); err != nil {
		return err
	}

	msg := Message{
		Type: PlayerMoveExpired,
		Payload: payloadPlayer{
			Player: player.id,
		},
	}
	r.players.SendBut(msg, player)

	return nil
}

func (r *Room) findPlayer(playerID games.PlayerID) (UUID, *Player) {
	for uuid, player := range r.players {
		if player.id == playerID {
			return uuid, player
		}
	}
	return UUID(""), nil
}
