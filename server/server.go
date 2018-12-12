package server

import (
	"net/http"
)

func New() *http.ServeMux {
	m := http.NewServeMux()

	m.HandleFunc("/rooms", listRooms)
	m.HandleFunc("/new", newRoom)
	m.HandleFunc("/join", joinRoom)

	return m
}
