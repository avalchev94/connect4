package server

import (
	"net/http"
)

func New() *http.ServeMux {
	m := http.NewServeMux()

	//m.Handle("/rooms", nil)
	m.HandleFunc("/new", newRoom)
	m.HandleFunc("/join", joinRoom)

	return m
}
