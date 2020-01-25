package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/avalchev94/tarantula"
	"github.com/avalchev94/tarantula/games"
	"github.com/avalchev94/tarantula/games/connect4"
	"github.com/gorilla/mux"
	"nhooyr.io/websocket"
)

var (
	rooms     = tarantula.NewRooms()
	wsOptions = websocket.AcceptOptions{
		InsecureSkipVerify: true,
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
		game = connect4.NewGame(6, 7)
	case "poker":
		//game = poker.NewGame("ip:port", body.Name)
	}

	room := tarantula.NewRoom(game)
	room.SetLogger(tarantula.NewLogger(body.Name))
	if err := rooms.Add(body.Name, room); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	go room.Run(context.TODO())

	w.WriteHeader(http.StatusCreated)
}

func handleJoinRoom(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	room, err := rooms.Get(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uuid, err := tarantula.NewUUID()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := room.Join(uuid); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	playerID, _ := room.PlayerExist(uuid)
	cookie := authCookie{uuid, playerID}
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

	if _, ok := room.PlayerExist(cookie.UUID); !ok {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	conn, err := websocket.Accept(w, r, &wsOptions)
	if err != nil {
		log.Printf("[Room %q] played %q failed to join: %v", name, cookie.UUID, err)
		return
	}

	if err := room.Connect(cookie.UUID, conn); err != nil {
		conn.Close(http.StatusBadRequest, err.Error())
	}
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

	if playerID, ok := room.PlayerExist(cookie.UUID); !ok || playerID != cookie.Player {
		http.Error(w, "player does not exist", http.StatusBadRequest)
		return
	}

	response := struct {
		Player   games.PlayerID `json:"player"`
		Settings games.Settings `json:"settings"`
	}{cookie.Player, room.GameSettings()}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
