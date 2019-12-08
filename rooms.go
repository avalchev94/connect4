package tarantula

import "sync"

type Rooms struct {
	rooms map[string]*Room
	mutex *sync.RWMutex
}

func NewRooms() *Rooms {
	return &Rooms{
		rooms: map[string]*Room{},
		mutex: &sync.RWMutex{},
	}
}

func (r *Rooms) Add(name string, room *Room) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if len(name) == 0 {
		return EmptyRoomName
	}

	if _, ok := r.rooms[name]; ok {
		return UsedRoomName
	}

	r.rooms[name] = room
	return nil
}

func (r *Rooms) Get(name string) (*Room, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if len(name) == 0 {
		return nil, EmptyRoomName
	}

	room, ok := r.rooms[name]
	if !ok {
		return nil, WrongRoomName
	}

	return room, nil
}

func (r *Rooms) ForEach(f func(name string, r *Room) error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for name, room := range r.rooms {
		if err := f(name, room); err != nil {
			return
		}
	}
}

func (r *Rooms) Delete(name string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.rooms, name)
}
