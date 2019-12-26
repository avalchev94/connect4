package tarantula

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/avalchev94/tarantula/games"
	"github.com/pkg/errors"
	"nhooyr.io/websocket"
)

type Room struct {
	game    games.Game
	mutex   *sync.Mutex
	players Players
	join    chan playerTuple
	leave   chan string
}

func NewRoom(game games.Game) *Room {
	return &Room{
		game:    game,
		mutex:   &sync.Mutex{},
		players: Players{},
		join:    make(chan playerTuple),
		leave:   make(chan string),
	}
}

func (r *Room) findPlayer(playerID games.PlayerID) *Player {
	for _, player := range r.players {
		if player.id == playerID {
			return player
		}
	}
	return nil
}

func (r *Room) GameSettings() games.Settings {
	return r.game.Settings()
}

func (r *Room) PlayersCount() int {
	return len(r.players)
}

func (r *Room) PlayerExist(uuid string, playerID games.PlayerID) bool {
	player, ok := r.players[uuid]
	if !ok {
		return false
	}

	return player.id == playerID
}

func (r *Room) Join(uuid string) (games.PlayerID, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	playerID, err := r.game.AddPlayer()
	if err != nil {
		return -1, err
	}

	r.players[uuid] = NewPlayer(playerID)
	// todo: check if player actually joined??
	return playerID, nil
}

type playerTuple struct {
	uuid string
	conn *websocket.Conn
}

func (r *Room) Connect(uuid string, conn *websocket.Conn) {
	r.join <- playerTuple{uuid, conn}
}

func (r *Room) Run() {
	var (
		moveTimer     = time.NewTimer(time.Second)
		currentPlayer = NewPlayer(0)
	)
	moveTimer.Stop()

	for {
		select {
		case player := <-r.join:
			if err := r.handleJoin(player); err != nil {
				log.Println(err)
			}

		case player := <-r.leave:
			if err := r.handleLeave(player); err != nil {
				log.Println(err)
			}

		case state := <-r.game.StateUpdated():
			switch state {
			case games.Running:
				currentPlayer = r.findPlayer(r.game.CurrentPlayer())
				moveTimer.Reset(30 * time.Second)
			case games.EndDraw, games.EndWin:
				r.handleEnd()
			}

		case <-moveTimer.C:
			r.handleTimeExpired(currentPlayer)

		case msg := <-currentPlayer.read:
			r.handleMove(msg, currentPlayer)

			currentPlayer = r.findPlayer(r.game.CurrentPlayer())
			moveTimer.Reset(30 * time.Second)
		}
	}
}

func (r *Room) handleJoin(data playerTuple) error {
	if r.game.State() == games.Running {
		return errors.New("the game is running")
	}

	// update the socket of that player
	player := r.players[data.uuid]
	if player != nil && player.socket != nil {
		return errors.New("the connection of this uuid is active")
	}

	// update player socket and start processing messages
	player.SetConnection(data.conn)

	log.Printf("[Room %q] player %q joined!", "name", data.uuid)

	go func() {
		if err := player.ProcessMessages(context.TODO()); err != nil {
			log.Printf("[Room %q] player %q disconected: %v", "name", data.uuid, err)
		}
		r.leave <- data.uuid
	}()

	// check if there are players with nil sockets
	for _, p := range r.players {
		if p.Connection() == nil {
			return nil
		}
	}

	// try start the game
	if err := r.game.Start(); err != nil {
		if err == games.PlayersNotEnough {
			return nil
		}
		return err
	}

	// send game starting message
	msg := Message{
		Type: GameStarting,
		Payload: payloadGameStarting{
			Staring: r.game.CurrentPlayer(),
		},
	}
	r.players.SendAll(msg)
	return nil
}

func (r *Room) handleLeave(uuid string) error {
	player := r.players[uuid]
	player.SetConnection(nil)

	msg := Message{
		Type: PlayerLeft,
		Payload: payloadPlayer{
			Player: player.id,
		},
	}
	r.players.SendBut(msg, player)

	// pause the game
	return r.game.Pause()
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
