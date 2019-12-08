package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/avalchev94/tarantula"
	"github.com/avalchev94/tarantula/games"
	"github.com/avalchev94/tarantula/games/connect4"
	"github.com/avalchev94/tarantula/games/poker"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	rooms    = tarantula.NewRooms()
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func handleListRooms(w http.ResponseWriter, r *http.Request) {
	roomsData := []interface{}{}
	rooms.ForEach(func(name string, r *tarantula.Room) error {
		data := struct {
			Name    string `json:"name"`
			Players int    `json:"players"`
			Game    string `json:"game"`
		}{name, r.PlayersCount(), r.GameSettings().Name()}
		roomsData = append(roomsData, data)

		return nil
	})

	if err := json.NewEncoder(w).Encode(roomsData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func handleNewRoom(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Name string `json:"name"`
		Game string `json:"game"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var game games.Game
	switch body.Game {
	case "connect4":
		game = connect4.NewGame(7, 6)
	case "poker":
		game = poker.NewGame("ip:port", body.Name)
	}

	room := tarantula.NewRoom(game)
	if err := rooms.Add(body.Name, room); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Room '%s' was created.", body.Name)
	w.WriteHeader(http.StatusCreated)

	go room.Run()
}

func handleJoinRoom(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	room, err := rooms.Get(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uuid, err := uuid.NewUUID()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	playerID, err := room.Join(uuid.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cookie := authCookie{uuid.String(), playerID}
	if err := encodeCookie(w, name, cookie); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleConnectRoom(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	room, err := rooms.Get(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cookie, err := decodeCookie(r, name)
	if err != nil {
		http.Error(w, "auth cookie not found", http.StatusUnauthorized)
		return
	}

	if !room.PlayerExist(cookie.UUID, cookie.Player) {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	room.Connect(cookie.UUID, conn)
}

func handleGameSettings(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	room, err := rooms.Get(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cookie, err := decodeCookie(r, name)
	if err != nil {
		http.Error(w, "auth cookie not found", http.StatusUnauthorized)
		return
	}

	if !room.PlayerExist(cookie.UUID, cookie.Player) {
		http.Error(w, "player does not exist", http.StatusBadRequest)
		return
	}

	response := struct {
		Player   games.PlayerID `json:"player"`
		Settings interface{}    `json:"settings"`
	}{cookie.Player, room.GameSettings()}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
