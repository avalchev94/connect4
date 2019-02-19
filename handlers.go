package tarantula

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/avalchev94/tarantula/games/connect4"
	"github.com/gorilla/websocket"
)

var (
	rooms    = map[string]Room{}
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	}
}

func handleListRooms(w http.ResponseWriter, r *http.Request) {
	roomsData := []interface{}{}
	for _, r := range rooms {
		data := struct {
			Name    string
			Players int
			Game    string
		}{r.Name, len(r.Players), r.Game.Name()}

		roomsData = append(roomsData, data)
	}

	if err := json.NewEncoder(w).Encode(roomsData); err != nil {
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

	rooms[name] = Room{
		Name:    name,
		Game:    connect4.NewGame(7, 6),
		Players: Players{},
	}
	log.Printf("New room %s was created.", name)

	w.WriteHeader(http.StatusCreated)
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

	if len(room.Players) == 2 {
		go func() {
			if err := room.Run(); err != nil {
				log.Fatalf("Game Run failed: %s", err)
			}
		}()
	}
}
