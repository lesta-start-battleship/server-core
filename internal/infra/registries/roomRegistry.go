package registries

import "lesta-battleship/server-core/internal/multiplayer/actors/rooms"

type RoomRegistry struct {
	rooms map[string]*rooms.Room
}

func NewRoomRegistry() *RoomRegistry {
	return &RoomRegistry{rooms: make(map[string]*rooms.Room)}
}

func (r *RoomRegistry) Track(id string, room *rooms.Room) {
	r.rooms[id] = room
}

func (r *RoomRegistry) Find(id string) *rooms.Room {
	return r.rooms[id]
}

func (r *RoomRegistry) Rooms() map[string]*rooms.Room {
	return r.rooms
}

func (r *RoomRegistry) Delete(id string) {
	delete(r.rooms, id)
}
