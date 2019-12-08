package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewServer() http.Handler {
	m := mux.NewRouter()
	m.StrictSlash(true)
	m.Use(enableCORS)

	rooms := m.PathPrefix("/rooms").Subrouter()
	rooms.HandleFunc("/", handleListRooms).Methods(http.MethodGet)
	rooms.HandleFunc("/new", handleNewRoom).Methods(http.MethodPost)
	rooms.HandleFunc("/{name}/join", handleJoinRoom).Methods(http.MethodPost)
	rooms.HandleFunc("/{name}/connect", handleConnectRoom).Methods(http.MethodGet)
	rooms.HandleFunc("/{name}/settings", handleGameSettings).Methods(http.MethodGet)

	return m
}
