package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/avalchev94/tarantula/server"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()

	if err := http.ListenAndServe(*addr, server.NewServer()); err != nil {
		log.Fatalln(err)
	}
}
