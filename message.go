package tarantula

import (
	"encoding/json"

	"github.com/avalchev94/tarantula/games"
)

type MessageType string

const (
	GameStarting       = MessageType("game_starting")
	GamePaused         = MessageType("game_paused")
	GameEnded          = MessageType("game_ended")
	PlayerMove         = MessageType("player_move")
	PlayerMoveExpired  = MessageType("player_move_expired")
	PlayerJoined       = MessageType("player_joined")
	PlayerConnected    = MessageType("player_connected")
	PlayerDisconnected = MessageType("player_disconnected")
	PlayerLeft         = MessageType("player_left")
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
	Move   moveData       `json:"move"`
}

type payloadPlayer struct {
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
	case PlayerMoveExpired, PlayerJoined, PlayerLeft:
		var payload payloadPlayer
		err = json.Unmarshal(msg["payload"], &payload)
		m.Payload = payload
	}
	return err
}
