package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	rooms    = map[string]room{}
	upgrader = websocket.Upgrader{}
)

func handleListRooms(w http.ResponseWriter, r *http.Request) {
	roomNames := []string{}
	for r := range rooms {
		roomNames = append(roomNames, r)
	}

	if err := json.NewEncoder(w).Encode(roomNames); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func handleNewRoom(w http.ResponseWriter, r *http.Request) {
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

	room := newRoom(name)
	room.AddPlayer(conn)
	rooms[name] = room

	log.Printf("New room %s was created.", name)
}

func handleJoinRoom(w http.ResponseWriter, r *http.Request) {
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
	room.AddPlayer(conn)
	log.Printf("Second player joined %s room.", name)
	go func() {
		if err := room.Run(); err != nil {
			log.Fatalf("Game Run failed: %s", err)
		}
	}()
}
