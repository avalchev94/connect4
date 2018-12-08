package client

import (
	"fmt"

	"github.com/avalchev94/connect4"
	"github.com/avalchev94/connect4/server"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn  *websocket.Conn
	Field connect4.Field
	Color connect4.Color
}

func (c *Client) Run() error {
	_, b, err := c.Conn.ReadMessage()
	if err != nil {
		return err
	}

	switch string(b) {
	case "1":
		fmt.Println("Start first...")
		render(c.Field)
		if err := c.ProcessTurn(); err != nil {
			return err
		}
	case "2":
		fmt.Println("Start second...")
	}

	for {
		render(c.Field)
		msg, err := c.Read()
		if err != nil {
			return err
		}

		switch msg.GameState {
		case connect4.Running:
			cell, err := c.Field.Update(msg.Cell.Col, c.Color.Next())
			if err != nil {
				return err
			}

			if cell != msg.Cell {
				return fmt.Errorf("client cell is not the same as server message")
			}

			render(c.Field)
			if err := c.ProcessTurn(); err != nil {
				return err
			}
		default:
			fmt.Printf("Game ended: %s", msg.GameState)
			return nil
		}
	}

}

func (c *Client) Read() (server.Message, error) {
	var msg server.Message
	return msg, c.Conn.ReadJSON(&msg)
}

func (c *Client) WaitColumn() int {
	fmt.Printf("Your move: ")

	var col int
	fmt.Scan(&col)
	return col
}

func (c *Client) ProcessTurn() error {
	cell, err := c.Field.Update(c.WaitColumn(), c.Color)
	if err != nil {
		return err
	}

	return c.Conn.WriteJSON(server.Message{
		Cell: cell,
	})
}
