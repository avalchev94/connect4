package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/avalchev94/connect4"

	"github.com/avalchev94/connect4/client"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "the address of the server")

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatal("no arguments")
	}

	switch flag.Arg(0) {
	case "list":
		listRooms()
	case "new":
		newRoom(flag.Args()[1:])
	case "join":
		joinRoom(flag.Args()[1:])
	}
}

func listRooms() {
	resp, err := http.Get(fmt.Sprintf("http://%s/rooms", *addr))
	if err != nil {
		log.Fatal(err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		rooms := []string{}
		if err := json.NewDecoder(resp.Body).Decode(&rooms); err != nil {
			log.Fatalf("Reading rooms failed: %s", err)
		}

		for i, r := range rooms {
			fmt.Printf("%d. %s\n", i+1, r)
		}
	default:
		log.Fatal(resp.Status)
	}
}

func newRoom(args []string) {
	fs := flag.NewFlagSet("new", flag.ContinueOnError)
	name := fs.String("name", "", "name of the room")
	fs.Parse(args)

	if len(*name) == 0 {
		log.Fatal("room name is empty")
	}

	url := fmt.Sprintf("ws://%s/new?name=%s", *addr, *name)
	conn, resp, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("Failed to dial server: %s", err)
	}

	switch resp.StatusCode {
	case http.StatusSwitchingProtocols:
		c := client.Client{
			Conn:  conn,
			Field: connect4.NewField(7, 6),
			Color: connect4.RedColor,
		}

		if err := c.Run(); err != nil {
			log.Fatalf("Failed to Run the game: %s", err)
		}
	default:
		data := []byte{}
		resp.Body.Read(data)
		log.Fatalf("Failed to create new room: %d %s", resp.StatusCode, string(data))
	}
}

func joinRoom(args []string) {
	fs := flag.NewFlagSet("join", flag.ContinueOnError) // todo?
	name := fs.String("name", "", "name of the room")
	fs.Parse(args)

	if len(*name) == 0 {
		log.Fatal("room name is empty")
	}

	url := fmt.Sprintf("ws://%s/join?name=%s", *addr, *name)
	conn, resp, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal(err)
	}

	switch resp.StatusCode {
	case http.StatusSwitchingProtocols:
		c := client.Client{
			Conn:  conn,
			Field: connect4.NewField(7, 6),
			Color: connect4.YellowColor,
		}

		if err := c.Run(); err != nil {
			log.Fatalf("Failed to Run the game: %s", err)
		}
	default:
		data := []byte{}
		resp.Body.Read(data)
		log.Fatalln(resp.StatusCode, string(data))
	}
}
