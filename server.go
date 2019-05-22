package tarantula

import (
	"net/http"
)

func NewServer() *http.ServeMux {
	m := http.NewServeMux()

	m.HandleFunc("/rooms", enableCORS(handleListRooms))
	m.HandleFunc("/new", enableCORS(handleNewRoom))
	m.HandleFunc("/join", enableCORS(handleJoinRoom))

	return m
}
