package tarantula

import (
	"log"
	"sync"

	"github.com/avalchev94/tarantula/games"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
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
	//moveTimer := time.NewTimer(time.Second)

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

		default:
			switch r.game.State() {
			case games.Running:
				if err := r.handleMove(); err != nil {
					log.Println(err)
				}
			case games.EndDraw, games.EndWin:
				if err := r.handleEnd(); err != nil {
					log.Println(err)
				}

				// temp -> no restart for now
				return
			}
		}
	}
}

func (r *Room) handleJoin(data playerTuple) error {
	if r.game.State() == games.Running {
		return errors.New("the game is running")
	}

	// update the socket of that player
	player := r.players[data.uuid]
	if player.socket != nil {
		return errors.New("the connection of this uuid is active")
	}

	// update player socket and start processing messages
	player.socket = data.conn

	go func() {
		if err := player.ProcessMessages(); err != nil {
			log.Println(err)
		}
		r.leave <- data.uuid
	}()

	// check if there are players with nil sockets
	for _, p := range r.players {
		if p.socket == nil {
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
	player.socket = nil

	msg := Message{
		Type: PlayerLeft,
		Payload: payloadPlayerLeft{
			Player: player.id,
		},
	}
	r.players.SendBut(msg, player)

	// pause the game
	return r.game.Pause()
}

func (r *Room) handleMove() error {
	currPlayer := r.findPlayer(r.game.CurrentPlayer())
	if currPlayer == nil {
		return errors.Errorf("failed to find current player: %v", r.game.CurrentPlayer())
	}

	msg := currPlayer.Read()
	movePayload, ok := msg.Payload.(payloadPlayerMove)
	if !ok {
		return errors.Errorf("failed to parse move payload: %v", msg)
	}

	// update game logic with the message data
	if err := r.game.Move(movePayload.Player, movePayload.Move); err != nil {
		return err
	}

	// send them message to the rest of the players
	r.players.SendBut(msg, currPlayer)

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
