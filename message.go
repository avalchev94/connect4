package tarantula

import (
	"encoding/json"

	"github.com/avalchev94/tarantula/games"
)

type MessageType int8

const (
	GameStarting MessageType = iota
	GameEnded
	PlayerMove
	PlayerJoined
	PlayerLeft
)

type Message struct {
	Type    MessageType `json:"type"`
	Payload interface{} `json:"payload"`
}

type payloadGameStarting struct {
	Staring games.PlayerID `json:"starting"`
}

type payloadGameEnded struct {
	State  games.GameState `json:"state"`
	Winner games.PlayerID  `json:"winner"`
}

type payloadPlayerMove struct {
	Player games.PlayerID `json:"player"`
	Move   games.MoveData `json:"move"`
}

type payloadPlayerJoined struct {
	Player games.PlayerID `json:"player"`
}

type payloadPlayerLeft struct {
	Player games.PlayerID `json:"player"`
}

func (m *Message) UnmarshalJSON(data []byte) error {
	msg := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}

	if err := json.Unmarshal(msg["type"], &m.Type); err != nil {
		return err
	}

	var err error
	switch m.Type {
	case GameStarting:
		var payload payloadGameStarting
		err = json.Unmarshal(msg["payload"], &payload)
		m.Payload = payload
	case GameEnded:
		var payload payloadGameEnded
		err = json.Unmarshal(msg["payload"], &payload)
		m.Payload = payload
	case PlayerMove:
		var payload payloadPlayerMove
		err = json.Unmarshal(msg["payload"], &payload)
		m.Payload = payload
	case PlayerJoined:
		var payload payloadPlayerJoined
		err = json.Unmarshal(msg["payload"], &payload)
		m.Payload = payload
	case PlayerLeft:
		var payload payloadPlayerLeft
		err = json.Unmarshal(msg["payload"], &payload)
		m.Payload = payload
	}
	return err
}
