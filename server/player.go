package server

import (
	"log"
	"strconv"

	"github.com/avalchev94/connect4"
	"github.com/gorilla/websocket"
)

type Message struct {
	GameState connect4.State
	Cell      connect4.Cell
	Error     string
}

type player struct {
	Conn  *websocket.Conn
	Color connect4.Color
}

func (p player) StartGame(order int) error {
	data := []byte(strconv.FormatInt(int64(order), 10))
	return p.Conn.WriteMessage(websocket.TextMessage, []byte(data))
}

func (p player) Write(msg Message) error {
	err := p.Conn.WriteJSON(msg)
	log.Printf("wrote to %s player: %v", p.Color, msg)

	return err
}

func (p player) Read(msg *Message) error {
	err := p.Conn.ReadJSON(msg)
	log.Printf("%s player move: %v", p.Color, msg)

	return err
}

func (p player) WriteError(err error) error {
	return p.Write(Message{
		Error: err.Error(),
	})
}

func (p player) EndGame(state connect4.State) error {
	return p.Write(Message{
		GameState: state,
	})
}

type players map[connect4.Player]player

func (p players) StartGame(firstPlayer connect4.Player) error {
	if err := p[firstPlayer].StartGame(1); err != nil {
		return err
	}
	return p[firstPlayer.Next()].StartGame(2)
}

func (p players) EndGame(state connect4.State) error {
	for _, v := range p {
		if err := v.EndGame(state); err != nil {
			return err
		}
	}
	return nil
}
