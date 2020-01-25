package tarantula

import "encoding/json"

type rawMessage struct {
	json.RawMessage
}

func (m rawMessage) String() string {
	return string(m.RawMessage)
}

type moveData struct {
	rawMessage
	Expired bool
}

func (m moveData) TimeExpired() bool {
	return m.Expired
}

func (m moveData) Decode(out interface{}) error {
	return json.Unmarshal(m.RawMessage, out)
}
