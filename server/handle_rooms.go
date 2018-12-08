package server

import (
	"log"
	"net/http"

	"github.com/avalchev94/connect4"

	"github.com/gorilla/websocket"
)

var (
	rooms    = map[string]Room{}
	upgrader = websocket.Upgrader{}
)

func newRoom(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if len(name) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Empty room name"))
		return
	}

	if _, ok := rooms[name]; ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Already used room name"))
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	rooms[name] = Room{
		Players: players{
			connect4.RedPlayer: player{
				Conn:  conn,
				Color: connect4.RedColor,
			},
		},
		Game: connect4.New(7, 6, connect4.RedPlayer),
	}
	log.Printf("New room %s was created.", name)
}

func joinRoom(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if len(name) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Empty room name"))
		return
	}

	if _, ok := rooms[name]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No such room name"))
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	room := rooms[name]
	room.Players[connect4.YellowPlayer] = player{
		Conn:  conn,
		Color: connect4.YellowColor,
	}

	log.Printf("Second player joined %s room.", name)
	go func() {
		if err := room.Run(); err != nil {
			log.Fatalf("Game Run failed: %s", err)
		}
	}()
}
