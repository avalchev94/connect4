package tarantula

import "encoding/json"

type moveData struct {
	json.RawMessage
	Expired bool
}

func (m moveData) TimeExpired() bool {
	return m.Expired
}

func (m moveData) Decode(out interface{}) error {
	return json.Unmarshal(m.RawMessage, out)
}
