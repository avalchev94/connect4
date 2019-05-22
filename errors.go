package tarantula

import "errors"

var (
	EmptyRoomName = errors.New("Empty room name")
	UsedRoomName  = errors.New("Used room name")
	WrongRoomName = errors.New("Room name does not exist")
	UUIDTaken     = errors.New("There is already player with that UUID")
)
