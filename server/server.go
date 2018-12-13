package server

import (
	"net/http"
)

// New returns server mux with following endpoints:
//  "/rooms" - list available rooms;
//	"/new" - creates new room; Socket
//	"/join" - joins available room; Socket
func New() *http.ServeMux {
	m := http.NewServeMux()

	m.HandleFunc("/rooms", handleListRooms)
	m.HandleFunc("/new", handleNewRoom)
	m.HandleFunc("/join", handleJoinRoom)

	return m
}
