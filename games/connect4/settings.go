package connect4

type Settings struct {
	Cols int `json:"cols"`
	Rows int `json:"rows"`
}

func (s Settings) Name() string {
	return "connect4"
}
