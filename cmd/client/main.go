package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":8081", "listen address")
	flag.Parse()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/rooms.html")
	})

	http.HandleFunc("/connect4/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/connect4.html")
	})

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalln(err)
	}
}
