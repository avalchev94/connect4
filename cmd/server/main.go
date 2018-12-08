package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/avalchev94/connect4/server"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()

	if err := http.ListenAndServe(*addr, server.New()); err != nil {
		log.Fatalln(err)
	}
}
