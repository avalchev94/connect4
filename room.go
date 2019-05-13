package tarantula

import (
	"sync"

	"github.com/avalchev94/tarantula/games"
	"github.com/gorilla/websocket"
)

type Message struct {
	Move   games.MoveData
	Player games.PlayerID
	State  games.GameState
}

type Room struct {
	Players Players
	Mutex   *sync.Mutex
	games.Game
}

func (r Room) Run() error {
	if err := r.Players.StartGame(r.CurrentPlayer()); err != nil {
		return err
	}

	for r.State() == games.Running {
		currentPlayer := r.CurrentPlayer()

		// get message from the current player
		var msg Message
		if err := r.Players[currentPlayer].ReadJSON(&msg); err != nil {
			return err
		}

		// update game logic with the message data
		if err := r.Move(msg.Player, msg.Move); err != nil {
			return err
		}

		// send them message to the rest of the players
		msg.State = games.Running
		if err := r.Players.Send(msg, currentPlayer); err != nil {
			return err
		}
	}

	// end game
	return r.Players.EndGame(r.State(), r.CurrentPlayer())
}

func (r Room) AddPlayer(conn *websocket.Conn) error {
	id, err := r.Game.AddPlayer()
	if err != nil {
		return err
	}

	r.Players[id] = conn

	return nil
}

type Rooms struct {
	rooms map[string]*Room
	mutex *sync.RWMutex
}

func NewRooms() *Rooms {
	return &Rooms{
		rooms: map[string]*Room{},
		mutex: &sync.RWMutex{},
	}
}

func (r *Rooms) Add(name string, room *Room) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if len(name) == 0 {
		return EmptyRoomName
	}

	if _, ok := r.rooms[name]; ok {
		return UsedRoomName
	}

	r.rooms[name] = room
	return nil
}

func (r *Rooms) Get(name string) (*Room, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if len(name) == 0 {
		return nil, EmptyRoomName
	}

	room, ok := r.rooms[name]
	if !ok {
		return nil, WrongRoomName
	}

	return room, nil
}

func (r *Rooms) Delete(name string) {
	r.mutex.Lock()
	defer r.mutex.RUnlock()

	delete(r.rooms, name)
}
